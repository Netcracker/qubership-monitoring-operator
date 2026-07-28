package grafana

import (
	"embed"
	"encoding/json"
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
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
)

//go:embed  assets/*.yaml
var assets embed.FS

// TODO(#377): getGrafanaRootURL — to be used when spec.grafana.config propagation is implemented.
// Once spec.grafana.config is properly forwarded to Grafana CR spec.config, this helper can be
// used to automatically populate server.root_url from spec.grafana.ingress.host as a convenience
// default (users can override via grafana.config in values).
//
// func getGrafanaRootURL(protocol string, host string) string {
// 	if protocol == "" {
// 		protocol = "http"
// 	}
// 	return fmt.Sprintf("%v://%v/", protocol, host)
// }

// configureGrafanaOperatorKubeAuth enables JWT auth so grafana-operator can call the Grafana API
// without GF_SECURITY_ADMIN_* environment variables (admin credentials are file-based in Grafana).
func configureGrafanaOperatorKubeAuth(graf *grafv1.Grafana, operatorNamespace string) {
	graf.Spec.Client = &grafv1.GrafanaClient{UseKubeAuth: true}

	operatorSA := operatorNamespace + "-" + utils.GrafanaOperatorComponentName
	rolePath := fmt.Sprintf(
		"contains(sub, 'system:serviceaccount:%s:%s') && 'GrafanaAdmin' || 'None'",
		operatorNamespace, operatorSA,
	)

	jwt := ensureGrafanaConfigSection(graf, "auth.jwt")
	jwt["enabled"] = "true"
	jwt["header_name"] = "Authorization"
	jwt["expect_claims"] = `{"aud": ["operator.grafana.com"]}`
	jwt["username_claim"] = "sub"
	jwt["email_claim"] = "sub"
	jwt["auto_sign_up"] = "true"
	jwt["role_attribute_strict"] = "true"
	jwt["role_attribute_path"] = rolePath
	jwt["jwk_set_url"] = "https://${KUBERNETES_SERVICE_HOST}:${KUBERNETES_SERVICE_PORT_HTTPS}/openid/v1/jwks"
	jwt["jwk_set_bearer_token_file"] = "/var/run/secrets/kubernetes.io/serviceaccount/token"
	jwt["tls_client_ca"] = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
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
// In v5, Template is a pointer (*DeploymentV1PodTemplateSpec), and Template.Spec is also a pointer (*DeploymentV1PodSpec)
func ensurePodSpecInitialized(graf *grafv1.Grafana) *grafv1.DeploymentV1PodSpec {
	ensureDeploymentInitialized(graf)
	deployment := graf.Spec.Deployment
	if deployment == nil {
		deployment = &grafv1.DeploymentV1{}
		graf.Spec.Deployment = deployment
	}

	// Initialize Template if nil (Template is a pointer)
	if deployment.Spec.Template == nil {
		deployment.Spec.Template = &grafv1.DeploymentV1PodTemplateSpec{}
	}

	// Initialize Template.Spec if nil (Spec is a pointer)
	if deployment.Spec.Template.Spec == nil {
		deployment.Spec.Template.Spec = &grafv1.DeploymentV1PodSpec{}
	}

	return deployment.Spec.Template.Spec
}

// ensureGrafanaContainerInitialized guarantees that PodSpec has at least one container
// and that this container has a non-empty name, so that container-level fields like
// EnvFrom and VolumeMounts can be configured safely.
func ensureGrafanaContainerInitialized(podSpec *grafv1.DeploymentV1PodSpec) *corev1.Container {
	if len(podSpec.Containers) == 0 {
		podSpec.Containers = []corev1.Container{{}}
	}
	if podSpec.Containers[0].Name == "" {
		podSpec.Containers[0].Name = "grafana"
	}
	return &podSpec.Containers[0]
}

// ensureGrafanaConfigSection ensures graf.Spec.Config and target section are initialized.
func ensureGrafanaConfigSection(graf *grafv1.Grafana, section string) map[string]string {
	if graf.Spec.Config == nil {
		graf.Spec.Config = map[string]map[string]string{}
	}
	if graf.Spec.Config[section] == nil {
		graf.Spec.Config[section] = map[string]string{}
	}
	return graf.Spec.Config[section]
}

func grafana(cr *monv1.PlatformMonitoring) (*grafv1.Grafana, error) {
	return grafanaWithAdminPasswordChecksum(cr, "")
}

func grafanaWithAdminPasswordChecksum(cr *monv1.PlatformMonitoring, adminPasswordChecksum string) (*grafv1.Grafana, error) {
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
		// Always instruct grafana-operator NOT to auto-generate its own admin secret (disableDefaultAdminSecret=true).
		// Admin credential management is handled at a higher level:
		//   - disableDefaultAdminSecret=false (default): Helm creates the secret with configured values.
		//   - disableDefaultAdminSecret=true: user creates the secret manually (optional).
		// In both cases we do not want grafana-operator to generate a random secret on its own.
		graf.Spec.DisableDefaultAdminSecret = true

		// IMPORTANT: propagate PlatformMonitoring.spec.grafana.image into Grafana.spec.version.
		// In grafana-operator v5, spec.version accepts either a tag (e.g. "12.3.3") or a full image reference
		if cr.Spec.Grafana.Image != "" {
			graf.Spec.Version = cr.Spec.Grafana.Image
		}

		// Do not set spec.ingress on the Grafana CR: Grafana Operator would then create its own Ingress.
		// In grafana-operator v5, there's no "enabled" field - if spec.ingress is present (even with empty spec),
		// Grafana Operator creates a catch-all (*) Ingress. Asset no longer contains spec.ingress to avoid this.
		// We explicitly set it to nil for safety (though it should already be nil after decoding the asset).
		graf.Spec.Ingress = nil

		// TODO(#377): spec.grafana.config (runtime.RawExtension) is not yet propagated to Grafana CR
		// spec.config (map[string]map[string]string). When implemented, user-provided config keys
		// should be merged after operator defaults so they can override them (e.g. server.root_url).
		// The configProvidedAsRawExtension gate (below) should also be restored to allow users to
		// opt out of operator-managed config sections when providing their own full config.
		// DataStorage removed in grafana-operator v5
		// EnvFrom configuration moved to container-level in v5
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
		// Configure container-level settings (EnvFrom, volumes, home dashboard, etc.).
		podSpec := ensurePodSpecInitialized(&graf)
		container := ensureGrafanaContainerInitialized(podSpec)

		// Attach envFrom so that grafana picks up extraVars / extraVarsSecret
		// (GF_DASHBOARDS_DEFAULT_HOME_DASHBOARD_PATH and other settings).
		extraVarsCmRef := corev1.EnvFromSource{
			ConfigMapRef: &corev1.ConfigMapEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: "grafana-extra-vars",
				},
			},
		}
		extraVarsSecretRef := corev1.EnvFromSource{
			SecretRef: &corev1.SecretEnvSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: "grafana-extra-vars-secret",
				},
			},
		}
		// Avoid duplicating entries if they already exist
		hasCmEnvFrom := false
		hasSecretEnvFrom := false
		for _, ef := range container.EnvFrom {
			if ef.ConfigMapRef != nil && ef.ConfigMapRef.Name == extraVarsCmRef.ConfigMapRef.Name {
				hasCmEnvFrom = true
			}
			if ef.SecretRef != nil && ef.SecretRef.Name == extraVarsSecretRef.SecretRef.Name {
				hasSecretEnvFrom = true
			}
		}
		if !hasCmEnvFrom {
			container.EnvFrom = append(container.EnvFrom, extraVarsCmRef)
		}
		if !hasSecretEnvFrom {
			container.EnvFrom = append(container.EnvFrom, extraVarsSecretRef)
		}

		// Grafana home dashboard: in v4 this was configured via graf.Spec.ConfigMaps.
		// In v5 there is no ConfigMaps field, so we mount the ConfigMap explicitly
		// and rely on GF_DASHBOARDS_DEFAULT_HOME_DASHBOARD_PATH from extraVars.
		if cr.Spec.Grafana.GrafanaHomeDashboard {
			volumeName := "configmap-grafana-home-dashboard"
			mountPath := "/etc/grafana-configmaps/grafana-home-dashboard"

			// Ensure volume exists
			hasVolume := false
			for _, v := range podSpec.Volumes {
				if v.Name == volumeName {
					hasVolume = true
					break
				}
			}
			if !hasVolume {
				podSpec.Volumes = append(podSpec.Volumes, corev1.Volume{
					Name: volumeName,
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: "grafana-home-dashboard",
							},
						},
					},
				})
			}

			// Ensure volume mount exists on main container
			hasMount := false
			for _, vm := range container.VolumeMounts {
				if vm.Name == volumeName && vm.MountPath == mountPath {
					hasMount = true
					break
				}
			}
			if !hasMount {
				container.VolumeMounts = append(container.VolumeMounts, corev1.VolumeMount{
					Name:      volumeName,
					MountPath: mountPath,
					ReadOnly:  true,
				})
			}
		}

		// Grafana plugins init container: in v4 this was injected automatically by grafana-operator
		// via --grafana-plugins-init-container-image CLI arg. In v5 that feature was removed,
		// so we set up the init container explicitly via spec.deployment.spec.template.spec.initContainers.
		// The init container pre-installs bundled plugins into a shared emptyDir volume that is
		// then mounted into the main Grafana container at /var/lib/grafana/plugins.
		if cr.Spec.Grafana.Operator.InitContainerImage != "" {
			pluginsVolumeName := "grafana-plugins"
			pluginsInitMountPath := "/opt/plugins"
			pluginsMainMountPath := "/var/lib/grafana/plugins"

			// Ensure the shared plugins emptyDir volume exists
			hasPluginsVolume := false
			for _, v := range podSpec.Volumes {
				if v.Name == pluginsVolumeName {
					hasPluginsVolume = true
					break
				}
			}
			if !hasPluginsVolume {
				podSpec.Volumes = append(podSpec.Volumes, corev1.Volume{
					Name: pluginsVolumeName,
					VolumeSource: corev1.VolumeSource{
						EmptyDir: &corev1.EmptyDirVolumeSource{},
					},
				})
			}

			// Mount plugins volume in main Grafana container
			hasPluginsMount := false
			for _, vm := range container.VolumeMounts {
				if vm.Name == pluginsVolumeName {
					hasPluginsMount = true
					break
				}
			}
			if !hasPluginsMount {
				container.VolumeMounts = append(container.VolumeMounts, corev1.VolumeMount{
					Name:      pluginsVolumeName,
					MountPath: pluginsMainMountPath,
				})
			}

			// Add the plugins init container if not already present
			hasInitContainer := false
			for _, ic := range podSpec.InitContainers {
				if ic.Name == "grafana-plugins-init" {
					hasInitContainer = true
					break
				}
			}
			if !hasInitContainer {
				podSpec.InitContainers = append(podSpec.InitContainers, corev1.Container{
					Name:            "grafana-plugins-init",
					Image:           cr.Spec.Grafana.Operator.InitContainerImage,
					ImagePullPolicy: corev1.PullIfNotPresent,
					Env: []corev1.EnvVar{
						{
							Name:  "GRAFANA_PLUGINS",
							Value: "",
						},
					},
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      pluginsVolumeName,
							MountPath: pluginsInitMountPath,
						},
					},
				})
			}
		}

		// DashboardLabelSelector and DashboardNamespaceSelector removed or renamed in v5
		// Secrets removed or renamed in v5 - handle secrets differently if needed

		// TODO(#376): cr.Spec.Auth OAuth fields (LoginURL, TokenURL, UserInfoURL, TLSConfig) are not
		// yet applied to Grafana CR spec.config["auth.generic_oauth"]. In grafana-operator v4 these
		// were mapped to graf.Spec.Config.AuthGenericOauth. In v5 the target is
		// spec.config["auth.generic_oauth"][key] = value. Needs reimplementation.
		// TLS secrets (CASecret, CertSecret, KeySecret) also require volume mounts.
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
		// Note: We only set resources if containers already exist. We don't create containers here
		// because grafana-operator should manage container creation. If we create an empty container,
		// it will be missing required fields (name, image) and deployment will fail.
		if cr.Spec.Grafana.Resources.Size() > 0 {
			podSpec := ensurePodSpecInitialized(&graf)
			// Only set resources if container already exists (created by grafana-operator)
			// Don't create empty container - let grafana-operator handle container creation
			if len(podSpec.Containers) > 0 && podSpec.Containers[0].Name != "" {
				podSpec.Containers[0].Resources = cr.Spec.Grafana.Resources
			}
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

		// Mount secrets and configure Grafana to read sensitive values from files.
		// This avoids passing credentials via environment variables.
		{
			podSpec := ensurePodSpecInitialized(&graf)
			container := ensureGrafanaContainerInitialized(podSpec)

			const (
				adminSecretVolumeName = "grafana-admin-secret"
				adminSecretMountPath  = "/etc/grafana-admin"
				oauthSecretVolumeName = "grafana-oauth-secret"
				oauthSecretMountPath  = "/etc/grafana-oauth"
			)

			adminSecretName := fmt.Sprintf("%s-admin-credentials", graf.GetName())
			// optional=true when the user manages the secret; optional=false when Helm manages it.
			userManaged := cr.Spec.Grafana.DisableDefaultAdminSecret != nil && *cr.Spec.Grafana.DisableDefaultAdminSecret
			adminSecretOptional := userManaged

			// Ensure admin secret volume exists.
			hasAdminVolume := false
			for _, v := range podSpec.Volumes {
				if v.Name == adminSecretVolumeName {
					hasAdminVolume = true
					break
				}
			}
			if !hasAdminVolume {
				podSpec.Volumes = append(podSpec.Volumes, corev1.Volume{
					Name: adminSecretVolumeName,
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: adminSecretName,
							Optional:   &adminSecretOptional,
						},
					},
				})
			}

			// Ensure admin secret mount exists on main container.
			hasAdminMount := false
			for _, vm := range container.VolumeMounts {
				if vm.Name == adminSecretVolumeName && vm.MountPath == adminSecretMountPath {
					hasAdminMount = true
					break
				}
			}
			if !hasAdminMount {
				container.VolumeMounts = append(container.VolumeMounts, corev1.VolumeMount{
					Name:      adminSecretVolumeName,
					MountPath: adminSecretMountPath,
					ReadOnly:  true,
				})
			}

			// Configure Grafana admin credentials via file provider.
			securitySection := ensureGrafanaConfigSection(&graf, "security")
			securitySection["admin_user"] = "$__file{/etc/grafana-admin/GF_SECURITY_ADMIN_USER}"
			securitySection["admin_password"] = "$__file{/etc/grafana-admin/GF_SECURITY_ADMIN_PASSWORD}"

			// grafana-operator authenticates to Grafana via ServiceAccount JWT, not admin env vars.
			configureGrafanaOperatorKubeAuth(&graf, cr.GetNamespace())

			// Configure OAuth client secret via file provider when auth is enabled.
			// The secret is managed by Helm template oauth2-configs/secret-grafana-oauth-client-secret.yaml.
			if cr.Spec.Auth != nil {
				hasOAuthVolume := false
				for _, v := range podSpec.Volumes {
					if v.Name == oauthSecretVolumeName {
						hasOAuthVolume = true
						break
					}
				}
				if !hasOAuthVolume {
					optional := true
					podSpec.Volumes = append(podSpec.Volumes, corev1.Volume{
						Name: oauthSecretVolumeName,
						VolumeSource: corev1.VolumeSource{
							Secret: &corev1.SecretVolumeSource{
								SecretName: "grafana-oauth-client-secret",
								Optional:   &optional,
							},
						},
					})
				}

				hasOAuthMount := false
				for _, vm := range container.VolumeMounts {
					if vm.Name == oauthSecretVolumeName && vm.MountPath == oauthSecretMountPath {
						hasOAuthMount = true
						break
					}
				}
				if !hasOAuthMount {
					container.VolumeMounts = append(container.VolumeMounts, corev1.VolumeMount{
						Name:      oauthSecretVolumeName,
						MountPath: oauthSecretMountPath,
						ReadOnly:  true,
					})
				}

				authSection := ensureGrafanaConfigSection(&graf, "auth.generic_oauth")
				authSection["client_secret"] = "$__file{/etc/grafana-oauth/GF_AUTH_GENERIC_OAUTH_CLIENT_SECRET}"
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
		deployment := graf.Spec.Deployment
		if deployment != nil && deployment.Spec.Template != nil {
			if deployment.Spec.Template.Labels == nil {
				deployment.Spec.Template.Labels = make(map[string]string)
			}
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

			if deployment.Spec.Template.Annotations == nil {
				deployment.Spec.Template.Annotations = make(map[string]string)
			}
			if cr.Spec.Grafana.Annotations != nil {
				for k, v := range cr.Spec.Grafana.Annotations {
					deployment.Spec.Template.Annotations[k] = v
				}
			}
			// Propagate the admin secret checksum so grafana-operator triggers a rolling
			// restart of the Grafana Deployment when the secret changes.
			if adminPasswordChecksum != "" {
				deployment.Spec.Template.Annotations[adminSecretChecksumAnnotation] = adminPasswordChecksum
			}
		}

		// ServiceAccount in v5 uses different structure - Annotations and Labels may be in different location
		// Note: ServiceAccount configuration may need to be handled differently in v5
		// TODO: investigate ServiceAccountV1 structure in grafana-operator v5 and implement
		// propagation of cr.Spec.Grafana.ServiceAccount annotations/labels if still applicable.
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
	grafanaDatasourceInterval := "30s"
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
				dataSource.Spec.Datasource.URL = vmCluster.AsURL(vmetricsv1b1.ClusterComponentSelect) + "/select/0/prometheus"
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
	grafanaDatasourceInterval := "30s"
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

func grafanaIngressV1(cr *monv1.PlatformMonitoring) (*networkingv1.Ingress, error) {
	ingress := networkingv1.Ingress{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.GrafanaIngressAsset), 100).Decode(&ingress); err != nil {
		return nil, err
	}
	//Set metadata
	ingress.SetGroupVersionKind(schema.GroupVersionKind{Group: "networking.k8s.io", Version: "v1", Kind: "Ingress"})
	ingress.SetName(cr.GetNamespace() + "-" + utils.GrafanaComponentName)
	ingress.SetNamespace(cr.GetNamespace())

	if cr.Spec.Grafana != nil && cr.Spec.Grafana.Ingress != nil && cr.Spec.Grafana.Ingress.IsInstall() {
		var rules []networkingv1.IngressRule
		pathType := networkingv1.PathTypePrefix
		ing := cr.Spec.Grafana.Ingress

		switch {
		// 1. If custom ingress rules provided
		case len(ing.Rules) > 0:
			for _, r := range ing.Rules {
				// fallback if HTTP is not set
				if r.HTTP == nil || len(r.HTTP.Paths) == 0 {
					r.HTTP = &monv1.HTTPIngressRuleValue{
						Paths: []monv1.IngressPath{
							{
								Path:     "/",
								PathType: string(pathType),
								Backend: monv1.IngressPathBackend{
									Service: monv1.IngressPathBackendService{
										Name: utils.GrafanaServiceName,
										Port: monv1.ServiceBackendPort{
											Number: utils.GrafanaServicePort,
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

		// 2. If Host is provided
		case ing.Host != "":
			rules = append(rules, networkingv1.IngressRule{
				Host: ing.Host,
				IngressRuleValue: networkingv1.IngressRuleValue{
					HTTP: &networkingv1.HTTPIngressRuleValue{
						Paths: []networkingv1.HTTPIngressPath{defaultGrafanaPath(pathType)},
					},
				},
			})

		// 3. fallback: if no custom ingress rules or Host provided
		default:
			rules = append(rules, networkingv1.IngressRule{
				IngressRuleValue: networkingv1.IngressRuleValue{
					HTTP: &networkingv1.HTTPIngressRuleValue{
						Paths: []networkingv1.HTTPIngressPath{defaultGrafanaPath(pathType)},
					},
				},
			})
		}
		ingress.Spec.Rules = rules

		tlsConfigured := false
		// Configure tls if TLS config is defined
		if !tlsConfigured && len(cr.Spec.Grafana.Ingress.TLS) > 0 {
			for _, hostgroup := range cr.Spec.Grafana.Ingress.TLS {
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
				// fallback: if secretName is empty - use ingress TLSSecretName only
				secret := hostgroup.SecretName
				if secret == "" {
					secret = cr.Spec.Grafana.Ingress.TLSSecretName
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
		if !tlsConfigured && cr.Spec.Grafana.Ingress.Host != "" {
			secret := cr.Spec.Grafana.Ingress.TLSSecretName
			if secret != "" {
				ingress.Spec.TLS = []networkingv1.IngressTLS{
					{
						Hosts:      []string{cr.Spec.Grafana.Ingress.Host},
						SecretName: secret,
					},
				}
				tlsConfigured = true
			}
		}
		// Fallback: use ingress rules to configure tls hosts and TLSSecretName
		if !tlsConfigured && len(cr.Spec.Grafana.Ingress.Rules) > 0 {
			tlsHosts := []string{}
			secret := cr.Spec.Grafana.Ingress.TLSSecretName
			for _, rule := range cr.Spec.Grafana.Ingress.Rules {
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

		if cr.Spec.Grafana.Ingress.IngressClassName != nil {
			ingress.Spec.IngressClassName = cr.Spec.Grafana.Ingress.IngressClassName
		}

		// Set annotations
		ingress.SetAnnotations(utils.GetIngressAnnotationsForGateway(cr, cr.Spec.Grafana.Ingress.Annotations))

		// Set labels via centralized API (Ingress: base only per spec)
		in := utils.BaseOnlyLabelInput(ingress.GetName(), utils.GrafanaComponentName)
		if len(cr.Spec.Grafana.Ingress.Labels) > 0 {
			in.ComponentLabels = cr.Spec.Grafana.Ingress.Labels
		}
		utils.SetLabelsForResource(&ingress, in, nil)
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

	utils.SetLabelsForResource(&podMonitor, utils.LabelInput{
		Name:      podMonitor.GetName(),
		Component: utils.GrafanaComponentName,
		ComponentLabels: utils.MergeLabels(
			map[string]string{"app.kubernetes.io/processed-by-operator": "victoriametrics-operator"},
			cr.GetLabels(),
		),
	}, nil)

	return &podMonitor, nil
}

func defaultGrafanaPath(pathType networkingv1.PathType) networkingv1.HTTPIngressPath {
	return networkingv1.HTTPIngressPath{
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
	}
}
