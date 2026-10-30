#!/usr/bin/env python3
"""Collect storage optimization statistics from VictoriaMetrics."""

from __future__ import annotations

import sys

from vmreport.cli import duration_seconds, parser, validate_args
from vmreport.client import HttpClient, VictoriaMetricsClient
from vmreport.common import query_time
from vmreport.data import collect_results, collect_tsdb_results, queries
from vmreport.html import write_html_report
from vmreport.report import build_report, write_report


def main() -> int:
    args = parser(__doc__).parse_args()
    validate_args(args)
    if args.full_label_scan and not args.dry_run:
        cap_suffix = (
            f" Current MAX_FULL_SCAN_METRICS hard cap: {args.max_full_scan_metrics}."
            if args.max_full_scan_metrics > 0
            else " No MAX_FULL_SCAN_METRICS hard cap is configured."
        )
        print(
            "FULL_LABEL_SCAN is enabled. The script will request all metric names and then fetch "
            "/api/v1/series for each of them, which can be slow on large installations."
            + cap_suffix,
            file=sys.stderr,
            flush=True,
        )
    snapshot_time = query_time(args.time_offset, duration_seconds)
    query_map = queries(args)
    client = VictoriaMetricsClient(
        HttpClient(
            args.victoriametrics_url,
            user=args.vm_user,
            password=args.vm_pass,
            insecure_skip_verify=args.insecure_skip_verify,
        ),
        snapshot_time,
    )
    metric_names_stats = client.safe_metric_names_stats(args.metric_usage_limit) if not args.dry_run else None
    raw_results = collect_results(client, query_map, dry_run=args.dry_run)
    tsdb_results = collect_tsdb_results(
        client,
        args,
        raw_results,
        dry_run=args.dry_run,
        metric_names_stats=metric_names_stats,
    )
    report = build_report(args, snapshot_time, raw_results, tsdb_results, query_map, client)
    write_report(report, args.output)
    if args.html_output:
        write_html_report(report, args.html_output)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
