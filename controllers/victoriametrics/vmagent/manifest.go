package vmagent

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

func vmAgentServiceAccount(cr *v1alpha1.PlatformMonitoring) (*corev1.ServiceAccount, error) {
	sa := corev1.ServiceAccount{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmAgentServiceAccountAsset), 100).Decode(&sa); err != nil {
		return nil, err
	}
	//Set parameters
	sa.SetGroupVersionKind(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "ServiceAccount"})
	sa.SetName(cr.GetNamespace() + "-" + utils.VmAgentComponentName)
	sa.SetNamespace(cr.GetNamespace())

	return &sa, nil
}

func vmAgentClusterRole(cr *v1alpha1.PlatformMonitoring, hasPsp, hasScc bool) (*rbacv1.ClusterRole, error) {
	clusterRole := rbacv1.ClusterRole{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmAgentClusterRoleAsset), 100).Decode(&clusterRole); err != nil {
		return nil, err
	}
	//Set parameters
	clusterRole.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRole"})
	clusterRole.SetName(cr.GetNamespace() + "-" + utils.VmAgentComponentName)
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

func vmAgentClusterRoleBinding(cr *v1alpha1.PlatformMonitoring) (*rbacv1.ClusterRoleBinding, error) {
	clusterRoleBinding := rbacv1.ClusterRoleBinding{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmAgentClusterRoleBindingAsset), 100).Decode(&clusterRoleBinding); err != nil {
		return nil, err
	}
	//Set parameters
	clusterRoleBinding.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRoleBinding"})
	clusterRoleBinding.SetName(cr.GetNamespace() + "-" + utils.VmAgentComponentName)
	clusterRoleBinding.RoleRef.Name = cr.GetNamespace() + "-" + utils.VmAgentComponentName

	// Set namespace for all subjects
	for it := range clusterRoleBinding.Subjects {
		sub := &clusterRoleBinding.Subjects[it]
		sub.Namespace = cr.GetNamespace()
		sub.Name = cr.GetNamespace() + "-" + utils.VmAgentComponentName
	}
	return &clusterRoleBinding, nil
}

func vmAgentRole(cr *v1alpha1.PlatformMonitoring) (*rbacv1.Role, error) {
	role := rbacv1.Role{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmAgentRoleAsset), 100).Decode(&role); err != nil {
		return nil, err
	}
	//Set parameters
	role.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "Role"})
	role.SetName(cr.GetNamespace() + "-" + utils.VmAgentComponentName)
	role.SetNamespace(cr.GetNamespace())

	return &role, nil
}

func vmAgentRoleBinding(cr *v1alpha1.PlatformMonitoring) (*rbacv1.RoleBinding, error) {
	roleBinding := rbacv1.RoleBinding{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmAgentRoleBindingAsset), 100).Decode(&roleBinding); err != nil {
		return nil, err
	}
	//Set parameters
	roleBinding.SetGroupVersionKind(schema.GroupVersionKind{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "RoleBinding"})
	roleBinding.SetName(cr.GetNamespace() + "-" + utils.VmAgentComponentName)
	roleBinding.SetNamespace(cr.GetNamespace())
	roleBinding.RoleRef.Name = cr.GetNamespace() + "-" + utils.VmAgentComponentName

	// Set namespace for all subjects
	for it := range roleBinding.Subjects {
		sub := &roleBinding.Subjects[it]
		sub.Namespace = cr.GetNamespace()
		sub.Name = cr.GetNamespace() + "-" + utils.VmAgentComponentName
	}
	return &roleBinding, nil
}

