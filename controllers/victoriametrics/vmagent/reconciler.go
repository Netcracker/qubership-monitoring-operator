package vmagent

import (
	"context"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	vmetricsv1b1 "github.com/VictoriaMetrics/operator/api/operator/v1beta1"
	secv1 "github.com/openshift/api/security/v1"
	pspApi "k8s.io/api/policy/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/discovery"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// VmAgentReconciler provides methods to reconcile VmAgent
type VmAgentReconciler struct {
	*utils.ComponentReconciler
}

// NewVmAgentReconciler creates an instance of VmAgentReconciler
func NewVmAgentReconciler(c client.Client, s *runtime.Scheme, dc discovery.DiscoveryInterface) *VmAgentReconciler {
	return &VmAgentReconciler{
		ComponentReconciler: &utils.ComponentReconciler{
			Client: c,
			Scheme: s,
			Dc:     dc,
			Log:    utils.Logger("vmagent_reconciler"),
		},
	}
}

// Run reconciles VmAgent.
// Creates new VmAgent CR if it doesn't exist.
// Updates VmAgent CR in case of any changes.
// Returns true if need to requeue, false otherwise.
func (r *VmAgentReconciler) Run(ctx context.Context, cr *monv1.PlatformMonitoring) error {
	r.Log.Info("Reconciling component")

	if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.VmAgent.IsInstall() {
		if !cr.Spec.Victoriametrics.VmAgent.Paused {
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
			// Reconcile resources with creation
			if err := r.handleRole(cr); err != nil {
				return err
			}
			if err := r.handleRoleBinding(cr); err != nil {
				return err
			}

			// Reconcile vmAgent with creation and update
			if err := r.handleVmAgent(cr); err != nil {
				return err
			}

			// Reconcile Ingress (version v1beta1) if necessary and the cluster is has such API
			// This API unavailable in k8s v1.22+
			if r.HasIngressV1beta1Api() {
				if cr.Spec.Victoriametrics.VmAgent.Ingress != nil &&
					cr.Spec.Victoriametrics.VmAgent.Ingress.IsInstall() &&
					!cr.Spec.Victoriametrics.VmAuth.IsInstall() {
					if err := r.handleIngressV1beta1(cr); err != nil {
						return err
					}
				} else {
					if err := r.deleteIngressV1beta1(cr); err != nil {
						r.Log.Error(err, "Can not delete Ingress")
					}
				}
			}
			// Reconcile Ingress (version v1) if necessary and the cluster is has such API
			// This API available in k8s v1.19+
			if r.HasIngressV1Api() {
				if cr.Spec.Victoriametrics.VmAgent.Ingress != nil &&
					cr.Spec.Victoriametrics.VmAgent.Ingress.IsInstall() &&
					!cr.Spec.Victoriametrics.VmAuth.IsInstall() {
					if err := r.handleIngressV1(cr); err != nil {
						return err
					}
				} else {
					if err := r.deleteIngressV1(cr); err != nil {
						r.Log.Error(err, "Can not delete Ingress")
					}
				}
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
func (r *VmAgentReconciler) uninstall(cr *monv1.PlatformMonitoring) {
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

	// Fetch the VMAgent instance
	vmagent, err := vmAgent(r, cr)
	if err != nil {
		r.Log.Error(err, "Failed creating vmagent manifest. Can not delete vmagent.")
		return
	}

	e := &vmetricsv1b1.VMAgent{ObjectMeta: vmagent.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			return
		}
		r.Log.Error(err, "Can not get vmagent resource")
		return
	}

	if err = r.deleteVmAgent(cr); err != nil {
		r.Log.Error(err, "Can not delete vmagent.")
	}

	// Try to delete Ingress (version v1beta1) is there is such API
	// This API unavailable in k8s v1.22+
	if r.HasIngressV1beta1Api() {
		if err = r.deleteIngressV1beta1(cr); err != nil {
			r.Log.Error(err, "Can not delete Ingress.")
		}
	}
	// Try to delete Ingress (version v1) is there is such API
	// This API available in k8s v1.19+
	if r.HasIngressV1Api() {
		if err = r.deleteIngressV1(cr); err != nil {
			r.Log.Error(err, "Can not delete Ingress.")
		}
	}

	if err := r.deleteServiceAccount(cr); err != nil {
		r.Log.Error(err, "Can not delete ServiceAccount")
	}
}

func (r *VmAgentReconciler) hasPodSecurityPolicyAPI() bool {
	return r.HasApi(pspApi.SchemeGroupVersion, "PodSecurityPolicy")
}

// hasSecurityContextConstraintsAPI checks that the cluster API has security.openshift.io.v1.SecurityContextConstraints API.
func (r *VmAgentReconciler) hasSecurityContextConstraintsAPI() bool {
	return r.HasApi(secv1.GroupVersion, "SecurityContextConstraints")
}
