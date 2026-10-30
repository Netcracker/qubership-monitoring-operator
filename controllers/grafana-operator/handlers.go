package grafana_operator

import (
	"context"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	grafv1 "github.com/grafana/grafana-operator/v5/api/v1beta1"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/util/retry"
)

func (r *GrafanaOperatorReconciler) handleServiceAccount(cr *monv1.PlatformMonitoring) error {
	m, err := grafanaOperatorServiceAccount(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ServiceAccount manifest")
		return err
	}

	e := &corev1.ServiceAccount{ObjectMeta: m.ObjectMeta}
	return utils.CreateOrUpdate(r.ComponentReconciler, cr, m, e, utils.SyncLabels)
}

func (r *GrafanaOperatorReconciler) handleClusterRole(cr *monv1.PlatformMonitoring) error {
	m, err := grafanaOperatorClusterRole(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ClusterRole manifest")
		return err
	}

	e := &rbacv1.ClusterRole{ObjectMeta: m.ObjectMeta}
	return utils.CreateOrUpdate(r.ComponentReconciler, cr, m, e, utils.SyncClusterRole)
}

func (r *GrafanaOperatorReconciler) handleClusterRoleBinding(cr *monv1.PlatformMonitoring) error {
	m, err := grafanaOperatorClusterRoleBinding(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ClusterRoleBinding manifest")
		return err
	}

	e := &rbacv1.ClusterRoleBinding{ObjectMeta: m.ObjectMeta}
	return utils.CreateOrUpdate(r.ComponentReconciler, cr, m, e, utils.SyncLabels)
}

func (r *GrafanaOperatorReconciler) handleRole(cr *monv1.PlatformMonitoring) error {
	m, err := grafanaOperatorRole(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Role manifest")
		return err
	}

	e := &rbacv1.Role{ObjectMeta: m.ObjectMeta}
	return utils.CreateOrUpdate(r.ComponentReconciler, cr, m, e, utils.SyncRole)
}

func (r *GrafanaOperatorReconciler) handleRoleBinding(cr *monv1.PlatformMonitoring) error {
	m, err := grafanaOperatorRoleBinding(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating RoleBinding manifest")
		return err
	}

	e := &rbacv1.RoleBinding{ObjectMeta: m.ObjectMeta}
	return utils.CreateOrUpdate(r.ComponentReconciler, cr, m, e, utils.SyncLabels)
}

func (r *GrafanaOperatorReconciler) handleDeployment(cr *monv1.PlatformMonitoring) error {
	isOpenShift, err := r.IsOpenShift()
	if err != nil {
		return err
	}
	m, err := grafanaOperatorDeployment(cr, isOpenShift)
	if err != nil {
		r.Log.Error(err, "Failed creating Deployment manifest")
		return err
	}
	e := &appsv1.Deployment{ObjectMeta: m.ObjectMeta}
	return utils.CreateOrUpdateWorkload(r.ComponentReconciler, cr, m, e, syncDeployment)
}

// syncDeployment copies the Deployment fields this controller manages.
// Deployment.spec.selector is immutable, so it keeps the selector already stored in the
// API server and merges template labels so matchLabels stay satisfied (upgrades / Helm).
func syncDeployment(existing, desired *appsv1.Deployment) bool {
	selectorChanged := existing.Spec.Selector == nil && utils.SetIfChanged(&existing.Spec.Selector, desired.Spec.Selector)
	templateLabels := deploymentTemplateLabelsForUpdate(
		existing.Spec.Template.Labels, desired.Spec.Template.Labels, existing.Spec.Selector)
	return utils.AnyChanged(
		utils.SetLabelsIfChanged(existing, desired.GetLabels()),
		selectorChanged,
		utils.SetLabelsIfChanged(&existing.Spec.Template, templateLabels),
		utils.SyncPodSpec(&existing.Spec.Template.Spec, &desired.Spec.Template.Spec),
		utils.SetIfChanged(&existing.Spec.Template.Spec.NodeSelector, desired.Spec.Template.Spec.NodeSelector),
	)
}