func vmAgent(r *VmAgentReconciler, cr *v1alpha1.PlatformMonitoring) (*vmetricsv1b1.VMAgent, error) {
	var err error
	vmagent := vmetricsv1b1.VMAgent{}
	if err = yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmAgentAsset), 100).Decode(&vmagent); err != nil {
		return nil, err
	}

	// Set parameters
	vmagent.SetNamespace(cr.GetNamespace())

	if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.VmAgent.IsInstall() {

		// Set Vmagent image
		vmagent.Spec.Image.Repository, vmagent.Spec.Image.Tag = utils.SplitImage(cr.Spec.Victoriametrics.VmAgent.Image)

		if r != nil {
			// Set security context
			if cr.Spec.Victoriametrics.VmAgent.SecurityContext != nil {
				if vmagent.Spec.SecurityContext == nil {
					vmagent.Spec.SecurityContext = &vmetricsv1b1.SecurityContext{}
				}
				if cr.Spec.Victoriametrics.VmAgent.SecurityContext.RunAsUser != nil {
					vmagent.Spec.SecurityContext.RunAsUser = cr.Spec.Victoriametrics.VmAgent.SecurityContext.RunAsUser
				}
				if cr.Spec.Victoriametrics.VmAgent.SecurityContext.FSGroup != nil {
					vmagent.Spec.SecurityContext.FSGroup = cr.Spec.Victoriametrics.VmAgent.SecurityContext.FSGroup
				}
			}
		}

		if cr.Spec.Victoriametrics.VmAgent.Replicas != nil {
			vmagent.Spec.ReplicaCount = cr.Spec.Victoriametrics.VmAgent.Replicas
		}
		if cr.Spec.Victoriametrics.VmCluster.IsInstall() && cr.Spec.Victoriametrics.VmReplicas != nil {
			vmagent.Spec.ReplicaCount = cr.Spec.Victoriametrics.VmReplicas
		}

		vmagent.Spec.ServiceAccountName = cr.GetNamespace() + "-" + utils.VmAgentComponentName

		// Set resources for Vmagent cr
		if cr.Spec.Victoriametrics.VmAgent.Resources.Size() > 0 {
			vmagent.Spec.Resources = cr.Spec.Victoriametrics.VmAgent.Resources
		}
		// Set secrets for VmAgent deployment
		if len(cr.Spec.Victoriametrics.VmAgent.Secrets) > 0 {
			vmagent.Spec.Secrets = cr.Spec.Victoriametrics.VmAgent.Secrets
		}

		// Set additional volumes
		if cr.Spec.Victoriametrics.VmAgent.Volumes != nil {
			vmagent.Spec.Volumes = cr.Spec.Victoriametrics.VmAgent.Volumes
		}
		// Set additional volumeMounts for each VmAgent container. The current container names are:
		// `vmagent`, `config-reloader`
		if cr.Spec.Victoriametrics.VmAgent.VolumeMounts != nil {
			for it := range vmagent.Spec.Containers {
				c := &vmagent.Spec.Containers[it]

				// Set additional volumeMounts only for VmAgent container
				if c.Name == utils.VmAgentComponentName {
					c.VolumeMounts = cr.Spec.Victoriametrics.VmAgent.VolumeMounts
				}
			}
		}
		// Set nodeSelector for Vmagent cr
		if cr.Spec.Victoriametrics.VmAgent.NodeSelector != nil {
			vmagent.Spec.NodeSelector = cr.Spec.Victoriametrics.VmAgent.NodeSelector
		}
		// Set affinity for Vmagent cr
		if cr.Spec.Victoriametrics.VmAgent.Affinity != nil {
			vmagent.Spec.Affinity = cr.Spec.Victoriametrics.VmAgent.Affinity
		}

		if len(strings.TrimSpace(cr.Spec.Victoriametrics.VmAgent.ScrapeInterval)) > 0 {
			vmagent.Spec.ScrapeInterval = cr.Spec.Victoriametrics.VmAgent.ScrapeInterval
		}

		if cr.Spec.Victoriametrics.VmAgent.MaxScrapeInterval != nil {
			vmagent.Spec.MaxScrapeInterval = cr.Spec.Victoriametrics.VmAgent.MaxScrapeInterval
		}

		if cr.Spec.Victoriametrics.VmAgent.MinScrapeInterval != nil {
			vmagent.Spec.MinScrapeInterval = cr.Spec.Victoriametrics.VmAgent.MinScrapeInterval
		}

		if cr.Spec.Victoriametrics.VmAgent.VMAgentExternalLabelName != nil {
			vmagent.Spec.VMAgentExternalLabelName = cr.Spec.Victoriametrics.VmAgent.VMAgentExternalLabelName
		}

		// Set external labels
		vmagent.Spec.ExternalLabels = cr.Spec.Victoriametrics.VmAgent.ExternalLabels

		if cr.Spec.Victoriametrics.VmAgent.ExtraArgs != nil {
			vmagent.Spec.ExtraArgs = cr.Spec.Victoriametrics.VmAgent.ExtraArgs
		}

		if cr.Spec.Victoriametrics.VmAgent.ExtraEnvs != nil {
			vmagent.Spec.ExtraEnvs = cr.Spec.Victoriametrics.VmAgent.ExtraEnvs
		}

		// Set external URL
		if cr.Spec.Victoriametrics.VmAgent.RemoteWrite != nil {
			vmagent.Spec.RemoteWrite = cr.Spec.Victoriametrics.VmAgent.RemoteWrite
		}

		if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.VmOperator.IsInstall() && cr.Spec.Victoriametrics.VmSingle.IsInstall() {
			vmSingle := vmetricsv1b1.VMSingle{}
			vmSingle.SetName(utils.VmComponentName)
			vmSingle.SetNamespace(cr.GetNamespace())
			vmagentRemoteWrite := vmetricsv1b1.VMAgentRemoteWriteSpec{}
			if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.TLSEnabled {
				if vmSingle.Spec.ExtraArgs == nil {
					vmSingle.Spec.ExtraArgs = make(map[string]string)
				}
				maps.Copy(vmSingle.Spec.ExtraArgs, map[string]string{"tls": "true"})
				vmagentRemoteWrite.TLSConfig = &vmetricsv1b1.TLSConfig{
					CAFile:   "/etc/vm/secrets/" + victoriametrics.GetVmagentTLSSecretName(cr.Spec.Victoriametrics.VmAgent) + "/ca.crt",
					CertFile: "/etc/vm/secrets/" + victoriametrics.GetVmagentTLSSecretName(cr.Spec.Victoriametrics.VmAgent) + "/tls.crt",
					KeyFile:  "/etc/vm/secrets/" + victoriametrics.GetVmagentTLSSecretName(cr.Spec.Victoriametrics.VmAgent) + "/tls.key",
				}
			}
			vmagentRemoteWrite.URL = vmSingle.AsURL() + "/api/v1/write"

			addVmSingle := true
			for _, rw := range vmagent.Spec.RemoteWrite {
				if rw.URL == vmagentRemoteWrite.URL {
					addVmSingle = false
					break
				}
			}
			if addVmSingle {
				vmagent.Spec.RemoteWrite = append(vmagent.Spec.RemoteWrite, vmagentRemoteWrite)
			}
		}

		if cr.Spec.Victoriametrics.VmOperator.IsInstall() && cr.Spec.Victoriametrics.VmCluster.IsInstall() {
			vmCluster := vmetricsv1b1.VMCluster{}
			vmCluster.SetName(utils.VmComponentName)
			vmCluster.SetNamespace(cr.GetNamespace())
			vmCluster.Spec.VMInsert = cr.Spec.Victoriametrics.VmCluster.VmInsert
			vmagentRemoteWrite := vmetricsv1b1.VMAgentRemoteWriteSpec{}
			if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.TLSEnabled {
				if vmCluster.Spec.VMInsert.ExtraArgs == nil {
					vmCluster.Spec.VMInsert.ExtraArgs = make(map[string]string)
				}
				maps.Copy(vmCluster.Spec.VMInsert.ExtraArgs, map[string]string{"tls": "true"})
				vmagentRemoteWrite.TLSConfig = &vmetricsv1b1.TLSConfig{
					CAFile:   "/etc/vm/secrets/" + victoriametrics.GetVmagentTLSSecretName(cr.Spec.Victoriametrics.VmAgent) + "/ca.crt",
					CertFile: "/etc/vm/secrets/" + victoriametrics.GetVmagentTLSSecretName(cr.Spec.Victoriametrics.VmAgent) + "/tls.crt",
					KeyFile:  "/etc/vm/secrets/" + victoriametrics.GetVmagentTLSSecretName(cr.Spec.Victoriametrics.VmAgent) + "/tls.key",
				}
			}
			vmagentRemoteWrite.URL = vmCluster.VMInsertURL() + "/insert/0/prometheus/api/v1/write"
			addVmInsert := true
			for _, rw := range vmagent.Spec.RemoteWrite {
				if rw.URL == vmagentRemoteWrite.URL {
					addVmInsert = false
					break
				}
			}
			if addVmInsert {
				vmagent.Spec.RemoteWrite = append(vmagent.Spec.RemoteWrite, vmagentRemoteWrite)
			}
		}

		if cr.Spec.Victoriametrics.VmAgent.RemoteWriteSettings != nil {
			vmagent.Spec.RemoteWriteSettings = cr.Spec.Victoriametrics.VmAgent.RemoteWriteSettings
		}

		if cr.Spec.Victoriametrics.VmAgent.RelabelConfig != nil {
			vmagent.Spec.RelabelConfig = cr.Spec.Victoriametrics.VmAgent.RelabelConfig
		}

		if cr.Spec.Victoriametrics.VmAgent.InlineRelabelConfig != nil {
			vmagent.Spec.InlineRelabelConfig = cr.Spec.Victoriametrics.VmAgent.InlineRelabelConfig
		}

		// Set additional containers
		if cr.Spec.Victoriametrics.VmAgent.Containers != nil {
			vmagent.Spec.Containers = cr.Spec.Victoriametrics.VmAgent.Containers
		}

		// Set tolerations for Vmagent
		if cr.Spec.Victoriametrics.VmAgent.Tolerations != nil {
			vmagent.Spec.Tolerations = cr.Spec.Victoriametrics.VmAgent.Tolerations
		}

		// Set selector for podMonitors
		if cr.Spec.Victoriametrics.VmAgent.PodMonitorSelector != nil {
			vmagent.Spec.PodScrapeSelector = cr.Spec.Victoriametrics.VmAgent.PodMonitorSelector
		}

		// Set selector for serviceMonitor
		if cr.Spec.Victoriametrics.VmAgent.ServiceMonitorSelector != nil {
			vmagent.Spec.ServiceScrapeSelector = cr.Spec.Victoriametrics.VmAgent.ServiceMonitorSelector
		}

		// Set namespaceSelector for podMonitors
		if cr.Spec.Victoriametrics.VmAgent.PodMonitorNamespaceSelector != nil {
			vmagent.Spec.PodScrapeNamespaceSelector = cr.Spec.Victoriametrics.VmAgent.PodMonitorNamespaceSelector
		}

		// Set namespaceSelector for serviceMonitor
		if cr.Spec.Victoriametrics.VmAgent.ServiceMonitorNamespaceSelector != nil {
			vmagent.Spec.ServiceScrapeNamespaceSelector = cr.Spec.Victoriametrics.VmAgent.ServiceMonitorNamespaceSelector
		}

		if cr.Spec.Victoriametrics.VmAgent.TerminationGracePeriodSeconds != nil {
			vmagent.Spec.TerminationGracePeriodSeconds = cr.Spec.Victoriametrics.VmAgent.TerminationGracePeriodSeconds
		}

		// Set labels
		vmagent.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(vmagent.GetName(), vmagent.GetNamespace())
		vmagent.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Victoriametrics.VmAgent.Image)

		vmagent.Spec.PodMetadata = &vmetricsv1b1.EmbeddedObjectMetadata{Labels: map[string]string{
			"name":                         utils.TruncLabel(vmagent.GetName()),
			"app.kubernetes.io/name":       utils.TruncLabel(vmagent.GetName()),
			"app.kubernetes.io/instance":   utils.GetInstanceLabel(vmagent.GetName(), vmagent.GetNamespace()),
			"app.kubernetes.io/component":  "victoriametrics",
			"app.kubernetes.io/part-of":    "monitoring",
			"app.kubernetes.io/managed-by": "monitoring-operator",
			"app.kubernetes.io/version":    utils.GetTagFromImage(cr.Spec.Victoriametrics.VmAgent.Image),
		}}

		if vmagent.Spec.PodMetadata != nil {
			if cr.Spec.Victoriametrics.VmAgent.Labels != nil {
				for k, v := range cr.Spec.Victoriametrics.VmAgent.Labels {
					vmagent.Spec.PodMetadata.Labels[k] = v
				}
			}

			if vmagent.Spec.PodMetadata.Annotations == nil && cr.Spec.Victoriametrics.VmAgent.Annotations != nil {
				vmagent.Spec.PodMetadata.Annotations = cr.Spec.Victoriametrics.VmAgent.Annotations
			} else {
				for k, v := range cr.Spec.Victoriametrics.VmAgent.Annotations {
					vmagent.Spec.PodMetadata.Annotations[k] = v
				}
			}
		}

		if cr.Spec.Victoriametrics.VmAgent.EnforcedNamespaceLabel != nil {
			vmagent.Spec.EnforcedNamespaceLabel = *cr.Spec.Victoriametrics.VmAgent.EnforcedNamespaceLabel
		}

		if len(strings.TrimSpace(cr.Spec.Victoriametrics.VmAgent.PriorityClassName)) > 0 {
			vmagent.Spec.PriorityClassName = cr.Spec.Victoriametrics.VmAgent.PriorityClassName
		}

		if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.TLSEnabled {
			vmagent.Spec.Secrets = append(vmagent.Spec.Secrets, victoriametrics.GetVmagentTLSSecretName(cr.Spec.Victoriametrics.VmAgent))

			if vmagent.Spec.ExtraArgs == nil {
				vmagent.Spec.ExtraArgs = make(map[string]string)
			}
			maps.Copy(vmagent.Spec.ExtraArgs, map[string]string{"tls": "true"})
			maps.Copy(vmagent.Spec.ExtraArgs, map[string]string{"tlsCertFile": "/etc/vm/secrets/" + victoriametrics.GetVmagentTLSSecretName(cr.Spec.Victoriametrics.VmAgent) + "/tls.crt"})
			maps.Copy(vmagent.Spec.ExtraArgs, map[string]string{"tlsKeyFile": "/etc/vm/secrets/" + victoriametrics.GetVmagentTLSSecretName(cr.Spec.Victoriametrics.VmAgent) + "/tls.key"})
		}
	}
	return &vmagent, nil
}

