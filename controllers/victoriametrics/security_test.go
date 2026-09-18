package victoriametrics

import (
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
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

	securityContext, err := HardenedSecurityContextFromPlatformSpec(false, configured)

	require.NoError(t, err)
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

	securityContext, err := HardenedSecurityContextFromPlatformSpec(true, configured)

	require.NoError(t, err)
	require.NotNil(t, securityContext.PodSecurityContext)
	assert.Equal(t, configured.RunAsUser, securityContext.RunAsUser)
	assert.Equal(t, configured.RunAsGroup, securityContext.RunAsGroup)
	assert.Equal(t, configured.FSGroup, securityContext.FSGroup)
	assert.Equal(t, ptr.To(true), securityContext.RunAsNonRoot)
	require.NotNil(t, securityContext.SeccompProfile)
	assert.Equal(t, corev1.SeccompProfileTypeRuntimeDefault, securityContext.SeccompProfile.Type)
	assertHardenedVictoriaMetricsContainerContext(t, securityContext.ContainerSecurityContext)
}

func TestHardenedSecurityContextRejectsRootUser(t *testing.T) {
	securityContext, err := HardenedSecurityContextFromPlatformSpec(false, &monv1.SecurityContext{
		RunAsUser: ptr.To[int64](0),
	})
	require.Error(t, err)
	assert.Nil(t, securityContext)

	securityContext, err = HardenedSecurityContext(false, &corev1.PodSecurityContext{
		RunAsNonRoot: ptr.To(false),
	})
	require.Error(t, err)
	assert.Nil(t, securityContext)
}

func TestEnsureTmpVolumeKeepsAnExistingUserVolume(t *testing.T) {
	volumes := []corev1.Volume{{
		Name: "tmp",
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{},
		},
	}}

	result, err := EnsureTmpVolume(volumes)

	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, volumes[0], result[0], "the user volume must be kept as is")
	require.NotNil(t, result[1].EmptyDir)
	require.NotNil(t, result[1].EmptyDir.SizeLimit)
	assert.Equal(t, resource.MustParse("100Mi"), *result[1].EmptyDir.SizeLimit)
	assert.Nil(t, volumes[0].EmptyDir.SizeLimit, "the source CR must not be mutated")
}

func TestEnsureTmpVolumeRejectsReservedName(t *testing.T) {
	volumes := []corev1.Volume{{
		Name:         utils.TmpVolumeMount().Name,
		VolumeSource: corev1.VolumeSource{Secret: &corev1.SecretVolumeSource{SecretName: "creds"}},
	}}

	result, err := EnsureTmpVolume(volumes)

	require.Error(t, err)
	assert.Nil(t, result)
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

	result, err := HardenContainers(containers)

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, ptr.To(false), result[0].SecurityContext.AllowPrivilegeEscalation)
	assert.Equal(t, ptr.To(true), result[0].SecurityContext.ReadOnlyRootFilesystem)
	require.NotNil(t, result[0].SecurityContext.Capabilities)
	assert.Equal(t, []corev1.Capability{"ALL"}, result[0].SecurityContext.Capabilities.Drop)
	assert.Equal(t, []corev1.VolumeMount{{Name: "other", MountPath: "/tmp"}}, result[0].VolumeMounts,
		"an existing /tmp mount must be kept")
	assert.Equal(t, ptr.To(true), containers[0].SecurityContext.AllowPrivilegeEscalation,
		"the source CR must not be mutated")
}

func TestHardenContainersAddsTmpMountWhenMissing(t *testing.T) {
	containers := []corev1.Container{{Name: "sidecar"}}

	result, err := HardenContainers(containers)

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, []corev1.VolumeMount{utils.TmpVolumeMount()}, result[0].VolumeMounts)
}

func TestHardenContainersRejectsPrivilegedContainer(t *testing.T) {
	containers := []corev1.Container{{
		Name:            "sidecar",
		SecurityContext: &corev1.SecurityContext{Privileged: ptr.To(true)},
	}}

	result, err := HardenContainers(containers)

	require.Error(t, err)
	assert.Nil(t, result)
}

func TestHardenContainerSecurityLeavesVolumeMountsUnchanged(t *testing.T) {
	containers := []corev1.Container{{
		Name:         "sidecar",
		VolumeMounts: []corev1.VolumeMount{{Name: "data", MountPath: "/data"}},
	}}

	result, err := HardenContainerSecurity(containers)

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, ptr.To(true), result[0].SecurityContext.ReadOnlyRootFilesystem)
	assert.Equal(t, containers[0].VolumeMounts, result[0].VolumeMounts)
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
