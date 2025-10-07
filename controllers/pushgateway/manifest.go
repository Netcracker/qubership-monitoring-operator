package pushgateway

import (
	"embed"
	"strings"

	v1alpha1 "github.com/Netcracker/qubership-monitoring-operator/api/v1alpha1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/yaml"
)

//go:embed  assets/*.yaml
var assets embed.FS

func pushgatewayDeployment(cr *v1alpha1.PlatformMonitoring) (*appsv1.Deployment, error) {
	d := appsv1.Deployment{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.PushgatewayDeploymentAsset), 100).Decode(&d); err != nil {
		return nil, err
	}
	//Set parameters
	d.SetGroupVersionKind(schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"})
	d.SetName(utils.PushgatewayComponentName)
	d.SetNamespace(cr.GetNamespace())

	if cr.Spec.Pushgateway != nil {
		// Set Pushgateway replicas
		if cr.Spec.Pushgateway.Replicas != nil {
			d.Spec.Replicas = cr.Spec.Pushgateway.Replicas
		}
		// Find container with c.Name as name and set Image from custom resource
		for it := range d.Spec.Template.Spec.Containers {
			c := &d.Spec.Template.Spec.Containers[it]
			if c.Name == utils.PushgatewayComponentName {
				c.Image = cr.Spec.Pushgateway.Image
				if len(cr.Spec.Pushgateway.ExtraArgs) > 0 {
					c.Args = cr.Spec.Pushgateway.ExtraArgs
				}
				portValue := cr.Spec.Pushgateway.Port
				for p := range c.Ports {
					port := &c.Ports[p]
					if port.Name == utils.PushgatewayPortName {
						port.HostPort = portValue
						port.ContainerPort = portValue
					}
				}
				if cr.Spec.Pushgateway.Resources.Size() > 0 {
					c.Resources = cr.Spec.Pushgateway.Resources
				}
				// Set volumeMounts
				if cr.Spec.Pushgateway.VolumeMounts != nil {
					c.VolumeMounts = cr.Spec.Pushgateway.VolumeMounts
				}
				// Set volumeMount and flags for work with PVC
				if cr.Spec.Pushgateway.Storage.Size() > 0 {
					mountPath := utils.PushgatewayPVCVolumeMountMountPath
					pvcVolumeMount := corev1.VolumeMount{
						Name:      utils.PushgatewayStorageVolumeName,
						MountPath: mountPath,
					}
					// If volumeMounts already contains mount with the same name as mount for PVC, replace it
					if c.VolumeMounts != nil {
						isStorageVolumePresent := false
						for i, vm := range c.VolumeMounts {
							if vm.Name == utils.PushgatewayStorageVolumeName {
								c.VolumeMounts[i] = pvcVolumeMount
								isStorageVolumePresent = true
								break
							}
						}
						if !isStorageVolumePresent {
							c.VolumeMounts = append(c.VolumeMounts, pvcVolumeMount)
						}
					} else {
						c.VolumeMounts = []corev1.VolumeMount{
							pvcVolumeMount,
						}
					}
					// Set persistence flags if these flags are not presented
					isPersistenceFilePresent := false
					isPersistenceIntervalPresent := false
					for _, arg := range c.Args {
						if strings.Contains(arg, "persistence.file") {
							isPersistenceFilePresent = true
						}
						if strings.Contains(arg, "persistence.interval") {
							isPersistenceIntervalPresent = true
						}
					}
					if !isPersistenceFilePresent {
						c.Args = append(c.Args, "--persistence.file="+mountPath+"/"+utils.PushgatewayPersistenceFile)
					}
					if !isPersistenceIntervalPresent {
						c.Args = append(c.Args, "--persistence.interval="+utils.PushgatewayPersistenceInterval)
					}
				}
				break
			}
		}
		// Set security context
		if cr.Spec.Pushgateway.SecurityContext != nil {
			if d.Spec.Template.Spec.SecurityContext == nil {
				d.Spec.Template.Spec.SecurityContext = &corev1.PodSecurityContext{}
			}
			if cr.Spec.Pushgateway.SecurityContext.RunAsUser != nil {
				d.Spec.Template.Spec.SecurityContext.RunAsUser = cr.Spec.Pushgateway.SecurityContext.RunAsUser
			}
			if cr.Spec.Pushgateway.SecurityContext.FSGroup != nil {
				d.Spec.Template.Spec.SecurityContext.FSGroup = cr.Spec.Pushgateway.SecurityContext.FSGroup
			}
		}
		// Set tolerations for pushgateway
		if cr.Spec.Pushgateway.Tolerations != nil {
			d.Spec.Template.Spec.Tolerations = cr.Spec.Pushgateway.Tolerations
		}
		// Set NodeSelector for pushgateway
		if cr.Spec.Pushgateway.NodeSelector != nil {
			d.Spec.Template.Spec.NodeSelector = cr.Spec.Pushgateway.NodeSelector
		}
		// Set affinity for pushgateway
		if cr.Spec.Pushgateway.Affinity != nil {
			d.Spec.Template.Spec.Affinity = cr.Spec.Pushgateway.Affinity
		}
		// Set volumes
		if cr.Spec.Pushgateway.Volumes != nil {
			d.Spec.Template.Spec.Volumes = cr.Spec.Pushgateway.Volumes
		}
		// Set volume for work with PVC
		if cr.Spec.Pushgateway.Storage.Size() > 0 {
			pvcVolume := corev1.Volume{
				Name: utils.PushgatewayStorageVolumeName,
				VolumeSource: corev1.VolumeSource{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: utils.PushgatewayComponentName,
					},
				},
			}
			// If volumes already contains volume with the same name as volume for PVC, replace it
			if d.Spec.Template.Spec.Volumes != nil {
				isStorageVolumePresent := false
				for i, v := range d.Spec.Template.Spec.Volumes {
					if v.Name == utils.PushgatewayStorageVolumeName {
						d.Spec.Template.Spec.Volumes[i] = pvcVolume
						isStorageVolumePresent = true
						break
					}
				}
				if !isStorageVolumePresent {
					d.Spec.Template.Spec.Volumes = append(d.Spec.Template.Spec.Volumes, pvcVolume)
				}
			} else {
				d.Spec.Template.Spec.Volumes = []corev1.Volume{
					pvcVolume,
				}
			}
		}
		// Set annotations and labels
		if cr.Spec.Pushgateway.Annotations != nil {
			if d.Annotations == nil {
				d.SetAnnotations(cr.Spec.Pushgateway.Annotations)
			} else {
				for k, v := range cr.Spec.Pushgateway.Annotations {
					d.Annotations[k] = v
				}
			}
			if d.Spec.Template.Annotations == nil {
				d.Spec.Template.Annotations = cr.Spec.Pushgateway.Annotations
			} else {
				for k, v := range cr.Spec.Pushgateway.Annotations {
					d.Spec.Template.Annotations[k] = v
				}
			}
		}

		d.Labels["name"] = utils.TruncLabel(d.GetName())
		d.Labels["app.kubernetes.io/name"] = utils.TruncLabel(d.GetName())
		d.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(d.GetName(), d.GetNamespace())
		d.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Pushgateway.Image)

		d.Spec.Template.Labels["name"] = utils.TruncLabel(d.GetName())
		d.Spec.Template.Labels["app.kubernetes.io/name"] = utils.TruncLabel(d.GetName())
		d.Spec.Template.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(d.GetName(), d.GetNamespace())
		d.Spec.Template.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Pushgateway.Image)

		if cr.Spec.Pushgateway.Labels != nil {
			for k, v := range cr.Spec.Pushgateway.Labels {
				d.Labels[k] = v
				d.Spec.Template.Labels[k] = v
			}
		}

		if len(strings.TrimSpace(cr.Spec.Pushgateway.PriorityClassName)) > 0 {
			d.Spec.Template.Spec.PriorityClassName = cr.Spec.Pushgateway.PriorityClassName
		}
	}

	return &d, nil
}

