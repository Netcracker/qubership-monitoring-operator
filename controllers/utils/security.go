package utils

import (
	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	secv1 "github.com/openshift/api/security/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/utils/ptr"
)

const (
	defaultSecurityContextID int64 = 1000
	tmpVolumeName                  = "tmp"
	tmpMountPath                   = "/tmp"
)

// HardenedPodSecurityContext returns the baseline pod security settings for the target platform.
func HardenedPodSecurityContext(isOpenShift bool) *corev1.PodSecurityContext {
	securityContext := &corev1.PodSecurityContext{
		RunAsNonRoot: ptr.To(true),
		SeccompProfile: &corev1.SeccompProfile{
			Type: corev1.SeccompProfileTypeRuntimeDefault,
		},
	}
	if !isOpenShift {
		securityContext.RunAsUser = ptr.To(defaultSecurityContextID)
		securityContext.RunAsGroup = ptr.To(defaultSecurityContextID)
		securityContext.FSGroup = ptr.To(defaultSecurityContextID)
	}
	return securityContext
}

// HardenedContainerSecurityContext returns the baseline security settings for a workload container.
func HardenedContainerSecurityContext() *corev1.SecurityContext {
	return &corev1.SecurityContext{
		AllowPrivilegeEscalation: ptr.To(false),
		ReadOnlyRootFilesystem:   ptr.To(true),
		Capabilities: &corev1.Capabilities{
			Drop: []corev1.Capability{"ALL"},
		},
	}
}

// TmpVolume returns a size-limited emptyDir volume for the pod's temporary directory.
func TmpVolume(sizeLimit string) corev1.Volume {
	quantity := resource.MustParse(sizeLimit)
	return corev1.Volume{
		Name: tmpVolumeName,
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{
				SizeLimit: &quantity,
			},
		},
	}
}

// TmpVolumeMount returns the standard temporary-directory mount for a workload container.
func TmpVolumeMount() corev1.VolumeMount {
	return corev1.VolumeMount{
		Name:      tmpVolumeName,
		MountPath: tmpMountPath,
	}
}

// HardenedPodSecurityContextWithOverrides returns the baseline pod security settings for the
// target platform, keeping any explicitly configured RunAsUser/RunAsGroup/FSGroup outside OpenShift.
func HardenedPodSecurityContextWithOverrides(isOpenShift bool, configured *monv1.SecurityContext) *corev1.PodSecurityContext {
	securityContext := HardenedPodSecurityContext(isOpenShift)
	if isOpenShift || configured == nil {
		return securityContext
	}
	if configured.RunAsUser != nil {
		securityContext.RunAsUser = configured.RunAsUser
	}
	if configured.RunAsGroup != nil {
		securityContext.RunAsGroup = configured.RunAsGroup
	}
	if configured.FSGroup != nil {
		securityContext.FSGroup = configured.FSGroup
	}
	return securityContext
}

// MergeContainerSecurityContext returns a security context that keeps any explicitly configured
// fields of existing while enforcing the required hardening fields.
func MergeContainerSecurityContext(existing *corev1.SecurityContext) *corev1.SecurityContext {
	required := HardenedContainerSecurityContext()
	if existing == nil {
		return required
	}
	merged := existing.DeepCopy()
	merged.AllowPrivilegeEscalation = required.AllowPrivilegeEscalation
	merged.ReadOnlyRootFilesystem = required.ReadOnlyRootFilesystem
	merged.Capabilities = required.Capabilities
	return merged
}

// EnsureTmpVolume returns a copy of volumes with the required size-limited temporary volume
// present, replacing any existing volume of the same name.
func EnsureTmpVolume(volumes []corev1.Volume, sizeLimit string) []corev1.Volume {
	result := make([]corev1.Volume, len(volumes))
	for i := range volumes {
		volumes[i].DeepCopyInto(&result[i])
	}

	required := TmpVolume(sizeLimit)
	for i := range result {
		if result[i].Name == required.Name {
			result[i] = required
			return result
		}
	}
	return append(result, required)
}

// EnsureTmpVolumeMount returns a copy of volumeMounts with the required temporary-directory mount
// present, dropping any existing mount that collides with it by name or path.
func EnsureTmpVolumeMount(volumeMounts []corev1.VolumeMount) []corev1.VolumeMount {
	result := make([]corev1.VolumeMount, 0, len(volumeMounts)+1)
	required := TmpVolumeMount()
	for i := range volumeMounts {
		if volumeMounts[i].Name == required.Name || volumeMounts[i].MountPath == required.MountPath {
			continue
		}
		result = append(result, *volumeMounts[i].DeepCopy())
	}
	return append(result, required)
}

// ApplySecurityContextConstraintsPolicy copies the policy fields of desired onto existing,
// leaving object metadata (name, labels, resourceVersion, ...) untouched.
func ApplySecurityContextConstraintsPolicy(existing, desired *secv1.SecurityContextConstraints) {
	existing.AllowPrivilegedContainer = desired.AllowPrivilegedContainer
	existing.DefaultAddCapabilities = desired.DefaultAddCapabilities
	existing.RequiredDropCapabilities = desired.RequiredDropCapabilities
	existing.AllowedCapabilities = desired.AllowedCapabilities
	existing.AllowHostDirVolumePlugin = desired.AllowHostDirVolumePlugin
	existing.Volumes = desired.Volumes
	existing.AllowedFlexVolumes = desired.AllowedFlexVolumes
	existing.AllowHostNetwork = desired.AllowHostNetwork
	existing.AllowHostPorts = desired.AllowHostPorts
	existing.AllowHostPID = desired.AllowHostPID
	existing.AllowHostIPC = desired.AllowHostIPC
	existing.DefaultAllowPrivilegeEscalation = desired.DefaultAllowPrivilegeEscalation
	existing.AllowPrivilegeEscalation = desired.AllowPrivilegeEscalation
	existing.SELinuxContext = desired.SELinuxContext
	existing.RunAsUser = desired.RunAsUser
	existing.SupplementalGroups = desired.SupplementalGroups
	existing.FSGroup = desired.FSGroup
	existing.ReadOnlyRootFilesystem = desired.ReadOnlyRootFilesystem
	existing.SeccompProfiles = desired.SeccompProfiles
	existing.AllowedUnsafeSysctls = desired.AllowedUnsafeSysctls
	existing.ForbiddenSysctls = desired.ForbiddenSysctls
}
