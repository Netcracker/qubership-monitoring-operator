#!/usr/bin/env bash

set -euo pipefail

chart_dir="${1:-charts/qubership-monitoring-operator}"
temporary_dir="$(mktemp -d)"
rendered_manifest="${temporary_dir}/manifest.yaml"
operator_rbac_manifest="${temporary_dir}/operator-rbac.yaml"
root_crd_dir="${chart_dir}/crds"
prometheus_crd_dir="${chart_dir}/charts/prometheus-operator/crds"
victoriametrics_crd_dir="${chart_dir}/charts/victoriametrics-operator/crds"
trap 'rm -rf "${temporary_dir}"' EXIT

shopt -s nullglob
root_crds=("${root_crd_dir}"/*.yaml "${root_crd_dir}"/*.yml)
if [[ "${#root_crds[@]}" -ne 1 ]]; then
    echo "The root chart contains ${#root_crds[@]} CRDs; expected only the PlatformMonitoring CRD." >&2
    exit 1
fi

if ! grep -q "^  group: monitoring.netcracker.com$" "${root_crds[0]}"; then
    echo "The root chart CRD does not belong to monitoring.netcracker.com." >&2
    exit 1
fi

prometheus_crds=("${prometheus_crd_dir}"/monitoring.coreos.com_*.yaml)
victoriametrics_prometheus_crds=("${victoriametrics_crd_dir}"/monitoring.coreos.com_*.yaml)
if [[ "${#prometheus_crds[@]}" -eq 0 ]]; then
    echo "The Prometheus subchart does not contain monitoring.coreos.com CRDs." >&2
    exit 1
fi
if [[ "${#prometheus_crds[@]}" -ne "${#victoriametrics_prometheus_crds[@]}" ]]; then
    echo "The Prometheus and VictoriaMetrics subcharts contain different Prometheus CRD counts." >&2
    exit 1
fi

for prometheus_crd in "${prometheus_crds[@]}"; do
    crd_name="$(basename "${prometheus_crd}")"
    if ! cmp -s "${prometheus_crd}" "${victoriametrics_crd_dir}/${crd_name}"; then
        echo "The Prometheus CRD ${crd_name} is not synchronized between subcharts." >&2
        exit 1
    fi
done

verify_servicemonitor_crd_count() {
    local prometheus_install="$1"
    local victoriametrics_install="$2"
    local expected_count="$3"

    helm template monitoring "${chart_dir}" \
        --include-crds \
        --set "prometheus.install=${prometheus_install}" \
        --set "victoriametrics.vmOperator.install=${victoriametrics_install}" \
        >"${rendered_manifest}"

    local actual_count
    actual_count="$(grep -c "^  name: servicemonitors.monitoring.coreos.com$" "${rendered_manifest}" || true)"
    if [[ "${actual_count}" -ne "${expected_count}" ]]; then
        echo "The chart renders ${actual_count} ServiceMonitor CRDs with prometheus.install=${prometheus_install}" \
            "and victoriametrics.vmOperator.install=${victoriametrics_install}; expected ${expected_count}." >&2
        exit 1
    fi
}

verify_servicemonitor_crd_count false false 0
verify_servicemonitor_crd_count true false 1
verify_servicemonitor_crd_count false true 1
verify_servicemonitor_crd_count true true 2

helm template monitoring "${chart_dir}" \
    --include-crds \
    --set prometheus.install=false \
    --set victoriametrics.vmOperator.install=true \
    >"${rendered_manifest}"

for resource_kind in PodMonitor ServiceMonitor; do
    if ! grep -q "^kind: ${resource_kind}$" "${rendered_manifest}"; then
        echo "The VM-only chart does not render a ${resource_kind} resource." >&2
        exit 1
    fi
done

for crd_name in podmonitors.monitoring.coreos.com servicemonitors.monitoring.coreos.com; do
    if ! grep -q "^  name: ${crd_name}$" "${rendered_manifest}"; then
        echo "The VM-only chart renders monitoring resources without the ${crd_name} CRD." >&2
        exit 1
    fi
done

helm template monitoring "${chart_dir}" \
    --show-only templates/operator/clusterrole.yaml \
    >"${operator_rbac_manifest}"

for vm_resource in vmanomalyconfigs vmanomalyconfigs/finalizers vmanomalyconfigs/status; do
    if ! grep -q "^      - ${vm_resource}$" "${operator_rbac_manifest}"; then
        echo "The monitoring operator ClusterRole cannot delegate ${vm_resource} permissions." >&2
        exit 1
    fi
done

helm package "${chart_dir}" --destination "${temporary_dir}" >/dev/null
chart_packages=("${temporary_dir}"/*.tgz)
chart_size="$(wc -c <"${chart_packages[0]}")"
maximum_chart_size=750000
if [[ "${chart_size}" -gt "${maximum_chart_size}" ]]; then
    echo "The packaged chart is ${chart_size} bytes; the limit is ${maximum_chart_size} bytes." >&2
    echo "Large charts can exceed the 1 MiB Kubernetes Secret limit for Helm releases." >&2
    exit 1
fi
