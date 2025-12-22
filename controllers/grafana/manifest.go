package grafana

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"strings"

	v1beta1 "github.com/Netcracker/qubership-monitoring-operator/api"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/prometheus"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	vmetricsv1b1 "github.com/VictoriaMetrics/operator/api/operator/v1beta1"
	grafv1 "github.com/grafana/grafana-operator/v5/api/v1beta1"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
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

// ensureDeploymentInitialized ensures that Deployment and its Template.Spec are properly initialized
// This function safely initializes the Deployment structure to avoid nil pointer dereference
func ensureDeploymentInitialized(graf *grafv1.Grafana) {
	if graf.Spec.Deployment == nil {
		graf.Spec.Deployment = &grafv1.DeploymentV1{}
	}
	// In v5, Template is a PodTemplateSpec struct (not a pointer), so it should be initialized
	// as a zero-value struct. However, we need to ensure SecurityContext pointer is initialized
	// if we're going to access it. Initialize SecurityContext if it's nil
	// Note: We access Template.Spec.SecurityContext safely - if Template or Spec are structs (not pointers),
	// they will be zero-value initialized, so we only need to check SecurityContext pointer
	// Double-check Deployment is not nil before accessing nested fields
	if graf.Spec.Deployment == nil {
		return
	}
	// Initialize SecurityContext if it's nil to avoid nil pointer dereference
	// Access nested fields safely - Spec and Template are structs (not pointers), so they're zero-value initialized
	if graf.Spec.Deployment.Spec.Template.Spec.SecurityContext == nil {
		graf.Spec.Deployment.Spec.Template.Spec.SecurityContext = &corev1.PodSecurityContext{}
	}
}

