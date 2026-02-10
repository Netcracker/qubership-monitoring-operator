package prometheus

import (
	"context"
	"fmt"

	monv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/discovery"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// PrometheusReconciler provides methods to reconcile Prometheus
type PrometheusReconciler struct {
	*utils.ComponentReconciler
}

// NewPrometheusReconciler creates an instance of PrometheusReconciler
func NewPrometheusReconciler(c client.Client, s *runtime.Scheme, dc discovery.DiscoveryInterface) *PrometheusReconciler {
	return &PrometheusReconciler{
		ComponentReconciler: &utils.ComponentReconciler{
			Client: c,
			Scheme: s,
			Dc:     dc,
			Log:    utils.Logger("prometheus_reconciler"),
		},
	}
}

// Run reconciles prometheus-operator.
// Creates new deployment, service, service account, cluster role and cluster role binding if its don't exists.
// Updates deployment and service in case of any changes.
// Returns true if need to requeue, false otherwise.
func (r *PrometheusReconciler) Run(cr *monv1.PlatformMonitoring) error {
	r.Log.Info("Reconciling component")

	if err := r.removePrometheusPVC(cr); err != nil {
		return err
	}

	if cr.Spec.Prometheus != nil && cr.Spec.Prometheus.IsInstall() {
		if !cr.Spec.Prometheus.Paused {

			// Reconcile resources with creation
			if err := r.handleServiceAccount(cr); err != nil {
				return err
			}

			// Reconcile ClusterRole and ClusterRoleBinding only if privileged mode used
			if utils.PrivilegedRights {
				if err := r.handleClusterRole(cr); err != nil {
					return err
				}
				if err := r.handleClusterRoleBinding(cr); err != nil {
					return err
				}
			} else {
				r.Log.Info("Skip ClusterRole and ClusterRoleBinding resources reconciliation because privilegedRights=false")
			}
			// Reconcile Prometheus with creation and update
			if err := r.handlePrometheus(cr); err != nil {
				return err
			}

			// Reconcile Ingress (version v1beta1) if necessary and the cluster is has such API
			// This API unavailable in k8s v1.22+
			if r.HasIngressV1beta1Api() {
				if cr.Spec.Prometheus.Ingress != nil && cr.Spec.Prometheus.Ingress.IsInstall() {
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
				if cr.Spec.Prometheus.Ingress != nil && cr.Spec.Prometheus.Ingress.IsInstall() {
					if err := r.handleIngressV1(cr); err != nil {
						return err
					}
				} else {
					if err := r.deleteIngressV1(cr); err != nil {
						r.Log.Error(err, "Can not delete Ingress")
					}
				}
			}
			// Reconcile PodMonitor if necessary
			if cr.Spec.Prometheus.PodMonitor != nil && cr.Spec.Prometheus.PodMonitor.IsInstall() {
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

// removePrometheusPVC deletes PVC for Prometheus pod if it isn't needed anymore.
func (r *PrometheusReconciler) removePrometheusPVC(cr *monv1.PlatformMonitoring) error {
	if cr.Spec.Prometheus == nil || !cr.Spec.Prometheus.IsInstall() || cr.Spec.Prometheus.Storage == nil {
		// We can use hardcoded value for now because prometheus-operator < v0.40.0
		// doesn't allow setting metadata for PVC template.
		// See https://github.com/prometheus-operator/prometheus-operator/issues/3093
		// ToDo: Implement deleting PVC by special value when prometheus-operator will be upgraded to v0.40.0+
		pvcName := "prometheus-k8s-db-prometheus-k8s-0"

		// Find PVC resource
		pvc := &corev1.PersistentVolumeClaim{}
		err := r.Client.Get(
			context.TODO(),
			types.NamespacedName{Name: pvcName, Namespace: cr.GetNamespace()},
			pvc,
		)
		if err != nil {
			if errors.IsNotFound(err) {
				return nil
			}

			return err
		}
		// Delete PVC
		r.Log.Info(fmt.Sprintf("delete Prometheus PVC with name %q as it is not needed anymore", pvc.GetName()))
		err = r.Client.Delete(context.TODO(), pvc)
		if err != nil {
			r.Log.Error(err, fmt.Sprintf("can not delete PVC %q", pvc.GetName()))
			return err
		}
	}
	return nil
}

// uninstall deletes all resources related to the component
func (r *PrometheusReconciler) uninstall(cr *monv1.PlatformMonitoring) {
	if utils.PrivilegedRights {
		if err := r.deleteClusterRole(cr); err != nil {
			r.Log.Error(err, "Can not delete ClusterRole")
		}
		if err := r.deleteClusterRoleBinding(cr); err != nil {
			r.Log.Error(err, "Can not delete ClusterRoleBinding")
		}
	}
	if err := r.deleteServiceAccount(cr); err != nil {
		r.Log.Error(err, "Can not delete ServiceAccount")
	}
	if err := r.deletePrometheus(cr); err != nil {
		r.Log.Error(err, "Can not delete Prometheus")
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
