### prometheus-adapter

<!-- markdownlint-disable line-length -->
| Field                                | Description                                                                                                                                                                                                                                                                                                                                                           | Scheme                                                                                                                          |
| ------------------------------------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------- |
| install                              | Allow to enable deploy Prometheus Adapter.                                                                                                                                                                                                                                                                                                                            | bool                                                                                                                            |
| image                                | The image to be used for the `prometheus-operator` deployment. The `prometheus-operator` makes the Prometheus configuration Kubernetes native, and manages and operates Prometheus and Alertmanager clusters. For more information, refer to [https://github.com/prometheus-operator/prometheus-operator](https://github.com/prometheus-operator/prometheus-operator) | string                                                                                                                          |
| resources                            | Resources defines resources requests and limits for single Pods.                                                                                                                                                                                                                                                                                                      | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcerequirements-v1-core)    |
| securityContext                      | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                                                                                                                                                                    | *[*v1.SecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#securitycontext-v1-core) |
| paused                               | Set paused to reconciliation.                                                                                                                                                                                                                                                                                                                                         | bool                                                                                                                            |
| prometheusUrl                        | PrometheusURL used to connect to any tool with Prometheus compatible API. It will eventually contain query parameters to configure the connection                                                                                                                                                                                                                     | string                                                                                                                          |
| metricsRelistInterval                | MetricsRelistInterval is the interval at which to update the cache of available metrics from Prometheus                                                                                                                                                                                                                                                               | string                                                                                                                          |
| tolerations                          | Tolerations allow the pods to schedule onto nodes with matching taints.                                                                                                                                                                                                                                                                                               | []v1.Toleration                                                                                                                 |
| nodeSelector                         | NodeSelector defines which nodes the pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\"                                                                                                                                                                                                                                       | map[string]string                                                                                                               |
| affinity                   | If specified, the pod's scheduling constraints                                                                                                                                                                         | *v1.Affinity                                                                                                                 |
| annotations                          | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value"                                                                                                                                                | map[string]string                                                                                                               |
| labels                               | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                                                                                                                                                                      | map[string]string                                                                                                               |
| enableResourceMetrics                | Enable adapter for `metrics.k8s.io`. By default - `false`                                                                                                                                                                                                                                                                                                             | bool                                                                                                                            |
| enableCustomMetrics                  | Enable adapter for `custom.metrics.k8s.io`. By default - `true`                                                                                                                                                                                                                                                                                                       | bool                                                                                                                            |
| customScaleMetricRulesSelector       | CustomScaleMetricRulesSelector defines label selectors to select CustomScaleMetricRule resources across the cluster                                                                                                                                                                                                                                                   | []*[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#labelselector-v1-meta)           |
| APIService.resourceMetrics           | Enable/disable creating APIServices for `metrics.k8s.io`                                                                                                                                                                                                                                                                                                              | bool                                                                                                                            |
| APIService.customMetrics             | Enable/disable creating APIServices for `custom.metrics.k8s.io`                                                                                                                                                                                                                                                                                                       | bool                                                                                                                            |
| auth                                 | Client credentials to connect to Prometheus or Victoriametrics endpoints. (Only basic authentication is supported)                                                                                                                                                                                                                                                    | map[string]string                                                                                                               |
| auth.basicAuth                       | Allow to specify client auth configuration as secret reference                                                                                                                                                                                                                                                                                                        | *[v1.SecretKeySelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#secretkeyselector-v1-core)         |
| auth.basicAuth.createSecret          | Allow to create secret with auth credentials automatically during deploy                                                                                                                                                                                                                                                                                              | map[string]string                                                                                                               |
| auth.basicAuth.createSecret.username | Allow to specify client username                                                                                                                                                                                                                                                                                                                                      | string                                                                                                                          |
| auth.basicAuth.createSecret.password | Allow to specify client password                                                                                                                                                                                                                                                                                                                                      | string                                                                                                                          |
| priorityClassName                    | PriorityClassName assigned to the Pods to prevent them from evicting.                                                                                                                                                                                                                                                                                                 | string                                                                                                                          |
<!-- markdownlint-enable line-length -->

Example (automatically created secrets for auth):

```yaml
prometheusAdapter:
  install: true
  image: "k8s-prometheus-adapter-amd64:v0.6.0"
  securityContext:
    runAsUser: 2000
    fsGroup: 2000
  resources:
    limits:
      cpu: 200m
      memory: 200Mi
    requests:
      cpu: 100m
      memory: 100Mi
  nodeSelector:
    node-role.kubernetes.io/worker: worker
  labels:
    label.key: label-value
  annotations:
    annotation.key: annotation-value
  priorityClassName: priority-class
  metricsRelistInterval: "1m"
  prometheusUrl: "http://prometheus-operated.monitoring.svc:9090"
  APIService:
    resourceMetrics: true
    customMetrics: true
  enableResourceMetrics: true
  enableCustomMetrics: true
  customScaleMetricRulesSelector:
    - matchExpressions:
        - key: app.kubernetes.io/component
          operator: In
          values: [ "monitoring" ]
  auth:
    createSecret:
      basicAuth:
        username: prometheus
        password: prometheus
  operator:
    ...see example by link...
```

Example (precreated secrets for auth):

```yaml
prometheusAdapter:
  install: true
  image: "k8s-prometheus-adapter-amd64:v0.6.0"
  securityContext:
    runAsUser: 2000
    fsGroup: 2000
  resources:
    limits:
      cpu: 200m
      memory: 200Mi
    requests:
      cpu: 100m
      memory: 100Mi
  nodeSelector:
    node-role.kubernetes.io/worker: worker
  labels:
    label.key: label-value
  annotations:
    annotation.key: annotation-value
  metricsRelistInterval: "1m"
  prometheusUrl: "http://prometheus-operated.monitoring.svc:9090"
  APIService:
    resourceMetrics: true
    customMetrics: true
  enableResourceMetrics: true
  enableCustomMetrics: true
  customScaleMetricRulesSelector:
    - matchExpressions:
        - key: app.kubernetes.io/component
          operator: In
          values: [ "monitoring" ]
  auth:
    basicAuth:
      username:
        name: secret
        key: username
      password:
        name: secret
        key: password
  operator:
    ...see example by link...
```


#### prometheus-adapter-operator

<!-- markdownlint-disable line-length -->
| Field             | Description                                                                                                                                                                                                            | Scheme                                                                                                                       |
| ----------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| image             | A docker image to use for prometheus-adapter-operator deployment                                                                                                                                                       | string                                                                                                                       |
| resources         | Resources defines resources requests and limits for single Pods.                                                                                                                                                       | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcerequirements-v1-core) |
| securityContext   | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                     | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#podsecuritycontext-v1-core)    |
| nodeSelector      | NodeSelector defines which nodes the pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\"                                                                                        | map[string]string                                                                                                            |
| affinity                   | If specified, the pod's scheduling constraints                                                                                                                                                                         | *v1.Affinity                                                                                                                 |
| annotations       | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value" | map[string]string                                                                                                            |
| labels            | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                       | map[string]string                                                                                                            |
| priorityClassName | PriorityClassName assigned to the Pods to prevent them from evicting.                                                                                                                                                  | string                                                                                                                       |
| tlsEnabled        | TLS configuration is enabled/disabled. By default, it is disabled.                                                                                                                                                     | boolean                                                                                                                      |
| tlsConfig         | Allow to specify client TLS configuration.                                                                                                                                                                             | [*TLSConfig](#prometheus-adapter-operator-tls-config)                                                                        |
<!-- markdownlint-enable line-length -->

```yaml
prometheusAdapter:
  operator:
    image: "ghcr.io/netcracker/qubership-prometheus-adapter-operator:latest"
    securityContext:
      runAsUser: 2000
      fsGroup: 2000
    resources:
      limits:
        cpu: 200m
        memory: 200Mi
      requests:
        cpu: 100m
        memory: 100Mi
    nodeSelector:
      node-role.kubernetes.io/worker: worker
    labels:
      label.key: label-value
    annotations:
      annotation.key: annotation-value
    priorityClassName: priority-class
    tlsEnabled: true
    tlsConfig:
      generateCerts:
        enabled: true
        duration: 365
        renewBefore: 15
        clusterIssuerName: "dev-cluster-issuer"
        secretName: "prometheus-adapter-client-tls-secret"
```


#### prometheus-adapter-operator-tls-config

TLSConfig holds SSL/TLS configuration attributes.
The parameters are required if SSL/TLS connection is required between Kubernetes cluster and prometheus-adapter-operator.
This section is applicable only if `tlsEnabled` is set to `true`.

<!-- markdownlint-disable line-length -->
| Parameter                         | Type                  | Mandatory | Default value                          | Description                                                                                                                                                                                                                                                                                                                                                                                                                 |
| --------------------------------- | --------------------- | --------- | -------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `caSecret`                        | *v1.SecretKeySelector | no        | `-`                                    | Secret containing the CA certificate to use for the targets.                                                                                                                                                                                                                                                                                                                                                                |
| `certSecret`                      | *v1.SecretKeySelector | no        | `-`                                    | Secret containing the client certificate file for the targets.                                                                                                                                                                                                                                                                                                                                                              |
| `keySecret`                       | *v1.SecretKeySelector | no        | `-`                                    | Secret containing the client key file for the targets.                                                                                                                                                                                                                                                                                                                                                                      |
| `existingSecret`                  | string                | no        | `-`                                    | Name of the pre-existing secret that contains TLS configuration for prometheus-adapter. If specified, `generateCerts.enabled` must be set to `false`. The `existingSecret` is expected to contain CA certificate, TLS key and TLS certificate in `ca.crt`, `tls.key` and `tls.crt` fields respectively. Use either `existingSecret` or the combination of `caSecret`, `certSecret` and `keySecret`. Do not use it together. |
| `generateCerts.enabled`           | boolean               | no        | `true`                                 | Generation of certificate is enabled by default. If `tlsConfig.existingSecret` or the combination of `tlsConfig.caSecret`, `tlsConfig.certSecret` and `tlsConfig.keySecret` is specified, `tlsConfig.generateCerts` section will be skipped. `cert-manager` will generate certificate with the name configured using `generateCerts.secretName`, if it doesn't exist already.                                               |
| `generateCerts.clusterIssuerName` | string                | no        | `-`                                    | Cluster issuer name for generated certificate. This is a mandatory field if `generateCerts.enabled` is set to `true`.                                                                                                                                                                                                                                                                                                       |
| `generateCerts.duration`          | integer               | no        | `365`                                  | Duration in days, until which issued certificate will be valid.                                                                                                                                                                                                                                                                                                                                                             |
| `generateCerts.renewBefore`       | integer               | no        | `15`                                   | Number of days before which certificate must be renewed.                                                                                                                                                                                                                                                                                                                                                                    |
| `generateCerts.secretName`        | string                | no        | `prometheus-adapter-client-tls-secret` | Name of the new secret that needs to be created for storing TLS configuration of prometheus-adapter.                                                                                                                                                                                                                                                                                                                        |
| `createSecret`                    | object                | no        | `-`                                    | New secret with the name `tlsConfig.createSecret.secretName` will be created using already known certificate content. If `tlsConfig.existingSecret` or the combination of `tlsConfig.caSecret`, `tlsConfig.certSecret` and `tlsConfig.keySecret` is specified, `tlsConfig.createSecret` section will be skipped.                                                                                                            |
| `createSecret.ca`                 | string                | no        | `-`                                    | Already known CA certificate will be added to newly created secret.                                                                                                                                                                                                                                                                                                                                                         |
| `createSecret.key`                | string                | no        | `-`                                    | Already known TLS key will be added to newly created secret.                                                                                                                                                                                                                                                                                                                                                                |
| `createSecret.cert`               | string                | no        | `-`                                    | Already known TLS certificate will be added to newly created secret.                                                                                                                                                                                                                                                                                                                                                        |
| `createSecret.secretName`         | string                | no        | `prometheus-adapter-client-tls-secret` | Already known TLS certificate will be added to newly created secret.                                                                                                                                                                                                                                                                                                                                                        |
<!-- markdownlint-enable line-length -->


