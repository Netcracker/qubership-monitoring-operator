This document describes how to configure Monitoring to collect metrics (user metrics or all metrics) from a specific
list of namespaces.

# Types of metrics

All metrics collected by Monitoring can be classified into next categories:

* Kubernetes node metrics (VM specific, like node CPU, Memory, FS, Disk)
* Kubernetes components metrics (metrics from `kube-apiserver`, `etcd`, and so on)
* Kubernetes microservices metrics (metrics about pods, containers, deployments and so on)
* Metrics about microservice runtime (like JVM, and Golang runtime)
* Microservice metrics (how microservice do any work)
* Custom metrics (for example, from any VMs outside the Cloud)

This classification needs to answer on question of which metrics can be limited by namespace and which not.

Obviously, Kubernetes node metrics can't be limited by namespace name because they are related to node and different
pods from one namespace can be run on different nodes.

Other metrics can be limited or filtered by namespace. But absent some metrics can affect dashboards or default alerts.
For example, if disable Kubernetes microservices metrics in the specific namespace
(from `kubelet` and/or `kube-state-metrics`) it will lead to a situation when you can't check pods resource usage.

So please think twice before disabling any metrics collection.

# How to configure filters for Monitoring CRs and metrics collection

## How configured filters will apply

All our operators allow us to configure two types of selectors:

* namespace selector
* label selector

Namespace selector allows to configure in which namespaces the operator will discover this type of Custom Resource
using differing criteria.

For example, the namespace selector can be used to filter the list of using namespaces like the following:

* filter by namespace name, using the special meta label `kubernetes.io/metadata.name`
* filter by namespace labels, using any labels or their keys

Selectors can used together or separately. In case, if you want to use them together they have an order
in which they will apply and work:

1. Firstly, will apply a namespace selector and select a scope of namespaces where the operator will discover
2. Second, will apply a label selector and the operator will select Custom Resources using specified filters
   (from the label selector)

If the namespace of the label selector doesn't have any conditions the operator will discover all
Custom Resources (CRs).

For example, if the namespace selector is empty the operator will discover CRs in all namespaces using to
filter label selector. If the label selector is empty, so the operator will discover all CRs using the configured
namespace selector.

## Filters by labels

The Monitoring allows set label selectors that will apply during the discovery for
`ServiceMonitor`/`PodMonitor`/`Probe` CRs.

To do it need to specify:

```yaml
victoriametrics:
  vmAgent:
    # You can use matchLabels of matchExpressions
    serviceMonitorLabelSelector:
      # Example of matchLabels
      matchLabels:
        <key>: <value>
    podMonitorLabelSelector:
      # Example of matchExpressions
      matchExpressions:
        - key: <key>
          operator: In
          values:
            - <value>
    probeLabelSelector:
      # Example of matchLabels
      matchLabels:
        <key>: <value>
```

For example:

```yaml
victoriametrics:
  vmAgent:
    serviceMonitorLabelSelector:
      matchLabels:
        monitoring.qubership.org/enable: true
    podMonitorLabelSelector:
      matchExpressions:
        - key: monitoring.qubership.org/enable
          operator: In
          values:
            - true
    probeLabelSelector:
      matchLabels:
        monitoring.qubership.org/enable: true
```

**Warning!** Here you may want to use multiple selectors, to select CRs using two or more different labels.
For example:

```yaml
victoriametrics:
  vmAgent:
    serviceMonitorLabelSelector:
      matchLabels:
        monitoring.qubership.org/enable: true
        app.kubernetes.io/component: monitoring
```

But these conditions work using **AND** logic (not OR logic). The conditions from the example about means:

* Need to discover all CRs with a type of `ServiceMonitor`
* Filter and keep only CRs that have labels
  `monitoring.qubership.org/enable: true` AND `app.kubernetes.io/component: monitoring`

How to use matchLabels and `matchExpressions` you can read:

