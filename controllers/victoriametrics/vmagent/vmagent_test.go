package vmagent

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

func TestVmAgentManifests(t *testing.T) {
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
				VmAgent: monv1.VmAgent{},
			},
		},
	}
	t.Run("Test Vmagent manifest with nil labels and annotation", func(t *testing.T) {
		m, err := vmAgent(nil, cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Vmagent manifest should not be empty")
		assert.NotNil(t, m.GetLabels())
		assert.Nil(t, m.GetAnnotations())
	})
	t.Run("Test Vmagent manifest with nil maxScrapeInternal and minScrapeInternal", func(t *testing.T) {
		m, err := vmAgent(nil, cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Vmagent manifest should not be empty")
		assert.Nil(t, m.Spec.MaxScrapeInterval)
		assert.Nil(t, m.Spec.MinScrapeInterval)
	})
	cr = &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
	}

}
