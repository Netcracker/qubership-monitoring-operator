package prometheus_operator

import (
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils/labelsassert"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/discovery/fake"
	"k8s.io/utils/ptr"
)

var (
	cr              *monv1.PlatformMonitoring
	labelKey        = "label.key"
	labelValue      = "label-value"
	annotationKey   = "annotation.key"
	annotationValue = "annotation-value"
)

func TestNewPrometheusOperatorReconcilerStoresDiscoveryClient(t *testing.T) {
	discoveryClient := &fake.FakeDiscovery{}

	reconciler := NewPrometheusOperatorReconciler(nil, nil, discoveryClient)

	assert.Same(t, discoveryClient, reconciler.Dc)
}

func TestPrometheusOperatorManifests(t *testing.T) {
	cr = &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
		Spec: monv1.PlatformMonitoringSpec{
			Prometheus: &monv1.Prometheus{
				Operator: monv1.PrometheusOperator{
					Annotations: map[string]string{annotationKey: annotationValue},
					Labels:      map[string]string{labelKey: labelValue},
				},
			},
		},
	}
	t.Run("Test Deployment manifest", func(t *testing.T) {
		m, err := prometheusOperatorDeployment(cr, false)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Deployment manifest should not be empty")
		assert.NotNil(t, m.GetLabels())
		assert.Equal(t, labelValue, m.GetLabels()[labelKey])
		assert.NotNil(t, m.Spec.Template.Labels)
		assert.Equal(t, labelValue, m.Spec.Template.Labels[labelKey])
		assert.NotNil(t, m.GetAnnotations())
		assert.Equal(t, annotationValue, m.GetAnnotations()[annotationKey])
		assert.NotNil(t, m.Spec.Template.Annotations)
		assert.Equal(t, annotationValue, m.Spec.Template.Annotations[annotationKey])
		assertPrometheusOperatorHardening(t, m, false)
	})
	cr = &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
		Spec: monv1.PlatformMonitoringSpec{
			Prometheus: &monv1.Prometheus{
				Operator: monv1.PrometheusOperator{},
			},
		},
	}
	t.Run("Test Deployment manifest with nil labels and annotation", func(t *testing.T) {
		m, err := prometheusOperatorDeployment(cr, false)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Deployment manifest should not be empty")
		assert.NotNil(t, m.GetLabels())
		assert.NotNil(t, m.Spec.Template.Labels)
		assert.Nil(t, m.GetAnnotations())
		assert.Nil(t, m.Spec.Template.Annotations)
	})
	t.Run("Test OpenShift Deployment security context", func(t *testing.T) {
		m, err := prometheusOperatorDeployment(cr, true)
		require.NoError(t, err)
		assertPrometheusOperatorHardening(t, m, true)
	})
	t.Run("Test configured IDs are preserved", func(t *testing.T) {
		configuredCR := cr.DeepCopy()
		configuredCR.Spec.Prometheus.Operator.SecurityContext = &monv1.SecurityContext{
			RunAsUser:  ptr.To(int64(3000)),
			RunAsGroup: ptr.To(int64(3000)),
			FSGroup:    ptr.To(int64(3000)),
		}

		m, err := prometheusOperatorDeployment(configuredCR, false)
		require.NoError(t, err)
		assert.Equal(t, int64(3000), *m.Spec.Template.Spec.SecurityContext.RunAsUser)
		assert.Equal(t, int64(3000), *m.Spec.Template.Spec.SecurityContext.RunAsGroup)
		assert.Equal(t, int64(3000), *m.Spec.Template.Spec.SecurityContext.FSGroup)
		assertPrometheusOperatorHardening(t, m, false)
	})
	t.Run("Test Role manifest", func(t *testing.T) {
		m, err := prometheusOperatorRole(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Role manifest should not be empty")
	})
	t.Run("Test RoleBinding manifest", func(t *testing.T) {
		m, err := prometheusOperatorRoleBinding(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "RoleBinding manifest should not be empty")
	})
	t.Run("Test ServiceAccount manifest", func(t *testing.T) {
		crWithSALabels := &monv1.PlatformMonitoring{
			ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"},
			Spec: monv1.PlatformMonitoringSpec{
				Prometheus: &monv1.Prometheus{
					Operator: monv1.PrometheusOperator{
						ServiceAccount: &monv1.EmbeddedObjectMetadata{
							Labels: map[string]string{labelKey: labelValue},
						},
					},
				},
			},
		}
		m, err := prometheusOperatorServiceAccount(crWithSALabels)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "ServiceAccount manifest should not be empty")
		assert.NotNil(t, m.GetLabels())
		assert.Equal(t, labelValue, m.GetLabels()[labelKey], "ServiceAccount.Labels should be merged")
	})
	t.Run("Test ClusterRole manifest", func(t *testing.T) {
		m, err := prometheusOperatorClusterRole(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "ClusterRole manifest should not be empty")
	})
	t.Run("Test ClusterRoleBinding manifest", func(t *testing.T) {
		m, err := prometheusOperatorClusterRoleBinding(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "ClusterRoleBinding manifest should not be empty")
	})
	t.Run("Test Service manifest", func(t *testing.T) {
		m, err := prometheusOperatorService(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Service manifest should not be empty")
	})
	t.Run("Test PodMonitor manifest", func(t *testing.T) {
		crWithLabels := &monv1.PlatformMonitoring{
			ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring", Labels: map[string]string{labelKey: labelValue}},
			Spec: monv1.PlatformMonitoringSpec{
				Prometheus: &monv1.Prometheus{
					Operator: monv1.PrometheusOperator{
						Annotations: map[string]string{annotationKey: annotationValue},
						Labels:      map[string]string{labelKey: labelValue},
					},
				},
			},
		}
		m, err := prometheusOperatorPodMonitor(crWithLabels)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "PodMonitor manifest should not be empty")
		labelsassert.AssertCRLabels(t, m.GetLabels(), utils.PrometheusOperatorComponentName, "prometheus-operator", map[string]string{labelKey: labelValue})
	})
}

