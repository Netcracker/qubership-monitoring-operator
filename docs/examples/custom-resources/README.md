---
icon: lucide/blocks
---

# Custom Resource Examples

Working manifests for the custom resources that the monitoring stack consumes. Each example is self-contained: apply it
in a namespace watched by the operator and check the result in Grafana or in the alert list.

## Metric Collection

* [Service Monitor](service-monitor/README.md) — scrape a service through its endpoints.
* [Pod Monitor](pod-monitor/README.md) — scrape pods directly, without a service.
* [Custom Endpoint](custom-endpoint/README.md) — scrape a target that lives outside the cluster.

## Alerting

* [Prometheus Rule](prometheus-rule/README.md) — alerting and recording rules.
* [AlertManager Config](alertmanagerconfig/README.md) — routes and receivers for a single namespace.

## Visualization

* [Grafana Dashboard](grafana-dashboard/README.md) — ship a dashboard with your application.
* [Grafana DataSource](grafana-datasource/README.md) — register an additional data source.

## Complete Services

* [Service with Alarms](service-with-alarms/README.md) — a service that ships its own alerts.
* [Service with Dashboard](service-with-dashboard/README.md) — a service that ships its own dashboard.
* [Service with TLS](service-with-tls/README.md) — a service scraped over TLS.
