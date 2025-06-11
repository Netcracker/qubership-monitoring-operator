package vmoperator

import (
	"embed"
	"strings"

	v1alpha1 "github.com/Netcracker/qubership-monitoring-operator/api/v1alpha1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/victoriametrics"
	secv1 "github.com/openshift/api/security/v1"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/utils/ptr"
)

//go:embed  assets/*.yaml
var assets embed.FS

// VmOperatorRole builds the ServiceAccount resource manifest
// and fill it with parameters from the CR.
func vmOperatorRole(cr *v1alpha1.PlatformMonitoring) (*rbacv1.Role, error) {
	role := rbacv1.Role{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmOperatorRoleAsset), 100).Decode(&role); err != nil {
		return nil, err
	}
	//Set parameters
	role.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "Role"})
	role.SetName(cr.GetNamespace() + "-" + utils.VmOperatorComponentName)
	role.SetNamespace(cr.GetNamespace())

	return &role, nil
}

func vmOperatorServiceAccount(cr *v1alpha1.PlatformMonitoring) (*corev1.ServiceAccount, error) {
	sa := corev1.ServiceAccount{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmOperatorServiceAccountAsset), 100).Decode(&sa); err != nil {
		return nil, err
	}
	//Set parameters
	sa.SetGroupVersionKind(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "ServiceAccount"})
	sa.SetName(cr.GetNamespace() + "-" + utils.VmOperatorComponentName)
	sa.SetNamespace(cr.GetNamespace())

	return &sa, nil
}

func vmOperatorRoleBinding(cr *v1alpha1.PlatformMonitoring) (*rbacv1.RoleBinding, error) {
	roleBinding := rbacv1.RoleBinding{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmOperatorRoleBindingAsset), 100).Decode(&roleBinding); err != nil {
		return nil, err
	}
	//Set parameters
	roleBinding.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "RoleBinding"})
	roleBinding.SetName(cr.GetNamespace() + "-" + utils.VmOperatorComponentName)
	roleBinding.SetNamespace(cr.GetNamespace())
	roleBinding.RoleRef.Name = cr.GetNamespace() + "-" + utils.VmOperatorComponentName

	// Set namespace for all subjects
	for it := range roleBinding.Subjects {
		sub := &roleBinding.Subjects[it]
		sub.Namespace = cr.GetNamespace()
		sub.Name = cr.GetNamespace() + "-" + utils.VmOperatorComponentName
	}
	return &roleBinding, nil
}

func vmOperatorClusterRole(cr *v1alpha1.PlatformMonitoring) (*rbacv1.ClusterRole, error) {
	clusterRole := rbacv1.ClusterRole{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmOperatorClusterRoleAsset), 100).Decode(&clusterRole); err != nil {
		return nil, err
	}
	//Set parameters
	clusterRole.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRole"})
	clusterRole.SetName(cr.GetNamespace() + "-" + utils.VmOperatorComponentName)

	// Add permissions to kubelet service
	clusterRole.Rules = append(clusterRole.Rules, rbacv1.PolicyRule{
		APIGroups: []string{""},
		Resources: []string{"services", "services/finalizers", "endpoints"},
		Verbs:     []string{"get", "create", "list", "update", "watch"},
	})
	return &clusterRole, nil
}

func vmOperatorClusterRoleBinding(cr *v1alpha1.PlatformMonitoring) (*rbacv1.ClusterRoleBinding, error) {
	clusterRoleBinding := rbacv1.ClusterRoleBinding{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmOperatorClusterRoleBindingAsset), 100).Decode(&clusterRoleBinding); err != nil {
		return nil, err
	}
	//Set parameters
	clusterRoleBinding.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRoleBinding"})
	clusterRoleBinding.SetName(cr.GetNamespace() + "-" + utils.VmOperatorComponentName)
	clusterRoleBinding.RoleRef.Name = cr.GetNamespace() + "-" + utils.VmOperatorComponentName

	// Set namespace for all subjects
	for it := range clusterRoleBinding.Subjects {
		sub := &clusterRoleBinding.Subjects[it]
		sub.Namespace = cr.GetNamespace()
		sub.Name = cr.GetNamespace() + "-" + utils.VmOperatorComponentName
	}
	return &clusterRoleBinding, nil
}

