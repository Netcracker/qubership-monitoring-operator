package grafana

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"maps"
	"strings"

	v1alpha1 "github.com/Netcracker/qubership-monitoring-operator/api/v1alpha1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/prometheus"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	vmetricsv1b1 "github.com/VictoriaMetrics/operator/api/operator/v1beta1"
	grafv1 "github.com/grafana-operator/grafana-operator/v4/api/integreatly/v1alpha1"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
)

//go:embed  assets/*.yaml
var assets embed.FS

func getGrafanaRootURL(protocol string, host string) string {
	if protocol == "" {
		protocol = "http"
	}
	return fmt.Sprintf("%v://%v/", protocol, host)
}

func grafana(cr *v1alpha1.PlatformMonitoring) (*grafv1.Grafana, error) {
	graf := grafv1.Grafana{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.GrafanaAsset), 100).Decode(&graf); err != nil {
		return nil, err
	}
	//Set parameters
	graf.SetGroupVersionKind(schema.GroupVersionKind{Group: "integreatly.org", Version: "v1alpha1", Kind: "Grafana"})
	graf.SetNamespace(cr.GetNamespace())

	if cr.Spec.Grafana != nil {
		if (cr.Spec.Grafana.Config != grafv1.GrafanaConfig{}) {
			graf.Spec.Config = cr.Spec.Grafana.Config
		}
		if cr.Spec.Grafana.DataStorage != nil {
			graf.Spec.DataStorage = cr.Spec.Grafana.DataStorage
		}
		if graf.Spec.Deployment != nil {
			graf.Spec.Deployment.EnvFrom = []corev1.EnvFromSource{
				{ConfigMapRef: &corev1.ConfigMapEnvSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: "grafana-extra-vars",
					},
					Optional: nil,
				},
				},
				{SecretRef: &corev1.SecretEnvSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: "grafana-extra-vars-secret",
					},
					Optional: nil,
				},
				},
			}
		} else {
			graf.Spec.Deployment = &grafv1.GrafanaDeployment{EnvFrom: []corev1.EnvFromSource{
				{ConfigMapRef: &corev1.ConfigMapEnvSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: "grafana-extra-vars",
					},
					Optional: nil,
				},
				},
				{SecretRef: &corev1.SecretEnvSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: "grafana-extra-vars-secret",
					},
					Optional: nil,
				},
				},
			},
			}
		}
		if cr.Spec.Grafana.Replicas != nil {
			graf.Spec.Deployment.Replicas = cr.Spec.Grafana.Replicas
		}
		// Set the configmap (with custom home dashboard) for grafana container
		if cr.Spec.Grafana.GrafanaHomeDashboard {
			graf.Spec.ConfigMaps = []string{"grafana-home-dashboard"}
		}
		graf.Spec.DashboardLabelSelector = cr.Spec.Grafana.DashboardLabelSelector
		graf.Spec.DashboardNamespaceSelector = cr.Spec.Grafana.DashboardNamespaceSelector
		// Add parameter to mount ldap config as secret
		graf.Spec.Secrets = append(graf.Spec.Secrets, "grafana-ldap-config")

		if cr.Spec.Auth != nil {
			if cr.Spec.Grafana.Ingress != nil && cr.Spec.Grafana.Ingress.IsInstall() && cr.Spec.Grafana.Ingress.Host != "" {
				if graf.Spec.Config.Server == nil {
					graf.Spec.Config.Server = &grafv1.GrafanaConfigServer{
						RootUrl: getGrafanaRootURL("", cr.Spec.Grafana.Ingress.Host),
					}
				} else {
					graf.Spec.Config.Server.RootUrl = getGrafanaRootURL(graf.Spec.Config.Server.Protocol, cr.Spec.Grafana.Ingress.Host)
				}
			}
			// Find all secrets names and add them under .secret.[]
			secrets := make(map[string]struct{})
			if cr.Spec.Auth.TLSConfig != nil {
				if cr.Spec.Auth.TLSConfig.CASecret != nil && cr.Spec.Auth.TLSConfig.CASecret.Name != "" && cr.Spec.Auth.TLSConfig.CASecret.Key != "" {
					secrets[cr.Spec.Auth.TLSConfig.CASecret.Name] = struct{}{}
				}
				if cr.Spec.Auth.TLSConfig.CertSecret != nil && cr.Spec.Auth.TLSConfig.CertSecret.Name != "" && cr.Spec.Auth.TLSConfig.CertSecret.Key != "" {
					secrets[cr.Spec.Auth.TLSConfig.CertSecret.Name] = struct{}{}
				}
				if cr.Spec.Auth.TLSConfig.KeySecret != nil && cr.Spec.Auth.TLSConfig.KeySecret.Name != "" && cr.Spec.Auth.TLSConfig.KeySecret.Key != "" {
					secrets[cr.Spec.Auth.TLSConfig.KeySecret.Name] = struct{}{}
				}
			}
			// Add only unique secrets names
			for k := range secrets {
				graf.Spec.Secrets = append(graf.Spec.Secrets, k)
			}
			// Create OAuth config
			if graf.Spec.Config.AuthGenericOauth == nil {
				graf.Spec.Config.AuthGenericOauth = &grafv1.GrafanaConfigAuthGenericOauth{}
			}
			if graf.Spec.Config.AuthGenericOauth != nil {
				ago := graf.Spec.Config.AuthGenericOauth

				enable := true
				ago.Enabled = &enable

				ago.AuthUrl = cr.Spec.Auth.LoginURL
				ago.TokenUrl = cr.Spec.Auth.TokenURL
				ago.ApiUrl = cr.Spec.Auth.UserInfoURL
				ago.Scopes = "openid profile"
			}
			// Set TLS config
			if cr.Spec.Auth.TLSConfig != nil {
				if cr.Spec.Auth.TLSConfig.InsecureSkipVerify != nil {
					graf.Spec.Config.AuthGenericOauth.TLSSkipVerifyInsecure = cr.Spec.Auth.TLSConfig.InsecureSkipVerify
				}
				CASecret := cr.Spec.Auth.TLSConfig.CASecret
				if CASecret != nil && CASecret.Name != "" && CASecret.Key != "" {
					graf.Spec.Config.AuthGenericOauth.TLSClientCa = fmt.Sprintf("/etc/grafana-secrets/%s/%s", CASecret.Name, CASecret.Key)
				}
				certSecret := cr.Spec.Auth.TLSConfig.CertSecret
				keySecret := cr.Spec.Auth.TLSConfig.KeySecret
				if certSecret != nil && keySecret != nil && certSecret.Name != "" && certSecret.Key != "" && keySecret.Name != "" && keySecret.Key != "" {
					graf.Spec.Config.AuthGenericOauth.TLSClientCert = fmt.Sprintf("/etc/grafana-secrets/%s/%s", certSecret.Name, certSecret.Key)
					graf.Spec.Config.AuthGenericOauth.TLSClientKey = fmt.Sprintf("/etc/grafana-secrets/%s/%s", keySecret.Name, keySecret.Key)
				}
			}
		}
		// Set security context
		if cr.Spec.Grafana.SecurityContext != nil {
			if graf.Spec.Deployment == nil {
				graf.Spec.Deployment = &grafv1.GrafanaDeployment{}
			}
			if graf.Spec.Deployment.SecurityContext == nil {
				graf.Spec.Deployment.SecurityContext = &corev1.PodSecurityContext{}
			}
			if cr.Spec.Grafana.SecurityContext.RunAsUser != nil {
				graf.Spec.Deployment.SecurityContext.RunAsUser = cr.Spec.Grafana.SecurityContext.RunAsUser
			}
			if cr.Spec.Grafana.SecurityContext.FSGroup != nil {
				graf.Spec.Deployment.SecurityContext.FSGroup = cr.Spec.Grafana.SecurityContext.FSGroup
			}
		}
		// Set resources for Grafana deployment
		if cr.Spec.Grafana.Resources.Size() > 0 {
			graf.Spec.Resources = &cr.Spec.Grafana.Resources
		}
		// Set tolerations for Grafana deployment
		if cr.Spec.Grafana.Tolerations != nil {
			graf.Spec.Deployment.Tolerations = cr.Spec.Grafana.Tolerations
		}
		// Set nodeSelector for Grafana deployment
		if cr.Spec.Grafana.NodeSelector != nil {
			graf.Spec.Deployment.NodeSelector = cr.Spec.Grafana.NodeSelector
		}
		// Set affinity for Grafana deployment
		if cr.Spec.Grafana.Affinity != nil {
			graf.Spec.Deployment.Affinity = cr.Spec.Grafana.Affinity
		}

		if len(strings.TrimSpace(cr.Spec.Grafana.PriorityClassName)) > 0 {
			graf.Spec.Deployment.PriorityClassName = cr.Spec.Grafana.PriorityClassName
		}

		// Set labels
		graf.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(graf.GetName(), graf.GetNamespace())
		graf.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Grafana.Image)
		if graf.Labels == nil && cr.Spec.Grafana.Labels != nil {
			graf.SetLabels(cr.Spec.Grafana.Labels)
		} else {
			for k, v := range cr.Spec.Grafana.Labels {
				graf.Labels[k] = v
			}
		}

		if graf.Annotations == nil && cr.Spec.Grafana.Annotations != nil {
			graf.SetAnnotations(cr.Spec.Grafana.Annotations)
		} else {
			for k, v := range cr.Spec.Grafana.Annotations {
				graf.Annotations[k] = v
			}
		}
		// Set labels
		graf.Spec.Deployment.Labels = map[string]string{
			"name":                         utils.TruncLabel(graf.GetName()),
			"app.kubernetes.io/name":       utils.TruncLabel(graf.GetName()),
			"app.kubernetes.io/instance":   utils.GetInstanceLabel(graf.GetName(), graf.GetNamespace()),
			"app.kubernetes.io/component":  "grafana",
			"app.kubernetes.io/part-of":    "monitoring",
			"app.kubernetes.io/version":    utils.GetTagFromImage(cr.Spec.Grafana.Image),
			"app.kubernetes.io/managed-by": "monitoring-operator",
		}
		if cr.Spec.Grafana.Labels != nil {
			for k, v := range cr.Spec.Grafana.Labels {
				graf.Spec.Deployment.Labels[k] = v
			}
		}

		if graf.Spec.Deployment.Annotations == nil && cr.Spec.Grafana.Annotations != nil {
			graf.Spec.Deployment.Annotations = cr.Spec.Grafana.Annotations
		} else {
			for k, v := range cr.Spec.Grafana.Annotations {
				graf.Spec.Deployment.Annotations[k] = v
			}
		}

		if graf.Spec.ServiceAccount != nil {
			if graf.Spec.ServiceAccount.Annotations == nil && cr.Spec.Grafana.ServiceAccount.Annotations != nil {
				graf.Spec.ServiceAccount.Annotations = cr.Spec.Grafana.ServiceAccount.Annotations
			} else {
				for k, v := range cr.Spec.Grafana.ServiceAccount.Annotations {
					graf.Spec.ServiceAccount.Annotations[k] = v
				}
			}

			if graf.Spec.ServiceAccount.Labels == nil && cr.Spec.Grafana.ServiceAccount.Labels != nil {
				graf.Spec.ServiceAccount.Labels = cr.Spec.Grafana.ServiceAccount.Labels
			} else {
				for k, v := range cr.Spec.Grafana.ServiceAccount.Labels {
					graf.Spec.ServiceAccount.Labels[k] = v
				}
			}
		}
	}
	return &graf, nil
}

