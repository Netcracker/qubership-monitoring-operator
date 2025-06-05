This document describes openshift cluster version operator metrics list and how to collect them.

# Metrics

Consul already can exposes its metrics in Prometheus format and doesn't require to use of specific exporters.

| Name       | Metrics Port | Metrics Endpoint    | Need Exporter? | Auth?          | Is Exporter Third Party? |
| ---------- | ------------ | ------------------- | -------------- | -------------- | ------------------------ |
| Prometheus | `9099`       | `/metrics`          | No             | Require, token | N/A                      |

## How to Collect

Metrics expose on port `9099` and endpoint `/metrics`. By default, Openshift cluster version operator has authentication
by token.

Config `ServiceMonitor` for `prometheus-operator` to collect Openshift cluster version operator:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: cluster-version-operator
  labels:
    k8s-app: cluster-version-operator
    app.kubernetes.io/name: cluster-version-operator
    app.kubernetes.io/component: monitoring
    app.kubernetes.io/managed-by: monitoring-operator
spec:
  endpoints:
    - bearerTokenFile: /var/run/secrets/kubernetes.io/serviceaccount/token
      interval: 30s
      port: metrics
      scheme: https
      tlsConfig:
        insecureSkipVerify: true
  namespaceSelector:
    matchNames:
      - openshift-cluster-version
  selector:
    matchLabels:
      k8s-app: cluster-version-operator
