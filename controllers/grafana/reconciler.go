package grafana

import (
	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var isSecretUpdated = false

// isManageAdminSecret returns true when monitoring-operator manages grafana-admin-credentials
// (DisableDefaultAdminSecret is true).
func isManageAdminSecret(cr *monv1.PlatformMonitoring) bool {
	if cr.Spec.Grafana == nil {
		return false
	}
	return cr.Spec.Grafana.DisableDefaultAdminSecret != nil && *cr.Spec.Grafana.DisableDefaultAdminSecret
}

type GrafanaReconciler struct {
	KubeClient kubernetes.Interface
	config     *rest.Config
	*utils.ComponentReconciler
}

func NewGrafanaReconciler(c client.Client, s *runtime.Scheme, dc discovery.DiscoveryInterface, r *rest.Config) *GrafanaReconciler {
	cl, _ := kubernetes.NewForConfig(r)
	return &GrafanaReconciler{
		ComponentReconciler: &utils.ComponentReconciler{
			Client: c,
			Scheme: s,
			Dc:     dc,
			Log:    utils.Logger("grafana_reconciler"),
		},
		KubeClient: cl,
		config:     r,
	}
}

// Run reconciles grafana custom resource.
// Creates new custom resources: Grafana and GrafanaDataSource if its don't exists.
// Updates custom resources in case of any changes.
// Returns true if need to requeue, false otherwise.
func (r *GrafanaReconciler) Run(cr *monv1.PlatformMonitoring) error {
	r.Log.Info("Reconciling component")

	if cr.Spec.Grafana != nil && cr.Spec.Grafana.IsInstall() {
		if !cr.Spec.Grafana.Paused {
			if isManageAdminSecret(cr) {
				if err := r.handleGrafanaCredentialsSecret(cr); err != nil {
					return err
				}
			}
			// Reconcile resources with creation and update
			if err := r.handleGrafana(cr); err != nil {
				return err
			}
			if err := r.handleGrafanaDataSource(cr); err != nil {
				return err
			}
			// Reconcile Promxy datasource only when Promxy is installed (otherwise Grafana hangs on missing service)
			if cr.Spec.Promxy != nil && cr.Spec.Promxy.IsInstall() {
				if err := r.handleGrafanaPromxyDataSource(cr); err != nil {
					return err
				}
			} else {
				if err := r.deleteGrafanaPromxyDataSource(cr); err != nil {
					r.Log.Error(err, "Can not delete GrafanaPromxyDataSource")
				}
			}

			// Reconcile Ingress (version v1beta1) if necessary and the cluster is has such API
			// This API unavailable in k8s v1.22+
			if r.HasIngressV1beta1Api() {
				if cr.Spec.Grafana.Ingress != nil && cr.Spec.Grafana.Ingress.IsInstall() {
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
				if cr.Spec.Grafana.Ingress != nil && cr.Spec.Grafana.Ingress.IsInstall() {
					if err := r.handleIngressV1(cr); err != nil {
						return err
					}
				} else {
					if err := r.deleteIngressV1(cr); err != nil {
						r.Log.Error(err, "Can not delete Ingress")
					}
				}
			}
			// Reconcile Pod Monitor
			if cr.Spec.Grafana.PodMonitor != nil && cr.Spec.Grafana.PodMonitor.IsInstall() {
				if err := r.handlePodMonitor(cr); err != nil {
					return err
				}
			} else {
				if err := r.deletePodMonitor(cr); err != nil {
					r.Log.Error(err, "Can not delete PodMonitor")
				}
			}
			// resetGrafanaCredentials is not used when disableDefaultAdminSecret=true:
			// we do not create/update the secret; user manages it and restarts the Grafana pod if needed.
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
func (r *GrafanaReconciler) uninstall(cr *monv1.PlatformMonitoring) {
	if err := r.deleteGrafana(cr); err != nil {
		r.Log.Error(err, "Can not delete Grafana")
	}
	if err := r.deleteGrafanaDataSource(cr); err != nil {
		r.Log.Error(err, "Can not delete GrafanaDataSource")
	}
	if err := r.deleteGrafanaPromxyDataSource(cr); err != nil {
		r.Log.Error(err, "Can not delete GrafanaPromxyDataSource")
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
