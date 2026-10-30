package grafana

import (
	"testing"

	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils/reconciletest"
	grafv1 "github.com/grafana/grafana-operator/v5/api/v1beta1"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	fakediscovery "k8s.io/client-go/discovery/fake"
	kubernetesfake "k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

func TestHandlersSkipIdleUpdates(t *testing.T) {
	scheme := reconciletest.NewScheme(t, promv1.AddToScheme, grafv1.AddToScheme)
	counter := &reconciletest.UpdateCounter{}
	r := &GrafanaReconciler{
		KubeClient: kubernetesfake.NewSimpleClientset(),
		ComponentReconciler: &utils.ComponentReconciler{
			Client: reconciletest.NewClient(scheme, counter),
			Scheme: scheme,
			Dc:     &fakediscovery.FakeDiscovery{Fake: &k8stesting.Fake{}},
			Log:    utils.Logger("grafana_idle_test"),
		},
	}

	reconciletest.AssertIdleReconcileSkipsUpdates(t, grafanaComparisonPlatformMonitoring(nil), counter, map[string]reconciletest.Handler{
		"GrafanaDatasource":       r.handleGrafanaDataSource,
		"PromxyGrafanaDatasource": r.handleGrafanaPromxyDataSource,
		"Ingress":                 r.handleIngressV1,
		"PodMonitor":              r.handlePodMonitor,
	})
}
