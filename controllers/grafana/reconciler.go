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

type GrafanaReconciler struct {
	KubeClient kubernetes.Interface
	config     *rest.Config
	// newExecutor overrides remotecommand.NewSPDYExecutor in tests.
	newExecutor executorFactory
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
// Returns an error wrapping ErrAdminPasswordSyncPending or ErrDatasourceMigrationPending when
// Grafana is not ready yet; the caller requeues the reconcile.
func (r *GrafanaReconciler) Run(cr *monv1.PlatformMonitoring) error {
	r.Log.Info("Reconciling component")

	if cr.Spec.Grafana == nil || !cr.Spec.Grafana.IsInstall() {
		r.Log.Info("Uninstalling component if exists")
		r.uninstall(cr)
		r.Log.Info("Component reconciled")
		return nil
	}
	if cr.Spec.Grafana.Paused {
		r.Log.Info("Reconciling paused")
		r.Log.Info("Component NOT reconciled")
		return nil
	}

	// Reconcile resources with creation and update
	if err := r.handleGrafana(cr); err != nil {
		return err
	}
	// Datasource reconciliation authenticates with the admin password, so it waits for the sync.
	// The routing resources do not depend on it and are reconciled while the sync is pending.
	syncErr := r.syncGrafanaAdminPassword(cr)
	if syncErr != nil && !IsPending(syncErr) {
		r.Log.Error(syncErr, "Cannot sync Grafana admin password")
		return syncErr
	}
	if syncErr == nil {
		if err := r.reconcileDataSources(cr); err != nil {
			return err
		}
	}
	if err := r.reconcileRouting(cr); err != nil {
		return err
	}
	if syncErr != nil {
		r.Log.Info("Grafana admin password sync is pending; datasources are reconciled later", "reason", syncErr.Error())
		return syncErr
	}
	r.Log.Info("Component reconciled")
	return nil
}

// reconcileDataSources reconciles the Prometheus, Jaeger, ClickHouse, and Promxy datasources.
func (r *GrafanaReconciler) reconcileDataSources(cr *monv1.PlatformMonitoring) error {
	if err := r.handleGrafanaDataSource(cr); err != nil {
		return err
	}
	if err := r.handleJaegerDataSources(cr); err != nil {
		return err
	}
	if err := r.handleClickHouseDataSources(cr); err != nil {
		return err
	}
	// Reconcile Promxy datasource only when Promxy is installed (otherwise Grafana hangs on missing service)
	if cr.Spec.Promxy != nil && cr.Spec.Promxy.IsInstall() {
		return r.handleGrafanaPromxyDataSource(cr)
	}
	if err := r.deleteGrafanaPromxyDataSource(cr); err != nil {
		r.Log.Error(err, "Can not delete GrafanaPromxyDataSource")
	}
	return nil
}

// reconcileRouting reconciles the Ingress, the Gateway API routes, and the PodMonitor.
func (r *GrafanaReconciler) reconcileRouting(cr *monv1.PlatformMonitoring) error {
	// Reconcile Ingress (version v1) if necessary and the cluster is has such API
	// This API available in k8s v1.19+
	if r.HasIngressV1Api() {
		if cr.Spec.Grafana.Ingress != nil && cr.Spec.Grafana.Ingress.IsInstall() {
			if err := r.handleIngressV1(cr); err != nil {
				return err
			}
		} else if err := r.deleteIngressV1(cr); err != nil {
			r.Log.Error(err, "Can not delete Ingress")
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
		ServiceName:    grafanaServiceName(cr),
		ServicePort:    int32(utils.GrafanaServicePort),
		Labels:         map[string]string{"name": utils.TruncLabel(cr.GetNamespace() + "-" + utils.GrafanaComponentName), "app.kubernetes.io/name": utils.TruncLabel(cr.GetNamespace() + "-" + utils.GrafanaComponentName), "app.kubernetes.io/instance": utils.GetInstanceLabel(cr.GetNamespace()+"-"+utils.GrafanaComponentName, cr.GetNamespace()), "app.kubernetes.io/version": utils.GetTagFromImage(cr.Spec.Grafana.Image)},
		ParentRefs:     parentRefs,
		ComponentRoute: cr.Spec.Grafana.HTTPRoute,
	}); err != nil {
		return err
	}
	// Reconcile Pod Monitor
	if cr.Spec.Grafana.PodMonitor != nil && cr.Spec.Grafana.PodMonitor.IsInstall() {
		return r.handlePodMonitor(cr)
	}
	if err := r.deletePodMonitor(cr); err != nil {
		r.Log.Error(err, "Can not delete PodMonitor")
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
	if err := r.deleteDiscoveredDataSources(cr, jaegerDatasourceComponent); err != nil {
		r.Log.Error(err, "Can not delete Jaeger GrafanaDatasource")
	}
	if err := r.deleteDiscoveredDataSources(cr, clickHouseDatasourceComponent); err != nil {
		r.Log.Error(err, "Can not delete ClickHouse GrafanaDatasource")
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
		ServiceName:    grafanaServiceName(cr),
		ServicePort:    int32(utils.GrafanaServicePort),
		ParentRefs:     parentRefs,
		ComponentRoute: componentRoute,
	}); err != nil {
		r.Log.Error(err, "Can not delete Gateway API routes.")
	}
}
