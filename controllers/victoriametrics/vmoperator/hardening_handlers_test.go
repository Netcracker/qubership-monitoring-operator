package vmoperator

import (
	"context"
	"errors"
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
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
	require.NoError(t, appsv1.AddToScheme(scheme))
	c := fake.NewClientBuilder().WithScheme(scheme).Build()
	dc := failingDiscovery()
	reconciler := &VmOperatorReconciler{ComponentReconciler: &utils.ComponentReconciler{
		Client: c, Scheme: scheme, Dc: dc, Log: utils.Logger("vmoperator_test"),
	}}
	cr := &monv1.PlatformMonitoring{ObjectMeta: metav1.ObjectMeta{Name: "platformmonitoring", Namespace: "monitoring"}}

	err := reconciler.handleDeployment(cr)

	require.Error(t, err)
	err = c.Get(context.Background(), types.NamespacedName{Name: utils.VmOperatorComponentName, Namespace: "monitoring"}, &appsv1.Deployment{})
	assert.True(t, apierrors.IsNotFound(err), "no workload must be created while the platform is unknown")
}

// Uninstall must succeed even when platform discovery fails, because it only needs the stable
// name and namespace of the resource.
func TestDeleteDoesNotDependOnDiscovery(t *testing.T) {
	scheme := runtime.NewScheme()
	require.NoError(t, appsv1.AddToScheme(scheme))
	existing := &appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Name: utils.VmOperatorComponentName, Namespace: "monitoring"}}
	c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(existing).Build()
	dc := failingDiscovery()
	reconciler := &VmOperatorReconciler{ComponentReconciler: &utils.ComponentReconciler{
		Client: c, Scheme: scheme, Dc: dc, Log: utils.Logger("vmoperator_test"),
	}}
	cr := &monv1.PlatformMonitoring{ObjectMeta: metav1.ObjectMeta{Name: "platformmonitoring", Namespace: "monitoring"}}

	require.NoError(t, reconciler.deleteVmOperatorDeployment(cr))

	err := c.Get(context.Background(), types.NamespacedName{Name: utils.VmOperatorComponentName, Namespace: "monitoring"}, &appsv1.Deployment{})
	assert.True(t, apierrors.IsNotFound(err), "the resource must be deleted")
	require.NoError(t, reconciler.deleteVmOperatorDeployment(cr), "a missing resource is not an error")
}

func TestManifestRejectsConflictingHardeningSettings(t *testing.T) {
	tests := map[string]monv1.VmOperator{
		"root user":            {SecurityContext: &monv1.SecurityContext{RunAsUser: ptr.To[int64](0)}},
		"privileged container": {ContainerSecurityContext: &corev1.SecurityContext{Privileged: ptr.To(true)}},
	}
	for name, spec := range tests {
		t.Run(name, func(t *testing.T) {
			cr := &monv1.PlatformMonitoring{
				ObjectMeta: metav1.ObjectMeta{Name: "platformmonitoring", Namespace: "monitoring"},
				Spec:       monv1.PlatformMonitoringSpec{Victoriametrics: &monv1.Victoriametrics{VmOperator: spec}},
			}

			m, err := vmOperatorDeployment(nil, cr, false)

			require.Error(t, err)
			assert.Nil(t, m)
		})
	}
}
