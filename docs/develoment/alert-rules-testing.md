# Alert Rules Unit Testing with vmalert-tool

A practical guide to writing unit tests for Prometheus/VictoriaMetrics alert
rules using `vmalert-tool unittest`. This document focuses on what is **not**
obvious from the official documentation: how to design good test data, how to
avoid common pitfalls, and how to wire the tool into CI in a way that scales
across many rules.

For reference material - flags, file schema, CLI options - see the official
docs:

- vmalert-tool: <https://docs.victoriametrics.com/victoriametrics/vmalert-tool/>
- Prometheus unit testing format (vmalert-tool is compatible):
  <https://prometheus.io/docs/prometheus/latest/configuration/unit_testing_rules/>
- Series notation (the `start+stepxcount` syntax):
  <https://prometheus.io/docs/prometheus/latest/configuration/unit_testing_rules/#series>

A working example using everything described here lives under
[`test/alerts-tests/`](../test/alerts-tests/) and is wired into CI via
[`.github/workflows/test-alert-rules-unit-tests.yaml`](../.github/workflows/test-alert-rules-unit-tests.yaml).

## 1. Getting started

### 1.1 Install the tool

`vmalert-tool` ships inside the VictoriaMetrics utils bundle. Pin a version -
the test file schema and matcher behavior have shifted between releases:

```bash
VMALERT_TOOL_VERSION=v1.142.0
wget -q https://github.com/VictoriaMetrics/VictoriaMetrics/releases/download/${VMALERT_TOOL_VERSION}/vmutils-linux-amd64-${VMALERT_TOOL_VERSION}.tar.gz
tar -xf vmutils-linux-amd64-${VMALERT_TOOL_VERSION}.tar.gz
mv vmalert-tool-prod vmalert-tool && chmod +x vmalert-tool
```

### 1.2 Minimum viable test file

A test file references a **rules file** and a list of test cases. Each test
case feeds synthetic series into a tiny in-memory evaluator, advances the clock
to `eval_time`, and asserts that the alert either fires with specific
labels/annotations or does not fire at all.

```yaml
# my-alerts-tests.yaml
rule_files:
  - rules.yaml
evaluation_interval: 1m
tests:
  - interval: 1m
    input_series:
      - series: 'up{job="api"}'
        values: "0x10"          # ten 0s = service is down for 10 minutes
    alert_rule_test:
      - eval_time: 6m
        groupname: ApiAvailability
        alertname: ApiDown
        exp_alerts:
          - exp_labels:
              severity: critical
            exp_annotations:
              summary: "API is down"
```

Run it:

```bash
./vmalert-tool unittest --file my-alerts-tests.yaml
```

### 1.3 Preparing the rules file

`vmalert-tool` expects the **bare** rules YAML - the same shape Prometheus
itself loads, i.e. starting at the `groups:` key. If your rules live inside a
Kubernetes `PrometheusRule` CR, strip the wrapper before running tests. In this
repo the source is at
[`controllers/prometheus-rules/assets/prometheus-rules.yaml`](../controllers/prometheus-rules/assets/prometheus-rules.yaml)
and CI renders it via:

```bash
cp controllers/prometheus-rules/assets/prometheus-rules.yaml test/alerts-tests/rules.yaml
sed -i '1,11d' test/alerts-tests/rules.yaml   # drop the PrometheusRule wrapper
```

The number of lines to drop depends on the wrapper - keep it in sync with the
file's `metadata:` block. A safer alternative is `yq`:

```bash
yq '{"groups": .spec.groups}' prometheus-rules.yaml > rules.yaml
```

---

## 2. Generic recommendations

These are the rules of thumb we landed on after writing tests for dozens of
alerts in this repo. Most of them push back against a tempting shortcut that
silently produces a green-but-useless test.

### 2.1 Every alert rule needs at least two tests

