package vmalertmanager

import (
	"embed"
	"errors"
	"strings"

	v1alpha1 "github.com/Netcracker/qubership-monitoring-operator/api/v1alpha1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/victoriametrics"
	vmetricsv1b1 "github.com/VictoriaMetrics/operator/api/operator/v1beta1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/api/networking/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/yaml"
)

//go:embed  assets/*.yaml
var assets embed.FS

func vmAlertManagerServiceAccount(cr *v1alpha1.PlatformMonitoring) (*corev1.ServiceAccount, error) {
	sa := corev1.ServiceAccount{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmAlertManagerServiceAccountAsset), 100).Decode(&sa); err != nil {
		return nil, err
	}
	//Set parameters
	sa.SetGroupVersionKind(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "ServiceAccount"})
	sa.SetName(cr.GetNamespace() + "-" + utils.VmAlertManagerComponentName)
	sa.SetNamespace(cr.GetNamespace())

	return &sa, nil
}

func vmAlertManagerClusterRole(cr *v1alpha1.PlatformMonitoring, hasPsp, hasScc bool) (*rbacv1.ClusterRole, error) {
	clusterRole := rbacv1.ClusterRole{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmAlertManagerClusterRoleAsset), 100).Decode(&clusterRole); err != nil {
		return nil, err
	}
	//Set parameters
	clusterRole.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRole"})
	clusterRole.SetName(cr.GetNamespace() + "-" + utils.VmAlertManagerComponentName)
	if hasPsp {
		clusterRole.Rules = append(clusterRole.Rules, rbacv1.PolicyRule{
			Resources:     []string{"podsecuritypolicies"},
			Verbs:         []string{"use"},
			APIGroups:     []string{"policy"},
			ResourceNames: []string{utils.VmOperatorComponentName},
		})
	}
	if hasScc {
		clusterRole.Rules = append(clusterRole.Rules, rbacv1.PolicyRule{
			Resources:     []string{"securitycontextconstraints"},
			Verbs:         []string{"use"},
			APIGroups:     []string{"security.openshift.io"},
			ResourceNames: []string{utils.VmOperatorComponentName},
		})
	}

	return &clusterRole, nil
}

func vmAlertManagerClusterRoleBinding(cr *v1alpha1.PlatformMonitoring) (*rbacv1.ClusterRoleBinding, error) {
	clusterRoleBinding := rbacv1.ClusterRoleBinding{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmAlertManagerClusterRoleBindingAsset), 100).Decode(&clusterRoleBinding); err != nil {
		return nil, err
	}
	//Set parameters
	clusterRoleBinding.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRoleBinding"})
	clusterRoleBinding.SetName(cr.GetNamespace() + "-" + utils.VmAlertManagerComponentName)
	clusterRoleBinding.RoleRef.Name = cr.GetNamespace() + "-" + utils.VmAlertManagerComponentName

	// Set namespace for all subjects
	for it := range clusterRoleBinding.Subjects {
		sub := &clusterRoleBinding.Subjects[it]
		sub.Namespace = cr.GetNamespace()
		sub.Name = cr.GetNamespace() + "-" + utils.VmAlertManagerComponentName
	}
	return &clusterRoleBinding, nil
}

