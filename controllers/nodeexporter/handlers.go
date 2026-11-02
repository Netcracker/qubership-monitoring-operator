package nodeexporter

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	secv1 "github.com/openshift/api/security/v1"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
)

func (r *NodeExporterReconciler) handleServiceAccount(cr *monv1.PlatformMonitoring) error {
	m, err := nodeExporterServiceAccount(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ServiceAccount manifest")
		return err
	}

	e := &corev1.ServiceAccount{ObjectMeta: m.ObjectMeta}
	return utils.CreateOrUpdate(r.ComponentReconciler, cr, m, e, utils.SyncLabels)
}

func (r *NodeExporterReconciler) handleClusterRole(cr *monv1.PlatformMonitoring) error {
	m, err := nodeExporterClusterRole(cr, r.hasPodSecurityPolicyAPI(), r.hasSecurityContextConstraintsAPI())
	if err != nil {
		r.Log.Error(err, "Failed creating ClusterRole manifest")
		return err
	}

	in := utils.BaseOnlyLabelInput(m.GetName(), utils.NodeExporterComponentName)
	utils.SetLabelsForResource(m, in, nil)

	e := &rbacv1.ClusterRole{ObjectMeta: m.ObjectMeta}
	return utils.CreateOrUpdate(r.ComponentReconciler, cr, m, e, utils.SyncClusterRole)
}

func (r *NodeExporterReconciler) handleClusterRoleBinding(cr *monv1.PlatformMonitoring) error {
	m, err := nodeExporterClusterRoleBinding(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ClusterRoleBinding manifest")
		return err
	}

	in := utils.BaseOnlyLabelInput(m.GetName(), utils.NodeExporterComponentName)
	utils.SetLabelsForResource(m, in, nil)

	e := &rbacv1.ClusterRoleBinding{ObjectMeta: m.ObjectMeta}
	return utils.CreateOrUpdate(r.ComponentReconciler, cr, m, e, utils.SyncLabels)
}

func (r *NodeExporterReconciler) handleSecurityContextConstraints(cr *monv1.PlatformMonitoring) error {
	m, err := nodeExporterSecurityContextConstraints(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating SecurityContextConstraints manifest")
		return err
	}

	e := &secv1.SecurityContextConstraints{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			if err = r.CreateResource(cr, m, false); err != nil {
				return err
			}
			return nil
		}
		return err
	}
	if !isNodeExporterSCCOwnedBy(e, cr) {
		return fmt.Errorf(
			"cannot reconcile SecurityContextConstraints %q: it is not owned by PlatformMonitoring %s/%s",
			e.GetName(),
			cr.GetNamespace(),
			cr.GetName(),
		)
	}
	return r.UpdateIfChanged(e, utils.SyncSecurityContextConstraints(e, m))
}

func isNodeExporterSCCOwnedBy(scc *secv1.SecurityContextConstraints, cr *monv1.PlatformMonitoring) bool {
	annotations := scc.GetAnnotations()
	return annotations[nodeExporterSCCOwnerNameAnnotation] == cr.GetName() &&
		annotations[nodeExporterSCCOwnerNamespaceAnnotation] == cr.GetNamespace()
}

func (r *NodeExporterReconciler) deleteSecurityContextConstraints(cr *monv1.PlatformMonitoring) error {
	m, err := nodeExporterSecurityContextConstraints(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating SecurityContextConstraints manifest")
		return err
	}
	e := &secv1.SecurityContextConstraints{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if !isNodeExporterSCCOwnedBy(e, cr) {
		r.Log.Info(
			"Skip deleting SecurityContextConstraints because it is owned by another installation",
			"name", e.GetName(),
		)
		return nil
	}
	if err = r.DeleteResource(e); err != nil {
		return err
	}
	return nil
}

func (r *NodeExporterReconciler) handleDaemonSet(cr *monv1.PlatformMonitoring) error {
	isOpenShift, err := r.IsOpenShift()
	if err != nil {
		return err
	}
	m, err := nodeExporterDaemonSet(cr, isOpenShift)
	if err != nil {
		r.Log.Error(err, "Failed creating DaemonSet manifest")
		return err
	}
	e := &appsv1.DaemonSet{ObjectMeta: m.ObjectMeta}
	return utils.CreateOrUpdateWorkload(r.ComponentReconciler, cr, m, e, syncDaemonSet)
}

// syncDaemonSet copies the DaemonSet fields this controller manages.
func syncDaemonSet(existing, desired *appsv1.DaemonSet) bool {
	return utils.AnyChanged(
		utils.SetLabelsIfChanged(existing, desired.GetLabels()),
		utils.SetLabelsIfChanged(&existing.Spec.Template, desired.Spec.Template.GetLabels()),
		utils.SyncPodSpec(&existing.Spec.Template.Spec, &desired.Spec.Template.Spec),
		utils.SetIfChanged(&existing.Spec.Template.Spec.NodeSelector, desired.Spec.Template.Spec.NodeSelector),
		utils.SetIfChanged(&existing.Spec.Template.Spec.Tolerations, desired.Spec.Template.Spec.Tolerations),
		utils.SetIfChanged(&existing.Spec.Template.Spec.HostNetwork, desired.Spec.Template.Spec.HostNetwork),
		utils.SetIfChanged(&existing.Spec.Template.Spec.HostPID, desired.Spec.Template.Spec.HostPID),
	)
}

func (r *NodeExporterReconciler) handleService(cr *monv1.PlatformMonitoring) error {
	m, err := nodeExporterService(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Service manifest")
		return err
	}

	m.Spec.Selector = map[string]string{"app.kubernetes.io/name": utils.TruncLabel(utils.NodeExporterComponentName)}

	e := &corev1.Service{ObjectMeta: m.ObjectMeta}
	return utils.CreateOrUpdate(r.ComponentReconciler, cr, m, e, utils.SyncService)
}

func (r *NodeExporterReconciler) handleServiceMonitor(cr *monv1.PlatformMonitoring) error {
	m, err := nodeExporterServiceMonitor(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ServiceMonitor manifest")
		return err
	}

	e := &promv1.ServiceMonitor{ObjectMeta: m.ObjectMeta}
	return utils.CreateOrUpdate(r.ComponentReconciler, cr, m, e, utils.SyncServiceMonitor)
}

func (r *NodeExporterReconciler) deleteServiceAccount(cr *monv1.PlatformMonitoring) error {
	m, err := nodeExporterServiceAccount(cr)
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

func (r *NodeExporterReconciler) deleteClusterRole(cr *monv1.PlatformMonitoring) error {
	m, err := nodeExporterClusterRole(cr, r.hasPodSecurityPolicyAPI(), r.hasSecurityContextConstraintsAPI())
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

func (r *NodeExporterReconciler) deleteClusterRoleBinding(cr *monv1.PlatformMonitoring) error {
	m, err := nodeExporterClusterRoleBinding(cr)
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

func (r *NodeExporterReconciler) deleteDaemonSet(cr *monv1.PlatformMonitoring) error {
	// Deletion targets are addressed by their stable name and namespace so that uninstall does not
	// depend on platform discovery or on validation of the desired state.
	e := &appsv1.DaemonSet{ObjectMeta: metav1.ObjectMeta{Name: utils.NodeExporterComponentName, Namespace: cr.GetNamespace()}}
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

func (r *NodeExporterReconciler) deleteService(cr *monv1.PlatformMonitoring) error {
	m, err := nodeExporterService(cr)
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

func (r *NodeExporterReconciler) deleteServiceMonitor(cr *monv1.PlatformMonitoring) error {
	m, err := nodeExporterServiceMonitor(cr)
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
