package gateway

import (
	"context"
	"encoding/json"
	"fmt"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	gatewayAPIGroup         = "gateway.networking.k8s.io"
	httpRouteKind           = "HTTPRoute"
	httpRouteNameSuffix     = "-http-route"
	pathPrefixPathMatchType = "PathPrefix"
)

type GatewayRouteConfig struct {
	NamePrefix     string
	Namespace      string
	Host           string
	ServiceName    string
	ServicePort    int32
	Labels         map[string]string
	ParentRefs     []monv1.GatewayParentRef
	ComponentRoute *monv1.GatewayHTTPRoute
}

func ReconcileGatewayRoutes(r *utils.ComponentReconciler, cr *monv1.PlatformMonitoring, cfg GatewayRouteConfig) error {
	routeCfg := cfg.ComponentRoute

	needHTTPRoute := false
	if routeCfg != nil && routeCfg.Install != nil {
		needHTTPRoute = *routeCfg.Install
	}

	routeGroup := gatewayAPIGroupForConfig(cfg, routeCfg)
	if !needHTTPRoute {
		// Skip cleanup when Gateway API was never configured for this component,
		// to avoid unnecessary API-discovery calls on every reconcile.
		if routeCfg == nil && len(cfg.ParentRefs) == 0 {
			return nil
		}
		desired, err := buildHTTPRoute(cfg, nil, nil)
		if err != nil {
			return err
		}
		return deleteGatewayResource(r, newGatewayAPIResolver(r), desired, routeGroup)
	}

	resolver := newGatewayAPIResolver(r)
	if err := validateGatewayParentRefs(resolver, cfg, routeCfg); err != nil {
		return err
	}

	hostnames := getHTTPRouteHostnames(cfg, routeCfg)
	if len(hostnames) == 0 {
		return errors.NewBadRequest("can not reconcile HTTPRoute without hostnames: set ingress.host or component httpRoute.hostnames")
	}

	desired, err := buildHTTPRoute(cfg, routeCfg, hostnames)
	if err != nil {
		return err
	}
	return applyGatewayResource(r, resolver, cr, desired, routeGroup)
}

func DeleteGatewayRoutes(r *utils.ComponentReconciler, cfg GatewayRouteConfig) error {
	routeGroup := gatewayAPIGroupForConfig(cfg, cfg.ComponentRoute)
	desired, err := buildHTTPRoute(cfg, nil, nil)
	if err != nil {
		return err
	}
	return deleteGatewayResource(r, newGatewayAPIResolver(r), desired, routeGroup)
}

func applyGatewayResource(r *utils.ComponentReconciler, resolver *gatewayAPIResolver, cr *monv1.PlatformMonitoring, desired *unstructured.Unstructured, group string) error {
	gvk, ok := resolver.ResolveGroupVersionKind(group, desired.GetKind())
	if !ok {
		r.Log.Info(
			"Skip Gateway API resource reconciliation because API is not available",
			utils.ResourceKey, desired.GetKind(),
			"apiGroup", group,
		)
		return nil
	}
	desired.SetGroupVersionKind(gvk)

	current := &unstructured.Unstructured{}
	current.SetGroupVersionKind(gvk)
	current.SetName(desired.GetName())
	current.SetNamespace(desired.GetNamespace())
	if err := r.Client.Get(context.TODO(), client.ObjectKeyFromObject(current), current); err != nil {
		if errors.IsNotFound(err) {
			if err := controllerutil.SetControllerReference(cr, desired, r.Scheme); err != nil {
				r.Log.Error(err, "Failed to set controller reference", utils.ResourceKey, desired.GetKind())
				return err
			}
			if err := r.Client.Create(context.TODO(), desired); err != nil {
				return err
			}
			r.Log.Info("Successful creating", utils.ResourceKey, desired.GetKind())
			return nil
		}
		return err
	}

	logHTTPRouteParentStatus(r, current)
	current.SetLabels(desired.GetLabels())
	current.SetAnnotations(desired.GetAnnotations())
	current.Object["spec"] = desired.Object["spec"]
	if err := r.Client.Update(context.TODO(), current); err != nil {
		return err
	}
	r.Log.Info("Successful updating", utils.ResourceKey, current.GetKind())
	return nil
}

