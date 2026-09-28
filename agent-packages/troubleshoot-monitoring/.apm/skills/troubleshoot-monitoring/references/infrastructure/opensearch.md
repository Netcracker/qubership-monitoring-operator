# Infrastructure component troubleshooting


## opensearch

### OpenSearch monitoring is not installed

**Symptoms:**

- OpenSearch Grafana dashboards have no `opensearch_*` series from this chart.
- Deployment `opensearch-monitoring` is absent.

**Root cause:**

Telegraf is a dedicated Deployment `{fullname}-monitoring` (default `opensearch-monitoring`), not a sidecar on
OpenSearch pods. Chart `monitoring.enabled` defaults **false**. Helper uses `MONITORING_ENABLED` only when
`global.cloudIntegrationEnabled` is true. The installation.md table that lists default `true` is not the chart default
at this tag.

**Applies to:** opensearch when OpenSearch Prometheus metrics are expected. Image and chart versions are supporting
facts.

**Failed layer:** `component configuration missing`

**How to check:**

1. Confirm that Deployment `opensearch-monitoring` is absent.
2. Name `monitoring.enabled` or `MONITORING_ENABLED` only when effective values show monitoring is off. Prefer chart
   values over the installation.md default table.

**How to fix:**

1. Set `monitoring.enabled: true` when that key is present, then upgrade.
2. Confirm that Deployment and Service `opensearch-monitoring` exist.

Do not change Monitoring Operator selectors until the monitoring workload exists.

**Owner:** opensearch (`Netcracker/qubership-opensearch`)

**Risk:** Enabling monitoring deploys Telegraf and starts exec collection against OpenSearch.

**Provenance:** Maintainer-verified at `Netcracker/qubership-opensearch` tag `2.6.7`
(`7ebb413d4ddfcaa6aad6e969f41496a88ea21daa`). That tag is not an applicability bound.

- `operator/charts/helm/opensearch-service/values.yaml`
- `operator/charts/helm/opensearch-service/templates/_helpers.tpl`
- `operator/charts/helm/opensearch-service/templates/monitoring/deployment.yaml`
- `docs/public/installation.md`


### OpenSearch monitoring type is InfluxDB

**Symptoms:**

- Deployment `opensearch-monitoring` exists.
- ServiceMonitor `opensearch-service-monitor` and GrafanaDashboard CRs are absent because Telegraf writes InfluxDB
  output instead of Prometheus on `:8096`.

**Root cause:**

SM and Grafana templates require `monitoring.monitoringType` not equal to `influxdb` (default `prometheus`). The
workload can exist without a Prometheus monitor.

**Applies to:** opensearch when the monitoring Deployment exists. Image and chart versions are supporting facts.

**Failed layer:** `monitor missing or misconfigured`

**How to check:**

1. Confirm that `opensearch-monitoring` exists and `opensearch-service-monitor` is absent.
2. Name `monitoring.monitoringType` only when effective values show `influxdb`.

**How to fix:**

1. Set `monitoring.monitoringType: prometheus` when Prometheus scrape is required and that key is present, then
   upgrade.
2. Confirm that `opensearch-service-monitor` scrapes port `prometheus-cli`.

**Owner:** opensearch (`Netcracker/qubership-opensearch`)

**Risk:** Switching to Prometheus creates a scrape target and dashboard CRs.

**Provenance:** Maintainer-verified at `Netcracker/qubership-opensearch` tag `2.6.7`
(`7ebb413d4ddfcaa6aad6e969f41496a88ea21daa`). That tag is not an applicability bound.

- `operator/charts/helm/opensearch-service/templates/monitoring/service-monitor.yaml`
- `operator/charts/helm/opensearch-service/templates/monitoring/configuration.yaml`


### OpenSearch Grafana dashboard is absent

**Symptoms:**

- Monitoring Deployment and ServiceMonitor exist.
- GrafanaDashboard CR `opensearch-grafana-dashboard` is absent, so bundled OpenSearch Monitoring never loads.

**Root cause:**

The main dashboard CR is gated by `monitoring.installDashboard` (default true), monitoring enabled, and type not
`influxdb`. This is a missing dashboard CR, not a scrape failure.

