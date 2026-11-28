{{- define "defaultAlerts" -}}
KubebernetesAlerts:
  labels:
    group_name: KubebernetesAlerts
  interval: 30s
  concurrency: 2
  rules:
    KubernetesNodeReady:
      expr: kube_node_status_condition{condition="Ready",status="true"} == 0
      for: 5m
      labels:
        severity: critical
      annotations:
        summary: "Kubernetes Node ready (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "Node {{ "{{" }} $labels.node {{ "}}" }} has been unready for a long time\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    KubernetesMemoryPressure:
      expr: kube_node_status_condition{condition="MemoryPressure",status="true"} == 1
      for: 5m
      labels:
        severity: critical
      annotations:
        summary: "Kubernetes memory pressure (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "{{ "{{" }} $labels.node {{ "}}" }} has MemoryPressure condition\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    KubernetesDiskPressure:
      expr: kube_node_status_condition{condition="DiskPressure",status="true"} == 1
      for: 5m
      labels:
        severity: critical
      annotations:
        summary: "Kubernetes disk pressure (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "{{ "{{" }} $labels.node {{ "}}" }} has DiskPressure condition\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    KubernetesOutOfDisk:
      expr: kube_node_status_condition{condition="OutOfDisk",status="true"} == 1
      for: 5m
      labels:
        severity: critical
      annotations:
        summary: "Kubernetes out of disk (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "{{ "{{" }} $labels.node {{ "}}" }} has OutOfDisk condition\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    KubernetesJobFailed:
      expr: kube_job_status_failed > 0
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "Kubernetes Job failed (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "Job {{ "{{" }} $labels.namespace {{ "}}" }}/{{ "{{" }} $labels.exported_job {{ "}}" }} failed to complete\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    KubernetesCronjobSuspended:
      expr: kube_cronjob_spec_suspend != 0
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "Kubernetes CronJob suspended (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "CronJob {{ "{{" }} $labels.namespace {{ "}}" }}/{{ "{{" }} $labels.cronjob {{ "}}" }} is suspended\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    KubernetesPersistentvolumeclaimPending:
      expr: kube_persistentvolumeclaim_status_phase{phase="Pending"} == 1
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "Kubernetes PersistentVolumeClaim pending (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "PersistentVolumeClaim {{ "{{" }} $labels.namespace {{ "}}" }}/{{ "{{" }} $labels.persistentvolumeclaim {{ "}}" }} is pending\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    KubernetesPersistentvolumeError:
      expr: (kube_persistentvolume_status_phase{phase=~"Failed|Pending",job="kube-state-metrics"}) > 0
      for: 5m
      labels:
        severity: critical
      annotations:
        summary: "Kubernetes PersistentVolume error (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "Persistent volume is in bad state\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    KubernetesVolumeOutOfDiskSpaceWarning:
      expr: (kubelet_volume_stats_available_bytes / kubelet_volume_stats_capacity_bytes) * 100 < 25
      for: 2m
      labels:
        severity: warning
      annotations:
        summary: Kubernetes Volume out of disk space (instance {{ "{{" }} $labels.instance {{ "}}" }})
        description: "Volume is almost full (< 25% left)\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    KubernetesVolumeOutOfDiskSpaceHigh:
      expr: (kubelet_volume_stats_available_bytes / kubelet_volume_stats_capacity_bytes) * 100 < 10
      for: 2m
      labels:
        severity: high
      annotations:
        summary: Kubernetes Volume out of disk space (instance {{ "{{" }} $labels.instance {{ "}}" }})
        description: "Volume is almost full (< 10% left)\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    KubernetesVolumeFullInFourDays:
      expr: predict_linear(kubelet_volume_stats_available_bytes[6h], 345600) < 0
      for: 10m
      labels:
        severity: warning
      annotations:
        summary: Kubernetes Volume full in four days (instance {{ "{{" }} $labels.instance {{ "}}" }})
        description: "{{ "{{" }} $labels.namespace {{ "}}" }}/{{ "{{" }} $labels.persistentvolumeclaim {{ "}}" }} is expected to fill up within four days. Currently {{ "{{" }} $value | humanize {{ "}}" }}% is available.\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    KubernetesStatefulsetDown:
      expr: kube_statefulset_replicas - kube_statefulset_status_replicas_ready != 0
      for: 5m
      labels:
        severity: critical
      annotations:
        summary: "Kubernetes StatefulSet down (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "A StatefulSet went down\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    KubernetesPodNotHealthy:
      expr: min_over_time(sum by (exported_namespace, exported_pod) (kube_pod_status_phase{phase=~"Pending|Unknown|Failed"})[1h:1m]) > 0
      for: 5m
      labels:
        severity: critical
      annotations:
        summary: "Kubernetes Pod not healthy (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "Pod has been in a non-ready state for longer than an hour.\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    KubernetesPodCrashLooping:
      expr: (rate(kube_pod_container_status_restarts_total[15m]) * 60) * 5 > 5
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "Kubernetes pod crash looping (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "Pod {{ "{{" }} $labels.pod {{ "}}" }} is crash looping\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    KubernetesReplicassetMismatch:
      expr: kube_replicaset_spec_replicas - kube_replicaset_status_ready_replicas != 0
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "Kubernetes ReplicasSet mismatch (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "Deployment Replicas mismatch\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    KubernetesDeploymentReplicasMismatch:
      expr: kube_deployment_spec_replicas - kube_deployment_status_replicas_available != 0
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "Kubernetes Deployment replicas mismatch (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "Deployment Replicas mismatch\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    KubernetesStatefulsetReplicasMismatch:
      expr: kube_statefulset_status_replicas_ready - kube_statefulset_status_replicas != 0
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "Kubernetes StatefulSet replicas mismatch (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "A StatefulSet has not matched the expected number of replicas for longer than 15 minutes.\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    KubernetesDeploymentGenerationMismatch:
      expr: kube_deployment_status_observed_generation - kube_deployment_metadata_generation != 0
      for: 5m
      labels:
        severity: critical
      annotations:
        summary: "Kubernetes Deployment generation mismatch (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "A Deployment has failed but has not been rolled back.\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    KubernetesStatefulsetGenerationMismatch:
      expr: kube_statefulset_status_observed_generation - kube_statefulset_metadata_generation != 0
      for: 5m
      labels:
        severity: critical
      annotations:
        summary: "Kubernetes StatefulSet generation mismatch (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "A StatefulSet has failed but has not been rolled back.\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    KubernetesStatefulsetUpdateNotRolledOut:
      expr: max without (revision) (kube_statefulset_status_current_revision unless kube_statefulset_status_update_revision) * (kube_statefulset_replicas != kube_statefulset_status_replicas_updated)
      for: 5m
      labels:
        severity: critical
      annotations:
        summary: "Kubernetes StatefulSet update not rolled out (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "StatefulSet update has not been rolled out.\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    KubernetesDaemonsetRolloutStuck:
      expr: (((kube_daemonset_status_number_ready / kube_daemonset_status_desired_number_scheduled) * 100) < 100) or (kube_daemonset_status_desired_number_scheduled - kube_daemonset_status_current_number_scheduled > 0)
      for: 5m
      labels:
        severity: critical
      annotations:
        summary: "Kubernetes DaemonSet rollout stuck (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "Some Pods of DaemonSet are not scheduled or not ready\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    KubernetesDaemonsetMisscheduled:
      expr: kube_daemonset_status_number_misscheduled > 0
      for: 5m
      labels:
        severity: critical
      annotations:
        summary: "Kubernetes DaemonSet misscheduled (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "Some DaemonSet Pods are running where they are not supposed to run\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    KubernetesCronjobTooLong:
      expr: time() - kube_cronjob_next_schedule_time > 3600
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "Kubernetes CronJob too long (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "CronJob {{ "{{" }} $labels.namespace {{ "}}" }}/{{ "{{" }} $labels.cronjob {{ "}}" }} is taking more than 1h to complete.\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    KubernetesJobCompletion:
      expr: (kube_job_spec_completions - kube_job_status_succeeded > 0) or (kube_job_status_failed > 0)
      for: 5m
      labels:
        severity: critical
      annotations:
        summary: "Kubernetes job completion (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "Kubernetes Job failed to complete\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    KubernetesApiServerErrors:
      expr: (sum(rate(apiserver_request_count{job="kube-apiserver",code=~"(?:5..)$"}[2m])) / sum(rate(apiserver_request_count{job="kube-apiserver"}[2m]))) * 100 > 3
      for: 5m
      labels:
        severity: critical
      annotations:
        summary: "Kubernetes API server errors (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "Kubernetes API server is experiencing high error rate\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    ApiServerRequestsSlow:
      expr: histogram_quantile(0.99, rate(apiserver_request_duration_seconds_bucket{verb!="WATCH"}[5m])) > 0.5
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "API Server requests are slow(instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "HTTP requests slowing down, 99th quantile is over 0.5s for 5 minutes\\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    ControllerWorkQueueDepth:
      expr: sum(workqueue_depth) > 10
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "Controller work queue depth is more than 10 (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "Controller work queue depth is more than 10\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    KubernetesApiClientErrors:
      expr: (sum(rate(rest_client_requests_total{code=~"(4|5).."}[2m])) by (instance, job) / sum(rate(rest_client_requests_total[2m])) by (instance, job)) * 100 > 5
      for: 5m
      labels:
        severity: critical
      annotations:
        summary: "Kubernetes API client errors (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "Kubernetes API client is experiencing high error rate\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    KubernetesClientCertificateExpiresNextWeek:
      expr: (apiserver_client_certificate_expiration_seconds_count{job="kubelet"}) > 0 and histogram_quantile(0.01, sum by (job, le) (rate(apiserver_client_certificate_expiration_seconds_bucket{job="kubelet"}[5m]))) < 604800
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "Kubernetes client certificate expires next week (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "A client certificate used to authenticate to the apiserver is expiring next week.\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    KubernetesClientCertificateExpiresSoon:
      expr: (apiserver_client_certificate_expiration_seconds_count{job="kubelet"}) > 0 and histogram_quantile(0.01, sum by (job, le) (rate(apiserver_client_certificate_expiration_seconds_bucket{job="kubelet"}[5m]))) < 86400
      for: 5m
      labels:
        severity: critical
      annotations:
        summary: "Kubernetes client certificate expires soon (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "A client certificate used to authenticate to the apiserver is expiring in less than 24.0 hours.\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

