package vmauth

import (
	"embed"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"maps"

	v1alpha1 "github.com/Netcracker/qubership-monitoring-operator/api/v1alpha1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/victoriametrics"
	vmetricsv1b1 "github.com/VictoriaMetrics/operator/api/operator/v1beta1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/yaml"
)

//go:embed  assets/*.yaml
var assets embed.FS

func vmAuthServiceAccount(cr *v1alpha1.PlatformMonitoring) (*corev1.ServiceAccount, error) {
	sa := corev1.ServiceAccount{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmAuthServiceAccountAsset), 100).Decode(&sa); err != nil {
		return nil, err
	}
	//Set parameters
	sa.SetGroupVersionKind(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "ServiceAccount"})
	sa.SetName(cr.GetNamespace() + "-" + utils.VmAuthComponentName)
	sa.SetNamespace(cr.GetNamespace())

	return &sa, nil
}

func vmAuthClusterRole(cr *v1alpha1.PlatformMonitoring, hasPsp, hasScc bool) (*rbacv1.ClusterRole, error) {
	clusterRole := rbacv1.ClusterRole{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmAuthClusterRoleAsset), 100).Decode(&clusterRole); err != nil {
		return nil, err
	}
	//Set parameters
	clusterRole.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRole"})
	clusterRole.SetName(cr.GetNamespace() + "-" + utils.VmAuthComponentName)
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

func vmAuthClusterRoleBinding(cr *v1alpha1.PlatformMonitoring) (*rbacv1.ClusterRoleBinding, error) {
	clusterRoleBinding := rbacv1.ClusterRoleBinding{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmAuthClusterRoleBindingAsset), 100).Decode(&clusterRoleBinding); err != nil {
		return nil, err
	}
	//Set parameters
	clusterRoleBinding.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRoleBinding"})
	clusterRoleBinding.SetName(cr.GetNamespace() + "-" + utils.VmAuthComponentName)
	clusterRoleBinding.RoleRef.Name = cr.GetNamespace() + "-" + utils.VmAuthComponentName

	// Set namespace for all subjects
	for it := range clusterRoleBinding.Subjects {
		sub := &clusterRoleBinding.Subjects[it]
		sub.Namespace = cr.GetNamespace()
		sub.Name = cr.GetNamespace() + "-" + utils.VmAuthComponentName
	}
	return &clusterRoleBinding, nil
}

