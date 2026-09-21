import argparse
import io
import logging
import os
import sys
import urllib.request
from pathlib import Path

from ruamel.yaml import YAML

log = logging.getLogger("crd-update")

REPOSITORY_ROOT = Path(__file__).resolve().parents[2]

OPERATORS = {
    "prometheus": {
        "url": (
            "https://github.com/prometheus-operator/prometheus-operator/"
            "releases/download/"
            "v{version}/stripped-down-crds.yaml"
        ),
        "version_annotation": None,  # already present in CRDs
    },
    "victoriametrics": {
        "url": (
            "https://github.com/VictoriaMetrics/operator/releases/download/"
            "v{version}/install-no-webhook.yaml"
        ),
        "version_annotation": "operator.victoriametrics.com/version",
    },
    "grafana": {
        "url": (
            "https://github.com/grafana/grafana-operator/releases/download/"
            "v{version}/crds.yaml"
        ),
        "version_annotation": "operator.grafana.com/version",
    },
}

# Map of common non-UTF-8 / "smart" characters that appear in CRD
# descriptions and break YAML parsers or downstream tooling.
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

LEGACY_GRAFANA_CRDS = {
    "v1alpha1.integreatly.org_grafanadashboards.yaml",
}

GRAFANA_CRDS_TO_REPLACE = {
    "grafana.integreatly.org_grafanaalertrulegroups.yaml",
    "grafana.integreatly.org_grafanacontactpoints.yaml",
    "grafana.integreatly.org_grafanadashboards.yaml",
    "grafana.integreatly.org_grafanadatasources.yaml",
    "grafana.integreatly.org_grafanafolders.yaml",
    "grafana.integreatly.org_grafanalibrarypanels.yaml",
    "grafana.integreatly.org_grafanamanifests.yaml",
    "grafana.integreatly.org_grafanamutetimings.yaml",
    "grafana.integreatly.org_grafananotificationpolicies.yaml",
    "grafana.integreatly.org_grafananotificationpolicyroutes.yaml",
    "grafana.integreatly.org_grafananotificationtemplates.yaml",
    "grafana.integreatly.org_grafanas.yaml",
    "grafana.integreatly.org_grafanaserviceaccounts.yaml",
    "integreatly.org_grafanadashboards.yaml",
    "integreatly.org_grafanadatasources.yaml",
    "integreatly.org_grafanafolders.yaml",
    "integreatly.org_grafananotificationchannels.yaml",
    "integreatly.org_grafanas.yaml",
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


def remove_validation_rules(node):
    if isinstance(node, dict):
        node.pop("x-kubernetes-validations", None)
        for value in node.values():
            remove_validation_rules(value)
    elif isinstance(node, list):
        for value in node:
            remove_validation_rules(value)


def preserve_kubernetes_125_compatibility(node, operator):
    if isinstance(node, dict):
        if operator == "victoriametrics":
            startup_boost = node.get("startupBoost")
            if isinstance(startup_boost, dict):
                properties = startup_boost.get("properties") or {}
                cpu = properties.get("cpu")
                if isinstance(cpu, dict):
                    cpu.pop("x-kubernetes-validations", None)
        elif operator == "grafana":
            http_route = node.get("httpRoute")
            if isinstance(http_route, dict):
                remove_validation_rules(http_route)
        for value in node.values():
            preserve_kubernetes_125_compatibility(value, operator)
    elif isinstance(node, list):
        for value in node:
            preserve_kubernetes_125_compatibility(value, operator)


def filename_for(crd):
    group = crd["spec"]["group"].lower()
    plural = crd["spec"]["names"]["plural"].lower()
    return f"{group}_{plural}.yaml"


def resolve_output_dir(output_dir):
    resolved = Path(output_dir).resolve()
    if resolved == REPOSITORY_ROOT or REPOSITORY_ROOT not in resolved.parents:
        raise ValueError(
            f"output directory must be inside the repository: {resolved}"
        )
    return resolved


def clear_output_dir(output_dir, operator):
    for name in os.listdir(output_dir):
        if operator == "grafana" and name in LEGACY_GRAFANA_CRDS:
            continue
        if operator == "grafana" and name not in GRAFANA_CRDS_TO_REPLACE:
            continue
        prefix = CRD_FILE_PREFIXES.get(operator)
        if prefix and not name.startswith(prefix):
            continue
        if name.endswith((".yaml", ".yml")):
            os.remove(os.path.join(output_dir, name))


def process(operator, version, output_dir):
    output_dir = resolve_output_dir(output_dir)
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

        if operator in {"victoriametrics", "grafana"}:
            # Kubernetes 1.25 rejects high-cost CEL rules nested under
            # unbounded third-party API lists.
            preserve_kubernetes_125_compatibility(doc, operator)

        annotations = {}
        if spec["version_annotation"]:
            annotations[spec["version_annotation"]] = version
        if annotations:
            ensure_annotations(doc, annotations)

        out_path = os.path.join(output_dir, filename_for(doc))
        with open(out_path, "w", encoding="utf-8") as f:
            yaml.dump(doc, f)
        log.info("Wrote %s", out_path)
        written += 1

    if written == 0:
        log.warning("No CRDs found in %s", url)
    else:
        log.info(
            "Wrote %d CRD(s) for %s %s into %s",
            written,
            operator,
            version,
            output_dir,
        )


def main():
    logging.basicConfig(
        level=logging.INFO,
        format="%(asctime)s %(levelname)-5s[%(name)s] - %(message)s",
        datefmt="%Y-%m-%d %H:%M:%S",
    )

    parser = argparse.ArgumentParser(
        description="Download and prepare operator CRDs for Helm packaging."
    )
    parser.add_argument(
        "--operator",
        "-o",
        required=True,
        choices=sorted(OPERATORS),
        help="Operator to fetch: prometheus, victoriametrics, or grafana.",
    )
    parser.add_argument(
        "--version",
        "-v",
        required=True,
        help="Operator release version without the leading 'v'.",
    )
    parser.add_argument(
        "--output-dir",
        "-d",
        default="output",
        help="Directory to write split CRD files into.",
    )

    args = parser.parse_args()
    version = args.version.lstrip("v")

    try:
        process(args.operator, version, args.output_dir)
    except Exception as e:
        log.error("Failed: %s", e)
        sys.exit(1)


if __name__ == "__main__":
    main()
