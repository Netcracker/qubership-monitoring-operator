# Grafana MCP

The `grafana.mcp` section deploys an in-cluster
[`mcp-grafana`](https://github.com/grafana/mcp-grafana) server together with
the Grafana stack.

The MCP server is installed only when both `grafana.install=true` and
`grafana.mcp.install=true`. Standalone deployment against an external Grafana
instance is not supported by this umbrella chart. `grafana.mcp.grafanaUrl` can
override the connection URL, but it does not remove the Grafana installation
requirement.

The MCP server exposes Grafana capabilities to MCP clients. It does not provide
its own inbound authentication for the HTTP endpoint. If the MCP server is
published through HTTPRoute or Ingress, protect it on the gateway, ingress, or
service mesh layer.

By default, the chart keeps the advertised tool set focused on monitoring
troubleshooting:

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

Set `grafana.mcp.enabledTools: []` to use the upstream default tool categories.

## Parameters

All parameters listed below are configured under the `grafana.mcp` level.


<!-- markdownlint-disable line-length -->
| Parameter                       | Description                                                                                                                                                                                           | Default                                                              |
| ------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------- |
| `install`                       | Enable deployment of the Grafana MCP server.                                                                                                                                                          | `false`                                                              |
| `name`                          | Base name for MCP resources.                                                                                                                                                                          | `mcp-grafana`                                                        |
| `image.registry`                | Image registry for `mcp-grafana`.                                                                                                                                                                     | `docker.io`                                                          |
| `image.repository`              | Image repository for `mcp-grafana`.                                                                                                                                                                   | `grafana/mcp-grafana`                                                |
| `image.tag`                     | Image tag.                                                                                                                                                                                            | `0.17.2`                                                             |
| `image.pullPolicy`              | Image pull policy.                                                                                                                                                                                    | `IfNotPresent`                                                       |
| `imagePullPolicy`               | Deprecated image pull policy field. Use `image.pullPolicy`.                                                                                                                                           | `""`                                                                 |
| `imagePullSecrets`              | Image pull secrets for the MCP pod.                                                                                                                                                                   | `[]`                                                                 |
| `replicas`                      | Number of MCP server replicas.                                                                                                                                                                        | `1`                                                                  |
| `grafanaUrl`                    | Grafana URL reachable from the MCP pod. Empty value defaults to `http://grafana-service.<namespace>.svc:3000`.                                                                                        | `""`                                                                 |
| `existingSecret.name`           | Existing Secret that contains the Grafana service account token.                                                                                                                                      | `""`                                                                 |
| `existingSecret.key`            | Key in `existingSecret.name` that contains the token.                                                                                                                                                 | `token`                                                              |
| `serviceAccountToken`           | Inline Grafana service account token. If set, Helm creates a Secret. Prefer `existingSecret` for production.                                                                                          | `""`                                                                 |
| `tokenSecret.name`              | Name of the Secret created for `serviceAccountToken`. Empty value defaults to `<mcp-name>-token`.                                                                                                     | `""`                                                                 |
| `tokenSecret.key`               | Key in the token Secret.                                                                                                                                                                              | `token`                                                              |
| `basicAuth.username`            | Grafana username for basic auth. Prefer service account token authentication.                                                                                                                         | `""`                                                                 |
| `basicAuth.password`            | Grafana password for basic auth.                                                                                                                                                                      | `""`                                                                 |
| `orgId`                         | Numeric Grafana organization ID.                                                                                                                                                                      | `""`                                                                 |
| `debug`                         | Enable debug logging in `mcp-grafana`.                                                                                                                                                                | `false`                                                              |
| `disableWrite`                  | Pass `--disable-write` to disable Grafana write operations.                                                                                                                                           | `true`                                                               |
| `enabledTools`                  | Comma-separated allow-list rendered as `--enabled-tools`. Set to `[]` to use the upstream default tool categories.                                                                                    | `["search","datasource","dashboard","prometheus","loki","alerting"]` |
| `disabledCategories`            | Additional categories rendered as `--disable-<category>`. Usually not needed when `enabledTools` is set.                                                                                              | `[]`                                                                 |
| `metrics.enabled`               | Pass `--metrics` and enable the `/metrics` endpoint. Required for `serviceMonitor.install`.                                                                                                           | `false`                                                              |
| `transport`                     | MCP transport for an in-cluster HTTP server: `streamable-http` or `sse`.                                                                                                                              | `streamable-http`                                                    |
| `endpointPath`                  | Endpoint path for `streamable-http` transport. Ignored for `sse`.                                                                                                                                     | `/`                                                                  |
| `allowedHosts`                  | Additional HTTP `Host` headers rendered through `tpl` and passed as `--allowed-hosts`. `localhost:<targetPort>` and hosts from enabled HTTPRoute or Ingress configuration are included automatically. | `[]`                                                                 |
| `allowedOrigins`                | Allowed HTTP `Origin` headers passed as `--allowed-origins`.                                                                                                                                          | `[]`                                                                 |
| `tls.caFile`                    | CA file used by MCP when connecting to Grafana.                                                                                                                                                       | `""`                                                                 |
| `tls.skipVerify`                | Skip TLS verification when connecting to Grafana. Use only for temporary testing.                                                                                                                     | `false`                                                              |
| `extraArgs`                     | Additional command-line arguments appended to the container args.                                                                                                                                     | `[]`                                                                 |
| `command`                       | Override container command.                                                                                                                                                                           | `[]`                                                                 |
| `extraEnv`                      | Additional environment variables as key-value pairs.                                                                                                                                                  | `{}`                                                                 |
| `envValueFrom`                  | Additional environment variables from Kubernetes `valueFrom` references.                                                                                                                              | `{}`                                                                 |
| `envFrom`                       | Additional `envFrom` sources.                                                                                                                                                                         | `[]`                                                                 |
| `resources`                     | Resource requests and limits for the MCP container.                                                                                                                                                   | `{}`                                                                 |
| `securityContext`               | Pod security context. Empty value uses the subchart default security context.                                                                                                                         | `{}`                                                                 |
| `containerSecurityContext`      | Container security context.                                                                                                                                                                           | Drops all capabilities, runs as non-root user `1000`                 |
| `automountServiceAccountToken`  | Mount the Kubernetes ServiceAccount token in the MCP pod. Keep disabled unless Kubernetes API access is explicitly needed.                                                                            | `false`                                                              |
| `readinessProbe`                | Readiness probe configuration.                                                                                                                                                                        | HTTP GET `/healthz`                                                  |
| `livenessProbe`                 | Liveness probe configuration.                                                                                                                                                                         | HTTP GET `/healthz`                                                  |
| `startupProbe`                  | Optional startup probe.                                                                                                                                                                               | `{}`                                                                 |
| `nodeSelector`                  | Node selector for the MCP pod.                                                                                                                                                                        | `{}`                                                                 |
| `affinity`                      | Pod affinity rules.                                                                                                                                                                                   | `{}`                                                                 |
| `tolerations`                   | Pod tolerations.                                                                                                                                                                                      | `[]`                                                                 |
| `topologySpreadConstraints`     | Topology spread constraints.                                                                                                                                                                          | `[]`                                                                 |
| `priorityClassName`             | PriorityClass name.                                                                                                                                                                                   | `""`                                                                 |
| `runtimeClassName`              | RuntimeClass name.                                                                                                                                                                                    | `""`                                                                 |
| `schedulerName`                 | Scheduler name.                                                                                                                                                                                       | `""`                                                                 |
| `terminationGracePeriodSeconds` | Pod termination grace period.                                                                                                                                                                         | `""`                                                                 |
| `hostAliases`                   | Pod host aliases.                                                                                                                                                                                     | `[]`                                                                 |
| `dnsPolicy`                     | Pod DNS policy.                                                                                                                                                                                       | `""`                                                                 |
| `dnsConfig`                     | Pod DNS config.                                                                                                                                                                                       | `{}`                                                                 |
| `lifecycle`                     | Container lifecycle hooks.                                                                                                                                                                            | `{}`                                                                 |
| `initContainers`                | Extra init containers rendered as YAML.                                                                                                                                                               | `[]`                                                                 |
| `extraInitContainers`           | Extra init containers rendered through `tpl`.                                                                                                                                                         | `[]`                                                                 |
| `extraContainers`               | Extra sidecar containers.                                                                                                                                                                             | `[]`                                                                 |
| `volumes`                       | Extra pod volumes.                                                                                                                                                                                    | `[]`                                                                 |
| `volumeMounts`                  | Extra volume mounts for the MCP container.                                                                                                                                                            | `[]`                                                                 |
| `labels`                        | Extra labels added to MCP Deployment and pod template.                                                                                                                                                | `{}`                                                                 |
| `annotations`                   | Extra annotations added to MCP Deployment.                                                                                                                                                            | `{}`                                                                 |
| `podAnnotations`                | Extra annotations added to MCP pods.                                                                                                                                                                  | `{}`                                                                 |
| `podLabels`                     | Extra labels added to MCP pods.                                                                                                                                                                       | `{}`                                                                 |
<!-- markdownlint-enable line-length -->

`existingSecret.name` and `serviceAccountToken` are mutually exclusive. Use
the existing Secret for production credentials; Helm fails rendering if both
sources are configured.

Authentication must be considered separately for two connections:

* From `mcp-grafana` to Grafana: configure static credentials through
  `existingSecret.name`, `serviceAccountToken`, or basic authentication; or use
  request-time credentials with `GRAFANA_FORWARD_HEADERS` as described below.
  Credentials are not required only when Grafana actually permits anonymous
  access to the requested APIs.
* From an MCP client to `mcp-grafana`: the MCP server does not provide its own
  inbound authentication. Protect its HTTPRoute or Ingress at the API Gateway,
  ingress controller, or service mesh layer, for example through integration
  with an identity provider. This protection is also required when MCP uses
  static credentials to access Grafana.

When request-time credentials are used, the MCP client must send them with each
request and every proxy in the path must preserve the forwarded headers.

The chart is intended to expose the MCP server through HTTPRoute or Ingress.
Their configured hosts are added to `--allowed-hosts` automatically. For a
non-standard proxy that rewrites the `Host` header, or for direct access through
the Kubernetes Service, add the resulting host (including its port when present)
to `allowedHosts` explicitly.

Instead of configuring a static service account token or username/password, an
HTTP MCP client can provide credentials with each request. Configure
`mcp-grafana` to forward selected incoming headers to Grafana through
`extraEnv`:

```yaml
grafana:
  mcp:
    extraEnv:
      GRAFANA_FORWARD_HEADERS: "Authorization"
```

Multiple header names can be specified as a comma-separated list, for example
`Authorization,Cookie,X-Session-Id`. This approach has the following
limitations:

* It is supported only with `streamable-http` or `sse` transport and has no
  effect in `stdio` mode.
* The MCP client must support custom HTTP headers and send the required header
  with each request.
* `GRAFANA_FORWARD_HEADERS` is an allow-list of header names; it does not
  contain or store their values.
* Forwarding headers does not authenticate or authorize access to the MCP
  endpoint itself.
* If an API Gateway authenticates requests to the MCP endpoint, verify that it
  preserves the forwarded headers.
* Avoid combining a forwarded `Authorization` header with static Grafana
  credentials unless the resulting precedence has been verified. Incoming
  forwarded headers take precedence over headers with the same names configured
  through `GRAFANA_EXTRA_HEADERS`.

## Service

All parameters listed below are configured under the `grafana.mcp.service`
level.

<!-- markdownlint-disable line-length -->
| Parameter                  | Description                                    | Default     |
| -------------------------- | ---------------------------------------------- | ----------- |
| `type`                     | Kubernetes Service type.                       | `ClusterIP` |
| `port`                     | Service port.                                  | `8000`      |
| `targetPort`               | Container listen port and Service target port. | `8000`      |
| `annotations`              | Extra Service annotations.                     | `{}`        |
| `labels`                   | Extra Service labels.                          | `{}`        |
| `clusterIP`                | Static ClusterIP.                              | `""`        |
| `externalIPs`              | External IPs.                                  | `[]`        |
| `loadBalancerIP`           | LoadBalancer IP.                               | `""`        |
| `loadBalancerSourceRanges` | LoadBalancer source ranges.                    | `[]`        |
| `externalName`             | ExternalName target.                           | `""`        |
| `nodePort`                 | NodePort value when Service type requires it.  | `""`        |
| `extraPorts`               | Additional Service ports.                      | `[]`        |
| `sessionAffinity`          | Service session affinity.                      | `""`        |
| `sessionAffinityConfig`    | Service session affinity config.               | `{}`        |
<!-- markdownlint-enable line-length -->

## ServiceAccount

All parameters listed below are configured under the
`grafana.mcp.serviceAccount` level.

<!-- markdownlint-disable line-length -->
| Parameter                      | Description                                                    | Default |
| ------------------------------ | -------------------------------------------------------------- | ------- |
| `create`                       | Create a ServiceAccount for the MCP pod.                       | `true`  |
| `name`                         | ServiceAccount name. Empty value defaults to the MCP fullname. | `""`    |
| `annotations`                  | Extra ServiceAccount annotations.                              | `{}`    |
| `labels`                       | Extra ServiceAccount labels.                                   | `{}`    |
| `automountServiceAccountToken` | ServiceAccount-level `automountServiceAccountToken`.           | `false` |
<!-- markdownlint-enable line-length -->

## ServiceMonitor

The ServiceMonitor is rendered only when both
`metrics.enabled=true` and
`serviceMonitor.install=true`.

All parameters listed below are configured under the
`grafana.mcp.serviceMonitor` level.

<!-- markdownlint-disable line-length -->
| Parameter           | Description                              | Default    |
| ------------------- | ---------------------------------------- | ---------- |
| `install`           | Create a ServiceMonitor for MCP metrics. | `false`    |
| `labels`            | Extra ServiceMonitor labels.             | `{}`       |
| `annotations`       | Extra ServiceMonitor annotations.        | `{}`       |
| `path`              | Metrics path.                            | `/metrics` |
| `namespaceSelector` | ServiceMonitor namespace selector.       | `{}`       |
| `interval`          | Scrape interval.                         | `30s`      |
| `scrapeTimeout`     | Scrape timeout.                          | `10s`      |
| `metricRelabelings` | Metric relabeling rules.                 | `[]`       |
| `relabelings`       | Target relabeling rules.                 | `[]`       |
<!-- markdownlint-enable line-length -->

## HTTPRoute

All parameters listed below are configured under the `grafana.mcp.httpRoute`
level.

<!-- markdownlint-disable line-length -->
| Parameter         | Description                                                                                                | Default                        |
| ----------------- | ---------------------------------------------------------------------------------------------------------- | ------------------------------ |
| `install`         | Create a Gateway API HTTPRoute for MCP.                                                                    | `false`                        |
| `apiVersion`      | HTTPRoute API version.                                                                                     | `gateway.networking.k8s.io/v1` |
| `kind`            | Route kind.                                                                                                | `HTTPRoute`                    |
| `hostnames`       | HTTPRoute hostnames. Values are rendered through `tpl`.                                                    | `[]`                           |
| `parentRefs`      | Gateway parent references. At least one entry is required when `install=true`.                             | `[]`                           |
| `rules`           | Full custom HTTPRoute rules. If set, these rules are used instead of `matches`, `filters`, and `timeouts`. | `[]`                           |
| `matches`         | Default route matches.                                                                                     | PathPrefix `/`                 |
| `timeouts`        | HTTPRoute timeouts.                                                                                        | `{}`                           |
| `filters`         | HTTPRoute filters.                                                                                         | `[]`                           |
| `additionalRules` | Additional rules prepended before the default rule.                                                        | `[]`                           |
| `httpsRedirect`   | Render an HTTPS redirect rule instead of the backend route.                                                | `false`                        |
| `annotations`     | Extra HTTPRoute annotations.                                                                               | `{}`                           |
| `labels`          | Extra HTTPRoute labels.                                                                                    | `{}`                           |
<!-- markdownlint-enable line-length -->

## Ingress

All parameters listed below are configured under the `grafana.mcp.ingress`
level.

<!-- markdownlint-disable line-length -->
| Parameter     | Description                                                                                                                                                            | Default |
| ------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------- |
| `install`     | Create an Ingress for MCP.                                                                                                                                             | `false` |
| `className`   | IngressClass name.                                                                                                                                                     | `""`    |
| `annotations` | Extra Ingress annotations.                                                                                                                                             | `{}`    |
| `labels`      | Extra Ingress labels.                                                                                                                                                  | `{}`    |
| `tls`         | Ingress TLS configuration.                                                                                                                                             | `[]`    |
| `hosts`       | Ingress hosts and paths. At least one entry is required when `install=true`; `host` values are rendered through `tpl`, and missing paths default to `/` with `Prefix`. | `[]`    |
<!-- markdownlint-enable line-length -->

## Example

```yaml
grafana:
  mcp:
    install: true
    existingSecret:
      name: grafana-mcp-token
      key: token
    grafanaUrl: http://grafana-service.monitoring.svc:3000
    metrics:
      enabled: true
    serviceMonitor:
      install: true
    httpRoute:
      install: true
      hostnames:
        - mcp-grafana.{{ .Release.Namespace }}.example.com
      parentRefs:
        - name: monitoring-gateway
```
