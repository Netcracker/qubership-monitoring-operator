package victoriametrics

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

	got := buildParentRefs(parentRefs, nil, "gateway.networking.k8s.io")

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
}

func TestRawJSONSliceToInterfaces(t *testing.T) {
	t.Parallel()

	items := []monv1.GatewayJSON{
		{Raw: []byte(`{"type":"PathPrefix","value":"/"}`)},
		{Raw: []byte(`invalid-json`)},
		{Raw: nil},
	}

	got := rawJSONSliceToInterfaces(items)
	if len(got) != 1 {
		t.Fatalf("expected 1 parsed item, got=%d", len(got))
	}
	parsed := got[0].(map[string]interface{})
	if parsed["type"] != "PathPrefix" {
		t.Fatalf("unexpected parsed content: %v", parsed)
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
			Filters: []monv1.GatewayJSON{{Raw: []byte(`{"type":"URLRewrite"}`)}},
		},
	}

	got := buildHTTPRouteRules(cfg, rules)
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
	if _, ok := rule["filters"]; !ok {
		t.Fatalf("expected filters to be set")
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

	obj := buildHTTPRoute(cfg, routeCfg, []string{"route.example.com"})
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
