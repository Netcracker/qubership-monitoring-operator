package kubernetes_monitors

import (
	"embed"

	"maps"

	v1beta1 "github.com/Netcracker/qubership-monitoring-operator/api"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

//go:embed  assets/*.yaml
var assets embed.FS

// Common K8s monitors
func kubernetesMonitorsApiServerServiceMonitor(cr *v1beta1.PlatformMonitoring) (*promv1.ServiceMonitor, error) {
	sm := promv1.ServiceMonitor{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.ApiServerServiceMonitorAsset), 100).Decode(&sm); err != nil {
		return nil, err
	}
	//Set parameters
	sm.SetName(cr.GetNamespace() + "-" + "kube-apiserver-service-monitor")
	sm.SetNamespace(cr.GetNamespace())

	if cr.Spec.KubernetesMonitors != nil {
		monitor, ok := cr.Spec.KubernetesMonitors[utils.ApiserverServiceMonitorName]
		if ok && monitor.IsInstall() {
			monitor.OverrideServiceMonitor(&sm)
		}
	}

	// Set labels
	sm.Labels["name"] = utils.TruncLabel(sm.GetName())
	sm.Labels["app.kubernetes.io/name"] = utils.TruncLabel(sm.GetName())
	sm.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(sm.GetName(), sm.GetNamespace())
	if label, ok := cr.Labels["app.kubernetes.io/version"]; ok {
		sm.Labels["app.kubernetes.io/version"] = label
	}
	if cr.GetLabels() != nil {
		for k, v := range cr.GetLabels() {
			if _, ok := sm.Labels[k]; !ok {
				sm.Labels[k] = v
			}
		}
	}

	if sm.Annotations == nil && cr.GetAnnotations() != nil {
		sm.SetAnnotations(cr.GetAnnotations())
	} else {
		for k, v := range cr.GetAnnotations() {
			sm.Annotations[k] = v
		}
	}

	return &sm, nil
}

func kubernetesMonitorsControllerManagerServiceMonitor(cr *v1beta1.PlatformMonitoring) (*promv1.ServiceMonitor, error) {
	sm := promv1.ServiceMonitor{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.KubeControllerManagerServiceMonitorAsset), 100).Decode(&sm); err != nil {
		return nil, err
	}
	//Set parameters
	sm.SetName(cr.GetNamespace() + "-" + "kube-controller-manager-service-monitor")
	sm.SetNamespace(cr.GetNamespace())

	if cr.Spec.KubernetesMonitors != nil {
		monitor, ok := cr.Spec.KubernetesMonitors[utils.KubeControllerManagerServiceMonitorName]
		if ok && monitor.IsInstall() {
			monitor.OverrideServiceMonitor(&sm)
		}
	}

	// Set labels
	sm.Labels["name"] = utils.TruncLabel(sm.GetName())
	sm.Labels["app.kubernetes.io/name"] = utils.TruncLabel(sm.GetName())
	sm.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(sm.GetName(), sm.GetNamespace())
	if label, ok := cr.Labels["app.kubernetes.io/version"]; ok {
		sm.Labels["app.kubernetes.io/version"] = label
	}
	if cr.GetLabels() != nil {
		for k, v := range cr.GetLabels() {
			if _, ok := sm.Labels[k]; !ok {
				sm.Labels[k] = v
			}
		}
	}

	if sm.Annotations == nil && cr.GetAnnotations() != nil {
		sm.SetAnnotations(cr.GetAnnotations())
	} else {
		for k, v := range cr.GetAnnotations() {
			sm.Annotations[k] = v
		}
	}

	return &sm, nil
}

