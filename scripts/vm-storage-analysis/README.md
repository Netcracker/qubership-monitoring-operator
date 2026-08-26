# VictoriaMetrics Storage Analysis

`vm_storage_report.py` collects a focused storage-optimization report from
VictoriaMetrics and writes it to JSON. It can also generate a small
self-contained HTML report for sharing.

The script is intentionally close in spirit to the logging analysis tooling:

- standalone Python with only the standard library
- `--dry-run` support for inspecting generated PromQL
- optional `.env`-driven execution
- report sections for both raw numbers and quick findings

## What It Collects

- active series from `vm_cache_entries{type="storage/hour_metric_ids"}`
- new series total and churn rate from `vm_new_timeseries_created_total`
- ingestion rate from `vm_rows_inserted_total`
- optional ingestion breakdown by cluster label from `vm_rows_inserted_total`
- total datapoints and bytes-per-sample from `vm_rows` and `vm_data_size_bytes`
- `indexdb` size and `indexdb/data` ratio
- free disk space from `vm_free_disk_space_bytes`
- linear storage full ETA from `vm_free_disk_space_bytes` and
  `rate(vm_data_size_bytes[STORAGE_ETA_LOOKBACK])`
- slow inserts ratio from `vm_slow_row_inserts_total`
- `vm_slow_queries_total` rate per second
- top service groups by visible series footprint from `count(...) by (...)`
- top service groups by new series from `scrape_series_added`
- observed metric/query API requests per second from `vm_request_duration_seconds_count`
- peak average metric/query API request duration from a single ratio expression over `vm_request_duration_seconds_sum / vm_request_duration_seconds_count`
- VictoriaMetrics flags from `flag`
- top high-cardinality metrics
- top queries from `/api/v1/status/top_queries`

The report also adds simple findings based on thresholds for:

- `indexdb/data` ratio
- max per-target slow insert ratio, with a default warning threshold of `0.02`
  because sustained slow inserts above 1-2% can indicate disk or system pressure
  (`SLOW_INSERTS_WARNING_THRESHOLD` or `--slow-inserts-warning-threshold` overrides it)
  The ratio is calculated per target after summing without `type`, so
  `vm_rows_inserted_total{type=...}` rows are combined for each target while
  per-instance degradation is still preserved.

## Usage

Basic example:

```bash
python scripts/vm-storage-analysis/vm_storage_report.py \
  --victoriametrics-url http://victoria-metrics:8428 \
  --selector 'cluster="prod",job=~"vmsingle"' \
  --cardinality-selector 'cluster="prod"' \
  --scrape-selector 'namespace!="monitoring"' \
  --output vm-storage-report.json \
  --html-output vm-storage-report.html
```

If VictoriaMetrics is behind `vmauth`, pass Basic Auth credentials:

```bash
python scripts/vm-storage-analysis/vm_storage_report.py \
  --victoriametrics-url https://vmauth.example.com/select/0/prometheus \
  --vm-user user \
  --vm-pass password \
  --selector 'cluster="prod",job=~"vmsingle"' \
  --output vm-storage-report.json
```

Use `--time-offset` to query a completed point in the past instead of the
latest available snapshot:

```bash
python scripts/vm-storage-analysis/vm_storage_report.py \
  --victoriametrics-url http://victoria-metrics:8428 \
  --selector 'cluster="prod",job=~"vmsingle"' \
  --time-offset 2h \
  --output vm-storage-report-2h-ago.json
```

Use `--dry-run` to inspect the generated PromQL without calling VictoriaMetrics:

```bash
python scripts/vm-storage-analysis/vm_storage_report.py \
  --victoriametrics-url http://victoria-metrics:8428 \
  --selector 'cluster="prod",job=~"vmsingle"' \
  --output dry-run.json \
  --dry-run
```

## Selectors

The script has three selector knobs because the relevant metrics often come
from different components:

- `--selector` for VictoriaMetrics self-monitoring metrics such as `vm_rows`,
  `vm_data_size_bytes`, `vm_cache_entries`
- `--scrape-selector` for `scrape_series_added`
- `--cardinality-selector` for the broad `topk(count by (__name__))` query
- `--service-group-by-labels` for how stored active-series service groups are defined
- `--monitoring-namespace` for container-scoped `vmsingle` / `vmagent` checks
  in the monitoring namespace; this selector is used for CPU throttling, CPU
  pressure, `flag`, `vm_persistentqueue_*`, and `vmagent_remotewrite_*`
  metrics. These checks are not automatically intersected with `--selector`,
  because their label sets often differ from storage self-metrics; in shared
  monitoring namespaces they should be treated as namespace-scoped component
  signals unless you override them with environment-specific queries.
- `--enable-cluster-ingestion-table` enables an optional table with
  `vm_rows_inserted_total` rate grouped by `--cluster-ingestion-label`
- `--query-requests-path-regex` for which VictoriaMetrics API paths count as metric/query traffic
- `--storage-eta-lookback` for the `vm_data_size_bytes` growth window used by
  `summary.storage_full_eta_days`; default is `24h`