func vmOperatorDeployment(r *VmOperatorReconciler, cr *v1alpha1.PlatformMonitoring) (*appsv1.Deployment, error) {
	d := appsv1.Deployment{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmOperatorDeploymentAsset), 100).Decode(&d); err != nil {
		return nil, err
	}
	//Set parameters
	d.SetGroupVersionKind(schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"})
	d.SetName(utils.VmOperatorComponentName)
	d.SetNamespace(cr.GetNamespace())

	if cr.Spec.Victoriametrics != nil {
		if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.TLSEnabled {
			var volume corev1.Volume
			volume.Name = victoriametrics.GetVmoperatorTLSSecretName(cr.Spec.Victoriametrics.VmOperator)
			volume.Secret = &corev1.SecretVolumeSource{
				SecretName: victoriametrics.GetVmoperatorTLSSecretName(cr.Spec.Victoriametrics.VmOperator),
			}
			d.Spec.Template.Spec.Volumes = append(d.Spec.Template.Spec.Volumes, volume)
		}

		// Set correct images and any parameters to containers spec
		for it := range d.Spec.Template.Spec.Containers {
			c := &d.Spec.Template.Spec.Containers[it]
			if c.Name == utils.VmOperatorComponentName {
				// Set vm-operator image
				c.Image = cr.Spec.Victoriametrics.VmOperator.Image

				if cr.Spec.Victoriametrics.VmOperator.Resources.Size() > 0 {
					c.Resources = cr.Spec.Victoriametrics.VmOperator.Resources
				}

				if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.TLSEnabled {
					c.VolumeMounts = []corev1.VolumeMount{
						{
							Name:      victoriametrics.GetVmoperatorTLSSecretName(cr.Spec.Victoriametrics.VmOperator),
							MountPath: "/tmp/k8s-webhook-server/serving-certs/ca.crt",
							SubPath:   "ca.crt",
							ReadOnly:  true,
						},
						{
							Name:      victoriametrics.GetVmoperatorTLSSecretName(cr.Spec.Victoriametrics.VmOperator),
							MountPath: "/tmp/k8s-webhook-server/serving-certs/tls.crt",
							SubPath:   "tls.crt",
							ReadOnly:  true,
						},
						{
							Name:      victoriametrics.GetVmoperatorTLSSecretName(cr.Spec.Victoriametrics.VmOperator),
							MountPath: "/tmp/k8s-webhook-server/serving-certs/tls.key",
							SubPath:   "tls.key",
							ReadOnly:  true,
						},
					}
					c.Args = append(c.Args, "--tls.enable=true")
				}

				if r != nil && r.HasRouteApi() {
					c.Args = append(c.Args, "--controller.disableCRDOwnership=true")
				} else {
					c.Args = append(c.Args, "--leader-elect")
				}
				break
			}
		}
		// Set security context
		if cr.Spec.Victoriametrics.VmOperator.SecurityContext != nil {
			if d.Spec.Template.Spec.SecurityContext == nil {
				d.Spec.Template.Spec.SecurityContext = &corev1.PodSecurityContext{}
			}
			if cr.Spec.Victoriametrics.VmOperator.SecurityContext.RunAsUser != nil {
				d.Spec.Template.Spec.SecurityContext.RunAsUser = cr.Spec.Victoriametrics.VmOperator.SecurityContext.RunAsUser
			}
			if cr.Spec.Victoriametrics.VmOperator.SecurityContext.FSGroup != nil {
				d.Spec.Template.Spec.SecurityContext.FSGroup = cr.Spec.Victoriametrics.VmOperator.SecurityContext.FSGroup
			}
		}
		if cr.Spec.Victoriametrics.VmOperator.Replicas != nil {
			d.Spec.Replicas = cr.Spec.Victoriametrics.VmOperator.Replicas
		}
		if cr.Spec.Victoriametrics.VmCluster.IsInstall() && cr.Spec.Victoriametrics.VmReplicas != nil {
			d.Spec.Replicas = cr.Spec.Victoriametrics.VmReplicas
		}
		// Set tolerations for VmOperator
		if cr.Spec.Victoriametrics.VmOperator.Tolerations != nil {
			d.Spec.Template.Spec.Tolerations = cr.Spec.Victoriametrics.VmOperator.Tolerations
		}
		// Set nodeSelector for VmOperator
		if cr.Spec.Victoriametrics.VmOperator.NodeSelector != nil {
			d.Spec.Template.Spec.NodeSelector = cr.Spec.Victoriametrics.VmOperator.NodeSelector
		}
		// Set affinity for VmOperator
		if cr.Spec.Victoriametrics.VmOperator.Affinity != nil {
			d.Spec.Template.Spec.Affinity = cr.Spec.Victoriametrics.VmOperator.Affinity
		}

		// Add extraEnvs for VmOperator
		if cr.Spec.Victoriametrics.VmOperator.ExtraEnvs != nil {
			for it := range d.Spec.Template.Spec.Containers {
				c := &d.Spec.Template.Spec.Containers[it]
				if c.Name == utils.VmOperatorComponentName {
					c.Env = append(c.Env, cr.Spec.Victoriametrics.VmOperator.ExtraEnvs...)
				}
			}
		}

		// Set security context for VmOperator container.
		if cr.Spec.Victoriametrics.VmOperator.ContainerSecurityContext != nil {
			for it := range d.Spec.Template.Spec.Containers {
				c := &d.Spec.Template.Spec.Containers[it]
				if c.Name == utils.VmOperatorComponentName {
					c.SecurityContext = cr.Spec.Victoriametrics.VmOperator.ContainerSecurityContext
					c.SecurityContext.RunAsUser = cr.Spec.Victoriametrics.VmOperator.ContainerSecurityContext.RunAsUser
					c.SecurityContext.RunAsGroup = cr.Spec.Victoriametrics.VmOperator.ContainerSecurityContext.RunAsGroup
				}
			}
		}

		// Set labels
		d.Labels["name"] = utils.TruncLabel(d.GetName())
		d.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Victoriametrics.VmOperator.Image)

		if cr.Spec.Victoriametrics.VmOperator.Labels != nil {
			for k, v := range cr.Spec.Victoriametrics.VmOperator.Labels {
				d.Labels[k] = v
			}
		}

		if d.Annotations == nil && cr.Spec.Victoriametrics.VmOperator.Annotations != nil {
			d.SetAnnotations(cr.Spec.Victoriametrics.VmOperator.Annotations)
		} else {
			for k, v := range cr.Spec.Victoriametrics.VmOperator.Annotations {
				d.Annotations[k] = v
			}
		}

		// Set labels
		d.Spec.Template.Labels["name"] = utils.TruncLabel(d.GetName())
		d.Spec.Template.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Victoriametrics.VmOperator.Image)

		if cr.Spec.Victoriametrics.VmOperator.Labels != nil {
			for k, v := range cr.Spec.Victoriametrics.VmOperator.Labels {
				d.Spec.Template.Labels[k] = v
			}
		}

		if d.Spec.Template.Annotations == nil && cr.Spec.Victoriametrics.VmOperator.Annotations != nil {
			d.Spec.Template.Annotations = cr.Spec.Victoriametrics.VmOperator.Annotations
		} else {
			for k, v := range cr.Spec.Victoriametrics.VmOperator.Annotations {
				d.Spec.Template.Annotations[k] = v
			}
		}
		if len(strings.TrimSpace(cr.Spec.Victoriametrics.VmOperator.PriorityClassName)) > 0 {
			d.Spec.Template.Spec.PriorityClassName = cr.Spec.Victoriametrics.VmOperator.PriorityClassName
		}
	}
	d.Spec.Template.Spec.ServiceAccountName = cr.GetNamespace() + "-" + utils.VmOperatorComponentName

	return &d, nil
}