func deleteGatewayResource(r *utils.ComponentReconciler, resolver *gatewayAPIResolver, desired *unstructured.Unstructured, group string) error {
	gvk, ok := resolver.ResolveGroupVersionKind(group, desired.GetKind())
	if !ok {
		return nil
	}
	desired.SetGroupVersionKind(gvk)

	current := &unstructured.Unstructured{}
	current.SetGroupVersionKind(gvk)
	current.SetName(desired.GetName())
	current.SetNamespace(desired.GetNamespace())
	if err := r.Client.Get(context.TODO(), client.ObjectKeyFromObject(current), current); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if err := r.Client.Delete(context.TODO(), current); err != nil {
		return err
	}
	r.Log.Info("Successful deleting", utils.ResourceKey, current.GetKind())
	return nil
}

type gatewayAPIResolver struct {
	r        *utils.ComponentReconciler
	apiLists []*metav1.APIResourceList
	loaded   bool
}

func newGatewayAPIResolver(r *utils.ComponentReconciler) *gatewayAPIResolver {
	return &gatewayAPIResolver{r: r}
}

func (resolver *gatewayAPIResolver) ResolveGroupVersionKind(group, kind string) (schema.GroupVersionKind, bool) {
	apiVersion, ok := resolver.GetAPIVersionForKind(group, kind)
	if !ok {
		return schema.GroupVersionKind{}, false
	}
	gv, err := schema.ParseGroupVersion(apiVersion)
	if err != nil {
		return schema.GroupVersionKind{}, false
	}
	return gv.WithKind(kind), true
}

func (resolver *gatewayAPIResolver) HasAPIGroupKind(group, kind string) bool {
	_, ok := resolver.GetAPIVersionForKind(group, kind)
	return ok
}

func (resolver *gatewayAPIResolver) GetAPIVersionForKind(group, kind string) (string, bool) {
	if !resolver.loaded {
		resolver.loaded = true
		_, apiLists, err := resolver.r.Dc.ServerGroupsAndResources()
		if err != nil {
			resolver.r.Log.Error(err, "Error while check hasAPI")
			return "", false
		}
		resolver.apiLists = apiLists
	}
	for _, apiList := range resolver.apiLists {
		gv, err := schema.ParseGroupVersion(apiList.GroupVersion)
		if err != nil {
			continue
		}
		if gv.Group != group {
			continue
		}
		for _, resource := range apiList.APIResources {
			if resource.Kind == kind {
				return apiList.GroupVersion, true
			}
		}
	}
	return "", false
}

func buildHTTPRoute(cfg GatewayRouteConfig, routeCfg *monv1.GatewayHTTPRoute, hostnames []string) (*unstructured.Unstructured, error) {
	group := gatewayAPIGroupForConfig(cfg, routeCfg)
	spec := map[string]interface{}{
		"hostnames": toStringInterfaces(hostnames),
		"rules": []interface{}{
			map[string]interface{}{
				"matches": []interface{}{
					map[string]interface{}{
						"path": map[string]interface{}{
							"type":  pathPrefixPathMatchType,
							"value": "/",
						},
					},
				},
				"backendRefs": []interface{}{
					map[string]interface{}{
						"name": cfg.ServiceName,
						"port": cfg.ServicePort,
					},
				},
			},
		},
	}

	if routeCfg != nil {
		if parentRefs := buildParentRefs(routeCfg.ParentRefs, cfg.ParentRefs); len(parentRefs) > 0 {
			spec["parentRefs"] = parentRefs
		}
		if len(routeCfg.Rules) > 0 {
			rules, err := buildHTTPRouteRules(cfg, routeCfg.Rules)
			if err != nil {
				return nil, err
			}
			spec["rules"] = rules
		}
	}

	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": group + "/v1",
			"kind":       httpRouteKind,
			"metadata": map[string]interface{}{
				"name":      cfg.NamePrefix + httpRouteNameSuffix,
				"namespace": cfg.Namespace,
				"labels":    cfg.Labels,
			},
			"spec": spec,
		},
	}, nil
}

func buildParentRefs(routeParentRefs, defaultParentRefs []monv1.GatewayParentRef) []interface{} {
	parentRefs := gatewayParentRefsForConfig(routeParentRefs, defaultParentRefs)
	result := make([]interface{}, 0, len(parentRefs))
	for _, parentRef := range parentRefs {
		if parentRef.Name == "" {
			continue
		}
		entry := map[string]interface{}{
			"name": parentRef.Name,
		}
		if parentRef.Group != "" {
			entry["group"] = parentRef.Group
		} else {
			entry["group"] = gatewayAPIGroup
		}
		if parentRef.Kind != "" {
			entry["kind"] = parentRef.Kind
		} else {
			entry["kind"] = "Gateway"
		}
		if parentRef.Namespace != "" {
			entry["namespace"] = parentRef.Namespace
		}
		if parentRef.SectionName != "" {
			entry["sectionName"] = parentRef.SectionName
		}
		result = append(result, entry)
	}
	return result
}

