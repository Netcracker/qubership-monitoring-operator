package gateway

import (
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestApplyGatewayResourceSkipsNoOpUpdate(t *testing.T) {
	cr := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "platformmonitoring",
			Namespace: "monitoring",
			UID:       types.UID("platform-monitoring-uid"),
		},
	}
	cfg := GatewayRouteConfig{
		NamePrefix:  "grafana",
		Namespace:   "monitoring",
		Host:        "grafana.example.com",
		ServiceName: "grafana",
		ServicePort: 3000,
		Labels:      map[string]string{"app.kubernetes.io/name": "grafana"},
	}
	desired, err := buildHTTPRoute(cfg, nil, []string{cfg.Host})
	require.NoError(t, err)
	gvk := schema.GroupVersionKind{Group: gatewayAPIGroup, Version: "v1", Kind: httpRouteKind}
	desired.SetGroupVersionKind(gvk)
	raw, err := desired.MarshalJSON()
	require.NoError(t, err)
	require.NoError(t, desired.UnmarshalJSON(raw))
	desired.SetGroupVersionKind(gvk)
	desired.SetResourceVersion("1")

	scheme := runtime.NewScheme()
	require.NoError(t, monv1.AddToScheme(scheme))
	scheme.AddKnownTypeWithName(gvk, &unstructured.Unstructured{})
	scheme.AddKnownTypeWithName(gvk.GroupVersion().WithKind(httpRouteKind+"List"), &unstructured.UnstructuredList{})
	metav1.AddToGroupVersion(scheme, gvk.GroupVersion())

	r := &utils.ComponentReconciler{
		Client: fake.NewClientBuilder().WithScheme(scheme).WithObjects(cr, desired.DeepCopy()).Build(),
		Scheme: scheme,
		Log:    utils.Logger("gateway_skip_test"),
	}
	resolver := &gatewayAPIResolver{
		r:      r,
		loaded: true,
		apiLists: []*metav1.APIResourceList{{
			GroupVersion: gvk.GroupVersion().String(),
			APIResources: []metav1.APIResource{{Kind: httpRouteKind, Name: "httproutes"}},
		}},
	}

	require.NoError(t, applyGatewayResource(r, resolver, cr, desired.DeepCopy(), gatewayAPIGroup))

	got := &unstructured.Unstructured{}
	got.SetGroupVersionKind(gvk)
	require.NoError(t, r.Client.Get(t.Context(), client.ObjectKeyFromObject(desired), got))
	assert.Equal(t, "1", got.GetResourceVersion(), "no-op reconcile must not bump resourceVersion")
}

func TestApplyGatewayResourceUpdatesChangedSpec(t *testing.T) {
	cr := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "platformmonitoring",
			Namespace: "monitoring",
			UID:       types.UID("platform-monitoring-uid"),
		},
	}
	cfg := GatewayRouteConfig{
		NamePrefix:  "grafana",
		Namespace:   "monitoring",
		Host:        "grafana.example.com",
		ServiceName: "grafana",
		ServicePort: 3000,
		Labels:      map[string]string{"app.kubernetes.io/name": "grafana"},
	}
	desired, err := buildHTTPRoute(cfg, nil, []string{cfg.Host})
	require.NoError(t, err)
	gvk := schema.GroupVersionKind{Group: gatewayAPIGroup, Version: "v1", Kind: httpRouteKind}
	desired.SetGroupVersionKind(gvk)
	raw, err := desired.MarshalJSON()
	require.NoError(t, err)
	require.NoError(t, desired.UnmarshalJSON(raw))
	desired.SetGroupVersionKind(gvk)

	live := desired.DeepCopy()
	live.SetResourceVersion("1")
	require.NoError(t, unstructured.SetNestedStringSlice(live.Object, []string{"stale.example.com"}, "spec", "hostnames"))

	scheme := runtime.NewScheme()
	require.NoError(t, monv1.AddToScheme(scheme))
	scheme.AddKnownTypeWithName(gvk, &unstructured.Unstructured{})
	scheme.AddKnownTypeWithName(gvk.GroupVersion().WithKind(httpRouteKind+"List"), &unstructured.UnstructuredList{})
	metav1.AddToGroupVersion(scheme, gvk.GroupVersion())

	r := &utils.ComponentReconciler{
		Client: fake.NewClientBuilder().WithScheme(scheme).WithObjects(cr, live).Build(),
		Scheme: scheme,
		Log:    utils.Logger("gateway_skip_test"),
	}
	resolver := &gatewayAPIResolver{
		r:      r,
		loaded: true,
		apiLists: []*metav1.APIResourceList{{
			GroupVersion: gvk.GroupVersion().String(),
			APIResources: []metav1.APIResource{{Kind: httpRouteKind, Name: "httproutes"}},
		}},
	}

	require.NoError(t, applyGatewayResource(r, resolver, cr, desired.DeepCopy(), gatewayAPIGroup))

	got := &unstructured.Unstructured{}
	got.SetGroupVersionKind(gvk)
	require.NoError(t, r.Client.Get(t.Context(), client.ObjectKeyFromObject(desired), got))
	assert.NotEqual(t, "1", got.GetResourceVersion(), "spec change must trigger Update")
	hostnames, found, err := unstructured.NestedStringSlice(got.Object, "spec", "hostnames")
	require.NoError(t, err)
	require.True(t, found)
	assert.Equal(t, []string{cfg.Host}, hostnames)
}