func vmOperatorService(cr *v1alpha1.PlatformMonitoring) (*corev1.Service, error) {
	service := corev1.Service{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmOperatorServiceAsset), 100).Decode(&service); err != nil {
		return nil, err
	}
	//Set parameters
	service.SetGroupVersionKind(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Service"})
	service.SetName(utils.VmOperatorComponentName)
	service.SetNamespace(cr.GetNamespace())

	return &service, nil
}

func vmKubeletService(cr *v1alpha1.PlatformMonitoring) (*corev1.Service, error) {
	service := corev1.Service{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmKubeletServiceAsset), 100).Decode(&service); err != nil {
		return nil, err
	}
	//Set parameters
	service.SetGroupVersionKind(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Service"})
	service.SetName(utils.VmKubeletName)
	service.SetNamespace(cr.GetNamespace())

	return &service, nil
}

func vmKubeletServiceEndpoints(cr *v1alpha1.PlatformMonitoring) (*corev1.Endpoints, error) {
	endpoints := corev1.Endpoints{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmKubeletServiceEndpointsAsset), 100).Decode(&endpoints); err != nil {
		return nil, err
	}
	//Set parameters
	endpoints.SetGroupVersionKind(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Endpoints"})
	endpoints.SetName(utils.VmKubeletName)
	endpoints.SetNamespace(cr.GetNamespace())

	return &endpoints, nil
}

func vmOperatorServiceMonitor(cr *v1alpha1.PlatformMonitoring) (*promv1.ServiceMonitor, error) {
	sm := promv1.ServiceMonitor{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmOperatorServiceMonitorAsset), 100).Decode(&sm); err != nil {
		return nil, err
	}
	//Set parameters
	sm.SetGroupVersionKind(schema.GroupVersionKind{Group: "monitoring.coreos.com", Version: "v1", Kind: "ServiceMonitor"})
	sm.SetName(cr.GetNamespace() + "-" + utils.VmOperatorComponentName)
	sm.SetNamespace(cr.GetNamespace())

	if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.VmOperator.ServiceMonitor != nil && cr.Spec.Victoriametrics.VmOperator.ServiceMonitor.IsInstall() {
		cr.Spec.Victoriametrics.VmOperator.ServiceMonitor.OverrideServiceMonitor(&sm)
	}
	sm.Spec.NamespaceSelector.MatchNames = []string{cr.GetNamespace()}
	if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.TLSEnabled {
		for i := range sm.Spec.Endpoints {
			sm.Spec.Endpoints[i].Scheme = "https"
			sm.Spec.Endpoints[i].TLSConfig = &promv1.TLSConfig{
				SafeTLSConfig: promv1.SafeTLSConfig{
					InsecureSkipVerify: ptr.To(true),
				},
			}
		}
	}

	return &sm, nil
}

func vmOperatorSecurityContextConstraints() (*secv1.SecurityContextConstraints, error) {
	scc := secv1.SecurityContextConstraints{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmOperatorSecurityContextConstraintsAsset), 100).Decode(&scc); err != nil {
		return nil, err
	}
	//Set parameters
	scc.SetGroupVersionKind(schema.GroupVersionKind{Group: "security.openshift.io", Version: "v1", Kind: "SecurityContextConstraints"})
	scc.SetName(utils.VmOperatorComponentName)

	return &scc, nil
}