One that **fires** and one that **does not**. Without the negative test you
have no signal that the threshold is meaningful - a rule like
`expr: up == 0` will happily "pass" a positive test even if you accidentally
change it to `up >= 0`.

This is enforced in CI by
[`tests-checker.sh`](../test/alerts-tests/tests-checker.sh):

```bash
expected_tests_count=2
# fails the job if any alert has fewer than 2 cases
```

We recommend adopting the same check. It is 30 lines of `bash` + `yq` and
catches the most common omission.

### 2.2 One test file per rule group

Mirror the structure of the rules file. We split tests by component
(`etcd-alerts-tests.yaml`, `node-alerts-tests.yaml`, …) and the CI step picks
them up with a glob:

```bash
mapfile -t test_files < <(find ${{ env.tests_dir }} -name "*-tests.yaml" | sort)
args=()
for f in "${test_files[@]}"; do
  args+=(--files "$f")
done
./vmalert-tool unittest "${args[@]}"
```

This keeps diffs small when a single rule changes and makes it obvious which
group owns which test.

### 2.3 Test the threshold, not the implementation

Write the inputs so they cross the threshold by a **small, intentional**
margin. Inputs that produce values 100× over the threshold pass even if you
later loosen the rule by an order of magnitude.

### 2.4 Pin the tool version

The test schema and matching semantics are stable, but error messages,
templating edge cases, and histogram-quantile rounding have changed between
releases. Pin `vmalert-tool` to a specific tag and bump it deliberately.

### 2.5 Keep `evaluation_interval` and `interval` aligned with reality

Real Prometheus typically scrapes at 15–30 s and evaluates rules at 30 s–1 min.
Tests should not use absurdly short intervals (e.g. `1s`) to "speed things up" -
some rules use `rate()` or `increase()` over windows like `5m` and need a
realistic step to behave the way they will in production. `1m` for both is a
safe default.

## 3. Writing test data - the hard part

This is what the official docs cover the least. The `input_series.values`
syntax looks innocent but encodes a small DSL that you have to fit your test
data into.

### 3.1 The series notation cheat sheet

```bash
"v"           → single sample with value v
"v+sxN"       → v, v+s, v+2s, ... (N+1 samples)
"v-sxN"       → v, v-s, v-2s, ...
"vxN"         → v repeated N+1 times
"_"           → missing sample (gap)
"stale"       → explicit staleness marker
```

A few worked examples:

| Notation        | Expands to                          | Use when…                              |
| --------------- | ----------------------------------- | -------------------------------------- |
| `"0x1440"`      | `0, 0, 0, ...` (1441 samples)       | Metric is steady at zero for 24h@1m    |
| `"1x1440"`      | `1, 1, 1, ...`                      | Boolean "healthy" gauge                |
| `"0+1x120"`     | `0, 1, 2, ..., 120`                 | Counter increasing by 1 per step       |
| `"100+30000x60"`| `100, 30100, 60100, ...`            | Counter increasing by 30000 per step   |
| `"5 _ _ 8"`     | `5, gap, gap, 8`                    | Reproducing a scrape miss              |

**Tip:** the number after `x` is the count of **additional** samples, not the
total. `0x10` produces eleven zeros - long enough to satisfy a `for: 5m` clause
with a 1-minute interval.

### 3.2 Make `eval_time` longer than `for:` plus one step

If the rule says `for: 5m`, evaluating at exactly `eval_time: 5m` is a race -
some versions of the evaluator consider the alert "pending" at that boundary.
Add a margin:

```yaml
# rule has  for: 5m
alert_rule_test:
  - eval_time: 6m       # 5m for + 1m breathing room
```

The negative test (where the alert should NOT fire) can use the same
`eval_time` for symmetry.

### 3.3 Counters: think in terms of `rate()`

Most "high error rate" alerts use `rate(errors_total[5m]) / rate(total[5m])`.
For the test, choose two parallel counters and pick rates that produce the
ratio you want:

