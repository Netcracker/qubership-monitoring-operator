from __future__ import annotations

import argparse
import math
import os
import re

DURATION_PATTERN = re.compile(r"^(?:0|[1-9][0-9]*)(?:s|m|h|d|w)$")


def env_value(name: str, default: str = "") -> str:
    value = os.getenv(name)
    if value is None:
        return default
    if value == "" and default:
        return default
    return value


def env_required(name: str) -> bool:
    return not env_value(name)


def env_bool(name: str, default: bool = False) -> bool:
    default_value = "true" if default else "false"
    return env_value(name, default_value).lower() in ("1", "true", "yes", "on")


def positive_int(value: str) -> int:
    number = int(value)
    if number < 1:
        raise argparse.ArgumentTypeError("value must be greater than zero")
    return number


def env_positive_int(name: str, default: str) -> int:
    return positive_int(env_value(name, default))


def label_name(value: str) -> str:
    stripped = value.strip()
    if not re.fullmatch(r"[a-zA-Z_][a-zA-Z0-9_]*", stripped):
        raise argparse.ArgumentTypeError(f"invalid label name: {value}")
    return stripped


def non_negative_int(value: str) -> int:
    number = int(value)
    if number < 0:
        raise argparse.ArgumentTypeError("value must be greater than or equal to zero")
    return number


def env_non_negative_int(name: str, default: str) -> int:
    return non_negative_int(env_value(name, default))


def ratio(value: str) -> float:
    number = float(value)
    if not math.isfinite(number):
        raise argparse.ArgumentTypeError("value must be a finite number")
    if number < 0:
        raise argparse.ArgumentTypeError("value must be greater than or equal to zero")
    return number


def env_ratio(name: str, default: str) -> float:
    return ratio(env_value(name, default))


def duration_seconds(value: str, *, allow_zero: bool = False) -> int:
    if not DURATION_PATTERN.fullmatch(value):
        raise argparse.ArgumentTypeError("duration must look like 5m, 1h, or 7d")
    amount = int(value[:-1])
    if amount == 0 and not allow_zero:
        raise argparse.ArgumentTypeError("duration must be greater than zero")
    multiplier = {"s": 1, "m": 60, "h": 3600, "d": 86400, "w": 604800}[value[-1]]
    return amount * multiplier


def output_path(value: str) -> str:
    if value == "-":
        raise argparse.ArgumentTypeError("output file path is required; stdout output is not supported")
    return value


def csv_labels(value: str) -> str:
    labels = [item.strip() for item in value.split(",") if item.strip()]
    if not labels:
        raise argparse.ArgumentTypeError("at least one group-by label is required")
    invalid = [label for label in labels if not re.fullmatch(r"[a-zA-Z_][a-zA-Z0-9_]*", label)]
    if invalid:
        raise argparse.ArgumentTypeError(f"invalid label names: {', '.join(invalid)}")
    return ",".join(labels)


