This document provides information about best practices for working
with components and resources that can be used in the monitoring-operator.

# Grafana dashboards

Best practices for Grafana dashboards.

## Dashboard UID

The unique identifier (UID) of a dashboard can be used for uniquely identify a dashboard between multiple Grafana
installs. It’s automatically generated if not provided when creating a dashboard. The uid allows having consistent URLs
for accessing dashboards and when syncing dashboards between multiple Grafana installs, see dashboard provisioning for
more information. This means that changing the title of a dashboard will not break any bookmarked links to that
dashboard.

**ATTENTION:** The uid can have a maximum length of 40 characters.

If you create your own dashboard, we highly recommend to use human-readable UIDs in the following format:
`<dashboard-namespace>-<short-dashboard-name>`. It allows using all advantages of consistent URLs for dashboards, e.g.
to create stable drill down links. Namespace should be presented in the UID to avoid situations when two or more
dashboards with the same names have equal UIDs.

UIDs for OOB monitoring dashboards (managed by monitoring-operator):

| Dashboard name                       | UID (namespace + short name)              |
|--------------------------------------|-------------------------------------------|
| alertmanager-overview                | `<namespace>-alertmanager-overview`       |
| alerts-overview                      | `<namespace>-alerts-overview`             |
| core-dns-dashboard                   | `<namespace>-core-dns`                    |
| etcd-dashboard                       | `<namespace>-etcd`                        |
| govm-processes                       | `<namespace>-govm-processes`              |
| grafana-overview                     | `<namespace>-grafana-overview`            |
| ha-services                          | `<namespace>-ha-services`                 |
| home-dashboard                       | `<namespace>-home-dashboard`              |
| ingress-list-of-ingresses            | `<namespace>-ing-list-of-ingresses`       |
| ingress-nginx-controller             | `<namespace>-ing-nginx-controller`        |
| ingress-request-handling-performance | `<namespace>-ing-req-handl-perform`       |
| jvm-processes                        | `<namespace>-jvm-processes`               |
| kubernetes-apiserver                 | `<namespace>-k8s-apiserver`               |
| kubernetes-cluster-overview          | `<namespace>-k8s-cluster-overview`        |
| kubernetes-distribution-by-labels    | `<namespace>-k8s-distr-by-labels`         |
| kubernetes-kubelet                   | `<namespace>-k8s-kubelet`                 |
| kubernetes-namespace-resources       | `<namespace>-k8s-namespace-resources`     |
| kubernetes-nodes-resources           | `<namespace>-k8s-nodes-resources`         |
| kubernetes-pod-resources             | `<namespace>-k8s-pod-resources`           |
| kubernetes-pods-distribution-by-node | `<namespace>-k8s-pods-distr-by-node`      |
| kubernetes-pods-distribution-by-zone | `<namespace>-k8s-pods-distr-by-zone`      |
| kubernetes-top-resources             | `<namespace>-k8s-top-resources`           |
| node-details                         | `<namespace>-node-details`                |
| openshift-apiserver                  | `<namespace>-os-apiserver`                |
| openshift-cluster-version-operator   | `<namespace>-os-cluster-version-operator` |
| openshift-state-metrics              | `<namespace>-os-state-metrics`            |
| openshift-haproxy                    | `<namespace>-os-haproxy`                  |
| operators-overview                   | `<namespace>-operators-overview`          |
| overall-platform-health              | `<namespace>-overall-platform-health`     |
| prometheus-cardinality-explorer      | `<namespace>-prom-cardinality`            |
| prometheus-self-monitoring           | `<namespace>-prom-self-monitoring`        |
| tls-status                           | `<namespace>-tls-status`                  |
| victoriametrics-vmagent              | `<namespace>-vm-vmagent`                  |
| victoriametrics-vmalert              | `<namespace>-vm-vmalert`                  |
| victoriametrics-vmoperator           | `<namespace>-vm-vmoperator`               |
| victoriametrics-vmsingle             | `<namespace>-vm-vmsingle`                 |

