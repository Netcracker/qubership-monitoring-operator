This document describes the monitoring architecture.

The **Monitoring Operator** is an operator (Kubernetes native application) that watches and manages the following
components:

* `prometheus-operator`
* `grafana-operator`
* `kube-state-metrics`
* `node-exporter`

The main function of this application is to collect metrics through prometheus from services and pods, internal
Kubernetes metrics, and from deployed exporters, `kube-state-metrics` and `node-exporter`.

# Architecture

![Architecture](../images/prometheus_k8s_control_plane.png)

## Monitoring-operator

It is an operator which watches and manages (create, update, remove) such services into a current namespace as:

* [Monitoring-operator](#monitoring-operator)
* [Prometheus-operator](#prometheus-operator)
  * [Prometheus](#prometheus)
  * [AlertManager](#alertmanager)
* [Grafana-operator](#grafana-operator)
  * [Grafana](#grafana)
* [Kube-state-metrics](#kube-state-metrics)
* [Node-exporter](#node-exporter)
* [Custom Resources](#custom-resources)

Monitoring operator executes a monitoring deploy with specified parameters. It is also possible to determine which cloud
is running and create cloud specific resources. For example, during deploy into OpenShift, >3.11 creates
SecurityContextConstraints whereas for Kubernetes >1.15 creates Ingress and PodSecurityPolicy.

Also it manages the following custom resources:

* In-built ServiceMonitors CRs for collect internal Kubernetes metrics such as kubelet, kube-apiserver, etcd, so on.
* In-built PrometheusRule CRs with default generic alerts for Kubernetes and common alerts for services.
* In-built GrafanaDashboard CRs with default generic Grafana dashboards for show Kubernetes, Etcd, Go, JVM, metrics.

## Prometheus-operator

The Prometheus Operator provides Kubernetes native deployment and management of [Prometheus](#prometheus)
and related monitoring components. The purpose of this project is to simplify and automate the configuration of a
Prometheus based monitoring stack for Kubernetes clusters.

The Prometheus operator includes, but is not limited to, the following features:

* Kubernetes Custom Resources: Use Kubernetes custom resources to deploy and manage Prometheus, Alertmanager, and
  related components.
* Simplified Deployment Configuration: Configure the fundamentals of Prometheus like versions, persistence, retention
  policies, and replicas from a native Kubernetes resource.
* Prometheus Target Configuration: Automatically generate monitoring target configurations based on familiar Kubernetes
  label queries; no need to learn a Prometheus specific configuration language.

In the current namespace, Prometheus Operator executes managing of:

* [Prometheus](#prometheus)
* [AlertManager](#alertmanager)

For read more information, see the official documentation
at [https://github.com/prometheus-operator/prometheus-operator](https://github.com/prometheus-operator/prometheus-operator)

### Prometheus

Prometheus, a Cloud Native Computing Foundation project, is a systems and service monitoring system. It collects metrics
from configured targets at given intervals, evaluates rule expressions, displays the results, and can trigger alerts if
some condition is observed to be true.

Prometheus's main distinguishing features as compared to other monitoring systems are:

* A multi-dimensional data model (timeseries defined by metric name and set of key/value dimensions).
* PromQL, a powerful and flexible query language to leverage this dimensionality.
* No dependency on distributed storage; single server nodes are autonomous.
* Timeseries collection happens through a pull model over HTTP.
* Pushing timeseries is supported through an intermediary gateway.
* Targets are discovered through service discovery or static configuration.
* Multiple modes of graphing and dashboarding support.
* Support for hierarchical and horizontal federation.

For more information, refer to the official documentation at [https://prometheus.io](https://prometheus.io).

### AlertManager

The Alertmanager handles alerts sent by client applications such as the Prometheus server. It takes care of
deduplicating, grouping, and routing them to the correct receiver integration such as email, PagerDuty, or OpsGenie. It
also takes care of silencing and inhibition of alerts.

For more information, refer to the official documentation at
[https://prometheus.io/docs/alerting/latest/alertmanager](https://prometheus.io/docs/alerting/latest/alertmanager)

## Grafana-operator

A Kubernetes Operator based on the Operator SDK for creating and managing Grafana instances.

In the current namespace, the Grafana Operator executes the managing of:

* [Grafana](#grafana)

For more information, refer to the official documentation at
[https://github.com/grafana/grafana-operator](https://github.com/grafana/grafana-operator)

### Grafana

Grafana is an open source visualization and analytics software. It allows you to query, visualize, alert on, and explore
your metrics no matter where they are stored. In plain language, it provides you with tools to turn your time-series
database (TSDB) data into beautiful graphs and visualizations.

For more information, refer to the official documentation at
[https://grafana.com/docs/grafana/latest/fundamentals/](https://grafana.com/docs/grafana/latest/fundamentals/)

## Kube-state-metrics

It is a simple service that listens to the Kubernetes API server and generates metrics about the state of the objects.
It is not focused on the health of the individual Kubernetes components, but rather on the health of the various objects
inside, such as deployments, nodes, and pods.

kube-state-metrics is about generating metrics from Kubernetes API objects without modification. This ensures that the
features provided by kube-state-metrics have the same grade of stability as the Kubernetes API objects themselves. In
turn, this specifies that kube-state-metrics in certain situations may not show the exact same values as kubectl, as
kubectl applies certain heuristics to display comprehensible messages. kube-state-metrics exposes raw data unmodified
from the Kubernetes API. This way users have all the data they require and perform heuristics as they see fit.

For more information, refer to the official documentation at
[https://github.com/kubernetes/kube-state-metrics](https://github.com/kubernetes/kube-state-metrics)

## Node-exporter

It is a simple service that collects hardware and OS metrics exposed by *NIX kernels.

The Monitoring Operator provides unify endpoint to specify the desired state of the Monitoring application through
custom k8s resource PlatformMonitoring.

For more information, refer to the official documentation at
[https://github.com/prometheus/node_exporter](https://github.com/prometheus/node_exporter)

# Custom Resources

Monitoring operator has one custom resource which contains all deploy settings and by which the operator manages the
deployment:

* `PlatformMonitoring` - Allows to declaratively specify the monitoring deployment settings. The Monitoring Operator
  automatically generates CRs `AlertManager`, `Prometheus`, `Grafana` to manage the deployment of these services.

Operators into monitoring deploy introduces several custom resource definitions which services can added itself:

* `ServiceMonitor` - Allows to declaratively specify how groups of Kubernetes services should be monitored. The
  Prometheus Operator automatically generates Prometheus scrape configuration based on the current state of the objects
  in the API server.
* `PodMonitor` - Allows to declaratively specify how group of pods should be monitored. The Prometheus Operator
  automatically generates Prometheus scrape configuration based on the current state of the objects in the API server.
* `PrometheusRule` - Allows to define a desired set of Prometheus alerting and/or recording rules. The Prometheus
  Operator generates a rule file, which can be used by Prometheus instances.
* `GrafanaDashboard` - Allows to specify the Grafana dashboard as a json config which should be added into Grafana. The
  Grafana Operator automatically creates a dashboard into Grafana.

All the above CRs have a namespaced scope and should have unique names only into one namespace.

Operators automatically discover all the above CRs on all available namespaces and filter them by the following
specified labels:

* By default, the custom resources `ServiceMonitor`, `PodMonitor`, `PrometheusRule` filter by
  label `app.kubernetes.io/component=monitoring`.
* By default, the custom resource `GrafanaDashboard` filter by label `app=grafana`.

This means that each CR should have such a label to be discovered by an operator.

Also, there are some custom resources which setup the service state and settings, but currently Monitoring Operator does
not allow to change them (it means that if you manually changed a CR, the operator overrides it by the original
version):

* `AlertManager` - Allows to define a desired Prometheus deployment. The Prometheus Operator ensures at all times that a
  deployment matching the resource definition is running.
* `Prometheus` - Allows to define a desired Alertmanager deployment. The Prometheus Operator ensures at all times that a
  deployment matching the resource definition is running.
* `Grafana` - Allows to define a desired Grafana deployment. The Grafana Operator ensures at all times that a deployment
  matching the resource definition is running.
