This document describes the metrics list and how to collect them from Alertmanager.

# Metrics

Alertmanager already exposes its metrics in Prometheus format and doesn't require to use of specific exporters.

| Application | Metrics Port | Metrics Endpoint      | Need Exporter? | Auth? | Is Exporter Third Party? |
| ----------- | ------------ | --------------------- | -------------- | ----- | ------------------------ |
| Prometheus  | `7979`       | `/metrics`            | No             | No    | N/A                      |
| Prometheus  | `7979`       | `/probe` + parameters | No             | No    | N/A                      |

# Use cases

This exporter can be use and be useful in the following cases.

First, you have any json structure that contains any useful information in metrics.

For example, if you have a catalog with items and them versions in the catalog. And you want to track
items from the catalog and them versions as metrics.

Second, for services that expose metrics, but doesn't support the Prometheus format. Using the `json-exporter`
can allow to avoid using third party exporters (and in some cases such exporters can absent).

Third, to get data from any APIs, parse output and get useful data as metrics.

## How to Collect

The `json-exporter` provides two endpoints on port `7979` with metrics:

* `/metrics` with self metrics
* `/probe` that allows to send some parameters, execute a probe and receive metrics as result of probe

By default, `json-exporter` has no authentication for it endpoints.

### Self metrics

To collect self-metrics need configure `ServiceMonitor`:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  labels:
    app.kubernetes.io/component: monitoring
  name: json-exporter-self-service-monitor
spec:
  endpoints:
    - interval: 30s
      scrapeTimeout: 10s
      path: /metrics
      port: 7979
      scheme: http
  jobLabel: json-exporter-self
  selector:
    matchExpressions:
      - key: app
        operator: In
        values:
          - json-exporter
```

To collect (or just to check) metrics manually you can use the following command:

```bash
curl -v -k -L http://<json_exporter_ip_or_dns>:7979/metrics
```

or

```bash
wget -O - http://<json_exporter_ip_or_dns>:7979/metrics
```

### Probes

The `prometheus-operator` provide a special type of Custom Resources for exporters that execute some probes
(like `json-exporter` or `blackbox-exporter`) with name `Probe`.

For example:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: Probe
metadata:
  name: json-exporter-probe
  labels:
    app.kubernetes.io/component: monitoring
spec:
  jobName: http-get
  interval: 30s
  module: default
  prober:
    url: json-exporter.monitoring.svc:7979
    scheme: http
    path: /probe
  targets:
    staticConfig:
      static:
        - 'http://example.com/data.json'
        - 'http://test-service.com/example/data.json'
```

But before using `json-exporter` you need to configure one or some modules. For example,
with using Monitoring deploy parameters it will looks like as following:

```yaml
jsonExporter:
 config:
    modules:
      default:
        metrics:
        - name: example_global_value
          path: "{ .counter }"
          help: Example of a top-level global value scrape in the json
          labels:
            environment: beta                       # static label
            location: "planet-{.location}"          # dynamic label
        - name: example_timestamped_value
          path: '{ .values[?(@.state == "INACTIVE")] }'
          epochTimestamp: "{ .timestamp }"
          help: Example of a timestamped value scrape in the json
          labels:
            environment: beta                       # static label
        - name: example_value
          type: object
          help: Example of sub-level value scrapes from a json
          path: '{.values[?(@.state == "ACTIVE")]}'
          labels:
            environment: beta                        # static label
            id: '{.id}'                              # dynamic label
          values:
            active: 1                                # static value
            count: '{.count}'                        # dynamic value
            boolean: '{.some_boolean}'
        headers:
          X-Dummy: my-test-header
```

More examples you can find in official repository:

