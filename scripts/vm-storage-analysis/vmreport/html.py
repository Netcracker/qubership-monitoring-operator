from __future__ import annotations

import json
import math
import re
from pathlib import Path
from typing import Any
from xml.sax.saxutils import escape


def html_escape(value: Any) -> str:
    return escape(str(value), {'"': "&quot;", "'": "&#x27;"})


def format_number(value: Any) -> str:
    if value is None:
        return ""
    if isinstance(value, bool):
        return str(value)
    if isinstance(value, int):
        return f"{value:,}".replace(",", " ")
    if isinstance(value, float):
        if math.isfinite(value) and value.is_integer():
            return f"{int(value):,}".replace(",", " ")
        return f"{value:,.3f}".rstrip("0").rstrip(".").replace(",", " ")
    return html_escape(value)


def format_percentage(value: float) -> str:
    percent = value * 100
    if percent == 0:
        return "0%"
    abs_percent = abs(percent)
    if abs_percent < 0.001:
        digits = 6
    elif abs_percent < 0.01:
        digits = 4
    elif abs_percent < 1:
        digits = 3
    else:
        digits = 2
    formatted = f"{percent:.{digits}f}".rstrip("0").rstrip(".")
    return f"{formatted}%"


def format_precise_rate(value: float) -> str:
    if value == 0:
        return "0"
    abs_value = abs(value)
    if abs_value < 0.000001:
        digits = 9
    elif abs_value < 0.001:
        digits = 6
    elif abs_value < 1:
        digits = 4
    else:
        digits = 3
    return f"{value:,.{digits}f}".rstrip("0").rstrip(".").replace(",", " ")


def format_summary_value(metric_name: str, value: Any) -> str:
    if metric_name in {
        "slow_inserts_ratio",
        "remote_write_http_error_ratio_max",
        "vmagent_persistentqueue_bytes_dropped_ratio",
    } and isinstance(value, (int, float)):
        return html_escape(format_percentage(float(value)))
    if metric_name in {
        "container_cpu_cfs_throttled_seconds_rate_max",
        "process_pressure_cpu_stalled_seconds_rate_max",
        "vmagent_container_cpu_cfs_throttled_seconds_rate_max",
        "vmagent_process_pressure_cpu_stalled_seconds_rate_max",
    } and isinstance(value, (int, float)):
        return html_escape(format_precise_rate(float(value)))
    if metric_name.endswith("_rate_per_second") and isinstance(value, (int, float)):
        return html_escape(format_precise_rate(float(value)))
    return format_cell(value)


def summary_metric_labels() -> dict[str, str]:
    return {
        "active_series": "Active Series",
        "total_datapoints": "Total Datapoints",
        "data_size_gb": "Data Size GB",
        "indexdb_size_gb": "IndexDB Size GB",
        "bytes_per_sample": "Bytes Per Sample",
        "index_to_data_ratio": "Index To Data Ratio",
        "min_free_disk_space_gb": "Min Free Disk Space GB",
        "total_free_disk_space_gb": "Total Free Disk Space GB",
        "storage_full_eta_days": "Storage Full ETA Days",
        "churn_rate_per_second": "Churn Rate Per Second",
        "new_series_total": "New Series Total",
        "new_series_per_active_series": "Lookback New Series / Active Series",
        "ingestion_rate_per_second": "Ingestion Rate Per Second",
        "remote_write_requests_max": "Remote Write Requests Max",
        "remote_write_http_errors_max": "Remote Write HTTP Errors Max",
        "remote_write_http_error_ratio_max": "Remote Write HTTP Error Ratio Max",
        "remote_write_parser_read_errors_max": "Remote Write Parser Read Errors Max",
        "remote_write_parser_unmarshal_errors_max": "Remote Write Parser Unmarshal Errors Max",
        "rows_ignored_too_many_labels_total": "Rows Ignored Too Many Labels",
        "rows_ignored_too_long_label_name_total": "Rows Ignored Too Long Label Name",
        "rows_ignored_too_long_label_value_total": "Rows Ignored Too Long Label Value",
        "insert_limit_reached_total": "Insert Limit Reached",
        "select_limit_reached_total": "Select Limit Reached",
        "select_limit_timeout_total": "Select Limit Timeout",
        "query_requests_rate_per_second": "Query Requests Per Second",
        "avg_query_request_duration_seconds": "Avg Query Request Duration Seconds",
        "query_concurrency_limit": "Query Concurrency Limit",
        "vmalert_requests_rate_per_second": "VMAlert Requests Per Second",
        "slow_inserts_ratio": "Slow Inserts Ratio",
        "container_cpu_cfs_throttled_seconds_rate_max": "VMSingle CPU Throttling Max",
        "process_pressure_cpu_stalled_seconds_rate_max": "VMSingle CPU Pressure Max",
        "vmagent_container_cpu_cfs_throttled_seconds_rate_max": "VMAgent CPU Throttling Max",
        "vmagent_process_pressure_cpu_stalled_seconds_rate_max": "VMAgent CPU Pressure Max",
        "vmagent_persistentqueue_bytes_dropped_ratio": "VMAgent Queue Drop Ratio",
        "remote_write_traffic_mbit_per_second": "Remote Write Traffic Mbit Per Second",
        "vm_slow_queries_total_rate_per_second": "Slow Queries Rate Per Second",
        "storage_growth_rows_per_second": "Storage Growth Rows Per Second",
    }


def format_cell(value: Any) -> str:
    if isinstance(value, dict):
        return "<code>" + html_escape(json.dumps(value, ensure_ascii=True)) + "</code>"
    if isinstance(value, list):
        return "<code>" + html_escape(json.dumps(value, ensure_ascii=True)) + "</code>"
    return html_escape(format_number(value))


def format_table_cell(table_name: str, header: str, value: Any) -> str:
    if table_name == "label_distribution" and header in {"Series Share", "Metrics Coverage"} and isinstance(value, (int, float)):
        return html_escape(f"{value:.1%}")
    return format_cell(value)


