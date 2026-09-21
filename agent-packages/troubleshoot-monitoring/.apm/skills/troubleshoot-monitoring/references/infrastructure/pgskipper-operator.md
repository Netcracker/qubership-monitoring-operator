# Infrastructure component troubleshooting


## pgskipper-operator

### Collector is not installed

**Symptoms:**

- PostgreSQL collector workload `monitoring-collector` is absent while the ticket expects collector metrics or the
  `PostgreSQL Cluster Prometheus` dashboard.

**Root cause:**

When the chart does not install the collector, no `monitoring-collector` Deployment or Service exists. The documented
enablement is `metricCollector.install`. This is a missing workload, not a scrape failure. Do not treat an absent
ServiceMonitor as this case.

**Applies to:** pgskipper-operator when collector metrics are expected. Image and chart versions are supporting facts.

**Failed layer:** `component configuration missing`

**How to check:**

1. Confirm from the ticket or cluster that `monitoring-collector` is absent.
2. Name a Helm key only when effective or rendered values show `metricCollector.install` is false, or that the chart
   did not install monitoring.
3. If the expected dashboard or series is unknown, request that title or series name, then still stop while the
   workload is missing.

If those artifacts are missing, request them as Data required. Do not guess the key from a version.

**How to fix:**

1. Set `metricCollector.install: true` in the pgskipper-operator deployment values when that key is present in the
   evidence.
2. Re-render or upgrade the pgskipper-operator release using the component team's normal deployment process.
3. Confirm that Deployment and Service `monitoring-collector` exist, with labels `app: monitoring-collector` and port
   `prometheus-port`.

Do not change Monitoring Operator discovery selectors until the collector workload exists.

**Owner:** pgskipper-operator (`Netcracker/pgskipper-operator`)

**Risk:** Setting `metricCollector.install: true` and upgrading the release deploys `monitoring-collector` and starts
collector queries against the database.

**Provenance:** Maintainer-verified at `Netcracker/pgskipper-operator` tag `1.54.3`
(`a572be2516932784b53ed70fc176873d629507a6`). That tag is not an applicability bound.

- `operator/charts/patroni-services/values.yaml`
- `operator/charts/patroni-services/templates/_helpers.tpl`
- `operator/charts/patroni-services/templates/cr.yaml`
- `operator/controllers/postgresservice_controller.go`
- `docs/public/installation.md`


### Collector monitor is absent

**Symptoms:**

- PostgreSQL dashboard panels or PostgreSQL collector panels have no data.
- PostgreSQL collector workloads are installed, but `postgres-service-monitor` is absent.

**Root cause:**

The collector workload exists, so this is a missing scrape resource, not a missing component. The documented
enablement for `postgres-service-monitor` is `metricCollector.prometheusMonitoring`. Cluster-status series such as
`ma_pg_patroni_cluster_status` never appear if nothing scrapes `/metrics` on `prometheus-port`.

**Applies to:** pgskipper-operator when the PostgreSQL collector workload exists. Image and chart versions are
supporting facts.

**Failed layer:** `monitor missing or misconfigured`

**How to check:**

1. Confirm that `monitoring-collector` exists.
2. Confirm that `postgres-service-monitor` is absent.
3. Name `metricCollector.prometheusMonitoring` only when effective or rendered values show it is false, or the
   rendered monitor templates omit that ServiceMonitor.

If those artifacts are missing, request them as Data required. Do not guess the key from a version.

**How to fix:**

1. Set `metricCollector.prometheusMonitoring: true` in the pgskipper-operator deployment values when that key is
   present in the evidence.
2. Re-render or upgrade the pgskipper-operator release using the component team's normal deployment process.
3. Confirm that `postgres-service-monitor` selects `app: monitoring-collector` and port `prometheus-port`.

Do not change Monitoring Operator discovery selectors until this ServiceMonitor exists. Do not use this case when
`monitoring-collector` itself is absent.

