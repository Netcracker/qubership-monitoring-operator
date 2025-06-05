# Architecture

This document describes the detailed architecture of the Qubership Monitoring Operator, a Kubernetes operator that manages the deployment and configuration of a comprehensive monitoring stack. It covers the core components, their relationships, control flows, and integration points.

## Overview

The Qubership Monitoring Operator serves as a centralized controller for managing multiple monitoring components within a Kubernetes environment. It orchestrates the deployment and configuration of Prometheus, VictoriaMetrics, Grafana, AlertManager, and various exporters to create a complete monitoring solution.

```mermaid
graph TB
    subgraph "Kubernetes Cluster"
        MO[Monitoring Operator]
        PM[PlatformMonitoring CR]
        TSDB[VictoriaMetrics OR Prometheus]
        GRAF[Grafana]
        GO[Grafana Operator]
        AM[AlertManager OR VMAlert]
        KSM[kube-state-metrics]
        NE[node-exporter]
        EXPORTERS[Various Exporters]
        
        subgraph "Application Namespaces"
            APP1[Application 1]
            APP2[Application 2]
            SM[ServiceMonitors]
        end
    end
    
    MO -->|Manages| PM
    PM -->|Configures| TSDB
    PM -->|Configures| GRAF
    PM -->|Configures| AM
    GO -->|Manages| GRAF
    
    TSDB -->|Scrapes| KSM
    TSDB -->|Scrapes| NE
    TSDB -->|Scrapes| EXPORTERS
    TSDB -->|Scrapes| APP1
    TSDB -->|Scrapes| APP2
    
    GRAF -->|Queries| TSDB
    
    SM -->|Configures| TSDB
```

## Operator Architecture

The Monitoring Operator is the central component that manages the entire monitoring stack. It watches for and processes a custom resource called `PlatformMonitoring`, which defines the desired state of the monitoring setup.

### Operator Controller Pattern

The operator follows the Kubernetes operator pattern, using the controller-runtime library to watch for changes to the `PlatformMonitoring` resource and reconcile the current state with the desired state.

```mermaid
graph LR
    subgraph "Operator Controller"
        WATCH[Watch Events]
        RECONCILE[Reconcile Logic]
        APPLY[Apply Changes]
    end
    
    CR[PlatformMonitoring CR] -->|Change Event| WATCH
    WATCH --> RECONCILE
    RECONCILE --> APPLY
    APPLY -->|Creates/Updates| RESOURCES[K8s Resources]
```

## Component Architecture

### Time Series Databases

#### Prometheus Stack

The Prometheus stack includes Prometheus itself, AlertManager, and related components. The operator deploys and configures these components using the Prometheus Operator.

```mermaid
graph TB
    subgraph "Prometheus Stack"
        PO[Prometheus Operator]
        PROM[Prometheus Server]
        AM[AlertManager]
        CR[Config Reloader]
        
        subgraph "Monitoring Configuration"
            SM[ServiceMonitor CRs]
            PM[PodMonitor CRs]
            PR[PrometheusRule CRs]
        end
    end
    
    PO -->|Manages| PROM
    PO -->|Manages| AM
    SM -->|Configures| PROM
    PM -->|Configures| PROM
    PR -->|Configures| PROM
    CR -->|Reloads Config| PROM
```

The Prometheus Operator handles the deployment and configuration of Prometheus and AlertManager instances. It automatically generates scrape configurations based on ServiceMonitor and PodMonitor custom resources.

#### VictoriaMetrics Integration

VictoriaMetrics can be used as an alternative or complement to Prometheus for storing metrics.

```mermaid
graph TB
    subgraph "VictoriaMetrics Stack"
        VMO[VM Operator]
        VMAGENT[VMAgent]
        VMSINGLE[VMSingle]
        VMALERT[VMAlert]
        VMAUTH[VMAuth]
        VMALERTMGR[VMAlertManager]
    end
    
    VMO -->|Manages| VMAGENT
    VMO -->|Manages| VMSINGLE
    VMO -->|Manages| VMALERT
    VMO -->|Manages| VMAUTH
    VMO -->|Manages| VMALERTMGR
    
    VMAGENT -->|Writes| VMSINGLE
    VMALERT -->|Queries| VMSINGLE
```

