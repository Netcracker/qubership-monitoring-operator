# Infrastructure component troubleshooting


## clickhouse-operator-helm

### ClickHouse ServiceMonitor is not enabled

**Symptoms:**

- ClickHouse Metrics Dashboard panels have no data.
- ServiceMonitor `clickhouse-service-monitor` and GrafanaDashboard CRs are absent, while Deployment `clickhouse` with
  container `metrics-exporter` may still run.

**Root cause:**

The Altinity metrics-exporter sidecar is always on Deployment `clickhouse`. Prometheus scrape and dashboard CRs render
only when helper `monitoring.install` is true: `clickhouseCluster.serviceMonitor` is `true` or `enable`, or
`MONITORING_ENABLED` with `global.cloudIntegrationEnabled`. Chart default `clickhouseCluster.serviceMonitor` is false.
This is scrape and dashboard off, not a missing exporter container.

**Applies to:** clickhouse-operator-helm when ClickHouse Prometheus metrics are expected. Image and chart versions are
supporting facts.

**Failed layer:** `monitor missing or misconfigured`

**How to check:**

1. Confirm that `clickhouse-service-monitor` and the ClickHouse GrafanaDashboard CRs are absent.
2. Name `clickhouseCluster.serviceMonitor` or `MONITORING_ENABLED` only when effective or rendered values show
   monitoring.install is false.

**How to fix:**

1. Set `clickhouseCluster.serviceMonitor: true` when that key is present in the evidence, then upgrade.
2. Confirm that `clickhouse-service-monitor` selects `clickhouse.altinity.com/app: chop` and port `metrics`.

Do not change Monitoring Operator selectors until the ServiceMonitor exists.

**Owner:** clickhouse-operator-helm (`Netcracker/qubership-clickhouse-operator-helm`)

**Risk:** Enabling the ServiceMonitor adds scrape of `clickhouse-metrics` and loads bundled ClickHouse dashboards.

**Provenance:** Maintainer-verified at `Netcracker/qubership-clickhouse-operator-helm` tag `0.49.1`
(`e1feb432807803c872b9bbbf45fe03225c19e34d`). That tag is not an applicability bound.

- `helm/clickhouse/templates/_helpers.tpl`
- `helm/clickhouse/values.yaml`
- `helm/clickhouse/templates/clickhouse-operator/clickhouse-service-monitor.yaml`


### ClickHouse cluster is not installed

**Symptoms:**

- Monitoring is enabled, but ServiceMonitor `clickhouse-service-monitor` is absent.
- GrafanaDashboard CRs may still exist because they gate on `monitoring.install` only.

**Root cause:**

ServiceMonitors require `clickhouseCluster.install` as well as `monitoring.install`. Cluster default is `yes`. This is
a missing cluster scrape resource, not ServiceMonitor disabled.

**Applies to:** clickhouse-operator-helm when monitoring is enabled. Image and chart versions are supporting facts.

**Failed layer:** `monitor missing or misconfigured`

**How to check:**

1. Confirm that `monitoring.install` is true and `clickhouse-service-monitor` is absent.
2. Name `clickhouseCluster.install` only when effective values show the cluster is off.

**How to fix:**

1. Set `clickhouseCluster.install` to the chart's enabled form when that key is present, then upgrade.
2. Confirm that Service `clickhouse-metrics` and `clickhouse-service-monitor` exist.

**Owner:** clickhouse-operator-helm (`Netcracker/qubership-clickhouse-operator-helm`)

**Risk:** Installing the cluster deploys ClickHouse and adds scrape targets.

**Provenance:** Maintainer-verified at `Netcracker/qubership-clickhouse-operator-helm` tag `0.49.1`
(`e1feb432807803c872b9bbbf45fe03225c19e34d`). That tag is not an applicability bound.

- `helm/clickhouse/templates/clickhouse-operator/clickhouse-service-monitor.yaml`


### ClickHouse collection scrape timeout

**Symptoms:**

- ServiceMonitor `clickhouse-service-monitor` exists and is discovered, but the scrape ends before the exporter
  finishes.
- The target is down or partial.

**Root cause:**

This case is chart `clickhouseCluster.serviceMonitorScrapeTimeout` on the ClickHouse ServiceMonitor endpoint (template
default `30s`). It is not the shared-discovery target-unavailable diagnosis.

**Applies to:** clickhouse-operator-helm when `clickhouse-service-monitor` is discovered. Image and chart versions are
supporting facts.

**Failed layer:** `target unavailable`

**How to check:**

1. Confirm the expected ServiceMonitor exists and is in discovery scope.
2. Confirm the scrape timeout from the live ServiceMonitor YAML or effective values
   (`clickhouseCluster.serviceMonitorScrapeTimeout`).
