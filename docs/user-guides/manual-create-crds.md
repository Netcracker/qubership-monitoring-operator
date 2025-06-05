This document describes how to manually create, update or delete Monitoring Custom Resource Definitions (CRDs).

# When is it needed?

Almost all Qubership applications and microservices are integrated with Monitoring. This integration means
that almost all microservices during deploy can create `ServiceMonitor`/`PodMonitor`/`Probe` objects from the API
`monitoring.coreos.com`.

Such objects allow for Monitoring to understand from which microservices and how to collect metrics.
Similar objects allow providing alerts, recording rules and Grafana dashboards.

Also, it means that before deploying applications and microservices in Kubernetes all Custom Resource Definitions
(CRDs) must already created in Kubernetes. Otherwise, the deployment will fail.

The deployment process (helm chart) of Monitoring will automatically install or update all required CRDs.
But there are environments where Monitoring can't or doesn't make sense to deploy. In such cases, there is
an ability to create or update CRDs manually.

In the list of Monitoring artifacts, you can find an archive with the name

```bash
monitoring-operator-<version>-crds.zip
```

that contains all CRDs required for Monitoring. How to use this archive and CRDs inside you can read below.

# Before you begin

* You should have cluster-wide permissions enough to operate with CRDs (cluster admin is not required)
* You should configure a context for your `kubectl` and make sure the connection configured to correct Kubernetes

# How to manage CRDs

This section describes different cases of manual manipulation with CRDs.

## Create

To create CRDs for Monitoring you need to execute the command:

```bash
kubectl create -f path/to/crds/directory/
```

**Warning!** Never use the `kubectl apply` command! This command will generate the annotation
`kubectl.kubernetes.io/last-applied-configuration` that will contain the whole content of CRD.
It means that the object size will be increasing two times and it may lead to problems with storing
this CRD in Etcd.

```bash
mkdir /tmp/crds/
unzip -d /tmp/crds/ monitoring-operator-<version>-crds.zip
kubectl create -f /tmp/crds/*
```

## Upgrade

To upgrade CRDs for Monitoring you need to execute the command:

```bash
kubectl replace -f path/to/crds/directory/
```

**Warning!** Never use the `kubectl apply` command! This command will generate the annotation
`kubectl.kubernetes.io/last-applied-configuration` that will contain the whole content of CRD.
It means that the object size will be increasing two times and it may lead to problems with storing
this CRD in Etcd.

For example, if you will use the archive with CRDs providing Monitoring:

```bash
mkdir /tmp/crds/
unzip -d /tmp/crds/ monitoring-operator-<version>-crds.zip
kubectl replace -f /tmp/crds/*
```

## Remove

**Warning!** This step removes **all CRDs** for Monitoring and deleting these CRDs causes the deletion of
**all resources** of their type in **all namespaces**.
It means that all resources like `ServiceMonitor` and `GrafanaDashboard` from applications are removed.

To remove CRDs and all Custom Resources (CRs) like ServiceMonitor or GrafanaDashboards from all namespaces
for Monitoring you need to execute the command:

```bash
# monitoring-operator CRDs
kubectl delete crd platformmonitorings.monitoring.qubership.org

# prometheus-adapter-operator CRDs
kubectl delete crd customscalemetricrules.monitoring.qubership.org
kubectl delete crd prometheusadapters.monitoring.qubership.org

# grafana-operator CRDs
kubectl delete crd grafanas.integreatly.org
kubectl delete crd grafanadashboards.integreatly.org
kubectl delete crd grafanafolders.integreatly.org
kubectl delete crd grafanadatasources.integreatly.org
kubectl delete crd grafananotificationchannels.integreatly.org

# prometheus-operator CRDs
kubectl delete crd alertmanagers.monitoring.coreos.com
kubectl delete crd alertmanagerconfigs.monitoring.coreos.com
kubectl delete crd podmonitors.monitoring.coreos.com
kubectl delete crd probes.monitoring.coreos.com
kubectl delete crd prometheusagents.monitoring.coreos.com
kubectl delete crd prometheuses.monitoring.coreos.com
kubectl delete crd prometheusrules.monitoring.coreos.com
kubectl delete crd scrapeconfigs.monitoring.coreos.com
kubectl delete crd servicemonitors.monitoring.coreos.com
kubectl delete crd thanosrulers.monitoring.coreos.com

# victoriametrics-operator CRDs
kubectl delete crd vmagents.operator.victoriametrics.com
kubectl delete crd vmalertmanagerconfigs..operator.victoriametrics.com
kubectl delete crd vmalerts.operator.victoriametrics.com
kubectl delete crd vmalertmanagers.operator.victoriametrics.com
kubectl delete crd vmauths.operator.victoriametrics.com
kubectl delete crd vmclusters.operator.victoriametrics.com
kubectl delete crd vmnodescrapes.operator.victoriametrics.com
kubectl delete crd vmpodscrapes.operator.victoriametrics.com
kubectl delete crd vmprobes.operator.victoriametrics.com
kubectl delete crd vmrules.operator.victoriametrics.com
kubectl delete crd vmservicescrapes.operator.victoriametrics.com
kubectl delete crd vmsingles.operator.victoriametrics.com
kubectl delete crd vmstaticscrapes.operator.victoriametrics.com
kubectl delete crd vmusers.operator.victoriametrics.com
```