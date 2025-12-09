package vmalert

import (
	"embed"
	"errors"
	"maps"
	"strings"

	v1alpha1 "github.com/Netcracker/qubership-monitoring-operator/api/v1alpha1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/victoriametrics"
	vmetricsv1b1 "github.com/VictoriaMetrics/operator/api/operator/v1beta1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/api/networking/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/yaml"
)

const (
	vmuiAlertSource    = "vmui/#/?g0.expr={{.Expr|queryEscape}}"
	grafanaAlertSource = "explore?orgId=1&left=[\"now-1h\",\"now\",\"VictoriaMetrics\",{\"expr\":{{$expr|jsonEscape|queryEscape}} },{\"mode\":\"Metrics\"},{\"ui\":[true,true,true,\"none\"]}]"
	vmalertAlertSource = "vmalert/alert?group_id={{.GroupID}}&alert_id={{.AlertID}}"
)

//go:embed  assets/*.yaml
var assets embed.FS

func vmAlertServiceAccount(cr *v1alpha1.PlatformMonitoring) (*corev1.ServiceAccount, error) {
	sa := corev1.ServiceAccount{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmAlertServiceAccountAsset), 100).Decode(&sa); err != nil {
		return nil, err
	}
	//Set parameters
	sa.SetGroupVersionKind(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "ServiceAccount"})
	sa.SetName(cr.GetNamespace() + "-" + utils.VmAlertComponentName)
	sa.SetNamespace(cr.GetNamespace())

	return &sa, nil
}

func vmAlertClusterRole(cr *v1alpha1.PlatformMonitoring, hasPsp, hasScc bool) (*rbacv1.ClusterRole, error) {
	clusterRole := rbacv1.ClusterRole{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmAlertClusterRoleAsset), 100).Decode(&clusterRole); err != nil {
		return nil, err
	}
	//Set parameters
	clusterRole.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRole"})
	clusterRole.SetName(cr.GetNamespace() + "-" + utils.VmAlertComponentName)
	if hasPsp {
		clusterRole.Rules = append(clusterRole.Rules, rbacv1.PolicyRule{
			Resources:     []string{"podsecuritypolicies"},
			Verbs:         []string{"use"},
			APIGroups:     []string{"policy"},
			ResourceNames: []string{utils.VmOperatorComponentName},
		})
	}
	if hasScc {
		clusterRole.Rules = append(clusterRole.Rules, rbacv1.PolicyRule{
			Resources:     []string{"securitycontextconstraints"},
			Verbs:         []string{"use"},
			APIGroups:     []string{"security.openshift.io"},
			ResourceNames: []string{utils.VmOperatorComponentName},
		})
	}

	return &clusterRole, nil
}

func vmAlertClusterRoleBinding(cr *v1alpha1.PlatformMonitoring) (*rbacv1.ClusterRoleBinding, error) {
	clusterRoleBinding := rbacv1.ClusterRoleBinding{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmAlertClusterRoleBindingAsset), 100).Decode(&clusterRoleBinding); err != nil {
		return nil, err
	}
	//Set parameters
	clusterRoleBinding.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRoleBinding"})
	clusterRoleBinding.SetName(cr.GetNamespace() + "-" + utils.VmAlertComponentName)
	clusterRoleBinding.RoleRef.Name = cr.GetNamespace() + "-" + utils.VmAlertComponentName

	// Set namespace for all subjects
	for it := range clusterRoleBinding.Subjects {
		sub := &clusterRoleBinding.Subjects[it]
		sub.Namespace = cr.GetNamespace()
		sub.Name = cr.GetNamespace() + "-" + utils.VmAlertComponentName
	}
	return &clusterRoleBinding, nil
}

