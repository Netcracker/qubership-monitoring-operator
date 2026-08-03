#!/usr/bin/env bash

set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
chart_dir="$(cd "${script_dir}/../../charts/qubership-monitoring-operator" && pwd)"

render_prometheus_spec() {
    local platform_api="$1"
    shift

    helm template monitoring-operator "${chart_dir}" \
        --api-versions "${platform_api}" \
        --set prometheus.install=true \
        --show-only templates/operator/platformmonitoring.yaml \
        "$@" \
        | awk '
            /^  prometheus:$/ { in_prometheus = 1 }
            in_prometheus && /^    operator:$/ { exit }
            in_prometheus { print }
        '
}

assert_contains() {
    local manifest="$1"
    local expected="$2"

    if ! grep -Fq -- "${expected}" <<<"${manifest}"; then
        echo "Expected the Prometheus security context to contain: ${expected}" >&2
        exit 1
    fi
}

assert_not_contains() {
    local manifest="$1"
    local unexpected="$2"

    if grep -Fq -- "${unexpected}" <<<"${manifest}"; then
        echo "Expected the Prometheus security context not to contain: ${unexpected}" >&2
        exit 1
    fi
}

kubernetes_manifest="$(render_prometheus_spec "")"

assert_contains "${kubernetes_manifest}" "runAsNonRoot: true"
assert_contains "${kubernetes_manifest}" "runAsUser: 2000"
assert_contains "${kubernetes_manifest}" "runAsGroup: 2000"
assert_contains "${kubernetes_manifest}" "fsGroup: 2000"
assert_contains "${kubernetes_manifest}" "type: RuntimeDefault"

openshift_manifest="$(
    render_prometheus_spec \
        "security.openshift.io/v1/SecurityContextConstraints"
)"

assert_contains "${openshift_manifest}" "runAsNonRoot: true"
assert_contains "${openshift_manifest}" "type: RuntimeDefault"
assert_not_contains "${openshift_manifest}" "runAsUser:"
assert_not_contains "${openshift_manifest}" "runAsGroup:"
assert_not_contains "${openshift_manifest}" "fsGroup:"

configured_manifest="$(
    render_prometheus_spec "" \
        --set prometheus.securityContext.runAsUser=3000 \
        --set prometheus.securityContext.runAsGroup=3001 \
        --set prometheus.securityContext.runAsNonRoot=false \
        --set prometheus.securityContext.seccompProfile.type=Unconfined
)"

assert_contains "${configured_manifest}" "runAsUser: 3000"
assert_contains "${configured_manifest}" "runAsGroup: 3001"
assert_contains "${configured_manifest}" "runAsNonRoot: true"
assert_contains "${configured_manifest}" "type: RuntimeDefault"

echo "Prometheus security hardening checks passed"
