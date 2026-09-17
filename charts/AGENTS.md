# The charts

Read the root `AGENTS.md` first. This file covers the traps specific to the chart tree.

## The parent chart is a fraction of the chart

The largest components are vendored subchart directories, each with its own values file, templates, and lifecycle
hooks, merged under the parent key of the same name. The parent values file does not even mention the biggest
components. Agents concluded three separate times that a knob was invented, that a dependency condition was dead,
and that a pre-delete hook did not exist — all from grepping the parent alone.

Search the whole chart tree. A dependency `condition:` resolves against coalesced parent-plus-subchart values, so
absence from the parent values file never means "not a knob".

Two mechanics worth knowing before you reason about a toggle:

- A `condition:` holding several comma-separated paths is **not** an OR. Helm stops at the first path that exists in
  the coalesced values, so a second path is unreachable while the first ships with a default.
- A subchart is enabled by an exact key-path match against its `condition:`. A differently-cased or differently
  punctuated block in the parent values is a dead knob that renders without error.

## A render is a conditional projection, not an inventory

Never conclude "this resource does not exist" from a render. Four independent mechanisms shrink the output:

- **Cluster capabilities.** Many templates and several subchart helpers branch on `Capabilities.APIVersions.Has`,
  so every OpenShift security-constraint object and every pod-security-policy object disappears unless you pass the
  capability explicitly. Audits of platform-specific behavior must set it.
- **`lookup`.** Credential templates branch on a cluster read that returns nil offline, so a client-side render
  always takes the "does not exist yet" branch. An object present in the render may never be created by a real
  install, and vice versa. Confirm ownership on a cluster through the release annotations and field managers.
- **Sibling guards.** At least one CR block is nested inside an unrelated component's guard, so setting it alone
  renders nothing.
- **Defaults already on.** A zero diff after enabling a toggle usually means the default was already true.

Local rendering also evaluates no cluster state, so it cannot reveal a conflict over a cluster-scoped name already
owned by another release. Use a server-side dry run for that.

## Values are not validated, in either direction

No values schema sets `additionalProperties: false`, and several top-level keys are absent from the schema
entirely. Misspelled keys render clean, lint clean, and do nothing — including a misspelling in the *disable*
direction, which leaves the component on. Verify an override by diffing the render against the unset baseline,
never by an exit code.

Two shipped artifacts demonstrate the failure mode: the resource profiles are single-line files written in dotted
key form, which Helm expands only for `--set` and never for `-f`, so all of them render byte-identical to the
default; and at least one example overlay sets a key spelled differently from the condition that reads it.

The same concept is also spelled differently across sibling subcharts — the nested toggle that gates a monitor
object is `.enabled` in about half of them and `.install` in the rest. Grep the specific subchart's template for
the condition before setting anything nested.

Numeric values pass through Helm's `default`, which treats zero as unset, so an explicit scale-to-zero silently
becomes one replica.

## What the chart renders is not what the operator receives

The chart renders CR fields the CRD does not declare, because a Go struct tag excludes them from serialization
while the template still emits them. What happens then depends on the apply path: a server-side strict apply
rejects the whole object and fails the release, while a client-side apply prunes the field in silence. Either way
the operator sees the zero value.

Beyond that, a defaulting pass rewrites component enable flags in memory before reconciliation and only the status
subresource is persisted — so the stored CR shows what the user asked for, not what the operator does.

Before believing a values key works, check the CRD for the property and the Go field's tag. `helm template` showing
a field proves nothing.

## Installing, upgrading, and removing

`--wait` covers only the chart's own objects. Nearly everything that matters is created afterwards by the
operator, so a successful upgrade or rollback exit code says nothing about convergence — gate on the custom
resource's status conditions and the managed workloads instead.

Lifecycle hooks are the main source of grief:

- The post-install and post-upgrade hook assumes a control-plane topology a development cluster does not have. It
  can stall for the whole timeout and leave the release `failed` while every workload is fine. Expect to pass
  `--no-hooks` on upgrades as well as installs.
- Pre-delete hooks act on the whole namespace rather than on release-owned objects, so an uninstall destroys
  matching custom resources somebody else created by hand. Never uninstall in a namespace holding anything you
  care about, and read the hook pod logs — Helm's own "kept due to resource policy" summary does not reflect what
  the hooks did.
- Hook resources are retained across uninstall and must be deleted by hand before reinstalling; a failed hook can
  leave the release neither installed nor removable except with `--no-hooks`.

Rollback additionally trips over immutable fields changed between versions, and reports success while the stack
never converges.

Several cluster-scoped names are hardcoded without a release or namespace prefix, so a second release in the same
cluster collides — and the managed operators watch all namespaces by default, so the first release's operators
reconcile and finalize the second one's resources. Side-by-side experiments are contaminated unless you isolate
them deliberately.

## CRDs

The same CRD is shipped in more than one chart and the copies have drifted by weeks in the past, so any schema
comparison must state which copy it used. Files are named for the fully-qualified CRD, not the component, and a
subchart's CRD directory is not restricted to that subchart.

The shipped CRDs are generator output post-processed by path-addressed YAML patches. Every field description is
stripped to stay under the object-size limit, so cluster-side introspection teaches nothing about field meaning —
read the Go types. The patches address absolute paths and no-op silently when an upstream path moves, and a
deletion against a missing path exits zero, so the guard passes vacuously.

## Chart tooling

The chart-lint step packages every vendored subchart into the source tree and neither the archives nor the lock
file are ignored, so it leaves the working tree dirty beside directories that are already tracked. Clean up before
reading a diff.

That step also copies the example profiles with a recursive glob that does not recurse in CI's shell and then
selects them by filename suffix, so well under half of the documented profiles ever reach the linter. It runs
without an install phase and with the chart-version gate switched off. A green chart job is not evidence that a
profile, a rendered manifest, or a version bump was checked.

The locally installed Helm is likely a major version ahead of the one CI installs, and the majors differ on flags,
on `lookup`, on hook handling, and on whether an unknown field is pruned or fails the release. Check the version
before reasoning from documentation or from a local result.
