package utils

import (
	corev1 "k8s.io/api/core/v1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/utils/ptr"
)

const (
	defaultVolumeMode                    int32 = 0644
	defaultServiceAccountTokenExpiration int64 = 3600
)

// SetVolumesIfChanged assigns src to *dst unless *dst equals src either as
// written or after the API server defaults are applied to src. Volumes from
// the PlatformMonitoring resource usually omit fields such as defaultMode, and
// the stored object always has them. Returns true when *dst was updated.
func SetVolumesIfChanged(dst *[]corev1.Volume, src []corev1.Volume) bool {
	if apiequality.Semantic.DeepEqual(*dst, src) ||
		apiequality.Semantic.DeepEqual(*dst, VolumesWithDefaults(src)) {
		return false
	}
	*dst = src
	return true
}

// VolumesWithDefaults returns a copy of volumes with the defaults that the
// API server applies to pod volumes. Explicitly set values are kept.
func VolumesWithDefaults(volumes []corev1.Volume) []corev1.Volume {
	if volumes == nil {
		return nil
	}
	result := make([]corev1.Volume, len(volumes))
	for i := range volumes {
		volumes[i].DeepCopyInto(&result[i])
		setVolumeSourceDefaults(&result[i].VolumeSource)
	}
	return result
}

func setVolumeSourceDefaults(source *corev1.VolumeSource) {
	if apiequality.Semantic.DeepEqual(*source, corev1.VolumeSource{}) {
		source.EmptyDir = &corev1.EmptyDirVolumeSource{}
		return
	}
	if source.Secret != nil {
		setPtrDefault(&source.Secret.DefaultMode, defaultVolumeMode)
	}
	if source.ConfigMap != nil {
		setPtrDefault(&source.ConfigMap.DefaultMode, defaultVolumeMode)
	}
	if source.DownwardAPI != nil {
		setPtrDefault(&source.DownwardAPI.DefaultMode, defaultVolumeMode)
		setDownwardAPIItemDefaults(source.DownwardAPI.Items)
	}
	if source.Projected != nil {
		setProjectedVolumeDefaults(source.Projected)
	}
	if source.HostPath != nil {
		setPtrDefault(&source.HostPath.Type, corev1.HostPathUnset)
	}
	if source.ISCSI != nil {
		setValueDefault(&source.ISCSI.ISCSIInterface, "default")
	}
	if source.RBD != nil {
		setValueDefault(&source.RBD.RBDPool, "rbd")
		setValueDefault(&source.RBD.RadosUser, "admin")
		setValueDefault(&source.RBD.Keyring, "/etc/ceph/keyring")
	}
	if source.AzureDisk != nil {
		setPtrDefault(&source.AzureDisk.CachingMode, corev1.AzureDataDiskCachingReadWrite)
		setPtrDefault(&source.AzureDisk.Kind, corev1.AzureSharedBlobDisk)
		setPtrDefault(&source.AzureDisk.FSType, "ext4")
		setPtrDefault(&source.AzureDisk.ReadOnly, false)
	}
	if source.ScaleIO != nil {
		setValueDefault(&source.ScaleIO.StorageMode, "ThinProvisioned")
		setValueDefault(&source.ScaleIO.FSType, "xfs")
	}
	if source.Ephemeral != nil && source.Ephemeral.VolumeClaimTemplate != nil {
		setPtrDefault(&source.Ephemeral.VolumeClaimTemplate.Spec.VolumeMode, corev1.PersistentVolumeFilesystem)
	}
}

func setProjectedVolumeDefaults(projected *corev1.ProjectedVolumeSource) {
	setPtrDefault(&projected.DefaultMode, defaultVolumeMode)
	for i := range projected.Sources {
		projection := &projected.Sources[i]
		if projection.DownwardAPI != nil {
			setDownwardAPIItemDefaults(projection.DownwardAPI.Items)
		}
		if projection.ServiceAccountToken != nil {
			setPtrDefault(&projection.ServiceAccountToken.ExpirationSeconds, defaultServiceAccountTokenExpiration)
		}
	}
}

func setDownwardAPIItemDefaults(items []corev1.DownwardAPIVolumeFile) {
	for i := range items {
		if items[i].FieldRef != nil {
			setValueDefault(&items[i].FieldRef.APIVersion, "v1")
		}
	}
}

func setPtrDefault[T any](field **T, value T) {
	if *field == nil {
		*field = ptr.To(value)
	}
}

func setValueDefault[T comparable](field *T, value T) {
	var zero T
	if *field == zero {
		*field = value
	}
}
