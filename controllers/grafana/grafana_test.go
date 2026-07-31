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
	"github.com/go-logr/logr"
	grafv1 "github.com/grafana/grafana-operator/v5/api/v1beta1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	kubernetesfake "k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestAddGrafanaExtraVarsResourceVersions(t *testing.T) {
	manifest, err := grafana(&monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"},
		Spec: monv1.PlatformMonitoringSpec{Grafana: &monv1.Grafana{
			Annotations: map[string]string{"example.com/user-annotation": "retained"},
		}},
	})
	assert.NoError(t, err)

	reconciler := &GrafanaReconciler{KubeClient: kubernetesfake.NewSimpleClientset(
		&corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{
			Name: "grafana-extra-vars", Namespace: "monitoring", ResourceVersion: "config-version",
		}},
		&corev1.Secret{ObjectMeta: metav1.ObjectMeta{
			Name: "grafana-extra-vars-secret", Namespace: "monitoring", ResourceVersion: "secret-version",
		}},
	)}

	err = reconciler.addGrafanaExtraVarsResourceVersions(context.Background(), "monitoring", manifest, nil)
	assert.NoError(t, err)
	assert.Equal(t, "config-version",
		manifest.Spec.Deployment.Spec.Template.Annotations[grafanaExtraVarsConfigMapResourceVersionAnnotation])
	assert.Equal(t, "secret-version",
		manifest.Spec.Deployment.Spec.Template.Annotations[grafanaExtraVarsSecretResourceVersionAnnotation])
	assert.Equal(t, "retained", manifest.Spec.Deployment.Spec.Template.Annotations["example.com/user-annotation"])
}

func TestAddGrafanaExtraVarsResourceVersionsWhenResourcesAreMissing(t *testing.T) {
	manifest, err := grafana(&monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"},
		Spec: monv1.PlatformMonitoringSpec{Grafana: &monv1.Grafana{
			Annotations: map[string]string{"example.com/user-annotation": "retained"},
		}},
	})
	assert.NoError(t, err)

	reconciler := &GrafanaReconciler{KubeClient: kubernetesfake.NewSimpleClientset()}

	err = reconciler.addGrafanaExtraVarsResourceVersions(context.Background(), "monitoring", manifest, nil)
	assert.NoError(t, err)
	assert.NotContains(t, manifest.Spec.Deployment.Spec.Template.Annotations,
		grafanaExtraVarsConfigMapResourceVersionAnnotation)
	assert.NotContains(t, manifest.Spec.Deployment.Spec.Template.Annotations,
		grafanaExtraVarsSecretResourceVersionAnnotation)
	assert.Equal(t, "retained", manifest.Spec.Deployment.Spec.Template.Annotations["example.com/user-annotation"])
}

func TestAddGrafanaExtraVarsResourceVersionsKeepsAnnotationsNilWhenResourcesAreMissing(t *testing.T) {
	manifest, err := grafana(&monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"},
		Spec:       monv1.PlatformMonitoringSpec{Grafana: &monv1.Grafana{}},
	})
	require.NoError(t, err)
	require.Nil(t, manifest.Spec.Deployment.Spec.Template.Annotations)
	reconciler := &GrafanaReconciler{KubeClient: kubernetesfake.NewSimpleClientset()}

	err = reconciler.addGrafanaExtraVarsResourceVersions(context.Background(), "monitoring", manifest, nil)

	require.NoError(t, err)
	assert.Nil(t, manifest.Spec.Deployment.Spec.Template.Annotations)
}

func TestAddGrafanaExtraVarsResourceVersionsPreservesVersionForMissingResource(t *testing.T) {
	manifest, err := grafana(&monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"},
		Spec:       monv1.PlatformMonitoringSpec{Grafana: &monv1.Grafana{}},
	})
	assert.NoError(t, err)

	reconciler := &GrafanaReconciler{KubeClient: kubernetesfake.NewSimpleClientset(
		&corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{
			Name: "grafana-extra-vars", Namespace: "monitoring", ResourceVersion: "new-config-version",
		}},
	)}
	existingAnnotations := map[string]string{
		grafanaExtraVarsConfigMapResourceVersionAnnotation: "old-config-version",
		grafanaExtraVarsSecretResourceVersionAnnotation:    "old-secret-version",
	}

	err = reconciler.addGrafanaExtraVarsResourceVersions(
		context.Background(), "monitoring", manifest, existingAnnotations)
	assert.NoError(t, err)
	assert.Equal(t, "new-config-version",
		manifest.Spec.Deployment.Spec.Template.Annotations[grafanaExtraVarsConfigMapResourceVersionAnnotation])
	assert.Equal(t, "old-secret-version",
		manifest.Spec.Deployment.Spec.Template.Annotations[grafanaExtraVarsSecretResourceVersionAnnotation])
}

