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


## mongodb-operator

### Prometheus exporter is not installed

**Symptoms:**

- MongoDB 7 Prometheus or MongoDB 8 Prometheus dashboard panels have no data.
- Deployment and Service `mongodb-prometheus-exporter` are absent.

**Root cause:**

The mongodb-services chart does not install the exporter workload, so no scrape target exists. The documented
enablement is `prometheusExporter.install`, with helper fallback `MONITORING_ENABLED` (default true when both are
unset). This is a missing workload, not a scrape failure.

**Applies to:** mongodb-operator when MongoDB exporter metrics are expected. Image and chart versions are supporting
facts.

**Failed layer:** `component configuration missing`

**How to check:**

1. Confirm that `mongodb-prometheus-exporter` is absent.
2. Name `prometheusExporter.install` or `MONITORING_ENABLED` only when effective or rendered values show monitoring is
   off.
3. If the expected dashboard or series is unknown, request that title or series name, then still stop while the
   workload is missing.

If those artifacts are missing, request them as Data required. Do not guess the key from a version.

**How to fix:**

1. Set `prometheusExporter.install: true` in the mongodb-services deployment values when that key is present in the
   evidence.
2. Re-render or upgrade the mongodb-services release using the component team's normal deployment process.
3. Confirm that Deployment and Service `mongodb-prometheus-exporter` exist, with labels
   `microservice: mongodb-prometheus-exporter`.

Do not change Monitoring Operator discovery selectors until the exporter workload exists.

**Owner:** mongodb-operator (`Netcracker/qubership-mongodb-operator`)

**Risk:** Enabling the exporter deploys `mongodb-prometheus-exporter` and starts metric queries against MongoDB.

**Provenance:** Maintainer-verified at `Netcracker/qubership-mongodb-operator` tag `2.16.9`
(`3c4aac847e3e87e88638a796148f1851492309f8`). That tag is not an applicability bound.

- `services/service/charts/helm/mongodb-services/templates/cr_and_operator.yaml`
- `services/service/charts/helm/mongodb-services/templates/_helper.tpl`
- `services/service/pkg/prometheus_exporter/templates.go`
- `docs/public/installation_guide.md`


### Backup exporter monitor is absent

**Symptoms:**

- MongoDB backup dashboard panels such as Last Backup Status or Last Backup Size have no data.
- ServiceMonitor `mongo-backup-exporter-service-monitor` is absent.

**Root cause:**

The backup scrape resource is rendered only when backup is installed and monitoring is on. The documented gate is
`backup.install`. This is a missing backup exporter scrape, not the main MongoDB exporter.

**Applies to:** mongodb-operator when backup exporter metrics are expected. Image and chart versions are supporting
facts.

**Failed layer:** `component configuration missing`

**How to check:**

1. Confirm that `mongo-backup-exporter-service-monitor` is absent.
2. Name `backup.install` only when effective or rendered values show it is false, or the rendered templates omit that
   ServiceMonitor.

If those artifacts are missing, request them as Data required. Do not invent the key from a version label.

**How to fix:**

1. Set `backup.install: true` in the mongodb-services deployment values when that key is present in the evidence.
2. Re-render or upgrade the mongodb-services release using the component team's normal deployment process.
3. Confirm that `mongo-backup-exporter-service-monitor` selects `microservice: mongodb-backup-daemon` and path
   `/health/prometheus`.

**Owner:** mongodb-operator (`Netcracker/qubership-mongodb-operator`)

**Risk:** Enabling backup creates `mongodb-backup-daemon` and an additional scrape target.

**Provenance:** Maintainer-verified at `Netcracker/qubership-mongodb-operator` tag `2.16.9`
(`3c4aac847e3e87e88638a796148f1851492309f8`). That tag is not an applicability bound.

- `services/service/charts/helm/mongodb-services/templates/service_monitor.yaml`
- `services/service/charts/helm/mongodb-services/values.yaml`
- `docs/public/grafana_dashboard.md`


### DBaaS adapter monitor is absent

**Symptoms:**

- MongoDB DBaaS Adapter dashboard row has no data.
- ServiceMonitor `mongo-dbaas-exporter-service-monitor` is absent.

**Root cause:**

The DBaaS adapter scrape resource is rendered only when the adapter is installed and monitoring is on. The documented
gates are `dbaas.install` and `DBAAS_ENABLED`.

**Applies to:** mongodb-operator when DBaaS adapter metrics are expected. Image and chart versions are supporting facts.

**Failed layer:** `component configuration missing`

**How to check:**

1. Confirm that `mongo-dbaas-exporter-service-monitor` is absent.
2. Name `dbaas.install` or `DBAAS_ENABLED` only when effective or rendered values show the adapter is off.

If those artifacts are missing, request them as Data required.

**How to fix:**

1. Set `dbaas.install: true` in the mongodb-services deployment values when that key is present in the evidence.
2. Re-render or upgrade the mongodb-services release using the component team's normal deployment process.
3. Confirm that `mongo-dbaas-exporter-service-monitor` selects the DBaaS adapter and path `/metrics`.

**Owner:** mongodb-operator (`Netcracker/qubership-mongodb-operator`)

**Risk:** Enabling the DBaaS adapter adds a scrape target for `dbaas-mongo-adapter`.

**Provenance:** Maintainer-verified at `Netcracker/qubership-mongodb-operator` tag `2.16.9`
(`3c4aac847e3e87e88638a796148f1851492309f8`). That tag is not an applicability bound.

- `services/service/charts/helm/mongodb-services/templates/service_monitor.yaml`
- `services/service/charts/helm/mongodb-services/templates/_helper.tpl`
- `docs/public/grafana_dashboard.md`


### MongoDB collection scrape timeout

**Symptoms:**

- ServiceMonitor `mongodb-prometheus-exporter-service-monitor` exists and is discovered, but the scrape ends before the
  exporter finishes.
- The target is down or partial.

**Root cause:**

This case is the chart `prometheusExporter.collectionScrapeTimeout` on the MongoDB exporter ServiceMonitor endpoint. It
is not the shared-discovery target-unavailable diagnosis (Service selector, endpoint port, or ready endpoints).

**Applies to:** mongodb-operator when a MongoDB exporter ServiceMonitor is discovered. Image and chart versions are
supporting facts.

**Failed layer:** `target unavailable`

**How to check:**

1. Confirm the expected ServiceMonitor exists and is in discovery scope.
2. Confirm the scrape timeout on that monitor from the live ServiceMonitor YAML or effective values
   (`prometheusExporter.collectionScrapeTimeout`).
3. Confirm the scrape ends before gathering finishes (target down or partial), not that endpoints are missing.

**How to fix:**

1. Raise the ServiceMonitor `scrapeTimeout` that the evidence named, then upgrade.
2. Do not use `prometheusExporter.collectionInterval` as the fix.
3. Do not change Monitoring Operator selectors for a scrape that already discovers the monitor.

**Owner:** mongodb-operator (`Netcracker/qubership-mongodb-operator`)

**Risk:** Raising scrape timeout keeps scrape workers busy longer on slow targets.

**Provenance:** Maintainer-verified at `Netcracker/qubership-mongodb-operator` tag `2.16.9`
(`3c4aac847e3e87e88638a796148f1851492309f8`). That tag is not an applicability bound.

- `services/service/charts/helm/mongodb-services/templates/service_monitor.yaml`
- `services/service/charts/helm/mongodb-services/values.yaml`
- `docs/public/installation_guide.md`


### Exporter cannot authenticate to MongoDB

**Symptoms:**

- MongoDB exporter is installed and scraped, but MongoDB series stay empty.
- The monitoring user is missing or lacks the documented role.

**Root cause:**

The exporter is scraped, so this is failed query execution, not a missing ServiceMonitor. Docs name
`prometheusExporter.monitoringUser` and `monitoringUserRole` (default user `monitoring`, role `readWrite` on `test`
plus `clusterMonitor`). Do not read Secret values during diagnosis.

**Applies to:** mongodb-operator when the MongoDB exporter is installed and scraped. Image and chart versions are
supporting facts.

**Failed layer:** `metric series absent`

**How to check:**

1. Confirm that `mongodb-prometheus-exporter` is installed and scraped.
2. Request non-secret evidence that the monitoring user exists and holds the documented role. Do not read Secret
   values.

If those artifacts are missing, request them as Data required.

**How to fix:**

1. Confirm without reading Secrets that the monitoring user exists and holds the documented role.
2. Reconcile mongodb-services so that user is created with `prometheusExporter.monitoringUserRole` when that key is
   present.

**Owner:** mongodb-operator (`Netcracker/qubership-mongodb-operator`)

**Risk:** Granting clusterMonitor and readWrite widens that user's catalog access. Do not copy Secret values into the
diagnosis report.

**Provenance:** Maintainer-verified at `Netcracker/qubership-mongodb-operator` tag `2.16.9`
(`3c4aac847e3e87e88638a796148f1851492309f8`). That tag is not an applicability bound.

- `services/service/charts/helm/mongodb-services/templates/mongo_monitoring_user.yaml`
- `docs/public/configuration.md`
- `docs/public/troubleshooting.md`
- `docs/public/installation_guide.md`


### MongoDB ServiceMonitor basicAuth mismatch

**Symptoms:**

- ServiceMonitor `mongodb-prometheus-exporter-service-monitor` exists.
- Scrape of the exporter fails authentication.

**Root cause:**

The ServiceMonitor `basicAuth` references Secret `mongodb-prom-exporter-credentials.v1` (values
`prometheusExporter.exporterUser`, `exporterPassword`, `prometheusExporterSecretName`). A missing or mismatched
exporter credential is a scrape-auth failure, not a namespace selector miss.

**Applies to:** mongodb-operator when the MongoDB exporter ServiceMonitor exists. Image and chart versions are
supporting facts.

**Failed layer:** `monitor missing or misconfigured`

