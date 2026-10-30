# Infrastructure component troubleshooting


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