**Applies to:** opensearch when bundled OpenSearch Grafana panels never load. Image and chart versions are supporting
facts.

**Failed layer:** `dashboard query mismatch`

**How to check:**

1. Confirm that `opensearch-grafana-dashboard` is absent while `opensearch-service-monitor` exists.
2. Name `monitoring.installDashboard` only when effective values show it is false.

**How to fix:**

1. Set `monitoring.installDashboard: true` when that key is present, then upgrade.
2. Confirm that CR `opensearch-grafana-dashboard` exists.

**Owner:** opensearch (`Netcracker/qubership-opensearch`)

**Risk:** Creating the GrafanaDashboard CR loads bundled OpenSearch panels into Grafana.

**Provenance:** Maintainer-verified at `Netcracker/qubership-opensearch` tag `2.6.7`
(`7ebb413d4ddfcaa6aad6e969f41496a88ea21daa`). That tag is not an applicability bound.

- `operator/charts/helm/opensearch-service/templates/monitoring/grafana-dashboard.yaml`
- `docs/public/monitoring.md`


### OpenSearch indices dashboard is not enabled

**Symptoms:**

- Main OpenSearch monitoring works.
- OpenSearch Indices dashboard is missing, or series `opensearch_indices_stats_*` are absent.

**Root cause:**

`monitoring.includeIndices` (default false) gates both GrafanaDashboard `opensearch-indices-grafana-dashboard` and
Telegraf `indices_include = ["_all"]`. Series `opensearch_indices_version_created_count` can still exist.

**Applies to:** opensearch when per-index stats are expected. Image and chart versions are supporting facts.

**Failed layer:** `dashboard query mismatch`

**How to check:**

1. Confirm that main monitoring is scraped and `opensearch-indices-grafana-dashboard` is absent or
   `opensearch_indices_stats_*` is absent.
2. Name `monitoring.includeIndices` only when effective values show it is false.

**How to fix:**

1. Set `monitoring.includeIndices: true` when that key is present, then upgrade.
2. Confirm that CR `opensearch-indices-grafana-dashboard` exists and `opensearch_indices_stats_*` series appear.

**Owner:** opensearch (`Netcracker/qubership-opensearch`)

**Risk:** Per-index stats increase Telegraf collection cost and metric cardinality.

**Provenance:** Maintainer-verified at `Netcracker/qubership-opensearch` tag `2.6.7`
(`7ebb413d4ddfcaa6aad6e969f41496a88ea21daa`). That tag is not an applicability bound.

- `operator/charts/helm/opensearch-service/templates/monitoring/indices-dashboard.yaml`
- `operator/charts/helm/opensearch-service/templates/monitoring/configuration.yaml`
- `docs/public/monitoring/indices-dashboard.md`


### OpenSearch slow queries are not enabled

**Symptoms:**

- OpenSearch Slow Queries dashboard is empty.
- GrafanaDashboard CR `opensearch-slow-queries-grafana-dashboard` is absent, or series
  `opensearch_slow_query_took_millis` is absent.

**Root cause:**

`monitoring.slowQueries.enabled` defaults false and requires `global.externalOpensearch.enabled` to be false. Docs
state slow queries do not work on AWS or external OpenSearch. An enabled dashboard that is empty because no slow logs
exist in the interval is not this case.

**Applies to:** opensearch when slow-query metrics are expected. Image and chart versions are supporting facts.

**Failed layer:** `dashboard query mismatch`

**How to check:**

1. Confirm that `opensearch-slow-queries-grafana-dashboard` is absent or `opensearch_slow_query_took_millis` is absent.
2. Name `monitoring.slowQueries.enabled` only when effective values show it is false. Do not use this case for
   external OpenSearch.

**How to fix:**

1. Set `monitoring.slowQueries.enabled: true` when that key is present and the cluster is not external, then upgrade.
2. Confirm that CR `opensearch-slow-queries-grafana-dashboard` exists.

**Owner:** opensearch (`Netcracker/qubership-opensearch`)

**Risk:** Enabling slow-query collection reads slow logs on a schedule and adds dashboard load.

