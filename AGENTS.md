# Repository agent instructions

## Scope

- This repository builds the Kubernetes operator and Helm charts that reconcile a `PlatformMonitoring` resource into
  a monitoring stack.
- These instructions apply repository-wide. Keep component-specific details near the affected controller, chart, or
  test suite.

## Repository map

- `api/v1/` defines the custom-resource types and generation markers.
- `controllers/platformmonitoring_controller.go` sequences component reconcilers from `controllers/`.
- `charts/qubership-monitoring-operator/` is the deployable chart, while `charts/qubership-monitoring-crds/` packages
  CRDs for separate installation.
- `test/envtests/` contains controller-runtime envtests; `test/robot-tests/` contains cluster-level integration tests.
- `agent-packages/troubleshoot-monitoring/` is the source of the troubleshooting skill and the
  `docs/troubleshooting.md` symlink target.

## Commands

- Run every command in this section from the repository root. Download Go modules with `go mod download`.
- Run a focused controller package with
  `go test -race -vet=off --shuffle=on ./controllers/prometheus/... -count=1 -v`; replace `prometheus` with the
  affected package.
- Run the core local Go gate with `make fmt vet test`.
- Build the operator binary with `make build-binary`; this also runs generation, formatting, and vetting.
- After setting `KUBEBUILDER_ASSETS`, run envtests with `ginkgo ./test/envtests/...`.

## Non-obvious invariants

- Do not edit `api/v1/zz_generated.deepcopy.go` or the operator's generated `PlatformMonitoring` CRD directly. Change
  the API types or markers, run `make generate`, then inspect the generated diff.
- After changing CRD sources, run `make docs` to synchronize chart CRDs into `docs/crds/` and the dedicated CRD chart;
  review the diff for unrelated generated changes.
- Preserve the component order in `PlatformMonitoringReconciler`: operators must create required CRDs before dependent
  resources. Component failures must not stop later reconcilers; preserve status conditions where the component owns
  one. Add targeted regression coverage when changing this flow, then run
  `go test -race -vet=off --shuffle=on ./controllers/... -count=1 -v`.
- Keep `ignoreDeletionPredicate` filtering updates with an unchanged generation so status writes do not retrigger
  reconciliation. When changing controller setup, cover the predicate with a focused test and run the controller suite.
- `make test` intentionally excludes `test/envtests/`; do not report it as envtest coverage.

## Done when

- Focused tests for the affected package and the applicable `make fmt vet test` checks pass.
- Generated API and CRD outputs are refreshed when their sources change, with no unrelated generated drift.
- Chart, alert-rule, envtest, or cluster-integration changes pass their applicable checks, or the final report names the
  checks that could not run and why.
- The final report separates repository-source inspection from commands actually executed.

## Context routing

- Before changing envtests, read `test/envtests/README.md` for required assets and
  `.github/workflows/test-sonar-go-coverage.yaml` for CI's pinned envtest setup.
- Before changing Helm charts, read `.github/workflows/chart-test-linter.yaml` for CI's staging and `ct lint` command.
- Before changing alert rules, read `.github/workflows/test-alert-rules-unit-tests.yaml` for rule preprocessing and the
  `vmalert-tool` invocation.
- Before changing cluster-level test coverage, read `.github/workflows/integration-tests.yaml` and
  `test/robot-tests/README.md` for the external pipeline and supported Robot tags.
