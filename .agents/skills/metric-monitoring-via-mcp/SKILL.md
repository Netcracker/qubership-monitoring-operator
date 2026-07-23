---
name: metric-monitoring-via-mcp
description: >-
  Read-only investigation of live Prometheus-compatible metrics through
  connected Grafana and VictoriaMetrics MCP servers. Trigger only when the task
  asks for PromQL or MetricsQL evidence, VictoriaMetrics metric samples or
  labels, Prometheus-compatible scrape targets, metric alert expressions,
  recording rules, or dependencies between recorded numeric time-series. A
  Grafana dashboard or panel qualifies only when the prompt identifies it as
  metric-backed or asks for its Prometheus-compatible metric query.
---

# Metric Monitoring via MCP

Use the connected monitoring MCP servers to investigate the live system. Prefer
read-only evidence from dashboards, rules, metadata, and queries over assumptions.

## Route the Request

- Use `mcp-grafana` for dashboards, folders, panels, dashboard variables,
  datasources, Grafana-managed alerts, and queries made through Grafana
  datasources.
- Use `mcp-victoriametrics` for direct PromQL or MetricsQL queries, metric and
  label discovery, time series, vmalert alerts and rules, cardinality, and query
  diagnostics. Use its vmagent target or service-discovery capabilities when
  they are exposed by the connected server.
- Use only connected monitoring MCP tools. Do not invoke `kubectl`, read the
  Kubernetes API, proxy through a Kubernetes Service, exec into pods, inspect
  local cluster configuration, or read repository manifests to supplement an
  investigation. This boundary applies even when such access is technically
  available.
- All live monitoring evidence and investigation calls must come from MCP.
  Outside loading this skill and its bundled references, do not use shell
  commands, filesystem reads, web requests, or other non-MCP tools to obtain or
  corroborate monitoring facts.
- Do not use shell utilities such as `date` merely to convert MCP-returned
  timestamps. Preserve the exact timestamp and timezone from MCP when possible;
  if a second timezone is useful, convert it in the report without another tool
  call or state that the conversion was not independently verified.
- Treat `ServiceMonitor`, `PodMonitor`, `ScrapeConfig`, `VMServiceScrape`,
  `VMPodScrape`, `VMStaticScrape`, Service, EndpointSlice, and collector
  configuration as unavailable unless a connected MCP tool explicitly returns
  that information. If the MCP exposes only active or dropped targets, report
  only the discovery and relabeling facts present in that response.
- Start with Grafana when the request names a dashboard or panel. Identify the
  actual datasource backend before using a direct backend MCP. Use
  `mcp-victoriametrics` for lower-level evidence only when the datasource or
  deployment context identifies VictoriaMetrics.
- Inspect the tools that are actually available. Tool names and optional
  capabilities vary by MCP server version and configuration.
- If a required MCP server or tool is unavailable, state the missing capability
  and continue with any useful read-only checks that remain.
- When the task intentionally compares access through multiple connections,
  define the facts assigned to each connection before calling either one. If a
  task-specific read is rejected, mark only that connection's facts unavailable;
  do not compensate by expanding the successful connection or replaying the
  rejected query through it unless the user explicitly asks for a fallback.
- For a task that expects one MCP connection to reject access, read the
  `Partial MCP Access or Authorization Failure` playbook before making live
  calls. That closed ledger overrides optional checks in the dashboard, metric,
  and alert playbooks.
- Treat the first successful task-specific read as the access check. Do not add
  a standalone health or `up` query when the required dashboard, rule, or metric
  lookup has already proved access.
- Do not assume where MCP authentication is implemented. A deployment may
  configure credentials in the MCP server or its proxy, inject them into the
  agent connection, or require no authentication. On `401` or `403`, stop
  repeated retries, identify the failing MCP connection, and report that its
  access configuration must be checked at the appropriate layer. Never print or
  persist tokens, passwords, or authorization headers, and do not ask the user
  to paste them into the conversation.

## Follow the Investigation Workflow

1. Identify the target, desired outcome, cluster or tenant, and time window. Use
   an explicit time window and timezone for every range query and state any
   inferred context.
2. Verify access with the first task-specific read. If the task has no natural
   initial lookup, inspect Grafana datasource health or use a narrow
   VictoriaMetrics metadata lookup or simple query such as `up`.
3. Discover the resource instead of guessing identifiers. Search dashboards by
   title or tags, metrics by name or prefix, and alerts or rules by group, name,
   labels, or state.
