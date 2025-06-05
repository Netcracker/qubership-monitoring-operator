# Qubership Monitoring Operator

Production-ready monitoring stack for Kubernetes in a single Custom Resource.

## Why Qubership Monitoring Operator?

* **One-Click Stack** – deploy metrics storage, dashboards and alerting with a single `PlatformMonitoring` object.
* **Backend Choice** – pick **VictoriaMetrics** or **Prometheus** depending on scale & preferences.
* **Open-Source Exporters Included** – 20+ exporters wired out of the box: node, kube-state, blackbox, cert, JSON, network-latency, Pushgateway, and more.
* **Public-Cloud Ready** – native integrations with **AWS CloudWatch**, **Azure Monitor / Promitor**, **Google Cloud Operations (Stackdriver)**.
* **Pre-built Dashboards & Alerts** – Grafana dashboards and alert rules for Kubernetes, infrastructure and common services.
* **HA & Zero-Downtime** – supports rolling upgrades, replication and persistent storage.
* **Operator Pattern** – automatic reconciliation, drift protection, GitOps-friendly.

## Core Modules

| Module | Purpose |
| --- | --- |
| Monitoring Operator | Reconciles the entire stack from the `PlatformMonitoring` CR |
| VictoriaMetrics / Prometheus | Scalable time-series database |
| Grafana & Grafana Operator | Visualization, dashboards, data-sources |
| AlertManager / VMAlert | Alert routing & notifications |
| Exporters | node-exporter, kube-state-metrics, blackbox-exporter, cert-exporter, json-exporter, network-latency-exporter, version-exporter, etc. |
| Cloud Exporters | cloudwatch-exporter, promitor-agent, stackdriver-exporter |
| Prometheus Adapter | Custom & external metrics for HPA/KEDA |
| Optional Tools | promxy, graphite-remote-adapter, Pushgateway |

## Get Started in 3 Steps

```bash
# 1. Install the operator (Helm)
helm repo add qubership-monitoring https://<your-repo>
helm repo update
helm install monitoring qubership-monitoring/monitoring-operator \
  --namespace monitoring --create-namespace

# 2. Deploy a monitoring stack with defaults
kubectl apply -f - <<EOF
apiVersion: monitoring.qubership.org/v1alpha1
kind: PlatformMonitoring
metadata:
  name: monitoring-stack
  namespace: monitoring
spec: {}
EOF

# 3. Access Grafana
kubectl port-forward svc/monitoring-grafana -n monitoring 3000:3000
```

## Documentation

* **Install / Upgrade** – [Installation Guide](installation/README.md)
* **Configure** – [Configuration Reference](configuration.md)
* **Architecture** – [Architecture Overview](architecture.md)
* **Cookbook & Examples** – [Recipes](cookbook/README.md)
* **API** – [PlatformMonitoring CRD](api/platform-monitoring.md)
* **Troubleshooting** – [Common Issues](troubleshooting.md)

## Community & Support

* GitHub Issues – <https://github.com/Netcracker/qubership-monitoring-operator/issues>
* GitHub Discussions – <https://github.com/Netcracker/qubership-monitoring-operator>
* Contribute – [CONTRIBUTING.md](https://github.com/Netcracker/qubership-monitoring-operator/blob/main/CONTRIBUTING.md)

---


Start monitoring your cluster in minutes with **Qubership Monitoring Operator**! 