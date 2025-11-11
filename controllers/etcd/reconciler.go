package etcd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	v1alpha1 "github.com/Netcracker/qubership-monitoring-operator/api/v1alpha1"
	kubernetesmonitors "github.com/Netcracker/qubership-monitoring-operator/controllers/kubernetes-monitors"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type EtcdMonitorReconciler struct {
	KubeClient kubernetes.Interface
	config     *rest.Config
	*utils.ComponentReconciler
}

func NewEtcdMonitorReconciler(c client.Client, s *runtime.Scheme, dc discovery.DiscoveryInterface, r *rest.Config) *EtcdMonitorReconciler {
	clientSet, err := kubernetes.NewForConfig(r)

	if err != nil {
		panic(err.Error())
	}

	return &EtcdMonitorReconciler{
		ComponentReconciler: &utils.ComponentReconciler{
			Client: c,
			Scheme: s,
			Dc:     dc,
			Log:    utils.Logger("etcd_monitor_reconciler"),
		},
		KubeClient: clientSet,
		config:     r,
	}
}

// Run reconciliation for etcd-monitor configuration.
// Creates new service monitor and secret with certificates if its don't exists.
// Updates monitor and secret in case of any changes.
// Returns true if need to requeue, false otherwise.
func (r *EtcdMonitorReconciler) Run(ctx context.Context, cr *v1alpha1.PlatformMonitoring) error {
	r.Log.Info("Reconciling component")

	// Try to get Route is there to check is it OpenShift or Kubernetes
	isOpenshift := r.HasRouteApi()

	isOpenshiftV4, err := r.IsOpenShiftV4()
	if err != nil {
		r.Log.Error(err, "Failed to recognise OpenShift V4, continue for OpenShift v3.11")
	}

	var etcdServiceNamespace = utils.EtcdServiceComponentNamespace
	// Etcd pods in OpenShift v4.x have a different namespace
	if isOpenshiftV4 {
		etcdServiceNamespace = utils.EtcdServiceComponentNamespaceOpenshiftV4
	}

	// If monitor should not be installed in public cloud, the rest of this code will not work too
	affected, installed := kubernetesmonitors.IsMonitorPresentInPublicCloud(cr, utils.EtcdServiceMonitorName)
	if affected && !installed {
		r.uninstall(cr, isOpenshift, etcdServiceNamespace, isOpenshiftV4)
		r.Log.Info("Skip Secret resource handling and certificates updating for public cloud")
		return nil
	}

	if len(cr.Spec.KubernetesMonitors) == 0 || !kubernetesmonitors.IsMonitorInstall(cr, utils.EtcdServiceMonitorName) {
		r.Log.Info("Uninstalling component if exists")
		r.uninstall(cr, isOpenshift, etcdServiceNamespace, isOpenshiftV4)
		r.Log.Info("Component reconciled")
		return nil
	}

	// Create Secret
	if err := r.handleSecret(cr); err != nil {
		return err
	}

	// Get minor server version to recognise Kubernetes v1.19 or higher to get etcd certs correctly
	minorServerVersion, err := r.GetMinorServerVersion()
	if err != nil {
		r.Log.Error(err, "Failed to get minor server version")
		return err
	}

	// If monitor should not be installed in public cloud, the rest of this code will not work too
	// Update certificate in the Secret
	if err = r.updateCertificates(ctx, cr, isOpenshift, minorServerVersion, isOpenshiftV4); err != nil {
		if apierrors.IsForbidden(err) && !utils.PrivilegedRights {
			slog.New(logr.ToSlogHandler(r.Log)).Warn(err.Error(), "details", "Unable to update etcd certificates due to a lack of permission to access the requested etcd resource.")
		} else {
			r.Log.Error(err, "Failed to update etcd certificates")
		}
	}

	if err = r.handleServiceAccount(cr); err != nil {
		return err
	}
	if err = r.handleServiceMonitor(cr, etcdServiceNamespace, isOpenshiftV4); err != nil {
		return err
	}
	if utils.PrivilegedRights {
		if err = r.handleService(cr, isOpenshift, etcdServiceNamespace, isOpenshiftV4); err != nil {
			return err
		}
	} else {
		r.Log.Info("Skip Service resource reconciliation because privilegedRights=false")
	}

	r.Log.Info("Component reconciled")
	return nil
}

// uninstall deletes all resources related to the component
func (r *EtcdMonitorReconciler) uninstall(cr *v1alpha1.PlatformMonitoring, isOpenshift bool, etcdServiceNamespace string, isOpenshiftV4 bool) {
	if err := r.deleteServiceMonitor(cr, etcdServiceNamespace, isOpenshiftV4); err != nil {
		r.Log.Error(err, "Can not delete ServiceMonitor")
	}
	r.Log.Info("ServiceMonitor for ETCD successful deleted.")
}

