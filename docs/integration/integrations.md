This document describes the integration capabilities with various third party monitoring systems.

# Overview

Platform monitoring can be integrated with other monitoring systems which allows customers use unified interface for
application and infrastructure monitoring.

Platform monitoring integration scope can be described by the following points:

* `Metrics` - sending metrics collected by Prometheus to target monitoring system
* `Dashboards` - support Prometheus metrics visualization in target monitoring system
* `Alerting` - support alerting based on Prometheus metrics in target monitoring system
* `Autoscaling` - support horizontal or vertical pods autoscaling based on custom metrics

The table below describes platform monitoring integration capabilities with other popular monitoring systems.

| **Monitoring System**     | **Metrics** | **Dashboards** | **Alerting** | **Autoscaling** |
| ------------------------- | :---------: | :------------: | :----------: | :-------------: |
| Google Cloud Operations   |    ✓ Yes    |      ✗ No      |     ✗ No     |        ?        |
| Azure Monitor             |    ✓ Yes    |      ✗ No      |     ✗ No     |      ✗ No       |
| Amazon CloudWatch         |    ✓ Yes    |      ✗ No      |     ✗ No     |      ✗ No       |
| Amazon Managed Prometheus |    ✓ Yes    |       -        |    ✓ Yes     |        -        |
| Amazon Managed Grafana    |      -      |     ✓ Yes      |      -       |        -        |

**Legend:**

* `✓ Yes` - feature is supported and implemented
* `✗ No` - feature not supported by target monitoring system
* `?` - not analyzed yet
* `-` - not applicable

# Integrations

The sections below describe integration details for the certain monitoring system.

## Public Clouds

* [Amazon Cloud Watch](amazon-aws.md#aws-cloudwatch)
* [Amazon Managed Prometheus](amazon-aws.md#aws-managed-prometheus)
* [Azure Monitor](azure-monitor.md)
* [Google Cloud Operations](google-cloud.md)

## Vendor observability solutions

* [IBM Netcool](ibm-netcool.md)

## Enabling/disabling specific monitoring components

We have access not to all components in public clouds. For example in AWS/GKE/Azure we have no access to ETCD metrics.
So using ETCD dashboard, Service Monitor and Prometheus rules is useful in this case.

The monitoring-operator has a feature that allows enabling or disabling particular
dashboards, Prometheus rule groups and Service or Pod monitors which must be installed or skipped in the public cloud.

You can set the value of the `publicCloudName` parameter depending on which public cloud you are using.
The following values are currently available:

* `aws` - Amazon Web Services;
* `azure` - Microsoft Azure;
* `google` - Google Cloud Platform;
* `""` (default) - the monitoring-operator will manage dashboards, rules and monitors according to parameters
  `grafanaDashboards`, `prometheusRules` and `kubernetesMonitors` respectively.

Tables below describe which dashboards, rules and monitors will be installed or skipped by
the `publicCloudName` parameter.

**Legend:**

* `✓ Install` - component will be enabled by the parameter with this value
* `✗ Skip` - component will be disabled by the parameter with this value
* `-` - not affected

Grafana dashboards:

| **Dashboard**     | **`aws`** | **`azure`** | **`google`** |
|-------------------|:---------:|:-----------:|:------------:|
| Kubernetes / Etcd |  ✗ Skip   |   ✗ Skip    |    ✗ Skip    |
| CoreDNS           |  ✗ Skip   |   ✗ Skip    |    ✗ Skip    |

Groups of Prometheus rules:

| **Rule group** | **`aws`** | **`azure`** | **`google`** |
|----------------|:---------:|:-----------:|:------------:|
| Etcd           |  ✗ Skip   |   ✗ Skip    |    ✗ Skip    |
| CoreDnsAlerts  |  ✗ Skip   |   ✗ Skip    |    ✗ Skip    |

Service and Pod monitors:

| **Monitor**              | **`aws`** | **`azure`** | **`google`** |
|--------------------------|:---------:|:-----------:|:------------:|
| ETCD Service monitor     |  ✗ Skip   |   ✗ Skip    |    ✗ Skip    |
| Core DNS Service monitor |  ✗ Skip   |   ✗ Skip    |    ✗ Skip    |

Components that not presented in the tables are not affected by these settings.

Please notice that if you set the `publicCloudName` parameter then affected components will be skipped or
installed regardless of other parameters. If you want to customize lists of dashboards, rules or monitors, please, set
the `publicCloudName` to default empty value and use `grafanaDashboards`, `prometheusRules`
and `kubernetesMonitors` parameters instead.