**Owner:** pgskipper-operator (`Netcracker/pgskipper-operator`)

**Risk:** Enabling `metricCollector.prometheusMonitoring` and upgrading the pgskipper-operator release creates
`postgres-service-monitor` in the component namespace.

**Provenance:** Maintainer-verified at `Netcracker/pgskipper-operator` tag `1.54.3`
(`a572be2516932784b53ed70fc176873d629507a6`). That tag is not an applicability bound.

- `operator/charts/patroni-services/values.yaml`
- `operator/charts/patroni-services/templates/monitoring-templates/service-monitor.yaml`
- `operator/charts/patroni-services/templates/_helpers.tpl`
- `docs/public/installation.md`


### Collector profile does not provide the panel metrics

**Symptoms:**

- `postgres-service-monitor` exists and its target is healthy.
- Query or table performance dashboard panels remain empty.

**Root cause:**

The scrape target is healthy, so this is a missing or disabled series, not a scrape failure. The collector emits
query-statistics and table-statistics series such as `ma_pg_queries_stat` only when env `METRICS_PROFILE` is `dev`.
The chart value is `metricCollector.metricsProfile`. Connection panels named `ma_pg_metrics_current_connections` are
not this case.

**Applies to:** pgskipper-operator when the collector is scraped. Image and chart versions are supporting facts.

**Failed layer:** `metric series absent`

**How to check:**

1. Identify the empty panel and its required metric family. If that mapping is not known, request the panel title
   and its PromQL or series name as Data required.
2. Confirm that `postgres-service-monitor` exists and the collector target is healthy.
3. Name `metricsProfile` only when effective values or collector env `METRICS_PROFILE` show `prod` (or another
   non-`dev` profile) while the empty panel needs a `dev` series.

**How to fix:**

1. Keep `metricsProfile: prod` when the panel does not require the additional `dev`-profile metrics.
2. Set `metricsProfile: dev` only when the required panel metric is provided by that profile and the component team
   accepts the collection cost.

**Owner:** pgskipper-operator (`Netcracker/pgskipper-operator`)

**Risk:** `metricsProfile: dev` increases query and table statistics collection cost on the database.

**Provenance:** Maintainer-verified at `Netcracker/pgskipper-operator` tag `1.54.3`
(`a572be2516932784b53ed70fc176873d629507a6`). That tag is not an applicability bound.

- `operator/charts/patroni-services/values.yaml`
- `operator/pkg/deployment/monitoring.go`
- `services/monitoring-agent/collector/main.go`
- `operator/charts/patroni-services/monitoring/grafana-dashboard.json`
- `docs/public/installation.md`


### PostgreSQL dashboard CR is absent

**Symptoms:**

- GrafanaDashboard CR `postgresql-grafana-dashboard` is absent, so bundled PostgreSQL panels never load.
- For managed DB, the type-specific CR (`aws-postgresql-grafana-dashboard`, `azure-postgresql-grafana-dashboard`, or
  `cloudsql-postgresql-grafana-dashboard`) is absent.

**Root cause:**

This is a missing dashboard CR, not a scrape failure and not a missing series on an existing dashboard. On-cluster
dashboards are gated by `metricCollector.applyGrafanaDashboard`. Managed-DB dashboards are gated by
`externalDataBase.type` plus `externalDataBase.applyGrafanaDashboard`. Query Exporter's dashboard CR is a different
case.

**Applies to:** pgskipper-operator when bundled PostgreSQL Grafana panels never load. Image and chart versions are
supporting facts.

**Failed layer:** `dashboard query mismatch`

**How to check:**

1. Confirm whether the expected GrafanaDashboard CR is absent.
2. For on-cluster PostgreSQL, name `metricCollector.applyGrafanaDashboard` only when effective values show it is
   false.
3. For managed DB, name `externalDataBase.applyGrafanaDashboard` only when type (`rds`, `azure`, or `cloudsql`) and
   that flag are in the evidence.
4. If the expected dashboard name is not known, request that CR name or the dashboard title in Grafana as Data
   required.

