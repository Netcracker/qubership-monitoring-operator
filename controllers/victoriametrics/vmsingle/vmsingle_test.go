package vmsingle

import (
	"context"
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	vmetricsv1b1 "github.com/VictoriaMetrics/operator/api/operator/v1beta1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

type deleteOptionsRecordingClient struct {
	client.Client
	deleteOptions client.DeleteOptions
}

func (c *deleteOptionsRecordingClient) Delete(
	ctx context.Context,
	obj client.Object,
	opts ...client.DeleteOption,
) error {
	c.deleteOptions.ApplyOptions(opts)
	return c.Client.Delete(ctx, obj, opts...)
}

var (
	cr *monv1.PlatformMonitoring
	// labelKey        = "label.key"
	// labelValue      = "label-value"
	// annotationKey   = "annotation.key"
	// annotationValue = "annotation-value"
)

func TestVmSingleManifests(t *testing.T) {
	// cr = &monv1.PlatformMonitoring{
	// 	ObjectMeta: metav1.ObjectMeta{
	// 		Namespace: "monitoring",
	// 	},
	// 	Spec: monv1.PlatformMonitoringSpec{
	// 		Victoriametrics: &v1.Victoriametrics{
	// 			vmSingle: v1.vmSingle{
	// 				Annotations: map[string]string{annotationKey: annotationValue},
	// 				Labels:      map[string]string{labelKey: labelValue},
	// 			},
	// 		},
	// 	},
	// }
	// t.Run("Test vmSingle manifest", func(t *testing.T) {
	// 	m, err := vmsingle(cr)
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}
	// 	assert.NotNil(t, m, "vmSingle manifest should not be empty")
	// })
	cr = &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
		Spec: monv1.PlatformMonitoringSpec{
			Victoriametrics: &monv1.Victoriametrics{
				VmSingle: monv1.VmSingle{},
			},
		},
	}
	t.Run("Test vmSingle manifest with nil labels and annotation", func(t *testing.T) {
		m, err := vmSingle(nil, cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "vmSingle manifest should not be empty")
		assert.NotNil(t, m.GetLabels())
		assert.Nil(t, m.GetAnnotations())
	})
	cr = &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
	}

}

func TestVmSingleUsesOperatorVmAlertURL(t *testing.T) {
	cr := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
		Spec: monv1.PlatformMonitoringSpec{
			Victoriametrics: &monv1.Victoriametrics{
				VmOperator: monv1.VmOperator{
					Image: "vmoperator:current",
				},
				VmSingle: monv1.VmSingle{
					Image: "vmsingle:current",
				},
				VmAlert: monv1.VmAlert{
					Image: "vmalert:current",
				},
			},
		},
	}

	manifest, err := vmSingle(nil, cr)
	require.NoError(t, err)

	assert.Equal(
		t,
		"http://vmalert-k8s.monitoring.svc:8080",
		manifest.Spec.ExtraArgs["vmalert.proxyURL"],
	)
}

func TestVmSinglePreservesPVCByDefault(t *testing.T) {
	cr := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
		Spec: monv1.PlatformMonitoringSpec{
			Victoriametrics: &monv1.Victoriametrics{
				VmSingle: monv1.VmSingle{
					Image: "vmsingle:current",
				},
			},
		},
	}

	manifest, err := vmSingle(nil, cr)
	require.NoError(t, err)

	assert.False(t, manifest.Spec.RemovePvcAfterDelete)
}

func TestHandleVmSingleDisablesPVCRemovalOnUpgrade(t *testing.T) {
	scheme := runtime.NewScheme()
	require.NoError(t, vmetricsv1b1.AddToScheme(scheme))

	existing := &vmetricsv1b1.VMSingle{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "k8s",
			Namespace: "monitoring",
		},
		Spec: vmetricsv1b1.VMSingleSpec{
			RemovePvcAfterDelete: true,
		},
	}
	controllerClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(existing).Build()
	reconciler := NewVmSingleReconciler(controllerClient, scheme, nil)
	cr := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "platformmonitoring",
			Namespace: "monitoring",
		},
		Spec: monv1.PlatformMonitoringSpec{
			Victoriametrics: &monv1.Victoriametrics{
				VmSingle: monv1.VmSingle{
					Image: "vmsingle:current",
				},
			},
		},
	}

	require.NoError(t, reconciler.handleVmSingle(cr))

	updated := &vmetricsv1b1.VMSingle{}
	require.NoError(t, controllerClient.Get(
		context.Background(),
		types.NamespacedName{Name: "k8s", Namespace: "monitoring"},
		updated,
	))
	assert.False(t, updated.Spec.RemovePvcAfterDelete)
}

func TestDeleteVmSingleUsesOrphanPropagation(t *testing.T) {
	scheme := runtime.NewScheme()
	require.NoError(t, vmetricsv1b1.AddToScheme(scheme))

	existing := &vmetricsv1b1.VMSingle{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "k8s",
			Namespace: "monitoring",
		},
	}
	controllerClient := &deleteOptionsRecordingClient{
		Client: fake.NewClientBuilder().WithScheme(scheme).WithObjects(existing).Build(),
	}
	reconciler := NewVmSingleReconciler(controllerClient, scheme, nil)
	cr := &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "platformmonitoring",
			Namespace: "monitoring",
		},
	}

	require.NoError(t, reconciler.deleteVmSingle(cr))
	require.NotNil(t, controllerClient.deleteOptions.PropagationPolicy)
	assert.Equal(t, metav1.DeletePropagationOrphan, *controllerClient.deleteOptions.PropagationPolicy)
}

func TestVmSingleClusterRBACUsesInstallationNamespace(t *testing.T) {
	cr := &monv1.PlatformMonitoring{ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"}}

	role, err := vmSingleClusterRole(cr, false, false)
	require.NoError(t, err)
	binding, err := vmSingleClusterRoleBinding(cr)
	require.NoError(t, err)

	assert.Equal(t, "monitoring", role.Labels["monitoring.netcracker.com/installation-namespace"])
	assert.Equal(t, "monitoring", binding.Labels["monitoring.netcracker.com/installation-namespace"])
}
