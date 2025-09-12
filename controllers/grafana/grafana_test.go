package grafana

import (
	"testing"

	v1alpha1 "github.com/Netcracker/qubership-monitoring-operator/api/v1alpha1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	cr              *v1alpha1.PlatformMonitoring
	labelKey        = "label.key"
	labelValue      = "label-value"
	annotationKey   = "annotation.key"
	annotationValue = "annotation-value"
)

func TestGrafanaManifests(t *testing.T) {
	cr = &v1alpha1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
		Spec: v1alpha1.PlatformMonitoringSpec{
			Grafana: &v1alpha1.Grafana{
				Annotations: map[string]string{annotationKey: annotationValue},
				Labels:      map[string]string{labelKey: labelValue},
			},
		},
	}
	t.Run("Test Grafana manifest", func(t *testing.T) {
		m, err := grafana(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Grafana manifest should not be empty")
		assert.NotNil(t, m.GetLabels())
		assert.Equal(t, labelValue, m.GetLabels()[labelKey])
		assert.NotNil(t, m.Spec.Deployment.Labels)
		assert.Equal(t, labelValue, m.Spec.Deployment.Labels[labelKey])
		assert.NotNil(t, m.GetAnnotations())
		assert.Equal(t, annotationValue, m.GetAnnotations()[annotationKey])
		assert.NotNil(t, m.Spec.Deployment.Annotations)
		assert.Equal(t, annotationValue, m.Spec.Deployment.Annotations[annotationKey])
	})
	cr = &v1alpha1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
		Spec: v1alpha1.PlatformMonitoringSpec{
			Grafana: &v1alpha1.Grafana{},
		},
	}
	//t.Run("Test Grafana manifest with nil annotation", func(t *testing.T) {
	//	m, err := grafana(cr)
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//	assert.NotNil(t, m, "Grafana manifest should not be empty")
	//	assert.NotNil(t, m.GetLabels())
	//	assert.Nil(t, m.Spec.Deployment.Labels)
	//	assert.Nil(t, m.GetAnnotations())
	//	assert.Nil(t, m.Spec.Deployment.Annotations)
	//})
	t.Run("Test GrafanaDataSource manifest", func(t *testing.T) {
		m, err := grafanaDataSource(cr, nil, nil, nil)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "GrafanaDataSource manifest should not be empty")
	})
	t.Run("Test Ingress v1beta1 manifest", func(t *testing.T) {
		m, err := grafanaIngressV1beta1(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Ingress v1beta1 manifest should not be empty")
	})
	t.Run("Test Ingress v1 manifest", func(t *testing.T) {
		m, err := grafanaIngressV1(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Ingress v1 manifest should not be empty")
	})
	t.Run("Test PodMonitor manifest", func(t *testing.T) {
		m, err := grafanaPodMonitor(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "PodMonitor manifest should not be empty")
	})
}