NodeProcesses:
  labels:
    group_name: NodeProcesses
  rules:
    CountPidsAndThreadOutOfLimit:
      expr: (sum(container_processes) by (node) +  on (node) label_replace(node_processes_threads * on(instance) group_left(nodename) (node_uname_info), "node", "$1", "nodename", "(.+)")) / on (node) label_replace(node_processes_max_processes * on(instance) group_left(nodename) (node_uname_info), "node", "$1", "nodename", "(.+)") * 100 > 80
      for: 5m
      labels:
        severity: high
      annotations:
        summary: "Host high PIDs and Threads usage (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "Sum of node's pids and threads is filling up (< 20% left)\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

NodeExporters:
  labels:
    group_name: NodeExporters
  rules:
    NodeDiskUsageIsMoreThanWarningThreshold:
      annotations:
        description: "Node {{ "{{" }} $labels.node {{ "}}" }} disk usage of {{ "{{" }} $labels.mountpoint {{ "}}" }} is\n  VALUE = {{ "{{" }} $value {{ "}}" }}%"
        summary: "Disk usage on node > 70% (instance {{ "{{" }} $labels.node {{ "}}" }})"
      expr: (node_filesystem_size_bytes{fstype=~"ext.*|xfs", mountpoint !~".*pod.*"} - node_filesystem_free_bytes{fstype=~"ext.*|xfs", mountpoint !~".*pod.*"}) * 100 / (node_filesystem_avail_bytes{fstype=~"ext.*|xfs", mountpoint !~".*pod.*"} + (node_filesystem_size_bytes{fstype=~"ext.*|xfs", mountpoint !~".*pod.*"} - node_filesystem_free_bytes{fstype=~"ext.*|xfs", mountpoint !~".*pod.*"})) > 70
      for: 5m
      labels:
        severity: warning

    NodeDiskUsageIsMoreThanCriticalThreshold:
      annotations:
        description: "Node {{ "{{" }} $labels.node {{ "}}" }} disk usage of {{ "{{" }} $labels.mountpoint {{ "}}" }} is\n VALUE = {{ "{{" }} $value {{ "}}" }}%"
        summary: "Disk usage on node > 90% (instance {{ "{{" }} $labels.node {{ "}}" }})"
      expr: (node_filesystem_size_bytes{fstype=~"ext.*|xfs", mountpoint !~".*pod.*"} - node_filesystem_free_bytes{fstype=~"ext.*|xfs", mountpoint !~".*pod.*"}) * 100 / (node_filesystem_avail_bytes{fstype=~"ext.*|xfs", mountpoint !~".*pod.*"} + (node_filesystem_size_bytes{fstype=~"ext.*|xfs", mountpoint !~".*pod.*"} - node_filesystem_free_bytes{fstype=~"ext.*|xfs", mountpoint !~".*pod.*"})) > 90
      for: 5m
      labels:
        severity: high

    HostOutOfMemory:
      expr: ((node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes) * 100) * on(instance) group_left(nodename) node_uname_info < 10
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "Host out of memory (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "Node memory is filling up (< 10% left)\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    HostMemoryUnderMemoryPressure:
      expr: rate(node_vmstat_pgmajfault[2m]) * on(instance) group_left(nodename) node_uname_info > 1000
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "Host memory under memory pressure (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "The node is under heavy memory pressure. High rate of major page faults\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    HostUnusualNetworkThroughputIn:
      expr: ((sum by (instance) (irate(node_network_receive_bytes_total[2m])) * on(instance) group_left(nodename) node_uname_info) / 1024) / 1024 > 100
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "Host unusual network throughput in (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "Host network interfaces are probably receiving too much data (> 100 MB/s)\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    HostUnusualNetworkThroughputOut:
      expr: ((sum by (instance) (irate(node_network_transmit_bytes_total[2m])) * on(instance) group_left(nodename) node_uname_info) / 1024) / 1024 > 100
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "Host unusual network throughput out (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "Host network interfaces are probably sending too much data (> 100 MB/s)\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    HostUnusualDiskReadRate:
      expr: (sum by (instance) (irate(node_disk_read_bytes_total[2m])) * on(instance) group_left(nodename) node_uname_info) / 1024 / 1024 > 50
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "Host unusual disk read rate (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "Disk is probably reading too much data (> 50 MB/s)\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    HostUnusualDiskWriteRate:
      expr: ((sum by (instance) (irate(node_disk_written_bytes_total[2m])) * on(instance) group_left(nodename) node_uname_info) / 1024) / 1024 > 50
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "Host unusual disk write rate (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "Disk is probably writing too much data (> 50 MB/s)\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    HostOutOfDiskSpace:
      expr: ((node_filesystem_avail_bytes{mountpoint="/"}  * 100) / node_filesystem_size_bytes{mountpoint="/"}) * on(instance) group_left(nodename) node_uname_info < 10
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "Host out of disk space (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "Disk is almost full (< 10% left)\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    HostDiskWillFillIn4Hours:
      expr: predict_linear(node_filesystem_free_bytes{fstype!~"tmpfs"}[1h], 14400) * on(instance) group_left(nodename) node_uname_info < 0
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "Host disk will fill in 4 hours (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "Disk will fill in 4 hours at current write rate\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    HostOutOfInodes:
      expr: ((node_filesystem_files_free{mountpoint ="/"} / node_filesystem_files{mountpoint ="/"}) * 100) * on(instance) group_left(nodename) node_uname_info < 10
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "Host out of inodes (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "Disk is almost running out of available inodes (< 10% left)\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    HostUnusualDiskReadLatency:
      expr: (rate(node_disk_read_time_seconds_total[2m]) / rate(node_disk_reads_completed_total[2m])) * on(instance) group_left(nodename) node_uname_info > 100
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "Host unusual disk read latency (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "Disk latency is growing (read operations > 100ms)\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    HostUnusualDiskWriteLatency:
      expr: (rate(node_disk_write_time_seconds_total[2m]) / rate(node_disk_writes_completed_total[2m])) * on(instance) group_left(nodename) node_uname_info > 100
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "Host unusual disk write latency (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "Disk latency is growing (write operations > 100ms)\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    HostHighCpuLoad:
      expr: 100 - ((avg(irate(node_cpu_seconds_total{mode="idle"}[5m])) by (instance) * 100) * on (instance) group_left (nodename) node_uname_info) > 80
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "Host high CPU load (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "CPU load is > 80%\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

