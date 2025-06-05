# Deployment Guide

This guide covers the deployment process for the Qubership Monitoring Operator using Helm.

!!! warning "Legacy Kubernetes Support"
    If you want to deploy the `monitoring-operator` into **Kubernetes v1.15** or lower or **OpenShift v3.11** or lower, you must work with _v1beta1_ CRDs manually. For more information see [Maintenance Guide: Work with legacy CRDs](../maintenance.md#work-with-legacy-crds).

## Overview

This chart installs Monitoring Operator which can create/configure/manage Prometheus/VictoriaMetrics and related components in Kubernetes/OpenShift.

The default installation includes VictoriaMetrics Operator, AlertManager, Exporters, and configuration for scraping the Kubernetes/OpenShift infrastructure.

## Quick Start

### Basic Installation

To install the chart with the release name `monitoring-operator`:

```bash
helm install monitoring-operator charts/monitoring-operator
```

### Installation with Custom Namespace

```bash
helm install monitoring-operator charts/monitoring-operator \
  --namespace monitoring \
  --create-namespace
```

### Installation with Custom Values

```bash
helm install monitoring-operator charts/monitoring-operator \
  --namespace monitoring \
  --create-namespace \
  --values custom-values.yaml
```

## Ingress Configuration

Ingress is enabled by default. You have several options for configuration:

### Automatic Host Configuration

If you want `host` to be installed automatically, specify these parameters:

```yaml
CLOUD_PUBLIC_URL: <public_url.com>
NAMESPACE: <monitoring>
```

Ingress `host` will be set as `<component>-{{ .Values.NAMESPACE }}.{{ .Values.CLOUD_PUBLIC_URL }}`.

Examples:
- grafana-monitoring.public_url.com
- victoriametrics-monitoring.public_url.com
- alertmanager-monitoring.public_url.com

### Manual Host Configuration

You can specify ingress configuration for each component individually:

```yaml
grafana:
  ingress:
    install: true
    host: grafana.example.com
    annotations:
      kubernetes.io/ingress.class: nginx
      cert-manager.io/cluster-issuer: letsencrypt-prod
    tls:
      - secretName: grafana-tls
        hosts:
          - grafana.example.com

victoriametrics:
  vmsingle:
    ingress:
      install: true
      host: victoriametrics.example.com

alertmanager:
  ingress:
    install: false  # Disable ingress for AlertManager
```

## Deployment Examples

### Production Deployment

```yaml
# production-values.yaml
global:
  privilegedRights: true

victoriametrics:
  vmSingle:
    storage:
      storageClassName: fast-ssd
      accessModes:
        - ReadWriteOnce
      resources:
        requests:
          storage: 100Gi

grafana:
  persistence:
    enabled: true
    storageClassName: fast-ssd
    size: 10Gi
  ingress:
    install: true
    host: grafana.production.com
    annotations:
      kubernetes.io/ingress.class: nginx
      cert-manager.io/cluster-issuer: letsencrypt-prod

alertmanager:
  replicas: 3
  ingress:
    install: true
    host: alertmanager.production.com

# Enable additional exporters
blackboxExporter:
  install: true

certExporter:
  install: true
```

Deploy with:

```bash
helm install monitoring-operator charts/monitoring-operator \
  --namespace monitoring \
  --create-namespace \
  --values production-values.yaml
```

### Development Deployment

```yaml
# development-values.yaml
global:
  privilegedRights: true

# Minimal resource allocation
victoriametrics:
  vmSingle:
    resources:
      requests:
        cpu: 500m
        memory: 1Gi
      limits:
        cpu: 1000m
        memory: 2Gi

grafana:
  ingress:
    install: true
    host: grafana.dev.local

# Disable some components
blackboxExporter:
  install: false

certExporter:
  install: false
```

Deploy with:

```bash
helm install monitoring-operator charts/monitoring-operator \
  --namespace monitoring-dev \
  --create-namespace \
  --values development-values.yaml
```

### Cloud-Specific Deployments

#### AWS EKS

```yaml
# aws-values.yaml
publicCloudName: "aws"

victoriametrics:
  vmSingle:
    storage:
      storageClassName: gp3
      resources:
        requests:
          storage: 50Gi

cloudwatchExporter:
  install: true
  
# Use AWS Load Balancer Controller
grafana:
  ingress:
    install: true
    annotations:
      kubernetes.io/ingress.class: alb
      alb.ingress.kubernetes.io/scheme: internet-facing
      alb.ingress.kubernetes.io/target-type: ip
```

#### Azure AKS

```yaml
# azure-values.yaml
publicCloudName: "azure"

victoriametrics:
  vmSingle:
    storage:
      storageClassName: managed-premium
      resources:
        requests:
          storage: 50Gi

promitorAgentScraper:
  install: true

grafana:
  ingress:
    install: true
    annotations:
      kubernetes.io/ingress.class: azure/application-gateway
```

#### Google GKE

```yaml
# gcp-values.yaml
publicCloudName: "google"

victoriametrics:
  vmSingle:
    storage:
      storageClassName: ssd
      resources:
        requests:
          storage: 50Gi

stackdriverExporter:
  install: true

grafana:
  ingress:
    install: true
    annotations:
      kubernetes.io/ingress.class: gce
      kubernetes.io/ingress.global-static-ip-name: monitoring-ip
```

## Upgrading

To upgrade the chart with the release name `monitoring-operator`:

```bash
helm upgrade monitoring-operator charts/monitoring-operator
```

### Upgrade with New Values

```bash
helm upgrade monitoring-operator charts/monitoring-operator \
  --values updated-values.yaml
```

### Upgrade from Specific Version

```bash
helm upgrade monitoring-operator charts/monitoring-operator \
  --version 1.2.3
```

## Uninstalling

To uninstall the `monitoring-operator` deployment:

```bash
helm uninstall monitoring-operator
```

!!! warning "CRD Cleanup Required"
    This command removes all Kubernetes components associated with the chart but **does not remove CRDs**. Deleting CRDs causes the deletion of all resources of their type, including resources from other applications.

### Manual CRD Cleanup

CRDs created by this chart should be manually cleaned up if needed:

#### Kubernetes

```bash
kubectl delete crd grafanas.integreatly.org
kubectl delete crd grafanadashboards.integreatly.org
kubectl delete crd grafanadatasources.integreatly.org
kubectl delete crd grafananotificationchannels.integreatly.org
kubectl delete crd alertmanagers.monitoring.coreos.com
kubectl delete crd alertmanagerconfigs.monitoring.coreos.com
kubectl delete crd podmonitors.monitoring.coreos.com
kubectl delete crd probes.monitoring.coreos.com
kubectl delete crd servicemonitors.monitoring.coreos.com
kubectl delete crd thanosrulers.monitoring.coreos.com
kubectl delete crd customscalemetricrules.monitoring.qubership.org
kubectl delete crd platformmonitorings.monitoring.qubership.org
kubectl delete crd vmsingles.operator.victoriametrics.com
kubectl delete crd vmagents.operator.victoriametrics.com
kubectl delete crd vmalertmanagers.operator.victoriametrics.com
kubectl delete crd vmalerts.operator.victoriametrics.com
```

#### OpenShift

```bash
oc delete crd grafanas.integreatly.org
oc delete crd grafanadashboards.integreatly.org
oc delete crd grafanadatasources.integreatly.org
oc delete crd grafananotificationchannels.integreatly.org
oc delete crd alertmanagers.monitoring.coreos.com
oc delete crd alertmanagerconfigs.monitoring.coreos.com
oc delete crd podmonitors.monitoring.coreos.com
oc delete crd probes.monitoring.coreos.com
oc delete crd servicemonitors.monitoring.coreos.com
oc delete crd thanosrulers.monitoring.coreos.com
oc delete crd customscalemetricrules.monitoring.qubership.org
oc delete crd platformmonitorings.monitoring.qubership.org
oc delete crd vmsingles.operator.victoriametrics.com
oc delete crd vmagents.operator.victoriametrics.com
oc delete crd vmalertmanagers.operator.victoriametrics.com
oc delete crd vmalerts.operator.victoriametrics.com
```

## Troubleshooting

### Common Issues

1. **CRD Installation Failures**: Ensure you have sufficient permissions to create CRDs
2. **Storage Issues**: Verify StorageClass exists and has sufficient space
3. **Network Policies**: Check that network policies allow required communication
4. **Resource Constraints**: Ensure cluster has sufficient CPU/Memory resources

### Verification Commands

```bash
# Check pod status
kubectl get pods -n monitoring

# Check CRDs
kubectl get crd | grep -E "(monitoring|victoriametrics|grafana)"

# Check services
kubectl get svc -n monitoring

# Check ingress
kubectl get ingress -n monitoring
```

## Next Steps

After successful deployment:

1. **[Post-Deploy Checks](post-deploy-checks.md)** - Verify installation
2. **[Configuration](../configuration.md)** - Customize your setup
3. **[Storage](storage.md)** - Configure persistent storage
4. **[Component Configuration](components/)** - Fine-tune individual components 