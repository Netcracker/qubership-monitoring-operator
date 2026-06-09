# Post-Deploy Checks

There are many ways to check whether Monitoring is working correctly and collecting metrics after deployment. This guide covers both common verification steps and specific validation techniques.

## Pod Status Verification

First, check that all components are running using the following command:

```bash
kubectl get pods -n <monitoring_namespace>
```

### Expected Pod List

For a typical Monitoring deployment, you should see pods similar to:

```bash
$ kubectl get pods -n monitoring
NAME                                          READY   STATUS    RESTARTS   AGE
alertmanager-k8s-0                            2/2     Running   0          9d
grafana-deployment-7d869f989d-24657           1/1     Running   0          9d
grafana-operator-5857f4d47d-f7nv4             1/1     Running   0          9d
kube-state-metrics-6c89f47c94-tbx7l           1/1     Running   0          9d
monitoring-operator-monito-6d6c78dd86-wqfcx   1/1     Running   0          9d
node-exporter-g5hsr                           1/1     Running   0          9d
node-exporter-jq7gq                           1/1     Running   0          9d
node-exporter-zppgl                           1/1     Running   0          9d
victoriametrics-operator-c5649d646-7tnw8      1/1     Running   0          9d
vmagent-k8s-8dd4ffccc-msdmx                   2/2     Running   0          9d
vmalert-k8s-6c4b989889-dndmv                  2/2     Running   0          9d
vmalertmanager-k8s-0                          2/2     Running   0          9d
vmauth-k8s-76cc654596-gn8pl                   1/1     Running   0          9d
vmsingle-k8s-6cf5b84846-c22xs                 1/1     Running   0          9d
```

### Troubleshooting Pod Issues

If pods are not in `Running` status, investigate further:

```bash
# Check specific pod details
kubectl describe pod <pod-name> -n monitoring

# Check pod logs
kubectl logs <pod-name> -n monitoring

# For containers with multiple containers
kubectl logs <pod-name> -c <container-name> -n monitoring
```

## Service Health Verification

### VictoriaMetrics Health Check

Check that all VictoriaMetrics services have `OK` status at the `/health` endpoint:

```bash
# Using curl
curl "http://<victoriametrics-pod-ip-or-service-name>/health"

# Using wget
wget -O - "http://<victoriametrics-pod-ip-or-service-name>/health"

# Port-forward to access from local machine
kubectl port-forward svc/vmsingle-k8s 8428:8428 -n monitoring
curl "http://localhost:8428/health"
```

Expected response:
```
OK
```

### Component-Specific Health Checks

#### Grafana Health Check

```bash
kubectl port-forward svc/grafana-service 3000:3000 -n monitoring
curl "http://localhost:3000/api/health"
```

Expected response:
```json
{
  "commit": "...",
  "database": "ok",
  "version": "..."
}
```

#### AlertManager Health Check

```bash
kubectl port-forward svc/alertmanager-k8s 9093:9093 -n monitoring
curl "http://localhost:9093/-/healthy"
```

Expected response:
```
Alertmanager is Healthy.
```

## Metrics Collection Verification

### List Available Metrics

Get the list of metrics collected by VictoriaMetrics:

```bash
# Port-forward VictoriaMetrics
kubectl port-forward svc/vmsingle-k8s 8428:8428 -n monitoring

# Get all metric labels
curl "http://localhost:8428/api/v1/labels" | jq '.'
```

### Verify Specific Metrics

Check that key metrics are being collected:

```bash
# Check node metrics
curl "http://localhost:8428/api/v1/label/__name__/values" | jq '.data[]' | grep node

# Check Kubernetes metrics
curl "http://localhost:8428/api/v1/label/__name__/values" | jq '.data[]' | grep kube

# Check specific metric values
curl "http://localhost:8428/api/v1/query?query=up" | jq '.'
```

### Validate Data Collection

Ensure metrics have recent timestamps and values:

```bash
# Query a basic metric to verify data collection
curl "http://localhost:8428/api/v1/query?query=up{job=\"node-exporter\"}" | jq '.data.result[] | {metric: .metric, value: .value}'

# Check metric count
curl "http://localhost:8428/api/v1/query?query=count(count by (__name__)({__name__!=\"\"}))" | jq '.data.result[0].value[1]'
```

## Service Discovery Verification

### Check Service Monitors

```bash
# List ServiceMonitors
kubectl get servicemonitors -n monitoring

# Check specific ServiceMonitor
kubectl describe servicemonitor <servicemonitor-name> -n monitoring
```

### Check Prometheus Targets

If using Prometheus instead of VictoriaMetrics:

```bash
kubectl port-forward svc/prometheus-k8s 9090:9090 -n monitoring
```

Navigate to `http://localhost:9090/targets` to verify target discovery.

### Check VictoriaMetrics Targets

```bash
kubectl port-forward svc/vmagent-k8s 8429:8429 -n monitoring
```

Navigate to `http://localhost:8429/targets` to verify target discovery.

## UI Access Verification

### Grafana Dashboard Access

```bash
# Port-forward Grafana
kubectl port-forward svc/grafana-service 3000:3000 -n monitoring
```

1. Open `http://localhost:3000`
2. Login with default credentials (admin/admin) or configured credentials
3. Verify dashboards are loaded and displaying data

