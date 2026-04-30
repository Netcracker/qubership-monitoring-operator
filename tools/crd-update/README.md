# CRD Update

## CLI

```bash
python crd-update.py --operator {prometheus,victoriametrics,grafana} --version xxx --output-dir path/to/output
```

## Logic

- Download: builds the GitHub release URL from a per-operator template (bundle.yaml / crd.yaml / crds.yaml).
- Sanitize: replaces smart quotes (incl. U+201D ”), en/em dashes, ellipsis, NBSP — runs on the raw text before YAML parsing so malformed
descriptions don't crash the loader.
- Split: iterates all docs, keeps only kind: CustomResourceDefinition, and writes one file per CRD as <group>_<plural>.yaml.

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