func grafana(cr *v1beta1.PlatformMonitoring) (*grafv1.Grafana, error) {
	graf := grafv1.Grafana{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.GrafanaAsset), 100).Decode(&graf); err != nil {
		return nil, err
	}
	//Set parameters
	graf.SetGroupVersionKind(schema.GroupVersionKind{Group: "integreatly.org", Version: "v1beta1", Kind: "Grafana"})
	graf.SetNamespace(cr.GetNamespace())

	if cr.Spec.Grafana != nil {
		// Config is now runtime.RawExtension in grafana-operator v5
		// Only set Config if it's provided as RawExtension, otherwise configure Config fields below
		configProvidedAsRawExtension := cr.Spec.Grafana.Config != nil
		if configProvidedAsRawExtension {
			// Config is runtime.RawExtension, but graf.Spec.Config expects map[string]map[string]string
			// We need to unmarshal RawExtension to set Config properly
			// For now, skip setting Config if provided as RawExtension - it will be handled by grafana-operator
		}
		// DataStorage removed in grafana-operator v5
		// EnvFrom configuration moved to different structure in v5
		// Note: EnvFrom configuration may need to be handled differently in v5
		if cr.Spec.Grafana.Replicas != nil {
			ensureDeploymentInitialized(&graf)
			// Replicas moved to Deployment.Spec.Replicas in v5
			// Deployment.Spec is DeploymentV1Spec (not a pointer), so we work with it directly
			graf.Spec.Deployment.Spec.Replicas = cr.Spec.Grafana.Replicas
		}
		// ConfigMaps removed in grafana-operator v5 - use different approach if needed
		// Note: GrafanaHomeDashboard functionality may need alternative implementation
		// DashboardLabelSelector and DashboardNamespaceSelector removed or renamed in v5
		// Secrets removed or renamed in v5 - handle secrets differently if needed

		// Only configure Config fields if Config was not provided as RawExtension
		// Note: In v5, Config is map[string]map[string]string or runtime.RawExtension
		// We need to work with Config as runtime.RawExtension and marshal/unmarshal JSON
		if cr.Spec.Auth != nil && !configProvidedAsRawExtension {
			// In grafana-operator v5, Config structure changed significantly
			// OAuth configuration needs to be handled via runtime.RawExtension or removed
			// Note: This functionality may need to be reimplemented for v5 API
		}
		// Set security context
		if cr.Spec.Grafana.SecurityContext != nil {
			ensureDeploymentInitialized(&graf)
			// SecurityContext is already initialized in ensureDeploymentInitialized
			// Double-check that Deployment and SecurityContext are initialized before accessing nested fields
			if graf.Spec.Deployment != nil && graf.Spec.Deployment.Spec.Template.Spec.SecurityContext != nil {
				// Now we can safely set its fields
				if cr.Spec.Grafana.SecurityContext.RunAsUser != nil {
					graf.Spec.Deployment.Spec.Template.Spec.SecurityContext.RunAsUser = cr.Spec.Grafana.SecurityContext.RunAsUser
				}
				if cr.Spec.Grafana.SecurityContext.FSGroup != nil {
					graf.Spec.Deployment.Spec.Template.Spec.SecurityContext.FSGroup = cr.Spec.Grafana.SecurityContext.FSGroup
				}
			}
		}
		// Set resources for Grafana deployment
		// Resources moved to Deployment.Spec.Template.Spec.Containers[0].Resources in v5
		if cr.Spec.Grafana.Resources.Size() > 0 {
			ensureDeploymentInitialized(&graf)
			if len(graf.Spec.Deployment.Spec.Template.Spec.Containers) == 0 {
				graf.Spec.Deployment.Spec.Template.Spec.Containers = []corev1.Container{{}}
			}
			graf.Spec.Deployment.Spec.Template.Spec.Containers[0].Resources = cr.Spec.Grafana.Resources
		}
		// Set tolerations for Grafana deployment
		if cr.Spec.Grafana.Tolerations != nil {
			ensureDeploymentInitialized(&graf)
			graf.Spec.Deployment.Spec.Template.Spec.Tolerations = cr.Spec.Grafana.Tolerations
		}
		// Set nodeSelector for Grafana deployment
		if cr.Spec.Grafana.NodeSelector != nil {
			ensureDeploymentInitialized(&graf)
			graf.Spec.Deployment.Spec.Template.Spec.NodeSelector = cr.Spec.Grafana.NodeSelector
		}
		// Set affinity for Grafana deployment
		if cr.Spec.Grafana.Affinity != nil {
			ensureDeploymentInitialized(&graf)
			graf.Spec.Deployment.Spec.Template.Spec.Affinity = cr.Spec.Grafana.Affinity
		}

		if len(strings.TrimSpace(cr.Spec.Grafana.PriorityClassName)) > 0 {
			ensureDeploymentInitialized(&graf)
			graf.Spec.Deployment.Spec.Template.Spec.PriorityClassName = cr.Spec.Grafana.PriorityClassName
		}

		// Set labels on Grafana resource
		// Initialize Labels map if it's nil to avoid nil pointer dereference
		if graf.Labels == nil {
			graf.Labels = make(map[string]string)
		}
		graf.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(graf.GetName(), graf.GetNamespace())
		graf.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Grafana.Image)
		if cr.Spec.Grafana.Labels != nil {
			for k, v := range cr.Spec.Grafana.Labels {
				graf.Labels[k] = v
			}
		}

		if graf.Annotations == nil && cr.Spec.Grafana.Annotations != nil {
			graf.SetAnnotations(cr.Spec.Grafana.Annotations)
		} else if cr.Spec.Grafana.Annotations != nil {
			for k, v := range cr.Spec.Grafana.Annotations {
				graf.Annotations[k] = v
			}
		}
		// Set labels on Deployment pod template - in v5, labels are in Deployment.Spec.Template.Labels
		ensureDeploymentInitialized(&graf)
		// Deployment.Spec is DeploymentV1Spec (not a pointer), so we work with it directly
		if graf.Spec.Deployment.Spec.Template.Labels == nil {
			graf.Spec.Deployment.Spec.Template.Labels = make(map[string]string)
		}
		graf.Spec.Deployment.Spec.Template.Labels["name"] = utils.TruncLabel(graf.GetName())
		graf.Spec.Deployment.Spec.Template.Labels["app.kubernetes.io/name"] = utils.TruncLabel(graf.GetName())
		graf.Spec.Deployment.Spec.Template.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(graf.GetName(), graf.GetNamespace())
		graf.Spec.Deployment.Spec.Template.Labels["app.kubernetes.io/component"] = "grafana"
		graf.Spec.Deployment.Spec.Template.Labels["app.kubernetes.io/part-of"] = "monitoring"
		graf.Spec.Deployment.Spec.Template.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Grafana.Image)
		graf.Spec.Deployment.Spec.Template.Labels["app.kubernetes.io/managed-by"] = "monitoring-operator"
		if cr.Spec.Grafana.Labels != nil {
			for k, v := range cr.Spec.Grafana.Labels {
				graf.Spec.Deployment.Spec.Template.Labels[k] = v
			}
		}

		if graf.Spec.Deployment.Spec.Template.Annotations == nil && cr.Spec.Grafana.Annotations != nil {
			graf.Spec.Deployment.Spec.Template.Annotations = cr.Spec.Grafana.Annotations
		} else if cr.Spec.Grafana.Annotations != nil {
			if graf.Spec.Deployment.Spec.Template.Annotations == nil {
				graf.Spec.Deployment.Spec.Template.Annotations = make(map[string]string)
			}
			for k, v := range cr.Spec.Grafana.Annotations {
				graf.Spec.Deployment.Spec.Template.Annotations[k] = v
			}
		}

		// ServiceAccount in v5 uses different structure - Annotations and Labels may be in different location
		// Note: ServiceAccount configuration may need to be handled differently in v5
		if graf.Spec.ServiceAccount != nil && cr.Spec.Grafana.ServiceAccount != nil {
			// In v5, ServiceAccountV1 structure changed - handle accordingly
			// Annotations and Labels may need to be set via ServiceAccount metadata
		}
	}
	return &graf, nil
}