### VictoriaMetrics UI Access

```bash
# Port-forward VictoriaMetrics
kubectl port-forward svc/vmsingle-k8s 8428:8428 -n monitoring
```

1. Open `http://localhost:8428/vmui`
2. Try executing queries like `up` or `node_cpu_seconds_total`

### AlertManager UI Access

```bash
# Port-forward AlertManager
kubectl port-forward svc/alertmanager-k8s 9093:9093 -n monitoring
```

1. Open `http://localhost:9093`
2. Verify AlertManager interface loads
3. Check for any active alerts

## Configuration Verification

### Check Custom Resources

```bash
# Check PlatformMonitoring
kubectl get platformmonitoring -n monitoring -o yaml

# Check VictoriaMetrics resources
kubectl get vmsingle -n monitoring
kubectl get vmagent -n monitoring
kubectl get vmalert -n monitoring

# Check Grafana resources
kubectl get grafana -n monitoring
kubectl get grafanadashboard -n monitoring
```

### Verify Ingress Configuration

If ingress is enabled:

```bash
# Check ingress resources
kubectl get ingress -n monitoring

# Check ingress details
kubectl describe ingress <ingress-name> -n monitoring

# Test external access (if DNS is configured)
curl -H "Host: grafana.example.com" http://<ingress-ip>
```

## Storage Verification

### Check Persistent Volumes

```bash
# Check PVCs
kubectl get pvc -n monitoring

# Check PV details
kubectl get pv

# Verify storage usage
kubectl exec -it <victoriametrics-pod> -n monitoring -- df -h /victoria-metrics-data
```

## Alerts Verification

### Check Alert Rules

```bash
# Check PrometheusRules
kubectl get prometheusrules -n monitoring

# Verify specific rule
kubectl get prometheusrule <rule-name> -n monitoring -o yaml
```

### Test Alert Generation

Trigger a test alert to verify the alerting pipeline:

```bash
# Create a test metric that should trigger an alert
kubectl run test-alert --image=busybox --restart=Never -- /bin/sh -c "while true; do echo 'test'; sleep 60; done"

# Check if alert appears in AlertManager
curl "http://localhost:9093/api/v1/alerts" | jq '.data[] | select(.labels.alertname=="TestAlert")'
```

## Performance Verification

### Resource Usage Check

```bash
# Check resource usage of monitoring components
kubectl top pods -n monitoring

# Check node resource availability
kubectl top nodes
```

### Query Performance

```bash
# Test query performance
time curl "http://localhost:8428/api/v1/query?query=up"

# Check heavy queries
curl "http://localhost:8428/api/v1/query?query=node_memory_MemTotal_bytes" | jq '.data.result | length'
```

## Common Issues and Solutions

### Issue: Pods in CrashLoopBackOff

**Solutions:**
1. Check resource limits and requests
2. Verify persistent volume permissions
3. Check configuration syntax

```bash
kubectl describe pod <failing-pod> -n monitoring
kubectl logs <failing-pod> -n monitoring --previous
```

### Issue: No Metrics Appearing

**Solutions:**
1. Verify ServiceMonitors are correctly configured
2. Check network policies
3. Verify RBAC permissions

```bash
kubectl get servicemonitor -n monitoring -o yaml
kubectl get networkpolicy -n monitoring
```

### Issue: Grafana Dashboards Empty

**Solutions:**
1. Verify data source configuration
2. Check VictoriaMetrics connectivity
3. Verify time range settings

```bash
kubectl exec -it <grafana-pod> -n monitoring -- grafana-cli admin data-source list
```

## Automated Health Check Script

Create a comprehensive health check script:

```bash
#!/bin/bash
# monitoring-health-check.sh

NAMESPACE="monitoring"

echo "=== Monitoring Health Check ==="

echo "1. Checking pod status..."
kubectl get pods -n $NAMESPACE

echo "2. Checking services..."
kubectl get svc -n $NAMESPACE

echo "3. Checking VictoriaMetrics health..."
VM_POD=$(kubectl get pods -n $NAMESPACE -l app.kubernetes.io/name=vmsingle -o jsonpath='{.items[0].metadata.name}')
kubectl port-forward pod/$VM_POD 8428:8428 -n $NAMESPACE &
sleep 5
curl -s http://localhost:8428/health || echo "VictoriaMetrics health check failed"
pkill -f "port-forward.*8428"

echo "4. Checking metrics count..."
kubectl port-forward pod/$VM_POD 8428:8428 -n $NAMESPACE &
sleep 5
METRIC_COUNT=$(curl -s "http://localhost:8428/api/v1/query?query=count(count%20by%20(__name__)({__name__!=\"\"}))" | jq -r '.data.result[0].value[1]')
echo "Total metrics: $METRIC_COUNT"
pkill -f "port-forward.*8428"

echo "=== Health check complete ==="
```

Run with:
```bash
chmod +x monitoring-health-check.sh
./monitoring-health-check.sh
```

## Next Steps

After successful verification:

1. **[Configuration](../configuration.md)** - Customize monitoring setup
2. **[Component Configuration](components/)** - Fine-tune individual components  
3. **[Troubleshooting](../troubleshooting.md)** - Handle common issues
4. **[Maintenance](../maintenance.md)** - Ongoing operations 