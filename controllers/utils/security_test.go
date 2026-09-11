package utils

import (
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	secv1 "github.com/openshift/api/security/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/ptr"
)

func TestHardenedPodSecurityContextForKubernetes(t *testing.T) {
	securityContext := HardenedPodSecurityContext(false)

	require.NotNil(t, securityContext.RunAsNonRoot)
	assert.True(t, *securityContext.RunAsNonRoot)
	require.NotNil(t, securityContext.RunAsUser)
	assert.Equal(t, int64(1000), *securityContext.RunAsUser)
	require.NotNil(t, securityContext.RunAsGroup)
	assert.Equal(t, int64(1000), *securityContext.RunAsGroup)
	require.NotNil(t, securityContext.FSGroup)
	assert.Equal(t, int64(1000), *securityContext.FSGroup)
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

func TestHardenedPodSecurityContextWithOverridesOnOpenShift(t *testing.T) {
	configured := &monv1.SecurityContext{RunAsUser: ptr.To(int64(5000))}

	securityContext := HardenedPodSecurityContextWithOverrides(true, configured)

	assert.Nil(t, securityContext.RunAsUser, "OpenShift assigns UIDs, so overrides must not apply")
}

func TestHardenedPodSecurityContextWithOverridesNilConfigured(t *testing.T) {
	securityContext := HardenedPodSecurityContextWithOverrides(false, nil)

	require.NotNil(t, securityContext.RunAsUser)
	assert.Equal(t, int64(1000), *securityContext.RunAsUser)
}

func TestHardenedPodSecurityContextWithOverridesAppliesConfiguredIDs(t *testing.T) {
	configured := &monv1.SecurityContext{
		RunAsUser:  ptr.To(int64(2000)),
		RunAsGroup: ptr.To(int64(3000)),
		FSGroup:    ptr.To(int64(4000)),
	}

	securityContext := HardenedPodSecurityContextWithOverrides(false, configured)

	require.NotNil(t, securityContext.RunAsUser)
	assert.Equal(t, int64(2000), *securityContext.RunAsUser)
	require.NotNil(t, securityContext.RunAsGroup)
	assert.Equal(t, int64(3000), *securityContext.RunAsGroup)
	require.NotNil(t, securityContext.FSGroup)
	assert.Equal(t, int64(4000), *securityContext.FSGroup)
}

func TestMergeContainerSecurityContextNilExisting(t *testing.T) {
	merged := MergeContainerSecurityContext(nil)

	assert.Equal(t, HardenedContainerSecurityContext(), merged)
}

func TestMergeContainerSecurityContextKeepsUnrelatedFields(t *testing.T) {
	existing := &corev1.SecurityContext{
		RunAsUser:                ptr.To(int64(42)),
		AllowPrivilegeEscalation: ptr.To(true),
		ReadOnlyRootFilesystem:   ptr.To(false),
		Capabilities: &corev1.Capabilities{
			Add: []corev1.Capability{"NET_BIND_SERVICE"},
		},
	}

	merged := MergeContainerSecurityContext(existing)

	require.NotNil(t, merged.RunAsUser)
	assert.Equal(t, int64(42), *merged.RunAsUser, "unrelated fields must be preserved")
	require.NotNil(t, merged.AllowPrivilegeEscalation)
	assert.False(t, *merged.AllowPrivilegeEscalation, "hardening fields must be enforced")
	require.NotNil(t, merged.ReadOnlyRootFilesystem)
	assert.True(t, *merged.ReadOnlyRootFilesystem)
	assert.Equal(t, []corev1.Capability{"ALL"}, merged.Capabilities.Drop)
}

func TestEnsureTmpVolumeAppendsWhenMissing(t *testing.T) {
	volumes := []corev1.Volume{{Name: "data"}}

	result := EnsureTmpVolume(volumes, "16Mi")

	require.Len(t, result, 2)
	assert.Equal(t, "data", result[0].Name)
	assert.Equal(t, "tmp", result[1].Name)
	require.NotNil(t, result[1].EmptyDir)
	assert.Equal(t, "16Mi", result[1].EmptyDir.SizeLimit.String())
	assert.Len(t, volumes, 1, "input slice must not be mutated in place")
}

func TestEnsureTmpVolumeReplacesExisting(t *testing.T) {
	volumes := []corev1.Volume{TmpVolume("1Mi")}

	result := EnsureTmpVolume(volumes, "100Mi")

	require.Len(t, result, 1)
	assert.Equal(t, "100Mi", result[0].EmptyDir.SizeLimit.String())
}

func TestEnsureTmpVolumeMountDropsCollidingEntries(t *testing.T) {
	mounts := []corev1.VolumeMount{
		{Name: "tmp", MountPath: "/other"},
		{Name: "other", MountPath: "/tmp"},
		{Name: "data", MountPath: "/data"},
	}

	result := EnsureTmpVolumeMount(mounts)

	require.Len(t, result, 2)
	assert.Equal(t, "data", result[0].Name)
	assert.Equal(t, TmpVolumeMount(), result[1])
}

func TestApplySecurityContextConstraintsPolicy(t *testing.T) {
	existing := &secv1.SecurityContextConstraints{
		AllowPrivilegedContainer: true,
		AllowHostNetwork:         false,
	}
	desired := &secv1.SecurityContextConstraints{
		AllowPrivilegedContainer: false,
		AllowHostNetwork:         true,
		AllowHostPID:             true,
		Volumes:                  []secv1.FSType{secv1.FSTypeEmptyDir},
	}

	ApplySecurityContextConstraintsPolicy(existing, desired)

	assert.False(t, existing.AllowPrivilegedContainer)
	assert.True(t, existing.AllowHostNetwork)
	assert.True(t, existing.AllowHostPID)
	assert.Equal(t, []secv1.FSType{secv1.FSTypeEmptyDir}, existing.Volumes)
}
