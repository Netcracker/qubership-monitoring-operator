package vmuser

import (
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
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

func TestVmUserManifests(t *testing.T) {
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
				VmUser: monv1.VmUser{},
			},
		},
	}
	t.Run("Test VmUser manifest with nil labels and annotation", func(t *testing.T) {
		m, err := vmUser(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "VmUser manifest should not be empty")
		assert.NotNil(t, m.GetLabels())
		assert.Nil(t, m.GetAnnotations())
	})
	cr = &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
	}

}

func TestVmUserTargetsUseNamespacedReferences(t *testing.T) {
	install := true
	username := "reader"
	cr := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
		Spec: monv1.PlatformMonitoringSpec{
			Victoriametrics: &monv1.Victoriametrics{
				VmUser: monv1.VmUser{
					Install:  &install,
					Image:    "vmuser:current",
					UserName: &username,
				},
				VmAuth: monv1.VmAuth{
					Install: &install,
					Image:   "vmauth:current",
				},
				VmSingle: monv1.VmSingle{
					Install: &install,
					Image:   "vmsingle:current",
				},
				VmAgent: monv1.VmAgent{
					Install: &install,
					Image:   "vmagent:current",
				},
				VmAlert: monv1.VmAlert{
					Install: &install,
					Image:   "vmalert:current",
				},
				VmAlertManager: monv1.VmAlertManager{
					Install: &install,
					Image:   "vmalertmanager:current",
				},
				VmCluster: monv1.VmCluster{
					Install:        &install,
					VmSelect:       &vmetricsv1b1.VMSelect{},
					VmSelectImage:  "vmselect:current",
					VmInsert:       &vmetricsv1b1.VMInsert{},
					VmInsertImage:  "vminsert:current",
					VmStorage:      &vmetricsv1b1.VMStorage{},
					VmStorageImage: "vmstorage:current",
				},
			},
		},
	}

	manifest, err := vmUser(cr)
	require.NoError(t, err)
	assert.Equal(t, &username, manifest.Spec.Username)

	kinds := make([]string, 0, len(manifest.Spec.TargetRefs))
	for _, target := range manifest.Spec.TargetRefs {
		require.NotNil(t, target.CRD)
		assert.Equal(t, "k8s", target.CRD.Name)
		assert.Equal(t, "monitoring", target.CRD.Namespace)
		kinds = append(kinds, target.CRD.Kind)
	}
	assert.ElementsMatch(t, []string{
		"VMSingle",
		"VMCluster/vmselect",
		"VMCluster/vmstorage",
		"VMCluster/vminsert",
		"VMAlert",
		"VMAlertmanager",
		"VMAgent",
		"VMSingle",
	}, kinds)

	cr.Spec.Victoriametrics.VmSingle.Image = ""
	manifest, err = vmUser(cr)
	require.NoError(t, err)
	require.NotEmpty(t, manifest.Spec.TargetRefs)
	require.NotNil(t, manifest.Spec.TargetRefs[0].CRD)
	assert.Equal(t, "VMAgent", manifest.Spec.TargetRefs[0].CRD.Kind)
	assert.Equal(t, "k8s", manifest.Spec.TargetRefs[0].CRD.Name)
	assert.Equal(t, "monitoring", manifest.Spec.TargetRefs[0].CRD.Namespace)
}
