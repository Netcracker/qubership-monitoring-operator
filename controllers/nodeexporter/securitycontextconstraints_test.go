package nodeexporter

import (
	"context"
	"errors"
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	secv1 "github.com/openshift/api/security/v1"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	fakediscovery "k8s.io/client-go/discovery/fake"
	k8stesting "k8s.io/client-go/testing"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
)

func TestHandleSecurityContextConstraints(t *testing.T) {
	t.Run("creates SCC", func(t *testing.T) {
		reconciler, kubeClient := newNodeExporterTestReconciler(t, true)
		platformMonitoring := newNodeExporterPlatformMonitoring(true)

		require.NoError(t, reconciler.handleSecurityContextConstraints(platformMonitoring))

		actual := &secv1.SecurityContextConstraints{}
		require.NoError(t, kubeClient.Get(
			context.Background(),
			client.ObjectKey{Name: nodeExporterSecurityContextConstraintsName(platformMonitoring)},
			actual,
		))
		assert.True(t, actual.AllowHostNetwork)
		assert.True(t, actual.AllowHostPID)
		assert.Contains(t, actual.Labels, "app.kubernetes.io/managed-by")
		assert.True(t, isNodeExporterSCCOwnedBy(actual, platformMonitoring))
	})

	t.Run("updates policy and preserves subjects", func(t *testing.T) {
		platformMonitoring := newNodeExporterPlatformMonitoring(true)
		desired, err := nodeExporterSecurityContextConstraints(platformMonitoring)
		require.NoError(t, err)

		existing := desired.DeepCopy()
		existing.AllowHostNetwork = false
		existing.AllowHostPID = false
		existing.AllowHostDirVolumePlugin = false
		existing.Volumes = nil
		existing.SeccompProfiles = nil
		existing.Users = []string{"system:serviceaccount:custom:node-exporter"}
		existing.Groups = []string{"system:serviceaccounts:custom"}

		reconciler, kubeClient := newNodeExporterTestReconciler(t, true, existing)
		require.NoError(t, reconciler.handleSecurityContextConstraints(platformMonitoring))

		actual := &secv1.SecurityContextConstraints{}
		require.NoError(t, kubeClient.Get(
			context.Background(),
			client.ObjectKey{Name: nodeExporterSecurityContextConstraintsName(platformMonitoring)},
			actual,
		))
		assert.Equal(t, desired.AllowHostNetwork, actual.AllowHostNetwork)
		assert.Equal(t, desired.AllowHostPID, actual.AllowHostPID)
		assert.Equal(t, desired.AllowHostDirVolumePlugin, actual.AllowHostDirVolumePlugin)
		assert.Equal(t, desired.Volumes, actual.Volumes)
		assert.Equal(t, desired.SeccompProfiles, actual.SeccompProfiles)
		assert.Equal(t, existing.Users, actual.Users)
		assert.Equal(t, existing.Groups, actual.Groups)
	})

	t.Run("refuses to update a foreign SCC", func(t *testing.T) {
		platformMonitoring := newNodeExporterPlatformMonitoring(true)
		foreign, err := nodeExporterSecurityContextConstraints(platformMonitoring)
		require.NoError(t, err)
		foreign.SetAnnotations(nil)
		foreign.AllowHostNetwork = false

		reconciler, kubeClient := newNodeExporterTestReconciler(t, true, foreign)
		err = reconciler.handleSecurityContextConstraints(platformMonitoring)
		require.ErrorContains(t, err, "is not owned by PlatformMonitoring monitoring/platform-monitoring")

		actual := &secv1.SecurityContextConstraints{}
		require.NoError(t, kubeClient.Get(context.Background(), client.ObjectKeyFromObject(foreign), actual))
		assert.False(t, actual.AllowHostNetwork)
		assert.Empty(t, actual.Annotations)
	})

	t.Run("deletes SCC idempotently", func(t *testing.T) {
		platformMonitoring := newNodeExporterPlatformMonitoring(true)
		existing, err := nodeExporterSecurityContextConstraints(platformMonitoring)
		require.NoError(t, err)

		reconciler, kubeClient := newNodeExporterTestReconciler(t, true, existing)
		require.NoError(t, reconciler.deleteSecurityContextConstraints(platformMonitoring))
		require.NoError(t, reconciler.deleteSecurityContextConstraints(platformMonitoring))

		err = kubeClient.Get(
			context.Background(),
			client.ObjectKey{Name: nodeExporterSecurityContextConstraintsName(platformMonitoring)},
			&secv1.SecurityContextConstraints{},
		)
		assert.True(t, apierrors.IsNotFound(err))
	})

	t.Run("preserves a foreign SCC during deletion", func(t *testing.T) {
		platformMonitoring := newNodeExporterPlatformMonitoring(true)
		foreign, err := nodeExporterSecurityContextConstraints(platformMonitoring)
		require.NoError(t, err)
		foreign.SetAnnotations(map[string]string{
			nodeExporterSCCOwnerNameAnnotation:      "another-platform-monitoring",
			nodeExporterSCCOwnerNamespaceAnnotation: platformMonitoring.GetNamespace(),
		})

		reconciler, kubeClient := newNodeExporterTestReconciler(t, true, foreign)
		require.NoError(t, reconciler.deleteSecurityContextConstraints(platformMonitoring))
		assertResourceExists(t, kubeClient, foreign)
	})

	t.Run("returns create error", func(t *testing.T) {
		reconciler, kubeClient := newNodeExporterTestReconciler(t, true)
		expectedErr := errors.New("create SCC")
		reconciler.Client = interceptNodeExporterClient(t, kubeClient, interceptor.Funcs{
			Create: func(
				_ context.Context,
				_ client.WithWatch,
				_ client.Object,
				_ ...client.CreateOption,
			) error {
				return expectedErr
			},
		})

		assert.ErrorIs(
			t,
			reconciler.handleSecurityContextConstraints(newNodeExporterPlatformMonitoring(true)),
			expectedErr,
		)
	})

	t.Run("returns get error while reconciling", func(t *testing.T) {
		reconciler, kubeClient := newNodeExporterTestReconciler(t, true)
		expectedErr := errors.New("get SCC")
		reconciler.Client = interceptNodeExporterClient(t, kubeClient, interceptor.Funcs{
			Get: func(
				_ context.Context,
				_ client.WithWatch,
				_ client.ObjectKey,
				_ client.Object,
				_ ...client.GetOption,
			) error {
				return expectedErr
			},
		})

		assert.ErrorIs(
			t,
			reconciler.handleSecurityContextConstraints(newNodeExporterPlatformMonitoring(true)),
			expectedErr,
		)
	})

	t.Run("returns update error", func(t *testing.T) {
		platformMonitoring := newNodeExporterPlatformMonitoring(true)
		existing, err := nodeExporterSecurityContextConstraints(platformMonitoring)
		require.NoError(t, err)
		existing.AllowHostNetwork = false
		reconciler, kubeClient := newNodeExporterTestReconciler(t, true, existing)
		expectedErr := errors.New("update SCC")
		reconciler.Client = interceptNodeExporterClient(t, kubeClient, interceptor.Funcs{
			Update: func(
				_ context.Context,
				_ client.WithWatch,
				_ client.Object,
				_ ...client.UpdateOption,
			) error {
				return expectedErr
			},
		})

		assert.ErrorIs(
			t,
			reconciler.handleSecurityContextConstraints(platformMonitoring),
			expectedErr,
		)
	})

	t.Run("returns get error while deleting", func(t *testing.T) {
		reconciler, kubeClient := newNodeExporterTestReconciler(t, true)
		expectedErr := errors.New("get SCC")
		reconciler.Client = interceptNodeExporterClient(t, kubeClient, interceptor.Funcs{
			Get: func(
				_ context.Context,
				_ client.WithWatch,
				_ client.ObjectKey,
				_ client.Object,
				_ ...client.GetOption,
			) error {
				return expectedErr
			},
		})

		assert.ErrorIs(
			t,
			reconciler.deleteSecurityContextConstraints(newNodeExporterPlatformMonitoring(true)),
			expectedErr,
		)
	})

	t.Run("returns delete error", func(t *testing.T) {
		platformMonitoring := newNodeExporterPlatformMonitoring(true)
		existing, err := nodeExporterSecurityContextConstraints(platformMonitoring)
		require.NoError(t, err)
		reconciler, kubeClient := newNodeExporterTestReconciler(t, true, existing)
		expectedErr := errors.New("delete SCC")
		reconciler.Client = interceptNodeExporterClient(t, kubeClient, interceptor.Funcs{
			Delete: func(
				_ context.Context,
				_ client.WithWatch,
				_ client.Object,
				_ ...client.DeleteOption,
			) error {
				return expectedErr
			},
		})

		assert.ErrorIs(
			t,
			reconciler.deleteSecurityContextConstraints(platformMonitoring),
			expectedErr,
		)
	})
}

