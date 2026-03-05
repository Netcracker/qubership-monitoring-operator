package labelsassert

import (
	"testing"

	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	"github.com/stretchr/testify/assert"
)

// CRLabelKeys are the required label keys for Custom Resources (ServiceMonitor, PodMonitor, etc.)
// per the label specification: base + app.kubernetes.io/processed-by-operator.
var CRLabelKeys = []string{
	"name",
	"app.kubernetes.io/name",
	"app.kubernetes.io/component",
	"app.kubernetes.io/part-of",
	"app.kubernetes.io/managed-by",
	"app.kubernetes.io/managed-by-operator",
	"app.kubernetes.io/processed-by-operator",
}

// AssertCRLabels verifies that labels on a CR (ServiceMonitor, PodMonitor, etc.) meet the label specification.
// It asserts all required keys are present, common values are correct, and optional CR labels are merged.
func AssertCRLabels(t *testing.T, labels map[string]string, expectedComponent, expectedProcessedBy string, expectedCRLabels map[string]string) {
	t.Helper()
	assert.NotNil(t, labels, "labels should not be nil")

	for _, key := range CRLabelKeys {
		val, ok := labels[key]
		assert.True(t, ok, "required label %q should be present", key)
		assert.NotEmpty(t, val, "label %q should not be empty", key)
	}

	assert.Equal(t, utils.PartOfMonitoring, labels["app.kubernetes.io/part-of"], "part-of should match")
	assert.Equal(t, utils.ManagedByOperator, labels["app.kubernetes.io/managed-by"], "managed-by should match")
	assert.Equal(t, utils.OperatorDeploymentName, labels["app.kubernetes.io/managed-by-operator"], "managed-by-operator should match")
	assert.Equal(t, expectedComponent, labels["app.kubernetes.io/component"], "component should match")
	assert.Equal(t, expectedProcessedBy, labels["app.kubernetes.io/processed-by-operator"], "processed-by-operator should match")

	for k, expected := range expectedCRLabels {
		assert.Equal(t, expected, labels[k], "CR label %q should be merged from cr.GetLabels()", k)
	}
}
