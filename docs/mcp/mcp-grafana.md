# Installing mcp-grafana

This guide describes how to install and configure
[`mcp-grafana`](https://github.com/grafana/mcp-grafana), the Grafana Model
Context Protocol (MCP) server.

The MCP server gives MCP-compatible clients access to Grafana dashboards,
datasources, alerting, annotations, snapshots, and query tools.

This guide uses upstream distributions for local installation with `uvx`,
Docker, or a binary. In-cluster Helm installation uses the Qubership Monitoring
Operator chart and its `grafana.mcp` values. Upstream Helm chart values use a
different schema and cannot be copied into the monitoring chart.

## Table of Contents

- [Installing mcp-grafana](#installing-mcp-grafana)
  - [Table of Contents](#table-of-contents)
  - [Prerequisites](#prerequisites)
  - [Choose Grafana Endpoint](#choose-grafana-endpoint)
  - [Authentication](#authentication)
  - [Protect MCP Endpoint](#protect-mcp-endpoint)
  - [Self-Signed Certificates](#self-signed-certificates)
  - [Installation](#installation)
    - [UVX Recommended](#uvx-recommended)
    - [Docker](#docker)
    - [Binary Release](#binary-release)
    - [Go Install from Source](#go-install-from-source)
    - [Helm with Monitoring Chart](#helm-with-monitoring-chart)
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

* Grafana enabled with `grafana.install: true`.
* Grafana 9.0 or later.
* An MCP-compatible client.
* A Grafana service account token, or Grafana username and password.
* `uv`, only if you plan to run local `stdio` mode through `uvx`.

Prefer a Grafana service account token. Use username/password only for local
debugging or when service accounts are not available.

## Choose Grafana Endpoint

`mcp-grafana` must be able to reach the root Grafana URL.

For MCP installed inside the same Kubernetes cluster, use the internal Grafana
service, usually:

```text
http://grafana-service.monitoring.svc:3000
```

Replace `monitoring` if the monitoring stack is installed in another
namespace.

For local MCP installation, use the external Grafana HTTPRoute or Ingress URL:

```text
https://grafana.example.com
```

If Grafana is not exposed externally, use port-forwarding while creating a
token or testing a local MCP server:

```bash
kubectl port-forward svc/grafana-service -n monitoring 3000:3000
```

Then use:

```text
http://localhost:3000
```

## Authentication

`GRAFANA_SERVICE_ACCOUNT_TOKEN` is not an operator parameter. It is an
environment variable read by `mcp-grafana`. The value must be a Grafana service
account token: a generated token that allows automated tools to call the
Grafana HTTP API without using a human user's password.

Create a Grafana service account token with enough permissions for the tools
you plan to use. In Grafana UI:

1. Sign in to Grafana as an administrator.
2. Open **Administration** -> **Users and access** -> **Service accounts**.
3. Click **Add service account**.
4. Set a name, for example `mcp-grafana`.
5. Assign a role. Use `Viewer` for read-only inspection; use `Editor` only if
   write operations are required.
6. Open the created service account and click **Add service account token**.
7. Copy the generated token and store it in a secret manager or Kubernetes
   Secret. Grafana shows the token only once.

For read-only usage, grant permissions for reading dashboards, folders,
datasources, alerting configuration, and querying datasources. If you need a
simple setup for testing, the built-in `Editor` role is usually enough, but it
is broader than least-privilege access.

Set the token for local usage with:

```bash
export GRAFANA_URL="https://grafana.example.com"
export GRAFANA_SERVICE_ACCOUNT_TOKEN="<service-account-token>"
```

If using basic authentication, set:

```bash
export GRAFANA_URL="https://grafana.example.com"
export GRAFANA_USERNAME="<username>"
export GRAFANA_PASSWORD="<password>"
```

For the operator-managed Grafana admin credentials, the secret is usually
`grafana-admin-credentials`:

```bash
kubectl get secret grafana-admin-credentials -n monitoring -o jsonpath='{.data.admin-user}' | base64 -d
kubectl get secret grafana-admin-credentials -n monitoring -o jsonpath='{.data.admin-password}' | base64 -d
```

Use these credentials to create a service account token in Grafana instead of
putting the admin password into long-lived MCP configuration.

If the Grafana UI is not available, create the service account and token
through the Grafana HTTP API. First port-forward Grafana if needed:

```bash
kubectl port-forward svc/grafana-service -n monitoring 3000:3000
```

Then create the service account:

```bash
curl -sS -u '<admin-user>:<admin-password>' \
  -H 'Content-Type: application/json' \
  -X POST http://localhost:3000/api/serviceaccounts \
  -d '{"name":"mcp-grafana","role":"Viewer","isDisabled":false}'
```

The response contains the service account `id`. Use this `id` to create a
token:

```bash
curl -sS -u '<admin-user>:<admin-password>' \
  -H 'Content-Type: application/json' \
  -X POST http://localhost:3000/api/serviceaccounts/<service-account-id>/tokens \
  -d '{"name":"mcp-grafana-token"}'
```

The response contains `key`. Use this value as
`GRAFANA_SERVICE_ACCOUNT_TOKEN`.

## Protect MCP Endpoint

`mcp-grafana` does not provide built-in authentication for incoming MCP client
requests to `/mcp` or `/sse`.

`GRAFANA_SERVICE_ACCOUNT_TOKEN`, `GRAFANA_USERNAME`, and `GRAFANA_PASSWORD`
authenticate `mcp-grafana` to Grafana. They do not authenticate Claude
Desktop, Codex, Cursor, or another MCP client to the MCP server itself.

The `--allowed-hosts` and `--allowed-origins` options are also not user
authentication. They protect the HTTP/SSE listener from DNS rebinding and
unexpected browser-origin requests.

If the MCP server is exposed through HTTPRoute, Ingress, or API Gateway,
protect it at that layer. Use one of the mechanisms already supported by your
gateway, for example:

* OAuth2/OIDC authentication.
* Basic Auth.
* mTLS client certificates.
* IP allowlists or private network access.
* Kubernetes NetworkPolicy when the endpoint is internal.

Before choosing gateway-level authentication, verify that the MCP client can
send the required credentials to a remote HTTP MCP server. If the client cannot
send the required headers or client certificate, prefer `stdio` mode or a
local-only HTTP endpoint such as `http://localhost:8000/`.

## Self-Signed Certificates

There are two TLS connection directions to consider.

For `mcp-grafana` connecting to Grafana, use `--tls-ca-file` when the Grafana
HTTPRoute or Ingress certificate is signed by a private CA:

```bash
mcp-grafana --tls-ca-file /path/to/api-gateway-ca.crt
```

For temporary local testing only, `mcp-grafana` also supports:

```bash
mcp-grafana --tls-skip-verify
```

Do not use `--tls-skip-verify` in a shared or production setup.

For Docker on SELinux-enabled systems, mount the CA with the `z` suffix:

```bash
docker run -d --name mcp-grafana \
  -v /path/to/api-gateway-ca.crt:/certs/api-gateway-ca.crt:ro,z \
  -e GRAFANA_URL=https://grafana.example.com \
  -e GRAFANA_SERVICE_ACCOUNT_TOKEN=<service-account-token> \
  -p 8000:8000 \
  grafana/mcp-grafana \
  -t streamable-http \
  --address 0.0.0.0:8000 \
  --endpoint-path / \
  --allowed-hosts localhost:8000 \
  --tls-ca-file /certs/api-gateway-ca.crt \
  --disable-write
```

The `z` suffix relabels the mounted file so the container process can read it
on SELinux-enabled systems.

For an MCP client connecting to `mcp-grafana` through an HTTPRoute or Ingress
with a self-signed certificate, trust the CA on the machine where the MCP
client runs. The CA mounted into the `mcp-grafana` container affects only
outbound requests from MCP to Grafana.

For temporary local debugging, expose the MCP server on
`http://localhost:8000/` with Docker or `kubectl port-forward` instead of
using a self-signed HTTPS route.

## Installation

### UVX Recommended

Use `uvx` for local `stdio` mode. This is the upstream recommended way to run
`mcp-grafana` when the MCP client starts the server process itself.

Install `uv` first. The cross-platform upstream installer is:

```bash
curl -LsSf https://astral.sh/uv/install.sh | sh
```

If your operating system provides a packaged `uv`, you can use the native
package manager instead, for example `sudo dnf install uv` on Fedora.

If you prefer Python tooling, install `uv` with `pipx`:

```bash
pipx install uv
```

Or with `pip`:

```bash
python3 -m pip install --user uv
```

After installation, verify that `mcp-grafana` can be started:

```bash
uvx mcp-grafana --help
```

For normal usage, do not keep `uvx mcp-grafana` running manually. Configure the
MCP client with `command: uvx` and `args: ["mcp-grafana", "--disable-write"]`.
The client will start the process itself when it needs the MCP server.

For a quick local smoke test, run:

```bash
GRAFANA_URL=https://grafana.example.com \
GRAFANA_SERVICE_ACCOUNT_TOKEN=<service-account-token> \
uvx mcp-grafana --disable-write
```

### Docker

Use Docker when the MCP server should run locally as a standalone streamable
HTTP server. Configure the MCP client with the local URL.

```bash
docker run -d --name mcp-grafana \
  -e GRAFANA_URL=https://grafana.example.com \
  -e GRAFANA_SERVICE_ACCOUNT_TOKEN=<service-account-token> \
  -p 8000:8000 \
  grafana/mcp-grafana \
  -t streamable-http \
  --address 0.0.0.0:8000 \
  --endpoint-path / \
  --allowed-hosts localhost:8000 \
  --disable-write
```

After startup, the streamable HTTP endpoint is available at:

```text
http://localhost:8000/
```

The server also exposes:

* `/healthz` - health endpoint
* `/metrics` - MCP server metrics, when started with `--metrics`

### Binary Release

Use a downloaded binary when you want local `stdio` mode without `uvx`, or when
the host should run a standalone HTTP MCP server without Docker.

Download the binary for your operating system and architecture from the
upstream releases page:

```text
https://github.com/grafana/mcp-grafana/releases
```

Then make it executable and place it somewhere stable, for example:

```bash
chmod +x ./mcp-grafana
sudo mv ./mcp-grafana /usr/local/bin/mcp-grafana
```

For `stdio` usage, do not start the binary manually. Configure the MCP client
with `command: /usr/local/bin/mcp-grafana` and the required environment
variables. The client will start the process itself.

For standalone local HTTP usage, start the binary explicitly:

```bash
export GRAFANA_URL="https://grafana.example.com"
export GRAFANA_SERVICE_ACCOUNT_TOKEN="<service-account-token>"

mcp-grafana \
  -t streamable-http \
  --address 0.0.0.0:8000 \
  --endpoint-path /mcp \
  --allowed-hosts localhost:8000 \
  --disable-write
```

### Go Install from Source

Use Go install when you need to build the latest `mcp-grafana` from the
upstream source module.

```bash
GOBIN="$HOME/go/bin" go install github.com/grafana/mcp-grafana/cmd/mcp-grafana@latest
```

The binary is installed to:

```text
$HOME/go/bin/mcp-grafana
```

For `stdio` usage, configure the MCP client with this binary path. For
standalone local HTTP usage, start it explicitly:

```bash
export GRAFANA_URL="https://grafana.example.com"
export GRAFANA_SERVICE_ACCOUNT_TOKEN="<service-account-token>"

"$HOME/go/bin/mcp-grafana" \
  -t streamable-http \
  --address 0.0.0.0:8000 \
  --endpoint-path /mcp \
  --allowed-hosts localhost:8000 \
  --disable-write
```

### Helm with Monitoring Chart

Use the Qubership Monitoring Operator Helm chart when the MCP server should run
in the same Kubernetes cluster as the operator-managed Grafana instance. The
MCP server is installed as part of the Grafana subchart and is not a standalone
Helm release.

The values in this section belong to the Qubership Monitoring Operator chart.
They are not compatible with the values of the upstream
`grafana-community/grafana-mcp` chart.

Create a Kubernetes secret with the Grafana service account token:

```bash
kubectl create namespace monitoring --dry-run=client -o yaml | kubectl apply -f -
kubectl create secret generic grafana-mcp-token \
  -n monitoring \
  --from-literal=token='<service-account-token>'
```

Add and adjust the following example parameters in the monitoring chart values
file:

```yaml
grafana:
  install: true
  mcp:
    install: true
    existingSecret:
      name: grafana-mcp-token
      key: token
    disableWrite: true
    httpRoute:
      install: true
      hostnames:
        - mcp-grafana.example.com
      parentRefs:
        - name: gateway
          namespace: istio-system
          sectionName: default
```

This snippet is an example, not a set of universal values. Adjust it for your
cluster:

* `grafana.mcp.existingSecret` must point to the Secret that contains the
  Grafana service account token.
* `grafana.mcp.httpRoute.hostnames` must contain the desired MCP server host.
* `grafana.mcp.httpRoute.parentRefs` must point to the Gateway that serves
  HTTPRoute traffic in your cluster.
* HTTPRoute and Ingress hosts are added to `--allowed-hosts` automatically.
  Use `grafana.mcp.allowedHosts` only for additional Host values, such as a
  hostname produced by a proxy that rewrites the original Host header.
* Keep `grafana.mcp.disableWrite: true` unless the MCP client is expected to
  create or modify Grafana resources.
* Use `grafana.mcp.ingress` instead of `grafana.mcp.httpRoute` when the cluster
  exposes applications through Ingress.

See the complete
[Grafana MCP Helm parameters](../installation/components/grafana-stack/mcp.md)
reference before enabling authentication, forwarded headers, TLS, metrics, or
advanced pod settings.

Apply the values through the normal monitoring chart installation or upgrade
process. For a source checkout and a release named `monitoring-operator`:

```bash
helm upgrade --install monitoring-operator \
  charts/qubership-monitoring-operator \
  -f values.yaml \
  -n monitoring \
  --create-namespace
```

Check the deployed resources:

```bash
kubectl get deployment,service,httproute -n monitoring | grep mcp-grafana
```

The streamable HTTP endpoint for the example is:

```text
https://mcp-grafana.example.com/
```

When the monitoring chart is obtained from a Helm repository instead of a
source checkout, replace `charts/qubership-monitoring-operator` with the chart
reference used by your environment.

The MCP server requires `grafana.install: true`. Setting only
`grafana.mcp.install: true` does not enable the Grafana subchart or create a
standalone MCP deployment.

If request-time authentication to Grafana is used instead of static MCP pod
credentials, configure the forwarded headers described in the Authentication
section and make every MCP client send the required header.

Protect the external MCP HTTPRoute or Ingress independently at the API Gateway,
Ingress controller, or service mesh layer.

The upstream `grafana-community/grafana-mcp` chart remains an option for a
standalone deployment, but it has a different values schema and is outside the
scope of this monitoring chart example.

## Configuration Parameters

`mcp-grafana` is configured with environment variables and command-line flags.
These are upstream server runtime parameters used directly by local and Docker
installations. The monitoring Helm chart maps its `grafana.mcp` values to these
parameters; use the
[Grafana MCP Helm parameters](../installation/components/grafana-stack/mcp.md)
reference when configuring an in-cluster deployment.

| Parameter                            | Description                                                                                                                                             | Required                            |
| ------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------- |
| `GRAFANA_URL`                        | Root Grafana URL. For in-cluster MCP use `http://grafana-service.monitoring.svc:3000`. For local MCP use the external Grafana HTTPRoute or Ingress URL. | Yes                                 |
| `GRAFANA_SERVICE_ACCOUNT_TOKEN`      | Grafana service account token. Preferred authentication method.                                                                                         | Yes, unless using username/password |
| `GRAFANA_SERVICE_ACCOUNT_TOKEN_FILE` | File that contains the Grafana service account token. Useful for Kubernetes Secrets mounted as files.                                                   | No                                  |
| `GRAFANA_USERNAME`                   | Grafana username for basic authentication.                                                                                                              | Yes, unless using token             |
| `GRAFANA_PASSWORD`                   | Grafana password for basic authentication.                                                                                                              | Yes, unless using token             |
| `GRAFANA_ORG_ID`                     | Numeric Grafana organization ID.                                                                                                                        | No                                  |
| `GRAFANA_EXTRA_HEADERS`              | JSON object with extra headers sent to Grafana API requests.                                                                                            | No                                  |
| `GRAFANA_FORWARD_HEADERS`            | Comma-separated list of incoming headers to forward to Grafana. Applies only to SSE or streamable HTTP transports.                                      | No                                  |
| `--transport`, `-t`                  | MCP transport: `stdio`, `sse`, or `streamable-http`. Defaults to `stdio` for the binary; the Docker image entrypoint defaults to SSE.                   | No                                  |
| `--address`                          | Listen address for SSE or streamable HTTP. Defaults to `localhost:8000`.                                                                                | No                                  |
| `--endpoint-path`                    | Endpoint path for streamable HTTP. Defaults to `/`.                                                                                                     | No                                  |
| `--allowed-hosts`                    | Comma-separated allowlist for HTTP `Host` headers. Required when exposing the server through HTTPRoute or Ingress with an external hostname.            | No                                  |
| `--allowed-origins`                  | Comma-separated allowlist for HTTP `Origin` headers. Usually not needed for desktop MCP clients.                                                        | No                                  |
| `--disable-write`                    | Disable write operations against Grafana. Recommended by default.                                                                                       | No                                  |
| `--enabled-tools`                    | Comma-separated list of enabled tool categories. Some categories are disabled by default upstream.                                                      | No                                  |
| `--tls-ca-file`                      | CA certificate file used by MCP when connecting to Grafana.                                                                                             | No                                  |
| `--tls-skip-verify`                  | Skip Grafana TLS certificate verification. Use only for temporary local testing.                                                                        | No                                  |
| `--server.tls-cert-file`             | Server TLS certificate for streamable HTTP.                                                                                                             | No                                  |
| `--server.tls-key-file`              | Server TLS private key for streamable HTTP.                                                                                                             | No                                  |

Use `stdio` mode when the MCP client is configured with a `command` that points
to a local `mcp-grafana` binary or to `uvx`. In this mode the client starts the
server process and exchanges MCP messages through stdin/stdout.

Use `streamable-http` mode when `mcp-grafana` should run separately,
for example as a Docker container, a manually started binary, or a Helm
deployment in Kubernetes. In this mode the MCP client connects to an endpoint
such as `http://localhost:8000/` or `https://mcp-grafana.example.com/`.

Use `sse` mode only if the MCP client explicitly expects an SSE endpoint such
as `http://localhost:8000/sse`.

## Configure MCP Clients

### Claude Desktop

For local binary usage in `stdio` mode, add the server command to
`claude_desktop_config.json`. In this case Claude Desktop starts
`mcp-grafana` itself:

```json
{
  "mcpServers": {
    "grafana": {
      "command": "/path/to/mcp-grafana",
      "args": ["--disable-write"],
      "env": {
        "GRAFANA_URL": "https://grafana.example.com",
        "GRAFANA_SERVICE_ACCOUNT_TOKEN": "<service-account-token>"
      }
    }
  }
}
```

For `uvx` usage:

```json
{
  "mcpServers": {
    "grafana": {
      "command": "uvx",
      "args": ["mcp-grafana", "--disable-write"],
      "env": {
        "GRAFANA_URL": "https://grafana.example.com",
        "GRAFANA_SERVICE_ACCOUNT_TOKEN": "<service-account-token>"
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
    "grafana": {
      "url": "https://mcp-grafana.example.com/"
    }
  }
}
```

For Docker installation running on the same workstation, use the local HTTP MCP
endpoint:

```json
{
  "mcpServers": {
    "grafana": {
      "url": "http://localhost:8000/"
    }
  }
}
```

### Claude Code

For local binary usage:

```bash
claude mcp add grafana -- /path/to/mcp-grafana \
  --disable-write \
  -e GRAFANA_URL=https://grafana.example.com \
  -e GRAFANA_SERVICE_ACCOUNT_TOKEN=<service-account-token>
```

For `uvx` usage:

```bash
claude mcp add grafana -- uvx mcp-grafana \
  --disable-write \
  -e GRAFANA_URL=https://grafana.example.com \
  -e GRAFANA_SERVICE_ACCOUNT_TOKEN=<service-account-token>
```

### Cursor

Add the server to `~/.cursor/mcp.json`:

```json
{
  "mcpServers": {
    "grafana": {
      "command": "/path/to/mcp-grafana",
      "args": ["--disable-write"],
      "env": {
        "GRAFANA_URL": "https://grafana.example.com",
        "GRAFANA_SERVICE_ACCOUNT_TOKEN": "<service-account-token>"
      }
    }
  }
}
```

For `uvx` usage:

```json
{
  "mcpServers": {
    "grafana": {
      "command": "uvx",
      "args": ["mcp-grafana", "--disable-write"],
      "env": {
        "GRAFANA_URL": "https://grafana.example.com",
        "GRAFANA_SERVICE_ACCOUNT_TOKEN": "<service-account-token>"
      }
    }
  }
}
```

### Codex

For local binary usage in `stdio` mode, add the server to
`~/.codex/config.toml`. In this mode Codex starts `mcp-grafana` itself:

```toml
[mcp_servers.grafana]
command = "/path/to/mcp-grafana"
args = ["--disable-write"]

[mcp_servers.grafana.env]
GRAFANA_URL = "https://grafana.example.com"
GRAFANA_SERVICE_ACCOUNT_TOKEN = "<service-account-token>"
```

For `uvx` usage:

```toml
[mcp_servers.grafana]
command = "uvx"
args = ["mcp-grafana", "--disable-write"]

[mcp_servers.grafana.env]
GRAFANA_URL = "https://grafana.example.com"
GRAFANA_SERVICE_ACCOUNT_TOKEN = "<service-account-token>"
```

For a separately running HTTP MCP server:

```toml
[mcp_servers.grafana]
url = "https://mcp-grafana.example.com/"
```

## Available Tools

The exact set of tools depends on the `mcp-grafana` version, enabled tool
categories, Grafana plugins, and permissions granted to the Grafana service
account token.

The upstream server exposes many tool categories. Some datasource-specific
categories are disabled by default upstream and must be explicitly included in
`--enabled-tools` if they are needed.

| Tool category   | Example tools or capability                                                                                               |
| --------------- | ------------------------------------------------------------------------------------------------------------------------- |
| `admin`         | List teams, users, roles, role assignments, and resource permissions.                                                     |
| `search`        | Search dashboards and folders.                                                                                            |
| `dashboard`     | Read dashboards, inspect panel queries, extract dashboard properties, and update dashboards when write access is enabled. |
| `runpanelquery` | Execute stored dashboard panel queries. Disabled by default upstream.                                                     |
| `datasource`    | List datasources, get datasource details, and check datasource health.                                                    |
| `examples`      | Get datasource query examples. Disabled by default upstream.                                                              |
| `prometheus`    | Run PromQL queries and inspect metric names, labels, label values, metadata, and histograms.                              |
| `loki`          | Run LogQL queries, inspect labels, query stats and patterns, and analyze Loki label usage.                                |
| `alerting`      | Read and manage alert rules, notification policies, contact points, and time intervals.                                   |
| `incident`      | List, create, inspect, and update Grafana Incident objects.                                                               |
| `oncall`        | List OnCall schedules, shifts, teams, users, and alert groups.                                                            |
| `sift`          | Inspect Sift investigations and run Sift analyses when available.                                                         |
| `asserts`       | Read Grafana Asserts summaries when the plugin is available.                                                              |
| `pyroscope`     | Inspect and query profiling data from Pyroscope datasources.                                                              |
| `annotations`   | Read and manage Grafana annotations and annotation tags.                                                                  |
| `snapshot`      | List, inspect, create, and delete dashboard snapshots.                                                                    |
| `rendering`     | Render dashboard or panel images when rendering is configured.                                                            |
| `navigation`    | Generate Grafana deeplink URLs.                                                                                           |
| `provisioning`  | List provisioning repositories and validate provisioning files.                                                           |
| `influxdb`      | Query InfluxDB datasources.                                                                                               |
| `elasticsearch` | Query Elasticsearch or OpenSearch datasources. Disabled by default upstream.                                              |
| `quickwit`      | Query Quickwit datasources. Disabled by default upstream.                                                                 |
| `clickhouse`    | Inspect and query ClickHouse datasources. Disabled by default upstream.                                                   |
| `cloudwatch`    | Inspect and query AWS CloudWatch datasources. Disabled by default upstream.                                               |
| `athena`        | Inspect and query Athena datasources. Disabled by default upstream.                                                       |
| `snowflake`     | Inspect and query Snowflake datasources. Disabled by default upstream.                                                    |
| `graphite`      | Query Graphite datasources. Disabled by default upstream.                                                                 |

For this chart, the default `grafana.mcp.enabledTools` value keeps the MCP
server focused on monitoring troubleshooting and reduces the number of tools
advertised to MCP clients:

```yaml
grafana:
  mcp:
    enabledTools:
      - search
      - datasource
      - dashboard
      - prometheus
      - loki
      - alerting
```

Set `grafana.mcp.enabledTools: []` to use the upstream default set, or provide
a custom list if your Grafana instance needs additional categories such as
`pyroscope`, `oncall`, `incident`, or datasource-specific SQL tools.

Useful tools for smoke testing are:

| Tool                          | Purpose                                                                   |
| ----------------------------- | ------------------------------------------------------------------------- |
| `list_datasources`            | Verify that MCP can authenticate to Grafana and read datasource metadata. |
| `get_datasource`              | Inspect one datasource by name or UID.                                    |
| `check_datasources_health`    | Check whether Grafana can reach configured datasources.                   |
| `list_prometheus_label_names` | Verify access to a Prometheus-compatible datasource.                      |
| `list_loki_label_names`       | Verify access to a Loki-compatible datasource.                            |

If `--disable-write` is enabled, read-only tools continue to work, but tools
that create or modify Grafana resources are disabled or rejected.

## Verify Installation

For streamable HTTP or SSE mode, check the health endpoint:

```bash
curl http://localhost:8000/healthz
```

For `stdio` mode, there is no readiness URL because the MCP server is attached
to the client's process through stdin/stdout. Restart the MCP client and check
that the `grafana` server appears in the client's MCP server list.

Run a simple prompt in the MCP client, for example:

```text
List Grafana datasources and search for dashboards related to Kubernetes.
```

## Uninstall

For Docker:

```bash
docker rm -f mcp-grafana
```

For an MCP server installed with the monitoring chart, set
`grafana.mcp.install: false` and upgrade the monitoring release. Do not
uninstall the monitoring release only to remove MCP.

For binary or Go install, remove the local `mcp-grafana` binary if it is no
longer needed.

## References

* [`mcp-grafana` upstream repository](https://github.com/grafana/mcp-grafana)
* [`mcp-grafana` releases](https://github.com/grafana/mcp-grafana/releases)
* [Upstream standalone `grafana-mcp` Helm chart](https://github.com/grafana-community/helm-charts/tree/main/charts/grafana-mcp)
* [Grafana service account documentation](https://grafana.com/docs/grafana/latest/administration/service-accounts/)
* [Model Context Protocol](https://modelcontextprotocol.io/)
