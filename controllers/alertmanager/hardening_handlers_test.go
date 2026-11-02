package alertmanager

import (
	"context"
	"errors"
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	fakediscovery "k8s.io/client-go/discovery/fake"
	ktesting "k8s.io/client-go/testing"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func failingDiscovery() *fakediscovery.FakeDiscovery {
	discoveryClient := &fakediscovery.FakeDiscovery{Fake: &ktesting.Fake{}}
	discoveryClient.PrependReactor("get", "resource", func(ktesting.Action) (bool, runtime.Object, error) {
		return true, nil, errors.New("aggregated API unavailable")
	})
	return discoveryClient
}

// A workload must not be created or updated while the platform is unknown.
func TestHandlerReturnsDiscoveryError(t *testing.T) {
	scheme := runtime.NewScheme()
	require.NoError(t, promv1.AddToScheme(scheme))
	c := fake.NewClientBuilder().WithScheme(scheme).Build()
	dc := failingDiscovery()
	reconciler := NewAlertManagerReconciler(c, scheme, dc)
	cr := &monv1.PlatformMonitoring{ObjectMeta: metav1.ObjectMeta{Name: "platformmonitoring", Namespace: "monitoring"}}

	err := reconciler.handleAlertmanager(cr)

	require.Error(t, err)
	err = c.Get(context.Background(), types.NamespacedName{Name: utils.ManagedCustomResourceName, Namespace: "monitoring"}, &promv1.Alertmanager{})
	assert.True(t, apierrors.IsNotFound(err), "no workload must be created while the platform is unknown")
}

// Uninstall must succeed even when platform discovery fails, because it only needs the stable
// name and namespace of the resource.
func TestDeleteDoesNotDependOnDiscovery(t *testing.T) {
	scheme := runtime.NewScheme()
	require.NoError(t, promv1.AddToScheme(scheme))
	existing := &promv1.Alertmanager{ObjectMeta: metav1.ObjectMeta{Name: utils.ManagedCustomResourceName, Namespace: "monitoring"}}
	c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(existing).Build()
	dc := failingDiscovery()
	reconciler := NewAlertManagerReconciler(c, scheme, dc)
	cr := &monv1.PlatformMonitoring{ObjectMeta: metav1.ObjectMeta{Name: "platformmonitoring", Namespace: "monitoring"}}

	require.NoError(t, reconciler.deleteAlertmanager(cr))

	err := c.Get(context.Background(), types.NamespacedName{Name: utils.ManagedCustomResourceName, Namespace: "monitoring"}, &promv1.Alertmanager{})
	assert.True(t, apierrors.IsNotFound(err), "the resource must be deleted")
	require.NoError(t, reconciler.deleteAlertmanager(cr), "a missing resource is not an error")
}