**How to check:**

1. Confirm that `mongodb-prometheus-exporter-service-monitor` exists and scrape fails with an auth error.
2. Name exporter credential keys only when non-secret evidence shows `prometheusExporter.exporterUser` or the Secret
   name, never Secret values.

**How to fix:**

1. Align `prometheusExporter.exporterUser` and the ServiceMonitor basicAuth Secret with the exporter credentials
   without copying Secret values into the report.
2. Do not change Monitoring Operator namespace selectors for an auth failure on an existing monitor.

**Owner:** mongodb-operator (`Netcracker/qubership-mongodb-operator`)

**Risk:** Changing exporter basicAuth credentials can briefly fail scrapes until both the exporter and the
ServiceMonitor use the same Secret.

**Provenance:** Maintainer-verified at `Netcracker/qubership-mongodb-operator` tag `2.16.9`
(`3c4aac847e3e87e88638a796148f1851492309f8`). That tag is not an applicability bound.

- `services/service/charts/helm/mongodb-services/templates/service_monitor.yaml`
- `services/service/charts/helm/mongodb-services/templates/secrets/prometheus-exporter-secret.yaml`


## cassandra-operator

### Cassandra monitoring is not enabled

**Symptoms:**

- Cassandra Monitoring dashboard panels have no data.
- ServiceMonitors `cassandra-exporter-service-monitor` and GrafanaDashboard CR `cassandra-grafana-dashboard` are
  absent.

**Root cause:**

The cassandra-services chart does not render Prometheus scrape or dashboard objects when monitoring is off. The
documented enablement is `monitoringAgent.install`, with helper fallback `MONITORING_ENABLED` (chart default true).
This is missing component monitoring configuration, not a scrape failure of an existing monitor.

**Applies to:** cassandra-operator when Cassandra metrics are expected. Image and chart versions are supporting facts.

**Failed layer:** `component configuration missing`

**How to check:**

1. Confirm that `cassandra-exporter-service-monitor` and `cassandra-grafana-dashboard` are absent.
2. Name `monitoringAgent.install` or `MONITORING_ENABLED` only when effective or rendered values show monitoring is
   off.

If those artifacts are missing, request them as Data required. Do not guess the key from a version.

**How to fix:**

1. Set `monitoringAgent.install: true` in the cassandra-services deployment values when that key is present in the
   evidence.
2. Re-render or upgrade the cassandra-services release using the component team's normal deployment process.
3. Confirm that `cassandra-exporter-service-monitor` and `cassandra-grafana-dashboard` exist.

Do not change Monitoring Operator discovery selectors until the scrape resource exists.

**Owner:** cassandra-operator (`Netcracker/qubership-cassandra-operator`)

**Risk:** Enabling monitoring creates Cassandra ServiceMonitors and starts scrape of cluster metrics.

**Provenance:** Maintainer-verified at `Netcracker/qubership-cassandra-operator` tag `2.16.11`
(`773b2e3954b5d6fbdc3c6d6d1136f367ca6f2d3e`). That tag is not an applicability bound.

- `services/service/charts/helm/cassandra-services/templates/service_monitor.yaml`
- `services/service/charts/helm/cassandra-services/templates/grafana_dashboard.yaml`
- `services/service/charts/helm/cassandra-services/values.yaml`


### Cassandra metric collector is not Prometheus

**Symptoms:**

- `monitoringAgent.install` is true, but `cassandra-exporter-service-monitor` and `cassandra-grafana-dashboard` are
  still absent.

**Root cause:**

ServiceMonitor and GrafanaDashboard templates render only when `monitoringAgent.metricCollector` is `prometheus`.
Another collector type leaves monitoring installed without Prometheus scrape resources.

**Applies to:** cassandra-operator when monitoring is installed. Image and chart versions are supporting facts.

**Failed layer:** `monitor missing or misconfigured`

**How to check:**

1. Confirm that monitoring is installed and the Prometheus ServiceMonitor is absent.
2. Name `monitoringAgent.metricCollector` only when effective values show a value other than `prometheus`.

**How to fix:**

1. Set `monitoringAgent.metricCollector: prometheus` when that key is present in the evidence, then upgrade.
2. Do not treat this as a missing `monitoringAgent.install` flag.

**Owner:** cassandra-operator (`Netcracker/qubership-cassandra-operator`)

**Risk:** Switching the collector to Prometheus creates scrape targets and dashboard CRs.

**Provenance:** Maintainer-verified at `Netcracker/qubership-cassandra-operator` tag `2.16.11`
(`773b2e3954b5d6fbdc3c6d6d1136f367ca6f2d3e`). That tag is not an applicability bound.

- `services/service/charts/helm/cassandra-services/templates/service_monitor.yaml`
- `services/service/charts/helm/cassandra-services/templates/grafana_dashboard.yaml`


### Cassandra metrics ServiceMonitor is absent

**Symptoms:**

- Cassandra Overview or pod panels have no data.
- ServiceMonitor `cassandra-exporter-service-monitor` is absent while monitoring is on and the collector is Prometheus.

**Root cause:**

That ServiceMonitor is nested under `cassandra.install` as well as monitoring. Cluster metrics come from Service
`cassandra-metrics` created by cassandra-operator. This is a missing scrape resource for the cluster metrics Service,
not a missing monitoring-agent image.

**Applies to:** cassandra-operator when cluster metrics are expected. Image and chart versions are supporting facts.

**Failed layer:** `monitor missing or misconfigured`

**How to check:**

1. Confirm that monitoring is on and `cassandra-exporter-service-monitor` is absent.
2. Name `cassandra.install` only when effective or rendered values show it is false.

**How to fix:**

1. Set `cassandra.install: true` when that key is present in the evidence, then upgrade.
2. Confirm that Service `cassandra-metrics` exists on port `metrics` and that `cassandra-exporter-service-monitor`
   selects it.

**Owner:** cassandra-operator (`Netcracker/qubership-cassandra-operator`)

**Risk:** Enabling the Cassandra cluster creates the metrics Service and an additional scrape target.

**Provenance:** Maintainer-verified at `Netcracker/qubership-cassandra-operator` tag `2.16.11`
(`773b2e3954b5d6fbdc3c6d6d1136f367ca6f2d3e`). That tag is not an applicability bound.

- `services/service/charts/helm/cassandra-services/templates/service_monitor.yaml`
- `operator/pkg/impl/cassandra/service_step.go`
- `operator/pkg/impl/utils/const.go`


### Cassandra backup exporter monitor is absent

**Symptoms:**

- Cassandra dashboard row Backup Daemon has no data.
- ServiceMonitor `cassandra-backup-exporter-service-monitor` is absent.

**Root cause:**

The backup scrape resource is rendered only when `backupDaemon.install` is true and monitoring is Prometheus. This is
a missing backup scrape, not the cluster metrics monitor.

**Applies to:** cassandra-operator when backup exporter metrics are expected. Image and chart versions are supporting
facts.

**Failed layer:** `component configuration missing`

**How to check:**

1. Confirm that `cassandra-backup-exporter-service-monitor` is absent.
2. Name `backupDaemon.install` only when effective or rendered values show it is false.

**How to fix:**

1. Set `backupDaemon.install: true` when that key is present in the evidence, then upgrade.
2. Confirm that `cassandra-backup-exporter-service-monitor` selects `cassandra-backup-daemon` and path
   `/health/prometheus`.

**Owner:** cassandra-operator (`Netcracker/qubership-cassandra-operator`)

**Risk:** Enabling the backup daemon adds a scrape target.

**Provenance:** Maintainer-verified at `Netcracker/qubership-cassandra-operator` tag `2.16.11`
(`773b2e3954b5d6fbdc3c6d6d1136f367ca6f2d3e`). That tag is not an applicability bound.

- `services/service/charts/helm/cassandra-services/templates/service_monitor.yaml`
- `services/service/charts/helm/cassandra-services/values.yaml`


### Cassandra DBaaS adapter monitor is absent

**Symptoms:**

- Cassandra dashboard row DBaaS Adapter has no data.
- ServiceMonitor `cassandra-dbaas-exporter-service-monitor` is absent.

**Root cause:**

The DBaaS adapter scrape resource is rendered only when `dbaas.install` or `DBAAS_ENABLED` is true and monitoring is
Prometheus.

**Applies to:** cassandra-operator when DBaaS adapter metrics are expected. Image and chart versions are supporting
facts.

**Failed layer:** `component configuration missing`

**How to check:**

1. Confirm that `cassandra-dbaas-exporter-service-monitor` is absent.
2. Name `dbaas.install` or `DBAAS_ENABLED` only when effective or rendered values show the adapter is off.

**How to fix:**

1. Set `dbaas.install: true` when that key is present in the evidence, then upgrade.
2. Confirm that `cassandra-dbaas-exporter-service-monitor` selects `dbaas-cassandra-adapter` and path `/metrics`.

**Owner:** cassandra-operator (`Netcracker/qubership-cassandra-operator`)

**Risk:** Enabling the DBaaS adapter adds a scrape target.

**Provenance:** Maintainer-verified at `Netcracker/qubership-cassandra-operator` tag `2.16.11`
(`773b2e3954b5d6fbdc3c6d6d1136f367ca6f2d3e`). That tag is not an applicability bound.

- `services/service/charts/helm/cassandra-services/templates/service_monitor.yaml`


### Cassandra table debug series dropped by relabeling

**Symptoms:**

- ServiceMonitor `cassandra-exporter-service-monitor` is scraped.
- Table Metrics debug panels that use `cassandra_table_live_rows_scanned` or `cassandra_table_sstables_per_read` stay
  empty.

**Root cause:**

The ServiceMonitor metric relabelings drop those table debug series. Other Cassandra series can exist. This is a
disabled series, not a scrape failure.

**Applies to:** cassandra-operator when the cluster metrics monitor is scraped. Image and chart versions are supporting
facts.

