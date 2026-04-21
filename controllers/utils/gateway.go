package utils

import (
	"maps"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
)

const GatewayAPIConverterIgnoreAnnotation = "gateway-api-converter.netcracker.com/ignore"

// GetIngressAnnotationsForGateway returns annotations with the gateway-api-converter ignore annotation
// added when the GatewayAPI spec requests it.
func GetIngressAnnotationsForGateway(cr *monv1.PlatformMonitoring, annotations map[string]string) map[string]string {
	result := maps.Clone(annotations)
	if result == nil {
		result = make(map[string]string)
	}
	if cr != nil && cr.Spec.GatewayAPI != nil && cr.Spec.GatewayAPI.AddIngressIgnoreAnnotation {
		result[GatewayAPIConverterIgnoreAnnotation] = "true"
	}
	return result
}
