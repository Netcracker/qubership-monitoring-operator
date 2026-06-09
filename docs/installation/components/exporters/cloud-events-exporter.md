# Cloud Events Exporter

The Cloud Events Exporter monitors Kubernetes events and exposes them as Prometheus
metrics, enabling monitoring and alerting on cluster events such as pod failures,
scaling events, and other operational activities.

## Overview

Cloud Events Exporter collects Kubernetes events from specified namespaces and
converts them into metrics that can be scraped by Prometheus. This enables
monitoring teams to:

- Track event patterns and trends
- Alert on critical events (warnings, errors)  
- Monitor cluster health through event analysis
- Create dashboards showing event statistics

The exporter supports filtering by event type, kind, namespace, and other
attributes to reduce noise and focus on relevant events.

## Configuration Parameters

<!-- markdownlint-disable line-length -->
| Field                                   | Description                                                                                                                                                                                                            | Scheme                                                                                                                       |
| --------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| install                                 | Allows to enable or disable deployment of cloud-events-exporter.                                                                                                                                                      | bool                                                                                                                         |
| name                                    | A deployment name for cloud-events-exporter. Used in deployment name and labels.                                                                                                                                     | string                                                                                                                       |
| image                                   | Docker image for cloud-events-exporter. If not specified, uses default image.                                                                                                                                        | string                                                                                                                       |
| port                                    | Port to scrape cloud events metrics.                                                                                                                                                                                  | int                                                                                                                          |
| namespaces                              | List of namespaces to monitor for events. Empty list monitors all namespaces.                                                                                                                                         | list[string]                                                                                                                 |
| rbac.createClusterRole                  | Allow creating ClusterRole. If set to `false`, ClusterRole must be created manually.                                                                                                                                 | bool                                                                                                                         |
| rbac.createClusterRoleBinding           | Allow creating ClusterRoleBinding. If set to `false`, ClusterRoleBinding must be created manually.                                                                                                                   | bool                                                                                                                         |
| filtering                               | Event filtering configuration to include/exclude specific events.                                                                                                                                                     | object                                                                                                                       |
| filtering.match                         | List of event patterns to include. Events matching these patterns will be exported as metrics.                                                                                                                        | list[object]                                                                                                                 |
| filtering.match[N].type                 | Event type to match (e.g., "Warning", "Normal").                                                                                                                                                                      | string                                                                                                                       |
| filtering.match[N].kind                 | Resource kind to match (e.g., "Pod", "Deployment").                                                                                                                                                                   | string                                                                                                                       |
| filtering.match[N].reason               | Event reason to match.                                                                                                                                                                                                 | string                                                                                                                       |
| filtering.match[N].message              | Event message pattern to match (supports regex).                                                                                                                                                                      | string                                                                                                                       |
| filtering.match[N].namespace            | Namespace to match events from.                                                                                                                                                                                       | string                                                                                                                       |
| filtering.match[N].reportingController  | Reporting controller to match.                                                                                                                                                                                         | string                                                                                                                       |
| filtering.match[N].reportingInstance    | Reporting instance to match.                                                                                                                                                                                           | string                                                                                                                       |
| filtering.exclude                       | List of event patterns to exclude. Events matching these patterns will NOT be exported as metrics.                                                                                                                    | list[object]                                                                                                                 |
| filtering.exclude[N].type               | Event type to exclude (e.g., "Normal").                                                                                                                                                                               | string                                                                                                                       |
| filtering.exclude[N].reason             | Event reason to exclude (supports regex, e.g., "Completed OR Pulled OR Started").                                                                                                                                          | string                                                                                                                       |
| filtering.exclude[N].message            | Event message pattern to exclude (supports regex).                                                                                                                                                                    | string                                                                                                                       |
| serviceMonitor                          | Service monitor configuration for scraping metrics.                                                                                                                                                                   | object                                                                                                                       |
| serviceMonitor.install                  | Enables ServiceMonitor creation for automatic metrics scraping.                                                                                                                                                       | bool                                                                                                                         |
| serviceMonitor.interval                 | Scraping interval for metrics collection.                                                                                                                                                                             | string                                                                                                                       |
| serviceMonitor.scrapeTimeout            | Timeout for scraping metrics. Must be less than interval.                                                                                                                                                             | string                                                                                                                       |
| serviceMonitor.metricRelabelings        | Metric relabeling rules applied to scraped metrics.                                                                                                                                                                   | list[object]                                                                                                                 |
| serviceMonitor.relabelings              | Target relabeling rules applied before scraping.                                                                                                                                                                      | list[object]                                                                                                                 |
| extraArgs                               | Additional arguments for cloud-events-exporter container.                                                                                                                                                             | list[string]                                                                                                                 |
| extraEnvs                               | Additional environment variables for cloud-events-exporter container.                                                                                                                                                 | object                                                                                                                       |
| resources                               | Resource requests and limits for the container.                                                                                                                                                                       | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcerequirements-v1-core) |
| securityContext                         | SecurityContext holds pod-level security attributes.                                                                                                                                                                 | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#podsecuritycontext-v1-core)    |
| nodeSelector                            | Defines which nodes the pods are scheduled on.                                                                                                                                                                       | map[string]string                                                                                                            |
| tolerations                             | Tolerations allow the pods to schedule onto nodes with matching taints.                                                                                                                                              | []v1.Toleration                                                                                                              |
| affinity                                | Pod scheduling constraints.                                                                                                                                                                                           | *v1.Affinity                                                                                                                 |
| annotations                             | Annotations to be added to pods.                                                                                                                                                                                      | map[string]string                                                                                                            |
| labels                                  | Labels to be added to pods.                                                                                                                                                                                           | map[string]string                                                                                                            |
| priorityClassName                       | PriorityClassName assigned to the pods to prevent them from evicting.                                                                                                                                                | string                                                                                                                       |
| serviceAccount.install                  | Specifies whether a ServiceAccount should be created.                                                                                                                                                                | bool                                                                                                                         |
| serviceAccount.annotations              | Annotations for the ServiceAccount.                                                                                                                                                                                   | map[string]string                                                                                                            |
| serviceAccount.labels                   | Labels for the ServiceAccount.                                                                                                                                                                                        | map[string]string                                                                                                            |
<!-- markdownlint-enable line-length -->

