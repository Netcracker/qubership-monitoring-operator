package grafana

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"reflect"
	"slices"
	"sort"
	"strings"
	"time"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	grafv1 "github.com/grafana/grafana-operator/v5/api/v1beta1"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
)

func (r *GrafanaReconciler) handleGrafana(cr *monv1.PlatformMonitoring) error {
	m, err := grafana(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Grafana manifest")
		return err
	}

	// Note: Config.AuthGenericOauth access removed as Config is now runtime.RawExtension in grafana-operator v5
	// OAuth configuration is handled in manifest.go during Grafana creation
	// Explicit GVK ensures correct API group (grafana.integreatly.org/v1beta1) for v5
	e := &grafv1.Grafana{ObjectMeta: m.ObjectMeta}
	e.SetGroupVersionKind(schema.GroupVersionKind{Group: "grafana.integreatly.org", Version: "v1beta1", Kind: "Grafana"})
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
	// Only update if something actually changed to avoid unnecessary updates
	needsUpdate := false
	if !reflect.DeepEqual(e.Spec, m.Spec) {
		e.Spec = m.Spec
		needsUpdate = true
	}
	if !reflect.DeepEqual(e.GetLabels(), m.GetLabels()) {
		e.SetLabels(m.GetLabels())
		needsUpdate = true
	}

	if needsUpdate {
		if err = r.UpdateResource(e); err != nil {
			return err
		}
	}
	// WA for https://github.com/grafana-operator/grafana-operator/issues/652
	r.Log.Info("Waiting grafana-deployment")
	time.Sleep(30 * time.Second)
	return nil
}

func (r *GrafanaReconciler) handleGrafanaDataSource(cr *monv1.PlatformMonitoring) error {
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
		r.Log.Error(err, "Failed creating GrafanaDatasource manifest")
		return err
	}

	// Set labels (asset has metadata.labels so m.Labels is non-nil)
	if m.Labels == nil {
		m.Labels = make(map[string]string)
	}
	m.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(m.GetName(), m.GetNamespace())
	m.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Grafana.Image)

	// Explicit GVK ensures correct API group (grafana.integreatly.org/v1beta1) for v5
	checkObj := &grafv1.GrafanaDatasource{}
	checkObj.SetName(m.GetName())
	checkObj.SetNamespace(m.GetNamespace())
	checkObj.SetGroupVersionKind(schema.GroupVersionKind{Group: "grafana.integreatly.org", Version: "v1beta1", Kind: "GrafanaDatasource"})
	if err = r.GetResource(checkObj); err != nil {
		if errors.IsNotFound(err) {
			if err = r.CreateResource(cr, m); err != nil {
				return err
			}
			return nil
		}
		return err
	}

	// Only update if something actually changed to avoid unnecessary updates
	needsUpdate := false
	if !reflect.DeepEqual(checkObj.Spec, m.Spec) {
		checkObj.Spec = m.Spec
		needsUpdate = true
	}
	if !reflect.DeepEqual(checkObj.GetLabels(), m.GetLabels()) {
		checkObj.SetLabels(m.GetLabels())
		needsUpdate = true
	}

	if needsUpdate {
		if err = r.UpdateResource(checkObj); err != nil {
			return err
		}
	}
	return nil
}

func (r *GrafanaReconciler) handleGrafanaPromxyDataSource(cr *monv1.PlatformMonitoring) error {
	m, err := grafanaPromxyDataSource(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating GrafanaPromxyDataSource manifest")
		return err
	}

	// Set labels (asset has metadata.labels so m.Labels is non-nil)
	if m.Labels == nil {
		m.Labels = make(map[string]string)
	}
	m.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(m.GetName(), m.GetNamespace())
	m.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Grafana.Image)

	// Explicit GVK ensures correct API group (grafana.integreatly.org/v1beta1) for v5
	checkObj := &grafv1.GrafanaDatasource{}
	checkObj.SetName(m.GetName())
	checkObj.SetNamespace(m.GetNamespace())
	checkObj.SetGroupVersionKind(schema.GroupVersionKind{Group: "grafana.integreatly.org", Version: "v1beta1", Kind: "GrafanaDatasource"})
	if err = r.GetResource(checkObj); err != nil {
		if errors.IsNotFound(err) {
			if err = r.CreateResource(cr, m); err != nil {
				return err
			}
			return nil
		}
		return err
	}

	// Set parameters
	// Only update if something actually changed to avoid unnecessary updates
	needsUpdate := false
	if !reflect.DeepEqual(checkObj.Spec, m.Spec) {
		checkObj.Spec = m.Spec
		needsUpdate = true
	}
	if !reflect.DeepEqual(checkObj.GetLabels(), m.GetLabels()) {
		checkObj.SetLabels(m.GetLabels())
		needsUpdate = true
	}

	if needsUpdate {
		if err = r.UpdateResource(checkObj); err != nil {
			return err
		}
	}
	return nil
}