func vmAgentIngressV1beta1(cr *v1alpha1.PlatformMonitoring) (*v1beta1.Ingress, error) {
	ingress := v1beta1.Ingress{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmAgentIngressAsset), 100).Decode(&ingress); err != nil {
		return nil, err
	}
	//Set parameters
	ingress.SetGroupVersionKind(schema.GroupVersionKind{Group: "networking.k8s.io", Version: "v1beta1", Kind: "Ingress"})
	ingress.SetName(cr.GetNamespace() + "-" + utils.VmAgentServiceName)
	ingress.SetNamespace(cr.GetNamespace())

	if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.VmAgent.Ingress != nil && cr.Spec.Victoriametrics.VmAgent.Ingress.IsInstall() {
		// Check that ingress host is specified.
		if cr.Spec.Victoriametrics.VmAgent.Ingress.Host == "" {
			return nil, errors.New("host for ingress can not be empty")
		}

		// Add rule for vmagent UI
		rule := v1beta1.IngressRule{Host: cr.Spec.Victoriametrics.VmAgent.Ingress.Host}
		serviceName := utils.VmAgentServiceName
		servicePort := intstr.FromInt(utils.VmAgentServicePort)

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
		if cr.Spec.Victoriametrics.VmAgent.Ingress.TLSSecretName != "" {
			ingress.Spec.TLS = []v1beta1.IngressTLS{
				{
					Hosts:      []string{cr.Spec.Victoriametrics.VmAgent.Ingress.Host},
					SecretName: cr.Spec.Victoriametrics.VmAgent.Ingress.TLSSecretName,
				},
			}
		}

		if cr.Spec.Victoriametrics.VmAgent.Ingress.IngressClassName != nil {
			ingress.Spec.IngressClassName = cr.Spec.Victoriametrics.VmAgent.Ingress.IngressClassName
		}

		// Set annotations
		ingress.SetAnnotations(cr.Spec.Victoriametrics.VmAgent.Ingress.Annotations)
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
		ingress.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Victoriametrics.VmAgent.Image)

		for lKey, lValue := range cr.Spec.Victoriametrics.VmAgent.Ingress.Labels {
			ingress.GetLabels()[lKey] = lValue
		}
	}
	return &ingress, nil
}

