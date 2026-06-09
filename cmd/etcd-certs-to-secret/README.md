# etcd-certs-to-secret

A small one-shot binary that discovers etcd client certificates in a Kubernetes
or OpenShift v4 cluster and writes them to a Secret so Prometheus can scrape
etcd over TLS.

It is run by the Helm chart as:

- a **Job** with `helm.sh/hook: post-install,post-upgrade` — runs once on
  install/upgrade to seed the Secret.
- a **CronJob** — re-runs on a schedule to refresh the Secret when the
  underlying etcd certificates rotate.

The ServiceMonitor that scrapes etcd is rendered directly by the Helm chart
([`templates/etcd-certs-to-secret-job/servicemonitor.yaml`](../../charts/qubership-monitoring-operator/templates/etcd-certs-to-secret-job/servicemonitor.yaml)),
not by this binary.

## What it does

On each run:

1. Detects the cluster flavor by probing for the `security.openshift.io/v1`
   API group (present ⇒ OpenShift v4, absent ⇒ plain Kubernetes).
2. Loads etcd client certs:
   - **OpenShift v4**: reads ConfigMap `etcd-metric-serving-ca` (key
     `ca-bundle.crt`) and Secret `etcd-metric-client` (keys `tls.crt` /
     `tls.key`) from namespace `openshift-etcd-operator`.
   - **Kubernetes**: locates a running etcd pod (label `component=etcd` in
     `kube-system`), extracts cert paths from its `--peer-key-file`,
     `--peer-trusted-ca-file`, `--peer-cert-file` arguments (falling back
     to `/etc/kubernetes/pki/etcd/{peer.key,ca.crt,peer.crt}`), and reads
     the files from the hostPath-mounted `/etc`.
3. Verifies that each PEM block has the expected headers.
4. Creates/updates a Secret (default name `kube-etcd-client-certs`) in the
   target namespace with keys `etcd-client.key`, `etcd-client-ca.crt`,
   `etcd-client.crt`.
5. **Kubernetes only**: creates/updates a headless `Service` named `etcd` in
   `kube-system` selecting `component=etcd`, port `2379/metrics`, so the
   chart's ServiceMonitor has something to target. Skipped on OpenShift v4
   (the `openshift-etcd/etcd` Service already exists).

The binary exits after one pass. Periodic refresh is the CronJob's job.

## Flags and environment

| Flag           | Env         | Default                  | Description                                          |
|----------------|-------------|--------------------------|------------------------------------------------------|
| `--secret`     | —           | `kube-etcd-client-certs` | Name of the Secret to create/update.                 |
| `--namespace`  | `NAMESPACE` | `monitoring`             | Namespace in which to create the Secret.             |
| `--log-level`  | —           | `info`                   | `debug` / `info` / `warn` / `error`.                 |

The Job/CronJob templates inject `NAMESPACE` via the downward API
(`metadata.namespace`).

## Permissions

See the chart's
[`role.yaml`](../../charts/qubership-monitoring-operator/templates/etcd-certs-to-secret-job/role.yaml)
and
[`clusterrole.yaml`](../../charts/qubership-monitoring-operator/templates/etcd-certs-to-secret-job/clusterrole.yaml).
The binary needs:

- Namespace-scoped: `secrets` create/update/patch/get/list/watch; `pods`
  get/list.
- Cluster-scoped: `services` create/get/list/update/watch (kube-system only
  in practice); `pods` get/list/watch; `configmaps` get/list/watch on
  `etcd-metric-serving-ca` / `etcd-serving-ca`; `secrets` get/list/watch on
  `etcd-metric-client` / `etcd-client`.

On OpenShift it additionally needs `use` / `update` on the chart-supplied
SecurityContextConstraint (the Job mounts `/etc` as hostPath read-only).

## Build

The binary is a sub-package of the repository root module and uses the root
`go.mod`. From the repository root:

```sh
go build ./cmd/etcd-certs-to-secret
go test  ./cmd/etcd-certs-to-secret
```

The container image is built from the repository root as build context:

```sh
docker build -f cmd/etcd-certs-to-secret/Dockerfile -t etcd-certs-to-secret:dev .
```

The CI workflow does the same — see
[`.github/workflows/build.yaml`](../../.github/workflows/build.yaml).
