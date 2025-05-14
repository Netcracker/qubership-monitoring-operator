package vmcluster

import (
	"embed"
	"errors"
	"maps"
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

func vmClusterServiceAccount(cr *v1alpha1.PlatformMonitoring) (*corev1.ServiceAccount, error) {
	sa := corev1.ServiceAccount{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmClusterServiceAccountAsset), 100).Decode(&sa); err != nil {
		return nil, err
	}
	//Set parameters
	sa.SetGroupVersionKind(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "ServiceAccount"})
	sa.SetName(cr.GetNamespace() + "-" + utils.VmClusterComponentName)
	sa.SetNamespace(cr.GetNamespace())

	return &sa, nil
}

func vmClusterClusterRole(cr *v1alpha1.PlatformMonitoring, hasPsp, hasScc bool) (*rbacv1.ClusterRole, error) {
	clusterRole := rbacv1.ClusterRole{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmClusterClusterRoleAsset), 100).Decode(&clusterRole); err != nil {
		return nil, err
	}
	//Set parameters
	clusterRole.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRole"})
	clusterRole.SetName(cr.GetNamespace() + "-" + utils.VmClusterComponentName)
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

func vmClusterClusterRoleBinding(cr *v1alpha1.PlatformMonitoring) (*rbacv1.ClusterRoleBinding, error) {
	clusterRoleBinding := rbacv1.ClusterRoleBinding{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmClusterClusterRoleBindingAsset), 100).Decode(&clusterRoleBinding); err != nil {
		return nil, err
	}
	//Set parameters
	clusterRoleBinding.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRoleBinding"})
	clusterRoleBinding.SetName(cr.GetNamespace() + "-" + utils.VmClusterComponentName)
	clusterRoleBinding.RoleRef.Name = cr.GetNamespace() + "-" + utils.VmClusterComponentName

	// Set namespace for all subjects
	for it := range clusterRoleBinding.Subjects {
		sub := &clusterRoleBinding.Subjects[it]
		sub.Namespace = cr.GetNamespace()
		sub.Name = cr.GetNamespace() + "-" + utils.VmClusterComponentName
	}
	return &clusterRoleBinding, nil
}

