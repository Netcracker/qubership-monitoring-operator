#!/usr/bin/env bash

set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_dir="$(cd "${script_dir}/../.." && pwd)"
fixture_dir="$(mktemp -d)"
trap 'rm -rf "${fixture_dir}"' EXIT
manifest_path="charts/qubership-monitoring-operator/crds/monitoring.netcracker.com_platformmonitorings.yaml"

# Copy only what the package needs, so a sibling checkout cannot mask a bad path.
for checkout_name in qubership-monitoring-operator renamed-checkout; do
    checkout_dir="${fixture_dir}/${checkout_name}"
    mkdir -p "${checkout_dir}/$(dirname "${manifest_path}")"
    cp "${repo_dir}/go.mod" "${repo_dir}/go.sum" "${checkout_dir}/"
    cp -R "${repo_dir}/api" "${checkout_dir}/api"
    cp "${repo_dir}/${manifest_path}" "${checkout_dir}/${manifest_path}"
    (
        cd "${checkout_dir}"
        go test ./api/v1 -run '^TestPlatformMonitoringCRDManifest$' -count=1
    )
    # Remove the conventional checkout before testing another name.
    rm -rf "${checkout_dir}"
done

echo "PlatformMonitoring CRD manifest checkout checks passed"
