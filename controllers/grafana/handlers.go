package grafana

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"time"

	v1alpha1 "github.com/Netcracker/qubership-monitoring-operator/api/v1alpha1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	appsv1 "k8s.io/api/apps/v1"
	grafv1 "github.com/grafana-operator/grafana-operator/v4/api/integreatly/v1alpha1"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

func (r *GrafanaReconciler) handleGrafana(cr *v1alpha1.PlatformMonitoring) error {
	m, err := grafana(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Grafana manifest")
		return err
	}

	if m.Spec.Config.AuthGenericOauth != nil {
		secret, err := r.KubeClient.CoreV1().Secrets(m.Namespace).Get(context.TODO(), utils.GrafanaExtraVarsSecret, metav1.GetOptions{})
		if err != nil {
			r.Log.Error(err, "auth.generic_oauth is configured but clientId and clientSecret is not stored in secret")
			return err
		}
		if len(secret.Data["GF_AUTH_GENERIC_OAUTH_CLIENT_ID"]) == 0 || len(secret.Data["GF_AUTH_GENERIC_OAUTH_CLIENT_SECRET"]) == 0 {
			r.Log.Error(err, "auth.generic_oauth is configured but clientId and/or clientSecret is not stored in secret")
			return err
		}
		m.Spec.Config.AuthGenericOauth.ClientId = secret.StringData["GF_AUTH_GENERIC_OAUTH_CLIENT_ID"]
		m.Spec.Config.AuthGenericOauth.ClientSecret = secret.StringData["GF_AUTH_GENERIC_OAUTH_CLIENT_SECRET"]
	}
	e := &grafv1.Grafana{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			if err = r.CreateResource(cr, m); err != nil {
				return err
			}
			return nil
		}
		return err
	}

	//Set parameters
	e.SetLabels(m.GetLabels())
	e.Spec = m.Spec

	if err = r.UpdateResource(e); err != nil {
		return err
	}
	// WA for https://github.com/grafana-operator/grafana-operator/issues/652
	r.Log.Info("Waiting grafana-deployment")
	time.Sleep(30 * time.Second)
	return nil
}

func (r *GrafanaReconciler) handleGrafanaDataSource(cr *v1alpha1.PlatformMonitoring) error {
	jaegerServices, err := r.getJaegerServices(cr)
	if err != nil {
		r.Log.Error(err, "Failed getting Jaeger services")
	}
	clickHouseServices, err := r.getClickhouseServices(cr)
	if err != nil {
		r.Log.Error(err, "Failed getting ClickHouse services")
	}
	m, err := grafanaDataSource(cr, r.KubeClient, jaegerServices, clickHouseServices)
	if err != nil {
		r.Log.Error(err, "Failed creating GrafanaDataSource manifest")
		return err
	}

	// Set labels
	m.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(m.GetName(), m.GetNamespace())
	m.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Grafana.Image)

	e := &grafv1.GrafanaDataSource{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			if err = r.CreateResource(cr, m); err != nil {
				return err
			}
			return nil
		}
		return err
	}

	//Set parameters
	e.SetLabels(m.GetLabels())
	e.Spec = m.Spec

	if err = r.UpdateResource(e); err != nil {
		return err
	}
	return nil
}

func (r *GrafanaReconciler) handleIngressV1beta1(cr *v1alpha1.PlatformMonitoring) error {
	m, err := grafanaIngressV1beta1(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Ingress manifest")
		return err
	}
	e := &v1beta1.Ingress{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			if err = r.CreateResource(cr, m); err != nil {
				return err
			}
			return nil
		}
		return err
	}

	//Set parameters
	e.SetLabels(m.GetLabels())
	e.SetAnnotations(m.GetAnnotations())
	e.Spec.Rules = m.Spec.Rules
	e.Spec.TLS = m.Spec.TLS

	if err = r.UpdateResource(e); err != nil {
		return err
	}
	return nil
}

