package prometheus

import (
	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
)

func (r *PrometheusReconciler) handleServiceAccount(cr *monv1.PlatformMonitoring) error {
	m, err := prometheusServiceAccount(cr)
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

func (r *PrometheusReconciler) handleClusterRole(cr *monv1.PlatformMonitoring) error {
	m, err := prometheusClusterRole(cr)
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

func (r *PrometheusReconciler) handleClusterRoleBinding(cr *monv1.PlatformMonitoring) error {
	m, err := prometheusClusterRoleBinding(cr)
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

func (r *PrometheusReconciler) handlePrometheus(cr *monv1.PlatformMonitoring) error {
	m, err := prometheus(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Prometheus manifest")
		return err
	}
	e := &promv1.Prometheus{ObjectMeta: m.ObjectMeta}
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
	changed = utils.SetIfChanged(&e.Spec, m.Spec) || changed
	if !changed {
		return nil
	}

	if err = r.UpdateResource(e); err != nil {
		return err
	}
	return nil
}

func (r *PrometheusReconciler) handleIngressV1(cr *monv1.PlatformMonitoring) error {
	m, err := prometheusIngressV1(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Ingress manifest")
		return err
	}
	e := &networkingv1.Ingress{ObjectMeta: m.ObjectMeta}
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
	changed = utils.SetAnnotationsIfChanged(e, m.GetAnnotations()) || changed
	changed = utils.SetIfChanged(&e.Spec.Rules, m.Spec.Rules) || changed
	changed = utils.SetIfChanged(&e.Spec.TLS, m.Spec.TLS) || changed
	if !changed {
		return nil
	}

	if err = r.UpdateResource(e); err != nil {
		return err
	}
	return nil
}

func (r *PrometheusReconciler) handlePodMonitor(cr *monv1.PlatformMonitoring) error {
	m, err := prometheusPodMonitor(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating PodMonitor manifest")
		return err
	}

	e := &promv1.PodMonitor{ObjectMeta: m.ObjectMeta}
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
	changed = utils.SetIfChanged(&e.Spec.PodMetricsEndpoints, m.Spec.PodMetricsEndpoints) || changed
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

func (r *PrometheusReconciler) deleteServiceAccount(cr *monv1.PlatformMonitoring) error {
	m, err := prometheusServiceAccount(cr)
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

func (r *PrometheusReconciler) deleteClusterRole(cr *monv1.PlatformMonitoring) error {
	m, err := prometheusClusterRole(cr)
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

func (r *PrometheusReconciler) deleteClusterRoleBinding(cr *monv1.PlatformMonitoring) error {
	m, err := prometheusClusterRoleBinding(cr)
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

func (r *PrometheusReconciler) deletePrometheus(cr *monv1.PlatformMonitoring) error {
	m, err := prometheus(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Prometheus manifest")
		return err
	}
	e := &promv1.Prometheus{ObjectMeta: m.ObjectMeta}
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

func (r *PrometheusReconciler) deleteIngressV1(cr *monv1.PlatformMonitoring) error {
	m, err := prometheusIngressV1(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Ingress manifest")
		return err
	}
	e := &networkingv1.Ingress{ObjectMeta: m.ObjectMeta}
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

func (r *PrometheusReconciler) deletePodMonitor(cr *monv1.PlatformMonitoring) error {
	m, err := prometheusPodMonitor(cr)
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
