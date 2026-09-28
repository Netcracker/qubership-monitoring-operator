# Infrastructure component troubleshooting


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
