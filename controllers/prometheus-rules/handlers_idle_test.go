package prometheus_rules

import (
	"testing"

	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils/reconciletest"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
)

func TestHandlersSkipIdleUpdates(t *testing.T) {
	scheme := reconciletest.NewScheme(t, promv1.AddToScheme)
	counter := &reconciletest.UpdateCounter{}
	r := NewPrometheusRulesReconciler(reconciletest.NewClient(scheme, counter), scheme)

	reconciletest.AssertIdleReconcileSkipsUpdates(t, skipTestPlatformMonitoring(), counter, map[string]reconciletest.Handler{
		"PrometheusRule": r.handlePrometheusRules,
	})
}