**Failed layer:** `metric series absent`

**How to check:**

1. Confirm that `cassandra-exporter-service-monitor` is scraped and other Cassandra series exist.
2. Confirm the empty panel series name. If that mapping is not known, request the panel title and its PromQL as Data
   required.
3. Name relabelings only when the live ServiceMonitor YAML drops those series.

**How to fix:**

1. Remove the drop relabeling for the required series from `monitoringAgent.prometheus.metricRelabelings.cassandra`
   when that key is present, then upgrade.
2. Do not enable `monitoringAgent.install` as the fix for a dropped series on a healthy scrape.

**Owner:** cassandra-operator (`Netcracker/qubership-cassandra-operator`)

**Risk:** Keeping table debug series increases metric cardinality on Cassandra scrapes.

**Provenance:** Maintainer-verified at `Netcracker/qubership-cassandra-operator` tag `2.16.11`
(`773b2e3954b5d6fbdc3c6d6d1136f367ca6f2d3e`). That tag is not an applicability bound.

- `services/service/charts/helm/cassandra-services/templates/service_monitor.yaml`


## redis

### Redis monitoring agent is not installed

**Symptoms:**

- Redis Monitoring dashboard panels have no data.
- Deployment and Service `redis-monitoring-agent` are absent.

**Root cause:**

The redis-operator chart does not install the monitoring agent, so no scrape target exists. The documented enablement
is `monitoringAgent.install`, with helper fallback `MONITORING_ENABLED`. Chart values default `true`; the installation
guide table that lists default `false` is not the chart default. This is a missing workload, not a scrape failure.

**Applies to:** redis when Redis exporter metrics are expected. Image and chart versions are supporting facts.

**Failed layer:** `component configuration missing`

**How to check:**

1. Confirm that `redis-monitoring-agent` is absent.
2. Name `monitoringAgent.install` or `MONITORING_ENABLED` only when effective or rendered values show monitoring is
   off.

If those artifacts are missing, request them as Data required. Prefer the live values over the installation-guide
default table.

**How to fix:**

1. Set `monitoringAgent.install: true` in the redis-operator deployment values when that key is present in the
   evidence.
2. Re-render or upgrade the redis-operator release using the component team's normal deployment process.
3. Confirm that Deployment and Service `redis-monitoring-agent` exist, with labels `app: redis-monitoring-agent`.

Do not change Monitoring Operator discovery selectors until the agent exists.

**Owner:** redis (`Netcracker/qubership-redis`)

**Risk:** Enabling the agent deploys `redis-monitoring-agent` and starts collection against Redis.

**Provenance:** Maintainer-verified at `Netcracker/qubership-redis` tag `4.2.9`
(`790aeeee8d0ca5237a7bc9748935392a29d9c5f1`). That tag is not an applicability bound.

- `redis-operator/charts/helm/redis-operator/values.yaml`
- `redis-operator/charts/helm/redis-operator/templates/cr_and_operator.yaml`
- `redis-operator/api/v2/impl/monitoring/templates.go`
- `redis-operator/docs/public/installation_guide.md`


### Redis metric collector is not Prometheus

**Symptoms:**

- `redis-monitoring-agent` may exist, but ServiceMonitor `redis-prometheus-exporter-service-monitor` and
  GrafanaDashboard CR `redis-grafana-dashboard` are absent.

**Root cause:**

Prometheus ServiceMonitor and dashboard templates render only when `monitoringAgent.metricCollector` is `prometheus`.
Another collector type writes Influx output instead of a Prometheus scrape resource.

**Applies to:** redis when the monitoring agent is installed. Image and chart versions are supporting facts.

**Failed layer:** `monitor missing or misconfigured`

**How to check:**

1. Confirm that the agent exists or monitoring is installed, and the Prometheus ServiceMonitor is absent.
2. Name `monitoringAgent.metricCollector` only when effective values show a value other than `prometheus`.

**How to fix:**

1. Set `monitoringAgent.metricCollector: prometheus` when that key is present in the evidence, then upgrade.
2. Confirm that `redis-prometheus-exporter-service-monitor` selects `app: redis-monitoring-agent`.

**Owner:** redis (`Netcracker/qubership-redis`)

**Risk:** Switching the collector to Prometheus creates a scrape target and dashboard CR.

**Provenance:** Maintainer-verified at `Netcracker/qubership-redis` tag `4.2.9`
(`790aeeee8d0ca5237a7bc9748935392a29d9c5f1`). That tag is not an applicability bound.

- `redis-operator/charts/helm/redis-operator/templates/service_monitor.yaml`
- `redis-operator/charts/helm/redis-operator/templates/grafana_dashboard.yaml`
- `redis-operator/charts/helm/redis-operator/templates/monitoring_configmap.yaml`


### Redis agent cannot collect from Redis

**Symptoms:**

- ServiceMonitor `redis-prometheus-exporter-service-monitor` is scraped and the agent target is healthy.
- Redis panels stay empty; docs describe the agent unable to collect (Redis Node Down).

**Root cause:**

The scrape target is the monitoring agent, not the Redis pod. Empty Redis series with a healthy agent scrape is Redis
AUTH, TLS, or network failure into Redis. Telegraf `inputs.redis` uses `redis.password` / `redis.secretName` and TLS
`redis.tls.enabled`. Do not read Secret values. If the agent Service has no ready endpoints, that is shared-discovery
target-unavailable.

**Applies to:** redis when the monitoring agent is scraped. Image and chart versions are supporting facts.

**Failed layer:** `metric series absent`

**How to check:**

1. Confirm that `redis-prometheus-exporter-service-monitor` is scraped and the agent target is healthy.
2. Request non-secret evidence that Redis AUTH or TLS matches the agent config. Do not read Secret values.

**How to fix:**

1. Align Redis password or TLS settings with the agent configuration without copying Secret values into the report.
2. Do not change Monitoring Operator selectors for a healthy agent scrape.

**Owner:** redis (`Netcracker/qubership-redis`)

**Risk:** Changing Redis AUTH or TLS can reconnect clients. Do not copy Secret values into the diagnosis report.

**Provenance:** Maintainer-verified at `Netcracker/qubership-redis` tag `4.2.9`
(`790aeeee8d0ca5237a7bc9748935392a29d9c5f1`). That tag is not an applicability bound.

- `redis-operator/charts/helm/redis-operator/templates/monitoring_configmap.yaml`
- `redis-operator/docs/public/troubleshooting.md`


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

## kafka

### Kafka Monitoring is not installed

**Symptoms:**

- Kafka Monitoring dashboard panels have no data, or series `kafka_cluster_status` is absent.
- Deployment and Service `kafka-monitoring` are absent.

**Root cause:**

The kafka-service chart does not install the Telegraf monitoring workload. The documented enablement is
`monitoring.install` (default true). Helper `MONITORING_ENABLED` applies only when `global.cloudIntegrationEnabled` is
true. This is a missing workload, not a scrape failure.

**Applies to:** kafka when Kafka cluster-state metrics are expected. Image and chart versions are supporting facts.

**Failed layer:** `component configuration missing`

**How to check:**

1. Confirm that Deployment `kafka-monitoring` is absent.
2. Name `monitoring.install` or `MONITORING_ENABLED` only when effective or rendered values show monitoring is off.

If those artifacts are missing, request them as Data required. Do not guess the key from a version.

**How to fix:**

1. Set `monitoring.install: true` in the kafka-service deployment values when that key is present in the evidence.
2. Re-render or upgrade the kafka-service release using the component team's normal deployment process.
3. Confirm that Deployment and Service `kafka-monitoring` exist, with labels `component: kafka-monitoring`.

Do not change Monitoring Operator discovery selectors until the monitoring workload exists.

**Owner:** kafka (`Netcracker/qubership-kafka`)

**Risk:** Enabling monitoring deploys `kafka-monitoring` and starts Telegraf collection against Kafka.

**Provenance:** Maintainer-verified at `Netcracker/qubership-kafka` tag `2.6.4`
(`a3e27ea4a1d51e8dc81b438c8bc7f23cdfde811a`). That tag is not an applicability bound.

- `operator/charts/helm/kafka-service/values.yaml`
- `operator/charts/helm/kafka-service/templates/_helpers.tpl`
- `operator/controllers/provider/monitoring_provider.go`
- `docs/public/installation.md`


### Kafka Telegraf ServiceMonitor is absent

**Symptoms:**

- Deployment `kafka-monitoring` exists.
- ServiceMonitor `kafka-service-monitor` is absent, and cluster-status panels stay empty.

**Root cause:**

The Telegraf ServiceMonitor renders when `monitoring.serviceMonitorEnabled` is true and `monitoring.monitoringType` /
`global.monitoringType` is not `influxdb`. This is a missing scrape resource, not a missing Telegraf workload.

**Applies to:** kafka when the Telegraf monitoring workload exists. Image and chart versions are supporting facts.

**Failed layer:** `monitor missing or misconfigured`

**How to check:**

1. Confirm that `kafka-monitoring` exists and `kafka-service-monitor` is absent.
2. Name `monitoring.serviceMonitorEnabled` or `monitoring.monitoringType` only when effective values show the
   Prometheus monitor is off or type is `influxdb`.

**How to fix:**

1. Set `monitoring.serviceMonitorEnabled: true` and Prometheus monitoring type when those keys are present, then
   upgrade.
2. Confirm that `kafka-service-monitor` selects `component: kafka-monitoring` and port `prometheus-cli`.

**Owner:** kafka (`Netcracker/qubership-kafka`)

**Risk:** Enabling the ServiceMonitor adds a scrape target for Telegraf.

**Provenance:** Maintainer-verified at `Netcracker/qubership-kafka` tag `2.6.4`
(`a3e27ea4a1d51e8dc81b438c8bc7f23cdfde811a`). That tag is not an applicability bound.

