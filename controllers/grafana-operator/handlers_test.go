package grafana_operator

import (
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	grafv1 "github.com/grafana/grafana-operator/v5/api/v1beta1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

const testDashboardAsset = "home-dashboard.yaml"

func testPlatformMonitoring() *monv1.PlatformMonitoring {
	return &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "platformmonitoring",
			Namespace: "monitoring",
			UID:       types.UID("platform-monitoring-uid"),
		},
		Spec: monv1.PlatformMonitoringSpec{
			Grafana: &monv1.Grafana{},
		},
	}
}

func newGrafanaDashboardTestReconciler(t *testing.T, objs ...client.Object) *GrafanaOperatorReconciler {
	t.Helper()
	scheme := runtime.NewScheme()
	require.NoError(t, monv1.AddToScheme(scheme))
	require.NoError(t, grafv1.AddToScheme(scheme))
	return &GrafanaOperatorReconciler{
		ComponentReconciler: &utils.ComponentReconciler{
			Client: fake.NewClientBuilder().WithScheme(scheme).WithObjects(objs...).Build(),
			Scheme: scheme,
			Log:    utils.Logger("grafanaoperator_dashboard_test"),
		},
	}
}

func TestHandleGrafanaDashboardCreate(t *testing.T) {
	cr := testPlatformMonitoring()
	desired, err := grafanaDashboard(cr, testDashboardAsset)
	require.NoError(t, err)

	r := newGrafanaDashboardTestReconciler(t, cr)
	require.NoError(t, r.handleGrafanaDashboard(testDashboardAsset, cr))

	got := &grafv1.GrafanaDashboard{}
	require.NoError(t, r.Client.Get(t.Context(), client.ObjectKeyFromObject(desired), got))
	assert.Equal(t, desired.Name, got.Name)
	assert.Equal(t, desired.Spec.ResyncPeriod, got.Spec.ResyncPeriod)
}

func TestHandleGrafanaDashboardSkipsNoOpUpdate(t *testing.T) {
	cr := testPlatformMonitoring()
	desired, err := grafanaDashboard(cr, testDashboardAsset)
	require.NoError(t, err)
	desired.ResourceVersion = "1"

	r := newGrafanaDashboardTestReconciler(t, cr, desired.DeepCopy())
	require.NoError(t, r.handleGrafanaDashboard(testDashboardAsset, cr))

	got := &grafv1.GrafanaDashboard{}
	require.NoError(t, r.Client.Get(t.Context(), client.ObjectKeyFromObject(desired), got))
	assert.Equal(t, "1", got.ResourceVersion, "no-op reconcile must not bump resourceVersion")
}

func TestHandleGrafanaDashboardUpdatesChangedSpec(t *testing.T) {
	cr := testPlatformMonitoring()
	desired, err := grafanaDashboard(cr, testDashboardAsset)
	require.NoError(t, err)

	stale := desired.DeepCopy()
	stale.ResourceVersion = "1"
	stale.Spec.ResyncPeriod = metav1.Duration{Duration: desired.Spec.ResyncPeriod.Duration * 2}

	r := newGrafanaDashboardTestReconciler(t, cr, stale)
	require.NoError(t, r.handleGrafanaDashboard(testDashboardAsset, cr))

	got := &grafv1.GrafanaDashboard{}
	require.NoError(t, r.Client.Get(t.Context(), client.ObjectKeyFromObject(desired), got))
	assert.Equal(t, desired.Spec.ResyncPeriod, got.Spec.ResyncPeriod)
	assert.NotEqual(t, "1", got.ResourceVersion, "spec change must trigger Update")
}

