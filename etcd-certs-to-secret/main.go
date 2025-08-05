package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"log/slog"
	"maps"
	"os"
	"strconv"
	"strings"

	qubershiporgv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1alpha1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

//go:embed  assets/*.yaml
var assets embed.FS

func main() {
	var namespace string
	var secretName string
	//var etcdPodLabel string
	var keyData, caData, crtData string
	// 0. Parsing flags
	flag.StringVar(&secretName, "secret", "kube-etcd-client-certs", "Name of the secret to create/update")
	//flag.StringVar(&etcdPodLabel, "label", "component=etcd", "Label selector for etcd pod")
	flag.Parse()
	namespace, found := os.LookupEnv("WATCH_NAMESPACE")
	if !found {
		namespace = "monitoring"
	}

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(log)

	namespacedName := types.NamespacedName{
		Namespace: namespace,
		Name:      "platformmonitoring",
	}

	scheme := runtime.NewScheme()

	if err := clientgoscheme.AddToScheme(scheme); err != nil {
		log.Error("Failed to add core Kubernetes types to scheme", "error", err)
		os.Exit(1)
	}
	if err := qubershiporgv1.AddToScheme(scheme); err != nil {
		log.Error("Failed to add PlatformMonitoring to scheme", "error", err)
		os.Exit(1)
	}
	if err := promv1.AddToScheme(scheme); err != nil {
		log.Error("Failed to add promv1 to scheme", "error", err)
		os.Exit(1)
	}

	cfg, err := config.GetConfig()
	if err != nil {
		log.Error("Couldn't get config", "error", err)
		os.Exit(1)
	}
	cl, err := client.New(cfg, client.Options{
		Scheme: scheme})
	if err != nil {
		log.Error("Couldn't get client", "error", err)
		os.Exit(1)
	}
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		log.Error("Couldn't get clientset", "error", err)
		os.Exit(1)
	}

	dc, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		log.Error("Couldn't get discovery client", "error", err)
		os.Exit(1)
	}

	isOpenshift, err := hasRouteApi(clientset)
	if err != nil {
		log.Error("Couldn't check if cluster has Route API", "error", err)
		os.Exit(1)
	}

	isOpenshiftV4, err := isOpenshiftV4(dc, isOpenshift, log)
	if err != nil {
		log.Error("Couldn't check cluster version", "error", err)
		os.Exit(1)
	}

	if isOpenshiftV4 {
		keyData, caData, crtData, err = getCertsFromConfigmapAndSecret(clientset, log, utils.EtcdCertificatesSourceNamespaceOpenshiftV4, utils.EtcdCertificatesSourceConfigmapOpenshiftV4, utils.EtcdCertificatesSourceSecretOpenshiftV4)
		if err != nil {
			if apierrors.IsForbidden(err) {
				log.Error("Unable to update etcd certificates due to a lack of permission to access the requested etcd resource.", "error", err)
			} else {
				log.Error("Failed to get etcd certificates from configmap and secret", "error", err)
			}
			os.Exit(1)
		}
		log.Info("Extracting etcd certificates from configmap and secret (Openshift v4)", "etcdNamespace", utils.EtcdCertificatesSourceNamespaceOpenshiftV4, "EtcdCertsSourceConfigmap", utils.EtcdCertificatesSourceConfigmapOpenshiftV4, "etcdCertsSourceSecret", utils.EtcdCertificatesSourceSecretOpenshiftV4)
	} else {
		keyData, caData, crtData, err = getCertsFromHostpath(clientset, log, utils.EtcdServiceComponentNamespace, utils.EtcdPodLabelSelector, isOpenshift)
		if err != nil {
			if apierrors.IsForbidden(err) {
				log.Error("Unable to update etcd certificates due to a lack of permission to access the requested etcd resource.", "error", err)
			} else {
				log.Error("Failed to get etcd certificates from hostpath", "error", err)
			}
			os.Exit(1)
		}
	}
	if err := certVerify(keyData, caData, crtData, log); err != nil {
		log.Error("Failed to verify etcd certificates", "error", err)
		os.Exit(1)
	}
	customResourceInstance := &qubershiporgv1.PlatformMonitoring{}
	if err := cl.Get(context.TODO(), namespacedName, customResourceInstance); err != nil {
		log.Error("Failed to get PlatformMonitoring custom resource", "error", err)
		os.Exit(1)
	}
	certData := make(map[string][]byte)
	certData["etcd-client-ca.crt"] = []byte(caData)
	certData["etcd-client.crt"] = []byte(crtData)
	certData["etcd-client.key"] = []byte(keyData)

	// Creating / updating etcd Secret
	secret := &corev1.Secret{}
	if secret, err = etcdSecret(customResourceInstance, &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
		},
		Data: certData,
	}); err != nil {
		log.Error("Failed creating Secret manifest", "error", err)
		os.Exit(1)
	}

	if err := createOrUpdateSecret(clientset, secret, log); err != nil {
		log.Error("Secret operation failed", "secret", secretName, "error", err)
		os.Exit(1)
	}

	if err := createOrUpdateServiceMonitor(customResourceInstance, cl, namespace, isOpenshiftV4, log); err != nil {
		log.Error("Failed to create/update etcd ServiceMonitor", "error", err)
		os.Exit(1)
	}

	if err := CreateOrUpdateService(customResourceInstance, clientset, isOpenshift, isOpenshiftV4, log); err != nil {
		log.Error("Failed to create/update etcd Service", "error", err)
		os.Exit(1)
	}
}

