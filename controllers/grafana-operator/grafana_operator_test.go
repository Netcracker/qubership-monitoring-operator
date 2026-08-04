package grafana_operator

import (
	"testing"
	"time"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils/labelsassert"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
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
	t.Run("Test Deployment manifest is namespace scoped without privileged rights", func(t *testing.T) {
		privilegedRights := utils.PrivilegedRights
		utils.PrivilegedRights = false
		t.Cleanup(func() {
			utils.PrivilegedRights = privilegedRights
		})

		m, err := grafanaOperatorDeployment(cr)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, cr.GetNamespace(), containerEnvValue(m.Spec.Template.Spec.Containers, utils.GrafanaOperatorComponentName, "WATCH_NAMESPACE"))
		assert.Empty(t, containerEnvValue(m.Spec.Template.Spec.Containers, utils.GrafanaOperatorComponentName, "WATCH_NAMESPACE_SELECTOR"))
	})
	t.Run("Test ServiceAccount manifest", func(t *testing.T) {
		crWithSALabels := &monv1.PlatformMonitoring{
			ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"},
			Spec: monv1.PlatformMonitoringSpec{
				Grafana: &monv1.Grafana{
					Operator: monv1.GrafanaOperator{
						ServiceAccount: &monv1.EmbeddedObjectMetadata{
							Labels: map[string]string{labelKey: labelValue},
						},
					},
				},
			},
		}
		m, err := grafanaOperatorServiceAccount(crWithSALabels)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "ServiceAccount manifest should not be empty")
		assert.NotNil(t, m.GetLabels())
		assert.Equal(t, labelValue, m.GetLabels()[labelKey], "ServiceAccount.Labels should be merged")
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
		assert.True(t, roleContainsRule(m, "", "events", []string{"create", "patch"}))
		assert.True(t, roleContainsRule(m, "events.k8s.io", "events", []string{"create", "patch"}))
		assert.True(t, roleContainsRule(m, "gateway.networking.k8s.io", "httproutes", []string{"get", "list", "watch", "create", "update", "patch", "delete"}))
		assert.True(t, roleContainsRule(m, "grafana.integreatly.org", "grafanamanifests", []string{"get", "list", "watch", "create", "update", "patch", "delete"}))
		assert.True(t, roleContainsRule(m, "grafana.integreatly.org", "grafanamanifests/status", []string{"get", "list", "watch", "create", "update", "patch", "delete"}))
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
			assert.Equal(t, 10*time.Minute, m.Spec.ResyncPeriod.Duration)
		}
	})
	crWithLabels := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring", Labels: map[string]string{labelKey: labelValue}},
		Spec: monv1.PlatformMonitoringSpec{
			Grafana: &monv1.Grafana{Operator: monv1.GrafanaOperator{Labels: map[string]string{labelKey: labelValue}}},
		},
	}
	t.Run("Test PodMonitor manifest", func(t *testing.T) {
		m, err := grafanaOperatorPodMonitor(crWithLabels)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "PodMonitor manifest should not be empty")
		labelsassert.AssertCRLabels(t, m.GetLabels(), utils.GrafanaOperatorComponentName, "victoriametrics-operator", map[string]string{labelKey: labelValue})
	})
}

func containerEnvValue(containers []corev1.Container, containerName, envName string) string {
	for _, container := range containers {
		if container.Name != containerName {
			continue
		}
		for _, env := range container.Env {
			if env.Name == envName {
				return env.Value
			}
		}
	}
	return ""
}

func roleContainsRule(role *rbacv1.Role, apiGroup, resource string, verbs []string) bool {
	for _, rule := range role.Rules {
		if !assert.ObjectsAreEqual([]string{apiGroup}, rule.APIGroups) {
			continue
		}
		for _, candidate := range rule.Resources {
			if candidate == resource {
				return assert.ObjectsAreEqualValues(verbs, rule.Verbs)
			}
		}
	}
	return false
}
