import importlib.util
import tempfile
import unittest
from pathlib import Path
from unittest.mock import patch

MODULE_PATH = Path(__file__).with_name("crd-update.py")
SPEC = importlib.util.spec_from_file_location("crd_update", MODULE_PATH)
crd_update = importlib.util.module_from_spec(SPEC)
SPEC.loader.exec_module(crd_update)


class OutputDirectoryValidationTest(unittest.TestCase):
    def test_accepts_directory_inside_repository(self):
        repository_root = MODULE_PATH.parents[2]
        with tempfile.TemporaryDirectory(dir=repository_root) as output_dir:
            self.assertEqual(
                Path(output_dir).resolve(),
                crd_update.resolve_output_dir(output_dir),
            )

    def test_rejects_directory_outside_repository(self):
        with tempfile.TemporaryDirectory() as output_dir:
            with self.assertRaisesRegex(
                ValueError,
                "must be inside the repository",
            ):
                crd_update.resolve_output_dir(output_dir)

    def test_rejects_symlink_escape_from_repository(self):
        repository_root = MODULE_PATH.parents[2]
        with (
            tempfile.TemporaryDirectory(dir=repository_root) as local_dir,
            tempfile.TemporaryDirectory() as external_dir,
        ):
            output_dir = Path(local_dir) / "output"
            output_dir.symlink_to(external_dir, target_is_directory=True)

            with self.assertRaisesRegex(
                ValueError,
                "must be inside the repository",
            ):
                crd_update.resolve_output_dir(output_dir)


class ClearOutputDirectoryTest(unittest.TestCase):
    def test_prometheus_cleanup_preserves_other_api_groups(self):
        with tempfile.TemporaryDirectory() as output_dir:
            output_path = Path(output_dir)
            project_crd = (
                output_path
                / "monitoring.netcracker.com_platformmonitorings.yaml"
            )
            victoriametrics_crd = (
                output_path
                / "operator.victoriametrics.com_vmagents.yaml"
            )
            prometheus_crd = (
                output_path
                / "monitoring.coreos.com_podmonitors.yaml"
            )
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
            victoriametrics_crd = (
                output_path
                / "operator.victoriametrics.com_vmagents.yaml"
            )
            prometheus_crd = (
                output_path
                / "monitoring.coreos.com_podmonitors.yaml"
            )
            victoriametrics_crd.touch()
            prometheus_crd.touch()

            crd_update.clear_output_dir(output_dir, operator="victoriametrics")

            self.assertFalse(victoriametrics_crd.exists())
            self.assertTrue(prometheus_crd.exists())

    def test_grafana_cleanup_preserves_legacy_and_unrelated_yaml(self):
        with tempfile.TemporaryDirectory() as output_dir:
            output_path = Path(output_dir)
            legacy_dashboard_crd = (
                output_path
                / "v1alpha1.integreatly.org_grafanadashboards.yaml"
            )
            stale_duplicate_crd = (
                output_path
                / "integreatly.org_grafanadatasources.yaml"
            )
            current_crd = (
                output_path
                / "grafana.integreatly.org_grafanadatasources.yaml"
            )
            unrelated_yaml = output_path / "custom-resource.yaml"
            legacy_dashboard_crd.touch()
            stale_duplicate_crd.touch()
            current_crd.touch()
            unrelated_yaml.touch()

            crd_update.clear_output_dir(output_dir, operator="grafana")

            self.assertTrue(legacy_dashboard_crd.exists())
            self.assertTrue(unrelated_yaml.exists())
            self.assertFalse(stale_duplicate_crd.exists())
            self.assertFalse(current_crd.exists())


class OperatorSourceTest(unittest.TestCase):
    def test_operators_use_required_upstream_crd_artifacts(self):
        self.assertTrue(
            crd_update.OPERATORS["prometheus"]["url"].endswith(
                "/stripped-down-crds.yaml"
            )
        )
        self.assertTrue(
            crd_update.OPERATORS["victoriametrics"]["url"].endswith(
                "/install-no-webhook.yaml"
            )
        )
        self.assertTrue(
            crd_update.OPERATORS["grafana"]["url"].endswith(
                "/crds.yaml"
            )
        )


