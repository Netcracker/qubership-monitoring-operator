package victoriametrics

import (
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	vmetricsv1b1 "github.com/VictoriaMetrics/operator/api/operator/v1beta1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/utils/ptr"
)

func TestHardenedSecurityContextForKubernetes(t *testing.T) {
	configured := &monv1.SecurityContext{
		RunAsUser:  ptr.To[int64](3000),
		RunAsGroup: ptr.To[int64](3001),
		FSGroup:    ptr.To[int64](3002),
	}

	securityContext := HardenedSecurityContextFromPlatformSpec(false, configured)

	require.NotNil(t, securityContext.PodSecurityContext)
	assert.Equal(t, configured.RunAsUser, securityContext.RunAsUser)
	assert.Equal(t, configured.RunAsGroup, securityContext.RunAsGroup)
	assert.Equal(t, configured.FSGroup, securityContext.FSGroup)
	assert.Equal(t, ptr.To(true), securityContext.RunAsNonRoot)
	require.NotNil(t, securityContext.SeccompProfile)
	assert.Equal(t, corev1.SeccompProfileTypeRuntimeDefault, securityContext.SeccompProfile.Type)
	assertHardenedVictoriaMetricsContainerContext(t, securityContext.ContainerSecurityContext)
}

func TestHardenedSecurityContextForOpenShift(t *testing.T) {
	configured := &monv1.SecurityContext{
		RunAsUser:  ptr.To[int64](3000),
		RunAsGroup: ptr.To[int64](3001),
		FSGroup:    ptr.To[int64](3002),
	}

	securityContext := HardenedSecurityContextFromPlatformSpec(true, configured)

	require.NotNil(t, securityContext.PodSecurityContext)
	assert.Nil(t, securityContext.RunAsUser)
	assert.Nil(t, securityContext.RunAsGroup)
	assert.Nil(t, securityContext.FSGroup)
	assert.Equal(t, ptr.To(true), securityContext.RunAsNonRoot)
	require.NotNil(t, securityContext.SeccompProfile)
	assert.Equal(t, corev1.SeccompProfileTypeRuntimeDefault, securityContext.SeccompProfile.Type)
	assertHardenedVictoriaMetricsContainerContext(t, securityContext.ContainerSecurityContext)
}

func TestEnsureTmpVolumeReplacesAnExistingVolume(t *testing.T) {
	volumes := []corev1.Volume{{
		Name: "tmp",
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
	}}

	result := EnsureTmpVolume(volumes)

	require.Len(t, result, 1)
	require.NotNil(t, result[0].EmptyDir)
	require.NotNil(t, result[0].EmptyDir.SizeLimit)
	assert.Equal(t, resource.MustParse("100Mi"), *result[0].EmptyDir.SizeLimit)
	assert.Nil(t, volumes[0].EmptyDir.SizeLimit, "the source CR must not be mutated")
}

func TestHardenContainers(t *testing.T) {
	containers := []corev1.Container{{
		Name: "sidecar",
		SecurityContext: &corev1.SecurityContext{
			AllowPrivilegeEscalation: ptr.To(true),
			ReadOnlyRootFilesystem:   ptr.To(false),
		},
		VolumeMounts: []corev1.VolumeMount{{Name: "other", MountPath: "/tmp"}},
	}}

	result := HardenContainers(containers)

	require.Len(t, result, 1)
	assert.Equal(t, ptr.To(false), result[0].SecurityContext.AllowPrivilegeEscalation)
	assert.Equal(t, ptr.To(true), result[0].SecurityContext.ReadOnlyRootFilesystem)
	require.NotNil(t, result[0].SecurityContext.Capabilities)
	assert.Equal(t, []corev1.Capability{"ALL"}, result[0].SecurityContext.Capabilities.Drop)
	assert.Equal(t, []corev1.VolumeMount{{Name: "tmp", MountPath: "/tmp"}}, result[0].VolumeMounts)
	assert.Equal(t, ptr.To(true), containers[0].SecurityContext.AllowPrivilegeEscalation,
		"the source CR must not be mutated")
}

func assertHardenedVictoriaMetricsContainerContext(
	t *testing.T,
	securityContext *vmetricsv1b1.ContainerSecurityContext,
) {
	t.Helper()
	require.NotNil(t, securityContext)
	assert.Equal(t, ptr.To(false), securityContext.AllowPrivilegeEscalation)
	assert.Equal(t, ptr.To(true), securityContext.ReadOnlyRootFilesystem)
	require.NotNil(t, securityContext.Capabilities)
	assert.Equal(t, []corev1.Capability{"ALL"}, securityContext.Capabilities.Drop)
}
