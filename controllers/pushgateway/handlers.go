package pushgateway

import (
	v1alpha1 "github.com/Netcracker/qubership-monitoring-operator/api/v1alpha1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
)

func (r *PushgatewayReconciler) handleDeployment(cr *v1alpha1.PlatformMonitoring) error {
	m, err := pushgatewayDeployment(cr)
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

	//Set parameters
	e.SetLabels(m.GetLabels())
	e.Spec.Template.SetLabels(m.Spec.Template.GetLabels())
	e.SetAnnotations(m.GetAnnotations())
	e.Spec.Template.SetAnnotations(m.Spec.Template.GetAnnotations())
	e.Spec.Template.Spec.SecurityContext = m.Spec.Template.Spec.SecurityContext
	e.Spec.Template.Spec.Containers = m.Spec.Template.Spec.Containers
	e.Spec.Template.Spec.ServiceAccountName = m.Spec.Template.Spec.ServiceAccountName
	e.Spec.Template.Spec.NodeSelector = m.Spec.Template.Spec.NodeSelector
	e.Spec.Template.Spec.Affinity = m.Spec.Template.Spec.Affinity
	e.Spec.Replicas = m.Spec.Replicas
	e.Spec.Template.Spec.PriorityClassName = m.Spec.Template.Spec.PriorityClassName

	if err = r.UpdateResource(e); err != nil {
		return err
	}
	return nil
}

func (r *PushgatewayReconciler) handlePVC(cr *v1alpha1.PlatformMonitoring) error {
	m, err := pushgatewayPVC(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating PVC manifest")
		return err
	}
	e := &corev1.PersistentVolumeClaim{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err == nil {
		r.Log.Info("PVC with name pushgateway is already exist. Skip reconciling")
	} else {
		if errors.IsNotFound(err) {
			if err = r.CreateResource(cr, m); err != nil {
				return err
			}
			return nil
		}
		return err
	}
	return nil
}

func (r *PushgatewayReconciler) handleService(cr *v1alpha1.PlatformMonitoring) error {
	m, err := pushgatewayService(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Service manifest")
		return err
	}

	// Set labels
	m.Labels["name"] = utils.TruncLabel(m.GetName())
	m.Labels["app.kubernetes.io/name"] = utils.TruncLabel(m.GetName())
	m.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(m.GetName(), m.GetNamespace())
	m.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Pushgateway.Image)

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
	e.SetLabels(m.GetLabels())
	e.Spec.Ports = m.Spec.Ports
	e.Spec.Selector = m.Spec.Selector

	if err = r.UpdateResource(e); err != nil {
		return err
	}
	return nil
}

func (r *PushgatewayReconciler) handleIngressV1beta1(cr *v1alpha1.PlatformMonitoring) error {
	m, err := pushgatewayIngressV1beta1(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Ingress manifest")
		return err
	}
	e := &v1beta1.Ingress{ObjectMeta: m.ObjectMeta}
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
	e.SetLabels(m.GetLabels())
	e.SetAnnotations(m.GetAnnotations())
	e.Spec.Rules = m.Spec.Rules
	e.Spec.TLS = m.Spec.TLS

	if err = r.UpdateResource(e); err != nil {
		return err
	}
	return nil
}

func (r *PushgatewayReconciler) handleIngressV1(cr *v1alpha1.PlatformMonitoring) error {
	m, err := pushgatewayIngressV1(cr)
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
	e.SetLabels(m.GetLabels())
	e.SetAnnotations(m.GetAnnotations())
	e.Spec.Rules = m.Spec.Rules
	e.Spec.TLS = m.Spec.TLS

	if err = r.UpdateResource(e); err != nil {
		return err
	}
	return nil
}

func (r *PushgatewayReconciler) handleServiceMonitor(cr *v1alpha1.PlatformMonitoring) error {
	m, err := pushgatewayServiceMonitor(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ServiceMonitor manifest")
		return err
	}

	// Set labels
	m.Labels["name"] = utils.TruncLabel(m.GetName())
	m.Labels["app.kubernetes.io/name"] = utils.TruncLabel(m.GetName())
	m.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(m.GetName(), m.GetNamespace())
	m.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Pushgateway.Image)

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
	e.SetLabels(m.GetLabels())
	e.Spec.JobLabel = m.Spec.JobLabel
	e.Spec.Endpoints = m.Spec.Endpoints
	e.Spec.NamespaceSelector = m.Spec.NamespaceSelector
	e.Spec.Selector = m.Spec.Selector

	if err = r.UpdateResource(e); err != nil {
		return err
	}
	return nil
}

func (r *PushgatewayReconciler) deleteDeployment(cr *v1alpha1.PlatformMonitoring) error {
	m, err := pushgatewayDeployment(cr)
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

func (r *PushgatewayReconciler) deleteService(cr *v1alpha1.PlatformMonitoring) error {
	m, err := pushgatewayService(cr)
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

func (r *PushgatewayReconciler) deleteIngressV1beta1(cr *v1alpha1.PlatformMonitoring) error {
	m, err := pushgatewayIngressV1beta1(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Ingress manifest")
		return err
	}
	e := &v1beta1.Ingress{ObjectMeta: m.ObjectMeta}
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

func (r *PushgatewayReconciler) deleteIngressV1(cr *v1alpha1.PlatformMonitoring) error {
	m, err := pushgatewayIngressV1(cr)
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

func (r *PushgatewayReconciler) deleteServiceMonitor(cr *v1alpha1.PlatformMonitoring) error {
	m, err := pushgatewayServiceMonitor(cr)
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
