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

- Download: builds the GitHub release URL from a per-operator template (bundle.yaml / crd.yaml / crds.yaml).
- Sanitize: replaces smart quotes (including U+201D ”), en/em dashes, ellipsis, and NBSP before YAML parsing so
  malformed descriptions don't crash the loader.
- Split: iterates all documents, keeps only `kind: CustomResourceDefinition`, and writes one file per CRD as
  `<group>_<plural>.yaml`.
- Compact: removes OpenAPI `description` fields while preserving validation schemas. This keeps the packaged chart
  below the Kubernetes Secret size limit used by Helm release storage.
- Store: `make update-prometheus-crds` writes the canonical Prometheus CRDs into the Prometheus subchart and
  synchronizes them into the VictoriaMetrics subchart. This preserves independent installations of either operator.
- Clean: each updater removes only the CRDs owned by its API group, so updating VictoriaMetrics preserves the shared
  Prometheus CRDs.
- Compatibility: Grafana updates retain the existing `integreatly.org/v1alpha1` dashboard CRD used by the dashboard
  converter. Other files from older Grafana layouts are removed before writing the current CRDs.

Run `make docs` after updating operator CRDs. It rebuilds `docs/crds` and updates the dedicated
`qubership-monitoring-crds` chart.

Add next annotations:

- Always adds `helm.sh/hook: crd-install` and `helm.sh/hook-weight: "-5"`
- For VictoriaMetrics adds `operator.victoriametrics.com/version: <version>`
- For Grafana adds `operator.grafana.com/version: <version>`
- For Prometheus, leaves the existing `operator.prometheus.io/version` alone since it already ships in the bundle

## Usage examples

```bash
python crd-update.py -o victoriametrics -v 0.69.0 -d output/vm
python crd-update.py -o prometheus     -v 0.90.1 -d output/prom
python crd-update.py -o grafana        -v 5.22.2 -d output/grafana
```
