package vmauth

import (
	"context"

	v1alpha1 "github.com/Netcracker/qubership-monitoring-operator/api/v1alpha1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	vmetricsv1b1 "github.com/VictoriaMetrics/operator/api/operator/v1beta1"
	secv1 "github.com/openshift/api/security/v1"
	pspApi "k8s.io/api/policy/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/discovery"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// VmAuthReconciler provides methods to reconcile VmAuth
type VmAuthReconciler struct {
	*utils.ComponentReconciler
}

// NewVmAuthReconciler creates an instance of VmAuthReconciler
func NewVmAuthReconciler(c client.Client, s *runtime.Scheme, dc discovery.DiscoveryInterface) *VmAuthReconciler {
	return &VmAuthReconciler{
		ComponentReconciler: &utils.ComponentReconciler{
			Client: c,
			Scheme: s,
			Dc:     dc,
			Log:    utils.Logger("vmauth_reconciler"),
		},
	}
}

// Run reconciles vmauth.
// Creates VmAuth CR if it doesn't exist.
// Updates VmAuth CR in case of any changes.
// Returns true if need to requeue, false otherwise.
func (r *VmAuthReconciler) Run(ctx context.Context, cr *v1alpha1.PlatformMonitoring) error {
	r.Log.Info("Reconciling component")

	if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.VmAuth.IsInstall() {
		if !cr.Spec.Victoriametrics.VmAuth.Paused {
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
			// Reconcile Role and RoleBinding
			if err := r.handleRole(cr); err != nil {
				return err
			}
			if err := r.handleRoleBinding(cr); err != nil {
				return err
			}

			// Reconcile VmSingle with creation and update
			if err := r.handleVmAuth(cr); err != nil {
				return err
			}

			// Reconcile Ingress (version v1) if necessary and the cluster is has such API
			// This API available in k8s v1.19+
			if r.HasIngressV1Api() {
				if cr.Spec.Victoriametrics.VmAuth.Ingress != nil && cr.Spec.Victoriametrics.VmAuth.Ingress.IsInstall() {
					if err := r.handleIngress(cr); err != nil {
						return err
					}
				} else {
					if err := r.deleteIngress(cr); err != nil {
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
func (r *VmAuthReconciler) uninstall(cr *v1alpha1.PlatformMonitoring) {
	if utils.PrivilegedRights {
		if err := r.deleteClusterRole(cr); err != nil {
			r.Log.Error(err, "Can not delete ClusterRole")
		}
		if err := r.deleteClusterRoleBinding(cr); err != nil {
			r.Log.Error(err, "Can not delete ClusterRoleBinding")
		}
	}

	// Fetch the VMSingle instance
	vmSingle, err := vmAuth(r, cr)
	if err != nil {
		r.Log.Error(err, "Failed creating vmauth manifest")
	}

	e := &vmetricsv1b1.VMAuth{ObjectMeta: vmSingle.ObjectMeta}
	r.Log.Info("Get resource VMAuth instance")
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			return
		}
		r.Log.Error(err, "Can not get vmauth resource")
	}

	if err = r.deleteVmAuth(cr); err != nil {
		r.Log.Error(err, "Can not delete vmauth")
	}

	if r.HasIngressV1Api() {
		if err = r.deleteIngress(cr); err != nil {
			r.Log.Error(err, "Can not delete Ingress")
		}
	}

	if err := r.deleteServiceAccount(cr); err != nil {
		r.Log.Error(err, "Can not delete ServiceAccount")
	}
}

func (r *VmAuthReconciler) hasPodSecurityPolicyAPI() bool {
	return r.HasApi(pspApi.SchemeGroupVersion, "PodSecurityPolicy")
}

// hasSecurityContextConstraintsAPI checks that the cluster API has security.openshift.io.v1.SecurityContextConstraints API.
func (r *VmAuthReconciler) hasSecurityContextConstraintsAPI() bool {
	return r.HasApi(secv1.GroupVersion, "SecurityContextConstraints")
}
