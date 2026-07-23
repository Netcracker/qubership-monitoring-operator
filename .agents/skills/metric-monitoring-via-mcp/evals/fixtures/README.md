# Monitoring eval fixtures

These manifests create the deterministic live resources required by the
`metric-monitoring-via-mcp` task evals. They contain no credentials and are not input
files for the agent under evaluation.

## Mapping

| Eval | Fixture | Expected live state |
| --- | --- | --- |
| 1 — recording-rule chain | `recording-rule-chain.yaml` | Two healthy recording rules produce `skill_eval:up:sum_by_job` and `skill_eval:up:double_by_job`. |
| 3 — dashboard No Data | `dashboard-no-data.yaml` | The panel selects the deliberately absent job value `skill-eval-nonexistent`. |
| 4 — alert investigation | `alert-rule.yaml` | `SkillEvalApiServerReplicaShortage` becomes pending and then firing when fewer than two healthy `kube-apiserver` targets are visible. |
| 7 — alert error versus No Data | `alert-error-no-data.yaml` | One alert rule has a runtime vector-matching error; the other is healthy but has no matching input series. |
| 8 — counter reset and gaps | `metric-temporal-semantics.yaml` | One synthetic counter resets every minute; another is absent during a repeatable part of each 30-second cycle. |
| 9 — expensive query | `dashboard-expensive-query.yaml` | A panel stores an intentionally broad all-metrics query over 24 hours with a one-second interval. |
| 10 — display mismatch | `dashboard-display-mismatch.yaml` | The backend stores 1.5 seconds while the Grafana panel formats the value as milliseconds. |
| 11 — recovered incident | `historical-recovered-alert.yaml` | A periodic signal is true for 90 seconds and then recovers, leaving historical metric evidence after the active instance disappears. |
| 12 — broken rule dependencies | `recording-rule-broken-chain.yaml` | One chain terminates at a missing producer and another contains a two-rule cycle. |
| 13 — Grafana-managed alert | `grafana-managed-alert.yaml` | A Grafana expression rule evaluates to one and becomes active independently of vmalert. |

Eval 5 was retired because logs and traces are outside the skill scope. Evals 2,
6, and 14 do not use Kubernetes fixtures. Eval 2 requires the deployment's
normal Grafana and VictoriaMetrics MCP access. Eval 6 must be launched with the
MCP authentication environment variables intentionally absent so that the
connections return `401` or `403`. Eval 14 uses normal Grafana access but must
launch only the VictoriaMetrics MCP without its required authentication
environment. Eval 15 was retired because its required ServiceMonitor, Service,
EndpointSlice, and vmagent target evidence was available only through Kubernetes
access in this deployment, which is outside the MCP-only skill scope. Never put
credentials in this directory.

## Prerequisites

- Namespace `monitoring` exists.
- Prometheus Operator `monitoring.coreos.com/v1` `PrometheusRule` CRD exists.
- Grafana Operator `grafana.integreatly.org/v1beta1` `GrafanaDashboard` CRD exists.
- A Grafana instance matches labels `app.kubernetes.io/component=grafana` and
  `app.kubernetes.io/part-of=monitoring`.
- A datasource named `Platform Monitoring Prometheus` exists.
- vmalert selects `PrometheusRule` resources in the `monitoring` namespace and
  evaluates them against a backend containing the raw `up` metric.
- Grafana Unified Alerting and `GrafanaAlertRuleGroup` synchronization are
  enabled for eval 13.

The fixture harness may use Kubernetes manifests to create deterministic live
state, but the agent under evaluation must receive evidence only through the
connected Grafana and VictoriaMetrics MCP servers. Do not grant it `kubectl`,
Kubernetes API, pod exec, service-proxy, or repository-fixture access.

## Apply

Apply only the fixtures needed by the selected eval, for example:

```bash
kubectl apply -f .agents/skills/metric-monitoring-via-mcp/evals/fixtures/recording-rule-chain.yaml
kubectl apply -f .agents/skills/metric-monitoring-via-mcp/evals/fixtures/dashboard-no-data.yaml
kubectl apply -f .agents/skills/metric-monitoring-via-mcp/evals/fixtures/alert-rule.yaml
kubectl apply -f .agents/skills/metric-monitoring-via-mcp/evals/fixtures/alert-error-no-data.yaml
kubectl apply -f .agents/skills/metric-monitoring-via-mcp/evals/fixtures/metric-temporal-semantics.yaml
kubectl apply -f .agents/skills/metric-monitoring-via-mcp/evals/fixtures/dashboard-expensive-query.yaml
kubectl apply -f .agents/skills/metric-monitoring-via-mcp/evals/fixtures/dashboard-display-mismatch.yaml
kubectl apply -f .agents/skills/metric-monitoring-via-mcp/evals/fixtures/historical-recovered-alert.yaml
kubectl apply -f .agents/skills/metric-monitoring-via-mcp/evals/fixtures/recording-rule-broken-chain.yaml
kubectl apply -f .agents/skills/metric-monitoring-via-mcp/evals/fixtures/grafana-managed-alert.yaml
```

