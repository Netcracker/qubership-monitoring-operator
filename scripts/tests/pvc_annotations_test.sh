#!/usr/bin/env bash

set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
chart_dir="$(cd "${script_dir}/../../charts/qubership-monitoring-operator" && pwd)"

render_platformmonitoring() {
    helm template monitoring-operator "${chart_dir}" \
        --set prometheus.install=true \
        --set victoriametrics.vmSingle.install=true \
        --set victoriametrics.vmAlertManager.install=true \
        --set grafana.install=true \
        --show-only templates/operator/platformmonitoring.yaml \
        "$@"
}

assert_count() {
    local manifest="$1"
    local expected_text="$2"
    local expected_count="$3"
    local actual_count

    actual_count="$(grep -Fc -- "${expected_text}" <<<"${manifest}" || true)"
    if [[ "${actual_count}" != "${expected_count}" ]]; then
        echo "Expected ${expected_count} occurrence(s) of '${expected_text}', found ${actual_count}" >&2
        exit 1
    fi
}

if [[ "$(render_platformmonitoring)" != "$(render_platformmonitoring --set pvc=null)" ]]; then
    echo "Expected pvc=null to render the same PlatformMonitoring as the default values" >&2
    exit 1
fi

storage_args=(
    --set prometheus.storage.volumeClaimTemplate.spec.resources.requests.storage=1Gi
    --set victoriametrics.vmSingle.storage.resources.requests.storage=1Gi
    --set victoriametrics.vmAlertManager.storage.volumeClaimTemplate.spec.resources.requests.storage=1Gi
    --set grafana.dataStorage.size=1Gi
)
global_args=(
    --set-string 'pvc.metadata.annotations.example\.com/keep=global'
    --set-string 'pvc.metadata.annotations.example\.com/global-only=yes'
)

annotated_manifest="$(
    render_platformmonitoring "${storage_args[@]}" "${global_args[@]}" \
        --set-string 'prometheus.storage.volumeClaimTemplate.metadata.annotations.example\.com/keep=component' \
        --set-string 'grafana.dataStorage.annotations.example\.com/keep=component'
)"

assert_count "${annotated_manifest}" "example.com/global-only: \"yes\"" 4
assert_count "${annotated_manifest}" "example.com/keep: global" 2
assert_count "${annotated_manifest}" "example.com/keep: component" 2

storage_manifest="$(render_platformmonitoring "${storage_args[@]}")"
if [[ "${storage_manifest}" != "$(render_platformmonitoring "${storage_args[@]}" --set pvc=null)" ]]; then
    echo "Expected pvc=null to render the same storage settings as the default values" >&2
    exit 1
fi
assert_count "${storage_manifest}" "annotations: {}" 0
assert_count "$(render_platformmonitoring "${global_args[@]}" --set prometheus.storage.emptyDir.medium=Memory)" "example.com/" 0

no_storage_manifest="$(render_platformmonitoring "${global_args[@]}")"

assert_count "${no_storage_manifest}" "volumeClaimTemplate:" 0
assert_count "${no_storage_manifest}" "storageMetadata:" 0
assert_count "${no_storage_manifest}" "dataStorage:" 0
assert_count "${no_storage_manifest}" "example.com/" 0

echo "PVC annotation checks passed"
