package prometheus_rules

import (
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func skipTestPlatformMonitoring() *monv1.PlatformMonitoring {
	install := true
	return &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "platformmonitoring",
			Namespace: "monitoring",
			UID:       types.UID("platform-monitoring-uid"),
		},
		Spec: monv1.PlatformMonitoringSpec{
			PrometheusRules: &monv1.PrometheusRules{
				Install:    &install,
				RuleGroups: []string{"SelfMonitoring"},
			},
		},
	}
}

func newPrometheusRulesSkipTestReconciler(t *testing.T, objs ...client.Object) *PrometheusRulesReconciler {
	t.Helper()
	scheme := runtime.NewScheme()
	require.NoError(t, monv1.AddToScheme(scheme))
	require.NoError(t, promv1.AddToScheme(scheme))
	return NewPrometheusRulesReconciler(
		fake.NewClientBuilder().WithScheme(scheme).WithObjects(objs...).Build(),
		scheme,
	)
}

func TestHandlePrometheusRulesSkipsNoOpUpdate(t *testing.T) {
	cr := skipTestPlatformMonitoring()
	desired, err := prometheusRules(cr)
	require.NoError(t, err)
	desired.ResourceVersion = "1"

	r := newPrometheusRulesSkipTestReconciler(t, cr, desired.DeepCopy())
	require.NoError(t, r.handlePrometheusRules(cr))

	got := &promv1.PrometheusRule{}
	require.NoError(t, r.Client.Get(t.Context(), client.ObjectKeyFromObject(desired), got))
	assert.Equal(t, "1", got.ResourceVersion, "no-op reconcile must not bump resourceVersion")
}

func TestHandlePrometheusRulesUpdatesChangedSpec(t *testing.T) {
	cr := skipTestPlatformMonitoring()
	desired, err := prometheusRules(cr)
	require.NoError(t, err)
	live := desired.DeepCopy()
	live.ResourceVersion = "1"
	require.NotEmpty(t, live.Spec.Groups)
	live.Spec.Groups[0].Name = "stale"

	r := newPrometheusRulesSkipTestReconciler(t, cr, live)
	require.NoError(t, r.handlePrometheusRules(cr))

	got := &promv1.PrometheusRule{}
	require.NoError(t, r.Client.Get(t.Context(), client.ObjectKeyFromObject(desired), got))
	assert.NotEqual(t, "1", got.ResourceVersion, "spec change must trigger Update")
	require.NotEmpty(t, got.Spec.Groups)
	assert.Equal(t, desired.Spec.Groups[0].Name, got.Spec.Groups[0].Name)
}
