### kube-state-metrics

<!-- markdownlint-disable line-length -->
| Field                 | Description                                                                                                                                                                                                                                                     | Scheme                                                                                                                       |
| --------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| install               | Allows to disable create kube-state-metrics during the deployment.                                                                                                                                                                                              | bool                                                                                                                         |
| paused                | Set paused to reconciliation.                                                                                                                                                                                                                                   | bool                                                                                                                         |
| image                 | A docker image to be used for the kube-state-metrics deployment.                                                                                                                                                                                                | string                                                                                                                       |
| resources             | The resources that describe the compute resource requests and limits for single pods.                                                                                                                                                                           | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcerequirements-v1-core) |
| securityContext       | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                                                              | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#podsecuritycontext-v1-core)    |
| affinity                                            | If specified, the pod's scheduling constraints                                                                                                                                                                                      | *v1.Affinity                                                                                                                   |
| annotations           | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value"                                          | map[string]string                                                                                                            |
| labels                | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                                                                | map[string]string                                                                                                            |
| priorityClassName     | PriorityClassName assigned to the Pods to prevent them from evicting.                                                                                                                                                                                           | string                                                                                                                       |
| namespaces            | Comma separated list of namespaces to monitor in non-privileged mode. This parameter is unnecessary if the kube-state-metrics has ClusterRole on all namespaces                                                                                                 | string                                                                                                                       |
| scrapeResources       | Comma-separated list of Resources to be enabled. Empty means that metrics will collect for all type of resources                                                                                                                                                | string                                                                                                                       |
| metricLabelsAllowlist | Comma-separated list of additional Kubernetes label keys that will be used in the resource labels metric. Default value is `nodes=[*],pods=[*],namespaces=[*],deployments=[*],statefulsets=[*],daemonsets=[*],cronjobs=[*],jobs=[*],ingresses=[*],services=[*]` | string                                                                                                                       |
<!-- markdownlint-enable line-length -->

Example:

```yaml
kubeStateMetrics:
  install: true
  paused: false
  image: coreos/kube-state-metrics:latest
  resources:
    limits:
      cpu: 200m
      memory: 200Mi
    requests:
      cpu: 100m
      memory: 100Mi
  securityContext:
    runAsUser: 2000
    fsGroup: 2000
  labels:
    label.key: label-value
  annotations:
    annotation.key: annotation-value
  priorityClassName: priority-class

  namespaces: monitoring,logging,tracing
  scrapeResources: configmaps,cronjobs,daemonsets,deployments,endpoints,jobs,limitranges,persistentvolumeclaims,poddisruptionbudgets,namespaces,nodes,pods,persistentvolumes,replicasets,replicationcontrollers,resourcequotas,services,statefulsets
  metricLabelsAllowlist: nodes=[*],pods=[*],namespaces=[*],deployments=[*],statefulsets=[*],daemonsets=[*],cronjobs=[*],jobs=[*],ingresses=[*],services=[*]
```


