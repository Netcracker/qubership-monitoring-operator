package vmcluster

import (
	"embed"
	"maps"
	"strings"

	v1alpha1 "github.com/Netcracker/qubership-monitoring-operator/api/v1alpha1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/victoriametrics"
	vmetricsv1b1 "github.com/VictoriaMetrics/operator/api/operator/v1beta1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/yaml"
)

//go:embed  assets/*.yaml
var assets embed.FS

func vmClusterServiceAccount(cr *v1alpha1.PlatformMonitoring) (*corev1.ServiceAccount, error) {
	sa := corev1.ServiceAccount{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmClusterServiceAccountAsset), 100).Decode(&sa); err != nil {
		return nil, err
	}
	//Set parameters
	sa.SetGroupVersionKind(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "ServiceAccount"})
	sa.SetName(cr.GetNamespace() + "-" + utils.VmClusterComponentName)
	sa.SetNamespace(cr.GetNamespace())

	return &sa, nil
}

func vmClusterClusterRole(cr *v1alpha1.PlatformMonitoring, hasPsp, hasScc bool) (*rbacv1.ClusterRole, error) {
	clusterRole := rbacv1.ClusterRole{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmClusterClusterRoleAsset), 100).Decode(&clusterRole); err != nil {
		return nil, err
	}
	//Set parameters
	clusterRole.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRole"})
	clusterRole.SetName(cr.GetNamespace() + "-" + utils.VmClusterComponentName)
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

func vmClusterClusterRoleBinding(cr *v1alpha1.PlatformMonitoring) (*rbacv1.ClusterRoleBinding, error) {
	clusterRoleBinding := rbacv1.ClusterRoleBinding{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmClusterClusterRoleBindingAsset), 100).Decode(&clusterRoleBinding); err != nil {
		return nil, err
	}
	//Set parameters
	clusterRoleBinding.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRoleBinding"})
	clusterRoleBinding.SetName(cr.GetNamespace() + "-" + utils.VmClusterComponentName)
	clusterRoleBinding.RoleRef.Name = cr.GetNamespace() + "-" + utils.VmClusterComponentName

	// Set namespace for all subjects
	for it := range clusterRoleBinding.Subjects {
		sub := &clusterRoleBinding.Subjects[it]
		sub.Namespace = cr.GetNamespace()
		sub.Name = cr.GetNamespace() + "-" + utils.VmClusterComponentName
	}
	return &clusterRoleBinding, nil
}

