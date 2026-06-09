# Kubernetes / Pod Resources

There is no description on this dashboard

## Tags

* `k8s`
* `pods`

## Panels

### Overview (Namespace)

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Pods Count (Namespace) | Show the number of pods running in the namespace |  |  |
| Container status (Namespace) | Show the count of containers in  Waiting, and Terminated statuses. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Last Restart Time | Shows the time of the last restart of the pod in the selected time range |  |  |
| Last Deployment Time | Shows the last time of pod deployed |  |  |
| Pod Status | Shows pod states |  |  |
| Container Restarts (Namespace) | Shows restarts of each container in the namespace |  |  |
<!-- markdownlint-enable line-length -->

### Container status (Namespace)

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Not healthy pods (Namespace) | Show information about the reason the container is currently in waiting or terminated state | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Last terminated status (Namespace) | Show information about the last reason the container was in terminated state |  |  |
<!-- markdownlint-enable line-length -->

### CPU

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| CPU Usage $pod | Shows CPU usage per container in pod against resource requests and limits |  | **Panel is multiplied by parameter `pod`** |
| CPU Quota $pod | Shows CPU quota usage by pod |  | **Panel is multiplied by parameter `pod`** |
<!-- markdownlint-enable line-length -->

### Memory

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Memory Usage $pod | Shows memory usage per container in pod against resource requests and limits |  | **Panel is multiplied by parameter `pod`** |
| Memory Quota $pod | Shows memory quota usage by pod | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> | **Panel is multiplied by parameter `pod`** |
<!-- markdownlint-enable line-length -->

### Disk

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Disk I/O utilization $pod | Shows data read/write per container |  | **Panel is multiplied by parameter `pod`** |
| Conteiner IOps $pod | Count of writes/reads completed |  | **Panel is multiplied by parameter `pod`** |
<!-- markdownlint-enable line-length -->

### Network

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Receive/Transmit Bandwidth | Shows network traffic in bytes per second for the pod |  |  |
| Rate of Received/Transmitted Packets | Shows network packets for pod |  |  |
| Rate of Received/Transmitted Packets Dropped | Shows dropped packets for pod |  |  |
| Open Sockets | Number of open sockets for the container.<br/>NOTE: Panel shows information about open sockets only for Kubernetes v1.21+ and Openshift v4.x. Openshift v3.11 and less doesn't support metrics about sockets. |  |  |
<!-- markdownlint-enable line-length -->

### File descriptor

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Help | Show information about panels in current section |  |  |
| File Descriptors | Number of open file descriptors for the container.  |  |  |
| Rate of File Descriptors | Rate of open file descriptors for the container.  |  |  |
| Ulimit Root Process | Soft ulimit values for the container root process. Unlimited if -1. |  |  |
| Threads | Number of threads running inside the container. |  |  |
| Rate of Threads | Rate of threads running inside the container. |  |  |
| Maximum Threads | Maximum number of threads allowed inside the container, infinity if value is zero. |  |  |
<!-- markdownlint-enable line-length -->

### Throttled pods

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Throttled pods, rate by seconds | Show pods which were throttled by CPU limits.<br/><br/>Show as rate by `container_cpu_cfs_throttled_seconds_total`. It means that this graph should show how many time pod spend in throttling per second. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
<!-- markdownlint-enable line-length -->
