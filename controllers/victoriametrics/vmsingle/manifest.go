package vmsingle

import (
	"embed"
	"strings"

	"maps"

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

func vmSingleServiceAccount(cr *v1alpha1.PlatformMonitoring) (*corev1.ServiceAccount, error) {
	sa := corev1.ServiceAccount{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmSingleServiceAccountAsset), 100).Decode(&sa); err != nil {
		return nil, err
	}
	//Set parameters
	sa.SetGroupVersionKind(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "ServiceAccount"})
	sa.SetName(cr.GetNamespace() + "-" + utils.VmSingleComponentName)
	sa.SetNamespace(cr.GetNamespace())

	return &sa, nil
}

func vmSingleClusterRole(cr *v1alpha1.PlatformMonitoring, hasPsp, hasScc bool) (*rbacv1.ClusterRole, error) {
	clusterRole := rbacv1.ClusterRole{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmSingleClusterRoleAsset), 100).Decode(&clusterRole); err != nil {
		return nil, err
	}
	//Set parameters
	clusterRole.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRole"})
	clusterRole.SetName(cr.GetNamespace() + "-" + utils.VmSingleComponentName)
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

func vmSingleClusterRoleBinding(cr *v1alpha1.PlatformMonitoring) (*rbacv1.ClusterRoleBinding, error) {
	clusterRoleBinding := rbacv1.ClusterRoleBinding{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmSingleClusterRoleBindingAsset), 100).Decode(&clusterRoleBinding); err != nil {
		return nil, err
	}
	//Set parameters
	clusterRoleBinding.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRoleBinding"})
	clusterRoleBinding.SetName(cr.GetNamespace() + "-" + utils.VmSingleComponentName)
	clusterRoleBinding.RoleRef.Name = cr.GetNamespace() + "-" + utils.VmSingleComponentName

	// Set namespace for all subjects
	for it := range clusterRoleBinding.Subjects {
		sub := &clusterRoleBinding.Subjects[it]
		sub.Namespace = cr.GetNamespace()
		sub.Name = cr.GetNamespace() + "-" + utils.VmSingleComponentName
	}
	return &clusterRoleBinding, nil
}

