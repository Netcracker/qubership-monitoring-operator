# CRD Update

## CLI

Install the script dependency before running the updater:

```bash
python -m pip install -r requirements.txt
```

```bash
python crd-update.py --operator {prometheus,victoriametrics,grafana} --version xxx --output-dir path/to/output
```

## Logic

- Download: uses the compact CRD artifacts published by Prometheus Operator and VictoriaMetrics Operator, and the
  standard CRD artifact published by Grafana Operator.
- Sanitize: replaces smart quotes (including U+201D ”), en/em dashes, ellipsis, and NBSP before YAML parsing so
  malformed descriptions don't crash the loader.
- Split: iterates all documents, keeps only `kind: CustomResourceDefinition`, and writes one file per CRD as
  `<group>_<plural>.yaml`.
- Compact: uses the pinned `bin/yq-4.53.2` tool to remove schema documentation after every update while preserving CRD
  properties named `description`. Prometheus Operator publishes `stripped-down-crds.yaml`, and VictoriaMetrics
  Operator includes compact CRDs in `install-no-webhook.yaml`.
- Store: `make update-prometheus-crds` writes the canonical Prometheus CRDs into the Prometheus subchart and
  synchronizes them into the VictoriaMetrics subchart. This preserves independent installations of either operator.
- Clean: each updater removes only the CRDs owned by its API group, so updating VictoriaMetrics preserves the shared
  Prometheus CRDs.
- Compatibility: Grafana updates retain the existing `integreatly.org/v1alpha1` dashboard CRD used by the dashboard
  converter. Other files from older Grafana layouts are removed before writing the current CRDs.
- Kubernetes 1.25 compatibility: removes only the high-cost CEL rules from embedded Grafana HTTPRoute schemas and
  VictoriaMetrics `startupBoost.cpu` schemas. The destination HTTPRoute CRD still validates generated routes.

Managed-operator CRDs are canonical in their subchart `crds/` directories. Run `make update-crds` after updating
them to generate the combined CRD set in `charts/qubership-monitoring-crds/crds/`.

Run `make test-crd-compaction` to validate the yq transform before updating CRDs. The fixture verifies that schema
documentation is removed while a CRD property named `description` remains a string.

Preserve or add these version annotations:

- For VictoriaMetrics adds `operator.victoriametrics.com/version: <version>`
- For Grafana adds `operator.grafana.com/version: <version>`
- For Prometheus, leaves the existing `operator.prometheus.io/version` alone since it already ships in the bundle

The updater does not add the obsolete Helm 2 `crd-install` hook. Helm 3 loads CRDs from chart `crds/` directories.

## Usage examples

```bash
python crd-update.py -o victoriametrics -v 0.69.0 -d output/vm
python crd-update.py -o prometheus     -v 0.90.1 -d output/prom
python crd-update.py -o grafana        -v 5.22.2 -d output/grafana
```
