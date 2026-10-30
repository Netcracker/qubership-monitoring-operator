package prometheus_rules

import (
	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
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
	return utils.CreateOrUpdate(r.ComponentReconciler, cr, m, e, func(existing, desired *promv1.PrometheusRule) bool {
		return utils.AnyChanged(
			utils.SetLabelsIfChanged(existing, desired.GetLabels()),
			utils.SetIfChanged(&existing.Spec, desired.Spec),
		)
	})
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