```

To collect (or just to check) metrics manually you can use the following command:

```bash
curl -v -k -L -H "Authorization: Bearer <token>" "http://<api_server_ip_or_dns>:88443/metrics"
```

Token usually you can find in the prometheus/vmagent pod by path `/var/run/secrets/kubernetes.io/serviceaccount/token`.

You can't use `wget` because it doesn't allow to add headers for authorization.

## Metrics List

### Openshift state metrics

```prometheus
# HELP cluster_installer Reports info about the installation process and, if applicable, the install tool. The type is either 'openshift-install', indicating that openshift-install was used to install the cluster, or 'other', indicating that an unknown process installed the cluster. The invoker is 'user' by default, but it may be overridden by a consuming tool. The version reported is that of the openshift-install that was used to generate the manifests and, if applicable, provision the infrastructure.
# TYPE cluster_installer gauge
cluster_installer{invoker="user",type="other",version="v4.12.0"} 1
# HELP cluster_operator_condition_transitions Reports the number of times that a condition on a cluster operator changes status
# TYPE cluster_operator_condition_transitions gauge
cluster_operator_condition_transitions{condition="Available",name="authentication"} 2372
cluster_operator_condition_transitions{condition="Available",name="console"} 212
cluster_operator_condition_transitions{condition="Available",name="kube-storage-version-migrator"} 4
cluster_operator_condition_transitions{condition="Available",name="machine-config"} 10
cluster_operator_condition_transitions{condition="Available",name="monitoring"} 38
cluster_operator_condition_transitions{condition="Available",name="openshift-apiserver"} 1940
cluster_operator_condition_transitions{condition="Degraded",name="authentication"} 14
cluster_operator_condition_transitions{condition="Degraded",name="cloud-controller-manager"} 18
cluster_operator_condition_transitions{condition="Degraded",name="console"} 12
cluster_operator_condition_transitions{condition="Degraded",name="etcd"} 18
cluster_operator_condition_transitions{condition="Degraded",name="ingress"} 4
cluster_operator_condition_transitions{condition="Degraded",name="kube-apiserver"} 7
cluster_operator_condition_transitions{condition="Degraded",name="kube-controller-manager"} 8
cluster_operator_condition_transitions{condition="Degraded",name="kube-scheduler"} 2
cluster_operator_condition_transitions{condition="Degraded",name="machine-config"} 11
cluster_operator_condition_transitions{condition="Degraded",name="monitoring"} 38
cluster_operator_condition_transitions{condition="Degraded",name="network"} 128
cluster_operator_condition_transitions{condition="Progressing",name="csi-snapshot-controller"} 2
cluster_operator_condition_transitions{condition="Progressing",name="dns"} 50
cluster_operator_condition_transitions{condition="Progressing",name="image-registry"} 48
cluster_operator_condition_transitions{condition="Progressing",name="ingress"} 8
cluster_operator_condition_transitions{condition="Progressing",name="kube-apiserver"} 29
cluster_operator_condition_transitions{condition="Progressing",name="kube-storage-version-migrator"} 4
cluster_operator_condition_transitions{condition="Progressing",name="machine-config"} 1
cluster_operator_condition_transitions{condition="Progressing",name="monitoring"} 38
cluster_operator_condition_transitions{condition="Progressing",name="network"} 391
cluster_operator_condition_transitions{condition="Progressing",name="node-tuning"} 8
cluster_operator_condition_transitions{condition="Progressing",name="openshift-controller-manager"} 42
cluster_operator_condition_transitions{condition="TrustedCABundleControllerControllerAvailable",name="cloud-controller-manager"} 32
cluster_operator_condition_transitions{condition="TrustedCABundleControllerControllerDegraded",name="cloud-controller-manager"} 32
cluster_operator_condition_transitions{condition="Upgradeable",name="cloud-controller-manager"} 18
cluster_operator_condition_transitions{condition="Upgradeable",name="console"} 464
cluster_operator_condition_transitions{condition="Upgradeable",name="machine-config"} 463
# HELP cluster_operator_conditions Report the conditions for active cluster operators. 0 is False and 1 is True.
# TYPE cluster_operator_conditions gauge
cluster_operator_conditions{condition="Available",name="authentication",reason="AsExpected"} 1
cluster_operator_conditions{condition="Available",name="baremetal",reason="WaitingForProvisioningCR"} 1
cluster_operator_conditions{condition="Available",name="cloud-controller-manager",reason="AsExpected"} 1
cluster_operator_conditions{condition="Available",name="cloud-credential",reason="AsExpected"} 1
cluster_operator_conditions{condition="Available",name="cluster-autoscaler",reason="AsExpected"} 1
cluster_operator_conditions{condition="Available",name="config-operator",reason="AsExpected"} 1
cluster_operator_conditions{condition="Available",name="console",reason="AsExpected"} 1
cluster_operator_conditions{condition="Available",name="control-plane-machine-set",reason="AsExpected"} 1
cluster_operator_conditions{condition="Available",name="csi-snapshot-controller",reason="AsExpected"} 1
cluster_operator_conditions{condition="Available",name="dns",reason="AsExpected"} 1
cluster_operator_conditions{condition="Available",name="etcd",reason="AsExpected"} 1
cluster_operator_conditions{condition="Available",name="image-registry",reason="Removed"} 1
cluster_operator_conditions{condition="Available",name="ingress",reason="IngressAvailable"} 1
cluster_operator_conditions{condition="Available",name="insights",reason="AsExpected"} 1
cluster_operator_conditions{condition="Available",name="kube-apiserver",reason="AsExpected"} 1
cluster_operator_conditions{condition="Available",name="kube-controller-manager",reason="AsExpected"} 1
cluster_operator_conditions{condition="Available",name="kube-scheduler",reason="AsExpected"} 1
cluster_operator_conditions{condition="Available",name="kube-storage-version-migrator",reason="AsExpected"} 1
cluster_operator_conditions{condition="Available",name="machine-api",reason="AsExpected"} 1
cluster_operator_conditions{condition="Available",name="machine-approver",reason="AsExpected"} 1
cluster_operator_conditions{condition="Available",name="machine-config",reason="AsExpected"} 1
cluster_operator_conditions{condition="Available",name="marketplace",reason="OperatorAvailable"} 1
cluster_operator_conditions{condition="Available",name="monitoring",reason="RollOutDone"} 1
cluster_operator_conditions{condition="Available",name="network",reason=""} 1
cluster_operator_conditions{condition="Available",name="node-tuning",reason="AsExpected"} 1
cluster_operator_conditions{condition="Available",name="openshift-apiserver",reason="AsExpected"} 1
cluster_operator_conditions{condition="Available",name="openshift-controller-manager",reason="AsExpected"} 1
cluster_operator_conditions{condition="Available",name="openshift-samples",reason="CurrentlyRemoved"} 1
cluster_operator_conditions{condition="Available",name="operator-lifecycle-manager",reason=""} 1
cluster_operator_conditions{condition="Available",name="operator-lifecycle-manager-catalog",reason=""} 1
cluster_operator_conditions{condition="Available",name="operator-lifecycle-manager-packageserver",reason="ClusterServiceVersionSucceeded"} 1
cluster_operator_conditions{condition="Available",name="service-ca",reason="AsExpected"} 1
cluster_operator_conditions{condition="Available",name="storage",reason="AsExpected"} 1
cluster_operator_conditions{condition="Available",name="version",reason=""} 1
cluster_operator_conditions{condition="CloudConfigControllerAvailable",name="cloud-controller-manager",reason="AsExpected"} 1
cluster_operator_conditions{condition="CloudConfigControllerDegraded",name="cloud-controller-manager",reason="AsExpected"} 0
cluster_operator_conditions{condition="CloudConfigControllerUpgradeable",name="cloud-controller-manager",reason="AsExpected"} 1
cluster_operator_conditions{condition="ClusterTransferAvailable",name="insights",reason="Disconnected"} 0
cluster_operator_conditions{condition="Degraded",name="authentication",reason="AsExpected"} 0
cluster_operator_conditions{condition="Degraded",name="baremetal",reason=""} 0
cluster_operator_conditions{condition="Degraded",name="cloud-controller-manager",reason="AsExpected"} 0
cluster_operator_conditions{condition="Degraded",name="cloud-credential",reason="AsExpected"} 0
cluster_operator_conditions{condition="Degraded",name="cluster-autoscaler",reason=""} 0
cluster_operator_conditions{condition="Degraded",name="config-operator",reason="AsExpected"} 0
cluster_operator_conditions{condition="Degraded",name="console",reason="AsExpected"} 0
cluster_operator_conditions{condition="Degraded",name="control-plane-machine-set",reason="AsExpected"} 0
cluster_operator_conditions{condition="Degraded",name="csi-snapshot-controller",reason="AsExpected"} 0
cluster_operator_conditions{condition="Degraded",name="dns",reason="DNSNotDegraded"} 0
cluster_operator_conditions{condition="Degraded",name="etcd",reason="AsExpected"} 0
cluster_operator_conditions{condition="Degraded",name="image-registry",reason="Removed"} 0
cluster_operator_conditions{condition="Degraded",name="ingress",reason="IngressNotDegraded"} 0
cluster_operator_conditions{condition="Degraded",name="insights",reason="AsExpected"} 0
cluster_operator_conditions{condition="Degraded",name="kube-apiserver",reason="MissingStaticPodController_SyncError::StaticPods_Error"} 1
cluster_operator_conditions{condition="Degraded",name="kube-controller-manager",reason="AsExpected"} 0
cluster_operator_conditions{condition="Degraded",name="kube-scheduler",reason="AsExpected"} 0
cluster_operator_conditions{condition="Degraded",name="kube-storage-version-migrator",reason="AsExpected"} 0
cluster_operator_conditions{condition="Degraded",name="machine-api",reason="AsExpected"} 0
cluster_operator_conditions{condition="Degraded",name="machine-approver",reason="AsExpected"} 0
cluster_operator_conditions{condition="Degraded",name="machine-config",reason=""} 0
cluster_operator_conditions{condition="Degraded",name="marketplace",reason="OperatorAvailable"} 0
cluster_operator_conditions{condition="Degraded",name="monitoring",reason="PrometheusDataPersistenceNotConfigured"} 0
cluster_operator_conditions{condition="Degraded",name="network",reason=""} 0
cluster_operator_conditions{condition="Degraded",name="node-tuning",reason="AsExpected"} 0
cluster_operator_conditions{condition="Degraded",name="openshift-apiserver",reason="AsExpected"} 0
cluster_operator_conditions{condition="Degraded",name="openshift-controller-manager",reason="AsExpected"} 0
cluster_operator_conditions{condition="Degraded",name="openshift-samples",reason=""} 0
cluster_operator_conditions{condition="Degraded",name="operator-lifecycle-manager",reason=""} 0
cluster_operator_conditions{condition="Degraded",name="operator-lifecycle-manager-catalog",reason=""} 0
cluster_operator_conditions{condition="Degraded",name="operator-lifecycle-manager-packageserver",reason=""} 0
cluster_operator_conditions{condition="Degraded",name="service-ca",reason="AsExpected"} 0
cluster_operator_conditions{condition="Degraded",name="storage",reason="AsExpected"} 0
cluster_operator_conditions{condition="Disabled",name="baremetal",reason=""} 0
cluster_operator_conditions{condition="Disabled",name="insights",reason="NoToken"} 1
cluster_operator_conditions{condition="EvaluationConditionsDetected",name="ingress",reason="AsExpected"} 0
cluster_operator_conditions{condition="Failing",name="version",reason="ClusterOperatorDegraded"} 1
cluster_operator_conditions{condition="ImplicitlyEnabledCapabilities",name="version",reason="AsExpected"} 0
cluster_operator_conditions{condition="ManagementStateDegraded",name="network",reason=""} 0
cluster_operator_conditions{condition="Progressing",name="authentication",reason="AsExpected"} 0
cluster_operator_conditions{condition="Progressing",name="baremetal",reason="WaitingForProvisioningCR"} 0
cluster_operator_conditions{condition="Progressing",name="cloud-controller-manager",reason="AsExpected"} 0
cluster_operator_conditions{condition="Progressing",name="cloud-credential",reason="AsExpected"} 0
cluster_operator_conditions{condition="Progressing",name="cluster-autoscaler",reason=""} 0
cluster_operator_conditions{condition="Progressing",name="config-operator",reason="AsExpected"} 0
cluster_operator_conditions{condition="Progressing",name="console",reason="AsExpected"} 0
cluster_operator_conditions{condition="Progressing",name="control-plane-machine-set",reason="AsExpected"} 0
cluster_operator_conditions{condition="Progressing",name="csi-snapshot-controller",reason="AsExpected"} 0
cluster_operator_conditions{condition="Progressing",name="dns",reason="AsExpected"} 0
cluster_operator_conditions{condition="Progressing",name="etcd",reason="AsExpected"} 0
cluster_operator_conditions{condition="Progressing",name="image-registry",reason="Removed"} 0
cluster_operator_conditions{condition="Progressing",name="ingress",reason="AsExpected"} 0
cluster_operator_conditions{condition="Progressing",name="insights",reason="AsExpected"} 0
cluster_operator_conditions{condition="Progressing",name="kube-apiserver",reason="NodeInstaller"} 1
cluster_operator_conditions{condition="Progressing",name="kube-controller-manager",reason="AsExpected"} 0
cluster_operator_conditions{condition="Progressing",name="kube-scheduler",reason="AsExpected"} 0
cluster_operator_conditions{condition="Progressing",name="kube-storage-version-migrator",reason="AsExpected"} 0
cluster_operator_conditions{condition="Progressing",name="machine-api",reason="AsExpected"} 0
cluster_operator_conditions{condition="Progressing",name="machine-approver",reason="AsExpected"} 0
cluster_operator_conditions{condition="Progressing",name="machine-config",reason=""} 0
cluster_operator_conditions{condition="Progressing",name="marketplace",reason="OperatorAvailable"} 0
cluster_operator_conditions{condition="Progressing",name="monitoring",reason=""} 0
cluster_operator_conditions{condition="Progressing",name="network",reason=""} 0
cluster_operator_conditions{condition="Progressing",name="node-tuning",reason="AsExpected"} 0
cluster_operator_conditions{condition="Progressing",name="openshift-apiserver",reason="AsExpected"} 0
cluster_operator_conditions{condition="Progressing",name="openshift-controller-manager",reason="AsExpected"} 0
cluster_operator_conditions{condition="Progressing",name="openshift-samples",reason="CurrentlyRemoved"} 0
cluster_operator_conditions{condition="Progressing",name="operator-lifecycle-manager",reason=""} 0
cluster_operator_conditions{condition="Progressing",name="operator-lifecycle-manager-catalog",reason=""} 0
cluster_operator_conditions{condition="Progressing",name="operator-lifecycle-manager-packageserver",reason=""} 0
cluster_operator_conditions{condition="Progressing",name="service-ca",reason="AsExpected"} 0
cluster_operator_conditions{condition="Progressing",name="storage",reason="AsExpected"} 0
cluster_operator_conditions{condition="Progressing",name="version",reason="ClusterOperatorDegraded"} 0
cluster_operator_conditions{condition="RecentBackup",name="etcd",reason="UpgradeBackupSuccessful"} 1
cluster_operator_conditions{condition="ReleaseAccepted",name="version",reason="PayloadLoaded"} 1
cluster_operator_conditions{condition="RetrievedUpdates",name="version",reason="RemoteFailed"} 0
cluster_operator_conditions{condition="SCAAvailable",name="insights",reason="NonHTTPError"} 0
cluster_operator_conditions{condition="TrustedCABundleControllerControllerAvailable",name="cloud-controller-manager",reason="AsExpected"} 1
cluster_operator_conditions{condition="TrustedCABundleControllerControllerDegraded",name="cloud-controller-manager",reason="AsExpected"} 0
cluster_operator_conditions{condition="Upgradeable",name="authentication",reason="AsExpected"} 1
cluster_operator_conditions{condition="Upgradeable",name="baremetal",reason=""} 1
cluster_operator_conditions{condition="Upgradeable",name="cloud-controller-manager",reason="AsExpected"} 1
cluster_operator_conditions{condition="Upgradeable",name="cloud-credential",reason=""} 1
cluster_operator_conditions{condition="Upgradeable",name="cluster-autoscaler",reason=""} 1
cluster_operator_conditions{condition="Upgradeable",name="config-operator",reason="AsExpected"} 1
cluster_operator_conditions{condition="Upgradeable",name="console",reason="AsExpected"} 1
cluster_operator_conditions{condition="Upgradeable",name="control-plane-machine-set",reason="AsExpected"} 1
cluster_operator_conditions{condition="Upgradeable",name="csi-snapshot-controller",reason="AsExpected"} 1
cluster_operator_conditions{condition="Upgradeable",name="dns",reason="DNSUpgradeable"} 1
cluster_operator_conditions{condition="Upgradeable",name="etcd",reason="AsExpected"} 1
cluster_operator_conditions{condition="Upgradeable",name="ingress",reason="IngressControllersUpgradeable"} 1
cluster_operator_conditions{condition="Upgradeable",name="insights",reason="InsightsUpgradeable"} 1
cluster_operator_conditions{condition="Upgradeable",name="kube-apiserver",reason="AsExpected"} 1
cluster_operator_conditions{condition="Upgradeable",name="kube-controller-manager",reason="AsExpected"} 1
cluster_operator_conditions{condition="Upgradeable",name="kube-scheduler",reason="AsExpected"} 1
cluster_operator_conditions{condition="Upgradeable",name="kube-storage-version-migrator",reason="AsExpected"} 1
cluster_operator_conditions{condition="Upgradeable",name="machine-api",reason=""} 1
cluster_operator_conditions{condition="Upgradeable",name="machine-approver",reason="AsExpected"} 1
cluster_operator_conditions{condition="Upgradeable",name="machine-config",reason="AsExpected"} 1
cluster_operator_conditions{condition="Upgradeable",name="marketplace",reason="OperatorAvailable"} 1
cluster_operator_conditions{condition="Upgradeable",name="monitoring",reason=""} 1
cluster_operator_conditions{condition="Upgradeable",name="network",reason=""} 1
cluster_operator_conditions{condition="Upgradeable",name="openshift-apiserver",reason="AsExpected"} 1
cluster_operator_conditions{condition="Upgradeable",name="openshift-controller-manager",reason="AsExpected"} 1
cluster_operator_conditions{condition="Upgradeable",name="operator-lifecycle-manager",reason=""} 1
cluster_operator_conditions{condition="Upgradeable",name="operator-lifecycle-manager-catalog",reason=""} 1
cluster_operator_conditions{condition="Upgradeable",name="operator-lifecycle-manager-packageserver",reason=""} 1
cluster_operator_conditions{condition="Upgradeable",name="service-ca",reason="AsExpected"} 1
cluster_operator_conditions{condition="Upgradeable",name="storage",reason="AsExpected"} 1
cluster_operator_conditions{condition="Upgradeable",name="version",reason="AdminAckRequired"} 0
# HELP cluster_operator_payload_errors Report the number of errors encountered applying the payload.
# TYPE cluster_operator_payload_errors gauge
cluster_operator_payload_errors{version="4.13.9"} 4122
# HELP cluster_operator_up 1 if a cluster operator is Available=True.  0 otherwise, including if a cluster operator sets no Available condition.  The 'version' label tracks the 'operator' version.  The 'reason' label is passed through from the Available condition, unless the cluster operator sets no Available condition, in which case NoAvailableCondition is used.
# TYPE cluster_operator_up gauge
cluster_operator_up{name="authentication",reason="AsExpected",version="4.13.9"} 1
cluster_operator_up{name="baremetal",reason="WaitingForProvisioningCR",version="4.13.9"} 1
cluster_operator_up{name="cloud-controller-manager",reason="AsExpected",version="4.13.9"} 1
cluster_operator_up{name="cloud-credential",reason="AsExpected",version="4.13.9"} 1
cluster_operator_up{name="cluster-autoscaler",reason="AsExpected",version="4.13.9"} 1
cluster_operator_up{name="config-operator",reason="AsExpected",version="4.13.9"} 1
cluster_operator_up{name="console",reason="AsExpected",version="4.13.9"} 1
cluster_operator_up{name="control-plane-machine-set",reason="AsExpected",version="4.13.9"} 1
cluster_operator_up{name="csi-snapshot-controller",reason="AsExpected",version="4.13.9"} 1
cluster_operator_up{name="dns",reason="AsExpected",version="4.13.9"} 1
cluster_operator_up{name="etcd",reason="AsExpected",version="4.13.9"} 1
cluster_operator_up{name="image-registry",reason="Removed",version="4.13.9"} 1
cluster_operator_up{name="ingress",reason="IngressAvailable",version="4.13.9"} 1
cluster_operator_up{name="insights",reason="AsExpected",version="4.13.9"} 1
cluster_operator_up{name="kube-apiserver",reason="AsExpected",version="4.13.9"} 1
cluster_operator_up{name="kube-controller-manager",reason="AsExpected",version="4.13.9"} 1
cluster_operator_up{name="kube-scheduler",reason="AsExpected",version="4.13.9"} 1
cluster_operator_up{name="kube-storage-version-migrator",reason="AsExpected",version="4.13.9"} 1
cluster_operator_up{name="machine-api",reason="AsExpected",version="4.13.9"} 1
cluster_operator_up{name="machine-approver",reason="AsExpected",version="4.13.9"} 1
cluster_operator_up{name="machine-config",reason="AsExpected",version="4.13.9"} 1
cluster_operator_up{name="marketplace",reason="OperatorAvailable",version="4.13.9"} 1
cluster_operator_up{name="monitoring",reason="RollOutDone",version="4.13.9"} 1
cluster_operator_up{name="network",reason="",version="4.13.9"} 1
cluster_operator_up{name="node-tuning",reason="AsExpected",version="4.13.9"} 1
cluster_operator_up{name="openshift-apiserver",reason="AsExpected",version="4.13.9"} 1
cluster_operator_up{name="openshift-controller-manager",reason="AsExpected",version="4.13.9"} 1
cluster_operator_up{name="openshift-samples",reason="CurrentlyRemoved",version="4.13.9"} 1
cluster_operator_up{name="operator-lifecycle-manager",reason="",version="4.13.9"} 1
cluster_operator_up{name="operator-lifecycle-manager-catalog",reason="",version="4.13.9"} 1
cluster_operator_up{name="operator-lifecycle-manager-packageserver",reason="ClusterServiceVersionSucceeded",version="4.13.9"} 1
cluster_operator_up{name="service-ca",reason="AsExpected",version="4.13.9"} 1
cluster_operator_up{name="storage",reason="AsExpected",version="4.13.9"} 1
# HELP cluster_version Reports the version of the cluster in terms of seconds since\nthe epoch. Type 'current' is the version being applied and\nthe value is the creation date of the payload. The type\n'desired' is returned if spec.desiredUpdate is set but the\noperator has not yet updated and the value is the most \nrecent status transition time. The type 'failure' is set \nif an error is preventing sync or upgrade with the last \ntransition timestamp of the condition. The type 'completed' \nis the timestamp when the last image was successfully\napplied. The type 'cluster' is the creation date of the\ncluster version object and the current version. The type\n'updating' is set when the cluster is transitioning to a\nnew version but has not reached the completed state and\nis the time the update was started. The type 'initial' is\nset to the oldest entry in the history. The from_version label\nwill be set to the last completed version, the initial\nversion for 'cluster', or empty for 'initial'.\n.
# TYPE cluster_version gauge
cluster_version{from_version="",image="quay.io/openshift-release-dev/ocp-release@sha256:ba7956f5c2aae61c8ff3ab1ab2ee7e625db9b1c8964a65339764db79c148e4e6",type="initial",version="4.12.22"} 1.689347834e+09
cluster_version{from_version="4.12.22",image="openshift-release-dev/ocp-release@sha256:a266d3d65c433b460cdef7ab5d6531580f5391adbe85d9c475208a56452e4c0b",type="cluster",version="4.13.9"} 1.689347834e+09
cluster_version{from_version="4.12.22",image="openshift-release-dev/ocp-release@sha256:a266d3d65c433b460cdef7ab5d6531580f5391adbe85d9c475208a56452e4c0b",type="completed",version="4.13.9"} 1.692368146e+09
cluster_version{from_version="4.13.9",image="openshift-release-dev/ocp-release@sha256:a266d3d65c433b460cdef7ab5d6531580f5391adbe85d9c475208a56452e4c0b",type="current",version="4.13.9"} 1.691660236e+09
cluster_version{from_version="4.13.9",image="openshift-release-dev/ocp-release@sha256:a266d3d65c433b460cdef7ab5d6531580f5391adbe85d9c475208a56452e4c0b",type="failure",version="4.13.9"} 1.700449562e+09
# HELP cluster_version_capability Report currently enabled cluster capabilities.  0 is disabled, and 1 is enabled.
# TYPE cluster_version_capability gauge
cluster_version_capability{name="CSISnapshot"} 1
cluster_version_capability{name="Console"} 1
cluster_version_capability{name="Insights"} 1
cluster_version_capability{name="NodeTuning"} 1
cluster_version_capability{name="Storage"} 1
cluster_version_capability{name="baremetal"} 1
cluster_version_capability{name="marketplace"} 1
cluster_version_capability{name="openshift-samples"} 1
# HELP cluster_version_operator_update_retrieval_timestamp_seconds Reports when updates were last successfully retrieved.
# TYPE cluster_version_operator_update_retrieval_timestamp_seconds gauge
cluster_version_operator_update_retrieval_timestamp_seconds{name=""} 1.692367394e+09
# HELP cluster_version_payload Report the number of entries in the payload.
# TYPE cluster_version_payload gauge
cluster_version_payload{type="applied",version="4.13.9"} 837
cluster_version_payload{type="pending",version="4.13.9"} 5
# HELP go_cgo_go_to_c_calls_calls_total Count of calls made from Go to C by the current process.
# TYPE go_cgo_go_to_c_calls_calls_total counter
go_cgo_go_to_c_calls_calls_total 6
# HELP go_gc_cycles_automatic_gc_cycles_total Count of completed GC cycles generated by the Go runtime.
# TYPE go_gc_cycles_automatic_gc_cycles_total counter
go_gc_cycles_automatic_gc_cycles_total 488426
# HELP go_gc_cycles_forced_gc_cycles_total Count of completed GC cycles forced by the application.
# TYPE go_gc_cycles_forced_gc_cycles_total counter
go_gc_cycles_forced_gc_cycles_total 0
# HELP go_gc_cycles_total_gc_cycles_total Count of all completed GC cycles.
# TYPE go_gc_cycles_total_gc_cycles_total counter
go_gc_cycles_total_gc_cycles_total 488426
# HELP go_gc_duration_seconds A summary of the pause duration of garbage collection cycles.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 5.4778e-05
go_gc_duration_seconds{quantile="0.25"} 7.9409e-05
go_gc_duration_seconds{quantile="0.5"} 0.000106051
go_gc_duration_seconds{quantile="0.75"} 0.00016854
go_gc_duration_seconds{quantile="1"} 0.046073051
go_gc_duration_seconds_sum 198.518708721
go_gc_duration_seconds_count 488426
# HELP go_gc_heap_allocs_by_size_bytes_total Distribution of heap allocations by approximate size. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks.
# TYPE go_gc_heap_allocs_by_size_bytes_total histogram
go_gc_heap_allocs_by_size_bytes_total_bucket{le="8.999999999999998"} 4.82474013e+08
go_gc_heap_allocs_by_size_bytes_total_bucket{le="24.999999999999996"} 1.8080876464e+10
go_gc_heap_allocs_by_size_bytes_total_bucket{le="64.99999999999999"} 2.9100532123e+10
go_gc_heap_allocs_by_size_bytes_total_bucket{le="144.99999999999997"} 3.3554110642e+10
go_gc_heap_allocs_by_size_bytes_total_bucket{le="320.99999999999994"} 3.5983546659e+10
go_gc_heap_allocs_by_size_bytes_total_bucket{le="704.9999999999999"} 4.0325420111e+10
go_gc_heap_allocs_by_size_bytes_total_bucket{le="1536.9999999999998"} 4.0616253514e+10
go_gc_heap_allocs_by_size_bytes_total_bucket{le="3200.9999999999995"} 4.069938606e+10
go_gc_heap_allocs_by_size_bytes_total_bucket{le="6528.999999999999"} 4.0735569022e+10
go_gc_heap_allocs_by_size_bytes_total_bucket{le="13568.999999999998"} 4.1265698649e+10
go_gc_heap_allocs_by_size_bytes_total_bucket{le="27264.999999999996"} 4.127497466e+10
go_gc_heap_allocs_by_size_bytes_total_bucket{le="+Inf"} 4.1288698108e+10
go_gc_heap_allocs_by_size_bytes_total_sum 1.0637402843264e+13
go_gc_heap_allocs_by_size_bytes_total_count 4.1288698108e+10
# HELP go_gc_heap_allocs_bytes_total Cumulative sum of memory allocated to the heap by the application.
# TYPE go_gc_heap_allocs_bytes_total counter
go_gc_heap_allocs_bytes_total 1.0637402843264e+13
# HELP go_gc_heap_allocs_objects_total Cumulative count of heap allocations triggered by the application. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks.
# TYPE go_gc_heap_allocs_objects_total counter
go_gc_heap_allocs_objects_total 4.1288698108e+10
# HELP go_gc_heap_frees_by_size_bytes_total Distribution of freed heap allocations by approximate size. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks.
# TYPE go_gc_heap_frees_by_size_bytes_total histogram
go_gc_heap_frees_by_size_bytes_total_bucket{le="8.999999999999998"} 4.82468573e+08
go_gc_heap_frees_by_size_bytes_total_bucket{le="24.999999999999996"} 1.8080715199e+10
go_gc_heap_frees_by_size_bytes_total_bucket{le="64.99999999999999"} 2.9100304738e+10
go_gc_heap_frees_by_size_bytes_total_bucket{le="144.99999999999997"} 3.355386747e+10
go_gc_heap_frees_by_size_bytes_total_bucket{le="320.99999999999994"} 3.59832735e+10
go_gc_heap_frees_by_size_bytes_total_bucket{le="704.9999999999999"} 4.032514151e+10
go_gc_heap_frees_by_size_bytes_total_bucket{le="1536.9999999999998"} 4.0615973778e+10
go_gc_heap_frees_by_size_bytes_total_bucket{le="3200.9999999999995"} 4.0699105943e+10
go_gc_heap_frees_by_size_bytes_total_bucket{le="6528.999999999999"} 4.0735288668e+10
go_gc_heap_frees_by_size_bytes_total_bucket{le="13568.999999999998"} 4.1265417617e+10
go_gc_heap_frees_by_size_bytes_total_bucket{le="27264.999999999996"} 4.1274693567e+10
go_gc_heap_frees_by_size_bytes_total_bucket{le="+Inf"} 4.1288416912e+10
go_gc_heap_frees_by_size_bytes_total_sum 1.0637363382712e+13
go_gc_heap_frees_by_size_bytes_total_count 4.1288416912e+10
# HELP go_gc_heap_frees_bytes_total Cumulative sum of heap memory freed by the garbage collector.
# TYPE go_gc_heap_frees_bytes_total counter
go_gc_heap_frees_bytes_total 1.0637363382712e+13
# HELP go_gc_heap_frees_objects_total Cumulative count of heap allocations whose storage was freed by the garbage collector. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks.
# TYPE go_gc_heap_frees_objects_total counter
go_gc_heap_frees_objects_total 4.1288416912e+10
# HELP go_gc_heap_goal_bytes Heap size target for the end of the GC cycle.
# TYPE go_gc_heap_goal_bytes gauge
go_gc_heap_goal_bytes 5.1911216e+07
# HELP go_gc_heap_objects_objects Number of objects, live or unswept, occupying heap memory.
# TYPE go_gc_heap_objects_objects gauge
go_gc_heap_objects_objects 281196
# HELP go_gc_heap_tiny_allocs_objects_total Count of small allocations that are packed together into blocks. These allocations are counted separately from other allocations because each individual allocation is not tracked by the runtime, only their block. Each block is already accounted for in allocs-by-size and frees-by-size.
# TYPE go_gc_heap_tiny_allocs_objects_total counter
go_gc_heap_tiny_allocs_objects_total 5.057530995e+09
# HELP go_gc_limiter_last_enabled_gc_cycle GC cycle the last time the GC CPU limiter was enabled. This metric is useful for diagnosing the root cause of an out-of-memory error, because the limiter trades memory for CPU time when the GC's CPU time gets too high. This is most likely to occur with use of SetMemoryLimit. The first GC cycle is cycle 1, so a value of 0 indicates that it was never enabled.
# TYPE go_gc_limiter_last_enabled_gc_cycle gauge
go_gc_limiter_last_enabled_gc_cycle 0
# HELP go_gc_pauses_seconds_total Distribution individual GC-related stop-the-world pause latencies.
# TYPE go_gc_pauses_seconds_total histogram
go_gc_pauses_seconds_total_bucket{le="-5e-324"} 0
go_gc_pauses_seconds_total_bucket{le="9.999999999999999e-10"} 0
go_gc_pauses_seconds_total_bucket{le="9.999999999999999e-09"} 0
go_gc_pauses_seconds_total_bucket{le="9.999999999999998e-08"} 0
go_gc_pauses_seconds_total_bucket{le="1.0239999999999999e-06"} 0
go_gc_pauses_seconds_total_bucket{le="1.0239999999999999e-05"} 379328
go_gc_pauses_seconds_total_bucket{le="0.00010239999999999998"} 732155
go_gc_pauses_seconds_total_bucket{le="0.0010485759999999998"} 949366
go_gc_pauses_seconds_total_bucket{le="0.010485759999999998"} 974775
go_gc_pauses_seconds_total_bucket{le="0.10485759999999998"} 976823
go_gc_pauses_seconds_total_bucket{le="+Inf"} 976852
go_gc_pauses_seconds_total_sum NaN
go_gc_pauses_seconds_total_count 976852
# HELP go_gc_stack_starting_size_bytes The stack size of new goroutines.
# TYPE go_gc_stack_starting_size_bytes gauge
go_gc_stack_starting_size_bytes 4096
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 89
# HELP go_info Information about the Go environment.
# TYPE go_info gauge
go_info{version="go1.19.10 X:strictfipsruntime"} 1
# HELP go_memory_classes_heap_free_bytes Memory that is completely free and eligible to be returned to the underlying system, but has not been. This metric is the runtime's estimate of free address space that is backed by physical memory.
# TYPE go_memory_classes_heap_free_bytes gauge
go_memory_classes_heap_free_bytes 8.609792e+06
# HELP go_memory_classes_heap_objects_bytes Memory occupied by live objects and dead objects that have not yet been marked free by the garbage collector.
# TYPE go_memory_classes_heap_objects_bytes gauge
go_memory_classes_heap_objects_bytes 3.9460552e+07
# HELP go_memory_classes_heap_released_bytes Memory that is completely free and has been returned to the underlying system. This metric is the runtime's estimate of free address space that is still mapped into the process, but is not backed by physical memory.
# TYPE go_memory_classes_heap_released_bytes gauge
go_memory_classes_heap_released_bytes 1.908736e+07
# HELP go_memory_classes_heap_stacks_bytes Memory allocated from the heap that is reserved for stack space, whether or not it is currently in-use.
# TYPE go_memory_classes_heap_stacks_bytes gauge
go_memory_classes_heap_stacks_bytes 2.031616e+06
# HELP go_memory_classes_heap_unused_bytes Memory that is reserved for heap objects but is not currently used to hold heap objects.
# TYPE go_memory_classes_heap_unused_bytes gauge
go_memory_classes_heap_unused_bytes 1.0502456e+07
# HELP go_memory_classes_metadata_mcache_free_bytes Memory that is reserved for runtime mcache structures, but not in-use.
# TYPE go_memory_classes_metadata_mcache_free_bytes gauge
go_memory_classes_metadata_mcache_free_bytes 10800
# HELP go_memory_classes_metadata_mcache_inuse_bytes Memory that is occupied by runtime mcache structures that are currently being used.
# TYPE go_memory_classes_metadata_mcache_inuse_bytes gauge
go_memory_classes_metadata_mcache_inuse_bytes 4800
# HELP go_memory_classes_metadata_mspan_free_bytes Memory that is reserved for runtime mspan structures, but not in-use.
# TYPE go_memory_classes_metadata_mspan_free_bytes gauge
go_memory_classes_metadata_mspan_free_bytes 440640
# HELP go_memory_classes_metadata_mspan_inuse_bytes Memory that is occupied by runtime mspan structures that are currently being used.
# TYPE go_memory_classes_metadata_mspan_inuse_bytes gauge
go_memory_classes_metadata_mspan_inuse_bytes 649584
# HELP go_memory_classes_metadata_other_bytes Memory that is reserved for or used to hold runtime metadata.
# TYPE go_memory_classes_metadata_other_bytes gauge
go_memory_classes_metadata_other_bytes 1.3159296e+07
# HELP go_memory_classes_os_stacks_bytes Stack memory allocated by the underlying operating system.
# TYPE go_memory_classes_os_stacks_bytes gauge
go_memory_classes_os_stacks_bytes 0
# HELP go_memory_classes_other_bytes Memory used by execution trace buffers, structures for debugging the runtime, finalizer and profiler specials, and more.
# TYPE go_memory_classes_other_bytes gauge
go_memory_classes_other_bytes 1.050192e+06
# HELP go_memory_classes_profiling_buckets_bytes Memory that is used by the stack trace hash map used for profiling.
# TYPE go_memory_classes_profiling_buckets_bytes gauge
go_memory_classes_profiling_buckets_bytes 5808
# HELP go_memory_classes_total_bytes All memory mapped by the Go runtime into the current process as read-write. Note that this does not include memory mapped by code called via cgo or via the syscall package. Sum of all metrics in /memory/classes.
# TYPE go_memory_classes_total_bytes gauge
go_memory_classes_total_bytes 9.5012896e+07
# HELP go_memstats_alloc_bytes Number of bytes allocated and still in use.
# TYPE go_memstats_alloc_bytes gauge
go_memstats_alloc_bytes 3.9460552e+07
# HELP go_memstats_alloc_bytes_total Total number of bytes allocated, even if freed.
# TYPE go_memstats_alloc_bytes_total counter
go_memstats_alloc_bytes_total 1.0637402843264e+13
# HELP go_memstats_buck_hash_sys_bytes Number of bytes used by the profiling bucket hash table.
# TYPE go_memstats_buck_hash_sys_bytes gauge
go_memstats_buck_hash_sys_bytes 5808
# HELP go_memstats_frees_total Total number of frees.
# TYPE go_memstats_frees_total counter
go_memstats_frees_total 4.6345947907e+10
# HELP go_memstats_gc_cpu_fraction The fraction of this program's available CPU time used by the GC since the program started.
# TYPE go_memstats_gc_cpu_fraction gauge
go_memstats_gc_cpu_fraction 0
# HELP go_memstats_gc_sys_bytes Number of bytes used for garbage collection system metadata.
# TYPE go_memstats_gc_sys_bytes gauge
go_memstats_gc_sys_bytes 1.3159296e+07
# HELP go_memstats_heap_alloc_bytes Number of heap bytes allocated and still in use.
# TYPE go_memstats_heap_alloc_bytes gauge
go_memstats_heap_alloc_bytes 3.9460552e+07
# HELP go_memstats_heap_idle_bytes Number of heap bytes waiting to be used.
# TYPE go_memstats_heap_idle_bytes gauge
go_memstats_heap_idle_bytes 2.7697152e+07
# HELP go_memstats_heap_inuse_bytes Number of heap bytes that are in use.
# TYPE go_memstats_heap_inuse_bytes gauge
go_memstats_heap_inuse_bytes 4.9963008e+07
# HELP go_memstats_heap_objects Number of allocated objects.
# TYPE go_memstats_heap_objects gauge
go_memstats_heap_objects 281196
# HELP go_memstats_heap_released_bytes Number of heap bytes released to OS.
# TYPE go_memstats_heap_released_bytes gauge
go_memstats_heap_released_bytes 1.908736e+07
# HELP go_memstats_heap_sys_bytes Number of heap bytes obtained from system.
# TYPE go_memstats_heap_sys_bytes gauge
go_memstats_heap_sys_bytes 7.766016e+07
# HELP go_memstats_last_gc_time_seconds Number of seconds since 1970 of last garbage collection.
# TYPE go_memstats_last_gc_time_seconds gauge
go_memstats_last_gc_time_seconds 1.7004618572108061e+09
# HELP go_memstats_lookups_total Total number of pointer lookups.
# TYPE go_memstats_lookups_total counter
go_memstats_lookups_total 0
# HELP go_memstats_mallocs_total Total number of mallocs.
# TYPE go_memstats_mallocs_total counter
go_memstats_mallocs_total 4.6346229103e+10
# HELP go_memstats_mcache_inuse_bytes Number of bytes in use by mcache structures.
# TYPE go_memstats_mcache_inuse_bytes gauge
go_memstats_mcache_inuse_bytes 4800
# HELP go_memstats_mcache_sys_bytes Number of bytes used for mcache structures obtained from system.
# TYPE go_memstats_mcache_sys_bytes gauge
go_memstats_mcache_sys_bytes 15600
# HELP go_memstats_mspan_inuse_bytes Number of bytes in use by mspan structures.
# TYPE go_memstats_mspan_inuse_bytes gauge
go_memstats_mspan_inuse_bytes 649584
# HELP go_memstats_mspan_sys_bytes Number of bytes used for mspan structures obtained from system.
# TYPE go_memstats_mspan_sys_bytes gauge
go_memstats_mspan_sys_bytes 1.090224e+06
# HELP go_memstats_next_gc_bytes Number of heap bytes when next garbage collection will take place.
# TYPE go_memstats_next_gc_bytes gauge
go_memstats_next_gc_bytes 5.1911216e+07
# HELP go_memstats_other_sys_bytes Number of bytes used for other system allocations.
# TYPE go_memstats_other_sys_bytes gauge
go_memstats_other_sys_bytes 1.050192e+06
# HELP go_memstats_stack_inuse_bytes Number of bytes in use by the stack allocator.
# TYPE go_memstats_stack_inuse_bytes gauge
go_memstats_stack_inuse_bytes 2.031616e+06
# HELP go_memstats_stack_sys_bytes Number of bytes obtained from system for stack allocator.
# TYPE go_memstats_stack_sys_bytes gauge
go_memstats_stack_sys_bytes 2.031616e+06
# HELP go_memstats_sys_bytes Number of bytes obtained from system.
# TYPE go_memstats_sys_bytes gauge
go_memstats_sys_bytes 9.5012896e+07
# HELP go_sched_gomaxprocs_threads The current runtime.GOMAXPROCS setting, or the number of operating system threads that can execute user-level Go code simultaneously.
# TYPE go_sched_gomaxprocs_threads gauge
go_sched_gomaxprocs_threads 4
# HELP go_sched_goroutines_goroutines Count of live goroutines.
# TYPE go_sched_goroutines_goroutines gauge
go_sched_goroutines_goroutines 90
# HELP go_sched_latencies_seconds Distribution of the time goroutines have spent in the scheduler in a runnable state before actually running.
# TYPE go_sched_latencies_seconds histogram
go_sched_latencies_seconds_bucket{le="-5e-324"} 0
go_sched_latencies_seconds_bucket{le="9.999999999999999e-10"} 2.09033415e+08
go_sched_latencies_seconds_bucket{le="9.999999999999999e-09"} 2.09033415e+08
go_sched_latencies_seconds_bucket{le="9.999999999999998e-08"} 2.11586382e+08
go_sched_latencies_seconds_bucket{le="1.0239999999999999e-06"} 5.18157946e+08
go_sched_latencies_seconds_bucket{le="1.0239999999999999e-05"} 5.3220367e+08
go_sched_latencies_seconds_bucket{le="0.00010239999999999998"} 5.47874447e+08
go_sched_latencies_seconds_bucket{le="0.0010485759999999998"} 5.50538585e+08
go_sched_latencies_seconds_bucket{le="0.010485759999999998"} 5.50741317e+08
go_sched_latencies_seconds_bucket{le="0.10485759999999998"} 5.50746078e+08
go_sched_latencies_seconds_bucket{le="+Inf"} 5.50746129e+08
go_sched_latencies_seconds_sum NaN
go_sched_latencies_seconds_count 5.50746129e+08
# HELP go_threads Number of OS threads created.
# TYPE go_threads gauge
go_threads 14
# HELP process_cpu_seconds_total Total user and system CPU time spent in seconds.
# TYPE process_cpu_seconds_total counter
process_cpu_seconds_total 92942.1
# HELP process_max_fds Maximum number of open file descriptors.
# TYPE process_max_fds gauge
process_max_fds 1.048576e+06
# HELP process_open_fds Number of open file descriptors.
# TYPE process_open_fds gauge
process_open_fds 16
# HELP process_resident_memory_bytes Resident memory size in bytes.
# TYPE process_resident_memory_bytes gauge
process_resident_memory_bytes 9.7079296e+07
# HELP process_start_time_seconds Start time of the process since unix epoch in seconds.
# TYPE process_start_time_seconds gauge
process_start_time_seconds 1.69236739071e+09
# HELP process_virtual_memory_bytes Virtual memory size in bytes.
# TYPE process_virtual_memory_bytes gauge
process_virtual_memory_bytes 1.767079936e+09
# HELP process_virtual_memory_max_bytes Maximum amount of virtual memory available in bytes.
# TYPE process_virtual_memory_max_bytes gauge
process_virtual_memory_max_bytes 1.8446744073709552e+19
# HELP promhttp_metric_handler_requests_in_flight Current number of scrapes being served.
# TYPE promhttp_metric_handler_requests_in_flight gauge
promhttp_metric_handler_requests_in_flight 1
# HELP promhttp_metric_handler_requests_total Total number of scrapes by HTTP status code.
# TYPE promhttp_metric_handler_requests_total counter
promhttp_metric_handler_requests_total{code="200"} 588020
promhttp_metric_handler_requests_total{code="500"} 0
promhttp_metric_handler_requests_total{code="503"} 0
```
