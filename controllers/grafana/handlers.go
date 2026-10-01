package grafana

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"slices"
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
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	grafanaExtraVarsConfigMapResourceVersionAnnotation = "monitoring.netcracker.com/grafana-extra-vars-configmap-resource-version"
	grafanaExtraVarsSecretResourceVersionAnnotation    = "monitoring.netcracker.com/grafana-extra-vars-secret-resource-version"
)

func grafanaPodTemplateAnnotations(manifest *grafv1.Grafana) map[string]string {
	if manifest == nil || manifest.Spec.Deployment == nil || manifest.Spec.Deployment.Spec.Template == nil {
		return nil
	}
	return manifest.Spec.Deployment.Spec.Template.Annotations
}

func (r *GrafanaReconciler) addGrafanaExtraVarsResourceVersions(
	ctx context.Context,
	namespace string,
	manifest *grafv1.Grafana,
	existingAnnotations map[string]string,
) error {
	configMap, err := r.KubeClient.CoreV1().ConfigMaps(namespace).Get(ctx, "grafana-extra-vars", metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return fmt.Errorf("cannot get Grafana extra-vars ConfigMap: %w", err)
	}
	configMapFound := err == nil
	secret, err := r.KubeClient.CoreV1().Secrets(namespace).Get(ctx, "grafana-extra-vars-secret", metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return fmt.Errorf("cannot get Grafana extra-vars Secret: %w", err)
	}
	secretFound := err == nil

	annotations := manifest.Spec.Deployment.Spec.Template.Annotations
	setAnnotation := func(key, value string) {
		if annotations == nil {
			annotations = make(map[string]string)
			manifest.Spec.Deployment.Spec.Template.Annotations = annotations
		}
		annotations[key] = value
	}
	if configMapFound {
		setAnnotation(grafanaExtraVarsConfigMapResourceVersionAnnotation, configMap.ResourceVersion)
	} else if resourceVersion, ok := existingAnnotations[grafanaExtraVarsConfigMapResourceVersionAnnotation]; ok {
		setAnnotation(grafanaExtraVarsConfigMapResourceVersionAnnotation, resourceVersion)
	}
	if secretFound {
		setAnnotation(grafanaExtraVarsSecretResourceVersionAnnotation, secret.ResourceVersion)
	} else if resourceVersion, ok := existingAnnotations[grafanaExtraVarsSecretResourceVersionAnnotation]; ok {
		setAnnotation(grafanaExtraVarsSecretResourceVersionAnnotation, resourceVersion)
	}
	return nil
}

// grafanaSecretNamespace returns the namespace holding Grafana's credential Secrets.
func grafanaSecretNamespace(cr *monv1.PlatformMonitoring) string {
	if cr.Spec.Grafana != nil && cr.Spec.Grafana.Namespace != "" {
		return cr.Spec.Grafana.Namespace
	}
	return cr.GetNamespace()
}

// secretHasKeys reports whether the named Secret exists and carries every key. A lookup error
// other than NotFound is treated as absent: emitting a $__file{} reference for a Secret we
// cannot confirm risks wedging Grafana at startup, while omitting it only loses the credential
// until the next reconcile.
func (r *GrafanaReconciler) secretHasKeys(name, namespace string, keys ...string) bool {
	secret := &corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace}}
	if err := r.GetResource(secret); err != nil {
		if !errors.IsNotFound(err) {
			r.Log.Error(err, "Cannot read Secret; treating it as absent", "secret", name, "namespace", namespace)
		}
		return false
	}
	for _, key := range keys {
		if _, ok := secret.Data[key]; !ok {
			r.Log.Info("Secret is missing a required key", "secret", name, "namespace", namespace, "key", key)
			return false
		}
	}
	return true
}

// observeCredentialSources records which optional credential Secrets exist right now, so the
// manifest only references files that Grafana can actually expand.
func (r *GrafanaReconciler) observeCredentialSources(cr *monv1.PlatformMonitoring) grafanaCredentialSources {
	namespace := grafanaSecretNamespace(cr)
	sources := grafanaCredentialSources{}
	// Helm guarantees the admin Secret unless the user opted out, so only that case needs a lookup.
	if grafanaAdminSecretUserManaged(cr) {
		sources.AdminSecretPresent = r.secretHasKeys(
			getGrafanaAdminSecretName(cr), namespace,
			"GF_SECURITY_ADMIN_USER", "GF_SECURITY_ADMIN_PASSWORD",
		)
	}
	if cr.Spec.Auth != nil {
		sources.OAuthSecretPresent = r.secretHasKeys(
			grafanaOAuthClientSecretName, namespace,
			"GF_AUTH_GENERIC_OAUTH_CLIENT_SECRET",
		)
	}
	return sources
}

