package grafana

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	grafv1 "github.com/grafana/grafana-operator/v5/api/v1beta1"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	discoveryfake "k8s.io/client-go/discovery/fake"
	"k8s.io/client-go/kubernetes"
	kubernetesfake "k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
)

func TestGrafanaJaegerDataSourceSingleService(t *testing.T) {
	cr := discoveredPlatformMonitoring("15s")
	service := jaegerService("tracing", "jaeger-query", 16686)

	datasources := grafanaJaegerDataSources(cr, []corev1.Service{service})

	require.Len(t, datasources, 1)
	assertDiscoveredDatasource(t, datasources[0], "platform-monitoring-jaeger-tracing.jaeger-query", jaegerDatasourceComponent)
	assert.Equal(t, "Jaeger", datasources[0].Spec.Datasource.Name)
	assert.Equal(t, "jaeger", datasources[0].Spec.Datasource.Type)
	assert.Equal(t, "proxy", datasources[0].Spec.Datasource.Access)
	assert.Equal(t, "http://jaeger-query.tracing.svc.cluster.local:16686", datasources[0].Spec.Datasource.URL)
	assert.Equal(t, false, *datasources[0].Spec.Datasource.IsDefault)
	assert.Equal(t, true, *datasources[0].Spec.Datasource.Editable)
	assert.JSONEq(t, `{"nodeGraph":{"enabled":true},"timeInterval":"15s","tlsSkipVerify":true}`, string(datasources[0].Spec.Datasource.JSONData))
	assert.Empty(t, datasources[0].Spec.Plugins)
}

func TestGrafanaJaegerDataSourceUsesDefaultIntervalWhenVmAgentIsNotInstalled(t *testing.T) {
	cr := discoveredPlatformMonitoring("15s")
	installed := false
	cr.Spec.Victoriametrics.VmAgent.Install = &installed

	datasources := grafanaJaegerDataSources(cr, []corev1.Service{jaegerService("tracing", "jaeger-query", 16686)})

	require.Len(t, datasources, 1)
	assert.JSONEq(t, `{"nodeGraph":{"enabled":true},"timeInterval":"30s","tlsSkipVerify":true}`, string(datasources[0].Spec.Datasource.JSONData))
}

func TestGrafanaJaegerDataSourcesUseServiceIdentityWhenSeveralExist(t *testing.T) {
	cr := discoveredPlatformMonitoring("30s")
	datasources := grafanaJaegerDataSources(cr, []corev1.Service{
		jaegerService("alpha", "query", 16686),
		jaegerService("beta", "query", 16687),
	})

	require.Len(t, datasources, 2)
	assert.Equal(t, "Jaeger alpha/query", datasources[0].Spec.Datasource.Name)
	assert.Equal(t, "Jaeger beta/query", datasources[1].Spec.Datasource.Name)
	assert.Equal(t, "platform-monitoring-jaeger-alpha.query", datasources[0].GetName())
	assert.Equal(t, "platform-monitoring-jaeger-beta.query", datasources[1].GetName())
}

func TestGrafanaJaegerDataSourcesUseDistinctNamesForHyphenatedServiceIdentities(t *testing.T) {
	cr := discoveredPlatformMonitoring("30s")
	datasources := grafanaJaegerDataSources(cr, []corev1.Service{
		jaegerService("alpha-beta", "query", 16686),
		jaegerService("alpha", "beta-query", 16687),
	})

	require.Len(t, datasources, 2)
	assert.Equal(t, "platform-monitoring-jaeger-alpha-beta.query", datasources[0].GetName())
	assert.Equal(t, "platform-monitoring-jaeger-alpha.beta-query", datasources[1].GetName())
	assert.NotEqual(t, datasources[0].GetName(), datasources[1].GetName())
}

func TestGrafanaJaegerDataSourceTruncatesKubernetesLabelValues(t *testing.T) {
	cr := discoveredPlatformMonitoring("30s")
	service := jaegerService("namespace-with-a-long-name", "service-with-a-long-name", 16686)

	datasources := grafanaJaegerDataSources(cr, []corev1.Service{service})

	require.Len(t, datasources, 1)
	assert.LessOrEqual(t, len(datasources[0].GetLabels()["name"]), 63)
	assert.LessOrEqual(t, len(datasources[0].GetLabels()["app.kubernetes.io/name"]), 63)
}

