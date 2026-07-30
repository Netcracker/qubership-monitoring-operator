package grafana_operator

import (
	"strings"
	"testing"
	"time"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils/labelsassert"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
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
		m, err := grafanaOperatorDeployment(cr, false)
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
		assertGrafanaOperatorHardening(t, m, false)
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
		m, err := grafanaOperatorDeployment(cr, false)
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

		m, err := grafanaOperatorDeployment(cr, false)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, cr.GetNamespace(), containerEnvValue(m.Spec.Template.Spec.Containers, utils.GrafanaOperatorComponentName, "WATCH_NAMESPACE"))
		assert.Empty(t, containerEnvValue(m.Spec.Template.Spec.Containers, utils.GrafanaOperatorComponentName, "WATCH_NAMESPACE_SELECTOR"))
	})
	t.Run("Test OpenShift Deployment security context", func(t *testing.T) {
		m, err := grafanaOperatorDeployment(cr, true)
		require.NoError(t, err)
		assertGrafanaOperatorHardening(t, m, true)
	})
	t.Run("Test configured IDs are preserved", func(t *testing.T) {
		configuredCR := cr.DeepCopy()
		configuredCR.Spec.Grafana.Operator.SecurityContext = &monv1.SecurityContext{
			RunAsUser:  ptr.To(int64(3000)),
			RunAsGroup: ptr.To(int64(3000)),
			FSGroup:    ptr.To(int64(3000)),
		}

		m, err := grafanaOperatorDeployment(configuredCR, false)
		require.NoError(t, err)
		require.NotNil(t, m.Spec.Template.Spec.SecurityContext)
		assert.Equal(t, int64(3000), *m.Spec.Template.Spec.SecurityContext.RunAsUser)
		assert.Equal(t, int64(3000), *m.Spec.Template.Spec.SecurityContext.RunAsGroup)
		assert.Equal(t, int64(3000), *m.Spec.Template.Spec.SecurityContext.FSGroup)
		assertGrafanaOperatorHardening(t, m, false)
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

func TestGrafanaOperatorLeaderElect(t *testing.T) {
	base := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"},
		Spec: monv1.PlatformMonitoringSpec{
			Grafana: &monv1.Grafana{Operator: monv1.GrafanaOperator{}},
		},
	}

	t.Run("omits flag when unset", func(t *testing.T) {
		m, err := grafanaOperatorDeployment(base)
		if err != nil {
			t.Fatal(err)
		}
		assert.Empty(t, leaderElectFlags(containerArgs(m.Spec.Template.Spec.Containers, utils.GrafanaOperatorComponentName)))
	})
	t.Run("passes --leader-elect=true", func(t *testing.T) {
		cr := base.DeepCopy()
		cr.Spec.Grafana.Operator.LeaderElect = ptr.To(true)
		m, err := grafanaOperatorDeployment(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, []string{"--leader-elect=true"}, leaderElectFlags(containerArgs(m.Spec.Template.Spec.Containers, utils.GrafanaOperatorComponentName)))
	})
	t.Run("passes --leader-elect=false", func(t *testing.T) {
		cr := base.DeepCopy()
		cr.Spec.Grafana.Operator.LeaderElect = ptr.To(false)
		m, err := grafanaOperatorDeployment(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, []string{"--leader-elect=false"}, leaderElectFlags(containerArgs(m.Spec.Template.Spec.Containers, utils.GrafanaOperatorComponentName)))
	})
}

func containerArgs(containers []corev1.Container, containerName string) []string {
	for _, container := range containers {
		if container.Name == containerName {
			return container.Args
		}
	}
	return nil
}

func leaderElectFlags(args []string) []string {
	var out []string
	for _, arg := range args {
		if arg == "--leader-elect" || strings.HasPrefix(arg, "--leader-elect=") {
			out = append(out, arg)
		}
	}
	return out
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

func assertGrafanaOperatorHardening(t *testing.T, deployment *appsv1.Deployment, isOpenShift bool) {
	t.Helper()

	podSecurityContext := deployment.Spec.Template.Spec.SecurityContext
	require.NotNil(t, podSecurityContext)
	require.NotNil(t, podSecurityContext.RunAsNonRoot)
	assert.True(t, *podSecurityContext.RunAsNonRoot)
	require.NotNil(t, podSecurityContext.SeccompProfile)
	assert.Equal(t, corev1.SeccompProfileTypeRuntimeDefault, podSecurityContext.SeccompProfile.Type)
	if isOpenShift {
		assert.Nil(t, podSecurityContext.RunAsUser)
		assert.Nil(t, podSecurityContext.RunAsGroup)
		assert.Nil(t, podSecurityContext.FSGroup)
	} else {
		require.NotNil(t, podSecurityContext.RunAsUser)
		require.NotNil(t, podSecurityContext.RunAsGroup)
		require.NotNil(t, podSecurityContext.FSGroup)
	}

	require.Len(t, deployment.Spec.Template.Spec.Containers, 1)
	container := deployment.Spec.Template.Spec.Containers[0]
	require.NotNil(t, container.SecurityContext)
	require.NotNil(t, container.SecurityContext.AllowPrivilegeEscalation)
	assert.False(t, *container.SecurityContext.AllowPrivilegeEscalation)
	require.NotNil(t, container.SecurityContext.ReadOnlyRootFilesystem)
	assert.True(t, *container.SecurityContext.ReadOnlyRootFilesystem)
	require.NotNil(t, container.SecurityContext.Capabilities)
	assert.Equal(t, []corev1.Capability{"ALL"}, container.SecurityContext.Capabilities.Drop)
	require.Len(t, container.VolumeMounts, 1)
	assert.Equal(t, "tmp", container.VolumeMounts[0].Name)
	assert.Equal(t, "/tmp", container.VolumeMounts[0].MountPath)

	require.Len(t, deployment.Spec.Template.Spec.Volumes, 1)
	volume := deployment.Spec.Template.Spec.Volumes[0]
	assert.Equal(t, "tmp", volume.Name)
	require.NotNil(t, volume.EmptyDir)
	require.NotNil(t, volume.EmptyDir.SizeLimit)
	assert.Equal(t, "16Mi", volume.EmptyDir.SizeLimit.String())
}
