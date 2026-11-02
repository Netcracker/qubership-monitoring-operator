// Package reconciletest checks that component handlers create their objects
// once and skip Updates on idle reconciles.
package reconciletest

import (
	"context"
	"sync/atomic"
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
)

// Handler reconciles one object for a PlatformMonitoring resource.
type Handler func(cr *monv1.PlatformMonitoring) error

// UpdateCounter counts the Update calls sent through a client built by NewClient.
type UpdateCounter struct {
	updates atomic.Int64
}

// Updates returns the number of Update calls since the last Reset.
func (c *UpdateCounter) Updates() int64 {
	return c.updates.Load()
}

// Reset sets the count to zero.
func (c *UpdateCounter) Reset() {
	c.updates.Store(0)
}

// NewScheme returns a scheme with the PlatformMonitoring and built-in Kubernetes
// types plus the types that addToScheme registers.
func NewScheme(t *testing.T, addToScheme ...func(*runtime.Scheme) error) *runtime.Scheme {
	t.Helper()
	scheme := runtime.NewScheme()
	require.NoError(t, monv1.AddToScheme(scheme))
	require.NoError(t, clientgoscheme.AddToScheme(scheme))
	for _, add := range addToScheme {
		require.NoError(t, add(scheme))
	}
	return scheme
}

// NewClient returns a fake client that holds objs and records every Update in
// counter.
func NewClient(scheme *runtime.Scheme, counter *UpdateCounter, objs ...client.Object) client.WithWatch {
	return fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(objs...).
		WithInterceptorFuncs(interceptor.Funcs{
			Update: func(ctx context.Context, c client.WithWatch, obj client.Object, opts ...client.UpdateOption) error {
				counter.updates.Add(1)
				return c.Update(ctx, obj, opts...)
			},
		}).
		Build()
}

// AssertIdleReconcileSkipsUpdates runs each handler twice against the client
// that counter observes. The first run creates the objects; the second run
// must succeed without sending an Update.
func AssertIdleReconcileSkipsUpdates(t *testing.T, cr *monv1.PlatformMonitoring, counter *UpdateCounter, handlers map[string]Handler) {
	t.Helper()
	for name, handle := range handlers {
		t.Run(name, func(t *testing.T) {
			require.NoError(t, handle(cr), "first reconcile")
			counter.Reset()
			require.NoError(t, handle(cr), "idle reconcile")
			require.Zero(t, counter.Updates(), "idle reconcile must not send an Update")
		})
	}
}
