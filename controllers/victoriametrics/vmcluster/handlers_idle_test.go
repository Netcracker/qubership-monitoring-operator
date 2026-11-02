package vmcluster

import (
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils/reconciletest"
	vmetricsv1b1 "github.com/VictoriaMetrics/operator/api/operator/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	fakediscovery "k8s.io/client-go/discovery/fake"
	ktesting "k8s.io/client-go/testing"
)

func TestHandlersSkipIdleUpdates(t *testing.T) {
	scheme := reconciletest.NewScheme(t, vmetricsv1b1.AddToScheme)
	counter := &reconciletest.UpdateCounter{}
	r := NewVmClusterReconciler(reconciletest.NewClient(scheme, counter), scheme, &fakediscovery.FakeDiscovery{Fake: &ktesting.Fake{}})
	cr := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Name: "platformmonitoring", Namespace: "monitoring"},
		Spec: monv1.PlatformMonitoringSpec{
			Victoriametrics: &monv1.Victoriametrics{VmCluster: monv1.VmCluster{}},
		},
	}

	reconciletest.AssertIdleReconcileSkipsUpdates(t, cr, counter, map[string]reconciletest.Handler{
		"ServiceAccount":     r.handleServiceAccount,
		"ClusterRole":        r.handleClusterRole,
		"ClusterRoleBinding": r.handleClusterRoleBinding,
		"VMCluster":          r.handleVmCluster,
		"Ingress":            r.handleIngressV1,
	})
}
