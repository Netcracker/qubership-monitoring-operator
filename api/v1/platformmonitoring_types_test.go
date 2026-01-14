package v1

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