func TestGrafanaJaegerDataSourceTruncatesSeparatorFromKubernetesLabelValues(t *testing.T) {
	cr := discoveredPlatformMonitoring("30s")
	service := jaegerService(strings.Repeat("a", 35), "query", 16686)

	datasources := grafanaJaegerDataSources(cr, []corev1.Service{service})
	reconciler := newDiscoveredReconciler(t, "", kubernetesfake.NewSimpleClientset())

	require.Len(t, datasources, 1)
	require.NoError(t, reconciler.applyDiscoveredDataSource(cr, datasources[0]))
	stored := discoveredDatasourceObject(datasources[0].GetName(), jaegerDatasourceComponent)
	require.NoError(t, reconciler.GetResource(stored))
	for _, labelKey := range []string{"name", "app.kubernetes.io/name", "app.kubernetes.io/instance"} {
		labelValue := stored.GetLabels()[labelKey]
		assert.LessOrEqual(t, len(labelValue), 63)
		assert.NotEqual(t, ".", labelValue[len(labelValue)-1:])
	}
}

func TestGrafanaClickHouseDataSourceWithoutCredentials(t *testing.T) {
	cr := discoveredPlatformMonitoring("30s")
	datasources, err := grafanaClickHouseDataSources(cr, nil, []corev1.Service{clickHouseService("analytics")})

	require.NoError(t, err)
	require.Len(t, datasources, 1)
	assertDiscoveredDatasource(t, datasources[0], "platform-monitoring-clickhouse-analytics", clickHouseDatasourceComponent)
	assert.Equal(t, "ClickHouse", datasources[0].Spec.Datasource.Name)
	assert.Equal(t, clickHousePluginName, datasources[0].Spec.Datasource.Type)
	assert.Equal(t, "http://clickhouse-cluster.analytics.svc.cluster.local:8123", datasources[0].Spec.Datasource.URL)
	assert.Nil(t, datasources[0].Spec.Datasource.BasicAuth)
	assert.Empty(t, datasources[0].Spec.Datasource.BasicAuthUser)
	assert.Empty(t, datasources[0].Spec.Datasource.SecureJSONData)
	require.Len(t, datasources[0].Spec.Plugins, 1)
	assert.Equal(t, clickHousePluginName, datasources[0].Spec.Plugins[0].Name)
	assert.Equal(t, clickHousePluginVersion, datasources[0].Spec.Plugins[0].Version)
}

func TestGrafanaClickHouseDataSourceUsesCredentials(t *testing.T) {
	cr := discoveredPlatformMonitoring("30s")
	kubeClient := kubernetesfake.NewSimpleClientset(&corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: utils.ClickHouseSecret, Namespace: "analytics"},
		Data: map[string][]byte{
			"username": []byte("metrics"),
			"password": []byte("secret"),
		},
	})

	datasources, err := grafanaClickHouseDataSources(cr, kubeClient, []corev1.Service{clickHouseService("analytics")})

	require.NoError(t, err)
	require.Len(t, datasources, 1)
	assert.Equal(t, true, *datasources[0].Spec.Datasource.BasicAuth)
	assert.Equal(t, "metrics", datasources[0].Spec.Datasource.BasicAuthUser)
	secure := map[string]string{}
	require.NoError(t, json.Unmarshal(datasources[0].Spec.Datasource.SecureJSONData, &secure))
	assert.Equal(t, "secret", secure["basicAuthPassword"])
}

func TestGrafanaClickHouseDataSourceIgnoresIncompleteCredentials(t *testing.T) {
	cr := discoveredPlatformMonitoring("30s")
	kubeClient := kubernetesfake.NewSimpleClientset(&corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: utils.ClickHouseSecret, Namespace: "analytics"},
		Data:       map[string][]byte{"username": []byte("metrics")},
	})

	datasources, err := grafanaClickHouseDataSources(cr, kubeClient, []corev1.Service{clickHouseService("analytics")})

	require.NoError(t, err)
	require.Len(t, datasources, 1)
	assert.Nil(t, datasources[0].Spec.Datasource.BasicAuth)
	assert.Empty(t, datasources[0].Spec.Datasource.BasicAuthUser)
	assert.Empty(t, datasources[0].Spec.Datasource.SecureJSONData)
}

