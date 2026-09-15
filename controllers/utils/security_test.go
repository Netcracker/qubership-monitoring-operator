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

	assert.Equal(t, "monitoring-tmp", volume.Name)
	require.NotNil(t, volume.EmptyDir)
	require.NotNil(t, volume.EmptyDir.SizeLimit)
	assert.Equal(t, "16Mi", volume.EmptyDir.SizeLimit.String())
}

func TestTmpVolumeMount(t *testing.T) {
	volumeMount := TmpVolumeMount()

	assert.Equal(t, "monitoring-tmp", volumeMount.Name)
	assert.Equal(t, "/tmp", volumeMount.MountPath)
}

func TestHardenedPodSecurityContextWithOverridesOnOpenShift(t *testing.T) {
	configured := &monv1.SecurityContext{
		RunAsUser:  ptr.To(int64(5000)),
		RunAsGroup: ptr.To(int64(5001)),
		FSGroup:    ptr.To(int64(5002)),
	}

	securityContext, err := HardenedPodSecurityContextWithOverrides(true, configured)

	require.NoError(t, err)
	assert.Equal(t, configured.RunAsUser, securityContext.RunAsUser)
	assert.Equal(t, configured.RunAsGroup, securityContext.RunAsGroup)
	assert.Equal(t, configured.FSGroup, securityContext.FSGroup)
}

func TestHardenedPodSecurityContextWithOverridesNilConfigured(t *testing.T) {
	securityContext, err := HardenedPodSecurityContextWithOverrides(false, nil)

	require.NoError(t, err)
	require.NotNil(t, securityContext.RunAsUser)
	assert.Equal(t, int64(2000), *securityContext.RunAsUser)
}

func TestHardenedPodSecurityContextWithOverridesAppliesConfiguredIDs(t *testing.T) {
	configured := &monv1.SecurityContext{
		RunAsUser:  ptr.To(int64(2000)),
		RunAsGroup: ptr.To(int64(3000)),
		FSGroup:    ptr.To(int64(4000)),
	}

	securityContext, err := HardenedPodSecurityContextWithOverrides(false, configured)

	require.NoError(t, err)
	require.NotNil(t, securityContext.RunAsUser)
	assert.Equal(t, int64(2000), *securityContext.RunAsUser)
	require.NotNil(t, securityContext.RunAsGroup)
	assert.Equal(t, int64(3000), *securityContext.RunAsGroup)
	require.NotNil(t, securityContext.FSGroup)
	assert.Equal(t, int64(4000), *securityContext.FSGroup)
}

func TestHardenedPodSecurityContextWithPodOverrides(t *testing.T) {
	configured := &corev1.PodSecurityContext{
		RunAsUser:      ptr.To(int64(3000)),
		SELinuxOptions: &corev1.SELinuxOptions{Level: "s0"},
	}

	securityContext, err := HardenedPodSecurityContextWithPodOverrides(true, configured)

	require.NoError(t, err)
	assert.Equal(t, int64(3000), *securityContext.RunAsUser)
	assert.Nil(t, securityContext.RunAsGroup)
	assert.True(t, *securityContext.RunAsNonRoot)
	assert.Nil(t, securityContext.SELinuxOptions, "only numeric IDs are taken from the configured context")
}

func TestHardenedPodSecurityContextRejectsConflictingIdentity(t *testing.T) {
	tests := map[string]*corev1.PodSecurityContext{
		"runAsUser 0":        {RunAsUser: ptr.To(int64(0))},
		"runAsNonRoot false": {RunAsNonRoot: ptr.To(false)},
	}
	for name, configured := range tests {
		t.Run(name, func(t *testing.T) {
			securityContext, err := HardenedPodSecurityContextWithPodOverrides(false, configured)

			require.Error(t, err)
			assert.Nil(t, securityContext)
		})
	}

	securityContext, err := HardenedPodSecurityContextWithOverrides(false, &monv1.SecurityContext{RunAsUser: ptr.To(int64(0))})
	require.Error(t, err)
	assert.Nil(t, securityContext)
	assert.Error(t, ValidateSecurityContextSpec(&monv1.SecurityContext{RunAsUser: ptr.To(int64(0))}))
	assert.NoError(t, ValidateSecurityContextSpec(nil))
}

func TestHardenContainersWithTmpKeepsInputUnchanged(t *testing.T) {
	containers := []corev1.Container{{
		Name:            "managed",
		SecurityContext: &corev1.SecurityContext{RunAsUser: ptr.To(int64(3000))},
		VolumeMounts:    []corev1.VolumeMount{{Name: "data", MountPath: "/data"}},
	}}

	result, err := HardenContainersWithTmp(containers)

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, int64(3000), *result[0].SecurityContext.RunAsUser)
	assert.False(t, *result[0].SecurityContext.AllowPrivilegeEscalation)
	assert.Contains(t, result[0].VolumeMounts, TmpVolumeMount())
	assert.Nil(t, containers[0].SecurityContext.AllowPrivilegeEscalation)
	assert.Len(t, containers[0].VolumeMounts, 1)
}

