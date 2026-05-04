/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"strconv"
	"time"

	qubershiporgv1 "github.com/Netcracker/qubership-monitoring-operator/api/v1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/alertmanager"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/grafana"
	grafanaoperator "github.com/Netcracker/qubership-monitoring-operator/controllers/grafana-operator"
	kubernetesmonitors "github.com/Netcracker/qubership-monitoring-operator/controllers/kubernetes-monitors"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/kubestatemetrics"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/nodeexporter"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/prometheus"
	prometheusoperator "github.com/Netcracker/qubership-monitoring-operator/controllers/prometheus-operator"
	prometheusrules "github.com/Netcracker/qubership-monitoring-operator/controllers/prometheus-rules"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/pushgateway"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/victoriametrics/vmagent"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/victoriametrics/vmalert"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/victoriametrics/vmalertmanager"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/victoriametrics/vmauth"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/victoriametrics/vmcluster"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/victoriametrics/vmoperator"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/victoriametrics/vmsingle"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/victoriametrics/vmuser"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// PlatformMonitoringReconciler reconciles a PlatformMonitoring object
type PlatformMonitoringReconciler struct {
	Log logr.Logger
	client.Client
	Scheme *runtime.Scheme
	// Cluster Config
	Config *rest.Config
	// Client to discovery cluster API
	DiscoveryClient discovery.DiscoveryInterface
}