// Retrieve etcd certificates paths from etcd pods' command line arguments
func extractEtcdCertPaths(ctx context.Context, clientset *kubernetes.Clientset, log *slog.Logger, etcdNamespace string, etcdPodLabel string) (peerKey, caCrt, peerCrt string, err error) {
	// Default cert paths
	peerKey = "/etc/kubernetes/pki/etcd/peer.key"
	caCrt = "/etc/kubernetes/pki/etcd/ca.crt"
	peerCrt = "/etc/kubernetes/pki/etcd/peer.crt"

	pods, err := clientset.CoreV1().Pods(etcdNamespace).List(ctx, metav1.ListOptions{
		LabelSelector: etcdPodLabel,
	})

	if err != nil || len(pods.Items) == 0 {
		log.Error("Failed to retrieve pods to get etcd certificates", "error", err)
		return "", "", "", err
	}

	var podNames []string
	var podIndex int

	for i, p := range pods.Items {
		if p.Status.Phase == corev1.PodRunning {
			podNames = append(podNames, p.ObjectMeta.Name)
			podIndex = i
		}
	}
	if len(podNames) == 0 {
		log.Error("Failed to find etcd pods among pods to get etcd certificates", "error", err)
		return "", "", "", err
	}

	etcdPod := pods.Items[podIndex]

	for _, container := range etcdPod.Spec.Containers {
		if container.Name == "etcd" {
			for _, arg := range container.Command {
				if strings.HasPrefix(arg, "--peer-key-file=") {
					peerKey = strings.SplitN(arg, "=", 2)[1]
				} else if strings.HasPrefix(arg, "--peer-trusted-ca-file=") {
					caCrt = strings.SplitN(arg, "=", 2)[1]
				} else if strings.HasPrefix(arg, "--peer-cert-file=") {
					peerCrt = strings.SplitN(arg, "=", 2)[1]
				}
			}
		}
	}
	return peerKey, caCrt, peerCrt, nil
}

func hasRouteApi(clientset *kubernetes.Clientset) (bool, error) {
	// get all available API resources
	resources, err := clientset.Discovery().ServerPreferredResources()
	if err != nil {
		return false, fmt.Errorf("failed to get server api resources: %v", err)
	}
	for _, resourceList := range resources {
		if resourceList.GroupVersion == "route.openshift.io/v1" {
			return true, nil
		}
	}
	return false, nil
}

func isOpenshiftV4(dc discovery.DiscoveryInterface, isOpenshift bool, log *slog.Logger) (bool, error) {
	serverVersion, err := dc.ServerVersion()
	if err != nil {
		return false, fmt.Errorf("failed to get server version: %v", err)
	}
	log.Info("Server version", "minor", serverVersion.Minor)
	minor, err := strconv.Atoi(serverVersion.Minor)
	if err != nil {
		return false, fmt.Errorf("failed to convert minor server version %s to integer: %v", serverVersion.Minor, err)
	}
	return minor >= 18 && isOpenshift, nil
}

func getCertsFromConfigmapAndSecret(clientset *kubernetes.Clientset, log *slog.Logger, etcdNamespace string, configmapName string, etcdCertsSourceSecret string) (string, string, string, error) {
	configMap, err := clientset.CoreV1().ConfigMaps(etcdNamespace).Get(context.TODO(), configmapName, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsForbidden(err) {
			log.Error("Failed to get configmap due to insufficient permissions", "namespace", etcdNamespace, "configmap", configmapName, "error", err)
		} else {
			log.Error("Failed to get configmap", "namespace", etcdNamespace, "configmap", configmapName, "error", err)
		}
		return "", "", "", err
	}
	caData := configMap.Data["ca-bundle.crt"]

	secret, err := clientset.CoreV1().Secrets(etcdNamespace).Get(context.TODO(), etcdCertsSourceSecret, metav1.GetOptions{})
	if err != nil {
		if !apierrors.IsForbidden(err) {
			log.Error(fmt.Sprintf("Failed to get secret %s (namespace: %s) to get etcd certificates", etcdCertsSourceSecret, etcdNamespace), "error", err)
		}
		return "", "", "", err
	}

	secretData := secret.Data
	crtData := string(secretData["tls.crt"])
	keyData := string(secretData["tls.key"])
	return keyData, caData, crtData, nil
}

