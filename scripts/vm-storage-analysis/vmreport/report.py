from __future__ import annotations

import json
import math
from datetime import UTC, datetime
from pathlib import Path
from typing import Any

from .common import (
    bytes_to_gb,
    filter_flag_rows,
    find_flag_value,
    flag_rows,
    grouped_sample_rows,
    maybe_round,
    parse_group_by_labels,
    positive_finite,
    safe_div,
    sample_value,
    sum_flag_values,
    table_rows,
    utc_iso,
)
from .data import (
    build_high_cardinality_usage_rows,
    build_label_distribution_rows,
    effective_metric_label_analysis_limit,
    merged_errors,
    metric_label_driver_rows,
    parse_metric_usage_rows,
    parse_top_queries_rows,
    tsdb_queries,
    tsdb_rows,
)
from .html import format_percentage

INFRA_LABELS = {
    "job",
    "instance",
    "namespace",
    "service",
    "endpoint",
    "pod",
    "container",
    "node",
    "mountpoint",
    "cluster",
    "prometheus",
}
LABEL_FINDINGS_EXCLUDED = INFRA_LABELS | {"le"}
LABEL_FINDINGS_UNIQUE_VALUES_THRESHOLD = 1000
TOP_LABELS_BY_UNIQUE_VALUES_EXCLUDED = {"__name__"}


