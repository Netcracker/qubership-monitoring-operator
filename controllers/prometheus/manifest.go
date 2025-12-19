package prometheus

import (
	"crypto/tls"
	"embed"
	"errors"
	"fmt"
	"strings"

	v1beta1 "github.com/Netcracker/qubership-monitoring-operator/api"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/utils/ptr"
)

//go:embed  assets/*.yaml
var assets embed.FS

// getPrometheusexternalURL add to host name http schema
// TODO: Currently hardcode unsecured http schema. If we want to keep the following logic in future
// need to get security settings from Ingress and according by settings select http or https.
// But currently we doesn't allow create secure Ingress.
func getPrometheusexternalURL(tls bool, host string) string {
	if tls {
		return fmt.Sprintf("https://%s", host)
	}
	return fmt.Sprintf("http://%s", host)
}

func prometheusServiceAccount(cr *v1beta1.PlatformMonitoring) (*corev1.ServiceAccount, error) {
	sa := corev1.ServiceAccount{}
	err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.PrometheusServiceAccountAsset), 100).Decode(&sa)

	if err != nil {
		return nil, err
	}
	//Set parameters
	sa.SetGroupVersionKind(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "ServiceAccount"})
	sa.SetName(cr.GetNamespace() + "-" + utils.PrometheusComponentName)
	sa.SetNamespace(cr.GetNamespace())

	// Set annotations and labels for ServiceAccount in case
	if cr.Spec.Prometheus != nil && cr.Spec.Prometheus.ServiceAccount != nil {
		if sa.Annotations == nil && cr.Spec.Prometheus.ServiceAccount.Annotations != nil {
			sa.SetAnnotations(cr.Spec.Prometheus.ServiceAccount.Annotations)
		} else {
			for k, v := range cr.Spec.Prometheus.ServiceAccount.Annotations {
				sa.Annotations[k] = v
			}
		}

		if sa.Labels == nil && cr.Spec.Prometheus.ServiceAccount.Labels != nil {
			sa.SetLabels(cr.Spec.Prometheus.ServiceAccount.Labels)
		} else {
			for k, v := range cr.Spec.Prometheus.ServiceAccount.Labels {
				sa.Labels[k] = v
			}
		}
	}

	return &sa, nil
}

func prometheusClusterRole(cr *v1beta1.PlatformMonitoring) (*rbacv1.ClusterRole, error) {
	clusterRole := rbacv1.ClusterRole{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.PrometheusClusterRoleAsset), 100).Decode(&clusterRole); err != nil {
		return nil, err
	}
	//Set parameters
	clusterRole.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRole"})
	clusterRole.SetName(cr.GetNamespace() + "-" + utils.PrometheusComponentName)

	return &clusterRole, nil
}

func prometheusClusterRoleBinding(cr *v1beta1.PlatformMonitoring) (*rbacv1.ClusterRoleBinding, error) {
	clusterRoleBinding := rbacv1.ClusterRoleBinding{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.PrometheusClusterRoleBindingAsset), 100).Decode(&clusterRoleBinding); err != nil {
		return nil, err
	}
	//Set parameters
	clusterRoleBinding.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRoleBinding"})
	clusterRoleBinding.SetName(cr.GetNamespace() + "-" + utils.PrometheusComponentName)
	clusterRoleBinding.RoleRef.Name = cr.GetNamespace() + "-" + utils.PrometheusComponentName

	// Set namespace and name for all subjects
	for it := range clusterRoleBinding.Subjects {
		sub := &clusterRoleBinding.Subjects[it]
		sub.Name = cr.GetNamespace() + "-" + utils.PrometheusComponentName
		sub.Namespace = cr.GetNamespace()
	}
	return &clusterRoleBinding, nil
}