func vmAlertManager(r *VmAlertManagerReconciler, cr *v1alpha1.PlatformMonitoring) (*vmetricsv1b1.VMAlertmanager, error) {
	var err error
	vmalertmgr := vmetricsv1b1.VMAlertmanager{}
	if err = yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmAlertManagerAsset), 100).Decode(&vmalertmgr); err != nil {
		return nil, err
	}

	// Set parameters
	vmalertmgr.SetNamespace(cr.GetNamespace())

	if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.VmAlertManager.IsInstall() {
		vmalertmgr.Spec.Retention = cr.Spec.Victoriametrics.VmAlertManager.Retention
		vmalertmgr.Spec.Image.Repository, vmalertmgr.Spec.Image.Tag = utils.SplitImage(cr.Spec.Victoriametrics.VmAlertManager.Image)

		if r != nil {
			// Set security context
			if cr.Spec.Victoriametrics.VmAlertManager.SecurityContext != nil {
				if vmalertmgr.Spec.SecurityContext == nil {
					vmalertmgr.Spec.SecurityContext = &vmetricsv1b1.SecurityContext{}
				}
				if cr.Spec.Victoriametrics.VmAlertManager.SecurityContext.RunAsUser != nil {
					vmalertmgr.Spec.SecurityContext.RunAsUser = cr.Spec.Victoriametrics.VmAlertManager.SecurityContext.RunAsUser
				}
				if cr.Spec.Victoriametrics.VmAlertManager.SecurityContext.FSGroup != nil {
					vmalertmgr.Spec.SecurityContext.FSGroup = cr.Spec.Victoriametrics.VmAlertManager.SecurityContext.FSGroup
				}
			}
		}

		// Set resources for vmAlertManager
		vmalertmgr.Spec.ServiceAccountName = cr.GetNamespace() + "-" + utils.VmAlertManagerComponentName

		// Set resources for vmAlertManager
		if cr.Spec.Victoriametrics.VmAlertManager.Resources.Size() > 0 {
			vmalertmgr.Spec.Resources = cr.Spec.Victoriametrics.VmAlertManager.Resources
		}
		// Set secrets for vmAlertManager deployment
		if len(cr.Spec.Victoriametrics.VmAlertManager.Secrets) > 0 {
			vmalertmgr.Spec.Secrets = cr.Spec.Victoriametrics.VmAlertManager.Secrets
		}

		// Set replicas
		if cr.Spec.Victoriametrics.VmAlertManager.Replicas != nil {
			vmalertmgr.Spec.ReplicaCount = cr.Spec.Victoriametrics.VmAlertManager.Replicas
		}
		if cr.Spec.Victoriametrics.VmCluster.IsInstall() && cr.Spec.Victoriametrics.VmReplicas != nil {
			vmalertmgr.Spec.ReplicaCount = cr.Spec.Victoriametrics.VmReplicas
		}

		if cr.Spec.Victoriametrics.VmAlertManager.ConfigRawYaml != "" {
			vmalertmgr.Spec.ConfigRawYaml = cr.Spec.Victoriametrics.VmAlertManager.ConfigRawYaml
		}

		if cr.Spec.Victoriametrics.VmAlertManager.ConfigSecret != "" {
			vmalertmgr.Spec.ConfigSecret = cr.Spec.Victoriametrics.VmAlertManager.ConfigSecret
		}

		// Set additional containers
		if cr.Spec.Victoriametrics.VmAlertManager.Containers != nil {
			vmalertmgr.Spec.Containers = cr.Spec.Victoriametrics.VmAlertManager.Containers
		}

		if cr.Spec.Victoriametrics.VmAlertManager.SelectAllByDefault {
			vmalertmgr.Spec.SelectAllByDefault = true
		}

		// Set storage spec to specify how storage shall be used
		if cr.Spec.Victoriametrics.VmAlertManager.Storage != nil {
			vmalertmgr.Spec.Storage = cr.Spec.Victoriametrics.VmAlertManager.Storage
		}

		if cr.Spec.Victoriametrics.VmAlertManager.ExtraArgs != nil {
			vmalertmgr.Spec.ExtraArgs = cr.Spec.Victoriametrics.VmAlertManager.ExtraArgs
		}

		if cr.Spec.Victoriametrics.VmAlertManager.WebConfig != nil {
			vmalertmgr.Spec.WebConfig = cr.Spec.Victoriametrics.VmAlertManager.WebConfig
		}

		if cr.Spec.Victoriametrics.VmAlertManager.GossipConfig != nil {
			vmalertmgr.Spec.GossipConfig = cr.Spec.Victoriametrics.VmAlertManager.GossipConfig
		}

		if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.TLSEnabled {
			vmalertmgr.Spec.Secrets = append(vmalertmgr.Spec.Secrets, victoriametrics.GetVmalertmanagerTLSSecretName(cr.Spec.Victoriametrics.VmAlertManager))

			if cr.Spec.Victoriametrics.VmAlertManager.WebConfig == nil || cr.Spec.Victoriametrics.VmAlertManager.WebConfig.TLSServerConfig == nil {
				vmalertmgr.Spec.WebConfig = &vmetricsv1b1.AlertmanagerWebConfig{
					TLSServerConfig: &vmetricsv1b1.TLSServerConfig{
						Certs: vmetricsv1b1.Certs{
							CertSecretRef: &corev1.SecretKeySelector{
								Key: "tls.crt",
								LocalObjectReference: corev1.LocalObjectReference{
									Name: victoriametrics.GetVmalertmanagerTLSSecretName(cr.Spec.Victoriametrics.VmAlertManager),
								},
							},
							KeySecretRef: &corev1.SecretKeySelector{
								Key: "tls.key",
								LocalObjectReference: corev1.LocalObjectReference{
									Name: victoriametrics.GetVmalertmanagerTLSSecretName(cr.Spec.Victoriametrics.VmAlertManager),
								},
							},
						},
					},
				}
			}

			vmalertmgr.Spec.EmbeddedProbes = &vmetricsv1b1.EmbeddedProbes{
				LivenessProbe: &corev1.Probe{
					TimeoutSeconds:   5,
					PeriodSeconds:    5,
					SuccessThreshold: 1,
					FailureThreshold: 10,
					ProbeHandler: corev1.ProbeHandler{
						HTTPGet: &corev1.HTTPGetAction{
							Path:   "/-/healthy",
							Port:   intstr.FromString("web"),
							Scheme: "HTTPS",
						},
					},
				},
				ReadinessProbe: &corev1.Probe{
					TimeoutSeconds:   5,
					PeriodSeconds:    5,
					SuccessThreshold: 1,
					FailureThreshold: 10,
					ProbeHandler: corev1.ProbeHandler{
						HTTPGet: &corev1.HTTPGetAction{
							Path:   "/-/healthy",
							Port:   intstr.FromString("web"),
							Scheme: "HTTPS",
						},
					},
				},
			}
		}

		if cr.Spec.Victoriametrics.VmAlertManager.ExtraEnvs != nil {
			vmalertmgr.Spec.ExtraEnvs = cr.Spec.Victoriametrics.VmAlertManager.ExtraEnvs
		}

		// Set additional volumes
		if cr.Spec.Victoriametrics.VmAlertManager.Volumes != nil {
			vmalertmgr.Spec.Volumes = cr.Spec.Victoriametrics.VmAlertManager.Volumes
		}
		// Set additional volumeMounts for each vmAlertManager container. The current container names are:
		// `vmalertmanager`, `config-reloader`
		if cr.Spec.Victoriametrics.VmAlertManager.VolumeMounts != nil {
			for it := range vmalertmgr.Spec.Containers {
				c := &vmalertmgr.Spec.Containers[it]

				// Set additional volumeMounts only for vmAlertManager container
				if c.Name == utils.VmAlertManagerComponentName {
					copy(c.VolumeMounts, cr.Spec.Victoriametrics.VmAlertManager.VolumeMounts)
				}
			}
		}
		// Set nodeSelector for vmAlertManager
		if cr.Spec.Victoriametrics.VmAlertManager.NodeSelector != nil {
			vmalertmgr.Spec.NodeSelector = cr.Spec.Victoriametrics.VmAlertManager.NodeSelector
		}

		// Set affinity for vmAlertManager
		if cr.Spec.Victoriametrics.VmAlertManager.Affinity != nil {
			vmalertmgr.Spec.Affinity = cr.Spec.Victoriametrics.VmAlertManager.Affinity
		}

		// Set tolerations for vmAlertManager
		if cr.Spec.Victoriametrics.VmAlertManager.Tolerations != nil {
			vmalertmgr.Spec.Tolerations = cr.Spec.Victoriametrics.VmAlertManager.Tolerations
		}

		if cr.Spec.Victoriametrics.VmAlertManager.ConfigSelector != nil {
			vmalertmgr.Spec.ConfigSelector = cr.Spec.Victoriametrics.VmAlertManager.ConfigSelector
		}

		if cr.Spec.Victoriametrics.VmAlertManager.ConfigNamespaceSelector != nil {
			vmalertmgr.Spec.ConfigNamespaceSelector = cr.Spec.Victoriametrics.VmAlertManager.ConfigNamespaceSelector
		}

		// Set disableNamespaceMatcher
		if cr.Spec.Victoriametrics.VmAlertManager.DisableNamespaceMatcher != nil {
			vmalertmgr.Spec.DisableNamespaceMatcher = *cr.Spec.Victoriametrics.VmAlertManager.DisableNamespaceMatcher
		}

		if cr.Spec.Victoriametrics.VmAlertManager.TerminationGracePeriodSeconds != nil {
			vmalertmgr.Spec.TerminationGracePeriodSeconds = cr.Spec.Victoriametrics.VmAlertManager.TerminationGracePeriodSeconds
		}

		// Set labels
		vmalertmgr.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(vmalertmgr.GetName(), vmalertmgr.GetNamespace())
		vmalertmgr.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Victoriametrics.VmAlertManager.Image)

		vmalertmgr.Spec.PodMetadata = &vmetricsv1b1.EmbeddedObjectMetadata{Labels: map[string]string{
			"name":                         utils.TruncLabel(vmalertmgr.GetName()),
			"app.kubernetes.io/name":       utils.TruncLabel(vmalertmgr.GetName()),
			"app.kubernetes.io/instance":   utils.GetInstanceLabel(vmalertmgr.GetName(), vmalertmgr.GetNamespace()),
			"app.kubernetes.io/component":  "victoriametrics",
			"app.kubernetes.io/part-of":    "monitoring",
			"app.kubernetes.io/managed-by": "monitoring-operator",
			"app.kubernetes.io/version":    utils.GetTagFromImage(cr.Spec.Victoriametrics.VmAlertManager.Image),
		}}

		if vmalertmgr.Spec.PodMetadata != nil {
			if cr.Spec.Victoriametrics.VmAlertManager.Labels != nil {
				for k, v := range cr.Spec.Victoriametrics.VmAlertManager.Labels {
					vmalertmgr.Spec.PodMetadata.Labels[k] = v
				}
			}

			if vmalertmgr.Spec.PodMetadata.Annotations == nil && cr.Spec.Victoriametrics.VmAlertManager.Annotations != nil {
				vmalertmgr.Spec.PodMetadata.Annotations = cr.Spec.Victoriametrics.VmAlertManager.Annotations
			} else {
				for k, v := range cr.Spec.Victoriametrics.VmAlertManager.Annotations {
					vmalertmgr.Spec.PodMetadata.Annotations[k] = v
				}
			}
		}

		if len(strings.TrimSpace(cr.Spec.Victoriametrics.VmAlertManager.PriorityClassName)) > 0 {
			vmalertmgr.Spec.PriorityClassName = cr.Spec.Victoriametrics.VmAlertManager.PriorityClassName
		}
	}

	return &vmalertmgr, nil
}

