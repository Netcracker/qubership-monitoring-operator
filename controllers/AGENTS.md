# Controllers

Read the root `AGENTS.md` first. This file covers the traps specific to the reconcilers.

## Enablement is not the boolean it looks like

The `IsInstall()` helpers are not uniform. Some check only the install flag. The VictoriaMetrics ones also require
the component's image field to be non-empty, and the `VmCluster` variant requires all three of its images. The doc
comments do not say so, and several agents built probes that set only the flag and silently exercised the disabled
path. Read the specific component's predicate instead of generalizing from a sibling.

The reconcilers branch on the same predicate, and a component that fails it is uninstalled, not skipped. A
hand-built CR without images deletes live resources. Set images in every fixture, and do not read `install: true` as
"this component is on".

## The embedded assets are templates, not manifests

The Go handlers load each file under a component's `assets/` directory and change it before they apply it, so
reading an asset tells you little:

- Several `cluster-role.yaml` assets ship with `rules: []`; the handlers add the rules at reconcile time. "No
  permissions" is the wrong conclusion.
- `utils.MustAssetReader` and `utils.MustAssetReaderToString` do not panic, despite their names and doc comments.
  They discard the read error and return empty content, so a missing asset surfaces later as an opaque decode error.
- The handlers can overwrite values from an asset with values from the CR, including nil. A default install can end
  up without a field that the asset sets.

Read the matching `manifest.go` or `handlers.go`, or exercise the builder from a test, before you state what the
operator applies. Then confirm against the deployed object, not the asset.

## Writing probes

Two harnesses are available, and both have edges:

- **The controller-runtime fake client alone panics.** The reconcilers hold a discovery client and probe for
  optional APIs before they create anything. Inject a `fakediscovery.FakeDiscovery` with the relevant API groups, as
  the existing tests do.
- **The fake client also invents failures.** It deep-copies where the real client path does not, which can panic on
  unstructured objects that a real cluster handles. Confirm a fake-client result in either direction against the
  envtest suite.

Trace callers to `Reconcile` before you treat a function as live. Some helpers are reachable only from other unused
code, while the RBAC grants and status semantics they motivated still ship.

## Reconcile timing and cancellation

A single pass blocks for tens of seconds by design: the Grafana handler sleeps for 30 seconds unconditionally. The
requeue interval starts at the end of the pass, so the effective period is longer than `RECONCILIATION_INTERVAL`.
Most sub-reconcilers receive no context and call `context.TODO()`, so a pass largely ignores cancellation, and pod
termination consumes most of the grace period.

Budget for this when you poll for an effect, and do not read a slow status update as a hang. `Reconcile` also names
its context parameter `context`, which shadows the package, so grep-based reasoning about which calls receive a
cancellable context is unreliable in that file.

## RBAC is hand-written, not generated

Unlike a stock kubebuilder project, the operator's `Role` and `ClusterRole` are hand-maintained chart templates. The
tree has only a few `+kubebuilder:rbac` markers and no generated RBAC directory, so adding an API call to a reconciler
does not update permissions, and no generator or drift check catches the omission. On a cluster, the omission
appears as a watch that never syncs: the operator hangs instead of failing, and the reconciler logs say nothing.

When you add a call, update the chart templates by hand. `tools/verify-helm-crds.sh` checks that the parent role can
delegate the permissions the child roles need.
