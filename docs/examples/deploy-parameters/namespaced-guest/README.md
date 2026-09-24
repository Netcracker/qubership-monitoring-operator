# Namespaced guest install

Two additive overlays on chart defaults. Host and guest are the same
monitoring-operator product.

|File|Release|What it changes|
|----|-------|---------------|
|[host-values.yaml](host-values.yaml)|Host (cluster-wide)|Host VMAgent / VMAlert / VMAlertmanager skip the guest namespace. VM and Grafana operators stay cluster-wide.|
|[values.yaml](values.yaml)|Guest|`namespaceScope: true`, unique etcd RBAC/SCC names, and `nodeExporter.port: 9901`|

Guest controllers that manage CRs stay in the release namespace. Guest
workloads may still discover host resources. Host VM Operator may still
reconcile guest VM CRs (it has no namespaceSelector). Host scrape/rules do
not take guest ServiceMonitors.

## Contract

|Side|Control|Effect|
|----|-------|------|
|Guest Helm|`global.namespaceScope: true`|Guest VM Operator and Grafana Operator `WATCH_NAMESPACE` is `.Release.Namespace`|
|Guest Helm|`etcdCertsJob.rbac.*Name`|Unique ClusterRole, ClusterRoleBinding, and OpenShift SCC names|
|Core operator|Pod-derived `WATCH_NAMESPACE`|Already namespaced|
|Host VM / Grafana operators|Chart default (empty)|Cluster-wide|
|Host VMAgent / VMAlert / VMAlertmanager|`NotIn` guest namespace name|Discover all namespaces except the guest|
|Guest Helm|`--skip-crds`|Host owns shared CRDs|

Leader election stays at the default `false`. Change `monitoring-test` in
`host-values.yaml` if the guest namespace is different.

## Temporary coexistence limitation

This example is for a short-lived guest installation beside a persistent
cluster-wide host, for example a CI or overlay scenario. It is not a
long-lived multi-tenant topology.

The host VictoriaMetrics Operator has an empty `WATCH_NAMESPACE`, so it
observes the guest `VMAgent`, `VMAlert`, `VMAlertmanager`, and other VM custom
resources. Namespace selectors on those resources do not control which
VictoriaMetrics Operator instance reconciles them; they only select
configuration inputs:

- `VMAgent` selectors choose `VMServiceScrape` and `VMPodScrape` objects;
- `VMAlert` selectors choose `VMRule` objects;
- `VMAlertmanager` selectors choose `VMAlertmanagerConfig` objects.

The host overlay applies selectors to the *host* workload CRs, so its
VMAgent, VMAlert, and VMAlertmanager omit configuration inputs from the guest
namespace. That prevents normal host consumption of the guest's scrape jobs,
rules, and alert-routing configuration.

It does not prevent the host VictoriaMetrics Operator from attempting to
reconcile the guest VM workload CRs. The host service account has broad
cluster-scoped access, including to guest VM CRs, Secrets, and ConfigMaps, but
does not have create or update permission for the guest Deployments and
StatefulSets. Those workload writes are denied and appear as reconciliation
errors in the host operator logs. The guest VictoriaMetrics Operator, whose
`WATCH_NAMESPACE` is the guest namespace, reconciles the guest workloads
successfully. The errors are expected and time-bounded: they stop when the
guest release is uninstalled.

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
|`vmAgent` / `vmAlert` / `vmAlertManager` `NotIn monitoring-test`|host|Host scrape and rules skip the guest namespace|
|`global.namespaceScope: true`|guest|Guest VM and Grafana operators watch only the release namespace|
|`etcdCertsJob.rbac.*Name`|guest|Avoid collisions with the host etcd certificate ClusterRole, ClusterRoleBinding, and OpenShift SCC|
|`nodeExporter.port: 9901`|guest|Avoid hostPort 9900 collision with the host node-exporter|
