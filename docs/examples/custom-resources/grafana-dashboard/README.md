# Simple Grafana dashboard

* [Simple Grafana dashboard](#simple-grafana-dashboard)
  * [Dashboard](#dashboard)
  * [Files](#files)
  * [How to apply example](#how-to-apply-example)
  * [Links](#links)

**[Back](../../README.md)**

This example show how to add Grafana dashboard for show collected metrics.

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

* [Simple Dashboard](simple-dashboard.yaml)

## How to apply example

Kubernetes:

```bash
kubectl apply -f dashboard.yaml
```

OpenShift:

```bash
oc apply -f dashboard.yaml
```

## Links

* Prometheus operator
  * [API Documentation](https://github.com/prometheus-operator/prometheus-operator/blob/master/Documentation/api.md)
* Grafana operator
  * [Dashboard API](https://github.com/integr8ly/grafana-operator/blob/master/documentation/dashboards.md)
* Victoriametrics operator
  * [API Documentation](https://docs.victoriametrics.com/operator/api.html)
