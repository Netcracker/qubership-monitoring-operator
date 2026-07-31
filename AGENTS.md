# Repository agent instructions

## Scope

- This repository builds the Go-based Qubership Monitoring Operator, its Helm charts, CRDs, documentation, and tests.
- These instructions apply repository-wide; keep component-specific guidance next to the affected component.

## Repository map

- `api/v1/` defines the `PlatformMonitoring` API and generated deepcopy code.
- `cmd/operator/` contains the manager entry point; `controllers/` contains the main reconciler and component
  reconcilers.
- `charts/qubership-monitoring-operator/` is the deployable chart; `charts/qubership-monitoring-crds/` packages CRDs
  separately.
- `test/envtests/`, `test/robot-tests/`, and `test/alerts-tests/` cover envtest, deployed integration, and alert-rule
  behavior respectively.
- `docs/` is the published documentation source; `site/mkdocs.yml` configures the documentation build.

## Commands

- Run commands from the repository root unless an instruction says otherwise.
- Regenerate project CRDs and deepcopy methods after API type or marker changes: `make generate`.
- Format, statically check, and run unit tests: `make fmt`, `make vet`, and `make test`.
- Run one Go test: `go test -v -run '<TestName>' ./<package>/...`.
- Build the operator binary: `make build-binary`; this also runs generation, formatting, and vetting.
- Run the operator against the cluster in the active kubeconfig: `make run`; use it only for explicit live-cluster work.
- Build the operator image tagged `qubership-monitoring-operator`: `make image`.
- Refresh copied CRDs after CRD changes: `make docs`; it updates `docs/crds/` and the separate CRD chart.
- Run the complete build pipeline, including image and archives, only when those artifacts are in scope: `make all`.
- After setting `KUBEBUILDER_ASSETS`, run envtests with `ginkgo ./test/envtests/...`; use
  `test/envtests/README.md` only for workstation or container setup.

## Non-obvious invariants

- `controllers/platformmonitoring_controller.go` deliberately runs operator reconcilers before dependent custom
  resources. Preserve CRD-before-CR ordering when adding or moving a component.
- Component errors do not stop later reconcilers. Most add a `Failed` condition, but `kubernetes-monitors` only logs its
  error. Any remaining `Failed` condition causes immediate requeue; otherwise the controller uses
  `RECONCILIATION_INTERVAL`, which defaults to 60 seconds.
- `FillEmptyWithDefaults()` supplies selected images and disables Prometheus when both Prometheus and the
  VictoriaMetrics operator are enabled. Optional component specs remain nullable, so keep the existing nil checks in
  component reconcilers.
- An empty or unset `WATCH_NAMESPACE` makes the manager watch all namespaces; a non-empty value scopes its cache to
  that namespace. Keep this behavior aligned with cross-namespace imports.
- Status-only updates are filtered by `metadata.generation` in `ignoreDeletionPredicate`; preserve that guard when
  changing event handling so status writes do not self-trigger reconciliation.
- The manager registers Prometheus, VictoriaMetrics, Grafana, OpenShift SCC, and extension schemes. Use the discovery
  client for APIs that may be absent from the target cluster instead of assuming every optional API is installed.
- `make build-binary` injects revision, user, date, branch, and version through `GO_BUILD_LDFLAGS`; the operator
  Dockerfile does not use those Makefile flags.

## Done when

- Run `make fmt`, `make vet`, and the smallest relevant Go tests for Go changes; run `make test` when the full unit
  suite is applicable.
- After API or CRD source changes, run the relevant generation command and review all regenerated files.
- For chart changes, follow the preparation steps in `.github/workflows/chart-test-linter.yaml`, then run
  `ct lint --check-version-increment=false --validate-maintainers=false --target-branch=main --chart-dirs=charts`.
- Behavior changes have focused test coverage, and the final response lists checks run and checks that could not run.

## Context routing

- Before changing reconciliation flow, read `controllers/platformmonitoring_controller.go` and `docs/architecture.md`
  to preserve component dependencies and status behavior.
- Before changing installation or chart defaults, read the "Quick Start" section in `README.md` for the supported
  CRD-first installation flow and `docs/configuration.md` for the user-facing configuration contract.
- Before changing alert rules, read `docs/develoment/alert-rules-testing.md`; it defines the repository's rule rendering
  and `vmalert-tool` test conventions.
- Before running envtests, read `test/envtests/README.md` because they require external Kubernetes API server and etcd
  binaries and are excluded from `make test`.
