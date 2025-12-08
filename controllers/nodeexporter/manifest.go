package nodeexporter

import (
	"embed"
	"fmt"
	"strings"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/yaml"
)

//go:embed  assets/*.yaml
var assets embed.FS

func nodeExporterClusterRole(cr *monv1.PlatformMonitoring, hasPsp, hasScc bool) (*rbacv1.ClusterRole, error) {
	clusterRole := rbacv1.ClusterRole{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.NodeExporterClusterRoleAsset), 100).Decode(&clusterRole); err != nil {
		return nil, err
	}
	//Set parameters
	clusterRole.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRole"})
	clusterRole.SetName(cr.GetNamespace() + "-" + utils.NodeExporterComponentName)
	if hasPsp {
		clusterRole.Rules = append(clusterRole.Rules, rbacv1.PolicyRule{
			Resources:     []string{"podsecuritypolicies"},
			Verbs:         []string{"use"},
			APIGroups:     []string{"policy"},
			ResourceNames: []string{utils.NodeExporterComponentName},
		})
	}
	if hasScc {
		clusterRole.Rules = append(clusterRole.Rules, rbacv1.PolicyRule{
			Resources:     []string{"securitycontextconstraints"},
			Verbs:         []string{"use"},
			APIGroups:     []string{"security.openshift.io"},
			ResourceNames: []string{utils.NodeExporterComponentName},
		})
	}

	return &clusterRole, nil
}

func nodeExporterClusterRoleBinding(cr *monv1.PlatformMonitoring) (*rbacv1.ClusterRoleBinding, error) {
	clusterRoleBinding := rbacv1.ClusterRoleBinding{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.NodeExporterClusterRoleBindingAsset), 100).Decode(&clusterRoleBinding); err != nil {
		return nil, err
	}
	//Set parameters
	clusterRoleBinding.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRoleBinding"})
	clusterRoleBinding.SetName(cr.GetNamespace() + "-" + utils.NodeExporterComponentName)
	clusterRoleBinding.RoleRef.Name = cr.GetNamespace() + "-" + utils.NodeExporterComponentName

	// Set namespace for all subjects
	for it := range clusterRoleBinding.Subjects {
		sub := &clusterRoleBinding.Subjects[it]
		sub.Namespace = cr.GetNamespace()
		sub.Name = cr.GetNamespace() + "-" + utils.NodeExporterComponentName
	}
	return &clusterRoleBinding, nil
}

func nodeExporterDaemonSet(cr *monv1.PlatformMonitoring) (*appsv1.DaemonSet, error) {
	daemonSet := appsv1.DaemonSet{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.NodeExporterDaemonSetAsset), 100).Decode(&daemonSet); err != nil {
		return nil, err
	}
	//Set parameters
	daemonSet.SetGroupVersionKind(schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "DaemonSet"})
	daemonSet.SetName(utils.NodeExporterComponentName)
	daemonSet.SetNamespace(cr.GetNamespace())

	if cr.Spec.NodeExporter != nil {
		// Find container with b.Name as name and set Image from custom resource
		for it := range daemonSet.Spec.Template.Spec.Containers {
			c := &daemonSet.Spec.Template.Spec.Containers[it]
			if c.Name == utils.NodeExporterComponentName {
				c.Image = cr.Spec.NodeExporter.Image
				portValue := cr.Spec.NodeExporter.Port
				c.Args[len(c.Args)-1] = fmt.Sprintf("--web.listen-address=:%d", portValue)
				if len(cr.Spec.NodeExporter.ExtraArgs) > 0 {
					c.Args = append(c.Args, cr.Spec.NodeExporter.ExtraArgs...)
				}
				for p := range c.Ports {
					port := &c.Ports[p]
					if port.Name == utils.NodeExporterMetricsPortName {
						port.HostPort = portValue
						port.ContainerPort = portValue
					}
				}
				if cr.Spec.NodeExporter.Resources.Size() > 0 {
					c.Resources = cr.Spec.NodeExporter.Resources
				}
				break
			}
		}
		// Set NodeSelector for nodeExporter
		if cr.Spec.NodeExporter.NodeSelector != nil {
			daemonSet.Spec.Template.Spec.NodeSelector = cr.Spec.NodeExporter.NodeSelector
		}
		// Set security context
		if cr.Spec.NodeExporter.SecurityContext != nil {
			if daemonSet.Spec.Template.Spec.SecurityContext == nil {
				daemonSet.Spec.Template.Spec.SecurityContext = &corev1.PodSecurityContext{}
			}
			if cr.Spec.NodeExporter.SecurityContext.RunAsUser != nil {
				daemonSet.Spec.Template.Spec.SecurityContext.RunAsUser = cr.Spec.NodeExporter.SecurityContext.RunAsUser
			}
			if cr.Spec.NodeExporter.SecurityContext.FSGroup != nil {
				daemonSet.Spec.Template.Spec.SecurityContext.FSGroup = cr.Spec.NodeExporter.SecurityContext.FSGroup
			}
		}
		// Set tolerations for NodeExporter
		if cr.Spec.NodeExporter.Tolerations != nil {
			daemonSet.Spec.Template.Spec.Tolerations = cr.Spec.NodeExporter.Tolerations
		}

		// Set affinity for NodeExporter
		if cr.Spec.NodeExporter.Affinity != nil {
			daemonSet.Spec.Template.Spec.Affinity = cr.Spec.NodeExporter.Affinity
		}

		// Set labels
		daemonSet.Labels["name"] = utils.TruncLabel(daemonSet.GetName())
		daemonSet.Labels["app.kubernetes.io/name"] = utils.TruncLabel(daemonSet.GetName())
		daemonSet.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(daemonSet.GetName(), daemonSet.GetNamespace())
		daemonSet.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.NodeExporter.Image)

		if cr.Spec.NodeExporter.Labels != nil {
			for k, v := range cr.Spec.NodeExporter.Labels {
				daemonSet.Labels[k] = v
			}
		}

		if daemonSet.Annotations == nil && cr.Spec.NodeExporter.Annotations != nil {
			daemonSet.SetAnnotations(cr.Spec.NodeExporter.Annotations)
		} else {
			for k, v := range cr.Spec.NodeExporter.Annotations {
				daemonSet.Annotations[k] = v
			}
		}

		// Set labels
		daemonSet.Spec.Template.Labels["name"] = utils.TruncLabel(daemonSet.GetName())
		daemonSet.Spec.Template.Labels["app.kubernetes.io/name"] = utils.TruncLabel(daemonSet.GetName())
		daemonSet.Spec.Template.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(daemonSet.GetName(), daemonSet.GetNamespace())
		daemonSet.Spec.Template.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.NodeExporter.Image)

		if cr.Spec.NodeExporter.Labels != nil {
			for k, v := range cr.Spec.NodeExporter.Labels {
				daemonSet.Spec.Template.Labels[k] = v
			}
		}

		if daemonSet.Spec.Template.Annotations == nil && cr.Spec.NodeExporter.Annotations != nil {
			daemonSet.Spec.Template.Annotations = cr.Spec.NodeExporter.Annotations
		} else {
			for k, v := range cr.Spec.NodeExporter.Annotations {
				daemonSet.Spec.Template.Annotations[k] = v
			}
		}
		if cr.Spec.NodeExporter.CollectorTextfileDirectory != "" {
			for it := range daemonSet.Spec.Template.Spec.Volumes {
				v := &daemonSet.Spec.Template.Spec.Volumes[it]
				if v.Name == utils.NodeExporterTextfileVolumeName {
					v.HostPath.Path = cr.Spec.NodeExporter.CollectorTextfileDirectory
				}
			}
		}

		if len(strings.TrimSpace(cr.Spec.NodeExporter.PriorityClassName)) > 0 {
			daemonSet.Spec.Template.Spec.PriorityClassName = cr.Spec.NodeExporter.PriorityClassName
		}
	}
	daemonSet.Spec.Template.Spec.ServiceAccountName = cr.GetNamespace() + "-" + utils.NodeExporterComponentName

	return &daemonSet, nil
}

