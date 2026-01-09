package nodeexporter

import (
	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	secv1 "github.com/openshift/api/security/v1"
	pspApi "k8s.io/api/policy/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/discovery"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type NodeExporterReconciler struct {
	*utils.ComponentReconciler
}

func NewNodeExporterReconciler(c client.Client, s *runtime.Scheme, dc discovery.DiscoveryInterface) *NodeExporterReconciler {
	return &NodeExporterReconciler{
		ComponentReconciler: &utils.ComponentReconciler{
			Client: c,
			Scheme: s,
			Dc:     dc,
			Log:    utils.Logger("nodeexporter_reconciler"),
		},
	}
}

// Run reconciliation for node-exporter configuration.
// Creates new daemonset, service, service account, cluster role and cluster role binding if its don't exists.
// Updates deployment and service in case of any changes.
// Returns true if need to requeue, false otherwise.
func (r *NodeExporterReconciler) Run(cr *monv1.PlatformMonitoring) error {
	r.Log.Info("Reconciling component")

	if cr.Spec.NodeExporter != nil && cr.Spec.NodeExporter.IsInstall() {
		if !cr.Spec.NodeExporter.Paused {
			if err := r.handleServiceAccount(cr); err != nil {
				return err
			}

			// Reconcile ClusterRole and ClusterRoleBinding only if privileged mode used
			if utils.PrivilegedRights {
				if err := r.handleClusterRole(cr); err != nil {
					return err
				}
				if err := r.handleClusterRoleBinding(cr); err != nil {
					return err
				}
			} else {
				r.Log.Info("Skip ClusterRole and ClusterRoleBinding resources reconciliation because privilegedRights=false")
			}
			if err := r.handleService(cr); err != nil {
				return err
			}
			if err := r.handleDaemonSet(cr); err != nil {
				return err
			}

			// Reconcile ServiceMonitor if necessary
			if cr.Spec.NodeExporter.ServiceMonitor != nil && cr.Spec.NodeExporter.ServiceMonitor.IsInstall() {
				if err := r.handleServiceMonitor(cr); err != nil {
					return err
				}
			} else {
				r.Log.Info("Uninstalling ServiceMonitor")
				if err := r.deleteServiceMonitor(cr); err != nil {
					r.Log.Error(err, "Can not delete ServiceMonitor")
				}
			}
			r.Log.Info("Component reconciled")
		} else {
			r.Log.Info("Reconciling paused")
			r.Log.Info("component NOT reconciled")
		}
	} else {
		r.Log.Info("Uninstalling component if exists")
		r.uninstall(cr)
		r.Log.Info("Component reconciled")
	}
	return nil
}

// uninstall deletes all resources related to the component
func (r *NodeExporterReconciler) uninstall(cr *monv1.PlatformMonitoring) {
	if utils.PrivilegedRights {
		if err := r.deleteClusterRole(cr); err != nil {
			r.Log.Error(err, "Can not delete ClusterRole")
		}
		if err := r.deleteClusterRoleBinding(cr); err != nil {
			r.Log.Error(err, "Can not delete ClusterRoleBinding")
		}
	}
	if err := r.deleteServiceAccount(cr); err != nil {
		r.Log.Error(err, "Can not delete ServiceAccount")
	}
	if err := r.deleteDaemonSet(cr); err != nil {
		r.Log.Error(err, "Can not delete DaemonSet")
	}
	if err := r.deleteService(cr); err != nil {
		r.Log.Error(err, "Can not delete Service")
	}
	if err := r.deleteServiceMonitor(cr); err != nil {
		r.Log.Error(err, "Can not delete ServiceMonitor")
	}

}

// hasPodSecurityPolicyAPI checks that the cluster API has policy.v1beta.PodSecurityPolicy API.
func (r *NodeExporterReconciler) hasPodSecurityPolicyAPI() bool {
	return r.HasApi(pspApi.SchemeGroupVersion, "PodSecurityPolicy")
}

// hasSecurityContextConstraintsAPI checks that the cluster API has security.openshift.io.v1.SecurityContextConstraints API.
func (r *NodeExporterReconciler) hasSecurityContextConstraintsAPI() bool {
	return r.HasApi(secv1.GroupVersion, "SecurityContextConstraints")
}
