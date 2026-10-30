# Infrastructure component troubleshooting


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
