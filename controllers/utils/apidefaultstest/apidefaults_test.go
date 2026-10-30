package apidefaultstest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/ptr"
)

func TestPodSpecAppliesDefaults(t *testing.T) {
	spec := &corev1.PodSpec{
		ServiceAccountName: "sa",
		Volumes: []corev1.Volume{
			{Name: "host", VolumeSource: corev1.VolumeSource{HostPath: &corev1.HostPathVolumeSource{Path: "/proc"}}},
			{Name: "empty", VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}}},
		},
		Containers: []corev1.Container{{
			Name:           "app",
			Ports:          []corev1.ContainerPort{{ContainerPort: 80}},
			LivenessProbe:  &corev1.Probe{ProbeHandler: corev1.ProbeHandler{HTTPGet: &corev1.HTTPGetAction{Path: "/"}}},
			ReadinessProbe: &corev1.Probe{PeriodSeconds: 5},
		}},
	}

	PodSpec(spec)

	assert.Equal(t, corev1.DNSClusterFirst, spec.DNSPolicy)
	assert.Equal(t, corev1.RestartPolicyAlways, spec.RestartPolicy)
	assert.Equal(t, corev1.DefaultSchedulerName, spec.SchedulerName)
	assert.Equal(t, ptr.To(int64(30)), spec.TerminationGracePeriodSeconds)
	assert.Equal(t, "sa", spec.DeprecatedServiceAccount)
	assert.Equal(t, ptr.To(corev1.HostPathUnset), spec.Volumes[0].HostPath.Type)

	container := spec.Containers[0]
	assert.Equal(t, corev1.ProtocolTCP, container.Ports[0].Protocol)
	assert.Equal(t, TerminationMessagePath, container.TerminationMessagePath)
	assert.Equal(t, corev1.TerminationMessageReadFile, container.TerminationMessagePolicy)
	assert.Equal(t, int32(10), container.LivenessProbe.PeriodSeconds)
	assert.Equal(t, int32(1), container.LivenessProbe.SuccessThreshold)
	assert.Equal(t, int32(3), container.LivenessProbe.FailureThreshold)
	assert.Equal(t, corev1.URISchemeHTTP, container.LivenessProbe.HTTPGet.Scheme)
	assert.Equal(t, int32(5), container.ReadinessProbe.PeriodSeconds, "explicit values are kept")
}

func TestPodSpecKeepsExplicitValues(t *testing.T) {
	spec := &corev1.PodSpec{
		DNSPolicy:                     corev1.DNSDefault,
		TerminationGracePeriodSeconds: ptr.To(int64(5)),
	}

	PodSpec(spec)

	assert.Equal(t, corev1.DNSDefault, spec.DNSPolicy)
	assert.Equal(t, ptr.To(int64(5)), spec.TerminationGracePeriodSeconds)
}

func TestPortDefaults(t *testing.T) {
	ports := []corev1.ServicePort{{Port: 80}, {Port: 53, Protocol: corev1.ProtocolUDP}}
	ServicePorts(ports)
	assert.Equal(t, corev1.ProtocolTCP, ports[0].Protocol)
	assert.Equal(t, corev1.ProtocolUDP, ports[1].Protocol)

	subsets := []corev1.EndpointSubset{{Ports: []corev1.EndpointPort{{Port: 80}}}} //nolint:staticcheck // SA1019: tests the v1 Endpoints helper.
	EndpointPorts(subsets)
	assert.Equal(t, corev1.ProtocolTCP, subsets[0].Ports[0].Protocol)
}