func pushgatewayPVC(cr *v1alpha1.PlatformMonitoring) (*corev1.PersistentVolumeClaim, error) {
	pvc := corev1.PersistentVolumeClaim{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.PushgatewayPVCAsset), 100).Decode(&pvc); err != nil {
		return nil, err
	}
	//Set parameters
	pvc.SetGroupVersionKind(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "PersistentVolumeClaim"})
	pvc.SetName(utils.PushgatewayComponentName)
	pvc.SetNamespace(cr.GetNamespace())

	// Set labels
	pvc.Labels["name"] = utils.TruncLabel(pvc.GetName())
	pvc.Labels["app.kubernetes.io/name"] = utils.TruncLabel(pvc.GetName())
	pvc.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(pvc.GetName(), pvc.GetNamespace())
	pvc.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Pushgateway.Image)

	if cr.Spec.Pushgateway != nil && cr.Spec.Pushgateway.Storage != nil {
		// Set PVC spec
		pvc.Spec = *cr.Spec.Pushgateway.Storage
	}
	return &pvc, nil
}

func pushgatewayService(cr *v1alpha1.PlatformMonitoring) (*corev1.Service, error) {
	service := corev1.Service{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.PushgatewayServiceAsset), 100).Decode(&service); err != nil {
		return nil, err
	}
	//Set parameters
	service.SetGroupVersionKind(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Service"})
	service.SetName(utils.PushgatewayComponentName)
	service.SetNamespace(cr.GetNamespace())

	if cr.Spec.Pushgateway != nil {
		// Set port
		for p := range service.Spec.Ports {
			port := &service.Spec.Ports[p]
			if port.Name == utils.PushgatewayPortName {
				port.Port = cr.Spec.Pushgateway.Port
				port.TargetPort = intstr.FromInt(int(cr.Spec.Pushgateway.Port))
			}
		}
	}
	return &service, nil
}

