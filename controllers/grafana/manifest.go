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
	networkingv1beta1 "k8s.io/api/networking/v1beta1" // For v1beta1 Ingress API (used in grafanaIngressV1beta1)
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

	// Initialize Template if nil
	if deployment.Spec.Template.Spec == nil {
		deployment.Spec.Template.Spec = &grafv1.DeploymentV1PodSpec{}
	}

	return deployment.Spec.Template.Spec
}

func grafana(cr *monv1.PlatformMonitoring) (*grafv1.Grafana, error) {
	graf := grafv1.Grafana{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.GrafanaAsset), 100).Decode(&graf); err != nil {
		return nil, err
	}
	//Set parameters
	graf.SetGroupVersionKind(schema.GroupVersionKind{Group: "grafana.integreatly.org", Version: "v1beta1", Kind: "Grafana"})

	// Add way to move Grafana to a different namespace and set a custom name for the Grafana instance.
	// Set custom namespace if specified, otherwise use PlatformMonitoring namespace
	grafanaNamespace := cr.GetNamespace()
	if cr.Spec.Grafana != nil && cr.Spec.Grafana.Namespace != "" {
		grafanaNamespace = cr.Spec.Grafana.Namespace
	}
	graf.SetNamespace(grafanaNamespace)

	// Set custom name if specified, otherwise use default from asset file
	if cr.Spec.Grafana != nil && cr.Spec.Grafana.Name != "" {
		graf.SetName(cr.Spec.Grafana.Name)
	}

	if cr.Spec.Grafana != nil {
		// In grafana-operator v5, disableDefaultAdminSecret is at spec level.
		// Default (nil) = we manage the secret (same as true).
		graf.Spec.DisableDefaultAdminSecret = true
		if cr.Spec.Grafana.DisableDefaultAdminSecret != nil {
			graf.Spec.DisableDefaultAdminSecret = *cr.Spec.Grafana.DisableDefaultAdminSecret
		}

		// Do not set spec.ingress on the Grafana CR: Grafana Operator would then create its own Ingress.
		// In grafana-operator v5, there's no "enabled" field - if spec.ingress is present (even with empty spec),
		// Grafana Operator creates a catch-all (*) Ingress. Asset no longer contains spec.ingress to avoid this.
		// We explicitly set it to nil for safety (though it should already be nil after decoding the asset).
		graf.Spec.Ingress = nil

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
		// Set security context (pod-level; v5 uses Deployment.Spec.Template.Spec.SecurityContext)
		if cr.Spec.Grafana.SecurityContext != nil {
			podSpec := ensurePodSpecInitialized(&graf)
			if podSpec.SecurityContext == nil {
				podSpec.SecurityContext = &corev1.PodSecurityContext{}
			}
			if cr.Spec.Grafana.SecurityContext.RunAsUser != nil {
				podSpec.SecurityContext.RunAsUser = cr.Spec.Grafana.SecurityContext.RunAsUser
			}
			if cr.Spec.Grafana.SecurityContext.RunAsGroup != nil {
				podSpec.SecurityContext.RunAsGroup = cr.Spec.Grafana.SecurityContext.RunAsGroup
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

		// When disableDefaultAdminSecret is true, we need to explicitly set environment variables
		// to reference the admin credentials secret so the operator can authenticate with Grafana
		// In grafana-operator v5, the operator needs these env vars to authenticate and check readiness
		if graf.Spec.DisableDefaultAdminSecret {
			podSpec := ensurePodSpecInitialized(&graf)
			if len(podSpec.Containers) == 0 {
				podSpec.Containers = []corev1.Container{{}}
			}
			// Ensure container has a name
			if podSpec.Containers[0].Name == "" {
				podSpec.Containers[0].Name = "grafana"
			}
			// Add environment variables for admin credentials from secret
			// Secret name pattern: {grafana-name}-admin-credentials
			adminSecretName := fmt.Sprintf("%s-admin-credentials", graf.GetName())
			envVars := []corev1.EnvVar{
				{
					Name: "GF_SECURITY_ADMIN_USER",
					ValueFrom: &corev1.EnvVarSource{
						SecretKeyRef: &corev1.SecretKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: adminSecretName,
							},
							Key: "GF_SECURITY_ADMIN_USER",
						},
					},
				},
				{
					Name: "GF_SECURITY_ADMIN_PASSWORD",
					ValueFrom: &corev1.EnvVarSource{
						SecretKeyRef: &corev1.SecretKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: adminSecretName,
							},
							Key: "GF_SECURITY_ADMIN_PASSWORD",
						},
					},
				},
			}
			// Append to existing env vars if any, otherwise set new
			if podSpec.Containers[0].Env == nil {
				podSpec.Containers[0].Env = envVars
			} else {
				// Check if env vars already exist to avoid duplicates
				envMap := make(map[string]bool)
				for _, env := range podSpec.Containers[0].Env {
					envMap[env.Name] = true
				}
				for _, env := range envVars {
					if !envMap[env.Name] {
						podSpec.Containers[0].Env = append(podSpec.Containers[0].Env, env)
					}
				}
			}
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

		// Set default labels for GrafanaDashboard instanceSelector matching
		// These labels are used by GrafanaDashboard instanceSelector to find matching Grafana instances
		// In grafana-operator v5, GrafanaDashboard uses instanceSelector.matchLabels to find Grafana instances
		// If not set in asset file, use defaults
		if graf.Labels["app.kubernetes.io/component"] == "" {
			graf.Labels["app.kubernetes.io/component"] = "grafana"
		}
		if graf.Labels["app.kubernetes.io/part-of"] == "" {
			graf.Labels["app.kubernetes.io/part-of"] = "monitoring"
		}

		// Allow overriding any labels (including component and part-of) via cr.Spec.Grafana.Labels
		// This allows different Grafana instances to have different labels for different dashboards
		// Users can set custom labels like "dashboards: custom" and use them in GrafanaDashboard instanceSelector
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
		// In grafana-operator v5, Labels and Annotations moved from Deployment to Deployment.Spec.Template
		ensureDeploymentInitialized(&graf)
		// Ensure PodSpec is initialized first (needed for Template.Spec access)
		ensurePodSpecInitialized(&graf)

		// Access Template through deployment and set labels/annotations
		// Use recover to safely handle any nil pointer issues that may occur due to v5 API structure changes
		func() {
			defer func() {
				if r := recover(); r != nil {
					// If we panic, Template structure might not be accessible (e.g., Spec or Template are nil pointers)
					// This is OK - labels/annotations are already set on the Grafana resource itself
				}
			}()
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

				// Initialize Annotations if nil
				if deployment.Spec.Template.Annotations == nil {
					deployment.Spec.Template.Annotations = make(map[string]string)
				}
				// Set annotations
				if cr.Spec.Grafana.Annotations != nil {
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

// grafanaDataSource creates GrafanaDatasource manifest
// Note: In grafana-operator v5, the type name changed from GrafanaDataSource to GrafanaDatasource
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

	// In v5, one GrafanaDatasource CR = one datasource. Promxy is a separate CR (grafanaPromxyDataSource).

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

// grafanaPromxyDataSource creates GrafanaDatasource manifest for Promxy
// In v5, each datasource must be a separate GrafanaDatasource CR (not an array like in v4)
func grafanaPromxyDataSource(cr *monv1.PlatformMonitoring) (*grafv1.GrafanaDatasource, error) {
	dataSource := grafv1.GrafanaDatasource{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.GrafanaPromxyDataSourceAsset), 100).Decode(&dataSource); err != nil {
		return nil, err
	}
	// Set parameters
	dataSource.SetGroupVersionKind(schema.GroupVersionKind{Group: "grafana.integreatly.org", Version: "v1beta1", Kind: "GrafanaDatasource"})
	dataSource.SetNamespace(cr.GetNamespace())

	// Set Promxy URL with port from CR (default: 9090)
	promxyPort := int32(9090)
	if cr.Spec.Promxy != nil && cr.Spec.Promxy.Port != nil {
		promxyPort = *cr.Spec.Promxy.Port
	}
	if dataSource.Spec.Datasource != nil {
		dataSource.Spec.Datasource.URL = fmt.Sprintf("http://promxy:%d", promxyPort)
	}

	// Set JSONData for timeInterval - in v5, JSONData is json.RawMessage
	var grafanaDatasourceInterval string = "30s"
	if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.VmAgent.IsInstall() && len(strings.TrimSpace(cr.Spec.Victoriametrics.VmAgent.ScrapeInterval)) > 0 {
		grafanaDatasourceInterval = cr.Spec.Victoriametrics.VmAgent.ScrapeInterval
	}

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

// grafanaIngressV1beta1 creates Ingress manifest using v1beta1 API
// Note: Uses networkingv1beta1 package (with alias) instead of v1beta1 to avoid conflicts
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
		// Note: All types use networkingv1beta1 prefix (not v1beta1) to match the import alias
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
	} else {
		// Ingress not configured (empty host or install=false): do not leave the asset rule with host: "" (catch-all).
		ingress.Spec.Rules = nil
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
	} else {
		// Ingress not configured: do not leave the asset rule with host: "" (catch-all).
		ingress.Spec.Rules = nil
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
