package kubernetes_monitors

import (
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestHandleApiServerServiceMonitorSkipsNoOpUpdate(t *testing.T) {
	cr := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "platformmonitoring",
			Namespace: "monitoring",
			UID:       types.UID("platform-monitoring-uid"),
		},
	}
	desired, err := kubernetesMonitorsApiServerServiceMonitor(cr)
	require.NoError(t, err)
	desired.ResourceVersion = "1"

	scheme := runtime.NewScheme()
	require.NoError(t, monv1.AddToScheme(scheme))
	require.NoError(t, promv1.AddToScheme(scheme))
	r := &KubernetesMonitorsReconciler{
		ComponentReconciler: &utils.ComponentReconciler{
			Client: fake.NewClientBuilder().WithScheme(scheme).WithObjects(cr, desired.DeepCopy()).Build(),
			Scheme: scheme,
			Log:    utils.Logger("kubernetes_monitors_skip_test"),
		},
	}
	require.NoError(t, r.handleApiServerServiceMonitor(cr))

	got := &promv1.ServiceMonitor{}
	require.NoError(t, r.Client.Get(t.Context(), client.ObjectKeyFromObject(desired), got))
	assert.Equal(t, "1", got.ResourceVersion, "no-op reconcile must not bump resourceVersion")
}
