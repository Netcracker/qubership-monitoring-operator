This document describes kubelet metrics list and how to collect them.

# Metrics

kubelet already exposes its metrics in Prometheus format and doesn't require using specific exporters.

<!-- markdownlint-disable line-length -->
| Name                  | Metrics Port | Metrics Endpoint    | Need Exporter? | Auth?         | Is Exporter Third Party? |
| --------------------- | ------------ | ------------------- | -------------- | ------------- | ------------------------ |
| Prometheus (kubelet)  | `10250`      | `/metrics`          | No             | Require, RBAC | N/A                      |
| Prometheus (cAdvisor) | `10250`      | `/metrics/cadvisor` | No             | Require, RBAC | N/A                      |
<!-- markdownlint-enable line-length -->

## How to Collect

Metrics expose on port `10250` and endpoint `/metrics` for kubelet and `/metrics/cadvisor` for cadvisor metrics.
It requires kubernetes authentication which placed in `/var/run/secrets/kubernetes.io/serviceaccount/token`.

Before configure Prometheus or VMAgent to collect these metrics you need to create a Service that will route to
`kubelet`. It can be created in any namespace. In case of using the `prometheus-operator` it automatically create
a Service in a namespace where will be deploy the Prometheus.

It config:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: kubelet
  labels:
    k8s-app: kubelet
spec:
  ports:
    - name: https-metrics
      protocol: TCP
      port: 10250
      targetPort: 10250
    - name: http-metrics
      protocol: TCP
      port: 10255
      targetPort: 10255
    - name: cadvisor
      protocol: TCP
      port: 4194
      targetPort: 4194
  clusterIP: None
  clusterIPs:
    - None
  type: ClusterIP
  sessionAffinity: None
  ipFamilies:
    - IPv4
    - IPv6
  ipFamilyPolicy: RequireDualStack
```

Config `ServiceMonitor` for `prometheus-operator` to collect kubelet metrics:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: kubelet-service-monitor
  labels:
    k8s-app: kubelet-service-monitor
    app.kubernetes.io/name: kubelet-service-monitor
    app.kubernetes.io/component: monitoring
    app.kubernetes.io/managed-by: monitoring-operator
spec:
  endpoints:
    - interval: 30s
      scrapeTimeout: 10s
      metricRelabelings:
        - sourceLabels: ['pod_name']
          targetLabel: pod
          regex: (.+)
        - sourceLabels: ['container_name']
          targetLabel: container
          regex: (.+)
        - action: labeldrop
          regex: pod_name
        - action: labeldrop
          regex: container_name
        - regex: 'kubelet_running_pods'
          replacement: 'kubelet_running_pod_count'
          sourceLabels: ['__name__']
          targetLabel: __name__
        - regex: 'kubelet_running_containers'
          replacement: 'kubelet_running_container_count'
          sourceLabels: ['__name__']
          targetLabel: __name__
      relabelings: []
      bearerTokenFile: /var/run/secrets/kubernetes.io/serviceaccount/token
      honorLabels: true
      port: https-metrics
      scheme: https
      tlsConfig:
        caFile: /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
        insecureSkipVerify: true
    - interval: 30s
      scrapeTimeout: 10s
      metricRelabelings:
        - sourceLabels: ['pod_name']
          targetLabel: pod
          regex: (.+)
        - sourceLabels: ['container_name']
          targetLabel: container
          regex: (.+)
        - action: labeldrop
          regex: pod_name
        - action: labeldrop
          regex: container_name
        - regex: 'kubelet_running_pods'
          replacement: 'kubelet_running_pod_count'
          sourceLabels: ['__name__']
          targetLabel: __name__
        - regex: 'kubelet_running_containers'
          replacement: 'kubelet_running_container_count'
          sourceLabels: ['__name__']
          targetLabel: __name__
      relabelings: []
      bearerTokenFile: /var/run/secrets/kubernetes.io/serviceaccount/token
      honorLabels: true
      path: /metrics/cadvisor
      port: https-metrics
      scheme: https
      tlsConfig:
        caFile: /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
        insecureSkipVerify: true
  jobLabel: k8s-app
  selector:
    matchLabels:
      k8s-app: kubelet
```

To collect (or just to check) metrics manually you can use the following command:

```bash
curl -v -k -L -H "Authorization: Bearer $(cat /var/run/secrets/kubernetes.io/serviceaccount/token)" https://<kube_node_ip_or_dns>:10250/metrics
```

You can't use `wget` because it doesn't allow to add headers for authorization.

## Metrics List

Metrics list of `kubelet` you can find in files:

