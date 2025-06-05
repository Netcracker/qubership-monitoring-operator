### node-exporter

<!-- markdownlint-disable line-length -->
| Field                      | Description                                                                                                                                                                                                            | Scheme                                                                                                                       |
| -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| install                    | Allows to disable create node-exporter during the deployment.                                                                                                                                                          | bool                                                                                                                         |
| paused                     | Set paused to reconciliation.                                                                                                                                                                                          | bool                                                                                                                         |
| image                      | A docker image to be used for the node-exporter deployment.                                                                                                                                                            | string                                                                                                                       |
| setupSecurityContext       | Allows to create PodSecurityPolicy or SecurityContextConstraints.                                                                                                                                                      | string                                                                                                                       |
| port                       | The port for node-exporter daemonset and service.                                                                                                                                                                      | int                                                                                                                          |
| resources                  | The resources that describe the compute resource requests and limits for single pods.                                                                                                                                  | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcerequirements-v1-core) |
| securityContext            | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                     | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#podsecuritycontext-v1-core)    |
| tolerations                | Tolerations allow the pods to schedule onto nodes with matching taints.                                                                                                                                                | []v1.Toleration                                                                                                              |
| affinity                   | If specified, the pod's scheduling constraints                                                                                                                                                                         | *v1.Affinity                                                                                                                 |
| annotations                | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value" | map[string]string                                                                                                            |
| labels                     | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                       | map[string]string                                                                                                            |
| collectorTextfileDirectory | Directory for textfile. For more information, refer to [https://github.com/prometheus/node_exporter#textfile-collector](https://github.com/prometheus/node_exporter#textfile-collector)                                | string                                                                                                                       |
| extraArgs                  | Additional arguments for node-exporter container. For example: "--collector.systemd".                                                                                                                                  | list[string]                                                                                                                 |
| priorityClassName          | PriorityClassName assigned to the Pods to prevent them from evicting.                                                                                                                                                  | string                                                                                                                       |
| serviceMonitor             | Service monitor configuration for pulling metrics.                                                                                                                                                                     | [Monitor](#monitor)                                                                                                          |
<!-- markdownlint-enable line-length -->

Example:

```yaml
nodeExporter:
  install: true
  paused: false
  setupSecurityContext: false
  image: prom/node-exporter:v0.18.1
  port: 19100
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
  tolerations: []
  affinity: {}
  annotations:
    annotation.key: annotation-value
  priorityClassName: priority-class
  collectorTextfileDirectory: /var/spool/monitoring
  extraArgs:
    - --collector.systemd
  
  serviceMonitor:
    interval: 60s
    ...
```


