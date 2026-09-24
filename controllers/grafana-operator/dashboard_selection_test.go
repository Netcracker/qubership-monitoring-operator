package grafana_operator

import (
	"context"
	"testing"

	"github.com/Netcracker/qubership-monitoring-operator/controllers/grafana"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	grafv1 "github.com/grafana/grafana-operator/v5/api/v1beta1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
)

func TestReconcileDashboardSelectionAdoptsMatchingDashboard(t *testing.T) {
	cr := testPlatformMonitoring()
	cr.Spec.Grafana.DashboardLabelSelector = []*metav1.LabelSelector{{}}
	dash := &grafv1.GrafanaDashboard{
		ObjectMeta: metav1.ObjectMeta{Name: "app-overview", Namespace: cr.Namespace},
	}

	r := newGrafanaDashboardTestReconciler(t, cr, dash)
	require.NoError(t, r.reconcileDashboardSelection(cr))

	got := &grafv1.GrafanaDashboard{}
	require.NoError(t, r.Client.Get(t.Context(), client.ObjectKeyFromObject(dash), got))
	require.NotNil(t, got.Spec.InstanceSelector)
	assert.Equal(t, grafana.GrafanaDashboardInstanceLabels(cr), got.Spec.InstanceSelector.MatchLabels)
	assert.Equal(t, cr.Name, got.Annotations[grafana.DashboardSelectorAnnotation])
	assert.Equal(t, false, got.Spec.AllowCrossNamespaceImport)
}

func TestReconcileDashboardSelectionIgnoresNamespaceSelectorWithoutPrivilegedRights(t *testing.T) {
	previous := utils.PrivilegedRights
	utils.PrivilegedRights = false
	t.Cleanup(func() { utils.PrivilegedRights = previous })

	cr := testPlatformMonitoring()
	cr.Spec.Grafana.DashboardLabelSelector = []*metav1.LabelSelector{{}}
	cr.Spec.Grafana.DashboardNamespaceSelector = &metav1.LabelSelector{
		MatchLabels: map[string]string{"openshift.io/cluster-monitoring": "true"},
	}
	dash := &grafv1.GrafanaDashboard{
		ObjectMeta: metav1.ObjectMeta{Name: "app-overview", Namespace: cr.Namespace},
	}

	r := newGrafanaDashboardTestReconciler(t, cr, dash)
	require.NoError(t, r.reconcileDashboardSelection(cr))

	got := &grafv1.GrafanaDashboard{}
	require.NoError(t, r.Client.Get(t.Context(), client.ObjectKeyFromObject(dash), got))
	require.NotNil(t, got.Spec.InstanceSelector)
	assert.Equal(t, cr.Name, got.Annotations[grafana.DashboardSelectorAnnotation])
}

func TestReconcileDashboardSelectionDetachesAdoptedDashboard(t *testing.T) {
	cr := testPlatformMonitoring()
	cr.Spec.Grafana.DashboardLabelSelector = []*metav1.LabelSelector{{
		MatchLabels: map[string]string{"dashboards": "selected"},
	}}
	dash := &grafv1.GrafanaDashboard{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "app-overview",
			Namespace: cr.Namespace,
			Annotations: map[string]string{
				grafana.DashboardSelectorAnnotation: cr.Name,
			},
		},
		Spec: grafv1.GrafanaDashboardSpec{
			GrafanaCommonSpec: grafv1.GrafanaCommonSpec{
				InstanceSelector: &metav1.LabelSelector{
					MatchLabels: grafana.GrafanaDashboardInstanceLabels(cr),
				},
			},
		},
	}

	r := newGrafanaDashboardTestReconciler(t, cr, dash)
	require.NoError(t, r.reconcileDashboardSelection(cr))

	got := &grafv1.GrafanaDashboard{}
	require.NoError(t, r.Client.Get(t.Context(), client.ObjectKeyFromObject(dash), got))
	assert.Nil(t, got.Spec.InstanceSelector)
	assert.Empty(t, got.Annotations[grafana.DashboardSelectorAnnotation])
}

func TestReconcileDashboardSelectionLeavesUnadoptedDashboard(t *testing.T) {
	cr := testPlatformMonitoring()
	cr.Spec.Grafana.DashboardLabelSelector = []*metav1.LabelSelector{{
		MatchLabels: map[string]string{"dashboards": "selected"},
	}}
	dash := &grafv1.GrafanaDashboard{
		ObjectMeta: metav1.ObjectMeta{
			Name:            "foreign",
			Namespace:       cr.Namespace,
			ResourceVersion: "7",
		},
		Spec: grafv1.GrafanaDashboardSpec{
			GrafanaCommonSpec: grafv1.GrafanaCommonSpec{
				InstanceSelector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"app.kubernetes.io/name": "other"},
				},
			},
		},
	}

	r := newGrafanaDashboardTestReconciler(t, cr, dash)
	require.NoError(t, r.reconcileDashboardSelection(cr))

	got := &grafv1.GrafanaDashboard{}
	require.NoError(t, r.Client.Get(t.Context(), client.ObjectKeyFromObject(dash), got))
	assert.Equal(t, "7", got.ResourceVersion)
	assert.Equal(t, dash.Spec.InstanceSelector, got.Spec.InstanceSelector)
	assert.Empty(t, got.Annotations)
}

