package prometheus_operator

import (
	v1beta1 "github.com/Netcracker/qubership-monitoring-operator/api"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PrometheusOperatorReconciler struct {
	*utils.ComponentReconciler
}

func NewPrometheusOperatorReconciler(c client.Client, s *runtime.Scheme) *PrometheusOperatorReconciler {
	return &PrometheusOperatorReconciler{
		ComponentReconciler: &utils.ComponentReconciler{
			Client: c,
			Scheme: s,
			Log:    utils.Logger("prometheusoperator_reconciler"),
		},
	}
}

// Run reconciles prometheus-operator.
// Creates new deployment, service, service account, cluster role and cluster role binding if its don't exists.
// Updates deployment and service in case of any changes.
// Returns true if need to requeue, false otherwise.
func (r *PrometheusOperatorReconciler) Run(cr *v1beta1.PlatformMonitoring) error {
	r.Log.Info("Reconciling component")

	if cr.Spec.Prometheus != nil && cr.Spec.Prometheus.IsInstall() {
		if !cr.Spec.Prometheus.Operator.Paused {
			if err := r.handleRole(cr); err != nil {
				return err
			}

			if !cr.Spec.Prometheus.Operator.Paused {
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
				if err := r.handleRoleBinding(cr); err != nil {
					return err
				}
				if err := r.handleService(cr); err != nil {
					return err
				}
				if err := r.handleDeployment(cr); err != nil {
					return err
				}
				// Reconcile Pod Monitor
				if cr.Spec.Prometheus.Operator.PodMonitor != nil && cr.Spec.Prometheus.Operator.PodMonitor.IsInstall() {
					if err := r.handlePodMonitor(cr); err != nil {
						return err
					}
				}
				r.Log.Info("Component reconciled")
			} else {
				r.Log.Info("Reconciling paused")
				r.Log.Info("Component NOT reconciled")
			}
		}
	} else {
		r.Log.Info("Uninstalling component if exists")
		r.uninstall(cr)
		r.Log.Info("Component reconciled")
	}

	return nil
}

// uninstall deletes all resources related to the component
func (r *PrometheusOperatorReconciler) uninstall(cr *v1beta1.PlatformMonitoring) {
	if err := r.deleteServiceAccount(cr); err != nil {
		r.Log.Error(err, "Can not delete ServiceAccount")
	}
	if utils.PrivilegedRights {
		if err := r.deleteClusterRole(cr); err != nil {
			r.Log.Error(err, "Can not delete ClusterRole")
		}
		if err := r.deleteClusterRoleBinding(cr); err != nil {
			r.Log.Error(err, "Can not delete ClusterRoleBinding")
		}
	}
	if err := r.deleteRole(cr); err != nil {
		r.Log.Error(err, "Can not delete Role")
	}
	if err := r.deleteRoleBinding(cr); err != nil {
		r.Log.Error(err, "Can not delete RoleBinding")
	}
	if err := r.deleteService(cr); err != nil {
		r.Log.Error(err, "Can not delete Service")
	}
	if err := r.deletePrometheusOperatorDeployment(cr); err != nil {
		r.Log.Error(err, "Can not delete Deployment")
	}
	if err := r.deletePodMonitor(cr); err != nil {
		r.Log.Error(err, "Can not delete PodMonitor")
	}
}