func grafanaDataSource(cr *v1alpha1.PlatformMonitoring, KubeClient kubernetes.Interface, jaegerServices []corev1.Service, clickHouseServices []corev1.Service) (*grafv1.GrafanaDataSource, error) {
	dataSource := grafv1.GrafanaDataSource{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.GrafanaDataSourceAsset), 100).Decode(&dataSource); err != nil {
		return nil, err
	}
	// Set Interval for Grafana datasource
	var grafanaDatasourceInterval string = "30s"
	if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.VmOperator.IsInstall() {
		if cr.Spec.Victoriametrics.VmSingle.IsInstall() {
			vmSingle := vmetricsv1b1.VMSingle{}
			vmSingle.SetName(utils.VmComponentName)
			vmSingle.SetNamespace(cr.GetNamespace())
			if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.TLSEnabled {
				vmSingle.Spec.ExtraArgs = make(map[string]string)
				maps.Copy(vmSingle.Spec.ExtraArgs, map[string]string{"tls": "true"})
			}
			dataSource.Spec.Datasources[0].Url = vmSingle.AsURL()
		}
		if cr.Spec.Victoriametrics.VmCluster.IsInstall() {
			vmCluster := &vmetricsv1b1.VMCluster{}
			vmCluster.SetName(utils.VmComponentName)
			vmCluster.SetNamespace(cr.GetNamespace())
			vmCluster.Spec.VMSelect = cr.Spec.Victoriametrics.VmCluster.VmSelect
			if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.TLSEnabled {
				vmCluster.Spec.VMSelect.ExtraArgs = make(map[string]string)
				maps.Copy(vmCluster.Spec.VMSelect.ExtraArgs, map[string]string{"tls": "true"})
			}
			dataSource.Spec.Datasources[0].Url = vmCluster.VMSelectURL() + "/select/0/prometheus"
		}
		if cr.Spec.Victoriametrics.VmAgent.IsInstall() && len(strings.TrimSpace(cr.Spec.Victoriametrics.VmAgent.ScrapeInterval)) > 0 {
			grafanaDatasourceInterval = cr.Spec.Victoriametrics.VmAgent.ScrapeInterval
		}
	}
	// Set parameters
	dataSource.SetGroupVersionKind(schema.GroupVersionKind{Group: "integreatly.org", Version: "v1alpha1", Kind: "GrafanaDatasource"})
	dataSource.SetNamespace(cr.GetNamespace())

	// Set additional datasource for Promxy
	if cr.Spec.Promxy != nil && cr.Spec.Promxy.IsInstall() {
		// Set port for Promxy
		if cr.Spec.Promxy.Port != nil {
			dataSource.Spec.Datasources[1].Url = "http://promxy:" + fmt.Sprint(*cr.Spec.Promxy.Port)
			dataSource.Spec.Datasources[1].JsonData.TimeInterval = grafanaDatasourceInterval
		}
	} else {
		// If Promxy is not install, remove Promxy datasource
		dataSource.Spec.Datasources = dataSource.Spec.Datasources[:1]
	}
	// Set Jaeger datasource if Jaeger services is found
	if len(jaegerServices) > 0 {
		for _, jaegerService := range jaegerServices {
			var jaegerServicePort int32
			for _, port := range jaegerService.Spec.Ports {
				if port.Name == "http-query" {
					jaegerServicePort = port.Port
				}
			}
			var jaegerDataSourceName string
			if len(jaegerServices) > 1 {
				jaegerDataSourceName = "Jaeger " + jaegerService.GetNamespace() + "/" + jaegerService.GetName()
			} else {
				jaegerDataSourceName = "Jaeger"
			}
			jaegerDataSource := grafv1.GrafanaDataSourceFields{
				Access:    "proxy",
				Editable:  true,
				IsDefault: false,
				JsonData: grafv1.GrafanaDataSourceJsonData{
					TimeInterval:  grafanaDatasourceInterval,
					TlsSkipVerify: true,
					NodeGraph: grafv1.GrafanaDatasourceJsonNodeGraph{
						Enabled: true,
					},
				},
				Name:    jaegerDataSourceName,
				Type:    "jaeger",
				Version: 1,
				Url:     fmt.Sprintf("http://%s.%s.svc.cluster.local:%d", jaegerService.GetName(), jaegerService.GetNamespace(), jaegerServicePort),
			}
			dataSource.Spec.Datasources = append(dataSource.Spec.Datasources, jaegerDataSource)
		}
	}
	if len(clickHouseServices) > 0 {
		for _, clickHouseService := range clickHouseServices {
			var clickHouseDataSourceName string
			if len(clickHouseServices) > 1 {
				clickHouseDataSourceName = "ClickHouse_" + clickHouseService.GetNamespace()
			} else {
				clickHouseDataSourceName = "ClickHouse"
			}
			clickHouseDataSource := grafv1.GrafanaDataSourceFields{
				Access:    "proxy",
				Editable:  true,
				IsDefault: false,
				JsonData:  grafv1.GrafanaDataSourceJsonData{},
				Name:      clickHouseDataSourceName,
				Type:      "vertamedia-clickhouse-datasource",
				Version:   1,
				Url:       fmt.Sprintf("http://%s.%s.svc.cluster.local:8123", clickHouseService.GetName(), clickHouseService.GetNamespace()),
			}
			secret, err := KubeClient.CoreV1().Secrets(clickHouseService.GetNamespace()).Get(context.TODO(), utils.ClickHouseSecret, metav1.GetOptions{})
			if err == nil && len(secret.Data["username"]) != 0 && len(secret.Data["password"]) != 0 {
				clickHouseDataSource.BasicAuth = true
				clickHouseDataSource.BasicAuthUser = string(secret.Data["username"])
				clickHouseDataSource.SecureJsonData = grafv1.GrafanaDataSourceSecureJsonData{
					BasicAuthPassword: string(secret.Data["password"]),
				}
			}
			dataSource.Spec.Datasources = append(dataSource.Spec.Datasources, clickHouseDataSource)
		}
	}

	if prometheus.IsPrometheusTLSEnabled(cr) {
		dataSource.Spec.Datasources[0].Url = "https://prometheus-operated:9090"
	}

	dataSource.Spec.Datasources[0].JsonData.TimeInterval = grafanaDatasourceInterval

	return &dataSource, nil
}

