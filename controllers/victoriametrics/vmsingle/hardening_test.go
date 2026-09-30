package vmsingle

import (
	"errors"
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	fakediscovery "k8s.io/client-go/discovery/fake"
	ktesting "k8s.io/client-go/testing"
	"k8s.io/utils/ptr"
)

func hardeningCR() *monv1.PlatformMonitoring {
	return &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Name: "platformmonitoring", Namespace: "monitoring"},
		Spec: monv1.PlatformMonitoringSpec{
			Victoriametrics: &monv1.Victoriametrics{
				VmSingle: monv1.VmSingle{Image: "example:v1"},
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
	tests := map[string]func(spec *monv1.VmSingle){
		"root user": func(spec *monv1.VmSingle) {
			spec.SecurityContext = &monv1.SecurityContext{RunAsUser: ptr.To[int64](0)}
		},
		"reserved volume name": func(spec *monv1.VmSingle) {
			spec.Volumes = []corev1.Volume{{
				Name:         utils.TmpVolumeMount().Name,
				VolumeSource: corev1.VolumeSource{Secret: &corev1.SecretVolumeSource{SecretName: "creds"}},
			}}
		},
		"privileged container": func(spec *monv1.VmSingle) {
			spec.Containers = []corev1.Container{{
				Name:            "sidecar",
				SecurityContext: &corev1.SecurityContext{Privileged: ptr.To(true)},
			}}
		},
	}
	for name, configure := range tests {
		t.Run(name, func(t *testing.T) {
			cr := hardeningCR()
			configure(&cr.Spec.Victoriametrics.VmSingle)

			m, err := vmSingle(nil, cr)

			require.Error(t, err)
			assert.Nil(t, m)
		})
	}
}

func TestManifestReturnsDiscoveryError(t *testing.T) {
	reconciler := NewVmSingleReconciler(nil, runtime.NewScheme(), failingDiscovery())

	m, err := vmSingle(reconciler, hardeningCR())

	require.Error(t, err, "an unknown platform must not produce a manifest")
	assert.Nil(t, m)
}
