import argparse
import io
import logging
import os
import sys
import urllib.request

from ruamel.yaml import YAML

log = logging.getLogger("crd-update")


OPERATORS = {
    "prometheus": {
        "url": "https://github.com/prometheus-operator/prometheus-operator/releases/download/v{version}/bundle.yaml",
        "version_annotation": None,  # already present in CRDs
    },
    "victoriametrics": {
        "url": "https://github.com/VictoriaMetrics/operator/releases/download/v{version}/crd.yaml",
        "version_annotation": "operator.victoriametrics.com/version",
    },
    "grafana": {
        "url": "https://github.com/grafana/grafana-operator/releases/download/v{version}/crds.yaml",
        "version_annotation": "operator.grafana.com/version",
    },
}

# Map of common non-UTF-8 / "smart" characters that appear in CRD descriptions
# and break YAML parsers or downstream tooling. Replace with safe ASCII equivalents.
CHAR_REPLACEMENTS = {
    "‘": "'",    # left single quotation mark
    "’": "'",    # right single quotation mark
    "“": '"',    # left double quotation mark
    "”": '"',    # right double quotation mark
    "–": "-",    # en dash
    "—": "-",    # em dash
    "…": "...",  # horizontal ellipsis
    " ": " ",    # non-breaking space
}

HELM_HOOK_ANNOTATIONS = {
    "helm.sh/hook": "crd-install",
    "helm.sh/hook-weight": "-5",
}

LEGACY_GRAFANA_CRDS = {
    "v1alpha1.integreatly.org_grafanadashboards.yaml",
}

CRD_FILE_PREFIXES = {
    "prometheus": "monitoring.coreos.com_",
    "victoriametrics": "operator.victoriametrics.com_",
}


def download(url):
    log.info("Downloading %s", url)
    with urllib.request.urlopen(url) as resp:
        return resp.read().decode("utf-8", errors="replace")


def sanitize(text):
    for bad, good in CHAR_REPLACEMENTS.items():
        text = text.replace(bad, good)
    return text


def ensure_annotations(doc, extra_annotations):
    metadata = doc.setdefault("metadata", {})
    annotations = metadata.get("annotations") or {}
    for key, value in extra_annotations.items():
        annotations[key] = value
    metadata["annotations"] = annotations


def remove_descriptions(value):
    versions = value.get("spec", {}).get("versions", [])
    for version in versions:
        schema = version.get("schema", {}).get("openAPIV3Schema")
        if schema:
            remove_schema_descriptions(schema)


def remove_schema_descriptions(schema):
    schema.pop("description", None)

    for child in (schema.get("properties") or {}).values():
        remove_schema_descriptions(child)

    for keyword in ("items", "additionalProperties", "not"):
        child = schema.get(keyword)
        if isinstance(child, dict):
            remove_schema_descriptions(child)
        elif isinstance(child, list):
            for item in child:
                remove_schema_descriptions(item)

    for keyword in ("allOf", "anyOf", "oneOf"):
        for child in schema.get(keyword) or []:
            remove_schema_descriptions(child)

    for keyword in ("$defs", "definitions", "dependentSchemas", "patternProperties"):
        for child in (schema.get(keyword) or {}).values():
            remove_schema_descriptions(child)


def filename_for(crd):
    group = crd["spec"]["group"].lower()
    plural = crd["spec"]["names"]["plural"].lower()
    return f"{group}_{plural}.yaml"


def clear_output_dir(output_dir, operator):
    for name in os.listdir(output_dir):
        if operator == "grafana" and name in LEGACY_GRAFANA_CRDS:
            continue
        prefix = CRD_FILE_PREFIXES.get(operator)
        if prefix and not name.startswith(prefix):
            continue
        if name.endswith((".yaml", ".yml")):
            os.remove(os.path.join(output_dir, name))


def compact_legacy_grafana_crds(output_dir, yaml):
    for name in LEGACY_GRAFANA_CRDS:
        path = os.path.join(output_dir, name)
        if not os.path.exists(path):
            continue

        with open(path, encoding="utf-8") as source:
            doc = yaml.load(source)
        remove_descriptions(doc)
        with open(path, "w", encoding="utf-8") as target:
            yaml.dump(doc, target)


def process(operator, version, output_dir):
    spec = OPERATORS[operator]
    url = spec["url"].format(version=version)
    raw = sanitize(download(url))

    yaml = YAML()
    yaml.width = 4096
    yaml.preserve_quotes = True

    docs = list(yaml.load_all(io.StringIO(raw)))

    os.makedirs(output_dir, exist_ok=True)
    clear_output_dir(output_dir, operator)

    written = 0
    for doc in docs:
        if not doc:
            continue
        if doc.get("kind", "").lower() != "customresourcedefinition":
            continue

        remove_descriptions(doc)
        annotations = dict(HELM_HOOK_ANNOTATIONS)
        if spec["version_annotation"]:
            annotations[spec["version_annotation"]] = version
        ensure_annotations(doc, annotations)

        out_path = os.path.join(output_dir, filename_for(doc))
        with open(out_path, "w", encoding="utf-8") as f:
            yaml.dump(doc, f)
        log.info("Wrote %s", out_path)
        written += 1

    if operator == "grafana":
        compact_legacy_grafana_crds(output_dir, yaml)

    if written == 0:
        log.warning("No CRDs found in %s", url)
    else:
        log.info("Wrote %d CRD(s) for %s %s into %s", written, operator, version, output_dir)


def main():
    logging.basicConfig(
        level=logging.INFO,
        format="%(asctime)s %(levelname)-5s[%(name)s] - %(message)s",
        datefmt="%Y-%m-%d %H:%M:%S",
    )

    parser = argparse.ArgumentParser(description="Download and prepare operator CRDs for Helm packaging.")
    parser.add_argument("--operator", "-o", required=True, choices=sorted(OPERATORS.keys()),
                        help="Operator whose CRDs to fetch. Available values: prometheus, victoriametrics, grafana.")
    parser.add_argument("--version", "-v", required=True,
                        help="Operator release version (without leading 'v', e.g. 0.90.1)")
    parser.add_argument("--output-dir", "-d", default="output",
                        help="Directory to write split CRD files into")

    args = parser.parse_args()
    version = args.version.lstrip("v")

    try:
        process(args.operator, version, args.output_dir)
    except Exception as e:
        log.error("Failed: %s", e)
        sys.exit(1)


if __name__ == "__main__":
    main()
