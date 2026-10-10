package prometheus

import (
	"context"
	"errors"
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	fakediscovery "k8s.io/client-go/discovery/fake"
	ktesting "k8s.io/client-go/testing"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

// Uninstall must succeed even when platform discovery fails and the desired state is invalid,
// because it only needs the stable name and namespace of the resource.
func TestDeletePrometheusDoesNotDependOnDiscoveryOrDesiredState(t *testing.T) {
	scheme := runtime.NewScheme()
	require.NoError(t, promv1.AddToScheme(scheme))
	existing := &promv1.Prometheus{ObjectMeta: metav1.ObjectMeta{Name: "k8s", Namespace: "monitoring"}}
	controllerClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(existing).Build()
	discoveryClient := &fakediscovery.FakeDiscovery{Fake: &ktesting.Fake{}}
	discoveryClient.PrependReactor("get", "resource", func(ktesting.Action) (bool, runtime.Object, error) {
		return true, nil, errors.New("aggregated API unavailable")
	})
	reconciler := NewPrometheusReconciler(controllerClient, scheme, discoveryClient)
	cr := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Name: "platformmonitoring", Namespace: "monitoring"},
		Spec: monv1.PlatformMonitoringSpec{
			Prometheus: &monv1.Prometheus{
				Install: ptr.To(false),
				Containers: []corev1.Container{{
					Name:            "sidecar",
					SecurityContext: &corev1.SecurityContext{Privileged: ptr.To(true)},
				}},
			},
		},
	}

	require.NoError(t, reconciler.deletePrometheus(cr))

	err := controllerClient.Get(context.Background(), types.NamespacedName{Name: "k8s", Namespace: "monitoring"}, &promv1.Prometheus{})
	assert.True(t, apierrors.IsNotFound(err), "the Prometheus resource must be deleted")
}

// A workload must not be created or updated while the platform is unknown.
func TestHandlePrometheusReturnsDiscoveryError(t *testing.T) {
	scheme := runtime.NewScheme()
	require.NoError(t, promv1.AddToScheme(scheme))
	controllerClient := fake.NewClientBuilder().WithScheme(scheme).Build()
	discoveryClient := &fakediscovery.FakeDiscovery{Fake: &ktesting.Fake{}}
	discoveryClient.PrependReactor("get", "resource", func(ktesting.Action) (bool, runtime.Object, error) {
		return true, nil, errors.New("aggregated API unavailable")
	})
	reconciler := NewPrometheusReconciler(controllerClient, scheme, discoveryClient)
	cr := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Name: "platformmonitoring", Namespace: "monitoring"},
		Spec:       monv1.PlatformMonitoringSpec{Prometheus: &monv1.Prometheus{}},
	}

	require.Error(t, reconciler.handlePrometheus(cr))

	err := controllerClient.Get(context.Background(), types.NamespacedName{Name: "k8s", Namespace: "monitoring"}, &promv1.Prometheus{})
	assert.True(t, apierrors.IsNotFound(err), "no Prometheus must be created while the platform is unknown")
}
