# The test suites

Read the root `AGENTS.md` first. This file covers the traps specific to the tests.

## A passing run frequently means nothing ran

Every file in the envtest directory carries a build tag, so a plain package-wildcard run compiles none of it,
matches no packages, and **exits 0**. Pass the tag explicitly. A CI job that merely provisions control-plane assets
is not evidence that the suite ran.

The unit suite never calls the top-level reconcile entry point, and the envtest specs call component reconcilers
directly rather than through it, so status conditions are never asserted by anything. Mutations to defaulting rules
and to component guards have survived the full suite. Mutate the line you care about and re-run before claiming a
behavior is covered.

The environment has no scheduler and no kubelet, so a custom resource can reach its success condition with zero
ready replicas and no pods — "the workloads came up" is never actually tested there.

## The suite is order-dependent

The privileged-rights mode is a mutable package-level global that one spec sets through a flag and never restores,
while another assumes the zero value. Container order is randomized, so which mode the specs exercise is decided by
the seed, and neither ordering asserts the negative case. Pin the seed and assert the mode explicitly; a single
green run does not establish which contract was tested.

## Reading the output

Diagnostics written through the test logger go to the Ginkgo writer and are invisible for passing specs unless
verbose Ginkgo output is requested on the test binary — `go test -v` alone is not enough. Several agents saw a
clean `ok` line with zero occurrences of their own probe output.

Focusing specs makes the package exit non-zero even when every executed spec passes, which reads as a failure.
Read the spec summary line rather than the exit code, and never leave focus in a committed test.

A Go coverage profile lists only packages the test binary links, so a package missing from the report is untested
rather than at zero — absence is not a measurement.

## Assets and stubs

The third-party CRDs used by the envtest suite are hand-written stubs that declare no properties and preserve
unknown fields, so the API server accepts any spec. A passing create there proves only that the request was
well-formed JSON. Assert on the object you read back.

The asset custom resource sets the non-default privilege mode, so the default path is not what the suite
exercises.

## A baseline failure is probably not yours

One test in the API-types package reconstructs the repository root by walking up and then re-appending the
project's own directory name, so the whole package fails in a git worktree, in a renamed clone, or on a CI path.
Establish a baseline run before attributing any failure to your change.

That same test's only assertion cannot fail — it checks a struct value for nil — so when it passes it proves only
that a file opened and decoded.

## Integration tests

The Robot suites are split between this repository and an upstream base image: keywords and the tag resolver live
on the image side, so an unresolved keyword or an apparently uncalled helper is usually defined there. Only
`.robot` files are executed — Python test files sitting beside them are copied into the image and run by nothing,
including the only automated regression test for at least one past fix.

The tag-exclusion logic treats an unset variable the same as a disabled one, and the exclusion list names
components the chart never exports, so those suites can never run even when the component is installed.

Keywords wrap calls in long retry loops, which makes a hard failure indistinguishable from a slow start. Read the
pod's actual environment rather than trusting elapsed time — a malformed value rendered into an environment
variable is a common cause, and the templates emit Go's formatting placeholder for an undefined value without any
render error.
