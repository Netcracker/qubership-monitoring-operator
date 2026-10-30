package vmalert

import (
	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	vmetricsv1b1 "github.com/VictoriaMetrics/operator/api/operator/v1beta1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *VmAlertReconciler) handleServiceAccount(cr *monv1.PlatformMonitoring) error {
	m, err := vmAlertServiceAccount(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ServiceAccount manifest")
		return err
	}

	e := &corev1.ServiceAccount{ObjectMeta: m.ObjectMeta}
	return utils.CreateOrUpdate(r.ComponentReconciler, cr, m, e, utils.SyncLabels)
}
func (r *VmAlertReconciler) handleClusterRole(cr *monv1.PlatformMonitoring) error {
	m, err := vmAlertClusterRole(cr, r.hasPodSecurityPolicyAPI(), r.hasSecurityContextConstraintsAPI())
	if err != nil {
		r.Log.Error(err, "Failed creating ClusterRole manifest")
		return err
	}

	e := &rbacv1.ClusterRole{ObjectMeta: m.ObjectMeta}
	return utils.CreateOrUpdate(r.ComponentReconciler, cr, m, e, utils.SyncClusterRole)
}

func (r *VmAlertReconciler) handleClusterRoleBinding(cr *monv1.PlatformMonitoring) error {
	m, err := vmAlertClusterRoleBinding(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ClusterRoleBinding manifest")
		return err
	}

	e := &rbacv1.ClusterRoleBinding{ObjectMeta: m.ObjectMeta}
	return utils.CreateOrUpdate(r.ComponentReconciler, cr, m, e, utils.SyncLabels)
}

func (r *VmAlertReconciler) handleVmAlert(cr *monv1.PlatformMonitoring) error {
	m, err := vmAlert(r, cr)
	if err != nil {
		r.Log.Error(err, "Failed creating vmalert manifest")
		return err
	}
	e := &vmetricsv1b1.VMAlert{ObjectMeta: m.ObjectMeta}
	legacy := &vmetricsv1b1.VMAlert{ObjectMeta: metav1.ObjectMeta{
		Name:      utils.VmAlertComponentName,
		Namespace: cr.GetNamespace(),
	}}
	return utils.CreateOrUpdateReplacingLegacy(r.ComponentReconciler, cr, m, e, legacy, func(existing, desired *vmetricsv1b1.VMAlert) bool {
		return utils.AnyChanged(
			utils.SetLabelsIfChanged(existing, desired.GetLabels()),
			utils.SetIfChanged(&existing.Spec, desired.Spec),
		)
	})
}

func (r *VmAlertReconciler) handleIngressV1(cr *monv1.PlatformMonitoring) error {
	m, err := vmAlertIngressV1(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Ingress manifest")
		return err
	}
	e := &networkingv1.Ingress{ObjectMeta: m.ObjectMeta}
	legacy := &networkingv1.Ingress{ObjectMeta: metav1.ObjectMeta{
		Name:      cr.GetNamespace() + "-" + utils.VmAlertComponentName,
		Namespace: cr.GetNamespace(),
	}}
	return utils.CreateOrUpdateReplacingLegacy(r.ComponentReconciler, cr, m, e, legacy, utils.SyncIngress)
}

func (r *VmAlertReconciler) deleteServiceAccount(cr *monv1.PlatformMonitoring) error {
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

func (r *VmAlertReconciler) deleteClusterRole(cr *monv1.PlatformMonitoring) error {
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

func (r *VmAlertReconciler) deleteClusterRoleBinding(cr *monv1.PlatformMonitoring) error {
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

func (r *VmAlertReconciler) deleteVmAlert(cr *monv1.PlatformMonitoring) error {
	// Deletion targets are addressed by their stable name and namespace so that uninstall does not
	// depend on platform discovery or on validation of the desired state.
	e := &vmetricsv1b1.VMAlert{ObjectMeta: metav1.ObjectMeta{Name: utils.ManagedCustomResourceName, Namespace: cr.GetNamespace()}}
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

func (r *VmAlertReconciler) deleteIngressV1(cr *monv1.PlatformMonitoring) error {
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
