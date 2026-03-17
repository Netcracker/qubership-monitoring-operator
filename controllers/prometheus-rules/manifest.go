package prometheus_rules

import (
	"embed"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/yaml"
)

//go:embed  assets/*.yaml
var assets embed.FS

func prepareOverrideConfigMap(cr *monv1.PlatformMonitoring) map[string]map[string]*monv1.PrometheusRule {
	// Init map with info about chosen groups and overridden rules from CR
	overrideConfigMap := make(map[string]map[string]*monv1.PrometheusRule)
	restrictedGroups := make(map[string]bool)
	if cr.Spec.PublicCloudName != "" {
		if pcMap, ok := utils.PublicCloudRulesEnabled[cr.Spec.PublicCloudName]; ok {
			for group, included := range pcMap {
				if included {
					cr.Spec.PrometheusRules.RuleGroups = append(cr.Spec.PrometheusRules.RuleGroups, group)
				} else {
					restrictedGroups[group] = true
				}
			}
		}
	}
	for _, group := range cr.Spec.PrometheusRules.RuleGroups {
		// Groups that should not be in the public cloud will be skipped
		if restrictedGroups[group] {
			continue
		}
		if overrideConfigMap[group] == nil {
			overrideConfigMap[group] = make(map[string]*monv1.PrometheusRule)
		}
		for i := range cr.Spec.PrometheusRules.Override {
			rule := cr.Spec.PrometheusRules.Override[i]
			if rule.Group == group {
				if rule.Alert != "" {
					overrideConfigMap[rule.Group][rule.Alert] = &rule
				} else if rule.Record != "" {
					overrideConfigMap[rule.Group][rule.Record] = &rule
				} else {
					continue
				}
			}

		}
	}
	return overrideConfigMap
}

func prometheusRules(cr *monv1.PlatformMonitoring) (*promv1.PrometheusRule, error) {
	rules := promv1.PrometheusRule{}

	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.PrometheusRulesAsset), 100).Decode(&rules); err != nil {
		return nil, err
	}
	//Set parameters
	rules.SetGroupVersionKind(schema.GroupVersionKind{Group: "monitoring.coreos.com", Version: "v1", Kind: "PrometheusRule"})
	rules.SetNamespace(cr.GetNamespace())

	if cr.Spec.PrometheusRules != nil {
		// Init map with info about chosen groups and overridden rules from CR
		overrideConfigMap := prepareOverrideConfigMap(cr)

		// Result spec contains only chosen groups with override alerts
		resultSpec := promv1.PrometheusRuleSpec{}

		// Iterate over Rules Groups from asset. Get those chosen in CR and apply overridden configs if needed
		for i := range rules.Spec.Groups {
			group := &rules.Spec.Groups[i]
			if overrideRulesByGroup, isGroupIncluded := overrideConfigMap[group.Name]; isGroupIncluded {

				// Do not include AlertManager group in case of AlertManager is not installed
				if group.Name == utils.AlertManagerGroupName && (cr.Spec.AlertManager == nil || !cr.Spec.AlertManager.IsInstall()) {
					continue
				}
				if overrideRulesByGroup != nil {
					// Search possible override config and apply it for each rule in manifest
					for j := range group.Rules {
						rule := &group.Rules[j]
						if rule.Alert != "" {
							if _, ok := overrideRulesByGroup[rule.Alert]; ok {
								overrideRulesByGroup[rule.Alert].OverridePrometheusRule(rule)
							}
						} else if rule.Record != "" {
							if _, ok := overrideRulesByGroup[rule.Record]; ok {
								overrideRulesByGroup[rule.Record].OverridePrometheusRule(rule)
							}
						} else {
							continue
						}
					}
				}
				resultSpec.Groups = append(resultSpec.Groups, *group)
			}
		}
		rules.Spec = resultSpec
	}

	if rules.Labels == nil && cr.GetLabels() != nil {
		rules.SetLabels(cr.Labels)
	} else {
		for k, v := range cr.GetLabels() {
			if _, ok := rules.Labels[k]; !ok {
				rules.Labels[k] = v
			}
		}
	}

	if rules.Annotations == nil && cr.GetAnnotations() != nil {
		rules.SetAnnotations(cr.GetAnnotations())
	} else {
		for k, v := range cr.GetAnnotations() {
			rules.Annotations[k] = v
		}
	}

	return &rules, nil
}
