package utils

import (
	"fmt"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	secv1 "github.com/openshift/api/security/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/utils/ptr"
)

const (
	defaultSecurityContextID int64 = 2000
	// tmpVolumeName is deliberately distinctive so that it cannot collide with a user-defined volume.
	tmpVolumeName = "monitoring-tmp"
	tmpMountPath  = "/tmp"
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
// target platform, keeping any explicitly configured RunAsUser/RunAsGroup/FSGroup.
// It returns an error when the configured identity conflicts with the non-root baseline.
func HardenedPodSecurityContextWithOverrides(isOpenShift bool, configured *monv1.SecurityContext) (*corev1.PodSecurityContext, error) {
	if configured == nil {
		return HardenedPodSecurityContext(isOpenShift), nil
	}
	return HardenedPodSecurityContextWithPodOverrides(isOpenShift, &corev1.PodSecurityContext{
		RunAsUser: configured.RunAsUser, RunAsGroup: configured.RunAsGroup, FSGroup: configured.FSGroup,
	})
}

// HardenedPodSecurityContextWithPodOverrides keeps explicitly configured numeric IDs while
// enforcing the pod security baseline.
// It returns an error when the configured identity conflicts with the non-root baseline.
func HardenedPodSecurityContextWithPodOverrides(isOpenShift bool, configured *corev1.PodSecurityContext) (*corev1.PodSecurityContext, error) {
	securityContext := HardenedPodSecurityContext(isOpenShift)
	if configured == nil {
		return securityContext, nil
	}
	if err := ValidatePodSecurityContext(configured); err != nil {
		return nil, err
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
	return securityContext, nil
}

// ValidatePodSecurityContext rejects pod settings that the enforced runAsNonRoot baseline cannot
// override: UID 0 and an explicit runAsNonRoot=false.
func ValidatePodSecurityContext(configured *corev1.PodSecurityContext) error {
	if configured == nil {
		return nil
	}
	if configured.RunAsUser != nil && *configured.RunAsUser == 0 {
		return fmt.Errorf("pod securityContext.runAsUser=0 conflicts with the enforced runAsNonRoot=true baseline: use a non-root UID")
	}
	if configured.RunAsNonRoot != nil && !*configured.RunAsNonRoot {
		return fmt.Errorf("pod securityContext.runAsNonRoot=false conflicts with the enforced security baseline: remove the field")
	}
	return nil
}

// ValidateSecurityContextSpec rejects a PlatformMonitoring security context whose UID conflicts
// with the enforced runAsNonRoot baseline.
func ValidateSecurityContextSpec(configured *monv1.SecurityContext) error {
	if configured == nil {
		return nil
	}
	return ValidatePodSecurityContext(&corev1.PodSecurityContext{RunAsUser: configured.RunAsUser})
}

// ValidateContainerSecurityContext rejects container settings that the enforced baseline cannot
// override: privileged mode, added capabilities, UID 0, and an explicit runAsNonRoot=false.
func ValidateContainerSecurityContext(name string, configured *corev1.SecurityContext) error {
	if configured == nil {
		return nil
	}
	if configured.Privileged != nil && *configured.Privileged {
		return fmt.Errorf("container %q: securityContext.privileged=true conflicts with the enforced security baseline: remove the field", name)
	}
	if configured.Capabilities != nil && len(configured.Capabilities.Add) > 0 {
		return fmt.Errorf("container %q: securityContext.capabilities.add conflicts with the enforced drop-all baseline: remove the added capabilities", name)
	}
	if configured.RunAsUser != nil && *configured.RunAsUser == 0 {
		return fmt.Errorf("container %q: securityContext.runAsUser=0 conflicts with the enforced runAsNonRoot=true baseline: use a non-root UID", name)
	}
	if configured.RunAsNonRoot != nil && !*configured.RunAsNonRoot {
		return fmt.Errorf("container %q: securityContext.runAsNonRoot=false conflicts with the enforced security baseline: remove the field", name)
	}
	return nil
}

// HardenContainer enforces the container baseline without touching volume mounts.
func HardenContainer(container *corev1.Container) error {
	merged, err := MergeContainerSecurityContext(container.Name, container.SecurityContext)
	if err != nil {
		return err
	}
	container.SecurityContext = merged
	return nil
}

// HardenContainerWithTmp enforces the container baseline and mounts the temporary volume.
func HardenContainerWithTmp(container *corev1.Container) error {
	if err := HardenContainer(container); err != nil {
		return err
	}
	container.VolumeMounts = EnsureTmpVolumeMount(container.VolumeMounts)
	return nil
}

// HardenContainers returns copies of all containers with the security baseline enforced and
// volume mounts left unchanged.
func HardenContainers(containers []corev1.Container) ([]corev1.Container, error) {
	return hardenContainers(containers, HardenContainer)
}

// HardenContainersWithTmp returns hardened copies of all containers.
func HardenContainersWithTmp(containers []corev1.Container) ([]corev1.Container, error) {
	return hardenContainers(containers, HardenContainerWithTmp)
}

func hardenContainers(containers []corev1.Container, harden func(*corev1.Container) error) ([]corev1.Container, error) {
	result := make([]corev1.Container, len(containers))
	for i := range containers {
		containers[i].DeepCopyInto(&result[i])
		if err := harden(&result[i]); err != nil {
			return nil, err
		}
	}
	return result, nil
}

// MergeContainerSecurityContext returns a security context that keeps any explicitly configured
// fields of existing while enforcing the required hardening fields. It returns an error when
// existing contains a field that the baseline cannot override.
func MergeContainerSecurityContext(name string, existing *corev1.SecurityContext) (*corev1.SecurityContext, error) {
	required := HardenedContainerSecurityContext()
	if existing == nil {
		return required, nil
	}
	if err := ValidateContainerSecurityContext(name, existing); err != nil {
		return nil, err
	}
	merged := existing.DeepCopy()
	merged.AllowPrivilegeEscalation = required.AllowPrivilegeEscalation
	merged.ReadOnlyRootFilesystem = required.ReadOnlyRootFilesystem
	merged.Capabilities = required.Capabilities
	return merged, nil
}

// EnsureTmpVolume returns a copy of volumes with the size-limited temporary volume present.
// The volume name is reserved: a user-defined volume with the same name is rejected, because
// the temporary-directory mount would otherwise expose that volume at /tmp.
func EnsureTmpVolume(volumes []corev1.Volume, sizeLimit string) ([]corev1.Volume, error) {
	result := make([]corev1.Volume, len(volumes))
	for i := range volumes {
		volumes[i].DeepCopyInto(&result[i])
	}

	required := TmpVolume(sizeLimit)
	for i := range result {
		if result[i].Name == required.Name {
			return nil, fmt.Errorf("volume name %q is reserved for the operator-managed temporary directory: rename the volume", required.Name)
		}
	}
	return append(result, required), nil
}

// EnsureTmpVolumeMount returns a copy of volumeMounts with a writable temporary directory.
// A user-defined mount at /tmp already satisfies that requirement and is kept as is.
func EnsureTmpVolumeMount(volumeMounts []corev1.VolumeMount) []corev1.VolumeMount {
	result := make([]corev1.VolumeMount, 0, len(volumeMounts)+1)
	for i := range volumeMounts {
		result = append(result, *volumeMounts[i].DeepCopy())
	}
	if HasTmpVolumeMount(volumeMounts) {
		return result
	}
	return append(result, TmpVolumeMount())
}

// HasTmpVolumeMount reports whether any mount already covers the temporary directory.
func HasTmpVolumeMount(volumeMounts []corev1.VolumeMount) bool {
	for i := range volumeMounts {
		if volumeMounts[i].MountPath == tmpMountPath {
			return true
		}
	}
	return false
}

// ApplySecurityContextConstraintsPolicy copies the policy fields of desired onto existing,
// leaving object metadata (name, labels, resourceVersion, ...) untouched.
// Returns true when existing was updated.
func ApplySecurityContextConstraintsPolicy(existing, desired *secv1.SecurityContextConstraints) bool {
	return AnyChanged(
		SetIfChanged(&existing.AllowPrivilegedContainer, desired.AllowPrivilegedContainer),
		SetIfChanged(&existing.DefaultAddCapabilities, desired.DefaultAddCapabilities),
		SetIfChanged(&existing.RequiredDropCapabilities, desired.RequiredDropCapabilities),
		SetIfChanged(&existing.AllowedCapabilities, desired.AllowedCapabilities),
		SetIfChanged(&existing.AllowHostDirVolumePlugin, desired.AllowHostDirVolumePlugin),
		SetIfChanged(&existing.Volumes, desired.Volumes),
		SetIfChanged(&existing.AllowedFlexVolumes, desired.AllowedFlexVolumes),
		SetIfChanged(&existing.AllowHostNetwork, desired.AllowHostNetwork),
		SetIfChanged(&existing.AllowHostPorts, desired.AllowHostPorts),
		SetIfChanged(&existing.AllowHostPID, desired.AllowHostPID),
		SetIfChanged(&existing.AllowHostIPC, desired.AllowHostIPC),
		SetIfChanged(&existing.DefaultAllowPrivilegeEscalation, desired.DefaultAllowPrivilegeEscalation),
		SetIfChanged(&existing.AllowPrivilegeEscalation, desired.AllowPrivilegeEscalation),
		SetIfChanged(&existing.SELinuxContext, desired.SELinuxContext),
		SetIfChanged(&existing.RunAsUser, desired.RunAsUser),
		SetIfChanged(&existing.SupplementalGroups, desired.SupplementalGroups),
		SetIfChanged(&existing.FSGroup, desired.FSGroup),
		SetIfChanged(&existing.ReadOnlyRootFilesystem, desired.ReadOnlyRootFilesystem),
		SetIfChanged(&existing.SeccompProfiles, desired.SeccompProfiles),
		SetIfChanged(&existing.AllowedUnsafeSysctls, desired.AllowedUnsafeSysctls),
		SetIfChanged(&existing.ForbiddenSysctls, desired.ForbiddenSysctls),
	)
}