- `operator/charts/helm/kafka-service/templates/service_monitor.yaml`
- `operator/charts/helm/kafka-service/values.yaml`


### Kafka JMX ServiceMonitor is absent

**Symptoms:**

- Kafka brokers are up, but JVM or JMX panels (`java_Memory_*`, `kafka_server_*`) stay empty.
- ServiceMonitor `kafka-service-monitor-jmx-exporter` is absent.

**Root cause:**

The JMX ServiceMonitor scrapes broker port `prometheus-http`. It is skipped when `global.externalKafka.enabled` is
true, and it also requires Prometheus ServiceMonitor enablement. This is a missing JMX scrape resource, not a missing
Telegraf Deployment.

**Applies to:** kafka when broker JMX metrics are expected. Image and chart versions are supporting facts.

**Failed layer:** `monitor missing or misconfigured`

**How to check:**

1. Confirm that `kafka-service-monitor-jmx-exporter` is absent.
2. Name `global.externalKafka.enabled` or ServiceMonitor enablement keys only when effective values show JMX scrape is
   skipped.

**How to fix:**

1. Disable `global.externalKafka.enabled` when that key is present and brokers are in-cluster, and enable the
   Prometheus ServiceMonitor, then upgrade.
2. Confirm that `kafka-service-monitor-jmx-exporter` selects `component: kafka` and port `prometheus-http`.

**Owner:** kafka (`Netcracker/qubership-kafka`)

**Risk:** Enabling JMX scrape adds one scrape per broker Service.

**Provenance:** Maintainer-verified at `Netcracker/qubership-kafka` tag `2.6.4`
(`a3e27ea4a1d51e8dc81b438c8bc7f23cdfde811a`). That tag is not an applicability bound.

- `operator/charts/helm/kafka-service/templates/service_monitor_jmx_exporter.yaml`
- `operator/controllers/provider/kafka_provider.go`


### Kafka JMX ServiceMonitor basicAuth breaks VMAgent

**Symptoms:**

- VMAgent or Prometheus fails to load scrape config with
  `job_name "serviceScrape/<namespace>/kafka-service-monitor-jmx-exporter/0": missing username in basic_auth`.
- VMAgent restarts.

**Root cause:**

At this chart tag, `kafka-service-monitor-jmx-exporter` emits `basicAuth` against Secret `kafka-secret` keys
`client-username` / `client-password` unless `global.secrets.kafka.disableSecurity` is true. Empty username in that
Secret makes VMAgent reject the job. This is a monitor misconfiguration, not shared-discovery monitor-excluded.

**Applies to:** kafka when `kafka-service-monitor-jmx-exporter` exists. Image and chart versions are supporting facts.

**Failed layer:** `monitor missing or misconfigured`

**How to check:**

1. Quote the VMAgent config error that names `kafka-service-monitor-jmx-exporter` and `basic_auth`.
2. Name `global.secrets.kafka.disableSecurity` or client username keys only when non-secret evidence shows empty
   basicAuth. Do not read Secret values.

**How to fix:**

1. Fill Kafka client username and password in `kafka-secret`, or set `global.secrets.kafka.disableSecurity: true` when
   that key is present and security is intentionally off, then upgrade.
2. Do not extend Monitoring Operator namespace selectors for a basicAuth parse error.

**Owner:** kafka (`Netcracker/qubership-kafka`)

**Risk:** Disabling Kafka security removes client auth. Filling credentials must not copy Secret values into the
diagnosis report.

**Provenance:** Maintainer-verified at `Netcracker/qubership-kafka` tag `2.6.4`
(`a3e27ea4a1d51e8dc81b438c8bc7f23cdfde811a`). That tag is not an applicability bound.

- `operator/charts/helm/kafka-service/templates/service_monitor_jmx_exporter.yaml`
- `operator/charts/helm/kafka-service/values.yaml`


### Kafka Grafana dashboards are absent

**Symptoms:**

- Bundled Kafka Monitoring or Kafka Topics panels never load.
- GrafanaDashboard CRs `kafka-grafana-dashboard` or `kafka-topics-grafana-dashboard` are absent.

**Root cause:**

Dashboard CRs are gated by `global.installDashboard` (default true), `monitoring.install`, and monitoring type not
`influxdb`. This is a missing dashboard CR, not a scrape failure.

**Applies to:** kafka when bundled Kafka Grafana panels never load. Image and chart versions are supporting facts.

**Failed layer:** `dashboard query mismatch`

**How to check:**

1. Confirm whether the expected GrafanaDashboard CR is absent.
2. Name `global.installDashboard` only when effective values show it is false.

**How to fix:**

1. Set `global.installDashboard: true` when that key is present, then upgrade.
2. Confirm that `kafka-grafana-dashboard` and `kafka-topics-grafana-dashboard` exist.

Empty panels after the CR exists are a later metric-path layer, not this case.

**Owner:** kafka (`Netcracker/qubership-kafka`)

**Risk:** Creating the GrafanaDashboard CRs loads bundled Kafka panels into Grafana.

**Provenance:** Maintainer-verified at `Netcracker/qubership-kafka` tag `2.6.4`
(`a3e27ea4a1d51e8dc81b438c8bc7f23cdfde811a`). That tag is not an applicability bound.

- `operator/charts/helm/kafka-service/templates/grafana_dashboard.yaml`
- `operator/charts/helm/kafka-service/templates/kafka-topics-dashboard.yaml`
- `docs/public/monitoring.md`


### Kafka additional Telegraf metrics are disabled

**Symptoms:**

- Telegraf is scraped and series `kafka_cluster_status` exists.
- Additional exec or cluster-derived series stay absent.

**Root cause:**

`monitoring.enableAdditionalMetrics` (default true) adds `/additional-metrics` to Telegraf `inputs.exec`. This is a
disabled series set, not a missing ServiceMonitor.

**Applies to:** kafka when Telegraf is scraped. Image and chart versions are supporting facts.

**Failed layer:** `metric series absent`

**How to check:**

1. Confirm that `kafka-service-monitor` is scraped and `kafka_cluster_status` exists.
2. Name `monitoring.enableAdditionalMetrics` only when effective values show it is false.

**How to fix:**

1. Set `monitoring.enableAdditionalMetrics: true` when that key is present, then upgrade.
2. Do not enable `monitoring.install` as the fix for missing additional series on a healthy scrape.

**Owner:** kafka (`Netcracker/qubership-kafka`)

**Risk:** Additional Telegraf exec metrics increase collection cost on Kafka.

**Provenance:** Maintainer-verified at `Netcracker/qubership-kafka` tag `2.6.4`
(`a3e27ea4a1d51e8dc81b438c8bc7f23cdfde811a`). That tag is not an applicability bound.

- `operator/charts/helm/kafka-service/templates/kafka-monitoring-configuration.yaml`
- `docs/public/installation.md`


### Kafka lag exporter is not installed

**Symptoms:**

- Kafka Exporter or lag dashboard is empty, and series `kafka_consumergroup_group_lag` is absent.
- Sidecar container `kafka-lag-exporter` is absent on `kafka-monitoring`.

**Root cause:**

Lag exporter is gated by `monitoring.lagExporter.enabled` (default false) and requires `monitoring.install`. This is a
missing component, not a scrape failure of Telegraf.

**Applies to:** kafka when consumer-lag metrics are expected. Image and chart versions are supporting facts.

**Failed layer:** `component configuration missing`

**How to check:**

1. Confirm that container `kafka-lag-exporter` is absent.
2. Name `monitoring.lagExporter.enabled` only when effective values show it is false.

**How to fix:**

1. Set `monitoring.lagExporter.enabled: true` when that key is present, then upgrade.
2. Confirm that sidecar `kafka-lag-exporter` and Service port `http` exist.

**Owner:** kafka (`Netcracker/qubership-kafka`)

**Risk:** Enabling lag exporter adds query load against Kafka consumer groups.

**Provenance:** Maintainer-verified at `Netcracker/qubership-kafka` tag `2.6.4`
(`a3e27ea4a1d51e8dc81b438c8bc7f23cdfde811a`). That tag is not an applicability bound.

- `operator/charts/helm/kafka-service/values.yaml`
- `operator/controllers/provider/lag_exporter_provider.go`
- `docs/public/monitoring.md`


### Kafka lag exporter ServiceMonitor is absent

**Symptoms:**

- Sidecar `kafka-lag-exporter` exists.
- ServiceMonitor `kafka-lag-exporter-service-monitor` is absent.

**Root cause:**

The lag exporter ServiceMonitor requires `monitoring.lagExporter.enabled`, `monitoring.serviceMonitorEnabled`, and
`monitoring.install`. This is a missing scrape resource for an installed lag exporter.

**Applies to:** kafka when the lag exporter sidecar exists. Image and chart versions are supporting facts.

**Failed layer:** `monitor missing or misconfigured`

**How to check:**

1. Confirm that `kafka-lag-exporter` exists and `kafka-lag-exporter-service-monitor` is absent.
2. Name `monitoring.serviceMonitorEnabled` only when effective values show it is false.

**How to fix:**

1. Set `monitoring.serviceMonitorEnabled: true` when that key is present, then upgrade.
2. Confirm that `kafka-lag-exporter-service-monitor` scrapes path `/metrics` on port `http`.

**Owner:** kafka (`Netcracker/qubership-kafka`)

**Risk:** Enabling the ServiceMonitor adds a scrape target for lag exporter.

**Provenance:** Maintainer-verified at `Netcracker/qubership-kafka` tag `2.6.4`
(`a3e27ea4a1d51e8dc81b438c8bc7f23cdfde811a`). That tag is not an applicability bound.

- `operator/charts/helm/kafka-service/templates/lag_exporter_service_monitor.yaml`


### Kafka Mirror Maker monitoring is not installed

**Symptoms:**

- Kafka Mirror Maker dashboard is empty, or series `kmm_health_status_code` is absent.
- Deployment `kafka-mirror-maker-monitoring` is absent.

