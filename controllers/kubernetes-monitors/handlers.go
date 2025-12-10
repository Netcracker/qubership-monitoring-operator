package kubernetes_monitors

import (
	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"k8s.io/apimachinery/pkg/api/errors"
)

func (r *KubernetesMonitorsReconciler) handleApiServerServiceMonitor(cr *monv1.PlatformMonitoring) error {
	m, err := kubernetesMonitorsApiServerServiceMonitor(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ApiServerServiceMonitor manifest")
		return err
	}
	e := &promv1.ServiceMonitor{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			if err = r.CreateResource(cr, m); err != nil {
				return err
			}
			return nil
		}
		return err
	}
	//Set parameters
	e.SetLabels(m.GetLabels())
	e.Spec.JobLabel = m.Spec.JobLabel
	e.Spec.Endpoints = m.Spec.Endpoints
	e.Spec.NamespaceSelector = m.Spec.NamespaceSelector
	e.Spec.Selector = m.Spec.Selector

	if err = r.UpdateResource(e); err != nil {
		return err
	}
	return nil
}

func (r *KubernetesMonitorsReconciler) handleKubeletServiceMonitor(cr *monv1.PlatformMonitoring) error {
	m, err := kubernetesMonitorsKubeletServiceMonitor(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating KubeletServiceMonitor manifest")
		return err
	}
	e := &promv1.ServiceMonitor{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			if err = r.CreateResource(cr, m); err != nil {
				return err
			}
			return nil
		}
		return err
	}

	//Set parameters
	e.SetLabels(m.GetLabels())
	e.Spec.JobLabel = m.Spec.JobLabel
	e.Spec.Endpoints = m.Spec.Endpoints
	e.Spec.NamespaceSelector = m.Spec.NamespaceSelector
	e.Spec.Selector = m.Spec.Selector

	if err = r.UpdateResource(e); err != nil {
		return err
	}
	return nil
}

func (r *KubernetesMonitorsReconciler) handleNginxIngressPodMonitor(cr *monv1.PlatformMonitoring) error {
	m, err := kubernetesMonitorsNginxIngressPodMonitor(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating NginxIngressPodMonitor manifest")
		return err
	}
	e := &promv1.PodMonitor{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			if err = r.CreateResource(cr, m); err != nil {
				return err
			}
			return nil
		}
		return err
	}

	//Set parameters
	e.SetLabels(m.GetLabels())
	e.Spec.JobLabel = m.Spec.JobLabel
	e.Spec.PodMetricsEndpoints = m.Spec.PodMetricsEndpoints
	e.Spec.NamespaceSelector = m.Spec.NamespaceSelector
	e.Spec.Selector = m.Spec.Selector

	if err = r.UpdateResource(e); err != nil {
		return err
	}
	return nil
}

func (r *KubernetesMonitorsReconciler) handleCoreDnsServiceMonitor(cr *monv1.PlatformMonitoring, isOpenshiftV4 bool) error {
	m, err := kubernetesMonitorsCoreDnsServiceMonitor(cr, isOpenshiftV4)
	if err != nil {
		r.Log.Error(err, "Failed creating CoreDnsServiceMonitor manifest")
		return err
	}
	e := &promv1.ServiceMonitor{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			if err = r.CreateResource(cr, m); err != nil {
				return err
			}
			return nil
		}
		return err
	}

	//Set parameters
	e.SetLabels(m.GetLabels())
	e.Spec.JobLabel = m.Spec.JobLabel
	e.Spec.Endpoints = m.Spec.Endpoints
	e.Spec.NamespaceSelector = m.Spec.NamespaceSelector
	e.Spec.Selector = m.Spec.Selector

	if err = r.UpdateResource(e); err != nil {
		return err
	}
	return nil
}

func (r *KubernetesMonitorsReconciler) handleOpenshiftServiceMonitor(cr *monv1.PlatformMonitoring,
	serviceMonitorAsset string,
	crServiceMonitorName string,
	metadataName string) error {
	m, err := openshiftServiceMonitor(cr, serviceMonitorAsset, crServiceMonitorName, metadataName)
	if err != nil {
		r.Log.Error(err, "Failed creating manifest for"+crServiceMonitorName)
		return err
	}
	e := &promv1.ServiceMonitor{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			if err = r.CreateResource(cr, m); err != nil {
				return err
			}
			return nil
		}
		return err
	}
	//Set parameters
	e.SetLabels(m.GetLabels())
	e.Spec.JobLabel = m.Spec.JobLabel
	e.Spec.Endpoints = m.Spec.Endpoints
	e.Spec.NamespaceSelector = m.Spec.NamespaceSelector
	e.Spec.Selector = m.Spec.Selector

	if err = r.UpdateResource(e); err != nil {
		return err
	}
	return nil
}

func (r *KubernetesMonitorsReconciler) deleteApiServerServiceMonitor(cr *monv1.PlatformMonitoring) error {
	m, err := kubernetesMonitorsApiServerServiceMonitor(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ApiServerServiceMonitor manifest")
		return err
	}
	e := &promv1.ServiceMonitor{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if err = r.DeleteResource(e); err != nil {
		return err
	}
	return nil
}

func (r *KubernetesMonitorsReconciler) deleteKubeletServiceMonitor(cr *monv1.PlatformMonitoring) error {
	m, err := kubernetesMonitorsKubeletServiceMonitor(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating KubeletServiceMonitor manifest")
		return err
	}
	e := &promv1.ServiceMonitor{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if err = r.DeleteResource(e); err != nil {
		return err
	}
	return nil
}

func (r *KubernetesMonitorsReconciler) deleteCoreDnsServiceMonitor(cr *monv1.PlatformMonitoring, isOpenshiftV4 bool) error {
	m, err := kubernetesMonitorsCoreDnsServiceMonitor(cr, isOpenshiftV4)
	if err != nil {
		r.Log.Error(err, "Failed creating CoreDnsServiceMonitor manifest")
		return err
	}
	e := &promv1.ServiceMonitor{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if err = r.DeleteResource(e); err != nil {
		return err
	}
	return nil
}

func (r *KubernetesMonitorsReconciler) deleteNginxIngressPodMonitor(cr *monv1.PlatformMonitoring) error {
	m, err := kubernetesMonitorsNginxIngressPodMonitor(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating NginxIngressPodMonitor manifest")
		return err
	}
	e := &promv1.PodMonitor{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if err = r.DeleteResource(e); err != nil {
		return err
	}
	return nil
}

func (r *KubernetesMonitorsReconciler) deleteOpenshiftServiceMonitor(cr *monv1.PlatformMonitoring,
	serviceMonitorAsset string,
	crServiceMonitorName string,
	metadataName string) error {
	m, err := openshiftServiceMonitor(cr, serviceMonitorAsset, crServiceMonitorName, metadataName)
	if err != nil {
		r.Log.Error(err, "Failed creating manifest for"+crServiceMonitorName)
		return err
	}
	e := &promv1.ServiceMonitor{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if err = r.DeleteResource(e); err != nil {
		return err
	}
	return nil
}