func kubernetesMonitorsSchedulerServiceMonitor(cr *v1beta1.PlatformMonitoring) (*promv1.ServiceMonitor, error) {
	sm := promv1.ServiceMonitor{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.KubeSchedulerServiceMonitorAsset), 100).Decode(&sm); err != nil {
		return nil, err
	}
	//Set parameters
	sm.SetName(cr.GetNamespace() + "-" + "kube-scheduler-service-monitor")
	sm.SetNamespace(cr.GetNamespace())

	if cr.Spec.KubernetesMonitors != nil {
		monitor, ok := cr.Spec.KubernetesMonitors[utils.KubeSchedulerServiceMonitorName]
		if ok && monitor.IsInstall() {
			monitor.OverrideServiceMonitor(&sm)
		}
	}

	// Set labels
	sm.Labels["name"] = utils.TruncLabel(sm.GetName())
	sm.Labels["app.kubernetes.io/name"] = utils.TruncLabel(sm.GetName())
	sm.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(sm.GetName(), sm.GetNamespace())
	if label, ok := cr.Labels["app.kubernetes.io/version"]; ok {
		sm.Labels["app.kubernetes.io/version"] = label
	}
	if cr.GetLabels() != nil {
		for k, v := range cr.GetLabels() {
			if _, ok := sm.Labels[k]; !ok {
				sm.Labels[k] = v
			}
		}
	}

	if sm.Annotations == nil && cr.GetAnnotations() != nil {
		sm.SetAnnotations(cr.GetAnnotations())
	} else {
		for k, v := range cr.GetAnnotations() {
			sm.Annotations[k] = v
		}
	}

	return &sm, nil
}

func kubernetesMonitorsKubeletServiceMonitor(cr *v1beta1.PlatformMonitoring) (*promv1.ServiceMonitor, error) {
	sm := promv1.ServiceMonitor{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.KubeletServiceMonitorAsset), 100).Decode(&sm); err != nil {
		return nil, err
	}
	//Set parameters
	sm.SetName(cr.GetNamespace() + "-" + "kubelet-service-monitor")
	sm.SetNamespace(cr.GetNamespace())

	if cr.Spec.KubernetesMonitors != nil {
		monitor, ok := cr.Spec.KubernetesMonitors[utils.KubeletServiceMonitorName]
		if ok && monitor.IsInstall() {
			monitor.OverrideServiceMonitor(&sm)
		}
	}
	sm.Spec.NamespaceSelector.MatchNames = []string{cr.GetNamespace()}

	// Set labels
	sm.Labels["name"] = utils.TruncLabel(sm.GetName())
	sm.Labels["app.kubernetes.io/name"] = utils.TruncLabel(sm.GetName())
	sm.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(sm.GetName(), sm.GetNamespace())
	if label, ok := cr.Labels["app.kubernetes.io/version"]; ok {
		sm.Labels["app.kubernetes.io/version"] = label
	}
	if cr.GetLabels() != nil {
		for k, v := range cr.GetLabels() {
			if _, ok := sm.Labels[k]; !ok {
				sm.Labels[k] = v
			}
		}
	}

	if sm.Annotations == nil && cr.GetAnnotations() != nil {
		sm.SetAnnotations(cr.GetAnnotations())
	} else {
		for k, v := range cr.GetAnnotations() {
			sm.Annotations[k] = v
		}
	}

	return &sm, nil
}

func kubernetesMonitorsCoreDnsServiceMonitor(cr *v1beta1.PlatformMonitoring, isOpenshiftV4 bool) (*promv1.ServiceMonitor, error) {
	sm := promv1.ServiceMonitor{}

	CoreDnsServiceMonitorAsset := utils.CoreDnsServiceMonitorAssetK8s

	if isOpenshiftV4 {
		CoreDnsServiceMonitorAsset = utils.CoreDnsServiceMonitorAssetOs4
	}

	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, CoreDnsServiceMonitorAsset), 100).Decode(&sm); err != nil {
		return nil, err
	}
	//Set parameters
	sm.SetName(cr.GetNamespace() + "-" + "core-dns-service-monitor")
	sm.SetNamespace(cr.GetNamespace())

	if cr.Spec.KubernetesMonitors != nil {
		monitor, ok := cr.Spec.KubernetesMonitors[utils.CoreDnsServiceMonitorName]
		if ok && monitor.IsInstall() {
			monitor.OverrideServiceMonitor(&sm)
		}
	}

	// Set labels
	sm.Labels["name"] = utils.TruncLabel(sm.GetName())
	sm.Labels["app.kubernetes.io/name"] = utils.TruncLabel(sm.GetName())
	sm.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(sm.GetName(), sm.GetNamespace())
	if label, ok := cr.Labels["app.kubernetes.io/version"]; ok {
		sm.Labels["app.kubernetes.io/version"] = label
	}
	if cr.GetLabels() != nil {
		for k, v := range cr.GetLabels() {
			if _, ok := sm.Labels[k]; !ok {
				sm.Labels[k] = v
			}
		}
	}

	if sm.Annotations == nil && cr.GetAnnotations() != nil {
		sm.SetAnnotations(cr.GetAnnotations())
	} else {
		for k, v := range cr.GetAnnotations() {
			sm.Annotations[k] = v
		}
	}

	return &sm, nil
}

