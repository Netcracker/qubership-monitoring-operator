package grafana

import (
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils/labelsassert"
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

func TestGrafanaManifests(t *testing.T) {
	cr = &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
		Spec: monv1.PlatformMonitoringSpec{
			Grafana: &monv1.Grafana{
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
		assert.NotNil(t, m.Spec.Client)
		assert.True(t, m.Spec.Client.UseKubeAuth)
		assert.NotNil(t, m.GetLabels())
		assert.Equal(t, labelValue, m.GetLabels()[labelKey])
		// In grafana-operator v5, Labels and Annotations are in Deployment.Spec.Template
		if m.Spec.Deployment != nil && m.Spec.Deployment.Spec.Template != nil {
			if m.Spec.Deployment.Spec.Template.Labels != nil {
				assert.Equal(t, labelValue, m.Spec.Deployment.Spec.Template.Labels[labelKey])
			}
			if m.Spec.Deployment.Spec.Template.Annotations != nil {
				assert.Equal(t, annotationValue, m.Spec.Deployment.Spec.Template.Annotations[annotationKey])
			}
		}
		assert.NotNil(t, m.GetAnnotations())
		assert.Equal(t, annotationValue, m.GetAnnotations()[annotationKey])
	})
	cr = &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
		Spec: monv1.PlatformMonitoringSpec{
			Grafana: &monv1.Grafana{},
		},
	}
	// Disabled for v5: in v5 labels/annotations live in Deployment.Spec.Template, not Deployment
	//t.Run("Test Grafana manifest with nil annotation", func(t *testing.T) {
	//	m, err := grafana(cr)
	//	...
	//})
	t.Run("Test GrafanaDatasource manifest", func(t *testing.T) {
		m, err := grafanaDataSource(cr, nil, nil, nil)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "GrafanaDatasource manifest should not be empty")
	})
	t.Run("Test GrafanaPromxyDatasource manifest", func(t *testing.T) {
		m, err := grafanaPromxyDataSource(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "GrafanaPromxyDatasource manifest should not be empty")
		assert.Equal(t, "platform-monitoring-promxy", m.GetName())
		if m.Spec.Datasource != nil {
			assert.Contains(t, m.Spec.Datasource.URL, "promxy")
		}
	})
	t.Run("Test Ingress v1 manifest", func(t *testing.T) {
		m, err := grafanaIngressV1(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Ingress v1 manifest should not be empty")
	})
	t.Run("Test PodMonitor manifest", func(t *testing.T) {
		crWithLabels := &monv1.PlatformMonitoring{
			ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring", Labels: map[string]string{labelKey: labelValue}},
			Spec: monv1.PlatformMonitoringSpec{
				Grafana: &monv1.Grafana{Labels: map[string]string{labelKey: labelValue}},
			},
		}
		m, err := grafanaPodMonitor(crWithLabels)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "PodMonitor manifest should not be empty")
		labelsassert.AssertCRLabels(t, m.GetLabels(), utils.GrafanaComponentName, "victoriametrics-operator", map[string]string{labelKey: labelValue})
	})
}
