package grafana_operator

import (
	"context"
	"reflect"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/grafana"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	grafv1 "github.com/grafana/grafana-operator/v5/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var grafanaDashboardGVK = schema.GroupVersionKind{
	Group:   "grafana.integreatly.org",
	Version: "v1beta1",
	Kind:    "GrafanaDashboard",
}

// bundledDashboardSelected reports whether a chart dashboard matches the PlatformMonitoring dashboard selectors.
// The reconciler creates the dashboard when this reports true and deletes it when this reports false.
// When PrivilegedRights is false, only dashboardLabelSelector is applied. The namespace selector needs
// permission to read Namespace objects, which the namespace Role does not have.
func (r *GrafanaOperatorReconciler) bundledDashboardSelected(cr *monv1.PlatformMonitoring, fileName string) (bool, error) {
	dash, err := grafanaDashboard(cr, fileName)
	if err != nil {
		return false, err
	}
	labelsMatch, err := grafana.DashboardLabelsMatch(cr.Spec.Grafana.DashboardLabelSelector, dash.GetLabels())
	if err != nil || !labelsMatch || !utils.PrivilegedRights {
		return labelsMatch, err
	}
	ns := &corev1.Namespace{}
	ns.SetName(cr.GetNamespace())
	if err := r.GetResource(ns); err != nil {
		return false, err
	}
	return grafana.NamespaceSelectorMatches(cr.Spec.Grafana.DashboardNamespaceSelector, ns.GetLabels())
}

func (r *GrafanaOperatorReconciler) reconcileDashboardSelection(cr *monv1.PlatformMonitoring) error {
	if cr.Spec.Grafana == nil {
		return nil
	}
	if utils.PrivilegedRights {
		return r.reconcileDashboardSelectionPrivileged(cr)
	}
	return r.reconcileDashboardSelectionInNamespace(cr, cr.GetNamespace(), nil, false)
}

func (r *GrafanaOperatorReconciler) reconcileDashboardSelectionPrivileged(cr *monv1.PlatformMonitoring) error {
	nsList := &corev1.NamespaceList{}
	if err := r.Client.List(context.TODO(), nsList); err != nil {
		return err
	}
	namespaceLabels := make(map[string]map[string]string, len(nsList.Items))
	for i := range nsList.Items {
		ns := &nsList.Items[i]
		namespaceLabels[ns.Name] = ns.GetLabels()
	}

	dashList := &grafv1.GrafanaDashboardList{}
	if err := r.Client.List(context.TODO(), dashList); err != nil {
		return err
	}
	var errs []error
	for i := range dashList.Items {
		dash := &dashList.Items[i]
		labels, found := namespaceLabels[dash.Namespace]
		if !found {
			if !dashboardSelectorAnnotated(dash) {
				continue
			}
			if err := r.writeDashboardSelection(cr, client.ObjectKeyFromObject(dash), grafana.DashboardSelectionDetach); err != nil {
				errs = append(errs, err)
			}
			continue
		}
		action, err := grafana.SelectDashboard(
			cr.Spec.Grafana.DashboardLabelSelector,
			cr.Spec.Grafana.DashboardNamespaceSelector,
			dash.GetLabels(),
			labels,
			dashboardSelectorAnnotated(dash),
			true,
		)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if action == grafana.DashboardSelectionLeave {
			continue
		}
		if err := r.writeDashboardSelection(cr, client.ObjectKeyFromObject(dash), action); err != nil {
			errs = append(errs, err)
		}
	}
	return utilerrors.NewAggregate(errs)
}

func (r *GrafanaOperatorReconciler) reconcileDashboardSelectionInNamespace(cr *monv1.PlatformMonitoring, namespace string, namespaceLabels map[string]string, honorNamespaceSelector bool) error {
	dashList := &grafv1.GrafanaDashboardList{}
	if err := r.Client.List(context.TODO(), dashList, client.InNamespace(namespace)); err != nil {
		return err
	}
	var errs []error
	for i := range dashList.Items {
		dash := &dashList.Items[i]
		action, err := grafana.SelectDashboard(
			cr.Spec.Grafana.DashboardLabelSelector,
			cr.Spec.Grafana.DashboardNamespaceSelector,
			dash.GetLabels(),
			namespaceLabels,
			dashboardSelectorAnnotated(dash),
			honorNamespaceSelector,
		)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if action == grafana.DashboardSelectionLeave {
			continue
		}
		if err := r.writeDashboardSelection(cr, client.ObjectKeyFromObject(dash), action); err != nil {
			errs = append(errs, err)
		}
	}
	return utilerrors.NewAggregate(errs)
}

func (r *GrafanaOperatorReconciler) writeDashboardSelection(cr *monv1.PlatformMonitoring, key client.ObjectKey, action grafana.DashboardSelectionAction) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		current := &grafv1.GrafanaDashboard{}
		current.SetName(key.Name)
		current.SetNamespace(key.Namespace)
		current.SetGroupVersionKind(grafanaDashboardGVK)
		if err := r.GetResource(current); err != nil {
			if apierrors.IsNotFound(err) {
				return nil
			}
			return err
		}
		updated := dashboardWithSelection(current, cr, action)
		if reflect.DeepEqual(current.Spec, updated.Spec) && reflect.DeepEqual(current.GetAnnotations(), updated.GetAnnotations()) {
			return nil
		}
		if err := r.UpdateResource(updated); err != nil {
			if apierrors.IsInvalid(err) {
				return r.updateDashboardSelectionFallback(current, cr, action)
			}
			return err
		}
		return nil
	})
}

