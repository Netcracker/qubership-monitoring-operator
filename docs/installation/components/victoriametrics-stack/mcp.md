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

When `victoriametrics.tlsEnabled=true` and `vm.entrypoint` is empty, the chart
uses HTTPS and mounts `ca.crt` from the TLS Secret of the selected VMSingle or
VMSelect backend. The MCP container receives only the CA file and its path in
`SSL_CERT_FILE`; the TLS private key is not mounted.

The MCP server does not provide its own inbound authentication for the HTTP
endpoint. If the MCP server is published through HTTPRoute or Ingress, protect
it on the gateway, ingress, or service mesh layer.

## Parameters

All parameters listed below are configured under the `victoriametrics.mcp` level.

The chart always runs the MCP container as non-root with the runtime-default
seccomp profile, a read-only root filesystem, no Linux capabilities, and no
privilege escalation. It mounts an internal `emptyDir` limited to `100Mi` at
`/tmp`. The chart treats an empty `victoriametrics.PAAS_PLATFORM` as Kubernetes
unless the OpenShift SCC API is available. Set it explicitly to `KUBERNETES`
when API discovery is unavailable and the pod must be pinned to UID and GID
`1000`.

<!-- markdownlint-disable line-length -->
| Parameter                   | Description                                                                                                                                                                                                                                                                                                  | Default                                                     |
| --------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ----------------------------------------------------------- |
| `install`                   | Enable deployment of the VictoriaMetrics MCP server.                                                                                                                                                                                                                                                         | `false`                                                     |
| `name`                      | Base name for MCP resources.                                                                                                                                                                                                                                                                                 | `mcp-victoriametrics`                                       |
| `image`                     | Container image for `mcp-victoriametrics`.                                                                                                                                                                                                                                                                   | `ghcr.io/victoriametrics/mcp-victoriametrics:v1.20.2`       |
| `imagePullPolicy`           | Image pull policy.                                                                                                                                                                                                                                                                                           | `IfNotPresent`                                              |
| `imagePullSecrets`          | Image pull secrets for the MCP pod.                                                                                                                                                                                                                                                                          | `[]`                                                        |
| `replicas`                  | Number of MCP server replicas.                                                                                                                                                                                                                                                                               | `1`                                                         |
| `mode`                      | MCP server mode for an in-cluster server: `http` for `/mcp`, or `sse` for SSE endpoints.                                                                                                                                                                                                                     | `http`                                                      |
| `listenPort`                | Port where `mcp-victoriametrics` listens inside the container.                                                                                                                                                                                                                                               | `8080`                                                      |
| `vm.entrypoint`             | Root VictoriaMetrics API URL reachable from the MCP pod. Empty value selects the operator-managed service matching `vm.type`: `vmsingle-k8s:8428` for `single` or `vmselect-k8s:8481` for `cluster`. The generated URL uses HTTPS when `victoriametrics.tlsEnabled=true`.                                    | `""`                                                        |
| `vm.type`                   | VictoriaMetrics instance type: `single` for VMSingle, or `cluster` for VMCluster/VMSelect.                                                                                                                                                                                                                   | `single`                                                    |
| `vm.tls.useBackendCA`       | Mount the CA from the selected operator-managed backend TLS Secret when TLS is enabled and `vm.entrypoint` is empty. Disable it when the certificate is already trusted through the container system trust store.                                                                                            | `true`                                                      |
| `vm.tls.caSecret.name`      | Override the backend CA Secret name. Set this parameter explicitly when a custom `vm.entrypoint` requires a private CA.                                                                                                                                                                                      | `""`                                                        |
| `vm.tls.caSecret.key`       | Key containing the CA certificate in the configured Secret.                                                                                                                                                                                                                                                  | `ca.crt`                                                    |
| `vm.bearerToken`            | Inline bearer token passed to MCP through an environment variable. It does not meet file-based secret hardening requirements.                                                                                                                                                                                | `""`                                                        |
| `vm.bearerTokenSecret.name` | Secret that contains the bearer token. Upstream MCP has no token-file option, so the selected value is exposed to the process through an environment variable.                                                                                                                                               | `""`                                                        |
| `vm.bearerTokenSecret.key`  | Key in the bearer token Secret.                                                                                                                                                                                                                                                                              | `token`                                                     |
| `vm.headers`                | Static headers passed to MCP through an environment variable, for example `Authorization=Basic <base64>`. It does not meet file-based secret hardening requirements.                                                                                                                                         | `""`                                                        |
| `vm.headersSecret.name`     | Secret that contains static headers. Upstream MCP has no headers-file option, so the selected value is exposed to the process through an environment variable.                                                                                                                                               | `""`                                                        |
| `vm.headersSecret.key`      | Key in the headers Secret.                                                                                                                                                                                                                                                                                   | `headers`                                                   |
| `passthroughHeaders`        | Incoming MCP request headers that should be forwarded to the configured VictoriaMetrics entrypoint. This does not protect the MCP endpoint itself.                                                                                                                                                           | `[]`                                                        |
| `disabledTools`             | Complete list passed through `MCP_DISABLED_TOOLS`. An empty list leaves the variable unset, so the upstream defaults remain disabled: `export`, `flags`, `metric_relabel_debug`, `downsampling_filters_debug`, `retention_filters_debug`, and `test_rules`. A non-empty list replaces that upstream default. | `[]` (upstream defaults apply)                              |
| `logLevel`                  | MCP server log level.                                                                                                                                                                                                                                                                                        | `info`                                                      |
| `extraEnv`                  | Additional environment variables as key-value pairs.                                                                                                                                                                                                                                                         | `{}`                                                        |
| `resources`                 | Resource requests and limits for the MCP container.                                                                                                                                                                                                                                                          | `{}`                                                        |
| `securityContext`           | Additional pod security context. Mandatory non-root and seccomp settings cannot be disabled.                                                                                                                                                                                                                 | `{}`                                                        |
| `containerSecurityContext`  | Additional container security context. Mandatory read-only filesystem, capability, and privilege-escalation settings cannot be disabled.                                                                                                                                                                     | Drops all capabilities and uses a read-only root filesystem |
| `nodeSelector`              | Node selector for the MCP pod.                                                                                                                                                                                                                                                                               | `{}`                                                        |
| `affinity`                  | Pod affinity rules.                                                                                                                                                                                                                                                                                          | `{}`                                                        |
| `tolerations`               | Pod tolerations.                                                                                                                                                                                                                                                                                             | `[]`                                                        |
| `volumes`                   | Extra pod volumes.                                                                                                                                                                                                                                                                                           | `[]`                                                        |
| `volumeMounts`              | Extra volume mounts for the MCP container.                                                                                                                                                                                                                                                                   | `[]`                                                        |
| `labels`                    | Extra labels added to MCP Deployment and pod template.                                                                                                                                                                                                                                                       | `{}`                                                        |
| `annotations`               | Extra annotations added to MCP Deployment.                                                                                                                                                                                                                                                                   | `{}`                                                        |
| `podAnnotations`            | Extra annotations added to MCP pods.                                                                                                                                                                                                                                                                         | `{}`                                                        |
<!-- markdownlint-enable line-length -->

