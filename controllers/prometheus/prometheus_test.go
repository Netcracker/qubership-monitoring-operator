package prometheus

import (
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils/labelsassert"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

var (
	cr              *monv1.PlatformMonitoring
	labelKey        = "label.key"
	labelValue      = "label-value"
	annotationKey   = "annotation.key"
	annotationValue = "annotation-value"
)

func TestPrometheusManifests(t *testing.T) {
	cr = &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
		Spec: monv1.PlatformMonitoringSpec{
			Prometheus: &monv1.Prometheus{
				Annotations: map[string]string{annotationKey: annotationValue},
				Labels:      map[string]string{labelKey: labelValue},
			},
		},
	}
	t.Run("Test Prometheus manifest", func(t *testing.T) {
		m, err := prometheus(cr, false)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Prometheus manifest should not be empty")
		assert.NotNil(t, m.Spec.PodMetadata.Labels)
		assert.Equal(t, labelValue, m.Spec.PodMetadata.Labels[labelKey])
		assert.NotNil(t, m.Spec.PodMetadata.Annotations)
		assert.Equal(t, annotationValue, m.Spec.PodMetadata.Annotations[annotationKey])
		assertPrometheusHardening(t, m, false)
	})
	cr = &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
		Spec: monv1.PlatformMonitoringSpec{
			Prometheus: &monv1.Prometheus{},
		},
	}
	t.Run("Test Prometheus manifest with nil labels and annotation", func(t *testing.T) {
		m, err := prometheus(cr, false)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Prometheus manifest should not be empty")
		assert.NotNil(t, m.GetLabels())
		assert.Nil(t, m.GetAnnotations())
	})
	t.Run("Test OpenShift security context", func(t *testing.T) {
		m, err := prometheus(cr, true)
		require.NoError(t, err)
		assertPrometheusHardening(t, m, true)
	})
	t.Run("Test configured IDs and sidecar settings", func(t *testing.T) {
		configuredCR := cr.DeepCopy()
		configuredCR.Spec.Prometheus.SecurityContext = &monv1.SecurityContext{
			RunAsUser:  ptr.To(int64(3000)),
			RunAsGroup: ptr.To(int64(3001)),
			FSGroup:    ptr.To(int64(3002)),
		}
		configuredCR.Spec.Prometheus.Volumes = []corev1.Volume{{
			Name: "tmp",
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		}}
		configuredCR.Spec.Prometheus.VolumeMounts = []corev1.VolumeMount{{Name: "custom", MountPath: "/custom-new"}}
		configuredCR.Spec.Prometheus.Containers = []corev1.Container{{
			Name: "sidecar",
			SecurityContext: &corev1.SecurityContext{
				AllowPrivilegeEscalation: ptr.To(true),
				ReadOnlyRootFilesystem:   ptr.To(false),
			},
			VolumeMounts: []corev1.VolumeMount{{Name: "other-tmp", MountPath: "/tmp"}},
		}, {
			Name:         "prometheus",
			VolumeMounts: []corev1.VolumeMount{{Name: "custom", MountPath: "/custom-old"}},
		}}

		m, err := prometheus(configuredCR, false)
		require.NoError(t, err)
		assert.Equal(t, int64(3000), *m.Spec.SecurityContext.RunAsUser)
		assert.Equal(t, int64(3001), *m.Spec.SecurityContext.RunAsGroup)
		assert.Equal(t, int64(3002), *m.Spec.SecurityContext.FSGroup)
		assertPrometheusHardening(t, m, false)
		assertPrometheusContainerHardening(t, findPrometheusContainer(t, m, "sidecar"))
		assert.Contains(t, findPrometheusContainer(t, m, "prometheus").VolumeMounts,
			corev1.VolumeMount{Name: "custom", MountPath: "/custom-new"})
		assert.NotContains(t, findPrometheusContainer(t, m, "prometheus").VolumeMounts,
			corev1.VolumeMount{Name: "custom", MountPath: "/custom-old"})
		assert.Nil(t, configuredCR.Spec.Prometheus.Volumes[0].EmptyDir.SizeLimit,
			"the source CR must not be mutated")
		assert.Equal(t, ptr.To(true), configuredCR.Spec.Prometheus.Containers[0].SecurityContext.AllowPrivilegeEscalation,
			"the source CR must not be mutated")

		openShiftManifest, err := prometheus(configuredCR, true)
		require.NoError(t, err)
		assert.Nil(t, openShiftManifest.Spec.SecurityContext.RunAsUser)
		assert.Nil(t, openShiftManifest.Spec.SecurityContext.RunAsGroup)
		assert.Nil(t, openShiftManifest.Spec.SecurityContext.FSGroup)
	})
	cr = &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
	}
	t.Run("Test ServiceAccount manifest", func(t *testing.T) {
		crWithSALabels := &monv1.PlatformMonitoring{
			ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"},
			Spec: monv1.PlatformMonitoringSpec{
				Prometheus: &monv1.Prometheus{
					ServiceAccount: &monv1.EmbeddedObjectMetadata{
						Labels: map[string]string{labelKey: labelValue},
					},
				},
			},
		}
		m, err := prometheusServiceAccount(crWithSALabels)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "ServiceAccount manifest should not be empty")
		assert.NotNil(t, m.GetLabels())
		assert.Equal(t, labelValue, m.GetLabels()[labelKey], "ServiceAccount.Labels should be merged")
	})
	t.Run("Test ClusterRole manifest", func(t *testing.T) {
		m, err := prometheusClusterRole(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "ClusterRole manifest should not be empty")
	})
	t.Run("Test ClusterRoleBinding manifest", func(t *testing.T) {
		m, err := prometheusClusterRoleBinding(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "ClusterRoleBinding manifest should not be empty")
	})
	t.Run("Test Ingress v1 manifest", func(t *testing.T) {
		m, err := prometheusIngressV1(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Ingress v1 manifest should not be empty")
	})
	t.Run("Test PodMonitor manifest", func(t *testing.T) {
		crWithLabels := &monv1.PlatformMonitoring{
			ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring", Labels: map[string]string{labelKey: labelValue}},
			Spec:       monv1.PlatformMonitoringSpec{Prometheus: &monv1.Prometheus{}},
		}
		m, err := prometheusPodMonitor(crWithLabels)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "PodMonitor manifest should not be empty")
		labelsassert.AssertCRLabels(t, m.GetLabels(), utils.PrometheusComponentName, "prometheus-operator", map[string]string{labelKey: labelValue})
	})
}

