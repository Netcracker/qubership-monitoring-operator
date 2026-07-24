package v1

import (
	"bufio"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sruntime "k8s.io/apimachinery/pkg/runtime"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
)

var (
	// Root folder of the project
	_, b, _, _                               = runtime.Caller(0)
	RootDir                                  = filepath.Join(filepath.Dir(b), "../../..")
	PlatformMonitoringCustomResourceManifest = filepath.Join(RootDir, "qubership-monitoring-operator",
		"charts", "qubership-monitoring-operator", "crds", "monitoring.netcracker.com_platformmonitorings.yaml")
)

func TestPlatformMonitoringCRDManifest(t *testing.T) {
	cr := PlatformMonitoring{}
	f, err := os.Open(PlatformMonitoringCustomResourceManifest)
	if err != nil {
		t.Fatal(err)
	}
	err = k8syaml.NewYAMLOrJSONDecoder(bufio.NewReader(f), 100).Decode(&cr)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, cr, "Custom resource manifest should not be empty")
}

func TestAddToSchemeRegistersAPITypes(t *testing.T) {
	scheme := k8sruntime.NewScheme()
	require.NoError(t, AddToScheme(scheme))

	testCases := []struct {
		kind     string
		expected k8sruntime.Object
	}{
		{kind: "PlatformMonitoring", expected: &PlatformMonitoring{}},
		{kind: "PlatformMonitoringList", expected: &PlatformMonitoringList{}},
		{kind: "ListOptions", expected: &metav1.ListOptions{}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.kind, func(t *testing.T) {
			actual, err := scheme.New(SchemeGroupVersion.WithKind(testCase.kind))
			require.NoError(t, err)
			assert.IsType(t, testCase.expected, actual)
		})
	}
}

func TestOverridePrometheusRuleAnnotations(t *testing.T) {
	t.Run("merge annotations", func(t *testing.T) {
		rule := promv1.Rule{
			Annotations: map[string]string{
				"summary":     "Original summary",
				"description": "Original description",
			},
		}
		override := PrometheusRule{
			Annotations: map[string]string{
				"summary":     "Custom summary",
				"runbook_url": "https://example.org/runbook",
			},
		}

		override.OverridePrometheusRule(&rule)

		assert.Equal(t, "Custom summary", rule.Annotations["summary"])
		assert.Equal(t, "Original description", rule.Annotations["description"])
		assert.Equal(t, "https://example.org/runbook", rule.Annotations["runbook_url"])
	})

	t.Run("initialize annotations", func(t *testing.T) {
		rule := promv1.Rule{}
		override := PrometheusRule{Annotations: map[string]string{"summary": "Custom summary"}}

		override.OverridePrometheusRule(&rule)

		assert.Equal(t, map[string]string{"summary": "Custom summary"}, rule.Annotations)
	})

	t.Run("preserve annotations when override is omitted", func(t *testing.T) {
		rule := promv1.Rule{Annotations: map[string]string{"summary": "Original summary"}}

		(&PrometheusRule{}).OverridePrometheusRule(&rule)

		assert.Equal(t, map[string]string{"summary": "Original summary"}, rule.Annotations)
	})
}