// +kubebuilder:rbac:groups=monitoring.netcracker.com,resources=platformmonitorings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=monitoring.netcracker.com,resources=platformmonitorings/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=monitoring.netcracker.com,resources=platformmonitorings/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PlatformMonitoring object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/reconcile
func (r *PlatformMonitoringReconciler) Reconcile(context context.Context, request ctrl.Request) (ctrl.Result, error) {
	// Fetch the PlatformMonitoring instance
	customResourceInstance := &qubershiporgv1.PlatformMonitoring{}
	err := r.Client.Get(context, request.NamespacedName, customResourceInstance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{Requeue: false}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}
	customResourceInstance.FillEmptyWithDefaults()

	r.Log.Info("Reconciliation started")
	if err = r.updateStatus(customResourceInstance, "In progress", "False", "ReconcileCycleStatus", "Monitoring service reconcile cycle in progress"); err != nil {
		r.Log.Error(err, "Error while update status")
	}

	// Prometheus Operator should create first because it should create CRDs:
	// * Prometheus
	// * ServiceMonitor
	// * PodMonitor
	// * Alertmanager
	// * PrometheusRule
	poReconciler := prometheusoperator.NewPrometheusOperatorReconciler(r.Client, r.Scheme)
	err = poReconciler.Run(customResourceInstance)
	if err != nil {
		r.Log.Error(err, "Reconciliation of prometheus-operator failed")
		r.prepareStatusForUpdate(customResourceInstance, "Failed", "False", "ReconcilePrometheusOperatorStatus", "Prometheus operator reconcile cycle failed")
	} else {
		r.removeStatus(customResourceInstance, "ReconcilePrometheusOperatorStatus")
	}

	kubernetesMonitorsReconciler := kubernetesmonitors.NewKubernetesMonitorsReconciler(r.Client, r.Scheme, r.DiscoveryClient)
	err = kubernetesMonitorsReconciler.Run(customResourceInstance)
	if err != nil {
		r.Log.Error(err, "Reconciliation of kubernetes-monitors failed")
	}

	// Reconcile Victoriametrics Operator custom resources
	vmoReconciler := vmoperator.NewVmOperatorReconciler(r.Client, r.Scheme, r.Config, r.DiscoveryClient)
	err = vmoReconciler.Run(customResourceInstance)
	if err != nil {
		r.Log.Error(err, "Reconciliation of vm-operator failed")
		r.prepareStatusForUpdate(customResourceInstance, "Failed", "False", "ReconcileVictoriametricsOperatorStatus", "Victoriametrics operator reconcile cycle failed")
	} else {
		r.removeStatus(customResourceInstance, "ReconcileVictoriametricsOperatorStatus")
	}

	// Reconcile vmSingle Operator custom resources
	vmSingleReconciler := vmsingle.NewVmSingleReconciler(r.Client, r.Scheme, r.DiscoveryClient)
	err = vmSingleReconciler.Run(context, customResourceInstance)
	if err != nil {
		r.Log.Error(err, "Reconciliation of vmsingle failed")
		r.prepareStatusForUpdate(customResourceInstance, "Failed", "False", "ReconcileVictoriametricsSingleStatus", "Victoriametrics Single reconcile cycle failed")
	} else {
		r.removeStatus(customResourceInstance, "ReconcileVictoriametricsSingleStatus")
	}

	// Reconcile VmCluster Operator custom resources
	vmclusterReconciler := vmcluster.NewVmClusterReconciler(r.Client, r.Scheme, r.DiscoveryClient)
	err = vmclusterReconciler.Run(context, customResourceInstance)
	if err != nil {
		r.Log.Error(err, "Reconciliation of vmcluster failed")
		r.prepareStatusForUpdate(customResourceInstance, "Failed", "False", "ReconcileVictoriametricsClusterStatus", "Victoriametrics Cluster reconcile cycle failed")
	} else {
		r.removeStatus(customResourceInstance, "ReconcileVictoriametricsClusterStatus")
	}

	// Reconcile vmUser Operator custom resources
	vmUserReconciler := vmuser.NewVmUserReconciler(r.Client, r.Scheme, r.DiscoveryClient)
	err = vmUserReconciler.Run(context, customResourceInstance)
	if err != nil {
		r.Log.Error(err, "Reconciliation of vmuser failed")
		r.prepareStatusForUpdate(customResourceInstance, "Failed", "False", "ReconcileVictoriametricsUserStatus", "Victoriametrics User reconcile cycle failed")
	} else {
		r.removeStatus(customResourceInstance, "ReconcileVictoriametricsUserStatus")
	}

	// Reconcile VmAgent Operator custom resources
	vmagentReconciler := vmagent.NewVmAgentReconciler(r.Client, r.Scheme, r.DiscoveryClient)
	err = vmagentReconciler.Run(context, customResourceInstance)
	if err != nil {
		r.Log.Error(err, "Reconciliation of vmagent failed")
		r.prepareStatusForUpdate(customResourceInstance, "Failed", "False", "ReconcileVictoriametricsAgentStatus", "Victoriametrics Agent reconcile cycle failed")
	} else {
		r.removeStatus(customResourceInstance, "ReconcileVictoriametricsAgentStatus")
	}

	// Reconcile vmAuth Operator custom resources
	vmAuthReconciler := vmauth.NewVmAuthReconciler(r.Client, r.Scheme, r.DiscoveryClient)
	err = vmAuthReconciler.Run(context, customResourceInstance)
	if err != nil {
		r.Log.Error(err, "Reconciliation of vmauth failed")
		r.prepareStatusForUpdate(customResourceInstance, "Failed", "False", "ReconcileVictoriametricsAuthStatus", "Victoriametrics Auth reconcile cycle failed")
	} else {
		r.removeStatus(customResourceInstance, "ReconcileVictoriametricsAuthStatus")
	}

	//Reconcile Prometheus Operator custom resources
	pReconciler := prometheus.NewPrometheusReconciler(r.Client, r.Scheme, r.DiscoveryClient)
	err = pReconciler.Run(customResourceInstance)
	if err != nil {
		r.Log.Error(err, "Reconciliation of prometheus custom resources failed")
		r.prepareStatusForUpdate(customResourceInstance, "Failed", "False", "ReconcilePrometheusStatus", "Prometheus reconcile cycle failed")
	} else {
		r.removeStatus(customResourceInstance, "ReconcilePrometheusStatus")
	}

	// Reconcile vmAlertManager custom resources
	vmAlertManagerReconciler := vmalertmanager.NewVmAlertManagerReconciler(r.Client, r.Scheme, r.DiscoveryClient)
	err = vmAlertManagerReconciler.Run(context, customResourceInstance)
	if err != nil {
		r.Log.Error(err, "Reconciliation of vmalertmanager failed")
		r.prepareStatusForUpdate(customResourceInstance, "Failed", "False", "ReconcileVictoriametricsAlertManagerStatus", "Victoriametrics Alert Manager reconcile cycle failed")
	} else {
		r.removeStatus(customResourceInstance, "ReconcileVictoriametricsAlertManagerStatus")
	}

	aReconciler := alertmanager.NewAlertManagerReconciler(r.Client, r.Scheme, r.DiscoveryClient)
	err = aReconciler.Run(customResourceInstance)
	if err != nil {
		r.Log.Error(err, "Reconciliation of alertmanager custom resources failed")
		r.prepareStatusForUpdate(customResourceInstance, "Failed", "False", "ReconcileAlertManagerStatus", "AlertManager reconcile cycle failed")
	} else {
		r.removeStatus(customResourceInstance, "ReconcileAlertManagerStatus")
	}

	// Reconcile vmAlert custom resources
	vmAlertReconciler := vmalert.NewVmAlertReconciler(r.Client, r.Scheme, r.DiscoveryClient)
	err = vmAlertReconciler.Run(context, customResourceInstance)
	if err != nil {
		r.Log.Error(err, "Reconciliation of vmalert failed")
		r.prepareStatusForUpdate(customResourceInstance, "Failed", "False", "ReconcileVictoriametricsAlertStatus", "Victoriametrics Alert reconcile cycle failed")
	} else {
		r.removeStatus(customResourceInstance, "ReconcileVictoriametricsAlertStatus")
	}

	//// Reconcile Exporters
	ksmReconciler := kubestatemetrics.NewKubeStateMetricsReconciler(r.Client, r.Scheme, r.DiscoveryClient)
	err = ksmReconciler.Run(customResourceInstance)
	if err != nil {
		r.Log.Error(err, "Reconciliation of kube-state-metrics failed")
		r.prepareStatusForUpdate(customResourceInstance, "Failed", "False", "ReconcileKubeStateMetricsStatus", "KubeStateMetrics reconcile cycle failed")
	} else {
		r.removeStatus(customResourceInstance, "ReconcileKubeStateMetricsStatus")
	}

	neReconciler := nodeexporter.NewNodeExporterReconciler(r.Client, r.Scheme, r.DiscoveryClient)
	err = neReconciler.Run(customResourceInstance)
	if err != nil {
		r.Log.Error(err, "Reconciliation of node-exporter failed")
		r.prepareStatusForUpdate(customResourceInstance, "Failed", "False", "ReconcileNodeExporterStatus", "NodeExporter reconcile cycle failed")
	} else {
		r.removeStatus(customResourceInstance, "ReconcileNodeExporterStatus")
	}

	//// Grafana Operator should create first because it should create CRDs:
	//// * Grafana
	//// * GrafanaDatasource
	//// * GrafanaDashboard
	goReconciler := grafanaoperator.NewGrafanaOperatorReconciler(r.Client, r.Scheme, r.DiscoveryClient)
	err = goReconciler.Run(customResourceInstance)
	if err != nil {
		r.Log.Error(err, "Reconciliation of grafana-operator failed")
		r.prepareStatusForUpdate(customResourceInstance, "Failed", "False", "ReconcileGrafanaOperatorStatus", "Grafana operator reconcile cycle failed")
	} else {
		r.removeStatus(customResourceInstance, "ReconcileGrafanaOperatorStatus")
	}

	gReconciler := grafana.NewGrafanaReconciler(r.Client, r.Scheme, r.DiscoveryClient, r.Config)
	err = gReconciler.Run(customResourceInstance)
	if err != nil {
		r.Log.Error(err, "Reconciliation of grafana custom resources failed")
		r.prepareStatusForUpdate(customResourceInstance, "Failed", "False", "ReconcileGrafanaStatus", "Grafana reconcile cycle failed")
	} else {
		r.removeStatus(customResourceInstance, "ReconcileGrafanaStatus")
	}

	prReconciler := prometheusrules.NewPrometheusRulesReconciler(r.Client, r.Scheme)
	err = prReconciler.Run(customResourceInstance)
	if err != nil {
		r.Log.Error(err, "Reconciliation of prometheus rules custom resources failed")
		r.prepareStatusForUpdate(customResourceInstance, "Failed", "False", "ReconcilePrometheusRulesStatus", "Prometheus Rules reconcile cycle failed")
	} else {
		r.removeStatus(customResourceInstance, "ReconcilePrometheusRulesStatus")
	}

	pgReconciler := pushgateway.NewPushgatewayReconciler(r.Client, r.Scheme, r.DiscoveryClient)
	err = pgReconciler.Run(customResourceInstance)
	if err != nil {
		r.Log.Error(err, "Reconciliation of pushgateway failed")
		r.prepareStatusForUpdate(customResourceInstance, "Failed", "False", "ReconcilePushgatewayStatus", "Pushgateway reconcile cycle failed")
	} else {
		r.removeStatus(customResourceInstance, "ReconcilePushgatewayStatus")
	}

	rInterval, err := strconv.ParseInt(utils.GetEnvWithDefaultValue("RECONCILIATION_INTERVAL"), 10, 64)
	if err != nil {
		return reconcile.Result{}, err
	}

	status := customResourceInstance.Status.Conditions
	failedStatus := false

	for i := range status {
		if status[i].Type == "Failed" {
			failedStatus = true
			break
		}
	}

	if failedStatus {
		r.prepareStatusForUpdate(customResourceInstance, "Failed", "False", "ReconcileCycleStatus", "Monitoring service reconcile cycle failed")

		if err = r.Client.Status().Update(context, customResourceInstance); err != nil {
			r.Log.Error(err, "Update status failed.")
		}

		r.Log.Info("Reconciliation failed. Run reconciliation again.")
		return reconcile.Result{Requeue: true}, nil
	}

	r.prepareStatusForUpdate(customResourceInstance, "Successful", "True", "ReconcileCycleStatus", "Monitoring service reconcile cycle succeeded")
	r.Log.Info("Reconciliation finished successful, next reconciliation after " + utils.GetEnvWithDefaultValue("RECONCILIATION_INTERVAL") + " seconds")

	if err = r.Client.Status().Update(context, customResourceInstance); err != nil {
		r.Log.Error(err, "Update status failed")
	}

	return reconcile.Result{RequeueAfter: time.Duration(rInterval) * time.Second}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PlatformMonitoringReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&qubershiporgv1.PlatformMonitoring{}).
		WithEventFilter(ignoreDeletionPredicate()).
		Complete(r)
}

