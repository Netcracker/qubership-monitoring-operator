package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestSetIfChangedNoOp(t *testing.T) {
	rules := []rbacv1.PolicyRule{{
		APIGroups: []string{""},
		Resources: []string{"pods"},
		Verbs:     []string{"get"},
	}}
	dst := append([]rbacv1.PolicyRule(nil), rules...)

	assert.False(t, SetIfChanged(&dst, rules))
	assert.Equal(t, rules, dst)
}

func TestSetIfChangedSpecOnly(t *testing.T) {
	dst := corev1.ServiceSpec{
		Ports: []corev1.ServicePort{{Name: "http", Port: 80}},
	}
	src := corev1.ServiceSpec{
		Ports: []corev1.ServicePort{{Name: "http", Port: 443}},
	}

	assert.True(t, SetIfChanged(&dst, src))
	assert.Equal(t, int32(443), dst.Ports[0].Port)
}

func TestSetIfChangedTreatsNilAndEmptySliceAsEqual(t *testing.T) {
	var dst []string
	src := []string{}

	assert.False(t, SetIfChanged(&dst, src), "k8s Semantic.DeepEqual treats nil and empty slices as equal")
	assert.Nil(t, dst)
}

func TestSetIfChangedTreatsNilAndEmptyMapAsEqual(t *testing.T) {
	var dst map[string]string
	src := map[string]string{}

	assert.False(t, SetIfChanged(&dst, src), "k8s Semantic.DeepEqual treats nil and empty maps as equal")
	assert.Nil(t, dst)
}

func TestSetIfChangedDoesNotEquateEmptyAndDefaultProtocol(t *testing.T) {
	dst := []corev1.ServicePort{{Name: "webhook", Port: 443}}
	src := []corev1.ServicePort{{Name: "webhook", Port: 443, Protocol: corev1.ProtocolTCP}}

	assert.True(t, SetIfChanged(&dst, src), "Semantic.DeepEqual does not treat omitted protocol as TCP")
}

func TestSetLabelsIfChangedNoOp(t *testing.T) {
	obj := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"app": "x"}},
	}

	assert.False(t, SetLabelsIfChanged(obj, map[string]string{"app": "x"}))
	assert.Equal(t, map[string]string{"app": "x"}, obj.GetLabels())
}

func TestSetLabelsIfChangedNilEqualsEmpty(t *testing.T) {
	obj := &corev1.Service{}

	assert.False(t, SetLabelsIfChanged(obj, map[string]string{}), "labels.Equals treats nil and empty maps as equal")
	assert.Nil(t, obj.GetLabels())
}

func TestSetLabelsIfChangedUpdates(t *testing.T) {
	obj := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"app": "old"}},
	}

	assert.True(t, SetLabelsIfChanged(obj, map[string]string{"app": "new"}))
	assert.Equal(t, map[string]string{"app": "new"}, obj.GetLabels())
}

func TestSetAnnotationsIfChangedNilEqualsEmpty(t *testing.T) {
	obj := &corev1.Service{}

	assert.False(t, SetAnnotationsIfChanged(obj, map[string]string{}))
	assert.Nil(t, obj.GetAnnotations())
}

func TestSetAnnotationsIfChangedUpdates(t *testing.T) {
	obj := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{Annotations: map[string]string{"k": "old"}},
	}

	assert.True(t, SetAnnotationsIfChanged(obj, map[string]string{"k": "new"}))
	assert.Equal(t, map[string]string{"k": "new"}, obj.GetAnnotations())
}

func TestSetAnnotationsIfChangedRemovesKeys(t *testing.T) {
	obj := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{Annotations: map[string]string{"keep": "v", "drop": "x"}},
	}

	assert.True(t, SetAnnotationsIfChanged(obj, map[string]string{"keep": "v"}))
	assert.Equal(t, map[string]string{"keep": "v"}, obj.GetAnnotations())
}

func TestSetUnstructuredFieldIfChangedNoOp(t *testing.T) {
	obj := map[string]interface{}{
		"spec": map[string]interface{}{"hostnames": []interface{}{"a.example"}},
	}
	src := map[string]interface{}{"hostnames": []interface{}{"a.example"}}

	assert.False(t, SetUnstructuredFieldIfChanged(obj, "spec", src))
}

func TestSetUnstructuredFieldIfChangedUpdates(t *testing.T) {
	obj := map[string]interface{}{
		"spec": map[string]interface{}{"hostnames": []interface{}{"a.example"}},
	}
	src := map[string]interface{}{"hostnames": []interface{}{"b.example"}}

	assert.True(t, SetUnstructuredFieldIfChanged(obj, "spec", src))
	assert.Equal(t, src, obj["spec"])
}

func TestSetUnstructuredFieldIfChangedJSONNumberCanonical(t *testing.T) {
	obj := map[string]interface{}{
		"spec": map[string]interface{}{"port": float64(3000)},
	}
	src := map[string]interface{}{"port": int32(3000)}

	assert.False(t, SetUnstructuredFieldIfChanged(obj, "spec", src))
}
