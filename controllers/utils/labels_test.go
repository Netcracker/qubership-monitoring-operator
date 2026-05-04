package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// objWithLabels implements metav1.Object for tests.
type objWithLabels struct {
	metav1.ObjectMeta
}

func (o *objWithLabels) GetLabels() map[string]string              { return o.Labels }
func (o *objWithLabels) SetLabels(labels map[string]string)          { o.Labels = labels }
func (o *objWithLabels) GetAnnotations() map[string]string           { return o.Annotations }
func (o *objWithLabels) SetAnnotations(annotations map[string]string) { o.Annotations = annotations }

func TestCommonLabels(t *testing.T) {
	got := CommonLabels()
	assert.Equal(t, PartOfMonitoring, got["app.kubernetes.io/part-of"])
	assert.Equal(t, ManagedByOperator, got["app.kubernetes.io/managed-by"])
	assert.Equal(t, OperatorDeploymentName, got["app.kubernetes.io/managed-by-operator"])
	assert.Len(t, got, 3)
}

func TestTruncLabel(t *testing.T) {
	t.Run("short label unchanged", func(t *testing.T) {
		assert.Equal(t, "short", TruncLabel("short"))
	})
	t.Run("exactly 63 chars", func(t *testing.T) {
		s := string(make([]byte, 63))
		for i := range s {
			s = s[:i] + "a" + s[i+1:]
		}
		assert.Len(t, TruncLabel(s), 63)
	})
	t.Run("long label truncated and trimmed", func(t *testing.T) {
		long := ""
		for i := 0; i < 70; i++ {
			long += "x"
		}
		got := TruncLabel(long)
		assert.LessOrEqual(t, len(got), 63)
		assert.Equal(t, 63, len(got))
	})
	t.Run("trailing hyphens trimmed", func(t *testing.T) {
		assert.Equal(t, "foo", TruncLabel("foo---"))
	})
}

func TestGetInstanceLabel(t *testing.T) {
	assert.Equal(t, "name-ns", GetInstanceLabel("name", "ns"))
	assert.LessOrEqual(t, len(GetInstanceLabel("a", "b")), 63)
}

func TestResourceLabels(t *testing.T) {
	got := ResourceLabels("my-name", "my-component")
	assert.Equal(t, "my-name", got["name"])
	assert.Equal(t, "my-name", got["app.kubernetes.io/name"])
	assert.Equal(t, "my-component", got["app.kubernetes.io/component"])
	assert.Equal(t, PartOfMonitoring, got["app.kubernetes.io/part-of"])
	assert.Equal(t, ManagedByOperator, got["app.kubernetes.io/managed-by"])
	assert.Equal(t, OperatorDeploymentName, got["app.kubernetes.io/managed-by-operator"])
	assert.Len(t, got, 6)
}

func TestMergeLabels(t *testing.T) {
	t.Run("single map", func(t *testing.T) {
		m := map[string]string{"a": "1"}
		got := MergeLabels(m)
		assert.Equal(t, map[string]string{"a": "1"}, got)
	})
	t.Run("two maps later overrides", func(t *testing.T) {
		got := MergeLabels(
			map[string]string{"a": "1", "b": "2"},
			map[string]string{"b": "over", "c": "3"},
		)
		assert.Equal(t, map[string]string{"a": "1", "b": "over", "c": "3"}, got)
	})
	t.Run("nil maps skipped", func(t *testing.T) {
		got := MergeLabels(nil, map[string]string{"a": "1"}, nil)
		assert.Equal(t, map[string]string{"a": "1"}, got)
	})
	t.Run("empty", func(t *testing.T) {
		got := MergeLabels()
		assert.NotNil(t, got)
		assert.Empty(t, got)
	})
}

func TestMergeInto(t *testing.T) {
	t.Run("nil src leaves dst unchanged", func(t *testing.T) {
		dst := map[string]string{"a": "1"}
		MergeInto(dst, nil)
		assert.Equal(t, map[string]string{"a": "1"}, dst)
	})
	t.Run("src merged into dst", func(t *testing.T) {
		dst := map[string]string{"a": "1"}
		MergeInto(dst, map[string]string{"b": "2"})
		assert.Equal(t, map[string]string{"a": "1", "b": "2"}, dst)
	})
}

