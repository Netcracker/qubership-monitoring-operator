package grafana

import (
	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/gateway"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)


// isSecretUpdated is set to true when the admin credentials secret changes
// between reconcile cycles and must be reset via grafana cli.
var isSecretUpdated = false

// currentAdminSecretChecksum holds the SHA256 of the admin credentials secret
// computed during the last reconcile. Written into the Grafana CR pod-template
// annotation so grafana-operator triggers a rolling restart on secret change.
var currentAdminSecretChecksum string

// adminSecretChecksumAnnotation is the pod-template annotation key used to
// propagate the admin secret checksum and drive rolling restarts.
const adminSecretChecksumAnnotation = "checksum/admin-secret"

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
			if err := r.handleGrafanaCredentialsSecret(cr); err != nil {
				return err
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
			ingressHost := ""
			if cr.Spec.Grafana.Ingress != nil {
				ingressHost = cr.Spec.Grafana.Ingress.Host
			}
			parentRefs := []monv1.GatewayParentRef(nil)
			if cr.Spec.GatewayAPI != nil {
				parentRefs = cr.Spec.GatewayAPI.ParentRefs
			}
			if err := gateway.ReconcileGatewayRoutes(r.ComponentReconciler, cr, gateway.GatewayRouteConfig{
				NamePrefix:     cr.GetNamespace() + "-" + utils.GrafanaComponentName,
				Namespace:      cr.GetNamespace(),
				Host:           ingressHost,
				ServiceName:    utils.GrafanaServiceName,
				ServicePort:    int32(utils.GrafanaServicePort),
				Labels:         map[string]string{"name": utils.TruncLabel(cr.GetNamespace() + "-" + utils.GrafanaComponentName), "app.kubernetes.io/name": utils.TruncLabel(cr.GetNamespace() + "-" + utils.GrafanaComponentName), "app.kubernetes.io/instance": utils.GetInstanceLabel(cr.GetNamespace()+"-"+utils.GrafanaComponentName, cr.GetNamespace()), "app.kubernetes.io/version": utils.GetTagFromImage(cr.Spec.Grafana.Image)},
				ParentRefs:     parentRefs,
				ComponentRoute: cr.Spec.Grafana.HTTPRoute,
			}); err != nil {
				return err
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
			// Reset Grafana admin credentials in the running database when the
			// secret has changed. This handles both PV and non-PV deployments:
			// the rolling restart triggered by the checksum annotation handles
			// the emptyDir case; the grafana cli exec handles the PV case.
			if isSecretUpdated {
				if err := r.resetGrafanaCredentials(cr); err != nil {
					r.Log.Error(err, "Cannot reset Grafana credentials")
					return err
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
	// Try to delete Ingress (version v1) is there is such API
	// This API available in k8s v1.19+
	if r.HasIngressV1Api() {
		if err := r.deleteIngressV1(cr); err != nil {
			r.Log.Error(err, "Can not delete Ingress")
		}
	}
	parentRefs := []monv1.GatewayParentRef(nil)
	if cr.Spec.GatewayAPI != nil {
		parentRefs = cr.Spec.GatewayAPI.ParentRefs
	}
	var componentRoute *monv1.GatewayHTTPRoute
	if cr.Spec.Grafana != nil {
		componentRoute = cr.Spec.Grafana.HTTPRoute
	}
	if err := gateway.DeleteGatewayRoutes(r.ComponentReconciler, gateway.GatewayRouteConfig{
		NamePrefix:     cr.GetNamespace() + "-" + utils.GrafanaComponentName,
		Namespace:      cr.GetNamespace(),
		ServiceName:    utils.GrafanaServiceName,
		ServicePort:    int32(utils.GrafanaServicePort),
		ParentRefs:     parentRefs,
		ComponentRoute: componentRoute,
	}); err != nil {
		r.Log.Error(err, "Can not delete Gateway API routes.")
	}
}
