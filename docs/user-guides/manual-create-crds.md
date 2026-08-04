# Manage CRDs Manually

## When is it needed?

This document describes how to manually create, update, or delete Monitoring Custom Resource Definitions (CRDs).

Almost all Qubership applications and microservices are integrated with Monitoring. This integration means
that almost all microservices during deploy can create `ServiceMonitor`/`PodMonitor`/`Probe` objects from the API
`monitoring.coreos.com`.

Such objects allow for Monitoring to understand from which microservices and how to collect metrics.
Similar objects allow providing alerts, recording rules and Grafana dashboards.

Also, it means that before deploying applications and microservices in Kubernetes all Custom Resource Definitions
(CRDs) must already created in Kubernetes. Otherwise, the deployment will fail.

Helm installs CRDs during a fresh Monitoring installation, but Helm does not upgrade existing CRDs. Apply the complete
CRD set manually before every Monitoring upgrade so existing CRDs are updated and newly introduced CRDs are created.
The same procedure can install CRDs in environments where Monitoring itself is not deployed.

In the list of Monitoring artifacts, you can find an archive with the name

```bash
monitoring-operator-<version>-crds.zip
```

that contains all CRDs required for Monitoring. How to use this archive and CRDs inside you can read below.

Sources: [project][p], [Prometheus][po], [VictoriaMetrics][vm], [Grafana][g], [Adapter][a]; [standalone][s] is generated.

[p]: ../../charts/qubership-monitoring-operator/crds/
[po]: ../../charts/qubership-monitoring-operator/charts/prometheus-operator/crds/
[vm]: ../../charts/qubership-monitoring-operator/charts/victoriametrics-operator/crds/
[g]: ../../charts/qubership-monitoring-operator/charts/grafana-operator/crds/
[a]: ../../charts/qubership-monitoring-operator/charts/prometheus-adapter-operator/crds/
[s]: ../../charts/qubership-monitoring-crds/crds/

## Before you begin

* You should have cluster-wide permissions enough to operate with CRDs (cluster admin is not required)
* You should configure a context for your `kubectl` and make sure the connection configured to correct Kubernetes

## How to manage CRDs

This section describes different cases of manual manipulation with CRDs.

### Create

To create CRDs for Monitoring you need to execute the command:

```bash
kubectl apply --server-side --recursive -f path/to/crds/directory/
```

The `--server-side` option avoids the large `kubectl.kubernetes.io/last-applied-configuration` annotation created by
client-side apply.

```bash
mkdir /tmp/crds/
unzip -d /tmp/crds/ monitoring-operator-<version>-crds.zip
kubectl apply --server-side --recursive -f /tmp/crds/
```

### Upgrade

To upgrade CRDs for Monitoring you need to execute the command:

```bash
kubectl apply --server-side --force-conflicts --recursive -f path/to/crds/directory/
```

Apply the complete CRD set before upgrading controller images. Unlike `kubectl replace`, server-side apply creates CRDs
introduced by a newer release and updates CRDs that already exist. `--force-conflicts` transfers fields previously
managed by Helm to this server-side apply operation.

For example, if you will use the archive with CRDs providing Monitoring:

```bash
mkdir /tmp/crds/
unzip -d /tmp/crds/ monitoring-operator-<version>-crds.zip
kubectl apply --server-side --force-conflicts --recursive -f /tmp/crds/
```

### Remove

**Warning!** This step removes **all CRDs** for Monitoring and deleting these CRDs causes the deletion of
**all resources** of their type in **all namespaces**.
It means that all resources like `ServiceMonitor` and `GrafanaDashboard` from applications are removed.

To remove CRDs and all Custom Resources (CRs) like ServiceMonitor or GrafanaDashboards from all namespaces
for Monitoring you need to execute the command:

```bash
kubectl delete --recursive -f path/to/extracted/crds/
```