func (r *GrafanaReconciler) handleGrafana(cr *monv1.PlatformMonitoring) error {
	isOpenShift, err := r.IsOpenShift()
	if err != nil {
		return err
	}
	m, err := grafana(cr, r.observeCredentialSources(cr), isOpenShift)
	if err != nil {
		r.Log.Error(err, "Failed creating Grafana manifest")
		return err
	}
	// Note: Config.AuthGenericOauth access removed as Config is now runtime.RawExtension in grafana-operator v5
	// OAuth configuration is handled in manifest.go during Grafana creation
	// Explicit GVK ensures correct API group (grafana.integreatly.org/v1beta1) for v5
	e := &grafv1.Grafana{ObjectMeta: m.ObjectMeta}
	e.SetGroupVersionKind(schema.GroupVersionKind{Group: "grafana.integreatly.org", Version: "v1beta1", Kind: "Grafana"})
	err = r.GetResource(e)
	grafanaExists := err == nil
	if err != nil && !errors.IsNotFound(err) {
		return err
	}

	var existingAnnotations map[string]string
	if grafanaExists {
		existingAnnotations = grafanaPodTemplateAnnotations(e)
	}
	if err = r.addGrafanaExtraVarsResourceVersions(
		context.TODO(), m.GetNamespace(), m, existingAnnotations,
	); err != nil {
		r.Log.Error(err, "Failed adding Grafana extra-vars resource versions")
		return err
	}
	if !grafanaExists {
		if err = r.CreateResource(cr, m); err != nil {
			return err
		}
		return r.migrateLegacyGrafanaResources(context.TODO(), cr, m)
	}

	if applyGrafanaDesiredState(e, m) {
		if err = r.UpdateResource(e); err != nil {
			return err
		}
	}
	if err = r.migrateLegacyGrafanaResources(context.TODO(), cr, e); err != nil {
		return err
	}
	// WA for https://github.com/grafana-operator/grafana-operator/issues/652
	r.Log.Info("Waiting grafana-deployment")
	time.Sleep(30 * time.Second)
	return nil
}

func applyGrafanaDesiredState(existing, desired *grafv1.Grafana) bool {
	needsUpdate := false
	if !reflect.DeepEqual(existing.Spec, desired.Spec) {
		existing.Spec = desired.Spec
		needsUpdate = true
	}
	if !reflect.DeepEqual(existing.GetLabels(), desired.GetLabels()) {
		existing.SetLabels(desired.GetLabels())
		needsUpdate = true
	}
	return needsUpdate
}

