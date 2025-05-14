package vmsingle

import (
	"embed"
	"errors"
	"strings"

	"maps"

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

func vmSingleServiceAccount(cr *v1alpha1.PlatformMonitoring) (*corev1.ServiceAccount, error) {
	sa := corev1.ServiceAccount{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmSingleServiceAccountAsset), 100).Decode(&sa); err != nil {
		return nil, err
	}
	//Set parameters
	sa.SetGroupVersionKind(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "ServiceAccount"})
	sa.SetName(cr.GetNamespace() + "-" + utils.VmSingleComponentName)
	sa.SetNamespace(cr.GetNamespace())

	return &sa, nil
}

func vmSingleClusterRole(cr *v1alpha1.PlatformMonitoring, hasPsp, hasScc bool) (*rbacv1.ClusterRole, error) {
	clusterRole := rbacv1.ClusterRole{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmSingleClusterRoleAsset), 100).Decode(&clusterRole); err != nil {
		return nil, err
	}
	//Set parameters
	clusterRole.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRole"})
	clusterRole.SetName(cr.GetNamespace() + "-" + utils.VmSingleComponentName)
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

func vmSingleClusterRoleBinding(cr *v1alpha1.PlatformMonitoring) (*rbacv1.ClusterRoleBinding, error) {
	clusterRoleBinding := rbacv1.ClusterRoleBinding{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmSingleClusterRoleBindingAsset), 100).Decode(&clusterRoleBinding); err != nil {
		return nil, err
	}
	//Set parameters
	clusterRoleBinding.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRoleBinding"})
	clusterRoleBinding.SetName(cr.GetNamespace() + "-" + utils.VmSingleComponentName)
	clusterRoleBinding.RoleRef.Name = cr.GetNamespace() + "-" + utils.VmSingleComponentName

	// Set namespace for all subjects
	for it := range clusterRoleBinding.Subjects {
		sub := &clusterRoleBinding.Subjects[it]
		sub.Namespace = cr.GetNamespace()
		sub.Name = cr.GetNamespace() + "-" + utils.VmSingleComponentName
	}
	return &clusterRoleBinding, nil
}

