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

package v1alpha1

import (
	vmetricsv1b1 "github.com/VictoriaMetrics/operator/api/operator/v1beta1"
	grafv1alpha1 "github.com/grafana-operator/grafana-operator/v4/api/integreatly/v1alpha1"
	promv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"
)

var (
	defaultAlertManagerImage = "prom/alertmanager:v0.19.0"

	defaultPrometheusOperatorImage       = "coreos/prometheus-operator:v0.34.0"
	defaultPrometheusConfigReloaderImage = "coreos/prometheus-config-reloader:v0.34.0"

	defaultPrometheusImage = "prom/prometheus:v2.1.0"

	defaultGrafanaOperatorImage              = "integreatly/grafana-operator:v3.1.0"
	defaultGrafanaOperatorInitContainerImage = "integreatly/grafana_plugins_init:0.0.2"

	defaultGrafanaImage = "grafana/grafana:10.4.10"
)

// PlatformMonitoringSpec defines the desired state of PlatformMonitoring
type PlatformMonitoringSpec struct {
	// Important: Run "make generate" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	AlertManager       *AlertManager      `json:"alertManager,omitempty"`
	KubeStateMetrics   *KubeStateMetrics  `json:"kubeStateMetrics,omitempty"`
	Prometheus         *Prometheus        `json:"prometheus,omitempty"`
	NodeExporter       *NodeExporter      `json:"nodeExporter,omitempty"`
	Grafana            *Grafana           `json:"grafana,omitempty"`
	Integration        *Integration       `json:"integration,omitempty"`
	Auth               *Auth              `json:"auth,omitempty"`
	OAuthProxy         *OAuthProxy        `json:"oAuthProxy,omitempty"`
	KubernetesMonitors map[string]Monitor `json:"kubernetesMonitors,omitempty"`
	GrafanaDashboards  *GrafanaDashboards `json:"grafanaDashboards,omitempty"`
	PrometheusRules    *PrometheusRules   `json:"prometheusRules,omitempty"`
	Promxy             *Promxy            `json:"promxy,omitempty"`
	Pushgateway        *Pushgateway       `json:"pushgateway,omitempty"`
	PublicCloudName    string             `json:"publicCloudName,omitempty"`
	Victoriametrics    *Victoriametrics   `json:"victoriametrics,omitempty"`
}

// AlertManager defines the desired state for some part of prometheus-operator deployment
type AlertManager struct {
	// Image to use for a `AlertManager` deployment.
	// The `AlertManager` is alerting system which read metrics from Prometheus
	// More info: https://prometheus.io/docs/alerting/alertmanager/
	Image string `json:"image"`
	// Port for alertManager service
	Port int32 `json:"port,omitempty"`
	// Install indicates is AlertManager will be installed.
	// Can be changed for already deployed service and the service
	// will be removed during next reconciliation iteration
	Install *bool `json:"install,omitempty"`
	// SecurityContext holds pod-level security attributes.
	SecurityContext *SecurityContext `json:"securityContext,omitempty"`
	// Resources defines resources requests and limits for single Pods.
	Resources v1.ResourceRequirements `json:"resources,omitempty"`
	// Containers allows injecting additional containers or modifying operator generated containers.
	// This can be used to allow adding an authentication proxy to a Prometheus pod or to change
	// the behavior of an operator generated container.
	Containers []v1.Container `json:"containers,omitempty"`
	// Ingress allows to create Ingress for AlertManager UI.
	Ingress *Ingress `json:"ingress,omitempty"`
	// Set paused to reconsilation
	Paused bool `json:"paused,omitempty"`
	// Set replicas
	Replicas *int32 `json:"replicas,omitempty"`
	// Tolerations allow the pods to schedule onto nodes with matching taints.
	Tolerations []v1.Toleration `json:"tolerations,omitempty"`
	// NodeSelector Define which Nodes the Pods are scheduled on.
	// Specified just as map[string]string. For example: "type: compute"
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// If specified, the pod's scheduling constraints.
	// More info: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core
	Affinity *v1.Affinity `json:"affinity,omitempty"`
	// Pod monitor for self monitoring
	PodMonitor *Monitor `json:"podMonitor,omitempty"`
	// Map of string keys and values that can be used to organize and categorize
	// (scope and select) objects. May match selectors of replication controllers
	// and services.
	// More info: https://kubernetes.io/docs/user-guide/labels
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
	// Annotations is an unstructured key value map stored with a resource that may be
	// set by external tools to store and retrieve arbitrary metadata. They are not
	// queryable and should be preserved when modifying objects.
	// More info: https://kubernetes.io/docs/user-guide/annotations
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// ServiceAccount is a structure which allow specify annotations and labels for Service Account
	// which will use by Alertmanager for work in Kubernetes. Cna be use by external tools to store
	// and retrieve arbitrary metadata.
	// +optional
	ServiceAccount *EmbeddedObjectMetadata `json:"serviceAccount,omitempty"`
	// PriorityClassName assigned to the Pods
	// +optional
	PriorityClassName string `json:"priorityClassName,omitempty"`
}

// EmbeddedObjectMetadata contains a subset of the fields included in k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta
// Only fields which are relevant to embedded resources are included.
type EmbeddedObjectMetadata struct {
	// Annotations is an unstructured key value map stored with a resource that may be set
	// by external tools to store and retrieve arbitrary metadata. They are not queryable and should be
	// preserved when modifying objects.
	// More info: http://kubernetes.io/docs/user-guide/annotations
	Annotations map[string]string `json:"annotations,omitempty"`
	// Map of string keys and values that can be used to organize and categorize (scope and select) objects.
	// May match selectors of replication controllers and services.
	// More info: http://kubernetes.io/docs/user-guide/labels
	Labels map[string]string `json:"labels,omitempty"`
}

// Grafana defines the desired state for some part of grafana-operator deployment
type Grafana struct {
	// Image to use for a `grafana` deployment.
	// The `grafana` is a web ui to show graphics.
	// More info: https://github.com/grafana/grafana
	Image string `json:"image"`
	// Operator parameters
	Operator GrafanaOperator `json:"operator"`
	// Install indicates is Grafana will be installed.
	// Can be changed for already deployed service and the service
	// will be removed during next reconciliation iteration
	Install *bool `json:"install,omitempty"`
	// SecurityContext holds pod-level security attributes.
	SecurityContext *SecurityContext `json:"securityContext,omitempty"`
	// Resources defines resources requests and limits for single Pods.
	Resources v1.ResourceRequirements `json:"resources,omitempty"`
	// Ingress allows to create Ingress for Grafana UI.
	Ingress *Ingress `json:"ingress,omitempty"`
	// Config allows to override Config for Grafana.
	Config grafv1alpha1.GrafanaConfig `json:"config,omitempty"`
	// Custom grafana home dashboard
	GrafanaHomeDashboard bool `json:"grafanaHomeDashboard,omitempty"`
	// DashboardLabelSelector allows to query over a set of resources according to labels
	DashboardLabelSelector []*metav1.LabelSelector `json:"dashboardLabelSelector,omitempty"`
	// DashboardNamespaceSelector allows to query over a set of resources in namespaces that fits label selector
	DashboardNamespaceSelector *metav1.LabelSelector `json:"dashboardNamespaceSelector,omitempty"`
	// Set paused to reconsilation
	Paused bool `json:"paused,omitempty"`
	// Tolerations allow the pods to schedule onto nodes with matching taints.
	Tolerations []v1.Toleration `json:"tolerations,omitempty"`
	// NodeSelector Define which Nodes the Pods are scheduled on.
	// Specified just as map[string]string. For example: "type: compute"
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// If specified, the pod's scheduling constraints.
	// More info: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core
	Affinity *v1.Affinity `json:"affinity,omitempty"`
	// Pod monitor for self monitoring
	PodMonitor *Monitor `json:"podMonitor,omitempty"`
	// DataStorage provides a means to configure the grafana data storage
	DataStorage *grafv1alpha1.GrafanaDataStorage `json:"dataStorage,omitempty"`
	// Map of string keys and values that can be used to organize and categorize
	// (scope and select) objects. May match selectors of replication controllers
	// and services.
	// More info: https://kubernetes.io/docs/user-guide/labels
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
	// Annotations is an unstructured key value map stored with a resource that may be
	// set by external tools to store and retrieve arbitrary metadata. They are not
	// queryable and should be preserved when modifying objects.
	// More info: https://kubernetes.io/docs/user-guide/annotations
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// Set replicas
	Replicas *int32 `json:"replicas,omitempty"`
	// ServiceAccount is a structure which allow specify annotations and labels for Service Account
	// which will use by Grafana for work in Kubernetes. Cna be use by external tools to store
	// and retrieve arbitrary metadata.
	// +optional
	ServiceAccount *EmbeddedObjectMetadata `json:"serviceAccount,omitempty"`
	// PriorityClassName assigned to the Pods
	// +optional
	PriorityClassName string `json:"priorityClassName,omitempty"`
}

// GrafanaOperator defines the desired state for some part of grafana-operator deployment
type GrafanaOperator struct {
	// Image to use for a `grafana-operator` deployment.
	// The `grafana-operator` is a control, deploy and process custom resources into Grafana entities.
	// More info: https://github.com/integr8ly/grafana-operator
	Image string `json:"image"`
	// Image to use to initialize Grafana deployment.
	InitContainerImage string `json:"initContainerImage"`
	// Namespaces to scope the interaction of the Grafana operator.
	Namespaces string `json:"namespaces,omitempty"`
	// SecurityContext holds pod-level security attributes.
	SecurityContext *SecurityContext `json:"securityContext,omitempty"`
	// Resources defines resources requests and limits for single Pods.
	Resources v1.ResourceRequirements `json:"resources,omitempty"`
	// Set paused to reconsilation
	Paused bool `json:"paused,omitempty"`
	// Tolerations allow the pods to schedule onto nodes with matching taints.
	Tolerations []v1.Toleration `json:"tolerations,omitempty"`
	// NodeSelector Define which Nodes the Pods are scheduled on.
	// Specified just as map[string]string. For example: "type: compute"
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// If specified, the pod's scheduling constraints.
	// More info: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core
	Affinity *v1.Affinity `json:"affinity,omitempty"`
	// Pod monitor for self monitoring
	PodMonitor *Monitor `json:"podMonitor,omitempty"`
	// Map of string keys and values that can be used to organize and categorize
	// (scope and select) objects. May match selectors of replication controllers
	// and services.
	// More info: https://kubernetes.io/docs/user-guide/labels
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
	// Annotations is an unstructured key value map stored with a resource that may be
	// set by external tools to store and retrieve arbitrary metadata. They are not
	// queryable and should be preserved when modifying objects.
	// More info: https://kubernetes.io/docs/user-guide/annotations
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// Zap log level (one of 'debug', 'info', 'error' or any integer value > 0) (default info)
	// More info: https://github.com/grafana-operator/grafana-operator/blob/master/documentation/deploy_grafana.md
	// +optional
	LogLevel string `json:"logLevel,omitempty"`
	// ServiceAccount is a structure which allow specify annotations and labels for Service Account
	// which will use by Alertmanager for work in Kubernetes. Cna be use by external tools to store
	// and retrieve arbitrary metadata.
	// +optional
	ServiceAccount *EmbeddedObjectMetadata `json:"serviceAccount,omitempty"`
	// PriorityClassName assigned to the Pods
	// +optional
	PriorityClassName string `json:"priorityClassName,omitempty"`
}

