# Infrastructure component troubleshooting


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
