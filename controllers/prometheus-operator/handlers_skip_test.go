package prometheus_operator

import (
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestHandleServiceSkipsAPIDefaultedPortProtocol(t *testing.T) {
	cr := skipTestPlatformMonitoring()
	desired, err := prometheusOperatorService(cr)
	require.NoError(t, err)
	desired.Spec.Selector = map[string]string{
		"app.kubernetes.io/name": utils.TruncLabel(utils.PrometheusOperatorComponentName),
	}
	live := desired.DeepCopy()
	live.ResourceVersion = "1"
	applyAPIDefaultedServicePorts(live.Spec.Ports)

	r := newPrometheusOperatorSkipTestReconciler(t, cr, live)
	require.NoError(t, r.handleService(cr))

	got := &corev1.Service{}
	require.NoError(t, r.Client.Get(t.Context(), client.ObjectKeyFromObject(desired), got))
	assert.Equal(t, "1", got.ResourceVersion, "API-defaulted Service port protocol must not trigger Update")
	require.NotEmpty(t, got.Spec.Ports)
	assert.Equal(t, corev1.ProtocolTCP, got.Spec.Ports[0].Protocol)
}

func TestHandleDeploymentSkipsAPIDefaultedContainers(t *testing.T) {
	cr := skipTestPlatformMonitoring()
	desired, err := prometheusOperatorDeployment(cr)
	require.NoError(t, err)
	live := desired.DeepCopy()
	live.ResourceVersion = "1"
	live.Spec.Template.Spec.SecurityContext = &corev1.PodSecurityContext{}
	applyAPIDefaultedContainerFields(live.Spec.Template.Spec.Containers)

	r := newPrometheusOperatorSkipTestReconciler(t, cr, live)
	require.NoError(t, r.handleDeployment(cr))

	got := &appsv1.Deployment{}
	require.NoError(t, r.Client.Get(t.Context(), client.ObjectKeyFromObject(desired), got))
	assert.Equal(t, "1", got.ResourceVersion, "API-defaulted container fields must not trigger Update")
}

func TestHandleServiceUpdatesChangedPort(t *testing.T) {
	cr := skipTestPlatformMonitoring()
	desired, err := prometheusOperatorService(cr)
	require.NoError(t, err)
	desired.Spec.Selector = map[string]string{
		"app.kubernetes.io/name": utils.TruncLabel(utils.PrometheusOperatorComponentName),
	}
	live := desired.DeepCopy()
	live.ResourceVersion = "1"
	require.NotEmpty(t, live.Spec.Ports)
	live.Spec.Ports[0].Port = 9999

	r := newPrometheusOperatorSkipTestReconciler(t, cr, live)
	require.NoError(t, r.handleService(cr))

	got := &corev1.Service{}
	require.NoError(t, r.Client.Get(t.Context(), client.ObjectKeyFromObject(desired), got))
	assert.NotEqual(t, "1", got.ResourceVersion, "port change must trigger Update")
	require.NotEmpty(t, got.Spec.Ports)
	assert.Equal(t, int32(8080), got.Spec.Ports[0].Port)
}

func TestHandleDeploymentUpdatesChangedImage(t *testing.T) {
	cr := skipTestPlatformMonitoring()
	desired, err := prometheusOperatorDeployment(cr)
	require.NoError(t, err)
	live := desired.DeepCopy()
	live.ResourceVersion = "1"
	require.NotEmpty(t, live.Spec.Template.Spec.Containers)
	live.Spec.Template.Spec.Containers[0].Image = "stale:old"

	r := newPrometheusOperatorSkipTestReconciler(t, cr, live)
	require.NoError(t, r.handleDeployment(cr))

	got := &appsv1.Deployment{}
	require.NoError(t, r.Client.Get(t.Context(), client.ObjectKeyFromObject(desired), got))
	assert.NotEqual(t, "1", got.ResourceVersion, "image change must trigger Update")
	require.NotEmpty(t, got.Spec.Template.Spec.Containers)
	assert.Equal(t, cr.Spec.Prometheus.Operator.Image, got.Spec.Template.Spec.Containers[0].Image)
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
			Prometheus: &monv1.Prometheus{
				Install:             &install,
				Image:               "docker.io/prom/prometheus:v3.13.1",
				ConfigReloaderImage: "quay.io/prometheus-operator/prometheus-config-reloader:v0.93.0",
				Operator: monv1.PrometheusOperator{
					Image: "quay.io/prometheus-operator/prometheus-operator:v0.93.0",
				},
			},
		},
	}
}

func newPrometheusOperatorSkipTestReconciler(t *testing.T, objs ...client.Object) *PrometheusOperatorReconciler {
	t.Helper()
	scheme := runtime.NewScheme()
	require.NoError(t, monv1.AddToScheme(scheme))
	require.NoError(t, appsv1.AddToScheme(scheme))
	require.NoError(t, corev1.AddToScheme(scheme))
	return NewPrometheusOperatorReconciler(
		fake.NewClientBuilder().WithScheme(scheme).WithObjects(objs...).Build(),
		scheme,
	)
}

func applyAPIDefaultedServicePorts(ports []corev1.ServicePort) {
	for i := range ports {
		if ports[i].Protocol == "" {
			ports[i].Protocol = corev1.ProtocolTCP
		}
	}
}

func applyAPIDefaultedContainerFields(containers []corev1.Container) {
	for i := range containers {
		c := &containers[i]
		for j := range c.Ports {
			if c.Ports[j].Protocol == "" {
				c.Ports[j].Protocol = corev1.ProtocolTCP
			}
		}
		if c.TerminationMessagePath == "" {
			c.TerminationMessagePath = "/dev/termination-log"
		}
		if c.TerminationMessagePolicy == "" {
			c.TerminationMessagePolicy = corev1.TerminationMessageReadFile
		}
	}
}