func grafanaIngressV1beta1(cr *v1alpha1.PlatformMonitoring) (*v1beta1.Ingress, error) {
	ingress := v1beta1.Ingress{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.GrafanaIngressAsset), 100).Decode(&ingress); err != nil {
		return nil, err
	}
	// Set parameters
	ingress.SetGroupVersionKind(schema.GroupVersionKind{Group: "networking.k8s.io", Version: "v1beta1", Kind: "Ingress"})
	ingress.SetName(cr.GetNamespace() + "-" + utils.GrafanaComponentName)
	ingress.SetNamespace(cr.GetNamespace())

	if cr.Spec.Grafana != nil && cr.Spec.Grafana.Ingress != nil && cr.Spec.Grafana.Ingress.IsInstall() {
		// Check that ingress host is specified.
		if cr.Spec.Grafana.Ingress.Host == "" {
			return nil, errors.New("host for ingress can not be empty")
		}
		// Add rule for grafana UI
		rule := v1beta1.IngressRule{Host: cr.Spec.Grafana.Ingress.Host}
		rule.HTTP = &v1beta1.HTTPIngressRuleValue{
			Paths: []v1beta1.HTTPIngressPath{
				{
					Path: "/",
					Backend: v1beta1.IngressBackend{
						ServiceName: utils.GrafanaServiceName,
						ServicePort: intstr.FromInt(utils.GrafanaServicePort),
					},
				},
			},
		}
		ingress.Spec.Rules = []v1beta1.IngressRule{rule}

		// Configure TLS if TLS secret name is set
		if cr.Spec.Grafana.Ingress.TLSSecretName != "" {
			ingress.Spec.TLS = []v1beta1.IngressTLS{
				{
					Hosts:      []string{cr.Spec.Grafana.Ingress.Host},
					SecretName: cr.Spec.Grafana.Ingress.TLSSecretName,
				},
			}
		}

		if cr.Spec.Grafana.Ingress.IngressClassName != nil {
			ingress.Spec.IngressClassName = cr.Spec.Grafana.Ingress.IngressClassName
		}

		// Set annotations
		ingress.SetAnnotations(cr.Spec.Grafana.Ingress.Annotations)

		// Set labels with saving default labels
		ingress.Labels["name"] = utils.TruncLabel(ingress.GetName())
		ingress.Labels["app.kubernetes.io/name"] = utils.TruncLabel(ingress.GetName())
		ingress.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(ingress.GetName(), ingress.GetNamespace())
		ingress.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Grafana.Image)
		for lKey, lValue := range cr.Spec.Grafana.Ingress.Labels {
			ingress.GetLabels()[lKey] = lValue
		}
	}
	return &ingress, nil
}

