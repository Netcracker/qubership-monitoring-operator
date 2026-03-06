package labelsassert

import (
	"testing"

	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
)

func TestAssertCRLabels(t *testing.T) {
	t.Run("valid labels with component and processed-by", func(t *testing.T) {
		labels := utils.MergeLabels(
			utils.ResourceLabels("my-cr", "my-component"),
			map[string]string{utils.ProcessedByOperatorKey: "my-operator"},
		)
		AssertCRLabels(t, labels, "my-component", "my-operator", nil)
	})
	t.Run("valid labels with optional CR labels", func(t *testing.T) {
		labels := utils.MergeLabels(
			utils.ResourceLabels("my-cr", "my-component"),
			map[string]string{
				utils.ProcessedByOperatorKey: "my-operator",
				"custom.key":                 "custom-value",
			},
		)
		AssertCRLabels(t, labels, "my-component", "my-operator", map[string]string{"custom.key": "custom-value"})
	})
}
