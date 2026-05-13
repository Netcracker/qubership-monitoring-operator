### ssl-exporter

SSL exporter allows probing SSL/TLS certificates for various targets (external/internal HTTPS endpoints, files on the host, Kubernetes secrets, and kubeconfig) and exposes metrics for Prometheus.

<!-- markdownlint-disable line-length -->
| Field                                                | Description                                                                                                                                                                                                                                       | Scheme |
|------------------------------------------------------| ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------ |
| install                                              | Enables or disables deployment of ssl-exporter.                                                                                                                                                                                                     | bool   |
| name                                                 | Microservice name used for object names and labels.                                                                                                                                                                                                  | string |
| installGrafanaDashboard                                | Creates a Grafana dashboard for ssl-exporter.                                                                                                                                                                                                       | bool   |
| additionalHostPathVolumes                            | List of HostPath volumes to mount files/directories from the host into the container (e.g., certificates or kubeconfig).                                                                                                                           | list[object] |
| additionalHostPathVolumes[N].volumeName              | Unique volume name.                                                                                                                                                                                                                                  | string |
| additionalHostPathVolumes[N].volumePath              | Path to the file/directory on the host. The same path is used as a mount point inside the container.                                                                                                                                                | string |
| service.type                                         | Kubernetes Service type.                                                                                                                                                                                                                             | string |
| service.port                                         | Service port.                                                                                                                                                                                                                                        | int    |
| service.targetPort                                   | Container target port.                                                                                                                                                                                                                               | int    |
| service.protocol                                     | Service protocol.                                                                                                                                                                                                                                    | string |
| service.name                                         | Service port name.                                                                                                                                                                                                                                   | string |
| service.labels                                       | Additional labels for the Service.                                                                                                                                                                                                                   | object |
| image                                                | Full container image (`repository:tag`). If unset, the chart default is defined in `charts/ssl-exporter/templates/_helpers.tpl` (`ribbybibby/ssl-exporter:2.4.3`).                                                                                     | string |
| imagePullPolicy                                      | Image pull policy.                                                                                                                                                                                                                                   | string |
| imagePullSecrets                                     | Pull secrets for the pod (same as Kubernetes `imagePullSecrets`, e.g. `- name: my-registry-secret`).                                                                                                                                               | list   |
| rbac.create                                          | Create RBAC objects (ClusterRole/Binding). Required when using the `kubernetes` module to read secrets.                                                                                                                                              | bool   |
| serviceAccount.create                                | Create a ServiceAccount.                                                                                                                                                                                                                             | bool   |
| serviceAccount.annotations                           | Annotations for ServiceAccount.                                                                                                                                                                                                                      | object |
| serviceAccount.name                                  | ServiceAccount name (if `create: false`, reference an existing one).                                                                                                                                                                                | string |
| podAnnotations                                       | Pod annotations.                                                                                                                                                                                                                                     | object |
| podSecurityContext                                   | Pod securityContext.                                                                                                                                                                                                                                 | object |
| securityContext                                      | Container securityContext.                                                                                                                                                                                                                           | object |
| resources                                            | Container resources.                                                                                                                                                                                                                                 | object |
| nodeSelector                                         | Node selector.                                                                                                                                                                                                                                       | object |
| tolerations                                          | Tolerations list.                                                                                                                                                                                                                                    | list   |
| affinity                                             | Affinity rules.                                                                                                                                                                                                                                      | object |
| modules                                              | Custom ssl-exporter modules configuration. You can override/disable defaults and set parameters like `timeout`, `tls_config`.                                                                                                                        | object |
| modules.https-selfsigned                             | Module for HTTPS checks with self-signed certificates. Enabled by default, `tls_config.insecure_skip_verify: true`.                                                                                                                                  | object |
| modules.https-external                               | Module for external HTTPS checks with system CA (`/etc/ssl/certs/ca-certificates.crt`). Enabled by default.                                                                                                                                         | object |
| modules.https-internal                               | Module for internal HTTPS checks with CA from serviceaccount (`/var/run/secrets/kubernetes.io/serviceaccount/ca.crt`). Enabled by default.                                                                                                          | object |
| modules.file                                         | Module for reading certificates from files inside the container filesystem. Enabled by default.                                                                                                                                                      | object |
| modules.kubernetes                                   | Module for reading certificates from Kubernetes secrets. Enabled by default. Requires RBAC permissions on `secrets` (get/list/watch).                                                                                                               | object |
| modules.kubeconfig                                   | Module for reading certificates from kubeconfig files. Enabled by default.                                                                                                                                                                           | object |
| serviceMonitor.enabled                               | Create a ServiceMonitor that scrapes the workload Service at `/metrics` (standard Prometheus exposition).                                                                                                                                             | bool   |
| serviceMonitor.scheme                                | Scrape scheme (`http`/`https`).                                                                                                                                                                                                                      | string |
| serviceMonitor.interval                              | Scrape interval for the `/metrics` endpoint.                                                                                                                                                                                                         | string |
| serviceMonitor.labels                                | Additional labels for ServiceMonitor metadata.                                                                                                                                                                                                       | object |
| serviceMonitor.scrapeTimeout                         | Scrape timeout for the `/metrics` endpoint.                                                                                                                                                                                                          | string |
| probes.enabled                                       | Create `Probe` resources for active checks via `/probe`.                                                                                                                                                                                             | bool   |
| probes.scheme                                        | Scheme used by Prometheus Operator when calling the ssl-exporter prober service.                                                                                                                                                                      | string |
| probes.path                                          | Path used by Prometheus Operator when calling the ssl-exporter prober service.                                                                                                                                                                        | string |
| probes.defaults.interval                             | Default scrape interval for generated `Probe` resources.                                                                                                                                                                                             | string |
| probes.defaults.labels                               | Additional metadata labels for generated `Probe` resources.                                                                                                                                                                                          | object |
| probes.defaults.scrapeTimeout                        | Default scrape timeout for generated `Probe` resources.                                                                                                                                                                                              | string |
| probes.defaults.additionalMetricsRelabels            | Extra `metricRelabelings` appended to every generated `Probe`.                                                                                                                                                                                       | list[object] |
| probes.targets                                       | List of active checks rendered as `Probe` resources.                                                                                                                                                                                                 | list[object] |
| probes.targets[N].name                               | Human-friendly probe name (also exposed as `target` label).                                                                                                                                                                                          | string |
| probes.targets[N].url                                | Target URL or path depending on module (`google.com:443`, `*/*`, `/etc/ssl/cert.pem`, etc.).                                                                                                                                                        | string |
| probes.targets[N].module                             | ssl-exporter module name (`https-external`, `https-selfsigned`, `https-internal`, `file`, `kubernetes`, `kubeconfig`).                                                                                                                              | string |
| probes.targets[N].interval                           | Probe interval for this target (overrides default).                                                                                                                                                                                                  | string |
| probes.targets[N].scrapeTimeout                      | Probe timeout for this target (overrides default).                                                                                                                                                                                                   | string |
| probes.targets[N].additionalMetricsRelabels          | Additional `metricRelabelings` for this target.                                                                                                                                                                                                      | list[object] |
| prometheusRule.enabled                               | Create a `PrometheusRule` in the cluster using the provided `rules`.                                                                                                                                                                                 | bool   |
| prometheusRule.namespace                             | Explicit namespace for `PrometheusRule`. Defaults to the release namespace.                                                                                                                                                                          | string |
| prometheusRule.labels                                | Additional labels for `PrometheusRule`.                                                                                                                                                                                                              | object |
| prometheusRule.rules                                 | List of alerting rules (same format as in the PrometheusRule CR).                                                                                                                                                                                    | list[object] |
<!-- markdownlint-enable line-length -->

