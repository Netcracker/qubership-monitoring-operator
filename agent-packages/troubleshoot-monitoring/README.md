# troubleshoot-monitoring

A single user-invoked skill that diagnoses problems with Qubership Monitoring Operator (a Kubernetes operator that
reconciles the `PlatformMonitoring` custom resource and installs a full monitoring stack: prometheus-operator, the
VictoriaMetrics stack, Prometheus, Alertmanager, Grafana, exporters, prometheus-adapter, and integrations).

The skill is **read-only and advisory**. It does not run `kubectl`, SSH, or Ansible, and it never changes a system. It
reads a pasted problem description plus any attached logs or configuration, matches the symptom against a curated
reference, and returns a diagnosis with remediation steps and a list of data to collect when the match is uncertain.

When no Monitoring Operator case matches, the skill checks the infrastructure catalog. Supported components are
pgskipper-operator, mongodb-operator, cassandra-operator, redis, clickhouse-operator-helm, kafka, zookeeper, consul,
rabbitmq, DRNavigator, and opensearch.

## Contents

| Path | Purpose |
| ---- | ------- |
| [`SKILL.md`](.apm/skills/troubleshoot-monitoring/SKILL.md) | The diagnosis procedure. |
| [`references/troubleshooting.md`](.apm/skills/troubleshoot-monitoring/references/troubleshooting.md) | Symptom-indexed Monitoring Operator catalog. |
| [`references/infrastructure/troubleshooting.md`](.apm/skills/troubleshoot-monitoring/references/infrastructure/troubleshooting.md) | Symptom-indexed infrastructure catalog. |
| [`scripts/show_cases.py`](.apm/skills/troubleshoot-monitoring/scripts/show_cases.py) | Symptom-catalog and section reader for `troubleshooting.md`. |

The Monitoring Operator reference is also exposed at `docs/troubleshooting.md` in the repository root via a symlink.
