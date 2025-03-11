package kubestatemetrics

import (
	"embed"
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

func kubeStateMetricsServiceAccount(cr *v1alpha1.PlatformMonitoring) (*corev1.ServiceAccount, error) {
	sa := corev1.ServiceAccount{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.KubestatemetricsServiceAccountAsset), 100).Decode(&sa); err != nil {
		return nil, err
	}
	//Set parameters
	sa.SetGroupVersionKind(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "ServiceAccount"})
	sa.SetName(cr.GetNamespace() + "-" + utils.KubestatemetricsComponentName)
	sa.SetNamespace(cr.GetNamespace())

	// Set annotations and labels for ServiceAccount in case
	if cr.Spec.KubeStateMetrics != nil && cr.Spec.KubeStateMetrics.ServiceAccount != nil {
		if sa.Annotations == nil && cr.Spec.KubeStateMetrics.ServiceAccount.Annotations != nil {
			sa.SetAnnotations(cr.Spec.KubeStateMetrics.ServiceAccount.Annotations)
		} else {
			for k, v := range cr.Spec.KubeStateMetrics.ServiceAccount.Annotations {
				sa.Annotations[k] = v
			}
		}

		if sa.Labels == nil && cr.Spec.KubeStateMetrics.ServiceAccount.Labels != nil {
			sa.SetLabels(cr.Spec.KubeStateMetrics.ServiceAccount.Labels)
		} else {
			for k, v := range cr.Spec.KubeStateMetrics.ServiceAccount.Labels {
				sa.Labels[k] = v
			}
		}
	}

	return &sa, nil
}

func kubeStateMetricsClusterRole(cr *v1alpha1.PlatformMonitoring) (*rbacv1.ClusterRole, error) {
	clusterRole := rbacv1.ClusterRole{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.KubestatemetricsClusterRoleAsset), 100).Decode(&clusterRole); err != nil {
		return nil, err
	}
	//Set parameters
	clusterRole.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRole"})
	clusterRole.SetName(cr.GetNamespace() + "-" + utils.KubestatemetricsComponentName)

	return &clusterRole, nil
}

func kubeStateMetricsClusterRoleBinding(cr *v1alpha1.PlatformMonitoring) (*rbacv1.ClusterRoleBinding, error) {
	clusterRoleBinding := rbacv1.ClusterRoleBinding{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.KubestatemetricsClusterRoleBindingAsset), 100).Decode(&clusterRoleBinding); err != nil {
		return nil, err
	}
	//Set parameters
	clusterRoleBinding.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRoleBinding"})
	clusterRoleBinding.SetName(cr.GetNamespace() + "-" + utils.KubestatemetricsComponentName)
	clusterRoleBinding.RoleRef.Name = cr.GetNamespace() + "-" + utils.KubestatemetricsComponentName

	// Set namespace and name for all subjects
	for it := range clusterRoleBinding.Subjects {
		sub := &clusterRoleBinding.Subjects[it]
		sub.Name = cr.GetNamespace() + "-" + utils.KubestatemetricsComponentName
		sub.Namespace = cr.GetNamespace()
	}
	return &clusterRoleBinding, nil
}

