# Service with dashboard

* [Service with dashboard](#service-with-dashboard)
  * [ServiceMonitor](#servicemonitor)
  * [Dashboard](#dashboard)
  * [Files](#files)
  * [How to apply example](#how-to-apply-example)
  * [Links](#links)

**[Back](../../README.md)**

This example show how to setup Monitoring for collect metrics from service and dashboard for show collected metrics.

## ServiceMonitor

```yaml
spec:
  endpoints:
  - interval: 30s
    port: http
  jobLabel: k8s-app
  selector:
    matchLabels:
      k8s-app: sample-service
```

It means that Prometheus will collect metrics from service with settings:

* metrics will collect with job with label `k8s-app`
* metrics will collect from all pods with label `k8s-app: sample-service`
* metrics will collect from all discovered pod from port with name `http` with interval `30s`

## Dashboard

```yaml
spec:
  name: simple-dashboard.json
  json: >
    {
      "id": null,
      "title": "Simple Dashboard",
      ...
      <dashboard content in json format>
      ...
    }
```

It means that in grafana will import dashboard with name `Simple Dashboard`.

## Files

* [Service Monitor](service-monitor.yaml)
* [Dashboard](dashboard.yaml)

## How to apply example

Kubernetes:

```bash
kubectl apply -f service-monitor.yaml
kubectl apply -f dashboard.yaml
```

OpenShift:

```bash
oc apply -f service-monitor.yaml
oc apply -f dashboard.yaml
```

## Links

* Prometheus operator
  * [API Documentation](https://github.com/prometheus-operator/prometheus-operator/blob/v0.79.2/Documentation/api.md)
* Grafana operator
  * [Dashboard API](https://github.com/grafana/grafana-operator/blob/v4/documentation/dashboards.md)
* Victoriametrics operator
  * [API Documentation](https://docs.victoriametrics.com/operator/api.html)