func TestApplyGatewayResourceSkipsAPIDefaultedSpec(t *testing.T) {
	cr, cfg := skipTestGatewayInput()
	desired, err := buildHTTPRoute(cfg, nil, []string{cfg.Host})
	require.NoError(t, err)
	live := skipTestDecodedHTTPRoute(t, desired)
	live.SetResourceVersion("1")
	// The standard HTTPRoute schema defaults these backend reference fields.
	live.Object["spec"].(map[string]interface{})["rules"] = []interface{}{
		map[string]interface{}{
			"matches": []interface{}{
				map[string]interface{}{"path": map[string]interface{}{"type": "PathPrefix", "value": "/"}},
			},
			"backendRefs": []interface{}{
				map[string]interface{}{
					"group":  "",
					"kind":   "Service",
					"name":   cfg.ServiceName,
					"port":   int64(cfg.ServicePort),
					"weight": int64(1),
				},
			},
		},
	}

	got := applySkipTestHTTPRoute(t, cr, desired, live)

	assert.Equal(t, "1", got.GetResourceVersion(), "API-defaulted HTTPRoute fields must not trigger Update")
	assert.Equal(t, cfg.Labels, got.GetLabels())
}

func TestApplyGatewayResourceUpdatesChangedBackendWeight(t *testing.T) {
	cr, cfg := skipTestGatewayInput()
	desired, err := buildHTTPRoute(cfg, nil, []string{cfg.Host})
	require.NoError(t, err)
	desired = skipTestDecodedHTTPRoute(t, desired)
	live := desired.DeepCopy()
	live.SetResourceVersion("1")
	backendRef := live.Object["spec"].(map[string]interface{})["rules"].([]interface{})[0].(map[string]interface{})["backendRefs"].([]interface{})[0].(map[string]interface{})
	backendRef["weight"] = int64(5)

	got := applySkipTestHTTPRoute(t, cr, desired, live)

	assert.NotEqual(t, "1", got.GetResourceVersion(), "a non-default weight must trigger Update")
	assert.Equal(t, cfg.Labels, got.GetLabels(), "an Update must keep the route labels")
}

