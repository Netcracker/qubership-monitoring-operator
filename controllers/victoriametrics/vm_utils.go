package victoriametrics

import (
	v1beta1 "github.com/Netcracker/qubership-monitoring-operator/api/v1beta1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
)

func GetVmalertTLSSecretName(vmalert v1beta1.VmAlert) string {
	if vmalert.TLSConfig != nil {
		return vmalert.TLSConfig.SecretName
	}
	return utils.VmAlertTLSSecret
}

func GetVmagentTLSSecretName(vmagent v1beta1.VmAgent) string {
	if vmagent.TLSConfig != nil {
		return vmagent.TLSConfig.SecretName
	}
	return utils.VmAgentTLSSecret
}

func GetVmalertmanagerTLSSecretName(vmalertmanager v1beta1.VmAlertManager) string {
	if vmalertmanager.TLSConfig != nil {
		return vmalertmanager.TLSConfig.SecretName
	}
	return utils.VmAlertManagerTLSSecret
}

func GetVmsingleTLSSecretName(vmsingle v1beta1.VmSingle) string {
	if vmsingle.TLSConfig != nil {
		return vmsingle.TLSConfig.SecretName
	}
	return utils.VmSingleTLSSecret
}

func GetVmauthTLSSecretName(vmauth v1beta1.VmAuth) string {
	if vmauth.TLSConfig != nil {
		return vmauth.TLSConfig.SecretName
	}
	return utils.VmAuthTLSSecret
}

func GetVmoperatorTLSSecretName(vmoperator v1beta1.VmOperator) string {
	if vmoperator.TLSConfig != nil {
		return vmoperator.TLSConfig.SecretName
	}
	return utils.VmOperatorTLSSecret
}

func GetVmselectTLSSecretName(vmcluster v1beta1.VmCluster) string {
	if vmcluster.VmSelectTLSConfig != nil {
		return vmcluster.VmSelectTLSConfig.SecretName
	}
	return utils.VmSelectTLSSecret
}

func GetVminsertTLSSecretName(vmcluster v1beta1.VmCluster) string {
	if vmcluster.VmInsertTLSConfig != nil {
		return vmcluster.VmInsertTLSConfig.SecretName
	}
	return utils.VmInsertTLSSecret
}

func GetVmstorageTLSSecretName(vmcluster v1beta1.VmCluster) string {
	if vmcluster.VmStorageTLSConfig != nil {
		return vmcluster.VmStorageTLSConfig.SecretName
	}
	return utils.VmStorageTLSSecret
}