func ignoreDeletionPredicate() predicate.Predicate {
	return predicate.Funcs{
		UpdateFunc: func(e event.UpdateEvent) bool {
			// Ignore updates to CR status in which case metadata.Generation does not change
			return e.ObjectOld.GetGeneration() != e.ObjectNew.GetGeneration()
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			// Evaluates to false if the object has been confirmed deleted.
			return !e.DeleteStateUnknown
		},
	}
}

// getCondition retrieves condition of custom resource instance by given reason.
// Returns found condition and it's index in the list of conditions.
func (r *PlatformMonitoringReconciler) getCondition(customResourceInstance *qubershiporgv1.PlatformMonitoring, newConditionReason string) (int, *qubershiporgv1.PlatformMonitoringCondition) {

	if len(newConditionReason) == 0 || len(customResourceInstance.Status.Conditions) == 0 {
		return -1, nil
	}
	for i := range customResourceInstance.Status.Conditions {
		if customResourceInstance.Status.Conditions[i].Reason == newConditionReason {
			return i, &customResourceInstance.Status.Conditions[i]
		}
	}
	return -1, nil
}

// removeStatus removes condition of custom resource instance by given reason.
func (r *PlatformMonitoringReconciler) removeStatus(customResourceInstance *qubershiporgv1.PlatformMonitoring, reason string) bool {
	idx, condition := r.getCondition(customResourceInstance, reason)
	if condition != nil {
		customResourceInstance.Status.Conditions[idx] = customResourceInstance.Status.Conditions[len(customResourceInstance.Status.Conditions)-1]
		customResourceInstance.Status.Conditions = customResourceInstance.Status.Conditions[:len(customResourceInstance.Status.Conditions)-1]
		return true
	}
	return false
}