func getCertsFromHostpath(clientset *kubernetes.Clientset, log *slog.Logger, etcdNamespace string, etcdPodLabel string, isOpenshift bool) (keyData string, caData string, crtData string, err error) {
	var peerKey, caCrt, peerCrt string
	if isOpenshift {
		peerKey = "/etc/etcd/peer.key"
		caCrt = "/etc/etcd/ca.crt"
		peerCrt = "/etc/etcd/peer.crt"
		log.Info("Using default Openshift prior to v4 etcd certificates paths", "peerKey", peerKey, "caCrt", caCrt, "peerCrt", peerCrt)
	} else {
		peerKey, caCrt, peerCrt, err = extractEtcdCertPaths(context.TODO(), clientset, log, etcdNamespace, etcdPodLabel)
		if err != nil {
			log.Error("Failed to get etcd certificates paths from etcd pods arguments", "error", err)
			return "", "", "", err
		}
		log.Info("Using etcd certificates paths from etcd pods", "peerKey", peerKey, "caCrt", caCrt, "peerCrt", peerCrt)
	}

	caData, err = readFileToString(log, caCrt)
	if err != nil {
		return "", "", "", err
	}

	keyData, err = readFileToString(log, peerKey)
	if err != nil {
		return "", "", "", err
	}

	crtData, err = readFileToString(log, peerCrt)
	if err != nil {
		return "", "", "", err
	}

	return keyData, caData, crtData, nil
}

func readFileToString(log *slog.Logger, filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Error("Failed to read file", "file", filename, "error", err)
		return "", err
	}
	return string(data), nil
}

func etcdServiceMonitor(cr *qubershiporgv1.PlatformMonitoring, namespace string, isOpenshiftV4 bool) (*promv1.ServiceMonitor, error) {
	sm := promv1.ServiceMonitor{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.EtcdServiceMonitorAsset), 100).Decode(&sm); err != nil {
		return nil, err
	}

	//Set parameters
	sm.SetGroupVersionKind(schema.GroupVersionKind{Group: "monitoring.coreos.com", Version: "v1", Kind: "ServiceMonitor"})
	sm.SetName(namespace + "-" + "etcd-service-monitor")
	sm.SetNamespace(namespace)
	if isOpenshiftV4 {
		sm.Spec.NamespaceSelector.MatchNames = []string{utils.EtcdServiceComponentNamespaceOpenshiftV4}
		// Port "etcd-metrics" is used in OpenShift v4.x
		sm.Spec.Endpoints[0].Port = "etcd-metrics"
	} else {
		sm.Spec.NamespaceSelector.MatchNames = []string{utils.EtcdServiceComponentNamespace}
	}

	if cr.Spec.KubernetesMonitors != nil {
		monitor, ok := cr.Spec.KubernetesMonitors[utils.EtcdServiceMonitorName]
		if ok && monitor.IsInstall() {
			monitor.OverrideServiceMonitor(&sm)
		}
	}

	if cr.GetLabels() != nil {
		maps.Copy(sm.Labels, cr.GetLabels())
	}
	// Set labels
	sm.Labels["name"] = utils.TruncLabel(sm.GetName())
	sm.Labels["app.kubernetes.io/name"] = utils.TruncLabel(sm.GetName())
	sm.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(sm.GetName(), sm.GetNamespace())
	sm.Labels["app.kubernetes.io/component"] = "monitoring-etcd"

	if sm.Annotations == nil && cr.GetAnnotations() != nil {
		sm.SetAnnotations(cr.GetAnnotations())
	} else {
		maps.Copy(sm.Annotations, cr.GetAnnotations())
	}
	sm.OwnerReferences = []metav1.OwnerReference{
		{
			APIVersion:         "monitoring.qubership.org/v1alpha1",
			Kind:               "PlatformMonitoring",
			Name:               cr.Name,
			UID:                cr.UID,
			Controller:         ptr.To(true),
			BlockOwnerDeletion: ptr.To(true),
		},
	}
	return &sm, nil
}

