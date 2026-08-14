# The API types

Read the root `AGENTS.md` first. This file covers the traps specific to the CR surface.

## The stored spec is not the effective spec

`FillEmptyWithDefaults()` rewrites the spec in memory on every pass, and only the status subresource is persisted.
It force-disables a component when a conflicting backend is enabled, so an explicitly requested component can be
off while the stored CR still shows it on, and the only signal is a deprecation condition whose message does not
say the request was overridden.

Read that function before concluding anything about which components are active from a chart render or from a live
CR.

## A field in the Go type is not necessarily in the API

At least one component field carries a struct tag that excludes it from serialization, with the real tag commented
out beside it, so the generated CRD has no such property while the chart still renders one. The field is pruned by
a client-side apply and rejected outright by a server-side strict apply, and the decoded spec always holds the zero
value.

A Go-constructed CR in a test bypasses the API server's structural-schema pruning exactly as a chart render does,
so either can "prove" a code path that no real cluster can reach. Check the generated CRD for the property and the
Go field's tag before believing a spec field matters.

## The install predicates are not uniform

Some `IsInstall()` helpers return the install flag alone; the VictoriaMetrics ones also require the component's
image fields to be non-empty, and the cluster variant requires all three of its images. Nothing in the doc comments
says so, and a component that fails the predicate is uninstalled rather than skipped. Read the specific
component's predicate rather than generalizing from a sibling — see `controllers/AGENTS.md`.

Relatedly, the schema marks a dozen-plus component fields required even though the defaulting pass would have
filled them, so a hand-written minimal CR is rejected by the API server. For dry-run probes, satisfy them with
empty strings — there is no minimum-length constraint.

## Generated artifacts

The CRD is generator output with all field descriptions stripped to stay under the object-size limit, then patched
further by path-addressed YAML expressions. Consequences:

- `kubectl explain` returns nothing for every field, and the CRD YAML carries structure only. These Go types are
  the documentation.
- The patches address absolute paths, no-op silently when a path moves, and can create phantom nodes when applied
  to a missing path. Nothing verifies that a patch landed.
- The same CRD ships in more than one chart, and the copies have drifted historically. State which copy you used.

Do not hand-edit the CRD. Regenerate, and expect the committed file to differ from raw generator output by the
post-processing.

## The status contract

The transition timestamp is written with Go's `time.Time` string form, including the monotonic-clock reading, and
the schema declares the field a plain string with no date format, so the API server accepts it and no date parser
will. The no-op detection compares that field against itself, so it can never report "unchanged".

The aggregate condition is a transcript of the current pass rather than a settled state. See the root file before
reading it.