class ProcessTest(unittest.TestCase):
    def test_process_does_not_add_obsolete_helm_hooks(self):
        crd = """\
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: vmagents.operator.victoriametrics.com
spec:
  group: operator.victoriametrics.com
  names:
    plural: vmagents
"""

        with tempfile.TemporaryDirectory(
            dir=MODULE_PATH.parents[2],
        ) as output_dir:
            with patch.object(crd_update, "download", return_value=crd):
                crd_update.process(
                    "victoriametrics",
                    "0.73.1",
                    output_dir,
                )

            output = (
                Path(output_dir)
                / "operator.victoriametrics.com_vmagents.yaml"
            ).read_text(encoding="utf-8")
            self.assertNotIn("helm.sh/hook", output)
            self.assertIn(
                "operator.victoriametrics.com/version: 0.73.1",
                output,
            )

    def test_process_keeps_victoriametrics_crds_valid_on_kubernetes_125(self):
        crd = """\
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: vmclusters.operator.victoriametrics.com
spec:
  group: operator.victoriametrics.com
  names:
    plural: vmclusters
  versions:
    - schema:
        openAPIV3Schema:
          properties:
            spec:
              properties:
                vpa:
                  properties:
                    resourcePolicy:
                      properties:
                        containerPolicies:
                          items:
                            properties:
                              startupBoost:
                                properties:
                                  cpu:
                                    x-kubernetes-validations:
                                      - rule: self.isQuantity()
                              retained:
                                x-kubernetes-validations:
                                  - rule: self != ''
"""

        with tempfile.TemporaryDirectory(
            dir=MODULE_PATH.parents[2],
        ) as output_dir:
            with patch.object(crd_update, "download", return_value=crd):
                crd_update.process(
                    "victoriametrics",
                    "0.73.1",
                    output_dir,
                )

            output_path = (
                Path(output_dir)
                / "operator.victoriametrics.com_vmclusters.yaml"
            )
            output = crd_update.YAML().load(
                output_path.read_text(encoding="utf-8")
            )
            properties = (
                output["spec"]["versions"][0]["schema"]
                ["openAPIV3Schema"]["properties"]["spec"]["properties"]
                ["vpa"]["properties"]["resourcePolicy"]["properties"]
                ["containerPolicies"]["items"]["properties"]
            )
            self.assertNotIn(
                "x-kubernetes-validations",
                properties["startupBoost"]["properties"]["cpu"],
            )
            self.assertIn(
                "x-kubernetes-validations",
                properties["retained"],
            )

    def test_process_keeps_embedded_httproute_valid_on_kubernetes_125(self):
        crd = """\
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: grafanas.grafana.integreatly.org
spec:
  group: grafana.integreatly.org
  names:
    plural: grafanas
  versions:
    - schema:
        openAPIV3Schema:
          properties:
            spec:
              properties:
                httpRoute:
                  properties:
                    spec:
                      properties:
                        rules:
                          items:
                            properties:
                              timeouts:
                                x-kubernetes-validations:
                                  - rule: self.request != ''
                retained:
                  x-kubernetes-validations:
                    - rule: self != ''
"""

        with tempfile.TemporaryDirectory(
            dir=MODULE_PATH.parents[2],
        ) as output_dir:
            with patch.object(crd_update, "download", return_value=crd):
                crd_update.process(
                    "grafana",
                    "5.24.0",
                    output_dir,
                )

            output_path = (
                Path(output_dir)
                / "grafana.integreatly.org_grafanas.yaml"
            )
            output = crd_update.YAML().load(
                output_path.read_text(encoding="utf-8")
            )
            properties = (
                output["spec"]["versions"][0]["schema"]
                ["openAPIV3Schema"]["properties"]["spec"]["properties"]
            )
            timeouts = (
                properties["httpRoute"]["properties"]["spec"]["properties"]
                ["rules"]["items"]["properties"]["timeouts"]
            )
            self.assertNotIn("x-kubernetes-validations", timeouts)
            self.assertIn(
                "x-kubernetes-validations",
                properties["retained"],
            )


if __name__ == "__main__":
    unittest.main()