func (r *GrafanaReconciler) handleIngressV1(cr *monv1.PlatformMonitoring) error {
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

func (r *GrafanaReconciler) handlePodMonitor(cr *monv1.PlatformMonitoring) error {
	m, err := grafanaPodMonitor(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating PodMonitor manifest")
		return err
	}

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

// getGrafanaAdminSecretName returns the name of the admin credentials secret for Grafana.
// The secret name pattern is: {grafana-name}-admin-credentials
func getGrafanaAdminSecretName(cr *monv1.PlatformMonitoring) string {
	grafanaName := utils.GrafanaComponentName // default name from asset
	if cr.Spec.Grafana != nil && cr.Spec.Grafana.Name != "" {
		grafanaName = cr.Spec.Grafana.Name
	}
	return fmt.Sprintf("%s-admin-credentials", grafanaName)
}

// getGrafanaNamespace returns the namespace where the Grafana instance lives,
// which may differ from the PlatformMonitoring namespace when spec.grafana.namespace is set.
func getGrafanaNamespace(cr *monv1.PlatformMonitoring) string {
	if cr.Spec.Grafana != nil && cr.Spec.Grafana.Namespace != "" {
		return cr.Spec.Grafana.Namespace
	}
	return cr.GetNamespace()
}

// getGrafanaName returns the Grafana CR name, defaulting to utils.GrafanaComponentName.
func getGrafanaName(cr *monv1.PlatformMonitoring) string {
	if cr.Spec.Grafana != nil && cr.Spec.Grafana.Name != "" {
		return cr.Spec.Grafana.Name
	}
	return utils.GrafanaComponentName
}

// computeSecretChecksum returns a deterministic SHA256 hex digest of a secret data map.
// Keys are sorted before hashing to guarantee a stable result regardless of map iteration order.
func computeSecretChecksum(data map[string][]byte) string {
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	h := sha256.New()
	for _, k := range keys {
		h.Write([]byte(k))
		h.Write(data[k])
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

// handleGrafanaCredentialsSecret computes a SHA256 checksum of the admin credentials secret
// and compares it with the checksum stored in the existing Grafana CR pod-template annotation.
// When the checksums differ (secret changed since last reconcile), isSecretUpdated is set to
// true so that resetGrafanaCredentials is invoked later in the same reconcile cycle.
// The new checksum is stored in currentAdminSecretChecksum and written into the Grafana CR
// annotation by grafana() in manifest.go, which causes grafana-operator to trigger a rolling
// restart of the Grafana Deployment automatically.
func (r *GrafanaReconciler) handleGrafanaCredentialsSecret(cr *monv1.PlatformMonitoring) error {
	secretName := getGrafanaAdminSecretName(cr)
	secretNamespace := getGrafanaNamespace(cr)

	adminSecret := &corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: secretName, Namespace: secretNamespace}}
	if err := r.GetResource(adminSecret); err != nil {
		if errors.IsNotFound(err) {
			r.Log.Info("Grafana admin credentials secret not found; skipping credential sync",
				"secret", secretName, "namespace", secretNamespace)
			return nil
		}
		return err
	}
	if len(adminSecret.Data) == 0 {
		r.Log.Info("Grafana admin credentials secret has no data; skipping credential sync",
			"secret", secretName, "namespace", secretNamespace)
		return nil
	}

	newChecksum := computeSecretChecksum(adminSecret.Data)

	// Read the last-applied checksum from the existing Grafana CR pod-template annotation.
	// This survives operator restarts without requiring additional persistent storage.
	// resetGrafanaCredentials is only triggered when the CR already exists (i.e. Grafana was
	// previously deployed). On first install the admin password is applied from the mounted
	// secret file automatically by Grafana on startup, so no exec is needed.
	existingCR := &grafv1.Grafana{}
	existingCR.SetName(getGrafanaName(cr))
	existingCR.SetNamespace(getGrafanaNamespace(cr))
	existingCR.SetGroupVersionKind(schema.GroupVersionKind{
		Group: "grafana.integreatly.org", Version: "v1beta1", Kind: "Grafana",
	})
	crExists := true
	lastChecksum := ""
	if err := r.GetResource(existingCR); err != nil {
		if !errors.IsNotFound(err) {
			r.Log.Info("Cannot read Grafana CR for checksum comparison; will not trigger credential reset", "err", err)
		}
		crExists = false
	} else {
		if existingCR.Spec.Deployment != nil &&
			existingCR.Spec.Deployment.Spec.Template != nil &&
			existingCR.Spec.Deployment.Spec.Template.Annotations != nil {
			lastChecksum = existingCR.Spec.Deployment.Spec.Template.Annotations[adminSecretChecksumAnnotation]
		}
	}

	currentAdminSecretChecksum = newChecksum
	if crExists && newChecksum != lastChecksum {
		r.Log.Info("Admin credentials secret changed; Grafana credentials will be reset",
			"secret", secretName)
		isSecretUpdated = true
	}
	return nil
}

// resetGrafanaCredentials resets the Grafana admin password in the running database via
// `grafana cli admin reset-admin-password`. This is required when a Persistent Volume is used
// (SQLite or external DB) because Grafana ignores admin_password from config once the admin user
// already exists in the database. For non-PV deployments the rolling restart triggered by the
// checksum annotation is sufficient; calling this function in that case is safe and idempotent.
func (r *GrafanaReconciler) resetGrafanaCredentials(cr *monv1.PlatformMonitoring) error {
	grafanaNamespace := getGrafanaNamespace(cr)
	grafanaName := getGrafanaName(cr)

	r.Log.Info("Waiting for Grafana pods readiness before credential reset",
		"deployment", utils.GrafanaDeploymentName, "namespace", grafanaNamespace)
	if err := r.WaitForPodsReadiness(&appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.GrafanaDeploymentName,
			Namespace: grafanaNamespace,
		},
	}); err != nil {
		return fmt.Errorf("grafana deployment not ready: %w", err)
	}
	r.Log.Info("Grafana pods are ready")

	secretName := getGrafanaAdminSecretName(cr)
	adminSecret := &corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: secretName, Namespace: grafanaNamespace}}
	if err := r.GetResource(adminSecret); err != nil {
		if errors.IsNotFound(err) {
			r.Log.Info("Admin credentials secret not found; skipping credential reset", "secret", secretName)
			isSecretUpdated = false
			return nil
		}
		return err
	}
	newPassword := string(adminSecret.Data["GF_SECURITY_ADMIN_PASSWORD"])
	if newPassword == "" {
		r.Log.Info("GF_SECURITY_ADMIN_PASSWORD is empty; skipping credential reset", "secret", secretName)
		isSecretUpdated = false
		return nil
	}

	// Find a running Grafana pod. Pods are labelled with app.kubernetes.io/name=<grafana-name>
	// via spec.deployment.spec.template.metadata.labels set in manifest.go.
	pods, err := r.KubeClient.CoreV1().Pods(grafanaNamespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app.kubernetes.io/name=%s", utils.TruncLabel(grafanaName)),
	})
	if err != nil {
		return fmt.Errorf("cannot list Grafana pods: %w", err)
	}
	var podName string
	for _, p := range pods.Items {
		if p.DeletionTimestamp == nil {
			podName = p.Name
			break
		}
	}
	if podName == "" {
		return fmt.Errorf("no running Grafana pod found with label app.kubernetes.io/name=%s in namespace %s",
			grafanaName, grafanaNamespace)
	}
	r.Log.Info("Found Grafana pod for credential reset", "pod", podName)

	// Official command per Grafana Operator docs for deployments with external/persistent databases.
	command := []string{
		"grafana", "cli",
		"--homepath", "/usr/share/grafana",
		"--config", "/etc/grafana/grafana.ini",
		"admin", "reset-admin-password", newPassword,
	}
	req := r.KubeClient.CoreV1().RESTClient().
		Post().
		Resource("pods").
		Name(podName).
		Namespace(grafanaNamespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: "grafana",
			Command:   command,
			Stdin:     false,
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(r.config, "POST", req.URL())
	if err != nil {
		return fmt.Errorf("cannot create executor for Grafana pod %s: %w", podName, err)
	}

	var stdout, stderr bytes.Buffer
	if err = exec.StreamWithContext(context.TODO(), remotecommand.StreamOptions{
		Stdout: &stdout,
		Stderr: &stderr,
	}); err != nil {
		return fmt.Errorf("grafana cli reset-admin-password failed: %v; stdout: %s; stderr: %s",
			err, stdout.String(), stderr.String())
	}

	isSecretUpdated = false
	r.Log.Info("Grafana admin credentials reset successfully")
	return nil
}