func prometheus(cr *v1beta1.PlatformMonitoring) (*promv1.Prometheus, error) {
	prom := promv1.Prometheus{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.PrometheusAsset), 100).Decode(&prom); err != nil {
		return nil, err
	}

	// Set parameters
	prom.SetNamespace(cr.GetNamespace())
	prom.Spec.ServiceAccountName = cr.GetNamespace() + "-" + utils.PrometheusComponentName

	if cr.Spec.Prometheus != nil {
		prom.Spec.RemoteWrite = cr.Spec.Prometheus.RemoteWrite
		prom.Spec.RemoteRead = cr.Spec.Prometheus.RemoteRead

		// Set Prometheus image
		prom.Spec.Image = &cr.Spec.Prometheus.Image

		// Set labels
		prom.Labels["app.kubernetes.io/instance"] = utils.GetTagFromImage(cr.Spec.Prometheus.Image)
		prom.Labels["app.kubernetes.io/version"] = utils.GetInstanceLabel(prom.GetName(), prom.GetNamespace())

		// Set Prometheus replicas
		if cr.Spec.Prometheus.Replicas != nil {
			prom.Spec.Replicas = cr.Spec.Prometheus.Replicas
		}
		// Set Prometheus Scrape Interval
		if cr.Spec.Prometheus.ScrapeInterval != nil {
			prom.Spec.ScrapeInterval = promv1.Duration(*cr.Spec.Prometheus.ScrapeInterval)
		}
		// Set Prometheus Scrape Timeout
		if cr.Spec.Prometheus.ScrapeTimeout != nil {
			prom.Spec.ScrapeTimeout = promv1.Duration(*cr.Spec.Prometheus.ScrapeTimeout)
		}
		// Set Prometheus Evaluation Interval
		if cr.Spec.Prometheus.EvaluationInterval != nil {
			prom.Spec.EvaluationInterval = promv1.Duration(*cr.Spec.Prometheus.ScrapeInterval)
		}

		// Add integration with AlertManager only when AlertManager should deploy
		if cr.Spec.AlertManager != nil {
			// prometheus-operator automatically should create service with name `alertmanager-operated`
			// when it deploy AlertManager, so we can use it
			ae := promv1.AlertmanagerEndpoints{
				Namespace: cr.GetNamespace(),
				Name:      "alertmanager-operated",
				Port:      intstr.FromString("web"),
			}
			prom.Spec.Alerting.Alertmanagers = append(prom.Spec.Alerting.Alertmanagers, ae)
		}
		// Set security context
		if cr.Spec.Prometheus.SecurityContext != nil {
			if prom.Spec.SecurityContext == nil {
				prom.Spec.SecurityContext = &corev1.PodSecurityContext{}
			}
			if cr.Spec.Prometheus.SecurityContext.RunAsUser != nil {
				prom.Spec.SecurityContext.RunAsUser = cr.Spec.Prometheus.SecurityContext.RunAsUser
			}
			if cr.Spec.Prometheus.SecurityContext.FSGroup != nil {
				prom.Spec.SecurityContext.FSGroup = cr.Spec.Prometheus.SecurityContext.FSGroup
			}
		}
		// Set resources for Prometheus deployment
		if cr.Spec.Prometheus.Resources.Size() > 0 {
			prom.Spec.Resources = cr.Spec.Prometheus.Resources
		}
		// Set secrets for Prometheus deployment
		if len(cr.Spec.Prometheus.Secrets) > 0 {
			prom.Spec.Secrets = cr.Spec.Prometheus.Secrets
		}
		// Set nodeSelector for Prometheus deployment
		if cr.Spec.Prometheus.NodeSelector != nil {
			prom.Spec.NodeSelector = cr.Spec.Prometheus.NodeSelector
		}
		// Set affinity for Prometheus deployment
		if cr.Spec.Prometheus.Affinity != nil {
			prom.Spec.Affinity = cr.Spec.Prometheus.Affinity
		}
		// Set storage spec to specify how storage shall be used
		if cr.Spec.Prometheus.Storage != nil {
			prom.Spec.Storage = cr.Spec.Prometheus.Storage
		}
		// Set additional volumes for StatefulSet
		if cr.Spec.Prometheus.Volumes != nil {
			prom.Spec.Volumes = cr.Spec.Prometheus.Volumes
		}
		// Set additional volumeMounts for each Prometheus container. The current container names are:
		// `prometheus`, `prometheus-config-reloader`, `rules-configmap-reloader`, and `thanos-sidecar`
		if cr.Spec.Prometheus.VolumeMounts != nil {
			for it := range prom.Spec.Containers {
				c := &prom.Spec.Containers[it]

				// Set additional volumeMounts only for prometheus container
				if c.Name == "prometheus" {
					c.VolumeMounts = cr.Spec.Prometheus.VolumeMounts
				}
			}
		}
		// Set Retention - determines when to remove old data
		if cr.Spec.Prometheus.Retention != "" {
			prom.Spec.Retention = promv1.Duration(cr.Spec.Prometheus.Retention)
		}
		// Set Query - defines the query command line flags when starting Prometheus
		// Set by specific properties to keep default LookbackDelta value, if user want to override another parameters
		if cr.Spec.Prometheus.Query != nil {
			prom.Spec.Query = cr.Spec.Prometheus.Query
		}
		// Set enableAdminAPI - enable access to prometheus web admin API
		if cr.Spec.Prometheus.EnableAdminAPI {
			prom.Spec.EnableAdminAPI = cr.Spec.Prometheus.EnableAdminAPI
		}
		// Set RetentionSize - determines the maximum number of bytes that storage blocks can use
		if cr.Spec.Prometheus.RetentionSize != "" {
			prom.Spec.RetentionSize = promv1.ByteSize(cr.Spec.Prometheus.RetentionSize)
		}
		// Set alerting
		if cr.Spec.Prometheus.Alerting != nil {
			prom.Spec.Alerting = cr.Spec.Prometheus.Alerting
		}
		// Set external labels
		prom.Spec.ExternalLabels = cr.Spec.Prometheus.ExternalLabels

		// Set external URL
		if cr.Spec.Prometheus.ExternalURL != "" {
			prom.Spec.ExternalURL = cr.Spec.Prometheus.ExternalURL
		} else {
			if cr.Spec.Prometheus.Ingress != nil && cr.Spec.Prometheus.Ingress.IsInstall() && cr.Spec.Prometheus.Ingress.Host != "" {
				prom.Spec.ExternalURL = getPrometheusexternalURL(IsPrometheusTLSEnabled(cr), cr.Spec.Prometheus.Ingress.Host)
			}
		}
		// Set additional containers
		if cr.Spec.Prometheus.Containers != nil {
			prom.Spec.Containers = cr.Spec.Prometheus.Containers
		}
		// Handle integrations
		if cr.Spec.Integration != nil && cr.Spec.Integration.StackDriverIntegration != nil {
			// Configure sidecar for integration with GCO
			sidecar := corev1.Container{
				Name:            utils.StackdriverPrometheusSidecarName,
				Image:           cr.Spec.Integration.StackDriverIntegration.Image,
				ImagePullPolicy: "IfNotPresent",
				Ports:           []corev1.ContainerPort{{Name: "sidecar", ContainerPort: 9091}},
				VolumeMounts:    []corev1.VolumeMount{{MountPath: "/data", Name: "prometheus-k8s-db"}},
				Args: []string{
					"--stackdriver.project-id=" + cr.Spec.Integration.StackDriverIntegration.ProjectID,
					"--prometheus.wal-directory=/data/wal",
					"--stackdriver.kubernetes.location=" + cr.Spec.Integration.StackDriverIntegration.Location,
					"--stackdriver.kubernetes.cluster-name=" + cr.Spec.Integration.StackDriverIntegration.Cluster,
				},
			}
			for _, filter := range cr.Spec.Integration.StackDriverIntegration.MetricsFilters {
				sidecar.Args = append(sidecar.Args, "--include="+filter)
			}
			containerIndex := -1
			for idx, c := range prom.Spec.Containers {
				if c.Name == utils.StackdriverPrometheusSidecarName {
					containerIndex = idx
					break
				}
			}
			if containerIndex >= 0 {
				prom.Spec.Containers[containerIndex] = sidecar
			} else {
				prom.Spec.Containers = append(prom.Spec.Containers, sidecar)
			}
		}

		if cr.Spec.Auth != nil && cr.Spec.OAuthProxy != nil {
			prom.Spec.Secrets = append(prom.Spec.Secrets, "oauth2-proxy-config")

			externalURL := "http://"
			if cr.Spec.Prometheus.Ingress != nil && cr.Spec.Prometheus.Ingress.IsInstall() && cr.Spec.Prometheus.Ingress.Host != "" {
				externalURL += cr.Spec.Prometheus.Ingress.Host
			}
			// Volume mounts for oauth2-proxy sidecar
			var vms []corev1.VolumeMount

			// Add oauth2-proxy config
			vms = append(vms, corev1.VolumeMount{MountPath: utils.OAuthProxySecretDir, Name: "secret-oauth2-proxy-config"})

			if cr.Spec.Auth.TLSConfig != nil {
				// Add CA secret
				if cr.Spec.Auth.TLSConfig.CASecret != nil {
					prom.Spec.Secrets = append(prom.Spec.Secrets, cr.Spec.Auth.TLSConfig.CASecret.Name)
					vms = append(vms, corev1.VolumeMount{
						MountPath: utils.TlsCertificatesSecretDir + "/" + cr.Spec.Auth.TLSConfig.CASecret.Name,
						Name:      "secret-" + cr.Spec.Auth.TLSConfig.CASecret.Name,
					})
				}
				// Add Cert secret
				if cr.Spec.Auth.TLSConfig.CertSecret != nil {
					prom.Spec.Secrets = append(prom.Spec.Secrets, cr.Spec.Auth.TLSConfig.CertSecret.Name)
					vms = append(vms, corev1.VolumeMount{
						MountPath: utils.TlsCertificatesSecretDir + "/" + cr.Spec.Auth.TLSConfig.CertSecret.Name,
						Name:      "secret-" + cr.Spec.Auth.TLSConfig.CertSecret.Name,
					})
				}
				// Add Key secret
				if cr.Spec.Auth.TLSConfig.KeySecret != nil {
					prom.Spec.Secrets = append(prom.Spec.Secrets, cr.Spec.Auth.TLSConfig.KeySecret.Name)
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
					"--upstream=http://localhost:9090",
					"--config=/etc/oauth-proxy/oauth2-proxy.cfg",
				},
			}
			containerIndex := -1
			for idx, c := range prom.Spec.Containers {
				if c.Name == utils.OAuthProxyName {
					containerIndex = idx
					break
				}
			}
			if containerIndex > 0 {
				prom.Spec.Containers[containerIndex] = sidecar
			} else {
				prom.Spec.Containers = append(prom.Spec.Containers, sidecar)
			}
		}
		// Set tolerations for Prometheus
		if cr.Spec.Prometheus.Tolerations != nil {
			prom.Spec.Tolerations = cr.Spec.Prometheus.Tolerations
		}

		// Set selector for rules
		if cr.Spec.Prometheus.RuleSelector != nil {
			prom.Spec.RuleSelector = cr.Spec.Prometheus.RuleSelector
		}

		// Set selector for podMonitors
		if cr.Spec.Prometheus.PodMonitorSelector != nil {
			prom.Spec.PodMonitorSelector = cr.Spec.Prometheus.PodMonitorSelector
		}

		// Set selector for serviceMonitor
		if cr.Spec.Prometheus.ServiceMonitorSelector != nil {
			prom.Spec.ServiceMonitorSelector = cr.Spec.Prometheus.ServiceMonitorSelector
		}

		// Set namespaceSelector for rules
		if cr.Spec.Prometheus.RuleNamespaceSelector != nil {
			prom.Spec.RuleNamespaceSelector = cr.Spec.Prometheus.RuleNamespaceSelector
		}

		// Set namespaceSelector for podMonitors
		if cr.Spec.Prometheus.PodMonitorNamespaceSelector != nil {
			prom.Spec.PodMonitorNamespaceSelector = cr.Spec.Prometheus.PodMonitorNamespaceSelector
		}

		// Set namespaceSelector for serviceMonitor
		if cr.Spec.Prometheus.ServiceMonitorNamespaceSelector != nil {
			prom.Spec.ServiceMonitorNamespaceSelector = cr.Spec.Prometheus.ServiceMonitorNamespaceSelector
		}

		// Set PodMetadata.Labels
		prom.Spec.PodMetadata = &promv1.EmbeddedObjectMetadata{Labels: map[string]string{
			"name":                         "prometheus",
			"app.kubernetes.io/name":       "prometheus",
			"app.kubernetes.io/instance":   utils.GetInstanceLabel("prometheus", prom.GetNamespace()),
			"app.kubernetes.io/component":  "prometheus",
			"app.kubernetes.io/part-of":    "monitoring",
			"app.kubernetes.io/version":    utils.GetTagFromImage(cr.Spec.Prometheus.Image),
			"app.kubernetes.io/managed-by": "monitoring-operator",
		}}

		if prom.Spec.PodMetadata != nil {
			if prom.Spec.PodMetadata.Labels == nil && cr.Spec.Prometheus.Labels != nil {
				prom.Spec.PodMetadata.Labels = cr.Spec.Prometheus.Labels
			} else {
				for k, v := range cr.Spec.Prometheus.Labels {
					prom.Spec.PodMetadata.Labels[k] = v
				}
			}

			if prom.Spec.PodMetadata.Annotations == nil && cr.Spec.Prometheus.Annotations != nil {
				prom.Spec.PodMetadata.Annotations = cr.Spec.Prometheus.Annotations
			} else {
				for k, v := range cr.Spec.Prometheus.Annotations {
					prom.Spec.PodMetadata.Annotations[k] = v
				}
			}

			if len(strings.TrimSpace(cr.Spec.Prometheus.PriorityClassName)) > 0 {
				prom.Spec.PriorityClassName = cr.Spec.Prometheus.PriorityClassName
			}
		}

		if IsPrometheusTLSEnabled(cr) {
			promWebSpec := &promv1.PrometheusWebSpec{}
			if cr.Spec.Prometheus.TLSConfig.WebTLSConfig != nil {
				promWebSpec.TLSConfig = cr.Spec.Prometheus.TLSConfig.WebTLSConfig
				if cr.Spec.Prometheus.TLSConfig.WebTLSConfig.ClientAuthType == "" &&
					(cr.Spec.Prometheus.TLSConfig.WebTLSConfig.ClientCA.ConfigMap != nil ||
						cr.Spec.Prometheus.TLSConfig.WebTLSConfig.ClientCA.Secret != nil) {
					promWebSpec.TLSConfig.ClientAuthType = tls.VerifyClientCertIfGiven.String()
				}
			}
			// If GenerateCerts is enabled, we have to use certificates generated by cert-manager
			if cr.Spec.Prometheus.TLSConfig.GenerateCerts != nil && cr.Spec.Prometheus.TLSConfig.GenerateCerts.Enabled {
				generatedSecretName := "prometheus-cert-manager-tls"
				if cr.Spec.Prometheus.TLSConfig.GenerateCerts.SecretName != "" {
					generatedSecretName = cr.Spec.Prometheus.TLSConfig.GenerateCerts.SecretName
				}
				caSecret := promv1.SecretOrConfigMap{
					Secret: &corev1.SecretKeySelector{
						Key: "ca.crt",
					},
				}
				caSecret.Secret.Name = generatedSecretName
				certSecret := promv1.SecretOrConfigMap{
					Secret: &corev1.SecretKeySelector{
						Key: "tls.crt",
					},
				}
				certSecret.Secret.Name = generatedSecretName
				keySecretKeySelector := corev1.SecretKeySelector{
					Key: "tls.key",
				}
				keySecretKeySelector.Name = generatedSecretName
				if promWebSpec.TLSConfig != nil {
					promWebSpec.TLSConfig.KeySecret = keySecretKeySelector
					promWebSpec.TLSConfig.Cert = certSecret
					promWebSpec.TLSConfig.ClientCA = caSecret
				} else {
					promWebSpec.TLSConfig = &promv1.WebTLSConfig{
						KeySecret: keySecretKeySelector,
						Cert:      certSecret,
						ClientCA:  caSecret,
					}
				}
				if promWebSpec.TLSConfig.ClientAuthType == "" {
					promWebSpec.TLSConfig.ClientAuthType = tls.VerifyClientCertIfGiven.String()
				}
			}
			prom.Spec.Web = promWebSpec
		}
		if cr.Spec.Prometheus.ReplicaExternalLabelName != nil {
			prom.Spec.ReplicaExternalLabelName = cr.Spec.Prometheus.ReplicaExternalLabelName
		}

		if cr.Spec.Prometheus.EnableFeatures != nil {
			prom.Spec.EnableFeatures = cr.Spec.Prometheus.EnableFeatures
		}
	}
	return &prom, nil
}