**How to fix:**

1. For on-cluster PostgreSQL, set `metricCollector.applyGrafanaDashboard: true` when that key is present, then
   upgrade so Helm creates `postgresql-grafana-dashboard`.
2. For managed DB, set `externalDataBase.applyGrafanaDashboard: true` with `externalDataBase.type` `rds`, `azure`,
   or `cloudsql`.
3. Confirm the matching GrafanaDashboard CR exists after the upgrade.

Empty panels after the CR exists are a later metric-path layer, not this case.

**Owner:** pgskipper-operator (`Netcracker/pgskipper-operator`)

**Risk:** Creating the GrafanaDashboard CR loads bundled PostgreSQL panels into Grafana. Prerequisites include
`integreatly.org/v1alpha1` create rights for the deploy user.

**Provenance:** Maintainer-verified at `Netcracker/pgskipper-operator` tag `1.54.3`
(`a572be2516932784b53ed70fc176873d629507a6`). That tag is not an applicability bound.

- `operator/charts/patroni-services/values.yaml`
- `operator/charts/patroni-services/templates/monitoring-templates/dashboards/grafana-dashboard.yaml`
- `docs/public/installation.md`


### Query Exporter monitor is absent

**Symptoms:**

- PostgreSQL performance dashboard panels have no data.
- `query-exporter-service-monitor` is absent.

**Root cause:**

This is a missing component, not a scrape failure of an installed exporter. The documented enablement is
`queryExporter.install`. That flag also gates the Query Exporter workload, ConfigMap, Secret, and GrafanaDashboard
`query-exporter-grafana-dashboard`. The Query Exporter image repository does not own this ServiceMonitor. Installing
Query Exporter and Postgres Exporter together can fail the Helm render.

**Applies to:** pgskipper-operator when PostgreSQL performance metrics from Query Exporter are expected. Image and
chart versions are supporting facts.

**Failed layer:** `component configuration missing`

**How to check:**

1. Confirm that `query-exporter-service-monitor` is absent, and that the Query Exporter workload is absent when that
   fact is available.
2. Name `queryExporter.install` only when effective or rendered values show it is false.

If those artifacts are missing, request them as Data required. Do not guess the key from a version.

**How to fix:**

1. Set `queryExporter.install: true` in the pgskipper-operator deployment values when that key is present in the
   evidence.
2. Re-render or upgrade the pgskipper-operator release using the component team's normal deployment process.
3. Confirm that `query-exporter-service-monitor` selects `app: query-exporter` and port `web`.
4. Do not set `postgresExporter.install` at the same time if the Secret template fails the render.

Do not change Monitoring Operator selectors, dashboard JSON, or the Query Exporter image repository.

**Owner:** pgskipper-operator (`Netcracker/pgskipper-operator`)

**Risk:** Setting `queryExporter.install: true` deploys Query Exporter and creates
`query-exporter-service-monitor`. That adds a scrape target, a PostgreSQL exporter user, and query load on the
database.

**Provenance:** Maintainer-verified at `Netcracker/pgskipper-operator` tag `1.54.3`
(`a572be2516932784b53ed70fc176873d629507a6`). That tag is not an applicability bound.

- `operator/charts/patroni-services/values.yaml`
- `operator/charts/patroni-services/templates/query-exporter/query-exporter-service-monitor.yaml`
- `docs/public/installation.md`


### pgBackRest Exporter monitor is absent

**Symptoms:**

- pgBackRest exporter dashboard panels have no data.
- `backrest-exporter-monitor` is absent.

**Root cause:**

The chart renders `backrest-exporter-monitor` only when `pgBackRestExporter.install` is true. Public docs at tag
`1.54.3` do not name this key, the monitor, or a pgBackRest Prometheus exporter; the gate is chart evidence.

**Applies to:** pgskipper-operator when pgBackRest exporter metrics are expected. Image and chart versions are
supporting facts.

