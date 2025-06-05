This document provides information about the maintenance of the Monitoring deployment and the configuration of its
services.

# Ways to Make Changes

The following sections describe the various ways to make changes in the Monitoring deployment.

## Redeploy/Update

### Manual Deploy using Helm

**NOTE:** if you want to deploy the `monitoring-operator` into **Kubernetes v1.15** or lower or
**OpenShift v3.11** or lower, you must work with _v1beta1_ CRDs manually. For more information see
[Work with legacy CRDs](#work-with-legacy-crds).

This chart installs Monitoring Operator which can create/configure/manage Prometheus and the components in
Kubernetes/OpenShift.

#### Installing the Chart

To install the chart with the `monitoring-operator` release name, use the following command.

```bash
helm install monitoring-operator charts/monitoring-operator
```

The command deploys monitoring-operator on the Kubernetes/OpenShift cluster in the default configuration. The
configuration section lists the parameters that can be configured during the installation.

The default installation includes Prometheus Operator, AlertManager, Exporters, and configuration for scraping the
Kubernetes/OpenShift infrastructure.

#### Upgrading the Chart

To upgrade the chart with the `monitoring-operator` release name, use the following command.

```bash
helm upgrade monitoring-operator charts/monitoring-operator
```

The command upgrades monitoring-operator on the Kubernetes/OpenShift cluster in the default configuration. The
configuration section lists the parameters that can be configured during the upgrade.

#### Uninstalling the Chart

To uninstall or delete the `monitoring-operator` deployment, use the following command.

```bash
helm delete monitoring-operator
```

The command removes all the Kubernetes/OpenShift components associated with the chart and deletes the release.

**Warning**: Note that this step removes CRDs for monitoring and deleting these CRDs deletes all resources of their
type. All resources like ServiceMonitor and GrafanaDashboard are removed from the applications.

CRDs created by this chart are not removed by default and should be manually cleaned up.

* For Kubernetes

    ```bash
    kubectl delete crd platformmonitoring.monitoring.qubership.org
    kubectl delete crd prometheuses.monitoring.coreos.com
    kubectl delete crd prometheusrules.monitoring.coreos.com
    kubectl delete crd servicemonitors.monitoring.coreos.com
    kubectl delete crd podmonitors.monitoring.coreos.com
    kubectl delete crd alertmanagers.monitoring.coreos.com
    kubectl delete crd grafana.integreatly.org
    kubectl delete crd grafanadashboard.integreatly.org
    kubectl delete crd grafanadatasource.integreatly.org
    ```

* For OpenShift

    ```bash
    oc delete crd platformmonitoring.monitoring.qubership.org
    oc delete crd prometheuses.monitoring.coreos.com
    oc delete crd prometheusrules.monitoring.coreos.com
    oc delete crd servicemonitors.monitoring.coreos.com
    oc delete crd podmonitors.monitoring.coreos.com
    oc delete crd alertmanagers.monitoring.coreos.com
    oc delete crd grafana.integreatly.org
    oc delete crd grafanadashboard.integreatly.org
    oc delete crd grafanadatasource.integreatly.org
    ```

## Change Parameters in Custom Resource in Runtime

For deploy and config, the Monitoring deployment uses monitoring-operator, which is based on operator-sdk.

It means that monitoring-operator controls most part of the Monitoring deployment. The settings to create/update/remove
any part of the deployment is used from the Custom Resource (CR) with type `PlatformMonitoring`.

So you can make changes in the Monitoring deployment just by changing the parameters in the `PlatformMonitoring` CR.

**Note**: To change the settings in the `PlatformMonitoring` CR, you must have permissions on the `PlatformMonitoring`
CR:

```yaml
- apiGroups:
    - "monitoring.qubership.org"
  resources:
    - platformmonitorings
  verbs:
    - 'create'
    - 'delete'
    - 'get'
    - 'list'
    - 'update'
```

Monitoring-operator watches resources with type `PlatformMonitoring` and can handle the following events:

* `create` - If the operator receives such an event, it runs the monitoring deploy in the namespace where the CR was
  created.
* `update` - If the operator receives such an event, it runs the update of monitoring deployment in the namespace where
  the CR was updated.
* `delete` - If the operator receives such an event, it runs remove of all objects (Deployments, StatefulSets,
  ConfigMaps, and so on), which it earlier created in the namespace where the CR was removed.

Thus using the operator, we can change the monitoring deployment parameters in runtime. To change the parameters:

1. Log in to Kubernetes through the `kubectl` cli or using any UI.
2. Find an object with type `PlatformMonitoring` and name `platformmonitoring` in the namespace where Monitoring is
   deployed. For example, use the following cli command.

   ```bash
   kubectl get platformmonitorings.monitoring.qubership.org platformmonitoring -n <namespace>
   ```

3. Use the following cli command to open the object for editing.

   ```bash
   kubectl edit platformmonitorings.monitoring.qubership.org platformmonitoring -n <namespace>
   ```

4. Change the necessary parameters and save the object.
5. Wait while the operator applies changes in runtime.

Using this procedure, you can:

* Change the component's images.
* Change any component specific settings, like retention period for Prometheus.
* Add or change persistence volumes for Prometheus and Grafana.
* Disable some components or stop tracking it by operator, and so on.

All parameters that you can specify in the `PlatformMonitoring` objects are described in the API documents. For more
details, refer to the _Platform Monitoring_ chapter in the _Cloud Platform Monitoring Guide_.

## Work with legacy CRDs

If you use **Kubernetes v1.15** or lower or **OpenShift v3.11** or lower, you must use CRDs with _v1beta1_ API version.
But the Helm chart uses CRDs _v1_ version, which are incompatible with _v1beta1_ CRDs.

So if you use Kubernetes or OpenShift of specified version or lower, you have to work with _v1beta1_ CRDs which
are contained in the [crds directory](crds/). Also, you have to use the features of deployment tools
to skip installing CRDs to avoid errors.

If you want to deploy the `monitoring-operator` in a cluster (Kubernetes v1.15 or lower or
OpenShift v3.11 or lower) for the first time, you have to install CRDs manually before deploying other components.

To create CRDs version _v1beta1_ need:

1. Login into cluster using `kubectl` client for Kubernetes or `oc` client for `OpenShift`
2. Navigate to `docs/crds` directory
3. Execute the following command for Kubernetes:

```bash
kubectl create -f ./v1beta1
```

or for OpenShift:

```bash
oc create -f ./v1beta1
```

To update CRDs version _v1beta1_ need:

1. Login into cluster using `kubectl` client for Kubernetes or `oc` client for `OpenShift`
2. Navigate to `docs/crds` directory
3. Execute the following command for Kubernetes:

```bash
kubectl replace -f ./v1beta1
```

or for OpenShift:

```bash
oc replace -f ./v1beta1
```

To remove CRDs you can use instructions from [Remove procedure](#remove-procedure).

**NOTE:** You must set the following parameters during deploy to skip CRDs creation if you use _v1beta1_ CRDs:

* for manually Helm installation - add `--skip-crds` to `helm install` or `helm upgrade` command

# Provided Procedures

The information about the procedures for monitoring deployment is described in the following sections.

## Update Procedure

Except for one important point, the update process is the same as the installation process.

**Important:** Helm doesn't have an ability to update/remove CRDs, but can create them. It means that CRDs will
create during first deploy. But during any updates and remove deployment Helm will not update or remove CRDs.
For more details you can read [Issue: CRD update during helm upgrade](https://github.com/helm/helm/issues/8668) and
[Design: CRD Handling in Helm](https://github.com/helm/community/blob/f9e06c16d89ccea1bea77c01a6a96ae3b309f823/architecture/crds.md).

The update procedure almost completely repeats the installation, with the exception of one important point.

During install/update you need:

* Create or use previously created configurations.
* Run any job for update or run update manually with using Helm.

But before run update need manually update CRDs.

**NOTE:** if you use the `monitoring-operator` into **Kubernetes v1.15** or lower or
**OpenShift v3.11** or lower, you must work with _v1beta1_ CRDs manually. For more information see
[Work with legacy CRDs](#work-with-legacy-crds).

To update CRDs need:

1. Download chart. How to find it see below.
2. Navigate to directory `charts/monitoring-operator`.
3. Execute a command:

    ```bash
    kubectl replace -f crds/
    ```

## Remove Procedure

There are some ways to uninstall or delete the `monitoring-operator` deployment:

* Remove only the `PlatformMonitoring` CR.
* Remove the `PlatformMonitoring` CR and operator.
* Remove all created CRDs.

### Remove PlatformMonitoring CR

Monitoring-operator watches the state of the `PlatformMonitoring` CR and tracks all events. So when you remove this, the
object operator should remove all objects (Deployments, ConfigMaps, Secrets, and so on), which it created during the
deployment.

This method allows you to remove Monitoring from namespace, but keeps the operator that redeploys all necessary
Monitoring components when you decide to create the `PlatformMonitoring` CR again.

To remove the `PlatformMonitoring` CR, you can use the following command.

* For Kubernetes

    ```bash
    kubectl delete platformmonitorings.monitoring.qubership.org platformmonitoring -n <namespace>
    ```

* For OpenShift

    ```bash
    oc delete platformmonitorings.monitoring.qubership.org platformmonitoring -n <namespace>
    ```

After executing the command, all Monitoring components are removed from the selected namespace. However, the following
objects are not removed:

* `monitoring-operator`
* Some cluster entities like ClusterRoles, ClusterRoleBindings, Security Context Constraints (SCC), and Pod Security
  Policy (PSP)

### Remove PlatformMonitoring CR and Operator

The simplest way to remove the CR and operator is by using the following helm delete command.

```bash
helm delete <release-name> -n <namespace>
```

To find the release name, use the following command.

```bash
helm list -n <namespace>
```

The command removes all the Kubernetes/OpenShift components associated with the chart and deletes the release.

### Remove all Created CRDs

**Warning**: Note that this step removes CRDs for monitoring, and deleting these CRDs causes the deletion of all
resources of their type. It means that all resources like ServiceMonitor and GrafanaDashboard are removed from the
applications.

CRDs created by this chart are not removed by default and should be manually cleaned up.

* For Kubernetes

    ```bash
    kubectl delete crd grafanas.integreatly.org
    kubectl delete crd grafanadashboards.integreatly.org
    kubectl delete crd grafanadatasources.integreatly.org
    kubectl delete crd grafananotificationchannels.integreatly.org
    kubectl delete crd alertmanagers.monitoring.coreos.com
    kubectl delete crd alertmanagerconfigs.monitoring.coreos.com
    kubectl delete crd podmonitors.monitoring.coreos.com
    kubectl delete crd probes.monitoring.coreos.com
    kubectl delete crd prometheuses.monitoring.coreos.com
    kubectl delete crd prometheusrules.monitoring.coreos.com
    kubectl delete crd servicemonitors.monitoring.coreos.com
    kubectl delete crd thanosrulers.monitoring.coreos.com
    kubectl delete crd customscalemetricrules.monitoring.qubership.org
    kubectl delete crd platformmonitorings.monitoring.qubership.org
    kubectl delete crd prometheusadapters.monitoring.qubership.org
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
    oc delete crd prometheuses.monitoring.coreos.com
    oc delete crd prometheusrules.monitoring.coreos.com
    oc delete crd servicemonitors.monitoring.coreos.com
    oc delete crd thanosrulers.monitoring.coreos.com
    oc delete crd customscalemetricrules.monitoring.qubership.org
    oc delete crd platformmonitorings.monitoring.qubership.org
    oc delete crd prometheusadapters.monitoring.qubership.org
    ```

## Deploy/Skip/Remove based components

Monitoring-operator control (deploy/update/remove) some base components which include into monitoring deployment
Components list which it control:

* `prometheus-operator` - control deployment and settings
* `prometheus` - control only settings in Prometheus CR, deployment control by `prometheus-operator`
* `alertmanager` - control only settings in AlertManager CR, deployment control by `prometheus-operator`
* `grafana-operator` - control deployment and settings
* `grafana` - control only settings in Grafana, GrafanaDashboard and GrafanaDatasource CRs, deployment control by `grafana-operator`
* `kube-state-metrics` - control deployment and settings
* `node-exporter` - control deployment and settings

All these components will be deploy or remove in runtime after deploy. For these purposes in `PlatformMonitoring` CR
exists special parameters.

Some of these component's section has a special parameter `.<section>.install: [true|false]`. For example:

```yaml
alertManager:
  install: true
grafana:
  install: true
prometheus:
  install: true
kubeStateMetrics:
  install: true
nodeExporter:
  install: true
```

**Note:** By default all `.<section>.install` parameters set as `true`.

So you can skip components deploy, or remove it after deploy, or deploy if it was skipped with using these parameters.
For example, if I want to skip Grafana deployment I can specify:

```yaml
grafana:
  install: false
```

For see more information about available deploy parameters please refer to [Installation Guide](installation/README.md).

## Pause reconciliation for components which control by monitoring-operator

As described in paragraph above monitoring-operator control components and can create/update/remove them.
Also it means that if you want change any settings manually you just can't do it.

Monitoring-operator watch all objects managed by it and will restore expected settings after you will change it.

But in some cases it may be necessary to change some settings in any Deployment directly, for example when CR doesn't
allow specify necessary parameter.

In this case almost all components in PlatformMonitoring CR has a special parameter `.<section>.paused: [true|false]`.
For example:

```yaml
alertManager:
  paused: false
grafana:
  paused: false
  operator:
    paused: false
prometheus:
  paused: false
  operator:
    paused: false
kubeStateMetrics:
  paused: false
nodeExporter:
  paused: false
```

**Note:** By default all `.<section>.paused` parameters set as `false`.

**Warning:** Please pay attention that we are not recommending use this parameter except of cases of real necessity
or cases of development.

So you can pause reconciliation process for some components and change them deployments as you want.
For example, if I want to change Node-Exporter DaemonSet manually I can:

```yaml
nodeExporter:
  paused: true
```

and next make changes in DaemonSet.

For see more information about available deploy parameters please refer to [Installation Guide](installation/README.md).

## Export Prometheus Data

The information about the various ways to export data from Prometheus is described in the following sections.

### Export using Snapshots

Snapshot creates a snapshot of all current data into `snapshots/<datetime>-<rand>` under the TSDB's data directory and
returns the directory as a response. It optionally skips snapshotting the data that is only present in the head block,
and which has not yet been compacted to a disk.

```bash
POST /api/v1/admin/tsdb/snapshot
PUT /api/v1/admin/tsdb/snapshot
```

Following are the URL query parameters:

* `skip_head=<bool>`: Skip data present in the head block. Optional.

```bash
$ curl -XPOST http://localhost:9090/api/v1/admin/tsdb/snapshot
{
  "status": "success",
  "data": {
    "name": "20171210T211224Z-2be650b6d019eb54"
  }
}
```

The snapshot now exists at `<data-dir>/snapshots/20171210T211224Z-2be650b6d019eb54`.

### Copy PV Content

The Prometheus data can be copied as a content of PV.

To do so, just copy the PV data to any external storage. For example,

```bash
cp -R <data-dir>/ <external-storage>/prometheus/
```

## Import Prometheus Data

The information about the various ways to import data from Prometheus is described in the following sections.

### Import from Snapshots

To use the data from a snapshot, just copy the snapshot to `--storage.tsdb.path=<dir>` and run Prometheus. By default,
it is `--storage.tsdb.path=/prometheus`.

### Use early copied PV Content

To use such data, just copy it to `--storage.tsdb.path=<dir>` and run Prometheus. By default, it
is `--storage.tsdb.path=/prometheus`.