func TestGrafanaClickHouseDataSourcesNameByNamespace(t *testing.T) {
	cr := discoveredPlatformMonitoring("30s")
	datasources, err := grafanaClickHouseDataSources(cr, nil, []corev1.Service{
		clickHouseService("alpha"),
		clickHouseService("beta"),
	})

	require.NoError(t, err)
	require.Len(t, datasources, 2)
	assert.Equal(t, "ClickHouse_alpha", datasources[0].Spec.Datasource.Name)
	assert.Equal(t, "ClickHouse_beta", datasources[1].Spec.Datasource.Name)
}

func TestHandleJaegerDataSourcesAdoptsLegacyUID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, "/api/datasources", request.URL.Path)
		writer.Header().Set("Content-Type", "application/json")
		_, err := writer.Write([]byte(`[{"name":"Jaeger","uid":"legacy-jaeger"}]`))
		assert.NoError(t, err)
	}))
	t.Cleanup(server.Close)

	withPrivilegedRights(t)
	cr := discoveredIntegrationPlatformMonitoring(true, false)
	kubeClient := kubernetesfake.NewSimpleClientset(
		&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "tracing"}},
		jaegerServicePtr("tracing", "jaeger-query", 16686),
	)
	reconciler := newDiscoveredReconciler(t, server.URL, kubeClient, legacyCombinedDatasource())

	require.NoError(t, reconciler.handleJaegerDataSources(cr))

	created := &grafv1.GrafanaDatasource{}
	created.SetName("platform-monitoring-jaeger-tracing.jaeger-query")
	created.SetNamespace("monitoring")
	require.NoError(t, reconciler.GetResource(created))
	assert.Equal(t, "legacy-jaeger", created.Spec.CustomUID)

	listed := &grafv1.GrafanaDatasourceList{}
	require.NoError(t, reconciler.Client.List(context.Background(), listed, client.InNamespace("monitoring")))
	assert.Len(t, listed.Items, 1)
}

func TestHandleClickHouseDataSourcesAdoptsLegacyUID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, "/api/datasources", request.URL.Path)
		writer.Header().Set("Content-Type", "application/json")
		_, err := writer.Write([]byte(`[{"name":"ClickHouse","uid":"legacy-clickhouse"}]`))
		assert.NoError(t, err)
	}))
	t.Cleanup(server.Close)

	withPrivilegedRights(t)
	cr := discoveredIntegrationPlatformMonitoring(false, true)
	kubeClient := kubernetesfake.NewSimpleClientset(
		&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "analytics"}},
		clickHouseServicePtr("analytics"),
	)
	reconciler := newDiscoveredReconciler(t, server.URL, kubeClient, legacyCombinedDatasource())

	require.NoError(t, reconciler.handleClickHouseDataSources(cr))

	created := discoveredDatasourceObject("platform-monitoring-clickhouse-analytics", clickHouseDatasourceComponent)
	require.NoError(t, reconciler.GetResource(created))
	assert.Equal(t, "legacy-clickhouse", created.Spec.CustomUID)
}

func TestHandleJaegerDataSourcesDeletesStaleAndKeepsOtherDatasources(t *testing.T) {
	withPrivilegedRights(t)
	cr := discoveredIntegrationPlatformMonitoring(true, false)
	kubeClient := kubernetesfake.NewSimpleClientset(
		&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "tracing"}},
		jaegerServicePtr("tracing", "jaeger-query", 16686),
	)
	stale := discoveredDatasourceObject("platform-monitoring-jaeger-old-query", jaegerDatasourceComponent)
	prometheus := discoveredDatasourceObject("platform-monitoring-prometheus", "grafana")
	promxy := discoveredDatasourceObject("platform-monitoring-promxy", "grafana")
	reconciler := newDiscoveredReconciler(t, "", kubeClient, stale, prometheus, promxy)

	require.NoError(t, reconciler.handleJaegerDataSources(cr))

	assert.Error(t, reconciler.GetResource(discoveredDatasourceObject("platform-monitoring-jaeger-old-query", jaegerDatasourceComponent)))
	assert.NoError(t, reconciler.GetResource(discoveredDatasourceObject("platform-monitoring-jaeger-tracing.jaeger-query", jaegerDatasourceComponent)))
	assert.NoError(t, reconciler.GetResource(prometheus))
	assert.NoError(t, reconciler.GetResource(promxy))
}

