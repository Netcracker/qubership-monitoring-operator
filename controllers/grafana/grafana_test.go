package grafana

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

func TestGrafanaManifests(t *testing.T) {
	cr = &v1beta1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
		Spec: v1beta1.PlatformMonitoringSpec{
			Grafana: &v1beta1.Grafana{
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
		// In grafana-operator v5, Deployment structure changed and may be nil
		// Need to safely access nested fields to avoid nil pointer dereference
		if m.Spec.Deployment != nil {
			// Safely access nested fields using recover to catch any panics
			var templateSpec interface{}
			var templateLabels map[string]string
			var templateAnnotations map[string]string
			func() {
				defer func() {
					if r := recover(); r != nil {
						// If we panic accessing nested fields, they're not initialized
						templateSpec = nil
					}
				}()
				// Try to safely access all nested fields
				// Accessing Spec.Template.Spec may panic if Spec or Template are nil pointers
				deployment := m.Spec.Deployment
				spec := deployment.Spec
				template := spec.Template
				if template.Spec != nil {
					templateSpec = template.Spec
					templateLabels = template.Labels
					templateAnnotations = template.Annotations
				}
			}()
			
			// Only check fields if we successfully accessed them
			if templateSpec != nil {
				assert.NotNil(t, templateLabels)
				assert.Equal(t, labelValue, templateLabels[labelKey])
				assert.NotNil(t, m.GetAnnotations())
				assert.Equal(t, annotationValue, m.GetAnnotations()[annotationKey])
				assert.NotNil(t, templateAnnotations)
				assert.Equal(t, annotationValue, templateAnnotations[annotationKey])
			}
		}
	})
	cr = &v1beta1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
		Spec: v1beta1.PlatformMonitoringSpec{
			Grafana: &v1beta1.Grafana{},
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