func TestAddGrafanaExtraVarsResourceVersionsPreservesMissingConfigMapVersion(t *testing.T) {
	manifest, err := grafana(&monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"},
		Spec:       monv1.PlatformMonitoringSpec{Grafana: &monv1.Grafana{}},
	})
	assert.NoError(t, err)

	reconciler := &GrafanaReconciler{KubeClient: kubernetesfake.NewSimpleClientset(
		&corev1.Secret{ObjectMeta: metav1.ObjectMeta{
			Name: "grafana-extra-vars-secret", Namespace: "monitoring", ResourceVersion: "new-secret-version",
		}},
	)}
	existingAnnotations := map[string]string{
		grafanaExtraVarsConfigMapResourceVersionAnnotation: "old-config-version",
	}

	err = reconciler.addGrafanaExtraVarsResourceVersions(
		context.Background(), "monitoring", manifest, existingAnnotations)
	assert.NoError(t, err)
	assert.Equal(t, "old-config-version",
		manifest.Spec.Deployment.Spec.Template.Annotations[grafanaExtraVarsConfigMapResourceVersionAnnotation])
	assert.Equal(t, "new-secret-version",
		manifest.Spec.Deployment.Spec.Template.Annotations[grafanaExtraVarsSecretResourceVersionAnnotation])
}

func TestAddGrafanaExtraVarsResourceVersionsInitializesAnnotations(t *testing.T) {
	manifest, err := grafana(&monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"},
		Spec:       monv1.PlatformMonitoringSpec{Grafana: &monv1.Grafana{}},
	})
	assert.NoError(t, err)
	manifest.Spec.Deployment.Spec.Template.Annotations = nil
	reconciler := &GrafanaReconciler{KubeClient: kubernetesfake.NewSimpleClientset(
		&corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{
			Name: "grafana-extra-vars", Namespace: "monitoring", ResourceVersion: "config-version",
		}},
		&corev1.Secret{ObjectMeta: metav1.ObjectMeta{
			Name: "grafana-extra-vars-secret", Namespace: "monitoring", ResourceVersion: "secret-version",
		}},
	)}

	err = reconciler.addGrafanaExtraVarsResourceVersions(context.Background(), "monitoring", manifest, nil)
	assert.NoError(t, err)
	assert.Equal(t, "config-version",
		manifest.Spec.Deployment.Spec.Template.Annotations[grafanaExtraVarsConfigMapResourceVersionAnnotation])
	assert.Equal(t, "secret-version",
		manifest.Spec.Deployment.Spec.Template.Annotations[grafanaExtraVarsSecretResourceVersionAnnotation])
}

func TestGrafanaPodTemplateAnnotationsHandlesMissingTemplate(t *testing.T) {
	manifest, err := grafana(&monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"},
		Spec:       monv1.PlatformMonitoringSpec{Grafana: &monv1.Grafana{}},
	})
	assert.NoError(t, err)
	manifest.Spec.Deployment.Spec.Template = nil

	assert.Nil(t, grafanaPodTemplateAnnotations(manifest))
	manifest, err = grafana(&monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"},
		Spec:       monv1.PlatformMonitoringSpec{Grafana: &monv1.Grafana{}},
	})
	assert.NoError(t, err)
	assert.Nil(t, grafanaPodTemplateAnnotations(manifest))
	manifest.Spec.Deployment.Spec.Template.Annotations = map[string]string{"custom": "value"}
	assert.Equal(t, map[string]string{"custom": "value"}, grafanaPodTemplateAnnotations(manifest))
}