func grafanaDataSource(cr *v1beta1.PlatformMonitoring, KubeClient kubernetes.Interface, jaegerServices []corev1.Service, clickHouseServices []corev1.Service) (*grafv1.GrafanaDatasource, error) {
	dataSource := grafv1.GrafanaDatasource{}
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
			if dataSource.Spec.Datasource != nil {
				dataSource.Spec.Datasource.URL = vmSingle.AsURL()
			}
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
			if dataSource.Spec.Datasource != nil {
				dataSource.Spec.Datasource.URL = vmCluster.VMSelectURL() + "/select/0/prometheus"
			}
		}
		if cr.Spec.Victoriametrics.VmAgent.IsInstall() && len(strings.TrimSpace(cr.Spec.Victoriametrics.VmAgent.ScrapeInterval)) > 0 {
			grafanaDatasourceInterval = cr.Spec.Victoriametrics.VmAgent.ScrapeInterval
		}
	}
	// Set parameters
	dataSource.SetGroupVersionKind(schema.GroupVersionKind{Group: "integreatly.org", Version: "v1beta1", Kind: "GrafanaDatasource"})
	dataSource.SetNamespace(cr.GetNamespace())

	// In grafana-operator v5, one GrafanaDatasource CR contains only one datasource
	// Promxy datasource would need to be a separate CR if needed
	// Note: Promxy datasource functionality may need to be reimplemented as separate CR

	// Set Jaeger datasource if Jaeger services is found
	// Note: In v5, each datasource needs to be a separate GrafanaDatasource CR
	// This functionality may need to be reimplemented to create multiple CRs

	// Set ClickHouse datasource if ClickHouse services is found
	// Note: In v5, each datasource needs to be a separate GrafanaDatasource CR
	// This functionality may need to be reimplemented to create multiple CRs

	if prometheus.IsPrometheusTLSEnabled(cr) && dataSource.Spec.Datasource != nil {
		dataSource.Spec.Datasource.URL = "https://prometheus-operated:9090"
	}

	// Set JSONData for timeInterval - in v5, JSONData is json.RawMessage
	if dataSource.Spec.Datasource != nil {
		jsonDataMap := make(map[string]interface{})
		if len(dataSource.Spec.Datasource.JSONData) > 0 {
			if err := json.Unmarshal(dataSource.Spec.Datasource.JSONData, &jsonDataMap); err == nil {
				jsonDataMap["timeInterval"] = grafanaDatasourceInterval
				if jsonBytes, err := json.Marshal(jsonDataMap); err == nil {
					dataSource.Spec.Datasource.JSONData = jsonBytes
				}
			}
		} else {
			jsonDataMap["timeInterval"] = grafanaDatasourceInterval
			if jsonBytes, err := json.Marshal(jsonDataMap); err == nil {
				dataSource.Spec.Datasource.JSONData = jsonBytes
			}
		}
	}

	return &dataSource, nil
}

func grafanaIngressV1beta1(cr *v1beta1.PlatformMonitoring) (*networkingv1beta1.Ingress, error) {
	ingress := networkingv1beta1.Ingress{}
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
		rule := networkingv1beta1.IngressRule{Host: cr.Spec.Grafana.Ingress.Host}
		rule.HTTP = &networkingv1beta1.HTTPIngressRuleValue{
			Paths: []networkingv1beta1.HTTPIngressPath{
				{
					Path: "/",
					Backend: networkingv1beta1.IngressBackend{
						ServiceName: utils.GrafanaServiceName,
						ServicePort: intstr.FromInt(utils.GrafanaServicePort),
					},
				},
			},
		}
		ingress.Spec.Rules = []networkingv1beta1.IngressRule{rule}

		// Configure TLS if TLS secret name is set
		if cr.Spec.Grafana.Ingress.TLSSecretName != "" {
			ingress.Spec.TLS = []networkingv1beta1.IngressTLS{
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
		// Initialize Labels map if it's nil to avoid nil pointer dereference
		if ingress.Labels == nil {
			ingress.Labels = make(map[string]string)
		}
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

func grafanaIngressV1(cr *v1beta1.PlatformMonitoring) (*networkingv1.Ingress, error) {
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
		// Initialize Labels map if it's nil to avoid nil pointer dereference
		if ingress.Labels == nil {
			ingress.Labels = make(map[string]string)
		}
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

func grafanaPodMonitor(cr *v1beta1.PlatformMonitoring) (*promv1.PodMonitor, error) {
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
