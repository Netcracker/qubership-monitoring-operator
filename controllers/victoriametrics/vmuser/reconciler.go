package vmuser

import (
	"context"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	vmetricsv1b1 "github.com/VictoriaMetrics/operator/api/operator/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/discovery"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// VmUserReconciler provides methods to reconcile vmSingle
type VmUserReconciler struct {
	*utils.ComponentReconciler
}

// NewVmUserReconciler creates an instance of VmUserReconciler
func NewVmUserReconciler(c client.Client, s *runtime.Scheme, dc discovery.DiscoveryInterface) *VmUserReconciler {
	return &VmUserReconciler{
		ComponentReconciler: &utils.ComponentReconciler{
			Client: c,
			Scheme: s,
			Dc:     dc,
			Log:    utils.Logger("vmuser_reconciler"),
		},
	}
}

// Run reconciles vmuser.
// Creates vmUser CR if it doesn't exist.
// Updates vmUser CR in case of any changes.
// Returns true if need to requeue, false otherwise.
func (r *VmUserReconciler) Run(ctx context.Context, cr *monv1.PlatformMonitoring) error {
	r.Log.Info("Reconciling component")

	if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.VmUser.IsInstall() && cr.Spec.Victoriametrics.VmAuth.IsInstall() {
		if !cr.Spec.Victoriametrics.VmUser.Paused {

			// Reconcile vmUser with creation and update
			if err := r.handleVmUser(cr); err != nil {
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
func (r *VmUserReconciler) uninstall(cr *monv1.PlatformMonitoring) {
	// Fetch the VMSingle instance
	m, err := vmUser(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating vmuser manifest. Can not delete vmuser.")
		return
	}

	e := &vmetricsv1b1.VMUser{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			return
		}
		r.Log.Error(err, "Can not get vmuser resource")
		return
	}

	if err = r.deleteVmUser(cr); err != nil {
		r.Log.Error(err, "Can not delete vmuser.")
	}
}
