from __future__ import annotations

import sys
from concurrent.futures import ThreadPoolExecutor, as_completed
from datetime import UTC, datetime
from typing import Any

from .client import VictoriaMetricsClient
from .common import (
    SECONDS_IN_DAY,
    maybe_round,
    merge_selector,
    merge_selector_matchers,
    parse_group_by_labels,
    parse_int_value,
    safe_div,
    selector_expr,
    utc_iso,
)


def effective_metric_label_analysis_limit(args) -> int:
    return max(args.metric_label_analysis_limit, args.top_limit)


def fetch_series_by_metric(
    client: VictoriaMetricsClient,
    metric_names: list[str],
    *,
    selector: str,
    limit: int,
    max_workers: int,
    progress_label: str | None = None,
    progress_every: int | None = None,
) -> dict[str, Any]:
    if not metric_names:
        return {}
    worker_count = max(1, min(max_workers, len(metric_names)))

    def load(metric_name: str) -> tuple[str, Any]:
        metric_selector = merge_selector(selector, f'__name__="{metric_name}"')
        return metric_name, client.safe_series(metric_selector, limit)

    results: dict[str, Any] = {}
    completed = 0
    total = len(metric_names)
    with ThreadPoolExecutor(max_workers=worker_count) as executor:
        futures = {executor.submit(load, metric_name): metric_name for metric_name in metric_names}
        for future in as_completed(futures):
            metric_name, value = future.result()
            results[metric_name] = value
            completed += 1
            if progress_label and progress_every and (completed % progress_every == 0 or completed == total):
                print(
                    f"{progress_label}: fetched {completed}/{total} metric series payloads",
                    file=sys.stderr,
                    flush=True,
                )
    return results


