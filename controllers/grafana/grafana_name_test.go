package grafana

import (
	"context"
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	grafv1 "github.com/grafana/grafana-operator/v5/api/v1beta1"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
)

func TestGrafanaServiceName(t *testing.T) {
	assert.Equal(t, utils.GrafanaServiceName, grafanaServiceName(adminCredentialsCR(nil)),
		"the default name must keep the existing Service name")
	assert.Equal(t, utils.GrafanaServiceName, grafanaServiceName(&monv1.PlatformMonitoring{}))
	assert.Equal(t, "custom-service", grafanaServiceName(adminCredentialsCR(&monv1.Grafana{Name: "custom"})))
	assert.Equal(t, utils.GrafanaDeploymentName, grafanaDeploymentName(utils.GrafanaComponentName))
}

func grafanaIngressBackends(t *testing.T, grafana *monv1.Grafana) []string {
	t.Helper()
	m, err := grafanaIngressV1(adminCredentialsCR(grafana))
	require.NoError(t, err)
	var names []string
	for _, rule := range m.Spec.Rules {
		for _, path := range rule.HTTP.Paths {
			names = append(names, path.Backend.Service.Name)
		}
	}
	return names
}

func TestGrafanaIngressTargetsGrafanaService(t *testing.T) {
	install := true
	ingress := func() *monv1.Ingress { return &monv1.Ingress{Install: &install, Host: "grafana.example.com"} }

	assert.Equal(t, []string{utils.GrafanaServiceName}, grafanaIngressBackends(t, &monv1.Grafana{Ingress: ingress()}))
	assert.Equal(t, []string{"custom-service"}, grafanaIngressBackends(t, &monv1.Grafana{Name: "custom", Ingress: ingress()}))

	withRules := ingress()
	withRules.Rules = []monv1.IngressRule{{Host: "grafana.example.com"}}
	assert.Equal(t, []string{"custom-service"}, grafanaIngressBackends(t, &monv1.Grafana{Name: "custom", Ingress: withRules}))
}

func TestGrafanaPodMonitorSelectsGrafanaPods(t *testing.T) {
	m, err := grafanaPodMonitor(adminCredentialsCR(nil))
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"app": utils.GrafanaComponentName}, m.Spec.Selector.MatchLabels,
		"the default name must keep the existing selector")

	m, err = grafanaPodMonitor(adminCredentialsCR(&monv1.Grafana{Name: "custom"}))
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"app": "custom"}, m.Spec.Selector.MatchLabels)
}

func TestRunDefersDatasourcesWhilePasswordSyncIsPending(t *testing.T) {
	install := true
	cr := adminCredentialsCR(&monv1.Grafana{Install: &install})
	r := newAdminCredentialsReconciler(t, nil, adminSecret(testNamespace, "grafana", "secret"))

	err := r.Run(cr)
	assert.ErrorIs(t, err, ErrAdminPasswordSyncPending)

	created := getGrafana(t, r, "grafana")
	assert.NotContains(t, created.GetAnnotations(), adminPasswordSyncedAnnotation)
	datasource := &grafv1.GrafanaDatasource{}
	err = r.Client.Get(context.Background(), types.NamespacedName{Name: "platform-monitoring-prometheus", Namespace: testNamespace}, datasource)
	assert.True(t, apierrors.IsNotFound(err), "datasources must wait for the admin password sync, got %v", err)
}

func TestRunReconcilesDatasourcesWithoutAdminSecret(t *testing.T) {
	install := true
	cr := adminCredentialsCR(&monv1.Grafana{Install: &install})
	r := newAdminCredentialsReconciler(t, nil)

	require.NoError(t, r.Run(cr))

	datasource := &grafv1.GrafanaDatasource{}
	require.NoError(t, r.Client.Get(context.Background(), types.NamespacedName{Name: "platform-monitoring-prometheus", Namespace: testNamespace}, datasource))
}

func TestRunReturnsPasswordSyncError(t *testing.T) {
	install := true
	cr := adminCredentialsCR(&monv1.Grafana{Install: &install})
	kubeClient := newExecClientset(t, grafanaDeployment("grafana"), readyPod("grafana-ready", true))
	r := newAdminCredentialsReconciler(t, kubeClient, adminSecret(testNamespace, "grafana", "secret"))
	useFakeExecutor(r, &fakeExecutor{err: assert.AnError})

	err := r.Run(cr)
	assert.ErrorIs(t, err, assert.AnError)
	assert.False(t, IsPending(err))
}