func createOrUpdateServiceMonitor(cr *qubershiporgv1.PlatformMonitoring, cl client.Client, namespace string, isOpenshiftV4 bool, log *slog.Logger) error {
	sm, err := etcdServiceMonitor(cr, namespace, isOpenshiftV4)
	if err != nil {
		log.Error("Failed creating ServiceMonitor manifest", "error", err)
	}
	// Creating / updating etcd ServiceMonitor
	existingSM := &promv1.ServiceMonitor{}
	err = cl.Get(context.TODO(), client.ObjectKey{Name: sm.Name, Namespace: sm.Namespace}, existingSM)
	if err != nil {
		if apierrors.IsNotFound(err) {
			// Create new ServiceMonitor
			if err := cl.Create(context.TODO(), sm); err != nil {
				return fmt.Errorf("failed to create etcd ServiceMonitor: %w", err)
			}
			log.Info("ServiceMonitor created", "name", sm.Name)
		} else {
			return fmt.Errorf("failed to check etcd servicemonitor existence: %w", err)
		}
	} else {
		// Update existing ServiceMonitor
		sm.SetResourceVersion(existingSM.GetResourceVersion())
		if err := cl.Update(context.TODO(), sm); err != nil {
			return fmt.Errorf("failed to update etcd ServiceMonitor: %w", err)
		}
		log.Info("ServiceMonitor updated", "name", sm.Name)
	}
	return nil
}

func CreateOrUpdateService(cr *qubershiporgv1.PlatformMonitoring, clientset *kubernetes.Clientset, isOpenshift bool, isOpenshiftV4 bool, log *slog.Logger) error {
	etcdServiceNamespace := utils.EtcdServiceComponentNamespace
	if isOpenshiftV4 {
		etcdServiceNamespace = utils.EtcdServiceComponentNamespaceOpenshiftV4
	}

	m, err := etcdService(isOpenshift, etcdServiceNamespace, isOpenshiftV4)
	if err != nil {
		log.Error("Failed creating Service manifest", "error", err)
		return err
	}

	e, err := clientset.CoreV1().Services(etcdServiceNamespace).Get(context.TODO(), m.Name, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			e = &corev1.Service{}
			e.Name = m.Name
			e.Namespace = etcdServiceNamespace
		} else {
			return fmt.Errorf("failed to get check if etcd service exists: %w", err)
		}
	}
	//Set parameters
	e.TypeMeta = m.TypeMeta
	e.SetLabels(m.GetLabels())
	e.Spec.Ports = m.Spec.Ports
	e.Spec.Selector = m.Spec.Selector
	e.Spec.ClusterIP = m.Spec.ClusterIP
	if cr.GetLabels() != nil {
		maps.Copy(e.Labels, cr.GetLabels())
	}
	if e.Labels == nil {
		e.Labels = make(map[string]string)
	}
	if e.Annotations == nil {
		e.Annotations = make(map[string]string)
	}
	if e.Annotations == nil && cr.GetAnnotations() != nil {
		maps.Copy(e.Annotations, cr.GetAnnotations())
	}
	if apierrors.IsNotFound(err) {
		if _, err = clientset.CoreV1().Services(etcdServiceNamespace).Create(context.TODO(), e, metav1.CreateOptions{}); err != nil {
			log.Error("Failed to create etcd service", "error", err, "service", m.Name)
			return err
		} else {
			log.Info("Service created", "service", m.Name)
			return nil
		}
	} else {
		if _, err = clientset.CoreV1().Services(etcdServiceNamespace).Update(context.TODO(), e, metav1.UpdateOptions{}); err != nil {
			if apierrors.IsForbidden(err) {
				log.Info("Failed to update etcd service", "error", err, "service", m.Name)
				return err
			}
			log.Error("Failed to update etcd service", "error", err, "service", m.Name)
			return fmt.Errorf("failed to update etcd service: %w", err)
		}
		log.Info("Service updated", "service", m.Name)
		return nil
	}
}
func etcdService(isOpenshift bool, etcdServiceNamespace string, isOpenshiftV4 bool) (*corev1.Service, error) {
	service := corev1.Service{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.EtcdServiceComponentAsset), 100).Decode(&service); err != nil {
		return nil, err
	}
	//Set parameters
	service.SetGroupVersionKind(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Service"})
	service.SetName(utils.EtcdServiceComponentName)
	service.SetNamespace(etcdServiceNamespace)

	// Kubernetes uses "component: etcd" selector
	// OpenShift v3.x uses "openshift.io/component: etcd" selector
	if isOpenshift && !isOpenshiftV4 {
		service.Spec.Selector = map[string]string{"openshift.io/component": "etcd"}
	}
	// OpenShift v4.x uses "etcd: 'true'" selector
	if isOpenshiftV4 {
		service.Spec.Selector = map[string]string{"etcd": "true"}
	}
	// If cluster is not OpenShift v4.x, remove port "etcd-metrics"
	if !isOpenshiftV4 {
		service.Spec.Ports = service.Spec.Ports[:1]
	}
	service.Spec.ClusterIP = ""
	// Set labels
	service.Labels["name"] = utils.TruncLabel(service.GetName())
	service.Labels["app.kubernetes.io/name"] = utils.TruncLabel(service.GetName())
	service.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(service.GetName(), service.GetNamespace())

	return &service, nil
}