**Failed layer:** `component configuration missing`

**How to check:**

1. Confirm that `backrest-exporter-monitor` is absent.
2. Name `pgBackRestExporter.install` only when effective or rendered values show that key is false, or the rendered
   templates omit that ServiceMonitor.

If those artifacts are missing, request them as Data required. Do not invent the key from a version label.

**How to fix:**

1. Set `pgBackRestExporter.install: true` in the pgskipper-operator deployment values when that key is present in
   the evidence.
2. Re-render or upgrade the pgskipper-operator release using the component team's normal deployment process.
3. Confirm that `backrest-exporter-monitor` selects `app: pgbackrest-exporter` and port `brexporter`.

**Owner:** pgskipper-operator (`Netcracker/pgskipper-operator`)

**Risk:** Enabling the exporter creates `backrest-exporter-monitor` and an additional scrape target.

**Provenance:** Maintainer-verified at `Netcracker/pgskipper-operator` tag `1.54.3`
(`a572be2516932784b53ed70fc176873d629507a6`). Chart evidence only; public docs do not document this path. That tag is
not an applicability bound.

- `operator/charts/patroni-services/values.yaml`
- `operator/charts/patroni-services/templates/monitoring-templates/backrest-exporter-monitor.yaml`


### Query Exporter extensions are missing

**Symptoms:**

- Query Exporter is installed and scraped.
- Extension-backed panels stay empty.
- Docs name `pg_stat_statements`, `pgsentinel`, `pg_wait_sampling`, `pg_buffercache`, and `pg_stat_kcache`.
- On AWS RDS, `pgsentinel` and `pg_wait_sampling` are not supported.

**Root cause:**

Query Exporter is scraped, so this is a missing series, not a scrape failure. Queries that need those catalogs fail or
return nothing if the extension is missing. There is no Helm key for the extension list. Typical series mapping:
`pgsentinel` to `pg_pash_*`; `pg_wait_sampling` to `pg_wait_type_count`; `pg_buffercache` to `pg_buffercache_*`;
`pg_stat_kcache` with `pg_stat_statements` to `pg_stat_statements_query`. Do not use a panel that queries
`pg_wait_query_count` unless that series exists on `/metrics`.

**Applies to:** pgskipper-operator when Query Exporter is installed and scraped. Image and chart versions are
supporting facts.

**Failed layer:** `metric series absent`

**How to check:**

1. Confirm Query Exporter is installed and `query-exporter-service-monitor` is scraped.
2. Identify the empty panel and its required series. If that mapping is not known, request the panel title and its
   PromQL or series name as Data required.
3. Confirm whether the required extension exists on the target database. Do not treat RDS missing `pgsentinel` or
   `pg_wait_sampling` as a misconfiguration.

**How to fix:**

1. On-cluster, keep Query Exporter installed so reconcile can create the documented extensions on database
   `postgres`.
2. On managed DB, enable those extensions on the instance manually.
3. Do not expect `pgsentinel` or `pg_wait_sampling` on AWS RDS.

**Owner:** pgskipper-operator (`Netcracker/pgskipper-operator`)

**Risk:** Creating extensions on the database changes PostgreSQL catalogs and can add shared-preload requirements.
RDS cannot enable the unsupported pair.

**Provenance:** Maintainer-verified at `Netcracker/pgskipper-operator` tag `1.54.3`
(`a572be2516932784b53ed70fc176873d629507a6`). That tag is not an applicability bound.

- `docs/public/features/query-exporter.md`
- `operator/pkg/queryexporter/query_exporter.go`
- `operator/charts/patroni-services/query-exporter/query-exporter-queries.yaml`


### Query Exporter query is excluded

**Symptoms:**

- Named Query Exporter series are absent while other query-exporter series exist.

**Root cause:**

