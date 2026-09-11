#!/usr/bin/env bash

set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
chart_dir="$(cd "${script_dir}/../../charts/qubership-monitoring-operator" && pwd)"

render_kube_state_metrics_spec() {
    local platform_api="$1"
    shift

    helm template monitoring-operator "${chart_dir}" \
        --api-versions "${platform_api}" \
        --set kubeStateMetrics.install=true \
        --show-only templates/operator/platformmonitoring.yaml \
        "$@" | awk '
            /^  kubeStateMetrics:$/ { in_kube_state_metrics = 1 }
            in_kube_state_metrics && /^  [a-zA-Z]/ && !/^  kubeStateMetrics:$/ { exit }
            in_kube_state_metrics { print }
        '
}

assert_contains() {
    local manifest="$1"
    local expected="$2"

    if ! grep -Fq -- "${expected}" <<<"${manifest}"; then
        echo "Expected the kube-state-metrics security context to contain: ${expected}" >&2
        exit 1
    fi
}

assert_not_contains() {
    local manifest="$1"
    local unexpected="$2"

    if grep -Fq -- "${unexpected}" <<<"${manifest}"; then
        echo "Expected the kube-state-metrics security context not to contain: ${unexpected}" >&2
        exit 1
    fi
}

kubernetes_manifest="$(render_kube_state_metrics_spec "")"

assert_contains "${kubernetes_manifest}" "runAsUser: 2000"
assert_contains "${kubernetes_manifest}" "runAsGroup: 2000"
assert_contains "${kubernetes_manifest}" "fsGroup: 2000"

openshift_manifest="$(
    render_kube_state_metrics_spec \
        "security.openshift.io/v1/SecurityContextConstraints"
)"

assert_not_contains "${openshift_manifest}" "runAsUser:"
assert_not_contains "${openshift_manifest}" "runAsGroup:"
assert_not_contains "${openshift_manifest}" "fsGroup:"

configured_manifest="$(
    render_kube_state_metrics_spec "" \
        --set kubeStateMetrics.securityContext.runAsUser=3000 \
        --set kubeStateMetrics.securityContext.runAsGroup=3001 \
        --set kubeStateMetrics.securityContext.fsGroup=3002
)"

assert_contains "${configured_manifest}" "runAsUser: 3000"
assert_contains "${configured_manifest}" "runAsGroup: 3001"
assert_contains "${configured_manifest}" "fsGroup: 3002"

echo "kube-state-metrics security hardening checks passed"
