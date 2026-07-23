import importlib.util
import tempfile
import unittest
from pathlib import Path

from ruamel.yaml import YAML

MODULE_PATH = Path(__file__).with_name("crd-update.py")
SPEC = importlib.util.spec_from_file_location("crd_update", MODULE_PATH)
crd_update = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(crd_update)


class ClearOutputDirectoryTest(unittest.TestCase):
    def test_prometheus_cleanup_preserves_other_api_groups(self):
        with tempfile.TemporaryDirectory() as output_dir:
            output_path = Path(output_dir)
            project_crd = output_path / "monitoring.netcracker.com_platformmonitorings.yaml"
            victoriametrics_crd = output_path / "operator.victoriametrics.com_vmagents.yaml"
            prometheus_crd = output_path / "monitoring.coreos.com_podmonitors.yaml"
            project_crd.touch()
            victoriametrics_crd.touch()
            prometheus_crd.touch()

            crd_update.clear_output_dir(output_dir, operator="prometheus")

            self.assertTrue(project_crd.exists())
            self.assertTrue(victoriametrics_crd.exists())
            self.assertFalse(prometheus_crd.exists())

    def test_victoriametrics_cleanup_preserves_prometheus_crds(self):
        with tempfile.TemporaryDirectory() as output_dir:
            output_path = Path(output_dir)
            victoriametrics_crd = output_path / "operator.victoriametrics.com_vmagents.yaml"
            prometheus_crd = output_path / "monitoring.coreos.com_podmonitors.yaml"
            victoriametrics_crd.touch()
            prometheus_crd.touch()

            crd_update.clear_output_dir(output_dir, operator="victoriametrics")

            self.assertFalse(victoriametrics_crd.exists())
            self.assertTrue(prometheus_crd.exists())

    def test_grafana_cleanup_preserves_only_supported_legacy_crd(self):
        with tempfile.TemporaryDirectory() as output_dir:
            output_path = Path(output_dir)
            legacy_dashboard_crd = output_path / "v1alpha1.integreatly.org_grafanadashboards.yaml"
            stale_duplicate_crd = output_path / "integreatly.org_grafanadatasources.yaml"
            current_crd = output_path / "grafana.integreatly.org_grafanadatasources.yaml"
            legacy_dashboard_crd.touch()
            stale_duplicate_crd.touch()
            current_crd.touch()

            crd_update.clear_output_dir(output_dir, operator="grafana")

            self.assertTrue(legacy_dashboard_crd.exists())
            self.assertFalse(stale_duplicate_crd.exists())
            self.assertFalse(current_crd.exists())


class CompactCrdTest(unittest.TestCase):
    def test_remove_descriptions_preserves_schema_fields(self):
        crd = {
            "spec": {
                "versions": [{
                    "schema": {
                        "openAPIV3Schema": {
                            "description": "root description",
                            "properties": {
                                "spec": {
                                    "description": "spec description",
                                    "type": "object",
                                    "properties": {
                                        "description": {
                                            "description": "user-provided description",
                                            "type": "string",
                                        },
                                        "replicas": {"description": "replicas", "type": "integer"},
                                    },
                                }
                            },
                        }
                    }
                }]
            }
        }

        crd_update.remove_descriptions(crd)

        schema = crd["spec"]["versions"][0]["schema"]["openAPIV3Schema"]
        spec_schema = schema["properties"]["spec"]
        self.assertNotIn("description", schema)
        self.assertNotIn("description", spec_schema)
        self.assertIn("description", spec_schema["properties"])
        self.assertNotIn("description", spec_schema["properties"]["description"])
        self.assertEqual("string", spec_schema["properties"]["description"]["type"])
        self.assertEqual("integer", spec_schema["properties"]["replicas"]["type"])

    def test_compact_legacy_grafana_crds_rewrites_preserved_files(self):
        with tempfile.TemporaryDirectory() as output_dir:
            crd_path = Path(output_dir) / "v1alpha1.integreatly.org_grafanadashboards.yaml"
            crd_path.write_text(
                "spec:\n  versions:\n  - schema:\n      openAPIV3Schema:\n"
                "        description: legacy\n        type: object\n",
                encoding="utf-8",
            )
            yaml = YAML()

            crd_update.compact_legacy_grafana_crds(output_dir, yaml)

            compacted_crd = yaml.load(crd_path)
            schema = compacted_crd["spec"]["versions"][0]["schema"]["openAPIV3Schema"]
            self.assertNotIn("description", schema)
            self.assertEqual("object", schema["type"])


if __name__ == "__main__":
    unittest.main()
