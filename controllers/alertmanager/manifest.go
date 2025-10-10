package alertmanager

import (
	"embed"
	"strings"

	v1alpha1 "github.com/Netcracker/qubership-monitoring-operator/api/v1alpha1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/yaml"
)

//go:embed  assets/*.yaml
var assets embed.FS

func alertmanagerServiceAccount(cr *v1alpha1.PlatformMonitoring) (*corev1.ServiceAccount, error) {
	sa := corev1.ServiceAccount{}
	err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.AlertManagerServiceAccountAsset), 100).Decode(&sa)

	if err != nil {
		return nil, err
	}
	//Set parameters
	sa.SetGroupVersionKind(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "ServiceAccount"})
	sa.SetName(cr.GetNamespace() + "-" + utils.AlertManagerComponentName)
	sa.SetNamespace(cr.GetNamespace())

	// Set annotations and labels for ServiceAccount in case
	if cr.Spec.AlertManager != nil && cr.Spec.AlertManager.ServiceAccount != nil {
		if sa.Annotations == nil && cr.Spec.AlertManager.ServiceAccount.Annotations != nil {
			sa.SetAnnotations(cr.Spec.AlertManager.ServiceAccount.Annotations)
		} else {
			for k, v := range cr.Spec.AlertManager.ServiceAccount.Annotations {
				sa.Annotations[k] = v
			}
		}

		if sa.Labels == nil && cr.Spec.AlertManager.ServiceAccount.Labels != nil {
			sa.SetLabels(cr.Spec.AlertManager.ServiceAccount.Labels)
		} else {
			for k, v := range cr.Spec.AlertManager.ServiceAccount.Labels {
				sa.Labels[k] = v
			}
		}
	}

	return &sa, nil
}

func alertmanagerSecret(cr *v1alpha1.PlatformMonitoring) (*corev1.Secret, error) {
	secret := corev1.Secret{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.AlertManagerSecretAsset), 100).Decode(&secret); err != nil {
		return nil, err
	}
	//Set parameters
	secret.SetGroupVersionKind(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Secret"})
	secret.SetNamespace(cr.GetNamespace())

	return &secret, nil
}