def queries(args) -> dict[str, str]:
    selector = selector_expr(args.selector)
    rate_window = args.rate_window
    churn_lookback = args.churn_lookback
    storage_eta_lookback = args.storage_eta_lookback
    vmsingle_cpu_selector = merge_selector_matchers(
        "",
        ['container="vmsingle"'] + ([f'namespace="{args.monitoring_namespace}"'] if args.monitoring_namespace else []),
    )
    vmagent_cpu_selector = merge_selector_matchers(
        "",
        ['container="vmagent"'] + ([f'namespace="{args.monitoring_namespace}"'] if args.monitoring_namespace else []),
    )
    service_group_labels = parse_group_by_labels(args.service_group_by_labels)
    active_group_by_clause = ", ".join(service_group_labels)
    active_present_matchers = [f'{label}!=""' for label in service_group_labels]
    scrape_group_by_clause = ", ".join(service_group_labels)
    scrape_present_matchers = [f'{label}!=""' for label in service_group_labels]
    cardinality_group_selector = merge_selector_matchers(args.cardinality_selector, active_present_matchers)
    scrape_group_selector = merge_selector_matchers(args.scrape_selector, scrape_present_matchers)
    query_requests_selector = merge_selector(args.selector, f'path=~"{args.query_requests_path_regex}"')
    remote_write_selector = merge_selector(args.selector, 'path="/api/v1/write",protocol="promremotewrite"')
    remote_write_parser_selector = merge_selector(args.selector, 'type="promremotewrite"')

    active_series_selector = merge_selector(args.selector, 'type="storage/hour_metric_ids"')
    non_index_selector = merge_selector(args.selector, 'type=~"storage/(big|small)"')
    index_selector = merge_selector(args.selector, 'type="indexdb/file"')
    dedup_selector = merge_selector(args.selector, 'type="merge"')

    result = {
        "active_series": f"sum(vm_cache_entries{active_series_selector})",
        "total_datapoints": f"sum(vm_rows{non_index_selector})",
        "data_size_bytes": f"sum(vm_data_size_bytes{non_index_selector})",
        "indexdb_size_bytes": f"sum(vm_data_size_bytes{index_selector})",
        "bytes_per_sample": f"sum(vm_data_size_bytes{non_index_selector}) / sum(vm_rows{non_index_selector})",
        "index_to_data_ratio": f"sum(vm_data_size_bytes{index_selector}) / sum(vm_data_size_bytes{non_index_selector})",
        "min_free_disk_space_bytes": f"min(vm_free_disk_space_bytes{selector})",
        "total_free_disk_space_bytes": f"sum(vm_free_disk_space_bytes{selector})",
        "storage_full_eta_days": (
            f"sum(vm_free_disk_space_bytes{selector}) / "
            f"sum(rate(vm_data_size_bytes{selector}[{storage_eta_lookback}])) / {SECONDS_IN_DAY}"
        ),
        "churn_rate_per_second": f"sum(rate(vm_new_timeseries_created_total{selector}[{rate_window}]))",
        "new_series_total": f"sum(increase(vm_new_timeseries_created_total{selector}[{churn_lookback}]))",
        "ingestion_rate_per_second": f"sum(rate(vm_rows_inserted_total{selector}[{rate_window}]))",
        "remote_write_requests_max": (
            f"max_over_time((sum(rate(vm_http_requests_total{remote_write_selector}[{rate_window}])))"
            f"[{churn_lookback}:{rate_window}])"
        ),
        "remote_write_http_errors_max": (
            f"max_over_time((sum(rate(vm_http_request_errors_total{remote_write_selector}[{rate_window}])))"
            f"[{churn_lookback}:{rate_window}])"
        ),
        "remote_write_http_error_ratio_max": (
            f"max_over_time(((sum(rate(vm_http_request_errors_total{remote_write_selector}[{rate_window}]))) / "
            f"clamp_min(sum(rate(vm_http_requests_total{remote_write_selector}[{rate_window}])), 1e-12))"
            f"[{churn_lookback}:{rate_window}])"
        ),
        "remote_write_parser_read_errors_max": (
            f"max_over_time((sum(rate(vm_protoparser_read_errors_total{remote_write_parser_selector}[{rate_window}])))"
            f"[{churn_lookback}:{rate_window}])"
        ),
        "remote_write_parser_unmarshal_errors_max": (
            f"max_over_time((sum(rate("
            f"vm_protoparser_unmarshal_errors_total{remote_write_parser_selector}[{rate_window}])))"
            f"[{churn_lookback}:{rate_window}])"
        ),
        "rows_ignored_too_many_labels_total": (
            f"sum(increase(vm_rows_ignored_total"
            f"{merge_selector(args.selector, 'reason=\"too_many_labels\"')}[{churn_lookback}]))"
        ),
        "rows_ignored_too_long_label_name_total": (
            f"sum(increase(vm_rows_ignored_total"
            f"{merge_selector(args.selector, 'reason=\"too_long_label_name\"')}[{churn_lookback}]))"
        ),
        "rows_ignored_too_long_label_value_total": (
            f"sum(increase(vm_rows_ignored_total"
            f"{merge_selector(args.selector, 'reason=\"too_long_label_value\"')}[{churn_lookback}]))"
        ),
        "insert_limit_reached_total": (
            f"sum(increase(vm_concurrent_insert_limit_reached_total{selector}[{churn_lookback}]))"
        ),
        "select_limit_reached_total": (
            f"sum(increase(vm_concurrent_select_limit_reached_total{selector}[{churn_lookback}]))"
        ),
        "select_limit_timeout_total": (
            f"sum(increase(vm_concurrent_select_limit_timeout_total{selector}[{churn_lookback}]))"
        ),
        "slow_inserts_ratio": (
            f"max((sum without(type)(rate(vm_slow_row_inserts_total{selector}[{rate_window}])) / "
            f"clamp_min(sum without(type)(rate(vm_rows_inserted_total{selector}[{rate_window}])), 1e-12)))"
        ),
        "vm_slow_queries_total_rate_per_second": f"sum(rate(vm_slow_queries_total{selector}[{rate_window}]))",
        "container_cpu_cfs_throttled_seconds_rate_max": (
            f"max_over_time((sum(rate("
            f"container_cpu_cfs_throttled_seconds_total{vmsingle_cpu_selector}[{rate_window}])))"
            f"[{churn_lookback}:{rate_window}])"
        ),
        "process_pressure_cpu_stalled_seconds_rate_max": (
            f"max_over_time((sum(rate("
            f"process_pressure_cpu_stalled_seconds_total{vmsingle_cpu_selector}[{rate_window}])))"
            f"[{churn_lookback}:{rate_window}])"
        ),
        "vmagent_container_cpu_cfs_throttled_seconds_rate_max": (
            f"max_over_time((sum(rate(container_cpu_cfs_throttled_seconds_total{vmagent_cpu_selector}[{rate_window}])))"
            f"[{churn_lookback}:{rate_window}])"
        ),
        "vmagent_process_pressure_cpu_stalled_seconds_rate_max": (
            f"max_over_time((sum(rate("
            f"process_pressure_cpu_stalled_seconds_total{vmagent_cpu_selector}[{rate_window}])))"
            f"[{churn_lookback}:{rate_window}])"
        ),
        "vmagent_persistentqueue_bytes_dropped_ratio": (
            f"max(((sum(rate(vm_persistentqueue_bytes_dropped_total"
            f"{vmagent_cpu_selector}[{rate_window}])) by (job, instance)) / "
            f"(sum(rate(vm_persistentqueue_bytes_written_total"
            f"{vmagent_cpu_selector}[{rate_window}])) by (job, instance))) "
            f"and on(job, instance) ((sum(rate(vm_persistentqueue_bytes_written_total"
            f"{vmagent_cpu_selector}[{rate_window}])) by (job, instance)) > 0))"
        ),
        "remote_write_traffic_mbit_per_second": (
            f"sum(rate(vmagent_remotewrite_conn_bytes_written_total"
            f"{vmagent_cpu_selector}[{rate_window}])) * 8 / 1000000"
        ),
        "storage_growth_rows_per_second": (
            f"sum(rate(vm_rows_added_to_storage_total{selector}[{rate_window}])) - "
            f"sum(rate(vm_deduplicated_samples_total{dedup_selector}[{rate_window}]))"
        ),
        "top_services_by_series": (
            f"topk({args.top_limit}, count({cardinality_group_selector}) by ({active_group_by_clause}))"
        ),
        "top_services_by_new_series": (
            f"topk({args.top_limit}, sum(increase(scrape_series_added{scrape_group_selector}[{churn_lookback}])) "
            f"by ({scrape_group_by_clause}))"
        ),
    }
    result["query_requests_rate_per_second"] = (
        f"max_over_time((sum(rate(vm_request_duration_seconds_count{query_requests_selector}[{rate_window}])))"
        f"[{churn_lookback}:{rate_window}])"
    )
    result["avg_query_request_duration_seconds"] = (
        f"max_over_time(((sum(rate(vm_request_duration_seconds_sum{query_requests_selector}[{rate_window}])) / "
        f"clamp_min(sum(rate(vm_request_duration_seconds_count{query_requests_selector}[{rate_window}])), 1e-12)))"
        f"[{churn_lookback}:{rate_window}])"
    )
    result["query_concurrency_limit"] = f"sum(vm_concurrent_select_capacity{selector})"
    result["vmsingle_flags"] = f"flag{vmsingle_cpu_selector}"
    result["vmagent_flags"] = f"flag{vmagent_cpu_selector}"
    if args.enable_cluster_ingestion_table:
        cluster_label = args.cluster_ingestion_label
        cluster_selector = merge_selector(args.selector, f'{cluster_label}!=""')
        result["ingestion_by_cluster"] = (
            f"topk({args.top_limit}, "
            f"sum(rate(vm_rows_inserted_total{cluster_selector}[{rate_window}])) by ({cluster_label}))"
        )
    vmalert_base_query = (
        args.vmalert_requests_query.strip()
        if args.vmalert_requests_query.strip()
        else (
            f"(sum(rate(vmalert_datasource_queries_total[{rate_window}])) "
            f"or sum(rate(vmalert_execution_total[{rate_window}])) "
            "or vector(0))"
        )
    )
    result["vmalert_requests_rate_per_second"] = (
        f"max_over_time(({vmalert_base_query})[{churn_lookback}:{rate_window}])"
    )
    return result