// Prometheus defines the link to PrometheusSpec objects from prometheus-operator
type Prometheus struct {
	// Image to use for a `prometheus` deployment.
	// The `prometheus` is a systems and service monitoring system.
	// It collects metrics from configured targets at given intervals.
	// More info: https://github.com/prometheus/prometheus
	Image string `json:"image"`
	// Operator parameters
	Operator PrometheusOperator `json:"operator"`
	// Install indicates is Prometheus will be installed.
	// Can be changed for already deployed service and the service
	// will be removed during next reconciliation iteration
	Install *bool `json:"install,omitempty"`
	// Image to use for a `prometheus-config-reloader`
	// The `prometheus-config-reloaded` is an add-on to prometheus that monitors changes in prometheus.yaml
	// and an HTTP request reloads the prometheus configuration.
	// More info: https://github.com/prometheus-operator/prometheus-operator/tree/master/cmd/prometheus-config-reloader
	ConfigReloaderImage string `json:"configReloaderImage,omitempty"`
	// RemoteWriteSpec defines the remote_write configuration for prometheus.
	// The `remote_write` allows transparently send samples to a long term storage.
	// More info: https://prometheus.io/docs/operating/integrations/#remote-endpoints-and-storage
	RemoteWrite []promv1.RemoteWriteSpec `json:"remoteWrite,omitempty"`
	// RemoteReadSpec defines the remote_read configuration for prometheus.
	// The `remote_read` allows transparently receive samples from a long term storage.
	// More info: https://prometheus.io/docs/operating/integrations/#remote-endpoints-and-storage
	RemoteRead []promv1.RemoteReadSpec `json:"remoteRead,omitempty"`
	// Secrets is a list of Secrets in the same namespace as the Prometheus
	// object, which shall be mounted into the Prometheus Pods.
	// The Secrets are mounted into /etc/prometheus/secrets/<secret-name>.
	Secrets []string `json:"secrets,omitempty"`
	// Define details regarding alerting.
	Alerting *promv1.AlertingSpec `json:"alerting,omitempty"`
	// The labels to add to any time series or alerts when communicating with
	// external systems (federation, remote storage, Alertmanager).
	ExternalLabels map[string]string `json:"externalLabels,omitempty"`
	// SecurityContext holds pod-level security attributes.
	SecurityContext *SecurityContext `json:"securityContext,omitempty"`
	// Define which Nodes the Pods are scheduled on.
	// Specified just as map[string]string. For example: "type: compute"
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// If specified, the pod's scheduling constraints.
	// More info: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.20/#affinity-v1-core
	Affinity *v1.Affinity `json:"affinity,omitempty"`
	// Resources defines resources requests and limits for single Pods.
	Resources v1.ResourceRequirements `json:"resources,omitempty"`
	// Storage spec to specify how storage shall be used.
	// More info: https://github.com/prometheus-operator/prometheus-operator/blob/master/Documentation/api.md#storagespec
	Storage *promv1.StorageSpec `json:"storage,omitempty"`
	// Volumes allows configuration of additional volumes on the output StatefulSet definition.
	// Volumes specified will be appended to other volumes that are generated as a result of StorageSpec objects.
	// More info: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.20/#volume-v1-core
	Volumes []v1.Volume `json:"volumes,omitempty"`
	// VolumeMounts allows configuration of additional VolumeMounts on the output StatefulSet definition.
	// VolumeMounts specified will be appended to other VolumeMounts in the prometheus container,
	// that are generated as a result of StorageSpec objects.
	// More info: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.20/#volumemount-v1-core
	VolumeMounts []v1.VolumeMount `json:"volumeMounts,omitempty"`
	// Ingress allows to create Ingress for Prometheus UI.
	Ingress *Ingress `json:"ingress,omitempty"`
	// Retention policy by time
	Retention string `json:"retention,omitempty"`
	// Retention policy by size [EXPERIMENTAL]
	RetentionSize string `json:"retentionsize,omitempty"`
	// Containers allows injecting additional containers or modifying operator generated containers.
	// This can be used to allow adding an authentication proxy to a Prometheus pod or to change
	// the behavior of an operator generated container.
	Containers []v1.Container `json:"containers,omitempty"`
	// The external URL the Prometheus instances will be available under. This is
	// necessary to generate correct URLs. This is necessary if Prometheus is not
	// served from root of a DNS name.
	ExternalURL string `json:"externalUrl,omitempty"`
	// Set paused to reconsilation
	Paused bool `json:"paused,omitempty"`
	// Set replicas
	Replicas *int32 `json:"replicas,omitempty"`
	// Tolerations allow the pods to schedule onto nodes with matching taints.
	Tolerations []v1.Toleration `json:"tolerations,omitempty"`
	// Pod monitor for self monitoring
	PodMonitor *Monitor `json:"podMonitor,omitempty"`
	// Namespace selector for rules
	RuleNamespaceSelector *metav1.LabelSelector `json:"ruleNamespaceSelector,omitempty"`
	// Namespace selector for PodMonitors
	PodMonitorNamespaceSelector *metav1.LabelSelector `json:"podMonitorNamespaceSelector,omitempty"`
	// Namespace selector for ServiceMonitors
	ServiceMonitorNamespaceSelector *metav1.LabelSelector `json:"serviceMonitorNamespaceSelector,omitempty"`
	// Selector for rules
	RuleSelector *metav1.LabelSelector `json:"ruleSelector,omitempty"`
	// Selector for PodMoniotors
	PodMonitorSelector *metav1.LabelSelector `json:"podMonitorSelector,omitempty"`
	// Selector for ServiceMonitors
	ServiceMonitorSelector *metav1.LabelSelector `json:"serviceMonitorSelector,omitempty"`
	// QuerySpec defines the query command line flags when starting Prometheus
	Query *promv1.QuerySpec `json:"query,omitempty"`
	// TLSConfig define TLS configuration for Prometheus.
	TLSConfig *PromTLSConfig `json:"tlsConfig,omitempty"`
	// Enable access to prometheus web admin API. Defaults to the value of false.
	// WARNING: Enabling the admin APIs enables mutating endpoints, to delete data, shutdown Prometheus, and more.
	// Enabling this should be done with care and the user is advised to add additional authentication authorization
	// via a proxy to ensure only clients authorized to perform these actions can do so.
	// For more information see https://prometheus.io/docs/prometheus/latest/querying/api/#tsdb-admin-apis
	EnableAdminAPI bool `json:"enableAdminAPI,omitempty"`
	// Map of string keys and values that can be used to organize and categorize
	// (scope and select) objects. May match selectors of replication controllers
	// and services.
	// More info: https://kubernetes.io/docs/user-guide/labels
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
	// Annotations is an unstructured key value map stored with a resource that may be
	// set by external tools to store and retrieve arbitrary metadata. They are not
	// queryable and should be preserved when modifying objects.
	// More info: https://kubernetes.io/docs/user-guide/annotations
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// ServiceAccount is a structure which allow specify annotations and labels for Service Account
	// which will use by Prometheus for work in Kubernetes. Cna be use by external tools to store
	// and retrieve arbitrary metadata.
	// +optional
	ServiceAccount *EmbeddedObjectMetadata `json:"serviceAccount,omitempty"`
	// Name of Prometheus external label used to denote replica name.
	// Defaults to the value of `prometheus_replica`. External label will
	// _not_ be added when value is set to empty string (`""`).
	ReplicaExternalLabelName *string `json:"replicaExternalLabelName,omitempty"`
	// Enable access to Prometheus disabled features. By default, no features are enabled. Enabling disabled features
	// is entirely outside the scope of what the maintainers will support and by doing so, you accept that this
	// behavior may break at any time without notice.
	// For more information see https://prometheus.io/docs/prometheus/latest/disabled_features/
	EnableFeatures []promv1.EnableFeature `json:"enableFeatures,omitempty"`
	// Interval between consecutive scrapes. Default: `30s`
	ScrapeInterval *string `json:"scrapeInterval,omitempty"`
	// Number of seconds to wait for target to respond before erroring. Default: `10s`
	ScrapeTimeout *string `json:"scrapeTimeout,omitempty"`
	// Interval between consecutive evaluations. Default: `30s`
	EvaluationInterval *string `json:"evaluationInterval,omitempty"`
	// PriorityClassName assigned to the Pods
	// +optional
	PriorityClassName string `json:"priorityClassName,omitempty"`
}

// KubeStateMetrics defines the desired state for some part of kube-state-metrics deployment
type KubeStateMetrics struct {
	// Install indicates is kube-state-metrics will be installed.
	// Can be changed for already deployed service and the service
	// will be removed during next reconciliation iteration
	Install *bool `json:"install,omitempty"`
	// Image to use for a `kube-state-metrics` deployment.
	// The `kube-state-metrics` is an exporter to collect Kubernetes metrics
	// More info: https://github.com/kubernetes/kube-state-metrics
	Image string `json:"image"`
	// Resources defines resources requests and limits for single Pods.
	Resources v1.ResourceRequirements `json:"resources,omitempty"`
	// SecurityContext holds pod-level security attributes.
	SecurityContext *SecurityContext `json:"securityContext,omitempty"`
	// List of comma-separated namespaces to scrape metrics in non-privileged mode.
	Namespaces string `json:"namespaces,omitempty"`
	//Comma-separated list of Resources to be enabled.
	ScrapeResources string `json:"scrapeResources,omitempty"`
	// Comma-separated list of additional Kubernetes label keys that will be used in the resource labels metric.
	MetricLabelsAllowlist string `json:"metricLabelsAllowlist,omitempty"`
	// Set paused to reconsilation.
	Paused bool `json:"paused,omitempty"`
	// Tolerations allow the pods to schedule onto nodes with matching taints.
	Tolerations []v1.Toleration `json:"tolerations,omitempty"`
	// NodeSelector Define which Nodes the Pods are scheduled on.
	// Specified just as map[string]string. For example: "type: compute"
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// If specified, the pod's scheduling constraints.
	// More info: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core
	Affinity *v1.Affinity `json:"affinity,omitempty"`
	// Service monitor for pulling metrics
	ServiceMonitor *Monitor `json:"serviceMonitor,omitempty"`
	// Map of string keys and values that can be used to organize and categorize
	// (scope and select) objects. May match selectors of replication controllers
	// and services.
	// More info: https://kubernetes.io/docs/user-guide/labels
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
	// Annotations is an unstructured key value map stored with a resource that may be
	// set by external tools to store and retrieve arbitrary metadata. They are not
	// queryable and should be preserved when modifying objects.
	// More info: https://kubernetes.io/docs/user-guide/annotations
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// ServiceAccount is a structure which allow specify annotations and labels for Service Account
	// which will use by Alertmanager for work in Kubernetes. Cna be use by external tools to store
	// and retrieve arbitrary metadata.
	// +optional
	ServiceAccount *EmbeddedObjectMetadata `json:"serviceAccount,omitempty"`
	// PriorityClassName assigned to the Pods
	// +optional
	PriorityClassName string `json:"priorityClassName,omitempty"`
}

