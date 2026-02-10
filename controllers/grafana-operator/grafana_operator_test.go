package grafana_operator

import (
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
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

func TestGrafanaOperatorManifests(t *testing.T) {
	cr = &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
		Spec: monv1.PlatformMonitoringSpec{
			Grafana: &monv1.Grafana{
				Operator: monv1.GrafanaOperator{
					Annotations: map[string]string{annotationKey: annotationValue},
					Labels:      map[string]string{labelKey: labelValue},
				},
			},
		},
	}
	t.Run("Test Deployment manifest", func(t *testing.T) {
		m, err := grafanaOperatorDeployment(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Deployment manifest should not be empty")
		assert.NotNil(t, m.GetLabels())
		assert.Equal(t, labelValue, m.GetLabels()[labelKey])
		assert.NotNil(t, m.Spec.Template.Labels)
		assert.Equal(t, labelValue, m.Spec.Template.Labels[labelKey])
		assert.NotNil(t, m.GetAnnotations())
		assert.Equal(t, annotationValue, m.GetAnnotations()[annotationKey])
		assert.NotNil(t, m.Spec.Template.Annotations)
		assert.Equal(t, annotationValue, m.Spec.Template.Annotations[annotationKey])
	})
	cr = &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
		Spec: monv1.PlatformMonitoringSpec{
			Grafana: &monv1.Grafana{
				Operator: monv1.GrafanaOperator{},
			},
		},
	}
	t.Run("Test Deployment manifest with nil labels and annotation", func(t *testing.T) {
		m, err := grafanaOperatorDeployment(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Deployment manifest should not be empty")
		assert.NotNil(t, m.GetLabels())
		assert.NotNil(t, m.Spec.Template.Labels)
		assert.Nil(t, m.GetAnnotations())
		assert.Nil(t, m.Spec.Template.Annotations)
	})
	t.Run("Test ServiceAccount manifest", func(t *testing.T) {
		m, err := grafanaOperatorServiceAccount(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "ServiceAccount manifest should not be empty")
	})
	t.Run("Test ClusterRole manifest", func(t *testing.T) {
		m, err := grafanaOperatorClusterRole(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "ClusterRole manifest should not be empty")
	})
	t.Run("Test ClusterRoleBinding manifest", func(t *testing.T) {
		m, err := grafanaOperatorClusterRoleBinding(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "ClusterRoleBinding manifest should not be empty")
	})
	t.Run("Test Role manifest", func(t *testing.T) {
		m, err := grafanaOperatorRole(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Role manifest should not be empty")
	})
	t.Run("Test RoleBinding manifest", func(t *testing.T) {
		m, err := grafanaOperatorRoleBinding(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "RoleBinding manifest should not be empty")
	})
	t.Run("Test GrafanaDashboard manifest", func(t *testing.T) {
		for _, mResource := range utils.GrafanaKubernetesDashboardsResources {
			m, err := grafanaDashboard(cr, mResource)
			if err != nil {
				t.Fatal(err)
			}
			assert.NotNil(t, m, "GrafanaDashboard manifest should not be empty")
		}
	})
}