func (r *GrafanaReconciler) handleGrafanaDataSource(cr *monv1.PlatformMonitoring) error {
	m, err := grafanaDataSource(cr)
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
			if err = r.adoptExistingDatasourceUID(context.TODO(), cr, m); err != nil {
				return err
			}
			if err = r.CreateResource(cr, m); err != nil {
				return err
			}
			return nil
		}
		return err
	}
	if m.Spec.CustomUID == "" {
		m.Spec.CustomUID = checkObj.Spec.CustomUID
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

// handleJaegerDataSources creates one GrafanaDatasource for each discovered Jaeger Service.
// A discovery error is returned before any owned datasource is created or deleted.
func (r *GrafanaReconciler) handleJaegerDataSources(cr *monv1.PlatformMonitoring) error {
	services, err := r.getJaegerServices(cr)
	if err != nil {
		r.Log.Error(err, "Failed getting Jaeger services")
		return err
	}
	return r.syncDiscoveredDataSources(cr, grafanaJaegerDataSources(cr, services), jaegerDatasourceComponent)
}

// handleClickHouseDataSources creates one GrafanaDatasource for each discovered ClickHouse Service.
// A discovery or credentials error is returned before any owned datasource is created or deleted.
func (r *GrafanaReconciler) handleClickHouseDataSources(cr *monv1.PlatformMonitoring) error {
	services, err := r.getClickhouseServices(cr)
	if err != nil {
		r.Log.Error(err, "Failed getting ClickHouse services")
		return err
	}
	desired, err := grafanaClickHouseDataSources(cr, r.KubeClient, services)
	if err != nil {
		r.Log.Error(err, "Failed creating ClickHouse GrafanaDatasource manifests")
		return err
	}
	return r.syncDiscoveredDataSources(cr, desired, clickHouseDatasourceComponent)
}

func (r *GrafanaReconciler) syncDiscoveredDataSources(cr *monv1.PlatformMonitoring, desired []*grafv1.GrafanaDatasource, component string) error {
	desiredNames := make(map[string]struct{}, len(desired))
	for _, datasource := range desired {
		desiredNames[datasource.GetName()] = struct{}{}
		if err := r.applyDiscoveredDataSource(cr, datasource); err != nil {
			return err
		}
	}
	return r.deleteStaleDiscoveredDataSources(cr, component, desiredNames)
}

func (r *GrafanaReconciler) applyDiscoveredDataSource(cr *monv1.PlatformMonitoring, desired *grafv1.GrafanaDatasource) error {
	if desired.Labels == nil {
		desired.Labels = make(map[string]string)
	}
	desired.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(desired.GetName(), desired.GetNamespace())
	desired.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Grafana.Image)

	current := &grafv1.GrafanaDatasource{}
	current.SetName(desired.GetName())
	current.SetNamespace(desired.GetNamespace())
	current.SetGroupVersionKind(desired.GroupVersionKind())
	if err := r.GetResource(current); err != nil {
		if errors.IsNotFound(err) {
			if err = r.adoptDiscoveredDatasourceUID(context.TODO(), cr, desired); err != nil {
				return err
			}
			return r.CreateResource(cr, desired)
		}
		return err
	}

	if current.GetLabels()[grafanaCleanupLabelKey] != grafanaCleanupLabelValue {
		r.Log.Info("Skipping unowned GrafanaDatasource", "name", current.GetName())
		return nil
	}
	desired.Spec.CustomUID = current.Spec.CustomUID
	desired.Spec.InstanceSelector = current.Spec.InstanceSelector
	needsUpdate := false
	if !reflect.DeepEqual(current.Spec, desired.Spec) {
		current.Spec = desired.Spec
		needsUpdate = true
	}
	if !reflect.DeepEqual(current.GetLabels(), desired.GetLabels()) {
		current.SetLabels(desired.GetLabels())
		needsUpdate = true
	}
	if needsUpdate {
		return r.UpdateResource(current)
	}
	return nil
}

func (r *GrafanaReconciler) deleteStaleDiscoveredDataSources(cr *monv1.PlatformMonitoring, component string, desiredNames map[string]struct{}) error {
	listed := &grafv1.GrafanaDatasourceList{}
	if err := r.Client.List(context.TODO(), listed,
		client.InNamespace(cr.GetNamespace()),
		client.MatchingLabels{
			"app.kubernetes.io/component": component,
			grafanaCleanupLabelKey:        grafanaCleanupLabelValue,
		},
	); err != nil {
		return err
	}
	for i := range listed.Items {
		datasource := listed.Items[i]
		if _, keep := desiredNames[datasource.GetName()]; keep {
			continue
		}
		datasource.SetGroupVersionKind(schema.GroupVersionKind{Group: "grafana.integreatly.org", Version: "v1beta1", Kind: "GrafanaDatasource"})
		if err := r.Client.Delete(context.TODO(), &datasource); err != nil && !errors.IsNotFound(err) {
			return err
		}
		r.Log.Info("Successful deleting", "resource", "GrafanaDatasource", "name", datasource.GetName())
	}
	return nil
}

