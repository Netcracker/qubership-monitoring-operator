#!/usr/bin/env bash

set -euo pipefail

chart_dir="${1:-charts/qubership-monitoring-operator}"
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
yq_binary="${script_dir}/../bin/yq-4.53.2"
temporary_dir="$(mktemp -d)"
default_manifest="${temporary_dir}/default.yaml"
override_manifest="${temporary_dir}/override.yaml"
secret_manifest="${temporary_dir}/secret.yaml"
custom_preinstall_manifest="${temporary_dir}/custom-preinstall.yaml"
empty_preinstall_manifest="${temporary_dir}/empty-preinstall.yaml"
precedence_manifest="${temporary_dir}/precedence.yaml"
conflicting_preinstall_manifest="${temporary_dir}/conflicting-preinstall.yaml"
conflicting_preinstall_error="${temporary_dir}/conflicting-preinstall.err"
invalid_config_value_error="${temporary_dir}/invalid-config-value.err"
invalid_secret_value_error="${temporary_dir}/invalid-secret-value.err"
config_override_values="${temporary_dir}/config-override-values.yaml"
secret_override_values="${temporary_dir}/secret-override-values.yaml"
custom_preinstall_values="${temporary_dir}/custom-preinstall-values.yaml"
empty_preinstall_values="${temporary_dir}/empty-preinstall-values.yaml"
precedence_values="${temporary_dir}/precedence-values.yaml"
conflicting_preinstall_values="${temporary_dir}/conflicting-preinstall-values.yaml"
invalid_config_value_values="${temporary_dir}/invalid-config-value-values.yaml"
invalid_secret_value_values="${temporary_dir}/invalid-secret-value-values.yaml"
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
    UNRELATED_CONFIG_NUMBER: 42
    UNRELATED_CONFIG_BOOLEAN: true
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
    UNRELATED_SECRET_NUMBER: 42
    UNRELATED_SECRET_BOOLEAN: true
EOF

cat >"${custom_preinstall_values}" <<'EOF'
grafana:
  extraVars:
    GF_PLUGINS_PREINSTALL_SYNC: private-plugin@1.2.3@https://plugins.internal.example/private-plugin.zip
EOF

cat >"${empty_preinstall_values}" <<'EOF'
grafana:
  extraVars:
    GF_PLUGINS_PREINSTALL_SYNC: ""
EOF

cat >"${precedence_values}" <<'EOF'
grafana:
  extraVars:
    GF_PLUGINS_PREINSTALL_DISABLED: true
    GF_PLUGINS_PREINSTALL_SYNC: private-plugin@1.2.3@https://plugins.internal.example/private-plugin.zip
  extraVarsSecret:
    GF_PLUGINS_PREINSTALL_SYNC: ""
EOF

cat >"${conflicting_preinstall_values}" <<'EOF'
grafana:
  extraVars:
    GF_PLUGINS_PREINSTALL_DISABLED: true
    GF_PLUGINS_PREINSTALL_SYNC: private-plugin@1.2.3@https://plugins.internal.example/private-plugin.zip
EOF

cat >"${invalid_config_value_values}" <<'EOF'
grafana:
  extraVars:
    INVALID_NESTED_VALUE:
      child: value
EOF

cat >"${invalid_secret_value_values}" <<'EOF'
grafana:
  extraVarsSecret:
    INVALID_NESTED_VALUE:
      child: value
EOF

helm template monitoring "${chart_dir}" --namespace monitoring \
    --show-only charts/grafana/templates/configmap-extra-vars.yaml >"${default_manifest}"
helm template monitoring "${chart_dir}" --namespace monitoring \
    --show-only charts/grafana/templates/configmap-extra-vars.yaml \
    --values "${config_override_values}" >"${override_manifest}"
helm template monitoring "${chart_dir}" --namespace monitoring \
    --show-only templates/oauth2-configs/secret-extra-vars-grafana.yaml \
    --values "${secret_override_values}" >"${secret_manifest}"
