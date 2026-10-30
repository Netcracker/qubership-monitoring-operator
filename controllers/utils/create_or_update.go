package utils

import (
	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"k8s.io/apimachinery/pkg/api/errors"
)

// SyncFunc copies the fields a controller manages from desired onto existing and
// reports whether existing changed.
type SyncFunc[T K8sResource] func(existing, desired T) bool

// CreateOrUpdate reads existing from the API server. It creates desired when the
// object does not exist. Otherwise it applies sync and sends an Update only when
// sync reports a change. existing must carry the name and namespace to read.
func CreateOrUpdate[T K8sResource](r *ComponentReconciler, cr *monv1.PlatformMonitoring, desired, existing T, sync SyncFunc[T]) error {
	found, err := r.getExisting(existing)
	if err != nil || !found {
		return r.createIfMissing(cr, desired, err)
	}
	return r.UpdateIfChanged(existing, sync(existing, desired))
}

// CreateOrUpdateWorkload works like CreateOrUpdate for a Deployment, DaemonSet, or
// StatefulSet. Because spec.selector is immutable, it deletes and recreates the
// workload when the selector differs instead of updating it.
func CreateOrUpdateWorkload[T K8sResource](r *ComponentReconciler, cr *monv1.PlatformMonitoring, desired, existing T, sync SyncFunc[T]) error {
	found, err := r.getExisting(existing)
	if err != nil || !found {
		return r.createIfMissing(cr, desired, err)
	}
	if WorkloadNeedsSelectorReplace(existing, desired) {
		if err := r.DeleteResource(existing); err != nil {
			return err
		}
		return r.CreateResource(cr, desired)
	}
	return r.UpdateIfChanged(existing, sync(existing, desired))
}

// CreateOrUpdateReplacingLegacy works like CreateOrUpdate, but before it creates
// desired it deletes legacy, an object stored under a previous name, if legacy
// exists. Errors from reading legacy are ignored.
func CreateOrUpdateReplacingLegacy[T K8sResource](r *ComponentReconciler, cr *monv1.PlatformMonitoring, desired, existing, legacy T, sync SyncFunc[T]) error {
	found, err := r.getExisting(existing)
	if err != nil {
		return err
	}
	if !found {
		if err := r.deleteIfPresent(legacy); err != nil {
			return err
		}
		return r.CreateResource(cr, desired)
	}
	return r.UpdateIfChanged(existing, sync(existing, desired))
}

// UpdateIfChanged sends an Update for o only when changed is true.
func (r *ComponentReconciler) UpdateIfChanged(o K8sResource, changed bool) error {
	if !changed {
		return nil
	}
	return r.UpdateResource(o)
}

// getExisting reads o. It returns false and no error when o does not exist.
func (r *ComponentReconciler) getExisting(o K8sResource) (bool, error) {
	err := r.GetResource(o)
	if errors.IsNotFound(err) {
		return false, nil
	}
	return err == nil, err
}

// createIfMissing returns getErr when reading the existing object failed and
// creates desired otherwise.
func (r *ComponentReconciler) createIfMissing(cr *monv1.PlatformMonitoring, desired K8sResource, getErr error) error {
	if getErr != nil {
		return getErr
	}
	return r.CreateResource(cr, desired)
}

// deleteIfPresent deletes o when it can be read.
func (r *ComponentReconciler) deleteIfPresent(o K8sResource) error {
	if r.GetResource(o) != nil {
		return nil
	}
	return r.DeleteResource(o)
}
