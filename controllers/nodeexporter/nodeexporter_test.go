package nodeexporter

import (
	"testing"

	v1beta1 "github.com/Netcracker/qubership-monitoring-operator/api"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	cr              *v1beta1.PlatformMonitoring
	labelKey        = "label.key"
	labelValue      = "label-value"
	annotationKey   = "annotation.key"
	annotationValue = "annotation-value"
)

func TestNodeExporterManifests(t *testing.T) {
	cr = &v1beta1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
		Spec: v1beta1.PlatformMonitoringSpec{
			NodeExporter: &v1beta1.NodeExporter{
				Annotations: map[string]string{annotationKey: annotationValue},
				Labels:      map[string]string{labelKey: labelValue},
			},
		},
	}
	t.Run("Test DaemonSet manifest", func(t *testing.T) {
		m, err := nodeExporterDaemonSet(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "DaemonSet manifest should not be empty")
		assert.NotNil(t, m.GetLabels())
		assert.Equal(t, labelValue, m.GetLabels()[labelKey])
		assert.NotNil(t, m.Spec.Template.Labels)
		assert.Equal(t, labelValue, m.Spec.Template.Labels[labelKey])
		assert.NotNil(t, m.GetAnnotations())
		assert.Equal(t, annotationValue, m.GetAnnotations()[annotationKey])
		assert.NotNil(t, m.Spec.Template.Annotations)
		assert.Equal(t, annotationValue, m.Spec.Template.Annotations[annotationKey])
	})
	cr = &v1beta1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
		Spec: v1beta1.PlatformMonitoringSpec{
			NodeExporter: &v1beta1.NodeExporter{},
		},
	}
	t.Run("Test DaemonSet manifest with nil labels and annotation", func(t *testing.T) {
		m, err := nodeExporterDaemonSet(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "DaemonSet manifest should not be empty")
		assert.NotNil(t, m.GetLabels())
		assert.NotNil(t, m.Spec.Template.Labels)
		assert.Nil(t, m.GetAnnotations())
		assert.Nil(t, m.Spec.Template.Annotations)
	})
	t.Run("Test ClusterRole manifest", func(t *testing.T) {
		m, err := nodeExporterClusterRole(cr, true, false)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "ClusterRole manifest should not be empty")
	})
	t.Run("Test ClusterRoleBinding manifest", func(t *testing.T) {
		m, err := nodeExporterClusterRoleBinding(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "ClusterRoleBinding manifest should not be empty")
	})
	t.Run("Test Service manifest", func(t *testing.T) {
		m, err := nodeExporterService(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Service manifest should not be empty")
	})
	t.Run("Test ServiceMonitor manifest", func(t *testing.T) {
		m, err := nodeExporterServiceMonitor(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "ServiceMonitor manifest should not be empty")
	})
	t.Run("Test ServiceAccount manifest", func(t *testing.T) {
		m, err := nodeExporterServiceAccount(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "ServiceAccount manifest should not be empty")
	})
}
