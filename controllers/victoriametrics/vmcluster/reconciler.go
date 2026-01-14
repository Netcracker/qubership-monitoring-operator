package vmcluster

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

// VmClusterReconciler provides methods to reconcile vmCluster
type VmClusterReconciler struct {
	*utils.ComponentReconciler
}

// NewVmClusterReconciler creates an instance of VmClusterReconciler
func NewVmClusterReconciler(c client.Client, s *runtime.Scheme, dc discovery.DiscoveryInterface) *VmClusterReconciler {
	return &VmClusterReconciler{
		ComponentReconciler: &utils.ComponentReconciler{
			Client: c,
			Scheme: s,
			Dc:     dc,
			Log:    utils.Logger("vmcluster_reconciler"),
		},
	}
}

// Run reconciles vmcluster.
// Creates vmCluster CR if it doesn't exist.
// Updates vmCluster CR in case of any changes.
// Returns true if need to requeue, false otherwise.
func (r *VmClusterReconciler) Run(ctx context.Context, cr *monv1.PlatformMonitoring) error {
	r.Log.Info("Reconciling component")

	if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.VmCluster.IsInstall() && !cr.Spec.Victoriametrics.VmSingle.IsInstall() {
		if !cr.Spec.Victoriametrics.VmCluster.Paused {
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
			// Reconcile vmCluster with creation and update
			if err := r.handleVmCluster(cr); err != nil {
				return err
			}

			// Reconcile Ingress (version v1beta1) if necessary and the cluster is has such API
			// This API unavailable in k8s v1.22+
			if r.HasIngressV1beta1Api() {
				if cr.Spec.Victoriametrics.VmCluster.VmSelectIngress != nil &&
					cr.Spec.Victoriametrics.VmCluster.VmSelectIngress.IsInstall() {
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
				if cr.Spec.Victoriametrics.VmCluster.VmSelectIngress != nil &&
					cr.Spec.Victoriametrics.VmCluster.VmSelectIngress.IsInstall() {
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
func (r *VmClusterReconciler) uninstall(cr *monv1.PlatformMonitoring) {
	if utils.PrivilegedRights {
		if err := r.deleteClusterRole(cr); err != nil {
			r.Log.Error(err, "Can not delete ClusterRole")
		}
		if err := r.deleteClusterRoleBinding(cr); err != nil {
			r.Log.Error(err, "Can not delete ClusterRoleBinding")
		}
	}

	// Fetch the VMCluster instance
	vmCluster, err := vmCluster(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating vmcluster manifest. Can not delete vmcluster.")
		return
	}

	e := &vmetricsv1b1.VMCluster{ObjectMeta: vmCluster.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			return
		}
		r.Log.Error(err, "Can not get vmcluster resource")
		return
	}

	if err = r.deleteVmCluster(cr); err != nil {
		r.Log.Error(err, "Can not delete vmcluster.")
	}

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
func (r *VmClusterReconciler) hasSecurityContextConstraintsAPI() bool {
	return r.HasApi(secv1.GroupVersion, "SecurityContextConstraints")
}

// hasPodSecurityPolicyAPI checks that the cluster API has policy.v1beta.PodSecurityPolicy API.
func (r *VmClusterReconciler) hasPodSecurityPolicyAPI() bool {
	return r.HasApi(pspApi.SchemeGroupVersion, "PodSecurityPolicy")
}
