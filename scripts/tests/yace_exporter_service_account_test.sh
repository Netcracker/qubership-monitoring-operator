#!/usr/bin/env bash

set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
chart_dir="$(cd "${script_dir}/../../charts/qubership-monitoring-operator" && pwd)"
render_dir="$(mktemp -d)"
trap 'rm -rf "${render_dir}"' EXIT

templates_dir="qubership-monitoring-operator/charts/yaceExporter/templates"

render() {
    local output_dir="$1"
    shift

    helm template monitoring "${chart_dir}" \
        --namespace monitoring \
        --output-dir "${output_dir}" \
        --set yaceExporter.install=true \
        --set grafana.operator.install=false \
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

assert_missing() {
    local manifest="$1"

    if [[ -e "${manifest}" ]]; then
        echo "Expected ${manifest} not to be rendered" >&2
        exit 1
    fi
}

# Defaults: the chart creates the ServiceAccount, and neither it nor the pod mounts the API token.
default_dir="${render_dir}/default"
render "${default_dir}"
deployment="${default_dir}/${templates_dir}/deployment.yaml"
service_account="${default_dir}/${templates_dir}/serviceaccount.yaml"
assert_contains "${deployment}" "serviceAccountName: yace-exporter"
assert_contains "${deployment}" "automountServiceAccountToken: false"
assert_contains "${service_account}" "name: yace-exporter"
assert_contains "${service_account}" "automountServiceAccountToken: false"

# Reuse: an existing ServiceAccount, for example one bound to an IRSA role, must still be used by the pod.
reuse_dir="${render_dir}/reuse"
render "${reuse_dir}" \
    --set yaceExporter.serviceAccount.install=false \
    --set yaceExporter.serviceAccount.name=existing-irsa \
    --set yaceExporter.rbac.createClusterRole=true \
    --set yaceExporter.rbac.createClusterRoleBinding=true
assert_contains "${reuse_dir}/${templates_dir}/deployment.yaml" "serviceAccountName: existing-irsa"
assert_contains "${reuse_dir}/${templates_dir}/clusterrolebinding.yaml" "name: existing-irsa"
assert_missing "${reuse_dir}/${templates_dir}/serviceaccount.yaml"

# An empty name falls back to the exporter name.
empty_name_dir="${render_dir}/empty-name"
render "${empty_name_dir}" \
    --set yaceExporter.serviceAccount.install=false \
    --set yaceExporter.serviceAccount.name=
assert_contains "${empty_name_dir}/${templates_dir}/deployment.yaml" "serviceAccountName: yace-exporter"

# Opting back in to the API token applies to the pod as well as to the created ServiceAccount.
automount_dir="${render_dir}/automount"
render "${automount_dir}" \
    --set yaceExporter.serviceAccount.automountServiceAccountToken=true
assert_contains "${automount_dir}/${templates_dir}/deployment.yaml" "automountServiceAccountToken: true"
assert_contains "${automount_dir}/${templates_dir}/serviceaccount.yaml" "automountServiceAccountToken: true"

echo "YACE exporter ServiceAccount checks passed"