**Root cause:**

Mirror Maker monitoring is gated by `mirrorMakerMonitoring.install` (default false) and `monitoring.install`.

**Applies to:** kafka when Mirror Maker metrics are expected. Image and chart versions are supporting facts.

**Failed layer:** `component configuration missing`

**How to check:**

1. Confirm that `kafka-mirror-maker-monitoring` is absent.
2. Name `mirrorMakerMonitoring.install` only when effective values show it is false.

**How to fix:**

1. Set `mirrorMakerMonitoring.install: true` when that key is present, then upgrade.
2. Confirm that Deployment `kafka-mirror-maker-monitoring` and ServiceMonitor `kafka-mirror-maker-service-monitor`
   exist.

**Owner:** kafka (`Netcracker/qubership-kafka`)

**Risk:** Enabling Mirror Maker monitoring deploys another Telegraf workload and scrape target.

**Provenance:** Maintainer-verified at `Netcracker/qubership-kafka` tag `2.6.4`
(`a3e27ea4a1d51e8dc81b438c8bc7f23cdfde811a`). That tag is not an applicability bound.

- `operator/charts/helm/kafka-service/values.yaml`
- `operator/charts/helm/kafka-service/templates/kmm_service_monitor.yaml`
- `docs/public/monitoring.md`


### Kafka backup daemon monitor is absent

**Symptoms:**

- Kafka backup health series are missing.
- ServiceMonitor `kafka-backup-daemon-service-monitor` is absent.

**Root cause:**

That ServiceMonitor renders when `backupDaemon.install` (default false) and `monitoring.install` are true.

**Applies to:** kafka when backup daemon metrics are expected. Image and chart versions are supporting facts.

**Failed layer:** `component configuration missing`

**How to check:**

1. Confirm that `kafka-backup-daemon-service-monitor` is absent.
2. Name `backupDaemon.install` only when effective values show it is false.

**How to fix:**

1. Set `backupDaemon.install: true` when that key is present, then upgrade.
2. Confirm that `kafka-backup-daemon-service-monitor` selects `name: kafka-backup-daemon` and path
   `/health/prometheus`.

**Owner:** kafka (`Netcracker/qubership-kafka`)

**Risk:** Enabling the backup daemon adds a scrape target.

**Provenance:** Maintainer-verified at `Netcracker/qubership-kafka` tag `2.6.4`
(`a3e27ea4a1d51e8dc81b438c8bc7f23cdfde811a`). That tag is not an applicability bound.

- `operator/charts/helm/kafka-service/templates/backup-daemon/backup-daemon-service-monitor.yaml`


### Kafka monitoring type is InfluxDB

**Symptoms:**

- Prometheus ServiceMonitors and Prometheus GrafanaDashboard CRs are never created.
- InfluxDB dashboards are expected instead.

**Root cause:**

`global.monitoringType` or `monitoring.monitoringType` set to `influxdb` skips Prometheus CR templates. Influx output
needs `global.smDbHost`, `global.smDbName`, and monitoring DB credentials.

**Applies to:** kafka when Prometheus scrape was expected. Image and chart versions are supporting facts.

**Failed layer:** `monitor missing or misconfigured`

**How to check:**

1. Confirm that Prometheus ServiceMonitors are absent while monitoring is installed.
2. Name `monitoring.monitoringType` or `global.monitoringType` only when effective values show `influxdb`.

**How to fix:**

1. Set monitoring type to `prometheus` when Prometheus scrape is required and that key is present, then upgrade.
2. Do not treat this as `monitoring.install: false`.

**Owner:** kafka (`Netcracker/qubership-kafka`)

**Risk:** Switching to Prometheus creates scrape targets and dashboard CRs.

**Provenance:** Maintainer-verified at `Netcracker/qubership-kafka` tag `2.6.4`
(`a3e27ea4a1d51e8dc81b438c8bc7f23cdfde811a`). That tag is not an applicability bound.

- `operator/charts/helm/kafka-service/templates/_helpers.tpl`
- `docs/public/installation.md`


## zookeeper

### ZooKeeper Monitoring is not installed

**Symptoms:**

- ZooKeeper dashboard panels have no data, or series `zookeeper_status_code` is absent.
- Deployment and Service `zookeeper-monitoring` are absent.

**Root cause:**

The zookeeper-service chart does not install the Telegraf monitoring workload. The documented enablement is
`monitoring.install` (default true). Helper `MONITORING_ENABLED` applies when `global.cloudIntegrationEnabled` is true.
This is a missing workload, not a scrape failure.

**Applies to:** zookeeper when ZooKeeper metrics are expected. Image and chart versions are supporting facts.

**Failed layer:** `component configuration missing`

**How to check:**

1. Confirm that Deployment `zookeeper-monitoring` is absent.
2. Name `monitoring.install` or `MONITORING_ENABLED` only when effective or rendered values show monitoring is off.

**How to fix:**

1. Set `monitoring.install: true` when that key is present in the evidence, then upgrade.
2. Confirm that Deployment and Service `zookeeper-monitoring` exist, with labels `component: zookeeper-monitoring`.

Do not change Monitoring Operator discovery selectors until the monitoring workload exists.

**Owner:** zookeeper (`Netcracker/qubership-zookeeper`)

**Risk:** Enabling monitoring deploys `zookeeper-monitoring` and starts Telegraf collection against ZooKeeper.

**Provenance:** Maintainer-verified at `Netcracker/qubership-zookeeper` tag `0.14.0`
(`c1ae51bfcaab220bb779540db796093a503269c8`). That tag is not an applicability bound.

- `operator/charts/helm/zookeeper-service/values.yaml`
- `operator/charts/helm/zookeeper-service/templates/_helpers.tpl`
- `operator/controllers/provider/monitoring_provider.go`
- `docs/public/monitoring.md`


### ZooKeeper Grafana dashboard is absent

**Symptoms:**

- Bundled ZooKeeper panels never load.
- GrafanaDashboard CR `zookeeper-grafana-dashboard` is absent.

**Root cause:**

The dashboard CR is gated by `monitoring.installGrafanaDashboard` (default true) and `monitoring.install`. This is a
missing dashboard CR, not a scrape failure. JMX ServiceMonitor `zookeeper-service-monitor-jmx-exporter` does not emit
basicAuth at this tag.

**Applies to:** zookeeper when bundled ZooKeeper Grafana panels never load. Image and chart versions are supporting
facts.

**Failed layer:** `dashboard query mismatch`

**How to check:**

1. Confirm that `zookeeper-grafana-dashboard` is absent.
2. Name `monitoring.installGrafanaDashboard` only when effective values show it is false.

**How to fix:**

1. Set `monitoring.installGrafanaDashboard: true` when that key is present, then upgrade.
2. Confirm that CR `zookeeper-grafana-dashboard` exists.

**Owner:** zookeeper (`Netcracker/qubership-zookeeper`)

**Risk:** Creating the GrafanaDashboard CR loads bundled ZooKeeper panels into Grafana.

**Provenance:** Maintainer-verified at `Netcracker/qubership-zookeeper` tag `0.14.0`
(`c1ae51bfcaab220bb779540db796093a503269c8`). That tag is not an applicability bound.

- `operator/charts/helm/zookeeper-service/templates/grafana_dashboard.yaml`
- `docs/public/zookeeper-dashboard.md`
- `operator/charts/helm/zookeeper-service/values.yaml`


### ZooKeeper backup metrics are missing

**Symptoms:**

- ZooKeeper backup dashboard section is empty, or series `zookeeper_backup_metric_*` is absent.
- Deployment `zookeeper-backup-daemon` is often absent.

**Root cause:**

There is no backup ServiceMonitor. Backup series come from Telegraf exec `backup_metric.py` when
`backupDaemon.install` is true (default false) so `ZOOKEEPER_BACKUP_DAEMON_HOST` is set. This is missing backup
component configuration, not a missing Telegraf ServiceMonitor.

**Applies to:** zookeeper when backup metrics are expected. Image and chart versions are supporting facts.

**Failed layer:** `component configuration missing`

**How to check:**

1. Confirm that Telegraf is scraped and `zookeeper_backup_metric_*` is absent.
2. Name `backupDaemon.install` only when effective values show it is false.

**How to fix:**

1. Set `backupDaemon.install: true` when that key is present, then upgrade.
2. Confirm that Deployment `zookeeper-backup-daemon` exists and backup series appear on the Telegraf scrape.

**Owner:** zookeeper (`Netcracker/qubership-zookeeper`)

**Risk:** Enabling the backup daemon deploys backup workloads and adds Telegraf exec against it.

**Provenance:** Maintainer-verified at `Netcracker/qubership-zookeeper` tag `0.14.0`
(`c1ae51bfcaab220bb779540db796093a503269c8`). That tag is not an applicability bound.

- `operator/charts/helm/zookeeper-service/templates/zookeeper-monitoring-configuration.yaml`
- `monitoring/exec-scripts/backup_metric.py`
- `docs/public/zookeeper-dashboard.md`


## consul

### Consul monitoring CRs are not installed

**Symptoms:**

- Consul dashboard never appears.
- ServiceMonitor `consul-server-service-monitor` and GrafanaDashboard CR `consul-grafana-dashboard` are absent.

**Root cause:**

Prometheus CRs render when helper `monitoring.enabled` is true (`monitoring.enabled` default true, or
`MONITORING_ENABLED` when `global.cloudIntegrationEnabled`). Consul servers can still expose `/v1/agent/metrics` via
pod annotations. This is missing scrape and dashboard CRs, not missing Consul servers. There is no Telegraf sidecar at
this tag.

**Applies to:** consul when Consul Prometheus monitoring CRs are expected. Image and chart versions are supporting
facts.

**Failed layer:** `monitor missing or misconfigured`

**How to check:**