def table_cell_class(table_name: str, header: str, value: Any) -> str:
    classes = [f"col-{css_class_name(header)}"]
    if table_name == "label_distribution" and header == "Classification" and isinstance(value, str):
        normalized = value.strip().lower()
        if normalized == "observe":
            classes.append("classification-observe")
        elif normalized == "review":
            classes.append("classification-review")
        elif normalized == "normal":
            classes.append("classification-normal")
    return " ".join(classes)


def threshold_status(value: Any, rule: dict[str, Any] | None) -> tuple[str, str]:
    if not rule or not isinstance(value, (int, float)):
        return ("n/a", "status-na")
    threshold = rule.get("value")
    kind = rule.get("kind")
    if not isinstance(threshold, (int, float)):
        return ("n/a", "status-na")
    if kind == "max":
        return ("alert", "status-alert") if value > threshold else ("ok", "status-ok")
    if kind == "min":
        return ("alert", "status-alert") if value < threshold else ("ok", "status-ok")
    return ("n/a", "status-na")


def threshold_text(rule: dict[str, Any] | None, *, metric_name: str = "") -> str:
    if not rule:
        return ""
    threshold = rule.get("value")
    kind = rule.get("kind")
    if not isinstance(threshold, (int, float)) or kind not in {"max", "min"}:
        return ""
    sign = ">" if kind == "max" else "<"
    if metric_name == "slow_inserts_ratio":
        return f"{sign} {format_percentage(float(threshold))}"
    if metric_name in {
        "remote_write_http_error_ratio_max",
    }:
        return f"{sign} {format_percentage(float(threshold))}"
    if metric_name == "vmagent_persistentqueue_bytes_dropped_ratio":
        return "> 0"
    if metric_name in {
        "remote_write_http_errors_max",
        "remote_write_parser_read_errors_max",
        "remote_write_parser_unmarshal_errors_max",
        "rows_ignored_too_many_labels_total",
        "rows_ignored_too_long_label_name_total",
        "rows_ignored_too_long_label_value_total",
        "insert_limit_reached_total",
        "select_limit_reached_total",
        "select_limit_timeout_total",
    } and kind == "max" and float(threshold) == 0:
        return "> 0"
    if metric_name in {
        "container_cpu_cfs_throttled_seconds_rate_max",
        "process_pressure_cpu_stalled_seconds_rate_max",
        "vmagent_container_cpu_cfs_throttled_seconds_rate_max",
        "vmagent_process_pressure_cpu_stalled_seconds_rate_max",
    }:
        return f"> {format_number(threshold)}"
    return f"{sign} {format_number(threshold)}"


def summary_metric_descriptions() -> dict[str, str]:
    return {
        "active_series": "Number of active time series with new data points during the last hour.",
        "total_datapoints": "Total stored raw samples excluding indexdb rows.",
        "data_size_gb": "Total data size on disk for VictoriaMetrics storage parts `storage/big` and `storage/small`, in GB.",
        "indexdb_size_gb": "Total on-disk size of VictoriaMetrics `indexdb/file` parts, in GB.",
        "bytes_per_sample": "Average storage footprint per stored sample.",
        "index_to_data_ratio": "Ratio of `indexdb/file` size to `storage/big + storage/small` size, following the VictoriaMetrics FAQ guidance. Elevated values can reflect high cardinality or churn, but they should be interpreted as a heuristic rather than a hard failure boundary.",
        "min_free_disk_space_gb": "Minimum free disk space across matched VictoriaMetrics targets, in GB.",
        "total_free_disk_space_gb": "Sum of free disk space across matched VictoriaMetrics targets, in GB.",
        "storage_full_eta_days": "Linear disk-full ETA in days across matched targets, computed from total vm_free_disk_space_bytes divided by total vm_data_size_bytes growth over the storage ETA lookback. Empty when growth is non-positive or unavailable.",
        "churn_rate_per_second": "Rate of newly created time series per second.",
        "new_series_total": "Total number of new time series created during the selected lookback window.",
        "new_series_per_active_series": "Reference ratio of new series created during the lookback window to the current active-series count. Useful for spotting sudden churn growth, but it should be interpreted together with service and label breakdowns.",
        "ingestion_rate_per_second": "Rate of inserted samples per second.",
        "remote_write_requests_max": "Maximum observed per-second rate of remote write HTTP requests reaching vmsingle on /api/v1/write during the lookback window.",
        "remote_write_http_errors_max": "Maximum observed per-second rate of remote write HTTP request errors on vmsingle for /api/v1/write during the lookback window.",
        "remote_write_http_error_ratio_max": "Maximum observed remote write HTTP error ratio during the lookback window, computed as max_over_time(rate(errors) / rate(requests)).",
        "remote_write_parser_read_errors_max": "Maximum observed per-second rate of vmsingle protoparser read errors for promremotewrite payloads during the lookback window.",
        "remote_write_parser_unmarshal_errors_max": "Maximum observed per-second rate of vmsingle unmarshal errors for promremotewrite payloads during the lookback window.",
        "rows_ignored_too_many_labels_total": "Total number of rows ignored by vmsingle because a time series contains too many labels during the lookback window.",
        "rows_ignored_too_long_label_name_total": "Total number of rows ignored by vmsingle because a label name is too long during the lookback window.",
        "rows_ignored_too_long_label_value_total": "Total number of rows ignored by vmsingle because a label value is too long during the lookback window.",
        "insert_limit_reached_total": "Total number of times vmsingle hit the concurrent insert limit during the lookback window.",
        "select_limit_reached_total": "Total number of times vmsingle hit the concurrent select/query limit during the lookback window.",
        "select_limit_timeout_total": "Total number of times read/query requests timed out while waiting for a concurrent select slot during the lookback window.",
        "query_requests_rate_per_second": "Maximum observed per-second rate of VictoriaMetrics metric/query API requests matching the configured or default request-path regex during the lookback window.",
        "avg_query_request_duration_seconds": "Maximum observed average duration of query-related VictoriaMetrics API requests during the lookback window, computed from one ratio expression so numerator and denominator come from the same timestamp.",
        "query_concurrency_limit": "VictoriaMetrics concurrent query capacity from vm_concurrent_select_capacity, with fallback to effective flag search.maxConcurrentRequests.",
        "vmalert_requests_rate_per_second": "Maximum observed per-second rate of vmalert datasource requests during the lookback window when vmalert_datasource_queries_total is available; otherwise uses a rule-execution rate proxy from vmalert_execution_total, not exact datasource request count. By default this query is global and is not automatically scoped by SELECTOR unless a custom VMALERT_REQUESTS_QUERY is provided.",
        "slow_inserts_ratio": "Maximum per-target share of inserted rows classified by VictoriaMetrics as slow, after summing across row types for each target and shown as a percentage. Values above the threshold can indicate disk or system pressure.",
        "container_cpu_cfs_throttled_seconds_rate_max": "Maximum summed per-second rate of container_cpu_cfs_throttled_seconds_total for container=vmsingle in the configured MONITORING_NAMESPACE during the lookback window.",
        "process_pressure_cpu_stalled_seconds_rate_max": "Maximum summed per-second rate of process_pressure_cpu_stalled_seconds_total for container=vmsingle in the configured MONITORING_NAMESPACE during the lookback window.",
        "vmagent_container_cpu_cfs_throttled_seconds_rate_max": "Maximum summed per-second rate of container_cpu_cfs_throttled_seconds_total for container=vmagent in the configured MONITORING_NAMESPACE during the lookback window.",
        "vmagent_process_pressure_cpu_stalled_seconds_rate_max": "Maximum summed per-second rate of process_pressure_cpu_stalled_seconds_total for container=vmagent in the configured MONITORING_NAMESPACE during the lookback window.",
        "vmagent_persistentqueue_bytes_dropped_ratio": "Maximum vmagent persistent-queue loss ratio across job and instance, computed as dropped bytes divided by written bytes. Empty value means no matching series or no computable ratio for the selected snapshot.",
        "remote_write_traffic_mbit_per_second": "Estimated remote write traffic from vmagent to the configured remoteWrite.url, calculated from vmagent_remotewrite_conn_bytes_written_total and shown in Mbit/s.",
        "vm_slow_queries_total_rate_per_second": "Per-second rate computed from the VictoriaMetrics counter vm_slow_queries_total.",
        "storage_growth_rows_per_second": "Estimated net growth of stored rows per second after deduplication adjustments.",
    }