VictoriaMetrics provides a similar but more resource-efficient alternative to Prometheus, with its own set of custom resources for configuration.

### Grafana Stack

The Grafana stack is responsible for visualization of metrics collected by Prometheus or VictoriaMetrics.

```mermaid
graph TB
    subgraph "Grafana Stack"
        GO[Grafana Operator]
        GRAF[Grafana Instance]
        
        subgraph "Grafana Resources"
            GD[GrafanaDashboard CRs]
            GDS[GrafanaDataSource CRs]
        end
    end
    
    GO -->|Manages| GRAF
    GD -->|Provides| GRAF
    GDS -->|Configures| GRAF
```

The Grafana Operator manages Grafana instances, datasources, and dashboards. It automatically discovers and applies GrafanaDashboard custom resources.

## Custom Resource Architecture

The monitoring system uses various custom resources to configure its components:

### PlatformMonitoring

This is the main custom resource that defines the overall monitoring setup. It's watched by the Monitoring Operator.

```yaml
apiVersion: monitoring.qubership.org/v1alpha1
kind: PlatformMonitoring
metadata:
  name: monitoring-stack
spec:
  prometheus:
    install: true
    retention: "7d"
  grafana:
    install: true
    persistence:
      enabled: true
  victoriametrics:
    vmOperator:
      install: true
```

### ServiceMonitor and PodMonitor

These custom resources define what metrics should be collected by Prometheus or VictoriaMetrics:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: example-app
spec:
  selector:
    matchLabels:
      app: example-app
  endpoints:
  - port: metrics
    interval: 30s
```

### PrometheusRule and AlertmanagerConfig

These resources define alerting and recording rules:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: example-alerts
spec:
  groups:
  - name: example.rules
    rules:
    - alert: HighErrorRate
      expr: rate(http_requests_total{status="500"}[5m]) > 0.1
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "High error rate detected"
```

## Metrics Collection Architecture

The system collects metrics from various sources:

```mermaid
graph TB
    subgraph "Metric Sources"
        KUBELET[Kubelet]
        KSM[kube-state-metrics]
        NE[node-exporter]
        APPS[Applications]
        EXPORTERS[External Exporters]
    end
    
    subgraph "Time Series Database"
        TSDB[VictoriaMetrics OR Prometheus]
    end
    
    KUBELET -->|/metrics| TSDB
    KSM -->|/metrics| TSDB
    NE -->|/metrics| TSDB
    APPS -->|/metrics| TSDB
    EXPORTERS -->|/metrics| TSDB
```

## Cloud Provider Integration

The system integrates with various cloud providers for metrics collection:

```mermaid
graph TB
    subgraph "Cloud Providers"
        AWS[AWS CloudWatch]
        AZURE[Azure Monitor]
        GCP[Google Cloud Operations]
    end
    
    subgraph "Cloud Exporters"
        CWE[CloudWatch Exporter]
        PROMITOR[Promitor Agent]
        SDE[Stackdriver Exporter]
    end
    
    subgraph "Monitoring Stack"
        TSDB[VictoriaMetrics OR Prometheus]
    end
    
    AWS -->|API| CWE
    AZURE -->|API| PROMITOR
    GCP -->|API| SDE
    
    CWE -->|/metrics| TSDB
    PROMITOR -->|/metrics| TSDB
    SDE -->|/metrics| TSDB
```

The monitoring operator can deploy specialized exporters for each cloud platform to collect metrics from cloud services and make them available to Prometheus/VictoriaMetrics.

## Deployment Architecture

The system is deployed using Helm charts with a set of configurable values:

```mermaid
graph TB
    subgraph "Deployment Process"
        HELM[Helm Chart]
        VALUES[values.yaml]
        
        subgraph "Generated Resources"
            PM[PlatformMonitoring CR]
            OPERATORS[Operator Deployments]
            CONFIGS[ConfigMaps/Secrets]
        end
    end
    
    subgraph "Operator Controllers"
        MO[Monitoring Operator]
        TSDB_OP[VM Operator OR Prometheus Operator]
        GO[Grafana Operator]
    end
    
    subgraph "Monitoring Components"
        TSDB[VictoriaMetrics OR Prometheus]
        GRAF[Grafana]
        AM[AlertManager]
    end
    
    HELM -->|Creates| PM
    HELM -->|Creates| OPERATORS
    HELM -->|Creates| CONFIGS
    VALUES -->|Configures| HELM
    
    MO -->|Watches| PM
    MO -->|Manages| TSDB
    MO -->|Manages| GRAF
    MO -->|Manages| AM
    
    TSDB_OP -->|Manages| TSDB
    GO -->|Manages| GRAF
```

