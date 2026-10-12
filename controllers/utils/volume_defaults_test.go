package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/ptr"
)

func TestVolumesWithDefaults(t *testing.T) {
	fieldRef := func(apiVersion string) []corev1.DownwardAPIVolumeFile {
		return []corev1.DownwardAPIVolumeFile{
			{Path: "labels", FieldRef: &corev1.ObjectFieldSelector{APIVersion: apiVersion, FieldPath: "metadata.labels"}},
			{Path: "cpu", ResourceFieldRef: &corev1.ResourceFieldSelector{Resource: "limits.cpu"}},
		}
	}
	tests := []struct {
		name string
		in   corev1.VolumeSource
		want corev1.VolumeSource
	}{
		{
			name: "empty source",
			in:   corev1.VolumeSource{},
			want: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}},
		},
		{
			name: "secret",
			in:   corev1.VolumeSource{Secret: &corev1.SecretVolumeSource{SecretName: "tls"}},
			want: corev1.VolumeSource{Secret: &corev1.SecretVolumeSource{SecretName: "tls", DefaultMode: ptr.To(int32(420))}},
		},
		{
			name: "secret with explicit mode",
			in:   corev1.VolumeSource{Secret: &corev1.SecretVolumeSource{SecretName: "tls", DefaultMode: ptr.To(int32(0400))}},
			want: corev1.VolumeSource{Secret: &corev1.SecretVolumeSource{SecretName: "tls", DefaultMode: ptr.To(int32(0400))}},
		},
		{
			name: "config map",
			in:   corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{}},
			want: corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{DefaultMode: ptr.To(int32(420))}},
		},
		{
			name: "downward API",
			in:   corev1.VolumeSource{DownwardAPI: &corev1.DownwardAPIVolumeSource{Items: fieldRef("")}},
			want: corev1.VolumeSource{DownwardAPI: &corev1.DownwardAPIVolumeSource{
				DefaultMode: ptr.To(int32(420)),
				Items:       fieldRef("v1"),
			}},
		},
		{
			name: "projected",
			in: corev1.VolumeSource{Projected: &corev1.ProjectedVolumeSource{Sources: []corev1.VolumeProjection{
				{ServiceAccountToken: &corev1.ServiceAccountTokenProjection{Path: "token"}},
				{ServiceAccountToken: &corev1.ServiceAccountTokenProjection{Path: "short", ExpirationSeconds: ptr.To(int64(600))}},
				{DownwardAPI: &corev1.DownwardAPIProjection{Items: fieldRef("")}},
				{ConfigMap: &corev1.ConfigMapProjection{}},
			}}},
			want: corev1.VolumeSource{Projected: &corev1.ProjectedVolumeSource{
				DefaultMode: ptr.To(int32(420)),
				Sources: []corev1.VolumeProjection{
					{ServiceAccountToken: &corev1.ServiceAccountTokenProjection{Path: "token", ExpirationSeconds: ptr.To(int64(3600))}},
					{ServiceAccountToken: &corev1.ServiceAccountTokenProjection{Path: "short", ExpirationSeconds: ptr.To(int64(600))}},
					{DownwardAPI: &corev1.DownwardAPIProjection{Items: fieldRef("v1")}},
					{ConfigMap: &corev1.ConfigMapProjection{}},
				},
			}},
		},
		{
			name: "host path",
			in:   corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: "/proc"}},
			want: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: "/proc", Type: ptr.To(corev1.HostPathUnset)}},
		},
		{
			name: "iSCSI",
			in:   corev1.VolumeSource{ISCSI: &corev1.ISCSIVolumeSource{IQN: "iqn"}},
			want: corev1.VolumeSource{ISCSI: &corev1.ISCSIVolumeSource{IQN: "iqn", ISCSIInterface: "default"}},
		},
		{
			name: "RBD",
			in:   corev1.VolumeSource{RBD: &corev1.RBDVolumeSource{RBDImage: "image", RadosUser: "user"}},
			want: corev1.VolumeSource{RBD: &corev1.RBDVolumeSource{
				RBDImage: "image", RBDPool: "rbd", RadosUser: "user", Keyring: "/etc/ceph/keyring",
			}},
		},
		{
			name: "Azure disk",
			in:   corev1.VolumeSource{AzureDisk: &corev1.AzureDiskVolumeSource{DiskName: "disk"}},
			want: corev1.VolumeSource{AzureDisk: &corev1.AzureDiskVolumeSource{
				DiskName:    "disk",
				CachingMode: ptr.To(corev1.AzureDataDiskCachingReadWrite),
				Kind:        ptr.To(corev1.AzureSharedBlobDisk),
				FSType:      ptr.To("ext4"),
				ReadOnly:    ptr.To(false),
			}},
		},
		{
			name: "ScaleIO",
			in:   corev1.VolumeSource{ScaleIO: &corev1.ScaleIOVolumeSource{System: "system"}},
			want: corev1.VolumeSource{ScaleIO: &corev1.ScaleIOVolumeSource{System: "system", StorageMode: "ThinProvisioned", FSType: "xfs"}},
		},
		{
			name: "ephemeral",
			in:   corev1.VolumeSource{Ephemeral: &corev1.EphemeralVolumeSource{VolumeClaimTemplate: &corev1.PersistentVolumeClaimTemplate{}}},
			want: corev1.VolumeSource{Ephemeral: &corev1.EphemeralVolumeSource{VolumeClaimTemplate: &corev1.PersistentVolumeClaimTemplate{
				Spec: corev1.PersistentVolumeClaimSpec{VolumeMode: ptr.To(corev1.PersistentVolumeFilesystem)},
			}}},
		},
		{
			name: "persistent volume claim",
			in:   corev1.VolumeSource{PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{ClaimName: "data"}},
			want: corev1.VolumeSource{PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{ClaimName: "data"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := []corev1.Volume{{Name: "volume", VolumeSource: tt.in}}
			original := in[0].DeepCopy()

			got := VolumesWithDefaults(in)

			assert.Equal(t, []corev1.Volume{{Name: "volume", VolumeSource: tt.want}}, got)
			assert.Equal(t, *original, in[0], "input must not be modified")
		})
	}

	assert.Nil(t, VolumesWithDefaults(nil))
}