def render_metric_table(
    title: str,
    summary: dict[str, Any],
    thresholds: dict[str, Any],
    metric_names: list[str],
) -> str:
    descriptions = summary_metric_descriptions()
    labels = summary_metric_labels()
    header_descriptions = table_column_descriptions("summary")
    rows = []
    for key in metric_names:
        if key not in summary:
            continue
        value = summary.get(key)
        rule = thresholds.get(key)
        status_label, status_class = threshold_status(value, rule if isinstance(rule, dict) else None)
        rows.append(
            "<tr>"
            f"<td class=\"col-metric\">{html_escape(labels.get(key, key))}</td>"
            f"<td class=\"col-value\">{format_summary_value(key, value)}</td>"
            f"<td class=\"col-threshold\">{html_escape(threshold_text(rule if isinstance(rule, dict) else None, metric_name=key))}</td>"
            f"<td class=\"col-status {status_class}\">{html_escape(status_label)}</td>"
            f"<td class=\"col-description\">{html_escape(descriptions.get(key, ''))}</td>"
            "</tr>"
        )
    body_html = "".join(rows)
    return (
        f"<section><h2>{html_escape(title)}</h2>"
        "<div class=\"table-wrap\"><table class=\"table-summary\"><thead><tr>"
        f"{header_html_with_tooltip('metric', header_descriptions.get('metric'), class_name='col-metric')}"
        f"{header_html_with_tooltip('value', header_descriptions.get('value'), class_name='col-value')}"
        f"{header_html_with_tooltip('alert_when', header_descriptions.get('threshold'), class_name='col-threshold')}"
        f"{header_html_with_tooltip('status', header_descriptions.get('status'), class_name='col-status')}"
        f"{header_html_with_tooltip('description', header_descriptions.get('description'), class_name='col-description')}"
        "</tr></thead>"
        f"<tbody>{body_html}</tbody></table></div></section>"
    )


def render_key_signals_table(summary: dict[str, Any], thresholds: dict[str, Any]) -> str:
    preferred_order = [
        "data_size_gb",
        "indexdb_size_gb",
        "index_to_data_ratio",
        "total_free_disk_space_gb",
        "storage_full_eta_days",
        "new_series_per_active_series",
        "query_requests_rate_per_second",
        "vmalert_requests_rate_per_second",
        "remote_write_requests_max",
        "remote_write_http_errors_max",
        "remote_write_http_error_ratio_max",
        "remote_write_parser_read_errors_max",
        "remote_write_parser_unmarshal_errors_max",
        "rows_ignored_too_many_labels_total",
        "rows_ignored_too_long_label_name_total",
        "rows_ignored_too_long_label_value_total",
        "insert_limit_reached_total",
        "select_limit_reached_total",
        "select_limit_timeout_total",
        "remote_write_traffic_mbit_per_second",
        "slow_inserts_ratio",
        "container_cpu_cfs_throttled_seconds_rate_max",
        "process_pressure_cpu_stalled_seconds_rate_max",
        "vmagent_container_cpu_cfs_throttled_seconds_rate_max",
        "vmagent_process_pressure_cpu_stalled_seconds_rate_max",
        "vmagent_persistentqueue_bytes_dropped_ratio",
        "active_series",
        "total_datapoints",
        "new_series_total",
    ]
    return render_metric_table("Key Signals", summary, thresholds, preferred_order)


