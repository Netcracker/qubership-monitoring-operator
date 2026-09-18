#!/usr/bin/env bash

set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
chart_dir="$(cd "${script_dir}/../../charts/qubership-monitoring-operator" && pwd)"

render_alertmanager_spec() {
    local platform_api="$1"
    shift

    helm template monitoring-operator "${chart_dir}" \
        --api-versions "${platform_api}" \
        --set alertManager.install=true \
        --show-only templates/operator/platformmonitoring.yaml \
        "$@" | awk '
            /^  alertManager:$/ { in_alertmanager = 1 }
            in_alertmanager && /^  [a-zA-Z]/ && !/^  alertManager:$/ { exit }
            in_alertmanager { print }
        '
}

assert_contains() {
    local manifest="$1"
    local expected="$2"

    if ! grep -Fq -- "${expected}" <<<"${manifest}"; then
        echo "Expected the Alertmanager security context to contain: ${expected}" >&2
        exit 1
    fi
}

assert_not_contains() {
    local manifest="$1"
    local unexpected="$2"

    if grep -Fq -- "${unexpected}" <<<"${manifest}"; then
        echo "Expected the Alertmanager security context not to contain: ${unexpected}" >&2
        exit 1
    fi
}

assert_render_fails() {
    local expected_error="$1"
    shift
    local output

    if output="$(helm template monitoring-operator "${chart_dir}" "$@" 2>&1)"; then
        echo "Expected rendering to fail with: ${expected_error}" >&2
        exit 1
    fi
    if ! grep -Fq -- "${expected_error}" <<<"${output}"; then
        echo "Expected the rendering error to contain: ${expected_error}" >&2
        echo "${output}" >&2
        exit 1
    fi
}

kubernetes_manifest="$(render_alertmanager_spec "")"

assert_contains "${kubernetes_manifest}" "runAsUser: 2000"
assert_contains "${kubernetes_manifest}" "runAsGroup: 2000"
assert_contains "${kubernetes_manifest}" "fsGroup: 2000"
assert_not_contains "${kubernetes_manifest}" "runAsNonRoot:"
assert_not_contains "${kubernetes_manifest}" "seccompProfile:"

openshift_manifest="$(
    render_alertmanager_spec \
        "security.openshift.io/v1/SecurityContextConstraints"
)"

assert_not_contains "${openshift_manifest}" "runAsUser:"
assert_not_contains "${openshift_manifest}" "runAsGroup:"
assert_not_contains "${openshift_manifest}" "fsGroup:"
assert_not_contains "${openshift_manifest}" "runAsNonRoot:"
assert_not_contains "${openshift_manifest}" "seccompProfile:"

configured_manifest="$(
    render_alertmanager_spec "" \
        --set alertManager.securityContext.runAsUser=3000 \
        --set alertManager.securityContext.runAsGroup=3001 \
        --set alertManager.securityContext.seccompProfile.type=Unconfined
)"

assert_contains "${configured_manifest}" "runAsUser: 3000"
assert_contains "${configured_manifest}" "runAsGroup: 3001"
assert_not_contains "${configured_manifest}" "runAsNonRoot:"
assert_not_contains "${configured_manifest}" "seccompProfile:"

assert_render_fails "securityContext.runAsUser=0 conflicts" \
    --set alertManager.install=true \
    --set alertManager.securityContext.runAsUser=0
assert_render_fails "securityContext.runAsNonRoot=false conflicts" \
    --set alertManager.install=true \
    --set alertManager.securityContext.runAsNonRoot=false

configured_openshift_manifest="$(
    render_alertmanager_spec "security.openshift.io/v1/SecurityContextConstraints" \
        --set alertManager.securityContext.runAsUser=3000 \
        --set alertManager.securityContext.runAsGroup=3001 \
        --set alertManager.securityContext.fsGroup=3002
)"

assert_contains "${configured_openshift_manifest}" "runAsUser: 3000"
assert_contains "${configured_openshift_manifest}" "runAsGroup: 3001"
assert_contains "${configured_openshift_manifest}" "fsGroup: 3002"

echo "Alertmanager security hardening checks passed"