// PrometheusOperator defines the desired state for some part of prometheus-operator deployment
type PrometheusOperator struct {
	// Image to use for a `prometheus-operator` deployment.
	// The `prometheus-operator` makes the Prometheus configuration Kubernetes native and manages and operates
	// Prometheus and Alertmanager clusters.
	// More info: https://github.com/prometheus-operator/prometheus-operator
	Image string `json:"image"`
	// Resources defines resources requests and limits for single Pods.
	Resources v1.ResourceRequirements `json:"resources,omitempty"`
	// SecurityContext holds pod-level security attributes.
	SecurityContext *SecurityContext `json:"securityContext,omitempty"`
	// Set paused to reconsilation
	Paused bool `json:"paused,omitempty"`
	// Tolerations allow the pods to schedule onto nodes with matching taints.
	Tolerations []v1.Toleration `json:"tolerations,omitempty"`
	// NodeSelector Define which Nodes the Pods are scheduled on.
	// Specified just as map[string]string. For example: "type: compute"
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// If specified, the pod's scheduling constraints.
	// More info: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core
	Affinity *v1.Affinity `json:"affinity,omitempty"`
	// Pod monitor for self monitoring
	PodMonitor *Monitor `json:"podMonitor,omitempty"`
	//Namespaces to scope the interaction of the Prometheus Operator and the apiserver.
	Namespaces string `json:"namespaces,omitempty"`
	// Map of string keys and values that can be used to organize and categorize
	// (scope and select) objects. May match selectors of replication controllers
	// and services.
	// More info: https://kubernetes.io/docs/user-guide/labels
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
	// Annotations is an unstructured key value map stored with a resource that may be
	// set by external tools to store and retrieve arbitrary metadata. They are not
	// queryable and should be preserved when modifying objects.
	// More info: https://kubernetes.io/docs/user-guide/annotations
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// ServiceAccount is a structure which allow specify annotations and labels for Service Account
	// which will use by Alertmanager for work in Kubernetes. Cna be use by external tools to store
	// and retrieve arbitrary metadata.
	// +optional
	ServiceAccount *EmbeddedObjectMetadata `json:"serviceAccount,omitempty"`
	// PriorityClassName assigned to the Pods
	// +optional
	PriorityClassName string `json:"priorityClassName,omitempty"`
}

type Victoriametrics struct {
	TLSEnabled bool `json:"tlsEnabled,omitempty"`
	//VmReplicas     *int32         `json:"vmReplicas,omitempty"`	//TODO: Revert this line when vmCluster is actualized
	VmReplicas     *int32         `json:"-"`
	VmOperator     VmOperator     `json:"vmOperator,omitempty"`
	VmSingle       VmSingle       `json:"vmSingle,omitempty"`
	VmAgent        VmAgent        `json:"vmAgent,omitempty"`
	VmAlertManager VmAlertManager `json:"vmAlertManager,omitempty"`
	VmAlert        VmAlert        `json:"vmAlert,omitempty"`
	VmAuth         VmAuth         `json:"vmAuth,omitempty"`
	VmUser         VmUser         `json:"vmUser,omitempty"`
	//VmCluster      VmCluster      `json:"vmCluster,omitempty"`	//TODO: Revert this line when vmCluster is actualized
	VmCluster VmCluster `json:"-"`
}
type VmOperator struct {
	// Install indicates is victoriametrics-operator will be installed.
	// Can be changed for already deployed service and the service
	// will be removed during next reconciliation iteration
	Install *bool `json:"install,omitempty"`
	//Number of replicas of victoriametrics-operator
	Replicas *int32 `json:"replicas,omitempty"`
	// Image to use for a `victoriametrics-operator` deployment.
	// The `victoriametrics-operator` makes the vmoperator configuration Kubernetes native and manages and operates
	// More info: https://github.com/VictoriaMetrics/operator
	Image string `json:"image"`
	// Resources defines resources requests and limits for single Pods.
	Resources v1.ResourceRequirements `json:"resources,omitempty"`
	// SecurityContext holds pod-level security attributes.
	SecurityContext *SecurityContext `json:"securityContext,omitempty"`
	// ContainerSecurityContext holds container-level security attributes.
	ContainerSecurityContext *v1.SecurityContext `json:"containerSecurityContext,omitempty"`
	// Set paused to reconcilation
	Paused bool `json:"paused,omitempty"`
	// ExtraEnvs that will be added to VMOperator pod
	// +optional
	ExtraEnvs []v1.EnvVar `json:"extraEnvs,omitempty"`
	// Tolerations allow the pods to schedule onto nodes with matching taints.
	Tolerations []v1.Toleration `json:"tolerations,omitempty"`
	// NodeSelector Define which Nodes the Pods are scheduled on.
	// Specified just as map[string]string. For example: "type: compute"
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// If specified, the pod's scheduling constraints.
	// More info: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core
	Affinity *v1.Affinity `json:"affinity,omitempty"`
	// Map of string keys and values that can be used to organize and categorize
	// (scope and select) objects. May match selectors of replication controllers
	// and services.
	// More info: https://kubernetes.io/docs/user-guide/labels
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
	// Annotations is an unstructured key value map stored with a resource that may be
	// set by external tools to store and retrieve arbitrary metadata. They are not
	// queryable and should be preserved when modifying objects.
	// More info: https://kubernetes.io/docs/user-guide/annotations
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// Service monitor for pulling metrics
	ServiceMonitor *Monitor `json:"serviceMonitor,omitempty"`
	// PriorityClassName assigned to the Pods
	// +optional
	PriorityClassName string `json:"priorityClassName,omitempty"`
	// TLS Configuration
	// +optional
	TLSConfig *VmTLSConfig `json:"tlsConfig,omitempty"`
}

type VmTLSConfig struct {
	SecretName string `json:"secretName,omitempty"`
}

type VmSingle struct {
	// Install indicates is vmsingle will be installed.
	// Can be changed for already deployed CR
	// will be removed during next reconciliation iteration
	Install *bool `json:"install,omitempty"`
	// Image to use for a `vmsingle` deployment.
	// The `vmsingle` makes the vmsingle configuration Kubernetes native and manages and operates
	// More info: https://docs.victoriametrics.com/Single-server-VictoriaMetrics.html
	Image string `json:"image"`
	// Resources defines resources requests and limits for single Pods.
	Resources v1.ResourceRequirements `json:"resources,omitempty"`
	// SecurityContext holds pod-level security attributes.
	SecurityContext *SecurityContext `json:"securityContext,omitempty"`
	// Set paused to reconcilation
	Paused bool `json:"paused,omitempty"`
	// Tolerations allow the pods to schedule onto nodes with matching taints.
	Tolerations []v1.Toleration `json:"tolerations,omitempty"`
	// NodeSelector Define which Nodes the Pods are scheduled on.
	// Specified just as map[string]string. For example: "type: compute"
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// If specified, the pod's scheduling constraints.
	// More info: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core
	Affinity *v1.Affinity `json:"affinity,omitempty"`
	// Map of string keys and values that can be used to organize and categorize
	// (scope and select) objects. May match selectors of replication controllers
	// and services.
	// More info: https://kubernetes.io/docs/user-guide/labels
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
	// Annotations is an unstructured key value map stored with a resource that may be
	// set by external tools to store and retrieve arbitrary metadata. They are not
	// queryable and should be preserved when modifying objects.
	// More info: https://kubernetes.io/docs/user-guide/annotations
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// Containers property allows to inject additions sidecars or to patch existing containers.
	// It can be useful for proxies, backup, etc.
	Containers []v1.Container `json:"containers,omitempty"`
	// Ingress allows to create Ingress for VmSingle UI.
	Ingress *Ingress `json:"ingress,omitempty"`
	// RetentionPeriod for the stored metrics
	// Note VictoriaMetrics has data/ and indexdb/ folders
	// metrics from data/ removed eventually as soon as partition leaves retention period
	// reverse index data at indexdb rotates once at the half of configured retention period
	// https://docs.victoriametrics.com/Single-server-VictoriaMetrics.html#retention
	RetentionPeriod string `json:"retentionPeriod"`
	// ExtraArgs that will be passed to  VMSingle pod
	// for example remoteWrite.tmpDataPath: /tmp
	// +optional
	ExtraArgs map[string]string `json:"extraArgs,omitempty"`
	// ExtraEnvs that will be added to VMSingle pod
	// +optional
	ExtraEnvs []v1.EnvVar `json:"extraEnvs,omitempty"`
	// Secrets is a list of Secrets in the same namespace as the VMSingle
	// object, which shall be mounted into the VMSingle Pods.
	// +optional
	Secrets []string `json:"secrets,omitempty"`
	// StorageDataPath disables spec.storage option and overrides arg for victoria-metrics binary --storageDataPath,
	// its users responsibility to mount proper device into given path.
	// + optional
	StorageDataPath string `json:"storageDataPath,omitempty"`
	// Storage is the definition of how storage will be used by the VMSingle
	// by default it`s empty dir
	// +optional
	Storage *v1.PersistentVolumeClaimSpec `json:"storage,omitempty"`
	// StorageMeta defines annotations and labels attached to PVC for given vmsingle CR
	// +optional
	StorageMetadata *vmetricsv1b1.EmbeddedObjectMetadata `json:"storageMetadata,omitempty"`
	// Volumes allows configuration of additional volumes on the output deploy definition.
	// Volumes specified will be appended to other volumes that are generated as a result of
	// StorageSpec objects.
	// +optional
	Volumes []v1.Volume `json:"volumes,omitempty"`
	// VolumeMounts allows configuration of additional VolumeMounts on the output deploy definition.
	// VolumeMounts specified will be appended to other VolumeMounts in the vmsingle container,
	// that are generated as a result of StorageSpec objects.
	// +optional
	VolumeMounts []v1.VolumeMount `json:"volumeMounts,omitempty"`
	// TerminationGracePeriodSeconds period for container graceful termination
	// +optional
	TerminationGracePeriodSeconds *int64 `json:"terminationGracePeriodSeconds,omitempty"`
	// PriorityClassName assigned to the Pods
	// +optional
	PriorityClassName string `json:"priorityClassName,omitempty"`
	// TLS Configuration
	// +optional
	TLSConfig *VmTLSConfig `json:"tlsConfig,omitempty"`
}

