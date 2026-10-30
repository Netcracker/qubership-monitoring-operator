# Test suites

Read the root `AGENTS.md` first. This file covers the traps specific to the tests.

## A passing run can mean nothing ran

Every file in `test/envtests/` carries the `envtest` build tag, so `go test ./...` compiles none of it and exits 0.
Use `make envtest`, or pass `-tags envtest` yourself. A CI job that only provisions control-plane assets is not
evidence that the suite ran.

Before you claim that a behavior is covered, change the line you care about and re-run the tests. Changes to
defaulting rules and component guards have survived the full suite before.

The envtest environment has no scheduler and no kubelet. A custom resource can reach its success condition with zero
ready replicas and no pods, so the suite never tests that workloads come up.

## The envtest suite is order-dependent

`utils.PrivilegedRights` is a mutable package-level global. The `Reconcile with privilegedRights` spec sets it through
a flag and never restores it, while other specs assume the zero value. Ginkgo randomizes container order, so the seed
decides which mode each spec exercises. Pin the seed and assert the mode explicitly; a single green run does not show
which contract was tested.

## Reading the output

Diagnostics written through the test logger go to the Ginkgo writer and are hidden for passing specs unless you pass
`-ginkgo.v` to the test binary. `go test -v` alone is not enough. Several agents saw a clean `ok` line with none of
their own probe output.

Focused specs make the package exit non-zero even when every executed spec passes. Read the spec summary line instead
of the exit code, and never commit a focused spec.

A Go coverage profile lists only packages that the test binary links. A package missing from the report is untested,
not at 0%.

## A baseline failure is probably not yours

`TestPlatformMonitoringCRDManifest` in `api/v1` rebuilds the repository root by walking up three directories and
appending `qubership-monitoring-operator`. The test passes in the standard GitHub Actions checkout layout, but fails in
a Git worktree or in a clone with a different directory name. Run a baseline before you attribute a failure to your
change.

That test's only assertion checks a struct value for nil, so it cannot fail. A pass proves only that the file opened
and decoded.

## Integration tests

The Robot suites are split between this repository and an upstream base image. Keywords and the tag resolver live in
the image, so an unresolved keyword or an apparently uncalled helper is usually defined there. Only `.robot` files run.
Python test files beside them, such as `test/robot-tests/src/lib/test_check_json_object_*.py`, are copied into the
image and run by nothing.

`tags_exclusion.py` treats an unset variable the same as `false`. Its list includes `CONFIGURATIONS-STREAMER` and
`JSON-EXPORTER`, which `templates/integration-tests/deployment.yml` never sets, so those suites never run, even when
the component is installed.

Keywords wrap calls in long retry loops, which makes a hard failure look like a slow start. Read the pod's actual
environment instead of trusting elapsed time. A malformed value rendered into an environment variable is a common
cause; the templates emit a Go formatting placeholder for an undefined value without a render error.
