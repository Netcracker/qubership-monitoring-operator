package prometheus_operator

import (
	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *PrometheusOperatorReconciler) handleRole(cr *monv1.PlatformMonitoring) error {
	m, err := prometheusOperatorRole(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Role manifest")
		return err
	}

	e := &rbacv1.Role{ObjectMeta: m.ObjectMeta}
	return utils.CreateOrUpdate(r.ComponentReconciler, cr, m, e, utils.SyncRole)
}

func (r *PrometheusOperatorReconciler) handleServiceAccount(cr *monv1.PlatformMonitoring) error {
	m, err := prometheusOperatorServiceAccount(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ServiceAccount manifest")
		return err
	}

	e := &corev1.ServiceAccount{ObjectMeta: m.ObjectMeta}
	return utils.CreateOrUpdate(r.ComponentReconciler, cr, m, e, utils.SyncLabels)
}

func (r *PrometheusOperatorReconciler) handleRoleBinding(cr *monv1.PlatformMonitoring) error {
	m, err := prometheusOperatorRoleBinding(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating RoleBinding manifest")
		return err
	}

	e := &rbacv1.RoleBinding{ObjectMeta: m.ObjectMeta}
	return utils.CreateOrUpdate(r.ComponentReconciler, cr, m, e, utils.SyncLabels)
}

func (r *PrometheusOperatorReconciler) handleClusterRole(cr *monv1.PlatformMonitoring) error {
	m, err := prometheusOperatorClusterRole(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ClusterRole manifest")
		return err
	}

	e := &rbacv1.ClusterRole{ObjectMeta: m.ObjectMeta}
	return utils.CreateOrUpdate(r.ComponentReconciler, cr, m, e, utils.SyncClusterRole)
}

func (r *PrometheusOperatorReconciler) handleClusterRoleBinding(cr *monv1.PlatformMonitoring) error {
	m, err := prometheusOperatorClusterRoleBinding(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ClusterRoleBinding manifest")
		return err
	}

	e := &rbacv1.ClusterRoleBinding{ObjectMeta: m.ObjectMeta}
	return utils.CreateOrUpdate(r.ComponentReconciler, cr, m, e, utils.SyncLabels)
}

func (r *PrometheusOperatorReconciler) handleDeployment(cr *monv1.PlatformMonitoring) error {
	isOpenShift, err := r.IsOpenShift()
	if err != nil {
		return err
	}
	m, err := prometheusOperatorDeployment(cr, isOpenShift)
	if err != nil {
		r.Log.Error(err, "Failed creating Deployment manifest")
		return err
	}
	e := &appsv1.Deployment{ObjectMeta: m.ObjectMeta}
	return utils.CreateOrUpdateWorkload(r.ComponentReconciler, cr, m, e, syncDeployment)
}

// syncDeployment copies the Deployment fields this controller manages.
func syncDeployment(existing, desired *appsv1.Deployment) bool {
	return utils.AnyChanged(
		utils.SetLabelsIfChanged(existing, desired.GetLabels()),
		utils.SetIfChanged(&existing.Spec.Selector, desired.Spec.Selector),
		utils.SetLabelsIfChanged(&existing.Spec.Template, desired.Spec.Template.GetLabels()),
		utils.SyncPodSpec(&existing.Spec.Template.Spec, &desired.Spec.Template.Spec),
		utils.SetIfChanged(&existing.Spec.Template.Spec.NodeSelector, desired.Spec.Template.Spec.NodeSelector),
		utils.SetIfChanged(&existing.Spec.Template.Spec.Tolerations, desired.Spec.Template.Spec.Tolerations),
	)
}

func (r *PrometheusOperatorReconciler) handleService(cr *monv1.PlatformMonitoring) error {
	m, err := prometheusOperatorService(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Service manifest")
		return err
	}

	m.Spec.Selector = map[string]string{"app.kubernetes.io/name": utils.TruncLabel(utils.PrometheusOperatorComponentName)}

	e := &corev1.Service{ObjectMeta: m.ObjectMeta}
	return utils.CreateOrUpdate(r.ComponentReconciler, cr, m, e, utils.SyncService)
}

func (r *PrometheusOperatorReconciler) handlePodMonitor(cr *monv1.PlatformMonitoring) error {
	m, err := prometheusOperatorPodMonitor(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating PodMonitor manifest")
		return err
	}

	e := &promv1.PodMonitor{ObjectMeta: m.ObjectMeta}
	return utils.CreateOrUpdate(r.ComponentReconciler, cr, m, e, utils.SyncPodMonitor)
}

func (r *PrometheusOperatorReconciler) deleteServiceAccount(cr *monv1.PlatformMonitoring) error {
	m, err := prometheusOperatorServiceAccount(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ServiceAccount manifest")
		return err
	}
	e := &corev1.ServiceAccount{ObjectMeta: m.ObjectMeta}
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

func (r *PrometheusOperatorReconciler) deleteRole(cr *monv1.PlatformMonitoring) error {
	m, err := prometheusOperatorRole(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Role manifest")
		return err
	}
	e := &rbacv1.Role{ObjectMeta: m.ObjectMeta}
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

func (r *PrometheusOperatorReconciler) deleteRoleBinding(cr *monv1.PlatformMonitoring) error {
	m, err := prometheusOperatorRoleBinding(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating RoleBinding manifest")
		return err
	}
	e := &rbacv1.RoleBinding{ObjectMeta: m.ObjectMeta}
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

func (r *PrometheusOperatorReconciler) deleteClusterRole(cr *monv1.PlatformMonitoring) error {
	m, err := prometheusOperatorClusterRole(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ClusterRole manifest")
		return err
	}
	e := &rbacv1.ClusterRole{ObjectMeta: m.ObjectMeta}
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

func (r *PrometheusOperatorReconciler) deleteClusterRoleBinding(cr *monv1.PlatformMonitoring) error {
	m, err := prometheusOperatorClusterRoleBinding(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ClusterRoleBinding manifest")
		return err
	}
	e := &rbacv1.ClusterRoleBinding{ObjectMeta: m.ObjectMeta}
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

func (r *PrometheusOperatorReconciler) deletePrometheusOperatorDeployment(cr *monv1.PlatformMonitoring) error {
	// Deletion targets are addressed by their stable name and namespace so that uninstall does not
	// depend on platform discovery or on validation of the desired state.
	e := &appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Name: utils.PrometheusOperatorComponentName, Namespace: cr.GetNamespace()}}
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

func (r *PrometheusOperatorReconciler) deleteService(cr *monv1.PlatformMonitoring) error {
	m, err := prometheusOperatorService(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Service manifest")
		return err
	}
	e := &corev1.Service{ObjectMeta: m.ObjectMeta}
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

func (r *PrometheusOperatorReconciler) deletePodMonitor(cr *monv1.PlatformMonitoring) error {
	m, err := prometheusOperatorPodMonitor(cr)
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
