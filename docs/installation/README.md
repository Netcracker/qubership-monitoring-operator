# Installation

This page provides comprehensive instructions for installing the Qubership Monitoring Operator into Kubernetes or OpenShift environments. The documentation covers prerequisites, supported versions, installation procedures, and post-installation verification.

## Supported Versions

Kubernetes 1.25 or later is required by the managed Prometheus and
VictoriaMetrics operators. OpenShift 4.12 is the minimum equivalent release.

| Platform   | Minimum version |
| ---------- | --------------- |
| Kubernetes | 1.25            |
| OpenShift  | 4.12            |

**Cloud Platform Support:**

| Cloud Platform                       | Support |
| ------------------------------------ | ------- |
| AWS Elastic Kubernetes Service (EKS) | ✓       |
| Azure Kubernetes Service (AKS)       | ✓       |
| Google Kubernetes Engine (GKE)       | ✓       |
| On-premise Kubernetes >= 1.25        | ✓       |
| On-premise OpenShift >= 4.12         | ✓       |

## Prerequisites

### Before You Begin

Ensure you have the following prerequisites in place before installation:

- Kubernetes 1.25+ or OpenShift 4.12+ cluster
- `kubectl` or `oc` compatible with the target cluster
- Helm 3.0+
- Pre-created namespace for installation
- Appropriate permissions for deployment
- Sufficient hardware resources

**Important considerations:**

- Apply CRDs manually if you disable automated CRD upgrade
- For OpenShift deployments, ensure the namespace has an empty node-selector

For detailed prerequisites, see [Prerequisites](prerequisites.md).

For the security constraints of etcd certificate discovery, see
[Etcd certificate synchronization security constraints](etcd-certs-to-secret-security.md).

For the hardening controls and limitations of optional monitoring workloads, see
[Optional component security hardening](optional-components-security-hardening.md).

## Permissions

The monitoring operator requires cluster-level permissions to create and manage the following components:

- prometheus
- prometheus-operator
- grafana-operator
- kube-state-metrics
- node-exporter
- cert-exporter
- Various exporters and monitors

A ClusterRole must be granted to the monitoring-operator service account with all required permissions. Alternatively, the operator can be deployed with restricted privileges, but additional manual setup is required.

## Hardware Requirements

The monitoring stack requires resources depending on the components installed and the scale of your environment.

**Minimal recommended hardware for an average system:**

| Resource | Requirement |
| -------- | ----------- |
| CPU      | 3000m       |
| RAM      | 3000Mi      |
| Storage  | 10GB        |

For detailed hardware sizing information, refer to the [Prerequisites](prerequisites.md) guide.

## Default Deployment

By default, the following components are installed:

| Component           | Description                      |
| ------------------- | -------------------------------- |
| alertmanager        | Alert management and routing     |
| grafana             | Visualization and dashboards     |
| grafana-operator    | Grafana management               |
| kube-state-metrics  | Kubernetes metrics collection    |
| monitoring-operator | Core operator (always installed) |
| node-exporter       | Node-level metrics               |
| victoriametrics     | Time series database             |
| vm-operator         | VictoriaMetrics management       |

**Default configuration includes:**

- Ephemeral storage for VictoriaMetrics and Grafana
- 24-hour metric retention
- Multiple pre-configured dashboards and alerts
- No ingress configured by default
- No integrations enabled

For detailed information about components, see [Components](components/).

## Installation Steps

### Using Helm

1. **Add the Helm repository:**
   ```bash
   helm repo add qubership-monitoring https://your-repo.com/
   helm repo update
   ```

2. **Create values file (values.yaml)** with your customizations or use the default values.

3. **Install the chart:**
   ```bash
   helm install monitoring qubership-monitoring/monitoring-operator \
     --namespace monitoring \
     --create-namespace \
     --values values.yaml
   ```

4. **Upgrading the chart:**
   ```bash
   helm upgrade monitoring qubership-monitoring/monitoring-operator \
     --namespace monitoring \
     --values values.yaml
   ```

5. **Uninstalling the chart:**
   ```bash
   helm uninstall monitoring --namespace monitoring
   ```

!!! note "CRD Management"
    Helm does not update or remove CRDs. For upgrades involving CRD changes, manual updates are required.

## Post-Installation Verification

After installing the monitoring operator, verify that all components are running correctly:

1. **Check pod status:**
   ```bash
   kubectl get pods -n monitoring
   ```

2. **Verify the operator deployment:**
   ```bash
   kubectl get deployment monitoring-operator -n monitoring
   ```

3. **Check Custom Resources:**
   ```bash
   kubectl get crd | grep monitoring
   ```

4. **Access Grafana dashboard:**

   If ingress is configured, access via URL. Otherwise, use port-forwarding:
   ```bash
   kubectl port-forward svc/grafana 3000:3000 -n monitoring
   ```

   Default login for Grafana is typically admin/admin unless configured otherwise.

For complete verification procedures, see [Post-Deploy Checks](post-deploy-checks.md).

## What's Next

After successful installation, you can:

- Configure component-specific settings, see [Components](components/)
- Set up persistent storage, see [Storage Configuration](storage.md)
- Configure monitoring parameters, see [Configuration](../configuration.md)
- Configure authentication and security, see [Authentication](../monitoring-configuration/authentication.md)
- Set up alerting rules, see [Troubleshooting](../troubleshooting.md)
- Review maintenance procedures, see [Maintenance](../maintenance.md)
