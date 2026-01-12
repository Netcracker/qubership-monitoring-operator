package v1beta1

import (
	"bufio"
	"github.com/stretchr/testify/assert"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

var (
	// Root folder of the project
	_, b, _, _ = runtime.Caller(0)
	RootDir    = filepath.Join(filepath.Dir(b), "../../..")
)

func getPlatformMonitoringCRDPath() string {
	// Try multiple possible paths
	possiblePaths := []string{
		filepath.Join(RootDir, "charts", "qubership-monitoring-operator", "crds", "monitoring.qubership.org_platformmonitorings.yaml"),
		filepath.Join(RootDir, "qubership-monitoring-operator", "charts", "qubership-monitoring-operator", "crds", "monitoring.qubership.org_platformmonitorings.yaml"),
	}
	
	// Check which path exists
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	
	// Return the first path as default (will fail with clear error message)
	return possiblePaths[0]
}

func TestPlatformMonitoringCRDManifest(t *testing.T) {
	cr := PlatformMonitoring{}
	manifestPath := getPlatformMonitoringCRDPath()
	f, err := os.Open(manifestPath)
	if err != nil {
		t.Fatalf("Failed to open CRD manifest at %s: %v", manifestPath, err)
	}
	defer f.Close()
	
	err = k8syaml.NewYAMLOrJSONDecoder(bufio.NewReader(f), 100).Decode(&cr)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, cr, "Custom resource manifest should not be empty")
}