func buildHTTPRouteRules(cfg GatewayRouteConfig, rules []monv1.GatewayHTTPRouteRule) ([]interface{}, error) {
	result := make([]interface{}, 0, len(rules))
	for ruleIdx, rule := range rules {
		item := map[string]interface{}{
			"backendRefs": []interface{}{
				map[string]interface{}{
					"name": cfg.ServiceName,
					"port": cfg.ServicePort,
				},
			},
		}
		matches, err := rawJSONSliceToInterfaces(rule.Matches)
		if err != nil {
			return nil, fmt.Errorf("invalid HTTPRoute rule %d matches: %w", ruleIdx, err)
		}
		if len(matches) > 0 {
			item["matches"] = matches
		}
		filters, err := rawJSONSliceToInterfaces(rule.Filters)
		if err != nil {
			return nil, fmt.Errorf("invalid HTTPRoute rule %d filters: %w", ruleIdx, err)
		}
		if len(filters) > 0 {
			item["filters"] = filters
		}
		result = append(result, item)
	}
	return result, nil
}

func rawJSONSliceToInterfaces(items []monv1.GatewayJSON) ([]interface{}, error) {
	result := make([]interface{}, 0, len(items))
	for itemIdx, item := range items {
		if len(item.Raw) == 0 {
			continue
		}
		var parsed interface{}
		if err := json.Unmarshal(item.Raw, &parsed); err != nil {
			return nil, fmt.Errorf("item %d: %w", itemIdx, err)
		}
		result = append(result, parsed)
	}
	return result, nil
}

func toStringInterfaces(values []string) []interface{} {
	result := make([]interface{}, 0, len(values))
	for _, value := range values {
		result = append(result, value)
	}
	return result
}

func getHTTPRouteHostnames(cfg GatewayRouteConfig, routeCfg *monv1.GatewayHTTPRoute) []string {
	if routeCfg != nil && len(routeCfg.Hostnames) > 0 {
		return routeCfg.Hostnames
	}
	if cfg.Host == "" {
		return nil
	}
	return []string{cfg.Host}
}

func gatewayAPIGroupForConfig(cfg GatewayRouteConfig, routeCfg *monv1.GatewayHTTPRoute) string {
	var routeParentRefs []monv1.GatewayParentRef
	if routeCfg != nil {
		routeParentRefs = routeCfg.ParentRefs
	}
	parentRefs := gatewayParentRefsForConfig(routeParentRefs, cfg.ParentRefs)
	return gatewayAPIGroupForParentRefs(parentRefs)
}

func gatewayAPIGroupForParentRefs(parentRefs []monv1.GatewayParentRef) string {
	for _, parentRef := range parentRefs {
		if parentRef.Group != "" {
			return parentRef.Group
		}
	}
	return gatewayAPIGroup
}

func validateGatewayParentRefs(resolver *gatewayAPIResolver, cfg GatewayRouteConfig, routeCfg *monv1.GatewayHTTPRoute) error {
	var routeParentRefs []monv1.GatewayParentRef
	if routeCfg != nil {
		routeParentRefs = routeCfg.ParentRefs
	}
	parentRefs := gatewayParentRefsForConfig(routeParentRefs, cfg.ParentRefs)
	if err := validateGatewayParentRefGroups(parentRefs); err != nil {
		return err
	}
	for _, parentRef := range parentRefs {
		parentKind := parentRef.Kind
		if parentKind == "" {
			parentKind = "Gateway"
		}
		if parentRef.Group != "" && !resolver.HasAPIGroupKind(parentRef.Group, parentKind) {
			return errors.NewBadRequest("Gateway API group for parentRef is not available: " + parentRef.Group)
		}
	}
	return nil
}

type httpRouteParentStatusWarning struct {
	ParentName      string
	ParentNamespace string
	ParentGroup     string
	ParentKind      string
	SectionName     string
	Reason          string
	Message         string
}

