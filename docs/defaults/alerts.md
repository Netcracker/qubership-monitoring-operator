# Alerts OOB

* [Alerts OOB](#alerts-oob)
  * [Monitoring-operator](#monitoring-operator)
    * [Heartbeat](#heartbeat)
    * [SelfMonitoring](#selfmonitoring)
    * [AlertManager](#alertmanager)
    * [KubernetesAlerts](#kubernetesalerts)
    * [NodeProcesses](#nodeprocesses)
    * [NodeExporters](#nodeexporters)
    * [DockerContainers](#dockercontainers)
    * [HAmode](#hamode)
    * [HAproxy](#haproxy)
    * [Etcd](#etcd)
    * [NginxIngressAlerts](#nginxingressalerts)
    * [CoreDnsAlerts](#corednsalerts)
    * [DRAlerts](#dralerts)
    * [BackupAlerts](#backupalerts)
  * [Cert-exporter](#cert-exporter)

## Monitoring-operator

### Heartbeat

#### Alerting rules

<!-- markdownlint-disable line-length -->
| Name | Summary | For | Severity | Expression | Description | Other labels | Other annotations |
| ---- | ------- | --- | -------- | ---------- | ----------- | ------------ | ----------------- |
| DeadMansSwitch | An always-firing Dead Man's Switch alert (instance {{ $labels.instance }}) | 3m | information | vector(1) | This is an alert meant to ensure that the entire alerting pipeline is functional. This alert should always be firing.<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
<!-- markdownlint-enable line-length -->

### SelfMonitoring

#### Alerting rules

<!-- markdownlint-disable line-length -->
| Name | Summary | For | Severity | Expression | Description | Other labels | Other annotations |
| ---- | ------- | --- | -------- | ---------- | ----------- | ------------ | ----------------- |
| PrometheusJobMissing | Prometheus job missing (instance {{ $labels.instance }}) | 5m | warning | absent(up{job=\~".\*prometheus\-pod\-monitor"}) | A Prometheus job has disappeared<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| PrometheusTargetMissing | Prometheus target missing (instance {{ $labels.instance }}) | 5m | high | up == 0 | A Prometheus target has disappeared. An exporter might be crashed.<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| PrometheusAllTargetsMissing | Prometheus all targets missing (job {{ $labels.job }}) | 5m | critical | count by(job) (up) == count by(job) (up == 0) | A Prometheus job does not have living target anymore.<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| PrometheusConfigurationReloadFailure | Prometheus configuration reload failure (instance {{ $labels.instance }}) | 5m | warning | prometheus_config_last_reload_successful != 1 | Prometheus configuration reload error<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| PrometheusTooManyRestarts | Prometheus too many restarts (instance {{ $labels.instance }}) | 5m | warning | changes(process_start_time_seconds{job=\~".\*prometheus\-pod\-monitor"}[15m]) \> 2 | Prometheus has restarted more than twice in the last 15 minutes. It might be crashlooping.<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| PrometheusRuleEvaluationFailures | Prometheus rule evaluation failures (instance {{ $labels.instance }}) | 5m | critical | increase(prometheus_rule_evaluation_failures_total[3m]) \> 0 | Prometheus encountered {{ $value }} rule evaluation failures, leading to potentially ignored alerts.<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| PrometheusTemplateTextExpansionFailures | Prometheus template text expansion failures (instance {{ $labels.instance }}) | 5m | critical | increase(prometheus_template_text_expansion_failures_total[3m]) \> 0 | Prometheus encountered {{ $value }} template text expansion failures<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| PrometheusRuleEvaluationSlow | Prometheus rule evaluation slow (instance {{ $labels.instance }}) | 5m | warning | prometheus_rule_group_last_duration_seconds \> prometheus_rule_group_interval_seconds | Prometheus rule evaluation took more time than the scheduled interval. I indicates a slower storage backend access or too complex query.<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| PrometheusNotificationsBacklog | Prometheus notifications backlog (instance {{ $labels.instance }}) | 5m | warning | min_over_time(prometheus_notifications_queue_length[10m]) \> 0 | The Prometheus notification queue has not been empty for 10 minutes<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| PrometheusTargetEmpty | Prometheus target empty (instance {{ $labels.instance }}) | 5m | critical | prometheus_sd_discovered_targets == 0 | Prometheus has no target in service discovery<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| PrometheusTargetScrapingSlowTwoMinutes | Prometheus target scraping slow (instance {{ $labels.instance }}) for 2 minutes | 5m | warning | (prometheus_target_interval_length_seconds{interval="2m0s", quantile="0.9"}) \> 135 | Prometheus is scraping exporters slowly<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| PrometheusTargetScrapingSlowOneMinute | Prometheus target scraping slow (instance {{ $labels.instance }}) for 1 minute | 5m | warning | (prometheus_target_interval_length_seconds{interval="1m0s", quantile="0.9"}) \> 70 | Prometheus is scraping exporters slowly<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| PrometheusTargetScrapingSlowThirtySeconds | Prometheus target scraping slow (instance {{ $labels.instance }}) for 30 seconds | 5m | warning | (prometheus_target_interval_length_seconds{interval="30s", quantile="0.9"}) \> 35 | Prometheus is scraping exporters slowly<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| PrometheusLargeScrape | Prometheus large scrape (instance {{ $labels.instance }}) | 5m | warning | increase(prometheus_target_scrapes_exceeded_sample_limit_total[10m]) \> 10 | Prometheus has many scrapes that exceed the sample limit<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| PrometheusTargetScrapeDuplicate | Prometheus target scrape duplicate (instance {{ $labels.instance }}) | 5m | warning | increase(prometheus_target_scrapes_sample_duplicate_timestamp_total[5m]) \> 0 | Prometheus has many samples rejected due to duplicate timestamps but different values<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| PrometheusTsdbCheckpointCreationFailures | Prometheus TSDB checkpoint creation failures (instance {{ $labels.instance }}) | 5m | critical | increase(prometheus_tsdb_checkpoint_creations_failed_total[3m]) \> 0 | Prometheus encountered {{ $value }} checkpoint creation failures<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| PrometheusTsdbCheckpointDeletionFailures | Prometheus TSDB checkpoint deletion failures (instance {{ $labels.instance }}) | 5m | critical | increase(prometheus_tsdb_checkpoint_deletions_failed_total[3m]) \> 0 | Prometheus encountered {{ $value }} checkpoint deletion failures<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| PrometheusTsdbCompactionsFailed | Prometheus TSDB compactions failed (instance {{ $labels.instance }}) | 5m | critical | increase(prometheus_tsdb_compactions_failed_total[3m]) \> 0 | Prometheus encountered {{ $value }} TSDB compactions failures<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| PrometheusTsdbHeadTruncationsFailed | Prometheus TSDB head truncations failed (instance {{ $labels.instance }}) | 5m | critical | increase(prometheus_tsdb_head_truncations_failed_total[3m]) \> 0 | Prometheus encountered {{ $value }} TSDB head truncation failures<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| PrometheusTsdbReloadFailures | Prometheus TSDB reload failures (instance {{ $labels.instance }}) | 5m | critical | increase(prometheus_tsdb_reloads_failures_total[3m]) \> 0 | Prometheus encountered {{ $value }} TSDB reload failures<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| PrometheusTsdbWalCorruptions | Prometheus TSDB WAL corruptions (instance {{ $labels.instance }}) | 5m | critical | increase(prometheus_tsdb_wal_corruptions_total[3m]) \> 0 | Prometheus encountered {{ $value }} TSDB WAL corruptions<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| PrometheusTsdbWalTruncationsFailed | Prometheus TSDB WAL truncations failed (instance {{ $labels.instance }}) | 5m | critical | increase(prometheus_tsdb_wal_truncations_failed_total[3m]) \> 0 | Prometheus encountered {{ $value }} TSDB WAL truncation failures<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
<!-- markdownlint-enable line-length -->

### AlertManager

#### Alerting rules

<!-- markdownlint-disable line-length -->
| Name | Summary | For | Severity | Expression | Description | Other labels | Other annotations |
| ---- | ------- | --- | -------- | ---------- | ----------- | ------------ | ----------------- |
| PrometheusAlertmanagerConfigurationReloadFailure | Prometheus AlertManager configuration reload failure (instance {{ $labels.instance }}) | 5m | warning | alertmanager_config_last_reload_successful != 1 | AlertManager configuration reload error<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| PrometheusNotConnectedToAlertmanager | Prometheus not connected to alertmanager (instance {{ $labels.instance }}) | 5m | critical | prometheus_notifications_alertmanagers_discovered < 1 | Prometheus cannot connect the alertmanager<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| PrometheusAlertmanagerNotificationFailing | Prometheus AlertManager notification failing (instance {{ $labels.instance }}) | 5m | high | rate(alertmanager_notifications_failed_total[2m]) \> 0 | Alertmanager is failing sending notifications<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
<!-- markdownlint-enable line-length -->

### KubernetesAlerts

#### Alerting rules

<!-- markdownlint-disable line-length -->
| Name | Summary | For | Severity | Expression | Description | Other labels | Other annotations |
| ---- | ------- | --- | -------- | ---------- | ----------- | ------------ | ----------------- |
| KubernetesNodeReady | Kubernetes Node ready (instance {{ $labels.instance }}) | 5m | critical | kube_node_status_condition{condition="Ready",status="true"} == 0 | Node {{ $labels.node }} has been unready for a long time<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| KubernetesMemoryPressure | Kubernetes memory pressure (instance {{ $labels.instance }}) | 5m | critical | kube_node_status_condition{condition="MemoryPressure",status="true"} == 1 | {{ $labels.node }} has MemoryPressure condition<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| KubernetesDiskPressure | Kubernetes disk pressure (instance {{ $labels.instance }}) | 5m | critical | kube_node_status_condition{condition="DiskPressure",status="true"} == 1 | {{ $labels.node }} has DiskPressure condition<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| KubernetesOutOfDisk | Kubernetes out of disk (instance {{ $labels.instance }}) | 5m | critical | kube_node_status_condition{condition="OutOfDisk",status="true"} == 1 | {{ $labels.node }} has OutOfDisk condition<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| KubernetesJobFailed | Kubernetes Job failed (instance {{ $labels.instance }}) | 5m | warning | kube_job_status_failed \> 0 | Job {{$labels.namespace}}/{{$labels.exported_job}} failed to complete<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| KubernetesCronjobSuspended | Kubernetes CronJob suspended (instance {{ $labels.instance }}) | 5m | warning | kube_cronjob_spec_suspend != 0 | CronJob {{ $labels.namespace }}/{{ $labels.cronjob }} is suspended<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| KubernetesPersistentvolumeclaimPending | Kubernetes PersistentVolumeClaim pending (instance {{ $labels.instance }}) | 5m | warning | kube_persistentvolumeclaim_status_phase{phase="Pending"} == 1 | PersistentVolumeClaim {{ $labels.namespace }}/{{ $labels.persistentvolumeclaim }} is pending<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| KubernetesPersistentvolumeError | Kubernetes PersistentVolume error (instance {{ $labels.instance }}) | 5m | critical | (kube_persistentvolume_status_phase{phase=\~"Failed&#124;Pending",job="kube\-state\-metrics"}) \> 0 | Persistent volume is in bad state<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| KubernetesVolumeOutOfDiskSpaceWarning | Kubernetes Volume out of disk space (instance {{ $labels.instance }}) | 2m | warning | (kubelet_volume_stats_available_bytes / kubelet_volume_stats_capacity_bytes) \* 100 < 25 | Volume is almost full (< 25% left)<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| KubernetesVolumeOutOfDiskSpaceHigh | Kubernetes Volume out of disk space (instance {{ $labels.instance }}) | 2m | high | (kubelet_volume_stats_available_bytes / kubelet_volume_stats_capacity_bytes) \* 100 < 10 | Volume is almost full (< 10% left)<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| KubernetesVolumeFullInFourDays | Kubernetes Volume full in four days (instance {{ $labels.instance }}) | 10m | warning | predict_linear(kubelet_volume_stats_available_bytes[6h], 345600) < 0 | {{ $labels.namespace }}/{{ $labels.persistentvolumeclaim }} is expected to fill up within four days. Currently {{ $value | humanize }}% is available.<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} |
| KubernetesStatefulsetDown | Kubernetes StatefulSet down (instance {{ $labels.instance }}) | 5m | critical | kube_statefulset_replicas \- kube_statefulset_status_replicas_ready != 0 | A StatefulSet went down<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| KubernetesPodNotHealthy | Kubernetes Pod not healthy (instance {{ $labels.instance }}) | 5m | critical | min_over_time(sum by (exported_namespace, exported_pod) (kube_pod_status_phase{phase=\~"Pending&#124;Unknown&#124;Failed"})[1h:1m]) \> 0 | Pod has been in a non\-ready state for longer than an hour.<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| KubernetesPodCrashLooping | Kubernetes pod crash looping (instance {{ $labels.instance }}) | 5m | warning | (rate(kube_pod_container_status_restarts_total[15m]) \* 60) \* 5 \> 5 | Pod {{ $labels.pod }} is crash looping<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| KubernetesReplicassetMismatch | Kubernetes ReplicasSet mismatch (instance {{ $labels.instance }}) | 5m | warning | kube_replicaset_spec_replicas \- kube_replicaset_status_ready_replicas != 0 | Deployment Replicas mismatch<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| KubernetesDeploymentReplicasMismatch | Kubernetes Deployment replicas mismatch (instance {{ $labels.instance }}) | 5m | warning | kube_deployment_spec_replicas \- kube_deployment_status_replicas_available != 0 | Deployment Replicas mismatch<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| KubernetesStatefulsetReplicasMismatch | Kubernetes StatefulSet replicas mismatch (instance {{ $labels.instance }}) | 5m | warning | kube_statefulset_status_replicas_ready \- kube_statefulset_status_replicas != 0 | A StatefulSet has not matched the expected number of replicas for longer than 15 minutes.<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| KubernetesDeploymentGenerationMismatch | Kubernetes Deployment generation mismatch (instance {{ $labels.instance }}) | 5m | critical | kube_deployment_status_observed_generation \- kube_deployment_metadata_generation != 0 | A Deployment has failed but has not been rolled back.<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| KubernetesStatefulsetGenerationMismatch | Kubernetes StatefulSet generation mismatch (instance {{ $labels.instance }}) | 5m | critical | kube_statefulset_status_observed_generation \- kube_statefulset_metadata_generation != 0 | A StatefulSet has failed but has not been rolled back.<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| KubernetesStatefulsetUpdateNotRolledOut | Kubernetes StatefulSet update not rolled out (instance {{ $labels.instance }}) | 5m | critical | max without (revision) (kube_statefulset_status_current_revision unless kube_statefulset_status_update_revision) \* (kube_statefulset_replicas != kube_statefulset_status_replicas_updated) | StatefulSet update has not been rolled out.<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| KubernetesDaemonsetRolloutStuck | Kubernetes DaemonSet rollout stuck (instance {{ $labels.instance }}) | 5m | critical | (((kube_daemonset_status_number_ready / kube_daemonset_status_desired_number_scheduled) \* 100) < 100) or (kube_daemonset_status_desired_number_scheduled \- kube_daemonset_status_current_number_scheduled \> 0) | Some Pods of DaemonSet are not scheduled or not ready<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| KubernetesDaemonsetMisscheduled | Kubernetes DaemonSet misscheduled (instance {{ $labels.instance }}) | 5m | critical | kube_daemonset_status_number_misscheduled \> 0 | Some DaemonSet Pods are running where they are not supposed to run<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| KubernetesCronjobTooLong | Kubernetes CronJob too long (instance {{ $labels.instance }}) | 5m | warning | time() \- kube_cronjob_next_schedule_time \> 3600 | CronJob {{ $labels.namespace }}/{{ $labels.cronjob }} is taking more than 1h to complete.<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| KubernetesJobCompletion | Kubernetes job completion (instance {{ $labels.instance }}) | 5m | critical | (kube_job_spec_completions \- kube_job_status_succeeded \> 0) or (kube_job_status_failed \> 0) | Kubernetes Job failed to complete<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| KubernetesApiServerErrors | Kubernetes API server errors (instance {{ $labels.instance }}) | 5m | critical | (sum(rate(apiserver_request_count{job="kube\-apiserver",code=\~"(?:5..)$"}[2m])) / sum(rate(apiserver_request_count{job="kube\-apiserver"}[2m]))) \* 100 \> 3 | Kubernetes API server is experiencing high error rate<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| ApiServerRequestsSlow | API Server requests are slow(instance {{ $labels.instance }}) | 5m | warning | histogram_quantile(0.99, rate(apiserver_request_duration_seconds_bucket{verb!="WATCH"}[5m])) \> 0.2 | HTTP requests slowing down, 99th quantile is over 0.2s for 5 minutes\n  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| ControllerWorkQueueDepth | Controller work queue depth is more than 10 (instance {{ $labels.instance }}) | 5m | warning | sum(workqueue_depth) \> 10 | Controller work queue depth is more than 10<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| KubernetesApiClientErrors | Kubernetes API client errors (instance {{ $labels.instance }}) | 5m | critical | (sum(rate(rest_client_requests_total{code=\~"(4&#124;5).."}[2m])) by (instance, job) / sum(rate(rest_client_requests_total[2m])) by (instance, job)) \* 100 \> 5 | Kubernetes API client is experiencing high error rate<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| KubernetesClientCertificateExpiresNextWeek | Kubernetes client certificate expires next week (instance {{ $labels.instance }}) | 5m | warning | (apiserver_client_certificate_expiration_seconds_count{job="kubelet"}) \> 0 and histogram_quantile(0.01, sum by (job, le) (rate(apiserver_client_certificate_expiration_seconds_bucket{job="kubelet"}[5m]))) < 604800 | A client certificate used to authenticate to the apiserver is expiring next week.<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| KubernetesClientCertificateExpiresSoon | Kubernetes client certificate expires soon (instance {{ $labels.instance }}) | 5m | critical | (apiserver_client_certificate_expiration_seconds_count{job="kubelet"}) \> 0 and histogram_quantile(0.01, sum by (job, le) (rate(apiserver_client_certificate_expiration_seconds_bucket{job="kubelet"}[5m]))) < 86400 | A client certificate used to authenticate to the apiserver is expiring in less than 24.0 hours.<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
<!-- markdownlint-enable line-length -->

### NodeProcesses

#### Alerting rules

<!-- markdownlint-disable line-length -->
| Name | Summary | For | Severity | Expression | Description | Other labels | Other annotations |
| ---- | ------- | --- | -------- | ---------- | ----------- | ------------ | ----------------- |
| CountPidsAndThreadOutOfLimit | Host high PIDs and Threads usage (instance {{ $labels.instance }}) | 5m | high | (sum(container_processes) by (node) \+  on (node) label_replace(node_processes_threads \* on(instance) group_left(nodename) (node_uname_info), "node", "$1", "nodename", "(.\+)")) / on (node) label_replace(node_processes_max_processes \* on(instance) group_left(nodename) (node_uname_info), "node", "$1", "nodename", "(.\+)") \* 100 \> 80 | Sum of node's pids and threads is filling up (< 20% left)<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
<!-- markdownlint-enable line-length -->

### NodeExporters

#### Alerting rules

<!-- markdownlint-disable line-length -->
| Name | Summary | For | Severity | Expression | Description | Other labels | Other annotations |
| ---- | ------- | --- | -------- | ---------- | ----------- | ------------ | ----------------- |
| NodeDiskUsageIsMoreThanThreshold | Disk usage on node > 70% (instance {{ $labels.node }}) | 5m | warning | (node_filesystem_size_bytes{fstype=\~"ext.\*&#124;xfs", mountpoint !\~".\*pod.\*"} \- node_filesystem_free_bytes{fstype=\~"ext.\*&#124;xfs", mountpoint !\~".\*pod.\*"}) \* 100 / (node_filesystem_avail_bytes{fstype=\~"ext.\*&#124;xfs", mountpoint !\~".\*pod.\*"} \+ (node_filesystem_size_bytes{fstype=\~"ext.\*&#124;xfs", mountpoint !\~".\*pod.\*"} \- node_filesystem_free_bytes{fstype=\~"ext.\*&#124;xfs", mountpoint !\~".\*pod.\*"})) \> 70 | Node {{ $labels.node }} disk usage of {{ $labels.mountpoint }} is<br/>  VALUE = {{ $value }}% | {} | {} |
| NodeDiskUsageIsMoreThanThreshold | Disk usage on node > 90% (instance {{ $labels.node }}) | 5m | high | (node_filesystem_size_bytes{fstype=\~"ext.\*&#124;xfs", mountpoint !\~".\*pod.\*"} \- node_filesystem_free_bytes{fstype=\~"ext.\*&#124;xfs", mountpoint !\~".\*pod.\*"}) \* 100 / (node_filesystem_avail_bytes{fstype=\~"ext.\*&#124;xfs", mountpoint !\~".\*pod.\*"} \+ (node_filesystem_size_bytes{fstype=\~"ext.\*&#124;xfs", mountpoint !\~".\*pod.\*"} \- node_filesystem_free_bytes{fstype=\~"ext.\*&#124;xfs", mountpoint !\~".\*pod.\*"})) \> 90 | Node {{ $labels.node }} disk usage of {{ $labels.mountpoint }} is<br/> VALUE = {{ $value }}% | {} | {} |
| HostOutOfMemory | Host out of memory (instance {{ $labels.instance }}) | 5m | warning | ((node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes) \* 100) \* on(instance) group_left(nodename) node_uname_info < 10 | Node memory is filling up (< 10% left)<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| HostMemoryUnderMemoryPressure | Host memory under memory pressure (instance {{ $labels.instance }}) | 5m | warning | rate(node_vmstat_pgmajfault[2m]) \* on(instance) group_left(nodename) node_uname_info \> 1000 | The node is under heavy memory pressure. High rate of major page faults<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| HostUnusualNetworkThroughputIn | Host unusual network throughput in (instance {{ $labels.instance }}) | 5m | warning | ((sum by (instance) (irate(node_network_receive_bytes_total[2m])) \* on(instance) group_left(nodename) node_uname_info) / 1024) / 1024 \> 100 | Host network interfaces are probably receiving too much data (\> 100 MB/s)<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| HostUnusualNetworkThroughputOut | Host unusual network throughput out (instance {{ $labels.instance }}) | 5m | warning | ((sum by (instance) (irate(node_network_transmit_bytes_total[2m])) \* on(instance) group_left(nodename) node_uname_info) / 1024) / 1024 \> 100 | Host network interfaces are probably sending too much data (\> 100 MB/s)<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| HostUnusualDiskReadRate | Host unusual disk read rate (instance {{ $labels.instance }}) | 5m | warning | (sum by (instance) (irate(node_disk_read_bytes_total[2m])) \* on(instance) group_left(nodename) node_uname_info) / 1024 / 1024 \> 50 | Disk is probably reading too much data (\> 50 MB/s)<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| HostUnusualDiskWriteRate | Host unusual disk write rate (instance {{ $labels.instance }}) | 5m | warning | ((sum by (instance) (irate(node_disk_written_bytes_total[2m])) \* on(instance) group_left(nodename) node_uname_info) / 1024) / 1024 \> 50 | Disk is probably writing too much data (\> 50 MB/s)<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| HostOutOfDiskSpace | Host out of disk space (instance {{ $labels.instance }}) | 5m | warning | ((node_filesystem_avail_bytes{mountpoint="/"}  \* 100) / node_filesystem_size_bytes{mountpoint="/"}) \* on(instance) group_left(nodename) node_uname_info < 10 | Disk is almost full (< 10% left)<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| HostDiskWillFillIn4Hours | Host disk will fill in 4 hours (instance {{ $labels.instance }}) | 5m | warning | predict_linear(node_filesystem_free_bytes{fstype!\~"tmpfs"}[1h], 14400) \* on(instance) group_left(nodename) node_uname_info < 0 | Disk will fill in 4 hours at current write rate<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| HostOutOfInodes | Host out of inodes (instance {{ $labels.instance }}) | 5m | warning | ((node_filesystem_files_free{mountpoint ="/"} / node_filesystem_files{mountpoint ="/"}) \* 100) \* on(instance) group_left(nodename) node_uname_info < 10 | Disk is almost running out of available inodes (< 10% left)<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| HostUnusualDiskReadLatency | Host unusual disk read latency (instance {{ $labels.instance }}) | 5m | warning | (rate(node_disk_read_time_seconds_total[2m]) / rate(node_disk_reads_completed_total[2m])) \* on(instance) group_left(nodename) node_uname_info \> 100 | Disk latency is growing (read operations \> 100ms)<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| HostUnusualDiskWriteLatency | Host unusual disk write latency (instance {{ $labels.instance }}) | 5m | warning | (rate(node_disk_write_time_seconds_total[2m]) / rate(node_disk_writes_completed_total[2m])) \* on(instance) group_left(nodename) node_uname_info \> 100 | Disk latency is growing (write operations \> 100ms)<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| HostHighCpuLoad | Host high CPU load (instance {{ $labels.instance }}) | 5m | warning | 100 \- ((avg(irate(node_cpu_seconds_total{mode="idle"}[5m])) by (instance) \* 100) \* on (instance) group_left (nodename) node_uname_info) \> 80 | CPU load is \> 80%<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
<!-- markdownlint-enable line-length -->

### DockerContainers

#### Alerting rules

<!-- markdownlint-disable line-length -->
| Name | Summary | For | Severity | Expression | Description | Other labels | Other annotations |
| ---- | ------- | --- | -------- | ---------- | ----------- | ------------ | ----------------- |
| ContainerKilled | Container killed (instance {{ $labels.instance }}) | 5m | warning | time() \- container_last_seen \> 60 | A container has disappeared<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| ContainerVolumeUsage | Container Volume usage (instance {{ $labels.instance }}) | 5m | warning | (1 \- (sum(container_fs_inodes_free) BY (node) / sum(container_fs_inodes_total) BY (node))) \* 100 \> 80 | Container Volume usage is above 80%<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| ContainerVolumeIoUsage | Container Volume IO usage (instance {{ $labels.instance }}) | 5m | warning | (sum(container_fs_io_current) BY (node, name) \* 100) \> 80 | Container Volume IO usage is above 80%<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| ContainerHighThrottleRate | Container high throttle rate (instance {{ $labels.instance }}) | 5m | warning | rate(container_cpu_cfs_throttled_seconds_total[3m]) \> 1 | Container is being throttled<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
<!-- markdownlint-enable line-length -->

### HAmode

#### Alerting rules

<!-- markdownlint-disable line-length -->
| Name | Summary | For | Severity | Expression | Description | Other labels | Other annotations |
| ---- | ------- | --- | -------- | ---------- | ----------- | ------------ | ----------------- |
| NotHAKubernetesDeploymentAvailableReplicas | Not HA mode: Deployment Available Replicas < 2 (instance {{ $labels.instance }}) | 5m | warning | `kube_deployment_status_replicas_available < 2` | Not HA mode: Kubernetes Deployment has less than 2 available replicas<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| NotHAKubernetesStatefulSetAvailableReplicas | Not HA mode: StatefulSet Available Replicas < 2 (instance {{ $labels.instance }}) | 5m | warning | `kube_statefulset_status_replicas_available < 2` | Not HA mode: Kubernetes StatefulSet has less than 2 available replicas<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| NotHAKubernetesDeploymentDesiredReplicas | Not HA mode: Deployment Desired Replicas < 2 (instance {{ $labels.instance }}) | 5m | warning | `kube_deployment_status_replicas < 2` | Not HA mode: Kubernetes Deployment has less than 2 desired replicas<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| NotHAKubernetesStatefulSetDesiredReplicas | Not HA mode: StatefulSet Desired Replicas < 2 (instance {{ $labels.instance }}) | 5m | warning | `kube_statefulset_status_replicas < 2` | Not HA mode: Kubernetes StatefulSet has less than 2 desired replicas<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| NotHAKubernetesDeploymentMultiplePodsPerNode | Not HA mode: Deployment Has Multiple Pods per Node (instance {{ $labels.instance }}) | 5m | warning | `count(sum(kube_pod_info{node=\~".\+", created_by_kind="ReplicaSet"}) by (namespace, node, created_by_name) \> 1) \> 0` | Not HA mode: Kubernetes Deployment has 2 or more replicas on the same node<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| NotHAKubernetesStatefulSetMultiplePodsPerNode | Not HA mode: StatefulSet Has Multiple Pods per Node (instance {{ $labels.instance }}) | 5m | warning | `count(sum(kube_pod_info{node=\~".\+", created_by_kind="StatefulSet"}) by (namespace, node, created_by_name) \> 1) \> 0` | Not HA mode: Kubernetes StatefulSet has 2 or more replicas on the same node<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
<!-- markdownlint-enable line-length -->

### HAproxy

#### Alerting rules

<!-- markdownlint-disable line-length -->
| Name | Summary | For | Severity | Expression | Description | Other labels | Other annotations |
| ---- | ------- | --- | -------- | ---------- | ----------- | ------------ | ----------------- |
| HaproxyDown | HAProxy down (instance {{ $labels.instance }}) | 5m | critical | haproxy_up == 0 | HAProxy down<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| HaproxyBackendConnectionErrors | HAProxy backend connection errors (instance {{ $labels.instance }}) | 5m | critical | sum by (backend) (rate(haproxy_backend_connection_errors_total[2m])) \> 10 | Too many connection errors to {{ $labels.fqdn }}/{{ $labels.backend }} backend (\> 10 req/s). Request throughput may be to high.<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| HaproxyServerResponseErrors | HAProxy server response errors (instance {{ $labels.instance }}) | 5m | critical | sum by (server) (rate(haproxy_server_response_errors_total[2m])) \> 5 | Too many response errors to {{ $labels.server }} server (\> 5 req/s).<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| HaproxyServerConnectionErrors | HAProxy server connection errors (instance {{ $labels.instance }}) | 5m | critical | sum by (server) (rate(haproxy_server_connection_errors_total[2m])) \> 10 | Too many connection errors to {{ $labels.server }} server (\> 10 req/s). Request throughput may be to high.<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| HaproxyPendingRequests | HAProxy pending requests (instance {{ $labels.instance }}) | 5m | warning | sum by (backend) (haproxy_backend_current_queue) \> 0 | Some HAProxy requests are pending on {{ $labels.fqdn }}/{{ $labels.backend }} backend<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| HaproxyHttpSlowingDown | HAProxy HTTP slowing down (instance {{ $labels.instance }}) | 5m | warning | avg by (backend) (haproxy_backend_http_total_time_average_seconds) \> 2 | Average request time is increasing<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| HaproxyRetryHigh | HAProxy retry high (instance {{ $labels.instance }}) | 5m | warning | sum by (backend) (rate(haproxy_backend_retry_warnings_total[5m])) \> 10 | High rate of retry on {{ $labels.fqdn }}/{{ $labels.backend }} backend<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| HaproxyBackendDown | HAProxy backend down (instance {{ $labels.instance }}) | 5m | critical | haproxy_backend_up == 0 | HAProxy backend is down<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| HaproxyServerDown | HAProxy server down (instance {{ $labels.instance }}) | 5m | critical | haproxy_server_up == 0 | HAProxy server is down<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| HaproxyFrontendSecurityBlockedRequests | HAProxy frontend security blocked requests (instance {{ $labels.instance }}) | 5m | warning | sum by (frontend) (rate(haproxy_frontend_requests_denied_total[5m])) \> 10 | HAProxy is blocking requests for security reason<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| HaproxyServerHealthcheckFailure | HAProxy server healthcheck failure (instance {{ $labels.instance }}) | 5m | warning | increase(haproxy_server_check_failures_total[5m]) \> 0 | Some server healthcheck are failing on {{ $labels.server }}<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
<!-- markdownlint-enable line-length -->

### Etcd

#### Alerting rules

<!-- markdownlint-disable line-length -->
| Name | Summary | For | Severity | Expression | Description | Other labels | Other annotations |
| ---- | ------- | --- | -------- | ---------- | ----------- | ------------ | ----------------- |
| EtcdInsufficientMembers | Etcd insufficient Members (instance {{ $labels.instance }}) | 5m | critical | count(etcd_server_id{job="etcd"}) % 2 == 0 | Etcd cluster should have an odd number of members<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| EtcdNoLeader | Etcd no Leader (instance {{ $labels.instance }}) | 5m | critical | etcd_server_has_leader == 0 | Etcd cluster have no leader<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| EtcdHighNumberOfLeaderChanges | Etcd high number of leader changes (instance {{ $labels.instance }}) | 5m | warning | increase(etcd_server_leader_changes_seen_total[1h]) \> 3 | Etcd leader changed more than 3 times during last hour<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| EtcdHighNumberOfFailedGrpcRequests | Etcd high number of failed GRPC requests (instance {{ $labels.instance }}) | 5m | warning | sum(rate(grpc_server_handled_total{job="etcd",grpc_code!="OK"}[5m])) BY (grpc_service, grpc_method) / sum(rate(grpc_server_handled_total{job="etcd"}[5m])) BY (grpc_service, grpc_method) \> 0.01 | More than 1% GRPC request failure detected in Etcd for 5 minutes<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| EtcdHighNumberOfFailedGrpcRequests | Etcd high number of failed GRPC requests (instance {{ $labels.instance }}) | 5m | critical | sum(rate(grpc_server_handled_total{job="etcd",grpc_code!="OK"}[5m])) BY (grpc_service, grpc_method) / sum(rate(grpc_server_handled_total{job="etcd"}[5m])) BY (grpc_service, grpc_method) \> 0.05 | More than 5% GRPC request failure detected in Etcd for 5 minutes<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| EtcdGrpcRequestsSlow | Etcd GRPC requests slow (instance {{ $labels.instance }}) | 5m | warning | histogram_quantile(0.99, sum(rate(grpc_server_handling_seconds_bucket{job="etcd",grpc_type="unary"}[5m])) by (grpc_service, grpc_method, le)) \> 0.15 | GRPC requests slowing down, 99th percentil is over 0.15s for 5 minutes<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| EtcdMemberCommunicationSlow | Etcd member communication slow (instance {{ $labels.instance }}) | 5m | warning | histogram_quantile(0.99, rate(etcd_network_peer_round_trip_time_seconds_bucket{job="etcd"}[5m])) \> 0.15 | Etcd member communication slowing down, 99th percentil is over 0.15s for 5 minutes<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| EtcdHighNumberOfFailedProposals | Etcd high number of failed proposals (instance {{ $labels.instance }}) | 5m | warning | increase(etcd_server_proposals_failed_total[1h]) \> 5 | Etcd server got more than 5 failed proposals past hour<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| EtcdHighFsyncDurations | Etcd high fsync durations (instance {{ $labels.instance }}) | 5m | warning | histogram_quantile(0.99, rate(etcd_disk_wal_fsync_duration_seconds_bucket[5m])) \> 0.5 | Etcd WAL fsync duration increasing, 99th percentil is over 0.5s for 5 minutes<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
| EtcdHighCommitDurations | Etcd high commit durations (instance {{ $labels.instance }}) | 5m | warning | histogram_quantile(0.99, rate(etcd_disk_backend_commit_duration_seconds_bucket[5m])) \> 0.25 | Etcd commit duration increasing, 99th percentil is over 0.25s for 5 minutes<br/>  VALUE = {{ $value }}<br/>  LABELS: {{ $labels }} | {} | {} |
<!-- markdownlint-enable line-length -->

### NginxIngressAlerts

#### Alerting rules

<!-- markdownlint-disable line-length -->
| Name | Summary | For | Severity | Expression | Description | Other labels | Other annotations |
| ---- | ------- | --- | -------- | ---------- | ----------- | ------------ | ----------------- |
| NginxHighHttp4xxErrorRate | Nginx high HTTP 4xx error rate (node: {{ $labels.node }}, namespace: {{ $labels.exported_namespace }}, ingress: {{ $labels.ingress }}) | 1m | high | sum by (ingress, exported_namespace, node) (rate(nginx_ingress_controller_requests{status=\~"^4.."}[2m])) / sum by (ingress, exported_namespace, node)(rate(nginx_ingress_controller_requests[2m])) \* 100 \> 5 | Too many HTTP requests with status 4xx (\> 5%)<br/>  VALUE = {{ $value }}<br/>  LABELS = {{ $labels }} | {} | {} |
| NginxHighHttp5xxErrorRate | Nginx high HTTP 5xx error rate (node: {{ $labels.node }}, namespace: {{ $labels.exported_namespace }}, ingress: {{ $labels.ingress }}) | 1m | high | sum by (ingress, exported_namespace, node) (rate(nginx_ingress_controller_requests{status=\~"^5.."}[2m])) / sum by (ingress, exported_namespace, node) (rate(nginx_ingress_controller_requests[2m])) \* 100 \> 5 | Too many HTTP requests with status 5xx (\> 5%)<br/>  VALUE = {{ $value }}<br/>  LABELS = {{ $labels }} | {} | {} |
| NginxLatencyHigh | Nginx latency high (node: {{ $labels.node }}, host: {{ $labels.host }}) | 2m | warning | histogram_quantile(0.99, sum(rate(nginx_ingress_controller_request_duration_seconds_bucket[2m])) by (host, node, le)) \> 3 | Nginx p99 latency is higher than 3 seconds<br/>  VALUE = {{ $value }}<br/>  LABELS = {{ $labels }} | {} | {} |
<!-- markdownlint-enable line-length -->

### CoreDnsAlerts

#### Alerting rules
<!-- markdownlint-disable line-length -->
| Name | Summary | For | Severity | Expression | Description | Other labels | Other annotations |
| ---- | ------- | --- | -------- | ---------- | ----------- | ------------ | ----------------- |
| CorednsPanicCount | CoreDNS Panic Count (instance {{ $labels.instance }}) | 0m | critical | increase(coredns_panics_total[1m]) \> 0 | Number of CoreDNS panics encountered<br/>  VALUE = {{ $value }}<br/>  LABELS = {{ $labels }} | {} | {} |
| CoreDNSLatencyHigh | CoreDNS have High Latency | 5m | critical | histogram_quantile(0.99, sum(rate(coredns_dns_request_duration_seconds_bucket[2m])) by(server, zone, le)) \> 3 | CoreDNS has 99th percentile latency of {{ $value }} seconds for server {{ $labels.server }} zone {{ $labels.zone }} | {} | {} |
| CoreDNSForwardHealthcheckFailureCount | CoreDNS health checks have failed to upstream server | 5m | warning | sum(rate(coredns_forward_healthcheck_broken_total[2m])) \> 0 | CoreDNS health checks have failed to upstream server {{ $labels.to }} | {} | {} |
| CoreDNSForwardHealthcheckBrokenCount | CoreDNS health checks have failed for all upstream servers | 5m | warning | sum(rate(coredns_forward_healthcheck_broken_total[2m])) \> 0 | CoreDNS health checks failed for all upstream servers LABELS = {{ $labels }} | {} | {} |
| CoreDNSErrorsHigh | CoreDNS is returning SERVFAIL | 5m | critical | sum(rate(coredns_dns_responses_total{rcode="SERVFAIL"}[2m])) / sum(rate(coredns_dns_responses_total[2m])) \> 0.03 | CoreDNS is returning SERVFAIL for {{ $value | humanizePercentage }} of requests | {} |
| CoreDNSErrorsHigh | CoreDNS is returning SERVFAIL | 5m | warning | sum(rate(coredns_dns_responses_total{rcode="SERVFAIL"}[2m])) / sum(rate(coredns_dns_responses_total[2m])) \> 0.01 | CoreDNS is returning SERVFAIL for {{ $value | humanizePercentage }} of requests | {} |
| CoreDNSForwardLatencyHigh | CoreDNS has 99th percentile latency for forwarding requests | 5m | critical | histogram_quantile(0.99, sum(rate(coredns_forward_request_duration_seconds_bucket[2m])) by(to, le)) \> 3 | CoreDNS has 99th percentile latency of {{ $value }} seconds forwarding requests to {{ $labels.to }} | {} | {} |
| CoreDNSForwardErrorsHigh | CoreDNS is returning SERVFAIL for forward requests | 5m | critical | sum(rate(coredns_forward_responses_total{rcode="SERVFAIL"}[2m])) / sum(rate(coredns_forward_responses_total[2m])) \> 0.03 | CoreDNS is returning SERVFAIL for {{ $value | humanizePercentage }} of forward requests to {{ $labels.to }} | {} |
| CoreDNSForwardErrorsHigh | CoreDNS is returning SERVFAIL for forward requests | 5m | warning | sum(rate(coredns_forward_responses_total{rcode="SERVFAIL"}[2m])) / sum(rate(coredns_forward_responses_total[2m])) \> 0.01 | CoreDNS is returning SERVFAIL for {{ $value | humanizePercentage }} of forward requests to {{ $labels.to }} | {} |
<!-- markdownlint-enable line-length -->

### DRAlerts

#### Alerting rules

<!-- markdownlint-disable line-length -->
| Name | Summary | For | Severity | Expression | Description | Other labels | Other annotations |
| ---- | ------- | --- | -------- | ---------- | ----------- | ------------ | ----------------- |
| ProbeFailed | Probe failed (instance: {{ $labels.instance }}) | 5m | critical | probe_success == 0 | Probe failed\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }} | {} | {} |
| SlowProbe | Slow probe (instance: {{ $labels.instance }}) | 5m | warning | avg_over_time(probe_duration_seconds[1m]) > 1 | Blackbox probe took more than 1s to complete\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }} | {} | {} |
| HttpStatusCode | HTTP Status Code (instance: {{ $labels.instance }}) | 5m | high | probe_http_status_code <= 199 OR probe_http_status_code >= 400 | HTTP status code is not 200-399\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }} | {} | {} |
| HttpSlowRequests | HTTP slow requests (instance: {{ $labels.instance }}) | 5m | warning | avg_over_time(probe_http_duration_seconds[1m]) > 1 | HTTP request took more than 1s\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }} | {} | {} |
<!-- markdownlint-enable line-length -->

### BackupAlerts

#### Alerting rules

<!-- markdownlint-disable line-length -->
| Name | Summary | For | Severity | Expression | Description | Other labels | Other annotations |
| ---- | ------- | --- | -------- | ---------- | ----------- | ------------ | ----------------- |
| Last Backup Failed | Last backup made by pod {{ $labels.pod }} in namespace {{ $labels.namespace }} failed. | 1m | warning | backup_storage_last_failed != 0 | Last backup made by pod {{ $labels.pod }} in namespace {{ $labels.namespace }} failed.\n  VALUE = {{ $value }}\n  LABELS: {{ $labels }} | {} | {} |
<!-- markdownlint-enable line-length -->

## Cert-exporter

### Cert-exporter rules

#### Alerting rules

<!-- markdownlint-disable line-length -->
| Name | Summary | For | Severity | Expression | Description | Other labels | Other annotations |
| ---- | ------- | --- | -------- | ---------- | ----------- | ------------ | ----------------- |
| FileCerts30DaysRemaining | Certificates from files expire within 30 days | 10m | warning | count(86400 \* 7 < cert_exporter_cert_expires_in_seconds <= 86400 \* 30) \> 0 | Some certificates from files expire within 30 days. | {} | {} |
| FileCerts7DaysRemaining | Certificates from files expire within 7 days | 10m | high | count(0 < cert_exporter_cert_expires_in_seconds <= 86400 \* 7) \> 0 | Some certificates from files expire within 7 days. | {} | {} |
| FileCertsExpired | Certificates from files expired | 10m | critical | count(cert_exporter_cert_expires_in_seconds <= 0) \> 0 | Some certificates from files already expired. | {} | {} |
| KubeconfigCerts30DaysRemaining | Certificates from kubeconfig expire within 30 days | 10m | warning | count(86400 \* 7 < cert_exporter_kubeconfig_expires_in_seconds <= 86400 \* 30) \> 0 | Some certificates from kubeconfig expire within 30 days. | {} | {} |
| KubeconfigCerts7DaysRemaining | Certificates from kubeconfig expire within 7 days | 10m | high | count(0 < cert_exporter_kubeconfig_expires_in_seconds <= 86400 \* 7) \> 0 | Some certificates from kubeconfig expire within 7 days. | {} | {} |
| KubeconfigCertsExpired | Certificates from kubeconfig expired | 10m | critical | count(cert_exporter_kubeconfig_expires_in_seconds <= 0) \> 0 | Some certificates from kubeconfig already expired. | {} | {} |
| SecretCerts30DaysRemaining | Certificates from secrets expire within 30 days | 10m | warning | count(86400 \* 7 < cert_exporter_secret_expires_in_seconds <= 86400 \* 30) \> 0 | Some certificates from secrets expire within 30 days. | {} | {} |
| SecretCerts7DaysRemaining | Certificates from secrets expire within 7 days | 10m | high | count(0 < cert_exporter_secret_expires_in_seconds <= 86400 \* 7) \> 0 | Some certificates from secrets expire within 7 days. | {} | {} |
| SecretCertsExpired | Certificates from secrets expired | 10m | critical | count(cert_exporter_secret_expires_in_seconds <= 0) \> 0 | Some certificates from secrets already expired. | {} | {} |
<!-- markdownlint-enable line-length -->