func vmSingle(r *VmSingleReconciler, cr *v1alpha1.PlatformMonitoring) (*vmetricsv1b1.VMSingle, error) {
	var err error
	vmsingle := vmetricsv1b1.VMSingle{}
	if err = yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmSingleAsset), 100).Decode(&vmsingle); err != nil {
		return nil, err
	}

	// Set parameters
	vmsingle.SetNamespace(cr.GetNamespace())

	if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.VmSingle.IsInstall() {
		vmsingle.Spec.RetentionPeriod = cr.Spec.Victoriametrics.VmSingle.RetentionPeriod

		// Set vmsingle image
		vmsingle.Spec.Image.Repository, vmsingle.Spec.Image.Tag = utils.SplitImage(cr.Spec.Victoriametrics.VmSingle.Image)

		if r != nil {
			// Set security context
			if cr.Spec.Victoriametrics.VmSingle.SecurityContext != nil {
				if vmsingle.Spec.SecurityContext == nil {
					vmsingle.Spec.SecurityContext = &vmetricsv1b1.SecurityContext{}
				}
				if cr.Spec.Victoriametrics.VmSingle.SecurityContext.RunAsUser != nil {
					vmsingle.Spec.SecurityContext.RunAsUser = cr.Spec.Victoriametrics.VmSingle.SecurityContext.RunAsUser
				}
				if cr.Spec.Victoriametrics.VmSingle.SecurityContext.RunAsGroup != nil {
					vmsingle.Spec.SecurityContext.RunAsGroup = cr.Spec.Victoriametrics.VmSingle.SecurityContext.RunAsGroup
				}
				if cr.Spec.Victoriametrics.VmSingle.SecurityContext.FSGroup != nil {
					vmsingle.Spec.SecurityContext.FSGroup = cr.Spec.Victoriametrics.VmSingle.SecurityContext.FSGroup
				}
			}
		}

		vmsingle.Spec.ServiceAccountName = cr.GetNamespace() + "-" + utils.VmSingleComponentName

		// Set resources for vmsingle deployment
		if cr.Spec.Victoriametrics.VmSingle.Resources.Size() > 0 {
			vmsingle.Spec.Resources = cr.Spec.Victoriametrics.VmSingle.Resources
		}
		// Set additional containers
		if cr.Spec.Victoriametrics.VmSingle.Containers != nil {
			vmsingle.Spec.Containers = cr.Spec.Victoriametrics.VmSingle.Containers
		}
		// Set secrets for vmsingle deployment
		if len(cr.Spec.Victoriametrics.VmSingle.Secrets) > 0 {
			vmsingle.Spec.Secrets = cr.Spec.Victoriametrics.VmSingle.Secrets
		}

		// Set additional volumes
		if cr.Spec.Victoriametrics.VmSingle.Volumes != nil {
			vmsingle.Spec.Volumes = cr.Spec.Victoriametrics.VmSingle.Volumes
		}

		if cr.Spec.Victoriametrics.VmSingle.VolumeMounts != nil {
			for it := range vmsingle.Spec.Containers {
				c := &vmsingle.Spec.Containers[it]

				// Set additional volumeMounts only for vmsingle container
				if c.Name == utils.VmSingleComponentName {
					c.VolumeMounts = cr.Spec.Victoriametrics.VmSingle.VolumeMounts
				}
			}
		}
		if len(strings.TrimSpace(cr.Spec.Victoriametrics.VmSingle.StorageDataPath)) > 0 {
			vmsingle.Spec.StorageDataPath = cr.Spec.Victoriametrics.VmSingle.StorageDataPath
		}
		if cr.Spec.Victoriametrics.VmSingle.Storage != nil {
			vmsingle.Spec.Storage = cr.Spec.Victoriametrics.VmSingle.Storage
		}
		if cr.Spec.Victoriametrics.VmSingle.StorageMetadata != nil {
			vmsingle.Spec.StorageMetadata = *cr.Spec.Victoriametrics.VmSingle.StorageMetadata
		}
		// Set nodeSelector for vmsingle deployment
		if cr.Spec.Victoriametrics.VmSingle.NodeSelector != nil {
			vmsingle.Spec.NodeSelector = cr.Spec.Victoriametrics.VmSingle.NodeSelector
		}
		// Set affinity for vmsingle deployment
		if cr.Spec.Victoriametrics.VmSingle.Affinity != nil {
			vmsingle.Spec.Affinity = cr.Spec.Victoriametrics.VmSingle.Affinity
		}

		if cr.Spec.Victoriametrics.VmSingle.ExtraArgs != nil {
			maps.Copy(vmsingle.Spec.ExtraArgs, cr.Spec.Victoriametrics.VmSingle.ExtraArgs)
		}

		//A single-node VictoriaMetrics is capable of proxying requests to vmalert
		//https://docs.victoriametrics.com/Single-server-VictoriaMetrics.html#vmalert
		if cr.Spec.Victoriametrics.VmOperator.IsInstall() && cr.Spec.Victoriametrics.VmAlert.IsInstall() {
			vmAlert := vmetricsv1b1.VMAlert{}
			vmAlert.SetName(utils.VmComponentName)
			vmAlert.SetNamespace(cr.GetNamespace())
			if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.TLSEnabled {
				vmAlert.Spec.ExtraArgs = make(map[string]string)
				maps.Copy(vmAlert.Spec.ExtraArgs, map[string]string{"tls": "true"})
			}
			if cr.Spec.Victoriametrics.VmAlert.Port != "" {
				vmAlert.Spec.Port = cr.Spec.Victoriametrics.VmAlert.Port
			} else {
				vmAlert.Spec.Port = "8080"
			}
			maps.Copy(vmsingle.Spec.ExtraArgs, map[string]string{"vmalert.proxyURL": vmAlert.AsURL()})
		}

		if cr.Spec.Victoriametrics.VmAgent.Replicas != nil && *cr.Spec.Victoriametrics.VmAgent.Replicas > 1 {
			maps.Copy(vmsingle.Spec.ExtraArgs, map[string]string{"dedup.minScrapeInterval": "30s"})
		}

		if cr.Spec.Victoriametrics.VmSingle.ExtraEnvs != nil {
			vmsingle.Spec.ExtraEnvs = cr.Spec.Victoriametrics.VmSingle.ExtraEnvs
		}

		// Set tolerations for vmsingle
		if cr.Spec.Victoriametrics.VmSingle.Tolerations != nil {
			vmsingle.Spec.Tolerations = cr.Spec.Victoriametrics.VmSingle.Tolerations
		}

		if cr.Spec.Victoriametrics.VmSingle.TerminationGracePeriodSeconds != nil {
			vmsingle.Spec.TerminationGracePeriodSeconds = cr.Spec.Victoriametrics.VmSingle.TerminationGracePeriodSeconds
		}

		if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.TLSEnabled {
			vmsingle.Spec.Secrets = append(vmsingle.Spec.Secrets, victoriametrics.GetVmsingleTLSSecretName(cr.Spec.Victoriametrics.VmSingle))

			if vmsingle.Spec.ExtraArgs == nil {
				vmsingle.Spec.ExtraArgs = make(map[string]string)
			}
			maps.Copy(vmsingle.Spec.ExtraArgs, map[string]string{"tls": "true"})
			maps.Copy(vmsingle.Spec.ExtraArgs, map[string]string{"tlsCertFile": "/etc/vm/secrets/" + victoriametrics.GetVmsingleTLSSecretName(cr.Spec.Victoriametrics.VmSingle) + "/tls.crt"})
			maps.Copy(vmsingle.Spec.ExtraArgs, map[string]string{"tlsKeyFile": "/etc/vm/secrets/" + victoriametrics.GetVmsingleTLSSecretName(cr.Spec.Victoriametrics.VmSingle) + "/tls.key"})
		}

		// Set labels
		vmsingle.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(vmsingle.GetName(), vmsingle.GetNamespace())
		vmsingle.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Victoriametrics.VmSingle.Image)

		vmsingle.Spec.PodMetadata = &vmetricsv1b1.EmbeddedObjectMetadata{Labels: map[string]string{
			"name":                         utils.TruncLabel(vmsingle.GetName()),
			"app.kubernetes.io/name":       utils.TruncLabel(vmsingle.GetName()),
			"app.kubernetes.io/instance":   utils.GetInstanceLabel(vmsingle.GetName(), vmsingle.GetNamespace()),
			"app.kubernetes.io/component":  "victoriametrics",
			"app.kubernetes.io/part-of":    "monitoring",
			"app.kubernetes.io/managed-by": "monitoring-operator",
			"app.kubernetes.io/version":    utils.GetTagFromImage(cr.Spec.Victoriametrics.VmSingle.Image),
		}}

		if vmsingle.Spec.PodMetadata != nil {
			if cr.Spec.Victoriametrics.VmSingle.Labels != nil {
				for k, v := range cr.Spec.Victoriametrics.VmSingle.Labels {
					vmsingle.Spec.PodMetadata.Labels[k] = v
				}
			}

			if vmsingle.Spec.PodMetadata.Annotations == nil && cr.Spec.Victoriametrics.VmSingle.Annotations != nil {
				vmsingle.Spec.PodMetadata.Annotations = cr.Spec.Victoriametrics.VmSingle.Annotations
			} else {
				for k, v := range cr.Spec.Victoriametrics.VmSingle.Annotations {
					vmsingle.Spec.PodMetadata.Annotations[k] = v
				}
			}
		}

		if len(strings.TrimSpace(cr.Spec.Victoriametrics.VmSingle.PriorityClassName)) > 0 {
			vmsingle.Spec.PriorityClassName = cr.Spec.Victoriametrics.VmSingle.PriorityClassName
		}
	}

	return &vmsingle, nil
}