**Provenance:** Maintainer-verified at `Netcracker/qubership-opensearch` tag `2.6.7`
(`7ebb413d4ddfcaa6aad6e969f41496a88ea21daa`). That tag is not an applicability bound.

- `operator/charts/helm/opensearch-service/templates/monitoring/slow-queries-dashboard.yaml`
- `docs/public/monitoring/slow-queries-dashboard.md`
- `docs/public/installation.md`


### OpenSearch replication dashboard is not enabled

**Symptoms:**

- OpenSearch Replication dashboard is empty.
- GrafanaDashboard CR `opensearch-replication-grafana-dashboard` is absent, or series
  `opensearch_replication_metric_*` are absent.

**Root cause:**

Helper `opensearch.enableDisasterRecovery` is true only when `global.disasterRecovery.mode` is `active`, `standby`, or
`disabled`. Values default `mode: ""` is false. This gates the replication dashboard, `replication_metric.py`, and
replication alerts. Lag **alert** additionally needs `monitoring.thresholds.lagAlert` greater than -1; that is
alerting, not dashboard enablement.

**Applies to:** opensearch when replication metrics are expected. Image and chart versions are supporting facts.

**Failed layer:** `dashboard query mismatch`

**How to check:**

1. Confirm that `opensearch-replication-grafana-dashboard` is absent or `opensearch_replication_metric_*` is absent.
2. Name `global.disasterRecovery.mode` only when effective values show disaster recovery is unset.

**How to fix:**

1. Set `global.disasterRecovery.mode` to `active`, `standby`, or `disabled` when that key is present and DR is
   intended, then upgrade.
2. Confirm that CR `opensearch-replication-grafana-dashboard` exists.

**Owner:** opensearch (`Netcracker/qubership-opensearch`)

**Risk:** Enabling disaster-recovery mode changes OpenSearch DR behavior, not only the dashboard.

**Provenance:** Maintainer-verified at `Netcracker/qubership-opensearch` tag `2.6.7`
(`7ebb413d4ddfcaa6aad6e969f41496a88ea21daa`). That tag is not an applicability bound.

- `operator/charts/helm/opensearch-service/templates/_helpers.tpl`
- `operator/charts/helm/opensearch-service/templates/monitoring/replication-dashboard.yaml`
- `docs/public/monitoring/replication-dashboard.md`


### OpenSearch DBaaS adapter is not installed

**Symptoms:**

- DBaaS Adapter Status Down or alert `OpenSearchDBaaSIsDownAlert`.
- Deployment `dbaas-opensearch-adapter` is absent.

**Root cause:**

`dbaasAdapter.enabled` defaults false (helper `dbaas.enabled` uses `DBAAS_ENABLED` when cloud integration is on).
Telegraf still runs `dbaas_health_metric.py` whenever monitoring is on; a missing adapter yields present series
`opensearch_dbaas_health_status=1`, not an empty series. Do not treat status=1 as missing metrics.

**Applies to:** opensearch when DBaaS adapter health is expected. Image and chart versions are supporting facts.

**Failed layer:** `component configuration missing`

**How to check:**

1. Confirm that Deployment `dbaas-opensearch-adapter` is absent.
2. Name `dbaasAdapter.enabled` or `DBAAS_ENABLED` only when effective values show the adapter is off.

**How to fix:**

1. Set `dbaasAdapter.enabled: true` when that key is present, then upgrade.
2. Confirm that Deployment `dbaas-opensearch-adapter` exists.

**Owner:** opensearch (`Netcracker/qubership-opensearch`)

**Risk:** Enabling the DBaaS adapter deploys `dbaas-opensearch-adapter`.

**Provenance:** Maintainer-verified at `Netcracker/qubership-opensearch` tag `2.6.7`
(`7ebb413d4ddfcaa6aad6e969f41496a88ea21daa`). That tag is not an applicability bound.

- `operator/charts/helm/opensearch-service/templates/_helpers.tpl`
- `operator/charts/helm/opensearch-service/templates/opensearch-dbaas-adapter/dbaas-adapter-deployment.yaml`
- `docs/public/monitoring.md`
