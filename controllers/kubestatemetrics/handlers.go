package kubestatemetrics

import (
	v1alpha1 "github.com/Netcracker/qubership-monitoring-operator/api/v1alpha1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
)

func (r *KubeStateMetricsReconciler) handleServiceAccount(cr *v1alpha1.PlatformMonitoring) error {
	m, err := kubeStateMetricsServiceAccount(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ServiceAccount manifest")
		return err
	}

	// Set labels
	m.Labels["name"] = utils.TruncLabel(m.GetName())
	m.Labels["app.kubernetes.io/name"] = utils.TruncLabel(m.GetName())
	m.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(m.GetName(), m.GetNamespace())
	m.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.KubeStateMetrics.Image)

	e := &corev1.ServiceAccount{ObjectMeta: m.ObjectMeta}
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

	if err = r.UpdateResource(e); err != nil {
		return err
	}
	return nil
}

func (r *KubeStateMetricsReconciler) handleClusterRole(cr *v1alpha1.PlatformMonitoring) error {
	m, err := kubeStateMetricsClusterRole(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ClusterRole manifest")
		return err
	}

	// Set labels
	m.Labels["name"] = utils.TruncLabel(m.GetName())
	m.Labels["app.kubernetes.io/name"] = utils.TruncLabel(m.GetName())
	m.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(m.GetName(), m.GetNamespace())
	m.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.KubeStateMetrics.Image)

	e := &rbacv1.ClusterRole{ObjectMeta: m.ObjectMeta}
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
	e.SetName(m.GetName())
	e.Rules = m.Rules

	if err = r.UpdateResource(e); err != nil {
		return err
	}
	return nil
}

func (r *KubeStateMetricsReconciler) handleClusterRoleBinding(cr *v1alpha1.PlatformMonitoring) error {
	m, err := kubeStateMetricsClusterRoleBinding(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ClusterRoleBinding manifest")
		return err
	}

	// Set labels
	m.Labels["name"] = utils.TruncLabel(m.GetName())
	m.Labels["app.kubernetes.io/name"] = utils.TruncLabel(m.GetName())
	m.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(m.GetName(), m.GetNamespace())
	m.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.KubeStateMetrics.Image)

	e := &rbacv1.ClusterRoleBinding{ObjectMeta: m.ObjectMeta}
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

	if err = r.UpdateResource(e); err != nil {
		return err
	}
	return nil
}

func (r *KubeStateMetricsReconciler) handleDeployment(cr *v1alpha1.PlatformMonitoring) error {
	m, err := kubeStateMetricsDeployment(cr, r.HasIngressV1Api() || r.HasIngressV1beta1Api())
	if err != nil {
		r.Log.Error(err, "Failed creating Deployment manifest")
		return err
	}
	e := &appsv1.Deployment{ObjectMeta: m.ObjectMeta}
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
	e.Spec.Template.SetLabels(m.Spec.Template.GetLabels())
	e.Spec.Template.Spec.SecurityContext = m.Spec.Template.Spec.SecurityContext
	e.Spec.Template.Spec.Containers = m.Spec.Template.Spec.Containers
	e.Spec.Template.Spec.ServiceAccountName = m.Spec.Template.Spec.ServiceAccountName
	e.Spec.Template.Spec.NodeSelector = m.Spec.Template.Spec.NodeSelector
	e.Spec.Template.Spec.Affinity = m.Spec.Template.Spec.Affinity
	e.Spec.Template.Spec.PriorityClassName = m.Spec.Template.Spec.PriorityClassName

	if err = r.UpdateResource(e); err != nil {
		return err
	}
	return nil
}

func (r *KubeStateMetricsReconciler) handleService(cr *v1alpha1.PlatformMonitoring) error {
	m, err := kubeStateMetricsService(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Service manifest")
		return err
	}

	// Set labels
	m.Labels["name"] = utils.TruncLabel(m.GetName())
	m.Labels["app.kubernetes.io/name"] = utils.TruncLabel(m.GetName())
	m.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(m.GetName(), m.GetNamespace())
	m.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.KubeStateMetrics.Image)

	e := &corev1.Service{ObjectMeta: m.ObjectMeta}
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
	e.Spec.Ports = m.Spec.Ports
	e.Spec.Selector = m.Spec.Selector

	if err = r.UpdateResource(e); err != nil {
		return err
	}
	return nil
}

func (r *KubeStateMetricsReconciler) handleServiceMonitor(cr *v1alpha1.PlatformMonitoring) error {
	m, err := kubeStateMetricsServiceMonitor(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ServiceMonitor manifest")
		return err
	}

	// Set labels
	m.Labels["name"] = utils.TruncLabel(m.GetName())
	m.Labels["app.kubernetes.io/name"] = utils.TruncLabel(m.GetName())
	m.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(m.GetName(), m.GetNamespace())
	m.Labels["app.kubernetes.io/version"] = utils.GetTagFromImage(cr.Spec.KubeStateMetrics.Image)

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

func (r *KubeStateMetricsReconciler) deleteServiceAccount(cr *v1alpha1.PlatformMonitoring) error {
	m, err := kubeStateMetricsServiceAccount(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ServiceAccount manifest")
		return err
	}
	e := &corev1.ServiceAccount{ObjectMeta: m.ObjectMeta}
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

func (r *KubeStateMetricsReconciler) deleteClusterRole(cr *v1alpha1.PlatformMonitoring) error {
	m, err := kubeStateMetricsClusterRole(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ClusterRole manifest")
		return err
	}
	e := &rbacv1.ClusterRole{ObjectMeta: m.ObjectMeta}
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

func (r *KubeStateMetricsReconciler) deleteClusterRoleBinding(cr *v1alpha1.PlatformMonitoring) error {
	m, err := kubeStateMetricsClusterRoleBinding(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ClusterRoleBinding manifest")
		return err
	}
	e := &rbacv1.ClusterRoleBinding{ObjectMeta: m.ObjectMeta}
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

func (r *KubeStateMetricsReconciler) deleteDeployment(cr *v1alpha1.PlatformMonitoring) error {
	m, err := kubeStateMetricsDeployment(cr, r.HasIngressV1Api() || r.HasIngressV1beta1Api())
	if err != nil {
		r.Log.Error(err, "Failed creating Deployment manifest")
		return err
	}
	e := &appsv1.Deployment{ObjectMeta: m.ObjectMeta}
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

func (r *KubeStateMetricsReconciler) deleteService(cr *v1alpha1.PlatformMonitoring) error {
	m, err := kubeStateMetricsService(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Service manifest")
		return err
	}
	e := &corev1.Service{ObjectMeta: m.ObjectMeta}
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

func (r *KubeStateMetricsReconciler) deleteServiceMonitor(cr *v1alpha1.PlatformMonitoring) error {
	m, err := kubeStateMetricsServiceMonitor(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ServiceMonitor manifest")
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