func (r *GrafanaReconciler) handleIngressV1(cr *v1alpha1.PlatformMonitoring) error {
	m, err := grafanaIngressV1(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Ingress manifest")
		return err
	}
	e := &networkingv1.Ingress{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			if err = r.CreateResource(cr, m); err != nil {
				return err
			}
			return nil
		}
		return err
	}

	//Set parameters
	e.SetLabels(m.GetLabels())
	e.SetAnnotations(m.GetAnnotations())
	e.Spec.Rules = m.Spec.Rules
	e.Spec.TLS = m.Spec.TLS

	if err = r.UpdateResource(e); err != nil {
		return err
	}
	return nil
}

func (r *GrafanaReconciler) handlePodMonitor(cr *v1alpha1.PlatformMonitoring) error {
	m, err := grafanaPodMonitor(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating PodMonitor manifest")
		return err
	}

	// Set labels
	m.Labels["name"] = utils.TruncLabel(m.GetName())
	m.Labels["app.kubernetes.io/name"] = utils.TruncLabel(m.GetName())
	m.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(m.GetName(), m.GetNamespace())
	m.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Grafana.Image)

	e := &promv1.PodMonitor{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			if err = r.CreateResource(cr, m); err != nil {
				return err
			}
			return nil
		}
		return err
	}

	//Set parameters
	e.SetLabels(m.GetLabels())
	e.Spec.JobLabel = m.Spec.JobLabel
	e.Spec.PodMetricsEndpoints = m.Spec.PodMetricsEndpoints
	e.Spec.NamespaceSelector = m.Spec.NamespaceSelector
	e.Spec.Selector = m.Spec.Selector

	if err = r.UpdateResource(e); err != nil {
		return err
	}
	return nil
}

func (r *GrafanaReconciler) handleGrafanaCredentialsSecret(cr *v1alpha1.PlatformMonitoring) (err error) {
	e := &corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: "grafana-admin-credentials", Namespace: cr.GetNamespace()}}
	tmpSecret := &corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: "grafana-admin-credentials-temp", Namespace: cr.GetNamespace()}}

	// Set labels
	e.SetLabels(map[string]string{
		"name":                         utils.TruncLabel(e.GetName()),
		"app.kubernetes.io/name":       utils.TruncLabel(e.GetName()),
		"app.kubernetes.io/managed-by": "monitoring-operator",
		"app.kubernetes.io/part-of":    "monitoring",
		"app.kubernetes.io/instance":   utils.GetInstanceLabel(e.GetName(), e.GetNamespace()),
		"app.kubernetes.io/version":    utils.GetTagFromImage(cr.Spec.Grafana.Image),
	})

	if err = r.GetResource(e); err == nil {
		if err = r.GetResource(tmpSecret); err == nil {
			if !reflect.DeepEqual(e.Data, tmpSecret.Data) {
				e.Data = tmpSecret.Data
				err = r.UpdateResource(e)
				if err == nil {
					isSecretUpdated = true
				}
			}
		}
	}
	if errors.IsNotFound(err) {
		return nil
	}
	return
}

