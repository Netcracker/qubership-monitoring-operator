# JSON Exporter

The JSON Exporter allows you to scrape arbitrary JSON endpoints and convert them
into Prometheus metrics. This is useful for monitoring APIs, services that expose
JSON metrics, or any system that provides structured data in JSON format.

## Overview

JSON Exporter provides a way to:

- Extract metrics from JSON APIs and endpoints
- Convert JSON data structures into Prometheus metrics
- Monitor third-party services that don't natively support Prometheus format
- Create custom metrics from REST APIs and web services
- Track application status and business metrics exposed via JSON

The exporter supports complex JSON path expressions, custom labels, and multiple
target configurations.

## Configuration Parameters

<!-- markdownlint-disable line-length -->
| Field                                   | Description                                                                                                                                                                                                            | Scheme                                                                                                                       |
| --------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| install                                 | Allows to enable or disable deployment of json-exporter.                                                                                                                                                              | bool                                                                                                                         |
| name                                    | A name of the microservice to deploy with. Used in deployment name and labels.                                                                                                                                        | string                                                                                                                       |
| image                                   | Docker image for json-exporter.                                                                                                                                                                                       | string                                                                                                                       |
| imagePullPolicy                         | Image pull policy for json-exporter deployment.                                                                                                                                                                       | string                                                                                                                       |
| imagePullSecrets                        | Reference to secrets for pulling images from private registries.                                                                                                                                                      | list[object]                                                                                                                 |
| containerPort                           | Port number that the json-exporter container listens on.                                                                                                                                                              | int                                                                                                                          |
| replicaCount                            | Number of json-exporter pods to create.                                                                                                                                                                               | int                                                                                                                          |
| service                                 | Service configuration for json-exporter.                                                                                                                                                                              | object                                                                                                                       |
| service.type                            | Type of Kubernetes service (ClusterIP, NodePort, LoadBalancer).                                                                                                                                                       | string                                                                                                                       |
| service.port                            | Port number for the service.                                                                                                                                                                                          | int                                                                                                                          |
| service.name                            | Port name used in the service.                                                                                                                                                                                        | string                                                                                                                       |
| service.labels                          | Additional labels for the service.                                                                                                                                                                                    | map[string]string                                                                                                            |
| serviceMonitor                          | ServiceMonitor configuration for metrics scraping.                                                                                                                                                                    | object                                                                                                                       |
| serviceMonitor.enabled                  | Enable ServiceMonitor creation for automatic metrics scraping.                                                                                                                                                        | bool                                                                                                                         |
| serviceMonitor.scheme                   | HTTP scheme to use for scraping (http/https).                                                                                                                                                                         | string                                                                                                                       |
| serviceMonitor.defaults                 | Default scrape interval/timeout/labels for the self-metrics ServiceMonitor (`/metrics`).                                                                                                                                 | object                                                                                                                       |
| serviceMonitor.defaults.interval        | Default scraping interval.                                                                                                                                                                                            | string                                                                                                                       |
| serviceMonitor.defaults.scrapeTimeout   | Default scrape timeout.                                                                                                                                                                                               | string                                                                                                                       |
| serviceMonitor.defaults.labels          | Default labels for the ServiceMonitor.                                                                                                                                                                                | map[string]string                                                                                                            |
| serviceMonitor.defaults.additionalMetricsRelabels | Default metric relabeling rules.                                                                                                                                                                                  | map[string]string                                                                                                            |
| config                                  | Configuration of json-exporter modules and metrics extraction rules.                                                                                                                                                  | object                                                                                                                       |
| config.modules                          | Map of module configurations for different scraping scenarios.                                                                                                                                                        | map[string]object                                                                                                            |
| config.modules[name].metrics            | List of metrics to extract from JSON responses.                                                                                                                                                                       | list[object]                                                                                                                 |
| config.modules[name].metrics[N].name    | Prometheus metric name.                                                                                                                                                                                               | string                                                                                                                       |
| config.modules[name].metrics[N].path    | JSONPath expression to extract the value.                                                                                                                                                                             | string                                                                                                                       |
| config.modules[name].metrics[N].help    | Metric description.                                                                                                                                                                                                    | string                                                                                                                       |
| config.modules[name].metrics[N].type    | Metric type (gauge, counter, etc.).                                                                                                                                                                                   | string                                                                                                                       |
| config.modules[name].metrics[N].labels  | Static and dynamic labels for the metric.                                                                                                                                                                             | map[string]string                                                                                                            |
| config.modules[name].metrics[N].values  | For object-type metrics, map of value names to JSONPath expressions.                                                                                                                                                  | map[string]string                                                                                                            |
| config.modules[name].headers            | HTTP headers to send with requests.                                                                                                                                                                                   | map[string]string                                                                                                            |
| config.modules[name].body               | HTTP body configuration for POST requests.                                                                                                                                                                            | object                                                                                                                       |
| config.modules[name].http_client_config | HTTP client configuration (TLS, auth, etc.).                                                                                                                                                                          | object                                                                                                                       |
| resources                               | Resource requests and limits for the container.                                                                                                                                                                       | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcerequirements-v1-core) |
| securityContext                         | Pod-level security context.                                                                                                                                                                                           | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#podsecuritycontext-v1-core)    |
| containerSecurityContext               | Container-level security context.                                                                                                                                                                                     | [*v1.SecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#securitycontext-v1-core)         |
| nodeSelector                            | Node selector for pod scheduling.                                                                                                                                                                                     | map[string]string                                                                                                            |
| tolerations                             | Tolerations for pod scheduling.                                                                                                                                                                                       | []v1.Toleration                                                                                                              |
| affinity                                | Affinity rules for pod scheduling.                                                                                                                                                                                    | *v1.Affinity                                                                                                                 |
| annotations                             | Annotations to add to pods.                                                                                                                                                                                           | map[string]string                                                                                                            |
| labels                                  | Labels to add to pods.                                                                                                                                                                                                | map[string]string                                                                                                            |
| priorityClassName                       | Priority class name for pod scheduling.                                                                                                                                                                               | string                                                                                                                       |
| serviceAccount                          | ServiceAccount configuration.                                                                                                                                                                                         | object                                                                                                                       |
| serviceAccount.install                  | Whether to create a ServiceAccount.                                                                                                                                                                                   | bool                                                                                                                         |
| serviceAccount.annotations              | Annotations for the ServiceAccount.                                                                                                                                                                                   | map[string]string                                                                                                            |
| serviceAccount.name                     | Name of the ServiceAccount to use.                                                                                                                                                                                    | string                                                                                                                       |
| additionalVolumes                       | Additional volumes to mount.                                                                                                                                                                                           | list[object]                                                                                                                 |
| additionalVolumeMounts                  | Additional volume mounts for the container.                                                                                                                                                                           | list[object]                                                                                                                 |
<!-- markdownlint-enable line-length -->

## Endpoints

The json-exporter provides two main endpoints:

- `/metrics` - Self-monitoring metrics
- `/probe` - Probe endpoint for scraping JSON targets (configure scraping with Prometheus Operator `Probe` CRs, not Helm `ServiceMonitor` templates in this chart)

## Configuration Examples

### Basic Configuration

Enable json-exporter with basic setup:

```yaml
jsonExporter:
  install: true
  name: json-exporter
  containerPort: 7979
  service:
    port: 7979
    name: http
  serviceMonitor:
    enabled: true
```

### Simple JSON API Monitoring

Monitor a simple JSON API endpoint:

```yaml
jsonExporter:
  install: true
  name: json-exporter
  
  # Configure module for API monitoring
  config:
    modules:
      default:
        metrics:
          - name: api_status
            path: "{ .status }"
            help: API status indicator
            labels:
              service: "my-api"
          - name: api_response_time
            path: "{ .response_time_ms }"
            help: API response time in milliseconds
            type: gauge
        headers:
          User-Agent: "json-exporter/1.0"
  
  serviceMonitor:
    enabled: true
```

### Complex JSON Structure

Extract metrics from complex JSON with arrays and nested objects:

```yaml
jsonExporter:
  install: true
  config:
    modules:
      complex_api:
        metrics:
          # Extract global counter
          - name: total_requests
            path: "{ .stats.total_requests }"
            help: Total number of requests
            type: counter
            
          # Extract metrics from array of services
          - name: service_status
            type: object
            help: Status of individual services
            path: '{.services[*]}'
            labels:
              service_name: '{.name}'
              version: '{.version}'
              environment: "production"  # static label
            values:
              status: '{.status == "healthy" ? 1 : 0}'
              uptime_seconds: '{.uptime}'
              
          # Extract timestamped metric
          - name: last_update
            path: "{ .last_update }"
            epochTimestamp: "{ .timestamp }"
            help: Last update timestamp
  
  serviceMonitor:
    enabled: true
```

### POST Request with Body

Configure json-exporter to send POST requests with custom body:

```yaml
jsonExporter:
  install: true
  config:
    modules:
      post_api:
        metrics:
          - name: query_result_count
            path: "{ .data.count }"
            help: Number of results from API query
        headers:
          Content-Type: "application/json"
          X-API-Key: "secret-key"
        body:
          content: |
            {
              "query": "SELECT COUNT(*) FROM users WHERE active = true",
              "format": "json"
            }
  
  serviceMonitor:
    enabled: true
```

### Authentication and TLS

Configure authentication and TLS settings:

```yaml
jsonExporter:
  install: true
  config:
    modules:
      secure_api:
        metrics:
          - name: secure_metric
            path: "{ .value }"
            help: Metric from secure endpoint
        http_client_config:
          tls_config:
            insecure_skip_verify: false
            ca_file: /etc/ssl/certs/ca.pem
          basic_auth:
            username: monitoring_user
            password_file: /etc/secrets/password
            
  # Mount secrets for authentication
  additionalVolumes:
    - name: auth-secret
      secret:
        secretName: json-exporter-auth
    - name: ca-cert
      configMap:
        name: ca-certificates
        
  additionalVolumeMounts:
    - name: auth-secret
      mountPath: /etc/secrets
      readOnly: true
    - name: ca-cert
      mountPath: /etc/ssl/certs
      readOnly: true
  
  serviceMonitor:
    enabled: true
```

### Template-based Body

Use Go templates for dynamic request bodies:

```yaml
jsonExporter:
  install: true
  config:
    modules:
      templated_api:
        metrics:
          - name: dynamic_query_result
            path: "{ .result }"
            help: Result from templated query
        body:
          content: |
            {
              "timestamp": "{{ now.Unix }}",
              "query_id": "{{ .query_id | default `default_query` }}",
              "duration": "{{ duration `300` }}"
            }
          templatize: true
  
  serviceMonitor:
    enabled: true
```

### Production Configuration

Complete production-ready configuration:

```yaml
jsonExporter:
  install: true
  name: json-exporter
  image: prometheuscommunity/json-exporter:v0.7.0
  imagePullPolicy: IfNotPresent
  replicaCount: 2
  
  # Service configuration
  service:
    type: ClusterIP
    port: 7979
    name: http
  
  # Resource management
  resources:
    limits:
      cpu: 200m
      memory: 256Mi
    requests:
      cpu: 100m
      memory: 128Mi
  
  # Security context
  securityContext:
    runAsNonRoot: true
    runAsUser: 65534
    fsGroup: 65534
  
  containerSecurityContext:
    allowPrivilegeEscalation: false
    readOnlyRootFilesystem: true
    capabilities:
      drop:
        - ALL
  
  # Scheduling
  nodeSelector:
    node-role.kubernetes.io/worker: "true"
    
  tolerations:
    - key: "monitoring"
      operator: "Equal"
      value: "true"
      effect: "NoSchedule"
  
  # ServiceAccount
  serviceAccount:
    install: true
    name: json-exporter
    annotations:
      description: "ServiceAccount for JSON Exporter"
  
  # Modules configuration
  config:
    modules:
      api_monitoring:
        metrics:
          - name: api_health_status
            path: "{ .health }"
            help: API health status (1=healthy, 0=unhealthy)
            labels:
              service: '{.service_name}'
              version: '{.version}'
          - name: api_response_time_seconds
            path: "{ .response_time / 1000 }"
            help: API response time in seconds
            type: gauge
      
      business_metrics:
        metrics:
          - name: active_users_total
            path: "{ .metrics.active_users }"
            help: Total number of active users
            type: gauge
          - name: revenue_total
            path: "{ .metrics.revenue }"
            help: Total revenue in cents
            type: counter
        headers:
          Authorization: "Bearer {{ .token }}"
  
  serviceMonitor:
    enabled: true
    scheme: http
    defaults:
      interval: 30s
      scrapeTimeout: 10s
      labels:
        monitoring: "json-exporter"
  
  # Pod labels and annotations
  labels:
    component: monitoring
    tier: exporters
    
  annotations:
    prometheus.io/scrape: "true"
    prometheus.io/port: "7979"
```

## JSONPath Examples

Common JSONPath expressions for extracting data:

```json
{
  "status": "ok",
  "data": {
    "count": 42,
    "items": [
      {"name": "item1", "value": 10},
      {"name": "item2", "value": 20}
    ]
  },
  "meta": {
    "timestamp": 1640995200
  }
}
```

<!-- markdownlint-disable line-length -->
| JSONPath | Description | Result |
|----------|-------------|---------|
| `{ .status }` | Root level field | "ok" |
| `{ .data.count }` | Nested field | 42 |
| `{ .data.items[0].value }` | Array element field | 10 |
| `{ .data.items[*].value }` | All array elements | [10, 20] |
| `{ .data.items[?(@.value > 15)].name }` | Filtered array | ["item2"] |
| `{ .meta.timestamp }` | Timestamp field | 1640995200 |
<!-- markdownlint-enable line-length -->

## Troubleshooting

### Common Issues

1. **No metrics from `/probe` endpoints**
   - Check target URL accessibility from the pod
   - Verify module configuration syntax
   - Check JSONPath expressions with test data
   - Confirm `Probe` CRs and VM/Prometheus scrape configuration

2. **Authentication failures**
   - Verify credentials are mounted correctly
   - Check HTTP client configuration
   - Review target service authentication requirements

3. **Invalid JSON responses**
   - Check if target returns valid JSON
   - Verify content-type headers
   - Use curl to test endpoints manually

### Debug Commands

```bash
# Check json-exporter pod status
kubectl get pods -l app.kubernetes.io/component=json-exporter -n monitoring

# View exporter logs
kubectl logs -l app.kubernetes.io/component=json-exporter -n monitoring

# Test metrics endpoint
kubectl port-forward svc/json-exporter 7979:7979 -n monitoring
curl http://localhost:7979/metrics

# Test probe endpoint manually
curl "http://localhost:7979/probe?module=default&target=http://example.com/api"

# Check configuration
kubectl get configmap json-exporter -o yaml -n monitoring
```

## Security Considerations

- Store sensitive credentials in Kubernetes secrets
- Use RBAC to limit access to json-exporter configuration
- Implement network policies to restrict outbound connections
- Use TLS for external API communications
- Regularly rotate API keys and passwords
- Avoid exposing sensitive data in metric labels

## Performance Tips

- Use appropriate scraping intervals based on data freshness needs
- Limit the number of metrics extracted from complex JSON structures
- Implement timeouts for slow endpoints
- Monitor json-exporter resource usage
- Use metric relabeling to reduce cardinality
- Cache responses when possible using HTTP client configuration