type VmAgent struct {
	// Install indicates is vmagent will be installed.
	// Can be changed for already deployed service and the service
	// will be removed during next reconciliation iteration
	Install *bool `json:"install,omitempty"`
	// Image to use for a `vmagent` deployment.
	// The `victoriametrics-operator` makes the VmAgent configuration Kubernetes native and manages and operates
	// More info: https://docs.victoriametrics.com/vmalert.html
	Image string `json:"image"`
	// ReplicaCount is the expected size of the VMAgent cluster. The controller will
	// eventually make the size of the running cluster equal to the expected
	// size.
	// NOTE enable VMSingle deduplication for replica usage
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Number of pods",xDescriptors="urn:alm:descriptor:com.tectonic.ui:podCount,urn:alm:descriptor:io.kubernetes:custom"
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`
	// Resources defines resources requests and limits for single Pods.
	Resources v1.ResourceRequirements `json:"resources,omitempty"`
	// SecurityContext holds pod-level security attributes.
	SecurityContext *SecurityContext `json:"securityContext,omitempty"`
	// Set paused to reconcilation
	Paused bool `json:"paused,omitempty"`
	// Containers property allows to inject additions sidecars or to patch existing containers.
	// It can be useful for proxies, backup, etc.
	Containers []v1.Container `json:"containers,omitempty"`
	// Tolerations allow the pods to schedule onto nodes with matching taints.
	Tolerations []v1.Toleration `json:"tolerations,omitempty"`
	// NodeSelector Define which Nodes the Pods are scheduled on.
	// Specified just as map[string]string. For example: "type: compute"
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// If specified, the pod's scheduling constraints.
	// More info: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core
	Affinity *v1.Affinity `json:"affinity,omitempty"`
	// Map of string keys and values that can be used to organize and categorize
	// (scope and select) objects. May match selectors of replication controllers
	// and services.
	// More info: https://kubernetes.io/docs/user-guide/labels
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
	// Annotations is an unstructured key value map stored with a resource that may be
	// set by external tools to store and retrieve arbitrary metadata. They are not
	// queryable and should be preserved when modifying objects.
	// More info: https://kubernetes.io/docs/user-guide/annotations
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// Ingress allows to create Ingress for VM UI.
	Ingress *Ingress `json:"ingress,omitempty"`
	// ScrapeInterval defines how often scrape targets by default
	// +optional
	// +kubebuilder:validation:Pattern:="[0-9]+(ms|s|m|h)"
	ScrapeInterval string `json:"scrapeInterval,omitempty"`
	//MaxScrapeInterval allows limiting maximum scrape interval for VMServiceScrape, VMPodScrape and other scrapes
	MaxScrapeInterval *string `json:"maxScrapeInterval,omitempty"`
	//MinScrapeInterval allows limiting minimal scrape interval for VMServiceScrape, VMPodScrape and other scrapes
	MinScrapeInterval *string `json:"minScrapeInterval,omitempty"`
	// VMAgentExternalLabelName Name of vmAgent external label used to denote VmAgent instance
	// name. Defaults to the value of `vmagent`. External label will
	// _not_ be added when value is set to empty string (`""`).
	// +optional
	VMAgentExternalLabelName *string `json:"vmAgentExternalLabelName,omitempty"`
	// ExternalLabels The labels to add to any time series scraped by vmagent.
	// it doesn't affect metrics ingested directly by push API's
	// +optional
	ExternalLabels map[string]string `json:"externalLabels,omitempty"`
	// RemoteWrite list of victoria metrics to some other remote write system
	// for vm it must look like: http://victoria-metrics-single:8429/api/v1/write
	// or for cluster different url
	// https://github.com/VictoriaMetrics/VictoriaMetrics/tree/master/app/vmagent#splitting-data-streams-among-multiple-systems
	// +optional
	RemoteWrite []vmetricsv1b1.VMAgentRemoteWriteSpec `json:"remoteWrite,omitempty"`
	// RemoteWriteSettings defines global settings for all remoteWrite urls.
	// + optional
	RemoteWriteSettings *vmetricsv1b1.VMAgentRemoteWriteSettings `json:"remoteWriteSettings,omitempty"`
	// RelabelConfig ConfigMap with global relabel config -remoteWrite.relabelConfig
	// This relabeling is applied to all the collected metrics before sending them to remote storage.
	// +optional
	RelabelConfig *v1.ConfigMapKeySelector `json:"relabelConfig,omitempty"`
	// InlineRelabelConfig - defines GlobalRelabelConfig for vmagent, can be defined directly at CRD.
	// +optional
	InlineRelabelConfig []vmetricsv1b1.RelabelConfig `json:"inlineRelabelConfig,omitempty"`
	// Namespace selector for PodMonitors
	PodMonitorNamespaceSelector *metav1.LabelSelector `json:"podMonitorNamespaceSelector,omitempty"`
	// Namespace selector for ServiceMonitors
	ServiceMonitorNamespaceSelector *metav1.LabelSelector `json:"serviceMonitorNamespaceSelector,omitempty"`
	// Selector for PodMoniotors
	PodMonitorSelector *metav1.LabelSelector `json:"podMonitorSelector,omitempty"`
	// Selector for ServiceMonitors
	ServiceMonitorSelector *metav1.LabelSelector `json:"serviceMonitorSelector,omitempty"`
	// Secrets is a list of Secrets in the same namespace as the vmagent
	// object, which shall be mounted into the vmagent Pods.
	// will be mounted at path /etc/vm/secrets
	// +optional
	Secrets []string `json:"secrets,omitempty"`
	// Volumes allows configuration of additional volumes on the output deploy definition.
	// Volumes specified will be appended to other volumes that are generated as a result of
	// StorageSpec objects.
	// +optional
	Volumes []v1.Volume `json:"volumes,omitempty"`
	// VolumeMounts allows configuration of additional VolumeMounts on the output deploy definition.
	// VolumeMounts specified will be appended to other VolumeMounts in the vmagent container,
	// that are generated as a result of StorageSpec objects.
	// +optional
	VolumeMounts []v1.VolumeMount `json:"volumeMounts,omitempty"`
	// ExtraArgs that will be passed to  VMAgent pod
	// for example remoteWrite.tmpDataPath: /tmp
	// it would be converted to flag --remoteWrite.tmpDataPath=/tmp
	// +optional
	ExtraArgs map[string]string `json:"extraArgs,omitempty"`
	// ExtraEnvs that will be added to VMAgent pod
	// +optional
	ExtraEnvs []v1.EnvVar `json:"extraEnvs,omitempty"`
	// TerminationGracePeriodSeconds period for container graceful termination
	// +optional
	TerminationGracePeriodSeconds *int64 `json:"terminationGracePeriodSeconds,omitempty"`
	//EnforcedNamespaceLabel enforces adding a namespace label of origin for each alert and metric that is user created.
	//The label value will always be the namespace of the object that is being created.
	//+optional
	EnforcedNamespaceLabel *string `json:"enforcedNamespaceLabel,omitempty"`
	// PriorityClassName assigned to the Pods
	// +optional
	PriorityClassName string `json:"priorityClassName,omitempty"`
	// TLS Configuration
	// +optional
	TLSConfig *VmTLSConfig `json:"tlsConfig,omitempty"`
}

type VmAlertManager struct {
	// Install indicates is AlertManager will be installed.
	// Can be changed for already deployed service and the service
	// will be removed during next reconciliation iteration
	Install *bool `json:"install,omitempty"`
	// Image to use for a `AlertManager` deployment.
	// The `AlertManager` is alerting system which read metrics from Prometheus
	// More info: https://prometheus.io/docs/alerting/alertmanager/
	Image string `json:"image"`
	// Ingress allows to create Ingress for VM UI.
	Ingress *Ingress `json:"ingress,omitempty"`
	// Secrets is a list of Secrets in the same namespace as the VMAlertmanager
	// object, which shall be mounted into the VMAlertmanager Pods.
	// The Secrets are mounted into /etc/vm/secrets/<secret-name>
	// +optional
	Secrets []string `json:"secrets,omitempty"`
	// ConfigRawYaml - raw configuration for alertmanager,
	// it helps it to start without secret.
	// priority -> hardcoded ConfigRaw -> ConfigRaw, provided by user -> ConfigSecret.
	// +optional
	ConfigRawYaml string `json:"configRawYaml,omitempty"`
	// ConfigSecret is the name of a Kubernetes Secret in the same namespace as the
	// VMAlertmanager object, which contains configuration for this VMAlertmanager,
	// configuration must be inside secret key: alertmanager.yaml.
	// It must be created by user.
	// instance. Defaults to 'vmalertmanager-<alertmanager-name>'
	// The secret is mounted into /etc/alertmanager/config.
	// +optional
	ConfigSecret string `json:"configSecret,omitempty"`
	// ReplicaCount Size is the expected size of the alertmanager cluster. The controller will
	// eventually make the size of the running cluster equal to the expected
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`
	// Retention Time duration VMAlertmanager shall retain data for. Default is '120h',
	// and must match the regular expression `[0-9]+(ms|s|m|h)` (milliseconds seconds minutes hours).
	// +optional
	Retention string `json:"retention,omitempty"`
	// Set paused to reconcilation
	Paused bool `json:"paused,omitempty"`
	// NodeSelector Define which Nodes the Pods are scheduled on.
	// Specified just as map[string]string. For example: "type: compute"
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// Resources defines resources requests and limits for single Pods.
	Resources v1.ResourceRequirements `json:"resources,omitempty"`
	// Affinity If specified, the pod's scheduling constraints.
	// +optional
	Affinity *v1.Affinity `json:"affinity,omitempty"`
	// Tolerations allow the pods to schedule onto nodes with matching taints.
	Tolerations []v1.Toleration `json:"tolerations,omitempty"`
	// SecurityContext holds pod-level security attributes.
	SecurityContext *SecurityContext `json:"securityContext,omitempty"`
	// Containers allows injecting additional containers or patching existing containers.
	// This is meant to allow adding an authentication proxy to an VMAlertmanager pod.
	Containers []v1.Container `json:"containers,omitempty"`
	// SelectAllByDefault changes default behavior for empty CRD selectors, such ConfigSelector.
	// with selectAllScrapes: true and undefined ConfigSelector and ConfigNamespaceSelector
	// Operator selects all exist alertManagerConfigs
	// with selectAllScrapes: false - selects nothing
	// +optional
	SelectAllByDefault bool `json:"selectAllByDefault,omitempty"`
	// ConfigSelector defines selector for VMAlertmanagerConfig, result config will be merged with with Raw or Secret config.
	// Works in combination with NamespaceSelector.
	// NamespaceSelector nil - only objects at VMAlertmanager namespace.
	// Selector nil - only objects at NamespaceSelector namespaces.
	// If both nil - behaviour controlled by selectAllByDefault
	// +optional
	ConfigSelector *metav1.LabelSelector `json:"configSelector,omitempty"`
	// ConfigNamespaceSelector defines namespace selector for VMAlertmanagerConfig.
	// Works in combination with Selector.
	// NamespaceSelector nil - only objects at VMAlertmanager namespace.
	// Selector nil - only objects at NamespaceSelector namespaces.
	// If both nil - behaviour controlled by selectAllByDefault
	// +optional
	ConfigNamespaceSelector *metav1.LabelSelector `json:"configNamespaceSelector,omitempty"`
	// TerminationGracePeriodSeconds period for container graceful termination
	// +optional
	TerminationGracePeriodSeconds *int64 `json:"terminationGracePeriodSeconds,omitempty"`
	// Map of string keys and values that can be used to organize and categorize
	// (scope and select) objects. May match selectors of replication controllers
	// and services.
	// More info: https://kubernetes.io/docs/user-guide/labels
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
	// Annotations is an unstructured key value map stored with a resource that may be
	// set by external tools to store and retrieve arbitrary metadata. They are not
	// queryable and should be preserved when modifying objects.
	// More info: https://kubernetes.io/docs/user-guide/annotations
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// Storage is the definition of how storage will be used by the VMAlertmanager
	// instances.
	// +optional
	Storage *vmetricsv1b1.StorageSpec `json:"storage,omitempty"`
	// Volumes allows configuration of additional volumes on the output deploy definition.
	// Volumes specified will be appended to other volumes that are generated as a result of
	// StorageSpec objects.
	// +optional
	Volumes []v1.Volume `json:"volumes,omitempty"`
	// VolumeMounts allows configuration of additional VolumeMounts on the output deploy definition.
	// VolumeMounts specified will be appended to other VolumeMounts in the alertmanager container,
	// that are generated as a result of StorageSpec objects.
	// +optional
	VolumeMounts []v1.VolumeMount `json:"volumeMounts,omitempty"`
	// ExtraArgs that will be passed to  VMAlertmanager pod for example log.level: debug
	// +optional
	ExtraArgs map[string]string `json:"extraArgs,omitempty"`
	// ExtraEnvs that will be added to VMAlertmanager pod
	// +optional
	ExtraEnvs []v1.EnvVar `json:"extraEnvs,omitempty"`
	// PriorityClassName assigned to the Pods
	// +optional
	PriorityClassName string `json:"priorityClassName,omitempty"`
	// TLS Configuration
	// +optional
	TLSConfig *VmTLSConfig `json:"tlsConfig,omitempty"`
	// WebConfig defines configuration for webserver
	// https://github.com/prometheus/alertmanager/blob/main/docs/https.md
	WebConfig *vmetricsv1b1.AlertmanagerWebConfig `json:"webConfig,omitempty"`
	// GossipConfig defines gossip TLS configuration for Alertmanager cluster
	GossipConfig *vmetricsv1b1.AlertmanagerGossipConfig `json:"gossipConfig,omitempty"`
}

