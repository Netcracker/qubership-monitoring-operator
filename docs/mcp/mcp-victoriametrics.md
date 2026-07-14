# Installing mcp-victoriametrics

This guide describes how to install and configure
[`mcp-victoriametrics`](https://github.com/VictoriaMetrics/mcp-victoriametrics),
the VictoriaMetrics Model Context Protocol (MCP) server.

The MCP server exposes VictoriaMetrics read-only APIs and embedded
VictoriaMetrics documentation to MCP-compatible clients, such as Claude
Desktop, Claude Code, Cursor, and other tools that support MCP.

## Table of Contents

- [Installing mcp-victoriametrics](#installing-mcp-victoriametrics)
  - [Table of Contents](#table-of-contents)
  - [Prerequisites](#prerequisites)
  - [Choose VictoriaMetrics Endpoint](#choose-victoriametrics-endpoint)
  - [Self-Signed Certificates](#self-signed-certificates)
  - [Installation](#installation)
    - [Helm](#helm)
    - [Binaries](#binaries)
    - [Docker](#docker)
    - [Source](#source)
  - [Configuration Parameters](#configuration-parameters)
  - [Configure MCP Clients](#configure-mcp-clients)
    - [Claude Desktop](#claude-desktop)
    - [Claude Code](#claude-code)
    - [Cursor](#cursor)
    - [Codex](#codex)
  - [Available Tools](#available-tools)
  - [Verify Installation](#verify-installation)
  - [Uninstall](#uninstall)
  - [References](#references)

## Prerequisites

Before installation, prepare the following:

* VictoriaMetrics stack enabled with `victoriametrics.vmSingle.install: true`.
  In the current operator setup, `VMSingle` is the expected VictoriaMetrics
  backend.
* An MCP-compatible client.

## Choose VictoriaMetrics Endpoint

`mcp-victoriametrics` must be able to reach the root VictoriaMetrics API.

For MCP installed inside the same Kubernetes cluster, use the internal
`VMSingle` service, usually `http://vmsingle-k8s.monitoring.svc:8428` for the
`monitoring` namespace.

Upstream `mcp-victoriametrics` documentation describes
`VM_INSTANCE_ENTRYPOINT` as the root URL of `vmsingle` or `vmselect`. In this
operator setup, the external `VMAuth` URL can be used for local MCP only
because `VMAuth` acts as a transparent proxy to the same `VMSingle` API.

This is expected to work with the operator-generated `VMAuth` configuration:
`VMAuth` proxies VictoriaMetrics API paths such as `/api/v1/query.*`,
`/api/v1/label.*`, `/api/v1/series.*`, `/api/v1/metadata.*`, and
`/api/v1/status.*` to `http://vmsingle-k8s.monitoring.svc:8428`. It also
contains a root route to `VMSingle`, so the external `VMAuth` URL behaves as
the VictoriaMetrics API entrypoint for MCP.

You can verify the generated `VMAuth` config with:

```bash
kubectl get secret vmauth-config-k8s -n monitoring -o jsonpath='{.data.config\.yaml\.gz}' | base64 -d | gunzip
```

Example with Ingress:

```yaml
victoriametrics:
  vmAuth:
    install: true
    ingress:
      install: true
      host: vmauth.example.com
```

Example with HTTPRoute:

```yaml
victoriametrics:
  vmAuth:
    install: true
    httpRoute:
      install: true
      hostnames:
        - vmauth.example.com
```

For these examples, local MCP configuration should use:

```bash
export VM_INSTANCE_ENTRYPOINT="https://vmauth.example.com"
export VM_INSTANCE_TYPE="single"
```

If `VMAuth` requires Basic Auth, pass the credentials either in the URL:

```bash
export VM_INSTANCE_ENTRYPOINT="https://<username>:<password>@vmauth.example.com"
export VM_INSTANCE_TYPE="single"
```

or with an explicit `Authorization` header:

```bash
export VM_INSTANCE_ENTRYPOINT="https://vmauth.example.com"
export VM_INSTANCE_TYPE="single"
export VM_INSTANCE_HEADERS="Authorization=Basic <base64-username-password>"
```

`<base64-username-password>` is the base64-encoded `<username>:<password>`
string. For example, for `admin:admin`:

```bash
printf 'admin:admin' | base64
```

The result is:

```text
YWRtaW46YWRtaW4=
```

So the header value is:

```bash
export VM_INSTANCE_HEADERS="Authorization=Basic YWRtaW46YWRtaW4="
```

## Self-Signed Certificates

If the external API Gateway, HTTPRoute, or Ingress for `VMAuth` or `VMSingle` uses a
self-signed certificate, local `mcp-victoriametrics` will not connect to it
unless the certificate is trusted by the environment where the MCP server runs.

The upstream MCP server configuration and current source code do not provide
an `insecure-skip-verify` option. For local installation, add the issuing CA
certificate to the operating system trust store, or run the MCP server with a
CA bundle that contains this certificate.

For a local binary, one option is to use `SSL_CERT_FILE` env variable, e.g.:

```bash
export SSL_CERT_FILE="/path/to/ca-bundle-with-api-gateway-ca.crt"
export VM_INSTANCE_ENTRYPOINT="https://vmauth.example.com"
export VM_INSTANCE_TYPE="single"
./mcp-victoriametrics
```

For Docker, mount the CA bundle and pass `SSL_CERT_FILE` to the container, for example:

```bash
docker run -d --name mcp-victoriametrics \
  -v /path/to/ca-bundle-with-api-gateway-ca.crt:/etc/ssl/certs/api-gateway-ca.crt:ro,z \
  -e SSL_CERT_FILE=/etc/ssl/certs/api-gateway-ca.crt \
  -e VM_INSTANCE_ENTRYPOINT=https://vmauth.example.com \
  -e VM_INSTANCE_TYPE=single \
  -e MCP_SERVER_MODE=http \
  -e MCP_LISTEN_ADDR=:8080 \
  -p 8080:8080 \
  ghcr.io/victoriametrics/mcp-victoriametrics
```

The `z` suffix is required on SELinux-enabled systems when the certificate is
mounted from the host filesystem. It relabels the file so the container process
can read it.

If the certificate is signed by a company or cluster CA, use that CA
certificate, not the leaf server certificate.

The same rule applies to the opposite connection direction: if an MCP client
connects to the MCP server through an HTTPRoute or Ingress with a self-signed
certificate, the CA must be trusted by the machine where the MCP client runs.
Do not put this CA only into the `mcp-victoriametrics` container: that affects
only outbound requests from MCP to `VMAuth` or `VMSingle`, not inbound requests
from Claude Desktop, Codex, Cursor, or another client to the MCP server URL.

For Linux workstations, add the API Gateway or Ingress CA to the system trust
store. For RHEL, Fedora, or CentOS:

```bash
sudo cp api-gateway-ca.crt /etc/pki/ca-trust/source/anchors/
sudo update-ca-trust extract
```

For Debian or Ubuntu:

```bash
sudo cp api-gateway-ca.crt /usr/local/share/ca-certificates/api-gateway-ca.crt
sudo update-ca-certificates
```

Most MCP clients do not provide a portable per-server
`insecure-skip-verify` setting. Prefer trusting the CA or using a certificate
issued by a CA already trusted by the client machine. For temporary local
debugging, expose the MCP server on `http://localhost:8080/mcp` with Docker or
`kubectl port-forward` instead of using a self-signed HTTPS route.

## Installation

### Helm

Use Helm when the MCP server should run inside the same Kubernetes cluster as
the monitoring stack.

Add the VictoriaMetrics Helm repository:

```bash
helm repo add vm https://victoriametrics.github.io/helm-charts/
helm repo update
```

Export default chart values:

```bash
helm show values vm/victoria-metrics-mcp > values.yaml
```

Add and adjust the following example parameters in `values.yaml`:

```yaml
vm:
  entrypoint: "http://vmsingle-k8s.monitoring.svc:8428"
  type: "single"
  bearerToken: ""

mcp:
  mode: http

service:
  type: ClusterIP
  port: 8080

# Use HTTPRoute to expose MCP server outside the cluster.
route:
  enabled: true
  hostnames:
    - mcp-vm.example.com
  parentRefs:
    - name: gateway
      namespace: istio-system
      sectionName: default
  rules:
    - matches:
        - path:
            type: PathPrefix
            value: /

# Or use Ingress instead of HTTPRoute. In this case set route.enabled to false.
ingress:
  enabled: false
  className: ""
  hosts:
    - host: mcp-vm.example.com
      paths:
        - path: /
          pathType: ImplementationSpecific
```

This snippet is an example, not a set of universal values. Adjust it for your
cluster:

* `vm.entrypoint` must contain the actual monitoring namespace. The URL is
  formed as `http://<vmsingle_service>.<monitoring_namespace>.svc:8428`.
* `route.hostnames` or `ingress.hosts` must contain the desired MCP server host.
* `route.parentRefs` must point to the Gateway that serves HTTPRoute traffic in
  your cluster.
* `ingress.className`, `ingress.hosts`, and TLS settings depend on the Ingress
  controller used in your cluster.
* If VictoriaMetrics requires authentication, set the corresponding `vm`
  authentication parameters.

Install the chart:

replace `monitoring` with the target namespace to install MCP server to.

```bash
helm install vmm vm/victoria-metrics-mcp \
  -f values.yaml \
  -n monitoring
```

Check the release and pods:

```bash
helm list -f vmm -n monitoring
kubectl get pods -n monitoring | grep vmm
```

### Binaries

Use a downloaded binary when you want local `stdio` mode without Docker, or
when the host should run a standalone HTTP MCP server.

Download the binary for your operating system and architecture from the
upstream releases page:

```text
https://github.com/VictoriaMetrics/mcp-victoriametrics/releases
```

Then make it executable and place it somewhere stable, for example:

```bash
chmod +x ./mcp-victoriametrics
sudo mv ./mcp-victoriametrics /usr/local/bin/mcp-victoriametrics
```

For `stdio` usage, do not start the binary manually. Configure the MCP client
with `command: /usr/local/bin/mcp-victoriametrics` and the required environment
variables. The client will start the process itself.

For standalone local HTTP usage, start the binary explicitly:

```bash
export VM_INSTANCE_ENTRYPOINT="https://vmauth.example.com"
export VM_INSTANCE_TYPE="single"
export MCP_SERVER_MODE="http"
export MCP_LISTEN_ADDR=":8080"

mcp-victoriametrics
```

After startup, the streamable HTTP endpoint is available at:

```text
http://localhost:8080/mcp
```

### Docker

Use Docker when the MCP server should run locally or as a standalone HTTP/SSE
service. For local usage, `VM_INSTANCE_ENTRYPOINT` should point to the external
`VMAuth` URL that proxies to `VMSingle` or if `VMAuth` is disabled then `VMSingle` ingress or http route host.

```bash
docker run -d --name mcp-victoriametrics \
  -e VM_INSTANCE_ENTRYPOINT=https://vmauth.example.com \
  -e VM_INSTANCE_TYPE=single \
  -e MCP_SERVER_MODE=http \
  -e MCP_LISTEN_ADDR=:8080 \
  -p 8080:8080 \
  ghcr.io/victoriametrics/mcp-victoriametrics
```

If `VMAuth` requires Basic Auth, pass credentials with
`VM_INSTANCE_HEADERS`:

```bash
docker run -d --name mcp-victoriametrics \
  -e VM_INSTANCE_ENTRYPOINT=https://vmauth.example.com \
  -e VM_INSTANCE_TYPE=single \
  -e MCP_SERVER_MODE=http \
  -e MCP_LISTEN_ADDR=:8080 \
  -e VM_INSTANCE_HEADERS="Authorization=Basic <base64-username-password>" \
  -p 8080:8080 \
  ghcr.io/victoriametrics/mcp-victoriametrics
```

Replace `<base64-username-password>` with base64-encoded
`<username>:<password>`, for example `YWRtaW46YWRtaW4=` for `admin:admin`.

After startup, the streamable HTTP endpoint is available at:

```text
http://localhost:8080/mcp
```

The server also exposes:

* `/` - setup page and tool inspection in HTTP mode
* `/metrics` - MCP server metrics in Prometheus format
* `/health/liveness` - liveness probe
* `/health/readiness` - readiness probe

### Source

Use source installation when you need to build the MCP server yourself. The
source build requires a Go toolchain compatible with the current
`mcp-victoriametrics` repository.

Clone the repository:

```bash
git clone https://github.com/VictoriaMetrics/mcp-victoriametrics.git
cd mcp-victoriametrics
```

Build the binary:

```bash
make build
```

For `stdio` usage, do not start the binary manually. Configure the MCP client
with the path to the built `mcp-victoriametrics` binary and the required
environment variables. The client will start the process itself when it needs
the MCP server.

If you want to run the built binary as a standalone local HTTP server instead,
use the same approach as in the Docker section, but pass the environment
variables to the binary directly and set `MCP_SERVER_MODE=http`.

## Configuration Parameters

`mcp-victoriametrics` is configured with environment variables.

| Variable | Description | Required |
| --- | --- | --- |
| `VM_INSTANCE_ENTRYPOINT` | Root URL of the VictoriaMetrics API. For in-cluster MCP use the direct `VMSingle` URL, for example `http://vmsingle-k8s.monitoring.svc:8428`, and replace `monitoring` if the stack is installed in another namespace. For local MCP, the external `VMAuth` HTTPRoute or Ingress URL can be used when it transparently proxies the root VictoriaMetrics API to `VMSingle`. In case `VMAuth` is not installed direct HTTPRoute or Ingress URL of `VMSingle` can be used. If `VMCluster` is installed, then `VMSelect` HTTPRoute or Ingress URL can be used. | Yes |
| `VM_INSTANCE_TYPE` | VictoriaMetrics type. Use `single` in most cases and `cluster` if VMCluster is deployed. | Yes |
| `VM_INSTANCE_BEARER_TOKEN` | Bearer token sent to the configured VictoriaMetrics entrypoint, for example `vmauth` or an API Gateway. | No |
| `VM_INSTANCE_HEADERS` | Custom headers sent to the configured entrypoint, comma-separated as `Header=Value`. | No |
| `MCP_PASSTHROUGH_HEADERS` | Headers forwarded from incoming MCP requests to the configured entrypoint in `sse` or `http` mode. | No |
| `MCP_SERVER_MODE` | MCP transport mode: `stdio`, `sse`, or `http`. Defaults to `stdio`. | No |
| `MCP_LISTEN_ADDR` | Listen address for `sse` or `http` mode. Defaults to `localhost:8080`. | No |
| `MCP_DISABLED_TOOLS` | Comma-separated list of tools to disable. | No |
| `MCP_LOG_LEVEL` | Log level: `debug`, `info`, `warn`, or `error`. | No |

Use `stdio` mode when the MCP client is configured with a `command` that points
to the `mcp-victoriametrics` binary. In this mode the client starts the binary
as a child process and exchanges MCP messages through standard input and
standard output. There is no HTTP port and no `/mcp` URL.

Use `http` mode when `mcp-victoriametrics` is already running separately, for
example as a Docker container, a manually started binary, or a Helm deployment
in Kubernetes. In this mode the MCP client connects to an endpoint such as
`http://localhost:8080/mcp` or `https://mcp-vm.example.com/mcp`.

Use `sse` mode only if the MCP client explicitly expects an SSE endpoint.

## Configure MCP Clients

### Claude Desktop

For local binary usage in `stdio` mode, add the server command to
`claude_desktop_config.json`. In this case Claude Desktop starts
`mcp-victoriametrics` itself:

```json
{
  "mcpServers": {
    "victoriametrics": {
      "command": "/path/to/mcp-victoriametrics",
      "env": {
        "VM_INSTANCE_ENTRYPOINT": "https://vmauth.example.com",
        "VM_INSTANCE_TYPE": "single",
        "VM_INSTANCE_HEADERS": "Authorization=Basic <base64-username-password>"
      }
    }
  }
}
```

For a separately running HTTP MCP server, configure `url` in
`claude_desktop_config.json` instead of `command`. For Helm installation, use
the external HTTPRoute or Ingress URL:

```json
{
  "mcpServers": {
    "victoriametrics": {
      "url": "https://mcp-vm.example.com/mcp"
    }
  }
}
```

For Docker installation running on the same workstation, use the local HTTP MCP
endpoint:

```json
{
  "mcpServers": {
    "victoriametrics": {
      "url": "http://localhost:8080/mcp"
    }
  }
}
```

For temporary local debugging of the Helm installation without HTTPRoute or
Ingress, you can also use
`kubectl port-forward svc/vmm-victoria-metrics-mcp -n monitoring 8080:8080`
and set `url` in `claude_desktop_config.json` to `http://localhost:8080/mcp`.

### Claude Code

For local source usage:

```bash
claude mcp add victoriametrics -- /path/to/mcp-victoriametrics \
  -e VM_INSTANCE_ENTRYPOINT=https://vmauth.example.com \
  -e VM_INSTANCE_TYPE=single \
  -e VM_INSTANCE_HEADERS="Authorization=Basic <base64-username-password>"
```

### Cursor

Add the server to `~/.cursor/mcp.json`:

```json
{
  "mcpServers": {
    "victoriametrics": {
      "command": "/path/to/mcp-victoriametrics",
      "env": {
        "VM_INSTANCE_ENTRYPOINT": "https://vmauth.example.com",
        "VM_INSTANCE_TYPE": "single",
        "VM_INSTANCE_HEADERS": "Authorization=Basic <base64-username-password>"
      }
    }
  }
}
```

### Codex

For local binary usage in `stdio` mode, add the server to
`~/.codex/config.toml`. In this mode Codex starts `mcp-victoriametrics` itself:

```toml
[mcp_servers.victoriametrics]
command = "/path/to/mcp-victoriametrics"

[mcp_servers.victoriametrics.env]
VM_INSTANCE_ENTRYPOINT = "https://vmauth.example.com"
VM_INSTANCE_TYPE = "single"
VM_INSTANCE_HEADERS = "Authorization=Basic <base64-username-password>"
```

## Available Tools

The exact set of tools depends on the `mcp-victoriametrics` version,
VictoriaMetrics instance type, and `MCP_DISABLED_TOOLS` configuration.

Commonly available tools are:

| Tool | Description | Enabled by default |
| --- | --- | --- |
| `query` | Execute instant PromQL or MetricsQL queries. | Yes |
| `query_range` | Execute range PromQL or MetricsQL queries over a time period. | Yes |
| `metrics` | List available metric names. | Yes |
| `metrics_metadata` | Read stored metric metadata, such as type, help, and unit. | Yes |
| `labels` | List available label names. | Yes |
| `label_values` | List values for a specific label. | Yes |
| `series` | List available time series. | Yes |
| `rules` | View alerting and recording rules. | Yes |
| `alerts` | View current firing and pending alerts. | Yes |
| `metric_statistics` | Inspect metric usage in queries. | Yes |
| `active_queries` | View currently executing queries. | Yes |
| `top_queries` | View frequent or slow queries. | Yes |
| `tsdb_status` | Inspect TSDB cardinality statistics. | Yes |
| `tenants` | List tenants in a multi-tenant cluster setup. Mostly useful for `VMCluster`. | Yes |
| `documentation` | Search embedded VictoriaMetrics documentation without online access. | Yes |
| `prettify_query` | Format PromQL or MetricsQL queries. | Yes |
| `explain_query` | Parse and explain PromQL or MetricsQL queries. | Yes |
| `export` | Export raw time series data to JSON or CSV. | No |
| `flags` | View non-default VictoriaMetrics startup flags. | No |
| `metric_relabel_debug` | Debug Prometheus-compatible relabeling rules. | No |
| `downsampling_filters_debug` | Debug downsampling configuration. | No |
| `retention_filters_debug` | Debug retention filter configuration. | No |
| `test_rules` | Unit-test alerting and recording rules with `vmalert`. | No |

By default upstream disables several potentially heavier or more specialized
tools, including `export`, `flags`, `metric_relabel_debug`,
`downsampling_filters_debug`, `retention_filters_debug`, and `test_rules`.
Override `MCP_DISABLED_TOOLS` only when the MCP clients really need these
capabilities.

Useful tools for smoke testing are:

| Tool | Purpose |
| --- | --- |
| `metrics` | Verify that MCP can authenticate to VictoriaMetrics and list metric names. |
| `labels` | Verify access to label metadata. |
| `label_values` | Verify access to values for a known label, for example `job` or `namespace`. |
| `query` | Run a simple instant query, for example `up`. |
| `tsdb_status` | Check cardinality statistics for troubleshooting. |

VictoriaMetrics Cloud tools are intentionally not covered here because this
operator setup uses a Kubernetes-hosted VictoriaMetrics stack, usually
`VMSingle`.

## Verify Installation

For HTTP or SSE mode, check the readiness endpoint:

```bash
curl http://localhost:8080/health/readiness
```

For `stdio` mode, there is no readiness URL because the MCP server is attached
to the client's process through stdin/stdout. Restart the MCP client and check
that the `victoriametrics` server appears in the client's MCP server list.

Run a simple prompt in the MCP client, for example:

```text
List available VictoriaMetrics metrics and show the labels for the up metric.
```

## Uninstall

For Docker:

```bash
docker rm -f mcp-victoriametrics
```

For Helm:

```bash
helm uninstall vmm -n monitoring
```

For binary or source installation, remove the local `mcp-victoriametrics`
binary if it is no longer needed. For source installation, remove the cloned
source directory as well.

## References

* [VictoriaMetrics MCP Server GitHub repository](https://github.com/VictoriaMetrics/mcp-victoriametrics)
* [VictoriaMetrics MCP Server releases](https://github.com/VictoriaMetrics/mcp-victoriametrics/releases)
* [VictoriaMetrics MCP Helm chart documentation](https://docs.victoriametrics.com/helm/victoria-metrics-mcp/)