* [Kubelet](../../defaults/metrics.md#kubelet-kubernetes-v123)
* [cAdvisor](../../defaults/metrics.md#cadvisor-kubernetes-v123)

## Volumes usage metrics

Some kubelet metrics `kubelet_volume_stats_*` may not work properly because get information from different CSI drivers
of volumes.

The most CSI drivers, or Provisioners support collecting of the metrics by default. Some of its have ability to turn it on.

### Built-in storages

There are some type  of volumes that built-in in the Kubernetes.

<!-- markdownlint-disable line-length -->
| Storage Type | Support metrics | Turned on by default | How to turn on collecting |
| ------------ | --------------- | -------------------- | ------------------------- |
| HostPath     | yes             | yes                  | Can't be disabled         |
| ConfigMap    | no              | no                   | Metrics are not supported |
| Secret       | no              | no                   | Metrics are not supported |
| Projected    | no              | no                   | Metrics are not supported |
| EmptyDir     | no              | no                   | Metrics are not supported |
<!-- markdownlint-enable line-length -->

### CSI Drivers

There is a list of the CSI drivers and information about support the metrics of volume statistics collecting:

**NOTE:** Some CSI Drivers can expose only some of `kubelet_volume_stats_*` metrics. For example, `AWS EFS` expose
only `kubelet_volume_stats_used_bytes` (`capacity` and `available` metrics are not available). So please keep in
mind it and verify before using.

<!-- markdownlint-disable line-length -->
| Storage Type           | CSI Driver                 | Support metrics | Turned on by default | How to turn on collecting                                                                                                                                                                                                                                               |
| ---------------------- | -------------------------- | --------------- | -------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| AWS EBS                | `ebs.csi.aws.com`          | yes             | yes                  | -                                                                                                                                                                                                                                                                       |
| AWS EFS                | `efs.csi.aws.com`          | yes             | no                   | Set in helm [values.yaml](https://github.com/kubernetes-sigs/aws-efs-csi-driver/blob/master/charts/aws-efs-csi-driver/values.yaml#L56) `controller.volMetricsOptIn=true` set `--vol-metrics-opt-in=true` in DaemonSet sci node                                          |
| Azure Blob             | `blob.csi.azure.com`       | yes             | no                   | Set in helm [values.yaml](https://github.com/kubernetes-sigs/blob-csi-driver/blob/master/charts/latest/blob-csi-driver/values.yaml#L155) `feature.enableGetVolumeStats=true`, or flag in DaemonSet sci node container's args `--enable-get-volume-stats=true`           |
| Azure Disk             | `disk.csi.azure.com`       | yes             | yes                  | -                                                                                                                                                                                                                                                                       |
| Azure File             | `file.csi.azure.com`       | yes             | yes                  | Set in helm [values.yaml](https://github.com/kubernetes-sigs/azurefile-csi-driver/blob/master/charts/latest/azurefile-csi-driver/values.yaml#L158) `feature.enableGetVolumeStats=true`, or flag in DaemonSet sci node container's args `--enable-get-volume-stats=true` |
| Ceph RBD               | `rbd.csi.ceph.com`         | yes             | yes                  | -                                                                                                                                                                                                                                                                       |
| CephFS                 | `cephfs.csi.ceph.com`      | yes             | yes                  | -                                                                                                                                                                                                                                                                       |
| Cinder                 | `cinder.csi.openstack.org` | yes             | yes                  | -                                                                                                                                                                                                                                                                       |
| GCE Persistent Disk    | `pd.csi.storage.gke.io`    | yes             | yes                  | -                                                                                                                                                                                                                                                                       |
| Google Cloud Filestore | `com.google.csi.filestore` | yes             | yes                  | -                                                                                                                                                                                                                                                                       |
| Google Cloud Storage   | `gcs.csi.ofek.dev`         | no              | -                    | -                                                                                                                                                                                                                                                                       |
| NFS                    | `nfs.csi.k8s.io`           | yes             | yes                  | -                                                                                                                                                                                                                                                                       |
<!-- markdownlint-enable line-length -->

### External Provisioners

Expect of built-in volume providers and CDI Drivers also exists External Provisioners:

* [https://kubernetes.io/docs/concepts/storage/storage-classes/#provisioner](https://kubernetes.io/docs/concepts/storage/storage-classes/#provisioner)

<!-- markdownlint-disable line-length -->
| Storage Type | Provisioner             | Support metrics | Turned on by default | How to turn on collecting                                                                                                      |
| ------------ | ----------------------- | --------------- | -------------------- | ------------------------------------------------------------------------------------------------------------------------------ |
| LocalPath    | `rancher.io/local-path` | no              | no                   | Currently not supported. See the issue [volumes usage collecting](https://github.com/rancher/local-path-provisioner/issues/44) |
<!-- markdownlint-enable line-length -->