def render_window_table(report: dict[str, Any]) -> str:
    windows = report.get("windows", {})
    header_descriptions = table_column_descriptions("lookback_window")
    rows = [
        ("query_time", report.get("query_time", "")),
        ("time_offset", windows.get("time_offset", "")),
        ("rate_window", windows.get("rate_window", "")),
        ("churn_lookback", windows.get("churn_lookback", "")),
        ("storage_eta_lookback", windows.get("storage_eta_lookback", "")),
    ]
    body_html = "".join("<tr>" f"<td>{html_escape(name)}</td>" f"<td>{html_escape(value)}</td>" "</tr>" for name, value in rows)
    return (
        "<section><h2>Lookback Window</h2>"
        "<div class=\"table-wrap\"><table><thead><tr>"
        f"{header_html_with_tooltip('parameter', header_descriptions.get('parameter'), class_name='col-parameter')}"
        f"{header_html_with_tooltip('value', header_descriptions.get('value'), class_name='col-value')}"
        "</tr></thead>"
        f"<tbody>{body_html}</tbody></table></div></section>"
    )


def render_findings(findings: Any) -> str:
    if not isinstance(findings, list) or not findings:
        return '<section class="group problems ok"><h2>Findings</h2><p>No warnings in this snapshot.</p></section>'
    cards: list[str] = []
    for item in findings:
        if not isinstance(item, dict):
            continue
        if str(item.get("severity", "")).lower() == "info":
            continue
        severity = html_escape(item.get("severity", "info"))
        area = html_escape(str(item.get("area", "finding")).replace("_", " ").title())
        message = html_escape(item.get("message", ""))
        cards.append(
            f'<article class="panel problem {severity}"><div class="problem-head"><h3>{area}</h3><span>{severity}</span></div><p>{message}</p></article>'
        )
    if not cards:
        return '<section class="group problems ok"><h2>Findings</h2><p>No warnings in this snapshot.</p></section>'
    return f'<section class="group problems"><h2>Findings</h2>{"".join(cards)}</section>'


def css_class_name(value: str) -> str:
    return re.sub(r"[^a-z0-9]+", "-", value.strip().lower()).strip("-")


def tooltip_text(value: str) -> str:
    return html_escape(value)


def display_header_name(value: str) -> str:
    normalized = value.strip().strip("_")
    if not normalized:
        return ""
    normalized = re.sub(r"([A-Z]+)([A-Z][a-z])", r"\1 \2", normalized)
    normalized = re.sub(r"([a-z0-9])([A-Z])", r"\1 \2", normalized)
    normalized = normalized.replace("_", " ").replace("-", " ")
    return " ".join(word.capitalize() for word in normalized.split())


def table_column_descriptions(name: str) -> dict[str, str]:
    descriptions: dict[str, dict[str, str]] = {
        "top_labels_by_unique_values": {"name": "Label name.", "value": "Number of distinct values observed for this label."},
        "vmsingle_configured_flags": {
            "Flag": "VictoriaMetrics command-line flag exposed via the flag metric.",
            "Value": "Configured value for this flag according to VictoriaMetrics self-metrics.",
            "Is Set": "Whether this flag was explicitly set at startup or left at its default value.",
        },
        "vmagent_configured_flags": {
            "Flag": "VMAgent command-line flag exposed via the flag metric.",
            "Value": "Configured value for this flag according to vmagent self-metrics.",
            "Is Set": "Whether this flag was explicitly set at startup or left at its default value.",
        },
        "vmsingle_all_effective_flags": {
            "Flag": "VictoriaMetrics command-line flag exposed via the flag metric.",
            "Value": "Effective value for this flag according to VictoriaMetrics self-metrics.",
            "Is Set": "Whether this flag was explicitly set at startup or left at its default value.",
        },
        "vmagent_all_effective_flags": {
            "Flag": "VMAgent command-line flag exposed via the flag metric.",
            "Value": "Effective value for this flag according to vmagent self-metrics.",
            "Is Set": "Whether this flag was explicitly set at startup or left at its default value.",
        },
        "top_services_by_series": {
            "Service": "Service or microservice group composed from the configured SERVICE_GROUP_BY_LABELS labels.",
            "Series": "Number of visible series returned by the current cardinality instant query for this service group.",
        },
        "top_services_by_new_series": {
            "Service": "Service or microservice group composed from the configured SERVICE_GROUP_BY_LABELS labels.",
            "New Series": "Total number of newly discovered scrape-time series from vmagent targets for this service group during the lookback window.",
        },
        "ingestion_by_cluster": {
            "cluster": "Cluster label value from vm_rows_inserted_total.",
            "Rows Per Second": "Per-second rate of inserted rows from vm_rows_inserted_total for this cluster label value.",
        },
        "tsdb_top_labels_by_memory_bytes": {
            "name": "Label name.",
            "value": "Estimated memory used by this label in TSDB metadata structures, in bytes.",
        },
        "label_distribution": {
            "Label": "Label name being evaluated across the metric scope used for this report.",
            "Series": "Number of scanned series carrying this label within the current report scope.",
            "Series Share": "Share of current active series represented by the scanned series carrying this label.",
            "Metrics Coverage": "Share of metrics in the current scan scope that contain this label.",
            "Unique Values": "How many distinct values of this label were observed in the scanned series within the current report scope.",
            "Classification": "Heuristic outcome: normal, observe, or review.",
        },
        "high_cardinality_metric_usage": {
            "metric": "Metric name joined from TSDB cardinality and metric usage views.",
            "series": "Number of series for this metric.",
            "top_label": "Main label driver for this metric based on sampled /api/v1/series analysis.",
            "top_label_unique_values": "Number of distinct values found for the main label driver.",
            "queryRequestsCount": "How many query requests touched this metric.",
            "lastRequestTimestamp": "Timestamp of the latest recorded query request for this metric.",
            "cleanupCandidate": "Whether the metric matches the configured low-usage or old-usage cleanup heuristics.",
        },
        "top_queries_by_sum_duration": {
            "Query": "PromQL query text recorded by VictoriaMetrics top-queries statistics.",
            "Sum Duration Seconds": "Total execution time accumulated by this query across the selected top-queries lifetime window, in seconds.",
            "Query Time Interval": "Original range selector used by the query, if VictoriaMetrics exposes it; '-' means no explicit range was recorded.",
            "Count": "How many times this query was executed within the selected top-queries lifetime window.",
        },
        "top_queries_by_avg_duration": {
            "Query": "PromQL query text recorded by VictoriaMetrics top-queries statistics.",
            "Avg Duration Seconds": "Average execution time of this query, in seconds.",
            "Query Time Interval": "Original range selector used by the query, if VictoriaMetrics exposes it; '-' means no explicit range was recorded.",
            "Count": "How many times this query was executed within the selected top-queries lifetime window.",
        },
        "top_queries_by_count": {
            "Query": "PromQL query text recorded by VictoriaMetrics top-queries statistics.",
            "Count": "How many times this query was executed within the selected top-queries lifetime window.",
            "Query Time Interval": "Original range selector used by the query, if VictoriaMetrics exposes it; '-' means no explicit range was recorded.",
        },
        "top_queries_by_avg_memory": {
            "Query": "PromQL query text recorded by VictoriaMetrics top-queries statistics.",
            "Avg Memory Usage Bytes": "Average memory allocated for this query according to VictoriaMetrics top-queries statistics, in bytes.",
            "Query Time Interval": "Original range selector used by the query, if VictoriaMetrics exposes it; '-' means no explicit range was recorded.",
            "Count": "How many times this query was executed within the selected top-queries lifetime window.",
        },
        "lookback_window": {"parameter": "Report parameter name.", "value": "Configured value used when collecting this snapshot."},
        "errors": {"query": "Query or API call name that failed.", "error": "Captured error message for that failed query or API call."},
        "summary": {
            "metric": "Summary metric name.",
            "value": "Measured or calculated value for this summary metric.",
            "threshold": "Condition under which this metric switches to alert status, if configured.",
            "status": "Threshold evaluation result for this metric.",
            "description": "What this summary metric represents.",
        },
    }
    return descriptions.get(name, {})


