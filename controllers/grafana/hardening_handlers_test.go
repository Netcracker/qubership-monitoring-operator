package grafana

import (
	"context"
	"errors"
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	"github.com/go-logr/logr"
	grafv1 "github.com/grafana/grafana-operator/v5/api/v1beta1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	fakediscovery "k8s.io/client-go/discovery/fake"
	k8stesting "k8s.io/client-go/testing"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func failingDiscoveryReconciler(t *testing.T, objects ...client.Object) (*GrafanaReconciler, client.Client) {
	t.Helper()
	scheme := runtime.NewScheme()
	require.NoError(t, monv1.AddToScheme(scheme))
	require.NoError(t, grafv1.AddToScheme(scheme))
	controllerClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(objects...).Build()
	discoveryClient := &fakediscovery.FakeDiscovery{Fake: &k8stesting.Fake{}}
	discoveryClient.PrependReactor("get", "resource", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, errors.New("aggregated API unavailable")
	})
	return &GrafanaReconciler{
		ComponentReconciler: &utils.ComponentReconciler{
			Client: controllerClient,
			Scheme: scheme,
			Dc:     discoveryClient,
			Log:    logr.Discard(),
		},
	}, controllerClient
}

func TestHandleGrafanaReturnsDiscoveryError(t *testing.T) {
	reconciler, controllerClient := failingDiscoveryReconciler(t)
	cr := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Name: "monitoring", Namespace: "monitoring"},
		Spec:       monv1.PlatformMonitoringSpec{Grafana: &monv1.Grafana{}},
	}

	require.Error(t, reconciler.handleGrafana(cr))

	err := controllerClient.Get(context.Background(), types.NamespacedName{Name: "grafana", Namespace: "monitoring"}, &grafv1.Grafana{})
	assert.True(t, apierrors.IsNotFound(err), "no Grafana must be created while the platform is unknown")
}

func TestDeleteGrafanaUsesConfiguredIdentity(t *testing.T) {
	existing := &grafv1.Grafana{ObjectMeta: metav1.ObjectMeta{Name: "custom-grafana", Namespace: "observability"}}
	reconciler, controllerClient := failingDiscoveryReconciler(t, existing)
	cr := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Name: "monitoring", Namespace: "monitoring"},
		Spec: monv1.PlatformMonitoringSpec{Grafana: &monv1.Grafana{
			Name:      "custom-grafana",
			Namespace: "observability",
			// A setting that the hardened manifest rejects must not block uninstall.
			SecurityContext: &monv1.SecurityContext{RunAsUser: ptr(int64(0))},
		}},
	}

	require.NoError(t, reconciler.deleteGrafana(cr))

	err := controllerClient.Get(context.Background(), types.NamespacedName{Name: "custom-grafana", Namespace: "observability"}, &grafv1.Grafana{})
	assert.True(t, apierrors.IsNotFound(err), "the Grafana resource must be deleted")
	require.NoError(t, reconciler.deleteGrafana(cr), "a missing resource is not an error")
}

func TestGrafanaManifestRejectsConflictingHardeningSettings(t *testing.T) {
	tests := map[string]*monv1.Grafana{
		"root user": {SecurityContext: &monv1.SecurityContext{RunAsUser: ptr(int64(0))}},
	}
	for name, spec := range tests {
		t.Run(name, func(t *testing.T) {
			cr := &monv1.PlatformMonitoring{
				ObjectMeta: metav1.ObjectMeta{Name: "monitoring", Namespace: "monitoring"},
				Spec:       monv1.PlatformMonitoringSpec{Grafana: spec},
			}

			m, err := grafana(cr, false)

			require.Error(t, err)
			assert.Nil(t, m)
		})
	}
}
