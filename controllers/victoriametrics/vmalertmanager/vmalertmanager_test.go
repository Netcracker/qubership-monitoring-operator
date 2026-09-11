package vmalertmanager

import (
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	cr *monv1.PlatformMonitoring
)

func TestVmAlertManagerManifests(t *testing.T) {
	cr = &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
		Spec: monv1.PlatformMonitoringSpec{
			Victoriametrics: &monv1.Victoriametrics{
				VmAlertManager: monv1.VmAlertManager{Image: "example:v1"},
			},
		},
	}
	t.Run("Test vmAlertManager manifest with nil labels and annotation", func(t *testing.T) {
		m, err := vmAlertManager(nil, cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "vmAlertManager manifest should not be empty")
		assert.NotNil(t, m.GetLabels())
		assert.Nil(t, m.GetAnnotations())
		require.NotNil(t, m.Spec.SecurityContext)
		require.NotNil(t, m.Spec.SecurityContext.RunAsNonRoot)
		require.NotNil(t, m.Spec.SecurityContext.AllowPrivilegeEscalation)
		require.NotNil(t, m.Spec.SecurityContext.ReadOnlyRootFilesystem)
		assert.Equal(t, true, *m.Spec.SecurityContext.RunAsNonRoot)
		assert.Equal(t, false, *m.Spec.SecurityContext.AllowPrivilegeEscalation)
		assert.Equal(t, true, *m.Spec.SecurityContext.ReadOnlyRootFilesystem)
		assert.Contains(t, m.Spec.Volumes, utils.TmpVolume("100Mi"))
		assert.Contains(t, m.Spec.VolumeMounts, utils.TmpVolumeMount())
	})
	cr = &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
	}
}

func TestVmAlertManagerUsesTopLevelHTTPSProbes(t *testing.T) {
	cr := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
		Spec: monv1.PlatformMonitoringSpec{
			Victoriametrics: &monv1.Victoriametrics{
				TLSEnabled: true,
				VmAlertManager: monv1.VmAlertManager{
					Image: "vmalertmanager:current",
				},
			},
		},
	}

	manifest, err := vmAlertManager(nil, cr)
	require.NoError(t, err)
	require.NotNil(t, manifest.Spec.WebConfig)
	require.NotNil(t, manifest.Spec.WebConfig.TLSServerConfig)
	require.NotNil(t, manifest.Spec.LivenessProbe)
	require.NotNil(t, manifest.Spec.ReadinessProbe)

	assert.Equal(t, "/-/healthy", manifest.Spec.LivenessProbe.HTTPGet.Path)
	assert.Equal(t, "web", manifest.Spec.LivenessProbe.HTTPGet.Port.StrVal)
	assert.Equal(t, "HTTPS", string(manifest.Spec.LivenessProbe.HTTPGet.Scheme))
	assert.Equal(t, "/-/healthy", manifest.Spec.ReadinessProbe.HTTPGet.Path)
	assert.Equal(t, "web", manifest.Spec.ReadinessProbe.HTTPGet.Port.StrVal)
	assert.Equal(t, "HTTPS", string(manifest.Spec.ReadinessProbe.HTTPGet.Scheme))
}

func TestVmAlertManagerClusterRBACUsesInstallationNamespace(t *testing.T) {
	cr := &monv1.PlatformMonitoring{ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"}}

	role, err := vmAlertManagerClusterRole(cr, false, false)
	require.NoError(t, err)
	binding, err := vmAlertManagerClusterRoleBinding(cr)
	require.NoError(t, err)

	assert.Equal(t, "monitoring", role.Labels["monitoring.netcracker.com/installation-namespace"])
	assert.Equal(t, "monitoring", binding.Labels["monitoring.netcracker.com/installation-namespace"])
}