4. Inspect the source definition. Capture stable identifiers, datasource,
   query expressions, variables, label selectors, scrape selectors, relabeling,
   thresholds, evaluation interval, and `for` duration as applicable.
5. Reproduce the behavior with the smallest useful query. Use an instant query
   for a current value and a range query for trends or incident analysis.
   Preserve the original selectors and evaluate complex expressions in parts.
6. Correlate configuration with returned data. Distinguish observed facts from
   inferences, and do not treat an empty result as proof that a metric never
   exists.
7. Report the resource identifiers, datasource, exact query, time window,
   relevant labels or values, conclusion, limitations, whether any state was
   changed, and the next most useful check.

Read [investigation playbooks](references/investigation-playbooks.md) for
task-specific sequences.

## Control Investigation Cost

- Convert the user's requested facts into a closed evidence checklist before
  calling tools. Adjacent facts are not automatically useful: once every listed
  fact is proved or explicitly unavailable, stop instead of expanding into a
  general health audit.
- Make an evidence plan before querying: definition, current sanity check,
  range validation, and only then a drill-down if the validation fails. Reuse
  data already returned instead of asking for the same rule, labels, or range in
  another form.
- Treat a task-specific evidence ledger in the playbooks as a ceiling, not a
  target. Plan the normal path at least one call below the stated budget so one
  genuine tool failure can be recovered without turning the investigation into
  an exhaustive inventory.
- A failed evidence slot is still consumed when the failure itself answers the
  access question. Do not replace an authorization failure with a proxy query,
  panel execution, variable enumeration, or metric inventory from another
  connection; report the unverified result and stop that branch.
- Treat the first successful response for an evidence slot as consumed. Filter
  or summarize that response locally; do not call the same inventory or target
  endpoint again with a more convenient filter unless the first response omitted
  the required field.
- Assign each fact to one evidence source. For example, reproduce a panel through
  its Grafana datasource and use the direct metrics backend for one compact
  diagnostic summary; do not repeat the same selector, count, or control series
  through both interfaces unless their results conflict.
- For a simple investigation, aim for roughly 8-12 backend calls. This is a
  soft budget, not a reason to stop before the result is supported; if the work
  grows beyond it, narrow the scope or explain which unresolved question needs
  another call.
- Prefer server-side counts and aggregations for equality, missing-series, and
  lag checks. Fetch full matrices only for mismatching or representative series
  because large MCP responses consume context without improving the conclusion.
- When MetricsQL is available, batch independent scalar or aggregate checks
  into one query by adding a distinguishing label with `label_set(...)` and
  joining the results with `or`. Fewer, slightly richer queries are easier to
  audit than many calls that each return one number.
- Match the query step to the rule or panel interval. Do not increase resolution
  beyond the source evaluation cadence, and do not repeat a range query as a
  broad raw selector such as `metric[10m]` unless actual stored timestamps are
  essential and the selector is narrowed to a representative series.
- Read backend documentation only after a tool error or unresolved semantic
  ambiguity blocks the conclusion. Do not make a documentation call merely to
  cite, restate, or confirm behavior already established by successful queries.
- Summarize cardinality, mismatch counts, extrema, and a few representative
  labels by default. Enumerate every series or sample only when the user asks or
  when the outliers themselves are the result.
- Include at least one explicit limitation in the report, even when the diagnosis
  is conclusive, so the reader knows the time window, tenant, or evidence scope
  outside which the conclusion was not tested.

## Inspect Dashboards and Panels

- Find the dashboard before fetching its full definition; prefer UID over a
  mutable title once discovered.
- Identify the panel by title or ID and inspect its datasource, targets,
  transformations, template variables, repeat settings, units, thresholds, and
  time overrides.
- Resolve variables and datasource references before executing a panel query.
- Run the stored panel query when that capability is available. Otherwise run
  the equivalent query through the identified datasource and preserve the
  panel's time range and step.
- Do not infer panel behavior from its title or screenshot when the dashboard
  definition and query results are accessible.

## Inspect Metrics

- Confirm the metric name with metric discovery before composing a broad query.
- Read metadata when available to establish metric type, unit, and description.
- Treat a metric suffix, panel description, or legend as a semantic hint rather
  than authoritative metadata. When unit or scaling is material and metadata is
  absent or conflicts with the display, perform one exact producer-rule lookup.
  For a recorded metric, use the producer expression and health to identify
  scaling, normalization, or unit conversion before blaming visualization. If
  no producer exists, state which remaining evidence supports the assumed unit.