func prometheusIngressV1beta1(cr *v1beta1.PlatformMonitoring) (*networkingv1beta1.Ingress, error) {
	ingress := networkingv1beta1.Ingress{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.PrometheusIngressAsset), 100).Decode(&ingress); err != nil {
		return nil, err
	}
	//Set parameters
	ingress.SetGroupVersionKind(schema.GroupVersionKind{Group: "networking.k8s.io", Version: "v1beta1", Kind: "Ingress"})
	ingress.SetName(cr.GetNamespace() + "-" + utils.PrometheusComponentName)
	ingress.SetNamespace(cr.GetNamespace())

	if cr.Spec.Prometheus != nil && cr.Spec.Prometheus.Ingress != nil && cr.Spec.Prometheus.Ingress.IsInstall() {
		// Check that ingress host is specified.
		if cr.Spec.Prometheus.Ingress.Host == "" {
			return nil, errors.New("host for ingress can not be empty")
		}

		// Add rule for prometheus UI
		rule := networkingv1beta1.IngressRule{Host: cr.Spec.Prometheus.Ingress.Host}
		serviceName := utils.PrometheusServiceName
		servicePort := intstr.FromInt(utils.PrometheusServicePort)

		if cr.Spec.Auth != nil && cr.Spec.OAuthProxy != nil {
			serviceName = utils.PrometheusOAuthProxyServiceName
			servicePort = intstr.FromString(utils.OAuthProxyServicePortName)
		}
		rule.HTTP = &networkingv1beta1.HTTPIngressRuleValue{
			Paths: []networkingv1beta1.HTTPIngressPath{
				{
					Path: "/",
					Backend: networkingv1beta1.IngressBackend{
						ServiceName: serviceName,
						ServicePort: servicePort,
					},
				},
			},
		}
		ingress.Spec.Rules = []networkingv1beta1.IngressRule{rule}

		if cr.Spec.Prometheus.Ingress.IngressClassName != nil {
			ingress.Spec.IngressClassName = cr.Spec.Prometheus.Ingress.IngressClassName
		}

		// Set annotations
		ingress.SetAnnotations(cr.Spec.Prometheus.Ingress.Annotations)

		// Configure TLS if TLS secret name is set
		if cr.Spec.Prometheus.Ingress.TLSSecretName != "" {
			if ingress.Annotations == nil {
				ingress.SetAnnotations(map[string]string{"nginx.ingress.kubernetes.io/backend-protocol": "HTTPS"})
			} else {
				ingress.Annotations["nginx.ingress.kubernetes.io/backend-protocol"] = "HTTPS"
			}
			ingress.Spec.TLS = []networkingv1beta1.IngressTLS{
				{
					Hosts:      []string{cr.Spec.Prometheus.Ingress.Host},
					SecretName: cr.Spec.Prometheus.Ingress.TLSSecretName,
				},
			}
		}
		// If GenerateCerts is enabled, we have to use certificate generated by cert-manager
		if IsPrometheusTLSEnabled(cr) {
			if cr.Spec.Prometheus.TLSConfig.GenerateCerts != nil && cr.Spec.Prometheus.TLSConfig.GenerateCerts.Enabled {
				if ingress.Annotations == nil {
					ingress.SetAnnotations(map[string]string{"nginx.ingress.kubernetes.io/backend-protocol": "HTTPS"})
				} else {
					ingress.Annotations["nginx.ingress.kubernetes.io/backend-protocol"] = "HTTPS"
				}
				generatedSecretName := "prometheus-cert-manager-tls"
				if cr.Spec.Prometheus.TLSConfig.GenerateCerts.SecretName != "" {
					generatedSecretName = cr.Spec.Prometheus.TLSConfig.GenerateCerts.SecretName
				}
				ingress.Spec.TLS = []networkingv1beta1.IngressTLS{
					{
						Hosts:      []string{cr.Spec.Prometheus.Ingress.Host},
						SecretName: generatedSecretName,
					},
				}
			}
		}

		// Set labels with saving default labels
		ingress.Labels["name"] = utils.TruncLabel(ingress.GetName())
		ingress.Labels["app.kubernetes.io/name"] = utils.TruncLabel(ingress.GetName())
		ingress.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(ingress.GetName(), ingress.GetNamespace())
		ingress.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Prometheus.Image)

		for lKey, lValue := range cr.Spec.Prometheus.Ingress.Labels {
			ingress.GetLabels()[lKey] = lValue
		}
	}
	return &ingress, nil
}

