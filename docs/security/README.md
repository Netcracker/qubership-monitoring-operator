---
icon: lucide/shield
---

# Security

Authentication and TLS settings for the monitoring stack. Two boundaries are configured separately: access to the
monitoring services themselves, and the connections those services open to scrape targets.

## Authentication

* [Monitoring Authentication](../monitoring-configuration/authentication.md) — protect Grafana, Prometheus,
  VictoriaMetrics, and AlertManager endpoints.
* [Metrics Collection Authentication](../metrics-collection/authentication.md) — credentials used when scraping
  protected targets.

## TLS

* [Monitoring TLS](../monitoring-configuration/tls.md) — certificates for the monitoring services.
* [Metrics Collection TLS](../metrics-collection/tls.md) — certificates and CA bundles used for scraping.
* [Route and Ingress TLS](../monitoring-configuration/route-ingress-tls.md) — TLS on the external entry points.

For a scrape target that serves TLS, see the
[Service with TLS example](../examples/custom-resources/service-with-tls/README.md).
