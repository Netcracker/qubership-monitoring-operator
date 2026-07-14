# Monitoring MCP Servers

This document describes which Model Context Protocol (MCP) servers can be used
with monitoring components managed by the operator, what we recommend, and
where to find detailed setup instructions.

## Table of Contents

- [Monitoring MCP Servers](#monitoring-mcp-servers)
  - [Table of Contents](#table-of-contents)
  - [Recommended MCP Servers](#recommended-mcp-servers)
  - [How to Choose](#how-to-choose)
  - [Deployment Recommendations](#deployment-recommendations)
  - [Deploy with Monitoring Helm Chart](#deploy-with-monitoring-helm-chart)
  - [Detailed Setup Guides](#detailed-setup-guides)
  - [How to Use MCP Servers](#how-to-use-mcp-servers)
  - [Security Notes](#security-notes)

## Recommended MCP Servers

For the monitoring stack, we recommend two MCP servers:

| Monitoring component | MCP server | Recommendation |
| --- | --- | --- |
| Grafana | `mcp-grafana` | Recommended as the main MCP entrypoint for dashboards, datasources, alerts, folders, annotations, and querying metrics or logs through Grafana datasources. |
| VictoriaMetrics | `mcp-victoriametrics` | Recommended for direct VictoriaMetrics API access, MetricsQL or PromQL queries, labels, series, cardinality and query diagnostics, rules, alerts, and embedded VictoriaMetrics documentation. |

Use `mcp-grafana` when the user needs a broad observability interface through
Grafana. Use `mcp-victoriametrics` when the user needs direct access to
VictoriaMetrics or wants to diagnose VictoriaMetrics-specific behavior.

## How to Choose

| Use case | Recommended MCP |
| --- | --- |
| Explore dashboards, folders, panels, and datasource configuration. | `mcp-grafana` |
| Query Prometheus-compatible metrics through the Grafana datasource layer. | `mcp-grafana` |
| Query Loki logs through Grafana datasources. | `mcp-grafana` |
| Inspect Grafana alerting, contact points, annotations, snapshots, or plugins. | `mcp-grafana` |
| Query VictoriaMetrics directly with MetricsQL or PromQL. | `mcp-victoriametrics` |
| List VictoriaMetrics metrics, labels, label values, and series. | `mcp-victoriametrics` |
| Investigate VictoriaMetrics cardinality, active queries, top queries, rules, or alerts. | `mcp-victoriametrics` |
| Search VictoriaMetrics documentation offline from the MCP client. | `mcp-victoriametrics` |

In most cases, start with `mcp-grafana` because it matches the way users
usually navigate the monitoring stack: dashboards, datasources, alerts, and
queries from one place. Add `mcp-victoriametrics` when direct VictoriaMetrics
debugging or lower-level API access is required.

## Deployment Recommendations

There are two common deployment models.

| Deployment model | When to use it | Notes |
| --- | --- | --- |
| Local MCP server | Local development, debugging, or personal assistant setup. | The MCP server runs on the user's workstation. The monitored component must be reachable from that workstation through HTTPRoute, Ingress, port-forward, or another accessible endpoint. |
| In-cluster MCP server | Shared team usage or a stable endpoint for multiple MCP clients. | The MCP server runs in the Kubernetes cluster, usually installed with Helm and exposed through HTTPRoute or Ingress. Protect the exposed MCP endpoint at the gateway or ingress layer. |

For local `stdio` mode, the MCP client starts the MCP server process itself.
This is common for local binaries and `uvx`-based `mcp-grafana` usage.

For `http` or `streamable-http` mode, the MCP server runs separately and the
MCP client connects to a URL such as:

```text
http://localhost:8000/
https://mcp-grafana.example.com/
https://mcp-vm.example.com/mcp
```

## Deploy with Monitoring Helm Chart

The monitoring Helm chart can deploy MCP servers together with the monitoring
stack.

Enable `mcp-grafana` with:

```yaml
grafana:
  install: true
  mcp:
    install: true
    grafanaUrl: "http://grafana-service.monitoring.svc:3000"
    existingSecret:
      name: grafana-mcp-token
      key: token
    disableWrite: true
    allowedHosts:
      - mcp-grafana.example.com
    httpRoute:
      install: true
      hostnames:
        - mcp-grafana.example.com
      parentRefs:
        - name: gateway
          namespace: istio-system
```

Enable `mcp-victoriametrics` with:

```yaml
victoriametrics:
  vmSingle:
    install: true
  mcp:
    install: true
    vm:
      entrypoint: "http://vmsingle-k8s.monitoring.svc:8428"
      type: single
    httpRoute:
      install: true
      hostnames:
        - mcp-vm.example.com
      parentRefs:
        - name: gateway
          namespace: istio-system
```

Use Ingress instead of HTTPRoute by setting `*.mcp.ingress.install: true` and
providing the required hosts, paths, annotations, and TLS settings.

## Detailed Setup Guides

Use the detailed guides for component-specific installation and configuration:

* [Installing mcp-grafana](mcp-grafana.md)
* [Installing mcp-victoriametrics](mcp-victoriametrics.md)

The detailed guides include:

* prerequisites for the operator-managed monitoring stack;
* endpoint selection for local and in-cluster MCP deployments;
* authentication and TLS configuration;
* installation examples;
* MCP client configuration examples;
* available tools and smoke-test prompts.

## How to Use MCP Servers

After an MCP server is configured in the MCP client, use natural-language
requests that map to the component capabilities.

Example prompts for `mcp-grafana`:

```text
List Grafana datasources and check their health.
```

```text
Find dashboards related to Kubernetes nodes and show which datasource they use.
```

```text
Query the default Prometheus datasource for CPU usage during the last hour.
```

Example prompts for `mcp-victoriametrics`:

```text
List available VictoriaMetrics metrics matching kube_pod.*.
```

```text
Run the query up and group the result by job.
```

```text
Show VictoriaMetrics top queries and TSDB cardinality statistics.
```

```text
Search VictoriaMetrics documentation for MetricsQL rollup functions.
```

## Security Notes

Treat MCP servers as privileged integration points. They can expose monitoring
metadata, queries, dashboards, labels, alerts, and sometimes write operations
depending on configuration.

Follow these rules:

* Prefer read-only operation unless write access is explicitly required.
* Use Grafana service account tokens with the minimum required permissions.
* Protect MCP servers exposed through HTTPRoute, Ingress, or API Gateway with
  gateway-level authentication such as OAuth2/OIDC, Basic Auth, mTLS, IP
  allowlists, or private network access.
* Do not expose local-only MCP servers to shared networks.
* For self-signed certificates, trust the issuing CA on the machine or
  container that initiates the TLS connection.
* Rotate tokens if they were shared outside the intended local or test
  environment.
