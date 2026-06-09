# Kubernetes / Etcd

Etcd Dashboard for Prometheus metrics scraper

## Tags

* `k8s`
* `etcd`

## Panels

<!-- markdownlint-disable line-length -->
| Name | Description | Thresholds | Repeat |
| ---- | ----------- | ---------- | ------ |
| Etcd has a leader? | Indicates whether the members have a leader.<br/><br/>Yes - all members have a leader.<br/>No - one or more members haven't a leader. | Default:<br/>Mode: absolute<br/>Level 1: 1<br/><br/> |  |
| Leader name | Indicate which server is the leader | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Disk operations wal_fsync duration | The latency distributions of fsync called by wal |  |  |
| Count of leader changes per $interval | Show the number of leader changes per $interval |  |  |
| The sum of leader changes | Show the sum of leader changes per $interval | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| Etcd nodes status | Shows status of each etcd pod in the cluster. | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| The sum rate of failed proposals seen | Show the sum rate of failed proposals seen per $interval | Default:<br/>Mode: absolute<br/>Level 1: 80<br/><br/> |  |
| RPC Rate | Total number of gRPC's started and failed on the server |  |  |
| Active Streams | Show active streams |  |  |
| DB Size | Total size of the underlying database |  |  |
| Disk Sync Duration | The latency distributions of fsync called by wal and backend |  |  |
| Memory | Resident memory size |  |  |
| Client Traffic In | The total number of bytes received from grpc clients |  |  |
| Client Traffic Out | The total number of bytes sent to grpc clients |  |  |
| Peer Traffic In | The total number of bytes received from peers |  |  |
| Peer Traffic Out | The total number of bytes sent to peers |  |  |
| Raft Proposals | The total number of failed proposals seen<br/><br/>The current number of pending proposals to commit<br/><br/>The total number of consensus proposals applied, committed |  |  |
| Total Leader Elections Per Day | The number of leader changes seen per day |  |  |
| The total number of consensus proposals committed | proposals_committed_total records the total number of consensus proposals committed. This gauge should increase over time if the cluster is healthy. Several healthy members of an etcd cluster may have different total committed proposals at once. This discrepancy may be due to recovering from peers after starting, lagging behind the leader, or being the leader and therefore having the most commits. It is important to monitor this metric across all the members in the cluster; a consistently large lag between a single member and its leader indicates that member is slow or unhealthy.<br/><br/>proposals_applied_total records the total number of consensus proposals applied. The etcd server applies every committed proposal asynchronously. The difference between proposals_committed_total and proposals_applied_total should usually be small (within a few thousands even under high load). If the difference between them continues to rise, it indicates that the etcd server is overloaded. This might happen when applying expensive queries like heavy range queries or large txn operations. |  |  |
| Proposals pending | indicates how many proposals are queued to commit. Rising pending proposals suggests there is a high client load or the member cannot commit proposals. |  |  |
| Disks operations | The sum rate latency distributions of fsync called by wal and backend |  |  |
| Slow Operations | The count of slow apply and read index operations  |  |  |
| Network | The rate number of bytes received and sent from grpc clients |  |  |
| Snapshot Duration | Abnormally high snapshot duration (snapshot_save_total_duration_seconds) indicates disk issues and might cause the cluster to be unstable. |  |  |
| Snapshot Fsync Duration | The latency distributions of fsync called by snap |  |  |
| Snapshot DB File Duration | The latency distributions of saving and fsyncing `.snap.db` files.  |  |  |
<!-- markdownlint-enable line-length -->
