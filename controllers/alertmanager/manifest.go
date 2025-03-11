package alertmanager

import (
	"embed"
	"errors"
	"strings"

	v1alpha1 "github.com/Netcracker/qubership-monitoring-operator/api/v1alpha1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/api/networking/v1beta1"
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

func alertmanagerIngressV1beta1(cr *v1alpha1.PlatformMonitoring) (*v1beta1.Ingress, error) {
	ingress := v1beta1.Ingress{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.AlertManagerIngressAsset), 100).Decode(&ingress); err != nil {
		return nil, err
	}
	//Set parameters
	ingress.SetGroupVersionKind(schema.GroupVersionKind{Group: "networking.k8s.io", Version: "v1beta1", Kind: "Ingress"})
	ingress.SetName(cr.GetNamespace() + "-" + utils.AlertManagerComponentName)
	ingress.SetNamespace(cr.GetNamespace())

	if cr.Spec.AlertManager != nil && cr.Spec.AlertManager.Ingress != nil && cr.Spec.AlertManager.Ingress.IsInstall() {
		// Check that ingress host is specified.
		if cr.Spec.AlertManager.Ingress.Host == "" {
			return nil, errors.New("host for ingress can not be empty")
		}

		servicePort := intstr.FromInt(utils.AlertmanagerServicePort)
		serviceName := utils.AlertmanagerServiceName

		if cr.Spec.Auth != nil && cr.Spec.OAuthProxy != nil {
			servicePort = intstr.FromString(utils.OAuthProxyServicePortName)
			serviceName = utils.AlertmanagerOAuthProxyServiceName
		}
		// Add rule for alertmanager UI
		rule := v1beta1.IngressRule{Host: cr.Spec.AlertManager.Ingress.Host}
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
		if cr.Spec.AlertManager.Ingress.TLSSecretName != "" {
			ingress.Spec.TLS = []v1beta1.IngressTLS{
				{
					Hosts:      []string{cr.Spec.AlertManager.Ingress.Host},
					SecretName: cr.Spec.AlertManager.Ingress.TLSSecretName,
				},
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

func alertmanagerIngressV1(cr *v1alpha1.PlatformMonitoring) (*networkingv1.Ingress, error) {
	ingress := networkingv1.Ingress{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.AlertManagerIngressAsset), 100).Decode(&ingress); err != nil {
		return nil, err
	}
	//Set parameters
	ingress.SetGroupVersionKind(schema.GroupVersionKind{Group: "networking.k8s.io", Version: "v1", Kind: "Ingress"})
	ingress.SetName(cr.GetNamespace() + "-" + utils.AlertManagerComponentName)
	ingress.SetNamespace(cr.GetNamespace())

	if cr.Spec.AlertManager != nil && cr.Spec.AlertManager.Ingress != nil && cr.Spec.AlertManager.Ingress.IsInstall() {
		// Check that ingress host is specified.
		if cr.Spec.AlertManager.Ingress.Host == "" {
			return nil, errors.New("host for ingress can not be empty")
		}

		var ingressServiceBackend *networkingv1.IngressServiceBackend
		if cr.Spec.Auth != nil && cr.Spec.OAuthProxy != nil {
			ingressServiceBackend = &networkingv1.IngressServiceBackend{
				Name: utils.AlertmanagerOAuthProxyServiceName,
				Port: networkingv1.ServiceBackendPort{
					Name: utils.OAuthProxyServicePortName,
				},
			}
		} else {
			ingressServiceBackend = &networkingv1.IngressServiceBackend{
				Name: utils.AlertmanagerServiceName,
				Port: networkingv1.ServiceBackendPort{
					Number: utils.AlertmanagerServicePort,
				},
			}
		}
		pathType := networkingv1.PathTypePrefix
		// Add rule for alertmanager UI
		rule := networkingv1.IngressRule{Host: cr.Spec.AlertManager.Ingress.Host}
		rule.HTTP = &networkingv1.HTTPIngressRuleValue{
			Paths: []networkingv1.HTTPIngressPath{
				{
					Path:     "/",
					PathType: &pathType,
					Backend: networkingv1.IngressBackend{
						Service: ingressServiceBackend,
					},
				},
			},
		}
		ingress.Spec.Rules = []networkingv1.IngressRule{rule}

		// Configure TLS if TLS secret name is set
		if cr.Spec.AlertManager.Ingress.TLSSecretName != "" {
			ingress.Spec.TLS = []networkingv1.IngressTLS{
				{
					Hosts:      []string{cr.Spec.AlertManager.Ingress.Host},
					SecretName: cr.Spec.AlertManager.Ingress.TLSSecretName,
				},
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