func vmSingle(r *VmSingleReconciler, cr *v1alpha1.PlatformMonitoring) (*vmetricsv1b1.VMSingle, error) {
	var err error
	vmsingle := vmetricsv1b1.VMSingle{}
	if err = yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmSingleAsset), 100).Decode(&vmsingle); err != nil {
		return nil, err
	}

	// Set parameters
	vmsingle.SetNamespace(cr.GetNamespace())

	if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.VmSingle.IsInstall() {
		vmsingle.Spec.RetentionPeriod = cr.Spec.Victoriametrics.VmSingle.RetentionPeriod

		// Set vmsingle image
		vmsingle.Spec.Image.Repository, vmsingle.Spec.Image.Tag = utils.SplitImage(cr.Spec.Victoriametrics.VmSingle.Image)

		if r != nil {
			// Set security context
			if cr.Spec.Victoriametrics.VmSingle.SecurityContext != nil {
				if vmsingle.Spec.SecurityContext == nil {
					vmsingle.Spec.SecurityContext = &vmetricsv1b1.SecurityContext{}
				}
				if cr.Spec.Victoriametrics.VmSingle.SecurityContext.RunAsUser != nil {
					vmsingle.Spec.SecurityContext.RunAsUser = cr.Spec.Victoriametrics.VmSingle.SecurityContext.RunAsUser
				}
				if cr.Spec.Victoriametrics.VmSingle.SecurityContext.RunAsGroup != nil {
					vmsingle.Spec.SecurityContext.RunAsGroup = cr.Spec.Victoriametrics.VmSingle.SecurityContext.RunAsGroup
				}
				if cr.Spec.Victoriametrics.VmSingle.SecurityContext.FSGroup != nil {
					vmsingle.Spec.SecurityContext.FSGroup = cr.Spec.Victoriametrics.VmSingle.SecurityContext.FSGroup
				}
			}
		}

		vmsingle.Spec.ServiceAccountName = cr.GetNamespace() + "-" + utils.VmSingleComponentName

		// Set resources for vmsingle deployment
		if cr.Spec.Victoriametrics.VmSingle.Resources.Size() > 0 {
			vmsingle.Spec.Resources = cr.Spec.Victoriametrics.VmSingle.Resources
		}
		// Set additional containers
		if cr.Spec.Victoriametrics.VmSingle.Containers != nil {
			vmsingle.Spec.Containers = cr.Spec.Victoriametrics.VmSingle.Containers
		}
		// Set secrets for vmsingle deployment
		if len(cr.Spec.Victoriametrics.VmSingle.Secrets) > 0 {
			vmsingle.Spec.Secrets = cr.Spec.Victoriametrics.VmSingle.Secrets
		}

		// Set additional volumes
		if cr.Spec.Victoriametrics.VmSingle.Volumes != nil {
			vmsingle.Spec.Volumes = cr.Spec.Victoriametrics.VmSingle.Volumes
		}

		if cr.Spec.Victoriametrics.VmSingle.VolumeMounts != nil {
			for it := range vmsingle.Spec.Containers {
				c := &vmsingle.Spec.Containers[it]

				// Set additional volumeMounts only for vmsingle container
				if c.Name == utils.VmSingleComponentName {
					c.VolumeMounts = cr.Spec.Victoriametrics.VmSingle.VolumeMounts
				}
			}
		}
		if len(strings.TrimSpace(cr.Spec.Victoriametrics.VmSingle.StorageDataPath)) > 0 {
			vmsingle.Spec.StorageDataPath = cr.Spec.Victoriametrics.VmSingle.StorageDataPath
		}
		if cr.Spec.Victoriametrics.VmSingle.Storage != nil {
			vmsingle.Spec.Storage = cr.Spec.Victoriametrics.VmSingle.Storage
		}
		if cr.Spec.Victoriametrics.VmSingle.StorageMetadata != nil {
			vmsingle.Spec.StorageMetadata = *cr.Spec.Victoriametrics.VmSingle.StorageMetadata
		}
		// Set nodeSelector for vmsingle deployment
		if cr.Spec.Victoriametrics.VmSingle.NodeSelector != nil {
			vmsingle.Spec.NodeSelector = cr.Spec.Victoriametrics.VmSingle.NodeSelector
		}
		// Set affinity for vmsingle deployment
		if cr.Spec.Victoriametrics.VmSingle.Affinity != nil {
			vmsingle.Spec.Affinity = cr.Spec.Victoriametrics.VmSingle.Affinity
		}

		if cr.Spec.Victoriametrics.VmSingle.ExtraArgs != nil {
			maps.Copy(vmsingle.Spec.ExtraArgs, cr.Spec.Victoriametrics.VmSingle.ExtraArgs)
		}

		//A single-node VictoriaMetrics is capable of proxying requests to vmalert
		//https://docs.victoriametrics.com/Single-server-VictoriaMetrics.html#vmalert
		if cr.Spec.Victoriametrics.VmOperator.IsInstall() && cr.Spec.Victoriametrics.VmAlert.IsInstall() {
			vmAlert := vmetricsv1b1.VMAlert{}
			vmAlert.SetName(utils.VmComponentName)
			vmAlert.SetNamespace(cr.GetNamespace())
			if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.TLSEnabled {
				vmAlert.Spec.ExtraArgs = make(map[string]string)
				maps.Copy(vmAlert.Spec.ExtraArgs, map[string]string{"tls": "true"})
			}
			if cr.Spec.Victoriametrics.VmAlert.Port != "" {
				vmAlert.Spec.Port = cr.Spec.Victoriametrics.VmAlert.Port
			} else {
				vmAlert.Spec.Port = "8080"
			}
			maps.Copy(vmsingle.Spec.ExtraArgs, map[string]string{"vmalert.proxyURL": vmAlert.AsURL()})
		}

		if cr.Spec.Victoriametrics.VmAgent.Replicas != nil && *cr.Spec.Victoriametrics.VmAgent.Replicas > 1 {
			maps.Copy(vmsingle.Spec.ExtraArgs, map[string]string{"dedup.minScrapeInterval": "30s"})
		}

		if cr.Spec.Victoriametrics.VmSingle.ExtraEnvs != nil {
			vmsingle.Spec.ExtraEnvs = cr.Spec.Victoriametrics.VmSingle.ExtraEnvs
		}

		// Set tolerations for vmsingle
		if cr.Spec.Victoriametrics.VmSingle.Tolerations != nil {
			vmsingle.Spec.Tolerations = cr.Spec.Victoriametrics.VmSingle.Tolerations
		}

		if cr.Spec.Victoriametrics.VmSingle.TerminationGracePeriodSeconds != nil {
			vmsingle.Spec.TerminationGracePeriodSeconds = cr.Spec.Victoriametrics.VmSingle.TerminationGracePeriodSeconds
		}

		if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.TLSEnabled {
			vmsingle.Spec.Secrets = append(vmsingle.Spec.Secrets, victoriametrics.GetVmsingleTLSSecretName(cr.Spec.Victoriametrics.VmSingle))

			if vmsingle.Spec.ExtraArgs == nil {
				vmsingle.Spec.ExtraArgs = make(map[string]string)
			}
			maps.Copy(vmsingle.Spec.ExtraArgs, map[string]string{"tls": "true"})
			maps.Copy(vmsingle.Spec.ExtraArgs, map[string]string{"tlsCertFile": "/etc/vm/secrets/" + victoriametrics.GetVmsingleTLSSecretName(cr.Spec.Victoriametrics.VmSingle) + "/tls.crt"})
			maps.Copy(vmsingle.Spec.ExtraArgs, map[string]string{"tlsKeyFile": "/etc/vm/secrets/" + victoriametrics.GetVmsingleTLSSecretName(cr.Spec.Victoriametrics.VmSingle) + "/tls.key"})
		}

		// Set labels
		vmsingle.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(vmsingle.GetName(), vmsingle.GetNamespace())
		vmsingle.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Victoriametrics.VmSingle.Image)

		vmsingle.Spec.PodMetadata = &vmetricsv1b1.EmbeddedObjectMetadata{Labels: map[string]string{
			"name":                         utils.TruncLabel(vmsingle.GetName()),
			"app.kubernetes.io/name":       utils.TruncLabel(vmsingle.GetName()),
			"app.kubernetes.io/instance":   utils.GetInstanceLabel(vmsingle.GetName(), vmsingle.GetNamespace()),
			"app.kubernetes.io/component":  "victoriametrics",
			"app.kubernetes.io/part-of":    "monitoring",
			"app.kubernetes.io/managed-by": "monitoring-operator",
			"app.kubernetes.io/version":    utils.GetTagFromImage(cr.Spec.Victoriametrics.VmSingle.Image),
		}}

		if vmsingle.Spec.PodMetadata != nil {
			if cr.Spec.Victoriametrics.VmSingle.Labels != nil {
				for k, v := range cr.Spec.Victoriametrics.VmSingle.Labels {
					vmsingle.Spec.PodMetadata.Labels[k] = v
				}
			}

			if vmsingle.Spec.PodMetadata.Annotations == nil && cr.Spec.Victoriametrics.VmSingle.Annotations != nil {
				vmsingle.Spec.PodMetadata.Annotations = cr.Spec.Victoriametrics.VmSingle.Annotations
			} else {
				for k, v := range cr.Spec.Victoriametrics.VmSingle.Annotations {
					vmsingle.Spec.PodMetadata.Annotations[k] = v
				}
			}
		}

		if len(strings.TrimSpace(cr.Spec.Victoriametrics.VmSingle.PriorityClassName)) > 0 {
			vmsingle.Spec.PriorityClassName = cr.Spec.Victoriametrics.VmSingle.PriorityClassName
		}
	}

	return &vmsingle, nil
}

