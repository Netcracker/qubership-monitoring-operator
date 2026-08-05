# Optional Component Security Hardening

This document describes the security baseline and known limitations for Grafana Image Renderer, Graphite Remote
Adapter, Prometheus Adapter, Promxy, and the VictoriaMetrics Operator cleanup hook.

## Common baseline

The chart applies the following controls to the workloads that it creates directly:

- `runAsNonRoot: true`
- `seccompProfile.type: RuntimeDefault`
- numeric user and group IDs on Kubernetes, while OpenShift assigns IDs from the namespace range
- `allowPrivilegeEscalation: false`
- `readOnlyRootFilesystem: true`
- all Linux capabilities dropped
- a writable `/tmp` backed by an `emptyDir` with a `100Mi` size limit

The mandatory controls override weaker values supplied by a user. The numeric IDs are intentionally omitted on
OpenShift so that the namespace Security Context Constraints allocation remains authoritative.

## Grafana Image Renderer

The renderer runs Chromium and creates temporary browser profiles and rendered output. It therefore requires a writable
temporary directory. The chart mounts the size-limited `emptyDir` at `/tmp`; the root filesystem remains read-only.

Rendering large dashboards can exhaust the `100Mi` temporary-storage limit. Changing that limit requires a reviewed
chart change after measuring the workload and accepting the additional ephemeral-storage consumption. The renderer does
not require host access, privileged mode, or added Linux capabilities.

The v5 renderer disables Chromium's internal sandbox by default. Enabling the browser sandbox requires a separate
runtime test because Chromium may need a setuid sandbox helper or user-namespace support that is unavailable under the
restricted container settings. The pod remains isolated by Kubernetes controls, including non-root execution, seccomp,
dropped capabilities, and disabled privilege escalation.

Grafana reads the remote-renderer URL from a ConfigMap through `envFrom`. Kubernetes does not update environment
variables in an existing Pod. When the renderer is enabled after Grafana is already running, recreate the Grafana Pod
after the ConfigMap is updated. A fresh installation does not require this extra step.

## Graphite Remote Adapter

The adapter reads its configuration from a ConfigMap and keeps its metric-path cache in memory. Its configuration mount
is read-only from the application's perspective, and `/tmp` is the only writable filesystem provided by the chart. It
does not require host access, privileged mode, or added Linux capabilities.

The adapter still needs network access to the configured Graphite read and Carbon write endpoints. Container hardening
does not replace the NetworkPolicy rules required for those connections.

## Prometheus Adapter Operator and generated adapter

The Prometheus Adapter Operator Deployment is fully covered by the common baseline. The operator needs Kubernetes API
access to reconcile the `PrometheusAdapter` custom resource and its dependent resources.

The generated adapter has an upstream API limitation. The `PrometheusAdapter` CRD version shipped with this chart
supports only `runAsUser` and `fsGroup` under `spec.securityContext`. The upstream operator creates the adapter
Deployment without pod-level `runAsNonRoot` and `seccompProfile` fields and without a container security context. The
operator also creates an unbounded `/tmp` volume because the adapter generates its serving certificate under
`/tmp/cert`.

Consequently, this chart cannot apply the complete baseline to the generated adapter without changing the upstream CRD
and reconciliation code. Adding unsupported fields to the custom resource would make it fail schema validation. The
adapter image runs as a non-root user, and this chart continues to configure only the supported numeric user and group
fields. Full compliance requires an upstream operator version that exposes the missing pod, container, and volume
security settings.

## Promxy

Promxy and its config-reload sidecar use the common pod baseline and share the size-limited `/tmp` volume. The default
configuration proxies requests to remote Prometheus-compatible servers and does not enable local TSDB storage.

Do not add `--storage.tsdb.path` through `extraArgs` unless it points to a writable mounted path. This chart does not
currently expose a dedicated data volume for Promxy. Using a path on the read-only root filesystem prevents startup;
using `/tmp` makes the data ephemeral and subjects it to the `100Mi` limit. Persistent local storage requires a chart
extension that adds a dedicated volume and mount.

Promxy also needs network access to every configured server group. Any NetworkPolicy must allow those outbound
connections and the local reload sidecar request.

## VictoriaMetrics Operator pre-delete job

The pre-delete Job uses the complete common baseline on Kubernetes. On OpenShift, its rendered template omits
`runAsNonRoot` and numeric user and group IDs. The remaining controls include `RuntimeDefault`, a read-only root
filesystem, disabled privilege escalation, dropped capabilities, and a size-limited `/tmp` volume.

The default `rancher/kuberlr-kubectl` image declares a nonnumeric `kuberlr` user. Kubelet cannot verify that this user
is non-root before container startup. Setting `runAsNonRoot: true` can therefore prevent the Job from starting when an
OpenShift SCC does not assign a numeric user ID.

The OpenShift template requires the standard `restricted-v2` SCC. This prevents a more specific SCC from being
selected for the Job and makes OpenShift assign a numeric UID from the namespace range. The SCC can then add
`runAsNonRoot` safely because kubelet receives both the requirement and a numeric UID. The service account that creates
the hook must be allowed to use `restricted-v2`.

The hook is destructive by design. A normal installation or upgrade does not execute it; Helm runs it during release
deletion. Do not execute the hook as a security test because it scales down the monitoring operator and deletes all
matching VictoriaMetrics resources in the release namespace.