def findings(report: dict[str, Any], args) -> list[dict[str, Any]]:
    summary = report["summary"]
    result: list[dict[str, Any]] = []
    errors = report.get("errors", {})
    if isinstance(errors, dict) and errors:
        failed_queries = len(errors)
        sample_names = ", ".join(list(errors.keys())[:3])
        more_suffix = "" if failed_queries <= 3 else f", and {failed_queries - 3} more"
        result.append(
            {
                "severity": "warning",
                "area": "data_collection",
                "message": (
                    f"{failed_queries} VictoriaMetrics query/API call(s) failed while building this report "
                    f"({sample_names}{more_suffix}). Some sections may be incomplete or stale "
                    "until connectivity or API errors are resolved."
                ),
            }
        )
    index_ratio = summary.get("index_to_data_ratio")
    if isinstance(index_ratio, (int, float)) and index_ratio >= args.index_ratio_warning_threshold:
        result.append(
            {
                "severity": "warning",
                "area": "indexdb",
                "message": (
                    f"indexdb/file to storage data ratio is {index_ratio:.3f}, which is above the configured "
                    f"threshold {args.index_ratio_warning_threshold:.3f}. This can be a sign of elevated "
                    "cardinality or churn and is best reviewed together with the cardinality sections below."
                ),
            }
        )
    slow_inserts_ratio = summary.get("slow_inserts_ratio")
    if isinstance(slow_inserts_ratio, (int, float)) and slow_inserts_ratio >= args.slow_inserts_warning_threshold:
        result.append(
            {
                "severity": "warning",
                "area": "ingestion",
                "message": (
                    f"slow insert ratio is {format_percentage(float(slow_inserts_ratio))}, above threshold "
                    f"{format_percentage(float(args.slow_inserts_warning_threshold))}. "
                    "This can indicate disk or system pressure."
                ),
            }
        )
    remote_write_http_error_ratio = summary.get("remote_write_http_error_ratio_max")
    if isinstance(remote_write_http_error_ratio, (int, float)) and remote_write_http_error_ratio > 0.01:
        result.append(
            {
                "severity": "warning",
                "area": "remote_write_errors",
                "message": (
                    f"remote write HTTP error ratio peak reached "
                    f"{format_percentage(float(remote_write_http_error_ratio))} "
                    f"during the current lookback window ({args.churn_lookback}), above the 1% warning threshold."
                ),
            }
        )
    remote_write_http_errors_max = summary.get("remote_write_http_errors_max")
    if isinstance(remote_write_http_errors_max, (int, float)) and remote_write_http_errors_max > 0:
        result.append(
            {
                "severity": "info",
                "area": "remote_write_errors",
                "message": (
                    f"remote write HTTP errors peaked at {remote_write_http_errors_max:.3f} requests/s "
                    f"during the current lookback window ({args.churn_lookback})."
                ),
            }
        )
    remote_write_parser_read_errors_max = summary.get("remote_write_parser_read_errors_max")
    if isinstance(remote_write_parser_read_errors_max, (int, float)) and remote_write_parser_read_errors_max > 0:
        result.append(
            {
                "severity": "info",
                "area": "remote_write_errors",
                "message": (
                    f"remote write parser read errors peaked at {remote_write_parser_read_errors_max:.3f} requests/s "
                    f"during the current lookback window ({args.churn_lookback})."
                ),
            }
        )
    remote_write_parser_unmarshal_errors_max = summary.get("remote_write_parser_unmarshal_errors_max")
    if (
        isinstance(remote_write_parser_unmarshal_errors_max, (int, float))
        and remote_write_parser_unmarshal_errors_max > 0
    ):
        result.append(
            {
                "severity": "info",
                "area": "remote_write_errors",
                "message": (
                    f"remote write parser unmarshal errors peaked at "
                    f"{remote_write_parser_unmarshal_errors_max:.3f} requests/s "
                    f"during the current lookback window ({args.churn_lookback})."
                ),
            }
        )
    ignored_rows_metrics = [
        ("too_many_labels", summary.get("rows_ignored_too_many_labels_total")),
        ("too_long_label_name", summary.get("rows_ignored_too_long_label_name_total")),
        ("too_long_label_value", summary.get("rows_ignored_too_long_label_value_total")),
    ]
    for reason, value in ignored_rows_metrics:
        if isinstance(value, (int, float)) and value > 0:
            result.append(
                {
                    "severity": "warning",
                    "area": "rows_ignored",
                    "message": (
                        f"rows ignored with reason {reason!r} totalled {value:.0f} row(s) "
                        f"during the current lookback window ({args.churn_lookback})."
                    ),
                }
            )
    limit_reached_warning_threshold = 50
    insert_limit_reached = summary.get("insert_limit_reached_total")
    if isinstance(insert_limit_reached, (int, float)) and insert_limit_reached > 0:
        result.append(
            {
                "severity": "warning" if insert_limit_reached >= limit_reached_warning_threshold else "info",
                "area": "ingestion_pressure",
                "message": (
                    f"concurrent insert limit was reached {insert_limit_reached:.0f} time(s) "
                    f"during the current lookback window ({args.churn_lookback}). "
                    "Ingestion concurrency may be saturated."
                ),
            }
        )
    select_limit_timeout = summary.get("select_limit_timeout_total")
    if isinstance(select_limit_timeout, (int, float)) and select_limit_timeout > 0:
        result.append(
            {
                "severity": "warning",
                "area": "query_pressure",
                "message": (
                    f"concurrent select wait timeout occurred {select_limit_timeout:.0f} time(s) "
                    f"during the current lookback window ({args.churn_lookback}). "
                    "Some reads likely waited too long for a query slot."
                ),
            }
        )
    high_cardinality_metrics = report["tables"].get("high_cardinality_metric_usage", [])
    if high_cardinality_metrics:
        top_metric = high_cardinality_metrics[0]
        metric_name = top_metric.get("metric", "")
        metric_series = top_metric.get("series", 0.0)
        total_active = summary.get("active_series")
        share = safe_div(float(metric_series), float(total_active)) if total_active else None
        if isinstance(share, float) and share >= 0.2:
            result.append(
                {
                    "severity": "info",
                    "area": "cardinality",
                    "message": (
                        f"metric {metric_name!r} alone accounts for about {share:.1%} of active series in this view. "
                        "It is a good candidate for label review or selective relabeling."
                    ),
                }
            )
    cleanup_candidates = [
        row
        for row in report["tables"].get("high_cardinality_metric_usage", [])
        if row.get("cleanupCandidate")
    ]
    if cleanup_candidates:
        candidate = cleanup_candidates[0]
        result.append(
            {
                "severity": "warning",
                "area": "cleanup_candidates",
                "message": (
                    f"metric {candidate.get('metric')!r} has high series count and weak usage signals "
                    f"(requests={candidate.get('queryRequestsCount')}, "
                    f"last_request={candidate.get('lastRequestTimestamp')}, "
                    f"driver_label={candidate.get('top_label')!r}). Review it as a cleanup candidate."
                ),
            }
        )
    label_reviews = [
        row
        for row in report["tables"].get("label_distribution", [])
        if row.get("Classification") == "review"
        and isinstance(row.get("Label"), str)
        and row.get("Label") not in LABEL_FINDINGS_EXCLUDED
        and isinstance(row.get("Series Share"), (int, float))
        and float(row.get("Series Share")) >= 0.20
        and isinstance(row.get("Metrics Coverage"), (int, float))
        and float(row.get("Metrics Coverage")) <= min(args.label_low_metric_coverage_threshold, 0.15)
        and isinstance(row.get("Unique Values"), (int, float))
        and float(row.get("Unique Values")) >= LABEL_FINDINGS_UNIQUE_VALUES_THRESHOLD
    ]
    if label_reviews:
        candidate = label_reviews[0]
        scope = "all scanned metrics" if args.full_label_scan else "scanned top metrics"
        result.append(
            {
                "severity": "warning",
                "area": "label_localization",
                "message": (
                    f"label {candidate.get('Label')!r} occupies about "
                    f"{candidate.get('Series Share'):.1%} of active series, while its metrics coverage is only "
                    f"{candidate.get('Metrics Coverage'):.1%} across {scope}, "
                    f"with about {int(candidate.get('Unique Values')):,} unique values. "
                    "This is a localized but high-impact label and may be worth "
                    "reviewing, especially if it belongs to top high-cardinality metrics."
                ),
            }
        )
    if not result:
        result.append(
            {
                "severity": "info",
                "area": "summary",
                "message": "No configured thresholds were exceeded in this snapshot.",
            }
        )
    return result