func (r *GrafanaOperatorReconciler) handleGrafanaDashboard(fileName string, cr *monv1.PlatformMonitoring) error {
	m, err := grafanaDashboard(cr, fileName)
	if err != nil {
		r.Log.Error(err, "Failed creating GrafanaDashboard manifest")
		return err
	}
	// Retry on conflict: the grafana-operator (and other actors) write status/finalizers
	// on GrafanaDashboard objects concurrently, so a Get-mutate-Update window can race.
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// Explicitly set GVK to ensure correct API group (grafana.integreatly.org/v1beta1) is used.
		// Required for Grafana Operator v5 migration from integreatly.org/v1alpha1.
		checkObj := &grafv1.GrafanaDashboard{}
		checkObj.SetName(m.GetName())
		checkObj.SetNamespace(m.GetNamespace())
		checkObj.SetGroupVersionKind(schema.GroupVersionKind{Group: "grafana.integreatly.org", Version: "v1beta1", Kind: "GrafanaDashboard"})
		if err := r.GetResource(checkObj); err != nil {
			if errors.IsNotFound(err) {
				return r.CreateResource(cr, m)
			}
			return err
		}

		// Only update when Spec or labels actually changed to avoid rewriting large
		// dashboard CRs on every PlatformMonitoring reconcile (#447).
		return r.UpdateIfChanged(checkObj, utils.AnyChanged(
			utils.SetIfChanged(&checkObj.Spec, m.Spec),
			utils.SetLabelsIfChanged(checkObj, m.GetLabels()),
		))
	})
}

func (r *GrafanaOperatorReconciler) handlePodMonitor(cr *monv1.PlatformMonitoring) error {
	m, err := grafanaOperatorPodMonitor(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating PodMonitor manifest")
		return err
	}

	e := &promv1.PodMonitor{ObjectMeta: m.ObjectMeta}
	return utils.CreateOrUpdate(r.ComponentReconciler, cr, m, e, utils.SyncPodMonitor)
}

func (r *GrafanaOperatorReconciler) deleteGrafanaOperatorDeployment(cr *monv1.PlatformMonitoring) error {
	// Deletion targets are addressed by their stable name and namespace so that uninstall does not
	// depend on platform discovery or on validation of the desired state.
	e := &appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Name: utils.GrafanaOperatorComponentName, Namespace: cr.GetNamespace()}}
	if err := r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if err := r.DeleteResource(e); err != nil {
		return err
	}
	return nil
}

func (r *GrafanaOperatorReconciler) deleteGrafanaDashboard(fileName string, cr *monv1.PlatformMonitoring) error {
	m, err := grafanaDashboard(cr, fileName)
	if err != nil {
		r.Log.Error(err, "Failed creating GrafanaDashboard manifest")
		return err
	}
	// Check if resource exists first
	// Explicitly set GVK to ensure correct API group (grafana.integreatly.org/v1beta1) is used
	// This is required for Grafana Operator v5 migration from integreatly.org/v1alpha1
	checkObj := &grafv1.GrafanaDashboard{}
	checkObj.SetName(m.GetName())
	checkObj.SetNamespace(m.GetNamespace())
	checkObj.SetGroupVersionKind(schema.GroupVersionKind{Group: "grafana.integreatly.org", Version: "v1beta1", Kind: "GrafanaDashboard"})
	if err = r.GetResource(checkObj); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	// Use the manifest object (which has correct type) for deletion
	// The manifest object already has GVK set correctly
	if err = r.Client.Delete(context.TODO(), m); err != nil {
		return err
	}
	r.Log.Info("Successful deleting", "resource", "GrafanaDashboard", "name", m.GetName())
	return nil
}

func (r *GrafanaOperatorReconciler) deletePodMonitor(cr *monv1.PlatformMonitoring) error {
	m, err := grafanaOperatorPodMonitor(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating PodMonitor manifest")
		return err
	}
	e := &promv1.PodMonitor{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if err = r.DeleteResource(e); err != nil {
		return err
	}
	return nil
}

func mergeStringStringMaps(base, overlay map[string]string) map[string]string {
	out := make(map[string]string)
	for k, v := range base {
		out[k] = v
	}
	for k, v := range overlay {
		out[k] = v
	}
	return out
}

// deploymentTemplateLabelsForUpdate merges desired template labels into existing ones and
// forces keys from spec.selector.matchLabels so the pod template stays valid for the API.
func deploymentTemplateLabelsForUpdate(existing, desired map[string]string, sel *metav1.LabelSelector) map[string]string {
	merged := mergeStringStringMaps(existing, desired)
	if sel != nil {
		for k, v := range sel.MatchLabels {
			merged[k] = v
		}
	}
	if len(merged) == 0 {
		return nil
	}
	return merged
}
