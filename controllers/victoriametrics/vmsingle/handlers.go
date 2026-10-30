package vmsingle

import (
	"context"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	vmetricsv1b1 "github.com/VictoriaMetrics/operator/api/operator/v1beta1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *VmSingleReconciler) handleServiceAccount(cr *monv1.PlatformMonitoring) error {
	m, err := vmSingleServiceAccount(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ServiceAccount manifest")
		return err
	}

	e := &corev1.ServiceAccount{ObjectMeta: m.ObjectMeta}
	return utils.CreateOrUpdate(r.ComponentReconciler, cr, m, e, utils.SyncLabels)
}
func (r *VmSingleReconciler) handleClusterRole(cr *monv1.PlatformMonitoring) error {
	m, err := vmSingleClusterRole(cr, r.hasPodSecurityPolicyAPI(), r.hasSecurityContextConstraintsAPI())
	if err != nil {
		r.Log.Error(err, "Failed creating ClusterRole manifest")
		return err
	}

	e := &rbacv1.ClusterRole{ObjectMeta: m.ObjectMeta}
	return utils.CreateOrUpdate(r.ComponentReconciler, cr, m, e, utils.SyncClusterRole)
}

func (r *VmSingleReconciler) handleClusterRoleBinding(cr *monv1.PlatformMonitoring) error {
	m, err := vmSingleClusterRoleBinding(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ClusterRoleBinding manifest")
		return err
	}

	e := &rbacv1.ClusterRoleBinding{ObjectMeta: m.ObjectMeta}
	return utils.CreateOrUpdate(r.ComponentReconciler, cr, m, e, utils.SyncLabels)
}

func (r *VmSingleReconciler) handleVmSingle(cr *monv1.PlatformMonitoring) error {
	m, err := vmSingle(r, cr)
	if err != nil {
		r.Log.Error(err, "Failed creating vmsingle manifest")
		return err
	}
	e := &vmetricsv1b1.VMSingle{ObjectMeta: m.ObjectMeta}
	legacy := &vmetricsv1b1.VMSingle{ObjectMeta: metav1.ObjectMeta{
		Name:      utils.VmSingleComponentName,
		Namespace: cr.GetNamespace(),
	}}
	return utils.CreateOrUpdateReplacingLegacy(r.ComponentReconciler, cr, m, e, legacy, func(existing, desired *vmetricsv1b1.VMSingle) bool {
		return utils.AnyChanged(
			utils.SetLabelsIfChanged(existing, desired.GetLabels()),
			utils.SetIfChanged(&existing.Spec, desired.Spec),
		)
	})
}

func (r *VmSingleReconciler) handleIngressV1(cr *monv1.PlatformMonitoring) error {
	m, err := vmSingleIngressV1(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Ingress manifest")
		return err
	}
	e := &networkingv1.Ingress{ObjectMeta: m.ObjectMeta}
	legacy := &networkingv1.Ingress{ObjectMeta: metav1.ObjectMeta{
		Name:      cr.GetNamespace() + "-" + utils.VmSingleComponentName,
		Namespace: cr.GetNamespace(),
	}}
	return utils.CreateOrUpdateReplacingLegacy(r.ComponentReconciler, cr, m, e, legacy, utils.SyncIngress)
}

func (r *VmSingleReconciler) deleteServiceAccount(cr *monv1.PlatformMonitoring) error {
	m, err := vmSingleServiceAccount(cr)
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

func (r *VmSingleReconciler) deleteClusterRole(cr *monv1.PlatformMonitoring) error {
	m, err := vmSingleClusterRole(cr, r.hasPodSecurityPolicyAPI(), r.hasSecurityContextConstraintsAPI())
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

func (r *VmSingleReconciler) deleteClusterRoleBinding(cr *monv1.PlatformMonitoring) error {
	m, err := vmSingleClusterRoleBinding(cr)
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

func (r *VmSingleReconciler) deleteVmSingle(cr *monv1.PlatformMonitoring) error {
	// Deletion targets are addressed by their stable name and namespace so that uninstall does not
	// depend on platform discovery or on validation of the desired state.
	e := &vmetricsv1b1.VMSingle{ObjectMeta: metav1.ObjectMeta{Name: utils.ManagedCustomResourceName, Namespace: cr.GetNamespace()}}
	if err := r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if err := r.Client.Delete(
		context.TODO(),
		e,
		client.PropagationPolicy(metav1.DeletePropagationOrphan),
	); err != nil {
		return err
	}
	r.Log.Info("Successful deleting", utils.ResourceKey, e.GetObjectKind().GroupVersionKind().Kind)
	return nil
}

func (r *VmSingleReconciler) deleteIngressV1(cr *monv1.PlatformMonitoring) error {
	m, err := vmSingleIngressV1(cr)
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