func prometheusIngressV1(cr *v1beta1.PlatformMonitoring) (*networkingv1.Ingress, error) {
	ingress := networkingv1.Ingress{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.PrometheusIngressAsset), 100).Decode(&ingress); err != nil {
		return nil, err
	}
	//Set parameters
	ingress.SetGroupVersionKind(schema.GroupVersionKind{Group: "networking.k8s.io", Version: "v1", Kind: "Ingress"})
	ingress.SetName(cr.GetNamespace() + "-" + utils.PrometheusComponentName)
	ingress.SetNamespace(cr.GetNamespace())

	if cr.Spec.Prometheus != nil && cr.Spec.Prometheus.Ingress != nil && cr.Spec.Prometheus.Ingress.IsInstall() {
		// Check that ingress host is specified.
		if cr.Spec.Prometheus.Ingress.Host == "" {
			return nil, errors.New("host for ingress can not be empty")
		}

		var ingressServiceBackend *networkingv1.IngressServiceBackend
		if cr.Spec.Auth != nil && cr.Spec.OAuthProxy != nil {
			ingressServiceBackend = &networkingv1.IngressServiceBackend{
				Name: utils.PrometheusOAuthProxyServiceName,
				Port: networkingv1.ServiceBackendPort{
					Name: utils.OAuthProxyServicePortName,
				},
			}
		} else {
			ingressServiceBackend = &networkingv1.IngressServiceBackend{
				Name: utils.PrometheusServiceName,
				Port: networkingv1.ServiceBackendPort{
					Number: utils.PrometheusServicePort,
				},
			}
		}
		pathType := networkingv1.PathTypePrefix
		// Add rule for prometheus UI
		rule := networkingv1.IngressRule{Host: cr.Spec.Prometheus.Ingress.Host}
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

		if cr.Spec.Prometheus.Ingress.IngressClassName != nil {
			ingress.Spec.IngressClassName = cr.Spec.Prometheus.Ingress.IngressClassName
		}

		// Set annotations
		ingress.SetAnnotations(cr.Spec.Prometheus.Ingress.Annotations)

		// Configure TLS if TLS secret name is set
		if cr.Spec.Prometheus.Ingress.TLSSecretName != "" {
			if ingress.Annotations == nil {
				ingress.SetAnnotations(map[string]string{"nginx.ingress.kubernetes.io/backend-protocol": "HTTPS"})
			} else {
				ingress.Annotations["nginx.ingress.kubernetes.io/backend-protocol"] = "HTTPS"
			}
			ingress.Spec.TLS = []networkingv1.IngressTLS{
				{
					Hosts:      []string{cr.Spec.Prometheus.Ingress.Host},
					SecretName: cr.Spec.Prometheus.Ingress.TLSSecretName,
				},
			}
		}
		// If GenerateCerts is enabled, we have to use certificate generated by cert-manager
		if IsPrometheusTLSEnabled(cr) {
			if cr.Spec.Prometheus.TLSConfig.GenerateCerts != nil && cr.Spec.Prometheus.TLSConfig.GenerateCerts.Enabled {
				if ingress.Annotations == nil {
					ingress.SetAnnotations(map[string]string{"nginx.ingress.kubernetes.io/backend-protocol": "HTTPS"})
				} else {
					ingress.Annotations["nginx.ingress.kubernetes.io/backend-protocol"] = "HTTPS"
				}
				generatedSecretName := "prometheus-cert-manager-tls"
				if cr.Spec.Prometheus.TLSConfig.GenerateCerts.SecretName != "" {
					generatedSecretName = cr.Spec.Prometheus.TLSConfig.GenerateCerts.SecretName
				}
				ingress.Spec.TLS = []networkingv1.IngressTLS{
					{
						Hosts:      []string{cr.Spec.Prometheus.Ingress.Host},
						SecretName: generatedSecretName,
					},
				}
			}
		}

		// Set labels with saving default labels
		ingress.Labels["name"] = utils.TruncLabel(ingress.GetName())
		ingress.Labels["app.kubernetes.io/name"] = utils.TruncLabel(ingress.GetName())
		ingress.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(ingress.GetName(), ingress.GetNamespace())
		ingress.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Prometheus.Image)

		for lKey, lValue := range cr.Spec.Prometheus.Ingress.Labels {
			ingress.GetLabels()[lKey] = lValue
		}
	}
	return &ingress, nil
}