func alertmanager(cr *v1alpha1.PlatformMonitoring) (*promv1.Alertmanager, error) {
	am := promv1.Alertmanager{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.AlertManagerAsset), 100).Decode(&am); err != nil {
		return nil, err
	}
	//Set parameters
	am.SetGroupVersionKind(schema.GroupVersionKind{Group: "monitoring.coreos.com", Version: "v1", Kind: "Alertmanager"})
	am.SetNamespace(cr.GetNamespace())
	am.Spec.ServiceAccountName = cr.GetNamespace() + "-" + utils.AlertManagerComponentName

	// Set AlertManager image
	if cr.Spec.AlertManager != nil {
		am.Spec.Image = &cr.Spec.AlertManager.Image

		// Set labels
		am.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.AlertManager.Image)
		am.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(am.GetName(), am.GetNamespace())

		// Set Alertmanager replicas
		if cr.Spec.AlertManager.Replicas != nil {
			am.Spec.Replicas = cr.Spec.AlertManager.Replicas
		}
		// Set security context
		if cr.Spec.AlertManager.SecurityContext != nil {
			if am.Spec.SecurityContext == nil {
				am.Spec.SecurityContext = &corev1.PodSecurityContext{}
			}
			if cr.Spec.AlertManager.SecurityContext.RunAsUser != nil {
				am.Spec.SecurityContext.RunAsUser = cr.Spec.AlertManager.SecurityContext.RunAsUser
			}
			if cr.Spec.AlertManager.SecurityContext.FSGroup != nil {
				am.Spec.SecurityContext.FSGroup = cr.Spec.AlertManager.SecurityContext.FSGroup
			}
		}
		// Set resources for AlertManager deployment
		if cr.Spec.AlertManager.Resources.Size() > 0 {
			am.Spec.Resources = cr.Spec.AlertManager.Resources
		}
		// Set additional containers
		if cr.Spec.AlertManager.Containers != nil {
			am.Spec.Containers = cr.Spec.AlertManager.Containers
		}
		// Set Auth
		if cr.Spec.Auth != nil && cr.Spec.OAuthProxy != nil {
			am.Spec.Secrets = append(am.Spec.Secrets, "oauth2-proxy-config")

			externalURL := "http://"
			if cr.Spec.AlertManager.Ingress != nil &&
				cr.Spec.AlertManager.Ingress.IsInstall() &&
				cr.Spec.AlertManager.Ingress.Host != "" {
				externalURL += cr.Spec.AlertManager.Ingress.Host
			}
			// Volume mounts for oauth2-proxy sidecar
			var vms []corev1.VolumeMount

			// Add oauth2-proxy config
			vms = append(vms, corev1.VolumeMount{MountPath: utils.OAuthProxySecretDir, Name: "secret-oauth2-proxy-config"})

			if cr.Spec.Auth.TLSConfig != nil {
				// Add CA secret
				if cr.Spec.Auth.TLSConfig.CASecret != nil {
					am.Spec.Secrets = append(am.Spec.Secrets, cr.Spec.Auth.TLSConfig.CASecret.Name)
					vms = append(vms, corev1.VolumeMount{
						MountPath: utils.TlsCertificatesSecretDir + "/" + cr.Spec.Auth.TLSConfig.CASecret.Name,
						Name:      "secret-" + cr.Spec.Auth.TLSConfig.CASecret.Name,
					})
				}
				// Add Cert secret
				if cr.Spec.Auth.TLSConfig.CertSecret != nil {
					am.Spec.Secrets = append(am.Spec.Secrets, cr.Spec.Auth.TLSConfig.CertSecret.Name)
					vms = append(vms, corev1.VolumeMount{
						MountPath: utils.TlsCertificatesSecretDir + "/" + cr.Spec.Auth.TLSConfig.CertSecret.Name,
						Name:      "secret-" + cr.Spec.Auth.TLSConfig.CertSecret.Name,
					})
				}
				// Add Key secret
				if cr.Spec.Auth.TLSConfig.KeySecret != nil {
					am.Spec.Secrets = append(am.Spec.Secrets, cr.Spec.Auth.TLSConfig.KeySecret.Name)
					vms = append(vms, corev1.VolumeMount{
						MountPath: utils.TlsCertificatesSecretDir + "/" + cr.Spec.Auth.TLSConfig.KeySecret.Name,
						Name:      "secret-" + cr.Spec.Auth.TLSConfig.KeySecret.Name,
					})
				}
			}
			// Configure oauthProxy for support authentication
			sidecar := corev1.Container{
				Name:            utils.OAuthProxyName,
				Image:           cr.Spec.OAuthProxy.Image,
				ImagePullPolicy: "IfNotPresent",
				Ports:           []corev1.ContainerPort{{Name: utils.OAuthProxyName, ContainerPort: 9092, Protocol: "TCP"}},
				VolumeMounts:    vms,
				Args: []string{
					"--redirect-url=" + externalURL,
					"--upstream=http://localhost:9093",
					"--config=/etc/oauth-proxy/oauth2-proxy.cfg",
				},
			}

			containerIndex := -1
			for idx, c := range am.Spec.Containers {
				if c.Name == utils.OAuthProxyName {
					containerIndex = idx
					break
				}
			}
			if containerIndex > 0 {
				am.Spec.Containers[containerIndex] = sidecar
			} else {
				am.Spec.Containers = append(am.Spec.Containers, sidecar)
			}
		}
		// Set tolerations for AlertManager
		if cr.Spec.AlertManager.Tolerations != nil {
			am.Spec.Tolerations = cr.Spec.AlertManager.Tolerations
		}
		// Set nodeSelector for AlertManager
		if cr.Spec.AlertManager.NodeSelector != nil {
			am.Spec.NodeSelector = cr.Spec.AlertManager.NodeSelector
		}
		// Set affinity for AlertManager
		if cr.Spec.AlertManager.Affinity != nil {
			am.Spec.Affinity = cr.Spec.AlertManager.Affinity
		}

		// Set PodMetadata.Labels
		am.Spec.PodMetadata = &promv1.EmbeddedObjectMetadata{Labels: map[string]string{
			"name":                         "alertmanager",
			"app.kubernetes.io/name":       "alertmanager",
			"app.kubernetes.io/instance":   utils.GetInstanceLabel("alertmanager", am.GetNamespace()),
			"app.kubernetes.io/component":  "alertmanager",
			"app.kubernetes.io/part-of":    "monitoring",
			"app.kubernetes.io/version":    utils.GetTagFromImage(cr.Spec.AlertManager.Image),
			"app.kubernetes.io/managed-by": "monitoring-operator",
		}}

		if cr.Spec.AlertManager.Labels != nil {
			for k, v := range cr.Spec.AlertManager.Labels {
				am.Spec.PodMetadata.Labels[k] = v
			}
		}

		if am.Spec.PodMetadata.Annotations == nil && cr.Spec.AlertManager.Annotations != nil {
			am.Spec.PodMetadata.Annotations = cr.Spec.AlertManager.Annotations
		} else {
			for k, v := range cr.Spec.AlertManager.Annotations {
				am.Spec.PodMetadata.Annotations[k] = v
			}
		}

		if len(strings.TrimSpace(cr.Spec.AlertManager.PriorityClassName)) > 0 {
			am.Spec.PriorityClassName = cr.Spec.AlertManager.PriorityClassName
		}
	}
	return &am, nil
}