There are several dashboards managed by Helm. In some cases name of the file with the dashboard is not the same as
the title of the dashboard, so we'll use dashboard titles below. Also, remember that the `home-dashboard` is present
in both places: in the grafana-operator Helm chart as a ConfigMap and together with the rest of dashboards managed by
the operator:

| Dashboard title                      | Helm subchart            | UID (namespace + short name)           |
|--------------------------------------|--------------------------|----------------------------------------|
| Blackbox Probes                      | blackbox-exporter        | `<namespace>-blackbox-probes`          |
| SSL/TLS Certificates                 | cert-exporter            | `<namespace>-ssl-tls-certs`            |
| Kafka Java Clients Monitoring        | common-dashboards        | `<namespace>-kafka-java-clients`       |
| Configurations Streamer              | configurations-streamer  | `<namespace>-configurations-streamer`  |
| Backup Daemon                        | grafana-operator         | `<namespace>-backup-daemon`            |
| Home Dashboard                       | grafana-operator         | `<namespace>-home-dashboard`           |
| Prometheus / Graphite remote adapter | graphite-remote-adapter  | `<namespace>-graphite-remote-adapter`  |
| Network Latency Details              | network-latency-exporter | `<namespace>-network-latency-details`  |
| Network Latency Overview             | network-latency-exporter | `<namespace>-network-latency-overview` |
| DR Overview                          | promxy                   | `<namespace>-dr-overview`              |
| Version overview                     | version-exporter         | `<namespace>-version-overview`         |

If the name of the namespace is too long, the whole UID of the OOB dashboard will be cut to 40 symbols.

## Creating custom dashboard

Best practices for creating custom Grafana dashboards.

### Tags

If you create dashboard, you should add some tags that will be described, what you can see on this dashboard.

**NOTE**: All tags must be in lowercase. If the tag contains more than one word, words must be in
"kebab-case" (separated with hyphens).

The tags below should be added to dashboard, if it satisfies the following conditions:

* tag `k8s` - if the dashboard shows data about Kubernetes cluster;
* tag `prometheus` - if the dashboard shows information about services (e.g. kafka, postgresql, mongodb, etc.);
* tag `standalone` - if the dashboard shows information about standalone hosts (e.g. Graylog, balancers, etc.);
* tag `self-monitor` - if the dashboard shows information about the monitoring system.

Also, dashboard should contain tags that describe specific information that can be founded on it.
For example, dashboard that shows information about Kubernetes namespace resources can contain
tags `k8s` and `k8s-namespaces`, dashboard that shows information about PostgreSQL - `prometheus` and `postgres`.

# Recommendations for creating recording-rules
<!-- markdownlint-disable line-length -->
* **Do not create** recording rules without sense, recording rules must have a reason, for example:
  * The new aggregated metric with heavy calculation will be used in alert - **good case**, we will use already calculated value
  * The new aggregated metric will be used on the Grafana dashboard and will open very rare - **bad case**, better to calculate the value in runtime that spends CPU time on its calculation
* **Do not write** recording rules that can be used to calculate a big metrics scope, for example, to calculate CPU usage for the last 15 minutes for all pods in the Cloud
  * If a product or project wants to use such rules I offer to add them only once and from our side, as a part of Monitoring deployment. A good example is moving a rule to calculate CPU usage for the last 15 minutes to Monitoring deployment.
* **Do not duplicate** recording rules that should calculate metrics by the same scope
* Recording rules **must be used** only to calculate aggregations or to **prepare new metrics** aggregated from some metrics, for example:
  * Right case of usage - calculate CPU usage for the last 5-10-15 minutes
  * Right case of usage - calculate new metric that will include metrics values and labels from (kube_pod_labels)
  * Wrong case of usage - calculate any metric just because to calculate metrics need to use a big query
<!-- markdownlint-enable line-length -->