func assertPrometheusHardening(t *testing.T, prometheus *promv1.Prometheus, isOpenShift bool) {
	t.Helper()

	podSecurityContext := prometheus.Spec.SecurityContext
	require.NotNil(t, podSecurityContext)
	assert.Equal(t, ptr.To(true), podSecurityContext.RunAsNonRoot)
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

	tmpVolumes := 0
	for _, volume := range prometheus.Spec.Volumes {
		if volume.Name == "tmp" {
			tmpVolumes++
			require.NotNil(t, volume.EmptyDir)
			require.NotNil(t, volume.EmptyDir.SizeLimit)
			assert.Equal(t, resource.MustParse("100Mi"), *volume.EmptyDir.SizeLimit)
		}
	}
	assert.Equal(t, 1, tmpVolumes)

	assertPrometheusContainerHardening(t, findPrometheusContainer(t, prometheus, "prometheus"))
	assertPrometheusContainerHardening(t, findPrometheusContainer(t, prometheus, "config-reloader"))
}

func assertPrometheusContainerHardening(t *testing.T, container corev1.Container) {
	t.Helper()
	require.NotNil(t, container.SecurityContext)
	assert.Equal(t, ptr.To(false), container.SecurityContext.AllowPrivilegeEscalation)
	assert.Equal(t, ptr.To(true), container.SecurityContext.ReadOnlyRootFilesystem)
	require.NotNil(t, container.SecurityContext.Capabilities)
	assert.Equal(t, []corev1.Capability{"ALL"}, container.SecurityContext.Capabilities.Drop)
	assert.Contains(t, container.VolumeMounts, utils.TmpVolumeMount())
}

func findPrometheusContainer(t *testing.T, prometheus *promv1.Prometheus, name string) corev1.Container {
	t.Helper()
	for _, container := range prometheus.Spec.Containers {
		if container.Name == name {
			return container
		}
	}
	t.Fatalf("container %q not found", name)
	return corev1.Container{}
}
