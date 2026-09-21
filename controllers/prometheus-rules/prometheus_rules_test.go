package prometheus_rules

import (
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	cr              *monv1.PlatformMonitoring
	labelKey        = "label.key"
	labelValue      = "label-value"
	annotationKey   = "annotation.key"
	annotationValue = "annotation-value"
)

func TestPrometheusRuleManifests(t *testing.T) {
	cr = &monv1.PlatformMonitoring{
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
	cr = &monv1.PlatformMonitoring{
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

func TestPrometheusRuleManifestWithAnnotationOverride(t *testing.T) {
	install := true
	customSummary := "Custom notification backlog summary"
	customRunbook := "https://example.org/runbooks/prometheus-backlog"
	cr := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"},
		Spec: monv1.PlatformMonitoringSpec{
			PrometheusRules: &monv1.PrometheusRules{
				Install:    &install,
				RuleGroups: []string{"SelfMonitoring"},
				Override: []monv1.PrometheusRule{
					{
						Group: "SelfMonitoring",
						Alert: "PrometheusNotificationsBacklog",
						Annotations: map[string]string{
							"summary":     customSummary,
							"runbook_url": customRunbook,
						},
					},
				},
			},
		},
	}

	manifest, err := prometheusRules(cr)
	if err != nil {
		t.Fatal(err)
	}

	for _, group := range manifest.Spec.Groups {
		for _, rule := range group.Rules {
			if rule.Alert == "PrometheusNotificationsBacklog" {
				assert.Equal(t, customSummary, rule.Annotations["summary"])
				assert.Equal(t, customRunbook, rule.Annotations["runbook_url"])
				assert.NotEmpty(t, rule.Annotations["description"])
				return
			}
		}
	}

	t.Fatal("PrometheusNotificationsBacklog rule was not found")
}
