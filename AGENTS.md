# Repository instructions for coding agents

Use these instructions when changing this repository.

## Common commands

Build / test (from repository root):

- `make generate` — runs `controller-gen` to regenerate project CRDs (into
  `charts/qubership-monitoring-operator/crds/`) and deepcopy methods. Must be re-run after any change to `api/v1/*.go`.
- `make update-operators-crds` — refreshes managed-operator CRDs in their canonical subchart `crds/` directories.
- `make update-crds` — syncs canonical CRDs into `charts/qubership-monitoring-crds/crds/`.
- `make verify-generated` — regenerates all CRDs and fails if generated artifacts drift.
- `make build-binary` — builds `bin/manager` (CGO disabled). Runs `generate`, `fmt`, `vet` first.
- `make test` — unit tests. Runs `go test -race -vet=off --shuffle=on ./... -count=1` across all packages **except**
  `/test/envtests`.
- `make envtest` — runs envtests against Kubernetes 1.25 by default; override `ENVTEST_K8S_VERSION` to test another
  version.
- `make run` — runs the operator locally against the cluster in `~/.kube/config`.
- `make image` — builds the Docker image (tag `qubership-monitoring-operator`).
- `make docs` — compatibility alias for `make update-crds`.
- `make all` — standard build pipeline: `generate test build-binary image docs archives`.

Running a single Go test: `go test -v -run <TestName> ./controllers/<pkg>/...`.

Envtests (`test/envtests/`) require `kube-apiserver` and `etcd` binaries and are excluded from `make test`. CI runs them
against the minimum supported Kubernetes 1.25, and the build workflow also covers Kubernetes 1.30. Run `make envtest`;
the target installs the pinned `setup-envtest` tool and downloads the requested Kubernetes assets.

For a source install, follow the CRD ownership contract in the [quick-start guide](README.md#quick-start):

```bash
kubectl apply --server-side --force-conflicts -f charts/qubership-monitoring-crds/crds/
helm install monitoring-operator charts/qubership-monitoring-operator \
  --skip-crds --namespace monitoring --create-namespace
```

## Compatibility and upgrade rules

- Preserve existing `PlatformMonitoring` JSON fields, Helm values, defaults, installation parameters, and upgrade
  behavior.
- Do not silently rename, remove, or change the meaning or default of an existing parameter. For intentional breaking
  changes, stop and ask for approval; document the impact and migration path.
- Keep the legacy `integreatly.org/v1alpha1` `GrafanaDashboard` CRD while
  `qubership-grafana-operator-converter` supports migration of existing dashboards; remove it only after the converter
  migration path is retired and approved.
- When upgrading managed operators, review release notes, APIs, CRDs, RBAC, images, and chart templates together. Run
  `make generate-all`; do not hand-edit generated files.
- For API changes, run `make generate`, `make update-crds`, and `make verify-generated`.
- Fix generated-file lint findings in their source or generator first. Exclude only a confirmed generated-only false
  positive, narrowly, and with a rationale.
- Review RBAC when managed resources, owner references, or third-party controllers change.
- Never use `*` for RBAC verbs. Grant only the verbs required by the controller; if generated or imported RBAC
  introduces a wildcard, fix its source or generation process.
- With `global.privilegedRights=false`, watches and RBAC must be namespace-scoped, and the operator must not create
  operator-managed `ClusterRole` or `ClusterRoleBinding` resources. The parent `Role` must exactly delegate every API
  group, resource, and verb required by child Grafana and VictoriaMetrics `Role`s.
- Finalize managed custom resources while their CRDs, operators, and RBAC are still available. Remove cluster RBAC only
  afterward. Do not solve ordinary cleanup by manually removing finalizers.
- Run unit tests and Helm checks relevant to the diff. Also run envtests or Kind integration tests when changes affect
  API-server behavior, RBAC, lifecycle, installation, or upgrades.
- Keep plans and temporary agent artifacts outside the PR unless explicitly requested.

## Architecture

This is a Kubernetes operator (controller-runtime / kubebuilder-style) that reconciles a **single custom resource**,
`PlatformMonitoring` (group `monitoring.netcracker.com/v1`), and uses it as a single knob to install and configure an
entire monitoring stack. There is one controller, `PlatformMonitoringReconciler`, in
`controllers/platformmonitoring_controller.go`.

### How reconciliation is structured

`Reconcile` is not a typical per-object loop. On each event it runs sub-reconcilers in this order; the controller code
is the authority:

1. `prometheus-operator`, `kubernetes-monitors`, `vm-operator`, `vmsingle`, `vmcluster`, `vmuser`, `vmagent`, and
   `vmauth`.
2. `prometheus`, `vmalertmanager`, `alertmanager`, `vmalert`, `kube-state-metrics`, and `node-exporter`.
3. `grafana-operator`, `grafana`, `prometheus-rules`, and `pushgateway`.

The Prometheus, VictoriaMetrics, and Grafana operator reconcilers precede the reconcilers for their managed custom
resources so operator workloads, RBAC, and configuration are reconciled first. This ordering does not install CRDs.
All required CRDs must already be installed through Helm or the dedicated CRD chart before operator-chart or controller
reconciliation begins.

Each component has a sub-reconciler under `controllers/`. Reconciliation continues after component failures. Every
component except `kubernetes-monitors` records failures in `Status.Conditions`; `kubernetes-monitors` failures are
logged only. A remaining failed condition causes an immediate requeue; otherwise, the default interval is 60 seconds.

The chart sets `WATCH_NAMESPACE` to the release namespace. If the binary receives an unset or empty
`WATCH_NAMESPACE`, it leaves `cache.Options.DefaultNamespaces` unset and watches all namespaces. Leader election ID is
`b0cb59fe.netcracker.com`.

### CR shape

`api/v1/platformmonitoring_types.go` defines the optional component specifications. Components usually use
`Install *bool` and `Paused bool`; reconcilers skip disabled or paused components. `FillEmptyWithDefaults()`
initializes component defaults before reconciliation.

Modifying the CR surface: after editing `platformmonitoring_types.go`, run `make generate`, `make update-crds`, and
`make verify-generated`. These regenerate deepcopy methods and the CRD artifacts, then verify they do not drift.

### Event filtering

`SetupWithManager` installs a predicate (`ignoreDeletionPredicate`) that ignores status-only updates (skips if
`metadata.Generation` did not change) and unknown-final-state delete tombstones, while accepting confirmed delete
events. This matters because the reconciler itself patches status on every run; without the filter it would
self-trigger indefinitely.

### External schemes

`main.go` registers schemes for managed third-party resources. Reconcilers use discovery to detect optional APIs such
as OpenShift Routes and supported Ingress versions before creating resources.

### Helm chart layout

`charts/qubership-monitoring-operator/` contains the deployable chart and managed component subcharts. Their `crds/`
directories are canonical. `make update-crds` builds the separate CRD chart, and `make archive-crds` packages the
same sources for release.

### Version ldflags

Build version info is injected at link time via
`-X github.com/Netcracker/qubership-monitoring-operator/version.{Revision,BuildUser,BuildDate,Branch,Version}` (see
`Makefile` `GO_BUILD_LDFLAGS`). The `Dockerfile` does **not** pass these flags; only `make build-binary` does.
