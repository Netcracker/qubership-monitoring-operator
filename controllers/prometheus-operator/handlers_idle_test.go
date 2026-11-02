package prometheus_operator

import (
	"testing"

	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils/reconciletest"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	fakediscovery "k8s.io/client-go/discovery/fake"
	k8stesting "k8s.io/client-go/testing"
)

func TestHandlersSkipIdleUpdates(t *testing.T) {
	scheme := reconciletest.NewScheme(t, promv1.AddToScheme)
	counter := &reconciletest.UpdateCounter{}
	r := NewPrometheusOperatorReconciler(reconciletest.NewClient(scheme, counter), scheme, &fakediscovery.FakeDiscovery{Fake: &k8stesting.Fake{}})

	reconciletest.AssertIdleReconcileSkipsUpdates(t, skipTestPlatformMonitoring(), counter, map[string]reconciletest.Handler{
		"ServiceAccount":     r.handleServiceAccount,
		"Role":               r.handleRole,
		"RoleBinding":        r.handleRoleBinding,
		"ClusterRole":        r.handleClusterRole,
		"ClusterRoleBinding": r.handleClusterRoleBinding,
		"Deployment":         r.handleDeployment,
		"Service":            r.handleService,
		"PodMonitor":         r.handlePodMonitor,
	})
}
