# Examples

## Which example see at the first time

This topic should help you to chose which example you need to see firstly.

### ServiceMonitor

Custom resource which allow declaratively specify how for collect metrics from services.

ServiceMonitor use a label selector to select Service and their Endpoint objects. It means that for collect metrics with
using ServiceMonitor config your service deploy should have Service and any deploy config.

The Operator automatically generates Prometheus/api/v1/label/<label_name>/values scrape configuration based on the definition.

### PodMonitor

Custom resource which allow declaratively specify how for collect metrics from pods. The Operator automatically
generates Prometheus/api/v1/label/<label_name>/values scrape configuration based on the definition.

### PrometheusRule

Custom resource which allow declaratively specify rules for AlertManager for alerting.

### GrafanaDashboard

Custom resource which allow specifying grafana dashboard json with name and list of plugins (or link to download this
dashboard). The Grafana Operator automatically import all found dashboards to Grafana.

### GrafanaDataSource

Custom resource which allow specifying grafana datasources.
The Grafana Operator automatically import all found datasources _in current namespace_ to Grafana.

## Simple examples

* [ServiceMonitor](custom-resources/service-monitor/README.md)
* [PodMonitor](custom-resources/pod-monitor/README.md)
* [PrometheusRule](custom-resources/prometheus-rule/README.md)
* [GrafanaDashboard](custom-resources/grafana-dashboard/README.md)
* [GrafanaDataSource](custom-resources/grafana-datasource/README.md)

## Complex examples

* [Collect metrics and add grafana dashboard](custom-resources/service-with-dashboard/README.md)
* [Collect metrics and add alarms](custom-resources/service-with-alarms/README.md)
* [Collect metrics from custom endpoints](custom-resources/custom-endpoint/README.md)
* [Collect metrics from HTTPS endpoints](custom-resources/service-with-tls/README.md)



