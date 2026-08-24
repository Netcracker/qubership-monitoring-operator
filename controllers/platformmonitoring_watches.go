package controllers

import (
	"context"
	"os"
	"strconv"
	"strings"
	"time"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	vmetricsv1b1 "github.com/VictoriaMetrics/operator/api/operator/v1beta1"
	grafv1 "github.com/grafana/grafana-operator/v5/api/v1beta1"
	secv1 "github.com/openshift/api/security/v1"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	watchBasedReconcileEnv = "WATCH_BASED_RECONCILE"
	gatewayAPIGroup        = "gateway.networking.k8s.io"
	httpRouteKind          = "HTTPRoute"
	httpRouteVersion       = "v1"

	nodeExporterSCCOwnerNameAnnotation      = "monitoring.netcracker.com/platform-monitoring-name"
	nodeExporterSCCOwnerNamespaceAnnotation = "monitoring.netcracker.com/platform-monitoring-namespace"
)

// watchBasedReconcileEnabled is on by default. Set WATCH_BASED_RECONCILE=false|0
// to fall back to PlatformMonitoring For + periodic repair only.
func watchBasedReconcileEnabled() bool {
	v := strings.TrimSpace(os.Getenv(watchBasedReconcileEnv))
	return v != "0" && !strings.EqualFold(v, "false")
}

// jaegerServiceWatchAllowed is the privileged-mode gate for the mapped
// Service watch. Namespaced Helm still sets WATCH_NAMESPACE, so even when
// this is true the informer only sees Services in that namespace unless the
// manager cache is cluster-scoped.
func jaegerServiceWatchAllowed() bool {
	return utils.PrivilegedRights
}

// clusterScopedWatchAllowed gates ClusterRole / ClusterRoleBinding / SCC
// informers. Those objects are only created in privileged mode, and the
// namespaced Role lacks cluster-scoped watch verbs.
func clusterScopedWatchAllowed() bool {
	return utils.PrivilegedRights
}

// successfulReconcileResult is the per-PlatformMonitoring periodic repair
// requeue. Interval 0 returns an empty result (no hot loop). A positive value
// is a repair timer, not controller-runtime informer SyncPeriod.
func successfulReconcileResult() (ctrl.Result, error) {
	raw := utils.GetEnvWithDefaultValue("RECONCILIATION_INTERVAL")
	interval, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return ctrl.Result{}, err
	}
	if interval <= 0 {
		return ctrl.Result{}, nil
	}
	return ctrl.Result{RequeueAfter: time.Duration(interval) * time.Second}, nil
}

func (r *PlatformMonitoringReconciler) hasKind(groupVersion, kind string) bool {
	if r.DiscoveryClient == nil {
		return false
	}
	ok, err := utils.ResourceExists(r.DiscoveryClient, groupVersion, kind)
	return err == nil && ok
}

func (r *PlatformMonitoringReconciler) routeAPIServed() bool {
	return r.hasKind("route.openshift.io/v1", "Route")
}

func (r *PlatformMonitoringReconciler) own(b *builder.Builder, obj client.Object) *builder.Builder {
	return b.Owns(obj, builder.WithPredicates(ownedChildPredicate()))
}

func (r *PlatformMonitoringReconciler) ownIfServed(b *builder.Builder, groupVersion, kind string, obj client.Object) *builder.Builder {
	if !r.hasKind(groupVersion, kind) {
		r.Log.Info("Skipping Owns: GVK not served; install CRDs then restart the manager",
			"groupVersion", groupVersion, "kind", kind)
		return b
	}
	return r.own(b, obj)
}

func (r *PlatformMonitoringReconciler) watchClusterScopedIfAllowed(b *builder.Builder, obj client.Object, mapFn handler.MapFunc) *builder.Builder {
	if !clusterScopedWatchAllowed() {
		r.Log.Info("Skipping cluster-scoped watch: PrivilegedRights=false; drift stays on the periodic repair requeue",
			"kind", obj.GetObjectKind().GroupVersionKind().Kind)
		return b
	}
	return b.Watches(obj, handler.EnqueueRequestsFromMapFunc(mapFn), builder.WithPredicates(ownedChildPredicate()))
}

