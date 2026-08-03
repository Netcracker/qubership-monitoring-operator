package alertmanager

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

func TestAlertmanagerManifests(t *testing.T) {
	cr = &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
		Spec: monv1.PlatformMonitoringSpec{
			AlertManager: &monv1.AlertManager{
				Annotations: map[string]string{annotationKey: annotationValue},
				Labels:      map[string]string{labelKey: labelValue},
			},
		},
	}
	t.Run("Test Alert Manager manifest", func(t *testing.T) {
		m, err := alertmanager(cr, false)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Alert Manager manifest should not be empty")
		assert.NotNil(t, m.Spec.PodMetadata.Labels)
		assert.Equal(t, labelValue, m.Spec.PodMetadata.Labels[labelKey])
		assert.NotNil(t, m.Spec.PodMetadata.Annotations)
		assert.Equal(t, annotationValue, m.Spec.PodMetadata.Annotations[annotationKey])
		assertAlertmanagerHardening(t, m, false)
	})
	cr = &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
		Spec: monv1.PlatformMonitoringSpec{
			AlertManager: &monv1.AlertManager{},
		},
	}
	t.Run("Test Alert Manager manifest with nil labels and annotation", func(t *testing.T) {
		m, err := alertmanager(cr, false)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Alert Manager manifest should not be empty")
		assert.NotNil(t, m.GetLabels())
		assert.Nil(t, m.GetAnnotations())
	})
	t.Run("Test OpenShift security context", func(t *testing.T) {
		m, err := alertmanager(cr, true)
		require.NoError(t, err)
		assertAlertmanagerHardening(t, m, true)
	})
	t.Run("Test configured IDs and sidecar settings", func(t *testing.T) {
		configuredCR := cr.DeepCopy()
		configuredCR.Spec.AlertManager.SecurityContext = &monv1.SecurityContext{
			RunAsUser:  ptr.To(int64(3000)),
			RunAsGroup: ptr.To(int64(3001)),
			FSGroup:    ptr.To(int64(3002)),
		}
		configuredCR.Spec.AlertManager.Containers = []corev1.Container{{
			Name: "sidecar",
			SecurityContext: &corev1.SecurityContext{
				AllowPrivilegeEscalation: ptr.To(true),
				ReadOnlyRootFilesystem:   ptr.To(false),
			},
			VolumeMounts: []corev1.VolumeMount{{Name: "other-tmp", MountPath: "/tmp"}},
		}}

		m, err := alertmanager(configuredCR, false)
		require.NoError(t, err)
		assert.Equal(t, int64(3000), *m.Spec.SecurityContext.RunAsUser)
		assert.Equal(t, int64(3001), *m.Spec.SecurityContext.RunAsGroup)
		assert.Equal(t, int64(3002), *m.Spec.SecurityContext.FSGroup)
		assertAlertmanagerHardening(t, m, false)
		sidecar := findAlertmanagerContainer(t, m, "sidecar")
		assertAlertmanagerContainerHardening(t, sidecar)
		assert.NotContains(t, sidecar.VolumeMounts, corev1.VolumeMount{Name: "other-tmp", MountPath: "/tmp"})
		assert.Equal(t, ptr.To(true), configuredCR.Spec.AlertManager.Containers[0].SecurityContext.AllowPrivilegeEscalation,
			"the source CR must not be mutated")

		openShiftManifest, err := alertmanager(configuredCR, true)
		require.NoError(t, err)
		assert.Nil(t, openShiftManifest.Spec.SecurityContext.RunAsUser)
		assert.Nil(t, openShiftManifest.Spec.SecurityContext.RunAsGroup)
		assert.Nil(t, openShiftManifest.Spec.SecurityContext.FSGroup)
	})
	t.Run("Test existing temporary volume is replaced", func(t *testing.T) {
		volumes := []corev1.Volume{{
			Name: "tmp",
			VolumeSource: corev1.VolumeSource{
				EmptyDir: &corev1.EmptyDirVolumeSource{},
			},
		}}

		result := ensureAlertmanagerTmpVolume(volumes)

		require.Len(t, result, 1)
		require.NotNil(t, result[0].EmptyDir)
		require.NotNil(t, result[0].EmptyDir.SizeLimit)
		assert.Equal(t, resource.MustParse("100Mi"), *result[0].EmptyDir.SizeLimit)
		assert.Nil(t, volumes[0].EmptyDir.SizeLimit, "the source volume must not be mutated")
	})
	t.Run("Test ServiceAccount manifest", func(t *testing.T) {
		crWithSALabels := &monv1.PlatformMonitoring{
			ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"},
			Spec: monv1.PlatformMonitoringSpec{
				AlertManager: &monv1.AlertManager{
					ServiceAccount: &monv1.EmbeddedObjectMetadata{
						Labels: map[string]string{labelKey: labelValue},
					},
				},
			},
		}
		m, err := alertmanagerServiceAccount(crWithSALabels)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "ServiceAccount manifest should not be empty")
		assert.NotNil(t, m.GetLabels())
		assert.Equal(t, labelValue, m.GetLabels()[labelKey], "ServiceAccount.Labels should be merged")
	})
	t.Run("Test Secret manifest", func(t *testing.T) {
		m, err := alertmanagerSecret(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Secret manifest should not be empty")
	})

	t.Run("Test Service manifest", func(t *testing.T) {
		m, err := alertmanagerService(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Service manifest should not be empty")
	})
	t.Run("Test Ingress v1 manifest", func(t *testing.T) {
		m, err := alertmanagerIngressV1(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Ingress v1 manifest should not be empty")
	})
	t.Run("Test PodMonitor manifest", func(t *testing.T) {
		crWithLabels := &monv1.PlatformMonitoring{
			ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring", Labels: map[string]string{labelKey: labelValue}},
			Spec:       monv1.PlatformMonitoringSpec{AlertManager: &monv1.AlertManager{}},
		}
		m, err := alertmanagerPodMonitor(crWithLabels)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "PodMonitor manifest should not be empty")
		labelsassert.AssertCRLabels(t, m.GetLabels(), utils.AlertManagerComponentName, "prometheus-operator", map[string]string{labelKey: labelValue})
	})
}

func assertAlertmanagerHardening(t *testing.T, alertmanager *promv1.Alertmanager, isOpenShift bool) {
	t.Helper()

	podSecurityContext := alertmanager.Spec.SecurityContext
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
	for _, volume := range alertmanager.Spec.Volumes {
		if volume.Name == "tmp" {
			tmpVolumes++
			require.NotNil(t, volume.EmptyDir)
			require.NotNil(t, volume.EmptyDir.SizeLimit)
			assert.Equal(t, resource.MustParse("100Mi"), *volume.EmptyDir.SizeLimit)
		}
	}
	assert.Equal(t, 1, tmpVolumes)

	assertAlertmanagerContainerHardening(t, findAlertmanagerContainer(t, alertmanager, "alertmanager"))
	assertAlertmanagerContainerHardening(t, findAlertmanagerContainer(t, alertmanager, "config-reloader"))
}

func assertAlertmanagerContainerHardening(t *testing.T, container corev1.Container) {
	t.Helper()
	require.NotNil(t, container.SecurityContext)
	assert.Equal(t, ptr.To(false), container.SecurityContext.AllowPrivilegeEscalation)
	assert.Equal(t, ptr.To(true), container.SecurityContext.ReadOnlyRootFilesystem)
	require.NotNil(t, container.SecurityContext.Capabilities)
	assert.Equal(t, []corev1.Capability{"ALL"}, container.SecurityContext.Capabilities.Drop)
	assert.Contains(t, container.VolumeMounts, utils.TmpVolumeMount())
}

func findAlertmanagerContainer(t *testing.T, alertmanager *promv1.Alertmanager, name string) corev1.Container {
	t.Helper()
	for _, container := range alertmanager.Spec.Containers {
		if container.Name == name {
			return container
		}
	}
	t.Fatalf("container %q not found", name)
	return corev1.Container{}
}
