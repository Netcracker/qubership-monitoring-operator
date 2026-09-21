# Integration tests Helm values

Robot-based integration tests are optional. Enable them with `integrationTests.install: true` in `charts/qubership-monitoring-operator/values.yaml`.

The chart can pass settings to the test image for S3-compatible result upload and reporting
(same conventions as [qubership-docker-integration-tests](https://github.com/Netcracker/qubership-docker-integration-tests)).
**`ATP_*` S3-related environment variables are injected only when `integrationTests.atpReport.enabled` is `true`**
(values are under `integrationTests.atpReport.atpStorage`).
`ENVIRONMENT_NAME` is always set when integration tests are installed.

When `integrationTests.atpReport.enabled` is `true`, the chart creates Secret `{{ integrationTests.name }}-atp-storage-secret`
from `atpReport.atpStorage.username` and `atpReport.atpStorage.password`, and the pod reads `ATP_STORAGE_USERNAME` / `ATP_STORAGE_PASSWORD` via `secretKeyRef`.

<!-- markdownlint-disable line-length -->
| Helm value                                          | Environment variable        | Description                                                                     |
| --------------------------------------------------- | --------------------------- | ------------------------------------------------------------------------------- |
| `integrationTests.atpReport.enabled`                | `ATP_REPORT_ENABLED`        | When `false`, no ATP S3 env vars are injected.                                  |
| `integrationTests.atpReport.atpStorage.provider`    | `ATP_STORAGE_PROVIDER`      | Storage backend (for example `aws`, `minio`).                                   |
| `integrationTests.atpReport.atpStorage.serverUrl`   | `ATP_STORAGE_SERVER_URL`    | S3 API endpoint URL.                                                            |
| `integrationTests.atpReport.atpStorage.serverUiUrl` | `ATP_STORAGE_SERVER_UI_URL` | Optional storage web UI URL.                                                    |
| `integrationTests.atpReport.atpStorage.bucket`      | `ATP_STORAGE_BUCKET`        | Bucket for uploads; if empty, S3 integration is disabled in the shared scripts. |
| `integrationTests.atpReport.atpStorage.region`      | `ATP_STORAGE_REGION`        | Region for providers that require it.                                           |
| `integrationTests.atpReport.atpStorage.username`    | `ATP_STORAGE_USERNAME`      | Access key; Secret + env when `atpReport.enabled`.                              |
| `integrationTests.atpReport.atpStorage.password`    | `ATP_STORAGE_PASSWORD`      | Secret key; Secret + env when `atpReport.enabled`.                              |
| `integrationTests.atpReportViewUiUrl`               | `ATP_REPORT_VIEW_UI_URL`    | Base URL for viewing reports (for example Allure).                              |
| `integrationTests.environmentName`                  | `ENVIRONMENT_NAME`          | Label for organizing result paths (for example environment or product name).    |
<!-- markdownlint-enable line-length -->

Keep `atpReport.enabled` at `false` unless you enable report upload and supply credentials.
For S3 upload, set bucket, endpoint, region, and credentials as needed.
