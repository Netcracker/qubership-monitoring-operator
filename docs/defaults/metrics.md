# Metrics OOB

* [Metrics OOB](#metrics-oob)
  * [AlertManager (Kubernetes v1.21)](#alertmanager-kubernetes-v121)
    * [Statistics](#statistics)
    * [Metrics](#metrics)
  * [Core DNS service (Kubernetes v1.21)](#core-dns-service-kubernetes-v121)
    * [Statistics](#statistics-1)
    * [Metrics](#metrics-1)
  * [Grafana (Kubernetes v1.21)](#grafana-kubernetes-v121)
    * [Statistics](#statistics-2)
    * [Metrics](#metrics-2)
  * [Grafana-operator (Kubernetes v1.21)](#grafana-operator-kubernetes-v121)
    * [Statistics](#statistics-3)
    * [Metrics](#metrics-3)
  * [Kube-apiserver (Kubernetes v1.21)](#kube-apiserver-kubernetes-v121)
    * [Statistics](#statistics-4)
    * [Metrics](#metrics-4)
  * [Kube-state-metrics (Kubernetes v1.21)](#kube-state-metrics-kubernetes-v121)
    * [Statistics](#statistics-5)
    * [Metrics](#metrics-5)
  * [Kubelet (Kubernetes v1.21)](#kubelet-kubernetes-v121)
    * [Statistics](#statistics-6)
    * [Metrics](#metrics-6)
  * [cAdvisor (Kubernetes v1.21)](#cadvisor-kubernetes-v121)
    * [Statistics](#statistics-7)
    * [Metrics](#metrics-7)
  * [NGINX ingress (Kubernetes v1.21)](#nginx-ingress-kubernetes-v121)
    * [Statistics](#statistics-8)
    * [Metrics](#metrics-8)
  * [Node-exporter (Kubernetes v1.21)](#node-exporter-kubernetes-v121)
    * [Statistics](#statistics-9)
    * [Metrics](#metrics-9)
  * [Prometheus (Kubernetes v1.21)](#prometheus-kubernetes-v121)
    * [Statistics](#statistics-10)
    * [Metrics](#metrics-10)
  * [Prometheus-operator (Kubernetes v1.21)](#prometheus-operator-kubernetes-v121)
    * [Statistics](#statistics-11)
    * [Metrics](#metrics-11)
  * [AlertManager (Kubernetes v1.23)](#alertmanager-kubernetes-v123)
    * [Statistics](#statistics-12)
    * [Metrics](#metrics-12)
  * [Core DNS service (Kubernetes v1.23)](#core-dns-service-kubernetes-v123)
    * [Statistics](#statistics-13)
    * [Metrics](#metrics-13)
  * [Grafana (Kubernetes v1.23)](#grafana-kubernetes-v123)
    * [Statistics](#statistics-14)
    * [Metrics](#metrics-14)
  * [Grafana-operator (Kubernetes v1.23)](#grafana-operator-kubernetes-v123)
    * [Statistics](#statistics-15)
    * [Metrics](#metrics-15)
  * [Kube-apiserver (Kubernetes v1.23)](#kube-apiserver-kubernetes-v123)
    * [Statistics](#statistics-16)
    * [Metrics](#metrics-16)
  * [Kube-state-metrics (Kubernetes v1.23)](#kube-state-metrics-kubernetes-v123)
    * [Statistics](#statistics-17)
    * [Metrics](#metrics-17)
  * [Kubelet (Kubernetes v1.23)](#kubelet-kubernetes-v123)
    * [Statistics](#statistics-18)
    * [Metrics](#metrics-18)
  * [cAdvisor (Kubernetes v1.23)](#cadvisor-kubernetes-v123)
    * [Statistics](#statistics-19)
    * [Metrics](#metrics-19)
  * [NGINX ingress (Kubernetes v1.23)](#nginx-ingress-kubernetes-v123)
    * [Statistics](#statistics-20)
    * [Metrics](#metrics-20)
  * [Node-exporter (Kubernetes v1.23)](#node-exporter-kubernetes-v123)
    * [Statistics](#statistics-21)
    * [Metrics](#metrics-21)
  * [Prometheus (Kubernetes v1.23)](#prometheus-kubernetes-v123)
    * [Statistics](#statistics-22)
    * [Metrics](#metrics-22)
  * [Prometheus-operator (Kubernetes v1.23)](#prometheus-operator-kubernetes-v123)
    * [Statistics](#statistics-23)
    * [Metrics](#metrics-23)
  * [Operator-sdk](#operator-sdk)
    * [Metrics](#metrics-24)
  * [Openshift-state-metrics](#openshift-state-metrics)
    * [Statistics](#statistics-24)
    * [Metrics](#metrics-25)
  * [Openshift-api-server](#openshift-api-server)
    * [Statistics](#statistics-25)
    * [Metrics](#metrics-26)
  * [Openshift-api-server-operator](#openshift-api-server-operator)
    * [Statistics](#statistics-26)
    * [Metrics](#metrics-27)
  * [Cluster-version-operator](#cluster-version-operator)
    * [Statistics](#statistics-27)
    * [Metrics](#metrics-28)
  * [Openshift-HAProxy](#openshift-haproxy)
    * [Statistics](#statistics-28)
    * [Metrics](#metrics-29)

## AlertManager (Kubernetes v1.21)

### Statistics

<!-- markdownlint-disable line-length -->
| Number of metrics families | Number of time series | Average number of labels per time series | Number of unique label pairs |
| -------------------------- | --------------------- | ---------------------------------------- | ---------------------------- |
| 73 | 312 | 1.53 | 60 |
<!-- markdownlint-enable line-length -->

### Metrics

<!-- markdownlint-disable line-length -->
| Name | Type, Unit | Description |
| ---- | ---------- | ----------- |
| alertmanager_alerts | gauge | How many alerts by state. |
| alertmanager_alerts_invalid | counter | The total number of received alerts that were invalid. |
| alertmanager_alerts_received | counter | The total number of received alerts. |
| alertmanager_build_info | gauge | A metric with a constant '1' value labeled by version, revision, branch, and goversion from which alertmanager was built. |
| alertmanager_cluster_enabled | gauge | Indicates whether the clustering is enabled or not. |
| alertmanager_config_hash | gauge | Hash of the currently loaded alertmanager configuration. |
| alertmanager_config_last_reload_success_timestamp_seconds | gauge | Timestamp of the last successful configuration reload. |
| alertmanager_config_last_reload_successful | gauge | Whether the last configuration reload attempt was successful. |
| alertmanager_dispatcher_aggregation_groups | gauge | Number of active aggregation groups |
| alertmanager_dispatcher_alert_processing_duration_seconds | summary | Summary of latencies for the processing of alerts. |
| alertmanager_http_concurrency_limit_exceeded | counter | Total number of times an HTTP request failed because the concurrency limit was reached. |
| alertmanager_http_request_duration_seconds | histogram | Histogram of latencies for HTTP requests. |
| alertmanager_http_requests_in_flight | gauge | Current number of HTTP requests being processed. |
| alertmanager_http_response_size_bytes | histogram | Histogram of response size for HTTP requests. |
| alertmanager_integrations | gauge | Number of configured integrations. |
| alertmanager_nflog_gc_duration_seconds | summary | Duration of the last notification log garbage collection cycle. |
| alertmanager_nflog_gossip_messages_propagated | counter | Number of received gossip messages that have been further gossiped. |
| alertmanager_nflog_queries | counter | Number of notification log queries were received. |
| alertmanager_nflog_query_duration_seconds | histogram | Duration of notification log query evaluation. |
| alertmanager_nflog_query_errors | counter | Number notification log received queries that failed. |
| alertmanager_nflog_snapshot_duration_seconds | summary | Duration of the last notification log snapshot. |
| alertmanager_nflog_snapshot_size_bytes | gauge | Size of the last notification log snapshot in bytes. |
| alertmanager_notification_latency_seconds | histogram | The latency of notifications in seconds. |
| alertmanager_notification_requests_failed | counter | The total number of failed notification requests. |
| alertmanager_notification_requests | counter | The total number of attempted notification requests. |
| alertmanager_notifications_failed | counter | The total number of failed notifications. |
| alertmanager_notifications | counter | The total number of attempted notifications. |
| alertmanager_receivers | gauge | Number of configured receivers. |
| alertmanager_silences | gauge | How many silences by state. |
| alertmanager_silences_gc_duration_seconds | summary | Duration of the last silence garbage collection cycle. |
| alertmanager_silences_gossip_messages_propagated | counter | Number of received gossip messages that have been further gossiped. |
| alertmanager_silences_queries | counter | How many silence queries were received. |
| alertmanager_silences_query_duration_seconds | histogram | Duration of silence query evaluation. |
| alertmanager_silences_query_errors | counter | How many silence received queries did not succeed. |
| alertmanager_silences_snapshot_duration_seconds | summary | Duration of the last silence snapshot. |
| alertmanager_silences_snapshot_size_bytes | gauge | Size of the last silence snapshot in bytes. |
| go_gc_duration_seconds | summary | A summary of the pause duration of garbage collection cycles. |
| go_goroutines | gauge | Number of goroutines that currently exist. |
| go_info | gauge | Information about the Go environment. |
| go_memstats_alloc_bytes | gauge | Number of bytes allocated and still in use. |
| go_memstats_alloc_bytes | counter | Total number of bytes allocated, even if freed. |
| go_memstats_buck_hash_sys_bytes | gauge | Number of bytes used by the profiling bucket hash table. |
| go_memstats_frees | counter | Total number of frees. |
| go_memstats_gc_cpu_fraction | gauge | The fraction of this program's available CPU time used by the GC since the program started. |
| go_memstats_gc_sys_bytes | gauge | Number of bytes used for garbage collection system metadata. |
| go_memstats_heap_alloc_bytes | gauge | Number of heap bytes allocated and still in use. |
| go_memstats_heap_idle_bytes | gauge | Number of heap bytes waiting to be used. |
| go_memstats_heap_inuse_bytes | gauge | Number of heap bytes that are in use. |
| go_memstats_heap_objects | gauge | Number of allocated objects. |
| go_memstats_heap_released_bytes | gauge | Number of heap bytes released to OS. |
| go_memstats_heap_sys_bytes | gauge | Number of heap bytes obtained from system. |
| go_memstats_last_gc_time_seconds | gauge | Number of seconds since 1970 of last garbage collection. |
| go_memstats_lookups | counter | Total number of pointer lookups. |
| go_memstats_mallocs | counter | Total number of mallocs. |
| go_memstats_mcache_inuse_bytes | gauge | Number of bytes in use by mcache structures. |
| go_memstats_mcache_sys_bytes | gauge | Number of bytes used for mcache structures obtained from system. |
| go_memstats_mspan_inuse_bytes | gauge | Number of bytes in use by mspan structures. |
| go_memstats_mspan_sys_bytes | gauge | Number of bytes used for mspan structures obtained from system. |
| go_memstats_next_gc_bytes | gauge | Number of heap bytes when next garbage collection will take place. |
| go_memstats_other_sys_bytes | gauge | Number of bytes used for other system allocations. |
| go_memstats_stack_inuse_bytes | gauge | Number of bytes in use by the stack allocator. |
| go_memstats_stack_sys_bytes | gauge | Number of bytes obtained from system for stack allocator. |
| go_memstats_sys_bytes | gauge | Number of bytes obtained from system. |
| go_threads | gauge | Number of OS threads created. |
| process_cpu_seconds | counter | Total user and system CPU time spent in seconds. |
| process_max_fds | gauge | Maximum number of open file descriptors. |
| process_open_fds | gauge | Number of open file descriptors. |
| process_resident_memory_bytes | gauge | Resident memory size in bytes. |
| process_start_time_seconds | gauge | Start time of the process since unix epoch in seconds. |
| process_virtual_memory_bytes | gauge | Virtual memory size in bytes. |
| process_virtual_memory_max_bytes | gauge | Maximum amount of virtual memory available in bytes. |
| promhttp_metric_handler_requests_in_flight | gauge | Current number of scrapes being served. |
| promhttp_metric_handler_requests | counter | Total number of scrapes by HTTP status code. |
<!-- markdownlint-enable line-length -->

## Core DNS service (Kubernetes v1.21)

### Statistics

<!-- markdownlint-disable line-length -->
| Number of metrics families | Number of time series | Average number of labels per time series | Number of unique label pairs |
| -------------------------- | --------------------- | ---------------------------------------- | ---------------------------- |
| 59 | 366 | 2.93 | 71 |
<!-- markdownlint-enable line-length -->

### Metrics

<!-- markdownlint-disable line-length -->
| Name | Type, Unit | Description |
| ---- | ---------- | ----------- |
| coredns_build_info | gauge | A metric with a constant '1' value labeled by version, revision, and goversion from which CoreDNS was built. |
| coredns_cache_entries | gauge | The number of elements in the cache. |
| coredns_cache_hits | counter | The count of cache hits. |
| coredns_cache_misses | counter | The count of cache misses. |
| coredns_dns_request_duration_seconds | histogram | Histogram of the time (in seconds) each request took. |
| coredns_dns_request_size_bytes | histogram | Size of the EDNS0 UDP buffer in bytes (64K for TCP). |
| coredns_dns_requests | counter | Counter of DNS requests made per zone, protocol and family. |
| coredns_dns_response_size_bytes | histogram | Size of the returned response in bytes. |
| coredns_dns_responses | counter | Counter of response status codes. |
| coredns_forward_conn_cache_hits | counter | Counter of connection cache hits per upstream and protocol. |
| coredns_forward_conn_cache_misses | counter | Counter of connection cache misses per upstream and protocol. |
| coredns_forward_healthcheck_broken | counter | Counter of the number of complete failures of the healthchecks. |
| coredns_forward_healthcheck_failures | counter | Counter of the number of failed healthchecks. |
| coredns_forward_max_concurrent_rejects | counter | Counter of the number of queries rejected because the concurrent queries were at maximum. |
| coredns_forward_request_duration_seconds | histogram | Histogram of the time each request took. |
| coredns_forward_requests | counter | Counter of requests made per upstream. |
| coredns_forward_responses | counter | Counter of responses received per upstream. |
| coredns_health_request_duration_seconds | histogram | Histogram of the time (in seconds) each request took. |
| coredns_hosts_entries | gauge | The combined number of entries in hosts and Corefile. |
| coredns_hosts_reload_timestamp_seconds | gauge | The timestamp of the last reload of hosts file. |
| coredns_panics | counter | A metrics that counts the number of panics. |
| coredns_plugin_enabled | gauge | A metric that indicates whether a plugin is enabled on per server and zone basis. |
| coredns_reload_failed | counter | Counter of the number of failed reload attempts. |
| coredns_template_matches | counter | Counter of template regex matches. |
| go_gc_duration_seconds | summary | A summary of the pause duration of garbage collection cycles. |
| go_goroutines | gauge | Number of goroutines that currently exist. |
| go_info | gauge | Information about the Go environment. |
| go_memstats_alloc_bytes | gauge | Number of bytes allocated and still in use. |
| go_memstats_alloc_bytes | counter | Total number of bytes allocated, even if freed. |
| go_memstats_buck_hash_sys_bytes | gauge | Number of bytes used by the profiling bucket hash table. |
| go_memstats_frees | counter | Total number of frees. |
| go_memstats_gc_cpu_fraction | gauge | The fraction of this program's available CPU time used by the GC since the program started. |
| go_memstats_gc_sys_bytes | gauge | Number of bytes used for garbage collection system metadata. |
| go_memstats_heap_alloc_bytes | gauge | Number of heap bytes allocated and still in use. |
| go_memstats_heap_idle_bytes | gauge | Number of heap bytes waiting to be used. |
| go_memstats_heap_inuse_bytes | gauge | Number of heap bytes that are in use. |
| go_memstats_heap_objects | gauge | Number of allocated objects. |
| go_memstats_heap_released_bytes | gauge | Number of heap bytes released to OS. |
| go_memstats_heap_sys_bytes | gauge | Number of heap bytes obtained from system. |
| go_memstats_last_gc_time_seconds | gauge | Number of seconds since 1970 of last garbage collection. |
| go_memstats_lookups | counter | Total number of pointer lookups. |
| go_memstats_mallocs | counter | Total number of mallocs. |
| go_memstats_mcache_inuse_bytes | gauge | Number of bytes in use by mcache structures. |
| go_memstats_mcache_sys_bytes | gauge | Number of bytes used for mcache structures obtained from system. |
| go_memstats_mspan_inuse_bytes | gauge | Number of bytes in use by mspan structures. |
| go_memstats_mspan_sys_bytes | gauge | Number of bytes used for mspan structures obtained from system. |
| go_memstats_next_gc_bytes | gauge | Number of heap bytes when next garbage collection will take place. |
| go_memstats_other_sys_bytes | gauge | Number of bytes used for other system allocations. |
| go_memstats_stack_inuse_bytes | gauge | Number of bytes in use by the stack allocator. |
| go_memstats_stack_sys_bytes | gauge | Number of bytes obtained from system for stack allocator. |
| go_memstats_sys_bytes | gauge | Number of bytes obtained from system. |
| go_threads | gauge | Number of OS threads created. |
| process_cpu_seconds | counter | Total user and system CPU time spent in seconds. |
| process_max_fds | gauge | Maximum number of open file descriptors. |
| process_open_fds | gauge | Number of open file descriptors. |
| process_resident_memory_bytes | gauge | Resident memory size in bytes. |
| process_start_time_seconds | gauge | Start time of the process since unix epoch in seconds. |
| process_virtual_memory_bytes | gauge | Virtual memory size in bytes. |
| process_virtual_memory_max_bytes | gauge | Maximum amount of virtual memory available in bytes. |
<!-- markdownlint-enable line-length -->

## Grafana (Kubernetes v1.21)

### Statistics

<!-- markdownlint-disable line-length -->
| Number of metrics families | Number of time series | Average number of labels per time series | Number of unique label pairs |
| -------------------------- | --------------------- | ---------------------------------------- | ---------------------------- |
| 158 | 1110 | 2.62 | 269 |
<!-- markdownlint-enable line-length -->

### Metrics

<!-- markdownlint-disable line-length -->
| Name | Type, Unit | Description |
| ---- | ---------- | ----------- |
| access_evaluation_duration | histogram | Histogram for the runtime of evaluation function. |
| access_permissions_duration | histogram | Histogram for the runtime of permissions check function. |
| cortex_experimental_features_in_use | counter | The number of experimental features in use. |
| deprecated_flags_inuse | counter | The number of deprecated flags currently set. |
| go_gc_cycles_automatic_gc_cycles | counter | Count of completed GC cycles generated by the Go runtime. |
| go_gc_cycles_forced_gc_cycles | counter | Count of completed GC cycles forced by the application. |
| go_gc_cycles_total_gc_cycles | counter | Count of all completed GC cycles. |
| go_gc_duration_seconds | summary | A summary of the pause duration of garbage collection cycles. |
| go_gc_heap_allocs_by_size_bytes_total | histogram | Distribution of heap allocations by approximate size. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks. |
| go_gc_heap_allocs_bytes | counter | Cumulative sum of memory allocated to the heap by the application. |
| go_gc_heap_allocs_objects | counter | Cumulative count of heap allocations triggered by the application. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks. |
| go_gc_heap_frees_by_size_bytes_total | histogram | Distribution of freed heap allocations by approximate size. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks. |
| go_gc_heap_frees_bytes | counter | Cumulative sum of heap memory freed by the garbage collector. |
| go_gc_heap_frees_objects | counter | Cumulative count of heap allocations whose storage was freed by the garbage collector. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks. |
| go_gc_heap_goal_bytes | gauge | Heap size target for the end of the GC cycle. |
| go_gc_heap_objects_objects | gauge | Number of objects, live or unswept, occupying heap memory. |
| go_gc_heap_tiny_allocs_objects | counter | Count of small allocations that are packed together into blocks. These allocations are counted separately from other allocations because each individual allocation is not tracked by the runtime, only their block. Each block is already accounted for in allocs-by-size and frees-by-size. |
| go_gc_pauses_seconds_total | histogram | Distribution individual GC-related stop-the-world pause latencies. |
| go_goroutines | gauge | Number of goroutines that currently exist. |
| go_info | gauge | Information about the Go environment. |
| go_memory_classes_heap_free_bytes | gauge | Memory that is completely free and eligible to be returned to the underlying system, but has not been. This metric is the runtime's estimate of free address space that is backed by physical memory. |
| go_memory_classes_heap_objects_bytes | gauge | Memory occupied by live objects and dead objects that have not yet been marked free by the garbage collector. |
| go_memory_classes_heap_released_bytes | gauge | Memory that is completely free and has been returned to the underlying system. This metric is the runtime's estimate of free address space that is still mapped into the process, but is not backed by physical memory. |
| go_memory_classes_heap_stacks_bytes | gauge | Memory allocated from the heap that is reserved for stack space, whether or not it is currently in-use. |
| go_memory_classes_heap_unused_bytes | gauge | Memory that is reserved for heap objects but is not currently used to hold heap objects. |
| go_memory_classes_metadata_mcache_free_bytes | gauge | Memory that is reserved for runtime mcache structures, but not in-use. |
| go_memory_classes_metadata_mcache_inuse_bytes | gauge | Memory that is occupied by runtime mcache structures that are currently being used. |
| go_memory_classes_metadata_mspan_free_bytes | gauge | Memory that is reserved for runtime mspan structures, but not in-use. |
| go_memory_classes_metadata_mspan_inuse_bytes | gauge | Memory that is occupied by runtime mspan structures that are currently being used. |
| go_memory_classes_metadata_other_bytes | gauge | Memory that is reserved for or used to hold runtime metadata. |
| go_memory_classes_os_stacks_bytes | gauge | Stack memory allocated by the underlying operating system. |
| go_memory_classes_other_bytes | gauge | Memory used by execution trace buffers, structures for debugging the runtime, finalizer and profiler specials, and more. |
| go_memory_classes_profiling_buckets_bytes | gauge | Memory that is used by the stack trace hash map used for profiling. |
| go_memory_classes_total_bytes | gauge | All memory mapped by the Go runtime into the current process as read-write. Note that this does not include memory mapped by code called via cgo or via the syscall package. Sum of all metrics in /memory/classes. |
| go_memstats_alloc_bytes | gauge | Number of bytes allocated and still in use. |
| go_memstats_alloc_bytes | counter | Total number of bytes allocated, even if freed. |
| go_memstats_buck_hash_sys_bytes | gauge | Number of bytes used by the profiling bucket hash table. |
| go_memstats_frees | counter | Total number of frees. |
| go_memstats_gc_cpu_fraction | gauge | The fraction of this program's available CPU time used by the GC since the program started. |
| go_memstats_gc_sys_bytes | gauge | Number of bytes used for garbage collection system metadata. |
| go_memstats_heap_alloc_bytes | gauge | Number of heap bytes allocated and still in use. |
| go_memstats_heap_idle_bytes | gauge | Number of heap bytes waiting to be used. |
| go_memstats_heap_inuse_bytes | gauge | Number of heap bytes that are in use. |
| go_memstats_heap_objects | gauge | Number of allocated objects. |
| go_memstats_heap_released_bytes | gauge | Number of heap bytes released to OS. |
| go_memstats_heap_sys_bytes | gauge | Number of heap bytes obtained from system. |
| go_memstats_last_gc_time_seconds | gauge | Number of seconds since 1970 of last garbage collection. |
| go_memstats_lookups | counter | Total number of pointer lookups. |
| go_memstats_mallocs | counter | Total number of mallocs. |
| go_memstats_mcache_inuse_bytes | gauge | Number of bytes in use by mcache structures. |
| go_memstats_mcache_sys_bytes | gauge | Number of bytes used for mcache structures obtained from system. |
| go_memstats_mspan_inuse_bytes | gauge | Number of bytes in use by mspan structures. |
| go_memstats_mspan_sys_bytes | gauge | Number of bytes used for mspan structures obtained from system. |
| go_memstats_next_gc_bytes | gauge | Number of heap bytes when next garbage collection will take place. |
| go_memstats_other_sys_bytes | gauge | Number of bytes used for other system allocations. |
| go_memstats_stack_inuse_bytes | gauge | Number of bytes in use by the stack allocator. |
| go_memstats_stack_sys_bytes | gauge | Number of bytes obtained from system for stack allocator. |
| go_memstats_sys_bytes | gauge | Number of bytes obtained from system. |
| go_sched_goroutines_goroutines | gauge | Count of live goroutines. |
| go_sched_latencies_seconds | histogram | Distribution of the time goroutines have spent in the scheduler in a runnable state before actually running. |
| go_threads | gauge | Number of OS threads created. |
| grafana_access_evaluation_count | counter | number of evaluation calls |
| grafana_alerting_active_alerts | gauge | amount of active alerts |
| grafana_alerting_active_configurations | gauge | The number of active Alertmanager configurations. |
| grafana_alerting_alerts | gauge | How many alerts by state. |
| grafana_alerting_discovered_configurations | gauge | The number of organizations we've discovered that require an Alertmanager configuration. |
| grafana_alerting_execution_time_milliseconds | summary | summary of alert execution duration |
| grafana_alerting_get_alert_rules_duration_seconds | histogram | The time taken to get all alert rules. |
| grafana_alerting_request_duration_seconds | histogram | Histogram of requests to the Alerting API |
| grafana_alerting_schedule_periodic_duration_seconds | histogram | The time taken to run the scheduler. |
| grafana_alerting_scheduler_behind_seconds | gauge | The total number of seconds the scheduler is behind. |
| grafana_api_admin_user_created | counter | api admin user created counter |
| grafana_api_dashboard_get_milliseconds | summary | summary for dashboard get duration |
| grafana_api_dashboard_save_milliseconds | summary | summary for dashboard save duration |
| grafana_api_dashboard_search_milliseconds | summary | summary for dashboard search duration |
| grafana_api_dashboard_snapshot_create | counter | dashboard snapshots created |
| grafana_api_dashboard_snapshot_external | counter | external dashboard snapshots created |
| grafana_api_dashboard_snapshot_get | counter | loaded dashboards |
| grafana_api_dataproxy_request_all_milliseconds | summary | summary for dataproxy request duration |
| grafana_api_login_oauth | counter | api login oauth counter |
| grafana_api_login_post | counter | api login post counter |
| grafana_api_login_saml | counter | api login saml counter |
| grafana_api_models_dashboard_insert | counter | dashboards inserted |
| grafana_api_org_create | counter | api org created counter |
| grafana_api_response_status | counter | api http response status |
| grafana_api_user_signup_completed | counter | amount of users who completed the signup flow |
| grafana_api_user_signup_invite | counter | amount of users who have been invited |
| grafana_api_user_signup_started | counter | amount of users who started the signup flow |
| grafana_aws_cloudwatch_get_metric_data | counter | counter for getting metric data time series from aws |
| grafana_aws_cloudwatch_get_metric_statistics | counter | counter for getting metric statistics from aws |
| grafana_aws_cloudwatch_list_metrics | counter | counter for getting list of metrics from aws |
| grafana_build_info | gauge | A metric with a constant '1' value labeled by version, revision, branch, and goversion from which Grafana was built |
| grafana_datasource_request_duration_seconds | summary | summary of outgoing datasource requests sent from Grafana |
| grafana_datasource_request_in_flight | gauge | A gauge of outgoing datasource requests currently being sent by Grafana |
| grafana_datasource_request | counter | A counter for outgoing requests for a datasource |
| grafana_datasource_response_size_bytes | summary | summary of datasource response sizes returned to Grafana |
| grafana_db_datasource_query_by_id | counter | counter for getting datasource by id |
| grafana_emails_sent_failed | counter | Number of emails Grafana failed to send |
| grafana_emails_sent | counter | Number of emails sent by Grafana |
| grafana_feature_toggles_info | gauge | info metric that exposes what feature toggles are enabled or not |
| grafana_frontend_boot_css_time_seconds | histogram | Frontend boot initial css load |
| grafana_frontend_boot_first_contentful_paint_time_seconds | histogram | Frontend boot first contentful paint |
| grafana_frontend_boot_first_paint_time_seconds | histogram | Frontend boot first paint |
| grafana_frontend_boot_js_done_time_seconds | histogram | Frontend boot initial js load |
| grafana_frontend_boot_load_time_seconds | histogram | Frontend boot time measurement |
| grafana_http_request_duration_seconds | histogram | Histogram of latencies for HTTP requests. |
| grafana_http_request_in_flight | gauge | A gauge of requests currently being served by Grafana. |
| grafana_instance_start | counter | counter for started instances |
| grafana_ldap_users_sync_execution_time | summary | summary for LDAP users sync execution duration |
| grafana_live_client_command_duration_seconds | summary | clientID command duration summary. |
| grafana_live_client_recover | counter | Count of recover operations. |
| grafana_live_node_action_count | counter | Number of node actions called. |
| grafana_live_node_build | gauge | Node build info. |
| grafana_live_node_messages_received_count | counter | Number of messages received. |
| grafana_live_node_messages_sent_count | counter | Number of messages sent. |
| grafana_live_node_num_channels | gauge | Number of channels with one or more subscribers. |
| grafana_live_node_num_clients | gauge | Number of clients connected. |
| grafana_live_node_num_nodes | gauge | Number of nodes in cluster. |
| grafana_live_node_num_subscriptions | gauge | Number of subscriptions. |
| grafana_live_node_num_users | gauge | Number of unique users connected. |
| grafana_live_transport_connect_count | counter | Number of connections to specific transport. |
| grafana_live_transport_messages_sent | counter | Number of messages sent over specific transport. |
| grafana_page_response_status | counter | page http response status |
| grafana_plugin_build_info | gauge | A metric with a constant '1' value labeled by pluginId, pluginType and version from which Grafana plugin was built |
| grafana_plugin_request_duration_milliseconds | summary | Plugin request duration |
| grafana_plugin_request | counter | The total amount of plugin requests |
| grafana_proxy_response_status | counter | proxy http response status |
| grafana_rendering_queue_size | gauge | size of rendering queue |
| grafana_stat_active_users | gauge | number of active users |
| grafana_stat_total_orgs | gauge | total amount of orgs |
| grafana_stat_total_playlists | gauge | total amount of playlists |
| grafana_stat_total_users | gauge | total amount of users |
| grafana_stat_totals_active_admins | gauge | total amount of active admins |
| grafana_stat_totals_active_editors | gauge | total amount of active editors |
| grafana_stat_totals_active_viewers | gauge | total amount of viewers |
| grafana_stat_totals_admins | gauge | total amount of admins |
| grafana_stat_totals_annotations | gauge | total amount of annotations in the database |
| grafana_stat_totals_dashboard | gauge | total amount of dashboards |
| grafana_stat_totals_dashboard_versions | gauge | total amount of dashboard versions in the database |
| grafana_stat_totals_datasource | gauge | total number of defined datasources, labeled by pluginId |
| grafana_stat_totals_editors | gauge | total amount of editors |
| grafana_stat_totals_folder | gauge | total amount of folders |
| grafana_stat_totals_library_panels | gauge | total amount of library panels in the database |
| grafana_stat_totals_library_variables | gauge | total amount of library variables in the database |
| grafana_stat_totals_viewers | gauge | total amount of viewers |
| net_conntrack_dialer_conn_attempted | counter | Total number of connections attempted by the given dialer a given name. |
| net_conntrack_dialer_conn_closed | counter | Total number of connections closed which originated from the dialer of a given name. |
| net_conntrack_dialer_conn_established | counter | Total number of connections successfully established by the given dialer a given name. |
| net_conntrack_dialer_conn_failed | counter | Total number of connections failed to dial by the dialer a given name. |
| process_cpu_seconds | counter | Total user and system CPU time spent in seconds. |
| process_max_fds | gauge | Maximum number of open file descriptors. |
| process_open_fds | gauge | Number of open file descriptors. |
| process_resident_memory_bytes | gauge | Resident memory size in bytes. |
| process_start_time_seconds | gauge | Start time of the process since unix epoch in seconds. |
| process_virtual_memory_bytes | gauge | Virtual memory size in bytes. |
| process_virtual_memory_max_bytes | gauge | Maximum amount of virtual memory available in bytes. |
| prometheus_template_text_expansion_failures | counter | The total number of template text expansion failures. |
| prometheus_template_text_expansions | counter | The total number of template text expansions. |
<!-- markdownlint-enable line-length -->

## Grafana-operator (Kubernetes v1.21)

### Statistics

<!-- markdownlint-disable line-length -->
| Number of metrics families | Number of time series | Average number of labels per time series | Number of unique label pairs |
| -------------------------- | --------------------- | ---------------------------------------- | ---------------------------- |
| 49 | 1580 | 2.54 | 178 |
<!-- markdownlint-enable line-length -->

### Metrics

<!-- markdownlint-disable line-length -->
| Name | Type, Unit | Description |
| ---- | ---------- | ----------- |
| controller_runtime_active_workers | gauge | Number of currently used workers per controller |
| controller_runtime_max_concurrent_reconciles | gauge | Maximum number of concurrent reconciles per controller |
| controller_runtime_reconcile_errors | counter | Total number of reconciliation errors per controller |
| controller_runtime_reconcile_time_seconds | histogram | Length of time per reconciliation per controller |
| controller_runtime_reconcile | counter | Total number of reconciliations per controller |
| go_gc_duration_seconds | summary | A summary of the pause duration of garbage collection cycles. |
| go_goroutines | gauge | Number of goroutines that currently exist. |
| go_info | gauge | Information about the Go environment. |
| go_memstats_alloc_bytes | gauge | Number of bytes allocated and still in use. |
| go_memstats_alloc_bytes | counter | Total number of bytes allocated, even if freed. |
| go_memstats_buck_hash_sys_bytes | gauge | Number of bytes used by the profiling bucket hash table. |
| go_memstats_frees | counter | Total number of frees. |
| go_memstats_gc_cpu_fraction | gauge | The fraction of this program's available CPU time used by the GC since the program started. |
| go_memstats_gc_sys_bytes | gauge | Number of bytes used for garbage collection system metadata. |
| go_memstats_heap_alloc_bytes | gauge | Number of heap bytes allocated and still in use. |
| go_memstats_heap_idle_bytes | gauge | Number of heap bytes waiting to be used. |
| go_memstats_heap_inuse_bytes | gauge | Number of heap bytes that are in use. |
| go_memstats_heap_objects | gauge | Number of allocated objects. |
| go_memstats_heap_released_bytes | gauge | Number of heap bytes released to OS. |
| go_memstats_heap_sys_bytes | gauge | Number of heap bytes obtained from system. |
| go_memstats_last_gc_time_seconds | gauge | Number of seconds since 1970 of last garbage collection. |
| go_memstats_lookups | counter | Total number of pointer lookups. |
| go_memstats_mallocs | counter | Total number of mallocs. |
| go_memstats_mcache_inuse_bytes | gauge | Number of bytes in use by mcache structures. |
| go_memstats_mcache_sys_bytes | gauge | Number of bytes used for mcache structures obtained from system. |
| go_memstats_mspan_inuse_bytes | gauge | Number of bytes in use by mspan structures. |
| go_memstats_mspan_sys_bytes | gauge | Number of bytes used for mspan structures obtained from system. |
| go_memstats_next_gc_bytes | gauge | Number of heap bytes when next garbage collection will take place. |
| go_memstats_other_sys_bytes | gauge | Number of bytes used for other system allocations. |
| go_memstats_stack_inuse_bytes | gauge | Number of bytes in use by the stack allocator. |
| go_memstats_stack_sys_bytes | gauge | Number of bytes obtained from system for stack allocator. |
| go_memstats_sys_bytes | gauge | Number of bytes obtained from system. |
| go_threads | gauge | Number of OS threads created. |
| process_cpu_seconds | counter | Total user and system CPU time spent in seconds. |
| process_max_fds | gauge | Maximum number of open file descriptors. |
| process_open_fds | gauge | Number of open file descriptors. |
| process_resident_memory_bytes | gauge | Resident memory size in bytes. |
| process_start_time_seconds | gauge | Start time of the process since unix epoch in seconds. |
| process_virtual_memory_bytes | gauge | Virtual memory size in bytes. |
| process_virtual_memory_max_bytes | gauge | Maximum amount of virtual memory available in bytes. |
| rest_client_request_latency_seconds | histogram | Request latency in seconds. Broken down by verb and URL. |
| rest_client_requests | counter | Number of HTTP requests, partitioned by status code, method, and host. |
| workqueue_adds | counter | Total number of adds handled by workqueue |
| workqueue_depth | gauge | Current depth of workqueue |
| workqueue_longest_running_processor_seconds | gauge | How many seconds has the longest running processor for workqueue been running. |
| workqueue_queue_duration_seconds | histogram | How long in seconds an item stays in workqueue before being requested |
| workqueue_retries | counter | Total number of retries handled by workqueue |
| workqueue_unfinished_work_seconds | gauge | How many seconds of work has been done that is in progress and hasn't been observed by work_duration. Large values indicate stuck threads. One can deduce the number of stuck threads by observing the rate at which this increases. |
| workqueue_work_duration_seconds | histogram | How long in seconds processing an item from workqueue takes. |
<!-- markdownlint-enable line-length -->

## Kube-apiserver (Kubernetes v1.21)

### Statistics

<!-- markdownlint-disable line-length -->
| Number of metrics families | Number of time series | Average number of labels per time series | Number of unique label pairs |
| -------------------------- | --------------------- | ---------------------------------------- | ---------------------------- |
| 119 | 44198 | 7.5 | 1091 |
<!-- markdownlint-enable line-length -->

### Metrics

<!-- markdownlint-disable line-length -->
| Name | Type, Unit | Description |
| ---- | ---------- | ----------- |
| aggregator_openapi_v2_regeneration_count | counter | ALPHA Counter of OpenAPI v2 spec regeneration count broken down by causing APIService name and reason. |
| aggregator_openapi_v2_regeneration_duration | gauge | ALPHA Gauge of OpenAPI v2 spec regeneration duration in seconds. |
| aggregator_unavailable_apiservice | gauge | ALPHA Gauge of APIServices which are marked as unavailable broken down by APIService name. |
| apiextensions_openapi_v2_regeneration_count | counter | ALPHA Counter of OpenAPI v2 spec regeneration count broken down by causing CRD name and reason. |
| apiserver_admission_controller_admission_duration_seconds | histogram | ALPHA Admission controller latency histogram in seconds, identified by name and broken out for each operation and API resource and type (validate or admit). |
| apiserver_admission_step_admission_duration_seconds | histogram | ALPHA Admission sub-step latency histogram in seconds, broken out for each operation and API resource and step type (validate or admit). |
| apiserver_admission_step_admission_duration_seconds_summary | summary | ALPHA Admission sub-step latency summary in seconds, broken out for each operation and API resource and step type (validate or admit). |
| apiserver_admission_webhook_admission_duration_seconds | histogram | ALPHA Admission webhook latency histogram in seconds, identified by name and broken out for each operation and API resource and type (validate or admit). |
| apiserver_audit_event | counter | ALPHA Counter of audit events generated and sent to the audit backend. |
| apiserver_audit_requests_rejected | counter | ALPHA Counter of apiserver requests rejected due to an error in audit logging backend. |
| apiserver_client_certificate_expiration_seconds | histogram | ALPHA Distribution of the remaining lifetime on the certificate used to authenticate a request. |
| apiserver_current_inflight_requests | gauge | ALPHA Maximal number of currently used inflight request limit of this apiserver per request kind in last second. |
| apiserver_current_inqueue_requests | gauge | ALPHA Maximal number of queued requests in this apiserver per request kind in last second. |
| apiserver_envelope_encryption_dek_cache_fill_percent | gauge | ALPHA Percent of the cache slots currently occupied by cached DEKs. |
| apiserver_flowcontrol_current_executing_requests | gauge | ALPHA Number of requests currently executing in the API Priority and Fairness system |
| apiserver_flowcontrol_current_inqueue_requests | gauge | ALPHA Number of requests currently pending in queues of the API Priority and Fairness system |
| apiserver_flowcontrol_dispatched_requests | counter | ALPHA Number of requests released by API Priority and Fairness system for service |
| apiserver_flowcontrol_priority_level_request_count_samples | histogram | ALPHA Periodic observations of the number of requests |
| apiserver_flowcontrol_priority_level_request_count_watermarks | histogram | ALPHA Watermarks of the number of requests |
| apiserver_flowcontrol_read_vs_write_request_count_samples | histogram | ALPHA Periodic observations of the number of requests |
| apiserver_flowcontrol_read_vs_write_request_count_watermarks | histogram | ALPHA Watermarks of the number of requests |
| apiserver_flowcontrol_request_concurrency_limit | gauge | ALPHA Shared concurrency limit in the API Priority and Fairness system |
| apiserver_flowcontrol_request_execution_seconds | histogram | ALPHA Duration of request execution in the API Priority and Fairness system |
| apiserver_flowcontrol_request_queue_length_after_enqueue | histogram | ALPHA Length of queue in the API Priority and Fairness system, as seen by each request after it is enqueued |
| apiserver_flowcontrol_request_wait_duration_seconds | histogram | ALPHA Length of time a request spent waiting in its queue |
| apiserver_init_events | counter | ALPHA Counter of init events processed in watchcache broken by resource type. |
| apiserver_longrunning_gauge | gauge | ALPHA Gauge of all active long-running apiserver requests broken out by verb, group, version, resource, scope and component. Not all requests are tracked this way. |
| apiserver_registered_watchers | gauge | ALPHA Number of currently registered watchers for a given resources |
| apiserver_request_aborts | counter | ALPHA Number of requests which apiserver aborted possibly due to a timeout, for each group, version, verb, resource, subresource and scope |
| apiserver_request_duration_seconds | histogram | STABLE Response latency distribution in seconds for each verb, dry run value, group, version, resource, subresource, scope and component. |
| apiserver_request_filter_duration_seconds | histogram | ALPHA Request filter latency distribution in seconds, for each filter type |
| apiserver_request_terminations | counter | ALPHA Number of requests which apiserver terminated in self-defense. |
| apiserver_request | counter | STABLE Counter of apiserver requests broken out for each verb, dry run value, group, version, resource, scope, component, and HTTP response code. |
| apiserver_requested_deprecated_apis | gauge | ALPHA Gauge of deprecated APIs that have been requested, broken out by API group, version, resource, subresource, and removed_release. |
| apiserver_response_sizes | histogram | ALPHA Response size distribution in bytes for each group, version, verb, resource, subresource, scope and component. |
| apiserver_selfrequest | counter | ALPHA Counter of apiserver self-requests broken out for each verb, API resource and subresource. |
| apiserver_storage_data_key_generation_duration_seconds | histogram | ALPHA Latencies in seconds of data encryption key(DEK) generation operations. |
| apiserver_storage_data_key_generation_failures | counter | ALPHA Total number of failed data encryption key(DEK) generation operations. |
| apiserver_storage_envelope_transformation_cache_misses | counter | ALPHA Total number of cache misses while accessing key decryption key(KEK). |
| apiserver_storage_objects | gauge | STABLE Number of stored objects at the time of last check split by kind. |
| apiserver_tls_handshake_errors | counter | ALPHA Number of requests dropped with 'TLS handshake error from' error |
| apiserver_watch_events_sizes | histogram | ALPHA Watch event size distribution in bytes |
| apiserver_watch_events | counter | ALPHA Number of events sent in watch clients |
| authenticated_user_requests | counter | ALPHA Counter of authenticated requests broken out by username. |
| authentication_attempts | counter | ALPHA Counter of authenticated attempts. |
| authentication_duration_seconds | histogram | ALPHA Authentication duration in seconds broken out by result. |
| authentication_token_cache_active_fetch_count | gauge | ALPHA |
| authentication_token_cache_fetch | counter | ALPHA |
| authentication_token_cache_request_duration_seconds | histogram | ALPHA |
| authentication_token_cache_request | counter | ALPHA |
| etcd_bookmark_counts | gauge | ALPHA Number of etcd bookmarks (progress notify events) split by kind. |
| etcd_db_total_size_in_bytes | gauge | ALPHA Total size of the etcd database file physically allocated in bytes. |
| etcd_lease_object_counts | histogram | ALPHA Number of objects attached to a single etcd lease. |
| etcd_object_counts | gauge | ALPHA Number of stored objects at the time of last check split by kind. This metric is replaced by apiserver_storage_object_counts. |
| etcd_request_duration_seconds | histogram | ALPHA Etcd request latency in seconds for each operation and object type. |
| get_token_count | counter | ALPHA Counter of total Token() requests to the alternate token source |
| get_token_fail_count | counter | ALPHA Counter of failed Token() requests to the alternate token source |
| go_gc_duration_seconds | summary | A summary of the pause duration of garbage collection cycles. |
| go_goroutines | gauge | Number of goroutines that currently exist. |
| go_info | gauge | Information about the Go environment. |
| go_memstats_alloc_bytes | gauge | Number of bytes allocated and still in use. |
| go_memstats_alloc_bytes | counter | Total number of bytes allocated, even if freed. |
| go_memstats_buck_hash_sys_bytes | gauge | Number of bytes used by the profiling bucket hash table. |
| go_memstats_frees | counter | Total number of frees. |
| go_memstats_gc_cpu_fraction | gauge | The fraction of this program's available CPU time used by the GC since the program started. |
| go_memstats_gc_sys_bytes | gauge | Number of bytes used for garbage collection system metadata. |
| go_memstats_heap_alloc_bytes | gauge | Number of heap bytes allocated and still in use. |
| go_memstats_heap_idle_bytes | gauge | Number of heap bytes waiting to be used. |
| go_memstats_heap_inuse_bytes | gauge | Number of heap bytes that are in use. |
| go_memstats_heap_objects | gauge | Number of allocated objects. |
| go_memstats_heap_released_bytes | gauge | Number of heap bytes released to OS. |
| go_memstats_heap_sys_bytes | gauge | Number of heap bytes obtained from system. |
| go_memstats_last_gc_time_seconds | gauge | Number of seconds since 1970 of last garbage collection. |
| go_memstats_lookups | counter | Total number of pointer lookups. |
| go_memstats_mallocs | counter | Total number of mallocs. |
| go_memstats_mcache_inuse_bytes | gauge | Number of bytes in use by mcache structures. |
| go_memstats_mcache_sys_bytes | gauge | Number of bytes used for mcache structures obtained from system. |
| go_memstats_mspan_inuse_bytes | gauge | Number of bytes in use by mspan structures. |
| go_memstats_mspan_sys_bytes | gauge | Number of bytes used for mspan structures obtained from system. |
| go_memstats_next_gc_bytes | gauge | Number of heap bytes when next garbage collection will take place. |
| go_memstats_other_sys_bytes | gauge | Number of bytes used for other system allocations. |
| go_memstats_stack_inuse_bytes | gauge | Number of bytes in use by the stack allocator. |
| go_memstats_stack_sys_bytes | gauge | Number of bytes obtained from system for stack allocator. |
| go_memstats_sys_bytes | gauge | Number of bytes obtained from system. |
| go_threads | gauge | Number of OS threads created. |
| grpc_client_handled | counter | Total number of RPCs completed by the client, regardless of success or failure. |
| grpc_client_msg_received | counter | Total number of RPC stream messages received by the client. |
| grpc_client_msg_sent | counter | Total number of gRPC stream messages sent by the client. |
| grpc_client_started | counter | Total number of RPCs started on the client. |
| kube_apiserver_pod_logs_pods_logs_backend_tls_failure | counter | ALPHA Total number of requests for pods/logs that failed due to kubelet server TLS verification |
| kube_apiserver_pod_logs_pods_logs_insecure_backend | counter | ALPHA Total number of requests for pods/logs sliced by usage type: enforce_tls, skip_tls_allowed, skip_tls_denied |
| kubernetes_build_info | gauge | ALPHA A metric with a constant '1' value labeled by major, minor, git version, git commit, git tree state, build date, Go version, and compiler from which Kubernetes was built, and platform on which it is running. |
| node_authorizer_graph_actions_duration_seconds | histogram | ALPHA Histogram of duration of graph actions in node authorizer. |
| process_cpu_seconds | counter | Total user and system CPU time spent in seconds. |
| process_max_fds | gauge | Maximum number of open file descriptors. |
| process_open_fds | gauge | Number of open file descriptors. |
| process_resident_memory_bytes | gauge | Resident memory size in bytes. |
| process_start_time_seconds | gauge | Start time of the process since unix epoch in seconds. |
| process_virtual_memory_bytes | gauge | Virtual memory size in bytes. |
| process_virtual_memory_max_bytes | gauge | Maximum amount of virtual memory available in bytes. |
| rest_client_exec_plugin_certificate_rotation_age | histogram | ALPHA Histogram of the number of seconds the last auth exec plugin client certificate lived before being rotated. If auth exec plugin client certificates are unused, histogram will contain no data. |
| rest_client_exec_plugin_ttl_seconds | gauge | ALPHA Gauge of the shortest TTL (time-to-live) of the client certificate(s) managed by the auth exec plugin. The value is in seconds until certificate expiry (negative if already expired). If auth exec plugins are unused or manage no TLS certificates, the value will be +INF. |
| rest_client_request_duration_seconds | histogram | ALPHA Request latency in seconds. Broken down by verb and URL. |
| rest_client_requests | counter | ALPHA Number of HTTP requests, partitioned by status code, method, and host. |
| serviceaccount_legacy_tokens | counter | ALPHA Cumulative legacy service account tokens used |
| serviceaccount_stale_tokens | counter | ALPHA Cumulative stale projected service account tokens used |
| serviceaccount_valid_tokens | counter | ALPHA Cumulative valid projected service account tokens used |
| ssh_tunnel_open_count | counter | ALPHA Counter of ssh tunnel total open attempts |
| ssh_tunnel_open_fail_count | counter | ALPHA Counter of ssh tunnel failed open attempts |
| watch_cache_capacity | gauge | ALPHA Total capacity of watch cache broken by resource type. |
| watch_cache_capacity_decrease | counter | ALPHA Total number of watch cache capacity decrease events broken by resource type. |
| watch_cache_capacity_increase | counter | ALPHA Total number of watch cache capacity increase events broken by resource type. |
| workqueue_adds | counter | ALPHA Total number of adds handled by workqueue |
| workqueue_depth | gauge | ALPHA Current depth of workqueue |
| workqueue_longest_running_processor_seconds | gauge | ALPHA How many seconds has the longest running processor for workqueue been running. |
| workqueue_queue_duration_seconds | histogram | ALPHA How long in seconds an item stays in workqueue before being requested. |
| workqueue_retries | counter | ALPHA Total number of retries handled by workqueue |
| workqueue_unfinished_work_seconds | gauge | ALPHA How many seconds of work has done that is in progress and hasn't been observed by work_duration. Large values indicate stuck threads. One can deduce the number of stuck threads by observing the rate at which this increases. |
| workqueue_work_duration_seconds | histogram | ALPHA How long in seconds processing an item from workqueue takes. |
<!-- markdownlint-enable line-length -->

## Kube-state-metrics (Kubernetes v1.21)

### Statistics

<!-- markdownlint-disable line-length -->
| Number of metrics families | Number of time series | Average number of labels per time series | Number of unique label pairs |
| -------------------------- | --------------------- | ---------------------------------------- | ---------------------------- |
| 220 | 34425 | 3.1 | 5441 |
<!-- markdownlint-enable line-length -->

### Metrics

<!-- markdownlint-disable line-length -->
| Name | Type, Unit | Description |
| ---- | ---------- | ----------- |
| kube_configmap_annotations | gauge | Kubernetes annotations converted to Prometheus labels. |
| kube_configmap_labels | gauge | Kubernetes labels converted to Prometheus labels. |
| kube_configmap_info | gauge | Information about configmap. |
| kube_configmap_created | gauge | Unix creation timestamp |
| kube_configmap_metadata_resource_version | gauge | Resource version representing a specific version of the configmap. |
| kube_cronjob_annotations | gauge | Kubernetes annotations converted to Prometheus labels. |
| kube_cronjob_labels | gauge | Kubernetes labels converted to Prometheus labels. |
| kube_cronjob_info | gauge | Info about cronjob. |
| kube_cronjob_created | gauge | Unix creation timestamp |
| kube_cronjob_status_active | gauge | Active holds pointers to currently running jobs. |
| kube_cronjob_status_last_schedule_time | gauge | LastScheduleTime keeps information of when was the last time the job was successfully scheduled. |
| kube_cronjob_spec_suspend | gauge | Suspend flag tells the controller to suspend subsequent executions. |
| kube_cronjob_spec_starting_deadline_seconds | gauge | Deadline in seconds for starting the job if it misses scheduled time for any reason. |
| kube_cronjob_next_schedule_time | gauge | Next time the cronjob should be scheduled. The time after lastScheduleTime, or after the cron job's creation time if it's never been scheduled. Use this to determine if the job is delayed. |
| kube_cronjob_metadata_resource_version | gauge | Resource version representing a specific version of the cronjob. |
| kube_cronjob_spec_successful_job_history_limit | gauge | Successful job history limit tells the controller how many completed jobs should be preserved. |
| kube_cronjob_spec_failed_job_history_limit | gauge | Failed job history limit tells the controller how many failed jobs should be preserved. |
| kube_daemonset_created | gauge | Unix creation timestamp |
| kube_daemonset_status_current_number_scheduled | gauge | The number of nodes running at least one daemon pod and are supposed to. |
| kube_daemonset_status_desired_number_scheduled | gauge | The number of nodes that should be running the daemon pod. |
| kube_daemonset_status_number_available | gauge | The number of nodes that should be running the daemon pod and have one or more of the daemon pod running and available |
| kube_daemonset_status_number_misscheduled | gauge | The number of nodes running a daemon pod but are not supposed to. |
| kube_daemonset_status_number_ready | gauge | The number of nodes that should be running the daemon pod and have one or more of the daemon pod running and ready. |
| kube_daemonset_status_number_unavailable | gauge | The number of nodes that should be running the daemon pod and have none of the daemon pod running and available |
| kube_daemonset_status_observed_generation | gauge | The most recent generation observed by the daemon set controller. |
| kube_daemonset_status_updated_number_scheduled | gauge | The total number of nodes that are running updated daemon pod |
| kube_daemonset_metadata_generation | gauge | Sequence number representing a specific generation of the desired state. |
| kube_daemonset_annotations | gauge | Kubernetes annotations converted to Prometheus labels. |
| kube_daemonset_labels | gauge | Kubernetes labels converted to Prometheus labels. |
| kube_deployment_created | gauge | Unix creation timestamp |
| kube_deployment_status_replicas | gauge | The number of replicas per deployment. |
| kube_deployment_status_replicas_ready | gauge | The number of ready replicas per deployment. |
| kube_deployment_status_replicas_available | gauge | The number of available replicas per deployment. |
| kube_deployment_status_replicas_unavailable | gauge | The number of unavailable replicas per deployment. |
| kube_deployment_status_replicas_updated | gauge | The number of updated replicas per deployment. |
| kube_deployment_status_observed_generation | gauge | The generation observed by the deployment controller. |
| kube_deployment_status_condition | gauge | The current status conditions of a deployment. |
| kube_deployment_spec_replicas | gauge | Number of desired pods for a deployment. |
| kube_deployment_spec_paused | gauge | Whether the deployment is paused and will not be processed by the deployment controller. |
| kube_deployment_spec_strategy_rollingupdate_max_unavailable | gauge | Maximum number of unavailable replicas during a rolling update of a deployment. |
| kube_deployment_spec_strategy_rollingupdate_max_surge | gauge | Maximum number of replicas that can be scheduled above the desired number of replicas during a rolling update of a deployment. |
| kube_deployment_metadata_generation | gauge | Sequence number representing a specific generation of the desired state. |
| kube_deployment_annotations | gauge | Kubernetes annotations converted to Prometheus labels. |
| kube_deployment_labels | gauge | Kubernetes labels converted to Prometheus labels. |
| kube_endpoint_info | gauge | Information about endpoint. |
| kube_endpoint_created | gauge | Unix creation timestamp |
| kube_endpoint_annotations | gauge | Kubernetes annotations converted to Prometheus labels. |
| kube_endpoint_labels | gauge | Kubernetes labels converted to Prometheus labels. |
| kube_endpoint_address_available | gauge | Number of addresses available in endpoint. |
| kube_endpoint_address_not_ready | gauge | Number of addresses not ready in endpoint |
| kube_endpoint_ports | gauge | Information about the Endpoint ports. |
| kube_horizontalpodautoscaler_info | gauge | Information about this autoscaler. |
| kube_horizontalpodautoscaler_metadata_generation | gauge | The generation observed by the HorizontalPodAutoscaler controller. |
| kube_horizontalpodautoscaler_spec_max_replicas | gauge | Upper limit for the number of pods that can be set by the autoscaler; cannot be smaller than MinReplicas. |
| kube_horizontalpodautoscaler_spec_min_replicas | gauge | Lower limit for the number of pods that can be set by the autoscaler, default 1. |
| kube_horizontalpodautoscaler_spec_target_metric | gauge | The metric specifications used by this autoscaler when calculating the desired replica count. |
| kube_horizontalpodautoscaler_status_current_replicas | gauge | Current number of replicas of pods managed by this autoscaler. |
| kube_horizontalpodautoscaler_status_desired_replicas | gauge | Desired number of replicas of pods managed by this autoscaler. |
| kube_horizontalpodautoscaler_annotations | gauge | Kubernetes annotations converted to Prometheus labels. |
| kube_horizontalpodautoscaler_labels | gauge | Kubernetes labels converted to Prometheus labels. |
| kube_horizontalpodautoscaler_status_condition | gauge | The condition of this autoscaler. |
| kube_ingress_info | gauge | Information about ingress. |
| kube_ingress_annotations | gauge | Kubernetes annotations converted to Prometheus labels. |
| kube_ingress_labels | gauge | Kubernetes labels converted to Prometheus labels. |
| kube_ingress_created | gauge | Unix creation timestamp |
| kube_ingress_metadata_resource_version | gauge | Resource version representing a specific version of ingress. |
| kube_ingress_path | gauge | Ingress host, paths and backend service information. |
| kube_ingress_tls | gauge | Ingress TLS host and secret information. |
| kube_job_annotations | gauge | Kubernetes annotations converted to Prometheus labels. |
| kube_job_labels | gauge | Kubernetes labels converted to Prometheus labels. |
| kube_job_info | gauge | Information about job. |
| kube_job_created | gauge | Unix creation timestamp |
| kube_job_spec_parallelism | gauge | The maximum desired number of pods the job should run at any given time. |
| kube_job_spec_completions | gauge | The desired number of successfully finished pods the job should be run with. |
| kube_job_spec_active_deadline_seconds | gauge | The duration in seconds relative to the startTime that the job may be active before the system tries to terminate it. |
| kube_job_status_succeeded | gauge | The number of pods which reached Phase Succeeded. |
| kube_job_status_failed | gauge | The number of pods which reached Phase Failed and the reason for failure. |
| kube_job_status_active | gauge | The number of actively running pods. |
| kube_job_complete | gauge | The job has completed its execution. |
| kube_job_failed | gauge | The job has failed its execution. |
| kube_job_status_start_time | gauge | StartTime represents time when the job was acknowledged by the Job Manager. |
| kube_job_status_completion_time | gauge | CompletionTime represents time when the job was completed. |
| kube_job_owner | gauge | Information about the Job's owner. |
| kube_limitrange | gauge | Information about limit range. |
| kube_limitrange_created | gauge | Unix creation timestamp |
| kube_namespace_created | gauge | Unix creation timestamp |
| kube_namespace_annotations | gauge | Kubernetes annotations converted to Prometheus labels. |
| kube_namespace_labels | gauge | Kubernetes labels converted to Prometheus labels. |
| kube_namespace_status_phase | gauge | kubernetes namespace status phase. |
| kube_namespace_status_condition | gauge | The condition of a namespace. |
| kube_networkpolicy_created | gauge | Unix creation timestamp of network policy |
| kube_networkpolicy_annotations | gauge | Kubernetes annotations converted to Prometheus labels. |
| kube_networkpolicy_labels | gauge | Kubernetes labels converted to Prometheus labels. |
| kube_networkpolicy_spec_ingress_rules | gauge | Number of ingress rules on the networkpolicy |
| kube_networkpolicy_spec_egress_rules | gauge | Number of egress rules on the networkpolicy |
| kube_node_created | gauge | Unix creation timestamp |
| kube_node_info | gauge | Information about a cluster node. |
| kube_node_annotations | gauge | Kubernetes annotations converted to Prometheus labels. |
| kube_node_labels | gauge | Kubernetes labels converted to Prometheus labels. |
| kube_node_role | gauge | The role of a cluster node. |
| kube_node_spec_taint | gauge | The taint of a cluster node. |
| kube_node_spec_unschedulable | gauge | Whether a node can schedule new pods. |
| kube_node_status_allocatable | gauge | The allocatable for different resources of a node that are available for scheduling. |
| kube_node_status_capacity | gauge | The capacity for different resources of a node. |
| kube_node_status_condition | gauge | The condition of a cluster node. |
| kube_persistentvolumeclaim_labels | gauge | Kubernetes labels converted to Prometheus labels. |
| kube_persistentvolumeclaim_annotations | gauge | Kubernetes annotations converted to Prometheus labels. |
| kube_persistentvolumeclaim_info | gauge | Information about persistent volume claim. |
| kube_persistentvolumeclaim_status_phase | gauge | The phase the persistent volume claim is currently in. |
| kube_persistentvolumeclaim_resource_requests_storage_bytes | gauge | The capacity of storage requested by the persistent volume claim. |
| kube_persistentvolumeclaim_access_mode | gauge | The access mode(s) specified by the persistent volume claim. |
| kube_persistentvolumeclaim_status_condition | gauge | Information about status of different conditions of persistent volume claim. |
| kube_persistentvolume_claim_ref | gauge | Information about the Persistent Volume Claim Reference. |
| kube_persistentvolume_annotations | gauge | Kubernetes annotations converted to Prometheus labels. |
| kube_persistentvolume_labels | gauge | Kubernetes labels converted to Prometheus labels. |
| kube_persistentvolume_status_phase | gauge | The phase indicates if a volume is available, bound to a claim, or released by a claim. |
| kube_persistentvolume_info | gauge | Information about persistentvolume. |
| kube_persistentvolume_capacity_bytes | gauge | Persistentvolume capacity in bytes. |
| kube_poddisruptionbudget_annotations | gauge | Kubernetes annotations converted to Prometheus labels. |
| kube_poddisruptionbudget_labels | gauge | Kubernetes labels converted to Prometheus labels. |
| kube_poddisruptionbudget_created | gauge | Unix creation timestamp |
| kube_poddisruptionbudget_status_current_healthy | gauge | Current number of healthy pods |
| kube_poddisruptionbudget_status_desired_healthy | gauge | Minimum desired number of healthy pods |
| kube_poddisruptionbudget_status_pod_disruptions_allowed | gauge | Number of pod disruptions that are currently allowed |
| kube_poddisruptionbudget_status_expected_pods | gauge | Total number of pods counted by this disruption budget |
| kube_poddisruptionbudget_status_observed_generation | gauge | Most recent generation observed when updating this PDB status |
| kube_pod_completion_time | gauge | Completion time in unix timestamp for a pod. |
| kube_pod_container_info | gauge | Information about a container in a pod. |
| kube_pod_container_resource_limits | gauge | The number of requested limit resource by a container. |
| kube_pod_container_resource_requests | gauge | The number of requested request resource by a container. |
| kube_pod_container_state_started | gauge | Start time in unix timestamp for a pod container. |
| kube_pod_container_status_last_terminated_reason | gauge | Describes the last reason the container was in terminated state. |
| kube_pod_container_status_ready | gauge | Describes whether the containers readiness check succeeded. |
| kube_pod_container_status_restarts | counter | The number of container restarts per container. |
| kube_pod_container_status_running | gauge | Describes whether the container is currently in running state. |
| kube_pod_container_status_terminated | gauge | Describes whether the container is currently in terminated state. |
| kube_pod_container_status_terminated_reason | gauge | Describes the reason the container is currently in terminated state. |
| kube_pod_container_status_waiting | gauge | Describes whether the container is currently in waiting state. |
| kube_pod_container_status_waiting_reason | gauge | Describes the reason the container is currently in waiting state. |
| kube_pod_created | gauge | Unix creation timestamp |
| kube_pod_deletion_timestamp | gauge | Unix deletion timestamp |
| kube_pod_info | gauge | Information about pod. |
| kube_pod_init_container_info | gauge | Information about an init container in a pod. |
| kube_pod_init_container_resource_limits | gauge | The number of requested limit resource by an init container. |
| kube_pod_init_container_resource_requests | gauge | The number of requested request resource by an init container. |
| kube_pod_init_container_status_last_terminated_reason | gauge | Describes the last reason the init container was in terminated state. |
| kube_pod_init_container_status_ready | gauge | Describes whether the init containers readiness check succeeded. |
| kube_pod_init_container_status_restarts | counter | The number of restarts for the init container. |
| kube_pod_init_container_status_running | gauge | Describes whether the init container is currently in running state. |
| kube_pod_init_container_status_terminated | gauge | Describes whether the init container is currently in terminated state. |
| kube_pod_init_container_status_terminated_reason | gauge | Describes the reason the init container is currently in terminated state. |
| kube_pod_init_container_status_waiting | gauge | Describes whether the init container is currently in waiting state. |
| kube_pod_init_container_status_waiting_reason | gauge | Describes the reason the init container is currently in waiting state. |
| kube_pod_annotations | gauge | Kubernetes annotations converted to Prometheus labels. |
| kube_pod_labels | gauge | Kubernetes labels converted to Prometheus labels. |
| kube_pod_overhead_cpu_cores | gauge | The pod overhead in regards to cpu cores associated with running a pod. |
| kube_pod_overhead_memory_bytes | gauge | The pod overhead in regards to memory associated with running a pod. |
| kube_pod_owner | gauge | Information about the Pod's owner. |
| kube_pod_restart_policy | gauge | Describes the restart policy in use by this pod. |
| kube_pod_runtimeclass_name_info | gauge | The runtimeclass associated with the pod. |
| kube_pod_spec_volumes_persistentvolumeclaims_info | gauge | Information about persistentvolumeclaim volumes in a pod. |
| kube_pod_spec_volumes_persistentvolumeclaims_readonly | gauge | Describes whether a persistentvolumeclaim is mounted read only. |
| kube_pod_start_time | gauge | Start time in unix timestamp for a pod. |
| kube_pod_status_phase | gauge | The pods current phase. |
| kube_pod_status_ready | gauge | Describes whether the pod is ready to serve requests. |
| kube_pod_status_reason | gauge | The pod status reasons |
| kube_pod_status_scheduled | gauge | Describes the status of the scheduling process for the pod. |
| kube_pod_status_scheduled_time | gauge | Unix timestamp when pod moved into scheduled status |
| kube_pod_status_unschedulable | gauge | Describes the unschedulable status for the pod. |
| kube_replicaset_created | gauge | Unix creation timestamp |
| kube_replicaset_status_replicas | gauge | The number of replicas per ReplicaSet. |
| kube_replicaset_status_fully_labeled_replicas | gauge | The number of fully labeled replicas per ReplicaSet. |
| kube_replicaset_status_ready_replicas | gauge | The number of ready replicas per ReplicaSet. |
| kube_replicaset_status_observed_generation | gauge | The generation observed by the ReplicaSet controller. |
| kube_replicaset_spec_replicas | gauge | Number of desired pods for a ReplicaSet. |
| kube_replicaset_metadata_generation | gauge | Sequence number representing a specific generation of the desired state. |
| kube_replicaset_owner | gauge | Information about the ReplicaSet's owner. |
| kube_replicaset_annotations | gauge | Kubernetes annotations converted to Prometheus labels. |
| kube_replicaset_labels | gauge | Kubernetes labels converted to Prometheus labels. |
| kube_replicationcontroller_created | gauge | Unix creation timestamp |
| kube_replicationcontroller_status_replicas | gauge | The number of replicas per ReplicationController. |
| kube_replicationcontroller_status_fully_labeled_replicas | gauge | The number of fully labeled replicas per ReplicationController. |
| kube_replicationcontroller_status_ready_replicas | gauge | The number of ready replicas per ReplicationController. |
| kube_replicationcontroller_status_available_replicas | gauge | The number of available replicas per ReplicationController. |
| kube_replicationcontroller_status_observed_generation | gauge | The generation observed by the ReplicationController controller. |
| kube_replicationcontroller_spec_replicas | gauge | Number of desired pods for a ReplicationController. |
| kube_replicationcontroller_metadata_generation | gauge | Sequence number representing a specific generation of the desired state. |
| kube_replicationcontroller_owner | gauge | Information about the ReplicationController's owner. |
| kube_resourcequota_created | gauge | Unix creation timestamp |
| kube_resourcequota | gauge | Information about resource quota. |
| kube_secret_info | gauge | Information about secret. |
| kube_secret_type | gauge | Type about secret. |
| kube_secret_annotations | gauge | Kubernetes annotations converted to Prometheus labels. |
| kube_secret_labels | gauge | Kubernetes labels converted to Prometheus labels. |
| kube_secret_created | gauge | Unix creation timestamp |
| kube_secret_metadata_resource_version | gauge | Resource version representing a specific version of secret. |
| kube_service_info | gauge | Information about service. |
| kube_service_created | gauge | Unix creation timestamp |
| kube_service_spec_type | gauge | Type about service. |
| kube_service_annotations | gauge | Kubernetes annotations converted to Prometheus labels. |
| kube_service_labels | gauge | Kubernetes labels converted to Prometheus labels. |
| kube_service_spec_external_ip | gauge | Service external ips. One series for each ip |
| kube_service_status_load_balancer_ingress | gauge | Service load balancer ingress status |
| kube_statefulset_created | gauge | Unix creation timestamp |
| kube_statefulset_status_replicas | gauge | The number of replicas per StatefulSet. |
| kube_statefulset_status_replicas_available | gauge | The number of available replicas per StatefulSet. |
| kube_statefulset_status_replicas_current | gauge | The number of current replicas per StatefulSet. |
| kube_statefulset_status_replicas_ready | gauge | The number of ready replicas per StatefulSet. |
| kube_statefulset_status_replicas_updated | gauge | The number of updated replicas per StatefulSet. |
| kube_statefulset_status_observed_generation | gauge | The generation observed by the StatefulSet controller. |
| kube_statefulset_replicas | gauge | Number of desired pods for a StatefulSet. |
| kube_statefulset_metadata_generation | gauge | Sequence number representing a specific generation of the desired state for the StatefulSet. |
| kube_statefulset_annotations | gauge | Kubernetes annotations converted to Prometheus labels. |
| kube_statefulset_labels | gauge | Kubernetes labels converted to Prometheus labels. |
| kube_statefulset_status_current_revision | gauge | Indicates the version of the StatefulSet used to generate Pods in the sequence [0,currentReplicas). |
| kube_statefulset_status_update_revision | gauge | Indicates the version of the StatefulSet used to generate Pods in the sequence [replicas-updatedReplicas,replicas) |
| kube_storageclass_info | gauge | Information about storageclass. |
| kube_storageclass_created | gauge | Unix creation timestamp |
| kube_storageclass_annotations | gauge | Kubernetes annotations converted to Prometheus labels. |
| kube_storageclass_labels | gauge | Kubernetes labels converted to Prometheus labels. |
<!-- markdownlint-enable line-length -->

## Kubelet (Kubernetes v1.21)

### Statistics

<!-- markdownlint-disable line-length -->
| Number of metrics families | Number of time series | Average number of labels per time series | Number of unique label pairs |
| -------------------------- | --------------------- | ---------------------------------------- | ---------------------------- |
| 90 | 1866 | 2.98 | 319 |
<!-- markdownlint-enable line-length -->

### Metrics

<!-- markdownlint-disable line-length -->
| Name | Type, Unit | Description |
| ---- | ---------- | ----------- |
| apiserver_audit_event | counter | ALPHA Counter of audit events generated and sent to the audit backend. |
| apiserver_audit_requests_rejected | counter | ALPHA Counter of apiserver requests rejected due to an error in audit logging backend. |
| apiserver_client_certificate_expiration_seconds | histogram | ALPHA Distribution of the remaining lifetime on the certificate used to authenticate a request. |
| apiserver_envelope_encryption_dek_cache_fill_percent | gauge | ALPHA Percent of the cache slots currently occupied by cached DEKs. |
| apiserver_storage_data_key_generation_duration_seconds | histogram | ALPHA Latencies in seconds of data encryption key(DEK) generation operations. |
| apiserver_storage_data_key_generation_failures | counter | ALPHA Total number of failed data encryption key(DEK) generation operations. |
| apiserver_storage_envelope_transformation_cache_misses | counter | ALPHA Total number of cache misses while accessing key decryption key(KEK). |
| authentication_token_cache_active_fetch_count | gauge | ALPHA |
| authentication_token_cache_fetch | counter | ALPHA |
| authentication_token_cache_request_duration_seconds | histogram | ALPHA |
| authentication_token_cache_request | counter | ALPHA |
| get_token_count | counter | ALPHA Counter of total Token() requests to the alternate token source |
| get_token_fail_count | counter | ALPHA Counter of failed Token() requests to the alternate token source |
| go_gc_duration_seconds | summary | A summary of the pause duration of garbage collection cycles. |
| go_goroutines | gauge | Number of goroutines that currently exist. |
| go_info | gauge | Information about the Go environment. |
| go_memstats_alloc_bytes | gauge | Number of bytes allocated and still in use. |
| go_memstats_alloc_bytes | counter | Total number of bytes allocated, even if freed. |
| go_memstats_buck_hash_sys_bytes | gauge | Number of bytes used by the profiling bucket hash table. |
| go_memstats_frees | counter | Total number of frees. |
| go_memstats_gc_cpu_fraction | gauge | The fraction of this program's available CPU time used by the GC since the program started. |
| go_memstats_gc_sys_bytes | gauge | Number of bytes used for garbage collection system metadata. |
| go_memstats_heap_alloc_bytes | gauge | Number of heap bytes allocated and still in use. |
| go_memstats_heap_idle_bytes | gauge | Number of heap bytes waiting to be used. |
| go_memstats_heap_inuse_bytes | gauge | Number of heap bytes that are in use. |
| go_memstats_heap_objects | gauge | Number of allocated objects. |
| go_memstats_heap_released_bytes | gauge | Number of heap bytes released to OS. |
| go_memstats_heap_sys_bytes | gauge | Number of heap bytes obtained from system. |
| go_memstats_last_gc_time_seconds | gauge | Number of seconds since 1970 of last garbage collection. |
| go_memstats_lookups | counter | Total number of pointer lookups. |
| go_memstats_mallocs | counter | Total number of mallocs. |
| go_memstats_mcache_inuse_bytes | gauge | Number of bytes in use by mcache structures. |
| go_memstats_mcache_sys_bytes | gauge | Number of bytes used for mcache structures obtained from system. |
| go_memstats_mspan_inuse_bytes | gauge | Number of bytes in use by mspan structures. |
| go_memstats_mspan_sys_bytes | gauge | Number of bytes used for mspan structures obtained from system. |
| go_memstats_next_gc_bytes | gauge | Number of heap bytes when next garbage collection will take place. |
| go_memstats_other_sys_bytes | gauge | Number of bytes used for other system allocations. |
| go_memstats_stack_inuse_bytes | gauge | Number of bytes in use by the stack allocator. |
| go_memstats_stack_sys_bytes | gauge | Number of bytes obtained from system for stack allocator. |
| go_memstats_sys_bytes | gauge | Number of bytes obtained from system. |
| go_threads | gauge | Number of OS threads created. |
| kubelet_certificate_manager_client_expiration_renew_errors | counter | ALPHA Counter of certificate renewal errors. |
| kubelet_certificate_manager_client_ttl_seconds | gauge | ALPHA Gauge of the TTL (time-to-live) of the Kubelet's client certificate. The value is in seconds until certificate expiry (negative if already expired). If client certificate is invalid or unused, the value will be +INF. |
| kubelet_cgroup_manager_duration_seconds | histogram | ALPHA Duration in seconds for cgroup manager operations. Broken down by method. |
| kubelet_container_log_filesystem_used_bytes | gauge | ALPHA Bytes used by the container's logs on the filesystem. |
| kubelet_containers_per_pod_count | histogram | ALPHA The number of containers per pod. |
| kubelet_eviction_stats_age_seconds | histogram | ALPHA Time between when stats are collected, and when pod is evicted based on those stats by eviction signal |
| kubelet_evictions | counter | ALPHA Cumulative number of pod evictions by eviction signal |
| kubelet_http_inflight_requests | gauge | ALPHA Number of the inflight http requests |
| kubelet_http_requests_duration_seconds | histogram | ALPHA Duration in seconds to serve http requests |
| kubelet_http_requests | counter | ALPHA Number of the http requests received since the server started |
| kubelet_node_config_error | gauge | ALPHA This metric is true (1) if the node is experiencing a configuration-related error, false (0) otherwise. |
| kubelet_node_name | gauge | ALPHA The node's name. The count is always 1. |
| kubelet_pleg_discard_events | counter | ALPHA The number of discard events in PLEG. |
| kubelet_pleg_last_seen_seconds | gauge | ALPHA Timestamp in seconds when PLEG was last seen active. |
| kubelet_pleg_relist_duration_seconds | histogram | ALPHA Duration in seconds for relisting pods in PLEG. |
| kubelet_pleg_relist_interval_seconds | histogram | ALPHA Interval in seconds between relisting in PLEG. |
| kubelet_pod_start_duration_seconds | histogram | ALPHA Duration in seconds for a single pod to go from pending to running. |
| kubelet_pod_worker_duration_seconds | histogram | ALPHA Duration in seconds to sync a single pod. Broken down by operation type: create, update, or sync |
| kubelet_pod_worker_start_duration_seconds | histogram | ALPHA Duration in seconds from seeing a pod to starting a worker. |
| kubelet_run_podsandbox_duration_seconds | histogram | ALPHA Duration in seconds of the run_podsandbox operations. Broken down by RuntimeClass.Handler. |
| kubelet_run_podsandbox_errors | counter | ALPHA Cumulative number of the run_podsandbox operation errors by RuntimeClass.Handler. |
| kubelet_running_containers | gauge | ALPHA Number of containers currently running |
| kubelet_running_pods | gauge | ALPHA Number of pods currently running |
| kubelet_runtime_operations_duration_seconds | histogram | ALPHA Duration in seconds of runtime operations. Broken down by operation type. |
| kubelet_runtime_operations_errors | counter | ALPHA Cumulative number of runtime operation errors by operation type. |
| kubelet_runtime_operations | counter | ALPHA Cumulative number of runtime operations by operation type. |
| kubernetes_build_info | gauge | ALPHA A metric with a constant '1' value labeled by major, minor, git version, git commit, git tree state, build date, Go version, and compiler from which Kubernetes was built, and platform on which it is running. |
| process_cpu_seconds | counter | Total user and system CPU time spent in seconds. |
| process_max_fds | gauge | Maximum number of open file descriptors. |
| process_open_fds | gauge | Number of open file descriptors. |
| process_resident_memory_bytes | gauge | Resident memory size in bytes. |
| process_start_time_seconds | gauge | Start time of the process since unix epoch in seconds. |
| process_virtual_memory_bytes | gauge | Virtual memory size in bytes. |
| process_virtual_memory_max_bytes | gauge | Maximum amount of virtual memory available in bytes. |
| rest_client_exec_plugin_certificate_rotation_age | histogram | ALPHA Histogram of the number of seconds the last auth exec plugin client certificate lived before being rotated. If auth exec plugin client certificates are unused, histogram will contain no data. |
| rest_client_exec_plugin_ttl_seconds | gauge | ALPHA Gauge of the shortest TTL (time-to-live) of the client certificate(s) managed by the auth exec plugin. The value is in seconds until certificate expiry (negative if already expired). If auth exec plugins are unused or manage no TLS certificates, the value will be +INF. |
| rest_client_request_duration_seconds | histogram | ALPHA Request latency in seconds. Broken down by verb and URL. |
| rest_client_requests | counter | ALPHA Number of HTTP requests, partitioned by status code, method, and host. |
| storage_operation_duration_seconds | histogram | ALPHA Storage operation duration |
| storage_operation_errors | counter | ALPHA Storage operation errors (Deprecated since 1.21.0) |
| storage_operation_status_count | counter | ALPHA Storage operation return statuses count (Deprecated since 1.21.0) |
| volume_manager_total_volumes | gauge | ALPHA Number of volumes in Volume Manager |
| workqueue_adds | counter | ALPHA Total number of adds handled by workqueue |
| workqueue_depth | gauge | ALPHA Current depth of workqueue |
| workqueue_longest_running_processor_seconds | gauge | ALPHA How many seconds has the longest running processor for workqueue been running. |
| workqueue_queue_duration_seconds | histogram | ALPHA How long in seconds an item stays in workqueue before being requested. |
| workqueue_retries | counter | ALPHA Total number of retries handled by workqueue |
| workqueue_unfinished_work_seconds | gauge | ALPHA How many seconds of work has done that is in progress and hasn't been observed by work_duration. Large values indicate stuck threads. One can deduce the number of stuck threads by observing the rate at which this increases. |
| workqueue_work_duration_seconds | histogram | ALPHA How long in seconds processing an item from workqueue takes. |
<!-- markdownlint-enable line-length -->

## cAdvisor (Kubernetes v1.21)

### Statistics

<!-- markdownlint-disable line-length -->
| Number of metrics families | Number of time series | Average number of labels per time series | Number of unique label pairs |
| -------------------------- | --------------------- | ---------------------------------------- | ---------------------------- |
| 66 | 5975 | 6.86 | 359 |
<!-- markdownlint-enable line-length -->

### Metrics

<!-- markdownlint-disable line-length -->
| Name | Type, Unit | Description |
| ---- | ---------- | ----------- |
| cadvisor_version_info | gauge | A metric with a constant '1' value labeled by kernel version, OS version, docker version, cadvisor version & cadvisor revision. |
| container_blkio_device_usage | counter | Blkio Device bytes usage |
| container_cpu_cfs_periods | counter | Number of elapsed enforcement period intervals. |
| container_cpu_cfs_throttled_periods | counter | Number of throttled period intervals. |
| container_cpu_cfs_throttled_seconds | counter | Total time duration the container has been throttled. |
| container_cpu_load_average_10s | gauge | Value of container cpu load average over the last 10 seconds. |
| container_cpu_system_seconds | counter | Cumulative system cpu time consumed in seconds. |
| container_cpu_usage_seconds | counter | Cumulative cpu time consumed in seconds. |
| container_cpu_user_seconds | counter | Cumulative user cpu time consumed in seconds. |
| container_file_descriptors | gauge | Number of open file descriptors for the container. |
| container_fs_inodes_free | gauge | Number of available Inodes |
| container_fs_inodes_total | gauge | Number of Inodes |
| container_fs_io_current | gauge | Number of I/Os currently in progress |
| container_fs_io_time_seconds | counter | Cumulative count of seconds spent doing I/Os |
| container_fs_io_time_weighted_seconds | counter | Cumulative weighted I/O time in seconds |
| container_fs_limit_bytes | gauge | Number of bytes that can be consumed by the container on this filesystem. |
| container_fs_read_seconds | counter | Cumulative count of seconds spent reading |
| container_fs_reads_bytes | counter | Cumulative count of bytes read |
| container_fs_reads_merged | counter | Cumulative count of reads merged |
| container_fs_reads | counter | Cumulative count of reads completed |
| container_fs_sector_reads | counter | Cumulative count of sector reads completed |
| container_fs_sector_writes | counter | Cumulative count of sector writes completed |
| container_fs_usage_bytes | gauge | Number of bytes that are consumed by the container on this filesystem. |
| container_fs_write_seconds | counter | Cumulative count of seconds spent writing |
| container_fs_writes_bytes | counter | Cumulative count of bytes written |
| container_fs_writes_merged | counter | Cumulative count of writes merged |
| container_fs_writes | counter | Cumulative count of writes completed |
| container_last_seen | gauge | Last time a container was seen by the exporter |
| container_memory_cache | gauge | Number of bytes of page cache memory. |
| container_memory_failcnt | counter | Number of memory usage hits limits |
| container_memory_failures | counter | Cumulative count of memory allocation failures. |
| container_memory_mapped_file | gauge | Size of memory mapped files in bytes. |
| container_memory_max_usage_bytes | gauge | Maximum memory usage recorded in bytes |
| container_memory_rss | gauge | Size of RSS in bytes. |
| container_memory_swap | gauge | Container swap usage in bytes. |
| container_memory_usage_bytes | gauge | Current memory usage in bytes, including all memory regardless of when it was accessed |
| container_memory_working_set_bytes | gauge | Current working set in bytes. |
| container_network_receive_bytes | counter | Cumulative count of bytes received |
| container_network_receive_errors | counter | Cumulative count of errors encountered while receiving |
| container_network_receive_packets_dropped | counter | Cumulative count of packets dropped while receiving |
| container_network_receive_packets | counter | Cumulative count of packets received |
| container_network_transmit_bytes | counter | Cumulative count of bytes transmitted |
| container_network_transmit_errors | counter | Cumulative count of errors encountered while transmitting |
| container_network_transmit_packets_dropped | counter | Cumulative count of packets dropped while transmitting |
| container_network_transmit_packets | counter | Cumulative count of packets transmitted |
| container_processes | gauge | Number of processes running inside the container. |
| container_scrape_error | gauge | 1 if there was an error while getting container metrics, 0 otherwise |
| container_sockets | gauge | Number of open sockets for the container. |
| container_spec_cpu_period | gauge | CPU period of the container. |
| container_spec_cpu_quota | gauge | CPU quota of the container. |
| container_spec_cpu_shares | gauge | CPU share of the container. |
| container_spec_memory_limit_bytes | gauge | Memory limit for the container. |
| container_spec_memory_reservation_limit_bytes | gauge | Memory reservation limit for the container. |
| container_spec_memory_swap_limit_bytes | gauge | Memory swap limit for the container. |
| container_start_time_seconds | gauge | Start time of the container since unix epoch in seconds. |
| container_tasks_state | gauge | Number of tasks in given state |
| container_threads | gauge | Number of threads running inside the container |
| container_threads_max | gauge | Maximum number of threads allowed inside the container, infinity if value is zero |
| container_ulimits_soft | gauge | Soft ulimit values for the container root process. Unlimited if -1, except priority and nice |
| machine_cpu_cores | gauge | Number of logical CPU cores. |
| machine_cpu_physical_cores | gauge | Number of physical CPU cores. |
| machine_cpu_sockets | gauge | Number of CPU sockets. |
| machine_memory_bytes | gauge | Amount of memory installed on the machine. |
| machine_nvm_avg_power_budget_watts | gauge | NVM power budget. |
| machine_nvm_capacity | gauge | NVM capacity value labeled by NVM mode (memory mode or app direct mode). |
| machine_scrape_error | gauge | 1 if there was an error while getting machine metrics, 0 otherwise. |
<!-- markdownlint-enable line-length -->

## NGINX ingress (Kubernetes v1.21)

### Statistics

<!-- markdownlint-disable line-length -->
| Number of metrics families | Number of time series | Average number of labels per time series | Number of unique label pairs |
| -------------------------- | --------------------- | ---------------------------------------- | ---------------------------- |
| 59 | 4597 | 11.6 | 126 |
<!-- markdownlint-enable line-length -->

### Metrics

<!-- markdownlint-disable line-length -->
| Name | Type, Unit | Description |
| ---- | ---------- | ----------- |
| go_gc_duration_seconds | summary | A summary of the pause duration of garbage collection cycles. |
| go_goroutines | gauge | Number of goroutines that currently exist. |
| go_info | gauge | Information about the Go environment. |
| go_memstats_alloc_bytes | gauge | Number of bytes allocated and still in use. |
| go_memstats_alloc_bytes | counter | Total number of bytes allocated, even if freed. |
| go_memstats_buck_hash_sys_bytes | gauge | Number of bytes used by the profiling bucket hash table. |
| go_memstats_frees | counter | Total number of frees. |
| go_memstats_gc_cpu_fraction | gauge | The fraction of this program's available CPU time used by the GC since the program started. |
| go_memstats_gc_sys_bytes | gauge | Number of bytes used for garbage collection system metadata. |
| go_memstats_heap_alloc_bytes | gauge | Number of heap bytes allocated and still in use. |
| go_memstats_heap_idle_bytes | gauge | Number of heap bytes waiting to be used. |
| go_memstats_heap_inuse_bytes | gauge | Number of heap bytes that are in use. |
| go_memstats_heap_objects | gauge | Number of allocated objects. |
| go_memstats_heap_released_bytes | gauge | Number of heap bytes released to OS. |
| go_memstats_heap_sys_bytes | gauge | Number of heap bytes obtained from system. |
| go_memstats_last_gc_time_seconds | gauge | Number of seconds since 1970 of last garbage collection. |
| go_memstats_lookups | counter | Total number of pointer lookups. |
| go_memstats_mallocs | counter | Total number of mallocs. |
| go_memstats_mcache_inuse_bytes | gauge | Number of bytes in use by mcache structures. |
| go_memstats_mcache_sys_bytes | gauge | Number of bytes used for mcache structures obtained from system. |
| go_memstats_mspan_inuse_bytes | gauge | Number of bytes in use by mspan structures. |
| go_memstats_mspan_sys_bytes | gauge | Number of bytes used for mspan structures obtained from system. |
| go_memstats_next_gc_bytes | gauge | Number of heap bytes when next garbage collection will take place. |
| go_memstats_other_sys_bytes | gauge | Number of bytes used for other system allocations. |
| go_memstats_stack_inuse_bytes | gauge | Number of bytes in use by the stack allocator. |
| go_memstats_stack_sys_bytes | gauge | Number of bytes obtained from system for stack allocator. |
| go_memstats_sys_bytes | gauge | Number of bytes obtained from system. |
| go_threads | gauge | Number of OS threads created. |
| nginx_ingress_controller_build_info | gauge | A metric with a constant '1' labeled with information about the build. |
| nginx_ingress_controller_bytes_sent | histogram | The number of bytes sent to a client |
| nginx_ingress_controller_config_hash | gauge | Running configuration hash actually running |
| nginx_ingress_controller_config_last_reload_successful | gauge | Whether the last configuration reload attempt was successful |
| nginx_ingress_controller_config_last_reload_successful_timestamp_seconds | gauge | Timestamp of the last successful configuration reload. |
| nginx_ingress_controller_ingress_upstream_latency_seconds | summary | Upstream service latency per Ingress |
| nginx_ingress_controller_nginx_process_connections | gauge | current number of client connections with state {active, reading, writing, waiting} |
| nginx_ingress_controller_nginx_process_connections | counter | total number of connections with state {accepted, handled} |
| nginx_ingress_controller_nginx_process_cpu_seconds | counter | Cpu usage in seconds |
| nginx_ingress_controller_nginx_process_num_procs | gauge | number of processes |
| nginx_ingress_controller_nginx_process_oldest_start_time_seconds | gauge | start time in seconds since 1970/01/01 |
| nginx_ingress_controller_nginx_process_read_bytes | counter | number of bytes read |
| nginx_ingress_controller_nginx_process_requests | counter | total number of client requests |
| nginx_ingress_controller_nginx_process_resident_memory_bytes | gauge | number of bytes of memory in use |
| nginx_ingress_controller_nginx_process_virtual_memory_bytes | gauge | number of bytes of memory in use |
| nginx_ingress_controller_nginx_process_write_bytes | counter | number of bytes written |
| nginx_ingress_controller_request_duration_seconds | histogram | The request processing time in milliseconds |
| nginx_ingress_controller_request_size | histogram | The request length (including request line, header, and request body) |
| nginx_ingress_controller_requests | counter | The total number of client requests. |
| nginx_ingress_controller_response_duration_seconds | histogram | The time spent on receiving the response from the upstream server |
| nginx_ingress_controller_response_size | histogram | The response length (including request line, header, and request body) |
| nginx_ingress_controller_success | counter | Cumulative number of Ingress controller reload operations |
| process_cpu_seconds | counter | Total user and system CPU time spent in seconds. |
| process_max_fds | gauge | Maximum number of open file descriptors. |
| process_open_fds | gauge | Number of open file descriptors. |
| process_resident_memory_bytes | gauge | Resident memory size in bytes. |
| process_start_time_seconds | gauge | Start time of the process since unix epoch in seconds. |
| process_virtual_memory_bytes | gauge | Virtual memory size in bytes. |
| process_virtual_memory_max_bytes | gauge | Maximum amount of virtual memory available in bytes. |
| promhttp_metric_handler_requests_in_flight | gauge | Current number of scrapes being served. |
| promhttp_metric_handler_requests | counter | Total number of scrapes by HTTP status code. |
<!-- markdownlint-enable line-length -->

## Node-exporter (Kubernetes v1.21)

### Statistics

<!-- markdownlint-disable line-length -->
| Number of metrics families | Number of time series | Average number of labels per time series | Number of unique label pairs |
| -------------------------- | --------------------- | ---------------------------------------- | ---------------------------- |
| 321 | 2692 | 1.42 | 416 |
<!-- markdownlint-enable line-length -->

### Metrics

<!-- markdownlint-disable line-length -->
| Name | Type, Unit | Description |
| ---- | ---------- | ----------- |
| go_gc_duration_seconds | summary | A summary of the pause duration of garbage collection cycles. |
| go_goroutines | gauge | Number of goroutines that currently exist. |
| go_info | gauge | Information about the Go environment. |
| go_memstats_alloc_bytes | gauge | Number of bytes allocated and still in use. |
| go_memstats_alloc_bytes | counter | Total number of bytes allocated, even if freed. |
| go_memstats_buck_hash_sys_bytes | gauge | Number of bytes used by the profiling bucket hash table. |
| go_memstats_frees | counter | Total number of frees. |
| go_memstats_gc_cpu_fraction | gauge | The fraction of this program's available CPU time used by the GC since the program started. |
| go_memstats_gc_sys_bytes | gauge | Number of bytes used for garbage collection system metadata. |
| go_memstats_heap_alloc_bytes | gauge | Number of heap bytes allocated and still in use. |
| go_memstats_heap_idle_bytes | gauge | Number of heap bytes waiting to be used. |
| go_memstats_heap_inuse_bytes | gauge | Number of heap bytes that are in use. |
| go_memstats_heap_objects | gauge | Number of allocated objects. |
| go_memstats_heap_released_bytes | gauge | Number of heap bytes released to OS. |
| go_memstats_heap_sys_bytes | gauge | Number of heap bytes obtained from system. |
| go_memstats_last_gc_time_seconds | gauge | Number of seconds since 1970 of last garbage collection. |
| go_memstats_lookups | counter | Total number of pointer lookups. |
| go_memstats_mallocs | counter | Total number of mallocs. |
| go_memstats_mcache_inuse_bytes | gauge | Number of bytes in use by mcache structures. |
| go_memstats_mcache_sys_bytes | gauge | Number of bytes used for mcache structures obtained from system. |
| go_memstats_mspan_inuse_bytes | gauge | Number of bytes in use by mspan structures. |
| go_memstats_mspan_sys_bytes | gauge | Number of bytes used for mspan structures obtained from system. |
| go_memstats_next_gc_bytes | gauge | Number of heap bytes when next garbage collection will take place. |
| go_memstats_other_sys_bytes | gauge | Number of bytes used for other system allocations. |
| go_memstats_stack_inuse_bytes | gauge | Number of bytes in use by the stack allocator. |
| go_memstats_stack_sys_bytes | gauge | Number of bytes obtained from system for stack allocator. |
| go_memstats_sys_bytes | gauge | Number of bytes obtained from system. |
| go_threads | gauge | Number of OS threads created. |
| node_arp_entries | gauge | ARP entries by device |
| node_boot_time_seconds | gauge | Node boot time, in unixtime. |
| node_context_switches | counter | Total number of context switches. |
| node_cooling_device_cur_state | gauge | Current throttle state of the cooling device |
| node_cooling_device_max_state | gauge | Maximum throttle state of the cooling device |
| node_cpu_guest_seconds | counter | Seconds the CPUs spent in guests (VMs) for each mode. |
| node_cpu_seconds | counter | Seconds the CPUs spent in each mode. |
| node_disk_info | gauge | Info of /sys/block/<block_device>. |
| node_disk_io_now | gauge | The number of I/Os currently in progress. |
| node_disk_io_time_seconds | counter | Total seconds spent doing I/Os. |
| node_disk_io_time_weighted_seconds | counter | The weighted # of seconds spent doing I/Os. |
| node_disk_read_bytes | counter | The total number of bytes read successfully. |
| node_disk_read_time_seconds | counter | The total number of seconds spent by all reads. |
| node_disk_reads_completed | counter | The total number of reads completed successfully. |
| node_disk_reads_merged | counter | The total number of reads merged. |
| node_disk_write_time_seconds | counter | This is the total number of seconds spent by all writes. |
| node_disk_writes_completed | counter | The total number of writes completed successfully. |
| node_disk_writes_merged | counter | The number of writes merged. |
| node_disk_written_bytes | counter | The total number of bytes written successfully. |
| node_dmi_info | gauge | A metric with a constant '1' value labeled by bios_date, bios_release, bios_vendor, bios_version, board_asset_tag, board_name, board_serial, board_vendor, board_version, chassis_asset_tag, chassis_serial, chassis_vendor, chassis_version, product_family, product_name, product_serial, product_sku, product_uuid, product_version, system_vendor if provided by DMI. |
| node_entropy_available_bits | gauge | Bits of available entropy. |
| node_entropy_pool_size_bits | gauge | Bits of entropy pool. |
| node_exporter_build_info | gauge | A metric with a constant '1' value labeled by version, revision, branch, and goversion from which node_exporter was built. |
| node_filefd_allocated | gauge | File descriptor statistics: allocated. |
| node_filefd_maximum | gauge | File descriptor statistics: maximum. |
| node_filesystem_avail_bytes | gauge | Filesystem space available to non-root users in bytes. |
| node_filesystem_device_error | gauge | Whether an error occurred while getting statistics for the given device. |
| node_filesystem_files | gauge | Filesystem total file nodes. |
| node_filesystem_files_free | gauge | Filesystem total free file nodes. |
| node_filesystem_free_bytes | gauge | Filesystem free space in bytes. |
| node_filesystem_readonly | gauge | Filesystem read-only status. |
| node_filesystem_size_bytes | gauge | Filesystem size in bytes. |
| node_forks | counter | Total number of forks. |
| node_intr | counter | Total number of interrupts serviced. |
| node_ipvs_connections | counter | The total number of connections made. |
| node_ipvs_incoming_bytes | counter | The total amount of incoming data. |
| node_ipvs_incoming_packets | counter | The total number of incoming packets. |
| node_ipvs_outgoing_bytes | counter | The total amount of outgoing data. |
| node_ipvs_outgoing_packets | counter | The total number of outgoing packets. |
| node_load1 | gauge | 1m load average. |
| node_load15 | gauge | 15m load average. |
| node_load5 | gauge | 5m load average. |
| node_memory_Active_anon_bytes | gauge | Memory information field Active_anon_bytes. |
| node_memory_Active_bytes | gauge | Memory information field Active_bytes. |
| node_memory_Active_file_bytes | gauge | Memory information field Active_file_bytes. |
| node_memory_AnonHugePages_bytes | gauge | Memory information field AnonHugePages_bytes. |
| node_memory_AnonPages_bytes | gauge | Memory information field AnonPages_bytes. |
| node_memory_Bounce_bytes | gauge | Memory information field Bounce_bytes. |
| node_memory_Buffers_bytes | gauge | Memory information field Buffers_bytes. |
| node_memory_Cached_bytes | gauge | Memory information field Cached_bytes. |
| node_memory_CmaFree_bytes | gauge | Memory information field CmaFree_bytes. |
| node_memory_CmaTotal_bytes | gauge | Memory information field CmaTotal_bytes. |
| node_memory_CommitLimit_bytes | gauge | Memory information field CommitLimit_bytes. |
| node_memory_Committed_AS_bytes | gauge | Memory information field Committed_AS_bytes. |
| node_memory_DirectMap2M_bytes | gauge | Memory information field DirectMap2M_bytes. |
| node_memory_DirectMap4k_bytes | gauge | Memory information field DirectMap4k_bytes. |
| node_memory_Dirty_bytes | gauge | Memory information field Dirty_bytes. |
| node_memory_HardwareCorrupted_bytes | gauge | Memory information field HardwareCorrupted_bytes. |
| node_memory_HugePages_Free | gauge | Memory information field HugePages_Free. |
| node_memory_HugePages_Rsvd | gauge | Memory information field HugePages_Rsvd. |
| node_memory_HugePages_Surp | gauge | Memory information field HugePages_Surp. |
| node_memory_HugePages_Total | gauge | Memory information field HugePages_Total. |
| node_memory_Hugepagesize_bytes | gauge | Memory information field Hugepagesize_bytes. |
| node_memory_Inactive_anon_bytes | gauge | Memory information field Inactive_anon_bytes. |
| node_memory_Inactive_bytes | gauge | Memory information field Inactive_bytes. |
| node_memory_Inactive_file_bytes | gauge | Memory information field Inactive_file_bytes. |
| node_memory_KernelStack_bytes | gauge | Memory information field KernelStack_bytes. |
| node_memory_Mapped_bytes | gauge | Memory information field Mapped_bytes. |
| node_memory_MemAvailable_bytes | gauge | Memory information field MemAvailable_bytes. |
| node_memory_MemFree_bytes | gauge | Memory information field MemFree_bytes. |
| node_memory_MemTotal_bytes | gauge | Memory information field MemTotal_bytes. |
| node_memory_Mlocked_bytes | gauge | Memory information field Mlocked_bytes. |
| node_memory_NFS_Unstable_bytes | gauge | Memory information field NFS_Unstable_bytes. |
| node_memory_PageTables_bytes | gauge | Memory information field PageTables_bytes. |
| node_memory_Percpu_bytes | gauge | Memory information field Percpu_bytes. |
| node_memory_SReclaimable_bytes | gauge | Memory information field SReclaimable_bytes. |
| node_memory_SUnreclaim_bytes | gauge | Memory information field SUnreclaim_bytes. |
| node_memory_Shmem_bytes | gauge | Memory information field Shmem_bytes. |
| node_memory_Slab_bytes | gauge | Memory information field Slab_bytes. |
| node_memory_SwapCached_bytes | gauge | Memory information field SwapCached_bytes. |
| node_memory_SwapFree_bytes | gauge | Memory information field SwapFree_bytes. |
| node_memory_SwapTotal_bytes | gauge | Memory information field SwapTotal_bytes. |
| node_memory_Unevictable_bytes | gauge | Memory information field Unevictable_bytes. |
| node_memory_VmallocChunk_bytes | gauge | Memory information field VmallocChunk_bytes. |
| node_memory_VmallocTotal_bytes | gauge | Memory information field VmallocTotal_bytes. |
| node_memory_VmallocUsed_bytes | gauge | Memory information field VmallocUsed_bytes. |
| node_memory_WritebackTmp_bytes | gauge | Memory information field WritebackTmp_bytes. |
| node_memory_Writeback_bytes | gauge | Memory information field Writeback_bytes. |
| node_netstat_Icmp6_InErrors | unknown | Statistic Icmp6InErrors. |
| node_netstat_Icmp6_InMsgs | unknown | Statistic Icmp6InMsgs. |
| node_netstat_Icmp6_OutMsgs | unknown | Statistic Icmp6OutMsgs. |
| node_netstat_Icmp_InErrors | unknown | Statistic IcmpInErrors. |
| node_netstat_Icmp_InMsgs | unknown | Statistic IcmpInMsgs. |
| node_netstat_Icmp_OutMsgs | unknown | Statistic IcmpOutMsgs. |
| node_netstat_Ip6_InOctets | unknown | Statistic Ip6InOctets. |
| node_netstat_Ip6_OutOctets | unknown | Statistic Ip6OutOctets. |
| node_netstat_IpExt_InOctets | unknown | Statistic IpExtInOctets. |
| node_netstat_IpExt_OutOctets | unknown | Statistic IpExtOutOctets. |
| node_netstat_Ip_Forwarding | unknown | Statistic IpForwarding. |
| node_netstat_TcpExt_ListenDrops | unknown | Statistic TcpExtListenDrops. |
| node_netstat_TcpExt_ListenOverflows | unknown | Statistic TcpExtListenOverflows. |
| node_netstat_TcpExt_SyncookiesFailed | unknown | Statistic TcpExtSyncookiesFailed. |
| node_netstat_TcpExt_SyncookiesRecv | unknown | Statistic TcpExtSyncookiesRecv. |
| node_netstat_TcpExt_SyncookiesSent | unknown | Statistic TcpExtSyncookiesSent. |
| node_netstat_TcpExt_TCPSynRetrans | unknown | Statistic TcpExtTCPSynRetrans. |
| node_netstat_TcpExt_TCPTimeouts | unknown | Statistic TcpExtTCPTimeouts. |
| node_netstat_Tcp_ActiveOpens | unknown | Statistic TcpActiveOpens. |
| node_netstat_Tcp_CurrEstab | unknown | Statistic TcpCurrEstab. |
| node_netstat_Tcp_InErrs | unknown | Statistic TcpInErrs. |
| node_netstat_Tcp_InSegs | unknown | Statistic TcpInSegs. |
| node_netstat_Tcp_OutRsts | unknown | Statistic TcpOutRsts. |
| node_netstat_Tcp_OutSegs | unknown | Statistic TcpOutSegs. |
| node_netstat_Tcp_PassiveOpens | unknown | Statistic TcpPassiveOpens. |
| node_netstat_Tcp_RetransSegs | unknown | Statistic TcpRetransSegs. |
| node_netstat_Udp6_InDatagrams | unknown | Statistic Udp6InDatagrams. |
| node_netstat_Udp6_InErrors | unknown | Statistic Udp6InErrors. |
| node_netstat_Udp6_NoPorts | unknown | Statistic Udp6NoPorts. |
| node_netstat_Udp6_OutDatagrams | unknown | Statistic Udp6OutDatagrams. |
| node_netstat_Udp6_RcvbufErrors | unknown | Statistic Udp6RcvbufErrors. |
| node_netstat_Udp6_SndbufErrors | unknown | Statistic Udp6SndbufErrors. |
| node_netstat_UdpLite6_InErrors | unknown | Statistic UdpLite6InErrors. |
| node_netstat_UdpLite_InErrors | unknown | Statistic UdpLiteInErrors. |
| node_netstat_Udp_InDatagrams | unknown | Statistic UdpInDatagrams. |
| node_netstat_Udp_InErrors | unknown | Statistic UdpInErrors. |
| node_netstat_Udp_NoPorts | unknown | Statistic UdpNoPorts. |
| node_netstat_Udp_OutDatagrams | unknown | Statistic UdpOutDatagrams. |
| node_netstat_Udp_RcvbufErrors | unknown | Statistic UdpRcvbufErrors. |
| node_netstat_Udp_SndbufErrors | unknown | Statistic UdpSndbufErrors. |
| node_network_address_assign_type | gauge | address_assign_type value of /sys/class/net/<iface>. |
| node_network_carrier | gauge | carrier value of /sys/class/net/<iface>. |
| node_network_carrier_changes | counter | carrier_changes_total value of /sys/class/net/<iface>. |
| node_network_device_id | gauge | device_id value of /sys/class/net/<iface>. |
| node_network_dormant | gauge | dormant value of /sys/class/net/<iface>. |
| node_network_flags | gauge | flags value of /sys/class/net/<iface>. |
| node_network_iface_id | gauge | iface_id value of /sys/class/net/<iface>. |
| node_network_iface_link | gauge | iface_link value of /sys/class/net/<iface>. |
| node_network_iface_link_mode | gauge | iface_link_mode value of /sys/class/net/<iface>. |
| node_network_info | gauge | Non-numeric data from /sys/class/net/<iface>, value is always 1. |
| node_network_mtu_bytes | gauge | mtu_bytes value of /sys/class/net/<iface>. |
| node_network_net_dev_group | gauge | net_dev_group value of /sys/class/net/<iface>. |
| node_network_protocol_type | gauge | protocol_type value of /sys/class/net/<iface>. |
| node_network_receive_bytes | counter | Network device statistic receive_bytes. |
| node_network_receive_compressed | counter | Network device statistic receive_compressed. |
| node_network_receive_drop | counter | Network device statistic receive_drop. |
| node_network_receive_errs | counter | Network device statistic receive_errs. |
| node_network_receive_fifo | counter | Network device statistic receive_fifo. |
| node_network_receive_frame | counter | Network device statistic receive_frame. |
| node_network_receive_multicast | counter | Network device statistic receive_multicast. |
| node_network_receive_packets | counter | Network device statistic receive_packets. |
| node_network_speed_bytes | gauge | speed_bytes value of /sys/class/net/<iface>. |
| node_network_transmit_bytes | counter | Network device statistic transmit_bytes. |
| node_network_transmit_carrier | counter | Network device statistic transmit_carrier. |
| node_network_transmit_colls | counter | Network device statistic transmit_colls. |
| node_network_transmit_compressed | counter | Network device statistic transmit_compressed. |
| node_network_transmit_drop | counter | Network device statistic transmit_drop. |
| node_network_transmit_errs | counter | Network device statistic transmit_errs. |
| node_network_transmit_fifo | counter | Network device statistic transmit_fifo. |
| node_network_transmit_packets | counter | Network device statistic transmit_packets. |
| node_network_transmit_queue_length | gauge | transmit_queue_length value of /sys/class/net/<iface>. |
| node_network_up | gauge | Value is 1 if operstate is 'up', 0 otherwise. |
| node_nf_conntrack_entries | gauge | Number of currently allocated flow entries for connection tracking. |
| node_nf_conntrack_entries_limit | gauge | Maximum size of connection tracking table. |
| node_nf_conntrack_stat_drop | gauge | Number of packets dropped due to conntrack failure. |
| node_nf_conntrack_stat_early_drop | gauge | Number of dropped conntrack entries to make room for new ones, if maximum table size was reached. |
| node_nf_conntrack_stat_found | gauge | Number of searched entries which were successful. |
| node_nf_conntrack_stat_ignore | gauge | Number of packets seen which are already connected to a conntrack entry. |
| node_nf_conntrack_stat_insert | gauge | Number of entries inserted into the list. |
| node_nf_conntrack_stat_insert_failed | gauge | Number of entries for which list insertion was attempted but failed. |
| node_nf_conntrack_stat_invalid | gauge | Number of packets seen which can not be tracked. |
| node_nf_conntrack_stat_search_restart | gauge | Number of conntrack table lookups which had to be restarted due to hashtable resizes. |
| node_nfs_connections | counter | Total number of NFSd TCP connections. |
| node_nfs_packets | counter | Total NFSd network packets (sent+received) by protocol type. |
| node_nfs_requests | counter | Number of NFS procedures invoked. |
| node_nfs_rpc_authentication_refreshes | counter | Number of RPC authentication refreshes performed. |
| node_nfs_rpc_retransmissions | counter | Number of RPC transmissions performed. |
| node_nfs_rpcs | counter | Total number of RPCs performed. |
| node_os_info | gauge | A metric with a constant '1' value labeled by build_id, id, id_like, image_id, image_version, name, pretty_name, variant, variant_id, version, version_codename, version_id. |
| node_os_version | gauge | Metric containing the major.minor part of the OS version. |
| node_processes_max_processes | gauge | Number of max PIDs limit |
| node_processes_max_threads | gauge | Limit of threads in the system |
| node_processes_pids | gauge | Number of PIDs |
| node_processes_state | gauge | Number of processes in each state. |
| node_processes_threads | gauge | Allocated threads in system |
| node_processes_threads_state | gauge | Number of threads in each state. |
| node_procs_blocked | gauge | Number of processes blocked waiting for I/O to complete. |
| node_procs_running | gauge | Number of processes in runnable state. |
| node_schedstat_running_seconds | counter | Number of seconds CPU spent running a process. |
| node_schedstat_timeslices | counter | Number of timeslices executed by CPU. |
| node_schedstat_waiting_seconds | counter | Number of seconds spent by processing waiting for this CPU. |
| node_scrape_collector_duration_seconds | gauge | node_exporter: Duration of a collector scrape. |
| node_scrape_collector_success | gauge | node_exporter: Whether a collector succeeded. |
| node_sockstat_FRAG6_inuse | gauge | Number of FRAG6 sockets in state inuse. |
| node_sockstat_FRAG6_memory | gauge | Number of FRAG6 sockets in state memory. |
| node_sockstat_FRAG_inuse | gauge | Number of FRAG sockets in state inuse. |
| node_sockstat_FRAG_memory | gauge | Number of FRAG sockets in state memory. |
| node_sockstat_RAW6_inuse | gauge | Number of RAW6 sockets in state inuse. |
| node_sockstat_RAW_inuse | gauge | Number of RAW sockets in state inuse. |
| node_sockstat_TCP6_inuse | gauge | Number of TCP6 sockets in state inuse. |
| node_sockstat_TCP_alloc | gauge | Number of TCP sockets in state alloc. |
| node_sockstat_TCP_inuse | gauge | Number of TCP sockets in state inuse. |
| node_sockstat_TCP_mem | gauge | Number of TCP sockets in state mem. |
| node_sockstat_TCP_mem_bytes | gauge | Number of TCP sockets in state mem_bytes. |
| node_sockstat_TCP_orphan | gauge | Number of TCP sockets in state orphan. |
| node_sockstat_TCP_tw | gauge | Number of TCP sockets in state tw. |
| node_sockstat_UDP6_inuse | gauge | Number of UDP6 sockets in state inuse. |
| node_sockstat_UDPLITE6_inuse | gauge | Number of UDPLITE6 sockets in state inuse. |
| node_sockstat_UDPLITE_inuse | gauge | Number of UDPLITE sockets in state inuse. |
| node_sockstat_UDP_inuse | gauge | Number of UDP sockets in state inuse. |
| node_sockstat_UDP_mem | gauge | Number of UDP sockets in state mem. |
| node_sockstat_UDP_mem_bytes | gauge | Number of UDP sockets in state mem_bytes. |
| node_sockstat_sockets_used | gauge | Number of IPv4 sockets in use. |
| node_softnet_dropped | counter | Number of dropped packets |
| node_softnet_processed | counter | Number of processed packets |
| node_softnet_times_squeezed | counter | Number of times processing packets ran out of quota |
| node_textfile_scrape_error | gauge | 1 if there was an error opening or reading a file, 0 otherwise |
| node_time_clocksource_available_info | gauge | Available clocksources read from '/sys/devices/system/clocksource'. |
| node_time_clocksource_current_info | gauge | Current clocksource read from '/sys/devices/system/clocksource'. |
| node_time_seconds | gauge | System time in seconds since epoch (1970). |
| node_time_zone_offset_seconds | gauge | System time zone offset in seconds. |
| node_timex_estimated_error_seconds | gauge | Estimated error in seconds. |
| node_timex_frequency_adjustment_ratio | gauge | Local clock frequency adjustment. |
| node_timex_loop_time_constant | gauge | Phase-locked loop time constant. |
| node_timex_maxerror_seconds | gauge | Maximum error in seconds. |
| node_timex_offset_seconds | gauge | Time offset in between local system and reference clock. |
| node_timex_pps_calibration | counter | Pulse per second count of calibration intervals. |
| node_timex_pps_error | counter | Pulse per second count of calibration errors. |
| node_timex_pps_frequency_hertz | gauge | Pulse per second frequency. |
| node_timex_pps_jitter_seconds | gauge | Pulse per second jitter. |
| node_timex_pps_jitter | counter | Pulse per second count of jitter limit exceeded events. |
| node_timex_pps_shift_seconds | gauge | Pulse per second interval duration. |
| node_timex_pps_stability_exceeded | counter | Pulse per second count of stability limit exceeded events. |
| node_timex_pps_stability_hertz | gauge | Pulse per second stability, average of recent frequency changes. |
| node_timex_status | gauge | Value of the status array bits. |
| node_timex_sync_status | gauge | Is clock synchronized to a reliable server (1 = yes, 0 = no). |
| node_timex_tai_offset_seconds | gauge | International Atomic Time (TAI) offset. |
| node_timex_tick_seconds | gauge | Seconds between clock ticks. |
| node_udp_queues | gauge | Number of allocated memory in the kernel for UDP datagrams in bytes. |
| node_uname_info | gauge | Labeled system information as provided by the uname system call. |
| node_vmstat_pgfault | unknown | /proc/vmstat information field pgfault. |
| node_vmstat_pgmajfault | unknown | /proc/vmstat information field pgmajfault. |
| node_vmstat_pgpgin | unknown | /proc/vmstat information field pgpgin. |
| node_vmstat_pgpgout | unknown | /proc/vmstat information field pgpgout. |
| node_vmstat_pswpin | unknown | /proc/vmstat information field pswpin. |
| node_vmstat_pswpout | unknown | /proc/vmstat information field pswpout. |
| node_xfs_allocation_btree_compares | counter | Number of allocation B-tree compares for a filesystem. |
| node_xfs_allocation_btree_lookups | counter | Number of allocation B-tree lookups for a filesystem. |
| node_xfs_allocation_btree_records_deleted | counter | Number of allocation B-tree records deleted for a filesystem. |
| node_xfs_allocation_btree_records_inserted | counter | Number of allocation B-tree records inserted for a filesystem. |
| node_xfs_block_map_btree_compares | counter | Number of block map B-tree compares for a filesystem. |
| node_xfs_block_map_btree_lookups | counter | Number of block map B-tree lookups for a filesystem. |
| node_xfs_block_map_btree_records_deleted | counter | Number of block map B-tree records deleted for a filesystem. |
| node_xfs_block_map_btree_records_inserted | counter | Number of block map B-tree records inserted for a filesystem. |
| node_xfs_block_mapping_extent_list_compares | counter | Number of extent list compares for a filesystem. |
| node_xfs_block_mapping_extent_list_deletions | counter | Number of extent list deletions for a filesystem. |
| node_xfs_block_mapping_extent_list_insertions | counter | Number of extent list insertions for a filesystem. |
| node_xfs_block_mapping_extent_list_lookups | counter | Number of extent list lookups for a filesystem. |
| node_xfs_block_mapping_reads | counter | Number of block map for read operations for a filesystem. |
| node_xfs_block_mapping_unmaps | counter | Number of block unmaps (deletes) for a filesystem. |
| node_xfs_block_mapping_writes | counter | Number of block map for write operations for a filesystem. |
| node_xfs_directory_operation_create | counter | Number of times a new directory entry was created for a filesystem. |
| node_xfs_directory_operation_getdents | counter | Number of times the directory getdents operation was performed for a filesystem. |
| node_xfs_directory_operation_lookup | counter | Number of file name directory lookups which miss the operating systems directory name lookup cache. |
| node_xfs_directory_operation_remove | counter | Number of times an existing directory entry was created for a filesystem. |
| node_xfs_extent_allocation_blocks_allocated | counter | Number of blocks allocated for a filesystem. |
| node_xfs_extent_allocation_blocks_freed | counter | Number of blocks freed for a filesystem. |
| node_xfs_extent_allocation_extents_allocated | counter | Number of extents allocated for a filesystem. |
| node_xfs_extent_allocation_extents_freed | counter | Number of extents freed for a filesystem. |
| node_xfs_inode_operation_attempts | counter | Number of times the OS looked for an XFS inode in the inode cache. |
| node_xfs_inode_operation_attribute_changes | counter | Number of times the OS explicitly changed the attributes of an XFS inode. |
| node_xfs_inode_operation_duplicates | counter | Number of times the OS tried to add a missing XFS inode to the inode cache, but found it had already been added by another process. |
| node_xfs_inode_operation_found | counter | Number of times the OS looked for and found an XFS inode in the inode cache. |
| node_xfs_inode_operation_missed | counter | Number of times the OS looked for an XFS inode in the cache, but did not find it. |
| node_xfs_inode_operation_reclaims | counter | Number of times the OS reclaimed an XFS inode from the inode cache to free memory for another purpose. |
| node_xfs_inode_operation_recycled | counter | Number of times the OS found an XFS inode in the cache, but could not use it as it was being recycled. |
| node_xfs_read_calls | counter | Number of read(2) system calls made to files in a filesystem. |
| node_xfs_vnode_active | counter | Number of vnodes not on free lists for a filesystem. |
| node_xfs_vnode_allocate | counter | Number of times vn_alloc called for a filesystem. |
| node_xfs_vnode_get | counter | Number of times vn_get called for a filesystem. |
| node_xfs_vnode_hold | counter | Number of times vn_hold called for a filesystem. |
| node_xfs_vnode_reclaim | counter | Number of times vn_reclaim called for a filesystem. |
| node_xfs_vnode_release | counter | Number of times vn_rele called for a filesystem. |
| node_xfs_vnode_remove | counter | Number of times vn_remove called for a filesystem. |
| node_xfs_write_calls | counter | Number of write(2) system calls made to files in a filesystem. |
| process_cpu_seconds | counter | Total user and system CPU time spent in seconds. |
| process_max_fds | gauge | Maximum number of open file descriptors. |
| process_open_fds | gauge | Number of open file descriptors. |
| process_resident_memory_bytes | gauge | Resident memory size in bytes. |
| process_start_time_seconds | gauge | Start time of the process since unix epoch in seconds. |
| process_virtual_memory_bytes | gauge | Virtual memory size in bytes. |
| process_virtual_memory_max_bytes | gauge | Maximum amount of virtual memory available in bytes. |
| promhttp_metric_handler_errors | counter | Total number of internal errors encountered by the promhttp metric handler. |
| promhttp_metric_handler_requests_in_flight | gauge | Current number of scrapes being served. |
| promhttp_metric_handler_requests | counter | Total number of scrapes by HTTP status code. |
<!-- markdownlint-enable line-length -->

## Prometheus (Kubernetes v1.21)

### Statistics

<!-- markdownlint-disable line-length -->
| Number of metrics families | Number of time series | Average number of labels per time series | Number of unique label pairs |
| -------------------------- | --------------------- | ---------------------------------------- | ---------------------------- |
| 247 | 2718 | 1.39 | 429 |
<!-- markdownlint-enable line-length -->

### Metrics

<!-- markdownlint-disable line-length -->
| Name | Type, Unit | Description |
| ---- | ---------- | ----------- |
| go_gc_cycles_automatic_gc_cycles | counter | Count of completed GC cycles generated by the Go runtime. |
| go_gc_cycles_forced_gc_cycles | counter | Count of completed GC cycles forced by the application. |
| go_gc_cycles_total_gc_cycles | counter | Count of all completed GC cycles. |
| go_gc_duration_seconds | summary | A summary of the pause duration of garbage collection cycles. |
| go_gc_heap_allocs_by_size_bytes_total | histogram | Distribution of heap allocations by approximate size. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks. |
| go_gc_heap_allocs_bytes | counter | Cumulative sum of memory allocated to the heap by the application. |
| go_gc_heap_allocs_objects | counter | Cumulative count of heap allocations triggered by the application. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks. |
| go_gc_heap_frees_by_size_bytes_total | histogram | Distribution of freed heap allocations by approximate size. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks. |
| go_gc_heap_frees_bytes | counter | Cumulative sum of heap memory freed by the garbage collector. |
| go_gc_heap_frees_objects | counter | Cumulative count of heap allocations whose storage was freed by the garbage collector. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks. |
| go_gc_heap_goal_bytes | gauge | Heap size target for the end of the GC cycle. |
| go_gc_heap_objects_objects | gauge | Number of objects, live or unswept, occupying heap memory. |
| go_gc_heap_tiny_allocs_objects | counter | Count of small allocations that are packed together into blocks. These allocations are counted separately from other allocations because each individual allocation is not tracked by the runtime, only their block. Each block is already accounted for in allocs-by-size and frees-by-size. |
| go_gc_pauses_seconds_total | histogram | Distribution individual GC-related stop-the-world pause latencies. |
| go_goroutines | gauge | Number of goroutines that currently exist. |
| go_info | gauge | Information about the Go environment. |
| go_memory_classes_heap_free_bytes | gauge | Memory that is completely free and eligible to be returned to the underlying system, but has not been. This metric is the runtime's estimate of free address space that is backed by physical memory. |
| go_memory_classes_heap_objects_bytes | gauge | Memory occupied by live objects and dead objects that have not yet been marked free by the garbage collector. |
| go_memory_classes_heap_released_bytes | gauge | Memory that is completely free and has been returned to the underlying system. This metric is the runtime's estimate of free address space that is still mapped into the process, but is not backed by physical memory. |
| go_memory_classes_heap_stacks_bytes | gauge | Memory allocated from the heap that is reserved for stack space, whether or not it is currently in-use. |
| go_memory_classes_heap_unused_bytes | gauge | Memory that is reserved for heap objects but is not currently used to hold heap objects. |
| go_memory_classes_metadata_mcache_free_bytes | gauge | Memory that is reserved for runtime mcache structures, but not in-use. |
| go_memory_classes_metadata_mcache_inuse_bytes | gauge | Memory that is occupied by runtime mcache structures that are currently being used. |
| go_memory_classes_metadata_mspan_free_bytes | gauge | Memory that is reserved for runtime mspan structures, but not in-use. |
| go_memory_classes_metadata_mspan_inuse_bytes | gauge | Memory that is occupied by runtime mspan structures that are currently being used. |
| go_memory_classes_metadata_other_bytes | gauge | Memory that is reserved for or used to hold runtime metadata. |
| go_memory_classes_os_stacks_bytes | gauge | Stack memory allocated by the underlying operating system. |
| go_memory_classes_other_bytes | gauge | Memory used by execution trace buffers, structures for debugging the runtime, finalizer and profiler specials, and more. |
| go_memory_classes_profiling_buckets_bytes | gauge | Memory that is used by the stack trace hash map used for profiling. |
| go_memory_classes_total_bytes | gauge | All memory mapped by the Go runtime into the current process as read-write. Note that this does not include memory mapped by code called via cgo or via the syscall package. Sum of all metrics in /memory/classes. |
| go_memstats_alloc_bytes | gauge | Number of bytes allocated and still in use. |
| go_memstats_alloc_bytes | counter | Total number of bytes allocated, even if freed. |
| go_memstats_buck_hash_sys_bytes | gauge | Number of bytes used by the profiling bucket hash table. |
| go_memstats_frees | counter | Total number of frees. |
| go_memstats_gc_cpu_fraction | gauge | The fraction of this program's available CPU time used by the GC since the program started. |
| go_memstats_gc_sys_bytes | gauge | Number of bytes used for garbage collection system metadata. |
| go_memstats_heap_alloc_bytes | gauge | Number of heap bytes allocated and still in use. |
| go_memstats_heap_idle_bytes | gauge | Number of heap bytes waiting to be used. |
| go_memstats_heap_inuse_bytes | gauge | Number of heap bytes that are in use. |
| go_memstats_heap_objects | gauge | Number of allocated objects. |
| go_memstats_heap_released_bytes | gauge | Number of heap bytes released to OS. |
| go_memstats_heap_sys_bytes | gauge | Number of heap bytes obtained from system. |
| go_memstats_last_gc_time_seconds | gauge | Number of seconds since 1970 of last garbage collection. |
| go_memstats_lookups | counter | Total number of pointer lookups. |
| go_memstats_mallocs | counter | Total number of mallocs. |
| go_memstats_mcache_inuse_bytes | gauge | Number of bytes in use by mcache structures. |
| go_memstats_mcache_sys_bytes | gauge | Number of bytes used for mcache structures obtained from system. |
| go_memstats_mspan_inuse_bytes | gauge | Number of bytes in use by mspan structures. |
| go_memstats_mspan_sys_bytes | gauge | Number of bytes used for mspan structures obtained from system. |
| go_memstats_next_gc_bytes | gauge | Number of heap bytes when next garbage collection will take place. |
| go_memstats_other_sys_bytes | gauge | Number of bytes used for other system allocations. |
| go_memstats_stack_inuse_bytes | gauge | Number of bytes in use by the stack allocator. |
| go_memstats_stack_sys_bytes | gauge | Number of bytes obtained from system for stack allocator. |
| go_memstats_sys_bytes | gauge | Number of bytes obtained from system. |
| go_sched_goroutines_goroutines | gauge | Count of live goroutines. |
| go_sched_latencies_seconds | histogram | Distribution of the time goroutines have spent in the scheduler in a runnable state before actually running. |
| go_threads | gauge | Number of OS threads created. |
| net_conntrack_dialer_conn_attempted | counter | Total number of connections attempted by the given dialer a given name. |
| net_conntrack_dialer_conn_closed | counter | Total number of connections closed which originated from the dialer of a given name. |
| net_conntrack_dialer_conn_established | counter | Total number of connections successfully established by the given dialer a given name. |
| net_conntrack_dialer_conn_failed | counter | Total number of connections failed to dial by the dialer a given name. |
| net_conntrack_listener_conn_accepted | counter | Total number of connections opened to the listener of a given name. |
| net_conntrack_listener_conn_closed | counter | Total number of connections closed that were made to the listener of a given name. |
| process_cpu_seconds | counter | Total user and system CPU time spent in seconds. |
| process_max_fds | gauge | Maximum number of open file descriptors. |
| process_open_fds | gauge | Number of open file descriptors. |
| process_resident_memory_bytes | gauge | Resident memory size in bytes. |
| process_start_time_seconds | gauge | Start time of the process since unix epoch in seconds. |
| process_virtual_memory_bytes | gauge | Virtual memory size in bytes. |
| process_virtual_memory_max_bytes | gauge | Maximum amount of virtual memory available in bytes. |
| prometheus_api_remote_read_queries | gauge | The current number of remote read queries being executed or waiting. |
| prometheus_build_info | gauge | A metric with a constant '1' value labeled by version, revision, branch, and goversion from which prometheus was built. |
| prometheus_config_last_reload_success_timestamp_seconds | gauge | Timestamp of the last successful configuration reload. |
| prometheus_config_last_reload_successful | gauge | Whether the last configuration reload attempt was successful. |
| prometheus_engine_queries | gauge | The current number of queries being executed or waiting. |
| prometheus_engine_queries_concurrent_max | gauge | The max number of concurrent queries. |
| prometheus_engine_query_duration_seconds | summary | Query timings |
| prometheus_engine_query_log_enabled | gauge | State of the query log. |
| prometheus_engine_query_log_failures | counter | The number of query log failures. |
| prometheus_http_request_duration_seconds | histogram | Histogram of latencies for HTTP requests. |
| prometheus_http_requests | counter | Counter of HTTP requests. |
| prometheus_http_response_size_bytes | histogram | Histogram of response size for HTTP requests. |
| prometheus_notifications_alertmanagers_discovered | gauge | The number of alertmanagers discovered and active. |
| prometheus_notifications_dropped | counter | Total number of alerts dropped due to errors when sending to Alertmanager. |
| prometheus_notifications_errors | counter | Total number of errors sending alert notifications. |
| prometheus_notifications_latency_seconds | summary | Latency quantiles for sending alert notifications. |
| prometheus_notifications_queue_capacity | gauge | The capacity of the alert notifications queue. |
| prometheus_notifications_queue_length | gauge | The number of alert notifications in the queue. |
| prometheus_notifications_sent | counter | Total number of alerts sent. |
| prometheus_remote_storage_bytes | counter | The total number of bytes of data (not metadata) sent by the queue after compression. Note that when exemplars over remote write is enabled the exemplars included in a remote write request count towards this metric. |
| prometheus_remote_storage_enqueue_retries | counter | Total number of times enqueue has failed because a shards queue was full. |
| prometheus_remote_storage_exemplars_dropped | counter | Total number of exemplars which were dropped after being read from the WAL before being sent via remote write, either via relabelling or unintentionally because of an unknown reference ID. |
| prometheus_remote_storage_exemplars_failed | counter | Total number of exemplars which failed on send to remote storage, non-recoverable errors. |
| prometheus_remote_storage_exemplars_in | counter | Exemplars in to remote storage, compare to exemplars out for queue managers. |
| prometheus_remote_storage_exemplars_pending | gauge | The number of exemplars pending in the queues shards to be sent to the remote storage. |
| prometheus_remote_storage_exemplars_retried | counter | Total number of exemplars which failed on send to remote storage but were retried because the send error was recoverable. |
| prometheus_remote_storage_exemplars | counter | Total number of exemplars sent to remote storage. |
| prometheus_remote_storage_highest_timestamp_in_seconds | gauge | Highest timestamp that has come into the remote storage via the Appender interface, in seconds since epoch. |
| prometheus_remote_storage_max_samples_per_send | gauge | The maximum number of samples to be sent, in a single request, to the remote storage. Note that, when sending of exemplars over remote write is enabled, exemplars count towards this limt. |
| prometheus_remote_storage_metadata_bytes | counter | The total number of bytes of metadata sent by the queue after compression. |
| prometheus_remote_storage_metadata_failed | counter | Total number of metadata entries which failed on send to remote storage, non-recoverable errors. |
| prometheus_remote_storage_metadata_retried | counter | Total number of metadata entries which failed on send to remote storage but were retried because the send error was recoverable. |
| prometheus_remote_storage_metadata | counter | Total number of metadata entries sent to remote storage. |
| prometheus_remote_storage_queue_highest_sent_timestamp_seconds | gauge | Timestamp from a WAL sample, the highest timestamp successfully sent by this queue, in seconds since epoch. |
| prometheus_remote_storage_samples_dropped | counter | Total number of samples which were dropped after being read from the WAL before being sent via remote write, either via relabelling or unintentionally because of an unknown reference ID. |
| prometheus_remote_storage_samples_failed | counter | Total number of samples which failed on send to remote storage, non-recoverable errors. |
| prometheus_remote_storage_samples_in | counter | Samples in to remote storage, compare to samples out for queue managers. |
| prometheus_remote_storage_samples_pending | gauge | The number of samples pending in the queues shards to be sent to the remote storage. |
| prometheus_remote_storage_samples_retried | counter | Total number of samples which failed on send to remote storage but were retried because the send error was recoverable. |
| prometheus_remote_storage_samples | counter | Total number of samples sent to remote storage. |
| prometheus_remote_storage_sent_batch_duration_seconds | histogram | Duration of send calls to the remote storage. |
| prometheus_remote_storage_shard_capacity | gauge | The capacity of each shard of the queue used for parallel sending to the remote storage. |
| prometheus_remote_storage_shards | gauge | The number of shards used for parallel sending to the remote storage. |
| prometheus_remote_storage_shards_desired | gauge | The number of shards that the queues shard calculation wants to run based on the rate of samples in vs. samples out. |
| prometheus_remote_storage_shards_max | gauge | The maximum number of shards that the queue is allowed to run. |
| prometheus_remote_storage_shards_min | gauge | The minimum number of shards that the queue is allowed to run. |
| prometheus_remote_storage_string_interner_zero_reference_releases | counter | The number of times release has been called for strings that are not interned. |
| prometheus_rule_evaluation_duration_seconds | summary | The duration for a rule to execute. |
| prometheus_rule_evaluation_failures | counter | The total number of rule evaluation failures. |
| prometheus_rule_evaluations | counter | The total number of rule evaluations. |
| prometheus_rule_group_duration_seconds | summary | The duration of rule group evaluations. |
| prometheus_rule_group_interval_seconds | gauge | The interval of a rule group. |
| prometheus_rule_group_iterations_missed | counter | The total number of rule group evaluations missed due to slow rule group evaluation. |
| prometheus_rule_group_iterations | counter | The total number of scheduled rule group evaluations, whether executed or missed. |
| prometheus_rule_group_last_duration_seconds | gauge | The duration of the last rule group evaluation. |
| prometheus_rule_group_last_evaluation_samples | gauge | The number of samples returned during the last rule group evaluation. |
| prometheus_rule_group_last_evaluation_timestamp_seconds | gauge | The timestamp of the last rule group evaluation in seconds. |
| prometheus_rule_group_rules | gauge | The number of rules. |
| prometheus_sd_consul_rpc_duration_seconds | summary | The duration of a Consul RPC call in seconds. |
| prometheus_sd_consul_rpc_failures | counter | The number of Consul RPC call failures. |
| prometheus_sd_discovered_targets | gauge | Current number of discovered targets. |
| prometheus_sd_dns_lookup_failures | counter | The number of DNS-SD lookup failures. |
| prometheus_sd_dns_lookups | counter | The number of DNS-SD lookups. |
| prometheus_sd_failed_configs | gauge | Current number of service discovery configurations that failed to load. |
| prometheus_sd_file_read_errors | counter | The number of File-SD read errors. |
| prometheus_sd_file_scan_duration_seconds | summary | The duration of the File-SD scan in seconds. |
| prometheus_sd_kubernetes_events | counter | The number of Kubernetes events handled. |
| prometheus_sd_kubernetes_http_request_duration_seconds | summary | Summary of latencies for HTTP requests to the Kubernetes API by endpoint. |
| prometheus_sd_kubernetes_http_request | counter | Total number of HTTP requests to the Kubernetes API by status code. |
| prometheus_sd_kubernetes_workqueue_depth | gauge | Current depth of the work queue. |
| prometheus_sd_kubernetes_workqueue_items | counter | Total number of items added to the work queue. |
| prometheus_sd_kubernetes_workqueue_latency_seconds | summary | How long an item stays in the work queue. |
| prometheus_sd_kubernetes_workqueue_longest_running_processor_seconds | gauge | Duration of the longest running processor in the work queue. |
| prometheus_sd_kubernetes_workqueue_unfinished_work_seconds | gauge | How long an item has remained unfinished in the work queue. |
| prometheus_sd_kubernetes_workqueue_work_duration_seconds | summary | How long processing an item from the work queue takes. |
| prometheus_sd_kuma_fetch_duration_seconds | summary | The duration of a Kuma MADS fetch call. |
| prometheus_sd_kuma_fetch_failures | counter | The number of Kuma MADS fetch call failures. |
| prometheus_sd_kuma_fetch_skipped_updates | counter | The number of Kuma MADS fetch calls that result in no updates to the targets. |
| prometheus_sd_received_updates | counter | Total number of update events received from the SD providers. |
| prometheus_sd_updates | counter | Total number of update events sent to the SD consumers. |
| prometheus_target_interval_length_seconds | summary | Actual intervals between scrapes. |
| prometheus_target_metadata_cache_bytes | gauge | The number of bytes that are currently used for storing metric metadata in the cache |
| prometheus_target_metadata_cache_entries | gauge | Total number of metric metadata entries in the cache |
| prometheus_target_scrape_pool_exceeded_label_limits | counter | Total number of times scrape pools hit the label limits, during sync or config reload. |
| prometheus_target_scrape_pool_exceeded_target_limit | counter | Total number of times scrape pools hit the target limit, during sync or config reload. |
| prometheus_target_scrape_pool_reloads_failed | counter | Total number of failed scrape pool reloads. |
| prometheus_target_scrape_pool_reloads | counter | Total number of scrape pool reloads. |
| prometheus_target_scrape_pool_sync | counter | Total number of syncs that were executed on a scrape pool. |
| prometheus_target_scrape_pool_targets | gauge | Current number of targets in this scrape pool. |
| prometheus_target_scrape_pools_failed | counter | Total number of scrape pool creations that failed. |
| prometheus_target_scrape_pools | counter | Total number of scrape pool creation attempts. |
| prometheus_target_scrapes_cache_flush_forced | counter | How many times a scrape cache was flushed due to getting big while scrapes are failing. |
| prometheus_target_scrapes_exceeded_body_size_limit | counter | Total number of scrapes that hit the body size limit |
| prometheus_target_scrapes_exceeded_sample_limit | counter | Total number of scrapes that hit the sample limit and were rejected. |
| prometheus_target_scrapes_exemplar_out_of_order | counter | Total number of exemplar rejected due to not being out of the expected order. |
| prometheus_target_scrapes_sample_duplicate_timestamp | counter | Total number of samples rejected due to duplicate timestamps but different values. |
| prometheus_target_scrapes_sample_out_of_bounds | counter | Total number of samples rejected due to timestamp falling outside of the time bounds. |
| prometheus_target_scrapes_sample_out_of_order | counter | Total number of samples rejected due to not being out of the expected order. |
| prometheus_target_sync_failed | counter | Total number of target sync failures. |
| prometheus_target_sync_length_seconds | summary | Actual interval to sync the scrape pool. |
| prometheus_template_text_expansion_failures | counter | The total number of template text expansion failures. |
| prometheus_template_text_expansions | counter | The total number of template text expansions. |
| prometheus_treecache_watcher_goroutines | gauge | The current number of watcher goroutines. |
| prometheus_treecache_zookeeper_failures | counter | The total number of ZooKeeper failures. |
| prometheus_tsdb_blocks_loaded | gauge | Number of currently loaded data blocks |
| prometheus_tsdb_checkpoint_creations_failed | counter | Total number of checkpoint creations that failed. |
| prometheus_tsdb_checkpoint_creations | counter | Total number of checkpoint creations attempted. |
| prometheus_tsdb_checkpoint_deletions_failed | counter | Total number of checkpoint deletions that failed. |
| prometheus_tsdb_checkpoint_deletions | counter | Total number of checkpoint deletions attempted. |
| prometheus_tsdb_chunk_write_queue_operations | counter | Number of operations on the chunk_write_queue. |
| prometheus_tsdb_clean_start | gauge | -1: lockfile is disabled. 0: a lockfile from a previous execution was replaced. 1: lockfile creation was clean |
| prometheus_tsdb_compaction_chunk_range_seconds | histogram | Final time range of chunks on their first compaction |
| prometheus_tsdb_compaction_chunk_samples | histogram | Final number of samples on their first compaction |
| prometheus_tsdb_compaction_chunk_size_bytes | histogram | Final size of chunks on their first compaction |
| prometheus_tsdb_compaction_duration_seconds | histogram | Duration of compaction runs |
| prometheus_tsdb_compaction_populating_block | gauge | Set to 1 when a block is currently being written to the disk. |
| prometheus_tsdb_compactions_failed | counter | Total number of compactions that failed for the partition. |
| prometheus_tsdb_compactions_skipped | counter | Total number of skipped compactions due to disabled auto compaction. |
| prometheus_tsdb_compactions | counter | Total number of compactions that were executed for the partition. |
| prometheus_tsdb_compactions_triggered | counter | Total number of triggered compactions for the partition. |
| prometheus_tsdb_data_replay_duration_seconds | gauge | Time taken to replay the data on disk. |
| prometheus_tsdb_exemplar_exemplars_appended | counter | Total number of appended exemplars. |
| prometheus_tsdb_exemplar_exemplars_in_storage | gauge | Number of exemplars currently in circular storage. |
| prometheus_tsdb_exemplar_last_exemplars_timestamp_seconds | gauge | The timestamp of the oldest exemplar stored in circular storage. Useful to check for what timerange the current exemplar buffer limit allows. This usually means the last timestampfor all exemplars for a typical setup. This is not true though if one of the series timestamp is in future compared to rest series. |
| prometheus_tsdb_exemplar_max_exemplars | gauge | Total number of exemplars the exemplar storage can store, resizeable. |
| prometheus_tsdb_exemplar_out_of_order_exemplars | counter | Total number of out of order exemplar ingestion failed attempts. |
| prometheus_tsdb_exemplar_series_with_exemplars_in_storage | gauge | Number of series with exemplars currently in circular storage. |
| prometheus_tsdb_head_active_appenders | gauge | Number of currently active appender transactions |
| prometheus_tsdb_head_chunks | gauge | Total number of chunks in the head block. |
| prometheus_tsdb_head_chunks_created | counter | Total number of chunks created in the head |
| prometheus_tsdb_head_chunks_removed | counter | Total number of chunks removed in the head |
| prometheus_tsdb_head_gc_duration_seconds | summary | Runtime of garbage collection in the head block. |
| prometheus_tsdb_head_max_time | gauge | Maximum timestamp of the head block. The unit is decided by the library consumer. |
| prometheus_tsdb_head_max_time_seconds | gauge | Maximum timestamp of the head block. |
| prometheus_tsdb_head_min_time | gauge | Minimum time bound of the head block. The unit is decided by the library consumer. |
| prometheus_tsdb_head_min_time_seconds | gauge | Minimum time bound of the head block. |
| prometheus_tsdb_head_samples_appended | counter | Total number of appended samples. |
| prometheus_tsdb_head_series | gauge | Total number of series in the head block. |
| prometheus_tsdb_head_series_created | counter | Total number of series created in the head |
| prometheus_tsdb_head_series_not_found | counter | Total number of requests for series that were not found. |
| prometheus_tsdb_head_series_removed | counter | Total number of series removed in the head |
| prometheus_tsdb_head_truncations_failed | counter | Total number of head truncations that failed. |
| prometheus_tsdb_head_truncations | counter | Total number of head truncations attempted. |
| prometheus_tsdb_isolation_high_watermark | gauge | The highest TSDB append ID that has been given out. |
| prometheus_tsdb_isolation_low_watermark | gauge | The lowest TSDB append ID that is still referenced. |
| prometheus_tsdb_lowest_timestamp | gauge | Lowest timestamp value stored in the database. The unit is decided by the library consumer. |
| prometheus_tsdb_lowest_timestamp_seconds | gauge | Lowest timestamp value stored in the database. |
| prometheus_tsdb_mmap_chunk_corruptions | counter | Total number of memory-mapped chunk corruptions. |
| prometheus_tsdb_out_of_bound_samples | counter | Total number of out of bound samples ingestion failed attempts. |
| prometheus_tsdb_out_of_order_samples | counter | Total number of out of order samples ingestion failed attempts. |
| prometheus_tsdb_reloads_failures | counter | Number of times the database failed to reloadBlocks block data from disk. |
| prometheus_tsdb_reloads | counter | Number of times the database reloaded block data from disk. |
| prometheus_tsdb_retention_limit_bytes | gauge | Max number of bytes to be retained in the tsdb blocks, configured 0 means disabled |
| prometheus_tsdb_size_retentions | counter | The number of times that blocks were deleted because the maximum number of bytes was exceeded. |
| prometheus_tsdb_snapshot_replay_error | counter | Total number snapshot replays that failed. |
| prometheus_tsdb_storage_blocks_bytes | gauge | The number of bytes that are currently used for local storage by all blocks. |
| prometheus_tsdb_symbol_table_size_bytes | gauge | Size of symbol table in memory for loaded blocks |
| prometheus_tsdb_time_retentions | counter | The number of times that blocks were deleted because the maximum time limit was exceeded. |
| prometheus_tsdb_tombstone_cleanup_seconds | histogram | The time taken to recompact blocks to remove tombstones. |
| prometheus_tsdb_vertical_compactions | counter | Total number of compactions done on overlapping blocks. |
| prometheus_tsdb_wal_completed_pages | counter | Total number of completed pages. |
| prometheus_tsdb_wal_corruptions | counter | Total number of WAL corruptions. |
| prometheus_tsdb_wal_fsync_duration_seconds | summary | Duration of WAL fsync. |
| prometheus_tsdb_wal_page_flushes | counter | Total number of page flushes. |
| prometheus_tsdb_wal_segment_current | gauge | WAL segment index that TSDB is currently writing to. |
| prometheus_tsdb_wal_truncate_duration_seconds | summary | Duration of WAL truncation. |
| prometheus_tsdb_wal_truncations_failed | counter | Total number of WAL truncations that failed. |
| prometheus_tsdb_wal_truncations | counter | Total number of WAL truncations attempted. |
| prometheus_tsdb_wal_writes_failed | counter | Total number of WAL writes that failed. |
| prometheus_wal_watcher_current_segment | gauge | Current segment the WAL watcher is reading records from. |
| prometheus_wal_watcher_record_decode_failures | counter | Number of records read by the WAL watcher that resulted in an error when decoding. |
| prometheus_wal_watcher_records_read | counter | Number of records read by the WAL watcher from the WAL. |
| prometheus_wal_watcher_samples_sent_pre_tailing | counter | Number of sample records read by the WAL watcher and sent to remote write during replay of existing WAL. |
| prometheus_web_federation_errors | counter | Total number of errors that occurred while sending federation responses. |
| prometheus_web_federation_warnings | counter | Total number of warnings that occurred while sending federation responses. |
| promhttp_metric_handler_requests_in_flight | gauge | Current number of scrapes being served. |
| promhttp_metric_handler_requests | counter | Total number of scrapes by HTTP status code. |
<!-- markdownlint-enable line-length -->

## Prometheus-operator (Kubernetes v1.21)

### Statistics

<!-- markdownlint-disable line-length -->
| Number of metrics families | Number of time series | Average number of labels per time series | Number of unique label pairs |
| -------------------------- | --------------------- | ---------------------------------------- | ---------------------------- |
| 87 | 1824 | 1.01 | 851 |
<!-- markdownlint-enable line-length -->

### Metrics

<!-- markdownlint-disable line-length -->
| Name | Type, Unit | Description |
| ---- | ---------- | ----------- |
| go_gc_cycles_automatic_gc_cycles | counter | Count of completed GC cycles generated by the Go runtime. |
| go_gc_cycles_forced_gc_cycles | counter | Count of completed GC cycles forced by the application. |
| go_gc_cycles_total_gc_cycles | counter | Count of all completed GC cycles. |
| go_gc_duration_seconds | summary | A summary of the pause duration of garbage collection cycles. |
| go_gc_heap_allocs_by_size_bytes_total | histogram | Distribution of heap allocations by approximate size. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks. |
| go_gc_heap_allocs_bytes | counter | Cumulative sum of memory allocated to the heap by the application. |
| go_gc_heap_allocs_objects | counter | Cumulative count of heap allocations triggered by the application. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks. |
| go_gc_heap_frees_by_size_bytes_total | histogram | Distribution of freed heap allocations by approximate size. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks. |
| go_gc_heap_frees_bytes | counter | Cumulative sum of heap memory freed by the garbage collector. |
| go_gc_heap_frees_objects | counter | Cumulative count of heap allocations whose storage was freed by the garbage collector. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks. |
| go_gc_heap_goal_bytes | gauge | Heap size target for the end of the GC cycle. |
| go_gc_heap_objects_objects | gauge | Number of objects, live or unswept, occupying heap memory. |
| go_gc_heap_tiny_allocs_objects | counter | Count of small allocations that are packed together into blocks. These allocations are counted separately from other allocations because each individual allocation is not tracked by the runtime, only their block. Each block is already accounted for in allocs-by-size and frees-by-size. |
| go_gc_pauses_seconds_total | histogram | Distribution individual GC-related stop-the-world pause latencies. |
| go_goroutines | gauge | Number of goroutines that currently exist. |
| go_info | gauge | Information about the Go environment. |
| go_memory_classes_heap_free_bytes | gauge | Memory that is completely free and eligible to be returned to the underlying system, but has not been. This metric is the runtime's estimate of free address space that is backed by physical memory. |
| go_memory_classes_heap_objects_bytes | gauge | Memory occupied by live objects and dead objects that have not yet been marked free by the garbage collector. |
| go_memory_classes_heap_released_bytes | gauge | Memory that is completely free and has been returned to the underlying system. This metric is the runtime's estimate of free address space that is still mapped into the process, but is not backed by physical memory. |
| go_memory_classes_heap_stacks_bytes | gauge | Memory allocated from the heap that is reserved for stack space, whether or not it is currently in-use. |
| go_memory_classes_heap_unused_bytes | gauge | Memory that is reserved for heap objects but is not currently used to hold heap objects. |
| go_memory_classes_metadata_mcache_free_bytes | gauge | Memory that is reserved for runtime mcache structures, but not in-use. |
| go_memory_classes_metadata_mcache_inuse_bytes | gauge | Memory that is occupied by runtime mcache structures that are currently being used. |
| go_memory_classes_metadata_mspan_free_bytes | gauge | Memory that is reserved for runtime mspan structures, but not in-use. |
| go_memory_classes_metadata_mspan_inuse_bytes | gauge | Memory that is occupied by runtime mspan structures that are currently being used. |
| go_memory_classes_metadata_other_bytes | gauge | Memory that is reserved for or used to hold runtime metadata. |
| go_memory_classes_os_stacks_bytes | gauge | Stack memory allocated by the underlying operating system. |
| go_memory_classes_other_bytes | gauge | Memory used by execution trace buffers, structures for debugging the runtime, finalizer and profiler specials, and more. |
| go_memory_classes_profiling_buckets_bytes | gauge | Memory that is used by the stack trace hash map used for profiling. |
| go_memory_classes_total_bytes | gauge | All memory mapped by the Go runtime into the current process as read-write. Note that this does not include memory mapped by code called via cgo or via the syscall package. Sum of all metrics in /memory/classes. |
| go_memstats_alloc_bytes | gauge | Number of bytes allocated and still in use. |
| go_memstats_alloc_bytes | counter | Total number of bytes allocated, even if freed. |
| go_memstats_buck_hash_sys_bytes | gauge | Number of bytes used by the profiling bucket hash table. |
| go_memstats_frees | counter | Total number of frees. |
| go_memstats_gc_cpu_fraction | gauge | The fraction of this program's available CPU time used by the GC since the program started. |
| go_memstats_gc_sys_bytes | gauge | Number of bytes used for garbage collection system metadata. |
| go_memstats_heap_alloc_bytes | gauge | Number of heap bytes allocated and still in use. |
| go_memstats_heap_idle_bytes | gauge | Number of heap bytes waiting to be used. |
| go_memstats_heap_inuse_bytes | gauge | Number of heap bytes that are in use. |
| go_memstats_heap_objects | gauge | Number of allocated objects. |
| go_memstats_heap_released_bytes | gauge | Number of heap bytes released to OS. |
| go_memstats_heap_sys_bytes | gauge | Number of heap bytes obtained from system. |
| go_memstats_last_gc_time_seconds | gauge | Number of seconds since 1970 of last garbage collection. |
| go_memstats_lookups | counter | Total number of pointer lookups. |
| go_memstats_mallocs | counter | Total number of mallocs. |
| go_memstats_mcache_inuse_bytes | gauge | Number of bytes in use by mcache structures. |
| go_memstats_mcache_sys_bytes | gauge | Number of bytes used for mcache structures obtained from system. |
| go_memstats_mspan_inuse_bytes | gauge | Number of bytes in use by mspan structures. |
| go_memstats_mspan_sys_bytes | gauge | Number of bytes used for mspan structures obtained from system. |
| go_memstats_next_gc_bytes | gauge | Number of heap bytes when next garbage collection will take place. |
| go_memstats_other_sys_bytes | gauge | Number of bytes used for other system allocations. |
| go_memstats_stack_inuse_bytes | gauge | Number of bytes in use by the stack allocator. |
| go_memstats_stack_sys_bytes | gauge | Number of bytes obtained from system for stack allocator. |
| go_memstats_sys_bytes | gauge | Number of bytes obtained from system. |
| go_sched_goroutines_goroutines | gauge | Count of live goroutines. |
| go_sched_latencies_seconds | histogram | Distribution of the time goroutines have spent in the scheduler in a runnable state before actually running. |
| go_threads | gauge | Number of OS threads created. |
| process_cpu_seconds | counter | Total user and system CPU time spent in seconds. |
| process_max_fds | gauge | Maximum number of open file descriptors. |
| process_open_fds | gauge | Number of open file descriptors. |
| process_resident_memory_bytes | gauge | Resident memory size in bytes. |
| process_start_time_seconds | gauge | Start time of the process since unix epoch in seconds. |
| process_virtual_memory_bytes | gauge | Virtual memory size in bytes. |
| process_virtual_memory_max_bytes | gauge | Maximum amount of virtual memory available in bytes. |
| prometheus_operator_alertmanager_config_validation_errors | counter | Number of errors that occurred while validating a alertmanagerconfig object |
| prometheus_operator_alertmanager_config_validation_triggered | counter | Number of times an alertmanagerconfig object triggered validation |
| prometheus_operator_build_info | gauge | A metric with a constant '1' value labeled by version, revision, branch, and goversion from which prometheus_operator was built. |
| prometheus_operator_kubernetes_client_http_request_duration_seconds | summary | Summary of latencies for the Kubernetes client's requests by endpoint. |
| prometheus_operator_kubernetes_client_http_requests | counter | Total number of Kubernetes's client requests by status code. |
| prometheus_operator_kubernetes_client_rate_limiter_duration_seconds | summary | Summary of latencies for the Kuberntes client's rate limiter by endpoint. |
| prometheus_operator_list_operations_failed | counter | Total number of list operations that failed |
| prometheus_operator_list_operations | counter | Total number of list operations |
| prometheus_operator_managed_resources | gauge | Number of resources managed by the operator's controller per state (selected/rejected) |
| prometheus_operator_node_address_lookup_errors | counter | Number of times a node IP address could not be determined |
| prometheus_operator_node_syncs_failed | counter | Number of node endpoints synchronisation failures |
| prometheus_operator_node_syncs | counter | Number of node endpoints synchronisations |
| prometheus_operator_ready | gauge | 1 when the controller is ready to reconcile resources, 0 otherwise |
| prometheus_operator_reconcile_errors | counter | Number of errors that occurred during reconcile operations |
| prometheus_operator_reconcile_operations | counter | Total number of reconcile operations |
| prometheus_operator_reconcile_sts_delete_create | counter | Number of times that reconciling a statefulset required deleting and re-creating it |
| prometheus_operator_rule_validation_errors | counter | Number of errors that occurred while validating a prometheusRules object |
| prometheus_operator_rule_validation_triggered | counter | Number of times a prometheusRule object triggered validation |
| prometheus_operator_spec_replicas | gauge | Number of expected replicas for the object. |
| prometheus_operator_syncs | gauge | Number of objects per sync status (ok/failed) |
| prometheus_operator_triggered | counter | Number of times a Kubernetes object add, delete or update event triggered the Prometheus Operator to reconcile an object |
| prometheus_operator_watch_operations_failed | counter | Total number of watch operations that failed |
| prometheus_operator_watch_operations | counter | Total number of watch operations |
<!-- markdownlint-enable line-length -->

## AlertManager (Kubernetes v1.23)

### Statistics

<!-- markdownlint-disable line-length -->
| Number of metrics families | Number of time series | Average number of labels per time series | Number of unique label pairs |
| -------------------------- | --------------------- | ---------------------------------------- | ---------------------------- |
| 73 | 312 | 1.53 | 60 |
<!-- markdownlint-enable line-length -->

### Metrics

<!-- markdownlint-disable line-length -->
| Name | Type, Unit | Description |
| ---- | ---------- | ----------- |
| alertmanager_alerts | gauge | How many alerts by state. |
| alertmanager_alerts_invalid | counter | The total number of received alerts that were invalid. |
| alertmanager_alerts_received | counter | The total number of received alerts. |
| alertmanager_build_info | gauge | A metric with a constant '1' value labeled by version, revision, branch, and goversion from which alertmanager was built. |
| alertmanager_cluster_enabled | gauge | Indicates whether the clustering is enabled or not. |
| alertmanager_config_hash | gauge | Hash of the currently loaded alertmanager configuration. |
| alertmanager_config_last_reload_success_timestamp_seconds | gauge | Timestamp of the last successful configuration reload. |
| alertmanager_config_last_reload_successful | gauge | Whether the last configuration reload attempt was successful. |
| alertmanager_dispatcher_aggregation_groups | gauge | Number of active aggregation groups |
| alertmanager_dispatcher_alert_processing_duration_seconds | summary | Summary of latencies for the processing of alerts. |
| alertmanager_http_concurrency_limit_exceeded | counter | Total number of times an HTTP request failed because the concurrency limit was reached. |
| alertmanager_http_request_duration_seconds | histogram | Histogram of latencies for HTTP requests. |
| alertmanager_http_requests_in_flight | gauge | Current number of HTTP requests being processed. |
| alertmanager_http_response_size_bytes | histogram | Histogram of response size for HTTP requests. |
| alertmanager_integrations | gauge | Number of configured integrations. |
| alertmanager_nflog_gc_duration_seconds | summary | Duration of the last notification log garbage collection cycle. |
| alertmanager_nflog_gossip_messages_propagated | counter | Number of received gossip messages that have been further gossiped. |
| alertmanager_nflog_queries | counter | Number of notification log queries were received. |
| alertmanager_nflog_query_duration_seconds | histogram | Duration of notification log query evaluation. |
| alertmanager_nflog_query_errors | counter | Number notification log received queries that failed. |
| alertmanager_nflog_snapshot_duration_seconds | summary | Duration of the last notification log snapshot. |
| alertmanager_nflog_snapshot_size_bytes | gauge | Size of the last notification log snapshot in bytes. |
| alertmanager_notification_latency_seconds | histogram | The latency of notifications in seconds. |
| alertmanager_notification_requests_failed | counter | The total number of failed notification requests. |
| alertmanager_notification_requests | counter | The total number of attempted notification requests. |
| alertmanager_notifications_failed | counter | The total number of failed notifications. |
| alertmanager_notifications | counter | The total number of attempted notifications. |
| alertmanager_receivers | gauge | Number of configured receivers. |
| alertmanager_silences | gauge | How many silences by state. |
| alertmanager_silences_gc_duration_seconds | summary | Duration of the last silence garbage collection cycle. |
| alertmanager_silences_gossip_messages_propagated | counter | Number of received gossip messages that have been further gossiped. |
| alertmanager_silences_queries | counter | How many silence queries were received. |
| alertmanager_silences_query_duration_seconds | histogram | Duration of silence query evaluation. |
| alertmanager_silences_query_errors | counter | How many silence received queries did not succeed. |
| alertmanager_silences_snapshot_duration_seconds | summary | Duration of the last silence snapshot. |
| alertmanager_silences_snapshot_size_bytes | gauge | Size of the last silence snapshot in bytes. |
| go_gc_duration_seconds | summary | A summary of the pause duration of garbage collection cycles. |
| go_goroutines | gauge | Number of goroutines that currently exist. |
| go_info | gauge | Information about the Go environment. |
| go_memstats_alloc_bytes | gauge | Number of bytes allocated and still in use. |
| go_memstats_alloc_bytes | counter | Total number of bytes allocated, even if freed. |
| go_memstats_buck_hash_sys_bytes | gauge | Number of bytes used by the profiling bucket hash table. |
| go_memstats_frees | counter | Total number of frees. |
| go_memstats_gc_cpu_fraction | gauge | The fraction of this program's available CPU time used by the GC since the program started. |
| go_memstats_gc_sys_bytes | gauge | Number of bytes used for garbage collection system metadata. |
| go_memstats_heap_alloc_bytes | gauge | Number of heap bytes allocated and still in use. |
| go_memstats_heap_idle_bytes | gauge | Number of heap bytes waiting to be used. |
| go_memstats_heap_inuse_bytes | gauge | Number of heap bytes that are in use. |
| go_memstats_heap_objects | gauge | Number of allocated objects. |
| go_memstats_heap_released_bytes | gauge | Number of heap bytes released to OS. |
| go_memstats_heap_sys_bytes | gauge | Number of heap bytes obtained from system. |
| go_memstats_last_gc_time_seconds | gauge | Number of seconds since 1970 of last garbage collection. |
| go_memstats_lookups | counter | Total number of pointer lookups. |
| go_memstats_mallocs | counter | Total number of mallocs. |
| go_memstats_mcache_inuse_bytes | gauge | Number of bytes in use by mcache structures. |
| go_memstats_mcache_sys_bytes | gauge | Number of bytes used for mcache structures obtained from system. |
| go_memstats_mspan_inuse_bytes | gauge | Number of bytes in use by mspan structures. |
| go_memstats_mspan_sys_bytes | gauge | Number of bytes used for mspan structures obtained from system. |
| go_memstats_next_gc_bytes | gauge | Number of heap bytes when next garbage collection will take place. |
| go_memstats_other_sys_bytes | gauge | Number of bytes used for other system allocations. |
| go_memstats_stack_inuse_bytes | gauge | Number of bytes in use by the stack allocator. |
| go_memstats_stack_sys_bytes | gauge | Number of bytes obtained from system for stack allocator. |
| go_memstats_sys_bytes | gauge | Number of bytes obtained from system. |
| go_threads | gauge | Number of OS threads created. |
| process_cpu_seconds | counter | Total user and system CPU time spent in seconds. |
| process_max_fds | gauge | Maximum number of open file descriptors. |
| process_open_fds | gauge | Number of open file descriptors. |
| process_resident_memory_bytes | gauge | Resident memory size in bytes. |
| process_start_time_seconds | gauge | Start time of the process since unix epoch in seconds. |
| process_virtual_memory_bytes | gauge | Virtual memory size in bytes. |
| process_virtual_memory_max_bytes | gauge | Maximum amount of virtual memory available in bytes. |
| promhttp_metric_handler_requests_in_flight | gauge | Current number of scrapes being served. |
| promhttp_metric_handler_requests | counter | Total number of scrapes by HTTP status code. |
<!-- markdownlint-enable line-length -->

## Core DNS service (Kubernetes v1.23)

### Statistics

<!-- markdownlint-disable line-length -->
| Number of metrics families | Number of time series | Average number of labels per time series | Number of unique label pairs |
| -------------------------- | --------------------- | ---------------------------------------- | ---------------------------- |
| 63 | 322 | 2.53 | 80 |
<!-- markdownlint-enable line-length -->

### Metrics

<!-- markdownlint-disable line-length -->
| Name | Type, Unit | Description |
| ---- | ---------- | ----------- |
| coredns_build_info | gauge | A metric with a constant '1' value labeled by version, revision, and goversion from which CoreDNS was built. |
| coredns_cache_entries | gauge | The number of elements in the cache. |
| coredns_cache_hits | counter | The count of cache hits. |
| coredns_cache_misses | counter | The count of cache misses. Deprecated, derive misses from cache hits/requests counters. |
| coredns_cache_requests | counter | The count of cache requests. |
| coredns_dns_request_duration_seconds | histogram | Histogram of the time (in seconds) each request took per zone. |
| coredns_dns_request_size_bytes | histogram | Size of the EDNS0 UDP buffer in bytes (64K for TCP) per zone and protocol. |
| coredns_dns_requests | counter | Counter of DNS requests made per zone, protocol and family. |
| coredns_dns_response_size_bytes | histogram | Size of the returned response in bytes. |
| coredns_dns_responses | counter | Counter of response status codes. |
| coredns_forward_conn_cache_hits | counter | Counter of connection cache hits per upstream and protocol. |
| coredns_forward_conn_cache_misses | counter | Counter of connection cache misses per upstream and protocol. |
| coredns_forward_healthcheck_broken | counter | Counter of the number of complete failures of the healthchecks. |
| coredns_forward_healthcheck_failures | counter | Counter of the number of failed healthchecks. |
| coredns_forward_max_concurrent_rejects | counter | Counter of the number of queries rejected because the concurrent queries were at maximum. |
| coredns_forward_request_duration_seconds | histogram | Histogram of the time each request took. |
| coredns_forward_requests | counter | Counter of requests made per upstream. |
| coredns_forward_responses | counter | Counter of responses received per upstream. |
| coredns_health_request_duration_seconds | histogram | Histogram of the time (in seconds) each request took. |
| coredns_health_request_failures | counter | The number of times the health check failed. |
| coredns_hosts_entries | gauge | The combined number of entries in hosts and Corefile. |
| coredns_hosts_reload_timestamp_seconds | gauge | The timestamp of the last reload of hosts file. |
| coredns_kubernetes_dns_programming_duration_seconds | histogram | Histogram of the time (in seconds) it took to program a dns instance. |
| coredns_local_localhost_requests | counter | Counter of localhost.<domain> requests. |
| coredns_panics | counter | A metrics that counts the number of panics. |
| coredns_plugin_enabled | gauge | A metric that indicates whether a plugin is enabled on per server and zone basis. |
| coredns_reload_failed | counter | Counter of the number of failed reload attempts. |
| coredns_template_matches | counter | Counter of template regex matches. |
| go_gc_duration_seconds | summary | A summary of the pause duration of garbage collection cycles. |
| go_goroutines | gauge | Number of goroutines that currently exist. |
| go_info | gauge | Information about the Go environment. |
| go_memstats_alloc_bytes | gauge | Number of bytes allocated and still in use. |
| go_memstats_alloc_bytes | counter | Total number of bytes allocated, even if freed. |
| go_memstats_buck_hash_sys_bytes | gauge | Number of bytes used by the profiling bucket hash table. |
| go_memstats_frees | counter | Total number of frees. |
| go_memstats_gc_cpu_fraction | gauge | The fraction of this program's available CPU time used by the GC since the program started. |
| go_memstats_gc_sys_bytes | gauge | Number of bytes used for garbage collection system metadata. |
| go_memstats_heap_alloc_bytes | gauge | Number of heap bytes allocated and still in use. |
| go_memstats_heap_idle_bytes | gauge | Number of heap bytes waiting to be used. |
| go_memstats_heap_inuse_bytes | gauge | Number of heap bytes that are in use. |
| go_memstats_heap_objects | gauge | Number of allocated objects. |
| go_memstats_heap_released_bytes | gauge | Number of heap bytes released to OS. |
| go_memstats_heap_sys_bytes | gauge | Number of heap bytes obtained from system. |
| go_memstats_last_gc_time_seconds | gauge | Number of seconds since 1970 of last garbage collection. |
| go_memstats_lookups | counter | Total number of pointer lookups. |
| go_memstats_mallocs | counter | Total number of mallocs. |
| go_memstats_mcache_inuse_bytes | gauge | Number of bytes in use by mcache structures. |
| go_memstats_mcache_sys_bytes | gauge | Number of bytes used for mcache structures obtained from system. |
| go_memstats_mspan_inuse_bytes | gauge | Number of bytes in use by mspan structures. |
| go_memstats_mspan_sys_bytes | gauge | Number of bytes used for mspan structures obtained from system. |
| go_memstats_next_gc_bytes | gauge | Number of heap bytes when next garbage collection will take place. |
| go_memstats_other_sys_bytes | gauge | Number of bytes used for other system allocations. |
| go_memstats_stack_inuse_bytes | gauge | Number of bytes in use by the stack allocator. |
| go_memstats_stack_sys_bytes | gauge | Number of bytes obtained from system for stack allocator. |
| go_memstats_sys_bytes | gauge | Number of bytes obtained from system. |
| go_threads | gauge | Number of OS threads created. |
| process_cpu_seconds | counter | Total user and system CPU time spent in seconds. |
| process_max_fds | gauge | Maximum number of open file descriptors. |
| process_open_fds | gauge | Number of open file descriptors. |
| process_resident_memory_bytes | gauge | Resident memory size in bytes. |
| process_start_time_seconds | gauge | Start time of the process since unix epoch in seconds. |
| process_virtual_memory_bytes | gauge | Virtual memory size in bytes. |
| process_virtual_memory_max_bytes | gauge | Maximum amount of virtual memory available in bytes. |
<!-- markdownlint-enable line-length -->

## Grafana (Kubernetes v1.23)

### Statistics

<!-- markdownlint-disable line-length -->
| Number of metrics families | Number of time series | Average number of labels per time series | Number of unique label pairs |
| -------------------------- | --------------------- | ---------------------------------------- | ---------------------------- |
| 147 | 777 | 2.08 | 238 |
<!-- markdownlint-enable line-length -->

### Metrics

<!-- markdownlint-disable line-length -->
| Name | Type, Unit | Description |
| ---- | ---------- | ----------- |
| access_evaluation_duration | histogram | Histogram for the runtime of evaluation function. |
| access_permissions_duration | histogram | Histogram for the runtime of permissions check function. |
| cortex_experimental_features_in_use | counter | The number of experimental features in use. |
| deprecated_flags_inuse | counter | The number of deprecated flags currently set. |
| go_gc_cycles_automatic_gc_cycles | counter | Count of completed GC cycles generated by the Go runtime. |
| go_gc_cycles_forced_gc_cycles | counter | Count of completed GC cycles forced by the application. |
| go_gc_cycles_total_gc_cycles | counter | Count of all completed GC cycles. |
| go_gc_duration_seconds | summary | A summary of the pause duration of garbage collection cycles. |
| go_gc_heap_allocs_by_size_bytes_total | histogram | Distribution of heap allocations by approximate size. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks. |
| go_gc_heap_allocs_bytes | counter | Cumulative sum of memory allocated to the heap by the application. |
| go_gc_heap_allocs_objects | counter | Cumulative count of heap allocations triggered by the application. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks. |
| go_gc_heap_frees_by_size_bytes_total | histogram | Distribution of freed heap allocations by approximate size. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks. |
| go_gc_heap_frees_bytes | counter | Cumulative sum of heap memory freed by the garbage collector. |
| go_gc_heap_frees_objects | counter | Cumulative count of heap allocations whose storage was freed by the garbage collector. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks. |
| go_gc_heap_goal_bytes | gauge | Heap size target for the end of the GC cycle. |
| go_gc_heap_objects_objects | gauge | Number of objects, live or unswept, occupying heap memory. |
| go_gc_heap_tiny_allocs_objects | counter | Count of small allocations that are packed together into blocks. These allocations are counted separately from other allocations because each individual allocation is not tracked by the runtime, only their block. Each block is already accounted for in allocs-by-size and frees-by-size. |
| go_gc_pauses_seconds_total | histogram | Distribution individual GC-related stop-the-world pause latencies. |
| go_goroutines | gauge | Number of goroutines that currently exist. |
| go_info | gauge | Information about the Go environment. |
| go_memory_classes_heap_free_bytes | gauge | Memory that is completely free and eligible to be returned to the underlying system, but has not been. This metric is the runtime's estimate of free address space that is backed by physical memory. |
| go_memory_classes_heap_objects_bytes | gauge | Memory occupied by live objects and dead objects that have not yet been marked free by the garbage collector. |
| go_memory_classes_heap_released_bytes | gauge | Memory that is completely free and has been returned to the underlying system. This metric is the runtime's estimate of free address space that is still mapped into the process, but is not backed by physical memory. |
| go_memory_classes_heap_stacks_bytes | gauge | Memory allocated from the heap that is reserved for stack space, whether or not it is currently in-use. |
| go_memory_classes_heap_unused_bytes | gauge | Memory that is reserved for heap objects but is not currently used to hold heap objects. |
| go_memory_classes_metadata_mcache_free_bytes | gauge | Memory that is reserved for runtime mcache structures, but not in-use. |
| go_memory_classes_metadata_mcache_inuse_bytes | gauge | Memory that is occupied by runtime mcache structures that are currently being used. |
| go_memory_classes_metadata_mspan_free_bytes | gauge | Memory that is reserved for runtime mspan structures, but not in-use. |
| go_memory_classes_metadata_mspan_inuse_bytes | gauge | Memory that is occupied by runtime mspan structures that are currently being used. |
| go_memory_classes_metadata_other_bytes | gauge | Memory that is reserved for or used to hold runtime metadata. |
| go_memory_classes_os_stacks_bytes | gauge | Stack memory allocated by the underlying operating system. |
| go_memory_classes_other_bytes | gauge | Memory used by execution trace buffers, structures for debugging the runtime, finalizer and profiler specials, and more. |
| go_memory_classes_profiling_buckets_bytes | gauge | Memory that is used by the stack trace hash map used for profiling. |
| go_memory_classes_total_bytes | gauge | All memory mapped by the Go runtime into the current process as read-write. Note that this does not include memory mapped by code called via cgo or via the syscall package. Sum of all metrics in /memory/classes. |
| go_memstats_alloc_bytes | gauge | Number of bytes allocated and still in use. |
| go_memstats_alloc_bytes | counter | Total number of bytes allocated, even if freed. |
| go_memstats_buck_hash_sys_bytes | gauge | Number of bytes used by the profiling bucket hash table. |
| go_memstats_frees | counter | Total number of frees. |
| go_memstats_gc_cpu_fraction | gauge | The fraction of this program's available CPU time used by the GC since the program started. |
| go_memstats_gc_sys_bytes | gauge | Number of bytes used for garbage collection system metadata. |
| go_memstats_heap_alloc_bytes | gauge | Number of heap bytes allocated and still in use. |
| go_memstats_heap_idle_bytes | gauge | Number of heap bytes waiting to be used. |
| go_memstats_heap_inuse_bytes | gauge | Number of heap bytes that are in use. |
| go_memstats_heap_objects | gauge | Number of allocated objects. |
| go_memstats_heap_released_bytes | gauge | Number of heap bytes released to OS. |
| go_memstats_heap_sys_bytes | gauge | Number of heap bytes obtained from system. |
| go_memstats_last_gc_time_seconds | gauge | Number of seconds since 1970 of last garbage collection. |
| go_memstats_lookups | counter | Total number of pointer lookups. |
| go_memstats_mallocs | counter | Total number of mallocs. |
| go_memstats_mcache_inuse_bytes | gauge | Number of bytes in use by mcache structures. |
| go_memstats_mcache_sys_bytes | gauge | Number of bytes used for mcache structures obtained from system. |
| go_memstats_mspan_inuse_bytes | gauge | Number of bytes in use by mspan structures. |
| go_memstats_mspan_sys_bytes | gauge | Number of bytes used for mspan structures obtained from system. |
| go_memstats_next_gc_bytes | gauge | Number of heap bytes when next garbage collection will take place. |
| go_memstats_other_sys_bytes | gauge | Number of bytes used for other system allocations. |
| go_memstats_stack_inuse_bytes | gauge | Number of bytes in use by the stack allocator. |
| go_memstats_stack_sys_bytes | gauge | Number of bytes obtained from system for stack allocator. |
| go_memstats_sys_bytes | gauge | Number of bytes obtained from system. |
| go_sched_goroutines_goroutines | gauge | Count of live goroutines. |
| go_sched_latencies_seconds | histogram | Distribution of the time goroutines have spent in the scheduler in a runnable state before actually running. |
| go_threads | gauge | Number of OS threads created. |
| grafana_access_evaluation_count | counter | number of evaluation calls |
| grafana_alerting_active_alerts | gauge | amount of active alerts |
| grafana_alerting_active_configurations | gauge | The number of active Alertmanager configurations. |
| grafana_alerting_alerts | gauge | How many alerts by state. |
| grafana_alerting_discovered_configurations | gauge | The number of organizations we've discovered that require an Alertmanager configuration. |
| grafana_alerting_execution_time_milliseconds | summary | summary of alert execution duration |
| grafana_alerting_get_alert_rules_duration_seconds | histogram | The time taken to get all alert rules. |
| grafana_alerting_schedule_periodic_duration_seconds | histogram | The time taken to run the scheduler. |
| grafana_alerting_scheduler_behind_seconds | gauge | The total number of seconds the scheduler is behind. |
| grafana_api_admin_user_created | counter | api admin user created counter |
| grafana_api_dashboard_get_milliseconds | summary | summary for dashboard get duration |
| grafana_api_dashboard_save_milliseconds | summary | summary for dashboard save duration |
| grafana_api_dashboard_search_milliseconds | summary | summary for dashboard search duration |
| grafana_api_dashboard_snapshot_create | counter | dashboard snapshots created |
| grafana_api_dashboard_snapshot_external | counter | external dashboard snapshots created |
| grafana_api_dashboard_snapshot_get | counter | loaded dashboards |
| grafana_api_dataproxy_request_all_milliseconds | summary | summary for dataproxy request duration |
| grafana_api_login_oauth | counter | api login oauth counter |
| grafana_api_login_post | counter | api login post counter |
| grafana_api_login_saml | counter | api login saml counter |
| grafana_api_models_dashboard_insert | counter | dashboards inserted |
| grafana_api_org_create | counter | api org created counter |
| grafana_api_response_status | counter | api http response status |
| grafana_api_user_signup_completed | counter | amount of users who completed the signup flow |
| grafana_api_user_signup_invite | counter | amount of users who have been invited |
| grafana_api_user_signup_started | counter | amount of users who started the signup flow |
| grafana_aws_cloudwatch_get_metric_data | counter | counter for getting metric data time series from aws |
| grafana_aws_cloudwatch_get_metric_statistics | counter | counter for getting metric statistics from aws |
| grafana_aws_cloudwatch_list_metrics | counter | counter for getting list of metrics from aws |
| grafana_build_info | gauge | A metric with a constant '1' value labeled by version, revision, branch, and goversion from which Grafana was built |
| grafana_db_datasource_query_by_id | counter | counter for getting datasource by id |
| grafana_emails_sent_failed | counter | Number of emails Grafana failed to send |
| grafana_emails_sent | counter | Number of emails sent by Grafana |
| grafana_feature_toggles_info | gauge | info metric that exposes what feature toggles are enabled or not |
| grafana_frontend_boot_css_time_seconds | histogram | Frontend boot initial css load |
| grafana_frontend_boot_first_contentful_paint_time_seconds | histogram | Frontend boot first contentful paint |
| grafana_frontend_boot_first_paint_time_seconds | histogram | Frontend boot first paint |
| grafana_frontend_boot_js_done_time_seconds | histogram | Frontend boot initial js load |
| grafana_frontend_boot_load_time_seconds | histogram | Frontend boot time measurement |
| grafana_http_request_duration_seconds | histogram | Histogram of latencies for HTTP requests. |
| grafana_http_request_in_flight | gauge | A gauge of requests currently being served by Grafana. |
| grafana_instance_start | counter | counter for started instances |
| grafana_ldap_users_sync_execution_time | summary | summary for LDAP users sync execution duration |
| grafana_live_client_command_duration_seconds | summary | clientID command duration summary. |
| grafana_live_client_recover | counter | Count of recover operations. |
| grafana_live_node_action_count | counter | Number of node actions called. |
| grafana_live_node_build | gauge | Node build info. |
| grafana_live_node_messages_received_count | counter | Number of messages received. |
| grafana_live_node_messages_sent_count | counter | Number of messages sent. |
| grafana_live_node_num_channels | gauge | Number of channels with one or more subscribers. |
| grafana_live_node_num_clients | gauge | Number of clients connected. |
| grafana_live_node_num_nodes | gauge | Number of nodes in cluster. |
| grafana_live_node_num_subscriptions | gauge | Number of subscriptions. |
| grafana_live_node_num_users | gauge | Number of unique users connected. |
| grafana_live_transport_connect_count | counter | Number of connections to specific transport. |
| grafana_live_transport_messages_sent | counter | Number of messages sent over specific transport. |
| grafana_page_response_status | counter | page http response status |
| grafana_plugin_build_info | gauge | A metric with a constant '1' value labeled by pluginId, pluginType and version from which Grafana plugin was built |
| grafana_proxy_response_status | counter | proxy http response status |
| grafana_rendering_queue_size | gauge | size of rendering queue |
| grafana_stat_active_users | gauge | number of active users |
| grafana_stat_total_orgs | gauge | total amount of orgs |
| grafana_stat_total_playlists | gauge | total amount of playlists |
| grafana_stat_total_users | gauge | total amount of users |
| grafana_stat_totals_active_admins | gauge | total amount of active admins |
| grafana_stat_totals_active_editors | gauge | total amount of active editors |
| grafana_stat_totals_active_viewers | gauge | total amount of viewers |
| grafana_stat_totals_admins | gauge | total amount of admins |
| grafana_stat_totals_annotations | gauge | total amount of annotations in the database |
| grafana_stat_totals_dashboard | gauge | total amount of dashboards |
| grafana_stat_totals_dashboard_versions | gauge | total amount of dashboard versions in the database |
| grafana_stat_totals_datasource | gauge | total number of defined datasources, labeled by pluginId |
| grafana_stat_totals_editors | gauge | total amount of editors |
| grafana_stat_totals_folder | gauge | total amount of folders |
| grafana_stat_totals_library_panels | gauge | total amount of library panels in the database |
| grafana_stat_totals_library_variables | gauge | total amount of library variables in the database |
| grafana_stat_totals_viewers | gauge | total amount of viewers |
| process_cpu_seconds | counter | Total user and system CPU time spent in seconds. |
| process_max_fds | gauge | Maximum number of open file descriptors. |
| process_open_fds | gauge | Number of open file descriptors. |
| process_resident_memory_bytes | gauge | Resident memory size in bytes. |
| process_start_time_seconds | gauge | Start time of the process since unix epoch in seconds. |
| process_virtual_memory_bytes | gauge | Virtual memory size in bytes. |
| process_virtual_memory_max_bytes | gauge | Maximum amount of virtual memory available in bytes. |
| prometheus_template_text_expansion_failures | counter | The total number of template text expansion failures. |
| prometheus_template_text_expansions | counter | The total number of template text expansions. |
<!-- markdownlint-enable line-length -->

## Grafana-operator (Kubernetes v1.23)

### Statistics

<!-- markdownlint-disable line-length -->
| Number of metrics families | Number of time series | Average number of labels per time series | Number of unique label pairs |
| -------------------------- | --------------------- | ---------------------------------------- | ---------------------------- |
| 49 | 1476 | 2.52 | 170 |
<!-- markdownlint-enable line-length -->

### Metrics

<!-- markdownlint-disable line-length -->
| Name | Type, Unit | Description |
| ---- | ---------- | ----------- |
| controller_runtime_active_workers | gauge | Number of currently used workers per controller |
| controller_runtime_max_concurrent_reconciles | gauge | Maximum number of concurrent reconciles per controller |
| controller_runtime_reconcile_errors | counter | Total number of reconciliation errors per controller |
| controller_runtime_reconcile_time_seconds | histogram | Length of time per reconciliation per controller |
| controller_runtime_reconcile | counter | Total number of reconciliations per controller |
| go_gc_duration_seconds | summary | A summary of the pause duration of garbage collection cycles. |
| go_goroutines | gauge | Number of goroutines that currently exist. |
| go_info | gauge | Information about the Go environment. |
| go_memstats_alloc_bytes | gauge | Number of bytes allocated and still in use. |
| go_memstats_alloc_bytes | counter | Total number of bytes allocated, even if freed. |
| go_memstats_buck_hash_sys_bytes | gauge | Number of bytes used by the profiling bucket hash table. |
| go_memstats_frees | counter | Total number of frees. |
| go_memstats_gc_cpu_fraction | gauge | The fraction of this program's available CPU time used by the GC since the program started. |
| go_memstats_gc_sys_bytes | gauge | Number of bytes used for garbage collection system metadata. |
| go_memstats_heap_alloc_bytes | gauge | Number of heap bytes allocated and still in use. |
| go_memstats_heap_idle_bytes | gauge | Number of heap bytes waiting to be used. |
| go_memstats_heap_inuse_bytes | gauge | Number of heap bytes that are in use. |
| go_memstats_heap_objects | gauge | Number of allocated objects. |
| go_memstats_heap_released_bytes | gauge | Number of heap bytes released to OS. |
| go_memstats_heap_sys_bytes | gauge | Number of heap bytes obtained from system. |
| go_memstats_last_gc_time_seconds | gauge | Number of seconds since 1970 of last garbage collection. |
| go_memstats_lookups | counter | Total number of pointer lookups. |
| go_memstats_mallocs | counter | Total number of mallocs. |
| go_memstats_mcache_inuse_bytes | gauge | Number of bytes in use by mcache structures. |
| go_memstats_mcache_sys_bytes | gauge | Number of bytes used for mcache structures obtained from system. |
| go_memstats_mspan_inuse_bytes | gauge | Number of bytes in use by mspan structures. |
| go_memstats_mspan_sys_bytes | gauge | Number of bytes used for mspan structures obtained from system. |
| go_memstats_next_gc_bytes | gauge | Number of heap bytes when next garbage collection will take place. |
| go_memstats_other_sys_bytes | gauge | Number of bytes used for other system allocations. |
| go_memstats_stack_inuse_bytes | gauge | Number of bytes in use by the stack allocator. |
| go_memstats_stack_sys_bytes | gauge | Number of bytes obtained from system for stack allocator. |
| go_memstats_sys_bytes | gauge | Number of bytes obtained from system. |
| go_threads | gauge | Number of OS threads created. |
| process_cpu_seconds | counter | Total user and system CPU time spent in seconds. |
| process_max_fds | gauge | Maximum number of open file descriptors. |
| process_open_fds | gauge | Number of open file descriptors. |
| process_resident_memory_bytes | gauge | Resident memory size in bytes. |
| process_start_time_seconds | gauge | Start time of the process since unix epoch in seconds. |
| process_virtual_memory_bytes | gauge | Virtual memory size in bytes. |
| process_virtual_memory_max_bytes | gauge | Maximum amount of virtual memory available in bytes. |
| rest_client_request_latency_seconds | histogram | Request latency in seconds. Broken down by verb and URL. |
| rest_client_requests | counter | Number of HTTP requests, partitioned by status code, method, and host. |
| workqueue_adds | counter | Total number of adds handled by workqueue |
| workqueue_depth | gauge | Current depth of workqueue |
| workqueue_longest_running_processor_seconds | gauge | How many seconds has the longest running processor for workqueue been running. |
| workqueue_queue_duration_seconds | histogram | How long in seconds an item stays in workqueue before being requested |
| workqueue_retries | counter | Total number of retries handled by workqueue |
| workqueue_unfinished_work_seconds | gauge | How many seconds of work has been done that is in progress and hasn't been observed by work_duration. Large values indicate stuck threads. One can deduce the number of stuck threads by observing the rate at which this increases. |
| workqueue_work_duration_seconds | histogram | How long in seconds processing an item from workqueue takes. |
<!-- markdownlint-enable line-length -->

## Kube-apiserver (Kubernetes v1.23)

### Statistics

<!-- markdownlint-disable line-length -->
| Number of metrics families | Number of time series | Average number of labels per time series | Number of unique label pairs |
| -------------------------- | --------------------- | ---------------------------------------- | ---------------------------- |
| 144 | 49025 | 7.16 | 1151 |
<!-- markdownlint-enable line-length -->

### Metrics

<!-- markdownlint-disable line-length -->
| Name | Type, Unit | Description |
| ---- | ---------- | ----------- |
| aggregator_openapi_v2_regeneration_count | counter | ALPHA Counter of OpenAPI v2 spec regeneration count broken down by causing APIService name and reason. |
| aggregator_openapi_v2_regeneration_duration | gauge | ALPHA Gauge of OpenAPI v2 spec regeneration duration in seconds. |
| aggregator_unavailable_apiservice | gauge | ALPHA Gauge of APIServices which are marked as unavailable broken down by APIService name. |
| aggregator_unavailable_apiservice | counter | ALPHA Counter of APIServices which are marked as unavailable broken down by APIService name and reason. |
| apiextensions_openapi_v2_regeneration_count | counter | ALPHA Counter of OpenAPI v2 spec regeneration count broken down by causing CRD name and reason. |
| apiserver_admission_controller_admission_duration_seconds | histogram | STABLE Admission controller latency histogram in seconds, identified by name and broken out for each operation and API resource and type (validate or admit). |
| apiserver_admission_step_admission_duration_seconds | histogram | STABLE Admission sub-step latency histogram in seconds, broken out for each operation and API resource and step type (validate or admit). |
| apiserver_admission_step_admission_duration_seconds_summary | summary | ALPHA Admission sub-step latency summary in seconds, broken out for each operation and API resource and step type (validate or admit). |
| apiserver_admission_webhook_admission_duration_seconds | histogram | STABLE Admission webhook latency histogram in seconds, identified by name and broken out for each operation and API resource and type (validate or admit). |
| apiserver_admission_webhook_request | counter | ALPHA Admission webhook request total, identified by name and broken out for each admission type (validating or mutating) and operation. Additional labels specify whether the request was rejected or not and an HTTP status code. Codes greater than 600 are truncated to 600, to keep the metrics cardinality bounded. |
| apiserver_audit_event | counter | ALPHA Counter of audit events generated and sent to the audit backend. |
| apiserver_audit_requests_rejected | counter | ALPHA Counter of apiserver requests rejected due to an error in audit logging backend. |
| apiserver_cache_list_fetched_objects | counter | ALPHA Number of objects read from watch cache in the course of serving a LIST request |
| apiserver_cache_list_returned_objects | counter | ALPHA Number of objects returned for a LIST request from watch cache |
| apiserver_cache_list | counter | ALPHA Number of LIST requests served from watch cache |
| apiserver_client_certificate_expiration_seconds | histogram | ALPHA Distribution of the remaining lifetime on the certificate used to authenticate a request. |
| apiserver_current_inflight_requests | gauge | STABLE Maximal number of currently used inflight request limit of this apiserver per request kind in last second. |
| apiserver_current_inqueue_requests | gauge | ALPHA Maximal number of queued requests in this apiserver per request kind in last second. |
| apiserver_envelope_encryption_dek_cache_fill_percent | gauge | ALPHA Percent of the cache slots currently occupied by cached DEKs. |
| apiserver_flowcontrol_current_executing_requests | gauge | ALPHA Number of requests in regular execution phase in the API Priority and Fairness system |
| apiserver_flowcontrol_current_inqueue_requests | gauge | ALPHA Number of requests currently pending in queues of the API Priority and Fairness system |
| apiserver_flowcontrol_current_r | gauge | ALPHA R(time of last change) |
| apiserver_flowcontrol_dispatch_r | gauge | ALPHA R(time of last dispatch) |
| apiserver_flowcontrol_dispatched_requests | counter | ALPHA Number of requests released by API Priority and Fairness system for service |
| apiserver_flowcontrol_latest_s | gauge | ALPHA S(most recently dispatched request) |
| apiserver_flowcontrol_next_discounted_s_bounds | gauge | ALPHA min and max, over queues, of S(oldest waiting request in queue) - estimated work in progress |
| apiserver_flowcontrol_next_s_bounds | gauge | ALPHA min and max, over queues, of S(oldest waiting request in queue) |
| apiserver_flowcontrol_priority_level_request_count_samples | histogram | ALPHA Periodic observations of the number of requests |
| apiserver_flowcontrol_priority_level_request_count_watermarks | histogram | ALPHA Watermarks of the number of requests |
| apiserver_flowcontrol_priority_level_seat_count_samples | histogram | ALPHA Periodic observations of the number of requests |
| apiserver_flowcontrol_priority_level_seat_count_watermarks | histogram | ALPHA Watermarks of the number of requests |
| apiserver_flowcontrol_read_vs_write_request_count_samples | histogram | ALPHA Periodic observations of the number of requests |
| apiserver_flowcontrol_read_vs_write_request_count_watermarks | histogram | ALPHA Watermarks of the number of requests |
| apiserver_flowcontrol_request_concurrency_in_use | gauge | ALPHA Concurrency (number of seats) occupided by the currently executing (all phases count) requests in the API Priority and Fairness system |
| apiserver_flowcontrol_request_concurrency_limit | gauge | ALPHA Shared concurrency limit in the API Priority and Fairness system |
| apiserver_flowcontrol_request_execution_seconds | histogram | ALPHA Duration of regular phase of request execution in the API Priority and Fairness system |
| apiserver_flowcontrol_request_queue_length_after_enqueue | histogram | ALPHA Length of queue in the API Priority and Fairness system, as seen by each request after it is enqueued |
| apiserver_flowcontrol_request_wait_duration_seconds | histogram | ALPHA Length of time a request spent waiting in its queue |
| apiserver_init_events | counter | ALPHA Counter of init events processed in watchcache broken by resource type. |
| apiserver_kube_aggregator_x509_missing_san | counter | ALPHA Counts the number of requests to servers missing SAN extension in their serving certificate OR the number of connection failures due to the lack of x509 certificate SAN extension missing (either/or, based on the runtime environment) |
| apiserver_longrunning_gauge | gauge | ALPHA (Deprecated since 1.23.0) Gauge of all active long-running apiserver requests broken out by verb, group, version, resource, scope and component. Not all requests are tracked this way. |
| apiserver_longrunning_requests | gauge | ALPHA Gauge of all active long-running apiserver requests broken out by verb, group, version, resource, scope and component. Not all requests are tracked this way. |
| apiserver_registered_watchers | gauge | ALPHA (Deprecated since 1.23.0) Number of currently registered watchers for a given resources |
| apiserver_request_aborts | counter | ALPHA Number of requests which apiserver aborted possibly due to a timeout, for each group, version, verb, resource, subresource and scope |
| apiserver_request_duration_seconds | histogram | STABLE Response latency distribution in seconds for each verb, dry run value, group, version, resource, subresource, scope and component. |
| apiserver_request_filter_duration_seconds | histogram | ALPHA Request filter latency distribution in seconds, for each filter type |
| apiserver_request_post_timeout | counter | ALPHA Tracks the activity of the request handlers after the associated requests have been timed out by the apiserver |
| apiserver_request_slo_duration_seconds | histogram | ALPHA Response latency distribution (not counting webhook duration) in seconds for each verb, group, version, resource, subresource, scope and component. |
| apiserver_request_terminations | counter | ALPHA Number of requests which apiserver terminated in self-defense. |
| apiserver_request | counter | STABLE Counter of apiserver requests broken out for each verb, dry run value, group, version, resource, scope, component, and HTTP response code. |
| apiserver_requested_deprecated_apis | gauge | STABLE Gauge of deprecated APIs that have been requested, broken out by API group, version, resource, subresource, and removed_release. |
| apiserver_response_sizes | histogram | STABLE Response size distribution in bytes for each group, version, verb, resource, subresource, scope and component. |
| apiserver_selfrequest | counter | ALPHA Counter of apiserver self-requests broken out for each verb, API resource and subresource. |
| apiserver_storage_data_key_generation_duration_seconds | histogram | ALPHA Latencies in seconds of data encryption key(DEK) generation operations. |
| apiserver_storage_data_key_generation_failures | counter | ALPHA Total number of failed data encryption key(DEK) generation operations. |
| apiserver_storage_envelope_transformation_cache_misses | counter | ALPHA Total number of cache misses while accessing key decryption key(KEK). |
| apiserver_storage_list_evaluated_objects | counter | ALPHA Number of objects tested in the course of serving a LIST request from storage |
| apiserver_storage_list_fetched_objects | counter | ALPHA Number of objects read from storage in the course of serving a LIST request |
| apiserver_storage_list_returned_objects | counter | ALPHA Number of objects returned for a LIST request from storage |
| apiserver_storage_list | counter | ALPHA Number of LIST requests served from storage |
| apiserver_storage_objects | gauge | STABLE Number of stored objects at the time of last check split by kind. |
| apiserver_terminated_watchers | counter | ALPHA Counter of watchers closed due to unresponsiveness broken by resource type. |
| apiserver_tls_handshake_errors | counter | ALPHA Number of requests dropped with 'TLS handshake error from' error |
| apiserver_watch_events_sizes | histogram | ALPHA Watch event size distribution in bytes |
| apiserver_watch_events | counter | ALPHA Number of events sent in watch clients |
| apiserver_webhooks_x509_missing_san | counter | ALPHA Counts the number of requests to servers missing SAN extension in their serving certificate OR the number of connection failures due to the lack of x509 certificate SAN extension missing (either/or, based on the runtime environment) |
| authenticated_user_requests | counter | ALPHA Counter of authenticated requests broken out by username. |
| authentication_attempts | counter | ALPHA Counter of authenticated attempts. |
| authentication_duration_seconds | histogram | ALPHA Authentication duration in seconds broken out by result. |
| authentication_token_cache_active_fetch_count | gauge | ALPHA |
| authentication_token_cache_fetch | counter | ALPHA |
| authentication_token_cache_request_duration_seconds | histogram | ALPHA |
| authentication_token_cache_request | counter | ALPHA |
| etcd_bookmark_counts | gauge | ALPHA Number of etcd bookmarks (progress notify events) split by kind. |
| etcd_db_total_size_in_bytes | gauge | ALPHA Total size of the etcd database file physically allocated in bytes. |
| etcd_lease_object_counts | histogram | ALPHA Number of objects attached to a single etcd lease. |
| etcd_request_duration_seconds | histogram | ALPHA Etcd request latency in seconds for each operation and object type. |
| get_token_count | counter | ALPHA Counter of total Token() requests to the alternate token source |
| get_token_fail_count | counter | ALPHA Counter of failed Token() requests to the alternate token source |
| go_gc_duration_seconds | summary | A summary of the pause duration of garbage collection cycles. |
| go_goroutines | gauge | Number of goroutines that currently exist. |
| go_info | gauge | Information about the Go environment. |
| go_memstats_alloc_bytes | gauge | Number of bytes allocated and still in use. |
| go_memstats_alloc_bytes | counter | Total number of bytes allocated, even if freed. |
| go_memstats_buck_hash_sys_bytes | gauge | Number of bytes used by the profiling bucket hash table. |
| go_memstats_frees | counter | Total number of frees. |
| go_memstats_gc_cpu_fraction | gauge | The fraction of this program's available CPU time used by the GC since the program started. |
| go_memstats_gc_sys_bytes | gauge | Number of bytes used for garbage collection system metadata. |
| go_memstats_heap_alloc_bytes | gauge | Number of heap bytes allocated and still in use. |
| go_memstats_heap_idle_bytes | gauge | Number of heap bytes waiting to be used. |
| go_memstats_heap_inuse_bytes | gauge | Number of heap bytes that are in use. |
| go_memstats_heap_objects | gauge | Number of allocated objects. |
| go_memstats_heap_released_bytes | gauge | Number of heap bytes released to OS. |
| go_memstats_heap_sys_bytes | gauge | Number of heap bytes obtained from system. |
| go_memstats_last_gc_time_seconds | gauge | Number of seconds since 1970 of last garbage collection. |
| go_memstats_lookups | counter | Total number of pointer lookups. |
| go_memstats_mallocs | counter | Total number of mallocs. |
| go_memstats_mcache_inuse_bytes | gauge | Number of bytes in use by mcache structures. |
| go_memstats_mcache_sys_bytes | gauge | Number of bytes used for mcache structures obtained from system. |
| go_memstats_mspan_inuse_bytes | gauge | Number of bytes in use by mspan structures. |
| go_memstats_mspan_sys_bytes | gauge | Number of bytes used for mspan structures obtained from system. |
| go_memstats_next_gc_bytes | gauge | Number of heap bytes when next garbage collection will take place. |
| go_memstats_other_sys_bytes | gauge | Number of bytes used for other system allocations. |
| go_memstats_stack_inuse_bytes | gauge | Number of bytes in use by the stack allocator. |
| go_memstats_stack_sys_bytes | gauge | Number of bytes obtained from system for stack allocator. |
| go_memstats_sys_bytes | gauge | Number of bytes obtained from system. |
| go_threads | gauge | Number of OS threads created. |
| grpc_client_handled | counter | Total number of RPCs completed by the client, regardless of success or failure. |
| grpc_client_msg_received | counter | Total number of RPC stream messages received by the client. |
| grpc_client_msg_sent | counter | Total number of gRPC stream messages sent by the client. |
| grpc_client_started | counter | Total number of RPCs started on the client. |
| kube_apiserver_clusterip_allocator_allocated_ips | gauge | ALPHA Gauge measuring the number of allocated IPs for Services |
| kube_apiserver_clusterip_allocator_allocation | counter | ALPHA Number of Cluster IPs allocations |
| kube_apiserver_clusterip_allocator_available_ips | gauge | ALPHA Gauge measuring the number of available IPs for Services |
| kube_apiserver_pod_logs_pods_logs_backend_tls_failure | counter | ALPHA Total number of requests for pods/logs that failed due to kubelet server TLS verification |
| kube_apiserver_pod_logs_pods_logs_insecure_backend | counter | ALPHA Total number of requests for pods/logs sliced by usage type: enforce_tls, skip_tls_allowed, skip_tls_denied |
| kubernetes_build_info | gauge | ALPHA A metric with a constant '1' value labeled by major, minor, git version, git commit, git tree state, build date, Go version, and compiler from which Kubernetes was built, and platform on which it is running. |
| node_authorizer_graph_actions_duration_seconds | histogram | ALPHA Histogram of duration of graph actions in node authorizer. |
| pod_security_evaluations | counter | ALPHA Number of policy evaluations that occurred, not counting ignored or exempt requests. |
| pod_security_exemptions | counter | ALPHA Number of exempt requests, not counting ignored or out of scope requests. |
| process_cpu_seconds | counter | Total user and system CPU time spent in seconds. |
| process_max_fds | gauge | Maximum number of open file descriptors. |
| process_open_fds | gauge | Number of open file descriptors. |
| process_resident_memory_bytes | gauge | Resident memory size in bytes. |
| process_start_time_seconds | gauge | Start time of the process since unix epoch in seconds. |
| process_virtual_memory_bytes | gauge | Virtual memory size in bytes. |
| process_virtual_memory_max_bytes | gauge | Maximum amount of virtual memory available in bytes. |
| rest_client_exec_plugin_certificate_rotation_age | histogram | ALPHA Histogram of the number of seconds the last auth exec plugin client certificate lived before being rotated. If auth exec plugin client certificates are unused, histogram will contain no data. |
| rest_client_exec_plugin_ttl_seconds | gauge | ALPHA Gauge of the shortest TTL (time-to-live) of the client certificate(s) managed by the auth exec plugin. The value is in seconds until certificate expiry (negative if already expired). If auth exec plugins are unused or manage no TLS certificates, the value will be +INF. |
| rest_client_request_duration_seconds | histogram | ALPHA Request latency in seconds. Broken down by verb and URL. |
| rest_client_requests | counter | ALPHA Number of HTTP requests, partitioned by status code, method, and host. |
| serviceaccount_legacy_tokens | counter | ALPHA Cumulative legacy service account tokens used |
| serviceaccount_stale_tokens | counter | ALPHA Cumulative stale projected service account tokens used |
| serviceaccount_valid_tokens | counter | ALPHA Cumulative valid projected service account tokens used |
| watch_cache_capacity | gauge | ALPHA Total capacity of watch cache broken by resource type. |
| watch_cache_capacity_decrease | counter | ALPHA Total number of watch cache capacity decrease events broken by resource type. |
| watch_cache_capacity_increase | counter | ALPHA Total number of watch cache capacity increase events broken by resource type. |
| workqueue_adds | counter | ALPHA Total number of adds handled by workqueue |
| workqueue_depth | gauge | ALPHA Current depth of workqueue |
| workqueue_longest_running_processor_seconds | gauge | ALPHA How many seconds has the longest running processor for workqueue been running. |
| workqueue_queue_duration_seconds | histogram | ALPHA How long in seconds an item stays in workqueue before being requested. |
| workqueue_retries | counter | ALPHA Total number of retries handled by workqueue |
| workqueue_unfinished_work_seconds | gauge | ALPHA How many seconds of work has done that is in progress and hasn't been observed by work_duration. Large values indicate stuck threads. One can deduce the number of stuck threads by observing the rate at which this increases. |
| workqueue_work_duration_seconds | histogram | ALPHA How long in seconds processing an item from workqueue takes. |
<!-- markdownlint-enable line-length -->

## Kube-state-metrics (Kubernetes v1.23)

### Statistics

<!-- markdownlint-disable line-length -->
| Number of metrics families | Number of time series | Average number of labels per time series | Number of unique label pairs |
| -------------------------- | --------------------- | ---------------------------------------- | ---------------------------- |
| 220 | 28072 | 3.15 | 5043 |
<!-- markdownlint-enable line-length -->

### Metrics

<!-- markdownlint-disable line-length -->
| Name | Type, Unit | Description |
| ---- | ---------- | ----------- |
| kube_configmap_annotations | gauge | Kubernetes annotations converted to Prometheus labels. |
| kube_configmap_labels | gauge | Kubernetes labels converted to Prometheus labels. |
| kube_configmap_info | gauge | Information about configmap. |
| kube_configmap_created | gauge | Unix creation timestamp |
| kube_configmap_metadata_resource_version | gauge | Resource version representing a specific version of the configmap. |
| kube_cronjob_annotations | gauge | Kubernetes annotations converted to Prometheus labels. |
| kube_cronjob_labels | gauge | Kubernetes labels converted to Prometheus labels. |
| kube_cronjob_info | gauge | Info about cronjob. |
| kube_cronjob_created | gauge | Unix creation timestamp |
| kube_cronjob_status_active | gauge | Active holds pointers to currently running jobs. |
| kube_cronjob_status_last_schedule_time | gauge | LastScheduleTime keeps information of when was the last time the job was successfully scheduled. |
| kube_cronjob_spec_suspend | gauge | Suspend flag tells the controller to suspend subsequent executions. |
| kube_cronjob_spec_starting_deadline_seconds | gauge | Deadline in seconds for starting the job if it misses scheduled time for any reason. |
| kube_cronjob_next_schedule_time | gauge | Next time the cronjob should be scheduled. The time after lastScheduleTime, or after the cron job's creation time if it's never been scheduled. Use this to determine if the job is delayed. |
| kube_cronjob_metadata_resource_version | gauge | Resource version representing a specific version of the cronjob. |
| kube_cronjob_spec_successful_job_history_limit | gauge | Successful job history limit tells the controller how many completed jobs should be preserved. |
| kube_cronjob_spec_failed_job_history_limit | gauge | Failed job history limit tells the controller how many failed jobs should be preserved. |
| kube_daemonset_created | gauge | Unix creation timestamp |
| kube_daemonset_status_current_number_scheduled | gauge | The number of nodes running at least one daemon pod and are supposed to. |
| kube_daemonset_status_desired_number_scheduled | gauge | The number of nodes that should be running the daemon pod. |
| kube_daemonset_status_number_available | gauge | The number of nodes that should be running the daemon pod and have one or more of the daemon pod running and available |
| kube_daemonset_status_number_misscheduled | gauge | The number of nodes running a daemon pod but are not supposed to. |
| kube_daemonset_status_number_ready | gauge | The number of nodes that should be running the daemon pod and have one or more of the daemon pod running and ready. |
| kube_daemonset_status_number_unavailable | gauge | The number of nodes that should be running the daemon pod and have none of the daemon pod running and available |
| kube_daemonset_status_observed_generation | gauge | The most recent generation observed by the daemon set controller. |
| kube_daemonset_status_updated_number_scheduled | gauge | The total number of nodes that are running updated daemon pod |
| kube_daemonset_metadata_generation | gauge | Sequence number representing a specific generation of the desired state. |
| kube_daemonset_annotations | gauge | Kubernetes annotations converted to Prometheus labels. |
| kube_daemonset_labels | gauge | Kubernetes labels converted to Prometheus labels. |
| kube_deployment_created | gauge | Unix creation timestamp |
| kube_deployment_status_replicas | gauge | The number of replicas per deployment. |
| kube_deployment_status_replicas_ready | gauge | The number of ready replicas per deployment. |
| kube_deployment_status_replicas_available | gauge | The number of available replicas per deployment. |
| kube_deployment_status_replicas_unavailable | gauge | The number of unavailable replicas per deployment. |
| kube_deployment_status_replicas_updated | gauge | The number of updated replicas per deployment. |
| kube_deployment_status_observed_generation | gauge | The generation observed by the deployment controller. |
| kube_deployment_status_condition | gauge | The current status conditions of a deployment. |
| kube_deployment_spec_replicas | gauge | Number of desired pods for a deployment. |
| kube_deployment_spec_paused | gauge | Whether the deployment is paused and will not be processed by the deployment controller. |
| kube_deployment_spec_strategy_rollingupdate_max_unavailable | gauge | Maximum number of unavailable replicas during a rolling update of a deployment. |
| kube_deployment_spec_strategy_rollingupdate_max_surge | gauge | Maximum number of replicas that can be scheduled above the desired number of replicas during a rolling update of a deployment. |
| kube_deployment_metadata_generation | gauge | Sequence number representing a specific generation of the desired state. |
| kube_deployment_annotations | gauge | Kubernetes annotations converted to Prometheus labels. |
| kube_deployment_labels | gauge | Kubernetes labels converted to Prometheus labels. |
| kube_endpoint_info | gauge | Information about endpoint. |
| kube_endpoint_created | gauge | Unix creation timestamp |
| kube_endpoint_annotations | gauge | Kubernetes annotations converted to Prometheus labels. |
| kube_endpoint_labels | gauge | Kubernetes labels converted to Prometheus labels. |
| kube_endpoint_address_available | gauge | Number of addresses available in endpoint. |
| kube_endpoint_address_not_ready | gauge | Number of addresses not ready in endpoint |
| kube_endpoint_ports | gauge | Information about the Endpoint ports. |
| kube_horizontalpodautoscaler_info | gauge | Information about this autoscaler. |
| kube_horizontalpodautoscaler_metadata_generation | gauge | The generation observed by the HorizontalPodAutoscaler controller. |
| kube_horizontalpodautoscaler_spec_max_replicas | gauge | Upper limit for the number of pods that can be set by the autoscaler; cannot be smaller than MinReplicas. |
| kube_horizontalpodautoscaler_spec_min_replicas | gauge | Lower limit for the number of pods that can be set by the autoscaler, default 1. |
| kube_horizontalpodautoscaler_spec_target_metric | gauge | The metric specifications used by this autoscaler when calculating the desired replica count. |
| kube_horizontalpodautoscaler_status_current_replicas | gauge | Current number of replicas of pods managed by this autoscaler. |
| kube_horizontalpodautoscaler_status_desired_replicas | gauge | Desired number of replicas of pods managed by this autoscaler. |
| kube_horizontalpodautoscaler_annotations | gauge | Kubernetes annotations converted to Prometheus labels. |
| kube_horizontalpodautoscaler_labels | gauge | Kubernetes labels converted to Prometheus labels. |
| kube_horizontalpodautoscaler_status_condition | gauge | The condition of this autoscaler. |
| kube_ingress_info | gauge | Information about ingress. |
| kube_ingress_annotations | gauge | Kubernetes annotations converted to Prometheus labels. |
| kube_ingress_labels | gauge | Kubernetes labels converted to Prometheus labels. |
| kube_ingress_created | gauge | Unix creation timestamp |
| kube_ingress_metadata_resource_version | gauge | Resource version representing a specific version of ingress. |
| kube_ingress_path | gauge | Ingress host, paths and backend service information. |
| kube_ingress_tls | gauge | Ingress TLS host and secret information. |
| kube_job_annotations | gauge | Kubernetes annotations converted to Prometheus labels. |
| kube_job_labels | gauge | Kubernetes labels converted to Prometheus labels. |
| kube_job_info | gauge | Information about job. |
| kube_job_created | gauge | Unix creation timestamp |
| kube_job_spec_parallelism | gauge | The maximum desired number of pods the job should run at any given time. |
| kube_job_spec_completions | gauge | The desired number of successfully finished pods the job should be run with. |
| kube_job_spec_active_deadline_seconds | gauge | The duration in seconds relative to the startTime that the job may be active before the system tries to terminate it. |
| kube_job_status_succeeded | gauge | The number of pods which reached Phase Succeeded. |
| kube_job_status_failed | gauge | The number of pods which reached Phase Failed and the reason for failure. |
| kube_job_status_active | gauge | The number of actively running pods. |
| kube_job_complete | gauge | The job has completed its execution. |
| kube_job_failed | gauge | The job has failed its execution. |
| kube_job_status_start_time | gauge | StartTime represents time when the job was acknowledged by the Job Manager. |
| kube_job_status_completion_time | gauge | CompletionTime represents time when the job was completed. |
| kube_job_owner | gauge | Information about the Job's owner. |
| kube_limitrange | gauge | Information about limit range. |
| kube_limitrange_created | gauge | Unix creation timestamp |
| kube_namespace_created | gauge | Unix creation timestamp |
| kube_namespace_annotations | gauge | Kubernetes annotations converted to Prometheus labels. |
| kube_namespace_labels | gauge | Kubernetes labels converted to Prometheus labels. |
| kube_namespace_status_phase | gauge | kubernetes namespace status phase. |
| kube_namespace_status_condition | gauge | The condition of a namespace. |
| kube_networkpolicy_created | gauge | Unix creation timestamp of network policy |
| kube_networkpolicy_annotations | gauge | Kubernetes annotations converted to Prometheus labels. |
| kube_networkpolicy_labels | gauge | Kubernetes labels converted to Prometheus labels. |
| kube_networkpolicy_spec_ingress_rules | gauge | Number of ingress rules on the networkpolicy |
| kube_networkpolicy_spec_egress_rules | gauge | Number of egress rules on the networkpolicy |
| kube_node_created | gauge | Unix creation timestamp |
| kube_node_info | gauge | Information about a cluster node. |
| kube_node_annotations | gauge | Kubernetes annotations converted to Prometheus labels. |
| kube_node_labels | gauge | Kubernetes labels converted to Prometheus labels. |
| kube_node_role | gauge | The role of a cluster node. |
| kube_node_spec_taint | gauge | The taint of a cluster node. |
| kube_node_spec_unschedulable | gauge | Whether a node can schedule new pods. |
| kube_node_status_allocatable | gauge | The allocatable for different resources of a node that are available for scheduling. |
| kube_node_status_capacity | gauge | The capacity for different resources of a node. |
| kube_node_status_condition | gauge | The condition of a cluster node. |
| kube_persistentvolumeclaim_labels | gauge | Kubernetes labels converted to Prometheus labels. |
| kube_persistentvolumeclaim_annotations | gauge | Kubernetes annotations converted to Prometheus labels. |
| kube_persistentvolumeclaim_info | gauge | Information about persistent volume claim. |
| kube_persistentvolumeclaim_status_phase | gauge | The phase the persistent volume claim is currently in. |
| kube_persistentvolumeclaim_resource_requests_storage_bytes | gauge | The capacity of storage requested by the persistent volume claim. |
| kube_persistentvolumeclaim_access_mode | gauge | The access mode(s) specified by the persistent volume claim. |
| kube_persistentvolumeclaim_status_condition | gauge | Information about status of different conditions of persistent volume claim. |
| kube_persistentvolume_claim_ref | gauge | Information about the Persistent Volume Claim Reference. |
| kube_persistentvolume_annotations | gauge | Kubernetes annotations converted to Prometheus labels. |
| kube_persistentvolume_labels | gauge | Kubernetes labels converted to Prometheus labels. |
| kube_persistentvolume_status_phase | gauge | The phase indicates if a volume is available, bound to a claim, or released by a claim. |
| kube_persistentvolume_info | gauge | Information about persistentvolume. |
| kube_persistentvolume_capacity_bytes | gauge | Persistentvolume capacity in bytes. |
| kube_poddisruptionbudget_annotations | gauge | Kubernetes annotations converted to Prometheus labels. |
| kube_poddisruptionbudget_labels | gauge | Kubernetes labels converted to Prometheus labels. |
| kube_poddisruptionbudget_created | gauge | Unix creation timestamp |
| kube_poddisruptionbudget_status_current_healthy | gauge | Current number of healthy pods |
| kube_poddisruptionbudget_status_desired_healthy | gauge | Minimum desired number of healthy pods |
| kube_poddisruptionbudget_status_pod_disruptions_allowed | gauge | Number of pod disruptions that are currently allowed |
| kube_poddisruptionbudget_status_expected_pods | gauge | Total number of pods counted by this disruption budget |
| kube_poddisruptionbudget_status_observed_generation | gauge | Most recent generation observed when updating this PDB status |
| kube_pod_completion_time | gauge | Completion time in unix timestamp for a pod. |
| kube_pod_container_info | gauge | Information about a container in a pod. |
| kube_pod_container_resource_limits | gauge | The number of requested limit resource by a container. |
| kube_pod_container_resource_requests | gauge | The number of requested request resource by a container. |
| kube_pod_container_state_started | gauge | Start time in unix timestamp for a pod container. |
| kube_pod_container_status_last_terminated_reason | gauge | Describes the last reason the container was in terminated state. |
| kube_pod_container_status_ready | gauge | Describes whether the containers readiness check succeeded. |
| kube_pod_container_status_restarts | counter | The number of container restarts per container. |
| kube_pod_container_status_running | gauge | Describes whether the container is currently in running state. |
| kube_pod_container_status_terminated | gauge | Describes whether the container is currently in terminated state. |
| kube_pod_container_status_terminated_reason | gauge | Describes the reason the container is currently in terminated state. |
| kube_pod_container_status_waiting | gauge | Describes whether the container is currently in waiting state. |
| kube_pod_container_status_waiting_reason | gauge | Describes the reason the container is currently in waiting state. |
| kube_pod_created | gauge | Unix creation timestamp |
| kube_pod_deletion_timestamp | gauge | Unix deletion timestamp |
| kube_pod_info | gauge | Information about pod. |
| kube_pod_init_container_info | gauge | Information about an init container in a pod. |
| kube_pod_init_container_resource_limits | gauge | The number of requested limit resource by an init container. |
| kube_pod_init_container_resource_requests | gauge | The number of requested request resource by an init container. |
| kube_pod_init_container_status_last_terminated_reason | gauge | Describes the last reason the init container was in terminated state. |
| kube_pod_init_container_status_ready | gauge | Describes whether the init containers readiness check succeeded. |
| kube_pod_init_container_status_restarts | counter | The number of restarts for the init container. |
| kube_pod_init_container_status_running | gauge | Describes whether the init container is currently in running state. |
| kube_pod_init_container_status_terminated | gauge | Describes whether the init container is currently in terminated state. |
| kube_pod_init_container_status_terminated_reason | gauge | Describes the reason the init container is currently in terminated state. |
| kube_pod_init_container_status_waiting | gauge | Describes whether the init container is currently in waiting state. |
| kube_pod_init_container_status_waiting_reason | gauge | Describes the reason the init container is currently in waiting state. |
| kube_pod_annotations | gauge | Kubernetes annotations converted to Prometheus labels. |
| kube_pod_labels | gauge | Kubernetes labels converted to Prometheus labels. |
| kube_pod_overhead_cpu_cores | gauge | The pod overhead in regards to cpu cores associated with running a pod. |
| kube_pod_overhead_memory_bytes | gauge | The pod overhead in regards to memory associated with running a pod. |
| kube_pod_owner | gauge | Information about the Pod's owner. |
| kube_pod_restart_policy | gauge | Describes the restart policy in use by this pod. |
| kube_pod_runtimeclass_name_info | gauge | The runtimeclass associated with the pod. |
| kube_pod_spec_volumes_persistentvolumeclaims_info | gauge | Information about persistentvolumeclaim volumes in a pod. |
| kube_pod_spec_volumes_persistentvolumeclaims_readonly | gauge | Describes whether a persistentvolumeclaim is mounted read only. |
| kube_pod_start_time | gauge | Start time in unix timestamp for a pod. |
| kube_pod_status_phase | gauge | The pods current phase. |
| kube_pod_status_ready | gauge | Describes whether the pod is ready to serve requests. |
| kube_pod_status_reason | gauge | The pod status reasons |
| kube_pod_status_scheduled | gauge | Describes the status of the scheduling process for the pod. |
| kube_pod_status_scheduled_time | gauge | Unix timestamp when pod moved into scheduled status |
| kube_pod_status_unschedulable | gauge | Describes the unschedulable status for the pod. |
| kube_replicaset_created | gauge | Unix creation timestamp |
| kube_replicaset_status_replicas | gauge | The number of replicas per ReplicaSet. |
| kube_replicaset_status_fully_labeled_replicas | gauge | The number of fully labeled replicas per ReplicaSet. |
| kube_replicaset_status_ready_replicas | gauge | The number of ready replicas per ReplicaSet. |
| kube_replicaset_status_observed_generation | gauge | The generation observed by the ReplicaSet controller. |
| kube_replicaset_spec_replicas | gauge | Number of desired pods for a ReplicaSet. |
| kube_replicaset_metadata_generation | gauge | Sequence number representing a specific generation of the desired state. |
| kube_replicaset_owner | gauge | Information about the ReplicaSet's owner. |
| kube_replicaset_annotations | gauge | Kubernetes annotations converted to Prometheus labels. |
| kube_replicaset_labels | gauge | Kubernetes labels converted to Prometheus labels. |
| kube_replicationcontroller_created | gauge | Unix creation timestamp |
| kube_replicationcontroller_status_replicas | gauge | The number of replicas per ReplicationController. |
| kube_replicationcontroller_status_fully_labeled_replicas | gauge | The number of fully labeled replicas per ReplicationController. |
| kube_replicationcontroller_status_ready_replicas | gauge | The number of ready replicas per ReplicationController. |
| kube_replicationcontroller_status_available_replicas | gauge | The number of available replicas per ReplicationController. |
| kube_replicationcontroller_status_observed_generation | gauge | The generation observed by the ReplicationController controller. |
| kube_replicationcontroller_spec_replicas | gauge | Number of desired pods for a ReplicationController. |
| kube_replicationcontroller_metadata_generation | gauge | Sequence number representing a specific generation of the desired state. |
| kube_replicationcontroller_owner | gauge | Information about the ReplicationController's owner. |
| kube_resourcequota_created | gauge | Unix creation timestamp |
| kube_resourcequota | gauge | Information about resource quota. |
| kube_secret_info | gauge | Information about secret. |
| kube_secret_type | gauge | Type about secret. |
| kube_secret_annotations | gauge | Kubernetes annotations converted to Prometheus labels. |
| kube_secret_labels | gauge | Kubernetes labels converted to Prometheus labels. |
| kube_secret_created | gauge | Unix creation timestamp |
| kube_secret_metadata_resource_version | gauge | Resource version representing a specific version of secret. |
| kube_service_info | gauge | Information about service. |
| kube_service_created | gauge | Unix creation timestamp |
| kube_service_spec_type | gauge | Type about service. |
| kube_service_annotations | gauge | Kubernetes annotations converted to Prometheus labels. |
| kube_service_labels | gauge | Kubernetes labels converted to Prometheus labels. |
| kube_service_spec_external_ip | gauge | Service external ips. One series for each ip |
| kube_service_status_load_balancer_ingress | gauge | Service load balancer ingress status |
| kube_statefulset_created | gauge | Unix creation timestamp |
| kube_statefulset_status_replicas | gauge | The number of replicas per StatefulSet. |
| kube_statefulset_status_replicas_available | gauge | The number of available replicas per StatefulSet. |
| kube_statefulset_status_replicas_current | gauge | The number of current replicas per StatefulSet. |
| kube_statefulset_status_replicas_ready | gauge | The number of ready replicas per StatefulSet. |
| kube_statefulset_status_replicas_updated | gauge | The number of updated replicas per StatefulSet. |
| kube_statefulset_status_observed_generation | gauge | The generation observed by the StatefulSet controller. |
| kube_statefulset_replicas | gauge | Number of desired pods for a StatefulSet. |
| kube_statefulset_metadata_generation | gauge | Sequence number representing a specific generation of the desired state for the StatefulSet. |
| kube_statefulset_annotations | gauge | Kubernetes annotations converted to Prometheus labels. |
| kube_statefulset_labels | gauge | Kubernetes labels converted to Prometheus labels. |
| kube_statefulset_status_current_revision | gauge | Indicates the version of the StatefulSet used to generate Pods in the sequence [0,currentReplicas). |
| kube_statefulset_status_update_revision | gauge | Indicates the version of the StatefulSet used to generate Pods in the sequence [replicas-updatedReplicas,replicas) |
| kube_storageclass_info | gauge | Information about storageclass. |
| kube_storageclass_created | gauge | Unix creation timestamp |
| kube_storageclass_annotations | gauge | Kubernetes annotations converted to Prometheus labels. |
| kube_storageclass_labels | gauge | Kubernetes labels converted to Prometheus labels. |
<!-- markdownlint-enable line-length -->

## Kubelet (Kubernetes v1.23)

### Statistics

<!-- markdownlint-disable line-length -->
| Number of metrics families | Number of time series | Average number of labels per time series | Number of unique label pairs |
| -------------------------- | --------------------- | ---------------------------------------- | ---------------------------- |
| 104 | 2373 | 2.99 | 432 |
<!-- markdownlint-enable line-length -->

### Metrics

<!-- markdownlint-disable line-length -->
| Name | Type, Unit | Description |
| ---- | ---------- | ----------- |
| apiserver_audit_event | counter | ALPHA Counter of audit events generated and sent to the audit backend. |
| apiserver_audit_requests_rejected | counter | ALPHA Counter of apiserver requests rejected due to an error in audit logging backend. |
| apiserver_client_certificate_expiration_seconds | histogram | ALPHA Distribution of the remaining lifetime on the certificate used to authenticate a request. |
| apiserver_delegated_authn_request_duration_seconds | histogram | ALPHA Request latency in seconds. Broken down by status code. |
| apiserver_delegated_authn_request | counter | ALPHA Number of HTTP requests partitioned by status code. |
| apiserver_delegated_authz_request_duration_seconds | histogram | ALPHA Request latency in seconds. Broken down by status code. |
| apiserver_delegated_authz_request | counter | ALPHA Number of HTTP requests partitioned by status code. |
| apiserver_envelope_encryption_dek_cache_fill_percent | gauge | ALPHA Percent of the cache slots currently occupied by cached DEKs. |
| apiserver_storage_data_key_generation_duration_seconds | histogram | ALPHA Latencies in seconds of data encryption key(DEK) generation operations. |
| apiserver_storage_data_key_generation_failures | counter | ALPHA Total number of failed data encryption key(DEK) generation operations. |
| apiserver_storage_envelope_transformation_cache_misses | counter | ALPHA Total number of cache misses while accessing key decryption key(KEK). |
| apiserver_webhooks_x509_missing_san | counter | ALPHA Counts the number of requests to servers missing SAN extension in their serving certificate OR the number of connection failures due to the lack of x509 certificate SAN extension missing (either/or, based on the runtime environment) |
| authentication_token_cache_active_fetch_count | gauge | ALPHA |
| authentication_token_cache_fetch | counter | ALPHA |
| authentication_token_cache_request_duration_seconds | histogram | ALPHA |
| authentication_token_cache_request | counter | ALPHA |
| get_token_count | counter | ALPHA Counter of total Token() requests to the alternate token source |
| get_token_fail_count | counter | ALPHA Counter of failed Token() requests to the alternate token source |
| go_gc_duration_seconds | summary | A summary of the pause duration of garbage collection cycles. |
| go_goroutines | gauge | Number of goroutines that currently exist. |
| go_info | gauge | Information about the Go environment. |
| go_memstats_alloc_bytes | gauge | Number of bytes allocated and still in use. |
| go_memstats_alloc_bytes | counter | Total number of bytes allocated, even if freed. |
| go_memstats_buck_hash_sys_bytes | gauge | Number of bytes used by the profiling bucket hash table. |
| go_memstats_frees | counter | Total number of frees. |
| go_memstats_gc_cpu_fraction | gauge | The fraction of this program's available CPU time used by the GC since the program started. |
| go_memstats_gc_sys_bytes | gauge | Number of bytes used for garbage collection system metadata. |
| go_memstats_heap_alloc_bytes | gauge | Number of heap bytes allocated and still in use. |
| go_memstats_heap_idle_bytes | gauge | Number of heap bytes waiting to be used. |
| go_memstats_heap_inuse_bytes | gauge | Number of heap bytes that are in use. |
| go_memstats_heap_objects | gauge | Number of allocated objects. |
| go_memstats_heap_released_bytes | gauge | Number of heap bytes released to OS. |
| go_memstats_heap_sys_bytes | gauge | Number of heap bytes obtained from system. |
| go_memstats_last_gc_time_seconds | gauge | Number of seconds since 1970 of last garbage collection. |
| go_memstats_lookups | counter | Total number of pointer lookups. |
| go_memstats_mallocs | counter | Total number of mallocs. |
| go_memstats_mcache_inuse_bytes | gauge | Number of bytes in use by mcache structures. |
| go_memstats_mcache_sys_bytes | gauge | Number of bytes used for mcache structures obtained from system. |
| go_memstats_mspan_inuse_bytes | gauge | Number of bytes in use by mspan structures. |
| go_memstats_mspan_sys_bytes | gauge | Number of bytes used for mspan structures obtained from system. |
| go_memstats_next_gc_bytes | gauge | Number of heap bytes when next garbage collection will take place. |
| go_memstats_other_sys_bytes | gauge | Number of bytes used for other system allocations. |
| go_memstats_stack_inuse_bytes | gauge | Number of bytes in use by the stack allocator. |
| go_memstats_stack_sys_bytes | gauge | Number of bytes obtained from system for stack allocator. |
| go_memstats_sys_bytes | gauge | Number of bytes obtained from system. |
| go_threads | gauge | Number of OS threads created. |
| kubelet_certificate_manager_client_expiration_renew_errors | counter | ALPHA Counter of certificate renewal errors. |
| kubelet_certificate_manager_client_ttl_seconds | gauge | ALPHA Gauge of the TTL (time-to-live) of the Kubelet's client certificate. The value is in seconds until certificate expiry (negative if already expired). If client certificate is invalid or unused, the value will be +INF. |
| kubelet_cgroup_manager_duration_seconds | histogram | ALPHA Duration in seconds for cgroup manager operations. Broken down by method. |
| kubelet_container_log_filesystem_used_bytes | gauge | ALPHA Bytes used by the container's logs on the filesystem. |
| kubelet_containers_per_pod_count | histogram | ALPHA The number of containers per pod. |
| kubelet_eviction_stats_age_seconds | histogram | ALPHA Time between when stats are collected, and when pod is evicted based on those stats by eviction signal |
| kubelet_evictions | counter | ALPHA Cumulative number of pod evictions by eviction signal |
| kubelet_http_inflight_requests | gauge | ALPHA Number of the inflight http requests |
| kubelet_http_requests_duration_seconds | histogram | ALPHA Duration in seconds to serve http requests |
| kubelet_http_requests | counter | ALPHA Number of the http requests received since the server started |
| kubelet_managed_ephemeral_containers | gauge | ALPHA Current number of ephemeral containers in pods managed by this kubelet. Ephemeral containers will be ignored if disabled by the EphemeralContainers feature gate, and this number will be 0. |
| kubelet_node_name | gauge | ALPHA The node's name. The count is always 1. |
| kubelet_pleg_discard_events | counter | ALPHA The number of discard events in PLEG. |
| kubelet_pleg_last_seen_seconds | gauge | ALPHA Timestamp in seconds when PLEG was last seen active. |
| kubelet_pleg_relist_duration_seconds | histogram | ALPHA Duration in seconds for relisting pods in PLEG. |
| kubelet_pleg_relist_interval_seconds | histogram | ALPHA Interval in seconds between relisting in PLEG. |
| kubelet_pod_start_duration_seconds | histogram | ALPHA Duration in seconds for a single pod to go from pending to running. |
| kubelet_pod_worker_duration_seconds | histogram | ALPHA Duration in seconds to sync a single pod. Broken down by operation type: create, update, or sync |
| kubelet_pod_worker_start_duration_seconds | histogram | ALPHA Duration in seconds from seeing a pod to starting a worker. |
| kubelet_run_podsandbox_duration_seconds | histogram | ALPHA Duration in seconds of the run_podsandbox operations. Broken down by RuntimeClass.Handler. |
| kubelet_run_podsandbox_errors | counter | ALPHA Cumulative number of the run_podsandbox operation errors by RuntimeClass.Handler. |
| kubelet_running_containers | gauge | ALPHA Number of containers currently running |
| kubelet_running_pods | gauge | ALPHA Number of pods that have a running pod sandbox |
| kubelet_runtime_operations_duration_seconds | histogram | ALPHA Duration in seconds of runtime operations. Broken down by operation type. |
| kubelet_runtime_operations_errors | counter | ALPHA Cumulative number of runtime operation errors by operation type. |
| kubelet_runtime_operations | counter | ALPHA Cumulative number of runtime operations by operation type. |
| kubelet_started_containers_errors | counter | ALPHA Cumulative number of errors when starting containers |
| kubelet_started_containers | counter | ALPHA Cumulative number of containers started |
| kubelet_started_pods_errors | counter | ALPHA Cumulative number of errors when starting pods |
| kubelet_started_pods | counter | ALPHA Cumulative number of pods started |
| kubelet_volume_stats_available_bytes | gauge | ALPHA Number of available bytes in the volume |
| kubelet_volume_stats_capacity_bytes | gauge | ALPHA Capacity in bytes of the volume |
| kubelet_volume_stats_inodes | gauge | ALPHA Maximum number of inodes in the volume |
| kubelet_volume_stats_inodes_free | gauge | ALPHA Number of free inodes in the volume |
| kubelet_volume_stats_inodes_used | gauge | ALPHA Number of used inodes in the volume |
| kubelet_volume_stats_used_bytes | gauge | ALPHA Number of used bytes in the volume |
| kubernetes_build_info | gauge | ALPHA A metric with a constant '1' value labeled by major, minor, git version, git commit, git tree state, build date, Go version, and compiler from which Kubernetes was built, and platform on which it is running. |
| process_cpu_seconds | counter | Total user and system CPU time spent in seconds. |
| process_max_fds | gauge | Maximum number of open file descriptors. |
| process_open_fds | gauge | Number of open file descriptors. |
| process_resident_memory_bytes | gauge | Resident memory size in bytes. |
| process_start_time_seconds | gauge | Start time of the process since unix epoch in seconds. |
| process_virtual_memory_bytes | gauge | Virtual memory size in bytes. |
| process_virtual_memory_max_bytes | gauge | Maximum amount of virtual memory available in bytes. |
| rest_client_exec_plugin_certificate_rotation_age | histogram | ALPHA Histogram of the number of seconds the last auth exec plugin client certificate lived before being rotated. If auth exec plugin client certificates are unused, histogram will contain no data. |
| rest_client_exec_plugin_ttl_seconds | gauge | ALPHA Gauge of the shortest TTL (time-to-live) of the client certificate(s) managed by the auth exec plugin. The value is in seconds until certificate expiry (negative if already expired). If auth exec plugins are unused or manage no TLS certificates, the value will be +INF. |
| rest_client_rate_limiter_duration_seconds | histogram | ALPHA Client side rate limiter latency in seconds. Broken down by verb and URL. |
| rest_client_request_duration_seconds | histogram | ALPHA Request latency in seconds. Broken down by verb and URL. |
| rest_client_requests | counter | ALPHA Number of HTTP requests, partitioned by status code, method, and host. |
| storage_operation_duration_seconds | histogram | ALPHA Storage operation duration |
| volume_manager_total_volumes | gauge | ALPHA Number of volumes in Volume Manager |
| workqueue_adds | counter | ALPHA Total number of adds handled by workqueue |
| workqueue_depth | gauge | ALPHA Current depth of workqueue |
| workqueue_longest_running_processor_seconds | gauge | ALPHA How many seconds has the longest running processor for workqueue been running. |
| workqueue_queue_duration_seconds | histogram | ALPHA How long in seconds an item stays in workqueue before being requested. |
| workqueue_retries | counter | ALPHA Total number of retries handled by workqueue |
| workqueue_unfinished_work_seconds | gauge | ALPHA How many seconds of work has done that is in progress and hasn't been observed by work_duration. Large values indicate stuck threads. One can deduce the number of stuck threads by observing the rate at which this increases. |
| workqueue_work_duration_seconds | histogram | ALPHA How long in seconds processing an item from workqueue takes. |
<!-- markdownlint-enable line-length -->

## cAdvisor (Kubernetes v1.23)

### Statistics

<!-- markdownlint-disable line-length -->
| Number of metrics families | Number of time series | Average number of labels per time series | Number of unique label pairs |
| -------------------------- | --------------------- | ---------------------------------------- | ---------------------------- |
| 66 | 13057 | 6.77 | 679 |
<!-- markdownlint-enable line-length -->

### Metrics

<!-- markdownlint-disable line-length -->
| Name | Type, Unit | Description |
| ---- | ---------- | ----------- |
| cadvisor_version_info | gauge | A metric with a constant '1' value labeled by kernel version, OS version, docker version, cadvisor version & cadvisor revision. |
| container_blkio_device_usage | counter | Blkio Device bytes usage |
| container_cpu_cfs_periods | counter | Number of elapsed enforcement period intervals. |
| container_cpu_cfs_throttled_periods | counter | Number of throttled period intervals. |
| container_cpu_cfs_throttled_seconds | counter | Total time duration the container has been throttled. |
| container_cpu_load_average_10s | gauge | Value of container cpu load average over the last 10 seconds. |
| container_cpu_system_seconds | counter | Cumulative system cpu time consumed in seconds. |
| container_cpu_usage_seconds | counter | Cumulative cpu time consumed in seconds. |
| container_cpu_user_seconds | counter | Cumulative user cpu time consumed in seconds. |
| container_file_descriptors | gauge | Number of open file descriptors for the container. |
| container_fs_inodes_free | gauge | Number of available Inodes |
| container_fs_inodes_total | gauge | Number of Inodes |
| container_fs_io_current | gauge | Number of I/Os currently in progress |
| container_fs_io_time_seconds | counter | Cumulative count of seconds spent doing I/Os |
| container_fs_io_time_weighted_seconds | counter | Cumulative weighted I/O time in seconds |
| container_fs_limit_bytes | gauge | Number of bytes that can be consumed by the container on this filesystem. |
| container_fs_read_seconds | counter | Cumulative count of seconds spent reading |
| container_fs_reads_bytes | counter | Cumulative count of bytes read |
| container_fs_reads_merged | counter | Cumulative count of reads merged |
| container_fs_reads | counter | Cumulative count of reads completed |
| container_fs_sector_reads | counter | Cumulative count of sector reads completed |
| container_fs_sector_writes | counter | Cumulative count of sector writes completed |
| container_fs_usage_bytes | gauge | Number of bytes that are consumed by the container on this filesystem. |
| container_fs_write_seconds | counter | Cumulative count of seconds spent writing |
| container_fs_writes_bytes | counter | Cumulative count of bytes written |
| container_fs_writes_merged | counter | Cumulative count of writes merged |
| container_fs_writes | counter | Cumulative count of writes completed |
| container_last_seen | gauge | Last time a container was seen by the exporter |
| container_memory_cache | gauge | Number of bytes of page cache memory. |
| container_memory_failcnt | counter | Number of memory usage hits limits |
| container_memory_failures | counter | Cumulative count of memory allocation failures. |
| container_memory_mapped_file | gauge | Size of memory mapped files in bytes. |
| container_memory_max_usage_bytes | gauge | Maximum memory usage recorded in bytes |
| container_memory_rss | gauge | Size of RSS in bytes. |
| container_memory_swap | gauge | Container swap usage in bytes. |
| container_memory_usage_bytes | gauge | Current memory usage in bytes, including all memory regardless of when it was accessed |
| container_memory_working_set_bytes | gauge | Current working set in bytes. |
| container_network_receive_bytes | counter | Cumulative count of bytes received |
| container_network_receive_errors | counter | Cumulative count of errors encountered while receiving |
| container_network_receive_packets_dropped | counter | Cumulative count of packets dropped while receiving |
| container_network_receive_packets | counter | Cumulative count of packets received |
| container_network_transmit_bytes | counter | Cumulative count of bytes transmitted |
| container_network_transmit_errors | counter | Cumulative count of errors encountered while transmitting |
| container_network_transmit_packets_dropped | counter | Cumulative count of packets dropped while transmitting |
| container_network_transmit_packets | counter | Cumulative count of packets transmitted |
| container_processes | gauge | Number of processes running inside the container. |
| container_scrape_error | gauge | 1 if there was an error while getting container metrics, 0 otherwise |
| container_sockets | gauge | Number of open sockets for the container. |
| container_spec_cpu_period | gauge | CPU period of the container. |
| container_spec_cpu_quota | gauge | CPU quota of the container. |
| container_spec_cpu_shares | gauge | CPU share of the container. |
| container_spec_memory_limit_bytes | gauge | Memory limit for the container. |
| container_spec_memory_reservation_limit_bytes | gauge | Memory reservation limit for the container. |
| container_spec_memory_swap_limit_bytes | gauge | Memory swap limit for the container. |
| container_start_time_seconds | gauge | Start time of the container since unix epoch in seconds. |
| container_tasks_state | gauge | Number of tasks in given state |
| container_threads | gauge | Number of threads running inside the container |
| container_threads_max | gauge | Maximum number of threads allowed inside the container, infinity if value is zero |
| container_ulimits_soft | gauge | Soft ulimit values for the container root process. Unlimited if -1, except priority and nice |
| machine_cpu_cores | gauge | Number of logical CPU cores. |
| machine_cpu_physical_cores | gauge | Number of physical CPU cores. |
| machine_cpu_sockets | gauge | Number of CPU sockets. |
| machine_memory_bytes | gauge | Amount of memory installed on the machine. |
| machine_nvm_avg_power_budget_watts | gauge | NVM power budget. |
| machine_nvm_capacity | gauge | NVM capacity value labeled by NVM mode (memory mode or app direct mode). |
| machine_scrape_error | gauge | 1 if there was an error while getting machine metrics, 0 otherwise. |
<!-- markdownlint-enable line-length -->

## NGINX ingress (Kubernetes v1.23)

### Statistics

<!-- markdownlint-disable line-length -->
| Number of metrics families | Number of time series | Average number of labels per time series | Number of unique label pairs |
| -------------------------- | --------------------- | ---------------------------------------- | ---------------------------- |
| 52 | 64 | 1.17 | 21 |
<!-- markdownlint-enable line-length -->

### Metrics

<!-- markdownlint-disable line-length -->
| Name | Type, Unit | Description |
| ---- | ---------- | ----------- |
| go_gc_duration_seconds | summary | A summary of the pause duration of garbage collection cycles. |
| go_goroutines | gauge | Number of goroutines that currently exist. |
| go_info | gauge | Information about the Go environment. |
| go_memstats_alloc_bytes | gauge | Number of bytes allocated and still in use. |
| go_memstats_alloc_bytes | counter | Total number of bytes allocated, even if freed. |
| go_memstats_buck_hash_sys_bytes | gauge | Number of bytes used by the profiling bucket hash table. |
| go_memstats_frees | counter | Total number of frees. |
| go_memstats_gc_cpu_fraction | gauge | The fraction of this program's available CPU time used by the GC since the program started. |
| go_memstats_gc_sys_bytes | gauge | Number of bytes used for garbage collection system metadata. |
| go_memstats_heap_alloc_bytes | gauge | Number of heap bytes allocated and still in use. |
| go_memstats_heap_idle_bytes | gauge | Number of heap bytes waiting to be used. |
| go_memstats_heap_inuse_bytes | gauge | Number of heap bytes that are in use. |
| go_memstats_heap_objects | gauge | Number of allocated objects. |
| go_memstats_heap_released_bytes | gauge | Number of heap bytes released to OS. |
| go_memstats_heap_sys_bytes | gauge | Number of heap bytes obtained from system. |
| go_memstats_last_gc_time_seconds | gauge | Number of seconds since 1970 of last garbage collection. |
| go_memstats_lookups | counter | Total number of pointer lookups. |
| go_memstats_mallocs | counter | Total number of mallocs. |
| go_memstats_mcache_inuse_bytes | gauge | Number of bytes in use by mcache structures. |
| go_memstats_mcache_sys_bytes | gauge | Number of bytes used for mcache structures obtained from system. |
| go_memstats_mspan_inuse_bytes | gauge | Number of bytes in use by mspan structures. |
| go_memstats_mspan_sys_bytes | gauge | Number of bytes used for mspan structures obtained from system. |
| go_memstats_next_gc_bytes | gauge | Number of heap bytes when next garbage collection will take place. |
| go_memstats_other_sys_bytes | gauge | Number of bytes used for other system allocations. |
| go_memstats_stack_inuse_bytes | gauge | Number of bytes in use by the stack allocator. |
| go_memstats_stack_sys_bytes | gauge | Number of bytes obtained from system for stack allocator. |
| go_memstats_sys_bytes | gauge | Number of bytes obtained from system. |
| go_threads | gauge | Number of OS threads created. |
| nginx_ingress_controller_build_info | gauge | A metric with a constant '1' labeled with information about the build. |
| nginx_ingress_controller_config_hash | gauge | Running configuration hash actually running |
| nginx_ingress_controller_config_last_reload_successful | gauge | Whether the last configuration reload attempt was successful |
| nginx_ingress_controller_config_last_reload_successful_timestamp_seconds | gauge | Timestamp of the last successful configuration reload. |
| nginx_ingress_controller_nginx_process_connections | gauge | current number of client connections with state {active, reading, writing, waiting} |
| nginx_ingress_controller_nginx_process_connections | counter | total number of connections with state {accepted, handled} |
| nginx_ingress_controller_nginx_process_cpu_seconds | counter | Cpu usage in seconds |
| nginx_ingress_controller_nginx_process_num_procs | gauge | number of processes |
| nginx_ingress_controller_nginx_process_oldest_start_time_seconds | gauge | start time in seconds since 1970/01/01 |
| nginx_ingress_controller_nginx_process_read_bytes | counter | number of bytes read |
| nginx_ingress_controller_nginx_process_requests | counter | total number of client requests |
| nginx_ingress_controller_nginx_process_resident_memory_bytes | gauge | number of bytes of memory in use |
| nginx_ingress_controller_nginx_process_virtual_memory_bytes | gauge | number of bytes of memory in use |
| nginx_ingress_controller_nginx_process_write_bytes | counter | number of bytes written |
| nginx_ingress_controller_success | counter | Cumulative number of Ingress controller reload operations |
| process_cpu_seconds | counter | Total user and system CPU time spent in seconds. |
| process_max_fds | gauge | Maximum number of open file descriptors. |
| process_open_fds | gauge | Number of open file descriptors. |
| process_resident_memory_bytes | gauge | Resident memory size in bytes. |
| process_start_time_seconds | gauge | Start time of the process since unix epoch in seconds. |
| process_virtual_memory_bytes | gauge | Virtual memory size in bytes. |
| process_virtual_memory_max_bytes | gauge | Maximum amount of virtual memory available in bytes. |
| promhttp_metric_handler_requests_in_flight | gauge | Current number of scrapes being served. |
| promhttp_metric_handler_requests | counter | Total number of scrapes by HTTP status code. |
<!-- markdownlint-enable line-length -->

## Node-exporter (Kubernetes v1.23)

### Statistics

<!-- markdownlint-disable line-length -->
| Number of metrics families | Number of time series | Average number of labels per time series | Number of unique label pairs |
| -------------------------- | --------------------- | ---------------------------------------- | ---------------------------- |
| 320 | 3461 | 1.55 | 542 |
<!-- markdownlint-enable line-length -->

### Metrics

<!-- markdownlint-disable line-length -->
| Name | Type, Unit | Description |
| ---- | ---------- | ----------- |
| go_gc_duration_seconds | summary | A summary of the pause duration of garbage collection cycles. |
| go_goroutines | gauge | Number of goroutines that currently exist. |
| go_info | gauge | Information about the Go environment. |
| go_memstats_alloc_bytes | gauge | Number of bytes allocated and still in use. |
| go_memstats_alloc_bytes | counter | Total number of bytes allocated, even if freed. |
| go_memstats_buck_hash_sys_bytes | gauge | Number of bytes used by the profiling bucket hash table. |
| go_memstats_frees | counter | Total number of frees. |
| go_memstats_gc_cpu_fraction | gauge | The fraction of this program's available CPU time used by the GC since the program started. |
| go_memstats_gc_sys_bytes | gauge | Number of bytes used for garbage collection system metadata. |
| go_memstats_heap_alloc_bytes | gauge | Number of heap bytes allocated and still in use. |
| go_memstats_heap_idle_bytes | gauge | Number of heap bytes waiting to be used. |
| go_memstats_heap_inuse_bytes | gauge | Number of heap bytes that are in use. |
| go_memstats_heap_objects | gauge | Number of allocated objects. |
| go_memstats_heap_released_bytes | gauge | Number of heap bytes released to OS. |
| go_memstats_heap_sys_bytes | gauge | Number of heap bytes obtained from system. |
| go_memstats_last_gc_time_seconds | gauge | Number of seconds since 1970 of last garbage collection. |
| go_memstats_lookups | counter | Total number of pointer lookups. |
| go_memstats_mallocs | counter | Total number of mallocs. |
| go_memstats_mcache_inuse_bytes | gauge | Number of bytes in use by mcache structures. |
| go_memstats_mcache_sys_bytes | gauge | Number of bytes used for mcache structures obtained from system. |
| go_memstats_mspan_inuse_bytes | gauge | Number of bytes in use by mspan structures. |
| go_memstats_mspan_sys_bytes | gauge | Number of bytes used for mspan structures obtained from system. |
| go_memstats_next_gc_bytes | gauge | Number of heap bytes when next garbage collection will take place. |
| go_memstats_other_sys_bytes | gauge | Number of bytes used for other system allocations. |
| go_memstats_stack_inuse_bytes | gauge | Number of bytes in use by the stack allocator. |
| go_memstats_stack_sys_bytes | gauge | Number of bytes obtained from system for stack allocator. |
| go_memstats_sys_bytes | gauge | Number of bytes obtained from system. |
| go_threads | gauge | Number of OS threads created. |
| node_arp_entries | gauge | ARP entries by device |
| node_boot_time_seconds | gauge | Node boot time, in unixtime. |
| node_context_switches | counter | Total number of context switches. |
| node_cooling_device_cur_state | gauge | Current throttle state of the cooling device |
| node_cooling_device_max_state | gauge | Maximum throttle state of the cooling device |
| node_cpu_guest_seconds | counter | Seconds the CPUs spent in guests (VMs) for each mode. |
| node_cpu_seconds | counter | Seconds the CPUs spent in each mode. |
| node_disk_info | gauge | Info of /sys/block/<block_device>. |
| node_disk_io_now | gauge | The number of I/Os currently in progress. |
| node_disk_io_time_seconds | counter | Total seconds spent doing I/Os. |
| node_disk_io_time_weighted_seconds | counter | The weighted # of seconds spent doing I/Os. |
| node_disk_read_bytes | counter | The total number of bytes read successfully. |
| node_disk_read_time_seconds | counter | The total number of seconds spent by all reads. |
| node_disk_reads_completed | counter | The total number of reads completed successfully. |
| node_disk_reads_merged | counter | The total number of reads merged. |
| node_disk_write_time_seconds | counter | This is the total number of seconds spent by all writes. |
| node_disk_writes_completed | counter | The total number of writes completed successfully. |
| node_disk_writes_merged | counter | The number of writes merged. |
| node_disk_written_bytes | counter | The total number of bytes written successfully. |
| node_dmi_info | gauge | A metric with a constant '1' value labeled by bios_date, bios_release, bios_vendor, bios_version, board_asset_tag, board_name, board_serial, board_vendor, board_version, chassis_asset_tag, chassis_serial, chassis_vendor, chassis_version, product_family, product_name, product_serial, product_sku, product_uuid, product_version, system_vendor if provided by DMI. |
| node_entropy_available_bits | gauge | Bits of available entropy. |
| node_entropy_pool_size_bits | gauge | Bits of entropy pool. |
| node_exporter_build_info | gauge | A metric with a constant '1' value labeled by version, revision, branch, and goversion from which node_exporter was built. |
| node_filefd_allocated | gauge | File descriptor statistics: allocated. |
| node_filefd_maximum | gauge | File descriptor statistics: maximum. |
| node_filesystem_avail_bytes | gauge | Filesystem space available to non-root users in bytes. |
| node_filesystem_device_error | gauge | Whether an error occurred while getting statistics for the given device. |
| node_filesystem_files | gauge | Filesystem total file nodes. |
| node_filesystem_files_free | gauge | Filesystem total free file nodes. |
| node_filesystem_free_bytes | gauge | Filesystem free space in bytes. |
| node_filesystem_readonly | gauge | Filesystem read-only status. |
| node_filesystem_size_bytes | gauge | Filesystem size in bytes. |
| node_forks | counter | Total number of forks. |
| node_intr | counter | Total number of interrupts serviced. |
| node_ipvs_connections | counter | The total number of connections made. |
| node_ipvs_incoming_bytes | counter | The total amount of incoming data. |
| node_ipvs_incoming_packets | counter | The total number of incoming packets. |
| node_ipvs_outgoing_bytes | counter | The total amount of outgoing data. |
| node_ipvs_outgoing_packets | counter | The total number of outgoing packets. |
| node_load1 | gauge | 1m load average. |
| node_load15 | gauge | 15m load average. |
| node_load5 | gauge | 5m load average. |
| node_memory_Active_anon_bytes | gauge | Memory information field Active_anon_bytes. |
| node_memory_Active_bytes | gauge | Memory information field Active_bytes. |
| node_memory_Active_file_bytes | gauge | Memory information field Active_file_bytes. |
| node_memory_AnonHugePages_bytes | gauge | Memory information field AnonHugePages_bytes. |
| node_memory_AnonPages_bytes | gauge | Memory information field AnonPages_bytes. |
| node_memory_Bounce_bytes | gauge | Memory information field Bounce_bytes. |
| node_memory_Buffers_bytes | gauge | Memory information field Buffers_bytes. |
| node_memory_Cached_bytes | gauge | Memory information field Cached_bytes. |
| node_memory_CmaFree_bytes | gauge | Memory information field CmaFree_bytes. |
| node_memory_CmaTotal_bytes | gauge | Memory information field CmaTotal_bytes. |
| node_memory_CommitLimit_bytes | gauge | Memory information field CommitLimit_bytes. |
| node_memory_Committed_AS_bytes | gauge | Memory information field Committed_AS_bytes. |
| node_memory_DirectMap2M_bytes | gauge | Memory information field DirectMap2M_bytes. |
| node_memory_DirectMap4k_bytes | gauge | Memory information field DirectMap4k_bytes. |
| node_memory_Dirty_bytes | gauge | Memory information field Dirty_bytes. |
| node_memory_HardwareCorrupted_bytes | gauge | Memory information field HardwareCorrupted_bytes. |
| node_memory_HugePages_Free | gauge | Memory information field HugePages_Free. |
| node_memory_HugePages_Rsvd | gauge | Memory information field HugePages_Rsvd. |
| node_memory_HugePages_Surp | gauge | Memory information field HugePages_Surp. |
| node_memory_HugePages_Total | gauge | Memory information field HugePages_Total. |
| node_memory_Hugepagesize_bytes | gauge | Memory information field Hugepagesize_bytes. |
| node_memory_Inactive_anon_bytes | gauge | Memory information field Inactive_anon_bytes. |
| node_memory_Inactive_bytes | gauge | Memory information field Inactive_bytes. |
| node_memory_Inactive_file_bytes | gauge | Memory information field Inactive_file_bytes. |
| node_memory_KernelStack_bytes | gauge | Memory information field KernelStack_bytes. |
| node_memory_Mapped_bytes | gauge | Memory information field Mapped_bytes. |
| node_memory_MemAvailable_bytes | gauge | Memory information field MemAvailable_bytes. |
| node_memory_MemFree_bytes | gauge | Memory information field MemFree_bytes. |
| node_memory_MemTotal_bytes | gauge | Memory information field MemTotal_bytes. |
| node_memory_Mlocked_bytes | gauge | Memory information field Mlocked_bytes. |
| node_memory_NFS_Unstable_bytes | gauge | Memory information field NFS_Unstable_bytes. |
| node_memory_PageTables_bytes | gauge | Memory information field PageTables_bytes. |
| node_memory_SReclaimable_bytes | gauge | Memory information field SReclaimable_bytes. |
| node_memory_SUnreclaim_bytes | gauge | Memory information field SUnreclaim_bytes. |
| node_memory_Shmem_bytes | gauge | Memory information field Shmem_bytes. |
| node_memory_Slab_bytes | gauge | Memory information field Slab_bytes. |
| node_memory_SwapCached_bytes | gauge | Memory information field SwapCached_bytes. |
| node_memory_SwapFree_bytes | gauge | Memory information field SwapFree_bytes. |
| node_memory_SwapTotal_bytes | gauge | Memory information field SwapTotal_bytes. |
| node_memory_Unevictable_bytes | gauge | Memory information field Unevictable_bytes. |
| node_memory_VmallocChunk_bytes | gauge | Memory information field VmallocChunk_bytes. |
| node_memory_VmallocTotal_bytes | gauge | Memory information field VmallocTotal_bytes. |
| node_memory_VmallocUsed_bytes | gauge | Memory information field VmallocUsed_bytes. |
| node_memory_WritebackTmp_bytes | gauge | Memory information field WritebackTmp_bytes. |
| node_memory_Writeback_bytes | gauge | Memory information field Writeback_bytes. |
| node_netstat_Icmp6_InErrors | unknown | Statistic Icmp6InErrors. |
| node_netstat_Icmp6_InMsgs | unknown | Statistic Icmp6InMsgs. |
| node_netstat_Icmp6_OutMsgs | unknown | Statistic Icmp6OutMsgs. |
| node_netstat_Icmp_InErrors | unknown | Statistic IcmpInErrors. |
| node_netstat_Icmp_InMsgs | unknown | Statistic IcmpInMsgs. |
| node_netstat_Icmp_OutMsgs | unknown | Statistic IcmpOutMsgs. |
| node_netstat_Ip6_InOctets | unknown | Statistic Ip6InOctets. |
| node_netstat_Ip6_OutOctets | unknown | Statistic Ip6OutOctets. |
| node_netstat_IpExt_InOctets | unknown | Statistic IpExtInOctets. |
| node_netstat_IpExt_OutOctets | unknown | Statistic IpExtOutOctets. |
| node_netstat_Ip_Forwarding | unknown | Statistic IpForwarding. |
| node_netstat_TcpExt_ListenDrops | unknown | Statistic TcpExtListenDrops. |
| node_netstat_TcpExt_ListenOverflows | unknown | Statistic TcpExtListenOverflows. |
| node_netstat_TcpExt_SyncookiesFailed | unknown | Statistic TcpExtSyncookiesFailed. |
| node_netstat_TcpExt_SyncookiesRecv | unknown | Statistic TcpExtSyncookiesRecv. |
| node_netstat_TcpExt_SyncookiesSent | unknown | Statistic TcpExtSyncookiesSent. |
| node_netstat_TcpExt_TCPSynRetrans | unknown | Statistic TcpExtTCPSynRetrans. |
| node_netstat_TcpExt_TCPTimeouts | unknown | Statistic TcpExtTCPTimeouts. |
| node_netstat_Tcp_ActiveOpens | unknown | Statistic TcpActiveOpens. |
| node_netstat_Tcp_CurrEstab | unknown | Statistic TcpCurrEstab. |
| node_netstat_Tcp_InErrs | unknown | Statistic TcpInErrs. |
| node_netstat_Tcp_InSegs | unknown | Statistic TcpInSegs. |
| node_netstat_Tcp_OutRsts | unknown | Statistic TcpOutRsts. |
| node_netstat_Tcp_OutSegs | unknown | Statistic TcpOutSegs. |
| node_netstat_Tcp_PassiveOpens | unknown | Statistic TcpPassiveOpens. |
| node_netstat_Tcp_RetransSegs | unknown | Statistic TcpRetransSegs. |
| node_netstat_Udp6_InDatagrams | unknown | Statistic Udp6InDatagrams. |
| node_netstat_Udp6_InErrors | unknown | Statistic Udp6InErrors. |
| node_netstat_Udp6_NoPorts | unknown | Statistic Udp6NoPorts. |
| node_netstat_Udp6_OutDatagrams | unknown | Statistic Udp6OutDatagrams. |
| node_netstat_Udp6_RcvbufErrors | unknown | Statistic Udp6RcvbufErrors. |
| node_netstat_Udp6_SndbufErrors | unknown | Statistic Udp6SndbufErrors. |
| node_netstat_UdpLite6_InErrors | unknown | Statistic UdpLite6InErrors. |
| node_netstat_UdpLite_InErrors | unknown | Statistic UdpLiteInErrors. |
| node_netstat_Udp_InDatagrams | unknown | Statistic UdpInDatagrams. |
| node_netstat_Udp_InErrors | unknown | Statistic UdpInErrors. |
| node_netstat_Udp_NoPorts | unknown | Statistic UdpNoPorts. |
| node_netstat_Udp_OutDatagrams | unknown | Statistic UdpOutDatagrams. |
| node_netstat_Udp_RcvbufErrors | unknown | Statistic UdpRcvbufErrors. |
| node_netstat_Udp_SndbufErrors | unknown | Statistic UdpSndbufErrors. |
| node_network_address_assign_type | gauge | address_assign_type value of /sys/class/net/<iface>. |
| node_network_carrier | gauge | carrier value of /sys/class/net/<iface>. |
| node_network_carrier_changes | counter | carrier_changes_total value of /sys/class/net/<iface>. |
| node_network_device_id | gauge | device_id value of /sys/class/net/<iface>. |
| node_network_dormant | gauge | dormant value of /sys/class/net/<iface>. |
| node_network_flags | gauge | flags value of /sys/class/net/<iface>. |
| node_network_iface_id | gauge | iface_id value of /sys/class/net/<iface>. |
| node_network_iface_link | gauge | iface_link value of /sys/class/net/<iface>. |
| node_network_iface_link_mode | gauge | iface_link_mode value of /sys/class/net/<iface>. |
| node_network_info | gauge | Non-numeric data from /sys/class/net/<iface>, value is always 1. |
| node_network_mtu_bytes | gauge | mtu_bytes value of /sys/class/net/<iface>. |
| node_network_net_dev_group | gauge | net_dev_group value of /sys/class/net/<iface>. |
| node_network_protocol_type | gauge | protocol_type value of /sys/class/net/<iface>. |
| node_network_receive_bytes | counter | Network device statistic receive_bytes. |
| node_network_receive_compressed | counter | Network device statistic receive_compressed. |
| node_network_receive_drop | counter | Network device statistic receive_drop. |
| node_network_receive_errs | counter | Network device statistic receive_errs. |
| node_network_receive_fifo | counter | Network device statistic receive_fifo. |
| node_network_receive_frame | counter | Network device statistic receive_frame. |
| node_network_receive_multicast | counter | Network device statistic receive_multicast. |
| node_network_receive_packets | counter | Network device statistic receive_packets. |
| node_network_speed_bytes | gauge | speed_bytes value of /sys/class/net/<iface>. |
| node_network_transmit_bytes | counter | Network device statistic transmit_bytes. |
| node_network_transmit_carrier | counter | Network device statistic transmit_carrier. |
| node_network_transmit_colls | counter | Network device statistic transmit_colls. |
| node_network_transmit_compressed | counter | Network device statistic transmit_compressed. |
| node_network_transmit_drop | counter | Network device statistic transmit_drop. |
| node_network_transmit_errs | counter | Network device statistic transmit_errs. |
| node_network_transmit_fifo | counter | Network device statistic transmit_fifo. |
| node_network_transmit_packets | counter | Network device statistic transmit_packets. |
| node_network_transmit_queue_length | gauge | transmit_queue_length value of /sys/class/net/<iface>. |
| node_network_up | gauge | Value is 1 if operstate is 'up', 0 otherwise. |
| node_nf_conntrack_entries | gauge | Number of currently allocated flow entries for connection tracking. |
| node_nf_conntrack_entries_limit | gauge | Maximum size of connection tracking table. |
| node_nf_conntrack_stat_drop | gauge | Number of packets dropped due to conntrack failure. |
| node_nf_conntrack_stat_early_drop | gauge | Number of dropped conntrack entries to make room for new ones, if maximum table size was reached. |
| node_nf_conntrack_stat_found | gauge | Number of searched entries which were successful. |
| node_nf_conntrack_stat_ignore | gauge | Number of packets seen which are already connected to a conntrack entry. |
| node_nf_conntrack_stat_insert | gauge | Number of entries inserted into the list. |
| node_nf_conntrack_stat_insert_failed | gauge | Number of entries for which list insertion was attempted but failed. |
| node_nf_conntrack_stat_invalid | gauge | Number of packets seen which can not be tracked. |
| node_nf_conntrack_stat_search_restart | gauge | Number of conntrack table lookups which had to be restarted due to hashtable resizes. |
| node_nfs_connections | counter | Total number of NFSd TCP connections. |
| node_nfs_packets | counter | Total NFSd network packets (sent+received) by protocol type. |
| node_nfs_requests | counter | Number of NFS procedures invoked. |
| node_nfs_rpc_authentication_refreshes | counter | Number of RPC authentication refreshes performed. |
| node_nfs_rpc_retransmissions | counter | Number of RPC transmissions performed. |
| node_nfs_rpcs | counter | Total number of RPCs performed. |
| node_os_info | gauge | A metric with a constant '1' value labeled by build_id, id, id_like, image_id, image_version, name, pretty_name, variant, variant_id, version, version_codename, version_id. |
| node_os_version | gauge | Metric containing the major.minor part of the OS version. |
| node_processes_max_processes | gauge | Number of max PIDs limit |
| node_processes_max_threads | gauge | Limit of threads in the system |
| node_processes_pids | gauge | Number of PIDs |
| node_processes_state | gauge | Number of processes in each state. |
| node_processes_threads | gauge | Allocated threads in system |
| node_processes_threads_state | gauge | Number of threads in each state. |
| node_procs_blocked | gauge | Number of processes blocked waiting for I/O to complete. |
| node_procs_running | gauge | Number of processes in runnable state. |
| node_schedstat_running_seconds | counter | Number of seconds CPU spent running a process. |
| node_schedstat_timeslices | counter | Number of timeslices executed by CPU. |
| node_schedstat_waiting_seconds | counter | Number of seconds spent by processing waiting for this CPU. |
| node_scrape_collector_duration_seconds | gauge | node_exporter: Duration of a collector scrape. |
| node_scrape_collector_success | gauge | node_exporter: Whether a collector succeeded. |
| node_sockstat_FRAG6_inuse | gauge | Number of FRAG6 sockets in state inuse. |
| node_sockstat_FRAG6_memory | gauge | Number of FRAG6 sockets in state memory. |
| node_sockstat_FRAG_inuse | gauge | Number of FRAG sockets in state inuse. |
| node_sockstat_FRAG_memory | gauge | Number of FRAG sockets in state memory. |
| node_sockstat_RAW6_inuse | gauge | Number of RAW6 sockets in state inuse. |
| node_sockstat_RAW_inuse | gauge | Number of RAW sockets in state inuse. |
| node_sockstat_TCP6_inuse | gauge | Number of TCP6 sockets in state inuse. |
| node_sockstat_TCP_alloc | gauge | Number of TCP sockets in state alloc. |
| node_sockstat_TCP_inuse | gauge | Number of TCP sockets in state inuse. |
| node_sockstat_TCP_mem | gauge | Number of TCP sockets in state mem. |
| node_sockstat_TCP_mem_bytes | gauge | Number of TCP sockets in state mem_bytes. |
| node_sockstat_TCP_orphan | gauge | Number of TCP sockets in state orphan. |
| node_sockstat_TCP_tw | gauge | Number of TCP sockets in state tw. |
| node_sockstat_UDP6_inuse | gauge | Number of UDP6 sockets in state inuse. |
| node_sockstat_UDPLITE6_inuse | gauge | Number of UDPLITE6 sockets in state inuse. |
| node_sockstat_UDPLITE_inuse | gauge | Number of UDPLITE sockets in state inuse. |
| node_sockstat_UDP_inuse | gauge | Number of UDP sockets in state inuse. |
| node_sockstat_UDP_mem | gauge | Number of UDP sockets in state mem. |
| node_sockstat_UDP_mem_bytes | gauge | Number of UDP sockets in state mem_bytes. |
| node_sockstat_sockets_used | gauge | Number of IPv4 sockets in use. |
| node_softnet_dropped | counter | Number of dropped packets |
| node_softnet_processed | counter | Number of processed packets |
| node_softnet_times_squeezed | counter | Number of times processing packets ran out of quota |
| node_textfile_scrape_error | gauge | 1 if there was an error opening or reading a file, 0 otherwise |
| node_time_clocksource_available_info | gauge | Available clocksources read from '/sys/devices/system/clocksource'. |
| node_time_clocksource_current_info | gauge | Current clocksource read from '/sys/devices/system/clocksource'. |
| node_time_seconds | gauge | System time in seconds since epoch (1970). |
| node_time_zone_offset_seconds | gauge | System time zone offset in seconds. |
| node_timex_estimated_error_seconds | gauge | Estimated error in seconds. |
| node_timex_frequency_adjustment_ratio | gauge | Local clock frequency adjustment. |
| node_timex_loop_time_constant | gauge | Phase-locked loop time constant. |
| node_timex_maxerror_seconds | gauge | Maximum error in seconds. |
| node_timex_offset_seconds | gauge | Time offset in between local system and reference clock. |
| node_timex_pps_calibration | counter | Pulse per second count of calibration intervals. |
| node_timex_pps_error | counter | Pulse per second count of calibration errors. |
| node_timex_pps_frequency_hertz | gauge | Pulse per second frequency. |
| node_timex_pps_jitter_seconds | gauge | Pulse per second jitter. |
| node_timex_pps_jitter | counter | Pulse per second count of jitter limit exceeded events. |
| node_timex_pps_shift_seconds | gauge | Pulse per second interval duration. |
| node_timex_pps_stability_exceeded | counter | Pulse per second count of stability limit exceeded events. |
| node_timex_pps_stability_hertz | gauge | Pulse per second stability, average of recent frequency changes. |
| node_timex_status | gauge | Value of the status array bits. |
| node_timex_sync_status | gauge | Is clock synchronized to a reliable server (1 = yes, 0 = no). |
| node_timex_tai_offset_seconds | gauge | International Atomic Time (TAI) offset. |
| node_timex_tick_seconds | gauge | Seconds between clock ticks. |
| node_udp_queues | gauge | Number of allocated memory in the kernel for UDP datagrams in bytes. |
| node_uname_info | gauge | Labeled system information as provided by the uname system call. |
| node_vmstat_pgfault | unknown | /proc/vmstat information field pgfault. |
| node_vmstat_pgmajfault | unknown | /proc/vmstat information field pgmajfault. |
| node_vmstat_pgpgin | unknown | /proc/vmstat information field pgpgin. |
| node_vmstat_pgpgout | unknown | /proc/vmstat information field pgpgout. |
| node_vmstat_pswpin | unknown | /proc/vmstat information field pswpin. |
| node_vmstat_pswpout | unknown | /proc/vmstat information field pswpout. |
| node_xfs_allocation_btree_compares | counter | Number of allocation B-tree compares for a filesystem. |
| node_xfs_allocation_btree_lookups | counter | Number of allocation B-tree lookups for a filesystem. |
| node_xfs_allocation_btree_records_deleted | counter | Number of allocation B-tree records deleted for a filesystem. |
| node_xfs_allocation_btree_records_inserted | counter | Number of allocation B-tree records inserted for a filesystem. |
| node_xfs_block_map_btree_compares | counter | Number of block map B-tree compares for a filesystem. |
| node_xfs_block_map_btree_lookups | counter | Number of block map B-tree lookups for a filesystem. |
| node_xfs_block_map_btree_records_deleted | counter | Number of block map B-tree records deleted for a filesystem. |
| node_xfs_block_map_btree_records_inserted | counter | Number of block map B-tree records inserted for a filesystem. |
| node_xfs_block_mapping_extent_list_compares | counter | Number of extent list compares for a filesystem. |
| node_xfs_block_mapping_extent_list_deletions | counter | Number of extent list deletions for a filesystem. |
| node_xfs_block_mapping_extent_list_insertions | counter | Number of extent list insertions for a filesystem. |
| node_xfs_block_mapping_extent_list_lookups | counter | Number of extent list lookups for a filesystem. |
| node_xfs_block_mapping_reads | counter | Number of block map for read operations for a filesystem. |
| node_xfs_block_mapping_unmaps | counter | Number of block unmaps (deletes) for a filesystem. |
| node_xfs_block_mapping_writes | counter | Number of block map for write operations for a filesystem. |
| node_xfs_directory_operation_create | counter | Number of times a new directory entry was created for a filesystem. |
| node_xfs_directory_operation_getdents | counter | Number of times the directory getdents operation was performed for a filesystem. |
| node_xfs_directory_operation_lookup | counter | Number of file name directory lookups which miss the operating systems directory name lookup cache. |
| node_xfs_directory_operation_remove | counter | Number of times an existing directory entry was created for a filesystem. |
| node_xfs_extent_allocation_blocks_allocated | counter | Number of blocks allocated for a filesystem. |
| node_xfs_extent_allocation_blocks_freed | counter | Number of blocks freed for a filesystem. |
| node_xfs_extent_allocation_extents_allocated | counter | Number of extents allocated for a filesystem. |
| node_xfs_extent_allocation_extents_freed | counter | Number of extents freed for a filesystem. |
| node_xfs_inode_operation_attempts | counter | Number of times the OS looked for an XFS inode in the inode cache. |
| node_xfs_inode_operation_attribute_changes | counter | Number of times the OS explicitly changed the attributes of an XFS inode. |
| node_xfs_inode_operation_duplicates | counter | Number of times the OS tried to add a missing XFS inode to the inode cache, but found it had already been added by another process. |
| node_xfs_inode_operation_found | counter | Number of times the OS looked for and found an XFS inode in the inode cache. |
| node_xfs_inode_operation_missed | counter | Number of times the OS looked for an XFS inode in the cache, but did not find it. |
| node_xfs_inode_operation_reclaims | counter | Number of times the OS reclaimed an XFS inode from the inode cache to free memory for another purpose. |
| node_xfs_inode_operation_recycled | counter | Number of times the OS found an XFS inode in the cache, but could not use it as it was being recycled. |
| node_xfs_read_calls | counter | Number of read(2) system calls made to files in a filesystem. |
| node_xfs_vnode_active | counter | Number of vnodes not on free lists for a filesystem. |
| node_xfs_vnode_allocate | counter | Number of times vn_alloc called for a filesystem. |
| node_xfs_vnode_get | counter | Number of times vn_get called for a filesystem. |
| node_xfs_vnode_hold | counter | Number of times vn_hold called for a filesystem. |
| node_xfs_vnode_reclaim | counter | Number of times vn_reclaim called for a filesystem. |
| node_xfs_vnode_release | counter | Number of times vn_rele called for a filesystem. |
| node_xfs_vnode_remove | counter | Number of times vn_remove called for a filesystem. |
| node_xfs_write_calls | counter | Number of write(2) system calls made to files in a filesystem. |
| process_cpu_seconds | counter | Total user and system CPU time spent in seconds. |
| process_max_fds | gauge | Maximum number of open file descriptors. |
| process_open_fds | gauge | Number of open file descriptors. |
| process_resident_memory_bytes | gauge | Resident memory size in bytes. |
| process_start_time_seconds | gauge | Start time of the process since unix epoch in seconds. |
| process_virtual_memory_bytes | gauge | Virtual memory size in bytes. |
| process_virtual_memory_max_bytes | gauge | Maximum amount of virtual memory available in bytes. |
| promhttp_metric_handler_errors | counter | Total number of internal errors encountered by the promhttp metric handler. |
| promhttp_metric_handler_requests_in_flight | gauge | Current number of scrapes being served. |
| promhttp_metric_handler_requests | counter | Total number of scrapes by HTTP status code. |
<!-- markdownlint-enable line-length -->

## Prometheus (Kubernetes v1.23)

### Statistics

<!-- markdownlint-disable line-length -->
| Number of metrics families | Number of time series | Average number of labels per time series | Number of unique label pairs |
| -------------------------- | --------------------- | ---------------------------------------- | ---------------------------- |
| 219 | 1542 | 1.27 | 286 |
<!-- markdownlint-enable line-length -->

### Metrics

<!-- markdownlint-disable line-length -->
| Name | Type, Unit | Description |
| ---- | ---------- | ----------- |
| go_gc_cycles_automatic_gc_cycles | counter | Count of completed GC cycles generated by the Go runtime. |
| go_gc_cycles_forced_gc_cycles | counter | Count of completed GC cycles forced by the application. |
| go_gc_cycles_total_gc_cycles | counter | Count of all completed GC cycles. |
| go_gc_duration_seconds | summary | A summary of the pause duration of garbage collection cycles. |
| go_gc_heap_allocs_by_size_bytes_total | histogram | Distribution of heap allocations by approximate size. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks. |
| go_gc_heap_allocs_bytes | counter | Cumulative sum of memory allocated to the heap by the application. |
| go_gc_heap_allocs_objects | counter | Cumulative count of heap allocations triggered by the application. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks. |
| go_gc_heap_frees_by_size_bytes_total | histogram | Distribution of freed heap allocations by approximate size. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks. |
| go_gc_heap_frees_bytes | counter | Cumulative sum of heap memory freed by the garbage collector. |
| go_gc_heap_frees_objects | counter | Cumulative count of heap allocations whose storage was freed by the garbage collector. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks. |
| go_gc_heap_goal_bytes | gauge | Heap size target for the end of the GC cycle. |
| go_gc_heap_objects_objects | gauge | Number of objects, live or unswept, occupying heap memory. |
| go_gc_heap_tiny_allocs_objects | counter | Count of small allocations that are packed together into blocks. These allocations are counted separately from other allocations because each individual allocation is not tracked by the runtime, only their block. Each block is already accounted for in allocs-by-size and frees-by-size. |
| go_gc_pauses_seconds_total | histogram | Distribution individual GC-related stop-the-world pause latencies. |
| go_goroutines | gauge | Number of goroutines that currently exist. |
| go_info | gauge | Information about the Go environment. |
| go_memory_classes_heap_free_bytes | gauge | Memory that is completely free and eligible to be returned to the underlying system, but has not been. This metric is the runtime's estimate of free address space that is backed by physical memory. |
| go_memory_classes_heap_objects_bytes | gauge | Memory occupied by live objects and dead objects that have not yet been marked free by the garbage collector. |
| go_memory_classes_heap_released_bytes | gauge | Memory that is completely free and has been returned to the underlying system. This metric is the runtime's estimate of free address space that is still mapped into the process, but is not backed by physical memory. |
| go_memory_classes_heap_stacks_bytes | gauge | Memory allocated from the heap that is reserved for stack space, whether or not it is currently in-use. |
| go_memory_classes_heap_unused_bytes | gauge | Memory that is reserved for heap objects but is not currently used to hold heap objects. |
| go_memory_classes_metadata_mcache_free_bytes | gauge | Memory that is reserved for runtime mcache structures, but not in-use. |
| go_memory_classes_metadata_mcache_inuse_bytes | gauge | Memory that is occupied by runtime mcache structures that are currently being used. |
| go_memory_classes_metadata_mspan_free_bytes | gauge | Memory that is reserved for runtime mspan structures, but not in-use. |
| go_memory_classes_metadata_mspan_inuse_bytes | gauge | Memory that is occupied by runtime mspan structures that are currently being used. |
| go_memory_classes_metadata_other_bytes | gauge | Memory that is reserved for or used to hold runtime metadata. |
| go_memory_classes_os_stacks_bytes | gauge | Stack memory allocated by the underlying operating system. |
| go_memory_classes_other_bytes | gauge | Memory used by execution trace buffers, structures for debugging the runtime, finalizer and profiler specials, and more. |
| go_memory_classes_profiling_buckets_bytes | gauge | Memory that is used by the stack trace hash map used for profiling. |
| go_memory_classes_total_bytes | gauge | All memory mapped by the Go runtime into the current process as read-write. Note that this does not include memory mapped by code called via cgo or via the syscall package. Sum of all metrics in /memory/classes. |
| go_memstats_alloc_bytes | gauge | Number of bytes allocated and still in use. |
| go_memstats_alloc_bytes | counter | Total number of bytes allocated, even if freed. |
| go_memstats_buck_hash_sys_bytes | gauge | Number of bytes used by the profiling bucket hash table. |
| go_memstats_frees | counter | Total number of frees. |
| go_memstats_gc_cpu_fraction | gauge | The fraction of this program's available CPU time used by the GC since the program started. |
| go_memstats_gc_sys_bytes | gauge | Number of bytes used for garbage collection system metadata. |
| go_memstats_heap_alloc_bytes | gauge | Number of heap bytes allocated and still in use. |
| go_memstats_heap_idle_bytes | gauge | Number of heap bytes waiting to be used. |
| go_memstats_heap_inuse_bytes | gauge | Number of heap bytes that are in use. |
| go_memstats_heap_objects | gauge | Number of allocated objects. |
| go_memstats_heap_released_bytes | gauge | Number of heap bytes released to OS. |
| go_memstats_heap_sys_bytes | gauge | Number of heap bytes obtained from system. |
| go_memstats_last_gc_time_seconds | gauge | Number of seconds since 1970 of last garbage collection. |
| go_memstats_lookups | counter | Total number of pointer lookups. |
| go_memstats_mallocs | counter | Total number of mallocs. |
| go_memstats_mcache_inuse_bytes | gauge | Number of bytes in use by mcache structures. |
| go_memstats_mcache_sys_bytes | gauge | Number of bytes used for mcache structures obtained from system. |
| go_memstats_mspan_inuse_bytes | gauge | Number of bytes in use by mspan structures. |
| go_memstats_mspan_sys_bytes | gauge | Number of bytes used for mspan structures obtained from system. |
| go_memstats_next_gc_bytes | gauge | Number of heap bytes when next garbage collection will take place. |
| go_memstats_other_sys_bytes | gauge | Number of bytes used for other system allocations. |
| go_memstats_stack_inuse_bytes | gauge | Number of bytes in use by the stack allocator. |
| go_memstats_stack_sys_bytes | gauge | Number of bytes obtained from system for stack allocator. |
| go_memstats_sys_bytes | gauge | Number of bytes obtained from system. |
| go_sched_goroutines_goroutines | gauge | Count of live goroutines. |
| go_sched_latencies_seconds | histogram | Distribution of the time goroutines have spent in the scheduler in a runnable state before actually running. |
| go_threads | gauge | Number of OS threads created. |
| net_conntrack_dialer_conn_attempted | counter | Total number of connections attempted by the given dialer a given name. |
| net_conntrack_dialer_conn_closed | counter | Total number of connections closed which originated from the dialer of a given name. |
| net_conntrack_dialer_conn_established | counter | Total number of connections successfully established by the given dialer a given name. |
| net_conntrack_dialer_conn_failed | counter | Total number of connections failed to dial by the dialer a given name. |
| net_conntrack_listener_conn_accepted | counter | Total number of connections opened to the listener of a given name. |
| net_conntrack_listener_conn_closed | counter | Total number of connections closed that were made to the listener of a given name. |
| process_cpu_seconds | counter | Total user and system CPU time spent in seconds. |
| process_max_fds | gauge | Maximum number of open file descriptors. |
| process_open_fds | gauge | Number of open file descriptors. |
| process_resident_memory_bytes | gauge | Resident memory size in bytes. |
| process_start_time_seconds | gauge | Start time of the process since unix epoch in seconds. |
| process_virtual_memory_bytes | gauge | Virtual memory size in bytes. |
| process_virtual_memory_max_bytes | gauge | Maximum amount of virtual memory available in bytes. |
| prometheus_api_remote_read_queries | gauge | The current number of remote read queries being executed or waiting. |
| prometheus_build_info | gauge | A metric with a constant '1' value labeled by version, revision, branch, and goversion from which prometheus was built. |
| prometheus_config_last_reload_success_timestamp_seconds | gauge | Timestamp of the last successful configuration reload. |
| prometheus_config_last_reload_successful | gauge | Whether the last configuration reload attempt was successful. |
| prometheus_engine_queries | gauge | The current number of queries being executed or waiting. |
| prometheus_engine_queries_concurrent_max | gauge | The max number of concurrent queries. |
| prometheus_engine_query_duration_seconds | summary | Query timings |
| prometheus_engine_query_log_enabled | gauge | State of the query log. |
| prometheus_engine_query_log_failures | counter | The number of query log failures. |
| prometheus_http_request_duration_seconds | histogram | Histogram of latencies for HTTP requests. |
| prometheus_http_requests | counter | Counter of HTTP requests. |
| prometheus_http_response_size_bytes | histogram | Histogram of response size for HTTP requests. |
| prometheus_notifications_alertmanagers_discovered | gauge | The number of alertmanagers discovered and active. |
| prometheus_notifications_dropped | counter | Total number of alerts dropped due to errors when sending to Alertmanager. |
| prometheus_notifications_errors | counter | Total number of errors sending alert notifications. |
| prometheus_notifications_latency_seconds | summary | Latency quantiles for sending alert notifications. |
| prometheus_notifications_queue_capacity | gauge | The capacity of the alert notifications queue. |
| prometheus_notifications_queue_length | gauge | The number of alert notifications in the queue. |
| prometheus_notifications_sent | counter | Total number of alerts sent. |
| prometheus_remote_storage_exemplars_in | counter | Exemplars in to remote storage, compare to exemplars out for queue managers. |
| prometheus_remote_storage_highest_timestamp_in_seconds | gauge | Highest timestamp that has come into the remote storage via the Appender interface, in seconds since epoch. |
| prometheus_remote_storage_samples_in | counter | Samples in to remote storage, compare to samples out for queue managers. |
| prometheus_remote_storage_string_interner_zero_reference_releases | counter | The number of times release has been called for strings that are not interned. |
| prometheus_rule_evaluation_duration_seconds | summary | The duration for a rule to execute. |
| prometheus_rule_evaluation_failures | counter | The total number of rule evaluation failures. |
| prometheus_rule_evaluations | counter | The total number of rule evaluations. |
| prometheus_rule_group_duration_seconds | summary | The duration of rule group evaluations. |
| prometheus_rule_group_interval_seconds | gauge | The interval of a rule group. |
| prometheus_rule_group_iterations_missed | counter | The total number of rule group evaluations missed due to slow rule group evaluation. |
| prometheus_rule_group_iterations | counter | The total number of scheduled rule group evaluations, whether executed or missed. |
| prometheus_rule_group_last_duration_seconds | gauge | The duration of the last rule group evaluation. |
| prometheus_rule_group_last_evaluation_samples | gauge | The number of samples returned during the last rule group evaluation. |
| prometheus_rule_group_last_evaluation_timestamp_seconds | gauge | The timestamp of the last rule group evaluation in seconds. |
| prometheus_rule_group_rules | gauge | The number of rules. |
| prometheus_sd_consul_rpc_duration_seconds | summary | The duration of a Consul RPC call in seconds. |
| prometheus_sd_consul_rpc_failures | counter | The number of Consul RPC call failures. |
| prometheus_sd_discovered_targets | gauge | Current number of discovered targets. |
| prometheus_sd_dns_lookup_failures | counter | The number of DNS-SD lookup failures. |
| prometheus_sd_dns_lookups | counter | The number of DNS-SD lookups. |
| prometheus_sd_failed_configs | gauge | Current number of service discovery configurations that failed to load. |
| prometheus_sd_file_read_errors | counter | The number of File-SD read errors. |
| prometheus_sd_file_scan_duration_seconds | summary | The duration of the File-SD scan in seconds. |
| prometheus_sd_kubernetes_events | counter | The number of Kubernetes events handled. |
| prometheus_sd_kubernetes_http_request_duration_seconds | summary | Summary of latencies for HTTP requests to the Kubernetes API by endpoint. |
| prometheus_sd_kubernetes_http_request | counter | Total number of HTTP requests to the Kubernetes API by status code. |
| prometheus_sd_kubernetes_workqueue_depth | gauge | Current depth of the work queue. |
| prometheus_sd_kubernetes_workqueue_items | counter | Total number of items added to the work queue. |
| prometheus_sd_kubernetes_workqueue_latency_seconds | summary | How long an item stays in the work queue. |
| prometheus_sd_kubernetes_workqueue_longest_running_processor_seconds | gauge | Duration of the longest running processor in the work queue. |
| prometheus_sd_kubernetes_workqueue_unfinished_work_seconds | gauge | How long an item has remained unfinished in the work queue. |
| prometheus_sd_kubernetes_workqueue_work_duration_seconds | summary | How long processing an item from the work queue takes. |
| prometheus_sd_kuma_fetch_duration_seconds | summary | The duration of a Kuma MADS fetch call. |
| prometheus_sd_kuma_fetch_failures | counter | The number of Kuma MADS fetch call failures. |
| prometheus_sd_kuma_fetch_skipped_updates | counter | The number of Kuma MADS fetch calls that result in no updates to the targets. |
| prometheus_sd_received_updates | counter | Total number of update events received from the SD providers. |
| prometheus_sd_updates | counter | Total number of update events sent to the SD consumers. |
| prometheus_target_interval_length_seconds | summary | Actual intervals between scrapes. |
| prometheus_target_metadata_cache_bytes | gauge | The number of bytes that are currently used for storing metric metadata in the cache |
| prometheus_target_metadata_cache_entries | gauge | Total number of metric metadata entries in the cache |
| prometheus_target_scrape_pool_exceeded_label_limits | counter | Total number of times scrape pools hit the label limits, during sync or config reload. |
| prometheus_target_scrape_pool_exceeded_target_limit | counter | Total number of times scrape pools hit the target limit, during sync or config reload. |
| prometheus_target_scrape_pool_reloads_failed | counter | Total number of failed scrape pool reloads. |
| prometheus_target_scrape_pool_reloads | counter | Total number of scrape pool reloads. |
| prometheus_target_scrape_pool_sync | counter | Total number of syncs that were executed on a scrape pool. |
| prometheus_target_scrape_pool_targets | gauge | Current number of targets in this scrape pool. |
| prometheus_target_scrape_pools_failed | counter | Total number of scrape pool creations that failed. |
| prometheus_target_scrape_pools | counter | Total number of scrape pool creation attempts. |
| prometheus_target_scrapes_cache_flush_forced | counter | How many times a scrape cache was flushed due to getting big while scrapes are failing. |
| prometheus_target_scrapes_exceeded_body_size_limit | counter | Total number of scrapes that hit the body size limit |
| prometheus_target_scrapes_exceeded_sample_limit | counter | Total number of scrapes that hit the sample limit and were rejected. |
| prometheus_target_scrapes_exemplar_out_of_order | counter | Total number of exemplar rejected due to not being out of the expected order. |
| prometheus_target_scrapes_sample_duplicate_timestamp | counter | Total number of samples rejected due to duplicate timestamps but different values. |
| prometheus_target_scrapes_sample_out_of_bounds | counter | Total number of samples rejected due to timestamp falling outside of the time bounds. |
| prometheus_target_scrapes_sample_out_of_order | counter | Total number of samples rejected due to not being out of the expected order. |
| prometheus_target_sync_failed | counter | Total number of target sync failures. |
| prometheus_target_sync_length_seconds | summary | Actual interval to sync the scrape pool. |
| prometheus_template_text_expansion_failures | counter | The total number of template text expansion failures. |
| prometheus_template_text_expansions | counter | The total number of template text expansions. |
| prometheus_treecache_watcher_goroutines | gauge | The current number of watcher goroutines. |
| prometheus_treecache_zookeeper_failures | counter | The total number of ZooKeeper failures. |
| prometheus_tsdb_blocks_loaded | gauge | Number of currently loaded data blocks |
| prometheus_tsdb_checkpoint_creations_failed | counter | Total number of checkpoint creations that failed. |
| prometheus_tsdb_checkpoint_creations | counter | Total number of checkpoint creations attempted. |
| prometheus_tsdb_checkpoint_deletions_failed | counter | Total number of checkpoint deletions that failed. |
| prometheus_tsdb_checkpoint_deletions | counter | Total number of checkpoint deletions attempted. |
| prometheus_tsdb_chunk_write_queue_operations | counter | Number of operations on the chunk_write_queue. |
| prometheus_tsdb_clean_start | gauge | -1: lockfile is disabled. 0: a lockfile from a previous execution was replaced. 1: lockfile creation was clean |
| prometheus_tsdb_compaction_chunk_range_seconds | histogram | Final time range of chunks on their first compaction |
| prometheus_tsdb_compaction_chunk_samples | histogram | Final number of samples on their first compaction |
| prometheus_tsdb_compaction_chunk_size_bytes | histogram | Final size of chunks on their first compaction |
| prometheus_tsdb_compaction_duration_seconds | histogram | Duration of compaction runs |
| prometheus_tsdb_compaction_populating_block | gauge | Set to 1 when a block is currently being written to the disk. |
| prometheus_tsdb_compactions_failed | counter | Total number of compactions that failed for the partition. |
| prometheus_tsdb_compactions_skipped | counter | Total number of skipped compactions due to disabled auto compaction. |
| prometheus_tsdb_compactions | counter | Total number of compactions that were executed for the partition. |
| prometheus_tsdb_compactions_triggered | counter | Total number of triggered compactions for the partition. |
| prometheus_tsdb_data_replay_duration_seconds | gauge | Time taken to replay the data on disk. |
| prometheus_tsdb_exemplar_exemplars_appended | counter | Total number of appended exemplars. |
| prometheus_tsdb_exemplar_exemplars_in_storage | gauge | Number of exemplars currently in circular storage. |
| prometheus_tsdb_exemplar_last_exemplars_timestamp_seconds | gauge | The timestamp of the oldest exemplar stored in circular storage. Useful to check for what timerange the current exemplar buffer limit allows. This usually means the last timestampfor all exemplars for a typical setup. This is not true though if one of the series timestamp is in future compared to rest series. |
| prometheus_tsdb_exemplar_max_exemplars | gauge | Total number of exemplars the exemplar storage can store, resizeable. |
| prometheus_tsdb_exemplar_out_of_order_exemplars | counter | Total number of out of order exemplar ingestion failed attempts. |
| prometheus_tsdb_exemplar_series_with_exemplars_in_storage | gauge | Number of series with exemplars currently in circular storage. |
| prometheus_tsdb_head_active_appenders | gauge | Number of currently active appender transactions |
| prometheus_tsdb_head_chunks | gauge | Total number of chunks in the head block. |
| prometheus_tsdb_head_chunks_created | counter | Total number of chunks created in the head |
| prometheus_tsdb_head_chunks_removed | counter | Total number of chunks removed in the head |
| prometheus_tsdb_head_gc_duration_seconds | summary | Runtime of garbage collection in the head block. |
| prometheus_tsdb_head_max_time | gauge | Maximum timestamp of the head block. The unit is decided by the library consumer. |
| prometheus_tsdb_head_max_time_seconds | gauge | Maximum timestamp of the head block. |
| prometheus_tsdb_head_min_time | gauge | Minimum time bound of the head block. The unit is decided by the library consumer. |
| prometheus_tsdb_head_min_time_seconds | gauge | Minimum time bound of the head block. |
| prometheus_tsdb_head_samples_appended | counter | Total number of appended samples. |
| prometheus_tsdb_head_series | gauge | Total number of series in the head block. |
| prometheus_tsdb_head_series_created | counter | Total number of series created in the head |
| prometheus_tsdb_head_series_not_found | counter | Total number of requests for series that were not found. |
| prometheus_tsdb_head_series_removed | counter | Total number of series removed in the head |
| prometheus_tsdb_head_truncations_failed | counter | Total number of head truncations that failed. |
| prometheus_tsdb_head_truncations | counter | Total number of head truncations attempted. |
| prometheus_tsdb_isolation_high_watermark | gauge | The highest TSDB append ID that has been given out. |
| prometheus_tsdb_isolation_low_watermark | gauge | The lowest TSDB append ID that is still referenced. |
| prometheus_tsdb_lowest_timestamp | gauge | Lowest timestamp value stored in the database. The unit is decided by the library consumer. |
| prometheus_tsdb_lowest_timestamp_seconds | gauge | Lowest timestamp value stored in the database. |
| prometheus_tsdb_mmap_chunk_corruptions | counter | Total number of memory-mapped chunk corruptions. |
| prometheus_tsdb_out_of_bound_samples | counter | Total number of out of bound samples ingestion failed attempts. |
| prometheus_tsdb_out_of_order_samples | counter | Total number of out of order samples ingestion failed attempts. |
| prometheus_tsdb_reloads_failures | counter | Number of times the database failed to reloadBlocks block data from disk. |
| prometheus_tsdb_reloads | counter | Number of times the database reloaded block data from disk. |
| prometheus_tsdb_retention_limit_bytes | gauge | Max number of bytes to be retained in the tsdb blocks, configured 0 means disabled |
| prometheus_tsdb_size_retentions | counter | The number of times that blocks were deleted because the maximum number of bytes was exceeded. |
| prometheus_tsdb_snapshot_replay_error | counter | Total number snapshot replays that failed. |
| prometheus_tsdb_storage_blocks_bytes | gauge | The number of bytes that are currently used for local storage by all blocks. |
| prometheus_tsdb_symbol_table_size_bytes | gauge | Size of symbol table in memory for loaded blocks |
| prometheus_tsdb_time_retentions | counter | The number of times that blocks were deleted because the maximum time limit was exceeded. |
| prometheus_tsdb_tombstone_cleanup_seconds | histogram | The time taken to recompact blocks to remove tombstones. |
| prometheus_tsdb_vertical_compactions | counter | Total number of compactions done on overlapping blocks. |
| prometheus_tsdb_wal_completed_pages | counter | Total number of completed pages. |
| prometheus_tsdb_wal_corruptions | counter | Total number of WAL corruptions. |
| prometheus_tsdb_wal_fsync_duration_seconds | summary | Duration of WAL fsync. |
| prometheus_tsdb_wal_page_flushes | counter | Total number of page flushes. |
| prometheus_tsdb_wal_segment_current | gauge | WAL segment index that TSDB is currently writing to. |
| prometheus_tsdb_wal_truncate_duration_seconds | summary | Duration of WAL truncation. |
| prometheus_tsdb_wal_truncations_failed | counter | Total number of WAL truncations that failed. |
| prometheus_tsdb_wal_truncations | counter | Total number of WAL truncations attempted. |
| prometheus_tsdb_wal_writes_failed | counter | Total number of WAL writes that failed. |
| prometheus_web_federation_errors | counter | Total number of errors that occurred while sending federation responses. |
| prometheus_web_federation_warnings | counter | Total number of warnings that occurred while sending federation responses. |
| promhttp_metric_handler_requests_in_flight | gauge | Current number of scrapes being served. |
| promhttp_metric_handler_requests | counter | Total number of scrapes by HTTP status code. |
<!-- markdownlint-enable line-length -->

## Prometheus-operator (Kubernetes v1.23)

### Statistics

<!-- markdownlint-disable line-length -->
| Number of metrics families | Number of time series | Average number of labels per time series | Number of unique label pairs |
| -------------------------- | --------------------- | ---------------------------------------- | ---------------------------- |
| 87 | 1822 | 1.0 | 850 |
<!-- markdownlint-enable line-length -->

### Metrics

<!-- markdownlint-disable line-length -->
| Name | Type, Unit | Description |
| ---- | ---------- | ----------- |
| go_gc_cycles_automatic_gc_cycles | counter | Count of completed GC cycles generated by the Go runtime. |
| go_gc_cycles_forced_gc_cycles | counter | Count of completed GC cycles forced by the application. |
| go_gc_cycles_total_gc_cycles | counter | Count of all completed GC cycles. |
| go_gc_duration_seconds | summary | A summary of the pause duration of garbage collection cycles. |
| go_gc_heap_allocs_by_size_bytes_total | histogram | Distribution of heap allocations by approximate size. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks. |
| go_gc_heap_allocs_bytes | counter | Cumulative sum of memory allocated to the heap by the application. |
| go_gc_heap_allocs_objects | counter | Cumulative count of heap allocations triggered by the application. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks. |
| go_gc_heap_frees_by_size_bytes_total | histogram | Distribution of freed heap allocations by approximate size. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks. |
| go_gc_heap_frees_bytes | counter | Cumulative sum of heap memory freed by the garbage collector. |
| go_gc_heap_frees_objects | counter | Cumulative count of heap allocations whose storage was freed by the garbage collector. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks. |
| go_gc_heap_goal_bytes | gauge | Heap size target for the end of the GC cycle. |
| go_gc_heap_objects_objects | gauge | Number of objects, live or unswept, occupying heap memory. |
| go_gc_heap_tiny_allocs_objects | counter | Count of small allocations that are packed together into blocks. These allocations are counted separately from other allocations because each individual allocation is not tracked by the runtime, only their block. Each block is already accounted for in allocs-by-size and frees-by-size. |
| go_gc_pauses_seconds_total | histogram | Distribution individual GC-related stop-the-world pause latencies. |
| go_goroutines | gauge | Number of goroutines that currently exist. |
| go_info | gauge | Information about the Go environment. |
| go_memory_classes_heap_free_bytes | gauge | Memory that is completely free and eligible to be returned to the underlying system, but has not been. This metric is the runtime's estimate of free address space that is backed by physical memory. |
| go_memory_classes_heap_objects_bytes | gauge | Memory occupied by live objects and dead objects that have not yet been marked free by the garbage collector. |
| go_memory_classes_heap_released_bytes | gauge | Memory that is completely free and has been returned to the underlying system. This metric is the runtime's estimate of free address space that is still mapped into the process, but is not backed by physical memory. |
| go_memory_classes_heap_stacks_bytes | gauge | Memory allocated from the heap that is reserved for stack space, whether or not it is currently in-use. |
| go_memory_classes_heap_unused_bytes | gauge | Memory that is reserved for heap objects but is not currently used to hold heap objects. |
| go_memory_classes_metadata_mcache_free_bytes | gauge | Memory that is reserved for runtime mcache structures, but not in-use. |
| go_memory_classes_metadata_mcache_inuse_bytes | gauge | Memory that is occupied by runtime mcache structures that are currently being used. |
| go_memory_classes_metadata_mspan_free_bytes | gauge | Memory that is reserved for runtime mspan structures, but not in-use. |
| go_memory_classes_metadata_mspan_inuse_bytes | gauge | Memory that is occupied by runtime mspan structures that are currently being used. |
| go_memory_classes_metadata_other_bytes | gauge | Memory that is reserved for or used to hold runtime metadata. |
| go_memory_classes_os_stacks_bytes | gauge | Stack memory allocated by the underlying operating system. |
| go_memory_classes_other_bytes | gauge | Memory used by execution trace buffers, structures for debugging the runtime, finalizer and profiler specials, and more. |
| go_memory_classes_profiling_buckets_bytes | gauge | Memory that is used by the stack trace hash map used for profiling. |
| go_memory_classes_total_bytes | gauge | All memory mapped by the Go runtime into the current process as read-write. Note that this does not include memory mapped by code called via cgo or via the syscall package. Sum of all metrics in /memory/classes. |
| go_memstats_alloc_bytes | gauge | Number of bytes allocated and still in use. |
| go_memstats_alloc_bytes | counter | Total number of bytes allocated, even if freed. |
| go_memstats_buck_hash_sys_bytes | gauge | Number of bytes used by the profiling bucket hash table. |
| go_memstats_frees | counter | Total number of frees. |
| go_memstats_gc_cpu_fraction | gauge | The fraction of this program's available CPU time used by the GC since the program started. |
| go_memstats_gc_sys_bytes | gauge | Number of bytes used for garbage collection system metadata. |
| go_memstats_heap_alloc_bytes | gauge | Number of heap bytes allocated and still in use. |
| go_memstats_heap_idle_bytes | gauge | Number of heap bytes waiting to be used. |
| go_memstats_heap_inuse_bytes | gauge | Number of heap bytes that are in use. |
| go_memstats_heap_objects | gauge | Number of allocated objects. |
| go_memstats_heap_released_bytes | gauge | Number of heap bytes released to OS. |
| go_memstats_heap_sys_bytes | gauge | Number of heap bytes obtained from system. |
| go_memstats_last_gc_time_seconds | gauge | Number of seconds since 1970 of last garbage collection. |
| go_memstats_lookups | counter | Total number of pointer lookups. |
| go_memstats_mallocs | counter | Total number of mallocs. |
| go_memstats_mcache_inuse_bytes | gauge | Number of bytes in use by mcache structures. |
| go_memstats_mcache_sys_bytes | gauge | Number of bytes used for mcache structures obtained from system. |
| go_memstats_mspan_inuse_bytes | gauge | Number of bytes in use by mspan structures. |
| go_memstats_mspan_sys_bytes | gauge | Number of bytes used for mspan structures obtained from system. |
| go_memstats_next_gc_bytes | gauge | Number of heap bytes when next garbage collection will take place. |
| go_memstats_other_sys_bytes | gauge | Number of bytes used for other system allocations. |
| go_memstats_stack_inuse_bytes | gauge | Number of bytes in use by the stack allocator. |
| go_memstats_stack_sys_bytes | gauge | Number of bytes obtained from system for stack allocator. |
| go_memstats_sys_bytes | gauge | Number of bytes obtained from system. |
| go_sched_goroutines_goroutines | gauge | Count of live goroutines. |
| go_sched_latencies_seconds | histogram | Distribution of the time goroutines have spent in the scheduler in a runnable state before actually running. |
| go_threads | gauge | Number of OS threads created. |
| process_cpu_seconds | counter | Total user and system CPU time spent in seconds. |
| process_max_fds | gauge | Maximum number of open file descriptors. |
| process_open_fds | gauge | Number of open file descriptors. |
| process_resident_memory_bytes | gauge | Resident memory size in bytes. |
| process_start_time_seconds | gauge | Start time of the process since unix epoch in seconds. |
| process_virtual_memory_bytes | gauge | Virtual memory size in bytes. |
| process_virtual_memory_max_bytes | gauge | Maximum amount of virtual memory available in bytes. |
| prometheus_operator_alertmanager_config_validation_errors | counter | Number of errors that occurred while validating a alertmanagerconfig object |
| prometheus_operator_alertmanager_config_validation_triggered | counter | Number of times an alertmanagerconfig object triggered validation |
| prometheus_operator_build_info | gauge | A metric with a constant '1' value labeled by version, revision, branch, and goversion from which prometheus_operator was built. |
| prometheus_operator_kubernetes_client_http_request_duration_seconds | summary | Summary of latencies for the Kubernetes client's requests by endpoint. |
| prometheus_operator_kubernetes_client_http_requests | counter | Total number of Kubernetes's client requests by status code. |
| prometheus_operator_kubernetes_client_rate_limiter_duration_seconds | summary | Summary of latencies for the Kuberntes client's rate limiter by endpoint. |
| prometheus_operator_list_operations_failed | counter | Total number of list operations that failed |
| prometheus_operator_list_operations | counter | Total number of list operations |
| prometheus_operator_managed_resources | gauge | Number of resources managed by the operator's controller per state (selected/rejected) |
| prometheus_operator_node_address_lookup_errors | counter | Number of times a node IP address could not be determined |
| prometheus_operator_node_syncs_failed | counter | Number of node endpoints synchronisation failures |
| prometheus_operator_node_syncs | counter | Number of node endpoints synchronisations |
| prometheus_operator_ready | gauge | 1 when the controller is ready to reconcile resources, 0 otherwise |
| prometheus_operator_reconcile_errors | counter | Number of errors that occurred during reconcile operations |
| prometheus_operator_reconcile_operations | counter | Total number of reconcile operations |
| prometheus_operator_reconcile_sts_delete_create | counter | Number of times that reconciling a statefulset required deleting and re-creating it |
| prometheus_operator_rule_validation_errors | counter | Number of errors that occurred while validating a prometheusRules object |
| prometheus_operator_rule_validation_triggered | counter | Number of times a prometheusRule object triggered validation |
| prometheus_operator_spec_replicas | gauge | Number of expected replicas for the object. |
| prometheus_operator_syncs | gauge | Number of objects per sync status (ok/failed) |
| prometheus_operator_triggered | counter | Number of times a Kubernetes object add, delete or update event triggered the Prometheus Operator to reconcile an object |
| prometheus_operator_watch_operations_failed | counter | Total number of watch operations that failed |
| prometheus_operator_watch_operations | counter | Total number of watch operations |
<!-- markdownlint-enable line-length -->

## Operator-sdk

### Metrics
<!-- markdownlint-disable line-length -->
| Name | Type, Unit | Description |
| ---- | ---------- | ----------- |
| controller_runtime_active_workers | gauge | Number of currently used workers per controller |
| controller_runtime_max_concurrent_reconciles | gauge | Maximum number of concurrent reconciles per controller |
| controller_runtime_reconcile_errors_total | counter | Total number of reconciliation errors per controller |
| controller_runtime_reconcile_time_seconds | histogram | Length of time per reconciliation per controller |
| controller_runtime_reconcile_total | counter | Total number of reconciliations per controller |
| controller_runtime_reconcile_total | counter | Total number of reconciliations per controller |
<!-- markdownlint-enable line-length -->

## Openshift-state-metrics

### Statistics

<!-- markdownlint-disable line-length -->
| Number of metrics families | Number of time series | Average number of labels per time series | Number of unique label pairs |
| -------------------------- | --------------------- | ---------------------------------------- | ---------------------------- |
| 68 | 124 | 0.42 | 28 |
<!-- markdownlint-enable line-length -->

### Metrics

<!-- markdownlint-disable line-length -->
| Name | Type, Unit | Description |
| ---- | ---------- | ----------- |
| go_cgo_go_to_c_calls_calls | counter | Count of calls made from Go to C by the current process. |
| go_gc_cycles_automatic_gc_cycles | counter | Count of completed GC cycles generated by the Go runtime. |
| go_gc_cycles_forced_gc_cycles | counter | Count of completed GC cycles forced by the application. |
| go_gc_cycles_total_gc_cycles | counter | Count of all completed GC cycles. |
| go_gc_duration_seconds | summary | A summary of the pause duration of garbage collection cycles. |
| go_gc_heap_allocs_by_size_bytes_total | histogram | Distribution of heap allocations by approximate size. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks. |
| go_gc_heap_allocs_bytes | counter | Cumulative sum of memory allocated to the heap by the application. |
| go_gc_heap_allocs_objects | counter | Cumulative count of heap allocations triggered by the application. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks. |
| go_gc_heap_frees_by_size_bytes_total | histogram | Distribution of freed heap allocations by approximate size. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks. |
| go_gc_heap_frees_bytes | counter | Cumulative sum of heap memory freed by the garbage collector. |
| go_gc_heap_frees_objects | counter | Cumulative count of heap allocations whose storage was freed by the garbage collector. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks. |
| go_gc_heap_goal_bytes | gauge | Heap size target for the end of the GC cycle. |
| go_gc_heap_objects_objects | gauge | Number of objects, live or unswept, occupying heap memory. |
| go_gc_heap_tiny_allocs_objects | counter | Count of small allocations that are packed together into blocks. These allocations are counted separately from other allocations because each individual allocation is not tracked by the runtime, only their block. Each block is already accounted for in allocs-by-size and frees-by-size. |
| go_gc_limiter_last_enabled_gc_cycle | gauge | GC cycle the last time the GC CPU limiter was enabled. This metric is useful for diagnosing the root cause of an out-of-memory error, because the limiter trades memory for CPU time when the GC's CPU time gets too high. This is most likely to occur with use of SetMemoryLimit. The first GC cycle is cycle 1, so a value of 0 indicates that it was never enabled. |
| go_gc_pauses_seconds_total | histogram | Distribution individual GC-related stop-the-world pause latencies. |
| go_gc_stack_starting_size_bytes | gauge | The stack size of new goroutines. |
| go_goroutines | gauge | Number of goroutines that currently exist. |
| go_info | gauge | Information about the Go environment. |
| go_memory_classes_heap_free_bytes | gauge | Memory that is completely free and eligible to be returned to the underlying system, but has not been. This metric is the runtime's estimate of free address space that is backed by physical memory. |
| go_memory_classes_heap_objects_bytes | gauge | Memory occupied by live objects and dead objects that have not yet been marked free by the garbage collector. |
| go_memory_classes_heap_released_bytes | gauge | Memory that is completely free and has been returned to the underlying system. This metric is the runtime's estimate of free address space that is still mapped into the process, but is not backed by physical memory. |
| go_memory_classes_heap_stacks_bytes | gauge | Memory allocated from the heap that is reserved for stack space, whether or not it is currently in-use. |
| go_memory_classes_heap_unused_bytes | gauge | Memory that is reserved for heap objects but is not currently used to hold heap objects. |
| go_memory_classes_metadata_mcache_free_bytes | gauge | Memory that is reserved for runtime mcache structures, but not in-use. |
| go_memory_classes_metadata_mcache_inuse_bytes | gauge | Memory that is occupied by runtime mcache structures that are currently being used. |
| go_memory_classes_metadata_mspan_free_bytes | gauge | Memory that is reserved for runtime mspan structures, but not in-use. |
| go_memory_classes_metadata_mspan_inuse_bytes | gauge | Memory that is occupied by runtime mspan structures that are currently being used. |
| go_memory_classes_metadata_other_bytes | gauge | Memory that is reserved for or used to hold runtime metadata. |
| go_memory_classes_os_stacks_bytes | gauge | Stack memory allocated by the underlying operating system. |
| go_memory_classes_other_bytes | gauge | Memory used by execution trace buffers, structures for debugging the runtime, finalizer and profiler specials, and more. |
| go_memory_classes_profiling_buckets_bytes | gauge | Memory that is used by the stack trace hash map used for profiling. |
| go_memory_classes_total_bytes | gauge | All memory mapped by the Go runtime into the current process as read-write. Note that this does not include memory mapped by code called via cgo or via the syscall package. Sum of all metrics in /memory/classes. |
| go_memstats_alloc_bytes | gauge | Number of bytes allocated and still in use. |
| go_memstats_alloc_bytes | counter | Total number of bytes allocated, even if freed. |
| go_memstats_buck_hash_sys_bytes | gauge | Number of bytes used by the profiling bucket hash table. |
| go_memstats_frees | counter | Total number of frees. |
| go_memstats_gc_cpu_fraction | gauge | The fraction of this program's available CPU time used by the GC since the program started. |
| go_memstats_gc_sys_bytes | gauge | Number of bytes used for garbage collection system metadata. |
| go_memstats_heap_alloc_bytes | gauge | Number of heap bytes allocated and still in use. |
| go_memstats_heap_idle_bytes | gauge | Number of heap bytes waiting to be used. |
| go_memstats_heap_inuse_bytes | gauge | Number of heap bytes that are in use. |
| go_memstats_heap_objects | gauge | Number of allocated objects. |
| go_memstats_heap_released_bytes | gauge | Number of heap bytes released to OS. |
| go_memstats_heap_sys_bytes | gauge | Number of heap bytes obtained from system. |
| go_memstats_last_gc_time_seconds | gauge | Number of seconds since 1970 of last garbage collection. |
| go_memstats_lookups | counter | Total number of pointer lookups. |
| go_memstats_mallocs | counter | Total number of mallocs. |
| go_memstats_mcache_inuse_bytes | gauge | Number of bytes in use by mcache structures. |
| go_memstats_mcache_sys_bytes | gauge | Number of bytes used for mcache structures obtained from system. |
| go_memstats_mspan_inuse_bytes | gauge | Number of bytes in use by mspan structures. |
| go_memstats_mspan_sys_bytes | gauge | Number of bytes used for mspan structures obtained from system. |
| go_memstats_next_gc_bytes | gauge | Number of heap bytes when next garbage collection will take place. |
| go_memstats_other_sys_bytes | gauge | Number of bytes used for other system allocations. |
| go_memstats_stack_inuse_bytes | gauge | Number of bytes in use by the stack allocator. |
| go_memstats_stack_sys_bytes | gauge | Number of bytes obtained from system for stack allocator. |
| go_memstats_sys_bytes | gauge | Number of bytes obtained from system. |
| go_sched_gomaxprocs_threads | gauge | The current runtime.GOMAXPROCS setting, or the number of operating system threads that can execute user-level Go code simultaneously. |
| go_sched_goroutines_goroutines | gauge | Count of live goroutines. |
| go_sched_latencies_seconds | histogram | Distribution of the time goroutines have spent in the scheduler in a runnable state before actually running. |
| go_threads | gauge | Number of OS threads created. |
| process_cpu_seconds | counter | Total user and system CPU time spent in seconds. |
| process_max_fds | gauge | Maximum number of open file descriptors. |
| process_open_fds | gauge | Number of open file descriptors. |
| process_resident_memory_bytes | gauge | Resident memory size in bytes. |
| process_start_time_seconds | gauge | Start time of the process since unix epoch in seconds. |
| process_virtual_memory_bytes | gauge | Virtual memory size in bytes. |
| process_virtual_memory_max_bytes | gauge | Maximum amount of virtual memory available in bytes. |
<!-- markdownlint-enable line-length -->

## Openshift-api-server

### Statistics

<!-- markdownlint-disable line-length -->
| Number of metrics families | Number of time series | Average number of labels per time series | Number of unique label pairs |
| -------------------------- | --------------------- | ---------------------------------------- | ---------------------------- |
| 152 | 6528 | 6.41 | 326 |
<!-- markdownlint-enable line-length -->

### Metrics

<!-- markdownlint-disable line-length -->
| Name | Type, Unit | Description |
| ---- | ---------- | ----------- |
| aggregator_openapi_v2_regeneration_count | counter | ALPHA Counter of OpenAPI v2 spec regeneration count broken down by causing APIService name and reason. |
| aggregator_openapi_v2_regeneration_duration | gauge | ALPHA Gauge of OpenAPI v2 spec regeneration duration in seconds. |
| apiserver_admission_controller_admission_duration_seconds | histogram | STABLE Admission controller latency histogram in seconds, identified by name and broken out for each operation and API resource and type (validate or admit). |
| apiserver_admission_step_admission_duration_seconds | histogram | STABLE Admission sub-step latency histogram in seconds, broken out for each operation and API resource and step type (validate or admit). |
| apiserver_admission_step_admission_duration_seconds_summary | summary | ALPHA Admission sub-step latency summary in seconds, broken out for each operation and API resource and step type (validate or admit). |
| apiserver_audit_event | counter | ALPHA Counter of audit events generated and sent to the audit backend. |
| apiserver_audit_level | counter | ALPHA Counter of policy levels for audit events (1 per request). |
| apiserver_audit_requests_rejected | counter | ALPHA Counter of apiserver requests rejected due to an error in audit logging backend. |
| apiserver_cache_list_fetched_objects | counter | ALPHA Number of objects read from watch cache in the course of serving a LIST request |
| apiserver_cache_list_returned_objects | counter | ALPHA Number of objects returned for a LIST request from watch cache |
| apiserver_cache_list | counter | ALPHA Number of LIST requests served from watch cache |
| apiserver_cel_compilation_duration_seconds | histogram | ALPHA |
| apiserver_cel_evaluation_duration_seconds | histogram | ALPHA |
| apiserver_client_certificate_expiration_seconds | histogram | ALPHA Distribution of the remaining lifetime on the certificate used to authenticate a request. |
| apiserver_current_inflight_requests | gauge | STABLE Maximal number of currently used inflight request limit of this apiserver per request kind in last second. |
| apiserver_delegated_authn_request_duration_seconds | histogram | ALPHA Request latency in seconds. Broken down by status code. |
| apiserver_delegated_authn_request | counter | ALPHA Number of HTTP requests partitioned by status code. |
| apiserver_delegated_authz_request_duration_seconds | histogram | ALPHA Request latency in seconds. Broken down by status code. |
| apiserver_delegated_authz_request | counter | ALPHA Number of HTTP requests partitioned by status code. |
| apiserver_envelope_encryption_dek_cache_fill_percent | gauge | ALPHA Percent of the cache slots currently occupied by cached DEKs. |
| apiserver_flowcontrol_read_vs_write_current_requests | histogram | EXPERIMENTAL: ALPHA Observations, at the end of every nanosecond, of the number of requests (as a fraction of the relevant limit) waiting or in regular stage of execution |
| apiserver_flowcontrol_seat_fair_frac | gauge | ALPHA Fair fraction of server's concurrency to allocate to each priority level that can use it |
| apiserver_init_events | counter | ALPHA Counter of init events processed in watch cache broken by resource type. |
| apiserver_kube_aggregator_x509_insecure_sha1 | counter | ALPHA Counts the number of requests to servers with insecure SHA1 signatures in their serving certificate OR the number of connection failures due to the insecure SHA1 signatures (either/or, based on the runtime environment) |
| apiserver_kube_aggregator_x509_missing_san | counter | ALPHA Counts the number of requests to servers missing SAN extension in their serving certificate OR the number of connection failures due to the lack of x509 certificate SAN extension missing (either/or, based on the runtime environment) |
| apiserver_longrunning_requests | gauge | STABLE Gauge of all active long-running apiserver requests broken out by verb, group, version, resource, scope and component. Not all requests are tracked this way. |
| apiserver_request_aborts | counter | ALPHA Number of requests which apiserver aborted possibly due to a timeout, for each group, version, verb, resource, subresource and scope |
| apiserver_request_duration_seconds | histogram | STABLE Response latency distribution in seconds for each verb, dry run value, group, version, resource, subresource, scope and component. |
| apiserver_request_filter_duration_seconds | histogram | ALPHA Request filter latency distribution in seconds, for each filter type |
| apiserver_request_post_timeout | counter | ALPHA Tracks the activity of the request handlers after the associated requests have been timed out by the apiserver |
| apiserver_request_sli_duration_seconds | histogram | ALPHA Response latency distribution (not counting webhook duration) in seconds for each verb, group, version, resource, subresource, scope and component. |
| apiserver_request_slo_duration_seconds | histogram | ALPHA Response latency distribution (not counting webhook duration) in seconds for each verb, group, version, resource, subresource, scope and component. |
| apiserver_request_terminations | counter | ALPHA Number of requests which apiserver terminated in self-defense. |
| apiserver_request_timestamp_comparison_time | histogram | ALPHA Time taken for comparison of old vs new objects in UPDATE or PATCH requests |
| apiserver_request | counter | STABLE Counter of apiserver requests broken out for each verb, dry run value, group, version, resource, scope, component, and HTTP response code. |
| apiserver_response_sizes | histogram | STABLE Response size distribution in bytes for each group, version, verb, resource, subresource, scope and component. |
| apiserver_selfrequest | counter | ALPHA Counter of apiserver self-requests broken out for each verb, API resource and subresource. |
| apiserver_storage_data_key_generation_duration_seconds | histogram | ALPHA Latencies in seconds of data encryption key(DEK) generation operations. |
| apiserver_storage_data_key_generation_failures | counter | ALPHA Total number of failed data encryption key(DEK) generation operations. |
| apiserver_storage_db_total_size_in_bytes | gauge | ALPHA Total size of the storage database file physically allocated in bytes. |
| apiserver_storage_envelope_transformation_cache_misses | counter | ALPHA Total number of cache misses while accessing key decryption key(KEK). |
| apiserver_storage_list_evaluated_objects | counter | ALPHA Number of objects tested in the course of serving a LIST request from storage |
| apiserver_storage_list_fetched_objects | counter | ALPHA Number of objects read from storage in the course of serving a LIST request |
| apiserver_storage_list_returned_objects | counter | ALPHA Number of objects returned for a LIST request from storage |
| apiserver_storage_list | counter | ALPHA Number of LIST requests served from storage |
| apiserver_storage_objects | gauge | STABLE Number of stored objects at the time of last check split by kind. |
| apiserver_tls_handshake_errors | counter | ALPHA Number of requests dropped with 'TLS handshake error from' error |
| apiserver_watch_cache_events_dispatched | counter | ALPHA Counter of events dispatched in watch cache broken by resource type. |
| apiserver_watch_cache_initializations | counter | ALPHA Counter of watch cache initializations broken by resource type. |
| apiserver_watch_events_sizes | histogram | ALPHA Watch event size distribution in bytes |
| apiserver_watch_events | counter | ALPHA Number of events sent in watch clients |
| apiserver_webhooks_x509_insecure_sha1 | counter | ALPHA Counts the number of requests to servers with insecure SHA1 signatures in their serving certificate OR the number of connection failures due to the insecure SHA1 signatures (either/or, based on the runtime environment) |
| apiserver_webhooks_x509_missing_san | counter | ALPHA Counts the number of requests to servers missing SAN extension in their serving certificate OR the number of connection failures due to the lack of x509 certificate SAN extension missing (either/or, based on the runtime environment) |
| authenticated_user_requests | counter | ALPHA Counter of authenticated requests broken out by username. |
| authentication_attempts | counter | ALPHA Counter of authenticated attempts. |
| authentication_duration_seconds | histogram | ALPHA Authentication duration in seconds broken out by result. |
| authentication_token_cache_active_fetch_count | gauge | ALPHA |
| authentication_token_cache_fetch | counter | ALPHA |
| authentication_token_cache_request_duration_seconds | histogram | ALPHA |
| authentication_token_cache_request | counter | ALPHA |
| disabled_metric | counter | ALPHA The count of disabled metrics. |
| etcd_bookmark_counts | gauge | ALPHA Number of etcd bookmarks (progress notify events) split by kind. |
| etcd_request_duration_seconds | histogram | ALPHA Etcd request latency in seconds for each operation and object type. |
| field_validation_request_duration_seconds | histogram | ALPHA Response latency distribution in seconds for each field validation value and whether field validation is enabled or not |
| go_cgo_go_to_c_calls_calls | counter | Count of calls made from Go to C by the current process. |
| go_gc_cycles_automatic_gc_cycles | counter | Count of completed GC cycles generated by the Go runtime. |
| go_gc_cycles_forced_gc_cycles | counter | Count of completed GC cycles forced by the application. |
| go_gc_cycles_total_gc_cycles | counter | Count of all completed GC cycles. |
| go_gc_duration_seconds | summary | A summary of the pause duration of garbage collection cycles. |
| go_gc_heap_allocs_by_size_bytes | histogram | Distribution of heap allocations by approximate size. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks. |
| go_gc_heap_allocs_bytes | counter | Cumulative sum of memory allocated to the heap by the application. |
| go_gc_heap_allocs_objects | counter | Cumulative count of heap allocations triggered by the application. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks. |
| go_gc_heap_frees_by_size_bytes | histogram | Distribution of freed heap allocations by approximate size. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks. |
| go_gc_heap_frees_bytes | counter | Cumulative sum of heap memory freed by the garbage collector. |
| go_gc_heap_frees_objects | counter | Cumulative count of heap allocations whose storage was freed by the garbage collector. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks. |
| go_gc_heap_goal_bytes | gauge | Heap size target for the end of the GC cycle. |
| go_gc_heap_objects_objects | gauge | Number of objects, live or unswept, occupying heap memory. |
| go_gc_heap_tiny_allocs_objects | counter | Count of small allocations that are packed together into blocks. These allocations are counted separately from other allocations because each individual allocation is not tracked by the runtime, only their block. Each block is already accounted for in allocs-by-size and frees-by-size. |
| go_gc_limiter_last_enabled_gc_cycle | gauge | GC cycle the last time the GC CPU limiter was enabled. This metric is useful for diagnosing the root cause of an out-of-memory error, because the limiter trades memory for CPU time when the GC's CPU time gets too high. This is most likely to occur with use of SetMemoryLimit. The first GC cycle is cycle 1, so a value of 0 indicates that it was never enabled. |
| go_gc_pauses_seconds | histogram | Distribution individual GC-related stop-the-world pause latencies. |
| go_gc_stack_starting_size_bytes | gauge | The stack size of new goroutines. |
| go_goroutines | gauge | Number of goroutines that currently exist. |
| go_info | gauge | Information about the Go environment. |
| go_memory_classes_heap_free_bytes | gauge | Memory that is completely free and eligible to be returned to the underlying system, but has not been. This metric is the runtime's estimate of free address space that is backed by physical memory. |
| go_memory_classes_heap_objects_bytes | gauge | Memory occupied by live objects and dead objects that have not yet been marked free by the garbage collector. |
| go_memory_classes_heap_released_bytes | gauge | Memory that is completely free and has been returned to the underlying system. This metric is the runtime's estimate of free address space that is still mapped into the process, but is not backed by physical memory. |
| go_memory_classes_heap_stacks_bytes | gauge | Memory allocated from the heap that is reserved for stack space, whether or not it is currently in-use. |
| go_memory_classes_heap_unused_bytes | gauge | Memory that is reserved for heap objects but is not currently used to hold heap objects. |
| go_memory_classes_metadata_mcache_free_bytes | gauge | Memory that is reserved for runtime mcache structures, but not in-use. |
| go_memory_classes_metadata_mcache_inuse_bytes | gauge | Memory that is occupied by runtime mcache structures that are currently being used. |
| go_memory_classes_metadata_mspan_free_bytes | gauge | Memory that is reserved for runtime mspan structures, but not in-use. |
| go_memory_classes_metadata_mspan_inuse_bytes | gauge | Memory that is occupied by runtime mspan structures that are currently being used. |
| go_memory_classes_metadata_other_bytes | gauge | Memory that is reserved for or used to hold runtime metadata. |
| go_memory_classes_os_stacks_bytes | gauge | Stack memory allocated by the underlying operating system. |
| go_memory_classes_other_bytes | gauge | Memory used by execution trace buffers, structures for debugging the runtime, finalizer and profiler specials, and more. |
| go_memory_classes_profiling_buckets_bytes | gauge | Memory that is used by the stack trace hash map used for profiling. |
| go_memory_classes_total_bytes | gauge | All memory mapped by the Go runtime into the current process as read-write. Note that this does not include memory mapped by code called via cgo or via the syscall package. Sum of all metrics in /memory/classes. |
| go_memstats_alloc_bytes | gauge | Number of bytes allocated and still in use. |
| go_memstats_alloc_bytes | counter | Total number of bytes allocated, even if freed. |
| go_memstats_buck_hash_sys_bytes | gauge | Number of bytes used by the profiling bucket hash table. |
| go_memstats_frees | counter | Total number of frees. |
| go_memstats_gc_sys_bytes | gauge | Number of bytes used for garbage collection system metadata. |
| go_memstats_heap_alloc_bytes | gauge | Number of heap bytes allocated and still in use. |
| go_memstats_heap_idle_bytes | gauge | Number of heap bytes waiting to be used. |
| go_memstats_heap_inuse_bytes | gauge | Number of heap bytes that are in use. |
| go_memstats_heap_objects | gauge | Number of allocated objects. |
| go_memstats_heap_released_bytes | gauge | Number of heap bytes released to OS. |
| go_memstats_heap_sys_bytes | gauge | Number of heap bytes obtained from system. |
| go_memstats_last_gc_time_seconds | gauge | Number of seconds since 1970 of last garbage collection. |
| go_memstats_lookups | counter | Total number of pointer lookups. |
| go_memstats_mallocs | counter | Total number of mallocs. |
| go_memstats_mcache_inuse_bytes | gauge | Number of bytes in use by mcache structures. |
| go_memstats_mcache_sys_bytes | gauge | Number of bytes used for mcache structures obtained from system. |
| go_memstats_mspan_inuse_bytes | gauge | Number of bytes in use by mspan structures. |
| go_memstats_mspan_sys_bytes | gauge | Number of bytes used for mspan structures obtained from system. |
| go_memstats_next_gc_bytes | gauge | Number of heap bytes when next garbage collection will take place. |
| go_memstats_other_sys_bytes | gauge | Number of bytes used for other system allocations. |
| go_memstats_stack_inuse_bytes | gauge | Number of bytes in use by the stack allocator. |
| go_memstats_stack_sys_bytes | gauge | Number of bytes obtained from system for stack allocator. |
| go_memstats_sys_bytes | gauge | Number of bytes obtained from system. |
| go_sched_gomaxprocs_threads | gauge | The current runtime.GOMAXPROCS setting, or the number of operating system threads that can execute user-level Go code simultaneously. |
| go_sched_goroutines_goroutines | gauge | Count of live goroutines. |
| go_sched_latencies_seconds | histogram | Distribution of the time goroutines have spent in the scheduler in a runnable state before actually running. |
| go_threads | gauge | Number of OS threads created. |
| grpc_client_handled | counter | Total number of RPCs completed by the client, regardless of success or failure. |
| grpc_client_msg_received | counter | Total number of RPC stream messages received by the client. |
| grpc_client_msg_sent | counter | Total number of gRPC stream messages sent by the client. |
| grpc_client_started | counter | Total number of RPCs started on the client. |
| hidden_metric | counter | ALPHA The count of hidden metrics. |
| process_cpu_seconds | counter | Total user and system CPU time spent in seconds. |
| process_max_fds | gauge | Maximum number of open file descriptors. |
| process_open_fds | gauge | Number of open file descriptors. |
| process_resident_memory_bytes | gauge | Resident memory size in bytes. |
| process_start_time_seconds | gauge | Start time of the process since unix epoch in seconds. |
| process_virtual_memory_bytes | gauge | Virtual memory size in bytes. |
| process_virtual_memory_max_bytes | gauge | Maximum amount of virtual memory available in bytes. |
| registered_metric | counter | ALPHA The count of registered metrics broken by stability level and deprecation version. |
| rest_client_exec_plugin_certificate_rotation_age | histogram | ALPHA Histogram of the number of seconds the last auth exec plugin client certificate lived before being rotated. If auth exec plugin client certificates are unused, histogram will contain no data. |
| rest_client_exec_plugin_ttl_seconds | gauge | ALPHA Gauge of the shortest TTL (time-to-live) of the client certificate(s) managed by the auth exec plugin. The value is in seconds until certificate expiry (negative if already expired). If auth exec plugins are unused or manage no TLS certificates, the value will be +INF. |
| rest_client_rate_limiter_duration_seconds | histogram | ALPHA Client side rate limiter latency in seconds. Broken down by verb, and host. |
| rest_client_request_duration_seconds | histogram | ALPHA Request latency in seconds. Broken down by verb, and host. |
| rest_client_request_size_bytes | histogram | ALPHA Request size in bytes. Broken down by verb and host. |
| rest_client_requests | counter | ALPHA Number of HTTP requests, partitioned by status code, method, and host. |
| rest_client_response_size_bytes | histogram | ALPHA Response size in bytes. Broken down by verb and host. |
| watch_cache_capacity | gauge | ALPHA Total capacity of watch cache broken by resource type. |
| workqueue_adds | counter | ALPHA Total number of adds handled by workqueue |
| workqueue_depth | gauge | ALPHA Current depth of workqueue |
| workqueue_longest_running_processor_seconds | gauge | ALPHA How many seconds has the longest running processor for workqueue been running. |
| workqueue_queue_duration_seconds | histogram | ALPHA How long in seconds an item stays in workqueue before being requested. |
| workqueue_retries | counter | ALPHA Total number of retries handled by workqueue |
| workqueue_unfinished_work_seconds | gauge | ALPHA How many seconds of work has done that is in progress and hasn't been observed by work_duration. Large values indicate stuck threads. One can deduce the number of stuck threads by observing the rate at which this increases. |
| workqueue_work_duration_seconds | histogram | ALPHA How long in seconds processing an item from workqueue takes. |
<!-- markdownlint-enable line-length -->

## Openshift-api-server-operator

### Statistics

<!-- markdownlint-disable line-length -->
| Number of metrics families | Number of time series | Average number of labels per time series | Number of unique label pairs |
| -------------------------- | --------------------- | ---------------------------------------- | ---------------------------- |
| 114 | 709 | 1.67 | 177 |
<!-- markdownlint-enable line-length -->

### Metrics

<!-- markdownlint-disable line-length -->
| Name | Type, Unit | Description |
| ---- | ---------- | ----------- |
| apiserver_audit_event | counter | ALPHA Counter of audit events generated and sent to the audit backend. |
| apiserver_audit_requests_rejected | counter | ALPHA Counter of apiserver requests rejected due to an error in audit logging backend. |
| apiserver_client_certificate_expiration_seconds | histogram | ALPHA Distribution of the remaining lifetime on the certificate used to authenticate a request. |
| apiserver_current_inflight_requests | gauge | STABLE Maximal number of currently used inflight request limit of this apiserver per request kind in last second. |
| apiserver_delegated_authn_request_duration_seconds | histogram | ALPHA Request latency in seconds. Broken down by status code. |
| apiserver_delegated_authn_request | counter | ALPHA Number of HTTP requests partitioned by status code. |
| apiserver_delegated_authz_request_duration_seconds | histogram | ALPHA Request latency in seconds. Broken down by status code. |
| apiserver_delegated_authz_request | counter | ALPHA Number of HTTP requests partitioned by status code. |
| apiserver_envelope_encryption_dek_cache_fill_percent | gauge | ALPHA Percent of the cache slots currently occupied by cached DEKs. |
| apiserver_flowcontrol_read_vs_write_current_requests | histogram | EXPERIMENTAL: ALPHA Observations, at the end of every nanosecond, of the number of requests (as a fraction of the relevant limit) waiting or in regular stage of execution |
| apiserver_flowcontrol_seat_fair_frac | gauge | ALPHA Fair fraction of server's concurrency to allocate to each priority level that can use it |
| apiserver_request_aborts | counter | ALPHA Number of requests which apiserver aborted possibly due to a timeout, for each group, version, verb, resource, subresource and scope |
| apiserver_request_filter_duration_seconds | histogram | ALPHA Request filter latency distribution in seconds, for each filter type |
| apiserver_request_post_timeout | counter | ALPHA Tracks the activity of the request handlers after the associated requests have been timed out by the apiserver |
| apiserver_request_terminations | counter | ALPHA Number of requests which apiserver terminated in self-defense. |
| apiserver_storage_data_key_generation_duration_seconds | histogram | ALPHA Latencies in seconds of data encryption key(DEK) generation operations. |
| apiserver_storage_data_key_generation_failures | counter | ALPHA Total number of failed data encryption key(DEK) generation operations. |
| apiserver_storage_envelope_transformation_cache_misses | counter | ALPHA Total number of cache misses while accessing key decryption key(KEK). |
| apiserver_tls_handshake_errors | counter | ALPHA Number of requests dropped with 'TLS handshake error from' error |
| apiserver_webhooks_x509_insecure_sha1 | counter | ALPHA Counts the number of requests to servers with insecure SHA1 signatures in their serving certificate OR the number of connection failures due to the insecure SHA1 signatures (either/or, based on the runtime environment) |
| apiserver_webhooks_x509_missing_san | counter | ALPHA Counts the number of requests to servers missing SAN extension in their serving certificate OR the number of connection failures due to the lack of x509 certificate SAN extension missing (either/or, based on the runtime environment) |
| authenticated_user_requests | counter | ALPHA Counter of authenticated requests broken out by username. |
| authentication_attempts | counter | ALPHA Counter of authenticated attempts. |
| authentication_duration_seconds | histogram | ALPHA Authentication duration in seconds broken out by result. |
| authentication_token_cache_active_fetch_count | gauge | ALPHA |
| authentication_token_cache_fetch | counter | ALPHA |
| authentication_token_cache_request_duration_seconds | histogram | ALPHA |
| authentication_token_cache_request | counter | ALPHA |
| disabled_metric | counter | ALPHA The count of disabled metrics. |
| event_recorder_total_events_count | counter | ALPHA Total count of events processed by this event recorder per involved object |
| go_cgo_go_to_c_calls_calls | counter | Count of calls made from Go to C by the current process. |
| go_gc_cycles_automatic_gc_cycles | counter | Count of completed GC cycles generated by the Go runtime. |
| go_gc_cycles_forced_gc_cycles | counter | Count of completed GC cycles forced by the application. |
| go_gc_cycles_total_gc_cycles | counter | Count of all completed GC cycles. |
| go_gc_duration_seconds | summary | A summary of the pause duration of garbage collection cycles. |
| go_gc_heap_allocs_by_size_bytes | histogram | Distribution of heap allocations by approximate size. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks. |
| go_gc_heap_allocs_bytes | counter | Cumulative sum of memory allocated to the heap by the application. |
| go_gc_heap_allocs_objects | counter | Cumulative count of heap allocations triggered by the application. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks. |
| go_gc_heap_frees_by_size_bytes | histogram | Distribution of freed heap allocations by approximate size. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks. |
| go_gc_heap_frees_bytes | counter | Cumulative sum of heap memory freed by the garbage collector. |
| go_gc_heap_frees_objects | counter | Cumulative count of heap allocations whose storage was freed by the garbage collector. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks. |
| go_gc_heap_goal_bytes | gauge | Heap size target for the end of the GC cycle. |
| go_gc_heap_objects_objects | gauge | Number of objects, live or unswept, occupying heap memory. |
| go_gc_heap_tiny_allocs_objects | counter | Count of small allocations that are packed together into blocks. These allocations are counted separately from other allocations because each individual allocation is not tracked by the runtime, only their block. Each block is already accounted for in allocs-by-size and frees-by-size. |
| go_gc_limiter_last_enabled_gc_cycle | gauge | GC cycle the last time the GC CPU limiter was enabled. This metric is useful for diagnosing the root cause of an out-of-memory error, because the limiter trades memory for CPU time when the GC's CPU time gets too high. This is most likely to occur with use of SetMemoryLimit. The first GC cycle is cycle 1, so a value of 0 indicates that it was never enabled. |
| go_gc_pauses_seconds | histogram | Distribution individual GC-related stop-the-world pause latencies. |
| go_gc_stack_starting_size_bytes | gauge | The stack size of new goroutines. |
| go_goroutines | gauge | Number of goroutines that currently exist. |
| go_info | gauge | Information about the Go environment. |
| go_memory_classes_heap_free_bytes | gauge | Memory that is completely free and eligible to be returned to the underlying system, but has not been. This metric is the runtime's estimate of free address space that is backed by physical memory. |
| go_memory_classes_heap_objects_bytes | gauge | Memory occupied by live objects and dead objects that have not yet been marked free by the garbage collector. |
| go_memory_classes_heap_released_bytes | gauge | Memory that is completely free and has been returned to the underlying system. This metric is the runtime's estimate of free address space that is still mapped into the process, but is not backed by physical memory. |
| go_memory_classes_heap_stacks_bytes | gauge | Memory allocated from the heap that is reserved for stack space, whether or not it is currently in-use. |
| go_memory_classes_heap_unused_bytes | gauge | Memory that is reserved for heap objects but is not currently used to hold heap objects. |
| go_memory_classes_metadata_mcache_free_bytes | gauge | Memory that is reserved for runtime mcache structures, but not in-use. |
| go_memory_classes_metadata_mcache_inuse_bytes | gauge | Memory that is occupied by runtime mcache structures that are currently being used. |
| go_memory_classes_metadata_mspan_free_bytes | gauge | Memory that is reserved for runtime mspan structures, but not in-use. |
| go_memory_classes_metadata_mspan_inuse_bytes | gauge | Memory that is occupied by runtime mspan structures that are currently being used. |
| go_memory_classes_metadata_other_bytes | gauge | Memory that is reserved for or used to hold runtime metadata. |
| go_memory_classes_os_stacks_bytes | gauge | Stack memory allocated by the underlying operating system. |
| go_memory_classes_other_bytes | gauge | Memory used by execution trace buffers, structures for debugging the runtime, finalizer and profiler specials, and more. |
| go_memory_classes_profiling_buckets_bytes | gauge | Memory that is used by the stack trace hash map used for profiling. |
| go_memory_classes_total_bytes | gauge | All memory mapped by the Go runtime into the current process as read-write. Note that this does not include memory mapped by code called via cgo or via the syscall package. Sum of all metrics in /memory/classes. |
| go_memstats_alloc_bytes | gauge | Number of bytes allocated and still in use. |
| go_memstats_alloc_bytes | counter | Total number of bytes allocated, even if freed. |
| go_memstats_buck_hash_sys_bytes | gauge | Number of bytes used by the profiling bucket hash table. |
| go_memstats_frees | counter | Total number of frees. |
| go_memstats_gc_sys_bytes | gauge | Number of bytes used for garbage collection system metadata. |
| go_memstats_heap_alloc_bytes | gauge | Number of heap bytes allocated and still in use. |
| go_memstats_heap_idle_bytes | gauge | Number of heap bytes waiting to be used. |
| go_memstats_heap_inuse_bytes | gauge | Number of heap bytes that are in use. |
| go_memstats_heap_objects | gauge | Number of allocated objects. |
| go_memstats_heap_released_bytes | gauge | Number of heap bytes released to OS. |
| go_memstats_heap_sys_bytes | gauge | Number of heap bytes obtained from system. |
| go_memstats_last_gc_time_seconds | gauge | Number of seconds since 1970 of last garbage collection. |
| go_memstats_lookups | counter | Total number of pointer lookups. |
| go_memstats_mallocs | counter | Total number of mallocs. |
| go_memstats_mcache_inuse_bytes | gauge | Number of bytes in use by mcache structures. |
| go_memstats_mcache_sys_bytes | gauge | Number of bytes used for mcache structures obtained from system. |
| go_memstats_mspan_inuse_bytes | gauge | Number of bytes in use by mspan structures. |
| go_memstats_mspan_sys_bytes | gauge | Number of bytes used for mspan structures obtained from system. |
| go_memstats_next_gc_bytes | gauge | Number of heap bytes when next garbage collection will take place. |
| go_memstats_other_sys_bytes | gauge | Number of bytes used for other system allocations. |
| go_memstats_stack_inuse_bytes | gauge | Number of bytes in use by the stack allocator. |
| go_memstats_stack_sys_bytes | gauge | Number of bytes obtained from system for stack allocator. |
| go_memstats_sys_bytes | gauge | Number of bytes obtained from system. |
| go_sched_gomaxprocs_threads | gauge | The current runtime.GOMAXPROCS setting, or the number of operating system threads that can execute user-level Go code simultaneously. |
| go_sched_goroutines_goroutines | gauge | Count of live goroutines. |
| go_sched_latencies_seconds | histogram | Distribution of the time goroutines have spent in the scheduler in a runnable state before actually running. |
| go_threads | gauge | Number of OS threads created. |
| hidden_metric | counter | ALPHA The count of hidden metrics. |
| openshift_apiserver_build_info | gauge | ALPHA A metric with a constant '1' value labeled by major, minor, git version, git commit, git tree state, build date, Go version, and compiler from which check-endpoints was built, and platform on which it is running. |
| process_cpu_seconds | counter | Total user and system CPU time spent in seconds. |
| process_max_fds | gauge | Maximum number of open file descriptors. |
| process_open_fds | gauge | Number of open file descriptors. |
| process_resident_memory_bytes | gauge | Resident memory size in bytes. |
| process_start_time_seconds | gauge | Start time of the process since unix epoch in seconds. |
| process_virtual_memory_bytes | gauge | Virtual memory size in bytes. |
| process_virtual_memory_max_bytes | gauge | Maximum amount of virtual memory available in bytes. |
| registered_metric | counter | ALPHA The count of registered metrics broken by stability level and deprecation version. |
| rest_client_exec_plugin_certificate_rotation_age | histogram | ALPHA Histogram of the number of seconds the last auth exec plugin client certificate lived before being rotated. If auth exec plugin client certificates are unused, histogram will contain no data. |
| rest_client_exec_plugin_ttl_seconds | gauge | ALPHA Gauge of the shortest TTL (time-to-live) of the client certificate(s) managed by the auth exec plugin. The value is in seconds until certificate expiry (negative if already expired). If auth exec plugins are unused or manage no TLS certificates, the value will be +INF. |
| rest_client_rate_limiter_duration_seconds | histogram | ALPHA Client side rate limiter latency in seconds. Broken down by verb, and host. |
| rest_client_request_duration_seconds | histogram | ALPHA Request latency in seconds. Broken down by verb, and host. |
| rest_client_request_size_bytes | histogram | ALPHA Request size in bytes. Broken down by verb and host. |
| rest_client_requests | counter | ALPHA Number of HTTP requests, partitioned by status code, method, and host. |
| rest_client_response_size_bytes | histogram | ALPHA Response size in bytes. Broken down by verb and host. |
| workqueue_adds | counter | ALPHA Total number of adds handled by workqueue |
| workqueue_depth | gauge | ALPHA Current depth of workqueue |
| workqueue_longest_running_processor_seconds | gauge | ALPHA How many seconds has the longest running processor for workqueue been running. |
| workqueue_queue_duration_seconds | histogram | ALPHA How long in seconds an item stays in workqueue before being requested. |
| workqueue_retries | counter | ALPHA Total number of retries handled by workqueue |
| workqueue_unfinished_work_seconds | gauge | ALPHA How many seconds of work has done that is in progress and hasn't been observed by work_duration. Large values indicate stuck threads. One can deduce the number of stuck threads by observing the rate at which this increases. |
| workqueue_work_duration_seconds | histogram | ALPHA How long in seconds processing an item from workqueue takes. |
<!-- markdownlint-enable line-length -->

## Cluster-version-operator

### Statistics

<!-- markdownlint-disable line-length -->
| Number of metrics families | Number of time series | Average number of labels per time series | Number of unique label pairs |
| -------------------------- | --------------------- | ---------------------------------------- | ---------------------------- |
| 79 | 360 | 1.95 | 132 |
<!-- markdownlint-enable line-length -->

### Metrics

<!-- markdownlint-disable line-length -->
| Name | Type, Unit | Description |
| ---- | ---------- | ----------- |
| cluster_installer | gauge | Reports info about the installation process and, if applicable, the install tool. The type is either 'openshift-install', indicating that openshift-install was used to install the cluster, or 'other', indicating that an unknown process installed the cluster. The invoker is 'user' by default, but it may be overridden by a consuming tool. The version reported is that of the openshift-install that was used to generate the manifests and, if applicable, provision the infrastructure. |
| cluster_operator_condition_transitions | gauge | Reports the number of times that a condition on a cluster operator changes status |
| cluster_operator_conditions | gauge | Report the conditions for active cluster operators. 0 is False and 1 is True. |
| cluster_operator_payload_errors | gauge | Report the number of errors encountered applying the payload. |
| cluster_operator_up | gauge | 1 if a cluster operator is Available=True.  0 otherwise, including if a cluster operator sets no Available condition.  The 'version' label tracks the 'operator' version.  The 'reason' label is passed through from the Available condition, unless the cluster operator sets no Available condition, in which case NoAvailableCondition is used. |
| cluster_version | gauge | Reports the version of the cluster in terms of seconds since the epoch. Type 'current' is the version being applied and the value is the creation date of the payload. The type 'desired' is returned if s desiredUpdate is set but the operator has not yet updated and the value is the most recent status transition time. The type 'failure' is set if an error is preventing sync or upgrade with the last transition timestamp of the condition. The type 'completed' is the timestamp when the last image was successfully applied. The type 'cluster' is the creation date of the cluster version object and the current version. The type 'updating' is set when the cluster is transitioning to a new version but has not reached the completed state and is the time the update was started. The type 'initial' is set to the oldest entry in the history. The from_version label will be set to the last completed version, the initial version for 'cluster', or empty for 'initial'. |
| cluster_version_capability | gauge | Report currently enabled cluster capabilities.  0 is disabled, and 1 is enabled. |
| cluster_version_operator_update_retrieval_timestamp_seconds | gauge | Reports when updates were last successfully retrieved. |
| cluster_version_payload | gauge | Report the number of entries in the payload. |
| go_cgo_go_to_c_calls_calls | counter | Count of calls made from Go to C by the current process. |
| go_gc_cycles_automatic_gc_cycles | counter | Count of completed GC cycles generated by the Go runtime. |
| go_gc_cycles_forced_gc_cycles | counter | Count of completed GC cycles forced by the application. |
| go_gc_cycles_total_gc_cycles | counter | Count of all completed GC cycles. |
| go_gc_duration_seconds | summary | A summary of the pause duration of garbage collection cycles. |
| go_gc_heap_allocs_by_size_bytes_total | histogram | Distribution of heap allocations by approximate size. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks. |
| go_gc_heap_allocs_bytes | counter | Cumulative sum of memory allocated to the heap by the application. |
| go_gc_heap_allocs_objects | counter | Cumulative count of heap allocations triggered by the application. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks. |
| go_gc_heap_frees_by_size_bytes_total | histogram | Distribution of freed heap allocations by approximate size. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks. |
| go_gc_heap_frees_bytes | counter | Cumulative sum of heap memory freed by the garbage collector. |
| go_gc_heap_frees_objects | counter | Cumulative count of heap allocations whose storage was freed by the garbage collector. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks. |
| go_gc_heap_goal_bytes | gauge | Heap size target for the end of the GC cycle. |
| go_gc_heap_objects_objects | gauge | Number of objects, live or unswept, occupying heap memory. |
| go_gc_heap_tiny_allocs_objects | counter | Count of small allocations that are packed together into blocks. These allocations are counted separately from other allocations because each individual allocation is not tracked by the runtime, only their block. Each block is already accounted for in allocs-by-size and frees-by-size. |
| go_gc_limiter_last_enabled_gc_cycle | gauge | GC cycle the last time the GC CPU limiter was enabled. This metric is useful for diagnosing the root cause of an out-of-memory error, because the limiter trades memory for CPU time when the GC's CPU time gets too high. This is most likely to occur with use of SetMemoryLimit. The first GC cycle is cycle 1, so a value of 0 indicates that it was never enabled. |
| go_gc_pauses_seconds_total | histogram | Distribution individual GC-related stop-the-world pause latencies. |
| go_gc_stack_starting_size_bytes | gauge | The stack size of new goroutines. |
| go_goroutines | gauge | Number of goroutines that currently exist. |
| go_info | gauge | Information about the Go environment. |
| go_memory_classes_heap_free_bytes | gauge | Memory that is completely free and eligible to be returned to the underlying system, but has not been. This metric is the runtime's estimate of free address space that is backed by physical memory. |
| go_memory_classes_heap_objects_bytes | gauge | Memory occupied by live objects and dead objects that have not yet been marked free by the garbage collector. |
| go_memory_classes_heap_released_bytes | gauge | Memory that is completely free and has been returned to the underlying system. This metric is the runtime's estimate of free address space that is still mapped into the process, but is not backed by physical memory. |
| go_memory_classes_heap_stacks_bytes | gauge | Memory allocated from the heap that is reserved for stack space, whether or not it is currently in-use. |
| go_memory_classes_heap_unused_bytes | gauge | Memory that is reserved for heap objects but is not currently used to hold heap objects. |
| go_memory_classes_metadata_mcache_free_bytes | gauge | Memory that is reserved for runtime mcache structures, but not in-use. |
| go_memory_classes_metadata_mcache_inuse_bytes | gauge | Memory that is occupied by runtime mcache structures that are currently being used. |
| go_memory_classes_metadata_mspan_free_bytes | gauge | Memory that is reserved for runtime mspan structures, but not in-use. |
| go_memory_classes_metadata_mspan_inuse_bytes | gauge | Memory that is occupied by runtime mspan structures that are currently being used. |
| go_memory_classes_metadata_other_bytes | gauge | Memory that is reserved for or used to hold runtime metadata. |
| go_memory_classes_os_stacks_bytes | gauge | Stack memory allocated by the underlying operating system. |
| go_memory_classes_other_bytes | gauge | Memory used by execution trace buffers, structures for debugging the runtime, finalizer and profiler specials, and more. |
| go_memory_classes_profiling_buckets_bytes | gauge | Memory that is used by the stack trace hash map used for profiling. |
| go_memory_classes_total_bytes | gauge | All memory mapped by the Go runtime into the current process as read-write. Note that this does not include memory mapped by code called via cgo or via the syscall package. Sum of all metrics in /memory/classes. |
| go_memstats_alloc_bytes | gauge | Number of bytes allocated and still in use. |
| go_memstats_alloc_bytes | counter | Total number of bytes allocated, even if freed. |
| go_memstats_buck_hash_sys_bytes | gauge | Number of bytes used by the profiling bucket hash table. |
| go_memstats_frees | counter | Total number of frees. |
| go_memstats_gc_cpu_fraction | gauge | The fraction of this program's available CPU time used by the GC since the program started. |
| go_memstats_gc_sys_bytes | gauge | Number of bytes used for garbage collection system metadata. |
| go_memstats_heap_alloc_bytes | gauge | Number of heap bytes allocated and still in use. |
| go_memstats_heap_idle_bytes | gauge | Number of heap bytes waiting to be used. |
| go_memstats_heap_inuse_bytes | gauge | Number of heap bytes that are in use. |
| go_memstats_heap_objects | gauge | Number of allocated objects. |
| go_memstats_heap_released_bytes | gauge | Number of heap bytes released to OS. |
| go_memstats_heap_sys_bytes | gauge | Number of heap bytes obtained from system. |
| go_memstats_last_gc_time_seconds | gauge | Number of seconds since 1970 of last garbage collection. |
| go_memstats_lookups | counter | Total number of pointer lookups. |
| go_memstats_mallocs | counter | Total number of mallocs. |
| go_memstats_mcache_inuse_bytes | gauge | Number of bytes in use by mcache structures. |
| go_memstats_mcache_sys_bytes | gauge | Number of bytes used for mcache structures obtained from system. |
| go_memstats_mspan_inuse_bytes | gauge | Number of bytes in use by mspan structures. |
| go_memstats_mspan_sys_bytes | gauge | Number of bytes used for mspan structures obtained from system. |
| go_memstats_next_gc_bytes | gauge | Number of heap bytes when next garbage collection will take place. |
| go_memstats_other_sys_bytes | gauge | Number of bytes used for other system allocations. |
| go_memstats_stack_inuse_bytes | gauge | Number of bytes in use by the stack allocator. |
| go_memstats_stack_sys_bytes | gauge | Number of bytes obtained from system for stack allocator. |
| go_memstats_sys_bytes | gauge | Number of bytes obtained from system. |
| go_sched_gomaxprocs_threads | gauge | The current runtime.GOMAXPROCS setting, or the number of operating system threads that can execute user-level Go code simultaneously. |
| go_sched_goroutines_goroutines | gauge | Count of live goroutines. |
| go_sched_latencies_seconds | histogram | Distribution of the time goroutines have spent in the scheduler in a runnable state before actually running. |
| go_threads | gauge | Number of OS threads created. |
| process_cpu_seconds | counter | Total user and system CPU time spent in seconds. |
| process_max_fds | gauge | Maximum number of open file descriptors. |
| process_open_fds | gauge | Number of open file descriptors. |
| process_resident_memory_bytes | gauge | Resident memory size in bytes. |
| process_start_time_seconds | gauge | Start time of the process since unix epoch in seconds. |
| process_virtual_memory_bytes | gauge | Virtual memory size in bytes. |
| process_virtual_memory_max_bytes | gauge | Maximum amount of virtual memory available in bytes. |
| promhttp_metric_handler_requests_in_flight | gauge | Current number of scrapes being served. |
| promhttp_metric_handler_requests | counter | Total number of scrapes by HTTP status code. |
<!-- markdownlint-enable line-length -->

## Openshift-HAProxy

### Statistics

<!-- markdownlint-disable line-length -->
| Number of metrics families | Number of time series | Average number of labels per time series | Number of unique label pairs |
| -------------------------- | --------------------- | ---------------------------------------- | ---------------------------- |
| 92 | 1821 | 4.2 | 199 |
<!-- markdownlint-enable line-length -->

### Metrics

<!-- markdownlint-disable line-length -->
| Name | Type, Unit | Description |
| ---- | ---------- | ----------- |
| go_gc_duration_seconds | summary | A summary of the pause duration of garbage collection cycles. |
| go_goroutines | gauge | Number of goroutines that currently exist. |
| go_info | gauge | Information about the Go environment. |
| go_memstats_alloc_bytes | gauge | Number of bytes allocated and still in use. |
| go_memstats_alloc_bytes | counter | Total number of bytes allocated, even if freed. |
| go_memstats_buck_hash_sys_bytes | gauge | Number of bytes used by the profiling bucket hash table. |
| go_memstats_frees | counter | Total number of frees. |
| go_memstats_gc_sys_bytes | gauge | Number of bytes used for garbage collection system metadata. |
| go_memstats_heap_alloc_bytes | gauge | Number of heap bytes allocated and still in use. |
| go_memstats_heap_idle_bytes | gauge | Number of heap bytes waiting to be used. |
| go_memstats_heap_inuse_bytes | gauge | Number of heap bytes that are in use. |
| go_memstats_heap_objects | gauge | Number of allocated objects. |
| go_memstats_heap_released_bytes | gauge | Number of heap bytes released to OS. |
| go_memstats_heap_sys_bytes | gauge | Number of heap bytes obtained from system. |
| go_memstats_last_gc_time_seconds | gauge | Number of seconds since 1970 of last garbage collection. |
| go_memstats_lookups | counter | Total number of pointer lookups. |
| go_memstats_mallocs | counter | Total number of mallocs. |
| go_memstats_mcache_inuse_bytes | gauge | Number of bytes in use by mcache structures. |
| go_memstats_mcache_sys_bytes | gauge | Number of bytes used for mcache structures obtained from system. |
| go_memstats_mspan_inuse_bytes | gauge | Number of bytes in use by mspan structures. |
| go_memstats_mspan_sys_bytes | gauge | Number of bytes used for mspan structures obtained from system. |
| go_memstats_next_gc_bytes | gauge | Number of heap bytes when next garbage collection will take place. |
| go_memstats_other_sys_bytes | gauge | Number of bytes used for other system allocations. |
| go_memstats_stack_inuse_bytes | gauge | Number of bytes in use by the stack allocator. |
| go_memstats_stack_sys_bytes | gauge | Number of bytes obtained from system for stack allocator. |
| go_memstats_sys_bytes | gauge | Number of bytes obtained from system. |
| go_threads | gauge | Number of OS threads created. |
| haproxy_backend_bytes_in_total | gauge | Current total of incoming bytes. |
| haproxy_backend_bytes_out_total | gauge | Current total of outgoing bytes. |
| haproxy_backend_connection_errors_total | gauge | Total of connection errors. |
| haproxy_backend_connections_reused_total | gauge | Total number of connections reused. |
| haproxy_backend_connections_total | gauge | Total number of connections. |
| haproxy_backend_current_queue | gauge | Current number of queued requests not assigned to any server. |
| haproxy_backend_current_session_rate | gauge | Current number of sessions per second over last elapsed second. |
| haproxy_backend_current_sessions | gauge | Current number of active sessions. |
| haproxy_backend_http_average_connect_latency_milliseconds | gauge | Average connect latency of the last 1024 requests in milliseconds. |
| haproxy_backend_http_average_queue_latency_milliseconds | gauge | Average latency to be dequeued of the last 1024 requests in milliseconds. |
| haproxy_backend_http_average_response_latency_milliseconds | gauge | Average response latency of the last 1024 requests in milliseconds. |
| haproxy_backend_http_responses_total | gauge | Total of HTTP responses. |
| haproxy_backend_max_session_rate | gauge | Maximum number of sessions per second. |
| haproxy_backend_max_sessions | gauge | Maximum observed number of active sessions. |
| haproxy_backend_response_errors_total | gauge | Total of response errors. |
| haproxy_backend_up | gauge | Current health status of the backend (1 = UP, 0 = DOWN). |
| haproxy_exporter_csv_parse_failures | counter | Number of errors while parsing CSV. |
| haproxy_exporter_scrape_interval | gauge | The time in seconds before another scrape is allowed, proportional to size of data. |
| haproxy_exporter_server_threshold | gauge | Number of servers tracked and the current threshold value. |
| haproxy_exporter_total_scrapes | counter | Current total HAProxy scrapes. |
| haproxy_frontend_bytes_in_total | gauge | Current total of incoming bytes. |
| haproxy_frontend_bytes_out_total | gauge | Current total of outgoing bytes. |
| haproxy_frontend_connections_total | gauge | Total number of connections. |
| haproxy_frontend_current_session_rate | gauge | Current number of sessions per second over last elapsed second. |
| haproxy_frontend_current_sessions | gauge | Current number of active sessions. |
| haproxy_frontend_http_responses_total | gauge | Total of HTTP responses. |
| haproxy_frontend_max_session_rate | gauge | Maximum observed number of sessions per second. |
| haproxy_frontend_max_sessions | gauge | Maximum observed number of active sessions. |
| haproxy_process_cpu_seconds | counter | Total user and system CPU time spent in seconds. |
| haproxy_process_max_fds | gauge | Maximum number of open file descriptors. |
| haproxy_process_resident_memory_bytes | gauge | Resident memory size in bytes. |
| haproxy_process_start_time_seconds | gauge | Start time of the process since unix epoch in seconds. |
| haproxy_process_virtual_memory_bytes | gauge | Virtual memory size in bytes. |
| haproxy_process_virtual_memory_max_bytes | gauge | Maximum amount of virtual memory available in bytes. |
| haproxy_server_bytes_in_total | gauge | Current total of incoming bytes. |
| haproxy_server_bytes_out_total | gauge | Current total of outgoing bytes. |
| haproxy_server_check_failures_total | gauge | Total number of failed health checks. |
| haproxy_server_connection_errors_total | gauge | Total of connection errors. |
| haproxy_server_connections_reused_total | gauge | Total number of connections reused. |
| haproxy_server_connections_total | gauge | Total number of connections. |
| haproxy_server_current_queue | gauge | Current number of queued requests assigned to this server. |
| haproxy_server_current_session_rate | gauge | Current number of sessions per second over last elapsed second. |
| haproxy_server_current_sessions | gauge | Current number of active sessions. |
| haproxy_server_downtime_seconds_total | gauge | Total downtime in seconds. |
| haproxy_server_http_average_connect_latency_milliseconds | gauge | Average connect latency of the last 1024 requests in milliseconds. |
| haproxy_server_http_average_queue_latency_milliseconds | gauge | Average latency to be dequeued of the last 1024 requests in milliseconds. |
| haproxy_server_http_average_response_latency_milliseconds | gauge | Average response latency of the last 1024 requests in milliseconds. |
| haproxy_server_http_responses_total | gauge | Total of HTTP responses. |
| haproxy_server_max_session_rate | gauge | Maximum observed number of sessions per second. |
| haproxy_server_max_sessions | gauge | Maximum observed number of active sessions. |
| haproxy_server_response_errors_total | gauge | Total of response errors. |
| haproxy_server_up | gauge | Current health status of the server (1 = UP, 0 = DOWN). |
| haproxy_up | gauge | Was the last scrape of haproxy successful. |
| process_cpu_seconds | counter | Total user and system CPU time spent in seconds. |
| process_max_fds | gauge | Maximum number of open file descriptors. |
| process_open_fds | gauge | Number of open file descriptors. |
| process_resident_memory_bytes | gauge | Resident memory size in bytes. |
| process_start_time_seconds | gauge | Start time of the process since unix epoch in seconds. |
| process_virtual_memory_bytes | gauge | Virtual memory size in bytes. |
| process_virtual_memory_max_bytes | gauge | Maximum amount of virtual memory available in bytes. |
| promhttp_metric_handler_requests_in_flight | gauge | Current number of scrapes being served. |
| promhttp_metric_handler_requests | counter | Total number of scrapes by HTTP status code. |
| template_router_reload_failure | gauge | Metric to track the status of the most recent HAProxy reload |
| template_router_reload_seconds | summary | Measures the time spent reloading the router in seconds. |
| template_router_write_config_seconds | summary | Measures the time spent writing out the router configuration to disk in seconds. |
<!-- markdownlint-enable line-length -->