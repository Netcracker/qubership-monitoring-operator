package etcd

import (
	"context"

	v1beta1 "github.com/Netcracker/qubership-monitoring-operator/api/v1beta1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"

	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *EtcdMonitorReconciler) handleServiceAccount(cr *v1beta1.PlatformMonitoring) error {
	m, err := nodeExporterServiceAccount(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ServiceAccount manifest")
		return err
	}
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

func (r *EtcdMonitorReconciler) handleServiceMonitor(cr *v1beta1.PlatformMonitoring, etcdServiceNamespace string, isOpenshiftV4 bool) error {
	m, err := etcdServiceMonitor(cr, etcdServiceNamespace, isOpenshiftV4)
	if err != nil {
		r.Log.Error(err, "Failed creating ServiceMonitor manifest")
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

func (r *EtcdMonitorReconciler) handleService(cr *v1beta1.PlatformMonitoring, isOpenshift bool, etcdServiceNamespace string, isOpenshiftV4 bool) error {
	m, err := etcdService(isOpenshift, etcdServiceNamespace, isOpenshiftV4)
	if err != nil {
		r.Log.Error(err, "Failed creating Service manifest")
		return err
	}
	var e *corev1.Service
	if e, err = r.KubeClient.CoreV1().Services(etcdServiceNamespace).Get(context.TODO(), m.GetName(), metav1.GetOptions{}); err != nil {
		if errors.IsNotFound(err) {
			if err = r.CreateResource(cr, m); err != nil {
				return err
			}
			return nil
		}
		return err
	}

	//Set parameters
	e.TypeMeta = m.TypeMeta
	e.SetLabels(m.GetLabels())
	e.Spec.Ports = m.Spec.Ports
	e.Spec.Selector = m.Spec.Selector

	if err = r.UpdateResource(e); err != nil {
		return err
	}
	return nil
}

func (r *EtcdMonitorReconciler) handleSecret(cr *v1beta1.PlatformMonitoring) error {
	m, err := etcdSecret(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Secret manifest")
		return err
	}

	// Set labels
	m.Labels["app.kubernetes.io/instance"] = utils.GetInstanceLabel(m.GetName(), m.GetNamespace())
	if label, ok := cr.Labels["app.kubernetes.io/version"]; ok {
		m.Labels["app.kubernetes.io/version"] = label
	}

	e := &corev1.Secret{ObjectMeta: m.ObjectMeta}
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

func (r *EtcdMonitorReconciler) deleteServiceMonitor(cr *v1beta1.PlatformMonitoring, etcdServiceNamespace string, isOpenshiftV4 bool) error {
	m, err := etcdServiceMonitor(cr, etcdServiceNamespace, isOpenshiftV4)
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