`mcp-victoriametrics v1.20.2` supports only
`VM_INSTANCE_BEARER_TOKEN` and `VM_INSTANCE_HEADERS`; it does not support
file-based equivalents. Consequently, `vm.bearerTokenSecret` and
`vm.headersSecret` prevent credentials from appearing in the Deployment but
still inject them into the process environment. Inline `vm.bearerToken` and
`vm.headers` also render their values in the Deployment. Do not use these four
parameters where file-based secret handling is mandatory.

For a hardened HTTP deployment, leave static credentials empty and use
`passthroughHeaders` so the MCP client supplies request-time credentials. The
CA configuration is file-based and is not subject to this limitation: the
chart mounts the selected CA Secret and puts only its path in `SSL_CERT_FILE`.

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

For an automatically generated HTTPS endpoint, the selected VMSingle or
VMSelect TLS Secret must contain the configured CA key, `ca.crt` by default. If
the backend certificate is already trusted through the container system trust
store, disable the Secret mount:

```yaml
victoriametrics:
  mcp:
    vm:
      tls:
        useBackendCA: false
```

For an explicit endpoint that uses a private CA, specify the Secret yourself:

```yaml
victoriametrics:
  mcp:
    vm:
      entrypoint: https://vmauth-monitoring.example.com
      tls:
        caSecret:
          name: vmauth-ca
          key: ca.crt
```

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
    passthroughHeaders:
      - Authorization
    serviceMonitor:
      install: true
    httpRoute:
      install: true
      hostnames:
        - mcp-vm.{{ .Release.Namespace }}.example.com
      parentRefs:
        - name: monitoring-gateway
```

In this hardened example, every MCP client must send a valid `Authorization`
header and the HTTPRoute or Ingress path must preserve it.
