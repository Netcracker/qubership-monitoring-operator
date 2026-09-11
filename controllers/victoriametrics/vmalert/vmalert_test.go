package vmalert

import (
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	vmetricsv1b1 "github.com/VictoriaMetrics/operator/api/operator/v1beta1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	cr *monv1.PlatformMonitoring
)

func TestVmAlertManifests(t *testing.T) {
	cr = &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
		Spec: monv1.PlatformMonitoringSpec{
			Victoriametrics: &monv1.Victoriametrics{
				VmAlert: monv1.VmAlert{Image: "example:v1"},
			},
		},
	}
	t.Run("Test vmAlert manifest with nil labels and annotation", func(t *testing.T) {
		m, err := vmAlert(nil, cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "vmAlert manifest should not be empty")
		assert.NotNil(t, m.GetLabels())
		assert.NotNil(t, m.GetAnnotations())
		require.NotNil(t, m.Spec.SecurityContext)
		require.NotNil(t, m.Spec.SecurityContext.RunAsNonRoot)
		require.NotNil(t, m.Spec.SecurityContext.AllowPrivilegeEscalation)
		require.NotNil(t, m.Spec.SecurityContext.ReadOnlyRootFilesystem)
		assert.Equal(t, true, *m.Spec.SecurityContext.RunAsNonRoot)
		assert.Equal(t, false, *m.Spec.SecurityContext.AllowPrivilegeEscalation)
		assert.Equal(t, true, *m.Spec.SecurityContext.ReadOnlyRootFilesystem)
		assert.Contains(t, m.Spec.Volumes, utils.TmpVolume("100Mi"))
		assert.Contains(t, m.Spec.VolumeMounts, utils.TmpVolumeMount())
	})
	cr = &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
	}
}

func TestVmAlertUsesOperatorServiceURLs(t *testing.T) {
	for _, test := range []struct {
		name   string
		tls    bool
		scheme string
	}{
		{name: "HTTP", scheme: "http"},
		{name: "HTTPS", tls: true, scheme: "https"},
	} {
		t.Run(test.name, func(t *testing.T) {
			cr := &monv1.PlatformMonitoring{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "monitoring",
				},
				Spec: monv1.PlatformMonitoringSpec{
					Victoriametrics: &monv1.Victoriametrics{
						TLSEnabled: test.tls,
						VmOperator: monv1.VmOperator{
							Image: "vmoperator:current",
						},
						VmSingle: monv1.VmSingle{
							Image: "vmsingle:current",
						},
						VmCluster: monv1.VmCluster{
							VmSelect:       &vmetricsv1b1.VMSelect{},
							VmSelectImage:  "vmselect:current",
							VmInsert:       &vmetricsv1b1.VMInsert{},
							VmInsertImage:  "vminsert:current",
							VmStorageImage: "vmstorage:current",
						},
						VmAlert: monv1.VmAlert{
							Image: "vmalert:current",
						},
						VmAlertManager: monv1.VmAlertManager{
							Image: "vmalertmanager:current",
						},
					},
				},
			}

			manifest, err := vmAlert(nil, cr)
			require.NoError(t, err)

			assert.Equal(
				t,
				test.scheme+"://vminsert-k8s.monitoring.svc:8480/insert/0/prometheus",
				manifest.Spec.RemoteWrite.URL,
			)
			assert.Equal(
				t,
				test.scheme+"://vmselect-k8s.monitoring.svc:8481/select/0/prometheus",
				manifest.Spec.Datasource.URL,
			)
			assert.Equal(
				t,
				test.scheme+"://vmalertmanager-k8s.monitoring.svc:9093",
				manifest.Spec.Notifier.URL,
			)
			if test.tls {
				assert.NotNil(t, manifest.Spec.RemoteWrite.TLSConfig)
				assert.NotNil(t, manifest.Spec.Datasource.TLSConfig)
				assert.NotNil(t, manifest.Spec.Notifier.TLSConfig)
			} else {
				assert.Nil(t, manifest.Spec.RemoteWrite.TLSConfig)
				assert.Nil(t, manifest.Spec.Datasource.TLSConfig)
				assert.Nil(t, manifest.Spec.Notifier.TLSConfig)
			}
		})
	}
}