// In any Kubernetes clusters or OpenShift version under 4.5 (OpenShift with version of Kubernetes under 1.18)
// certificates can be gotten from etcd pods in kube-system namespace
func (r *EtcdMonitorReconciler) getCertsFromEtcdPods(ctx context.Context, isOpenshift bool, minorServerVersion int) (string, string, string, error) {
	pods, err := r.KubeClient.CoreV1().Pods(utils.EtcdServiceComponentNamespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: utils.EtcdPodLabelSelector,
	})
	if err != nil {
		if !apierrors.IsForbidden(err) || utils.PrivilegedRights {
			r.Log.Error(err, "Failed to retrieve pods to get etcd certificates")
		}
		return "", "", "", err
	}
	var podNames []string
	var podIndex int
	if len(pods.Items) != 0 {
		for i, p := range pods.Items {
			if p.Status.Phase == corev1.PodRunning {
				podNames = append(podNames, p.ObjectMeta.Name)
				podIndex = i
			}
		}
	}
	if len(podNames) == 0 {
		return "", "", "", errors.New("failed to find etcd pods among pods to get etcd certificates")
	}

	var pathsToCerts []string
	// Certificates in Kubernetes and OpenShift have different paths
	if isOpenshift {
		pathsToCerts = []string{
			"/etc/etcd/peer.key",
			"/etc/etcd/ca.crt",
			"/etc/etcd/peer.crt",
		}
	} else {
		// Paths in Kubernetes can be changed by arguments
		peerKey := "/etc/kubernetes/pki/etcd/peer.key"
		caCrt := "/etc/kubernetes/pki/etcd/ca.crt"
		peerCrt := "/etc/kubernetes/pki/etcd/peer.crt"
		etcdPod := pods.Items[podIndex]
		for _, container := range etcdPod.Spec.Containers {
			if container.Name == utils.EtcdServiceComponentName {
				commands := container.Command
				for _, command := range commands {
					if strings.HasPrefix(command, "--peer-key-file") {
						peerKey = strings.Split(command, "=")[1]
					} else if strings.HasPrefix(command, "--peer-trusted-ca-file") {
						caCrt = strings.Split(command, "=")[1]
					} else if strings.HasPrefix(command, "--peer-cert-file") {
						peerCrt = strings.Split(command, "=")[1]
					}
				}
			}
		}
		pathsToCerts = []string{
			peerKey,
			caCrt,
			peerCrt,
		}
	}

	// If it's Kubernetes cluster v1.19 or higher, printf "%s\n" is used
	// echo $(<path/to/peer.key); echo $(<path/to/ca.crt); echo $(<path/to/peer.crt)
	command := fmt.Sprintf("echo \"$(<%s)\"", strings.Join(pathsToCerts, ")\";echo \"$(<"))

	if isOpenshift || minorServerVersion < 19 {
		// Otherwise, cat is used
		// cat path/to/peer.key path/to/ca.crt path/to/peer.crt
		command = fmt.Sprintf("cat %s", strings.Join(pathsToCerts, " "))
	}

	cmd := []string{
		"/bin/sh",
		"-c",
		command,
	}

	buf := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}
	for _, podName := range podNames {
		req := r.KubeClient.CoreV1().RESTClient().Post().Resource("pods").Name(podName).Namespace(utils.EtcdServiceComponentNamespace).SubResource("exec").Param("container", utils.EtcdServiceComponentName)
		req.VersionedParams(&corev1.PodExecOptions{
			Container: utils.EtcdServiceComponentName,
			Command:   cmd,
			Stderr:    true,
			Stdout:    true,
		}, scheme.ParameterCodec)

		exec, err := remotecommand.NewSPDYExecutor(r.config, http.MethodPost, req.URL())
		if err != nil {
			r.Log.Error(err, "Failed to execute command in etcd pod")
			return "", "", "", err
		}

		streamOptions := remotecommand.StreamOptions{
			Stdout: buf,
			Stderr: errBuf,
			Tty:    false,
		}

		err = exec.StreamWithContext(ctx, streamOptions)
		if err != nil {
			r.Log.Error(err, fmt.Sprintf("Error during stream getting in %v. Try another pod..", podName))
		} else {
			r.Log.Info(fmt.Sprintf("Successfully getting certificates for etcd from pods in %s namespace.", utils.EtcdServiceComponentNamespace))
			break
		}
	}

	certs := buf.String()
	if certs == "" {
		return "", "", "", errors.New("failed to execute command in etcd pod")
	}

	peerKeyBeginIndex := strings.LastIndex(certs, "-----BEGIN RSA PRIVATE KEY-----")
	if peerKeyBeginIndex == -1 {
		peerKeyBeginIndex = strings.LastIndex(certs, "-----BEGIN PRIVATE KEY-----")
	}
	caCertBeginIndex := strings.Index(certs, "-----BEGIN CERTIFICATE-----")
	caCertEndIndex := strings.Index(certs, "-----END CERTIFICATE-----")
	peerCertBeginIndex := strings.LastIndex(certs, "-----BEGIN CERTIFICATE-----")
	peerCertEndIndex := strings.LastIndex(certs, "-----END CERTIFICATE-----")

	if peerKeyBeginIndex == -1 || caCertBeginIndex == -1 || caCertEndIndex == -1 ||
		peerCertBeginIndex == -1 || peerCertEndIndex == -1 || caCertBeginIndex == peerCertBeginIndex {
		return "", "", "", errors.New("failed to get certificate: bad certificates content")
	}

	keyFile := certs[peerKeyBeginIndex:caCertBeginIndex]
	caFile := certs[caCertBeginIndex:caCertEndIndex] + "-----END CERTIFICATE-----"
	certFile := certs[peerCertBeginIndex:peerCertEndIndex] + "-----END CERTIFICATE-----"

	return caFile, certFile, keyFile, nil
}

