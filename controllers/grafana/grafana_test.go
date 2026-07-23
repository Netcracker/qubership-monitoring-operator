package grafana

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils/labelsassert"
	grafv1 "github.com/grafana/grafana-operator/v5/api/v1beta1"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var (
	cr              *monv1.PlatformMonitoring
	labelKey        = "label.key"
	labelValue      = "label-value"
	annotationKey   = "annotation.key"
	annotationValue = "annotation-value"
)

func TestGrafanaManifests(t *testing.T) {
	cr = &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
		Spec: monv1.PlatformMonitoringSpec{
			Grafana: &monv1.Grafana{
				Annotations: map[string]string{annotationKey: annotationValue},
				Labels:      map[string]string{labelKey: labelValue},
			},
		},
	}
	t.Run("Test Grafana manifest", func(t *testing.T) {
		m, err := grafana(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Grafana manifest should not be empty")
		assert.NotNil(t, m.GetLabels())
		assert.Equal(t, labelValue, m.GetLabels()[labelKey])
		assert.NotNil(t, m.Spec.Deployment)
		assert.NotNil(t, m.Spec.Deployment.Spec.Template.Labels)
		assert.Equal(t, labelValue, m.Spec.Deployment.Spec.Template.Labels[labelKey])
		assert.NotNil(t, m.Spec.Deployment.Spec.Template.Annotations)
		assert.Equal(t, annotationValue, m.Spec.Deployment.Spec.Template.Annotations[annotationKey])
		assert.NotNil(t, m.GetAnnotations())
		assert.Equal(t, annotationValue, m.GetAnnotations()[annotationKey])
	})
	t.Run("Test Grafana manifest preserves cleanup label", func(t *testing.T) {
		cr.Spec.Grafana.Labels["app.kubernetes.io/managed-by"] = "custom-manager"
		cr.Spec.Grafana.Labels["app.kubernetes.io/managed-by-operator"] = "custom-manager"

		m, err := grafana(cr)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, "custom-manager", m.GetLabels()["app.kubernetes.io/managed-by"])
		assert.Equal(t, "monitoring-operator", m.GetLabels()["app.kubernetes.io/managed-by-operator"])
	})
	cr = &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
		Spec: monv1.PlatformMonitoringSpec{
			Grafana: &monv1.Grafana{},
		},
	}
	// Disabled for v5: in v5 labels/annotations live in Deployment.Spec.Template, not Deployment
	//t.Run("Test Grafana manifest with nil annotation", func(t *testing.T) {
	//	m, err := grafana(cr)
	//	...
	//})
	t.Run("Test GrafanaDatasource manifest", func(t *testing.T) {
		m, err := grafanaDataSource(cr, nil, nil, nil)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "GrafanaDatasource manifest should not be empty")
		assert.Equal(t, 10*time.Minute, m.Spec.ResyncPeriod.Duration)
	})
	t.Run("Test GrafanaPromxyDatasource manifest", func(t *testing.T) {
		m, err := grafanaPromxyDataSource(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "GrafanaPromxyDatasource manifest should not be empty")
		assert.Equal(t, "platform-monitoring-promxy", m.GetName())
		assert.Equal(t, 10*time.Minute, m.Spec.ResyncPeriod.Duration)
		if m.Spec.Datasource != nil {
			assert.Contains(t, m.Spec.Datasource.URL, "promxy")
		}
	})
	t.Run("Test Ingress v1 manifest", func(t *testing.T) {
		m, err := grafanaIngressV1(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Ingress v1 manifest should not be empty")
	})
	t.Run("Test PodMonitor manifest", func(t *testing.T) {
		crWithLabels := &monv1.PlatformMonitoring{
			ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring", Labels: map[string]string{labelKey: labelValue}},
			Spec: monv1.PlatformMonitoringSpec{
				Grafana: &monv1.Grafana{Labels: map[string]string{labelKey: labelValue}},
			},
		}
		m, err := grafanaPodMonitor(crWithLabels)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "PodMonitor manifest should not be empty")
		labelsassert.AssertCRLabels(t, m.GetLabels(), utils.GrafanaComponentName, "victoriametrics-operator", map[string]string{labelKey: labelValue})
	})
}