DockerContainers:
  labels:
    group_name: DockerContainers
  rules:
    ContainerKilled:
      expr: time() - container_last_seen > 60
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "Container killed (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "A container has disappeared\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    ContainerVolumeUsage:
      expr: (1 - (sum(container_fs_inodes_free) BY (node) / sum(container_fs_inodes_total) BY (node))) * 100 > 80
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "Container Volume usage (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "Container Volume usage is above 80%\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    ContainerVolumeIoUsage:
      expr: (sum(container_fs_io_current) BY (node, name) * 100) > 80
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "Container Volume IO usage (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "Container Volume IO usage is above 80%\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    ContainerHighThrottleRate:
      expr: rate(container_cpu_cfs_throttled_seconds_total[3m]) > 1
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "Container high throttle rate (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "Container is being throttled\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

HAmode:
  labels:
    group_name: HAmode
  rules:
    NotHAKubernetesDeploymentAvailableReplicas:
      expr: kube_deployment_status_replicas_available < 2
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "Not HA mode: Deployment Available Replicas < 2 (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "Not HA mode: Kubernetes Deployment has less than 2 available replicas\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    NotHAKubernetesStatefulSetAvailableReplicas:
      expr: kube_statefulset_status_replicas_available < 2
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "Not HA mode: StatefulSet Available Replicas < 2 (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "Not HA mode: Kubernetes StatefulSet has less than 2 available replicas\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    NotHAKubernetesDeploymentDesiredReplicas:
      expr: kube_deployment_status_replicas < 2
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "Not HA mode: Deployment Desired Replicas < 2 (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "Not HA mode: Kubernetes Deployment has less than 2 desired replicas\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    NotHAKubernetesStatefulSetDesiredReplicas:
      expr: kube_statefulset_status_replicas < 2
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "Not HA mode: StatefulSet Desired Replicas < 2 (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "Not HA mode: Kubernetes StatefulSet has less than 2 desired replicas\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    NotHAKubernetesDeploymentMultiplePodsPerNode:
      expr: count(sum(kube_pod_info{node=~".+", created_by_kind="ReplicaSet"}) by (namespace, node, created_by_name) > 1) > 0
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "Not HA mode: Deployment Has Multiple Pods per Node (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "Not HA mode: Kubernetes Deployment has 2 or more replicas on the same node\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    NotHAKubernetesStatefulSetMultiplePodsPerNode:
      expr: count(sum(kube_pod_info{node=~".+", created_by_kind="StatefulSet"}) by (namespace, node, created_by_name) > 1) > 0
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "Not HA mode: StatefulSet Has Multiple Pods per Node (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "Not HA mode: Kubernetes StatefulSet has 2 or more replicas on the same node\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

