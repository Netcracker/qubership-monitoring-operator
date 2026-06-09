# Kubernetes / Kubelet

There is no description on this dashboard

## Tags

* `k8s`
* `kubelet`

## Panels

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Up Nodes | Shows the number of alive nodes in the cluster by running kubelet instances | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Running Pods | Shows the total number of running pods on selected nodes | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Running Containers | Shows the total number of running containers on selected nodes | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Actual Volume Count | Shows the total number of volumes in Volume Manager for selected nodes.<br/><br/>**NOTE**:  This metric is available only in Kubernetes > 1.12 (Openshift > 4.1) | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Desired Volume Count | Show the desired number of volumes in Volume Manager for selected nodes.<br/><br/>**NOTE**:  This metric is available only in Kubernetes > 1.12 (Openshift > 4.1) | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Config Error Count | Shows the number of configuration-related errors for selected nodes | Default:<br/>Mode: absolute<br/>Level 1: 1<br/>Level 2: 1<br/><br/> |  |
| Operation Rate | Shows a total count of runtime operations of each type for selected nodes |  |  |
| Operation Error Rate | Shows a count of errors in the operations for selected nodes |  |  |
| Operation duration 99th quantile | Shows duration of the operations summary for selected nodes.<br/><br/>**NOTE**:  This metric is available only in Kubernetes > 1.12 (Openshift > 4.1) |  |  |
| Pod Start Rate | Shos number of pod start operations for selected nodes.<br/><br/>**NOTE**:  This metric is available only in Kubernetes > 1.12 (Openshift > 4.1) |  |  |
| Pod Start Duration | Shows pod start time summary for selected nodes.<br/><br/>**NOTE**:  This metric is available only in Kubernetes > 1.12 (Openshift > 4.1) |  |  |
| Storage Operation Rate | Shows storage operation duration for selected nodes |  |  |
| Storage Operation Error Rate | Shows the number of storage operations errors for selected nodes |  |  |
| Storage Operation Duration 99th quantile | Shows storage operations duration summary per opeartion for selected nodes |  |  |
| Memory | Shows memory usage by kubelet on selected nodes |  |  |
| CPU usage | Shows CPU usage by kubelet on selected nodes |  |  |
| Goroutines | Shows number of gorutines in kubelet on selected nodes |  |  |
| Cgroup manager operation rate | Shows duration in seconds for cgroup manager operations.<br/><br/>**NOTE**:  This metric is available only in Kubernetes > 1.12 (Openshift > 4.1) |  |  |
| Cgroup manager 99th quantile | Duration in seconds for cgroup manager operations.<br/><br/>**NOTE**:  This metric is available only in Kubernetes > 1.12 (Openshift > 4.1) |  |  |
<!-- markdownlint-enable line-length -->
