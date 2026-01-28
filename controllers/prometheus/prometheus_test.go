package prometheus

import (
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	cr              *monv1.PlatformMonitoring
	labelKey        = "label.key"
	labelValue      = "label-value"
	annotationKey   = "annotation.key"
	annotationValue = "annotation-value"
)

func TestPrometheusManifests(t *testing.T) {
	cr = &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
		Spec: monv1.PlatformMonitoringSpec{
			Prometheus: &monv1.Prometheus{
				Annotations: map[string]string{annotationKey: annotationValue},
				Labels:      map[string]string{labelKey: labelValue},
			},
		},
	}
	t.Run("Test Prometheus manifest", func(t *testing.T) {
		m, err := prometheus(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Prometheus manifest should not be empty")
		assert.NotNil(t, m.Spec.PodMetadata.Labels)
		assert.Equal(t, labelValue, m.Spec.PodMetadata.Labels[labelKey])
		assert.NotNil(t, m.Spec.PodMetadata.Annotations)
		assert.Equal(t, annotationValue, m.Spec.PodMetadata.Annotations[annotationKey])
	})
	cr = &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
		Spec: monv1.PlatformMonitoringSpec{
			Prometheus: &monv1.Prometheus{},
		},
	}
	t.Run("Test Prometheus manifest with nil labels and annotation", func(t *testing.T) {
		m, err := prometheus(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Prometheus manifest should not be empty")
		assert.NotNil(t, m.GetLabels())
		assert.Nil(t, m.GetAnnotations())
	})
	cr = &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
	}
	t.Run("Test ServiceAccount manifest", func(t *testing.T) {
		m, err := prometheusServiceAccount(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "ServiceAccount manifest should not be empty")
	})
	t.Run("Test ClusterRole manifest", func(t *testing.T) {
		m, err := prometheusClusterRole(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "ClusterRole manifest should not be empty")
	})
	t.Run("Test ClusterRoleBinding manifest", func(t *testing.T) {
		m, err := prometheusClusterRoleBinding(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "ClusterRoleBinding manifest should not be empty")
	})
	t.Run("Test Ingress v1beta1 manifest", func(t *testing.T) {
		m, err := prometheusIngressV1beta1(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Ingress v1beta1 manifest should not be empty")
	})
	t.Run("Test Ingress v1 manifest", func(t *testing.T) {
		m, err := prometheusIngressV1(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Ingress v1 manifest should not be empty")
	})
	t.Run("Test PodMonitor manifest", func(t *testing.T) {
		m, err := prometheusPodMonitor(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "PodMonitor manifest should not be empty")
	})
}
