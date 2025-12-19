package prometheus_rules

import (
	"testing"

	v1beta1 "github.com/Netcracker/qubership-monitoring-operator/api"
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

func TestPrometheusRuleManifests(t *testing.T) {
	cr = &v1beta1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:   "monitoring",
			Annotations: map[string]string{annotationKey: annotationValue},
			Labels:      map[string]string{labelKey: labelValue},
		},
	}
	t.Run("Test PrometheusRule manifest", func(t *testing.T) {
		m, err := prometheusRules(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "PrometheusRule manifest should not be empty")
		assert.NotNil(t, m.GetLabels())
		assert.Equal(t, labelValue, m.GetLabels()[labelKey])
		assert.NotNil(t, m.GetAnnotations())
		assert.Equal(t, annotationValue, m.GetAnnotations()[annotationKey])
	})
	cr = &v1beta1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
	}
	t.Run("Test PrometheusRule manifest with nil labels and annotation", func(t *testing.T) {
		m, err := prometheusRules(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "PrometheusRule manifest should not be empty")
		assert.NotNil(t, m.GetLabels())
		assert.Nil(t, m.GetAnnotations())
	})
}
