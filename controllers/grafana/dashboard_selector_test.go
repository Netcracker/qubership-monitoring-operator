package grafana

import (
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestDashboardLabelsMatch(t *testing.T) {
	t.Parallel()
	labels := map[string]string{"app": "payments"}

	match, err := DashboardLabelsMatch(nil, labels)
	require.NoError(t, err)
	assert.Equal(t, false, match, "DashboardLabelsMatch(nil)")

	match, err = DashboardLabelsMatch([]*metav1.LabelSelector{{}}, labels)
	require.NoError(t, err)
	assert.Equal(t, true, match, "DashboardLabelsMatch([{}])")

	match, err = DashboardLabelsMatch([]*metav1.LabelSelector{
		{MatchLabels: map[string]string{"app": "other"}},
		{MatchLabels: map[string]string{"app": "payments"}},
	}, labels)
	require.NoError(t, err)
	assert.Equal(t, true, match, "DashboardLabelsMatch(any entry)")

	match, err = DashboardLabelsMatch([]*metav1.LabelSelector{
		{MatchLabels: map[string]string{"app": "other"}},
	}, labels)
	require.NoError(t, err)
	assert.Equal(t, false, match, "DashboardLabelsMatch(no entry)")
}

func TestNamespaceSelectorMatches(t *testing.T) {
	t.Parallel()
	labels := map[string]string{"openshift.io/cluster-monitoring": "true"}

	match, err := NamespaceSelectorMatches(nil, labels)
	require.NoError(t, err)
	assert.Equal(t, false, match, "NamespaceSelectorMatches(nil)")

	match, err = NamespaceSelectorMatches(&metav1.LabelSelector{}, labels)
	require.NoError(t, err)
	assert.Equal(t, true, match, "NamespaceSelectorMatches({})")

	match, err = NamespaceSelectorMatches(&metav1.LabelSelector{
		MatchExpressions: []metav1.LabelSelectorRequirement{{
			Key:      "openshift.io/cluster-monitoring",
			Operator: metav1.LabelSelectorOpNotIn,
			Values:   []string{"true"},
		}},
	}, labels)
	require.NoError(t, err)
	assert.Equal(t, false, match, "NamespaceSelectorMatches(NotIn)")
}

func TestSelectDashboard(t *testing.T) {
	t.Parallel()
	all := []*metav1.LabelSelector{{}}
	none := []*metav1.LabelSelector{{MatchLabels: map[string]string{"app": "missing"}}}

	action, err := SelectDashboard(all, nil, map[string]string{"app": "payments"}, nil, false, false)
	require.NoError(t, err)
	assert.Equal(t, DashboardSelectionAdopt, action, "SelectDashboard(match, not adopted)")

	action, err = SelectDashboard(none, nil, map[string]string{"app": "payments"}, nil, true, false)
	require.NoError(t, err)
	assert.Equal(t, DashboardSelectionDetach, action, "SelectDashboard(no match, adopted)")

	action, err = SelectDashboard(none, nil, map[string]string{"app": "payments"}, nil, false, false)
	require.NoError(t, err)
	assert.Equal(t, DashboardSelectionLeave, action, "SelectDashboard(no match, not adopted)")

	action, err = SelectDashboard(all, nil, map[string]string{"app": "payments"}, nil, false, true)
	require.NoError(t, err)
	assert.Equal(t, DashboardSelectionLeave, action, "SelectDashboard(nil namespace selector)")

	namespaceMismatch := &metav1.LabelSelector{MatchLabels: map[string]string{"team": "other"}}
	namespaceLabels := map[string]string{"team": "platform"}
	action, err = SelectDashboard(all, namespaceMismatch, map[string]string{"app": "payments"}, namespaceLabels, true, true)
	require.NoError(t, err)
	assert.Equal(t, DashboardSelectionDetach, action, "SelectDashboard(namespace mismatch, adopted)")

	action, err = SelectDashboard(all, namespaceMismatch, map[string]string{"app": "payments"}, namespaceLabels, false, true)
	require.NoError(t, err)
	assert.Equal(t, DashboardSelectionLeave, action, "SelectDashboard(namespace mismatch, not adopted)")
}

func TestGrafanaNamespace(t *testing.T) {
	t.Parallel()
	cr := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"},
		Spec:       monv1.PlatformMonitoringSpec{Grafana: &monv1.Grafana{Namespace: "grafana"}},
	}
	assert.Equal(t, "grafana", GrafanaNamespace(cr))
	cr.Spec.Grafana.Namespace = ""
	assert.Equal(t, "monitoring", GrafanaNamespace(cr))
}

func TestSelectDashboardRejectsInvalidSelector(t *testing.T) {
	t.Parallel()
	bad := []*metav1.LabelSelector{{
		MatchExpressions: []metav1.LabelSelectorRequirement{{
			Key:      "app",
			Operator: "NoSuch",
		}},
	}}
	_, err := SelectDashboard(bad, nil, map[string]string{"app": "payments"}, nil, false, false)
	require.Error(t, err)
	_, err = DashboardLabelsMatch(bad, map[string]string{"app": "payments"})
	require.Error(t, err)
}

func TestGrafanaDashboardInstanceLabels(t *testing.T) {
	t.Parallel()
	cr := &monv1.PlatformMonitoring{
		Spec: monv1.PlatformMonitoringSpec{
			Grafana: &monv1.Grafana{
				Labels: map[string]string{
					"app.kubernetes.io/component": "custom-grafana",
					"app.kubernetes.io/part-of":   "platform",
					"team":                        "observability",
				},
			},
		},
	}

	assert.Equal(t, map[string]string{
		"app.kubernetes.io/component": "custom-grafana",
		"app.kubernetes.io/part-of":   "platform",
	}, GrafanaDashboardInstanceLabels(cr))
}