```yaml
# Reproduces "more than 1% errors"
input_series:
  - series: 'grpc_server_handled_total{job="etcd",grpc_code="Error"}'
    values: "400+400x1440"        # +400/min → 6.66/s
  - series: 'grpc_server_handled_total{job="etcd"}'
    values: "30000+30000x1440"    # +30000/min → 500/s, ratio = 1.33%
```

The negative test is just `0x1440` on both - `rate()` of a flat counter is 0,
and `0/0` evaluates to `NaN` (no alert), which is the behavior you want.

### 3.4 Histograms: respect cumulative bucket ordering

`histogram_quantile()` reads `_bucket` series and **requires** monotonically
non-decreasing counts across `le` buckets. If you put a higher count in a
smaller bucket, you get nonsense quantiles.

Pattern that works (from [`etcd-alerts-tests.yaml`](../test/alerts-tests/etcd-alerts-tests.yaml)):

```yaml
# Most samples are above 0.15 → quantile lands in the slow bucket
input_series:
  - series: 'grpc_server_handling_seconds_bucket{le="0.05"}'
    values: "0+1x1440"
  - series: 'grpc_server_handling_seconds_bucket{le="0.1"}'
    values: "0+2x1440"
  - series: 'grpc_server_handling_seconds_bucket{le="0.15"}'
    values: "0+3x1440"
  - series: 'grpc_server_handling_seconds_bucket{le="0.2"}'
    values: "0+1000x1440"   # the jump - slow requests live here
  - series: 'grpc_server_handling_seconds_bucket{le="+Inf"}'
    values: "0+2000x1440"   # +Inf is always required
```

For the negative case, **shift the buckets**: keep the same series and rates
but change the `le` labels so all the traffic lands in fast buckets:

```yaml
  - series: 'grpc_server_handling_seconds_bucket{le="0.005"}'
    values: "0+1x1440"
  - series: 'grpc_server_handling_seconds_bucket{le="0.01"}'
    values: "0+2x1440"
  # ...
  - series: 'grpc_server_handling_seconds_bucket{le="+Inf"}'
    values: "0+1001x1440"
```

Rules of thumb for histogram test data:

1. Always include a `+Inf` bucket. Quantile calculations need it.
2. Bucket counts must be non-decreasing as `le` increases.
3. To shift the quantile, change which `le` boundary holds "most" of the
   traffic - don't rescale the values.

### 3.5 Annotations: keep them tight or be ready to maintain them

`exp_annotations` is matched **exactly**. If your rule template emits something
like:

```bash
description: "Etcd cluster have no leader\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }}"
```

then your test must reproduce the whole string including the `LABELS:
map[...]` rendering and the literal `\n` line breaks. This is brittle - a
harmless template tweak forces test churn.

Two reasonable strategies:

- **Strict (used in this repo):** copy the full rendered string into
  `exp_annotations`. You catch every accidental template change, at the cost
  of edits when you intentionally reword.
- **Minimal:** only assert on the parts you care about (e.g. `summary` only,
  skip `description`). You get less coverage but a stabler test suite.

Pick one and apply it consistently in a given file.

### 3.6 The LABELS map ordering trap

When a template renders `{{ $labels }}`, the labels appear **sorted
alphabetically by key**, including the auto-added `alertname` and `alertgroup`
(VictoriaMetrics-specific). If you write the expected string with labels in a
different order the match fails. Run the test once, copy the actual rendering
from the diff output, then commit.

### 3.7 Empty-label gotcha in templates

If a rule references `{{ $labels.instance }}` but your input series doesn't
have an `instance` label, the template renders as the empty string - you'll
see `"… (instance )"` with a trailing space. That is **not** a bug to fix in
the rule; either add the label to your input series or accept the rendering
verbatim in `exp_annotations`.

### 3.8 `groupname` must match the source rule group

A common silent failure: the test "passes" but actually didn't run because
`groupname` didn't match. Always copy the value from the rules file's
`groups[].name` exactly - case-sensitive.