func TestGrafanaManifestPreservesDataStorage(t *testing.T) {
	monitoring := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"},
		Spec: monv1.PlatformMonitoringSpec{
			Grafana: &monv1.Grafana{
				DataStorage: &runtime.RawExtension{
					Raw: []byte(`{"accessModes":["ReadWriteOnce"],"class":"standard","size":"2Gi"}`),
				},
			},
		},
	}

	manifest, err := grafana(monitoring)
	if err != nil {
		t.Fatal(err)
	}

	if assert.NotNil(t, manifest.Spec.PersistentVolumeClaim) &&
		assert.NotNil(t, manifest.Spec.PersistentVolumeClaim.Spec) {
		assert.Equal(t, []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			manifest.Spec.PersistentVolumeClaim.Spec.AccessModes)
		assert.Equal(t, "standard", *manifest.Spec.PersistentVolumeClaim.Spec.StorageClassName)
		assert.Equal(t, "2Gi",
			manifest.Spec.PersistentVolumeClaim.Spec.Resources.Requests.Storage().String())
	}
	if assert.NotNil(t, manifest.Spec.Deployment) &&
		assert.NotNil(t, manifest.Spec.Deployment.Spec.Strategy) {
		assert.Equal(t, appsv1.RecreateDeploymentStrategyType, manifest.Spec.Deployment.Spec.Strategy.Type)
	}
	assert.Equal(t, "Container", manifest.Spec.DisableDefaultSecurityContext)

	podSpec := manifest.Spec.Deployment.Spec.Template.Spec
	assert.Contains(t, podSpec.Volumes, corev1.Volume{
		Name: "grafana-data",
		VolumeSource: corev1.VolumeSource{
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{ClaimName: "grafana-pvc"},
		},
	})
	assert.Contains(t, podSpec.Containers[0].VolumeMounts, corev1.VolumeMount{
		Name:      "grafana-data",
		MountPath: "/var/lib/grafana",
	})
}

func TestFindDatasourceUID(t *testing.T) {
	datasources := []grafanaDatasourceSummary{
		{Name: "Other", UID: "other"},
		{Name: "Platform Monitoring Prometheus", UID: "legacy-prometheus"},
	}

	uid, found := findDatasourceUID(datasources, "Platform Monitoring Prometheus")
	assert.True(t, found)
	assert.Equal(t, "legacy-prometheus", uid)

	_, found = findDatasourceUID(datasources, "Missing")
	assert.False(t, found)
}

func TestAdoptExistingDatasourceUID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		username, password, ok := request.BasicAuth()
		assert.True(t, ok)
		assert.Equal(t, "admin", username)
		assert.Equal(t, "password", password)
		assert.Equal(t, "/api/datasources", request.URL.Path)
		writer.Header().Set("Content-Type", "application/json")
		_, err := writer.Write([]byte(`[{"name":"Platform Monitoring Prometheus","uid":"legacy-prometheus"}]`))
		assert.NoError(t, err)
	}))
	t.Cleanup(server.Close)

	scheme := runtime.NewScheme()
	assert.NoError(t, grafv1.AddToScheme(scheme))
	assert.NoError(t, corev1.AddToScheme(scheme))

	currentGrafana := &grafv1.Grafana{
		ObjectMeta: metav1.ObjectMeta{Name: "grafana", Namespace: "monitoring"},
		Status:     grafv1.GrafanaStatus{AdminURL: server.URL},
	}
	credentials := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: "grafana-admin-credentials", Namespace: "monitoring"},
		Data: map[string][]byte{
			"GF_SECURITY_ADMIN_USER":     []byte("admin"),
			"GF_SECURITY_ADMIN_PASSWORD": []byte("password"),
		},
	}
	legacyDatasource := &unstructured.Unstructured{}
	legacyDatasource.SetGroupVersionKind(legacyGrafanaDatasourceGVK)
	legacyDatasource.SetName("platform-monitoring-prometheus")
	legacyDatasource.SetNamespace("monitoring")
	kubeClient := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(currentGrafana, credentials, legacyDatasource).
		Build()
	reconciler := &GrafanaReconciler{
		ComponentReconciler: &utils.ComponentReconciler{
			Client: kubeClient,
			Scheme: scheme,
			Log:    utils.Logger("grafana_test"),
		},
	}
	platformMonitoring := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Name: "platformmonitoring", Namespace: "monitoring"},
		Spec: monv1.PlatformMonitoringSpec{
			Grafana: &monv1.Grafana{},
		},
	}
	datasource, err := grafanaDataSource(platformMonitoring, nil, nil, nil)
	assert.NoError(t, err)

	assert.NoError(t, reconciler.adoptExistingDatasourceUID(context.Background(), platformMonitoring, datasource))
	assert.Equal(t, "legacy-prometheus", datasource.Spec.CustomUID)
}

func TestAdoptExistingDatasourceUIDSkipsFreshInstall(t *testing.T) {
	scheme := runtime.NewScheme()
	assert.NoError(t, grafv1.AddToScheme(scheme))

	reconciler := &GrafanaReconciler{
		ComponentReconciler: &utils.ComponentReconciler{
			Client: fake.NewClientBuilder().WithScheme(scheme).Build(),
			Scheme: scheme,
			Log:    utils.Logger("grafana_test"),
		},
	}
	platformMonitoring := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Name: "platformmonitoring", Namespace: "monitoring"},
		Spec: monv1.PlatformMonitoringSpec{
			Grafana: &monv1.Grafana{},
		},
	}
	datasource, err := grafanaDataSource(platformMonitoring, nil, nil, nil)
	assert.NoError(t, err)

	assert.NoError(t, reconciler.adoptExistingDatasourceUID(context.Background(), platformMonitoring, datasource))
	assert.Empty(t, datasource.Spec.CustomUID)
}

