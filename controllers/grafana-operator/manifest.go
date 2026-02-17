package grafana_operator

import (
	"embed"
	"fmt"
	"strings"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	grafv1 "github.com/grafana/grafana-operator/v5/api/v1beta1"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/yaml"
)

//go:embed  assets/*/*.yaml
var assets embed.FS

func grafanaOperatorServiceAccount(cr *monv1.PlatformMonitoring) (*corev1.ServiceAccount, error) {
	sa := corev1.ServiceAccount{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.GrafanaOperatorServiceAccountAsset), 100).Decode(&sa); err != nil {
		return nil, err
	}
	//Set parameters
	sa.SetGroupVersionKind(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "ServiceAccount"})
	sa.SetName(cr.GetNamespace() + "-" + utils.GrafanaOperatorComponentName)
	sa.SetNamespace(cr.GetNamespace())

	// Set annotations and labels for ServiceAccount in case
	if cr.Spec.Grafana != nil && cr.Spec.Grafana.Operator.ServiceAccount != nil {
		if sa.Annotations == nil && cr.Spec.Grafana.Operator.ServiceAccount.Annotations != nil {
			sa.SetAnnotations(cr.Spec.Grafana.Operator.ServiceAccount.Annotations)
		} else {
			for k, v := range cr.Spec.Grafana.Operator.ServiceAccount.Annotations {
				sa.Annotations[k] = v
			}
		}

		if cr.Spec.Grafana.Operator.ServiceAccount.Labels != nil {
			for k, v := range cr.Spec.Grafana.Operator.ServiceAccount.Labels {
				sa.Labels[k] = v
			}
		}
	}

	return &sa, nil
}

func grafanaOperatorClusterRole(cr *monv1.PlatformMonitoring) (*rbacv1.ClusterRole, error) {
	clusterRole := rbacv1.ClusterRole{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.GrafanaOperatorClusterRoleAsset), 100).Decode(&clusterRole); err != nil {
		return nil, err
	}
	//Set parameters
	clusterRole.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRole"})
	clusterRole.SetName(cr.GetNamespace() + "-" + utils.GrafanaOperatorComponentName)

	return &clusterRole, nil
}

func grafanaOperatorClusterRoleBinding(cr *monv1.PlatformMonitoring) (*rbacv1.ClusterRoleBinding, error) {
	clusterRoleBinding := rbacv1.ClusterRoleBinding{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.GrafanaOperatorClusterRoleBindingAsset), 100).Decode(&clusterRoleBinding); err != nil {
		return nil, err
	}
	//Set parameters
	clusterRoleBinding.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRoleBinding"})
	clusterRoleBinding.SetName(cr.GetNamespace() + "-" + utils.GrafanaOperatorComponentName)
	clusterRoleBinding.RoleRef.Name = cr.GetNamespace() + "-" + utils.GrafanaOperatorComponentName

	// Set namespace for all subjects
	for it := range clusterRoleBinding.Subjects {
		sub := &clusterRoleBinding.Subjects[it]
		sub.Namespace = cr.GetNamespace()
		sub.Name = cr.GetNamespace() + "-" + utils.GrafanaOperatorComponentName
	}
	return &clusterRoleBinding, nil
}

func grafanaOperatorRole(cr *monv1.PlatformMonitoring) (*rbacv1.Role, error) {
	role := rbacv1.Role{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.GrafanaOperatorRoleAsset), 100).Decode(&role); err != nil {
		return nil, err
	}
	//Set parameters
	role.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "Role"})
	role.SetName(cr.GetNamespace() + "-" + utils.GrafanaOperatorComponentName)
	role.SetNamespace(cr.GetNamespace())

	return &role, nil
}

func grafanaOperatorRoleBinding(cr *monv1.PlatformMonitoring) (*rbacv1.RoleBinding, error) {
	roleBinding := rbacv1.RoleBinding{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.GrafanaOperatorRoleBindingAsset), 100).Decode(&roleBinding); err != nil {
		return nil, err
	}
	//Set parameters
	roleBinding.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "RoleBinding"})
	roleBinding.SetName(cr.GetNamespace() + "-" + utils.GrafanaOperatorComponentName)
	roleBinding.SetNamespace(cr.GetNamespace())
	roleBinding.RoleRef.Name = cr.GetNamespace() + "-" + utils.GrafanaOperatorComponentName

	// Set namespace for all subjects
	for it := range roleBinding.Subjects {
		sub := &roleBinding.Subjects[it]
		sub.Namespace = cr.GetNamespace()
		sub.Name = cr.GetNamespace() + "-" + utils.GrafanaOperatorComponentName
	}
	return &roleBinding, nil
}

