package controllers

import (
	"context"
	"os"
	"strconv"
	"strings"
	"time"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	grafv1 "github.com/grafana/grafana-operator/v5/api/v1beta1"
	routev1 "github.com/openshift/api/route/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const watchBasedReconcileEnv = "WATCH_BASED_RECONCILE"

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

// routeAPIServed reports whether a typed Route watch would be safe to register.
// Route is not owned in this slice; drift waits on the periodic repair requeue
// until an OpenShift smoke justifies a discovery-gated watch.
func (r *PlatformMonitoringReconciler) routeAPIServed() bool {
	return r.hasKind(routev1.GroupVersion.String(), "Route")
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

func (r *PlatformMonitoringReconciler) addWatchBasedSources(b *builder.Builder) *builder.Builder {
	// Kind-proven same-namespace Owns only. Wider child types (Service, Role,
	// VM CRs, ServiceMonitors, …) still wait on the repair requeue until
	// skip-if-equal exists; each Owns event runs the full handler loop.
	// ClusterRole / ClusterRoleBinding / Node / Route are never Owned.
	b = r.own(b, &appsv1.Deployment{})

	grafGV := grafv1.SchemeGroupVersion.String()
	b = r.ownIfServed(b, grafGV, "Grafana", &grafv1.Grafana{})
	b = r.ownIfServed(b, grafGV, "GrafanaDashboard", &grafv1.GrafanaDashboard{})

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