func TestHandleJaegerDataSourcesFlagOffDeletesOwnedDatasource(t *testing.T) {
	withPrivilegedRights(t)
	cr := discoveredIntegrationPlatformMonitoring(false, false)
	owned := discoveredDatasourceObject("platform-monitoring-jaeger-tracing.jaeger-query", jaegerDatasourceComponent)
	prometheus := discoveredDatasourceObject("platform-monitoring-prometheus", "grafana")
	reconciler := newDiscoveredReconciler(t, "", kubernetesfake.NewSimpleClientset(), owned, prometheus)

	require.NoError(t, reconciler.handleJaegerDataSources(cr))

	assert.Error(t, reconciler.GetResource(discoveredDatasourceObject("platform-monitoring-jaeger-tracing.jaeger-query", jaegerDatasourceComponent)))
	assert.NoError(t, reconciler.GetResource(prometheus))
}

func TestHandleJaegerDataSourcesPreservesUIDWhenUpdatingURL(t *testing.T) {
	withPrivilegedRights(t)
	cr := discoveredIntegrationPlatformMonitoring(true, false)
	kubeClient := kubernetesfake.NewSimpleClientset(
		&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "tracing"}},
		jaegerServicePtr("tracing", "jaeger-query", 16686),
	)
	existing := discoveredDatasourceObject("platform-monitoring-jaeger-tracing.jaeger-query", jaegerDatasourceComponent)
	existing.Spec.CustomUID = "stable-uid"
	existing.Spec.InstanceSelector = &metav1.LabelSelector{MatchLabels: map[string]string{
		"app.kubernetes.io/component": "existing-grafana",
	}}
	existing.Spec.Datasource = &grafv1.GrafanaDatasourceInternal{Name: "old", URL: "http://old"}
	reconciler := newDiscoveredReconciler(t, "", kubeClient, existing)

	require.NoError(t, reconciler.handleJaegerDataSources(cr))
	require.NoError(t, reconciler.handleJaegerDataSources(cr))

	updated := discoveredDatasourceObject(existing.GetName(), jaegerDatasourceComponent)
	require.NoError(t, reconciler.GetResource(updated))
	assert.Equal(t, "http://jaeger-query.tracing.svc.cluster.local:16686", updated.Spec.Datasource.URL)
	assert.Equal(t, "stable-uid", updated.Spec.CustomUID)
	assert.Equal(t, map[string]string{"app.kubernetes.io/component": "existing-grafana"}, updated.Spec.InstanceSelector.MatchLabels)
	assert.Equal(t, "platform-monitoring-jaeger-tracing.jaeger-query-monitoring", updated.GetLabels()["app.kubernetes.io/instance"])
	assert.Equal(t, "12.4.3", updated.GetLabels()["app.kubernetes.io/version"])
}

func TestHandleJaegerDataSourcesGetErrorLeavesExistingDatasource(t *testing.T) {
	withPrivilegedRights(t)
	cr := discoveredIntegrationPlatformMonitoring(true, false)
	kubeClient := kubernetesfake.NewSimpleClientset(
		&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "tracing"}},
		jaegerServicePtr("tracing", "jaeger-query", 16686),
	)
	existing := discoveredDatasourceObject("platform-monitoring-jaeger-tracing.jaeger-query", jaegerDatasourceComponent)
	existing.Spec.Datasource = &grafv1.GrafanaDatasourceInternal{URL: "http://old"}
	reads := 0
	reconciler := newDiscoveredReconcilerIntercepting(t, "", kubeClient, interceptor.Funcs{
		Get: func(ctx context.Context, kubeClient client.WithWatch, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
			if _, ok := obj.(*grafv1.GrafanaDatasource); ok && reads == 0 {
				reads++
				return errors.New("datasource get failed")
			}
			return kubeClient.Get(ctx, key, obj, opts...)
		},
	}, existing)

	err := reconciler.handleJaegerDataSources(cr)

	require.Error(t, err)
	assert.ErrorContains(t, err, "datasource get failed")
	stored := discoveredDatasourceObject(existing.GetName(), jaegerDatasourceComponent)
	require.NoError(t, reconciler.GetResource(stored))
	assert.Equal(t, "http://old", stored.Spec.Datasource.URL)
}

