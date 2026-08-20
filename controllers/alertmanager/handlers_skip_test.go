package alertmanager

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

func TestHandleAlertmanagerSkipsCRDDefaultedSpec(t *testing.T) {
	cr := skipTestPlatformMonitoring()
	desired, err := alertmanager(cr)
	require.NoError(t, err)
	live := desired.DeepCopy()
	live.ResourceVersion = "1"
	applyCRDDefaultedAlertmanagerSpec(&live.Spec)

	r := newAlertmanagerSkipTestReconciler(t, cr, live)
	require.NoError(t, r.handleAlertmanager(cr))

	got := &promv1.Alertmanager{}
	require.NoError(t, r.Client.Get(t.Context(), client.ObjectKeyFromObject(desired), got))
	assert.Equal(t, "1", got.ResourceVersion, "CRD-defaulted Alertmanager spec must not trigger Update")
	assert.Equal(t, "web", got.Spec.PortName)
	assert.Equal(t, "120h", string(got.Spec.Retention))
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
			AlertManager: &monv1.AlertManager{
				Install: &install,
				Image:   "docker.io/prom/alertmanager:v0.33.1",
			},
		},
	}
}

func newAlertmanagerSkipTestReconciler(t *testing.T, objs ...client.Object) *AlertManagerReconciler {
	t.Helper()
	scheme := runtime.NewScheme()
	require.NoError(t, monv1.AddToScheme(scheme))
	require.NoError(t, promv1.AddToScheme(scheme))
	return NewAlertManagerReconciler(
		fake.NewClientBuilder().WithScheme(scheme).WithObjects(objs...).Build(),
		scheme,
		nil,
	)
}

func applyCRDDefaultedAlertmanagerSpec(spec *promv1.AlertmanagerSpec) {
	if spec.PortName == "" {
		spec.PortName = "web"
	}
	if spec.Retention == "" {
		spec.Retention = "120h"
	}
	if spec.AlertmanagerConfigMatcherStrategy.Type == "" {
		spec.AlertmanagerConfigMatcherStrategy.Type = promv1.OnNamespaceConfigMatcherStrategyType
	}
}
