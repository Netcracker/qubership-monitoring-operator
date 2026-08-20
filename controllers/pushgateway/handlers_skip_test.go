package pushgateway

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

func TestHandleDeploymentSkipsAPIDefaultedContainers(t *testing.T) {
	cr := skipTestPlatformMonitoring()
	desired, err := pushgatewayDeployment(cr)
	require.NoError(t, err)
	live := desired.DeepCopy()
	live.ResourceVersion = "1"
	live.SetAnnotations(map[string]string{"deployment.kubernetes.io/revision": "1"})
	live.Spec.Template.Spec.SecurityContext = &corev1.PodSecurityContext{}
	applyAPIDefaultedContainerFields(live.Spec.Template.Spec.Containers)

	r := newPushgatewaySkipTestReconciler(t, cr, live)
	require.NoError(t, r.handleDeployment(cr))

	got := &appsv1.Deployment{}
	require.NoError(t, r.Client.Get(t.Context(), client.ObjectKeyFromObject(desired), got))
	assert.Equal(t, "1", got.ResourceVersion, "API-defaulted container fields must not trigger Update")
	assert.Equal(t, map[string]string{"deployment.kubernetes.io/revision": "1"}, got.GetAnnotations())
	require.NotEmpty(t, got.Spec.Template.Spec.Containers)
	assert.Equal(t, int32(0), got.Spec.Template.Spec.Containers[0].Ports[0].HostPort)
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
			Pushgateway: &monv1.Pushgateway{
				Install: &install,
				Image:   "docker.io/prom/pushgateway:v1.11.2",
				Port:    9091,
			},
		},
	}
}

func newPushgatewaySkipTestReconciler(t *testing.T, objs ...client.Object) *PushgatewayReconciler {
	t.Helper()
	scheme := runtime.NewScheme()
	require.NoError(t, monv1.AddToScheme(scheme))
	require.NoError(t, appsv1.AddToScheme(scheme))
	return &PushgatewayReconciler{
		ComponentReconciler: &utils.ComponentReconciler{
			Client: fake.NewClientBuilder().WithScheme(scheme).WithObjects(objs...).Build(),
			Scheme: scheme,
			Log:    utils.Logger("pushgateway_skip_test"),
		},
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
		for _, probe := range []*corev1.Probe{c.LivenessProbe, c.ReadinessProbe} {
			if probe == nil {
				continue
			}
			if probe.PeriodSeconds == 0 {
				probe.PeriodSeconds = 10
			}
			if probe.SuccessThreshold == 0 {
				probe.SuccessThreshold = 1
			}
			if probe.FailureThreshold == 0 {
				probe.FailureThreshold = 3
			}
			if probe.HTTPGet != nil && probe.HTTPGet.Scheme == "" {
				probe.HTTPGet.Scheme = corev1.URISchemeHTTP
			}
		}
	}
}