func vmCluster(cr *v1alpha1.PlatformMonitoring) (*vmetricsv1b1.VMCluster, error) {
	var err error
	vmcluster := vmetricsv1b1.VMCluster{}
	if err = yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmClusterAsset), 100).Decode(&vmcluster); err != nil {
		return nil, err
	}

	// Set parameters
	vmcluster.SetNamespace(cr.GetNamespace())

	if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.VmCluster.IsInstall() {
		if len(cr.Spec.Victoriametrics.VmCluster.RetentionPeriod) != 0 {
			vmcluster.Spec.RetentionPeriod = cr.Spec.Victoriametrics.VmCluster.RetentionPeriod
		}
		if cr.Spec.Victoriametrics.VmCluster.ReplicationFactor != nil {
			vmcluster.Spec.ReplicationFactor = cr.Spec.Victoriametrics.VmCluster.ReplicationFactor
		}

		vmcluster.Spec.ServiceAccountName = cr.GetNamespace() + "-" + utils.VmClusterComponentName

		// Set labels
		vmcluster.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(vmcluster.GetName(), vmcluster.GetNamespace())

		if len(strings.TrimSpace(cr.Spec.Victoriametrics.VmCluster.ClusterVersion)) > 0 {
			vmcluster.Spec.ClusterVersion = cr.Spec.Victoriametrics.VmCluster.ClusterVersion
		}

		if cr.Spec.Victoriametrics.VmCluster.VmSelect != nil {
			vmcluster.Spec.VMSelect = cr.Spec.Victoriametrics.VmCluster.VmSelect
			if cr.Spec.Victoriametrics.VmReplicas != nil {
				vmcluster.Spec.VMSelect.ReplicaCount = cr.Spec.Victoriametrics.VmReplicas
			}
			vmcluster.Spec.VMSelect.Image.Repository, vmcluster.Spec.VMSelect.Image.Tag = utils.SplitImage(cr.Spec.Victoriametrics.VmCluster.VmSelectImage)
			if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.TLSEnabled {
				vmcluster.Spec.VMSelect.Secrets = append(vmcluster.Spec.VMSelect.Secrets, victoriametrics.GetVmselectTLSSecretName(cr.Spec.Victoriametrics.VmCluster))
				if vmcluster.Spec.VMSelect.ExtraArgs == nil {
					vmcluster.Spec.VMSelect.ExtraArgs = make(map[string]string)
				}
				maps.Copy(vmcluster.Spec.VMSelect.ExtraArgs, map[string]string{"tls": "true"})
				maps.Copy(vmcluster.Spec.VMSelect.ExtraArgs, map[string]string{"tlsCertFile": "/etc/vm/secrets/" + victoriametrics.GetVmselectTLSSecretName(cr.Spec.Victoriametrics.VmCluster) + "/tls.crt"})
				maps.Copy(vmcluster.Spec.VMSelect.ExtraArgs, map[string]string{"tlsKeyFile": "/etc/vm/secrets/" + victoriametrics.GetVmselectTLSSecretName(cr.Spec.Victoriametrics.VmCluster) + "/tls.key"})
			}
			vmcluster.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(vmcluster.Spec.VMSelect.GetNameWithPrefix(cr.Name), vmcluster.GetNamespace())
			vmcluster.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Victoriametrics.VmCluster.VmSelectImage)

			vmcluster.Spec.VMSelect.PodMetadata = &vmetricsv1b1.EmbeddedObjectMetadata{Labels: map[string]string{
				"name":                         utils.TruncLabel(vmcluster.Spec.VMSelect.GetNameWithPrefix(cr.Name)),
				"app.kubernetes.io/name":       utils.TruncLabel(vmcluster.Spec.VMSelect.GetNameWithPrefix(cr.Name)),
				"app.kubernetes.io/instance":   utils.GetInstanceLabel(vmcluster.Spec.VMSelect.GetNameWithPrefix(cr.Name), vmcluster.GetNamespace()),
				"app.kubernetes.io/component":  "victoriametrics",
				"app.kubernetes.io/part-of":    "monitoring",
				"app.kubernetes.io/managed-by": "monitoring-operator",
				"app.kubernetes.io/version":    utils.GetTagFromImage(cr.Spec.Victoriametrics.VmCluster.VmSelectImage),
			}}
		}

		if cr.Spec.Victoriametrics.VmCluster.VmStorage != nil {
			vmcluster.Spec.VMStorage = cr.Spec.Victoriametrics.VmCluster.VmStorage
			if cr.Spec.Victoriametrics.VmReplicas != nil {
				vmcluster.Spec.VMStorage.ReplicaCount = cr.Spec.Victoriametrics.VmReplicas
			}
			vmcluster.Spec.VMStorage.Image.Repository, vmcluster.Spec.VMStorage.Image.Tag = utils.SplitImage(cr.Spec.Victoriametrics.VmCluster.VmStorageImage)
			if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.TLSEnabled {
				vmcluster.Spec.VMStorage.Secrets = append(vmcluster.Spec.VMStorage.Secrets, victoriametrics.GetVmstorageTLSSecretName(cr.Spec.Victoriametrics.VmCluster))

				if vmcluster.Spec.VMStorage.ExtraArgs == nil {
					vmcluster.Spec.VMStorage.ExtraArgs = make(map[string]string)
				}
				maps.Copy(vmcluster.Spec.VMStorage.ExtraArgs, map[string]string{"tls": "true"})
				maps.Copy(vmcluster.Spec.VMStorage.ExtraArgs, map[string]string{"tlsCertFile": "/etc/vm/secrets/" + victoriametrics.GetVmstorageTLSSecretName(cr.Spec.Victoriametrics.VmCluster) + "/tls.crt"})
				maps.Copy(vmcluster.Spec.VMStorage.ExtraArgs, map[string]string{"tlsKeyFile": "/etc/vm/secrets/" + victoriametrics.GetVmstorageTLSSecretName(cr.Spec.Victoriametrics.VmCluster) + "/tls.key"})
			}
			vmcluster.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(vmcluster.Spec.VMSelect.GetNameWithPrefix(cr.Name), vmcluster.GetNamespace())
			vmcluster.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Victoriametrics.VmCluster.VmStorageImage)

			vmcluster.Spec.VMStorage.PodMetadata = &vmetricsv1b1.EmbeddedObjectMetadata{Labels: map[string]string{
				"name":                         utils.TruncLabel(vmcluster.Spec.VMStorage.GetNameWithPrefix(cr.Name)),
				"app.kubernetes.io/name":       utils.TruncLabel(vmcluster.Spec.VMStorage.GetNameWithPrefix(cr.Name)),
				"app.kubernetes.io/instance":   utils.GetInstanceLabel(vmcluster.Spec.VMStorage.GetNameWithPrefix(cr.Name), vmcluster.GetNamespace()),
				"app.kubernetes.io/component":  "victoriametrics",
				"app.kubernetes.io/part-of":    "monitoring",
				"app.kubernetes.io/managed-by": "monitoring-operator",
				"app.kubernetes.io/version":    utils.GetTagFromImage(cr.Spec.Victoriametrics.VmCluster.VmStorageImage),
			}}
		}

		if cr.Spec.Victoriametrics.VmCluster.VmInsert != nil {
			vmcluster.Spec.VMInsert = cr.Spec.Victoriametrics.VmCluster.VmInsert
			if cr.Spec.Victoriametrics.VmReplicas != nil {
				vmcluster.Spec.VMInsert.ReplicaCount = cr.Spec.Victoriametrics.VmReplicas
			}
			vmcluster.Spec.VMInsert.Image.Repository, vmcluster.Spec.VMInsert.Image.Tag = utils.SplitImage(cr.Spec.Victoriametrics.VmCluster.VmInsertImage)
			if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.TLSEnabled {
				vmcluster.Spec.VMInsert.Secrets = append(vmcluster.Spec.VMInsert.Secrets, victoriametrics.GetVminsertTLSSecretName(cr.Spec.Victoriametrics.VmCluster))

				if vmcluster.Spec.VMInsert.ExtraArgs == nil {
					vmcluster.Spec.VMInsert.ExtraArgs = make(map[string]string)
				}
				maps.Copy(vmcluster.Spec.VMInsert.ExtraArgs, map[string]string{"tls": "true"})
				maps.Copy(vmcluster.Spec.VMInsert.ExtraArgs, map[string]string{"tlsCertFile": "/etc/vm/secrets/" + victoriametrics.GetVminsertTLSSecretName(cr.Spec.Victoriametrics.VmCluster) + "/tls.crt"})
				maps.Copy(vmcluster.Spec.VMInsert.ExtraArgs, map[string]string{"tlsKeyFile": "/etc/vm/secrets/" + victoriametrics.GetVminsertTLSSecretName(cr.Spec.Victoriametrics.VmCluster) + "/tls.key"})
			}
			vmcluster.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(vmcluster.Spec.VMInsert.GetNameWithPrefix(cr.Name), vmcluster.GetNamespace())
			vmcluster.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Victoriametrics.VmCluster.VmInsertImage)

			vmcluster.Spec.VMInsert.PodMetadata = &vmetricsv1b1.EmbeddedObjectMetadata{Labels: map[string]string{
				"name":                         utils.TruncLabel(vmcluster.Spec.VMInsert.GetNameWithPrefix(cr.Name)),
				"app.kubernetes.io/name":       utils.TruncLabel(vmcluster.Spec.VMInsert.GetNameWithPrefix(cr.Name)),
				"app.kubernetes.io/instance":   utils.GetInstanceLabel(vmcluster.Spec.VMInsert.GetNameWithPrefix(cr.Name), vmcluster.GetNamespace()),
				"app.kubernetes.io/component":  "victoriametrics",
				"app.kubernetes.io/part-of":    "monitoring",
				"app.kubernetes.io/managed-by": "monitoring-operator",
				"app.kubernetes.io/version":    utils.GetTagFromImage(cr.Spec.Victoriametrics.VmCluster.VmInsertImage),
			}}
		}
	}

	return &vmcluster, nil
}