func vmSingleIngressV1beta1(cr *v1alpha1.PlatformMonitoring) (*v1beta1.Ingress, error) {
	ingress := v1beta1.Ingress{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmSingleIngressAsset), 100).Decode(&ingress); err != nil {
		return nil, err
	}
	//Set parameters
	ingress.SetGroupVersionKind(schema.GroupVersionKind{Group: "networking.k8s.io", Version: "v1beta1", Kind: "Ingress"})
	ingress.SetName(cr.GetNamespace() + "-" + utils.VmSingleServiceName)
	ingress.SetNamespace(cr.GetNamespace())

	if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.VmSingle.Ingress != nil && cr.Spec.Victoriametrics.VmSingle.Ingress.IsInstall() {
		// Check that ingress host is specified.
		if cr.Spec.Victoriametrics.VmSingle.Ingress.Host == "" {
			return nil, errors.New("host for ingress can not be empty")
		}

		// Add rule for vmsingle UI
		rule := v1beta1.IngressRule{Host: cr.Spec.Victoriametrics.VmSingle.Ingress.Host}
		serviceName := utils.VmSingleServiceName
		servicePort := intstr.FromInt(utils.VmSingleServicePort)

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
		if cr.Spec.Victoriametrics.VmSingle.Ingress.TLSSecretName != "" {
			ingress.Spec.TLS = []v1beta1.IngressTLS{
				{
					Hosts:      []string{cr.Spec.Victoriametrics.VmSingle.Ingress.Host},
					SecretName: cr.Spec.Victoriametrics.VmSingle.Ingress.TLSSecretName,
				},
			}
		}

		if cr.Spec.Victoriametrics.VmSingle.Ingress.IngressClassName != nil {
			ingress.Spec.IngressClassName = cr.Spec.Victoriametrics.VmSingle.Ingress.IngressClassName
		}

		// Set annotations
		ingress.SetAnnotations(cr.Spec.Victoriametrics.VmSingle.Ingress.Annotations)
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
		ingress.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Victoriametrics.VmSingle.Image)

		for lKey, lValue := range cr.Spec.Victoriametrics.VmSingle.Ingress.Labels {
			ingress.GetLabels()[lKey] = lValue
		}
	}
	return &ingress, nil
}

