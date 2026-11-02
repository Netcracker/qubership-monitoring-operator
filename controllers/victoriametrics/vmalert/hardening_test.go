package vmalert

import (
	"context"
	"errors"
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	vmetricsv1b1 "github.com/VictoriaMetrics/operator/api/operator/v1beta1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	fakediscovery "k8s.io/client-go/discovery/fake"
	ktesting "k8s.io/client-go/testing"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func hardeningCR() *monv1.PlatformMonitoring {
	return &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Name: "platformmonitoring", Namespace: "monitoring"},
		Spec: monv1.PlatformMonitoringSpec{
			Victoriametrics: &monv1.Victoriametrics{
				VmAlert: monv1.VmAlert{Image: "example:v1"},
			},
		},
	}
}

func failingDiscovery() *fakediscovery.FakeDiscovery {
	discoveryClient := &fakediscovery.FakeDiscovery{Fake: &ktesting.Fake{}}
	discoveryClient.PrependReactor("get", "resource", func(ktesting.Action) (bool, runtime.Object, error) {
		return true, nil, errors.New("aggregated API unavailable")
	})
	return discoveryClient
}

func TestManifestRejectsConflictingHardeningSettings(t *testing.T) {
	tests := map[string]func(spec *monv1.VmAlert){
		"root user": func(spec *monv1.VmAlert) {
			spec.SecurityContext = &corev1.PodSecurityContext{RunAsUser: ptr.To[int64](0)}
		},
		"reserved volume name": func(spec *monv1.VmAlert) {
			spec.Volumes = []corev1.Volume{{
				Name:         utils.TmpVolumeMount().Name,
				VolumeSource: corev1.VolumeSource{Secret: &corev1.SecretVolumeSource{SecretName: "creds"}},
			}}
		},
		"privileged container": func(spec *monv1.VmAlert) {
			spec.Containers = []corev1.Container{{
				Name:            "sidecar",
				SecurityContext: &corev1.SecurityContext{Privileged: ptr.To(true)},
			}}
		},
	}
	for name, configure := range tests {
		t.Run(name, func(t *testing.T) {
			cr := hardeningCR()
			configure(&cr.Spec.Victoriametrics.VmAlert)

			m, err := vmAlert(nil, cr)

			require.Error(t, err)
			assert.Nil(t, m)
		})
	}
}

func TestManifestReturnsDiscoveryError(t *testing.T) {
	reconciler := NewVmAlertReconciler(nil, runtime.NewScheme(), failingDiscovery())

	m, err := vmAlert(reconciler, hardeningCR())

	require.Error(t, err, "an unknown platform must not produce a manifest")
	assert.Nil(t, m)
}

// Uninstall must succeed even when platform discovery fails, because it only needs the stable
// name and namespace of the resource.
func TestDeleteVMAlertDoesNotDependOnDiscovery(t *testing.T) {
	scheme := runtime.NewScheme()
	require.NoError(t, vmetricsv1b1.AddToScheme(scheme))
	existing := &vmetricsv1b1.VMAlert{ObjectMeta: metav1.ObjectMeta{Name: "k8s", Namespace: "monitoring"}}
	controllerClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(existing).Build()
	reconciler := NewVmAlertReconciler(controllerClient, scheme, failingDiscovery())

	require.NoError(t, reconciler.deleteVmAlert(hardeningCR()))

	err := controllerClient.Get(context.Background(), types.NamespacedName{Name: "k8s", Namespace: "monitoring"}, &vmetricsv1b1.VMAlert{})
	assert.True(t, apierrors.IsNotFound(err), "the resource must be deleted")
	require.NoError(t, reconciler.deleteVmAlert(hardeningCR()), "a missing resource is not an error")
}