func vmCluster(cr *v1alpha1.PlatformMonitoring) (*vmetricsv1b1.VMCluster, error) {
	var err error
	vmcluster := vmetricsv1b1.VMCluster{}
	if err = yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmClusterAsset), 100).Decode(&vmcluster); err != nil {
		return nil, err
	}

	// Set parameters
	vmcluster.SetNamespace(cr.GetNamespace())

	if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.VmCluster.IsInstall() {
		if len(cr.Spec.Victoriametrics.VmCluster.RetentionPeriod) != 0 {
			vmcluster.Spec.RetentionPeriod = cr.Spec.Victoriametrics.VmCluster.RetentionPeriod
		}
		if cr.Spec.Victoriametrics.VmCluster.ReplicationFactor != nil {
			vmcluster.Spec.ReplicationFactor = cr.Spec.Victoriametrics.VmCluster.ReplicationFactor
		}

		vmcluster.Spec.ServiceAccountName = cr.GetNamespace() + "-" + utils.VmClusterComponentName

		// Set labels
		vmcluster.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(vmcluster.GetName(), vmcluster.GetNamespace())

		if len(strings.TrimSpace(cr.Spec.Victoriametrics.VmCluster.ClusterVersion)) > 0 {
			vmcluster.Spec.ClusterVersion = cr.Spec.Victoriametrics.VmCluster.ClusterVersion
		}

		if cr.Spec.Victoriametrics.VmCluster.VmSelect != nil {
			vmcluster.Spec.VMSelect = cr.Spec.Victoriametrics.VmCluster.VmSelect
			if cr.Spec.Victoriametrics.VmReplicas != nil {
				vmcluster.Spec.VMSelect.ReplicaCount = cr.Spec.Victoriametrics.VmReplicas
			}
			vmcluster.Spec.VMSelect.Image.Repository, vmcluster.Spec.VMSelect.Image.Tag = utils.SplitImage(cr.Spec.Victoriametrics.VmCluster.VmSelectImage)
			if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.TLSEnabled {
				vmcluster.Spec.VMSelect.Secrets = append(vmcluster.Spec.VMSelect.Secrets, victoriametrics.GetVmselectTLSSecretName(cr.Spec.Victoriametrics.VmCluster))
				if vmcluster.Spec.VMSelect.ExtraArgs == nil {
					vmcluster.Spec.VMSelect.ExtraArgs = make(map[string]string)
				}
				maps.Copy(vmcluster.Spec.VMSelect.ExtraArgs, map[string]string{"tls": "true"})
				maps.Copy(vmcluster.Spec.VMSelect.ExtraArgs, map[string]string{"tlsCertFile": "/etc/vm/secrets/" + victoriametrics.GetVmselectTLSSecretName(cr.Spec.Victoriametrics.VmCluster) + "/tls.crt"})
				maps.Copy(vmcluster.Spec.VMSelect.ExtraArgs, map[string]string{"tlsKeyFile": "/etc/vm/secrets/" + victoriametrics.GetVmselectTLSSecretName(cr.Spec.Victoriametrics.VmCluster) + "/tls.key"})
			}
			vmcluster.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(vmcluster.Spec.VMSelect.GetNameWithPrefix(cr.Name), vmcluster.GetNamespace())
			vmcluster.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Victoriametrics.VmCluster.VmSelectImage)

			vmcluster.Spec.VMSelect.PodMetadata = &vmetricsv1b1.EmbeddedObjectMetadata{Labels: map[string]string{
				"name":                         utils.TruncLabel(vmcluster.Spec.VMSelect.GetNameWithPrefix(cr.Name)),
				"app.kubernetes.io/name":       utils.TruncLabel(vmcluster.Spec.VMSelect.GetNameWithPrefix(cr.Name)),
				"app.kubernetes.io/instance":   utils.GetInstanceLabel(vmcluster.Spec.VMSelect.GetNameWithPrefix(cr.Name), vmcluster.GetNamespace()),
				"app.kubernetes.io/component":  "victoriametrics",
				"app.kubernetes.io/part-of":    "monitoring",
				"app.kubernetes.io/managed-by": "monitoring-operator",
				"app.kubernetes.io/version":    utils.GetTagFromImage(cr.Spec.Victoriametrics.VmCluster.VmSelectImage),
			}}
		}

		if cr.Spec.Victoriametrics.VmCluster.VmStorage != nil {
			vmcluster.Spec.VMStorage = cr.Spec.Victoriametrics.VmCluster.VmStorage
			if cr.Spec.Victoriametrics.VmReplicas != nil {
				vmcluster.Spec.VMStorage.ReplicaCount = cr.Spec.Victoriametrics.VmReplicas
			}
			vmcluster.Spec.VMStorage.Image.Repository, vmcluster.Spec.VMStorage.Image.Tag = utils.SplitImage(cr.Spec.Victoriametrics.VmCluster.VmStorageImage)
			if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.TLSEnabled {
				vmcluster.Spec.VMStorage.Secrets = append(vmcluster.Spec.VMStorage.Secrets, victoriametrics.GetVmstorageTLSSecretName(cr.Spec.Victoriametrics.VmCluster))

				if vmcluster.Spec.VMStorage.ExtraArgs == nil {
					vmcluster.Spec.VMStorage.ExtraArgs = make(map[string]string)
				}
				maps.Copy(vmcluster.Spec.VMStorage.ExtraArgs, map[string]string{"tls": "true"})
				maps.Copy(vmcluster.Spec.VMStorage.ExtraArgs, map[string]string{"tlsCertFile": "/etc/vm/secrets/" + victoriametrics.GetVmstorageTLSSecretName(cr.Spec.Victoriametrics.VmCluster) + "/tls.crt"})
				maps.Copy(vmcluster.Spec.VMStorage.ExtraArgs, map[string]string{"tlsKeyFile": "/etc/vm/secrets/" + victoriametrics.GetVmstorageTLSSecretName(cr.Spec.Victoriametrics.VmCluster) + "/tls.key"})
			}
			vmcluster.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(vmcluster.Spec.VMSelect.GetNameWithPrefix(cr.Name), vmcluster.GetNamespace())
			vmcluster.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Victoriametrics.VmCluster.VmStorageImage)

			vmcluster.Spec.VMStorage.PodMetadata = &vmetricsv1b1.EmbeddedObjectMetadata{Labels: map[string]string{
				"name":                         utils.TruncLabel(vmcluster.Spec.VMStorage.GetNameWithPrefix(cr.Name)),
				"app.kubernetes.io/name":       utils.TruncLabel(vmcluster.Spec.VMStorage.GetNameWithPrefix(cr.Name)),
				"app.kubernetes.io/instance":   utils.GetInstanceLabel(vmcluster.Spec.VMStorage.GetNameWithPrefix(cr.Name), vmcluster.GetNamespace()),
				"app.kubernetes.io/component":  "victoriametrics",
				"app.kubernetes.io/part-of":    "monitoring",
				"app.kubernetes.io/managed-by": "monitoring-operator",
				"app.kubernetes.io/version":    utils.GetTagFromImage(cr.Spec.Victoriametrics.VmCluster.VmStorageImage),
			}}
		}

		if cr.Spec.Victoriametrics.VmCluster.VmInsert != nil {
			vmcluster.Spec.VMInsert = cr.Spec.Victoriametrics.VmCluster.VmInsert
			if cr.Spec.Victoriametrics.VmReplicas != nil {
				vmcluster.Spec.VMInsert.ReplicaCount = cr.Spec.Victoriametrics.VmReplicas
			}
			vmcluster.Spec.VMInsert.Image.Repository, vmcluster.Spec.VMInsert.Image.Tag = utils.SplitImage(cr.Spec.Victoriametrics.VmCluster.VmInsertImage)
			if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.TLSEnabled {
				vmcluster.Spec.VMInsert.Secrets = append(vmcluster.Spec.VMInsert.Secrets, victoriametrics.GetVminsertTLSSecretName(cr.Spec.Victoriametrics.VmCluster))

				if vmcluster.Spec.VMInsert.ExtraArgs == nil {
					vmcluster.Spec.VMInsert.ExtraArgs = make(map[string]string)
				}
				maps.Copy(vmcluster.Spec.VMInsert.ExtraArgs, map[string]string{"tls": "true"})
				maps.Copy(vmcluster.Spec.VMInsert.ExtraArgs, map[string]string{"tlsCertFile": "/etc/vm/secrets/" + victoriametrics.GetVminsertTLSSecretName(cr.Spec.Victoriametrics.VmCluster) + "/tls.crt"})
				maps.Copy(vmcluster.Spec.VMInsert.ExtraArgs, map[string]string{"tlsKeyFile": "/etc/vm/secrets/" + victoriametrics.GetVminsertTLSSecretName(cr.Spec.Victoriametrics.VmCluster) + "/tls.key"})
			}
			vmcluster.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(vmcluster.Spec.VMInsert.GetNameWithPrefix(cr.Name), vmcluster.GetNamespace())
			vmcluster.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Victoriametrics.VmCluster.VmInsertImage)

			vmcluster.Spec.VMInsert.PodMetadata = &vmetricsv1b1.EmbeddedObjectMetadata{Labels: map[string]string{
				"name":                         utils.TruncLabel(vmcluster.Spec.VMInsert.GetNameWithPrefix(cr.Name)),
				"app.kubernetes.io/name":       utils.TruncLabel(vmcluster.Spec.VMInsert.GetNameWithPrefix(cr.Name)),
				"app.kubernetes.io/instance":   utils.GetInstanceLabel(vmcluster.Spec.VMInsert.GetNameWithPrefix(cr.Name), vmcluster.GetNamespace()),
				"app.kubernetes.io/component":  "victoriametrics",
				"app.kubernetes.io/part-of":    "monitoring",
				"app.kubernetes.io/managed-by": "monitoring-operator",
				"app.kubernetes.io/version":    utils.GetTagFromImage(cr.Spec.Victoriametrics.VmCluster.VmInsertImage),
			}}
		}
	}

	return &vmcluster, nil
}

