# Grafana v5 Upgrade Limitations

The first upgrade from a release that uses Grafana Operator v4, including Qubership Monitoring Operator `v0.88.0`,
does not preserve every existing configuration behavior. Some fields remain accepted by the `PlatformMonitoring` API
but are not applied to the Grafana Operator v5 resources. Check the configuration against this page before upgrading.

> **Dependency:** Legacy dashboard compatibility requires the dashboard-only converter integration tracked in
> [PR #322](https://github.com/Netcracker/qubership-monitoring-operator/pull/322). Use this upgrade profile only with
> a target release that contains that integration.

## Compatibility Boundary

The following settings and behaviors are not supported by the current Grafana Operator v5 integration:

- Runtime changes to the Grafana admin credentials Secret are not synchronized to an existing Grafana database.
  Initial credentials are still used when the database is initialized. Follow
  [issue #375](https://github.com/Netcracker/qubership-monitoring-operator/issues/375).
- `spec.auth`, including OAuth endpoints and TLS settings, is not propagated to Grafana. Follow
  [issue #376](https://github.com/Netcracker/qubership-monitoring-operator/issues/376).
- `spec.grafana.config` is not propagated to Grafana. Follow
  [issue #377](https://github.com/Netcracker/qubership-monitoring-operator/issues/377).
- Grafana LDAP configuration, Grafana ServiceAccount labels and annotations, `dashboardLabelSelector`, and
  `dashboardNamespaceSelector` are not propagated. Follow
  [issue #434](https://github.com/Netcracker/qubership-monitoring-operator/issues/434).
- Jaeger and ClickHouse integrations do not automatically create Grafana datasources. Follow
  [issue #435](https://github.com/Netcracker/qubership-monitoring-operator/issues/435).

These fields can remain stored in a `PlatformMonitoring` resource without affecting the generated Grafana v5
resources. The absence of a validation error does not mean that the setting is supported.

## Supported Upgrade Profile

The current upgrade path is intended for installations that:

- do not depend on the settings and automatic datasource generation listed above;
- use the default Grafana dashboard selection;
- do not rotate the managed Grafana admin credentials after its database has been initialized;
- keep `grafanaConverter.install=true` until legacy dashboard conversion is verified;
- migrate every legacy `integreatly.org/v1alpha1` `GrafanaDataSource` to a
  `grafana.integreatly.org/v1beta1` `GrafanaDatasource` before removing the legacy resources.

Legacy `integreatly.org/v1alpha1` `GrafanaDashboard` resources remain supported through the
[Grafana operator converter](https://github.com/Netcracker/qubership-grafana-operator-converter). Do not remove the
legacy dashboard CRD while products still publish those resources.

The Monitoring Operator integration enables dashboard conversion only. Datasource, folder, and notification
conversion remain follow-up work in
[converter issue #88](https://github.com/Netcracker/qubership-grafana-operator-converter/issues/88). Do not enable
datasource conversion while following the manual migration procedure in this guide.

## Before Upgrading

1. Review the current `PlatformMonitoring` resource and Helm values for every setting listed under
   [Compatibility Boundary](#compatibility-boundary).
2. Inventory and migrate legacy Grafana datasources as described in the
   [Grafana datasource migration guide](../examples/custom-resources/grafana-datasource/README.md#migrate-v4-resources).
3. Back up the Grafana database or PVC and the current `PlatformMonitoring` resource. Record datasource UIDs that
   must remain stable.
4. Apply the complete CRD bundle from the target release with server-side apply before upgrading the operator.
   Ordinary `helm upgrade` does not update installed CRDs.
5. Confirm that the target release contains the converter integration and that `grafanaConverter.install` remains
   enabled.

### Argo CD staged upgrade

When Argo CD manages the CRD and operator charts as separate Applications, stage the first Grafana v5 upgrade without
pruning:

1. Disable automated sync and pruning for both Applications. Wait for any in-progress sync to finish.
2. Update and manually sync the CRD Application without pruning. This installs the v5 CRDs while retaining the v4
   Grafana, GrafanaDataSource, and GrafanaDashboard CRDs needed during migration.
3. Set `source.helm.skipCrds: true` on the operator Application. Update it to the same target revision, and manually
   sync it without pruning.
4. Wait for `PlatformMonitoring` to report `Successful=True`, the v5 Grafana resource to become ready, and every
   expected `GrafanaDatasource` and converted `GrafanaDashboard` to report successful synchronization.
5. Re-enable automated sync and pruning. Sync the operator Application first and the CRD Application afterward.

Do not delete the CRD Application or enable pruning before the migration checks pass. Keep the legacy
GrafanaDashboard CRD while products continue to publish legacy dashboard resources.

## After Upgrading

Verify all of the following before pruning legacy resources or CRDs:

- `PlatformMonitoring` reports `Successful=True`;
- the Grafana v5 resource and Grafana workload are ready;
- every expected `GrafanaDatasource` reports successful synchronization;
- required datasource UIDs are unchanged;
- every legacy dashboard has its expected converter-managed `grafana.integreatly.org/v1beta1` resource;
- converted dashboards report successful synchronization and are available in Grafana;
- the configured Grafana login method works; and
- the Grafana PVC or external database contains the expected data.

Keep the backup and the legacy Grafana v4 CRDs until these checks pass. Each limitation on this page should be removed
only after its linked issue is implemented and verified.
