package utils

import (
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// PartOfMonitoring is the value for app.kubernetes.io/part-of on all monitoring-operator resources.
	PartOfMonitoring = "monitoring"
	// ManagedByOperator is the value for app.kubernetes.io/managed-by on resources created by the operator.
	ManagedByOperator = "monitoring-operator"
	// OperatorDeploymentName is the name of the monitoring-operator Deployment; used for app.kubernetes.io/managed-by-operator.
	OperatorDeploymentName = "monitoring-operator"
	// ProcessedByOperatorKey is the label key for the operator that processes the resource (e.g. ServiceMonitor, PodMonitor).
	ProcessedByOperatorKey = "app.kubernetes.io/processed-by-operator"
)

// CommonLabels returns the labels applied to all resources (part-of, managed-by, managed-by-operator).
// Single source of truth for operator-created resources; mirrors Helm commonLabels.
func CommonLabels() map[string]string {
	return map[string]string{
		"app.kubernetes.io/part-of":             PartOfMonitoring,
		"app.kubernetes.io/managed-by":          ManagedByOperator,
		"app.kubernetes.io/managed-by-operator": OperatorDeploymentName,
	}
}

// TruncLabel truncates label values to 63 chars (Kubernetes limit) and trims trailing hyphens.
func TruncLabel(label string) string {
	if len(label) >= 63 {
		return strings.Trim(label[:63], "-")
	}
	return strings.Trim(label, "-")
}

// GetInstanceLabel returns the instance label (name-namespace), truncated to 63 chars.
func GetInstanceLabel(name, namespace string) string {
	return TruncLabel(fmt.Sprintf("%s-%s", name, namespace))
}

// ResourceLabels returns name, app.kubernetes.io/name, component, plus CommonLabels.
// Names are truncated to fit the 63-char Kubernetes label limit.
// Use for any resource (Service, ConfigMap, ServiceAccount, workload, etc.) so labels stay consistent.
func ResourceLabels(name, component string) map[string]string {
	return MergeLabels(
		map[string]string{
			"name":                        TruncLabel(name),
			"app.kubernetes.io/name":      TruncLabel(name),
			"app.kubernetes.io/component": component,
		},
		CommonLabels(),
	)
}

// MergeLabels returns a new map with all key-value pairs from the given maps.
// Later maps override earlier ones on key conflict. Nil maps are skipped.
func MergeLabels(maps ...map[string]string) map[string]string {
	out := make(map[string]string)
	for _, m := range maps {
		if m == nil {
			continue
		}
		for k, v := range m {
			out[k] = v
		}
	}
	return out
}

// MergeInto copies all key-value pairs from src into dst. If src is nil, dst is unchanged.
// Caller must ensure dst is non-nil (e.g. after SetLabels or when building a map).
func MergeInto(dst, src map[string]string) {
	if src == nil {
		return
	}
	for k, v := range src {
		dst[k] = v
	}
}

// LabelInput holds parameters for setting labels on operator-managed resources.
// ComponentLabels are applied to both resource metadata and pod templates.
type LabelInput struct {
	Name            string
	Component       string
	Instance        string
	Version         string
	Technology      string
	ComponentLabels map[string]string
}

// BaseOnlyLabelInput returns LabelInput with base labels only (no instance, version, technology).
// Use for ServiceAccount, ClusterRole, Service, ServiceMonitor, etc. per label specification.
func BaseOnlyLabelInput(name, component string) LabelInput {
	return LabelInput{Name: name, Component: component}
}

func (in LabelInput) instanceVersionTechnologyMap() map[string]string {
	m := make(map[string]string)
	if in.Instance != "" {
		m["app.kubernetes.io/instance"] = in.Instance
	}
	if in.Version != "" {
		m["app.kubernetes.io/version"] = in.Version
	}
	if in.Technology != "" {
		m["app.kubernetes.io/technology"] = in.Technology
	}
	return m
}

// resourceLabelsFromBase returns the full resource label map: base maps merged, then component labels merged in.
func (in LabelInput) resourceLabelsFromBase(baseMaps ...map[string]string) map[string]string {
	out := MergeLabels(baseMaps...)
	MergeInto(out, in.ComponentLabels)
	return out
}

// resourceLabels returns the full resource label map. Merge order: existing, ResourceLabels, instance/version/technology, ComponentLabels.
// ComponentLabels override earlier layers on key conflict.
func (in LabelInput) resourceLabels(existing map[string]string) map[string]string {
	return in.resourceLabelsFromBase(existing, ResourceLabels(in.Name, in.Component), in.instanceVersionTechnologyMap())
}

// Labels returns the full label map for resource metadata or pod template (same in default case).
// Use for both SetLabelsForResource and pod template labels. Pass nil for existing when building from scratch.
func (in LabelInput) Labels(existing map[string]string) map[string]string {
	return in.resourceLabels(existing)
}

// SetLabelsForResource sets base + component labels on any resource (Service, ServiceAccount, ConfigMap, etc.).
// When existing is nil, the object's current labels (obj.GetLabels()) are used as the initial layer; when non-nil, existing is used instead.
func SetLabelsForResource(obj metav1.Object, in LabelInput, existing map[string]string) {
	initial := existing
	if initial == nil {
		initial = obj.GetLabels()
	}
	obj.SetLabels(in.resourceLabels(initial))
}

// SetLabelsForWorkload sets labels on the resource and pod template using the same procedure (resource metadata = pod template labels).
// Use for DaemonSet, StatefulSet, Deployment — pass the object and &obj.Spec.Template.Labels.
func SetLabelsForWorkload(obj metav1.Object, templateLabels *map[string]string, in LabelInput) {
	SetLabelsForResource(obj, in, nil)
	*templateLabels = in.Labels(*templateLabels)
}

// LabelsForPodTemplate returns labels for pod template using the same procedure as resource metadata.
// Use when building LabelInput from (name, component, instance, version, technology) without ComponentLabels.
func LabelsForPodTemplate(name, component, instance, version, technology string) map[string]string {
	in := LabelInput{Name: name, Component: component, Instance: instance, Version: version, Technology: technology}
	return in.Labels(nil)
}
