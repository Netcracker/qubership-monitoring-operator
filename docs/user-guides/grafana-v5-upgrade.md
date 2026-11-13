# Grafana v5 Upgrade Limitations

The first upgrade from a release that uses Grafana Operator v4, including Qubership Monitoring Operator `v0.88.0`,
does not preserve every existing configuration behavior. Some fields remain accepted by the `PlatformMonitoring` API
but are not applied to the Grafana Operator v5 resources. Check the configuration against this page before upgrading.

Legacy dashboards keep working through the bundled
[Grafana converter](../installation/components/grafana-stack/grafana-converter.md), which is installed by default.

## Compatibility Boundary

The following settings and behaviors are not supported by the current Grafana Operator v5 integration:

- `dashboardLabelSelector` and `dashboardNamespaceSelector` have no effect. To choose which legacy dashboards are
  converted, set the
  [Grafana converter selectors](../installation/components/grafana-stack/grafana-converter.md#dashboard-selection).

These fields can remain stored in a `PlatformMonitoring` resource without affecting the generated Grafana v5
resources. The absence of a validation error does not mean that the setting is supported.

## Grafana settings applied on v5

The operator resets the Grafana admin password in an existing Grafana database when the password in the admin
credentials Secret changes.

`spec.auth`, including the OAuth endpoints and TLS settings, is applied to Grafana.

`spec.grafana.config` is applied to `grafana.ini`. Each section must be an object of scalar values. The operator
keeps its own values for the `auth.jwt` section and for `security.admin_user` and `security.admin_password`.

`spec.grafana.serviceAccount` labels and annotations are copied onto the Grafana ServiceAccount. When both maps are
empty, that ServiceAccount metadata is left unset.

The `grafana-ldap-config` Secret is mounted read-only at `/etc/grafana-secrets/grafana-ldap-config`. The mount does
not enable LDAP. The chart's default `ldap.toml` is empty, so LDAP login stays off until you provide the LDAP
configuration and, in `spec.grafana.config`, enable `auth.ldap` and set its `config_file` to
`/etc/grafana-secrets/grafana-ldap-config/ldap.toml`.

## Supported Upgrade Profile

The current upgrade path is intended for installations that:

- do not depend on the settings listed above;
- use the default Grafana dashboard selection;
- keep `grafanaConverter.install=true` until legacy dashboard conversion is verified;
- migrate every legacy `integreatly.org/v1alpha1` `GrafanaDataSource` to a
  `grafana.integreatly.org/v1beta1` `GrafanaDatasource` before removing the legacy resources.

Legacy `integreatly.org/v1alpha1` `GrafanaDashboard` resources remain supported through the
[Grafana converter](../installation/components/grafana-stack/grafana-converter.md). Do not remove the legacy dashboard
CRD while products still publish those resources. In non-privileged installs, list every product namespace with legacy
dashboards in `grafanaConverter.watchNamespaces`.

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
5. Confirm that `grafanaConverter.install` remains enabled. In non-privileged installs, set
   `grafanaConverter.watchNamespaces` to every namespace that holds legacy dashboards.

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

## Grafana 13

An upgrade from `v0.88.x` changes three components in one step: Grafana Operator v4 to v5, Grafana 12 to 13, and the
Grafana image renderer 3 to 5.

### Roll back

On its first start, Grafana 13 moves short URLs and user preferences to unified storage. Grafana 12 reads them from
the legacy tables and does not see later changes. `helm rollback` alone does not restore the previous state: also
restore the Grafana database backup from [Before Upgrading](#before-upgrading).

### Image renderer authentication

Grafana 13 authenticates the image renderer with JSON Web Tokens (JWTs) by default and rejects an empty token and the
default token `-` in that mode. The chart does not set a shared token, so when `grafana.imageRenderer.install` is
`true`, it sets `GF_FEATURE_TOGGLES_renderAuthJWT=false` to keep the previous authentication.

To use JWT authentication, set the same token on both sides:

```yaml
grafana:
  extraVarsSecret:
    GF_FEATURE_TOGGLES_renderAuthJWT: "true"
    GF_RENDERING_RENDERER_TOKEN: <token>
  imageRenderer:
    install: true
    extraEnvs:
      AUTH_TOKEN: <token>
```

### Settings removed in Grafana 13

Check your own data sources and Grafana settings for the following changes:

- Elasticsearch and Zipkin are no longer core data sources. Install the corresponding plugins to keep using them.
- The core Prometheus data source no longer supports Azure and SigV4 authentication.
- The `[marketplace]` configuration section is renamed to `[plugins_marketplace]`. Rename the section in
  `spec.grafana.config`, and replace `GF_MARKETPLACE_*` variables in `extraVars` with `GF_PLUGINS_MARKETPLACE_*`.
- The `grafana_api_dashboard_get_milliseconds` and `grafana_api_dashboard_save_milliseconds` metrics are removed.
  The Grafana Overview dashboard uses `grafana_http_request_duration_seconds` for the dashboard API instead.

### Default image for a hand-written custom resource

The chart always sets `spec.grafana.image`. A `PlatformMonitoring` resource written by hand without
`spec.grafana.image` gets the operator fallback `grafana/grafana:11.6.5`, which is a downgrade after the Grafana 13
migration. Set `spec.grafana.image` explicitly in such resources.

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
