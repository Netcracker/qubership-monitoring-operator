#!/usr/bin/env bash

set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
chart_dir="$(cd "${script_dir}/../../charts/qubership-monitoring-operator" && pwd)"
render_dir="$(mktemp -d)"
trap 'rm -rf "${render_dir}"' EXIT

exporter_args=(
    --set blackboxExporter.install=true
    --set certExporter.install=true
    --set cloudEventsExporter.install=true
    --set cloudwatchExporter.install=true
    --set goldpingerExporter.install=true
    --set jsonExporter.install=true
    --set networkLatencyExporter.install=true
    --set promitorAgentResourceDiscovery.install=true
    --set promitorAgentScraper.install=true
    --set sslExporter.install=true
    --set stackdriverExporter.install=true
    --set versionExporter.install=true
    --set grafana.operator.install=false
)

render_exporters() {
    local output_dir="$1"
    shift

    helm template hardening "${chart_dir}" \
        --namespace monitoring \
        --output-dir "${output_dir}" \
        "${exporter_args[@]}" \
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

kubernetes_dir="${render_dir}/kubernetes"
openshift_dir="${render_dir}/openshift"
blackbox_daemonset_dir="${render_dir}/blackbox-daemonset"
enforcement_dir="${render_dir}/enforcement"

render_exporters "${kubernetes_dir}"
render_exporters \
    "${openshift_dir}" \
    --api-versions security.openshift.io/v1/SecurityContextConstraints

helm template hardening "${chart_dir}" \
    --namespace monitoring \
    --output-dir "${blackbox_daemonset_dir}" \
    --set blackboxExporter.install=true \
    --set blackboxExporter.asDaemonSet=true \
    --set grafana.operator.install=false \
    >/dev/null

helm template hardening "${chart_dir}" \
    --namespace monitoring \
    --api-versions security.openshift.io/v1/SecurityContextConstraints \
    --output-dir "${enforcement_dir}" \
    --set blackboxExporter.install=true \
    --set blackboxExporter.securityContext.runAsUser=3000 \
    --set blackboxExporter.securityContext.seccompProfile.type=Unconfined \
    --set blackboxExporter.containerSecurityContext.allowPrivilegeEscalation=true \
    --set blackboxExporter.containerSecurityContext.readOnlyRootFilesystem=false \
    --set 'blackboxExporter.containerSecurityContext.capabilities.drop={NET_RAW}' \
    --set grafana.operator.install=false \
    >/dev/null

relative_manifests=(
    blackboxExporter/templates/deployment.yaml
    certExporter/templates/deployment.yaml
    certExporter/templates/daemonset.yaml
    cloudEventsExporter/templates/deployment.yaml
    cloudwatchExporter/templates/deployment.yaml
    goldpingerExporter/templates/daemonset.yaml
    jsonExporter/templates/deployment.yaml
    promitorAgentResourceDiscovery/templates/deployment.yaml
    promitorAgentScraper/templates/deployment.yaml
    sslExporter/templates/daemonset.yaml
    stackdriverExporter/templates/deployment.yaml
    versionExporter/templates/deployment.yaml
)

for relative_manifest in "${relative_manifests[@]}"; do
    kubernetes_manifest="${kubernetes_dir}/qubership-monitoring-operator/charts/${relative_manifest}"
    openshift_manifest="${openshift_dir}/qubership-monitoring-operator/charts/${relative_manifest}"

    assert_contains "${kubernetes_manifest}" "runAsNonRoot: true"
    assert_contains "${kubernetes_manifest}" "type: RuntimeDefault"
    assert_contains "${kubernetes_manifest}" "allowPrivilegeEscalation: false"
    assert_contains "${kubernetes_manifest}" "readOnlyRootFilesystem: true"
    assert_contains "${kubernetes_manifest}" "- ALL"
    assert_contains "${kubernetes_manifest}" "sizeLimit: 100Mi"

    assert_contains "${openshift_manifest}" "runAsNonRoot: true"
    assert_contains "${openshift_manifest}" "type: RuntimeDefault"
    assert_not_contains "${openshift_manifest}" "runAsUser:"
    assert_not_contains "${openshift_manifest}" "runAsGroup:"
    assert_not_contains "${openshift_manifest}" "fsGroup:"
done

blackbox_daemonset_manifest="${blackbox_daemonset_dir}/qubership-monitoring-operator/charts/blackboxExporter/templates/daemonset.yaml"
assert_contains "${blackbox_daemonset_manifest}" "runAsNonRoot: true"
assert_contains "${blackbox_daemonset_manifest}" "allowPrivilegeEscalation: false"
assert_contains "${blackbox_daemonset_manifest}" "readOnlyRootFilesystem: true"
assert_contains "${blackbox_daemonset_manifest}" "sizeLimit: 100Mi"

enforcement_manifest="${enforcement_dir}/qubership-monitoring-operator/charts/blackboxExporter/templates/deployment.yaml"
assert_contains "${enforcement_manifest}" "runAsNonRoot: true"
assert_contains "${enforcement_manifest}" "type: RuntimeDefault"
assert_contains "${enforcement_manifest}" "allowPrivilegeEscalation: false"
assert_contains "${enforcement_manifest}" "readOnlyRootFilesystem: true"
assert_contains "${enforcement_manifest}" "- ALL"
assert_contains "${enforcement_manifest}" "runAsUser: 3000"
assert_not_contains "${enforcement_manifest}" "Unconfined"
assert_not_contains "${enforcement_manifest}" "- NET_RAW"

assert_render_fails "securityContext.runAsUser=0 conflicts" \
    --set blackboxExporter.install=true \
    --set blackboxExporter.securityContext.runAsUser=0
assert_render_fails "securityContext.runAsNonRoot=false conflicts" \
    --set blackboxExporter.install=true \
    --set blackboxExporter.securityContext.runAsNonRoot=false
assert_render_fails "containerSecurityContext.privileged=true conflicts" \
    --set blackboxExporter.install=true \
    --set blackboxExporter.containerSecurityContext.privileged=true
assert_render_fails "containerSecurityContext.capabilities.add conflicts" \
    --set blackboxExporter.install=true \
    --set 'blackboxExporter.containerSecurityContext.capabilities.add={NET_RAW}'
assert_render_fails "securityContext.runAsUser=0 conflicts" \
    --set sslExporter.install=true \
    --set sslExporter.podSecurityContext.runAsUser=0
assert_render_fails "securityContext.runAsNonRoot=false conflicts" \
    --set sslExporter.install=true \
    --set sslExporter.podSecurityContext.runAsNonRoot=false

# User-defined volume extensions must not collide with the chart-managed temporary volume.
assert_render_fails 'volume name "monitoring-tmp" is reserved' \
    --set jsonExporter.install=true \
    --set 'jsonExporter.additionalVolumes[0].name=monitoring-tmp' \
    --set 'jsonExporter.additionalVolumes[0].emptyDir.medium=Memory'
assert_render_fails 'volume name "monitoring-tmp" is reserved' \
    --set promitorAgentResourceDiscovery.install=true \
    --set 'promitorAgentResourceDiscovery.extraVolumes[0].name=monitoring-tmp' \
    --set 'promitorAgentResourceDiscovery.extraVolumes[0].emptyDir.medium=Memory'
assert_render_fails 'volume name "monitoring-tmp" is reserved' \
    --set sslExporter.install=true \
    --set 'sslExporter.additionalHostPathVolumes[0].volumeName=monitoring-tmp' \
    --set 'sslExporter.additionalHostPathVolumes[0].volumePath=/etc/ssl-extra'
assert_render_fails 'volume name "monitoring-tmp" is reserved' \
    --set certExporter.install=true \
    --set certExporter.certsInFiles.enabled=true \
    --set 'certExporter.additionalHostPathVolumes[0].volumeName=monitoring-tmp' \
    --set 'certExporter.additionalHostPathVolumes[0].volumePath=/etc/cert-extra'

user_tmp_dir="${render_dir}/user-tmp"
helm template hardening "${chart_dir}" \
    --namespace monitoring \
    --output-dir "${user_tmp_dir}" \
    --set jsonExporter.install=true \
    --set 'jsonExporter.additionalVolumes[0].name=scratch' \
    --set 'jsonExporter.additionalVolumes[0].emptyDir.medium=Memory' \
    --set 'jsonExporter.additionalVolumeMounts[0].name=scratch' \
    --set 'jsonExporter.additionalVolumeMounts[0].mountPath=/tmp' \
    --set promitorAgentResourceDiscovery.install=true \
    --set 'promitorAgentResourceDiscovery.extraVolumes[0].name=scratch' \
    --set 'promitorAgentResourceDiscovery.extraVolumes[0].emptyDir.medium=Memory' \
    --set 'promitorAgentResourceDiscovery.extraVolumeMounts[0].name=scratch' \
    --set 'promitorAgentResourceDiscovery.extraVolumeMounts[0].mountPath=/tmp/' \
    --set grafana.operator.install=false \
    >/dev/null
json_user_tmp_manifest="${user_tmp_dir}/qubership-monitoring-operator/charts/jsonExporter/templates/deployment.yaml"
promitor_user_tmp_manifest="${user_tmp_dir}/qubership-monitoring-operator/charts/promitorAgentResourceDiscovery/templates/deployment.yaml"
for manifest in "${json_user_tmp_manifest}" "${promitor_user_tmp_manifest}"; do
    # The user mount covers /tmp, so the chart must declare its volume without a second mount at /tmp.
    assert_contains "${manifest}" "name: monitoring-tmp"
    if [ "$(grep -c -- "- name: monitoring-tmp" "${manifest}")" -ne 1 ]; then
        echo "Expected ${manifest} to reference monitoring-tmp only as a volume" >&2
        exit 1
    fi
done

network_latency_manifest="${kubernetes_dir}/qubership-monitoring-operator/charts/networkLatencyExporter/templates/daemonset.yaml"
assert_contains "${network_latency_manifest}" "runAsUser: 2001"
assert_contains "${network_latency_manifest}" "runAsNonRoot: true"
assert_contains "${network_latency_manifest}" "type: RuntimeDefault"
assert_contains "${network_latency_manifest}" "allowPrivilegeEscalation: false"
assert_contains "${network_latency_manifest}" "readOnlyRootFilesystem: true"
assert_contains "${network_latency_manifest}" "- ALL"
assert_not_contains "${network_latency_manifest}" "- NET_RAW"
assert_contains "${network_latency_manifest}" "sizeLimit: 100Mi"
assert_not_contains "${network_latency_manifest}" "shareProcessNamespace:"

network_latency_openshift_manifest="${openshift_dir}/qubership-monitoring-operator/charts/networkLatencyExporter/templates/daemonset.yaml"
assert_not_contains "${network_latency_openshift_manifest}" "runAsUser:"
assert_not_contains "${network_latency_openshift_manifest}" "fsGroup:"
assert_contains "${network_latency_openshift_manifest}" "runAsNonRoot: true"
assert_not_contains "${network_latency_openshift_manifest}" "- NET_RAW"

cert_daemonset_manifest="${kubernetes_dir}/qubership-monitoring-operator/charts/certExporter/templates/daemonset.yaml"
cert_scc_manifest="${openshift_dir}/qubership-monitoring-operator/charts/certExporter/templates/securitycontextconstraints.yaml"
ssl_daemonset_manifest="${kubernetes_dir}/qubership-monitoring-operator/charts/sslExporter/templates/daemonset.yaml"
assert_contains "${cert_daemonset_manifest}" "hostPath:"
cert_scc_run_as_user_type="$(awk '/^runAsUser:$/ { getline; print $2; exit }' "${cert_scc_manifest}")"
if [[ "${cert_scc_run_as_user_type}" != "RunAsAny" ]]; then
    echo "Expected ${cert_scc_manifest} runAsUser.type to equal RunAsAny, got: ${cert_scc_run_as_user_type}" >&2
    exit 1
fi
assert_contains "${ssl_daemonset_manifest}" "hostPath:"
assert_not_contains "${cert_daemonset_manifest}" "hostNetwork: true"
assert_not_contains "${ssl_daemonset_manifest}" "hostNetwork: true"

echo "Exporter security hardening checks passed"