def tsdb_queries(args) -> dict[str, Any]:
    metric_label_analysis_limit = effective_metric_label_analysis_limit(args)
    metric_series_match = merge_selector(args.selector, '__name__="<metric_name>"')
    return {
        "tsdb_stats": {"path": "/api/v1/status/tsdb", "query": {"limit": str(args.top_limit)}},
        "metric_names_stats": {
            "path": "/api/v1/status/metric_names_stats",
            "query": {"limit": str(args.metric_usage_limit)},
        },
        "top_queries": {
            "path": "/api/v1/status/top_queries",
            "query": {"topN": str(args.top_queries_limit), "maxLifetime": args.top_queries_lookback},
        },
        "metric_label_analysis": {
            "path": "/api/v1/series",
            "query_template": {"match[]": metric_series_match, "limit": str(args.series_sample_limit)},
            "top_metrics_limit": metric_label_analysis_limit,
        },
        "global_label_distribution": {
            "metric_names_path": "/api/v1/label/__name__/values",
            "series_query_template": {
                "match[]": metric_series_match,
                "limit": str(args.global_series_fetch_limit),
            },
            "series_limit_per_metric": args.global_series_fetch_limit,
            "enabled": args.full_label_scan,
        },
    }


def metric_error(value: Any) -> str | None:
    if isinstance(value, dict) and "error" in value:
        return str(value["error"])
    return None