### 3.9 Patterns for common alert shapes

| Alert shape                               | Input recipe                                  |
| ----------------------------------------- | --------------------------------------------- |
| `metric == 0` for N minutes               | `values: "0xN+5"`                             |
| `metric == 1` for N minutes               | `values: "1xN+5"`                             |
| Counter rate above threshold              | `values: "0+RATExN"` where RATE/60 = req/s    |
| Counter rate below threshold              | `values: "0xN"` (flat)                        |
| Multi-member quorum (count series)        | Emit `metric{label=A}`, `metric{label=B}`, …  |
| Restart loop (`changes()`)                | Use `1+1xN` - each step increments the value  |
| Missing data / scrape failure             | Use `_` in the values string                  |

## 4. Recommended layout for a project

```bash
test/alerts-tests/
├── rules.yaml                        # rendered from the source, .gitignore-able
├── <component>-alerts-tests.yaml     # one file per group
├── tests-checker.sh                  # enforces "at least N tests per rule"
└── README.md                         # link to this guide
```

The `rules.yaml` file is a build artifact - generate it in CI rather than
committing it, so it can't drift from the source.

## 5. CI integration

The shape is the same in any CI system: install the tool, render the rules
file, run `vmalert-tool unittest` over a glob of test files.

### 5.1 GitHub Actions

A full working example lives at
[`.github/workflows/test-alert-rules-unit-tests.yaml`](../.github/workflows/test-alert-rules-unit-tests.yaml).
The essentials:

```yaml
name: "Tests: Alert rules unit-tests"
on:
  pull_request:
  workflow_run:
    workflows: ["Build Artifacts"]
    types: [completed]

permissions:
  contents: read

env:
  vmalert_tool_version: v1.142.0
  tests_dir: test/alerts-tests
  source_alert_rules_file: controllers/prometheus-rules/assets/prometheus-rules.yaml
  target_alert_rules_file: test/alerts-tests/rules.yaml

jobs:
  alert-rules-unit-tests:
    runs-on: ubuntu-latest
    timeout-minutes: 30
    steps:
      - uses: actions/checkout@v6
        with:
          persist-credentials: false

      - name: Install vmalert-tool
        run: |
          wget -q https://github.com/VictoriaMetrics/VictoriaMetrics/releases/download/${{ env.vmalert_tool_version }}/vmutils-linux-amd64-${{ env.vmalert_tool_version }}.tar.gz
          tar -xf vmutils-linux-amd64-${{ env.vmalert_tool_version }}.tar.gz
          mv vmalert-tool-prod vmalert-tool && chmod +x vmalert-tool

      - name: Render rules file
        run: |
          cp ${{ env.source_alert_rules_file }} ${{ env.target_alert_rules_file }}
          sed -i '1,11d' ${{ env.target_alert_rules_file }}

      - name: Enforce minimum test count per rule
        run: ./${{ env.tests_dir }}/tests-checker.sh
        continue-on-error: true        # warn-only; flip to false to enforce hard

      - name: Run unit tests
        run: |
          mapfile -t test_files < <(find ${{ env.tests_dir }} -name "*-tests.yaml" | sort)
          args=()
          for f in "${test_files[@]}"; do
            args+=(--files "$f")
          done
          ./vmalert-tool unittest "${args[@]}"
```

Notes:

- Pin `actions/checkout` to a SHA in production workflows.
- `continue-on-error: true` on the checker step is a good way to roll out the
  "min N tests per rule" policy without immediately failing existing PRs.
  Promote to a hard fail once coverage catches up.

### 5.2 GitLab CI