func grafanaOperatorDeployment(cr *monv1.PlatformMonitoring) (*appsv1.Deployment, error) {
	d := appsv1.Deployment{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.GrafanaOperatorDeploymentAsset), 100).Decode(&d); err != nil {
		return nil, err
	}
	//Set parameters
	d.SetGroupVersionKind(schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"})
	d.SetName(utils.GrafanaOperatorComponentName)
	d.SetNamespace(cr.GetNamespace())

	if cr.Spec.Grafana != nil {
		// Set correct images and any parameters to containers spec
		for it := range d.Spec.Template.Spec.Containers {
			c := &d.Spec.Template.Spec.Containers[it]
			if c.Name == utils.GrafanaOperatorComponentName {
				// Set grafana-operator image
				c.Image = cr.Spec.Grafana.Operator.Image

				// Remove command if present - grafana-operator v5 image has its own entrypoint
				// Setting command to /manager causes "stat /manager: no such file or directory" error
				c.Command = nil

				// Remove all already exist arguments from deployment asset
				c.Args = nil

				// Note: In grafana-operator v5, many command-line flags were removed:
				// - --grafana-image, --grafana-image-tag (removed - image is now specified in Grafana CR)
				// - --grafana-plugins-init-container-image, --grafana-plugins-init-container-tag (removed)
				// - --scan-all, --watch-namespaces, --namespace-scope (removed - use environment variables instead)
				// - --watch-namespace-selector, --watch-label-selectors (removed - use environment variables instead)
				// - --enforce-cache-labels, --cluster-domain (removed - use environment variables instead)
				// Only supported flags in v5:
				// - --max-concurrent-reconciles (supported via CR)
				// - --leader-elect (supported via CR)
				// - --zap-log-level (supported via CR)
				// - --zap-encoder, --zap-stacktrace-level, --zap-time-encoding, --zap-devel (not supported via CR - use deployment asset)
				// - --health-probe-bind-address, --metrics-bind-address, --pprof-addr (not supported via CR - use deployment asset)
				// - --kubeconfig (not supported via CR - use deployment asset)
				// Removed flags are now supported via environment variables:
				// - WATCH_NAMESPACE, WATCH_NAMESPACE_SELECTOR (officially supported, configured below)
				// Note: WATCH_LABEL_SELECTORS, ENFORCE_CACHE_LABELS, CLUSTER_DOMAIN are not supported in v5
				// These parameters are ignored if set in CR - they were removed in Grafana Operator v5
				// Add max concurrent reconciles if specified
				if cr.Spec.Grafana.Operator.MaxConcurrentReconciles != nil {
					c.Args = append(c.Args, fmt.Sprintf("--max-concurrent-reconciles=%d", *cr.Spec.Grafana.Operator.MaxConcurrentReconciles))
				}
				// Add leader elect if specified
				if cr.Spec.Grafana.Operator.LeaderElect != nil {
					c.Args = append(c.Args, fmt.Sprintf("--leader-elect=%t", *cr.Spec.Grafana.Operator.LeaderElect))
				}
				if cr.Spec.Grafana.Operator.LogLevel != "" {
					c.Args = append(c.Args, "--zap-log-level="+cr.Spec.Grafana.Operator.LogLevel)
				}
				// Set resources only if explicitly configured (not empty)
				// Size() checks if Resources has any requests or limits set
				// Empty Resources struct will have Size() == 0, preventing override of defaults
				if cr.Spec.Grafana.Operator.Resources.Size() > 0 {
					c.Resources = cr.Spec.Grafana.Operator.Resources
				}
				// Configure WATCH_NAMESPACE and WATCH_NAMESPACE_SELECTOR environment variables
				// Priority: NamespaceScope > WatchNamespaceSelector > WatchNamespaces > Namespaces (deprecated)
				// If none are set, WATCH_NAMESPACE will be empty (watch all namespaces) - filtering should be done via instanceSelector.matchExpressions
				// If NamespaceScope is true, limit to operator's namespace
				watchNamespace := ""
				watchNamespaceSelector := ""
				if cr.Spec.Grafana.Operator.NamespaceScope {
					watchNamespace = cr.GetNamespace()
				} else if cr.Spec.Grafana.Operator.WatchNamespaceSelector != "" {
					// Use label selector to dynamically discover namespaces
					watchNamespaceSelector = cr.Spec.Grafana.Operator.WatchNamespaceSelector
				} else if cr.Spec.Grafana.Operator.WatchNamespaces != "" {
					watchNamespace = cr.Spec.Grafana.Operator.WatchNamespaces
				} else if cr.Spec.Grafana.Operator.Namespaces != "" {
					// Support deprecated Namespaces field for backward compatibility
					watchNamespace = cr.Spec.Grafana.Operator.Namespaces
				}
				// If watchNamespace is still empty, operator will watch all namespaces (WATCH_NAMESPACE="")
				// In this case, filtering should be done via instanceSelector.matchExpressions on dashboards/datasources
				// Note: When WatchNamespaceSelector is used, watchNamespace remains empty ("") and WATCH_NAMESPACE_SELECTOR is set
				// This ensures only WATCH_NAMESPACE_SELECTOR is used, and WATCH_NAMESPACE from asset (if present) is overridden to ""
				// Update WATCH_NAMESPACE environment variable
				foundWatchNamespace := false
				for i := range c.Env {
					if c.Env[i].Name == "WATCH_NAMESPACE" {
						c.Env[i].Value = watchNamespace
						foundWatchNamespace = true
						break
					}
				}
				// Add WATCH_NAMESPACE if not found in deployment asset
				if !foundWatchNamespace {
					c.Env = append(c.Env, corev1.EnvVar{
						Name:  "WATCH_NAMESPACE",
						Value: watchNamespace,
					})
				}
				// Update or remove WATCH_NAMESPACE_SELECTOR environment variable
				// If watchNamespaceSelector is empty, remove the variable to ensure idempotent behavior
				foundWatchNamespaceSelector := false
				for i := range c.Env {
					if c.Env[i].Name == "WATCH_NAMESPACE_SELECTOR" {
						if watchNamespaceSelector != "" {
							c.Env[i].Value = watchNamespaceSelector
							foundWatchNamespaceSelector = true
						} else {
							// Remove WATCH_NAMESPACE_SELECTOR if it was set in asset but removed from CR
							c.Env = append(c.Env[:i], c.Env[i+1:]...)
						}
						break
					}
				}
				// Add WATCH_NAMESPACE_SELECTOR if not found and selector is set
				if !foundWatchNamespaceSelector && watchNamespaceSelector != "" {
					c.Env = append(c.Env, corev1.EnvVar{
						Name:  "WATCH_NAMESPACE_SELECTOR",
						Value: watchNamespaceSelector,
					})
				}
				// Note: WATCH_LABEL_SELECTORS, ENFORCE_CACHE_LABELS, CLUSTER_DOMAIN are not supported in Grafana Operator v5
				// These parameters (WatchLabelSelectors, EnforceCacheLabels, ClusterDomain) are ignored if set in CR
				break
			}
		}
		// Set pod-level security context
		if cr.Spec.Grafana.Operator.SecurityContext != nil {
			if d.Spec.Template.Spec.SecurityContext == nil {
				d.Spec.Template.Spec.SecurityContext = &corev1.PodSecurityContext{}
			}
			if cr.Spec.Grafana.Operator.SecurityContext.RunAsUser != nil {
				d.Spec.Template.Spec.SecurityContext.RunAsUser = cr.Spec.Grafana.Operator.SecurityContext.RunAsUser
			}
			if cr.Spec.Grafana.Operator.SecurityContext.RunAsGroup != nil {
				d.Spec.Template.Spec.SecurityContext.RunAsGroup = cr.Spec.Grafana.Operator.SecurityContext.RunAsGroup
			}
			if cr.Spec.Grafana.Operator.SecurityContext.FSGroup != nil {
				d.Spec.Template.Spec.SecurityContext.FSGroup = cr.Spec.Grafana.Operator.SecurityContext.FSGroup
			}
		}
		// Set tolerations for GrafanaOperator
		if cr.Spec.Grafana.Operator.Tolerations != nil {
			d.Spec.Template.Spec.Tolerations = cr.Spec.Grafana.Operator.Tolerations
		}
		// Set nodeSelector for GrafanaOperator
		if cr.Spec.Grafana.Operator.NodeSelector != nil {
			d.Spec.Template.Spec.NodeSelector = cr.Spec.Grafana.Operator.NodeSelector
		}
		// Set affinity for GrafanaOperator
		if cr.Spec.Grafana.Operator.Affinity != nil {
			d.Spec.Template.Spec.Affinity = cr.Spec.Grafana.Operator.Affinity
		}

		// Initialize labels and annotations maps if nil to prevent panic
		if d.Labels == nil {
			d.Labels = make(map[string]string)
		}
		if d.Annotations == nil {
			d.Annotations = make(map[string]string)
		}
		if d.Spec.Template.Labels == nil {
			d.Spec.Template.Labels = make(map[string]string)
		}
		if d.Spec.Template.Annotations == nil {
			d.Spec.Template.Annotations = make(map[string]string)
		}

		// Set labels
		d.Labels["app.kubernetes.io/name"] = utils.TruncLabel(d.GetName())
		d.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(d.GetName(), d.GetNamespace())
		d.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Grafana.Operator.Image)

		if cr.Spec.Grafana.Operator.Labels != nil {
			for k, v := range cr.Spec.Grafana.Operator.Labels {
				d.Labels[k] = v
			}
		}

		if cr.Spec.Grafana.Operator.Annotations != nil {
			for k, v := range cr.Spec.Grafana.Operator.Annotations {
				d.Annotations[k] = v
			}
		}

		// Set template labels
		d.Spec.Template.Labels["app.kubernetes.io/name"] = utils.TruncLabel(d.GetName())
		d.Spec.Template.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(d.GetName(), d.GetNamespace())
		d.Spec.Template.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Grafana.Operator.Image)

		if cr.Spec.Grafana.Operator.Labels != nil {
			for k, v := range cr.Spec.Grafana.Operator.Labels {
				d.Spec.Template.Labels[k] = v
			}
		}

		if cr.Spec.Grafana.Operator.Annotations != nil {
			for k, v := range cr.Spec.Grafana.Operator.Annotations {
				d.Spec.Template.Annotations[k] = v
			}
		}
		if len(strings.TrimSpace(cr.Spec.Grafana.Operator.PriorityClassName)) > 0 {
			d.Spec.Template.Spec.PriorityClassName = cr.Spec.Grafana.Operator.PriorityClassName
		}
	}
	d.Spec.Template.Spec.ServiceAccountName = cr.GetNamespace() + "-" + utils.GrafanaOperatorComponentName
	return &d, nil
}

