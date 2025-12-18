package vmoperator

import (
	v1beta1 "github.com/Netcracker/qubership-monitoring-operator/api/v1beta1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	secv1 "github.com/openshift/api/security/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type VmOperatorReconciler struct {
	KubeClient kubernetes.Interface
	*utils.ComponentReconciler
}

func NewVmOperatorReconciler(c client.Client, s *runtime.Scheme, r *rest.Config, dc discovery.DiscoveryInterface) *VmOperatorReconciler {
	clientSet, _ := kubernetes.NewForConfig(r)
	return &VmOperatorReconciler{
		ComponentReconciler: &utils.ComponentReconciler{
			Client: c,
			Scheme: s,
			Dc:     dc,
			Log:    utils.Logger("vmoperator_reconciler"),
		},
		KubeClient: clientSet,
	}
}

// Run reconciles victoriametrics-operator.
// Creates new deployment, service, service account, cluster role and cluster role binding if they don't exist.
// Updates deployment and service in case of any changes.
// Returns true if need to requeue, false otherwise.
func (r *VmOperatorReconciler) Run(cr *v1beta1.PlatformMonitoring) error {
	r.Log.Info("Reconciling component")

	if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.VmOperator.IsInstall() {
		if !cr.Spec.Victoriametrics.VmOperator.Paused {
			if err := r.handleServiceAccount(cr); err != nil {
				return err
			}
			// Reconcile PodSecurityPolicy, ClusterRole and ClusterRoleBinding only if privileged mode used
			if utils.PrivilegedRights {
				// Create SecurityContextConstraints if run in Openshift and
				// PodSecurityPolicy if run in Kubernetes.
				// Checks that necessary API exists and print message if no suitable API find.
				if r.hasSecurityContextConstraintsAPI() {
					if err := r.handleSecurityContextConstraints(cr); err != nil {
						return err
					}
				} else {
					r.Log.Info("there is no PodSecurityPolicy API found, skip creating PSP for vmoperator")
				}
				if err := r.handleClusterRole(cr); err != nil {
					return err
				}
				if err := r.handleClusterRoleBinding(cr); err != nil {
					return err
				}
			} else {
				r.Log.Info("Skip ClusterRole and ClusterRoleBinding resources reconciliation because privilegedRights=false")
			}
			if err := r.handleRole(cr); err != nil {
				return err
			}
			if err := r.handleRoleBinding(cr); err != nil {
				return err
			}
			if err := r.handleService(cr); err != nil {
				return err
			}
			// Reconcile ServiceMonitor if necessary
			if cr.Spec.Victoriametrics.VmOperator.ServiceMonitor != nil && cr.Spec.Victoriametrics.VmOperator.ServiceMonitor.IsInstall() {
				if err := r.handleServiceMonitor(cr); err != nil {
					return err
				}
			} else {
				r.Log.Info("uninstalling ServiceMonitor")
				if err := r.deleteServiceMonitor(cr); err != nil {
					r.Log.Error(err, "can not delete ServiceMonitor")
				}
			}
			if err := r.handleDeployment(cr); err != nil {
				return err
			}
			if err := r.handleKubeletService(cr); err != nil {
				return err
			}
			if err := r.handleKubeletServiceEndpoints(cr); err != nil {
				return err
			}
			r.Log.Info("Component reconciled")

		} else {
			r.Log.Info("Reconciling paused")
			r.Log.Info("Component NOT reconciled")
		}
	} else {
		r.Log.Info("Uninstalling component if exists")
		r.uninstall(cr)
		r.Log.Info("Component reconciled")
	}

	return nil
}

// uninstall deletes all resources related to the component
func (r *VmOperatorReconciler) uninstall(cr *v1beta1.PlatformMonitoring) {
	if err := r.deleteAllCRDObjects(cr); err != nil {
		r.Log.Error(err, "Can not delete CRD Objects")
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
	if err := r.deleteServiceMonitor(cr); err != nil {
		r.Log.Error(err, "Can not delete Service")
	}
	if err := r.deleteService(cr); err != nil {
		r.Log.Error(err, "Can not delete Service")
	}
	if err := r.deleteVmOperatorConfigMap(cr); err != nil {
		r.Log.Error(err, "Can not delete ConfigMap object")
	}
	// delete custom resource
	if err := r.deleteVmOperatorDeployment(cr); err != nil {
		r.Log.Error(err, "Can not delete Deployment")
	}
	if err := r.deleteServiceAccount(cr); err != nil {
		r.Log.Error(err, "Can not delete ServiceAccount")
	}
}

// hasSecurityContextConstraintsAPI checks that the cluster API has security.openshift.io.v1.SecurityContextConstraints API.
func (r *VmOperatorReconciler) hasSecurityContextConstraintsAPI() bool {
	return r.HasApi(secv1.GroupVersion, "SecurityContextConstraints")
}
