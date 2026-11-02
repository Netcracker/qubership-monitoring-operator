package vmsingle

import (
	"testing"

	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils/reconciletest"
	vmetricsv1b1 "github.com/VictoriaMetrics/operator/api/operator/v1beta1"
	fakediscovery "k8s.io/client-go/discovery/fake"
	ktesting "k8s.io/client-go/testing"
)

func TestHandlersSkipIdleUpdates(t *testing.T) {
	scheme := reconciletest.NewScheme(t, vmetricsv1b1.AddToScheme)
	counter := &reconciletest.UpdateCounter{}
	r := NewVmSingleReconciler(reconciletest.NewClient(scheme, counter), scheme, &fakediscovery.FakeDiscovery{Fake: &ktesting.Fake{}})

	reconciletest.AssertIdleReconcileSkipsUpdates(t, hardeningCR(), counter, map[string]reconciletest.Handler{
		"ServiceAccount":     r.handleServiceAccount,
		"ClusterRole":        r.handleClusterRole,
		"ClusterRoleBinding": r.handleClusterRoleBinding,
		"VMSingle":           r.handleVmSingle,
		"Ingress":            r.handleIngressV1,
	})
}