Do not start the eval immediately. Poll the live systems until:

- both recorded metrics return data and both producer rules report healthy;
- each GrafanaDashboard reports `DashboardSynchronized=True` and can be found by
  its expected title through mcp-grafana;
- the alert instance is visible through mcp-victoriametrics and has reached the
  state expected by the eval. With a 30-second evaluation interval and `for: 1m`,
  this normally requires more than one minute.
- for eval 7, both rules are visible in group
  `monitoring-skill-alert-failure-eval`, `SkillEvalEvaluationError` has unhealthy
  rule health and a nonempty vector-matching error, while
  `SkillEvalNoInputData` is healthy with zero output samples and no active alert
  instance.
- for eval 8, both recording rules are healthy and have evaluated for at least
  two minutes. Over a three-minute range, `skill_eval_sawtooth_counter_total`
  has at least one reset and `skill_eval_intermittent_counter_total` has missing
  evaluation slots.
- for eval 9, the GrafanaDashboard reports `DashboardSynchronized=True`, its
  title is discoverable through mcp-grafana, and nobody has opened the panel or
  otherwise executed its intentionally broad stored query before the eval.
- for eval 10, both the recording rule and dashboard are synchronized and the
  metric returns `1.5` at the current time.
- for eval 11, the rules have evaluated for more than five minutes. Start the
  agent only while `skill_eval_periodic_incident` currently equals zero and the
  previous five-minute window contains a continuous true segment of at least
  20 seconds; record the exact launch time for grading.
- for eval 12, all three rule definitions are visible and none of the three
  recorded outputs currently returns samples.
- for eval 13, `GrafanaAlertRuleGroup` reports successful synchronization and
  Grafana Alerting exposes rule UID `skill-eval-grafana-alert` as pending or
  firing. If the deployment does not expose Grafana alert reads through MCP,
  mark the fixture unsupported instead of grading the agent on invented data.

The alert threshold assumes this eval environment exposes fewer than two
healthy `kube-apiserver` `up` series. Verify that assumption before grading; if
the environment is highly available, use a separate isolated fixture namespace
or adjust the fixture and matching expected output together.

## Cleanup

Delete the same manifests after the eval:

```bash
kubectl delete --ignore-not-found -f .agents/skills/metric-monitoring-via-mcp/evals/fixtures/recording-rule-chain.yaml
kubectl delete --ignore-not-found -f .agents/skills/metric-monitoring-via-mcp/evals/fixtures/dashboard-no-data.yaml
kubectl delete --ignore-not-found -f .agents/skills/metric-monitoring-via-mcp/evals/fixtures/alert-rule.yaml
kubectl delete --ignore-not-found -f .agents/skills/metric-monitoring-via-mcp/evals/fixtures/alert-error-no-data.yaml
kubectl delete --ignore-not-found -f .agents/skills/metric-monitoring-via-mcp/evals/fixtures/metric-temporal-semantics.yaml
kubectl delete --ignore-not-found -f .agents/skills/metric-monitoring-via-mcp/evals/fixtures/dashboard-expensive-query.yaml
kubectl delete --ignore-not-found -f .agents/skills/metric-monitoring-via-mcp/evals/fixtures/dashboard-display-mismatch.yaml
kubectl delete --ignore-not-found -f .agents/skills/metric-monitoring-via-mcp/evals/fixtures/historical-recovered-alert.yaml
kubectl delete --ignore-not-found -f .agents/skills/metric-monitoring-via-mcp/evals/fixtures/recording-rule-broken-chain.yaml
kubectl delete --ignore-not-found -f .agents/skills/metric-monitoring-via-mcp/evals/fixtures/grafana-managed-alert.yaml
```

The shared label key `skill-eval` and the purpose annotation make these resources
easy to identify, but cleanup uses exact manifest names to avoid deleting other
tests.
