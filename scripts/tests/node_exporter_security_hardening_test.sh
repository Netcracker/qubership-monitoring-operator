#!/usr/bin/env bash

set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
chart_dir="$(cd "${script_dir}/../../charts/qubership-monitoring-operator" && pwd)"

render_node_exporter_spec() {
    local platform_api="$1"
    shift

    helm template monitoring-operator "${chart_dir}" \
        --api-versions "${platform_api}" \
        --set nodeExporter.install=true \
        --show-only templates/operator/platformmonitoring.yaml \
        "$@" | awk '
            /^  nodeExporter:$/ { in_node_exporter = 1 }
            in_node_exporter && /^  [a-zA-Z]/ && !/^  nodeExporter:$/ { exit }
            in_node_exporter { print }
        '
}

assert_contains() {
    local manifest="$1"
    local expected="$2"

    if ! grep -Fq -- "${expected}" <<<"${manifest}"; then
        echo "Expected the node-exporter security context to contain: ${expected}" >&2
        exit 1
    fi
}

assert_not_contains() {
    local manifest="$1"
    local unexpected="$2"

    if grep -Fq -- "${unexpected}" <<<"${manifest}"; then
        echo "Expected the node-exporter security context not to contain: ${unexpected}" >&2
        exit 1
    fi
}

kubernetes_manifest="$(render_node_exporter_spec "")"

assert_contains "${kubernetes_manifest}" "runAsUser: 2000"
assert_contains "${kubernetes_manifest}" "runAsGroup: 2000"
assert_contains "${kubernetes_manifest}" "fsGroup: 2000"

openshift_manifest="$(
    render_node_exporter_spec \
        "security.openshift.io/v1/SecurityContextConstraints"
)"

assert_not_contains "${openshift_manifest}" "runAsUser:"
assert_not_contains "${openshift_manifest}" "runAsGroup:"
assert_not_contains "${openshift_manifest}" "fsGroup:"

configured_manifest="$(
    render_node_exporter_spec "" \
        --set nodeExporter.securityContext.runAsUser=3000 \
        --set nodeExporter.securityContext.runAsGroup=3001 \
        --set nodeExporter.securityContext.fsGroup=3002
)"

assert_contains "${configured_manifest}" "runAsUser: 3000"
assert_contains "${configured_manifest}" "runAsGroup: 3001"
assert_contains "${configured_manifest}" "fsGroup: 3002"

echo "node-exporter security hardening checks passed"
