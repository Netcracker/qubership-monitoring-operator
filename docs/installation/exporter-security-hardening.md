# Exporter Security Hardening

This document describes the security settings applied to exporter workloads and the cases where an exporter cannot meet
the complete hardening baseline without losing its primary function.

## Applied baseline

The following settings are enforced in the Helm templates for exporter containers, unless an exception is listed below:

- `allowPrivilegeEscalation: false`
- `readOnlyRootFilesystem: true`
- all Linux capabilities are dropped
- `runAsNonRoot: true` at the pod level
- `seccompProfile.type: RuntimeDefault`
- a writable `/tmp` backed by an `emptyDir` with a `100Mi` size limit

On Kubernetes, charts supply a numeric non-root user and group. On OpenShift, numeric user and group settings are
omitted so that the active Security Context Constraints policy can assign IDs from the namespace range. User-supplied
security context fields are preserved, but they cannot weaken the enforced baseline fields.

The baseline is applied to these chart workloads:

| Exporter | Workload | Baseline status |
| --- | --- | --- |
| Blackbox Exporter | Deployment or DaemonSet | Applied |
| Cert Exporter | Deployment | Applied |
| Cert Exporter | DaemonSet | Applied with host filesystem limitations |
| Cloud Events Exporter | Deployment | Applied |
| CloudWatch Exporter | Deployment | Applied with credential-file considerations |
| Goldpinger Exporter | DaemonSet | Applied unless host networking is explicitly enabled |
| JSON Exporter | Deployment | Applied; additional volumes remain user-controlled |
| Network Latency Exporter | DaemonSet | Applied unless privileged mode is explicitly enabled |
| Promitor Resource Discovery | Deployment | Applied; additional volumes remain user-controlled |
| Promitor Scraper | Deployment | Applied |
| SSL Exporter | DaemonSet | Applied with host filesystem limitations |
| Stackdriver Exporter | Deployment | Applied |
| Version Exporter | Deployment | Applied |

Node Exporter, Kube State Metrics, and Pushgateway are covered by their component-specific hardening implementation and
documentation.

## Known limitations

### Cert Exporter DaemonSet

The DaemonSet inspects certificates and kubeconfig files stored on each node. It therefore uses read-only `hostPath`
volumes for paths such as `/etc/kubernetes`, `/var/lib/kubelet/pki`, `/etc/origin`, `/etc/etcd`, and
`/root/.kube/config`. Additional paths can also be configured through `additionalHostPathVolumes`.

This mode cannot comply with a policy that forbids all `hostPath` volumes. The mounts remain read-only, the container
baseline is enforced, and the chart's legacy PSP and SCC definitions are restricted to the capabilities and volume types
that the workload uses. Prefer the Deployment and secret-scanning mode when node-local certificate inspection is not
required.

The chart retains the legacy `RunAsAny` strategy in the Cert Exporter SCC. It does not encode the UID range that
OpenShift allocates to a namespace at runtime. `RunAsAny` is broader than the common hardening baseline and is a
documented exception. The pod template still requests a non-root user. The upstream image declares a nonnumeric user
named `app`; some runtime combinations may reject the pod because the kubelet cannot verify that this user is non-root.
Validate admission with the exact image and OpenShift version used by the deployment.

Even when the pod is admitted, the assigned non-root process must have read permission for every selected host file.
Kubernetes distributions commonly protect private keys and root kubeconfig files with root ownership and mode `0600`.
This includes some etcd certificates. Such files will be skipped or cause read errors.

If node-local certificate inspection is required, a cluster administrator must provide a distribution-specific,
narrowly scoped configuration and make only the required files readable by the selected service account. Reading
root-owned `0600` files requires either safe permission changes or an approved exception to the pod security context
and SCC; the chart does not enable that exception. Do not grant broad access to node certificate directories. Prefer
the Deployment and secret-scanning mode when the node files cannot be exposed safely.

### Network Latency Exporter