1. Confirm that `consul-server-service-monitor` and `consul-grafana-dashboard` are absent.
2. Name `monitoring.enabled` or `MONITORING_ENABLED` only when effective values show monitoring is off.

**How to fix:**

1. Set `monitoring.enabled: true` when that key is present, then upgrade.
2. Confirm that `consul-server-service-monitor` scrapes `/v1/agent/metrics?format=prometheus`.

Do not change Monitoring Operator selectors until the ServiceMonitor exists.

**Owner:** consul (`Netcracker/qubership-consul`)

**Risk:** Enabling monitoring CRs adds scrape of Consul agent metrics and loads the Consul dashboard.

**Provenance:** Maintainer-verified at `Netcracker/qubership-consul` tag `0.14.2`
(`f462eb4d92295048bb5c9d9ad1f4698480cbb8d6`). That tag is not an applicability bound.

- `charts/helm/consul-service/templates/service-monitor.yaml`
- `charts/helm/consul-service/templates/grafana-dashboard.yaml`
- `charts/helm/consul-service/templates/_helpers.tpl`
- `docs/public/installation.md`


### Consul agent Prometheus retention is disabled

**Symptoms:**

- ServiceMonitor `consul-server-service-monitor` exists and is scraped, but `/v1/agent/metrics` is empty.
- Consul telemetry panels (`consul.raft.*`, `consul.runtime.*`) stay empty.

**Root cause:**

Agent Prometheus telemetry is configured only when `global.metrics.enabled` and `global.metrics.enableAgentMetrics`
are true (defaults true) and `global.metrics.agentMetricsRetentionTime` is greater than `0` (default `24h`). This is
missing series at the agent, not a missing ServiceMonitor.

**Applies to:** consul when the server ServiceMonitor is scraped. Image and chart versions are supporting facts.

**Failed layer:** `metric series absent`

**How to check:**

1. Confirm that `consul-server-service-monitor` is scraped.
2. Name `global.metrics.enableAgentMetrics` or `global.metrics.agentMetricsRetentionTime` only when effective values
   show Prometheus retention is off.

**How to fix:**

1. Set `global.metrics.enabled` and `global.metrics.enableAgentMetrics` to true, and keep
   `global.metrics.agentMetricsRetentionTime` greater than `0`, when those keys are present, then upgrade.
2. Do not enable `monitoring.enabled` as the fix for empty agent metrics on a healthy scrape.

**Owner:** consul (`Netcracker/qubership-consul`)

**Risk:** Enabling agent Prometheus retention stores telemetry in the Consul process.

**Provenance:** Maintainer-verified at `Netcracker/qubership-consul` tag `0.14.2`
(`f462eb4d92295048bb5c9d9ad1f4698480cbb8d6`). That tag is not an applicability bound.

- `charts/helm/consul-service/templates/_helpers.tpl`
- `charts/helm/consul-service/templates/server-statefulset.yaml`
- `docs/public/installation.md`


### Consul ACL token scrape failure

**Symptoms:**

- ServiceMonitor `consul-server-service-monitor` exists.
- Scrape of `/v1/agent/metrics` returns 401 or 403, and Consul series stay absent.

**Root cause:**

`global.acls.manageSystemACLs` defaults true, so the ServiceMonitor sets `bearerTokenSecret` (default Secret
`consul-bootstrap-acl-token` key `token`). A missing or wrong token is scrape auth failure, not kubectl access-failure
and not monitor-excluded.

**Applies to:** consul when the server ServiceMonitor exists. Image and chart versions are supporting facts.

**Failed layer:** `target unavailable`

**How to check:**

1. Confirm that the ServiceMonitor exists and scrape fails with 401 or 403 on `/v1/agent/metrics`.
2. Name ACL Secret keys only from non-secret evidence (`global.acls.bootstrapToken.secretName`). Do not read Secret
   values.

**How to fix:**

1. Point `bearerTokenSecret` at a token Secret that can read agent metrics, without copying the token into the report.
2. Do not extend Monitoring Operator namespace selectors for an ACL scrape failure.

**Owner:** consul (`Netcracker/qubership-consul`)

**Risk:** Granting a scrape token widens ACL access. Do not copy Secret values into the diagnosis report.

**Provenance:** Maintainer-verified at `Netcracker/qubership-consul` tag `0.14.2`
(`f462eb4d92295048bb5c9d9ad1f4698480cbb8d6`). That tag is not an applicability bound.

- `charts/helm/consul-service/templates/service-monitor.yaml`
- `charts/helm/consul-service/values.yaml`


### Consul client ServiceMonitor is absent

**Symptoms:**

- Consul client metrics are expected.
- ServiceMonitor `consul-client-service-monitor` is absent.

**Root cause:**

Client scrape renders only when `client.enabled` (default false) and `monitoring.enabled` are true.

**Applies to:** consul when client agent metrics are expected. Image and chart versions are supporting facts.

**Failed layer:** `monitor missing or misconfigured`

**How to check:**

1. Confirm that `consul-client-service-monitor` is absent.
2. Name `client.enabled` only when effective values show clients are off.

**How to fix:**

1. Set `client.enabled: true` when that key is present and clients are required, then upgrade.
2. Confirm that `consul-client-service-monitor` selects `component: client`.

**Owner:** consul (`Netcracker/qubership-consul`)

**Risk:** Enabling Consul clients deploys a client DaemonSet and adds scrape of client agents.

**Provenance:** Maintainer-verified at `Netcracker/qubership-consul` tag `0.14.2`
(`f462eb4d92295048bb5c9d9ad1f4698480cbb8d6`). That tag is not an applicability bound.

- `charts/helm/consul-service/templates/service-monitor-client.yaml`
- `charts/helm/consul-service/values.yaml`


### Consul Grafana dashboard is absent

**Symptoms:**

- Bundled Consul panels never load.
- GrafanaDashboard CR `consul-grafana-dashboard` is absent.

**Root cause:**

The dashboard CR is gated by `monitoring.installDashboard` (default true) and `monitoring.enabled`. This is a missing
dashboard CR, not a scrape failure.

**Applies to:** consul when bundled Consul Grafana panels never load. Image and chart versions are supporting facts.

**Failed layer:** `dashboard query mismatch`

**How to check:**

1. Confirm that `consul-grafana-dashboard` is absent.
2. Name `monitoring.installDashboard` only when effective values show it is false.

**How to fix:**

1. Set `monitoring.installDashboard: true` when that key is present, then upgrade.
2. Confirm that CR `consul-grafana-dashboard` exists.

**Owner:** consul (`Netcracker/qubership-consul`)

**Risk:** Creating the GrafanaDashboard CR loads bundled Consul panels into Grafana.

**Provenance:** Maintainer-verified at `Netcracker/qubership-consul` tag `0.14.2`
(`f462eb4d92295048bb5c9d9ad1f4698480cbb8d6`). That tag is not an applicability bound.

- `charts/helm/consul-service/templates/grafana-dashboard.yaml`
- `docs/public/monitoring.md`


### Consul backup daemon monitor is absent

**Symptoms:**

- Consul backup daemon metrics are expected.
- ServiceMonitor `consul-backup-daemon-service-monitor` is absent.

**Root cause:**

That ServiceMonitor renders when `backupDaemon.enabled` (default false) and `monitoring.enabled` are true.

**Applies to:** consul when backup daemon metrics are expected. Image and chart versions are supporting facts.

**Failed layer:** `component configuration missing`

**How to check:**

1. Confirm that `consul-backup-daemon-service-monitor` is absent.
2. Name `backupDaemon.enabled` only when effective values show it is false.

**How to fix:**

1. Set `backupDaemon.enabled: true` when that key is present, then upgrade.
2. Confirm that `consul-backup-daemon-service-monitor` selects `component: backup-daemon` and path
   `/health/prometheus`.

**Owner:** consul (`Netcracker/qubership-consul`)

**Risk:** Enabling the backup daemon adds a scrape target.

**Provenance:** Maintainer-verified at `Netcracker/qubership-consul` tag `0.14.2`
(`f462eb4d92295048bb5c9d9ad1f4698480cbb8d6`). That tag is not an applicability bound.

- `charts/helm/consul-service/templates/backup-daemon-service-monitor.yaml`


## rabbitmq

### RabbitMQ Prometheus monitoring is not installed

**Symptoms:**

- RabbitMQ Overview dashboard is empty, or alert `NoMetrics` fires for missing `rabbitmq_identity_info`.
- ServiceMonitor `rabbitmq-service-monitor` is absent, and plugin `rabbitmq_prometheus` is not enabled.

**Root cause:**

Native Prometheus monitoring is gated by helper `monitoring.enabled`: `rabbitmqPrometheusMonitoring` (default false;
used in templates and docs, not present in values.yaml) or `MONITORING_ENABLED` when cloud integration is on. The
operator Service exposes port `15692-tcp`. This is a missing plugin and scrape resource, not external Telegraf.

**Applies to:** rabbitmq when in-cluster Prometheus metrics are expected. Image and chart versions are supporting
facts.

**Failed layer:** `component configuration missing`

**How to check:**

1. Confirm that `rabbitmq-service-monitor` is absent and `externalRabbitmq.enabled` is not the path in use.
2. Name `rabbitmqPrometheusMonitoring` or `MONITORING_ENABLED` only when effective or rendered values show monitoring
   is off.

**How to fix:**

1. Set `rabbitmqPrometheusMonitoring: true` when that key is present in the evidence, then upgrade.
2. Confirm that plugin `rabbitmq_prometheus` is enabled and `rabbitmq-service-monitor` selects `app: rmqlocal` and
   port `15692-tcp`.

Do not change Monitoring Operator selectors until the ServiceMonitor exists.

**Owner:** rabbitmq (`Netcracker/qubership-rabbitmq`)

**Risk:** Enabling the Prometheus plugin exposes `/metrics` on port 15692 and adds a scrape target.

