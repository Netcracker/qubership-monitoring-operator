from __future__ import annotations

import math
import re
from datetime import UTC, datetime, timedelta
from typing import Any

BYTES_IN_GB = 1024**3
SECONDS_IN_DAY = 86400


def utc_iso(value: datetime) -> str:
    return value.isoformat(timespec="milliseconds").replace("+00:00", "Z")


def query_time(time_offset: str, duration_seconds_fn) -> str:
    offset_seconds = duration_seconds_fn(time_offset, allow_zero=True)
    return utc_iso(datetime.now(UTC) - timedelta(seconds=offset_seconds))


def selector_expr(matchers: str) -> str:
    stripped = matchers.strip()
    return "{" + stripped + "}" if stripped else ""


def merge_selector(base_metric_selector: str, extra_matcher: str) -> str:
    stripped = base_metric_selector.strip()
    extra = extra_matcher.strip()
    if stripped and extra:
        return "{" + stripped + "," + extra + "}"
    if stripped:
        return "{" + stripped + "}"
    if extra:
        return "{" + extra + "}"
    return ""


def merge_selector_matchers(base_metric_selector: str, extra_matchers: list[str]) -> str:
    matchers = [base_metric_selector.strip()] if base_metric_selector.strip() else []
    matchers.extend(item.strip() for item in extra_matchers if item.strip())
    return "{" + ",".join(matchers) + "}" if matchers else ""


def bytes_to_gb(value: float | None) -> float | None:
    if value is None:
        return None
    return round(value / BYTES_IN_GB, 6)


def maybe_round(value: float | None, digits: int = 6) -> float | None:
    if value is None:
        return None
    return round(value, digits)


def positive_finite(value: float | None) -> float | None:
    if value is None or not math.isfinite(value) or value <= 0:
        return None
    return value


def sample_value(samples: Any) -> float | None:
    if not isinstance(samples, list) or not samples:
        return None
    try:
        return float(samples[0]["value"][1])
    except (KeyError, IndexError, TypeError, ValueError):
        return None


def table_rows(samples: Any, label_column: str, value_name: str) -> list[dict[str, Any]]:
    if not isinstance(samples, list):
        return []
    rows: list[dict[str, Any]] = []
    for sample in samples:
        metric = sample.get("metric", {})
        try:
            value = float(sample["value"][1])
        except (KeyError, IndexError, TypeError, ValueError):
            continue
        rows.append({label_column: metric.get(label_column, ""), value_name: value})
    rows.sort(key=lambda row: row[value_name], reverse=True)
    return rows


def parse_group_by_labels(value: str) -> list[str]:
    return [item.strip() for item in value.split(",") if item.strip()]


def format_group_value(metric: dict[str, Any], label_names: list[str]) -> str:
    parts = []
    for label_name in label_names:
        label_value = metric.get(label_name)
        parts.append(f"{label_name}={label_value}" if label_value not in (None, "") else f"{label_name}=<missing>")
    return ", ".join(parts)


def grouped_sample_rows(samples: Any, value_name: str, label_names: list[str], group_column: str = "Service") -> list[dict[str, Any]]:
    if not isinstance(samples, list):
        return []
    rows: list[dict[str, Any]] = []
    for sample in samples:
        metric = sample.get("metric", {})
        if not isinstance(metric, dict):
            continue
        try:
            value = float(sample["value"][1])
        except (KeyError, IndexError, TypeError, ValueError):
            continue
        rows.append({group_column: format_group_value(metric, label_names), value_name: value})
    rows.sort(key=lambda row: row[value_name], reverse=True)
    return rows


def parse_flag_value(value: Any) -> Any:
    if not isinstance(value, str):
        return value
    stripped = value.strip()
    if stripped == "":
        return ""
    lowered = stripped.lower()
    if lowered in {"true", "false"}:
        return lowered == "true"
    try:
        number = float(stripped)
    except ValueError:
        return stripped
    if math.isfinite(number) and number.is_integer():
        return int(number)
    return number


def flag_rows(samples: Any) -> list[dict[str, Any]]:
    if not isinstance(samples, list):
        return []
    rows: list[dict[str, Any]] = []
    for sample in samples:
        metric = sample.get("metric", {})
        if not isinstance(metric, dict):
            continue
        flag_name = metric.get("name")
        if not isinstance(flag_name, str) or not flag_name:
            continue
        rows.append(
            {
                "Flag": flag_name,
                "Value": parse_flag_value(metric.get("value", "")),
                "Is Set": parse_flag_value(metric.get("is_set", "")),
            }
        )
    rows.sort(key=lambda row: row["Flag"])
    return rows


def filter_flag_rows(rows: list[dict[str, Any]], *, is_set: bool | None = None, require_value: bool = False) -> list[dict[str, Any]]:
    result: list[dict[str, Any]] = []
    for row in rows:
        if is_set is not None and row.get("Is Set") is not is_set:
            continue
        value = row.get("Value")
        if require_value and value in (None, ""):
            continue
        result.append(row)
    return result


def find_flag_value(rows: list[dict[str, Any]], *flag_names: str) -> Any:
    lookup = {row.get("Flag"): row.get("Value") for row in rows if row.get("Flag")}
    for flag_name in flag_names:
        if flag_name in lookup:
            return lookup[flag_name]
    return None


def sum_flag_values(rows: list[dict[str, Any]], *flag_names: str) -> float | int | None:
    names = set(flag_names)
    values: list[float] = []
    for row in rows:
        if row.get("Flag") not in names:
            continue
        value = row.get("Value")
        if isinstance(value, bool):
            continue
        if isinstance(value, (int, float)) and math.isfinite(float(value)):
            values.append(float(value))
    if not values:
        return None
    total = sum(values)
    return int(total) if float(total).is_integer() else total


def safe_div(numerator: float | None, denominator: float | None) -> float | None:
    if numerator is None or denominator in (None, 0):
        return None
    return numerator / denominator


def parse_int_value(value: Any) -> int | Any:
    if isinstance(value, bool):
        return value
    if isinstance(value, int):
        return value
    if isinstance(value, float) and math.isfinite(value):
        return int(value)
    if isinstance(value, str):
        stripped = value.strip()
        if re.fullmatch(r"[0-9]+", stripped):
            return int(stripped)
    return value
