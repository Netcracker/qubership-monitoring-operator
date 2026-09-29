package prometheus

import (
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestHandlePrometheusSkipsCRDDefaultedSpec(t *testing.T) {
	cr := skipTestPlatformMonitoring()
	desired, err := prometheus(cr)
	require.NoError(t, err)
	live := desired.DeepCopy()
	live.ResourceVersion = "1"
	applyCRDDefaultedPrometheusSpec(&live.Spec)

	r := newPrometheusSkipTestReconciler(t, cr, live)
	require.NoError(t, r.handlePrometheus(cr))

	got := &promv1.Prometheus{}
	require.NoError(t, r.Client.Get(t.Context(), client.ObjectKeyFromObject(desired), got))
	assert.Equal(t, "1", got.ResourceVersion, "CRD-defaulted Prometheus spec must not trigger Update")
	assert.Equal(t, promv1.Duration("30s"), got.Spec.ScrapeInterval)
	assert.Equal(t, "web", got.Spec.PortName)
}

func TestHandlePrometheusUpdatesChangedSpec(t *testing.T) {
	cr := skipTestPlatformMonitoring()
	desired, err := prometheus(cr)
	require.NoError(t, err)
	live := desired.DeepCopy()
	live.ResourceVersion = "1"
	live.Spec.ScrapeInterval = "15s"

	r := newPrometheusSkipTestReconciler(t, cr, live)
	require.NoError(t, r.handlePrometheus(cr))

	got := &promv1.Prometheus{}
	require.NoError(t, r.Client.Get(t.Context(), client.ObjectKeyFromObject(desired), got))
	assert.NotEqual(t, "1", got.ResourceVersion, "spec change must trigger Update")
	assert.Equal(t, promv1.Duration("30s"), got.Spec.ScrapeInterval)
}

func skipTestPlatformMonitoring() *monv1.PlatformMonitoring {
	install := true
	return &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "platformmonitoring",
			Namespace: "monitoring",
			UID:       types.UID("platform-monitoring-uid"),
		},
		Spec: monv1.PlatformMonitoringSpec{
			Prometheus: &monv1.Prometheus{
				Install:             &install,
				Image:               "docker.io/prom/prometheus:v3.13.1",
				ConfigReloaderImage: "quay.io/prometheus-operator/prometheus-config-reloader:v0.93.0",
				Operator: monv1.PrometheusOperator{
					Image: "quay.io/prometheus-operator/prometheus-operator:v0.93.0",
				},
			},
			AlertManager: &monv1.AlertManager{
				Install: &install,
				Image:   "docker.io/prom/alertmanager:v0.33.1",
			},
		},
	}
}

func newPrometheusSkipTestReconciler(t *testing.T, objs ...client.Object) *PrometheusReconciler {
	t.Helper()
	scheme := runtime.NewScheme()
	require.NoError(t, monv1.AddToScheme(scheme))
	require.NoError(t, promv1.AddToScheme(scheme))
	return &PrometheusReconciler{
		ComponentReconciler: &utils.ComponentReconciler{
			Client: fake.NewClientBuilder().WithScheme(scheme).WithObjects(objs...).Build(),
			Scheme: scheme,
			Log:    utils.Logger("prometheus_skip_test"),
		},
	}
}

func applyCRDDefaultedPrometheusSpec(spec *promv1.PrometheusSpec) {
	if spec.ScrapeInterval == "" {
		spec.ScrapeInterval = "30s"
	}
	if spec.EvaluationInterval == "" {
		spec.EvaluationInterval = "30s"
	}
	if spec.PortName == "" {
		spec.PortName = "web"
	}
}
