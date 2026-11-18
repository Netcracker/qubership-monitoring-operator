package grafana_operator

import (
	"context"
	"errors"
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	grafv1 "github.com/grafana/grafana-operator/v5/api/v1beta1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
)

var (
	errInjected     = errors.New("injected API error")
	errLegacyNoKind = &meta.NoKindMatchError{GroupKind: legacyGrafanaDashboardGVK.GroupKind()}
	errNativeNoKind = &meta.NoKindMatchError{GroupKind: schema.GroupKind{Group: "grafana.integreatly.org", Kind: "GrafanaDashboard"}}
)

func isLegacyDashboard(obj client.Object) bool {
	u, ok := obj.(*unstructured.Unstructured)
	return ok && u.GroupVersionKind() == legacyGrafanaDashboardGVK
}

func isNativeDashboard(obj client.Object) bool {
	_, ok := obj.(*grafv1.GrafanaDashboard)
	return ok
}

// newInterceptedDashboardReconciler builds a reconciler whose fake client returns the errors that funcs inject.
func newInterceptedDashboardReconciler(t *testing.T, funcs interceptor.Funcs, objs ...client.Object) *GrafanaOperatorReconciler {
	t.Helper()
	scheme := runtime.NewScheme()
	require.NoError(t, monv1.AddToScheme(scheme))
	require.NoError(t, grafv1.AddToScheme(scheme))
	return &GrafanaOperatorReconciler{
		ComponentReconciler: &utils.ComponentReconciler{
			Client: fake.NewClientBuilder().WithScheme(scheme).WithObjects(objs...).WithInterceptorFuncs(funcs).Build(),
			Scheme: scheme,
			Log:    utils.Logger("grafanaoperator_dashboard_errors_test"),
		},
	}
}

// failGet returns a Get interceptor that fails for objects that match.
func failGet(match func(client.Object) bool, err error) func(context.Context, client.WithWatch, client.ObjectKey, client.Object, ...client.GetOption) error {
	return func(ctx context.Context, c client.WithWatch, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
		if match(obj) {
			return err
		}
		return c.Get(ctx, key, obj, opts...)
	}
}

func builtInCopy(cr *monv1.PlatformMonitoring, sourceName string) *grafv1.GrafanaDashboard {
	return &grafv1.GrafanaDashboard{ObjectMeta: metav1.ObjectMeta{
		Name:      sourceName + convertedDashboardNameSuffix,
		Namespace: cr.GetNamespace(),
		Labels: map[string]string{
			converterManagedByLabelKey: converterManagedByLabelValue,
			builtInDashboardLabelKey:   builtInDashboardLabelValue,
		},
	}}
}

func TestDeleteOrphanedConvertedDashboardsHandlesAPIErrors(t *testing.T) {
	failList := func(err error) func(context.Context, client.WithWatch, client.ObjectList, ...client.ListOption) error {
		return func(context.Context, client.WithWatch, client.ObjectList, ...client.ListOption) error { return err }
	}
	for name, tc := range map[string]struct {
		funcs       interceptor.Funcs
		wantErr     bool
		wantDeleted bool
	}{
		"v1beta1 API not served":   {funcs: interceptor.Funcs{List: failList(errNativeNoKind)}},
		"listing copies fails":     {funcs: interceptor.Funcs{List: failList(errInjected)}, wantErr: true},
		"legacy API not served":    {funcs: interceptor.Funcs{Get: failGet(isLegacyDashboard, errLegacyNoKind)}},
		"reading the source fails": {funcs: interceptor.Funcs{Get: failGet(isLegacyDashboard, errInjected)}, wantErr: true},
		"deleting the copy fails": {
			funcs: interceptor.Funcs{Delete: func(context.Context, client.WithWatch, client.Object, ...client.DeleteOption) error {
				return errInjected
			}},
			wantErr: true,
		},
		"copy already deleted": {
			funcs: interceptor.Funcs{Delete: func(_ context.Context, _ client.WithWatch, obj client.Object, _ ...client.DeleteOption) error {
				return apierrors.NewNotFound(grafv1.SchemeGroupVersion.WithResource("grafanadashboards").GroupResource(), obj.GetName())
			}},
		},
		"no errors": {wantDeleted: true},
	} {
		t.Run(name, func(t *testing.T) {
			cr := testPlatformMonitoring()
			orphan := builtInCopy(cr, "version-exporter")
			r := newInterceptedDashboardReconciler(t, tc.funcs, cr, orphan)

			err := r.deleteOrphanedConvertedDashboards(cr)

			if tc.wantErr {
				require.ErrorIs(t, err, errInjected)
			} else {
				require.NoError(t, err)
			}
			// The Get interceptors fail legacy reads only, so this read of the copy is real.
			err = r.Client.Get(t.Context(), client.ObjectKeyFromObject(orphan), &grafv1.GrafanaDashboard{})
			assert.Equal(t, tc.wantDeleted, apierrors.IsNotFound(err), "Get copy error = %v", err)
		})
	}
}