func vmAuth(r *VmAuthReconciler, cr *v1alpha1.PlatformMonitoring) (*vmetricsv1b1.VMAuth, error) {
	var err error
	vmauth := vmetricsv1b1.VMAuth{}
	if err = yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmAuthAsset), 100).Decode(&vmauth); err != nil {
		return nil, err
	}

	// Set parameters
	vmauth.SetNamespace(cr.GetNamespace())

	if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.VmAuth.IsInstall() {

		// Set VmAuth image
		vmauth.Spec.Image.Repository, vmauth.Spec.Image.Tag = utils.SplitImage(cr.Spec.Victoriametrics.VmAuth.Image)

		if r != nil {
			// Set security context
			if cr.Spec.Victoriametrics.VmAuth.SecurityContext != nil {
				if vmauth.Spec.SecurityContext == nil {
					vmauth.Spec.SecurityContext = &vmetricsv1b1.SecurityContext{}
				}
				if cr.Spec.Victoriametrics.VmAuth.SecurityContext.RunAsUser != nil {
					vmauth.Spec.SecurityContext.RunAsUser = cr.Spec.Victoriametrics.VmAuth.SecurityContext.RunAsUser
				}
				if cr.Spec.Victoriametrics.VmAuth.SecurityContext.RunAsGroup != nil {
					vmauth.Spec.SecurityContext.RunAsGroup = cr.Spec.Victoriametrics.VmAuth.SecurityContext.RunAsGroup
				}
				if cr.Spec.Victoriametrics.VmAuth.SecurityContext.FSGroup != nil {
					vmauth.Spec.SecurityContext.FSGroup = cr.Spec.Victoriametrics.VmAuth.SecurityContext.FSGroup
				}
			}
		}

		vmauth.Spec.ServiceAccountName = cr.GetNamespace() + "-" + utils.VmAuthComponentName

		// Set secrets for VmAuth deployment
		if len(cr.Spec.Victoriametrics.VmAuth.Secrets) > 0 {
			vmauth.Spec.Secrets = cr.Spec.Victoriametrics.VmAuth.Secrets
		}

		// Set additional containers
		if cr.Spec.Victoriametrics.VmAuth.Containers != nil {
			vmauth.Spec.Containers = cr.Spec.Victoriametrics.VmAuth.Containers
		}

		if len(cr.Spec.Victoriametrics.VmAuth.ConfigMaps) > 0 {
			vmauth.Spec.ConfigMaps = cr.Spec.Victoriametrics.VmAuth.ConfigMaps
		}

		// Set resources for VmAuth deployment
		if cr.Spec.Victoriametrics.VmAuth.Resources.Size() > 0 {
			vmauth.Spec.Resources = cr.Spec.Victoriametrics.VmAuth.Resources
		}

		if cr.Spec.Victoriametrics.VmAuth.ReplicaCount != nil {
			vmauth.Spec.ReplicaCount = cr.Spec.Victoriametrics.VmAuth.ReplicaCount
		}
		if cr.Spec.Victoriametrics.VmCluster.IsInstall() && cr.Spec.Victoriametrics.VmReplicas != nil {
			vmauth.Spec.ReplicaCount = cr.Spec.Victoriametrics.VmReplicas
		}

		// Set additional volumes
		if cr.Spec.Victoriametrics.VmAuth.Volumes != nil {
			vmauth.Spec.Volumes = cr.Spec.Victoriametrics.VmAuth.Volumes
		}

		// Set additional volumeMounts for vmauth container.
		if cr.Spec.Victoriametrics.VmAuth.VolumeMounts != nil {
			for it := range vmauth.Spec.Containers {
				c := &vmauth.Spec.Containers[it]

				// Set additional volumeMounts only for VmAuth container
				if c.Name == utils.VmAuthComponentName {
					c.VolumeMounts = cr.Spec.Victoriametrics.VmAuth.VolumeMounts
				}
			}
		}

		// Set affinity for VmAuth
		if cr.Spec.Victoriametrics.VmAuth.Affinity != nil {
			vmauth.Spec.Affinity = cr.Spec.Victoriametrics.VmAuth.Affinity
		}

		// Set tolerations for VmAuth
		if cr.Spec.Victoriametrics.VmAuth.Tolerations != nil {
			vmauth.Spec.Tolerations = cr.Spec.Victoriametrics.VmAuth.Tolerations
		}

		if len(cr.Spec.Victoriametrics.VmAuth.Port) > 0 {
			vmauth.Spec.Port = cr.Spec.Victoriametrics.VmAuth.Port
		}

		if cr.Spec.Victoriametrics.VmAuth.SelectAllByDefault {
			vmauth.Spec.SelectAllByDefault = true
		}

		if cr.Spec.Victoriametrics.VmAuth.UserSelector != nil {
			vmauth.Spec.UserSelector = cr.Spec.Victoriametrics.VmAuth.UserSelector
		}

		if cr.Spec.Victoriametrics.VmAuth.UserNamespaceSelector != nil {
			vmauth.Spec.UserNamespaceSelector = cr.Spec.Victoriametrics.VmAuth.UserNamespaceSelector
		}

		if cr.Spec.Victoriametrics.VmAuth.ExtraArgs != nil {
			vmauth.Spec.ExtraArgs = cr.Spec.Victoriametrics.VmAuth.ExtraArgs
		}

		if cr.Spec.Victoriametrics.VmAuth.ExtraEnvs != nil {
			vmauth.Spec.ExtraEnvs = cr.Spec.Victoriametrics.VmAuth.ExtraEnvs
		}

		if cr.Spec.Victoriametrics.VmAuth.Tolerations != nil {
			vmauth.Spec.Tolerations = cr.Spec.Victoriametrics.VmAuth.Tolerations
		}

		// Set Auth
		if cr.Spec.Auth != nil && cr.Spec.OAuthProxy != nil {
			externalURL := utils.ExternalURLPrefix
			if cr.Spec.Victoriametrics.VmAuth.Ingress != nil &&
				cr.Spec.Victoriametrics.VmAuth.Ingress.IsInstall() &&
				cr.Spec.Victoriametrics.VmAuth.Ingress.Host != "" {
				externalURL += cr.Spec.Victoriametrics.VmAuth.Ingress.Host
			}
			if externalURL == utils.ExternalURLPrefix {
				return nil, errors.New("host for ingress can not be empty")
			}

			vmauth.Spec.Secrets = append(vmauth.Spec.Secrets, utils.OAuthProxySecret)
			// Volume mounts for oauth2-proxy sidecar
			var vms []corev1.VolumeMount

			// Add oauth2-proxy config
			vms = append(vms, corev1.VolumeMount{MountPath: utils.OAuthProxySecretDir, Name: utils.OAuthProxySecretName})

			if cr.Spec.Auth.TLSConfig != nil {
				// Add CA secret
				if cr.Spec.Auth.TLSConfig.CASecret != nil {
					vmauth.Spec.Secrets = append(vmauth.Spec.Secrets, cr.Spec.Auth.TLSConfig.CASecret.Name)
					vms = append(vms, corev1.VolumeMount{
						MountPath: utils.TlsCertificatesSecretDir + "/" + cr.Spec.Auth.TLSConfig.CASecret.Name,
						Name:      utils.SecretNamePrefix + cr.Spec.Auth.TLSConfig.CASecret.Name,
					})
				}
				// Add Cert secret
				if cr.Spec.Auth.TLSConfig.CertSecret != nil {
					vmauth.Spec.Secrets = append(vmauth.Spec.Secrets, cr.Spec.Auth.TLSConfig.CertSecret.Name)
					vms = append(vms, corev1.VolumeMount{
						MountPath: utils.TlsCertificatesSecretDir + "/" + cr.Spec.Auth.TLSConfig.CertSecret.Name,
						Name:      utils.SecretNamePrefix + cr.Spec.Auth.TLSConfig.CertSecret.Name,
					})
				}
				// Add Key secret
				if cr.Spec.Auth.TLSConfig.KeySecret != nil {
					vmauth.Spec.Secrets = append(vmauth.Spec.Secrets, cr.Spec.Auth.TLSConfig.KeySecret.Name)
					vms = append(vms, corev1.VolumeMount{
						MountPath: utils.TlsCertificatesSecretDir + "/" + cr.Spec.Auth.TLSConfig.KeySecret.Name,
						Name:      utils.SecretNamePrefix + cr.Spec.Auth.TLSConfig.KeySecret.Name,
					})
				}
			}

			port := cr.Spec.Victoriametrics.VmAuth.Port
			if port == "" {
				port = strconv.Itoa(utils.VmAuthServicePort)
			}

			var upstream string
			if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.TLSEnabled {
				upstream = fmt.Sprintf("https://%s.%s.svc:%s", vmauth.PrefixedName(), cr.Namespace, port)
			} else {
				upstream = fmt.Sprintf("http://%s.%s.svc:%s", vmauth.PrefixedName(), cr.Namespace, port)
			}

			// Configure oauthProxy for support authentication
			sidecar := corev1.Container{
				Name:            utils.OAuthProxyName,
				Image:           cr.Spec.OAuthProxy.Image,
				ImagePullPolicy: corev1.PullIfNotPresent,
				Ports:           []corev1.ContainerPort{{Name: utils.OAuthProxyName, ContainerPort: utils.OAuthPort, Protocol: corev1.ProtocolTCP}},
				VolumeMounts:    vms,
				Args: []string{
					"--redirect-url=" + externalURL,
					"--upstream=" + upstream,
					"--config=" + utils.OAuthProxyCfg,
				},
			}

			containerIndex := -1
			for idx, c := range vmauth.Spec.Containers {
				if c.Name == utils.OAuthProxyName {
					containerIndex = idx
					break
				}
			}
			if containerIndex > 0 {
				vmauth.Spec.Containers[containerIndex] = sidecar
			} else {
				vmauth.Spec.Containers = append(vmauth.Spec.Containers, sidecar)
			}

			svc := &vmetricsv1b1.AdditionalServiceSpec{
				EmbeddedObjectMetadata: vmetricsv1b1.EmbeddedObjectMetadata{
					Name:        utils.VmAuthOAuthProxyServiceName,
					Labels:      vmauth.AllLabels(),
					Annotations: vmauth.AnnotationsFiltered(),
				},
				Spec: corev1.ServiceSpec{
					Type:     corev1.ServiceTypeClusterIP,
					Selector: vmauth.SelectorLabels(),
					Ports: []corev1.ServicePort{
						{
							Name:       utils.OAuthProxyName,
							Protocol:   corev1.ProtocolTCP,
							Port:       int32(utils.VmAuthServicePort),
							TargetPort: intstr.Parse(utils.OAuthProxyName),
						},
					},
				},
			}
			vmauth.Spec.ServiceSpec = svc
		}

		if cr.Spec.Victoriametrics.VmAuth.TerminationGracePeriodSeconds != nil {
			vmauth.Spec.TerminationGracePeriodSeconds = cr.Spec.Victoriametrics.VmAuth.TerminationGracePeriodSeconds
		}

		if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.TLSEnabled {
			vmauth.Spec.Secrets = append(vmauth.Spec.Secrets, victoriametrics.GetVmauthTLSSecretName(cr.Spec.Victoriametrics.VmAuth))

			if vmauth.Spec.ExtraArgs == nil {
				vmauth.Spec.ExtraArgs = make(map[string]string)
			}
			maps.Copy(vmauth.Spec.ExtraArgs, map[string]string{"tls": "true"})
			maps.Copy(vmauth.Spec.ExtraArgs, map[string]string{"tlsCertFile": "/etc/vm/secrets/" + victoriametrics.GetVmauthTLSSecretName(cr.Spec.Victoriametrics.VmAuth) + "/tls.crt"})
			maps.Copy(vmauth.Spec.ExtraArgs, map[string]string{"tlsKeyFile": "/etc/vm/secrets/" + victoriametrics.GetVmauthTLSSecretName(cr.Spec.Victoriametrics.VmAuth) + "/tls.key"})
		}

		// Set labels
		vmauth.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(vmauth.GetName(), vmauth.GetNamespace())
		vmauth.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Victoriametrics.VmAuth.Image)

		vmauth.Spec.PodMetadata = &vmetricsv1b1.EmbeddedObjectMetadata{Labels: map[string]string{
			"name":                         utils.TruncLabel(vmauth.GetName()),
			"app.kubernetes.io/name":       utils.TruncLabel(vmauth.GetName()),
			"app.kubernetes.io/instance":   utils.GetInstanceLabel(vmauth.GetName(), vmauth.GetNamespace()),
			"app.kubernetes.io/component":  "victoriametrics",
			"app.kubernetes.io/part-of":    "monitoring",
			"app.kubernetes.io/managed-by": "monitoring-operator",
			"app.kubernetes.io/version":    utils.GetTagFromImage(cr.Spec.Victoriametrics.VmAuth.Image),
		}}

		if vmauth.Spec.PodMetadata != nil {
			if cr.Spec.Victoriametrics.VmAuth.Labels != nil {
				for k, v := range cr.Spec.Victoriametrics.VmAuth.Labels {
					vmauth.Spec.PodMetadata.Labels[k] = v
				}
			}

			if vmauth.Spec.PodMetadata.Annotations == nil && cr.Spec.Victoriametrics.VmAuth.Annotations != nil {
				vmauth.Spec.PodMetadata.Annotations = cr.Spec.Victoriametrics.VmAuth.Annotations
			} else {
				for k, v := range cr.Spec.Victoriametrics.VmAuth.Annotations {
					vmauth.Spec.PodMetadata.Annotations[k] = v
				}
			}
		}

		if len(strings.TrimSpace(cr.Spec.Victoriametrics.VmAuth.PriorityClassName)) > 0 {
			vmauth.Spec.PriorityClassName = cr.Spec.Victoriametrics.VmAuth.PriorityClassName
		}
	}
	return &vmauth, nil
}

