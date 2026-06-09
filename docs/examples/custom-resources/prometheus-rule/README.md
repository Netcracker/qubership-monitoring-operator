# Service with alarms

* [Service with alarms](#service-with-alarms)
  * [PrometheusRule](#prometheusrule)
  * [Files](#files)
  * [How to apply example](#how-to-apply-example)
  * [Links](#links)

**[Back](../../README.md)**

This example show how to add alarm rule to Prometheus.

## PrometheusRule

```yaml
spec:
  groups:
  - name: general.rules
    rules:
    - alert: TargetDown-service-prom
      annotations:
        description: '{{ $value }}% of {{ $labels.job }} targets are down.'
        summary: Targets are down
      expr: 100 * (count(up == 0) BY (job) / count(up) BY (job)) > 10
      for: 10m
      labels:
        severity: warning
    - alert: DeadMansSwitch-service-prom
      annotations:
        description: This is a DeadMansSwitch meant to ensure that the entire Alerting pipeline is functional.
        summary: Alerting DeadMansSwitch
      expr: vector(1)
      labels:
        severity: none
```

It means that AlertManager will collect data from Prometheus and will check specified expressions. In this example:

* Alert: `Targets are down`
  * Description contains
    templates, [AlertManager Templating](https://prometheus.io/docs/prometheus/latest/configuration/alerting_rules/#templating)
  * Evaluate every 10 minutes (`for: 10m`)
  * Evaluate expression and check that result is true. If result is true alarm will be raise
  * If expression result will true alarm will raise with severity `warning` (`labels.severity: warning`)
* Alert: `Alerting DeadMansSwitch`
  * Simple alarm to verify that AlertManager is working

## Files

* [Prometheus Rule](prometheus-rule.yaml)

## How to apply example

Kubernetes:

```bash
kubectl apply -f prometheus-rule.yaml
```

OpenShift:

```bash
oc apply -f prometheus-rule.yaml
```

## Links

* Prometheus operator
  * [API Documentation](https://github.com/prometheus-operator/prometheus-operator/blob/v0.79.2/Documentation/api.md)
* [Alerting Rules](https://prometheus.io/docs/prometheus/latest/configuration/alerting_rules/)
* Victoriametrics operator
  * [API Documentation](https://docs.victoriametrics.com/operator/api.html)
