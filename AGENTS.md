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

Envtests (`test/envtests/`) require `kube-apiserver` and `etcd` binaries and are excluded from `make test`. Run
`make envtest`; the target installs the pinned `setup-envtest` tool and downloads the requested Kubernetes assets.

Only the dedicated minimum-version job runs them in CI. The coverage job provisions envtest assets but invokes a
plain package-wildcard test run, and every file in that directory carries a build tag — so it compiles none of the
suite and still reports success. See `test/AGENTS.md`.

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

`Makefile` `GO_BUILD_LDFLAGS` passes `-X <module>/version.{Revision,BuildUser,BuildDate,Branch,Version}`, and the
`Dockerfile` does not pass them at all. Note that **the injection is a no-op**: no `version` package exists in the
tree, and the linker accepts `-X` for an absent symbol without a warning, so the build succeeds and the binary
carries nothing. Verify version stamping by inspecting a built binary, never by reading these flags.

## Traps this repository sets

Distilled from reviews in which several dozen agents worked this repository independently. Everything here is
something more than one of them got wrong or lost significant time to. These are observations, not policy — if one
stops being true, correct it.

Scoped notes live next to the code they describe: `api/v1/`, `controllers/`, `charts/`, `test/`.

### A green test run is compatible with completely broken reconciliation

No test calls the top-level `Reconcile`, and the default unit target excludes the envtest suite. Agents inverted the
Prometheus-versus-VictoriaMetrics defaulting rule and removed `Paused` guards from several sub-reconcilers; the full
suite still passed. The race detector finds nothing because the non-test code has no goroutines, channels, or
mutexes at all.

Prove any reconciler claim with a probe you write. See `test/AGENTS.md` before running the envtest suite — it is
gated behind a build tag, so the obvious command runs none of it and still exits 0.

### The CR status is a transcript of the current pass, not a state

A healthy stack reports the aggregate condition as `In progress` / `False` for roughly a third to a half of every
cycle, because the controller writes that at the start of each pass and the terminal state at the end. A single
`kubectl get platformmonitoring -o jsonpath='{.status.conditions}'` is therefore a coin flip.

Sample across at least two full intervals and report the fraction. The transition timestamp is Go's `time.Time`
string form including the monotonic-clock reading, and the schema declares it a plain string — no date parser will
accept it, so never feed it to `jq` or `date`.

### Working on the live cluster

Treat every CR patch as destructive: setting a component's install flag to false does not idle the component, it
takes the uninstall branch and deletes the running workload within a reconcile cycle. Restore immediately.

Assume other writers. Reviews here run several agents against one cluster, so a deployment's environment, a bound
port, or a stray manager process may be someone else's. Re-read the deployment env before trusting any rate or
counter measurement, and capture the object before and after your own change rather than inferring causation.

Do not run the operator binary to inspect its flags. There is no early exit for informational flags — the process
falls through to loading the current kubeconfig and starts reconciling, so probing `--help` or `--version` starts a
second operator against whatever cluster you are pointed at. Most of its ~100 flags come from an imported metrics
library and control nothing here.

An incomplete RBAC grant makes the operator hang rather than fail: the cache waits on a watch that never syncs,
reconciler logs stay silent and the CR status says nothing. Look for the cache component's own log lines.

### Reading command output

Never take an exit status from the end of a pipeline, and never discard stderr from a chart render. This chart
genuinely fails to render under ordinary value combinations, so `helm template … 2>/dev/null | grep -c` cannot
distinguish "the resource is absent" from "nothing rendered at all". Several agents drew conclusions from a
truncated or piped result — including one who read `list … | head -20` of a descending list as the complete set.

Commands here that reach the network — the CRD refresh targets, the envtest asset download, image pulls — are slow
and flaky in proportion to the link, not to the repository. A single timing is not evidence of a build defect, and
a hang on first use is usually the registry. Retry before reporting.

### Tooling lives in a gitignored directory under version-suffixed names

The Makefile downloads pinned binaries into the local `bin/` directory as `<tool>-<version>`, so a tool absent from
`PATH` is not absent from the project — list or glob that directory before concluding a check cannot be run. In a
fresh clone or worktree the directory is empty and the repository's own verification scripts abort; run the make
target that provisions them rather than copying binaries between trees.

Several verification scripts under `tools/` use bash-4 builtins under `set -e`. On a default macOS shell they either
abort partway or, worse, iterate an empty input list and exit 0 — a partial run and a clean pass look identical. Run
them with a modern bash and confirm they printed results.

Much of the project's real verification lives in those scripts and is invoked only from workflows, never from the
Makefile. Grep the workflows as well before asserting that something is unverified.

### What CI does not check

The linter env file disables Go linting and Kubernetes schema validation, no workflow invokes the Go linter, and the
linter rules are fetched at run time from an organization repository on a moving branch. Merges are gated on no
required status check at all. The workflow that produces coverage does not use the documented test target and runs
without the race detector or shuffling.

Presence of a lint config in the tree is not evidence the linter has ever run; a green pipeline is not evidence the
documented flags ran.
