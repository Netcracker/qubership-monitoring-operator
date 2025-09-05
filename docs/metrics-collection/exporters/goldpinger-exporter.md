# Goldpinger Exporter

This document describes the purpose, use cases, configuration, and metrics for the Goldpinger exporter.

## Overview

Goldpinger is a Kubernetes DaemonSet that continuously tests and visualizes network connectivity between nodes and pods
in your cluster. It provides a web UI and exposes Prometheus metrics for monitoring network health and diagnosing
connectivity issues. More details in the [official documentation](https://github.com/bloomberg/goldpinger/blob/master/README.md)

| Application   | Metrics Port | Metrics Endpoint | Need Exporter? | Auth? | Is Exporter Third Party? |
| ------------- | ------------ | --------------- | -------------- | ----- | ------------------------ |
| Goldpinger    | `8081`       | `/metrics`      | Yes            | No    | Yes                      |

---

## Use Cases

Goldpinger is useful in the following scenarios:

- **Network diagnostics:** Quickly detect and visualize network connectivity issues between Kubernetes nodes and pods.
- **Latency monitoring:** Track network latency and packet loss between nodes.
- **Alerting:** Set up Prometheus alerts for network degradation or node isolation.
- **Post-upgrade validation:** Ensure network connectivity after CNI, firewall, or cluster upgrades.
- **Continuous assurance:** Provide ongoing visibility into the health of the Kubernetes network fabric.

---

## How to Collect

Goldpinger exposes Prometheus metrics on port `8081` by default.

### ServiceMonitor Example

To collect metrics using Prometheus Operator, configure a
[ServiceMonitor](../../../charts/qubership-monitoring-operator/charts/goldpinger-exporter/templates/servicemonitor.yaml)
as example:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: goldpinger-service-monitor
  labels:
    app.kubernetes.io/component: monitoring
spec:
  endpoints:
    - interval: 30s
      scrapeTimeout: 10s
      path: /metrics
      port: http
      scheme: http
  selector:
    matchLabels:
      app.kubernetes.io/name: goldpinger
```

### Manual Metrics Collection

You can manually check metrics with:

```bash
curl -v -k -L http://<goldpinger_pod_ip_or_dns>:8081/metrics
```

---

## Metrics List

Goldpinger exports the following key metrics:

- `goldpinger_cluster_health_total` — 1 if all checks pass, 0 otherwise.
- `goldpinger_errors_total` — Statistics of errors per instance.
- `goldpinger_kube_master_response_time_s` — Histogram of response times from Kubernetes API server.
- `goldpinger_nodes_health_total` — Number of nodes seen as healthy/unhealthy from this instance's POV.
- `goldpinger_peers_response_time_s` — Histogram of response times from other hosts (peer calls).
- `goldpinger_stats_total` — Statistics of calls made in goldpinger instances.

**Go runtime metrics:**
- `go_gc_duration_seconds`, `go_goroutines`, `go_info`, `go_memstats_*`, `go_threads` — standard Go application metrics.

**Process/Prometheus handler:**
- `process_cpu_seconds_total`, `process_max_fds`, `process_open_fds`, `process_resident_memory_bytes`, `process_start_time_seconds`, `process_virtual_memory_bytes`, `process_virtual_memory_max_bytes`, `promhttp_metric_handler_requests_in_flight`, `promhttp_metric_handler_requests_total`

---

## Example PrometheusRule

You can set up [alerts](../../../charts/qubership-monitoring-operator/charts/goldpinger-exporter/templates/prometheusrule.yaml)
for unhealthy nodes, as example:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: goldpinger-alerts
  labels:
    app.kubernetes.io/component: monitoring
spec:
  groups:
    - name: goldpinger.rules
      rules:
        - alert: GoldpingerNodeUnhealthy
          expr: sum(goldpinger_nodes_health_total{status="unhealthy"}) > 0
          for: 5m
          labels:
            severity: warning
          annotations:
            summary: "Some nodes are unhealthy according to Goldpinger"
            description: "Goldpinger has detected one or more unhealthy nodes in the cluster."
```

---

## References

- [Goldpinger GitHub Repository](https://github.com/bloomberg/goldpinger)
- [Goldpinger Metrics Documentation](https://github.com/bloomberg/goldpinger#prometheus)

---

**If you need to enable or configure Goldpinger, refer to the
[values.yaml](../../../charts/qubership-monitoring-operator/charts/goldpinger-exporter/values.yaml) in the Helm chart
for all available options.** 