func vmAuthIngress(cr *v1alpha1.PlatformMonitoring) (*networkingv1.Ingress, error) {
	ingress := networkingv1.Ingress{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmAuthIngressAsset), 100).Decode(&ingress); err != nil {
		return nil, err
	}
	//Set parameters
	ingress.SetGroupVersionKind(networkingv1.SchemeGroupVersion.WithKind("Ingress"))
	ingress.SetName(cr.GetNamespace() + "-" + utils.VmAuthServiceName)
	ingress.SetNamespace(cr.GetNamespace())

	if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.VmAuth.Ingress != nil && cr.Spec.Victoriametrics.VmAuth.Ingress.IsInstall() {
		var rules []networkingv1.IngressRule
		pathType := networkingv1.PathTypePrefix
		ing := cr.Spec.Victoriametrics.VmAuth.Ingress

		switch {
		// 1. If Host is provided
		case ing.Host != "":
			rules = append(rules, networkingv1.IngressRule{
				Host: ing.Host,
				IngressRuleValue: networkingv1.IngressRuleValue{
					HTTP: &networkingv1.HTTPIngressRuleValue{
						Paths: []networkingv1.HTTPIngressPath{defaultVmAuthPath(pathType)},
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
										Name: utils.VmAuthServiceName,
										Port: v1alpha1.ServiceBackendPort{
											Number: utils.VmAuthServicePort,
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
						Paths: []networkingv1.HTTPIngressPath{defaultVmAuthPath(pathType)},
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
		if !tlsConfigured && len(cr.Spec.Victoriametrics.VmAuth.Ingress.TLS) > 0 {
			for _, hostgroup := range cr.Spec.Victoriametrics.VmAuth.Ingress.TLS {
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
					secret = pickSecret(cr.Spec.Victoriametrics.VmAuth.Ingress.TLSSecretName, cr.Spec.Victoriametrics.VmAuth.TLSConfig)
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
		if !tlsConfigured && cr.Spec.Victoriametrics.VmAuth.Ingress.Host != "" {
			secret := pickSecret(cr.Spec.Victoriametrics.VmAuth.Ingress.TLSSecretName, cr.Spec.Victoriametrics.VmAuth.TLSConfig)
			if secret != "" {
				ingress.Spec.TLS = []networkingv1.IngressTLS{
					{
						Hosts:      []string{cr.Spec.Victoriametrics.VmAuth.Ingress.Host},
						SecretName: secret,
					},
				}
				tlsConfigured = true
			}
		}
		// Fallback: use ingress rules to configure tls hosts and TLSSecretName
		if !tlsConfigured && len(cr.Spec.Victoriametrics.VmAuth.Ingress.Rules) > 0 {
			tlsHosts := []string{}
			secret := pickSecret(cr.Spec.Victoriametrics.VmAuth.Ingress.TLSSecretName, cr.Spec.Victoriametrics.VmAuth.TLSConfig)
			for _, rule := range cr.Spec.Victoriametrics.VmAuth.Ingress.Rules {
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

		if cr.Spec.Victoriametrics.VmAuth.Ingress.IngressClassName != nil {
			ingress.Spec.IngressClassName = cr.Spec.Victoriametrics.VmAuth.Ingress.IngressClassName
		}

		// Set annotations
		ingress.SetAnnotations(cr.Spec.Victoriametrics.VmAuth.Ingress.Annotations)
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
		ingress.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Victoriametrics.VmAuth.Image)

		ingress.SetLabels(labels.Merge(ingress.GetLabels(), cr.Spec.Victoriametrics.VmAuth.Ingress.Labels))
	}
	return &ingress, nil
}

func defaultVmAuthPath(pathType networkingv1.PathType) networkingv1.HTTPIngressPath {
	return networkingv1.HTTPIngressPath{
		Path:     "/",
		PathType: &pathType,
		Backend: networkingv1.IngressBackend{
			Service: &networkingv1.IngressServiceBackend{
				Name: utils.VmAuthServiceName,
				Port: networkingv1.ServiceBackendPort{
					Number: utils.VmAuthServicePort,
				},
			},
		},
	}
}
