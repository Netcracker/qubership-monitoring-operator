# Charts

Read the root `AGENTS.md` first. This file covers the traps specific to the chart tree.

## The parent chart is a fraction of the chart

The largest components are vendored subchart directories, each with its own values file, templates, and lifecycle
hooks, merged under the parent key of the same name. The parent `values.yaml` does not mention the biggest
components. Agents have concluded that a knob was invented, that a dependency condition was dead, and that a
pre-delete hook did not exist, all from searching the parent chart alone.

Search the whole chart tree. A dependency `condition:` resolves against the coalesced parent and subchart values, so
a key missing from the parent `values.yaml` can still be a working knob. Helm enables a subchart by an exact key-path
match against its `condition:`. A block in the parent values with different casing or punctuation is a dead knob
that renders without error.

## A render is a conditional projection, not an inventory

Do not conclude that a resource does not exist from a render. Four mechanisms shrink the output:

- **Cluster capabilities.** Many templates branch on `.Capabilities.APIVersions.Has`, so OpenShift security context
  constraints and similar platform objects disappear unless you pass the API version with `--api-versions`. Audits of
  platform-specific behavior must set it.
- **`lookup`.** `grafana-admin-credentials-secret.yaml` branches on a cluster read that returns nil offline, so a
  client-side render always takes the "does not exist yet" branch. A real install may differ in either direction.
  Confirm ownership on a cluster through the release annotations and field managers.
- **Sibling guards.** Some CR blocks are nested inside another component's guard, so setting them alone renders
  nothing.
- **Defaults already on.** A zero diff after you enable a toggle usually means the default was already `true`.

A local render evaluates no cluster state, so it cannot reveal a conflict over a cluster-scoped name that another
release already owns. Use a server-side dry run for that.

## Values are not validated

`values.schema.json` never sets `additionalProperties: false`, and several top-level keys are missing from it.
Misspelled keys render cleanly, lint cleanly, and do nothing, including a misspelling meant to disable a component,
which leaves it on. Check an override by diffing the render against the unset baseline, not by an exit code.

The files under `resource-profiles/` show the failure mode. Each one is a single line in dotted form, such as
`global.profile: large`. Helm expands dotted keys only for `--set`, never for `-f`, so `-f resource-profiles/large.yaml`
renders byte-identical to the default. Use `--set global.profile=large`.

The nested toggle that gates a service monitor is `serviceMonitor.enabled` in some subcharts and
`serviceMonitor.install` in others. Search the
specific subchart's template for the condition before you set anything nested.

Numeric values pass through Helm's `default`, which treats zero as unset, so an explicit scale to zero silently
becomes one replica.

## What the chart renders is not what the operator receives

The chart renders CR fields that the CRD does not declare, because a Go struct tag excludes them from
serialization. See `api/v1/AGENTS.md`. A server-side strict apply rejects the whole object and fails the release,
while a client-side apply prunes the field silently. Either way, the operator sees the zero value.

The defaulting pass also rewrites component enable flags in memory before reconciliation, and only the status
subresource is persisted. The stored CR shows what the user asked for, not what the operator does.

Before you believe a values key works, check the CRD for the property and the Go field's tag. A field in the
`helm template` output proves nothing.

## Installing, upgrading, and removing

`--wait` covers only the chart's own objects. The operator creates most of the stack afterward, so a successful
upgrade or rollback exit code says nothing about convergence. Wait on the custom resource's status conditions and
the managed workloads instead.

The chart has four lifecycle hooks. Each hook `Job` uses the `before-hook-creation,hook-failed` delete policy, so a
reinstall replaces it and Helm deletes a failed one. The supporting `ServiceAccount`, `Role`, and binding objects use
`before-hook-creation,hook-succeeded`, so they remain only after a failure, until the next hook run replaces them. You
do not need to delete hook resources by hand before you reinstall. Because Helm deletes a failed `Job`, read its pod
logs while the hook runs.

- **`monitoring-grafana-cleanup-hook`** (`pre-delete`, weight `-10`) scales the operator to zero, then deletes Grafana
  dashboards, data sources, and `Grafana` objects in the release namespace that carry
  `app.kubernetes.io/managed-by-operator=monitoring-operator`, plus two dashboards by name.
- **The VictoriaMetrics cleanup hook** (`pre-delete`, weight `-5`) scales the operator to zero, then deletes every
  `VMSingle` (orphaning its dependents), `VMUser`, `VMAuth`, `VMAlert`, `VMAgent`, `VMAlertmanager`, and `VMCluster`
  in the release namespace with `--all`.
- **`monitoring-rbac-cleanup-hook`** (`post-delete`) deletes `ClusterRole` and `ClusterRoleBinding` objects that carry
  the managed-by-operator label and `monitoring.netcracker.com/installation-namespace=<release namespace>`, plus this
  installation's node-exporter SCC.
- **The `etcdCertsJob` hook** (`post-install`, `post-upgrade`) copies etcd client certificates into the
  `kube-etcd-client-certs` secret. It renders only when `global.privilegedRights` is `true` and `publicCloudName` is
  empty.

A finished cleanup `Job` is removed after `victoriametrics.cleanup.hook.ttlSecondsAfterFinished` (default 180
seconds), and a finished `etcdCertsJob` after `etcdCertsJob.ttlSecondsAfterFinished` (default 180 seconds).

The VictoriaMetrics cleanup hook uses `--all`, not a label selector, so an uninstall also deletes VictoriaMetrics
custom resources that someone created by hand in the release namespace. Do not uninstall in a namespace that holds
VictoriaMetrics resources you want to keep.

Do not skip these hooks with `--no-hooks` as a shortcut. The pre-delete hooks finalize managed custom resources while
their operators and RBAC still exist, and the post-delete hook removes cluster RBAC afterward. Skipping them leaves
that work undone. If a hook fails, read its logs and fix the cause.

A rollback can also fail on immutable fields that changed between versions, and it can report success while the
stack never converges.

## CRDs

The same CRD ships in more than one chart, and the copies have drifted in the past, so state which copy a schema
comparison used. CRD files are named for the fully qualified CRD, not the component, and a subchart's `crds/`
directory can hold CRDs for other subcharts.

The shipped CRDs are generator output post-processed by path-addressed `yq` expressions. See `api/v1/AGENTS.md`.

## Chart tooling

The chart-lint step packages every vendored subchart into the source tree, and `.gitignore` covers neither the
archives nor the lockfile. The step leaves the working tree dirty beside directories that are already tracked, so
clean up before you read a diff.

The same workflow copies example profiles with `docs/examples/deploy-parameters/**/*.yaml`, which does not recurse in
CI's shell, and `ct lint` then selects files by their `-values.yaml` suffix. Many documented profiles never reach the
linter. The step runs with no install phase and with `--check-version-increment=false`. A green chart job is not
evidence that a profile, a rendered manifest, or a version bump was checked.

Compare your local Helm version with the version CI installs before you reason from a local result. Helm major
versions differ in flags, in `lookup`, in hook handling, and in whether an unknown field is pruned or fails the
release.