func TestBundledDashboardKeptByEmptySelector(t *testing.T) {
	cr := testPlatformMonitoring()
	cr.Spec.Grafana.DashboardLabelSelector = []*metav1.LabelSelector{{}}

	r := newGrafanaDashboardTestReconciler(t, cr)
	selected, err := r.bundledDashboardSelected(cr, testDashboardAsset)
	require.NoError(t, err)
	assert.Equal(t, true, selected, "bundledDashboardSelected([{}])")
}

func TestBundledDashboardSkippedWhenLabelsDoNotMatch(t *testing.T) {
	cr := testPlatformMonitoring()
	cr.Spec.Grafana.DashboardLabelSelector = []*metav1.LabelSelector{{
		MatchLabels: map[string]string{"dashboards": "selected"},
	}}

	r := newGrafanaDashboardTestReconciler(t, cr)
	selected, err := r.bundledDashboardSelected(cr, testDashboardAsset)
	require.NoError(t, err)
	assert.Equal(t, false, selected, "bundledDashboardSelected(no match)")
}

func TestReconcileDashboardSelectionAllowsCrossNamespaceImport(t *testing.T) {
	previous := utils.PrivilegedRights
	utils.PrivilegedRights = true
	t.Cleanup(func() { utils.PrivilegedRights = previous })

	cr := testPlatformMonitoring()
	cr.Spec.Grafana.Namespace = "monitoring"
	cr.Spec.Grafana.DashboardLabelSelector = []*metav1.LabelSelector{{}}
	cr.Spec.Grafana.DashboardNamespaceSelector = &metav1.LabelSelector{}
	dash := &grafv1.GrafanaDashboard{
		ObjectMeta: metav1.ObjectMeta{Name: "app-overview", Namespace: "apps"},
	}
	namespaces := []client.Object{
		&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "monitoring"}},
		&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "apps"}},
	}

	r := newGrafanaDashboardTestReconciler(t, append(namespaces, cr, dash)...)
	require.NoError(t, r.reconcileDashboardSelection(cr))

	got := &grafv1.GrafanaDashboard{}
	require.NoError(t, r.Client.Get(t.Context(), client.ObjectKeyFromObject(dash), got))
	assert.Equal(t, true, got.Spec.AllowCrossNamespaceImport)
	assert.Equal(t, cr.Name, got.Annotations[grafana.DashboardSelectorAnnotation])
}

func TestBundledDashboardSelectedHonorsNamespaceSelectorWhenPrivileged(t *testing.T) {
	previous := utils.PrivilegedRights
	utils.PrivilegedRights = true
	t.Cleanup(func() { utils.PrivilegedRights = previous })

	cr := testPlatformMonitoring()
	cr.Spec.Grafana.DashboardLabelSelector = []*metav1.LabelSelector{{}}
	cr.Spec.Grafana.DashboardNamespaceSelector = &metav1.LabelSelector{
		MatchLabels: map[string]string{"openshift.io/cluster-monitoring": "true"},
	}
	ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: cr.Namespace}}

	r := newGrafanaDashboardTestReconciler(t, cr, ns)
	selected, err := r.bundledDashboardSelected(cr, testDashboardAsset)
	require.NoError(t, err)
	assert.Equal(t, false, selected, "bundledDashboardSelected(namespace excluded)")

	cr.Spec.Grafana.DashboardNamespaceSelector = &metav1.LabelSelector{}
	selected, err = r.bundledDashboardSelected(cr, testDashboardAsset)
	require.NoError(t, err)
	assert.Equal(t, true, selected, "bundledDashboardSelected(empty namespace selector)")
}

func TestBundledDashboardSelectedReportsMissingNamespace(t *testing.T) {
	previous := utils.PrivilegedRights
	utils.PrivilegedRights = true
	t.Cleanup(func() { utils.PrivilegedRights = previous })

	cr := testPlatformMonitoring()
	cr.Spec.Grafana.DashboardLabelSelector = []*metav1.LabelSelector{{}}
	r := newGrafanaDashboardTestReconciler(t, cr)
	_, err := r.bundledDashboardSelected(cr, testDashboardAsset)
	require.Error(t, err)
}

