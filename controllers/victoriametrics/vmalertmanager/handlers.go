package vmalertmanager

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

func (r *VmAlertManagerReconciler) handleServiceAccount(cr *monv1.PlatformMonitoring) error {
	m, err := vmAlertManagerServiceAccount(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ServiceAccount manifest")
		return err
	}

	e := &corev1.ServiceAccount{ObjectMeta: m.ObjectMeta}
	return utils.CreateOrUpdate(r.ComponentReconciler, cr, m, e, utils.SyncLabels)
}

func (r *VmAlertManagerReconciler) handleClusterRole(cr *monv1.PlatformMonitoring) error {
	m, err := vmAlertManagerClusterRole(cr, r.hasPodSecurityPolicyAPI(), r.hasSecurityContextConstraintsAPI())
	if err != nil {
		r.Log.Error(err, "Failed creating ClusterRole manifest")
		return err
	}

	e := &rbacv1.ClusterRole{ObjectMeta: m.ObjectMeta}
	return utils.CreateOrUpdate(r.ComponentReconciler, cr, m, e, utils.SyncClusterRole)
}

func (r *VmAlertManagerReconciler) handleClusterRoleBinding(cr *monv1.PlatformMonitoring) error {
	m, err := vmAlertManagerClusterRoleBinding(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ClusterRoleBinding manifest")
		return err
	}

	e := &rbacv1.ClusterRoleBinding{ObjectMeta: m.ObjectMeta}
	return utils.CreateOrUpdate(r.ComponentReconciler, cr, m, e, utils.SyncLabels)
}

func (r *VmAlertManagerReconciler) handleRole(cr *monv1.PlatformMonitoring) error {
	m, err := vmAlertManagerRole(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Role manifest")
		return err
	}

	e := &rbacv1.Role{ObjectMeta: m.ObjectMeta}
	return utils.CreateOrUpdate(r.ComponentReconciler, cr, m, e, utils.SyncRole)
}

func (r *VmAlertManagerReconciler) handleRoleBinding(cr *monv1.PlatformMonitoring) error {
	m, err := vmAlertManagerRoleBinding(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating RoleBinding manifest")
		return err
	}

	e := &rbacv1.RoleBinding{ObjectMeta: m.ObjectMeta}
	return utils.CreateOrUpdate(r.ComponentReconciler, cr, m, e, utils.SyncLabels)
}

func (r *VmAlertManagerReconciler) handleVmAlertManager(cr *monv1.PlatformMonitoring) error {
	m, err := vmAlertManager(r, cr)
	if err != nil {
		r.Log.Error(err, "Failed creating vmalertmanager manifest")
		return err
	}
	e := &vmetricsv1b1.VMAlertmanager{ObjectMeta: m.ObjectMeta}
	legacy := &vmetricsv1b1.VMAlertmanager{ObjectMeta: metav1.ObjectMeta{
		Name:      utils.VmAlertManagerComponentName,
		Namespace: cr.GetNamespace(),
	}}
	return utils.CreateOrUpdateReplacingLegacy(r.ComponentReconciler, cr, m, e, legacy, func(existing, desired *vmetricsv1b1.VMAlertmanager) bool {
		return utils.AnyChanged(
			utils.SetLabelsIfChanged(existing, desired.GetLabels()),
			utils.SetIfChanged(&existing.Spec, desired.Spec),
		)
	})
}

func (r *VmAlertManagerReconciler) handleSecret(cr *monv1.PlatformMonitoring) error {
	m, err := vmAlertmanagerSecret(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Secret manifest")
		return err
	}

	e := &corev1.Secret{ObjectMeta: m.ObjectMeta}
	legacy := &corev1.Secret{ObjectMeta: metav1.ObjectMeta{
		Name:      "vmalertmanager-vmalertmanager-config",
		Namespace: cr.GetNamespace(),
	}}
	return utils.CreateOrUpdateReplacingLegacy(r.ComponentReconciler, cr, m, e, legacy, utils.SyncLabels)
}

func (r *VmAlertManagerReconciler) handleIngressV1(cr *monv1.PlatformMonitoring) error {
	m, err := vmAlertManagerIngressV1(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Ingress manifest")
		return err
	}
	e := &networkingv1.Ingress{ObjectMeta: m.ObjectMeta}
	legacy := &networkingv1.Ingress{ObjectMeta: metav1.ObjectMeta{
		Name:      cr.GetNamespace() + "-" + utils.VmAlertManagerComponentName,
		Namespace: cr.GetNamespace(),
	}}
	return utils.CreateOrUpdateReplacingLegacy(r.ComponentReconciler, cr, m, e, legacy, utils.SyncIngress)
}

func (r *VmAlertManagerReconciler) deleteServiceAccount(cr *monv1.PlatformMonitoring) error {
	m, err := vmAlertManagerServiceAccount(cr)
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

func (r *VmAlertManagerReconciler) deleteClusterRole(cr *monv1.PlatformMonitoring) error {
	m, err := vmAlertManagerClusterRole(cr, r.hasPodSecurityPolicyAPI(), r.hasSecurityContextConstraintsAPI())
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

func (r *VmAlertManagerReconciler) deleteClusterRoleBinding(cr *monv1.PlatformMonitoring) error {
	m, err := vmAlertManagerClusterRoleBinding(cr)
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

func (r *VmAlertManagerReconciler) deleteVmAlertManager(cr *monv1.PlatformMonitoring) error {
	// Deletion targets are addressed by their stable name and namespace so that uninstall does not
	// depend on platform discovery or on validation of the desired state.
	e := &vmetricsv1b1.VMAlertmanager{ObjectMeta: metav1.ObjectMeta{Name: utils.ManagedCustomResourceName, Namespace: cr.GetNamespace()}}
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

func (r *VmAlertManagerReconciler) deleteSecret(cr *monv1.PlatformMonitoring) error {
	m, err := vmAlertmanagerSecret(cr)
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

func (r *VmAlertManagerReconciler) deleteIngressV1(cr *monv1.PlatformMonitoring) error {
	m, err := vmAlertManagerIngressV1(cr)
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