func TestHandleJaegerDataSourcesLegacyLookupErrorDoesNotCreate(t *testing.T) {
	withPrivilegedRights(t)
	cr := discoveredIntegrationPlatformMonitoring(true, false)
	kubeClient := kubernetesfake.NewSimpleClientset(
		&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "tracing"}},
		jaegerServicePtr("tracing", "jaeger-query", 16686),
	)
	reconciler := newDiscoveredReconcilerIntercepting(t, "", kubeClient, interceptor.Funcs{
		Get: func(ctx context.Context, kubeClient client.WithWatch, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
			if _, ok := obj.(*unstructured.Unstructured); ok && key.Name == legacyCombinedDatasourceName {
				return errors.New("legacy lookup failed")
			}
			return kubeClient.Get(ctx, key, obj, opts...)
		},
	})

	err := reconciler.handleJaegerDataSources(cr)

	require.Error(t, err)
	assert.ErrorContains(t, err, "legacy lookup failed")
	assert.Error(t, reconciler.GetResource(discoveredDatasourceObject("platform-monitoring-jaeger-tracing.jaeger-query", jaegerDatasourceComponent)))
}

func TestHandleJaegerDataSourcesCleanupListErrorLeavesStaleDatasource(t *testing.T) {
	withPrivilegedRights(t)
	cr := discoveredIntegrationPlatformMonitoring(true, false)
	kubeClient := kubernetesfake.NewSimpleClientset(
		&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "tracing"}},
		jaegerServicePtr("tracing", "jaeger-query", 16686),
	)
	stale := discoveredDatasourceObject("platform-monitoring-jaeger-old-query", jaegerDatasourceComponent)
	reconciler := newDiscoveredReconcilerIntercepting(t, "", kubeClient, interceptor.Funcs{
		List: func(ctx context.Context, kubeClient client.WithWatch, list client.ObjectList, opts ...client.ListOption) error {
			if _, ok := list.(*grafv1.GrafanaDatasourceList); ok {
				return errors.New("datasource list failed")
			}
			return kubeClient.List(ctx, list, opts...)
		},
	}, stale)

	err := reconciler.handleJaegerDataSources(cr)

	require.Error(t, err)
	assert.ErrorContains(t, err, "datasource list failed")
	assert.NoError(t, reconciler.GetResource(discoveredDatasourceObject(stale.GetName(), jaegerDatasourceComponent)))
}

func TestHandleJaegerDataSourcesDeleteErrorLeavesStaleDatasource(t *testing.T) {
	withPrivilegedRights(t)
	cr := discoveredIntegrationPlatformMonitoring(true, false)
	kubeClient := kubernetesfake.NewSimpleClientset(
		&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "tracing"}},
		jaegerServicePtr("tracing", "jaeger-query", 16686),
	)
	stale := discoveredDatasourceObject("platform-monitoring-jaeger-old-query", jaegerDatasourceComponent)
	reconciler := newDiscoveredReconcilerIntercepting(t, "", kubeClient, interceptor.Funcs{
		Delete: func(ctx context.Context, kubeClient client.WithWatch, obj client.Object, opts ...client.DeleteOption) error {
			if obj.GetName() == stale.GetName() {
				return errors.New("datasource delete failed")
			}
			return kubeClient.Delete(ctx, obj, opts...)
		},
	}, stale)

	err := reconciler.handleJaegerDataSources(cr)

	require.Error(t, err)
	assert.ErrorContains(t, err, "datasource delete failed")
	assert.NoError(t, reconciler.GetResource(discoveredDatasourceObject(stale.GetName(), jaegerDatasourceComponent)))
}

func TestApplyDiscoveredDataSourceSetsLabelsWhenMissing(t *testing.T) {
	cr := discoveredPlatformMonitoring("30s")
	reconciler := newDiscoveredReconciler(t, "", kubernetesfake.NewSimpleClientset())
	desired := grafanaJaegerDataSources(cr, []corev1.Service{jaegerService("tracing", "jaeger-query", 16686)})[0]
	desired.Labels = nil

	require.NoError(t, reconciler.applyDiscoveredDataSource(cr, desired))

	stored := discoveredDatasourceObject(desired.GetName(), jaegerDatasourceComponent)
	require.NoError(t, reconciler.GetResource(stored))
	assert.Equal(t, "platform-monitoring-jaeger-tracing.jaeger-query-monitoring", stored.GetLabels()["app.kubernetes.io/instance"])
	assert.Equal(t, "12.4.3", stored.GetLabels()["app.kubernetes.io/version"])
}

