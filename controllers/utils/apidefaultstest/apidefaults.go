// Package apidefaultstest fills in the defaults that the Kubernetes API server
// applies to stored objects. Tests use it to build realistic readbacks for
// fake clients, which do not apply API defaults themselves.
package apidefaultstest

import corev1 "k8s.io/api/core/v1"

// TerminationMessagePath is the terminationMessagePath that the API server sets
// on containers.
const TerminationMessagePath = "/dev/termination-log"

// PodSpec applies the pod-level and container-level API defaults to spec.
func PodSpec(spec *corev1.PodSpec) {
	setIfZero(&spec.DNSPolicy, corev1.DNSClusterFirst)
	setIfZero(&spec.RestartPolicy, corev1.RestartPolicyAlways)
	setIfZero(&spec.SchedulerName, corev1.DefaultSchedulerName)
	if spec.TerminationGracePeriodSeconds == nil {
		grace := int64(corev1.DefaultTerminationGracePeriodSeconds)
		spec.TerminationGracePeriodSeconds = &grace
	}
	setIfZero(&spec.DeprecatedServiceAccount, spec.ServiceAccountName)
	for i := range spec.Volumes {
		if hostPath := spec.Volumes[i].HostPath; hostPath != nil && hostPath.Type == nil {
			unset := corev1.HostPathUnset
			hostPath.Type = &unset
		}
	}
	Containers(spec.Containers)
}

// Containers applies the API defaults to container ports, termination message
// settings, and liveness and readiness probes.
func Containers(containers []corev1.Container) {
	for i := range containers {
		container := &containers[i]
		for j := range container.Ports {
			setIfZero(&container.Ports[j].Protocol, corev1.ProtocolTCP)
		}
		setIfZero(&container.TerminationMessagePath, TerminationMessagePath)
		setIfZero(&container.TerminationMessagePolicy, corev1.TerminationMessageReadFile)
		probe(container.LivenessProbe)
		probe(container.ReadinessProbe)
	}
}

// ServicePorts sets the default protocol on Service ports.
func ServicePorts(ports []corev1.ServicePort) {
	for i := range ports {
		setIfZero(&ports[i].Protocol, corev1.ProtocolTCP)
	}
}

// EndpointPorts sets the default protocol on the ports of every Endpoints subset.
func EndpointPorts(subsets []corev1.EndpointSubset) { //nolint:staticcheck // SA1019: the operator still manages v1 Endpoints.
	for i := range subsets {
		for j := range subsets[i].Ports {
			setIfZero(&subsets[i].Ports[j].Protocol, corev1.ProtocolTCP)
		}
	}
}

func probe(p *corev1.Probe) {
	if p == nil {
		return
	}
	setIfZero(&p.PeriodSeconds, 10)
	setIfZero(&p.SuccessThreshold, 1)
	setIfZero(&p.FailureThreshold, 3)
	if p.HTTPGet != nil {
		setIfZero(&p.HTTPGet.Scheme, corev1.URISchemeHTTP)
	}
}

func setIfZero[T comparable](field *T, value T) {
	var zero T
	if *field == zero {
		*field = value
	}
}