func assertPrometheusOperatorHardening(t *testing.T, deployment *appsv1.Deployment, isOpenShift bool) {
	t.Helper()

	podSecurityContext := deployment.Spec.Template.Spec.SecurityContext
	require.NotNil(t, podSecurityContext)
	require.NotNil(t, podSecurityContext.RunAsNonRoot)
	assert.True(t, *podSecurityContext.RunAsNonRoot)
	require.NotNil(t, podSecurityContext.SeccompProfile)
	assert.Equal(t, corev1.SeccompProfileTypeRuntimeDefault, podSecurityContext.SeccompProfile.Type)
	if isOpenShift {
		assert.Nil(t, podSecurityContext.RunAsUser)
		assert.Nil(t, podSecurityContext.RunAsGroup)
		assert.Nil(t, podSecurityContext.FSGroup)
	} else {
		require.NotNil(t, podSecurityContext.RunAsUser)
		require.NotNil(t, podSecurityContext.RunAsGroup)
		require.NotNil(t, podSecurityContext.FSGroup)
	}

	require.Len(t, deployment.Spec.Template.Spec.Containers, 1)
	container := deployment.Spec.Template.Spec.Containers[0]
	require.NotNil(t, container.SecurityContext)
	require.NotNil(t, container.SecurityContext.AllowPrivilegeEscalation)
	assert.False(t, *container.SecurityContext.AllowPrivilegeEscalation)
	require.NotNil(t, container.SecurityContext.ReadOnlyRootFilesystem)
	assert.True(t, *container.SecurityContext.ReadOnlyRootFilesystem)
	require.NotNil(t, container.SecurityContext.Capabilities)
	assert.Equal(t, []corev1.Capability{"ALL"}, container.SecurityContext.Capabilities.Drop)
	assert.Contains(t, container.VolumeMounts, utils.TmpVolumeMount())

	tmpVolumes := 0
	for _, volume := range deployment.Spec.Template.Spec.Volumes {
		if volume.Name == "tmp" {
			tmpVolumes++
			require.NotNil(t, volume.EmptyDir)
			require.NotNil(t, volume.EmptyDir.SizeLimit)
			assert.Equal(t, "100Mi", volume.EmptyDir.SizeLimit.String())
		}
	}
	assert.Equal(t, 1, tmpVolumes)
}
