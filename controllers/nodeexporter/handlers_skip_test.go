package nodeexporter

import (
	"encoding/json"
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils/apidefaultstest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestHandleDaemonSetSkipsAPIDefaultedPodSpec(t *testing.T) {
	cr := newNodeExporterPlatformMonitoring(true)
	cr.Spec.NodeExporter.Resources = corev1.ResourceRequirements{
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("500m"),
			corev1.ResourceMemory: resource.MustParse("100Mi"),
		},
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse("50m"),
			corev1.ResourceMemory: resource.MustParse("50Mi"),
		},
	}
	runAsUser, fsGroup := int64(2000), int64(2000)
	cr.Spec.NodeExporter.SecurityContext = &monv1.SecurityContext{
		RunAsUser: &runAsUser,
		FSGroup:   &fsGroup,
	}

	desired, err := nodeExporterDaemonSet(cr, false)
	require.NoError(t, err)

	live := desired.DeepCopy()
	apidefaultstest.PodSpec(&live.Spec.Template.Spec)
	raw, err := json.Marshal(live)
	require.NoError(t, err)
	live = &appsv1.DaemonSet{}
	require.NoError(t, json.Unmarshal(raw, live))
	live.ResourceVersion = "1"

	r, kubeClient := newNodeExporterTestReconciler(t, false, cr, live)
	require.NoError(t, r.handleDaemonSet(cr))

	got := &appsv1.DaemonSet{}
	require.NoError(t, kubeClient.Get(t.Context(), client.ObjectKeyFromObject(desired), got))
	assert.Equal(t, "1", got.ResourceVersion, "API-defaulted PodSpec fields must not trigger Update")
	assert.Equal(t, corev1.DNSClusterFirst, got.Spec.Template.Spec.DNSPolicy)
	assert.Equal(t, "default-scheduler", got.Spec.Template.Spec.SchedulerName)
}

func TestHandleDaemonSetUpdatesChangedImage(t *testing.T) {
	cr := newNodeExporterPlatformMonitoring(true)
	desired, err := nodeExporterDaemonSet(cr, false)
	require.NoError(t, err)
	live := desired.DeepCopy()
	require.NotEmpty(t, live.Spec.Template.Spec.Containers)
	live.Spec.Template.Spec.Containers[0].Image = "stale:old"
	live.ResourceVersion = "1"

	r, kubeClient := newNodeExporterTestReconciler(t, false, cr, live)
	require.NoError(t, r.handleDaemonSet(cr))

	got := &appsv1.DaemonSet{}
	require.NoError(t, kubeClient.Get(t.Context(), client.ObjectKeyFromObject(desired), got))
	assert.NotEqual(t, "1", got.ResourceVersion, "image change must trigger Update")
	require.NotEmpty(t, got.Spec.Template.Spec.Containers)
	assert.Equal(t, cr.Spec.NodeExporter.Image, got.Spec.Template.Spec.Containers[0].Image)
}
