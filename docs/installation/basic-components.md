# Basic Components

This document describes the monitoring components and their deployment scenarios.

## Component Overview

The Qubership Monitoring Operator includes many components that serve various monitoring functions. Not all projects need all functions, so you can specify which components to install.

### Always Installed

* **`monitoring-operator`** - Core operator that manages all other components

### Default Configuration

By default, the following components are installed:

* **`monitoring-operator`** - Core operator
* **`alertmanager`** - Alert management and routing
* **`grafana`** - Visualization and dashboards
* **`grafana-operator`** - Grafana management
* **`kube-state-metrics`** - Kubernetes state metrics collection
* **`node-exporter`** - Node-level system metrics
* **`victoriametrics`** - Time series database and storage
* **`vm-operator`** - VictoriaMetrics management

## Available Components

### Core Monitoring Stack

#### VictoriaMetrics Stack
* **`vmoperator`** - VictoriaMetrics operator
* **`vmagent`** - Metrics collection agent
* **`vmalert`** - Alerting component
* **`vmalertmanager`** - Alert management
* **`vmauth`** - Authentication proxy
* **`vmsingle`** - Single-node VictoriaMetrics
* **`vmuser`** - User management

#### Grafana Stack  
* **`grafana-operator`** - Grafana instance management
* **`grafana`** - Visualization and dashboard platform
* **`grafana-dashboards`** - Pre-built monitoring dashboards

#### AlertManager
* **`alertmanager`** - Alert routing, grouping, and notification management

#### Prometheus Stack
* **`prometheus-operator`** - Prometheus instance management
* **`prometheus`** - Time series database and monitoring engine
* **`prometheus-rules`** - Alerting and recording rules
* **`kubernetes-monitors`** - Kubernetes-specific monitoring targets


### Metrics Exporters

#### Core Exporters
* **`node-exporter`** - Hardware and OS metrics
* **`kube-state-metrics`** - Kubernetes object state metrics

#### Specialized Exporters
* **`blackbox-exporter`** - External endpoint probing
* **`cert-exporter`** - TLS certificate monitoring
* **`version-exporter`** - Version tracking
* **`network-latency-exporter`** - Network latency measurements

#### Cloud Exporters
* **`cloudwatch-exporter`** - AWS CloudWatch metrics
* **`promitor-agent-scraper`** - Azure Monitor metrics
* **`stackdriver-exporter`** - Google Cloud metrics

### Additional Components

* **`pushgateway`** - Push-based metrics collection
* **`promxy`** - Prometheus proxy and aggregator
* **`prometheus-adapter`** - Kubernetes custom metrics API
* **`prometheus-adapter-operator`** - Prometheus adapter management
* **`graphite-remote-adapter`** - Graphite protocol support

## Deployment Scenarios

### Scenario 1: Metrics Collection Only
Minimal setup for basic metrics collection:
```yaml
components:
  - monitoring-operator
  - node-exporter
  - kube-state-metrics
  - victoriametrics
```

### Scenario 2: Collection + Visualization
Add dashboards and visualization:
```yaml
components:
  - monitoring-operator
  - node-exporter
  - kube-state-metrics
  - victoriametrics
  - grafana
  - grafana-operator
```

### Scenario 3: Full Monitoring Stack
Complete monitoring with alerting:
```yaml
components:
  - monitoring-operator
  - alertmanager
  - grafana
  - grafana-operator
  - kube-state-metrics
  - node-exporter
  - victoriametrics
  - vm-operator
```

### Scenario 4: External Storage Integration
Send metrics to external systems:
```yaml
components:
  - monitoring-operator
  - node-exporter
  - kube-state-metrics
  - vmagent  # Collection agent
  - graphite-remote-adapter  # External storage
```

## Default Deployment Settings

### Resource Allocation

Default resource requests and limits:

| Component             | Requests                | Limits                  |
| --------------------- | ----------------------- | ----------------------- |
| `alertmanager`        | CPU: 100m, RAM: 100Mi   | CPU: 200m, RAM: 200Mi   |
| `grafana`             | CPU: 300m, RAM: 400Mi   | CPU: 500m, RAM: 800Mi   |
| `grafana-operator`    | CPU: 50m, RAM: 50Mi     | CPU: 100m, RAM: 100Mi   |
| `kube-state-metrics`  | CPU: 50m, RAM: 50Mi     | CPU: 100m, RAM: 256Mi   |
| `monitoring-operator` | CPU: 50m, RAM: 50Mi     | CPU: 100m, RAM: 150Mi   |
| `node-exporter`       | CPU: 50m, RAM: 50Mi     | CPU: 100m, RAM: 100Mi   |
| `victoriametrics`     | CPU: 1000m, RAM: 3Gi    | CPU: 1500m, RAM: 5Gi    |
| `vm-operator`         | CPU: 100m, RAM: 100Mi   | CPU: 200m, RAM: 200Mi   |

### Default Behavior

* **Storage**: VictoriaMetrics and Grafana use ephemeral storage (emptyDir)
* **Retention**: VictoriaMetrics stores metrics for 24 hours by default
* **Ingress**: No ingress configured by default (Kubernetes)
* **Integrations**: No external integrations enabled
* **Exporters**: Basic exporters only (no blackbox/cert exporters)

### Metrics Collection Sources

By default, metrics are collected from:

* Kubelet and cAdvisor
* Kubernetes API Server
* Nginx Ingress (if available)
* Etcd
* Kube-state-metrics
* Node Exporter
* Self-monitoring: VictoriaMetrics, AlertManager, Grafana

### Pre-built Resources

* **Dashboards**: See [Default Dashboards](../defaults/dashboards/)
* **Alerts**: See [Default Alerts](../defaults/alerts.md)

## Component Configuration

### Enabling/Disabling Components

You can selectively enable or disable:

#### Dashboards
Configure via `grafanaDashboards.list` parameter - specify which dashboards to install.

#### Prometheus Rules  
Configure via `prometheusRules.ruleGroups` parameter - specify alert rule groups to install.

#### Kubernetes Monitors
Configure via `kubernetesMonitors` parameter - specify monitoring targets and override scraping parameters.

### Automatic Rules

Some dashboards and monitors are automatically enabled/disabled based on component availability:

#### Dashboard Rules
* `kubernetes-nodes-resources` - disabled if `node-exporter` not installed
* `home-dashboard` - disabled if `grafana.grafanaHomeDashboard` is false
* `core-dns-dashboard` - disabled if CoreDNS monitor not available or OpenShift ≤ 3.11
* Ingress dashboards - disabled on OpenShift or if nginx ingress monitor unavailable
* VictoriaMetrics dashboards - enabled only with VM components
* Prometheus dashboards - enabled only with Prometheus components

## Example Configurations

### Basic Setup
```yaml
# Minimal installation - use all defaults
```

### Custom Ingress Setup
```yaml
grafana:
  ingress:
    install: true
    host: grafana.example.com

victoriametrics:
  vmsingle:
    install: true
    ingress:
      install: true
      host: victoriametrics.example.com

  vmAlertManager:
    install: true  
    ingress:
      install: true
      host: alertmanager.example.com
```

### Cloud-Specific Setup
```yaml
# AWS with CloudWatch
cloudwatch-exporter:
  install: true

# Azure with Promitor
promitor-agent-scraper:
  install: true

# GCP with Stackdriver
stackdriver-exporter:
  install: true
```

## Next Steps

1. Review [Storage](storage.md) configuration options
2. Plan your [Deployment](deploy.md) strategy
3. Configure individual [Components](components/) as needed 