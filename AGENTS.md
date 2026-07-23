# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Common commands

Build / test (from repository root):

- `make generate` — runs `controller-gen` to regenerate project CRDs (into `charts/qubership-monitoring-operator/crds/`) and deepcopy methods. Must be re-run after any change to `api/v1/*.go`.
- `make update-operators-crds` — refreshes managed-operator CRDs in their canonical subchart `crds/` directories.
- `make update-crds` — syncs canonical CRDs into `charts/qubership-monitoring-crds/crds/`.
- `make verify-generated` — regenerates all CRDs and fails if generated artifacts drift.
- `make build-binary` — builds `bin/manager` (CGO disabled). Runs `generate`, `fmt`, `vet` first.
- `make test` — unit tests. Runs `go test -race -vet=off --shuffle=on ./... -count=1` across all packages **except** `/test/envtests`.
- `make envtest` — runs envtests against Kubernetes 1.25 by default; override `ENVTEST_K8S_VERSION` to test another version.
- `make run` — runs the operator locally against the cluster in `~/.kube/config`.
- `make image` — builds the Docker image (tag `qubership-monitoring-operator`).
- `make docs` — compatibility alias for `make update-crds`.
- `make all` — full pipeline: `generate test build-binary image docs archives`.

Running a single Go test: `go test -v -run <TestName> ./controllers/<pkg>/...`.

Envtests (`test/envtests/`) require `kube-apiserver` and `etcd` binaries and are excluded from `make test`. CI runs them against the minimum supported Kubernetes 1.25 and the build workflow also covers Kubernetes 1.30. Run `make envtest`; the target installs the pinned `setup-envtest` tool and downloads the requested Kubernetes assets.

Helm install from source (see `README.md`): CRDs must be applied first (`kubectl apply -f charts/qubership-monitoring-crds/crds/ --server-side`), then `helm install monitoring-operator charts/qubership-monitoring-operator -n monitoring --create-namespace`.

## Compatibility and upgrade rules

- Preserve existing `PlatformMonitoring` JSON fields, Helm values, defaults, installation parameters, and upgrade behavior.
- Do not silently rename, remove, or change the meaning or default of an existing parameter. For intentional breaking changes, stop and ask for approval; document the impact and migration path.
- Keep the legacy `integreatly.org/v1alpha1` `GrafanaDashboard` CRD while `qubership-grafana-operator-converter` supports migration of existing dashboards; remove it only after the converter migration path is retired and approved.
- When upgrading managed operators, review release notes, APIs, CRDs, RBAC, images, and chart templates together. Run `make generate-all`; do not hand-edit generated files.
- If a linter reports an issue only in generated files, update the linter skip rules instead of editing the generated output.
- Review RBAC when managed resources, owner references, or third-party controllers change.
- Never use `*` for RBAC verbs. Grant only the verbs required by the controller; if generated or imported RBAC introduces a wildcard, fix its source or generation process.
- Validate fresh installations and upgrades with unit tests, envtests, Helm chart checks, and the repository’s Kind integration tests.
- Keep plans and temporary agent artifacts outside the PR unless explicitly requested.

## Architecture

This is a Kubernetes operator (controller-runtime / kubebuilder-style) that reconciles a **single custom resource**, `PlatformMonitoring` (group `monitoring.netcracker.com/v1`), and uses it as a single knob to install and configure an entire monitoring stack. There is one controller, `PlatformMonitoringReconciler` in `controllers/platformmonitoring_controller.go`.

### How reconciliation is structured

`Reconcile` is not a typical per-object loop — on each event it fans out to **many sub-reconcilers**, one per managed component, in a **deliberate order**:

1. `prometheus-operator` first (installs `Prometheus`, `ServiceMonitor`, `PodMonitor`, `Alertmanager`, `PrometheusRule` CRDs consumed by later steps).
2. `etcd`, `kubernetes-monitors` (ServiceMonitors for kube components).
3. VictoriaMetrics stack in order: `vm-operator` → `vmsingle` → `vmcluster` → `vmuser` → `vmagent` → `vmauth` → `vmalertmanager` → `vmalert`. The VM-operator must run first because it installs the `vmetricsv1b1` CRDs the others create.
4. `prometheus` (Prometheus CR), `alertmanager`.
5. Exporters: `kube-state-metrics`, `node-exporter`, `pushgateway`.
6. `grafana-operator` before `grafana` (same CRD-before-CR pattern).
7. `prometheus-rules` last.

Each component has a sub-reconciler under `controllers/`. Failures are recorded in `Status.Conditions`, but reconciliation continues for the remaining components. Failed conditions cause an immediate requeue; otherwise, the default interval is 60 seconds.

The manager is scoped to a **single namespace** via `WATCH_NAMESPACE` env (default `monitoring`), set in `main.go` through `cache.Options.DefaultNamespaces`. Leader election ID is `b0cb59fe.netcracker.com`.

### CR shape

`api/v1/platformmonitoring_types.go` defines the optional component specifications. Components usually use `Install *bool` and `Paused bool`; reconcilers skip disabled or paused components. `FillEmptyWithDefaults()` initializes component defaults before reconciliation.

Modifying the CR surface: after editing `platformmonitoring_types.go`, run `make generate` — this regenerates both `zz_generated.deepcopy.go` AND the CRD YAML under `charts/qubership-monitoring-operator/crds/`.

### Event filtering

`SetupWithManager` installs a predicate (`ignoreDeletionPredicate`) that ignores status-only updates (skips if `metadata.Generation` didn't change) and skips delete events whose state is already known. This matters because the reconciler itself patches status on every run — without the filter it would self-trigger infinitely.

### External schemes

`main.go` registers schemes for managed third-party resources. Reconcilers use discovery to detect optional APIs such as OpenShift Routes and supported Ingress versions before creating resources.

### Helm chart layout

`charts/qubership-monitoring-operator/` contains the deployable chart and managed component subcharts. Their `crds/` directories are canonical. `make update-crds` builds the separate CRD chart, and `make archive-crds` packages the same sources for release.

### Version ldflags

Build version info is injected at link time via `-X github.com/Netcracker/qubership-monitoring-operator/version.{Revision,BuildUser,BuildDate,Branch,Version}` (see `Makefile` `GO_BUILD_LDFLAGS`). The `Dockerfile` does **not** pass these flags — only `make build-binary` does.
