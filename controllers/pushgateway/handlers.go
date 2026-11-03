package pushgateway

import (
	"maps"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *PushgatewayReconciler) handleDeployment(cr *monv1.PlatformMonitoring) error {
	isOpenShift, err := r.IsOpenShift()
	if err != nil {
		return err
	}
	m, err := pushgatewayDeployment(cr, isOpenShift)
	if err != nil {
		r.Log.Error(err, "Failed creating Deployment manifest")
		return err
	}
	e := &appsv1.Deployment{ObjectMeta: m.ObjectMeta}
	return utils.CreateOrUpdateWorkload(r.ComponentReconciler, cr, m, e, syncDeployment)
}

// syncDeployment copies the Deployment fields this controller manages.
func syncDeployment(existing, desired *appsv1.Deployment) bool {
	return utils.AnyChanged(
		utils.SetLabelsIfChanged(existing, desired.GetLabels()),
		utils.SetLabelsIfChanged(&existing.Spec.Template, desired.Spec.Template.GetLabels()),
		utils.SetAnnotationsIfChanged(existing, annotationsForDeploymentUpdate(existing.GetAnnotations(), desired.GetAnnotations())),
		utils.SetAnnotationsIfChanged(&existing.Spec.Template, desired.Spec.Template.GetAnnotations()),
		utils.SyncPodSpec(&existing.Spec.Template.Spec, &desired.Spec.Template.Spec),
		utils.SetIfChanged(&existing.Spec.Template.Spec.NodeSelector, desired.Spec.Template.Spec.NodeSelector),
		utils.SetIfChanged(&existing.Spec.Replicas, desired.Spec.Replicas),
	)
}

const deploymentRevisionAnnotation = "deployment.kubernetes.io/revision"

// annotationsForDeploymentUpdate keeps full-replace semantics for controller-owned
// annotations (including removal) while copying kube's revision so idle skip holds.
func annotationsForDeploymentUpdate(live, desired map[string]string) map[string]string {
	ann := maps.Clone(desired)
	rev, ok := live[deploymentRevisionAnnotation]
	if !ok {
		return ann
	}
	if ann == nil {
		ann = map[string]string{}
	}
	ann[deploymentRevisionAnnotation] = rev
	return ann
}

func (r *PushgatewayReconciler) handlePVC(cr *monv1.PlatformMonitoring) error {
	m, err := pushgatewayPVC(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating PVC manifest")
		return err
	}
	e := &corev1.PersistentVolumeClaim{ObjectMeta: m.ObjectMeta}
	if err = r.GetResource(e); err == nil {
		r.Log.Info("PVC with name pushgateway is already exist. Skip reconciling")
	} else {
		if errors.IsNotFound(err) {
			if err = r.CreateResource(cr, m); err != nil {
				return err
			}
			return nil
		}
		return err
	}
	return nil
}

func (r *PushgatewayReconciler) handleService(cr *monv1.PlatformMonitoring) error {
	m, err := pushgatewayService(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Service manifest")
		return err
	}

	e := &corev1.Service{ObjectMeta: m.ObjectMeta}
	return utils.CreateOrUpdate(r.ComponentReconciler, cr, m, e, utils.SyncService)
}

func (r *PushgatewayReconciler) handleIngressV1(cr *monv1.PlatformMonitoring) error {
	m, err := pushgatewayIngressV1(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Ingress manifest")
		return err
	}
	e := &networkingv1.Ingress{ObjectMeta: m.ObjectMeta}
	return utils.CreateOrUpdate(r.ComponentReconciler, cr, m, e, utils.SyncIngress)
}

func (r *PushgatewayReconciler) handleServiceMonitor(cr *monv1.PlatformMonitoring) error {
	m, err := pushgatewayServiceMonitor(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating ServiceMonitor manifest")
		return err
	}

	e := &promv1.ServiceMonitor{ObjectMeta: m.ObjectMeta}
	return utils.CreateOrUpdate(r.ComponentReconciler, cr, m, e, utils.SyncServiceMonitor)
}

func (r *PushgatewayReconciler) deleteDeployment(cr *monv1.PlatformMonitoring) error {
	// Deletion targets are addressed by their stable name and namespace so that uninstall does not
	// depend on platform discovery or on validation of the desired state.
	e := &appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Name: utils.PushgatewayComponentName, Namespace: cr.GetNamespace()}}
	if err := r.GetResource(e); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if err := r.DeleteResource(e); err != nil {
		return err
	}
	return nil
}

func (r *PushgatewayReconciler) deleteService(cr *monv1.PlatformMonitoring) error {
	m, err := pushgatewayService(cr)
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

func (r *PushgatewayReconciler) deleteIngressV1(cr *monv1.PlatformMonitoring) error {
	m, err := pushgatewayIngressV1(cr)
	if err != nil {
		r.Log.Error(err, "Failed creating Ingress manifest")
		return err
	}
	e := &networkingv1.Ingress{ObjectMeta: m.ObjectMeta}
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

func (r *PushgatewayReconciler) deleteServiceMonitor(cr *monv1.PlatformMonitoring) error {
	m, err := pushgatewayServiceMonitor(cr)
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
