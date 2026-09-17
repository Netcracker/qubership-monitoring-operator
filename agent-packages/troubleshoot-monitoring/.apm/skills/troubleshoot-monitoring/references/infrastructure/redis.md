# Infrastructure component troubleshooting


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