func TestHandleGrafanaDashboardReturnsAPIErrors(t *testing.T) {
	for name, funcs := range map[string]interceptor.Funcs{
		"reading the v1beta1 dashboard fails": {Get: failGet(isNativeDashboard, errInjected)},
		"reading the legacy dashboard fails":  {Get: failGet(isLegacyDashboard, errInjected)},
	} {
		t.Run(name, func(t *testing.T) {
			cr := testPlatformMonitoring()
			r := newInterceptedDashboardReconciler(t, funcs, cr)
			require.ErrorIs(t, r.handleGrafanaDashboard(testDashboardAsset, cr), errInjected)
		})
	}
}

func TestHandleGrafanaDashboardFailsWithoutOwnerScheme(t *testing.T) {
	cr := testPlatformMonitoring()
	scheme := runtime.NewScheme()
	require.NoError(t, grafv1.AddToScheme(scheme))
	r := &GrafanaOperatorReconciler{
		ComponentReconciler: &utils.ComponentReconciler{
			Client: fake.NewClientBuilder().WithScheme(scheme).Build(),
			Scheme: scheme,
			Log:    utils.Logger("grafanaoperator_dashboard_errors_test"),
		},
	}
	// Without PlatformMonitoring in the scheme, the owner reference cannot be set.
	require.Error(t, r.handleGrafanaDashboard(testDashboardAsset, cr))
}

func TestDeleteGrafanaDashboardReturnsAPIErrors(t *testing.T) {
	for name, funcs := range map[string]interceptor.Funcs{
		"deleting the legacy dashboard fails": {
			Delete: func(_ context.Context, c client.WithWatch, obj client.Object, opts ...client.DeleteOption) error {
				if isLegacyDashboard(obj) {
					return errInjected
				}
				return c.Delete(context.TODO(), obj, opts...)
			},
		},
		"reading the v1beta1 dashboard fails": {Get: failGet(isNativeDashboard, errInjected)},
	} {
		t.Run(name, func(t *testing.T) {
			cr := testPlatformMonitoring()
			r := newInterceptedDashboardReconciler(t, funcs, cr)
			require.ErrorIs(t, r.deleteGrafanaDashboard(testDashboardAsset, cr), errInjected)
		})
	}
}

func TestRunDeletesOrphanedConvertedDashboards(t *testing.T) {
	cr := testPlatformMonitoring()
	cr.Spec.Grafana = nil
	orphan := builtInCopy(cr, "version-exporter")
	r := newGrafanaOperatorWorkloadTestReconciler(t, cr, orphan)

	require.NoError(t, r.Run(cr))

	err := r.Client.Get(t.Context(), client.ObjectKeyFromObject(orphan), &grafv1.GrafanaDashboard{})
	assert.True(t, apierrors.IsNotFound(err), "Get orphaned copy error = %v", err)
}

func TestRunLogsOrphanedCopyCleanupError(t *testing.T) {
	cr := testPlatformMonitoring()
	cr.Spec.Grafana = nil
	r := newInterceptedDashboardReconciler(t, interceptor.Funcs{
		List: func(context.Context, client.WithWatch, client.ObjectList, ...client.ListOption) error {
			return errInjected
		},
	}, cr)
	r.Dc = newGrafanaOperatorWorkloadTestReconciler(t).Dc

	// A failed cleanup is logged and does not fail the component.
	require.NoError(t, r.Run(cr))
}
