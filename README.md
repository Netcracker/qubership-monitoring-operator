# Qubership Monitoring Operator

[![Build](https://github.com/Netcracker/qubership-monitoring-operator/actions/workflows/build.yaml/badge.svg)](https://github.com/Netcracker/qubership-monitoring-operator/actions/workflows/build.yaml)
[![Check Links](https://github.com/Netcracker/qubership-monitoring-operator/actions/workflows/link-checker.yaml/badge.svg)](https://github.com/Netcracker/qubership-monitoring-operator/actions/workflows/link-checker.yaml)
[![Super-Linter](https://github.com/Netcracker/qubership-monitoring-operator/actions/workflows/super-linter.yaml/badge.svg)](https://github.com/Netcracker/qubership-monitoring-operator/actions/workflows/super-linter.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/Netcracker/qubership-monitoring-operator)](https://goreportcard.com/report/github.com/Netcracker/qubership-monitoring-operator)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

A comprehensive Kubernetes operator that simplifies the deployment and management of production-ready monitoring stacks. Built to handle complex monitoring environments with minimal operational overhead while providing maximum flexibility and scalability.

## What is Qubership Monitoring Operator?

The Qubership Monitoring Operator is a cloud-native solution that automates the deployment and management of complete monitoring infrastructure on Kubernetes. It orchestrates industry-standard monitoring tools and provides a unified interface for comprehensive observability.

### Key Benefits

- **Automated Management**: Deploy and manage complex monitoring stacks with a single custom resource
- **Production Ready**: Battle-tested configurations optimized for enterprise environments
- **Multi-Stack Support**: Choose between VictoriaMetrics or Prometheus based on your needs
- **Resource Efficient**: VictoriaMetrics uses 2-5x less RAM compared to Prometheus
- **Cloud Native**: Seamless integration with AWS, Azure, and Google Cloud platforms
- **Complete Observability**: Metrics collection, visualization, alerting, and autoscaling in one package
- **Zero Downtime**: Rolling updates and high availability configurations out of the box

## What You Get

### Core Components
- **Time Series Database**: VictoriaMetrics or Prometheus for metrics storage
- **Visualization**: Grafana with pre-built dashboards for Kubernetes and applications
- **Alerting**: AlertManager or VMAlert for intelligent alert management
- **Metrics Collection**: Automated discovery and scraping of application metrics
- **Autoscaling**: Horizontal Pod Autoscaler integration with custom metrics

### Included Exporters
- **Infrastructure**: node-exporter, kube-state-metrics for Kubernetes insights
- **Security**: cert-exporter for TLS certificate monitoring
- **Network**: blackbox-exporter for endpoint monitoring and network latency tracking
- **Cloud Platforms**: AWS CloudWatch, Azure Monitor, Google Cloud Operations exporters
- **Custom**: JSON exporter for REST APIs, version exporter for application versioning
- **Events**: cloud-events-exporter for CloudEvents monitoring

### Integrations
- **Graphite**: graphite-remote-adapter for Graphite integration
- **Load Balancing**: promxy for high availability and federation

## Architecture

```mermaid
graph TB
    subgraph "Deployment & Management"
        HELM[Helm Chart]
        MO[Monitoring Operator]
        PM[PlatformMonitoring CR]
    end
    
    subgraph "Core Monitoring Stack"
        VM[VictoriaMetrics OR Prometheus Stack]
        GRAF[Grafana]
        AM[AlertManager]
    end
    
    subgraph "Metrics Sources"
        CLOUDS[Public Clouds<br/>AWS CloudWatch, Azure Monitor, Google Cloud Operations]
        LOCAL[Local Metrics<br/>Kubernetes, Infrastructure, Network, Applications]
    end
    
    subgraph "External Integrations"
        NOTIF[Notifications<br/>Slack, Email, PagerDuty]
    end
    
    %% Deployment flow
    HELM -->|deploys| MO
    MO -->|watches| PM
    PM -->|configures| VM
    PM -->|configures| GRAF
    PM -->|configures| AM
    
    %% Data flow
    CLOUDS -->|metrics| VM
    LOCAL -->|metrics| VM
    
    %% Visualization & Alerting
    GRAF -->|queries| VM
    AM -->|alerts from| VM
    AM -->|sends| NOTIF
    
    %% Styling
    classDef management fill:#e3f2fd,stroke:#1976d2,stroke-width:2px
    classDef core fill:#f3e5f5,stroke:#7b1fa2,stroke-width:2px
    classDef sources fill:#e8f5e8,stroke:#388e3c,stroke-width:2px
    classDef external fill:#fff3e0,stroke:#f57c00,stroke-width:2px
    
    class HELM,MO,PM management
    class VM,GRAF,AM core
    class CLOUDS,LOCAL sources
    class NOTIF external
```

## Quick Start

### Prerequisites
- Kubernetes 1.19+ cluster
- Helm 3.0+
- kubectl configured for your cluster

### 1. Install the Operator

**Install from source:**
```bash
# Clone the repository
git clone https://github.com/Netcracker/qubership-monitoring-operator.git
cd qubership-monitoring-operator

# Install the operator from local charts
helm install monitoring-operator charts/qubership-monitoring-operator \
  --namespace monitoring \
  --create-namespace
```

### 2. Deploy Basic Monitoring Stack

The operator comes with sensible defaults. You can deploy with zero configuration:

```bash
# Create PlatformMonitoring resource with defaults
kubectl apply -f - <<EOF
apiVersion: monitoring.qubership.org/v1alpha1
kind: PlatformMonitoring
metadata:
  name: monitoring-stack
  namespace: monitoring
spec: {}
EOF
```

**What gets installed by default:**
- **VictoriaMetrics Operator** - enabled
- **VictoriaMetrics Single** - time series database with 24h retention
- **Grafana** - visualization with pre-built dashboards
- **Grafana Operator** - manages Grafana instances
- **kube-state-metrics** - Kubernetes metrics collector
- **node-exporter** - infrastructure metrics collector
- **Common Dashboards** - essential monitoring dashboards
- **Prometheus Rules** - basic alerting rules

**What's disabled by default:**
- AlertManager (use VMAlert instead)
- All cloud exporters (AWS, Azure, GCP)
- All optional exporters (blackbox, cert, json, etc.)
- Prometheus Adapter for HPA
- Integrations (Graphite, Promxy)

### 3. Access Your Monitoring

```bash
# Get Grafana admin password
kubectl get secret monitoring-grafana-admin -n monitoring -o jsonpath="{.data.password}" | base64 -d

# Port forward to access Grafana
kubectl port-forward -n monitoring svc/monitoring-grafana 3000:3000

# Open http://localhost:3000 (admin/password from above)
```

## Documentation

### Quick Guides
- **[Installation Guide](docs/installation/README.md)** - Detailed installation instructions
- **[Configuration Guide](docs/monitoring-configuration/)** - Complete configuration options
- **[Troubleshooting](docs/troubleshooting.md)** - Common issues and solutions

### API Reference
- **[PlatformMonitoring](docs/api/platform-monitoring.md)** - Main custom resource reference
- **[PrometheusAdapter](docs/api/prometheus-adapter.md)** - HPA metrics adapter configuration
- **[CustomScaleMetricRule](docs/api/custom-scale-metric-rule.md)** - Custom autoscaling metrics

### Default Monitoring
- **[Metrics](docs/defaults/metrics.md)** - Out-of-the-box metrics collection
- **[Alerts](docs/defaults/alerts.md)** - Pre-configured alerting rules
- **[Dashboards](docs/defaults/dashboards/)** - Built-in Grafana dashboards

### Examples
- **[Service Monitoring](docs/examples/custom-resources/)** - Monitor your applications
- **[Full Service Setup](docs/examples/full-service/)** - Complete monitoring setup examples
- **[Cloud Integration](docs/examples/components/)** - Cloud provider integrations

### Architecture
- **[Architecture Overview](docs/architecture.md)** - Detailed system architecture
- **[Component Guide](docs/cookbook/overview.md)** - Understanding the monitoring stack

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Development Setup

```bash
# Clone the repository
git clone https://github.com/Netcracker/qubership-monitoring-operator.git
cd qubership-monitoring-operator

# Install dependencies
go mod download

# Run tests
make test

# Run locally
make run
```

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Support

- **Documentation**: [Full Documentation](docs/)
- **Issues**: [GitHub Issues](https://github.com/Netcracker/qubership-monitoring-operator/issues)
- **Discussions**: [GitHub Discussions](https://github.com/Netcracker/qubership-monitoring-operator/discussions)

## Star History

If this project helped you, please consider giving it a star!

---

**Ready to get started?** Follow our [Quick Start guide](#quick-start) and have monitoring running in minutes!
