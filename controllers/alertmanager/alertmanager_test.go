package alertmanager

import (
	"testing"

	v1beta1 "github.com/Netcracker/qubership-monitoring-operator/api/v1beta1"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	cr              *v1beta1.PlatformMonitoring
	labelKey        = "label.key"
	labelValue      = "label-value"
	annotationKey   = "annotation.key"
	annotationValue = "annotation-value"
)

func TestAlertmanagerManifests(t *testing.T) {
	cr = &v1beta1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
		Spec: v1beta1.PlatformMonitoringSpec{
			AlertManager: &v1beta1.AlertManager{
				Annotations: map[string]string{annotationKey: annotationValue},
				Labels:      map[string]string{labelKey: labelValue},
			},
		},
	}
	t.Run("Test Alert Manager manifest", func(t *testing.T) {
		m, err := alertmanager(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Alert Manager manifest should not be empty")
		assert.NotNil(t, m.Spec.PodMetadata.Labels)
		assert.Equal(t, labelValue, m.Spec.PodMetadata.Labels[labelKey])
		assert.NotNil(t, m.Spec.PodMetadata.Annotations)
		assert.Equal(t, annotationValue, m.Spec.PodMetadata.Annotations[annotationKey])
	})
	cr = &v1beta1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "monitoring",
		},
		Spec: v1beta1.PlatformMonitoringSpec{
			AlertManager: &v1beta1.AlertManager{},
		},
	}
	t.Run("Test Alert Manager manifest with nil labels and annotation", func(t *testing.T) {
		m, err := alertmanager(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Alert Manager manifest should not be empty")
		assert.NotNil(t, m.GetLabels())
		assert.Nil(t, m.GetAnnotations())
	})
	t.Run("Test ServiceAccount manifest", func(t *testing.T) {
		m, err := alertmanagerServiceAccount(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "ServiceAccount manifest should not be empty")
	})
	t.Run("Test Secret manifest", func(t *testing.T) {
		m, err := alertmanagerSecret(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Secret manifest should not be empty")
	})

	t.Run("Test Service manifest", func(t *testing.T) {
		m, err := alertmanagerService(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Service manifest should not be empty")
	})
	t.Run("Test Ingress v1beta1 manifest", func(t *testing.T) {
		m, err := alertmanagerIngressV1beta1(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Ingress v1beta1 manifest should not be empty")
	})
	t.Run("Test Ingress v1 manifest", func(t *testing.T) {
		m, err := alertmanagerIngressV1(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Ingress v1 manifest should not be empty")
	})
	t.Run("Test PodMonitor manifest", func(t *testing.T) {
		m, err := alertmanagerPodMonitor(cr)
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "PodMonitor manifest should not be empty")
	})
}
