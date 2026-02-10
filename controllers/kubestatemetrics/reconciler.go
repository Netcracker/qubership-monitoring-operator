package kubestatemetrics

import (
	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/discovery"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type KubeStateMetricsReconciler struct {
	*utils.ComponentReconciler
}

func NewKubeStateMetricsReconciler(c client.Client, s *runtime.Scheme, dc discovery.DiscoveryInterface) *KubeStateMetricsReconciler {
	return &KubeStateMetricsReconciler{
		ComponentReconciler: &utils.ComponentReconciler{
			Client: c,
			Scheme: s,
			Dc:     dc,
			Log:    utils.Logger("kubestatemetrics_reconciler"),
		},
	}
}

// Run reconciles kube-state-metrics.
// Creates new deployment, service, service account, cluster role and cluster role binding if its don't exists.
// Updates deployment and service in case of any changes.
// Returns true if need to requeue, false otherwise.
func (r *KubeStateMetricsReconciler) Run(cr *monv1.PlatformMonitoring) error {
	r.Log.Info("reconciling component")

	if cr.Spec.KubeStateMetrics != nil && cr.Spec.KubeStateMetrics.IsInstall() {
		if !cr.Spec.KubeStateMetrics.Paused {
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
				r.Log.Info("skip ClusterRole and ClusterRoleBinding resources reconciliation because privilegedRights=false")
			}
			if err := r.handleService(cr); err != nil {
				return err
			}
			if err := r.handleDeployment(cr); err != nil {
				return err
			}

			// Reconcile ServiceMonitor if necessary
			if cr.Spec.KubeStateMetrics.ServiceMonitor != nil && cr.Spec.KubeStateMetrics.ServiceMonitor.IsInstall() {
				if err := r.handleServiceMonitor(cr); err != nil {
					return err
				}
			} else {
				r.Log.Info("uninstalling ServiceMonitor")
				if err := r.deleteServiceMonitor(cr); err != nil {
					r.Log.Error(err, "can not delete ServiceMonitor")
				}
			}
			r.Log.Info("component reconciled")
		} else {
			r.Log.Info("reconciling paused")
			r.Log.Info("component NOT reconciled")
		}
	} else {
		r.Log.Info("uninstalling component if exists")
		r.uninstall(cr)
		r.Log.Info("component reconciled")
	}
	return nil
}

// uninstall deletes all resources related to the component
func (r *KubeStateMetricsReconciler) uninstall(cr *monv1.PlatformMonitoring) {
	if utils.PrivilegedRights {
		if err := r.deleteClusterRole(cr); err != nil {
			r.Log.Error(err, "can not delete ClusterRole")
		}
		if err := r.deleteClusterRoleBinding(cr); err != nil {
			r.Log.Error(err, "can not delete ClusterRoleBinding")
		}
	}
	if err := r.deleteServiceAccount(cr); err != nil {
		r.Log.Error(err, "can not delete ServiceAccount")
	}
	if err := r.deleteDeployment(cr); err != nil {
		r.Log.Error(err, "can not delete Deployment")
	}
	if err := r.deleteService(cr); err != nil {
		r.Log.Error(err, "can not delete Service")
	}
	if err := r.deleteServiceMonitor(cr); err != nil {
		r.Log.Error(err, "can not delete ServiceMonitor")
	}

}
