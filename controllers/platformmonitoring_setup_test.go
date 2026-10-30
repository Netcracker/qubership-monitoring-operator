package controllers

import (
	"errors"
	"os"
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	vmetricsv1b1 "github.com/VictoriaMetrics/operator/api/operator/v1beta1"
	grafv1 "github.com/grafana/grafana-operator/v5/api/v1beta1"
	secv1 "github.com/openshift/api/security/v1"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	fakediscovery "k8s.io/client-go/discovery/fake"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ktesting "k8s.io/client-go/testing"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrlconfig "sigs.k8s.io/controller-runtime/pkg/config"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/log"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
)

func startControllerEnvtest(t *testing.T) *rest.Config {
	t.Helper()
	if os.Getenv("KUBEBUILDER_ASSETS") == "" {
		t.Skip("KUBEBUILDER_ASSETS is required to exercise SetupWithManager")
	}
	env := &envtest.Environment{}
	cfg, err := env.Start()
	if err != nil {
		t.Fatalf("start envtest: %v", err)
	}
	t.Cleanup(func() {
		if err := env.Stop(); err != nil {
			t.Errorf("stop envtest: %v", err)
		}
	})
	return cfg
}

func newWatchScheme(t *testing.T) *runtime.Scheme {
	t.Helper()
	s := runtime.NewScheme()
	for _, add := range []func(*runtime.Scheme) error{
		scheme.AddToScheme,
		monv1.AddToScheme,
		grafv1.AddToScheme,
		promv1.AddToScheme,
		vmetricsv1b1.AddToScheme,
		secv1.Install,
	} {
		if err := add(s); err != nil {
			t.Fatal(err)
		}
	}
	return s
}

func newFakeDiscovery(lists ...*metav1.APIResourceList) *fakediscovery.FakeDiscovery {
	dc := &fakediscovery.FakeDiscovery{Fake: &ktesting.Fake{}}
	dc.Resources = lists
	return dc
}

func apiResourceList(groupVersion string, kinds ...string) *metav1.APIResourceList {
	apis := make([]metav1.APIResource, 0, len(kinds))
	for _, kind := range kinds {
		apis = append(apis, metav1.APIResource{Kind: kind, Name: kind, Verbs: []string{"list", "watch", "get"}})
	}
	return &metav1.APIResourceList{GroupVersion: groupVersion, APIResources: apis}
}

func setupWithManager(t *testing.T, cfg *rest.Config, dc *fakediscovery.FakeDiscovery, privileged bool) {
	t.Helper()
	prev := utils.PrivilegedRights
	utils.PrivilegedRights = privileged
	t.Cleanup(func() { utils.PrivilegedRights = prev })

	sch := newWatchScheme(t)
	mgr, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme:     sch,
		Metrics:    metricsserver.Options{BindAddress: "0"},
		Controller: ctrlconfig.Controller{SkipNameValidation: ptr.To(true)},
	})
	if err != nil {
		t.Fatalf("NewManager: %v", err)
	}
	r := &PlatformMonitoringReconciler{
		Client:          mgr.GetClient(),
		Scheme:          sch,
		Log:             log.Log,
		Config:          cfg,
		DiscoveryClient: dc,
	}
	if err := r.SetupWithManager(mgr); err != nil {
		t.Fatalf("SetupWithManager: %v", err)
	}
}

func TestAddWatchBasedSourcesDiscoversOnce(t *testing.T) {
	dc := newFakeDiscovery(
		apiResourceList(grafv1.SchemeGroupVersion.String(), "Grafana", "GrafanaDashboard", "GrafanaDatasource"),
		apiResourceList(promv1.SchemeGroupVersion.String(), "Prometheus", "Alertmanager", "PrometheusRule", "ServiceMonitor", "PodMonitor"),
		apiResourceList(vmetricsv1b1.SchemeGroupVersion.String(), "VMSingle", "VMAgent", "VMAuth", "VMAlert", "VMAlertmanager", "VMCluster", "VMUser"),
		apiResourceList(networkingv1.SchemeGroupVersion.String(), "Ingress"),
		apiResourceList(gatewayAPIGroup+"/"+httpRouteVersion, httpRouteKind),
		apiResourceList(secv1.GroupVersion.String(), "SecurityContextConstraints"),
	)
	r := &PlatformMonitoringReconciler{Scheme: runtime.NewScheme(), Log: log.Log, DiscoveryClient: dc}
	if _, err := r.addWatchBasedSources(ctrl.NewControllerManagedBy(nil).For(&monv1.PlatformMonitoring{})); err != nil {
		t.Fatalf("addWatchBasedSources: %v", err)
	}

	var calls int
	for _, action := range dc.Actions() {
		if action.GetResource().Resource == "resource" {
			calls++
		}
	}
	if calls != 1 {
		t.Fatalf("ServerGroupsAndResources calls = %d, want 1", calls)
	}
}

func TestOwnHTTPRouteLeavesSchemeUnchanged(t *testing.T) {
	sch := runtime.NewScheme()
	gvk := schema.GroupVersion{Group: gatewayAPIGroup, Version: httpRouteVersion}.WithKind(httpRouteKind)
	dc := newFakeDiscovery(apiResourceList(gvk.GroupVersion().String(), httpRouteKind))
	r := &PlatformMonitoringReconciler{Scheme: sch, Log: log.Log, DiscoveryClient: dc}
	r.ownHTTPRouteIfServed(ctrl.NewControllerManagedBy(nil).For(&monv1.PlatformMonitoring{}))
	if sch.Recognizes(gvk) {
		t.Fatalf("scheme recognizes %s after HTTPRoute watch setup", gvk)
	}
}

