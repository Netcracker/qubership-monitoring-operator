# ssl-exporter-metrics

This document describes the metrics list and how to collect them from ssl-exporter.

## How to Collect

ssl-exporter exposes process metrics on port `9219` at `/metrics`, and provides a `/probe` endpoint to actively check certificates using `target` and `module` parameters.

The recommended way to collect metrics is via a ServiceMonitor created by the chart, see the [installation guide](../../installation/components/exporters/ssl-exporter.md). Example of a standalone ServiceMonitor:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: ssl-exporter-example
  labels:
    app.kubernetes.io/component: monitoring
spec:
  endpoints:
    - port: http
      scheme: http
      path: /probe
      interval: 60s
      params:
        target:
          - google.com:443
        module:
          - https-external
      metricRelabelings:
        - sourceLabels: [instance]
          targetLabel: instance
          replacement: google.com:443
        - sourceLabels: [target]
          targetLabel: target
          replacement: https-external-google
  jobLabel: k8s-app
  selector:
    matchLabels:
      app.kubernetes.io/name: ssl-exporter
```

Check metrics manually:

```bash
# Exporter process metrics endpoint
curl -s http://<ssl_exporter_service>:9219/metrics | head -n 40

# Run a probe for external HTTPS (returns probe metrics)
curl -G -s "http://<ssl_exporter_service>:9219/probe" \
  --data-urlencode target=google.com:443 \
  --data-urlencode module=https-external | head -n 80
```

## Metrics List

Below are typical metrics emitted by ssl-exporter. The exact set depends on the selected module and target.

```prometheus
# Time until certificate expiry (seconds)
# The lower the value, the closer to expiry
# TYPE ssl_cert_not_after gauge
ssl_cert_not_after{target="google.com:443",module="https-external"} 2.592e+06

# Certificate valid from timestamp (seconds since epoch)
# TYPE ssl_cert_not_before gauge
ssl_cert_not_before{target="google.com:443",module="https-external"} 1.700e+09

# Certificate age (seconds)
# TYPE ssl_cert_age_seconds gauge
ssl_cert_age_seconds{target="google.com:443",module="https-external"} 1.234e+07

# Days until expiry (if exported by the module)
# TYPE ssl_cert_days_until_expiry gauge
ssl_cert_days_until_expiry{target="google.com:443",module="https-external"} 30

# Certificate serial number exposed as label (value set to 1)
# TYPE ssl_cert_serial gauge
ssl_cert_serial{target="google.com:443",module="https-external",serial="03:AB:CD:..."} 1

# Chain validation result (0 — ok, 1 — error), if exported by the module
# TYPE ssl_cert_validation_result gauge
ssl_cert_validation_result{target="google.com:443",module="https-external"} 0

# Issuer/subject information as labels (value set to 1)
# TYPE ssl_cert_info gauge
ssl_cert_info{target="google.com:443",module="https-external",issuer_cn="GTS CA 1C3",subject_cn="*.google.com"} 1

# Exporter process metrics (examples)
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
# ...
```

Notes:
- For file and kubeconfig targets, use the corresponding modules (`file`, `kubeconfig`) and mount paths via `additionalHostPathVolumes`.
- To read certificates from Kubernetes secrets, use the `kubernetes` module and enable RBAC (`rbac.create: true`).
- For self-signed certificates, use the `https-selfsigned` module (`insecure_skip_verify: true` by default).

## Alerting (recommendations)

Example rules:

```yaml
- alert: SSLCertExpiringSoon
  expr: ssl_cert_not_after - time() < 86400 * 7
  for: 10m
  labels:
    severity: warning
  annotations:
    summary: "SSL certificate for {{ $labels.instance }} expires soon"
    description: "The SSL certificate for {{ $labels.instance }} will expire in less than 7 days."

- alert: SSLCertExpired
  expr: ssl_cert_not_after - time() < 0
  for: 10m
  labels:
    severity: critical
  annotations:
    summary: "SSL certificate for {{ $labels.instance }} has expired"
    description: "The SSL certificate for {{ $labels.instance }} has expired."
``` 