func vmSingleIngressV1(cr *v1alpha1.PlatformMonitoring) (*networkingv1.Ingress, error) {
	ingress := networkingv1.Ingress{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmSingleIngressAsset), 100).Decode(&ingress); err != nil {
		return nil, err
	}
	//Set metadata
	ingress.SetGroupVersionKind(schema.GroupVersionKind{Group: "networking.k8s.io", Version: "v1", Kind: "Ingress"})
	ingress.SetName(cr.GetNamespace() + "-" + utils.VmSingleServiceName)
	ingress.SetNamespace(cr.GetNamespace())

	if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.VmSingle.Ingress != nil && cr.Spec.Victoriametrics.VmSingle.Ingress.IsInstall() {
		var rules []networkingv1.IngressRule
		pathType := networkingv1.PathTypePrefix
		ing := cr.Spec.Victoriametrics.VmSingle.Ingress

		switch {
		// 1. If Host is provided
		case ing.Host != "":
			rules = append(rules, networkingv1.IngressRule{
				Host: ing.Host,
				IngressRuleValue: networkingv1.IngressRuleValue{
					HTTP: &networkingv1.HTTPIngressRuleValue{
						Paths: []networkingv1.HTTPIngressPath{defaultVmSinglePath(pathType)},
					},
				},
			})

		// 2. If custom ingress rules provided
		case len(ing.Rules) > 0:
			for _, r := range ing.Rules {
				// fallback if HTTP not set
				if r.HTTP == nil || len(r.HTTP.Paths) == 0 {
					r.HTTP = &v1alpha1.HTTPIngressRuleValue{
						Paths: []v1alpha1.IngressPath{
							{
								Path:     "/",
								PathType: string(pathType),
								Backend: v1alpha1.IngressPathBackend{
									Service: v1alpha1.IngressPathBackendService{
										Name: utils.VmSingleServiceName,
										Port: v1alpha1.ServiceBackendPort{
											Number: utils.VmSingleServicePort,
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
						Paths: []networkingv1.HTTPIngressPath{defaultVmSinglePath(pathType)},
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
		if !tlsConfigured && len(cr.Spec.Victoriametrics.VmSingle.Ingress.TLS) > 0 {
			for _, hostgroup := range cr.Spec.Victoriametrics.VmSingle.Ingress.TLS {
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
				// fallback: if secretName is empty - use TLSSecretName or secret from TLSConfig
				if secret == "" {
					secret = pickSecret(cr.Spec.Victoriametrics.VmSingle.Ingress.TLSSecretName, cr.Spec.Victoriametrics.VmAuth.TLSConfig)
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
		if !tlsConfigured && cr.Spec.Victoriametrics.VmSingle.Ingress.Host != "" {
			secret := pickSecret(cr.Spec.Victoriametrics.VmSingle.Ingress.TLSSecretName, cr.Spec.Victoriametrics.VmSingle.TLSConfig)
			if secret != "" {
				ingress.Spec.TLS = []networkingv1.IngressTLS{
					{
						Hosts:      []string{cr.Spec.Victoriametrics.VmSingle.Ingress.Host},
						SecretName: secret,
					},
				}
				tlsConfigured = true
			}
		}
		// Fallback: use ingress rules to configure tls hosts and TLSSecretName
		if !tlsConfigured && len(cr.Spec.Victoriametrics.VmSingle.Ingress.Rules) > 0 {
			tlsHosts := []string{}
			secret := pickSecret(cr.Spec.Victoriametrics.VmSingle.Ingress.TLSSecretName, cr.Spec.Victoriametrics.VmSingle.TLSConfig)
			for _, rule := range cr.Spec.Victoriametrics.VmSingle.Ingress.Rules {
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

		if cr.Spec.Victoriametrics.VmSingle.Ingress.IngressClassName != nil {
			ingress.Spec.IngressClassName = cr.Spec.Victoriametrics.VmSingle.Ingress.IngressClassName
		}

		// Set annotations
		ingress.SetAnnotations(cr.Spec.Victoriametrics.VmSingle.Ingress.Annotations)
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
		ingress.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Victoriametrics.VmSingle.Image)

		for lKey, lValue := range cr.Spec.Victoriametrics.VmSingle.Ingress.Labels {
			ingress.GetLabels()[lKey] = lValue
		}
	}
	return &ingress, nil
}

func defaultVmSinglePath(pathType networkingv1.PathType) networkingv1.HTTPIngressPath {
	return networkingv1.HTTPIngressPath{
		Path:     "/",
		PathType: &pathType,
		Backend: networkingv1.IngressBackend{
			Service: &networkingv1.IngressServiceBackend{
				Name: utils.VmSingleServiceName,
				Port: networkingv1.ServiceBackendPort{
					Number: utils.VmSingleServicePort,
				},
			},
		},
	}
}
