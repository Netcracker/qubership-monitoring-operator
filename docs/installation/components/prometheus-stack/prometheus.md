### prometheus

<!-- markdownlint-disable line-length -->
| Field                                         | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                          | Scheme                                                                                                                          |
| --------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------- |
| image                                         | The image to be used for the `prometheus` deployment. The `prometheus` is a systems and service monitoring system. It collects metrics from configured targets at given intervals. For more information, refer to [https://github.com/prometheus/prometheus](https://github.com/prometheus/prometheus)                                                                                                                                                               | string                                                                                                                          |
| operator                                      | The operator parameters.                                                                                                                                                                                                                                                                                                                                                                                                                                             | [PrometheusOperator](#prometheus-operator)                                                                                      |
| install                                       | Install indicates whether Prometheus is to be installed. It can be changed for an already deployed service and the service is removed during the next reconciliation iteration.                                                                                                                                                                                                                                                                                      | *bool                                                                                                                           |
| configReloaderImage                           | The image to be used for `prometheus-config-reloader`. The `prometheus-config-reloaded` is an add-on to prometheus that monitors changes in prometheus.yaml and an HTTP request reloads the prometheus configuration. For more information, refer to [https://github.com/prometheus-operator/prometheus-operator/tree/master/cmd/prometheus-config-reloader](https://github.com/prometheus-operator/prometheus-operator/tree/master/cmd/prometheus-config-reloader)  | string                                                                                                                          |
| remoteWrite                                   | RemoteWriteSpec defines the remote_write configuration for prometheus. The `remote_write` allows to transparently send samples to a long term storage. For more information, refer to [https://prometheus.io/docs/operating/integrations/#remote-endpoints-and-storage](https://prometheus.io/docs/operating/integrations/#remote-endpoints-and-storage)                                                                                                             | []promv1.RemoteWriteSpec                                                                                                        |
| remoteWrite.basicAuth.createSecret            | Allow create secret for basic authentification automatically.                                                                                                                                                                                                                                                                                                                                                                                                        | map[string]string                                                                                                               |
| remoteWrite.basicAuth.createSecret.secretName | Name of the secret with which will be created.                                                                                                                                                                                                                                                                                                                                                                                                                       | string                                                                                                                          |
| remoteWrite.basicAuth.createSecret.username   | Username for basic authentification.                                                                                                                                                                                                                                                                                                                                                                                                                                 | string                                                                                                                          |
| remoteWrite.basicAuth.createSecret.password   | Password for basic authentification.                                                                                                                                                                                                                                                                                                                                                                                                                                 | string                                                                                                                          |
| remoteRead                                    | RemoteReadSpec defines the remote_read configuration for prometheus. The `remote_read` allows to transparently receive samples from a long term storage. For more information, refer to [https://prometheus.io/docs/operating/integrations/#remote-endpoints-and-storage](https://prometheus.io/docs/operating/integrations/#remote-endpoints-and-storage)                                                                                                           | []promv1.RemoteReadSpec                                                                                                         |
| secrets                                       | Secrets is a list of secrets in the same namespace as the Prometheus object, which shall be mounted into the Prometheus pods. The Secrets are mounted into `/etc/prometheus/secrets/<secret-name>`.                                                                                                                                                                                                                                                                  | []string                                                                                                                        |
| alerting                                      | Defines the details regarding alerting.                                                                                                                                                                                                                                                                                                                                                                                                                              | *promv1.AlertingSpec                                                                                                            |
| externalLabels                                | The labels to add to any time series or alerts when communicating with external systems (federation, remote storage, Alertmanager).                                                                                                                                                                                                                                                                                                                                  | map[string]string                                                                                                               |
| securityContext                               | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                                                                                                                                                                                                                                                                   | *[*v1.SecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#securitycontext-v1-core) |
| nodeSelector                                  | Defines which nodes the pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\"                                                                                                                                                                                                                                                                                                                                                   | map[string]string                                                                                                               |
| annotations                                   | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value"                                                                                                                                                                                                                                               | map[string]string                                                                                                               |
| labels                                        | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                                                                                                                                                                                                                                                                     | map[string]string                                                                                                               |
| affinity                                      | It specifies the pod's scheduling constraints. For more information, refer to [https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core)                                                                                                                                                                                                         | *v1.Affinity                                                                                                                    |
| resources                                     | Resources defines resource requests and limits for single pods.                                                                                                                                                                                                                                                                                                                                                                                                      | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcerequirements-v1-core)    |
| storage                                       | Storage specifies how storage shall be used. For more information, refer to [https://github.com/prometheus-operator/prometheus-operator/blob/v0.79.2/Documentation/api.md#storagespec](https://github.com/prometheus-operator/prometheus-operator/blob/v0.79.2/Documentation/api.md#storagespec)                                                                                                                                                                       | *promv1.StorageSpec                                                                                                             |
| volumes                                       | Volumes allows configuration of additional volumes on the output StatefulSet definition. Volumes specified are appended to other volumes that are generated as a result of StorageSpec objects. For more information, refer to [https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#volume-v1-core](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#volume-v1-core)                                                            | []v1.Volume                                                                                                                     |
| volumeMounts                                  | VolumeMounts allows configuration of additional VolumeMounts on the output StatefulSet definition. VolumeMounts specified are appended to other VolumeMounts in the prometheus container, that are generated as a result of StorageSpec objects. For more informaiton, refer to [https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#volumemount-v1-core](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#volumemount-v1-core) | []v1.VolumeMount                                                                                                                |
| tlsConfig                                     | Defines the TLS parameters for prometheus. For more information, refer to [TLS guide](../../../monitoring-configuration/tls.md)                                                                                                                                                                                                                                                                                                                                                            | object                                                                                                                          |
| tlsConfig.webTLSConfig                        | For more information, refer to [https://github.com/prometheus-operator/prometheus-operator/blob/v0.79.2/Documentation/api.md#webtlsconfig](https://github.com/prometheus-operator/prometheus-operator/blob/v0.79.2/Documentation/api.md#webtlsconfig)                                                                                                                                                                                                                      | object                                                                                                                          |
| tlsConfig.createSecret                        | Specifies content for secret and create it.                                                                                                                                                                                                                                                                                                                                                                                                                          | object                                                                                                                          |
| tlsConfig.createSecret.name                   | Name of secret with cert, ca and key.                                                                                                                                                                                                                                                                                                                                                                                                                                | string                                                                                                                          |
| tlsConfig.createSecret.cert                   | TLS certificate for prometheus.                                                                                                                                                                                                                                                                                                                                                                                                                                      | string                                                                                                                          |
| tlsConfig.createSecret.key                    | TLS key for prometheus.                                                                                                                                                                                                                                                                                                                                                                                                                                              | string                                                                                                                          |
| tlsConfig.createSecret.ca                     | TLS CA for prometheus.                                                                                                                                                                                                                                                                                                                                                                                                                                               | string                                                                                                                          |
| tlsConfig.generateCerts                       | Allows to configure generation of TLS certificate for Prometheus by cert-manager.                                                                                                                                                                                                                                                                                                                                                                                    | object                                                                                                                          |
| tlsConfig.generateCerts.enabled               | Allows to enable work with cert-manager.                                                                                                                                                                                                                                                                                                                                                                                                                             | bool                                                                                                                            |
| tlsConfig.generateCerts.secretName            | Name of secret generated by cert-manager.                                                                                                                                                                                                                                                                                                                                                                                                                            | string                                                                                                                          |
| tlsConfig.generateCerts.clusterIssuerName     | Defines name of Cluster Issuer. Otherwise, if this parameter is empty, self-signed non-cluster Issuer will be created and used in the new Certificate resource.                                                                                                                                                                                                                                                                                                      | string                                                                                                                          |
| tlsConfig.generateCerts.duration              | Defines duration of the certificate in days.                                                                                                                                                                                                                                                                                                                                                                                                                         | int                                                                                                                             |
| tlsConfig.generateCerts.renewBefore           | Specifies how long before expiry a certificate should be renewed.                                                                                                                                                                                                                                                                                                                                                                                                    | int                                                                                                                             |
| ingress                                       | Ingress allows to create Ingress for the Prometheus UI.                                                                                                                                                                                                                                                                                                                                                                                                              | *[Ingress](#ingress)                                                                                                            |
| retention                                     | Retention policy by time.                                                                                                                                                                                                                                                                                                                                                                                                                                            | string                                                                                                                          |
| retentionsize                                 | Retention policy by size [EXPERIMENTAL].                                                                                                                                                                                                                                                                                                                                                                                                                             | string                                                                                                                          |
| containers                                    | Containers allows injecting additional containers or modifying operator generated containers. This can be used to allow adding an authentication proxy to a Prometheus pod or to change the behavior of an operator generated container.                                                                                                                                                                                                                             | [[]Kubernetes core/v1.Container](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#container-v1-core)        |
| externalUrl                                   | The external URL the Prometheus instances are available under. This is necessary to generate correct URLs. This is necessary if Prometheus is not served from the root of a DNS name.                                                                                                                                                                                                                                                                                | string                                                                                                                          |
| paused                                        | Set paused to reconciliation.                                                                                                                                                                                                                                                                                                                                                                                                                                        | bool                                                                                                                            |
| replicas                                      | Set replicas.                                                                                                                                                                                                                                                                                                                                                                                                                                                        | *int32                                                                                                                          |
| tolerations                                   | Tolerations allow the pods to schedule on nodes with matching taints.                                                                                                                                                                                                                                                                                                                                                                                                | []v1.Toleration                                                                                                                 |
| podMonitor                                    | Pod monitor for self monitoring.                                                                                                                                                                                                                                                                                                                                                                                                                                     | *[Monitor](#monitor)                                                                                                            |
| enableAdminAPI                                | Enable access to prometheus web admin API.                                                                                                                                                                                                                                                                                                                                                                                                                           | bool                                                                                                                            |
| query                                         | Defines the query command line flags when starting Prometheus. For more intoformation, refer to [https://github.com/prometheus-operator/prometheus-operator/blob/v0.79.2/Documentation/api.md#queryspec](https://github.com/prometheus-operator/prometheus-operator/blob/v0.79.2/Documentation/api.md#queryspec)                                                                                                                                                       | promv1.QuerySpec                                                                                                                |
| additionalScrape                              | The additionalScrape allows you to monitor hosts outside the cloud (balancers, graylog server, jenkins...)                                                                                                                                                                                                                                                                                                                                                           | object, [Promethues Config](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#scrape_config)            |
| additionalAlertManager                        | The additionalAlertManager allows you to use AlertManager outside the cloud                                                                                                                                                                                                                                                                                                                                                                                          | object, [Promethues Config](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#scrape_config)            |
| enableFeatures                                | Enable access to Prometheus disabled features. By default, enabled "auto-gomaxprox" feature.. Enabling disabled features is entirely outside the scope of what the maintainers will support and by doing so, you accept that this behavior may break at any time without notice. For more information see <https://prometheus.io/docs/prometheus/latest/disabled_features/>                                                                                          | []string                                                                                                                        |
| scrapeInterval                                | Interval between consecutive scrapes. Default: `30s`                                                                                                                                                                                                                                                                                                                                                                                                                 | string                                                                                                                          |
| scrapeTimeout                                 | Number of seconds to wait for target to respond before erroring. Default: `10s`                                                                                                                                                                                                                                                                                                                                                                                      | string                                                                                                                          |
| evaluationInterval                            | Interval between consecutive evaluations. Default: `30s`                                                                                                                                                                                                                                                                                                                                                                                                             | string                                                                                                                          |
| priorityClassName                             | PriorityClassName assigned to the Pods to prevent them from evicting.                                                                                                                                                                                                                                                                                                                                                                                                | string                                                                                                                          |
<!-- markdownlint-enable line-length -->

Example:

```yaml
prometheus:
  operator:
    ...
  image: prom/prometheus:v2.26.0
  install: true
  configReloaderImage: prometheus-operator/prometheus-config-reloader:v0.48.1
  enableAdminAPI: true
  query:
    lookbackDelta: 2m
  remoteWrite:
    - url: "http://graphite-remote-adapter:9201/write"
    - url: https://1.2.3.4:8429/api/v1/write
      tlsConfig:
        insecureSkipVerify: true
      basicAuth:
        createSecret:
          secretName: esm-vmagent
          username: prometheus
          password: prometheus
  remoteRead:
    - url: "http://graphite-remote-adapter:9201/read"
  secrets:
    - kube-etcd-client-certs
  alerting:
    alertmanagers:
      - namespace: default
        name: alertmanager-example
        port: web
  externalLabels:
    cluster: example_cloud
  securityContext:
    runAsUser: 2000
    fsGroup: 2000
  nodeSelector:
    node-role.kubernetes.io/worker: worker
  labels:
    label.key: label-value
  annotations:
    annotation.key: annotation-value
  affinity:
  resources:
    requests:
      cpu: 2000m
      memory: 8Gi
    limits:
      cpu: 6000m
      memory: 12Gi
  storage:
    volumeClaimTemplate:
      spec:
        storageClassName: nfs-dynamic-provisioning
        resources:
          requests:
            storage: 10Gi
  volumes:
    - name: additional-volume
      hostPath:
        path: /mnt/data/additional_volume
        type: Directory
  volumeMounts:
    - mountPath: /additional_volume
      name: additional-volume
  ingress:
    ...see example by link...
  retention: 24h
  retentionsize: 5Gi
  containers:
    ...see example by link...
  externalUrl: "prometheus.example.cloud.org"
  paused: false
  replicas: 1
  tolerations:
    - key: "example-key"
      operator: "Exists"
      effect: "NoSchedule"
  priorityClassName: priority-class
  podMonitor:
    ...see example by link...
  additionalAlertManager:
    - tls_config:
        insecure_skip_verify: true
      scheme: https
      static_configs:
      - targets:
        - "alertmaneger.example.outside.cloud.org"
  additionalScrape:
    - job_name: graylog
      honor_timestamps: true
      scrape_interval: 30s
      scrape_timeout: 10s
      metrics_path: /api/plugins/org.graylog.plugins.metrics.prometheus/metrics
      scheme: http
      static_configs:
        - targets:
            - 1.2.3.4
      basic_auth:
        username: admin
        password: admin
      tls_config:
        insecure_skip_verify: true
  enableFeatures:
    - auto-gomaxprocs
```

Example of Prometheus config for enabling [Agent mode](https://prometheus.io/blog/2021/11/16/agent/):

```yaml
prometheus:
  ...
  containers:
    - name: prometheus
      args:
        - '--config.file=/etc/prometheus/config_out/prometheus.env.yaml'
        - '--storage.agent.path=/prometheus'
        - '--enable-feature=agent'
        - '--web.enable-lifecycle'
  ...
```


#### prometheus-operator

<!-- markdownlint-disable line-length -->
| Field             | Description                                                                                                                                                                                                                                                                                                                                                           | Scheme                                                                                                                          |
| ----------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------- |
| image             | The image to be used for the `prometheus-operator` deployment. The `prometheus-operator` makes the Prometheus configuration Kubernetes native, and manages and operates Prometheus and Alertmanager clusters. For more information, refer to [https://github.com/prometheus-operator/prometheus-operator](https://github.com/prometheus-operator/prometheus-operator) | string                                                                                                                          |
| resources         | Resources defines resource requests and limits for single pods.                                                                                                                                                                                                                                                                                                       | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcerequirements-v1-core)    |
| securityContext   | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                                                                                                                                                                    | *[*v1.SecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#securitycontext-v1-core) |
| paused            | Set paused to reconciliation.                                                                                                                                                                                                                                                                                                                                         | bool                                                                                                                            |
| tolerations       | Tolerations allow the pods to schedule on nodes with matching taints.                                                                                                                                                                                                                                                                                                 | []v1.Toleration                                                                                                                 |
| nodeSelector      | NodeSelector defines which nodes the pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\"                                                                                                                                                                                                                                       | map[string]string                                                                                                               |
| affinity                   | If specified, the pod's scheduling constraints                                                                                                                                                                         | *v1.Affinity                                                                                                                 |
| annotations       | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value"                                                                                                                                                | map[string]string                                                                                                               |
| labels            | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                                                                                                                                                                      | map[string]string                                                                                                               |
| podMonitor        | Pod monitor for self monitoring.                                                                                                                                                                                                                                                                                                                                      | *[Monitor](#monitor)                                                                                                            |
| priorityClassName | PriorityClassName assigned to the prometheus-operator Pods to prevent them from evicting                                                                                                                                                                                                                                                                              | string                                                                                                                          |

<!-- markdownlint-enable line-length -->

Example:

```yaml
prometheus:
  operator:
    image: prometheus-operator/prometheus-operator:v0.48.1
    resources:
      requests:
        cpu: 20m
        memory: 20Mi
      limits:
        cpu: 50m
        memory: 50Mi
    securityContext:
      runAsUser: 2000
      fsGroup: 2000
    paused: false
    tolerations:
      - key: "example-key"
        operator: "Exists"
        effect: "NoSchedule"
    nodeSelector:
      node-role.kubernetes.io/worker: worker
    labels:
      label.key: label-value
    annotations:
      annotation.key: annotation-value
    priorityClassName: priority-class
    podMonitor:
      ...see example by link...
```


