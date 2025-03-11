package prometheus_operator

import (
	"embed"
	"fmt"
	"strings"

	v1alpha1 "github.com/Netcracker/qubership-monitoring-operator/api/v1alpha1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/yaml"
)

//go:embed  assets/*.yaml
var assets embed.FS

// PrometheusOperatorRole builds the ServiceAccount resource manifest
// and fill it with parameters from the CR.
func prometheusOperatorRole(cr *v1alpha1.PlatformMonitoring) (*rbacv1.Role, error) {
	role := rbacv1.Role{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.PrometheusOperatorRoleAsset), 100).Decode(&role); err != nil {
		return nil, err
	}
	//Set parameters
	role.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "Role"})
	role.SetName(cr.GetNamespace() + "-" + utils.PrometheusOperatorComponentName)
	role.SetNamespace(cr.GetNamespace())

	return &role, nil
}

func prometheusOperatorServiceAccount(cr *v1alpha1.PlatformMonitoring) (*corev1.ServiceAccount, error) {
	sa := corev1.ServiceAccount{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.PrometheusOperatorServiceAccountAsset), 100).Decode(&sa); err != nil {
		return nil, err
	}
	//Set parameters
	sa.SetGroupVersionKind(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "ServiceAccount"})
	sa.SetName(cr.GetNamespace() + "-" + utils.PrometheusOperatorComponentName)
	sa.SetNamespace(cr.GetNamespace())

	// Set annotations and labels for ServiceAccount in case
	if cr.Spec.Prometheus != nil {
		if cr.Spec.Prometheus.Operator.ServiceAccount != nil {
			if sa.Annotations == nil && cr.Spec.Prometheus.ServiceAccount.Annotations != nil {
				sa.SetAnnotations(cr.Spec.Prometheus.Operator.ServiceAccount.Annotations)
			} else {
				for k, v := range cr.Spec.Prometheus.Operator.ServiceAccount.Annotations {
					sa.Annotations[k] = v
				}
			}

			if sa.Labels == nil && cr.Spec.Prometheus.Operator.ServiceAccount.Labels != nil {
				sa.SetLabels(cr.Spec.Prometheus.Operator.ServiceAccount.Labels)
			} else {
				for k, v := range cr.Spec.Prometheus.Operator.ServiceAccount.Labels {
					sa.Labels[k] = v
				}
			}
		}
	}

	return &sa, nil
}

func prometheusOperatorRoleBinding(cr *v1alpha1.PlatformMonitoring) (*rbacv1.RoleBinding, error) {
	roleBinding := rbacv1.RoleBinding{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.PrometheusOperatorRoleBindingAsset), 100).Decode(&roleBinding); err != nil {
		return nil, err
	}
	//Set parameters
	roleBinding.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "RoleBinding"})
	roleBinding.SetName(cr.GetNamespace() + "-" + utils.PrometheusOperatorComponentName)
	roleBinding.SetNamespace(cr.GetNamespace())
	roleBinding.RoleRef.Name = cr.GetNamespace() + "-" + utils.PrometheusOperatorComponentName

	// Set namespace for all subjects
	for it := range roleBinding.Subjects {
		sub := &roleBinding.Subjects[it]
		sub.Namespace = cr.GetNamespace()
		sub.Name = cr.GetNamespace() + "-" + utils.PrometheusOperatorComponentName
	}
	return &roleBinding, nil
}

func prometheusOperatorClusterRole(cr *v1alpha1.PlatformMonitoring) (*rbacv1.ClusterRole, error) {
	clusterRole := rbacv1.ClusterRole{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.PrometheusOperatorClusterRoleAsset), 100).Decode(&clusterRole); err != nil {
		return nil, err
	}
	//Set parameters
	clusterRole.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRole"})
	clusterRole.SetName(cr.GetNamespace() + "-" + utils.PrometheusOperatorComponentName)

	// Add permissions to kubelet service
	clusterRole.Rules = append(clusterRole.Rules, rbacv1.PolicyRule{
		APIGroups: []string{""},
		Resources: []string{"services", "services/finalizers", "endpoints"},
		Verbs:     []string{"get", "create", "list", "update", "watch"},
	})
	return &clusterRole, nil
}

