# Kubernetes / Cluster Overview

There is no description on this dashboard

## Tags

* `k8s`
* `k8s-cluster`

## Panels

### Overview

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| API Servers | Shows the number of running API servers depending on the number of endpoints exposing metrics | Default:<br/>Mode: absolute<br/>Level 1: 1<br/>Level 2: 2<br/><br/> |  |
| ETCD Servers | Shows number of active ETCD servers. May contain no data for PaaS clouds. | Default:<br/>Mode: absolute<br/>Level 1: 1<br/>Level 2: 3<br/><br/> |  |
| kubelet Nodes | Shows number of active kubelet servers |  |  |
| Ingress Pods | Shows count of running monitoring pods<br/><br/>Note: Metrics may not be collected from Openshift cluster |  |  |
| Containers Status | Show the count of containers in Running, Waiting, and Terminated statuses<br/><br/>Note: Terminated pods not include "Completed" status |  |  |
| Nodes Status | Shows count of running and not running nodes in the cluster. May contain no data for PaaS clouds. |  |  |
| Nodes Transition Status | Shows transition status of each node in the cluster<br/><br/>1 - OK, 0 - Problem<br/><br/>! This panel doesn't work in restricted mode if kube-state-metrics doesn't have permission to scrape the following resources: nodes |  |  |
| Persistent Volumes Status | Shows persistent volumes count grouped by statuses |  |  |
| CPU Usage | Shows CPU usage and allocation for the whole cluster<br/><br/>! This panel doesn't work in restricted mode if kube-state-metrics doesn't have permission to scrape the following resources: nodes | Default:<br/>Mode: absolute<br/>Level 1: 0.75<br/>Level 2: 0.9<br/><br/> |  |
| Storage | Shows usage disk space for all mounted volumes to nodes in the cluster | Default:<br/>Mode: absolute<br/>Level 1: 0.6<br/>Level 2: 0.8<br/><br/> |  |
| Memory Usage | Shows memory utilization and allocation for the whole cluster<br/><br/>! This panel doesn't work in restricted mode if kube-state-metrics doesn't have permission to scrape the following resources: nodes | Default:<br/>Mode: absolute<br/>Level 1: 0.7<br/>Level 2: 0.9<br/><br/> |  |
| Not Healthy Pods | Show information about the reason the container is currently in waiting or terminated state | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Last Terminated Status | Show information about the last reason the container was in terminated state | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Resources

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Daemon Sets | Shows the counts of running/all daemon sets in the cluster |  |  |
| Stateful Sets | Shows the counts of running/all stateful sets in the cluster |  |  |
| Deployments | Shows the counts of running/all deployments in the cluster |  |  |
| Jobs | Shows the counts of running/all jobs in the cluster |  |  |
| Pods | Shows the counts of running/all pods in the cluster |  |  |
| Replica Sets | Shows the counts of running/all replica sets in the cluster |  |  |
| Replication Controllers | Shows the counts of running/all replication controllers in the cluster |  |  |
| Cron Jobs | Shows the counts of running/all cron jobs in the cluster |  |  |
<!-- markdownlint-enable line-length -->

### Ingress

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Controller Request Volume | Shows number of operation per second<br/><br/>Note: Metrics may not be collected from Openshift cluster |  |  |
| Ingress Request Volume | Show the total number of client requests by ingress<br/><br/>Note: Metrics may not be collected from Openshift cluster |  |  |
| Ingress Controller Success Rate | Shows the percentage of successful requests with status codes 2xx-3xx<br/><br/>Note: Metrics may not be collected from Openshift cluster | Default:<br/>Mode: percentage<br/>Level 1: 70<br/>Level 2: 85<br/>Level 3: 90<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Certificates

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Certificates from files | Shows count of total, expiring in 30 days and expired certificates from files.<br/>To get the information about certificates cert-exporter must be running. |  |  |
| Certificates from kubeconfig | Shows count of total, expiring in 30 days and expired certificates from kubeconfig. To get the information about certificates cert-exporter must be running. |  |  |
| Certificates from secrets | Shows count of total, expiring in 30 days and expired certificates from secrets. To get the information about certificates cert-exporter must be running. |  |  |
<!-- markdownlint-enable line-length -->

### CPU

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| CPU Usage | Shows CPU usage by each namespace in the cluster |  |  |
| CPU Quota | Shows CPU requests and limits quota allocation per namespace in the cluster |  |  |
<!-- markdownlint-enable line-length -->

### Memory

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Memory Usage | Shows memory usage by namespace in the cluster |  |  |
| Memory Quota | Shows memory requests and limits quota allocation per namespace in the cluster |  |  |
<!-- markdownlint-enable line-length -->

### Disk

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Disk I/O utilization  per namspaces | Shows data read/write per namespace |  |  |
|  IOps per namespace | Count of writes/reads completed |  |  |
| Disk I/O utilization  per node | Shows data read/write per node |  |  |
|  IOps per node | Count of writes/reads completed |  |  |
| Device storage usage | Shows the amounts of total and used device storage in the cluster |  |  |
<!-- markdownlint-enable line-length -->

### Volumes

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Count PVC | Show the number of pvc in the cluster |  |  |
| Error: Persistent Volumes Not Bound | Show the number of PV with status != "Bound" in the cluster |  |  |
| Error: Persistent Volumes Claim Not Bound | Show the number of pvc with status != "Bound" in the cluster |  |  |
| Count PV | Show the number of pv in the cluster |  |  |
<!-- markdownlint-enable line-length -->

### Network

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Current Network Usage | Shows network usage details per namespace in the cluster |  |  |
| Network RX/TX Total | Shows overall network input and output traffic |  |  |
| Receive/Transmit Bandwidth | Shows network traffic per namespace in the cluster |  |  |
| Average Container Bandwidth by Namespace: Received/Transmitted | Shows network traffic per container in the cluster |  |  |
| Rate of Received/Transmitted Packets | Shows average packets rate per namespace in the cluster |  |  |
| Rate of Received/Transmitted Packets Dropped | Shows dropped packets per namespaces in the cluster |  |  |
<!-- markdownlint-enable line-length -->