func (r *PlatformMonitoringReconciler) ownHTTPRouteIfServed(b *builder.Builder) *builder.Builder {
	gv := schema.GroupVersion{Group: gatewayAPIGroup, Version: httpRouteVersion}
	if !r.hasKind(gv.String(), httpRouteKind) {
		r.Log.Info("Skipping Owns: GVK not served; install CRDs then restart the manager",
			"groupVersion", gv.String(), "kind", httpRouteKind)
		return b
	}
	if r.Scheme != nil && !r.Scheme.Recognizes(gv.WithKind(httpRouteKind)) {
		r.Scheme.AddKnownTypeWithName(gv.WithKind(httpRouteKind), &unstructured.Unstructured{})
		r.Scheme.AddKnownTypeWithName(gv.WithKind(httpRouteKind+"List"), &unstructured.UnstructuredList{})
		metav1.AddToGroupVersion(r.Scheme, gv)
	}
	u := &unstructured.Unstructured{}
	u.SetGroupVersionKind(gv.WithKind(httpRouteKind))
	return r.own(b, u)
}

func (r *PlatformMonitoringReconciler) addWatchBasedSources(b *builder.Builder) *builder.Builder {
	b = r.own(b, &appsv1.Deployment{})
	b = r.own(b, &appsv1.DaemonSet{})
	b = r.own(b, &corev1.ServiceAccount{})
	b = r.own(b, &corev1.Service{})
	b = r.own(b, &corev1.Endpoints{})
	b = r.own(b, &corev1.Secret{})
	b = r.own(b, &corev1.PersistentVolumeClaim{})
	b = r.own(b, &rbacv1.Role{})
	b = r.own(b, &rbacv1.RoleBinding{})

	grafGV := grafv1.SchemeGroupVersion.String()
	b = r.ownIfServed(b, grafGV, "Grafana", &grafv1.Grafana{})
	b = r.ownIfServed(b, grafGV, "GrafanaDashboard", &grafv1.GrafanaDashboard{})
	b = r.ownIfServed(b, grafGV, "GrafanaDatasource", &grafv1.GrafanaDatasource{})

	promGV := promv1.SchemeGroupVersion.String()
	b = r.ownIfServed(b, promGV, "Prometheus", &promv1.Prometheus{})
	b = r.ownIfServed(b, promGV, "Alertmanager", &promv1.Alertmanager{})
	b = r.ownIfServed(b, promGV, "PrometheusRule", &promv1.PrometheusRule{})
	b = r.ownIfServed(b, promGV, "ServiceMonitor", &promv1.ServiceMonitor{})
	b = r.ownIfServed(b, promGV, "PodMonitor", &promv1.PodMonitor{})

	vmGV := vmetricsv1b1.SchemeGroupVersion.String()
	b = r.ownIfServed(b, vmGV, "VMSingle", &vmetricsv1b1.VMSingle{})
	b = r.ownIfServed(b, vmGV, "VMAgent", &vmetricsv1b1.VMAgent{})
	b = r.ownIfServed(b, vmGV, "VMAuth", &vmetricsv1b1.VMAuth{})
	b = r.ownIfServed(b, vmGV, "VMAlert", &vmetricsv1b1.VMAlert{})
	b = r.ownIfServed(b, vmGV, "VMAlertmanager", &vmetricsv1b1.VMAlertmanager{})
	b = r.ownIfServed(b, vmGV, "VMCluster", &vmetricsv1b1.VMCluster{})
	b = r.ownIfServed(b, vmGV, "VMUser", &vmetricsv1b1.VMUser{})

	b = r.ownIfServed(b, networkingv1.SchemeGroupVersion.String(), "Ingress", &networkingv1.Ingress{})
	b = r.ownHTTPRouteIfServed(b)

	b = r.watchClusterScopedIfAllowed(b, &rbacv1.ClusterRole{
		TypeMeta: metav1.TypeMeta{Kind: "ClusterRole", APIVersion: rbacv1.SchemeGroupVersion.String()},
	}, r.mapClusterScopedToPlatformMonitoring)
	b = r.watchClusterScopedIfAllowed(b, &rbacv1.ClusterRoleBinding{
		TypeMeta: metav1.TypeMeta{Kind: "ClusterRoleBinding", APIVersion: rbacv1.SchemeGroupVersion.String()},
	}, r.mapClusterScopedToPlatformMonitoring)
	if r.hasKind(secv1.GroupVersion.String(), "SecurityContextConstraints") {
		b = r.watchClusterScopedIfAllowed(b, &secv1.SecurityContextConstraints{}, r.mapSecurityContextConstraintsToPlatformMonitorings)
	} else {
		r.Log.Info("Skipping SCC watch: GVK not served")
	}

	if jaegerServiceWatchAllowed() {
		b = b.Watches(
			&corev1.Service{},
			handler.EnqueueRequestsFromMapFunc(r.mapJaegerServicesToPlatformMonitorings),
			builder.WithPredicates(jaegerServicePredicate()),
		)
	} else {
		r.Log.Info("Skipping Jaeger Service watch: PrivilegedRights=false; cross-namespace discovery stays on the periodic repair requeue")
	}
	return b
}

