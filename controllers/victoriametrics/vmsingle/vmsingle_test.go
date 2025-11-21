package vmsingle

import (
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/stretchr/testify/assert"
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
				VmSingle: monv1.VmSingle{},
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
	})
	cr = &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
	}

}
