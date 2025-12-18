package kubernetes_monitors

import (
	v1beta1 "github.com/Netcracker/qubership-monitoring-operator/api/v1beta1"
	"github.com/Netcracker/qubership-monitoring-operator/controllers/utils"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/discovery"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// KubernetesMonitorsReconciler is a reconciler to maintain configuration of k8s service monitors
type KubernetesMonitorsReconciler struct {
	*utils.ComponentReconciler
}

// NewKubernetesMonitorsReconciler  returns KubernetesMonitorsReconciler by specified parameters
func NewKubernetesMonitorsReconciler(c client.Client, s *runtime.Scheme, dc discovery.DiscoveryInterface) *KubernetesMonitorsReconciler {
	return &KubernetesMonitorsReconciler{
		ComponentReconciler: &utils.ComponentReconciler{
			Client: c,
			Scheme: s,
			Dc:     dc,
			Log:    utils.Logger("kubernetes_monitors_reconciler"),
		},
	}
}

// Run reconciles k8s service monitors
// Creates, updates and deletes service monitors for k8s monitoring depending of configurtion
func (r *KubernetesMonitorsReconciler) Run(cr *v1beta1.PlatformMonitoring) error {
	r.Log.Info("Reconciling component")

	if len(cr.Spec.KubernetesMonitors) > 0 {
		isOpenshiftV4, err := r.IsOpenShiftV4()
		if err != nil {
			r.Log.Error(err, "Failed to recognise OpenShift V4")
		}

		isOpenshiftV3, err := r.IsOpenShiftV3()
		if err != nil {
			r.Log.Error(err, "Failed to recognise OpenShift V3.11")
		}

		if IsMonitorInstall(cr, utils.ApiserverServiceMonitorName) {
			if err = r.handleApiServerServiceMonitor(cr); err != nil {
				return err
			}
		} else {
			if err = r.deleteApiServerServiceMonitor(cr); err != nil {
				r.Log.Error(err, "Can not delete ApiServerServiceMonitor")
			}
		}
		if IsMonitorInstall(cr, utils.KubeControllerManagerServiceMonitorName) {
			if err = r.handleControllerManagerServiceMonitor(cr); err != nil {
				return err
			}
		} else {
			if err = r.deleteControllerManagerServiceMonitor(cr); err != nil {
				r.Log.Error(err, "Can not delete ControllerManagerServiceMonitor")
			}
		}
		if IsMonitorInstall(cr, utils.KubeSchedulerServiceMonitorName) {
			if err = r.handleSchedulerServiceMonitor(cr); err != nil {
				return err
			}
		} else {
			if err = r.deleteSchedulerServiceMonitor(cr); err != nil {
				r.Log.Error(err, "Can not delete SchedulerServiceMonitor")
			}
		}
		if IsMonitorInstall(cr, utils.KubeletServiceMonitorName) {
			if err = r.handleKubeletServiceMonitor(cr); err != nil {
				return err
			}
		} else {
			if err = r.deleteKubeletServiceMonitor(cr); err != nil {
				r.Log.Error(err, "Can not delete KubeletServiceMonitor")
			}
		}
		if !isOpenshiftV3 && IsMonitorInstall(cr, utils.CoreDnsServiceMonitorName) {
			if err = r.handleCoreDnsServiceMonitor(cr, isOpenshiftV4); err != nil {
				return err
			}
		} else {
			if err = r.deleteCoreDnsServiceMonitor(cr, isOpenshiftV4); err != nil {
				r.Log.Error(err, "Can not delete CoreDnsServiceMonitor")
			}
		}
		if !r.HasRouteApi() && IsMonitorInstall(cr, utils.NginxIngressPodMonitorName) {
			if err = r.handleNginxIngressPodMonitor(cr); err != nil {
				return err
			}
		} else {
			if err = r.deleteNginxIngressPodMonitor(cr); err != nil {
				r.Log.Error(err, "Can not delete NginxIngressPodMonitor")
			}
		}
		if isOpenshiftV4 {
			if IsMonitorInstall(cr, utils.OpenshiftApiServerServiceMonitorName) {
				if err = r.handleOpenshiftServiceMonitor(cr, utils.OpenshiftApiServerServiceMonitorAsset,
					utils.OpenshiftApiServerServiceMonitorName,
					utils.OpenshiftApiserver); err != nil {
					return err
				}
			} else {
				if err = r.deleteOpenshiftServiceMonitor(cr, utils.OpenshiftApiServerServiceMonitorAsset,
					utils.OpenshiftApiServerServiceMonitorName,
					utils.OpenshiftApiserver); err != nil {
					r.Log.Error(err, "Can not delete OpenshiftApiServerServiceMonitor")
				}
			}
			if IsMonitorInstall(cr, utils.OpenshiftApiServerOperatorServiceMonitorName) {
				if err = r.handleOpenshiftServiceMonitor(cr, utils.OpenshiftApiServerOperatorServiceMonitorAsset,
					utils.OpenshiftApiServerOperatorServiceMonitorName,
					utils.OpenshiftApiServerOperator); err != nil {
					return err
				}
			} else {
				if err = r.deleteOpenshiftServiceMonitor(cr, utils.OpenshiftApiServerOperatorServiceMonitorAsset,
					utils.OpenshiftApiServerOperatorServiceMonitorName,
					utils.OpenshiftApiServerOperator); err != nil {
					r.Log.Error(err, "Can not delete OpenshiftApiServerOperatorServiceMonitor")
				}
			}
			if IsMonitorInstall(cr, utils.OpenshiftClusterVersionOperatorServiceMonitorName) {
				if err = r.handleOpenshiftServiceMonitor(cr, utils.OpenshiftClusterVersionOperatorServiceMonitorAsset,
					utils.OpenshiftClusterVersionOperatorServiceMonitorName,
					utils.OpenshiftClusterVersionOperator); err != nil {
					return err
				}
			} else {
				if err = r.deleteOpenshiftServiceMonitor(cr, utils.OpenshiftClusterVersionOperatorServiceMonitorAsset,
					utils.OpenshiftClusterVersionOperatorServiceMonitorName,
					utils.OpenshiftClusterVersionOperator); err != nil {
					r.Log.Error(err, "Can not delete OpenshiftClusterVersionOperatorServiceMonitor")
				}
			}
			if IsMonitorInstall(cr, utils.OpenshiftStatemetricsServiceMonitorName) {
				if err = r.handleOpenshiftServiceMonitor(cr, utils.OpenshiftStatemetricsServiceMonitorAsset,
					utils.OpenshiftStatemetricsServiceMonitorName,
					utils.OpenshiftStatemetrics); err != nil {
					return err
				}
			} else {
				if err = r.deleteOpenshiftServiceMonitor(cr, utils.OpenshiftStatemetricsServiceMonitorAsset,
					utils.OpenshiftStatemetricsServiceMonitorName,
					utils.OpenshiftStatemetrics); err != nil {
					r.Log.Error(err, "Can not delete OpenshiftStatemetricsServiceMonitor")
				}
			}
			if IsMonitorInstall(cr, utils.OpenshiftHAProxyServiceMonitorName) {
				if err = r.handleOpenshiftServiceMonitor(cr, utils.OpenshiftHAProxyServiceMonitorAsset,
					utils.OpenshiftHAProxyServiceMonitorName,
					utils.OpenshiftHAProxy); err != nil {
					return err
				}
			} else {
				if err = r.deleteOpenshiftServiceMonitor(cr, utils.OpenshiftHAProxyServiceMonitorAsset,
					utils.OpenshiftHAProxyServiceMonitorName,
					utils.OpenshiftHAProxy); err != nil {
					r.Log.Error(err, "Can not delete OpenshiftHAProxyServiceMonitor")
				}
			}
		}
	} else {
		r.Log.Info("Uninstalling component if exists")
		r.uninstall(cr)
	}

	r.Log.Info("Component reconciled")
	return nil
}

// uninstall deletes all resources related to the component
func (r *KubernetesMonitorsReconciler) uninstall(cr *v1beta1.PlatformMonitoring) {

	isOpenshiftV4, err := r.IsOpenShiftV4()
	if err != nil {
		r.Log.Error(err, "Failed to recognise OpenShift V4, continue for OpenShift v3.11")
	}

	if err = r.deleteApiServerServiceMonitor(cr); err != nil {
		r.Log.Error(err, "Can not delete ApiServerServiceMonitor")
	}
	if err = r.deleteControllerManagerServiceMonitor(cr); err != nil {
		r.Log.Error(err, "Can not delete ControllerManagerServiceMonitor")
	}
	if err = r.deleteSchedulerServiceMonitor(cr); err != nil {
		r.Log.Error(err, "Can not delete SchedulerServiceMonitor")
	}
	if err = r.deleteKubeletServiceMonitor(cr); err != nil {
		r.Log.Error(err, "Can not delete KubeletServiceMonitor")
	}
	if err = r.deleteCoreDnsServiceMonitor(cr, isOpenshiftV4); err != nil {
		r.Log.Error(err, "Can not delete CoreDnsServiceMonitor")
	}
	if err = r.deleteNginxIngressPodMonitor(cr); err != nil {
		r.Log.Error(err, "Can not delete NginxIngressPodMonitor")
	}

	if isOpenshiftV4 {
		if err = r.deleteOpenshiftServiceMonitor(cr, utils.OpenshiftApiServerServiceMonitorAsset,
			utils.OpenshiftApiServerServiceMonitorName,
			utils.OpenshiftApiserver); err != nil {
			r.Log.Error(err, "Can not delete OpenshiftApiServerServiceMonitor")
		}
		if err = r.deleteOpenshiftServiceMonitor(cr, utils.OpenshiftApiServerOperatorServiceMonitorAsset,
			utils.OpenshiftApiServerOperatorServiceMonitorName,
			utils.OpenshiftApiServerOperator); err != nil {
			r.Log.Error(err, "Can not delete OpenshiftApiServerOperatorServiceMonitor")
		}
		if err = r.deleteOpenshiftServiceMonitor(cr, utils.OpenshiftClusterVersionOperatorServiceMonitorAsset,
			utils.OpenshiftClusterVersionOperatorServiceMonitorAsset,
			utils.OpenshiftClusterVersionOperator); err != nil {
			r.Log.Error(err, "Can not delete OpenshiftClusterVersionOperatorServiceMonitor")
		}
		if err = r.deleteOpenshiftServiceMonitor(cr, utils.OpenshiftStatemetricsServiceMonitorAsset,
			utils.OpenshiftStatemetricsServiceMonitorName,
			utils.OpenshiftStatemetrics); err != nil {
			r.Log.Error(err, "Can not delete OpenshiftStatemetricsServiceMonitor")
		}
		if err = r.deleteOpenshiftServiceMonitor(cr, utils.OpenshiftHAProxyServiceMonitorAsset,
			utils.OpenshiftHAProxyServiceMonitorName,
			utils.OpenshiftHAProxy); err != nil {
			r.Log.Error(err, "Can not delete OpenshiftHAProxyServiceMonitor")
		}
	}
}

// IsMonitorPresentInPublicCloud gets an answer in format (does the public cloud affect the monitor?, should the monitor be installed?)
func IsMonitorPresentInPublicCloud(cr *v1beta1.PlatformMonitoring, monitorName string) (bool, bool) {
	if cr.Spec.PublicCloudName != "" {
		if pcMap, ok := utils.PublicCloudMonitorsEnabled[cr.Spec.PublicCloudName]; ok {
			if installed, included := pcMap[monitorName]; included {
				return true, installed
			}
		}
	}
	return false, false
}

// IsMonitorInstall returns "true" if monitor should be installed
func IsMonitorInstall(cr *v1beta1.PlatformMonitoring, monitorName string) bool {
	// If the monitor affected by installation in the public cloud, cr.Spec.KubernetesMonitors doesn't matter
	affected, installed := IsMonitorPresentInPublicCloud(cr, monitorName)
	if affected {
		return installed
	} else {
		monitor, ok := cr.Spec.KubernetesMonitors[monitorName]
		return ok && monitor.IsInstall()
	}
}
