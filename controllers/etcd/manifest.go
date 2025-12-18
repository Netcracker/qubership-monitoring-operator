package etcd

import (
	"embed"

	v1beta1 "github.com/Netcracker/qubership-monitoring-operator/api/v1beta1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/yaml"
)

//go:embed  assets/*.yaml
var assets embed.FS

func etcdServiceMonitor(cr *v1beta1.PlatformMonitoring, etcdServiceNamespace string, isOpenshiftV4 bool) (*promv1.ServiceMonitor, error) {
	sm := promv1.ServiceMonitor{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.EtcdServiceMonitorAsset), 100).Decode(&sm); err != nil {
		return nil, err
	}
	//Set parameters
	sm.SetGroupVersionKind(schema.GroupVersionKind{Group: "monitoring.coreos.com", Version: "v1", Kind: "ServiceMonitor"})
	sm.SetName(cr.GetNamespace() + "-" + "etcd-service-monitor")
	sm.SetNamespace(cr.GetNamespace())
	sm.Spec.NamespaceSelector.MatchNames = []string{etcdServiceNamespace}
	// Port "etcd-metrics" is used in OpenShift v4.x
	if isOpenshiftV4 {
		sm.Spec.Endpoints[0].Port = "etcd-metrics"
	}

	if cr.Spec.KubernetesMonitors != nil {
		monitor, ok := cr.Spec.KubernetesMonitors[utils.EtcdServiceMonitorName]
		if ok && monitor.IsInstall() {
			monitor.OverrideServiceMonitor(&sm)
		}
	}

	// Set labels
	sm.Labels["name"] = utils.TruncLabel(sm.GetName())
	sm.Labels["app.kubernetes.io/name"] = utils.TruncLabel(sm.GetName())
	sm.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(sm.GetName(), sm.GetNamespace())

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

func etcdService(isOpenshift bool, etcdServiceNamespace string, isOpenshiftV4 bool) (*corev1.Service, error) {
	service := corev1.Service{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.EtcdServiceComponentAsset), 100).Decode(&service); err != nil {
		return nil, err
	}
	//Set parameters
	service.SetGroupVersionKind(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Service"})
	service.SetName(utils.EtcdServiceComponentName)
	service.SetNamespace(etcdServiceNamespace)

	// Kubernetes uses "component: etcd" selector
	// OpenShift v3.x uses "openshift.io/component: etcd" selector
	if isOpenshift && !isOpenshiftV4 {
		service.Spec.Selector = map[string]string{"openshift.io/component": "etcd"}
	}
	// OpenShift v4.x uses "etcd: 'true'" selector
	if isOpenshiftV4 {
		service.Spec.Selector = map[string]string{"etcd": "true"}
	}
	// If cluster is not OpenShift v4.x, remove port "etcd-metrics"
	if !isOpenshiftV4 {
		service.Spec.Ports = service.Spec.Ports[:1]
	}
	// Set labels
	service.Labels["name"] = utils.TruncLabel(service.GetName())
	service.Labels["app.kubernetes.io/name"] = utils.TruncLabel(service.GetName())
	service.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(service.GetName(), service.GetNamespace())

	return &service, nil
}

func etcdSecret(cr *v1beta1.PlatformMonitoring) (*corev1.Secret, error) {
	secret := corev1.Secret{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.EtcdClientCertsSecretAsset), 100).Decode(&secret); err != nil {
		return nil, err
	}
	//Set parameters
	secret.SetGroupVersionKind(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Secret"})
	secret.SetNamespace(cr.GetNamespace())

	return &secret, nil
}

func nodeExporterServiceAccount(cr *v1beta1.PlatformMonitoring) (*corev1.ServiceAccount, error) {
	sa := corev1.ServiceAccount{}
	if err := yaml.NewYAMLOrJSONDecoder(utils.MustAssetReader(assets, utils.EtcdNodeExporterServiceAccountAsset), 100).Decode(&sa); err != nil {
		return nil, err
	}
	//Set parameters
	sa.SetGroupVersionKind(schema.GroupVersionKind{Group: "", Version: "v1", Kind: "ServiceAccount"})
	sa.SetName(cr.GetNamespace() + "-" + utils.NodeExporterComponentName)
	sa.SetNamespace(cr.GetNamespace())

	// Set labels
	sa.Labels["name"] = utils.TruncLabel(sa.GetName())
	sa.Labels["app.kubernetes.io/name"] = utils.TruncLabel(sa.GetName())
	sa.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(sa.GetName(), sa.GetNamespace())
	if cr.Spec.NodeExporter != nil {
		sa.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.NodeExporter.Image)
	}
	return &sa, nil
}
