# Grafana Converter

The Grafana converter keeps legacy dashboards working after the move to Grafana Operator v5. It watches
`integreatly.org/v1alpha1` `GrafanaDashboard` resources and creates a matching `grafana.integreatly.org/v1beta1`
`GrafanaDashboard` for each one. Grafana Operator v5 then loads the converted dashboard into Grafana.

The chart installs the converter by default as the `grafanaConverter` subchart, from
[qubership-grafana-operator-converter](https://github.com/Netcracker/qubership-grafana-operator-converter). Keep it
installed while any product still publishes legacy dashboards.

The Monitoring Operator and its chart also publish the built-in dashboards as legacy `GrafanaDashboard` resources, so
they reach Grafana through the converter. With `grafanaConverter.install=false`, or with the release namespace outside
the converter watch scope, Grafana has no built-in dashboards.

Uninstalling the release deletes every converted dashboard in the release namespace, because Grafana is removed with
it. Converted dashboards in other namespaces stay.

## What the converter does

* It converts dashboards only. Datasource, folder, and notification conversion stay disabled; migrate legacy
  datasources manually as described in the [Grafana datasource migration guide][datasource-migration].
* The converted dashboard is named `<legacy dashboard name>-converted` and is created in the namespace of the legacy
  dashboard. The suffix keeps it from colliding with a v1beta1 dashboard that has the name of the legacy dashboard. It
  copies the dashboard content
  (`json`, `gzipJson`, `url`, `jsonnet`, or `configMapRef`), plugins, and datasources, and maps `customFolderName` to
  `folderTitle`.
* The converted dashboard selects the monitoring Grafana instance through
  `instanceSelector` `app.kubernetes.io/component: grafana` and `app.kubernetes.io/part-of: monitoring`, and sets
  `allowCrossNamespaceImport: true`, so dashboards in product namespaces reach Grafana.
* The converter labels every resource it creates with
  `app.kubernetes.io/managed-by-operator=grafana-operator-converter` and updates only resources with that label. A
  v1beta1 dashboard without the label is left unchanged, and the converter reports a conflict.
* The legacy dashboard is never changed. Keep the legacy `GrafanaDashboard` CRD installed while products publish
  legacy dashboards.
* Updating a legacy dashboard updates its converted dashboard. Deleting a legacy dashboard keeps the converted one
  unless `deleteTargetOnSourceDeletion` is enabled.

## Watch scope

The converter watch scope follows the monitoring security mode. Its RBAC always matches the namespaces it watches.

<!-- markdownlint-disable line-length -->
| Settings                                                  | Watched namespaces                           | RBAC                                   |
| --------------------------------------------------------- | -------------------------------------------- | -------------------------------------- |
| `global.privilegedRights: true`, `watchNamespaces` empty  | All namespaces                               | ClusterRole and ClusterRoleBinding     |
| `watchNamespaces: "ns-a,ns-b"`                            | The listed namespaces                        | Role and RoleBinding in each namespace |
| `global.privilegedRights: false`, `watchNamespaces` empty | The release namespace                        | Role and RoleBinding                   |
| `grafanaConverter.namespaceScope: true`                   | The release namespace                        | Role and RoleBinding                   |
| `global.namespaceScope: true`                             | The release namespace, even if a list is set | Role and RoleBinding                   |
<!-- markdownlint-enable line-length -->

In non-privileged mode, legacy dashboards in product namespaces are converted only when you list those namespaces in
`grafanaConverter.watchNamespaces`. Include the release namespace in the list if it also contains legacy dashboards.

In a [namespaced guest install](../../../examples/deploy-parameters/namespaced-guest/README.md), set the host
`grafanaConverter.watchNamespaces` to the same list as the host Grafana operator. The guest converter watches only the
guest namespace.

## Dashboard selection

The chart defaults convert every legacy dashboard in the watched namespaces, except in namespaces labeled
`openshift.io/cluster-monitoring=true`. To convert only some of them, set the converter selectors in the monitoring
Helm values:

```yaml
grafanaConverter:
  config:
    dashboard:
      labelSelector:
        - matchLabels:
            app.kubernetes.io/component: monitoring
      namespaceSelector:
        matchLabels:
          monitoring: enabled
```

Helm writes the selectors to the `grafana-resources-converter` ConfigMap, and a change restarts the converter. The
converter does not read `grafana.dashboardLabelSelector` or `grafana.dashboardNamespaceSelector` from
`PlatformMonitoring`. If you set those for Grafana, set the same values here.

A dashboard matches when any one label selector matches. `namespaceSelector` applies only when the converter watches
all namespaces. Helm merges your `namespaceSelector` with the default one, so the
`openshift.io/cluster-monitoring` exclusion stays unless you set `matchExpressions: []`. When a legacy dashboard stops
matching, the converter deletes its converted dashboard and leaves the legacy one in place.

These selectors control which legacy dashboards are converted. They do not filter native
`grafana.integreatly.org/v1beta1` dashboards, which Grafana Operator selects through `instanceSelector`.

## Parameters

Parameters are set under `grafanaConverter` in the monitoring Helm values.

<!-- markdownlint-disable line-length -->
| Field                                                  | Description                                                                                                                                                                         | Scheme            |
| ------------------------------------------------------ | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------- |
| install                                                | Installs the converter. Default: `true`.                                                                                                                                            | bool              |
| watchNamespaces                                        | Comma-separated namespaces to watch. Empty watches all namespaces in privileged mode and the release namespace otherwise.                                                           | string            |
| namespaceScope                                         | Watches only the release namespace when `watchNamespaces` is empty. Default: `false`.                                                                                               | bool              |
| image                                                  | Converter image. The default pins a multi-architecture image by digest.                                                                                                             | string            |
| imagePullPolicy                                        | Image pull policy. Default: `IfNotPresent`.                                                                                                                                         | string            |
| imagePullSecrets                                       | Image pull secrets for the converter pod.                                                                                                                                           | []object          |
| resources                                              | Converter container resources. Unset keys default to requests `50m` CPU, `64Mi` memory, `16Mi` ephemeral storage, and limits `200m` CPU, `256Mi` memory, `256Mi` ephemeral storage. | object            |
| securityContext                                        | Container security context. Default drops all capabilities and uses a read-only root file system as a non-root user.                                                                | object            |
| podSecurityContext                                     | Pod security context.                                                                                                                                                               | object            |
| podAnnotations                                         | Extra pod annotations.                                                                                                                                                              | map[string]string |
| additionalLabels                                       | Extra labels for every converter resource.                                                                                                                                          | map[string]string |
| nodeSelector, tolerations, affinity, priorityClassName | Pod scheduling settings.                                                                                                                                                            | object            |
| env                                                    | Extra environment variables. `WATCH_NAMESPACE` and `WATCH_NAMESPACE_SELECTOR` are reserved and fail the render.                                                                     | []object          |
| serviceMonitor.enabled                                 | Creates a ServiceMonitor for converter metrics. Default: `false`.                                                                                                                   | bool              |
| config.instanceSelector                                | Labels of the Grafana instance that receives every converted resource.                                                                                                              | object            |
| config.dashboard.enabled                               | Converts legacy dashboards. Default: `true`.                                                                                                                                        | bool              |
| config.dashboard.labelSelector                         | Legacy dashboards to convert, as a list of label selectors. See [Dashboard selection](#dashboard-selection).                                                                        | []object          |
| config.dashboard.namespaceSelector                     | Namespace label selector for legacy dashboards. See [Dashboard selection](#dashboard-selection).                                                                                    | object            |
| config.dashboard.deleteTargetOnSourceDeletion          | Deletes a converted dashboard when its legacy dashboard is deleted, through a Kubernetes owner reference. Default: `false`.                                                         | bool              |
| config.dashboard.gzipConfigMapRef.enabled              | Converts legacy dashboards whose content is in a ConfigMap referenced by `spec.gzipConfigMapRef`. Grants read access to ConfigMaps. Default: `false`.                               | bool              |
| config.dashboard.gzipConfigMapRef.maxDecompressedSize  | Maximum decompressed size of content read through `gzipConfigMapRef`. Default: `32Mi`.                                                                                              | string            |
| config.datasource.enabled, config.folder.enabled       | Converts legacy data sources or folders. Needs the matching legacy CRD, which the chart does not install. Default: `false`.                                                         | bool              |
| config.notification.enabled                            | Converts legacy notification channels to contact points. Needs the legacy CRD, which the chart does not install. Default: `false`.                                                  | bool              |
<!-- markdownlint-enable line-length -->

`watchNamespaceSelector` is rejected. The converter cannot change its namespace list at runtime; use
`watchNamespaces`.

The converter runs one replica with the `Recreate` update strategy, so two converters never write the same dashboards
during a rollout. A change to `grafanaConverter.config` restarts the pod.

Enable at least one converter under `grafanaConverter.config`. To run none, set `grafanaConverter.install: false`.
The chart rejects unknown keys under `grafanaConverter.config`.

## Examples

* [Deploy parameters](../../../examples/deploy-parameters/grafana-converter/README.md) for privileged,
  non-privileged, and disabled installs.
* [Legacy Grafana dashboard](../../../examples/custom-resources/legacy-grafana-dashboard/README.md) with the
  converted dashboard it produces.

## Check the conversion

List the converted dashboards and their synchronization status:

```bash
kubectl get grafanadashboards.grafana.integreatly.org -A \
  -l app.kubernetes.io/managed-by-operator=grafana-operator-converter
```

Check one converted dashboard:

```bash
kubectl get grafanadashboards.grafana.integreatly.org <legacy-name>-converted -n <namespace> \
  -o jsonpath='{.status.conditions[?(@.type=="DashboardSynchronized")].status}'
```

The command prints `True` when Grafana has loaded the dashboard.

## Replace a converted dashboard with a v1beta1 dashboard

When the dashboard JSON sets `uid`, a v1beta1 dashboard that you publish yourself and the converted copy write the
same dashboard in Grafana. Remove the copy as part of the migration:

1. Publish the v1beta1 dashboard.
2. Delete the legacy dashboard. While it exists, the converter recreates the copy.
3. Delete the `<legacy dashboard name>-converted` dashboard.

Grafana Operator can remove the dashboard from Grafana when the copy is deleted. It restores the dashboard on the
next synchronization of the v1beta1 dashboard that remains.

To delete every converted dashboard at once, select them by label:

```bash
kubectl delete grafanadashboards.grafana.integreatly.org -A \
  -l app.kubernetes.io/managed-by-operator=grafana-operator-converter
```

The converter recreates the copies whose legacy dashboards still exist.

A build that predates the suffix named converted dashboards like their legacy dashboards. The converter does not
update or delete those copies, so each one diverges from its `-converted` replacement while both write the same
dashboard in Grafana. After an upgrade from such a build, delete them in every namespace that holds legacy dashboards:

```bash
kubectl get grafanadashboards.grafana.integreatly.org -n <namespace> -o name \
  -l app.kubernetes.io/managed-by-operator=grafana-operator-converter \
  | grep -v -- '-converted$' | xargs -r kubectl delete -n <namespace>
```

The filter keeps names that end with `-converted`. An earlier build can also have made such a copy for a legacy
dashboard that is itself named `<name>-converted`. If a legacy dashboard `<name>` exists, the converter reuses that copy
as the copy of `<name>`, so do not delete it. Delete it by name only when no legacy dashboard `<name>` exists.

A legacy dashboard with a name longer than 243 characters gets no `-converted` copy, because the converted name
exceeds the Kubernetes limit of 253 characters, so the command deletes its only copy. Before you run the command,
republish that dashboard under a shorter name and wait for its `-converted` copy.

When the dashboard JSON sets no `uid`, Grafana Operator uses the `metadata.uid` of the converted dashboard as the
dashboard UID in Grafana. The `-converted` replacement is a new resource, so it creates a second dashboard in Grafana,
and deleting the old copy removes the UID that existing links use. The upgraded converter can also delete an old
copy on startup, so do steps 1 and 2 before you upgrade the converter:

1. Read the UID from `metadata.uid` of the old copy:

   ```bash
   kubectl get grafanadashboards.grafana.integreatly.org <name> -n <namespace> -o jsonpath='{.metadata.uid}'
   ```

2. Set that value as `uid` in the dashboard JSON of the legacy dashboard: in `spec.json`, in `spec.gzipJson`, or in the
   ConfigMap that the dashboard references.
3. After the upgrade, check that `status.uid` of the `<name>-converted` dashboard has the value from step 1.

A dashboard that is loaded from `spec.url` or `spec.grafanaCom` keeps its UID only when the JSON at that source sets
`uid`.

## Troubleshooting

| Symptom                                                     | Check                                                                                                                        |
| ----------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| A legacy dashboard has no converted dashboard               | The namespace is in the watch scope, the dashboard matches the selectors, and the converter logs show no conversion error.   |
| Grafana has no built-in dashboards                          | The converter is installed, and the release namespace is in its watch scope and matches `namespaceSelector`.                 |
| The converter logs report a conflict                        | A `<legacy name>-converted` v1beta1 dashboard exists without the `app.kubernetes.io/managed-by-operator` label. Rename it.   |
| A dashboard with `gzipConfigMapRef` is reported as an error | Enable `grafanaConverter.config.dashboard.gzipConfigMapRef.enabled`.                                                         |
| The converted dashboard exists but is not in Grafana        | The Grafana instance has the labels from `instanceSelector`, and the dashboard's `DashboardSynchronized` condition.          |
| Selector changes have no effect                             | Set them in `grafanaConverter.config.dashboard`, not in `grafana`.                                                           |

[datasource-migration]: ../../../examples/custom-resources/grafana-datasource/README.md#migrate-v4-resources
