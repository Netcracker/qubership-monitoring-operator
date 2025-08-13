### ssl-exporter

SSL exporter allows probing SSL/TLS certificates for various targets (external/internal HTTPS endpoints, files on the host, Kubernetes secrets, and kubeconfig) and exposes metrics for Prometheus.

<!-- markdownlint-disable line-length -->
| Field                                                | Description                                                                                                                                                                                                                                       | Scheme |
|------------------------------------------------------| ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------ |
| install                                              | Enables or disables deployment of ssl-exporter.                                                                                                                                                                                                     | bool   |
| name                                                 | Microservice name used for object names and labels.                                                                                                                                                                                                  | string |
| setupSecurityContext                                 | Creates PodSecurityPolicy or SecurityContextConstraints (when applicable for the platform).                                                                                                                                                          | bool   |
| setupGrafanaDashboard                                | Creates a Grafana dashboard for ssl-exporter.                                                                                                                                                                                                       | bool   |
| setupAlertingRules                                   | Creates Prometheus alerting rules for ssl-exporter (as part of the operator).                                                                                                                                                                       | bool   |
| additionalHostPathVolumes                            | List of HostPath volumes to mount files/directories from the host into the container (e.g., certificates or kubeconfig).                                                                                                                           | list[object] |
| additionalHostPathVolumes[N].volumeName              | Unique volume name.                                                                                                                                                                                                                                  | string |
| additionalHostPathVolumes[N].volumePath              | Path to the file/directory on the host. The same path is used as a mount point inside the container.                                                                                                                                                | string |
| service.type                                         | Kubernetes Service type.                                                                                                                                                                                                                             | string |
| service.port                                         | Service port.                                                                                                                                                                                                                                        | int    |
| service.targetPort                                   | Container target port.                                                                                                                                                                                                                               | int    |
| service.protocol                                     | Service protocol.                                                                                                                                                                                                                                    | string |
| service.name                                         | Service port name.                                                                                                                                                                                                                                   | string |
| service.labels                                       | Additional labels for the Service.                                                                                                                                                                                                                   | object |
| image.repository                                     | Container image repository. Defaults to `ribbybibby/ssl-exporter`.                                                                                                                                                                                   | string |
| image.tag                                            | Image tag. By default equals to chart `appVersion`.                                                                                                                                                                                                  | string |
| image.pullPolicy                                     | Image pull policy.                                                                                                                                                                                                                                   | string |
| image.pullSecret                                     | Docker registry secret name (if required).                                                                                                                                                                                                           | string |
| preDeleteHook.enabled                                | Delete generated resources on chart uninstall.                                                                                                                                                                                                       | bool   |
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
| serviceMonitor.enabled                               | Create a ServiceMonitor for Prometheus Operator.                                                                                                                                                                                                     | bool   |
| serviceMonitor.scheme                                | Scrape scheme (`http`/`https`).                                                                                                                                                                                                                      | string |
| serviceMonitor.defaults.additionalMetricsRelabels    | Default `metricRelabelings` for all targets.                                                                                                                                                                                                         | object |
| serviceMonitor.defaults.interval                     | Default scrape interval.                                                                                                                                                                                                                             | string |
| serviceMonitor.defaults.labels                       | Additional labels for ServiceMonitor.                                                                                                                                                                                                                | object |
| serviceMonitor.defaults.scrapeTimeout                | Default scrape timeout.                                                                                                                                                                                                                              | string |
| serviceMonitor.targets                               | List of targets to create dedicated ServiceMonitors with `/probe` endpoint.                                                                                                                                                                          | list[object] |
| serviceMonitor.targets[N].name                       | Human-friendly target name (also exposed as `target` label).                                                                                                                                                                                         | string |
| serviceMonitor.targets[N].url                        | Target URL (e.g., `google.com:443`) or a file/secret path depending on selected module.                                                                                                                                                              | string |
| serviceMonitor.targets[N].module                     | Module name (`https-external`, `https-selfsigned`, `https-internal`, `file`, `kubernetes`, `kubeconfig`).                                                                                                                                            | string |
| serviceMonitor.targets[N].interval                   | Scrape interval for this target (overrides default).                                                                                                                                                                                                  | string |
| serviceMonitor.targets[N].scrapeTimeout              | Scrape timeout for this target (overrides default).                                                                                                                                                                                                   | string |
| serviceMonitor.targets[N].additionalMetricsRelabels  | Additional metric relabelings for this target.                                                                                                                                                                                                       | object |
| alerts.enabled                                       | Enable a set of preconfigured alerts (if supported by the operator).                                                                                                                                                                                 | bool   |
| prometheusRule.enabled                               | Create a `PrometheusRule` in the cluster using the provided `rules`.                                                                                                                                                                                 | bool   |
| prometheusRule.namespace                             | Explicit namespace for `PrometheusRule`. Defaults to the release namespace.                                                                                                                                                                          | string |
| prometheusRule.labels                                | Additional labels for `PrometheusRule`.                                                                                                                                                                                                              | object |
| prometheusRule.rules                                 | List of alerting rules (same format as in the PrometheusRule CR).                                                                                                                                                                                    | list[object] |
<!-- markdownlint-enable line-length -->

### Example: basic installation with multiple targets

```yaml
sslExporter:
  install: true
  name: ssl-exporter
  setupSecurityContext: true
  setupGrafanaDashboard: true
  setupAlertingRules: true

  # Optional host mounts
  additionalHostPathVolumes:
    - volumeName: host-ssl-cert
      volumePath: /etc/ssl/cert.pem
    - volumeName: host-ca-certs
      volumePath: /etc/ssl/certs
    - volumeName: host-kubeconfig
      volumePath: /etc/rancher/k3s/k3s.yaml

  image:
    repository: ribbybibby/ssl-exporter
    tag: 2.4.3
    pullPolicy: IfNotPresent

  service:
    type: ClusterIP
    port: 9219
    targetPort: 9219
    protocol: TCP
    name: http

  serviceMonitor:
    enabled: true
    scheme: http
    defaults:
      interval: 30s
      scrapeTimeout: 30s
      additionalMetricsRelabels: {}
    targets:
      - name: https-external-google
        url: google.com:443
        module: https-external
        interval: 60s

      - name: https-self-kubernetes-apiserver
        url: kubernetes.default.svc:443
        module: https-selfsigned
        interval: 60s

      - name: https-internal-kubernetes-apiserver
        url: kubernetes.default.svc:443
        module: https-internal
        interval: 60s

      - name: secret-service-ca
        url: kube-system/serving-ca
        module: kubernetes
        interval: 30s

      - name: file-cert-pem
        url: /etc/ssl/cert.pem
        module: file
        interval: 30s

      - name: kubeconfig-cert
        url: /etc/rancher/k3s/k3s.yaml
        module: kubeconfig
        interval: 300s
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
