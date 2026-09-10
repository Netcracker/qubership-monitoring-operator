package grafana

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	grafv1 "github.com/grafana/grafana-operator/v5/api/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type grafanaDatasourceSummary struct {
	Name string `json:"name"`
	UID  string `json:"uid"`
}

var (
	legacyGrafanaGVK = schema.GroupVersionKind{
		Group:   "integreatly.org",
		Version: "v1alpha1",
		Kind:    "Grafana",
	}
	legacyGrafanaDatasourceGVK = schema.GroupVersionKind{
		Group:   "integreatly.org",
		Version: "v1alpha1",
		Kind:    "GrafanaDataSource",
	}
)

func (r *GrafanaReconciler) adoptExistingDatasourceUID(
	ctx context.Context,
	platformMonitoring *monv1.PlatformMonitoring,
	datasource *grafv1.GrafanaDatasource,
) error {
	legacyDatasource := &unstructured.Unstructured{}
	legacyDatasource.SetGroupVersionKind(legacyGrafanaDatasourceGVK)
	legacyDatasource.SetName(datasource.Name)
	legacyDatasource.SetNamespace(datasource.Namespace)
	if err := r.Client.Get(ctx, client.ObjectKeyFromObject(legacyDatasource), legacyDatasource); err != nil {
		if meta.IsNoMatchError(err) || apierrors.IsNotFound(err) {
			return nil
		}
		return fmt.Errorf("checking legacy Grafana datasource: %w", err)
	}

	grafanaManifest, err := grafana(platformMonitoring)
	if err != nil {
		return err
	}

	currentGrafana := &grafv1.Grafana{}
	if err := r.Client.Get(ctx, client.ObjectKeyFromObject(grafanaManifest), currentGrafana); err != nil {
		return fmt.Errorf("getting Grafana instance for datasource migration: %w", err)
	}

	credentials := &corev1.Secret{}
	credentialsKey := client.ObjectKey{
		Name:      currentGrafana.Name + "-admin-credentials",
		Namespace: currentGrafana.Namespace,
	}
	if err := r.Client.Get(ctx, credentialsKey, credentials); err != nil {
		return fmt.Errorf("getting Grafana admin credentials for datasource migration: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet,
		strings.TrimRight(currentGrafana.Status.AdminURL, "/")+"/api/datasources", nil)
	if err != nil {
		return fmt.Errorf("creating Grafana datasource migration request: %w", err)
	}
	request.SetBasicAuth(
		string(credentials.Data["GF_SECURITY_ADMIN_USER"]),
		string(credentials.Data["GF_SECURITY_ADMIN_PASSWORD"]),
	)

	response, err := (&http.Client{Timeout: 10 * time.Second}).Do(request)
	if err != nil {
		return fmt.Errorf("listing Grafana datasources for migration: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("listing Grafana datasources for migration returned HTTP %d", response.StatusCode)
	}

	datasources := make([]grafanaDatasourceSummary, 0)
	if err := json.NewDecoder(response.Body).Decode(&datasources); err != nil {
		return fmt.Errorf("decoding Grafana datasources for migration: %w", err)
	}

	uid, found := findDatasourceUID(datasources, datasource.Spec.Datasource.Name)
	if !found {
		return nil
	}

	datasource.Spec.CustomUID = uid
	r.Log.Info("Preserving existing Grafana datasource UID",
		"name", datasource.Spec.Datasource.Name, "uid", uid)
	return nil
}

func findDatasourceUID(datasources []grafanaDatasourceSummary, name string) (string, bool) {
	for _, datasource := range datasources {
		if datasource.Name == name && datasource.UID != "" {
			return datasource.UID, true
		}
	}
	return "", false
}

func (r *GrafanaReconciler) migrateLegacyGrafanaResources(
	ctx context.Context,
	platformMonitoring *monv1.PlatformMonitoring,
	currentGrafana *grafv1.Grafana,
) error {
	legacyGrafana := &unstructured.Unstructured{}
	legacyGrafana.SetGroupVersionKind(legacyGrafanaGVK)
	key := client.ObjectKeyFromObject(currentGrafana)
	if err := r.Client.Get(ctx, key, legacyGrafana); err != nil {
		if meta.IsNoMatchError(err) {
			return nil
		}
		return client.IgnoreNotFound(err)
	}

	legacyDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: currentGrafana.Name + "-deployment", Namespace: currentGrafana.Namespace},
	}
	legacyServiceAccountName := ""
	if err := r.Client.Get(ctx, client.ObjectKeyFromObject(legacyDeployment), legacyDeployment); err != nil {
		if !apierrors.IsNotFound(err) {
			return err
		}
	} else {
		legacyServiceAccountName = legacyDeployment.Spec.Template.Spec.ServiceAccountName
	}

	adoptedResources := []client.Object{
		&corev1.Service{ObjectMeta: metav1.ObjectMeta{Name: currentGrafana.Name + "-service", Namespace: currentGrafana.Namespace}},
		legacyDeployment,
		&corev1.PersistentVolumeClaim{ObjectMeta: metav1.ObjectMeta{Name: currentGrafana.Name + "-pvc", Namespace: currentGrafana.Namespace}},
	}
	for _, resource := range adoptedResources {
		if err := r.reparentLegacyGrafanaResource(ctx, resource, legacyGrafana, currentGrafana); err != nil {
			return err
		}
	}

	preservedResources := []client.Object{
		&corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: currentGrafana.Name + "-admin-credentials", Namespace: currentGrafana.Namespace}},
		&corev1.Service{ObjectMeta: metav1.ObjectMeta{Name: currentGrafana.Name + "-alert", Namespace: currentGrafana.Namespace}},
		&corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: currentGrafana.Name + "-config", Namespace: currentGrafana.Namespace}},
		&corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: currentGrafana.Name + "-datasources", Namespace: currentGrafana.Namespace}},
	}
	for _, resource := range preservedResources {
		if err := r.reparentLegacyGrafanaResource(ctx, resource, legacyGrafana, platformMonitoring); err != nil {
			return err
		}
	}
	if err := r.adoptLegacyGrafanaServiceAccount(
		ctx, legacyServiceAccountName, currentGrafana.Namespace, legacyGrafana, platformMonitoring); err != nil {
		return err
	}

	propagation := metav1.DeletePropagationOrphan
	if err := r.Client.Delete(ctx, legacyGrafana, &client.DeleteOptions{PropagationPolicy: &propagation}); err != nil {
		return client.IgnoreNotFound(err)
	}

	r.Log.Info("Migrated legacy Grafana resources to grafana-operator v5", "name", key.Name, "namespace", key.Namespace)
	return nil
}

