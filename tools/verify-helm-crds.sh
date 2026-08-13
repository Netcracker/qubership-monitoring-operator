#!/usr/bin/env bash

set -euo pipefail

chart_dir="${1:-charts/qubership-monitoring-operator}"
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
yq_binary="${script_dir}/../bin/yq-4.53.2"
temporary_dir="$(mktemp -d)"
rendered_manifest="${temporary_dir}/manifest.yaml"
operator_rbac_manifest="${temporary_dir}/operator-rbac.yaml"
operator_role_manifest="${temporary_dir}/operator-role.yaml"
cleanup_role_manifest="${temporary_dir}/cleanup-role.yaml"
cleanup_command="${temporary_dir}/cleanup-command.sh"
grafana_cleanup_role_manifest="${temporary_dir}/grafana-cleanup-role.yaml"
grafana_cleanup_command="${temporary_dir}/grafana-cleanup-command.sh"
rbac_cleanup_role_manifest="${temporary_dir}/rbac-cleanup-clusterrole.yaml"
rbac_cleanup_command="${temporary_dir}/rbac-cleanup-command.sh"
other_namespace_manifest="${temporary_dir}/other-namespace.yaml"
other_namespace_rbac_cleanup_command="${temporary_dir}/other-namespace-rbac-cleanup-command.sh"
root_crd_dir="${chart_dir}/crds"
prometheus_crd_dir="${chart_dir}/charts/prometheus/crds"
victoriametrics_crd_dir="${chart_dir}/charts/victoriametrics/crds"
standalone_crd_dir="${chart_dir}/../qubership-monitoring-crds/crds"
grafana_cluster_role="${script_dir}/../controllers/grafana-operator/assets/grafana-operator/cluster-role.yaml"
grafana_role="${script_dir}/../controllers/grafana-operator/assets/grafana-operator/role.yaml"
prometheus_cluster_role="${script_dir}/../controllers/prometheus-operator/assets/cluster-role.yaml"
prometheus_role="${script_dir}/../controllers/prometheus-operator/assets/role.yaml"
victoriametrics_cluster_role="${script_dir}/../controllers/victoriametrics/vmoperator/assets/cluster-role.yaml"
victoriametrics_role="${script_dir}/../controllers/victoriametrics/vmoperator/assets/role.yaml"
victoriametrics_deployment="${script_dir}/../controllers/victoriametrics/vmoperator/assets/deployment.yaml"
trap 'rm -rf "${temporary_dir}"' EXIT

if [[ ! -x "${yq_binary}" ]]; then
    echo "The pinned yq binary is missing: ${yq_binary}. Run 'make yq' first." >&2
    exit 1
fi

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

