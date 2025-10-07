package pushgateway

import (
	v1alpha1 "github.com/Netcracker/qubership-monitoring-operator/api/v1alpha1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/discovery"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PushgatewayReconciler struct {
	*utils.ComponentReconciler
}

func NewPushgatewayReconciler(c client.Client, s *runtime.Scheme, dc discovery.DiscoveryInterface) *PushgatewayReconciler {
	return &PushgatewayReconciler{
		ComponentReconciler: &utils.ComponentReconciler{
			Client: c,
			Scheme: s,
			Dc:     dc,
			Log:    utils.Logger("pushgateway_reconciler"),
		},
	}
}

// Run reconciliation for pushgateway configuration.
// Creates new deployment, service and service monitor if its don't exists.
// Updates deployment, service and service monitor in case of any changes.
// Returns true if need to requeue, false otherwise.
func (r *PushgatewayReconciler) Run(cr *v1alpha1.PlatformMonitoring) error {
	r.Log.Info("Reconciling component")

	if cr.Spec.Pushgateway != nil && cr.Spec.Pushgateway.IsInstall() {
		if !cr.Spec.Pushgateway.Paused {

			if err := r.handleService(cr); err != nil {
				return err
			}
			if cr.Spec.Pushgateway.Storage != nil {
				if err := r.handlePVC(cr); err != nil {
					return err
				}
			}
			if err := r.handleDeployment(cr); err != nil {
				return err
			}
			// Reconcile Ingress (version v1) if necessary and the cluster is has such API
			// This API available in k8s v1.19+
			if r.HasIngressV1Api() {
				if cr.Spec.Pushgateway.Ingress != nil && cr.Spec.Pushgateway.Ingress.IsInstall() {
					if err := r.handleIngressV1(cr); err != nil {
						return err
					}
				} else {
					if err := r.deleteIngressV1(cr); err != nil {
						r.Log.Error(err, "Can not delete Ingress")
					}
				}
			}
			// Reconcile ServiceMonitor if necessary
			if cr.Spec.Pushgateway.ServiceMonitor != nil && cr.Spec.Pushgateway.ServiceMonitor.IsInstall() {
				if err := r.handleServiceMonitor(cr); err != nil {
					return err
				}
			} else {
				if err := r.deleteServiceMonitor(cr); err != nil {
					r.Log.Error(err, "Can not delete ServiceMonitor")
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
func (r *PushgatewayReconciler) uninstall(cr *v1alpha1.PlatformMonitoring) {
	if err := r.deleteDeployment(cr); err != nil {
		r.Log.Error(err, "Can not delete Deployment")
	}
	if err := r.deleteService(cr); err != nil {
		r.Log.Error(err, "Can not delete Service")
	}
	// Try to delete Ingress (version v1) is there is such API
	// This API available in k8s v1.19+
	if r.HasIngressV1Api() {
		if err := r.deleteIngressV1(cr); err != nil {
			r.Log.Error(err, "Can not delete Ingress")
		}
	}
}
