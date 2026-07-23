package vmagent

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

func TestVmAgentUsesOperatorRemoteWriteURLs(t *testing.T) {
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
				VmCluster: monv1.VmCluster{
					VmInsert:       &vmetricsv1b1.VMInsert{},
					VmSelectImage:  "vmselect:current",
					VmInsertImage:  "vminsert:current",
					VmStorageImage: "vmstorage:current",
				},
				VmAgent: monv1.VmAgent{
					Image: "vmagent:current",
				},
			},
		},
	}

	manifest, err := vmAgent(nil, cr)
	require.NoError(t, err)

	urls := make([]string, 0, len(manifest.Spec.RemoteWrite))
	for _, remoteWrite := range manifest.Spec.RemoteWrite {
		urls = append(urls, remoteWrite.URL)
	}
	assert.ElementsMatch(t, []string{
		"http://vmsingle-k8s.monitoring.svc:8428/api/v1/write",
		"http://vminsert-k8s.monitoring.svc:8480/insert/0/prometheus/api/v1/write",
	}, urls)
}