func TestHandleGrafanaCreatesResourceWithoutExtraVarsResources(t *testing.T) {
	testScheme := runtime.NewScheme()
	assert.NoError(t, monv1.AddToScheme(testScheme))
	assert.NoError(t, grafv1.AddToScheme(testScheme))
	controllerClient := fake.NewClientBuilder().WithScheme(testScheme).Build()
	reconciler := &GrafanaReconciler{
		KubeClient: kubernetesfake.NewSimpleClientset(),
		ComponentReconciler: &utils.ComponentReconciler{
			Client: controllerClient,
			Scheme: testScheme,
			Log:    logr.Discard(),
		},
	}
	platformMonitoring := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Name: "monitoring", Namespace: "monitoring"},
		Spec:       monv1.PlatformMonitoringSpec{Grafana: &monv1.Grafana{}},
	}

	err := reconciler.handleGrafana(platformMonitoring)
	assert.NoError(t, err)

	created := &grafv1.Grafana{}
	err = controllerClient.Get(context.Background(), types.NamespacedName{
		Name: "grafana", Namespace: "monitoring",
	}, created)
	assert.NoError(t, err)
	assert.NotNil(t, created.Spec.Deployment)
}

func TestHandleGrafanaReturnsExtraVarsAPIErrorForExistingResource(t *testing.T) {
	testScheme := runtime.NewScheme()
	assert.NoError(t, monv1.AddToScheme(testScheme))
	assert.NoError(t, grafv1.AddToScheme(testScheme))
	platformMonitoring := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Name: "monitoring", Namespace: "monitoring"},
		Spec:       monv1.PlatformMonitoringSpec{Grafana: &monv1.Grafana{}},
	}
	existing, err := grafana(platformMonitoring)
	assert.NoError(t, err)
	existing.Spec.Deployment.Spec.Template = nil
	controllerClient := fake.NewClientBuilder().WithScheme(testScheme).WithObjects(existing).Build()
	kubeClient := kubernetesfake.NewSimpleClientset()
	kubeClient.PrependReactor("get", "configmaps", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, apierrors.NewForbidden(
			schema.GroupResource{Resource: "configmaps"}, "grafana-extra-vars", assert.AnError)
	})
	reconciler := &GrafanaReconciler{
		KubeClient: kubeClient,
		ComponentReconciler: &utils.ComponentReconciler{
			Client: controllerClient,
			Scheme: testScheme,
			Log:    logr.Discard(),
		},
	}

	err = reconciler.handleGrafana(platformMonitoring)
	assert.ErrorContains(t, err, "cannot get Grafana extra-vars ConfigMap")
}

func TestHandleGrafanaReturnsClientError(t *testing.T) {
	reconciler := &GrafanaReconciler{
		KubeClient: kubernetesfake.NewSimpleClientset(),
		ComponentReconciler: &utils.ComponentReconciler{
			Client: fake.NewClientBuilder().WithScheme(runtime.NewScheme()).Build(),
			Scheme: runtime.NewScheme(),
			Log:    logr.Discard(),
		},
	}
	platformMonitoring := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Name: "monitoring", Namespace: "monitoring"},
		Spec:       monv1.PlatformMonitoringSpec{Grafana: &monv1.Grafana{}},
	}

	err := reconciler.handleGrafana(platformMonitoring)
	assert.Error(t, err)
}

