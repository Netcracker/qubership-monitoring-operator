# Hardware sizing guide

This document describes how to calculate resources requirements for Platform Monitoring components.

The Platform Monitoring application consists of a set of different components each of it requires different amount of
hardware resources for smooth work. Some components require permanent amount of resources the other requires
resources depends on system loads.
This document describes overall amount of hardware resources for a whole application and detailed information for each
component.

## Components of Platform Monitoring stack

Platform Monitoring components can be combined into some logical groups. The main groups it is:

* Count of metrics
* Cloud size
* Count of CR/CRD
* Count of users/requests

### 1. Count of metrics

Components in this group depend on count of metrics per minute. Resource consumption increases as the amount of metrics
are increased that need to be collected and processed. Points per minute, samples per minute, active targets, etc. -
affect to resource usage.

* [prometheus](#prometheus)
* [vmSingle](#vmsingle)
* [vmAgent](#vmagent)
* [prometheus-adapter](#prometheus-adapter)
* [graphite-remote-adapter](#graphite-remote-adapter)
* [cloudwatch-exporter](#cloudwatch-exporter)
* [blackbox-exporter](#blackbox-exporter)
* [stackdriver-exporter](#stackdriver-exporter)
* [pushgateway](#pushgateway)
* [promitor-agent-scrape](#promitor-agent-scrape)
* [json-exporter](#json-exporter)

### 2. Cloud size

Components in this group depend on count of object in cloud. Resource consumption increases on cloud with a lot of pods,
nodes, configMaps, secrets, etc.

* [kubeStateMetrics](#kubestatemetrics)
* [nodeExporter](#nodeexporter)
* [cert-exporter](#cert-exporter)
* [network-latency-exporter](#network-latency-exporter)
* [version-exporter](#version-exporter)
* [alertManager](#alertmanager)
* [vmAlertManager](#vmalertmanager)
* [vmAlert](#vmalert)

### 3. Count of CR/CRD

Components in this group depend on count of CR/CRD. Resource consumption increases when components have to compute a lot
of CR/CRD. The main part of components in this group it is operators. Operators work with CRD and create CR.

* [monitoring-operator](#monitoring-operator)
* [grafana-operator](#grafana-operator)
* [prometheus-operator](#prometheus-operator)
* [prometheus-adapter-operator](#prometheus-adapter-operator)
* [vmOperator](#vmoperator)
* [configuration-streamer](#configuration-streamer)

### 4. Count of users/requests

Components in this group have UI or handle requests. For example, grafana have UI and handle requests which view on
dashboards.

* [grafana](#grafana)
* [promxy](#promxy)
* [vmAuth](#vmauth)
* oAuth2-proxy

## HWE profiles

Platform Monitoring has several profiles which can be used - `small`, `medium` and `large`. Set `global.profile`
parameter to using one of these profiles else will be used `medium` values for each component of monitoring stack.

|            | Nodes | Pods     | Points per minute |
|----------- | ----- | -------- | ------------------|
| **Small**  | 1-6   | less 100 | less 1Mil         |
| **Medium** | 6-15  | 100-500  | 1-3Mil            |
| **Large**  | 15+   | 500+     | 3+Mil             |

* The `Small` profile is suitable for small cloud with nodes less 6, pods less 100 and points per minute less than 1Mil.
* The `Medium` profile is suitable for cloud with 6-15 nodes, 100-500 pods and points per minute about 3Mil.
* The `Large` profile is suitable for huge cloud with big count of nodes, pods and metrics. Nodes more than 15,
  pods more 500 and points per minute can be 10Mil.

Also, you can specify `resource` parameter for one or more components to override value from profile that uses for
deploy.

Example, overriding resources for monitoring-operator and prometheus, but resources for other components will be used
from `small` profile.

```yaml
global:
  profile: "small"
monitoringOperator:
  resources:
    limits:
      cpu: 100m
      memory: 150Mi
    requests:
      cpu: 50m
      memory: 50Mi
prometheus:
  install: true
  resources:
    requests:
      cpu: 1000m
      memory: 2Gi
    limits:
      cpu: 3000m
      memory: 8Gi
```

**NOTE:** These profiles don't guarantee a stable work of each component. You can increase/override resource parameter
if it needs. Our profiles can't cover all cases for all clouds. [Examples](#examples) with resource usage in different
clouds.

**NOTE:** If you do not set `profile` parameter value for resources will set such as `medium` value of resources.

## Hardware sizing

### monitoring-operator

The `monitoring-operator` is deploying as a single Pod with one Container. Hardware sizing for this service is a
constant and does not depend on system configuration because all handled resources are processing consistently.

<table>
  <tr>
    <th rowspan="2" colspan="2"></th>
    <th colspan="3">Profiles</th>
  </tr>
  <tr>
    <td><strong>Small</strong></td>
    <td><strong>Medium</strong></td>
    <td><strong>Large</strong></td>
  </tr>
  <tr>
    <th rowspan="2">CPU</th>
    <td><strong>requests</strong></td>
    <td>50m</td>
    <td>50m</td>
    <td>70m</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>70m</td>
    <td>100m</td>
    <td>200m</td>
  </tr>
  <tr>
    <th rowspan="2">RAM</th>
    <td><strong>requests</strong></td>
    <td>64Mi</td>
    <td>64Mi</td>
    <td>64Mi</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>256Mi</td>
    <td>256Mi</td>
    <td>256Mi</td>
  </tr>
</table>

You can override `resources` parameter for monitoring-operator:

```yaml
monitoringOperator:
  resources:
    requests:
      cpu: 50m
      memory: 50Mi
    limits:
      cpu: 100m
      memory: 100Mi
```

### prometheus

The `prometheus` collects metrics from configured targets at given intervals, evaluates rule expressions,
displays the results. Resource usage depends on count of metrics that prometheus have to scrape and compute.

<table>
  <tr>
    <th rowspan="2" colspan="2"></th>
    <th colspan="3">Profiles</th>
  </tr>
  <tr>
    <td><strong>Small</strong></td>
    <td><strong>Medium</strong></td>
    <td><strong>Large</strong></td>
  </tr>
  <tr>
    <th rowspan="2">CPU</th>
    <td><strong>requests</strong></td>
    <td>1000m</td>
    <td>2000m</td>
    <td>2500m</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>2000m</td>
    <td>3500m</td>
    <td>4000m</td>
  </tr>
  <tr>
    <th rowspan="2">RAM</th>
    <td><strong>requests</strong></td>
    <td>4Gi</td>
    <td>7Gi</td>
    <td>15Gi</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>6Gi</td>
    <td>12Gi</td>
    <td>25Gi</td>
  </tr>
</table>

You can override `resources` parameter for prometheus:

```yaml
prometheus:
  resources:
    requests:
      cpu: 1000m
      memory: 2Gi
    limits:
      cpu: 2000m
      memory: 3Gi
```

### prometheus-operator

The `prometheus-operator` provides Kubernetes native deployment and management of prometheus and related monitoring
components. Resource usage for prometheus-operator depends on count of prometheus custom resources - Prometheus,
Alertmanager, ThanosRuler, ServiceMonitor, PodMonitor, Probe, PrometheusRule and AlertmanagerConfig.

<table>
  <tr>
    <th rowspan="2" colspan="2"></th>
    <th colspan="3">Profiles</th>
  </tr>
  <tr>
    <td><strong>Small</strong></td>
    <td><strong>Medium</strong></td>
    <td><strong>Large</strong></td>
  </tr>
  <tr>
    <th rowspan="2">CPU</th>
    <td><strong>requests</strong></td>
    <td>30m</td>
    <td>50m</td>
    <td>50m</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>100m</td>
    <td>100m</td>
    <td>100m</td>
  </tr>
  <tr>
    <th rowspan="2">RAM</th>
    <td><strong>requests</strong></td>
    <td>100Mi</td>
    <td>150Mi</td>
    <td>150Mi</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>250Mi</td>
    <td>250Mi</td>
    <td>300Mi</td>
  </tr>
</table>

You can override `resources` parameter for prometheus-operator:

```yaml
prometheus:
  operator:
    resources:
      requests:
        cpu: 50m
        memory: 50Mi
      limits:
        cpu: 100m
        memory: 250Mi
```

### prometheus-adapter

The `prometheus-adapter` is therefore suitable for use with the autoscaling/v2 Horizontal Pod Autoscaler in
Kubernetes 1.6+. It can also replace the metrics server on clusters that already run Prometheus and collect the
appropriate metrics. Resource usage depends on amount of metrics.

<table>
  <tr>
    <th rowspan="2" colspan="2"></th>
    <th colspan="3">Profiles</th>
  </tr>
  <tr>
    <td><strong>Small</strong></td>
    <td><strong>Medium</strong></td>
    <td><strong>Large</strong></td>
  </tr>
  <tr>
    <th rowspan="2">CPU</th>
    <td><strong>requests</strong></td>
    <td>150m</td>
    <td>400m</td>
    <td>500m</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>250m</td>
    <td>500m</td>
    <td>700m</td>
  </tr>
  <tr>
    <th rowspan="2">RAM</th>
    <td><strong>requests</strong></td>
    <td>1000Mi</td>
    <td>2000Mi</td>
    <td>3000Mi</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>2000Mi</td>
    <td>3000Mi</td>
    <td>5000Mi</td>
  </tr>
</table>

You can override `resources` parameter for prometheus-adapter:

```yaml
prometheusAdapter:
  resources:
    requests:
      cpu: 100m
      memory: 256Mi
    limits:
      cpu: 200m
      memory: 384Mi
```

### prometheus-adapter-operator

The `prometheus-adapter-operator` provides Kubernetes native deployment and management of prometheus-adapter
and related monitoring components.

<table>
  <tr>
    <th rowspan="2" colspan="2"></th>
    <th colspan="3">Profiles</th>
  </tr>
  <tr>
    <td><strong>Small</strong></td>
    <td><strong>Medium</strong></td>
    <td><strong>Large</strong></td>
  </tr>
  <tr>
    <th rowspan="2">CPU</th>
    <td><strong>requests</strong></td>
    <td>20m</td>
    <td>30m</td>
    <td>50m</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>50m</td>
    <td>70m</td>
    <td>100m</td>
  </tr>
  <tr>
    <th rowspan="2">RAM</th>
    <td><strong>requests</strong></td>
    <td>20Mi</td>
    <td>30Mi</td>
    <td>30Mi</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>50Mi</td>
    <td>70Mi</td>
    <td>100Mi</td>
  </tr>
</table>

You can override `resources` parameter for prometheus-adapter:

```yaml
prometheusAdapter:
  operator:
    resources:
      requests:
        cpu: 20m
        memory: 20Mi
      limits:
        cpu: 50m
        memory: 100Mi
```

### vmOperator

The `vmOperator` allows you to manage VictoriaMetrics applications inside kubernetes cluster and simplifies this
process. It installs, upgrades and manages victoria metrics resources. Resources usage depends on count of victoria
metrics custom resources. E.g. VMServiceScrape, VMPodScrape, VMProbe, etc.

<table>
  <tr>
    <th rowspan="2" colspan="2"></th>
    <th colspan="3">Profiles</th>
  </tr>
  <tr>
    <td><strong>Small</strong></td>
    <td><strong>Medium</strong></td>
    <td><strong>Large</strong></td>
  </tr>
  <tr>
    <th rowspan="2">CPU</th>
    <td><strong>requests</strong></td>
    <td>50m</td>
    <td>70m</td>
    <td>150m</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>100m</td>
    <td>150m</td>
    <td>300m</td>
  </tr>
  <tr>
    <th rowspan="2">RAM</th>
    <td><strong>requests</strong></td>
    <td>100Mi</td>
    <td>150Mi</td>
    <td>300Mi</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>200Mi</td>
    <td>300Mi</td>
    <td>500Mi</td>
  </tr>
</table>

You can override `resources` parameter for vmOperator:

```yaml
victoriametrics:
  vmOperator:
    resources:
      requests:
        cpu: 200m
        memory: 100Mi
      limits:
        cpu: 400m
        memory: 200Mi
```

### vmAgent

The `vmAgent` is a tiny agent which helps you collect metrics from various sources, relabel and filter the collected
metrics and store them in VictoriaMetrics. Resource usage depends on count of metrics that vmAgent have to collect and
compute.

<table>
  <tr>
    <th rowspan="2" colspan="2"></th>
    <th colspan="3">Profiles</th>
  </tr>
  <tr>
    <td><strong>Small</strong></td>
    <td><strong>Medium</strong></td>
    <td><strong>Large</strong></td>
  </tr>
  <tr>
    <th rowspan="2">CPU</th>
    <td><strong>requests</strong></td>
    <td>100m</td>
    <td>500m</td>
    <td>1500m</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>200m</td>
    <td>750m</td>
    <td>2000m</td>
  </tr>
  <tr>
    <th rowspan="2">RAM</th>
    <td><strong>requests</strong></td>
    <td>100Mi</td>
    <td>400Mi</td>
    <td>2000Mi</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>250Mi</td>
    <td>800Mi</td>
    <td>3500Mi</td>
  </tr>
</table>

You can override `resources` parameter for vmAgent:

```yaml
victoriametrics:
  vmAgent:
    resources:
      requests:
        cpu: 200m
        memory: 100Mi
      limits:
        cpu: 400m
        memory: 200Mi
```

### vmSingle

The `vmSingle` is TSDB and resource usage depends on amount of metrics which have to write/read.

<table>
  <tr>
    <th rowspan="2" colspan="2"></th>
    <th colspan="3">Profiles</th>
  </tr>
  <tr>
    <td><strong>Small</strong></td>
    <td><strong>Medium</strong></td>
    <td><strong>Large</strong></td>
  </tr>
  <tr>
    <th rowspan="2">CPU</th>
    <td><strong>requests</strong></td>
    <td>300m</td>
    <td>1000m</td>
    <td>2000m</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>600m</td>
    <td>1500m</td>
    <td>3000m</td>
  </tr>
  <tr>
    <th rowspan="2">RAM</th>
    <td><strong>requests</strong></td>
    <td>1Gi</td>
    <td>3Gi</td>
    <td>7Gi</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>1.5Gi</td>
    <td>5Gi</td>
    <td>10Gi</td>
  </tr>
</table>

You can override `resources` parameter for vmSingle:

```yaml
victoriametrics:
  vmSingle:
    resources:
      requests:
        cpu: 500m
        memory: 1000Mi
      limits:
        cpu: 1000m
        memory: 2000Mi
```

### vmAlert

The `vmAlert` executes a list of the given alerting or recording rules. Resource usage depends on count of rules on
cloud.

<table>
  <tr>
    <th rowspan="2" colspan="2"></th>
    <th colspan="3">Profiles</th>
  </tr>
  <tr>
    <td><strong>Small</strong></td>
    <td><strong>Medium</strong></td>
    <td><strong>Large</strong></td>
  </tr>
  <tr>
    <th rowspan="2">CPU</th>
    <td><strong>requests</strong></td>
    <td>50m</td>
    <td>100m</td>
    <td>250m</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>100m</td>
    <td>150m</td>
    <td>400m</td>
  </tr>
  <tr>
    <th rowspan="2">RAM</th>
    <td><strong>requests</strong></td>
    <td>150Mi</td>
    <td>250Mi</td>
    <td>400Mi</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>200Mi</td>
    <td>400Mi</td>
    <td>700Mi</td>
  </tr>
</table>

You can override `resources` parameter for vmAlert:

```yaml
victoriametrics:
  vmAlert:
    resources:
      requests:
        cpu: 50m
        memory: 200Mi
      limits:
        cpu: 200m
        memory: 500Mi
```

### vmAlertManager

The `vmAlertManager` is deployment which uses alertmanager for handles alerts sent by client applications.

<table>
  <tr>
    <th rowspan="2" colspan="2"></th>
    <th colspan="3">Profiles</th>
  </tr>
  <tr>
    <td><strong>Small</strong></td>
    <td><strong>Medium</strong></td>
    <td><strong>Large</strong></td>
  </tr>
  <tr>
    <th rowspan="2">CPU</th>
    <td><strong>requests</strong></td>
    <td>30m</td>
    <td>100m</td>
    <td>150m</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>70m</td>
    <td>150m</td>
    <td>200m</td>
  </tr>
  <tr>
    <th rowspan="2">RAM</th>
    <td><strong>requests</strong></td>
    <td>50Mi</td>
    <td>100Mi</td>
    <td>150Mi</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>100Mi</td>
    <td>150Mi</td>
    <td>200Mi</td>
  </tr>
</table>

You can override `resources` parameter for vmAlertManager:

```yaml
victoriametrics:
  vmAlertManager:
    resources:
      requests:
        cpu: 30m
        memory: 56Mi
      limits:
        cpu: 100m
        memory: 256Mi
```

### vmAuth

The `vmAuth` is a simple auth proxy, router and load balancer for VictoriaMetrics. Resource usage increases with
increasing users and requests to proxy.

<table>
  <tr>
    <th rowspan="2" colspan="2"></th>
    <th colspan="3">Profiles</th>
  </tr>
  <tr>
    <td><strong>Small</strong></td>
    <td><strong>Medium</strong></td>
    <td><strong>Large</strong></td>
  </tr>
  <tr>
    <th rowspan="2">CPU</th>
    <td><strong>requests</strong></td>
    <td>50m</td>
    <td>100m</td>
    <td>200m</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>100m</td>
    <td>200m</td>
    <td>350m</td>
  </tr>
  <tr>
    <th rowspan="2">RAM</th>
    <td><strong>requests</strong></td>
    <td>100Mi</td>
    <td>150Mi</td>
    <td>250Mi</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>200Mi</td>
    <td>250Mi</td>
    <td>400Mi</td>
  </tr>
</table>

You can override `resources` parameter for vmAuth:

```yaml
victoriametrics:
  vmAuth:
    resources:
      requests:
        cpu: 50m
        memory: 200Mi
      limits:
        cpu: 200m
        memory: 500Mi
```

### vmSelect

The `vmSelect` is a cluster mode VictoriaMetrics TSDB instance used for reading data. Resource usage depends on
the amount of metrics which has to be read.

<table>
  <tr>
    <th rowspan="2" colspan="2"></th>
    <th colspan="3">Profiles</th>
  </tr>
  <tr>
    <td><strong>Small</strong></td>
    <td><strong>Medium</strong></td>
    <td><strong>Large</strong></td>
  </tr>
  <tr>
    <th rowspan="2">CPU</th>
    <td><strong>requests</strong></td>
    <td>30m</td>
    <td>100m</td>
    <td>150m</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>70m</td>
    <td>150m</td>
    <td>200m</td>
  </tr>
  <tr>
    <th rowspan="2">RAM</th>
    <td><strong>requests</strong></td>
    <td>50Mi</td>
    <td>100Mi</td>
    <td>150Mi</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>100Mi</td>
    <td>150Mi</td>
    <td>200Mi</td>
  </tr>
</table>

You can override `resources` parameter for vmSelect:

```yaml
victoriametrics:
  vmSelect:
    resources:
      requests:
        cpu: 50m
        memory: 64Mi
      limits:
        cpu: 40m
        memory: 64Mi
```

### vmInsert

The `vmInsert` is a cluster mode VictoriaMetrics TSDB instance used for writing data. Resource usage depends on
the amount of metrics which has to be written.

<table>
  <tr>
    <th rowspan="2" colspan="2"></th>
    <th colspan="3">Profiles</th>
  </tr>
  <tr>
    <td><strong>Small</strong></td>
    <td><strong>Medium</strong></td>
    <td><strong>Large</strong></td>
  </tr>
  <tr>
    <th rowspan="2">CPU</th>
    <td><strong>requests</strong></td>
    <td>30m</td>
    <td>100m</td>
    <td>150m</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>70m</td>
    <td>150m</td>
    <td>200m</td>
  </tr>
  <tr>
    <th rowspan="2">RAM</th>
    <td><strong>requests</strong></td>
    <td>50Mi</td>
    <td>100Mi</td>
    <td>150Mi</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>100Mi</td>
    <td>150Mi</td>
    <td>200Mi</td>
  </tr>
</table>

You can override `resources` parameter for vmInsert:

```yaml
victoriametrics:
  vmInsert:
    resources:
      requests:
        cpu: 50m
        memory: 64Mi
      limits:
        cpu: 40m
        memory: 64Mi
```

### vmStorage

The `vmStorage` is a cluster mode VictoriaMetrics TSDB instance used for store data. Resource usage depends on
the amount of metrics which has to be stored.

<table>
  <tr>
    <th rowspan="2" colspan="2"></th>
    <th colspan="3">Profiles</th>
  </tr>
  <tr>
    <td><strong>Small</strong></td>
    <td><strong>Medium</strong></td>
    <td><strong>Large</strong></td>
  </tr>
  <tr>
    <th rowspan="2">CPU</th>
    <td><strong>requests</strong></td>
    <td>300m</td>
    <td>500m</td>
    <td>1000m</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>300m</td>
    <td>500m</td>
    <td>1000m</td>
  </tr>
  <tr>
    <th rowspan="2">RAM</th>
    <td><strong>requests</strong></td>
    <td>256Mi</td>
    <td>512Mi</td>
    <td>1024Mi</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>256Mi</td>
    <td>512Mi</td>
    <td>1024Mi</td>
  </tr>
</table>

You can override `resources` parameter for vmInsert:

```yaml
victoriametrics:
  vmStorage:
    resources:
      requests:
        cpu: 500m
        memory: 512Mi
      limits:
        cpu: 500m
        memory: 512Mi
```

### grafana

The `grafana` queries, alerts and  visualizes metrics which was collected by prometheus or victoriaMetrics. The number
and complexity of dashboards affects resource usage.

<table>
  <tr>
    <th rowspan="2" colspan="2"></th>
    <th colspan="3">Profiles</th>
  </tr>
  <tr>
    <td><strong>Small</strong></td>
    <td><strong>Medium</strong></td>
    <td><strong>Large</strong></td>
  </tr>
  <tr>
    <th rowspan="2">CPU</th>
    <td><strong>requests</strong></td>
    <td>250m</td>
    <td>400m</td>
    <td>700m</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>400m</td>
    <td>500m</td>
    <td>900m</td>
  </tr>
  <tr>
    <th rowspan="2">RAM</th>
    <td><strong>requests</strong></td>
    <td>200Mi</td>
    <td>350Mi</td>
    <td>600Mi</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>300Mi</td>
    <td>450Mi</td>
    <td>700Mi</td>
  </tr>
</table>

You can override `resources` parameter for grafana:

```yaml
grafana:
  resources:
    requests:
      cpu: 300m
      memory: 400Mi
    limits:
      cpu: 500m
      memory: 800Mi
```

### grafana-operator

The `grafana-operator` is a Kubernetes operator built to help you manage your Grafana instances in and outside of
Kubernetes. Resource usage for grafana-operator depends on count of grafana custom resources - Grafana, GrafanaDashboard,
GrafanaDataSource, GrafanaFolder, GrafanaNotificationChannel.

<table>
  <tr>
    <th rowspan="2" colspan="2"></th>
    <th colspan="3">Profiles</th>
  </tr>
  <tr>
    <td><strong>Small</strong></td>
    <td><strong>Medium</strong></td>
    <td><strong>Large</strong></td>
  </tr>
  <tr>
    <th rowspan="2">CPU</th>
    <td><strong>requests</strong></td>
    <td>30m</td>
    <td>50m</td>
    <td>150m</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>70m</td>
    <td>100m</td>
    <td>250m</td>
  </tr>
  <tr>
    <th rowspan="2">RAM</th>
    <td><strong>requests</strong></td>
    <td>50Mi</td>
    <td>150Mi</td>
    <td>200Mi</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>100Mi</td>
    <td>250Mi</td>
    <td>350Mi</td>
  </tr>
</table>

You can override `resources` parameter for grafana-operator:

```yaml
grafana:
  operator:
    resources:
      requests:
        cpu: 50m
        memory: 50Mi
      limits:
        cpu: 100m
        memory: 100Mi
```

### grafana-image-renderer

The `grafana-image-renderer` handles rendering panels and dashboards to PNGs using a headless browser (Chromium).
Rendering images requires a lot of memory, mainly because Grafana creates browser instances in the background for the
actual rendering. We recommend a minimum of 16GB of free memory on the system rendering images on clouds with big count
of dashboards. Rendering multiple images in parallel requires an even bigger memory footprint.

<table>
  <tr>
    <th rowspan="2" colspan="2"></th>
    <th colspan="3">Profiles</th>
  </tr>
  <tr>
    <td><strong>Small</strong></td>
    <td><strong>Medium</strong></td>
    <td><strong>Large</strong></td>
  </tr>
  <tr>
    <th rowspan="2">CPU</th>
    <td><strong>requests</strong></td>
    <td>100m</td>
    <td>300m</td>
    <td>500m</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>200m</td>
    <td>500m</td>
    <td>800m</td>
  </tr>
  <tr>
    <th rowspan="2">RAM</th>
    <td><strong>requests</strong></td>
    <td>200Mi</td>
    <td>500Mi</td>
    <td>1000Mi</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>400Mi</td>
    <td>800Mi</td>
    <td>2000Mi</td>
  </tr>
</table>

You can override `resources` parameter for grafana-image-renderer:

```yaml
grafana:
  imageRenderer:
    resources:
      requests:
        cpu: 150m
        memory: 250Mi
      limits:
        cpu: 300m
        memory: 500Mi
```

### alertManager

The `alertmanager` handles alerts sent by client applications such as the Prometheus server.
It takes care of deduplicating, grouping, and routing them to the correct receiver integration. Resource usage depends
on alert which need to watch and compute.

<table>
  <tr>
    <th rowspan="2" colspan="2"></th>
    <th colspan="3">Profiles</th>
  </tr>
  <tr>
    <td><strong>Small</strong></td>
    <td><strong>Medium</strong></td>
    <td><strong>Large</strong></td>
  </tr>
  <tr>
    <th rowspan="2">CPU</th>
    <td><strong>requests</strong></td>
    <td>50m</td>
    <td>70m</td>
    <td>150m</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>100m</td>
    <td>120m</td>
    <td>200m</td>
  </tr>
  <tr>
    <th rowspan="2">RAM</th>
    <td><strong>requests</strong></td>
    <td>50Mi</td>
    <td>100Mi</td>
    <td>200Mi</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>100Mi</td>
    <td>150Mi</td>
    <td>300Mi</td>
  </tr>
</table>

You can override `resources` parameter for alertManager:

```yaml
alertManager:
  resources:
    requests:
      cpu: 100m
      memory: 100Mi
    limits:
      cpu: 200m
      memory: 200Mi
```

### kubeStateMetrics

The `kubeStateMetrics` is a simple service that listens to the Kubernetes API server and generates metrics about the
state of the objects. It is not focused on the health of the individual Kubernetes components, but rather on the health
of the various objects inside, such as deployments, nodes and pods. Resource usage depends on amount of objects on cloud.

<table>
  <tr>
    <th rowspan="2" colspan="2"></th>
    <th colspan="3">Profiles</th>
  </tr>
  <tr>
    <td><strong>Small</strong></td>
    <td><strong>Medium</strong></td>
    <td><strong>Large</strong></td>
  </tr>
  <tr>
    <th rowspan="2">CPU</th>
    <td><strong>requests</strong></td>
    <td>50m</td>
    <td>70m</td>
    <td>100m</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>100m</td>
    <td>150m</td>
    <td>200m</td>
  </tr>
  <tr>
    <th rowspan="2">RAM</th>
    <td><strong>requests</strong></td>
    <td>50Mi</td>
    <td>120Mi</td>
    <td>200Mi</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>100Mi</td>
    <td>200Mi</td>
    <td>300Mi</td>
  </tr>
</table>

You can override `resources` parameter for kubeStateMetrics:

```yaml
kubeStateMetrics:
  resources:
    requests:
      cpu: 50m
      memory: 50Mi
    limits:
      cpu: 100m
      memory: 256Mi
```

### pushgateway

The `pushgateway` exists to allow ephemeral and batch jobs to expose their metrics to Prometheus. Since these kinds of
jobs may not exist long enough to be scraped, they can instead push their metrics to a Pushgateway. The Pushgateway then
exposes these metrics to Prometheus.

<table>
  <tr>
    <th rowspan="2" colspan="2"></th>
    <th colspan="3">Profiles</th>
  </tr>
  <tr>
    <td><strong>Small</strong></td>
    <td><strong>Medium</strong></td>
    <td><strong>Large</strong></td>
  </tr>
  <tr>
    <th rowspan="2">CPU</th>
    <td><strong>requests</strong></td>
    <td>50m</td>
    <td>150m</td>
    <td>250m</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>70m</td>
    <td>250m</td>
    <td>400m</td>
  </tr>
  <tr>
    <th rowspan="2">RAM</th>
    <td><strong>requests</strong></td>
    <td>30Mi</td>
    <td>100Mi</td>
    <td>150Mi</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>50Mi</td>
    <td>150Mi</td>
    <td>250Mi</td>
  </tr>
</table>

You can override `resources` parameter for pushgateway:

```yaml
pushgateway:
  resources:
    requests:
      cpu: 100m
      memory: 30Mi
    limits:
      cpu: 200m
      memory: 50Mi
```

### promxy

The `promxy` is a prometheus proxy that makes many shards of prometheus appear as a single API endpoint to the user.
This significantly simplifies operations and use of prometheus at scale (when you have more than one prometheus host).
Promxy delivers this unified access endpoint without requiring any sidecars, custom-builds, or other changes to your
prometheus infrastructure.

<table>
  <tr>
    <th rowspan="2" colspan="2"></th>
    <th colspan="3">Profiles</th>
  </tr>
  <tr>
    <td><strong>Small</strong></td>
    <td><strong>Medium</strong></td>
    <td><strong>Large</strong></td>
  </tr>
  <tr>
    <th rowspan="2">CPU</th>
    <td><strong>requests</strong></td>
    <td>50m</td>
    <td>100m</td>
    <td>200m</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>100m</td>
    <td>150m</td>
    <td>250m</td>
  </tr>
  <tr>
    <th rowspan="2">RAM</th>
    <td><strong>requests</strong></td>
    <td>100Mi</td>
    <td>200Mi</td>
    <td>300Mi</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>150Mi</td>
    <td>250Mi</td>
    <td>350Mi</td>
  </tr>
</table>

You can override `resources` parameter for promxy:

```yaml
promxy:
  resources:
    requests:
      cpu: 50m
      memory: 128Mi
    limits:
      cpu: 150m
      memory: 256Mi
```

### promxy-configmap-reloader

<table>
  <tr>
    <th rowspan="2" colspan="2"></th>
    <th colspan="3">Profiles</th>
  </tr>
  <tr>
    <td><strong>Small</strong></td>
    <td><strong>Medium</strong></td>
    <td><strong>Large</strong></td>
  </tr>
  <tr>
    <th rowspan="2">CPU</th>
    <td><strong>requests</strong></td>
    <td>10m</td>
    <td>10m</td>
    <td>15m</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>15m</td>
    <td>15m</td>
    <td>20m</td>
  </tr>
  <tr>
    <th rowspan="2">RAM</th>
    <td><strong>requests</strong></td>
    <td>6Mi</td>
    <td>10Mi</td>
    <td>15Mi</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>15Mi</td>
    <td>15Mi</td>
    <td>20Mi</td>
  </tr>
</table>

You can override `resources` parameter for promxy-configmap-reloader:

```yaml
promxy:
  configmapReload:
    resources:
      requests:
        cpu: 5m
        memory: 3Mi
      limits:
        cpu: 10m
        memory: 20Mi
```

### graphite-remote-adapter

The `graphite-remote-adapter` is a read/write adapter that receives samples via Prometheus's remote write protocol and
stores them in remote storage like Graphite. Resource usage depends on count of samples that have to read/write.

<table>
  <tr>
    <th rowspan="2" colspan="2"></th>
    <th colspan="3">Profiles</th>
  </tr>
  <tr>
    <td><strong>Small</strong></td>
    <td><strong>Medium</strong></td>
    <td><strong>Large</strong></td>
  </tr>
  <tr>
    <th rowspan="2">CPU</th>
    <td><strong>requests</strong></td>
    <td>100m</td>
    <td>250m</td>
    <td>500m</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>150m</td>
    <td>400m</td>
    <td>750m</td>
  </tr>
  <tr>
    <th rowspan="2">RAM</th>
    <td><strong>requests</strong></td>
    <td>150Mi</td>
    <td>400Mi</td>
    <td>1000Mi</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>250Mi</td>
    <td>700Mi</td>
    <td>1500Mi</td>
  </tr>
</table>

You can override `resources` parameter for graphite-remote-adapter:

```yaml
graphite_remote_adapter:
  resources:
    requests:
      cpu: 200m
      memory: 300Mi
    limits:
      cpu: 500m
      memory: 1000Mi
```

### promitor-agent-scrape

The `promitor-agent-scrape` is an Azure Monitor scraper which makes the metrics available through a scraping endpoint
for Prometheus and resource usage depends on count of metrics that have to collect and compute.

<table>
  <tr>
    <th rowspan="2" colspan="2"></th>
    <th colspan="3">Profiles</th>
  </tr>
  <tr>
    <td><strong>Small</strong></td>
    <td><strong>Medium</strong></td>
    <td><strong>Large</strong></td>
  </tr>
  <tr>
    <th rowspan="2">CPU</th>
    <td><strong>requests</strong></td>
    <td>70m</td>
    <td>150m</td>
    <td>200m</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>100m</td>
    <td>200m</td>
    <td>400m</td>
  </tr>
  <tr>
    <th rowspan="2">RAM</th>
    <td><strong>requests</strong></td>
    <td>100Mi</td>
    <td>200Mi</td>
    <td>250Mi</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>150Mi</td>
    <td>200Mi</td>
    <td>500Mi</td>
  </tr>
</table>

You can override `resources` parameter for promitor-agent-scrape:

```yaml
promitorAgentScraper:
  resources:
    requests:
      cpu: 100m
      memory: 128Mi
    limits:
      cpu: 200m
      memory: 256Mi
```

### nodeExporter

The `nodeExporter` is Prometheus exporter for hardware and OS metrics exposed by *NIX kernels, written in Go with
pluggable metric collectors. Count and size of nodes to affects on resource usage.

<table>
  <tr>
    <th rowspan="2" colspan="2"></th>
    <th colspan="3">Profiles</th>
  </tr>
  <tr>
    <td><strong>Small</strong></td>
    <td><strong>Medium</strong></td>
    <td><strong>Large</strong></td>
  </tr>
  <tr>
    <th rowspan="2">CPU</th>
    <td><strong>requests</strong></td>
    <td>30m</td>
    <td>50m</td>
    <td>50m</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>50m</td>
    <td>70m</td>
    <td>100m</td>
  </tr>
  <tr>
    <th rowspan="2">RAM</th>
    <td><strong>requests</strong></td>
    <td>30Mi</td>
    <td>50Mi</td>
    <td>50Mi</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>50Mi</td>
    <td>70Mi</td>
    <td>100Mi</td>
  </tr>
</table>

You can override `resources` parameter for nodeExporter:

```yaml
nodeExporter:
  resources:
    requests:
      cpu: 50m
      memory: 50Mi
    limits:
      cpu: 100m
      memory: 100Mi
```

### cert-exporter

The `cert-exporter` is designed to parse certificates and export expiration information for Prometheus to scrape.
Kubernetes uses PKI certificates for authentication between all major components. These certs are critical for the
operation of your cluster but are often opaque to an administrator.

#### DaemonSet

The `cert-exporter` daemonset collects certs from files and/or kubeconfig.

<table>
  <tr>
    <th rowspan="2" colspan="2"></th>
    <th colspan="3">Profiles</th>
  </tr>
  <tr>
    <td><strong>Small</strong></td>
    <td><strong>Medium</strong></td>
    <td><strong>Large</strong></td>
  </tr>
  <tr>
    <th rowspan="2">CPU</th>
    <td><strong>requests</strong></td>
    <td>10m</td>
    <td>20m</td>
    <td>30m</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>20m</td>
    <td>40m</td>
    <td>50m</td>
  </tr>
  <tr>
    <th rowspan="2">RAM</th>
    <td><strong>requests</strong></td>
    <td>20Mi</td>
    <td>30Mi</td>
    <td>50Mi</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>30Mi</td>
    <td>50Mi</td>
    <td>70Mi</td>
  </tr>
</table>

You can override `resources` parameter for cert-exporter daemonset:

```yaml
certExporter:
  daemonset:
    resources:
      requests:
        cpu: 10m
        memory: 25Mi
      limits:
        cpu: 20m
        memory: 50Mi
```

#### Deployment

The `cert-exporter` deployment collects certs from secrets.

<table>
  <tr>
    <th rowspan="2" colspan="2"></th>
    <th colspan="3">Profiles</th>
  </tr>
  <tr>
    <td><strong>Small</strong></td>
    <td><strong>Medium</strong></td>
    <td><strong>Large</strong></td>
  </tr>
  <tr>
    <th rowspan="2">CPU</th>
    <td><strong>requests</strong></td>
    <td>10m</td>
    <td>20m</td>
    <td>30m</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>20m</td>
    <td>30m</td>
    <td>50m</td>
  </tr>
  <tr>
    <th rowspan="2">RAM</th>
    <td><strong>requests</strong></td>
    <td>30Mi</td>
    <td>70Mi</td>
    <td>100Mi</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>50Mi</td>
    <td>150Mi</td>
    <td>200Mi</td>
  </tr>
</table>

You can override `resources` parameter for cert-exporter deployment:

```yaml
certExporter:
  deployment:
    resources:
      requests:
        cpu: 10m
        memory: 50Mi
      limits:
        cpu: 20m
        memory: 150Mi
```

### blackbox-exporter

The `blackbox-exporter` allows blackbox probing of endpoints over HTTP, HTTPS, DNS, TCP, ICMP and gRPC.

<table>
  <tr>
    <th rowspan="2" colspan="2"></th>
    <th colspan="3">Profiles</th>
  </tr>
  <tr>
    <td><strong>Small</strong></td>
    <td><strong>Medium</strong></td>
    <td><strong>Large</strong></td>
  </tr>
  <tr>
    <th rowspan="2">CPU</th>
    <td><strong>requests</strong></td>
    <td>20m</td>
    <td>50m</td>
    <td>100m</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>30m</td>
    <td>70m</td>
    <td>150m</td>
  </tr>
  <tr>
    <th rowspan="2">RAM</th>
    <td><strong>requests</strong></td>
    <td>20Mi</td>
    <td>50Mi</td>
    <td>100Mi</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>50Mi</td>
    <td>100Mi</td>
    <td>250Mi</td>
  </tr>
</table>

You can override `resources` parameter for blackbox-exporter:

```yaml
blackboxExporter:
  resources:
    requests:
      cpu: 50m
      memory: 50Mi
    limits:
      cpu: 100m
      memory: 300Mi
```

### cloudwatch-exporter

The `cloudwatch-exporter` is exporter for Amazon CloudWatch. Count of metrics (points, samples) affects on resource
usage for `cloudwatch-exporter` deployment.

<table>
  <tr>
    <th rowspan="2" colspan="2"></th>
    <th colspan="3">Profiles</th>
  </tr>
  <tr>
    <td><strong>Small</strong></td>
    <td><strong>Medium</strong></td>
    <td><strong>Large</strong></td>
  </tr>
  <tr>
    <th rowspan="2">CPU</th>
    <td><strong>requests</strong></td>
    <td>50m</td>
    <td>100m</td>
    <td>150m</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>70m</td>
    <td>150m</td>
    <td>250m</td>
  </tr>
  <tr>
    <th rowspan="2">RAM</th>
    <td><strong>requests</strong></td>
    <td>100Mi</td>
    <td>150Mi</td>
    <td>200Mi</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>150Mi</td>
    <td>250Mi</td>
    <td>300Mi</td>
  </tr>
</table>

You can override `resources` parameter for cloudwatch-exporter:

```yaml
cloudwatchExporter:
  resources:
    requests:
      cpu: 100m
      memory: 128Mi
    limits:
      cpu: 200m
      memory: 256Mi
```

### json-exporter

The `json-exporter` is a prometheus exporter which scrapes remote JSON by JSONPath.

<table>
  <tr>
    <th rowspan="2" colspan="2"></th>
    <th colspan="3">Profiles</th>
  </tr>
  <tr>
    <td><strong>Small</strong></td>
    <td><strong>Medium</strong></td>
    <td><strong>Large</strong></td>
  </tr>
  <tr>
    <th rowspan="2">CPU</th>
    <td><strong>requests</strong></td>
    <td>50m</td>
    <td>100m</td>
    <td>200m</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>70m</td>
    <td>150m</td>
    <td>300m</td>
  </tr>
  <tr>
    <th rowspan="2">RAM</th>
    <td><strong>requests</strong></td>
    <td>100Mi</td>
    <td>150Mi</td>
    <td>250Mi</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>150Mi</td>
    <td>200Mi</td>
    <td>350Mi</td>
  </tr>
</table>

You can override `resources` parameter for json-exporter:

```yaml
jsonExporter:
  resources:
    requests:
      cpu: 100m
      memory: 128Mi
    limits:
      cpu: 100m
      memory: 128Mi
```

### network-latency-exporter

The `network-latency-exporter` is a service which collects RTT and TTL metrics for the list of target hosts and sends
collected data to Prometheus. It is possible to use UDP, TCP or ICMP network protocols to sent package during probes.
The service collects metrics with `mtr` tool which accumulates functionality of `ping` and `traceroute` tools.
Target hosts can be discovered automatically by retrieving all k8s cluster nodes. Resource usage depends on count of
targets. More targets have more metrics.

<table>
  <tr>
    <th rowspan="2" colspan="2"></th>
    <th colspan="3">Profiles</th>
  </tr>
  <tr>
    <td><strong>Small</strong></td>
    <td><strong>Medium</strong></td>
    <td><strong>Large</strong></td>
  </tr>
  <tr>
    <th rowspan="2">CPU</th>
    <td><strong>requests</strong></td>
    <td>70m</td>
    <td>150m</td>
    <td>200m</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>150m</td>
    <td>250m</td>
    <td>300m</td>
  </tr>
  <tr>
    <th rowspan="2">RAM</th>
    <td><strong>requests</strong></td>
    <td>100Mi</td>
    <td>200Mi</td>
    <td>250Mi</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>200Mi</td>
    <td>300Mi</td>
    <td>350Mi</td>
  </tr>
</table>

You can override `resources` parameter for network-latency-exporter:

```yaml
networkLatencyExporter:
  resources:
    requests:
      cpu: 100m
      memory: 128Mi
    limits:
      cpu: 200m
      memory: 256Mi
```

### stackdriver-exporter

The `stackdriver-exporter` is a Prometheus exporter for Google Stackdriver Monitoring metrics. It acts as a proxy that
requests Stackdriver API for the metric's time-series everytime prometheus scrapes it. Count of metrics (points,
samples) affects on resource usage for `stackdriver-exporter` deployment.

<table>
  <tr>
    <th rowspan="2" colspan="2"></th>
    <th colspan="3">Profiles</th>
  </tr>
  <tr>
    <td><strong>Small</strong></td>
    <td><strong>Medium</strong></td>
    <td><strong>Large</strong></td>
  </tr>
  <tr>
    <th rowspan="2">CPU</th>
    <td><strong>requests</strong></td>
    <td>50m</td>
    <td>100m</td>
    <td>250m</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>100m</td>
    <td>150m</td>
    <td>350m</td>
  </tr>
  <tr>
    <th rowspan="2">RAM</th>
    <td><strong>requests</strong></td>
    <td>70Mi</td>
    <td>150Mi</td>
    <td>300Mi</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>150Mi</td>
    <td>200Mi</td>
    <td>400Mi</td>
  </tr>
</table>

You can override `resources` parameter for stackdriver-exporter:

```yaml
stackdriverExporter:
  resources:
    requests:
      cpu: 100m
      memory: 128Mi
    limits:
      cpu: 100m
      memory: 128Mi
```

### version-exporter

The `version-exporter` is a useful tool that allows you to get product, project, third-party versions of an application
and store the results in custom Prometheus metrics.

<table>
  <tr>
    <th rowspan="2" colspan="2"></th>
    <th colspan="3">Profiles</th>
  </tr>
  <tr>
    <td><strong>Small</strong></td>
    <td><strong>Medium</strong></td>
    <td><strong>Large</strong></td>
  </tr>
  <tr>
    <th rowspan="2">CPU</th>
    <td><strong>requests</strong></td>
    <td>100m</td>
    <td>150m</td>
    <td>200m</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>150m</td>
    <td>200m</td>
    <td>300m</td>
  </tr>
  <tr>
    <th rowspan="2">RAM</th>
    <td><strong>requests</strong></td>
    <td>200Mi</td>
    <td>250Mi</td>
    <td>300Mi</td>
  </tr>
  <tr>
    <td><strong>limits</strong></td>
    <td>250Mi</td>
    <td>300Mi</td>
    <td>400Mi</td>
  </tr>
</table>

You can override `resources` parameter for version-exporter:

```yaml
versionExporter:
  resources:
    requests:
      cpu: 200m
      memory: 256Mi
    limits:
      cpu: 500m
      memory: 512Mi
```

## Examples

Examples show that count of nodes, pods and metrics don't affect all components equally. Each component resource usage
depends on a lot of facts.

These are examples from real clouds. It shows that equal number of pods and nodes doesn't guarantee equal resource usage.

For example, several clouds can be compared:

<table>
  <tr>
    <th rowspan="2"></th>
    <th rowspan="2">Nodes</th>
    <th rowspan="2">Pods</th>
    <th rowspan="2">Active targets</th>
    <th rowspan="2">Dashboards</th>
    <th colspan="2">Monitoring-operator</th>
    <th colspan="2">Prometheus</th>
    <th colspan="2">Grafana</th>
  </tr>
  <tr>
    <th>CPU</th>
    <th>RAM</th>
    <th>CPU</th>
    <th>RAM</th>
    <th>CPU</th>
    <th>RAM</th>
  </tr>
  <tr>
    <td>#1</td>
    <td>6</td>
    <td>1000</td>
    <td>550</td>
    <td>100</td>
    <td>23m</td>
    <td>50Mi</td>
    <td>2200m</td>
    <td>12.5Gi</td>
    <td>270m</td>
    <td>120Mi</td>
  </tr>
  <tr>
    <td>#2</td>
    <td>21</td>
    <td>400</td>
    <td>170</td>
    <td>30</td>
    <td>21m</td>
    <td>56Mi</td>
    <td>450m</td>
    <td>7Gi</td>
    <td>1100m</td>
    <td>240Mi</td>
  </tr>
  <tr>
    <td>#3</td>
    <td>40</td>
    <td>2000</td>
    <td>220</td>
    <td>30</td>
    <td>25m</td>
    <td>70Mi</td>
    <td>530m</td>
    <td>7.5Gi</td>
    <td>200m</td>
    <td>130Mi</td>
  </tr>
</table>

There is we can see absolute different clouds. What is it mean?:

* The #1 cloud has pods more that in #2 cloud, because on the #1 cloud size of node bigger than on the #2.
* On the #1 cloud a lot of active targets and resource usage on this cloud more the on the both other.
* Clouds #2 and #3 has equal amount of dashboards, but really different resource usage for grafana.
  This happened due to the fact that on the #2 cloud there are more complex dashboards(it requires
  more CPU and RAM).
* Resource usage of monitoring-operator almost equal between these clouds because all handled resources are processing
  consistently.
* The prometheus on #2 and #3 cloud use the same amount of resources. Number of active targets and metrics
  is the same on them.

### Example with prometheus

<table>
  <tr>
    <th rowspan="2">Nodes</th>
    <th rowspan="2">Pods</th>
    <th rowspan="2">Active targets</th>
    <th rowspan="2">Dashboards</th>
    <th colspan="2">Monitoring-operator</th>
    <th colspan="2">Prometheus</th>
    <th colspan="2">Prometheus-operator</th>
    <th colspan="2">Grafana</th>
    <th colspan="2">Grafana-operator</th>
    <th colspan="2">AlertManager</th>
    <th colspan="2">KubeStateMetrics</th>
    <th colspan="2">Graphite-remote-adapter</th>
    <th colspan="2">Cert-exporter</th>
    <th colspan="2">CloudWatch-exporter</th>
    <th colspan="2">Network-Latency-exporter</th>
  </tr>
  <tr>
    <th>CPU</th>
    <th>RAM</th>
    <th>CPU</th>
    <th>RAM</th>
    <th>CPU</th>
    <th>RAM</th>
    <th>CPU</th>
    <th>RAM</th>
    <th>CPU</th>
    <th>RAM</th>
    <th>CPU</th>
    <th>RAM</th>
    <th>CPU</th>
    <th>RAM</th>
    <th>CPU</th>
    <th>RAM</th>
    <th>CPU</th>
    <th>RAM</th>
    <th>CPU</th>
    <th>RAM</th>
    <th>CPU</th>
    <th>RAM</th>
  </tr>
  <tr>
    <td>6</td>
    <td>500</td>
    <td>289</td>
    <td>65</td>
    <td>7m</td>
    <td>70Mi</td>
    <td>400m</td>
    <td>10Gi</td>
    <td>4m</td>
    <td>120Mi</td>
    <td>230m</td>
    <td>120Mi</td>
    <td>20m</td>
    <td>100Mi</td>
    <td>6m</td>
    <td>50Mi</td>
    <td>5m</td>
    <td>60Mi</td>
    <td>410m</td>
    <td>3.5Gi</td>
    <td>6m</td>
    <td>100Mi</td>
    <td>5m</td>
    <td>160Mi</td>
    <td>10m</td>
    <td>250Mi</td>
  </tr>
  <tr>
    <td>12</td>
    <td>440</td>
    <td>226</td>
    <td>59</td>
    <td>25m</td>
    <td>60Mi</td>
    <td>1150m</td>
    <td>5.6Gi</td>
    <td>10m</td>
    <td>50Mi</td>
    <td>1150m</td>
    <td>85Mi</td>
    <td>100m</td>
    <td>70Mi</td>
    <td>10m</td>
    <td>30Mi</td>
    <td>5m</td>
    <td>55Mi</td>
    <td>325m</td>
    <td>800Mi</td>
  </tr>
  <tr>
    <td>16</td>
    <td>-</td>
    <td>121</td>
    <td>-</td>
    <td>21m</td>
    <td>105Mi</td>
    <td>290m</td>
    <td>6.2Gi</td>
    <td>15m</td>
    <td>250Mi</td>
    <td>-</td>
    <td>-</td>
    <td>-</td>
    <td>-</td>
    <td>-</td>
    <td>-</td>
    <td>10m</td>
    <td>125Mi</td>
    <td>600m</td>
    <td>950Mi</td>
  </tr>
  <tr>
    <td>19</td>
    <td>-</td>
    <td>153</td>
    <td>-</td>
    <td>22m</td>
    <td>110Mi</td>
    <td>450m</td>
    <td>5.5Gi</td>
    <td>15m</td>
    <td>190Mi</td>
    <td>-</td>
    <td>-</td>
    <td>-</td>
    <td>-</td>
    <td>-</td>
    <td>-</td>
    <td>5m</td>
    <td>100Mi</td>
    <td>360m</td>
    <td>750Mi</td>
  </tr>
  <tr>
    <td>28</td>
    <td>290</td>
    <td>156</td>
    <td>79</td>
    <td>100m</td>
    <td>195Mi</td>
    <td>1600m</td>
    <td>6.4Gi</td>
    <td>5m</td>
    <td>100Mi</td>
    <td>60m</td>
    <td>125Mi</td>
    <td>30m</td>
    <td>115Mi</td>
    <td>10m</td>
    <td>40Mi</td>
    <td>5m</td>
    <td>90Mi</td>
    <td>720m</td>
    <td>850Mi</td>
  </tr>
  <tr>
    <td>37</td>
    <td>228</td>
    <td>129</td>
    <td>30</td>
    <td>20m</td>
    <td>56Mi</td>
    <td>750m</td>
    <td>7.2Gi</td>
    <td>7m</td>
    <td>160Mi</td>
    <td>430m</td>
    <td>165Mi</td>
    <td>20m</td>
    <td>70Mi</td>
    <td>10m</td>
    <td>50Mi</td>
    <td>10m</td>
    <td>140Mi</td>
    <td>450m</td>
    <td>1150Mi</td>
  </tr>
</table>

### Example with victoriaMetrics

<table>
  <tr>
    <th rowspan="2">Nodes</th>
    <th rowspan="2">Pods</th>
    <th rowspan="2">Samples per second</th>
    <th colspan="2">vmOperator</th>
    <th colspan="2">vmAgent</th>
    <th colspan="2">NodeExporter</th>
    <th colspan="2">KubeStateMetrics</th>
  </tr>
  <tr>
    <th>CPU</th>
    <th>RAM</th>
    <th>CPU</th>
    <th>RAM</th>
    <th>CPU</th>
    <th>RAM</th>
    <th>CPU</th>
    <th>RAM</th>
  </tr>
  <tr>
    <td>15</td>
    <td>600</td>
    <td>6K</td>
    <td>100m</td>
    <td>52Mi</td>
    <td>90m</td>
    <td>200Mi</td>
    <td>30m</td>
    <td>16Mi</td>
    <td>10m</td>
    <td>50Mi</td>
  </tr>
  <tr>
    <td>16</td>
    <td>370</td>
    <td>4K</td>
    <td>100m</td>
    <td>53Mi</td>
    <td>70m</td>
    <td>180Mi</td>
    <td>20m</td>
    <td>25Mi</td>
    <td>10m</td>
    <td>35Mi</td>
  </tr>
  <tr>
    <td>26</td>
    <td>2700</td>
    <td>43K</td>
    <td>100m</td>
    <td>320Mi</td>
    <td>470m</td>
    <td>440Mi</td>
    <td>20m</td>
    <td>55Mi</td>
    <td>10m</td>
    <td>140Mi</td>
  </tr>
  <tr>
    <td>41</td>
    <td>3000</td>
    <td>47K</td>
    <td>100m</td>
    <td>180Mi</td>
    <td>480m</td>
    <td>550Mi</td>
    <td>10m</td>
    <td>30Mi</td>
    <td>10m</td>
    <td>230Mi</td>
  </tr>
  <tr>
    <td>57</td>
    <td>3700</td>
    <td>41K</td>
    <td>100m</td>
    <td>200Mi</td>
    <td>380m</td>
    <td>340Mi</td>
    <td>10m</td>
    <td>35Mi</td>
    <td>10m</td>
    <td>180Mi</td>
  </tr>
  <tr>
    <td>65</td>
    <td>4500</td>
    <td>67K</td>
    <td>100m</td>
    <td>270Mi</td>
    <td>710m</td>
    <td>1100Mi</td>
    <td>10m</td>
    <td>19Mi</td>
    <td>40m</td>
    <td>270Mi</td>
  </tr>
  <tr>
    <td>73</td>
    <td>7500</td>
    <td>114K</td>
    <td>100m</td>
    <td>340Mi</td>
    <td>1600m</td>
    <td>2200Mi</td>
    <td>20m</td>
    <td>22Mi</td>
    <td>30m</td>
    <td>520Mi</td>
  </tr>
</table>
