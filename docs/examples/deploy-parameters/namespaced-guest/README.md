# Namespaced guest install

Two additive overlays on chart defaults. Host and guest are the same
monitoring-operator product.

|File|Release|What it changes|
|----|-------|---------------|
|[host-values.yaml](host-values.yaml)|Host (privileged)|Host VM Operator and Grafana Operator watch an explicit namespace list. Host VMAgent, VMAlert, and VMAlertmanager skip the guest namespace.|
|[values.yaml](values.yaml)|Guest|`namespaceScope: true`, unique etcd RBAC/SCC names, and `nodeExporter.port: 9901`|

Guest controllers that manage custom resources stay in the release namespace.
The host operators watch only the namespaces named in the host allow-list.
Guest workloads may still discover host resources when those workloads keep
the chart-default selectors.

## Contract

|Side|Control|Effect|
|----|-------|------|
|Guest Helm|`global.namespaceScope: true`|Guest VM Operator and Grafana Operator `WATCH_NAMESPACE` is `.Release.Namespace`|
|Guest Helm|`etcdCertsJob.rbac.*Name`|Unique ClusterRole, ClusterRoleBinding, and OpenShift SCC names|
|Core operator|Pod-derived `WATCH_NAMESPACE`|Already namespaced|
|Host Helm|`victoriametrics.vmOperator.extraEnvs` `WATCH_NAMESPACE`|Comma-separated namespaces the privileged host VM Operator watches|
|Host Helm|`grafana.operator.watchNamespaces`|The same list for the privileged host Grafana Operator|
|Host VMAgent / VMAlert / VMAlertmanager|`NotIn` guest namespace name|Host scrape, rules, and alert routing skip the guest namespace|
|Guest Helm|`--skip-crds`|Host owns shared CRDs|

Leader election stays at the default `false`. The example allow-list is
`monitoring`. Replace it with the namespaces that host release must manage,
and keep the guest namespace out of the list. Change `monitoring-test` in
`host-values.yaml` when the guest namespace is different.

## Host namespace allow-list

The host release stays privileged. Set one comma-separated list, with no
spaces, in both of these values:

- `victoriametrics.vmOperator.extraEnvs`, name `WATCH_NAMESPACE`
- `grafana.operator.watchNamespaces`

You own the list. Name each namespace the host operators watch. The chart
does not turn the list into every namespace except the guest. Omit the guest
namespace. When a namespace is added or removed, edit the list and upgrade
the host release. The host keeps its ClusterRoles. This mode does not create
a Role or RoleBinding in each listed namespace.

The example list is `monitoring`, the host namespace in the install commands
below.

Leave `WATCH_NAMESPACE` empty and the host VictoriaMetrics Operator watches
every namespace, including the guest. It then reconciles guest `VMAgent`,
`VMAlert`, `VMAlertmanager`, and other VM custom resources. Selectors on
those resources choose configuration inputs:

- `VMAgent` selectors choose `VMServiceScrape` and `VMPodScrape` objects;
- `VMAlert` selectors choose `VMRule` objects;
- `VMAlertmanager` selectors choose `VMAlertmanagerConfig` objects.

The host `NotIn` selectors apply to the host workload custom resources, so
host VMAgent, VMAlert, and VMAlertmanager omit guest scrape jobs, rules, and
alert-routing configuration. Those selectors do not choose which
VictoriaMetrics Operator reconciles a custom resource.

An empty `WATCH_NAMESPACE` leaves that choice with the cluster-wide host
operator. Its service account can read guest VM custom resources, Secrets,
and ConfigMaps, and it cannot create or update Deployments and StatefulSets
in the guest namespace. Those writes fail, and the denials stay in the host
operator logs until the guest release is uninstalled. Put the guest namespace
outside the allow-list so the host operator does not watch it.

## Shared CRDs

CRDs are cluster-wide APIs owned by the host. The guest **always** uses
`--skip-crds`. Small host/guest content differences are acceptable. A
material CRD or compatibility conflict is handled manually when it arises.

## Install

```bash
kubectl create namespace monitoring-test

helm upgrade --install monitoring-operator charts/qubership-monitoring-operator \
  --namespace monitoring \
  --skip-crds \
  --values docs/examples/deploy-parameters/namespaced-guest/host-values.yaml

helm install monitoring-operator-guest charts/qubership-monitoring-operator \
  --namespace monitoring-test \
  --skip-crds \
  --values docs/examples/deploy-parameters/namespaced-guest/values.yaml
```

|Setting|File|Why|
|-------|----|---|
|`WATCH_NAMESPACE` and `grafana.operator.watchNamespaces`|host|Host operators watch the deployment-owned namespace list and omit the guest|
|`vmAgent` / `vmAlert` / `vmAlertManager` `NotIn monitoring-test`|host|Host scrape, rules, and alert routing skip the guest namespace|
|`global.namespaceScope: true`|guest|Guest VM and Grafana operators watch only the release namespace|
|`etcdCertsJob.rbac.*Name`|guest|Avoid collisions with the host etcd certificate ClusterRole, ClusterRoleBinding, and OpenShift SCC|
|`nodeExporter.port: 9901`|guest|Avoid hostPort 9900 collision with the host node-exporter|