```yaml
# .gitlab-ci.yml
variables:
  VMALERT_TOOL_VERSION: "v1.142.0"
  TESTS_DIR: "test/alerts-tests"
  SOURCE_RULES: "controllers/prometheus-rules/assets/prometheus-rules.yaml"
  TARGET_RULES: "test/alerts-tests/rules.yaml"

alert-rules-unit-tests:
  stage: test
  image: alpine:3.23
  before_script:
    - apk add --no-cache wget tar bash findutils sed
  script: |-
    wget -q https://github.com/VictoriaMetrics/VictoriaMetrics/releases/download/${VMALERT_TOOL_VERSION}/vmutils-linux-amd64-${VMALERT_TOOL_VERSION}.tar.gz
    tar -xf vmutils-linux-amd64-${VMALERT_TOOL_VERSION}.tar.gz
    mv vmalert-tool-prod vmalert-tool && chmod +x vmalert-tool
    cp "${SOURCE_RULES}" "${TARGET_RULES}"
    sed -i '1,11d' "${TARGET_RULES}"
    bash "${TESTS_DIR}/tests-checker.sh" || true     # warn-only
    mapfile -t test_files < <(find ${{ env.tests_dir }} -name "*-tests.yaml" | sort)
    args=()
    for f in "${test_files[@]}"; do
      args+=(--files "$f")
    done
    ./vmalert-tool unittest "${args[@]}"
  rules:
    - if: '$CI_PIPELINE_SOURCE == "merge_request_event"'
    - changes:
        - controllers/prometheus-rules/assets/prometheus-rules.yaml
        - test/alerts-tests/**/*
```

The `changes:` filter limits the job to PRs that actually touch rules or tests -
useful in monorepos where most pipelines should skip it.

### 5.3 Caching the binary

In both systems the tarball is ~30 MB and downloads in 2–3 s, so caching is
rarely worth the complexity. If you do want to cache, key the cache by
`vmalert_tool_version` so a bump invalidates it automatically.

### 5.4 Failure output and debugging in CI

When `vmalert-tool unittest` fails, it prints the difference between expected
and actual alerts. The most common failure modes:

1. **`alert is not firing`** - your input series didn't produce the value the
   rule expects. Re-check rates, histogram bucket layout, and whether
   `eval_time` is past `for:`.
2. **`unexpected metric ... in result`** - labels don't match. Check the
   `groupname` and any `keep_firing_for`/`labels:` blocks on the rule.
3. **`expected annotations do not match`** - almost always a templating
   change. Diff the two strings character-by-character; LABELS map ordering
   and trailing spaces are the usual culprits.

For local debugging, run a single file: `./vmalert-tool unittest --file
test/alerts-tests/etcd-alerts-tests.yaml`. The output is the same as in CI but
without the noise of unrelated groups.

## 6. Quick reference checklist

When adding a new alert rule, before you open the PR:

- [ ] At least one positive test (alert fires) and one negative test (alert
      does not fire), both pointing at the new `alertname`.
- [ ] `groupname` matches the source rule group exactly.
- [ ] `eval_time` is at least `for:` duration + one interval.
- [ ] Input rates picked to cross the threshold by a small, intentional
      margin - not 100× over.
- [ ] For histogram-based alerts: `+Inf` bucket present, counts non-decreasing
      by `le`.
- [ ] `exp_labels` includes every label the rule sets (typically `severity`).
- [ ] `exp_annotations` either matches the template output exactly, or is
      asserted on a stable subset.
- [ ] `tests-checker.sh` (or equivalent) passes locally.

## 7. Further reading

- [vmalert-tool official docs](https://docs.victoriametrics.com/victoriametrics/vmalert-tool/)
- [Prometheus unit testing for alerting rules](https://prometheus.io/docs/prometheus/latest/configuration/unit_testing_rules/)
- [Prometheus `histogram_quantile()`](https://prometheus.io/docs/prometheus/latest/querying/functions/#histogram_quantile)
- Working example in this repo: [`test/alerts-tests/`](../test/alerts-tests/)
  and [`.github/workflows/test-alert-rules-unit-tests.yaml`](../.github/workflows/test-alert-rules-unit-tests.yaml)