func TestHandleJaegerDataSourcesFlagOffKeepsUnownedDatasource(t *testing.T) {
	withPrivilegedRights(t)
	cr := discoveredIntegrationPlatformMonitoring(false, false)
	unowned := unownedDiscoveredDatasourceObject("user-jaeger", jaegerDatasourceComponent)
	reconciler := newDiscoveredReconciler(t, "", kubernetesfake.NewSimpleClientset(), unowned)

	require.NoError(t, reconciler.handleJaegerDataSources(cr))

	assert.NoError(t, reconciler.GetResource(unowned))
}

func TestHandleJaegerDataSourcesPreservesUnownedDatasourceWithGeneratedName(t *testing.T) {
	withPrivilegedRights(t)
	cr := discoveredIntegrationPlatformMonitoring(true, false)
	service := jaegerService("tracing", "jaeger-query", 16686)
	kubeClient := kubernetesfake.NewSimpleClientset(
		&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: service.Namespace}},
		&service,
	)
	unowned := unownedDiscoveredDatasourceObject(jaegerDatasourceName(service), jaegerDatasourceComponent)
	unowned.Spec.Datasource = &grafv1.GrafanaDatasourceInternal{Name: "user-managed", URL: "http://user-managed"}
	reconciler := newDiscoveredReconciler(t, "", kubeClient, unowned)

	require.NoError(t, reconciler.handleJaegerDataSources(cr))

	stored := unownedDiscoveredDatasourceObject(unowned.GetName(), jaegerDatasourceComponent)
	require.NoError(t, reconciler.GetResource(stored))
	assert.Equal(t, "user-managed", stored.Spec.Datasource.Name)
	assert.Equal(t, "http://user-managed", stored.Spec.Datasource.URL)
	assert.Empty(t, stored.GetLabels()[grafanaCleanupLabelKey])
}

func TestHandleJaegerDataSourcesDiscoveryErrorDoesNotDelete(t *testing.T) {
	withPrivilegedRights(t)
	cr := discoveredIntegrationPlatformMonitoring(true, false)
	kubeClient := kubernetesfake.NewSimpleClientset()
	kubeClient.PrependReactor("list", "namespaces", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, errors.New("namespace list failed")
	})
	owned := discoveredDatasourceObject("platform-monitoring-jaeger-tracing.jaeger-query", jaegerDatasourceComponent)
	reconciler := newDiscoveredReconciler(t, "", kubeClient, owned)

	err := reconciler.handleJaegerDataSources(cr)

	require.Error(t, err)
	assert.ErrorContains(t, err, "namespace list failed")
	assert.NoError(t, reconciler.GetResource(discoveredDatasourceObject("platform-monitoring-jaeger-tracing.jaeger-query", jaegerDatasourceComponent)))
}

func TestHandleClickHouseDataSourcesCredentialsErrorDoesNotDelete(t *testing.T) {
	withPrivilegedRights(t)
	cr := discoveredIntegrationPlatformMonitoring(false, true)
	kubeClient := kubernetesfake.NewSimpleClientset(
		&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "analytics"}},
		clickHouseServicePtr("analytics"),
	)
	kubeClient.PrependReactor("get", "secrets", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, errors.New("secret read failed")
	})
	owned := discoveredDatasourceObject("platform-monitoring-clickhouse-analytics", clickHouseDatasourceComponent)
	reconciler := newDiscoveredReconciler(t, "", kubeClient, owned)

	err := reconciler.handleClickHouseDataSources(cr)

	require.Error(t, err)
	assert.ErrorContains(t, err, "secret read failed")
	assert.NoError(t, reconciler.GetResource(discoveredDatasourceObject("platform-monitoring-clickhouse-analytics", clickHouseDatasourceComponent)))
}