def header_html_with_tooltip(header: str, description: str | None, *, class_name: str) -> str:
    title_attr = f' title="{tooltip_text(description)}"' if description else ""
    return f"<th class=\"{class_name}\"{title_attr}>{html_escape(display_header_name(header))}</th>"


def show_infra_toggle(name: str) -> bool:
    return name in {"top_labels_by_unique_values", "label_distribution"}


def render_infra_toggle(name: str) -> str:
    if not show_infra_toggle(name):
        return ""
    checkbox_id = f"show-infra-{css_class_name(name)}"
    return (
        "<label class=\"table-toggle\" for=\""
        f"{html_escape(checkbox_id)}"
        "\">"
        f"<input type=\"checkbox\" id=\"{html_escape(checkbox_id)}\" data-table-toggle=\"{html_escape(name)}\" autocomplete=\"off\"> "
        "<span>Show Infra Labels</span>"
        "</label>"
    )


def table_intro(name: str, report: dict[str, Any]) -> str:
    analysis_modes = report.get("analysis_modes", {})
    if not isinstance(analysis_modes, dict):
        analysis_modes = {}
    if name == "label_distribution":
        scope = analysis_modes.get("label_distribution_scope", "top metrics")
        scope_text = (
            "This check was run across all metrics because FULL_LABEL_SCAN is enabled."
            if scope == "all metrics"
            else "This check was run only across top metrics for speed. Enable FULL_LABEL_SCAN to scan all metrics."
        )
        return (
            "<p class=\"table-note\"><strong>How to read this table:</strong> <code>normal</code> means the label looks broadly expected for the current metric scope. <code>observe</code> means the label either has a noticeable series share with limited metric spread or already shows a non-trivial number of unique values. <code>review</code> means the label combines high series share, low metric coverage, and high unique-value count, which is a stronger high-cardinality signal. </p>"
            "<p class=\"table-note\"><strong>Note:</strong> This table is built only from the selector-scoped <code>/api/v1/series</code> scan used for this report. <code>Series</code>, <code>Series Share</code>, and <code>Unique Values</code> therefore reflect the scanned report scope rather than global TSDB totals. For fully global reports without <code>SELECTOR</code>, sampled metrics may be extrapolated with TSDB per-metric totals; for selector-scoped reports that extrapolation is intentionally disabled to avoid mixing in foreign series. Non-infra labels with <code>Metrics Coverage=100%</code> are usually omitted, but labels with high <code>Unique Values</code> are kept because globally present business labels can still be major cardinality drivers. Infra labels are hidden by default and can be shown with the checkbox below. "
            f"{html_escape(scope_text)}</p>"
        )
    if name == "top_labels_by_unique_values":
        return "<p class=\"table-note\"><strong>Note:</strong> Infra labels are hidden by default and can be shown with the checkbox below.</p>"
    if name in {"vmsingle_configured_flags", "vmagent_configured_flags"}:
        return "<p class=\"table-note\"><strong>Note:</strong> This table shows only flags with <code>Is Set=true</code>, meaning they were explicitly passed to VictoriaMetrics at startup.</p>"
    if name in {"vmsingle_all_effective_flags", "vmagent_all_effective_flags"}:
        return "<p class=\"table-note\"><strong>Note:</strong> This table shows all effective flag values reported by the <code>flag</code> metric, including defaults where <code>Is Set=false</code>.</p>"
    if name == "top_services_by_series":
        labels = analysis_modes.get("service_group_by_labels", [])
        label_text = ", ".join(labels) if isinstance(labels, list) and labels else "namespace, job"
        return (
            "<p class=\"table-note\"><strong>Note:</strong> Services in this table are grouped by "
            f"{html_escape(label_text)}. It shows the current visible series footprint from the cardinality instant query for each group, not vm_cache_entries-based active series.</p>"
        )
    if name == "top_services_by_new_series":
        labels = analysis_modes.get("service_group_by_labels", [])
        label_text = ", ".join(labels) if isinstance(labels, list) and labels else "namespace, job"
        churn_window = report.get("windows", {}).get("churn_lookback", "")
        churn_suffix = f" during the {html_escape(churn_window)} lookback window" if churn_window else ""
        return (
            "<p class=\"table-note\"><strong>Note:</strong> Scrape targets in this table are grouped by "
            f"{html_escape(label_text)}. It shows how many new series each group created{churn_suffix}.</p>"
        )
    if name == "ingestion_by_cluster":
        label_name_value = analysis_modes.get("cluster_ingestion_label", "cluster")
        label_text = label_name_value if isinstance(label_name_value, str) and label_name_value else "cluster"
        return "<p class=\"table-note\"><strong>Note:</strong> This optional table groups <code>vm_rows_inserted_total</code> rate by " f"<code>{html_escape(label_text)}</code>.</p>"
    if name == "high_cardinality_metric_usage":
        parts: list[str] = []
        if analysis_modes.get("top_metric_label_drivers_partial_sample"):
            parts.append(
                "Some Top Label values come from a partial <code>/api/v1/series</code> sample because <code>SERIES_SAMPLE_LIMIT</code> was reached."
            )
        top_limit = analysis_modes.get("top_limit")
        analysis_limit = analysis_modes.get("metric_label_analysis_limit")
        configured_analysis_limit = analysis_modes.get("configured_metric_label_analysis_limit")
        if (
            isinstance(top_limit, int)
            and isinstance(analysis_limit, int)
            and isinstance(configured_analysis_limit, int)
            and configured_analysis_limit < top_limit
        ):
            parts.append(
                "The configured <code>METRIC_LABEL_ANALYSIS_LIMIT</code> was lower than "
                f"<code>TOP_LIMIT</code> ({top_limit}), so the effective analysis limit was automatically raised to <code>{analysis_limit}</code>."
            )
        if parts:
            return "<p class=\"table-note\"><strong>Note:</strong> " + " ".join(parts) + "</p>"
    if name in {
        "top_queries_by_sum_duration",
        "top_queries_by_avg_duration",
        "top_queries_by_count",
        "top_queries_by_avg_memory",
    }:
        top_limit = analysis_modes.get("top_queries_limit", 10)
        lookback = analysis_modes.get("top_queries_lookback", "24h")
        return (
            "<p class=\"table-note\"><strong>Note:</strong> This table is read from "
            "<code>/api/v1/status/top_queries</code> with "
            f"<code>topN={html_escape(top_limit)}</code> and <code>maxLifetime={html_escape(lookback)}</code>.</p>"
        )
    return ""


