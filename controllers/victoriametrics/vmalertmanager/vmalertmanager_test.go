package vmalertmanager

import (
	"testing"

	v1beta1 "github.com/Netcracker/qubership-monitoring-operator/api"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	cr *v1beta1.PlatformMonitoring
)

func TestVmAlertManagerManifests(t *testing.T) {
	cr = &v1beta1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
		Spec: v1beta1.PlatformMonitoringSpec{
			Victoriametrics: &v1beta1.Victoriametrics{
				VmAlertManager: v1beta1.VmAlertManager{},
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
	})
	cr = &v1beta1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
	}
}