func vmAlertmanagerSecret(cr *v1alpha1.PlatformMonitoring) (*corev1.Secret, error) {
	secret := corev1.Secret{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmAlertManagerSecretAsset), 100).Decode(&secret); err != nil {
		return nil, err
	}
	//Set parameters
	secret.SetGroupVersionKind(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Secret"})
	secret.SetNamespace(cr.GetNamespace())
	return &secret, nil
}

func vmAlertManagerIngressV1beta1(cr *v1alpha1.PlatformMonitoring) (*v1beta1.Ingress, error) {
	ingress := v1beta1.Ingress{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmAlertManagerIngressAsset), 100).Decode(&ingress); err != nil {
		return nil, err
	}
	//Set parameters
	ingress.SetGroupVersionKind(schema.GroupVersionKind{Group: "networking.k8s.io", Version: "v1beta1", Kind: "Ingress"})
	ingress.SetName(cr.GetNamespace() + "-" + utils.VmAlertManagerServiceName)
	ingress.SetNamespace(cr.GetNamespace())

	if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.VmAlertManager.Ingress != nil && cr.Spec.Victoriametrics.VmAlertManager.Ingress.IsInstall() {
		// Check that ingress host is specified.
		if cr.Spec.Victoriametrics.VmAlertManager.Ingress.Host == "" {
			return nil, errors.New("host for ingress can not be empty")
		}

		// Add rule for vmagent UI
		rule := v1beta1.IngressRule{Host: cr.Spec.Victoriametrics.VmAlertManager.Ingress.Host}
		serviceName := utils.VmAlertManagerServiceName
		servicePort := intstr.FromInt(utils.VmAlertManagerServicePort)

		rule.HTTP = &v1beta1.HTTPIngressRuleValue{
			Paths: []v1beta1.HTTPIngressPath{
				{
					Path: "/",
					Backend: v1beta1.IngressBackend{
						ServiceName: serviceName,
						ServicePort: servicePort,
					},
				},
			},
		}
		ingress.Spec.Rules = []v1beta1.IngressRule{rule}

		// Configure TLS if TLS secret name is set
		if cr.Spec.Victoriametrics.VmAlertManager.Ingress.TLSSecretName != "" {
			ingress.Spec.TLS = []v1beta1.IngressTLS{
				{
					Hosts:      []string{cr.Spec.Victoriametrics.VmAlertManager.Ingress.Host},
					SecretName: cr.Spec.Victoriametrics.VmAlertManager.Ingress.TLSSecretName,
				},
			}
		}

		if cr.Spec.Victoriametrics.VmAlertManager.Ingress.IngressClassName != nil {
			ingress.Spec.IngressClassName = cr.Spec.Victoriametrics.VmAlertManager.Ingress.IngressClassName
		}

		// Set annotations
		ingress.SetAnnotations(cr.Spec.Victoriametrics.VmAlertManager.Ingress.Annotations)
		if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.TLSEnabled {
			if ingress.GetAnnotations() == nil {
				annotation := make(map[string]string)
				annotation["nginx.ingress.kubernetes.io/backend-protocol"] = "HTTPS"
				ingress.SetAnnotations(annotation)
			} else {
				ingress.GetAnnotations()["nginx.ingress.kubernetes.io/backend-protocol"] = "HTTPS"
			}
		}

		// Set labels with saving default labels
		ingress.Labels["name"] = utils.TruncLabel(ingress.GetName())
		ingress.Labels["app.kubernetes.io/name"] = utils.TruncLabel(ingress.GetName())
		ingress.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(ingress.GetName(), ingress.GetNamespace())
		ingress.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Victoriametrics.VmAlertManager.Image)

		for lKey, lValue := range cr.Spec.Victoriametrics.VmAlertManager.Ingress.Labels {
			ingress.GetLabels()[lKey] = lValue
		}
	}
	return &ingress, nil
}