The exporter performs ICMP and MTR-style network probes. The current image runs as non-root UID `2001`. ICMP, UDP, and
TCP MTR probes were validated with all capabilities dropped, privilege escalation disabled, and a read-only root
filesystem. The default template therefore applies the complete baseline.

Setting `rbac.privileged=true` remains a backward-compatible, explicit escape hatch. It should be used only when the
runtime, kernel, or a replacement image cannot perform the configured probes under the default restrictions. Privileged
mode does not comply with the hardening baseline and requires an administrator-provided OpenShift SCC. Numeric IDs are
not emitted by the normal OpenShift template, and the assigned UID remains non-root.

### SSL Exporter DaemonSet

SSL Exporter supports HTTPS, Kubernetes Secret, file, and kubeconfig probes. Only the file and kubeconfig probes need
access to the host filesystem. The default values mount `/etc/ssl/cert.pem` and `/etc/ssl/certs` read-only. These paths
usually contain public CA certificates that a non-root process can read.

The default mounts use `hostPath`, which violates policies that forbid host filesystem access. A restricted OpenShift
SCC does not allow these volumes, and the SSL Exporter chart does not create a component-specific SCC. A cluster
administrator must provide a narrowly scoped SCC when host file inspection is required. Set
`sslExporter.additionalHostPathVolumes=[]` when only HTTPS or Kubernetes Secret probes are needed.

Every selected host file must be readable by the UID assigned to the pod. Private keys, kubeconfig files, and some
node-local certificates may be root-owned with mode `0600`. The enforced `runAsNonRoot: true` setting prevents SSL
Exporter from reading those files. `RunAsAny` in an SCC would permit a requested UID, but it would not override the pod
security context or grant filesystem permissions.

Reading root-owned `0600` files requires safe permission changes or an approved, deployment-specific exception to the
pod security context and SCC. The chart cannot apply a portable exception because file ownership and OpenShift UID
ranges depend on the target environment. Do not grant broad access to host certificate directories.

### Goldpinger host-network mode

Goldpinger uses normal pod networking by default. Adding `USE_HOST_IP=true` to `extraEnv` enables `hostNetwork` and a
`hostPort` for backward compatibility. That opt-in mode does not comply with the restricted baseline. Keep it disabled
unless the deployment specifically requires node-address-based discovery.

### User-supplied volumes

JSON Exporter and Promitor Resource Discovery accept arbitrary extra volumes and mounts. These extension points are kept
for compatibility and are not rewritten by the templates. Do not supply `hostPath` or other privileged volume types.
Clusters that require enforcement should use Pod Security Admission and an admission policy to reject unsafe additions.

### CloudWatch credential file

Static AWS credentials are exposed as a read-only Secret volume at `/root/.aws/credentials` for compatibility with the
AWS SDK lookup path. The non-root process must be able to traverse the mounted directory and read the `0400` Secret
file. Validate this mode on the target Kubernetes or OpenShift version. Prefer workload identity, IRSA, or an instance
profile when available.

## Image boundary

Exporter images are external dependencies and do not have Dockerfiles in this repository. The Helm chart can enforce
runtime restrictions, but it cannot guarantee that an upstream image works with a read-only filesystem, an arbitrary
OpenShift UID, or a non-root Kubernetes UID. Test the exact pinned images before release and revalidate them after every
image upgrade.

The `100Mi` `/tmp` limit is a conservative default because exporter-specific peak temporary storage requirements are not
known. Monitor ephemeral-storage use under production-like load and change the limit only when measurements justify it.

## Validation checklist

For every enabled exporter:

1. Render the chart for both Kubernetes and OpenShift API capabilities.
2. Verify the effective pod and container security contexts in the rendered manifest.
3. Confirm that the pod reaches `Ready` without permission, read-only filesystem, or seccomp errors.
4. Inspect logs for failures to read configuration, credentials, certificates, or kubeconfig files.
5. Scrape the metrics endpoint and exercise each configured probe type.
6. Confirm that `/tmp` usage stays below `100Mi` during representative load.
7. For DaemonSets, verify readiness and logs on every supported node role and operating system.
