# Etcd certificate synchronization security constraints

The `etcd-certs-to-secret` workload copies the certificates required for etcd monitoring into the monitoring namespace.
Its certificate source and security model depend on the platform.

| Platform | Certificate source | Runs as root | Uses `hostPath` |
| --- | --- | --- | --- |
| Kubernetes | Files on a control-plane node | Yes | Yes, read-only `/etc` |
| OpenShift | A ConfigMap and Secret in `openshift-etcd-operator` | No | No |

## Kubernetes constraints

On Kubernetes, the workload finds a running etcd Pod in the `kube-system` namespace and reads certificate paths from
the `etcd` container's `command` and `args` fields. It recognizes these flags:

- `--peer-key-file`
- `--peer-trusted-ca-file`
- `--peer-cert-file`

If the flags are absent, the workload uses these paths:

```text
/etc/kubernetes/pki/etcd/peer.key
/etc/kubernetes/pki/etcd/ca.crt
/etc/kubernetes/pki/etcd/peer.crt
```

Certificate locations vary between Kubernetes distributions. The Helm chart therefore mounts the control-plane
node's entire `/etc` directory into the container at `/etc`. The mount is read-only.

The workload runs as UID and GID `0` because etcd private keys are commonly readable only by `root`. Kubernetes does
not apply `fsGroup` ownership changes to these host files, so a non-root process cannot reliably read them.

These requirements prevent the Kubernetes workload from complying with the following hardening controls:

- `runAsNonRoot: true`
- The prohibition of `hostPath` volumes
- The Kubernetes Pod Security Standards `Baseline` and `Restricted` profiles

Pod Security Admission cannot grant an exception to one Pod through workload labels. A cluster that enforces
`Baseline` or `Restricted` for the monitoring namespace must configure an admission exemption for this workload or
run it in a namespace with an appropriate policy. Prefer an exemption scoped to the workload's ServiceAccount over a
namespace-wide exemption when the admission implementation supports it.

The Kubernetes templates apply the following compensating controls:

- Schedule the workload on control-plane nodes.
- Mount host `/etc` read-only.
- Disable privilege escalation.
- Drop all Linux capabilities.
- Use the runtime-default seccomp profile.
- Mount the container root filesystem read-only.
- Provide a size-limited `emptyDir` at `/tmp`.

The ServiceAccount token must remain mounted. The workload uses the Kubernetes API to discover the cluster type, list
etcd Pods, create or update the destination Secret, and create or update the etcd Service.

The discovered certificate paths must be under `/etc`. A distribution that configures etcd certificates outside
`/etc` requires an additional read-only host mount and a corresponding path mapping before the workload can read the
files.

## OpenShift behavior

On OpenShift, the workload reads certificates through the Kubernetes API instead of the host filesystem:

- CA certificate: ConfigMap `openshift-etcd-operator/etcd-metric-serving-ca`
- Client key and certificate: Secret `openshift-etcd-operator/etcd-metric-client`

The OpenShift templates do not mount host `/etc`. They run the workload as a non-root user selected by OpenShift and
apply these controls:

- `runAsNonRoot: true`
- `allowPrivilegeEscalation: false`
- `readOnlyRootFilesystem: true`
- All Linux capabilities dropped
- Runtime-default seccomp
- Only ConfigMap, downward API, `emptyDir`, projected, and Secret volumes allowed by the dedicated SCC

The chart detects OpenShift through the `security.openshift.io/v1/SecurityContextConstraints` API.

## Residual Kubernetes risk

A process compromise in the Kubernetes workload provides read access to files under the control-plane node's `/etc`
directory. The read-only mount prevents modification but does not prevent disclosure. Keep the image, ServiceAccount,
RBAC rules, and admission exception narrowly scoped, and review them whenever the certificate discovery logic changes.