HAproxy:
  labels:
    group_name: HAproxy
  rules:
    HaproxyDown:
      expr: haproxy_up == 0
      for: 5m
      labels:
        severity: critical
      annotations:
        summary: "HAProxy down (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "HAProxy down\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    HaproxyBackendConnectionErrors:
      expr: sum by (backend) (rate(haproxy_backend_connection_errors_total[2m])) > 10
      for: 5m
      labels:
        severity: critical
      annotations:
        summary: "HAProxy backend connection errors (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "Too many connection errors to {{ "{{" }} $labels.fqdn {{ "}}" }}/{{ "{{" }} $labels.backend {{ "}}" }} backend (> 10 req/s). Request throughput may be to high.\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    HaproxyServerResponseErrors:
      expr: sum by (server) (rate(haproxy_server_response_errors_total[2m])) > 5
      for: 5m
      labels:
        severity: critical
      annotations:
        summary: "HAProxy server response errors (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "Too many response errors to {{ "{{" }} $labels.server {{ "}}" }} server (> 5 req/s).\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    HaproxyServerConnectionErrors:
      expr: sum by (server) (rate(haproxy_server_connection_errors_total[2m])) > 10
      for: 5m
      labels:
        severity: critical
      annotations:
        summary: "HAProxy server connection errors (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "Too many connection errors to {{ "{{" }} $labels.server {{ "}}" }} server (> 10 req/s). Request throughput may be to high.\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    HaproxyPendingRequests:
      expr: sum by (backend) (haproxy_backend_current_queue) > 0
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "HAProxy pending requests (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "Some HAProxy requests are pending on {{ "{{" }} $labels.fqdn {{ "}}" }}/{{ "{{" }} $labels.backend {{ "}}" }} backend\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    HaproxyHttpSlowingDown:
      expr: avg by (backend) (haproxy_backend_http_total_time_average_seconds) > 2
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "HAProxy HTTP slowing down (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "Average request time is increasing\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    HaproxyRetryHigh:
      expr: sum by (backend) (rate(haproxy_backend_retry_warnings_total[5m])) > 10
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "HAProxy retry high (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "High rate of retry on {{ "{{" }} $labels.fqdn {{ "}}" }}/{{ "{{" }} $labels.backend {{ "}}" }} backend\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    HaproxyBackendDown:
      expr: haproxy_backend_up == 0
      for: 5m
      labels:
        severity: critical
      annotations:
        summary: "HAProxy backend down (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "HAProxy backend is down\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    HaproxyServerDown:
      expr: haproxy_server_up == 0
      for: 5m
      labels:
        severity: critical
      annotations:
        summary: "HAProxy server down (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "HAProxy server is down\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    HaproxyFrontendSecurityBlockedRequests:
      expr: sum by (frontend) (rate(haproxy_frontend_requests_denied_total[5m])) > 10
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "HAProxy frontend security blocked requests (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "HAProxy is blocking requests for security reason\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    HaproxyServerHealthcheckFailure:
      expr: increase(haproxy_server_check_failures_total[5m]) > 0
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "HAProxy server healthcheck failure (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "Some server healthcheck are failing on {{ "{{" }} $labels.server {{ "}}" }}\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

