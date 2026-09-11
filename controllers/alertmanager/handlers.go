package alertmanager

import (
	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
)

func (r *AlertManagerReconciler) handleServiceAccount(cr *monv1.PlatformMonitoring) error {
	m, err := alertmanagerServiceAccount(cr)
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

func (r *AlertManagerReconciler) handleSecret(cr *monv1.PlatformMonitoring) error {
	m, err := alertmanagerSecret(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Secret manifest")
		return err
	}

	e := &corev1.Secret{ObjectMeta: m.ObjectMeta}
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

func (r *AlertManagerReconciler) handleAlertmanager(cr *monv1.PlatformMonitoring) error {
	m, err := alertmanager(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Alertmanager manifest")
		return err
	}

	e := &promv1.Alertmanager{ObjectMeta: m.ObjectMeta}
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

func (r *AlertManagerReconciler) handleService(cr *monv1.PlatformMonitoring) error {
	m, err := alertmanagerService(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Service manifest")
		return err
	}

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

func (r *AlertManagerReconciler) handleIngressV1(cr *monv1.PlatformMonitoring) error {
	m, err := alertmanagerIngressV1(cr)
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

func (r *AlertManagerReconciler) handlePodMonitor(cr *monv1.PlatformMonitoring) error {
	m, err := alertmanagerPodMonitor(cr)
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

func (r *AlertManagerReconciler) deleteServiceAccount(cr *monv1.PlatformMonitoring) error {

	m, err := alertmanagerServiceAccount(cr)
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

func (r *AlertManagerReconciler) deleteSecret(cr *monv1.PlatformMonitoring) error {
	m, err := alertmanagerSecret(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Secret manifest")
		return err
	}
	e := &corev1.Secret{ObjectMeta: m.ObjectMeta}
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

func (r *AlertManagerReconciler) deleteAlertmanager(cr *monv1.PlatformMonitoring) error {
	m, err := alertmanager(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Alertmanager manifest")
		return err
	}
	e := &promv1.Alertmanager{ObjectMeta: m.ObjectMeta}
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

func (r *AlertManagerReconciler) deleteService(cr *monv1.PlatformMonitoring) error {
	m, err := alertmanagerService(cr)
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

func (r *AlertManagerReconciler) deleteIngressV1(cr *monv1.PlatformMonitoring) error {
	m, err := alertmanagerIngressV1(cr)
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

func (r *AlertManagerReconciler) deletePodMonitor(cr *monv1.PlatformMonitoring) error {
	m, err := alertmanagerPodMonitor(cr)
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
