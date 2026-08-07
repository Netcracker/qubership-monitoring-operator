# Grafana datasource example

* [Grafana DataSource example](#grafana-datasource-example)
  * [Overview](#overview)
  * [DataSource](#datasource)
  * [Migrate v4 resources](#migrate-v4-resources)
  * [Files](#files)
  * [How to apply the simple example](#how-to-apply-the-simple-example)
  * [Links](#links)

**[Back](../../README.md)**

This example shows how to add a Grafana datasource that collects data from another source.

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
  instanceSelector:
    matchLabels:
      app.kubernetes.io/component: grafana
  datasource:
    name: Prometheus datasource
    type: prometheus
    access: proxy
    ...
```

It means that datasource with name `Prometheus datasource` and type `prometheus` will be added to Grafana.

Each GrafanaDatasource CR defines one datasource. Create a separate CR for each additional datasource.

Required fields for every datasource: `name`, `type`, `access`. Full list of parameters is unique for each
type of datasource.

## Migrate v4 resources

Grafana Operator v5 does not reconcile the former `integreatly.org/v1alpha1` `GrafanaDataSource` resource. Before
upgrading, inventory those resources:

```bash
kubectl get grafanadatasources.integreatly.org --all-namespaces
```

For every item in the old `spec.datasources` array, create one `grafana.integreatly.org/v1beta1`
`GrafanaDatasource`. Move the item to singular `spec.datasource`, give each resource a unique `metadata.name`, and add
an `instanceSelector` that matches the target Grafana instance. Apply and verify the replacement resources before
removing the old resources or CRD.

The [Grafana operator converter](https://github.com/Netcracker/qubership-grafana-operator-converter) migrates legacy
Grafana dashboards only; it does not migrate datasources.

## Files

* [Simple DataSource example](simple-datasource-example.yaml)
* [Full DataSource example](full-datasource-example.yaml)

See more examples in the
[`grafana-operator` documentation](https://grafana.github.io/grafana-operator/docs/examples/datasource/).

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
* Grafana Operator
  * [GrafanaDatasource examples](https://grafana.github.io/grafana-operator/docs/examples/datasource/)
  * [Plugin management](https://grafana.github.io/grafana-operator/docs/examples/datasource/plugins/readme/)
