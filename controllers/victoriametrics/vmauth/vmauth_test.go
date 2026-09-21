package vmauth

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
	// labelKey        = "label.key"
	// labelValue      = "label-value"
	// annotationKey   = "annotation.key"
	// annotationValue = "annotation-value"
)

func TestVmAuthManifests(t *testing.T) {
	// cr = &monv1.PlatformMonitoring{
	// 	ObjectMeta: metav1.ObjectMeta{
	// 		Namespace: "monitoring",
	// 	},
	// 	Spec: monv1.PlatformMonitoringSpec{
	// 		Victoriametrics: &v1.Victoriametrics{
	// 			VmAgent: v1.VmAgent{
	// 				Annotations: map[string]string{annotationKey: annotationValue},
	// 				Labels:      map[string]string{labelKey: labelValue},
	// 			},
	// 		},
	// 	},
	// }
	// t.Run("Test VmAgent manifest", func(t *testing.T) {
	// 	m, err := vmagent(cr)
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}
	// 	assert.NotNil(t, m, "VmAgent manifest should not be empty")
	// 	assert.NotNil(t, m.Spec.PodMetadata.Labels)
	// 	assert.Equal(t, labelValue, m.Spec.PodMetadata.Labels[labelKey])
	// 	assert.NotNil(t, m.Spec.PodMetadata.Annotations)
	// 	assert.Equal(t, annotationValue, m.Spec.PodMetadata.Annotations[annotationKey])
	// })
	cr = &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
		Spec: monv1.PlatformMonitoringSpec{
			Victoriametrics: &monv1.Victoriametrics{
				VmAuth: monv1.VmAuth{Image: "example:v1"},
			},
		},
	}
	t.Run("Test VmAuth manifest with nil labels and annotation", func(t *testing.T) {
		m, err := vmAuth(nil, cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "VmAuth manifest should not be empty")
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

func TestVmAuthClusterRBACUsesInstallationNamespace(t *testing.T) {
	cr := &monv1.PlatformMonitoring{ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"}}

	role, err := vmAuthClusterRole(cr, false, false)
	require.NoError(t, err)
	binding, err := vmAuthClusterRoleBinding(cr)
	require.NoError(t, err)

	assert.Equal(t, "monitoring", role.Labels["monitoring.netcracker.com/installation-namespace"])
	assert.Equal(t, "monitoring", binding.Labels["monitoring.netcracker.com/installation-namespace"])
}
