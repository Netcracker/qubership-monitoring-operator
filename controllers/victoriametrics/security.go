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
	return utils.EnsureTmpVolume(volumes, tmpVolumeSize)
}

// EnsureTmpVolumeMount returns a copy of the mounts containing the required temporary-directory mount.
func EnsureTmpVolumeMount(volumeMounts []corev1.VolumeMount) []corev1.VolumeMount {
	return utils.EnsureTmpVolumeMount(volumeMounts)
}

// HardenContainers returns hardened copies of explicitly configured application containers.
func HardenContainers(containers []corev1.Container) []corev1.Container {
	result := make([]corev1.Container, len(containers))
	for i := range containers {
		containers[i].DeepCopyInto(&result[i])
		result[i].SecurityContext = utils.MergeContainerSecurityContext(result[i].SecurityContext)
		result[i].VolumeMounts = EnsureTmpVolumeMount(result[i].VolumeMounts)
	}
	return result
}
