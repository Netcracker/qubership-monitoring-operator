---
icon: lucide/radio-tower
---

# Exporters

Exporters expose metrics that Prometheus or VictoriaMetrics scrapes. Each page below lists the Helm parameters for one
exporter, the resources the operator creates for it, and the scrape configuration it ships with.

## Cluster and Workload

* [Kube State Metrics](kube-state-metrics.md) — object state metrics for Kubernetes resources.
* [Node Exporter](node-exporter.md) — hardware and OS metrics from every node.
* [Version Exporter](version-exporter.md) — versions of deployed applications and their images.

## Probing and Connectivity

* [Blackbox Exporter](blackbox-exporter.md) — HTTP, TCP, DNS, and ICMP probes.
* [Network Latency Exporter](network-latency-exporter.md) — node-to-node latency measurements.

## Certificates and Endpoints

* [Cert Exporter](cert-exporter.md) — expiry of certificates stored in secrets and on nodes.
* [SSL Exporter](ssl-exporter.md) — expiry of certificates served by external endpoints.
* [JSON Exporter](json-exporter.md) — metrics parsed from arbitrary JSON endpoints.

## Cloud Providers

* [CloudWatch Exporter](cloudwatch-exporter.md) — Amazon CloudWatch metrics.
* [Promitor Agent Scraper](promitor-agent-scraper.md) — Azure Monitor metrics.
* [Stackdriver Exporter](stackdriver-exporter.md) — Google Cloud Operations metrics.

## Events

* [Cloud Events Exporter](cloud-events-exporter.md) — Kubernetes events as metrics and logs.
