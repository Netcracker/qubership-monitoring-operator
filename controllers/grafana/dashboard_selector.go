package grafana

import (
	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

const (
	labelComponent = "app.kubernetes.io/component"
	labelPartOf    = "app.kubernetes.io/part-of"

	// DashboardSelectorAnnotation marks a GrafanaDashboard whose instanceSelector this operator set.
	// A later reconcile removes that pointer only when the annotation is present.
	DashboardSelectorAnnotation = "monitoring.netcracker.com/dashboard-selector"
)

// DashboardSelectionAction is what to do with one GrafanaDashboard during selection.
type DashboardSelectionAction int

const (
	// DashboardSelectionLeave leaves the dashboard unchanged.
	DashboardSelectionLeave DashboardSelectionAction = iota
	// DashboardSelectionAdopt points the dashboard at this Grafana.
	DashboardSelectionAdopt
	// DashboardSelectionDetach removes a pointer this operator previously set.
	DashboardSelectionDetach
)

// DashboardLabelsMatch reports whether lbls satisfy any selector.
// A nil list matches nothing. An empty selector matches every label set.
// The dashboard matches when any entry matches.
func DashboardLabelsMatch(selectors []*metav1.LabelSelector, lbls map[string]string) (bool, error) {
	if selectors == nil {
		return false, nil
	}
	for _, selector := range selectors {
		ok, err := labelSelectorMatches(selector, lbls)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}
	return false, nil
}

// NamespaceSelectorMatches reports whether namespaceLabels satisfy selector.
// A nil selector matches no namespace. An empty selector matches every namespace.
func NamespaceSelectorMatches(selector *metav1.LabelSelector, namespaceLabels map[string]string) (bool, error) {
	if selector == nil {
		return false, nil
	}
	return labelSelectorMatches(selector, namespaceLabels)
}

// SelectDashboard decides how a dashboard in a candidate namespace should point at this Grafana.
// adopted is true when the dashboard carries DashboardSelectorAnnotation.
// When honorNamespaceSelector is false, namespaceSelector is ignored.
func SelectDashboard(labelSelectors []*metav1.LabelSelector, namespaceSelector *metav1.LabelSelector, dashboardLabels, namespaceLabels map[string]string, adopted, honorNamespaceSelector bool) (DashboardSelectionAction, error) {
	if honorNamespaceSelector {
		namespaceOK, err := NamespaceSelectorMatches(namespaceSelector, namespaceLabels)
		if err != nil {
			return DashboardSelectionLeave, err
		}
		if !namespaceOK {
			if adopted {
				return DashboardSelectionDetach, nil
			}
			return DashboardSelectionLeave, nil
		}
	}
	labelsOK, err := DashboardLabelsMatch(labelSelectors, dashboardLabels)
	if err != nil {
		return DashboardSelectionLeave, err
	}
	if labelsOK {
		return DashboardSelectionAdopt, nil
	}
	if adopted {
		return DashboardSelectionDetach, nil
	}
	return DashboardSelectionLeave, nil
}

// GrafanaDashboardInstanceLabels returns the Grafana labels a dashboard instanceSelector must match.
// Those labels are app.kubernetes.io/component and app.kubernetes.io/part-of, including overrides from
// spec.grafana.labels.
func GrafanaDashboardInstanceLabels(cr *monv1.PlatformMonitoring) map[string]string {
	component := "grafana"
	partOf := "monitoring"
	if cr != nil && cr.Spec.Grafana != nil && cr.Spec.Grafana.Labels != nil {
		if value := cr.Spec.Grafana.Labels[labelComponent]; value != "" {
			component = value
		}
		if value := cr.Spec.Grafana.Labels[labelPartOf]; value != "" {
			partOf = value
		}
	}
	return map[string]string{
		labelComponent: component,
		labelPartOf:    partOf,
	}
}

// GrafanaNamespace returns the namespace of the Grafana resource this operator creates.
func GrafanaNamespace(cr *monv1.PlatformMonitoring) string {
	if cr.Spec.Grafana != nil && cr.Spec.Grafana.Namespace != "" {
		return cr.Spec.Grafana.Namespace
	}
	return cr.Namespace
}

func labelSelectorMatches(selector *metav1.LabelSelector, lbls map[string]string) (bool, error) {
	parsed, err := metav1.LabelSelectorAsSelector(selector)
	if err != nil {
		return false, err
	}
	return parsed.Matches(labels.Set(lbls)), nil
}