func prometheusOperatorClusterRoleBinding(cr *v1alpha1.PlatformMonitoring) (*rbacv1.ClusterRoleBinding, error) {
	clusterRoleBinding := rbacv1.ClusterRoleBinding{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.PrometheusOperatorClusterRoleBindingAsset), 100).Decode(&clusterRoleBinding); err != nil {
		return nil, err
	}
	//Set parameters
	clusterRoleBinding.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRoleBinding"})
	clusterRoleBinding.SetName(cr.GetNamespace() + "-" + utils.PrometheusOperatorComponentName)
	clusterRoleBinding.RoleRef.Name = cr.GetNamespace() + "-" + utils.PrometheusOperatorComponentName

	// Set namespace for all subjects
	for it := range clusterRoleBinding.Subjects {
		sub := &clusterRoleBinding.Subjects[it]
		sub.Namespace = cr.GetNamespace()
		sub.Name = cr.GetNamespace() + "-" + utils.PrometheusOperatorComponentName
	}
	return &clusterRoleBinding, nil
}

func prometheusOperatorDeployment(cr *v1alpha1.PlatformMonitoring) (*appsv1.Deployment, error) {
	d := appsv1.Deployment{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.PrometheusOperatorDeploymentAsset), 100).Decode(&d); err != nil {
		return nil, err
	}
	//Set parameters
	d.SetGroupVersionKind(schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"})
	d.SetName(utils.PrometheusOperatorComponentName)
	d.SetNamespace(cr.GetNamespace())

	if cr.Spec.Prometheus != nil {

		// Set correct images and any parameters to containers spec
		for it := range d.Spec.Template.Spec.Containers {
			c := &d.Spec.Template.Spec.Containers[it]
			if c.Name == utils.PrometheusOperatorComponentName {
				// Set prometheus-operator image
				c.Image = cr.Spec.Prometheus.Operator.Image

				// Set prometheus config reloader image
				c.Args = append(c.Args, "--prometheus-config-reloader="+cr.Spec.Prometheus.ConfigReloaderImage)
				c.Args = append(c.Args, "--prometheus-instance-namespaces="+cr.GetNamespace())
				c.Args = append(c.Args, "--alertmanager-instance-namespaces="+cr.GetNamespace())
				c.Args = append(c.Args, fmt.Sprintf("--kubelet-service=%s/kubelet", cr.GetNamespace()))
				if cr.Spec.Prometheus.Operator.Namespaces != "" {
					c.Args = append(c.Args, "--namespaces="+cr.Spec.Prometheus.Operator.Namespaces)
				}
				if cr.Spec.Prometheus.Operator.Resources.Size() > 0 {
					c.Resources = cr.Spec.Prometheus.Operator.Resources
				}
				break
			}
		}
		// Set security context
		if cr.Spec.Prometheus.Operator.SecurityContext != nil {
			if d.Spec.Template.Spec.SecurityContext == nil {
				d.Spec.Template.Spec.SecurityContext = &corev1.PodSecurityContext{}
			}
			if cr.Spec.Prometheus.Operator.SecurityContext.RunAsUser != nil {
				d.Spec.Template.Spec.SecurityContext.RunAsUser = cr.Spec.Prometheus.Operator.SecurityContext.RunAsUser
			}
			if cr.Spec.Prometheus.Operator.SecurityContext.FSGroup != nil {
				d.Spec.Template.Spec.SecurityContext.FSGroup = cr.Spec.Prometheus.Operator.SecurityContext.FSGroup
			}
		}
		// Set tolerations for PrometheusOperator
		if cr.Spec.Prometheus.Operator.Tolerations != nil {
			d.Spec.Template.Spec.Tolerations = cr.Spec.Prometheus.Operator.Tolerations
		}
		// Set nodeSelector for PrometheusOperator
		if cr.Spec.Prometheus.Operator.NodeSelector != nil {
			d.Spec.Template.Spec.NodeSelector = cr.Spec.Prometheus.Operator.NodeSelector
		}
		// Set affinity for PrometheusOperator
		if cr.Spec.Prometheus.Operator.Affinity != nil {
			d.Spec.Template.Spec.Affinity = cr.Spec.Prometheus.Operator.Affinity
		}

		// Set labels
		d.Labels["name"] = utils.TruncLabel(d.GetName())
		d.Labels["app.kubernetes.io/name"] = utils.TruncLabel(d.GetName())
		d.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(d.GetName(), d.GetNamespace())
		d.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Prometheus.Operator.Image)

		if cr.Spec.Prometheus.Operator.Labels != nil {
			for k, v := range cr.Spec.Prometheus.Operator.Labels {
				d.Labels[k] = v
			}
		}

		if d.Annotations == nil && cr.Spec.Prometheus.Operator.Annotations != nil {
			d.SetAnnotations(cr.Spec.Prometheus.Operator.Annotations)
		} else {
			for k, v := range cr.Spec.Prometheus.Operator.Annotations {
				d.Annotations[k] = v
			}
		}

		// Set labels
		d.Spec.Template.Labels["name"] = utils.TruncLabel(d.GetName())
		d.Spec.Template.Labels["app.kubernetes.io/name"] = utils.TruncLabel(d.GetName())
		d.Spec.Template.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(d.GetName(), d.GetNamespace())
		d.Spec.Template.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Prometheus.Operator.Image)

		if cr.Spec.Prometheus.Operator.Labels != nil {
			for k, v := range cr.Spec.Prometheus.Operator.Labels {
				d.Spec.Template.Labels[k] = v
			}
		}

		if d.Spec.Template.Annotations == nil && cr.Spec.Prometheus.Operator.Annotations != nil {
			d.Spec.Template.Annotations = cr.Spec.Prometheus.Operator.Annotations
		} else {
			for k, v := range cr.Spec.Prometheus.Operator.Annotations {
				d.Spec.Template.Annotations[k] = v
			}
		}
		if len(strings.TrimSpace(cr.Spec.Prometheus.Operator.PriorityClassName)) > 0 {
			d.Spec.Template.Spec.PriorityClassName = cr.Spec.Prometheus.Operator.PriorityClassName
		}
	}
	d.Spec.Template.Spec.ServiceAccountName = cr.GetNamespace() + "-" + utils.PrometheusOperatorComponentName

	return &d, nil
}