func (r *GrafanaReconciler) resetGrafanaCredentials(cr *v1alpha1.PlatformMonitoring) (err error) {
	// Waiting Grafana Pods readiness
	r.Log.Info("Waiting for Grafana pods statuses", "kind", "Deployment", "name", utils.GrafanaDeploymentName)
	if err := r.WaitForPodsReadiness(
		&appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      utils.GrafanaDeploymentName,
				Namespace: cr.GetNamespace(),
			}}); err != nil {
		return err
	}
	r.Log.Info("Grafana Pods are ready", "kind", "Deployment", "name", utils.GrafanaDeploymentName)		
	// Getting Temp Secret
	r.Log.Info("Getting Temp Secret")
	tmpSecret := &corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: "grafana-admin-credentials-temp", Namespace: cr.GetNamespace()}}
	if err = r.GetResource(tmpSecret); err == nil {
		// Get Grafana Pod
		r.Log.Info("Getting Grafana Pod")
		config, err := rest.InClusterConfig()
		if err != nil {
			return fmt.Errorf("cannot load in-cluster config: %w", err)
		}
		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			return fmt.Errorf("cannot create clientset: %w", err)
		}
		pods, err := clientset.CoreV1().Pods(cr.GetNamespace()).List(context.TODO(), metav1.ListOptions{
			LabelSelector: "app=grafana",
		})
		if err != nil || len(pods.Items) == 0 {
			return fmt.Errorf("grafana deployment pod wasn't found: %w", err)
		}
		var podName *string = nil
		for _, p := range pods.Items {
			if p.DeletionTimestamp == nil {
				podName = &p.Name
				break
			}
		}
		if podName == nil {
			return fmt.Errorf("no suitable grafana deployment pod was found: %w", err)
		}
		r.Log.Info("Grafana Pod was found: " + *podName)

		// Prepare Grafana CLI request
		command := []string{"grafana", "cli", "admin", "reset-admin-password", string(tmpSecret.Data["GF_SECURITY_ADMIN_PASSWORD"])}
		req := r.KubeClient.CoreV1().RESTClient().
			Post().
			Resource("pods").
			Name(*podName).
			Namespace(cr.GetNamespace()).
			SubResource("exec").
			VersionedParams(&corev1.PodExecOptions{
				Container: "grafana",
				Command:   command,
				Stdin:     false,
				Stdout:    true,
				Stderr:    true,
				TTY:       false,
			}, scheme.ParameterCodec)

		// Set up a connection
		r.Log.Info("Setting Up a Connection with Grafana Pod")
		exec, err := remotecommand.NewSPDYExecutor(r.config, "POST", req.URL())
		if err != nil {
			return fmt.Errorf("grafana pod connection wasn't set up: %w", err)
		}

		// Execute Grafana CLI request
		r.Log.Info("Executing Grafana CLI command")
		var stdout, stderr bytes.Buffer
		err = exec.StreamWithContext(context.TODO(), remotecommand.StreamOptions{
			Stdout: &stdout,
			Stderr: &stderr,
		})
		if err != nil {
			return fmt.Errorf("error: %v; stdout: %s; stderr: %s;", err, stdout.String(), stderr.String())
		}

		isSecretUpdated = false
		r.Log.Info("Grafana Credentials Reset was finished")
	}
	if errors.IsNotFound(err) {
		r.Log.Info("Temp Secret wasn't found")
		return nil
	}
	return err
}

