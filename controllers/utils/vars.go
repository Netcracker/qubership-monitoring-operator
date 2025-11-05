package utils

import (
	"path/filepath"
	"runtime"
)

const (
	OAuthProxySecretDir      = "/etc/oauth-proxy"
	OAuthProxyCfg            = OAuthProxySecretDir + "/oauth2-proxy.cfg"
	OAuthProxySecret         = "oauth2-proxy-config"
	SecretNamePrefix         = "secret-"
	OAuthProxySecretName     = SecretNamePrefix + OAuthProxySecret
	TlsCertificatesSecretDir = "/etc/oauth-proxy/certificates"
	OAuthProxyName           = "oauth-proxy"
	OAuthPort                = 9092
	ExternalURLPrefix        = "http://"

	AlertmanagerServiceName           = "alertmanager-operated"
	AlertmanagerServicePort           = 9093
	AlertmanagerOAuthProxyServiceName = "alertmanager-oauth2-proxy"
	OAuthProxyServicePortName         = "oauth-proxy"

	StackdriverPrometheusSidecarName = "stackdriver-prometheus"
	PrometheusServiceName            = "prometheus-operated"
	PrometheusServicePort            = 9090
	PrometheusOAuthProxyServiceName  = "prometheus-oauth2-proxy"

	GrafanaServiceName     = "grafana-service"
	GrafanaServicePort     = 3000
	GrafanaExtraVarsSecret = "grafana-extra-vars-secret"

	ClickHouseServiceName = "clickhouse-cluster"
	ClickHouseSecret      = "clickhouse-operator-credentials"

	VmAuthOAuthProxyServiceName = "vmauth-oauth2-proxy"
	VmAuthServicePort           = 8427

	ScrapeResources = "configmaps,cronjobs,daemonsets,deployments,endpointslices,jobs,limitranges,persistentvolumeclaims,poddisruptionbudgets,namespaces,nodes,pods,persistentvolumes,replicasets,replicationcontrollers,resourcequotas,services,statefulsets"

	VmComponentName       = "k8s"
	VmSingleComponentName = "vmsingle"
	VmSingleServiceName   = "vmsingle-k8s"
	VmSingleServicePort   = 8429

	VmAgentComponentName = "vmagent"
	VmAgentServiceName   = "vmagent-k8s"
	VmAgentServicePort   = 8429

	VmAlertManagerComponentName = "vmalertmanager"
	VmAlertManagerServiceName   = "vmalertmanager-k8s"
	VmAlertManagerServicePort   = 9093

	VmAlertComponentName   = "vmalert"
	VmAlertServiceName     = "vmalert-k8s"
	VmAlertServicePort     = 8080
	VmAlertAppRootEndpoint = "/vmalert"

	VmAuthComponentName = "vmauth"
	VmAuthServiceName   = "vmauth-k8s"

	VmSelectServiceName = "vmselect-k8s"
	VmSelectServicePort = 8481

	VmClusterComponentName = "vmcluster"
	VmSelectComponentName  = "vmselect"

	NginxIngressAppRootAnnotation = "nginx.ingress.kubernetes.io/app-root"

	VmOperatorTLSSecret     = "vmoperator-tls-secret"
	VmAlertTLSSecret        = "vmalert-tls-secret"
	VmAlertManagerTLSSecret = "vmalertmanager-tls-secret"
	VmAgentTLSSecret        = "vmagent-tls-secret"
	VmSingleTLSSecret       = "vmsingle-tls-secret"
	VmAuthTLSSecret         = "vmauth-tls-secret"
	VmSelectTLSSecret       = "vmselect-tls-secret"
	VmInsertTLSSecret       = "vminsert-tls-secret"
	VmStorageTLSSecret      = "vmstorage-tls-secret"
)