func TestBaseOnlyLabelInput(t *testing.T) {
	in := BaseOnlyLabelInput("app", "comp")
	assert.Equal(t, "app", in.Name)
	assert.Equal(t, "comp", in.Component)
	assert.Empty(t, in.Instance)
	assert.Empty(t, in.Version)
	assert.Empty(t, in.Technology)
	assert.Nil(t, in.ComponentLabels)
}

func TestSetLabelsForResource(t *testing.T) {
	t.Run("existing nil uses empty and applies labels", func(t *testing.T) {
		obj := &objWithLabels{ObjectMeta: metav1.ObjectMeta{Name: "x", Namespace: "ns"}}
		in := BaseOnlyLabelInput("x", "comp")
		SetLabelsForResource(obj, in, nil)
		labels := obj.GetLabels()
		assert.Equal(t, "x", labels["name"])
		assert.Equal(t, "comp", labels["app.kubernetes.io/component"])
		assert.Equal(t, PartOfMonitoring, labels["app.kubernetes.io/part-of"])
	})
	t.Run("existing non-nil is merged", func(t *testing.T) {
		obj := &objWithLabels{
			ObjectMeta: metav1.ObjectMeta{
				Name: "x", Namespace: "ns",
				Labels: map[string]string{"keep": "value"},
			},
		}
		in := BaseOnlyLabelInput("x", "comp")
		SetLabelsForResource(obj, in, nil)
		labels := obj.GetLabels()
		assert.Equal(t, "value", labels["keep"])
		assert.Equal(t, "comp", labels["app.kubernetes.io/component"])
	})
	t.Run("explicit existing map", func(t *testing.T) {
		obj := &objWithLabels{ObjectMeta: metav1.ObjectMeta{Name: "x", Namespace: "ns"}}
		in := BaseOnlyLabelInput("x", "comp")
		SetLabelsForResource(obj, in, map[string]string{"existing": "val"})
		labels := obj.GetLabels()
		assert.Equal(t, "val", labels["existing"])
		assert.Equal(t, "comp", labels["app.kubernetes.io/component"])
	})
	t.Run("with instance version technology and component labels", func(t *testing.T) {
		obj := &objWithLabels{ObjectMeta: metav1.ObjectMeta{Name: "x", Namespace: "ns"}}
		in := LabelInput{
			Name:            "x",
			Component:       "comp",
			Instance:        "inst",
			Version:         "v1",
			Technology:      "go",
			ComponentLabels: map[string]string{"custom": "customVal"},
		}
		SetLabelsForResource(obj, in, nil)
		labels := obj.GetLabels()
		assert.Equal(t, "inst", labels["app.kubernetes.io/instance"])
		assert.Equal(t, "v1", labels["app.kubernetes.io/version"])
		assert.Equal(t, "go", labels["app.kubernetes.io/technology"])
		assert.Equal(t, "customVal", labels["custom"])
	})
}

func TestSetLabelsForWorkload(t *testing.T) {
	podLabels := make(map[string]string)
	obj := &objWithLabels{ObjectMeta: metav1.ObjectMeta{Name: "w", Namespace: "ns"}}
	in := LabelInput{
		Name: "w", Component: "comp",
		Instance: "i", Version: "1", Technology: "go",
	}
	SetLabelsForWorkload(obj, &podLabels, in)
	labels := obj.GetLabels()
	assert.Equal(t, "w", labels["name"])
	assert.Equal(t, "i", labels["app.kubernetes.io/instance"])
	assert.Equal(t, "1", labels["app.kubernetes.io/version"])
	assert.Equal(t, "go", labels["app.kubernetes.io/technology"])
	assert.Equal(t, labels["name"], podLabels["name"])
	assert.Equal(t, labels["app.kubernetes.io/instance"], podLabels["app.kubernetes.io/instance"])
}

func TestLabelsForPodTemplate(t *testing.T) {
	got := LabelsForPodTemplate("n", "c", "i", "v", "tech")
	assert.Equal(t, "n", got["name"])
	assert.Equal(t, "c", got["app.kubernetes.io/component"])
	assert.Equal(t, "i", got["app.kubernetes.io/instance"])
	assert.Equal(t, "v", got["app.kubernetes.io/version"])
	assert.Equal(t, "tech", got["app.kubernetes.io/technology"])
}

func TestLabelInput_Labels(t *testing.T) {
	in := LabelInput{Name: "n", Component: "c", Instance: "i"}
	got := in.Labels(nil)
	assert.Equal(t, "n", got["name"])
	assert.Equal(t, "i", got["app.kubernetes.io/instance"])
}