func vmSingleIngressV1(cr *v1alpha1.PlatformMonitoring) (*networkingv1.Ingress, error) {
	ingress := networkingv1.Ingress{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmSingleIngressAsset), 100).Decode(&ingress); err != nil {
		return nil, err
	}
	//Set parameters
	ingress.SetGroupVersionKind(schema.GroupVersionKind{Group: "networking.k8s.io", Version: "v1", Kind: "Ingress"})
	ingress.SetName(cr.GetNamespace() + "-" + utils.VmSingleServiceName)
	ingress.SetNamespace(cr.GetNamespace())

	if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.VmSingle.Ingress != nil && cr.Spec.Victoriametrics.VmSingle.Ingress.IsInstall() {
		// Check that ingress host is specified.
		if cr.Spec.Victoriametrics.VmSingle.Ingress.Host == "" {
			return nil, errors.New("host for ingress can not be empty")
		}

		pathType := networkingv1.PathTypePrefix
		// Add rule for vmsingle UI
		rule := networkingv1.IngressRule{Host: cr.Spec.Victoriametrics.VmSingle.Ingress.Host}
		rule.HTTP = &networkingv1.HTTPIngressRuleValue{
			Paths: []networkingv1.HTTPIngressPath{
				{
					Path:     "/",
					PathType: &pathType,
					Backend: networkingv1.IngressBackend{
						Service: &networkingv1.IngressServiceBackend{
							Name: utils.VmSingleServiceName,
							Port: networkingv1.ServiceBackendPort{
								Number: int32(utils.VmSingleServicePort),
							},
						},
					},
				},
			},
		}
		ingress.Spec.Rules = []networkingv1.IngressRule{rule}

		// Configure TLS if TLS secret name is set
		if cr.Spec.Victoriametrics.VmSingle.Ingress.TLSSecretName != "" {
			ingress.Spec.TLS = []networkingv1.IngressTLS{
				{
					Hosts:      []string{cr.Spec.Victoriametrics.VmSingle.Ingress.Host},
					SecretName: cr.Spec.Victoriametrics.VmSingle.Ingress.TLSSecretName,
				},
			}
		}

		if cr.Spec.Victoriametrics.VmSingle.Ingress.IngressClassName != nil {
			ingress.Spec.IngressClassName = cr.Spec.Victoriametrics.VmSingle.Ingress.IngressClassName
		}

		// Set annotations
		ingress.SetAnnotations(cr.Spec.Victoriametrics.VmSingle.Ingress.Annotations)
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
		ingress.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Victoriametrics.VmSingle.Image)

		for lKey, lValue := range cr.Spec.Victoriametrics.VmSingle.Ingress.Labels {
			ingress.GetLabels()[lKey] = lValue
		}
	}
	return &ingress, nil
}
