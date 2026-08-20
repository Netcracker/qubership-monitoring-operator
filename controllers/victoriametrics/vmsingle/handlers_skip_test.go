package vmsingle

import (
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	vmetricsv1b1 "github.com/VictoriaMetrics/operator/api/operator/v1beta1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestHandleVmSingleSkipsNoOpUpdate(t *testing.T) {
	cr := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "platformmonitoring",
			Namespace: "monitoring",
			UID:       types.UID("platform-monitoring-uid"),
		},
		Spec: monv1.PlatformMonitoringSpec{
			Victoriametrics: &monv1.Victoriametrics{
				VmSingle: monv1.VmSingle{Image: "vmsingle:current"},
			},
		},
	}
	r := newVmSingleSkipTestReconciler(t, cr)
	desired, err := vmSingle(r, cr)
	require.NoError(t, err)
	desired.ResourceVersion = "1"

	r = newVmSingleSkipTestReconciler(t, cr, desired.DeepCopy())
	require.NoError(t, r.handleVmSingle(cr))

	got := &vmetricsv1b1.VMSingle{}
	require.NoError(t, r.Client.Get(t.Context(), client.ObjectKeyFromObject(desired), got))
	assert.Equal(t, "1", got.ResourceVersion, "no-op reconcile must not bump resourceVersion")
}

func TestHandleVmSingleUpdatesChangedSpec(t *testing.T) {
	cr := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "platformmonitoring",
			Namespace: "monitoring",
			UID:       types.UID("platform-monitoring-uid"),
		},
		Spec: monv1.PlatformMonitoringSpec{
			Victoriametrics: &monv1.Victoriametrics{
				VmSingle: monv1.VmSingle{Image: "vmsingle:current"},
			},
		},
	}
	r := newVmSingleSkipTestReconciler(t, cr)
	desired, err := vmSingle(r, cr)
	require.NoError(t, err)
	live := desired.DeepCopy()
	live.ResourceVersion = "1"
	live.Spec.Image.Tag = "stale"

	r = newVmSingleSkipTestReconciler(t, cr, live)
	require.NoError(t, r.handleVmSingle(cr))

	got := &vmetricsv1b1.VMSingle{}
	require.NoError(t, r.Client.Get(t.Context(), client.ObjectKeyFromObject(desired), got))
	assert.NotEqual(t, "1", got.ResourceVersion, "spec change must trigger Update")
	assert.Equal(t, desired.Spec.Image, got.Spec.Image)
}

func newVmSingleSkipTestReconciler(t *testing.T, objs ...client.Object) *VmSingleReconciler {
	t.Helper()
	scheme := runtime.NewScheme()
	require.NoError(t, monv1.AddToScheme(scheme))
	require.NoError(t, vmetricsv1b1.AddToScheme(scheme))
	return &VmSingleReconciler{
		ComponentReconciler: &utils.ComponentReconciler{
			Client: fake.NewClientBuilder().WithScheme(scheme).WithObjects(objs...).Build(),
			Scheme: scheme,
			Log:    utils.Logger("vmsingle_skip_test"),
		},
	}
}
