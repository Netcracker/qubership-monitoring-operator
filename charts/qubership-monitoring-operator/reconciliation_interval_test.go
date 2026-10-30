package chart_test

import (
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestReconciliationIntervalZeroRendersZero(t *testing.T) {
	if _, err := exec.LookPath("helm"); err != nil {
		t.Skip("helm is required to render the operator chart")
	}
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	cmd := exec.Command("helm", "template", "m", filepath.Dir(file),
		"--set", "monitoringOperator.reconciliationInterval=0",
		"--show-only", "templates/operator/deployment.yaml")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("helm template: %v\n%s", err, out)
	}
	got := reconciliationIntervalValue(t, string(out))
	if got != `"0"` {
		t.Fatalf("RECONCILIATION_INTERVAL = %s, want \"0\"", got)
	}
}

func reconciliationIntervalValue(t *testing.T, rendered string) string {
	t.Helper()
	lines := strings.Split(rendered, "\n")
	for i, line := range lines {
		if strings.TrimSpace(line) != "- name: RECONCILIATION_INTERVAL" {
			continue
		}
		if i+1 >= len(lines) {
			break
		}
		valueLine := strings.TrimSpace(lines[i+1])
		const prefix = "value: "
		if !strings.HasPrefix(valueLine, prefix) {
			t.Fatalf("line after RECONCILIATION_INTERVAL = %q, want a value", valueLine)
		}
		return strings.TrimPrefix(valueLine, prefix)
	}
	t.Fatal("rendered Deployment has no RECONCILIATION_INTERVAL")
	return ""
}
