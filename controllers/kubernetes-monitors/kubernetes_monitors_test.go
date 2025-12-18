package kubernetes_monitors

import (
	"testing"

	v1beta1 "github.com/Netcracker/qubership-monitoring-operator/api/v1beta1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	cr              *v1beta1.PlatformMonitoring
	labelKey        = "label.key"
	labelValue      = "label-value"
	annotationKey   = "annotation.key"
	annotationValue = "annotation-value"
)

func TestKubernetesMonitorsManifests(t *testing.T) {
	cr = &v1beta1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
	}
	t.Run("Test ApiServerServiceMonitor manifest with nil labels and annotation", func(t *testing.T) {
		m, err := kubernetesMonitorsApiServerServiceMonitor(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "ApiServerServiceMonitor manifest should not be empty")
		assert.NotNil(t, m.GetLabels())
		assert.Nil(t, m.GetAnnotations())
	})
	cr = &v1beta1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:   "monitoring",
			Annotations: map[string]string{annotationKey: annotationValue},
			Labels:      map[string]string{labelKey: labelValue},
		},
	}
	t.Run("Test ApiServerServiceMonitor manifest", func(t *testing.T) {
		m, err := kubernetesMonitorsApiServerServiceMonitor(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "ApiServerServiceMonitor manifest should not be empty")
		assert.NotNil(t, m.GetLabels())
		assert.Equal(t, labelValue, m.GetLabels()[labelKey])
		assert.NotNil(t, m.GetAnnotations())
		assert.Equal(t, annotationValue, m.GetAnnotations()[annotationKey])
	})
	t.Run("Test ControllerManagerServiceMonitor manifest", func(t *testing.T) {
		m, err := kubernetesMonitorsControllerManagerServiceMonitor(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "ControllerManagerServiceMonitor manifest should not be empty")
		assert.NotNil(t, m.GetLabels())
		assert.Equal(t, labelValue, m.GetLabels()[labelKey])
		assert.NotNil(t, m.GetAnnotations())
		assert.Equal(t, annotationValue, m.GetAnnotations()[annotationKey])
	})
	t.Run("Test SchedulerServiceMonitor manifest", func(t *testing.T) {
		m, err := kubernetesMonitorsSchedulerServiceMonitor(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "SchedulerServiceMonitor manifest should not be empty")
		assert.NotNil(t, m.GetLabels())
		assert.Equal(t, labelValue, m.GetLabels()[labelKey])
		assert.NotNil(t, m.GetAnnotations())
		assert.Equal(t, annotationValue, m.GetAnnotations()[annotationKey])
	})
	t.Run("Test KubeletServiceMonitor manifest", func(t *testing.T) {
		m, err := kubernetesMonitorsKubeletServiceMonitor(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "KubeletServiceMonitor manifest should not be empty")
		assert.NotNil(t, m.GetLabels())
		assert.Equal(t, labelValue, m.GetLabels()[labelKey])
		assert.NotNil(t, m.GetAnnotations())
		assert.Equal(t, annotationValue, m.GetAnnotations()[annotationKey])
	})
	t.Run("Test CoreDnsServiceMonitor manifest for K8s", func(t *testing.T) {
		m, err := kubernetesMonitorsCoreDnsServiceMonitor(cr, false)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "CoreDnsServiceMonitor manifest should not be empty")
		assert.NotNil(t, m.GetLabels())
		assert.Equal(t, labelValue, m.GetLabels()[labelKey])
		assert.NotNil(t, m.GetAnnotations())
		assert.Equal(t, annotationValue, m.GetAnnotations()[annotationKey])
	})
	t.Run("Test CoreDnsServiceMonitor manifest for OS4", func(t *testing.T) {
		m, err := kubernetesMonitorsCoreDnsServiceMonitor(cr, true)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "CoreDnsServiceMonitor manifest should not be empty")
		assert.NotNil(t, m.GetLabels())
		assert.Equal(t, labelValue, m.GetLabels()[labelKey])
		assert.NotNil(t, m.GetAnnotations())
		assert.Equal(t, annotationValue, m.GetAnnotations()[annotationKey])
	})
	t.Run("Test NginxIngressPodMonitor manifest", func(t *testing.T) {
		m, err := kubernetesMonitorsNginxIngressPodMonitor(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "NginxIngressPodMonitor manifest should not be empty")
		assert.NotNil(t, m.GetLabels())
		assert.Equal(t, labelValue, m.GetLabels()[labelKey])
		assert.NotNil(t, m.GetAnnotations())
		assert.Equal(t, annotationValue, m.GetAnnotations()[annotationKey])
	})
	t.Run("Test OpenshiftApiserverServiceMonitor manifest", func(t *testing.T) {
		m, err := openshiftServiceMonitor(cr, utils.OpenshiftApiServerServiceMonitorAsset,
			utils.OpenshiftApiServerServiceMonitorName, utils.OpenshiftApiserver)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "OpenshiftApiserverServiceMonitor manifest should not be empty")
		assert.NotNil(t, m.GetLabels())
		assert.Equal(t, labelValue, m.GetLabels()[labelKey])
		assert.NotNil(t, m.GetAnnotations())
		assert.Equal(t, annotationValue, m.GetAnnotations()[annotationKey])
	})
	t.Run("Test OpenshiftApiserverOperatorServiceMonitor manifest", func(t *testing.T) {
		m, err := openshiftServiceMonitor(cr, utils.OpenshiftApiServerOperatorServiceMonitorAsset,
			utils.OpenshiftApiServerOperatorServiceMonitorName, utils.OpenshiftApiServerOperator)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "OpenshiftApiserverOperatorServiceMonitor manifest should not be empty")
		assert.NotNil(t, m.GetLabels())
		assert.Equal(t, labelValue, m.GetLabels()[labelKey])
		assert.NotNil(t, m.GetAnnotations())
		assert.Equal(t, annotationValue, m.GetAnnotations()[annotationKey])
	})
	t.Run("Test OpenshiftClusterVersionOperatorServiceMonitor", func(t *testing.T) {
		m, err := openshiftServiceMonitor(cr, utils.OpenshiftClusterVersionOperatorServiceMonitorAsset,
			utils.OpenshiftClusterVersionOperatorServiceMonitorName, utils.OpenshiftClusterVersionOperator)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "OpenshiftClusterVersionOperatorServiceMonitor manifest should not be empty")
		assert.NotNil(t, m.GetLabels())
		assert.Equal(t, labelValue, m.GetLabels()[labelKey])
		assert.NotNil(t, m.GetAnnotations())
		assert.Equal(t, annotationValue, m.GetAnnotations()[annotationKey])
	})
	t.Run("Test openshiftStatemetricsServiceMonitor", func(t *testing.T) {
		m, err := openshiftServiceMonitor(cr, utils.OpenshiftStatemetricsServiceMonitorAsset,
			utils.OpenshiftStatemetricsServiceMonitorName, utils.OpenshiftStatemetrics)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "openshiftStatemetricsServiceMonitor manifest should not be empty")
		assert.NotNil(t, m.GetLabels())
		assert.Equal(t, labelValue, m.GetLabels()[labelKey])
		assert.NotNil(t, m.GetAnnotations())
		assert.Equal(t, annotationValue, m.GetAnnotations()[annotationKey])
	})
}