def render_table(name: str, title: str, rows: list[dict[str, Any]], report: dict[str, Any]) -> str:
    if not rows:
        return ""
    headers = [header for header in rows[0].keys() if not header.startswith("__")]
    table_class = f"table-{css_class_name(name)}"
    column_descriptions = table_column_descriptions(name)
    header_html = "".join(header_html_with_tooltip(header, column_descriptions.get(header), class_name=f"col-{css_class_name(header)}") for header in headers)
    body_html = "".join(
        (
            "<tr"
            + (
                f' class="infra-label-row" data-table-name="{html_escape(name)}" style="display:none"'
                if row.get("__infraLabel") and show_infra_toggle(name)
                else ""
            )
            + ">"
            + "".join(
                f"<td class=\"{table_cell_class(name, header, row.get(header))}\">{format_table_cell(name, header, row.get(header))}</td>"
                for header in headers
            )
            + "</tr>"
        )
        for row in rows
    )
    table_html = (
        f"{table_intro(name, report)}"
        f"{render_infra_toggle(name)}"
        f"<div class=\"table-wrap\"><table class=\"{table_class}\"><thead><tr>{header_html}</tr></thead><tbody>{body_html}</tbody></table></div>"
    )
    if name in {"vmsingle_all_effective_flags", "vmagent_all_effective_flags"}:
        return f"<section><details class=\"collapsible-section\"><summary>{html_escape(title)}</summary>{table_html}</details></section>"
    return f"<section><h2>{html_escape(title)}</h2>{table_html}</section>"


def table_title(name: str) -> str:
    titles = {
        "vmsingle_configured_flags": "VMSingle Configured Flags",
        "vmagent_configured_flags": "VMAgent Configured Flags",
        "vmsingle_all_effective_flags": "VMSingle All Effective Flags",
        "vmagent_all_effective_flags": "VMAgent All Effective Flags",
        "top_services_by_series": "Top Services By Series",
        "top_services_by_new_series": "Top Services By New Series",
        "ingestion_by_cluster": "Ingestion By Cluster",
        "top_labels_by_unique_values": "Top Labels By Unique Values",
        "tsdb_top_labels_by_memory_bytes": "Tsdb Top Labels By Memory Bytes",
        "label_distribution": "Label Distribution",
        "high_cardinality_metric_usage": "High Cardinality Metric Usage",
        "top_queries_by_sum_duration": "Queries With Most Summary Time To Execute",
        "top_queries_by_avg_duration": "Most Heavy Queries",
        "top_queries_by_count": "Most Frequently Executed Queries",
        "top_queries_by_avg_memory": "Queries With Most Memory To Execute",
    }
    return titles.get(name, name.replace("_", " ").title())