func vmAlert(r *VmAlertReconciler, cr *v1alpha1.PlatformMonitoring) (*vmetricsv1b1.VMAlert, error) {
	var err error
	vmalert := vmetricsv1b1.VMAlert{}
	if err = yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmAlertAsset), 100).Decode(&vmalert); err != nil {
		return nil, err
	}

	// Set parameters
	vmalert.SetNamespace(cr.GetNamespace())

	if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.VmAlert.IsInstall() {

		vmalert.Spec.Image.Repository, vmalert.Spec.Image.Tag = utils.SplitImage(cr.Spec.Victoriametrics.VmAlert.Image)

		if r != nil {
			// Set security context
			if cr.Spec.Victoriametrics.VmAlert.SecurityContext != nil {
				if vmalert.Spec.SecurityContext == nil {
					vmalert.Spec.SecurityContext = &vmetricsv1b1.SecurityContext{}
				}
				if cr.Spec.Victoriametrics.VmAlert.SecurityContext.RunAsUser != nil {
					vmalert.Spec.SecurityContext.RunAsUser = cr.Spec.Victoriametrics.VmAlert.SecurityContext.RunAsUser
				}
				if cr.Spec.Victoriametrics.VmAlert.SecurityContext.FSGroup != nil {
					vmalert.Spec.SecurityContext.FSGroup = cr.Spec.Victoriametrics.VmAlert.SecurityContext.FSGroup
				}
			}
		}

		// Set resources for vmAlert
		vmalert.Spec.ServiceAccountName = cr.GetNamespace() + "-" + utils.VmAlertComponentName

		// Set resources for VmAlert
		if cr.Spec.Victoriametrics.VmAlert.Resources.Size() > 0 {
			vmalert.Spec.Resources = cr.Spec.Victoriametrics.VmAlert.Resources
		}
		// Set secrets for vmAlert deployment
		if len(cr.Spec.Victoriametrics.VmAlert.Secrets) > 0 {
			vmalert.Spec.Secrets = cr.Spec.Victoriametrics.VmAlert.Secrets
		}

		if cr.Spec.Victoriametrics.VmAlert.Replicas != nil {
			vmalert.Spec.ReplicaCount = cr.Spec.Victoriametrics.VmAlert.Replicas
		}
		if cr.Spec.Victoriametrics.VmCluster.IsInstall() && cr.Spec.Victoriametrics.VmReplicas != nil {
			vmalert.Spec.ReplicaCount = cr.Spec.Victoriametrics.VmReplicas
		}

		// Set additional volumes
		if cr.Spec.Victoriametrics.VmAlert.Volumes != nil {
			vmalert.Spec.Volumes = cr.Spec.Victoriametrics.VmAlert.Volumes
		}

		// Set additional volumeMounts for each vmAlert container. The current container names are:
		// `vmalert`, `config-reloader`
		if cr.Spec.Victoriametrics.VmAlert.VolumeMounts != nil {
			for it := range vmalert.Spec.Containers {
				c := &vmalert.Spec.Containers[it]

				// Set additional volumeMounts only for vmAlert container
				if c.Name == utils.VmAlertComponentName {
					copy(c.VolumeMounts, cr.Spec.Victoriametrics.VmAlert.VolumeMounts)
				}
			}
		}

		if cr.Spec.Victoriametrics.VmAlert.NodeSelector != nil {
			vmalert.Spec.NodeSelector = cr.Spec.Victoriametrics.VmAlert.NodeSelector
		}

		// Set affinity for vmAlert
		if cr.Spec.Victoriametrics.VmAlert.Affinity != nil {
			vmalert.Spec.Affinity = cr.Spec.Victoriametrics.VmAlert.Affinity
		}

		// Set tolerations for vmAlert
		if cr.Spec.Victoriametrics.VmAlert.Tolerations != nil {
			vmalert.Spec.Tolerations = cr.Spec.Victoriametrics.VmAlert.Tolerations
		}

		// Set additional containers
		if cr.Spec.Victoriametrics.VmAlert.Containers != nil {
			vmalert.Spec.Containers = cr.Spec.Victoriametrics.VmAlert.Containers
		}

		if cr.Spec.Victoriametrics.VmAlert.EvaluationInterval != "" {
			vmalert.Spec.EvaluationInterval = cr.Spec.Victoriametrics.VmAlert.EvaluationInterval
		}

		vmalert.Spec.SelectAllByDefault = cr.Spec.Victoriametrics.VmAlert.SelectAllByDefault

		if cr.Spec.Victoriametrics.VmAlert.RuleSelector != nil {
			vmalert.Spec.RuleSelector = cr.Spec.Victoriametrics.VmAlert.RuleSelector
		}

		if cr.Spec.Victoriametrics.VmAlert.RuleNamespaceSelector != nil {
			vmalert.Spec.RuleNamespaceSelector = cr.Spec.Victoriametrics.VmAlert.RuleNamespaceSelector
		}

		if cr.Spec.Victoriametrics.VmAlert.Port != "" {
			vmalert.Spec.Port = cr.Spec.Victoriametrics.VmAlert.Port
		}

		if cr.Spec.Victoriametrics.VmAlert.RemoteWrite != nil {
			vmalert.Spec.RemoteWrite = cr.Spec.Victoriametrics.VmAlert.RemoteWrite
		} else {
			if cr.Spec.Victoriametrics.VmOperator.IsInstall() && cr.Spec.Victoriametrics.VmSingle.IsInstall() {
				vmSingle := vmetricsv1b1.VMSingle{}
				vmSingle.SetName(utils.VmComponentName)
				vmSingle.SetNamespace(cr.GetNamespace())
				if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.TLSEnabled {
					if vmSingle.Spec.ExtraArgs == nil {
						vmSingle.Spec.ExtraArgs = make(map[string]string)
					}
					maps.Copy(vmSingle.Spec.ExtraArgs, map[string]string{"tls": "true"})
					vmalert.Spec.RemoteWrite = &vmetricsv1b1.VMAlertRemoteWriteSpec{URL: vmSingle.AsURL()}
					vmalert.Spec.RemoteWrite.TLSConfig = &vmetricsv1b1.TLSConfig{
						CAFile:   "/etc/vm/secrets/" + victoriametrics.GetVmalertTLSSecretName(cr.Spec.Victoriametrics.VmAlert) + "/ca.crt",
						CertFile: "/etc/vm/secrets/" + victoriametrics.GetVmalertTLSSecretName(cr.Spec.Victoriametrics.VmAlert) + "/tls.crt",
						KeyFile:  "/etc/vm/secrets/" + victoriametrics.GetVmalertTLSSecretName(cr.Spec.Victoriametrics.VmAlert) + "/tls.key",
					}
				} else {
					vmalert.Spec.RemoteWrite = &vmetricsv1b1.VMAlertRemoteWriteSpec{URL: vmSingle.AsURL()}
				}
			}
			if cr.Spec.Victoriametrics.VmOperator.IsInstall() && cr.Spec.Victoriametrics.VmCluster.IsInstall() {
				vmCluster := vmetricsv1b1.VMCluster{}
				vmCluster.SetName(utils.VmComponentName)
				vmCluster.SetNamespace(cr.GetNamespace())
				vmCluster.Spec.VMInsert = cr.Spec.Victoriametrics.VmCluster.VmInsert
				if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.TLSEnabled {
					if vmCluster.Spec.VMInsert.ExtraArgs == nil {
						vmCluster.Spec.VMInsert.ExtraArgs = make(map[string]string)
					}
					maps.Copy(vmCluster.Spec.VMInsert.ExtraArgs, map[string]string{"tls": "true"})
					vmalert.Spec.RemoteWrite = &vmetricsv1b1.VMAlertRemoteWriteSpec{URL: vmCluster.VMInsertURL() + "/insert/0/prometheus"}
					vmalert.Spec.RemoteWrite.TLSConfig = &vmetricsv1b1.TLSConfig{
						CAFile:   "/etc/vm/secrets/" + victoriametrics.GetVmalertTLSSecretName(cr.Spec.Victoriametrics.VmAlert) + "/ca.crt",
						CertFile: "/etc/vm/secrets/" + victoriametrics.GetVmalertTLSSecretName(cr.Spec.Victoriametrics.VmAlert) + "/tls.crt",
						KeyFile:  "/etc/vm/secrets/" + victoriametrics.GetVmalertTLSSecretName(cr.Spec.Victoriametrics.VmAlert) + "/tls.key",
					}
				} else {
					vmalert.Spec.RemoteWrite = &vmetricsv1b1.VMAlertRemoteWriteSpec{URL: vmCluster.VMInsertURL() + "/insert/0/prometheus"}
				}
			}
		}

		if cr.Spec.Victoriametrics.VmAlert.RemoteRead != nil {
			vmalert.Spec.RemoteRead = cr.Spec.Victoriametrics.VmAlert.RemoteRead
		}

		if len(cr.Spec.Victoriametrics.VmAlert.RulePath) > 0 {
			vmalert.Spec.RulePath = cr.Spec.Victoriametrics.VmAlert.RulePath
		}

		for _, notifier := range cr.Spec.Victoriametrics.VmAlert.Notifiers {
			if len(strings.TrimSpace(notifier.URL)) != 0 || notifier.Selector != nil {
				vmalert.Spec.Notifiers = append(vmalert.Spec.Notifiers, notifier)
			}
		}

		if cr.Spec.Victoriametrics.VmAlert.Notifier != nil &&
			(len(strings.TrimSpace(cr.Spec.Victoriametrics.VmAlert.Notifier.URL)) != 0 ||
				cr.Spec.Victoriametrics.VmAlert.Notifier.Selector != nil) {
			vmalert.Spec.Notifier = cr.Spec.Victoriametrics.VmAlert.Notifier
		} else {
			if len(cr.Spec.Victoriametrics.VmAlert.Notifiers) == 0 &&
				cr.Spec.Victoriametrics.VmAlert.NotifierConfigRef == nil &&
				cr.Spec.Victoriametrics.VmOperator.IsInstall() &&
				cr.Spec.Victoriametrics.VmAlertManager.IsInstall() {

				if cr.Spec.Victoriametrics.VmSingle.IsInstall() {
					vmAlertManager := vmetricsv1b1.VMAlertmanager{}
					vmAlertManager.SetName(utils.VmComponentName)
					vmAlertManager.SetNamespace(cr.GetNamespace())
					if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.TLSEnabled {
						vmAlertManager.Spec.WebConfig = &vmetricsv1b1.AlertmanagerWebConfig{
							TLSServerConfig: &vmetricsv1b1.TLSServerConfig{
								Certs: vmetricsv1b1.Certs{
									CertSecretRef: &corev1.SecretKeySelector{
										Key: "tls.crt",
										LocalObjectReference: corev1.LocalObjectReference{
											Name: victoriametrics.GetVmalertmanagerTLSSecretName(cr.Spec.Victoriametrics.VmAlertManager),
										},
									},
									KeySecretRef: &corev1.SecretKeySelector{
										Key: "tls.key",
										LocalObjectReference: corev1.LocalObjectReference{
											Name: victoriametrics.GetVmalertmanagerTLSSecretName(cr.Spec.Victoriametrics.VmAlertManager),
										},
									},
								},
							},
						}
						vmalert.Spec.Notifier = &vmetricsv1b1.VMAlertNotifierSpec{URL: vmAlertManager.AsURL()}
						vmalert.Spec.Notifier.TLSConfig = &vmetricsv1b1.TLSConfig{
							CAFile:   "/etc/vm/secrets/" + victoriametrics.GetVmalertTLSSecretName(cr.Spec.Victoriametrics.VmAlert) + "/ca.crt",
							CertFile: "/etc/vm/secrets/" + victoriametrics.GetVmalertTLSSecretName(cr.Spec.Victoriametrics.VmAlert) + "/tls.crt",
							KeyFile:  "/etc/vm/secrets/" + victoriametrics.GetVmalertTLSSecretName(cr.Spec.Victoriametrics.VmAlert) + "/tls.key",
						}
					} else {
						vmalert.Spec.Notifier = &vmetricsv1b1.VMAlertNotifierSpec{URL: vmAlertManager.AsURL()}
					}
				}

				if cr.Spec.Victoriametrics.VmCluster.IsInstall() {
					vmAlertManager := vmetricsv1b1.VMAlertmanager{}
					vmAlertManager.SetName(utils.VmComponentName)
					vmAlertManager.SetNamespace(cr.GetNamespace())
					vmAlertManager.Spec.ReplicaCount = cr.Spec.Victoriametrics.VmAlertManager.Replicas
					if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.TLSEnabled {
						vmAlertManager.Spec.WebConfig = &vmetricsv1b1.AlertmanagerWebConfig{
							TLSServerConfig: &vmetricsv1b1.TLSServerConfig{
								Certs: vmetricsv1b1.Certs{
									CertSecretRef: &corev1.SecretKeySelector{
										Key: "tls.crt",
										LocalObjectReference: corev1.LocalObjectReference{
											Name: victoriametrics.GetVmalertmanagerTLSSecretName(cr.Spec.Victoriametrics.VmAlertManager),
										},
									},
									KeySecretRef: &corev1.SecretKeySelector{
										Key: "tls.key",
										LocalObjectReference: corev1.LocalObjectReference{
											Name: victoriametrics.GetVmalertmanagerTLSSecretName(cr.Spec.Victoriametrics.VmAlertManager),
										},
									},
								},
							},
						}
						var notifiers []vmetricsv1b1.VMAlertNotifierSpec
						for _, notifier := range vmAlertManager.AsNotifiers() {
							notifier.TLSConfig = &vmetricsv1b1.TLSConfig{
								CAFile:   "/etc/vm/secrets/" + victoriametrics.GetVmalertTLSSecretName(cr.Spec.Victoriametrics.VmAlert) + "/ca.crt",
								CertFile: "/etc/vm/secrets/" + victoriametrics.GetVmalertTLSSecretName(cr.Spec.Victoriametrics.VmAlert) + "/tls.crt",
								KeyFile:  "/etc/vm/secrets/" + victoriametrics.GetVmalertTLSSecretName(cr.Spec.Victoriametrics.VmAlert) + "/tls.key",
							}
							notifiers = append(notifiers, notifier)
						}
						vmalert.Spec.Notifiers = notifiers
					} else {
						vmalert.Spec.Notifiers = vmAlertManager.AsNotifiers()
					}
				}
			}
		}

		if cr.Spec.Victoriametrics.VmAlert.NotifierConfigRef != nil {
			if cr.Spec.Victoriametrics.VmAlert.Notifier == nil &&
				len(cr.Spec.Victoriametrics.VmAlert.Notifiers) == 0 {
				vmalert.Spec.NotifierConfigRef = cr.Spec.Victoriametrics.VmAlert.NotifierConfigRef
			} else {
				return nil, errors.New("only one of notifier options could be chosen: notifierConfigRef or notifiers + notifier")
			}
		}

		if cr.Spec.Victoriametrics.VmAlert.Datasource != nil {
			vmalert.Spec.Datasource = *cr.Spec.Victoriametrics.VmAlert.Datasource
		} else {
			if cr.Spec.Victoriametrics.VmOperator.IsInstall() && cr.Spec.Victoriametrics.VmSingle.IsInstall() {
				vmSingle := vmetricsv1b1.VMSingle{}
				vmSingle.SetName(utils.VmComponentName)
				vmSingle.SetNamespace(cr.GetNamespace())
				if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.TLSEnabled {
					vmSingle.Spec.ExtraArgs = make(map[string]string)
					maps.Copy(vmSingle.Spec.ExtraArgs, map[string]string{"tls": "true"})
					vmalert.Spec.Datasource = vmetricsv1b1.VMAlertDatasourceSpec{URL: vmSingle.AsURL()}
					vmalert.Spec.Datasource.TLSConfig = &vmetricsv1b1.TLSConfig{
						CAFile:   "/etc/vm/secrets/" + victoriametrics.GetVmalertTLSSecretName(cr.Spec.Victoriametrics.VmAlert) + "/ca.crt",
						CertFile: "/etc/vm/secrets/" + victoriametrics.GetVmalertTLSSecretName(cr.Spec.Victoriametrics.VmAlert) + "/tls.crt",
						KeyFile:  "/etc/vm/secrets/" + victoriametrics.GetVmalertTLSSecretName(cr.Spec.Victoriametrics.VmAlert) + "/tls.key",
					}
				} else {
					vmalert.Spec.Datasource = vmetricsv1b1.VMAlertDatasourceSpec{URL: vmSingle.AsURL()}
				}
			}
			if cr.Spec.Victoriametrics.VmOperator.IsInstall() && cr.Spec.Victoriametrics.VmCluster.IsInstall() {
				vmCluster := vmetricsv1b1.VMCluster{}
				vmCluster.SetName(utils.VmComponentName)
				vmCluster.SetNamespace(cr.GetNamespace())
				vmCluster.Spec.VMSelect = cr.Spec.Victoriametrics.VmCluster.VmSelect
				if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.TLSEnabled {
					if vmCluster.Spec.VMSelect.ExtraArgs == nil {
						vmCluster.Spec.VMSelect.ExtraArgs = make(map[string]string)
					}
					maps.Copy(vmCluster.Spec.VMSelect.ExtraArgs, map[string]string{"tls": "true"})
					vmalert.Spec.Datasource = vmetricsv1b1.VMAlertDatasourceSpec{URL: vmCluster.VMSelectURL() + "/select/0/prometheus"}
					vmalert.Spec.Datasource.TLSConfig = &vmetricsv1b1.TLSConfig{
						CAFile:   "/etc/vm/secrets/" + victoriametrics.GetVmalertTLSSecretName(cr.Spec.Victoriametrics.VmAlert) + "/ca.crt",
						CertFile: "/etc/vm/secrets/" + victoriametrics.GetVmalertTLSSecretName(cr.Spec.Victoriametrics.VmAlert) + "/tls.crt",
						KeyFile:  "/etc/vm/secrets/" + victoriametrics.GetVmalertTLSSecretName(cr.Spec.Victoriametrics.VmAlert) + "/tls.key",
					}
				} else {
					vmalert.Spec.Datasource = vmetricsv1b1.VMAlertDatasourceSpec{URL: vmCluster.VMSelectURL() + "/select/0/prometheus"}
				}
			}
		}

		vmalert.Spec.ExternalLabels = cr.Spec.Victoriametrics.VmAlert.ExternalLabels
		vmalert.Spec.ExtraArgs = cr.Spec.Victoriametrics.VmAlert.ExtraArgs
		vmalert.Spec.ExtraEnvs = cr.Spec.Victoriametrics.VmAlert.ExtraEnvs

		alertExternalUrl := ""
		alertExternalSource := vmalertAlertSource

		if cr.Spec.Victoriametrics.VmAuth.IsInstall() {
			alertExternalSource = vmuiAlertSource
			if cr.Spec.Victoriametrics.VmAuth.Ingress != nil && cr.Spec.Victoriametrics.VmAuth.Ingress.IsInstall() {
				alertExternalUrl = "http://" + cr.Spec.Victoriametrics.VmAuth.Ingress.Host
			}
		} else if cr.Spec.Victoriametrics.VmSingle.IsInstall() {
			alertExternalSource = vmuiAlertSource
			if cr.Spec.Victoriametrics.VmSingle.Ingress != nil && cr.Spec.Victoriametrics.VmSingle.Ingress.IsInstall() {
				alertExternalUrl = "http://" + cr.Spec.Victoriametrics.VmSingle.Ingress.Host
			}
		} else if cr.Spec.Grafana != nil && cr.Spec.Grafana.IsInstall() {
			alertExternalSource = grafanaAlertSource
			if cr.Spec.Grafana.Ingress != nil && cr.Spec.Grafana.Ingress.IsInstall() {
				alertExternalUrl = "http://" + cr.Spec.Grafana.Ingress.Host
			}
		} else {
			if cr.Spec.Victoriametrics.VmAlert.Ingress != nil && cr.Spec.Victoriametrics.VmAlert.Ingress.IsInstall() {
				alertExternalUrl = "http://" + cr.Spec.Victoriametrics.VmAlert.Ingress.Host
			}
		}

		if alertExternalUrl != "" && alertExternalSource != "" {
			if vmalert.Spec.ExtraArgs == nil {
				vmalert.Spec.ExtraArgs = map[string]string{}
			}
			if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.TLSEnabled {
				alertExternalUrl = strings.Replace(alertExternalUrl, "http", "https", 1)
			}
			vmalert.Spec.ExtraArgs["external.url"] = alertExternalUrl
			vmalert.Spec.ExtraArgs["external.alert.source"] = alertExternalSource
		}

		if cr.Spec.Victoriametrics.VmAlert.TerminationGracePeriodSeconds != nil {
			vmalert.Spec.TerminationGracePeriodSeconds = cr.Spec.Victoriametrics.VmAlert.TerminationGracePeriodSeconds
		}

		if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.TLSEnabled {
			vmalert.Spec.Secrets = append(vmalert.Spec.Secrets, victoriametrics.GetVmalertTLSSecretName(cr.Spec.Victoriametrics.VmAlert))

			if vmalert.Spec.ExtraArgs == nil {
				vmalert.Spec.ExtraArgs = make(map[string]string)
			}
			maps.Copy(vmalert.Spec.ExtraArgs, map[string]string{"tls": "true"})
			maps.Copy(vmalert.Spec.ExtraArgs, map[string]string{"tlsCertFile": "/etc/vm/secrets/" + victoriametrics.GetVmalertTLSSecretName(cr.Spec.Victoriametrics.VmAlert) + "/tls.crt"})
			maps.Copy(vmalert.Spec.ExtraArgs, map[string]string{"tlsKeyFile": "/etc/vm/secrets/" + victoriametrics.GetVmalertTLSSecretName(cr.Spec.Victoriametrics.VmAlert) + "/tls.key"})
		}

		// Set labels
		vmalert.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(vmalert.GetName(), vmalert.GetNamespace())
		vmalert.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Victoriametrics.VmAlert.Image)

		vmalert.Spec.PodMetadata = &vmetricsv1b1.EmbeddedObjectMetadata{Labels: map[string]string{
			"name":                         utils.TruncLabel(vmalert.GetName()),
			"app.kubernetes.io/name":       utils.TruncLabel(vmalert.GetName()),
			"app.kubernetes.io/instance":   utils.GetInstanceLabel(vmalert.GetName(), vmalert.GetNamespace()),
			"app.kubernetes.io/component":  "victoriametrics",
			"app.kubernetes.io/part-of":    "monitoring",
			"app.kubernetes.io/managed-by": "monitoring-operator",
			"app.kubernetes.io/version":    utils.GetTagFromImage(cr.Spec.Victoriametrics.VmAlert.Image),
		}}

		if vmalert.Spec.PodMetadata != nil {
			if cr.Spec.Victoriametrics.VmAlert.Labels != nil {
				for k, v := range cr.Spec.Victoriametrics.VmAlert.Labels {
					vmalert.Spec.PodMetadata.Labels[k] = v
				}
			}

			if vmalert.Spec.PodMetadata.Annotations == nil && cr.Spec.Victoriametrics.VmAlert.Annotations != nil {
				vmalert.Spec.PodMetadata.Annotations = cr.Spec.Victoriametrics.VmAlert.Annotations
			} else {
				for k, v := range cr.Spec.Victoriametrics.VmAlert.Annotations {
					vmalert.Spec.PodMetadata.Annotations[k] = v
				}
			}
		}

		if cr.Spec.Victoriametrics.VmAlert.EnforcedNamespaceLabel != nil {
			vmalert.Spec.EnforcedNamespaceLabel = *cr.Spec.Victoriametrics.VmAlert.EnforcedNamespaceLabel
		}

		if len(strings.TrimSpace(cr.Spec.Victoriametrics.VmAlert.PriorityClassName)) > 0 {
			vmalert.Spec.PriorityClassName = cr.Spec.Victoriametrics.VmAlert.PriorityClassName
		}
	}

	return &vmalert, nil
}

