# Component Configuration

This section contains detailed configuration documentation for all monitoring components that can be deployed with the Qubership Monitoring Operator.

## Component Categories

### Core Monitoring Stack

* **[AlertManager](prometheus-stack/alertmanager.md)** - Alert routing, grouping, and notification management
* **[Grafana](grafana-stack/grafana.md)** - Visualization and dashboard platform 
* **[VictoriaMetrics](victoriametrics-stack/victoriametrics.md)** - Time series database and storage

### Prometheus Stack

* **[Prometheus](prometheus-stack/prometheus.md)** - Time series database and monitoring engine
* **[Prometheus Operator](prometheus-stack/prometheus.md#prometheus-operator)** - Prometheus instance management
* **[Prometheus Rules](prometheus-stack/prometheus-rules.md)** - Alerting and recording rules
* **[Prometheus Adapter](prometheus-stack/prometheus-adapter.md)** - Kubernetes custom metrics API

### Exporters

* **[Node Exporter](exporters/node-exporter.md)** - Hardware and OS metrics
* **[Kube State Metrics](exporters/kube-state-metrics.md)** - Kubernetes object state metrics
* **[Blackbox Exporter](exporters/blackbox-exporter.md)** - External endpoint probing
* **[Cert Exporter](exporters/cert-exporter.md)** - TLS certificate monitoring
* **[Version Exporter](exporters/version-exporter.md)** - Version tracking
* **[Network Latency Exporter](exporters/network-latency-exporter.md)** - Network latency measurements

### Cloud Exporters

* **[CloudWatch Exporter](exporters/cloudwatch-exporter.md)** - AWS CloudWatch metrics
* **[Promitor Agent](exporters/promitor-agent-scraper.md)** - Azure Monitor metrics
* **[Stackdriver Exporter](exporters/stackdriver-exporter.md)** - Google Cloud metrics

### VictoriaMetrics Operators

* **[VM Operator](victoriametrics-stack/vm-operator.md)** - VictoriaMetrics operator management
* **[VM Agent](victoriametrics-stack/vmagent.md)** - Metrics collection agent
* **[VM Alert](victoriametrics-stack/vmalert.md)** - Alerting component
* **[VM AlertManager](victoriametrics-stack/vmalertmanager.md)** - Alert management
* **[VM Auth](victoriametrics-stack/vmauth.md)** - Authentication proxy
* **[VM Single](victoriametrics-stack/vmsingle.md)** - Single-node VictoriaMetrics
* **[VM User](victoriametrics-stack/vmuser.md)** - User management

### Additional Components

* **[Pushgateway](pushgateway.md)** - Push-based metrics collection
* **[Promxy](promxy.md)** - Prometheus proxy and aggregator
* **[Graphite Remote Adapter](graphite-remote-adapter.md)** - Graphite protocol support

## Common Configuration Patterns

### Authentication
<!-- Most components support authentication configuration via the global `auth` section. See [Authentication](../monitoring-configuration/authentication.md) for details. -->

### TLS Configuration
<!-- Components can be configured with TLS certificates for secure communication. See [TLS Configuration](../monitoring-configuration/tls.md). -->

### Resource Management
All components support resource requests and limits configuration:

```yaml
resources:
  requests:
    cpu: 100m
    memory: 128Mi
  limits:
    cpu: 200m
    memory: 256Mi
```

### Ingress Configuration
Web-based components (Grafana, AlertManager, VictoriaMetrics) can be exposed via ingress:

```yaml
ingress:
  install: true
  host: component.example.com
  annotations:
    kubernetes.io/ingress.class: nginx
```

### Security Context
Components can run with specific security contexts:

```yaml
securityContext:
  runAsUser: 2000
  fsGroup: 2000
  runAsNonRoot: true
```

## Next Steps

1. Review component-specific configuration files above
2. See [Examples](../../examples/) for complete configuration examples
3. Check [Configuration Best Practices](../../user-guides/best-practices.md) 