3. Confirm the scrape ends before gathering finishes, not that endpoints are missing.

**How to fix:**

1. Raise the ServiceMonitor `scrapeTimeout` that the evidence named, then upgrade.
2. Do not use `clickhouseCluster.serviceMonitorScrapeInterval` as the fix.
3. Do not change Monitoring Operator selectors for a scrape that already discovers the monitor.

**Owner:** clickhouse-operator-helm (`Netcracker/qubership-clickhouse-operator-helm`)

**Risk:** Raising scrape timeout keeps scrape workers busy longer on slow targets.

**Provenance:** Maintainer-verified at `Netcracker/qubership-clickhouse-operator-helm` tag `0.49.1`
(`e1feb432807803c872b9bbbf45fe03225c19e34d`). That tag is not an applicability bound.

- `helm/clickhouse/templates/clickhouse-operator/clickhouse-service-monitor.yaml`
- `helm/clickhouse/values.yaml`


### ClickHouse metrics exporter cannot fetch from ClickHouse

**Symptoms:**

- `up` for container `metrics-exporter` is 1, but cluster panels stay empty.
- Series `chi_clickhouse_metric_fetch_errors` increments, or alerts `ClickHouseMetricsExporterFetchErrors` /
  `ClickHouseServerDown` fire.

**Root cause:**

The exporter sidecar is scraped, so this is failed fetch into ClickHouse, not a missing ServiceMonitor. Credentials
are `clickhouseOperator.credentialsToInstances.chUsername`, `chPassword`, `chPort`, and `chCredentialsSecretName`
(Secret `clickhouse-operator-credentials`). Do not read Secret values.

**Applies to:** clickhouse-operator-helm when the metrics exporter is scraped. Image and chart versions are supporting
facts.

**Failed layer:** `metric series absent`

**How to check:**

1. Confirm that `clickhouse-service-monitor` is scraped and the metrics-exporter target is up.
2. Request non-secret evidence that the ClickHouse metrics user exists. Do not read Secret values.

**How to fix:**

1. Align `clickhouseOperator.credentialsToInstances` with a user that can query system tables, without copying Secret
   values into the report.
2. Do not enable `clickhouseCluster.serviceMonitor` as the fix for a healthy scrape with fetch errors.

**Owner:** clickhouse-operator-helm (`Netcracker/qubership-clickhouse-operator-helm`)

**Risk:** Changing the metrics user changes ClickHouse access. Do not copy Secret values into the diagnosis report.

**Provenance:** Maintainer-verified at `Netcracker/qubership-clickhouse-operator-helm` tag `0.49.1`
(`e1feb432807803c872b9bbbf45fe03225c19e34d`). That tag is not an applicability bound.

- `helm/clickhouse/values.yaml`
- `helm/clickhouse/templates/clickhouse-operator/clickhouse-operator-credentials.yaml`
- `helm/clickhouse/templates/clickhouse-operator/clickhouse_prometheus_alert_rules.yaml`


### ClickHouse performance dashboard datasource is missing

**Symptoms:**

- GrafanaDashboard CR `clickhouse-performance-grafana-dashboard` exists.
- ClickHouse Performance Dashboard variables `$db` or `vertamedia-clickhouse-datasource` stay empty, while Prometheus
  `$datasource` panels may still work.

**Root cause:**

That dashboard queries Grafana datasource type `vertamedia-clickhouse-datasource`. Monitoring Operator can create a
ClickHouse datasource from `integration.clickHouse.createGrafanaDataSource`. This is a dashboard datasource mismatch,
not a missing ServiceMonitor.

**Applies to:** clickhouse-operator-helm when the performance dashboard CR exists. Image and chart versions are
supporting facts.

**Failed layer:** `dashboard query mismatch`

**How to check:**

1. Confirm that `clickhouse-performance-grafana-dashboard` exists and Prometheus panels are not the empty ones.
2. Confirm whether a Grafana datasource of type `vertamedia-clickhouse-datasource` is present. Name
   `integration.clickHouse.createGrafanaDataSource` only when Monitoring Operator values show it.

**How to fix:**

1. Ask Monitoring Operator owners to create the ClickHouse Grafana datasource when that integration key is present.
2. Do not enable `clickhouseCluster.serviceMonitor` as the fix for a missing ClickHouse datasource.

**Owner:** Monitoring Operator Grafana integration (`Netcracker/qubership-monitoring-operator`)

**Risk:** Creating a ClickHouse datasource exposes ClickHouse query access in Grafana.

