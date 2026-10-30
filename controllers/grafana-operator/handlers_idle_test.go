package grafana_operator

import (
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils/reconciletest"
	grafv1 "github.com/grafana/grafana-operator/v5/api/v1beta1"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	fakediscovery "k8s.io/client-go/discovery/fake"
	k8stesting "k8s.io/client-go/testing"
)

func TestHandlersSkipIdleUpdates(t *testing.T) {
	scheme := reconciletest.NewScheme(t, promv1.AddToScheme, grafv1.AddToScheme)
	counter := &reconciletest.UpdateCounter{}
	r := NewGrafanaOperatorReconciler(reconciletest.NewClient(scheme, counter), scheme, &fakediscovery.FakeDiscovery{Fake: &k8stesting.Fake{}})
	cr := testPlatformMonitoring()
	cr.Spec.Grafana.Operator.Image = "ghcr.io/grafana/grafana-operator:v5.18.0"

	reconciletest.AssertIdleReconcileSkipsUpdates(t, cr, counter, map[string]reconciletest.Handler{
		"ServiceAccount":     r.handleServiceAccount,
		"ClusterRole":        r.handleClusterRole,
		"ClusterRoleBinding": r.handleClusterRoleBinding,
		"Role":               r.handleRole,
		"RoleBinding":        r.handleRoleBinding,
		"Deployment":         r.handleDeployment,
		"PodMonitor":         r.handlePodMonitor,
		"GrafanaDashboard": func(cr *monv1.PlatformMonitoring) error {
			return r.handleGrafanaDashboard(testDashboardAsset, cr)
		},
	})
}