def parser(description: str) -> argparse.ArgumentParser:
    result = argparse.ArgumentParser(description=description)
    result.add_argument(
        "--victoriametrics-url",
        default=env_value("VICTORIAMETRICS_URL"),
        required=env_required("VICTORIAMETRICS_URL"),
        help="VictoriaMetrics base URL. Env: VICTORIAMETRICS_URL.",
    )
    result.add_argument("--vm-user", default=env_value("VM_USER"), help="VictoriaMetrics Basic Auth username.")
    result.add_argument("--vm-pass", default=env_value("VM_PASS"), help="VictoriaMetrics Basic Auth password.")
    result.add_argument(
        "--selector",
        default=env_value("SELECTOR"),
        help=(
            "Label matchers applied to VictoriaMetrics self-monitoring metrics, "
            'for example cluster="prod",job=~"vmsingle". Env: SELECTOR.'
        ),
    )
    result.add_argument(
        "--cardinality-selector",
        default=env_value("CARDINALITY_SELECTOR"),
        help=(
            "Optional label matchers for the high-cardinality metrics query. "
            'For example cluster="prod". Env: CARDINALITY_SELECTOR.'
        ),
    )
    result.add_argument(
        "--scrape-selector",
        default=env_value("SCRAPE_SELECTOR"),
        help=(
            "Optional label matchers for scrape_series_added queries. "
            'For example namespace="monitoring" or job!="vmagent-k8s". Env: SCRAPE_SELECTOR.'
        ),
    )
    result.add_argument(
        "--monitoring-namespace",
        default=env_value("MONITORING_NAMESPACE"),
        help=(
            "Namespace used for vmsingle and vmagent CPU-pressure checks. "
            'These queries are evaluated with selectors container="vmsingle",namespace="<value>" '
            'and container="vmagent",namespace="<value>". Env: MONITORING_NAMESPACE.'
        ),
    )
    result.add_argument(
        "--service-group-by-labels",
        type=csv_labels,
        default=csv_labels(env_value("SERVICE_GROUP_BY_LABELS", "namespace,job")),
        help=(
            "Comma-separated label names used to group stored active-series views, "
            "for example namespace,app or namespace,job. Default: namespace,job. "
            "Env: SERVICE_GROUP_BY_LABELS."
        ),
    )
    result.add_argument(
        "--enable-cluster-ingestion-table",
        action="store_true",
        default=env_bool("ENABLE_CLUSTER_INGESTION_TABLE", False),
        help="Include optional ingestion breakdown by cluster label. Env: ENABLE_CLUSTER_INGESTION_TABLE.",
    )
    result.add_argument(
        "--cluster-ingestion-label",
        type=label_name,
        default=env_value("CLUSTER_INGESTION_LABEL", "cluster"),
        help="Label used by the optional ingestion breakdown table. Default: cluster. Env: CLUSTER_INGESTION_LABEL.",
    )
    result.add_argument(
        "--query-requests-path-regex",
        default=env_value(
            "QUERY_REQUESTS_PATH_REGEX",
            ".*/api/v1/(query|query_range|series|labels|label/.+/values|status/tsdb|status/metric_names_stats)$",
        ),
        help=(
            "Regex for VictoriaMetrics request paths counted as metric/query API traffic. "
            "Env: QUERY_REQUESTS_PATH_REGEX."
        ),
    )
    result.add_argument(
        "--vmalert-requests-query",
        default=env_value("VMALERT_REQUESTS_QUERY"),
        help=(
            "Optional PromQL override for vmalert datasource requests per second. "
            "Default behavior uses vmalert_datasource_queries_total, falls back to vmalert_execution_total, "
            "and then to 0 when neither metric is present. "
            "Env: VMALERT_REQUESTS_QUERY."
        ),
    )
    result.add_argument(
        "--time-offset",
        default=env_value("TIME_OFFSET", "0s"),
        help=(
            "Shift the report time into the past, for example 2h means query the snapshot "
            "at now-2h. Env: TIME_OFFSET."
        ),
    )
    result.add_argument(
        "--rate-window",
        default=env_value("RATE_WINDOW", "5m"),
        help="Range window used for rate() queries. Default: 5m. Env: RATE_WINDOW.",
    )
    result.add_argument(
        "--churn-lookback",
        default=env_value("CHURN_LOOKBACK", "24h"),
        help="Range window used for new-series totals. Default: 24h. Env: CHURN_LOOKBACK.",
    )
    result.add_argument(
        "--storage-eta-lookback",
        default=env_value("STORAGE_ETA_LOOKBACK", "24h"),
        help=(
            "Range window for disk-full ETA based on vm_data_size_bytes growth. "
            "Default: 24h. Env: STORAGE_ETA_LOOKBACK."
        ),
    )
    result.add_argument(
        "--top-limit",
        type=positive_int,
        default=env_positive_int("TOP_LIMIT", "10"),
        help="Top N row limit. Default: 10. Env: TOP_LIMIT.",
    )
    result.add_argument(
        "--top-queries-limit",
        type=positive_int,
        default=env_positive_int("TOP_QUERIES_LIMIT", "10"),
        help="Top N row limit for /api/v1/status/top_queries tables. Default: 10. Env: TOP_QUERIES_LIMIT.",
    )
    result.add_argument(
        "--top-queries-lookback",
        default=env_value("TOP_QUERIES_LOOKBACK", "24h"),
        help=(
            "Lookback window passed as maxLifetime to /api/v1/status/top_queries. "
            "Default: 24h. Env: TOP_QUERIES_LOOKBACK."
        ),
    )
    result.add_argument(
        "--metric-label-analysis-limit",
        type=positive_int,
        default=env_positive_int("METRIC_LABEL_ANALYSIS_LIMIT", "10"),
        help=(
            "How many top TSDB metrics to inspect label cardinality for. "
            "Default: 10. Env: METRIC_LABEL_ANALYSIS_LIMIT."
        ),
    )
    result.add_argument(
        "--metric-usage-limit",
        type=positive_int,
        default=env_positive_int("METRIC_USAGE_LIMIT", "5000"),
        help=(
            "How many metric names to request from /api/v1/status/metric_names_stats. "
            "Default: 5000. Env: METRIC_USAGE_LIMIT."
        ),
    )
    result.add_argument(
        "--series-sample-limit",
        type=positive_int,
        default=env_positive_int("SERIES_SAMPLE_LIMIT", "20000"),
        help=(
            "Maximum number of series fetched from /api/v1/series for one metric during label-driver analysis. "
            "Default: 20000. Env: SERIES_SAMPLE_LIMIT."
        ),
    )
    result.add_argument(
        "--global-series-fetch-limit",
        type=non_negative_int,
        default=env_non_negative_int("GLOBAL_SERIES_FETCH_LIMIT", "0"),
        help=(
            "Maximum number of series fetched from /api/v1/series per metric "
            "during global label distribution analysis. "
            "Use 0 for no explicit limit. Default: 0. Env: GLOBAL_SERIES_FETCH_LIMIT."
        ),
    )
    result.add_argument(
        "--full-label-scan",
        action="store_true",
        default=env_bool("FULL_LABEL_SCAN"),
        help="Scan all metric names for label distribution instead of only top TSDB metrics. Env: FULL_LABEL_SCAN.",
    )
    result.add_argument(
        "--max-full-scan-metrics",
        type=non_negative_int,
        default=env_non_negative_int("MAX_FULL_SCAN_METRICS", "0"),
        help=(
            "Optional hard cap for FULL_LABEL_SCAN metric names. Use 0 to disable the cap. "
            "If the discovered metric count exceeds this value, the full scan is skipped with an error. "
            "Default: 0. Env: MAX_FULL_SCAN_METRICS."
        ),
    )
    result.add_argument(
        "--series-fetch-workers",
        type=positive_int,
        default=env_positive_int("SERIES_FETCH_WORKERS", "8"),
        help=(
            "Maximum number of parallel workers used for /api/v1/series and missing label-values fetches. "
            "Default: 8. Env: SERIES_FETCH_WORKERS."
        ),
    )
    result.add_argument(
        "--metric-request-low-threshold",
        type=ratio,
        default=env_ratio("METRIC_REQUEST_LOW_THRESHOLD", "0"),
        help="Treat metrics with queryRequestsCount at or below this value as low-usage. Default: 0.",
    )
    result.add_argument(
        "--metric-last-request-old-days",
        type=positive_int,
        default=env_positive_int("METRIC_LAST_REQUEST_OLD_DAYS", "30"),
        help="Treat metrics not requested for this many days as old-usage candidates. Default: 30.",
    )
    result.add_argument(
        "--label-series-share-threshold",
        type=ratio,
        default=env_ratio("LABEL_SERIES_SHARE_THRESHOLD", "0.05"),
        help="Inspect labels occupying at least this estimated share of active series. Default: 0.05.",
    )
    result.add_argument(
        "--label-low-metric-coverage-threshold",
        type=ratio,
        default=env_ratio("LABEL_LOW_METRIC_COVERAGE_THRESHOLD", "0.3"),
        help="Treat labels present in at most this share of inspected metrics as localized. Default: 0.3.",
    )
    result.add_argument(
        "--label-global-unique-values-threshold",
        type=positive_int,
        default=env_positive_int("LABEL_GLOBAL_UNIQUE_VALUES_THRESHOLD", "100"),
        help=(
            "Keep globally present non-infra labels in Label Distribution "
            "when their Unique Values reach this threshold. "
            "Default: 100. Env: LABEL_GLOBAL_UNIQUE_VALUES_THRESHOLD."
        ),
    )
    result.add_argument(
        "--index-ratio-warning-threshold",
        type=ratio,
        default=env_ratio("INDEX_RATIO_WARNING_THRESHOLD", "1.0"),
        help="Warn when indexdb/file to storage data ratio reaches this value. Default: 1.0.",
    )
    result.add_argument(
        "--slow-inserts-warning-threshold",
        type=ratio,
        default=env_ratio("SLOW_INSERTS_WARNING_THRESHOLD", "0.02"),
        help="Warn when max slow insert ratio reaches this value. Default: 0.02.",
    )
    result.add_argument(
        "--output",
        type=output_path,
        default=env_value("OUTPUT"),
        required=env_required("OUTPUT"),
        help="Required JSON output file path. Env: OUTPUT.",
    )
    result.add_argument(
        "--html-output",
        default=env_value("HTML_OUTPUT"),
        help="Optional self-contained HTML report output path. Env: HTML_OUTPUT.",
    )
    result.add_argument("--insecure-skip-verify", action="store_true", default=env_bool("INSECURE_SKIP_VERIFY"))
    result.add_argument(
        "--dry-run",
        action="store_true",
        default=env_bool("DRY_RUN"),
        help="Write generated queries without calling VictoriaMetrics.",
    )
    return result


def validate_args(args: argparse.Namespace) -> None:
    duration_seconds(args.time_offset, allow_zero=True)
    duration_seconds(args.rate_window)
    duration_seconds(args.churn_lookback)
    duration_seconds(args.storage_eta_lookback)
    duration_seconds(args.top_queries_lookback)
    args.cluster_ingestion_label = label_name(args.cluster_ingestion_label)