### Example: basic installation

The chart installs ssl-exporter and, when `serviceMonitor.enabled` is true, a single ServiceMonitor that scrapes **`/metrics`** on the workload Service. For **per-target** active checks via **`/probe`**, the chart renders **`Probe`** resources from `sslExporter.probes.targets`. You can also create your own manual `Probe` resources if you need something custom. See [ssl-exporter metrics](../../../metrics-collection/exporters/ssl-exporter.md).

```yaml
sslExporter:
  install: true
  name: ssl-exporter
  installGrafanaDashboard: true

  # Optional host mounts
  additionalHostPathVolumes:
    - volumeName: host-ssl-cert
      volumePath: /etc/ssl/cert.pem
    - volumeName: host-ca-certs
      volumePath: /etc/ssl/certs
    # - volumeName: host-kubeconfig
    #   volumePath: /etc/rancher/k3s/k3s.yaml

  # image: ribbybibby/ssl-exporter:2.4.3
  imagePullPolicy: IfNotPresent
  imagePullSecrets: []

  service:
    type: ClusterIP
    port: 9219
    targetPort: 9219
    protocol: TCP
    name: http

  serviceMonitor:
    enabled: true
    scheme: http
    interval: 30s
    scrapeTimeout: 30s
    labels: {}

  probes:
    enabled: true
    scheme: http
    path: /probe
    defaults:
      interval: 30s
      scrapeTimeout: 30s
      labels: {}
      additionalMetricsRelabels: []
    targets:
      # Example target for an external HTTPS endpoint.
      # Uncomment and adjust as needed.
      # - name: https-external-google
      #   url: google.com:443
      #   module: https-external
      #   interval: 60s
      - name: https-self-kubernetes-apiserver
        url: kubernetes.default.svc:443
        module: https-selfsigned
        interval: 60s
      - name: secret-tls-all-namespaces
        url: "*/*"
        module: kubernetes
        interval: 30s
```