def ordered_table_sections(report: dict[str, Any]) -> list[tuple[str, list[dict[str, Any]]]]:
    tables = report.get("tables", {})
    if not isinstance(tables, dict):
        return []
    configuration_names = {
        "vmsingle_configured_flags",
        "vmagent_configured_flags",
        "vmsingle_all_effective_flags",
        "vmagent_all_effective_flags",
    }
    preferred_order = [
        "high_cardinality_metric_usage",
        "top_labels_by_unique_values",
        "tsdb_top_labels_by_memory_bytes",
        "label_distribution",
        "top_services_by_series",
        "top_services_by_new_series",
        "ingestion_by_cluster",
        "top_queries_by_sum_duration",
        "top_queries_by_avg_duration",
        "top_queries_by_count",
        "top_queries_by_avg_memory",
    ]
    sections: list[tuple[str, list[dict[str, Any]]]] = []
    seen: set[str] = set()
    for name in preferred_order:
        rows = tables.get(name)
        if isinstance(rows, list):
            sections.append((name, rows))
            seen.add(name)
    for name, rows in tables.items():
        if name in seen or name in configuration_names or not isinstance(rows, list):
            continue
        sections.append((name, rows))
    return sections


def render_configuration_section(report: dict[str, Any]) -> str:
    tables = report.get("tables", {})
    if not isinstance(tables, dict):
        return ""
    names = [
        "vmsingle_configured_flags",
        "vmagent_configured_flags",
        "vmsingle_all_effective_flags",
        "vmagent_all_effective_flags",
    ]
    rendered = "".join(
        render_table(name, table_title(name), rows, report)
        for name in names
        if isinstance((rows := tables.get(name)), list) and rows
    )
    if not rendered:
        return ""
    return f"<section><h2>Configuration</h2>{rendered}</section>"


