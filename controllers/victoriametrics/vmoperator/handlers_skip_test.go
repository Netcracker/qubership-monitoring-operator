package vmoperator

import (
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	kubernetesfake "k8s.io/client-go/kubernetes/fake"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestHandleServiceSkipsAPIDefaultedPortProtocol(t *testing.T) {
	cr := skipTestPlatformMonitoring()
	desired, err := vmOperatorService(cr)
	require.NoError(t, err)

	live := desired.DeepCopy()
	live.ResourceVersion = "1"
	applyAPIDefaultedServicePorts(live.Spec.Ports)

	r := newVmOperatorSkipTestReconciler(t, nil, cr, live)
	require.NoError(t, r.handleService(cr))

	got := &corev1.Service{}
	require.NoError(t, r.Client.Get(t.Context(), client.ObjectKeyFromObject(desired), got))
	assert.Equal(t, "1", got.ResourceVersion, "API-defaulted Service port protocol must not trigger Update")
	assert.Equal(t, corev1.ProtocolTCP, portProtocol(got.Spec.Ports, "webhook"))
}

func TestHandleKubeletServiceEndpointsSkipsAPIDefaultedPortProtocol(t *testing.T) {
	cr := skipTestPlatformMonitoring()
	node := &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{Name: "kind-control-plane", UID: "node-uid"},
		Status: corev1.NodeStatus{
			Addresses: []corev1.NodeAddress{{Type: corev1.NodeInternalIP, Address: "172.18.0.2"}},
		},
	}
	kubeClient := kubernetesfake.NewSimpleClientset(node)

	nodes, err := kubeClient.CoreV1().Nodes().List(t.Context(), metav1.ListOptions{})
	require.NoError(t, err)
	addresses, addressErrs := getNodeAddresses(nodes)
	require.Empty(t, addressErrs)

	desired, err := vmKubeletServiceEndpoints(cr)
	require.NoError(t, err)
	desired.Subsets[0].Addresses = addresses
	desired.Labels["name"] = utils.TruncLabel(desired.GetName())
	desired.Labels["app.kubernetes.io/name"] = utils.TruncLabel(desired.GetName())
	desired.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(desired.GetName(), desired.GetNamespace())
	desired.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Victoriametrics.VmOperator.Image)
	applyAPIDefaultedEndpointPorts(desired.Subsets)
	desired.ResourceVersion = "1"

	r := newVmOperatorSkipTestReconciler(t, kubeClient, cr, desired.DeepCopy())
	require.NoError(t, r.handleKubeletServiceEndpoints(cr))

	got := &corev1.Endpoints{}
	require.NoError(t, r.Client.Get(t.Context(), client.ObjectKeyFromObject(desired), got))
	assert.Equal(t, "1", got.ResourceVersion, "API-defaulted Endpoints port protocol must not trigger Update")
	require.NotEmpty(t, got.Subsets)
	require.NotEmpty(t, got.Subsets[0].Ports)
	assert.Equal(t, corev1.ProtocolTCP, got.Subsets[0].Ports[0].Protocol)
}

func skipTestPlatformMonitoring() *monv1.PlatformMonitoring {
	install := true
	return &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "platformmonitoring",
			Namespace: "monitoring",
			UID:       types.UID("platform-monitoring-uid"),
		},
		Spec: monv1.PlatformMonitoringSpec{
			Victoriametrics: &monv1.Victoriametrics{
				VmOperator: monv1.VmOperator{
					Install: &install,
					Image:   "victoriametrics/operator:v0.73.1",
				},
			},
		},
	}
}

func newVmOperatorSkipTestReconciler(t *testing.T, kubeClient kubernetes.Interface, objs ...client.Object) *VmOperatorReconciler {
	t.Helper()
	scheme := runtime.NewScheme()
	require.NoError(t, monv1.AddToScheme(scheme))
	require.NoError(t, corev1.AddToScheme(scheme))
	return &VmOperatorReconciler{
		KubeClient: kubeClient,
		ComponentReconciler: &utils.ComponentReconciler{
			Client: fake.NewClientBuilder().WithScheme(scheme).WithObjects(objs...).Build(),
			Scheme: scheme,
			Log:    utils.Logger("vmoperator_skip_test"),
		},
	}
}

func applyAPIDefaultedServicePorts(ports []corev1.ServicePort) {
	for i := range ports {
		if ports[i].Protocol == "" {
			ports[i].Protocol = corev1.ProtocolTCP
		}
	}
}

func applyAPIDefaultedEndpointPorts(subsets []corev1.EndpointSubset) {
	for i := range subsets {
		for j := range subsets[i].Ports {
			if subsets[i].Ports[j].Protocol == "" {
				subsets[i].Ports[j].Protocol = corev1.ProtocolTCP
			}
		}
	}
}

func portProtocol(ports []corev1.ServicePort, name string) corev1.Protocol {
	for _, p := range ports {
		if p.Name == name {
			return p.Protocol
		}
	}
	return ""
}