func (r *GrafanaReconciler) deleteGrafana(cr *monv1.PlatformMonitoring) error {
	m, err := grafana(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Grafana manifest")
		return err
	}
	// Check if resource exists first
	checkObj := &grafv1.Grafana{}
	checkObj.SetName(m.GetName())
	checkObj.SetNamespace(m.GetNamespace())
	checkObj.SetGroupVersionKind(schema.GroupVersionKind{Group: "grafana.integreatly.org", Version: "v1beta1", Kind: "Grafana"})
	if err = r.GetResource(checkObj); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	// Use the manifest object (which has correct type) for deletion
	// The manifest object already has GVK set correctly
	if err = r.Client.Delete(context.TODO(), m); err != nil {
		return err
	}
	r.Log.Info("Successful deleting", "resource", "Grafana", "name", m.GetName())
	return nil
}

func (r *GrafanaReconciler) deleteGrafanaDataSource(cr *monv1.PlatformMonitoring) error {
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
		r.Log.Error(err, "Failed creating GrafanaDatasource manifest")
		return err
	}
	// Explicit GVK ensures correct API group (grafana.integreatly.org/v1beta1) for v5
	checkObj := &grafv1.GrafanaDatasource{}
	checkObj.SetName(m.GetName())
	checkObj.SetNamespace(m.GetNamespace())
	checkObj.SetGroupVersionKind(schema.GroupVersionKind{Group: "grafana.integreatly.org", Version: "v1beta1", Kind: "GrafanaDatasource"})
	if err = r.GetResource(checkObj); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	// Use the manifest object (which has correct type) for deletion
	// The manifest object already has GVK set correctly
	if err = r.Client.Delete(context.TODO(), m); err != nil {
		return err
	}
	r.Log.Info("Successful deleting", "resource", "GrafanaDatasource", "name", m.GetName())
	return nil
}