def write_html_report(report: dict[str, Any], output: str) -> None:
    summary = report.get("summary", {})
    thresholds = report.get("thresholds", {})
    findings_rows = report.get("findings", [])
    findings_html = render_findings(findings_rows)
    key_signals_html = render_key_signals_table(summary if isinstance(summary, dict) else {}, thresholds if isinstance(thresholds, dict) else {})
    window_html = render_window_table(report)
    configuration_html = render_configuration_section(report)
    tables_html = "".join(render_table(name, table_title(name), rows, report) for name, rows in ordered_table_sections(report))
    errors = report.get("errors", {})
    errors_html = ""
    if errors:
        error_headers = table_column_descriptions("errors")
        error_rows = "".join(
            f"<tr><td>{html_escape(name)}</td><td>{html_escape(message)}</td></tr>"
            for name, message in errors.items()
            if message
        )
        errors_html = (
            "<section><h2>Errors</h2><div class=\"table-wrap\"><table><thead><tr>"
            f"{header_html_with_tooltip('query', error_headers.get('query'), class_name='col-query')}"
            f"{header_html_with_tooltip('error', error_headers.get('error'), class_name='col-error')}"
            f"</tr></thead><tbody>{error_rows}</tbody></table></div></section>"
        )

    html = f"""<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>VictoriaMetrics Analysis Report</title>
  <style>
    :root {{
      --bg: #f5f1e8;
      --ink: #1d2528;
      --muted: #687276;
      --card: #fffaf0;
      --line: #dfd5c3;
      --accent: #b85c38;
      --accent-2: #264653;
      --error: #b42318;
      --skip: #8a5a00;
    }}
    * {{ box-sizing: border-box; }}
    body {{
      margin: 0;
      background:
        radial-gradient(circle at top left, rgba(184, 92, 56, .18), transparent 30rem),
        linear-gradient(135deg, #f5f1e8 0%, #ece2d0 100%);
      color: var(--ink);
      font: 15px/1.5 "Aptos", "Segoe UI", sans-serif;
    }}
    main {{ width: min(1180px, calc(100% - 32px)); margin: 0 auto; padding: 32px 0 56px; }}
    .hero {{
      border: 1px solid var(--line);
      border-radius: 28px;
      padding: 32px;
      background: linear-gradient(135deg, rgba(38, 70, 83, .95), rgba(38, 70, 83, .78));
      color: #fffaf0;
      box-shadow: 0 24px 70px rgba(38, 70, 83, .18);
    }}
    .hero p {{
      margin: 0 0 6px;
      letter-spacing: .12em;
      text-transform: uppercase;
      color: #f2c9a8;
    }}
    .hero h1 {{ margin: 0; font-size: clamp(32px, 6vw, 58px); line-height: .95; }}
    .subtle {{ display: block; margin-top: 14px; color: rgba(255, 250, 240, .72); word-break: break-all; }}
    .group {{ margin-top: 18px; padding: 18px; }}
    .panel {{
      margin: 12px 0;
      padding: 16px;
      overflow-x: auto;
      border: 1px solid var(--line);
      border-radius: 20px;
      background: rgba(255, 250, 240, .82);
      box-shadow: 0 16px 38px rgba(38, 70, 83, .08);
    }}
    section {{
      margin-top: 18px;
      padding: 18px;
      border: 1px solid var(--line);
      border-radius: 20px;
      background: rgba(255, 250, 240, .82);
      box-shadow: 0 16px 38px rgba(38, 70, 83, .08);
    }}
    section h2 {{ margin: 0 0 14px; font-size: 24px; color: var(--accent-2); }}
    p.meta, p.table-note {{ color: var(--muted); }}
    p.meta {{ margin: 0 0 12px; }}
    p.table-note {{ margin: 6px 0 12px; line-height: 1.5; }}
    .table-toggle {{
      display: inline-flex;
      align-items: center;
      gap: 8px;
      margin: 0 0 12px;
      color: var(--accent-2);
      font-weight: 700;
      cursor: pointer;
    }}
    .table-toggle input {{
      inline-size: 16px;
      block-size: 16px;
      accent-color: var(--accent);
    }}
    .table-wrap {{
      overflow-x: auto;
      margin: 12px 0;
      padding: 16px;
      border: 1px solid var(--line);
      border-radius: 20px;
      background: rgba(255, 250, 240, .82);
      box-shadow: 0 16px 38px rgba(38, 70, 83, .08);
    }}
    table {{ width: 100%; border-collapse: collapse; min-width: 520px; }}
    th, td {{
      padding: 9px 10px;
      border-bottom: 1px solid var(--line);
      text-align: left;
      vertical-align: top;
    }}
    th {{
      color: var(--accent-2);
      font-size: 12px;
      text-transform: uppercase;
      letter-spacing: .06em;
      background: rgba(38, 70, 83, .06);
    }}
    th[title] {{ cursor: help; text-decoration: underline dotted; text-underline-offset: 3px; }}
    tr:hover td {{ background: #e6dcd5; }}
    code, pre {{
      font-family: "Cascadia Code", "SFMono-Regular", Consolas, monospace;
      font-size: 12px;
      white-space: pre-wrap;
      word-break: break-word;
    }}
    .table-note code, .hero code {{
      display: inline-block;
      padding: 1px 6px;
      border: 1px solid rgba(38, 70, 83, .18);
      border-radius: 7px;
      background: rgba(38, 70, 83, .09);
      color: var(--accent-2);
      font-weight: 800;
      white-space: normal;
    }}
    details.collapsible-section {{
      margin-top: 12px;
      border: 1px solid var(--line);
      border-radius: 20px;
      background: rgba(255, 250, 240, .82);
      padding: 16px;
      box-shadow: 0 16px 38px rgba(38, 70, 83, .08);
    }}
    details.collapsible-section summary {{ cursor: pointer; color: var(--accent); font-weight: 700; }}
    .table-summary .col-metric {{ width: 1%; white-space: nowrap; }}
    .table-summary .col-value {{ width: 1%; white-space: nowrap; }}
    .table-summary .col-threshold {{ width: 1%; white-space: nowrap; }}
    .table-summary td.col-metric {{ background: rgba(184, 92, 56, .06); color: var(--accent); font-weight: 700; }}
    .table-summary tr:hover td.col-metric {{ background: #e6dcd5; }}
    .table-label-distribution {{ width: 100%; table-layout: auto; }}
    .table-label-distribution .col-label {{ width: 100%; }}
    .table-label-distribution .col-series {{ white-space: nowrap; width: 1%; }}
    .table-label-distribution .col-share {{ white-space: nowrap; width: 1%; }}
    .table-label-distribution .col-metrics-coverage {{ white-space: nowrap; width: 1%; }}
    .table-label-distribution .col-unique-values {{ white-space: nowrap; width: 1%; }}
    .table-label-distribution .col-classification {{ white-space: nowrap; width: 1%; }}
    .table-label-distribution td.classification-normal {{
      background: rgba(38, 70, 83, .06);
      color: var(--accent-2);
      font-weight: 700;
    }}
    .table-label-distribution td.classification-observe {{
      background: #e6d9bc;
      color: #834e21;
      font-weight: 700;
    }}
    .table-label-distribution td.classification-review {{
      background: rgba(180, 35, 24, .14);
      color: var(--error);
      font-weight: 800;
    }}
    .status-ok {{ background: #dce8ca; color: #344C2A; font-weight: 700; }}
    .status-alert {{ background: rgba(180, 35, 24, .12); color: var(--error); font-weight: 700; }}
    .status-na {{ color: var(--muted); }}
    .problems {{
      border-color: rgba(180, 35, 24, .32);
      background: linear-gradient(135deg, rgba(255, 244, 242, .96), rgba(255, 250, 240, .84));
    }}
    .problems.ok {{ border-color: rgba(38, 70, 83, .18); background: rgba(255, 250, 240, .82); }}
    .problem {{
      border-color: rgba(180, 35, 24, .42);
      border-left: 7px solid var(--error);
      background: #fff4f2;
      box-shadow: 0 18px 42px rgba(180, 35, 24, .1);
    }}
    .problem-head {{ display: flex; align-items: center; justify-content: space-between; gap: 12px; }}
    .problem-head h3 {{ color: var(--error); margin: 0; }}
    .problem-head span {{
      display: inline-block;
      padding: 3px 9px;
      border-radius: 999px;
      background: rgba(180, 35, 24, .12);
      color: var(--error);
      font-size: 11px;
      font-weight: 800;
      letter-spacing: .08em;
      text-transform: uppercase;
    }}
    ul {{ margin: 0; padding-left: 22px; }}
    pre {{
      margin: 0;
      background: #0f172a;
      color: #e2e8f0;
      padding: 16px;
      border-radius: 20px;
      overflow: auto;
    }}
    @media (max-width: 640px) {{
      main {{ width: min(100% - 20px, 1180px); padding-top: 16px; }}
      .hero {{ padding: 22px; border-radius: 22px; }}
      section {{ padding: 12px; }}
    }}
  </style>
</head>
<body>
  <main>
    <section class="hero">
      <p>VictoriaMetrics Analysis</p>
      <h1>VictoriaMetrics Analysis Report</h1>
      <span class="subtle">Generated: {html_escape(report.get("generated_at", ""))}</span>
    </section>
    {window_html}
    {findings_html}
    {key_signals_html}
    {configuration_html}
    {tables_html}
    {errors_html}
    <section>
      <details class="collapsible-section">
        <summary>Raw JSON</summary>
        <pre>{html_escape(json.dumps(report, ensure_ascii=True, indent=2))}</pre>
      </details>
    </section>
  </main>
  <script>
    document.querySelectorAll('[data-table-toggle]').forEach((checkbox) => {{
      checkbox.checked = false;
      const tableName = checkbox.getAttribute('data-table-toggle');
      document.querySelectorAll(`tr.infra-label-row[data-table-name="${{tableName}}"]`).forEach((row) => {{
        row.style.display = 'none';
      }});
      checkbox.addEventListener('change', () => {{
        document.querySelectorAll(`tr.infra-label-row[data-table-name="${{tableName}}"]`).forEach((row) => {{
          row.style.display = checkbox.checked ? 'table-row' : 'none';
        }});
      }});
    }});
  </script>
</body>
</html>
"""
    Path(output).write_text(html, encoding="utf-8")
