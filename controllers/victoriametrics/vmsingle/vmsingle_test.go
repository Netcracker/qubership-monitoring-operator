package vmsingle

import (
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	vmetricsv1b1 "github.com/VictoriaMetrics/operator/api/operator/v1beta1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	cr *monv1.PlatformMonitoring
	// labelKey        = "label.key"
	// labelValue      = "label-value"
	// annotationKey   = "annotation.key"
	// annotationValue = "annotation-value"
)

func TestVmSingleManifests(t *testing.T) {
	// cr = &monv1.PlatformMonitoring{
	// 	ObjectMeta: metav1.ObjectMeta{
	// 		Namespace: "monitoring",
	// 	},
	// 	Spec: monv1.PlatformMonitoringSpec{
	// 		Victoriametrics: &v1.Victoriametrics{
	// 			vmSingle: v1.vmSingle{
	// 				Annotations: map[string]string{annotationKey: annotationValue},
	// 				Labels:      map[string]string{labelKey: labelValue},
	// 			},
	// 		},
	// 	},
	// }
	// t.Run("Test vmSingle manifest", func(t *testing.T) {
	// 	m, err := vmsingle(cr)
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}
	// 	assert.NotNil(t, m, "vmSingle manifest should not be empty")
	// })
	cr = &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
		Spec: monv1.PlatformMonitoringSpec{
			Victoriametrics: &monv1.Victoriametrics{
				VmSingle: monv1.VmSingle{Image: "example:v1"},
			},
		},
	}
	t.Run("Test vmSingle manifest with nil labels and annotation", func(t *testing.T) {
		m, err := vmSingle(nil, cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "vmSingle manifest should not be empty")
		assert.NotNil(t, m.GetLabels())
		assert.Nil(t, m.GetAnnotations())
		assertVmSingleHardening(t, m)
	})
	cr = &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
	}

}

func TestVmSingleUsesOperatorVmAlertURL(t *testing.T) {
	cr := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
		Spec: monv1.PlatformMonitoringSpec{
			Victoriametrics: &monv1.Victoriametrics{
				VmOperator: monv1.VmOperator{
					Image: "vmoperator:current",
				},
				VmSingle: monv1.VmSingle{
					Image: "vmsingle:current",
				},
				VmAlert: monv1.VmAlert{
					Image: "vmalert:current",
				},
			},
		},
	}

	manifest, err := vmSingle(nil, cr)
	require.NoError(t, err)

	assert.Equal(
		t,
		"http://vmalert-k8s.monitoring.svc:8080",
		manifest.Spec.ExtraArgs["vmalert.proxyURL"],
	)
}

func TestVmSingleClusterRBACUsesInstallationNamespace(t *testing.T) {
	cr := &monv1.PlatformMonitoring{ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"}}

	role, err := vmSingleClusterRole(cr, false, false)
	require.NoError(t, err)
	binding, err := vmSingleClusterRoleBinding(cr)
	require.NoError(t, err)

	assert.Equal(t, "monitoring", role.Labels["monitoring.netcracker.com/installation-namespace"])
	assert.Equal(t, "monitoring", binding.Labels["monitoring.netcracker.com/installation-namespace"])
}

func assertVmSingleHardening(t *testing.T, m *vmetricsv1b1.VMSingle) {
	t.Helper()
	require.NotNil(t, m.Spec.SecurityContext)
	require.NotNil(t, m.Spec.SecurityContext.RunAsNonRoot)
	require.NotNil(t, m.Spec.SecurityContext.AllowPrivilegeEscalation)
	require.NotNil(t, m.Spec.SecurityContext.ReadOnlyRootFilesystem)
	assert.Equal(t, true, *m.Spec.SecurityContext.RunAsNonRoot)
	assert.Equal(t, false, *m.Spec.SecurityContext.AllowPrivilegeEscalation)
	assert.Equal(t, true, *m.Spec.SecurityContext.ReadOnlyRootFilesystem)
	assert.Contains(t, m.Spec.Volumes, utils.TmpVolume("100Mi"))
	assert.Contains(t, m.Spec.VolumeMounts, utils.TmpVolumeMount())
}
