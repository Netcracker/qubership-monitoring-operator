# VictoriaMetrics MCP

The `victoriametrics.mcp` section deploys an in-cluster
[`mcp-victoriametrics`](https://github.com/VictoriaMetrics/mcp-victoriametrics)
server together with the VictoriaMetrics stack.

The MCP server is installed only when:

* `victoriametrics.mcp.install=true`; and
* `victoriametrics.vmOperator.install=true`; and
* the backend matching `victoriametrics.mcp.vm.type` is installed: `vmSingle`
  for `single`, or `vmCluster` for `cluster`.

Standalone MCP deployment is not supported by this umbrella chart. An explicit
`victoriametrics.mcp.vm.entrypoint` overrides the connection URL, for example to
use VMAuth, but it does not remove the matching backend installation requirement.

The MCP server does not provide its own inbound authentication for the HTTP
endpoint. If the MCP server is published through HTTPRoute or Ingress, protect
it on the gateway, ingress, or service mesh layer.

## Parameters

All parameters listed below are configured under the `victoriametrics.mcp` level.

<!-- markdownlint-disable line-length -->
| Parameter                   | Description                                                                                                                                                                                                                                                                                                  | Default                                               |
| --------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ----------------------------------------------------- |
| `install`                   | Enable deployment of the VictoriaMetrics MCP server.                                                                                                                                                                                                                                                         | `false`                                               |
| `name`                      | Base name for MCP resources.                                                                                                                                                                                                                                                                                 | `mcp-victoriametrics`                                 |
| `image`                     | Container image for `mcp-victoriametrics`.                                                                                                                                                                                                                                                                   | `ghcr.io/victoriametrics/mcp-victoriametrics:v1.20.2` |
| `imagePullPolicy`           | Image pull policy.                                                                                                                                                                                                                                                                                           | `IfNotPresent`                                        |
| `imagePullSecrets`          | Image pull secrets for the MCP pod.                                                                                                                                                                                                                                                                          | `[]`                                                  |
| `replicas`                  | Number of MCP server replicas.                                                                                                                                                                                                                                                                               | `1`                                                   |
| `mode`                      | MCP server mode for an in-cluster server: `http` for `/mcp`, or `sse` for SSE endpoints.                                                                                                                                                                                                                     | `http`                                                |
| `listenPort`                | Port where `mcp-victoriametrics` listens inside the container.                                                                                                                                                                                                                                               | `8080`                                                |
| `vm.entrypoint`             | Root VictoriaMetrics API URL reachable from the MCP pod. Empty value selects the operator-managed service matching `vm.type`: `vmsingle-k8s:8428` for `single` or `vmselect-k8s:8481` for `cluster`.                                                                                                         | `""`                                                  |
| `vm.type`                   | VictoriaMetrics instance type: `single` for VMSingle, or `cluster` for VMCluster/VMSelect.                                                                                                                                                                                                                   | `single`                                              |
| `vm.bearerToken`            | Inline bearer token sent by MCP to the configured VictoriaMetrics entrypoint, for example `vmauth` or an API Gateway. Prefer Secret-based configuration for production.                                                                                                                                      | `""`                                                  |
| `vm.bearerTokenSecret.name` | Secret that contains the bearer token for the configured entrypoint.                                                                                                                                                                                                                                         | `""`                                                  |
| `vm.bearerTokenSecret.key`  | Key in the bearer token Secret.                                                                                                                                                                                                                                                                              | `token`                                               |
| `vm.headers`                | Static headers sent by MCP to the configured entrypoint, for example `Authorization=Basic <base64>` for `vmauth` or an API Gateway.                                                                                                                                                                          | `""`                                                  |
| `vm.headersSecret.name`     | Secret that contains static headers for the configured entrypoint.                                                                                                                                                                                                                                           | `""`                                                  |
| `vm.headersSecret.key`      | Key in the headers Secret.                                                                                                                                                                                                                                                                                   | `headers`                                             |
| `passthroughHeaders`        | Incoming MCP request headers that should be forwarded to the configured VictoriaMetrics entrypoint. This does not protect the MCP endpoint itself.                                                                                                                                                           | `[]`                                                  |
| `disabledTools`             | Complete list passed through `MCP_DISABLED_TOOLS`. An empty list leaves the variable unset, so the upstream defaults remain disabled: `export`, `flags`, `metric_relabel_debug`, `downsampling_filters_debug`, `retention_filters_debug`, and `test_rules`. A non-empty list replaces that upstream default. | `[]` (upstream defaults apply)                        |
| `logLevel`                  | MCP server log level.                                                                                                                                                                                                                                                                                        | `info`                                                |
| `extraEnv`                  | Additional environment variables as key-value pairs.                                                                                                                                                                                                                                                         | `{}`                                                  |
| `resources`                 | Resource requests and limits for the MCP container.                                                                                                                                                                                                                                                          | `{}`                                                  |
| `securityContext`           | Pod security context. Empty value uses the subchart default security context.                                                                                                                                                                                                                                | `{}`                                                  |
| `nodeSelector`              | Node selector for the MCP pod.                                                                                                                                                                                                                                                                               | `{}`                                                  |
| `affinity`                  | Pod affinity rules.                                                                                                                                                                                                                                                                                          | `{}`                                                  |
| `tolerations`               | Pod tolerations.                                                                                                                                                                                                                                                                                             | `[]`                                                  |
| `volumes`                   | Extra pod volumes.                                                                                                                                                                                                                                                                                           | `[]`                                                  |
| `volumeMounts`              | Extra volume mounts for the MCP container.                                                                                                                                                                                                                                                                   | `[]`                                                  |
| `labels`                    | Extra labels added to MCP Deployment and pod template.                                                                                                                                                                                                                                                       | `{}`                                                  |
| `annotations`               | Extra annotations added to MCP Deployment.                                                                                                                                                                                                                                                                   | `{}`                                                  |
| `podAnnotations`            | Extra annotations added to MCP pods.                                                                                                                                                                                                                                                                         | `{}`                                                  |
<!-- markdownlint-enable line-length -->

Instead of configuring static credentials through `vm.bearerToken` or
`vm.headers`, an HTTP MCP client can send an authentication header with each
request and the server can forward it to the configured VictoriaMetrics
entrypoint. For example:

```yaml
victoriametrics:
  mcp:
    passthroughHeaders:
      - Authorization
```

This approach has the following limitations:

* It is supported only in `http` or `sse` mode and does not apply to `stdio`.
* The MCP client must support custom HTTP headers and send the header with each
  request.
* `passthroughHeaders` contains header names only; it does not contain or store
  their values.
* Forwarding a header does not authenticate or authorize access to the MCP
  endpoint itself.
* If an API Gateway authenticates requests to the MCP endpoint, verify that it
  preserves the header before forwarding the request to the MCP server.
* Avoid configuring both a static authorization credential and a forwarded
  `Authorization` header unless the resulting precedence has been verified.

## Service

All parameters listed below are configured under the
`victoriametrics.mcp.service` level.

<!-- markdownlint-disable line-length -->
| Parameter     | Description                                                 | Default     |
| ------------- | ----------------------------------------------------------- | ----------- |
| `type`        | Kubernetes Service type.                                    | `ClusterIP` |
| `port`        | Service port. The Service targets the named container port. | `8080`      |
| `annotations` | Extra Service annotations.                                  | `{}`        |
| `labels`      | Extra Service labels.                                       | `{}`        |
<!-- markdownlint-enable line-length -->

## ServiceAccount

All parameters listed below are configured under the
`victoriametrics.mcp.serviceAccount` level.

<!-- markdownlint-disable line-length -->
| Parameter     | Description                                                    | Default |
| ------------- | -------------------------------------------------------------- | ------- |
| `create`      | Create a ServiceAccount for the MCP pod.                       | `true`  |
| `name`        | ServiceAccount name. Empty value defaults to the MCP fullname. | `""`    |
| `annotations` | Extra ServiceAccount annotations.                              | `{}`    |
| `labels`      | Extra ServiceAccount labels.                                   | `{}`    |
<!-- markdownlint-enable line-length -->

## ServiceMonitor

`mcp-victoriametrics` exposes `/metrics` in HTTP/SSE modes. The ServiceMonitor
is optional and disabled by default.

All parameters listed below are configured under the
`victoriametrics.mcp.serviceMonitor` level.

<!-- markdownlint-disable line-length -->
| Parameter           | Description                              | Default |
| ------------------- | ---------------------------------------- | ------- |
| `install`           | Create a ServiceMonitor for MCP metrics. | `false` |
| `labels`            | Extra ServiceMonitor labels.             | `{}`    |
| `annotations`       | Extra ServiceMonitor annotations.        | `{}`    |
| `interval`          | Scrape interval.                         | `30s`   |
| `scrapeTimeout`     | Scrape timeout.                          | `10s`   |
| `metricRelabelings` | Metric relabeling rules.                 | `[]`    |
| `relabelings`       | Target relabeling rules.                 | `[]`    |
<!-- markdownlint-enable line-length -->

## HTTPRoute

All parameters listed below are configured under the
`victoriametrics.mcp.httpRoute` level.

<!-- markdownlint-disable line-length -->
| Parameter     | Description                                                                             | Default |
| ------------- | --------------------------------------------------------------------------------------- | ------- |
| `install`     | Create a Gateway API HTTPRoute for MCP.                                                 | `false` |
| `hostnames`   | HTTPRoute hostnames. Values are rendered through `tpl`.                                 | `[]`    |
| `parentRefs`  | Gateway parent references. At least one entry is required when `install=true`.          | `[]`    |
| `rules`       | Full custom HTTPRoute rules. If empty, the chart renders a default PathPrefix `/` rule. | `[]`    |
| `annotations` | Extra HTTPRoute annotations.                                                            | `{}`    |
| `labels`      | Extra HTTPRoute labels.                                                                 | `{}`    |
<!-- markdownlint-enable line-length -->

## Ingress

All parameters listed below are configured under the
`victoriametrics.mcp.ingress` level.

<!-- markdownlint-disable line-length -->
| Parameter     | Description                                                                                                                                                            | Default |
| ------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------- |
| `install`     | Create an Ingress for MCP.                                                                                                                                             | `false` |
| `className`   | IngressClass name.                                                                                                                                                     | `""`    |
| `annotations` | Extra Ingress annotations.                                                                                                                                             | `{}`    |
| `tls`         | Ingress TLS configuration.                                                                                                                                             | `[]`    |
| `hosts`       | Ingress hosts and paths. At least one entry is required when `install=true`; `host` values are rendered through `tpl`, and missing paths default to `/` with `Prefix`. | `[]`    |
<!-- markdownlint-enable line-length -->

## Example

```yaml
victoriametrics:
  vmOperator:
    install: true
  vmSingle:
    install: true
  mcp:
    install: true
    vm:
      entrypoint: https://vmauth-monitoring.example.com
      type: single
      headersSecret:
        name: vmauth-headers
        key: headers
    serviceMonitor:
      install: true
    httpRoute:
      install: true
      hostnames:
        - mcp-vm.{{ .Release.Namespace }}.example.com
      parentRefs:
        - name: monitoring-gateway
```
