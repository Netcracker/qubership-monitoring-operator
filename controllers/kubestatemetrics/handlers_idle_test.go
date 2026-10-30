package kubestatemetrics

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
	r := NewKubeStateMetricsReconciler(reconciletest.NewClient(scheme, counter), scheme, &fakediscovery.FakeDiscovery{Fake: &k8stesting.Fake{}})

	reconciletest.AssertIdleReconcileSkipsUpdates(t, skipTestPlatformMonitoring(), counter, map[string]reconciletest.Handler{
		"ServiceAccount":     r.handleServiceAccount,
		"ClusterRole":        r.handleClusterRole,
		"ClusterRoleBinding": r.handleClusterRoleBinding,
		"Deployment":         r.handleDeployment,
		"Service":            r.handleService,
		"ServiceMonitor":     r.handleServiceMonitor,
	})
}
