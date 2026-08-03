package pushgateway

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
	"k8s.io/utils/ptr"
)

var (
	cr         *monv1.PlatformMonitoring
	labelKey   = "label.key"
	labelValue = "label-value"
)

func TestPushgatewayManifests(t *testing.T) {
	cr = &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
	}
	t.Run("Test Deployment manifest", func(t *testing.T) {
		m, err := pushgatewayDeployment(cr, false)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Deployment manifest should not be empty")
		assertPushgatewayHardening(t, m, false)
	})
	t.Run("Test OpenShift Deployment security context", func(t *testing.T) {
		m, err := pushgatewayDeployment(cr, true)
		require.NoError(t, err)
		assertPushgatewayHardening(t, m, true)
	})
	t.Run("Test configured IDs are preserved", func(t *testing.T) {
		configuredCR := &monv1.PlatformMonitoring{
			ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"},
			Spec: monv1.PlatformMonitoringSpec{
				Pushgateway: &monv1.Pushgateway{
					SecurityContext: &monv1.SecurityContext{
						RunAsUser:  ptr.To(int64(3000)),
						RunAsGroup: ptr.To(int64(3001)),
						FSGroup:    ptr.To(int64(3002)),
					},
				},
			},
		}

		m, err := pushgatewayDeployment(configuredCR, false)
		require.NoError(t, err)
		assert.Equal(t, int64(3000), *m.Spec.Template.Spec.SecurityContext.RunAsUser)
		assert.Equal(t, int64(3001), *m.Spec.Template.Spec.SecurityContext.RunAsGroup)
		assert.Equal(t, int64(3002), *m.Spec.Template.Spec.SecurityContext.FSGroup)
		assertPushgatewayHardening(t, m, false)

		openShiftDeployment, err := pushgatewayDeployment(configuredCR, true)
		require.NoError(t, err)
		assert.Nil(t, openShiftDeployment.Spec.Template.Spec.SecurityContext.RunAsUser)
		assert.Nil(t, openShiftDeployment.Spec.Template.Spec.SecurityContext.RunAsGroup)
		assert.Nil(t, openShiftDeployment.Spec.Template.Spec.SecurityContext.FSGroup)
	})
	t.Run("Test persistence and custom volumes are preserved", func(t *testing.T) {
		persistentCR := &monv1.PlatformMonitoring{
			ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"},
			Spec: monv1.PlatformMonitoringSpec{
				Pushgateway: &monv1.Pushgateway{
					Storage: &corev1.PersistentVolumeClaimSpec{
						AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
					},
					Volumes: []corev1.Volume{{
						Name: "custom",
						VolumeSource: corev1.VolumeSource{
							ConfigMap: &corev1.ConfigMapVolumeSource{},
						},
					}},
					VolumeMounts: []corev1.VolumeMount{{Name: "custom", MountPath: "/etc/custom"}},
				},
			},
		}

		m, err := pushgatewayDeployment(persistentCR, false)
		require.NoError(t, err)
		assertPushgatewayHardening(t, m, false)

		container := m.Spec.Template.Spec.Containers[0]
		assert.Contains(t, container.VolumeMounts,
			corev1.VolumeMount{Name: utils.PushgatewayStorageVolumeName, MountPath: "/data"})
		assert.Contains(t, container.VolumeMounts,
			corev1.VolumeMount{Name: "custom", MountPath: "/etc/custom"})
		assert.Contains(t, container.Args, "--persistence.file=/data/pushgateway.data")
		assert.Contains(t, container.Args, "--persistence.interval=5m")

		volumeNames := make(map[string]bool, len(m.Spec.Template.Spec.Volumes))
		for _, volume := range m.Spec.Template.Spec.Volumes {
			volumeNames[volume.Name] = true
		}
		assert.True(t, volumeNames[utils.PushgatewayStorageVolumeName])
		assert.True(t, volumeNames["custom"])
		assert.True(t, volumeNames["tmp"])
	})
	t.Run("Test existing temporary volume and mount are replaced", func(t *testing.T) {
		deployment := &appsv1.Deployment{
			Spec: appsv1.DeploymentSpec{
				Template: corev1.PodTemplateSpec{
					Spec: corev1.PodSpec{
						Volumes: []corev1.Volume{{
							Name: "tmp",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							},
						}},
						Containers: []corev1.Container{{
							Name:         "pushgateway",
							VolumeMounts: []corev1.VolumeMount{{Name: "other-tmp", MountPath: "/tmp"}},
						}},
					},
				},
			},
		}

		applyPushgatewayHardening(deployment, false, nil)

		assertPushgatewayHardening(t, deployment, false)
		assert.NotContains(t, deployment.Spec.Template.Spec.Containers[0].VolumeMounts,
			corev1.VolumeMount{Name: "other-tmp", MountPath: "/tmp"})
	})
	t.Run("Test Service manifest", func(t *testing.T) {
		m, err := pushgatewayService(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Service manifest should not be empty")
	})
	t.Run("Test Ingress v1 manifest", func(t *testing.T) {
		m, err := pushgatewayIngressV1(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Ingress v1 manifest should not be empty")
	})
	t.Run("Test ServiceMonitor manifest", func(t *testing.T) {
		crWithLabels := &monv1.PlatformMonitoring{
			ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring", Labels: map[string]string{labelKey: labelValue}},
			Spec:       monv1.PlatformMonitoringSpec{Pushgateway: &monv1.Pushgateway{}},
		}
		m, err := pushgatewayServiceMonitor(crWithLabels)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "ServiceMonitor manifest should not be empty")
		labelsassert.AssertCRLabels(t, m.GetLabels(), utils.PushgatewayComponentName, "prometheus-operator", map[string]string{labelKey: labelValue})
	})
}

func assertPushgatewayHardening(t *testing.T, deployment *appsv1.Deployment, isOpenShift bool) {
	t.Helper()

	podSecurityContext := deployment.Spec.Template.Spec.SecurityContext
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

	require.Len(t, deployment.Spec.Template.Spec.Containers, 1)
	container := deployment.Spec.Template.Spec.Containers[0]
	require.NotNil(t, container.SecurityContext)
	assert.Equal(t, ptr.To(false), container.SecurityContext.AllowPrivilegeEscalation)
	assert.Equal(t, ptr.To(true), container.SecurityContext.ReadOnlyRootFilesystem)
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
