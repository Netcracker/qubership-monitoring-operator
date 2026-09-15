package victoriametrics

import (
	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	vmetricsv1b1 "github.com/VictoriaMetrics/operator/api/operator/v1beta1"
	corev1 "k8s.io/api/core/v1"
)

const tmpVolumeSize = "100Mi"

// HardenedSecurityContext returns the required VictoriaMetrics pod and container security settings.
// It returns an error when the configured identity conflicts with the non-root baseline.
func HardenedSecurityContext(isOpenShift bool, configured *corev1.PodSecurityContext) (*vmetricsv1b1.SecurityContext, error) {
	return hardenedSecurityContext(isOpenShift, configured)
}

// HardenedSecurityContextFromPlatformSpec adapts the legacy platform security settings to VictoriaMetrics.
// It returns an error when the configured identity conflicts with the non-root baseline.
func HardenedSecurityContextFromPlatformSpec(
	isOpenShift bool,
	configured *monv1.SecurityContext,
) (*vmetricsv1b1.SecurityContext, error) {
	if configured == nil {
		return hardenedSecurityContext(isOpenShift, nil)
	}
	return hardenedSecurityContext(isOpenShift, &corev1.PodSecurityContext{
		RunAsUser: configured.RunAsUser, RunAsGroup: configured.RunAsGroup, FSGroup: configured.FSGroup,
	})
}

func hardenedSecurityContext(isOpenShift bool, configured *corev1.PodSecurityContext) (*vmetricsv1b1.SecurityContext, error) {
	podSecurityContext, err := utils.HardenedPodSecurityContextWithPodOverrides(isOpenShift, configured)
	if err != nil {
		return nil, err
	}

	containerSecurityContext := utils.HardenedContainerSecurityContext()
	return &vmetricsv1b1.SecurityContext{
		PodSecurityContext: podSecurityContext,
		ContainerSecurityContext: &vmetricsv1b1.ContainerSecurityContext{
			AllowPrivilegeEscalation: containerSecurityContext.AllowPrivilegeEscalation,
			ReadOnlyRootFilesystem:   containerSecurityContext.ReadOnlyRootFilesystem,
			Capabilities:             containerSecurityContext.Capabilities,
		},
	}, nil
}

// EnsureTmpVolume returns a copy of the volumes containing the required size-limited temporary volume.
// It returns an error when a user-defined volume uses the reserved volume name.
func EnsureTmpVolume(volumes []corev1.Volume) ([]corev1.Volume, error) {
	return utils.EnsureTmpVolume(volumes, tmpVolumeSize)
}

// EnsureTmpVolumeMount returns a copy of the mounts containing the required temporary-directory mount.
func EnsureTmpVolumeMount(volumeMounts []corev1.VolumeMount) []corev1.VolumeMount {
	return utils.EnsureTmpVolumeMount(volumeMounts)
}

// HardenContainers returns hardened copies of explicitly configured application containers,
// each with the temporary-directory mount.
func HardenContainers(containers []corev1.Container) ([]corev1.Container, error) {
	return utils.HardenContainersWithTmp(containers)
}

// HardenContainerSecurity returns copies of explicitly configured application containers with
// the security baseline enforced and volume mounts left unchanged. VMAgent uses it because the
// VictoriaMetrics Operator keeps the persistent queue below /tmp of the main container.
func HardenContainerSecurity(containers []corev1.Container) ([]corev1.Container, error) {
	return utils.HardenContainers(containers)
}
