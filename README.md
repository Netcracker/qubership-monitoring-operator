# Qubership Monitoring Operator

[![Build](https://github.com/Netcracker/qubership-monitoring-operator/actions/workflows/build.yaml/badge.svg)](https://github.com/Netcracker/qubership-monitoring-operator/actions/workflows/build.yaml)
[![Check Links](https://github.com/Netcracker/qubership-monitoring-operator/actions/workflows/link-checker.yaml/badge.svg)](https://github.com/Netcracker/qubership-monitoring-operator/actions/workflows/link-checker.yaml)
[![Super-Linter](https://github.com/Netcracker/qubership-monitoring-operator/actions/workflows/super-linter.yaml/badge.svg)](https://github.com/Netcracker/qubership-monitoring-operator/actions/workflows/super-linter.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/Netcracker/qubership-monitoring-operator)](https://goreportcard.com/report/github.com/Netcracker/qubership-monitoring-operator)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

A comprehensive Kubernetes operator that simplifies the deployment and management of production-ready monitoring stacks.
Built to handle complex monitoring environments with minimal operational overhead while providing maximum flexibility and scalability.

## What is Qubership Monitoring Operator?

The Qubership Monitoring Operator is a cloud-native solution that automates the deployment and management of complete
monitoring infrastructure on Kubernetes. It orchestrates industry-standard monitoring tools and provides a unified
interface for comprehensive observability.

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

- Kubernetes 1.25+ or OpenShift 4.12+ cluster
- Helm 3.0+
- kubectl configured for your cluster

### 1. Install or upgrade required CRDs

**Install from source:**

```bash
# Clone the repository
git clone https://github.com/Netcracker/qubership-monitoring-operator.git

cd qubership-monitoring-operator

# Run this before installing or upgrading the operator. Server-side apply creates new CRDs, updates existing CRDs, and
# avoids the last-applied-configuration annotation size limit.
kubectl apply --server-side --force-conflicts -f charts/qubership-monitoring-crds/crds/

```
Ordinary `helm upgrade` does not upgrade CRDs. Alternatively, an Argo CD Application pointed to the CRD Helm chart
can apply the complete CRD set. Before transferring existing Helm-managed CRDs to Argo CD, run the one-time
`kubectl apply --server-side --force-conflicts` command above so Argo CD does not encounter field-ownership conflicts.
Example application's spec:

```yaml
project: default
source:
  repoURL: https://github.com/Netcracker/qubership-monitoring-operator
  path: charts/qubership-monitoring-crds
  targetRevision: main # Use the same release branch or tag as the operator Application
destination:
  server: https://kubernetes.default.svc
  namespace: monitoring
syncPolicy:
  automated:
    prune: true
    selfHeal: true
    enabled: true
  syncOptions:
    - CreateNamespace=true
    - ApplyOutOfSyncOnly=true
    - ServerSideApply=true
ignoreDifferences:
  - group: apiextensions.k8s.io
    kind: CustomResourceDefinition
    jsonPointers:
      - /spec/preserveUnknownFields
```

With `ServerSideApply=true`, Argo CD applies CRDs in the same manner as the `kubectl apply --server-side` command
above. When the operator chart is deployed as a separate Application, set `source.helm.skipCrds: true` so that the CRD
and operator Applications do not compete for CRD ownership.

When upgrading from a release that contains Grafana Operator v4 CRDs, including `v0.88.0`, stage the first v5 upgrade
without pruning:

1. Disable automated sync and pruning for both the CRD and operator Applications. Wait for any in-progress sync to
   finish before starting the upgrade.
2. Update and manually sync the CRD Application without pruning. This installs the v5 CRDs while retaining the v4
   Grafana and GrafanaDataSource CRDs needed as migration inputs.
3. Set `source.helm.skipCrds: true` on the operator Application, update it to the same target revision, and manually sync
   it without pruning.
4. Wait for `PlatformMonitoring` to report `Successful=True`, the v5 Grafana resource to become ready, and every v5
   GrafanaDatasource to report successful synchronization. Verify any datasource UID that must remain stable.
5. Re-enable automated sync and pruning. Sync the operator Application first and the CRD Application afterward. The
   operator Application can remain `OutOfSync` while it retains resources that are pending pruning.

Do not delete the CRD Application or enable pruning for it before the v4-to-v5 migration succeeds. The legacy
GrafanaDashboard CRD remains supported for the converter and is not removed by this procedure.

CRDs in a Helm chart's `crds/` directory have a protected lifecycle. Argo CD updates CRDs that remain in the CRD chart,
but it does not automatically delete CRDs removed from that directory. Treat removal of obsolete CRDs as a separate
cluster-administration decision after confirming that no custom resources still use them.

Argo CD 3.3 or later is required for safe operator Application deletion. Starting with 3.3, Argo CD maps the chart's
Helm `pre-delete` cleanup hooks to the `PreDelete` phase.

Delete the operator Application first and wait until its cleanup hooks and foreground deletion complete. Delete the CRD
Application afterward if it is no longer needed. Deleting the CRD Application does not guarantee that Kubernetes CRDs
are removed. Deleting both Applications concurrently, or deleting the CRD Application first, can remove APIs that the
operator cleanup Jobs still need in deployments that track CRDs as ordinary manifests.

Keep managed custom-resource APIs, operators, and their RBAC available until managed resources finish finalizing.
Cluster RBAC cleanup runs afterward. Do not remove finalizers manually during ordinary cleanup.

### 2. Install the Operator

**Install from source:**

```bash
# Clone the repository
git clone https://github.com/Netcracker/qubership-monitoring-operator.git
cd qubership-monitoring-operator

# Install the operator from local charts
# This will automatically create a PlatformMonitoring resource with default configuration
helm install monitoring-operator charts/qubership-monitoring-operator \
  --skip-crds \
  --namespace monitoring \
  --create-namespace
```

With `global.privilegedRights=false`, the managed Grafana and VictoriaMetrics operators use namespace-scoped watches
and `Role` resources. The chart does not create operator-managed `ClusterRole` or `ClusterRoleBinding` resources for
those operators.

**What gets installed automatically:**

- **Monitoring Operator** - manages monitoring stack lifecycle
- **VictoriaMetrics Operator** - enabled
- **VictoriaMetrics Single** - time series database with 14d retention
- **VictoriaMetrics Agent** - metrics collector
- **VictoriaMetrics Alert** - alerting component
- **VictoriaMetrics AlertManager** - alert manager
- **VictoriaMetrics Auth** - authentication proxy
- **Grafana** - visualization with pre-built dashboards
- **Grafana Operator** - manages Grafana instances
- **kube-state-metrics** - Kubernetes metrics collector
- **node-exporter** - infrastructure metrics collector
- **Common Dashboards** - essential monitoring dashboards
- **Prometheus Rules** - basic alerting rules

**What's disabled by default:**

- All cloud exporters (AWS, Azure, GCP)
- All optional exporters (blackbox, cert, JSON, etc.)
- Prometheus Adapter for HPA
- Integrations (Graphite, Promxy)

### 3. Verify Installation

```bash
# Check that monitoring operator is running
kubectl get pods -n monitoring -l "app.kubernetes.io/part-of=monitoring"

# Check PlatformMonitoring resource (created automatically by Helm)
kubectl get platformmonitoring -n monitoring

# Wait for all components to be ready
kubectl get pods -n monitoring
```

### 4. Access Your Monitoring

```bash
# Get Grafana admin password
kubectl get secret monitoring-grafana-admin -n monitoring -o jsonpath="{.data.password}" | base64 -d

# Port forward to access Grafana
kubectl port-forward -n monitoring svc/monitoring-grafana 3000:3000

# Open http://localhost:3000 (admin/password from above)
```

## Documentation

### Quick Guides

- **[Installation Guide](https://netcracker.github.io/qubership-monitoring-operator/installation/)** - Detailed installation instructions
- **[Configuration Guide](https://netcracker.github.io/qubership-monitoring-operator/configuration/)** - Complete configuration options
- **[Troubleshooting](docs/troubleshooting.md)** - Common issues and solutions

### API Reference

- **[PlatformMonitoring](https://netcracker.github.io/qubership-monitoring-operator/api/platform-monitoring/)** - Main custom resource reference
- **[PrometheusAdapter](https://netcracker.github.io/qubership-monitoring-operator/api/prometheus-adapter/)** - HPA metrics adapter configuration

### Default Monitoring

- **[Metrics](https://netcracker.github.io/qubership-monitoring-operator/defaults/metrics/)** - Out-of-the-box metrics collection
- **[Alerts](https://netcracker.github.io/qubership-monitoring-operator/defaults/alerts/)** - Pre-configured alerting rules
- **[Dashboards](https://netcracker.github.io/qubership-monitoring-operator/defaults/dashboards/overall-platform-health/)** - Built-in Grafana dashboards

### Examples

- **[Service Monitoring](https://netcracker.github.io/qubership-monitoring-operator/examples/)** - Monitor your applications
- **[Cloud Watch Integration](https://netcracker.github.io/qubership-monitoring-operator/examples/components/cloudwatch-exporter-config/)** - Cloud provider integrations

### Architecture Overview

- **[Architecture Overview](https://netcracker.github.io/qubership-monitoring-operator/architecture/)** - Detailed system architecture

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

- **Documentation**: [Full Documentation](https://netcracker.github.io/qubership-monitoring-operator)
- **Issues**: [GitHub Issues](https://github.com/Netcracker/qubership-monitoring-operator/issues)

## Star History

If this project helped you, please consider giving it a star!

---

**Ready to get started?** Follow our [Quick Start guide](#quick-start) and have monitoring running in minutes!