`queryExporter.excludeQueries` skips named queries. All metrics for an excluded query are omitted. Query Exporter is
scraped and other series exist, so this is a disabled series, not a scrape failure. Typical mapping:
`pg_lock_tree_query` to `pg_lock_tree_info`; `connection_by_role_with_limit_query` to
`connection_by_role_with_limit_role_cnt`. Do not use `pg_wait_query_count` as a symptom unless that series exists on
`/metrics`.

**Applies to:** pgskipper-operator when Query Exporter is installed and scraped. Image and chart versions are
supporting facts.

**Failed layer:** `metric series absent`

**How to check:**

1. Confirm Query Exporter is installed and scraped, and that some query-exporter series exist.
2. Identify the empty panel's query name and series. If that mapping is not known, request the panel title and its
   PromQL or series name as Data required.
3. Name `excludeQueries` only when that query is in effective `queryExporter.excludeQueries` or env
   `EXCLUDED_QUERIES`.

Do not recommend removing exclusions unless that query is in the list. Circuit-breaker skip is a different case.

**How to fix:**

1. Remove the empty panel's query name from `queryExporter.excludeQueries`.
2. Re-render or upgrade the pgskipper-operator release so operator env `EXCLUDED_QUERIES` no longer lists it.

**Owner:** pgskipper-operator (`Netcracker/pgskipper-operator`)

**Risk:** Removing an exclusion resumes that query against the database on every collection interval.

**Provenance:** Maintainer-verified at `Netcracker/pgskipper-operator` tag `1.54.3`
(`a572be2516932784b53ed70fc176873d629507a6`). That tag is not an applicability bound.

- `docs/public/features/query-exporter.md`
- `operator/pkg/queryexporter/query_exporter.go`
- `operator/charts/patroni-services/query-exporter/query-exporter-queries.yaml`


### Query Exporter circuit breaker skipped a query

**Symptoms:**

- Series for one Query Exporter query stop updating after repeated timeouts, while other queries continue.

**Root cause:**

After `queryExporter.queryTimeout` is exceeded `queryExporter.maxFailedTimeouts` times, Query Exporter skips that
query. Other queries continue and the ServiceMonitor still scrapes, so this is a disabled series, not a scrape
failure. The skip list is Query Exporter process memory. A Query Exporter pod restart clears it. Docs that say restart
the postgres porter pod are wrong.

**Applies to:** pgskipper-operator when Query Exporter is installed and scraped. Image and chart versions are
supporting facts.

**Failed layer:** `metric series absent`

**How to check:**

1. Confirm Query Exporter is installed and scraped, and that other query-exporter series still update.
2. Identify the skipped query from logs such as `Query %s is skipped as long-running`. If the query or panel mapping
   is not known, request that log line or the panel's series name as Data required.
3. Name timeout settings only when effective `queryExporter.queryTimeout` and `queryExporter.maxFailedTimeouts` (or
   the matching exporter env) are in the evidence.

Do not fold this into `collection-timeout`. ServiceMonitor `scrapeTimeout` is a different case.

**How to fix:**

1. Restart the Query Exporter pod to clear the in-memory skip list.
2. If the query is slow, raise `queryExporter.queryTimeout` or `queryExporter.maxFailedTimeouts` when those keys are
   present, then upgrade.
3. Do not restart a postgres porter pod.

**Owner:** pgskipper-operator (`Netcracker/pgskipper-operator`)

**Risk:** Raising the timeout or failure count lets a slow query run longer and more often. Restarting Query Exporter
clears the skip list and briefly gaps scrapes.

**Provenance:** Maintainer-verified at `Netcracker/pgskipper-operator` tag `1.54.3`
(`a572be2516932784b53ed70fc176873d629507a6`). That tag is not an applicability bound.

- `docs/public/features/query-exporter.md`
- `operator/pkg/queryexporter/query_exporter.go`
- `services/query-exporter/main.go`


### Custom Query Exporter queries are not loaded

**Symptoms:**

- Custom Query Exporter series never appear.
- ConfigMaps lack label `query-exporter: custom-queries`, custom query loading is off, or the operator cannot list
  ConfigMaps in the configured namespaces.

