package nodeexporter

import (
	"encoding/json"
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestHandleDaemonSetSkipsAPIDefaultedPodSpec(t *testing.T) {
	cr := newNodeExporterPlatformMonitoring(true)
	cr.Spec.NodeExporter.Resources = corev1.ResourceRequirements{
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("500m"),
			corev1.ResourceMemory: resource.MustParse("100Mi"),
		},
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("50m"),
			corev1.ResourceMemory: resource.MustParse("50Mi"),
		},
	}
	runAsUser, fsGroup := int64(2000), int64(2000)
	cr.Spec.NodeExporter.SecurityContext = &monv1.SecurityContext{
		RunAsUser: &runAsUser,
		FSGroup:   &fsGroup,
	}

	desired, err := nodeExporterDaemonSet(cr)
	require.NoError(t, err)

	live := desired.DeepCopy()
	applyAPIDefaultedPodSpec(&live.Spec.Template.Spec)
	raw, err := json.Marshal(live)
	require.NoError(t, err)
	live = &appsv1.DaemonSet{}
	require.NoError(t, json.Unmarshal(raw, live))
	live.ResourceVersion = "1"

	r, kubeClient := newNodeExporterTestReconciler(t, false, cr, live)
	require.NoError(t, r.handleDaemonSet(cr))

	got := &appsv1.DaemonSet{}
	require.NoError(t, kubeClient.Get(t.Context(), client.ObjectKeyFromObject(desired), got))
	assert.Equal(t, "1", got.ResourceVersion, "API-defaulted PodSpec fields must not trigger Update")
	assert.Equal(t, corev1.DNSClusterFirst, got.Spec.Template.Spec.DNSPolicy)
	assert.Equal(t, "default-scheduler", got.Spec.Template.Spec.SchedulerName)
}

func applyAPIDefaultedPodSpec(spec *corev1.PodSpec) {
	if spec.DNSPolicy == "" {
		spec.DNSPolicy = corev1.DNSClusterFirst
	}
	if spec.RestartPolicy == "" {
		spec.RestartPolicy = corev1.RestartPolicyAlways
	}
	if spec.SchedulerName == "" {
		spec.SchedulerName = "default-scheduler"
	}
	if spec.TerminationGracePeriodSeconds == nil {
		tgs := int64(30)
		spec.TerminationGracePeriodSeconds = &tgs
	}
	if spec.DeprecatedServiceAccount == "" {
		spec.DeprecatedServiceAccount = spec.ServiceAccountName
	}
	for i := range spec.Volumes {
		if hp := spec.Volumes[i].HostPath; hp != nil && hp.Type == nil {
			t := corev1.HostPathUnset
			hp.Type = &t
		}
	}
	applyAPIDefaultedContainerFields(spec.Containers)
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
		if probe := c.ReadinessProbe; probe != nil {
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