func createOrUpdateSecret(clientset *kubernetes.Clientset, secret *corev1.Secret, log *slog.Logger) error {
	_, err := clientset.CoreV1().Secrets(secret.Namespace).Get(context.TODO(), secret.Name, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsForbidden(err) {
			return fmt.Errorf("failed to check secret existence due to insufficient permissions: %w", err)
		}
		if apierrors.IsNotFound(err) {

			// Create
			_, err = clientset.CoreV1().Secrets(secret.Namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
			if err != nil {
				return fmt.Errorf("failed to create secret: %w", err)
			}
			log.Info("Secret created", "secret", secret.Name)
			return nil
		}
		return fmt.Errorf("failed to check secret existence: %w", err)
	}

	// Update
	_, err = clientset.CoreV1().Secrets(secret.Namespace).Update(context.TODO(), secret, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update secret: %w", err)
	}
	log.Info("Secret updated", "secret", secret.Name)
	return nil
}

func certVerify(keyData string, caData string, crtData string, log *slog.Logger) error {
	if len(keyData) == 0 || len(caData) == 0 || len(crtData) == 0 {
		return fmt.Errorf("failed to get etcd certificates content, empty certificate data")
	}

	peerKeyBeginIndex := strings.Index(keyData, "-----BEGIN PRIVATE KEY-----")
	if peerKeyBeginIndex == -1 {
		peerKeyBeginIndex = strings.Index(keyData, "-----BEGIN RSA PRIVATE KEY-----")
	}

	caCertBeginIndex := strings.Index(caData, "-----BEGIN CERTIFICATE-----")
	caCertEndIndex := strings.Index(caData, "-----END CERTIFICATE-----")

	peerCertBeginIndex := strings.Index(crtData, "-----BEGIN CERTIFICATE-----")
	peerCertEndIndex := strings.Index(crtData, "-----END CERTIFICATE-----")

	if peerKeyBeginIndex == -1 {
		return fmt.Errorf("failed to find private key header")
	}
	if caCertBeginIndex == -1 || caCertEndIndex == -1 || caCertBeginIndex > caCertEndIndex {
		return fmt.Errorf("invalid CA certificate format")
	}
	if peerCertBeginIndex == -1 || peerCertEndIndex == -1 || peerCertBeginIndex > peerCertEndIndex {
		log.Error("invalid peer certificate format")
		os.Exit(1)
		return fmt.Errorf("invalid peer certificate format")
	}
	return nil
}
func etcdSecret(cr *qubershiporgv1.PlatformMonitoring, secret *corev1.Secret) (*corev1.Secret, error) {
	secret.Labels = make(map[string]string)
	secret.Annotations = make(map[string]string)
	//Set parameters
	if cr.GetLabels() != nil {
		for k, v := range cr.GetLabels() {
			if _, ok := secret.Labels[k]; !ok {
				secret.Labels[k] = v
			}
		}
	}
	secret.Labels["name"] = secret.Name
	secret.Labels["app.kubernetes.io/name"] = utils.TruncLabel(secret.Name)
	secret.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(secret.Name, secret.Namespace)
	secret.Labels["app.kubernetes.io/component"] = "monitoring-etcd"
	if secret.Annotations == nil && cr.GetAnnotations() != nil {
		secret.SetAnnotations(cr.GetAnnotations())
	} else {
		for k, v := range cr.GetAnnotations() {
			secret.Annotations[k] = v
		}
	}
	secret.OwnerReferences = []metav1.OwnerReference{
		{
			APIVersion:         "monitoring.qubership.org/v1alpha1",
			Kind:               "PlatformMonitoring",
			Name:               cr.Name,
			UID:                cr.UID,
			Controller:         ptr.To(true),
			BlockOwnerDeletion: ptr.To(true),
		},
	}
	return secret, nil
}
