# Restricted mode installation

Use this procedure when the deployment process installs the release with credentials that can manage every resource in
the release namespace and cannot create cluster-scoped resources. An administrator identity creates the cluster-scoped
resources first. The Monitoring Operator then runs with `global.privilegedRights: false` and does not create or manage
its cluster RBAC.

The collectors still read the cluster. kube-state-metrics and VMAgent use the administrator-created bindings, not the
Helm credential.

## Install prerequisites

The administrator identity creates the release namespace and installs the CRDs:

```bash
kubectl create namespace monitoring
kubectl apply --server-side --force-conflicts -f charts/qubership-monitoring-crds/crds/
```

The default profile needs administrator-owned RBAC for kube-state-metrics and VMAgent. Apply the manifests before the
deployment process runs Helm:

```bash
export MONITORING_NAMESPACE=monitoring

envsubst < docs/installation/restricted-rbac/kube-state-metrics.yaml | kubectl apply -f -
envsubst < docs/installation/restricted-rbac/vmagent.yaml | kubectl apply -f -
envsubst < docs/installation/restricted-rbac/restricted-helm.yaml | kubectl apply -f -
```

The manifests bind the cluster roles to their service accounts in `MONITORING_NAMESPACE`. Keep them aligned with the
chart version. When the deployment process enables components outside the default profile, review their cluster RBAC
before enabling them. The chart does not install the etcd certificate Job in restricted mode. To collect etcd metrics,
create the `kube-etcd-client-certs` Secret before installation as described in
[etcd metrics](../metrics-collection/metrics/etcd-metrics.md).

`restricted-helm.yaml` is an example. Bind its `Role` to the deployment identity when that identity is not the
`restricted-helm` ServiceAccount.

On OpenShift, `vmagent.yaml` grants VMAgent permission to use the `victoriametrics-operator` security context
constraint.

## Install a restricted release

Create a values file for the deployment process:

```yaml
global:
  privilegedRights: false
```

Run Helm with the namespace-only credential after the administrator has created the namespace, installed the CRDs, and
applied the manifests. That credential cannot create the namespace:

```bash
helm install monitoring-operator charts/qubership-monitoring-operator \
  --skip-crds \
  --namespace monitoring \
  --values restricted-values.yaml
```

Use the same chart version for the CRDs, administrator-provided RBAC, and Helm release.

## Verify the installation

Confirm that the collector service accounts have the required access:

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
