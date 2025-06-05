### pushgateway

<!-- markdownlint-disable line-length -->
| Field             | Description                                                                                                                                                                                                                                                                                                                                                | Scheme                                                                                                                                 |
| ----------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------- |
| install           | Allow to disable create Pushgateway during deploy.                                                                                                                                                                                                                                                                                                         | boolean                                                                                                                                |
| image             | Image of pushgateway.                                                                                                                                                                                                                                                                                                                                      | string                                                                                                                                 |
| replicas          | Number of created pods.                                                                                                                                                                                                                                                                                                                                    | int                                                                                                                                    |
| paused            | Set paused to reconciliation.                                                                                                                                                                                                                                                                                                                              | boolean                                                                                                                                |
| extraArgs         | Additional pushgateway container arguments.                                                                                                                                                                                                                                                                                                                | list[string]                                                                                                                           |
| port              | Port for pushgateway deployment and service.                                                                                                                                                                                                                                                                                                               | integer                                                                                                                                |
| serviceMonitor    | Service monitor for pulling metrics.                                                                                                                                                                                                                                                                                                                       | *[Monitor](#monitor)                                                                                                                   |
| ingress           | Ingress allows to create Ingress.                                                                                                                                                                                                                                                                                                                          | *[Ingress](#ingress)                                                                                                                   |
| resources         | The resources that describe the compute resource requests and limits for single pods.                                                                                                                                                                                                                                                                      | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcerequirements-v1-core)           |
| nodeSelector      | Defines which nodes the pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\"                                                                                                                                                                                                                                         | map[string]string                                                                                                                      |
| securityContext   | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                                                                                                                                                         | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#podsecuritycontext-v1-core)              |
| tolerations       | Tolerations allow the pods to schedule onto nodes with matching taints.                                                                                                                                                                                                                                                                                    | []v1.Toleration                                                                                                                        |
| affinity                                            | If specified, the pod's scheduling constraints                                                                                                                                                                                      | *v1.Affinity                                                                                                                   |
| annotations       | Annotations is an unstructured key value map stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. They are not queryable and should be preserved when modifying objects. For more information, refer to [https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/](https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/) | map[string]string                                                                                                                      |
| labels            | The map of string keys and values that can be used to organize and categorize (scope and select) objects. May match selectors of replication controllers and services. For more information, refer to [https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/)                                                           | map[string]string                                                                                                                      |
| volumes           | Volumes allows configuration of additional volumes on the output StatefulSet definition. Volumes specified will be appended to other volumes that are generated as a result of StorageSpec objects.                                                                                                                                                        | [v1.Volume](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#volume-v1-core)                                       |
| volumeMounts      | VolumeMounts allows configuration of additional VolumeMounts on the output StatefulSet definition. VolumeMounts specified will be appended to other VolumeMounts in the prometheus container, that are generated as a result of StorageSpec objects.                                                                                                       | [v1.VolumeMount](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#volumemount-v1-core)                             |
| storage           | PVC spec for Pushgateway. If specified, also adds flags --persistence.file=/data/pushgateway.data and --persistence.interval=5m, creates volume and volumeMount with name "storage-volume" in the deployment.                                                                                                                                              | [v1.PersistentVolumeClaimSpec](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#persistentvolumeclaimspec-v1-core) |
| priorityClassName | PriorityClassName assigned to the Pods to prevent them from evicting.                                                                                                                                                                                                                                                                                      | string                                                                                                                                 |
<!-- markdownlint-enable line-length -->

Parameter `extraArgs` can contain any flags that Pushgateway can handle.
You can find them in the table bellow (relevant for Pushgateway v1.4.1).

<!-- markdownlint-disable line-length -->
| Flag                               | Default  | Description                                                                                    |
| ---------------------------------- | -------- | ---------------------------------------------------------------------------------------------- |
| `-h`, `--help`                     |          | Show context-sensitive help.                                                                   |
| `--version`                        |          | Show application version.                                                                      |
| `--web.listen-address`             | :9091    | Address to listen on for the web interface, API, and telemetry.                                |
| `--web.telemetry-path`             | /metrics | Path under which to expose metrics.                                                            |
| `--web.external-url`               |          | The URL under which the Pushgateway is externally reachable.                                   |
| `--web.route-prefix`               |          | Prefix for the internal routes of web endpoints. Defaults to the path of `--web.external-url`. |
| `--web.enable-lifecycle`           | false    | Enable shutdown via HTTP request.                                                              |
| `--web.enable-admin-api`           | false    | Enable API endpoints for admin control actions.                                                |
| `--persistence.file`               |          | File to persist metrics. If empty, metrics are only kept in memory.                            |
| `--persistence.interval`           | 5m       | The minimum interval at which to write out the persistence file.                               |
| `--push.disable-consistency-check` | false    | Do not check consistency of pushed metrics. DANGEROUS.                                         |
| `--log.level`                      | info     | Only log messages with the given severity or above. One of: [debug, info, warn, error].        |
| `--log.format`                     | logfmt   | Output format of log messages. One of: [logfmt, json].                                         |
<!-- markdownlint-enable line-length -->

Example:

```yaml
pushgateway:
  install: true
  image: prom/pushgateway:v1.4.1
  replicas: 1
  paused: false
  extraArgs:
    - "--log.level=info"
  volumes: {}
  volumeMounts: {}
  storage: {}
  port: 9091
  serviceMonitor:
    install: true
    interval: 30s
    scrapeTimeout: 10s
    metricRelabelings: []
    relabelings: []
  ingress: {}
  nodeSelector: {}
  resources:
    limits:
      cpu: 200m
      memory: 50Mi
    requests:
      cpu: 100m
      memory: 30Mi
  securityContext:
    runAsUser: 2001
    fsGroup: 2001
  tolerations: []
  annotations: {}
  labels: {}
  priorityClassName: priority-class
```