func nodeExporterService(cr *monv1.PlatformMonitoring) (*corev1.Service, error) {
	service := corev1.Service{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.NodeExporterServiceAsset), 100).Decode(&service); err != nil {
		return nil, err
	}
	//Set parameters
	service.SetGroupVersionKind(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Service"})
	service.SetName(utils.NodeExporterComponentName)
	service.SetNamespace(cr.GetNamespace())

	if cr.Spec.NodeExporter != nil {
		// Set port
		for p := range service.Spec.Ports {
			port := &service.Spec.Ports[p]
			if port.Name == utils.NodeExporterMetricsPortName {
				port.Port = cr.Spec.NodeExporter.Port
				port.TargetPort = intstr.FromInt(int(cr.Spec.NodeExporter.Port))
			}
		}
	}
	return &service, nil
}

func nodeExporterServiceMonitor(cr *monv1.PlatformMonitoring) (*promv1.ServiceMonitor, error) {
	sm := promv1.ServiceMonitor{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.NodeExporterServiceMonitorAsset), 100).Decode(&sm); err != nil {
		return nil, err
	}
	//Set parameters
	sm.SetGroupVersionKind(schema.GroupVersionKind{Group: "monitoring.coreos.com", Version: "v1", Kind: "ServiceMonitor"})
	sm.SetName(cr.GetNamespace() + "-" + utils.NodeExporterComponentName)
	sm.SetNamespace(cr.GetNamespace())

	if cr.Spec.NodeExporter != nil && cr.Spec.NodeExporter.ServiceMonitor != nil && cr.Spec.NodeExporter.ServiceMonitor.IsInstall() {
		cr.Spec.NodeExporter.ServiceMonitor.OverrideServiceMonitor(&sm)
	}
	sm.Spec.NamespaceSelector.MatchNames = []string{cr.GetNamespace()}

	return &sm, nil
}

func nodeExporterServiceAccount(cr *monv1.PlatformMonitoring) (*corev1.ServiceAccount, error) {
	sa := corev1.ServiceAccount{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.NodeExporterServiceAccountAsset), 100).Decode(&sa); err != nil {
		return nil, err
	}
	//Set parameters
	sa.SetGroupVersionKind(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "ServiceAccount"})
	sa.SetName(cr.GetNamespace() + "-" + utils.NodeExporterComponentName)
	sa.SetNamespace(cr.GetNamespace())

	// Set annotations and labels for ServiceAccount in case
	if cr.Spec.NodeExporter != nil && cr.Spec.NodeExporter.ServiceAccount != nil {
		if sa.Annotations == nil && cr.Spec.NodeExporter.ServiceAccount.Annotations != nil {
			sa.SetAnnotations(cr.Spec.NodeExporter.ServiceAccount.Annotations)
		} else {
			for k, v := range cr.Spec.NodeExporter.ServiceAccount.Annotations {
				sa.Annotations[k] = v
			}
		}

		if sa.Labels == nil && cr.Spec.NodeExporter.ServiceAccount.Labels != nil {
			sa.SetLabels(cr.Spec.NodeExporter.ServiceAccount.Labels)
		} else {
			for k, v := range cr.Spec.NodeExporter.ServiceAccount.Labels {
				sa.Labels[k] = v
			}
		}
	}

	return &sa, nil
}