verify_crd_compaction() {
    local crd_dir
    local crd
    local -a crd_dirs=(
        "${root_crd_dir}"
        "${chart_dir}/charts/grafana/crds"
        "${prometheus_crd_dir}"
        "${chart_dir}/charts/prometheusAdapter/crds"
        "${victoriametrics_crd_dir}"
    )

    for crd_dir in "${crd_dirs[@]}"; do
        local -a crds=("${crd_dir}"/*.yaml "${crd_dir}"/*.yml)
        for crd in "${crds[@]}"; do
            if ! "${yq_binary}" eval -e \
                '(.metadata.annotations // {}) | has("helm.sh/hook") | not' \
                "${crd}" >/dev/null; then
                echo "The CRD ${crd} contains the obsolete Helm 2 crd-install hook." >&2
                exit 1
            fi
            if ! "${yq_binary}" eval -e \
                '[.spec.versions[]?.schema.openAPIV3Schema | .. | select(tag == "!!map" and has("description") and (.description | tag == "!!str"))] | length == 0' \
                "${crd}" >/dev/null; then
                echo "The CRD ${crd} contains OpenAPI schema documentation." >&2
                exit 1
            fi
        done
    done
}

verify_crd_compaction

verify_standalone_crds() {
    local expected_dir="${temporary_dir}/expected-standalone-crds"
    local duplicate_names
    local -a standalone_crds=("${standalone_crd_dir}"/*.yaml "${standalone_crd_dir}"/*.yml)
    local -a standalone_crd_names
    mkdir -p "${expected_dir}"

    cp "${root_crd_dir}"/* "${expected_dir}/"
    cp "${chart_dir}/charts/grafana/crds/"* "${expected_dir}/"
    cp "${prometheus_crd_dir}/"* "${expected_dir}/"
    cp "${chart_dir}/charts/prometheusAdapter/crds/"* "${expected_dir}/"
    cp "${victoriametrics_crd_dir}/"* "${expected_dir}/"

    if ! diff -qr "${expected_dir}" "${standalone_crd_dir}" >/dev/null; then
        echo "The standalone CRD chart differs from the generated aggregate." >&2
        diff -qr "${expected_dir}" "${standalone_crd_dir}" >&2 || true
        exit 1
    fi

    mapfile -t standalone_crd_names < <(
        "${yq_binary}" eval-all '.metadata.name' "${standalone_crds[@]}" |
            grep -v '^---$'
    )
    if [[ "${#standalone_crd_names[@]}" -ne "${#standalone_crds[@]}" ]]; then
        echo "The standalone chart must contain exactly one CRD document per file." >&2
        exit 1
    fi

    duplicate_names="$(
        printf '%s\n' "${standalone_crd_names[@]}" |
            sort |
            uniq -d
    )"
    if [[ -n "${duplicate_names}" ]]; then
        echo "The standalone chart contains duplicate CRD identities:" >&2
        printf '%s\n' "${duplicate_names}" >&2
        exit 1
    fi
}

verify_standalone_crds

verify_platform_monitoring_compatibility() {
    local platform_monitoring_crd="${root_crds[0]}"
    local prometheus_schema='.spec.versions[0].schema.openAPIV3Schema.properties.spec.properties.prometheus.properties'
    local vm_alertmanager_storage='.spec.versions[0].schema.openAPIV3Schema.properties.spec.properties.victoriametrics.properties.vmAlertManager.properties.storage.properties'

    if ! "${yq_binary}" eval -e \
        "${prometheus_schema}.remoteRead.items.properties.url | has(\"pattern\") | not" \
        "${platform_monitoring_crd}" >/dev/null; then
        echo "PlatformMonitoring remoteRead.url is stricter than the existing installation contract." >&2
        exit 1
    fi

    for remote_endpoint in remoteRead remoteWrite; do
        if ! "${yq_binary}" eval -e \
            "${prometheus_schema}.${remote_endpoint}.items.properties.oauth2.properties.tokenUrl | .minLength == 1 and (has(\"pattern\") | not)" \
            "${platform_monitoring_crd}" >/dev/null; then
            echo "PlatformMonitoring ${remote_endpoint}.oauth2.tokenUrl is incompatible with existing installations." >&2
            exit 1
        fi
    done

    if ! "${yq_binary}" eval -e \
        "${vm_alertmanager_storage}.disableMountSubPath.type == \"boolean\"" \
        "${platform_monitoring_crd}" >/dev/null; then
        echo "PlatformMonitoring no longer preserves vmAlertManager.storage.disableMountSubPath." >&2
        exit 1
    fi
}

verify_platform_monitoring_compatibility

verify_kubernetes_version_guard() {
    local chart="$1"

    if helm template version-check "${chart}" --kube-version 1.24.17 >/dev/null 2>&1; then
        echo "The chart ${chart} accepts unsupported Kubernetes 1.24." >&2
        exit 1
    fi

    if ! helm template version-check "${chart}" --kube-version 1.25.0 >/dev/null; then
        echo "The chart ${chart} rejects supported Kubernetes 1.25." >&2
        exit 1
    fi
}

verify_kubernetes_version_guard "${chart_dir}"
verify_kubernetes_version_guard "${chart_dir}/../qubership-monitoring-crds"

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
    --set prometheus.install=true \
    >"${rendered_manifest}"

expected_images=(
    "docker.io/grafana/grafana:12.4.3"
    "docker.io/prom/alertmanager:v0.33.1"
    "docker.io/prom/prometheus:v3.13.1"
    "docker.io/victoriametrics/operator:v0.73.1"
    "docker.io/victoriametrics/victoria-metrics:v1.148.0"
    "docker.io/victoriametrics/vmagent:v1.148.0"
    "docker.io/victoriametrics/vmalert:v1.148.0"
    "docker.io/victoriametrics/vmauth:v1.148.0"
    "quay.io/grafana-operator/grafana-operator:v5.24.0"
    "quay.io/prometheus-operator/prometheus-config-reloader:v0.93.0"
    "quay.io/prometheus-operator/prometheus-operator:v0.93.0"
)

for expected_image in "${expected_images[@]}"; do
    if ! grep -Fq "${expected_image}" "${rendered_manifest}"; then
        echo "The chart does not render the expected migrated image ${expected_image}." >&2
        exit 1
    fi
done

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
helm template monitoring "${chart_dir}" \
    --set global.privilegedRights=false \
    --show-only templates/operator/role.yaml \
    >"${operator_role_manifest}"

for vm_resource in vmanomalyconfigs vmanomalyconfigs/finalizers vmanomalyconfigs/status; do
    if ! grep -q "^      - ${vm_resource}$" "${operator_rbac_manifest}"; then
        echo "The monitoring operator ClusterRole cannot delegate ${vm_resource} permissions." >&2
        exit 1
    fi
done

verify_rbac_resource() {
    local manifest="$1"
    local api_group="$2"
    local resource="$3"
    local verb="$4"

    if ! "${yq_binary}" eval -e \
        ".rules[] | select(.apiGroups | any_c(. == \"${api_group}\")) | select(.resources | any_c(. == \"${resource}\")) | select(.verbs | any_c(. == \"${verb}\"))" \
        "${manifest}" >/dev/null; then
        echo "${manifest} does not grant ${verb} on ${api_group}/${resource}." >&2
        exit 1
    fi
}

verify_exact_rbac_verbs() {
    local manifest="$1"
    local api_group="$2"
    local resource="$3"
    local expected_verbs="$4"
    local matches
    local match_count
    local actual_verbs

    matches="[.rules[] | select(.apiGroups | any_c(. == \"${api_group}\")) | select(.resources | any_c(. == \"${resource}\"))]"
    match_count="$("${yq_binary}" eval "${matches} | length" "${manifest}")"
    if [[ "${match_count}" == "1" ]]; then
        actual_verbs="$("${yq_binary}" eval "${matches} | .[0].verbs | sort | join(\",\")" "${manifest}")"
    fi
    if [[ "${match_count}" != "1" || "${actual_verbs:-}" != "${expected_verbs}" ]]; then
        echo "${manifest} does not grant exactly ${expected_verbs} on ${api_group}/${resource}." >&2
        exit 1
    fi
}

verify_rbac_can_delegate() {
    local parent_manifest="$1"
    local delegated_manifest="$2"
    local api_group
    local resource
    local verb

    while IFS='|' read -r api_group resource verb; do
        verify_rbac_resource "${parent_manifest}" "${api_group}" "${resource}" "${verb}"
    done < <(
        "${yq_binary}" eval -r \
            ".rules[] | select(.resources != null) | .apiGroups[] as \$apiGroup |
            .resources[] as \$resource | .verbs[] | [\$apiGroup, \$resource, .] | join(\"|\")" \
            "${delegated_manifest}"
    )
}

verify_rendered_resource_count() {
    local manifest="$1"
    local selector="$2"
    local expected_count="$3"
    local description="$4"
    local actual_count

    actual_count="$("${yq_binary}" eval-all "[select(${selector})] | length" "${manifest}")"
    if [[ "${actual_count}" != "${expected_count}" ]]; then
        echo "${manifest} renders ${actual_count} ${description}; expected ${expected_count}." >&2
        exit 1
    fi
}

verify_text_contains() {
    local file="$1"
    local expected="$2"
    local description="$3"

    if ! grep -Fq -- "${expected}" "${file}"; then
        echo "${file} does not contain ${description}: ${expected}" >&2
        exit 1
    fi
}

verify_text_excludes() {
    local file="$1"
    local forbidden="$2"
    local description="$3"

    if grep -Fq -- "${forbidden}" "${file}"; then
        echo "${file} contains forbidden ${description}: ${forbidden}" >&2
        exit 1
    fi
}

for grafana_rbac in "${grafana_cluster_role}" "${grafana_role}"; do
    verify_rbac_resource "${grafana_rbac}" grafana.integreatly.org '*/finalizers' update
    verify_rbac_resource "${grafana_rbac}" grafana.integreatly.org '*/finalizers' patch
done
verify_rbac_resource "${grafana_cluster_role}" events.k8s.io events create
verify_rbac_resource "${grafana_cluster_role}" events.k8s.io events patch
verify_rbac_resource "${grafana_cluster_role}" gateway.networking.k8s.io httproutes create
verify_rbac_resource "${grafana_role}" route.openshift.io routes/custom-host create
verify_exact_rbac_verbs "${grafana_role}" "" events "create,patch"
verify_exact_rbac_verbs "${grafana_role}" events.k8s.io events "create,patch"
verify_exact_rbac_verbs "${grafana_role}" gateway.networking.k8s.io httproutes \
    "create,delete,get,list,patch,update,watch"
verify_exact_rbac_verbs "${grafana_role}" grafana.integreatly.org grafanamanifests \
    "create,delete,get,list,patch,update,watch"
verify_exact_rbac_verbs "${grafana_role}" grafana.integreatly.org grafanamanifests/status \
    "create,delete,get,list,patch,update,watch"

for prometheus_rbac in "${prometheus_cluster_role}" "${prometheus_role}"; do
    verify_rbac_resource "${prometheus_rbac}" monitoring.coreos.com prometheusagents/finalizers update
    verify_rbac_resource "${prometheus_rbac}" monitoring.coreos.com scrapeconfigs/status update
done
verify_rbac_resource "${prometheus_cluster_role}" events.k8s.io events create
verify_rbac_resource "${prometheus_cluster_role}" events.k8s.io events patch
verify_exact_rbac_verbs "${victoriametrics_cluster_role}" autoscaling.k8s.io verticalpodautoscalers \
    "create,delete,get,list,update,watch"
verify_exact_rbac_verbs "${victoriametrics_cluster_role}" gateway.networking.k8s.io httproutes \
    "create,delete,get,list,update,watch"
verify_exact_rbac_verbs "${victoriametrics_role}" autoscaling.k8s.io verticalpodautoscalers \
    "create,delete,get,list,update,watch"
verify_exact_rbac_verbs "${victoriametrics_role}" gateway.networking.k8s.io httproutes \
    "create,delete,get,list,update,watch"
verify_exact_rbac_verbs "${victoriametrics_role}" policy poddisruptionbudgets \
    "create,delete,get,list,patch,update,watch"
verify_exact_rbac_verbs "${victoriametrics_role}" policy poddisruptionbudgets/finalizers \
    "create,delete,get,list,patch,update,watch"
verify_exact_rbac_verbs "${victoriametrics_role}" autoscaling horizontalpodautoscalers \
    "create,delete,get,list,update,watch"
verify_exact_rbac_verbs "${victoriametrics_role}" discovery.k8s.io endpointslices "get,list,watch"
verify_exact_rbac_verbs "${victoriametrics_role}" apps replicasets "get,list,watch"
verify_exact_rbac_verbs "${victoriametrics_role}" apps daemonsets/status \
    "create,delete,get,list,patch,update,watch"
verify_exact_rbac_verbs "${victoriametrics_role}" apps deployments/status \
    "create,delete,get,list,patch,update,watch"
verify_exact_rbac_verbs "${victoriametrics_role}" operator.victoriametrics.com vmanomalyconfigs \
    "create,delete,get,list,patch,update,watch"
verify_exact_rbac_verbs "${victoriametrics_role}" operator.victoriametrics.com vmanomalyconfigs/finalizers \
    "create,delete,get,list,patch,update,watch"
verify_exact_rbac_verbs "${victoriametrics_role}" operator.victoriametrics.com vmanomalyconfigs/status \
    "create,delete,get,list,patch,update,watch"
for prometheus_resource in alertmanagerconfigs podmonitors probes prometheusrules scrapeconfigs servicemonitors; do
    verify_exact_rbac_verbs "${victoriametrics_role}" monitoring.coreos.com "${prometheus_resource}" \
        "get,list,watch"
done
for prometheus_finalizer in alertmanagerconfigs/finalizers podmonitors/finalizers probes/finalizers \
    prometheusrules/finalizers servicemonitors/finalizers thanosrulers/finalizers; do
    verify_exact_rbac_verbs "${victoriametrics_role}" monitoring.coreos.com "${prometheus_finalizer}" \
        "create,delete,get,list,patch,update,watch"
done
verify_rbac_can_delegate "${operator_role_manifest}" "${grafana_role}"
verify_rbac_can_delegate "${operator_role_manifest}" "${victoriametrics_role}"
verify_rbac_resource "${operator_rbac_manifest}" grafana.integreatly.org '*/finalizers' update
verify_rbac_resource "${operator_rbac_manifest}" grafana.integreatly.org '*/finalizers' patch
verify_rbac_resource "${operator_rbac_manifest}" events.k8s.io events create
verify_rbac_resource "${operator_rbac_manifest}" events.k8s.io events patch
verify_rbac_resource "${operator_rbac_manifest}" route.openshift.io routes patch
verify_rbac_resource "${operator_rbac_manifest}" route.openshift.io routes/custom-host patch
verify_rbac_resource "${operator_role_manifest}" grafana.integreatly.org '*/finalizers' update
verify_rbac_resource "${operator_role_manifest}" grafana.integreatly.org '*/finalizers' patch
verify_rbac_resource "${operator_role_manifest}" route.openshift.io routes patch
verify_rbac_resource "${operator_role_manifest}" route.openshift.io routes/custom-host patch
verify_exact_rbac_verbs "${operator_role_manifest}" events.k8s.io events "create,patch"
verify_exact_rbac_verbs "${operator_role_manifest}" grafana.integreatly.org grafanamanifests \
    "create,delete,get,list,patch,update,watch"
verify_exact_rbac_verbs "${operator_role_manifest}" grafana.integreatly.org grafanamanifests/status \
    "create,delete,get,list,patch,update,watch"
verify_exact_rbac_verbs "${operator_role_manifest}" autoscaling.k8s.io verticalpodautoscalers \
    "create,delete,get,list,update,watch"
verify_exact_rbac_verbs "${operator_role_manifest}" gateway.networking.k8s.io httproutes \
    "create,delete,get,list,patch,update,watch"
verify_exact_rbac_verbs "${operator_role_manifest}" operator.victoriametrics.com vmanomalyconfigs \
    "create,delete,get,list,patch,update,watch"
verify_exact_rbac_verbs "${operator_role_manifest}" operator.victoriametrics.com vmanomalyconfigs/finalizers \
    "create,delete,get,list,patch,update,watch"
verify_exact_rbac_verbs "${operator_role_manifest}" operator.victoriametrics.com vmanomalyconfigs/status \
    "create,delete,get,list,patch,update,watch"
verify_exact_rbac_verbs "${operator_rbac_manifest}" autoscaling.k8s.io verticalpodautoscalers \
    "create,delete,get,list,update,watch"
verify_exact_rbac_verbs "${operator_rbac_manifest}" gateway.networking.k8s.io httproutes \
    "create,delete,get,list,patch,update,watch"
verify_exact_rbac_verbs "${operator_rbac_manifest}" apps deployments/scale "get,patch,update"
verify_exact_rbac_verbs "${operator_role_manifest}" apps deployments/scale "get,patch,update"

helm template monitoring "${chart_dir}" >"${rendered_manifest}"
verify_rendered_resource_count "${rendered_manifest}" \
    '.kind == "Job" and
    .metadata.labels."app.kubernetes.io/component" == "etcd-certs-to-secret" and
    .metadata.annotations."helm.sh/hook-delete-policy" == "before-hook-creation,hook-failed"' \
    1 "revision-scoped etcd certificate hook Jobs with safe replacement policy"
"${yq_binary}" eval-all \
    'select(.kind == "Role" and .metadata.name == "default-cleanup-hook")' \
    "${rendered_manifest}" >"${cleanup_role_manifest}"
verify_exact_rbac_verbs "${cleanup_role_manifest}" apps deployments/scale "get,patch,update"
"${yq_binary}" eval-all \
    'select(.kind == "Role" and .metadata.name == "monitoring-grafana-cleanup-hook")' \
    "${rendered_manifest}" >"${grafana_cleanup_role_manifest}"
verify_exact_rbac_verbs "${grafana_cleanup_role_manifest}" apps deployments "get,list,patch,update,watch"
verify_exact_rbac_verbs "${grafana_cleanup_role_manifest}" apps deployments/scale "get,patch,update"
for grafana_cleanup_resource in grafanadashboards grafanadatasources grafanas; do
    verify_exact_rbac_verbs "${grafana_cleanup_role_manifest}" grafana.integreatly.org \
        "${grafana_cleanup_resource}" "delete,deletecollection,get,list,watch"
done
verify_rendered_resource_count "${rendered_manifest}" \
    '.kind == "Job" and .metadata.name == "monitoring-grafana-cleanup-hook" and
    .metadata.annotations."helm.sh/hook" == "pre-delete" and
    .metadata.annotations."helm.sh/hook-weight" == "-10" and
    .metadata.annotations."helm.sh/hook-delete-policy" == "before-hook-creation,hook-failed" and
    .metadata.annotations."helm.sh/resource-policy" == "keep" and
    .spec.template.spec.serviceAccountName == "monitoring-grafana-cleanup-hook" and
    .spec.ttlSecondsAfterFinished == 180' \
    1 "correctly configured Grafana cleanup Jobs"
verify_rendered_resource_count "${rendered_manifest}" \
    '.kind == "Job" and .metadata.annotations."helm.sh/hook-weight" == "-5" and
    .metadata.annotations."helm.sh/hook-delete-policy" == "before-hook-creation,hook-failed" and
    .metadata.annotations."helm.sh/resource-policy" == "keep" and
    .spec.ttlSecondsAfterFinished == 180' \
    1 "correctly configured VictoriaMetrics cleanup Jobs"
verify_rendered_resource_count "${rendered_manifest}" \
    '.metadata.name == "monitoring-grafana-cleanup-hook" and
    (.kind == "ServiceAccount" or .kind == "Role" or .kind == "RoleBinding") and
    .metadata.annotations."helm.sh/hook" == "pre-delete" and
    .metadata.annotations."helm.sh/hook-weight" == "-11" and
    .metadata.annotations."helm.sh/hook-delete-policy" == "before-hook-creation,hook-succeeded" and
    .metadata.annotations."helm.sh/resource-policy" == "keep"' \
    3 "Argo CD-safe Grafana cleanup RBAC hooks"
verify_rendered_resource_count "${rendered_manifest}" \
    '.metadata.name == "default-cleanup-hook" and
    (.kind == "ServiceAccount" or .kind == "Role" or .kind == "RoleBinding") and
    .metadata.annotations."helm.sh/hook" == "pre-delete" and
    .metadata.annotations."helm.sh/hook-weight" == "-6" and
    .metadata.annotations."helm.sh/hook-delete-policy" == "before-hook-creation,hook-succeeded" and
    .metadata.annotations."helm.sh/resource-policy" == "keep"' \
    3 "Argo CD-safe VictoriaMetrics cleanup RBAC hooks"
verify_rendered_resource_count "${rendered_manifest}" \
    '.metadata.name == "monitoring-rbac-cleanup-hook" and
    (.kind == "ServiceAccount" or .kind == "ClusterRole" or .kind == "ClusterRoleBinding") and
    .metadata.annotations."helm.sh/hook" == "post-delete" and
    .metadata.annotations."helm.sh/hook-weight" == "-1" and
    .metadata.annotations."helm.sh/hook-delete-policy" == "before-hook-creation,hook-succeeded"' \
    3 "post-delete cluster RBAC cleanup hooks"
verify_rendered_resource_count "${rendered_manifest}" \
    '.kind == "ClusterRoleBinding" and .metadata.name == "monitoring-rbac-cleanup-hook" and
    .roleRef.apiGroup == "rbac.authorization.k8s.io" and .roleRef.kind == "ClusterRole" and
    .roleRef.name == "monitoring-rbac-cleanup-hook" and
    (.subjects | any_c(.kind == "ServiceAccount" and .name == "monitoring-rbac-cleanup-hook" and
    .namespace == "default"))' \
    1 "post-delete cleanup ServiceAccount bindings"
verify_rendered_resource_count "${rendered_manifest}" \
    '.kind == "Job" and .metadata.name == "monitoring-rbac-cleanup-hook" and
    .metadata.annotations."helm.sh/hook" == "post-delete" and
    .metadata.annotations."helm.sh/hook-weight" == "0" and
    .metadata.annotations."helm.sh/hook-delete-policy" == "before-hook-creation,hook-failed" and
    .spec.template.spec.serviceAccountName == "monitoring-rbac-cleanup-hook" and
    (.spec.template.spec.containers |
    any_c(.name == "kubectl" and .resources.limits.memory == "256Mi")) and
    .spec.ttlSecondsAfterFinished == 180' \
    1 "post-delete cluster RBAC cleanup Jobs"
"${yq_binary}" eval-all \
    'select(.kind == "ClusterRole" and .metadata.name == "monitoring-rbac-cleanup-hook")' \
    "${rendered_manifest}" >"${rbac_cleanup_role_manifest}"
verify_exact_rbac_verbs "${rbac_cleanup_role_manifest}" rbac.authorization.k8s.io clusterrolebindings \
    "delete,get,list,watch"
verify_exact_rbac_verbs "${rbac_cleanup_role_manifest}" rbac.authorization.k8s.io clusterroles \
    "delete,get,list,watch"
"${yq_binary}" eval-all \
    'select(.kind == "Job" and .metadata.name == "default-cleanup-hook") |
    .spec.template.spec.containers[] | select(.name == "kubectl") | .command[-1]' \
    "${rendered_manifest}" >"${cleanup_command}"
verify_text_excludes "${cleanup_command}" "clusterrolebindings" "cluster RBAC deletion before operators stop"
verify_text_excludes "${cleanup_command}" "clusterroles" "cluster RBAC deletion before operators stop"
"${yq_binary}" eval-all \
    'select(.kind == "Job" and .metadata.name == "monitoring-rbac-cleanup-hook") |
    .spec.template.spec.containers[] | select(.name == "kubectl") | .command[-1]' \
    "${rendered_manifest}" >"${rbac_cleanup_command}"
verify_text_contains "${rbac_cleanup_command}" \
    "clusterrolebindings.rbac.authorization.k8s.io" "qualified ClusterRoleBinding cleanup"
verify_text_contains "${rbac_cleanup_command}" \
    "clusterroles.rbac.authorization.k8s.io" "qualified ClusterRole cleanup"
verify_text_contains "${rbac_cleanup_command}" \
    "app.kubernetes.io/managed-by-operator=monitoring-operator" "operator-managed RBAC selector"
verify_text_contains "${rbac_cleanup_command}" \
    "monitoring.netcracker.com/installation-namespace=default" "installation-scoped RBAC selector"
helm template monitoring "${chart_dir}" --namespace other-monitoring >"${other_namespace_manifest}"
"${yq_binary}" eval-all \
    'select(.kind == "Job" and .metadata.name == "monitoring-rbac-cleanup-hook") |
    .spec.template.spec.containers[] | select(.name == "kubectl") | .command[-1]' \
    "${other_namespace_manifest}" >"${other_namespace_rbac_cleanup_command}"
verify_text_contains "${other_namespace_rbac_cleanup_command}" \
    "monitoring.netcracker.com/installation-namespace=other-monitoring" "second installation RBAC selector"
verify_text_excludes "${other_namespace_rbac_cleanup_command}" \
    "monitoring.netcracker.com/installation-namespace=default" "first installation RBAC selector"
verify_rendered_resource_count "${rendered_manifest}" \
    '.kind == "GrafanaDashboard" and
    (.metadata.name == "backup-daemon" or .metadata.name == "kafka-java-clients") and
    .metadata.annotations."helm.sh/resource-policy" == "keep" and .spec.resyncPeriod == "10m"' \
    2 "retained Helm-owned Grafana dashboards"

helm template monitoring "${chart_dir}" \
    --set blackboxExporter.install=true \
    --set certExporter.install=true \
    --set cloudEventsExporter.install=true \
    --set commonDashboards.install=true \
    --set goldpingerExporter.install=true \
    --set graphite_remote_adapter.install=true \
    --set networkLatencyExporter.install=true \
    --set promxy.install=true \
    --set sslExporter.install=true \
    --set versionExporter.install=true \
    >"${rendered_manifest}"
verify_rendered_resource_count "${rendered_manifest}" \
    '.apiVersion == "grafana.integreatly.org/v1beta1" and
    .kind == "GrafanaDashboard"' \
    12 "optional Grafana dashboards"
verify_rendered_resource_count "${rendered_manifest}" \
    '.apiVersion == "grafana.integreatly.org/v1beta1" and
    .kind == "GrafanaDashboard" and .spec.instanceSelector == null' \
    0 "optional Grafana dashboards without an instance selector"

"${yq_binary}" eval-all \
    'select(.kind == "Job" and .metadata.name == "monitoring-grafana-cleanup-hook") |
    .spec.template.spec.containers[] | select(.name == "kubectl") | .command[-1]' \
    "${rendered_manifest}" >"${grafana_cleanup_command}"
verify_text_excludes "${grafana_cleanup_command}" "--all" "namespace-wide Grafana deletion"
verify_text_excludes "${grafana_cleanup_command}" "grafanadashboards," "unqualified Grafana resource list"
for user_owned_grafana_resource in \
    grafanaalertrulegroups \
    grafanacontactpoints \
    grafanafolders \
    grafanalibrarypanels \
    grafanamanifests \
    grafanamutetimings \
    grafananotificationpolicies \
    grafananotificationpolicyroutes \
    grafananotificationtemplates \
    grafanaserviceaccounts; do
    verify_text_excludes "${grafana_cleanup_command}" "${user_owned_grafana_resource}" \
        "user-owned Grafana resource deletion"
done
for required_grafana_cleanup_token in \
    grafanadashboards.grafana.integreatly.org \
    grafanadatasources.grafana.integreatly.org \
    grafanas.grafana.integreatly.org \
    app.kubernetes.io/managed-by-operator=monitoring-operator \
    app.kubernetes.io/component=grafana-operator \
    backup-daemon \
    kafka-java-clients; do
    verify_text_contains "${grafana_cleanup_command}" "${required_grafana_cleanup_token}" \
        "required scoped Grafana cleanup token"
done
verify_text_contains "${grafana_cleanup_command}" \
    "grafanadatasources.grafana.integreatly.org --selector=app.kubernetes.io/managed-by-operator=monitoring-operator" \
    "reserved Grafana datasource cleanup selector"
verify_text_excludes "${grafana_cleanup_command}" \
    "grafanadatasources.grafana.integreatly.org --selector=app.kubernetes.io/managed-by=monitoring-operator" \
    "generic Grafana datasource cleanup selector"

helm template monitoring "${chart_dir}" \
    --set grafana.install=false \
    --set grafana.operator.install=false \
    --set commonDashboards.install=true \
    >"${rendered_manifest}"
"${yq_binary}" eval-all \
    'select(.kind == "Job" and .metadata.name == "monitoring-grafana-cleanup-hook") |
    .spec.template.spec.containers[] | select(.name == "kubectl") | .command[-1]' \
    "${rendered_manifest}" >"${grafana_cleanup_command}"
verify_text_contains "${grafana_cleanup_command}" "grafanadashboards.grafana.integreatly.org" \
    "qualified dashboard cleanup in common-dashboard-only mode"
verify_text_contains "${grafana_cleanup_command}" "app.kubernetes.io/managed-by-operator=monitoring-operator" \
    "monitoring dashboard selector in common-dashboard-only mode"
verify_text_contains "${grafana_cleanup_command}" "backup-daemon" \
    "backup dashboard cleanup in common-dashboard-only mode"
verify_text_contains "${grafana_cleanup_command}" "kafka-java-clients" \
    "Kafka dashboard cleanup in common-dashboard-only mode"
verify_text_excludes "${grafana_cleanup_command}" "grafanadatasources.grafana.integreatly.org" \
    "datasource deletion in common-dashboard-only mode"
verify_text_excludes "${grafana_cleanup_command}" "grafanas.grafana.integreatly.org" \
    "Grafana instance deletion in common-dashboard-only mode"

helm template monitoring "${chart_dir}" \
    --set global.privilegedRights=false \
    >"${rendered_manifest}"
verify_rendered_resource_count "${rendered_manifest}" \
    '.kind == "Job" and .metadata.annotations."helm.sh/hook-weight" == "-5"' \
    1 "namespaced VictoriaMetrics cleanup Jobs in non-privileged mode"
verify_rendered_resource_count "${rendered_manifest}" \
    '.metadata.name == "default-cleanup-hook" and
    (.kind == "ClusterRole" or .kind == "ClusterRoleBinding")' \
    0 "VictoriaMetrics cleanup cluster RBAC resources in non-privileged mode"
verify_rendered_resource_count "${rendered_manifest}" \
    '.metadata.name == "monitoring-rbac-cleanup-hook"' \
    0 "post-delete cluster RBAC cleanup resources in non-privileged mode"

helm template monitoring "${chart_dir}" \
    --set victoriametrics.cleanup.deleteCRs=false \
    >"${rendered_manifest}"
verify_rendered_resource_count "${rendered_manifest}" \
    '.metadata.name == "default-cleanup-hook"' \
    0 "VictoriaMetrics cleanup resources when CR deletion is disabled"
verify_rendered_resource_count "${rendered_manifest}" \
    '.metadata.name == "monitoring-rbac-cleanup-hook"' \
    4 "post-delete cluster RBAC cleanup resources when VM CR deletion is disabled"

if ! "${yq_binary}" eval -e \
    '.spec.template.spec.containers[] | select(.name == "victoriametrics-operator") | .image == "victoriametrics/operator:v0.73.1"' \
    "${victoriametrics_deployment}" >/dev/null; then
    echo "The VictoriaMetrics operator asset does not use v0.73.1." >&2
    exit 1
fi

helm package "${chart_dir}" --destination "${temporary_dir}" >/dev/null
chart_packages=("${temporary_dir}"/*.tgz)
chart_size="$(wc -c <"${chart_packages[0]}")"
maximum_chart_size=750000
if [[ "${chart_size}" -gt "${maximum_chart_size}" ]]; then
    echo "The packaged chart is ${chart_size} bytes; the limit is ${maximum_chart_size} bytes." >&2
    echo "Large charts can exceed the 1 MiB Kubernetes Secret limit for Helm releases." >&2
    exit 1
fi
