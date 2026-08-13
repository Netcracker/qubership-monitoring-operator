#!/usr/bin/env bash

set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repository_root="$(cd "${script_dir}/.." && pwd)"
chart_root="charts/qubership-monitoring-operator/charts"
failures=0

renames=(
    "blackbox-exporter:blackboxExporter"
    "cert-exporter:certExporter"
    "cloud-events-exporter:cloudEventsExporter"
    "cloudwatch-exporter:cloudwatchExporter"
    "common-dashboards:commonDashboards"
    "goldpinger-exporter:goldpingerExporter"
    "grafana-operator:grafana"
    "graphite-remote-adapter:graphite_remote_adapter"
    "json-exporter:jsonExporter"
    "network-latency-exporter:networkLatencyExporter"
    "prometheus-operator:prometheus"
    "prometheus-adapter-operator:prometheusAdapter"
    "promitor-agent-resource-discovery:promitorAgentResourceDiscovery"
    "promitor-agent-scraper:promitorAgentScraper"
    "ssl-exporter:sslExporter"
    "stackdriver-exporter:stackdriverExporter"
    "version-exporter:versionExporter"
    "victoriametrics-operator:victoriametrics"
)

cd "${repository_root}"

for rename in "${renames[@]}"; do
    old_name="${rename%%:*}"
    new_name="${rename#*:}"
    old_directory="${chart_root}/${old_name}"
    new_directory="${chart_root}/${new_name}"

    if [[ -e "${old_directory}" ]]; then
        echo "The obsolete subchart directory still exists: ${old_directory}" >&2
        failures=1
    fi

    if [[ ! -d "${new_directory}" ]]; then
        echo "The renamed subchart directory is missing: ${new_directory}" >&2
        failures=1
    fi

    stale_references="$(
        git grep -n -F "charts/${old_name}" -- . \
            ':(exclude)tools/verify-subchart-paths.sh' || true
    )"
    if [[ -n "${stale_references}" ]]; then
        echo "Replace stale subchart path references from ${old_name} to ${new_name}:" >&2
        printf '%s\n' "${stale_references}" >&2
        failures=1
    fi
done

if [[ "${failures}" -ne 0 ]]; then
    exit 1
fi

echo "Subchart directory names and repository references are consistent."