func (r *GrafanaReconciler) adoptLegacyGrafanaServiceAccount(
	ctx context.Context,
	name string,
	namespace string,
	legacyGrafana client.Object,
	platformMonitoring client.Object,
) error {
	if name == "" {
		return nil
	}

	serviceAccount := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
	}
	if err := r.Client.Get(ctx, client.ObjectKeyFromObject(serviceAccount), serviceAccount); err != nil {
		return client.IgnoreNotFound(err)
	}

	controller := metav1.GetControllerOf(serviceAccount)
	if controller != nil {
		if controller.UID == legacyGrafana.GetUID() {
			return r.reparentLegacyGrafanaResource(ctx, serviceAccount, legacyGrafana, platformMonitoring)
		}
		return nil
	}

	if err := controllerutil.SetControllerReference(platformMonitoring, serviceAccount, r.Scheme); err != nil {
		return fmt.Errorf("cannot adopt legacy Grafana service account %s/%s: %w", namespace, name, err)
	}
	return r.Client.Update(ctx, serviceAccount)
}

func (r *GrafanaReconciler) reparentLegacyGrafanaResource(
	ctx context.Context,
	resource client.Object,
	legacyGrafana client.Object,
	newOwner client.Object,
) error {
	if err := r.Client.Get(ctx, client.ObjectKeyFromObject(resource), resource); err != nil {
		return client.IgnoreNotFound(err)
	}

	ownerReferences := resource.GetOwnerReferences()
	hasLegacyOwner := false
	filteredOwnerReferences := make([]metav1.OwnerReference, 0, len(ownerReferences))
	for _, ownerReference := range ownerReferences {
		if ownerReference.UID == legacyGrafana.GetUID() {
			hasLegacyOwner = true
			continue
		}
		filteredOwnerReferences = append(filteredOwnerReferences, ownerReference)
	}
	if !hasLegacyOwner {
		return nil
	}
	resource.SetOwnerReferences(filteredOwnerReferences)

	if _, adoptForGrafanaOperator := newOwner.(*grafv1.Grafana); adoptForGrafanaOperator {
		labels := resource.GetLabels()
		if labels == nil {
			labels = make(map[string]string)
		}
		labels["app.kubernetes.io/managed-by"] = "grafana-operator"
		resource.SetLabels(labels)
	}

	if err := controllerutil.SetControllerReference(newOwner, resource, r.Scheme); err != nil {
		return fmt.Errorf("cannot set owner for %T %s/%s: %w",
			resource, resource.GetNamespace(), resource.GetName(), err)
	}
	if err := r.Client.Update(ctx, resource); err != nil {
		if apierrors.IsConflict(err) {
			return fmt.Errorf("cannot migrate %T %s/%s due to a concurrent update: %w",
				resource, resource.GetNamespace(), resource.GetName(), err)
		}
		return err
	}
	return nil
}