func TestSetupWithManagerWatchBasedSources(t *testing.T) {
	cfg := startControllerEnvtest(t)

	t.Run("privileged mixed discovery registers Owns and cluster watches", func(t *testing.T) {
		dc := newFakeDiscovery(
			apiResourceList(grafv1.SchemeGroupVersion.String(), "Grafana", "GrafanaDashboard", "GrafanaDatasource"),
			apiResourceList(promv1.SchemeGroupVersion.String(), "Prometheus", "Alertmanager", "PrometheusRule", "ServiceMonitor", "PodMonitor"),
			apiResourceList(vmetricsv1b1.SchemeGroupVersion.String(), "VMSingle", "VMAgent", "VMAuth", "VMAlert", "VMAlertmanager", "VMCluster", "VMUser"),
			apiResourceList(networkingv1.SchemeGroupVersion.String(), "Ingress"),
			apiResourceList(gatewayAPIGroup+"/"+httpRouteVersion, httpRouteKind),
			apiResourceList(secv1.GroupVersion.String(), "SecurityContextConstraints"),
		)
		setupWithManager(t, cfg, dc, true)
	})

	t.Run("namespaced install skips cluster and Jaeger watches", func(t *testing.T) {
		dc := newFakeDiscovery(
			apiResourceList(grafv1.SchemeGroupVersion.String(), "Grafana", "GrafanaDashboard"),
		)
		setupWithManager(t, cfg, dc, false)
	})

	t.Run("watch-based reconcile disabled", func(t *testing.T) {
		t.Setenv(watchBasedReconcileEnv, "false")
		setupWithManager(t, cfg, newFakeDiscovery(), true)
	})
}

func withFastDiscoveryRetry(t *testing.T) {
	t.Helper()
	prevDelay := servedKindsRetryDelay
	servedKindsRetryDelay = 0
	t.Cleanup(func() { servedKindsRetryDelay = prevDelay })
}

// failDiscovery makes the first failures resource-discovery calls return err.
// A negative count fails every call.
func failDiscovery(dc *fakediscovery.FakeDiscovery, failures int, err error) {
	dc.AddReactor("get", "resource", func(ktesting.Action) (bool, runtime.Object, error) {
		if failures == 0 {
			return false, nil, nil
		}
		failures--
		return true, nil, err
	})
}

func groupDiscoveryFailed(gv schema.GroupVersion) error {
	return &discovery.ErrGroupDiscoveryFailed{Groups: map[schema.GroupVersion]error{gv: errors.New("service unavailable")}}
}

func TestLoadServedKindsKeepsHealthyListsOnUnrelatedGroupFailure(t *testing.T) {
	withFastDiscoveryRetry(t)
	grafGV := grafv1.SchemeGroupVersion.String()
	dc := newFakeDiscovery(apiResourceList(grafGV, "Grafana"))
	failDiscovery(dc, -1, groupDiscoveryFailed(schema.GroupVersion{Group: "metrics.k8s.io", Version: "v1beta1"}))

	kinds, err := loadServedKinds(dc, watchedGroupVersions())
	if err != nil {
		t.Fatalf("loadServedKinds: %v", err)
	}
	if _, ok := kinds[grafGV]["Grafana"]; !ok {
		t.Fatalf("Grafana missing from served kinds after an unrelated group failed: %v", kinds)
	}
}

func TestLoadServedKindsFailsWhenWatchedGroupStaysUnavailable(t *testing.T) {
	withFastDiscoveryRetry(t)
	dc := newFakeDiscovery()
	failDiscovery(dc, -1, groupDiscoveryFailed(grafv1.SchemeGroupVersion))

	if _, err := loadServedKinds(dc, watchedGroupVersions()); err == nil {
		t.Fatal("loadServedKinds succeeded while a watched group version failed discovery")
	}
	if got := len(dc.Actions()); got != servedKindsAttempts*2 {
		t.Fatalf("discovery actions = %d, want %d (groups and resources per attempt)", got, servedKindsAttempts*2)
	}
}

func TestLoadServedKindsRetriesTransientFailure(t *testing.T) {
	withFastDiscoveryRetry(t)
	promGV := promv1.SchemeGroupVersion.String()
	dc := newFakeDiscovery(apiResourceList(promGV, "Prometheus"))
	failDiscovery(dc, 1, errors.New("connection refused"))

	kinds, err := loadServedKinds(dc, watchedGroupVersions())
	if err != nil {
		t.Fatalf("loadServedKinds: %v", err)
	}
	if _, ok := kinds[promGV]["Prometheus"]; !ok {
		t.Fatalf("Prometheus missing from served kinds after a retry: %v", kinds)
	}
}

func TestAddWatchBasedSourcesReturnsDiscoveryError(t *testing.T) {
	withFastDiscoveryRetry(t)
	dc := newFakeDiscovery()
	failDiscovery(dc, -1, errors.New("connection refused"))
	r := &PlatformMonitoringReconciler{Scheme: runtime.NewScheme(), Log: log.Log, DiscoveryClient: dc}

	if _, err := r.addWatchBasedSources(ctrl.NewControllerManagedBy(nil).For(&monv1.PlatformMonitoring{})); err == nil {
		t.Fatal("addWatchBasedSources succeeded although discovery failed")
	}
}

func TestIndexServedKindsNilClient(t *testing.T) {
	kinds, failed, err := indexServedKinds(nil)
	if err != nil || len(kinds) != 0 || len(failed) != 0 {
		t.Fatalf("indexServedKinds(nil) = %v, %v, %v; want empty results", kinds, failed, err)
	}
}
