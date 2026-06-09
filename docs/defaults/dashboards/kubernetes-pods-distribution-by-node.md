# Kubernetes / Pods Distribution By Node

There is no description on this dashboard

## Tags

* `k8s`
* `nodes`
* `pods`

## Panels

### Node: $node

**Row Node: $node is multiplied by parameter `node`**

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Number of pods | Shows the number of pods running on the node | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| CPU | Shows CPU usage on the node |  |  |
| Memory | Shows memory usage on the node |  |  |
| CPU Usage | Shows detailed CPU usage per pod on the node |  |  |
| CPU Quota | Shows CPU usage, requests and limits per pod on the node | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Memory Usage | Shows detailed memory usage per pod on the node |  |  |
| Memory Quota | Shows memory usage, requests and limits per pod on the node | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Network

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Receive Bytes By Pods | Shows received bytes by pods in current namespace |  |  |
| Transmit Bytes By Pods | Shows transmitted bytes by pods in current namespace |  |  |
| Receive Packets / Drops By Pods | Shows received packets and received packets drops by pods in current namespace |  |  |
| Transmit Packets / Drops By Pods | Shows transmitted packets and transmitted packets drops by pods in current namespace |  |  |
| Receive/Transmit Bandwidth | Shows overall network receive/transmit bytes in current node |  |  |
| Network By Pod | Shows general information about the network on the current node |  |  |
| Open Sockets By Pod | Number of open sockets for the pod.<br/>NOTE: Panel shows information about open sockets only for Kubernetes v1.21+ and Openshift v4.x. Openshift v3.11 and less doesn't support metrics about sockets. |  |  |
<!-- markdownlint-enable line-length -->

### Disk

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Read Bytes Per Pods | Shows the average filesystem reads in bytes for the current node |  |  |
| Written Bytes Per Pods | Shows the average filesystem writes in bytes for the current node |  |  |
| Reads Per Pods | Shows the average successful reads operation for the current node |  |  |
| Writes Per Pods | Shows the average successful writes operation for the current node |  |  |
| Disk By Pod | Shows general information about the disk on the current node | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Space

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Help | Information about panels on this row |  |  |
| CSI PVs space usage | Shows information about space consumption for persistent volumes powered by CSI drivers (except NFS storages) | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Space Info | Shows the filesystem usage for each device and mountpoint |  |  |
<!-- markdownlint-enable line-length -->

### File descriptor

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Help | Show information about panels in current section |  |  |
| File Descriptors | Number of open file descriptors for the pod.  |  |  |
| Threads | Number of threads running inside the pod. |  |  |
| File Descriptors Quota | Shows file descriptors, threads and limits per pod and container on the node. |  |  |
<!-- markdownlint-enable line-length -->
