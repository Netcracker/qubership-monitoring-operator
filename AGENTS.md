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

Every file in `test/envtests/` carries the `envtest` build tag, so a plain `go test ./...` compiles none of the suite
and still exits 0. In CI, only the `Envtest with Kubernetes 1.25` job runs it. The coverage job provisions Kubernetes
1.30 envtest assets, but it runs `go test ./...` without the tag. See `test/AGENTS.md`.

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
`Dockerfile` does not pass them. The injection is a no-op: the tree has no `version` package, and the linker accepts
`-X` for an absent symbol without a warning. Check version stamping on a built binary, not by reading these flags.

## Traps this repository sets

These notes come from reviews in which many agents worked on this repository independently. Each item is something
more than one of them got wrong or spent significant time on. They describe observed behavior, not policy; if one
stops being true, correct it.

Scoped notes live next to the code they describe: `api/v1/`, `controllers/`, `charts/`, and `test/`.

### A green test run does not prove a behavior is covered

The unit tests call the top-level `Reconcile` for a few status and migration paths only, and the envtest suite is
behind a build tag. Before you claim that a test covers a reconciler behavior, change the line you care about and
confirm that a test fails. Otherwise, prove the claim with a probe you write.

### Reading the CR status

The transition timestamp comes from `metav1.Now().String()`, which is Go's `time.Time` string form with the
monotonic-clock reading. The schema declares the field a plain string, so the API server accepts it, but no date
parser does. Do not feed it to `jq` date functions or `date`.

### Working on the live cluster

Treat every CR patch as destructive. Setting a component's install flag to `false` does not idle the component: the
reconciler takes the uninstall branch and deletes the running workload within a cycle. Restore the value
immediately.

Assume other writers. Several agents may share one cluster, so a deployment's environment, a bound port, or a stray
manager process may belong to someone else. Re-read the deployment environment before you trust a rate or counter
measurement, and capture the object before and after your own change instead of inferring causation.

`--help` is safe: the flag parser prints usage and exits before the kubeconfig is loaded. `--version` is not. The
imported flag does not exit, so the process loads the current kubeconfig and starts a second operator against the
cluster you point at. Most of the operator's flags come from an imported metrics library and control nothing here.

An incomplete RBAC grant makes the operator hang instead of fail: the cache waits for a watch that never syncs, the
reconciler logs stay silent, and the CR status reports nothing. Look for the cache's own log lines.

### Reading command output

Do not take an exit status from the end of a pipeline, and do not discard stderr from a chart render.
`helm template … 2>/dev/null | grep -c` cannot tell "the resource is absent" from "nothing rendered". Several agents
drew conclusions from truncated output, including one who read `list … | head -20` of a descending list as the
complete set.

Commands that reach the network, such as the CRD refresh targets, the envtest asset download, and image pulls, are
slow and flaky in proportion to the link. A single slow run is not evidence of a build defect, and a hang on first
use is usually the registry. Retry before you report it.

### Tooling lives in the ignored `bin/` directory

The Makefile downloads pinned binaries into `bin/` as `<tool>-<version>`. A tool missing from `PATH` may still be
available there, so list that directory before you conclude that a check cannot run. In a fresh clone or worktree
the directory is empty, and the verification scripts under `tools/` stop with a message that names the `make` target
to run. Run that target instead of copying binaries between trees.

Some scripts under `tools/` use Bash 4 built-ins such as `mapfile` under `set -e`. The default macOS Bash 3.2 either
stops partway or iterates an empty list and exits 0, so a partial run looks like a clean pass. Run them with Bash 4
or later and confirm that they printed results.

Workflows call these scripts directly; the Makefile does not. Search the workflows as well before you claim that
something is unverified.

### What CI does not check

`.github/super-linter.env` disables Go linting and Kubernetes schema validation, and no workflow runs
`golangci-lint`. The super-linter workflow fetches shared linter rules from the `main` branch of the
`netcracker/.github` repository, so the rules can change without a commit here. The coverage job runs
`go test ./...` from a shared reusable workflow, not `make test`, so it runs without the race detector, shuffling, or
the envtest build tag.

A lint config in the tree does not prove that the linter runs, and a green pipeline does not prove that the
documented test flags ran.
