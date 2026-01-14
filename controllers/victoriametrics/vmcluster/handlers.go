package vmcluster

import (
	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	vmetricsv1b1 "github.com/VictoriaMetrics/operator/api/operator/v1beta1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/api/networking/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *VmClusterReconciler) handleServiceAccount(cr *monv1.PlatformMonitoring) error {
	m, err := vmClusterServiceAccount(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ServiceAccount manifest")
		return err
	}

	// Set labels
	m.Labels["name"] = utils.TruncLabel(m.GetName())
	m.Labels["app.kubernetes.io/name"] = utils.TruncLabel(m.GetName())
	m.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(m.GetName(), m.GetNamespace())

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
func (r *VmClusterReconciler) handleClusterRole(cr *monv1.PlatformMonitoring) error {
	m, err := vmClusterClusterRole(cr, r.hasPodSecurityPolicyAPI(), r.hasSecurityContextConstraintsAPI())
	if err != nil {
		r.Log.Error(err, "Failed creating ClusterRole manifest")
		return err
	}

	// Set labels
	m.Labels["name"] = utils.TruncLabel(m.GetName())
	m.Labels["app.kubernetes.io/name"] = utils.TruncLabel(m.GetName())
	m.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(m.GetName(), m.GetNamespace())

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

func (r *VmClusterReconciler) handleClusterRoleBinding(cr *monv1.PlatformMonitoring) error {
	m, err := vmClusterClusterRoleBinding(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ClusterRoleBinding manifest")
		return err
	}

	// Set labels
	m.Labels["name"] = utils.TruncLabel(m.GetName())
	m.Labels["app.kubernetes.io/name"] = utils.TruncLabel(m.GetName())
	m.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(m.GetName(), m.GetNamespace())

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

func (r *VmClusterReconciler) handleIngressV1beta1(cr *monv1.PlatformMonitoring) error {
	m, err := vmSelectIngressV1beta1(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Ingress manifest")
		return err
	}
	e := &v1beta1.Ingress{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			e = &v1beta1.Ingress{ObjectMeta: metav1.ObjectMeta{
				Name:      cr.GetNamespace() + "-" + utils.VmSelectComponentName,
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

func (r *VmClusterReconciler) handleIngressV1(cr *monv1.PlatformMonitoring) error {
	m, err := vmSelectIngressV1(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Ingress manifest")
		return err
	}
	e := &networkingv1.Ingress{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			e = &networkingv1.Ingress{ObjectMeta: metav1.ObjectMeta{
				Name:      cr.GetNamespace() + "-" + utils.VmSelectComponentName,
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

func (r *VmClusterReconciler) handleVmCluster(cr *monv1.PlatformMonitoring) error {
	m, err := vmCluster(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating vmcluster manifest")
		return err
	}
	e := &vmetricsv1b1.VMCluster{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			e = &vmetricsv1b1.VMCluster{ObjectMeta: metav1.ObjectMeta{
				Name:      utils.VmClusterComponentName,
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

func (r *VmClusterReconciler) deleteServiceAccount(cr *monv1.PlatformMonitoring) error {
	m, err := vmClusterServiceAccount(cr)
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

func (r *VmClusterReconciler) deleteClusterRole(cr *monv1.PlatformMonitoring) error {
	m, err := vmClusterClusterRole(cr, r.hasPodSecurityPolicyAPI(), r.hasSecurityContextConstraintsAPI())
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

func (r *VmClusterReconciler) deleteClusterRoleBinding(cr *monv1.PlatformMonitoring) error {
	m, err := vmClusterClusterRoleBinding(cr)
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

func (r *VmClusterReconciler) deleteVmCluster(cr *monv1.PlatformMonitoring) error {
	m, err := vmCluster(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating vmCluster manifest")
		return err
	}
	e := &vmetricsv1b1.VMCluster{ObjectMeta: m.ObjectMeta}
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

func (r *VmClusterReconciler) deleteIngressV1beta1(cr *monv1.PlatformMonitoring) error {
	m, err := vmSelectIngressV1beta1(cr)
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

func (r *VmClusterReconciler) deleteIngressV1(cr *monv1.PlatformMonitoring) error {
	m, err := vmSelectIngressV1(cr)
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
