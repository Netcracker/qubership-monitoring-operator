package controllers

import (
	"reflect"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// platformMonitoringPredicate enqueues PlatformMonitoring when Spec-driven
// generation changes, or when labels/annotations change. Those maps are copied
// onto managed children, so a generation-only filter would miss them. Status
// and other metadata-only churn is ignored.
func platformMonitoringPredicate() predicate.Predicate {
	return predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool { return true },
		DeleteFunc: func(e event.DeleteEvent) bool { return !e.DeleteStateUnknown },
		UpdateFunc: func(e event.UpdateEvent) bool {
			if e.ObjectOld == nil || e.ObjectNew == nil {
				return false
			}
			if e.ObjectOld.GetGeneration() != e.ObjectNew.GetGeneration() {
				return true
			}
			if !labels.Equals(e.ObjectOld.GetLabels(), e.ObjectNew.GetLabels()) {
				return true
			}
			return !reflect.DeepEqual(e.ObjectOld.GetAnnotations(), e.ObjectNew.GetAnnotations())
		},
	}
}

// ownedChildPredicate enqueues owned children on spec/generation or label
// changes, including deletes. Status-only updates (resourceVersion/status
// subresource) must not re-enter the full PlatformMonitoring loop.
func ownedChildPredicate() predicate.Predicate {
	return predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool { return true },
		DeleteFunc: func(e event.DeleteEvent) bool { return !e.DeleteStateUnknown },
		UpdateFunc: func(e event.UpdateEvent) bool {
			if e.ObjectOld == nil || e.ObjectNew == nil {
				return false
			}
			if e.ObjectOld.GetGeneration() != e.ObjectNew.GetGeneration() {
				return true
			}
			return !labels.Equals(e.ObjectOld.GetLabels(), e.ObjectNew.GetLabels())
		},
	}
}

// jaegerServicePredicate matches Services used as Grafana Jaeger datasources.
func jaegerServicePredicate() predicate.Predicate {
	selector := labels.SelectorFromSet(utils.JaegerServiceLabels)
	matches := func(obj interface{}) bool {
		svc, ok := obj.(*corev1.Service)
		if !ok || svc == nil {
			return false
		}
		return selector.Matches(labels.Set(svc.GetLabels()))
	}
	return predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool { return matches(e.Object) },
		DeleteFunc: func(e event.DeleteEvent) bool { return matches(e.Object) },
		UpdateFunc: func(e event.UpdateEvent) bool {
			if !matches(e.ObjectNew) && !matches(e.ObjectOld) {
				return false
			}
			oldSvc, _ := e.ObjectOld.(*corev1.Service)
			newSvc, _ := e.ObjectNew.(*corev1.Service)
			if oldSvc == nil || newSvc == nil {
				return true
			}
			if !labels.Equals(oldSvc.GetLabels(), newSvc.GetLabels()) {
				return true
			}
			return !reflect.DeepEqual(oldSvc.Spec, newSvc.Spec)
		},
	}
}

func jaegerDatasourceRequested(pm *monv1.PlatformMonitoring) bool {
	return pm != nil &&
		pm.Spec.Integration != nil &&
		pm.Spec.Integration.Jaeger != nil &&
		pm.Spec.Integration.Jaeger.CreateGrafanaDataSource
}
