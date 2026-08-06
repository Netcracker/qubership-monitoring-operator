import sys
import types
import unittest
from unittest import mock


# PlatformLibrary is provided by the integration-tests base image, not requirements.txt.
_fake_platform = types.ModuleType("PlatformLibrary")


class _FakePlatformLibrary:
    def __init__(self, *args, **kwargs):
        pass


_fake_platform.PlatformLibrary = _FakePlatformLibrary
sys.modules["PlatformLibrary"] = _fake_platform

import CheckJsonObject as cjo  # noqa: E402


class _LBEntry:
    def __init__(self, hostname=None, ip=None):
        self.hostname = hostname
        self.ip = ip


class _LoadBalancer:
    def __init__(self, ingress=None):
        self.ingress = ingress


class _Status:
    def __init__(self, load_balancer=None):
        self.load_balancer = load_balancer


class _Ingress:
    def __init__(self, status=None):
        self.status = status


class TestIsWildcardHost(unittest.TestCase):
    def test_bare_wildcard_host(self):
        self.assertTrue(cjo._is_wildcard_host("*.us-east-1.elb.amazonaws.com"))

    def test_http_url_with_wildcard(self):
        self.assertTrue(cjo._is_wildcard_host("http://*.us-east-1.elb.amazonaws.com"))

    def test_https_url_with_wildcard(self):
        self.assertTrue(cjo._is_wildcard_host("https://*.example.com"))

    def test_non_wildcard_host(self):
        self.assertFalse(cjo._is_wildcard_host("grafana.example.com"))

    def test_non_wildcard_url(self):
        self.assertFalse(cjo._is_wildcard_host("http://grafana.example.com"))

    def test_empty(self):
        self.assertFalse(cjo._is_wildcard_host(""))
        self.assertFalse(cjo._is_wildcard_host(None))

    def test_star_in_path_is_not_host_wildcard(self):
        self.assertFalse(cjo._is_wildcard_host("http://grafana.example.com/*/api"))


class TestResolveIngressUrl(unittest.TestCase):
    def test_non_wildcard_passthrough(self):
        k8s = mock.Mock()
        self.assertEqual(
            cjo._resolve_ingress_url(
                "http://grafana.example.com", k8s, "ns-grafana", "ns"
            ),
            "http://grafana.example.com",
        )
        k8s.get_ingress.assert_not_called()

    def test_wildcard_http_preserves_scheme_and_uses_hostname(self):
        k8s = mock.Mock()
        k8s.get_ingress.return_value = _Ingress(
            _Status(_LoadBalancer([_LBEntry(hostname="abc-123.elb.amazonaws.com")]))
        )
        self.assertEqual(
            cjo._resolve_ingress_url(
                "http://*.us-east-1.elb.amazonaws.com",
                k8s,
                "monitoring-grafana",
                "monitoring",
            ),
            "http://abc-123.elb.amazonaws.com",
        )
        k8s.get_ingress.assert_called_once_with("monitoring-grafana", "monitoring")

    def test_wildcard_uses_ip_when_hostname_absent(self):
        k8s = mock.Mock()
        k8s.get_ingress.return_value = _Ingress(
            _Status(_LoadBalancer([_LBEntry(ip="10.0.0.5")]))
        )
        self.assertEqual(
            cjo._resolve_ingress_url(
                "*.example.com", k8s, "ns-grafana", "ns"
            ),
            "10.0.0.5",
        )

    def test_wildcard_empty_status_returns_empty(self):
        k8s = mock.Mock()
        k8s.get_ingress.return_value = _Ingress(_Status(_LoadBalancer([])))
        self.assertEqual(
            cjo._resolve_ingress_url(
                "http://*.example.com", k8s, "ns-grafana", "ns"
            ),
            "",
        )

    def test_wildcard_missing_status_returns_empty(self):
        k8s = mock.Mock()
        k8s.get_ingress.return_value = _Ingress(None)
        self.assertEqual(
            cjo._resolve_ingress_url("*.example.com", k8s, "ns-grafana", "ns"),
            "",
        )


class TestResolveIngressHost(unittest.TestCase):
    def test_non_empty_non_wildcard_unchanged(self):
        with mock.patch.object(cjo, "PlatformLibrary") as pl_cls:
            result = cjo.resolve_ingress_host(
                "grafana.example.com", "ns-grafana", "ns"
            )
            self.assertEqual(result, "grafana.example.com")
            pl_cls.assert_called_once_with()

    def test_empty_host_resolves_from_lb(self):
        k8s = mock.Mock()
        k8s.get_ingress.return_value = _Ingress(
            _Status(_LoadBalancer([_LBEntry(hostname="lb.example.com")]))
        )
        with mock.patch.object(cjo, "PlatformLibrary", return_value=k8s):
            self.assertEqual(
                cjo.resolve_ingress_host("", "ns-grafana", "ns"),
                "lb.example.com",
            )

    def test_wildcard_host_resolves_from_lb(self):
        k8s = mock.Mock()
        k8s.get_ingress.return_value = _Ingress(
            _Status(_LoadBalancer([_LBEntry(hostname="lb.example.com")]))
        )
        with mock.patch.object(cjo, "PlatformLibrary", return_value=k8s):
            self.assertEqual(
                cjo.resolve_ingress_host(
                    "*.us-east-1.elb.amazonaws.com", "ns-grafana", "ns"
                ),
                "lb.example.com",
            )

    def test_empty_host_without_lb_raises(self):
        k8s = mock.Mock()
        k8s.get_ingress.return_value = _Ingress(_Status(_LoadBalancer([])))
        with mock.patch.object(cjo, "PlatformLibrary", return_value=k8s):
            with self.assertRaises(AssertionError):
                cjo.resolve_ingress_host("", "ns-grafana", "ns")


if __name__ == "__main__":
    unittest.main()