func (r *PlatformMonitoringReconciler) mapClusterScopedToPlatformMonitoring(ctx context.Context, obj client.Object) []reconcile.Request {
	if obj == nil {
		return nil
	}
	ns := obj.GetLabels()[utils.InstallationNamespaceLabelKey]
	if ns == "" {
		return nil
	}
	return r.requestsForPlatformMonitoringsInNamespace(ctx, ns)
}

func (r *PlatformMonitoringReconciler) mapSecurityContextConstraintsToPlatformMonitorings(ctx context.Context, obj client.Object) []reconcile.Request {
	if obj == nil {
		return nil
	}
	if reqs := r.mapClusterScopedToPlatformMonitoring(ctx, obj); len(reqs) > 0 {
		return reqs
	}
	annotations := obj.GetAnnotations()
	name := annotations[nodeExporterSCCOwnerNameAnnotation]
	namespace := annotations[nodeExporterSCCOwnerNamespaceAnnotation]
	if name != "" && namespace != "" {
		return []reconcile.Request{{NamespacedName: client.ObjectKey{Name: name, Namespace: namespace}}}
	}
	if obj.GetName() == utils.VmOperatorComponentName {
		return r.listAllPlatformMonitoringRequests(ctx)
	}
	return nil
}

func (r *PlatformMonitoringReconciler) requestsForPlatformMonitoringsInNamespace(ctx context.Context, namespace string) []reconcile.Request {
	list := &monv1.PlatformMonitoringList{}
	if err := r.List(ctx, list); err != nil {
		r.Log.Error(err, "Failed listing PlatformMonitoring for cluster-scoped mapping")
		return nil
	}
	var reqs []reconcile.Request
	for i := range list.Items {
		pm := &list.Items[i]
		if pm.GetNamespace() == namespace {
			reqs = append(reqs, reconcile.Request{NamespacedName: client.ObjectKeyFromObject(pm)})
		}
	}
	return reqs
}

func (r *PlatformMonitoringReconciler) listAllPlatformMonitoringRequests(ctx context.Context) []reconcile.Request {
	list := &monv1.PlatformMonitoringList{}
	if err := r.List(ctx, list); err != nil {
		r.Log.Error(err, "Failed listing PlatformMonitoring for SCC mapping")
		return nil
	}
	reqs := make([]reconcile.Request, 0, len(list.Items))
	for i := range list.Items {
		reqs = append(reqs, reconcile.Request{NamespacedName: client.ObjectKeyFromObject(&list.Items[i])})
	}
	return reqs
}

func (r *PlatformMonitoringReconciler) mapJaegerServicesToPlatformMonitorings(ctx context.Context, obj client.Object) []reconcile.Request {
	svc, ok := obj.(*corev1.Service)
	if !ok || svc == nil {
		return nil
	}
	if !labels.SelectorFromSet(utils.JaegerServiceLabels).Matches(labels.Set(svc.GetLabels())) {
		return nil
	}

	list := &monv1.PlatformMonitoringList{}
	if err := r.List(ctx, list); err != nil {
		r.Log.Error(err, "Failed listing PlatformMonitoring for Jaeger Service mapping")
		return nil
	}

	var reqs []reconcile.Request
	for i := range list.Items {
		pm := &list.Items[i]
		if !jaegerDatasourceRequested(pm) {
			continue
		}
		reqs = append(reqs, reconcile.Request{NamespacedName: client.ObjectKeyFromObject(pm)})
	}
	return reqs
}
