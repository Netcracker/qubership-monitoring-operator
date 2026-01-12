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

// findProjectRoot finds the project root by looking for go.mod file
func findProjectRoot() string {
	// Start from the directory of this test file
	_, b, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(b)
	
	// Walk up the directory tree looking for go.mod
	for {
		goModPath := filepath.Join(currentDir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return currentDir
		}
		
		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			// Reached filesystem root, fallback to relative path
			break
		}
		currentDir = parent
	}
	
	// Fallback: use relative path from test file location
	return filepath.Join(filepath.Dir(b), "../../..")
}

func getPlatformMonitoringCRDPath() string {
	projectRoot := findProjectRoot()
	relativePath := filepath.Join("charts", "qubership-monitoring-operator", "crds", "monitoring.qubership.org_platformmonitorings.yaml")
	
	// Get current working directory
	wd, _ := os.Getwd()
	
	// Try multiple possible paths
	possiblePaths := []string{
		// From project root (found via go.mod)
		filepath.Join(projectRoot, relativePath),
		// From current working directory
		filepath.Join(wd, relativePath),
		// Relative to current directory
		filepath.Join(".", relativePath),
		// Try if projectRoot has extra qubership-monitoring-operator subdirectory
		filepath.Join(projectRoot, "qubership-monitoring-operator", relativePath),
		// Try if wd has extra qubership-monitoring-operator subdirectory
		filepath.Join(wd, "qubership-monitoring-operator", relativePath),
	}
	
	// Check which path exists
	for _, path := range possiblePaths {
		if absPath, err := filepath.Abs(path); err == nil {
			if _, err := os.Stat(absPath); err == nil {
				return absPath
			}
		}
	}
	
	// Return the first path as default (will fail with clear error message)
	absPath, _ := filepath.Abs(possiblePaths[0])
	return absPath
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
