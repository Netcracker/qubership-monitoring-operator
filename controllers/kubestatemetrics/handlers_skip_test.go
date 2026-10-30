package kubestatemetrics

import (
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils/apidefaultstest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	fakediscovery "k8s.io/client-go/discovery/fake"
	k8stesting "k8s.io/client-go/testing"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestHandleDeploymentSkipsAPIDefaultedContainers(t *testing.T) {
	cr := skipTestPlatformMonitoring()
	desired, err := kubeStateMetricsDeployment(cr, false, false)
	require.NoError(t, err)

	live := desired.DeepCopy()
	live.ResourceVersion = "1"
	apidefaultstest.Containers(live.Spec.Template.Spec.Containers)
	live.Spec.Template.Spec.DNSPolicy = corev1.DNSClusterFirst
	live.Spec.Template.Spec.RestartPolicy = corev1.RestartPolicyAlways
	live.Spec.Template.Spec.SchedulerName = "default-scheduler"

	r := newKubeStateMetricsSkipTestReconciler(t, cr, live)
	require.NoError(t, r.handleDeployment(cr))

	got := &appsv1.Deployment{}
	require.NoError(t, r.Client.Get(t.Context(), client.ObjectKeyFromObject(desired), got))
	assert.Equal(t, "1", got.ResourceVersion, "API-defaulted container and PodSpec fields must not trigger Update")
}

func TestHandleDeploymentUpdatesChangedImage(t *testing.T) {
	cr := skipTestPlatformMonitoring()
	desired, err := kubeStateMetricsDeployment(cr, false, false)
	require.NoError(t, err)
	live := desired.DeepCopy()
	live.ResourceVersion = "1"
	require.NotEmpty(t, live.Spec.Template.Spec.Containers)
	live.Spec.Template.Spec.Containers[0].Image = "stale:old"

	r := newKubeStateMetricsSkipTestReconciler(t, cr, live)
	require.NoError(t, r.handleDeployment(cr))

	got := &appsv1.Deployment{}
	require.NoError(t, r.Client.Get(t.Context(), client.ObjectKeyFromObject(desired), got))
	assert.NotEqual(t, "1", got.ResourceVersion, "image change must trigger Update")
	require.NotEmpty(t, got.Spec.Template.Spec.Containers)
	assert.Equal(t, cr.Spec.KubeStateMetrics.Image, got.Spec.Template.Spec.Containers[0].Image)
}

func skipTestPlatformMonitoring() *monv1.PlatformMonitoring {
	install := true
	return &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "platformmonitoring",
			Namespace: "monitoring",
			UID:       types.UID("platform-monitoring-uid"),
		},
		Spec: monv1.PlatformMonitoringSpec{
			KubeStateMetrics: &monv1.KubeStateMetrics{
				Install: &install,
				Image:   "registry.k8s.io/kube-state-metrics/kube-state-metrics:v2.15.0",
			},
		},
	}
}

func newKubeStateMetricsSkipTestReconciler(t *testing.T, objs ...client.Object) *KubeStateMetricsReconciler {
	t.Helper()
	scheme := runtime.NewScheme()
	require.NoError(t, monv1.AddToScheme(scheme))
	require.NoError(t, appsv1.AddToScheme(scheme))
	return NewKubeStateMetricsReconciler(
		fake.NewClientBuilder().WithScheme(scheme).WithObjects(objs...).Build(),
		scheme,
		&fakediscovery.FakeDiscovery{Fake: &k8stesting.Fake{}},
	)
}
