#!/usr/bin/env bash

set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
chart_dir="$(cd "${script_dir}/../../charts/qubership-monitoring-operator" && pwd)"

assert_contains() {
    local manifest="$1"
    local expected="$2"

    if ! grep -Fq -- "${expected}" <<<"${manifest}"; then
        echo "Expected the monitoring-operator manifest to contain: ${expected}" >&2
        exit 1
    fi
}

assert_not_contains() {
    local manifest="$1"
    local unexpected="$2"

    if grep -Fq -- "${unexpected}" <<<"${manifest}"; then
        echo "Expected the monitoring-operator manifest not to contain: ${unexpected}" >&2
        exit 1
    fi
}

kubernetes_manifest="$(
    helm template monitoring-operator "${chart_dir}" \
        --show-only templates/operator/deployment.yaml
)"

assert_contains "${kubernetes_manifest}" "runAsNonRoot: true"
assert_contains "${kubernetes_manifest}" "runAsUser: 2000"
assert_contains "${kubernetes_manifest}" "runAsGroup: 2000"
assert_contains "${kubernetes_manifest}" "type: RuntimeDefault"
assert_contains "${kubernetes_manifest}" "allowPrivilegeEscalation: false"
assert_contains "${kubernetes_manifest}" "readOnlyRootFilesystem: true"
assert_contains "${kubernetes_manifest}" "- ALL"
assert_contains "${kubernetes_manifest}" "mountPath: /tmp"
assert_contains "${kubernetes_manifest}" "sizeLimit: 16Mi"

openshift_manifest="$(
    helm template monitoring-operator "${chart_dir}" \
        --api-versions security.openshift.io/v1/SecurityContextConstraints \
        --show-only templates/operator/deployment.yaml
)"

assert_contains "${openshift_manifest}" "runAsNonRoot: true"
assert_contains "${openshift_manifest}" "type: RuntimeDefault"
assert_not_contains "${openshift_manifest}" "runAsUser:"
assert_not_contains "${openshift_manifest}" "runAsGroup:"

configured_manifest="$(
    helm template monitoring-operator "${chart_dir}" \
        --show-only templates/operator/deployment.yaml \
        --set monitoringOperator.securityContext.runAsUser=3000 \
        --set monitoringOperator.securityContext.runAsGroup=3000 \
        --set monitoringOperator.securityContext.runAsNonRoot=false \
        --set monitoringOperator.securityContext.seccompProfile.type=Unconfined
)"

assert_contains "${configured_manifest}" "runAsUser: 3000"
assert_contains "${configured_manifest}" "runAsGroup: 3000"
assert_contains "${configured_manifest}" "runAsNonRoot: true"
assert_contains "${configured_manifest}" "type: RuntimeDefault"

echo "monitoring-operator security hardening checks passed"
