import sys
import types
import unittest
from unittest import mock
from urllib.parse import urlparse

# PlatformLibrary is provided by the integration-tests base image, not requirements.txt.
_fake_platform = types.ModuleType("PlatformLibrary")


class _FakePlatformLibrary:
    def __init__(self, *args, **kwargs):
        pass


_fake_platform.PlatformLibrary = _FakePlatformLibrary
sys.modules["PlatformLibrary"] = _fake_platform

import CheckJsonObject as cjo  # noqa: E402

WILDCARD = "*.us-east-1.elb.amazonaws.com"
MATCHING_LB = "abc-123.us-east-1.elb.amazonaws.com"
UNRELATED_LB = "other.example.com"


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


class _Rule:
    def __init__(self, host=None):
        self.host = host


class _Spec:
    def __init__(self, rules=None):
        self.rules = rules or []


class _Ingress:
    def __init__(self, status=None, spec=None):
        self.status = status
        self.spec = spec


def _ingress_with_lb(entries, spec_host=WILDCARD):
    return _Ingress(
        status=_Status(_LoadBalancer(entries)),
        spec=_Spec([_Rule(spec_host)]),
    )


class TestIsWildcardHost(unittest.TestCase):
    def test_bare_wildcard_host(self):
        self.assertTrue(cjo._is_wildcard_host(WILDCARD))

    def test_http_url_with_wildcard(self):
        self.assertTrue(cjo._is_wildcard_host(f"http://{WILDCARD}"))

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


class TestHostnameMatchesIngressWildcard(unittest.TestCase):
    def test_matching_single_label(self):
        self.assertTrue(
            cjo._hostname_matches_ingress_wildcard(MATCHING_LB, WILDCARD)
        )

    def test_rejects_multi_label_prefix(self):
        self.assertFalse(
            cjo._hostname_matches_ingress_wildcard(
                "a.b.us-east-1.elb.amazonaws.com", WILDCARD
            )
        )

    def test_rejects_unrelated_hostname(self):
        self.assertFalse(
            cjo._hostname_matches_ingress_wildcard(UNRELATED_LB, WILDCARD)
        )

    def test_rejects_ip(self):
        self.assertFalse(
            cjo._hostname_matches_ingress_wildcard("10.96.254.104", WILDCARD)
        )

    def test_rejects_bare_suffix(self):
        self.assertFalse(
            cjo._hostname_matches_ingress_wildcard(
                "us-east-1.elb.amazonaws.com", WILDCARD
            )
        )


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

    def test_matching_aws_hostname_preserves_https_scheme(self):
        k8s = mock.Mock()
        k8s.get_ingress.return_value = _ingress_with_lb(
            [_LBEntry(hostname=MATCHING_LB)]
        )
        self.assertEqual(
            cjo._resolve_ingress_url(
                f"https://{WILDCARD}",
                k8s,
                "monitoring-grafana",
                "monitoring",
            ),
            f"https://{MATCHING_LB}",
        )

    def test_skips_unrelated_hostname_and_selects_matching(self):
        k8s = mock.Mock()
        k8s.get_ingress.return_value = _ingress_with_lb(
            [
                _LBEntry(hostname=UNRELATED_LB),
                _LBEntry(ip="10.0.0.5"),
                _LBEntry(hostname=MATCHING_LB),
            ]
        )
        self.assertEqual(
            cjo._resolve_ingress_url(WILDCARD, k8s, "ns-grafana", "ns"),
            MATCHING_LB,
        )

    def test_ip_only_is_not_used_as_authority(self):
        k8s = mock.Mock()
        k8s.get_ingress.return_value = _ingress_with_lb(
            [_LBEntry(ip="10.0.0.5")]
        )
        with self.assertRaises(AssertionError):
            cjo._resolve_ingress_url(WILDCARD, k8s, "ns-grafana", "ns")

    def test_unrelated_hostname_only_raises(self):
        k8s = mock.Mock()
        k8s.get_ingress.return_value = _ingress_with_lb(
            [_LBEntry(hostname=UNRELATED_LB)]
        )
        with self.assertRaises(AssertionError):
            cjo._resolve_ingress_url(f"http://{WILDCARD}", k8s, "ns-grafana", "ns")

    def test_wildcard_empty_status_raises(self):
        k8s = mock.Mock()
        k8s.get_ingress.return_value = _ingress_with_lb([])
        with self.assertRaises(AssertionError):
            cjo._resolve_ingress_url(f"http://{WILDCARD}", k8s, "ns-grafana", "ns")

    def test_api_error_is_propagated(self):
        k8s = mock.Mock()
        k8s.get_ingress.side_effect = RuntimeError("forbidden")
        with self.assertRaises(RuntimeError):
            cjo._resolve_ingress_url(WILDCARD, k8s, "ns-grafana", "ns")


