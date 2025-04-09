# Simple PodMonitor

* [Simple PodMonitor](#simple-podmonitor)
  * [Files](#files)
  * [How to apply example](#how-to-apply-example)
  * [Links](#links)

**[Back](../../README.md)**

This example show a basic configuration for PodMonitor.

```yaml
spec:
  podMetricsEndpoints:
  - interval: 30s
    targetPort: 14269
  jobLabel: k8s-app
  selector:
    matchLabels:
      k8s-app: sample-service
```

It means that Prometheus/Victoriametrics will collect metrics from pods with settings:

* metrics will collect with job with label `k8s-app`
* metrics will collect from all pods with label `k8s-app: sample-service`
* metrics will collect from all discovered pod from port `14269` with interval `30s`

## Files

* [Pod monitor](pod-monitor.yaml)

## How to apply example

Kubernetes:

```bash
kubectl apply -f pod-monitor.yaml
```

OpenShift:

```bash
oc apply -f pod-monitor.yaml
```

## Links

Prometheus operator API

* [API Documentation](https://github.com/prometheus-operator/prometheus-operator/blob/v0.79.2/Documentation/api.md)

Victoriametrics operator API

* [API Documentation](https://docs.victoriametrics.com/operator/api.html)