func prometheusPodMonitor(cr *v1beta1.PlatformMonitoring) (*promv1.PodMonitor, error) {
	podMonitor := promv1.PodMonitor{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.PrometheusPodMonitorAsset), 100).Decode(&podMonitor); err != nil {
		return nil, err
	}
	//Set parameters
	podMonitor.SetGroupVersionKind(schema.GroupVersionKind{Group: "monitoring.coreos.com", Version: "v1", Kind: "PodMonitor"})
	podMonitor.SetName(cr.GetNamespace() + "-" + "prometheus-pod-monitor")
	podMonitor.SetNamespace(cr.GetNamespace())

	if cr.Spec.Prometheus != nil && cr.Spec.Prometheus.PodMonitor != nil && cr.Spec.Prometheus.PodMonitor.IsInstall() {
		cr.Spec.Prometheus.PodMonitor.OverridePodMonitor(&podMonitor)
	}
	if IsPrometheusTLSEnabled(cr) {
		var endpoints []promv1.PodMetricsEndpoint

		for _, item := range podMonitor.Spec.PodMetricsEndpoints {
			if item.Port == "web" {
				item.Scheme = "https"
				item.TLSConfig = &promv1.SafeTLSConfig{
					InsecureSkipVerify: ptr.To(true),
				}
			}
			endpoints = append(endpoints, item)
		}
		podMonitor.Spec.PodMetricsEndpoints = endpoints
	}
	return &podMonitor, nil
}

func IsPrometheusTLSEnabled(cr *v1beta1.PlatformMonitoring) bool {
	return cr.Spec.Prometheus != nil && cr.Spec.Prometheus.TLSConfig != nil &&
		(cr.Spec.Prometheus.TLSConfig.WebTLSConfig != nil ||
			cr.Spec.Prometheus.TLSConfig.GenerateCerts != nil && cr.Spec.Prometheus.TLSConfig.GenerateCerts.Enabled)
}