### Example: custom manual Probe

If you do not want to manage probe targets through chart values, create a manual `Probe` resource and point it at the ssl-exporter Service:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: Probe
metadata:
  name: ssl-exporter-external-https-example
  namespace: monitoring
  labels:
    app.kubernetes.io/component: monitoring
spec:
  jobName: ssl-exporter-probe
  interval: 60s
  module: https-external
  prober:
    url: ssl-exporter.monitoring.svc:9219
    scheme: http
    path: /probe
  targets:
    staticConfig:
      static:
        - google.com:443
```

### Example: overriding modules

```yaml
sslExporter:
  install: true
  modules:
    https-external:
      enabled: true
      timeout: 45s
      tls_config:
        ca_file: /etc/ssl/certs/ca-certificates.crt
    https-selfsigned:
      enabled: true
      timeout: 30s
      tls_config:
        insecure_skip_verify: true
    file:
      enabled: true
      timeout: 30s
```

### Example: custom PrometheusRule alerting rules

```yaml
sslExporter:
  install: true
  prometheusRule:
    enabled: true
    # namespace: monitoring
    # labels: {}
    rules:
      - alert: SSLCertExpiringSoon
        expr: ssl_cert_not_after - time() < 86400 * 7
        for: 10m
        labels:
          severity: warning
        annotations:
          summary: "SSL certificate for {{ $labels.instance }} expires soon"
          description: "The SSL certificate for {{ $labels.instance }} will expire in less than 7 days."

      - alert: SSLCertExpired
        expr: ssl_cert_not_after - time() < 0
        for: 10m
        labels:
          severity: critical
        annotations:
          summary: "SSL certificate for {{ $labels.instance }} has expired"
          description: "The SSL certificate for {{ $labels.instance }} has expired."
```

### Security and access notes
- When using the `kubernetes` module, RBAC permissions to read `secrets` (get/list/watch) are required.
- For reading files/crypto-material from the host, use `additionalHostPathVolumes` and ensure the pod has read-only access to those paths.
- Default modules and CA paths are safe; avoid weakening TLS verification unless absolutely necessary (`insecure_skip_verify: true`).
