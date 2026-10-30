package vmcluster

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

func (r *VmClusterReconciler) handleServiceAccount(cr *monv1.PlatformMonitoring) error {
	m, err := vmClusterServiceAccount(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ServiceAccount manifest")
		return err
	}

	e := &corev1.ServiceAccount{ObjectMeta: m.ObjectMeta}
	return utils.CreateOrUpdate(r.ComponentReconciler, cr, m, e, utils.SyncLabels)
}
func (r *VmClusterReconciler) handleClusterRole(cr *monv1.PlatformMonitoring) error {
	m, err := vmClusterClusterRole(cr, r.hasPodSecurityPolicyAPI(), r.hasSecurityContextConstraintsAPI())
	if err != nil {
		r.Log.Error(err, "Failed creating ClusterRole manifest")
		return err
	}

	e := &rbacv1.ClusterRole{ObjectMeta: m.ObjectMeta}
	return utils.CreateOrUpdate(r.ComponentReconciler, cr, m, e, utils.SyncClusterRole)
}

func (r *VmClusterReconciler) handleClusterRoleBinding(cr *monv1.PlatformMonitoring) error {
	m, err := vmClusterClusterRoleBinding(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ClusterRoleBinding manifest")
		return err
	}

	e := &rbacv1.ClusterRoleBinding{ObjectMeta: m.ObjectMeta}
	return utils.CreateOrUpdate(r.ComponentReconciler, cr, m, e, utils.SyncLabels)
}

func (r *VmClusterReconciler) handleIngressV1(cr *monv1.PlatformMonitoring) error {
	m, err := vmSelectIngressV1(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Ingress manifest")
		return err
	}
	e := &networkingv1.Ingress{ObjectMeta: m.ObjectMeta}
	legacy := &networkingv1.Ingress{ObjectMeta: metav1.ObjectMeta{
		Name:      cr.GetNamespace() + "-" + utils.VmSelectComponentName,
		Namespace: cr.GetNamespace(),
	}}
	return utils.CreateOrUpdateReplacingLegacy(r.ComponentReconciler, cr, m, e, legacy, utils.SyncIngress)
}

func (r *VmClusterReconciler) handleVmCluster(cr *monv1.PlatformMonitoring) error {
	m, err := vmCluster(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating vmcluster manifest")
		return err
	}
	e := &vmetricsv1b1.VMCluster{ObjectMeta: m.ObjectMeta}
	legacy := &vmetricsv1b1.VMCluster{ObjectMeta: metav1.ObjectMeta{
		Name:      utils.VmClusterComponentName,
		Namespace: cr.GetNamespace(),
	}}
	return utils.CreateOrUpdateReplacingLegacy(r.ComponentReconciler, cr, m, e, legacy, func(existing, desired *vmetricsv1b1.VMCluster) bool {
		return utils.AnyChanged(
			utils.SetLabelsIfChanged(existing, desired.GetLabels()),
			utils.SetIfChanged(&existing.Spec, desired.Spec),
		)
	})
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
