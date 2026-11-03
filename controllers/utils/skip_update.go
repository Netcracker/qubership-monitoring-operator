package utils

import (
	"encoding/json"
	"maps"

	apiequality "k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8slabels "k8s.io/apimachinery/pkg/labels"
)

// SetIfChanged assigns src to *dst when they are not equal. Equality uses
// Kubernetes Semantic.DeepEqual, which understands API types (Quantity, Time)
// and treats nil maps/slices as equal to empty ones. Returns true when *dst
// was updated.
func SetIfChanged[T any](dst *T, src T) bool {
	if apiequality.Semantic.DeepEqual(*dst, src) {
		return false
	}
	*dst = src
	return true
}

// SetLabelsIfChanged assigns labels when they differ according to
// labels.Equals (nil and empty maps are equal). Returns true when labels
// were updated.
func SetLabelsIfChanged(obj metav1.Object, labels map[string]string) bool {
	if k8slabels.Equals(obj.GetLabels(), labels) {
		return false
	}
	obj.SetLabels(labels)
	return true
}

// SetAnnotationsIfChanged assigns annotations when they differ. Nil and empty
// maps are treated as equal. Callers that must distinguish omitted versus empty
// should canonicalize before calling. Returns true when annotations were updated.
func SetAnnotationsIfChanged(obj metav1.Object, annotations map[string]string) bool {
	if maps.Equal(obj.GetAnnotations(), annotations) {
		return false
	}
	obj.SetAnnotations(annotations)
	return true
}

// SetUnstructuredFieldIfChanged assigns obj[key] when the values differ.
// Comparison uses Semantic equality, then a JSON round-trip so numbers that the
// apiserver decoded as float64 still match in-memory int/int32 builders.
// Returns true when the field was updated.
func SetUnstructuredFieldIfChanged(obj map[string]interface{}, key string, src interface{}) bool {
	if UnstructuredValuesEqual(obj[key], src) {
		return false
	}
	obj[key] = src
	return true
}

// UnstructuredValuesEqual reports whether a and b are equal, either
// semantically or after a JSON round-trip of both values.
func UnstructuredValuesEqual(a, b interface{}) bool {
	if apiequality.Semantic.DeepEqual(a, b) {
		return true
	}
	ao, errA := jsonRoundTrip(a)
	bo, errB := jsonRoundTrip(b)
	return errA == nil && errB == nil && apiequality.Semantic.DeepEqual(ao, bo)
}

// jsonRoundTrip encodes v as JSON and decodes it into generic values, the
// form in which the API server returns unstructured objects.
func jsonRoundTrip(v interface{}) (interface{}, error) {
	raw, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var out interface{}
	err = json.Unmarshal(raw, &out)
	return out, err
}