func grafanaIngressV1(cr *v1alpha1.PlatformMonitoring) (*networkingv1.Ingress, error) {
	ingress := networkingv1.Ingress{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.GrafanaIngressAsset), 100).Decode(&ingress); err != nil {
		return nil, err
	}
	//Set parameters
	ingress.SetGroupVersionKind(schema.GroupVersionKind{Group: "networking.k8s.io", Version: "v1", Kind: "Ingress"})
	ingress.SetName(cr.GetNamespace() + "-" + utils.GrafanaComponentName)
	ingress.SetNamespace(cr.GetNamespace())

	if cr.Spec.Grafana != nil && cr.Spec.Grafana.Ingress != nil && cr.Spec.Grafana.Ingress.IsInstall() {
		// Check that ingress host is specified.
		if cr.Spec.Grafana.Ingress.Host == "" {
			return nil, errors.New("host for ingress can not be empty")
		}

		pathType := networkingv1.PathTypePrefix
		// Add rule for grafana UI
		rule := networkingv1.IngressRule{Host: cr.Spec.Grafana.Ingress.Host}
		rule.HTTP = &networkingv1.HTTPIngressRuleValue{
			Paths: []networkingv1.HTTPIngressPath{
				{
					Path:     "/",
					PathType: &pathType,
					Backend: networkingv1.IngressBackend{
						Service: &networkingv1.IngressServiceBackend{
							Name: utils.GrafanaServiceName,
							Port: networkingv1.ServiceBackendPort{
								Number: utils.GrafanaServicePort,
							},
						},
					},
				},
			},
		}
		ingress.Spec.Rules = []networkingv1.IngressRule{rule}

		// Configure TLS if TLS secret name is set
		if cr.Spec.Grafana.Ingress.TLSSecretName != "" {
			ingress.Spec.TLS = []networkingv1.IngressTLS{
				{
					Hosts:      []string{cr.Spec.Grafana.Ingress.Host},
					SecretName: cr.Spec.Grafana.Ingress.TLSSecretName,
				},
			}
		}

		if cr.Spec.Grafana.Ingress.IngressClassName != nil {
			ingress.Spec.IngressClassName = cr.Spec.Grafana.Ingress.IngressClassName
		}

		// Set annotations
		ingress.SetAnnotations(cr.Spec.Grafana.Ingress.Annotations)

		// Set labels with saving default labels
		ingress.Labels["name"] = utils.TruncLabel(ingress.GetName())
		ingress.Labels["app.kubernetes.io/name"] = utils.TruncLabel(ingress.GetName())
		ingress.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(ingress.GetName(), ingress.GetNamespace())
		ingress.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Grafana.Image)
		for lKey, lValue := range cr.Spec.Grafana.Ingress.Labels {
			ingress.GetLabels()[lKey] = lValue
		}
	}
	return &ingress, nil
}

func grafanaPodMonitor(cr *v1alpha1.PlatformMonitoring) (*promv1.PodMonitor, error) {
	podMonitor := promv1.PodMonitor{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.GrafanaPodMonitorAsset), 100).Decode(&podMonitor); err != nil {
		return nil, err
	}
	//Set parameters
	podMonitor.SetGroupVersionKind(schema.GroupVersionKind{Group: "monitoring.coreos.com", Version: "v1", Kind: "PodMonitor"})
	podMonitor.SetName(cr.GetNamespace() + "-" + "grafana-pod-monitor")
	podMonitor.SetNamespace(cr.GetNamespace())

	if cr.Spec.Grafana != nil && cr.Spec.Grafana.PodMonitor != nil && cr.Spec.Grafana.PodMonitor.IsInstall() {
		cr.Spec.Grafana.PodMonitor.OverridePodMonitor(&podMonitor)
	}
	return &podMonitor, nil
}
