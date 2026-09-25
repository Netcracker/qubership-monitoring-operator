package grafana_operator

import (
	"testing"

	"github.com/Netcracker/qubership-monitoring-operator/controllers/grafana"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	grafv1 "github.com/grafana/grafana-operator/v5/api/v1beta1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
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

func TestReconcileDashboardSelectionDeletesAdoptedDashboardWhenSelectorIsImmutable(t *testing.T) {
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

	r := newGrafanaDashboardTestReconciler(t, dash)
	require.NoError(t, r.updateDashboardSelectionFallback(dash, grafana.DashboardSelectionDetach))

	got := &grafv1.GrafanaDashboard{}
	err := r.Client.Get(t.Context(), client.ObjectKeyFromObject(dash), got)
	assert.Equal(t, true, apierrors.IsNotFound(err), "dashboard must be removed after immutable deselection")
}

func TestUpdateDashboardSelectionFallbackPreservesDashboardForImmutableAdoption(t *testing.T) {
	dash := &grafv1.GrafanaDashboard{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "app-overview",
			Namespace: "apps",
		},
		Spec: grafv1.GrafanaDashboardSpec{
			GrafanaCommonSpec: grafv1.GrafanaCommonSpec{
				InstanceSelector: &metav1.LabelSelector{
					MatchLabels: map[string]string{"app.kubernetes.io/component": "other-grafana"},
				},
			},
		},
	}
	r := newGrafanaDashboardTestReconciler(t, dash)

	require.NoError(t, r.updateDashboardSelectionFallback(dash, grafana.DashboardSelectionLeave))
	require.NoError(t, r.updateDashboardSelectionFallback(dash, grafana.DashboardSelectionAdopt))
	got := &grafv1.GrafanaDashboard{}
	require.NoError(t, r.Client.Get(t.Context(), client.ObjectKeyFromObject(dash), got))
	assert.Equal(t, dash.Spec, got.Spec)
	assert.Equal(t, dash.Annotations, got.Annotations)
}
