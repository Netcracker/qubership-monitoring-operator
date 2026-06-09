# Overall Cloud Status

Dashboard shows health status of applications are deployed into cloud platform, k8s/OpenShift nodes,
applications are deployed out of cloud.

## Tags

* `k8s`
* `health`

## Panels

### Kubernetes overview

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| API server status | Shows status of Kubernetes API server. | Default:<br/>Mode: absolute<br/>Level 1: 1<br/><br/> |  |
| API servers | Shows number of API servers. | Default:<br/>Mode: absolute<br/>Level 1: 2<br/>Level 2: 3<br/><br/> |  |
| API server requests | Shows count of requests to API server, requests per minute. |  |  |
| API server errors | Shows errors in requests to API server. | Default:<br/>Mode: absolute<br/>Level 1: 1<br/>Level 2: 3<br/><br/> |  |
| ETCD status | Show status of etcd cluster. May contain no data for PaaS clouds. |  |  |
| ETCD servers | Shows number of  active ETCD servers. May contain no data for PaaS clouds. | Default:<br/>Mode: absolute<br/>Level 1: 1<br/>Level 2: 3<br/><br/> |  |
| ETCD requests | Shows number of requests per second to ETCD servers. May contain no data for PaaS clouds. | Default:<br/>Mode: absolute<br/>Level 1: 1<br/>Level 2: 500<br/><br/> |  |
| ETCD server request error | Shows percent of error requests to ETCD server. May contain no data for PaaS clouds. | Default:<br/>Mode: absolute<br/>Level 1: 1<br/>Level 2: 3<br/><br/> |  |
| API server nodes status | Shows status of each API server in the cluster.  1 - OK, 0 - Problem |  |  |
| API server failed requests | Shows errors  in requests to API server, operations per minute. |  |  |
| Etcd nodes status | Shows status of each etcd pod in the cluster. 1 - OK, 0 - Problem.<br/>May contain no data for PaaS clouds. |  |  |
| ETCD failed requests | Shows number  of errors per minute in requests to ETCD server. May contain no data for PaaS clouds. |  |  |
| Total CPU usage | Shows overall CPU usage | Default:<br/>Mode: absolute<br/>Level 1: 75<br/>Level 2: 90<br/><br/> |  |
| Total Memory usage | Shows overall RAM usage for all nodes against total available RAM on all nodes. | Default:<br/>Mode: absolute<br/>Level 1: 75<br/>Level 2: 90<br/><br/> |  |
| Total Filesystem usage | Shows summary file system usage on Kubernetes cluster nodes | Default:<br/>Mode: absolute<br/>Level 1: 75<br/>Level 2: 90<br/><br/> |  |
| Used cores | Show used cores for cloud in cores (1 core = 1000 millicores) |  |  |
| Total cores | Show total cores available for cloud |  |  |
| Used memory | Show total used memory for cloud |  |  |
| Total memory | Show total available memory for cloud |  |  |
| Used space | Show sum by used space for directories and files on all nodes in cloud where `fstype == xfs &#124; ext.`. It means that all FS like `tmpfs`, `rootfs` will be exclude from value. |  |  |
| Total space | Show total available space for directories and files on all nodes in cloud where `fstype == xfs &#124; ext.`. It means that all FS like `tmpfs`, `rootfs` will be exclude from value. |  |  |
| Number of nodes | Shows number of active Kubernetes cluster nodes | Default:<br/>Mode: absolute<br/>Level 1: 1<br/>Level 2: 3<br/><br/> |  |
| Nodes Unavailable | Shows number of unavailable nodes. | Default:<br/>Mode: absolute<br/>Level 1: 1<br/><br/> |  |
| Running Pods | Shows the total number of running pods in cluster. Show only pods with `status = ready` |  |  |
| Running containers | Shows the total number of running containers in pods. |  |  |
<!-- markdownlint-enable line-length -->

### Node health

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Node State | Show running state of all nodes in selected cloud | Default:<br/>Mode: absolute<br/>Level 1: 1<br/>Level 2: 2<br/><br/> |  |
| Nodes Overview | Shows cluster nodes overview: <br/>  * Node Uptime<br/>  * Total available CPU and RAM on node<br/>  * Overall resources usage on node<br/>  * Can be grouped by node_label | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->

### Applications health

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Total pods | Shows the total number of pods in cluster. |  |  |
| Running pods | Shows the total number of running pods in cluster. Show only pods with `status = ready` |  |  |
| Not runnning pods | Shows the total number of not running / not healthy pods in cluster. | Default:<br/>Mode: absolute<br/>Level 1: 1<br/><br/> |  |
| Help | Show information about panels in current section |  |  |
| Not Healthy Pods | Show information about the reason the container is currently in waiting or terminated state | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Last Terminated Status | Show information about the last reason the container was in terminated state | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->