def merged_errors(raw_results: dict[str, Any], tsdb_results: dict[str, Any]) -> dict[str, str]:
    result = {name: message for name, value in raw_results.items() if (message := metric_error(value))}
    if message := metric_error(tsdb_results.get("tsdb_stats")):
        result["tsdb_stats"] = message
    if message := metric_error(tsdb_results.get("metric_names_stats")):
        result["metric_names_stats"] = message
    if message := metric_error(tsdb_results.get("top_queries")):
        result["top_queries"] = message
    if message := metric_error(tsdb_results.get("all_metric_names")):
        result["all_metric_names"] = message
    metric_series = tsdb_results.get("metric_series")
    if isinstance(metric_series, dict):
        for metric_name, value in metric_series.items():
            if message := metric_error(value):
                result[f"metric_series:{metric_name}"] = message
    global_metric_series = tsdb_results.get("global_metric_series")
    if isinstance(global_metric_series, dict):
        for metric_name, value in global_metric_series.items():
            if message := metric_error(value):
                result[f"global_metric_series:{metric_name}"] = message
    return result


def collect_results(client: VictoriaMetricsClient, query_map: dict[str, str], *, dry_run: bool) -> dict[str, Any]:
    if dry_run:
        return {"queries": query_map}
    return {name: client.safe_query(expression) for name, expression in query_map.items()}


