package utils

import (
	"slices"

	secv1 "github.com/openshift/api/security/v1"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AnyChanged reports whether any result is true. Go evaluates every argument
// before the call, so each SetIfChanged-style call passed here runs, unlike
// with a short-circuiting || chain.
func AnyChanged(results ...bool) bool {
	return slices.Contains(results, true)
}

// SyncLabels copies the labels of desired onto existing. Use it for objects
// whose only managed field is metadata.labels, such as ServiceAccounts and
// role bindings.
func SyncLabels[T metav1.Object](existing, desired T) bool {
	return SetLabelsIfChanged(existing, desired.GetLabels())
}

// SyncClusterRole copies labels and rules.
func SyncClusterRole(existing, desired *rbacv1.ClusterRole) bool {
	return AnyChanged(
		SetLabelsIfChanged(existing, desired.GetLabels()),
		SetIfChanged(&existing.Rules, desired.Rules),
	)
}

// SyncRole copies labels and rules.
func SyncRole(existing, desired *rbacv1.Role) bool {
	return AnyChanged(
		SetLabelsIfChanged(existing, desired.GetLabels()),
		SetIfChanged(&existing.Rules, desired.Rules),
	)
}

// SyncService copies labels, ports, and the selector. Other spec fields, such as
// clusterIP, are assigned by the API server and stay untouched.
func SyncService(existing, desired *corev1.Service) bool {
	return AnyChanged(
		SetLabelsIfChanged(existing, desired.GetLabels()),
		SetIfChanged(&existing.Spec.Ports, desired.Spec.Ports),
		SetIfChanged(&existing.Spec.Selector, desired.Spec.Selector),
	)
}

// SyncIngress copies labels, annotations, rules, and TLS settings.
func SyncIngress(existing, desired *networkingv1.Ingress) bool {
	return AnyChanged(
		SetLabelsIfChanged(existing, desired.GetLabels()),
		SetAnnotationsIfChanged(existing, desired.GetAnnotations()),
		SetIfChanged(&existing.Spec.Rules, desired.Spec.Rules),
		SetIfChanged(&existing.Spec.TLS, desired.Spec.TLS),
	)
}

// SyncServiceMonitor copies labels, the job label, endpoints, the namespace
// selector, and the selector.
func SyncServiceMonitor(existing, desired *promv1.ServiceMonitor) bool {
	return AnyChanged(
		SetLabelsIfChanged(existing, desired.GetLabels()),
		SetIfChanged(&existing.Spec.JobLabel, desired.Spec.JobLabel),
		SetIfChanged(&existing.Spec.Endpoints, desired.Spec.Endpoints),
		SetIfChanged(&existing.Spec.NamespaceSelector, desired.Spec.NamespaceSelector),
		SetIfChanged(&existing.Spec.Selector, desired.Spec.Selector),
	)
}

// SyncPodMonitor copies labels, the job label, pod metrics endpoints, the
// namespace selector, and the selector.
func SyncPodMonitor(existing, desired *promv1.PodMonitor) bool {
	return AnyChanged(
		SetLabelsIfChanged(existing, desired.GetLabels()),
		SetIfChanged(&existing.Spec.JobLabel, desired.Spec.JobLabel),
		SetIfChanged(&existing.Spec.PodMetricsEndpoints, desired.Spec.PodMetricsEndpoints),
		SetIfChanged(&existing.Spec.NamespaceSelector, desired.Spec.NamespaceSelector),
		SetIfChanged(&existing.Spec.Selector, desired.Spec.Selector),
	)
}

// SyncSecurityContextConstraints copies labels and the policy fields.
func SyncSecurityContextConstraints(existing, desired *secv1.SecurityContextConstraints) bool {
	return AnyChanged(
		SetLabelsIfChanged(existing, desired.GetLabels()),
		ApplySecurityContextConstraintsPolicy(existing, desired),
	)
}

// SyncPodSpec copies the pod spec fields that every managed workload sets: the
// security context, containers, volumes, service account name, affinity, and
// priority class. Fields defaulted by the API server, such as dnsPolicy and
// restartPolicy, stay untouched so idle reconciles do not send Updates. Callers
// sync any additional fields themselves.
func SyncPodSpec(existing, desired *corev1.PodSpec) bool {
	return AnyChanged(
		SetIfChanged(&existing.SecurityContext, desired.SecurityContext),
		SetIfChanged(&existing.Containers, desired.Containers),
		SetVolumesIfChanged(&existing.Volumes, desired.Volumes),
		SetIfChanged(&existing.ServiceAccountName, desired.ServiceAccountName),
		SetIfChanged(&existing.Affinity, desired.Affinity),
		SetIfChanged(&existing.PriorityClassName, desired.PriorityClassName),
	)
}