def build_report(
    args,
    snapshot_time: str,
    raw_results: dict[str, Any],
    tsdb_results: dict[str, Any],
    query_map: dict[str, str],
    client=None,
) -> dict[str, Any]:
    if args.dry_run:
        return {
            "report_type": "victoriametrics_analysis_report",
            "generated_at": utc_iso(datetime.now(UTC)),
            "query_time": snapshot_time,
            "victoriametrics_url": args.victoriametrics_url,
            "selectors": {"selector": args.selector},
            "windows": {
                "rate_window": args.rate_window,
                "churn_lookback": args.churn_lookback,
                "storage_eta_lookback": args.storage_eta_lookback,
                "time_offset": args.time_offset,
            },
            "tsdb_queries": tsdb_queries(args),
            "queries": query_map,
        }

    vmsingle_flag_rows = flag_rows(raw_results.get("vmsingle_flags"))
    vmagent_flag_rows = flag_rows(raw_results.get("vmagent_flags"))
    vmsingle_configured_flags = filter_flag_rows(vmsingle_flag_rows, is_set=True)
    vmagent_configured_flags = filter_flag_rows(vmagent_flag_rows, is_set=True)
    vmsingle_effective_flags = filter_flag_rows(vmsingle_flag_rows, require_value=True)
    vmagent_effective_flags = filter_flag_rows(vmagent_flag_rows, require_value=True)
    storage_full_eta_days_value = positive_finite(sample_value(raw_results.get("storage_full_eta_days")))
    summary = {
        "active_series": sample_value(raw_results.get("active_series")),
        "total_datapoints": sample_value(raw_results.get("total_datapoints")),
        "data_size_bytes": sample_value(raw_results.get("data_size_bytes")),
        "data_size_gb": bytes_to_gb(sample_value(raw_results.get("data_size_bytes"))),
        "indexdb_size_bytes": sample_value(raw_results.get("indexdb_size_bytes")),
        "indexdb_size_gb": bytes_to_gb(sample_value(raw_results.get("indexdb_size_bytes"))),
        "bytes_per_sample": sample_value(raw_results.get("bytes_per_sample")),
        "index_to_data_ratio": sample_value(raw_results.get("index_to_data_ratio")),
        "min_free_disk_space_bytes": sample_value(raw_results.get("min_free_disk_space_bytes")),
        "min_free_disk_space_gb": bytes_to_gb(sample_value(raw_results.get("min_free_disk_space_bytes"))),
        "total_free_disk_space_bytes": sample_value(raw_results.get("total_free_disk_space_bytes")),
        "total_free_disk_space_gb": bytes_to_gb(sample_value(raw_results.get("total_free_disk_space_bytes"))),
        "storage_full_eta_days": (
            math.floor(storage_full_eta_days_value) if storage_full_eta_days_value is not None else None
        ),
        "churn_rate_per_second": sample_value(raw_results.get("churn_rate_per_second")),
        "new_series_total": sample_value(raw_results.get("new_series_total")),
        "ingestion_rate_per_second": sample_value(raw_results.get("ingestion_rate_per_second")),
        "remote_write_requests_max": sample_value(raw_results.get("remote_write_requests_max")),
        "remote_write_http_errors_max": sample_value(raw_results.get("remote_write_http_errors_max")),
        "remote_write_http_error_ratio_max": sample_value(raw_results.get("remote_write_http_error_ratio_max")),
        "remote_write_parser_read_errors_max": sample_value(raw_results.get("remote_write_parser_read_errors_max")),
        "remote_write_parser_unmarshal_errors_max": sample_value(
            raw_results.get("remote_write_parser_unmarshal_errors_max")
        ),
        "rows_ignored_too_many_labels_total": sample_value(raw_results.get("rows_ignored_too_many_labels_total")),
        "rows_ignored_too_long_label_name_total": sample_value(
            raw_results.get("rows_ignored_too_long_label_name_total")
        ),
        "rows_ignored_too_long_label_value_total": sample_value(
            raw_results.get("rows_ignored_too_long_label_value_total")
        ),
        "insert_limit_reached_total": sample_value(raw_results.get("insert_limit_reached_total")),
        "select_limit_reached_total": sample_value(raw_results.get("select_limit_reached_total")),
        "select_limit_timeout_total": sample_value(raw_results.get("select_limit_timeout_total")),
        "slow_inserts_ratio": sample_value(raw_results.get("slow_inserts_ratio")),
        "vm_slow_queries_total_rate_per_second": sample_value(raw_results.get("vm_slow_queries_total_rate_per_second")),
        "container_cpu_cfs_throttled_seconds_rate_max": sample_value(
            raw_results.get("container_cpu_cfs_throttled_seconds_rate_max")
        ),
        "process_pressure_cpu_stalled_seconds_rate_max": sample_value(
            raw_results.get("process_pressure_cpu_stalled_seconds_rate_max")
        ),
        "vmagent_container_cpu_cfs_throttled_seconds_rate_max": sample_value(
            raw_results.get("vmagent_container_cpu_cfs_throttled_seconds_rate_max")
        ),
        "vmagent_process_pressure_cpu_stalled_seconds_rate_max": sample_value(
            raw_results.get("vmagent_process_pressure_cpu_stalled_seconds_rate_max")
        ),
        "vmagent_persistentqueue_bytes_dropped_ratio": sample_value(
            raw_results.get("vmagent_persistentqueue_bytes_dropped_ratio")
        ),
        "remote_write_traffic_mbit_per_second": sample_value(raw_results.get("remote_write_traffic_mbit_per_second")),
        "storage_growth_rows_per_second": sample_value(raw_results.get("storage_growth_rows_per_second")),
        "query_requests_rate_per_second": sample_value(raw_results.get("query_requests_rate_per_second")),
        "avg_query_request_duration_seconds": sample_value(raw_results.get("avg_query_request_duration_seconds")),
        "query_concurrency_limit": sample_value(raw_results.get("query_concurrency_limit")),
        "vmalert_requests_rate_per_second": sample_value(raw_results.get("vmalert_requests_rate_per_second")),
    }
    if summary["query_concurrency_limit"] is None:
        summary["query_concurrency_limit"] = sum_flag_values(
            vmsingle_flag_rows,
            "search.maxConcurrentRequests",
            "maxConcurrentRequests",
        )
    if summary["query_concurrency_limit"] is None:
        summary["query_concurrency_limit"] = find_flag_value(
            vmsingle_flag_rows,
            "search.maxConcurrentRequests",
            "maxConcurrentRequests",
        )
    if not isinstance(summary["remote_write_http_error_ratio_max"], (int, float)) or not math.isfinite(
        summary["remote_write_http_error_ratio_max"]
    ):
        summary["remote_write_http_error_ratio_max"] = None
    if not isinstance(summary["slow_inserts_ratio"], (int, float)) or not math.isfinite(
        summary["slow_inserts_ratio"]
    ):
        summary["slow_inserts_ratio"] = None
    summary["new_series_per_active_series"] = maybe_round(
        safe_div(summary["new_series_total"], summary["active_series"]),
        6,
    )
    service_group_labels = parse_group_by_labels(args.service_group_by_labels)
    metric_label_analysis_limit = effective_metric_label_analysis_limit(args)

    top_metrics_by_cardinality = [
        {"__name__": row["name"], "series": row["value"]}
        for row in tsdb_rows(tsdb_results.get("tsdb_stats", {}).get("seriesCountByMetricName"))
    ]
    top_metric_label_drivers = metric_label_driver_rows(
        tsdb_results.get("tsdb_stats"),
        tsdb_results.get("metric_series"),
    )
    top_metric_label_drivers_sampled = any(
        row.get("sampled") for row in top_metric_label_drivers if isinstance(row, dict)
    )
    top_metric_label_drivers = [
        {key: value for key, value in row.items() if key != "sampled"}
        for row in top_metric_label_drivers
        if isinstance(row, dict)
    ]
    metric_usage_rows = parse_metric_usage_rows(tsdb_results.get("metric_names_stats"))
    top_queries_payload = tsdb_results.get("top_queries")
    top_queries_by_sum_duration = parse_top_queries_rows(
        top_queries_payload.get("topBySumDuration") if isinstance(top_queries_payload, dict) else [],
        value_field="sumDurationSeconds",
        value_key="Sum Duration Seconds",
    )
    top_queries_by_avg_duration = parse_top_queries_rows(
        top_queries_payload.get("topByAvgDuration") if isinstance(top_queries_payload, dict) else [],
        value_field="avgDurationSeconds",
        value_key="Avg Duration Seconds",
    )
    top_queries_by_count = parse_top_queries_rows(
        top_queries_payload.get("topByCount") if isinstance(top_queries_payload, dict) else [],
        value_field="count",
        value_key="Count",
    )
    top_queries_by_avg_memory = parse_top_queries_rows(
        top_queries_payload.get("topByAvgMemoryUsage") if isinstance(top_queries_payload, dict) else [],
        value_field="avgMemoryBytes",
        value_key="Avg Memory Usage Bytes",
    )
    high_cardinality_metric_usage = build_high_cardinality_usage_rows(
        top_metrics_by_cardinality,
        top_metric_label_drivers,
        metric_usage_rows,
        snapshot_time,
        args.metric_request_low_threshold,
        args.metric_last_request_old_days,
    )
    global_label_distribution_source = tsdb_results.get("global_metric_series")
    if args.full_label_scan and isinstance(global_label_distribution_source, dict) and global_label_distribution_source:
        label_distribution_source = global_label_distribution_source
        label_distribution_scope = "all metrics"
    else:
        label_distribution_source = tsdb_results.get("metric_series")
        label_distribution_scope = "top metrics"
    label_distribution = build_label_distribution_rows(
        tsdb_results.get("tsdb_stats"),
        label_distribution_source,
        summary.get("active_series"),
        args.label_series_share_threshold,
        args.label_low_metric_coverage_threshold,
        args.label_global_unique_values_threshold,
        allow_global_series_extrapolation=not bool(args.selector.strip()),
    )
    errors = merged_errors(raw_results, tsdb_results)

    report = {
        "report_type": "victoriametrics_analysis_report",
        "generated_at": utc_iso(datetime.now(UTC)),
        "query_time": snapshot_time,
        "victoriametrics_url": args.victoriametrics_url,
        "selectors": {"selector": args.selector},
        "windows": {
            "rate_window": args.rate_window,
            "churn_lookback": args.churn_lookback,
            "storage_eta_lookback": args.storage_eta_lookback,
            "time_offset": args.time_offset,
        },
        "thresholds": {
            "index_to_data_ratio": {"kind": "max", "value": args.index_ratio_warning_threshold},
            "slow_inserts_ratio": {"kind": "max", "value": args.slow_inserts_warning_threshold},
            "remote_write_http_errors_max": {"kind": "max", "value": 0},
            "remote_write_http_error_ratio_max": {"kind": "max", "value": 0.01},
            "remote_write_parser_read_errors_max": {"kind": "max", "value": 0},
            "remote_write_parser_unmarshal_errors_max": {"kind": "max", "value": 0},
            "rows_ignored_too_many_labels_total": {"kind": "max", "value": 0},
            "rows_ignored_too_long_label_name_total": {"kind": "max", "value": 0},
            "rows_ignored_too_long_label_value_total": {"kind": "max", "value": 0},
            "insert_limit_reached_total": {"kind": "max", "value": 50},
            "select_limit_reached_total": {"kind": "max", "value": 50},
            "select_limit_timeout_total": {"kind": "max", "value": 0},
            "container_cpu_cfs_throttled_seconds_rate_max": {"kind": "max", "value": 0.05},
            "process_pressure_cpu_stalled_seconds_rate_max": {"kind": "max", "value": 0.05},
            "vmagent_container_cpu_cfs_throttled_seconds_rate_max": {"kind": "max", "value": 0.05},
            "vmagent_process_pressure_cpu_stalled_seconds_rate_max": {"kind": "max", "value": 0.05},
            "vmagent_persistentqueue_bytes_dropped_ratio": {"kind": "max", "value": 0},
        },
        "usage_thresholds": {
            "metric_request_low_threshold": args.metric_request_low_threshold,
            "metric_last_request_old_days": args.metric_last_request_old_days,
            "label_series_share_threshold": args.label_series_share_threshold,
            "label_low_metric_coverage_threshold": args.label_low_metric_coverage_threshold,
            "label_global_unique_values_threshold": args.label_global_unique_values_threshold,
            "full_label_scan": args.full_label_scan,
        },
        "analysis_modes": {
            "label_distribution_scope": label_distribution_scope,
            "top_metric_label_drivers_partial_sample": top_metric_label_drivers_sampled,
            "service_group_by_labels": service_group_labels,
            "vmalert_requests_query_overridden": bool(args.vmalert_requests_query.strip()),
            "cluster_ingestion_table_enabled": args.enable_cluster_ingestion_table,
            "cluster_ingestion_label": args.cluster_ingestion_label,
            "top_limit": args.top_limit,
            "metric_label_analysis_limit": metric_label_analysis_limit,
            "configured_metric_label_analysis_limit": args.metric_label_analysis_limit,
            "metric_usage_limit": args.metric_usage_limit,
            "series_fetch_workers": args.series_fetch_workers,
            "top_queries_limit": args.top_queries_limit,
            "top_queries_lookback": args.top_queries_lookback,
            "infra_labels": sorted(INFRA_LABELS),
            "label_global_unique_values_threshold": args.label_global_unique_values_threshold,
        },
        "summary": summary,
        "tables": {
            "top_queries_by_sum_duration": top_queries_by_sum_duration,
            "top_queries_by_avg_duration": top_queries_by_avg_duration,
            "top_queries_by_count": top_queries_by_count,
            "top_queries_by_avg_memory": top_queries_by_avg_memory,
            "top_labels_by_unique_values": [
                {
                    "name": row.get("name"),
                    "value": row.get("value"),
                    "__infraLabel": row.get("name") in INFRA_LABELS,
                }
                for row in tsdb_rows(
                    tsdb_results.get("tsdb_stats", {}).get("labelValueCountByLabelName")
                )
                if row.get("name") not in TOP_LABELS_BY_UNIQUE_VALUES_EXCLUDED
            ],
            "vmsingle_configured_flags": vmsingle_configured_flags,
            "vmagent_configured_flags": vmagent_configured_flags,
            "vmsingle_all_effective_flags": vmsingle_effective_flags,
            "vmagent_all_effective_flags": vmagent_effective_flags,
            "top_services_by_series": grouped_sample_rows(
                raw_results.get("top_services_by_series"),
                "Series",
                service_group_labels,
            ),
            "top_services_by_new_series": grouped_sample_rows(
                raw_results.get("top_services_by_new_series"),
                "New Series",
                service_group_labels,
            ),
            "ingestion_by_cluster": (
                table_rows(
                    raw_results.get("ingestion_by_cluster"),
                    args.cluster_ingestion_label,
                    "Rows Per Second",
                )
                if args.enable_cluster_ingestion_table
                else []
            ),
            "tsdb_top_labels_by_memory_bytes": tsdb_rows(
                tsdb_results.get("tsdb_stats", {}).get("memoryInBytesByLabelName")
            ),
            "label_distribution": [
                {
                    **row,
                    "__infraLabel": row.get("Label") in INFRA_LABELS,
                }
                for row in label_distribution
                if not (
                    row.get("Label") in INFRA_LABELS
                    and isinstance(row.get("Unique Values"), (int, float))
                    and float(row.get("Unique Values")) <= 1
                )
            ],
            "high_cardinality_metric_usage": high_cardinality_metric_usage,
        },
        "errors": errors,
        "limitations": [
            "TSDB stats are global cardinality views from /api/v1/status/tsdb.",
            (
                "Metric usage view is read from /api/v1/status/metric_names_stats and joined "
                "with high-cardinality metrics by metric name. Increase METRIC_USAGE_LIMIT "
                "if usage columns are empty for high-cardinality metrics."
            ),
            (
                "Top query tables are read from /api/v1/status/top_queries with the "
                "configured TOP_QUERIES_LIMIT and TOP_QUERIES_LOOKBACK values."
            ),
            (
                "Top metric label-driver analysis uses /api/v1/series per metric and may be "
                "sampled when series count exceeds the configured SERIES_SAMPLE_LIMIT."
            ),
            (
                "Label distribution is calculated from selector-scoped /api/v1/series scans "
                "only. In the default mode it covers top TSDB metrics for speed; with "
                "FULL_LABEL_SCAN it covers all discovered metric names that still have "
                "series in the current selector scope."
                if args.full_label_scan
                else (
                    "Label distribution is calculated only from selector-scoped "
                    "/api/v1/series scans across top TSDB metrics for speed. This makes "
                    "Series, Series Share, Metrics Coverage, and Unique Values consistent "
                    "within the scanned scope, but still heuristic because non-top metrics "
                    "are not scanned unless FULL_LABEL_SCAN is enabled."
                )
            ),
            (
                "If FULL_LABEL_SCAN is requested but aborted by MAX_FULL_SCAN_METRICS, "
                "the report falls back to the already collected top-metrics scan instead "
                "of dropping label distribution completely."
            ),
            (
                "FULL_LABEL_SCAN can still be expensive in time and memory because it "
                "first discovers metric names globally and then stores per-metric "
                "/api/v1/series payloads locally; use MAX_FULL_SCAN_METRICS, "
                "SERIES_FETCH_WORKERS, and GLOBAL_SERIES_FETCH_LIMIT to constrain the "
                "workload on large installations."
            ),
            (
                "If FULL_LABEL_SCAN is enabled and GLOBAL_SERIES_FETCH_LIMIT is greater "
                "than zero, label distribution becomes sample-based for very large metrics."
            ),
            (
                "When SELECTOR is empty, label distribution can extrapolate sampled label "
                "presence by using TSDB seriesCountByMetricName for each metric. This "
                "extrapolation is intentionally disabled for selector-scoped reports, "
                "because TSDB per-metric totals are global and would otherwise mix "
                "foreign series into scoped estimates."
            ),
            (
                "The metric_names_stats payload shape may vary by VictoriaMetrics version; "
                "this script normalizes common field names such as metricName, "
                "queryRequestsCount, and lastRequestTimestamp."
            ),
            (
                "Query Requests Per Second and Avg Query Request Duration Seconds are "
                "peak values over the lookback window, derived from "
                "vm_request_duration_seconds_* over the configured or default request-path regex."
            ),
            (
                "Top Services By Series is derived from the current cardinality instant "
                "query using count(...) by (...); it is a visible series-footprint view, "
                "not vm_cache_entries-based active-series accounting."
            ),
            (
                "Query Concurrency Limit is read from vm_concurrent_select_capacity when "
                "available, with a fallback to the effective search.maxConcurrentRequests "
                "flag value."
            ),
            (
                "VMAlert Requests Per Second is a direct datasource request rate when "
                "vmalert_datasource_queries_total is available; otherwise it uses "
                "vmalert_execution_total as a rule execution rate proxy, not as proof "
                "of exact datasource requests. By default this query is global and is "
                "not automatically scoped by SELECTOR; use VMALERT_REQUESTS_QUERY if "
                "environment-specific scoping is required."
                if not args.vmalert_requests_query.strip()
                else (
                    "VMAlert Requests Per Second is computed from the custom "
                    "VMALERT_REQUESTS_QUERY expression supplied for this environment."
                )
            ),
            (
                "Container-scoped vmsingle/vmagent checks such as flags, CPU pressure, "
                "throttling, persistent-queue ratio, and remote-write traffic are "
                "scoped by container name plus MONITORING_NAMESPACE, not by SELECTOR. "
                "In shared monitoring namespaces, add a more specific deployment-level "
                "selector through custom queries if you need per-stack isolation."
            ),
        ],
        "tsdb_queries": tsdb_queries(args),
        "queries": query_map,
    }
    report["findings"] = findings(report, args)
    return report


def write_report(report: dict[str, Any], output: str) -> None:
    Path(output).write_text(json.dumps(report, ensure_ascii=True, indent=2) + "\n", encoding="utf-8")
