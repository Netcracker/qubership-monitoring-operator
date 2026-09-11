# Monitoring Investigation Playbooks

Use the playbook matching the user's outcome. Skip a step only when its answer
is already established by evidence.

## Contents

- [Dashboard or Panel](#dashboard-or-panel)
- [Partial MCP Access or Authorization Failure](#partial-mcp-access-or-authorization-failure)
- [Metric](#metric)
- [Alert State or Evaluation Failure](#alert-state-or-evaluation-failure)
- [Scrape Discovery or Relabeling](#scrape-discovery-or-relabeling)
- [Recording Rule or Recorded Metric](#recording-rule-or-recorded-metric)
- [Dashboard Shows No Data](#dashboard-shows-no-data)
- [Noisy or Expensive Query](#noisy-or-expensive-query)

## Partial MCP Access or Authorization Failure

When a task intentionally checks multiple MCP connections and one rejects
access, preserve evidence from the working connection but do not widen that
side to compensate for the failure. Before calling tools, assign every requested
outcome to exactly one slot. For a dashboard, datasource-health, and
direct-backend access check, use this closed ledger:

1. Search for the exact dashboard once and record its stable UID.
2. Read one panel or dashboard datasource reference only if the search response
   does not identify the concrete datasource. If panels refer to a datasource
   variable, read the dashboard's complete, small templating list once and
   filter it locally. Do not first try a JSONPath predicate and then fetch the
   same list again when server-side predicate support is uncertain.
3. Read datasource detail only for identity fields still missing, and call its
   health endpoint once when health is requested. If one datasource response
   supplies both identity and health, it fills both slots.
4. Make exactly one task-specific direct-backend query. A `401` or `403` fills
   this slot as an access result even though it supplies no metric value.
5. Stop and report which backend value remains unverified. Do not resolve
   unrelated variables, enumerate labels or metric names, execute a saved or
   arbitrary panel, or replay the rejected expression through Grafana merely to
   compensate for the failure. Use a working-connection fallback only when the
   user explicitly requests one.

Phrases such as "complete every useful check" remain bounded by these named
outcomes, not every adjacent dashboard fact. This path normally takes 3-5
Grafana calls plus the single rejected backend call; a blank response may use
one recovery slot, but an authorization rejection may not.

## Dashboard or Panel

1. Search dashboards by title, tags, or folder.
2. Fetch the selected dashboard by UID.
3. Find the panel by title or ID.
4. Record its datasource UID and type, targets, variables, transformations,
   units, thresholds, and time overrides.
5. Resolve variable values for the requested cluster, namespace, workload, or
   other scope.
6. Execute the stored panel query if supported; otherwise execute its expression
   against the identified datasource.
7. Compare returned values with panel transformations, units, and thresholds.
8. Report whether the issue is in the data, selector or variable resolution,
   datasource, transformation, or visualization configuration.

For a panel whose backend value exists but its displayed magnitude, unit, or
threshold interpretation looks wrong, use this bounded ledger:

1. Search the exact dashboard once, then read the target panel configuration
   through one narrow path. Prefer one panel property that includes target,
   datasource, transformations, field unit, reduction, thresholds, and time
   overrides. Do not also read the dashboard summary, panel index, full
   dashboard, and the same panel through a second path.
2. Resolve the datasource through one detail call unless the panel response
   already supplies its UID, type, and backend. Execute the stored query once in
   its saved instant or range mode; this establishes the value entering panel
   processing.
3. Establish semantic type and unit with at most one metadata lookup. A suffix
   such as `_seconds`, dashboard prose, or a legend is a useful hint but not a
   substitute for source evidence when unit or scaling is the disputed fact.
4. If metadata is absent or conflicts and the queried name may be a recorded
   metric, make one exact producer-rule lookup. Record its expression and
   health, and check whether it scales, normalizes, aggregates, or converts the
   input. Do not replace this lookup with broad metric discovery or multiple
   alternative value queries.
5. Compare the one returned value with transformations, reduction, unit,
   decimals, thresholds, and overrides locally. Add one direct-backend range
   only when temporal stability or a Grafana/backend disagreement is material;
   do not run instant, range, `last_over_time`, and a second interface merely to
   prove the same value.
6. Stop when source semantics, query value, and panel processing identify the
   first conversion or presentation mismatch. Report any remaining semantic
   uncertainty rather than inferring authority from the metric name alone.

This normally takes 5-8 monitoring calls: dashboard discovery, one panel read,
one datasource read, one stored-query execution, metadata and—only when that
metadata is insufficient—one producer lookup, plus at most one justified
backend comparison. Leave one recovery slot instead of filling every optional
branch.

## Metric

1. Search metric names with the narrowest known prefix or pattern.
2. Read metric metadata to determine type, unit, and meaning.
3. Inspect label names, then values for only the labels needed to scope the
   query.
4. Run a narrow instant query to confirm current availability.
5. Run a range query for the requested interval using an appropriate step.
6. Compare related series using identical selectors and time settings.
7. Explain the result, including aggregation, rate or increase semantics,
   counter resets, gaps, staleness, and units when relevant.

For a bounded temporal comparison of a few known metrics, use this evidence
ledger:

1. Fetch all known recording-rule definitions in one filtered rule call when
   the tool accepts multiple names. The definitions already establish
   expressions, interval, health, and rule-added labels; do not repeat one call
   per metric.
2. Read metric metadata only when type, unit, or meaning remains material and
   cannot be established from the rule or user context. Do not make identical
   metadata calls for every generated metric merely to confirm that metadata is
   absent.
3. Freeze one range and use one combined range query for raw values plus the
   decisive reset, rate, increase, sample-count, zero, or stale-marker checks.
   Add distinguishing labels so multiple metrics and checks can share the call.
4. Use at most one compact instant query when the range response does not
   establish current presence or freshness. Do not first request range vectors
   through an instant-query endpoint and then repeat them with a range query.
5. Stop when rule health and expression explain the observed pattern and the
   stored range agrees. Remote-write health, a second expression comparison, or
   another raw matrix is justified only by conflicting or missing evidence.

This normally takes 3-6 backend calls. Documentation is a recovery path for an
unresolved function or response format, not a routine confirmation of reset or
staleness semantics already demonstrated by query results.

## Alert State or Evaluation Failure

1. Determine whether Grafana Alerting or vmalert evaluates the alert.
2. Fetch the rule definition and check for an active alert instance. An absent
   instance does not by itself distinguish a false condition, valid empty
   result, evaluation error, or expired alert.
3. Record state, health, last error, labels, annotations, expression, threshold,
   evaluation interval, `for` duration, and error or no-data behavior. Give each
   inspected rule one compact definition row in the report. Include an empty
   last error explicitly, and preserve every MCP-returned annotation field that
   carries user-authored meaning. When both `summary` and `description` exist,
   label and include both; do not select one or replace either with a diagnosis
   or inferred intent. This keeps severity, scope, and stated meaning attached
   to the correct rule without dumping unrelated metadata.
4. Freeze one explicit time and evaluate the full expression there. Classify
   the outcome before drilling down:
   - unhealthy rule or query error: evaluation failed and produced no alert
     vector; inspect the error and only the operands needed to isolate it;
   - healthy rule and successful empty vector: valid no-result evaluation;
     inspect the narrow selector or aggregation that removed the input;
   - nonempty vector: compare its values with the threshold and `for` duration.
5. Evaluate meaningful operands or intermediate aggregations separately. Use a
   combined labelled count or summary when it can prove both input availability
   and selector mismatch without downloading full vectors.
6. Check whether missing data, changed labels, delayed ingestion, vector
   matching, or datasource errors explain the state. Do not label a healthy
   empty vector as an evaluator failure, and do not label an unhealthy rule as
   ordinary No Data or a false threshold.
7. Recommend the smallest correction for the established cause. Do not execute
   alternative corrected expressions merely to demonstrate that they parse
   unless the user asks to verify a proposed correction or evidence remains
   ambiguous.
8. Report evidence for why the alert is firing, pending, inactive, absent, or
   unknown, including the point-in-time or historical limitation.

For a comparison of two known vmalert rules, normally use one filtered rule
lookup, one full-expression query per rule, up to two compact input summaries,
and one active-instance lookup. Reuse rule fields already returned. A failed
paginated alert lookup may be retried once without pagination, but do not try
multiple limits. Exceed this 5-6 call ledger only for a concrete tool error,
conflicting evidence, or required historical investigation.

For a reported incident that has already recovered, freeze the reported
incident window before checking the current state. Query the full expression
and decisive operands over that historical window, then use one current check
of the exact complete alert expression to establish recovery. A raw input value
or empty current-instance list does not replace that expression check. If the
alert expression references a metric that could be recorded, use one exact
producer-rule lookup and report its expression and health when found; otherwise
the investigation proves the input's values but not its generation semantics.
A missing current alert instance is evidence only about the current state; it
is not proof that the alert never fired. If the alerting API exposes no retained
state history, say so and use rule definitions plus historical metric values
without inventing `activeAt`, notification, or continuity evidence.

Use a closed recovered-incident ledger: one alert-rule lookup; one exact
producer lookup when the input may be recorded; one historical range combining
the full expression and decisive input; one exact current full-expression
query; and one current active-instance lookup. Add at most one historical alert
state range when it materially distinguishes inferred condition truth from an
observed pending or firing state. Reuse the original rule definition and current
check even if wall-clock time advances during the investigation; report a newly
started incident or changed snapshot instead of repeating the rule lookup or
the same full-expression query at successive times.

For one named Grafana-managed alert, use this evidence ledger:

1. Perform one exact or name-filtered Grafana alert-rule lookup. If it returns a
   matching UID and definition, ownership is established; do not add a vmalert
   lookup. If the response is `null`, blank, or lacks the named rule, make at
   most one unfiltered Grafana rule-list call because some server versions do
   not implement filters consistently. Do not try separate firing, pending,
   all-state, and limit variants.
2. Reuse the list response when it already contains condition, queries,
   interval, `for`, no-data/error behavior, labels, annotations, and paused
   state. Otherwise fetch that UID once. Do not list all rules or all alert
   states after the exact match.
3. Use current instances or state included by the Grafana alerting tool. When
   state fields are absent, choose exactly one supported fallback:
   - one Grafana alert-instance or Alertmanager read through the connected MCP;
   - or one compact query of Grafana alerting metrics that can be tied to the
     rule UID or name.
   Do not combine global metric inventory, individual state queries, a range
   query, and an internal state-store read. If no supported fallback identifies
   the instance, report current state as unavailable while still explaining the
   rule definition.
4. Use provenance or synchronization fields only when a connected Grafana MCP
   response includes them. Do not inspect `GrafanaAlertRuleGroup` resources,
   Kubernetes CRDs, pods, internal databases, repository manifests, or direct
   Grafana endpoints outside MCP. If Grafana lookups do not identify the rule,
   one exact vmalert lookup may test the remaining evaluator hypothesis; if it
   is also empty, report ownership as unresolved.
5. For `__expr__` queries, evaluate the stored math or reduce condition locally
   when its inputs are already literal or present in the rule response. Do not
   query VictoriaMetrics for an expression that has no VictoriaMetrics input.
6. Stop after definition and one current-state path agree. A normal path is 3-6
   MCP calls; reserve a seventh for one failed or blank state response. When
   state or provenance remains unavailable after the one MCP fallback, say so—
   the evidence ledger is complete even though the desired fact is unavailable.

## Scrape Discovery or Relabeling

1. Inspect the target and service-discovery capabilities actually exposed by
   `mcp-victoriametrics`. If neither is available, state that this skill cannot
   localize the missing target and stop the discovery branch. Do not use
   `kubectl`, Kubernetes APIs, service proxies, pod exec, or local manifests.
2. Request active and dropped targets once, scoped by the narrowest stable
   identity the MCP accepts. Prefer scrape-pool/job identity or discovered
   resource labels over a mutable final `job` label. Reuse this response rather
   than retrying with progressively looser filters.
3. Separate discovered labels from final labels. Record any MCP-returned scrape
   URL, health, last scrape, last error, dropped reason, and relabeling or scrape
   configuration. Do not claim that a collector selected a scrape resource, a
   Service matched, or an endpoint was ready unless the MCP response explicitly
   proves that fact.
4. Replay relabeling only when the MCP exposes both the relevant ordered rules
   and their source labels, or an explicit dropped-target explanation. Account
   for defaults, separators, missing labels, and earlier mutations. Otherwise
   report that relabel causality is unavailable rather than inferring it from a
   suggestive label value.
5. For an active target, use its MCP-reported health and error to distinguish a
   scrape failure from discovery loss. Add at most one narrow `up` query for the
   same final labels when stored presence is material. If `up` exists but a
   requested metric does not, diagnose metric relabeling only when MCP-returned
   configuration proves it; otherwise keep sample filtering and endpoint output
   as unresolved alternatives.
6. Check remote-write or ingestion evidence only after MCP data proves a
   successful scrape retained the relevant sample. Stop at the first stage the
   MCP evidence actually establishes and state every earlier configuration
   stage that remains outside visibility.

This normally takes 2-5 MCP calls: one capability or target-discovery read, one
active-and-dropped target snapshot, an optional returned-config read, one narrow
metric query, and at most one justified ingestion check. A missing MCP target
capability is a valid limitation, not permission to switch tools.

## Recording Rule or Recorded Metric

1. Search for the producer by exact output metric name. Use a filtered rule
   lookup when available; scan the full inventory only if filtering is
   unavailable or multiple producers must be disambiguated.
2. Record the group, expression, rule-added labels, evaluation interval, health,
   last error, and last evaluation. The successful rule read also proves access,
   so do not add a separate connectivity query.
3. Parse metric references from the expression. Search each candidate by exact
   rule output name, track visited names, and build the dependency chain before
   running range queries. Perform this exact producer lookup for the terminal
   candidate too: only an empty producer result establishes that it is raw or
   has a missing rule. Do not classify a familiar metric such as `up` as raw by
   convention. Stop after the verified empty lookup or at a cycle.
4. Confirm current availability with one combined instant summary when the
   backend supports MetricsQL. Label and join per-level counts, sums, minima, or
   maxima in one query instead of fetching each full vector separately. Query a
   representative series only if labels or values need drill-down. Choose one
   composition form (`or` with distinguishing labels, or `union(...)`) before
   calling the backend. If the response is nonempty and contains every expected
   `level` or `source` label, accept it: do not retry with the other composition
   form or query the same stored and expression vectors individually.
   Before range comparison, classify each branch as populated, one-sided, or
   empty. When both a stored output and its expression are proven empty, their
   equality is vacuous: report the empty branch and its dependency defect, then
   skip mismatch matrices, freshness checks, and healthy-chain validation for
   that branch. When only one side exists, retain the missing-series comparison
   because the asymmetry is evidence.
5. Freeze one explicit range and use the rule interval as the step for every
   comparison. Compare stored output with its expression using server-side
   mismatch counts where vector matching permits it:
   - count unequal values after ignoring only rule-added labels;
   - count stored series absent from the expression;
   - count expression series absent from the stored output.
   Combine the three checks for one level into a single MetricsQL range query.
   `or vector(0)` makes a successful zero-mismatch result explicit, while the
   `check` label keeps the three results distinct:

   ```promql
   label_set(
     (count(<stored> != ignoring(<rule_added_labels>) (<expression>)) or vector(0)),
     "check", "unequal"
   )
   or label_set(
     (count(<stored> unless ignoring(<rule_added_labels>) (<expression>)) or vector(0)),
     "check", "stored_missing"
   )
   or label_set(
     (count((<expression>) unless ignoring(<rule_added_labels>) <stored>) or vector(0)),
     "check", "expression_missing"
   )
   ```

   Omit `ignoring(...)` when the rule adds no labels. If either side is a scalar
   or vector matching is otherwise invalid, use equivalent server-side
   aggregations or narrow both sides to the same representative label set.
   Use one combined range call per recorded level; split it only when the
   backend rejects the combined syntax or a nonzero check requires drill-down.
6. If all counts are zero, record the counts and avoid downloading full output
   and expression matrices. If matching is ambiguous or a count is nonzero,
   fetch only the affected label set or one representative series before
   widening the query.
7. Check freshness with one combined instant age query for all levels, or a
   server-side `tlast_over_time` aggregation when actual sample time is
   available. A range query's timestamps may be evaluation grid timestamps
   rather than raw storage timestamps, so do not infer stored sample lag from
   them alone.
8. Repeat the same compact validation for every recorded input. Query raw input
   series only when needed to explain an aggregation, gap, or mismatch; prefer
   counts and grouped summaries over an unfiltered matrix.
9. Report the dependency chain first, then a compact table of rule health and
   per-level mismatch, missing-series, label, and lag results. Include exact
   queries for reproducibility, but summarize samples instead of listing them.

For a mature two-rule chain, the normal target is about 7-10 backend calls:
three producer lookups, one combined instant summary, one combined range
comparison per recorded level, and one combined freshness check. Exceed this
when a nonzero check, rejected combined query, or genuine ambiguity requires
drill-down, not merely to restate already established facts.

Use this evidence ledger for a healthy mature two-rule chain; each slot may be
filled at most once:

1. Three exact producer lookups, including the terminal candidate.
2. Two combined mismatch range queries, one per recorded level.
3. At most two endpoint summaries: one for raw/expression labels and values and
   one for stored recording values. Prefer one combined aggregate summary when
   it answers both.
4. One combined freshness query for every level.
5. One combined range-statistics query only when the endpoint summaries do not
   establish the requested values across the window.

This gives 7-9 calls in the normal case. After a successful zero-mismatch result,
do not add alternative union syntax, individual stored or expression queries,
change-count probes, or another freshness form. Those calls are justified only
by a rejected query, a missing expected source label, or a nonzero mismatch.

## Dashboard Shows No Data

1. Search once with the exact supplied dashboard title. After a match, use its
   UID and do not try title variants unless the result is genuinely ambiguous.
2. Read the dashboard summary, target panel query, templating values, and time
   settings with the narrowest available calls. Treat these reads as a fixed
   evidence ledger: one summary read, one panel-query read, one templating read,
   and a time-property read only when the summary or query response lacks the
   window. Reuse fields already returned. When the panel-query tool identifies
   the panel and target, do not also read `$.panels[*]`, `$.panels[0]`, or the
   full dashboard.
3. Resolve the panel's datasource variable and identify its UID, type, and
   backend before using a direct backend MCP. Choose one identity path: reuse a
   datasource-list result when it contains the needed name, UID, type, and URL,
   otherwise use one datasource-detail call. Do not call both merely to restate
   identity, and never reread datasource details after executing the panel.
   Check health once.
4. Substitute the saved variable values without changing them and execute the
   effective panel expression as one range query over the dashboard window. For
   simple string variables, substitute locally from the templating read; do not
   call the panel-query tool again only to resolve them. If Grafana-side variable
   interpolation is necessary, make the resolved panel-query call the single
   panel-query read for the investigation. The final range point also establishes
   current absence, so do not repeat it as an instant query unless the range
   response omits the endpoint or is ambiguous.
5. Diagnose metric and selector availability with one compact backend query when
   MetricsQL is available. Combine a zero-explicit count for the selected
   expression with counts grouped by the suspected label, for example:

   ```promql
   label_set(
     (count(<resolved_panel_expression>) or vector(0)),
     "check", "selected_series"
   )
   or label_set(
     count by (<suspected_label>) (<referenced_metric>),
     "check", "available_label_value"
   )
   ```

   This single result establishes whether the metric exists, whether the saved
   selector matches, and which label values are present. If the backend lacks
   this syntax, use at most one compact metric count and one scoped label-values
   lookup. Do not additionally fetch the raw metric vector or perform metric-name
   discovery after these results establish existence.
6. Stop when the datasource is healthy, the resolved panel range is empty, the
   referenced metric exists, and the selected label value is absent. These facts
   are sufficient to attribute No data to selector or variable resolution; a
   second reproduction through another interface and a control range for a
   known-good label merely restate the conclusion.
7. Drill down only if those facts conflict or remain ambiguous. Remove selectors
   one at a time, preserving the same time window, and query one representative
   series rather than widening to an unfiltered matrix.
8. Report the dashboard and panel identifiers, datasource, stored and resolved
   expressions, time window, decisive evidence, root-cause category, smallest
   read-only correction, and at least one explicit limitation on the scope of
   the conclusion.

For a straightforward invalid-selector case, target 8-10 monitoring MCP calls:
dashboard discovery and targeted properties, one datasource identity/health
check, one panel range reproduction, and one compact backend diagnosis. Exceed
this only for a tool error, ambiguous discovery, or conflicting evidence.

## Noisy or Expensive Query

1. Capture the exact expression, caller or panel, time range, step, and scope.
2. Choose one bounded diagnostic appropriate to the question: active queries,
   top queries, metric usage, or TSDB cardinality.
3. Use the stored expression and that diagnostic to identify which selector,
   grouping, or join expands the work. Evaluate a safe subexpression only when
   those facts leave the cause ambiguous.
4. Preserve semantics when proposing narrower selectors, earlier aggregation,
   or a recording rule.
5. Recommend a recording rule only when MCP evidence shows that the corrected
   aggregation is reused or remains expensive; otherwise prefer the narrower
   query itself.

For a dashboard-backed expensive-query investigation, use this fixed normal
ledger and stop when it is filled:

1. Grafana, at most six calls: exact dashboard search; one panel-query read; at
   most one numeric-ID or index-based panel property for execution settings;
   one templating read; one time-property read; and one datasource-detail read.
   Skip any slot whose fact is already returned. Do not add dashboard summary,
   full-dashboard, title-filter property, datasource-list, or datasource-health
   calls to the normal path.
2. Backend diagnostic, exactly one family: choose recent top queries, current
   active queries, one metric-usage view, or one current TSDB-cardinality view.
   The word "or" is exclusive here. Do not combine these families or retry the
   same family with alternative dates, matches, focus labels, or limits after a
   successful response. An absent query in bounded retention is not proof that
   it never ran.
3. Candidate validation, at most one call: execute one safe narrow replacement.
   For a current-state intent, omit `time` when the tool defaults to now. Do not
   test it at midnight, the saved range start, or another control time, and do
   not run multiple replacement variants.
4. Estimate range cost from the stored selector, grouping, window, step, and the
   single bounded diagnostic. Stop after these facts explain the risk and the
   candidate answers the stated intent.

This is an eight-call normal path with one spare recovery slot. Do not call
documentation, metric metadata, datasource health, a query explainer, or the
unsafe original unless a concrete error or unresolved backend-specific semantic
question makes that exact call necessary. If a successful result is merely less
detailed than hoped, report the limitation instead of widening the inventory.
