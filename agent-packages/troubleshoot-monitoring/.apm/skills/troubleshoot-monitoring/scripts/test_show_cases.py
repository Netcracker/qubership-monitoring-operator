"""CLI contract for show_cases.py: list catalog symptoms or print one section."""

import subprocess
import sys
import unittest
from pathlib import Path

SCRIPT = Path(__file__).resolve().parent / "show_cases.py"
CATALOG = Path(__file__).resolve().parent.parent / "references" / "troubleshooting.md"
INFRASTRUCTURE_CATALOG = (
    Path(__file__).resolve().parent.parent / "references" / "infrastructure" / "troubleshooting.md"
)
FIXTURE = Path(__file__).resolve().parent / "testdata" / "catalog" / "troubleshooting.md"


def run_helper(*args: str) -> subprocess.CompletedProcess[str]:
    return subprocess.run(
        [sys.executable, str(SCRIPT), *args],
        check=False,
        capture_output=True,
        text=True,
    )


class ShowCasesTests(unittest.TestCase):
    def test_list_prints_symptoms_and_omits_root_cause(self):
        result = run_helper(str(FIXTURE))
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertIn("### Existing resource conflict", result.stdout)
        self.assertIn("helm.go:75", result.stdout)
        self.assertNotIn("Helm refuses to adopt", result.stdout)
        self.assertNotIn("kubectl annotate", result.stdout)

    def test_print_loads_one_section_by_heading(self):
        result = run_helper(str(FIXTURE), "Existing resource conflict")
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertIn("### Existing resource conflict", result.stdout)
        self.assertIn("Helm refuses to adopt", result.stdout)
        self.assertNotIn("### Other case", result.stdout)

    def test_live_catalog_still_lists_monitoring_operator_cases(self):
        result = run_helper(str(CATALOG))
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertIn("Rendered manifests contain a new resource that already exists", result.stdout)
        self.assertIn("Metrics absent or errors during metrics collection", result.stdout)
        self.assertNotIn("Helm refuses to adopt or overwrite", result.stdout)

    def test_infrastructure_catalog_lists_pgskipper_symptoms(self):
        result = run_helper(str(INFRASTRUCTURE_CATALOG))
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertIn("### Collection scrape timeout", result.stdout)
        self.assertNotIn("This case is the chart", result.stdout)

    def test_infrastructure_catalog_loads_one_case_by_heading(self):
        result = run_helper(str(INFRASTRUCTURE_CATALOG), "Collection scrape timeout")
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertIn("### Collection scrape timeout", result.stdout)
        self.assertIn("This case is the chart", result.stdout)
        self.assertNotIn("### Target unavailable", result.stdout)

    def test_infrastructure_catalog_lists_each_component_group(self):
        result = run_helper(str(INFRASTRUCTURE_CATALOG))
        self.assertEqual(result.returncode, 0, result.stderr)
        for heading in (
            "### Collector is not installed",
            "### Prometheus exporter is not installed",
            "### Cassandra monitoring is not enabled",
            "### Redis monitoring agent is not installed",
            "### ClickHouse ServiceMonitor is not enabled",
            "### Kafka Monitoring is not installed",
            "### ZooKeeper Monitoring is not installed",
            "### Consul monitoring CRs are not installed",
            "### RabbitMQ Prometheus monitoring is not installed",
            "### Site-manager monitor is not rendered",
            "### OpenSearch monitoring is not installed",
            "### Monitor excluded",
        ):
            with self.subTest(heading=heading):
                self.assertIn(heading, result.stdout)

    def test_infrastructure_catalog_loads_one_case_from_each_component_group(self):
        samples = (
            ("Prometheus exporter is not installed", "mongodb-prometheus-exporter"),
            ("Cassandra monitoring is not enabled", "monitoringAgent.install"),
            ("Redis monitoring agent is not installed", "redis-monitoring-agent"),
            ("ClickHouse ServiceMonitor is not enabled", "clickhouseCluster.serviceMonitor"),
            ("Kafka Monitoring is not installed", "kafka-monitoring"),
            ("ZooKeeper Grafana dashboard is absent", "zookeeper-grafana-dashboard"),
            ("Consul ACL token scrape failure", "bearerTokenSecret"),
            ("RabbitMQ per-queue metrics are disabled", "perQueueMetrics"),
            ("paas-geo-monitor is not installed", "paasGeoMonitor.install"),
            ("OpenSearch indices dashboard is not enabled", "includeIndices"),
        )
        for heading, marker in samples:
            with self.subTest(heading=heading):
                result = run_helper(str(INFRASTRUCTURE_CATALOG), heading)
                self.assertEqual(result.returncode, 0, result.stderr)
                self.assertIn("### " + heading, result.stdout)
                self.assertIn(marker, result.stdout)
                self.assertNotIn("### Target unavailable", result.stdout)


if __name__ == "__main__":
    unittest.main()