func (r *EtcdMonitorReconciler) getCertsFromConfigmapAndSecret() (string, string, string, error) {
	configMap, err := r.KubeClient.CoreV1().ConfigMaps(utils.EtcdCertificatesSourceNamespaceOpenshiftV4).Get(context.TODO(), utils.EtcdCertificatesSourceConfigmapOpenshiftV4, metav1.GetOptions{})
	if err != nil {
		if !apierrors.IsForbidden(err) || utils.PrivilegedRights {
			r.Log.Error(err, fmt.Sprintf("Failed to get configmap %s (namespace: %s) to get etcd certificates", utils.EtcdCertificatesSourceConfigmapOpenshiftV4, utils.EtcdCertificatesSourceNamespaceOpenshiftV4))
		}
		return "", "", "", err
	}
	caFile := configMap.Data["ca-bundle.crt"]

	secret, err := r.KubeClient.CoreV1().Secrets(utils.EtcdCertificatesSourceNamespaceOpenshiftV4).Get(context.TODO(), utils.EtcdCertificatesSourceSecretOpenshiftV4, metav1.GetOptions{})
	if err != nil {
		if !apierrors.IsForbidden(err) || utils.PrivilegedRights {
			r.Log.Error(err, fmt.Sprintf("Failed to get secret %s (namespace: %s) to get etcd certificates", utils.EtcdCertificatesSourceSecretOpenshiftV4, utils.EtcdCertificatesSourceNamespaceOpenshiftV4))
		}
		return "", "", "", err
	}

	secretData := secret.Data
	certFile := string(secretData["tls.crt"])
	keyFile := string(secretData["tls.key"])

	r.Log.Info(fmt.Sprintf("Successfully getting certificates for etcd from configmap and secret in %s namespace.", utils.EtcdCertificatesSourceNamespaceOpenshiftV4))
	return caFile, certFile, keyFile, nil
}

func (r *EtcdMonitorReconciler) updateCertificates(ctx context.Context, cr *v1alpha1.PlatformMonitoring, isOpenshift bool, minorServerVersion int, isOpenshiftV4 bool) error {
	var caFile, certFile, keyFile string
	var err error

	// In OpenShift v4.x we get certificates from configmap and secret instead of pods
	if isOpenshiftV4 {
		caFile, certFile, keyFile, err = r.getCertsFromConfigmapAndSecret()
	} else {
		caFile, certFile, keyFile, err = r.getCertsFromEtcdPods(ctx, isOpenshift, minorServerVersion)
	}

	if caFile == "" || certFile == "" || keyFile == "" || err != nil {
		if !apierrors.IsForbidden(err) || utils.PrivilegedRights {
			r.Log.Error(err, "Failed to get etcd certificates from cluster")
		}
		return err
	}

	secret, err := r.KubeClient.CoreV1().Secrets(cr.GetNamespace()).Get(context.TODO(), utils.KubeEtcdClientCertsSecretName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	certData := make(map[string][]byte)

	certData["etcd-client-ca.crt"] = []byte(caFile)
	certData["etcd-client.crt"] = []byte(certFile)
	certData["etcd-client.key"] = []byte(keyFile)

	secret.Data = certData
	_, err = r.KubeClient.CoreV1().Secrets(cr.GetNamespace()).Update(context.TODO(), secret, metav1.UpdateOptions{})

	if err != nil {
		return err
	}
	return nil
}
