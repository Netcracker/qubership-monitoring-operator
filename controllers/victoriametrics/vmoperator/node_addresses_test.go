package vmoperator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGetNodeAddressesSortsByIPThenNodeName(t *testing.T) {
	node := func(name, ip string) corev1.Node {
		return corev1.Node{
			ObjectMeta: metav1.ObjectMeta{Name: name},
			Status: corev1.NodeStatus{
				Addresses: []corev1.NodeAddress{{Type: corev1.NodeInternalIP, Address: ip}},
			},
		}
	}
	nodes := &corev1.NodeList{Items: []corev1.Node{
		node("node-c", "10.0.0.2"),
		node("node-b", "10.0.0.1"),
		node("node-a", "10.0.0.1"),
	}}

	addresses, errs := getNodeAddresses(nodes)

	require.Empty(t, errs)
	var got []string
	for _, address := range addresses {
		got = append(got, address.IP+"/"+address.TargetRef.Name)
	}
	assert.Equal(t, []string{"10.0.0.1/node-a", "10.0.0.1/node-b", "10.0.0.2/node-c"}, got)
}