func TestHandleClickHouseDataSourcesDiscoveryErrorDoesNotDelete(t *testing.T) {
	withPrivilegedRights(t)
	cr := discoveredIntegrationPlatformMonitoring(false, true)
	kubeClient := kubernetesfake.NewSimpleClientset(
		&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: "analytics"}},
	)
	kubeClient.PrependReactor("list", "services", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, errors.New("service list failed")
	})
	owned := discoveredDatasourceObject("platform-monitoring-clickhouse-analytics", clickHouseDatasourceComponent)
	reconciler := newDiscoveredReconciler(t, "", kubeClient, owned)

	err := reconciler.handleClickHouseDataSources(cr)

	require.Error(t, err)
	assert.ErrorContains(t, err, "service list failed")
	assert.NoError(t, reconciler.GetResource(discoveredDatasourceObject("platform-monitoring-clickhouse-analytics", clickHouseDatasourceComponent)))
}

func TestUninstallDeletesDiscoveredDataSources(t *testing.T) {
	cr := discoveredPlatformMonitoring("30s")
	cr.Spec.Grafana.Install = boolPtr(false)
	jaegerDatasource := discoveredDatasourceObject("platform-monitoring-jaeger-tracing.jaeger-query", jaegerDatasourceComponent)
	clickHouseDatasource := discoveredDatasourceObject("platform-monitoring-clickhouse-analytics", clickHouseDatasourceComponent)
	reconciler := newDiscoveredReconciler(t, "", kubernetesfake.NewSimpleClientset(), jaegerDatasource, clickHouseDatasource)

	reconciler.uninstall(cr)

	assert.Error(t, reconciler.GetResource(discoveredDatasourceObject(jaegerDatasource.GetName(), jaegerDatasourceComponent)))
	assert.Error(t, reconciler.GetResource(discoveredDatasourceObject(clickHouseDatasource.GetName(), clickHouseDatasourceComponent)))
}

func TestUninstallKeepsDiscoveredDataSourcesWhenListFails(t *testing.T) {
	cr := discoveredPlatformMonitoring("30s")
	cr.Spec.Grafana.Install = boolPtr(false)
	jaegerDatasource := discoveredDatasourceObject("platform-monitoring-jaeger-tracing.jaeger-query", jaegerDatasourceComponent)
	clickHouseDatasource := discoveredDatasourceObject("platform-monitoring-clickhouse-analytics", clickHouseDatasourceComponent)
	reconciler := newDiscoveredReconcilerIntercepting(t, "", kubernetesfake.NewSimpleClientset(), interceptor.Funcs{
		List: func(ctx context.Context, kubeClient client.WithWatch, list client.ObjectList, opts ...client.ListOption) error {
			if _, ok := list.(*grafv1.GrafanaDatasourceList); ok {
				return errors.New("datasource list failed")
			}
			return kubeClient.List(ctx, list, opts...)
		},
	}, jaegerDatasource, clickHouseDatasource)

	reconciler.uninstall(cr)

	assert.NoError(t, reconciler.GetResource(discoveredDatasourceObject(jaegerDatasource.GetName(), jaegerDatasourceComponent)))
	assert.NoError(t, reconciler.GetResource(discoveredDatasourceObject(clickHouseDatasource.GetName(), clickHouseDatasourceComponent)))
}

func discoveredPlatformMonitoring(scrapeInterval string) *monv1.PlatformMonitoring {
	install := true
	return &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Name: "platformmonitoring", Namespace: "monitoring"},
		Spec: monv1.PlatformMonitoringSpec{
			Grafana: &monv1.Grafana{Image: "grafana:12.4.3"},
			Victoriametrics: &monv1.Victoriametrics{
				VmAgent: monv1.VmAgent{Install: &install, Image: "vmagent", ScrapeInterval: scrapeInterval},
			},
		},
	}
}

func discoveredIntegrationPlatformMonitoring(jaeger, clickHouse bool) *monv1.PlatformMonitoring {
	cr := discoveredPlatformMonitoring("30s")
	cr.Spec.Integration = &monv1.Integration{
		Jaeger:     &monv1.Jaeger{CreateGrafanaDataSource: jaeger},
		ClickHouse: &monv1.ClickHouse{CreateGrafanaDataSource: clickHouse},
	}
	return cr
}

func jaegerService(namespace, name string, port int32) corev1.Service {
	return corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    utils.JaegerServiceLabels,
		},
		Spec: corev1.ServiceSpec{Ports: []corev1.ServicePort{{Name: "http-query", Port: port}}},
	}
}

func jaegerServicePtr(namespace, name string, port int32) *corev1.Service {
	service := jaegerService(namespace, name, port)
	return &service
}

func clickHouseService(namespace string) corev1.Service {
	return corev1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: utils.ClickHouseServiceName, Namespace: namespace},
	}
}