func vmAlertManagerIngressV1(cr *v1alpha1.PlatformMonitoring) (*networkingv1.Ingress, error) {
	ingress := networkingv1.Ingress{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmAlertManagerIngressAsset), 100).Decode(&ingress); err != nil {
		return nil, err
	}
	//Set parameters
	ingress.SetGroupVersionKind(schema.GroupVersionKind{Group: "networking.k8s.io", Version: "v1", Kind: "Ingress"})
	ingress.SetName(cr.GetNamespace() + "-" + utils.VmAlertManagerServiceName)
	ingress.SetNamespace(cr.GetNamespace())

	if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.VmAlertManager.Ingress != nil && cr.Spec.Victoriametrics.VmAlertManager.Ingress.IsInstall() {
		// Check that ingress host is specified.
		if cr.Spec.Victoriametrics.VmAlertManager.Ingress.Host == "" {
			return nil, errors.New("host for ingress can not be empty")
		}

		pathType := networkingv1.PathTypePrefix
		// Add rule for vmagent UI
		rule := networkingv1.IngressRule{Host: cr.Spec.Victoriametrics.VmAlertManager.Ingress.Host}
		rule.HTTP = &networkingv1.HTTPIngressRuleValue{
			Paths: []networkingv1.HTTPIngressPath{
				{
					Path:     "/",
					PathType: &pathType,
					Backend: networkingv1.IngressBackend{
						Service: &networkingv1.IngressServiceBackend{
							Name: utils.VmAlertManagerServiceName,
							Port: networkingv1.ServiceBackendPort{
								Number: int32(utils.VmAlertManagerServicePort),
							},
						},
					},
				},
			},
		}

		ingress.Spec.Rules = []networkingv1.IngressRule{rule}

		// Configure TLS if TLS secret name is set
		if cr.Spec.Victoriametrics.VmAlertManager.Ingress.TLSSecretName != "" {
			ingress.Spec.TLS = []networkingv1.IngressTLS{
				{
					Hosts:      []string{cr.Spec.Victoriametrics.VmAlertManager.Ingress.Host},
					SecretName: cr.Spec.Victoriametrics.VmAlertManager.Ingress.TLSSecretName,
				},
			}
		}

		if cr.Spec.Victoriametrics.VmAlertManager.Ingress.IngressClassName != nil {
			ingress.Spec.IngressClassName = cr.Spec.Victoriametrics.VmAlertManager.Ingress.IngressClassName
		}

		// Set annotations
		ingress.SetAnnotations(cr.Spec.Victoriametrics.VmAlertManager.Ingress.Annotations)
		if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.TLSEnabled {
			if ingress.GetAnnotations() == nil {
				annotation := make(map[string]string)
				annotation["nginx.ingress.kubernetes.io/backend-protocol"] = "HTTPS"
				ingress.SetAnnotations(annotation)
			} else {
				ingress.GetAnnotations()["nginx.ingress.kubernetes.io/backend-protocol"] = "HTTPS"
			}
		}

		// Set labels with saving default labels
		ingress.Labels["name"] = utils.TruncLabel(ingress.GetName())
		ingress.Labels["app.kubernetes.io/name"] = utils.TruncLabel(ingress.GetName())
		ingress.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(ingress.GetName(), ingress.GetNamespace())
		ingress.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Victoriametrics.VmAlertManager.Image)

		for lKey, lValue := range cr.Spec.Victoriametrics.VmAlertManager.Ingress.Labels {
			ingress.GetLabels()[lKey] = lValue
		}
	}
	return &ingress, nil
}
