# ServiceMonitor with custom endpoint

* [ServiceMonitor with custom endpoint](#servicemonitor-with-custom-endpoint)
  * [PodMonitor](#podmonitor)
  * [Files](#files)
  * [How to apply example](#how-to-apply-example)
  * [Links](#links)

**[Back](../../README.md)**

This example show a configuration of ServiceMonitor with custom-endpoint.

## PodMonitor

```yaml
spec:
  endpoints:
  - interval: 30s
    port: https
    scheme: https
    path: /v1/my/awesome/endpoint
    params:
        format:
          - prometheus
  jobLabel: k8s-app
  selector:
    matchLabels:
      k8s-app: sample-service
```

It means that Prometheus will collect metrics from service with settings:

* metrics will collect with job with label `k8s-app`
* metrics will collect from all pods with label `k8s-app: sample-service`
* metrics will collect from all discovered pod from port with name `https` with interval `30s`
* metrics will collect from endpoint `/v1/my/awesome/endpoint`
* request for collect metrics will contains URL parameters: `format=prometheus`

## Files

* [Service Monitor](service-monitor.yaml)

## How to apply example

Kubernetes:

```bash
kubectl apply -f service-monitor.yaml
```

OpenShift:

```bash
oc apply -f service-monitor.yaml
```

## Links

Prometheus operator API

* [API Documentation](https://github.com/prometheus-operator/prometheus-operator/blob/v0.79.2/Documentation/api.md)

Victoriametrics operator API

* [API Documentation](https://docs.victoriametrics.com/operator/api.html)
