// NOTE: Boilerplate only.  Ignore this file.

// Package v1 contains API Schema definitions for the monitoring v1 API group
// +k8s:deepcopy-gen=package,register
// +groupName=monitoring.netcracker.com
package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	// SchemeGroupVersion is group version used to register these objects
	SchemeGroupVersion = schema.GroupVersion{Group: "monitoring.netcracker.com", Version: "v1"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)

	// AddToScheme adds the types in this group-version to the given scheme.
	AddToScheme = SchemeBuilder.AddToScheme
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion, &PlatformMonitoring{}, &PlatformMonitoringList{})
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