func TestNodeExporterReconcilesSecurityContextConstraints(t *testing.T) {
	previousPrivilegedRights := utils.PrivilegedRights
	utils.PrivilegedRights = true
	t.Cleanup(func() {
		utils.PrivilegedRights = previousPrivilegedRights
	})

	t.Run("creates and uninstalls SCC when API is available", func(t *testing.T) {
		reconciler, kubeClient := newNodeExporterTestReconciler(t, true)
		platformMonitoring := newNodeExporterPlatformMonitoring(true)

		require.NoError(t, reconciler.Run(platformMonitoring))
		require.NoError(t, reconciler.Run(platformMonitoring))
		assertResourceExists(t, kubeClient, &secv1.SecurityContextConstraints{
			ObjectMeta: metav1.ObjectMeta{Name: nodeExporterSecurityContextConstraintsName(platformMonitoring)},
		})

		install := false
		platformMonitoring.Spec.NodeExporter.Install = &install
		require.NoError(t, reconciler.Run(platformMonitoring))
		assertResourceDoesNotExist(t, kubeClient, &secv1.SecurityContextConstraints{
			ObjectMeta: metav1.ObjectMeta{Name: nodeExporterSecurityContextConstraintsName(platformMonitoring)},
		})
	})

	t.Run("returns SCC creation error", func(t *testing.T) {
		reconciler, kubeClient := newNodeExporterTestReconciler(t, true)
		expectedErr := errors.New("create SCC")
		reconciler.Client = interceptNodeExporterClient(t, kubeClient, interceptor.Funcs{
			Create: func(
				ctx context.Context,
				underlyingClient client.WithWatch,
				object client.Object,
				opts ...client.CreateOption,
			) error {
				if _, ok := object.(*secv1.SecurityContextConstraints); ok {
					return expectedErr
				}
				return underlyingClient.Create(ctx, object, opts...)
			},
		})

		assert.ErrorIs(t, reconciler.Run(newNodeExporterPlatformMonitoring(true)), expectedErr)
	})

	t.Run("continues uninstall after SCC deletion error", func(t *testing.T) {
		platformMonitoring := newNodeExporterPlatformMonitoring(true)
		existing, err := nodeExporterSecurityContextConstraints(platformMonitoring)
		require.NoError(t, err)
		serviceAccount, err := nodeExporterServiceAccount(platformMonitoring)
		require.NoError(t, err)
		reconciler, kubeClient := newNodeExporterTestReconciler(t, true, existing, serviceAccount)
		reconciler.Client = interceptNodeExporterClient(t, kubeClient, interceptor.Funcs{
			Delete: func(
				ctx context.Context,
				underlyingClient client.WithWatch,
				object client.Object,
				opts ...client.DeleteOption,
			) error {
				if _, ok := object.(*secv1.SecurityContextConstraints); ok {
					return errors.New("delete SCC")
				}
				return underlyingClient.Delete(ctx, object, opts...)
			},
		})

		reconciler.uninstall(platformMonitoring)
		assertResourceExists(t, kubeClient, existing)
		assertResourceDoesNotExist(t, kubeClient, serviceAccount)
	})

	t.Run("deletes an owned SCC after setup is disabled", func(t *testing.T) {
		reconciler, kubeClient := newNodeExporterTestReconciler(t, true)
		platformMonitoring := newNodeExporterPlatformMonitoring(true)

		require.NoError(t, reconciler.Run(platformMonitoring))
		assertResourceExists(t, kubeClient, &secv1.SecurityContextConstraints{
			ObjectMeta: metav1.ObjectMeta{Name: nodeExporterSecurityContextConstraintsName(platformMonitoring)},
		})

		platformMonitoring.Spec.NodeExporter.SetupSecurityContext = false
		require.NoError(t, reconciler.Run(platformMonitoring))
		assertResourceDoesNotExist(t, kubeClient, &secv1.SecurityContextConstraints{
			ObjectMeta: metav1.ObjectMeta{Name: nodeExporterSecurityContextConstraintsName(platformMonitoring)},
		})
	})

	t.Run("preserves a foreign SCC when setup is disabled", func(t *testing.T) {
		platformMonitoring := newNodeExporterPlatformMonitoring(false)
		foreign, err := nodeExporterSecurityContextConstraints(platformMonitoring)
		require.NoError(t, err)
		foreign.SetAnnotations(nil)
		reconciler, kubeClient := newNodeExporterTestReconciler(t, true, foreign)

		require.NoError(t, reconciler.Run(platformMonitoring))
		assertResourceExists(t, kubeClient, foreign)
	})

	t.Run("preserves a foreign SCC during uninstall", func(t *testing.T) {
		platformMonitoring := newNodeExporterPlatformMonitoring(true)
		foreign, err := nodeExporterSecurityContextConstraints(platformMonitoring)
		require.NoError(t, err)
		foreign.SetAnnotations(nil)
		reconciler, kubeClient := newNodeExporterTestReconciler(t, true, foreign)

		install := false
		platformMonitoring.Spec.NodeExporter.Install = &install
		require.NoError(t, reconciler.Run(platformMonitoring))
		assertResourceExists(t, kubeClient, foreign)
	})

	t.Run("skips SCC when API is unavailable", func(t *testing.T) {
		reconciler, kubeClient := newNodeExporterTestReconciler(t, false)
		platformMonitoring := newNodeExporterPlatformMonitoring(true)

		require.NoError(t, reconciler.Run(platformMonitoring))
		assertResourceDoesNotExist(t, kubeClient, &secv1.SecurityContextConstraints{
			ObjectMeta: metav1.ObjectMeta{Name: nodeExporterSecurityContextConstraintsName(platformMonitoring)},
		})
	})
}