// updateStatus updates condition of custom resource instance
func (r *PlatformMonitoringReconciler) updateStatus(customResourceInstance *qubershiporgv1.PlatformMonitoring, statusType string, status string, reason string, message string) error {
	if r.prepareStatusForUpdate(customResourceInstance, statusType, status, reason, message) {
		// Update status if not equal to the last one
		if err := r.Client.Status().Update(context.TODO(), customResourceInstance); err != nil {
			r.Log.Error(err, "Update status failed")
			return err
		}
	}
	return nil
}

// prepareStatusForUpdate checks if old condition with same reason exist and change it on the new condition
func (r *PlatformMonitoringReconciler) prepareStatusForUpdate(customResourceInstance *qubershiporgv1.PlatformMonitoring, statusType string, status string, reason string, message string) bool {
	// get status timestamp
	transitionTime := metav1.Now()
	newCondition := qubershiporgv1.PlatformMonitoringCondition{
		Type:               statusType,
		Status:             status,
		Reason:             reason,
		Message:            message,
		LastTransitionTime: transitionTime.String(),
	}
	idx, oldCondition := r.getCondition(customResourceInstance, newCondition.Reason)

	if oldCondition == nil {
		customResourceInstance.Status.Conditions = append(customResourceInstance.Status.Conditions, newCondition)
		return true
	}

	isEqual := newCondition.Type == oldCondition.Type &&
		newCondition.Status == oldCondition.Status &&
		newCondition.Reason == oldCondition.Reason &&
		newCondition.Message == oldCondition.Message &&
		newCondition.LastTransitionTime == oldCondition.LastTransitionTime

	if !isEqual {
		//replace old status
		customResourceInstance.Status.Conditions[idx] = newCondition
	}
	return !isEqual
}
