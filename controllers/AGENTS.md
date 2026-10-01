# The controllers

Read the root `AGENTS.md` first. This file covers the traps specific to the reconcilers.

## Enablement is not the boolean it looks like

The `IsInstall()` helpers are not uniform. Some check only the install flag; the VictoriaMetrics ones are
conjunctions that also require the component's image fields to be non-empty, and the cluster variant requires all
three of its images. The doc comments do not say so, and several agents built probes that set only the flag and
silently exercised the disabled path. Read the specific component's predicate rather than generalizing from a
sibling.

Worse, the reconcilers branch on the same predicate and a component that fails it is **actively uninstalled**, not
skipped — a hand-built CR without images deletes live resources. Populate images in every fixture, and never read
`install: true` as "this component is on".

The same applies on a cluster: flipping a component's install flag to false in a live CR deletes its workload
within a cycle.

## The embedded assets are templates, not manifests

Files under each component's `assets/` directory are loaded and then mutated by the Go handlers before they are
applied, so reading an asset tells you little:

- Several cluster-role assets ship with an empty rule list; the rules are appended at reconcile time. "No
  permissions" is the wrong conclusion.
- At least one asset is a zero-length file that is still pulled in by the embed directive, and at least one
  constant points at an asset that does not exist. A `Must…` accessor here does not panic — it discards the error
  and returns empty, so a missing asset surfaces far away as an opaque decode error.
- Values visible in an asset can be overwritten unconditionally from the CR, including with nil. A default install
  can therefore end up without a field the asset clearly sets.

Read the matching `manifest.go` / `handlers.go`, or exercise the builder from a test, before asserting what the
operator applies. And confirm against the deployed object rather than the asset.

## Writing probes

Nothing exercises reconciliation, so any behavioral claim needs a probe you write. Two harnesses, both with edges:

- **The controller-runtime fake client alone panics.** The reconcilers hold a discovery client and probe for
  optional APIs before creating anything; the fake provides none. Inject a fake clientset with the relevant API
  groups populated.
- **The fake client also invents failures.** It deep-copies where the real client path does not, which raises
  panics on unstructured objects that a real cluster never triggers. Confirm both directions against the
  API-server-backed suite before believing a fake-client result.

Trace callers to a top-level entry point before treating a function as live. Several helpers here are retained for
manual recovery and are reachable only from other unused code — while the RBAC grants and status semantics they
motivated are still shipped.

## Reconcile timing and cancellation

A single pass blocks for tens of seconds by design — one component handler sleeps unconditionally — and the requeue
is measured from the end of the pass, so the effective period is substantially longer than the configured interval.
The context is not propagated into the sub-reconcilers, so a pass ignores cancellation and pod termination consumes
most of the grace period.

Budget for this when polling for an effect, and do not read a slow status update as a hang. Note also that the
reconcile method names its context parameter after the `context` package, so grep-based reasoning about which calls
receive a cancellable context is unreliable in that file.

## RBAC is hand-written, not generated

Unlike a stock kubebuilder project, the operator's `Role` and `ClusterRole` are hand-maintained chart templates.
The tree holds only a handful of RBAC markers and there is no generated RBAC directory, so adding an API call to a
reconciler does not update permissions and no generator or drift check will catch the omission. It surfaces on a
cluster as a watch that never syncs — the operator hangs rather than failing, and the reconciler logs say nothing.

Whenever you add a call, update the chart templates by hand, and check the delegation script under `tools/` that
verifies the parent role covers what the child roles need.