func TestHandleGrafanaDashboardUpdatesChangedLabels(t *testing.T) {
	cr := testPlatformMonitoring()
	desired, err := grafanaDashboard(cr, testDashboardAsset)
	require.NoError(t, err)

	stale := desired.DeepCopy()
	stale.ResourceVersion = "1"
	labels := stale.GetLabels()
	if labels == nil {
		labels = map[string]string{}
	}
	labels["test.example.com/stale"] = "true"
	stale.SetLabels(labels)

	r := newGrafanaDashboardTestReconciler(t, cr, stale)
	require.NoError(t, r.handleGrafanaDashboard(testDashboardAsset, cr))

	got := &grafv1.GrafanaDashboard{}
	require.NoError(t, r.Client.Get(t.Context(), client.ObjectKeyFromObject(desired), got))
	assert.Equal(t, desired.GetLabels(), got.GetLabels())
	assert.NotEqual(t, "1", got.ResourceVersion, "label change must trigger Update")
}

func newGrafanaOperatorWorkloadTestReconciler(t *testing.T, objs ...client.Object) *GrafanaOperatorReconciler {
	t.Helper()
	scheme := runtime.NewScheme()
	require.NoError(t, monv1.AddToScheme(scheme))
	require.NoError(t, grafv1.AddToScheme(scheme))
	require.NoError(t, appsv1.AddToScheme(scheme))
	require.NoError(t, rbacv1.AddToScheme(scheme))
	return &GrafanaOperatorReconciler{
		ComponentReconciler: &utils.ComponentReconciler{
			Client: fake.NewClientBuilder().WithScheme(scheme).WithObjects(objs...).Build(),
			Scheme: scheme,
			Log:    utils.Logger("grafanaoperator_skip_test"),
		},
	}
}

func TestHandleClusterRoleSkipsNoOpUpdate(t *testing.T) {
	cr := testPlatformMonitoring()
	desired, err := grafanaOperatorClusterRole(cr)
	require.NoError(t, err)
	desired.ResourceVersion = "1"

	r := newGrafanaOperatorWorkloadTestReconciler(t, cr, desired.DeepCopy())
	require.NoError(t, r.handleClusterRole(cr))

	got := &rbacv1.ClusterRole{}
	require.NoError(t, r.Client.Get(t.Context(), client.ObjectKeyFromObject(desired), got))
	assert.Equal(t, "1", got.ResourceVersion, "no-op reconcile must not bump resourceVersion")
}

func TestHandleDeploymentSkipsNoOpUpdateWithPreservedSelector(t *testing.T) {
	cr := testPlatformMonitoring()
	desired, err := grafanaOperatorDeployment(cr)
	require.NoError(t, err)
	desired.ResourceVersion = "1"

	r := newGrafanaOperatorWorkloadTestReconciler(t, cr, desired.DeepCopy())
	require.NoError(t, r.handleDeployment(cr))

	got := &appsv1.Deployment{}
	require.NoError(t, r.Client.Get(t.Context(), client.ObjectKeyFromObject(desired), got))
	assert.Equal(t, "1", got.ResourceVersion, "no-op reconcile must not bump resourceVersion")
	assert.Equal(t, desired.Spec.Selector, got.Spec.Selector)
}

func TestHandleDeploymentUpdatesChangedImage(t *testing.T) {
	cr := testPlatformMonitoring()
	desired, err := grafanaOperatorDeployment(cr)
	require.NoError(t, err)
	live := desired.DeepCopy()
	live.ResourceVersion = "1"
	require.NotEmpty(t, live.Spec.Template.Spec.Containers)
	live.Spec.Template.Spec.Containers[0].Image = "stale:old"

	r := newGrafanaOperatorWorkloadTestReconciler(t, cr, live)
	require.NoError(t, r.handleDeployment(cr))

	got := &appsv1.Deployment{}
	require.NoError(t, r.Client.Get(t.Context(), client.ObjectKeyFromObject(desired), got))
	assert.NotEqual(t, "1", got.ResourceVersion, "image change must trigger Update")
	require.NotEmpty(t, got.Spec.Template.Spec.Containers)
	assert.Equal(t, desired.Spec.Template.Spec.Containers[0].Image, got.Spec.Template.Spec.Containers[0].Image)
}
