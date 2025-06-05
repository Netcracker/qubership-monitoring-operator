<!-- markdownlint-disable MD013 -->
This section describes the types introduced by the Monitoring Operator.


> Note this document is generated from code comments. When contributing a change to this document please do so by changing the code comments.

## AlertManager

AlertManager defines the desired state for some part of the prometheus-operator deployment.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| image | Image to use for a `AlertManager` deployment. The `AlertManager` is alerting system which read metrics from Prometheus More info: [https://prometheus.io/docs/alerting/alertmanager/](https://prometheus.io/docs/alerting/alertmanager/) | string | true |
| port | Port for alertManager service | int32 | false |
| install | Install indicates is AlertManager will be installed. Can be changed for already deployed service and the service will be removed during next reconciliation iteration | *bool | false |
| securityContext | SecurityContext holds pod-level security attributes. | *[SecurityContext](#securitycontext) | false |
| resources | Resources defines resources requests and limits for single Pods. | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#resourcerequirements-v1-core) | false |
| containers | Containers allows injecting additional containers or modifying operator generated containers. This can be used to allow adding an authentication proxy to a Prometheus pod or to change the behavior of an operator generated container. | []v1.Container | false |
| ingress | Ingress allows to create Ingress for AlertManager UI. | *[Ingress](#ingress) | false |
| paused | Set paused to reconsilation | bool | false |
| replicas | Set replicas | *int32 | false |
| tolerations | Tolerations allow the pods to schedule onto nodes with matching taints. | []v1.Toleration | false |
| nodeSelector | NodeSelector Define which Nodes the Pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\" | map[string]string | false |
| podMonitor | Pod monitor for self monitoring | *[Monitor](#monitor) | false |
| labels | Map of string keys and values that can be used to organize and categorize (scope and select) objects. May match selectors of replication controllers and services. More info: [https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/) | map[string]string | false |
| annotations | Annotations is an unstructured key value map stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. They are not queryable and should be preserved when modifying objects. More info: [https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/](https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/) | map[string]string | false |
| serviceAccount | ServiceAccount is a structure which allow specify annotations and labels for Service Account which will use by Alertmanager for work in Kubernetes. Cna be use by external tools to store and retrieve arbitrary metadata. | *[EmbeddedObjectMetadata](#embeddedobjectmetadata) | false |




## Auth

Auth handles parameters to set up Platform Monitoring auth for services. It currently supports:\n * IDP

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| clientId |  | string | true |
| clientSecret |  | string | true |
| loginUrl |  | string | true |
| tokenUrl |  | string | true |
| userInfoUrl |  | string | true |
| tlsConfig |  | *[TLSConfig](#tlsconfig) | false |




## EmbeddedObjectMetadata

EmbeddedObjectMetadata contains a subset of the fields included in k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta Only fields which are relevant to embedded resources are included.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| annotations | Annotations is an unstructured key value map stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. They are not queryable and should be preserved when modifying objects. More info: [https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/](https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/) | map[string]string | false |
| labels | Map of string keys and values that can be used to organize and categorize (scope and select) objects. May match selectors of replication controllers and services. More info: [https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/) | map[string]string | false |




## Grafana

Grafana defines the desired state for some part of the grafana-operator deployment.

| Field                      | Description                                                                                                                                                                                                                                                                                                                         | Scheme                                                                                                                       | Required |
|----------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------|----------|
| image                      | Image to use for a `grafana` deployment. The `grafana` is a web ui to show graphics. More info: [https://github.com/grafana/grafana](https://github.com/grafana/grafana)                                                                                                                                                            | string                                                                                                                       | true     |
| operator                   | Operator parameters                                                                                                                                                                                                                                                                                                                 | [GrafanaOperator](#grafanaoperator)                                                                                          | true     |
| install                    | Install indicates is Grafana will be installed. Can be changed for already deployed service and the service will be removed during next reconciliation iteration                                                                                                                                                                    | *bool                                                                                                                        | false    |
| securityContext            | SecurityContext holds pod-level security attributes.                                                                                                                                                                                                                                                                                | *[SecurityContext](#securitycontext)                                                                                         | false    |
| resources                  | Resources defines resources requests and limits for single Pods.                                                                                                                                                                                                                                                                    | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#resourcerequirements-v1-core) | false    |
| ingress                    | Ingress allows to create Ingress for Grafana UI.                                                                                                                                                                                                                                                                                    | *[Ingress](#ingress)                                                                                                         | false    |
| config                     | Config allows to override Config for Grafana.                                                                                                                                                                                                                                                                                       | grafv1alpha1.GrafanaConfig                                                                                                   | false    |
| grafanaHomeDashboard       | Custom grafana home dashboard                                                                                                                                                                                                                                                                                                       | bool                                                                                                                         | false    |
| backupDaemonDashboard      | Enables Backup Daemon Dashboard installation.                                                                                                                                                                                                                                                                                       | bool                                                                                                                         | false    |
| dashboardLabelSelector     | Allows to query over a set of resources according to labels.<br/>The result of matchLabels and matchExpressions are ANDed.<br/>An empty label selector matches all objects. A null label selector matches no objects.                                                                                                               | []*[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#labelselector-v1-meta)        | false    |
| dashboardNamespaceSelector | Allows to query over a set of resources in namespaces that fits label selector.                                                                                                                                                                                                                                                     | *[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#labelselector-v1-meta)          | false    |
| paused                     | Set paused to reconciliation                                                                                                                                                                                                                                                                                                        | bool                                                                                                                         | false    |
| tolerations                | Tolerations allow the pods to schedule onto nodes with matching taints.                                                                                                                                                                                                                                                             | []v1.Toleration                                                                                                              | false    |
| nodeSelector               | NodeSelector Define which Nodes the Pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\"                                                                                                                                                                                                      | map[string]string                                                                                                            | false    |
| podMonitor                 | Pod monitor for self monitoring                                                                                                                                                                                                                                                                                                     | *[Monitor](#monitor)                                                                                                         | false    |
| dataStorage                | DataStorage provides a means to configure the grafana data storage                                                                                                                                                                                                                                                                  | *grafv1alpha1.GrafanaDataStorage                                                                                             | false    |
| labels                     | Map of string keys and values that can be used to organize and categorize (scope and select) objects. May match selectors of replication controllers and services. More info: [https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/)                                                          | map[string]string                                                                                                            | false    |
| annotations                | Annotations is an unstructured key value map stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. They are not queryable and should be preserved when modifying objects. More info: [https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/) | map[string]string                                                                                                            | false    |
| replicas                   | Set replicas                                                                                                                                                                                                                                                                                                                        | *int32                                                                                                                       | false    |
| serviceAccount             | ServiceAccount is a structure which allow specify annotations and labels for Service Account which will use by Grafana for work in Kubernetes. Cna be use by external tools to store and retrieve arbitrary metadata.                                                                                                               | *[EmbeddedObjectMetadata](#embeddedobjectmetadata)                                                                           | false    |




## GrafanaDashboards

GrafanaDashboards contains parameters for specifying the dashboards to install.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| install |  | *bool | false |
| list |  | []string | false |




## GrafanaOperator

GrafanaOperator defines the desired state for some part of the grafana-operator deployment.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| image | Image to use for a `grafana-operator` deployment. The `grafana-operator` is a control, deploy and process custom resources into Grafana entities. More info: [https://github.com/grafana/grafana-operator](https://github.com/grafana/grafana-operator) | string | true |
| initContainerImage | Image to use to initialize Grafana deployment. | string | true |
| namespaces | Namespaces to scope the interaction of the Grafana operator. | string | false |
| securityContext | SecurityContext holds pod-level security attributes. | *[SecurityContext](#securitycontext) | false |
| resources | Resources defines resources requests and limits for single Pods. | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#resourcerequirements-v1-core) | false |
| paused | Set paused to reconsilation | bool | false |
| tolerations | Tolerations allow the pods to schedule onto nodes with matching taints. | []v1.Toleration | false |
| nodeSelector | NodeSelector Define which Nodes the Pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\" | map[string]string | false |
| podMonitor | Pod monitor for self monitoring | *[Monitor](#monitor) | false |
| labels | Map of string keys and values that can be used to organize and categorize (scope and select) objects. May match selectors of replication controllers and services. More info: [https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/) | map[string]string | false |
| annotations | Annotations is an unstructured key value map stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. They are not queryable and should be preserved when modifying objects. More info: [https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/](https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/) | map[string]string | false |
| logLevel | Zap log level (one of 'debug', 'info', 'error' or any integer value > 0) (default info) More info: [https://github.com/grafana/grafana-operator/blob/v4/documentation/deploy_grafana.md](https://github.com/grafana/grafana-operator/blob/v4/documentation/deploy_grafana.md) | string | false |
| serviceAccount | ServiceAccount is a structure which allow specify annotations and labels for Service Account which will use by Alertmanager for work in Kubernetes. Cna be use by external tools to store and retrieve arbitrary metadata. | *[EmbeddedObjectMetadata](#embeddedobjectmetadata) | false |




## Ingress

Ingress holds parameters to configure Ingress. It allows to set Ingress annotation to use. For example, ingress-nginx
for services authentication. For more information, refer
to [https://github.com/kubernetes/ingress-nginx](https://github.com/kubernetes/ingress-nginx).

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| install | Install indicates is Ingress will be installed. | *bool | false |
| host | Host for routing. | string | false |
| labels | Labels allows to set additional labels to the Ingress. Basic labels will be saved. | map[string]string | false |
| annotations | Annotations allows to set annotations for the Ingress. | map[string]string | false |
| tlsSecretName | TlsSecretName allows to set secret name which will be used for TLS setting for the Ingress for specified host. | string | false |




## Integration

Integration handles parameters to set up the Platform Monitoring integration with public clouds. Currently it supports:
\n * Google Cloud Platform (integration with Google Cloud Operations).

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| stackdriver |  | *[StackDriverIntegrationConfig](#stackdriverintegrationconfig) | false |
| jaeger |  | *[Jaeger](#jaeger) | false |




## Jaeger

Jaeger holds parameters to set up Platform Monitoring integration with Jaeger.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| createGrafanaDataSource | If true, looking for Jaeger Service in all namespaces and add Grafana DataSource for it service if it is found. | bool | false |




## KubeStateMetrics

KubeStateMetrics defines the desired state for some part of the kube-state-metrics deployment.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| install | Install indicates is kube-state-metrics will be installed. Can be changed for already deployed service and the service will be removed during next reconciliation iteration | *bool | false |
| image | Image to use for a `kube-state-metrics` deployment. The `kube-state-metrics` is an exporter to collect Kubernetes metrics More info: [https://github.com/kubernetes/kube-state-metrics](https://github.com/kubernetes/kube-state-metrics) | string | true |
| resources | Resources defines resources requests and limits for single Pods. | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#resourcerequirements-v1-core) | false |
| securityContext | SecurityContext holds pod-level security attributes. | *[SecurityContext](#securitycontext) | false |
| namespaces | List of comma-separated namespaces to scrape metrics in non-privileged mode. | string | false |
| scrapeResources | Comma-separated list of Resources to be enabled. | string | false |
| metricLabelsAllowlist | Comma-separated list of additional Kubernetes label keys that will be used in the resource labels metric. | string | false |
| paused | Set paused to reconsilation. | bool | false |
| tolerations | Tolerations allow the pods to schedule onto nodes with matching taints. | []v1.Toleration | false |
| nodeSelector | NodeSelector Define which Nodes the Pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\" | map[string]string | false |
| serviceMonitor | Service monitor for pulling metrics | *[Monitor](#monitor) | false |
| labels | Map of string keys and values that can be used to organize and categorize (scope and select) objects. May match selectors of replication controllers and services. More info: [https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/) | map[string]string | false |
| annotations | Annotations is an unstructured key value map stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. They are not queryable and should be preserved when modifying objects. More info: [https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/](https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/) | map[string]string | false |
| serviceAccount | ServiceAccount is a structure which allow specify annotations and labels for Service Account which will use by Alertmanager for work in Kubernetes. Cna be use by external tools to store and retrieve arbitrary metadata. | *[EmbeddedObjectMetadata](#embeddedobjectmetadata) | false |




## Monitor

Monitor handles parameters to set up Service or Pod Monitor.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| install |  | *bool | false |
| interval |  | string | false |
| scrapeTimeout |  | string | false |
| relabelings |  | []*promv1.RelabelConfig | false |
| metricRelabelings |  | []*promv1.RelabelConfig | false |
| Selector |  | *[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#labelselector-v1-meta) | false |




## NodeExporter

NodeExporter defines the desired state for some part of the node-exporter deployment.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| install | Install indicates is node-exporter will be installed. Can be changed for already deployed service and the service will be removed during next reconciliation iteration | *bool | false |
| setupSecurityContext | SetupSecurityContext indicates is PSP or SCC (depends on cluster type) need to be created. | bool | false |
| image | Image to use for a `node-exporter` deployment. The `node-exporter` is an exporter to collect metrics from VM More info: [https://github.com/prometheus/node_exporter](https://github.com/prometheus/node_exporter) | string | true |
| port | Port for `node-exporter` daemonset and service | int32 | true |
| resources | Resources defines resources requests and limits for single Pods. | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#resourcerequirements-v1-core) | false |
| securityContext | SecurityContext holds pod-level security attributes. | *[SecurityContext](#securitycontext) | false |
| nodeSelector | NodeSelector select nodes for deploy | map[string]string | false |
| paused | Set paused to reconsilation | bool | false |
| serviceMonitor | Service monitor for pulling metrics | *[Monitor](#monitor) | false |
| tolerations | Tolerations allow the pods to schedule onto nodes with matching taints. | []v1.Toleration | false |
| labels | Map of string keys and values that can be used to organize and categorize (scope and select) objects. May match selectors of replication controllers and services. More info: [https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/) | map[string]string | false |
| annotations | Annotations is an unstructured key value map stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. They are not queryable and should be preserved when modifying objects. More info: [https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/](https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/) | map[string]string | false |
| serviceAccount | ServiceAccount is a structure which allow specify annotations and labels for Service Account which will use by Alertmanager for work in Kubernetes. Cna be use by external tools to store and retrieve arbitrary metadata. | *[EmbeddedObjectMetadata](#embeddedobjectmetadata) | false |
| CollectorTextfileDirectory | Directory for textfile collector. More info: [https://github.com/prometheus/node_exporter#textfile-collector](https://github.com/prometheus/node_exporter#textfile-collector) | string | false |
| extraArgs | Additional arguments for node-exporter container. For example: "--collector.systemd". | list[string] | false |




## OAuthProxy

OAuthProxy handles parameters to set up Platform Monitoring oauth proxy for services. It is currently used in:

* Prometheus
* AlertManager

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| image |  | string | true |




## PlatformMonitoring

PlatformMonitoring is the Schema for the platformmonitorings API.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#objectmeta-v1-meta) | false |
| spec |  | [PlatformMonitoringSpec](#platformmonitoringspec) | false |
| status |  | [PlatformMonitoringStatus](#platformmonitoringstatus) | false |




## PlatformMonitoringCondition

PlatformMonitoringCondition contains the description of the status of PlatformMonitoring.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| type |  | string | true |
| status |  | string | true |
| reason |  | string | true |
| message |  | string | true |
| lastTransitionTime |  | string | true |




## PlatformMonitoringList

PlatformMonitoringList contains a list of PlatformMonitoring.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#listmeta-v1-meta) | false |
| items |  | \[\][PlatformMonitoring](#platformmonitoring) | true |




## PlatformMonitoringSpec

PlatformMonitoringSpec defines the desired state of PlatformMonitoring.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| alertManager |  | *[AlertManager](#alertmanager) | false |
| kubeStateMetrics |  | *[KubeStateMetrics](#kubestatemetrics) | false |
| prometheus |  | *[Prometheus](#prometheus) | false |
| nodeExporter |  | *[NodeExporter](#nodeexporter) | false |
| grafana |  | *[Grafana](#grafana) | false |
| integration |  | *[Integration](#integration) | false |
| auth |  | *[Auth](#auth) | false |
| oAuthProxy |  | *[OAuthProxy](#oauthproxy) | false |
| kubernetesMonitors |  | map\[string\][Monitor](#monitor) | false |
| grafanaDashboards |  | *[GrafanaDashboards](#grafanadashboards) | false |
| prometheusRules |  | *[PrometheusRules](#prometheusrules) | false |
| promxy |  | *[Promxy](#promxy) | false |
| pushgateway |  | *[Pushgateway](#pushgateway) | false |
| publicCloudName |  | string | false |
| victoriametrics |  | \*[Victoriametrics](#victoriametrics) | false |




## PlatformMonitoringStatus

PlatformMonitoringStatus defines the observed state of PlatformMonitoring.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| conditions |  | \[\][PlatformMonitoringCondition](#platformmonitoringcondition) | true |




## Prometheus

Prometheus defines the link to PrometheusSpec objects from prometheus-operator.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| image | Image to use for a `prometheus` deployment. The `prometheus` is a systems and service monitoring system. It collects metrics from configured targets at given intervals. More info: [https://github.com/prometheus/prometheus](https://github.com/prometheus/prometheus) | string | true |
| operator | Operator parameters | [PrometheusOperator](#prometheusoperator) | true |
| install | Install indicates is Prometheus will be installed. Can be changed for already deployed service and the service will be removed during next reconciliation iteration | *bool | false |
| configReloaderImage | Image to use for a `prometheus-config-reloader` The `prometheus-config-reloaded` is an add-on to prometheus that monitors changes in prometheus.yaml and an HTTP request reloads the prometheus configuration. More info: [https://github.com/prometheus-operator/prometheus-operator/tree/master/cmd/prometheus-config-reloader](https://github.com/prometheus-operator/prometheus-operator/tree/master/cmd/prometheus-config-reloader) | string | false |
| remoteWrite | RemoteWriteSpec defines the remote_write configuration for prometheus. The `remote_write` allows transparently send samples to a long term storage. More info: [https://prometheus.io/docs/operating/integrations/#remote-endpoints-and-storage](https://prometheus.io/docs/operating/integrations/#remote-endpoints-and-storage) | []promv1.RemoteWriteSpec | false |
| remoteRead | RemoteReadSpec defines the remote_read configuration for prometheus. The `remote_read` allows transparently receive samples from a long term storage. More info: [https://prometheus.io/docs/operating/integrations/#remote-endpoints-and-storage](https://prometheus.io/docs/operating/integrations/#remote-endpoints-and-storage) | []promv1.RemoteReadSpec | false |
| secrets | Secrets is a list of Secrets in the same namespace as the Prometheus object, which shall be mounted into the Prometheus Pods. The Secrets are mounted into /etc/prometheus/secrets/\<secret-name\>. | []string | false |
| alerting | Define details regarding alerting. | *promv1.AlertingSpec | false |
| externalLabels | The labels to add to any time series or alerts when communicating with external systems (federation, remote storage, Alertmanager). | map[string]string | false |
| securityContext | SecurityContext holds pod-level security attributes. | *[SecurityContext](#securitycontext) | false |
| nodeSelector | Define which Nodes the Pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\" | map[string]string | false |
| affinity | If specified, the pod's scheduling constraints. More info: [https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#affinity-v1-core](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#affinity-v1-core) | *v1.Affinity | false |
| resources | Resources defines resources requests and limits for single Pods. | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#resourcerequirements-v1-core) | false |
| storage | Storage spec to specify how storage shall be used. More info: [https://github.com/prometheus-operator/prometheus-operator/blob/v0.79.2/Documentation/api.md#storagespec](https://github.com/prometheus-operator/prometheus-operator/blob/v0.79.2/Documentation/api.md#storagespec) | *promv1.StorageSpec | false |
| volumes | Volumes allows configuration of additional volumes on the output StatefulSet definition. Volumes specified will be appended to other volumes that are generated as a result of StorageSpec objects. More info: [https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#volume-v1-core](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#volume-v1-core) | []v1.Volume | false |
| volumeMounts | VolumeMounts allows configuration of additional VolumeMounts on the output StatefulSet definition. VolumeMounts specified will be appended to other VolumeMounts in the prometheus container, that are generated as a result of StorageSpec objects. More info: [https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#volumemount-v1-core](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#volumemount-v1-core) | []v1.VolumeMount | false |
| ingress | Ingress allows to create Ingress for Prometheus UI. | *[Ingress](#ingress) | false |
| retention | Retention policy by time | string | false |
| retentionsize | Retention policy by size [EXPERIMENTAL] | string | false |
| containers | Containers allows injecting additional containers or modifying operator generated containers. This can be used to allow adding an authentication proxy to a Prometheus pod or to change the behavior of an operator generated container. | []v1.Container | false |
| externalUrl | The external URL the Prometheus instances will be available under. This is necessary to generate correct URLs. This is necessary if Prometheus is not served from root of a DNS name. | string | false |
| paused | Set paused to reconsilation | bool | false |
| replicas | Set replicas | *int32 | false |
| tolerations | Tolerations allow the pods to schedule onto nodes with matching taints. | []v1.Toleration | false |
| podMonitor | Pod monitor for self monitoring | *[Monitor](#monitor) | false |
| ruleNamespaceSelector | Namespace selector for rules | *[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#labelselector-v1-meta) | false |
| podMonitorNamespaceSelector | Namespace selector for PodMonitors | *[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#labelselector-v1-meta) | false |
| serviceMonitorNamespaceSelector | Namespace selector for ServiceMonitors | *[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#labelselector-v1-meta) | false |
| ruleSelector | Selector for rules | *[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#labelselector-v1-meta) | false |
| podMonitorSelector | Selector for PodMoniotors | *[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#labelselector-v1-meta) | false |
| serviceMonitorSelector | Selector for ServiceMonitors | *[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#labelselector-v1-meta) | false |
| query | QuerySpec defines the query command line flags when starting Prometheus | *promv1.QuerySpec | false |
| tlsConfig | TLS define TLS configuration for Prometheus. | *promv1.WebTLSConfig | false |
| enableAdminAPI | Enable access to prometheus web admin API. Defaults to the value of false. WARNING: Enabling the admin APIs enables mutating endpoints, to delete data, shutdown Prometheus, and more. Enabling this should be done with care and the user is advised to add additional authentication authorization via a proxy to ensure only clients authorized to perform these actions can do so. For more information see [https://prometheus.io/docs/prometheus/latest/querying/api/#tsdb-admin-apis](https://prometheus.io/docs/prometheus/latest/querying/api/#tsdb-admin-apis) | bool | false |
| labels | Map of string keys and values that can be used to organize and categorize (scope and select) objects. May match selectors of replication controllers and services. More info: [https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/) | map[string]string | false |
| annotations | Annotations is an unstructured key value map stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. They are not queryable and should be preserved when modifying objects. More info: [https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/](https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/) | map[string]string | false |
| serviceAccount | ServiceAccount is a structure which allow specify annotations and labels for Service Account which will use by Prometheus for work in Kubernetes. Cna be use by external tools to store and retrieve arbitrary metadata. | *[EmbeddedObjectMetadata](#embeddedobjectmetadata) | false |
| replicaExternalLabelName | Name of Prometheus external label used to denote replica name. Defaults to the value of `prometheus_replica`. External label will _not_ be added when value is set to empty string (`\"\"`). | *string | false |
| enableFeatures | Enable access to Prometheus disabled features. By default, no features are enabled. Enabling disabled features is entirely outside the scope of what the maintainers will support and by doing so, you accept that this behavior may break at any time without notice. For more information see [https://prometheus.io/docs/prometheus/latest/disabled_features/](https://prometheus.io/docs/prometheus/latest/disabled_features/) | []string | false |




## PrometheusOperator

PrometheusOperator defines the desired state for some part of the prometheus-operator deployment.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| image | Image to use for a `prometheus-operator` deployment. The `prometheus-operator` makes the Prometheus configuration Kubernetes native and manages and operates Prometheus and Alertmanager clusters. More info: [https://github.com/prometheus-operator/prometheus-operator](https://github.com/prometheus-operator/prometheus-operator) | string | true |
| resources | Resources defines resources requests and limits for single Pods. | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#resourcerequirements-v1-core) | false |
| securityContext | SecurityContext holds pod-level security attributes. | *[SecurityContext](#securitycontext) | false |
| paused | Set paused to reconsilation | bool | false |
| tolerations | Tolerations allow the pods to schedule onto nodes with matching taints. | []v1.Toleration | false |
| nodeSelector | NodeSelector Define which Nodes the Pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\" | map[string]string | false |
| podMonitor | Pod monitor for self monitoring | *[Monitor](#monitor) | false |
| namespaces | Namespaces to scope the interaction of the Prometheus Operator and the apiserver. | string | false |
| labels | Map of string keys and values that can be used to organize and categorize (scope and select) objects. May match selectors of replication controllers and services. More info: [https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/) | map[string]string | false |
| annotations | Annotations is an unstructured key value map stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. They are not queryable and should be preserved when modifying objects. More info: [https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/](https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/) | map[string]string | false |
| serviceAccount | ServiceAccount is a structure which allow specify annotations and labels for Service Account which will use by Alertmanager for work in Kubernetes. Cna be use by external tools to store and retrieve arbitrary metadata. | *[EmbeddedObjectMetadata](#embeddedobjectmetadata) | false |




## PrometheusRule

PrometheusRule handles parameters to override PrometheusRule: alerts of recording rules.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| group |  | string | false |
| alert |  | string | false |
| record |  | string | false |
| for |  | string | false |
| expr |  | string | false |
| severity |  | string | false |




## PrometheusRules

PrometheusRules help to add and override Prometheus rules.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| install |  | *bool | false |
| ruleGroups |  | []string | false |
| override |  | \[\][PrometheusRule](#prometheusrule) | false |




## Promxy

Promxy handles parameters to set up Platform Monitoring with Prometheus proxy.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| install |  | *bool | false |
| port |  | *int32 | false |




## Pushgateway

Pushgateway defines the desired state for some part of pushgateway deployment

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| install | Install indicates is pushgateway will be installed. Can be changed for already deployed service and the service will be removed during next reconciliation iteration | *bool | false |
| image | Image to use for a `pushgateway` deployment. The `pushgateway` is an exporter to collect metrics from VM More info: [https://github.com/prometheus/pushgateway](https://github.com/prometheus/pushgateway) | string | true |
| replicas | Set replicas | *int32 | false |
| extraArgs | Additional pushgateway container arguments. | []string | false |
| volumes | Volumes allows configuration of additional volumes on the output StatefulSet definition. Volumes specified will be appended to other volumes that are generated as a result of StorageSpec objects. More info: [https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#volume-v1-core](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#volume-v1-core) | []v1.Volume | false |
| volumeMounts | VolumeMounts allows configuration of additional VolumeMounts on the output StatefulSet definition. VolumeMounts specified will be appended to other VolumeMounts in the prometheus container, that are generated as a result of StorageSpec objects. More info: [https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#volumemount-v1-core](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#volumemount-v1-core) | []v1.VolumeMount | false |
| storage | PVC spec for Pushgateway. If specified, also adds flags --persistence.file and --persistence.interval with default values, creates volume and volumeMount with name \"storage-volume\" in the deployment. | *v1.PersistentVolumeClaimSpec | false |
| port | Port for `pushgateway` deployment and service | int32 | true |
| ingress | Ingress allows to create Ingress. | *[Ingress](#ingress) | false |
| resources | Resources defines resources requests and limits for single Pods | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#resourcerequirements-v1-core) | false |
| securityContext | SecurityContext holds pod-level security attributes. | *[SecurityContext](#securitycontext) | false |
| nodeSelector | NodeSelector select nodes for deploy | map[string]string | false |
| paused | Set paused to reconsilation | bool | false |
| serviceMonitor | Service monitor for pulling metrics | *[Monitor](#monitor) | false |
| tolerations | Tolerations allow the pods to schedule onto nodes with matching taints. | []v1.Toleration | false |
| labels | Map of string keys and values that can be used to organize and categorize (scope and select) objects. May match selectors of replication controllers and services. More info: [https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/) | map[string]string | false |
| annotations | Annotations is an unstructured key value map stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. They are not queryable and should be preserved when modifying objects. More info: [https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/](https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/) | map[string]string | false |







## SecurityContext

SecurityContext holds pod-level security attributes. The parameters are required if a Pod Security Policy is enabled for
the Kubernetes cluster and required if a Security Context Constraints is enabled for the OpenShift cluster.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| runAsUser | The UID to run the entrypoint of the container process. Defaults to user specified in image metadata if unspecified. | *int64 | false |
| fsGroup | A special supplemental group that applies to all containers in a pod. Some volume types allow the Kubelet to change the ownership of that volume to be owned by the pod:\n\n1. The owning GID will be the FSGroup 2. The setgid bit is set (new files created in the volume will be owned by FSGroup) 3. The permission bits are OR'd with rw-rw  | *int64 | false |




## StackDriverIntegrationConfig

StackDriverIntegrationConfig holds parameters to set up Platform Monitoring integration with Google Cloud Operations (
GCO). Integration schema:\n * Send metrics from Prometheus to GCO by deploying 'stackdriver-prometheus-sidecar'
container\n as sidecar to Prometheus pod. It allows to specify filters for metrics to be sent.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| image | Image of 'stackdriver-prometheus-sidecar'. This service is deploying as sidecar container to Prometheus pod and send metrics from Prometheus to GCO. | string | true |
| projectId | Identificator of project in Google Cloud | string | true |
| location | Location where project is deployed in Google Cloud | string | true |
| cluster | Name of Kubernetes cluster in Google Cloud which will be monitored | string | true |
| metricsFilter | List of filters for metrics which will be sent to GCO. Filters use the same syntax as Prometheus instant vector selectors: [https://prometheus.io/docs/prometheus/latest/querying/basics/#instant-vector-selectors](https://prometheus.io/docs/prometheus/latest/querying/basics/#instant-vector-selectors) | []string | false |




## TLSConfig

TLSConfig extends the safe TLS configuration with file parameters.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| caSecret |  | *[v1.SecretKeySelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#secretkeyselector-v1-core) | false |
| certSecret |  | *[v1.SecretKeySelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#secretkeyselector-v1-core) | false |
| keySecret |  | *[v1.SecretKeySelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#secretkeyselector-v1-core) | false |
| insecureSkipVerify |  | *bool | false |




## Victoriametrics

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| vmOperator |  | [VmOperator](#vmoperator) | false |
| vmSingle |  | [VmSingle](#vmsingle) | false |
| vmAgent |  | [VmAgent](#vmagent) | false |
| vmAlertManager |  | [VmAlertManager](#vmalertmanager) | false |
| vmAlert |  | [VmAlert](#vmalert) | false |




## VmOperator

VmOperator defines the desired state for some part of the victoriametrics-operator deployment.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| install | Install indicates is victoriametrics-operator will be installed. Can be changed for already deployed service and the service will be removed during next reconciliation iteration | *bool | false |
| image | Image to use for a `victoriametrics-operator` deployment. The `victoriametrics-operator` makes the vmoperator configuration Kubernetes native and manages and operates More info: [https://github.com/VictoriaMetrics/operator](https://github.com/VictoriaMetrics/operator) | string | true |
| resources | Resources defines resources requests and limits for single Pods. | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#resourcerequirements-v1-core) | false |
| securityContext | SecurityContext holds pod-level security attributes. | *[SecurityContext](#securitycontext) | false |
| paused | Set paused to reconsilation | bool | false |
| extraEnvs | ExtraEnvs that will be added to VMOperator pod | []v1.EnvVar | false |
| tolerations | Tolerations allow the pods to schedule onto nodes with matching taints. | []v1.Toleration | false |
| nodeSelector | NodeSelector Define which Nodes the Pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\" | map[string]string | false |
| labels | Map of string keys and values that can be used to organize and categorize (scope and select) objects. May match selectors of replication controllers and services. More info: [https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/) | map[string]string | false |
| annotations | Annotations is an unstructured key value map stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. They are not queryable and should be preserved when modifying objects. More info: [https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/) | map[string]string | false |
| serviceMonitor | Service monitor for pulling metrics | *[Monitor](#monitor) | false |




## VmAgent

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| install | Install indicates is vmagent will be installed. Can be changed for already deployed service and the service will be removed during next reconciliation iteration | *bool | false |
| image | Image to use for a `vmagent` deployment. The `victoriametrics-operator` makes the VmAgent configuration Kubernetes native and manages and operates More info: [https://docs.victoriametrics.com/vmalert.html](https://docs.victoriametrics.com/vmalert.html) | string | true |
| resources | Resources defines resources requests and limits for single Pods. | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#resourcerequirements-v1-core) | false |
| securityContext | SecurityContext holds pod-level security attributes. | *[SecurityContext](#securitycontext) | false |
| containers | Containers property allows to inject additions sidecars or to patch existing containers. It can be useful for proxies, backup, etc. | [[]v1.Container](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#container-v1-core) | false |
| paused | Set paused to reconsilation | bool | false |
| tolerations | Tolerations allow the pods to schedule onto nodes with matching taints. | []v1.Toleration | false |
| nodeSelector | NodeSelector Define which Nodes the Pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\" | map[string]string | false |
| labels | Map of string keys and values that can be used to organize and categorize (scope and select) objects. May match selectors of replication controllers and services. More info: [https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/) | map[string]string | false |
| annotations | Annotations is an unstructured key value map stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. They are not queryable and should be preserved when modifying objects. More info: [https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/](https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/) | map[string]string | false |
| ingress | Ingress allows to create Ingress for VM UI. | *[Ingress](#ingress) | false |
| scrapeInterval | ScrapeInterval defines how often scrape targets by default | string | false |
| vmAgentExternalLabelName | VMAgentExternalLabelName Name of vmAgent external label used to denote VmAgent instance name. Defaults to the value of `vmagent`. External label will _not_ be added when value is set to empty string (`\"\"`). | *string | false |
| externalLabels | ExternalLabels The labels to add to any time series scraped by vmagent. it doesn't affect metrics ingested directly by push API's | map[string]string | false |
| remoteWrite | RemoteWrite list of victoria metrics to some other remote write system for vm it must look like: [http://victoria-metrics-single:8429/api/v1/write](http://victoria-metrics-single:8429/api/v1/write) or for cluster different url [https://github.com/VictoriaMetrics/VictoriaMetrics/tree/master/app/vmagent#splitting-data-streams-among-multiple-systems](https://github.com/VictoriaMetrics/VictoriaMetrics/tree/master/app/vmagent#splitting-data-streams-among-multiple-systems) | []vmetricsv1b1.VMAgentRemoteWriteSpec | false |
| remoteWriteSettings | RemoteWriteSettings defines global settings for all remoteWrite urls. | *vmetricsv1b1.VMAgentRemoteWriteSettings | false |
| relabelConfig | RelabelConfig ConfigMap with global relabel config -remoteWrite.relabelConfig This relabeling is applied to all the collected metrics before sending them to remote storage. | *v1.ConfigMapKeySelector | false |
| podMonitorNamespaceSelector | Namespace selector for PodMonitors | *[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#labelselector-v1-meta) | false |
| serviceMonitorNamespaceSelector | Namespace selector for ServiceMonitors | *[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#labelselector-v1-meta) | false |
| podMonitorSelector | Selector for PodMoniotors | *[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#labelselector-v1-meta) | false |
| serviceMonitorSelector | Selector for ServiceMonitors | *[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#labelselector-v1-meta) | false |
| secrets | Secrets is a list of Secrets in the same namespace as the vmagent object, which shall be mounted into the vmagent Pods. will be mounted at path /etc/vm/secrets | []string | false |
| volumes | Volumes allows configuration of additional volumes on the output deploy definition. Volumes specified will be appended to other volumes that are generated as a result of StorageSpec objects. | []v1.Volume | false |
| volumeMounts | VolumeMounts allows configuration of additional VolumeMounts on the output deploy definition. VolumeMounts specified will be appended to other VolumeMounts in the vmagent container, that are generated as a result of StorageSpec objects. | []v1.VolumeMount | false |
| extraArgs | ExtraArgs that will be passed to  VMAgent pod for example remoteWrite.tmpDataPath: /tmp it would be converted to flag --remoteWrite.tmpDataPath=/tmp | map[string]string | false |
| extraEnvs | ExtraEnvs that will be added to VMAgent pod | []v1.EnvVar | false |




## VmSingle

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| install | Install indicates is vmsingle will be installed. Can be changed for already deployed CR will be removed during next reconciliation iteration | *bool | false |
| image | Image to use for a `vmsingle` deployment. The `vmsingle` makes the vmsingle configuration Kubernetes native and manages and operates More info: [https://docs.victoriametrics.com/Single-server-VictoriaMetrics.html](https://docs.victoriametrics.com/Single-server-VictoriaMetrics.html) | string | true |
| resources | Resources defines resources requests and limits for single Pods. | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#resourcerequirements-v1-core) | false |
| securityContext | SecurityContext holds pod-level security attributes. | *[SecurityContext](#securitycontext) | false |
| containers | Containers property allows to inject additions sidecars or to patch existing containers. It can be useful for proxies, backup, etc. | [[]v1.Container](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#container-v1-core) | false |
| paused | Set paused to reconsilation | bool | false |
| tolerations | Tolerations allow the pods to schedule onto nodes with matching taints. | []v1.Toleration | false |
| nodeSelector | NodeSelector Define which Nodes the Pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\" | map[string]string | false |
| labels | Map of string keys and values that can be used to organize and categorize (scope and select) objects. May match selectors of replication controllers and services. More info: [https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/) | map[string]string | false |
| annotations | Annotations is an unstructured key value map stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. They are not queryable and should be preserved when modifying objects. More info: [https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/](https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/) | map[string]string | false |
| ingress | Ingress allows to create Ingress for VmSingle UI. | *[Ingress](#ingress) | false |
| retentionPeriod | RetentionPeriod for the stored metrics Note VictoriaMetrics has data/ and indexdb/ folders metrics from data/ removed eventually as soon as partition leaves retention period reverse index data at indexdb rotates once at the half of configured retention period [https://docs.victoriametrics.com/Single-server-VictoriaMetrics.html#retention](https://docs.victoriametrics.com/Single-server-VictoriaMetrics.html#retention) | string | true |
| extraArgs | ExtraArgs that will be passed to  VMSingle pod for example remoteWrite.tmpDataPath: /tmp | map[string]string | false |
| extraEnvs | ExtraEnvs that will be added to VMSingle pod | []v1.EnvVar | false |
| secrets | Secrets is a list of Secrets in the same namespace as the VMSingle object, which shall be mounted into the VMSingle Pods. | []string | false |
| storage | Storage is the definition of how storage will be used by the VMSingle by default it`s empty dir | *v1.PersistentVolumeClaimSpec | false |
| volumes | Volumes allows configuration of additional volumes on the output deploy definition. Volumes specified will be appended to other volumes that are generated as a result of StorageSpec objects. | []v1.Volume | false |
| volumeMounts | VolumeMounts allows configuration of additional VolumeMounts on the output deploy definition. VolumeMounts specified will be appended to other VolumeMounts in the vmsingle container, that are generated as a result of StorageSpec objects. | []v1.VolumeMount | false |




## VmAlert

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| install | Install indicates is VmAlert will be installed. Can be changed for already deployed service and the service will be removed during next reconciliation iteration | *bool | false |
| image | Image - docker image settings for VMAlert if no specified operator uses default config version | string | true |
| ingress | Ingress allows to create Ingress for VM UI. | *[Ingress](#ingress) | false |
| secrets | Secrets is a list of Secrets in the same namespace as the VMAlert object, which shall be mounted into the VMAlert Pods. The Secrets are mounted into /etc/vm/secrets/\<secret-name\>. | []string | false |
| replicas | ReplicaCount is the expected size of the VMAlert cluster. The controller will eventually make the size of the running cluster equal to the expected size. | *int32 | false |
| resources | Resources container resource request and limits, [https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/](https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/) | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#resourcerequirements-v1-core) | false |
| affinity | Affinity If specified, the pod's scheduling constraints. | *v1.Affinity | false |
| tolerations | Tolerations If specified, the pod's tolerations. | []v1.Toleration | false |
| securityContext | SecurityContext holds pod-level security attributes and common container settings. This defaults to the default PodSecurityContext. | *v1.PodSecurityContext | false |
| containers | Containers property allows to inject additions sidecars or to patch existing containers. It can be useful for proxies, backup, etc. | [[]v1.Container](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#container-v1-core) | false |
| evaluationInterval | EvaluationInterval how often evalute rules by default. Pattern:=\"[0-9]+(ms\|s\|m\|h) | string | false |
| selectAllByDefault | SelectAllByDefault changes default behavior for empty CRD selectors, such RuleSelector. with selectAllByDefault: true and empty serviceScrapeSelector and RuleNamespaceSelector Operator selects all exist serviceScrapes with selectAllByDefault: false - selects nothing | bool | false |
| ruleSelector | RuleSelector selector to select which VMRules to mount for loading alerting rules from. Works in combination with NamespaceSelector. If both nil - behaviour controlled by selectAllByDefault NamespaceSelector nil - only objects at VMAlert namespace. | *[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#labelselector-v1-meta) | false |
| ruleNamespaceSelector | RuleNamespaceSelector to be selected for VMRules discovery. Works in combination with Selector. If both nil - behaviour controlled by selectAllByDefault NamespaceSelector nil - only objects at VMAlert namespace. | *[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#labelselector-v1-meta) | false |
| port | Port for listen | string | false |
| remoteWrite | RemoteWrite Optional URL to remote-write compatible storage to persist vmalert state and rule results to. Rule results will be persisted according to each rule. Alerts state will be persisted in the form of time series named ALERTS and ALERTS_FOR_STATE see -remoteWrite.url docs in vmalerts for details. E.g. `http://127.0.0.1:8428` | *vmetricsv1b1.VMAlertRemoteWriteSpec | false |
| remoteRead | RemoteRead Optional URL to read vmalert state (persisted via RemoteWrite) This configuration only makes sense if alerts state has been successfully persisted (via RemoteWrite) before. see -remoteRead.url docs in vmalerts for details. E.g. `http://127.0.0.1:8428` | *vmetricsv1b1.VMAlertRemoteReadSpec | false |
| rulePath | RulePath to the file with alert rules. Supports patterns. Flag can be specified multiple times. Examples: -rule /path/to/file. Path to a single file with alerting rules -rule dir/_*_.yaml -rule /_*_.yaml. Relative path to all .yaml files in folder, absolute path to all .yaml files in root. by default operator adds /etc/vmalert/configs/base/vmalert.yaml | []string | false |
| notifier | Notifier prometheus alertmanager endpoint spec. Required at least one of  notifier or notifiers. e.g. `http://127.0.0.1:9093` If specified both notifier and notifiers, notifier will be added as last element to notifiers. | *vmetricsv1b1.VMAlertNotifierSpec | false |
| datasource | Datasource Victoria Metrics or VMSelect url. Required parameter. e.g. `http://127.0.0.1:8428` | vmetricsv1b1.VMAlertDatasourceSpec | true |
| extraArgs | ExtraArgs that will be passed to  VMAlert pod for example -remoteWrite.tmpDataPath=/tmp | map[string]string | false |
| extraEnvs | ExtraEnvs that will be added to VMAlert pod | []v1.EnvVar | false |
| externalLabels | ExternalLabels in the form 'name: value' to add to all generated recording rules and alerts. | map[string]string | false |
| nodeSelector | NodeSelector Define which Nodes the Pods are scheduled on. | map[string]string | false |
| terminationGracePeriodSeconds | TerminationGracePeriodSeconds period for container graceful termination | *int64 | false |
| paused | Set paused to reconsilation | bool | false |
| labels | Map of string keys and values that can be used to organize and categorize (scope and select) objects. May match selectors of replication controllers and services. More info: [https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/) | map[string]string | false |
| annotations | Annotations is an unstructured key value map stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. They are not queryable and should be preserved when modifying objects. More info: [https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/](https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/) | map[string]string | false |
| volumes | Volumes allows configuration of additional volumes on the output Deployment definition. Volumes specified will be appended to other volumes that are generated as a result of StorageSpec objects. | []v1.Volume | false |
| volumeMounts | VolumeMounts allows configuration of additional VolumeMounts on the output Deployment definition. VolumeMounts specified will be appended to other VolumeMounts in the VMAlert container, that are generated as a result of StorageSpec objects. | []v1.VolumeMount | false |




## VmAlertManager

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| install | Install indicates is AlertManager will be installed. Can be changed for already deployed service and the service will be removed during next reconciliation iteration | *bool | false |
| image | Image to use for a `AlertManager` deployment. The `AlertManager` is alerting system which read metrics from Prometheus More info: [https://prometheus.io/docs/alerting/alertmanager/](https://prometheus.io/docs/alerting/alertmanager/) | string | true |
| ingress | Ingress allows to create Ingress for VM UI. | *[Ingress](#ingress) | false |
| secrets | Secrets is a list of Secrets in the same namespace as the VMAlertmanager object, which shall be mounted into the VMAlertmanager Pods. The Secrets are mounted into /etc/vm/secrets/\<secret-name\> | []string | false |
| configRawYaml | ConfigRawYaml - raw configuration for alertmanager, it helps it to start without secret. priority -> hardcoded ConfigRaw -> ConfigRaw, provided by user -> ConfigSecret. | string | false |
| configSecret | ConfigSecret is the name of a Kubernetes Secret in the same namespace as the VMAlertmanager object, which contains configuration for this VMAlertmanager, configuration must be inside secret key: alertmanager.yaml. It must be created by user. instance. Defaults to 'vmalertmanager-\<alertmanager-name\>' The secret is mounted into /etc/alertmanager/config. | string | false |
| replicas | ReplicaCount Size is the expected size of the alertmanager cluster. The controller will eventually make the size of the running cluster equal to the expected | *int32 | false |
| retention | Retention Time duration VMAlertmanager shall retain data for. Default is '120h', and must match the regular expression `[0-9]+(ms\|s\|m\|h)` (milliseconds seconds minutes hours). | string | false |
| paused | Set paused to reconsilation | bool | false |
| nodeSelector | NodeSelector Define which Nodes the Pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\" | map[string]string | false |
| resources | Resources defines resources requests and limits for single Pods. | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#resourcerequirements-v1-core) | false |
| affinity | Affinity If specified, the pod's scheduling constraints. | *v1.Affinity | false |
| tolerations | Tolerations allow the pods to schedule onto nodes with matching taints. | []v1.Toleration | false |
| securityContext | SecurityContext holds pod-level security attributes. | *[SecurityContext](#securitycontext) | false |
| containers | Containers allows injecting additional containers or patching existing containers. This is meant to allow adding an authentication proxy to an VMAlertmanager pod. | [[]v1.Container](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#container-v1-core) | false |
| selectAllByDefault | SelectAllByDefault changes default behavior for empty CRD selectors, such ConfigSelector. with selectAllScrapes: true and undefined ConfigSelector and ConfigNamespaceSelector Operator selects all exist alertManagerConfigs with selectAllScrapes: false - selects nothing | bool | false |
| configSelector | ConfigSelector defines selector for VMAlertmanagerConfig, result config will be merged with with Raw or Secret config. Works in combination with NamespaceSelector. NamespaceSelector nil - only objects at VMAlertmanager namespace. Selector nil - only objects at NamespaceSelector namespaces. If both nil - behaviour controlled by selectAllByDefault | *[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#labelselector-v1-meta) | false |
| configNamespaceSelector | ConfigNamespaceSelector defines namespace selector for VMAlertmanagerConfig. Works in combination with Selector. NamespaceSelector nil - only objects at VMAlertmanager namespace. Selector nil - only objects at NamespaceSelector namespaces. If both nil - behaviour controlled by selectAllByDefault | *[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#labelselector-v1-meta) | false |
| terminationGracePeriodSeconds | TerminationGracePeriodSeconds period for container graceful termination | *int64 | false |
| labels | Map of string keys and values that can be used to organize and categorize (scope and select) objects. May match selectors of replication controllers and services. More info: [https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/) | map[string]string | false |
| annotations | Annotations is an unstructured key value map stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. They are not queryable and should be preserved when modifying objects. More info: [https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/](https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/) | map[string]string | false |
| storage | Storage is the definition of how storage will be used by the VMAlertmanager instances. | *vmetricsv1b1.StorageSpec | false |
| volumes | Volumes allows configuration of additional volumes on the output deploy definition. Volumes specified will be appended to other volumes that are generated as a result of StorageSpec objects. | []v1.Volume | false |
| volumeMounts | VolumeMounts allows configuration of additional VolumeMounts on the output deploy definition. VolumeMounts specified will be appended to other VolumeMounts in the alertmanager container, that are generated as a result of StorageSpec objects. | []v1.VolumeMount | false |
| extraArgs | ExtraArgs that will be passed to  VMAlertmanager pod for example log.level: debug | map[string]string | false |
| extraEnvs | ExtraEnvs that will be added to VMAlertmanager pod | []v1.EnvVar | false |