- `--metric-label-analysis-limit` controls how many top metrics are inspected via
  `/api/v1/series`; the effective value is never lower than `--top-limit` so
  `High Cardinality Metric Usage` can populate its top-label columns.
- `--series-fetch-workers` controls the parallelism used for `/api/v1/series`
  and for filling missing label-value counts. Default is `8`; lower it on very
  large VictoriaMetrics installations if the API becomes too busy.
- `--metric-usage-limit` controls how many metric names are requested from
  `/api/v1/status/metric_names_stats`; default is `5000`. Increase it if
  request usage columns are empty for high-cardinality metrics.
- `--full-label-scan` scans all discovered metric names for label distribution.
  This is intentionally expensive because metric-name discovery is global via
  `/api/v1/label/__name__/values`, after which the script performs
  selector-scoped `/api/v1/series` calls per metric and prints a warning
  before starting this scan.
- `--label-global-unique-values-threshold` keeps globally present non-infra
  labels in `Label Distribution` when they still have high value cardinality,
  instead of hiding them just because they appear across the full scanned
  metric scope.
- `--max-full-scan-metrics` adds an optional hard cap for FULL_LABEL_SCAN. If
  the discovered metric-name count exceeds this value, the scan is skipped with
  an error instead of continuing into an unbounded fetch.
- `--top-queries-limit` controls how many rows are requested for each
  `/api/v1/status/top_queries` ranking table.
- `--top-queries-lookback` is passed as `maxLifetime` to
  `/api/v1/status/top_queries`; default is `24h`.
- `--vmalert-requests-query` for an optional override of the default vmalert traffic query

Typical setup:

- `--selector 'cluster="prod",job=~"vmsingle"'`
- `--scrape-selector ''`
- `--cardinality-selector 'cluster="prod"'`
- `--service-group-by-labels 'namespace,job'`
- `--query-requests-path-regex '.*/api/v1/(query|query_range|series|labels|label/.+/values|status/tsdb|status/metric_names_stats)$'`

## Environment Variables

All CLI options can also be provided with environment variables. Use the
example file as a template:

```bash
cp scripts/vm-storage-analysis/vm_storage_report.env.example vm-storage.env
set -a
. ./vm-storage.env
set +a
python scripts/vm-storage-analysis/vm_storage_report.py
```

## Output Notes

- `summary` contains the main storage and churn numbers.
- `summary.query_requests_rate_per_second` shows the observed per-second rate of
  VictoriaMetrics metric/query API requests matching the configured or default
  request-path regular expression.
- `summary.storage_full_eta_days` is a linear disk-full projection across
  matched targets. It is empty when recent `vm_data_size_bytes` growth is
  non-positive or unavailable.
- `summary.avg_query_request_duration_seconds` shows the peak average duration for those query-related requests during the lookback window.
- `summary.query_concurrency_limit` is read from `vm_concurrent_select_capacity` and falls back to the effective `search.maxConcurrentRequests` flag value if needed.
- `summary.vmalert_requests_rate_per_second` is a direct datasource request rate when `vmalert_datasource_queries_total` is available; otherwise it uses `vmalert_execution_total` as a rule execution rate proxy, then falls back to `0`, and can be overridden with `VMALERT_REQUESTS_QUERY`.
- `tables.label_distribution` is built only from selector-scoped `/api/v1/series`
  scans. In the default mode this is limited to top TSDB metrics for speed; with
  `FULL_LABEL_SCAN` it expands to all discovered metric names. The table is
  therefore internally scope-consistent, but still heuristic when not every
  metric is scanned. If `FULL_LABEL_SCAN` is aborted by `MAX_FULL_SCAN_METRICS`,
  the report falls back to the already collected top-metrics scan instead of
  dropping the table entirely.
- For global reports without `SELECTOR`, sampled label coverage can be
  extrapolated with TSDB per-metric totals. This extrapolation is intentionally
  disabled for selector-scoped reports, because TSDB metric totals are global
  and would otherwise overmix foreign series into scoped estimates.
- `tables.top_queries_by_*` are read from `/api/v1/status/top_queries` using
  the configured `TOP_QUERIES_LIMIT` and `TOP_QUERIES_LOOKBACK`.
- `tables.vmsingle_configured_flags` and `tables.vmagent_configured_flags` list only explicitly set flags with `Is Set=true`.
- `tables.vmsingle_all_effective_flags` and `tables.vmagent_all_effective_flags` contain the full effective flag sets and are rendered collapsed in the HTML report.
- `tables.top_services_by_series` shows which service groups currently occupy
  the largest visible series footprint in the cardinality instant query. It is
  not based on `vm_cache_entries{type="storage/hour_metric_ids"}` active-series
  accounting.
- `tables.top_services_by_new_series` helps locate churn-heavy service groups.
- `errors` captures failed queries without aborting the whole report.
