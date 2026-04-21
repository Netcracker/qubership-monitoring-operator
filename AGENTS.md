# AGENTS.md

This file provides guidance to AI agents when working with code in this repository.

## Common commands

Build / test (from repo root):

- `make generate` — runs `controller-gen` to regenerate CRDs (into `charts/qubership-monitoring-operator/crds/`) and deepcopy methods. Must be re-run after any change to `api/v1/*.go`.
- `make build-binary` — builds `bin/manager` (CGO disabled). Runs `generate`, `fmt`, `vet` first.
- `make test` — unit tests. Runs `go test -race -vet=off --shuffle=on ./... -count=1` across all packages **except** `/test/envtests`.
- `make run` — runs the operator locally against the cluster in `~/.kube/config`.
- `make image` — builds the Docker image (tag `qubership-monitoring-operator`).
- `make docs` — copies CRDs from charts into `docs/crds/`; `make update-crds` then syncs them into `charts/qubership-monitoring-crds/crds/`.
- `make all` — full pipeline: `generate test build-binary image docs archives`.

Running a single Go test: `go test -v -run <TestName> ./controllers/<pkg>/...`.

Envtests (`test/envtests/`) are Ginkgo suites that require real `kube-apiserver` + `etcd` binaries (controller-runtime `envtest`) and are excluded from the `make test` package list. CI runs them via `Netcracker/qubership-core-infra/.github/workflows/generic-go-build.yaml` with `install-envtest: true`, `kube-version: '1.30.0'`, `envtest-version: 'release-0.19'`. Locally use `setup-envtest` or the container workflow described in `test/envtests/README.md`. Run with `ginkgo ./test/envtests/...` once `KUBEBUILDER_ASSETS` is set.

Helm install from source (see `README.md`): CRDs must be applied first (`kubectl apply -f charts/qubership-monitoring-crds/crds/ --server-side`), then `helm install monitoring-operator charts/qubership-monitoring-operator -n monitoring --create-namespace`.

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

Each sub-reconciler lives in its own package under `controllers/` and exposes `NewXxxReconciler(...)` returning a struct that embeds `utils.ComponentReconciler` and implements `Run(cr *monv1.PlatformMonitoring) error` (a handful take `context.Context` as first arg). Failures in sub-reconcilers are **not fatal** — they are logged, the CR's `Status.Conditions` get a `Reason`-keyed entry via `prepareStatusForUpdate`, and the outer reconciler still continues through the remaining components. If any `Failed` condition remains at the end, the reconciler requeues immediately; otherwise it requeues after `RECONCILIATION_INTERVAL` seconds (default `60`, see `controllers/utils/env.go`).

The manager is scoped to a **single namespace** via `WATCH_NAMESPACE` env (default `monitoring`), set in `main.go` through `cache.Options.DefaultNamespaces`. Leader election ID is `b0cb59fe.netcracker.com`.

### CR shape

`api/v1/platformmonitoring_types.go` defines `PlatformMonitoringSpec` as a big flat struct of optional component sub-specs (`AlertManager`, `Prometheus`, `Victoriametrics`, `Grafana`, `NodeExporter`, `KubeStateMetrics`, `KubernetesMonitors`, `GrafanaDashboards`, `PrometheusRules`, `Pushgateway`, `Promxy`, `Integration`, `OAuthProxy`, `Auth`, …). Each component typically has an `Install *bool` toggle and a `Paused bool` flag — sub-reconcilers early-return unless `IsInstall()` is true and `Paused` is false. `FillEmptyWithDefaults()` is called at the top of every reconcile to populate nil sub-structs with defaults, so downstream code can assume non-nil pointers for anything the CR "touches."

Modifying the CR surface: after editing `platformmonitoring_types.go`, run `make generate` — this regenerates both `zz_generated.deepcopy.go` AND the CRD YAML under `charts/qubership-monitoring-operator/crds/`. The `generate` target also injects `helm.sh/hook: crd-install` / `hook-weight: "-5"` annotations into the generated CRDs via `sed`.

### Event filtering

`SetupWithManager` installs a predicate (`ignoreDeletionPredicate`) that ignores status-only updates (skips if `metadata.Generation` didn't change) and skips delete events whose state is already known. This matters because the reconciler itself patches status on every run — without the filter it would self-trigger infinitely.

### External schemes

`main.go` registers a long list of schemes into the manager so the client can read/write third-party CRs directly: `vmetricsv1b1` (VictoriaMetrics operator), `grafv1` (grafana-operator v4 integreatly/v1alpha1), `promv1` (prometheus-operator monitoring/v1), OpenShift `secv1` (SCC), and `apiextensions/v1beta1`. Sub-reconcilers use `DiscoveryClient` to probe which of these APIs are actually present in the cluster (e.g., Ingress v1 vs v1beta1, OpenShift Routes) and branch accordingly — see `controllers/grafana/reconciler.go` for the `HasIngressV1Api()` / `HasIngressV1beta1Api()` pattern.

### Helm chart layout

`charts/qubership-monitoring-operator/` is the deployable chart; it contains the operator templates plus **sub-charts** for each managed component (`victoriametrics-operator`, `prometheus-operator`, `grafana-operator`, `prometheus-adapter-operator`, and the exporters: `blackbox-exporter`, `cert-exporter`, `cloudwatch-exporter`, `json-exporter`, `node-exporter`-equivalents, `promxy`, etc.). The CRDs live in the sub-charts' `crds/` dirs; `make docs` / `make update-crds` / `make archive-crds` copy them into `docs/crds/`, `charts/qubership-monitoring-crds/crds/`, and the release zip respectively. The separate `qubership-monitoring-crds` chart exists so CRDs can be installed independently (and via ArgoCD with `ServerSideApply`) — this is the recommended way to avoid the kubectl last-applied-configuration annotation size limit on large CRDs.

### Version ldflags

Build version info is injected at link time via `-X github.com/Netcracker/qubership-monitoring-operator/version.{Revision,BuildUser,BuildDate,Branch,Version}` (see `Makefile` `GO_BUILD_LDFLAGS`). The `Dockerfile` does **not** pass these flags — only `make build-binary` does.