func clickHouseServicePtr(namespace string) *corev1.Service {
	service := clickHouseService(namespace)
	return &service
}

func assertDiscoveredDatasource(t *testing.T, datasource *grafv1.GrafanaDatasource, name, component string) {
	t.Helper()
	assert.Equal(t, name, datasource.GetName())
	assert.Equal(t, "monitoring", datasource.GetNamespace())
	assert.Equal(t, component, datasource.GetLabels()["app.kubernetes.io/component"])
	assert.Equal(t, grafanaCleanupLabelValue, datasource.GetLabels()[grafanaCleanupLabelKey])
	assert.Equal(t, grafanaDatasourceResync, datasource.Spec.ResyncPeriod.Duration)
	require.NotNil(t, datasource.Spec.InstanceSelector)
	assert.Equal(t, map[string]string{
		"app.kubernetes.io/component": "grafana",
		"app.kubernetes.io/part-of":   "monitoring",
	}, datasource.Spec.InstanceSelector.MatchLabels)
}

func withPrivilegedRights(t *testing.T) {
	t.Helper()
	previous := utils.PrivilegedRights
	utils.PrivilegedRights = true
	t.Cleanup(func() { utils.PrivilegedRights = previous })
}

func boolPtr(value bool) *bool {
	return &value
}

func legacyCombinedDatasource() client.Object {
	legacy := &unstructured.Unstructured{}
	legacy.SetGroupVersionKind(legacyGrafanaDatasourceGVK)
	legacy.SetName(legacyCombinedDatasourceName)
	legacy.SetNamespace("monitoring")
	return legacy
}

func discoveredDatasourceObject(name, component string) *grafv1.GrafanaDatasource {
	datasource := &grafv1.GrafanaDatasource{}
	datasource.SetName(name)
	datasource.SetNamespace("monitoring")
	datasource.SetLabels(map[string]string{
		"app.kubernetes.io/component": component,
		grafanaCleanupLabelKey:        grafanaCleanupLabelValue,
	})
	return datasource
}

func unownedDiscoveredDatasourceObject(name, component string) *grafv1.GrafanaDatasource {
	datasource := &grafv1.GrafanaDatasource{}
	datasource.SetName(name)
	datasource.SetNamespace("monitoring")
	datasource.SetLabels(map[string]string{"app.kubernetes.io/component": component})
	return datasource
}

func newDiscoveredReconciler(t *testing.T, adminURL string, kubeClient kubernetes.Interface, objects ...client.Object) *GrafanaReconciler {
	t.Helper()
	return newDiscoveredReconcilerIntercepting(t, adminURL, kubeClient, interceptor.Funcs{}, objects...)
}

func newDiscoveredReconcilerIntercepting(t *testing.T, adminURL string, kubeClient kubernetes.Interface, intercept interceptor.Funcs, objects ...client.Object) *GrafanaReconciler {
	t.Helper()
	scheme := runtime.NewScheme()
	require.NoError(t, monv1.AddToScheme(scheme))
	require.NoError(t, grafv1.AddToScheme(scheme))
	require.NoError(t, corev1.AddToScheme(scheme))
	require.NoError(t, promv1.AddToScheme(scheme))
	currentGrafana := &grafv1.Grafana{
		ObjectMeta: metav1.ObjectMeta{Name: "grafana", Namespace: "monitoring"},
		Status:     grafv1.GrafanaStatus{AdminURL: adminURL},
	}
	credentials := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{Name: "grafana-admin-credentials", Namespace: "monitoring"},
		Data: map[string][]byte{
			"GF_SECURITY_ADMIN_USER":     []byte("admin"),
			"GF_SECURITY_ADMIN_PASSWORD": []byte("password"),
		},
	}
	stored := []client.Object{currentGrafana, credentials}
	stored = append(stored, objects...)
	return &GrafanaReconciler{
		KubeClient: kubeClient,
		ComponentReconciler: &utils.ComponentReconciler{
			Client: fake.NewClientBuilder().WithScheme(scheme).WithObjects(stored...).WithInterceptorFuncs(intercept).Build(),
			Scheme: scheme,
			Log:    utils.Logger("grafana_test"),
			Dc:     &discoveryfake.FakeDiscovery{Fake: &k8stesting.Fake{}},
		},
	}
}