def collect_tsdb_results(
    client: VictoriaMetricsClient,
    args,
    raw_results: dict[str, Any],
    *,
    dry_run: bool,
    metric_names_stats: Any = None,
) -> dict[str, Any]:
    query_plan = tsdb_queries(args)
    if dry_run:
        return {"queries": query_plan}

    result: dict[str, Any] = {}
    with ThreadPoolExecutor(max_workers=3) as executor:
        futures = {
            "tsdb_stats": executor.submit(client.safe_tsdb_stats, args.top_limit),
            "top_queries": executor.submit(client.safe_top_queries, args.top_queries_limit, args.top_queries_lookback),
        }
        if metric_names_stats is None:
            futures["metric_names_stats"] = executor.submit(client.safe_metric_names_stats, args.metric_usage_limit)
        else:
            result["metric_names_stats"] = metric_names_stats
        for name, future in futures.items():
            result[name] = future.result()

    top_metrics_rows = tsdb_rows(
        result["tsdb_stats"].get("seriesCountByMetricName")
        if isinstance(result["tsdb_stats"], dict)
        else []
    )
    metric_names = [
        row["name"]
        for row in top_metrics_rows[: effective_metric_label_analysis_limit(args)]
        if row.get("name")
    ]
    result["metric_series"] = fetch_series_by_metric(
        client,
        metric_names,
        selector=args.selector,
        limit=args.series_sample_limit,
        max_workers=args.series_fetch_workers,
    )

    if args.full_label_scan:
        result["all_metric_names"] = client.safe_label_values("__name__")
        global_series_by_metric: dict[str, Any] = {
            metric_name: value
            for metric_name, value in result.get("metric_series", {}).items()
            if isinstance(metric_name, str)
        }
        all_metric_names = result.get("all_metric_names")
        if isinstance(all_metric_names, list):
            valid_metric_names = [
                metric_name
                for metric_name in all_metric_names
                if isinstance(metric_name, str) and metric_name
            ]
            if args.max_full_scan_metrics > 0 and len(valid_metric_names) > args.max_full_scan_metrics:
                result["all_metric_names"] = {
                    "error": (
                        f"FULL_LABEL_SCAN aborted: discovered {len(valid_metric_names)} metric names, "
                        f"which exceeds MAX_FULL_SCAN_METRICS={args.max_full_scan_metrics}. "
                        "Increase the cap or disable FULL_LABEL_SCAN."
                    )
                }
            else:
                remaining_metric_names = [
                    metric_name
                    for metric_name in valid_metric_names
                    if metric_name not in global_series_by_metric
                ]
                print(
                    f"FULL_LABEL_SCAN enabled: fetching /api/v1/series for "
                    f"{len(remaining_metric_names)} additional metric names. "
                    "This may take a long time on large installations.",
                    file=sys.stderr,
                    flush=True,
                )
                global_series_by_metric.update(
                    fetch_series_by_metric(
                        client,
                        remaining_metric_names,
                        selector=args.selector,
                        limit=args.global_series_fetch_limit,
                        max_workers=args.series_fetch_workers,
                        progress_label="FULL_LABEL_SCAN",
                        progress_every=250,
                    )
                )
            global_series_by_metric = {
                metric_name: value
                for metric_name, value in global_series_by_metric.items()
                if not isinstance(value, list) or value
            }
        result["global_metric_series"] = global_series_by_metric
    return result


def tsdb_rows(items: Any) -> list[dict[str, Any]]:
    if not isinstance(items, list):
        return []
    rows: list[dict[str, Any]] = []
    for item in items:
        if not isinstance(item, dict):
            continue
        name = item.get("name")
        value = item.get("value")
        if not isinstance(name, str) or not isinstance(value, (int, float)):
            continue
        rows.append({"name": name, "value": value})
    rows.sort(key=lambda row: row["value"], reverse=True)
    return rows


def first_present(item: dict[str, Any], keys: list[str]) -> Any:
    for key in keys:
        if key in item and item.get(key) is not None:
            return item.get(key)
    return None


def parse_timestamp(value: Any) -> str | None:
    if value in (None, ""):
        return None
    if isinstance(value, str):
        return value
    if isinstance(value, (int, float)):
        timestamp = float(value)
        if timestamp > 1e17:
            timestamp /= 1_000_000_000.0
        elif timestamp > 1e14:
            timestamp /= 1_000_000.0
        elif timestamp > 1e11:
            timestamp /= 1000.0
        return utc_iso(datetime.fromtimestamp(timestamp, tz=UTC))
    return str(value)


def compact_timestamp(value: str | None) -> str | None:
    if value is None:
        return None
    if value.endswith(".000Z"):
        return value[:-5] + "Z"
    return value


def top_queries_time_range(seconds: Any) -> str:
    if not isinstance(seconds, (int, float)):
        return "-"
    total_seconds = int(seconds)
    if total_seconds <= 0:
        return "-"
    units = (
        (86400, "d"),
        (3600, "h"),
        (60, "m"),
        (1, "s"),
    )
    parts: list[str] = []
    remainder = total_seconds
    for unit_seconds, suffix in units:
        if remainder >= unit_seconds:
            value, remainder = divmod(remainder, unit_seconds)
            parts.append(f"{value}{suffix}")
        if len(parts) == 2:
            break
    return "".join(parts) if parts else "-"