* [https://kubernetes.io/docs/concepts/overview/working-with-objects/labels](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels)
* [https://riyafa.wordpress.com/2020/06/07/kubernetes-matchexpressions-explained/](https://riyafa.wordpress.com/2020/06/07/kubernetes-matchexpressions-explained/)
* [https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#labelselector-v1-meta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#labelselector-v1-meta)

There are two approaches how to use label selectors:

* White list
  * Specify the label/labels that must be set on each CRs
  * Specify the label/labels key that must be set on each CRs
* Black list
  * Specify the label/labels that shouldn't be on each CRs
  * Specify the label/labels key that shouldn't be set on each CRs

**Warning!** Do not forget to add the namespace name of the namespace where deployed Monitoring or label to your
white or do NOT add to the black list. Almost all configurations to collect default metrics are created
in the Monitoring namespace.

## Filters by namespaces

The Monitoring allows set namespace selectors that will apply during the discovery for
`ServiceMonitor`/`PodMonitor`/`Probe` CRs.

To do it need to specify:

```yaml
victoriametrics:
  vmAgent:
    # You can use matchLabels of matchExpressions
    serviceMonitorNamespaceSelector:
      # Example of matchLabels
      matchLabels:
        <key>: <value>
    podMonitorNamespaceSelector:
      # Example of matchExpressions
      matchExpressions:
        - key: <key>
          operator: In
          values:
            - <value>
    probeNamespaceSelector:
      # Example of matchLabels
      matchLabels:
        <key>: <value>
```

For example:

```yaml
victoriametrics:
  vmAgent:
    serviceMonitorNamespaceSelector:
      matchLabels:
        monitoring.qubership.org/enable: true
    podMonitorNamespaceSelector:
      matchExpressions:
        - key: monitoring.qubership.org/enable
          operator: In
          values:
            - true
    probeNamespaceSelector:
      matchLabels:
        monitoring.qubership.org/enable: true
```

**Warning!** Here you may want to use multiple selectors, to select CRs using two or more different labels.
For example:

```yaml
victoriametrics:
  vmAgent:
    serviceMonitorNamespaceSelector:
      matchLabels:
        monitoring.qubership.org/enable: true
        app.kubernetes.io/component: monitoring
```

But these conditions work using AND logic (not OR logic). The conditions from the example about means:

* Need to all namespaces to discover CRs with a type of `ServiceMonitor`
* Filter and keep only CRs from namespaces that have labels
  `monitoring.qubership.org/enable: true` AND `app.kubernetes.io/component: monitoring`

How to use matchLabels and `matchExpressions` you can read:

* [https://kubernetes.io/docs/concepts/overview/working-with-objects/labels](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels)
* [https://riyafa.wordpress.com/2020/06/07/kubernetes-matchexpressions-explained/](https://riyafa.wordpress.com/2020/06/07/kubernetes-matchexpressions-explained/)
* [https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#labelselector-v1-meta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#labelselector-v1-meta)

There are two approaches how to use namespaceSelectors:

* White list
  * Specify the label selector that must be set on the namespace labels
  * Specify the list of namespace names in namespaceSelector
* Black list
  * Specify the label selector that shouldn't be present in namespace labels
  * Specify the list of namespace names in namespaceSelector

**Warning!** Do not forget to add the namespace name of the namespace where deployed Monitoring or label to your
white or do NOT add to the black list. Almost all configurations to collect default metrics are created
in the Monitoring namespace.

## Disable metric collection for default metric sources

This section describes how to disable metrics collection or filter by namespace default metrics.

### kube-state-metrics

This exporter provides metrics about Kubernetes and objects in it. By default, it collects metrics from namespaces.

If you want to limit it and specify a list of namespaces from which it should collect metrics you
should use the parameter:

```yaml
kubeStateMetrics:
  # Comma-separated values
  namespaces: <name1>,<name2>,...,<nameN>
```

**Warning!** If you set a list of namespaces, the kube-state-metrics won't discover new namespaces
and will work only with a list of specified labels. Also, you will lose all metrics about Deployment,
DaemonSet, Pod, Ingresses, and other Kubernetes objects from namespaces not in the list.

### kubelet

You can't exclude any namespace from the metrics collection. You can only add relabeling rules to drop metrics
for specific namespaces.

How to do it you can read:

* [Monitoring: Kubernetes Monitors](../installation/README.md#kubernetes-monitors)
* [https://prometheus.io/docs/prometheus/latest/configuration/configuration/#scrape_config](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#scrape_config)

For example:

```yaml
kubernetesMonitors:
  kubeletServiceMonitor:
    install: true
    interval: 30s
    scrapeTimeout: 10s
    metricRelabelings:
      - action: drop
        regex: <your_namespace_name>
        sourceLabels:
          - __meta_kubernetes_namespace
```

**Warning!** If you configure relabeling it will drop all metrics about Pod and Container runtime, like CPU, memory,
and filesystem usage, for all pods from configured namespaces.

### kube-apiserver

You can't exclude any namespace from the metrics collection because all metrics are produced with
`namespace=default` label.

### etcd

You can't exclude any namespace from the metrics collection because all metrics are produced with
`namespace=<etcd_namespace>` label.

### node-exporter

You can't exclude any namespace from the metrics collection because it metrics have no `namespace` label.

## Examples

This section provide different examples how user can configure namespace and label selectors.

### Whitelist by the namespace name

**Warning!** Do not forget to add the namespace name of the namespace where deployed Monitoring or label to your
white or do NOT add to the black list. Almost all configurations to collect default metrics are created
in the Monitoring namespace.

For example, let's offer that you have the following namespaces:

```yaml
kind: Namespace
apiVersion: v1
metadata:
  name: opensearch
  labels:
    kubernetes.io/metadata.name: opensearch
```

and

```yaml
kind: Namespace
apiVersion: v1
metadata:
  name: cpq
  labels:
    kubernetes.io/metadata.name: cpq
```

and the namespace where installed Monitoring:

```yaml
kind: Namespace
apiVersion: v1
metadata:
  name: monitoring
  labels:
    kubernetes.io/metadata.name: monitoring
```

So to disable discovery Monitoring CRs for namespace `opensearch` (and metrics collections from all microservices
in this namespace) you can configure the whitelist using which Monitoring will discover ServiceMonitor/PodMonitors
**only** in `cpq` and `monitoring` namespaces:

```yaml
victoriametrics:
  vmAgent:
    serviceMonitorNamespaceSelector:
      matchExpressions:
        - key: kubernetes.io/metadata.name
          operator: In
          values:
            - cqp
            - monitoring
    podMonitorNamespaceSelector:
      matchExpressions:
        - key: kubernetes.io/metadata.name
          operator: In
          values:
            - cpq
            - monitoring
```

### Blacklist by the namespace name

For example, let's offer that you have the following namespaces:

```yaml
kind: Namespace
apiVersion: v1
metadata:
  name: opensearch
  labels:
    kubernetes.io/metadata.name: opensearch
```

and

```yaml
kind: Namespace
apiVersion: v1
metadata:
  name: cpq
  labels:
    kubernetes.io/metadata.name: cpq
```

So to disable discovery Monitoring CRs for namespace `opensearch` (and metrics collections from all microservices
in this namespace) you can configure the blacklist:

```yaml
victoriametrics:
  vmAgent:
    serviceMonitorNamespaceSelector:
      matchExpressions:
        - key: kubernetes.io/metadata.name
          operator: NotIn
          values:
            - opensearch
    podMonitorNamespaceSelector:
      matchExpressions:
        - key: kubernetes.io/metadata.name
          operator: NotIn
          values:
            - opensearch
```

### Whitelist by the special label

Let's assume that you want to configure Monitoring to collect metrics from all namespaces except some with specific
labels. For example with the label:

```yaml
monitoring.qubership.org/enable: true
```

**Note:** You can any label key and value or even some label combinations that you want. But do not use any common
or often used key pairs to avoid intersection with common labels.

The selector for VMAgent will look as follows if you want to check the label value:

```yaml
victoriametrics:
  vmAgent:
    serviceMonitorNamespaceSelector:
      matchExpressions:
        - key: monitoring.qubership.org/enable
          operator: In
          values:
            - true
    podMonitorNamespaceSelector:
      matchExpressions:
        - key: monitoring.qubership.org/enable
          operator: In
          values:
            - true
```

or like as follows if you want to check only the label key presenting in namespace labels:

```yaml
victoriametrics:
  vmAgent:
    serviceMonitorNamespaceSelector:
      matchExpressions:
        - key: monitoring.qubership.org/enable
          operator: Exists
    podMonitorNamespaceSelector:
      matchExpressions:
        - key: monitoring.qubership.org/enable
          operator: Exists
```

**Warning!** Do not forget to add the label of the namespace where deployed Monitoring or label to your
white or do NOT add to the black list. Almost all configurations to collect default metrics are created
in the Monitoring namespace.

Next, you should add the selected label on namespaces that should be ignored:

```yaml
kind: Namespace
apiVersion: v1
metadata:
  name: opensearch
  labels:
    kubernetes.io/metadata.name: opensearch
    monitoring.qubership.org/enable: true
```

and

```yaml
kind: Namespace
apiVersion: v1
metadata:
  name: monitoring
  labels:
    kubernetes.io/metadata.name: monitoring
    monitoring.qubership.org/enable: true
```

### Blacklist by the special label

Let's assume that you want to configure Monitoring to collect metrics from all namespaces except some with specific
labels. For example with the label:

```yaml
monitoring.qubership.org/disable: true
```

**Note:** You can any label key and value or even some label combinations that you want. But do not use any common
or often used key pairs to avoid intersection with common labels.

The selector for VMAgent will look as follows if you want to check the label value:

```yaml
victoriametrics:
  vmAgent:
    serviceMonitorNamespaceSelector:
      matchExpressions:
        - key: monitoring.qubership.org/disable
          operator: NotIn
          values:
            - true
    podMonitorNamespaceSelector:
      matchExpressions:
        - key: monitoring.qubership.org/disable
          operator: NotIn
          values:
            - true
```

or like as follows if you want to check only the label key presenting in namespace labels:

```yaml
victoriametrics:
  vmAgent:
    serviceMonitorNamespaceSelector:
      matchExpressions:
        - key: monitoring.qubership.org/disable
          operator: DoesNotExist
    podMonitorNamespaceSelector:
      matchExpressions:
        - key: monitoring.qubership.org/disable
          operator: DoesNotExist
```

Next, you should add the selected label on namespaces that should be ignored:

```yaml
kind: Namespace
apiVersion: v1
metadata:
  name: opensearch
  labels:
    kubernetes.io/metadata.name: opensearch
    monitoring.qubership.org/disable: true
```
