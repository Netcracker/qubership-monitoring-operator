package grafana_operator

import (
	v1alpha1 "github.com/Netcracker/qubership-monitoring-operator/api/v1alpha1"
	kubernetes_monitors "github.com/Netcracker/qubership-monitoring-operator/controllers/kubernetes-monitors"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/discovery"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type GrafanaOperatorReconciler struct {
	*utils.ComponentReconciler
}

func NewGrafanaOperatorReconciler(c client.Client, s *runtime.Scheme, dc discovery.DiscoveryInterface) *GrafanaOperatorReconciler {
	return &GrafanaOperatorReconciler{
		ComponentReconciler: &utils.ComponentReconciler{
			Client: c,
			Scheme: s,
			Dc:     dc,
			Log:    utils.Logger("grafanaoperator_reconciler"),
		},
	}
}

// Run reconciles grafana-operator.
// Creates new cluster role, cluster role binding and grafana custom resource if its don't exists.
// Updates deployment and service in case of any changes.
// Returns true if need to requeue, false otherwise.
func (r *GrafanaOperatorReconciler) Run(cr *v1alpha1.PlatformMonitoring) error {
	r.Log.Info("Reconciling component")

	restrictedDashboards := make(map[string]bool)
	if cr.Spec.PublicCloudName != "" {
		if pcMap, ok := utils.PublicCloudDashboardsEnabled[cr.Spec.PublicCloudName]; ok {
			for dashboard, included := range pcMap {
				if included {
					cr.Spec.GrafanaDashboards.List = append(cr.Spec.GrafanaDashboards.List, dashboard)
				} else {
					restrictedDashboards[dashboard] = true
				}
			}
		}
	}
	// Create dashboards after we have created CRD-s
	if cr.Spec.GrafanaDashboards != nil && cr.Spec.GrafanaDashboards.IsInstall() && len(cr.Spec.GrafanaDashboards.List) > 0 {

		// Nginx dashboards should be installed only if NginxIngressPodMonitor is installed
		isNginxIngressPodMonitorInstalled := isMonitorInstall(cr, utils.NginxIngressPodMonitorName)

		// CoreDns dashboards should be installed only if CoreDnsServiceMonitor is installed
		isCoreDnsServiceMonitorInstalled := isMonitorInstall(cr, utils.CoreDnsServiceMonitorName)

		if cr.Spec.Prometheus != nil && cr.Spec.Prometheus.IsInstall() {
			cr.Spec.GrafanaDashboards.List = append(cr.Spec.GrafanaDashboards.List, "prometheus-self-monitoring")
		}

		if cr.Spec.Victoriametrics != nil && cr.Spec.Victoriametrics.VmOperator.IsInstall() {
			cr.Spec.GrafanaDashboards.List = append(cr.Spec.GrafanaDashboards.List, "victoriametrics-vmoperator")
			if cr.Spec.Victoriametrics.VmAlert.IsInstall() {
				cr.Spec.GrafanaDashboards.List = append(cr.Spec.GrafanaDashboards.List, "victoriametrics-vmalert")
			}
			if cr.Spec.Victoriametrics.VmAgent.IsInstall() {
				cr.Spec.GrafanaDashboards.List = append(cr.Spec.GrafanaDashboards.List, "victoriametrics-vmagent")
			}
			if cr.Spec.Victoriametrics.VmSingle.IsInstall() {
				cr.Spec.GrafanaDashboards.List = append(cr.Spec.GrafanaDashboards.List, "victoriametrics-vmsingle")
			}
		}
		// Pod distribution by zones dashboard
		if cr.Spec.PublicCloudName != "" && cr.Spec.PublicCloudName == "aws" {
			cr.Spec.GrafanaDashboards.List = append(cr.Spec.GrafanaDashboards.List, "kubernetes-pods-distribution-by-zone")
		}

		isOpenshiftV4, err := r.IsOpenShiftV4()
		if err != nil {
			r.Log.Error(err, "Failed to recognize OpenShift V4")
		}

		if isOpenshiftV4 {
			cr.Spec.GrafanaDashboards.List = append(cr.Spec.GrafanaDashboards.List, "openshift-apiserver")
			cr.Spec.GrafanaDashboards.List = append(cr.Spec.GrafanaDashboards.List, "openshift-state-metrics")
			cr.Spec.GrafanaDashboards.List = append(cr.Spec.GrafanaDashboards.List, "openshift-cluster-version-operator")
			cr.Spec.GrafanaDashboards.List = append(cr.Spec.GrafanaDashboards.List, "openshift-haproxy")
		}

		// Create map for fast search for dashboards chosen in cr
		dashboardsToInstall := make(map[string]string)
		for _, dashboardName := range cr.Spec.GrafanaDashboards.List {
			dashboardsToInstall[dashboardName+".yaml"] = dashboardName
		}

		isOpenshiftV3, err := r.IsOpenShiftV3()
		if err != nil {
			r.Log.Error(err, "Failed to recognize OpenShift V3.11")
		}

		for _, mResource := range utils.GrafanaKubernetesDashboardsResources {
			if _, ok := dashboardsToInstall[mResource]; ok {
				// Dashboards that should not be in the public cloud will be skipped
				if restrictedDashboards[dashboardsToInstall[mResource]] {
					if err = r.deleteGrafanaDashboard(mResource, cr); err != nil {
						r.Log.Error(err, "Can not delete GrafanaDashboard")
					}
					continue
				}
				// Create node-exporter dashboard only if we install it
				if mResource == utils.GrafanaNodeExporterDashboardResource && (cr.Spec.NodeExporter == nil || !cr.Spec.NodeExporter.IsInstall()) {
					r.Log.Info("Delete dashboard " + utils.GrafanaNodeExporterDashboardResource + " if exists and NodeExporter is not installed")
					if err = r.deleteGrafanaDashboard(mResource, cr); err != nil {
						r.Log.Error(err, "Can not delete GrafanaDashboard")
					}
					continue
				}
				// Create default home dashboard only if grafanaHomeDashboard parameter set to true
				if mResource == utils.GrafanaHomeDashboardResource && cr.Spec.Grafana != nil && !cr.Spec.Grafana.GrafanaHomeDashboard {
					r.Log.Info("Delete dashboard " + utils.GrafanaHomeDashboardResource + " if exists and grafanaHomeDashboard set to false")
					if err = r.deleteGrafanaDashboard(mResource, cr); err != nil {
						r.Log.Error(err, "Can not delete GrafanaDashboard")
					}
					continue
				}
				// Create Grafana Self Monitoring dashboard if grafana.install parameter set to true
				if mResource == utils.GrafanaSelfMonitoringDashboardResource && (cr.Spec.Grafana == nil || !cr.Spec.Grafana.IsInstall()) {
					r.Log.Info("Delete dashboard " + utils.GrafanaSelfMonitoringDashboardResource + " if exists and Grafana Overview set to false")
					if err = r.deleteGrafanaDashboard(mResource, cr); err != nil {
						r.Log.Error(err, "Can not delete GrafanaDashboard")
					}
					continue
				}
				// Create Alertmanager dashboard if alertmanager.install parameter set to true
				if mResource == utils.GrafanaAlertmanagerDashboardResource && ((cr.Spec.AlertManager == nil || !cr.Spec.AlertManager.IsInstall()) && (cr.Spec.Victoriametrics == nil || !cr.Spec.Victoriametrics.VmAlertManager.IsInstall())) {
					r.Log.Info("Delete dashboard " + utils.GrafanaAlertmanagerDashboardResource + " if exists and Alertmanager Overview set to false")
					if err = r.deleteGrafanaDashboard(mResource, cr); err != nil {
						r.Log.Error(err, "Can not delete GrafanaDashboard")
					}
					continue
				}
				// Create CoreDNS dashboards only if we install CoreDnsServiceMonitor and not in the OpenShift v3.11
				if mResource == utils.GrafanaCoreDnsDashboardResource {
					if !isCoreDnsServiceMonitorInstalled {
						r.Log.Info("Delete dashboard " + utils.GrafanaCoreDnsDashboardResource + " if exists and CoreDnsServiceMonitor is not installed")
						if err = r.deleteGrafanaDashboard(mResource, cr); err != nil {
							r.Log.Error(err, "Can not delete GrafanaDashboard")
						}
						continue
					}

					if isOpenshiftV3 {
						r.Log.Info("Delete dashboard " + utils.GrafanaCoreDnsDashboardResource + " if exists and cluster based on OpenShift < 4.1")
						if err = r.deleteGrafanaDashboard(mResource, cr); err != nil {
							r.Log.Error(err, "Can not delete GrafanaDashboard")
						}
						continue
					}
				}
				// Create Nginx Ingress dashboards only if we install NginxIngressPodMonitor and not in the OpenShift
				if mResource == utils.GrafanaNginxRequestHandlingPerformanceDashboardResource {
					if !isNginxIngressPodMonitorInstalled {
						r.Log.Info("Delete dashboard " + utils.GrafanaNginxRequestHandlingPerformanceDashboardResource + " if exists and NginxIngressPodMonitor is not installed")
						if err = r.deleteGrafanaDashboard(mResource, cr); err != nil {
							r.Log.Error(err, "Can not delete GrafanaDashboard")
						}
						continue
					}
					if r.HasRouteApi() {
						r.Log.Info("Delete dashboard " + utils.GrafanaNginxRequestHandlingPerformanceDashboardResource + " if exists and cluster based on OpenShift")
						if err = r.deleteGrafanaDashboard(mResource, cr); err != nil {
							r.Log.Error(err, "Can not delete GrafanaDashboard")
						}
						continue
					}
				}
				if mResource == utils.GrafanaNginxIngressDashboardResource {
					if !isNginxIngressPodMonitorInstalled {
						r.Log.Info("Delete dashboard " + utils.GrafanaNginxIngressDashboardResource + " if exists and NginxIngressPodMonitor is not installed")
						if err = r.deleteGrafanaDashboard(mResource, cr); err != nil {
							r.Log.Error(err, "Can not delete GrafanaDashboard")
						}
						continue
					}
					if r.HasRouteApi() {
						r.Log.Info("Delete dashboard " + utils.GrafanaNginxIngressDashboardResource + " if exists and cluster based on OpenShift")
						if err = r.deleteGrafanaDashboard(mResource, cr); err != nil {
							r.Log.Error(err, "Can not delete GrafanaDashboard")
						}
						continue
					}
				}
				if mResource == utils.GrafanaNginxIngressListOfIngresses {
					if !isNginxIngressPodMonitorInstalled {
						r.Log.Info("Delete dashboard " + utils.GrafanaNginxIngressListOfIngresses + " if exists and NginxIngressPodMonitor is not installed")
						if err = r.deleteGrafanaDashboard(mResource, cr); err != nil {
							r.Log.Error(err, "Can not delete GrafanaDashboard")
						}
						continue
					}
					if r.HasRouteApi() {
						r.Log.Info("Delete dashboard " + utils.GrafanaNginxIngressListOfIngresses + " if exists and cluster based on OpenShift")
						if err = r.deleteGrafanaDashboard(mResource, cr); err != nil {
							r.Log.Error(err, "Can not delete GrafanaDashboard")
						}
						continue
					}
				}

				if err = r.handleGrafanaDashboard(mResource, cr); err != nil {
					return err
				}
			} else {
				if err = r.deleteGrafanaDashboard(mResource, cr); err != nil {
					r.Log.Error(err, "Can not delete GrafanaDashboard")
				}
			}
		}
	} else {
		r.Log.Info("Remove all grafana dashboards")
		r.uninstallGrafanaDashboards(cr)
	}
	if cr.Spec.Grafana != nil && cr.Spec.Grafana.IsInstall() {
		if !cr.Spec.Grafana.Operator.Paused {
			err := r.handleServiceAccount(cr)
			if err != nil {
				return err
			}
			// Reconcile resources with creation
			if utils.PrivilegedRights {
				if err = r.handleClusterRole(cr); err != nil {
					return err
				}
				if err = r.handleClusterRoleBinding(cr); err != nil {
					return err
				}
			}
			// Reconcile resources with creation
			if err = r.handleRole(cr); err != nil {
				return err
			}
			if err = r.handleRoleBinding(cr); err != nil {
				return err
			}

			// Reconcile resources with creation and update
			if err = r.handleDeployment(cr); err != nil {
				return err
			}

			// Reconcile Pod Monitor
			if cr.Spec.Grafana.Operator.PodMonitor != nil && cr.Spec.Grafana.Operator.PodMonitor.IsInstall() {
				if err = r.handlePodMonitor(cr); err != nil {
					return err
				}
			} else {
				r.Log.Info("Uninstalling PodMonitor")
				if err = r.deletePodMonitor(cr); err != nil {
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

func (r *GrafanaOperatorReconciler) uninstallGrafanaDashboards(cr *v1alpha1.PlatformMonitoring) {
	for _, mResource := range utils.GrafanaKubernetesDashboardsResources {
		if err := r.deleteGrafanaDashboard(mResource, cr); err != nil {
			r.Log.Error(err, "Can not delete GrafanaDashboard")
		}
	}
}

// uninstall deletes all resources related to the component
func (r *GrafanaOperatorReconciler) uninstall(cr *v1alpha1.PlatformMonitoring) {
	if err := r.deleteGrafanaOperatorDeployment(cr); err != nil {
		r.Log.Error(err, "Can not delete Deployment")
	}
	if err := r.deletePodMonitor(cr); err != nil {
		r.Log.Error(err, "Can not delete PodMonitor")
	}
}

// isMonitorInstall returns "true" if the monitor installed
func isMonitorInstall(cr *v1alpha1.PlatformMonitoring, monitorName string) bool {
	// If the monitor affected by installation in the public cloud, cr.Spec.KubernetesMonitors doesn't matter
	affected, installed := kubernetes_monitors.IsMonitorPresentInPublicCloud(cr, monitorName)
	if affected {
		return installed
	} else if len(cr.Spec.KubernetesMonitors) > 0 {
		monitor, ok := cr.Spec.KubernetesMonitors[monitorName]
		if ok && monitor.IsInstall() {
			return true
		}
	}
	return false
}
