package nodeexporter

import (
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils/labelsassert"
	secv1 "github.com/openshift/api/security/v1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	cr              *monv1.PlatformMonitoring
	labelKey        = "label.key"
	labelValue      = "label-value"
	annotationKey   = "annotation.key"
	annotationValue = "annotation-value"
)

func TestNodeExporterManifests(t *testing.T) {
	cr = &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
		Spec: monv1.PlatformMonitoringSpec{
			NodeExporter: &monv1.NodeExporter{
				Annotations: map[string]string{annotationKey: annotationValue},
				Labels:      map[string]string{labelKey: labelValue},
			},
		},
	}
	t.Run("Test DaemonSet manifest", func(t *testing.T) {
		m, err := nodeExporterDaemonSet(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "DaemonSet manifest should not be empty")
		assert.NotNil(t, m.GetLabels())
		assert.Equal(t, labelValue, m.GetLabels()[labelKey])
		assert.NotNil(t, m.Spec.Template.Labels)
		assert.Equal(t, labelValue, m.Spec.Template.Labels[labelKey])
		assert.NotNil(t, m.GetAnnotations())
		assert.Equal(t, annotationValue, m.GetAnnotations()[annotationKey])
		assert.NotNil(t, m.Spec.Template.Annotations)
		assert.Equal(t, annotationValue, m.Spec.Template.Annotations[annotationKey])
	})
	cr = &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
		Spec: monv1.PlatformMonitoringSpec{
			NodeExporter: &monv1.NodeExporter{},
		},
	}
	t.Run("Test DaemonSet manifest with nil labels and annotation", func(t *testing.T) {
		m, err := nodeExporterDaemonSet(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "DaemonSet manifest should not be empty")
		assert.NotNil(t, m.GetLabels())
		assert.NotNil(t, m.Spec.Template.Labels)
		assert.Nil(t, m.GetAnnotations())
		assert.Nil(t, m.Spec.Template.Annotations)
	})
	t.Run("Test ClusterRole manifest", func(t *testing.T) {
		m, err := nodeExporterClusterRole(cr, true, false)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "ClusterRole manifest should not be empty")
	})
	t.Run("Test ClusterRoleBinding manifest", func(t *testing.T) {
		m, err := nodeExporterClusterRoleBinding(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "ClusterRoleBinding manifest should not be empty")
	})
	t.Run("Test Service manifest", func(t *testing.T) {
		m, err := nodeExporterService(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Service manifest should not be empty")
	})
	t.Run("Test ServiceMonitor manifest", func(t *testing.T) {
		crWithLabels := &monv1.PlatformMonitoring{
			ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring", Labels: map[string]string{labelKey: labelValue}},
			Spec: monv1.PlatformMonitoringSpec{
				NodeExporter: &monv1.NodeExporter{
					Annotations: map[string]string{annotationKey: annotationValue},
					Labels:      map[string]string{labelKey: labelValue},
				},
			},
		}
		m, err := nodeExporterServiceMonitor(crWithLabels)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "ServiceMonitor manifest should not be empty")
		labelsassert.AssertCRLabels(t, m.GetLabels(), utils.NodeExporterComponentName, "victoriametrics-operator", map[string]string{labelKey: labelValue})
	})
	t.Run("Test ServiceAccount manifest", func(t *testing.T) {
		crWithSALabels := &monv1.PlatformMonitoring{
			ObjectMeta: metav1.ObjectMeta{Namespace: "monitoring"},
			Spec: monv1.PlatformMonitoringSpec{
				NodeExporter: &monv1.NodeExporter{
					ServiceAccount: &monv1.EmbeddedObjectMetadata{
						Labels: map[string]string{labelKey: labelValue},
					},
				},
			},
		}
		m, err := nodeExporterServiceAccount(crWithSALabels)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "ServiceAccount manifest should not be empty")
		assert.NotNil(t, m.GetLabels())
		assert.Equal(t, labelValue, m.GetLabels()[labelKey], "ServiceAccount.Labels should be merged")
	})
	t.Run("Test SecurityContextConstraints manifest", func(t *testing.T) {
		m, err := nodeExporterSecurityContextConstraints()
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "SecurityContextConstraints manifest should not be empty")
		assert.Equal(t, utils.NodeExporterComponentName, m.GetName())
		assert.True(t, m.AllowHostNetwork, "node-exporter needs hostNetwork to scrape host-level metrics")
		assert.True(t, m.AllowHostPID, "node-exporter needs hostPID to scrape host-level metrics")
		assert.True(t, m.AllowHostPorts)
		assert.True(t, m.AllowHostDirVolumePlugin, "node-exporter mounts /proc, /sys, and / as hostPath volumes")
		assert.False(t, m.AllowPrivilegedContainer)
		// Without an explicit seccompProfiles entry, the SCC admission plugin silently drops this
		// SCC from the candidate list for pods that set securityContext.seccompProfile - no error,
		// it's just never tried, and the DaemonSet falls back to the built-in SCCs and fails.
		assert.Contains(t, m.SeccompProfiles, "runtime/default")
		// MustRunAsRange (not RunAsAny) makes the SCC assign a UID from the namespace's allocated
		// range. The node-exporter image sets a non-numeric USER (nobody), so without an assigned
		// UID the kubelet can't verify runAsNonRoot and refuses to start the container.
		assert.Equal(t, secv1.RunAsUserStrategyMustRunAsRange, m.RunAsUser.Type)
	})
}
