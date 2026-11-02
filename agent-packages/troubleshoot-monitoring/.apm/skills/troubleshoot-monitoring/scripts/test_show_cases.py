"""CLI contract for show_cases.py: list catalog symptoms or print one section."""

import subprocess
import sys
import unittest
from pathlib import Path

SCRIPT = Path(__file__).resolve().parent / "show_cases.py"
CATALOG = Path(__file__).resolve().parent.parent / "references" / "troubleshooting.md"
INFRASTRUCTURE_CATALOGS = Path(__file__).resolve().parent.parent / "references" / "infrastructure"
FIXTURE = Path(__file__).resolve().parent / "testdata" / "catalog" / "troubleshooting.md"
COMPONENT_CATALOGS = (
    ("pgskipper-operator.md", "Collector is not installed", "metricCollector.install"),
    ("mongodb-operator.md", "Prometheus exporter is not installed", "mongodb-prometheus-exporter"),
    ("cassandra-operator.md", "Cassandra monitoring is not enabled", "monitoringAgent.install"),
    ("redis.md", "Redis monitoring agent is not installed", "redis-monitoring-agent"),
    ("clickhouse-operator-helm.md", "ClickHouse ServiceMonitor is not enabled", "clickhouseCluster.serviceMonitor"),
    ("kafka.md", "Kafka Monitoring is not installed", "kafka-monitoring"),
    ("zookeeper.md", "ZooKeeper Monitoring is not installed", "zookeeper-monitoring"),
    ("consul.md", "Consul monitoring CRs are not installed", "monitoring"),
    ("rabbitmq.md", "RabbitMQ Prometheus monitoring is not installed", "rabbitmq_prometheus"),
    ("drnavigator.md", "paas-geo-monitor is not installed", "paasGeoMonitor.install"),
    ("opensearch.md", "OpenSearch monitoring is not installed", "opensearch-monitoring"),
)
SHARED_DISCOVERY_CATALOG = INFRASTRUCTURE_CATALOGS / "shared-discovery.md"


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

    def test_component_catalogs_list_component_symptoms(self):
        for filename, heading, _ in COMPONENT_CATALOGS:
            with self.subTest(filename=filename):
                result = run_helper(str(INFRASTRUCTURE_CATALOGS / filename))
                self.assertEqual(result.returncode, 0, result.stderr)
                self.assertIn(heading, result.stdout)
                self.assertNotIn("**Root cause:**", result.stdout)

    def test_component_catalogs_load_one_component_case(self):
        for filename, heading, marker in COMPONENT_CATALOGS:
            with self.subTest(filename=filename, heading=heading):
                result = run_helper(str(INFRASTRUCTURE_CATALOGS / filename), heading)
                self.assertEqual(result.returncode, 0, result.stderr)
                self.assertIn("### " + heading, result.stdout)
                self.assertIn(marker, result.stdout)

    def test_shared_discovery_catalog_lists_and_loads_cases(self):
        result = run_helper(str(SHARED_DISCOVERY_CATALOG))
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertIn("### Monitor excluded", result.stdout)
        self.assertNotIn("**Root cause:**", result.stdout)

        result = run_helper(str(SHARED_DISCOVERY_CATALOG), "Monitor excluded")
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertIn("### Monitor excluded", result.stdout)
        self.assertIn("serviceMonitorNamespaceSelector", result.stdout)


if __name__ == "__main__":
    unittest.main()