def parse_top_queries_rows(items: Any, *, value_field: str, value_key: str) -> list[dict[str, Any]]:
    if not isinstance(items, list):
        return []
    rows: list[dict[str, Any]] = []
    for item in items:
        if not isinstance(item, dict):
            continue
        query = item.get("query")
        if not isinstance(query, str) or not query:
            continue
        value = item.get(value_field)
        count = item.get("count")
        time_range_seconds = item.get("timeRangeSeconds")
        row = {
            "Query": query,
            value_key: value if isinstance(value, (int, float)) else None,
            "Query Time Interval": top_queries_time_range(time_range_seconds),
            "Count": parse_int_value(count),
        }
        rows.append(row)
    rows.sort(
        key=lambda row: row.get(value_key) if isinstance(row.get(value_key), (int, float)) else -1,
        reverse=True,
    )
    return rows


def metric_usage_items(items: Any) -> list[dict[str, Any]]:
    if isinstance(items, dict):
        nested = first_present(
            items,
            [
                "stats",
                "metricNamesStats",
                "metric_names_stats",
                "rows",
                "result",
                "data",
                "items",
                "metrics",
                "records",
            ],
        )
        if nested is not None and nested is not items:
            return metric_usage_items(nested)
        rows: list[dict[str, Any]] = []
        for metric_name, value in items.items():
            if not isinstance(metric_name, str):
                continue
            if isinstance(value, dict):
                row = dict(value)
                row.setdefault("metricName", metric_name)
                rows.append(row)
        return rows
    if isinstance(items, list):
        return [item for item in items if isinstance(item, dict)]
    return []


def parse_metric_usage_rows(items: Any) -> list[dict[str, Any]]:
    rows: list[dict[str, Any]] = []
    for item in metric_usage_items(items):
        metric_name = first_present(item, ["metricName", "metric", "__name__", "name"])
        if not isinstance(metric_name, str) or not metric_name:
            continue
        query_requests_count = first_present(
            item,
            [
                "queryRequestsCount",
                "requestsCount",
                "requestCount",
                "requests",
                "query_requests_count",
                "requests_count",
                "request_count",
            ],
        )
        last_request_timestamp = first_present(
            item,
            [
                "lastRequestTimestamp",
                "lastQueryRequestTimestamp",
                "lastQueryTimestamp",
                "lastRequestTime",
                "lastRequestTs",
                "lastQueryRequestTs",
                "last_request_timestamp",
                "last_query_request_timestamp",
                "last_query_timestamp",
                "last_request_time",
                "last_request_ts",
                "last_query_request_ts",
            ],
        )
        rows.append(
            {
                "metric": metric_name,
                "queryRequestsCount": parse_int_value(query_requests_count),
                "lastRequestTimestamp": compact_timestamp(parse_timestamp(last_request_timestamp)),
            }
        )
    rows.sort(
        key=lambda row: (
            row["queryRequestsCount"] if isinstance(row.get("queryRequestsCount"), (int, float)) else -1
        ),
        reverse=True,
    )
    return rows


def parse_iso_timestamp(value: Any) -> datetime | None:
    if not isinstance(value, str) or not value:
        return None
    normalized = value.replace("Z", "+00:00")
    try:
        return datetime.fromisoformat(normalized)
    except ValueError:
        return None


def metric_usage_age_days(value: Any, query_time_value: str) -> float | None:
    last_requested = parse_iso_timestamp(value)
    query_time_dt = parse_iso_timestamp(query_time_value)
    if last_requested is None or query_time_dt is None:
        return None
    return round((query_time_dt - last_requested).total_seconds() / SECONDS_IN_DAY, 3)


