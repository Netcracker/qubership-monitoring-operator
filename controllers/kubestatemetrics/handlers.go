package kubestatemetrics

import (
	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
)

func (r *KubeStateMetricsReconciler) handleServiceAccount(cr *monv1.PlatformMonitoring) error {
	m, err := kubeStateMetricsServiceAccount(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ServiceAccount manifest")
		return err
	}

	e := &corev1.ServiceAccount{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			if err = r.CreateResource(cr, m); err != nil {
				return err
			}
			return nil
		}
		return err
	}
	//Set parameters
	changed := utils.SetLabelsIfChanged(e, m.GetLabels())
	if !changed {
		return nil
	}

	if err = r.UpdateResource(e); err != nil {
		return err
	}
	return nil
}

func (r *KubeStateMetricsReconciler) handleClusterRole(cr *monv1.PlatformMonitoring) error {
	m, err := kubeStateMetricsClusterRole(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ClusterRole manifest")
		return err
	}

	e := &rbacv1.ClusterRole{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			if err = r.CreateResource(cr, m); err != nil {
				return err
			}
			return nil
		}
		return err
	}

	//Set parameters
	changed := utils.SetLabelsIfChanged(e, m.GetLabels())
	changed = utils.SetIfChanged(&e.Rules, m.Rules) || changed
	if !changed {
		return nil
	}

	if err = r.UpdateResource(e); err != nil {
		return err
	}
	return nil
}

func (r *KubeStateMetricsReconciler) handleClusterRoleBinding(cr *monv1.PlatformMonitoring) error {
	m, err := kubeStateMetricsClusterRoleBinding(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ClusterRoleBinding manifest")
		return err
	}

	e := &rbacv1.ClusterRoleBinding{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			if err = r.CreateResource(cr, m); err != nil {
				return err
			}
			return nil
		}
		return err
	}
	//Set parameters
	changed := utils.SetLabelsIfChanged(e, m.GetLabels())
	if !changed {
		return nil
	}

	if err = r.UpdateResource(e); err != nil {
		return err
	}
	return nil
}

func (r *KubeStateMetricsReconciler) handleDeployment(cr *monv1.PlatformMonitoring) error {
	m, err := kubeStateMetricsDeployment(cr, r.HasIngressV1Api() || r.HasIngressV1beta1Api())
	if err != nil {
		r.Log.Error(err, "Failed creating Deployment manifest")
		return err
	}
	e := &appsv1.Deployment{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			if err = r.CreateResource(cr, m); err != nil {
				return err
			}
			return nil
		}
		return err
	}
	if utils.WorkloadNeedsSelectorReplace(e, m) {
		if err := r.DeleteResource(e); err != nil {
			return err
		}
		return r.CreateResource(cr, m)
	}

	//Set parameters
	changed := utils.SetLabelsIfChanged(e, m.GetLabels())
	changed = utils.SetLabelsIfChanged(&e.Spec.Template, m.Spec.Template.GetLabels()) || changed
	changed = utils.SetIfChanged(&e.Spec.Template.Spec.SecurityContext, m.Spec.Template.Spec.SecurityContext) || changed
	changed = utils.SetIfChanged(&e.Spec.Template.Spec.Containers, m.Spec.Template.Spec.Containers) || changed
	changed = utils.SetIfChanged(&e.Spec.Template.Spec.ServiceAccountName, m.Spec.Template.Spec.ServiceAccountName) || changed
	changed = utils.SetIfChanged(&e.Spec.Template.Spec.NodeSelector, m.Spec.Template.Spec.NodeSelector) || changed
	changed = utils.SetIfChanged(&e.Spec.Template.Spec.Affinity, m.Spec.Template.Spec.Affinity) || changed
	changed = utils.SetIfChanged(&e.Spec.Template.Spec.PriorityClassName, m.Spec.Template.Spec.PriorityClassName) || changed
	if !changed {
		return nil
	}

	if err = r.UpdateResource(e); err != nil {
		return err
	}
	return nil
}

func (r *KubeStateMetricsReconciler) handleService(cr *monv1.PlatformMonitoring) error {
	m, err := kubeStateMetricsService(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Service manifest")
		return err
	}

	m.Spec.Selector = map[string]string{"app.kubernetes.io/name": utils.TruncLabel(utils.KubestatemetricsComponentName)}

	e := &corev1.Service{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			if err = r.CreateResource(cr, m); err != nil {
				return err
			}
			return nil
		}
		return err
	}

	//Set parameters
	changed := utils.SetLabelsIfChanged(e, m.GetLabels())
	changed = utils.SetIfChanged(&e.Spec.Ports, m.Spec.Ports) || changed
	changed = utils.SetIfChanged(&e.Spec.Selector, m.Spec.Selector) || changed
	if !changed {
		return nil
	}

	if err = r.UpdateResource(e); err != nil {
		return err
	}
	return nil
}

func (r *KubeStateMetricsReconciler) handleServiceMonitor(cr *monv1.PlatformMonitoring) error {
	m, err := kubeStateMetricsServiceMonitor(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ServiceMonitor manifest")
		return err
	}

	e := &promv1.ServiceMonitor{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			if err = r.CreateResource(cr, m); err != nil {
				return err
			}
			return nil
		}
		return err
	}

	//Set parameters
	changed := utils.SetLabelsIfChanged(e, m.GetLabels())
	changed = utils.SetIfChanged(&e.Spec.JobLabel, m.Spec.JobLabel) || changed
	changed = utils.SetIfChanged(&e.Spec.Endpoints, m.Spec.Endpoints) || changed
	changed = utils.SetIfChanged(&e.Spec.NamespaceSelector, m.Spec.NamespaceSelector) || changed
	changed = utils.SetIfChanged(&e.Spec.Selector, m.Spec.Selector) || changed
	if !changed {
		return nil
	}

	if err = r.UpdateResource(e); err != nil {
		return err
	}
	return nil
}

func (r *KubeStateMetricsReconciler) deleteServiceAccount(cr *monv1.PlatformMonitoring) error {
	m, err := kubeStateMetricsServiceAccount(cr)
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

func (r *KubeStateMetricsReconciler) deleteClusterRole(cr *monv1.PlatformMonitoring) error {
	m, err := kubeStateMetricsClusterRole(cr)
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

func (r *KubeStateMetricsReconciler) deleteClusterRoleBinding(cr *monv1.PlatformMonitoring) error {
	m, err := kubeStateMetricsClusterRoleBinding(cr)
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

func (r *KubeStateMetricsReconciler) deleteDeployment(cr *monv1.PlatformMonitoring) error {
	m, err := kubeStateMetricsDeployment(cr, r.HasIngressV1Api() || r.HasIngressV1beta1Api())
	if err != nil {
		r.Log.Error(err, "Failed creating Deployment manifest")
		return err
	}
	e := &appsv1.Deployment{ObjectMeta: m.ObjectMeta}
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

func (r *KubeStateMetricsReconciler) deleteService(cr *monv1.PlatformMonitoring) error {
	m, err := kubeStateMetricsService(cr)
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

func (r *KubeStateMetricsReconciler) deleteServiceMonitor(cr *monv1.PlatformMonitoring) error {
	m, err := kubeStateMetricsServiceMonitor(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ServiceMonitor manifest")
		return err
	}
	e := &promv1.ServiceMonitor{ObjectMeta: m.ObjectMeta}
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