func vmAlertIngressV1beta1(cr *v1alpha1.PlatformMonitoring) (*v1beta1.Ingress, error) {
	ingress := v1beta1.Ingress{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmAlertIngressAsset), 100).Decode(&ingress); err != nil {
		return nil, err
	}
	//Set parameters
	ingress.SetGroupVersionKind(schema.GroupVersionKind{Group: "networking.k8s.io", Version: "v1beta1", Kind: "Ingress"})
	ingress.SetName(cr.GetNamespace() + "-" + utils.VmAlertServiceName)
	ingress.SetNamespace(cr.GetNamespace())

	if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.VmAlert.Ingress != nil && cr.Spec.Victoriametrics.VmAlert.Ingress.IsInstall() {
		// Check that ingress host is specified.
		if cr.Spec.Victoriametrics.VmAlert.Ingress.Host == "" {
			return nil, errors.New("host for ingress can not be empty")
		}

		// Add rule for vmalert UI
		rule := v1beta1.IngressRule{Host: cr.Spec.Victoriametrics.VmAlert.Ingress.Host}
		serviceName := utils.VmAlertServiceName
		servicePort := intstr.FromInt(utils.VmAlertServicePort)
		// If VMAuth is enabled, move routing to the VMAuth service to make VMAlert UI available from this Ingress
		if cr.Spec.Victoriametrics.VmAuth.IsInstall() {
			serviceName = utils.VmAuthServiceName
			servicePort = intstr.FromInt(utils.VmAuthServicePort)
		}

		rule.HTTP = &v1beta1.HTTPIngressRuleValue{
			Paths: []v1beta1.HTTPIngressPath{
				{
					Path: "/",
					Backend: v1beta1.IngressBackend{
						ServiceName: serviceName,
						ServicePort: servicePort,
					},
				},
			},
		}
		ingress.Spec.Rules = []v1beta1.IngressRule{rule}

		// Configure TLS if TLS secret name is set
		if cr.Spec.Victoriametrics.VmAlert.Ingress.TLSSecretName != "" {
			ingress.Spec.TLS = []v1beta1.IngressTLS{
				{
					Hosts:      []string{cr.Spec.Victoriametrics.VmAlert.Ingress.Host},
					SecretName: cr.Spec.Victoriametrics.VmAlert.Ingress.TLSSecretName,
				},
			}
		}

		if cr.Spec.Victoriametrics.VmAlert.Ingress.IngressClassName != nil {
			ingress.Spec.IngressClassName = cr.Spec.Victoriametrics.VmAlert.Ingress.IngressClassName
		}

		vmAlertAnnotations := cr.Spec.Victoriametrics.VmAlert.Ingress.Annotations
		// If VMAuth is enabled, add "nginx.ingress.kubernetes.io/app-root: /vmalert" annotation
		// to make ONLY VMAlert UI available from this Ingress
		if cr.Spec.Victoriametrics.VmAuth.IsInstall() {

			if vmAlertAnnotations == nil {
				vmAlertAnnotations = make(map[string]string)
			}
			vmAlertAnnotations[utils.NginxIngressAppRootAnnotation] = utils.VmAlertAppRootEndpoint
		}
		// Set annotations
		ingress.SetAnnotations(vmAlertAnnotations)
		if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.TLSEnabled {
			if ingress.GetAnnotations() == nil {
				annotation := make(map[string]string)
				annotation["nginx.ingress.kubernetes.io/backend-protocol"] = "HTTPS"
				ingress.SetAnnotations(annotation)
			} else {
				ingress.GetAnnotations()["nginx.ingress.kubernetes.io/backend-protocol"] = "HTTPS"
			}
		}

		// Set labels with saving default labels
		ingress.Labels["name"] = utils.TruncLabel(ingress.GetName())
		ingress.Labels["app.kubernetes.io/name"] = utils.TruncLabel(ingress.GetName())
		ingress.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(ingress.GetName(), ingress.GetNamespace())
		ingress.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Victoriametrics.VmAlert.Image)

		for lKey, lValue := range cr.Spec.Victoriametrics.VmAlert.Ingress.Labels {
			ingress.GetLabels()[lKey] = lValue
		}
	}
	return &ingress, nil
}