func (r *GrafanaReconciler) deleteGrafana(cr *v1alpha1.PlatformMonitoring) error {
	m, err := grafana(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Grafana manifest")
		return err
	}
	e := &grafv1.Grafana{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if err = r.DeleteResource(e); err != nil {
		return err
	}
	return nil
}

func (r *GrafanaReconciler) deleteGrafanaDataSource(cr *v1alpha1.PlatformMonitoring) error {
	jaegerServices, err := r.getJaegerServices(cr)
	if err != nil {
		r.Log.Error(err, "Failed getting Jaeger services")
	}
	clickHouseServices, err := r.getClickhouseServices(cr)
	if err != nil {
		r.Log.Error(err, "Failed getting ClickHouse services")
	}
	m, err := grafanaDataSource(cr, r.KubeClient, jaegerServices, clickHouseServices)
	if err != nil {
		r.Log.Error(err, "Failed creating GrafanaDataSource manifest")
		return err
	}
	e := &grafv1.GrafanaDataSource{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if err = r.DeleteResource(e); err != nil {
		return err
	}
	return nil
}

func (r *GrafanaReconciler) deleteIngressV1beta1(cr *v1alpha1.PlatformMonitoring) error {
	m, err := grafanaIngressV1beta1(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Ingress manifest")
		return err
	}
	e := &v1beta1.Ingress{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if err = r.DeleteResource(e); err != nil {
		return err
	}
	return nil
}

func (r *GrafanaReconciler) deleteIngressV1(cr *v1alpha1.PlatformMonitoring) error {
	m, err := grafanaIngressV1(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Ingress manifest")
		return err
	}
	e := &networkingv1.Ingress{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if err = r.DeleteResource(e); err != nil {
		return err
	}
	return nil
}

func (r *GrafanaReconciler) deletePodMonitor(cr *v1alpha1.PlatformMonitoring) error {
	m, err := grafanaPodMonitor(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating PodMonitor manifest")
		return err
	}
	e := &promv1.PodMonitor{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if err = r.DeleteResource(e); err != nil {
		return err
	}
	return nil
}

// Looking for Jaeger Services in all namespaces except current using a label selector and return list of them or nil
func (r *GrafanaReconciler) getJaegerServices(cr *v1alpha1.PlatformMonitoring) ([]corev1.Service, error) {
	if !utils.PrivilegedRights || cr.Spec.Integration == nil || cr.Spec.Integration.Jaeger == nil || !cr.Spec.Integration.Jaeger.CreateGrafanaDataSource {
		return nil, nil
	}
	allNamespaces, err := r.KubeClient.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		r.Log.Error(err, "Failed getting namespaces")
		return nil, err
	}
	// map with "namespace/service-name" as keys and services as values
	uniqeServices := make(map[string]corev1.Service)
	// Make list options with label selector
	listOptions := metav1.ListOptions{
		LabelSelector: labels.Set(utils.JaegerServiceLabels).String(),
	}
	for _, namespace := range allNamespaces.Items {
		if namespace.GetNamespace() == cr.GetNamespace() {
			continue
		}
		serviceList, err := r.KubeClient.CoreV1().Services(namespace.GetNamespace()).List(context.TODO(), listOptions)
		if err != nil {
			if errors.IsNotFound(err) {
				continue
			}
			r.Log.Error(err, "Failed getting Jaeger services")
			return nil, err
		}
		if serviceList != nil {
			for _, service := range serviceList.Items {
				uniqeServices[fmt.Sprintf("%s/%s", service.GetNamespace(), service.GetName())] = service
			}
		}
	}
	var services []corev1.Service
	for _, v := range uniqeServices {
		services = append(services, v)
	}
	if len(services) == 0 {
		r.Log.Info("Jaeger services is not found. Additional datasource will not be created")
	}
	sortServices(services)
	return services, nil
}

// Looking for Clickhouse Services in all namespaces except current using a label selector and return list of them or nil
func (r *GrafanaReconciler) getClickhouseServices(cr *v1alpha1.PlatformMonitoring) ([]corev1.Service, error) {
	if !utils.PrivilegedRights || cr.Spec.Integration == nil || cr.Spec.Integration.ClickHouse == nil || !cr.Spec.Integration.ClickHouse.CreateGrafanaDataSource {
		r.Log.Info(fmt.Sprintf("neto, utils.PrivilegedRights: %+v, cr.Spec.Integration: %+v", utils.PrivilegedRights, cr.Spec.Integration))
		return nil, nil
	}
	allNamespaces, err := r.KubeClient.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		r.Log.Error(err, "Failed getting namespaces")
		return nil, err
	}
	var services []corev1.Service
	for _, namespace := range allNamespaces.Items {
		if namespace.GetName() == cr.GetNamespace() {
			continue
		}
		serviceList, err := r.KubeClient.CoreV1().Services(namespace.GetName()).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			r.Log.Info(fmt.Sprintf("Error getting services in namespace:%s Error: %v", namespace.GetNamespace(), err))
			continue
		}
		for _, service := range serviceList.Items {
			if service.GetName() == utils.ClickHouseServiceName {
				services = append(services, service)
			}
		}
	}
	if len(services) == 0 {
		r.Log.Info("ClickHouse services is not found. Additional datasource will not be created")
	}
	sortServices(services)
	return services, nil
}

func sortServices(services []corev1.Service) {
	slices.SortFunc(services, func(a, b corev1.Service) int {
		// Order services by namespace
		if n := strings.Compare(a.Namespace, b.Namespace); n != 0 {
			return n
		}
		// If namespaces are equal, order services by name
		return strings.Compare(a.Name, b.Name)
	})
}
