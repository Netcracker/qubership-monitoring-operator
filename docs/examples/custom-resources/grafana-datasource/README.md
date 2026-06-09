# Grafana DataSource example

* [Grafana DataSource example](#grafana-datasource-example)
  * [Overview](#overview)
  * [DataSource](#datasource)
  * [Files](#files)
  * [How to apply the simple example](#how-to-apply-the-simple-example)
  * [Links](#links)

**[Back](../../README.md)**

This example show how to add Grafana DataSource to collect data from different sources.

## Overview

Grafana uses datasources to retrieve data from various sources and display it on dashboards.
Every datasource has a type to specify which source of data
will be used (Prometheus, InfluxDB, Graphite, Jaeger, etc.).
Some the most popular types are available by default, but to support the rest,
the corresponding plugins must be installed in Grafana.

Datasource has information about connection to source of data and some support information
(credentials, time intervals between scraping metrics, HTTP headers, etc.).

You can add datasource via Grafana UI (see more
[in official documentation](https://grafana.com/docs/grafana/latest/datasources/add-a-data-source/)),
but the datasources created in this way will be deleted as soon as the Grafana instance is rebooted.

This document describes how to add Grafana datasource as Custom Resource for the `grafana-operator`.

## DataSource

```yaml
...
spec:
  name: simple-datasource-example.yaml
  datasources:
  - name: Prometheus datasource
    type: prometheus
    access: proxy
    ...
```

It means that datasource with name `Prometheus datasource` and type `prometheus` will be added to Grafana.

One DataSource CR can handle multiple datasources with unique names.

Required fields for every datasource: `name`, `type`, `access`. Full list of parameters is unique for each
type of datasource.

## Files

* [Simple DataSource example](simple-datasource-example.yaml)
* [Full DataSource example](full-datasource-example.yaml)

See more examples in the
[`grafana-operator` repository](https://github.com/grafana/grafana-operator/tree/v4.10.1/deploy/examples/datasources).

## How to apply the simple example

Kubernetes:

```bash
kubectl apply -f simple-datasource-example.yaml
```

OpenShift:

```bash
oc apply -f simple-datasource-example.yaml
```

## Links

* Grafana official documentation
  * [Configuration of datasource example](https://grafana.com/docs/grafana/latest/administration/provisioning/#data-sources)
  * [Add a data source (via UI)](https://grafana.com/docs/grafana/latest/datasources/add-a-data-source)
* Grafana-operator
  * [Working with data sources](https://github.com/grafana/grafana-operator/blob/v4/documentation/datasources.md)
  * [Plugins](https://github.com/grafana/grafana-operator/blob/v4/documentation/plugins.md)
  * [Examples of data sources](https://github.com/grafana/grafana-operator/tree/v4.10.1/deploy/examples/datasources)