func vmAgentIngressV1(cr *v1alpha1.PlatformMonitoring) (*networkingv1.Ingress, error) {
	ingress := networkingv1.Ingress{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.VmAgentIngressAsset), 100).Decode(&ingress); err != nil {
		return nil, err
	}
	//Set parameters
	ingress.SetGroupVersionKind(schema.GroupVersionKind{Group: "networking.k8s.io", Version: "v1", Kind: "Ingress"})
	ingress.SetName(cr.GetNamespace() + "-" + utils.VmAgentServiceName)
	ingress.SetNamespace(cr.GetNamespace())

	if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.VmAgent.Ingress != nil && cr.Spec.Victoriametrics.VmAgent.Ingress.IsInstall() {
		// Check that ingress host is specified.
		if cr.Spec.Victoriametrics.VmAgent.Ingress.Host == "" {
			return nil, errors.New("host for ingress can not be empty")
		}

		pathType := networkingv1.PathTypePrefix
		// Add rule for vmagent UI
		rule := networkingv1.IngressRule{Host: cr.Spec.Victoriametrics.VmAgent.Ingress.Host}
		rule.HTTP = &networkingv1.HTTPIngressRuleValue{
			Paths: []networkingv1.HTTPIngressPath{
				{
					Path:     "/",
					PathType: &pathType,
					Backend: networkingv1.IngressBackend{
						Service: &networkingv1.IngressServiceBackend{
							Name: utils.VmAgentServiceName,
							Port: networkingv1.ServiceBackendPort{
								Number: int32(utils.VmAgentServicePort),
							},
						},
					},
				},
			},
		}

		ingress.Spec.Rules = []networkingv1.IngressRule{rule}

		// Configure TLS if TLS secret name is set
		if cr.Spec.Victoriametrics.VmAgent.Ingress.TLSSecretName != "" {
			ingress.Spec.TLS = []networkingv1.IngressTLS{
				{
					Hosts:      []string{cr.Spec.Victoriametrics.VmAgent.Ingress.Host},
					SecretName: cr.Spec.Victoriametrics.VmAgent.Ingress.TLSSecretName,
				},
			}
		}

		if cr.Spec.Victoriametrics.VmAgent.Ingress.IngressClassName != nil {
			ingress.Spec.IngressClassName = cr.Spec.Victoriametrics.VmAgent.Ingress.IngressClassName
		}

		// Set annotations
		ingress.SetAnnotations(cr.Spec.Victoriametrics.VmAgent.Ingress.Annotations)
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
		ingress.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Victoriametrics.VmAgent.Image)

		for lKey, lValue := range cr.Spec.Victoriametrics.VmAgent.Ingress.Labels {
			ingress.GetLabels()[lKey] = lValue
		}
	}
	return &ingress, nil
}
