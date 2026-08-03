#!/usr/bin/env bash

set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
chart_dir="$(cd "${script_dir}/../../charts/qubership-monitoring-operator" && pwd)"

render_pushgateway_spec() {
    local platform_api="$1"
    shift

    helm template monitoring-operator "${chart_dir}" \
        --api-versions "${platform_api}" \
        --set pushgateway.install=true \
        --show-only templates/operator/platformmonitoring.yaml \
        "$@" \
        | awk '
            /^  pushgateway:$/ { in_pushgateway = 1 }
            in_pushgateway && /^  [a-zA-Z]/ && !/^  pushgateway:$/ { exit }
            in_pushgateway { print }
        '
}

assert_contains() {
    local manifest="$1"
    local expected="$2"

    if ! grep -Fq -- "${expected}" <<<"${manifest}"; then
        echo "Expected the pushgateway security context to contain: ${expected}" >&2
        exit 1
    fi
}

assert_not_contains() {
    local manifest="$1"
    local unexpected="$2"

    if grep -Fq -- "${unexpected}" <<<"${manifest}"; then
        echo "Expected the pushgateway security context not to contain: ${unexpected}" >&2
        exit 1
    fi
}

kubernetes_manifest="$(render_pushgateway_spec "")"

assert_contains "${kubernetes_manifest}" "runAsUser: 2000"
assert_contains "${kubernetes_manifest}" "runAsGroup: 2000"
assert_contains "${kubernetes_manifest}" "fsGroup: 2000"

openshift_manifest="$(
    render_pushgateway_spec \
        "security.openshift.io/v1/SecurityContextConstraints"
)"

assert_not_contains "${openshift_manifest}" "runAsUser:"
assert_not_contains "${openshift_manifest}" "runAsGroup:"
assert_not_contains "${openshift_manifest}" "fsGroup:"

configured_manifest="$(
    render_pushgateway_spec "" \
        --set pushgateway.securityContext.runAsUser=3000 \
        --set pushgateway.securityContext.runAsGroup=3001 \
        --set pushgateway.securityContext.fsGroup=3002
)"

assert_contains "${configured_manifest}" "runAsUser: 3000"
assert_contains "${configured_manifest}" "runAsGroup: 3001"
assert_contains "${configured_manifest}" "fsGroup: 3002"

echo "pushgateway security hardening checks passed"