func TestHardenContainersLeavesVolumeMountsUnchanged(t *testing.T) {
	containers := []corev1.Container{{
		Name:         "sidecar",
		VolumeMounts: []corev1.VolumeMount{{Name: "data", MountPath: "/data"}},
	}}

	result, err := HardenContainers(containers)

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, HardenedContainerSecurityContext(), result[0].SecurityContext)
	assert.Equal(t, containers[0].VolumeMounts, result[0].VolumeMounts)
}

func TestHardenContainersRejectsConflictingSettings(t *testing.T) {
	tests := map[string]*corev1.SecurityContext{
		"privileged":         {Privileged: ptr.To(true)},
		"added capabilities": {Capabilities: &corev1.Capabilities{Add: []corev1.Capability{"NET_RAW"}}},
		"runAsUser 0":        {RunAsUser: ptr.To(int64(0))},
		"runAsNonRoot false": {RunAsNonRoot: ptr.To(false)},
	}
	for name, securityContext := range tests {
		t.Run(name, func(t *testing.T) {
			containers := []corev1.Container{{Name: "sidecar", SecurityContext: securityContext}}

			result, err := HardenContainersWithTmp(containers)

			require.Error(t, err)
			assert.Contains(t, err.Error(), `container "sidecar"`)
			assert.Nil(t, result)

			result, err = HardenContainers(containers)

			require.Error(t, err)
			assert.Nil(t, result)
		})
	}
}

func TestMergeContainerSecurityContextNilExisting(t *testing.T) {
	merged, err := MergeContainerSecurityContext("app", nil)

	require.NoError(t, err)
	assert.Equal(t, HardenedContainerSecurityContext(), merged)
}

func TestMergeContainerSecurityContextKeepsUnrelatedFields(t *testing.T) {
	existing := &corev1.SecurityContext{
		RunAsUser:                ptr.To(int64(42)),
		AllowPrivilegeEscalation: ptr.To(true),
		ReadOnlyRootFilesystem:   ptr.To(false),
		Capabilities: &corev1.Capabilities{
			Drop: []corev1.Capability{"NET_RAW"},
		},
	}

	merged, err := MergeContainerSecurityContext("app", existing)

	require.NoError(t, err)
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
	assert.Equal(t, "monitoring-tmp", result[1].Name)
	require.NotNil(t, result[1].EmptyDir)
	assert.Equal(t, "16Mi", result[1].EmptyDir.SizeLimit.String())
	assert.Len(t, volumes, 1, "input slice must not be mutated in place")
}

func TestEnsureTmpVolumeKeepsExistingVolumeOfTheSameName(t *testing.T) {
	volumes := []corev1.Volume{{
		Name: "monitoring-tmp",
		VolumeSource: corev1.VolumeSource{
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{ClaimName: "scratch"},
		},
	}}

	result := EnsureTmpVolume(volumes, "100Mi")

	require.Len(t, result, 1)
	assert.Equal(t, volumes[0], result[0], "a user-defined volume must not be replaced")
}

func TestEnsureTmpVolumeKeepsUserVolumeNamedTmp(t *testing.T) {
	volumes := []corev1.Volume{{
		Name:         "tmp",
		VolumeSource: corev1.VolumeSource{Secret: &corev1.SecretVolumeSource{SecretName: "creds"}},
	}}

	result := EnsureTmpVolume(volumes, "100Mi")

	require.Len(t, result, 2)
	assert.Equal(t, volumes[0], result[0])
	assert.Equal(t, TmpVolume("100Mi"), result[1])
}

func TestEnsureTmpVolumeMountAppendsWhenMissing(t *testing.T) {
	mounts := []corev1.VolumeMount{
		{Name: "tmp", MountPath: "/other"},
		{Name: "data", MountPath: "/data"},
	}

	result := EnsureTmpVolumeMount(mounts)

	require.Len(t, result, 3)
	assert.Equal(t, mounts[0], result[0], "a user mount named tmp must be preserved")
	assert.Equal(t, mounts[1], result[1])
	assert.Equal(t, TmpVolumeMount(), result[2])
	assert.Len(t, mounts, 2, "input slice must not be mutated in place")
}

func TestEnsureTmpVolumeMountKeepsUserMountAtTmp(t *testing.T) {
	mounts := []corev1.VolumeMount{
		{Name: "scratch", MountPath: "/tmp"},
		{Name: "data", MountPath: "/data"},
	}

	result := EnsureTmpVolumeMount(mounts)

	assert.Equal(t, mounts, result, "an existing /tmp mount already provides the writable directory")
	assert.True(t, HasTmpVolumeMount(mounts))
	assert.False(t, HasTmpVolumeMount(mounts[1:]))
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