The deployment can be customized through various configuration options in the Helm chart's values.yaml file, which controls aspects like storage, authentication, resource limits, and cloud provider integration.

## Extension Architecture

The system can be extended through various custom resources and configuration options:

```mermaid
graph TB
    subgraph "User Extensions"
        SM[ServiceMonitor]
        PM[PodMonitor]
        GD[GrafanaDashboard]
        PR[PrometheusRule]
        AC[AlertmanagerConfig]
    end
    
    subgraph "Admin Extensions"
        PLATFORMMON[PlatformMonitoring]
        HELM[Helm Values]
        CRD[Custom CRDs]
    end
    
    subgraph "Monitoring Stack"
        TSDB[VictoriaMetrics OR Prometheus]
        GRAF[Grafana]
        AM[AlertManager]
    end
    
    SM -->|Configures| TSDB
    PM -->|Configures| TSDB
    GD -->|Provides| GRAF
    PR -->|Configures| TSDB
    AC -->|Configures| AM
    
    PLATFORMMON -->|Controls| TSDB
    PLATFORMMON -->|Controls| GRAF
    PLATFORMMON -->|Controls| AM
    HELM -->|Configures| PLATFORMMON
```

This architecture allows for flexible extensions by both users (who can add monitoring for their applications) and administrators (who can configure the overall monitoring system).

## Security Architecture

The system includes security features such as authentication for the monitoring components:

```mermaid
graph TB
    subgraph "Authentication Layer"
        OAUTH[OAuth2/OIDC]
        BASIC[Basic Auth]
        LDAP[LDAP]
    end
    
    subgraph "Monitoring UIs"
        GRAF[Grafana]
        TSDB[VictoriaMetrics OR Prometheus]
        AM[AlertManager]
    end
    
    subgraph "Proxy Layer"
        OAUTH2PROXY[oauth2-proxy]
    end
    
    OAUTH -->|Native| GRAF
    LDAP -->|Native| GRAF
    BASIC -->|Native| GRAF
    
    OAUTH -->|via Proxy| OAUTH2PROXY
    OAUTH2PROXY -->|Protects| TSDB
    OAUTH2PROXY -->|Protects| AM
```

The system supports various authentication methods, including OAuth/OIDC, basic auth, and token-based authentication, as well as TLS encryption for secure communications.

## Component Relationships

The following table shows the relationships between different components:

| Component | Managed By | Configures | Provides Data To |
|-----------|------------|------------|------------------|
| Prometheus | Prometheus Operator | ServiceMonitor, PodMonitor | Grafana, AlertManager |
| VictoriaMetrics | VM Operator | VMServiceScrape, VMPodScrape | Grafana, VMAlert |
| Grafana | Grafana Operator | GrafanaDashboard, GrafanaDataSource | Users |
| AlertManager | Prometheus Operator | AlertmanagerConfig | Notification channels |
| kube-state-metrics | Monitoring Operator | Built-in config | Prometheus, VictoriaMetrics |
| node-exporter | Monitoring Operator | Built-in config | Prometheus, VictoriaMetrics |

## Benefits of This Architecture

The Qubership Monitoring Operator architecture provides several key benefits:

1. **Simplified Management**: Automates the deployment and configuration of complex monitoring components
2. **Comprehensive Monitoring**: Collects metrics from various sources including Kubernetes, applications, and cloud providers
3. **Scalability**: Supports both Prometheus and VictoriaMetrics for metrics storage, allowing for scalable monitoring solutions
4. **Visualization**: Integrates with Grafana for metrics visualization and dashboarding
5. **Alerting**: Provides alerting capabilities through AlertManager and VMAlertManager
6. **Cloud Integration**: Supports integration with major cloud providers
7. **Extensibility**: Allows users and administrators to extend functionality through custom resources
8. **Security**: Provides multiple authentication and authorization options

This architecture enables organizations to deploy and maintain a production-ready monitoring stack with minimal operational overhead while providing the flexibility to customize and extend the system as needed. 