func vmSelectIngressV1(cr *v1alpha1.PlatformMonitoring) (*networkingv1.Ingress, error) {
	ingress := networkingv1.Ingress{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmSelectIngressAsset), 100).Decode(&ingress); err != nil {
		return nil, err
	}
	//Set metadata
	ingress.SetGroupVersionKind(schema.GroupVersionKind{Group: "networking.k8s.io", Version: "v1", Kind: "Ingress"})
	ingress.SetName(cr.GetNamespace() + "-" + utils.VmSelectServiceName)
	ingress.SetNamespace(cr.GetNamespace())

	if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.VmCluster.VmSelectIngress != nil && cr.Spec.Victoriametrics.VmCluster.VmSelectIngress.IsInstall() {
		var rules []networkingv1.IngressRule
		pathType := networkingv1.PathTypePrefix
		ing := cr.Spec.Victoriametrics.VmCluster.VmSelectIngress

		switch {
		// 1. If Host is provided
		case ing.Host != "":
			rules = append(rules, networkingv1.IngressRule{
				Host: ing.Host,
				IngressRuleValue: networkingv1.IngressRuleValue{
					HTTP: &networkingv1.HTTPIngressRuleValue{
						Paths: []networkingv1.HTTPIngressPath{defaultVmSelectPath(pathType)},
					},
				},
			})

		// 2. If custom ingress rules provided
		case len(ing.Rules) > 0:
			for _, r := range ing.Rules {
				// fallback if HTTP is not set
				if r.HTTP == nil || len(r.HTTP.Paths) == 0 {
					r.HTTP = &v1alpha1.HTTPIngressRuleValue{
						Paths: []v1alpha1.IngressPath{
							{
								Path:     "/",
								PathType: string(pathType),
								Backend: v1alpha1.IngressPathBackend{
									Service: v1alpha1.IngressPathBackendService{
										Name: utils.VmSelectServiceName,
										Port: v1alpha1.ServiceBackendPort{
											Number: utils.VmSelectServicePort,
										},
									},
								},
							},
						},
					}
				}

				// converting to k8s networkingv1
				var paths []networkingv1.HTTPIngressPath
				for _, p := range r.HTTP.Paths {
					pt := networkingv1.PathTypePrefix
					if p.PathType != "" {
						pt = networkingv1.PathType(p.PathType)
					}

					backendPort := networkingv1.ServiceBackendPort{}
					if p.Backend.Service.Port.Number != 0 {
						backendPort.Number = p.Backend.Service.Port.Number
					} else {
						backendPort.Name = p.Backend.Service.Port.Name
					}

					paths = append(paths, networkingv1.HTTPIngressPath{
						Path:     p.Path,
						PathType: &pt,
						Backend: networkingv1.IngressBackend{
							Service: &networkingv1.IngressServiceBackend{
								Name: p.Backend.Service.Name,
								Port: backendPort,
							},
						},
					})
				}

				rules = append(rules, networkingv1.IngressRule{
					Host: r.Host,
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{Paths: paths},
					},
				})
			}

		// 3. fallback: if no Host or no custom ingress rules provided
		default:
			rules = append(rules, networkingv1.IngressRule{
				IngressRuleValue: networkingv1.IngressRuleValue{
					HTTP: &networkingv1.HTTPIngressRuleValue{
						Paths: []networkingv1.HTTPIngressPath{defaultVmSelectPath(pathType)},
					},
				},
			})
		}
		ingress.Spec.Rules = rules

		tlsConfigured := false
		pickSecret := func(ingressTLSSecret string, tlsCfg *v1alpha1.VmTLSConfig) string {
			if ingressTLSSecret != "" {
				return ingressTLSSecret
			}
			if tlsCfg != nil {
				return tlsCfg.SecretName
			}
			return ""
		}
		// Configure tls if TLS config is defined
		if !tlsConfigured && len(cr.Spec.Victoriametrics.VmCluster.VmSelectIngress.TLS) > 0 {
			for _, hostgroup := range cr.Spec.Victoriametrics.VmCluster.VmSelectIngress.TLS {
				if len(hostgroup.Hosts) == 0 {
					continue
				}
				validHosts := make([]string, 0, len(hostgroup.Hosts))
				for _, h := range hostgroup.Hosts {
					if strings.TrimSpace(h) != "" {
						validHosts = append(validHosts, h)
					}
				}
				if len(validHosts) == 0 {
					continue
				}
				secret := hostgroup.SecretName
				// fallback: if secretName is empty - use TLSSecretName
				if secret == "" {
					secret = pickSecret(cr.Spec.Victoriametrics.VmCluster.VmSelectIngress.TLSSecretName, cr.Spec.Victoriametrics.VmCluster.VmSelectTLSConfig)
				}
				if secret != "" {
					ingress.Spec.TLS = append(ingress.Spec.TLS, networkingv1.IngressTLS{
						Hosts:      validHosts,
						SecretName: secret,
					})
				}
			}
			if len(ingress.Spec.TLS) > 0 {
				tlsConfigured = true
			}
		}
		// Configure TLS if TLS secret name and host is set
		if !tlsConfigured && cr.Spec.Victoriametrics.VmCluster.VmSelectIngress.Host != "" {
			secret := pickSecret(cr.Spec.Victoriametrics.VmCluster.VmSelectIngress.TLSSecretName, cr.Spec.Victoriametrics.VmCluster.VmSelectTLSConfig)
			if secret != "" {
				ingress.Spec.TLS = []networkingv1.IngressTLS{
					{
						Hosts:      []string{cr.Spec.Victoriametrics.VmCluster.VmSelectIngress.Host},
						SecretName: secret,
					},
				}
				tlsConfigured = true
			}
		}
		// Fallback: use ingress rules to configure tls hosts and TLSSecretName
		if !tlsConfigured && len(cr.Spec.Victoriametrics.VmCluster.VmSelectIngress.Rules) > 0 {
			tlsHosts := []string{}
			secret := pickSecret(cr.Spec.Victoriametrics.VmCluster.VmSelectIngress.TLSSecretName, cr.Spec.Victoriametrics.VmCluster.VmSelectTLSConfig)
			for _, rule := range cr.Spec.Victoriametrics.VmCluster.VmSelectIngress.Rules {
				if rule.Host != "" {
					tlsHosts = append(tlsHosts, rule.Host)
				}
			}
			if len(tlsHosts) > 0 && secret != "" {
				ingress.Spec.TLS = []networkingv1.IngressTLS{
					{
						Hosts:      tlsHosts,
						SecretName: secret,
					},
				}
			}
		}

		if cr.Spec.Victoriametrics.VmCluster.VmSelectIngress.IngressClassName != nil {
			ingress.Spec.IngressClassName = cr.Spec.Victoriametrics.VmCluster.VmSelectIngress.IngressClassName
		}

		// Set annotations
		ingress.SetAnnotations(cr.Spec.Victoriametrics.VmCluster.VmSelectIngress.Annotations)
		if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.TLSEnabled {
			if ingress.GetAnnotations() == nil {
				annotation := make(map[string]string)
				annotation["nginx.ingress.kubernetes.io/backend-protocol"] = "HTTPS"
				ingress.SetAnnotations(annotation)
			} else {
				ingress.GetAnnotations()["nginx.ingress.kubernetes.io/backend-protocol"] = "HTTPS"
			}
		}

		if ingress.GetAnnotations() == nil {
			annotation := make(map[string]string)
			annotation["nginx.ingress.kubernetes.io/app-root"] = "/select/0/vmui"
			ingress.SetAnnotations(annotation)
		} else {
			ingress.GetAnnotations()["nginx.ingress.kubernetes.io/app-root"] = "/select/0/vmui"
		}

		// Set labels with saving default labels
		ingress.Labels["name"] = utils.TruncLabel(ingress.GetName())
		ingress.Labels["app.kubernetes.io/name"] = utils.TruncLabel(ingress.GetName())
		ingress.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(ingress.GetName(), ingress.GetNamespace())
		ingress.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Victoriametrics.VmCluster.VmSelectImage)

		for lKey, lValue := range cr.Spec.Victoriametrics.VmCluster.VmSelectIngress.Labels {
			ingress.GetLabels()[lKey] = lValue
		}
	}
	return &ingress, nil
}

func defaultVmSelectPath(pathType networkingv1.PathType) networkingv1.HTTPIngressPath {
	return networkingv1.HTTPIngressPath{
		Path:     "/",
		PathType: &pathType,
		Backend: networkingv1.IngressBackend{
			Service: &networkingv1.IngressServiceBackend{
				Name: utils.VmSelectServiceName,
				Port: networkingv1.ServiceBackendPort{
					Number: utils.VmSelectServicePort,
				},
			},
		},
	}
}