**Root cause:**

Custom query loading is off or the ConfigMaps are unlabeled, so the merged query set never includes them. This is
missing component configuration, not a scrape failure of bundled queries. The documented enablement is
`queryExporter.customQueries.enabled`. Default labels include `query-exporter: custom-queries`. Custom series are not
in the bundled dashboard unless the custom metric names match existing panels.

**Applies to:** pgskipper-operator when Query Exporter is installed. Image and chart versions are supporting facts.

**Failed layer:** `component configuration missing`

**How to check:**

1. Confirm Query Exporter is installed.
2. Name `queryExporter.customQueries.enabled` only when effective values show it is false.
3. Confirm ConfigMaps in the configured namespaces carry labels matching `customQueries.labels`.
4. If the expected custom series name is not known, request that metric name or ConfigMap as Data required.

**How to fix:**

1. Set `queryExporter.customQueries.enabled: true` when that key is present in the evidence.
2. Put ConfigMaps in `customQueries.namespacesList` with labels matching `customQueries.labels` (default includes
   `query-exporter: custom-queries`).
3. Grant ConfigMap get/watch RBAC, then upgrade.

**Owner:** pgskipper-operator (`Netcracker/pgskipper-operator`)

**Risk:** Enabling custom queries restarts Query Exporter pods when watched ConfigMaps change, which briefly gaps
scrapes. Merged queries add load on the database.

**Provenance:** Maintainer-verified at `Netcracker/pgskipper-operator` tag `1.54.3`
(`a572be2516932784b53ed70fc176873d629507a6`). That tag is not an applicability bound.

- `docs/public/features/query-exporter.md`
- `operator/charts/patroni-services/values.yaml`
- `operator/pkg/queryexporter/query_exporter.go`


### Query Exporter access failed

**Symptoms:**

- Query Exporter is installed and scraped, but query execution fails because the exporter user or grants are wrong.
- Series stay empty or self-monitor `query_status{status="failed"}` increments.

**Root cause:**

Query Exporter is scraped, so this is failed query execution, not a missing ServiceMonitor. The exporter user needs
`pg_read_all_data` and `pg_monitor`. Do not read Secret values during diagnosis.

**Applies to:** pgskipper-operator when Query Exporter is installed and scraped. Image and chart versions are
supporting facts.

**Failed layer:** `metric series absent`

**How to check:**

1. Confirm Query Exporter is installed and scraped.
2. Prefer `query_status{status="failed"}` when self-monitor is enabled. If that series is absent, request non-secret
   evidence that the exporter user exists and holds `pg_read_all_data` and `pg_monitor`.
3. If the empty panel's series is not known, request the panel title and its PromQL or series name, then still
   confirm grants before naming a dashboard mismatch.

Do not use this case when the exporter is not installed.

**How to fix:**

1. Confirm without reading Secrets that the exporter user exists and holds `pg_read_all_data` and `pg_monitor`.
2. Reconcile Query Exporter so those roles are granted.

**Owner:** pgskipper-operator (`Netcracker/pgskipper-operator`)

**Risk:** Granting `pg_monitor` and `pg_read_all_data` widens that user's catalog access. Do not copy Secret values
into the diagnosis report.

**Provenance:** Maintainer-verified at `Netcracker/pgskipper-operator` tag `1.54.3`
(`a572be2516932784b53ed70fc176873d629507a6`). That tag is not an applicability bound.

- `docs/public/features/query-exporter.md`
- `operator/pkg/queryexporter/query_exporter.go`
- `operator/charts/patroni-services/templates/query-exporter/query-exporter-user-credentials.yaml`


### Query Exporter self-monitoring is disabled

**Symptoms:**

- `query_status` and `query_latency` are absent while other Query Exporter series exist.

**Root cause:**