func logHTTPRouteParentStatus(r *utils.ComponentReconciler, route *unstructured.Unstructured) {
	for _, warning := range httpRouteParentStatusWarnings(route) {
		args := []interface{}{
			utils.ResourceKey, route.GetKind(),
			"name", route.GetName(),
			"namespace", route.GetNamespace(),
			"reason", warning.Reason,
			"statusMessage", warning.Message,
		}
		if warning.ParentName != "" {
			args = append(args, "parentName", warning.ParentName)
		}
		if warning.ParentNamespace != "" {
			args = append(args, "parentNamespace", warning.ParentNamespace)
		}
		if warning.ParentGroup != "" {
			args = append(args, "parentGroup", warning.ParentGroup)
		}
		if warning.ParentKind != "" {
			args = append(args, "parentKind", warning.ParentKind)
		}
		if warning.SectionName != "" {
			args = append(args, "sectionName", warning.SectionName)
		}
		utils.Warn(r.Log, "Gateway API HTTPRoute parent status is not healthy", args...)
	}
}

func httpRouteParentStatusWarnings(route *unstructured.Unstructured) []httpRouteParentStatusWarning {
	if _, found, _ := unstructured.NestedMap(route.Object, "status"); !found {
		return nil
	}
	parents, found, _ := unstructured.NestedSlice(route.Object, "status", "parents")
	if !found || len(parents) == 0 {
		return []httpRouteParentStatusWarning{
			{
				Reason:  "NoParents",
				Message: "HTTPRoute status.parents is empty",
			},
		}
	}

	result := make([]httpRouteParentStatusWarning, 0)
	for _, parent := range parents {
		parentMap, ok := parent.(map[string]interface{})
		if !ok {
			continue
		}
		warning := httpRouteParentStatusWarning{
			Reason:  "AcceptedConditionMissing",
			Message: "HTTPRoute parent status has no Accepted condition",
		}
		if parentRef, ok := parentMap["parentRef"].(map[string]interface{}); ok {
			warning.ParentName = stringFromMap(parentRef, "name")
			warning.ParentNamespace = stringFromMap(parentRef, "namespace")
			warning.ParentGroup = stringFromMap(parentRef, "group")
			warning.ParentKind = stringFromMap(parentRef, "kind")
			warning.SectionName = stringFromMap(parentRef, "sectionName")
		}
		conditions, ok := parentMap["conditions"].([]interface{})
		if !ok {
			result = append(result, warning)
			continue
		}
		acceptedFound := false
		for _, condition := range conditions {
			conditionMap, ok := condition.(map[string]interface{})
			if !ok {
				continue
			}
			switch stringFromMap(conditionMap, "type") {
			case "Accepted":
				acceptedFound = true
				if stringFromMap(conditionMap, "status") != "True" {
					result = append(result, warningForCondition(warning, conditionMap, "AcceptedFalse", "HTTPRoute parent Accepted condition is not True"))
				}
			case "ResolvedRefs":
				if stringFromMap(conditionMap, "status") != "True" {
					result = append(result, warningForCondition(warning, conditionMap, "ResolvedRefsFalse", "HTTPRoute parent ResolvedRefs condition is not True"))
				}
			}
		}
		if !acceptedFound {
			result = append(result, warning)
		}
	}
	return result
}

func warningForCondition(base httpRouteParentStatusWarning, condition map[string]interface{}, defaultReason, defaultMessage string) httpRouteParentStatusWarning {
	base.Reason = stringFromMap(condition, "reason")
	if base.Reason == "" {
		base.Reason = defaultReason
	}
	base.Message = stringFromMap(condition, "message")
	if base.Message == "" {
		base.Message = defaultMessage
	}
	return base
}

func stringFromMap(values map[string]interface{}, key string) string {
	value, _ := values[key].(string)
	return value
}

func validateGatewayParentRefGroups(parentRefs []monv1.GatewayParentRef) error {
	routeGroup := gatewayAPIGroupForParentRefs(parentRefs)
	for _, parentRef := range parentRefs {
		parentGroup := gatewayAPIGroup
		if parentRef.Group != "" {
			parentGroup = parentRef.Group
		}
		if parentGroup != routeGroup {
			return errors.NewBadRequest("all Gateway API parentRefs in HTTPRoute must use the same group")
		}
	}
	return nil
}

func gatewayParentRefsForConfig(routeParentRefs, defaultParentRefs []monv1.GatewayParentRef) []monv1.GatewayParentRef {
	if len(routeParentRefs) > 0 {
		return routeParentRefs
	}
	return defaultParentRefs
}