func TestHTTPRouteSpecWithDefaults(t *testing.T) {
	spec := map[string]interface{}{
		"parentRefs": []interface{}{map[string]interface{}{"name": "gateway"}},
		"rules": []interface{}{
			map[string]interface{}{
				"matches": []interface{}{
					map[string]interface{}{
						"headers":     []interface{}{map[string]interface{}{"name": "x-env", "value": "dev"}},
						"queryParams": []interface{}{map[string]interface{}{"name": "debug", "value": "1"}},
					},
				},
				"filters": []interface{}{
					map[string]interface{}{"type": "RequestRedirect", "requestRedirect": map[string]interface{}{"scheme": "https"}},
				},
			},
			map[string]interface{}{
				"backendRefs": []interface{}{map[string]interface{}{"name": "grafana", "port": 3000, "kind": "Custom"}},
			},
		},
	}

	got, err := httpRouteSpecWithDefaults(spec)
	require.NoError(t, err)

	defaultPath := map[string]interface{}{"type": "PathPrefix", "value": "/"}
	assert.Equal(t, map[string]interface{}{
		"parentRefs": []interface{}{
			map[string]interface{}{"name": "gateway", "group": gatewayAPIGroup, "kind": "Gateway"},
		},
		"rules": []interface{}{
			map[string]interface{}{
				"matches": []interface{}{
					map[string]interface{}{
						"path":        defaultPath,
						"headers":     []interface{}{map[string]interface{}{"name": "x-env", "value": "dev", "type": "Exact"}},
						"queryParams": []interface{}{map[string]interface{}{"name": "debug", "value": "1", "type": "Exact"}},
					},
				},
				"filters": []interface{}{
					map[string]interface{}{
						"type":            "RequestRedirect",
						"requestRedirect": map[string]interface{}{"scheme": "https", "statusCode": float64(302)},
					},
				},
			},
			map[string]interface{}{
				"matches": []interface{}{map[string]interface{}{"path": defaultPath}},
				"backendRefs": []interface{}{
					map[string]interface{}{"name": "grafana", "port": float64(3000), "group": "", "kind": "Custom", "weight": float64(1)},
				},
			},
		},
	}, got)
	assert.NotContains(t, spec["parentRefs"].([]interface{})[0], "group", "the input spec must stay unchanged")
}

var skipTestHTTPRouteGVK = schema.GroupVersionKind{Group: gatewayAPIGroup, Version: "v1", Kind: httpRouteKind}

func skipTestGatewayInput() (*monv1.PlatformMonitoring, GatewayRouteConfig) {
	cr := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "platformmonitoring",
			Namespace: "monitoring",
			UID:       types.UID("platform-monitoring-uid"),
		},
	}
	cfg := GatewayRouteConfig{
		NamePrefix:  "grafana",
		Namespace:   "monitoring",
		Host:        "grafana.example.com",
		ServiceName: "grafana",
		ServicePort: 3000,
		Labels:      map[string]string{"app.kubernetes.io/name": "grafana"},
	}
	return cr, cfg
}

// skipTestDecodedHTTPRoute returns route as decoded from JSON, which is the
// form the API server returns and the only form the fake client can store.
func skipTestDecodedHTTPRoute(t *testing.T, route *unstructured.Unstructured) *unstructured.Unstructured {
	t.Helper()
	raw, err := route.MarshalJSON()
	require.NoError(t, err)
	decoded := &unstructured.Unstructured{}
	require.NoError(t, decoded.UnmarshalJSON(raw))
	decoded.SetGroupVersionKind(skipTestHTTPRouteGVK)
	return decoded
}

// applySkipTestHTTPRoute reconciles desired against a cluster that holds live
// and returns the stored route.
func applySkipTestHTTPRoute(t *testing.T, cr *monv1.PlatformMonitoring, desired, live *unstructured.Unstructured) *unstructured.Unstructured {
	t.Helper()
	gvk := skipTestHTTPRouteGVK
	scheme := runtime.NewScheme()
	require.NoError(t, monv1.AddToScheme(scheme))
	scheme.AddKnownTypeWithName(gvk, &unstructured.Unstructured{})
	scheme.AddKnownTypeWithName(gvk.GroupVersion().WithKind(httpRouteKind+"List"), &unstructured.UnstructuredList{})
	metav1.AddToGroupVersion(scheme, gvk.GroupVersion())

	r := &utils.ComponentReconciler{
		Client: fake.NewClientBuilder().WithScheme(scheme).WithObjects(cr, live).Build(),
		Scheme: scheme,
		Log:    utils.Logger("gateway_skip_test"),
	}
	resolver := &gatewayAPIResolver{
		r:      r,
		loaded: true,
		apiLists: []*metav1.APIResourceList{{
			GroupVersion: gvk.GroupVersion().String(),
			APIResources: []metav1.APIResource{{Kind: httpRouteKind, Name: "httproutes"}},
		}},
	}
	require.NoError(t, applyGatewayResource(r, resolver, cr, desired, gatewayAPIGroup))

	got := &unstructured.Unstructured{}
	got.SetGroupVersionKind(gvk)
	require.NoError(t, r.Client.Get(t.Context(), client.ObjectKeyFromObject(live), got))
	return got
}
