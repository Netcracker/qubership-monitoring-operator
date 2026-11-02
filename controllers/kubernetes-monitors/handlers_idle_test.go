package kubernetes_monitors

import (
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils/reconciletest"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	fakediscovery "k8s.io/client-go/discovery/fake"
	k8stesting "k8s.io/client-go/testing"
)

func TestHandlersSkipIdleUpdates(t *testing.T) {
	scheme := reconciletest.NewScheme(t, promv1.AddToScheme)
	counter := &reconciletest.UpdateCounter{}
	r := NewKubernetesMonitorsReconciler(reconciletest.NewClient(scheme, counter), scheme, &fakediscovery.FakeDiscovery{Fake: &k8stesting.Fake{}})
	cr := &monv1.PlatformMonitoring{ObjectMeta: metav1.ObjectMeta{Name: "platformmonitoring", Namespace: "monitoring"}}

	reconciletest.AssertIdleReconcileSkipsUpdates(t, cr, counter, map[string]reconciletest.Handler{
		"ApiServerServiceMonitor": r.handleApiServerServiceMonitor,
		"KubeletServiceMonitor":   r.handleKubeletServiceMonitor,
		"NginxIngressPodMonitor":  r.handleNginxIngressPodMonitor,
		"CoreDnsServiceMonitor": func(cr *monv1.PlatformMonitoring) error {
			return r.handleCoreDnsServiceMonitor(cr, false)
		},
		"OpenshiftServiceMonitor": func(cr *monv1.PlatformMonitoring) error {
			return r.handleOpenshiftServiceMonitor(cr, utils.OpenshiftApiServerServiceMonitorAsset,
				utils.OpenshiftApiServerServiceMonitorName, utils.OpenshiftApiserver)
		},
	})
}
