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
