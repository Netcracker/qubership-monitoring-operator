### blackbox-exporter

Blackbox exporter allow to execute probe to URL (static or for Ingress) as to black box and meter state of answer,
code and response time. Also Blackbox exporter check SSL certificate state (if probe configure for URL
with SSL certificate).

Details, how to configure Probes for use with Blackbox exporter please read in
[Configuration documentation](../../../configuration.md#probe).

<!-- markdownlint-disable line-length -->
| Field                                               | Description                                                                                                                                                                                                                         | Scheme                                                                                                                         |
| --------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------ |
| install                                             | Allows to enable or disable deploy blackbox-exporter.                                                                                                                                                                               | bool                                                                                                                           |
| name                                                | A name of the microservice to deploy with. This name is used as the name of the microservice deployment and in labels.                                                                                                              | string                                                                                                                         |
| image                                               | A Docker image to deploy the blackbox-exporter.                                                                                                                                                                                     | string                                                                                                                         |
| asDaemonSet                                         | Allows deploying blackbox-exporter as DaemonSet instead of Deployment.                                                                                                                                                              | bool                                                                                                                           |
| containerSecurityContext                            | Security Context for a container.                                                                                                                                                                                                   | [*v1.SecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#securitycontext-v1-core) |
| configExistingSecretName                            | If the configuration is managed as secret outside the chart, using SealedSecret for example, provide the name of the secret here. If secretConfig is set to true, configExistingSecretName is ignored in favor of the config value. | string                                                                                                                         |
| secretConfig                                        | Store the configuration as a `Secret` instead of a `ConfigMap`, useful in case it contains sensitive data.                                                                                                                          | boolean                                                                                                                        |
| config                                              | Configuration of blackbox-exporter modules.                                                                                                                                                                                         | [Modules Configuration](https://github.com/prometheus/blackbox_exporter/blob/master/CONFIGURATION.md)                          |
| resources                                           | The resources that describe the compute resource requests and limits for single pods.                                                                                                                                               | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcerequirements-v1-core)   |
| servicePort                                         | Port for blackbox-exporter service.                                                                                                                                                                                                 | integer                                                                                                                        |
| containerPort                                       | Only changes container port. Application port can be changed with extraArgs (--web.listen-address=:9115).                                                                                                                           | integer                                                                                                                        |
| createServiceAccount                                | Specifies whether a ServiceAccount should be created.                                                                                                                                                                               | boolean                                                                                                                        |
| extraArgs                                           | Additional arguments for blackbox-exporter container.                                                                                                                                                                               | list[string]                                                                                                                   |
| serviceMonitor.enabled                              | If true, a ServiceMonitor CRD is created for a prometheus operator.                                                                                                                                                                 | boolean                                                                                                                        |
| serviceMonitor.interval                             | Scrape interval.                                                                                                                                                                                                                    | string                                                                                                                         |
| serviceMonitor.scrapeTimeout                        | Scrape timeout.                                                                                                                                                                                                                     | string                                                                                                                         |
| serviceMonitor.scheme                               | HTTP scheme to use for scraping. Can be used with `tlsConfig` for example if using istio mTLS.                                                                                                                                      | string                                                                                                                         |
| serviceMonitor.defaults                             | DEPRECATED! Please use `Probe` instead. Default values that are used for all ServiceMonitors created by `targets`.                                                                                                                  | object                                                                                                                         |
| serviceMonitor.defaults.additionalMetricsRelabels   | DEPRECATED! Please use `Probe` instead. Default additional metrics relabels.                                                                                                                                                        | object                                                                                                                         |
| serviceMonitor.defaults.interval                    | DEPRECATED! Please use `Probe` instead. Default interval.                                                                                                                                                                           | string                                                                                                                         |
| serviceMonitor.defaults.scrapeTimeout               | DEPRECATED! Please use `Probe` instead. Default scrape timeout.                                                                                                                                                                     | string                                                                                                                         |
| serviceMonitor.defaults.module                      | DEPRECATED! Please use `Probe` instead. Default module name.                                                                                                                                                                        | string                                                                                                                         |
| serviceMonitor.targets                              | DEPRECATED! Please use `Probe` instead. Parameters for each targets that are created.                                                                                                                                               | list[object]                                                                                                                   |
| serviceMonitor.targets[N].name                      | DEPRECATED! Please use `Probe` instead. Human readable URL that appears in Prometheus / AlertManager                                                                                                                                | string                                                                                                                         |
| serviceMonitor.targets[N].url                       | DEPRECATED! Please use `Probe` instead. The URL that blackbox scrapes                                                                                                                                                               | string                                                                                                                         |
| serviceMonitor.targets[N].interval                  | DEPRECATED! Please use `Probe` instead. Scraping interval. Overrides value set in `defaults`                                                                                                                                        | string                                                                                                                         |
| serviceMonitor.targets[N].scrapeTimeout             | DEPRECATED! Please use `Probe` instead. Scrape timeout. Overrides value set in `defaults`                                                                                                                                           | string                                                                                                                         |
| serviceMonitor.targets[N].module                    | DEPRECATED! Please use `Probe` instead. Module used for scraping. Overrides value set in `defaults`                                                                                                                                 | string                                                                                                                         |
| serviceMonitor.targets[N].additionalMetricsRelabels | DEPRECATED! Please use `Probe` instead. Map of metric labels and values to add                                                                                                                                                      | object                                                                                                                         |
| grafanaDashboard                                    | Allows to create Grafana dashboard for blackbox-exporter.                                                                                                                                                                           | boolean                                                                                                                        |
| securityContext                                     | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                                  | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#podsecuritycontext-v1-core)      |
| tolerations                                         | Tolerations allow the pods to schedule onto nodes with matching taints.                                                                                                                                                             | []v1.Toleration                                                                                                                |
| nodeSelector                                        | Defines which nodes the pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\"                                                                                                                  | map[string]string                                                                                                              |
| annotations                                         | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value"              | map[string]string                                                                                                              |
| labels                                              | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                                    | map[string]string                                                                                                              |
| affinity                                            | If specified, the pod's scheduling constraints                                                                                                                                                                                      | *v1.Affinity                                                                                                                   |
| priorityClassName                                   | PriorityClassName assigned to the Pods to prevent them from evicting.                                                                                                                                                               | string                                                                                                                         |
<!-- markdownlint-enable line-length -->

Example:

```yaml
blackboxExporter:
  install: true
  name: blackbox-exporter
  image: prom/blackbox-exporter:v0.19.0
  containerSecurityContext:
    runAsUser: 1000
    runAsNonRoot: true
  configExistingSecretName: ""
  secretConfig: false
  config:
    modules:
      http_2xx:
        prober: http
        timeout: 5s
        tls_config:
          insecure_skip_verify: true
        http:
          valid_http_versions: [ "HTTP/1.1", "HTTP/2.0" ]
          no_follow_redirects: false
          preferred_ip_protocol: "ip4"
  resources:
    limits:
      memory: 300Mi
    requests:
      memory: 50Mi
  servicePort: 9115
  containerPort: 9115
  createServiceAccount: true
  extraArgs:
    - "--web.listen-address=:9115"
    - "--timeout-offset=0.5"
    - "--config.check=false"
    - "--history.limit=100"
    - "--web.external-url=http://example.com"
    - "--web.route-prefix=/example/path"
  serviceMonitor:
    enabled: true
    interval: 30s
    scrapeTimeout: 30s
    scheme: http
  grafanaDashboard: true
  securityContext:
    runAsUser: 2001
    fsGroup: 2001
  tolerations: []
  nodeSelector: {}
  labels:
    label.key: label-value
  annotations:
    annotation.key: annotation-value
  priorityClassName: priority-class
```


