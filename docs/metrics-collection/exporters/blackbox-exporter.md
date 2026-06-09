This document describes the metrics list and how to collect them from the Blackbox-exporter.

# Metrics

The Blackbox-exporter is the exporter and already exposes its metrics in Prometheus format.

| Name          | Metrics Port | Metrics Endpoint      | Need Exporter? | Is Exporter Third Party? |
| ------------- | ------------ | --------------------- | -------------- | ------------------------ |
| Self metrics  | `9115`       | `/metrics`            | No             | N/A                      |
| Probe metrics | `9115`       | `/probe` + parameters | No             | N/A                      |

## How to Collect

The `blackbox-exporter` provides two endpoints on port `9115` with metrics:

* `/metrics` with self metrics
* `/probe` that allows to send some parameters, execute a probe and receive metrics as result of probe

By default, `blackbox-exporter` has no authentication for it endpoints.

### Self metrics

Config `ServiceMonitor` for `prometheus-operator` to collect Blackbox-exporter self metrics:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  labels:
    app.kubernetes.io/component: monitoring
  name: blackbox-exporter-self-service-monitor
spec:
  endpoints:
    - interval: 30s
      scrapeTimeout: 10s
      path: /metrics
      port: http
      scheme: http
  jobLabel: blackbox-exporter-self
  selector:
    matchExpressions:
      - key: app
        operator: In
        values:
          - blackbox-exporter
```

To collect (or just to check) self-metrics manually you can use the following command:

```bash
curl -v -k -L http://<blackbox_ip_or_dns>:9115/metrics
```

or

```bash
wget -O - http://<blackbox_ip_or_dns>:9115/metrics
```

### Probe

To configure the probe for URL need to configure `Probe` Custom Resource for `prometheus-operator` like the following.

For static list of URLs:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: Probe
metadata:
  name: blackbox-static-urls-probe
  labels:
    app.kubernetes.io/component: monitoring  # Mandatory label
spec:
  jobName: http-get
  interval: 30s
  module: http_2xx
  prober:
    url: blackbox-exporter.monitoring.svc:9115
    scheme: http
    path: /probe
  targets:
    staticConfig:
      static:
        - 'https://example.com'
        - 'https://google.com'
```

For Ingresses (Ingress host will be discover automatically):

```yaml
apiVersion: monitoring.coreos.com/v1
kind: Probe
metadata:
  name: blackbox-ingress-probe
  labels:
    app.kubernetes.io/component: monitoring  # Mandatory label
spec:
  jobName: http-get
  interval: 30s
  module: http_2xx
  prober:
    url: blackbox-exporter.monitoring.svc:9115
    scheme: http
    path: /probe
  targets:
    ingress:
      selector:
        matchLabels:
          name: prometheus   # label that will be use to discover Ingresses
      namespaceSelector:
        matchNames:
          - monitoring       # namespace in which Ingresses will be discover
```

To collect (or just to check) probe metrics manually you can use the following command:

```bash
curl -v -k -L http://<blackbox_ip_or_dns>:9115/probe?target=google.com&module=http_2xx
```

or

```bash
wget -O - http://<blackbox_ip_or_dns>:9115/probe?target=google.com&module=http_2xx
```

## Metrics List

This section describes metrics which be can be collected from `blackbox-exporter`.

### Self Metrics

