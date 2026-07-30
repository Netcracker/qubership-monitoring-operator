package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
)

func TestHardenedPodSecurityContextForKubernetes(t *testing.T) {
	securityContext := HardenedPodSecurityContext(false)

	require.NotNil(t, securityContext.RunAsNonRoot)
	assert.True(t, *securityContext.RunAsNonRoot)
	require.NotNil(t, securityContext.RunAsUser)
	assert.Equal(t, int64(2000), *securityContext.RunAsUser)
	require.NotNil(t, securityContext.RunAsGroup)
	assert.Equal(t, int64(2000), *securityContext.RunAsGroup)
	require.NotNil(t, securityContext.FSGroup)
	assert.Equal(t, int64(2000), *securityContext.FSGroup)
	require.NotNil(t, securityContext.SeccompProfile)
	assert.Equal(t, corev1.SeccompProfileTypeRuntimeDefault, securityContext.SeccompProfile.Type)
}

func TestHardenedPodSecurityContextForOpenShift(t *testing.T) {
	securityContext := HardenedPodSecurityContext(true)

	require.NotNil(t, securityContext.RunAsNonRoot)
	assert.True(t, *securityContext.RunAsNonRoot)
	assert.Nil(t, securityContext.RunAsUser)
	assert.Nil(t, securityContext.RunAsGroup)
	assert.Nil(t, securityContext.FSGroup)
	require.NotNil(t, securityContext.SeccompProfile)
	assert.Equal(t, corev1.SeccompProfileTypeRuntimeDefault, securityContext.SeccompProfile.Type)
}

func TestHardenedContainerSecurityContext(t *testing.T) {
	securityContext := HardenedContainerSecurityContext()

	require.NotNil(t, securityContext.AllowPrivilegeEscalation)
	assert.False(t, *securityContext.AllowPrivilegeEscalation)
	require.NotNil(t, securityContext.ReadOnlyRootFilesystem)
	assert.True(t, *securityContext.ReadOnlyRootFilesystem)
	require.NotNil(t, securityContext.Capabilities)
	assert.Equal(t, []corev1.Capability{"ALL"}, securityContext.Capabilities.Drop)
}

func TestTmpVolume(t *testing.T) {
	volume := TmpVolume("16Mi")

	assert.Equal(t, "tmp", volume.Name)
	require.NotNil(t, volume.EmptyDir)
	require.NotNil(t, volume.EmptyDir.SizeLimit)
	assert.Equal(t, "16Mi", volume.EmptyDir.SizeLimit.String())
}

func TestTmpVolumeMount(t *testing.T) {
	volumeMount := TmpVolumeMount()

	assert.Equal(t, "tmp", volumeMount.Name)
	assert.Equal(t, "/tmp", volumeMount.MountPath)
}