func grafanaDashboard(cr *monv1.PlatformMonitoring, fileName string) (*grafv1.GrafanaDashboard, error) {
	dashboard := grafv1.GrafanaDashboard{}
	fullPath := utils.BasePath + utils.DashboardsFolder + fileName
	crParams := cr.ToParams()
	// Add a map contains human-readable UIDs for Grafana dashboards to the current Custom Resource
	crParams.DashboardsUIDs = utils.DashboardsUIDsMap
	fileContent, err := utils.ParseTemplate(utils.MustAssetReaderToString(assets, fullPath), fullPath, utils.DashboardTemplateLeftDelim, utils.DashboardTemplateRightDelim, crParams)
	if err != nil {
		return nil, err
	}
	if err := yaml.NewYAMLOrJSONDecoder(strings.NewReader(fileContent), 100).Decode(&dashboard); err != nil {
		return nil, err
	}
	// Set parameters
	// Explicitly set GVK to ensure correct API group (grafana.integreatly.org/v1beta1) is used
	// This is required for Grafana Operator v5 migration from integreatly.org/v1alpha1
	dashboard.SetGroupVersionKind(schema.GroupVersionKind{Group: "grafana.integreatly.org", Version: "v1beta1", Kind: "GrafanaDashboard"})
	dashboard.SetNamespace(cr.GetNamespace())

	// Initialize labels map if nil to prevent panic
	if dashboard.Labels == nil {
		dashboard.Labels = make(map[string]string)
	}

	// Set labels
	dashboard.Labels["name"] = utils.TruncLabel(dashboard.GetName())
	dashboard.Labels["app.kubernetes.io/name"] = utils.TruncLabel(dashboard.GetName())
	dashboard.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(dashboard.GetName(), dashboard.GetNamespace())
	dashboard.Labels["app.kubernetes.io/part-of"] = "monitoring"
	dashboard.Labels["app.kubernetes.io/managed-by"] = "monitoring-operator"
	if cr.Spec.Grafana != nil {
		dashboard.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Grafana.Operator.Image)
	}

	// Configure instanceSelector dynamically based on Grafana labels
	// In grafana-operator v5, instanceSelector is optional - if not specified, dashboard applies to all Grafana instances in namespace
	// If instanceSelector is already set in dashboard YAML, we can override it with labels from Grafana resource
	// This allows using custom labels set via cr.Spec.Grafana.Labels for instanceSelector matching
	if cr.Spec.Grafana != nil {
		// Build instanceSelector matchLabels from Grafana labels
		// Use default labels if custom labels are not set
		instanceLabels := make(map[string]string)

		// Get component label - use custom if set, otherwise default
		if cr.Spec.Grafana.Labels != nil && cr.Spec.Grafana.Labels["app.kubernetes.io/component"] != "" {
			instanceLabels["app.kubernetes.io/component"] = cr.Spec.Grafana.Labels["app.kubernetes.io/component"]
		} else {
			instanceLabels["app.kubernetes.io/component"] = "grafana"
		}

		// Get part-of label - use custom if set, otherwise default
		if cr.Spec.Grafana.Labels != nil && cr.Spec.Grafana.Labels["app.kubernetes.io/part-of"] != "" {
			instanceLabels["app.kubernetes.io/part-of"] = cr.Spec.Grafana.Labels["app.kubernetes.io/part-of"]
		} else {
			instanceLabels["app.kubernetes.io/part-of"] = "monitoring"
		}

		// Add any other custom labels that might be useful for instanceSelector
		// Users can set custom labels like "dashboards: custom" and use them in instanceSelector
		if cr.Spec.Grafana.Labels != nil {
			for k, v := range cr.Spec.Grafana.Labels {
				// Skip standard labels that are already set above
				if k != "app.kubernetes.io/component" && k != "app.kubernetes.io/part-of" &&
					k != "app.kubernetes.io/instance" && k != "app.kubernetes.io/version" &&
					k != "app.kubernetes.io/name" && k != "app.kubernetes.io/managed-by" {
					instanceLabels[k] = v
				}
			}
		}

		// Configure instanceSelector for dashboard
		// In v5, if instanceSelector is not set, dashboard applies to ALL Grafana instances in namespace
		// We only update instanceSelector if it was already set in dashboard YAML (vanilla v5 behavior)
		// This preserves the ability to have dashboards apply to all Grafana instances by omitting instanceSelector
		if dashboard.Spec.InstanceSelector != nil {
			// Update matchLabels from Grafana resource
			// Note: matchExpressions (if present in YAML) are preserved and combined with matchLabels via AND logic
			// This allows combining labels from Grafana resource with expressions from dashboard YAML
			// Example: matchLabels from Grafana + matchExpressions from YAML = both conditions must be met
			// If you need to use only matchLabels from Grafana, ensure matchExpressions is not set in dashboard YAML
			dashboard.Spec.InstanceSelector.MatchLabels = make(map[string]string)
			// Set labels from Grafana resource
			for k, v := range instanceLabels {
				dashboard.Spec.InstanceSelector.MatchLabels[k] = v
			}
			// matchExpressions are preserved if they exist in YAML (not modified)
		}
		// If instanceSelector is nil, dashboard will apply to all Grafana instances in namespace (vanilla v5 behavior)
	}

	return &dashboard, nil
}

func grafanaOperatorPodMonitor(cr *monv1.PlatformMonitoring) (*promv1.PodMonitor, error) {
	podMonitor := promv1.PodMonitor{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.GrafanaOperatorPodMonitorAsset), 100).Decode(&podMonitor); err != nil {
		return nil, err
	}
	//Set parameters
	podMonitor.SetGroupVersionKind(schema.GroupVersionKind{Group: "monitoring.coreos.com", Version: "v1", Kind: "PodMonitor"})
	podMonitor.SetName(cr.GetNamespace() + "-" + "grafana-operator-pod-monitor")
	podMonitor.SetNamespace(cr.GetNamespace())

	if cr.Spec.Grafana != nil && cr.Spec.Grafana.Operator.PodMonitor != nil && cr.Spec.Grafana.Operator.PodMonitor.IsInstall() {
		cr.Spec.Grafana.Operator.PodMonitor.OverridePodMonitor(&podMonitor)
	}
	return &podMonitor, nil
}
