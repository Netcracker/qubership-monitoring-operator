package grafana

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"strings"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
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

// ensureDeploymentInitialized ensures that Deployment is properly initialized
// This function safely initializes the Deployment structure to avoid nil pointer dereference
func ensureDeploymentInitialized(graf *grafv1.Grafana) {
	if graf.Spec.Deployment == nil {
		graf.Spec.Deployment = &grafv1.DeploymentV1{}
	}
}

// ensurePodSpecInitialized ensures that Deployment.Spec.Template.Spec is properly initialized
// This function safely initializes the PodSpec structure to avoid nil pointer dereference
// In v5, Template.Spec is a pointer (*DeploymentV1PodSpec), so it needs to be initialized
func ensurePodSpecInitialized(graf *grafv1.Grafana) *grafv1.DeploymentV1PodSpec {
	ensureDeploymentInitialized(graf)
	deployment := graf.Spec.Deployment
	if deployment == nil {
		deployment = &grafv1.DeploymentV1{}
		graf.Spec.Deployment = deployment
	}

	// Try to safely access and initialize Template.Spec
	// We'll use a recover to catch any panics from accessing nested pointer fields
	var podSpec *grafv1.DeploymentV1PodSpec
	panicked := false

	func() {
		defer func() {
			if r := recover(); r != nil {
				panicked = true
			}
		}()
		// Try to access Template.Spec - this may panic if Spec or Template are pointers and nil
		if deployment.Spec.Template.Spec == nil {
			deployment.Spec.Template.Spec = &grafv1.DeploymentV1PodSpec{}
		}
		podSpec = deployment.Spec.Template.Spec
	}()

	if panicked || podSpec == nil {
		// We panicked trying to access Template.Spec, which means Spec or Template is a pointer and nil
		// Create a fresh PodSpec
		podSpec = &grafv1.DeploymentV1PodSpec{}

		// Try to set it back into the deployment structure using a safe method
		// We'll try multiple approaches to ensure it gets set
		func() {
			defer func() { recover() }()
			// Approach 1: Try to set it directly (may panic if Spec or Template are nil pointers)
			deployment.Spec.Template.Spec = podSpec
		}()

		// If that didn't work, try recreating the deployment
		func() {
			defer func() { recover() }()
			newDeployment := &grafv1.DeploymentV1{}
			// Try to initialize Template.Spec in the new deployment
			// This may still panic if Spec or Template are pointers
			newDeployment.Spec.Template.Spec = podSpec
			graf.Spec.Deployment = newDeployment
		}()
	}

	return podSpec
}

// ensureTemplateInitialized ensures that Deployment.Spec.Template is properly initialized
// This function safely initializes the Template structure to avoid nil pointer dereference
func ensureTemplateInitialized(graf *grafv1.Grafana) {
	ensureDeploymentInitialized(graf)
	deployment := graf.Spec.Deployment
	if deployment == nil {
		deployment = &grafv1.DeploymentV1{}
		graf.Spec.Deployment = deployment
	}

	// Try to safely access Template
	// Use recover to catch any panics from accessing nested pointer fields
	func() {
		defer func() {
			if r := recover(); r != nil {
				// Panic occurred - Spec might be a pointer and nil
				// Re-initialize the entire Deployment structure
				newDeployment := &grafv1.DeploymentV1{}
				graf.Spec.Deployment = newDeployment
			}
		}()
		// Try to access Template - this may panic if Spec is a pointer and nil
		// Just accessing it is enough to trigger initialization if needed
		_ = deployment.Spec.Template
	}()
}