Etcd:
  labels:
    group_name: Etcd
  rules:
    EtcdInsufficientMembers:
      expr: count(etcd_server_id{job="etcd"}) % 2 == 0
      for: 5m
      labels:
        severity: critical
      annotations:
        summary: "Etcd insufficient Members (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "Etcd cluster should have an odd number of members\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    EtcdNoLeader:
      expr: etcd_server_has_leader == 0
      for: 5m
      labels:
        severity: critical
      annotations:
        summary: "Etcd no Leader (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "Etcd cluster have no leader\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    EtcdHighNumberOfLeaderChanges:
      expr: increase(etcd_server_leader_changes_seen_total[1h]) > 3
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "Etcd high number of leader changes (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "Etcd leader changed more than 3 times during last hour\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    EtcdWarningNumberOfFailedGrpcRequests:
      expr: sum(rate(grpc_server_handled_total{job="etcd",grpc_code!="OK", grpc_method!="Watch"}[5m])) BY (grpc_service, grpc_method) / sum(rate(grpc_server_handled_total{job="etcd"}[5m])) BY (grpc_service, grpc_method) > 0.01
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "Etcd high number of failed GRPC requests (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "More than 1% GRPC request failure detected in Etcd for 5 minutes\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    EtcdCriticalNumberOfFailedGrpcRequests:
      expr: sum(rate(grpc_server_handled_total{job="etcd",grpc_code!="OK", grpc_method!="Watch"}[5m])) BY (grpc_service, grpc_method) / sum(rate(grpc_server_handled_total{job="etcd"}[5m])) BY (grpc_service, grpc_method) > 0.05
      for: 5m
      labels:
        severity: critical
      annotations:
        summary: "Etcd high number of failed GRPC requests (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "More than 5% GRPC request failure detected in Etcd for 5 minutes\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    EtcdGrpcRequestsSlow:
      expr: histogram_quantile(0.99, sum(rate(grpc_server_handling_seconds_bucket{job="etcd",grpc_type="unary"}[5m])) by (grpc_service, grpc_method, le)) > 0.15
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "Etcd GRPC requests slow (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "GRPC requests slowing down, 99th percentil is over 0.15s for 5 minutes\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    EtcdMemberCommunicationSlow:
      expr: histogram_quantile(0.99, rate(etcd_network_peer_round_trip_time_seconds_bucket{job="etcd"}[5m])) > 0.15
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "Etcd member communication slow (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "Etcd member communication slowing down, 99th percentil is over 0.15s for 5 minutes\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    EtcdHighNumberOfFailedProposals:
      expr: increase(etcd_server_proposals_failed_total[1h]) > 5
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "Etcd high number of failed proposals (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "Etcd server got more than 5 failed proposals past hour\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    EtcdHighFsyncDurations:
      expr: histogram_quantile(0.99, rate(etcd_disk_wal_fsync_duration_seconds_bucket[5m])) > 0.5
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "Etcd high fsync durations (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "Etcd WAL fsync duration increasing, 99th percentil is over 0.5s for 5 minutes\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    EtcdHighCommitDurations:
      expr: histogram_quantile(0.99, rate(etcd_disk_backend_commit_duration_seconds_bucket[5m])) > 0.25
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "Etcd high commit durations (instance {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "Etcd commit duration increasing, 99th percentil is over 0.25s for 5 minutes\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

