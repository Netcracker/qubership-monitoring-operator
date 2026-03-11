package pushgateway

import (
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils/labelsassert"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	cr         *monv1.PlatformMonitoring
	labelKey   = "label.key"
	labelValue = "label-value"
)

func TestPushgatewayManifests(t *testing.T) {
	cr = &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
	}
	t.Run("Test Deployment manifest", func(t *testing.T) {
		m, err := pushgatewayDeployment(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Deployment manifest should not be empty")
	})
	t.Run("Test Service manifest", func(t *testing.T) {
		m, err := pushgatewayService(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Service manifest should not be empty")
	})
	t.Run("Test Ingress v1beta1 manifest", func(t *testing.T) {
		m, err := pushgatewayIngressV1beta1(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Ingress v1beta1 manifest should not be empty")
	})
	t.Run("Test Ingress v1 manifest", func(t *testing.T) {
		m, err := pushgatewayIngressV1(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Ingress v1 manifest should not be empty")
	})
	t.Run("Test ServiceMonitor manifest", func(t *testing.T) {
		crWithLabels := &monv1.PlatformMonitoring{
			ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring", Labels: map[string]string{labelKey: labelValue}},
			Spec:      monv1.PlatformMonitoringSpec{Pushgateway: &monv1.Pushgateway{}},
		}
		m, err := pushgatewayServiceMonitor(crWithLabels)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "ServiceMonitor manifest should not be empty")
		labelsassert.AssertCRLabels(t, m.GetLabels(), utils.PushgatewayComponentName, "prometheus-operator", map[string]string{labelKey: labelValue})
	})
}