func TestReconcileDashboardSelectionDetachesAdoptedDashboardInMissingNamespace(t *testing.T) {
	previous := utils.PrivilegedRights
	utils.PrivilegedRights = true
	t.Cleanup(func() { utils.PrivilegedRights = previous })

	cr := testPlatformMonitoring()
	cr.Spec.Grafana.DashboardLabelSelector = []*metav1.LabelSelector{{}}
	cr.Spec.Grafana.DashboardNamespaceSelector = &metav1.LabelSelector{}
	dash := &grafv1.GrafanaDashboard{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "app-overview",
			Namespace: "deleted",
			Annotations: map[string]string{
				grafana.DashboardSelectorAnnotation: cr.Name,
			},
		},
		Spec: grafv1.GrafanaDashboardSpec{
			GrafanaCommonSpec: grafv1.GrafanaCommonSpec{
				InstanceSelector: &metav1.LabelSelector{MatchLabels: map[string]string{"app.kubernetes.io/component": "grafana"}},
			},
		},
	}

	r := newGrafanaDashboardTestReconciler(t, cr, dash)
	require.NoError(t, r.reconcileDashboardSelection(cr))

	got := &grafv1.GrafanaDashboard{}
	require.NoError(t, r.Client.Get(t.Context(), client.ObjectKeyFromObject(dash), got))
	assert.Nil(t, got.Spec.InstanceSelector)
}

func TestReconcileDashboardSelectionSuspendsWhenSelectorIsImmutable(t *testing.T) {
	previous := utils.PrivilegedRights
	utils.PrivilegedRights = false
	t.Cleanup(func() { utils.PrivilegedRights = previous })

	cr := testPlatformMonitoring()
	cr.Spec.Grafana.DashboardLabelSelector = []*metav1.LabelSelector{{
		MatchLabels: map[string]string{"dashboards": "selected"},
	}}
	dash := &grafv1.GrafanaDashboard{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "app-overview",
			Namespace: cr.Namespace,
			Annotations: map[string]string{
				grafana.DashboardSelectorAnnotation: cr.Name,
			},
		},
		Spec: grafv1.GrafanaDashboardSpec{
			GrafanaCommonSpec: grafv1.GrafanaCommonSpec{
				InstanceSelector: &metav1.LabelSelector{
					MatchLabels: grafana.GrafanaDashboardInstanceLabels(cr),
				},
			},
		},
	}

	scheme := runtime.NewScheme()
	require.NoError(t, grafv1.AddToScheme(scheme))
	require.NoError(t, corev1.AddToScheme(scheme))
	updates := 0
	base := fake.NewClientBuilder().WithScheme(scheme).WithObjects(dash).Build()
	r := &GrafanaOperatorReconciler{
		ComponentReconciler: &utils.ComponentReconciler{
			Client: interceptor.NewClient(base, interceptor.Funcs{
				Update: func(ctx context.Context, c client.WithWatch, obj client.Object, opts ...client.UpdateOption) error {
					updates++
					if updates == 1 {
						return apierrors.NewInvalid(schema.GroupKind{Group: "grafana.integreatly.org", Kind: "GrafanaDashboard"}, obj.GetName(), nil)
					}
					return c.Update(ctx, obj, opts...)
				},
			}),
			Scheme: scheme,
			Log:    utils.Logger("grafanaoperator_dashboard_test"),
		},
	}
	require.NoError(t, r.reconcileDashboardSelection(cr))

	got := &grafv1.GrafanaDashboard{}
	require.NoError(t, r.Client.Get(t.Context(), client.ObjectKeyFromObject(dash), got))
	assert.Equal(t, true, got.Spec.Suspend)
	assert.Equal(t, cr.Name, got.Annotations[grafana.DashboardSelectorAnnotation])
	require.NotNil(t, got.Spec.InstanceSelector)
}

func TestUpdateDashboardSelectionFallback(t *testing.T) {
	cr := testPlatformMonitoring()
	dash := &grafv1.GrafanaDashboard{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "app-overview",
			Namespace: "apps",
			Annotations: map[string]string{
				grafana.DashboardSelectorAnnotation: cr.Name,
			},
		},
		Spec: grafv1.GrafanaDashboardSpec{
			GrafanaCommonSpec: grafv1.GrafanaCommonSpec{
				Suspend: true,
				InstanceSelector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"app.kubernetes.io/component": "grafana"},
				},
			},
		},
	}
	r := newGrafanaDashboardTestReconciler(t, dash)

	require.NoError(t, r.updateDashboardSelectionFallback(dash, cr, grafana.DashboardSelectionLeave))
	require.NoError(t, r.updateDashboardSelectionFallback(dash, cr, grafana.DashboardSelectionDetach))

	require.NoError(t, r.updateDashboardSelectionFallback(dash, cr, grafana.DashboardSelectionAdopt))
	got := &grafv1.GrafanaDashboard{}
	require.NoError(t, r.Client.Get(t.Context(), client.ObjectKeyFromObject(dash), got))
	assert.Equal(t, false, got.Spec.Suspend)
	assert.Equal(t, true, got.Spec.AllowCrossNamespaceImport)
}