func (r *GrafanaReconciler) deleteGrafanaPromxyDataSource(cr *monv1.PlatformMonitoring) error {
	m, err := grafanaPromxyDataSource(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating GrafanaPromxyDataSource manifest")
		return err
	}
	checkObj := &grafv1.GrafanaDatasource{}
	checkObj.SetName(m.GetName())
	checkObj.SetNamespace(m.GetNamespace())
	checkObj.SetGroupVersionKind(schema.GroupVersionKind{Group: "grafana.integreatly.org", Version: "v1beta1", Kind: "GrafanaDatasource"})
	if err = r.GetResource(checkObj); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if err = r.Client.Delete(context.TODO(), m); err != nil {
		return err
	}
	r.Log.Info("Successful deleting", "resource", "GrafanaDatasource", "name", m.GetName())
	return nil
}

func (r *GrafanaReconciler) deleteIngressV1(cr *monv1.PlatformMonitoring) error {
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

func (r *GrafanaReconciler) deletePodMonitor(cr *monv1.PlatformMonitoring) error {
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
func (r *GrafanaReconciler) getJaegerServices(cr *monv1.PlatformMonitoring) ([]corev1.Service, error) {
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
func (r *GrafanaReconciler) getClickhouseServices(cr *monv1.PlatformMonitoring) ([]corev1.Service, error) {
	if !utils.PrivilegedRights || cr.Spec.Integration == nil || cr.Spec.Integration.ClickHouse == nil || !cr.Spec.Integration.ClickHouse.CreateGrafanaDataSource {
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
			r.Log.Info(fmt.Sprintf("Error getting services in namespace:%s Error: %v", namespace.GetName(), err))
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