def build_high_cardinality_usage_rows(
    top_metrics: list[dict[str, Any]],
    metric_drivers: list[dict[str, Any]],
    metric_usage: list[dict[str, Any]],
    query_time_value: str,
    low_request_threshold: float,
    old_request_days: int,
) -> list[dict[str, Any]]:
    drivers_by_metric = {row["metric"]: row for row in metric_drivers if row.get("metric")}
    usage_by_metric = {row["metric"]: row for row in metric_usage if row.get("metric")}
    rows: list[dict[str, Any]] = []
    for top_metric in top_metrics:
        metric_name = top_metric.get("__name__")
        if not isinstance(metric_name, str) or not metric_name:
            continue
        driver = drivers_by_metric.get(metric_name, {})
        usage = usage_by_metric.get(metric_name, {})
        query_requests_count = usage.get("queryRequestsCount")
        last_request_timestamp = usage.get("lastRequestTimestamp")
        last_request_age_days = metric_usage_age_days(last_request_timestamp, query_time_value)
        low_usage = isinstance(query_requests_count, (int, float)) and query_requests_count <= low_request_threshold
        old_usage = isinstance(last_request_age_days, (int, float)) and last_request_age_days >= old_request_days
        rows.append(
            {
                "metric": metric_name,
                "series": top_metric.get("series"),
                "top_label": driver.get("top_label"),
                "top_label_unique_values": driver.get("unique_values"),
                "queryRequestsCount": query_requests_count,
                "lastRequestTimestamp": last_request_timestamp,
                "cleanupCandidate": bool(low_usage or old_usage),
            }
        )
    rows.sort(key=lambda row: row.get("series") if isinstance(row.get("series"), (int, float)) else 0, reverse=True)
    return rows


def metric_label_driver_rows(tsdb_stats: Any, metric_series: Any) -> list[dict[str, Any]]:
    series_count_by_metric = {
        row["name"]: int(row["value"])
        for row in tsdb_rows(tsdb_stats.get("seriesCountByMetricName") if isinstance(tsdb_stats, dict) else [])
    }
    rows: list[dict[str, Any]] = []
    if not isinstance(metric_series, dict):
        return rows
    for metric_name, series_list in metric_series.items():
        if not isinstance(series_list, list):
            continue
        unique_by_label: dict[str, set[str]] = {}
        for series in series_list:
            if not isinstance(series, dict):
                continue
            for label_name, label_value in series.items():
                if label_name == "__name__":
                    continue
                unique_by_label.setdefault(label_name, set()).add(str(label_value))
        if not unique_by_label:
            continue
        top_label, top_values = max(unique_by_label.items(), key=lambda item: len(item[1]))
        total_series = series_count_by_metric.get(metric_name, len(series_list))
        sampled_series = len(series_list)
        sampled = total_series > sampled_series
        rows.append(
            {
                "metric": metric_name,
                "series": total_series,
                "top_label": top_label,
                "unique_values": len(top_values),
                "sampled_series": sampled_series,
                "sampled": sampled,
            }
        )
    rows.sort(key=lambda row: row["series"], reverse=True)
    return rows


