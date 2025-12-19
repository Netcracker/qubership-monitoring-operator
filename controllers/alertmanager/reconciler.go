package alertmanager

import (
	v1beta1 "github.com/Netcracker/qubership-monitoring-operator/api"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/discovery"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type AlertManagerReconciler struct {
	*utils.ComponentReconciler
}

func NewAlertManagerReconciler(c client.Client, s *runtime.Scheme, dc discovery.DiscoveryInterface) *AlertManagerReconciler {
	return &AlertManagerReconciler{
		ComponentReconciler: &utils.ComponentReconciler{
			Client: c,
			Scheme: s,
			Dc:     dc,
			Log:    utils.Logger("alertmanager_reconciler"),
		},
	}
}

// Run reconciles alertmanager custom resource.
// Creates new deployment, service, service account, cluster role and cluster role binding if its don't exists.
// Updates deployment and service in case of any changes.
// Returns true if need to requeue, false otherwise.
func (r *AlertManagerReconciler) Run(cr *v1beta1.PlatformMonitoring) error {
	r.Log.Info("Reconciling component")

	if cr.Spec.AlertManager != nil && cr.Spec.AlertManager.IsInstall() {
		if !cr.Spec.AlertManager.Paused {
			if err := r.handleSecret(cr); err != nil {
				return err
			}
			if err := r.handleServiceAccount(cr); err != nil {
				return err
			}
			if err := r.handleService(cr); err != nil {
				return err
			}
			if err := r.handleAlertmanager(cr); err != nil {
				return err
			}

			// Reconcile Ingress (version v1beta1) if necessary and the cluster is has such API
			// This API unavailable in k8s v1.22+
			if r.HasIngressV1beta1Api() {
				if cr.Spec.AlertManager.Ingress != nil && cr.Spec.AlertManager.Ingress.IsInstall() {
					if err := r.handleIngressV1beta1(cr); err != nil {
						return err
					}
				} else {
					if err := r.deleteIngressV1beta1(cr); err != nil {
						r.Log.Error(err, "Can not delete Ingress")
					}
				}
			}
			// Reconcile Ingress (version v1) if necessary and the cluster is has such API
			// This API available in k8s v1.19+
			if r.HasIngressV1Api() {
				if cr.Spec.AlertManager.Ingress != nil && cr.Spec.AlertManager.Ingress.IsInstall() {
					if err := r.handleIngressV1(cr); err != nil {
						return err
					}
				} else {
					if err := r.deleteIngressV1(cr); err != nil {
						r.Log.Error(err, "Can not delete Ingress")
					}
				}
			}

			if cr.Spec.AlertManager.PodMonitor != nil && cr.Spec.AlertManager.PodMonitor.IsInstall() {
				if err := r.handlePodMonitor(cr); err != nil {
					return err
				}
			} else {
				if err := r.deletePodMonitor(cr); err != nil {
					r.Log.Error(err, "Can not delete PodMonitor")
				}
			}
			r.Log.Info("Component reconciled")
		} else {
			r.Log.Info("Reconciling paused")
			r.Log.Info("Component NOT reconciled")
		}
	} else {
		r.Log.Info("Uninstalling component if exists")
		r.uninstall(cr)
		r.Log.Info("Component reconciled")
	}
	return nil
}

// uninstall deletes all resources related to the component
func (r *AlertManagerReconciler) uninstall(cr *v1beta1.PlatformMonitoring) {
	if err := r.deleteServiceAccount(cr); err != nil {
		r.Log.Error(err, "Can not delete ServiceAccount")
	}
	if err := r.deleteService(cr); err != nil {
		r.Log.Error(err, "Can not delete Service")
	}
	if err := r.deleteSecret(cr); err != nil {
		r.Log.Error(err, "Can not delete Secret")
	}
	if err := r.deleteAlertmanager(cr); err != nil {
		r.Log.Error(err, "Can not delete AlertManager")
	}
	if err := r.deletePodMonitor(cr); err != nil {
		r.Log.Error(err, "Can not delete PodMonitor")
	}
	// Try to delete Ingress (version v1beta1) is there is such API
	// This API unavailable in k8s v1.22+
	if r.HasIngressV1beta1Api() {
		if err := r.deleteIngressV1beta1(cr); err != nil {
			r.Log.Error(err, "Can not delete Ingress")
		}
	}
	// Try to delete Ingress (version v1) is there is such API
	// This API available in k8s v1.19+
	if r.HasIngressV1Api() {
		if err := r.deleteIngressV1(cr); err != nil {
			r.Log.Error(err, "Can not delete Ingress")
		}
	}

}