NginxIngressAlerts:
  labels:
    group_name: NginxIngressAlerts
  rules:
    NginxHighHttp4xxErrorRate:
      expr: sum by (ingress, exported_namespace, node) (rate(nginx_ingress_controller_requests{status=~"^4.."}[2m])) / sum by (ingress, exported_namespace, node)(rate(nginx_ingress_controller_requests[2m])) * 100 > 5
      for: 1m
      labels:
        severity: high
      annotations:
        summary: "Nginx high HTTP 4xx error rate (node: {{ "{{" }} $labels.node {{ "}}" }}, namespace: {{ "{{" }} $labels.exported_namespace {{ "}}" }}, ingress: {{ "{{" }} $labels.ingress {{ "}}" }})"
        description: "Too many HTTP requests with status 4xx (> 5%)\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS = {{ "{{" }} $labels {{ "}}" }}"

    NginxHighHttp5xxErrorRate:
      expr: sum by (ingress, exported_namespace, node) (rate(nginx_ingress_controller_requests{status=~"^5.."}[2m])) / sum by (ingress, exported_namespace, node) (rate(nginx_ingress_controller_requests[2m])) * 100 > 5
      for: 1m
      labels:
        severity: high
      annotations:
        summary: "Nginx high HTTP 5xx error rate (node: {{ "{{" }} $labels.node {{ "}}" }}, namespace: {{ "{{" }} $labels.exported_namespace {{ "}}" }}, ingress: {{ "{{" }} $labels.ingress {{ "}}" }})"
        description: "Too many HTTP requests with status 5xx (> 5%)\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS = {{ "{{" }} $labels {{ "}}" }}"

    NginxLatencyHigh:
      expr: histogram_quantile(0.99, sum(rate(nginx_ingress_controller_request_duration_seconds_bucket[2m])) by (host, node, le)) > 3
      for: 2m
      labels:
        severity: warning
      annotations:
        summary: "Nginx latency high (node: {{ "{{" }} $labels.node {{ "}}" }}, host: {{ "{{" }} $labels.host {{ "}}" }})"
        description: "Nginx p99 latency is higher than 3 seconds\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS = {{ "{{" }} $labels {{ "}}" }}"

