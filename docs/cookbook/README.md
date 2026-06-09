For more information, refer to the following:

* [Monitoring Overview](overview.md)
* Architecture Overview
* Best Practices
* Features
* Clients

# Business Value

Monitoring is an important part of the solution, which allows you collect information about various aspects of the system,
the application's behavior and displays it (as in raw view, as in as various aggregations) to the user.

In Monitoring stack, we are using the following main tools to collect, analyze, visualize, and check metrics:

* `Prometheus` - Prometheus collects and stores its metrics as time series data, i.e., metrics information
  is stored with the timestamp at which it was recorded, alongside optional key-value pairs called labels.
* `Grafana` - Grafana is an open source visualization and analytics software. It allows you to query, visualize,
  alert on, and explore your metrics no matter where they are stored. In plain language, it provides
  you with tools to turn your time-series database (TSDB) data into beautiful graphs and visualizations.

Prometheus's main features are:

* A multi-dimensional data model with time series data identified by metric name and key/value pairs.
* PromQL, a flexible query language to leverage this dimensionality.
* No reliance on distributed storage; single server nodes are autonomous.
* Time series collection happens through a pull model over HTTP.
* Pushing time series is supported through an intermediary gateway.
* Targets are discovered through service discovery or static configuration.
* Multiple modes of graphing and dashboarding support.

# Links

The important links are described in the following table.

<!-- markdownlint-disable line-length -->
| Name                                       | URL                                                                                                                                                                                                                |
| ------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| Prometheus Official Documentation          | [https://prometheus.io/docs/introduction/overview/](https://prometheus.io/docs/introduction/overview/)                                                                                                             |
| Prometheus-operator Official Documentation | [https://github.com/prometheus-operator/prometheus-operator/tree/master/Documentation](https://github.com/prometheus-operator/prometheus-operator/tree/master/Documentation)                                       |
| Grafana Official Documentation             | [https://grafana.com/docs/grafana/latest/](https://grafana.com/docs/grafana/latest/)                                                                                                                               |
| Grafana-operator Official Documentation    | [https://github.com/grafana/grafana-operator/tree/v4.10.1/documentation](https://github.com/grafana/grafana-operator/tree/v4.10.1/documentation)                                                   |
| OpenMetrics / OpenTelemetry specification  | [https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/metrics/README.md](https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/metrics/README.md) |
<!-- markdownlint-enable line-length -->