func vmSelectIngressV1beta1(cr *v1alpha1.PlatformMonitoring) (*v1beta1.Ingress, error) {
	ingress := v1beta1.Ingress{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmSelectIngressAsset), 100).Decode(&ingress); err != nil {
		return nil, err
	}
	//Set parameters
	ingress.SetGroupVersionKind(schema.GroupVersionKind{Group: "networking.k8s.io", Version: "v1beta1", Kind: "Ingress"})
	ingress.SetName(cr.GetNamespace() + "-" + utils.VmSelectServiceName)
	ingress.SetNamespace(cr.GetNamespace())

	if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.VmCluster.VmSelectIngress != nil && cr.Spec.Victoriametrics.VmCluster.VmSelectIngress.IsInstall() {
		// Check that ingress host is specified.
		if cr.Spec.Victoriametrics.VmCluster.VmSelectIngress.Host == "" {
			return nil, errors.New("host for ingress can not be empty")
		}

		// Add rule for vmselect UI
		rule := v1beta1.IngressRule{Host: cr.Spec.Victoriametrics.VmCluster.VmSelectIngress.Host}
		serviceName := utils.VmSelectServiceName
		servicePort := intstr.FromInt(utils.VmSelectServicePort)

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
		if cr.Spec.Victoriametrics.VmCluster.VmSelectIngress.TLSSecretName != "" {
			ingress.Spec.TLS = []v1beta1.IngressTLS{
				{
					Hosts:      []string{cr.Spec.Victoriametrics.VmCluster.VmSelectIngress.Host},
					SecretName: cr.Spec.Victoriametrics.VmCluster.VmSelectIngress.TLSSecretName,
				},
			}
		}

		if cr.Spec.Victoriametrics.VmCluster.VmSelectIngress.IngressClassName != nil {
			ingress.Spec.IngressClassName = cr.Spec.Victoriametrics.VmCluster.VmSelectIngress.IngressClassName
		}

		// Set annotations
		ingress.SetAnnotations(cr.Spec.Victoriametrics.VmCluster.VmSelectIngress.Annotations)
		if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.TLSEnabled {
			if ingress.GetAnnotations() == nil {
				annotation := make(map[string]string)
				annotation["nginx.ingress.kubernetes.io/backend-protocol"] = "HTTPS"
				ingress.SetAnnotations(annotation)
			} else {
				ingress.GetAnnotations()["nginx.ingress.kubernetes.io/backend-protocol"] = "HTTPS"
			}
		}

		if ingress.GetAnnotations() == nil {
			annotation := make(map[string]string)
			annotation["nginx.ingress.kubernetes.io/app-root"] = "/select/0/vmui"
			ingress.SetAnnotations(annotation)
		} else {
			ingress.GetAnnotations()["nginx.ingress.kubernetes.io/app-root"] = "/select/0/vmui"
		}

		// Set labels with saving default labels
		ingress.Labels["name"] = utils.TruncLabel(ingress.GetName())
		ingress.Labels["app.kubernetes.io/name"] = utils.TruncLabel(ingress.GetName())
		ingress.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(ingress.GetName(), ingress.GetNamespace())
		ingress.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Victoriametrics.VmCluster.VmSelectImage)

		for lKey, lValue := range cr.Spec.Victoriametrics.VmCluster.VmSelectIngress.Labels {
			ingress.GetLabels()[lKey] = lValue
		}
	}
	return &ingress, nil
}