- Discover label names and values before adding selectors. Avoid unbounded
  series enumeration and high-cardinality range queries.
- Compare queries with the same time window, step, tenant, and selectors. Call
  out counter resets, aggregation, missing series, stale data, and unit
  conversion when they affect interpretation.

## Inspect Scrape Discovery and Relabeling

- Inspect only target, service-discovery, relabeling, and scrape-configuration
  evidence returned by `mcp-victoriametrics`. Do not reconstruct missing stages
  by reading Kubernetes resources or calling vmagent through a Kubernetes proxy.
- Trace the target pipeline only as far as MCP evidence permits: discovered or
  dropped target; discovered labels; applied relabeling evidence; final labels;
  scrape URL, health, and last error; stored `up` or requested samples.
- Distinguish these failure classes: the scrape resource was not selected, it
  selected no object, the Service has no ready endpoint, target relabeling
  dropped the target, the target is present but down, metric relabeling dropped
  samples, and remote write or ingestion failed after a successful scrape. Make
  one of these diagnoses only when the MCP evidence reaches that stage; an
  absent active target alone does not identify which earlier stage failed.
- Apply relabel rules in order using the discovered `__meta_*` labels. Do not
  infer the result from a single regex unless the MCP also exposes the relevant
  rule or an explicit dropped-target reason. Account for action defaults,
  missing source labels, separators, and earlier label mutations when available.
- Use `up` only for targets that survived discovery and target relabeling. An
  absent `up` series cannot distinguish an undiscovered target from a dropped
  target or an expired series without configuration and target-state evidence.
- If the connected MCP lacks target or service-discovery capabilities, state
  that the scrape stage cannot be localized through this skill. Do not fall back
  to shell commands, cluster APIs, or configuration files.

## Inspect Alerts and Rules

- First distinguish a Grafana-managed alert from a Prometheus or vmalert rule;
  inspect it through the system that evaluates it.
- Once an exact Grafana rule UID and definition establish Grafana ownership, do
  not query vmalert for a namesake merely to prove absence. Use the other rule
  system only when the first exact lookup is empty or ownership remains
  genuinely ambiguous.
- For an alert, inspect both the current alert instance and its rule definition.
  Capture state, health, expression, labels, annotations, threshold, evaluation
  interval, `for` duration, and no-data or error behavior when available.
- Give every inspected alert rule one compact report row containing its health,
  last error (including an explicitly empty error), expression, interval, `for`,
  labels, and annotations. Preserve every MCP-returned annotation field that
  carries user-authored meaning; when both `summary` and `description` are
  present, report both by name rather than choosing one or replacing them with
  an inferred intent. Keeping definition metadata together prevents a correct
  diagnosis from losing the rule's stated meaning in the final summary.
- Evaluate the alert expression at the incident time and at the current time.
  Break compound expressions into operands to locate the failing assumption.
- When an alert expression references a metric that could be a recording-rule
  output, perform one exact producer lookup. Include the producer expression and
  health when it exists; an alert's historical values alone do not establish
  how that input was generated.
- For a recording rule, find the rule that produces the recorded metric, inspect
  its expression and health, then compare its inputs and output over the same
  time window. Follow dependent recording rules until reaching raw metrics, a
  missing rule, or a cycle. Use the staged, low-volume comparison in the
  recording-rule playbook instead of downloading every series at every level.
- Do not assume that a Grafana alert and a similarly named vmalert rule are the
  same object.
- Do not exec into Grafana pods, copy or inspect Grafana's internal database, or
  create a port-forward merely to recover alert state that the connected
  read-only APIs omit. Use one supported state fallback or report the state as
  unavailable; internal implementation access is disproportionate and brittle.

## Respect Safety and Ownership

- Default to read-only operations. Evaluating an expression is a read-only test
  and is allowed. Do not create, update, delete, or mute resources, change
  evaluation state, or send test notifications unless the user explicitly
  requests that effect and the relevant write tools are enabled.
- Before proposing or applying a change, determine whether the resource is
  user-managed or provisioned from Helm, an operator custom resource, a
  ConfigMap, or a repository file only when MCP-returned provenance establishes
  that ownership. Otherwise report ownership as unknown; do not inspect the
  cluster or repository to fill the gap.
- Keep production queries narrow. Expand their time range or label scope only
  when the narrower query cannot answer the question.
- Never reveal service-account tokens, passwords, authorization headers, or
  secret contents in queries or responses.
- Never claim that a resource was changed when only a recommendation was made or
  when the write operation was unavailable.
