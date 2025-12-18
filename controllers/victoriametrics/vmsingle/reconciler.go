package vmsingle

import (
	"context"

	v1beta1 "github.com/Netcracker/qubership-monitoring-operator/api/v1beta1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	vmetricsv1b1 "github.com/VictoriaMetrics/operator/api/operator/v1beta1"
	secv1 "github.com/openshift/api/security/v1"
	pspApi "k8s.io/api/policy/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/discovery"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// VmSingleReconciler provides methods to reconcile vmSingle
type VmSingleReconciler struct {
	*utils.ComponentReconciler
}

// NewVmSingleReconciler creates an instance of VmSingleReconciler
func NewVmSingleReconciler(c client.Client, s *runtime.Scheme, dc discovery.DiscoveryInterface) *VmSingleReconciler {
	return &VmSingleReconciler{
		ComponentReconciler: &utils.ComponentReconciler{
			Client: c,
			Scheme: s,
			Dc:     dc,
			Log:    utils.Logger("vmsingle_reconciler"),
		},
	}
}

// Run reconciles vmsingle.
// Creates vmSingle CR if it doesn't exist.
// Updates vmSingle CR in case of any changes.
// Returns true if need to requeue, false otherwise.
func (r *VmSingleReconciler) Run(ctx context.Context, cr *v1beta1.PlatformMonitoring) error {
	r.Log.Info("Reconciling component")

	if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.VmSingle.IsInstall() && !cr.Spec.Victoriametrics.VmCluster.IsInstall() {
		if !cr.Spec.Victoriametrics.VmSingle.Paused {
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
			// Reconcile vmSingle with creation and update
			if err := r.handleVmSingle(cr); err != nil {
				return err
			}

			// Reconcile Ingress (version v1beta1) if necessary and the cluster is has such API
			// This API unavailable in k8s v1.22+
			if r.HasIngressV1beta1Api() {
				if cr.Spec.Victoriametrics.VmSingle.Ingress != nil &&
					cr.Spec.Victoriametrics.VmSingle.Ingress.IsInstall() &&
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
				if cr.Spec.Victoriametrics.VmSingle.Ingress != nil &&
					cr.Spec.Victoriametrics.VmSingle.Ingress.IsInstall() &&
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
func (r *VmSingleReconciler) uninstall(cr *v1beta1.PlatformMonitoring) {
	if utils.PrivilegedRights {
		if err := r.deleteClusterRole(cr); err != nil {
			r.Log.Error(err, "Can not delete ClusterRole")
		}
		if err := r.deleteClusterRoleBinding(cr); err != nil {
			r.Log.Error(err, "Can not delete ClusterRoleBinding")
		}
	}

	// Fetch the VMSingle instance
	vmSingle, err := vmSingle(r, cr)
	if err != nil {
		r.Log.Error(err, "Failed creating vmsingle manifest. Can not delete vmsingle.")
		return
	}

	e := &vmetricsv1b1.VMSingle{ObjectMeta: vmSingle.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			return
		}
		r.Log.Error(err, "Can not get vmsingle resource")
		return
	}

	if err = r.deleteVmSingle(cr); err != nil {
		r.Log.Error(err, "Can not delete vmsingle.")
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

// hasSecurityContextConstraintsAPI checks that the cluster API has security.openshift.io.v1.SecurityContextConstraints API.
func (r *VmSingleReconciler) hasSecurityContextConstraintsAPI() bool {
	return r.HasApi(secv1.GroupVersion, "SecurityContextConstraints")
}

// hasPodSecurityPolicyAPI checks that the cluster API has policy.v1beta.PodSecurityPolicy API.
func (r *VmSingleReconciler) hasPodSecurityPolicyAPI() bool {
	return r.HasApi(pspApi.SchemeGroupVersion, "PodSecurityPolicy")
}