```prometheus
# HELP blackbox_exporter_build_info A metric with a constant '1' value labeled by version, revision, branch, and goversion from which blackbox_exporter was built.
# TYPE blackbox_exporter_build_info gauge
blackbox_exporter_build_info{branch="HEAD",goversion="go1.18.5",revision="0bbd65d1264722f7afb87a72ec4128b9214e5840",version="0.22.0"} 1
# HELP blackbox_exporter_config_last_reload_success_timestamp_seconds Timestamp of the last successful configuration reload.
# TYPE blackbox_exporter_config_last_reload_success_timestamp_seconds gauge
blackbox_exporter_config_last_reload_success_timestamp_seconds 1.6760364991695104e+09
# HELP blackbox_exporter_config_last_reload_successful Blackbox exporter config loaded successfully.
# TYPE blackbox_exporter_config_last_reload_successful gauge
blackbox_exporter_config_last_reload_successful 1
# HELP blackbox_module_unknown_total Count of unknown modules requested by probes
# TYPE blackbox_module_unknown_total counter
blackbox_module_unknown_total 0
# HELP go_gc_duration_seconds A summary of the pause duration of garbage collection cycles.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 3.5996e-05
go_gc_duration_seconds{quantile="0.25"} 7.6105e-05
go_gc_duration_seconds{quantile="0.5"} 0.000128373
go_gc_duration_seconds{quantile="0.75"} 0.000308447
go_gc_duration_seconds{quantile="1"} 0.025982248
go_gc_duration_seconds_sum 107.893983213
go_gc_duration_seconds_count 101280
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 12
# HELP go_info Information about the Go environment.
# TYPE go_info gauge
go_info{version="go1.18.5"} 1
# HELP go_memstats_alloc_bytes Number of bytes allocated and still in use.
# TYPE go_memstats_alloc_bytes gauge
go_memstats_alloc_bytes 9.65204e+06
# HELP go_memstats_alloc_bytes_total Total number of bytes allocated, even if freed.
# TYPE go_memstats_alloc_bytes_total counter
go_memstats_alloc_bytes_total 2.80681919632e+11
# HELP go_memstats_buck_hash_sys_bytes Number of bytes used by the profiling bucket hash table.
# TYPE go_memstats_buck_hash_sys_bytes gauge
go_memstats_buck_hash_sys_bytes 1.823603e+06
# HELP go_memstats_frees_total Total number of frees.
# TYPE go_memstats_frees_total counter
go_memstats_frees_total 8.89636693e+08
# HELP go_memstats_gc_sys_bytes Number of bytes used for garbage collection system metadata.
# TYPE go_memstats_gc_sys_bytes gauge
go_memstats_gc_sys_bytes 5.188568e+06
# HELP go_memstats_heap_alloc_bytes Number of heap bytes allocated and still in use.
# TYPE go_memstats_heap_alloc_bytes gauge
go_memstats_heap_alloc_bytes 9.65204e+06
# HELP go_memstats_heap_idle_bytes Number of heap bytes waiting to be used.
# TYPE go_memstats_heap_idle_bytes gauge
go_memstats_heap_idle_bytes 9.551872e+06
# HELP go_memstats_heap_inuse_bytes Number of heap bytes that are in use.
# TYPE go_memstats_heap_inuse_bytes gauge
go_memstats_heap_inuse_bytes 1.0174464e+07
# HELP go_memstats_heap_objects Number of allocated objects.
# TYPE go_memstats_heap_objects gauge
go_memstats_heap_objects 32833
# HELP go_memstats_heap_released_bytes Number of heap bytes released to OS.
# TYPE go_memstats_heap_released_bytes gauge
go_memstats_heap_released_bytes 9.125888e+06
# HELP go_memstats_heap_sys_bytes Number of heap bytes obtained from system.
# TYPE go_memstats_heap_sys_bytes gauge
go_memstats_heap_sys_bytes 1.9726336e+07
# HELP go_memstats_last_gc_time_seconds Number of seconds since 1970 of last garbage collection.
# TYPE go_memstats_last_gc_time_seconds gauge
go_memstats_last_gc_time_seconds 1.6836131746554873e+09
# HELP go_memstats_lookups_total Total number of pointer lookups.
# TYPE go_memstats_lookups_total counter
go_memstats_lookups_total 0
# HELP go_memstats_mallocs_total Total number of mallocs.
# TYPE go_memstats_mallocs_total counter
go_memstats_mallocs_total 8.89669526e+08
# HELP go_memstats_mcache_inuse_bytes Number of bytes in use by mcache structures.
# TYPE go_memstats_mcache_inuse_bytes gauge
go_memstats_mcache_inuse_bytes 7200
# HELP go_memstats_mcache_sys_bytes Number of bytes used for mcache structures obtained from system.
# TYPE go_memstats_mcache_sys_bytes gauge
go_memstats_mcache_sys_bytes 15600
# HELP go_memstats_mspan_inuse_bytes Number of bytes in use by mspan structures.
# TYPE go_memstats_mspan_inuse_bytes gauge
go_memstats_mspan_inuse_bytes 138720
# HELP go_memstats_mspan_sys_bytes Number of bytes used for mspan structures obtained from system.
# TYPE go_memstats_mspan_sys_bytes gauge
go_memstats_mspan_sys_bytes 163200
# HELP go_memstats_next_gc_bytes Number of heap bytes when next garbage collection will take place.
# TYPE go_memstats_next_gc_bytes gauge
go_memstats_next_gc_bytes 1.0929376e+07
# HELP go_memstats_other_sys_bytes Number of bytes used for other system allocations.
# TYPE go_memstats_other_sys_bytes gauge
go_memstats_other_sys_bytes 1.280589e+06
# HELP go_memstats_stack_inuse_bytes Number of bytes in use by the stack allocator.
# TYPE go_memstats_stack_inuse_bytes gauge
go_memstats_stack_inuse_bytes 1.245184e+06
# HELP go_memstats_stack_sys_bytes Number of bytes obtained from system for stack allocator.
# TYPE go_memstats_stack_sys_bytes gauge
go_memstats_stack_sys_bytes 1.245184e+06
# HELP go_memstats_sys_bytes Number of bytes obtained from system.
# TYPE go_memstats_sys_bytes gauge
go_memstats_sys_bytes 2.944308e+07
# HELP go_threads Number of OS threads created.
# TYPE go_threads gauge
go_threads 13
# HELP net_conntrack_dialer_conn_attempted_total Total number of connections attempted by the given dialer a given name.
# TYPE net_conntrack_dialer_conn_attempted_total counter
net_conntrack_dialer_conn_attempted_total{dialer_name="http_probe"} 376520
# HELP net_conntrack_dialer_conn_closed_total Total number of connections closed which originated from the dialer of a given name.
# TYPE net_conntrack_dialer_conn_closed_total counter
net_conntrack_dialer_conn_closed_total{dialer_name="http_probe"} 376380
# HELP net_conntrack_dialer_conn_established_total Total number of connections successfully established by the given dialer a given name.
# TYPE net_conntrack_dialer_conn_established_total counter
net_conntrack_dialer_conn_established_total{dialer_name="http_probe"} 376379
# HELP net_conntrack_dialer_conn_failed_total Total number of connections failed to dial by the dialer a given name.
# TYPE net_conntrack_dialer_conn_failed_total counter
net_conntrack_dialer_conn_failed_total{dialer_name="http_probe",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="http_probe",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="http_probe",reason="timeout"} 141
net_conntrack_dialer_conn_failed_total{dialer_name="http_probe",reason="unknown"} 141
# HELP process_cpu_seconds_total Total user and system CPU time spent in seconds.
# TYPE process_cpu_seconds_total counter
process_cpu_seconds_total 4170.78
# HELP process_max_fds Maximum number of open file descriptors.
# TYPE process_max_fds gauge
process_max_fds 1.048576e+06
# HELP process_open_fds Number of open file descriptors.
# TYPE process_open_fds gauge
process_open_fds 11
# HELP process_resident_memory_bytes Resident memory size in bytes.
# TYPE process_resident_memory_bytes gauge
process_resident_memory_bytes 2.0688896e+07
# HELP process_start_time_seconds Start time of the process since unix epoch in seconds.
# TYPE process_start_time_seconds gauge
process_start_time_seconds 1.67603649734e+09
# HELP process_virtual_memory_bytes Virtual memory size in bytes.
# TYPE process_virtual_memory_bytes gauge
process_virtual_memory_bytes 7.37415168e+08
# HELP process_virtual_memory_max_bytes Maximum amount of virtual memory available in bytes.
# TYPE process_virtual_memory_max_bytes gauge
process_virtual_memory_max_bytes 1.8446744073709552e+19
# HELP promhttp_metric_handler_requests_in_flight Current number of scrapes being served.
# TYPE promhttp_metric_handler_requests_in_flight gauge
promhttp_metric_handler_requests_in_flight 1
# HELP promhttp_metric_handler_requests_total Total number of scrapes by HTTP status code.
# TYPE promhttp_metric_handler_requests_total counter
promhttp_metric_handler_requests_total{code="200"} 252555
promhttp_metric_handler_requests_total{code="500"} 0
promhttp_metric_handler_requests_total{code="503"} 0
```

### Probe Metrics

The list of custom metrics depends on json-exporter configuration.