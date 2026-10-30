# Infrastructure component troubleshooting


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
