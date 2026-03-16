package victoriametrics

import (
	"maps"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
)

const GatewayAPIConverterIgnoreAnnotation = "gateway-api-converter.netcracker.com/ignore"

func GetIngressAnnotationsForGateway(cr *monv1.PlatformMonitoring, annotations map[string]string) map[string]string {
	result := maps.Clone(annotations)
	if result == nil {
		result = make(map[string]string)
	}
	if cr != nil && shouldAddIngressIgnoreAnnotation(cr) {
		result[GatewayAPIConverterIgnoreAnnotation] = "true"
	}
	return result
}

func shouldAddIngressIgnoreAnnotation(cr *monv1.PlatformMonitoring) bool {
	return cr.Spec.GatewayAPI != nil && cr.Spec.GatewayAPI.AddIngressIgnoreAnnotation
}

func GetVmalertTLSSecretName(vmalert monv1.VmAlert) string {
	if vmalert.TLSConfig != nil {
		return vmalert.TLSConfig.SecretName
	}
	return utils.VmAlertTLSSecret
}

func GetVmagentTLSSecretName(vmagent monv1.VmAgent) string {
	if vmagent.TLSConfig != nil {
		return vmagent.TLSConfig.SecretName
	}
	return utils.VmAgentTLSSecret
}

func GetVmalertmanagerTLSSecretName(vmalertmanager monv1.VmAlertManager) string {
	if vmalertmanager.TLSConfig != nil {
		return vmalertmanager.TLSConfig.SecretName
	}
	return utils.VmAlertManagerTLSSecret
}

func GetVmsingleTLSSecretName(vmsingle monv1.VmSingle) string {
	if vmsingle.TLSConfig != nil {
		return vmsingle.TLSConfig.SecretName
	}
	return utils.VmSingleTLSSecret
}

func GetVmauthTLSSecretName(vmauth monv1.VmAuth) string {
	if vmauth.TLSConfig != nil {
		return vmauth.TLSConfig.SecretName
	}
	return utils.VmAuthTLSSecret
}

func GetVmoperatorTLSSecretName(vmoperator monv1.VmOperator) string {
	if vmoperator.TLSConfig != nil {
		return vmoperator.TLSConfig.SecretName
	}
	return utils.VmOperatorTLSSecret
}

func GetVmselectTLSSecretName(vmcluster monv1.VmCluster) string {
	if vmcluster.VmSelectTLSConfig != nil {
		return vmcluster.VmSelectTLSConfig.SecretName
	}
	return utils.VmSelectTLSSecret
}

func GetVminsertTLSSecretName(vmcluster monv1.VmCluster) string {
	if vmcluster.VmInsertTLSConfig != nil {
		return vmcluster.VmInsertTLSConfig.SecretName
	}
	return utils.VmInsertTLSSecret
}

func GetVmstorageTLSSecretName(vmcluster monv1.VmCluster) string {
	if vmcluster.VmStorageTLSConfig != nil {
		return vmcluster.VmStorageTLSConfig.SecretName
	}
	return utils.VmStorageTLSSecret
}
