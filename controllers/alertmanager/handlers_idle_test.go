package alertmanager

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
	r := NewAlertManagerReconciler(reconciletest.NewClient(scheme, counter), scheme, &fakediscovery.FakeDiscovery{Fake: &k8stesting.Fake{}})

	reconciletest.AssertIdleReconcileSkipsUpdates(t, skipTestPlatformMonitoring(), counter, map[string]reconciletest.Handler{
		"ServiceAccount": r.handleServiceAccount,
		"Secret":         r.handleSecret,
		"Alertmanager":   r.handleAlertmanager,
		"Service":        r.handleService,
		"Ingress":        r.handleIngressV1,
		"PodMonitor":     r.handlePodMonitor,
	})
}
