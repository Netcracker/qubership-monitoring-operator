package victoriametrics

import (
	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	vmetricsv1b1 "github.com/VictoriaMetrics/operator/api/operator/v1beta1"
	corev1 "k8s.io/api/core/v1"
)

const tmpVolumeSize = "100Mi"

// HardenedSecurityContext returns the required VictoriaMetrics pod and container security settings.
func HardenedSecurityContext(isOpenShift bool, configured *corev1.PodSecurityContext) *vmetricsv1b1.SecurityContext {
	if configured == nil {
		return hardenedSecurityContext(isOpenShift, nil, nil, nil)
	}
	return hardenedSecurityContext(isOpenShift, configured.RunAsUser, configured.RunAsGroup, configured.FSGroup)
}

// HardenedSecurityContextFromPlatformSpec adapts the legacy platform security settings to VictoriaMetrics.
func HardenedSecurityContextFromPlatformSpec(
	isOpenShift bool,
	configured *monv1.SecurityContext,
) *vmetricsv1b1.SecurityContext {
	if configured == nil {
		return hardenedSecurityContext(isOpenShift, nil, nil, nil)
	}
	return hardenedSecurityContext(isOpenShift, configured.RunAsUser, configured.RunAsGroup, configured.FSGroup)
}

func hardenedSecurityContext(isOpenShift bool, runAsUser, runAsGroup, fsGroup *int64) *vmetricsv1b1.SecurityContext {
	podSecurityContext := utils.HardenedPodSecurityContext(isOpenShift)
	if !isOpenShift {
		if runAsUser != nil {
			podSecurityContext.RunAsUser = runAsUser
		}
		if runAsGroup != nil {
			podSecurityContext.RunAsGroup = runAsGroup
		}
		if fsGroup != nil {
			podSecurityContext.FSGroup = fsGroup
		}
	}

	containerSecurityContext := utils.HardenedContainerSecurityContext()
	return &vmetricsv1b1.SecurityContext{
		PodSecurityContext: podSecurityContext,
		ContainerSecurityContext: &vmetricsv1b1.ContainerSecurityContext{
			AllowPrivilegeEscalation: containerSecurityContext.AllowPrivilegeEscalation,
			ReadOnlyRootFilesystem:   containerSecurityContext.ReadOnlyRootFilesystem,
			Capabilities:             containerSecurityContext.Capabilities,
		},
	}
}

// EnsureTmpVolume returns a copy of the volumes containing the required size-limited temporary volume.
func EnsureTmpVolume(volumes []corev1.Volume) []corev1.Volume {
	result := make([]corev1.Volume, len(volumes))
	for i := range volumes {
		volumes[i].DeepCopyInto(&result[i])
	}

	required := utils.TmpVolume(tmpVolumeSize)
	for i := range result {
		if result[i].Name == required.Name {
			result[i] = required
			return result
		}
	}
	return append(result, required)
}

// EnsureTmpVolumeMount returns a copy of the mounts containing the required temporary-directory mount.
func EnsureTmpVolumeMount(volumeMounts []corev1.VolumeMount) []corev1.VolumeMount {
	result := make([]corev1.VolumeMount, 0, len(volumeMounts)+1)
	required := utils.TmpVolumeMount()
	for i := range volumeMounts {
		if volumeMounts[i].Name == required.Name || volumeMounts[i].MountPath == required.MountPath {
			continue
		}
		result = append(result, *volumeMounts[i].DeepCopy())
	}
	return append(result, required)
}

// HardenContainers returns hardened copies of explicitly configured application containers.
func HardenContainers(containers []corev1.Container) []corev1.Container {
	result := make([]corev1.Container, len(containers))
	for i := range containers {
		containers[i].DeepCopyInto(&result[i])
		securityContext := utils.HardenedContainerSecurityContext()
		if result[i].SecurityContext != nil {
			securityContext = result[i].SecurityContext.DeepCopy()
			required := utils.HardenedContainerSecurityContext()
			securityContext.AllowPrivilegeEscalation = required.AllowPrivilegeEscalation
			securityContext.ReadOnlyRootFilesystem = required.ReadOnlyRootFilesystem
			securityContext.Capabilities = required.Capabilities
		}
		result[i].SecurityContext = securityContext
		result[i].VolumeMounts = EnsureTmpVolumeMount(result[i].VolumeMounts)
	}
	return result
}
