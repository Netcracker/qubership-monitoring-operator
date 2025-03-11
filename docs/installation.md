This document provides information about the requirements, configuration, and steps to install Monitoring to an
environment.

# Table of Content

* [Table of Content](#table-of-content)
* [Supported Versions](#supported-versions)
  * [Oldest Version and Deploy in Old Kubernetes/OpenShift Versions](#oldest-version-and-deploy-in-old-kubernetesopenshift-versions)
* [Versioning](#versioning)
  * [Components](#components)
  * [Kubernetes API](#kubernetes-api)
* [Prerequisites](#prerequisites)
  * [Before You Begin](#before-you-begin)
    * [AWS](#aws)
      * [AWS EBS persistence volume](#aws-ebs-persistence-volume)
      * [AWS EFS persistence volume](#aws-efs-persistence-volume)
  * [Permissions](#permissions)
  * [Hardware Requirements](#hardware-requirements)
    * [Overall Minimal Requirements](#overall-minimal-requirements)
    * [Detailed Hardware Requirements](#detailed-hardware-requirements)
* [Inventory](#inventory)
  * [Basic Components](#basic-components)
    * [Default deploy settings](#default-deploy-settings)
    * [Full Stack Deploy](#full-stack-deploy)
    * [Enabling and disabling dashboards, rules and monitors](#enabling-and-disabling-dashboards-rules-and-monitors)
    * [Deploy with cert-exporter](#deploy-with-cert-exporter)
      * [Types of Checked Certificates](#types-of-checked-certificates)
        * [Certificates from Files](#certificates-from-files)
        * [Certificates from Kubeconfig](#certificates-from-kubeconfig)
        * [Certificates from Kubernetes Secrets](#certificates-from-kubernetes-secrets)
      * [Installation with Root Rights](#installation-with-root-rights)
  * [Integrations](#integrations)
    * [Metrics Collectors and Long Term Storages](#metrics-collectors-and-long-term-storages)
      * [Common Restrictions and Recommendations](#common-restrictions-and-recommendations)
        * [Filter Pushed Metrics](#filter-pushed-metrics)
        * [RemoteWrite Increases Memory Usage](#remotewrite-increases-memory-usage)
        * [Add External Labels](#add-external-labels)
        * [Use Authentication for RemoteWrite and RemoteRead if Remote Storage Supports it](#use-authentication-for-remotewrite-and-remoteread-if-remote-storage-supports-it)
      * [Graphite](#graphite)
      * [InfluxDB](#influxdb)
      * [VictoriaMetrics](#victoriametrics)
      * [Prometheus Federation](#prometheus-federation)
    * [Alerting Systems](#alerting-systems)
      * [External AlertManager](#external-alertmanager)
        * [AlertManager in the same Cloud](#alertmanager-in-the-same-cloud)
        * [AlertManager outside the Cloud](#alertmanager-outside-the-cloud)
  * [Storage](#storage)
    * [Known issues](#known-issues)
    * [Storage with Dynamic Provisioning](#storage-with-dynamic-provisioning)
    * [Manual Storage Provisioning](#manual-storage-provisioning)
    * [Disabling Default StorageClasses](#disabling-default-storageclasses)
  * [Configuration](#configuration)
    * [root parameters](#root-parameters)
    * [alertmanager](#alertmanager)
    * [auth](#auth)
    * [blackbox-exporter](#blackbox-exporter)
    * [cert-exporter](#cert-exporter)
    * [cloud-events-exporter](#cloud-events-exporter)
    * [cloudwatch-exporter](#cloudwatch-exporter)
    * [grafana](#grafana)
      * [grafana-operator](#grafana-operator)
      * [grafana-image-renderer](#grafana-image-renderer)
    * [grafana-dashboards](#grafana-dashboards)
    * [graphite-remote-adapter](#graphite-remote-adapter)
    * [ingress](#ingress)
    * [integration](#integration)
      * [jaeger](#jaeger)
      * [clickHouse](#clickhouse)
    * [json-exporter](#json-exporter)
    * [kube-state-metrics](#kube-state-metrics)
    * [kubernetes-monitors](#kubernetes-monitors)
    * [monitor](#monitor)
    * [monitoring-operator](#monitoring-operator)
    * [network-latency-exporter](#network-latency-exporter)
    * [node-exporter](#node-exporter)
    * [oauth-proxy](#oauth-proxy)
    * [prometheus](#prometheus)
      * [prometheus-operator](#prometheus-operator)
    * [prometheus-adapter](#prometheus-adapter)
      * [prometheus-adapter-operator](#prometheus-adapter-operator)
      * [prometheus-adapter-operator-tls-config](#prometheus-adapter-operator-tls-config)
    * [prometheus-rules](#prometheus-rules)
    * [promitor-agent-scraper](#promitor-agent-scraper)
    * [promitor-agent-resource-discovery](#promitor-agent-resource-discovery)
    * [publicCloudName](#publiccloudname)
    * [promxy](#promxy)
    * [pushgateway](#pushgateway)
    * [stackdriver](#stackdriver)
    * [stackdriver-exporter](#stackdriver-exporter)
    * [victoria-metrics](#victoria-metrics)
      * [vmoperator](#vmoperator)
      * [vmsingle](#vmsingle)
      * [vmagent](#vmagent)
      * [vmalertmanager](#vmalertmanager)
      * [vmalert](#vmalert)
      * [vmauth](#vmauth)
      * [vmuser](#vmuser)
      * [vmcluster](#vmcluster)
      * [cleanup](#cleanup)
      * [TLS Config](#tls-config)
    * [version-exporter](#version-exporter)
      * [install parameters](#install-parameters)
      * [arguments](#arguments)
      * [postgres collector parameters](#postgres-collector-parameters)
      * [http requests collector parameters](#http-requests-collector-parameters)
      * [ConfigMap collector parameters](#configmap-collector-parameters)
    * [SSH collector parameters](#ssh-collector-parameters)
    * [Integration tests](#integration-tests)
* [Deploy](#deploy)
  * [Admission Migration](#admission-migration)
  * [CRD Upgrade](#crd-upgrade)
  * [Cleanup VictoriaMetrics Objects](#cleanup-victoriametrics-objects)
  * [Upgrade](#upgrade)
    * [Upgrade from 0.50.0 version](#upgrade-from-0500-version)
  * [Manual Deploy using Helm](#manual-deploy-using-helm)
    * [Ingress using Helm](#ingress-using-helm)
    * [Installing the Chart](#installing-the-chart)
    * [Upgrading the Chart](#upgrading-the-chart)
    * [Uninstalling the Chart](#uninstalling-the-chart)
* [Post Deploy Checks](#post-deploy-checks)
  * [Common Post Deploy Checks](#common-post-deploy-checks)
  * [Jobs Post Deploy Check](#jobs-post-deploy-check)
* [What's next](#whats-next)


# Supported Versions

According to platform third party support policy now we are supporting deploy in Kubernetes N +- 2.

Current recommended Kubernetes version is `1.26.x`, so we support:

| Kubernetes version     | Status of `0.48.0`    |
| ---------------------- | --------------------- |
| `1.24.x`               | Tested                |
| `1.25.x`               | Tested                |
| `1.26.x` (recommended) | Tested                |
| `1.27.x`               | Forward compatibility |
| `1.28.x`               | Forward compatibility |

**Note:** `Forward compatibility` means that the latest version don't use any Kubernetes APIs that will be remove
in next Kubernetes versions. All Kubernetes APIs verified by official documentation
[Deprecated API Migration Guide](https://kubernetes.io/docs/reference/using-api/deprecation-guide/)

**Note:** OpenShift 4.x based on Kubernetes and regularly picks up new Kubernetes releases. So compatibility with
OpenShift you can track by Kubernetes version. To understand which Kubernetes version OpenShift use you can read
it release notes. For example, `OpenShift 4.12` based on `Kubernetes v1.25.0` about which wrote in
[Release notes](https://docs.openshift.com/container-platform/4.12/release_notes/ocp-4-12-release-notes.html).

If we will consider support by types of Clouds we can create the following table:

| Clouds                                   | `0.48.0` |
| ---------------------------------------- | -------- |
| AWS Elastic Kubernetes Service (AWS EKS) | ✓        |
| Azure Kubernetes                         | ✓        |
| Google Kubernetes Engine (GKE)           | ✓        |
| On-premise Kubernetes >= 1.25            | ✓        |
| On-premise OpenShift >= 4.10             | ✓        |

[Back to TOC](#table-of-content)

## Oldest Version and Deploy in Old Kubernetes/OpenShift Versions

This section contains information about old versions and which Kubernetes or OpenShift them supported.

Last Monitoring version, that supported specified Kubernetes versions:

| Kubernetes  | Monitoring |
| ----------- | ---------- |
| `< v1.25.0` | `0.46.0`   |
| `< v1.18.0` | `0.26.0`   |

Last Monitoring version, that supported specified OpenShift versions:

| OpenShift | Monitoring        |
| --------- | ----------------- |
| `v3.11.0` | `0.46.0`          |
| `1.x`     | `Never supported` |

[Back to TOC](#table-of-content)

# Versioning

This section describes how Monitoring components and API versioning.

## Components

For versioning all components we are using `Semantic Versioning` (SemVer).

It means that a version number `MAJOR.MINOR.PATCH`, increment the:

* `MAJOR` version when you make incompatible API changes,
* `MINOR` version when you add functionality in a backwards compatible manner, and
* `PATCH` version when you make backwards compatible bug fixes.

More information about SemVer you can find by link [https://semver.org](https://semver.org/).

[Back to TOC](#table-of-content)

## Kubernetes API

Except of components we also provide Kubernetes API, such as Custom Resource Definitions (CRDs).
Kubernetes API (and CRDs in particular) usually versioning as described in
[The Kubernetes API](https://kubernetes.io/docs/concepts/overview/kubernetes-api/#api-groups-and-versioning).

**Note**: If you want to deploy the `monitoring-operator` into **Kubernetes v1.15** or lower or
**OpenShift v3.11** or lower, you must work with _v1beta1_ CRDs manually. For more information see
[Maintenance Guide: Work with legacy CRDs](maintenance.md#work-with-legacy-crds).

For example different Kubernetes API can cycle from alpha to release with versions:

* `/<api>/v1alpha1`
* `/<api>/v1beta1`
* `/<api>/v1`

**Note**: The versions may be different, not only those that were listed above.

One important point of the deployment is associated with versioning the CRDs.

Helm doesn't have an ability to update/remove CRDs, but can create them. It means that CRDs will
create during first deploy. But during any updates and remove deployment Helm will not update or remove CRDs.

For more details you can read [Issue: CRD update during helm upgrade](https://github.com/helm/helm/issues/8668) and
[Design: CRD Handling in Helm](https://github.com/helm/community/blob/f9e06c16d89ccea1bea77c01a6a96ae3b309f823/architecture/crds.md).

So during the update, you need to manually update all CRDs. For more information, refer to
[Maintenance Guide: Update Procedure](maintenance.md#update-procedure).

[Back to TOC](#table-of-content)

# Prerequisites

This section describes which prerequisites should be checked before deploy Monitoring.

## Before You Begin

* Verify that you have cloud with Kubernetes 1.19+ or OpenShift 3.11+ versions.
* Verify that you use CLI `kubectl` 1.19+ or `oc` 3.11+.
* Verify that you use Helm 3.0+.
* Verify that namespace is already created before the installation.
* Verify that you have enough rights to deploy in the target project: [Permissions](#permissions)
* Verify the hardware requirements: [Hardware Sizing](../internal/install/hardware-sizing.md)
* Verify the installation parameters. The list of parameters is available in [Inventory](#inventory).

Also it is important to note the following:

* If you disable automation [CRD upgrade](#crd-upgrade) you need to apply ALL [CRDs](./crds) before installation.
* If you want to update your deployment firstly read [Maintenance Guide: Update Procedure](maintenance.md#update-procedure).
* If you have Pod Security Policy Restricted by default on namespace, it can lead to unreachable `node-exporter` pods.
  Then you need to add privileged PSP ClusterRole to ServiceAccount `<namespace>-node-exporter`.
* If another prometheus/victoriametrics was installed before, for example, default prometheus/victoriametrics
  in `monitoring` namespace, which is included in Kubernetes, you need to check its `namespaces` scope.
  Add to another prometheus/victoriametrics-operator Deployment restriction on namespaces  `- '--namespaces=monitoring'`.

  ```yaml
  containers:
    - args:
      - '--kubelet-service=kube-system/kubelet'
      - '--namespaces=monitoring'
        ...
  ```

* If you deploy on OpenShift, verify that the namespace has empty `node-selector`.

  ```bash
  $ oc describe namespace <namespace>
  Name:         <namespace>
  Labels:       <none>
  Annotations:  openshift.io/description=
                openshift.io/display-name=
                openshift.io/node-selector=
                openshift.io/requester=admin
                openshift.io/sa.scc.mcs=s0:c16,c15
                openshift.io/sa.scc.supplemental-groups=1000270000/10000
                openshift.io/sa.scc.uid-range=1000270000/10000
  Status:        Active
  ```

[Back to TOC](#table-of-content)

### AWS

In case of using Persistence Volume in AWS for VictoriaMetrics or Prometheus we **recommended** use `AWS EBS`.

#### AWS EBS persistence volume

Selecting between `AWS EBS` disk types you should select at least HDD with type `Throughput Optimized HDD volumes`.

For Prometheus we recommended use SSD volumes.
For VictoriaMetrics you can use at least HDD with type `Throughput Optimized HDD volumes`. If you want you also cause
SSD volumes.

About `AWS EBS` disk types you can read in the documentation
[https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ebs-volume-types.html](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ebs-volume-types.html).

#### AWS EFS persistence volume

**Warning!** Prometheus **doesn't support** non-POSIX compliant storages and **we strongly doesn't recommended**
try to install on `AWS EFS`. Prometheus work very bad with all NFS-like storages and especially with `AWS EFS`.
Information about NFS in the official documentation
[https://prometheus.io/docs/prometheus/latest/storage/#operational-aspects](https://prometheus.io/docs/prometheus/latest/storage/#operational-aspects).

VictoriaMetrics officially support all volume types, see a documentation
[https://docs.victoriametrics.com/#prominent-features](https://docs.victoriametrics.com/#prominent-features).

But if you want use `AWS EFS` you must plan you volume by throughput. You shouldn't use volume with `Bursting`
throughput type. Need to select `Provisioned Throughput` mode for EFS volume or at least `Elastic Throughput`.
The first option will charge some additional money. Minimal value of `Provisioned Throughput` is 10 Mb/s
(of course better to use much more).

## Permissions

The `monitoring-operator` service account requires ClusterRole and permissions for a lot of Cluster APIs because the
following roles work through ClusterRole, and the operator should have the ability to create ClusterRole and
ClusterRoleBinding resources.

Components for which ClusterRoles are created by default:

* `prometheus`
* `prometheus-operator`
* `grafana-operator`
* `kube-state-metrics`
* `node-exporter`
* `cert-exporter`
* `cloudwatch-exporter`
* `configuration-streamer`
* `network-latency-exporter`
* `prometheus-adapter-operator`
* `version-exporter`
* `vmoperator`

Thus a ClusterRole of `monitoring-operator` should include a set of all the required permissions from all ClusterRoles
that `monitoring-operator` creates during the deploy.

Also the `monitoring-operator` supports to deploy with restricted privileges when the operator has access only to
resources inside the namespace and restricted access to cluster scoped resources. In this case, all the necessary
privileges for cluster scoped resources and resources in another namespaces must be granted manually before the deploy.

[Back to TOC](#table-of-content)

## Hardware Requirements

### Overall Minimal Requirements

The overall hardware requirements can be calculated as sum of requirements for the platform monitoring components:

* `alertmanager` resources
* `grafana` resources
* `grafana-operator` resources
* `kube-state-metrics` resources
* `monitoring-operator` resources
* `node-exporter` resources multiplied by number of nodes in cluster
* `victoriametris` resources calculated for specific system load
* `vm-operator` resources

The **minimal** recommended hardware requirements for average system is:

| Resource | Requirement |
| -------- | ----------- |
| CPU      | 3000m       |
| RAM      | 3000Mi      |
| IOPS     | -           |
| Storage  | 10GB        |

[Back to TOC](#table-of-content)

### Detailed Hardware Requirements

Monitoring can be deployed in various schemes and with additional components for integrations. Also it can collect
various metrics count, and different components can use a different number of resources.

There is a separate guide [Hardware Sizing](../internal/install/hardware-sizing.md) if you need to calculate
the exact amount of resources for each component.

[Back to TOC](#table-of-content)

# Inventory

This section describes how to create an inventory and which parameters are used to deploy Monitoring stack in one form
or another.

## Basic Components

Monitoring includes a lot of components which can execute various functions. But not all projects need all functions. So
before installing Monitoring into the inventory, the components that should be installed can be specified.

There is only one component that is always installed: `monitoring-operator`

The default configuration is:

* `monitoring-operator`
* `alertmanager`
* `grafana`
* `grafana-operator`
* `kube-state-metrics`
* `node-exporter`
* `victoriametrics`
* `vm-operator`

The available components for enabling/disabling are:

* [`alertmanager`](#alertmanager)
* Prometheus:
  * [`prometheus-operator`](#prometheus-operator)
  * [`prometheus`](#prometheus)
  * [`prometheus rules`](#prometheus-rules)
  * [`kubernetes monitors`](#kubernetes-monitors)
* Grafana:
  * [`grafana-operator`](#grafana-operator)
  * [`grafana`](#grafana)
  * [`grafana dashboards`](#grafana-dashboards)
* [`pushgateway`](#pushgateway)
* Exporters:
  * [`node-exporter`](#node-exporter)
  * [`kube-state-metrics`](#kube-state-metrics)
  * [`blackbox-exporter`](#blackbox-exporter)
  * [`cert-exporter`](#cert-exporter)
  * [`version-exporter`](#version-exporter)
  * [`network-latency-exporter`](#network-latency-exporter)
  * [`cloudwatch-exporter`](#cloudwatch-exporter) (for AWS public cloud)
  * [`promitor`](#promitor-agent-scraper) (for Azure public cloud)
  * [`stackdriver-exporter`](#stackdriver-exporter) (for Google public cloud)
* [`promxy`](#promxy)
* Prometheus-adapter:
  * [`prometheus-adapter-operator`](#prometheus-adapter-operator)
  * [`prometheus-adapter`](#prometheus-adapter)
* [`graphite-remote-adapter`](#graphite-remote-adapter)
* [`configuration-streamer`](#configurations-streamer)
* VictoriaMetrics stack:
  * [`vmoperator`](#vmoperator)
  * [`vmagent`](#vmagent)
  * [`vmalert`](#vmalert)
  * [`vmalertmanager`](#vmalertmanager)
  * [`vmauth`](#vmauth)
  * [`vmsingle`](#vmsingle)
  * [`vmuser`](#vmuser)

You can install these components in various combinations.

But the most common deployment options can have the following combined scenarios:

* Deploy Monitoring only for collect metrics.
* Deploy Monitoring for collect and watch metrics.
* Deploy Monitoring for collect metrics, watch them, and alert them states.
* Deploy Monitoring for collect metrics and send them to external storage.

[Back to TOC](#table-of-content)

### Default deploy settings

Monitoring can be deploy without specifying any parameters for deploy. Even if you will specify any parameters
for all other parameters will be using defaults from `values.yaml` file.

By default Monitoring deploy with components:

| Component             | Requests                | Limits                  |
| --------------------- | ----------------------- | ----------------------- |
| `alertmanager`        | CPU: 100m, RAM 100Mi    | CPU: 200m, RAM: 200Mi   |
| `grafana`             | CPU: 300m, RAM: 400Mi   | CPU: 500m, RAM: 800Mi   |
| `grafana-operator`    | CPU: 50m, RAM: 50Mi     | CPU: 100m, RAM: 100Mi   |
| `kube-state-metrics`  | CPU: 50m, RAM: 50Mi     | CPU: 100m, RAM: 256Mi   |
| `monitoring-operator` | CPU: 50m, RAM: 50Mi     | CPU: 100m, RAM: 150Mi   |
| `node-exporter`       | CPU: 50m, RAM: 50Mi     | CPU: 100m, RAM: 100Mi   |
| `victoriametris`      | CPU: 1000m, RAM 3 Gi    | CPU: 1500m, RAM: 5Gi    |
| `vm-operator`         | CPU: 100m, RAM: 100Mi   | CPU: 200m, RAM: 200Mi   |
| **Total**             | CPU: 1700m, RAM: 3872Mi | CPU: 2800m, RAM: 6926Mi |

**Note**: We are strongly recommending check these requests and limits, and change if necessary for your environment.

Also by default:

* Victoriametrics and Grafana will deploy with **ephemeral storage** (with using emptyDir)
* Victoriametrics will collect and store metrics **24 hours**
* Victoriametrics will collect metrics from (full list see in documents [Metrics OOB](metrics-oob.md))
  * Kubelet and cAdviser
  * Kubernetes APIServer
  * Nginx Ingress
  * Etcd
  * Kube-state-metrics
  * Node Exporter
  * Self metrics: Victoriametrics, AlertManager, Grafana
* In Grafana will create dashboards, dashboard list and descriptions see in documents [Dashboard OOB](dashboards-oob)
* In Victoriametrics will create alerts, list of alerts see in
  document [Alerts OOB](alerts-oob.md)
* Without any ingress (in case deploy in Kubernetes)
* Without any enabled and configured integrations
* Without blackbox and cert exporters
* Without promxy
* Without autoscaling

Full list of default values you can find in `values.yaml` file in `monitoring-operator` helm chart.

To read which parameters can be changed during deploy or update, refer to the [Configuration](#configuration) section.

[Back to TOC](#table-of-content)

### Full Stack Deploy

By default, Monitoring is deployed with components without any parameters:

* `alertmanager`
* `grafana`
* `grafana-operator`
* `kube-state-metrics`
* `monitoring-operator`
* `node-exporter`
* `victoriametrics`
* `vm-operator`

To deploy all the components, it is not required to specify any parameters. However, it may be useful to specify
parameters to:

* Create ingresses during the deployment.
* Specify SecurityContext to deploy with a restricted PSP.

For example, to deploy Monitoring into Kubernetes, the following parameters can be specified:

```yaml
grafana:
  ingress:
    install: true
    host: grafana.cloud.org
victoriametrics:
  vmsingle:
    install: true
    ingress:
      install: true
      host: victoriametrics.cloud.org
  vmAlertManager:
    install: true
    ingress:
      install: true
      host: alertmanager.cloud.org
```

All other parameter will be used from default `values.yaml`. The part of default settings described in section
[Default deploy settings](#default-deploy-settings)

[Back to TOC](#table-of-content)

### Enabling and disabling dashboards, rules and monitors

You can enable or disable Grafana dashboards, Prometheus rules and Service or Pod monitors manually by using the
following parameters:

* [`grafanaDashboards`](#grafana-dashboards): parameter `grafanaDashboards.list` contains names of Grafana dashboards
  that will be installed by the monitoring-operator;
* [`prometheusRules`](#prometheus-rules): parameter `prometheusRules.ruleGroups` contains group names of
  Prometheus rules that will be installed by the monitoring-operator. You can find these groups in the
  [Alerts OOB document](alerts-oob.md#monitoring-operator).
  Also, you can override each group by `prometheusRules.override` parameter.
  You can find more information in the [alerts configuration document](user-guides/monitoring-configuration/alerts.md);
* [`kubernetesMonitors`](#kubernetes-monitors): parameter `kubernetesMonitors` contains names of
  Kubernetes monitors that will be installed by the monitoring-operator. Also, you can override `interval`,
  `scrapeTimeout`, `metricRelabelings` and `relabelings` for each monitor. Values will be overridden for each endpoint.

You can enable required and disable useless dashboards, rules and monitors for installation in _public cloud_ with one
parameter. See more in
[Enabling/disabling specific monitoring components](integrations.md#enablingdisabling-specific-monitoring-components).

Also, there is automatic features for disabling specific dashboards and monitors in the monitoring-operator in some
cases. The following rules work for dashboards:

* `kubernetes-nodes-resources` dashboard isn't installed if the `node-exporter` isn't installed;
* `home-dashboard` dashboard isn't installed if `grafana.grafanaHomeDashboard` set as false;
* `core-dns-dashboard` dashboard isn't installed if `coreDnsServiceMonitor` Service monitor isn't installed or
  the monitoring-operator is installed in the OpenShift v3.11 or below;
* `ingress-request-handling-performance`, `ingress-nginx-controller` and `ingress-list-of-ingresses` dashboards aren't
  installed if `nginxIngressPodMonitor` Pod monitor isn't installed or
  the monitoring-operator is installed in the OpenShift.
* `Victoria Metrics` dashboards are installed with VmOperator, VmSingle, VmAlert and VmAgent only.
* `Prometheus Self Monitoring` dashboard is installed with Prometheus only.

The following rules work for monitors:

* `coreDnsServiceMonitor` Service monitor isn't installed if the monitoring-operator is installed
  in the OpenShift v3.11 or below;
* `nginxIngressPodMonitor` Pod monitor isn't installed if the monitoring-operator is installed in the OpenShift.

[Back to TOC](#table-of-content)

### Deploy with cert-exporter

#### Types of Checked Certificates

The deployment with `cert-exporter` allows checking ssl/tls certificates and their expiration dates from the following
sources:

* Files on the host file system
* Kubeconfig files
* Kubernetes secrets

##### Certificates from Files

The parsing of certificates from the host file system can be configured by the `certExporter.certsInFiles`
section of the configuration. It is enabled by default.

**Note**: If the collecting of certificates from files, the part of `cert-exporter` will be deployed as a
DaemonSet. You can specify some parameters especially for it by using `certExporter.daemonset` section of parameters.

The exporter can parse default internal certificates from cluster nodes. These types of certificates require Kubernetes
for authentication over TLS and for the Kubernetes cluster to work properly. The parsing of default internal
certificates can be enabled by the `certExporter.certsInFiles.defaultCerts` parameter. By default, the parameter is set
to "true."
This type contains the following certificates:

* For the Kubernetes cluster, files with `.crt` and `.pem` extensions in the `/etc/kubernetes`
  and `/var/lib/kubelet/pki` directories and their child directories.
* For the OpenShift cluster, files with `.crt` and `.pem` extensions in the `/etc/origin` and `/etc/etcd` directories
  and their child directories.

**Note**: The directories with default internal certificates on OpenShift nodes , and on Kubernetes nodes in some cases,
can be read only by the _root_ user. So `cert-exporter` without root rights will not be able to access this part of the
certificates. For more information about this limitation,
see [Installation with Root Rights](#installation-with-root-rights).

Also, you can monitor other certificates from the host file system. To do this, mount custom host path volumes with
the `certExporter.additionalHostPathVolumes` parameter and make regex for matching the certificates of interest with
the `certExporter.certsInFiles.includeCerts` parameter.

For more information about the format of regex, see
[https://github.com/bmatcuk/doublestar#usage](https://github.com/bmatcuk/doublestar#usage).

For more information about the configuration of `cert-exporter`, refer to the [cert-exporter](#cert-exporter) section.

##### Certificates from Kubeconfig

The parsing of certificates from the kubeconfig files from the host file system can be configured by the
`certExporter.certsInKubeconfig` section of the configuration. It is disabled by default because kubeconfig files are
usually root readable.

**Note**: If the collecting of certificates from Kubeconfig, the part of `cert-exporter` will be deployed as a
DaemonSet. You can specify some parameters especially for it by using `certExporter.daemonset` section of parameters.

The exporter can parse certificates from the default kubeconfig files. The parsing of kubeconfig files on the default
paths can be enabled by the `certExporter.certsInKubeconfig.defaultCerts` parameter. By default, it is set to "true."
This type contains the following files:

* For the Kubernetes cluster, files with the `.conf` extension in the `/etc/kubernetes` directory and its child
  directories.
* For the OpenShift cluster, files with the `.kubeconfig` extension in the `/etc/origin` directory and its child
  directories.
* For both the clusters, the `/root/.kube/config` file.

**Note**: Kubeconfig files usually can be read only by the _root_ user, so `cert-exporter` without root rights will not
be able to access these files. For more information about this limitation, see
[Installation with Root Rights](#installation-with-root-rights).

Also, you can monitor the kubeconfig files from other paths on the host file system. To do this, mount the custom host
path volumes with the `certExporter.additionalHostPathVolumes` parameter and make regex for matching the files of
interest with the `certExporter.certsInFiles.includeCerts` parameter.

For more information about the format of regex,
see [https://github.com/bmatcuk/doublestar#usage](https://github.com/bmatcuk/doublestar#usage).

For more information about the configuration of `cert-exporter`, refer to the [cert-exporter](#cert-exporter) section.

##### Certificates from Kubernetes Secrets

The parsing of certificates from Kubernetes secrets can be configured by the `certExporter.certsInSecrets`
section of the configuration. It is enabled by default.

**Note**: If the collecting of certificates from secrets, the part of `cert-exporter` will be deployed as a Deployment.
You can specify some parameters especially for it by using `certExporter.deployment` section of parameters.

The exporter checks the secret content and chooses the fields that match the given regex (`*.crt` by default). You can
specify the regex for the matching fields in secrets with the `certExporter.certsInSecrets.includeCerts` parameter.

**Important!** The exporter **doesn't use doublestar module to match secrets** for some reason, so patterns we can use
for `certsInFiles` or `certsInKubeconfig` don't work here. It uses [Golang filepath](https://pkg.go.dev/path/filepath)
notation instead. That means we can't use one parameter to set 2 or more comma separated patterns in `includeCerts`
to match files with different file extensions. But there is the way to avoid this restriction: you can use
`certExporter.deployment.extraArgs` to set additional `--secrets-include-glob` parameters. You can set several secrets
include parameters, and they will work as if they're connected via boolean AND.

The note above means that if you want to watch cert files with `.crt`, `.cer` and `.pem` extensions, you can use the
following configuration:

```yaml
certExporter:
  ...
  install: true
  certsInSecrets:
    enabled: true
    includeCerts: "*.crt"
  ...
  deployment:
    extraArgs:
      - '--secrets-include-glob=*.cer'
      - '--secrets-include-glob=*.pem'
  ...
```

By default, the exporter can parse any types of secrets from all the available namespaces on the cluster. The types of
secrets for checking can be limited by the `certExporter.certsInSecrets.types` parameter. The namespaces for parsing
secrets can be specified by `certExporter.certsInSecrets.namespaces`.

**Note**: You cannot specify more than one namespace at the same time. You can choose only one namespace or all of them.

The exporter can find secrets with specific annotations or labels that can be configured with
the `certExporter.certsInSecrets.annotationSelector` or `certExporter.certsInSecrets.labelSelector` parameters. Also,
the exporter can find and export cert-manager certificates. For more information, see
[https://cert-manager.io/docs/](https://cert-manager.io/docs/).

Different parameters are required for different versions of the cert-manager API due to changes in the way the secrets
were stored. To monitor such certificates, configure the appropriate fields:

* **cert-manager.k8s.io/v1alpha1** `certExporter.certsInSecrets.labelSelector: "certmanager.k8s.io/certificate-name"`
* **cert-manager.io/v1alpha2** `certExporter.certsInSecrets.annotationSelector: "cert-manager.io/certificate-name"`

For more information about the configuration of the `cert-exporter`, refer to the [cert-exporter](#cert-exporter)
section.

#### Installation with Root Rights

For security reasons, some directories and files on the host file system are readable only for root users. In such cases
to collect information about certificates, you need to run the `cert-exporter` daemonset with root rights.

The deployment of `cert-exporter` doesn't need for root rights because it is aimed only on certificates in secrets.

To deploy the exporter with root rights, you should specify `daemonset.securityContext.runAsUser: "0"` in the parameters
in the `certExporter` section:

```yaml
...
certExporter:
  ...
  daemonset:
    securityContext:
      runAsUser: "0"
```

In some cases, the `cert-exporter` without root rights may not have access to some files, although the rest of the files
from the same directory are available. In this case, the exporter signals an error in the pod logs on the corresponding
node, as well as on the Errors panel on the dashboard.

[Back to TOC](#table-of-content)

## Integrations

This section describes integrations with various external systems for the following purposes:

* Store metrics for a long term
* Alerting systems
* Visualization of systems

You can find more information about integrations in the [integrations guide](integrations.md).

### Metrics Collectors and Long Term Storages

Prometheus is a system for collecting metrics and getting a prompt response on any unexpected situations. It can store
metrics into an internal TSBD (Prometheus use OpenTSDB). But in the case of a large number of metrics, it is recommended
to use external storage if it is required to store metrics for a longer period, as Prometheus can store metrics only for
a short period like a couple of weeks or months.

#### Common Restrictions and Recommendations

Before you start reading about various integrations and how to configure them, it is important to know the common
restrictions and recommendations.

##### Filter Pushed Metrics

Firstly, if integration uses `Remote Write`, it is strongly recommended adding relabeling config and send only the
metrics that you really want to see. Prometheus can collect 10k or 100k metrics, but it is not required to send them all
to the external storage. Moreover, `graphite-remote-adapter` or the external storage may not withstand such a load, and
start to use a lot of resources or start to skip some metrics.

For example, you can filter sending metrics as follows:

```yaml
prometheus:
  remoteWrite:
    - url: 'http://graphite-remote-adapter:9201/write'
      remoteTimeout: 45s
      writeRelabelConfigs:
        - action: keep
          sourceLabels:
            - job
          regex: (node-exporter|kubelet|kube-state-metrics)
```

In this example, `Remote Write` sends only the metrics with the following labels:

* `job=node-exporter` - Metrics that contains information about VMs.
* `job=kubelet` - Metrics that contains information about pods and nodes.
* `job=kube-state-metrics` - Metrics that contains information about Deployments, DaemonSets, Pods, and so on.

[Back to TOC](#table-of-content)

##### RemoteWrite Increases Memory Usage

Secondly, using remote write increases the memory footprint of Prometheus. Most users report ~25% increase in the memory
usage, but that number is dependent on the shape of the data. For each series in the WAL, the remote write code caches a
mapping of series ID to label values, causing large amounts of series churn to significantly increase the memory usage.

In addition to the series cache, each shard and its queue increases the memory usage. Shard memory is proportional to
the `number of shards * (capacity + max_samples_per_send)`. When tuning, consider reducing `max_shards` alongside
increases to the capacity and `max_samples_per_send`, to avoid inadvertently running out of memory. The default values
for `capacity: 2500` and `max_samples_per_send: 500` constrain the shard memory usage to less than `500 kB` per shard.

For more information about `Remote Write` tuning and which parameters can be used, refer to the official Prometheus
documentation
at [https://prometheus.io/docs/practices/remote_write/](https://prometheus.io/docs/practices/remote_write/).

[Back to TOC](#table-of-content)

##### Add External Labels

Thirdly, it is strongly recommended adding additional labels when using RemoveWrite. These labels allow filtering
metrics for the necessary environment into a central storage.

For example, you can add labels as follows:

```yaml
prometheus:
  externalLabels:
    <key_1>: <value_1>
    <key_2>: <value_2>
    ...
    <key_N>: <value_N>
```

[Back to TOC](#table-of-content)

##### Use Authentication for RemoteWrite and RemoteRead if Remote Storage Supports it

Some remote storages support authentication on their endpoints. And if a remote storage support authentication, it is
strongly recommended to configure it to avoid store and send data by unauthorized clients.

Monitoring supports two common authentication scenarios:

* Basic Authentication
* Authentication with token

To specify the credentials for basic authentication, specify the following parameters:

```yaml
prometheus:
  remoteWrite:
    - url: 'http://<remote_storage_addr>:<port>/api/v1/write'
      basicAuth:
        username:
          name: your-secret-name
          key: username
        password:
          name: your-secret-name
          key: password
  remoteRead:
    - url: 'http://<remote_storage_addr>:<port>/api/v1/read'
      basicAuth:
        username:
          name: your-secret-name
          key: username
        password:
          name: your-secret-name
          key: password
```

Where, Secret contains the following fields:

```yaml
kind: Secret
apiVersion: v1
metadata:
  name: prometheus-remote-storage-basic-auth
type: Opaque
data:
  user: "<username>"
  password: "<password>"
```

**Warning**: Currently monitoring deploy does not allow creating secret with credentials for the external storage. You
should create it manually before the deployment.

To specify a token for authentication, specify the following parameters:

```yaml
prometheus:
  remoteWrite:
    - url: 'http://<remote_storage_addr>:<port>/api/v1/write'
      bearerToken: <token>
  remoteRead:
    - url: 'http://<remote_storage_addr>:<port>/api/v1/read'
      bearerToken: <token>
```

[Back to TOC](#table-of-content)

#### Graphite

For integration with Graphite into Monitoring, the chart has the ability to deploy the `graphite-remote-adapter`
service.

The following image illustrates how the metrics are sent during the use of this service in an external server.

![Prometheus to Graphite](/docs/public/images/prometheus_k8s_integrations_graphite_clickhouse.png)

To setup this integration, you need to configure:

* Remote write section into Prometheus. It needs to send the metrics from Prometheus to `graphite-remote-adapter`.
* The `graphite-remote-adapter` section. It contains deploy and runtime parameters for `graphite-remote-adapter`.

For example:

```yaml
prometheus:
  externalLabels:
    cluster: my_awesome_cluster
  remoteWrite:
    - url: "http://graphite-remote-adapter:9201/write"
      queueConfig:
        maxSamplesPerSend: 10000
        capacity: 20000
        maxShards: 30
...
graphite_remote_adapter:
  install: true
  writeCarbonAddress: <carbonapi_addr>:9999
```

Where:

* `http://graphite-remote-adapter:9201/write` is a service name of `graphite-remote-adapter` which Monitoring creates
  during the deployment. So you can always use it to deploy `graphite-remote-adapter`.

The full list of parameters can be viewed at:

* [Prometheus](#prometheus)
* [Graphite-remote-adapter](#graphite-remote-adapter)

[Back to TOC](#table-of-content)

#### InfluxDB

For integration with InfluxDB into Monitoring, the chart has the ability to specify RemoteWrite and RemoteRead settings.

The following image illustrates how the metrics are sent during the use of this service in an external server.

![Prometheus to InfluxDB](/docs/public/images/prometheus_k8s_integration_influxdb.png)

For more information about the integration with InfluxDB, refer to the official documentation at
[https://docs.influxdata.com/influxdb/v1.8/supported_protocols/prometheus/](https://docs.influxdata.com/influxdb/v1.8/supported_protocols/prometheus/).

To setup this integration, you need to configure the `remoteWrite` or/and `remoteRead` section into Prometheus.
It needs to exchange metrics between Prometheus and InfluxDB.

For example:

```yaml
prometheus:
  externalLabels:
    cluster: my_awesome_cluster
  remoteWrite:
    - url: "http://<influxdb_addr>:8086/api/v1/prom/write?db=<database>"
  remoteRead:
    - url: "http://<influxdb_addr>:8086/api/v1/prom/read?db=<database>"
```

**Note**:

* InfluxDB had some [bugs](https://github.com/prometheus/prometheus/issues/5610) related to remote read.
  It was fixed since 1.8.x version.
* Prometheus can read only these metrics which it writes or in format which satisfy it conditions.
* Prometheus can't read all other metrics which stored in InfluxDB.
* Prometheus can read metrics only one database.

The full list of parameters can be viewed at:

* [Prometheus](#prometheus)

[Back to TOC](#table-of-content)

#### VictoriaMetrics

For integration with VictoriaMetrics into Monitoring, the chart has the ability to specify RemoteWrite settings.

The following image illustrates how the metrics are sent during the use of this service in an external server.

![Prometheus to VictoriaMetrics](/docs/public/images/prometheus_k8s_integration_victoriametrics.png)

To setup this integration, you need to configure:

* The Remote write section into Prometheus. It needs to send the metrics from Prometheus to VictoriaMetrics.

For example:

```yaml
prometheus:
  externalLabels:
    cluster: my_awesome_cluster
  remoteWrite:
    - url: "http://<victoriametrics_addr>:8428/api/v1/write"
      queueConfig:
        maxSamplesPerSend: 10000
        capacity: 20000
        maxShards: 30
```

The full list of parameters can be viewed at:

* [Prometheus](#prometheus)

[Back to TOC](#table-of-content)

#### Prometheus Federation

The Prometheus Federation allows a Prometheus server to scrape a selected time series from another Prometheus server. So
there are two options to setup Prometheus Federation:

* To collect metrics by Prometheus in Monitoring from any other Prometheus.
* To collect metrics by any external Prometheus from Prometheus in Monitoring.

To configure metrics collection by any external Prometheus from Prometheus in Monitoring, no actions are required during
the deployment. All parameters should be configured on an external Prometheus. For an external Prometheus, use any
vanilla image without a configuration manager. For more information, refer to the official documentation at
[https://prometheus.io/docs/prometheus/latest/federation/#configuring-federation](https://prometheus.io/docs/prometheus/latest/federation/#configuring-federation).

Also need add settings for external labels in Prometheus. Prometheus will add these labels for all metrics which
try to collect external systems.

For example:

```yaml
prometheus:
  externalLabels:
    cluster: my_awesome_cluster
```

To configure metrics collection by Prometheus in Monitoring from any other Prometheus, create a ServiceMonitor custom
resource.

The following image illustrate how the metrics is scraped from any other Prometheus:

![Prometheus Federation](/docs/public/images/prometheus_k8s_integration_prometheus_federation.png)

For example:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  labels:
    app.kubernetes.io/component: monitoring
  name: AnyInternalPrometheus
spec:
  endpoints:
    - interval: 30s
      port: web
      path: /federate
      params:
        - 'match[]':
            - '{job="prometheus"}'
            - '{__name__=~"job:.*"}'
  jobLabel: federate
  selector:
    matchLabels:
      app: prometheus
```

The full list of parameters can be viewed at:

* [Prometheus](#prometheus)

[Back to TOC](#table-of-content)

### Alerting Systems

#### External AlertManager

##### AlertManager in the same Cloud

**Note**: AlertManager may not be deployed with Monitoring stack for integration with an external AlertManager.

The AlertManager handles alerts sent by client applications such as the Prometheus/Victoriametrics server.
It takes care of deduplicating, grouping, and routing them to the correct receiver integration such as email,
PagerDuty, or OpsGenie. It also takes care of silencing and inhibition of alerts.

There are two options to deploy AlertManager:

* Deploy it into cloud.
* Deploy on any external server as a service or a docker.

During work with external AlertManager, Prometheus calculates the rules and sends only events to the AlertManager. The
AlertManager should handle these events and send alerts as notifications by any configured channel.

To configure the integration with AlertManager, which deploys into the same cloud, specify the parameters into
the `alerting` section in the Prometheus and into the `vmAlertManager` section in the Victoriametrics. For example:

```yaml
prometheus:
  alerting:
    alertmanagers:
      - namespace: upstream-monitoring
        name: alertmanager-example
        port: web
```

```yaml
victoriametrics:
  vmAlertManager:
    install: true
```

Where:

* `namespace` is a namespace where external AlertManager is deployed.
* `name` is a name of an external AlertManager's service.
* `port` is a name of a port in an external AlertManager's service.

##### AlertManager outside the Cloud

To configure the integration with AlertManager, which deploys on any external machine, use
`AdditionalAlertManagerConfigs`. It is a Secret with name `additional-alertmanager-configs`, which is created into the
namespace with monitoring with deploy Monitoring stack.

**Note**: Currently, a deployment process of monitoring stack does not provide the ability to fill this Secret. It is
always created empty and may be filled manually after the deployment.

In Secret, you can specify a part of Prometheus config file with AlertManager settings. For example:

```yaml
- scheme: https
  static_configs:
    - targets:
        - "alertmanager.example.outside.cloud.org"
  tls_config:
    insecure_skip_verify: true
```

Remember that the Secret data should be encoded by base64 to store into secret. For example, in Linux the following
command can be used to encode content of the file with config:

```bash
cat alertmanager.yaml | base64
```

All configurations from `AdditionalAlertManagerConfigs` are applied automatically and without restarting Prometheus.

[Back to TOC](#table-of-content)

## Storage

**Note**: Ability to specify storage was added since version `0.4.0`

To maintain data across deployments and version upgrades, the data must be persisted to some volume other
than `emptyDir`, allowing it to be reused by Pods after an upgrade.

Kubernetes supports several kinds of storage volumes. The Monitoring Operator uses Prometheus/Victoriametrics Operator
to control the Prometheus/Vicoriametrics deployment.
The Prometheus/Vicoriametrics Operator works with `PersistentVolumeClaims`, which supports the underlying
`PersistentVolume` to be provisioned when requested.

This document assumes a basic understanding of `PersistentVolumes`, `PersistentVolumeClaims`, and their provisioning.

### Known issues

If you deploy monitoring with Prometheus Operation with version `v0.34.0` or less, you may get problems described in the
following GitHub issue:

* Creation/Update of Prometheus object throws error on volume claim template creationTimeStamp being null. For more
  information, refer
  to [https://github.com/prometheus-operator/prometheus-operator/issues/2824](https://github.com/prometheus-operator/prometheus-operator/issues/2824)
* Setting podMetadata labels without creation timestamp fails from golang. For more information, refer
  to [https://github.com/prometheus-operator/prometheus-operator/issues/2399](https://github.com/prometheus-operator/prometheus-operator/issues/2399)

This issue was fixed since version `v0.35.0`. To create volume in an older version, you need to manually
add `creationTimeStamp` to `volumeClaimTemplate`.

For example:

```yaml
prometheus:
  ...
  storage:
    volumeClaimTemplate:
      metadata:
        creationTimestamp: '2020-06-23T18:55:02Z'
      spec:
        ...
```

Also if you want to use a `hostPath` with non default users (which you can specify with using `securityContext`), you
may get problems described in the following GitHub issues:

* Permission denied writing to mount using volumeClaimTemplate. For more information, refer
  to [https://github.com/prometheus-operator/prometheus-operator/issues/966](https://github.com/prometheus-operator/prometheus-operator/issues/966)
* mkdir /prometheus/wal: permission denied. For more information, refer
  to [https://github.com/helm/charts/issues/12176](https://github.com/helm/charts/issues/12176).

These problems occur because the Dockerfile of vanilla Prometheus image contains the following logic:

* Without `USER` directive (it means under `root`), create a directory for the Prometheus data.
* Changes ownership and grants it to a user and group `nobody` (uid:gid = 65534:65534).
* Execute directive `USER nobody` and run entrypoint under user `nobody`.

It means that the directory is available only for the user `nobody` (65534).

But if you still need to run with a custom user, you must create `prometheus-db` in PV folder and own it to a custom
user. For example, if hostPath PV is created in the directory `/mnt/data/prometheus/` you should create a directory as:

```bash
/mnt/data/prometheus/prometheus-db
```

and own it to user, for example with uid and gid `2001:2001`:

```bash
chown -R 2001:2001 /mnt/data/prometheus
```

[Back to TOC](#table-of-content)

### Storage with Dynamic Provisioning

Automatic provisioning of storage requires an already existing `StorageClass`.

For best results, use volumes that have high I/O throughput.

The `StorageClass` that was created can be specified in the storage section in the Prometheus/Victoriametrics resource:

```yaml
prometheus:
  storage:
    volumeClaimTemplate:
      spec:
        # Specify storage class to create volume
        storageClassName: nfs-dynamic-provisioning
        resources:
          requests:
            storage: 10Gi
```

```yaml
victoriametrics:
  vmSingle:
    storage:
      # Specify storage class to create volume
      storageClassName: nfs-dynamic-provisioning
      accessModes:
        - ReadWriteOnce
      resources:
        requests:
          storage: 1Gi
      selector:
        matchLabels:
          app.kubernetes.io/name: vmsingle
      volumeName: data-victoriametrics-pv
```

[Back to TOC](#table-of-content)

### Manual Storage Provisioning

The monitoring deploy parameters allow you to support arbitrary storage through a PersistentVolumeClaim.

The easiest way to use a volume that cannot be automatically provisioned (for whatever reason) is to use a label
selector alongside a manually created PersistentVolume.

For example, using `hostPath` persistence volume might be accomplished with the following specifications:

```yaml
prometheus:
  # Vanilla Prometheus image use user and group nobody = 65534
  # So for use PV it is better to use user nobody = 65534
  securityContext:
    fsGroup: 65534
    runAsUser: 65534
  # Because hostPath PV created on specific node, we must bind Prometheus on this node
  nodeSelector:
    kubernetes.io/hostname: worker1
  storage:
    volumeClaimTemplate:
      spec:
        resources:
          requests:
            storage: 10Gi
        selector:
          # Match PV by label on PV
          matchLabels:
            app: prometheus
```

```yaml
victoriametrics:
  vmSingle:
    storage:
      storageClassName: manual
      resources:
        requests:
          storage: 10Gi
      selector:
        # Match PV by label on PV
        matchLabels:
          app: victoriametrics
```

[Back to TOC](#table-of-content)

### Disabling Default StorageClasses

To manually provision volumes (as of Kubernetes 1.6.0), you may need to disable the default StorageClass that is
automatically created for certain Cloud Providers.

The default StorageClass behavior overrides manual storage provisioning, preventing PersistentVolumeClaims from
automatically binding to manually created PersistentVolumes.

To override this behavior, you must explicitly create the same resource, but set it to not be default.

For example, to disable default StorageClasses on a minikube, edit the following StorageClass:

```yaml
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  namespace: kube-system
  name: standard
  annotations:
    # disable this default storage class by setting this annotation to false.
    storageclass.beta.kubernetes.io/is-default-class: "false"
  labels:
    addonmanager.kubernetes.io/mode: Reconcile
provisioner: k8s.io/minikube-hostpath
```

[Back to TOC](#table-of-content)

## Configuration

This section describes all the parameters that can be used to configure Monitoring.

### root parameters

<!-- markdownlint-disable line-length -->
| Field                   | Description                                                                                                                                                                                                                                                                                                                                                                                                            | Scheme |
| ----------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------ |
| name                    | Provide a name in place of monitoring-operator for `app:` labels.                                                                                                                                                                                                                                                                                                                                                      | string |
| fullnameOverride        | Provide a name to substitute for the full names of resources.                                                                                                                                                                                                                                                                                                                                                          | string |
| global.privilegedRights | It indicates whether monitoring-operator has permissions to deploy ClusterRole and ClusterRoleBinding resources. If set to `true`, it deploys all resources. Otherwise, it deploys only Role and RoleBinding resources. The ClusterRole and ClusterRoleBinding resources were deployed manually. Also skip ClusterRole and ClusterRole binding resources deployed for the monitoring-operator itself if set to `false` | bool   |
| role.install            | Allow to disable create Role and ClusterRole for `monitoring-operator` during deploy. If `global.privilegedRights` parameter is set to false, ClusterRole will not be installed in any case.                                                                                                                                                                                                                           | bool   |
| roleBinding.install     | Allow to disable create RoleBinding and ClusterRoleBinding for `monitoring-operator` during deploy. If `global.privilegedRights` parameter is set to false, ClusterRoleBinding will not be installed in any case.                                                                                                                                                                                                      | bool   |
| publicCloudName         | Allows to disable or enable specific dashboards, rules and service monitors for the specified public cloud. Allowed values: `aws`, `azure`, `google`, `""`                                                                                                                                                                                                                                                             | string |
<!-- markdownlint-enable line-length -->

Example:

```yaml
nameOverride: "monitoring-operator"
fullnameOverride: "monitoring-operator"
role:
  install: true
roleBinding:
  install: true
publicCloudName: ""
```

[Back to TOC](#table-of-content)

### alertmanager

AlertManager is a specification of the desired behavior of the AlertManager cluster.

<!-- markdownlint-disable line-length -->
| Field             | Description                                                                                                                                                                                                            | Scheme                                                                                                                       |
| ----------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| install           | Allows to disable deploy AlertManager. If AlertManager was not deployed during the deployment using helm, it can be deployed using change custom resource PlatformMonitoring.                                          | bool                                                                                                                         |
| paused            | Set paused to reconciliation.                                                                                                                                                                                          | bool                                                                                                                         |
| image             | A docker image to use for AlertManager deployment.                                                                                                                                                                     | string                                                                                                                       |
| ingress           | Ingress allows to create Ingress for the AlertManager UI.                                                                                                                                                              | *[Ingress](#ingress)                                                                                                         |
| nodeSelector      | Defines which nodes the pods are scheduled on. Specified just as map[string]string. For example: "type: compute"                                                                                                       | map[string]string                                                                                                            |
| affinity                                            | If specified, the pod's scheduling constraints                                                                                                                                                                                      | *v1.Affinity                                                                                                                   |
| annotations       | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value" | map[string]string                                                                                                            |
| labels            | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                       | map[string]string                                                                                                            |
| port              | The port for AlertManager service.                                                                                                                                                                                     | int                                                                                                                          |
| replicas          | Set replicas.                                                                                                                                                                                                          | *int32                                                                                                                       |
| resources         | The resources that describe the compute resource requests and limits for single Pods.                                                                                                                                  | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#resourcerequirements-v1-core) |
| securityContext   | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                     | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#podsecuritycontext-v1-core)    |
| podMonitor        | Pod monitor for self monitoring.                                                                                                                                                                                       | *[Monitor](#monitor)                                                                                                         |
| priorityClassName | PriorityClassName assigned to the Pods to prevent them from evicting.                                                                                                                                                  | string                                                                                                                       |
<!-- markdownlint-enable line-length -->

Example:

```yaml
alertManager:
  install: true
  paused: false
  image: prom/alertmanager:v0.19.0
  port: 30903
  resources:
    limits:
      cpu: 200m
      memory: 200Mi
    requests:
      cpu: 100m
      memory: 100Mi
  ingress:
    ...see example by link...
  nodeSelector:
    node-role.kubernetes.io/worker: worker
  labels:
    label.key: label-value
  annotations:
    annotation.key: annotation-value 
  replicas: 1
  securityContext:
    runAsUser: 2000
    fsGroup: 2000
  podMonitor:
    ...see example by link...
```

[Back to TOC](#table-of-content)

### auth

Auth is a specification of the desired behavior of the authentication for Grafana, Prometheus, and AlertManager.

**Important**: Old versions of Prometheus and AlertManager has no authentication. So by default (without setup oAuth)
Prometheus and AlertManager will not have any authentication, and in Grafana will be enable Basic Auth with default
credentials (admin / admin).

<!-- markdownlint-disable line-length -->
| Field                        | Description                                                                                        | Scheme                                                                                                                  |
| ---------------------------- | -------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------- |
| clientId                     | A public identifier of client for OAuth.                                                           | string                                                                                                                  |
| clientSecret                 | A secret known only to the client and the authorization server.                                    | string                                                                                                                  |
| loginUrl                     | URL of identity-provider to log in. Usually configured as: `http://<ip>/login`                     | string                                                                                                                  |
| tokenUrl                     | URL of identity-provider to request by client. Usually configured as: `http://<ip>/token`          | string                                                                                                                  |
| userInfoUrl                  | URL of identity-provider to get user-info by client. Usually configured as: `http://<ip>/userinfo` | string                                                                                                                  |
| tlsConfig                    | TLSConfig extends the safe TLS configuration with file parameters.                                 | object                                                                                                                  |
| tlsConfig.caSecret           | Secret containing the CA certificate to use for the targets.                                       | *[v1.SecretKeySelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#secretkeyselector-v1-core) |
| tlsConfig.certSecret         | Secret containing the client certificate file for the targets.                                     | *[v1.SecretKeySelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#secretkeyselector-v1-core) |
| tlsConfig.keySecret          | Secret containing the client key file for the targets.                                             | *[v1.SecretKeySelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#secretkeyselector-v1-core) |
| tlsConfig.insecureSkipVerify | Disable certificate validation.                                                                    | bool                                                                                                                    |
<!-- markdownlint-enable line-length -->

**Note**: Pay attention that Secret which will store certificates should be created manually. For example you
can create it with using a template:

```yaml
apiVersion: v1
data:
  trusted-ca.crt: <ca_certificate_content>
  client-cert.crt: <client_certificate_content>
  client-key.key: <client_key_content>
kind: Secret
metadata:
  labels:
    app.kubernetes.io/name: grafana-client-certificates
  name: grafana-client-certificates
type: Opaque
```

Example:

```yaml
auth:
  clientId: xxx
  clientSecret: xxx
  loginUrl: http://1.2.3.4/authorize
  tokenUrl: http://1.2.3.4/token
  userInfoUrl: http://1.2.3.4/userinfo

  ## Settings for TLS
  tlsConfig:
    caSecret:
      name: grafana-client-certificates
      key: trusted-ca.crt
    certSecret:
      name: grafana-client-certificates
      key: client-cert.crt
    keySecret:
      name: grafana-client-certificates
      key: client-key.key
    insecureSkipVerify: false
```

Grafana reads `clientId` and `clientSecret` from Grafana Custom Resource or from environment variables.
Since storing credentials information in Custom Resource is not recommended,
it is necessary to read `clientId` and `clientSecret` from grafana environment variables.
System will automatically populate `GF_AUTH_GENERIC_OAUTH_CLIENT_ID` and `GF_AUTH_GENERIC_OAUTH_CLIENT_SECRET`
environment variables via `grafana-extra-vars-secret` using the values from `.Values.auth` section.

[Back to TOC](#table-of-content)

### blackbox-exporter

Blackbox exporter allow to execute probe to URL (static or for Ingress) as to black box and meter state of answer,
code and response time. Also Blackbox exporter check SSL certificate state (if probe configure for URL
with SSL certificate).

Details, how to configure Probes for use with Blackbox exporter please read in
[Configuration documentation](configuration.md#probe).

<!-- markdownlint-disable line-length -->
| Field                                               | Description                                                                                                                                                                                                                         | Scheme                                                                                                                         |
| --------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------ |
| install                                             | Allows to enable or disable deploy blackbox-exporter.                                                                                                                                                                               | bool                                                                                                                           |
| name                                                | A name of the microservice to deploy with. This name is used as the name of the microservice deployment and in labels.                                                                                                              | string                                                                                                                         |
| image                                               | A Docker image to deploy the blackbox-exporter.                                                                                                                                                                                     | string                                                                                                                         |
| asDaemonSet                                         | Allows deploying blackbox-exporter as DaemonSet instead of Deployment.                                                                                                                                                              | bool                                                                                                                           |
| containerSecurityContext                            | Security Context for a container.                                                                                                                                                                                                   | [*v1.SecurityContext](https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#securitycontext-v1-core) |
| configExistingSecretName                            | If the configuration is managed as secret outside the chart, using SealedSecret for example, provide the name of the secret here. If secretConfig is set to true, configExistingSecretName is ignored in favor of the config value. | string                                                                                                                         |
| secretConfig                                        | Store the configuration as a `Secret` instead of a `ConfigMap`, useful in case it contains sensitive data.                                                                                                                          | boolean                                                                                                                        |
| config                                              | Configuration of blackbox-exporter modules.                                                                                                                                                                                         | [Modules Configuration](https://github.com/prometheus/blackbox_exporter/blob/master/CONFIGURATION.md)                          |
| resources                                           | The resources that describe the compute resource requests and limits for single pods.                                                                                                                                               | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#resourcerequirements-v1-core)   |
| servicePort                                         | Port for blackbox-exporter service.                                                                                                                                                                                                 | integer                                                                                                                        |
| containerPort                                       | Only changes container port. Application port can be changed with extraArgs (--web.listen-address=:9115).                                                                                                                           | integer                                                                                                                        |
| createServiceAccount                                | Specifies whether a ServiceAccount should be created.                                                                                                                                                                               | boolean                                                                                                                        |
| extraArgs                                           | Additional arguments for blackbox-exporter container.                                                                                                                                                                               | list[string]                                                                                                                   |
| serviceMonitor.enabled                              | If true, a ServiceMonitor CRD is created for a prometheus operator.                                                                                                                                                                 | boolean                                                                                                                        |
| serviceMonitor.interval                             | Scrape interval.                                                                                                                                                                                                                    | string                                                                                                                         |
| serviceMonitor.scrapeTimeout                        | Scrape timeout.                                                                                                                                                                                                                     | string                                                                                                                         |
| serviceMonitor.scheme                               | HTTP scheme to use for scraping. Can be used with `tlsConfig` for example if using istio mTLS.                                                                                                                                      | string                                                                                                                         |
| serviceMonitor.defaults                             | DEPRECATED! Please use `Probe` instead. Default values that are used for all ServiceMonitors created by `targets`.                                                                                                                  | object                                                                                                                         |
| serviceMonitor.defaults.additionalMetricsRelabels   | DEPRECATED! Please use `Probe` instead. Default additional metrics relabels.                                                                                                                                                        | object                                                                                                                         |
| serviceMonitor.defaults.interval                    | DEPRECATED! Please use `Probe` instead. Default interval.                                                                                                                                                                           | string                                                                                                                         |
| serviceMonitor.defaults.scrapeTimeout               | DEPRECATED! Please use `Probe` instead. Default scrape timeout.                                                                                                                                                                     | string                                                                                                                         |
| serviceMonitor.defaults.module                      | DEPRECATED! Please use `Probe` instead. Default module name.                                                                                                                                                                        | string                                                                                                                         |
| serviceMonitor.targets                              | DEPRECATED! Please use `Probe` instead. Parameters for each targets that are created.                                                                                                                                               | list[object]                                                                                                                   |
| serviceMonitor.targets[N].name                      | DEPRECATED! Please use `Probe` instead. Human readable URL that appears in Prometheus / AlertManager                                                                                                                                | string                                                                                                                         |
| serviceMonitor.targets[N].url                       | DEPRECATED! Please use `Probe` instead. The URL that blackbox scrapes                                                                                                                                                               | string                                                                                                                         |
| serviceMonitor.targets[N].interval                  | DEPRECATED! Please use `Probe` instead. Scraping interval. Overrides value set in `defaults`                                                                                                                                        | string                                                                                                                         |
| serviceMonitor.targets[N].scrapeTimeout             | DEPRECATED! Please use `Probe` instead. Scrape timeout. Overrides value set in `defaults`                                                                                                                                           | string                                                                                                                         |
| serviceMonitor.targets[N].module                    | DEPRECATED! Please use `Probe` instead. Module used for scraping. Overrides value set in `defaults`                                                                                                                                 | string                                                                                                                         |
| serviceMonitor.targets[N].additionalMetricsRelabels | DEPRECATED! Please use `Probe` instead. Map of metric labels and values to add                                                                                                                                                      | object                                                                                                                         |
| grafanaDashboard                                    | Allows to create Grafana dashboard for blackbox-exporter.                                                                                                                                                                           | boolean                                                                                                                        |
| securityContext                                     | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                                  | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#podsecuritycontext-v1-core)      |
| tolerations                                         | Tolerations allow the pods to schedule onto nodes with matching taints.                                                                                                                                                             | []v1.Toleration                                                                                                                |
| nodeSelector                                        | Defines which nodes the pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\"                                                                                                                  | map[string]string                                                                                                              |
| annotations                                         | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value"              | map[string]string                                                                                                              |
| labels                                              | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                                    | map[string]string                                                                                                              |
| affinity                                            | If specified, the pod's scheduling constraints                                                                                                                                                                                      | *v1.Affinity                                                                                                                   |
| priorityClassName                                   | PriorityClassName assigned to the Pods to prevent them from evicting.                                                                                                                                                               | string                                                                                                                         |
<!-- markdownlint-enable line-length -->

Example:

```yaml
blackboxExporter:
  install: true
  name: blackbox-exporter
  image: prom/blackbox-exporter:v0.19.0
  containerSecurityContext:
    runAsUser: 1000
    runAsNonRoot: true
  configExistingSecretName: ""
  secretConfig: false
  config:
    modules:
      http_2xx:
        prober: http
        timeout: 5s
        tls_config:
          insecure_skip_verify: true
        http:
          valid_http_versions: [ "HTTP/1.1", "HTTP/2.0" ]
          no_follow_redirects: false
          preferred_ip_protocol: "ip4"
  resources:
    limits:
      memory: 300Mi
    requests:
      memory: 50Mi
  servicePort: 9115
  containerPort: 9115
  createServiceAccount: true
  extraArgs:
    - "--web.listen-address=:9115"
    - "--timeout-offset=0.5"
    - "--config.check=false"
    - "--history.limit=100"
    - "--web.external-url=http://example.com"
    - "--web.route-prefix=/example/path"
  serviceMonitor:
    enabled: true
    interval: 30s
    scrapeTimeout: 30s
    scheme: http
  grafanaDashboard: true
  securityContext:
    runAsUser: 2001
    fsGroup: 2001
  tolerations: []
  nodeSelector: {}
  labels:
    label.key: label-value
  annotations:
    annotation.key: annotation-value
  priorityClassName: priority-class
```

[Back to TOC](#table-of-content)

### cert-exporter

**Note**: Pay attention to the specifics of deploy with `cert-exporter`. For more information, refer to
the [Deploy with cert-exporter](#deploy-with-cert-exporter) section.

<!-- markdownlint-disable line-length -->
| Field                                   | Description                                                                                                                                                                                                            | Scheme                                                                                                                       |
| --------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| install                                 | Allows to enable or disable deploy cert-exporter.                                                                                                                                                                      | bool                                                                                                                         |
| name                                    | A name of the microservice to deploy with. This name is used as the name of the microservice deployment and in labels.                                                                                                 | string                                                                                                                       |
| setupSecurityContext                    | Allows to create PodSecurityPolicy or SecurityContextConstraints.                                                                                                                                                      | bool                                                                                                                         |
| setupGrafanaDashboard                   | Allows to create Grafana dashboard for cert-exporter.                                                                                                                                                                  | bool                                                                                                                         |
| setupAlertingRules                      | Allows to create Prometheus alerting rules for cert-exporter.                                                                                                                                                          | bool                                                                                                                         |
| additionalHostPathVolumes               | Allows to mount additional directories or files from the host file system to container. Should be used when the exporter needs to collect information from files from non-default paths.                               | list[object]                                                                                                                 |
| additionalHostPathVolumes[N].volumeName | Name of volume, must be unique. Reserved names (shouldn't be used): kube, openshift-origin, openshift-etcd, kubelet-pki, root-kube-config.                                                                             | string                                                                                                                       |
| additionalHostPathVolumes[N].volumePath | Path to the directory or file on the host system. The same path is used in the container.                                                                                                                              | string                                                                                                                       |
| certsInFiles                            | Settings for parsing certificates from host file system.                                                                                                                                                               | object                                                                                                                       |
| certsInFiles.enabled                    | Enables parsing certificates from host file system. If true, the part of the exporter will be deployed as a DaemonSet.                                                                                                 | bool                                                                                                                         |
| certsInFiles.defaultCerts               | Allows to check internal certificates from default paths which is default for Kubernetes or Openshift clusters.                                                                                                        | bool                                                                                                                         |
| certsInFiles.includeCerts               | Allows to check certificates with custom regex.                                                                                                                                                                        | string                                                                                                                       |
| certsInFiles.excludeCerts               | Allows to check only certificates that NOT match regex.                                                                                                                                                                | string                                                                                                                       |
| certsInKubeconfig                       | Settings for parsing certificates from Kubeconfig files on the host file system.                                                                                                                                       | object                                                                                                                       |
| certsInKubeconfig.enabled               | Enables parsing certificates from Kubeconfig. If true, the part of the exporter will be deployed as a DaemonSet.                                                                                                       | bool                                                                                                                         |
| certsInKubeconfig.defaultCerts          | Allows to check Kubeconfig files from default paths on Kubernetes or Openshift clusters.                                                                                                                               | bool                                                                                                                         |
| certsInKubeconfig.includeCerts          | Allows to check Kubeconfig files with custom regex.                                                                                                                                                                    | string                                                                                                                       |
| certsInKubeconfig.excludeCerts          | Allows to check only Kubeconfig files that NOT match regex.                                                                                                                                                            | string                                                                                                                       |
| certsInSecrets                          | Settings for parsing certificates from Kubernetes secrets.                                                                                                                                                             | object                                                                                                                       |
| certsInSecrets.enabled                  | Enables parsing certificates from Kubernetes secrets. If true, the part of the exporter will be deployed as a Deployment.                                                                                              | bool                                                                                                                         |
| certsInSecrets.includeCerts             | Allows to check fields in the secrets with keys that match regex.                                                                                                                                                      | string                                                                                                                       |
| certsInSecrets.excludeCerts             | Allows to check only fields in the secrets with keys NOT match regex.                                                                                                                                                  | string                                                                                                                       |
| certsInSecrets.annotationSelector       | Allows to match secrets by annotation.                                                                                                                                                                                 | string                                                                                                                       |
| certsInSecrets.labelSelector            | Allows to match secrets by label.                                                                                                                                                                                      | string                                                                                                                       |
| certsInSecrets.namespaces               | Allows to find secrets in the selected namespaces. Kubernetes comma-delimited list of namespaces to search for secrets. Empty string specifies that the exporter checks all available namespaces.                      | string                                                                                                                       |
| certsInSecrets.types                    | Allows to select only specific secret type. An empty list specifies that the exporter checks all available secrets.                                                                                                    | list[string]                                                                                                                 |
| certsInSecrets.kubeconfigPath           | Allows to specify path to kubeconfig file for getting access to secrets via kubectl. Only required if out-of-cluster installation.                                                                                     | string                                                                                                                       |
| pollingPeriod                           | Periodic interval in which to check certs. Format: [time.Duration](https://golang.org/pkg/time/#ParseDuration) from GoLang                                                                                             | string                                                                                                                       |
| image                                   | A Docker image to deploy the cert-exporter.                                                                                                                                                                            | string                                                                                                                       |
| serviceMonitor                          | Service monitor for pulling metrics.                                                                                                                                                                                   | object                                                                                                                       |
| serviceMonitor.install                  | Allows to install serviceMonitor.                                                                                                                                                                                      | bool                                                                                                                         |
| serviceMonitor.interval                 | Allow change metrics scrape interval.                                                                                                                                                                                  | string                                                                                                                       |
| servicePort                             | Port for cert-exporter service.                                                                                                                                                                                        | int                                                                                                                          |
| daemonset                               | Pod-specific settings for cert-exporter as a daemonset. Works if collecting from files and/or kubeconfig is enabled.                                                                                                   | object                                                                                                                       |
| daemonset.resources                     | The resources that describe the compute resource requests and limits for single pods. Affects daemonset pods.                                                                                                          | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#resourcerequirements-v1-core) |
| daemonset.extraArgs                     | Additional arguments for cert-exporter containers from daemonset.                                                                                                                                                      | list[string]                                                                                                                 |
| daemonset.securityContext               | SecurityContext holds pod-level security attributes.Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                      | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#podsecuritycontext-v1-core)    |
| daemonset.tolerations                   | Tolerations allow the pods to schedule onto nodes with matching taints.                                                                                                                                                | []v1.Toleration                                                                                                              |
| daemonset.nodeSelector                  | Defines which nodes the pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\"                                                                                                     | map[string]string                                                                                                            |
| daemonset.affinity                      | If specified, the pod's scheduling constraints                                                                                                                                                                         | *v1.Affinity                                                                                                                 |
| daemonset.annotations                   | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value" | map[string]string                                                                                                            |
| daemonset.labels                        | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                       | map[string]string                                                                                                            |
| daemonset.priorityClassName             | PriorityClassName assigned to the Pods to prevent them from evicting.                                                                                                                                                  | string                                                                                                                       |
| deployment                              | Settings for parsing certificates from Kubernetes secrets.                                                                                                                                                             | object                                                                                                                       |
| deployment.resources                    | The resources that describe the compute resource requests and limits for single pods. Affects deployment pods.                                                                                                         | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#resourcerequirements-v1-core) |
| deployment.extraArgs                    | Additional arguments for cert-exporter containers from deployment.                                                                                                                                                     | list[string]                                                                                                                 |
| deployment.securityContext              | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                     | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#podsecuritycontext-v1-core)    |
| deployment.tolerations                  | Tolerations allow the pods to schedule onto nodes with matching taints.                                                                                                                                                | []v1.Toleration                                                                                                              |
| deployment.nodeSelector                 | Defines which nodes the pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\"                                                                                                     | map[string]string                                                                                                            |
| deployment.affinity                     | If specified, the pod's scheduling constraints                                                                                                                                                                         | *v1.Affinity                                                                                                                 |
| deployment.annotations                  | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value" | map[string]string                                                                                                            |
| deployment.labels                       | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                       | map[string]string                                                                                                            |
| deployment.priorityClassName            | PriorityClassName assigned to the Pods to prevent them from evicting.                                                                                                                                                  | string                                                                                                                       |
<!-- markdownlint-enable line-length -->

Example of standard installation without root rights:

```yaml
certExporter:
  ...
  install: true
  name: cert-exporter
  setupSecurityContext: true
  setupGrafanaDashboard: true
  setupAlertingRules: true
  additionalHostPathVolumes: []
  certsInFiles:
    enabled: true
    defaultCerts: true
  certsInKubeconfig:
    enabled: false
  certsInSecrets:
    enabled: true
  pollingPeriod: 1h0m0s
  image: joeelliott/cert-exporter:v2.7.0
  serviceMonitor:
    install: true
    interval: 30s
  servicePort: 9219
  daemonset:
    resources:
      limits:
        cpu: 20m
        memory: 50Mi
      requests:
        cpu: 10m
        memory: 25Mi
    extraArgs: []
    securityContext:
      runAsUser: "0"
      fsGroup: "0"
    tolerations:
      - operator: "Exists"
    nodeSelector:
      node-role.kubernetes.io/worker: worker
    affinity: {}
    labels:
      label.key: label-value
    annotations:
      annotation.key: annotation-value
    priorityClassName: priority-class
  deployment:
    resources:
      limits:
        cpu: 20m
        memory: 150Mi
      requests:
        cpu: 10m
        memory: 50Mi
    extraArgs: []
    securityContext:
      runAsUser: 2000
      fsGroup: 2000
    tolerations: []
    nodeSelector:
      node-role.kubernetes.io/worker: worker
    affinity: {}
    labels:
      label.key: label-value
    annotations:
      annotation.key: annotation-value
    priorityClassName: priority-class
```

Example of installation with root rights:

```yaml
certExporter:
  ...
  install: true
  name: cert-exporter
  setupSecurityContext: true
  setupGrafanaDashboard: true
  setupAlertingRules: true
  certsInFiles:
    enabled: true
    defaultCerts: true
  certsInKubeconfig:
    enabled: true
    defaultCerts: true
  certsInSecrets:
    enabled: true
    kubeconfigPath: "/root/.kube/config"
  pollingPeriod: 1h0m0s
  image: joeelliott/cert-exporter:v2.7.0
  securityContext:
    runAsUser: "0"
```

Example of installation with non-default paths of files:

```yaml
certExporter:
  ...
  additionalHostPathVolumes:
            - volumeName: unique-volume-name-1
              volumePath: /path/to/certificates
                            - volumeName: unique-volume-name-2
                            volumePath: /path/to/kubeconfig
  certsInFiles:
    enabled: true
    defaultCerts: true
    includeCerts: "/path/to/certificates/*.{crt,cert}"
  certsInKubeconfig:
    enabled: true
    defaultCerts: false
    includeCerts: "/path/to/kubeconfig/**/*.conf"
  certsInSecrets:
    enabled: true
    kubeconfigPath: "/path/to/kubeconfig/config"
  ...
```

Example of installation with specifying parameters for parsing secrets:

```yaml
certExporter:
  ...
  certsInSecrets:
    enabled: true
    includeCerts: "*.crt"
    excludeCerts: "do-not-check.crt"
    annotationSelector: "some-annotation"
    labelSelector: "some-label"
    namespaces: "monitoring, namespace_1"
    types:
      - "Opaque"
      - "kubernetes.io/tls"
    kubeconfigPath: "/root/.kube/config"
  ...
  deployment:
    extraArgs:
      - '--secrets-include-glob=*.cer'
      - '--secrets-include-glob=*.pem'
  ...
```

[Back to TOC](#table-of-content)

### cloud-events-exporter

<!-- markdownlint-disable line-length -->

| Field                            | Description                                                                                                                                                                                                            | Scheme                                                                                                                       |
| -------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| install                          | Allows to enable or disable deploy cloud-events-exporter                                                                                                                                                               | bool                                                                                                                         |
| name                             | A deployment name for cloud-events-exporter                                                                                                                                                                            | string                                                                                                                       |
| image                            | A docker image to be used for deployment                                                                                                                                                                               | string                                                                                                                       |
| tolerations                      | Tolerations allow the pods to schedule onto nodes with matching taints                                                                                                                                                 | []v1.Toleration                                                                                                              |
| nodeSelector                     | Defines which nodes the pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\"                                                                                                     | map[string]string                                                                                                            |
| affinity                   | If specified, the pod's scheduling constraints                                                                                                                                                                         | *v1.Affinity                                                                                                                 |
| annotations                      | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value" | map[string]string                                                                                                            |
| labels                           | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                       | map[string]string                                                                                                            |
| resources                        | The resources that describe the compute resource requests and limits for single pods                                                                                                                                   | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#resourcerequirements-v1-core) |
| securityContext                  | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`                                                                                      | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#podsecuritycontext-v1-core)    |
| priorityClassName                | PriorityClassName assigned to the Pods to prevent them from evicting                                                                                                                                                   | string                                                                                                                       |
| namespaces                       | Namespace to watch for events. The parameter can be used multiple times.If parameter is not set events of all namespaces will be watched                                                                               | list[string]                                                                                                                 |
| port                             | Port to scrape cloud events metrics                                                                                                                                                                                    | string                                                                                                                       |
| filtering.match                  | Regular expressions for matching event by fields values. Only matching events will be recorded to metrics                                                                                                              | object                                                                                                                       |
| filtering.exclude                | Regular expressions for excluding event by fields values. Only events that do not match these expressions will be recorded to metrics                                                                                  | object                                                                                                                       |
| serviceMonitor.install           | Allow to enable or disable deploy service monitor                                                                                                                                                                      | boolean                                                                                                                      |
| serviceMonitor.interval          | Allow to change metrics scrape interval                                                                                                                                                                                | string                                                                                                                       |
| serviceMonitor.scrapeTimeout     | Allow to change metrics scrape timeout. Note that scrapeTimeout must be less the interval                                                                                                                              | string                                                                                                                       |
| serviceMonitor.metricRelabelings | Set metricRelabelings for the ServiceMonitor, use to apply to samples for ingestion                                                                                                                                    | []*promv1.RelabelConfig                                                                                                      |
| serviceMonitor.relabelings       | Set relabelings for the ServiceMonitor, use to apply to samples before scraping                                                                                                                                        | []*promv1.RelabelConfig                                                                                                      |
| serviceAccount.install           | Allow to disable create ServiceAccount during deploy                                                                                                                                                                   | bool                                                                                                                         |
| serviceAccount.annotations       | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value" | map[string]string                                                                                                            |
| serviceAccount.labels            | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                       | map[string]string                                                                                                            |
| rbac.createClusterRole           | Allow creating ClusterRole. If set to `false`, ClusterRole must be created manually. Default: `true`                                                                                                                   | bool                                                                                                                         |
| rbac.createClusterRoleBinding    | Allow creating ClusterRoleBinding. If set to `false`, ClusterRoleBinding must be created manually. Default: `true`                                                                                                     | bool                                                                                                                         |
| extraArgs                        | Additional arguments for cloud-events-exporter container. For example: "-workers=2"                                                                                                                                    | list[string]                                                                                                                 |
| extraEnvs                        | Additional environment variables for cloud-events-exporter container                                                                                                                                                   | object                                                                                                                       |

<!-- markdownlint-enable line-length -->

Example:

```yaml
cloudEventsExporter:
  install: true
  name: cloud-events-exporter
  image: product/prod.platform.logging_cloud-events-reader:latest
  tolerations:
    - key: "example-key"
      operator: "Exists"
      effect: "NoSchedule"
  nodeSelector:
    node-role.kubernetes.io/worker: worker
  labels:
    label.key: label-value
  annotations:
    annotation.key: annotation-value
  resources:
    limits:
      cpu: 200m
      memory: 200Mi
    requests:
      cpu: 100m
      memory: 100Mi
  securityContext:
    runAsUser: 2000
    fsGroup: 2000
  port: 9999
  metricsPath: "/metrics"
  namespaces:
    - mongo
    - postgres
  filtering:
    match:
      - type: "Warning"
        kind: "Pod|Deployment"
        reason: ""
        message: ""
        namespace: ""
        reportingController: ""
        reportingInstance: ""
    exclude:
      - type: "Normal"
        message: ".*image.*"
      - reason: "Completed|Pulled|Started"
  extraArgs:
    - "-workers=2"
  extraEnvs:
    ENV1: VALUE
```

[Back to TOC](#table-of-content)

### cloudwatch-exporter

cloudwatchExporter is a specification of the desired deployment of cloudwatch-exporter.

**Note**: Pay attention to settings that can be specified under the `cloudwatchExporter.config` parameter, such as,

```yaml
cloudwatchExporter:
  serviceMonitor:
    interval: 2m
  config:
    period_seconds: 120
    delay_seconds: 60
```

These parameters are very important to get the actual metrics from CloudWatch without a delay. The default values of
`delay_seconds` is `600 s`, and this parameter is used to avoid collecting data that has not fully converged.
`600 s` specifies that the exporter fetches data only after `600 s = 10 m`. To decrease the delay, it is recommended
to decrease the value to `60 s`. This time is enough for CloudWatch to collect the metrics and allows
cloudwatch_exporter to fetch them.

Refer to the official documentation of cloudwatch_exporter for full descriptions of all parameters at
[https://github.com/prometheus/cloudwatch_exporter#configuration](https://github.com/prometheus/cloudwatch_exporter#configuration).

<!-- markdownlint-disable line-length -->
| Field                                       | Description                                                                                                                                                                                                                                                                                                           | Scheme                                                                                                                       |
| ------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| install                                     | Allows to disable deploy cloudwatch-exporter.                                                                                                                                                                                                                                                                         | bool                                                                                                                         |
| replicas                                    | Number of created pods.                                                                                                                                                                                                                                                                                               | int                                                                                                                          |
| name                                        | A deployment name for cloudwatch-exporter                                                                                                                                                                                                                                                                             | string                                                                                                                       |
| image                                       | A docker image to use for cloudwatch-exporter deployment                                                                                                                                                                                                                                                              | string                                                                                                                       |
| imagePullPolicy                             | Image pull policy to use for cloudwatch-exporter deployment                                                                                                                                                                                                                                                           | string                                                                                                                       |
| command                                     | Allow override command to run docker container                                                                                                                                                                                                                                                                        | []string                                                                                                                     |
| resources                                   | The resources that describe the compute resource requests and limits for single pods.                                                                                                                                                                                                                                 | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#resourcerequirements-v1-core) |
| serviceAccount.install                      | Allow to disable create ServiceAccount during deploy                                                                                                                                                                                                                                                                  | bool                                                                                                                         |
| serviceAccount.name                         | Provide a name in place of cloudwatch-exporter for ServiceAccount                                                                                                                                                                                                                                                     | bool                                                                                                                         |
| serviceAccount.automountServiceAccountToken | Specifies whether to automount API credentials for the ServiceAccount.                                                                                                                                                                                                                                                | bool                                                                                                                         |
| rbac.createClusterRole                      | Allow creating ClusterRole. If set to `false`, ClusterRole must be created manually. Default: `true`                                                                                                                                                                                                                  | bool                                                                                                                         |
| rbac.createClusterRoleBinding               | Allow creating ClusterRoleBinding. If set to `false`, ClusterRoleBinding must be created manually. Default: `true`                                                                                                                                                                                                    | bool                                                                                                                         |
| nodeSelector                                | Defines which nodes the pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\"                                                                                                                                                                                                    | map[string]string                                                                                                            |
| annotations                                 | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value"                                                                                                | map[string]string                                                                                                            |
| labels                                      | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                                                                                                                      | map[string]string                                                                                                            |
| securityContext                             | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 65534, fsGroup: 65534 }`.                                                                                                                                                                                  | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#podsecuritycontext-v1-core)    |
| tolerations                                 | Tolerations allow the pods to schedule onto nodes with matching taints.                                                                                                                                                                                                                                               | []v1.Toleration                                                                                                              |
| affinity                                    | It specifies the pod's scheduling constraints. For more information, refer to [https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core)                                                          | *v1.Affinity                                                                                                                 |
| serviceMonitor                              | ServiceMonitor holds configuration attributes for cloudwatch-exporter.                                                                                                                                                                                                                                                | object                                                                                                                       |
| aws.aws_access_key_id                       | AWS Access Key ID for programmatic access. Do not specify the `aws_access_key_id` and `aws_secret_access_key` if you specified role or `.aws.secret.name` before.                                                                                                                                                     | string                                                                                                                       |
| aws.aws_secret_access_key                   | AWS Secret Access Key for programmatic access. Do not specify the `aws_access_key_id` and `aws_secret_access_key` if you specified role or `.aws.secret.name` before.                                                                                                                                                 | string                                                                                                                       |
| aws.secret.name                             | The name of a pre-created secret in which AWS credentials are stored. When set, `aws_access_key_id` is assumed to be in a field called `access_key`, `aws_secret_access_key` is assumed to be in a field called `secret_key`, and the session token, if it exists, is assumed to be in a field called security_token. | string                                                                                                                       |
| aws.secret.includesSessionToken             | Allow specify manually generated token and exporter will not try to use STS for get token. When set, token is assumed to be in a field called `security_token`.                                                                                                                                                       | bool                                                                                                                         |
| config                                      | Configuration is rendered with `tpl` function, therefore you can use any Helm variables and/or templates.                                                                                                                                                                                                             | string                                                                                                                       |
| priorityClassName                           | PriorityClassName assigned to the Pods to prevent them from evicting.                                                                                                                                                                                                                                                 | string                                                                                                                       |
<!-- markdownlint-enable line-length -->

Example:

```yaml
cloudwatchExporter:
  install: true
  replicas: 1
  name: cloudwatch-exporter
  image: prom/cloudwatch-exporter:cloudwatch_exporter-0.11.0
  imagePullPolicy: IfNotPresent
  command:
    - 'java'
    - '-Dhttp.proxyHost=proxy.example.com'
    - '-Dhttp.proxyPort=3128'
    - '-Dhttps.proxyHost=proxy.example.com'
    - '-Dhttps.proxyPort=3128'
    - '-jar'
    - '/cloudwatch_exporter.jar'
    - '9106'
    - '/config/config.yml'
  resources:
    limits:
      cpu: 200m
      memory: 256Mi
    requests:
      cpu: 100m
      memory: 128Mi
  serviceAccount:
    install: true
    name: cloudwatch-exporter
    annotations:
      eks.amazonaws.com/role-arn: arn:aws:iam::1234567890:role/prom-cloudwatch-exporter-oidc
    automountServiceAccountToken: true
  rbac:
    createClusterRole: true
    createClusterRoleBinding: true
  nodeSelector:
    node-role.kubernetes.io/worker: worker
  labels:
    label.key: label-value
  annotations:
    annotation.key: annotation-value
  securityContext:
    runAsUser: 65534  # run as nobody user instead of root
    fsGroup: 65534  # necessary to be able to read the EKS IAM token
  tolerations:
    - key: "example-key"
      operator: "Exists"
      effect: "NoSchedule"
  affinity: {}
  serviceMonitor:
    install: true
    interval: 5m
    telemetryPath: /metrics
    labels: {}
    timeout: 10s
    relabelings: []
    metricRelabelings:
      - sourceLabels: [dbinstance_identifier]
        action: replace
        replacement: mydbname
        targetLabel: dbname
  aws:
    secret:
      name: <secret_name>
      includesSessionToken: false
    aws_access_key_id: <access_key_id>
    aws_secret_access_key: <secret_access_key>
  config: |-
    # This is the default configuration for prometheus-cloudwatch-exporter
    region: us-east-1
    period_seconds: 240
    metrics:
    ### PostgreSQL
    - aws_namespace: AWS/RDS
      aws_metric_name: DatabaseConnections
      aws_dimensions: [DBInstanceIdentifier]
      aws_statistics: [Average]
    - aws_namespace: AWS/RDS
      aws_metric_name: DiskQueueDepth
      aws_dimensions: [DBInstanceIdentifier]
      aws_statistics: [Average]
```

[Back to TOC](#table-of-content)

### grafana

<!-- markdownlint-disable line-length -->
| Field                      | Description                                                                                                                                                                                                            | Scheme                                                                                                                                                |
| -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- |
| install                    | Allows to disable deploy Grafana. If Grafana was not deployed during the deployment using helm, it can be deployed using the change custom resource PlatformMonitoring.                                                | bool                                                                                                                                                  |
| paused                     | Set paused to reconciliation.                                                                                                                                                                                          | bool                                                                                                                                                  |
| image                      | A docker image to be used for the grafana deployment.                                                                                                                                                                  | string                                                                                                                                                |
| ingress                    | Allows to create Ingress for Grafana UI using monitoring-operator.                                                                                                                                                     | [v1.Ingress](#ingress)                                                                                                                                |
| resources                  | The resources that describe the compute resource requests and limits for single pods.                                                                                                                                  | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#resourcerequirements-v1-core)                          |
| securityContext            | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                     | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#podsecuritycontext-v1-core)                             |
| dataStorage                | Allows set a means to configure the grafana data storage.                                                                                                                                                              | [grafv1alpha1.GrafanaDataStorage](https://github.com/integr8ly/grafana-operator/blob/master/documentation/deploy_grafana.md#configuring-data-storage) |
| extraVars                  | Allows set extra system environment variables for grafana.                                                                                                                                                             | map[string]string                                                                                                                                     |
| grafanaHomeDashboard       | Allows set custom home dashboard for grafana. Dependence: `extraVars: GF_DASHBOARDS_DEFAULT_HOME_DASHBOARD_PATH: /etc/grafana-configmaps/grafana-home-dashboard/grafana-home-dashboard.json`                           | bool                                                                                                                                                  |
| backupDaemonDashboard      | Enables Backup Daemon Dashboard installation.                                                                                                                                                                          | bool                                                                                                                                                  |
| dashboardLabelSelector     | Allows to query over a set of resources according to labels.<br/>The result of matchLabels and matchExpressions are ANDed.<br/>An empty label selector matches all objects. A null label selector matches no objects.  | []*[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#labelselector-v1-meta)                                 |
| dashboardNamespaceSelector | Allows to query over a set of resources in namespaces that fits label selector.                                                                                                                                        | *[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#labelselector-v1-meta)                                   |
| podMonitor                 | Pod monitor for self monitoring.                                                                                                                                                                                       | *[Monitor](#monitor)                                                                                                                                  |
| config                     | Allows set configuration for grafana. The properties used to generate grafana.ini.                                                                                                                                     | [grafv1alpha1.GrafanaConfig](https://github.com/grafana-operator/grafana-operator/blob/master/documentation/deploy_grafana.md#config-reconciliation)  |
| affinity                                            | If specified, the pod's scheduling constraints                                                                                                                                                                                      | *v1.Affinity                                                                                                                   |
| annotations                | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value" | map[string]string                                                                                                                                     |
| labels                     | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                       | map[string]string                                                                                                                                     |
| priorityClassName          | PriorityClassName assigned to the Pods to prevent them from evicting                                                                                                                                                   | string                                                                                                                                                |
<!-- markdownlint-enable line-length -->

Example:

```yaml
grafana:
  install: true
  paused: false
  image: grafana/grafana:10.4.10
  ingress:
    ...
  resources:
    limits:
      cpu: 200m
      memory: 200Mi
    requests:
      cpu: 100m
      memory: 100Mi
  securityContext:
    runAsUser: 2000
    fsGroup: 2000
  config:
    auth:
      disable_login_form: false
      disable_signout_menu: true
    auth.anonymous:
      enabled: false
    log:
      level: warn
      mode: console
  extraVars:
    GF_DASHBOARDS_DEFAULT_HOME_DASHBOARD_PATH: /etc/grafana-configmaps/grafana-home-dashboard/grafana-home-dashboard.json
    GF_LIVE_ALLOWED_ORIGINS: "*"
    GF_FEATURE_TOGGLES_ENABLE: ngalert
  grafanaHomeDashboard: true
  backupDaemonDashboard: true
  dashboardLabelSelector:
    - matchLabels:
        app.kubernetes.io/component: monitoring
      matchExpressions:
        - key: openshift.io/cluster-monitoring
          operator: NotIn
          values: [ "true" ]
    - matchExpressions:
        - key: app.kubernetes.io/instance
          operator: Exists
        - key: app.kubernetes.io/version
          operator: Exists
  dashboardNamespaceSelector:
    matchLabels:
      label-key: label-value
    matchExpressions:
      - key: openshift.io/cluster-monitoring
        operator: NotIn
        values: [ "true" ]
      - key: kubernetes.io/metadata.name
        operator: In
        values:
          - monitoring
          - cassandra
  podMonitor:
    ...see example by link...
  dataStorage:
    labels:
      app: grafana
    annotations:
      app: grafana
    accessModes:
      - ReadWriteOnce
    size: 2Gi
    class: local-storage
  labels:
    label.key: label-value
  annotations:
    annotation.key: annotation-value
  priorityClassName: priority-class
```

[Back to TOC](#table-of-content)

#### grafana-operator

<!-- markdownlint-disable line-length -->
| Field              | Description                                                                                                                                                                                                            | Scheme                                                                                                                       |
| ------------------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| image              | A docker image to be used for the grafana-operator deployment.                                                                                                                                                         | string                                                                                                                       |
| paused             | Set paused to reconciliation.                                                                                                                                                                                          | bool                                                                                                                         |
| initContainerImage | A docker image to be used into initContainer in the Grafana deployment.                                                                                                                                                | string                                                                                                                       |
| resources          | The resources that describe the compute resource requests and limits for single Pods.                                                                                                                                  | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#resourcerequirements-v1-core) |
| securityContext    | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                     | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#podsecuritycontext-v1-core)    |
| podMonitor         | Pod monitor for self monitoring.                                                                                                                                                                                       | *[Monitor](#monitor)                                                                                                         |
| affinity                                            | If specified, the pod's scheduling constraints                                                                                                                                                                                      | *v1.Affinity                                                                                                                   |
| annotations        | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value" | map[string]string                                                                                                            |
| labels             | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                       | map[string]string                                                                                                            |
| priorityClassName  | PriorityClassName assigned to the Pods to prevent them from evicting                                                                                                                                                   | string                                                                                                                       |
<!-- markdownlint-enable line-length -->

Example:

```yaml
grafana:
  operator:
    image: integreatly/grafana-operator:latest
    paused: false
    initContainerImage: integreatly/grafana_plugins_init:latest
    resources:
      limits:
        cpu: 200m
        memory: 200Mi
      requests:
        cpu: 100m
        memory: 100Mi
    securityContext:
      runAsUser: 2000
      fsGroup: 2000
    podMonitor:
      ...see example by link...
    labels:
      label.key: label-value
    annotations:
      annotation.key: annotation-value
    priorityClassName: priority-class
```

[Back to TOC](#table-of-content)

#### grafana-image-renderer

<!-- markdownlint-disable line-length -->
**Warning**: The grafana-image-renderer requires two extra environment variables in Grafana:

* GF_RENDERING_SERVER_URL - `http://<image-renderer-address>:<port>/render`
* GF_RENDERING_CALLBACK_URL - `http://<grafana-adderss>:<port>/`

These variables have been set by default for local renderer and Grafana services. You don't have to override them. You
need change them in case if youare yousing external renderer.

**Warning**: Rendering images requires a lot of memory, mainly because Grafana creates browser instances in the
background for the actual rendering. If you are going to render a lot of panels it make sense allocate much more memory
than default value(developers of plugin suggest 16GB ram).

| Field             | Description                                                                                                                                                                                                                                                                     | Scheme                                                                                                                       |
| ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| install           | Allows to enable deploy Grafana image renderer.                                                                                                                                                                                                                                 | *bool                                                                                                                        |
| image             | A docker image to use for grafana-image-renderer deployment.                                                                                                                                                                                                                    | string                                                                                                                       |
| name              | This name is used as the name of the microservice deployment and in labels.                                                                                                                                                                                                     | []string                                                                                                                     |
| securityContext   | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                                                                              | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#podsecuritycontext-v1-core)    |
| annotations       | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value"                                                          | map[string]string                                                                                                            |
| labels            | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                                                                                | map[string]string                                                                                                            |
| resources         | The resources that describe the compute resource requests and limits for single Pods.                                                                                                                                                                                           | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#resourcerequirements-v1-core) |
| tolerations       | Tolerations allow the pods to schedule onto nodes with matching taints.                                                                                                                                                                                                         | []v1.Toleration                                                                                                              |
| port              | Port for grafana-image-renderer service.                                                                                                                                                                                                                                        | integer                                                                                                                      |
| extraEnvs         | Allow to set extra system environment variables for grafana-image-renderer. More information about env  variables in [Configuration guide](https://grafana.com/docs/grafana/v9.0/setup-grafana/image-rendering/?src=your_stories_page---------------------------#configuration) | map[string]string                                                                                                            |
| nodeSelector      | Defines which nodes the pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\"                                                                                                                                                              | map[string]string                                                                                                            |
| affinity                   | If specified, the pod's scheduling constraints                                                                                                                                                                         | *v1.Affinity                                                                                                                 |
| priorityClassName | PriorityClassName assigned to the Pods to prevent them from evicting.                                                                                                                                                                                                           | string                                                                                                                       |
<!-- markdownlint-enable line-length -->

Example:

```yaml
grafana:
  imageRenderer:
    install: true
    image: grafana/grafana-image-renderer:3.6.2
    name: grafana-image-renderer
    resources:
      limits:
        cpu: 300m
        memory: 500Mi
      requests:
        cpu: 150m
        memory: 250Mi
    securityContext:
      runAsUser: 2000
      fsGroup: 2000
    labels:
      label.key: label-value
    annotations:
      annotation.key: annotation-value
    port: 8282
    priorityClassName: priority-class
```

[Back to TOC](#table-of-content)

### grafana-dashboards

<!-- markdownlint-disable line-length -->
| Field   | Description                                       | Scheme   |
| ------- | ------------------------------------------------- | -------- |
| install | Allows to disable deploy all Grafana Dashboards.  | *bool    |
| list    | Allows to specify the list of Grafana dashboards. | []string |
<!-- markdownlint-enable line-length -->

Example:

```yaml
grafanaDashboards:
  install: true
  list:
    - alerts-overview
    - core-dns-dashboard
    - etcd-dashboard
    - govm-processes
    - home-dashboard
    - ingress-list-of-ingresses
    - ingress-nginx-controller
    - ingress-request-handling-performance
    - jvm-processes
    - kubernetes-cluster-overview
    - kubernetes-kubelet
    - kubernetes-apiserver
    - kubernetes-distribution-by-labels
    - kubernetes-namespace-resources
    - kubernetes-nodes-resources
    - kubernetes-pod-resources
    - kubernetes-pods-distribution-by-node
    - kubernetes-top-resources
    - node-details
    - overall-platform-health
    - prometheus-cardinality-explorer
    - prometheus-self-monitoring
```

[Back to TOC](#table-of-content)

### graphite-remote-adapter

<!-- markdownlint-disable line-length -->
| Field                    | Description                                                                                                                                                                                                                | Scheme                                                                                                                       |
| ------------------------ | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| install                  | Allows to enable or disable deploy graphite-remote-adapter.                                                                                                                                                                | *bool                                                                                                                        |
| name                     | A name of the microservice to deploy with. This name is used as the name of the microservice deployment and in labels.                                                                                                     | string                                                                                                                       |
| image                    | A Docker image to deploy the graphite-remote-adapter.                                                                                                                                                                      | string                                                                                                                       |
| replicas                 | Number of created pods.                                                                                                                                                                                                    | int                                                                                                                          |
| resources                | The resources that describe the compute resource requests and limits for single pods.                                                                                                                                      | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#resourcerequirements-v1-core) |
| securityContext          | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                         | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#podsecuritycontext-v1-core)    |
| tolerations              | Tolerations allow the pods to schedule onto nodes with matching taints.                                                                                                                                                    | []v1.Toleration                                                                                                              |
| nodeSelector             | Defines which nodes the pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\"                                                                                                         | map[string]string                                                                                                            |
| affinity                                            | If specified, the pod's scheduling constraints                                                                                                                                                                                      | *v1.Affinity                                                                                                                   |
| annotations              | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value"     | map[string]string                                                                                                            |
| labels                   | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                           | map[string]string                                                                                                            |
| servicemonitor           | ServiceMonitor holds configuration attributes for graphite-remote-adapter.                                                                                                                                                 | object                                                                                                                       |
| servicemonitor.install   | Allows to disable create ServiceMonitor CR for graphite-remote-adapter during the deployment.                                                                                                                              | bool                                                                                                                         |
| servicePort              | The port for graphite-remote-adapter service.                                                                                                                                                                              | int                                                                                                                          |
| writeCarbonAddress       | The `host:port` of the Graphite server to send samples to.                                                                                                                                                                 | string                                                                                                                       |
| readUrl                  | The URL of the remote Graphite Web server to send samples to.                                                                                                                                                              | string                                                                                                                       |
| defaultPrefix            | The prefix to prepends to all metrics exported to Graphite.                                                                                                                                                                | string                                                                                                                       |
| enableTags               | Enable using Graphite tags.                                                                                                                                                                                                | string                                                                                                                       |
| priorityClassName        | PriorityClassName assigned to the Pods to prevent them from evicting.                                                                                                                                                      | string                                                                                                                       |
| additionalGraphiteConfig | Additional Graphite Config.                                                                                                                                                                                                | object                                                                                                                       |
| graphite                 | Graphite configuration                                                                                                                                                                                                     | object                                                                                                                       |
| write                    | Write graphite configuration configuration                                                                                                                                                                                 | object                                                                                                                       |
| compress_type            | Write graphite configuration configuration compress type. Supported values: "", plain, lz4                                                                                                                                 | string                                                                                                                       |
| lz4_preferences          | Parameters for lz4 streaming compression                                                                                                                                                                                   | object                                                                                                                       |
| compression_level        | LZ4 streaming compression level. Min value 3, max 12, default 9                                                                                                                                                            | int                                                                                                                          |
| auto_flush               | LZ4 streaming compression always flush; reduces usage of internal buffers. Default - false                                                                                                                                 | bool                                                                                                                         |
| decompression_speed      | LZ4 streaming compression parser favors decompression speed vs compression ratio. Works for high compression modes (compression_level >= 10) only. Default - false                                                         | bool                                                                                                                         |
| frame                    | Parameters for lz4 streaming compression frame.                                                                                                                                                                            | object                                                                                                                       |
| block_size               | The larger the block size, the (slightly) better the compression ratio. Larger blocks also increase memory usage on both compression and decompression sides. Values: max64KB, max256KB, max1MB, max4MB. Default: max64KB. | string                                                                                                                       |
| block_mode               | Linked blocks sharply reduce inefficiencies when using small blocks, they compress better. Default - false, i.e. disabled.                                                                                                 | bool                                                                                                                         |
| content_checksum         | Add a 32-bit checksum of frame's decompressed data. Default - false, i.e. disabled.                                                                                                                                        | bool                                                                                                                         |
| block_checksum           | Each block followed by a checksum of block's compressed data. Default - false, i.e. disabled.                                                                                                                              | bool                                                                                                                         |
<!-- markdownlint-enable line-length -->

Example:

```yaml
graphite_remote_adapter:
  install: true
  name: graphite-remote-adapter
  image: product/prod.platform.system.graphite-remote-adapter:latest
  replicas: 1
  
    limits:
      cpu: 200m
      memory: 200Mi
    requests:
      cpu: 100m
      memory: 100Mi
  securityContext:
    runAsUser: 2000
    fsGroup: 2000
  tolerations:
    - key: "example-key"
      operator: "Exists"
      effect: "NoSchedule"
  nodeSelector:
    node-role.kubernetes.io/worker: worker
  labels:
    label.key: label-value
  annotations:
    annotation.key: annotation-value
  priorityClassName: priority-class
  servicemonitor:
    install: no
  servicePort: 9201
  writeCarbonAddress: localhost:9999
  readUrl: "http://guest:guest@localhost:8080"
  defaultPrefix: ""
  enableTags: true
  additionalGraphiteConfig:
    web:
      telemetry_path: "/metrics"
    write:
      timeout: 5m
    read:
      timeout: 5m
      delay: 1h
      ignore_error: true
    graphite:
      write:
        compress_type: lz4
        lz4_preferences:
          frame:
            block_size: max64KB
            block_mode: false
            content_checksum: false
            block_checksum: false
          compression_level: 9
          auto_flush: false
          decompression_speed: false
        carbon_transport: tcp
        carbon_reconnect_interval: 5m
        enable_paths_cache: true
        paths_cache_ttl: 4h
        paths_cache_purge_interval: 4h
        template_data:
          var1:
            foo: bar
          var2: foobar
```

[Back to TOC](#table-of-content)

### ingress

<!-- markdownlint-disable line-length -->
| Field         | Description                                                                                                                                                                                                                                                                                                                                                | Scheme            |
| ------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------- |
| install       | Indicates whether Ingress is needed to be installed. Marked as true if a host is specified.                                                                                                                                                                                                                                                                | bool              |
| host          | FQDN of the ingress host.                                                                                                                                                                                                                                                                                                                                  | string            |
| annotations   | Annotations is an unstructured key value map stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. They are not queryable and should be preserved when modifying objects. For more information, refer to [http://kubernetes.io/docs/user-guide/annotations](http://kubernetes.io/docs/user-guide/annotations) | map[string]string |
| labels        | Map of string keys and values that can be used to organize and categorize (scope and select) objects. May match selectors of replication controllers and services. For more information, refer to [http://kubernetes.io/docs/user-guide/labels](http://kubernetes.io/docs/user-guide/labels)                                                               | map[string]string |
| tlsSecretName | TlsSecretName allows to set secret name which will be used for TLS setting for the Ingress for specified host.                                                                                                                                                                                                                                             | string            |

<!-- markdownlint-enable line-length -->
You can also refer to [Ingress configuration](/docs/public/user-guides/monitoring-configuration/tls.md#ingress-configuration).

Example:

```yaml
<key>:
  ingress:
    install: true
    host: "my-awesome-ingress.cloud.org"
    annotations:
      nginx.ingress.kubernetes.io/proxy-http-version: 1.1
      nginx.ingress.kubernetes.io/ssl-redirect: true
    labels:
      app.kubernetes.io/part-of: monitoring
      app.kubernetes.io/managed-by: monitoring-operator
    tlsSecretName: prometheus-tls-secret
```

[Back to TOC](#table-of-content)

### integration

<!-- markdownlint-disable line-length -->
| Field       | Description                                                                                                                                                                                                                                                          | Scheme                       |
| ----------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------- |
| stackdriver | Set up Platform Monitoring integration with Google Cloud Operations (GCO). Integration schema: Send metrics from Prometheus to GCO by deploying 'stackdriver-prometheus-sidecar' container as sidecar to Prometheus pod. Allows specify filters for metrics to send. | *[stackdriver](#stackdriver) |
| jaeger      | Set up Platform Monitoring integration with Jaeger.                                                                                                                                                                                                                  | *[jaeger](#jaeger)           |
| clickHouse  | Create Grafana DateSource for ClickHouse Performance Dashboard.                                                                                                                                                                                                      | *[clickHouse](#clickhouse)   |
<!-- markdownlint-enable line-length -->

Example:

```yaml
integration:
  stackdriver:
    ...see example by link...
  jaeger:
    ...see example by link...
  clickHouse:
    ...see example by link...
```

[Back to TOC](#table-of-content)

#### jaeger

<!-- markdownlint-disable line-length -->
| Field                   | Description                                                                                                     | Scheme |
| ----------------------- | --------------------------------------------------------------------------------------------------------------- | ------ |
| createGrafanaDataSource | If true, looking for Jaeger Service in all namespaces and add Grafana DataSource for it service if it is found. | bool   |
<!-- markdownlint-enable line-length -->

Example:

```yaml
integration:
  ...
  jaeger:
    createGrafanaDataSource: true
  ...
```

#### clickHouse

<!-- markdownlint-disable line-length -->
| Field                   | Description                                                                                                         | Scheme |
| ----------------------- | ------------------------------------------------------------------------------------------------------------------- | ------ |
| createGrafanaDataSource | If true, looking for ClickHouse Service in all namespaces and add Grafana DataSource for it service if it is found. | bool   |
<!-- markdownlint-enable line-length -->

Example:

```yaml
integration:
  ...
  clickHouse:
    createGrafanaDataSource: true
  ...
```

### json-exporter

Json exporter allow to execute probe to URL (static or for Ingress) to scrape remote JSON by JSONPath.

<!-- markdownlint-disable line-length -->
| Field                                               | Description                                                                                                                                                                                                                                                                                                                                                                                                                     | Scheme                                                                                                                         |
| --------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------ |
| install                                             | Allows to enable or disable deploy json-exporter.                                                                                                                                                                                                                                                                                                                                                                               | bool                                                                                                                           |
| name                                                | A name of the microservice to deploy with. This name is used as the name of the microservice deployment and in labels.                                                                                                                                                                                                                                                                                                          | string                                                                                                                         |
| image                                               | A Docker image to deploy the json-exporter.                                                                                                                                                                                                                                                                                                                                                                                     | string                                                                                                                         |
| imagePullPolicy                                     | Image pull policy to use for promitor-agent-scraper deployment.                                                                                                                                                                                                                                                                                                                                                                 | string                                                                                                                         |
| imagePullSecrets                                    | Reference to one or more secrets to be used when pulling images, refer to [https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry](https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry)                                                                                                                                                                          | string                                                                                                                         |
| containerPort                                       | Only changes container port.                                                                                                                                                                                                                                                                                                                                                                                                    | integer                                                                                                                        |
| replicaCount                                        | Number of created pods.                                                                                                                                                                                                                                                                                                                                                                                                         | integer                                                                                                                        |
| serviceAccount.install                              | Allow to disable create ServiceAccount during deploy                                                                                                                                                                                                                                                                                                                                                                            | bool                                                                                                                           |
| serviceAccount.name                                 | Provide a name in place of json-exporter for ServiceAccount                                                                                                                                                                                                                                                                                                                                                                     | bool                                                                                                                           |
| serviceAccount.annotations                          | Map of string keys and values stored with a resource for the ServiceAccount.                                                                                                                                                                                                                                                                                                                                                    | map[string]string                                                                                                              |
| securityContext                                     | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                                                                                                                                                                                                                              | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#podsecuritycontext-v1-core)      |
| containerSecurityContext                            | Security Context for a container.                                                                                                                                                                                                                                                                                                                                                                                               | [*v1.SecurityContext](https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#securitycontext-v1-core) |
| service.type                                        | Type of json-exporter service.                                                                                                                                                                                                                                                                                                                                                                                                  | string                                                                                                                         |
| service.port                                        | Port of json-exporter which use in service.                                                                                                                                                                                                                                                                                                                                                                                     | integer                                                                                                                        |
| service.name                                        | Port name of json-exporter which use in service.                                                                                                                                                                                                                                                                                                                                                                                | string                                                                                                                         |
| service.labels                                      | Labels set which will create in service.                                                                                                                                                                                                                                                                                                                                                                                        | map[string]string                                                                                                              |
| serviceMonitor.enabled                              | If true, a ServiceMonitor CRD is created for a prometheus operator.                                                                                                                                                                                                                                                                                                                                                             | boolean                                                                                                                        |
| serviceMonitor.interval                             | Scrape interval.                                                                                                                                                                                                                                                                                                                                                                                                                | string                                                                                                                         |
| serviceMonitor.scrapeTimeout                        | Scrape timeout.                                                                                                                                                                                                                                                                                                                                                                                                                 | string                                                                                                                         |
| serviceMonitor.scheme                               | HTTP scheme to use for scraping. Can be used with `tlsConfig` for example if using istio mTLS.                                                                                                                                                                                                                                                                                                                                  | string                                                                                                                         |
| serviceMonitor.defaults                             | Default values that are used for all ServiceMonitors created by `targets`.                                                                                                                                                                                                                                                                                                                                                      | object                                                                                                                         |
| serviceMonitor.defaults.additionalMetricsRelabels   | Default additional metrics relabels.                                                                                                                                                                                                                                                                                                                                                                                            | object                                                                                                                         |
| serviceMonitor.defaults.interval                    | Default interval.                                                                                                                                                                                                                                                                                                                                                                                                               | string                                                                                                                         |
| serviceMonitor.defaults.scrapeTimeout               | Default scrape timeout.                                                                                                                                                                                                                                                                                                                                                                                                         | string                                                                                                                         |
| serviceMonitor.defaults.module                      | Default module name.                                                                                                                                                                                                                                                                                                                                                                                                            | string                                                                                                                         |
| serviceMonitor.targets                              | Parameters for each targets that are created.                                                                                                                                                                                                                                                                                                                                                                                   | list[object]                                                                                                                   |
| serviceMonitor.targets[N].name                      | Human readable URL that appears in Prometheus / AlertManager                                                                                                                                                                                                                                                                                                                                                                    | string                                                                                                                         |
| serviceMonitor.targets[N].url                       | The URL that blackbox scrapes                                                                                                                                                                                                                                                                                                                                                                                                   | string                                                                                                                         |
| serviceMonitor.targets[N].interval                  | Scraping interval. Overrides value set in `defaults`                                                                                                                                                                                                                                                                                                                                                                            | string                                                                                                                         |
| serviceMonitor.targets[N].scrapeTimeout             | Scrape timeout. Overrides value set in `defaults`                                                                                                                                                                                                                                                                                                                                                                               | string                                                                                                                         |
| serviceMonitor.targets[N].module                    | Module used for scraping. Overrides value set in `defaults`                                                                                                                                                                                                                                                                                                                                                                     | string                                                                                                                         |
| serviceMonitor.targets[N].additionalMetricsRelabels | Map of metric labels and values to add                                                                                                                                                                                                                                                                                                                                                                                          | object                                                                                                                         |
| resources                                           | The resources that describe the compute resource requests and limits for single pods.                                                                                                                                                                                                                                                                                                                                           | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#resourcerequirements-v1-core)   |
| nodeSelector                                        | Defines which nodes the pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\"                                                                                                                                                                                                                                                                                                              | map[string]string                                                                                                              |
| tolerations                                         | Tolerations allow the pods to schedule onto nodes with matching taints.                                                                                                                                                                                                                                                                                                                                                         | []v1.Toleration                                                                                                                |
| affinity                                            | It specifies the pod's scheduling constraints. For more information, refer to [https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core)                                                                                                                                                                    | *v1.Affinity                                                                                                                   |
| config                                              | Configuration of json-exporter modules.                                                                                                                                                                                                                                                                                                                                                                                         | [Modules Configuration](https://github.com/prometheus-community/json_exporter/tree/master/examples)                            |
| additionalVolumes                                   | Volumes allows configuration of additional volumes. Volumes specified are appended to other volumes that are generated as a result of StorageSpec objects. For more information, refer to [https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#volume-v1-core](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#volume-v1-core)                                                            | []v1.Volume                                                                                                                    |
| additionalVolumeMounts                              | VolumeMounts allows configuration of additional VolumeMounts. VolumeMounts specified are appended to other VolumeMounts in the prometheus container, that are generated as a result of StorageSpec objects. For more informaiton, refer to [https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#volumemount-v1-core](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#volumemount-v1-core) | []v1.VolumeMount                                                                                                               |
| annotations                                         | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value"                                                                                                                                                                                                          | map[string]string                                                                                                              |
| labels                                              | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                                                                                                                                                                                                                                | map[string]string                                                                                                              |
| priorityClassName                                   | PriorityClassName assigned to the Pods to prevent them from evicting.                                                                                                                                                                                                                                                                                                                                                           | string                                                                                                                         |
<!-- markdownlint-enable line-length -->

Example:

```yaml
jsonExporter:
  install: true
  config: |
    modules:
      tsdbstatus:
        metrics:
          - name: "tsdb_series_count_by_metric"
            type: object
            help: "General vmsingle tsdb status of the series count by metric name"
            path: '{.data.seriesCountByMetricName[*]}'
            labels:
              name: '{.name}'
            values:
              value: '{.value}'
          - name: "tsdb_series_count_by_label"
            type: object
            help: "General vmsingle tsdb status of the series count by label name"
            path: '{.data.seriesCountByLabelName[*]}'
            labels:
              name: '{.name}'
            values:
              value: '{.value}'
          - name: "tsdb_series_count_by_label_value_pair"
            type: object
            help: "General vmsingle tsdb status of the series count by label value pair"
            path: '{.data.seriesCountByLabelValuePair[*]}'
            labels:
              name: '{.name}'
            values:
              value: '{.value}'
          - name: "tsdb_label_value_count_by_label_name"
            type: object
            help: "General vmsingle tsdb status of the label value count by label name"
            path: '{.data.labelValueCountByLabelName[*]}'
            labels:
              name: '{.name}'
            values:
              value: '{.value}'
  serviceMonitor:
    targets:
      - name: vmsingle-tsdb-status
        url: http://vmsingle-k8s.monitoring.svc:8429/api/v1/status/tsdb
        labels: {}
        interval: 60s
        scrapeTimeout: 60s
        additionalMetricsRelabels: {}
        module: tsdbstatus    
  securityContext:
    runAsUser: 2000
    fsGroup: 2000
```

[Back to TOC](#table-of-content)

### kube-state-metrics

<!-- markdownlint-disable line-length -->
| Field                 | Description                                                                                                                                                                                                                                                     | Scheme                                                                                                                       |
| --------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| install               | Allows to disable create kube-state-metrics during the deployment.                                                                                                                                                                                              | bool                                                                                                                         |
| paused                | Set paused to reconciliation.                                                                                                                                                                                                                                   | bool                                                                                                                         |
| image                 | A docker image to be used for the kube-state-metrics deployment.                                                                                                                                                                                                | string                                                                                                                       |
| resources             | The resources that describe the compute resource requests and limits for single pods.                                                                                                                                                                           | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#resourcerequirements-v1-core) |
| securityContext       | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                                                              | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#podsecuritycontext-v1-core)    |
| affinity                                            | If specified, the pod's scheduling constraints                                                                                                                                                                                      | *v1.Affinity                                                                                                                   |
| annotations           | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value"                                          | map[string]string                                                                                                            |
| labels                | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                                                                | map[string]string                                                                                                            |
| priorityClassName     | PriorityClassName assigned to the Pods to prevent them from evicting.                                                                                                                                                                                           | string                                                                                                                       |
| namespaces            | Comma separated list of namespaces to monitor in non-privileged mode. This parameter is unnecessary if the kube-state-metrics has ClusterRole on all namespaces                                                                                                 | string                                                                                                                       |
| scrapeResources       | Comma-separated list of Resources to be enabled. Empty means that metrics will collect for all type of resources                                                                                                                                                | string                                                                                                                       |
| metricLabelsAllowlist | Comma-separated list of additional Kubernetes label keys that will be used in the resource labels metric. Default value is `nodes=[*],pods=[*],namespaces=[*],deployments=[*],statefulsets=[*],daemonsets=[*],cronjobs=[*],jobs=[*],ingresses=[*],services=[*]` | string                                                                                                                       |
<!-- markdownlint-enable line-length -->

Example:

```yaml
kubeStateMetrics:
  install: true
  paused: false
  image: coreos/kube-state-metrics:latest
  resources:
    limits:
      cpu: 200m
      memory: 200Mi
    requests:
      cpu: 100m
      memory: 100Mi
  securityContext:
    runAsUser: 2000
    fsGroup: 2000
  labels:
    label.key: label-value
  annotations:
    annotation.key: annotation-value
  priorityClassName: priority-class

  namespaces: monitoring,logging,tracing
  scrapeResources: configmaps,cronjobs,daemonsets,deployments,endpoints,jobs,limitranges,persistentvolumeclaims,poddisruptionbudgets,namespaces,nodes,pods,persistentvolumes,replicasets,replicationcontrollers,resourcequotas,services,statefulsets
  metricLabelsAllowlist: nodes=[*],pods=[*],namespaces=[*],deployments=[*],statefulsets=[*],daemonsets=[*],cronjobs=[*],jobs=[*],ingresses=[*],services=[*]
```

[Back to TOC](#table-of-content)

### kubernetes-monitors

<!-- markdownlint-disable line-length -->
| Field                                        | Description | Scheme               |
| -------------------------------------------- | ----------- | -------------------- |
| apiserverServiceMonitor                      |             | *[Monitor](#monitor) |
| coreDnsServiceMonitor                        |             | *[Monitor](#monitor) |
| etcdServiceMonitor                           |             | *[Monitor](#monitor) |
| kubeControllerManagerServiceMonitor          |             | *[Monitor](#monitor) |
| kubeletServiceMonitor                        |             | *[Monitor](#monitor) |
| kubeSchedulerServiceMonitor                  |             | *[Monitor](#monitor) |
| nginxIngressPodMonitor                       |             | *[Monitor](#monitor) |
| openshiftApiserverServiceMonitor             |             | *[Monitor](#monitor) |
| openshiftApiserverOperatorServiceMonitor     |             | *[Monitor](#monitor) |
| openshiftClusterVersionOperatorServiceMonito |             | *[Monitor](#monitor) |
| openshiftStatemetricsServiceMonitor          |             | *[Monitor](#monitor) |
| openshiftHAProxyServiceMonitor               |             | *[Monitor](#monitor) |
<!-- markdownlint-enable line-length -->

```yaml
kubernetesMonitors:
  etcdServiceMonitor:
    install: true
    interval: 30s
    scrapeTimeout: 10s
    metricRelabelings: []
    relabelings: []
  apiserverServiceMonitor:
    install: true
    interval: 30s
    scrapeTimeout: 10s
    metricRelabelings:
      - sourceLabels: ['__name__']
        regex: 'apiserver_request_(.*)|workqueue_(.*)|etcd_(.*)|process_(.*)|go_(.*)'
        action: keep
    relabelings:
      - sourceLabels: ['job']
        targetLabel: 'job'
        replacement: 'kube-apiserver'
  kubeControllerManagerServiceMonitor:
    install: false
    interval: 30s
    scrapeTimeout: 10s
    metricRelabelings:
      - action: drop
        regex: etcd_(debugging|disk|request|server).*
        sourceLabels: ['__name__']
    relabelings:
      - action: replace
        regex: (.+)(?::\d+)
        replacement: $1:10252
        sourceLabels: ['__address__']
        targetLabel: __address__
  kubeSchedulerServiceMonitor:
    install: false
    interval: 30s
    scrapeTimeout: 10s
    metricRelabelings: []
    relabelings:
      - action: replace
        regex: (.+)(?::\d+)
        replacement: $1:10251
        sourceLabels: ['__address__']
        targetLabel: __address__
  kubeletServiceMonitor:
    install: true
    interval: 30s
    scrapeTimeout: 10s
    metricRelabelings:
      - sourceLabels: ['pod_name']
        targetLabel: pod
        regex: (.+)
      - sourceLabels: ['container_name']
        targetLabel: container
        regex: (.+)
      - action: labeldrop
        regex: pod_name
      - action: labeldrop
        regex: container_name
      - regex: 'kubelet_running_pods'
        replacement: 'kubelet_running_pod_count'
        sourceLabels: ['__name__']
        targetLabel: __name__
      - regex: 'kubelet_running_containers'
        replacement: 'kubelet_running_container_count'
        sourceLabels: ['__name__']
        targetLabel: __name__
    relabelings: []
  coreDnsServiceMonitor:
    install: true
    interval: 30s
    scrapeTimeout: 10s
    metricRelabelings: []
    relabelings:
      - action: replace
        regex: (.*)
        replacement: $1
        sourceLabels: ['__meta_kubernetes_pod_node_name']
        targetLabel: node
  nginxIngressPodMonitor:
    install: true
    interval: 30s
    scrapeTimeout: 10s
    metricRelabelings: []
    relabelings:
      - action: replace
        regex: (.*)
        replacement: $1
        sourceLabels: ['__meta_kubernetes_pod_node_name']
        targetLabel: node
  openshiftApiserverServiceMonitor:
    install: true
    interval: 30s
    scrapeTimeout: 10s
    metricRelabelings:
      - action: drop
        regex: etcd_(debugging|disk|server).*
        sourceLabels:
          - __name__
      - action: drop
        regex: apiserver_admission_controller_admission_latencies_seconds_.*
        sourceLabels:
          - __name__
      - action: drop
        regex: apiserver_admission_step_admission_latencies_seconds_.*
        sourceLabels:
          - __name__
      - action: drop
        regex: >-
          apiserver_request_duration_seconds_bucket;(0.15|0.25|0.3|0.35|0.4|0.45|0.6|0.7|0.8|0.9|1.25|1.5|1.75|2.5|3|3.5|4.5|6|7|8|9|15|25|30|50)
        sourceLabels:
          - __name__
          - le
    relabelings:
      - action: replace
        replacement: openshift-apiserver
        targetLabel: apiserver
  openshiftApiserverOperatorServiceMonitor:
    install: true
    interval: 30s
    scrapeTimeout: 10s
    metricRelabelings: []
    relabelings: []
  openshiftClusterVersionOperatorServiceMonitor:
    install: true
    interval: 30s
    scrapeTimeout: 10s
    metricRelabelings: []
    relabelings: []
  openshiftStatemetricsServiceMonitor:
    install: true
    interval: 30s
    scrapeTimeout: 10s
    metricRelabelings: []
    relabelings: []
  openshiftHAProxyServiceMonitor:
    install: true
    interval: 30s
    scrapeTimeout: 10s
    metricRelabelings: []
    relabelings: []
```

[Back to TOC](#table-of-content)

### monitor

<!-- markdownlint-disable line-length -->
| Field             | Description | Scheme                  |
| ----------------- | ----------- | ----------------------- |
| install           |             | *bool                   |
| interval          |             | string                  |
| scrapeTimeout     |             | string                  |
| relabelings       |             | []*promv1.RelabelConfig |
| metricRelabelings |             | []*promv1.RelabelConfig |
<!-- markdownlint-enable line-length -->

**Note**: Pay attention that to enable metrics collection from ingress-nginx, you need open port 10254, for example
with `prometheus` name.

```yaml

kind: DaemonSet
apiVersion: apps/v1
metadata:
  name: ingress-nginx-controller
...
- name: prometheus
  ...
  containerPort: 10254
  protocol: TCP
...
```

Example:

```yaml
[serviceMonitor|podMonitor]:
install: true
interval: 30s
metricRelabelings:
  - sourceLabels:
      - pod_name
    targetLabel: pod
    regex: (.+)
  - sourceLabels:
      - container_name
    targetLabel: container
    regex: (.+)
  - action: labeldrop
    regex: pod_name
  - action: labeldrop
    regex: container_name
  - regex: 'kubelet_running_pods'
    replacement: 'kubelet_running_pod_count'
    sourceLabels: [__name__]
    targetLabel: __name__
  - regex: 'kubelet_running_containers'
    replacement: 'kubelet_running_container_count'
    sourceLabels: [__name__]
    targetLabel: __name__
relabelings: []
```

[Back to TOC](#table-of-content)

### monitoring-operator

<!-- markdownlint-disable line-length -->
| Field                  | Description                                                                                                                                                                                                            | Scheme                                                                                                                       |
| ---------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| image                  | A docker image to be used for the monitoring-operator deployment.                                                                                                                                                      | string                                                                                                                       |
| resources              | The resources that describe the compute resource requests and limits for single pods.                                                                                                                                  | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#resourcerequirements-v1-core) |
| securityContext        | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                     | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#podsecuritycontext-v1-core)    |
| serviceAccount.install | Allow to disable create ServiceAccount during deploy                                                                                                                                                                   | bool                                                                                                                         |
| serviceAccount.name    | Provide a name in place of monitoring-operator for ServiceAccount                                                                                                                                                      | bool                                                                                                                         |
| affinity                                            | If specified, the pod's scheduling constraints                                                                                                                                                                                      | *v1.Affinity                                                                                                                   |
| annotations            | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value" | map[string]string                                                                                                            |
| labels                 | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                       | map[string]string                                                                                                            |
| priorityClassName      | PriorityClassName assigned to the Pods to prevent them from evicting.                                                                                                                                                  | string                                                                                                                       |
<!-- markdownlint-enable line-length -->

Example:

```yaml
monitoringOperator:
  image: prod.platform.system.monitoring_monitoring-operator:latest
  resources:
    limits:
      cpu: 200m
      memory: 200Mi
    requests:
      cpu: 100m
      memory: 100Mi
  securityContext:
    runAsUser: 2000
    fsGroup: 2000
  serviceAccount:
    install: true
    name: "monitoring-operator"
  labels:
    label.key: label-value
  annotations:
    annotation.key: annotation-value
  priorityClassName: priority-class
```

[Back to TOC](#table-of-content)

### network-latency-exporter

**Important**: In some cases, TCP and UDP probes can show large packet loss not due to network problems, but due to system
settings. The `net.ipv4.icmp_ratelimit` kernel parameter with a default value of `1000` prevents a node from sending "
ICMP time in-transit exceeded" reply immediately, so `mtr` utility which is used in `network-latency-exporter` shows
losses.

Address **Important**: Please note that memory consumption by network-latency-exporter depends on number of nodes where
exporter will be installed and depends on number of using protocols.

The main formula is:

```bash
10Mi + (5Mi * (<number_of_nodes_in_cluster> - 1) * <number_of_protocols>)
```

_To avoid this behavior, set kernel variable `net.ipv4.icmp_ratelimit = 0` on all nodes._

<!-- markdownlint-disable line-length -->
| Field                         | Description                                                                                                                                                                                                               | Scheme                                                                                                                       |
| ----------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| install                       | Allows to disable create network-latency-exporter during the deployment.                                                                                                                                                  | bool                                                                                                                         |
| name                          | A name of the microservice to deploy with. This name is used as the name of the microservice deployment and in labels.                                                                                                    | string                                                                                                                       |
| image                         | A docker image to be used for the network-latency-exporter deployment.                                                                                                                                                    | string                                                                                                                       |
| rbac.createClusterRole        | Allow creating ClusterRole. If set to `false`, ClusterRole must be created manually.                                                                                                                                      | bool                                                                                                                         |
| rbac.createClusterRoleBinding | Allow creating ClusterRoleBinding. If set to `false`, ClusterRoleBinding must be created manually.                                                                                                                        | bool                                                                                                                         |
| rbac.setupSecurityContext     | Allow creating PodSecurityPolicy or SecurityContextConstraints. If set to `false`, PodSecurityPolicy / SecurityContextConstraints must be created manually.                                                               | bool                                                                                                                         |
| rbac.privileged               | If `true`, set parameters in PSP or SCC for ability to running container in the privileged mode and set `privileged: true` to the security context in the exporter's container.                                           | bool                                                                                                                         |
| createGrafanaDashboards       | Allow creating Grafana Dashboards `Network Latency Overview` and `Network Latency Details`.                                                                                                                               | bool                                                                                                                         |
| serviceAccount.create         | Allow creating ServiceAccount. If set to `false`, ServiceAccount must be created manually.                                                                                                                                | bool                                                                                                                         |
| serviceAccount.name           | Provide a name in place of network-latency-exporter for ServiceAccount.                                                                                                                                                   | string                                                                                                                       |
| extraArgs                     | The level of application logging may be set using extraArgs. Set `INFO` or `DEBUG` to see informative messages about metrics collection.                                                                                  | object                                                                                                                       |
| discoverEnable                | Allow enabling/disabling script for discovering nodes IP.                                                                                                                                                                 | bool                                                                                                                         |
| targets                       | List nodes with name and IP address. For example, [{"name": "node1", "ipAddress": "1.2.3.1"}, {"name": "node2", "ipAddress": "1.2.3.3"}]                                                                                  | object                                                                                                                       |
| requestTimeout                | The response time for each packet sent which the application waits response, in seconds.                                                                                                                                  | integer                                                                                                                      |
| timeout                       | The metrics collection timeout. Can be calculated as TIMEOUT = 10s + (REQUEST_TIMEOUT * PACKETS_NUM * <NUMBER_OF_PROTOCOLS>)                                                                                              | string                                                                                                                       |
| packetsNum                    | The number of packets to send per probe.                                                                                                                                                                                  | integer                                                                                                                      |
| packetSize                    | The size of packet to sent in bytes.                                                                                                                                                                                      | integer                                                                                                                      |
| checkTarget                   | The comma-separated list of network protocols and ports (separated by ':') through which packets are sent. The supported protocols are: `UDP`, `TCP`, `ICMP`. If no port is specified for the protocol, port `1` is used. | string                                                                                                                       |
| serviceMonitor.enabled        | If true, a ServiceMonitor is created for a `prometheus-operator`.                                                                                                                                                         | string                                                                                                                       |
| serviceMonitor.interval       | Scraping interval for Prometheus.                                                                                                                                                                                         | bool                                                                                                                         |
| additionalLabels              | Allows specifying custom labels for DaemonSet of `network-latency-exporter`.                                                                                                                                              | object                                                                                                                       |
| setupSecurityContext          | Allows to create PodSecurityPolicy or SecurityContextConstraints.                                                                                                                                                         | string                                                                                                                       |
| port                          | The port for node-exporter daemonset and service.                                                                                                                                                                         | int                                                                                                                          |
| resources                     | The resources that describe the compute resource requests and limits for single pods.                                                                                                                                     | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#resourcerequirements-v1-core) |
| securityContext               | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 0, fsGroup: 2000 }`.                                                                                           | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#podsecuritycontext-v1-core)    |
| tolerations                   | Tolerations allow the pods to schedule onto nodes with matching taints.                                                                                                                                                   | []v1.Toleration                                                                                                              |
| nodeSelector                  | Defines which nodes the pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\"                                                                                                        | map[string]string                                                                                                            |
| affinity                      | If specified, the pod's scheduling constraints                                                                                                                                                                            | *v1.Affinity                                                                                                                 |
| annotations                   | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value"    | map[string]string                                                                                                            |
| labels                        | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                          | map[string]string                                                                                                            |
| priorityClassName             | PriorityClassName assigned to the Pods to prevent them from evicting.                                                                                                                                                     | string                                                                                                                       |
<!-- markdownlint-enable line-length -->

Example:

```yaml
networkLatencyExporter:
  install: true
  name: "network-latency-exporter"
  rbac:
    createClusterRole: true
    createClusterRoleBinding: true
    setupSecurityContext: true
    privileged: false
  createGrafanaDashboards: true
  serviceAccount:
    create: true
    name: "network-latency-exporter"
  image: prod.platform.system.monitoring_network-latency-exporter:latest
  resources:
    limits:
      cpu: 200m
      memory: 256Mi
    requests:
      cpu: 100m
      memory: 128Mi
  securityContext:
    runAsUser: "0"
  tolerations: []
  nodeSelector: {}
  affinity: {}
  labels:
    label.key: label-value
  annotations:
    annotation.key: annotation-value
  priorityClassName: priority-class
  extraArgs:
    - "--log.level=debug"
  requestTimeout: 5
  timeout: 60s
  packetsNum: 10
  packetSize: 64
  checkTarget: "UDP:80,TCP:80,ICMP"
  serviceMonitor:
    enabled: true
    interval: 30s
```

[Back to TOC](#table-of-content)

### node-exporter

<!-- markdownlint-disable line-length -->
| Field                      | Description                                                                                                                                                                                                            | Scheme                                                                                                                       |
| -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| install                    | Allows to disable create node-exporter during the deployment.                                                                                                                                                          | bool                                                                                                                         |
| paused                     | Set paused to reconciliation.                                                                                                                                                                                          | bool                                                                                                                         |
| image                      | A docker image to be used for the node-exporter deployment.                                                                                                                                                            | string                                                                                                                       |
| setupSecurityContext       | Allows to create PodSecurityPolicy or SecurityContextConstraints.                                                                                                                                                      | string                                                                                                                       |
| port                       | The port for node-exporter daemonset and service.                                                                                                                                                                      | int                                                                                                                          |
| resources                  | The resources that describe the compute resource requests and limits for single pods.                                                                                                                                  | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#resourcerequirements-v1-core) |
| securityContext            | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                     | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#podsecuritycontext-v1-core)    |
| tolerations                | Tolerations allow the pods to schedule onto nodes with matching taints.                                                                                                                                                | []v1.Toleration                                                                                                              |
| affinity                   | If specified, the pod's scheduling constraints                                                                                                                                                                         | *v1.Affinity                                                                                                                 |
| annotations                | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value" | map[string]string                                                                                                            |
| labels                     | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                       | map[string]string                                                                                                            |
| collectorTextfileDirectory | Directory for textfile. For more information, refer to [https://github.com/prometheus/node_exporter#textfile-collector](https://github.com/prometheus/node_exporter#textfile-collector)                                | string                                                                                                                       |
| extraArgs                  | Additional arguments for node-exporter container. For example: "--collector.systemd".                                                                                                                                  | list[string]                                                                                                                 |
| priorityClassName          | PriorityClassName assigned to the Pods to prevent them from evicting.                                                                                                                                                  | string                                                                                                                       |
| serviceMonitor             | Service monitor configuration for pulling metrics.                                                                                                                                                                     | [Monitor](#monitor)                                                                                                          |
<!-- markdownlint-enable line-length -->

Example:

```yaml
nodeExporter:
  install: true
  paused: false
  setupSecurityContext: false
  image: prom/node-exporter:v0.18.1
  port: 19100
  resources:
    limits:
      cpu: 200m
      memory: 200Mi
    requests:
      cpu: 100m
      memory: 100Mi
  securityContext:
    runAsUser: 2000
    fsGroup: 2000
  labels:
    label.key: label-value
  tolerations: []
  affinity: {}
  annotations:
    annotation.key: annotation-value
  priorityClassName: priority-class
  collectorTextfileDirectory: /var/spool/monitoring
  extraArgs:
    - --collector.systemd
  
  serviceMonitor:
    interval: 60s
    ...
```

[Back to TOC](#table-of-content)

### oauth-proxy

<!-- markdownlint-disable line-length -->
| Field | Description                                               | Scheme |
| ----- | --------------------------------------------------------- | ------ |
| image | A docker image to be used for the oauth-proxy deployment. | string |
<!-- markdownlint-enable line-length -->

Example:

```yaml
oAuthProxy:
  image: oauth2-proxy/oauth2-proxy:v7.0.1
```

[Back to TOC](#table-of-content)

### prometheus

<!-- markdownlint-disable line-length -->
| Field                                         | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                          | Scheme                                                                                                                          |
| --------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------- |
| image                                         | The image to be used for the `prometheus` deployment. The `prometheus` is a systems and service monitoring system. It collects metrics from configured targets at given intervals. For more information, refer to [https://github.com/prometheus/prometheus](https://github.com/prometheus/prometheus)                                                                                                                                                               | string                                                                                                                          |
| operator                                      | The operator parameters.                                                                                                                                                                                                                                                                                                                                                                                                                                             | [PrometheusOperator](#prometheus-operator)                                                                                      |
| install                                       | Install indicates whether Prometheus is to be installed. It can be changed for an already deployed service and the service is removed during the next reconciliation iteration.                                                                                                                                                                                                                                                                                      | *bool                                                                                                                           |
| configReloaderImage                           | The image to be used for `prometheus-config-reloader`. The `prometheus-config-reloaded` is an add-on to prometheus that monitors changes in prometheus.yaml and an HTTP request reloads the prometheus configuration. For more information, refer to [https://github.com/prometheus-operator/prometheus-operator/tree/master/cmd/prometheus-config-reloader](https://github.com/prometheus-operator/prometheus-operator/tree/master/cmd/prometheus-config-reloader)  | string                                                                                                                          |
| remoteWrite                                   | RemoteWriteSpec defines the remote_write configuration for prometheus. The `remote_write` allows to transparently send samples to a long term storage. For more information, refer to [https://prometheus.io/docs/operating/integrations/#remote-endpoints-and-storage](https://prometheus.io/docs/operating/integrations/#remote-endpoints-and-storage)                                                                                                             | []promv1.RemoteWriteSpec                                                                                                        |
| remoteWrite.basicAuth.createSecret            | Allow create secret for basic authentification automatically.                                                                                                                                                                                                                                                                                                                                                                                                        | map[string]string                                                                                                               |
| remoteWrite.basicAuth.createSecret.secretName | Name of the secret with which will be created.                                                                                                                                                                                                                                                                                                                                                                                                                       | string                                                                                                                          |
| remoteWrite.basicAuth.createSecret.username   | Username for basic authentification.                                                                                                                                                                                                                                                                                                                                                                                                                                 | string                                                                                                                          |
| remoteWrite.basicAuth.createSecret.password   | Password for basic authentification.                                                                                                                                                                                                                                                                                                                                                                                                                                 | string                                                                                                                          |
| remoteRead                                    | RemoteReadSpec defines the remote_read configuration for prometheus. The `remote_read` allows to transparently receive samples from a long term storage. For more information, refer to [https://prometheus.io/docs/operating/integrations/#remote-endpoints-and-storage](https://prometheus.io/docs/operating/integrations/#remote-endpoints-and-storage)                                                                                                           | []promv1.RemoteReadSpec                                                                                                         |
| secrets                                       | Secrets is a list of secrets in the same namespace as the Prometheus object, which shall be mounted into the Prometheus pods. The Secrets are mounted into `/etc/prometheus/secrets/<secret-name>`.                                                                                                                                                                                                                                                                  | []string                                                                                                                        |
| alerting                                      | Defines the details regarding alerting.                                                                                                                                                                                                                                                                                                                                                                                                                              | *promv1.AlertingSpec                                                                                                            |
| externalLabels                                | The labels to add to any time series or alerts when communicating with external systems (federation, remote storage, Alertmanager).                                                                                                                                                                                                                                                                                                                                  | map[string]string                                                                                                               |
| securityContext                               | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                                                                                                                                                                                                                                                                   | *[*v1.SecurityContext](https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#securitycontext-v1-core) |
| nodeSelector                                  | Defines which nodes the pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\"                                                                                                                                                                                                                                                                                                                                                   | map[string]string                                                                                                               |
| annotations                                   | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value"                                                                                                                                                                                                                                               | map[string]string                                                                                                               |
| labels                                        | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                                                                                                                                                                                                                                                                     | map[string]string                                                                                                               |
| affinity                                      | It specifies the pod's scheduling constraints. For more information, refer to [https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core)                                                                                                                                                                                                         | *v1.Affinity                                                                                                                    |
| resources                                     | Resources defines resource requests and limits for single pods.                                                                                                                                                                                                                                                                                                                                                                                                      | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#resourcerequirements-v1-core)    |
| storage                                       | Storage specifies how storage shall be used. For more information, refer to [https://github.com/prometheus-operator/prometheus-operator/blob/master/Documentation/api.md#storagespec](https://github.com/prometheus-operator/prometheus-operator/blob/master/Documentation/api.md#storagespec)                                                                                                                                                                       | *promv1.StorageSpec                                                                                                             |
| volumes                                       | Volumes allows configuration of additional volumes on the output StatefulSet definition. Volumes specified are appended to other volumes that are generated as a result of StorageSpec objects. For more information, refer to [https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#volume-v1-core](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#volume-v1-core)                                                            | []v1.Volume                                                                                                                     |
| volumeMounts                                  | VolumeMounts allows configuration of additional VolumeMounts on the output StatefulSet definition. VolumeMounts specified are appended to other VolumeMounts in the prometheus container, that are generated as a result of StorageSpec objects. For more informaiton, refer to [https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#volumemount-v1-core](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#volumemount-v1-core) | []v1.VolumeMount                                                                                                                |
| tlsConfig                                     | Defines the TLS parameters for prometheus. For more information, refer to [TLS guide](user-guides/tls.md)                                                                                                                                                                                                                                                                                                                                                            | object                                                                                                                          |
| tlsConfig.webTLSConfig                        | For more information, refer to [https://github.com/prometheus-operator/prometheus-operator/blob/main/Documentation/api.md#webtlsconfig](https://github.com/prometheus-operator/prometheus-operator/blob/main/Documentation/api.md#webtlsconfig)                                                                                                                                                                                                                      | object                                                                                                                          |
| tlsConfig.createSecret                        | Specifies content for secret and create it.                                                                                                                                                                                                                                                                                                                                                                                                                          | object                                                                                                                          |
| tlsConfig.createSecret.name                   | Name of secret with cert, ca and key.                                                                                                                                                                                                                                                                                                                                                                                                                                | string                                                                                                                          |
| tlsConfig.createSecret.cert                   | TLS certificate for prometheus.                                                                                                                                                                                                                                                                                                                                                                                                                                      | string                                                                                                                          |
| tlsConfig.createSecret.key                    | TLS key for prometheus.                                                                                                                                                                                                                                                                                                                                                                                                                                              | string                                                                                                                          |
| tlsConfig.createSecret.ca                     | TLS CA for prometheus.                                                                                                                                                                                                                                                                                                                                                                                                                                               | string                                                                                                                          |
| tlsConfig.generateCerts                       | Allows to configure generation of TLS certificate for Prometheus by cert-manager.                                                                                                                                                                                                                                                                                                                                                                                    | object                                                                                                                          |
| tlsConfig.generateCerts.enabled               | Allows to enable work with cert-manager.                                                                                                                                                                                                                                                                                                                                                                                                                             | bool                                                                                                                            |
| tlsConfig.generateCerts.secretName            | Name of secret generated by cert-manager.                                                                                                                                                                                                                                                                                                                                                                                                                            | string                                                                                                                          |
| tlsConfig.generateCerts.clusterIssuerName     | Defines name of Cluster Issuer. Otherwise, if this parameter is empty, self-signed non-cluster Issuer will be created and used in the new Certificate resource.                                                                                                                                                                                                                                                                                                      | string                                                                                                                          |
| tlsConfig.generateCerts.duration              | Defines duration of the certificate in days.                                                                                                                                                                                                                                                                                                                                                                                                                         | int                                                                                                                             |
| tlsConfig.generateCerts.renewBefore           | Specifies how long before expiry a certificate should be renewed.                                                                                                                                                                                                                                                                                                                                                                                                    | int                                                                                                                             |
| ingress                                       | Ingress allows to create Ingress for the Prometheus UI.                                                                                                                                                                                                                                                                                                                                                                                                              | *[Ingress](#ingress)                                                                                                            |
| retention                                     | Retention policy by time.                                                                                                                                                                                                                                                                                                                                                                                                                                            | string                                                                                                                          |
| retentionsize                                 | Retention policy by size [EXPERIMENTAL].                                                                                                                                                                                                                                                                                                                                                                                                                             | string                                                                                                                          |
| containers                                    | Containers allows injecting additional containers or modifying operator generated containers. This can be used to allow adding an authentication proxy to a Prometheus pod or to change the behavior of an operator generated container.                                                                                                                                                                                                                             | [[]Kubernetes core/v1.Container](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.24/#container-v1-core)        |
| externalUrl                                   | The external URL the Prometheus instances are available under. This is necessary to generate correct URLs. This is necessary if Prometheus is not served from the root of a DNS name.                                                                                                                                                                                                                                                                                | string                                                                                                                          |
| paused                                        | Set paused to reconciliation.                                                                                                                                                                                                                                                                                                                                                                                                                                        | bool                                                                                                                            |
| replicas                                      | Set replicas.                                                                                                                                                                                                                                                                                                                                                                                                                                                        | *int32                                                                                                                          |
| tolerations                                   | Tolerations allow the pods to schedule on nodes with matching taints.                                                                                                                                                                                                                                                                                                                                                                                                | []v1.Toleration                                                                                                                 |
| podMonitor                                    | Pod monitor for self monitoring.                                                                                                                                                                                                                                                                                                                                                                                                                                     | *[Monitor](#monitor)                                                                                                            |
| enableAdminAPI                                | Enable access to prometheus web admin API.                                                                                                                                                                                                                                                                                                                                                                                                                           | bool                                                                                                                            |
| query                                         | Defines the query command line flags when starting Prometheus. For more intoformation, refer to [https://github.com/prometheus-operator/prometheus-operator/blob/master/Documentation/api.md#queryspec](https://github.com/prometheus-operator/prometheus-operator/blob/master/Documentation/api.md#queryspec)                                                                                                                                                       | promv1.QuerySpec                                                                                                                |
| additionalScrape                              | The additionalScrape allows you to monitor hosts outside the cloud (balancers, graylog server, jenkins...)                                                                                                                                                                                                                                                                                                                                                           | object, [Promethues Config](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#scrape_config)            |
| additionalAlertManager                        | The additionalAlertManager allows you to use AlertManager outside the cloud                                                                                                                                                                                                                                                                                                                                                                                          | object, [Promethues Config](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#scrape_config)            |
| enableFeatures                                | Enable access to Prometheus disabled features. By default, enabled "auto-gomaxprox" feature.. Enabling disabled features is entirely outside the scope of what the maintainers will support and by doing so, you accept that this behavior may break at any time without notice. For more information see <https://prometheus.io/docs/prometheus/latest/disabled_features/>                                                                                          | []string                                                                                                                        |
| scrapeInterval                                | Interval between consecutive scrapes. Default: `30s`                                                                                                                                                                                                                                                                                                                                                                                                                 | string                                                                                                                          |
| scrapeTimeout                                 | Number of seconds to wait for target to respond before erroring. Default: `10s`                                                                                                                                                                                                                                                                                                                                                                                      | string                                                                                                                          |
| evaluationInterval                            | Interval between consecutive evaluations. Default: `30s`                                                                                                                                                                                                                                                                                                                                                                                                             | string                                                                                                                          |
| priorityClassName                             | PriorityClassName assigned to the Pods to prevent them from evicting.                                                                                                                                                                                                                                                                                                                                                                                                | string                                                                                                                          |
<!-- markdownlint-enable line-length -->

Example:

```yaml
prometheus:
  operator:
    ...
  image: prom/prometheus:v2.26.0
  install: true
  configReloaderImage: prometheus-operator/prometheus-config-reloader:v0.48.1
  enableAdminAPI: true
  query:
    lookbackDelta: 2m
  remoteWrite:
    - url: "http://graphite-remote-adapter:9201/write"
    - url: https://1.2.3.4:8429/api/v1/write
      tlsConfig:
        insecureSkipVerify: true
      basicAuth:
        createSecret:
          secretName: esm-vmagent
          username: prometheus
          password: prometheus
  remoteRead:
    - url: "http://graphite-remote-adapter:9201/read"
  secrets:
    - kube-etcd-client-certs
  alerting:
    alertmanagers:
      - namespace: default
        name: alertmanager-example
        port: web
  externalLabels:
    cluster: example_cloud
  securityContext:
    runAsUser: 2000
    fsGroup: 2000
  nodeSelector:
    node-role.kubernetes.io/worker: worker
  labels:
    label.key: label-value
  annotations:
    annotation.key: annotation-value
  affinity:
  resources:
    requests:
      cpu: 2000m
      memory: 8Gi
    limits:
      cpu: 6000m
      memory: 12Gi
  storage:
    volumeClaimTemplate:
      spec:
        storageClassName: nfs-dynamic-provisioning
        resources:
          requests:
            storage: 10Gi
  volumes:
    - name: additional-volume
      hostPath:
        path: /mnt/data/additional_volume
        type: Directory
  volumeMounts:
    - mountPath: /additional_volume
      name: additional-volume
  ingress:
    ...see example by link...
  retention: 24h
  retentionsize: 5Gi
  containers:
    ...see example by link...
  externalUrl: "prometheus.example.cloud.org"
  paused: false
  replicas: 1
  tolerations:
    - key: "example-key"
      operator: "Exists"
      effect: "NoSchedule"
  priorityClassName: priority-class
  podMonitor:
    ...see example by link...
  additionalAlertManager:
    - tls_config:
        insecure_skip_verify: true
      scheme: https
      static_configs:
      - targets:
        - "alertmaneger.example.outside.cloud.org"
  additionalScrape:
    - job_name: graylog
      honor_timestamps: true
      scrape_interval: 30s
      scrape_timeout: 10s
      metrics_path: /api/plugins/org.graylog.plugins.metrics.prometheus/metrics
      scheme: http
      static_configs:
        - targets:
            - 1.2.3.4
      basic_auth:
        username: admin
        password: admin
      tls_config:
        insecure_skip_verify: true
  enableFeatures:
    - auto-gomaxprocs
```

Example of Prometheus config for enabling [Agent mode](https://prometheus.io/blog/2021/11/16/agent/):

```yaml
prometheus:
  ...
  containers:
    - name: prometheus
      args:
        - '--config.file=/etc/prometheus/config_out/prometheus.env.yaml'
        - '--storage.agent.path=/prometheus'
        - '--enable-feature=agent'
        - '--web.enable-lifecycle'
  ...
```

[Back to TOC](#table-of-content)

#### prometheus-operator

<!-- markdownlint-disable line-length -->
| Field             | Description                                                                                                                                                                                                                                                                                                                                                           | Scheme                                                                                                                          |
| ----------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------- |
| image             | The image to be used for the `prometheus-operator` deployment. The `prometheus-operator` makes the Prometheus configuration Kubernetes native, and manages and operates Prometheus and Alertmanager clusters. For more information, refer to [https://github.com/prometheus-operator/prometheus-operator](https://github.com/prometheus-operator/prometheus-operator) | string                                                                                                                          |
| resources         | Resources defines resource requests and limits for single pods.                                                                                                                                                                                                                                                                                                       | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#resourcerequirements-v1-core)    |
| securityContext   | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                                                                                                                                                                    | *[*v1.SecurityContext](https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#securitycontext-v1-core) |
| paused            | Set paused to reconciliation.                                                                                                                                                                                                                                                                                                                                         | bool                                                                                                                            |
| tolerations       | Tolerations allow the pods to schedule on nodes with matching taints.                                                                                                                                                                                                                                                                                                 | []v1.Toleration                                                                                                                 |
| nodeSelector      | NodeSelector defines which nodes the pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\"                                                                                                                                                                                                                                       | map[string]string                                                                                                               |
| affinity                   | If specified, the pod's scheduling constraints                                                                                                                                                                         | *v1.Affinity                                                                                                                 |
| annotations       | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value"                                                                                                                                                | map[string]string                                                                                                               |
| labels            | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                                                                                                                                                                      | map[string]string                                                                                                               |
| podMonitor        | Pod monitor for self monitoring.                                                                                                                                                                                                                                                                                                                                      | *[Monitor](#monitor)                                                                                                            |
| priorityClassName | PriorityClassName assigned to the prometheus-operator Pods to prevent them from evicting                                                                                                                                                                                                                                                                              | string                                                                                                                          |

<!-- markdownlint-enable line-length -->

Example:

```yaml
prometheus:
  operator:
    image: prometheus-operator/prometheus-operator:v0.48.1
    resources:
      requests:
        cpu: 20m
        memory: 20Mi
      limits:
        cpu: 50m
        memory: 50Mi
    securityContext:
      runAsUser: 2000
      fsGroup: 2000
    paused: false
    tolerations:
      - key: "example-key"
        operator: "Exists"
        effect: "NoSchedule"
    nodeSelector:
      node-role.kubernetes.io/worker: worker
    labels:
      label.key: label-value
    annotations:
      annotation.key: annotation-value
    priorityClassName: priority-class
    podMonitor:
      ...see example by link...
```

[Back to TOC](#table-of-content)

### prometheus-adapter

<!-- markdownlint-disable line-length -->
| Field                                | Description                                                                                                                                                                                                                                                                                                                                                           | Scheme                                                                                                                          |
| ------------------------------------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------- |
| install                              | Allow to enable deploy Prometheus Adapter.                                                                                                                                                                                                                                                                                                                            | bool                                                                                                                            |
| image                                | The image to be used for the `prometheus-operator` deployment. The `prometheus-operator` makes the Prometheus configuration Kubernetes native, and manages and operates Prometheus and Alertmanager clusters. For more information, refer to [https://github.com/prometheus-operator/prometheus-operator](https://github.com/prometheus-operator/prometheus-operator) | string                                                                                                                          |
| resources                            | Resources defines resources requests and limits for single Pods.                                                                                                                                                                                                                                                                                                      | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#resourcerequirements-v1-core)    |
| securityContext                      | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                                                                                                                                                                    | *[*v1.SecurityContext](https://v1-18.docs.kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#securitycontext-v1-core) |
| paused                               | Set paused to reconciliation.                                                                                                                                                                                                                                                                                                                                         | bool                                                                                                                            |
| prometheusUrl                        | PrometheusURL used to connect to any tool with Prometheus compatible API. It will eventually contain query parameters to configure the connection                                                                                                                                                                                                                     | string                                                                                                                          |
| metricsRelistInterval                | MetricsRelistInterval is the interval at which to update the cache of available metrics from Prometheus                                                                                                                                                                                                                                                               | string                                                                                                                          |
| tolerations                          | Tolerations allow the pods to schedule onto nodes with matching taints.                                                                                                                                                                                                                                                                                               | []v1.Toleration                                                                                                                 |
| nodeSelector                         | NodeSelector defines which nodes the pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\"                                                                                                                                                                                                                                       | map[string]string                                                                                                               |
| affinity                   | If specified, the pod's scheduling constraints                                                                                                                                                                         | *v1.Affinity                                                                                                                 |
| annotations                          | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value"                                                                                                                                                | map[string]string                                                                                                               |
| labels                               | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                                                                                                                                                                      | map[string]string                                                                                                               |
| enableResourceMetrics                | Enable adapter for `metrics.k8s.io`. By default - `false`                                                                                                                                                                                                                                                                                                             | bool                                                                                                                            |
| enableCustomMetrics                  | Enable adapter for `custom.metrics.k8s.io`. By default - `true`                                                                                                                                                                                                                                                                                                       | bool                                                                                                                            |
| customScaleMetricRulesSelector       | CustomScaleMetricRulesSelector defines label selectors to select CustomScaleMetricRule resources across the cluster                                                                                                                                                                                                                                                   | []*[metav1.LabelSelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.23/#labelselector-v1-meta)           |
| APIService.resourceMetrics           | Enable/disable creating APIServices for `metrics.k8s.io`                                                                                                                                                                                                                                                                                                              | bool                                                                                                                            |
| APIService.customMetrics             | Enable/disable creating APIServices for `custom.metrics.k8s.io`                                                                                                                                                                                                                                                                                                       | bool                                                                                                                            |
| auth                                 | Client credentials to connect to Prometheus or Victoriametrics endpoints. (Only basic authentication is supported)                                                                                                                                                                                                                                                    | map[string]string                                                                                                               |
| auth.basicAuth                       | Allow to specify client auth configuration as secret reference                                                                                                                                                                                                                                                                                                        | *[v1.SecretKeySelector](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#secretkeyselector-v1-core)         |
| auth.basicAuth.createSecret          | Allow to create secret with auth credentials automatically during deploy                                                                                                                                                                                                                                                                                              | map[string]string                                                                                                               |
| auth.basicAuth.createSecret.username | Allow to specify client username                                                                                                                                                                                                                                                                                                                                      | string                                                                                                                          |
| auth.basicAuth.createSecret.password | Allow to specify client password                                                                                                                                                                                                                                                                                                                                      | string                                                                                                                          |
| priorityClassName                    | PriorityClassName assigned to the Pods to prevent them from evicting.                                                                                                                                                                                                                                                                                                 | string                                                                                                                          |
<!-- markdownlint-enable line-length -->

Example (automatically created secrets for auth):

```yaml
prometheusAdapter:
  install: true
  image: "k8s-prometheus-adapter-amd64:v0.6.0"
  securityContext:
    runAsUser: 2000
    fsGroup: 2000
  resources:
    limits:
      cpu: 200m
      memory: 200Mi
    requests:
      cpu: 100m
      memory: 100Mi
  nodeSelector:
    node-role.kubernetes.io/worker: worker
  labels:
    label.key: label-value
  annotations:
    annotation.key: annotation-value
  priorityClassName: priority-class
  metricsRelistInterval: "1m"
  prometheusUrl: "http://prometheus-operated.monitoring.svc:9090"
  APIService:
    resourceMetrics: true
    customMetrics: true
  enableResourceMetrics: true
  enableCustomMetrics: true
  customScaleMetricRulesSelector:
    - matchExpressions:
        - key: app.kubernetes.io/component
          operator: In
          values: [ "monitoring" ]
  auth:
    createSecret:
      basicAuth:
        username: prometheus
        password: prometheus
  operator:
    ...see example by link...
```

Example (precreated secrets for auth):

```yaml
prometheusAdapter:
  install: true
  image: "k8s-prometheus-adapter-amd64:v0.6.0"
  securityContext:
    runAsUser: 2000
    fsGroup: 2000
  resources:
    limits:
      cpu: 200m
      memory: 200Mi
    requests:
      cpu: 100m
      memory: 100Mi
  nodeSelector:
    node-role.kubernetes.io/worker: worker
  labels:
    label.key: label-value
  annotations:
    annotation.key: annotation-value
  metricsRelistInterval: "1m"
  prometheusUrl: "http://prometheus-operated.monitoring.svc:9090"
  APIService:
    resourceMetrics: true
    customMetrics: true
  enableResourceMetrics: true
  enableCustomMetrics: true
  customScaleMetricRulesSelector:
    - matchExpressions:
        - key: app.kubernetes.io/component
          operator: In
          values: [ "monitoring" ]
  auth:
    basicAuth:
      username:
        name: secret
        key: username
      password:
        name: secret
        key: password
  operator:
    ...see example by link...
```

[Back to TOC](#table-of-content)

#### prometheus-adapter-operator

<!-- markdownlint-disable line-length -->
| Field             | Description                                                                                                                                                                                                            | Scheme                                                                                                                       |
| ----------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| image             | A docker image to use for prometheus-adapter-operator deployment                                                                                                                                                       | string                                                                                                                       |
| resources         | Resources defines resources requests and limits for single Pods.                                                                                                                                                       | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.17/#resourcerequirements-v1-core) |
| securityContext   | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                     | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#podsecuritycontext-v1-core)    |
| nodeSelector      | NodeSelector defines which nodes the pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\"                                                                                        | map[string]string                                                                                                            |
| affinity                   | If specified, the pod's scheduling constraints                                                                                                                                                                         | *v1.Affinity                                                                                                                 |
| annotations       | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value" | map[string]string                                                                                                            |
| labels            | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                       | map[string]string                                                                                                            |
| priorityClassName | PriorityClassName assigned to the Pods to prevent them from evicting.                                                                                                                                                  | string                                                                                                                       |
| tlsEnabled        | TLS configuration is enabled/disabled. By default, it is disabled.                                                                                                                                                     | boolean                                                                                                                      |
| tlsConfig         | Allow to specify client TLS configuration.                                                                                                                                                                             | [*TLSConfig](#prometheus-adapter-operator-tls-config)                                                                        |
<!-- markdownlint-enable line-length -->

```yaml
prometheusAdapter:
  operator:
    image: "product/prod.platform.system.monitoring_k8s-scaling_prometheus-adapter-operator:latest"
    securityContext:
      runAsUser: 2000
      fsGroup: 2000
    resources:
      limits:
        cpu: 200m
        memory: 200Mi
      requests:
        cpu: 100m
        memory: 100Mi
    nodeSelector:
      node-role.kubernetes.io/worker: worker
    labels:
      label.key: label-value
    annotations:
      annotation.key: annotation-value
    priorityClassName: priority-class
    tlsEnabled: true
    tlsConfig:
      generateCerts:
        enabled: true
        duration: 365
        renewBefore: 15
        clusterIssuerName: "dev-cluster-issuer"
        secretName: "prometheus-adapter-client-tls-secret"
```

[Back to TOC](#table-of-content)

#### prometheus-adapter-operator-tls-config

TLSConfig holds SSL/TLS configuration attributes.
The parameters are required if SSL/TLS connection is required between Kubernetes cluster and prometheus-adapter-operator.
This section is applicable only if `tlsEnabled` is set to `true`.

<!-- markdownlint-disable line-length -->
| Parameter                         | Type                  | Mandatory | Default value                          | Description                                                                                                                                                                                                                                                                                                                                                                                                                 |
| --------------------------------- | --------------------- | --------- | -------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `caSecret`                        | *v1.SecretKeySelector | no        | `-`                                    | Secret containing the CA certificate to use for the targets.                                                                                                                                                                                                                                                                                                                                                                |
| `certSecret`                      | *v1.SecretKeySelector | no        | `-`                                    | Secret containing the client certificate file for the targets.                                                                                                                                                                                                                                                                                                                                                              |
| `keySecret`                       | *v1.SecretKeySelector | no        | `-`                                    | Secret containing the client key file for the targets.                                                                                                                                                                                                                                                                                                                                                                      |
| `existingSecret`                  | string                | no        | `-`                                    | Name of the pre-existing secret that contains TLS configuration for prometheus-adapter. If specified, `generateCerts.enabled` must be set to `false`. The `existingSecret` is expected to contain CA certificate, TLS key and TLS certificate in `ca.crt`, `tls.key` and `tls.crt` fields respectively. Use either `existingSecret` or the combination of `caSecret`, `certSecret` and `keySecret`. Do not use it together. |
| `generateCerts.enabled`           | boolean               | no        | `true`                                 | Generation of certificate is enabled by default. If `tlsConfig.existingSecret` or the combination of `tlsConfig.caSecret`, `tlsConfig.certSecret` and `tlsConfig.keySecret` is specified, `tlsConfig.generateCerts` section will be skipped. `cert-manager` will generate certificate with the name configured using `generateCerts.secretName`, if it doesn't exist already.                                               |
| `generateCerts.clusterIssuerName` | string                | no        | `-`                                    | Cluster issuer name for generated certificate. This is a mandatory field if `generateCerts.enabled` is set to `true`.                                                                                                                                                                                                                                                                                                       |
| `generateCerts.duration`          | integer               | no        | `365`                                  | Duration in days, until which issued certificate will be valid.                                                                                                                                                                                                                                                                                                                                                             |
| `generateCerts.renewBefore`       | integer               | no        | `15`                                   | Number of days before which certificate must be renewed.                                                                                                                                                                                                                                                                                                                                                                    |
| `generateCerts.secretName`        | string                | no        | `prometheus-adapter-client-tls-secret` | Name of the new secret that needs to be created for storing TLS configuration of prometheus-adapter.                                                                                                                                                                                                                                                                                                                        |
| `createSecret`                    | object                | no        | `-`                                    | New secret with the name `tlsConfig.createSecret.secretName` will be created using already known certificate content. If `tlsConfig.existingSecret` or the combination of `tlsConfig.caSecret`, `tlsConfig.certSecret` and `tlsConfig.keySecret` is specified, `tlsConfig.createSecret` section will be skipped.                                                                                                            |
| `createSecret.ca`                 | string                | no        | `-`                                    | Already known CA certificate will be added to newly created secret.                                                                                                                                                                                                                                                                                                                                                         |
| `createSecret.key`                | string                | no        | `-`                                    | Already known TLS key will be added to newly created secret.                                                                                                                                                                                                                                                                                                                                                                |
| `createSecret.cert`               | string                | no        | `-`                                    | Already known TLS certificate will be added to newly created secret.                                                                                                                                                                                                                                                                                                                                                        |
| `createSecret.secretName`         | string                | no        | `prometheus-adapter-client-tls-secret` | Already known TLS certificate will be added to newly created secret.                                                                                                                                                                                                                                                                                                                                                        |
<!-- markdownlint-enable line-length -->

[Back to TOC](#table-of-content)

### prometheus-rules

<!-- markdownlint-disable line-length -->
| Field       | Description                                                                                                                                                                                                            | Scheme            |
| ----------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------- |
| install     | Allows to install Prometheus Rules for monitoring-operator                                                                                                                                                             | bool              |
| ruleGroups  | List of groups to be installed                                                                                                                                                                                         | list[string]      |
| override    | Allows overriding of Prometheus Rules for monitoring-operator                                                                                                                                                          | list[object]      |
| annotations | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value" | map[string]string |
| labels      | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                       | map[string]string |
<!-- markdownlint-enable line-length -->

Example:

```yaml
prometheusRules:
  install: true
  ruleGroups:
    - Heartbeat
    - SelfMonitoring
    - AlertManager
    - KubebernetesAlerts
    - NodeExporters
    - DockerContainers
    - HAproxy
    - Etcd
  override:
      - group: SelfMonitoring
        alert: PrometheusNotificationsBacklog
        for: 0m
        expr: min_over_time(prometheus_notifications_queue_length[20m]) > 0
        severity: high
  labels:
    label.key: label-value
  annotations:
    annotation.key: annotation-value
```

[Back to TOC](#table-of-content)

### promitor-agent-scraper

<!-- markdownlint-disable line-length -->
| Field                                                         | Description                                                                                                                                                                                                                                                  | Scheme                                                                                                                       |
| ------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ---------------------------------------------------------------------------------------------------------------------------- |
| install                                                       | Allows to enable or disable deploy promitor-agent-scraper.                                                                                                                                                                                                   | boolean                                                                                                                      |
| replicas                                                      | Number of created pods.                                                                                                                                                                                                                                      | integer                                                                                                                      |
| name                                                          | A name of the microservice to deploy with. This name is used as the name of the microservice deployment and in labels.                                                                                                                                       | string                                                                                                                       |
| image                                                         | A Docker image to deploy the promitor-agent-scraper.                                                                                                                                                                                                         | string                                                                                                                       |
| imagePullPolicy                                               | Image pull policy to use for promitor-agent-scraper deployment.                                                                                                                                                                                              | string                                                                                                                       |
| azureAuthentication.mode                                      | Authentication type to use to authenticate. Options are ServicePrincipal (default), UserAssignedManagedIdentity or SystemAssignedManagedIdentity.                                                                                                            | string                                                                                                                       |
| azureAuthentication.identity.id                               | Id of the Azure AD entity to authenticate with.                                                                                                                                                                                                              | string                                                                                                                       |
| azureAuthentication.identity.key                              | Secret of the Azure AD entity to authenticate with. Sets the environment variable PROMITOR_AUTH_APPKEY through the secrets.appKeySecret field in the Secret.                                                                                                 | string                                                                                                                       |
| azureAuthentication.identity.binding                          | Aad Pod Identity name, when using UserAssignedManagedIdentity or SystemAssignedManagedIdentity as mode.                                                                                                                                                      | string                                                                                                                       |
| metricSinks.prometheusScrapingEndpoint.baseUriPath            | Controls the path where the scraping endpoint for Prometheus is being exposed.                                                                                                                                                                               | string                                                                                                                       |
| metricSinks.prometheusScrapingEndpoint.enableMetricTimestamps | Defines whether or not a timestamp should be included when the value was scraped on Azure Monitor. Supported values are True to opt-in & False to opt-out.                                                                                                   | boolean                                                                                                                      |
| metricSinks.prometheusScrapingEndpoint.metricUnavailableValue | Defines the value that will be reported if a metric is unavailable.                                                                                                                                                                                          | string                                                                                                                       |
| metricSinks.prometheusScrapingEndpoint.labelTransformation    | Controls how label values are reported to Prometheus by using transformation. Options are None & Lowercase.                                                                                                                                                  | string                                                                                                                       |
| resourceDiscovery.enabled                                     | Indication whether or not resource discovery is enabled through the Promitor Resource Discovery agent.                                                                                                                                                       | boolean                                                                                                                      |
| resourceDiscovery.host                                        | DNS name or IP address of the Promitor Resource Discovery agent.                                                                                                                                                                                             | string                                                                                                                       |
| resourceDiscovery.port                                        | Port (UDP) address of the Promitor Resource Discovery agent.                                                                                                                                                                                                 | integer                                                                                                                      |
| telemetry.defaultLogLevel                                     | Defines the default minimum log level that should be logged if a sink does not provide one. Allowed values are Trace, Debug, Information, Warning, Error, Critical, None ordered from most to least verbose.                                                 | string                                                                                                                       |
| telemetry.applicationInsights.enabled                         | Determines if the sink is used or not.                                                                                                                                                                                                                       | boolean                                                                                                                      |
| telemetry.applicationInsights.key                             | Defines the instrumentation key to use when sending telemetry to Azure Application Insights.                                                                                                                                                                 | string                                                                                                                       |
| telemetry.applicationInsights.logLevel                        | Verbosity to use for this sink, if not specified then the telemetry.defaultLogLevel will be used.                                                                                                                                                            | string                                                                                                                       |
| telemetry.containerLogs.enabled                               | Determines if the sink is used or not.                                                                                                                                                                                                                       | boolean                                                                                                                      |
| telemetry.containerLogs.logLevel                              | Verbosity to use for this sink, if not specified then the telemetry.defaultLogLevel will be used.                                                                                                                                                            | string                                                                                                                       |
| azureMetadata.tenantId                                        | The id of the Azure tenant that will be queried.                                                                                                                                                                                                             | string                                                                                                                       |
| azureMetadata.subscriptionId                                  | The id of the default subscription to query.                                                                                                                                                                                                                 | string                                                                                                                       |
| azureMetadata.resourceGroupName                               | The name of the default resource group to query.                                                                                                                                                                                                             | string                                                                                                                       |
| azureMetadata.cloud                                           | The name of the Azure cloud to use. Options are Global (default), China, UsGov & Germany.                                                                                                                                                                    | string                                                                                                                       |
| metricDefaults.aggregation.interval                           | The default interval which defines over what period measurements of a metric should be aggregated. A cron that fits your needs.                                                                                                                              | string                                                                                                                       |
| metricDefaults.scraping.schedule                              | A cron expression that controls the frequency of which all the configured metrics will be scraped from Azure Monitor. You can use crontab-generator.org to generate a cron that fits your needs.                                                             | string                                                                                                                       |
| metrics                                                       | List of metrics to scrape.                                                                                                                                                                                                                                   | list[[Metric](https://docs.promitor.io/configuration/v2.x/metrics/)]                                                         |
| secrets.createSecret                                          | To use your own secret, set createSecret to false and define the name/keys that your secret uses. Indication if you want to bring your own secret level of logging.                                                                                          | boolean                                                                                                                      |
| secrets.secretName                                            | Name of the secret for Azure AD identity secret.                                                                                                                                                                                                             | string                                                                                                                       |
| secrets.appKeySecret                                          | Name of the field for Azure AD identity secret in the Secret resource.                                                                                                                                                                                       | string                                                                                                                       |
| service.type                                                  | Type of promitor-agent-scraper service.                                                                                                                                                                                                                      | string                                                                                                                       |
| service.port                                                  | Port of promitor-agent-scraper which use in service.                                                                                                                                                                                                         | integer                                                                                                                      |
| service.targetPort                                            | Target port of promitor-agent-scraper which use in service and container.                                                                                                                                                                                    | integer                                                                                                                      |
| service.labels                                                | Labels set which will create in service.                                                                                                                                                                                                                     | map[string]string                                                                                                            |
| service.azureLoadBalancer.dnsPrefix                           | Prefix for DNS name to expose the service on using \<name\>.\<location\>.cloudapp.azure.com format. This setting is specific to Azure Kubernetes Service.                                                                                                    | string                                                                                                                       |
| service.azureLoadBalancer.exposeInternally                    | To restrict access to Promitor by exposing it through an internal load balancer. This setting is specific to Azure Kubernetes Service.                                                                                                                       | string                                                                                                                       |
| serviceAccount.install                                        | Allow to disable create ServiceAccount during deploy.                                                                                                                                                                                                        | boolean                                                                                                                      |
| serviceAccount.name                                           | Provide a name in place of promitor-agent-scraper for ServiceAccount.                                                                                                                                                                                        | string                                                                                                                       |
| serviceAccount.annotations                                    | Will add the provided map to the annotations for the created serviceAccount.                                                                                                                                                                                 | map[string]string                                                                                                            |
| serviceAccount.labels                                         | Will add the provided map to the labels for the created serviceAccount.                                                                                                                                                                                      | map[string]string                                                                                                            |
| serviceMonitor                                                | ServiceMonitor holds configuration attributes for promitor-agent-scraper.                                                                                                                                                                                    | object                                                                                                                       |
| resources                                                     | The resources that describe the compute resource requests and limits for single pods.                                                                                                                                                                        | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#resourcerequirements-v1-core) |
| securityContext                                               | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                                                           | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#podsecuritycontext-v1-core)    |
| affinity                                                      | It specifies the pod's scheduling constraints. For more information, refer to [https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core) | *v1.Affinity                                                                                                                 |
| nodeSelector                                                  | Defines which nodes the pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\"                                                                                                                                           | map[string]string                                                                                                            |
| tolerations                                                   | Tolerations allow the pods to schedule onto nodes with matching taints.                                                                                                                                                                                      | []v1.Toleration                                                                                                              |
| annotations                                                   | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value"                                       | map[string]string                                                                                                            |
| labels                                                        | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                                                             | map[string]string                                                                                                            |
| priorityClassName                                             | PriorityClassName assigned to the Pods to prevent them from evicting.                                                                                                                                                                                        | string                                                                                                                       |
<!-- markdownlint-enable line-length -->

Example:

```yaml
promitorAgentScraper:
  install: true
  replicaCount: 1
  name: promitor-agent-scraper
  image: ghcr.io/tomkerkhove/promitor-agent-scraper:2.5.1
  imagePullPolicy: IfNotPresent
  azureAuthentication:
    mode: "ServicePrincipal"
    identity:
      id: "<azure-ad-id>"
      key: "<azure-ad-secret>"
      binding: ""
  metricSinks:
    prometheusScrapingEndpoint:
      baseUriPath: /metrics
      enableMetricTimestamps: true
      metricUnavailableValue: NaN
      labelTransformation: None
      enableServiceDiscovery: true
  resourceDiscovery:
    enabled: false
    host: ""
    port: 80
  telemetry:
    defaultLogLevel: "Error"
    applicationInsights:
      enabled: false
      key: ""
      logLevel: ""
    containerLogs:
      enabled: true
      logLevel: ""
  azureMetadata:
    tenantId: "<tenant-id>"
    subscriptionId: "<subscription-id>"
    resourceGroupName: "<resource-group-name>"
    cloud: "Global"
  metricDefaults:
    aggregation:
      interval: 00:05:00
    scraping:
      schedule: "*/5 * * * *"
  metrics:
    - name: azure_node_cpu_usage_millicores
      description: "Node CPU usage millicores"
      resourceType: KubernetesService
      azureMetricConfiguration:
        metricName: node_cpu_usage_millicores
        aggregation:
          type: Average
      resources:
        - clusterName: <cluster-name>
  secrets:
    createSecret: true
    secretName: promitor-agent-scraper
    appKeySecret: azure-app-key
  service:
    type: ClusterIP
    port: 8888
    targetPort: 5000
    labels: {}
    azureLoadBalancer:
      dnsPrefix:
      exposeInternally: false
  resources:
    limits:
      cpu: 200m
      memory: 256Mi
    requests:
      cpu: 100m
      memory: 128Mi
  serviceAccount:
    install: true
    name: promitor-agent-scraper
    annotations: {}
    labels: {}
  serviceMonitor:
    install: true
    interval: 5m
    telemetryPath: /metrics
    labels:
      app.kubernetes.io/component: monitoring
  securityContext:
    runAsUser: 2000
    fsGroup: 2000
  affinity: {}
  nodeSelector: {}
  tolerations: []
  annotations: {}
  labels: {}
  priorityClassName: priority-class
```

[Back to TOC](#table-of-content)

### promitor-agent-resource-discovery

<!-- markdownlint-disable line-length -->
| Field                                       | Description                                                                                                                                                                                                                                                  | Scheme                                                                                                                        |
| ------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ----------------------------------------------------------------------------------------------------------------------------- |
| install                                     | Allows to enable or disable deploy promitor-agent-resource-discovery.                                                                                                                                                                                        | boolean                                                                                                                       |
| replicas                                    | Number of created pods.                                                                                                                                                                                                                                      | integer                                                                                                                       |
| name                                        | A name of the microservice to deploy with. This name is used as the name of the microservice deployment and in labels.                                                                                                                                       | string                                                                                                                        |
| image                                       | A Docker image to deploy the promitor-agent-resource-discovery.                                                                                                                                                                                              | string                                                                                                                        |
| imagePullPolicy                             | Image pull policy to use for promitor-agent-resource-discovery deployment.                                                                                                                                                                                   | string                                                                                                                        |
| imagePullSecrets                            | Reference to one or more secrets to be used when pulling images.                                                                                                                                                                                             | [list](https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/)                                 |
| azureAuthentication.mode                    | Authentication type to use to authenticate. Options are ServicePrincipal (default), UserAssignedManagedIdentity or SystemAssignedManagedIdentity.                                                                                                            | string                                                                                                                        |
| azureAuthentication.identity.id             | Id of the Azure AD entity to authenticate with.                                                                                                                                                                                                              | string                                                                                                                        |
| azureAuthentication.identity.key            | Secret of the Azure AD entity to authenticate with. Sets the environment variable PROMITOR_AUTH_APPKEY through the secrets.appKeySecret field in the Secret.                                                                                                 | string                                                                                                                        |
| azureAuthentication.identity.binding        | Aad Pod Identity name, when using UserAssignedManagedIdentity or SystemAssignedManagedIdentity as mode.                                                                                                                                                      | string                                                                                                                        |
| azureLandscape.tenantId                     | The id of the Azure tenant that will be queried.                                                                                                                                                                                                             | string                                                                                                                        |
| azureLandscape.subscriptions                | The ids of the subscriptions to query.                                                                                                                                                                                                                       | list[string]                                                                                                                  |
| azureLandscape.cloud                        | The name of the Azure cloud to use. Options are Global (default), China, UsGov & Germany.                                                                                                                                                                    | string                                                                                                                        |
| resourceDiscoveryGroups                     | Resource groups to query.                                                                                                                                                                                                                                    | list[object]                                                                                                                  |
| cache.enabled                               | Indication whether or not discovered resources should be cached in-memory to avoid Azure throttling.                                                                                                                                                         | boolean                                                                                                                       |
| cache.durationInMinutes                     | Amount of minutes to cache discovered resources.                                                                                                                                                                                                             | integer                                                                                                                       |
| telemetry.defaultLogLevel                   | Defines the default minimum log level that should be logged if a sink does not provide one. Allowed values are Trace, Debug, Information, Warning, Error, Critical, None ordered from most to least verbose.                                                 | string                                                                                                                        |
| telemetry.applicationInsights.enabled       | Determines if the sink is used or not.                                                                                                                                                                                                                       | boolean                                                                                                                       |
| telemetry.applicationInsights.key           | Defines the instrumentation key to use when sending telemetry to Azure Application Insights.                                                                                                                                                                 | string                                                                                                                        |
| telemetry.applicationInsights.logLevel      | Verbosity to use for this sink, if not specified then the telemetry.defaultLogLevel will be used.                                                                                                                                                            | string                                                                                                                        |
| telemetry.containerLogs.enabled             | Determines if the sink is used or not.                                                                                                                                                                                                                       | boolean                                                                                                                       |
| telemetry.containerLogs.logLevel            | Verbosity to use for this sink, if not specified then the telemetry.defaultLogLevel will be used.                                                                                                                                                            | string                                                                                                                        |
| secrets.createSecret                        | To use your own secret, set createSecret to false and define the name/keys that your secret uses. Indication if you want to bring your own secret level of logging.                                                                                          | boolean                                                                                                                       |
| secrets.secretName                          | Name of the secret for Azure AD identity secret.                                                                                                                                                                                                             | string                                                                                                                        |
| secrets.appKeySecret                        | Name of the field for Azure AD identity secret in the Secret resource.                                                                                                                                                                                       | string                                                                                                                        |
| secrets.extra                               | Add extra secrets.                                                                                                                                                                                                                                           | map[string]string                                                                                                             |
| service.type                                | Type of promitor-agent-scraper service.                                                                                                                                                                                                                      | string                                                                                                                        |
| service.port                                | Port of promitor-agent-scraper which use in service.                                                                                                                                                                                                         | integer                                                                                                                       |
| service.targetPort                          | Target port of promitor-agent-scraper which use in service and container.                                                                                                                                                                                    | integer                                                                                                                       |
| service.labels                              | Labels set which will create in service.                                                                                                                                                                                                                     | map[string]string                                                                                                             |
| service.azureLoadBalancer.dnsPrefix         | Prefix for DNS name to expose the service on using \<name\>.\<location\>.cloudapp.azure.com format. This setting is specific to Azure Kubernetes Service.                                                                                                    | string                                                                                                                        |
| service.azureLoadBalancer.exposeInternally  | To restrict access to Promitor by exposing it through an internal load balancer. This setting is specific to Azure Kubernetes Service.                                                                                                                       | string                                                                                                                        |
| serviceAccount.install                      | Allow to disable create ServiceAccount during deploy.                                                                                                                                                                                                        | boolean                                                                                                                       |
| serviceAccount.name                         | Provide a name in place of promitor-agent-scraper for ServiceAccount.                                                                                                                                                                                        | string                                                                                                                        |
| serviceAccount.annotations                  | Will add the provided map to the annotations for the created serviceAccount.                                                                                                                                                                                 | map[string]string                                                                                                             |
| serviceAccount.labels                       | Will add the provided map to the labels for the created serviceAccount.                                                                                                                                                                                      | map[string]string                                                                                                             |
| serviceAccount.automountServiceAccountToken | If service account token should be mounted inside pod. Set this to true if you plan on using Pod Security Policy.                                                                                                                                            | boolean                                                                                                                       |
| serviceMonitor                              | ServiceMonitor holds configuration attributes for promitor-agent-resource-discovery.                                                                                                                                                                         | object                                                                                                                        |
| resources                                   | The resources that describe the compute resource requests and limits for single pods.                                                                                                                                                                        | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#resourcerequirements-v1-core)  |
| securityContext                             | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                                                           | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#podsecuritycontext-v1-core)     |
| containerSecurityContext                    | Custom security context object for container                                                                                                                                                                                                                 | [object](https://kubernetes.io/docs/tasks/configure-pod-container/security-context/#set-the-security-context-for-a-container) |
| affinity                                    | It specifies the pod's scheduling constraints. For more information, refer to [https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core) | *v1.Affinity                                                                                                                  |
| nodeSelector                                | Defines which nodes the pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\"                                                                                                                                           | map[string]string                                                                                                             |
| tolerations                                 | Tolerations allow the pods to schedule onto nodes with matching taints.                                                                                                                                                                                      | []v1.Toleration                                                                                                               |
| annotations                                 | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value"                                       | map[string]string                                                                                                             |
| labels                                      | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                                                             | map[string]string                                                                                                             |
| extraVolumeMounts                           | Pass extra volumeMounts to the promitor-agent-resource-discovery container.                                                                                                                                                                                  | list                                                                                                                          |
| extraVolumes                                | Pass extra volumes to the promitor-agent-resource-discovery deployment.                                                                                                                                                                                      | list                                                                                                                          |
| priorityClassName                           | PriorityClassName assigned to the Pods to prevent them from evicting.                                                                                                                                                                                        | string                                                                                                                        |
<!-- markdownlint-enable line-length -->

Example:

```yaml
promitorAgentResourceDiscovery:
  install: true
  azureAuthentication:
    mode: ServicePrincipal
    identity:
      id: <azure_entity_id>
      key: <azure_entity_key>
      binding: ""
  azureLandscape:
    cloud: Global
    tenantId: <tenant_id>
    subscriptions:
      - SUBSCRIPTON-ID-ABC
      - SUBSCRIPTON-ID-GHI
  resourceDiscoveryGroups:
    - name: psql-rg
      type: PostgreSql
      criteria:
        include:
          subscriptionIds:
          - SUBSCRIPTON-ID-ABC
          - SUBSCRIPTON-ID-GHI
          resourceGroups:
          - promitor-resource-group-1
          - promitor-resource-group-2
          tags:
            app: promitor
            region: europe
          regions:
          - northeurope
          - westeurope
```

[Back to TOC](#table-of-content)

### publicCloudName

Allows to disable or enable specific dashboards, rules and service monitors for the specified public cloud.
Allowed values: aws, azure, google, ""

See
[Enabling/disabling specific monitoring components](integrations.md#enablingdisabling-specific-monitoring-components)
for additional information about `publicCloudName` parameter work.

[Back to TOC](#table-of-content)

### promxy

<!-- markdownlint-disable line-length -->
| Field                          | Description                                                                                                                                                                                                                                                  | Scheme                                                                                                                       |
| ------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ---------------------------------------------------------------------------------------------------------------------------- |
| install                        | Allows to enable or disable deploy promxy.                                                                                                                                                                                                                   | boolean                                                                                                                      |
| name                           | A name of the microservice to deploy with. This name is used as the name of the microservice deployment and in labels.                                                                                                                                       | string                                                                                                                       |
| image                          | A Docker image to deploy the promxy.                                                                                                                                                                                                                         | string                                                                                                                       |
| extraArgs                      | Additional arguments for promxy container.                                                                                                                                                                                                                   | list[string]                                                                                                                 |
| config                         | Configuration for promxy.                                                                                                                                                                                                                                    | object                                                                                                                       |
| config.serverGroups            | Base configuration of targets (server groups) for promxy. Every group should contain address (service or ingress for scraping Prometheus metrics), scheme (http or https) and unique label for identification in Grafana. Scheme "http" is using by default. | list[object]                                                                                                                 |
| config.serverGroups[N].address | URL (service or ingress) for scraping Prometheus metrics.                                                                                                                                                                                                    | string                                                                                                                       |
| config.serverGroups[N].scheme  | Configures the protocol scheme used for requests (http or https). Defaults to http.                                                                                                                                                                          | string                                                                                                                       |
| config.serverGroups[N].label   | Unique label to be added to metrics retrieved from this server group for identification in Grafana.                                                                                                                                                          | string                                                                                                                       |
| config.detailedConfig          | Detailed configuration for promxy in the YAML format. If this parameter is specified, the rest of parameters from the config section will not be used.                                                                                                       | [Promxy Configuration](https://github.com/jacksontj/promxy/blob/v0.0.71/cmd/promxy/config.yaml)                              |
| resources                      | The resources that describe the compute resource requests and limits for single pods.                                                                                                                                                                        | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#resourcerequirements-v1-core) |
| port                           | Port for promxy container and service.                                                                                                                                                                                                                       | integer                                                                                                                      |
| configmapReload                | Additional container with small sidecar that reload promxy when configmap with configuration is changing.                                                                                                                                                    | object                                                                                                                       |
| configmapReload.install        | Allow to disable deploy promxy-configmap-reload container.                                                                                                                                                                                                   | boolean                                                                                                                      |
| configmapReload.image          | A docker image to use for promxy-configmap-reload deployment.                                                                                                                                                                                                | string                                                                                                                       |
| configmapReload.resources      | The resources describes the compute resource requests and limits for promxy-configmap-reload.                                                                                                                                                                | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#resourcerequirements-v1-core) |
| securityContext                | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                                                           | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#podsecuritycontext-v1-core)    |
| tolerations                    | Tolerations allow the pods to schedule onto nodes with matching taints.                                                                                                                                                                                      | []v1.Toleration                                                                                                              |
| nodeSelector                   | Defines which nodes the pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\"                                                                                                                                           | map[string]string                                                                                                            |
| affinity                                            | If specified, the pod's scheduling constraints                                                                                                                                                                                      | *v1.Affinity                                                                                                                   |
| annotations                    | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value"                                       | map[string]string                                                                                                            |
| labels                         | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                                                             | map[string]string                                                                                                            |
| priorityClassName              | PriorityClassName assigned to the Pods to prevent them from evicting.                                                                                                                                                                                        | string                                                                                                                       |
<!-- markdownlint-enable line-length -->

Example:

```yaml
promxy:
  install: true
  name: "promxy"
  image: jacksontj/promxy:v0.0.71
  extraArgs:
    - '--log-level=info'
  config:
    serverGroups:
      - address: "prometheus-operated:9090"
        label: "k8s-1"
      - address: "prometheus.k8s-2.cloud.com"
        scheme: https
        label: "k8s-2"
      - address: "victoriametrics.cloud.com"
        label: "k8s-itdpl"
    detailedConfig: ""
  resources:
    limits:
      cpu: 150m
      memory: 256Mi
    requests:
      cpu: 50m
      memory: 128Mi
  port: 9090
  configmapReload:
    install: true
    image: jimmidyson/configmap-reload:v0.5.0
    resources:
      limits:
        cpu: 10m
        memory: 20Mi
      requests:
        cpu: 5m
        memory: 3Mi
  labels:
    label.key: label-value
  annotations:
    annotation.key: annotation-value
  priorityClassName: priority-class
```

[Back to TOC](#table-of-content)

### pushgateway

<!-- markdownlint-disable line-length -->
| Field             | Description                                                                                                                                                                                                                                                                                                                                                | Scheme                                                                                                                                 |
| ----------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------- |
| install           | Allow to disable create Pushgateway during deploy.                                                                                                                                                                                                                                                                                                         | boolean                                                                                                                                |
| image             | Image of pushgateway.                                                                                                                                                                                                                                                                                                                                      | string                                                                                                                                 |
| replicas          | Number of created pods.                                                                                                                                                                                                                                                                                                                                    | int                                                                                                                                    |
| paused            | Set paused to reconciliation.                                                                                                                                                                                                                                                                                                                              | boolean                                                                                                                                |
| extraArgs         | Additional pushgateway container arguments.                                                                                                                                                                                                                                                                                                                | list[string]                                                                                                                           |
| port              | Port for pushgateway deployment and service.                                                                                                                                                                                                                                                                                                               | integer                                                                                                                                |
| serviceMonitor    | Service monitor for pulling metrics.                                                                                                                                                                                                                                                                                                                       | *[Monitor](#monitor)                                                                                                                   |
| ingress           | Ingress allows to create Ingress.                                                                                                                                                                                                                                                                                                                          | *[Ingress](#ingress)                                                                                                                   |
| resources         | The resources that describe the compute resource requests and limits for single pods.                                                                                                                                                                                                                                                                      | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#resourcerequirements-v1-core)           |
| nodeSelector      | Defines which nodes the pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\"                                                                                                                                                                                                                                         | map[string]string                                                                                                                      |
| securityContext   | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                                                                                                                                                         | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#podsecuritycontext-v1-core)              |
| tolerations       | Tolerations allow the pods to schedule onto nodes with matching taints.                                                                                                                                                                                                                                                                                    | []v1.Toleration                                                                                                                        |
| affinity                                            | If specified, the pod's scheduling constraints                                                                                                                                                                                      | *v1.Affinity                                                                                                                   |
| annotations       | Annotations is an unstructured key value map stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. They are not queryable and should be preserved when modifying objects. For more information, refer to [http://kubernetes.io/docs/user-guide/annotations](http://kubernetes.io/docs/user-guide/annotations) | map[string]string                                                                                                                      |
| labels            | The map of string keys and values that can be used to organize and categorize (scope and select) objects. May match selectors of replication controllers and services. For more information, refer to [http://kubernetes.io/docs/user-guide/labels](http://kubernetes.io/docs/user-guide/labels)                                                           | map[string]string                                                                                                                      |
| volumes           | Volumes allows configuration of additional volumes on the output StatefulSet definition. Volumes specified will be appended to other volumes that are generated as a result of StorageSpec objects.                                                                                                                                                        | [v1.Volume](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#volume-v1-core)                                       |
| volumeMounts      | VolumeMounts allows configuration of additional VolumeMounts on the output StatefulSet definition. VolumeMounts specified will be appended to other VolumeMounts in the prometheus container, that are generated as a result of StorageSpec objects.                                                                                                       | [v1.VolumeMount](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#volumemount-v1-core)                             |
| storage           | PVC spec for Pushgateway. If specified, also adds flags --persistence.file=/data/pushgateway.data and --persistence.interval=5m, creates volume and volumeMount with name "storage-volume" in the deployment.                                                                                                                                              | [v1.PersistentVolumeClaimSpec](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.22/#persistentvolumeclaimspec-v1-core) |
| priorityClassName | PriorityClassName assigned to the Pods to prevent them from evicting.                                                                                                                                                                                                                                                                                      | string                                                                                                                                 |
<!-- markdownlint-enable line-length -->

Parameter `extraArgs` can contain any flags that Pushgateway can handle.
You can find them in the table bellow (relevant for Pushgateway v1.4.1).

<!-- markdownlint-disable line-length -->
| Flag                               | Default  | Description                                                                                    |
| ---------------------------------- | -------- | ---------------------------------------------------------------------------------------------- |
| `-h`, `--help`                     |          | Show context-sensitive help.                                                                   |
| `--version`                        |          | Show application version.                                                                      |
| `--web.listen-address`             | :9091    | Address to listen on for the web interface, API, and telemetry.                                |
| `--web.telemetry-path`             | /metrics | Path under which to expose metrics.                                                            |
| `--web.external-url`               |          | The URL under which the Pushgateway is externally reachable.                                   |
| `--web.route-prefix`               |          | Prefix for the internal routes of web endpoints. Defaults to the path of `--web.external-url`. |
| `--web.enable-lifecycle`           | false    | Enable shutdown via HTTP request.                                                              |
| `--web.enable-admin-api`           | false    | Enable API endpoints for admin control actions.                                                |
| `--persistence.file`               |          | File to persist metrics. If empty, metrics are only kept in memory.                            |
| `--persistence.interval`           | 5m       | The minimum interval at which to write out the persistence file.                               |
| `--push.disable-consistency-check` | false    | Do not check consistency of pushed metrics. DANGEROUS.                                         |
| `--log.level`                      | info     | Only log messages with the given severity or above. One of: [debug, info, warn, error].        |
| `--log.format`                     | logfmt   | Output format of log messages. One of: [logfmt, json].                                         |
<!-- markdownlint-enable line-length -->

Example:

```yaml
pushgateway:
  install: true
  image: prom/pushgateway:v1.4.1
  replicas: 1
  paused: false
  extraArgs:
    - "--log.level=info"
  volumes: {}
  volumeMounts: {}
  storage: {}
  port: 9091
  serviceMonitor:
    install: true
    interval: 30s
    scrapeTimeout: 10s
    metricRelabelings: []
    relabelings: []
  ingress: {}
  nodeSelector: {}
  resources:
    limits:
      cpu: 200m
      memory: 50Mi
    requests:
      cpu: 100m
      memory: 30Mi
  securityContext:
    runAsUser: 2001
    fsGroup: 2001
  tolerations: []
  annotations: {}
  labels: {}
  priorityClassName: priority-class
```

[Back to TOC](#table-of-content)

### stackdriver

<!-- markdownlint-disable line-length -->
| Field         | Description                                                                                                                                                                            | Scheme                                                                                                                         |
| ------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------ |
| image         | Image of 'stackdriver-prometheus-sidecar'. This service is deploying as sidecar container to Prometheus pod and send metrics from Prometheus to GCO.                                   | string                                                                                                                         |
| projectId     | Identificator of project in Google Cloud.                                                                                                                                              | string                                                                                                                         |
| location      | Location where project is deployed in Google Cloud.                                                                                                                                    | string                                                                                                                         |
| cluster       | Name of Kubernetes cluster in Google Cloud which will be monitored.                                                                                                                    | string                                                                                                                         |
| metricsFilter | List of filters for metrics which will be sent to GCO. Filters use the same syntax as Prometheus instant vector selectors. Example: `["{__name__=~'prometheus_example_app_metric'}"]`. | *[Prometheus instant vector selectors](https://prometheus.io/docs/prometheus/latest/querying/basics/#instant-vector-selectors) |
<!-- markdownlint-enable line-length -->

Example:

```yaml
integration:
  ...
  stackdriver:
    image: gcr.io/stackdriver-prometheus/stackdriver-prometheus-sidecar:0.8.2
    projectId: cloud-bss-sandbox
    location: us-central1-c
    cluster: cluster
    metricsFilter:
      - "{job!='kubelet'}"
      - "{job!='kube-apiserver'}"
      - "{job!='node-exporter'}"
  ...
```

[Back to TOC](#table-of-content)

### stackdriver-exporter

<!-- markdownlint-disable line-length -->
| Field                               | Description                                                                                                                                                                                                                                                  | Scheme                                                                                                                       |
| ----------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ---------------------------------------------------------------------------------------------------------------------------- |
| install                             | Allows to disable deploy stackdriver-exporter.                                                                                                                                                                                                               | bool                                                                                                                         |
| replicas                            | Number of created pods.                                                                                                                                                                                                                                      | int                                                                                                                          |
| name                                | A deployment name for stackdriver-exporter                                                                                                                                                                                                                   | string                                                                                                                       |
| image                               | A docker image to use for stackdriver-exporter deployment                                                                                                                                                                                                    | string                                                                                                                       |
| imagePullPolicy                     | Image pull policy to use for stackdriver-exporter deployment                                                                                                                                                                                                 | string                                                                                                                       |
| resources                           | The resources that describe the compute resource requests and limits for single pods.                                                                                                                                                                        | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#resourcerequirements-v1-core) |
| serviceAccount.create               | Allow to disable create ServiceAccount during deploy                                                                                                                                                                                                         | bool                                                                                                                         |
| serviceAccount.name                 | Provide a name in place of stackdriver-exporter for ServiceAccount                                                                                                                                                                                           | string                                                                                                                       |
| serviceAccount.annotations          | Will add the provided map to the annotations for the created serviceAccount.                                                                                                                                                                                 | map[string]string                                                                                                            |
| service.type                        | Type of promitor-agent-scraper service.                                                                                                                                                                                                                      | string                                                                                                                       |
| service.port                        | Port of promitor-agent-scraper which use in service.                                                                                                                                                                                                         | integer                                                                                                                      |
| service.labels                      | Annotations set which will be created in service.                                                                                                                                                                                                            | map[string]string                                                                                                            |
| nodeSelector                        | Defines which nodes the pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\"                                                                                                                                           | map[string]string                                                                                                            |
| annotations                         | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value"                                       | map[string]string                                                                                                            |
| labels                              | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                                                             | map[string]string                                                                                                            |
| securityContext                     | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                                                           | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#podsecuritycontext-v1-core)    |
| tolerations                         | Tolerations allow the pods to schedule onto nodes with matching taints.                                                                                                                                                                                      | []v1.Toleration                                                                                                              |
| affinity                            | It specifies the pod's scheduling constraints. For more information, refer to [https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core) | *v1.Affinity                                                                                                                 |
| priorityClassName                   | PriorityClassName assigned to the Pods to prevent them from evicting.                                                                                                                                                                                        | string                                                                                                                       |
| serviceMonitor                      | ServiceMonitor holds configuration attributes for stackdriver-exporter.                                                                                                                                                                                      | object                                                                                                                       |
| stackdriver.projectId               | The Google Project ID to gather metrics for                                                                                                                                                                                                                  | string                                                                                                                       |
| stackdriver.serviceAccountSecret    | An existing secret which contains credentials.json                                                                                                                                                                                                           | string                                                                                                                       |
| stackdriver.serviceAccountSecretKey | Provide custom key for the existing secret to load credentials.json from                                                                                                                                                                                     | string                                                                                                                       |
| stackdriver.serviceAccountKey       | A service account key JSON file. Must be provided when no existing secret is used, in this case a new secret will be created holding this service account                                                                                                    | string                                                                                                                       |
| stackdriver.maxRetries              | Max number of retries that should be attempted on 503 errors from Stackdriver                                                                                                                                                                                | string                                                                                                                       |
| stackdriver.httpTimeout             | How long should Stackdriver_exporter wait for a result from the Stackdriver API                                                                                                                                                                              | string                                                                                                                       |
| stackdriver.maxBackoff              | Max time between each request in an exp backoff scenario                                                                                                                                                                                                     | string                                                                                                                       |
| stackdriver.backoffJitter           | The amount of jitter to introduce in an exp backoff scenario                                                                                                                                                                                                 | string                                                                                                                       |
| stackdriver.retryStatuses           | The HTTP statuses that should trigger a retry                                                                                                                                                                                                                | string                                                                                                                       |
| stackdriver.dropDelegatedProjects   | Drop metrics from attached projects and fetch `project_id` only                                                                                                                                                                                              | string                                                                                                                       |
| stackdriver.metrics.typePrefixes    | The prefixes to gather metrics for, we default to just CPU metrics                                                                                                                                                                                           | string                                                                                                                       |
| stackdriver.metrics.interval        | The frequency to request                                                                                                                                                                                                                                     | string                                                                                                                       |
| stackdriver.metrics.offset          | How far into the past to offset                                                                                                                                                                                                                              | string                                                                                                                       |
| stackdriver.metrics.ingestDelay     | Offset for the Google Stackdriver Monitoring Metrics interval into the past by the ingest delay from the metric's metadata.                                                                                                                                  | string                                                                                                                       |
<!-- markdownlint-enable line-length -->

Example:

<!-- markdownlint-disable line-length -->
```yaml
stackdriverExporter:
  install: true
  name: "stackdriver-exporter"
  replicas: 1
  image: "prometheuscommunity/stackdriver-exporter:v0.12.0"
  imagePullPolicy: IfNotPresent
  resources:
    requests:
      cpu: 100m
      memory: 128Mi
    limits:
      cpu: 100m
      memory: 128Mi
  securityContext:
    runAsUser: 2001
    fsGroup: 2001
  service:
    type: ClusterIP
    port: 9255
    annotations: {}
  stackdriver:
    projectId: "my-test-project"
    serviceAccountSecret: ""
    serviceAccountSecretKey: ""
    serviceAccountKey: |
      {
        "type": "service_account",
        "project_id": "dummy-acc",
        "private_key_id": "00000000000000000000000000000000000",
        "private_key": "",
        "client_email": "dummy@dummy-acc.iam.gserviceaccount.com",
        "client_id": "0000000000000000000000000",
        "auth_uri": "https://accounts.google.com/o/oauth2/auth",
        "token_uri": "https://oauth2.googleapis.com/token",
        "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
        "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/dummy%dummy-acc.iam.gserviceaccount.com"
      }
    maxRetries: 0
    httpTimeout: 10s
    maxBackoff: 5s
    backoffJitter: 1s
    retryStatuses: 503
    dropDelegatedProjects: false
    metrics:
      typePrefixes: 'compute.googleapis.com/instance/cpu,compute.googleapis.com/instance/disk,pubsub.googleapis.com/subscription'
      interval: '5m'
      offset: '0s'
      ingestDelay: false
  affinity: {}
  extraArgs: []
    - --monitoring.filters='pubsub.googleapis.com/subscription:resource.labels.subscription_id=monitoring.regex.full_match("us-west4.*my-team-subs.*")'
  nodeSelector: {}
  tolerations: []
  serviceAccount:
    create: true
    name: stackdriver-exporter
    annotations: {}
  serviceMonitor:
    install: true
    interval: 2m
    telemetryPath: /metrics
    labels:
      app.kubernetes.io/component: monitoring
    timeout: 30s
    relabelings: []
    metricRelabelings: []
  annotations: {}
  labels: {}
  priorityClassName: priority-class
```
<!-- markdownlint-enable line-length -->

[Back to TOC](#table-of-content)

### victoria-metrics

**NOTE:** Before clean installation on cloud with victoria-metrics stack [read this.](#cleanup-victoriametrics-objects)

<!-- markdownlint-disable line-length -->
| Field             | Description                                                                                                                                                                                                        | Scheme  |
| ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ------- |
| tlsEnabled        | TLS configuration is enabled/disabled. By default, it is disabled.                                                                                                                                                 | boolean |
| clusterIssuerName | Cluster issuer name for generating certificates. This is a mandatory field if `tlsEnabled` is set to `true` and `generateCerts.enabled` is set to `true` for individual Victoriametrics components.                | string  |
| vmReplicas        | Defines the number of replicas all Victoriametrics components must have in cluster mode. If specified, individual Victoriametrics components' replicas will be overridden by `vmReplicas`. Must be greater than 1. | integer |

#### vmoperator

<!-- markdownlint-disable line-length -->
| Field                            | Description                                                                                                                                                                                                            | Scheme                                                                                                                       |
| -------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| install                          | Allow to enable or disable deploy vmOperator via monitoring-operator.                                                                                                                                                  | boolean                                                                                                                      |
| paused                           | Set paused to reconciliation for vmOperator                                                                                                                                                                            | boolean                                                                                                                      |
| serviceMonitor                   | ServiceMonitor holds configuration attributes for vmoperator.                                                                                                                                                          | object                                                                                                                       |
| serviceMonitor.install           | Allow to enable or disable deploy vmOperator service monitor                                                                                                                                                           | boolean                                                                                                                      |
| serviceMonitor.interval          | Allow to change metrics scrape interval.                                                                                                                                                                               | string                                                                                                                       |
| serviceMonitor.scrapeTimeout     | Allow to change metrics scrape timeout. Note that scrapeTimeout must be less the interval                                                                                                                              | string                                                                                                                       |
| serviceMonitor.metricRelabelings | Set metricRelabelings for the ServiceMonitor, use to apply to samples for ingestion                                                                                                                                    | []*promv1.RelabelConfig                                                                                                      |
| serviceMonitor.relabelings       | Set relabelings for the ServiceMonitor, use to apply to samples before scraping                                                                                                                                        | []*promv1.RelabelConfig                                                                                                      |
| extraEnvs                        | Allow to set extra system environment variables for vmoperator, refer to [https://docs.victoriametrics.com/operator/vars.html](https://docs.victoriametrics.com/operator/vars.html)                                    | map[string]string                                                                                                            |
| image                            | A Docker image to deploy the vmoperator.                                                                                                                                                                               | string                                                                                                                       |
| containerSecurityContext         | Holds container-level security attributes. Default for Kubernetes, `containerSecurityContext: { capabilities.drop: ALL, allowPrivilegeEscalation: false }`.                                                            | [core/v1.SecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#securitycontext-v1-core)      |
| securityContext                  | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                     | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#podsecuritycontext-v1-core)    |
| resources                        | The resources that describe the compute resource requests and limits for single pods.                                                                                                                                  | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#resourcerequirements-v1-core) |
| nodeSelector                     | Define which Nodes the Pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\"                                                                                                      | map[string]string                                                                                                            |
| tolerations                      | Allow the pods to schedule onto nodes with matching taints.                                                                                                                                                            | []v1.Toleration                                                                                                              |
| affinity                                            | If specified, the pod's scheduling constraints                                                                                                                                                                                      | *v1.Affinity                                                                                                                   |
| annotations                      | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value" | map[string]string                                                                                                            |
| labels                           | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                       | map[string]string                                                                                                            |
| priorityClassName                | PriorityClassName assigned to the Pods to prevent them from evicting.                                                                                                                                                  | string                                                                                                                       |
| tlsConfig                        | TLS configuration for VMOperator. Must be specified if `victoriametrics.tlsEnabled` is set to `true`                                                                                                                   | [TLSConfig](#tls-config)                                                                                                     |
<!-- markdownlint-enable line-length -->

Example:

```yaml
victoriametrics:
  tlsEnabled: false
  vmOperator:
    install: true
    paused: false
    serviceMonitor:
      install: true
      interval: 30s
      scrapeTimeout: 30s
      metricRelabelings: []
      relabelings: []
    extraEnvs: {}
    securityContext:
      runAsUser: 2000
      fsGroup: 2000
    resources:
      limits:
        cpu: 400m
        memory: 200Mi
      requests:
        cpu: 200m
        memory: 100Mi
    priorityClassName: priority-class
```

If TLS is enabled for Victoriametrics:

```yaml
victoriametrics:
  tlsEnabled: true
  vmOperator:
    install: true
    paused: false
    serviceMonitor:
      install: true
      interval: 30s
      scrapeTimeout: 30s
      metricRelabelings: []
      relabelings: []
    extraEnvs: {}
    securityContext:
      runAsUser: 2000
      fsGroup: 2000
    resources:
      limits:
        cpu: 400m
        memory: 200Mi
      requests:
        cpu: 200m
        memory: 100Mi
    priorityClassName: priority-class
    tlsConfig:
      generateCerts:
        enabled: true
        duration: 365
        renewBefore: 15
        secretName: "vmoperator-tls-secret"
```

#### vmsingle

<!-- markdownlint-disable line-length -->
| Field                      | Description                                                                                                                                                                                                                                 | Scheme                                                                                                                       |
| -------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| install                    | Allow to enable or disable deploy vmSingle via monitoring-operator.                                                                                                                                                                         | boolean                                                                                                                      |
| image                      | A Docker image to deploy the vmoperator.                                                                                                                                                                                                    | string                                                                                                                       |
| resources                  | The resources that describe the compute resource requests and limits for single pods.                                                                                                                                                       | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#resourcerequirements-v1-core) |
| securityContext            | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000, runAsGroup: 2000 }`.                                                                                        | [*v1.PodSecurityContext](https://kubernetes.io/docs/tasks/configure-pod-container/security-context/)                         |
| containers                 | Containers property allows to inject additions sidecars or to patch existing containers. It can be useful for proxies, backup, etc.                                                                                                         | [[]v1.Container](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#container-v1-core)                     |
| securityContext.runAsUser  | Specifies that for any Containers in the Pod, all processes run with user ID.                                                                                                                                                               | [*v1.PodSecurityContext](https://kubernetes.io/docs/tasks/configure-pod-container/security-context/)                         |
| securityContext.runAsGroup | Specifies the primary group ID for all processes within any containers of the Pod.                                                                                                                                                          | [*v1.PodSecurityContext](https://kubernetes.io/docs/tasks/configure-pod-container/security-context/)                         |
| securityContext.fsGroup    | Specifies that all processes of the container are also part of the supplementary group ID. The owner for volume and any files created in that volume will be Group ID.                                                                      | [*v1.PodSecurityContext](https://kubernetes.io/docs/tasks/configure-pod-container/security-context/)                         |
| paused                     | Set paused to reconciliation for vmsingle                                                                                                                                                                                                   | boolean                                                                                                                      |
| tolerations                | Allow the pods to schedule onto nodes with matching taints.                                                                                                                                                                                 | []v1.Toleration                                                                                                              |
| nodeSelector               | Define which Nodes the Pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\"                                                                                                                           | map[string]string                                                                                                            |
| affinity                                            | If specified, the pod's scheduling constraints                                                                                                                                                                                      | *v1.Affinity                                                                                                                   |
| labels                     | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                                            | map[string]string                                                                                                            |
| annotations                | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value"                      | map[string]string                                                                                                            |
| ingress                    | Ingress allows to create Ingress.                                                                                                                                                                                                           | *[Ingress](#ingress)                                                                                                         |
| retentionPeriod            | RetentionPeriod for the stored metrics. [https://docs.victoriametrics.com/Single-server-VictoriaMetrics.html#retention](https://docs.victoriametrics.com/Single-server-VictoriaMetrics.html#retention)                                      | string                                                                                                                       |
| extraArgs                  | List of key-value with additional arguments for run vmsingle [https://docs.victoriametrics.com/#list-of-command-line-flags](https://docs.victoriametrics.com/#list-of-command-line-flags)                                                   | map[string]string                                                                                                            |
| extraEnvs                  | Allow to set extra system environment variables for vmsingle.                                                                                                                                                                               | map[string]string                                                                                                            |
| secrets                    | Secrets is a list of Secrets in the same namespace as the VMSingle object, which shall be mounted into the VMSingle Pods                                                                                                                    | []string                                                                                                                     |
| storageDataPath            | Disables spec.storage option and overrides arg for victoria-metrics binary --storageDataPath, its users responsibility to mount proper device into given path.                                                                              | string                                                                                                                       |
| storage                    | Storage is the definition of how storage will be used by the VMSingle by default it`s empty dir                                                                                                                                             | *v1.PersistentVolumeClaimSpec                                                                                                |
| storageMetadata            | Defines annotations and labels attached to PVC for given vmsingle CR                                                                                                                                                                        | *github.com/VictoriaMetrics/operator/api/operator/v1beta1.EmbeddedObjectMetadata                                             |
| volumes                    | Volumes allows configuration of additional volumes on the output deploy definition. Volumes specified will be appended to other volumes that are generated as a result of StorageSpec objects                                               | []v1.Volume                                                                                                                  |
| volumeMounts               | VolumeMounts allows configuration of additional VolumeMounts on the output deploy definition. VolumeMounts specified will be appended to other VolumeMounts in the vmsingle container that are generated as a result of StorageSpec objects | []v1.VolumeMount                                                                                                             |
| priorityClassName          | PriorityClassName assigned to the Pods to prevent them from evicting                                                                                                                                                                        | string                                                                                                                       |
| tlsConfig                  | TLS configuration for VMSingle. Must be specified if `victoriametrics.tlsEnabled` is set to `true`                                                                                                                                          | [TLSConfig](#tls-config)                                                                                                     |
<!-- markdownlint-enable line-length -->

Example:

```yaml
victoriametrics:
  vmSingle:
    install: true
    ingress:
      host: vmsingle.test.org
      install: true
    extraArgs:
      search.maxPointsPerTimeseries: "150000"
    paused: false
    retentionPeriod: 1d
    securityContext:
      runAsUser: 2000
      runAsGroup: 2000
      fsGroup: 2000
    resources:
      limits:
        cpu: 1000m
        memory: 2000Mi
      requests:
        cpu: 500m
        memory: 1000Mi
    storage:
      accessModes:
        - ReadWriteOnce
      resources:
        requests:
          storage: 1Gi
      selector:
        matchLabels:
          app.kubernetes.io/name: vmsingle
      storageClassName: openstack-cinder
```

If TLS is enabled for Victoriametrics:

```yaml
victoriametrics:
  tlsEnabled: true
  clusterIssuerName: dev-cluster-issuer
  vmSingle:
    install: true
    ingress:
      host: vmsingle.test.org
      install: true
    extraArgs:
      search.maxPointsPerTimeseries: "150000"
    paused: false
    retentionPeriod: 1d
    securityContext:
      runAsUser: 2000
      runAsGroup: 2000
      fsGroup: 2000
    resources:
      limits:
        cpu: 1000m
        memory: 2000Mi
      requests:
        cpu: 500m
        memory: 1000Mi
    storage:
      accessModes:
        - ReadWriteOnce
      resources:
        requests:
          storage: 1Gi
      selector:
        matchLabels:
          app.kubernetes.io/name: vmsingle
      storageClassName: openstack-cinder
    tlsConfig:
      generateCerts:
        enabled: true
        duration: 365
        renewBefore: 15
        secretName: "vmsingle-tls-secret"
```

#### vmagent

<!-- markdownlint-disable line-length -->
| Field                                         | Description                                                                                                                                                                                                                                                                                                                                  | Scheme                                                                                                                       |
| --------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| install                                       | Allows to enable or disable deploy vmagent.                                                                                                                                                                                                                                                                                                  | boolean                                                                                                                      |
| image                                         | A Docker image to deploy the vmagent.                                                                                                                                                                                                                                                                                                        | string                                                                                                                       |
| resources                                     | The resources that describe the compute resource requests and limits for single pods.                                                                                                                                                                                                                                                        | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#resourcerequirements-v1-core) |
| securityContext                               | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                                                                                                                                           | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#podsecuritycontext-v1-core)    |
| containers                                    | Containers property allows to inject additions sidecars or to patch existing containers. It can be useful for proxies, backup, etc.                                                                                                                                                                                                          | [[]v1.Container](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#container-v1-core)                     |
| paused                                        | Set paused to reconciliation for vmagent                                                                                                                                                                                                                                                                                                     | boolean                                                                                                                      |
| tolerations                                   | Tolerations allow the pods to schedule onto nodes with matching taints.                                                                                                                                                                                                                                                                      | []v1.Toleration                                                                                                              |
| nodeSelector                                  | Defines which nodes the pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\"                                                                                                                                                                                                                           | map[string]string                                                                                                            |
| affinity                                            | If specified, the pod's scheduling constraints                                                                                                                                                                                      | *v1.Affinity                                                                                                                   |
| labels                                        | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                                                                                                                                             | map[string]string                                                                                                            |
| annotations                                   | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value"                                                                                                                       | map[string]string                                                                                                            |
| ingress                                       | Ingress allows to create Ingress for the vmagent UI.                                                                                                                                                                                                                                                                                         | *[Ingress](#ingress)                                                                                                         |
| scrapeInterval                                | Defines how often scrape targets by default.                                                                                                                                                                                                                                                                                                 | string                                                                                                                       |
| additionalScrape                              | The additionalScrape allows you to monitor hosts outside the cloud (balancers, graylog server, jenkins...)                                                                                                                                                                                                                                   | *v1.SecretKeySelector, [Victoria Metrics Config](https://docs.victoriametrics.com/operator/additional-scrape.html)           |
| vmAgentExternalLabelName                      | Name of vmAgent external label used to denote VmAgent instance name. Defaults to the value of `vmagent`. External label will _not_ be added when value is set to empty string (`""`).                                                                                                                                                        | string                                                                                                                       |
| externalLabels                                | The labels to add to any time series scraped by vmagent. It doesn't affect metrics ingested directly by push API's                                                                                                                                                                                                                           | map[string]string                                                                                                            |
| remoteWriteSettings                           | Defines global settings for all remoteWrite urls. For more information, refer to [https://docs.victoriametrics.com/vmagent.html](https://docs.victoriametrics.com/vmagent.html)                                                                                                                                                              | *v1beta1.VMAgentRemoteWriteSettings                                                                                          |
| remoteWriteSettings.maxDiskUsagePerURL        | The maximum file-based buffer size in bytes at -remoteWrite.tmpDataPath for each -remoteWrite.url. Disk usage is unlimited if the value is set to 0. For more information, refer to [https://docs.victoriametrics.com/vmagent.html](https://docs.victoriametrics.com/vmagent.html)                                                           | string                                                                                                                       |
| remoteWriteSettings.showURL                   | Whether to show -remoteWrite.url in the exported metrics. It is hidden by default, since it can contain sensitive auth info                                                                                                                                                                                                                  | string                                                                                                                       |
| remoteWriteSettings.label                     | Optional labels in the form 'name=value' to add to all the metrics before sending them                                                                                                                                                                                                                                                       | string                                                                                                                       |
| remoteWrite                                   | RemoteWriteSpec defines the remote_write configuration for vmagent. The `remote_write` allows to transparently send samples to a long term storage. For more information, refer to [https://docs.victoriametrics.com/vmagent.html](https://docs.victoriametrics.com/vmagent.html)                                                            | []v1beta1.VMAgentRemoteWriteSpec                                                                                             |
| remoteWrite.url                               | URL of the endpoint to send samples to.                                                                                                                                                                                                                                                                                                      | string                                                                                                                       |
| remoteWrite.basicAuth                         | Allow an endpoint to authenticate over basic authentication.                                                                                                                                                                                                                                                                                 | *v1beta1.BasicAuth                                                                                                           |
| remoteWrite.basicAuth.createSecret            | Allow create secret for basic authentification automatically.                                                                                                                                                                                                                                                                                | map[string]string                                                                                                            |
| remoteWrite.basicAuth.createSecret.secretName | Name of the secret with which will be created.                                                                                                                                                                                                                                                                                               | string                                                                                                                       |
| remoteWrite.basicAuth.createSecret.username   | Username for basic authentification.                                                                                                                                                                                                                                                                                                         | string                                                                                                                       |
| remoteWrite.basicAuth.createSecret.password   | Password for basic authentification.                                                                                                                                                                                                                                                                                                         | string                                                                                                                       |
| remoteWrite.bearerTokenSecret                 | Optional bearer auth token to use for -remoteWrite.url.                                                                                                                                                                                                                                                                                      | *v1.SecretKeySelector                                                                                                        |
| remoteWrite.oauth2                            | Defines auth configuration.                                                                                                                                                                                                                                                                                                                  | *v1beta1.OAuth2                                                                                                              |
| remoteWrite.tlsConfig                         | Describes tls configuration for remote write target.                                                                                                                                                                                                                                                                                         | *v1beta1.OAuth2                                                                                                              |
| remoteWrite.urlRelabelConfig                  | ConfigMap with relabeling config which is applied to metrics before sending them to the corresponding -remoteWrite.url [https://github.com/cybozu-go/VictoriaMetrics-operator/blob/master/docs/relabeling.MD](https://github.com/cybozu-go/VictoriaMetrics-operator/blob/master/docs/relabeling.MD)                                          | *v1.ConfigMapKeySelector                                                                                                     |
| remoteWrite.inlineUrlRelabelConfig            | Defines relabeling config for remoteWriteURL, it can be defined at crd spec. [https://github.com/cybozu-go/VictoriaMetrics-operator/blob/master/docs/relabeling.MD](https://github.com/cybozu-go/VictoriaMetrics-operator/blob/master/docs/relabeling.MD)                                                                                    | []RelabelConfig                                                                                                              |
| urlRelabelConfig                              | ConfigMap with global relabel config -remoteWrite.relabelConfig. This relabeling is applied to all the collected metrics before sending them to remote storage. [https://github.com/cybozu-go/VictoriaMetrics-operator/blob/master/docs/relabeling.MD](https://github.com/cybozu-go/VictoriaMetrics-operator/blob/master/docs/relabeling.MD) | *v1.ConfigMapKeySelector                                                                                                     |
| inlineUrlRelabelConfig                        | Defines GlobalRelabelConfig for vmagent, can be defined directly at CRD. [https://github.com/cybozu-go/VictoriaMetrics-operator/blob/master/docs/relabeling.MD](https://github.com/cybozu-go/VictoriaMetrics-operator/blob/master/docs/relabeling.MD)                                                                                        | []RelabelConfig                                                                                                              |
| podMonitorNamespaceSelector                   | Defines Namespaces to be selected for VMPodScrape discovery.                                                                                                                                                                                                                                                                                 | *metav1.LabelSelector                                                                                                        |
| serviceMonitorNamespaceSelector               | Defines Namespaces to be selected for VMServiceScrape discovery.                                                                                                                                                                                                                                                                             | *metav1.LabelSelector                                                                                                        |
| podMonitorSelector                            | Defines PodScrapes to be selected for target discovery.                                                                                                                                                                                                                                                                                      | *metav1.LabelSelector                                                                                                        |
| serviceMonitorSelector                        | Defines ServiceScrapes to be selected for target discovery.                                                                                                                                                                                                                                                                                  | *metav1.LabelSelector                                                                                                        |
| extraArgs                                     | List of key-value with additional arguments which will be passed to VMAgent pod [https://docs.victoriametrics.com/#list-of-command-line-flags](https://docs.victoriametrics.com/#list-of-command-line-flags)                                                                                                                                 | map[string]string                                                                                                            |
| extraEnvs                                     | Allow to set extra system environment variables for vmagent.                                                                                                                                                                                                                                                                                 | []v1.EnvVar                                                                                                                  |
| secrets                                       | Secrets is a list of Secrets in the same namespace as the VMSingle object, which shall be mounted into the vmagent pods                                                                                                                                                                                                                      | []string                                                                                                                     |
| volumes                                       | Volumes allows configuration of additional volumes on the output deploy definition. Volumes specified will be appended to other volumes that are generated as a result of StorageSpec objects                                                                                                                                                | []v1.Volume                                                                                                                  |
| volumeMounts                                  | VolumeMounts allows configuration of additional VolumeMounts on the output deploy definition. VolumeMounts specified will be appended to other VolumeMounts in the vmsingle container that are generated as a result of StorageSpec objects                                                                                                  | []v1.VolumeMount                                                                                                             |
| maxScrapeInterval                             | MaxScrapeInterval allows limiting maximum scrape interval for VMServiceScrape, VMPodScrape and other scrapes                                                                                                                                                                                                                                 | string                                                                                                                       |
| minScrapeInterval                             | MinScrapeInterval allows limiting minimal scrape interval for VMServiceScrape, VMPodScrape and other scrapes                                                                                                                                                                                                                                 | string                                                                                                                       |
| priorityClassName                             | PriorityClassName assigned to the Pods to prevent them from evicting                                                                                                                                                                                                                                                                         | string                                                                                                                       |
| tlsConfig                                     | TLS configuration for VMAgent. Must be specified if `victoriametrics.tlsEnabled` is set to `true`                                                                                                                                                                                                                                            | [TLSConfig](#tls-config)                                                                                                     |
<!-- markdownlint-enable line-length -->

Example:

```yaml
  vmAgent:
    install: true
    route:
      install: true
      host: vmagent.test.org
    logLevel: WARN
    serviceMonitor:
      interval: 29s
      scrapeTimeout: 9s
    extraArgs:
      loggerFormat: default
    remoteWriteSettings:
      label:
        project: development
        environment: qa-kuber-test
        cluster: qa-kuber-test
        team: test_team-test
      showURL: true
    remoteWrite:
      - url: https://remotehost-1:8429/api/v1/write
        basicAuth:
          username:
            key: username
            name: vmsecret
          password:
            key: password
            name: vmsecret
        tlsConfig:
          insecureSkipVerify: true
        urlRelabelConfig:
          name: "vmagent-relabel"
          key: "target-1-relabel.yaml"          
        inlineUrlRelabelConfig:
          - source_labels: [__name__]
            separator: ;
            regex: kube.*
            replacement: $1
            action: keep
      - url: https://remotehost-2:8429/api/v1/write
        basicAuth:
          createSecret:
            username: user
            password: pass
            secretName: vmsecret-basicauth
        tlsConfig:
          insecureSkipVerify: true
        urlRelabelConfig:
          name: "vmagent-relabel"
          key: "target-1-relabel.yaml"          
        inlineUrlRelabelConfig:
          - source_labels: [__name__]
            separator: ;
            regex: kube.*
            replacement: $1
            action: keep
    nodeSelector:
      role: compute
    labels:
      label: label_vmagent
    annotations:
      annotation: test_vmagent
    secrets:
      - kube-etcd-client-certs
```

If TLS is enabled for Victoriametrics:

```yaml
victoriametrics:
  tlsEnabled: true
  clusterIssuerName: dev-cluster-issuer
  vmAgent:
    install: true
    route:
      install: true
      host: vmagent.test.org
    logLevel: WARN
    serviceMonitor:
      interval: 29s
      scrapeTimeout: 9s
    extraArgs:
      loggerFormat: default
    remoteWriteSettings:
      label:
        project: development
        environment: qa-kuber-test
        cluster: qa-kuber-test
        team: test_team-test
      showURL: true
    remoteWrite:
      - url: https://remotehost-1:8429/api/v1/write
        basicAuth:
          username:
            key: username
            name: vmsecret
          password:
            key: password
            name: vmsecret
        tlsConfig:
          insecureSkipVerify: true
        urlRelabelConfig:
          name: "vmagent-relabel"
          key: "target-1-relabel.yaml"
        inlineUrlRelabelConfig:
          - source_labels: [__name__]
            separator: ;
            regex: kube.*
            replacement: $1
            action: keep
      - url: https://remotehost-2:8429/api/v1/write
        basicAuth:
          createSecret:
            username: user
            password: pass
            secretName: vmsecret-basicauth
        tlsConfig:
          insecureSkipVerify: true
        urlRelabelConfig:
          name: "vmagent-relabel"
          key: "target-1-relabel.yaml"          
        inlineUrlRelabelConfig:
          - source_labels: [__name__]
            separator: ;
            regex: kube.*
            replacement: $1
            action: keep
    nodeSelector:
      role: compute
    labels:
      label: label_vmagent
    annotations:
      annotation: test_vmagent
    secrets:
      - kube-etcd-client-certs
    tlsConfig:
      generateCerts:
        enabled: false
      createSecret:
        secretName: "vmagent-tls-secret"
        ca: |-
          -----BEGIN CERTIFICATE-----
          ...
          -----END CERTIFICATE-----
        key: |-
          -----BEGIN PRIVATE KEY-----
          ...
          -----END PRIVATE KEY-----
        cert: |-
          -----BEGIN CERTIFICATE-----
          ...
          ------END CERTIFICATE-----
```

#### vmalertmanager

<!-- markdownlint-disable line-length -->
| Field                         | Description                                                                                                                                                                                                                                                                                        | Scheme                                                                                                                       |
| ----------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| install                       | Allows to enable or disable deploy vmalertmanager.                                                                                                                                                                                                                                                 | boolean                                                                                                                      |
| image                         | A Docker image to deploy the vmalertmanager.                                                                                                                                                                                                                                                       | string                                                                                                                       |
| ingress                       | Ingress allows to create Ingress for the vmalertmanager UI.                                                                                                                                                                                                                                        | *[Ingress](#ingress)                                                                                                         |
| secrets                       | Secrets is a list of Secrets in the same namespace as the VMSingle object, which shall be mounted into the vmalertmanager pods                                                                                                                                                                     | []string                                                                                                                     |
| configRawYaml                 | Raw configuration for vmalertmanager, it helps it to start without secret.                                                                                                                                                                                                                         | string                                                                                                                       |
| configSecret                  | The name of a Kubernetes Secret in the same namespace as the VMAlertmanager object, which contains configuration for this VMAlertmanager, configuration must be inside secret key: alertmanager.yaml.                                                                                              | string                                                                                                                       |
| retention                     | Time duration VMAlertmanager shall retain data for. Default is '120h'                                                                                                                                                                                                                              | string                                                                                                                       |
| paused                        | Set paused to reconciliation for vmalertmanager                                                                                                                                                                                                                                                    | boolean                                                                                                                      |
| nodeSelector                  | Defines which nodes the pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\"                                                                                                                                                                                 | map[string]string                                                                                                            |
| resources                     | The resources that describe the compute resource requests and limits for single pods.                                                                                                                                                                                                              | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#resourcerequirements-v1-core) |
| affinity                      | It specifies the pod's scheduling constraints. For more information, refer to [https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core)                                       | *v1.Affinity                                                                                                                 |
| tolerations                   | Tolerations allow the pods to schedule onto nodes with matching taints.                                                                                                                                                                                                                            | []v1.Toleration                                                                                                              |
| securityContext               | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                                                                                                 | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#podsecuritycontext-v1-core)    |
| containers                    | Allows injecting additional containers or patching existing containers. This is meant to allow adding an authentication proxy to an VMAlertmanager pod.                                                                                                                                            | [[]v1.Container](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#container-v1-core)                     |
| selectAllByDefault            | Changes default behavior for empty CRD selectors, such ConfigSelector.                                                                                                                                                                                                                             | boolean                                                                                                                      |
| configSelector                | Defines selector for VMAlertmanagerConfig, result config will be merged with with Raw or Secret config.                                                                                                                                                                                            | *metav1.LabelSelector                                                                                                        |
| configNamespaceSelector       | Defines namespace selector for VMAlertmanagerConfig.                                                                                                                                                                                                                                               | *metav1.LabelSelector                                                                                                        |
| terminationGracePeriodSeconds | Defines period for container graceful termination.                                                                                                                                                                                                                                                 | string                                                                                                                       |
| labels                        | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                                                                                                   | map[string]string                                                                                                            |
| annotations                   | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value"                                                                             | map[string]string                                                                                                            |
| storage                       | Definition of how storage will be used by the VMAlertmanager                                                                                                                                                                                                                                       | *v1beta1.StorageSpec                                                                                                         |
| extraArgs                     | List of key-value with additional arguments which will be passed to VMAlertmanager pod [https://docs.victoriametrics.com/#list-of-command-line-flags](https://docs.victoriametrics.com/#list-of-command-line-flags)                                                                                | map[string]string                                                                                                            |
| extraEnvs                     | Allow to set extra system environment variables for VMAlertmanager.                                                                                                                                                                                                                                | map[string]string                                                                                                            |
| volumes                       | Volumes allows configuration of additional volumes on the output deploy definition. Volumes specified will be appended to other volumes that are generated as a result of StorageSpec objects                                                                                                      | []v1.Volume                                                                                                                  |
| volumeMounts                  | VolumeMounts allows configuration of additional VolumeMounts on the output deploy definition. VolumeMounts specified will be appended to other VolumeMounts in the vmsingle container that are generated as a result of StorageSpec objects                                                        | []v1.VolumeMount                                                                                                             |
| priorityClassName             | PriorityClassName assigned to the Pods to prevent them from evicting                                                                                                                                                                                                                               | string                                                                                                                       |
| tlsConfig                     | TLS configuration for VMAlertManager. Must be specified if `victoriametrics.tlsEnabled` is set to `true`                                                                                                                                                                                           | [TLSConfig](#tls-config)                                                                                                     |
| webConfig                     | Web configuration for VMAlertManager. Parameter is optional and this configuration is auto-populated by the operator. However, it can be used to customize/override TLS settings. More details [here](https://github.com/prometheus/alertmanager/blob/main/docs/https.md#https-and-authentication) | object                                                                                                                       |
| gossipConfig                  | Gossip configuration for VMAlertManager. Can be used to specify whether to use mutual TLS for gossip. More details [here](https://github.com/prometheus/alertmanager/blob/main/docs/https.md#gossip-traffic)                                                                                       | object                                                                                                                       |
<!-- markdownlint-enable line-length -->

Example:

```yaml
  vmAlertManager:
    install: true
    ingress:
      host: vmalertmanager.test.org
      install: true
    paused: false
    selectAllByDefault: true
    securityContext:
      runAsUser: 2001
      fsGroup: 2001
    secrets:
      - kube-etcd-client-certs
    volumes: {}
    volumeMounts: {}
```

If TLS is enabled for Victoriametrics:

```yaml
victoriametrics:
  tlsEnabled: true
  vmAlertManager:
    install: true
    ingress:
      host: vmalertmanager.test.org
      install: true
    paused: false
    selectAllByDefault: true
    securityContext:
      runAsUser: 2001
      fsGroup: 2001
    secrets:
      - kube-etcd-client-certs
    volumes: {}
    volumeMounts: {}
    tlsConfig:
      existingSecret: "vmalertmanager-tls-secret"
    webConfig: # optional
      tls_server_config:
        cert_secret_ref:
          key: tls.crt
          name: vmalertmanager-tls-secret
        key_secret_ref:
          key: tls.key
          name: vmalertmanager-tls-secret
```

**Note**
In above case, it is assumed that there is already a secret with the name `vmalertmanager-tls-secret`
present in the server with valid certificates.

#### vmalert

<!-- markdownlint-disable line-length -->
| Field                         | Description                                                                                                                                                                                                                                                  | Scheme                                                                                                                       |
| ----------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ---------------------------------------------------------------------------------------------------------------------------- |
| install                       | Allows to enable or disable deploy vmalert.                                                                                                                                                                                                                  | boolean                                                                                                                      |
| image                         | A Docker image to deploy the vmalert.                                                                                                                                                                                                                        | string                                                                                                                       |
| ingress                       | Ingress allows to create Ingress for the vmalert UI.                                                                                                                                                                                                         | *[Ingress](#ingress)                                                                                                         |
| secrets                       | Secrets is a list of Secrets in the same namespace as the VMAlert object, which shall be mounted into the vmalert pods                                                                                                                                       | []string                                                                                                                     |
| resources                     | The resources that describe the compute resource requests and limits for single pods.                                                                                                                                                                        | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#resourcerequirements-v1-core) |
| affinity                      | It specifies the pod's scheduling constraints. For more information, refer to [https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core) | *v1.Affinity                                                                                                                 |
| tolerations                   | Tolerations allow the pods to schedule onto nodes with matching taints.                                                                                                                                                                                      | []v1.Toleration                                                                                                              |
| securityContext               | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                                                           | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#podsecuritycontext-v1-core)    |
| containers                    | Containers property allows to inject additions sidecars or to patch existing containers. It can be useful for proxies, backup, etc.                                                                                                                          | [[]v1.Container](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#container-v1-core)                     |
| evaluationInterval            | How often evalute rules by default.                                                                                                                                                                                                                          | string                                                                                                                       |
| selectAllByDefault            | Changes default behavior for empty CRD selectors, such RuleSelector.                                                                                                                                                                                         | boolean                                                                                                                      |
| ruleSelector                  | Defines which VMRules to mount for loading alerting                                                                                                                                                                                                          | *metav1.LabelSelector                                                                                                        |
| ruleNamespaceSelector         | Defines namespace selector for VMRules discovery.                                                                                                                                                                                                            | *metav1.LabelSelector                                                                                                        |
| port                          | Port for listen.                                                                                                                                                                                                                                             | string                                                                                                                       |
| remoteWrite                   | Optional URL to remote-write compatible storage to persist vmalert state and rule results to.                                                                                                                                                                | []v1beta1.VMAlertRemoteWriteSpec                                                                                             |
| remoteWrite.url               | URL of the endpoint to send samples to.                                                                                                                                                                                                                      | string                                                                                                                       |
| remoteWrite.flushInterval     | Defines interval of flushes to remote write endpoint (default 5s).                                                                                                                                                                                           | string                                                                                                                       |
| remoteWrite.basicAuth         | Allow an endpoint to authenticate over basic authentication.                                                                                                                                                                                                 | *v1beta1.BasicAuth                                                                                                           |
| remoteWrite.tlsConfig         | Describes tls configuration for remote write target.                                                                                                                                                                                                         | *v1beta1.TLSConfig                                                                                                           |
| remoteRead                    | Optional URL to read vmalert state (persisted via RemoteWrite) This configuration only makes sense if alerts state has been successfully persisted (via RemoteWrite) before.                                                                                 | []v1beta1.VMAlertRemoteReadSpec                                                                                              |
| remoteRead.url                | URL of the endpoint to send samples to.                                                                                                                                                                                                                      | string                                                                                                                       |
| remoteRead.basicAuth          | Allow an endpoint to authenticate over basic authentication.                                                                                                                                                                                                 | *v1beta1.BasicAuth                                                                                                           |
| remoteRead.tlsConfig          | Describes tls configuration for remote write target.                                                                                                                                                                                                         | *v1beta1.TLSConfig                                                                                                           |
| rulePath                      | RulePath to the file with alert rules. Supports patterns. Flag can be specified multiple times.                                                                                                                                                              | []string                                                                                                                     |
| notifier                      | Notifier prometheus alertmanager endpoint spec.                                                                                                                                                                                                              | *v1beta1.VMAlertNotifierSpec                                                                                                 |
| notifier.url                  | AlertManager url.                                                                                                                                                                                                                                            | string                                                                                                                       |
| notifier.basicAuth            | Allow notifier to authenticate over basic authentication.                                                                                                                                                                                                    | *v1beta1.BasicAuth                                                                                                           |
| notifier.tlsConfig            | Describes describes tls configuration for notifier.                                                                                                                                                                                                          | *v1beta1.TLSConfig                                                                                                           |
| notifier.selector             | Allows service discovery for alertmanager in this case all matched vmalertmanager replicas will be added into vmalert notifier.url as statefulset pod.fqdn                                                                                                   | *metav1.LabelSelector                                                                                                        |
| notifiers                     | Notifier prometheus alertmanager endpoint specs.                                                                                                                                                                                                             | []v1beta1.VMAlertNotifierSpec                                                                                                |
| notifierConfigRef             | NotifierConfigRef reference for secret with notifier configuration for vmalert. Only one of notifier options could be chosen: notifierConfigRef or notifiers +  notifier.                                                                                    | *v1.SecretKeySelector                                                                                                        |
| datasource                    | Datasource Victoria Metrics.                                                                                                                                                                                                                                 | *vmetricsv1b1.VMAlertDatasourceSpec                                                                                          |
| datasource.url                | Victoria Metrics url.                                                                                                                                                                                                                                        | string                                                                                                                       |
| datasource.basicAuth          | Allow datasource to authenticate over basic authentication.                                                                                                                                                                                                  | *v1beta1.BasicAuth                                                                                                           |
| datasource.tlsConfig          | Describes describes tls configuration for datasource target.                                                                                                                                                                                                 | *v1beta1.TLSConfig                                                                                                           |
| paused                        | Set paused to reconciliation for vmalert                                                                                                                                                                                                                     | boolean                                                                                                                      |
| terminationGracePeriodSeconds | Defines period for container graceful termination.                                                                                                                                                                                                           | string                                                                                                                       |
| externalLabels                | The labels to add to any time series. Example: ["team=monitoring", "project=development", "environment=dr311qa", "cluster=dr311qa"]                                                                                                                          | list                                                                                                                         |
| nodeSelector                  | Defines which nodes the pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\"                                                                                                                                           | map[string]string                                                                                                            |
| labels                        | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                                                             | map[string]string                                                                                                            |
| annotations                   | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value"                                       | map[string]string                                                                                                            |
| extraArgs                     | List of key-value with additional arguments which will be passed to VMAgent pod [https://docs.victoriametrics.com/#list-of-command-line-flags](https://docs.victoriametrics.com/#list-of-command-line-flags)                                                 | map[string]string                                                                                                            |
| extraEnvs                     | Allow to set extra system environment variables for vmagent.                                                                                                                                                                                                 | map[string]string                                                                                                            |
| volumes                       | Volumes allows configuration of additional volumes on the output deploy definition. Volumes specified will be appended to other volumes that are generated as a result of StorageSpec objects                                                                | []v1.Volume                                                                                                                  |
| volumeMounts                  | VolumeMounts allows configuration of additional VolumeMounts on the output deploy definition. VolumeMounts specified will be appended to other VolumeMounts in the vmsingle container that are generated as a result of StorageSpec objects                  | []v1.VolumeMount                                                                                                             |
| priorityClassName             | PriorityClassName assigned to the Pods to prevent them from evicting                                                                                                                                                                                         | string                                                                                                                       |
| tlsConfig                     | TLS configuration for VMAlert. Must be specified if `victoriametrics.tlsEnabled` is set to `true`                                                                                                                                                            | [TLSConfig](#tls-config)                                                                                                     |
<!-- markdownlint-enable line-length -->

```yaml
  vmAlert:
    install: true
    ingress:
      host: vmalert.test.org
      install: true
    paused: false
    selectAllByDefault: true
    datasource:
      url: http://vmsingle-vmsingle.monitoring:8429
    notifier:
      url: http://vmalertmanager-vmalertmanager.monitoring:9093
    remoteWrite:
      url: http://vmsingle-vmsingle.monitoring:8429
    securityContext:
      runAsUser: 2001
      fsGroup: 2001
    secrets:
      - kube-etcd-client-certs
    volumes: {}
    volumeMounts: {}
```

If TLS is enabled for Victoriametrics:

```yaml
victoriametrics:
  tlsEnabled: true
  clusterIssuerName: dev-cluster-issuer
  vmAlert:
    install: true
    ingress:
      host: vmalert.test.org
      install: true
    paused: false
    selectAllByDefault: true
    datasource:
      url: http://vmsingle-vmsingle.monitoring:8429
    notifier:
      url: http://vmalertmanager-vmalertmanager.monitoring:9093
    remoteWrite:
      url: http://vmsingle-vmsingle.monitoring:8429
    securityContext:
      runAsUser: 2001
      fsGroup: 2001
    secrets:
      - kube-etcd-client-certs
    volumes: {}
    volumeMounts: {}
    tlsConfig:
      generateCerts:
        enabled: true
        duration: 365
        renewBefore: 15
        secretName: "vmalert-tls-secret"
```

#### vmauth

<!-- markdownlint-disable line-length -->
| Field                         | Description                                                                                                                                                                                                                                                  | Scheme                                                                                                                       | Required |
| ----------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ---------------------------------------------------------------------------------------------------------------------------- | -------- |
| install                       | Allows to enable or disable deploy vmauth.                                                                                                                                                                                                                   | boolean                                                                                                                      | false    |
| image                         | A Docker image to deploy the vmauth.                                                                                                                                                                                                                         | string                                                                                                                       | false    |
| paused                        | Set paused to reconciliation for vmauth.                                                                                                                                                                                                                     | boolean                                                                                                                      | false    |
| ingress                       | Allows Ingress configuration for the vmauth.                                                                                                                                                                                                                 | *v1beta1.EmbeddedIngress                                                                                                     | false    |
| secrets                       | Secrets is a list of Secrets in the same namespace as the VMAuth object, which shall be mounted into the vmauth pods                                                                                                                                         | []string                                                                                                                     | false    |
| configMaps                    | ConfigMaps is a list of ConfigMaps in the same namespace as the VMAuth object, which shall be mounted into the VMAuth Pods.                                                                                                                                  | []string                                                                                                                     | false    |
| volumes                       | Volumes allows configuration of additional volumes on the output deploy definition. Volumes specified will be appended to other volumes that are generated as a result of StorageSpec objects                                                                | []v1.Volume                                                                                                                  | false    |
| volumeMounts                  | VolumeMounts allows configuration of additional VolumeMounts on the output deploy definition. VolumeMounts specified will be appended to other VolumeMounts in the vmauth container that are generated as a result of StorageSpec objects                    | []v1.VolumeMount                                                                                                             | false    |
| resources                     | The resources that describe the compute resource requests and limits for single pods.                                                                                                                                                                        | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#resourcerequirements-v1-core) | false    |
| affinity                      | It specifies the pod's scheduling constraints. For more information, refer to [https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core) | *v1.Affinity                                                                                                                 | false    |
| tolerations                   | Tolerations allow the pods to schedule onto nodes with matching taints.                                                                                                                                                                                      | []v1.Toleration                                                                                                              | false    |
| securityContext               | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                                                           | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#podsecuritycontext-v1-core)    | false    |
| containers                    | Containers property allows to inject additions sidecars or to patch existing containers. It can be useful for proxies, backup, etc.                                                                                                                          | [[]v1.Container](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#container-v1-core)                     | false    |
| port                          | Port for listen.                                                                                                                                                                                                                                             | string                                                                                                                       | false    |
| selectAllByDefault            | Changes default behavior for empty CRD selectors, such UserSelector. Default - true.                                                                                                                                                                         | boolean                                                                                                                      | false    |
| userSelector                  | Defines VMUser to be selected for config file generation                                                                                                                                                                                                     | *metav1.LabelSelector                                                                                                        | false    |
| userNamespaceSelector         | Defines namespace selector for VMAuth discovery.                                                                                                                                                                                                             | *metav1.LabelSelector                                                                                                        | false    |
| extraArgs                     | List of key-value with additional arguments which will be passed to VMAuth pod [https://docs.victoriametrics.com/#list-of-command-line-flags](https://docs.victoriametrics.com/#list-of-command-line-flags)                                                  | map[string]string                                                                                                            | false    |
| extraEnvs                     | Allow to set extra system environment variables for vmauth.                                                                                                                                                                                                  | map[string]string                                                                                                            | false    |
| labels                        | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                                                             | map[string]string                                                                                                            | false    |
| annotations                   | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value"                                       | map[string]string                                                                                                            | false    |
| nodeSelector                  | Defines which nodes the pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\"                                                                                                                                           | map[string]string                                                                                                            | false    |
| terminationGracePeriodSeconds | Defines period for container graceful termination.                                                                                                                                                                                                           | string                                                                                                                       | false    |
| priorityClassName             | PriorityClassName assigned to the Pods to prevent them from evicting                                                                                                                                                                                         | string                                                                                                                       | false    |
| tlsConfig                     | TLS configuration for VMAuth. Must be specified if `victoriametrics.tlsEnabled` is set to `true`                                                                                                                                                             | [TLSConfig](#tls-config)                                                                                                     | false    |
<!-- markdownlint-enable line-length -->

```yaml
  vmAuth:
    install: true
    ingress:
      host: vmauth.test.org
    paused: false
    secretName: vmauth-secret
    extraVarsSecret:
      pass: vmauth
    securityContext:
      runAsUser: 2000
      fsGroup: 2000
```

If TLS is enabled for Victoriametrics:

```yaml
victoriametrics:
  tlsEnabled: true
  clusterIssuerName: dev-cluster-issuer
  vmAuth:
    install: true
    ingress:
      host: vmauth.test.org
    paused: false
    secretName: vmauth-secret
    extraVarsSecret:
      pass: vmauth
    securityContext:
      runAsUser: 2000
      fsGroup: 2000
    tlsConfig:
      generateCerts:
        enabled: true
        duration: 365
        renewBefore: 15
        secretName: "vmauth-tls-secret"
```

#### vmuser

<!-- markdownlint-disable line-length -->
| Field            | Description                                                                                  | Scheme                | Required |
| ---------------- | -------------------------------------------------------------------------------------------- | --------------------- | -------- |
| install          | Allows to enable or disable deploy vmuser.                                                   | boolean               | false    |
| image            | A Docker image to deploy the vmuser.                                                         | string                | false    |
| paused           | Set paused to reconciliation for vmuser.                                                     | boolean               | false    |
| username         | UserName basic auth username for accessing protected endpoint. Default value: `admin`        | *string               | false    |
| password         | Password basic auth password for accessing protected endpoint. Default value: `admin`        | *string               | false    |
| passwordRef      | PasswordRef allows fetching password from user-create secret by its name and key.            | *v1.SecretKeySelector | false    |
| tokenRef         | TokenRef allows fetching token from user-created secrets by its name and key.                | *v1.SecretKeySelector | false    |
| generatePassword | GeneratePassword instructs operator to generate password for user if spec.password if empty. | *v1.SecretKeySelector | false    |
| bearerToken      | BearerToken Authorization header value for accessing protected endpoint.                     | *string               | false    |
| targetRefs       | TargetRefs - reference to endpoints, which user may access.                                  | []v1beta1.TargetRef   | false    |
<!-- markdownlint-enable line-length -->

```yaml
  vmUser:
    install: true
    paused: false
    username: vmauth
    passwordRef:
      key: pass
      name: vmauth-secret
```

If targetRefs is empty - targetRefs for vmalert, vmalertmanager, vmagent and vmsingle will be automatically created.
Default paths for vmalert:

* /api/v1/rules
* /api/v1/alerts
* /api/v1/alert.*
* /vmalert.*

Default paths for vmalertmanager:

* /api/v2/alerts.*
* /api/v2/receivers.*
* /api/v2/silences.*
* /api/v2/status.*

**TODO:** At this moment, you can navigate to the VMAlertManager UI through its Ingress, without using VMAuth, because
VMAlertManager is unavailable via VMAuth Ingress. We need to configure separated routing for it based on another Ingress
or on changing the UI with references to VM components.

Default paths for vmagent:

* /target.*
* /service-discovery.*
* /api/v1/write
* /api/v1/import.*
* /api/v1/target.*

Default paths for vmsingle:

* /vmui.*
* /graph.*
* /api/v1/label.*
* /api/v1/query.*
* /api/v1/rules
* /api/v1/alerts
* /api/v1/metadata.*
* /api/v1/format.*
* /api/v1/series.*
* /api/v1/status.*
* /api/v1/export.*
* /api/v1/admin/tsdb.*
* /prometheus/graph.*
* /prometheus/api/v1/label.*
* /graphite.*
* /prometheus/api/v1/query.*
* /prometheus/api/v1/rules
* /prometheus/api/v1/alerts
* /prometheus/api/v1/metadata
* /prometheus/api/v1/series.*
* /prometheus/api/v1/status.*
* /prometheus/api/v1/export.*
* /prometheus/federate
* /prometheus/api/v1/admin/tsdb.*

In case of OAuth + vmauth `Basic Auth` is supported only.
If OAuth is going to be installed - username (field "username: vmauth" in the example above) is to match the OAuth one,
auth `basicAuthPwd` field is to contain vmauth password.
Here is the config example.

> **Note:**
>
> If you skip parameters `username` and `password` they should be automatically set with default values `admin/admin`.

```yaml
   vmAuth:
     install: true
     ingress:
       install: true
       host: vmauth.test.org
     paused: false
     secretName: vmauth-secret
     extraVarsSecret:
       pass: vmauth
     securityContext:
       runAsUser: 2000
       fsGroup: 2000
   vmUser:
     install: true
     paused: false
     username: test
     passwordRef:
       key: pass
       name: vmauth-secret
 auth:
   clientId: vmtest
   clientSecret: vmtest
   loginUrl: https://idp.test.org/realms/myrealm/protocol/openid-connect/auth
   tokenUrl: https://idp.test.org/realms/myrealm/protocol/openid-connect/token
   userInfoUrl: https://idp.test.org/realms/myrealm/protocol/openid-connect/userinfo
   tlsConfig:
     insecureSkipVerify: true
   basicAuthPwd: vmauth
 oAuthProxy:
   image: oauth2-proxy/oauth2-proxy:v7.4.0
```

Default paths for vmselect:

* `/select/.*`

[Back to TOC](#table-of-content)

#### vmcluster

**Note:**
Currently, `vmCluster` support is not available.

<!-- markdownlint-disable line-length -->
| Field              | Description                                                                                         | Scheme                                                                        | Required |
| ------------------ | --------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------- | -------- |
| install            | Allows to enable or disable deploy vmcluster.                                                       | boolean                                                                       | false    |
| vmSelectImage      | A Docker image to deploy the vmselect.                                                              | string                                                                        | false    |
| vmStorageImage     | A Docker image to deploy the vmstorage.                                                             | string                                                                        | false    |
| vmInsertImage      | A Docker image to deploy the vminsert.                                                              | string                                                                        | false    |
| paused             | Set paused to reconciliation for vmcluster.                                                         | boolean                                                                       | false    |
| vmSelect           | Parameters for `vmselect`                                                                           | [v1beta1.VMSelect](https://docs.victoriametrics.com/operator/api/#vmselect)   | false    |
| vmStorage          | Parameters for `vmstorage`                                                                          | [v1beta1.VMStorage](https://docs.victoriametrics.com/operator/api/#vmstorage) | false    |
| vmInsert           | Parameters for `vminsert`                                                                           | [v1beta1.VMInsert](https://docs.victoriametrics.com/operator/api/#vminsert)   | false    |
| vmSelectIngress    | Ingress configuration for VMSelect.                                                                 | *[Ingress](#ingress)                                                          | false    |
| vmSelectTlsConfig  | TLS configuration for VMSelect. Must be specified if `victoriametrics.tlsEnabled` is set to `true`  | [TLSConfig](#tls-config)                                                      | false    |
| vmInsertTlsConfig  | TLS configuration for VMInsert. Must be specified if `victoriametrics.tlsEnabled` is set to `true`  | [TLSConfig](#tls-config)                                                      | false    |
| vmStorageTlsConfig | TLS configuration for VMStorage. Must be specified if `victoriametrics.tlsEnabled` is set to `true` | [TLSConfig](#tls-config)                                                      | false    |

<!-- markdownlint-enable line-length -->

```yaml
  vmCluster:
    install: true
    paused: false
    vmSelectImage: victoriametrics/vmselect:v1.101.0-cluster
    vmInsertImage: victoriametrics/vminsert:v1.101.0-cluster
    vmStorageImage: victoriametrics/vmstorage:v1.101.0-cluster
    vmSelect:
      replicaCount: 2
    vmSelectIngress:
      host: vmselect.test.org
      install: true
    vmStorage:
      replicaCount: 2
    vmInsert:
      replicaCount: 2
    vmSelectTlsConfig:
      generateCerts:
        enabled: true
        duration: 365
        renewBefore: 15
        secretName: "vmselect-tls-secret"
    vmInsertTlsConfig:
      generateCerts:
        enabled: true
        duration: 365
        renewBefore: 15
        secretName: "vminsert-tls-secret"
    vmStorageTlsConfig:
      generateCerts:
        enabled: true
        duration: 365
        renewBefore: 15
        secretName: "vmstorage-tls-secret"
```

**Note:**
`monitoring-operator` can install/manage ingress for `vmselect` when victoriametrics is installed in cluster mode.
Example configuration below:

```yaml
  vmCluster:
    install: true
    paused: false
    vmSelectImage: victoriametrics/vmselect:v1.101.0-cluster
    vmInsertImage: victoriametrics/vminsert:v1.101.0-cluster
    vmStorageImage: victoriametrics/vmstorage:v1.101.0-cluster
    vmSelect:
      replicaCount: 2
    vmSelectIngress:
      host: vmselect.test.org
      install: true
```

If necessary, `vminsert` and `vmstorage` ingresses can be installed manually.

[Back to TOC](#table-of-content)

#### cleanup

<!-- markdownlint-disable line-length -->
| Field                                 | Description                                                                                                                                                                                                                               | Scheme            | Required |
| ------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------- | -------- |
| cleanup.deleteCRs                     | Enable of disable pre-delete cleanup job that should remove victoriametrics-operator Custom Resources (CRs). This job allow to correct remove VM CRs with processing it remove by victoriametrics-operator.                               | bool              | yes      |
| cleanup.hook.name                     | The name of cleanup job                                                                                                                                                                                                                   | string            | no       |
| cleanup.hook.timeout                  | Timeout between scale and delete commands in cleanup job. can be specified without units or with units. If you want yo use units, check which units support the Linux command "sleep".                                                    | string            | no       |
| cleanup.hook.image                    | A docker image to use for cleanup job. Calculate automatically from deployDescriptor section, but this parameter allow to override the value.                                                                                             | string            | no       |
| cleanup.hook.affinity                 | If specified, the pod's scheduling constraints. More info: [https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#affinity-v1-core](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.28/#affinity-v1-core) | object            | no       |
| cleanup.hook.nodeSelector             | Allow define which Nodes the Pods are scheduled on                                                                                                                                                                                        | map[string]string | no       |
| cleanup.hook.tolerations              | Tolerations allow the pods to schedule onto nodes with matching taints                                                                                                                                                                    | object            | no       |
| cleanup.hook.securityContext          | SecurityContext holds pod-level security attributes                                                                                                                                                                                       | object            | no       |
| cleanup.hook.containerSecurityContext | SecurityContext holds container-level security attributes                                                                                                                                                                                 | object            | no       |
| cleanup.hook.resources                | The resources describes the compute resource requests and limits for single Pods. More info: [https://kubernetes.io/docs/user-guide/compute-resources/](https://kubernetes.io/docs/user-guide/compute-resources/)                         | object            | no       |
<!-- markdownlint-enable line-length -->

Example:

```yaml
victoriametrics:
  cleanup:
    deleteCRs: true

    hook:
      name: cleanup-hook
      image: bitnami/kubectl:1.28.10

      timeout: 10s

      affinity: {}
      nodeSelector: {}
      tolerations: {}

      securityContext: {}
      containerSecurityContext:
        capabilities:
          drop:
            - ALL
        allowPrivilegeEscalation: false

      resources:
        requests:
          cpu: 100m
          memory: 56Mi
        limits:
          cpu: 500m
          memory: 256Mi
```

[Back to TOC](#table-of-content)

#### TLS Config

TLSConfig holds SSL/TLS configuration attributes.
The parameters are required if SSL/TLS connection is required between Kubernetes cluster and Victoriametrics.
This section is applicable only if `victoriametrics.tlsEnabled` is set to `true`.

<!-- markdownlint-disable line-length -->
| Parameter                   | Type    | Mandatory | Default value  | Description                                                                                                                                                                                                                                                                                                                                                                                   |
| --------------------------- | ------- | --------- | -------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `existingSecret`            | string  | no        | `-`            | Name of the pre-existing secret that contains TLS configuration for Victoriametrics. If specified, `generateCerts.enabled` must be set to `false`. The `existingSecret` is expected to contain CA certificate, TLS key and TLS certificate in `ca.crt`, `tls.key` and `tls.crt` fields respectively. Example can be found for [vmalertmanager](#vmalertmanager)                               |
| `generateCerts.enabled`     | boolean | no        | `true`         | Generation of certificate is enabled by default. If `tlsConfig.existingSecret` is specified, `tlsConfig.generateCerts` section will be skipped. `cert-manager` will generate certificate with the name configured using `generateCerts.secretName`, if it doesn't exist already. Example can be found for [vmauth](#vmauth)                                                                   |
| `generateCerts.duration`    | integer | no        | `365`          | Duration in days, until which issued certificate will be valid.                                                                                                                                                                                                                                                                                                                               |
| `generateCerts.renewBefore` | integer | no        | `15`           | Number of days before which certificate must be renewed.                                                                                                                                                                                                                                                                                                                                      |
| `generateCerts.secretName`  | string  | no        | `*-tls-secret` | Name of the new secret that needs to be created for storing TLS configuration of Victoriametrics components.                                                                                                                                                                                                                                                                                  |
| `createSecret`              | object  | no        | `-`            | New secret with the name `tlsConfig.createSecret.secretName` will be created using already known certificate content. If `tlsConfig.existingSecret` is specified, `tlsConfig.createSecret` section will be skipped. If `createSecret` section is specified, then it is necessary to set `tlsConfig.generateCerts.enabled` to `false` explicitly. Example can be found for [vmagent](#vmagent) |
| `createSecret.ca`           | string  | no        | `-`            | Already known CA certificate will be added to newly created secret.                                                                                                                                                                                                                                                                                                                           |
| `createSecret.key`          | string  | no        | `-`            | Already known TLS key will be added to newly created secret.                                                                                                                                                                                                                                                                                                                                  |
| `createSecret.cert`         | string  | no        | `-`            | Already known TLS certificate will be added to newly created secret.                                                                                                                                                                                                                                                                                                                          |
| `createSecret.secretName`   | string  | no        | `*-tls-secret` | Already known TLS certificates will be added to newly created secret with this name. Note that there should NOT be any secret with this name already.                                                                                                                                                                                                                                         |
<!-- markdownlint-enable line-length -->

**Note:**
`*-tls-secret` in the above section represents various names of secrets depending on the Victoriametrics components.
For example, when TLSConfig is configured for `vmagent`, the default name of the secret would be `vmagent-tls-secret`.

[Back to TOC](#table-of-content)

### version-exporter

#### install parameters
<!-- markdownlint-disable line-length -->
| Field                      | Description                                                                                                                                                                                                                                                  | Scheme                                                                                                                                     |
| -------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------ |
| install                    | Allows to enable or disable deploy version-exporter.                                                                                                                                                                                                         | bool                                                                                                                                       |
| name                       | A name of the microservice to deploy with. This name is used as the name of the microservice deployment and in labels.                                                                                                                                       | string                                                                                                                                     |
| containerPort              | Only changes container port. Application port can be changed with extraArgs (--web.listen-address=:9115).                                                                                                                                                    | integer                                                                                                                                    |
| image                      | A Docker image to deploy the version-exporter.                                                                                                                                                                                                               | string                                                                                                                                     |
| imagePullPolicy            | Image pull policy to use for version-exporter deployment                                                                                                                                                                                                     | string                                                                                                                                     |
| replicas                   | Number of created pods.                                                                                                                                                                                                                                      | *int32                                                                                                                                     |
| serviceAccount.create      | Allow to disable create ServiceAccount during deploy.                                                                                                                                                                                                        | bool                                                                                                                                       |
| serviceAccount.annotations | Annotations to add to the ServiceAccount during deploy.                                                                                                                                                                                                      | map[string]string                                                                                                                          |
| serviceAccount.name        | Provide a name of ServiceAccount to use.                                                                                                                                                                                                                     | bool                                                                                                                                       |
| resources                  | The resources that describe computed resource requests and limits for single pods.                                                                                                                                                                           | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#resourcerequirements-v1-core)               |
| securityContext            | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                                                           | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.18/#podsecuritycontext-v1-core)                  |
| extraArgs                  | Additional arguments for version-exporter container.                                                                                                                                                                                                         | [arguments](#arguments)                                                                                                                    |
| extraVarsSecret            | Allows set extra system environment variables for Version-exporter into the Secret.                                                                                                                                                                          | list[string]                                                                                                                               |
| serviceMonitor.enabled     | If true, a ServiceMonitor CRD is created for version-exporter.                                                                                                                                                                                               | boolean                                                                                                                                    |
| serviceMonitor.interval    | Default interval.                                                                                                                                                                                                                                            | string                                                                                                                                     |
| exporterConfig             | Config for exporters.                                                                                                                                                                                                                                        | [postgres collector parameters](#postgres-collector-parameters)  [http requests collector parameters](#http-requests-collector-parameters) |
| tolerations                | Tolerations allow the pods to schedule onto nodes with matching taints.                                                                                                                                                                                      | []v1.Toleration                                                                                                                            |
| nodeSelector               | Defines which nodes the pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\"                                                                                                                                           | map[string]string                                                                                                                          |
| affinity                   | It specifies the pod's scheduling constraints. For more information, refer to [https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core) | *v1.Affinity                                                                                                                               |
| annotations                | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value"                                       | map[string]string                                                                                                                          |
| labels                     | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                                                             | map[string]string                                                                                                                          |
| priorityClassName          | PriorityClassName assigned to the Pods to prevent them from evicting.                                                                                                                                                                                        | string                                                                                                                                     |
<!-- markdownlint-enable line-length -->

#### arguments

<!-- markdownlint-disable line-length -->
| Name                   | Description                                                  | Default value                             |
| ---------------------- | ------------------------------------------------------------ | ----------------------------------------- |
| `--log.level`          | Logging verbosity. Possible values: error, warn, info, debug | info                                      |
| `--web.listen-address` | Address on which to expose metrics and web interface.        | :9100                                     |
| `--web.telemetry-path` | Path under which to expose metrics.                          | /metrics                                  |
| `--web.max-requests`   | Maximum number of parallel scrape requests.                  | 40 (0 means no limit is applied)          |
| `--config.file`        | Path to a exporter configuration file.                       | /config/exporterConfig.yaml               |
| `--config.watch`       | Directory for watching change config map events.             | /config ("" means watching is turned off) |
<!-- markdownlint-enable line-length -->

#### postgres collector parameters

<!-- markdownlint-disable line-length -->
| Field                        | Description                                                                                                                                                                                     | Scheme       |
| ---------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------ |
| host                         | Postgres host name or ip to connect to. A part of postgres url. Must be unique.                                                                                                                 | string       |
| port                         | Postgres port number to connect to at the server host.                                                                                                                                          | string       |
| credentials.username.key     | Credentials for basic authentication. Secret key                                                                                                                                                | string       |
| credentials.username.name    | Credentials for basic authentication. Secret name                                                                                                                                               | string       |
| credentials.password.key     | Credentials for basic authentication. Secret key                                                                                                                                                | string       |
| credentials.password.name    | Credentials for basic authentication. Secret name                                                                                                                                               | string       |
| db                           | The database name.                                                                                                                                                                              | string       |
| timeout                      | Max connection life time is the duration since creation after which a connection will be automatically closed.                                                                                  | string       |
| requests.sql                 | Postgres sql request. Read-Only operations are exclusively supported. To avoid prometheus client errors about metric value duplicate use "select distinct" to selects only the distinct values. | string       |
| requests.metricName          | Name of new Prometheus metric.                                                                                                                                                                  | string       |
| requests.description         | Description of new Prometheus metric. Limit 100 symbols.                                                                                                                                        | string       |
| requests.metrics             | List of the field returned by sql request that is to be used in the metric. The "metrics" list fields will be included into the prometheus metric only.                                         | list[Metric] |
| requests.metrics.fieldName   | Name of the field returned by sql request that is to be used in the metric. It is to be unique.                                                                                                 | string       |
| requests.metrics.label       | Name of new Prometheus metric label. It is to be unique if defined.                                                                                                                             | string       |
| requests.metrics.valueRegexp | Regular expression applied to results of sql request                                                                                                                                            | string       |
<!-- markdownlint-enable line-length -->

If some requests.sql query is incorrect(e.g. request to non-existent table)
the appropriate error will be printed and no queries will be executed.

If requests.metrics.label is not defined, use as labels:

1. Named group of regexp
2. Column name of sql request results

Example:

```yaml
postgres_collector:
  connections:
    # Form postgres url, like
    # postgres://postgres:password@pg-patroni.postgres12.svc:5432/postgres?sslmode=disable
    - host: pg-patroni.postgres12.svc
      port: 5432
      credentials:
        username:
          key: pgUsername
          name: version-exporter-extra-vars-secret
        password:
          key: pgPassword
          name: version-exporter-extra-vars-secret
      db: postgres
      # Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
      timeout: 10s
      requests:
        - sql: select version()
          metricName: postgres_build_info
          metrics:
            # If metric exists, 'fieldName' must be present and either 'label' or 'valueRegexp' field must not be empty
            - fieldName: version
              label:
              valueRegexp: \s((?P<major>\d+).(?P<minor>\d+))\son\s(?P<platform>.*?),
        - sql: select extname, extversion, extnamespace from pg_extension
          metricName: postgres_build_info
          description: "Information about postgres version"
          metrics:
            - fieldName: extname
              label: extension_name
              valueRegexp:
            - fieldName: extversion
              label: extension_version
            - fieldName: extnamespace
        - sql: show server_version
          metricName: postgres_build_info
          metrics:
            - fieldName: server_version
              label: pg_server_version
              valueRegexp: ((\d+).(\d+))
    - host: 1.2.3.4
      port: 30063
      credentials:
        username:
          key: pgUsername
          name: version-exporter-extra-vars-secret
        password:
          key: pgPassword
          name: version-exporter-extra-vars-secret
      db: postgres
      # Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
      timeout: 10s
      requests:
        - sql: select version()
          metricName: postgres_build_info
          metrics:
            - fieldName: version
              valueRegexp: \s((?P<major>\d+).(?P<minor>\d+))\son\s(?P<platform>.*?),
        - sql: show server_version
          metricName: postgres_build_info
          metrics:
            - fieldName: server_version
              valueRegexp: ((?P<major>\d+).(?P<minor>\d+))
        - sql: select * from pg_extension
          metricName: postgres_build_info
          metrics:
            - fieldName: extname
              label: extension_name
              valueRegexp:
            - fieldName: extversion
              label: extension_version
            - fieldName: extnamespace
```

#### http requests collector parameters

HTTP requests collector can send `GET` and `POST` requests and parse responses with `Content-Type` of `application/json`
and `text/plain`.

To one endpoint (`url`) can be sent one or more requests. A Prometheus metric can be different for different requests.
If `metricName` is not filled, the results will be stored in metric with default name `http_collected_versions` with
configured `labels`. `labels` section is required and `label.name` must be unique.

If `Content-Type` of response is `application/json`, `label.jsonPath` value applied to this json and regular expression
`label.valueRegexp` applied to result of jsonpath search.
If `Content-Type` of response is `text/plain`, `label.valueRegexp` applied to result.

The parameters of `HTTP collector` should be provided under the section:

```yaml
version_exporter:
  #...
  exporterConfig:
    #...
    http_collector:
      connections:
      # list of connections
```

<!-- markdownlint-disable line-length -->

| Field                       | Description                                                                                                          | Scheme | Required |
| --------------------------- | -------------------------------------------------------------------------------------------------------------------- | ------ | -------- |
| url                         | A part of REST request (protocol, domain ant port). Must be unique.                                                  | string | true     |
| tlsConfig.tlsSkip           | Allow to disable certificates verification. Default: `false`.                                                        | bool   | false    |
| tlsConfig.ca                | Secret name and key where Certificate Authority is stored. Ignored if `tlsConfig.tlsSkip=true`.                      | object | true     |
| tlsConfig.cert              | Secret name and key where Certificate signing request is stored. Ignored if `tlsConfig.tlsSkip=true`.                | object | false    |
| tlsConfig.key               | Secret name and key where private key is stored. Ignored if `tlsConfig.tlsSkip=true`.                                | object | false    |
| credentials.username        | Credentials for basic authentication. Basic auth or token data should be provided if necessary.                      | object | false    |
| credentials.password        | Credentials for basic authentication.                                                                                | object | false    |
| credentials.token           | Credentials for token-based authentication.                                                                          | object | false    |
| requests.path               | Path and parameters of REST request (without protocol, domain ant port). Must be unique.                             | string | true     |
| requests.method             | Method of REST request. Possible values: `get/post`                                                                  | string | true     |
| requests.metricName         | Name of new Prometheus metric.                                                                                       | string | true     |
| requests.description        | Description of new Prometheus metric. Limit 100 symbols.                                                             | string | false    |
| requests.metrics.jsonPath   | JsonPath expressions applied to response of Content-Type "application/json". Must be unique.                         | string | false    |
| requests.labels.name        | Name of label of new Prometheus metric. Must be unique.                                                              | string | true     |
| requests.labels.valueRegexp | Regular expression applied to results of JsonPath search or to response of Content-Type "text/plain". Default: `.*`. | string | false    |

<!-- markdownlint-enable line-length -->

Example of configuration:

```yaml
http_collector:
  connections:
    - url: https://www.example.com
      tlsConfig:
        tlsSkip: false
        ca:
          key: cert-ca.pem
          name: secret-certificate-authority
      credentials:
        username:
          key: user
          name: version-exporter-extra-vars-secret
        password:
          key: password
          name: version-exporter-extra-vars-secret
      requests:
        - path: /version
          method: get
          metricName: "versions_ci_metric"
          metrics:
            - jsonPath: $.buildDate
              labels:
                - name: buildDate
                  valueRegexp: "[a-z0-9.]*"
            - jsonPath: $.gitCommit
              labels:
                - name: gitCommit
        - path: /product_versions
          method: get
          metricName: "versions_ci_product_metric"
          description: "Metric shows versions of product components"
          metrics:
            - jsonPath: $.go.goVersion
              labels:
                - name: goVersion
    - url: https://255.255.255.255:6300
      tlsConfig:
        tlsSkip: true
      credentials:
        token:
          key: token
          name: version-exporter-extra-vars-secret
      requests:
        - path: /version
          method: get
          metricName: "versions_ci_product_metric"
          metrics:
            - jsonPath: $.platform
              labels:
                - name: platform
                  valueRegexp: "[a-z0-9.]*"
        - path: /thirdparty_versions
          method: get
          metricName: "thirdparty_versions_ci_metric"
          metrics:
            - jsonPath: "{range .plugins[*]}{.name}{.version}{end}"
              labels:
                - name: component_name
                  valueRegexp: "[a-z0-9.]*"
                - name: version
```

The jsonpath library used here is `https://kubernetes.io/docs/reference/kubectl/jsonpath`

If you need to collect metrics with some empty values, use "{range ...} ... {end}" function.
Some examples of range-end jsonpath function:

```json
{range .plugins[*]}{.unique_id}{.name}{.version}{end}
```

or

```json
{range @[*]}{.name}{.version}{end}
```

Be careful, jsonpath with range function should not have plain text symbols, because version-exporter can not parse it.
Example of unsupported case:

```json
{range .items[*]}[{.metadata.name}, {.status.capacity}] {end}
```

#### ConfigMap collector parameters

ConfigMap collector is able to collect data from K8s resources such ConfigMaps and Secrets and
exposes it as prometheus metrics.

The parameters of `ConfigMap collector` should be provided under the section:

```yaml
version_exporter:
  #...
  exporterConfig:
    #...
    configmap_collector:
      defaults:
        # Default values to search k8s resources and take data from them
      resources:
        - name: # Unique name of searching resource (or group of resources)
          # Values that overrides the default section
        - #...
```

<!-- markdownlint-disable line-length -->
| Field                              | Description                                                                                                                                                                                                                                   | Scheme            |
| ---------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------- |
| defaults                           | Default values for configuration. Can be overrided in the resources section for each resource. Most of fields in this section are mandatory.                                                                                                  | object            |
| defaults.type                      | Type of k8s resource. Allowed values: configmap, secret.                                                                                                                                                                                      | string            |
| defaults.namespaces                | Namespaces in which to search for resources. Use empty list to search in all namespaces.                                                                                                                                                      | list[string]      |
| defaults.resourceLabels            | Allows to search resources by k8s labels instead of names. If resourceLabels is not empty, the search is carried out by labels, and not by the name of the resource.                                                                          | map[string]string |
| defaults.metricName                | Name of new Prometheus metric.                                                                                                                                                                                                                | string            |
| requests.description               | Description of new Prometheus metric. Limit 100 symbols.                                                                                                                                                                                      | string            |
| defaults.labels                    | Labels of new Prometheus metric.                                                                                                                                                                                                              | list[object]      |
| defaults.labels[N].name            | Name of label of new Prometheus metric. Value for this label will be taken with keyRegexp or valueRegexp. Must be unique.                                                                                                                     | string            |
| defaults.labels[N].keyRegexp       | Regular expression that will be used to find version information in keys of fields in Data from found ConfigMap or Secret. Each labels item must contain either the keyRegexp or the valueRegexp, but not both.                               | string            |
| defaults.labels[N].valueRegexp     | Regular expression that will be used to find version information in values of fields in Data from found ConfigMap or Secret. Each labels item must contain either the keyRegexp or the valueRegexp, but not both.                             | string            |
| resources                          | Each resource must contain the name parameter. Other fields are optional and override values from defaults section. If resourceLabels is empty, resource will be searched by name. Configuration must contain at least one item in this list. | list[object]      |
| resources[N].name                  | Name of resource. If resourceLabels is NOT empty, name parameter will not be used to find resources.                                                                                                                                          | string            |
| resources[N].type                  | Override value from defaults section. Allowed values: configmap, secret.                                                                                                                                                                      | string            |
| resources[N].namespaces            | Override value from defaults section.                                                                                                                                                                                                         | list[string]      |
| resources[N].resourceLabels        | Override value from defaults section.                                                                                                                                                                                                         | map[string]string |
| resources[N].metricName            | Override value from defaults section.                                                                                                                                                                                                         | string            |
| resources[N].description           | Override value from defaults section.                                                                                                                                                                                                         | string            |
| resources[N].labels                | Override value from defaults section.                                                                                                                                                                                                         | list[object]      |
| resources[N].labels[N].name        | Override value from defaults section.                                                                                                                                                                                                         | string            |
| resources[N].labels[N].keyRegexp   | Override value from defaults section.                                                                                                                                                                                                         | string            |
| resources[N].labels[N].valueRegexp | Override value from defaults section.                                                                                                                                                                                                         | string            |
<!-- markdownlint-enable line-length -->

**Note**: Almost all the default parameters are required fields.

You can use `name` or `labels` to find ConfigMaps or Secrets you are looking for.
If you declare the empty `resourceLabels` parameter, collector will search resources with the same name as in the
`resources[N].name` field in each of the specified namespaces.

If the `resourceLabels` is not empty, collector will search resources that contain labels from this parameter.
You should declare labels as `map[string]string`. Each label is AND connected.
In this case the `resources[N].name` will NOT be used for searching resources (if `resourceLabels` is not empty,
you can use fake `resources[N].name`), but every item of the `resources`
list must contain this field, and it should be unique anyway.

Parameter `namespaces` allows to specify namespaces for searching resources. Empty list means that resources will be
searched in all namespaces.

Collector parses and processes data from ConfigMaps and Secrets with regular expressions from the `labels` sections.
Configuration must contain at least one item in the `labels` list. Each item must contain unique `name` that will be
used as label name for Prometheus metrics, and must contain either `keyRegexp` or `valueRegexp`.

The data in ConfigMaps and Secrets is a `map[string]string` where each field has a key and a value.
When you specify parameter `keyRegexp` for label item, it means that collector will parse _keys_ from the data with
regexp from the parameter. If you specify `valueRegexp`, collector will parse _values_ with this regexp.

For example, collector found ConfigMap with follow content:

```yaml
...
data:
  example-app: v0.0.1
  ...
```

Configuration of the collector contains follow statements:

```yaml
...
  metricName: configmap_collected_versions
  labels:
    - name: app_name
      keyRegexp: ".*"
    - name: app_version
      valueRegexp: ".*"
```

Then exposed metrics will contain the follow one:

```bash
configmap_collected_versions{...app_name="example-app",app_version="v0.0.1"}
```

Example of configuration:

```yaml
configmap_collector:
  defaults:
    type: configmap
    namespaces: ["monitoring"]
    resourceLabels: {}
    metricName: "configmap_collected_versions"
    description: "Metric shows version collected from configmaps"
    labels:
      - name: key_label
        keyRegexp: "[a-z0-9.]*"
      - name: value_label
        valueRegexp: ".*"
  resources:
    # Collector will search configmaps with label "label-key: label-value" in all namespaces
    - name: version-configmap-test
      namespaces: []
      resourceLabels:
        label-key: label-value
      metricName: "configmap_versions"
      description: "Help for Prometheus metric"
      labels:
        - name: version_from_key
          keyRegexp: "version-app-*"
    # Collector will search secret with name "version-secret-test" in the "monitoring" namespace
    - name: version-secret-test
      type: secret
    # Collector will search configmap with name "version-default-configmap" in the "monitoring" namespace
    - name: version-default-configmap
```

### SSH collector parameters

SSH collector is able to collect data over ssh by requests provided in configuration and
exposes it as prometheus metrics.

The parameters of `SSH collector` should be provided under the section:

```yaml
version_exporter:
  #...
  exporterConfig:
    #...
    ssh_collector:
      connections:
      # list of connections
```

<!-- markdownlint-disable line-length -->
| Field                       | Description                                                                                                    | Scheme |
| --------------------------- | -------------------------------------------------------------------------------------------------------------- | ------ |
| host                        | SSH server host name or ip to connect to. Must be unique.                                                      | string |
| port                        | Port number to connect to the server host.                                                                     | string |
| network                     | Network type to connect to the server host. ("tcp" is only supported)                                          | string |
| k8sCredentials.login.key    | Credentials for basic authentication. Secret key                                                               | string |
| k8sCredentials.login.name   | Credentials for basic authentication. Secret name                                                              | string |
| k8sCredentials.pkey.key     | Credentials for basic authentication. Secret key                                                               | string |
| k8sCredentials.pkey.name    | Credentials for basic authentication. Secret name                                                              | string |
| timeout                     | Max connection life time is the duration since creation after which a connection will be automatically closed. | string |
| requests.cmd                | cmd command request. Supported ones: cat, nl, head, tail, echo, hostname, uname                                | string |
| requests.metricName         | Name of new Prometheus metric.                                                                                 | string |
| requests.description        | Description of new Prometheus metric. Limit 100 symbols.                                                       | string |
| requests.labels.name        | Name of new Prometheus metric label. It is to be unique if defined.                                            | string |
| requests.labels.valueRegexp | Regular expression applied to results of cmd request                                                           | string |
<!-- markdownlint-enable line-length -->

Example of configuration:

```yaml
ssh_collector:
  connections:
    - host: 1.2.3.4
      port: 22
      network: tcp
      # Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
      timeout: 5s
      k8sCredentials:
        login:
          key: sshLogin
          name: version-exporter-extra-vars-secret
        pkey:
          key: privateKey
          name: version-exporter-extra-vars-secret
      requests:
        - cmd: "head /etc/os-release"
          metricName: "os_versions_metric"
          labels:
            - name: os_version
              valueRegexp: 'VERSION="(?P<version>[a-zA-Z0-9. ()]*)"'
            - name: name
              valueRegexp: 'NAME="(?P<name>[a-zA-Z ]*)"'
        - cmd: 'tail /etc/ssh/ssh_config'
          metricName: "ssh_versions_metric"
          description: "Metric shows versions of ssh component"
          labels:
            - name: ssh_version
              valueRegexp: "[^a-zA-Z0-9]"
        - cmd: 'cat /etc/ssh/ssh_config'
          metricName: "ssh_ver_metric"
          description: "Metric shows versions of ssh component"
          labels:
            - name: version
              valueRegexp: "[^a-zA-Z0-9]"
    - host: 1.2.3.4
      port: 22
      network: tcp
      # Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
      timeout: 5s
      k8sCredentials:
        login:
          key: sshLogin
          name: version-exporter-extra-vars-secret
        pkey:
          key: privateKey
          name: version-exporter-extra-vars-secret
      requests:
        - cmd: 'echo `ssh -V`'
          metricName: "ssh_version_metric"
          labels:
            - name: ssh_version
              valueRegexp: ".*[^\n]"
        - cmd: 'hostname --fqdn'
          metricName: "hostname_metric"
          labels:
            - name: version
              valueRegexp: "[^a-zA-Z0-9]"
```

### Integration tests

<!-- markdownlint-disable line-length -->
| Field                              | Description                                                                                                                                                                                                                                                  | Scheme            |
| ---------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ----------------- |
| install                            | The parameter specifies whether integration tests are installed or not. Set to `true` to run tests.                                                                                                                                                          | bool              |
| image                              | The parameter specifies the Docker image of integration tests.                                                                                                                                                                                               | string            |
| resources.requests.memory          | The parameter specifies the minimum amount of memory the container should use.                                                                                                                                                                               | string            |
| resources.requests.cpu             | The parameter specifies the minimum number of CPUs the container should use.                                                                                                                                                                                 | string            |
| resources.limits.memory            | The parameter specifies the maximum amount of memory the container can use.                                                                                                                                                                                  | string            |
| resources.limits.cpu               | The parameter specifies the maximum number of CPUs the container can use.                                                                                                                                                                                    | string            |
| securityContext                    | The monitoring tests container security context.                                                                                                                                                                                                             | string            |
| tolerations                        | Tolerations allow the pods to schedule onto nodes with matching taints.                                                                                                                                                                                      | []v1.Toleration   |
| nodeSelector                       | Defines which nodes the pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\"                                                                                                                                           | map[string]string |
| affinity                           | It specifies the pod's scheduling constraints. For more information, refer to [https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core) | *v1.Affinity      |
| annotations                        | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value"                                       | map[string]string |
| labels                             | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                                                             | map[string]string |
| timeoutBeforeStart                 | The parameter specifies the timeout before the start of integration tests. Default value is `120`                                                                                                                                                            | string            |
| tags                               | The parameter specifies the tags combined together with AND, OR, and NOT operators that select the test cases to run. Default value is `grafanaORtest-appORprometheus-adapterORsmoke`. To run all tests, set `full`.                                         | string            |
| priorityClassName                  | PriorityClassName assigned to the Pods to prevent them from evicting.                                                                                                                                                                                        | string            |
| statusWriting.enabled              | The parameter specifies enable/disable writing status of Integration tests' execution to the specified Custom Resource.                                                                                                                                      | bool              |
| statusWriting.isShortStatusMessage | The parameter specifies the size of the integration test status message.                                                                                                                                                                                     | bool              |
| statusWriting.onlyIntegrationTests | The parameter specifies whether to deploy only integration tests without any component (component was installed before).                                                                                                                                     | bool              |
| statusWriting.customResourcePath   | The parameter (optional) specifies the path to Custom Resource that should be used for the write status of auto-tests execution.                                                                                                                             | string            |
<!-- markdownlint-enable line-length -->

Minimal configuration to run tests

```yaml
integrationTests:
  install: true
```

Examples of configuration:

```yaml
integrationTests:
  install: true
  tags: full
  statusWriting:
    enabled: true
    isShortStatusMessage: true
    onlyIntegrationTests: true
    customResourcePath: "monitoring.qubership.org/v1alpha1/monitoring/platformmonitorings/platformmonitoring"
```

```yaml
integrationTests:
  name: monitoring-tests
  install: false
  resources:
    limits:
      cpu: 300m
      memory: 400Mi
    requests:
      cpu: 100m
      memory: 200Mi
  tags: grafanaORtest-appORprometheus-adapterORsmoke
  timeoutBeforeStart: 300
  statusWriting:
    enabled: true
    isShortStatusMessage: true
    onlyIntegrationTests: true
    customResourcePath: "monitoring.qubership.org/v1alpha1/monitoring/platformmonitorings/platformmonitoring"
```

[Back to TOC](#table-of-content)

# Deploy

## Admission Migration

Since Kubernetes v1.20 Pod Security Policy (PSP) has been deprecated and will be deleted in Kubernetes 1.25
the migration procedure from PSP to another solution is very important.

The Pod Security Standard privileged policy is required for the monitoring system.
It can be defined on namespace level:

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: my-privileged-namespace
  labels:
    pod-security.kubernetes.io/enforce: privileged
    pod-security.kubernetes.io/enforce-version: latest
```

After migrating but before upgrade monitoring system, please delete all secrets with name like "sh.helm*"

## CRD Upgrade

The upgrade of CRD happens automatically via pre-deploy scripts.

**NOTE:** Automation CRD upgrade requires the corresponding cluster rights for deploy user:

```yaml
- apiGroups: ["apiextensions.k8s.io"]
  resources: ["customresourcedefinitions"]
  verbs: ["get", "create", "patch"]
```

To disable this feature you need to add parameter `DISABLE_CRD: true`.

## Cleanup VictoriaMetrics Objects

Objects were created by victoria-metrics-operator have finalizers. They must be removed manually before clean
installation on cloud with victoria-metrics stack.

You can use this commands. Don't forget change namespace `-n <namespace>`:

```shell
kubectl patch deployment victoriametrics-operator --type json -p '[{"op": "replace", "path": "/metadata/finalizers", "value": null}]' -n <namespace>
kubectl patch deployment vmsingle-k8s --type json -p '[{"op": "replace", "path": "/metadata/finalizers", "value": null}]' -n <namespace>
kubectl patch deployment vmalert-k8s --type json -p '[{"op": "replace", "path": "/metadata/finalizers", "value": null}]' -n <namespace>
kubectl patch deployment vmagent-k8s --type json -p '[{"op": "replace", "path": "/metadata/finalizers", "value": null}]' -n <namespace>
kubectl patch deployment vmauth-k8s --type json -p '[{"op": "replace", "path": "/metadata/finalizers", "value": null}]' -n <namespace>
kubectl patch deployment vminsert-k8s --type json -p '[{"op": "replace", "path": "/metadata/finalizers", "value": null}]' -n monitoring
kubectl patch statefulset vmalertmanager-k8s --type json -p '[{"op": "replace", "path": "/metadata/finalizers", "value": null}]' -n <namespace>
kubectl patch statefulset vmselect-k8s --type json -p '[{"op": "replace", "path": "/metadata/finalizers", "value": null}]' -n monitoring
kubectl patch statefulset vmstorage-k8s --type json -p '[{"op": "replace", "path": "/metadata/finalizers", "value": null}]' -n monitoring

kubectl patch serviceaccount vmalertmanager-k8s --type json -p '[{"op": "replace", "path": "/metadata/finalizers", "value": null}]' -n <namespace>
kubectl patch serviceaccount vmalert-k8s --type json -p '[{"op": "replace", "path": "/metadata/finalizers", "value": null}]' -n <namespace>
kubectl patch serviceaccount vmagent-k8s --type json -p '[{"op": "replace", "path": "/metadata/finalizers", "value": null}]' -n <namespace>
kubectl patch serviceaccount vmsingle-k8s --type json -p '[{"op": "replace", "path": "/metadata/finalizers", "value": null}]' -n <namespace>
oc patch secret vmalertmanager-config-secret --type json -p '[{"op": "replace", "path": "/metadata/finalizers", "value": null}]' -n <namespace>
kubectl patch serviceaccount vmauth-k8s --type json -p '[{"op": "replace", "path": "/metadata/finalizers", "value": null}]' -n <namespace>
oc patch secret vmalertmanager-config-secret --type json -p '[{"op": "replace", "path": "/metadata/finalizers", "value": null}]' -n <namespace>

kubectl patch secret vmalertmanager-k8s-config --type json -p '[{"op": "replace", "path": "/metadata/finalizers", "value": null}]' -n monitoring
kubectl patch secret tls-assets-vmalert-k8s --type json -p '[{"op": "replace", "path": "/metadata/finalizers", "value": null}]' -n <namespace>
kubectl patch secret tls-assets-vmagent-k8s --type json -p '[{"op": "replace", "path": "/metadata/finalizers", "value": null}]' -n <namespace>
kubectl patch secret vmagent-k8s --type json -p '[{"op": "replace", "path": "/metadata/finalizers", "value": null}]' -n <namespace>
kubectl patch configmap vm-vmalert-rulefiles-0 --type json -p '[{"op": "replace", "path": "/metadata/finalizers", "value": null}]' -n <namespace>
kubectl patch configmap vm-k8s-rulefiles-0 --type json -p '[{"op": "replace", "path": "/metadata/finalizers", "value": null}]' -n <namespace>
kubectl patch configmap relabelings-assets-vmagent-k8s --type json -p '[{"op": "replace", "path": "/metadata/finalizers", "value": null}]' -n <namespace>
kubectl patch pvc vmsingle-vmsingle --type json -p '[{"op": "replace", "path": "/metadata/finalizers", "value": null}]' -n <namespace>

kubectl patch service vmalertmanager-k8s --type json -p '[{"op": "replace", "path": "/metadata/finalizers", "value": null}]' -n <namespace>
kubectl patch service vmalert-k8s --type json -p '[{"op": "replace", "path": "/metadata/finalizers", "value": null}]' -n <namespace>
kubectl patch service vmauth-k8s --type json -p '[{"op": "replace", "path": "/metadata/finalizers", "value": null}]' -n <namespace>
kubectl patch service vmsingle-k8s --type json -p '[{"op": "replace", "path": "/metadata/finalizers", "value": null}]' -n <namespace>
kubectl patch service vmagent-k8s --type json -p '[{"op": "replace", "path": "/metadata/finalizers", "value": null}]' -n <namespace>
kubectl patch service vmselect-k8s --type json -p '[{"op": "replace", "path": "/metadata/finalizers", "value": null}]' -n monitoring
kubectl patch service vminsert-k8s --type json -p '[{"op": "replace", "path": "/metadata/finalizers", "value": null}]' -n monitoring
kubectl patch service vmstorage-k8s --type json -p '[{"op": "replace", "path": "/metadata/finalizers", "value": null}]' -n monitoring

kubectl patch vmagents k8s --type json -p '[{"op": "replace", "path": "/metadata/finalizers", "value": null}]' -n <namespace>
kubectl patch vmalertmanagers  k8s --type json -p '[{"op": "replace", "path": "/metadata/finalizers", "value": null}]' -n <namespace>
kubectl patch vmalerts k8s --type json -p '[{"op": "replace", "path": "/metadata/finalizers", "value": null}]' -n <namespace>
kubectl patch vmauths k8s --type json -p '[{"op": "replace", "path": "/metadata/finalizers", "value": null}]' -n <namespace>
kubectl patch vmsingles k8s --type json -p '[{"op": "replace", "path": "/metadata/finalizers", "value": null}]' -n <namespace>
kubectl patch vmclusters k8s --type json -p '[{"op": "replace", "path": "/metadata/finalizers", "value": null}]' -n monitoring

kubectl patch psp vmalertmanager-k8s --type json -p '[{"op": "replace", "path": "/metadata/finalizers", "value": null}]'
kubectl patch psp vmalert-k8s --type json -p '[{"op": "replace", "path": "/metadata/finalizers", "value": null}]'
kubectl patch psp vmagent-k8s --type json -p '[{"op": "replace", "path": "/metadata/finalizers", "value": null}]'
kubectl patch psp vmsingle-k8s --type json -p '[{"op": "replace", "path": "/metadata/finalizers", "value": null}]'
kubectl patch psp vmauths-k8s --type json -p '[{"op": "replace", "path": "/metadata/finalizers", "value": null}]'
kubectl patch psp vmclusters-k8s --type json -p '[{"op": "replace", "path": "/metadata/finalizers", "value": null}]'
```

## Upgrade

This section contains information about some actions that you may do during upgrades.

### Upgrade from 0.50.0 version

If you are trying to upgrade `monitoring-operator` with `prometheus-adapter-operator`
from older version to `0.50.0` or higher, it will be necessary to perform some manual steps.

Older version of `prometheus-adapter-operator` installed `APIService` object from the operator.
Since `0.50.0`, `APIService` object will be created by `Helm`.
During rolling update, `Helm` cannot replace already created `APIService` with new one.
It is necessary to manually delete the `APIService` object created by the operator using below command.

```bash
kubectl delete apiservice v1beta1.custom.metrics.k8s.io
```

After this, rolling update can be performed and `Helm` is going to create new `APIService` object.

## Manual Deploy using Helm

**Note**: If you want to deploy the `monitoring-operator` into **Kubernetes v1.15** or lower or
**OpenShift v3.11** or lower, you must work with _v1beta1_ CRDs manually. For more information see
[Maintenance Guide: Work with legacy CRDs](maintenance.md#work-with-legacy-crds).

This chart installs Monitoring Operator which can create/configure/manage Prometheus/Victoriametrics and the components
in Kubernetes/OpenShift.

### Ingress using Helm

Ingress is enabled by default.
If you want `host` to be installed automatically, you have to specify next parameters:

```yaml
CLOUD_PUBLIC_URL: <public_url.com>
NAMESPACE: <monitoring>
```

Ingress `host` will be set such as `<component>-{{ .Values.NAMESPACE }}.{{ .Values.CLOUD_PUBLIC_URL }}`.
For example, grafana-monitoring.public_url.com, victoriametrics-monitoring.public_url.com, etc.

Also, you can specify ingress structure for each component. For example:

```yaml
grafana:
  ingress:
    host: <grafana.example.com>

alertmanager:
  ingress:
    install: false
```

### Installing the Chart

To install the chart with the release name `monitoring-operator`:

```bash
helm install monitoring-operator charts/monitoring-operator
```

The command deploys monitoring-operator on the Kubernetes/OpenShift cluster in the default configuration. The
configuration section lists the parameters that can be configured during the installation.

The default installation includes Victoriametrics Operator, AlertManager, Exporters, and configuration for scraping the
Kubernetes/OpenShift infrastructure.

[Back to TOC](#table-of-content)

### Upgrading the Chart

To upgrade the chart with the release name `monitoring-operator`:

```bash
helm upgrade monitoring-operator charts/monitoring-operator
```

The command upgrades monitoring-operator on the Kubernetes/OpenShift cluster in the default configuration. The
configuration section lists the parameters that can be configured during the upgrade.

[Back to TOC](#table-of-content)

### Uninstalling the Chart

To uninstall the `monitoring-operator` deployment:

```bash
helm uninstall monitoring-operator
```

The command removes all the Kubernetes/OpenShift components associated with the chart and deletes the release.

**Warning**: Note that this step removes CRDs for monitoring and deleting these CRDs causes the deletion of all
resources of their type. It means that all resources like ServiceMonitor and GrafanaDashboard from applications are
removed.

CRDs created by this chart are not removed by default and should be manually cleaned up:

* Kubernetes

    ```bash
    kubectl delete crd grafanas.integreatly.org
    kubectl delete crd grafanadashboards.integreatly.org
    kubectl delete crd grafanadatasources.integreatly.org
    kubectl delete crd grafananotificationchannels.integreatly.org
    kubectl delete crd alertmanagers.monitoring.coreos.com
    kubectl delete crd alertmanagerconfigs.monitoring.coreos.com
    kubectl delete crd podmonitors.monitoring.coreos.com
    kubectl delete crd probes.monitoring.coreos.com
    kubectl delete crd servicemonitors.monitoring.coreos.com
    kubectl delete crd thanosrulers.monitoring.coreos.com
    kubectl delete crd customscalemetricrules.monitoring.qubership.org
    kubectl delete crd platformmonitorings.monitoring.qubership.org
    kubectl delete crd vmsingles.operator.victoriametrics.com
    kubectl delete crd vmagents.operator.victoriametrics.com
    kubectl delete crd vmalertmanagers.operator.victoriametrics.com
    kubectl delete crd vmalerts.operator.victoriametrics.com
    ```

* For OpenShift

    ```bash
    oc delete crd grafanas.integreatly.org
    oc delete crd grafanadashboards.integreatly.org
    oc delete crd grafanadatasources.integreatly.org
    oc delete crd grafananotificationchannels.integreatly.org
    oc delete crd alertmanagers.monitoring.coreos.com
    oc delete crd alertmanagerconfigs.monitoring.coreos.com
    oc delete crd podmonitors.monitoring.coreos.com
    oc delete crd probes.monitoring.coreos.com
    oc delete crd servicemonitors.monitoring.coreos.com
    oc delete crd thanosrulers.monitoring.coreos.com
    oc delete crd customscalemetricrules.monitoring.qubership.org
    oc delete crd platformmonitorings.monitoring.qubership.org
    oc delete crd vmsingles.operator.victoriametrics.com
    oc delete crd vmagents.operator.victoriametrics.com
    oc delete crd vmalertmanagers.operator.victoriametrics.com
    oc delete crd vmalerts.operator.victoriametrics.com
    ```

[Back to TOC](#table-of-content)

# Post Deploy Checks

There are many ways to check whether Monitoring is working correctly and collecting metrics after the deployment.
Depending on the selected deploy tool, some common options and some
specific for deploy tool options are available.

## Common Post Deploy Checks

Firstly, you can check that all components are up using the following command:

```bash
kubectl get pods -n <monitoring_namespace>
```

For example, a typical Monitoring deployment contains the following pods:

```bash
$ kubectl get pods -n monitoring
NAME                                          READY   STATUS    RESTARTS   AGE
alertmanager-k8s-0                            2/2     Running   0          9d
grafana-deployment-7d869f989d-24657           1/1     Running   0          9d
grafana-operator-5857f4d47d-f7nv4             1/1     Running   0          9d
kube-state-metrics-6c89f47c94-tbx7l           1/1     Running   0          9d
monitoring-operator-monito-6d6c78dd86-wqfcx   1/1     Running   0          9d
node-exporter-g5hsr                           1/1     Running   0          9d
node-exporter-jq7gq                           1/1     Running   0          9d
node-exporter-zppgl                           1/1     Running   0          9d
victoriametrics-operator-c5649d646-7tnw8      1/1     Running   0          9d
vmagent-k8s-8dd4ffccc-msdmx                   2/2     Running   0          9d
vmalert-k8s-6c4b989889-dndmv                  2/2     Running   0          9d
vmalertmanager-k8s-0                          2/2     Running   0          9d
vmauth-k8s-76cc654596-gn8pl                   1/1     Running   0          9d
vmsingle-k8s-6cf5b84846-c22xs                 0/1     Pending   0          9d
```

Or you can use the Kubernetes Dashboard to see pods and their statuses in the UI.

Secondly, you can check that all Victoriametrics services have the `OK` status at the `/health` probe,
using the following command:

```bash
curl "http://<victoriametrics pod ip or service name>/health
```

or

```bash
wget -O - "http://<victoriametrics pod ip or service name>/health"
```

Thirdly, you can try to select any metrics that are collected by Victoriametrics for a typical deployment.
You can get the list of metrics by the next command:

```bash
curl "http://<victoriametrics pod ip or service name>/api/v1/labels
```

To select the metrics, you can use the next command:

```bash
curl "http://<victoriametrics pod ip or service name>/api/v1/label/<label_name>/values
```

with the metric name in the field `label_name`. For each metric, a non-empty result should be returned.

[Back to TOC](#table-of-content)

# What's next

* Read more about [available metrics which collected after deploy](metrics-oob.md)
* Read more about [available alerts after deploy](alerts-oob.md)
* Read more about [available dashboards after deploy](dashboards-oob)
* Read more about [maintenance](maintenance.md)
* Read more about [troubleshooting](troubleshooting.md)
