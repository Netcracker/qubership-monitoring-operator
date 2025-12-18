package vmuser

import (
	v1beta1 "github.com/Netcracker/qubership-monitoring-operator/api/v1beta1"
	vmetricsv1b1 "github.com/VictoriaMetrics/operator/api/operator/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
)

func (r *VmUserReconciler) handleVmUser(cr *v1beta1.PlatformMonitoring) error {
	m, err := vmUser(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating vmuser manifest")
		return err
	}
	e := &vmetricsv1b1.VMUser{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			if err = r.CreateResource(cr, m); err != nil {
				return err
			}
			return nil
		}
		return err
	}

	//Set parameters
	e.SetLabels(m.GetLabels())
	e.Spec = m.Spec

	if err = r.UpdateResource(e); err != nil {
		return err
	}
	return nil
}

func (r *VmUserReconciler) deleteVmUser(cr *v1beta1.PlatformMonitoring) error {
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