func kubeStateMetricsDeployment(cr *v1alpha1.PlatformMonitoring, hasIngress bool) (*appsv1.Deployment, error) {
	d := appsv1.Deployment{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.KubestatemetricsDeploymentAsset), 100).Decode(&d); err != nil {
		return nil, err
	}
	//Set parameters
	d.SetGroupVersionKind(schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"})
	d.SetName(utils.KubestatemetricsComponentName)
	d.SetNamespace(cr.GetNamespace())

	if cr.Spec.KubeStateMetrics != nil {
		// Find container with b.Name as name and set Image from custom resource
		for it := range d.Spec.Template.Spec.Containers {
			c := &d.Spec.Template.Spec.Containers[it]
			if c.Name == utils.KubestatemetricsComponentName {
				c.Image = cr.Spec.KubeStateMetrics.Image
				if cr.Spec.KubeStateMetrics.Resources.Size() > 0 {
					c.Resources = cr.Spec.KubeStateMetrics.Resources
				}
				// Set collectors and namespaces
				if !utils.PrivilegedRights {
					if cr.Spec.KubeStateMetrics.ScrapeResources != "" {
						c.Args = append(c.Args, "--resources="+cr.Spec.KubeStateMetrics.ScrapeResources)
					} else {
						if hasIngress {
							c.Args = append(c.Args, "--resources="+utils.ScrapeResources+",horizontalpodautoscalers,ingresses")
						} else {
							c.Args = append(c.Args, "--resources="+utils.ScrapeResources)
						}
					}
					if cr.Spec.KubeStateMetrics.Namespaces != "" {
						c.Args = append(c.Args, "--namespaces="+cr.Spec.KubeStateMetrics.Namespaces)
					}
				} else {
					if cr.Spec.KubeStateMetrics.ScrapeResources != "" {
						c.Args = append(c.Args, "--resources="+cr.Spec.KubeStateMetrics.ScrapeResources)
					} else {
						if hasIngress {
							c.Args = append(c.Args, "--resources="+utils.ScrapeResources+",secrets,storageclasses,horizontalpodautoscalers,ingresses,networkpolicies")
						} else {
							c.Args = append(c.Args, "--resources="+utils.ScrapeResources+",secrets,storageclasses,networkpolicies")
						}
					}
				}
				if cr.Spec.KubeStateMetrics.MetricLabelsAllowlist != "" {
					c.Args = append(c.Args, "--metric-labels-allowlist="+cr.Spec.KubeStateMetrics.MetricLabelsAllowlist)
				} else {
					c.Args = append(c.Args, "--metric-labels-allowlist=nodes=[*],pods=[*],namespaces=[*],deployments=[*],statefulsets=[*],daemonsets=[*],cronjobs=[*],jobs=[*],ingresses=[*],services=[*]")
				}
				break
			}
		}
		// Set security context
		if cr.Spec.KubeStateMetrics.SecurityContext != nil {
			if d.Spec.Template.Spec.SecurityContext == nil {
				d.Spec.Template.Spec.SecurityContext = &corev1.PodSecurityContext{}
			}
			if cr.Spec.KubeStateMetrics.SecurityContext.RunAsUser != nil {
				d.Spec.Template.Spec.SecurityContext.RunAsUser = cr.Spec.KubeStateMetrics.SecurityContext.RunAsUser
			}
			if cr.Spec.KubeStateMetrics.SecurityContext.FSGroup != nil {
				d.Spec.Template.Spec.SecurityContext.FSGroup = cr.Spec.KubeStateMetrics.SecurityContext.FSGroup
			}
		}
		// Set tolerations for KubeStateMetrics
		if cr.Spec.KubeStateMetrics.Tolerations != nil {
			d.Spec.Template.Spec.Tolerations = cr.Spec.KubeStateMetrics.Tolerations
		}
		// Set nodeSelector for KubeStateMetrics
		if cr.Spec.KubeStateMetrics.NodeSelector != nil {
			d.Spec.Template.Spec.NodeSelector = cr.Spec.KubeStateMetrics.NodeSelector
		}
		// Set affinity for KubeStateMetrics
		if cr.Spec.KubeStateMetrics.Affinity != nil {
			d.Spec.Template.Spec.Affinity = cr.Spec.KubeStateMetrics.Affinity
		}

		// Set labels
		d.Labels["name"] = utils.TruncLabel(d.GetName())
		d.Labels["app.kubernetes.io/name"] = utils.TruncLabel(d.GetName())
		d.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(d.GetName(), d.GetNamespace())
		d.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.KubeStateMetrics.Image)

		if cr.Spec.KubeStateMetrics.Labels != nil {
			for k, v := range cr.Spec.KubeStateMetrics.Labels {
				d.Labels[k] = v
			}
		}

		if d.Annotations == nil && cr.Spec.KubeStateMetrics.Annotations != nil {
			d.SetAnnotations(cr.Spec.KubeStateMetrics.Annotations)
		} else {
			for k, v := range cr.Spec.KubeStateMetrics.Annotations {
				d.Annotations[k] = v
			}
		}

		// Set labels
		d.Spec.Template.Labels["name"] = utils.TruncLabel(d.GetName())
		d.Spec.Template.Labels["app.kubernetes.io/name"] = utils.TruncLabel(d.GetName())
		d.Spec.Template.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(d.GetName(), d.GetNamespace())
		d.Spec.Template.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.KubeStateMetrics.Image)

		if cr.Spec.KubeStateMetrics.Labels != nil {
			for k, v := range cr.Spec.KubeStateMetrics.Labels {
				d.Spec.Template.Labels[k] = v
			}
		}

		if d.Spec.Template.Annotations == nil && cr.Spec.KubeStateMetrics.Annotations != nil {
			d.Spec.Template.Annotations = cr.Spec.KubeStateMetrics.Annotations
		} else {
			for k, v := range cr.Spec.KubeStateMetrics.Annotations {
				d.Spec.Template.Annotations[k] = v
			}
		}

		if len(strings.TrimSpace(cr.Spec.KubeStateMetrics.PriorityClassName)) > 0 {
			d.Spec.Template.Spec.PriorityClassName = cr.Spec.KubeStateMetrics.PriorityClassName
		}
	}
	d.Spec.Template.Spec.ServiceAccountName = cr.GetNamespace() + "-" + utils.KubestatemetricsComponentName

	return &d, nil
}

func kubeStateMetricsService(cr *v1alpha1.PlatformMonitoring) (*corev1.Service, error) {
	service := corev1.Service{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.KubestatemetricsServiceAsset), 100).Decode(&service); err != nil {
		return nil, err
	}
	//Set parameters
	service.SetGroupVersionKind(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Service"})
	service.SetName(utils.KubestatemetricsComponentName)
	service.SetNamespace(cr.GetNamespace())

	return &service, nil
}

func kubeStateMetricsServiceMonitor(cr *v1alpha1.PlatformMonitoring) (*promv1.ServiceMonitor, error) {
	sm := promv1.ServiceMonitor{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.KubestatemetricsServiceMonitorAsset), 100).Decode(&sm); err != nil {
		return nil, err
	}
	//Set parameters
	sm.SetGroupVersionKind(schema.GroupVersionKind{Group: "monitoring.coreos.com", Version: "v1", Kind: "ServiceMonitor"})
	sm.SetName(cr.GetNamespace() + "-" + utils.KubestatemetricsComponentName)
	sm.SetNamespace(cr.GetNamespace())

	if cr.Spec.KubeStateMetrics != nil && cr.Spec.KubeStateMetrics.ServiceMonitor != nil && cr.Spec.KubeStateMetrics.ServiceMonitor.IsInstall() {
		cr.Spec.KubeStateMetrics.ServiceMonitor.OverrideServiceMonitor(&sm)
	}
	sm.Spec.NamespaceSelector.MatchNames = []string{cr.GetNamespace()}

	return &sm, nil
}
