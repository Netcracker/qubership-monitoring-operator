# Grafana Converter

The Grafana converter keeps legacy dashboards working after the move to Grafana Operator v5. It watches
`integreatly.org/v1alpha1` `GrafanaDashboard` resources and creates a matching `grafana.integreatly.org/v1beta1`
`GrafanaDashboard` for each one. Grafana Operator v5 then loads the converted dashboard into Grafana.

The chart installs the converter by default as the `grafanaConverter` subchart, from
[qubership-grafana-operator-converter](https://github.com/Netcracker/qubership-grafana-operator-converter). Keep it
installed while any product still publishes legacy dashboards.

## What the converter does

* It converts dashboards only. Datasource, folder, and notification conversion stay disabled; migrate legacy
  datasources manually as described in the [Grafana datasource migration guide][datasource-migration].
* The converted dashboard has the same name and namespace as the legacy dashboard. It copies the dashboard content
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
`openshift.io/cluster-monitoring=true`. To convert only some of them, set the Grafana dashboard selectors in the
monitoring Helm values:

```yaml
grafana:
  dashboardLabelSelector:
    - matchLabels:
        app.kubernetes.io/component: monitoring
  dashboardNamespaceSelector:
    matchLabels:
      monitoring: enabled
```

Monitoring Operator copies `dashboardLabelSelector` and `dashboardNamespaceSelector` from `PlatformMonitoring` into
the `grafana-resources-converter` ConfigMap and restarts the converter when the selectors change. This happens only
while Grafana is installed and not paused. A Helm upgrade can rewrite the ConfigMap, and the next reconciliation
restores the selectors. Do not set the selectors under `grafanaConverter.grafana.converter`, because Monitoring
Operator replaces them.

A dashboard matches when any one label selector matches. `dashboardNamespaceSelector` applies only when the converter
watches all namespaces. Helm merges your `dashboardNamespaceSelector` with the default one, so the
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
| namespaceOverride                                      | Must be empty or the release namespace. Any other value fails the render, because Monitoring Operator updates the converter in the release namespace.                               | string            |
| image.repository                                       | Converter image repository.                                                                                                                                                         | string            |
| image.tag                                              | Converter image tag. The default pins a multi-architecture image by digest.                                                                                                         | string            |
| image.pullPolicy                                       | Image pull policy. Default: `IfNotPresent`.                                                                                                                                         | string            |
| imagePullSecrets                                       | Image pull secrets for the converter pod.                                                                                                                                           | []object          |
| resources                                              | Converter container resources. Unset keys default to requests `50m` CPU, `64Mi` memory, `16Mi` ephemeral storage, and limits `200m` CPU, `256Mi` memory, `256Mi` ephemeral storage. | object            |
| securityContext                                        | Container security context. Default drops all capabilities and uses a read-only root file system as a non-root user.                                                                | object            |
| podSecurityContext                                     | Pod security context.                                                                                                                                                               | object            |
| podAnnotations                                         | Extra pod annotations.                                                                                                                                                              | map[string]string |
| additionalLabels                                       | Extra labels for every converter resource.                                                                                                                                          | map[string]string |
| nodeSelector, tolerations, affinity, priorityClassName | Pod scheduling settings.                                                                                                                                                            | object            |
| env                                                    | Extra environment variables. `WATCH_NAMESPACE` and `WATCH_NAMESPACE_SELECTOR` are reserved and fail the render.                                                                     | []object          |
| serviceMonitor.enabled                                 | Creates a ServiceMonitor for converter metrics. Default: `false`.                                                                                                                   | bool              |
| grafana.converter.deleteTargetOnSourceDeletion         | Deletes a converted dashboard when its legacy dashboard is deleted, through a Kubernetes owner reference. Default: `false`.                                                         | bool              |
| grafana.converter.resolveGzipConfigMapRef              | Converts legacy dashboards whose content is in a ConfigMap referenced by `spec.gzipConfigMapRef`. Grants read access to ConfigMaps. Default: `false`.                               | bool              |
| grafana.converter.gzipConfigMapMaxDecompressedSize     | Maximum decompressed size of content read through `gzipConfigMapRef`. Default: `32Mi`.                                                                                              | string            |
<!-- markdownlint-enable line-length -->

`watchNamespaceSelector` is rejected. The converter cannot change its namespace list at runtime; use
`watchNamespaces`.

The converter runs one replica with the `Recreate` update strategy, so two converters never write the same dashboards
during a rollout. A change to `grafanaConverter.grafana.converter` restarts the pod.

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
kubectl get grafanadashboards.grafana.integreatly.org <name> -n <namespace> \
  -o jsonpath='{.status.conditions[?(@.type=="DashboardSynchronized")].status}'
```

The command prints `True` when Grafana has loaded the dashboard.

## Troubleshooting

| Symptom                                                     | Check                                                                                                                        |
| ----------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| A legacy dashboard has no converted dashboard               | The namespace is in the watch scope, the dashboard matches the selectors, and the converter logs show no conversion error.   |
| The converter logs report a conflict                        | A v1beta1 dashboard with the same name exists without the `app.kubernetes.io/managed-by-operator` label. Rename one of them. |
| A dashboard with `gzipConfigMapRef` is reported as an error | Enable `grafanaConverter.grafana.converter.resolveGzipConfigMapRef`.                                                         |
| The converted dashboard exists but is not in Grafana        | The Grafana instance has the labels from `instanceSelector`, and the dashboard's `DashboardSynchronized` condition.          |
| Selector changes have no effect                             | Set them in `grafana.dashboardLabelSelector` and `grafana.dashboardNamespaceSelector`, and check that Grafana is not paused. |

[datasource-migration]: ../../../examples/custom-resources/grafana-datasource/README.md#migrate-v4-resources
