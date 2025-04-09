# Simple AlertmanagerConfig

* [Simple AlertmanagerConfig](#simple-alertmanagerconfig)
  * [Config](#config)
  * [Files](#files)
  * [How to apply example](#how-to-apply-example)
  * [Links](#links)

**[Back](../../README.md)**

This example shows how to add config into Alertmanager for sending metrics and grouping them.

## Config

```yaml
spec:
  route:
    groupBy: ['job']
    groupWait: 30s
    groupInterval: 5m
    repeatInterval: 12h
    receiver: 'wechat-example'
  receivers:
  - name: 'wechat-example'
    wechatConfigs:
    - apiURL: 'http://wechatserver:8080/'
      corpID: 'wechat-corpid'
      apiSecret:
        name: 'wechat-config'
        key: 'apiSecret'
```

It means that the Alertmanager will:

* group alerts by label `job`
* before group alert wait `30s` to be sure that alerts will not close quickly
* group by interval `5n`
* repeat alert (send new notification about raised alerts) every `12h`
* send alerts to receiver with name `wechat-example`
* create receiver with name `wechat-example` and specified settings

## Files

* [Simple AlertmanagerConfig](alertmanager-config.yaml)

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
  * [API Documentation](https://github.com/prometheus-operator/prometheus-operator/blob/v0.79.2/Documentation/api.md)
  * [Alerting user guide](https://github.com/prometheus-operator/prometheus-operator/blob/v0.79.2/Documentation/user-guides/alerting.md)
* Victoriametrics operator
  * [API Documentation](https://docs.victoriametrics.com/operator/api.html)
