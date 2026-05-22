package gateway

import (
	"reflect"
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestGetHTTPRouteHostnames(t *testing.T) {
	t.Parallel()

	cfg := GatewayRouteConfig{Host: "ingress.example.com"}
	routeCfg := &monv1.GatewayHTTPRoute{Hostnames: []string{"route.example.com"}}

	got := getHTTPRouteHostnames(cfg, routeCfg)
	want := []string{"route.example.com"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected route hostnames to win, got=%v want=%v", got, want)
	}

	routeCfg.Hostnames = nil
	got = getHTTPRouteHostnames(cfg, routeCfg)
	want = []string{"ingress.example.com"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("expected ingress hostname fallback, got=%v want=%v", got, want)
	}

	cfg.Host = ""
	got = getHTTPRouteHostnames(cfg, routeCfg)
	if got != nil {
		t.Fatalf("expected nil hostnames when none set, got=%v", got)
	}
}

func TestBuildParentRefs(t *testing.T) {
	t.Parallel()

	parentRefs := []monv1.GatewayParentRef{
		{
			Name:        "gw",
			Namespace:   "infra",
			SectionName: "http",
		},
		{
			Name:  "gw2",
			Group: "example.io",
			Kind:  "Gateway",
		},
		{
			Name: "",
		},
	}

	got := buildParentRefs(parentRefs, nil)

	if len(got) != 2 {
		t.Fatalf("expected 2 parentRefs, got=%d", len(got))
	}

	first := got[0].(map[string]interface{})
	if first["name"] != "gw" {
		t.Fatalf("expected name gw, got=%v", first["name"])
	}
	if first["group"] != "gateway.networking.k8s.io" {
		t.Fatalf("expected default group, got=%v", first["group"])
	}
	if first["kind"] != "Gateway" {
		t.Fatalf("expected default kind Gateway, got=%v", first["kind"])
	}
	if first["namespace"] != "infra" {
		t.Fatalf("expected namespace infra, got=%v", first["namespace"])
	}
	if first["sectionName"] != "http" {
		t.Fatalf("expected sectionName http, got=%v", first["sectionName"])
	}

	second := got[1].(map[string]interface{})
	if second["group"] != "example.io" {
		t.Fatalf("expected custom group, got=%v", second["group"])
	}
	if second["kind"] != "Gateway" {
		t.Fatalf("expected kind Gateway, got=%v", second["kind"])
	}
}

func TestGatewayAPIGroupForConfig(t *testing.T) {
	t.Parallel()

	cfg := GatewayRouteConfig{
		ParentRefs: []monv1.GatewayParentRef{{Group: "example.io"}},
	}
	got := gatewayAPIGroupForConfig(cfg, nil)
	if got != "example.io" {
		t.Fatalf("expected group example.io, got=%s", got)
	}

	cfg.ParentRefs = nil
	got = gatewayAPIGroupForConfig(cfg, nil)
	if got != gatewayAPIGroup {
		t.Fatalf("expected default group %s, got=%s", gatewayAPIGroup, got)
	}

	routeCfg := &monv1.GatewayHTTPRoute{
		ParentRefs: []monv1.GatewayParentRef{{Group: "route.example.io"}},
	}
	cfg.ParentRefs = []monv1.GatewayParentRef{{Group: "default.example.io"}}
	got = gatewayAPIGroupForConfig(cfg, routeCfg)
	if got != "route.example.io" {
		t.Fatalf("expected route parentRefs to override defaults, got=%s", got)
	}
}

func TestBuildParentRefsDefaultsEmptyGroupToGatewayAPIGroup(t *testing.T) {
	t.Parallel()

	got := buildParentRefs(
		[]monv1.GatewayParentRef{
			{Name: "gw", Group: "example.io"},
			{Name: "default-gw"},
		},
		nil,
	)

	if len(got) != 2 {
		t.Fatalf("expected 2 parentRefs, got=%d", len(got))
	}

	first := got[0].(map[string]interface{})
	if first["group"] != "example.io" {
		t.Fatalf("expected explicit group example.io, got=%v", first["group"])
	}

	second := got[1].(map[string]interface{})
	if second["group"] != gatewayAPIGroup {
		t.Fatalf("expected empty group to default to %s, got=%v", gatewayAPIGroup, second["group"])
	}
}

func TestGatewayParentRefsForConfig(t *testing.T) {
	t.Parallel()

	routeParentRefs := []monv1.GatewayParentRef{{Name: "route-gw"}}
	defaultParentRefs := []monv1.GatewayParentRef{{Name: "default-gw"}}

	got := gatewayParentRefsForConfig(routeParentRefs, defaultParentRefs)
	if !reflect.DeepEqual(got, routeParentRefs) {
		t.Fatalf("expected route parentRefs to win, got=%v", got)
	}

	got = gatewayParentRefsForConfig(nil, defaultParentRefs)
	if !reflect.DeepEqual(got, defaultParentRefs) {
		t.Fatalf("expected default parentRefs fallback, got=%v", got)
	}
}

func TestHasHTTPRouteParentRefs(t *testing.T) {
	t.Parallel()

	cfg := GatewayRouteConfig{}
	routeCfg := &monv1.GatewayHTTPRoute{}
	if hasHTTPRouteParentRefs(cfg, routeCfg) {
		t.Fatalf("expected no parentRefs when neither route nor defaults configure a named parent")
	}

	cfg.ParentRefs = []monv1.GatewayParentRef{{Name: ""}}
	if hasHTTPRouteParentRefs(cfg, routeCfg) {
		t.Fatalf("expected empty parentRef name to be ignored")
	}

	cfg.ParentRefs = []monv1.GatewayParentRef{{Name: "default-gw"}}
	if !hasHTTPRouteParentRefs(cfg, routeCfg) {
		t.Fatalf("expected default parentRefs to be detected")
	}

	routeCfg.ParentRefs = []monv1.GatewayParentRef{{Name: "route-gw"}}
	cfg.ParentRefs = nil
	if !hasHTTPRouteParentRefs(cfg, routeCfg) {
		t.Fatalf("expected route parentRefs to be detected")
	}
}

func TestValidateGatewayParentRefsRejectsMixedGroups(t *testing.T) {
	t.Parallel()

	err := validateGatewayParentRefGroups([]monv1.GatewayParentRef{
		{Name: "custom-gw", Group: "example.io"},
		{Name: "default-gw"},
	})
	if err == nil {
		t.Fatalf("expected mixed groups to be rejected")
	}

	err = validateGatewayParentRefGroups([]monv1.GatewayParentRef{
		{Name: "custom-gw", Group: "example.io"},
		{Name: "custom-gw-2", Group: "example.io"},
	})
	if err != nil {
		t.Fatalf("expected matching custom groups to pass, got=%v", err)
	}
}

func TestRawJSONSliceToInterfaces(t *testing.T) {
	t.Parallel()

	items := []monv1.GatewayJSON{
		{Raw: []byte(`{"type":"PathPrefix","value":"/"}`)},
		{Raw: nil},
	}

	got, err := rawJSONSliceToInterfaces(items)
	if err != nil {
		t.Fatalf("expected valid raw JSON to parse, got %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 parsed item, got=%d", len(got))
	}
	parsed := got[0].(map[string]interface{})
	if parsed["type"] != "PathPrefix" {
		t.Fatalf("unexpected parsed content: %v", parsed)
	}
}

func TestRawJSONSliceToInterfacesRejectsInvalidJSON(t *testing.T) {
	t.Parallel()

	_, err := rawJSONSliceToInterfaces([]monv1.GatewayJSON{
		{Raw: []byte(`invalid-json`)},
	})
	if err == nil {
		t.Fatalf("expected invalid raw JSON to be rejected")
	}
}

func TestBuildHTTPRouteRules(t *testing.T) {
	t.Parallel()

	cfg := GatewayRouteConfig{
		ServiceName: "svc",
		ServicePort: 8080,
	}
	rules := []monv1.GatewayHTTPRouteRule{
		{
			Matches: []monv1.GatewayJSON{{Raw: []byte(`{"path":{"type":"PathPrefix","value":"/a"}}`)}},
		},
	}

	got, err := buildHTTPRouteRules(cfg, rules)
	if err != nil {
		t.Fatalf("expected valid rules, got %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 rule, got=%d", len(got))
	}
	rule := got[0].(map[string]interface{})
	backendRefs := rule["backendRefs"].([]interface{})
	if len(backendRefs) != 1 {
		t.Fatalf("expected 1 backendRef, got=%d", len(backendRefs))
	}
	backend := backendRefs[0].(map[string]interface{})
	if backend["name"] != "svc" || backend["port"] != int32(8080) {
		t.Fatalf("unexpected backendRef: %v", backend)
	}
	if _, ok := rule["matches"]; !ok {
		t.Fatalf("expected matches to be set")
	}
}

func TestBuildHTTPRouteRulesOmitsBackendRefsForTerminalFilters(t *testing.T) {
	t.Parallel()

	cfg := GatewayRouteConfig{
		ServiceName: "svc",
		ServicePort: 8080,
	}

	rules := []monv1.GatewayHTTPRouteRule{
		{
			Matches: []monv1.GatewayJSON{{Raw: []byte(`{"path":{"type":"PathPrefix","value":"/a"}}`)}},
			Filters: []monv1.GatewayJSON{{Raw: []byte(`{"type":"RequestRedirect"}`)}},
		},
	}
	got, err := buildHTTPRouteRules(cfg, rules)
	if err != nil {
		t.Fatalf("expected valid rules, got %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 rule, got=%d", len(got))
	}
	rule := got[0].(map[string]interface{})
	if _, ok := rule["backendRefs"]; ok {
		t.Fatalf("expected backendRefs to be absent for RequestRedirect filter")
	}
	if _, ok := rule["matches"]; !ok {
		t.Fatalf("expected matches to be set")
	}
	if _, ok := rule["filters"]; !ok {
		t.Fatalf("expected filters to be set")
	}
}

func TestBuildHTTPRouteRulesKeepsBackendRefsForURLRewrite(t *testing.T) {
	t.Parallel()

	cfg := GatewayRouteConfig{
		ServiceName: "svc",
		ServicePort: 8080,
	}
	rules := []monv1.GatewayHTTPRouteRule{
		{
			Matches: []monv1.GatewayJSON{{Raw: []byte(`{"path":{"type":"PathPrefix","value":"/a"}}`)}},
			Filters: []monv1.GatewayJSON{{Raw: []byte(`{"type":"URLRewrite"}`)}},
		},
	}

	got, err := buildHTTPRouteRules(cfg, rules)
	if err != nil {
		t.Fatalf("expected valid rules, got %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 rule, got=%d", len(got))
	}
	rule := got[0].(map[string]interface{})
	backendRefs, ok := rule["backendRefs"].([]interface{})
	if !ok || len(backendRefs) != 1 {
		t.Fatalf("expected backendRefs to be present for URLRewrite filter, got=%v", rule["backendRefs"])
	}
}

func TestBuildHTTPRouteRulesRejectsInvalidJSON(t *testing.T) {
	t.Parallel()

	cfg := GatewayRouteConfig{
		ServiceName: "svc",
		ServicePort: 8080,
	}
	_, err := buildHTTPRouteRules(cfg, []monv1.GatewayHTTPRouteRule{
		{
			Matches: []monv1.GatewayJSON{{Raw: []byte(`invalid-json`)}},
		},
	})
	if err == nil {
		t.Fatalf("expected invalid rule JSON to be rejected")
	}
}

func TestBuildHTTPRoute(t *testing.T) {
	t.Parallel()

	cfg := GatewayRouteConfig{
		NamePrefix:  "ns-svc",
		Namespace:   "ns",
		ServiceName: "svc",
		ServicePort: 80,
		Labels:      map[string]string{"app": "test"},
		ParentRefs:  []monv1.GatewayParentRef{{Name: "gw", Namespace: "infra"}},
	}
	routeCfg := &monv1.GatewayHTTPRoute{
		Hostnames: []string{"route.example.com"},
	}

	obj, err := buildHTTPRoute(cfg, routeCfg, []string{"route.example.com"})
	if err != nil {
		t.Fatalf("expected valid HTTPRoute, got %v", err)
	}
	if obj.GetKind() != httpRouteKind {
		t.Fatalf("expected kind %s, got=%s", httpRouteKind, obj.GetKind())
	}
	if obj.GetName() != "ns-svc"+httpRouteNameSuffix {
		t.Fatalf("unexpected name: %s", obj.GetName())
	}

	spec, ok := obj.Object["spec"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected spec map")
	}
	hostnames := spec["hostnames"].([]interface{})
	if len(hostnames) != 1 || hostnames[0] != "route.example.com" {
		t.Fatalf("unexpected hostnames: %v", hostnames)
	}
	parentRefs := spec["parentRefs"].([]interface{})
	if len(parentRefs) != 1 {
		t.Fatalf("expected 1 parentRef, got=%d", len(parentRefs))
	}
	parent := parentRefs[0].(map[string]interface{})
	if parent["name"] != "gw" || parent["namespace"] != "infra" {
		t.Fatalf("unexpected parentRef: %v", parent)
	}

	// Ensure object is a valid unstructured
	if _, ok := any(obj).(*unstructured.Unstructured); !ok {
		t.Fatalf("expected unstructured object")
	}
}

func TestHTTPRouteParentStatusWarningsNoStatus(t *testing.T) {
	t.Parallel()

	route := &unstructured.Unstructured{Object: map[string]interface{}{
		"kind": httpRouteKind,
	}}

	got := httpRouteParentStatusWarnings(route)
	if len(got) != 0 {
		t.Fatalf("expected no warnings without status, got=%v", got)
	}
}

func TestHTTPRouteParentStatusWarningsEmptyParents(t *testing.T) {
	t.Parallel()

	route := &unstructured.Unstructured{Object: map[string]interface{}{
		"kind": httpRouteKind,
		"status": map[string]interface{}{
			"parents": []interface{}{},
		},
	}}

	got := httpRouteParentStatusWarnings(route)
	if len(got) != 1 {
		t.Fatalf("expected one warning, got=%v", got)
	}
	if got[0].Reason != "NoParents" {
		t.Fatalf("expected NoParents warning, got=%v", got[0])
	}
}

func TestHTTPRouteParentStatusWarningsAccepted(t *testing.T) {
	t.Parallel()

	route := &unstructured.Unstructured{Object: map[string]interface{}{
		"kind": httpRouteKind,
		"status": map[string]interface{}{
			"parents": []interface{}{
				map[string]interface{}{
					"parentRef": map[string]interface{}{"name": "gw"},
					"conditions": []interface{}{
						map[string]interface{}{
							"type":   "Accepted",
							"status": "True",
						},
					},
				},
			},
		},
	}}

	got := httpRouteParentStatusWarnings(route)
	if len(got) != 0 {
		t.Fatalf("expected no warnings for accepted route, got=%v", got)
	}
}

func TestHTTPRouteParentStatusWarningsAcceptedAndResolvedRefs(t *testing.T) {
	t.Parallel()

	route := &unstructured.Unstructured{Object: map[string]interface{}{
		"kind": httpRouteKind,
		"status": map[string]interface{}{
			"parents": []interface{}{
				map[string]interface{}{
					"parentRef": map[string]interface{}{"name": "gw"},
					"conditions": []interface{}{
						map[string]interface{}{
							"type":   "Accepted",
							"status": "True",
							"reason": "Accepted",
						},
						map[string]interface{}{
							"type":   "ResolvedRefs",
							"status": "True",
							"reason": "ResolvedRefs",
						},
					},
				},
			},
		},
	}}

	got := httpRouteParentStatusWarnings(route)
	if len(got) != 0 {
		t.Fatalf("expected no warnings for accepted route with resolved refs, got=%v", got)
	}
}

func TestHTTPRouteParentStatusWarningsRejected(t *testing.T) {
	t.Parallel()

	route := &unstructured.Unstructured{Object: map[string]interface{}{
		"kind": httpRouteKind,
		"status": map[string]interface{}{
			"parents": []interface{}{
				map[string]interface{}{
					"parentRef": map[string]interface{}{
						"group":       gatewayAPIGroup,
						"kind":        "Gateway",
						"name":        "missing-gw",
						"namespace":   "infra",
						"sectionName": "http",
					},
					"conditions": []interface{}{
						map[string]interface{}{
							"type":    "Accepted",
							"status":  "False",
							"reason":  "ParentNotFound",
							"message": "Gateway does not exist",
						},
					},
				},
			},
		},
	}}

	got := httpRouteParentStatusWarnings(route)
	if len(got) != 1 {
		t.Fatalf("expected one warning, got=%v", got)
	}
	if got[0].Reason != "ParentNotFound" || got[0].Message != "Gateway does not exist" {
		t.Fatalf("unexpected warning: %v", got[0])
	}
	if got[0].ParentName != "missing-gw" || got[0].ParentNamespace != "infra" || got[0].SectionName != "http" {
		t.Fatalf("expected parentRef details, got=%v", got[0])
	}
}

func TestHTTPRouteParentStatusWarningsUnresolvedRefs(t *testing.T) {
	t.Parallel()

	route := &unstructured.Unstructured{Object: map[string]interface{}{
		"kind": httpRouteKind,
		"status": map[string]interface{}{
			"parents": []interface{}{
				map[string]interface{}{
					"parentRef": map[string]interface{}{"name": "gw"},
					"conditions": []interface{}{
						map[string]interface{}{
							"type":   "Accepted",
							"status": "True",
							"reason": "Accepted",
						},
						map[string]interface{}{
							"type":    "ResolvedRefs",
							"status":  "False",
							"reason":  "BackendNotFound",
							"message": "backend service does not exist",
						},
					},
				},
			},
		},
	}}

	got := httpRouteParentStatusWarnings(route)
	if len(got) != 1 {
		t.Fatalf("expected one warning, got=%v", got)
	}
	if got[0].Reason != "BackendNotFound" || got[0].Message != "backend service does not exist" {
		t.Fatalf("unexpected warning: %v", got[0])
	}
}
