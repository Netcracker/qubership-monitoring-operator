package nodeexporter

import (
	"fmt"

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

func (r *NodeExporterReconciler) handleClusterRole(cr *monv1.PlatformMonitoring) error {
	m, err := nodeExporterClusterRole(cr, r.hasPodSecurityPolicyAPI(), r.hasSecurityContextConstraintsAPI())
	if err != nil {
		r.Log.Error(err, "Failed creating ClusterRole manifest")
		return err
	}

	in := utils.BaseOnlyLabelInput(m.GetName(), utils.NodeExporterComponentName)
	utils.SetLabelsForResource(m, in, nil)

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

func (r *NodeExporterReconciler) handleClusterRoleBinding(cr *monv1.PlatformMonitoring) error {
	m, err := nodeExporterClusterRoleBinding(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ClusterRoleBinding manifest")
		return err
	}

	in := utils.BaseOnlyLabelInput(m.GetName(), utils.NodeExporterComponentName)
	utils.SetLabelsForResource(m, in, nil)

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
	//Set parameters
	changed := utils.SetLabelsIfChanged(e, m.GetLabels())
	changed = applyNodeExporterSCCPolicy(e, m) || changed
	if !changed {
		return nil
	}

	if err = r.UpdateResource(e); err != nil {
		return err
	}
	return nil
}

func isNodeExporterSCCOwnedBy(scc *secv1.SecurityContextConstraints, cr *monv1.PlatformMonitoring) bool {
	annotations := scc.GetAnnotations()
	return annotations[nodeExporterSCCOwnerNameAnnotation] == cr.GetName() &&
		annotations[nodeExporterSCCOwnerNamespaceAnnotation] == cr.GetNamespace()
}

func applyNodeExporterSCCPolicy(existing, desired *secv1.SecurityContextConstraints) bool {
	changed := false
	changed = utils.SetIfChanged(&existing.AllowPrivilegedContainer, desired.AllowPrivilegedContainer) || changed
	changed = utils.SetIfChanged(&existing.DefaultAddCapabilities, desired.DefaultAddCapabilities) || changed
	changed = utils.SetIfChanged(&existing.RequiredDropCapabilities, desired.RequiredDropCapabilities) || changed
	changed = utils.SetIfChanged(&existing.AllowedCapabilities, desired.AllowedCapabilities) || changed
	changed = utils.SetIfChanged(&existing.AllowHostDirVolumePlugin, desired.AllowHostDirVolumePlugin) || changed
	changed = utils.SetIfChanged(&existing.Volumes, desired.Volumes) || changed
	changed = utils.SetIfChanged(&existing.AllowedFlexVolumes, desired.AllowedFlexVolumes) || changed
	changed = utils.SetIfChanged(&existing.AllowHostNetwork, desired.AllowHostNetwork) || changed
	changed = utils.SetIfChanged(&existing.AllowHostPorts, desired.AllowHostPorts) || changed
	changed = utils.SetIfChanged(&existing.AllowHostPID, desired.AllowHostPID) || changed
	changed = utils.SetIfChanged(&existing.AllowHostIPC, desired.AllowHostIPC) || changed
	changed = utils.SetIfChanged(&existing.DefaultAllowPrivilegeEscalation, desired.DefaultAllowPrivilegeEscalation) || changed
	changed = utils.SetIfChanged(&existing.AllowPrivilegeEscalation, desired.AllowPrivilegeEscalation) || changed
	changed = utils.SetIfChanged(&existing.SELinuxContext, desired.SELinuxContext) || changed
	changed = utils.SetIfChanged(&existing.RunAsUser, desired.RunAsUser) || changed
	changed = utils.SetIfChanged(&existing.SupplementalGroups, desired.SupplementalGroups) || changed
	changed = utils.SetIfChanged(&existing.FSGroup, desired.FSGroup) || changed
	changed = utils.SetIfChanged(&existing.ReadOnlyRootFilesystem, desired.ReadOnlyRootFilesystem) || changed
	changed = utils.SetIfChanged(&existing.SeccompProfiles, desired.SeccompProfiles) || changed
	changed = utils.SetIfChanged(&existing.AllowedUnsafeSysctls, desired.AllowedUnsafeSysctls) || changed
	changed = utils.SetIfChanged(&existing.ForbiddenSysctls, desired.ForbiddenSysctls) || changed
	return changed
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
	m, err := nodeExporterDaemonSet(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating DaemonSet manifest")
		return err
	}
	e := &appsv1.DaemonSet{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			if err = r.CreateResource(cr, m); err != nil {
				return err
			}
			return nil
		}
		return err
	}
	if utils.WorkloadNeedsSelectorReplace(e, m) {
		if err := r.DeleteResource(e); err != nil {
			return err
		}
		return r.CreateResource(cr, m)
	}
	//Set parameters
	changed := utils.SetLabelsIfChanged(e, m.GetLabels())
	changed = utils.SetLabelsIfChanged(&e.Spec.Template, m.Spec.Template.GetLabels()) || changed
	changed = utils.SetIfChanged(&e.Spec.Template.Spec, m.Spec.Template.Spec) || changed
	if !changed {
		return nil
	}

	if err = r.UpdateResource(e); err != nil {
		return err
	}
	return nil
}

func (r *NodeExporterReconciler) handleService(cr *monv1.PlatformMonitoring) error {
	m, err := nodeExporterService(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Service manifest")
		return err
	}

	m.Spec.Selector = map[string]string{"app.kubernetes.io/name": utils.TruncLabel(utils.NodeExporterComponentName)}

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

func (r *NodeExporterReconciler) handleServiceMonitor(cr *monv1.PlatformMonitoring) error {
	m, err := nodeExporterServiceMonitor(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ServiceMonitor manifest")
		return err
	}

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
	changed := utils.SetLabelsIfChanged(e, m.GetLabels())
	changed = utils.SetIfChanged(&e.Spec.JobLabel, m.Spec.JobLabel) || changed
	changed = utils.SetIfChanged(&e.Spec.Endpoints, m.Spec.Endpoints) || changed
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
	m, err := nodeExporterDaemonSet(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating DaemonSet manifest")
		return err
	}
	e := &appsv1.DaemonSet{ObjectMeta: m.ObjectMeta}
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
