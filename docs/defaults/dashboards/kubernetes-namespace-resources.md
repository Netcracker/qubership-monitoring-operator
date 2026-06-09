# Kubernetes / Namespace Resources

There is no description on this dashboard

## Tags

* `k8s`
* `k8s-namespaces`

## Panels

### Overview

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Pods Count | Shows the total number of pods in the namespace |  |  |
| Not Running Pods Count | Shows the number of pods in Failed, Pending, and Unknown state in the namespace | Default:<br/>Mode: absolute<br/>Level 1: 1<br/><br/> |  |
| Containers Count | Shows the total number of containers in the namespace |  |  |
| Container Status | Show the count of containers in Running, Waiting, and Terminated statuses<br/><br/>Note: Terminated not include "Completed" status | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Last Terminated Status | Show information about the last reason the container was in terminated state |  |  |
| Not Healthy Pods | Show information about the reason the container is currently in waiting or terminated state | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| CPU Resources Summary | Shows a general summary of CPU resources |  |  |
| Memory Resources Summary | Shows a general summary of Memory resources |  |  |
<!-- markdownlint-enable line-length -->

### Counts Of Resources

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Cron Jobs | Shows the counts of running/all cron jobs in the current namespace |  |  |
| Daemon Sets | Shows the counts of running/all daemon sets in the current namespace |  |  |
| Deployments | Shows the counts of running/all deployments in the current namespace |  |  |
| Jobs | Shows the counts of running/all jobs in the current namespace |  |  |
| Pods | Shows the counts of running/all pods in the current namespace |  |  |
| Replica Sets | Shows the counts of running/all replica sets in the current namespace |  |  |
| Replication Controllers | Shows the counts of running/all replication controllers in the current namespace |  |  |
| Stateful Sets | Shows the counts of running/all stateful sets in the current namespace |  |  |
| Config And Storage | Shows the counts of configs and storages in the current namespace |  |  |
| Service | Shows the counts of services in the current namespace<br/><br/>Note: Do not show the routes |  |  |
<!-- markdownlint-enable line-length -->

### CPU And Memory Quota

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| CPU Limits Quota Allocated | Shows percentage of allocated CPU Limits Resource Quota in the namespace<br/><br/>Shows N/A if the quota is not specified | Default:<br/>Mode: absolute<br/>Level 1: 0.5<br/>Level 2: 0.8<br/><br/> |  |
| CPU Requests Quota Allocated | Shows percentage of allocated CPU Requests Resource Quota in the namespace<br/><br/>Shows N/A if the quota is not specified | Default:<br/>Mode: absolute<br/>Level 1: 0.9<br/><br/> |  |
| CPU Limits Quota Actual Usage | Shows percentage of CPU usage against CPU Limits Resource Quota in the namespace<br/><br/>Shows N/A if the quota is not specified | Default:<br/>Mode: absolute<br/>Level 1: 0.5<br/>Level 2: 0.8<br/><br/> |  |
| CPU Requests Quota Actual Usage | Shows percentage of CPU usage against CPU Requests Resource Quota in the namespace<br/><br/>Shows N/A if the quota is not specified | Default:<br/>Mode: absolute<br/>Level 1: 0.5<br/>Level 2: 0.8<br/><br/> |  |
| Memory Limits Quota Allocated | Shows percentage of allocated Memory Limits Resource Quota in the namespace<br/><br/>Shows N/A if the quota is not specified | Default:<br/>Mode: absolute<br/>Level 1: 0.5<br/>Level 2: 0.8<br/><br/> |  |
| Memory Requests Quota Allocated | Shows percentage of allocated Memory Requests Resource Quota in the namespace<br/><br/>Shows N/A if the quota is not specified | Default:<br/>Mode: absolute<br/>Level 1: 0.9<br/><br/> |  |
| Memory Limits Quota Actual Usage | Shows percentage of memory usage against Memory Limits Resource Quota in the namespace<br/><br/>Shows N/A if the quota is not specified | Default:<br/>Mode: absolute<br/>Level 1: 0.5<br/>Level 2: 0.8<br/><br/> |  |
| Memory Requests Quota Actual Usage | Shows percentage of memory usage against Memory Requests Resource Quota in the namespace<br/><br/>Shows N/A if the quota is not specified | Default:<br/>Mode: absolute<br/>Level 1: 0.5<br/>Level 2: 0.8<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### CPU Usage

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| CPU Usage | Shows CPU usage by pods in the namespace |  |  |
| CPU Resources | Shows CPU usage in the namespace against CPU quota in the namespace |  |  |
| CPU Usage By Pods | Shows detailed information about CPU usage by pods in the namespace<br/><br/>Click on the pod name to go to the dashboard with the pod detailed information |  |  |
<!-- markdownlint-enable line-length -->

### Memory Usage

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Memory Usage | Shows memory usage by pods in the namespace |  |  |
| Memory Resources | Shows memory usage in the namespace against memory quota in the namespace |  |  |
| Memory Usage By Pods | Shows detailed information about memory usage by pods in the namespace<br/><br/>Click on the pod name to go to the dashboard with the pod detailed information |  |  |
<!-- markdownlint-enable line-length -->

### Network

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Receive/Transmit Bandwidth | Shows the overall incoming and outgoing network traffic in bytes per second |  |  |
| Receive/Transmit Bandwidth Usage By Pods | Shows the overall incoming and outgoing network traffic in bytes per second by pods |  |  |
| Receive/Transmit Bandwidth Usage By Pods | Shows the overall incoming and outgoing network traffic in bytes per second by pods<br/><br/>Click on the pod name to go to the dashboard with the pod detailed information |  |  |
| Rate of Received/Transmitted Packets | Shows the overall incoming and outgoing packets per second for the namespace |  |  |
| Rate of Received/Transmitted Packets By Pods | Shows the overall incoming and outgoing packets per second by pods for the namespace |  |  |
| Rate of Received/Transmitted/Drops Packets By Pods | Shows overall incoming and outgoing network traffic in bytes per second by pods<br/><br/>Click on the pod name to go to the dashboard with the pod detailed information |  |  |
<!-- markdownlint-enable line-length -->

### Disk

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Disk I/O Utilization | Shows data read/write in bytes |  |  |
| Disk I/O Utilization By Pods | Shows data read/write in bytes by pods |  |  |
| Disk I/O Utilization By Pods | Shows data read/write in bytes by pods |  |  |
| IOps | Count of writes/reads completed |  |  |
| IOps By Pods | Count of writes/reads completed by pods |  |  |
| IOps By Pods | Count of writes/reads completed by pods |  |  |
<!-- markdownlint-enable line-length -->

### Volumes

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| $pvc | Show information about use volumes space in percentage. Cinder plugin doesn't provide "kubelet_volume_stats_*" metrics if Kubernetes version < 1.19v  | Default:<br/>Mode: percentage<br/>Level 1: 70<br/>Level 2: 90<br/><br/> | **Panel is multiplied by parameter `pvc`** |
| Volume Used in bytes | Show information about use volumes space. Cinder plugin doesn't provide "kubelet_volume_stats_*" metrics if Kubernetes version < 1.19v  |  |  |
<!-- markdownlint-enable line-length -->