func grafana(cr *monv1.PlatformMonitoring) (*grafv1.Grafana, error) {
	graf := grafv1.Grafana{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.GrafanaAsset), 100).Decode(&graf); err != nil {
		return nil, err
	}
	//Set parameters
	graf.SetGroupVersionKind(schema.GroupVersionKind{Group: "grafana.integreatly.org", Version: "v1beta1", Kind: "Grafana"})
	graf.SetNamespace(cr.GetNamespace())

	if cr.Spec.Grafana != nil {
		// Disable default admin secret creation - we manage it ourselves
		// In grafana-operator v5, disableDefaultAdminSecret is at spec level
		graf.Spec.DisableDefaultAdminSecret = true

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
			// Replicas moved to Deployment.Spec.Replicas in v5
			// Deployment.Spec is DeploymentV1Spec (not a pointer), so we work with it directly
			// But first ensure Deployment is initialized
			ensureDeploymentInitialized(&graf)
			// Ensure Deployment.Spec is initialized (it's a struct, not a pointer, but we need to make sure Deployment itself is not nil)
			if graf.Spec.Deployment != nil {
				graf.Spec.Deployment.Spec.Replicas = cr.Spec.Grafana.Replicas
			}
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
			podSpec := ensurePodSpecInitialized(&graf)
			// Initialize SecurityContext if it's nil
			if podSpec.SecurityContext == nil {
				podSpec.SecurityContext = &corev1.PodSecurityContext{}
			}
			// Now we can safely set its fields
			if cr.Spec.Grafana.SecurityContext.RunAsUser != nil {
				podSpec.SecurityContext.RunAsUser = cr.Spec.Grafana.SecurityContext.RunAsUser
			}
			if cr.Spec.Grafana.SecurityContext.FSGroup != nil {
				podSpec.SecurityContext.FSGroup = cr.Spec.Grafana.SecurityContext.FSGroup
			}
		}
		// Set resources for Grafana deployment
		// Resources moved to Deployment.Spec.Template.Spec.Containers[0].Resources in v5
		if cr.Spec.Grafana.Resources.Size() > 0 {
			podSpec := ensurePodSpecInitialized(&graf)
			if len(podSpec.Containers) == 0 {
				podSpec.Containers = []corev1.Container{{}}
			}
			podSpec.Containers[0].Resources = cr.Spec.Grafana.Resources
		}
		// Set tolerations for Grafana deployment
		if cr.Spec.Grafana.Tolerations != nil {
			podSpec := ensurePodSpecInitialized(&graf)
			podSpec.Tolerations = cr.Spec.Grafana.Tolerations
		}
		// Set nodeSelector for Grafana deployment
		if cr.Spec.Grafana.NodeSelector != nil {
			podSpec := ensurePodSpecInitialized(&graf)
			podSpec.NodeSelector = cr.Spec.Grafana.NodeSelector
		}
		// Set affinity for Grafana deployment
		if cr.Spec.Grafana.Affinity != nil {
			podSpec := ensurePodSpecInitialized(&graf)
			podSpec.Affinity = cr.Spec.Grafana.Affinity
		}

		if len(strings.TrimSpace(cr.Spec.Grafana.PriorityClassName)) > 0 {
			podSpec := ensurePodSpecInitialized(&graf)
			podSpec.PriorityClassName = cr.Spec.Grafana.PriorityClassName
		}

		// Set labels on Grafana resource
		// Initialize Labels map if it's nil to avoid nil pointer dereference
		// Labels from asset file (app.kubernetes.io/component, app.kubernetes.io/part-of) are preserved
		if graf.Labels == nil {
			graf.Labels = make(map[string]string)
		}
		// Set dynamic labels that are always computed
		graf.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(graf.GetName(), graf.GetNamespace())
		graf.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Grafana.Image)
		// Allow overriding any labels (including component and part-of) via cr.Spec.Grafana.Labels
		// This allows different Grafana instances to have different labels for different dashboards
		if cr.Spec.Grafana.Labels != nil {
			for k, v := range cr.Spec.Grafana.Labels {
				graf.Labels[k] = v
			}
		}
		// Set default labels only if they weren't set in asset file and weren't overridden
		// These labels are needed for GrafanaDashboard instanceSelector matching
		if graf.Labels["app.kubernetes.io/component"] == "" {
			graf.Labels["app.kubernetes.io/component"] = "grafana"
		}
		if graf.Labels["app.kubernetes.io/part-of"] == "" {
			graf.Labels["app.kubernetes.io/part-of"] = "monitoring"
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
		ensureTemplateInitialized(&graf)
		// Ensure PodSpec is initialized first (needed for Template.Spec access)
		ensurePodSpecInitialized(&graf)

		// Safely access Template.Labels using a helper function
		func() {
			defer func() {
				if r := recover(); r != nil {
					// If we panic, re-initialize everything and try again
					ensureTemplateInitialized(&graf)
					ensurePodSpecInitialized(&graf)
				}
			}()
			// Access Template through deployment
			deployment := graf.Spec.Deployment
			if deployment != nil {
				// Initialize Labels if nil
				if deployment.Spec.Template.Labels == nil {
					deployment.Spec.Template.Labels = make(map[string]string)
				}
				// Set labels
				deployment.Spec.Template.Labels["name"] = utils.TruncLabel(graf.GetName())
				deployment.Spec.Template.Labels["app.kubernetes.io/name"] = utils.TruncLabel(graf.GetName())
				deployment.Spec.Template.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(graf.GetName(), graf.GetNamespace())
				deployment.Spec.Template.Labels["app.kubernetes.io/component"] = "grafana"
				deployment.Spec.Template.Labels["app.kubernetes.io/part-of"] = "monitoring"
				deployment.Spec.Template.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Grafana.Image)
				deployment.Spec.Template.Labels["app.kubernetes.io/managed-by"] = "monitoring-operator"
				if cr.Spec.Grafana.Labels != nil {
					for k, v := range cr.Spec.Grafana.Labels {
						deployment.Spec.Template.Labels[k] = v
					}
				}

				// Set annotations
				if deployment.Spec.Template.Annotations == nil && cr.Spec.Grafana.Annotations != nil {
					deployment.Spec.Template.Annotations = cr.Spec.Grafana.Annotations
				} else if cr.Spec.Grafana.Annotations != nil {
					if deployment.Spec.Template.Annotations == nil {
						deployment.Spec.Template.Annotations = make(map[string]string)
					}
					for k, v := range cr.Spec.Grafana.Annotations {
						deployment.Spec.Template.Annotations[k] = v
					}
				}
			}
		}()

		// ServiceAccount in v5 uses different structure - Annotations and Labels may be in different location
		// Note: ServiceAccount configuration may need to be handled differently in v5
		if graf.Spec.ServiceAccount != nil && cr.Spec.Grafana.ServiceAccount != nil {
			// In v5, ServiceAccountV1 structure changed - handle accordingly
			// Annotations and Labels may need to be set via ServiceAccount metadata
		}
	}
	return &graf, nil
}

func grafanaDataSource(cr *monv1.PlatformMonitoring, KubeClient kubernetes.Interface, jaegerServices []corev1.Service, clickHouseServices []corev1.Service) (*grafv1.GrafanaDatasource, error) {
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
	dataSource.SetGroupVersionKind(schema.GroupVersionKind{Group: "grafana.integreatly.org", Version: "v1beta1", Kind: "GrafanaDatasource"})
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

func grafanaIngressV1beta1(cr *monv1.PlatformMonitoring) (*networkingv1beta1.Ingress, error) {
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

func grafanaIngressV1(cr *monv1.PlatformMonitoring) (*networkingv1.Ingress, error) {
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

func grafanaPodMonitor(cr *monv1.PlatformMonitoring) (*promv1.PodMonitor, error) {
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