func vmAlertIngressV1(cr *v1alpha1.PlatformMonitoring) (*networkingv1.Ingress, error) {
	ingress := networkingv1.Ingress{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmAlertIngressAsset), 100).Decode(&ingress); err != nil {
		return nil, err
	}
	//Set parameters
	ingress.SetGroupVersionKind(schema.GroupVersionKind{Group: "networking.k8s.io", Version: "v1", Kind: "Ingress"})
	ingress.SetName(cr.GetNamespace() + "-" + utils.VmAlertServiceName)
	ingress.SetNamespace(cr.GetNamespace())

	if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.VmAlert.Ingress != nil && cr.Spec.Victoriametrics.VmAlert.Ingress.IsInstall() {
		// Check that ingress host is specified.
		if cr.Spec.Victoriametrics.VmAlert.Ingress.Host == "" {
			return nil, errors.New("host for ingress can not be empty")
		}

		serviceName := utils.VmAlertServiceName
		servicePort := int32(utils.VmAlertServicePort)
		// If VMAuth is enabled, move routing to the VMAuth service to make VMAlert UI available from this Ingress
		if cr.Spec.Victoriametrics.VmAuth.IsInstall() {
			serviceName = utils.VmAuthServiceName
			servicePort = int32(utils.VmAuthServicePort)
		}

		pathType := networkingv1.PathTypePrefix
		// Add rule for vmalert UI
		rule := networkingv1.IngressRule{Host: cr.Spec.Victoriametrics.VmAlert.Ingress.Host}
		rule.HTTP = &networkingv1.HTTPIngressRuleValue{
			Paths: []networkingv1.HTTPIngressPath{
				{
					Path:     "/",
					PathType: &pathType,
					Backend: networkingv1.IngressBackend{
						Service: &networkingv1.IngressServiceBackend{
							Name: serviceName,
							Port: networkingv1.ServiceBackendPort{
								Number: servicePort,
							},
						},
					},
				},
			},
		}

		ingress.Spec.Rules = []networkingv1.IngressRule{rule}

		// Configure TLS if TLS secret name is set
		if cr.Spec.Victoriametrics.VmAlert.Ingress.TLSSecretName != "" {
			ingress.Spec.TLS = []networkingv1.IngressTLS{
				{
					Hosts:      []string{cr.Spec.Victoriametrics.VmAlert.Ingress.Host},
					SecretName: cr.Spec.Victoriametrics.VmAlert.Ingress.TLSSecretName,
				},
			}
		}

		if cr.Spec.Victoriametrics.VmAlert.Ingress.IngressClassName != nil {
			ingress.Spec.IngressClassName = cr.Spec.Victoriametrics.VmAlert.Ingress.IngressClassName
		}

		vmAlertAnnotations := cr.Spec.Victoriametrics.VmAlert.Ingress.Annotations
		// If VMAuth is enabled, add "nginx.ingress.kubernetes.io/app-root: /vmalert" annotation
		// to make ONLY VMAlert UI available from this Ingress
		if cr.Spec.Victoriametrics.VmAuth.IsInstall() {

			if vmAlertAnnotations == nil {
				vmAlertAnnotations = make(map[string]string)
			}
			vmAlertAnnotations[utils.NginxIngressAppRootAnnotation] = utils.VmAlertAppRootEndpoint
		}
		// Set annotations
		ingress.SetAnnotations(vmAlertAnnotations)
		if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.TLSEnabled {
			if ingress.GetAnnotations() == nil {
				annotation := make(map[string]string)
				annotation["nginx.ingress.kubernetes.io/backend-protocol"] = "HTTPS"
				ingress.SetAnnotations(annotation)
			} else {
				ingress.GetAnnotations()["nginx.ingress.kubernetes.io/backend-protocol"] = "HTTPS"
			}
		}

		// Set labels with saving default labels
		ingress.Labels["name"] = ingress.GetName()
		ingress.Labels["app.kubernetes.io/name"] = ingress.GetName()
		ingress.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(ingress.GetName(), ingress.GetNamespace())
		ingress.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Victoriametrics.VmAlert.Image)

		for lKey, lValue := range cr.Spec.Victoriametrics.VmAlert.Ingress.Labels {
			ingress.GetLabels()[lKey] = lValue
		}
	}
	return &ingress, nil
}