func (r *GrafanaReconciler) deleteDiscoveredDataSources(cr *monv1.PlatformMonitoring, component string) error {
	return r.deleteStaleDiscoveredDataSources(cr, component, nil)
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
			if err = r.adoptExistingDatasourceUID(context.TODO(), cr, m); err != nil {
				return err
			}
			if err = r.CreateResource(cr, m); err != nil {
				return err
			}
			return nil
		}
		return err
	}
	if m.Spec.CustomUID == "" {
		m.Spec.CustomUID = checkObj.Spec.CustomUID
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

// handleGrafanaCredentialsSecret checks the Helm-managed admin credentials secret when
// disableDefaultAdminSecret=false. It is only used to log status; it never fails reconciliation.
// If the secret is missing or incomplete, Grafana still has a fallback (built-in admin/admin when
// env vars are optional, or the pod will use the secret once Helm creates it). The user can always
// authenticate one way or another.
func (r *GrafanaReconciler) handleGrafanaCredentialsSecret(cr *monv1.PlatformMonitoring) (err error) {
	secretName := getGrafanaAdminSecretName(cr)
	secretNamespace := cr.GetNamespace()
	if cr.Spec.Grafana != nil && cr.Spec.Grafana.Namespace != "" {
		secretNamespace = cr.Spec.Grafana.Namespace
	}

	e := &corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: secretName, Namespace: secretNamespace}}
	err = r.GetResource(e)
	if err != nil {
		if errors.IsNotFound(err) {
			r.Log.Info("Grafana admin credentials secret not found; Grafana may use built-in default until Helm creates the secret",
				"secret", secretName, "namespace", secretNamespace)
			return nil
		}
		return err
	}

	if e.Data == nil {
		r.Log.Info("Grafana admin credentials secret has no data", "secret", secretName, "namespace", secretNamespace)
		return nil
	}
	if _, ok := e.Data["GF_SECURITY_ADMIN_USER"]; !ok {
		r.Log.Info("Grafana admin credentials secret missing GF_SECURITY_ADMIN_USER", "secret", secretName, "namespace", secretNamespace)
		return nil
	}
	if _, ok := e.Data["GF_SECURITY_ADMIN_PASSWORD"]; !ok {
		r.Log.Info("Grafana admin credentials secret missing GF_SECURITY_ADMIN_PASSWORD", "secret", secretName, "namespace", secretNamespace)
		return nil
	}

	r.Log.Info("Grafana admin credentials secret validated", "secret", secretName, "namespace", secretNamespace)
	return nil
}

//nolint:unused // Kept for manual Grafana credential recovery.
func (r *GrafanaReconciler) resetGrafanaCredentials(cr *monv1.PlatformMonitoring) (err error) {
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
	// Getting Admin Credentials Secret
	r.Log.Info("Getting Admin Credentials Secret")
	secretName := getGrafanaAdminSecretName(cr)
	secretNamespace := cr.GetNamespace()
	if cr.Spec.Grafana != nil && cr.Spec.Grafana.Namespace != "" {
		secretNamespace = cr.Spec.Grafana.Namespace
	}
	adminSecret := &corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: secretName, Namespace: secretNamespace}}
	if err = r.GetResource(adminSecret); err == nil {
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
		command := []string{"grafana", "cli", "admin", "reset-admin-password", string(adminSecret.Data["GF_SECURITY_ADMIN_PASSWORD"])}
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

		r.Log.Info("Grafana Credentials Reset was finished")
	}
	if errors.IsNotFound(err) {
		r.Log.Info("Admin Credentials Secret wasn't found")
		return nil
	}
	return err
}

func (r *GrafanaReconciler) deleteGrafana(cr *monv1.PlatformMonitoring) error {
	// Address the deletion target by its stable identity so uninstall does not depend on
	// platform discovery, credential Secret lookup, or desired-state validation.
	existing := &grafv1.Grafana{}
	existing.SetName(grafanaName(cr))
	existing.SetNamespace(grafanaNamespace(cr))
	existing.SetGroupVersionKind(schema.GroupVersionKind{Group: "grafana.integreatly.org", Version: "v1beta1", Kind: "Grafana"})
	if err := r.GetResource(existing); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if err := r.Client.Delete(context.TODO(), existing); err != nil {
		return err
	}
	r.Log.Info("Successful deleting", "resource", "Grafana", "name", existing.GetName())
	return nil
}

func (r *GrafanaReconciler) deleteGrafanaDataSource(cr *monv1.PlatformMonitoring) error {
	m, err := grafanaDataSource(cr)
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
			return nil, fmt.Errorf("listing services in namespace %q: %w", namespace.GetName(), err)
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