**Provenance:** Maintainer-verified at `Netcracker/qubership-clickhouse-operator-helm` tag `0.49.1`
(`e1feb432807803c872b9bbbf45fe03225c19e34d`). That tag is not an applicability bound.

- `helm/clickhouse/monitoring/clickhouse_performance_grafana_dashboard.json`


### ClickHouse backup orchestrator is not installed

**Symptoms:**

- ClickHouse dashboard row CH Backup Orchestrator or series `backup_storage_last_successful_size` is empty.
- Deployment `clickhouse-backup-orchestrator` is absent.

**Root cause:**

Backup orchestrator is gated by `backupDaemon.install` (clickhouse chart default `no`; clickhouse-services default
`yes`). The backup ServiceMonitor can still render with the cluster monitor and then have no target; that later layer
is shared-discovery target-unavailable.

**Applies to:** clickhouse-operator-helm when backup orchestrator metrics are expected. Image and chart versions are
supporting facts.

**Failed layer:** `component configuration missing`

**How to check:**

1. Confirm that Deployment `clickhouse-backup-orchestrator` is absent.
2. Name `backupDaemon.install` only when effective values show it is false.

**How to fix:**

1. Set `backupDaemon.install: true` when that key is present in the evidence, then upgrade.
2. Confirm that Deployment and Service `clickhouse-backup-orchestrator` exist.

**Owner:** clickhouse-operator-helm (`Netcracker/qubership-clickhouse-operator-helm`)

**Risk:** Enabling the backup orchestrator deploys backup workloads and adds scrape load.

**Provenance:** Maintainer-verified at `Netcracker/qubership-clickhouse-operator-helm` tag `0.49.1`
(`e1feb432807803c872b9bbbf45fe03225c19e34d`). That tag is not an applicability bound.

- `helm/clickhouse/templates/clickhouse-operator/clickhouse-service-monitor.yaml`
- `helm/clickhouse/values.yaml`


### ClickHouse DR replicator monitor is absent

**Symptoms:**

- ClickHouse DR replicator metrics are expected.
- ServiceMonitor `clickhouse-replicator-service-monitor` is absent.

**Root cause:**

That ServiceMonitor renders only when `.Values.disasterRecovery` is set and `monitoring.install` is true.

**Applies to:** clickhouse-operator-helm when DR replicator metrics are expected. Image and chart versions are
supporting facts.

**Failed layer:** `component configuration missing`

**How to check:**

1. Confirm that `clickhouse-replicator-service-monitor` is absent.
2. Name disaster-recovery and monitoring keys only when effective values show DR or monitoring is off.

**How to fix:**

1. Enable disaster recovery and monitoring.install when those keys are present in the evidence, then upgrade.
2. Confirm that `clickhouse-replicator-service-monitor` selects `app: clickhouse-replicator`.

**Owner:** clickhouse-operator-helm (`Netcracker/qubership-clickhouse-operator-helm`)

**Risk:** Enabling DR replicator monitoring adds a scrape target.

**Provenance:** Maintainer-verified at `Netcracker/qubership-clickhouse-operator-helm` tag `0.49.1`
(`e1feb432807803c872b9bbbf45fe03225c19e34d`). That tag is not an applicability bound.

- `helm/clickhouse-services/templates/disaster-recovery/replicator-service-monitor.yaml`


### ClickHouse site-manager monitor is absent

**Symptoms:**

- Site-manager metrics are expected for ClickHouse DR.
- ServiceMonitor `clickhouse-sm-metrics` is absent.

**Root cause:**

That ServiceMonitor renders when `disasterRecovery.siteManager` is set. It is not gated by `monitoring.install`.

**Applies to:** clickhouse-operator-helm when ClickHouse site-manager metrics are expected. Image and chart versions
are supporting facts.

**Failed layer:** `monitor missing or misconfigured`

**How to check:**

1. Confirm that `clickhouse-sm-metrics` is absent.
2. Name `disasterRecovery.siteManager` only when effective values show site-manager is off.

**How to fix:**

1. Enable `disasterRecovery.siteManager` when that key is present in the evidence, then upgrade.
2. Confirm that `clickhouse-sm-metrics` selects `app: site-manager` and path `/metrics`.

**Owner:** clickhouse-operator-helm (`Netcracker/qubership-clickhouse-operator-helm`)

**Risk:** Enabling site-manager metrics adds a scrape target.

**Provenance:** Maintainer-verified at `Netcracker/qubership-clickhouse-operator-helm` tag `0.49.1`
(`e1feb432807803c872b9bbbf45fe03225c19e34d`). That tag is not an applicability bound.

- `helm/clickhouse-services/templates/disaster-recovery/site-manager-metrics.yaml`
