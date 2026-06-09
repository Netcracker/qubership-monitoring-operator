package vmalert

import (
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	cr *monv1.PlatformMonitoring
)

func TestVmAlertManifests(t *testing.T) {
	cr = &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
		Spec: monv1.PlatformMonitoringSpec{
			Victoriametrics: &monv1.Victoriametrics{
				VmAlert: monv1.VmAlert{},
			},
		},
	}
	t.Run("Test vmAlert manifest with nil labels and annotation", func(t *testing.T) {
		m, err := vmAlert(nil, cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "vmAlert manifest should not be empty")
		assert.NotNil(t, m.GetLabels())
		assert.NotNil(t, m.GetAnnotations())
	})
	cr = &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
	}
}