type VmAlert struct {
	// Install indicates is VmAlert will be installed.
	// Can be changed for already deployed service and the service
	// will be removed during next reconciliation iteration
	Install *bool `json:"install,omitempty"`
	// Image - docker image settings for VMAlert
	// if no specified operator uses default config version
	Image string `json:"image"`
	// Ingress allows to create Ingress for VM UI.
	Ingress *Ingress `json:"ingress,omitempty"`
	// Secrets is a list of Secrets in the same namespace as the VMAlert
	// object, which shall be mounted into the VMAlert Pods.
	// The Secrets are mounted into /etc/vm/secrets/<secret-name>.
	// +optional
	Secrets []string `json:"secrets,omitempty"`
	// ReplicaCount is the expected size of the VMAlert cluster. The controller will
	// eventually make the size of the running cluster equal to the expected size.
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`
	// Resources container resource request and limits, https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
	// +optional
	Resources v1.ResourceRequirements `json:"resources,omitempty"`
	// Affinity If specified, the pod's scheduling constraints.
	// +optional
	Affinity *v1.Affinity `json:"affinity,omitempty"`
	// Tolerations If specified, the pod's tolerations.
	// +optional
	Tolerations []v1.Toleration `json:"tolerations,omitempty"`
	// SecurityContext holds pod-level security attributes and common container settings.
	// This defaults to the default PodSecurityContext.
	// +optional
	SecurityContext *v1.PodSecurityContext `json:"securityContext,omitempty"`
	// Containers property allows to inject additions sidecars or to patch existing containers.
	// It can be useful for proxies, backup, etc.
	// +optional
	Containers []v1.Container `json:"containers,omitempty"`
	// EvaluationInterval how often evalute rules by default. Pattern:="[0-9]+(ms|s|m|h)
	// +optional
	EvaluationInterval string `json:"evaluationInterval,omitempty"`
	// SelectAllByDefault changes default behavior for empty CRD selectors, such RuleSelector.
	// with selectAllByDefault: true and empty serviceScrapeSelector and RuleNamespaceSelector
	// Operator selects all exist serviceScrapes
	// with selectAllByDefault: false - selects nothing
	// +optional
	SelectAllByDefault bool `json:"selectAllByDefault,omitempty"`
	// RuleSelector selector to select which VMRules to mount for loading alerting
	// rules from.
	// Works in combination with NamespaceSelector.
	// If both nil - behaviour controlled by selectAllByDefault
	// NamespaceSelector nil - only objects at VMAlert namespace.
	// +optional
	RuleSelector *metav1.LabelSelector `json:"ruleSelector,omitempty"`
	// RuleNamespaceSelector to be selected for VMRules discovery.
	// Works in combination with Selector.
	// If both nil - behaviour controlled by selectAllByDefault
	// NamespaceSelector nil - only objects at VMAlert namespace.
	// +optional
	RuleNamespaceSelector *metav1.LabelSelector `json:"ruleNamespaceSelector,omitempty"`

	// Port for listen
	// +optional
	Port string `json:"port,omitempty"`
	// RemoteWrite Optional URL to remote-write compatible storage to persist
	// vmalert state and rule results to.
	// Rule results will be persisted according to each rule.
	// Alerts state will be persisted in the form of time series named ALERTS and ALERTS_FOR_STATE
	// see -remoteWrite.url docs in vmalerts for details.
	// E.g. http://127.0.0.1:8428
	// +optional
	RemoteWrite *vmetricsv1b1.VMAlertRemoteWriteSpec `json:"remoteWrite,omitempty"`
	// RemoteRead Optional URL to read vmalert state (persisted via RemoteWrite)
	// This configuration only makes sense if alerts state has been successfully
	// persisted (via RemoteWrite) before.
	// see -remoteRead.url docs in vmalerts for details.
	// E.g. http://127.0.0.1:8428
	// +optional
	RemoteRead *vmetricsv1b1.VMAlertRemoteReadSpec `json:"remoteRead,omitempty"`
	// RulePath to the file with alert rules.
	// Supports patterns. Flag can be specified multiple times.
	// Examples:
	// -rule /path/to/file. Path to a single file with alerting rules
	// -rule dir/*.yaml -rule /*.yaml. Relative path to all .yaml files in folder,
	// absolute path to all .yaml files in root.
	// by default operator adds /etc/vmalert/configs/base/vmalert.yaml
	// +optional
	RulePath []string `json:"rulePath,omitempty"`
	// Notifier prometheus alertmanager endpoint spec. Required at least one of  notifier or notifiers. e.g. http://127.0.0.1:9093
	// If specified both notifier and notifiers, notifier will be added as last element to notifiers.
	Notifier *vmetricsv1b1.VMAlertNotifierSpec `json:"notifier,omitempty"`
	// Notifiers prometheus alertmanager endpoints. Required at least one of  notifier or notifiers. e.g. http://127.0.0.1:9093
	// If specified both notifier and notifiers, notifier will be added as last element to notifiers.
	// only one of notifier options could be chosen: notifierConfigRef or notifiers +  notifier
	// +optional
	Notifiers []vmetricsv1b1.VMAlertNotifierSpec `json:"notifiers,omitempty"`
	// NotifierConfigRef reference for secret with notifier configuration for vmalert
	// only one of notifier options could be chosen: notifierConfigRef or notifiers +  notifier
	// +optional
	NotifierConfigRef *v1.SecretKeySelector `json:"notifierConfigRef,omitempty"`
	// Datasource Victoria Metrics or VMSelect url. Required parameter. e.g. http://127.0.0.1:8428
	Datasource *vmetricsv1b1.VMAlertDatasourceSpec `json:"datasource,omitempty"`

	// ExtraArgs that will be passed to  VMAlert pod
	// for example -remoteWrite.tmpDataPath=/tmp
	// +optional
	ExtraArgs map[string]string `json:"extraArgs,omitempty"`
	// ExtraEnvs that will be added to VMAlert pod
	// +optional
	ExtraEnvs []v1.EnvVar `json:"extraEnvs,omitempty"`
	// ExternalLabels in the form 'name: value' to add to all generated recording rules and alerts.
	// +optional
	ExternalLabels map[string]string `json:"externalLabels,omitempty"`
	// NodeSelector Define which Nodes the Pods are scheduled on.
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// TerminationGracePeriodSeconds period for container graceful termination
	// +optional
	TerminationGracePeriodSeconds *int64 `json:"terminationGracePeriodSeconds,omitempty"`
	// Set paused to reconcilation
	Paused bool `json:"paused,omitempty"`
	// Map of string keys and values that can be used to organize and categorize
	// (scope and select) objects. May match selectors of replication controllers
	// and services.
	// More info: https://kubernetes.io/docs/user-guide/labels
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
	// Annotations is an unstructured key value map stored with a resource that may be
	// set by external tools to store and retrieve arbitrary metadata. They are not
	// queryable and should be preserved when modifying objects.
	// More info: https://kubernetes.io/docs/user-guide/annotations
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// Volumes allows configuration of additional volumes on the output Deployment definition.
	// Volumes specified will be appended to other volumes that are generated as a result of
	// StorageSpec objects.
	// +optional
	Volumes []v1.Volume `json:"volumes,omitempty"`
	// VolumeMounts allows configuration of additional VolumeMounts on the output Deployment definition.
	// VolumeMounts specified will be appended to other VolumeMounts in the VMAlert container,
	// that are generated as a result of StorageSpec objects.
	// +optional
	VolumeMounts []v1.VolumeMount `json:"volumeMounts,omitempty"`
	//EnforcedNamespaceLabel enforces adding a namespace label of origin for each alert and metric that is user created.
	//The label value will always be the namespace of the object that is being created.
	//+optional
	EnforcedNamespaceLabel *string `json:"enforcedNamespaceLabel,omitempty"`
	// PriorityClassName assigned to the Pods
	// +optional
	PriorityClassName string `json:"priorityClassName,omitempty"`
	// TLS Configuration
	// +optional
	TLSConfig *VmTLSConfig `json:"tlsConfig,omitempty"`
}

type VmAuth struct {
	// Install indicates is VmAuth will be installed.
	// Can be changed for already deployed service and the service
	// will be removed during next reconciliation iteration
	Install *bool `json:"install,omitempty"`
	// Image - docker image settings for VMAuth
	// if no specified operator uses default config version
	Image string `json:"image,omitempty"`
	// Set paused to reconcilation
	Paused bool `json:"paused,omitempty"`
	// Secrets is a list of Secrets in the same namespace as the VMAuth
	// object, which shall be mounted into the VMAuth Pods.
	Secrets []string `json:"secrets,omitempty"`
	// ConfigMaps is a list of ConfigMaps in the same namespace as the VMAuth
	// object, which shall be mounted into the VMAuth Pods.
	ConfigMaps []string `json:"configMaps,omitempty"`
	// ReplicaCount is the expected size of the VMAuth
	ReplicaCount *int32 `json:"replicaCount,omitempty"`
	// Volumes allows configuration of additional volumes on the output deploy definition.
	// Volumes specified will be appended to other volumes that are generated as a result of StorageSpec objects.
	Volumes []v1.Volume `json:"volumes,omitempty"`
	// VolumeMounts allows configuration of additional VolumeMounts on the output Deployment definition.
	// VolumeMounts specified will be appended to other VolumeMounts in the VMAuth container,
	// that are generated as a result of StorageSpec objects.
	VolumeMounts []v1.VolumeMount `json:"volumeMounts,omitempty"`
	// Resources container resource request and limits, https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
	// if not defined default resources from operator config will be used
	Resources v1.ResourceRequirements `json:"resources,omitempty"`
	// Affinity If specified, the pod's scheduling constraints.
	Affinity *v1.Affinity `json:"affinity,omitempty"`
	// Tolerations If specified, the pod's tolerations.
	Tolerations []v1.Toleration `json:"tolerations,omitempty"`
	// SecurityContext holds pod-level security attributes and common container settings.
	// This defaults to the default PodSecurityContext.
	SecurityContext *v1.PodSecurityContext `json:"securityContext,omitempty"`
	// Containers property allows to inject additions sidecars or to patch existing containers.
	// It can be useful for proxies, backup, etc.
	// +optional
	Containers []v1.Container `json:"containers,omitempty"`
	// Port listen port
	// +optional
	Port string `json:"port,omitempty"`
	// SelectAllByDefault changes default behavior for empty CRD selectors, such userSelector.
	// with selectAllByDefault: true and empty userSelector and userNamespaceSelector
	// Operator selects all exist users
	// with selectAllByDefault: false - selects nothing
	// +optional
	SelectAllByDefault bool `json:"selectAllByDefault,omitempty"`
	// UserSelector defines VMUser to be selected for config file generation.
	// Works in combination with NamespaceSelector.
	// NamespaceSelector nil - only objects at VMAuth namespace.
	// If both nil - behaviour controlled by selectAllByDefault
	// +optional
	UserSelector *metav1.LabelSelector `json:"userSelector,omitempty"`
	// UserNamespaceSelector Namespaces to be selected for  VMAuth discovery.
	// Works in combination with Selector.
	// NamespaceSelector nil - only objects at VMAuth namespace.
	// Selector nil - only objects at NamespaceSelector namespaces.
	// If both nil - behaviour controlled by selectAllByDefault
	// +optional
	UserNamespaceSelector *metav1.LabelSelector `json:"userNamespaceSelector,omitempty"`
	// ExtraArgs that will be passed to  VMAuth pod
	// for example remoteWrite.tmpDataPath: /tmp
	// +optional
	ExtraArgs map[string]string `json:"extraArgs,omitempty"`
	// ExtraEnvs that will be added to VMAuth pod
	// +optional
	ExtraEnvs []v1.EnvVar `json:"extraEnvs,omitempty"`
	// Ingress enables ingress configuration for VMAuth.
	Ingress *Ingress `json:"ingress,omitempty"`
	// Map of string keys and values that can be used to organize and categorize
	// (scope and select) objects. May match selectors of replication controllers
	// and services.
	// More info: https://kubernetes.io/docs/user-guide/labels
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
	// Annotations is an unstructured key value map stored with a resource that may be
	// set by external tools to store and retrieve arbitrary metadata. They are not
	// queryable and should be preserved when modifying objects.
	// More info: https://kubernetes.io/docs/user-guide/annotations
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// NodeSelector Define which Nodes the Pods are scheduled on.
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// TerminationGracePeriodSeconds period for container graceful termination
	// +optional
	TerminationGracePeriodSeconds *int64 `json:"terminationGracePeriodSeconds,omitempty"`
	// PriorityClassName assigned to the Pods
	// +optional
	PriorityClassName string `json:"priorityClassName,omitempty"`
	// TLS Configuration
	// +optional
	TLSConfig *VmTLSConfig `json:"tlsConfig,omitempty"`
}

type VmUser struct {
	// Install indicates is VmUser will be installed.
	// Can be changed for already deployed service and the service
	// will be removed during next reconciliation iteration
	Install *bool `json:"install,omitempty"`
	// Image - docker image settings for VMUser
	// if no specified operator uses default config version
	Image string `json:"image,omitempty"`
	// Set paused to reconcilation
	Paused bool `json:"paused,omitempty"`
	// Name of the VMUser object.
	// +optional
	Name *string `json:"name,omitempty"`
	// UserName basic auth username for accessing protected endpoint,
	// will be replaced with metadata.name of VMUser if omitted.
	// +optional
	UserName *string `json:"username,omitempty"`
	// Password basic auth password for accessing protected endpoint.
	// +optional
	Password *string `json:"password,omitempty"`
	// PasswordRef allows fetching password from user-create secret by its name and key.
	// +optional
	PasswordRef *v1.SecretKeySelector `json:"passwordRef,omitempty"`
	// TokenRef allows fetching token from user-created secrets by its name and key.
	// +optional
	TokenRef *v1.SecretKeySelector `json:"tokenRef,omitempty"`
	// GeneratePassword instructs operator to generate password for user
	// if spec.password if empty.
	// +optional
	GeneratePassword bool `json:"generatePassword,omitempty"`
	// BearerToken Authorization header value for accessing protected endpoint.
	// +optional
	BearerToken *string `json:"bearerToken,omitempty"`
	// TargetRefs - reference to endpoints, which user may access.
	TargetRefs []vmetricsv1b1.TargetRef `json:"targetRefs,omitempty"`
}

// VMClusterSpec defines the desired state of VMCluster
// +k8s:openapi-gen=true
type VmCluster struct {
	// Install indicates is VmCluster will be installed.
	// Can be changed for already deployed service and the service
	// will be removed during next reconciliation iteration
	Install *bool `json:"install,omitempty"`
	// ParsingError contents error with context if operator was failed to parse json object from kubernetes api server
	ParsingError string `json:"-" yaml:"-"`
	// RetentionPeriod for the stored metrics
	// Note VictoriaMetrics has data/ and indexdb/ folders
	// metrics from data/ removed eventually as soon as partition leaves retention period
	// reverse index data at indexdb rotates once at the half of configured retention period
	// https://docs.victoriametrics.com/Single-server-VictoriaMetrics.html#retention
	RetentionPeriod string `json:"retentionPeriod"`
	// ReplicationFactor defines how many copies of data make among
	// distinct storage nodes
	// +optional
	ReplicationFactor *int32 `json:"replicationFactor,omitempty"`

	// ServiceAccountName is the name of the ServiceAccount to use to run the
	// VMSelect, VMStorage and VMInsert Pods.
	// +optional
	ServiceAccountName string `json:"serviceAccountName,omitempty"`

	// ImagePullSecrets An optional list of references to secrets in the same namespace
	// to use for pulling images from registries
	// see https://kubernetes.io/docs/concepts/containers/images/#referring-to-an-imagepullsecrets-on-a-pod
	// +optional
	ImagePullSecrets []v1.LocalObjectReference `json:"imagePullSecrets,omitempty"`

	// ClusterVersion defines default images tag for all components.
	// it can be overwritten with component specific image.tag value.
	ClusterVersion string `json:"clusterVersion,omitempty"`

	// +optional
	VmSelect *vmetricsv1b1.VMSelect `json:"vmselect,omitempty"`
	//Image for VMSelect
	VmSelectImage string `json:"vmSelectImage,omitempty"`
	// TLS Configuration
	// +optional
	VmSelectTLSConfig *VmTLSConfig `json:"vmSelectTlsConfig,omitempty"`
	// Ingress enables ingress configuration for VMSelect.
	VmSelectIngress *Ingress `json:"vmSelectIngress,omitempty"`
	// +optional
	VmInsert *vmetricsv1b1.VMInsert `json:"vminsert,omitempty"`
	//Image for VMInsert
	VmInsertImage string `json:"vmInsertImage,omitempty"`
	// TLS Configuration
	// +optional
	VmInsertTLSConfig *VmTLSConfig `json:"vmInsertTlsConfig,omitempty"`
	// +optional
	VmStorage *vmetricsv1b1.VMStorage `json:"vmstorage,omitempty"`
	//Image for VMStorage
	VmStorageImage string `json:"vmStorageImage,omitempty"`
	// TLS Configuration
	// +optional
	VmStorageTLSConfig *VmTLSConfig `json:"vmStorageTlsConfig,omitempty"`
	// Paused If set to true all actions on the underlying managed objects are not
	// going to be performed, except for delete actions.
	// +optional
	Paused bool `json:"paused,omitempty"`
	// UseStrictSecurity enables strict security mode for component
	// it restricts disk writes access
	// uses non-root user out of the box
	// drops not needed security permissions
	// +optional
	UseStrictSecurity *bool `json:"useStrictSecurity,omitempty"`
}

// NodeExporter defines the desired state for some part of node-exporter deployment
type NodeExporter struct {
	// Install indicates is node-exporter will be installed.
	// Can be changed for already deployed service and the service
	// will be removed during next reconciliation iteration
	Install *bool `json:"install,omitempty"`
	// SetupSecurityContext indicates is PSP or SCC (depends on cluster type) need to be created.
	SetupSecurityContext bool `json:"setupSecurityContext,omitempty"`
	// Image to use for a `node-exporter` deployment.
	// The `node-exporter` is an exporter to collect metrics from VM
	// More info: https://github.com/prometheus/node_exporter
	Image string `json:"image"`
	// Port for `node-exporter` daemonset and service
	Port int32 `json:"port"`
	// Resources defines resources requests and limits for single Pods.
	Resources v1.ResourceRequirements `json:"resources,omitempty"`
	// SecurityContext holds pod-level security attributes.
	SecurityContext *SecurityContext `json:"securityContext,omitempty"`
	// NodeSelector select nodes for deploy
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// Set paused to reconsilation
	Paused bool `json:"paused,omitempty"`
	// Service monitor for pulling metrics
	ServiceMonitor *Monitor `json:"serviceMonitor,omitempty"`
	// Tolerations allow the pods to schedule onto nodes with matching taints.
	Tolerations []v1.Toleration `json:"tolerations,omitempty"`
	// If specified, the pod's scheduling constraints.
	// More info: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core
	Affinity *v1.Affinity `json:"affinity,omitempty"`
	// Map of string keys and values that can be used to organize and categorize
	// (scope and select) objects. May match selectors of replication controllers
	// and services.
	// More info: https://kubernetes.io/docs/user-guide/labels
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
	// Annotations is an unstructured key value map stored with a resource that may be
	// set by external tools to store and retrieve arbitrary metadata. They are not
	// queryable and should be preserved when modifying objects.
	// More info: https://kubernetes.io/docs/user-guide/annotations
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// ServiceAccount is a structure which allow specify annotations and labels for Service Account
	// which will use by Alertmanager for work in Kubernetes. Cna be use by external tools to store
	// and retrieve arbitrary metadata.
	// +optional
	ServiceAccount *EmbeddedObjectMetadata `json:"serviceAccount,omitempty"`
	// Directory for textfile collector
	// More info: https://github.com/prometheus/node_exporter#textfile-collector
	// +optional
	CollectorTextfileDirectory string `json:"collectorTextfileDirectory,omitempty"`
	// Additional node-exporter container arguments.
	// for example --collector.systemd
	// +optional
	ExtraArgs []string `json:"extraArgs,omitempty"`
	// PriorityClassName assigned to the Pods
	// +optional
	PriorityClassName string `json:"priorityClassName,omitempty"`
}

// Pushgateway defines the desired state for some part of pushgateway deployment
type Pushgateway struct {
	// Install indicates is pushgateway will be installed.
	// Can be changed for already deployed service and the service
	// will be removed during next reconciliation iteration
	Install *bool `json:"install,omitempty"`
	// Image to use for a `pushgateway` deployment.
	// The `pushgateway` is an exporter to collect metrics from VM
	// More info: https://github.com/prometheus/pushgateway
	Image string `json:"image"`
	// Set replicas
	Replicas *int32 `json:"replicas,omitempty"`
	// Additional pushgateway container arguments.
	ExtraArgs []string `json:"extraArgs,omitempty"`
	// Volumes allows configuration of additional volumes on the output StatefulSet definition.
	// Volumes specified will be appended to other volumes that are generated as a result of StorageSpec objects.
	// More info: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.20/#volume-v1-core
	Volumes []v1.Volume `json:"volumes,omitempty"`
	// VolumeMounts allows configuration of additional VolumeMounts on the output StatefulSet definition.
	// VolumeMounts specified will be appended to other VolumeMounts in the prometheus container,
	// that are generated as a result of StorageSpec objects.
	// More info: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.20/#volumemount-v1-core
	VolumeMounts []v1.VolumeMount `json:"volumeMounts,omitempty"`
	// PVC spec for Pushgateway. If specified, also adds flags
	// --persistence.file and --persistence.interval with default values,
	// creates volume and volumeMount with name "storage-volume" in the deployment.
	Storage *v1.PersistentVolumeClaimSpec `json:"storage,omitempty"`
	// Port for `pushgateway` deployment and service
	Port int32 `json:"port"`
	// Ingress allows to create Ingress.
	Ingress *Ingress `json:"ingress,omitempty"`
	// Resources defines resources requests and limits for single Pods
	Resources v1.ResourceRequirements `json:"resources,omitempty"`
	// SecurityContext holds pod-level security attributes.
	SecurityContext *SecurityContext `json:"securityContext,omitempty"`
	// NodeSelector select nodes for deploy
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// If specified, the pod's scheduling constraints.
	// More info: https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core
	Affinity *v1.Affinity `json:"affinity,omitempty"`
	// Set paused to reconsilation
	Paused bool `json:"paused,omitempty"`
	// Service monitor for pulling metrics
	ServiceMonitor *Monitor `json:"serviceMonitor,omitempty"`
	// Tolerations allow the pods to schedule onto nodes with matching taints.
	Tolerations []v1.Toleration `json:"tolerations,omitempty"`
	// Map of string keys and values that can be used to organize and categorize
	// (scope and select) objects. May match selectors of replication controllers
	// and services.
	// More info: https://kubernetes.io/docs/user-guide/labels
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
	// Annotations is an unstructured key value map stored with a resource that may be
	// set by external tools to store and retrieve arbitrary metadata. They are not
	// queryable and should be preserved when modifying objects.
	// More info: https://kubernetes.io/docs/user-guide/annotations
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// PriorityClassName assigned to the Pods
	// +optional
	PriorityClassName string `json:"priorityClassName,omitempty"`
}

// Integration handles parameters to set up Platform Monitoring integration with other monitoring tools and public clouds.
// Currently supports:
//   - Google Cloud Platform (integration with Google Cloud Operations)
//   - Jaeger (a distributed tracing platform)
type Integration struct {
	StackDriverIntegration *StackDriverIntegrationConfig `json:"stackdriver,omitempty"`
	Jaeger                 *Jaeger                       `json:"jaeger,omitempty"`
	ClickHouse             *ClickHouse                   `json:"clickHouse,omitempty"`
}

// Auth handles parameters to set up Platform Monitoring auth for services.
// Currently supports:
//   - IDP
type Auth struct {
	// Deprecated field. ClientID is not expected to be stored in CRD.
	// +kubebuilder:validation:Deprecated=true
	ClientID string `json:"clientId,omitempty"`
	// Deprecated field. ClientSecret is not expected to be stored in CRD.
	// +kubebuilder:validation:Deprecated=true
	ClientSecret string     `json:"clientSecret,omitempty"`
	LoginURL     string     `json:"loginUrl"`
	TokenURL     string     `json:"tokenUrl"`
	UserInfoURL  string     `json:"userInfoUrl"`
	TLSConfig    *TLSConfig `json:"tlsConfig,omitempty"`
}

// TLSConfig extends the safe TLS configuration with file parameters.
type TLSConfig struct {
	CASecret           *v1.SecretKeySelector `json:"caSecret,omitempty"`
	CertSecret         *v1.SecretKeySelector `json:"certSecret,omitempty"`
	KeySecret          *v1.SecretKeySelector `json:"keySecret,omitempty"`
	InsecureSkipVerify *bool                 `json:"insecureSkipVerify,omitempty"`
}

// OAuthProxy handles parameters to set up Platform Monitoring oauth proxy for services.
// Currently used in:
//   - Prometheus
//   - AlertManager
type OAuthProxy struct {
	Image string `json:"image"`
}

// PromTLSConfig define TLS configuration for Prometheus.
type PromTLSConfig struct {
	// GenerateCerts allows to use cert-manager to generate certificates.
	GenerateCerts *GenerateCerts `json:"generateCerts,omitempty"`
	// WebTLSConfig allows to configure paths to certificates and other parameters manually.
	WebTLSConfig *promv1.WebTLSConfig `json:"webTLSConfig,omitempty"`
}

// GenerateCerts define settings for cert-manager.
type GenerateCerts struct {
	Enabled    bool   `json:"enabled"`
	SecretName string `json:"secretName,omitempty"`
}

// StackDriverIntegrationConfig holds parameters to set up Platform Monitoring integration with Google Cloud Operations (GCO).
// Integration schema:
//   - Send metrics from Prometheus to GCO by deploying 'stackdriver-prometheus-sidecar' container
//     as sidecar to Prometheus pod. Allows specify filters for metrics to send.
type StackDriverIntegrationConfig struct {
	// Image of 'stackdriver-prometheus-sidecar'.
	// This service is deploying as sidecar container to Prometheus pod and
	// send metrics from Prometheus to GCO.
	Image string `json:"image"`
	// Identificator of project in Google Cloud
	ProjectID string `json:"projectId"`
	// Location where project is deployed in Google Cloud
	Location string `json:"location"`
	// Name of Kubernetes cluster in Google Cloud which will be monitored
	Cluster string `json:"cluster"`
	// List of filters for metrics which will be sent to GCO.
	// Filters use the same syntax as Prometheus instant vector selectors:
	// https://prometheus.io/docs/prometheus/latest/querying/basics/#instant-vector-selectors.
	MetricsFilters []string `json:"metricsFilter,omitempty"`
}

// Jaeger holds parameters to set up Platform Monitoring integration with Jaeger.
type Jaeger struct {
	// If true, looking for Jaeger Service in all namespaces and add Grafana DataSource for it service if it is found.
	CreateGrafanaDataSource bool `json:"createGrafanaDataSource,omitempty"`
}

type ClickHouse struct {
	// If true, looking for ClickHouse Service in all namespaces and add Grafana DataSource for it service if it is found.
	CreateGrafanaDataSource bool `json:"createGrafanaDataSource,omitempty"`
}

// SecurityContext holds pod-level security attributes.
// The parameters are required if a Pod Security Policy is enabled
// for Kubernetes cluster and required if a Security Context Constraints is enabled
// for Openshift cluster.
type SecurityContext struct {
	// The UID to run the entrypoint of the container process.
	// Defaults to user specified in image metadata if unspecified.
	RunAsUser *int64 `json:"runAsUser,omitempty"`
	// The GID to run the entrypoint of the container process.
	// Uses runtime default if unset.
	// May also be set in SecurityContext.  If set in both SecurityContext and
	// PodSecurityContext, the value specified in SecurityContext takes precedence
	// for that container.
	// +optional
	RunAsGroup *int64 `json:"runAsGroup,omitempty"`
	// A special supplemental group that applies to all containers in a pod.
	// Some volume types allow the Kubelet to change the ownership of that volume
	// to be owned by the pod:
	//
	// 1. The owning GID will be the FSGroup
	// 2. The setgid bit is set (new files created in the volume will be owned by FSGroup)
	// 3. The permission bits are OR'd with rw-rw----
	//
	// If unset, the Kubelet will not modify the ownership and permissions of any volume.
	FSGroup *int64 `json:"fsGroup,omitempty"`
}

// Ingress holds parameters to configure Ingress.
// Allows to set Ingress annotation to use e.g. ingress-nginx for
// services authentication: https://github.com/kubernetes/ingress-nginx
type Ingress struct {
	// Install indicates is Ingress will be installed.
	Install *bool `json:"install,omitempty"`
	// Host for routing.
	Host string `json:"host,omitempty"`
	// Labels allows to set additional labels to the Ingress.
	// Basic labels will be saved.
	Labels map[string]string `json:"labels,omitempty"`
	// Annotations allows to set annotations for the Ingress.
	Annotations map[string]string `json:"annotations,omitempty"`
	// IngressClassName allows to set name for the IngressClass cluster resource.
	IngressClassName *string `json:"ingressClassName,omitempty"`
	// TlsSecretName allows to set secret name which will be used for TLS setting for the Ingress for specified host.
	TLSSecretName string `json:"tlsSecretName,omitempty"`
}

// PlatformMonitoringCondition contains description of status of PlatformMonitoring
// +k8s:openapi-gen=true
type PlatformMonitoringCondition struct {
	Type               string `json:"type"`
	Status             string `json:"status"`
	Reason             string `json:"reason"`
	Message            string `json:"message"`
	LastTransitionTime string `json:"lastTransitionTime"`
}

// PlatformMonitoringStatus defines the observed state of PlatformMonitoring
type PlatformMonitoringStatus struct {
	Conditions []PlatformMonitoringCondition `json:"conditions"`
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:resource:path=platformmonitorings,scope=Namespaced
// PlatformMonitoring is the Schema for the platformmonitorings API
type PlatformMonitoring struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PlatformMonitoringSpec   `json:"spec,omitempty"`
	Status PlatformMonitoringStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// PlatformMonitoringList contains a list of PlatformMonitoring
type PlatformMonitoringList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PlatformMonitoring `json:"items"`
}

// Monitor handles parameters to set up Service or Pod Monitor
type Monitor struct {
	Install              *bool                  `json:"install,omitempty"`
	Interval             string                 `json:"interval,omitempty"`
	ScrapeTimeout        string                 `json:"scrapeTimeout,omitempty"`
	RelabelConfigs       []promv1.RelabelConfig `json:"relabelings,omitempty"`
	MetricRelabelConfigs []promv1.RelabelConfig `json:"metricRelabelings,omitempty"`
	Selector             *metav1.LabelSelector  `json:"Selector,omitempty"`
}

// GrafanaDashboards contains parameters for specifying dashboards to install
type GrafanaDashboards struct {
	Install *bool    `json:"install,omitempty"`
	List    []string `json:"list,omitempty"`
}

// PrometheusRule handles parameters to override PrometheusRule: alerts of recording rules
type PrometheusRule struct {
	Group    string `json:"group,omitempty"`
	Alert    string `json:"alert,omitempty"`
	Record   string `json:"record,omitempty"`
	For      string `json:"for,omitempty"`
	Expr     string `json:"expr,omitempty"`
	Severity string `json:"severity,omitempty"`
}

// PrometheusRules help to add and override Prometheus rules
type PrometheusRules struct {
	Install    *bool            `json:"install,omitempty"`
	RuleGroups []string         `json:"ruleGroups,omitempty"`
	Override   []PrometheusRule `json:"override,omitempty"`
}

// Promxy handles parameters to set up Platform Monitoring with Prometheus proxy.
type Promxy struct {
	Install *bool  `json:"install,omitempty"`
	Port    *int32 `json:"port,omitempty"`
}

func init() {
	SchemeBuilder.Register(&PlatformMonitoring{}, &PlatformMonitoringList{})
}

// IsInstall check if AlertManager should be installed
// Returns false if parameter `install` is false or not set
func (am AlertManager) IsInstall() bool {
	if am.Install != nil {
		return *am.Install
	}
	return false
}

// IsInstall check is Grafana need to be installed.
// Returns false if parameter `install` is false or all other parameters
// are empty what means that section isn't presented in CR.
func (g Grafana) IsInstall() bool {
	if g.Install == nil || *g.Install {
		return true
	} else if !*g.Install {
		return false
	}

	return true
}

// IsInstall check if Prometheus should be installed
// Returns false if parameter `install` is false or not set
func (p Prometheus) IsInstall() bool {
	if p.Install != nil {
		return *p.Install
	}
	return false
}

// IsInstall check is kube-state-metrics need to be installed.
// Returns false if parameter `install` is false or all other parameters
// are empty what means that section isn't presented in CR.
func (ksm KubeStateMetrics) IsInstall() bool {
	if ksm.Install == nil || *ksm.Install {
		return true
	} else if !*ksm.Install {
		return false
	}

	return true
}

// IsInstall check is node-exporter need to be installed.
// Returns false if parameter `install` is false or all other parameters
// are empty what means that section isn't presented in CR.
func (ne NodeExporter) IsInstall() bool {
	if ne.Install == nil || *ne.Install {
		return true
	} else if !*ne.Install {
		return false
	}

	return true
}

// IsInstall check is ingress need to be installed.
// Returns false if parameter `install` is false or host is empty.
func (i Ingress) IsInstall() bool {
	if i.Install == nil || *i.Install {
		return i.Host != ""
	}

	return false
}

// IsInstall check if GrafanaDashboards should be installed
// Returns false if parameter `install` is false or not set
func (gd GrafanaDashboards) IsInstall() bool {
	if gd.Install != nil {
		return *gd.Install
	}
	return false
}

// IsInstall check if PrometheusRules should be installed
// Returns false if parameter `install` is false or not set
func (pr PrometheusRules) IsInstall() bool {
	if pr.Install != nil {
		return *pr.Install
	}
	return false
}

// IsInstall check if Monitor should be installed
// Returns false if parameter `install` is false or not set
func (m Monitor) IsInstall() bool {
	if m.Install != nil {
		return *m.Install
	}
	return false
}

// IsInstall check if Promxy should be installed
// Returns false if parameter `install` is false or not set
func (px Promxy) IsInstall() bool {
	if px.Install != nil {
		return *px.Install
	}
	return false
}

// IsInstall check if Pushgateway should be installed
// Returns false if parameter `install` is false or not set
func (pg Pushgateway) IsInstall() bool {
	if pg.Install != nil {
		return *pg.Install
	}
	return false
}

// IsInstall check if VmOperator should be installed
// Returns true if parameter `install` is true or not set
func (vo VmOperator) IsInstall() bool {
	if vo.Install == nil || *vo.Install {
		return vo.Image != ""
	}

	return false
}

// IsInstall check if VmAgent should be installed
// Returns true if parameter `install` is true or not set
func (va VmAgent) IsInstall() bool {
	if va.Install == nil || *va.Install {
		return va.Image != ""
	}

	return false
}

// IsInstall check is VmSingle need to be installed.
// Returns false if parameter `install` is true or all other parameters
// are empty what means that section isn't presented in CR.
func (vs VmSingle) IsInstall() bool {
	if vs.Install == nil || *vs.Install {
		return vs.Image != ""
	}

	return false
}

// IsInstall check if VmAlertManager should be installed
// Returns true if parameter `install` is true or not set
func (am VmAlertManager) IsInstall() bool {
	if am.Install == nil || *am.Install {
		return am.Image != ""
	}

	return false
}

// IsInstall check if VmAlert should be installed
// Returns true if parameter `install` is true or not set
func (va VmAlert) IsInstall() bool {
	if va.Install == nil || *va.Install {
		return va.Image != ""
	}

	return false
}

// IsInstall check if VmAuth should be installed
// Returns true if parameter `install` is true or not set
func (va VmAuth) IsInstall() bool {
	if va.Install == nil || *va.Install {
		return va.Image != ""
	}

	return false
}

// IsInstall check if VmUser should be installed
// Returns true if parameter `install` is true or not set
func (vu VmUser) IsInstall() bool {
	if vu.Install == nil || *vu.Install {
		return true
	} else if !*vu.Install {
		return false
	}

	return true
}

// IsInstall check if VmCluster should be installed
// Returns true if parameter `install` is true or not set
func (vc VmCluster) IsInstall() bool {
	if vc.Install == nil || *vc.Install {
		return vc.VmInsertImage != "" && vc.VmSelectImage != "" && vc.VmStorageImage != ""
	}

	return false
}

// OverridePodMonitor overrides specified fields of PodMonitor
func (m Monitor) OverridePodMonitor(podMonitor *promv1.PodMonitor) {

	for i := range podMonitor.Spec.PodMetricsEndpoints {
		if m.Interval != "" {
			podMonitor.Spec.PodMetricsEndpoints[i].Interval = promv1.Duration(m.Interval)
		}
		if m.ScrapeTimeout != "" {
			podMonitor.Spec.PodMetricsEndpoints[i].ScrapeTimeout = promv1.Duration(m.ScrapeTimeout)
		}
		if len(m.RelabelConfigs) > 0 {
			podMonitor.Spec.PodMetricsEndpoints[i].RelabelConfigs = m.RelabelConfigs
		}
		if len(m.MetricRelabelConfigs) > 0 {
			podMonitor.Spec.PodMetricsEndpoints[i].MetricRelabelConfigs = m.MetricRelabelConfigs
		}
	}
}

// OverrideServiceMonitor overrides specified fields of serviceMonitor
func (m Monitor) OverrideServiceMonitor(serviceMonitor *promv1.ServiceMonitor) {

	for i := range serviceMonitor.Spec.Endpoints {
		if m.Interval != "" {
			serviceMonitor.Spec.Endpoints[i].Interval = promv1.Duration(m.Interval)
		}
		if m.ScrapeTimeout != "" {
			serviceMonitor.Spec.Endpoints[i].ScrapeTimeout = promv1.Duration(m.ScrapeTimeout)
		}
		if len(m.RelabelConfigs) > 0 {
			serviceMonitor.Spec.Endpoints[i].RelabelConfigs = m.RelabelConfigs
		}
		if len(m.MetricRelabelConfigs) > 0 {
			serviceMonitor.Spec.Endpoints[i].MetricRelabelConfigs = m.MetricRelabelConfigs
		}
	}
}

// OverridePrometheusRule overrides specified fields of prometheus rules
func (pr *PrometheusRule) OverridePrometheusRule(rule *promv1.Rule) {

	if pr.Expr != "" {
		rule.Expr = intstr.FromString(pr.Expr)
	}
	if pr.For != "" {
		*rule.For = promv1.Duration(pr.For)
	}
	if pr.Severity != "" {
		rule.Labels["severity"] = pr.Severity
	}
}

// FillEmptyWithDefaults fill empty fields with default values
func (pm *PlatformMonitoring) FillEmptyWithDefaults() {
	// Fill AlertManager
	if pm.Spec.AlertManager != nil {
		if pm.Spec.AlertManager.Image == "" {
			pm.Spec.AlertManager.Image = defaultAlertManagerImage
		}
	}

	// Fill Grafana
	if pm.Spec.Grafana != nil {
		if pm.Spec.Grafana.Image == "" {
			pm.Spec.Grafana.Image = defaultGrafanaImage
		}

		// Fill GrafanaOperator
		if pm.Spec.Grafana.Operator.Image == "" {
			pm.Spec.Grafana.Operator.Image = defaultGrafanaOperatorImage
		}
		if pm.Spec.Grafana.Operator.InitContainerImage == "" {
			pm.Spec.Grafana.Operator.InitContainerImage = defaultGrafanaOperatorInitContainerImage
		}
	}

	// Fill Prometheus
	if pm.Spec.Prometheus != nil {
		if pm.Spec.Prometheus.Image == "" {
			pm.Spec.Prometheus.Image = defaultPrometheusImage
		}
		if pm.Spec.Prometheus.ConfigReloaderImage == "" {
			pm.Spec.Prometheus.ConfigReloaderImage = defaultPrometheusConfigReloaderImage
		}

		// Fill PrometheusOperator
		if pm.Spec.Prometheus.Operator.Image == "" {
			pm.Spec.Prometheus.Operator.Image = defaultPrometheusOperatorImage
		}
	}

	if (pm.Spec.Prometheus != nil && pm.Spec.Prometheus.IsInstall()) &&
		(pm.Spec.Victoriametrics != nil && pm.Spec.Victoriametrics.VmOperator.IsInstall()) {
		// If prometheus and VM stack are installed together, consider VM stack as default
		pm.Spec.Prometheus.Install = ptr.To(false)
	}

	// if pm.Spec.Victoriametrics != nil && pm.Spec.Victoriametrics.VmCluster.IsInstall() {
	// 	// Disable VmSingle installation if VMCluster has to be installed
	// 	pm.Spec.Victoriametrics.VmSingle.Install = ptr.To(false)
	// }
}

type PlatformMonitoringTemplatingParameters struct {
	Values PlatformMonitoringSpec
	Release
	DashboardsUIDs map[string]string
}

type Release struct {
	Namespace string
}

func (in *PlatformMonitoring) ToParams() PlatformMonitoringTemplatingParameters {
	return PlatformMonitoringTemplatingParameters{
		Values: in.Spec,
		Release: Release{
			Namespace: in.GetNamespace(),
		},
	}
}
