package controllers

import (
	"reflect"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	corev1 "k8s.io/api/core/v1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
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

// ownedChildPredicate enqueues an owned child on create, confirmed delete, a
// generation or label change, or a body change while generation stays 0.
// Status, resourceVersion, and managedFields do not enqueue.
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
			if !labels.Equals(e.ObjectOld.GetLabels(), e.ObjectNew.GetLabels()) {
				return true
			}
			if e.ObjectOld.GetGeneration() != 0 {
				return false
			}
			return generationlessChildChanged(e.ObjectOld, e.ObjectNew)
		},
	}
}

// generationlessChildChanged reports whether two generation-0 objects differ
// outside status, resourceVersion, and managedFields. A conversion failure
// enqueues, so a body edit is not dropped when the objects cannot be compared.
func generationlessChildChanged(oldObj, newObj client.Object) bool {
	oldContent, err := runtime.DefaultUnstructuredConverter.ToUnstructured(oldObj)
	if err != nil {
		return true
	}
	newContent, err := runtime.DefaultUnstructuredConverter.ToUnstructured(newObj)
	if err != nil {
		return true
	}
	stripGenerationlessChildNoise(oldContent)
	stripGenerationlessChildNoise(newContent)
	return !apiequality.Semantic.DeepEqual(oldContent, newContent)
}

func stripGenerationlessChildNoise(obj map[string]any) {
	delete(obj, "status")
	meta, _ := obj["metadata"].(map[string]any)
	if meta == nil {
		return
	}
	delete(meta, "resourceVersion")
	delete(meta, "managedFields")
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