func TestMigrateLegacyGrafanaResources(t *testing.T) {
	scheme := runtime.NewScheme()
	assert.NoError(t, monv1.AddToScheme(scheme))
	assert.NoError(t, grafv1.AddToScheme(scheme))
	assert.NoError(t, corev1.AddToScheme(scheme))
	assert.NoError(t, appsv1.AddToScheme(scheme))

	legacyUID := types.UID("legacy-grafana")
	currentUID := types.UID("current-grafana")
	platformMonitoringUID := types.UID("platform-monitoring")
	legacyOwner := metav1.OwnerReference{
		APIVersion: "integreatly.org/v1alpha1",
		Kind:       "Grafana",
		Name:       "grafana",
		UID:        legacyUID,
		Controller: ptr(true),
	}

	legacyGrafana := &unstructured.Unstructured{}
	legacyGrafana.SetGroupVersionKind(legacyGrafanaGVK)
	legacyGrafana.SetName("grafana")
	legacyGrafana.SetNamespace("monitoring")
	legacyGrafana.SetUID(legacyUID)
	currentGrafana := &grafv1.Grafana{
		ObjectMeta: metav1.ObjectMeta{Name: "grafana", Namespace: "monitoring", UID: currentUID},
	}
	platformMonitoring := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Name: "platformmonitoring", Namespace: "monitoring", UID: platformMonitoringUID},
	}
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:            "grafana-service",
			Namespace:       "monitoring",
			OwnerReferences: []metav1.OwnerReference{legacyOwner},
		},
	}
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:            "grafana-deployment",
			Namespace:       "monitoring",
			OwnerReferences: []metav1.OwnerReference{legacyOwner},
		},
		Spec: appsv1.DeploymentSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{ServiceAccountName: "grafana-serviceaccount"},
			},
		},
	}
	legacyServiceAccount := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{Name: "grafana-serviceaccount", Namespace: "monitoring"},
	}
	credentials := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:            "grafana-admin-credentials",
			Namespace:       "monitoring",
			OwnerReferences: []metav1.OwnerReference{legacyOwner},
		},
	}

	kubeClient := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(
			legacyGrafana,
			currentGrafana,
			platformMonitoring,
			service,
			deployment,
			legacyServiceAccount,
			credentials,
		).
		Build()
	reconciler := &GrafanaReconciler{
		ComponentReconciler: &utils.ComponentReconciler{Client: kubeClient, Scheme: scheme},
	}

	assert.NoError(t, reconciler.migrateLegacyGrafanaResources(
		context.Background(), platformMonitoring, currentGrafana))

	migratedService := &corev1.Service{}
	assert.NoError(t, kubeClient.Get(context.Background(), client.ObjectKeyFromObject(service), migratedService))
	assert.Equal(t, "grafana-operator", migratedService.Labels["app.kubernetes.io/managed-by"])
	assert.Equal(t, currentUID, metav1.GetControllerOf(migratedService).UID)

	migratedCredentials := &corev1.Secret{}
	assert.NoError(t, kubeClient.Get(context.Background(), client.ObjectKeyFromObject(credentials), migratedCredentials))
	assert.Equal(t, platformMonitoringUID, metav1.GetControllerOf(migratedCredentials).UID)

	migratedServiceAccount := &corev1.ServiceAccount{}
	assert.NoError(t, kubeClient.Get(
		context.Background(), client.ObjectKeyFromObject(legacyServiceAccount), migratedServiceAccount))
	if assert.NotNil(t, metav1.GetControllerOf(migratedServiceAccount)) {
		assert.Equal(t, platformMonitoringUID, metav1.GetControllerOf(migratedServiceAccount).UID)
	}

	err := kubeClient.Get(context.Background(), client.ObjectKeyFromObject(legacyGrafana), &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "integreatly.org/v1alpha1",
			"kind":       "Grafana",
		},
	})
	assert.True(t, apierrors.IsNotFound(err))
}

func TestMigrateLegacyGrafanaResourcesSkipsFreshInstall(t *testing.T) {
	scheme := runtime.NewScheme()
	assert.NoError(t, grafv1.AddToScheme(scheme))

	kubeClient := fake.NewClientBuilder().WithScheme(scheme).Build()
	reconciler := &GrafanaReconciler{
		ComponentReconciler: &utils.ComponentReconciler{
			Client: kubeClient,
			Scheme: scheme,
		},
	}
	platformMonitoring := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "platformmonitoring",
			Namespace: "monitoring",
		},
	}
	currentGrafana := &grafv1.Grafana{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "grafana",
			Namespace: "monitoring",
		},
	}

	assert.NoError(t, reconciler.migrateLegacyGrafanaResources(
		context.Background(),
		platformMonitoring,
		currentGrafana,
	))
}

func ptr[T any](value T) *T {
	return &value
}
