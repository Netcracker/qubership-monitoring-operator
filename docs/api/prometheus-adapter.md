<!-- markdownlint-disable MD013 -->
This section describes the types introduced by the Prometheus Adapter Operator.


> Note this document is generated from code comments. When contributing a change to this document please do so by changing the code comments.


## PrometheusAdapter

PrometheusAdapter is the Schema for the prometheusadapters API.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#objectmeta-v1-meta) | false |
| spec |  | [PrometheusAdapterSpec](#prometheusadapterspec) | false |
| status |  | PrometheusAdapterStatus | false |




## PrometheusAdapterList

PrometheusAdapterList contains a list of PrometheusAdapter.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata |  | [metav1.ListMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#listmeta-v1-meta) | false |
| items |  | [PrometheusAdapter](#prometheusadapter) | true |




## PrometheusAdapterSpec

PrometheusAdapterSpec defines the desired state of PrometheusAdapter.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| image | Image to use for a `prometheus-adapter` deployment. | string | true |
| prometheusUrl | PrometheusURL used to connect to any tool with Prometheus compatible API. It will eventually contain query parameters to configure the connection. | string | false |
| metricsRelistInterval | MetricsRelistInterval is the interval at which to update the cache of available metrics from Prometheus | string | false |
| enableResourceMetrics | Enable adapter for `metrics.k8s.io`. By default - `false` | boolean | false |
| enableCustomMetrics   | Enable adapter for `custom.metrics.k8s.io`. By default - `true` | boolean | false |
| customScaleMetricRulesSelector | CustomScaleMetricRulesSelector defines label selectors to select CustomScaleMetricRule resources across the cluster. | []*[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#labelselector-v1-meta) | false |
| resources | Resources defines resources requests and limits for single Pods. | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#resourcerequirements-v1-core) | false |
| securityContext | SecurityContext holds pod-level security attributes. | *[SecurityContext](#securitycontext) | false |
| nodeSelector | Define which Nodes the Pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\" | map[string]string | false |
| labels | Map of string keys and values that can be used to organize and categorize (scope and select) objects. May match selectors of replication controllers and services. More info: [https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/) | map[string]string | false |
| annotations | Annotations is an unstructured key value map stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. They are not queryable and should be preserved when modifying objects. More info: [https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/](https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/) | map[string]string | false |
| tolerations | Tolerations allow the pods to schedule onto nodes with matching taints. More info: [https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration](https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration) | []v1.Toleration | false |




## SecurityContext

SecurityContext holds pod-level security attributes. The parameters are required if a Pod Security Policy is enabled for
the Kubernetes cluster and required if a Security Context Constraints is enabled for the OpenShift cluster.

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| runAsUser | The UID to run the entrypoint of the container process. Defaults to user specified in image metadata if unspecified. | *int64 | false |
| fsGroup | A special supplemental group that applies to all containers in a pod. Some volume types allow the Kubelet to change the ownership of that volume to be owned by the pod:\n\n1. The owning GID will be the FSGroup 2. The setgid bit is set (new files created in the volume will be owned by FSGroup) 3. The permission bits are OR'd with rw-rw  | *int64 | false |