func TestSetVolumesIfChanged(t *testing.T) {
	desired := []corev1.Volume{{
		Name:         "config",
		VolumeSource: corev1.VolumeSource{ConfigMap: &corev1.ConfigMapVolumeSource{LocalObjectReference: corev1.LocalObjectReference{Name: "config"}}},
	}}
	withMode := func(mode int32) []corev1.Volume {
		volumes := VolumesWithDefaults(desired)
		volumes[0].ConfigMap.DefaultMode = ptr.To(mode)
		return volumes
	}

	t.Run("skips an identical list", func(t *testing.T) {
		live := VolumesWithDefaults(desired)
		live[0].ConfigMap.DefaultMode = nil
		assert.False(t, SetVolumesIfChanged(&live, desired))
	})

	t.Run("skips an API-defaulted list", func(t *testing.T) {
		live := withMode(420)
		assert.False(t, SetVolumesIfChanged(&live, desired))
		assert.Equal(t, withMode(420), live)
	})

	t.Run("updates a changed mode", func(t *testing.T) {
		live := withMode(0400)
		assert.True(t, SetVolumesIfChanged(&live, desired))
		assert.Equal(t, desired, live)
	})

	t.Run("updates when an explicit mode differs from the default", func(t *testing.T) {
		live := withMode(420)
		assert.True(t, SetVolumesIfChanged(&live, withMode(0400)))
		assert.Equal(t, withMode(0400), live)
	})

	t.Run("updates a removed volume", func(t *testing.T) {
		live := withMode(420)
		assert.True(t, SetVolumesIfChanged(&live, nil))
		assert.Empty(t, live)
	})
}
