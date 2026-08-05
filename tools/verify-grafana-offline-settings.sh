#!/usr/bin/env bash

set -euo pipefail

chart_dir="${1:-charts/qubership-monitoring-operator}"
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
yq_binary="${script_dir}/../bin/yq-4.53.2"
temporary_dir="$(mktemp -d)"
default_manifest="${temporary_dir}/default.yaml"
override_manifest="${temporary_dir}/override.yaml"
secret_manifest="${temporary_dir}/secret.yaml"
config_override_values="${temporary_dir}/config-override-values.yaml"
secret_override_values="${temporary_dir}/secret-override-values.yaml"
trap 'rm -rf "${temporary_dir}"' EXIT

if [[ ! -x "${yq_binary}" ]]; then
    echo "The pinned yq binary is missing: ${yq_binary}. Run 'make yq' first." >&2
    exit 1
fi

cat >"${config_override_values}" <<'EOF'
grafana:
  imageRenderer:
    install: true
  extraVars:
    GF_PLUGINS_PREINSTALL_DISABLED: false
    GF_PLUGINS_PUBLIC_KEY_RETRIEVAL_DISABLED: false
    GF_ANALYTICS_CHECK_FOR_UPDATES: true
    GF_ANALYTICS_CHECK_FOR_PLUGIN_UPDATES: true
    GF_ANALYTICS_REPORTING_ENABLED: true
    GF_NEWS_NEWS_FEED_ENABLED: true
    GF_RENDERING_CALLBACK_URL: ""
    GF_RENDERING_SERVER_URL: ""
    UNRELATED_CONFIG_VALUE: retained
EOF

cat >"${secret_override_values}" <<'EOF'
grafana:
  extraVarsSecret:
    GF_PLUGINS_PREINSTALL_DISABLED: false
    GF_PLUGINS_PUBLIC_KEY_RETRIEVAL_DISABLED: false
    GF_ANALYTICS_CHECK_FOR_UPDATES: true
    GF_ANALYTICS_CHECK_FOR_PLUGIN_UPDATES: true
    GF_ANALYTICS_REPORTING_ENABLED: true
    GF_NEWS_NEWS_FEED_ENABLED: true
    UNRELATED_SECRET_VALUE: retained
EOF

helm template monitoring "${chart_dir}" --namespace monitoring \
    --show-only charts/grafana/templates/configmap-extra-vars.yaml >"${default_manifest}"
helm template monitoring "${chart_dir}" --namespace monitoring \
    --show-only charts/grafana/templates/configmap-extra-vars.yaml \
    --values "${config_override_values}" >"${override_manifest}"
helm template monitoring "${chart_dir}" --namespace monitoring \
    --show-only templates/oauth2-configs/secret-extra-vars-grafana.yaml \
    --values "${secret_override_values}" >"${secret_manifest}"

required_config_expression='select(.kind == "ConfigMap" and .metadata.name == "grafana-extra-vars") |
    .data.GF_PLUGINS_PREINSTALL_DISABLED == "true" and
    .data.GF_PLUGINS_PUBLIC_KEY_RETRIEVAL_DISABLED == "true" and
    .data.GF_ANALYTICS_CHECK_FOR_UPDATES == "false" and
    .data.GF_ANALYTICS_CHECK_FOR_PLUGIN_UPDATES == "false" and
    .data.GF_ANALYTICS_REPORTING_ENABLED == "false" and
    .data.GF_NEWS_NEWS_FEED_ENABLED == "false"'

for manifest in "${default_manifest}" "${override_manifest}"; do
    if ! "${yq_binary}" eval-all -e "${required_config_expression}" "${manifest}" >/dev/null; then
        echo "The rendered grafana-extra-vars ConfigMap does not enforce the offline settings." >&2
        exit 1
    fi
done

if ! "${yq_binary}" eval-all -e \
    'select(.kind == "ConfigMap" and .metadata.name == "grafana-extra-vars") |
        .data.UNRELATED_CONFIG_VALUE == "retained" and
        .data.GF_RENDERING_SERVER_URL == "http://grafana-image-renderer:8081/render" and
        .data.GF_RENDERING_CALLBACK_URL == "http://grafana-service:3000/"' \
    "${override_manifest}" >/dev/null; then
    echo "The rendered grafana-extra-vars ConfigMap does not preserve unrelated values." >&2
    exit 1
fi

if ! "${yq_binary}" eval-all -e \
    'select(.kind == "Secret" and .metadata.name == "grafana-extra-vars-secret") |
        .stringData.UNRELATED_SECRET_VALUE == "retained" and
        (.stringData | has("GF_PLUGINS_PREINSTALL_DISABLED") | not) and
        (.stringData | has("GF_PLUGINS_PUBLIC_KEY_RETRIEVAL_DISABLED") | not) and
        (.stringData | has("GF_ANALYTICS_CHECK_FOR_UPDATES") | not) and
        (.stringData | has("GF_ANALYTICS_CHECK_FOR_PLUGIN_UPDATES") | not) and
        (.stringData | has("GF_ANALYTICS_REPORTING_ENABLED") | not) and
        (.stringData | has("GF_NEWS_NEWS_FEED_ENABLED") | not)' \
    "${secret_manifest}" >/dev/null; then
    echo "The rendered grafana-extra-vars-secret Secret contains reserved keys or loses unrelated values." >&2
    exit 1
fi

echo "Grafana offline settings render verification passed."