def build_label_distribution_rows(
    tsdb_stats: Any,
    metric_series: Any,
    active_series: float | None,
    series_share_threshold: float,
    low_metric_coverage_threshold: float,
    global_unique_values_threshold: int,
    *,
    allow_global_series_extrapolation: bool = False,
) -> list[dict[str, Any]]:
    infra_labels = {
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
    normal_share_threshold = 0.10
    observe_share_threshold = 0.10
    observe_coverage_threshold = low_metric_coverage_threshold
    observe_unique_values_threshold = global_unique_values_threshold
    review_share_threshold = 0.20
    review_coverage_threshold = min(low_metric_coverage_threshold, 0.15)
    review_unique_values_threshold = 1000

    if not isinstance(metric_series, dict):
        return []

    series_count_by_metric = (
        {
            row["name"]: float(row["value"])
            for row in tsdb_rows(tsdb_stats.get("seriesCountByMetricName") if isinstance(tsdb_stats, dict) else [])
        }
        if allow_global_series_extrapolation
        else {}
    )
    inspected_metrics = [
        metric_name for metric_name, series_list in metric_series.items() if isinstance(series_list, list)
    ]
    inspected_metric_count = len(inspected_metrics)
    metrics_by_label: dict[str, set[str]] = {}
    local_series_by_label: dict[str, float] = {}
    local_unique_values_by_label: dict[str, set[str]] = {}
    for metric_name, series_list in metric_series.items():
        if not isinstance(series_list, list):
            continue
        total_metric_series = series_count_by_metric.get(metric_name, float(len(series_list)))
        sampled_series = len(series_list)
        label_presence_count: dict[str, int] = {}
        for series in series_list:
            if not isinstance(series, dict):
                continue
            seen_labels_for_series: set[str] = set()
            for label_name in series:
                if label_name == "__name__":
                    continue
                seen_labels_for_series.add(label_name)
                local_unique_values_by_label.setdefault(label_name, set()).add(str(series[label_name]))
            for label_name in seen_labels_for_series:
                label_presence_count[label_name] = label_presence_count.get(label_name, 0) + 1
        for label_name, count in label_presence_count.items():
            metrics_by_label.setdefault(label_name, set()).add(metric_name)
            estimated_series = (
                total_metric_series * (count / sampled_series)
                if sampled_series > 0 and total_metric_series > sampled_series
                else float(count)
            )
            local_series_by_label[label_name] = local_series_by_label.get(label_name, 0.0) + estimated_series

    rows: list[dict[str, Any]] = []
    candidate_labels = sorted(set(metrics_by_label) - {"__name__"})
    for label_name in candidate_labels:
        unique_values = len(local_unique_values_by_label.get(label_name, set()))
        series = local_series_by_label.get(label_name, 0.0)
        raw_series_share = safe_div(float(series), active_series) if active_series else None
        series_share = (
            maybe_round(min(float(raw_series_share), 1.0), 6)
            if isinstance(raw_series_share, (int, float))
            else None
        )
        if not isinstance(series_share, float) or series_share < series_share_threshold:
            continue
        if label_name in infra_labels and isinstance(unique_values, (int, float)) and float(unique_values) <= 1:
            continue
        metric_count = len(metrics_by_label.get(label_name, set()))
        metric_coverage_ratio = (
            maybe_round(safe_div(float(metric_count), float(inspected_metric_count)), 6)
            if inspected_metric_count
            else None
        )
        if (
            metric_coverage_ratio == 1
            and label_name not in infra_labels
            and not (
                isinstance(unique_values, (int, float))
                and float(unique_values) >= observe_unique_values_threshold
            )
        ):
            continue
        classification = "normal"
        if (
            isinstance(series_share, float)
            and isinstance(metric_coverage_ratio, float)
            and series_share >= review_share_threshold
            and metric_coverage_ratio <= review_coverage_threshold
            and isinstance(unique_values, (int, float))
            and float(unique_values) >= review_unique_values_threshold
        ):
            classification = "review"
        elif (
            isinstance(series_share, float)
            and isinstance(metric_coverage_ratio, float)
            and series_share >= observe_share_threshold
            and metric_coverage_ratio < observe_coverage_threshold
        ):
            classification = "observe"
        elif (
            isinstance(unique_values, (int, float))
            and float(unique_values) >= observe_unique_values_threshold
            and isinstance(metric_coverage_ratio, float)
            and metric_coverage_ratio <= 1
        ):
            classification = "observe"
        elif (
            isinstance(series_share, float)
            and series_share >= normal_share_threshold
            and isinstance(metric_coverage_ratio, float)
            and metric_coverage_ratio < 1
        ):
            classification = "observe"
        rows.append(
            {
                "Label": label_name,
                "Series": round(series, 3),
                "Series Share": series_share,
                "Metrics Coverage": metric_coverage_ratio,
                "Unique Values": unique_values,
                "Classification": classification,
            }
        )
    rows.sort(key=lambda item: item["Series"], reverse=True)
    return rows
