# Kubernetes / Nodes Resources

There is no description on this dashboard

## Tags

* `k8s`
* `k8s-nodes`

## Panels

### Overview

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Nodes Overview | Shows cluster nodes overview: <br/>  * Node Uptime<br/>  * Total available CPU and RAM on node<br/>  * Overall resources usage on node<br/>  * Can be grouped by node_label<br/><br/>! This panel doesn't work in restricted mode if kube-state-metrics doesn't have permission to scrape the following resources: nodes | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Overall CPU Load And Usage | Shows overall CPU average load on all nodes and CPU usage against total number of CPU cores<br/>  * Can be grouped by node_label<br/><br/>! This panel doesn't work in restricted mode if kube-state-metrics doesn't have permission to scrape the following resources: nodes |  |  |
| Overall Memory Usage | Shows overall RAM usage for all nodes against total available RAM on all nodes<br/>  * Can be grouped by node_label<br/><br/>! This panel doesn't work in restricted mode if kube-state-metrics doesn't have permission to scrape the following resources: nodes |  |  |
| Overall Disk Usage | Shows overall disk space usage for all nodes  against total available disk space on all nodes<br/>  * Can be grouped by node_label<br/><br/>! This panel doesn't work in restricted mode if kube-state-metrics doesn't have permission to scrape the following resources: nodes |  |  |
| Disk I/O utilization  per node | Shows data  read/write per node<br/>  * Can be grouped by node_label<br/><br/>! This panel doesn't work in restricted mode if kube-state-metrics doesn't have permission to scrape the following resources: nodes |  |  |
| IOps per node | Count of writes/reads completed<br/>  * Can be grouped by node_label<br/><br/>! This panel doesn't work in restricted mode if kube-state-metrics doesn't have permission to scrape the following resources: nodes |  |  |
<!-- markdownlint-enable line-length -->

### Node: $node

**Row Node: $node is multiplied by parameter `node`**

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Uptime | Show the node uptime | Default:<br/>Mode: absolute<br/>Level 1: 1<br/>Level 2: 3<br/><br/> |  |
| CPU Cores Quota | Shows total available CPU cores on the node | Default:<br/>Mode: absolute<br/>Level 1: 1<br/>Level 2: 2<br/><br/> |  |
| *TODO: Add panel name* | Shows brief info about resources usage on the node: <br/>* CPU<br/>* RAM<br/>* Disk Space<br/>* Swap | Default:<br/>Mode: absolute<br/>Level 1: 70<br/>Level 2: 90<br/><br/> |  |
| Disk Space Used Basic(EXT?/XFS) | Shows the node disk info: <br/>* Space usage<br/>* Available disk space,<br/>* Disk usage rate<br/><br/>The data is coming from the 'df' command. The 'Used' column algorithm in df is: (size-free) * 100 / (avail + (size-free)), the result is divisible by this value, non-divisible by this value is +1, and the unit of the result is % | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| CPU iowait | Shows the node CPU iowait as a percentage of total available CPU on the node | Default:<br/>Mode: absolute<br/>Level 1: 20<br/>Level 2: 50<br/><br/> |  |
| Memory Quota | Shows total available RAM on the node | Default:<br/>Mode: absolute<br/>Level 1: 2<br/>Level 2: 3<br/><br/> |  |
| Used PIDs / Threads | Shows the node PIDs and threads current used count | Default:<br/>Mode: absolute<br/>Level 1: 32000<br/><br/> |  |
| Max PIDs / Threads | Shows the node PIDs and threads limits per node | Default:<br/>Mode: absolute<br/>Level 1: 10000000<br/><br/> |  |
| CPU Usage | Shows detailed information about CPU usage on the node |  |  |
| Memory Usage | Shows detailed information about memory usage on the node |  |  |
| System Load | Shows average system load for 1m, 5m, 15m on the node |  |  |
| Internet Traffic Per $interval For $interface | Shows incoming and outcoming internet traffic in bytes on the node grouped by the time specified in the 'Interval' variable. |  |  |
| Network Bandwidth Per Second For $interface | Shows network outcoming and incoming traffic in bytes per seconds on the node |  |  |
| Network Sockstat | Shows sockets information on the node:<br/><br/>* Sockets_used - sockets currently in use<br/>* CurrEstab - TCP connections for which the current state is either ESTABLISHED or CLOSE- WAIT<br/>* TCP_alloc - Allocated sockets<br/>* TCP_tw - Sockets wating close<br/>* UDP_inuse - Udp sockets currently in use<br/>* RetransSegs - TCP retransmission packets<br/>* OutSegs - Number of packets sent by TCP<br/>* InSegs - Number of packets received by TCP |  |  |
| Disk Space Usage | Shows disk space usage by the mounts selected in the 'Mount' variable on the node |  |  |
| Disk R/W Data (bytes / sec) | Shows per second disk read / write bytes on the node |  |  |
| Open  File Descriptor / Context switches | Shows current amount of opened filed descriptors and context switches for the node |  |  |
| Time Spent Doing I/Os | Shows time spent on I/O in the natural time of each second |  |  |
| Disk IOps | Shows amount of physical IOs (reads and writes) on different devices. |  |  |
| I/O Usage Read / Write | Shows the volume of disks operations (read and write) |  |  |
| I/O Usage Times | Total seconds spent doing I/Os |  |  |
| PIDs Number and Limit | Shows amount of used processes (PIDs) and process limit on node |  |  |
| Threads Number and Limit | Shows amount of used threads and threads limit on node |  |  |
<!-- markdownlint-enable line-length -->
