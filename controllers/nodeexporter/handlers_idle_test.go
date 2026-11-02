package nodeexporter

import (
	"testing"

	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils/reconciletest"
	secv1 "github.com/openshift/api/security/v1"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	fakediscovery "k8s.io/client-go/discovery/fake"
	k8stesting "k8s.io/client-go/testing"
)

func TestHandlersSkipIdleUpdates(t *testing.T) {
	scheme := reconciletest.NewScheme(t, promv1.AddToScheme, secv1.Install)
	counter := &reconciletest.UpdateCounter{}
	r := NewNodeExporterReconciler(reconciletest.NewClient(scheme, counter), scheme, &fakediscovery.FakeDiscovery{Fake: &k8stesting.Fake{}})

	reconciletest.AssertIdleReconcileSkipsUpdates(t, newNodeExporterPlatformMonitoring(true), counter, map[string]reconciletest.Handler{
		"ServiceAccount":             r.handleServiceAccount,
		"ClusterRole":                r.handleClusterRole,
		"ClusterRoleBinding":         r.handleClusterRoleBinding,
		"SecurityContextConstraints": r.handleSecurityContextConstraints,
		"DaemonSet":                  r.handleDaemonSet,
		"Service":                    r.handleService,
		"ServiceMonitor":             r.handleServiceMonitor,
	})
}
