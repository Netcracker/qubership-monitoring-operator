# API types

Read the root `AGENTS.md` first. This file covers the traps specific to the CR surface.

## The stored spec is not the effective spec

`FillEmptyWithDefaults()` rewrites the spec in memory on every pass, and the controller persists only the status
subresource. It force-disables a component when a conflicting backend is enabled, so an explicitly requested
component can be off while the stored CR still shows it on. The only signal is a deprecation condition whose message
does not say that the request was overridden.

Read that function before you conclude which components are active from a chart render or a live CR.

## A field in the Go type is not necessarily in the API

`Victoriametrics.VmCluster` and `Victoriametrics.VmReplicas` carry `json:"-"`, with the real tag commented out beside
them. The generated CRD has no such property, but the chart still renders `vmCluster` when
`victoriametrics.vmCluster.install` is `true`. A client-side apply prunes the field, a server-side strict apply
rejects the object, and the decoded spec always holds the zero value.

A Go-constructed CR in a test bypasses the API server's structural-schema pruning, just as a chart render does, so
either can "prove" a code path that no real cluster can reach. Check the generated CRD for the property and the Go
field's tag before you believe a spec field matters.

## The install predicates are not uniform

Some `IsInstall()` helpers return the install flag alone. The VictoriaMetrics ones also require the component's
image field to be non-empty, and the `VmCluster` variant requires all three of its images. The doc comments do not
say so, and a component that fails the predicate is uninstalled, not skipped. Read the specific component's
predicate instead of generalizing from a sibling. See `controllers/AGENTS.md`.

The schema also marks many component fields required even though the defaulting pass would fill them, so the API
server rejects a hand-written minimal CR. For dry-run probes, set them to empty strings; the schema has no
minimum-length constraint on them.

## Generated artifacts

`make generate` runs `controller-gen` with `maxDescLen=0`, which strips every field description to stay under the
object-size limit, and then patches the CRD with path-addressed `yq` expressions. As a result:

- `kubectl explain` returns no description for any field. Read these Go types or the generated
  `docs/api/platform-monitoring.md` instead.
- The `yq` patches address absolute paths. A `del(...)` on a moved path does nothing and exits 0, and an assignment
  to a missing path creates a new node. Nothing verifies that a patch landed.
- The same CRD ships in more than one chart, and the copies have drifted in the past. State which copy you used.

Do not hand-edit the CRD. Regenerate it, and expect the committed file to differ from raw `controller-gen` output by
the post-processing.

## The status contract

The transition timestamp is written with `metav1.Now().String()`, which is Go's `time.Time` string form with the
monotonic-clock reading. The schema declares the field a plain string with no date format, so the API server accepts
it, but no date parser does.