**Provenance:** Maintainer-verified at `Netcracker/qubership-rabbitmq` tag `2.4.0`
(`48f0931e17ddb378b688c6cb72e0d6418bbff1d3`). That tag is not an applicability bound.

- `operator/charts/helm/rabbitmq/templates/_helpers.tpl`
- `operator/charts/helm/rabbitmq/templates/rabbitmq_service_monitor.yaml`
- `docs/public/monitoring.md`
- `docs/public/troubleshooting.md`


### RabbitMQ per-queue metrics are disabled

**Symptoms:**

- RabbitMQ Queues dashboard is empty (`rabbitmq_queue_messages_*` per queue).
- Overview aggregated metrics may still work. ServiceMonitor `rabbitmq-per-object-service-monitor` is absent.

**Root cause:**

Per-object scrape is gated by `rabbitmq.perQueueMetrics` (default false; in templates and docs, not in values.yaml)
and requires native monitoring. Docs warn of resource cost. This is a disabled series path, not missing Overview
monitoring.

**Applies to:** rabbitmq when per-queue Prometheus metrics are expected. Image and chart versions are supporting
facts.

**Failed layer:** `metric series absent`

**How to check:**

1. Confirm that `rabbitmq-service-monitor` is scraped and `rabbitmq-per-object-service-monitor` is absent.
2. Name `rabbitmq.perQueueMetrics` only when effective values show it is false.

**How to fix:**

1. Set `rabbitmq.perQueueMetrics: true` when that key is present and the component team accepts the cost, then
   upgrade.
2. Confirm that `rabbitmq-per-object-service-monitor` scrapes `/metrics/per-object` on port `15692-tcp`.

**Owner:** rabbitmq (`Netcracker/qubership-rabbitmq`)

**Risk:** Per-queue metrics increase scrape cardinality and RabbitMQ load.

**Provenance:** Maintainer-verified at `Netcracker/qubership-rabbitmq` tag `2.4.0`
(`48f0931e17ddb378b688c6cb72e0d6418bbff1d3`). That tag is not an applicability bound.

- `operator/charts/helm/rabbitmq/templates/rabbitmq_per_object_service_monitor.yaml`
- `docs/public/monitoring.md`


### RabbitMQ external Telegraf monitoring is not installed

**Symptoms:**

- `externalRabbitmq.enabled` is true.
- Deployment `rabbit-monitoring-external` and ServiceMonitor `rabbitmq-service-monitor-external` are absent.

**Root cause:**

External Telegraf monitoring requires both `externalRabbitmq.enabled` (default false) and `telegraf.install` (default
false). In-cluster `telegraf.install` without external RabbitMQ creates an Influx-oriented Telegraf Deployment with no
Prometheus ServiceMonitor. This is a missing external monitoring component.

**Applies to:** rabbitmq when external RabbitMQ Prometheus metrics are expected. Image and chart versions are
supporting facts.

**Failed layer:** `component configuration missing`

**How to check:**

1. Confirm that `externalRabbitmq.enabled` is the path in use and `rabbit-monitoring-external` is absent.
2. Name `telegraf.install` only when effective values show it is false.

**How to fix:**

1. Set `telegraf.install: true` with `externalRabbitmq.enabled: true` when those keys are present, then upgrade.
2. Confirm that Deployment `rabbit-monitoring-external` and ServiceMonitor `rabbitmq-service-monitor-external` exist.

**Owner:** rabbitmq (`Netcracker/qubership-rabbitmq`)

**Risk:** Enabling external Telegraf deploys a monitoring workload against the external cluster.

**Provenance:** Maintainer-verified at `Netcracker/qubership-rabbitmq` tag `2.4.0`
(`48f0931e17ddb378b688c6cb72e0d6418bbff1d3`). That tag is not an applicability bound.

- `operator/charts/helm/rabbitmq/templates/rabbit_monitoring_external_deployment.yaml`
- `operator/charts/helm/rabbitmq/templates/rabbitmq_monitoring_external_servicemonitor.yaml`
- `docs/public/installation.md`


### RabbitMQ backup daemon monitor is absent

**Symptoms:**

- RabbitMQ backup daemon metrics are expected.
- ServiceMonitor `rabbitmq-backup-daemon-service-monitor` is absent.

**Root cause:**

That ServiceMonitor renders when `backupDaemon.enabled` (default false) and `monitoring.enabled` are true.

**Applies to:** rabbitmq when backup daemon metrics are expected. Image and chart versions are supporting facts.

**Failed layer:** `component configuration missing`

**How to check:**

1. Confirm that `rabbitmq-backup-daemon-service-monitor` is absent.
2. Name `backupDaemon.enabled` only when effective values show it is false.

**How to fix:**

1. Set `backupDaemon.enabled: true` when that key is present, then upgrade.
2. Confirm that `rabbitmq-backup-daemon-service-monitor` selects `app: rabbitmq-backup-daemon` and path
   `/health/prometheus`.

**Owner:** rabbitmq (`Netcracker/qubership-rabbitmq`)

**Risk:** Enabling the backup daemon adds a scrape target.

**Provenance:** Maintainer-verified at `Netcracker/qubership-rabbitmq` tag `2.4.0`
(`48f0931e17ddb378b688c6cb72e0d6418bbff1d3`). That tag is not an applicability bound.

- `operator/charts/helm/rabbitmq/templates/backup-daemon-service-monitor.yaml`

## DRNavigator

### Site-manager monitor is not rendered

**Symptoms:**

- SM Monitoring dashboard is empty.
- Deployment `site-manager` exists, but ServiceMonitor `site-manager-service-monitor` and/or GrafanaDashboard CR
  `site-manager-grafana-dashboard` are absent.

**Root cause:**

The site-manager workload always binds metrics on port `metrics` (9000). `MONITORING_ENABLED` (default true) only
creates the ServiceMonitor and GrafanaDashboard CRs. This is missing scrape and dashboard resources, not a missing
site-manager Deployment.

**Applies to:** DRNavigator when site-manager Prometheus metrics are expected. Image and chart versions are supporting
facts.

**Failed layer:** `monitor missing or misconfigured`

**How to check:**

1. Confirm that `site-manager` exists and `site-manager-service-monitor` is absent.
2. Name `MONITORING_ENABLED` only when effective or rendered values show it is false.

**How to fix:**

1. Set `MONITORING_ENABLED: true` when that key is present, then upgrade.
2. Confirm that `site-manager-service-monitor` selects `app: site-manager`, port `metrics`, and path `/metrics`.

Do not change Monitoring Operator selectors until the ServiceMonitor exists.

**Owner:** DRNavigator (`Netcracker/DRNavigator`)

**Risk:** Enabling monitoring CRs adds scrape of site-manager `/metrics` and loads the SM Monitoring dashboard.

**Provenance:** Maintainer-verified at `Netcracker/DRNavigator` tag `0.33.0`
(`49e7d4c02ab384ef36e098d89c5944181ea166a0`). That tag is not an applicability bound.

- `charts/site-manager/values.yaml`
- `charts/site-manager/templates/service-monitor.yaml`
- `charts/site-manager/templates/grafana-dashboard.yaml`
- `documentation/public/installation.md`


### paas-geo-monitor is not installed

**Symptoms:**

- Paas-Geo-Monitor Monitoring dashboard is empty, or series `peer_*` / `paas_geo_monitor_*` are absent.
- Deployment `paas-geo-monitor` is absent.

**Root cause:**

paas-geo-monitor is a second workload gated by `paasGeoMonitor.install` (default false). Stop here; do not treat an
absent geo ServiceMonitor as a scrape failure while the workload is missing.

**Applies to:** DRNavigator when paas-geo-monitor metrics are expected. Image and chart versions are supporting facts.

**Failed layer:** `component configuration missing`

**How to check:**

1. Confirm that Deployment `paas-geo-monitor` is absent.
2. Name `paasGeoMonitor.install` only when effective values show it is false.

**How to fix:**

1. Set `paasGeoMonitor.install: true` when that key is present, then upgrade.
2. Confirm that Deployment and Service `paas-geo-monitor` exist, with port `web`.

**Owner:** DRNavigator (`Netcracker/DRNavigator`)

**Risk:** Enabling paas-geo-monitor deploys a connectivity probe workload.

**Provenance:** Maintainer-verified at `Netcracker/DRNavigator` tag `0.33.0`
(`49e7d4c02ab384ef36e098d89c5944181ea166a0`). That tag is not an applicability bound.

- `charts/site-manager/values.yaml`
- `charts/site-manager/templates/paas-geo-monitor-deployment.yaml`
- `documentation/public/installation.md`


### paas-geo-monitor scrape is not rendered

**Symptoms:**

- Deployment `paas-geo-monitor` exists.
- ServiceMonitor `paas-geo-monitor-service-monitor` and/or GrafanaDashboard CR `paas-geo-monitor-grafana-dashboard` are
  absent.

**Root cause:**

Geo ServiceMonitor and dashboard templates require both `paasGeoMonitor.install` and `MONITORING_ENABLED`. This is
missing scrape and dashboard resources for an installed geo workload.

**Applies to:** DRNavigator when paas-geo-monitor is installed. Image and chart versions are supporting facts.

**Failed layer:** `monitor missing or misconfigured`

**How to check:**

1. Confirm that `paas-geo-monitor` exists and `paas-geo-monitor-service-monitor` is absent.
2. Name `MONITORING_ENABLED` only when effective values show it is false.

**How to fix:**

1. Set `MONITORING_ENABLED: true` when that key is present, then upgrade.
2. Confirm that `paas-geo-monitor-service-monitor` selects `app: paas-geo-monitor`, port `web`, and path `/metrics`.

**Owner:** DRNavigator (`Netcracker/DRNavigator`)

**Risk:** Enabling geo monitoring CRs adds a scrape target for paas-geo-monitor.