func newNodeExporterTestReconciler(
	t *testing.T,
	hasSecurityContextConstraints bool,
	objects ...client.Object,
) (*NodeExporterReconciler, client.Client) {
	t.Helper()

	scheme := runtime.NewScheme()
	require.NoError(t, monv1.AddToScheme(scheme))
	require.NoError(t, secv1.Install(scheme))
	require.NoError(t, appsv1.AddToScheme(scheme))
	require.NoError(t, corev1.AddToScheme(scheme))
	require.NoError(t, rbacv1.AddToScheme(scheme))
	require.NoError(t, promv1.AddToScheme(scheme))

	kubeClient := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(objects...).
		Build()
	discoveryClient := &fakediscovery.FakeDiscovery{Fake: &k8stesting.Fake{}}
	if hasSecurityContextConstraints {
		discoveryClient.Resources = []*metav1.APIResourceList{{
			GroupVersion: secv1.GroupVersion.String(),
			APIResources: []metav1.APIResource{{Kind: "SecurityContextConstraints"}},
		}}
	}

	return NewNodeExporterReconciler(kubeClient, scheme, discoveryClient), kubeClient
}

func newNodeExporterPlatformMonitoring(setupSecurityContext bool) *monv1.PlatformMonitoring {
	install := true
	return &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "platform-monitoring",
			Namespace: "monitoring",
			UID:       types.UID("platform-monitoring"),
		},
		Spec: monv1.PlatformMonitoringSpec{
			NodeExporter: &monv1.NodeExporter{
				Install:              &install,
				SetupSecurityContext: setupSecurityContext,
				Image:                "quay.io/prometheus/node-exporter:v1.9.1",
				Port:                 9100,
			},
		},
	}
}

func assertResourceExists(t *testing.T, kubeClient client.Client, object client.Object) {
	t.Helper()
	assert.NoError(t, kubeClient.Get(context.Background(), client.ObjectKeyFromObject(object), object))
}

func assertResourceDoesNotExist(t *testing.T, kubeClient client.Client, object client.Object) {
	t.Helper()
	err := kubeClient.Get(context.Background(), client.ObjectKeyFromObject(object), object)
	assert.True(t, apierrors.IsNotFound(err), "expected resource to be absent, got %v", err)
}

func interceptNodeExporterClient(
	t *testing.T,
	kubeClient client.Client,
	funcs interceptor.Funcs,
) client.Client {
	t.Helper()
	clientWithWatch, ok := kubeClient.(client.WithWatch)
	require.True(t, ok)
	return interceptor.NewClient(clientWithWatch, funcs)
}