class TestResolveIngressHost(unittest.TestCase):
    def test_non_empty_non_wildcard_unchanged(self):
        with mock.patch.object(cjo, "PlatformLibrary") as pl_cls:
            result = cjo.resolve_ingress_host(
                "grafana.example.com", "ns-grafana", "ns"
            )
            self.assertEqual(result, "grafana.example.com")
            pl_cls.assert_called_once_with()

    def test_empty_host_uses_spec_wildcard_pattern(self):
        k8s = mock.Mock()
        k8s.get_ingress.return_value = _ingress_with_lb(
            [_LBEntry(hostname=MATCHING_LB)],
            spec_host=WILDCARD,
        )
        with mock.patch.object(cjo, "PlatformLibrary", return_value=k8s):
            self.assertEqual(
                cjo.resolve_ingress_host("", "ns-grafana", "ns"),
                MATCHING_LB,
            )

    def test_wildcard_host_resolves_matching_lb(self):
        k8s = mock.Mock()
        k8s.get_ingress.return_value = _ingress_with_lb(
            [_LBEntry(hostname=MATCHING_LB)]
        )
        with mock.patch.object(cjo, "PlatformLibrary", return_value=k8s):
            self.assertEqual(
                cjo.resolve_ingress_host(WILDCARD, "ns-grafana", "ns"),
                MATCHING_LB,
            )

    def test_no_matching_lb_raises(self):
        k8s = mock.Mock()
        k8s.get_ingress.return_value = _ingress_with_lb(
            [_LBEntry(ip="10.96.254.104")]
        )
        with mock.patch.object(cjo, "PlatformLibrary", return_value=k8s):
            with self.assertRaises(AssertionError):
                cjo.resolve_ingress_host(WILDCARD, "ns-grafana", "ns")

    def test_api_error_is_propagated(self):
        k8s = mock.Mock()
        k8s.get_ingress.side_effect = RuntimeError("timeout")
        with mock.patch.object(cjo, "PlatformLibrary", return_value=k8s):
            with self.assertRaises(RuntimeError):
                cjo.resolve_ingress_host(WILDCARD, "ns-grafana", "ns")


class TestRequestAuthority(unittest.TestCase):
    """Resolved URL is used as real HTTPS Host and TLS SNI authority."""

    @staticmethod
    def _make_self_signed_cert(hostname):
        import os
        import subprocess
        import tempfile

        tmp = tempfile.mkdtemp()
        cert_path = os.path.join(tmp, "cert.pem")
        key_path = os.path.join(tmp, "key.pem")
        subprocess.check_call(
            [
                "openssl",
                "req",
                "-x509",
                "-newkey",
                "rsa:2048",
                "-keyout",
                key_path,
                "-out",
                cert_path,
                "-days",
                "1",
                "-nodes",
                "-subj",
                f"/CN={hostname}",
                "-addext",
                f"subjectAltName=DNS:{hostname}",
            ],
            stdout=subprocess.DEVNULL,
            stderr=subprocess.DEVNULL,
        )
        return cert_path, key_path

    def test_matching_hostname_is_https_host_and_sni(self):
        import http.client
        import socket
        import ssl
        import threading
        from http.server import BaseHTTPRequestHandler, HTTPServer

        k8s = mock.Mock()
        k8s.get_ingress.return_value = _ingress_with_lb(
            [
                _LBEntry(ip="10.96.254.104"),
                _LBEntry(hostname=UNRELATED_LB),
                _LBEntry(hostname=MATCHING_LB),
            ]
        )
        resolved = cjo._resolve_ingress_url(
            f"https://{WILDCARD}", k8s, "ns-grafana", "ns"
        )
        authority = urlparse(resolved).hostname
        self.assertEqual(authority, MATCHING_LB)

        captured = {"host": None, "sni": None, "path": None}
        cert_path, key_path = self._make_self_signed_cert(MATCHING_LB)

        class Handler(BaseHTTPRequestHandler):
            def do_GET(self):
                captured["host"] = self.headers.get("Host")
                captured["path"] = self.path
                self.send_response(200)
                self.end_headers()
                self.wfile.write(b"ok")

            def log_message(self, format, *args):
                return

        server = HTTPServer(("127.0.0.1", 0), Handler)
        ctx = ssl.SSLContext(ssl.PROTOCOL_TLS_SERVER)
        ctx.load_cert_chain(certfile=cert_path, keyfile=key_path)

        def sni_callback(sock, server_name, _ctx):
            captured["sni"] = server_name
            return None

        ctx.sni_callback = sni_callback
        server.socket = ctx.wrap_socket(server.socket, server_side=True)
        port = server.server_address[1]
        thread = threading.Thread(target=server.serve_forever, daemon=True)
        thread.start()
        try:
            client_ctx = ssl.create_default_context()
            client_ctx.check_hostname = False
            client_ctx.verify_mode = ssl.CERT_NONE

            class LocalHTTPSConnection(http.client.HTTPSConnection):
                def connect(self):
                    sock = socket.create_connection(
                        ("127.0.0.1", self.port), self.timeout
                    )
                    self.sock = self._context.wrap_socket(
                        sock, server_hostname=self.host
                    )

            conn = LocalHTTPSConnection(
                MATCHING_LB, port=port, context=client_ctx, timeout=5
            )
            conn.request("GET", "/login")
            response = conn.getresponse()
            self.assertEqual(response.status, 200)
            self.assertEqual(response.read(), b"ok")
            conn.close()
        finally:
            server.shutdown()
            server.server_close()

        self.assertEqual(captured["path"], "/login")
        self.assertEqual(captured["host"].split(":", 1)[0], MATCHING_LB)
        self.assertEqual(captured["sni"], MATCHING_LB)

    def test_ip_is_never_request_authority(self):
        k8s = mock.Mock()
        k8s.get_ingress.return_value = _ingress_with_lb(
            [_LBEntry(ip="10.96.254.104")]
        )
        with self.assertRaises(AssertionError):
            cjo._resolve_ingress_url(f"https://{WILDCARD}", k8s, "ns-grafana", "ns")

    def test_unrelated_hostname_is_never_request_authority(self):
        k8s = mock.Mock()
        k8s.get_ingress.return_value = _ingress_with_lb(
            [_LBEntry(hostname=UNRELATED_LB)]
        )
        with self.assertRaises(AssertionError):
            cjo._resolve_ingress_url(f"https://{WILDCARD}", k8s, "ns-grafana", "ns")


if __name__ == "__main__":
    unittest.main()