func TestRunSkipsPausedGrafana(t *testing.T) {
	install := true
	r := newAdminCredentialsReconciler(t, nil)

	require.NoError(t, r.Run(adminCredentialsCR(&monv1.Grafana{Install: &install, Paused: true})))

	err := r.Client.Get(context.Background(), types.NamespacedName{Name: "grafana", Namespace: testNamespace}, &grafv1.Grafana{})
	assert.True(t, apierrors.IsNotFound(err))
}

// newRunReconciler builds a reconciler whose client fails Create for the matching objects.
func newRunReconciler(t *testing.T, failCreate func(client.Object) bool) *GrafanaReconciler {
	t.Helper()
	testScheme := runtime.NewScheme()
	require.NoError(t, corev1.AddToScheme(testScheme))
	require.NoError(t, monv1.AddToScheme(testScheme))
	require.NoError(t, grafv1.AddToScheme(testScheme))
	require.NoError(t, promv1.AddToScheme(testScheme))
	r := newAdminCredentialsReconciler(t, nil)
	r.Client = fake.NewClientBuilder().WithScheme(testScheme).WithInterceptorFuncs(interceptor.Funcs{
		Create: func(ctx context.Context, c client.WithWatch, obj client.Object, opts ...client.CreateOption) error {
			if failCreate(obj) {
				return assert.AnError
			}
			return c.Create(ctx, obj, opts...)
		},
	}).Build()
	r.Scheme = testScheme
	return r
}

func TestRunReturnsReconcileErrors(t *testing.T) {
	install := true
	tests := map[string]struct {
		grafana    *monv1.Grafana
		promxy     bool
		failCreate func(client.Object) bool
	}{
		"Grafana": {
			grafana: &monv1.Grafana{Install: &install},
			failCreate: func(obj client.Object) bool {
				_, ok := obj.(*grafv1.Grafana)
				return ok
			},
		},
		"Prometheus datasource": {
			grafana: &monv1.Grafana{Install: &install},
			failCreate: func(obj client.Object) bool {
				_, ok := obj.(*grafv1.GrafanaDatasource)
				return ok
			},
		},
		"Promxy datasource": {
			grafana: &monv1.Grafana{Install: &install},
			promxy:  true,
			failCreate: func(obj client.Object) bool {
				return obj.GetName() == "platform-monitoring-promxy"
			},
		},
		"PodMonitor": {
			grafana: &monv1.Grafana{Install: &install, PodMonitor: &monv1.Monitor{Install: &install}},
			failCreate: func(obj client.Object) bool {
				_, ok := obj.(*promv1.PodMonitor)
				return ok
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			cr := adminCredentialsCR(tc.grafana)
			if tc.promxy {
				cr.Spec.Promxy = &monv1.Promxy{Install: &install}
			}
			r := newRunReconciler(t, tc.failCreate)

			assert.ErrorIs(t, r.Run(cr), assert.AnError)
		})
	}
}

func TestRunUninstallsDisabledGrafana(t *testing.T) {
	install := false
	cr := adminCredentialsCR(&monv1.Grafana{Install: &install})
	r := newAdminCredentialsReconciler(t, nil, existingGrafana(t, cr, ""))

	require.NoError(t, r.Run(cr))

	err := r.Client.Get(context.Background(), types.NamespacedName{Name: "grafana", Namespace: testNamespace}, &grafv1.Grafana{})
	assert.True(t, apierrors.IsNotFound(err))
}

func TestHandleGrafanaReturnsAdminSecretReadError(t *testing.T) {
	r := newAdminCredentialsReconciler(t, nil)
	testScheme := runtime.NewScheme()
	require.NoError(t, grafv1.AddToScheme(testScheme))
	require.NoError(t, monv1.AddToScheme(testScheme))
	r.Client = fake.NewClientBuilder().WithScheme(testScheme).Build()

	assert.ErrorContains(t, r.handleGrafana(adminCredentialsCR(nil)), "cannot read Grafana admin credentials Secret")
}
