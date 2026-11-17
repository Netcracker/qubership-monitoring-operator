package grafana_operator

import (
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	grafv1 "github.com/grafana/grafana-operator/v5/api/v1beta1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	fakediscovery "k8s.io/client-go/discovery/fake"
	k8stesting "k8s.io/client-go/testing"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
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

func getLegacyDashboard(t *testing.T, c client.Client, key client.ObjectKey) *unstructured.Unstructured {
	t.Helper()
	got := &unstructured.Unstructured{}
	got.SetGroupVersionKind(legacyGrafanaDashboardGVK)
	require.NoError(t, c.Get(t.Context(), key, got))
	return got
}

func TestHandleGrafanaDashboardCreatesLegacyDashboard(t *testing.T) {
	cr := testPlatformMonitoring()
	desired, err := grafanaDashboard(cr, testDashboardAsset)
	require.NoError(t, err)

	r := newGrafanaDashboardTestReconciler(t, cr)
	require.NoError(t, r.handleGrafanaDashboard(testDashboardAsset, cr))

	got := getLegacyDashboard(t, r.Client, client.ObjectKeyFromObject(desired))
	assert.Equal(t, desired.Object["spec"], got.Object["spec"])
	assert.True(t, metav1.IsControlledBy(got, cr), "the legacy dashboard must be owned by PlatformMonitoring")
}

func TestHandleGrafanaDashboardSkipsNoOpUpdate(t *testing.T) {
	cr := testPlatformMonitoring()
	desired, err := grafanaDashboard(cr, testDashboardAsset)
	require.NoError(t, err)
	desired.SetResourceVersion("1")

	r := newGrafanaDashboardTestReconciler(t, cr, desired.DeepCopy())
	require.NoError(t, r.handleGrafanaDashboard(testDashboardAsset, cr))

	got := getLegacyDashboard(t, r.Client, client.ObjectKeyFromObject(desired))
	assert.Equal(t, "1", got.GetResourceVersion(), "no-op reconcile must not bump resourceVersion")
}

func TestHandleGrafanaDashboardUpdatesChangedSpec(t *testing.T) {
	cr := testPlatformMonitoring()
	desired, err := grafanaDashboard(cr, testDashboardAsset)
	require.NoError(t, err)

	stale := desired.DeepCopy()
	stale.SetResourceVersion("1")
	require.NoError(t, unstructured.SetNestedField(stale.Object, "{}", "spec", "json"))

	r := newGrafanaDashboardTestReconciler(t, cr, stale)
	require.NoError(t, r.handleGrafanaDashboard(testDashboardAsset, cr))

	got := getLegacyDashboard(t, r.Client, client.ObjectKeyFromObject(desired))
	assert.Equal(t, desired.Object["spec"], got.Object["spec"])
	assert.NotEqual(t, "1", got.GetResourceVersion(), "spec change must trigger Update")
}

func TestHandleGrafanaDashboardUpdatesChangedLabels(t *testing.T) {
	cr := testPlatformMonitoring()
	desired, err := grafanaDashboard(cr, testDashboardAsset)
	require.NoError(t, err)

	stale := desired.DeepCopy()
	stale.SetResourceVersion("1")
	labels := stale.GetLabels()
	labels["test.example.com/stale"] = "true"
	stale.SetLabels(labels)

	r := newGrafanaDashboardTestReconciler(t, cr, stale)
	require.NoError(t, r.handleGrafanaDashboard(testDashboardAsset, cr))

	got := getLegacyDashboard(t, r.Client, client.ObjectKeyFromObject(desired))
	assert.Equal(t, desired.GetLabels(), got.GetLabels())
	assert.NotEqual(t, "1", got.GetResourceVersion(), "label change must trigger Update")
}

func TestHandleGrafanaDashboardReplacesOnlyOwnedV1beta1Dashboard(t *testing.T) {
	foreignOwner := testPlatformMonitoring()
	foreignOwner.Name = "other-platformmonitoring"
	foreignOwner.UID = types.UID("other-platform-monitoring-uid")

	for name, owner := range map[string]*monv1.PlatformMonitoring{
		"controlled by this PlatformMonitoring": testPlatformMonitoring(),
		"controlled by another object":          foreignOwner,
		"without an owner":                      nil,
	} {
		t.Run(name, func(t *testing.T) {
			cr := testPlatformMonitoring()
			desired, err := grafanaDashboard(cr, testDashboardAsset)
			require.NoError(t, err)
			native := &grafv1.GrafanaDashboard{ObjectMeta: metav1.ObjectMeta{Name: desired.GetName(), Namespace: desired.GetNamespace()}}
			r := newGrafanaDashboardTestReconciler(t, cr)
			if owner != nil {
				require.NoError(t, controllerutil.SetControllerReference(owner, native, r.Scheme))
			}
			require.NoError(t, r.Client.Create(t.Context(), native))

			require.NoError(t, r.handleGrafanaDashboard(testDashboardAsset, cr))

			getLegacyDashboard(t, r.Client, client.ObjectKeyFromObject(desired))
			got := &grafv1.GrafanaDashboard{}
			err = r.Client.Get(t.Context(), client.ObjectKeyFromObject(native), got)
			if owner != nil && owner.UID == cr.UID {
				assert.True(t, apierrors.IsNotFound(err), "Get v1beta1 dashboard error = %v", err)
				return
			}
			require.NoError(t, err, "a v1beta1 dashboard that this PlatformMonitoring does not control must stay")
			assert.Equal(t, native.GetOwnerReferences(), got.GetOwnerReferences())
		})
	}
}

func TestDeleteGrafanaDashboardDeletesOnlyBuiltInConvertedCopy(t *testing.T) {
	for name, builtIn := range map[string]bool{"built-in copy": true, "copy without the built-in label": false} {
		t.Run(name, func(t *testing.T) {
			cr := testPlatformMonitoring()
			desired, err := grafanaDashboard(cr, testDashboardAsset)
			require.NoError(t, err)
			labels := map[string]string{converterManagedByLabelKey: converterManagedByLabelValue}
			if builtIn {
				labels[builtInDashboardLabelKey] = builtInDashboardLabelValue
			}
			converted := &grafv1.GrafanaDashboard{ObjectMeta: metav1.ObjectMeta{
				Name:      desired.GetName() + convertedDashboardNameSuffix,
				Namespace: desired.GetNamespace(),
				Labels:    labels,
			}}
			r := newGrafanaDashboardTestReconciler(t, cr, desired.DeepCopy(), converted)

			require.NoError(t, r.deleteGrafanaDashboard(testDashboardAsset, cr))

			err = r.Client.Get(t.Context(), client.ObjectKeyFromObject(converted), &grafv1.GrafanaDashboard{})
			assert.Equal(t, builtIn, apierrors.IsNotFound(err), "Get converted dashboard error = %v", err)
		})
	}
}

func TestDeleteGrafanaDashboardDeletesLegacyDashboard(t *testing.T) {
	cr := testPlatformMonitoring()
	desired, err := grafanaDashboard(cr, testDashboardAsset)
	require.NoError(t, err)

	r := newGrafanaDashboardTestReconciler(t, cr, desired.DeepCopy())
	require.NoError(t, r.deleteGrafanaDashboard(testDashboardAsset, cr))
	require.NoError(t, r.deleteGrafanaDashboard(testDashboardAsset, cr), "deleting an absent dashboard must succeed")

	got := &unstructured.Unstructured{}
	got.SetGroupVersionKind(legacyGrafanaDashboardGVK)
	err = r.Client.Get(t.Context(), client.ObjectKeyFromObject(desired), got)
	assert.True(t, apierrors.IsNotFound(err), "Get legacy dashboard error = %v", err)
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
			Dc:     &fakediscovery.FakeDiscovery{Fake: &k8stesting.Fake{}},
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
	desired, err := grafanaOperatorDeployment(cr, false)
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
	desired, err := grafanaOperatorDeployment(cr, false)
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

func TestHandleDeploymentSkipsAPIDefaultedTokenVolume(t *testing.T) {
	cr := testPlatformMonitoring()
	desired, err := grafanaOperatorDeployment(cr, false)
	require.NoError(t, err)
	live := desired.DeepCopy()
	live.ResourceVersion = "1"
	// The API server stores defaultMode 420 on every projected volume.
	tokenVolume := projectedVolume(t, live)
	tokenVolume.DefaultMode = ptr.To(int32(420))

	r := newGrafanaOperatorWorkloadTestReconciler(t, cr, live)
	require.NoError(t, r.handleDeployment(cr))

	got := &appsv1.Deployment{}
	require.NoError(t, r.Client.Get(t.Context(), client.ObjectKeyFromObject(desired), got))
	assert.Equal(t, "1", got.ResourceVersion, "API-defaulted volume fields must not trigger Update")
}

func TestHandleDeploymentUpdatesChangedTokenVolume(t *testing.T) {
	cr := testPlatformMonitoring()
	desired, err := grafanaOperatorDeployment(cr, false)
	require.NoError(t, err)
	live := desired.DeepCopy()
	live.ResourceVersion = "1"
	projectedVolume(t, live).Sources[0].ServiceAccountToken.Audience = "stale"

	r := newGrafanaOperatorWorkloadTestReconciler(t, cr, live)
	require.NoError(t, r.handleDeployment(cr))

	got := &appsv1.Deployment{}
	require.NoError(t, r.Client.Get(t.Context(), client.ObjectKeyFromObject(desired), got))
	assert.NotEqual(t, "1", got.ResourceVersion, "volume change must trigger Update")
	assert.Equal(t, projectedVolume(t, desired).Sources, projectedVolume(t, got).Sources)
}

func projectedVolume(t *testing.T, deployment *appsv1.Deployment) *corev1.ProjectedVolumeSource {
	t.Helper()
	for i := range deployment.Spec.Template.Spec.Volumes {
		if projected := deployment.Spec.Template.Spec.Volumes[i].Projected; projected != nil {
			return projected
		}
	}
	require.FailNow(t, "the Deployment has no projected volume")
	return nil
}