helm template monitoring "${chart_dir}" --namespace monitoring \
    --show-only charts/grafana/templates/configmap-extra-vars.yaml \
    --values "${custom_preinstall_values}" >"${custom_preinstall_manifest}"
helm template monitoring "${chart_dir}" --namespace monitoring \
    --show-only charts/grafana/templates/configmap-extra-vars.yaml \
    --values "${empty_preinstall_values}" >"${empty_preinstall_manifest}"
helm template monitoring "${chart_dir}" --namespace monitoring \
    --show-only charts/grafana/templates/configmap-extra-vars.yaml \
    --values "${precedence_values}" >"${precedence_manifest}"

if helm template monitoring "${chart_dir}" --namespace monitoring \
    --show-only charts/grafana/templates/configmap-extra-vars.yaml \
    --values "${invalid_config_value_values}" >/dev/null 2>"${invalid_config_value_error}"; then
    echo "grafana.extraVars accepts a nested object instead of a scalar value." >&2
    exit 1
fi
if [[ "$(<"${invalid_config_value_error}")" != *"/extraVars/INVALID_NESTED_VALUE"* ]]; then
    echo "grafana.extraVars failed for an unexpected reason:" >&2
    sed 's/^/  /' "${invalid_config_value_error}" >&2
    exit 1
fi

if helm template monitoring "${chart_dir}" --namespace monitoring \
    --show-only templates/oauth2-configs/secret-extra-vars-grafana.yaml \
    --values "${invalid_secret_value_values}" >/dev/null 2>"${invalid_secret_value_error}"; then
    echo "grafana.extraVarsSecret accepts a nested object instead of a scalar value." >&2
    exit 1
fi
if [[ "$(<"${invalid_secret_value_error}")" != *"/extraVarsSecret/INVALID_NESTED_VALUE"* ]]; then
    echo "grafana.extraVarsSecret failed for an unexpected reason:" >&2
    sed 's/^/  /' "${invalid_secret_value_error}" >&2
    exit 1
fi

required_config_expression='select(.kind == "ConfigMap" and .metadata.name == "grafana-extra-vars") |
    .data.GF_PLUGINS_PREINSTALL_DISABLED == "true" and
    .data.GF_PLUGINS_PUBLIC_KEY_RETRIEVAL_DISABLED == "true" and
    .data.GF_ANALYTICS_CHECK_FOR_UPDATES == "false" and
    .data.GF_ANALYTICS_CHECK_FOR_PLUGIN_UPDATES == "false" and
    .data.GF_ANALYTICS_REPORTING_ENABLED == "false" and
    .data.GF_NEWS_NEWS_FEED_ENABLED == "false"'

if ! "${yq_binary}" eval-all -e "${required_config_expression}" "${default_manifest}" >/dev/null; then
    echo "The default grafana-extra-vars ConfigMap does not enable the offline settings." >&2
    exit 1
fi

export EXPECTED_GRAFANA_RENDERER_URL="http://grafana-image-renderer:8081/render" # NOSONAR -- cluster-local service
export EXPECTED_GRAFANA_CALLBACK_URL="http://grafana-service:3000/"              # NOSONAR -- cluster-local service
if ! "${yq_binary}" eval-all -e \
    'select(.kind == "ConfigMap" and .metadata.name == "grafana-extra-vars") |
        .data.UNRELATED_CONFIG_VALUE == "retained" and
        .data.UNRELATED_CONFIG_NUMBER == "42" and
        .data.UNRELATED_CONFIG_BOOLEAN == "true" and
        .data.GF_PLUGINS_PREINSTALL_DISABLED == "false" and
        .data.GF_PLUGINS_PUBLIC_KEY_RETRIEVAL_DISABLED == "false" and
        .data.GF_ANALYTICS_CHECK_FOR_UPDATES == "true" and
        .data.GF_ANALYTICS_CHECK_FOR_PLUGIN_UPDATES == "true" and
        .data.GF_ANALYTICS_REPORTING_ENABLED == "true" and
        .data.GF_NEWS_NEWS_FEED_ENABLED == "true" and
        .data.GF_RENDERING_SERVER_URL == strenv(EXPECTED_GRAFANA_RENDERER_URL) and
        .data.GF_RENDERING_CALLBACK_URL == strenv(EXPECTED_GRAFANA_CALLBACK_URL)' \
    "${override_manifest}" >/dev/null; then
    echo "The rendered grafana-extra-vars ConfigMap does not preserve unrelated values." >&2
    exit 1