func prometheusOperatorService(cr *v1alpha1.PlatformMonitoring) (*corev1.Service, error) {
	service := corev1.Service{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.PrometheusOperatorServiceAsset), 100).Decode(&service); err != nil {
		return nil, err
	}
	//Set parameters
	service.SetGroupVersionKind(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Service"})
	service.SetName(utils.PrometheusOperatorComponentName)
	service.SetNamespace(cr.GetNamespace())

	return &service, nil
}

func prometheusOperatorPodMonitor(cr *v1alpha1.PlatformMonitoring) (*promv1.PodMonitor, error) {
	podMonitor := promv1.PodMonitor{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.PrometheusOperatorPodMonitorAsset), 100).Decode(&podMonitor); err != nil {
		return nil, err
	}
	//Set parameters
	podMonitor.SetGroupVersionKind(schema.GroupVersionKind{Group: "monitoring.coreos.com", Version: "v1", Kind: "PodMonitor"})
	podMonitor.SetName(cr.GetNamespace() + "-" + "prometheus-operator-pod-monitor")
	podMonitor.SetNamespace(cr.GetNamespace())

	if cr.Spec.Prometheus != nil && cr.Spec.Prometheus.Operator.PodMonitor != nil && cr.Spec.Prometheus.Operator.PodMonitor.IsInstall() {
		cr.Spec.Prometheus.Operator.PodMonitor.OverridePodMonitor(&podMonitor)
	}
	return &podMonitor, nil
}
