---
icon: lucide/flame
---

# Prometheus Stack

Prometheus is the alternative to the VictoriaMetrics backend. The Prometheus operator reconciles the custom resources
described below, and `PlatformMonitoring` configures each of them through its own section.

* [Prometheus](prometheus.md) — scraping, storage, and retention.
* [AlertManager](alertmanager.md) — alert routing and delivery.
* [Prometheus Rules](prometheus-rules.md) — alerting and recording rules shipped with the operator.
* [Prometheus Adapter](prometheus-adapter.md) — custom and external metrics for HPA and KEDA.
