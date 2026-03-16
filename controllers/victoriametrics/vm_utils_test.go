package victoriametrics

import (
	"testing"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
)

func TestGetIngressAnnotationsForGateway(t *testing.T) {
	t.Parallel()

	base := map[string]string{"existing": "true"}
	cr := &monv1.PlatformMonitoring{
		Spec: monv1.PlatformMonitoringSpec{
			GatewayAPI: &monv1.GatewayAPI{AddIngressIgnoreAnnotation: true},
		},
	}

	got := GetIngressAnnotationsForGateway(cr, base)
	if got[GatewayAPIConverterIgnoreAnnotation] != "true" {
		t.Fatalf("expected %s annotation to be set", GatewayAPIConverterIgnoreAnnotation)
	}
	if got["existing"] != "true" {
		t.Fatalf("expected existing annotations to be preserved")
	}
}

func TestGetIngressAnnotationsForGateway_NoGatewayAPI(t *testing.T) {
	t.Parallel()

	cr := &monv1.PlatformMonitoring{}
	got := GetIngressAnnotationsForGateway(cr, nil)
	if len(got) != 0 {
		t.Fatalf("expected no annotations, got=%v", got)
	}
}