func pushgatewayIngressV1(cr *v1alpha1.PlatformMonitoring) (*networkingv1.Ingress, error) {
	ingress := networkingv1.Ingress{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.PushgatewayIngressAsset), 100).Decode(&ingress); err != nil {
		return nil, err
	}
	//Set metadata
	ingress.SetGroupVersionKind(schema.GroupVersionKind{Group: "networking.k8s.io", Version: "v1", Kind: "Ingress"})
	ingress.SetName(cr.GetNamespace() + "-" + utils.PushgatewayComponentName)
	ingress.SetNamespace(cr.GetNamespace())

	if cr.Spec.Pushgateway != nil && cr.Spec.Pushgateway.Ingress != nil && cr.Spec.Pushgateway.Ingress.IsInstall() {
		var rules []networkingv1.IngressRule
		pathType := networkingv1.PathTypePrefix
		ing := cr.Spec.Pushgateway.Ingress

		switch {
		// 1. If Host is provided
		case ing.Host != "":
			rules = append(rules, networkingv1.IngressRule{
				Host: ing.Host,
				IngressRuleValue: networkingv1.IngressRuleValue{
					HTTP: &networkingv1.HTTPIngressRuleValue{
						Paths: []networkingv1.HTTPIngressPath{defaultPushgatewayPath(pathType)},
					},
				},
			})
			// 2. If custom ingress rules provided
		case len(ing.Rules) > 0:
			for _, r := range ing.Rules {
				// fallback if HTTP is not set
				if r.HTTP == nil || len(r.HTTP.Paths) == 0 {
					r.HTTP = &v1alpha1.HTTPIngressRuleValue{
						Paths: []v1alpha1.IngressPath{
							{
								Path:     "/",
								PathType: string(pathType),
								Backend: v1alpha1.IngressPathBackend{
									Service: v1alpha1.IngressPathBackendService{
										Name: utils.PushgatewayComponentName,
										Port: v1alpha1.ServiceBackendPort{
											Name: utils.PushgatewayPortName,
										},
									},
								},
							},
						},
					}
				}
				// converting to k8s networkingv1
				var paths []networkingv1.HTTPIngressPath

				for _, p := range r.HTTP.Paths {
					pt := networkingv1.PathTypePrefix
					if p.PathType != "" {
						pt = networkingv1.PathType(p.PathType)
					}

					backendPort := networkingv1.ServiceBackendPort{}
					if p.Backend.Service.Port.Number != 0 {
						backendPort.Number = p.Backend.Service.Port.Number
					} else {
						backendPort.Name = p.Backend.Service.Port.Name
					}
					paths = append(paths, networkingv1.HTTPIngressPath{
						Path:     p.Path,
						PathType: &pt,
						Backend: networkingv1.IngressBackend{
							Service: &networkingv1.IngressServiceBackend{
								Name: p.Backend.Service.Name,
								Port: backendPort,
							},
						},
					})
				}
				rules = append(rules, networkingv1.IngressRule{
					Host: r.Host,
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{Paths: paths},
					},
				})
			}
			// 3. fallback: if no Host or no custom ingress rules provided
		default:
			rules = append(rules, networkingv1.IngressRule{
				IngressRuleValue: networkingv1.IngressRuleValue{
					HTTP: &networkingv1.HTTPIngressRuleValue{
						Paths: []networkingv1.HTTPIngressPath{defaultPushgatewayPath(pathType)},
					},
				},
			})
		}

		ingress.Spec.Rules = rules

		tlsConfigured := false
		pickSecret := func(ingressTLSSecret string, tlsCfg *v1alpha1.CommonTLSConfig) string {
			if ingressTLSSecret != "" {
				return ingressTLSSecret
			}
			if tlsCfg != nil {
				return tlsCfg.SecretName
			}
			return ""
		}
		// Configure tls if TLS config is defined
		if !tlsConfigured && len(cr.Spec.Pushgateway.Ingress.TLS) > 0 {
			for _, hostgroup := range cr.Spec.Pushgateway.Ingress.TLS {
				if len(hostgroup.Hosts) == 0 {
					continue
				}
				validHosts := make([]string, 0, len(hostgroup.Hosts))
				for _, h := range hostgroup.Hosts {
					if strings.TrimSpace(h) != "" {
						validHosts = append(validHosts, h)
					}
				}
				if len(validHosts) == 0 {
					continue
				}
				secret := hostgroup.SecretName
				// fallback: if secretName is empty - use TLSSecretName
				if secret == "" {
					secret = pickSecret(cr.Spec.Pushgateway.Ingress.TLSSecretName, cr.Spec.Pushgateway.TLSConfig)
				}
				if secret != "" {
					ingress.Spec.TLS = append(ingress.Spec.TLS, networkingv1.IngressTLS{
						Hosts:      validHosts,
						SecretName: secret,
					})
				}
			}
			if len(ingress.Spec.TLS) > 0 {
				tlsConfigured = true
			}
		}
		// Configure TLS if TLS secret name and host is set
		if !tlsConfigured && cr.Spec.Pushgateway.Ingress.Host != "" {
			secret := pickSecret(cr.Spec.Pushgateway.Ingress.TLSSecretName, cr.Spec.Pushgateway.TLSConfig)
			if secret != "" {
				ingress.Spec.TLS = []networkingv1.IngressTLS{
					{
						Hosts:      []string{cr.Spec.Pushgateway.Ingress.Host},
						SecretName: secret,
					},
				}
				tlsConfigured = true
			}
		}
		// Fallback: use ingress rules to configure tls hosts and TLSSecretName
		if !tlsConfigured && len(cr.Spec.Pushgateway.Ingress.Rules) > 0 {
			tlsHosts := []string{}
			secret := pickSecret(cr.Spec.Pushgateway.Ingress.TLSSecretName, cr.Spec.Pushgateway.TLSConfig)
			for _, rule := range cr.Spec.Pushgateway.Ingress.Rules {
				if rule.Host != "" {
					tlsHosts = append(tlsHosts, rule.Host)
				}
			}
			if len(tlsHosts) > 0 && secret != "" {
				ingress.Spec.TLS = []networkingv1.IngressTLS{
					{
						Hosts:      tlsHosts,
						SecretName: secret,
					},
				}
			}
		}

		if cr.Spec.Pushgateway.Ingress.IngressClassName != nil {
			ingress.Spec.IngressClassName = cr.Spec.Pushgateway.Ingress.IngressClassName
		}

		// Set annotations
		ingress.SetAnnotations(cr.Spec.Pushgateway.Ingress.Annotations)

		// Set labels with saving default labels
		ingress.Labels["name"] = utils.TruncLabel(ingress.GetName())
		ingress.Labels["app.kubernetes.io/name"] = utils.TruncLabel(ingress.GetName())
		ingress.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(ingress.GetName(), ingress.GetNamespace())
		ingress.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.Pushgateway.Image)

		for lKey, lValue := range cr.Spec.Pushgateway.Ingress.Labels {
			ingress.GetLabels()[lKey] = lValue
		}
	}
	return &ingress, nil
}

