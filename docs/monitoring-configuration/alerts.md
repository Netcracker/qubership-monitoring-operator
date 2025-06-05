This document describes how to configure Prometheus alerts in Monitoring.

# Platform Monitoring alerts

Out-of-box alerts delivered along with Monitoring.

## Enabling and disabling alerts

OOB alerts installation can be customised by changing
[`prometheusRules` section](../installation/README.md#prometheus-rules).

You can enable or disable installation of all alerts by changing `prometheusRules.install` option. The default value is
`true` which means that OOB alerts will be enabled in general by default.

Alerts divided into groups. Alerts are grouped by meaning for monitoring purposes: SelfMonitoring, NodeExporters, Etcd,
etc. Each alert group can be enabled or disabled individually by changing `prometheusRules.ruleGroups` parameter.
This parameter contains a list of alert groups that should be installed.

OOB alerts include the groups mentioned in the table below, but **not all alert groups are enabled by default**,
therefore you should pay attention if you want to use alerts disabled by default (e.g. Heartbeat alert):

| Alert group        | Enabled by default |
|--------------------|--------------------|
| Heartbeat          | ✗ No               |
| SelfMonitoring     | ✓ Yes              |
| AlertManager       | ✓ Yes              |
| KubebernetesAlerts | ✓ Yes              |
| NodeProcesses      | ✓ Yes              |
| NodeExporters      | ✓ Yes              |
| DockerContainers   | ✓ Yes              |
| HAmode             | ✗ No               |
| HAproxy            | ✗ No               |
| Etcd               | ✓ Yes              |
| NginxIngressAlerts | ✓ Yes              |
| CoreDnsAlerts      | ✓ Yes              |
| DRAlerts           | ✓ Yes              |
| BackupAlerts       | ✓ Yes              |

Full list of all alerts can be found in the [alerts-oob document](../defaults/alerts.md).

If you want to enable alerts for HAmode or HAproxy,
you should add `HAmode` or `HAproxy` respectively to `prometheusRules.ruleGroups` parameter.

If you want to enable Dead Man's Switch (Heartbeat) alert, you should add `Heartbeat` to `prometheusRules.ruleGroups`
parameter.

You can find examples of configuration [in the appropriate section](#examples).

### Dead Man's Switch alert

[Dead Man's Switch](https://en.wikipedia.org/wiki/Dead_man%27s_switch) alert is a special always-firing alert that meant
to ensure that the entire alerting pipeline is functional. If this alert stops firing, it means that some of
the critical monitoring and/or alerting components have failed.

Platform Monitoring's Dead Man's Switch alert is placed in the `Heartbeat` alert group and called `DeadMansSwitch`.
It uses the simplest expression possible under the hood: `vector(1)`, and the lowest severity: `information`.
Even the simplest expressions are calculated on the monitoring back-end (Prometheus/VMSingle),
so the alert checks that side of the monitoring. The alert has `for: 0s`, so it should start fire immediately since
all base monitoring and alerting components are installed.

Fields `for`, `expr` and `severity` can be [overridden](#alerts-overriding) as well as every other OOB alert.

This type of alert is **disabled by default**. You can enable the Dead Man's Switch alert by adding `Heartbeat`
alert group to the `prometheusRules.ruleGroups` parameter. Example of configuration with enabled Dead Man's Switch can
be found [here](#configuration-with-enabled-additional-groups).

**Attention:** If you want to enable Dead Man's Switch alert and have a connected notification system for alerts,
proactively make sure that your system is ready for a constantly firing alert and will not encounter
a flood of notifications because of this. This can be done, for example, by ignoring alerts with a severity
below warning.

### HAmode alerts

`HAmode` is an alert group includes alert rules that report that some Deployments and StatefulSets do not comply
with working conditions in HA (High Availability) mode. Services in HA mode should have 2 or more available or at least
desired replicas, and these replicas should be placed on different nodes.

`HAmode` will be triggered if:

* Some Deployments or StatefulSets have less than 2 desired replicas
* Some Deployments or StatefulSets have less than 2 available replicas
* Some Deployments or StatefulSets have 2 or more replicas placed on the same node

This alert group is **disabled by default** in case if clusters have a lot of services that don't follow conditions
for the HA mode.

Platform Monitoring has a dashboard called [`HA services`](../defaults/dashboards/ha-services.md) for the same
purposes as the `HAmode` alert group.

## Alerts overriding

Monitoring has a mechanism that allows changing specific alert(-s) from the OOB set. If you want to override an alert,
you can use `prometheusRules.override` parameter. This parameter includes list of objects with definitions for alerts
that should be overridden.

You can override fields `for`, `expr` and `severity`. The `prometheusRules.override` parameter looks like this:

```yaml
prometheusRules:
  override:
    - group: SelfMonitoring
      alert: PrometheusNotificationsBacklog
      for: 0s
      expr: "min_over_time(prometheus_notifications_queue_length[20m]) > 0"
      severity: high
    - ...
```

You can either override all 3 fields for the alert or only one of them. Every item in the `prometheusRules.override`
parameter may have the following fields:

<!-- markdownlint-disable line-length -->
| Parameter | Description                                                                                                                                      | Required |
|-----------|--------------------------------------------------------------------------------------------------------------------------------------------------|----------|
| group     | Name of alert group where the overridden alert from.                                                                                             | true     |
| alert     | Name of the overridden alert.                                                                                                                    | true     |
| for       | Alerts are considered firing once they have been returned for this long. Alerts which have not yet fired for long enough are considered pending. | false    |
| expr      | How long an alert will continue firing after the condition that triggered it has cleared.                                                        | false    |
| severity  | Shows the level of importance for the alert. Recommended levels: critical, high, warning, information.                                           | false    |
<!-- markdownlint-enable line-length -->

You can find more information about alert configuring process in the
[alert best practice document](../user-guides/alert-best-practice.md).

# Examples

## Default configuration example

```yaml
prometheusRules:
  install: true
  ruleGroups:
    - SelfMonitoring
    - AlertManager
    - KubebernetesAlerts
    - NodeProcesses
    - NodeExporters
    - DockerContainers
    - Etcd
    - NginxIngressAlerts
    - CoreDnsAlerts
    - DRAlerts
    - BackupAlerts
```

## Configuration with enabled additional groups

```yaml
prometheusRules:
  install: true
  ruleGroups:
    - Heartbeat
    - HAmode
    - HAproxy
    - SelfMonitoring
    - AlertManager
    - KubebernetesAlerts
    - NodeProcesses
    - NodeExporters
    - DockerContainers
    - Etcd
    - NginxIngressAlerts
    - CoreDnsAlerts
    - DRAlerts
    - BackupAlerts
```

## Configuration with overridden alert

The following configuration changes `for` to `0s`, `expr`
to `min_over_time(prometheus_notifications_queue_length[20m]) > 0` and `severity` to `high` for
alert `PrometheusNotificationsBacklog` name from `SelfMonitoring` group.

```yaml
prometheusRules:
  install: true
  ruleGroups:
    - SelfMonitoring
    - AlertManager
    - KubebernetesAlerts
    - NodeProcesses
    - NodeExporters
    - DockerContainers
    - Etcd
    - NginxIngressAlerts
    - CoreDnsAlerts
    - DRAlerts
    - BackupAlerts
  override:
    - group: SelfMonitoring
      alert: PrometheusNotificationsBacklog
      for: 0s
      expr: "min_over_time(prometheus_notifications_queue_length[20m]) > 0"
      severity: high
```

