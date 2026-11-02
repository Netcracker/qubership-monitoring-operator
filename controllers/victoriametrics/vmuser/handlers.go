package vmuser

import (
	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	vmetricsv1b1 "github.com/VictoriaMetrics/operator/api/operator/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
)

func (r *VmUserReconciler) handleVmUser(cr *monv1.PlatformMonitoring) error {
	m, err := vmUser(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating vmuser manifest")
		return err
	}
	e := &vmetricsv1b1.VMUser{ObjectMeta: m.ObjectMeta}
	return utils.CreateOrUpdate(r.ComponentReconciler, cr, m, e, func(existing, desired *vmetricsv1b1.VMUser) bool {
		return utils.AnyChanged(
			utils.SetLabelsIfChanged(existing, desired.GetLabels()),
			utils.SetIfChanged(&existing.Spec, desired.Spec),
		)
	})
}

func (r *VmUserReconciler) deleteVmUser(cr *monv1.PlatformMonitoring) error {
	m, err := vmUser(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating vmUser manifest")
		return err
	}
	e := &vmetricsv1b1.VMUser{ObjectMeta: m.ObjectMeta}
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