## Metrics

The cloud-events-exporter exposes the following metrics:

- `cloud_events_total{type, kind, reason, namespace, ...}` - Total count of events
  by type, kind, reason, and namespace
- `cloud_events_warning_total{...}` - Total count of warning events
- `cloud_events_normal_total{...}` - Total count of normal events
- `cloud_events_reporting_controller_warning_total{controller, ...}` - Warning
  events by reporting controller
- `cloud_events_reporting_controller_normal_total{controller, ...}` - Normal
  events by reporting controller

## Configuration Examples

### Basic Configuration

Enable cloud-events-exporter with default settings:

```yaml
cloudEventsExporter:
  install: true
  name: cloud-events-exporter
  port: 9999
  serviceMonitor:
    install: true
    interval: 30s
    scrapeTimeout: 10s
```

### Monitor Specific Namespaces

Monitor events only from specified namespaces:

```yaml
cloudEventsExporter:
  install: true
  namespaces:
    - monitoring
    - kube-system
    - default
```

### Event Filtering

Filter events to focus on warnings and exclude noisy normal events:

```yaml
cloudEventsExporter:
  install: true
  filtering:
    match:
      - type: "Warning"
        kind: "Pod|Deployment"
      - type: "Warning"
        namespace: "production"
    exclude:
      - type: "Normal"
        message: ".*image.*"
      - reason: "Completed|Pulled|Started"
```

### Custom Resource Configuration

Configure resource limits and scheduling preferences:

```yaml
cloudEventsExporter:
  install: true
  resources:
    limits:
      cpu: 100m
      memory: 128Mi
    requests:
      cpu: 50m
      memory: 64Mi
  nodeSelector:
    node-role.kubernetes.io/worker: "true"
  tolerations:
    - key: "monitoring"
      operator: "Equal"
      value: "true"
      effect: "NoSchedule"
```

## Grafana Dashboard

The cloud-events-exporter includes a pre-built Grafana dashboard that provides:

- Event trends by type (Warning/Normal)
- Top events by namespace and object
- Event breakdown by reporting controllers
- Cluster vs namespace-scoped events
- Event frequency over time

The dashboard is automatically deployed when `install: true` is set.

## Troubleshooting

### Common Issues

1. **No metrics appearing**
   - Check that `install: true` is set
   - Verify RBAC permissions for events access
   - Check pod logs for connection errors

2. **Too many metrics**
   - Configure `filtering.exclude` to reduce noise
   - Limit `namespaces` to specific ones
   - Use metric relabeling in ServiceMonitor

3. **Missing events**
   - Check `filtering.match` configuration
   - Verify namespace access permissions
   - Review event retention in cluster

### Debug Commands

```bash
# Check cloud-events-exporter pod status
kubectl get pods -l app.kubernetes.io/component=cloud-events-exporter -n monitoring

# View exporter logs
kubectl logs -l app.kubernetes.io/component=cloud-events-exporter -n monitoring

# Check metrics endpoint
kubectl port-forward svc/cloud-events-exporter 9999:9999 -n monitoring
curl http://localhost:9999/metrics
```

## Security Considerations

- The exporter requires cluster-wide `get`, `list`, and `watch` permissions on events
- Consider limiting namespace access by configuring `namespaces` parameter
- Use security contexts to run with minimal privileges
- Filter sensitive events that might contain confidential information