func vmSelectIngressV1(cr *v1alpha1.PlatformMonitoring) (*networkingv1.Ingress, error) {
	ingress := networkingv1.Ingress{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmSelectIngressAsset), 100).Decode(&ingress); err != nil {
		return nil, err
	}
	//Set parameters
	ingress.SetGroupVersionKind(schema.GroupVersionKind{Group: "networking.k8s.io", Version: "v1", Kind: "Ingress"})
	ingress.SetName(cr.GetNamespace() + "-" + utils.VmSelectServiceName)
	ingress.SetNamespace(cr.GetNamespace())

	if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.VmCluster.VmSelectIngress != nil && cr.Spec.Victoriametrics.VmCluster.VmSelectIngress.IsInstall() {
		// Check that ingress host is specified.
		if cr.Spec.Victoriametrics.VmCluster.VmSelectIngress.Host == "" {
			return nil, errors.New("host for ingress can not be empty")
		}

		pathType := networkingv1.PathTypePrefix
		// Add rule for vmselect UI
		rule := networkingv1.IngressRule{Host: cr.Spec.Victoriametrics.VmCluster.VmSelectIngress.Host}
		rule.HTTP = &networkingv1.HTTPIngressRuleValue{
			Paths: []networkingv1.HTTPIngressPath{
				{
					Path:     "/",
					PathType: &pathType,
					Backend: networkingv1.IngressBackend{
						Service: &networkingv1.IngressServiceBackend{
							Name: utils.VmSelectServiceName,
							Port: networkingv1.ServiceBackendPort{
								Number: int32(utils.VmSelectServicePort),
							},
						},
					},
				},
			},
		}
		ingress.Spec.Rules = []networkingv1.IngressRule{rule}

		// Configure TLS if TLS secret name is set
		if cr.Spec.Victoriametrics.VmCluster.VmSelectIngress.TLSSecretName != "" {
			ingress.Spec.TLS = []networkingv1.IngressTLS{
				{
					Hosts:      []string{cr.Spec.Victoriametrics.VmCluster.VmSelectIngress.Host},
					SecretName: cr.Spec.Victoriametrics.VmCluster.VmSelectIngress.TLSSecretName,
				},
			}
		}

		if cr.Spec.Victoriametrics.VmCluster.VmSelectIngress.IngressClassName != nil {
			ingress.Spec.IngressClassName = cr.Spec.Victoriametrics.VmCluster.VmSelectIngress.IngressClassName
		}

		// Set annotations
		ingress.SetAnnotations(cr.Spec.Victoriametrics.VmCluster.VmSelectIngress.Annotations)
		if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.TLSEnabled {
			if ingress.GetAnnotations() == nil {
				annotation := make(map[string]string)
				annotation["nginx.ingress.kubernetes.io/backend-protocol"] = "HTTPS"
				ingress.SetAnnotations(annotation)
			} else {
				ingress.GetAnnotations()["nginx.ingress.kubernetes.io/backend-protocol"] = "HTTPS"
			}
		}

		if ingress.GetAnnotations() == nil {
			annotation := make(map[string]string)
			annotation["nginx.ingress.kubernetes.io/app-root"] = "/select/0/vmui"
			ingress.SetAnnotations(annotation)
		} else {
			ingress.GetAnnotations()["nginx.ingress.kubernetes.io/app-root"] = "/select/0/vmui"
		}

		// Set labels with saving default labels
		ingress.Labels["name"] = utils.TruncLabel(ingress.GetName())
		ingress.Labels["app.kubernetes.io/name"] = utils.TruncLabel(ingress.GetName())
		ingress.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(ingress.GetName(), ingress.GetNamespace())
		ingress.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Victoriametrics.VmCluster.VmSelectImage)

		for lKey, lValue := range cr.Spec.Victoriametrics.VmCluster.VmSelectIngress.Labels {
			ingress.GetLabels()[lKey] = lValue
		}
	}
	return &ingress, nil
}
