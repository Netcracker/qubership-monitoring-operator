package vmalert

import (
	v1alpha1 "github.com/Netcracker/qubership-monitoring-operator/api/v1alpha1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	vmetricsv1b1 "github.com/VictoriaMetrics/operator/api/operator/v1beta1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *VmAlertReconciler) handleServiceAccount(cr *v1alpha1.PlatformMonitoring) error {
	m, err := vmAlertServiceAccount(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ServiceAccount manifest")
		return err
	}

	// Set labels
	m.Labels["name"] = utils.TruncLabel(m.GetName())
	m.Labels["app.kubernetes.io/name"] = utils.TruncLabel(m.GetName())
	m.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(m.GetName(), m.GetNamespace())
	m.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Victoriametrics.VmAlert.Image)

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
	e.SetLabels(m.GetLabels())

	if err = r.UpdateResource(e); err != nil {
		return err
	}
	return nil
}
func (r *VmAlertReconciler) handleClusterRole(cr *v1alpha1.PlatformMonitoring) error {
	m, err := vmAlertClusterRole(cr, r.hasPodSecurityPolicyAPI(), r.hasSecurityContextConstraintsAPI())
	if err != nil {
		r.Log.Error(err, "Failed creating ClusterRole manifest")
		return err
	}

	// Set labels
	m.Labels["name"] = utils.TruncLabel(m.GetName())
	m.Labels["app.kubernetes.io/name"] = utils.TruncLabel(m.GetName())
	m.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(m.GetName(), m.GetNamespace())
	m.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Victoriametrics.VmAlert.Image)

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
	e.SetLabels(m.GetLabels())
	e.SetName(m.GetName())
	e.Rules = m.Rules

	if err = r.UpdateResource(e); err != nil {
		return err
	}
	return nil
}

func (r *VmAlertReconciler) handleClusterRoleBinding(cr *v1alpha1.PlatformMonitoring) error {
	m, err := vmAlertClusterRoleBinding(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ClusterRoleBinding manifest")
		return err
	}

	// Set labels
	m.Labels["name"] = utils.TruncLabel(m.GetName())
	m.Labels["app.kubernetes.io/name"] = utils.TruncLabel(m.GetName())
	m.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(m.GetName(), m.GetNamespace())
	m.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Victoriametrics.VmAlert.Image)

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
	e.SetLabels(m.GetLabels())

	if err = r.UpdateResource(e); err != nil {
		return err
	}
	return nil
}

func (r *VmAlertReconciler) handleVmAlert(cr *v1alpha1.PlatformMonitoring) error {
	m, err := vmAlert(r, cr)
	if err != nil {
		r.Log.Error(err, "Failed creating vmalert manifest")
		return err
	}
	e := &vmetricsv1b1.VMAlert{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			e = &vmetricsv1b1.VMAlert{ObjectMeta: metav1.ObjectMeta{
				Name:      utils.VmAlertComponentName,
				Namespace: cr.GetNamespace(),
			}}
			if err = r.GetResource(e); err == nil {
				if err = r.DeleteResource(e); err != nil {
					return err
				}
			}
			if err = r.CreateResource(cr, m); err != nil {
				return err
			}
			return nil
		}
		return err
	}

	//Set parameters
	e.SetLabels(m.GetLabels())
	e.Spec = m.Spec

	if err = r.UpdateResource(e); err != nil {
		return err
	}
	return nil
}

// func (r *VmAlertReconciler) handleIngressV1beta1(cr *v1alpha1.PlatformMonitoring) error {
// 	m, err := vmAlertIngressV1beta1(cr)
// 	if err != nil {
// 		r.Log.Error(err, "Failed creating Ingress manifest")
// 		return err
// 	}
// 	e := &v1beta1.Ingress{ObjectMeta: m.ObjectMeta}
// 	if err = r.GetResource(e); err != nil {
// 		if errors.IsNotFound(err) {
// 			e = &v1beta1.Ingress{ObjectMeta: metav1.ObjectMeta{
// 				Name:      cr.GetNamespace() + "-" + utils.VmAlertComponentName,
// 				Namespace: cr.GetNamespace(),
// 			}}
// 			if err = r.GetResource(e); err == nil {
// 				if err = r.DeleteResource(e); err != nil {
// 					return err
// 				}
// 			}
// 			if err = r.CreateResource(cr, m); err != nil {
// 				return err
// 			}
// 			return nil
// 		}
// 		return err
// 	}

// 	//Set parameters
// 	e.SetLabels(m.GetLabels())
// 	e.SetAnnotations(m.GetAnnotations())
// 	e.Spec.Rules = m.Spec.Rules
// 	e.Spec.TLS = m.Spec.TLS

// 	if err = r.UpdateResource(e); err != nil {
// 		return err
// 	}
// 	return nil
// }

func (r *VmAlertReconciler) handleIngressV1(cr *v1alpha1.PlatformMonitoring) error {
	m, err := vmAlertIngressV1(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Ingress manifest")
		return err
	}
	e := &networkingv1.Ingress{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			e = &networkingv1.Ingress{ObjectMeta: metav1.ObjectMeta{
				Name:      cr.GetNamespace() + "-" + utils.VmAlertComponentName,
				Namespace: cr.GetNamespace(),
			}}
			if err = r.GetResource(e); err == nil {
				if err = r.DeleteResource(e); err != nil {
					return err
				}
			}
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

func (r *VmAlertReconciler) deleteServiceAccount(cr *v1alpha1.PlatformMonitoring) error {
	m, err := vmAlertServiceAccount(cr)
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

func (r *VmAlertReconciler) deleteClusterRole(cr *v1alpha1.PlatformMonitoring) error {
	m, err := vmAlertClusterRole(cr, r.hasPodSecurityPolicyAPI(), r.hasSecurityContextConstraintsAPI())
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

func (r *VmAlertReconciler) deleteClusterRoleBinding(cr *v1alpha1.PlatformMonitoring) error {
	m, err := vmAlertClusterRoleBinding(cr)
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

func (r *VmAlertReconciler) deleteVmAlert(cr *v1alpha1.PlatformMonitoring) error {
	m, err := vmAlert(r, cr)
	if err != nil {
		r.Log.Error(err, "Failed creating vmalert manifest")
		return err
	}
	e := &vmetricsv1b1.VMAlert{ObjectMeta: m.ObjectMeta}
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

// func (r *VmAlertReconciler) deleteIngressV1beta1(cr *v1alpha1.PlatformMonitoring) error {
// 	m, err := vmAlertIngressV1beta1(cr)
// 	if err != nil {
// 		r.Log.Error(err, "Failed creating Ingress manifest")
// 		return err
// 	}
// 	e := &v1beta1.Ingress{ObjectMeta: m.ObjectMeta}
// 	if err = r.GetResource(e); err != nil {
// 		if errors.IsNotFound(err) {
// 			return nil
// 		}
// 		return err
// 	}
// 	if err = r.DeleteResource(e); err != nil {
// 		return err
// 	}
// 	return nil
// }

func (r *VmAlertReconciler) deleteIngressV1(cr *v1alpha1.PlatformMonitoring) error {
	m, err := vmAlertIngressV1(cr)
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