// updateDashboardSelectionFallback keeps an immutable instanceSelector in place.
// Grafana Operator v5 rejects changes to spec.instanceSelector and rejects turning
// spec.allowCrossNamespaceImport off. Suspending the dashboard stops import when the
// pointer cannot be removed. Adoption still records the annotation and unsuspends.
func (r *GrafanaOperatorReconciler) updateDashboardSelectionFallback(current *grafv1.GrafanaDashboard, cr *monv1.PlatformMonitoring, action grafana.DashboardSelectionAction) error {
	updated := current.DeepCopy()
	switch action {
	case grafana.DashboardSelectionAdopt:
		ensureDashboardSelectorAnnotation(updated, cr)
		if updated.Namespace != grafana.GrafanaNamespace(cr) {
			updated.Spec.AllowCrossNamespaceImport = true
		}
		updated.Spec.Suspend = false
	case grafana.DashboardSelectionDetach:
		ensureDashboardSelectorAnnotation(updated, cr)
		updated.Spec.Suspend = true
	default:
		return nil
	}
	if reflect.DeepEqual(current.Spec, updated.Spec) && reflect.DeepEqual(current.GetAnnotations(), updated.GetAnnotations()) {
		return nil
	}
	return r.UpdateResource(updated)
}

func dashboardWithSelection(current *grafv1.GrafanaDashboard, cr *monv1.PlatformMonitoring, action grafana.DashboardSelectionAction) *grafv1.GrafanaDashboard {
	updated := current.DeepCopy()
	switch action {
	case grafana.DashboardSelectionAdopt:
		updated.Spec.InstanceSelector = &metav1.LabelSelector{
			MatchLabels: grafana.GrafanaDashboardInstanceLabels(cr),
		}
		updated.Spec.AllowCrossNamespaceImport = updated.Namespace != grafana.GrafanaNamespace(cr)
		updated.Spec.Suspend = false
		ensureDashboardSelectorAnnotation(updated, cr)
	case grafana.DashboardSelectionDetach:
		updated.Spec.InstanceSelector = nil
		updated.Spec.AllowCrossNamespaceImport = false
		updated.Spec.Suspend = false
		removeDashboardSelectorAnnotation(updated)
	}
	return updated
}

func dashboardSelectorAnnotated(dash *grafv1.GrafanaDashboard) bool {
	if dash.Annotations == nil {
		return false
	}
	_, ok := dash.Annotations[grafana.DashboardSelectorAnnotation]
	return ok
}

func ensureDashboardSelectorAnnotation(dash *grafv1.GrafanaDashboard, cr *monv1.PlatformMonitoring) {
	if dash.Annotations == nil {
		dash.Annotations = map[string]string{}
	}
	dash.Annotations[grafana.DashboardSelectorAnnotation] = cr.GetName()
}

func removeDashboardSelectorAnnotation(dash *grafv1.GrafanaDashboard) {
	if dash.Annotations == nil {
		return
	}
	delete(dash.Annotations, grafana.DashboardSelectorAnnotation)
	if len(dash.Annotations) == 0 {
		dash.Annotations = nil
	}
}
