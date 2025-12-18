package prometheus_rules

import (
	v1beta1 "github.com/Netcracker/qubership-monitoring-operator/api/v1beta1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// PrometheusRulesReconciler is a reconciler to maintain configuration of prometheus rules
type PrometheusRulesReconciler struct {
	*utils.ComponentReconciler
}

// NewPrometheusRulesReconciler  returns PrometheusRulesReconciler by specified parameters
func NewPrometheusRulesReconciler(c client.Client, s *runtime.Scheme) *PrometheusRulesReconciler {
	return &PrometheusRulesReconciler{
		ComponentReconciler: &utils.ComponentReconciler{
			Client: c,
			Scheme: s,
			Log:    utils.Logger("prometheus_rules_reconciler"),
		},
	}
}

// Run reconciles k8s prometheus rules
// Creates, updates and deletes prometheus rules depending of configuration
func (r *PrometheusRulesReconciler) Run(cr *v1beta1.PlatformMonitoring) error {
	r.Log.Info("Reconciling component")

	if cr.Spec.PrometheusRules != nil && cr.Spec.PrometheusRules.IsInstall() && len(cr.Spec.PrometheusRules.RuleGroups) > 0 {
		if err := r.handlePrometheusRules(cr); err != nil {
			return err
		}
	} else {
		r.Log.Info("Uninstalling PrometheusRules")
		if err := r.deletePrometheusRules(cr); err != nil {
			r.Log.Error(err, "Can not delete PrometheusRules")
		}
	}
	r.Log.Info("Component reconciled")
	return nil
}