var (
	DashboardsFolder = "grafana-dashboards/"
	BasePath         = "assets/"

	// PrometheusOperatorComponentName contains name of prometheus-operator pod
	PrometheusOperatorComponentName           = "prometheus-operator"
	PrometheusOperatorClusterRoleAsset        = BasePath + "cluster-role.yaml"
	PrometheusOperatorRoleAsset               = BasePath + "role.yaml"
	PrometheusOperatorClusterRoleBindingAsset = BasePath + "cluster-role-binding.yaml"
	PrometheusOperatorRoleBindingAsset        = BasePath + "role-binding.yaml"
	PrometheusOperatorServiceAccountAsset     = BasePath + "service-account.yaml"
	PrometheusOperatorDeploymentAsset         = BasePath + "deployment.yaml"
	PrometheusOperatorServiceAsset            = BasePath + "service.yaml"
	PrometheusOperatorPodMonitorAsset         = BasePath + "pod-monitor.yaml"

	// PrometheusComponentName contains name of prometheus pod
	PrometheusComponentName           = "prometheus"
	PrometheusAsset                   = BasePath + "prometheus.yaml"
	PrometheusServiceAccountAsset     = BasePath + "service-account.yaml"
	PrometheusClusterRoleAsset        = BasePath + "cluster-role.yaml"
	PrometheusClusterRoleBindingAsset = BasePath + "cluster-role-binding.yaml"
	PrometheusIngressAsset            = BasePath + "ingress.yaml"
	PrometheusPodMonitorAsset         = BasePath + "pod-monitor.yaml"

	PrometheusRulesAsset = "assets/prometheus-rules.yaml"

	// VmOperatorComponentName contains name of victoriametrics-operator pod
	VmOperatorComponentName                   = "victoriametrics-operator"
	VmKubeletName                             = "kubelet"
	VmOperatorClusterRoleAsset                = BasePath + "cluster-role.yaml"
	VmOperatorRoleAsset                       = BasePath + "role.yaml"
	VmOperatorClusterRoleBindingAsset         = BasePath + "cluster-role-binding.yaml"
	VmOperatorRoleBindingAsset                = BasePath + "role-binding.yaml"
	VmOperatorServiceAccountAsset             = BasePath + "service-account.yaml"
	VmOperatorDeploymentAsset                 = BasePath + "deployment.yaml"
	VmOperatorServiceAsset                    = BasePath + "service.yaml"
	VmKubeletServiceAsset                     = BasePath + "kubelet-service.yaml"
	VmKubeletServiceEndpointsAsset            = BasePath + "kubelet-endpoints.yaml"
	VmOperatorServiceMonitorAsset             = BasePath + "service-monitor.yaml"
	VmOperatorPodSecurityPolicyAsset          = BasePath + "podsecuritypolicy.yaml"
	VmOperatorSecurityContextConstraintsAsset = BasePath + "securitycontextconstraints.yaml"

	VmSingleAsset                   = BasePath + "vmsingle.yaml"
	VmSingleIngressAsset            = BasePath + "ingress.yaml"
	VmSingleServiceAccountAsset     = BasePath + "service-account.yaml"
	VmSingleClusterRoleAsset        = BasePath + "cluster-role.yaml"
	VmSingleClusterRoleBindingAsset = BasePath + "cluster-role-binding.yaml"

	VmAgentAsset                   = BasePath + "vmagent.yaml"
	VmAgentIngressAsset            = BasePath + "ingress.yaml"
	VmAgentServiceAccountAsset     = BasePath + "service-account.yaml"
	VmAgentClusterRoleAsset        = BasePath + "cluster-role.yaml"
	VmAgentClusterRoleBindingAsset = BasePath + "cluster-role-binding.yaml"
	VmAgentRoleAsset               = BasePath + "role.yaml"
	VmAgentRoleBindingAsset        = BasePath + "role-binding.yaml"

	VmAlertManagerAsset                   = BasePath + "vmalertmanager.yaml"
	VmAlertManagerIngressAsset            = BasePath + "ingress.yaml"
	VmAlertManagerSecretAsset             = BasePath + "secret.yaml"
	VmAlertManagerServiceAccountAsset     = BasePath + "service-account.yaml"
	VmAlertManagerClusterRoleAsset        = BasePath + "cluster-role.yaml"
	VmAlertManagerClusterRoleBindingAsset = BasePath + "cluster-role-binding.yaml"

	VmAlertAsset                   = BasePath + "vmalert.yaml"
	VmAlertIngressAsset            = BasePath + "ingress.yaml"
	VmAlertServiceAccountAsset     = BasePath + "service-account.yaml"
	VmAlertClusterRoleAsset        = BasePath + "cluster-role.yaml"
	VmAlertClusterRoleBindingAsset = BasePath + "cluster-role-binding.yaml"

	VmAuthAsset                   = BasePath + "vmauth.yaml"
	VmAuthIngressAsset            = BasePath + "ingress.yaml"
	VmAuthServiceAccountAsset     = BasePath + "service-account.yaml"
	VmAuthClusterRoleAsset        = BasePath + "cluster-role.yaml"
	VmAuthClusterRoleBindingAsset = BasePath + "cluster-role-binding.yaml"

	VmClusterAsset                   = BasePath + "vmcluster.yaml"
	VmClusterServiceAccountAsset     = BasePath + "service-account.yaml"
	VmClusterClusterRoleAsset        = BasePath + "cluster-role.yaml"
	VmClusterClusterRoleBindingAsset = BasePath + "cluster-role-binding.yaml"
	VmSelectIngressAsset             = BasePath + "vmselect-ingress.yaml"

	VmUserAsset = BasePath + "vmuser.yaml"

	// AlertManagerComponentName contains name of alertmanager pod
	AlertManagerComponentName       = "alertmanager"
	AlertManagerGroupName           = "AlertManager"
	AlertManagerAsset               = BasePath + "alertmanager.yaml"
	AlertManagerServiceAccountAsset = BasePath + "service-account.yaml"
	AlertManagerServiceAsset        = BasePath + "service.yaml"
	AlertManagerSecretAsset         = BasePath + "secret.yaml"
	AlertManagerIngressAsset        = BasePath + "ingress.yaml"
	AlertManagerPodMonitorAsset     = BasePath + "pod-monitor.yaml"

	OpenshiftApiServerServiceMonitorAsset              = BasePath + "service-monitor-openshift-apiserver.yaml"
	OpenshiftApiServerOperatorServiceMonitorAsset      = BasePath + "service-monitor-openshift-apiserver-operator.yaml"
	OpenshiftClusterVersionOperatorServiceMonitorAsset = BasePath + "service-monitor-openshift-cluster-version-operator.yaml"
	OpenshiftStatemetricsServiceMonitorAsset           = BasePath + "service-monitor-openshift-state-metrics.yaml"
	OpenshiftHAProxyServiceMonitorAsset                = BasePath + "service-monitor-openshift-haproxy.yaml"

	OpenshiftApiserver              = "openshift-apiserver-service-monitor"
	OpenshiftApiServerOperator      = "openshift-apiserver-operator-service-monitor"
	OpenshiftClusterVersionOperator = "openshift-cluster-version-operator-service-monitor"
	OpenshiftStatemetrics           = "openshift-state-metrics-service-monitor"
	OpenshiftHAProxy                = "openshift-haproxy-service-monitor"

	ApiServerServiceMonitorAsset             = BasePath + "service-monitor-apiserver.yaml"
	KubeControllerManagerServiceMonitorAsset = BasePath + "service-monitor-kube-controller-manager.yaml"
	KubeSchedulerServiceMonitorAsset         = BasePath + "service-monitor-kube-scheduler.yaml"
	KubeletServiceMonitorAsset               = BasePath + "service-monitor-kubelet.yaml"
	CoreDnsServiceMonitorAssetK8s            = BasePath + "service-monitor-core-dns-k8s.yaml"
	CoreDnsServiceMonitorAssetOs4            = BasePath + "service-monitor-core-dns-os4.yaml"
	NginxIngressPodMonitorAsset              = BasePath + "pod-monitor-nginx-ingress.yaml"

	EtcdServiceMonitorAsset             = BasePath + "service-monitor-etcd.yaml"
	EtcdClientCertsSecretAsset          = BasePath + "kube-etcd-client-certs-secret.yaml"
	EtcdServiceComponentAsset           = BasePath + "service-etcd.yaml"
	EtcdNodeExporterServiceAccountAsset = BasePath + "service-account.yaml"

	EtcdServiceComponentName                   = "etcd"
	EtcdServiceComponentNamespace              = "kube-system"
	EtcdPodLabelSelector                       = "component=etcd"
	EtcdServiceComponentNamespaceOpenshiftV4   = "openshift-etcd"
	EtcdCertificatesSourceNamespaceOpenshiftV4 = "openshift-etcd-operator"
	EtcdCertificatesSourceConfigmapOpenshiftV4 = "etcd-metric-serving-ca"
	EtcdCertificatesSourceSecretOpenshiftV4    = "etcd-metric-client"

	// EtcdServiceMonitorName contains the name of etcd service monitor
	EtcdServiceMonitorName = "etcdServiceMonitor"
	// ApiserverServiceMonitorName contains the name of k8s apiserver service monitor
	ApiserverServiceMonitorName = "apiserverServiceMonitor"
	// KubeControllerManagerServiceMonitorName contains the name of  kube controller manager service monitor
	KubeControllerManagerServiceMonitorName = "kubeControllerManagerServiceMonitor"
	// KubeSchedulerServiceMonitorName contains the name of  kube scheduler service monitor
	KubeSchedulerServiceMonitorName = "kubeSchedulerServiceMonitor"
	// KubeletServiceMonitorName contains the name of kubelet service monitor
	KubeletServiceMonitorName = "kubeletServiceMonitor"
	// KubeEtcdClientCertsSecretName contains the name of etcd client certificate
	KubeEtcdClientCertsSecretName = "kube-etcd-client-certs"
	// CoreDnsServiceMonitorName contains the name of coreDNS service monitor
	CoreDnsServiceMonitorName = "coreDnsServiceMonitor"
	// NginxIngressPodMonitorName contains the name of nginx-ingress pod monitor
	NginxIngressPodMonitorName = "nginxIngressPodMonitor"
	// OcpApiServerServiceMonitorName contains the name of service monitor for the openshift apiserver
	OpenshiftApiServerServiceMonitorName = "openshiftApiserverServiceMonitor"
	// OcpApiServerOperatorServiceMonitorName contains the name of service monitor for openshift apiserver operator
	OpenshiftApiServerOperatorServiceMonitorName = "openshiftApiserverOperatorServiceMonitor"
	// OcpClusterVersionOperatorServiceMonitorName contains the name of service monitor for openshift cluster version operator
	OpenshiftClusterVersionOperatorServiceMonitorName = "openshiftClusterVersionOperatorServiceMonitor"
	// OcpStatemetricsServiceMonitorName contains the name of service monitor for openshift state metrics exporter
	OpenshiftStatemetricsServiceMonitorName = "openshiftStatemetricsServiceMonitor"
	// OpenshiftHAProxyServiceMonitorName contains the name of service monitor for openshift haproxy
	OpenshiftHAProxyServiceMonitorName = "openshiftHAProxyServiceMonitor"

	// NodeExporterComponentName contains name of node-exporter pod
	NodeExporterComponentName                   = "node-exporter"
	NodeExporterMetricsPortName                 = "metrics"
	NodeExporterTextfileVolumeName              = "node-exporter-textfile"
	NodeExporterServiceAccountAsset             = BasePath + "service-account.yaml"
	NodeExporterClusterRoleAsset                = BasePath + "cluster-role.yaml"
	NodeExporterClusterRoleBindingAsset         = BasePath + "cluster-role-binding.yaml"
	NodeExporterServiceAsset                    = BasePath + "service.yaml"
	NodeExporterDaemonSetAsset                  = BasePath + "daemonset.yaml"
	NodeExporterServiceMonitorAsset             = BasePath + "service-monitor.yaml"
	NodeExporterSecurityContextConstraintsAsset = BasePath + "securitycontextconstraints.yaml"
	NodeExporterPodSecurityPolicyAsset          = BasePath + "podsecuritypolicy.yaml"

	// KubestatemetricsComponentName contains name of kube-state-metrics pod
	KubestatemetricsComponentName           = "kube-state-metrics"
	KubestatemetricsClusterRoleAsset        = BasePath + "cluster-role.yaml"
	KubestatemetricsClusterRoleBindingAsset = BasePath + "cluster-role-binding.yaml"
	KubestatemetricsServiceAccountAsset     = BasePath + "service-account.yaml"
	KubestatemetricsDeploymentAsset         = BasePath + "deployment.yaml"
	KubestatemetricsServiceAsset            = BasePath + "service.yaml"
	KubestatemetricsServiceMonitorAsset     = BasePath + "service-monitor.yaml"

	// GrafanaOperatorComponentName contains name of alertmanager pod
	GrafanaOperatorComponentName           = "grafana-operator"
	GrafanaOperatorClusterRoleAsset        = BasePath + GrafanaOperatorComponentName + "/cluster-role.yaml"
	GrafanaOperatorRoleAsset               = BasePath + GrafanaOperatorComponentName + "/role.yaml"
	GrafanaOperatorClusterRoleBindingAsset = BasePath + GrafanaOperatorComponentName + "/cluster-role-binding.yaml"
	GrafanaOperatorRoleBindingAsset        = BasePath + GrafanaOperatorComponentName + "/role-binding.yaml"
	GrafanaOperatorServiceAccountAsset     = BasePath + GrafanaOperatorComponentName + "/service-account.yaml"
	GrafanaOperatorDeploymentAsset         = BasePath + GrafanaOperatorComponentName + "/deployment.yaml"
	GrafanaOperatorPodMonitorAsset         = BasePath + GrafanaOperatorComponentName + "/pod-monitor.yaml"

	// GrafanaComponentName contains name of alertmanager pod
	GrafanaComponentName   = "grafana"
	GrafanaAsset           = BasePath + "grafana.yaml"
	GrafanaDataSourceAsset = BasePath + "grafana-datasource.yaml"
	GrafanaIngressAsset    = BasePath + "ingress.yaml"
	GrafanaPodMonitorAsset = BasePath + "pod-monitor.yaml"

	// JaegerServiceLabels contains labels for Jaeger Service label selector
	JaegerServiceLabels = map[string]string{
		"app":                         "jaeger",
		"app.kubernetes.io/component": "query",
		"app.kubernetes.io/part-of":   "jaeger",
	}

	// PushgatewayComponentName contains name of pushgateway pod
	PushgatewayComponentName           = "pushgateway"
	PushgatewayPortName                = "http"
	PushgatewayStorageVolumeName       = "storage-volume"
	PushgatewayPVCVolumeMountMountPath = "/data"
	PushgatewayPersistenceFile         = "pushgateway.data"
	PushgatewayPersistenceInterval     = "5m"
	PushgatewayDeploymentAsset         = BasePath + "deployment.yaml"
	PushgatewayServiceAsset            = BasePath + "service.yaml"
	PushgatewayPVCAsset                = BasePath + "pvc.yaml"
	PushgatewayIngressAsset            = BasePath + "ingress.yaml"
	PushgatewayServiceMonitorAsset     = BasePath + "service-monitor.yaml"

	// GrafanaKubernetesDashboardsResources is a list of common dashboards which will be work with default installation
	GrafanaKubernetesDashboardsResources = []string{
		"alerts-overview.yaml",
		"core-dns-dashboard.yaml",
		"etcd-dashboard.yaml",
		"govm-processes.yaml",
		"ingress-list-of-ingresses.yaml",
		"ingress-nginx-controller.yaml",
		"ingress-request-handling-performance.yaml",
		"jvm-processes.yaml",
		"home-dashboard.yaml",
		"kubernetes-cluster-overview.yaml",
		"kubernetes-kubelet.yaml",
		"kubernetes-apiserver.yaml",
		"kubernetes-distribution-by-labels.yaml",
		"kubernetes-namespace-resources.yaml",
		"kubernetes-nodes-resources.yaml",
		"kubernetes-pod-resources.yaml",
		"kubernetes-pods-distribution-by-node.yaml",
		"kubernetes-pods-distribution-by-zone.yaml",
		"kubernetes-top-resources.yaml",
		"node-details.yaml",
		"overall-platform-health.yaml",
		"prometheus-cardinality-explorer.yaml",
		"prometheus-self-monitoring.yaml",
		"victoriametrics-vmsingle.yaml",
		"victoriametrics-vmoperator.yaml",
		"victoriametrics-vmalert.yaml",
		"victoriametrics-vmagent.yaml",
		"operators-overview.yaml",
		"grafana-overview.yaml",
		"alertmanager-overview.yaml",
		"openshift-state-metrics.yaml",
		"openshift-apiserver.yaml",
		"openshift-cluster-version-operator.yaml",
		"openshift-haproxy.yaml",
		"tls-status.yaml",
		"ha-services.yaml",
	}

	// GrafanaNodeExporterDashboardResource should be installed only if node-exporter is installing
	GrafanaNodeExporterDashboardResource = "kubernetes-nodes-resources.yaml"
	// GrafanaHomeDashboardResource should be installed only if grafanaHomeDashboard parameter set to true
	GrafanaHomeDashboardResource = "home-dashboard.yaml"
	// GrafanaNginxRequestHandlingPerformanceDashboardResource should be installed only if
	// NginxIngressPodMonitor is installing and should NOT be installed in OpenShift 3.11
	GrafanaNginxRequestHandlingPerformanceDashboardResource = "ingress-request-handling-performance.yaml"
	// GrafanaNginxIngressDashboardResource should be installed only if
	// NginxIngressPodMonitor is installing and should NOT be installed in OpenShift 3.11
	GrafanaNginxIngressDashboardResource = "ingress-nginx-controller.yaml"
	// GrafanaNginxIngressListOfIngresses should be installed only if
	// NginxIngressPodMonitor is installing and should NOT be installed in OpenShift 3.11
	GrafanaNginxIngressListOfIngresses = "ingress-list-of-ingresses.yaml"
	// GrafanaCoreDnsDashboardResource should be installed only if CoreDnsServiceMonitor is installing
	// and should NOT be installed in public clouds
	GrafanaCoreDnsDashboardResource = "core-dns-dashboard.yaml"
	// GrafanaSelfMonitoringDashboardResource should be installed only if Grafana is installing
	GrafanaSelfMonitoringDashboardResource = "grafana-overview.yaml"
	// GrafanaAlertmanagerDashboardResource should be installed only if Grafana is installing
	GrafanaAlertmanagerDashboardResource = "alertmanager-overview.yaml"

	// DashboardsUIDsMap contains human-readable UIDs for Grafana dashboards as VALUES.
	// The controller get UIDs by key from this map and concatenates it with the current Namespace,
	// then put the result as a UID to the dashboard by Go templates.
	// NOTE: The uid can have a maximum length of 40 characters. So uid must be as short as possible.
	// Ref to this limit:
	// https://grafana.com/docs/grafana/latest/developers/http_api/dashboard/#identifier-id-vs-unique-identifier-uid
	DashboardsUIDsMap = map[string]string{
		"alertmanager-overview":                "alertmanager-overview",
		"alerts-overview":                      "alerts-overview",
		"core-dns-dashboard":                   "core-dns",
		"etcd-dashboard":                       "etcd",
		"govm-processes":                       "govm-processes",
		"grafana-overview":                     "grafana-overview",
		"ha-services":                          "ha-services",
		"home-dashboard":                       "home-dashboard",
		"ingress-list-of-ingresses":            "ing-list-of-ingresses",
		"ingress-nginx-controller":             "ing-nginx-controller",
		"ingress-request-handling-performance": "ing-req-handl-perform",
		"jvm-processes":                        "jvm-processes",
		"kubernetes-apiserver":                 "k8s-apiserver",
		"kubernetes-cluster-overview":          "k8s-cluster-overview",
		"kubernetes-distribution-by-labels":    "k8s-distr-by-labels",
		"kubernetes-kubelet":                   "k8s-kubelet",
		"kubernetes-namespace-resources":       "k8s-namespace-resources",
		"kubernetes-nodes-resources":           "k8s-nodes-resources",
		"kubernetes-pod-resources":             "k8s-pod-resources",
		"kubernetes-pods-distribution-by-node": "k8s-pods-distr-by-node",
		"kubernetes-pods-distribution-by-zone": "k8s-pods-distr-by-zone",
		"kubernetes-top-resources":             "k8s-top-resources",
		"node-details":                         "node-details",
		"openshift-apiserver":                  "os-apiserver",
		"openshift-cluster-version-operator":   "os-cluster-version-operator",
		"openshift-state-metrics":              "os-state-metrics",
		"openshift-haproxy":                    "os-haproxy",
		"operators-overview":                   "operators-overview",
		"overall-platform-health":              "overall-platform-health",
		"prometheus-cardinality-explorer":      "prom-cardinality",
		"prometheus-self-monitoring":           "prom-self-monitoring",
		"tls-status":                           "tls-status",
		"victoriametrics-vmagent":              "vm-vmagent",
		"victoriametrics-vmalert":              "vm-vmalert",
		"victoriametrics-vmoperator":           "vm-vmoperator",
		"victoriametrics-vmsingle":             "vm-vmsingle",
	}
	// DashboardTemplateLeftDelim defines left delimiter for Grafana dashboards to avoid using {{
	DashboardTemplateLeftDelim = "{%"
	// DashboardTemplateRightDelim defines right delimiter for Grafana dashboards to avoid using }}
	DashboardTemplateRightDelim = "%}"

	// PublicCloudDashboardsEnabled is a map[public-cloud-name]map[dashboard-name]enabling/disabling
	// that configures enabling and disabling dashboards for specific public clouds.
	// Set bool value to "true" if this dashboard must be always enabled in this public cloud.
	// Set bool value to "false" if this dashboard must be always disabled in this public cloud.
	PublicCloudDashboardsEnabled = map[string]map[string]bool{
		"aws": {
			"core-dns-dashboard": false,
			"etcd-dashboard":     false,
		},
		"azure": {
			"core-dns-dashboard": false,
			"etcd-dashboard":     false,
		},
		"google": {
			"core-dns-dashboard": false,
			"etcd-dashboard":     false,
		},
	}
	// PublicCloudRulesEnabled is a map[public-cloud-name]map[rule-group]enabling/disabling
	// that configures enabling and disabling rule groups for specific public clouds.
	// Set bool value to "true" if this group must be always enabled in this public cloud.
	// Set bool value to "false" if this group must be always disabled in this public cloud.
	PublicCloudRulesEnabled = map[string]map[string]bool{
		"aws": {
			"CoreDnsAlerts": false,
			"Etcd":          false,
		},
		"azure": {
			"CoreDnsAlerts": false,
			"Etcd":          false,
		},
		"google": {
			"CoreDnsAlerts": false,
			"Etcd":          false,
		},
	}
	// PublicCloudMonitorsEnabled is a map[public-cloud-name]map[monitor-name]enabling/disabling
	// that configures enabling and disabling service and pod monitors for specific public clouds.
	// Set bool value to "true" if this monitor must be always enabled in this public cloud.
	// Set bool value to "false" if this monitor must be always disabled in this public cloud.
	PublicCloudMonitorsEnabled = map[string]map[string]bool{
		"aws": {
			CoreDnsServiceMonitorName: false,
			EtcdServiceMonitorName:    false,
		},
		"azure": {
			CoreDnsServiceMonitorName: false,
			EtcdServiceMonitorName:    false,
		},
		"google": {
			CoreDnsServiceMonitorName: false,
			EtcdServiceMonitorName:    false,
		},
	}

	// PrivilegedRights indicates is extended privileges should be used for the monitoring components.
	// If set to true, creates ClusterRole resources for services which needs it.
	// If set to false, creates Role resources where it is possible and expects that ClusterRole resources
	// were created manually. Also reconfigure some components to use RoleBinding resources in other namespaces for
	// access to necessary custom resources.
	PrivilegedRights bool

	// Root folder of the project
	_, b, _, _ = runtime.Caller(0)
	RootDir    = filepath.Join(filepath.Dir(b), "../../..")
)
