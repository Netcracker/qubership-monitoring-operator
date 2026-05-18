package prometheus_rules

import (
	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"k8s.io/apimachinery/pkg/api/errors"
)

func (r *PrometheusRulesReconciler) handlePrometheusRules(cr *monv1.PlatformMonitoring) error {
	m, err := prometheusRules(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating PrometheusRules manifest")
		return err
	}

	e := &promv1.PrometheusRule{ObjectMeta: m.ObjectMeta}
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

func (r *PrometheusRulesReconciler) deletePrometheusRules(cr *monv1.PlatformMonitoring) error {
	m, err := prometheusRules(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating PrometheusRules manifest")
		return err
	}
	e := &promv1.PrometheusRule{ObjectMeta: m.ObjectMeta}
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
