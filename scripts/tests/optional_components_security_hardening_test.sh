#!/usr/bin/env bash

set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
chart_dir="$(cd "${script_dir}/../../charts/qubership-monitoring-operator" && pwd)"
render_dir="$(mktemp -d)"
trap 'rm -rf "${render_dir}"' EXIT

component_args=(
    --set global.privilegedRights=true
    --set grafana.install=true
    --set grafana.imageRenderer.install=true
    --set graphite_remote_adapter.install=true
    --set prometheusAdapter.install=true
    --set promxy.install=true
    --set victoriametrics.vmOperator.install=true
)

render_components() {
    local output_dir="$1"
    shift

    helm template hardening "${chart_dir}" \
        --namespace monitoring \
        --output-dir "${output_dir}" \
        "${component_args[@]}" \
        "$@" \
        >/dev/null
}

assert_contains() {
    local manifest="$1"
    local expected="$2"

    if ! grep -Fq -- "${expected}" "${manifest}"; then
        echo "Expected ${manifest} to contain: ${expected}" >&2
        exit 1
    fi
}

assert_not_contains() {
    local manifest="$1"
    local unexpected="$2"

    if grep -Fq -- "${unexpected}" "${manifest}"; then
        echo "Expected ${manifest} not to contain: ${unexpected}" >&2
        exit 1
    fi
}

kubernetes_dir="${render_dir}/kubernetes"
openshift_dir="${render_dir}/openshift"
enforcement_dir="${render_dir}/enforcement"

render_components "${kubernetes_dir}"
render_components \
    "${openshift_dir}" \
    --api-versions security.openshift.io/v1/SecurityContextConstraints
render_components \
    "${enforcement_dir}" \
    --api-versions security.openshift.io/v1/SecurityContextConstraints \
    --set grafana.imageRenderer.securityContext.runAsUser=0 \
    --set grafana.imageRenderer.securityContext.runAsNonRoot=false \
    --set grafana.imageRenderer.securityContext.seccompProfile.type=Unconfined \
    --set victoriametrics.cleanup.hook.containerSecurityContext.allowPrivilegeEscalation=true \
    --set victoriametrics.cleanup.hook.containerSecurityContext.readOnlyRootFilesystem=false \
    --set 'victoriametrics.cleanup.hook.containerSecurityContext.capabilities.drop={NET_RAW}'

relative_manifests=(
    grafana/templates/image-renderer/deployment.yaml
    graphite_remote_adapter/templates/deployment.yaml
    prometheusAdapter/templates/deployment.yaml
    promxy/templates/deployment.yaml
)

for relative_manifest in "${relative_manifests[@]}"; do
    kubernetes_manifest="${kubernetes_dir}/qubership-monitoring-operator/charts/${relative_manifest}"
    openshift_manifest="${openshift_dir}/qubership-monitoring-operator/charts/${relative_manifest}"

    assert_contains "${kubernetes_manifest}" "runAsNonRoot: true"
    assert_contains "${kubernetes_manifest}" "runAsUser: 2000"
    assert_contains "${kubernetes_manifest}" "runAsGroup: 2000"
    assert_contains "${kubernetes_manifest}" "fsGroup: 2000"
    assert_contains "${kubernetes_manifest}" "type: RuntimeDefault"
    assert_contains "${kubernetes_manifest}" "allowPrivilegeEscalation: false"
    assert_contains "${kubernetes_manifest}" "readOnlyRootFilesystem: true"
    assert_contains "${kubernetes_manifest}" "- ALL"
    assert_contains "${kubernetes_manifest}" "sizeLimit: 100Mi"
    assert_not_contains "${kubernetes_manifest}" "hostNetwork: true"
    assert_not_contains "${kubernetes_manifest}" "hostPID: true"
    assert_not_contains "${kubernetes_manifest}" "hostIPC: true"
    assert_not_contains "${kubernetes_manifest}" "hostPath:"

    assert_contains "${openshift_manifest}" "runAsNonRoot: true"
    assert_contains "${openshift_manifest}" "type: RuntimeDefault"
    assert_not_contains "${openshift_manifest}" "runAsUser:"
    assert_not_contains "${openshift_manifest}" "runAsGroup:"
    assert_not_contains "${openshift_manifest}" "fsGroup:"