* [https://github.com/prometheus-community/json_exporter](https://github.com/prometheus-community/json_exporter)

To collect (or just to check) metrics manually you can use the following command:

```bash
curl -v -k -L http://<json_exporter_ip_or_dns>:7979/probe?module=default&target=http://example.com/data.json
```

or

```bash
wget -O - http://<json_exporter_ip_or_dns>:7979/probe?module=default&target=http://example.com/data.json
```

## Metrics list

This section contains list of raw metrics of this exporter.

### Self metrics

```prometheus
# HELP go_gc_duration_seconds A summary of the pause duration of garbage collection cycles.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 3.6671e-05
go_gc_duration_seconds{quantile="0.25"} 8.3025e-05
go_gc_duration_seconds{quantile="0.5"} 0.000128741
go_gc_duration_seconds{quantile="0.75"} 0.000206546
go_gc_duration_seconds{quantile="1"} 0.032284481
go_gc_duration_seconds_sum 9.216077959
go_gc_duration_seconds_count 15273
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 7
# HELP go_info Information about the Go environment.
# TYPE go_info gauge
go_info{version="go1.18.3"} 1
# HELP go_memstats_alloc_bytes Number of bytes allocated and still in use.
# TYPE go_memstats_alloc_bytes gauge
go_memstats_alloc_bytes 3.037272e+06
# HELP go_memstats_alloc_bytes_total Total number of bytes allocated, even if freed.
# TYPE go_memstats_alloc_bytes_total counter
go_memstats_alloc_bytes_total 2.2296786008e+10
# HELP go_memstats_buck_hash_sys_bytes Number of bytes used by the profiling bucket hash table.
# TYPE go_memstats_buck_hash_sys_bytes gauge
go_memstats_buck_hash_sys_bytes 4650
# HELP go_memstats_frees_total Total number of frees.
# TYPE go_memstats_frees_total counter
go_memstats_frees_total 8.1974596e+07
# HELP go_memstats_gc_sys_bytes Number of bytes used for garbage collection system metadata.
# TYPE go_memstats_gc_sys_bytes gauge
go_memstats_gc_sys_bytes 4.94004e+06
# HELP go_memstats_heap_alloc_bytes Number of heap bytes allocated and still in use.
# TYPE go_memstats_heap_alloc_bytes gauge
go_memstats_heap_alloc_bytes 3.037272e+06
# HELP go_memstats_heap_idle_bytes Number of heap bytes waiting to be used.
# TYPE go_memstats_heap_idle_bytes gauge
go_memstats_heap_idle_bytes 7.430144e+06
# HELP go_memstats_heap_inuse_bytes Number of heap bytes that are in use.
# TYPE go_memstats_heap_inuse_bytes gauge
go_memstats_heap_inuse_bytes 4.3008e+06
# HELP go_memstats_heap_objects Number of allocated objects.
# TYPE go_memstats_heap_objects gauge
go_memstats_heap_objects 18626
# HELP go_memstats_heap_released_bytes Number of heap bytes released to OS.
# TYPE go_memstats_heap_released_bytes gauge
go_memstats_heap_released_bytes 5.750784e+06
# HELP go_memstats_heap_sys_bytes Number of heap bytes obtained from system.
# TYPE go_memstats_heap_sys_bytes gauge
go_memstats_heap_sys_bytes 1.1730944e+07
# HELP go_memstats_last_gc_time_seconds Number of seconds since 1970 of last garbage collection.
# TYPE go_memstats_last_gc_time_seconds gauge
go_memstats_last_gc_time_seconds 1.683610762165975e+09
# HELP go_memstats_lookups_total Total number of pointer lookups.
# TYPE go_memstats_lookups_total counter
go_memstats_lookups_total 0
# HELP go_memstats_mallocs_total Total number of mallocs.
# TYPE go_memstats_mallocs_total counter
go_memstats_mallocs_total 8.1993222e+07
# HELP go_memstats_mcache_inuse_bytes Number of bytes in use by mcache structures.
# TYPE go_memstats_mcache_inuse_bytes gauge
go_memstats_mcache_inuse_bytes 7200
# HELP go_memstats_mcache_sys_bytes Number of bytes used for mcache structures obtained from system.
# TYPE go_memstats_mcache_sys_bytes gauge
go_memstats_mcache_sys_bytes 15600
# HELP go_memstats_mspan_inuse_bytes Number of bytes in use by mspan structures.
# TYPE go_memstats_mspan_inuse_bytes gauge
go_memstats_mspan_inuse_bytes 113016
# HELP go_memstats_mspan_sys_bytes Number of bytes used for mspan structures obtained from system.
# TYPE go_memstats_mspan_sys_bytes gauge
go_memstats_mspan_sys_bytes 146880
# HELP go_memstats_next_gc_bytes Number of heap bytes when next garbage collection will take place.
# TYPE go_memstats_next_gc_bytes gauge
go_memstats_next_gc_bytes 4.194304e+06
# HELP go_memstats_other_sys_bytes Number of bytes used for other system allocations.
# TYPE go_memstats_other_sys_bytes gauge
go_memstats_other_sys_bytes 1.40011e+06
# HELP go_memstats_stack_inuse_bytes Number of bytes in use by the stack allocator.
# TYPE go_memstats_stack_inuse_bytes gauge
go_memstats_stack_inuse_bytes 851968
# HELP go_memstats_stack_sys_bytes Number of bytes obtained from system for stack allocator.
# TYPE go_memstats_stack_sys_bytes gauge
go_memstats_stack_sys_bytes 851968
# HELP go_memstats_sys_bytes Number of bytes obtained from system.
# TYPE go_memstats_sys_bytes gauge
go_memstats_sys_bytes 1.9090192e+07
# HELP go_threads Number of OS threads created.
# TYPE go_threads gauge
go_threads 11
# HELP net_conntrack_dialer_conn_attempted_total Total number of connections attempted by the given dialer a given name.
# TYPE net_conntrack_dialer_conn_attempted_total counter
net_conntrack_dialer_conn_attempted_total{dialer_name="fetch_json"} 6855
# HELP net_conntrack_dialer_conn_closed_total Total number of connections closed which originated from the dialer of a given name.
# TYPE net_conntrack_dialer_conn_closed_total counter
net_conntrack_dialer_conn_closed_total{dialer_name="fetch_json"} 0
# HELP net_conntrack_dialer_conn_established_total Total number of connections successfully established by the given dialer a given name.
# TYPE net_conntrack_dialer_conn_established_total counter
net_conntrack_dialer_conn_established_total{dialer_name="fetch_json"} 0
# HELP net_conntrack_dialer_conn_failed_total Total number of connections failed to dial by the dialer a given name.
# TYPE net_conntrack_dialer_conn_failed_total counter
net_conntrack_dialer_conn_failed_total{dialer_name="fetch_json",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="fetch_json",reason="resolution"} 6855
net_conntrack_dialer_conn_failed_total{dialer_name="fetch_json",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="fetch_json",reason="unknown"} 0
# HELP process_cpu_seconds_total Total user and system CPU time spent in seconds.
# TYPE process_cpu_seconds_total counter
process_cpu_seconds_total 265.67
# HELP process_max_fds Maximum number of open file descriptors.
# TYPE process_max_fds gauge
process_max_fds 1.048576e+06
# HELP process_open_fds Number of open file descriptors.
# TYPE process_open_fds gauge
process_open_fds 10
# HELP process_resident_memory_bytes Resident memory size in bytes.
# TYPE process_resident_memory_bytes gauge
process_resident_memory_bytes 1.3000704e+07
# HELP process_start_time_seconds Start time of the process since unix epoch in seconds.
# TYPE process_start_time_seconds gauge
process_start_time_seconds 1.6831994495e+09
# HELP process_virtual_memory_bytes Virtual memory size in bytes.
# TYPE process_virtual_memory_bytes gauge
process_virtual_memory_bytes 7.32217344e+08
# HELP process_virtual_memory_max_bytes Maximum amount of virtual memory available in bytes.
# TYPE process_virtual_memory_max_bytes gauge
process_virtual_memory_max_bytes 1.8446744073709552e+19
# HELP promhttp_metric_handler_requests_in_flight Current number of scrapes being served.
# TYPE promhttp_metric_handler_requests_in_flight gauge
promhttp_metric_handler_requests_in_flight 1
# HELP promhttp_metric_handler_requests_total Total number of scrapes by HTTP status code.
# TYPE promhttp_metric_handler_requests_total counter
promhttp_metric_handler_requests_total{code="200"} 95978
promhttp_metric_handler_requests_total{code="500"} 0
promhttp_metric_handler_requests_total{code="503"} 0
```

### Probe metrics

The list of custom metrics depends on json-exporter configuration.