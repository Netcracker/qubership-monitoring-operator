package utils

import (
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