func TestAddGrafanaExtraVarsResourceVersionsReturnsSecretAPIError(t *testing.T) {
	manifest, err := grafana(&monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"},
		Spec:       monv1.PlatformMonitoringSpec{Grafana: &monv1.Grafana{}},
	})
	assert.NoError(t, err)
	kubeClient := kubernetesfake.NewSimpleClientset(&corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{
		Name: "grafana-extra-vars", Namespace: "monitoring",
	}})
	kubeClient.PrependReactor("get", "secrets", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, apierrors.NewForbidden(
			schema.GroupResource{Resource: "secrets"}, "grafana-extra-vars-secret", assert.AnError)
	})
	reconciler := &GrafanaReconciler{KubeClient: kubeClient}

	err = reconciler.addGrafanaExtraVarsResourceVersions(context.Background(), "monitoring", manifest, nil)
	assert.ErrorContains(t, err, "cannot get Grafana extra-vars Secret")
}

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
		m, err := grafana(cr, false)
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
		assertGrafanaHardening(t, m, false)
	})
	t.Run("Test Grafana manifest preserves cleanup label", func(t *testing.T) {
		cr.Spec.Grafana.Labels["app.kubernetes.io/managed-by"] = "custom-manager"
		cr.Spec.Grafana.Labels["app.kubernetes.io/managed-by-operator"] = "custom-manager"

		m, err := grafana(cr, false)
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
	//	m, err := grafana(cr, false)
	//	...
	//})
	t.Run("Test OpenShift Grafana security context", func(t *testing.T) {
		m, err := grafana(cr, true)
		require.NoError(t, err)
		assertGrafanaHardening(t, m, true)
	})
	t.Run("Test configured IDs are preserved", func(t *testing.T) {
		configuredCR := cr.DeepCopy()
		configuredCR.Spec.Grafana.SecurityContext = &monv1.SecurityContext{
			RunAsUser:  ptr(int64(3000)),
			RunAsGroup: ptr(int64(3000)),
			FSGroup:    ptr(int64(3000)),
		}

		m, err := grafana(configuredCR, false)
		require.NoError(t, err)
		podSpec := grafanaPodSpec(t, m)
		require.NotNil(t, podSpec.SecurityContext)
		assert.Equal(t, int64(3000), *podSpec.SecurityContext.RunAsUser)
		assert.Equal(t, int64(3000), *podSpec.SecurityContext.RunAsGroup)
		assert.Equal(t, int64(3000), *podSpec.SecurityContext.FSGroup)
		assertGrafanaHardening(t, m, false)
	})
	t.Run("Test hardening is merged into the generated Deployment", func(t *testing.T) {
		m, err := grafana(cr, false)
		require.NoError(t, err)

		deployment := &appsv1.Deployment{
			Spec: appsv1.DeploymentSpec{
				Template: corev1.PodTemplateSpec{
					Spec: corev1.PodSpec{
						Volumes: []corev1.Volume{{Name: "grafana-data"}},
						Containers: []corev1.Container{{
							Name:  "grafana",
							Image: "grafana/grafana:latest",
						}},
					},
				},
			},
		}
		require.NoError(t, grafv1.Merge(deployment, m.Spec.Deployment))

		require.NotNil(t, deployment.Spec.Template.Spec.SecurityContext)
		assert.True(t, *deployment.Spec.Template.Spec.SecurityContext.RunAsNonRoot)
		require.Len(t, deployment.Spec.Template.Spec.Containers, 1)
		container := deployment.Spec.Template.Spec.Containers[0]
		assert.Equal(t, "grafana/grafana:latest", container.Image)
		assert.True(t, *container.SecurityContext.ReadOnlyRootFilesystem)
		assert.Contains(t, container.VolumeMounts, utils.TmpVolumeMount())
		assert.Contains(t, deployment.Spec.Template.Spec.Volumes, corev1.Volume{Name: "grafana-data"})
		assert.Contains(t, deployment.Spec.Template.Spec.Volumes, utils.TmpVolume("100Mi"))
	})
	t.Run("Test GrafanaDatasource manifest", func(t *testing.T) {
		m, err := grafanaDataSource(cr, nil, nil, nil)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "GrafanaDatasource manifest should not be empty")
		assert.Equal(t, grafanaCleanupLabelValue, m.GetLabels()[grafanaCleanupLabelKey])
		assert.Equal(t, 10*time.Minute, m.Spec.ResyncPeriod.Duration)
	})
	t.Run("Test GrafanaPromxyDatasource manifest", func(t *testing.T) {
		m, err := grafanaPromxyDataSource(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "GrafanaPromxyDatasource manifest should not be empty")
		assert.Equal(t, "platform-monitoring-promxy", m.GetName())
		assert.Equal(t, grafanaCleanupLabelValue, m.GetLabels()[grafanaCleanupLabelKey])
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

func TestGrafanaPodTemplateAnnotations(t *testing.T) {
	t.Run("keeps annotations nil when none are configured", func(t *testing.T) {
		manifest, err := grafana(grafanaComparisonPlatformMonitoring(nil))
		require.NoError(t, err)

		assert.Nil(t, manifest.Spec.Deployment.Spec.Template.Annotations)
	})

	t.Run("preserves configured annotations", func(t *testing.T) {
		annotations := map[string]string{"example.com/key": "value"}
		manifest, err := grafana(grafanaComparisonPlatformMonitoring(annotations))
		require.NoError(t, err)

		assert.Equal(t, annotations, manifest.Spec.Deployment.Spec.Template.Annotations)
	})
}

func TestGrafanaManifestPreservesDataStorage(t *testing.T) {
	monitoring := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"},
		Spec: monv1.PlatformMonitoringSpec{
			Grafana: &monv1.Grafana{
				DataStorage: &runtime.RawExtension{
					Raw: []byte(`{"accessModes":["ReadWriteOnce"],"annotations":{"storage.example.com/owner":"monitoring"},"class":"standard","labels":{"app.example.com/tier":"data"},"size":"2Gi","volumeName":"grafana-static-pv"}`),
				},
			},
		},
	}

	manifest, err := grafana(monitoring, false)
	if err != nil {
		t.Fatal(err)
	}

	if assert.NotNil(t, manifest.Spec.PersistentVolumeClaim) &&
		assert.NotNil(t, manifest.Spec.PersistentVolumeClaim.Spec) {
		assert.Equal(t, map[string]string{"storage.example.com/owner": "monitoring"},
			manifest.Spec.PersistentVolumeClaim.ObjectMeta.Annotations)
		assert.Equal(t, map[string]string{"app.example.com/tier": "data"},
			manifest.Spec.PersistentVolumeClaim.ObjectMeta.Labels)
		assert.Equal(t, []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			manifest.Spec.PersistentVolumeClaim.Spec.AccessModes)
		assert.Equal(t, "standard", *manifest.Spec.PersistentVolumeClaim.Spec.StorageClassName)
		assert.Equal(t, "2Gi",
			manifest.Spec.PersistentVolumeClaim.Spec.Resources.Requests.Storage().String())
		assert.Equal(t, "grafana-static-pv", manifest.Spec.PersistentVolumeClaim.Spec.VolumeName)
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

func assertGrafanaHardening(t *testing.T, manifest *grafv1.Grafana, isOpenShift bool) {
	t.Helper()

	podSpec := grafanaPodSpec(t, manifest)
	require.NotNil(t, podSpec.SecurityContext)
	require.NotNil(t, podSpec.SecurityContext.RunAsNonRoot)
	assert.True(t, *podSpec.SecurityContext.RunAsNonRoot)
	require.NotNil(t, podSpec.SecurityContext.SeccompProfile)
	assert.Equal(t, corev1.SeccompProfileTypeRuntimeDefault, podSpec.SecurityContext.SeccompProfile.Type)
	if isOpenShift {
		assert.Nil(t, podSpec.SecurityContext.RunAsUser)
		assert.Nil(t, podSpec.SecurityContext.RunAsGroup)
		assert.Nil(t, podSpec.SecurityContext.FSGroup)
	} else {
		require.NotNil(t, podSpec.SecurityContext.RunAsUser)
		require.NotNil(t, podSpec.SecurityContext.RunAsGroup)
		require.NotNil(t, podSpec.SecurityContext.FSGroup)
	}

	require.NotEmpty(t, podSpec.Containers)
	container := podSpec.Containers[0]
	require.NotNil(t, container.SecurityContext)
	require.NotNil(t, container.SecurityContext.AllowPrivilegeEscalation)
	assert.False(t, *container.SecurityContext.AllowPrivilegeEscalation)
	require.NotNil(t, container.SecurityContext.ReadOnlyRootFilesystem)
	assert.True(t, *container.SecurityContext.ReadOnlyRootFilesystem)
	require.NotNil(t, container.SecurityContext.Capabilities)
	assert.Equal(t, []corev1.Capability{"ALL"}, container.SecurityContext.Capabilities.Drop)

	tmpMounts := 0
	for _, volumeMount := range container.VolumeMounts {
		if volumeMount.Name == "tmp" && volumeMount.MountPath == "/tmp" {
			tmpMounts++
		}
	}
	assert.Equal(t, 1, tmpMounts)

	tmpVolumes := 0
	for _, volume := range podSpec.Volumes {
		if volume.Name == "tmp" {
			tmpVolumes++
			require.NotNil(t, volume.EmptyDir)
			require.NotNil(t, volume.EmptyDir.SizeLimit)
			assert.Equal(t, "100Mi", volume.EmptyDir.SizeLimit.String())
		}
	}
	assert.Equal(t, 1, tmpVolumes)
}

func grafanaPodSpec(t *testing.T, manifest *grafv1.Grafana) *grafv1.DeploymentV1PodSpec {
	t.Helper()
	require.NotNil(t, manifest.Spec.Deployment)
	require.NotNil(t, manifest.Spec.Deployment.Spec.Template)
	require.NotNil(t, manifest.Spec.Deployment.Spec.Template.Spec)
	return manifest.Spec.Deployment.Spec.Template.Spec
}
