#!/usr/bin/env bash

set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
chart_dir="$(cd "${script_dir}/../../charts/qubership-monitoring-operator" && pwd)"

manifest="$(
    helm template monitoring-operator "${chart_dir}" \
        --show-only charts/grafana/templates/configmap-extra-vars.yaml
)"

assert_contains() {
    local expected="$1"

    if ! grep -Fq -- "${expected}" <<<"${manifest}"; then
        echo "Expected the Grafana extra-vars ConfigMap to contain: ${expected}" >&2
        exit 1
    fi
}

assert_contains '"GF_ANALYTICS_CHECK_FOR_PLUGIN_UPDATES": "false"'
assert_contains '"GF_PLUGINS_PLUGIN_ADMIN_ENABLED": "false"'
assert_contains '"GF_PLUGINS_PREINSTALL_AUTO_UPDATE": "false"'
assert_contains '"GF_PLUGINS_PREINSTALL_DISABLED": "true"'

echo "Grafana plugin update checks passed"
