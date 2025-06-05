### prometheus-rules

<!-- markdownlint-disable line-length -->
| Field       | Description                                                                                                                                                                                                            | Scheme            |
| ----------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------- |
| install     | Allows to install Prometheus Rules for monitoring-operator                                                                                                                                                             | bool              |
| ruleGroups  | List of groups to be installed                                                                                                                                                                                         | list[string]      |
| override    | Allows overriding of Prometheus Rules for monitoring-operator                                                                                                                                                          | list[object]      |
| annotations | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value" | map[string]string |
| labels      | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                       | map[string]string |
<!-- markdownlint-enable line-length -->

Example:

```yaml
prometheusRules:
  install: true
  ruleGroups:
    - Heartbeat
    - SelfMonitoring
    - AlertManager
    - KubebernetesAlerts
    - NodeExporters
    - DockerContainers
    - HAproxy
    - Etcd
  override:
      - group: SelfMonitoring
        alert: PrometheusNotificationsBacklog
        for: 0m
        expr: min_over_time(prometheus_notifications_queue_length[20m]) > 0
        severity: high
  labels:
    label.key: label-value
  annotations:
    annotation.key: annotation-value
```


