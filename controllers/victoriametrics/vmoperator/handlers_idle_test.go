package vmoperator

import (
	"testing"

	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils/reconciletest"
	secv1 "github.com/openshift/api/security/v1"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	fakediscovery "k8s.io/client-go/discovery/fake"
	kubernetesfake "k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

func TestHandlersSkipIdleUpdates(t *testing.T) {
	scheme := reconciletest.NewScheme(t, promv1.AddToScheme, secv1.Install)
	counter := &reconciletest.UpdateCounter{}
	node := &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{Name: "node", UID: "node-uid"},
		Status: corev1.NodeStatus{
			Addresses: []corev1.NodeAddress{{Type: corev1.NodeInternalIP, Address: "172.18.0.2"}},
		},
	}
	r := &VmOperatorReconciler{
		KubeClient: kubernetesfake.NewSimpleClientset(node),
		ComponentReconciler: &utils.ComponentReconciler{
			Client: reconciletest.NewClient(scheme, counter),
			Scheme: scheme,
			Dc:     &fakediscovery.FakeDiscovery{Fake: &k8stesting.Fake{}},
			Log:    utils.Logger("vmoperator_idle_test"),
		},
	}

	reconciletest.AssertIdleReconcileSkipsUpdates(t, skipTestPlatformMonitoring(), counter, map[string]reconciletest.Handler{
		"ServiceAccount":                        r.handleServiceAccount,
		"Role":                                  r.handleRole,
		"RoleBinding":                           r.handleRoleBinding,
		"ClusterRole":                           r.handleClusterRole,
		"ClusterRoleBinding":                    r.handleClusterRoleBinding,
		"Deployment":                            r.handleDeployment,
		"Service":                               r.handleService,
		"KubeletService":                        r.handleKubeletService,
		"KubeletServiceEndpoints":               r.handleKubeletServiceEndpoints,
		"KubeSchedulerService":                  r.handleKubeSchedulerService,
		"KubeSchedulerServiceEndpoints":         r.handleKubeSchedulerServiceEndpoints,
		"KubeControllerManagerService":          r.handleKubeControllerManagerService,
		"KubeControllerManagerServiceEndpoints": r.handleKubeControllerManagerServiceEndpoints,
		"ServiceMonitor":                        r.handleServiceMonitor,
		"SecurityContextConstraints":            r.handleSecurityContextConstraints,
	})
}