**Provenance:** Maintainer-verified at `Netcracker/DRNavigator` tag `0.33.0`
(`49e7d4c02ab384ef36e098d89c5944181ea166a0`). That tag is not an applicability bound.

- `charts/site-manager/templates/paas-geo-monitor-service-monitor.yaml`
- `charts/site-manager/templates/grafana-dashboard.yaml`


### paas-geo-monitor peer series are absent

**Symptoms:**

- paas-geo-monitor ServiceMonitor exists and is scraped.
- Series `peer_dns_status`, `peer_svc_status`, or `peer_pod_status` are absent, while health or CPU panels may work.

**Root cause:**

Labeled peer gauges are written only by the ping loop over configured peers. Docs mark
`paasGeoMonitor.env.PAAS_PING_PEERS` as required for monitoring (default true). `paasGeoMonitor.config.peers` defaults
to `[]`, so no peer series are emitted until peers are configured.

**Applies to:** DRNavigator when paas-geo-monitor is scraped. Image and chart versions are supporting facts.

**Failed layer:** `metric series absent`

**How to check:**

1. Confirm that `paas-geo-monitor-service-monitor` is scraped.
2. Name `PAAS_PING_PEERS` or `paasGeoMonitor.config.peers` only when effective values show ping is off or the peer list
   is empty.

**How to fix:**

1. Set `paasGeoMonitor.env.PAAS_PING_PEERS: true` and configure `paasGeoMonitor.config.peers` when those keys are
   present, then upgrade.
2. Do not set `paasGeoMonitor.install: true` as the fix for missing peer series on a healthy scrape.

**Owner:** DRNavigator (`Netcracker/DRNavigator`)

**Risk:** Enabling peer pings sends DNS, Service, and pod probes to each configured peer.

**Provenance:** Maintainer-verified at `Netcracker/DRNavigator` tag `0.33.0`
(`49e7d4c02ab384ef36e098d89c5944181ea166a0`). That tag is not an applicability bound.

- `charts/site-manager/values.yaml`
- `paas-geo-monitor/pkg/app/app.go`
- `documentation/public/installation.md`


### paas-geo-monitor BGP series are absent

**Symptoms:**

- paas-geo-monitor target is healthy.
- Series `paas_geo_monitor_bgp_peer` or `paas_geo_monitor_bgp_route` are absent, and BGP or Service Network panels stay
  empty.

**Root cause:**

BGP series register only when env `PAAS_BGP_METRICS` is `true` (`paasGeoMonitor.env.PAAS_BGP_METRICS`, values default
`"true"`). This is a disabled series set, not missing peer ping.

**Applies to:** DRNavigator when paas-geo-monitor is scraped. Image and chart versions are supporting facts.

**Failed layer:** `metric series absent`

**How to check:**

1. Confirm that the geo target is healthy and BGP series are absent.
2. Name `PAAS_BGP_METRICS` only when effective values or env show it is not `true`.

**How to fix:**

1. Set `paasGeoMonitor.env.PAAS_BGP_METRICS: "true"` when that key is present, then upgrade.
2. Do not treat this as an empty `paasGeoMonitor.config.peers` list.

**Owner:** DRNavigator (`Netcracker/DRNavigator`)

**Risk:** Enabling BGP metrics queries BGP state on the node.

**Provenance:** Maintainer-verified at `Netcracker/DRNavigator` tag `0.33.0`
(`49e7d4c02ab384ef36e098d89c5944181ea166a0`). That tag is not an applicability bound.

- `charts/site-manager/values.yaml`
- `charts/site-manager/templates/paas-geo-monitor-deployment.yaml`
- `paas-geo-monitor/pkg/app/app.go`


### paas-geo-monitor service CIDR dashboard regex mismatch

**Symptoms:**

- BGP series exist.
- Service Network Availability variable or panels stay empty because `$service_cidr` does not match.

**Root cause:**

Dashboard variable `service_cidr` uses regex from `paasGeoMonitor.monitoring.serviceCIDRRegex` (default `172.*`). A
cluster service subnet outside that regex leaves those panels empty. This is a dashboard query mismatch, not missing
BGP series.

**Applies to:** DRNavigator when paas-geo-monitor BGP series exist. Image and chart versions are supporting facts.

**Failed layer:** `dashboard query mismatch`

**How to check:**

1. Confirm that BGP series exist and the empty panels use `$service_cidr`.
2. Name `paasGeoMonitor.monitoring.serviceCIDRRegex` only when effective values or the dashboard JSON show a regex that
   does not match the service subnet.

**How to fix:**

1. Set `paasGeoMonitor.monitoring.serviceCIDRRegex` to the service subnet regex when that key is present, then
   upgrade.
2. Do not enable `PAAS_BGP_METRICS` as the fix for a regex mismatch.

**Owner:** DRNavigator (`Netcracker/DRNavigator`)

**Risk:** A broader CIDR regex shows more BGP service networks in Grafana.

**Provenance:** Maintainer-verified at `Netcracker/DRNavigator` tag `0.33.0`
(`49e7d4c02ab384ef36e098d89c5944181ea166a0`). That tag is not an applicability bound.

- `charts/site-manager/values.yaml`
- `charts/site-manager/templates/grafana-dashboard.yaml`
- `charts/site-manager/monitoring/paas-geo-monitor-grafana-dashboard.json`


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

## Shared discovery

### Monitor excluded

**Symptoms:**

- A component `ServiceMonitor` or `PodMonitor` exists in the component namespace.
- The namespace label does not match the Prometheus or VMAgent `serviceMonitorNamespaceSelector`.
- The selector allowlist includes namespaces such as `monitoring` and `cqp` and omits the component namespace.

**Root cause:**

Monitoring Operator installation values restrict which namespaces VMAgent or Prometheus may discover. A valid
component monitor outside that scope is never scraped.

**Applies to:** any supported component after its scrape resource exists.

**Failed layer:** `monitor excluded`

**How to check:**

1. Confirm from the ticket that the component scrape resource exists. If the operator still needs to collect it,
   they can run:

   ```bash
   kubectl --context <context> -n <component-namespace> get servicemonitors,podmonitors
   kubectl --context <context> -n <component-namespace> get servicemonitor <name> -o yaml
   ```

2. Locate the `PlatformMonitoring` resource and read its VMAgent or Prometheus namespace and label selectors from
   pasted YAML, or ask the operator to collect:

   ```bash
   kubectl --context <context> get platformmonitoring -A
   kubectl --context <context> -n <monitoring-namespace> get platformmonitoring <name> -o yaml
   ```

3. Compare those selectors with the monitor namespace labels. Stop here when they do not match.

If the ticket already quotes the selector and the monitor namespace, that is enough. Do not enable a component
exporter flag to fix this layer. Do not run `kubectl` from this skill.

**How to fix:**

1. Read `victoriametrics.vmAgent.serviceMonitorNamespaceSelector` (or the Prometheus equivalent) in the
   Monitoring Operator installation values.
2. Ask the Monitoring Operator owners to extend that selector, or to apply a namespace label the selector
   already requires.
3. After they apply the values, re-read the PlatformMonitoring spec and confirm it selects the component
   namespace.

**Owner:** Monitoring Operator installation values

**Risk:** Expanding the namespace selector causes VMAgent to discover ServiceMonitors in every extra matching
namespace. Confirm that discovery scope before applying the values change.

**Provenance:** Monitoring Operator docs `docs/monitoring-configuration/limits-metric-collection.md` and
`docs/installation/components/victoriametrics-stack/vmagent.md`.


### Target unavailable

**Symptoms:**

- The expected `ServiceMonitor` or `PodMonitor` exists and is in discovery scope.
- The selected Service has no ready endpoints, or the endpoint port does not match the monitor.

**Root cause:**

The monitor does not select a ready Service port. The scrape target never becomes healthy.

**Applies to:** any supported component after the monitor is discovered.

**Failed layer:** `target unavailable`

**How to check:**

Read the pasted Service and EndpointSlice YAML. If the operator still needs to collect it, they can run:

```bash
kubectl --context <context> -n <component-namespace> get service <name> -o yaml
kubectl --context <context> -n <component-namespace> get endpointslice -l kubernetes.io/service-name=<name> -o yaml
```

Stop when the Service selector, endpoint port, or ready endpoints do not match the monitor. A scrape timeout on a
ready target is the component collection-timeout case (pgskipper `Collection scrape timeout`, MongoDB `MongoDB
collection scrape timeout`, or ClickHouse `ClickHouse collection scrape timeout`), not this one. Do not run `kubectl`
from this skill.

**How to fix:**

1. Ask the component team to align the Service selector and the named metrics port with the monitor.
2. Do not change Monitoring Operator selectors for a target that is not ready.

**Owner:** component chart or workload

**Risk:** None — every recommended step is read-only

**Provenance:** none beyond the live Service and EndpointSlice YAML.


### Access failure

**Symptoms:**

- The pasted output shows that the named Kubernetes context, namespace, resource type, or permission is unavailable.
- Example: `error: context "production-eu" does not exist`.

**Root cause:**

The diagnosis cannot inspect the named context. Continue from ticket attachments or pasted non-secret evidence.

**Applies to:** after the ticket quotes a named `--context` and namespace, before cluster resources can be read.

**Failed layer:** `access failure`

**How to check:**

Quote the command and the error verbatim from the ticket. Do not run `kubectl`. Do not run
`kubectl config use-context`. Do not retry against another context.

**How to fix:**

Request one of:

- The ticket description and non-secret attachments.
- The complete image reference from the component workload.
- The non-secret workload YAML showing `.spec.template.spec.containers[*].image`.
- The same read-only `kubectl --context <name>` output from a kubeconfig that already contains that context.

Do not ask for Secret data, credentials, or kubeconfig content.

**Owner:** none. Diagnosis stopped before a metric-path layer.

**Risk:** None — every recommended step is read-only

**Provenance:** none. Load no component case until evidence from the ticket is available.
