package vmoperator

import (
	"strings"
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	routev1 "github.com/openshift/api/route/v1"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	fakediscovery "k8s.io/client-go/discovery/fake"
	ktesting "k8s.io/client-go/testing"
	"k8s.io/utils/ptr"
)

var (
	cr              *monv1.PlatformMonitoring
	labelKey        = "label.key"
	labelValue      = "label-value"
	annotationKey   = "annotation.key"
	annotationValue = "annotation-value"
)

func TestVmOperatorManifests(t *testing.T) {
	cr = &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
		Spec: monv1.PlatformMonitoringSpec{
			Victoriametrics: &monv1.Victoriametrics{
				VmOperator: monv1.VmOperator{
					Annotations: map[string]string{annotationKey: annotationValue},
					Labels:      map[string]string{labelKey: labelValue},
					Replicas:    ptr.To[int32](1),
				},
			},
		},
	}
	t.Run("Test Deployment manifest", func(t *testing.T) {
		m, err := vmOperatorDeployment(nil, cr)
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
			Victoriametrics: &monv1.Victoriametrics{
				VmOperator: monv1.VmOperator{},
			},
		},
	}
	t.Run("Test Deployment manifest with nil labels and annotation", func(t *testing.T) {
		m, err := vmOperatorDeployment(nil, cr)
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

		m, err := vmOperatorDeployment(nil, cr)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, cr.GetNamespace(), vmOperatorContainerEnvValue(m.Spec.Template.Spec.Containers, "WATCH_NAMESPACE"))
	})
	t.Run("Test non-privileged namespace scope overrides extra environment variables", func(t *testing.T) {
		privilegedRights := utils.PrivilegedRights
		utils.PrivilegedRights = false
		t.Cleanup(func() {
			utils.PrivilegedRights = privilegedRights
		})

		scopedCR := cr.DeepCopy()
		scopedCR.Spec.Victoriametrics.VmOperator.ExtraEnvs = []corev1.EnvVar{
			{Name: "WATCH_NAMESPACE", Value: "another-namespace"},
		}
		m, err := vmOperatorDeployment(nil, scopedCR)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, scopedCR.GetNamespace(), vmOperatorContainerEnvValue(m.Spec.Template.Spec.Containers, "WATCH_NAMESPACE"))
	})
	t.Run("Test Role manifest", func(t *testing.T) {
		m, err := vmOperatorRole(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Role manifest should not be empty")
		assert.True(t, roleContainsRule(m, "autoscaling.k8s.io", "verticalpodautoscalers", []string{"create", "delete", "get", "list", "update", "watch"}))
		assert.True(t, roleContainsRule(m, "gateway.networking.k8s.io", "httproutes", []string{"create", "delete", "get", "list", "update", "watch"}))
		assert.True(t, roleContainsResource(m, "operator.victoriametrics.com", "vmanomalyconfigs"))
		assert.True(t, roleContainsResource(m, "operator.victoriametrics.com", "vmanomalyconfigs/finalizers"))
		assert.True(t, roleContainsResource(m, "operator.victoriametrics.com", "vmanomalyconfigs/status"))
	})
	t.Run("Test RoleBinding manifest", func(t *testing.T) {
		m, err := vmOperatorRoleBinding(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "RoleBinding manifest should not be empty")
	})
	t.Run("Test ServiceAccount manifest", func(t *testing.T) {
		m, err := vmOperatorServiceAccount(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "ServiceAccount manifest should not be empty")
	})
	t.Run("Test ClusterRole manifest", func(t *testing.T) {
		m, err := vmOperatorClusterRole(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "ClusterRole manifest should not be empty")
		assert.True(t, clusterRoleContains(m, "vmanomalyconfigs"), "ClusterRole should manage VMAnomalyConfig resources")
		assert.True(t, clusterRoleContains(m, "vmanomalyconfigs/finalizers"), "ClusterRole should manage VMAnomalyConfig finalizers")
		assert.True(t, clusterRoleContains(m, "vmanomalyconfigs/status"), "ClusterRole should manage VMAnomalyConfig status")
	})
	t.Run("Test ClusterRoleBinding manifest", func(t *testing.T) {
		m, err := vmOperatorClusterRoleBinding(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "ClusterRoleBinding manifest should not be empty")
	})
	t.Run("Test Service manifest", func(t *testing.T) {
		m, err := vmOperatorService(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Service manifest should not be empty")
	})
	t.Run("Test Kubelet service manifest", func(t *testing.T) {
		m, err := vmKubeletService(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Kubelet service manifest should not be empty")
	})
	t.Run("Test Kubelet service endpoints manifest", func(t *testing.T) {
		m, err := vmKubeletServiceEndpoints(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Kubelet service manifest should not be empty")
	})
	t.Run("Test KubeScheduler service manifest", func(t *testing.T) {
		m, err := vmKubeSchedulerService(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "KubeScheduler service manifest should not be empty")
	})
	t.Run("Test KubeScheduler service endpoints manifest", func(t *testing.T) {
		m, err := vmKubeSchedulerServiceEndpoints(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "KubeScheduler service manifest should not be empty")
	})
	t.Run("Test KubeControllerManager service manifest", func(t *testing.T) {
		m, err := vmKubeControllerManagerService(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "KubeScheduler service manifest should not be empty")
	})
	t.Run("Test KubeControllerManager service endpoints manifest", func(t *testing.T) {
		m, err := vmKubeControllerManagerServiceEndpoints(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "KubeControllerManager service manifest should not be empty")
	})
	// t.Run("Test PodMonitor manifest", func(t *testing.T) {
	// 	m, err := vmOperatorPodMonitor(cr)
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}
	// 	assert.NotNil(t, m, "PodMonitor manifest should not be empty")
	// })
}

func TestVmOperatorLeaderElect(t *testing.T) {
	base := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"},
		Spec: monv1.PlatformMonitoringSpec{
			Victoriametrics: &monv1.Victoriametrics{VmOperator: monv1.VmOperator{}},
		},
	}

	t.Run("omits flag when unset", func(t *testing.T) {
		m, err := vmOperatorDeployment(nil, base)
		if err != nil {
			t.Fatal(err)
		}
		assert.Empty(t, leaderElectFlags(vmOperatorContainerArgs(m.Spec.Template.Spec.Containers)))
	})
	t.Run("passes --leader-elect=true", func(t *testing.T) {
		cr := base.DeepCopy()
		cr.Spec.Victoriametrics.VmOperator.LeaderElect = ptr.To(true)
		m, err := vmOperatorDeployment(nil, cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, []string{"--leader-elect=true"}, leaderElectFlags(vmOperatorContainerArgs(m.Spec.Template.Spec.Containers)))
	})
	t.Run("passes --leader-elect=false", func(t *testing.T) {
		cr := base.DeepCopy()
		cr.Spec.Victoriametrics.VmOperator.LeaderElect = ptr.To(false)
		m, err := vmOperatorDeployment(nil, cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, []string{"--leader-elect=false"}, leaderElectFlags(vmOperatorContainerArgs(m.Spec.Template.Spec.Containers)))
	})
	t.Run("keeps disableCRDOwnership on Route API independently of leaderElect", func(t *testing.T) {
		r := routeAPIReconciler()
		cr := base.DeepCopy()
		cr.Spec.Victoriametrics.VmOperator.LeaderElect = ptr.To(true)
		m, err := vmOperatorDeployment(r, cr)
		if err != nil {
			t.Fatal(err)
		}
		args := vmOperatorContainerArgs(m.Spec.Template.Spec.Containers)
		assert.Contains(t, args, "--controller.disableCRDOwnership=true")
		assert.Equal(t, []string{"--leader-elect=true"}, leaderElectFlags(args))
	})
}

func clusterRoleContains(role *rbacv1.ClusterRole, resource string) bool {
	for _, rule := range role.Rules {
		for _, candidate := range rule.Resources {
			if candidate == resource {
				return true
			}
		}
	}
	return false
}

func routeAPIReconciler() *VmOperatorReconciler {
	return &VmOperatorReconciler{
		ComponentReconciler: &utils.ComponentReconciler{
			Dc: &fakediscovery.FakeDiscovery{Fake: &ktesting.Fake{Resources: []*metav1.APIResourceList{{
				GroupVersion: routev1.GroupVersion.String(),
				APIResources: []metav1.APIResource{{Name: "routes", Kind: "Route"}},
			}}}},
			Log: utils.Logger("vmoperator_test"),
		},
	}
}

func vmOperatorContainerArgs(containers []corev1.Container) []string {
	for _, container := range containers {
		if container.Name == utils.VmOperatorComponentName {
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

func vmOperatorContainerEnvValue(containers []corev1.Container, envName string) string {
	for _, container := range containers {
		if container.Name != utils.VmOperatorComponentName {
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

func roleContainsResource(role *rbacv1.Role, apiGroup, resource string) bool {
	for _, rule := range role.Rules {
		if !assert.ObjectsAreEqual([]string{apiGroup}, rule.APIGroups) {
			continue
		}
		for _, candidate := range rule.Resources {
			if candidate == resource {
				return true
			}
		}
	}
	return false
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