func pushgatewayServiceMonitor(cr *v1alpha1.PlatformMonitoring) (*promv1.ServiceMonitor, error) {
	sm := promv1.ServiceMonitor{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.PushgatewayServiceMonitorAsset), 100).Decode(&sm); err != nil {
		return nil, err
	}
	//Set parameters
	sm.SetGroupVersionKind(schema.GroupVersionKind{Group: "monitoring.coreos.com", Version: "v1", Kind: "ServiceMonitor"})
	sm.SetName(cr.GetNamespace() + "-" + utils.PushgatewayComponentName)
	sm.SetNamespace(cr.GetNamespace())

	if cr.Spec.Pushgateway != nil && cr.Spec.Pushgateway.ServiceMonitor != nil && cr.Spec.Pushgateway.ServiceMonitor.IsInstall() {
		cr.Spec.Pushgateway.ServiceMonitor.OverrideServiceMonitor(&sm)
	}
	sm.Spec.NamespaceSelector.MatchNames = []string{cr.GetNamespace()}

	return &sm, nil
}

func defaultPushgatewayPath(pathType networkingv1.PathType) networkingv1.HTTPIngressPath {
	return networkingv1.HTTPIngressPath{
		Path:     "/",
		PathType: &pathType,
		Backend: networkingv1.IngressBackend{
			Service: &networkingv1.IngressServiceBackend{
				Name: utils.PushgatewayComponentName,
				Port: networkingv1.ServiceBackendPort{
					Name: utils.PushgatewayPortName,
				},
			},
		},
	}
}