func alertmanagerService(cr *v1alpha1.PlatformMonitoring) (*corev1.Service, error) {
	service := corev1.Service{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.AlertManagerServiceAsset), 100).Decode(&service); err != nil {
		return nil, err
	}
	//Set parameters
	service.SetGroupVersionKind(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Service"})
	service.SetName(utils.AlertManagerComponentName)
	service.SetNamespace(cr.GetNamespace())

	// Set port
	if cr.Spec.AlertManager != nil {
		for p := range service.Spec.Ports {
			port := &service.Spec.Ports[p]
			if port.Name == "web" {
				port.NodePort = cr.Spec.AlertManager.Port
			}
		}
		if cr.Spec.Auth != nil && cr.Spec.OAuthProxy != nil {
			port := corev1.ServicePort{
				Name:       utils.OAuthProxyName,
				TargetPort: intstr.FromString(utils.OAuthProxyName),
				Port:       9092,
				Protocol:   corev1.ProtocolTCP,
			}
			service.Spec.Ports = append(service.Spec.Ports, port)
		}
	}
	return &service, nil
}

func alertmanagerIngressV1(cr *v1alpha1.PlatformMonitoring) (*networkingv1.Ingress, error) {
	ingress := networkingv1.Ingress{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.AlertManagerIngressAsset), 100).Decode(&ingress); err != nil {
		return nil, err
	}
	//Set metadata
	ingress.SetGroupVersionKind(schema.GroupVersionKind{Group: "networking.k8s.io", Version: "v1", Kind: "Ingress"})
	ingress.SetName(cr.GetNamespace() + "-" + utils.AlertManagerComponentName)
	ingress.SetNamespace(cr.GetNamespace())

	if cr.Spec.AlertManager != nil && cr.Spec.AlertManager.Ingress != nil && cr.Spec.AlertManager.Ingress.IsInstall() {
		var rules []networkingv1.IngressRule
		pathType := networkingv1.PathTypePrefix
		ing := cr.Spec.AlertManager.Ingress

		switch {
		// 1. If Host is provided
		case ing.Host != "":
			rules = append(rules, networkingv1.IngressRule{
				Host: ing.Host,
				IngressRuleValue: networkingv1.IngressRuleValue{
					HTTP: &networkingv1.HTTPIngressRuleValue{
						Paths: []networkingv1.HTTPIngressPath{defaultAlertManagerPath(pathType)},
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
										Name: utils.AlertManagerComponentName,
										Port: v1alpha1.ServiceBackendPort{
											Number: utils.AlertmanagerServicePort,
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
						Paths: []networkingv1.HTTPIngressPath{defaultAlertManagerPath(pathType)},
					},
				},
			})
		}
		ingress.Spec.Rules = rules

		tlsConfigured := false
		pickSecret := func(ingressTLSSecret string, tlsCfg *v1alpha1.CommonTLSConfig) string {
			if ingressTLSSecret != "" {
				return ingressTLSSecret
			}
			if tlsCfg != nil {
				return tlsCfg.SecretName
			}
			return ""
		}
		// Configure tls if TLS config is defined
		if !tlsConfigured && len(cr.Spec.AlertManager.Ingress.TLS) > 0 {
			for _, hostgroup := range cr.Spec.AlertManager.Ingress.TLS {
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
					secret = pickSecret(cr.Spec.AlertManager.Ingress.TLSSecretName, cr.Spec.AlertManager.TLSConfig)
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
		if !tlsConfigured && cr.Spec.AlertManager.Ingress.Host != "" {
			secret := pickSecret(cr.Spec.AlertManager.Ingress.TLSSecretName, cr.Spec.AlertManager.TLSConfig)
			if secret != "" {
				ingress.Spec.TLS = []networkingv1.IngressTLS{
					{
						Hosts:      []string{cr.Spec.AlertManager.Ingress.Host},
						SecretName: secret,
					},
				}
				tlsConfigured = true
			}
		}
		// Fallback: use ingress rules to configure tls hosts and TLSSecretName
		if !tlsConfigured && len(cr.Spec.AlertManager.Ingress.Rules) > 0 {
			tlsHosts := []string{}
			secret := pickSecret(cr.Spec.AlertManager.Ingress.TLSSecretName, cr.Spec.AlertManager.TLSConfig)
			for _, rule := range cr.Spec.AlertManager.Ingress.Rules {
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

		if cr.Spec.AlertManager.Ingress.IngressClassName != nil {
			ingress.Spec.IngressClassName = cr.Spec.AlertManager.Ingress.IngressClassName
		}

		// Set annotations
		ingress.SetAnnotations(cr.Spec.AlertManager.Ingress.Annotations)

		// Set labels with saving default labels
		ingress.Labels["name"] = utils.TruncLabel(ingress.GetName())
		ingress.Labels["app.kubernetes.io/name"] = utils.TruncLabel(ingress.GetName())
		ingress.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(ingress.GetName(), ingress.GetNamespace())
		ingress.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.AlertManager.Image)

		for lKey, lValue := range cr.Spec.AlertManager.Ingress.Labels {
			ingress.GetLabels()[lKey] = lValue
		}
	}
	return &ingress, nil
}

func alertmanagerPodMonitor(cr *v1alpha1.PlatformMonitoring) (*promv1.PodMonitor, error) {
	podMonitor := promv1.PodMonitor{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.AlertManagerPodMonitorAsset), 100).Decode(&podMonitor); err != nil {
		return nil, err
	}
	//Set parameters
	podMonitor.SetGroupVersionKind(schema.GroupVersionKind{Group: "monitoring.coreos.com", Version: "v1", Kind: "PodMonitor"})
	podMonitor.SetName(cr.GetNamespace() + "-" + "alertmanager-pod-monitor")
	podMonitor.SetNamespace(cr.GetNamespace())

	if cr.Spec.AlertManager != nil && cr.Spec.AlertManager.PodMonitor != nil && cr.Spec.AlertManager.PodMonitor.IsInstall() {
		cr.Spec.AlertManager.PodMonitor.OverridePodMonitor(&podMonitor)
	}

	return &podMonitor, nil
}

func defaultAlertManagerPath(pathType networkingv1.PathType) networkingv1.HTTPIngressPath {
	return networkingv1.HTTPIngressPath{
		Path:     "/",
		PathType: &pathType,
		Backend: networkingv1.IngressBackend{
			Service: &networkingv1.IngressServiceBackend{
				Name: utils.AlertManagerComponentName,
				Port: networkingv1.ServiceBackendPort{
					Number: utils.AlertmanagerServicePort,
				},
			},
		},
	}
}