Other Query Exporter series exist and the exporter is scraped, so this is a disabled series, not a scrape failure.
Env `QUERY_EXPORTER_DISABLE_SELF_MONITOR` comes from `queryExporter.selfMonitorDisabled`. When that flag is true,
self-monitor does not register `query_status` or `query_latency`. Default `false` means self-monitor is enabled. The
feature-page sentence that says set the parameter to false to disable is wrong. The bundled dashboard does not panel
those series; the symptom is missing series on `/metrics`.

**Applies to:** pgskipper-operator when Query Exporter is installed and scraped. Image and chart versions are
supporting facts.

**Failed layer:** `metric series absent`

**How to check:**

1. Confirm Query Exporter is installed and scraped, and that other query-exporter series exist.
2. Confirm `query_status` and `query_latency` are absent on the exporter `/metrics`.
3. Name `queryExporter.selfMonitorDisabled` only when effective values or env `QUERY_EXPORTER_DISABLE_SELF_MONITOR`
   show `true`. Use that verified rule; do not follow the feature-page disable sentence.

**How to fix:**

1. Set `queryExporter.selfMonitorDisabled: false` when that key is present, then upgrade so
   `QUERY_EXPORTER_DISABLE_SELF_MONITOR` is `false`.
2. Set the flag to `true` only when disabling is intended.

**Owner:** pgskipper-operator (`Netcracker/pgskipper-operator`)

**Risk:** Re-enabling self-monitor adds `query_status` and `query_latency` on the exporter `/metrics` endpoint.

**Provenance:** Maintainer-verified at `Netcracker/pgskipper-operator` tag `1.54.3`
(`a572be2516932784b53ed70fc176873d629507a6`). That tag is not an applicability bound.

- `operator/charts/patroni-services/values.yaml`
- `operator/pkg/queryexporter/query_exporter.go`
- `services/query-exporter/selfmonitor/self-monitor.go`
- `docs/public/installation.md`
- `docs/public/features/query-exporter.md` (feature-page disable sentence is inconsistent with the code)


### Collection scrape timeout

**Symptoms:**

- ServiceMonitor exists and is discovered, but the scrape ends before the exporter finishes.
- The target is down or partial.

**Root cause:**

This case is the chart `scrapeTimeout` on the pgskipper ServiceMonitor endpoint. It is not the shared-discovery
target-unavailable diagnosis (Service selector, endpoint port, or ready endpoints). Collector Go `scrapeTimeout` is
the collection interval, not this field. Query Exporter `queryTimeout` is the circuit-breaker case. A documented
collection interval delay is not an empty-metrics diagnosis.

**Applies to:** pgskipper-operator when a pgskipper ServiceMonitor is discovered. Image and chart versions are
supporting facts.

**Failed layer:** `target unavailable`

**How to check:**

1. Confirm the expected ServiceMonitor exists and is in discovery scope.
2. Confirm the scrape timeout on that monitor from the live ServiceMonitor YAML or effective values
   (`metricCollector.scrapeTimeout` or `queryExporter.scrapeTimeout`).
3. Confirm the scrape ends before gathering finishes (target down or partial), not that endpoints are missing.
4. If only a collection interval is documented, request time-series evidence as Data required.

**How to fix:**

1. Raise the ServiceMonitor `scrapeTimeout` that the evidence named, then upgrade.
2. Do not use `collectionInterval` as the fix.
3. Do not change Monitoring Operator selectors for a scrape that already discovers the monitor.

**Owner:** pgskipper-operator (`Netcracker/pgskipper-operator`)

**Risk:** Raising scrape timeout keeps scrape workers busy longer on slow targets.

**Provenance:** Maintainer-verified at `Netcracker/pgskipper-operator` tag `1.54.3`
(`a572be2516932784b53ed70fc176873d629507a6`). That tag is not an applicability bound.

- `operator/charts/patroni-services/values.yaml`
- `operator/charts/patroni-services/templates/monitoring-templates/service-monitor.yaml`
- `operator/charts/patroni-services/templates/query-exporter/query-exporter-service-monitor.yaml`
- `docs/public/installation.md`
