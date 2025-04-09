# Monitoring Operator

The `monitoring-operator` is an operator (Kubernetes native application) which watches and manages following components:

* `victoriametrics-operator`
* `vmSingle`
* `vmAgent`
* `vmAlertManager`
* `vmAlert`
* `vmAuth`
* `vmUser`
* `grafana-operator`
* `grafana`
* `kube-state-metrics`
* `node-exporter`
* `configuration-streamer`
* `graphite-remote-adapter`

And provides ability to collect metrics via prometheus from `kube-state-metrics`, `node-exporter` and different
applications and store this metrics to remote LTS.

## Documentation

### Public documents

This section contains documents of directories which may be provide to customers.

Guides:

* [Installation](docs/installation.md)
* [Maintenance](docs/maintenance.md)
* [Troubleshooting](docs/troubleshooting.md)

Documents described Custom Resources Definitions (CRDs):

* [PlatformMonitoring](docs/api/platform-monitoring.md)
* [PrometheusAdapter](docs/api/prometheus-adapter.md)
* [CustomScaleMetricRule](docs/api/custom-scale-metric-rule.md)

Documents described metrics, alerts, dashboards which deploy with monitoring out-of-box (OOB):

* [Metrics](docs/defaults/metrics.md)
* [Alerts](docs/defaults/alerts.md)
* [Dashboards](docs/defaults/dashboards)

Examples:

* [Custom resources](docs/examples/custom-resources)

### Internal documents

This section contains documents or directories which should not be provided to customers and should be used only
within the company.

### Images

All images are stored in [docs/images](docs/images).

And all the sources of these images are stored in following directories:

* [docs/sources/draw.io](docs/sources/draw.io) - for diagrams which are made using [https://draw.io](https://draw.io)
* [docs/sources/plantuml](docs/sources/plantuml) - for diagrams which are made using PlantUML syntax