fi

if ! "${yq_binary}" eval-all -e \
    'select(.kind == "ConfigMap" and .metadata.name == "grafana-extra-vars") |
        .data.GF_PLUGINS_PREINSTALL_DISABLED == "true"' \
    "${precedence_manifest}" >/dev/null; then
    echo "An empty Secret preinstall value does not override the ConfigMap list." >&2
    exit 1
fi

if ! "${yq_binary}" eval-all -e \
    'select(.kind == "Secret" and .metadata.name == "grafana-extra-vars-secret") |
        .stringData.UNRELATED_SECRET_VALUE == "retained" and
        .stringData.UNRELATED_SECRET_NUMBER == "42" and
        .stringData.UNRELATED_SECRET_BOOLEAN == "true" and
        .stringData.GF_PLUGINS_PREINSTALL_DISABLED == "false" and
        .stringData.GF_PLUGINS_PUBLIC_KEY_RETRIEVAL_DISABLED == "false" and
        .stringData.GF_ANALYTICS_CHECK_FOR_UPDATES == "true" and
        .stringData.GF_ANALYTICS_CHECK_FOR_PLUGIN_UPDATES == "true" and
        .stringData.GF_ANALYTICS_REPORTING_ENABLED == "true" and
        .stringData.GF_NEWS_NEWS_FEED_ENABLED == "true"' \
    "${secret_manifest}" >/dev/null; then
    echo "The rendered grafana-extra-vars-secret Secret does not preserve explicit overrides." >&2
    exit 1
fi

if ! "${yq_binary}" eval-all -e \
    'select(.kind == "ConfigMap" and .metadata.name == "grafana-extra-vars") |
        .data.GF_PLUGINS_PREINSTALL_DISABLED == "false" and
        .data.GF_PLUGINS_PREINSTALL_SYNC == "private-plugin@1.2.3@https://plugins.internal.example/private-plugin.zip"' \
    "${custom_preinstall_manifest}" >/dev/null; then
    echo "A custom Grafana preinstall list remains disabled." >&2
    exit 1
fi

if ! "${yq_binary}" eval-all -e \
    'select(.kind == "ConfigMap" and .metadata.name == "grafana-extra-vars") |
        .data.GF_PLUGINS_PREINSTALL_DISABLED == "true" and
        .data.GF_PLUGINS_PREINSTALL_SYNC == ""' \
    "${empty_preinstall_manifest}" >/dev/null; then
    echo "An empty Grafana preinstall list enables suggested-plugin preinstallation." >&2
    exit 1
fi

if helm template monitoring "${chart_dir}" --namespace monitoring \
    --show-only charts/grafana/templates/configmap-extra-vars.yaml \
    --values "${conflicting_preinstall_values}" \
    >"${conflicting_preinstall_manifest}" 2>"${conflicting_preinstall_error}"; then
    echo "Grafana accepted a preinstall list while plugin preinstallation was explicitly disabled." >&2
    exit 1
fi

conflicting_preinstall_message="$(<"${conflicting_preinstall_error}")"
case "${conflicting_preinstall_message}" in
*"GF_PLUGINS_PREINSTALL_DISABLED=true conflicts with GF_PLUGINS_PREINSTALL or GF_PLUGINS_PREINSTALL_SYNC"*) ;;
*)
    echo "Grafana rejected conflicting preinstall settings without the expected error." >&2
    exit 1
    ;;
esac

echo "Grafana offline settings render verification passed."