CoreDnsAlerts:
  labels:
    group_name: CoreDnsAlerts
  rules:
    CorednsPanicCount:
      expr: increase(coredns_panics_total[1m]) > 0
      for: 0m
      labels:
        severity: critical
      annotations:
        summary: CoreDNS Panic Count (instance {{ "{{" }} $labels.instance {{ "}}" }})
        description: "Number of CoreDNS panics encountered\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS = {{ "{{" }} $labels {{ "}}" }}"

    CoreDNSLatencyHigh:
      annotations:
        description: CoreDNS has 99th percentile latency of {{ "{{" }} $value {{ "}}" }} seconds for server {{ "{{" }} $labels.server {{ "}}" }} zone {{ "{{" }} $labels.zone {{ "}}" }}
        summary: CoreDNS have High Latency
      expr: histogram_quantile(0.99, sum(rate(coredns_dns_request_duration_seconds_bucket[2m])) by(server, zone, le)) > 3
      for: 5m
      labels:
        severity: critical

    CoreDNSForwardHealthcheckFailureCount:
      annotations:
        summary: CoreDNS health checks have failed to upstream server
        description: CoreDNS health checks have failed to upstream server {{ "{{" }} $labels.to {{ "}}" }}
      expr: sum(rate(coredns_forward_healthcheck_broken_total[2m])) > 0
      for: 5m
      labels:
        severity: warning

    CoreDNSForwardHealthcheckBrokenCount:
      annotations:
        summary: CoreDNS health checks have failed for all upstream servers
        description: "CoreDNS health checks failed for all upstream servers LABELS = {{ "{{" }} $labels {{ "}}" }}"
      expr: sum(rate(coredns_forward_healthcheck_broken_total[2m])) > 0
      for: 5m
      labels:
        severity: warning

    CoreDNSErrorsCritical:
      annotations:
        description: CoreDNS is returning SERVFAIL for {{ "{{" }} $value | humanizePercentage {{ "}}" }} of requests
        summary: CoreDNS is returning SERVFAIL
      expr: sum(rate(coredns_dns_responses_total{rcode="SERVFAIL"}[2m])) / sum(rate(coredns_dns_responses_total[2m])) > 0.03
      for: 5m
      labels:
        severity: critical
        
    CoreDNSErrorsWarning:
      annotations:
        description: CoreDNS is returning SERVFAIL for {{ "{{" }} $value | humanizePercentage {{ "}}" }} of requests
        summary: CoreDNS is returning SERVFAIL
      expr: sum(rate(coredns_dns_responses_total{rcode="SERVFAIL"}[2m])) / sum(rate(coredns_dns_responses_total[2m])) > 0.01
      for: 5m
      labels:
        severity: warning

    CoreDNSForwardLatencyHigh:
      annotations:
        description: CoreDNS has 99th percentile latency of {{ "{{" }} $value {{ "}}" }} seconds forwarding requests to {{ "{{" }} $labels.to {{ "}}" }}
        summary: CoreDNS has 99th percentile latency for forwarding requests
      expr: histogram_quantile(0.99, sum(rate(coredns_forward_request_duration_seconds_bucket[2m])) by(to, le)) > 3
      for: 5m
      labels:
        severity: critical
        
    CoreDNSForwardErrorsCritical:
      annotations:
        description: CoreDNS is returning SERVFAIL for {{ "{{" }} $value | humanizePercentage {{ "}}" }} of forward requests to {{ "{{" }} $labels.to {{ "}}" }}
        summary: CoreDNS is returning SERVFAIL for forward requests
      expr: sum(rate(coredns_forward_responses_total{rcode="SERVFAIL"}[2m])) / sum(rate(coredns_forward_responses_total[2m])) > 0.03
      for: 5m
      labels:
        severity: critical
        
    CoreDNSForwardErrorsWarning:
      annotations:
        description: CoreDNS is returning SERVFAIL for {{ "{{" }} $value | humanizePercentage {{ "}}" }} of forward requests to {{ "{{" }} $labels.to {{ "}}" }}
        summary: CoreDNS is returning SERVFAIL for forward requests
      expr: sum(rate(coredns_forward_responses_total{rcode="SERVFAIL"}[2m])) / sum(rate(coredns_forward_responses_total[2m])) > 0.01
      for: 5m
      labels:
        severity: warning

