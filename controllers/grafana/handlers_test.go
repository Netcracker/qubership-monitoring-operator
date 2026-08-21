package grafana

import (
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestApplyGrafanaDesiredStateSkipsAPICanonicalEmptyAnnotations(t *testing.T) {
	cr := grafanaComparisonPlatformMonitoring(nil)
	desired, err := grafana(cr, false)
	require.NoError(t, err)

	live := desired.DeepCopy()
	live.Spec.Deployment.Spec.Template.Annotations = nil

	assert.False(t, applyGrafanaDesiredState(live, desired))
}

func TestApplyGrafanaDesiredStateUpdatesChangedAnnotation(t *testing.T) {
	desired, err := grafana(grafanaComparisonPlatformMonitoring(map[string]string{"example.com/key": "desired"}), false)
	require.NoError(t, err)
	live := desired.DeepCopy()
	live.Spec.Deployment.Spec.Template.Annotations["example.com/key"] = "stale"

	assert.True(t, applyGrafanaDesiredState(live, desired))
	assert.Equal(t, map[string]string{"example.com/key": "desired"}, live.Spec.Deployment.Spec.Template.Annotations)
}

func TestApplyGrafanaDesiredStateRemovesStaleAnnotation(t *testing.T) {
	desired, err := grafana(grafanaComparisonPlatformMonitoring(nil), false)
	require.NoError(t, err)
	live := desired.DeepCopy()
	live.Spec.Deployment.Spec.Template.Annotations = map[string]string{"example.com/stale": "value"}

	assert.True(t, applyGrafanaDesiredState(live, desired))
	assert.Nil(t, live.Spec.Deployment.Spec.Template.Annotations)
}

func TestApplyGrafanaDesiredStateUpdatesChangedLabels(t *testing.T) {
	desired, err := grafana(grafanaComparisonPlatformMonitoring(nil), false)
	require.NoError(t, err)
	live := desired.DeepCopy()
	live.Labels = map[string]string{"example.com/stale": "value"}

	assert.True(t, applyGrafanaDesiredState(live, desired))
	assert.Equal(t, desired.Labels, live.Labels)
}

func grafanaComparisonPlatformMonitoring(annotations map[string]string) *monv1.PlatformMonitoring {
	return &monv1.PlatformMonitoring{
		ObjectMeta: metav1.ObjectMeta{Name: "platformmonitoring", Namespace: "monitoring"},
		Spec: monv1.PlatformMonitoringSpec{
			Grafana: &monv1.Grafana{Annotations: annotations},
		},
	}
}