done

enforced_renderer_manifest="${enforcement_dir}/qubership-monitoring-operator/charts/grafana/templates/image-renderer/deployment.yaml"
assert_contains "${enforced_renderer_manifest}" "docker.io/grafana/grafana-image-renderer:v5.8.3"
assert_contains "${enforced_renderer_manifest}" "BROWSER_FLAGS"
assert_contains "${enforced_renderer_manifest}" "runAsNonRoot: true"
assert_contains "${enforced_renderer_manifest}" "type: RuntimeDefault"
assert_not_contains "${enforced_renderer_manifest}" "runAsUser:"
assert_not_contains "${enforced_renderer_manifest}" "Unconfined"

renderer_config="${kubernetes_dir}/qubership-monitoring-operator/charts/grafana/templates/configmap-extra-vars.yaml"
assert_contains "${renderer_config}" "GF_RENDERING_SERVER_URL"
assert_contains "${renderer_config}" "http://grafana-image-renderer:8081/render"

cleanup_kubernetes_manifest="${kubernetes_dir}/qubership-monitoring-operator/charts/victoriametrics/templates/pre-delete/job.yaml"
cleanup_openshift_manifest="${openshift_dir}/qubership-monitoring-operator/charts/victoriametrics/templates/pre-delete/job.yaml"
cleanup_enforced_manifest="${enforcement_dir}/qubership-monitoring-operator/charts/victoriametrics/templates/pre-delete/job.yaml"

for cleanup_manifest in "${cleanup_kubernetes_manifest}" "${cleanup_openshift_manifest}"; do
    assert_contains "${cleanup_manifest}" "type: RuntimeDefault"
    assert_contains "${cleanup_manifest}" "allowPrivilegeEscalation: false"
    assert_contains "${cleanup_manifest}" "readOnlyRootFilesystem: true"
    assert_contains "${cleanup_manifest}" "- ALL"
    assert_contains "${cleanup_manifest}" "sizeLimit: 100Mi"
done

assert_contains "${cleanup_kubernetes_manifest}" "runAsNonRoot: true"
assert_contains "${cleanup_kubernetes_manifest}" "runAsUser: 2000"
assert_contains "${cleanup_kubernetes_manifest}" "runAsGroup: 2000"
assert_contains "${cleanup_kubernetes_manifest}" "fsGroup: 2000"
assert_not_contains "${cleanup_kubernetes_manifest}" "openshift.io/required-scc"
assert_not_contains "${cleanup_openshift_manifest}" "runAsNonRoot:"
assert_not_contains "${cleanup_openshift_manifest}" "runAsUser:"
assert_not_contains "${cleanup_openshift_manifest}" "runAsGroup:"
assert_not_contains "${cleanup_openshift_manifest}" "fsGroup:"
assert_contains "${cleanup_openshift_manifest}" "openshift.io/required-scc: restricted-v2"
assert_contains "${cleanup_enforced_manifest}" "allowPrivilegeEscalation: false"
assert_contains "${cleanup_enforced_manifest}" "readOnlyRootFilesystem: true"
assert_contains "${cleanup_enforced_manifest}" "- ALL"
assert_not_contains "${cleanup_enforced_manifest}" "- NET_RAW"

promxy_manifest="${kubernetes_dir}/qubership-monitoring-operator/charts/promxy/templates/deployment.yaml"
if [[ "$(grep -Fc -- "allowPrivilegeEscalation: false" "${promxy_manifest}")" -ne 2 ]]; then
    echo "Expected both Promxy containers to disable privilege escalation" >&2
    exit 1
fi

adapter_cr="${kubernetes_dir}/qubership-monitoring-operator/charts/prometheusAdapter/templates/prometheusadapter.yaml"
assert_contains "${adapter_cr}" "runAsUser: 2000"
assert_contains "${adapter_cr}" "fsGroup: 2000"
assert_not_contains "${adapter_cr}" "runAsNonRoot:"
assert_not_contains "${adapter_cr}" "seccompProfile:"
assert_not_contains "${adapter_cr}" "allowPrivilegeEscalation:"

echo "Optional component security hardening checks passed"