func kubernetesMonitorsNginxIngressPodMonitor(cr *v1beta1.PlatformMonitoring) (*promv1.PodMonitor, error) {
	pm := promv1.PodMonitor{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.NginxIngressPodMonitorAsset), 100).Decode(&pm); err != nil {
		return nil, err
	}
	//Set parameters
	pm.SetName(cr.GetNamespace() + "-" + "nginx-ingress-pod-monitor")
	pm.SetNamespace(cr.GetNamespace())

	if cr.Spec.KubernetesMonitors != nil {
		monitor, ok := cr.Spec.KubernetesMonitors[utils.NginxIngressPodMonitorName]
		if ok && monitor.IsInstall() {
			monitor.OverridePodMonitor(&pm)
		}
	}
	// Set labels
	pm.Labels["name"] = utils.TruncLabel(pm.GetName())
	pm.Labels["app.kubernetes.io/name"] = utils.TruncLabel(pm.GetName())
	pm.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(pm.GetName(), pm.GetNamespace())
	if label, ok := cr.Labels["app.kubernetes.io/version"]; ok {
		pm.Labels["app.kubernetes.io/version"] = label
	}
	if cr.GetLabels() != nil {
		for k, v := range cr.GetLabels() {
			if _, ok := pm.Labels[k]; !ok {
				pm.Labels[k] = v
			}
		}
	}

	if pm.Annotations == nil && cr.GetAnnotations() != nil {
		pm.SetAnnotations(cr.GetAnnotations())
	} else {
		for k, v := range cr.GetAnnotations() {
			pm.Annotations[k] = v
		}
	}

	return &pm, nil
}

// Openshift specific monitors
func openshiftServiceMonitor(cr *v1beta1.PlatformMonitoring,
	serviceMonitorAsset string,
	crServiceMonitorName string,
	metadataName string) (*promv1.ServiceMonitor, error) {
	sm := promv1.ServiceMonitor{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, serviceMonitorAsset), 100).Decode(&sm); err != nil {
		return nil, err
	}
	//Set parameters
	sm.SetName(cr.GetNamespace() + "-" + metadataName)
	sm.SetNamespace(cr.GetNamespace())

	if cr.Spec.KubernetesMonitors != nil {
		monitor, ok := cr.Spec.KubernetesMonitors[crServiceMonitorName]
		if ok && monitor.IsInstall() {
			monitor.OverrideServiceMonitor(&sm)
		}
	}

	// Set labels
	sm.Labels["name"] = utils.TruncLabel(sm.GetName())
	sm.Labels["app.kubernetes.io/name"] = utils.TruncLabel(sm.GetName())
	sm.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(sm.GetName(), sm.GetNamespace())
	if label, ok := cr.Labels["app.kubernetes.io/version"]; ok {
		sm.Labels["app.kubernetes.io/version"] = label
	}
	if cr.GetLabels() != nil {
		for k, v := range cr.GetLabels() {
			if _, ok := sm.Labels[k]; !ok {
				sm.Labels[k] = v
			}
		}
	}

	if sm.Annotations == nil && cr.GetAnnotations() != nil {
		sm.SetAnnotations(cr.GetAnnotations())
	} else {
		maps.Copy(sm.Annotations, cr.GetAnnotations())
	}

	return &sm, nil
}
