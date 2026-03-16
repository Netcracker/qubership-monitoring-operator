package victoriametrics

import (
	"context"
	"encoding/json"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	"k8s.io/apimachinery/pkg/api/errors"
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
		return deleteGatewayResource(r, buildHTTPRoute(cfg, nil, nil), routeGroup)
	}

	if err := validateGatewayParentRefs(r, cfg, routeCfg); err != nil {
		return err
	}

	hostnames := getHTTPRouteHostnames(cfg, routeCfg)
	if len(hostnames) == 0 {
		return errors.NewBadRequest("can not reconcile HTTPRoute without hostnames: set ingress.host or component httpRoute.hostnames")
	}

	return applyGatewayResource(r, cr, buildHTTPRoute(cfg, routeCfg, hostnames), routeGroup)
}

func DeleteGatewayRoutes(r *utils.ComponentReconciler, cfg GatewayRouteConfig) error {
	return deleteGatewayResource(r, buildHTTPRoute(cfg, nil, nil), gatewayAPIGroupForConfig(cfg, nil))
}

func applyGatewayResource(r *utils.ComponentReconciler, cr *monv1.PlatformMonitoring, desired *unstructured.Unstructured, group string) error {
	gvk, ok := resolveGroupVersionKind(r, group, desired.GetKind())
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
			}
			return r.Client.Create(context.TODO(), desired)
		}
		return err
	}

	current.SetLabels(desired.GetLabels())
	current.SetAnnotations(desired.GetAnnotations())
	current.Object["spec"] = desired.Object["spec"]
	return r.Client.Update(context.TODO(), current)
}

func deleteGatewayResource(r *utils.ComponentReconciler, desired *unstructured.Unstructured, group string) error {
	gvk, ok := resolveGroupVersionKind(r, group, desired.GetKind())
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
	return r.Client.Delete(context.TODO(), current)
}

func resolveGroupVersionKind(r *utils.ComponentReconciler, group, kind string) (schema.GroupVersionKind, bool) {
	if apiVersion, ok := r.GetApiVersionForKind(group, kind); ok {
		gv, err := schema.ParseGroupVersion(apiVersion)
		if err == nil {
			return gv.WithKind(kind), true
		}
	}
	return schema.GroupVersionKind{}, false
}

func buildHTTPRoute(cfg GatewayRouteConfig, routeCfg *monv1.GatewayHTTPRoute, hostnames []string) *unstructured.Unstructured {
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
		if parentRefs := buildParentRefs(routeCfg.ParentRefs, cfg.ParentRefs, group); len(parentRefs) > 0 {
			spec["parentRefs"] = parentRefs
		}
		if len(routeCfg.Rules) > 0 {
			spec["rules"] = buildHTTPRouteRules(cfg, routeCfg.Rules)
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
	}
}

func buildParentRefs(routeParentRefs, defaultParentRefs []monv1.GatewayParentRef, defaultGroup string) []interface{} {
	parentRefs := routeParentRefs
	if len(parentRefs) == 0 {
		parentRefs = defaultParentRefs
	}
	result := make([]interface{}, 0, len(parentRefs))
	for _, parentRef := range parentRefs {
		if parentRef.Name == "" {
			continue
		}
		entry := map[string]interface{}{
			"name": parentRef.Name,
		}
		entry["group"] = defaultGroup
		if parentRef.Group != "" {
			entry["group"] = parentRef.Group
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

func buildHTTPRouteRules(cfg GatewayRouteConfig, rules []monv1.GatewayHTTPRouteRule) []interface{} {
	result := make([]interface{}, 0, len(rules))
	for _, rule := range rules {
		item := map[string]interface{}{
			"backendRefs": []interface{}{
				map[string]interface{}{
					"name": cfg.ServiceName,
					"port": cfg.ServicePort,
				},
			},
		}
		if matches := rawJSONSliceToInterfaces(rule.Matches); len(matches) > 0 {
			item["matches"] = matches
		}
		if filters := rawJSONSliceToInterfaces(rule.Filters); len(filters) > 0 {
			item["filters"] = filters
		}
		result = append(result, item)
	}
	return result
}

func rawJSONSliceToInterfaces(items []monv1.GatewayJSON) []interface{} {
	result := make([]interface{}, 0, len(items))
	for _, item := range items {
		if len(item.Raw) == 0 {
			continue
		}
		var parsed interface{}
		if err := json.Unmarshal(item.Raw, &parsed); err != nil {
			continue
		}
		result = append(result, parsed)
	}
	return result
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
	parentRefs := cfg.ParentRefs
	if routeCfg != nil && len(routeCfg.ParentRefs) > 0 {
		parentRefs = routeCfg.ParentRefs
	}
	for _, parentRef := range parentRefs {
		if parentRef.Group != "" {
			return parentRef.Group
		}
	}
	return gatewayAPIGroup
}

func validateGatewayParentRefs(r *utils.ComponentReconciler, cfg GatewayRouteConfig, routeCfg *monv1.GatewayHTTPRoute) error {
	parentRefs := cfg.ParentRefs
	if routeCfg != nil && len(routeCfg.ParentRefs) > 0 {
		parentRefs = routeCfg.ParentRefs
	}
	for _, parentRef := range parentRefs {
		if parentRef.Group == "" {
			continue
		}
		if !r.HasApiGroupKind(parentRef.Group, "Gateway") {
			return errors.NewBadRequest("Gateway API group for parentRef is not available: " + parentRef.Group)
		}
	}
	return nil
}