DRAlerts:
  labels:
    group_name: DRAlerts
  rules:
    ProbeFailed:
      expr: probe_success == 0
      for: 5m
      labels:
        severity: critical
      annotations:
        summary: "Probe failed (instance: {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "Probe failed\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    SlowProbe:
      expr: avg_over_time(probe_duration_seconds[1m]) > 1
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "Slow probe (instance: {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "Blackbox probe took more than 1s to complete\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    HttpStatusCode:
      expr: probe_http_status_code <= 199 OR probe_http_status_code >= 400
      for: 5m
      labels:
        severity: high
      annotations:
        summary: "HTTP Status Code (instance: {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "HTTP status code is not 200-399\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

    HttpSlowRequests:
      expr: avg_over_time(probe_http_duration_seconds[1m]) > 1
      for: 5m
      labels:
        severity: warning
      annotations:
        summary: "HTTP slow requests (instance: {{ "{{" }} $labels.instance {{ "}}" }})"
        description: "HTTP request took more than 1s\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"

BackupAlerts:
  labels:
    group_name: BackupAlerts
  rules:
    Last Backup Failed:
      expr: backup_storage_last_failed != 0
      for: 1m
      labels:
        severity: warning
      annotations:
        summary: "Last backup made by pod {{ "{{" }} $labels.pod {{ "}}" }} in namespace {{ "{{" }} $labels.namespace {{ "}}" }} failed.\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"
        description: "Last backup made by pod {{ "{{" }} $labels.pod {{ "}}" }} in namespace {{ "{{" }} $labels.namespace {{ "}}" }} failed.\n  VALUE = {{ "{{" }} $value {{ "}}" }}\n  LABELS: {{ "{{" }} $labels {{ "}}" }}"
 {{- end }}


