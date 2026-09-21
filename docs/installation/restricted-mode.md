# Restricted mode installation

Set `global.privilegedRights: false` when a cluster administrator owns cluster-scoped RBAC and the
Monitoring Operator must not create or manage it. The setting controls which RBAC resources the chart and the
operator create. It does not restrict the Helm identity that installs or upgrades the release.

## Install prerequisites

The cluster administrator installs CRDs and the required cluster RBAC. The Helm identity installs the release.

Install or update the CRDs:

```bash
kubectl apply --server-side --force-conflicts -f charts/qubership-monitoring-crds/crds/
```

The default profile requires administrator-owned RBAC for kube-state-metrics and VMAgent. Apply both manifests before
installing the release:

```bash
export MONITORING_NAMESPACE=monitoring

envsubst < docs/installation/restricted-rbac/kube-state-metrics.yaml | kubectl apply -f -
envsubst < docs/installation/restricted-rbac/vmagent.yaml | kubectl apply -f -
```

The manifests bind the cluster roles to their component service accounts in `MONITORING_NAMESPACE`. Keep them aligned
with the chart version. When enabling components outside the default profile, have the cluster administrator review
their RBAC requirements first.

With the default empty `publicCloudName`, Helm also creates the etcd certificate Job's `ClusterRole` and
`ClusterRoleBinding`. The Helm identity must have permission to create those resources.

On OpenShift, `vmagent.yaml` grants VMAgent permission to use the `victoriametrics-operator` SCC.

## Install a restricted release

Create a values file that enables restricted mode:

```yaml
global:
  privilegedRights: false
```

Install the chart after the administrator has installed the CRDs and selected exception RBAC:

```bash
helm install monitoring-operator charts/qubership-monitoring-operator \
  --skip-crds \
  --namespace monitoring \
  --create-namespace \
  --values restricted-values.yaml
```

Use the same chart version for the CRDs, administrator-provided RBAC, and Helm release.

## Transition an existing release

1. Apply the required administrator-owned RBAC before changing the release.
2. Set `global.privilegedRights: false` in the release values.
3. Upgrade the release:

   ```bash
   helm upgrade monitoring-operator charts/qubership-monitoring-operator \
     --skip-crds \
     --namespace monitoring \
     --values restricted-values.yaml
   ```

4. Verify the installation, then remove obsolete cluster RBAC.

## Verify the installation

Confirm that the component service accounts have the required access:

```bash
kubectl auth can-i --as=system:serviceaccount:monitoring:monitoring-kube-state-metrics list namespaces
kubectl auth can-i --as=system:serviceaccount:monitoring:monitoring-kube-state-metrics list pods --all-namespaces
kubectl auth can-i --as=system:serviceaccount:monitoring:monitoring-vmagent get nodes
kubectl auth can-i --as=system:serviceaccount:monitoring:monitoring-vmagent list pods --all-namespaces
kubectl auth can-i --as=system:serviceaccount:monitoring:monitoring-vmagent list endpoints --all-namespaces
```

Every command must return `yes`. Confirm that the workloads are ready and their logs do not contain authorization
errors:

```bash
kubectl get pods --namespace monitoring
kubectl logs deployment/kube-state-metrics --namespace monitoring
kubectl logs deployment/vmagent-k8s --namespace monitoring
```
