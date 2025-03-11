package grafana_operator

import (
	"embed"
	"strings"

	v1alpha1 "github.com/Netcracker/qubership-monitoring-operator/api/v1alpha1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	grafv1 "github.com/grafana-operator/grafana-operator/v4/api/integreatly/v1alpha1"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/yaml"
)

//go:embed  assets/*/*.yaml
var assets embed.FS

func grafanaOperatorServiceAccount(cr *v1alpha1.PlatformMonitoring) (*corev1.ServiceAccount, error) {
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

func grafanaOperatorClusterRole(cr *v1alpha1.PlatformMonitoring) (*rbacv1.ClusterRole, error) {
	clusterRole := rbacv1.ClusterRole{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.GrafanaOperatorClusterRoleAsset), 100).Decode(&clusterRole); err != nil {
		return nil, err
	}
	//Set parameters
	clusterRole.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRole"})
	clusterRole.SetName(cr.GetNamespace() + "-" + utils.GrafanaOperatorComponentName)

	return &clusterRole, nil
}

func grafanaOperatorClusterRoleBinding(cr *v1alpha1.PlatformMonitoring) (*rbacv1.ClusterRoleBinding, error) {
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

func grafanaOperatorRole(cr *v1alpha1.PlatformMonitoring) (*rbacv1.Role, error) {
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

func grafanaOperatorRoleBinding(cr *v1alpha1.PlatformMonitoring) (*rbacv1.RoleBinding, error) {
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

func grafanaOperatorDeployment(cr *v1alpha1.PlatformMonitoring) (*appsv1.Deployment, error) {
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

				// Remove all already exist arguments from deployment asset
				c.Args = nil

				grafanaImage, grafanaTag := utils.SplitImage(cr.Spec.Grafana.Image)

				// Set config reloader image
				c.Args = append(c.Args, "--grafana-image="+grafanaImage)

				// Set prometheus config reloader image
				if grafanaTag != "" {
					c.Args = append(c.Args, "--grafana-image-tag="+grafanaTag)
				}

				initContainerImage, initContainerTag := utils.SplitImage(cr.Spec.Grafana.Operator.InitContainerImage)

				// Set prometheus config reloader image
				c.Args = append(c.Args, "--grafana-plugins-init-container-image="+initContainerImage)

				// Set prometheus config reloader image
				if initContainerTag != "" {
					c.Args = append(c.Args, "--grafana-plugins-init-container-tag="+initContainerTag)
				}
				// Set flag for scan namespaces
				if cr.Spec.Grafana.Operator.Namespaces != "" {
					c.Args = append(c.Args, "--namespaces="+cr.Spec.Grafana.Operator.Namespaces)
				} else {
					c.Args = append(c.Args, "--scan-all")
				}
				if cr.Spec.Grafana.Operator.LogLevel != "" {
					c.Args = append(c.Args, "--zap-log-level="+cr.Spec.Grafana.Operator.LogLevel)
				}
				if cr.Spec.Grafana.Operator.Resources.Size() > 0 {
					c.Resources = cr.Spec.Grafana.Operator.Resources
				}
				break
			}
		}
		// Set security context
		if cr.Spec.Grafana.Operator.SecurityContext != nil {
			if d.Spec.Template.Spec.SecurityContext == nil {
				d.Spec.Template.Spec.SecurityContext = &corev1.PodSecurityContext{}
			}
			if cr.Spec.Grafana.Operator.SecurityContext.RunAsUser != nil {
				d.Spec.Template.Spec.SecurityContext.RunAsUser = cr.Spec.Grafana.Operator.SecurityContext.RunAsUser
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

		// Set labels
		d.Labels["app.kubernetes.io/name"] = utils.TruncLabel(d.GetName())
		d.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(d.GetName(), d.GetNamespace())
		d.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Grafana.Operator.Image)

		if cr.Spec.Grafana.Operator.Labels != nil {
			for k, v := range cr.Spec.Grafana.Operator.Labels {
				d.Labels[k] = v
			}
		}

		if d.Annotations == nil && cr.Spec.Grafana.Operator.Annotations != nil {
			d.SetAnnotations(cr.Spec.Grafana.Operator.Annotations)
		} else {
			for k, v := range cr.Spec.Grafana.Operator.Annotations {
				d.Annotations[k] = v
			}
		}

		// Set labels
		d.Spec.Template.Labels["app.kubernetes.io/name"] = utils.TruncLabel(d.GetName())
		d.Spec.Template.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(d.GetName(), d.GetNamespace())
		d.Spec.Template.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Grafana.Operator.Image)

		if cr.Spec.Grafana.Operator.Labels != nil {
			for k, v := range cr.Spec.Grafana.Operator.Labels {
				d.Spec.Template.Labels[k] = v
			}
		}

		if d.Spec.Template.Annotations == nil && cr.Spec.Grafana.Operator.Annotations != nil {
			d.Spec.Template.Annotations = cr.Spec.Grafana.Operator.Annotations
		} else {
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

func grafanaDashboard(cr *v1alpha1.PlatformMonitoring, fileName string) (*grafv1.GrafanaDashboard, error) {
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
	//Set parameters
	dashboard.SetGroupVersionKind(schema.GroupVersionKind{Group: "integreatly.org", Version: "v1alpha1", Kind: "Grafana"})
	dashboard.SetNamespace(cr.GetNamespace())

	// Set labels
	dashboard.Labels["name"] = utils.TruncLabel(dashboard.GetName())
	dashboard.Labels["app.kubernetes.io/name"] = utils.TruncLabel(dashboard.GetName())
	dashboard.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(dashboard.GetName(), dashboard.GetNamespace())
	dashboard.Labels["app.kubernetes.io/part-of"] = "monitoring"
	dashboard.Labels["app.kubernetes.io/managed-by"] = "monitoring-operator"
	if cr.Spec.Grafana != nil {
		dashboard.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Grafana.Operator.Image)
	}

	return &dashboard, nil
}

func grafanaOperatorPodMonitor(cr *v1alpha1.PlatformMonitoring) (*promv1.PodMonitor, error) {
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
