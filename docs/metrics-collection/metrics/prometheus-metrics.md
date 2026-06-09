This document describes the metrics list and how to collect them from Prometheus.

# Metrics

Prometheus already exposes its metrics in Prometheus format and doesn't require to use of specific exporters.

| Name       | Metrics Port | Metrics Endpoint | Need Exporter? | Auth? | Is Exporter Third Party? |
| ---------- | ------------ | ---------------- | -------------- | ----- | ------------------------ |
| Prometheus | `9090`       | `/metrics`       | No             | No    | N/A                      |

## How to Collect

Metrics expose on port `9090` and endpoint `/metrics`. By default, Prometheus has no authentication.

Config `PodMonitor` for `prometheus-operator` to collect Prometheus metrics:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: PodMonitor
metadata:
  labels:
    app.kubernetes.io/component: monitoring
    app.kubernetes.io/managed-by: monitoring-operator
    app.kubernetes.io/name: prometheus-pod-monitor
    k8s-app: prometheus-pod-monitor
  name: prometheus-pod-monitor
spec:
  podMetricsEndpoints:
  - interval: 30s
    port: web
    scheme: http
  jobLabel: k8s-app
  namespaceSelector: {}
  selector:
    matchLabels:
      app.kubernetes.io/name: prometheus
```

To collect (or just to check) metrics manually you can use the following command:

```bash
curl -v -k -L http://<prometheus_ip_or_dns>:9090/metrics
```

or

```bash
wget -O - http://<prometheus_ip_or_dns>:9090/metrics
```

## Metrics list

```prometheus
# HELP go_gc_duration_seconds A summary of the pause duration of garbage collection cycles.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 0.000102787
go_gc_duration_seconds{quantile="0.25"} 0.000135784
go_gc_duration_seconds{quantile="0.5"} 0.000160478
go_gc_duration_seconds{quantile="0.75"} 0.00020981
go_gc_duration_seconds{quantile="1"} 0.005808876
go_gc_duration_seconds_sum 0.181144477
go_gc_duration_seconds_count 730
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 1450
# HELP go_info Information about the Go environment.
# TYPE go_info gauge
go_info{version="go1.19.3"} 1
# HELP go_memstats_alloc_bytes Number of bytes allocated and still in use.
# TYPE go_memstats_alloc_bytes gauge
go_memstats_alloc_bytes 1.457800744e+09
# HELP go_memstats_alloc_bytes_total Total number of bytes allocated, even if freed.
# TYPE go_memstats_alloc_bytes_total counter
go_memstats_alloc_bytes_total 7.0011457432e+11
# HELP go_memstats_buck_hash_sys_bytes Number of bytes used by the profiling bucket hash table.
# TYPE go_memstats_buck_hash_sys_bytes gauge
go_memstats_buck_hash_sys_bytes 5.531244e+06
# HELP go_memstats_frees_total Total number of frees.
# TYPE go_memstats_frees_total counter
go_memstats_frees_total 4.749523382e+09
# HELP go_memstats_gc_sys_bytes Number of bytes used for garbage collection system metadata.
# TYPE go_memstats_gc_sys_bytes gauge
go_memstats_gc_sys_bytes 1.20569256e+08
# HELP go_memstats_heap_alloc_bytes Number of heap bytes allocated and still in use.
# TYPE go_memstats_heap_alloc_bytes gauge
go_memstats_heap_alloc_bytes 1.457800744e+09
# HELP go_memstats_heap_idle_bytes Number of heap bytes waiting to be used.
# TYPE go_memstats_heap_idle_bytes gauge
go_memstats_heap_idle_bytes 1.218199552e+09
# HELP go_memstats_heap_inuse_bytes Number of heap bytes that are in use.
# TYPE go_memstats_heap_inuse_bytes gauge
go_memstats_heap_inuse_bytes 1.589559296e+09
# HELP go_memstats_heap_objects Number of allocated objects.
# TYPE go_memstats_heap_objects gauge
go_memstats_heap_objects 8.52274e+06
# HELP go_memstats_heap_released_bytes Number of heap bytes released to OS.
# TYPE go_memstats_heap_released_bytes gauge
go_memstats_heap_released_bytes 2.21749248e+08
# HELP go_memstats_heap_sys_bytes Number of heap bytes obtained from system.
# TYPE go_memstats_heap_sys_bytes gauge
go_memstats_heap_sys_bytes 2.807758848e+09
# HELP go_memstats_last_gc_time_seconds Number of seconds since 1970 of last garbage collection.
# TYPE go_memstats_last_gc_time_seconds gauge
go_memstats_last_gc_time_seconds 1.6726472220755763e+09
# HELP go_memstats_lookups_total Total number of pointer lookups.
# TYPE go_memstats_lookups_total counter
go_memstats_lookups_total 0
# HELP go_memstats_mallocs_total Total number of mallocs.
# TYPE go_memstats_mallocs_total counter
go_memstats_mallocs_total 4.758046122e+09
# HELP go_memstats_mcache_inuse_bytes Number of bytes in use by mcache structures.
# TYPE go_memstats_mcache_inuse_bytes gauge
go_memstats_mcache_inuse_bytes 2400
# HELP go_memstats_mcache_sys_bytes Number of bytes used for mcache structures obtained from system.
# TYPE go_memstats_mcache_sys_bytes gauge
go_memstats_mcache_sys_bytes 15600
# HELP go_memstats_mspan_inuse_bytes Number of bytes in use by mspan structures.
# TYPE go_memstats_mspan_inuse_bytes gauge
go_memstats_mspan_inuse_bytes 1.8980296e+07
# HELP go_memstats_mspan_sys_bytes Number of bytes used for mspan structures obtained from system.
# TYPE go_memstats_mspan_sys_bytes gauge
go_memstats_mspan_sys_bytes 3.388032e+07
# HELP go_memstats_next_gc_bytes Number of heap bytes when next garbage collection will take place.
# TYPE go_memstats_next_gc_bytes gauge
go_memstats_next_gc_bytes 2.3774022e+09
# HELP go_memstats_other_sys_bytes Number of bytes used for other system allocations.
# TYPE go_memstats_other_sys_bytes gauge
go_memstats_other_sys_bytes 4.773668e+06
# HELP go_memstats_stack_inuse_bytes Number of bytes in use by the stack allocator.
# TYPE go_memstats_stack_inuse_bytes gauge
go_memstats_stack_inuse_bytes 1.5007744e+07
# HELP go_memstats_stack_sys_bytes Number of bytes obtained from system for stack allocator.
# TYPE go_memstats_stack_sys_bytes gauge
go_memstats_stack_sys_bytes 1.5007744e+07
# HELP go_memstats_sys_bytes Number of bytes obtained from system.
# TYPE go_memstats_sys_bytes gauge
go_memstats_sys_bytes 2.98753668e+09
# HELP go_threads Number of OS threads created.
# TYPE go_threads gauge
go_threads 9
# HELP net_conntrack_dialer_conn_attempted_total Total number of connections attempted by the given dialer a given name.
# TYPE net_conntrack_dialer_conn_attempted_total counter
net_conntrack_dialer_conn_attempted_total{dialer_name="alertmanager"} 22
net_conntrack_dialer_conn_attempted_total{dialer_name="default"} 0
net_conntrack_dialer_conn_attempted_total{dialer_name="podMonitor/monitoring/monitoring-alertmanager-pod-monitor/0"} 1
net_conntrack_dialer_conn_attempted_total{dialer_name="podMonitor/monitoring/monitoring-grafana-operator-pod-monitor/0"} 1
net_conntrack_dialer_conn_attempted_total{dialer_name="podMonitor/monitoring/monitoring-grafana-pod-monitor/0"} 1
net_conntrack_dialer_conn_attempted_total{dialer_name="podMonitor/monitoring/monitoring-nginx-ingress-pod-monitor/0"} 6
net_conntrack_dialer_conn_attempted_total{dialer_name="podMonitor/monitoring/monitoring-prometheus-operator-pod-monitor/0"} 1
net_conntrack_dialer_conn_attempted_total{dialer_name="podMonitor/monitoring/monitoring-prometheus-pod-monitor/0"} 1
net_conntrack_dialer_conn_attempted_total{dialer_name="podMonitor/spark-operator-gcp/spark-operator-podmonitor/0"} 1
net_conntrack_dialer_conn_attempted_total{dialer_name="podMonitor/vault-service/vault-operator-pod-monitor/0"} 1
net_conntrack_dialer_conn_attempted_total{dialer_name="probe/monitoring-examples/blackbox-ingress-probe"} 0
net_conntrack_dialer_conn_attempted_total{dialer_name="probe/monitoring-examples/blackbox-static-probe"} 1523
net_conntrack_dialer_conn_attempted_total{dialer_name="remote_storage_write_client"} 36572
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/anko0714airflowbkp/airflow-anko0714airflowbkp/0"} 0
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/arango-db-a/arangodb-service-monitor/0"} 0
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/arangodb-graph-test/arangodb-service-monitor/0"} 0
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/arangodb/arangodb-service-monitor/0"} 0
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/consul-service/consul-service-monitor/0"} 3
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/core/cdc-control-monitor/0"} 16
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/core/cdc-stream-processor-monitor/0"} 16
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/core/cdc-streaming-service-monitor-jmx-exporter/0"} 762
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/core/cdc-test-api-monitor/0"} 16
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/core/cdc-test-db-monitor/0"} 0
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/core/cloud-streaming-service-monitor-jmx-exporter/0"} 762
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/core/mistral-service-monitor/0"} 1
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/default/arangodb-service-monitor/0"} 0
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/dev-helm-install/mistral-service-monitor/0"} 0
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/dev-mistral-demo/mistral-service-monitor/0"} 1
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/dev-mistral/mistral-service-monitor/0"} 1
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/domain-search/opensearch-service-monitor/0"} 1
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/elasticsearch-cluster/elasticsearch-service-monitor/0"} 1
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/elasticsearch-service/elasticsearch-service-monitor/0"} 1
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/envoy/monitoring-envoy/0"} 0
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/envoy/monitoring-keycloak/0"} 0
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/kafka-service/kafka-lag-exporter-service-monitor/0"} 761
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/kafka-service/kafka-service-monitor-jmx-exporter/0"} 2284
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/kafka-service/kafka-service-monitor/0"} 1
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/mistral/mistral-service-monitor/0"} 1
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/mongodb/mongo-backup-exporter-service-monitor/0"} 1523
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/mongodb/mongo-dbaas-exporter-service-monitor/0"} 1
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/mongodb/mongodb-prometheus-exporter-service-monitor/0"} 6095
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/monitoring-examples/go-example-service-service-monitor/0"} 0
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/monitoring-examples/go-example-service-service-monitor/1"} 0
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/monitoring-examples/java-example-service-monitor/0"} 0
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/monitoring-examples/monitoring-example-service-monitor/0"} 0
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/monitoring/monitoring-cert-exporter/0"} 1
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/monitoring/monitoring-core-dns-service-monitor/0"} 2
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/monitoring/monitoring-etcd-service-monitor/0"} 3
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/monitoring/monitoring-kube-apiserver-service-monitor/0"} 3
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/monitoring/monitoring-kube-state-metrics/0"} 1
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/monitoring/monitoring-kube-state-metrics/1"} 1
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/monitoring/monitoring-kube-state-metrics/2"} 1
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/monitoring/monitoring-kubelet-service-monitor/0"} 9
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/monitoring/monitoring-kubelet-service-monitor/1"} 9
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/monitoring/monitoring-node-exporter/0"} 9
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/opensearch-service/opensearch-service-monitor/0"} 1
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/pg-metr/postgres-service-monitor/0"} 1
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/postgres-upgrade/clickhouse-backup-daemon-service-monitor/0"} 0
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/postgres-upgrade/clickhouse-service-monitor/0"} 1
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/profiler-cassandra/cassandra-exporter-service-monitor/0"} 6855
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/profiler-test/esc-collector-service-monitor/0"} 16
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/profiler-test/esc-test-service-monitor/0"} 16
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/profiler-test/esc-ui-service-monitor/0"} 16
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/profiler/esc-collector-service-monitor/0"} 16
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/profiler/esc-test-service-monitor/0"} 16
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/profiler/esc-ui-service-monitor/0"} 16
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/rabbitmq/rabbitmq-per-object-service-monitor/0"} 6
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/rabbitmq/rabbitmq-service-monitor-old/0"} 6
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/rabbitmq/rabbitmq-service-monitor/0"} 6
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/redis/redis-prometheus-exporter-service-monitor/0"} 1
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/streaming-service/streaming-service-monitor-jmx-exporter/0"} 2284
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/streaming-service/streaming-service-monitor/0"} 1
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/vault-service/vault-service-monitor/0"} 2
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/zookeeper-service/zookeeper-service-monitor-jmx-exporter/0"} 2283
net_conntrack_dialer_conn_attempted_total{dialer_name="serviceMonitor/zookeeper-service/zookeeper-service-monitor/0"} 1
# HELP net_conntrack_dialer_conn_closed_total Total number of connections closed which originated from the dialer of a given name.
# TYPE net_conntrack_dialer_conn_closed_total counter
net_conntrack_dialer_conn_closed_total{dialer_name="alertmanager"} 20
net_conntrack_dialer_conn_closed_total{dialer_name="default"} 0
net_conntrack_dialer_conn_closed_total{dialer_name="podMonitor/monitoring/monitoring-alertmanager-pod-monitor/0"} 0
net_conntrack_dialer_conn_closed_total{dialer_name="podMonitor/monitoring/monitoring-grafana-operator-pod-monitor/0"} 0
net_conntrack_dialer_conn_closed_total{dialer_name="podMonitor/monitoring/monitoring-grafana-pod-monitor/0"} 0
net_conntrack_dialer_conn_closed_total{dialer_name="podMonitor/monitoring/monitoring-nginx-ingress-pod-monitor/0"} 0
net_conntrack_dialer_conn_closed_total{dialer_name="podMonitor/monitoring/monitoring-prometheus-operator-pod-monitor/0"} 0
net_conntrack_dialer_conn_closed_total{dialer_name="podMonitor/monitoring/monitoring-prometheus-pod-monitor/0"} 0
net_conntrack_dialer_conn_closed_total{dialer_name="podMonitor/spark-operator-gcp/spark-operator-podmonitor/0"} 0
net_conntrack_dialer_conn_closed_total{dialer_name="podMonitor/vault-service/vault-operator-pod-monitor/0"} 0
net_conntrack_dialer_conn_closed_total{dialer_name="probe/monitoring-examples/blackbox-ingress-probe"} 0
net_conntrack_dialer_conn_closed_total{dialer_name="probe/monitoring-examples/blackbox-static-probe"} 0
net_conntrack_dialer_conn_closed_total{dialer_name="remote_storage_write_client"} 18302
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/anko0714airflowbkp/airflow-anko0714airflowbkp/0"} 0
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/arango-db-a/arangodb-service-monitor/0"} 0
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/arangodb-graph-test/arangodb-service-monitor/0"} 0
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/arangodb/arangodb-service-monitor/0"} 0
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/consul-service/consul-service-monitor/0"} 0
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/core/cdc-control-monitor/0"} 15
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/core/cdc-stream-processor-monitor/0"} 15
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/core/cdc-streaming-service-monitor-jmx-exporter/0"} 761
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/core/cdc-test-api-monitor/0"} 15
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/core/cdc-test-db-monitor/0"} 0
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/core/cloud-streaming-service-monitor-jmx-exporter/0"} 761
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/core/mistral-service-monitor/0"} 0
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/default/arangodb-service-monitor/0"} 0
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/dev-helm-install/mistral-service-monitor/0"} 0
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/dev-mistral-demo/mistral-service-monitor/0"} 0
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/dev-mistral/mistral-service-monitor/0"} 0
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/domain-search/opensearch-service-monitor/0"} 0
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/elasticsearch-cluster/elasticsearch-service-monitor/0"} 0
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/elasticsearch-service/elasticsearch-service-monitor/0"} 0
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/envoy/monitoring-envoy/0"} 0
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/envoy/monitoring-keycloak/0"} 0
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/kafka-service/kafka-lag-exporter-service-monitor/0"} 761
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/kafka-service/kafka-service-monitor-jmx-exporter/0"} 2283
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/kafka-service/kafka-service-monitor/0"} 0
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/mistral/mistral-service-monitor/0"} 0
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/mongodb/mongo-backup-exporter-service-monitor/0"} 1523
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/mongodb/mongo-dbaas-exporter-service-monitor/0"} 0
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/mongodb/mongodb-prometheus-exporter-service-monitor/0"} 1522
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/monitoring-examples/go-example-service-service-monitor/0"} 0
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/monitoring-examples/go-example-service-service-monitor/1"} 0
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/monitoring-examples/java-example-service-monitor/0"} 0
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/monitoring-examples/monitoring-example-service-monitor/0"} 0
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/monitoring/monitoring-cert-exporter/0"} 0
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/monitoring/monitoring-core-dns-service-monitor/0"} 0
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/monitoring/monitoring-etcd-service-monitor/0"} 0
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/monitoring/monitoring-kube-apiserver-service-monitor/0"} 0
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/monitoring/monitoring-kube-state-metrics/0"} 0
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/monitoring/monitoring-kube-state-metrics/1"} 0
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/monitoring/monitoring-kube-state-metrics/2"} 0
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/monitoring/monitoring-kubelet-service-monitor/0"} 0
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/monitoring/monitoring-kubelet-service-monitor/1"} 0
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/monitoring/monitoring-node-exporter/0"} 0
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/opensearch-service/opensearch-service-monitor/0"} 0
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/pg-metr/postgres-service-monitor/0"} 0
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/postgres-upgrade/clickhouse-backup-daemon-service-monitor/0"} 0
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/postgres-upgrade/clickhouse-service-monitor/0"} 0
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/profiler-cassandra/cassandra-exporter-service-monitor/0"} 6854
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/profiler-test/esc-collector-service-monitor/0"} 15
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/profiler-test/esc-test-service-monitor/0"} 15
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/profiler-test/esc-ui-service-monitor/0"} 15
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/profiler/esc-collector-service-monitor/0"} 15
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/profiler/esc-test-service-monitor/0"} 15
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/profiler/esc-ui-service-monitor/0"} 15
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/rabbitmq/rabbitmq-per-object-service-monitor/0"} 3
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/rabbitmq/rabbitmq-service-monitor-old/0"} 3
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/rabbitmq/rabbitmq-service-monitor/0"} 3
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/redis/redis-prometheus-exporter-service-monitor/0"} 0
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/streaming-service/streaming-service-monitor-jmx-exporter/0"} 2283
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/streaming-service/streaming-service-monitor/0"} 0
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/vault-service/vault-service-monitor/0"} 0
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/zookeeper-service/zookeeper-service-monitor-jmx-exporter/0"} 2283
net_conntrack_dialer_conn_closed_total{dialer_name="serviceMonitor/zookeeper-service/zookeeper-service-monitor/0"} 0
# HELP net_conntrack_dialer_conn_established_total Total number of connections successfully established by the given dialer a given name.
# TYPE net_conntrack_dialer_conn_established_total counter
net_conntrack_dialer_conn_established_total{dialer_name="alertmanager"} 22
net_conntrack_dialer_conn_established_total{dialer_name="default"} 0
net_conntrack_dialer_conn_established_total{dialer_name="podMonitor/monitoring/monitoring-alertmanager-pod-monitor/0"} 1
net_conntrack_dialer_conn_established_total{dialer_name="podMonitor/monitoring/monitoring-grafana-operator-pod-monitor/0"} 1
net_conntrack_dialer_conn_established_total{dialer_name="podMonitor/monitoring/monitoring-grafana-pod-monitor/0"} 1
net_conntrack_dialer_conn_established_total{dialer_name="podMonitor/monitoring/monitoring-nginx-ingress-pod-monitor/0"} 6
net_conntrack_dialer_conn_established_total{dialer_name="podMonitor/monitoring/monitoring-prometheus-operator-pod-monitor/0"} 1
net_conntrack_dialer_conn_established_total{dialer_name="podMonitor/monitoring/monitoring-prometheus-pod-monitor/0"} 1
net_conntrack_dialer_conn_established_total{dialer_name="podMonitor/spark-operator-gcp/spark-operator-podmonitor/0"} 1
net_conntrack_dialer_conn_established_total{dialer_name="podMonitor/vault-service/vault-operator-pod-monitor/0"} 1
net_conntrack_dialer_conn_established_total{dialer_name="probe/monitoring-examples/blackbox-ingress-probe"} 0
net_conntrack_dialer_conn_established_total{dialer_name="probe/monitoring-examples/blackbox-static-probe"} 0
net_conntrack_dialer_conn_established_total{dialer_name="remote_storage_write_client"} 18303
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/anko0714airflowbkp/airflow-anko0714airflowbkp/0"} 0
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/arango-db-a/arangodb-service-monitor/0"} 0
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/arangodb-graph-test/arangodb-service-monitor/0"} 0
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/arangodb/arangodb-service-monitor/0"} 0
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/consul-service/consul-service-monitor/0"} 3
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/core/cdc-control-monitor/0"} 16
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/core/cdc-stream-processor-monitor/0"} 16
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/core/cdc-streaming-service-monitor-jmx-exporter/0"} 762
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/core/cdc-test-api-monitor/0"} 16
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/core/cdc-test-db-monitor/0"} 0
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/core/cloud-streaming-service-monitor-jmx-exporter/0"} 762
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/core/mistral-service-monitor/0"} 1
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/default/arangodb-service-monitor/0"} 0
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/dev-helm-install/mistral-service-monitor/0"} 0
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/dev-mistral-demo/mistral-service-monitor/0"} 1
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/dev-mistral/mistral-service-monitor/0"} 1
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/domain-search/opensearch-service-monitor/0"} 1
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/elasticsearch-cluster/elasticsearch-service-monitor/0"} 1
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/elasticsearch-service/elasticsearch-service-monitor/0"} 1
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/envoy/monitoring-envoy/0"} 0
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/envoy/monitoring-keycloak/0"} 0
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/kafka-service/kafka-lag-exporter-service-monitor/0"} 761
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/kafka-service/kafka-service-monitor-jmx-exporter/0"} 2284
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/kafka-service/kafka-service-monitor/0"} 1
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/mistral/mistral-service-monitor/0"} 1
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/mongodb/mongo-backup-exporter-service-monitor/0"} 1523
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/mongodb/mongo-dbaas-exporter-service-monitor/0"} 1
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/mongodb/mongodb-prometheus-exporter-service-monitor/0"} 1524
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/monitoring-examples/go-example-service-service-monitor/0"} 0
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/monitoring-examples/go-example-service-service-monitor/1"} 0
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/monitoring-examples/java-example-service-monitor/0"} 0
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/monitoring-examples/monitoring-example-service-monitor/0"} 0
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/monitoring/monitoring-cert-exporter/0"} 1
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/monitoring/monitoring-core-dns-service-monitor/0"} 2
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/monitoring/monitoring-etcd-service-monitor/0"} 3
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/monitoring/monitoring-kube-apiserver-service-monitor/0"} 3
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/monitoring/monitoring-kube-state-metrics/0"} 1
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/monitoring/monitoring-kube-state-metrics/1"} 1
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/monitoring/monitoring-kube-state-metrics/2"} 1
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/monitoring/monitoring-kubelet-service-monitor/0"} 9
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/monitoring/monitoring-kubelet-service-monitor/1"} 9
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/monitoring/monitoring-node-exporter/0"} 9
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/opensearch-service/opensearch-service-monitor/0"} 1
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/pg-metr/postgres-service-monitor/0"} 1
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/postgres-upgrade/clickhouse-backup-daemon-service-monitor/0"} 0
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/postgres-upgrade/clickhouse-service-monitor/0"} 1
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/profiler-cassandra/cassandra-exporter-service-monitor/0"} 6854
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/profiler-test/esc-collector-service-monitor/0"} 16
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/profiler-test/esc-test-service-monitor/0"} 16
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/profiler-test/esc-ui-service-monitor/0"} 16
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/profiler/esc-collector-service-monitor/0"} 16
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/profiler/esc-test-service-monitor/0"} 16
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/profiler/esc-ui-service-monitor/0"} 16
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/rabbitmq/rabbitmq-per-object-service-monitor/0"} 6
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/rabbitmq/rabbitmq-service-monitor-old/0"} 6
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/rabbitmq/rabbitmq-service-monitor/0"} 6
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/redis/redis-prometheus-exporter-service-monitor/0"} 1
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/streaming-service/streaming-service-monitor-jmx-exporter/0"} 2284
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/streaming-service/streaming-service-monitor/0"} 1
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/vault-service/vault-service-monitor/0"} 2
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/zookeeper-service/zookeeper-service-monitor-jmx-exporter/0"} 2283
net_conntrack_dialer_conn_established_total{dialer_name="serviceMonitor/zookeeper-service/zookeeper-service-monitor/0"} 1
# HELP net_conntrack_dialer_conn_failed_total Total number of connections failed to dial by the dialer a given name.
# TYPE net_conntrack_dialer_conn_failed_total counter
net_conntrack_dialer_conn_failed_total{dialer_name="alertmanager",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="alertmanager",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="alertmanager",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="alertmanager",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="default",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="default",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="default",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="default",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="podMonitor/monitoring/monitoring-alertmanager-pod-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="podMonitor/monitoring/monitoring-alertmanager-pod-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="podMonitor/monitoring/monitoring-alertmanager-pod-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="podMonitor/monitoring/monitoring-alertmanager-pod-monitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="podMonitor/monitoring/monitoring-grafana-operator-pod-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="podMonitor/monitoring/monitoring-grafana-operator-pod-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="podMonitor/monitoring/monitoring-grafana-operator-pod-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="podMonitor/monitoring/monitoring-grafana-operator-pod-monitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="podMonitor/monitoring/monitoring-grafana-pod-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="podMonitor/monitoring/monitoring-grafana-pod-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="podMonitor/monitoring/monitoring-grafana-pod-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="podMonitor/monitoring/monitoring-grafana-pod-monitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="podMonitor/monitoring/monitoring-nginx-ingress-pod-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="podMonitor/monitoring/monitoring-nginx-ingress-pod-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="podMonitor/monitoring/monitoring-nginx-ingress-pod-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="podMonitor/monitoring/monitoring-nginx-ingress-pod-monitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="podMonitor/monitoring/monitoring-prometheus-operator-pod-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="podMonitor/monitoring/monitoring-prometheus-operator-pod-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="podMonitor/monitoring/monitoring-prometheus-operator-pod-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="podMonitor/monitoring/monitoring-prometheus-operator-pod-monitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="podMonitor/monitoring/monitoring-prometheus-pod-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="podMonitor/monitoring/monitoring-prometheus-pod-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="podMonitor/monitoring/monitoring-prometheus-pod-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="podMonitor/monitoring/monitoring-prometheus-pod-monitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="podMonitor/spark-operator-gcp/spark-operator-podmonitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="podMonitor/spark-operator-gcp/spark-operator-podmonitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="podMonitor/spark-operator-gcp/spark-operator-podmonitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="podMonitor/spark-operator-gcp/spark-operator-podmonitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="podMonitor/vault-service/vault-operator-pod-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="podMonitor/vault-service/vault-operator-pod-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="podMonitor/vault-service/vault-operator-pod-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="podMonitor/vault-service/vault-operator-pod-monitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="probe/monitoring-examples/blackbox-ingress-probe",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="probe/monitoring-examples/blackbox-ingress-probe",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="probe/monitoring-examples/blackbox-ingress-probe",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="probe/monitoring-examples/blackbox-ingress-probe",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="probe/monitoring-examples/blackbox-static-probe",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="probe/monitoring-examples/blackbox-static-probe",reason="resolution"} 1523
net_conntrack_dialer_conn_failed_total{dialer_name="probe/monitoring-examples/blackbox-static-probe",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="probe/monitoring-examples/blackbox-static-probe",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="remote_storage_write_client",reason="refused"} 18269
net_conntrack_dialer_conn_failed_total{dialer_name="remote_storage_write_client",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="remote_storage_write_client",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="remote_storage_write_client",reason="unknown"} 18269
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/anko0714airflowbkp/airflow-anko0714airflowbkp/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/anko0714airflowbkp/airflow-anko0714airflowbkp/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/anko0714airflowbkp/airflow-anko0714airflowbkp/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/anko0714airflowbkp/airflow-anko0714airflowbkp/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/arango-db-a/arangodb-service-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/arango-db-a/arangodb-service-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/arango-db-a/arangodb-service-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/arango-db-a/arangodb-service-monitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/arangodb-graph-test/arangodb-service-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/arangodb-graph-test/arangodb-service-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/arangodb-graph-test/arangodb-service-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/arangodb-graph-test/arangodb-service-monitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/arangodb/arangodb-service-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/arangodb/arangodb-service-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/arangodb/arangodb-service-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/arangodb/arangodb-service-monitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/consul-service/consul-service-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/consul-service/consul-service-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/consul-service/consul-service-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/consul-service/consul-service-monitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/core/cdc-control-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/core/cdc-control-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/core/cdc-control-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/core/cdc-control-monitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/core/cdc-stream-processor-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/core/cdc-stream-processor-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/core/cdc-stream-processor-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/core/cdc-stream-processor-monitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/core/cdc-streaming-service-monitor-jmx-exporter/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/core/cdc-streaming-service-monitor-jmx-exporter/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/core/cdc-streaming-service-monitor-jmx-exporter/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/core/cdc-streaming-service-monitor-jmx-exporter/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/core/cdc-test-api-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/core/cdc-test-api-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/core/cdc-test-api-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/core/cdc-test-api-monitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/core/cdc-test-db-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/core/cdc-test-db-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/core/cdc-test-db-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/core/cdc-test-db-monitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/core/cloud-streaming-service-monitor-jmx-exporter/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/core/cloud-streaming-service-monitor-jmx-exporter/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/core/cloud-streaming-service-monitor-jmx-exporter/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/core/cloud-streaming-service-monitor-jmx-exporter/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/core/mistral-service-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/core/mistral-service-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/core/mistral-service-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/core/mistral-service-monitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/default/arangodb-service-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/default/arangodb-service-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/default/arangodb-service-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/default/arangodb-service-monitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/dev-helm-install/mistral-service-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/dev-helm-install/mistral-service-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/dev-helm-install/mistral-service-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/dev-helm-install/mistral-service-monitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/dev-mistral-demo/mistral-service-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/dev-mistral-demo/mistral-service-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/dev-mistral-demo/mistral-service-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/dev-mistral-demo/mistral-service-monitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/dev-mistral/mistral-service-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/dev-mistral/mistral-service-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/dev-mistral/mistral-service-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/dev-mistral/mistral-service-monitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/domain-search/opensearch-service-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/domain-search/opensearch-service-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/domain-search/opensearch-service-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/domain-search/opensearch-service-monitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/elasticsearch-cluster/elasticsearch-service-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/elasticsearch-cluster/elasticsearch-service-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/elasticsearch-cluster/elasticsearch-service-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/elasticsearch-cluster/elasticsearch-service-monitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/elasticsearch-service/elasticsearch-service-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/elasticsearch-service/elasticsearch-service-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/elasticsearch-service/elasticsearch-service-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/elasticsearch-service/elasticsearch-service-monitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/envoy/monitoring-envoy/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/envoy/monitoring-envoy/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/envoy/monitoring-envoy/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/envoy/monitoring-envoy/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/envoy/monitoring-keycloak/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/envoy/monitoring-keycloak/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/envoy/monitoring-keycloak/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/envoy/monitoring-keycloak/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/kafka-service/kafka-lag-exporter-service-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/kafka-service/kafka-lag-exporter-service-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/kafka-service/kafka-lag-exporter-service-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/kafka-service/kafka-lag-exporter-service-monitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/kafka-service/kafka-service-monitor-jmx-exporter/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/kafka-service/kafka-service-monitor-jmx-exporter/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/kafka-service/kafka-service-monitor-jmx-exporter/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/kafka-service/kafka-service-monitor-jmx-exporter/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/kafka-service/kafka-service-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/kafka-service/kafka-service-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/kafka-service/kafka-service-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/kafka-service/kafka-service-monitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/mistral/mistral-service-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/mistral/mistral-service-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/mistral/mistral-service-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/mistral/mistral-service-monitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/mongodb/mongo-backup-exporter-service-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/mongodb/mongo-backup-exporter-service-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/mongodb/mongo-backup-exporter-service-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/mongodb/mongo-backup-exporter-service-monitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/mongodb/mongo-dbaas-exporter-service-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/mongodb/mongo-dbaas-exporter-service-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/mongodb/mongo-dbaas-exporter-service-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/mongodb/mongo-dbaas-exporter-service-monitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/mongodb/mongodb-prometheus-exporter-service-monitor/0",reason="refused"} 4570
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/mongodb/mongodb-prometheus-exporter-service-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/mongodb/mongodb-prometheus-exporter-service-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/mongodb/mongodb-prometheus-exporter-service-monitor/0",reason="unknown"} 4571
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring-examples/go-example-service-service-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring-examples/go-example-service-service-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring-examples/go-example-service-service-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring-examples/go-example-service-service-monitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring-examples/go-example-service-service-monitor/1",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring-examples/go-example-service-service-monitor/1",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring-examples/go-example-service-service-monitor/1",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring-examples/go-example-service-service-monitor/1",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring-examples/java-example-service-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring-examples/java-example-service-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring-examples/java-example-service-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring-examples/java-example-service-monitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring-examples/monitoring-example-service-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring-examples/monitoring-example-service-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring-examples/monitoring-example-service-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring-examples/monitoring-example-service-monitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring/monitoring-cert-exporter/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring/monitoring-cert-exporter/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring/monitoring-cert-exporter/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring/monitoring-cert-exporter/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring/monitoring-core-dns-service-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring/monitoring-core-dns-service-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring/monitoring-core-dns-service-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring/monitoring-core-dns-service-monitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring/monitoring-etcd-service-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring/monitoring-etcd-service-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring/monitoring-etcd-service-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring/monitoring-etcd-service-monitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring/monitoring-kube-apiserver-service-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring/monitoring-kube-apiserver-service-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring/monitoring-kube-apiserver-service-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring/monitoring-kube-apiserver-service-monitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring/monitoring-kube-state-metrics/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring/monitoring-kube-state-metrics/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring/monitoring-kube-state-metrics/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring/monitoring-kube-state-metrics/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring/monitoring-kube-state-metrics/1",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring/monitoring-kube-state-metrics/1",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring/monitoring-kube-state-metrics/1",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring/monitoring-kube-state-metrics/1",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring/monitoring-kube-state-metrics/2",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring/monitoring-kube-state-metrics/2",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring/monitoring-kube-state-metrics/2",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring/monitoring-kube-state-metrics/2",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring/monitoring-kubelet-service-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring/monitoring-kubelet-service-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring/monitoring-kubelet-service-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring/monitoring-kubelet-service-monitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring/monitoring-kubelet-service-monitor/1",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring/monitoring-kubelet-service-monitor/1",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring/monitoring-kubelet-service-monitor/1",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring/monitoring-kubelet-service-monitor/1",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring/monitoring-node-exporter/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring/monitoring-node-exporter/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring/monitoring-node-exporter/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/monitoring/monitoring-node-exporter/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/opensearch-service/opensearch-service-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/opensearch-service/opensearch-service-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/opensearch-service/opensearch-service-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/opensearch-service/opensearch-service-monitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/pg-metr/postgres-service-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/pg-metr/postgres-service-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/pg-metr/postgres-service-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/pg-metr/postgres-service-monitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/postgres-upgrade/clickhouse-backup-daemon-service-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/postgres-upgrade/clickhouse-backup-daemon-service-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/postgres-upgrade/clickhouse-backup-daemon-service-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/postgres-upgrade/clickhouse-backup-daemon-service-monitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/postgres-upgrade/clickhouse-service-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/postgres-upgrade/clickhouse-service-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/postgres-upgrade/clickhouse-service-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/postgres-upgrade/clickhouse-service-monitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/profiler-cassandra/cassandra-exporter-service-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/profiler-cassandra/cassandra-exporter-service-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/profiler-cassandra/cassandra-exporter-service-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/profiler-cassandra/cassandra-exporter-service-monitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/profiler-test/esc-collector-service-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/profiler-test/esc-collector-service-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/profiler-test/esc-collector-service-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/profiler-test/esc-collector-service-monitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/profiler-test/esc-test-service-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/profiler-test/esc-test-service-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/profiler-test/esc-test-service-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/profiler-test/esc-test-service-monitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/profiler-test/esc-ui-service-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/profiler-test/esc-ui-service-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/profiler-test/esc-ui-service-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/profiler-test/esc-ui-service-monitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/profiler/esc-collector-service-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/profiler/esc-collector-service-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/profiler/esc-collector-service-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/profiler/esc-collector-service-monitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/profiler/esc-test-service-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/profiler/esc-test-service-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/profiler/esc-test-service-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/profiler/esc-test-service-monitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/profiler/esc-ui-service-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/profiler/esc-ui-service-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/profiler/esc-ui-service-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/profiler/esc-ui-service-monitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/rabbitmq/rabbitmq-per-object-service-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/rabbitmq/rabbitmq-per-object-service-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/rabbitmq/rabbitmq-per-object-service-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/rabbitmq/rabbitmq-per-object-service-monitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/rabbitmq/rabbitmq-service-monitor-old/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/rabbitmq/rabbitmq-service-monitor-old/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/rabbitmq/rabbitmq-service-monitor-old/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/rabbitmq/rabbitmq-service-monitor-old/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/rabbitmq/rabbitmq-service-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/rabbitmq/rabbitmq-service-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/rabbitmq/rabbitmq-service-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/rabbitmq/rabbitmq-service-monitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/redis/redis-prometheus-exporter-service-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/redis/redis-prometheus-exporter-service-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/redis/redis-prometheus-exporter-service-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/redis/redis-prometheus-exporter-service-monitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/streaming-service/streaming-service-monitor-jmx-exporter/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/streaming-service/streaming-service-monitor-jmx-exporter/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/streaming-service/streaming-service-monitor-jmx-exporter/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/streaming-service/streaming-service-monitor-jmx-exporter/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/streaming-service/streaming-service-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/streaming-service/streaming-service-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/streaming-service/streaming-service-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/streaming-service/streaming-service-monitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/vault-service/vault-service-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/vault-service/vault-service-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/vault-service/vault-service-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/vault-service/vault-service-monitor/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/zookeeper-service/zookeeper-service-monitor-jmx-exporter/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/zookeeper-service/zookeeper-service-monitor-jmx-exporter/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/zookeeper-service/zookeeper-service-monitor-jmx-exporter/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/zookeeper-service/zookeeper-service-monitor-jmx-exporter/0",reason="unknown"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/zookeeper-service/zookeeper-service-monitor/0",reason="refused"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/zookeeper-service/zookeeper-service-monitor/0",reason="resolution"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/zookeeper-service/zookeeper-service-monitor/0",reason="timeout"} 0
net_conntrack_dialer_conn_failed_total{dialer_name="serviceMonitor/zookeeper-service/zookeeper-service-monitor/0",reason="unknown"} 0
# HELP net_conntrack_listener_conn_accepted_total Total number of connections opened to the listener of a given name.
# TYPE net_conntrack_listener_conn_accepted_total counter
net_conntrack_listener_conn_accepted_total{listener_name="http"} 18282
# HELP net_conntrack_listener_conn_closed_total Total number of connections closed that were made to the listener of a given name.
# TYPE net_conntrack_listener_conn_closed_total counter
net_conntrack_listener_conn_closed_total{listener_name="http"} 18280
# HELP process_cpu_seconds_total Total user and system CPU time spent in seconds.
# TYPE process_cpu_seconds_total counter
process_cpu_seconds_total 10828.95
# HELP process_max_fds Maximum number of open file descriptors.
# TYPE process_max_fds gauge
process_max_fds 1.048576e+06
# HELP process_open_fds Number of open file descriptors.
# TYPE process_open_fds gauge
process_open_fds 131
# HELP process_resident_memory_bytes Resident memory size in bytes.
# TYPE process_resident_memory_bytes gauge
process_resident_memory_bytes 2.910109696e+09
# HELP process_start_time_seconds Start time of the process since unix epoch in seconds.
# TYPE process_start_time_seconds gauge
process_start_time_seconds 1.67260154866e+09
# HELP process_virtual_memory_bytes Virtual memory size in bytes.
# TYPE process_virtual_memory_bytes gauge
process_virtual_memory_bytes 4.921503744e+09
# HELP process_virtual_memory_max_bytes Maximum amount of virtual memory available in bytes.
# TYPE process_virtual_memory_max_bytes gauge
process_virtual_memory_max_bytes 1.8446744073709552e+19
# HELP prometheus_api_remote_read_queries The current number of remote read queries being executed or waiting.
# TYPE prometheus_api_remote_read_queries gauge
prometheus_api_remote_read_queries 0
# HELP prometheus_build_info A metric with a constant '1' value labeled by version, revision, branch, and goversion from which prometheus was built.
# TYPE prometheus_build_info gauge
prometheus_build_info{branch="HEAD",goversion="go1.19.3",revision="44af4716c86138869aa621737139e6dacf0e2550",version="2.40.5"} 1
# HELP prometheus_config_last_reload_success_timestamp_seconds Timestamp of the last successful configuration reload.
# TYPE prometheus_config_last_reload_success_timestamp_seconds gauge
prometheus_config_last_reload_success_timestamp_seconds 1.6726015500208852e+09
# HELP prometheus_config_last_reload_successful Whether the last configuration reload attempt was successful.
# TYPE prometheus_config_last_reload_successful gauge
prometheus_config_last_reload_successful 1
# HELP prometheus_engine_queries The current number of queries being executed or waiting.
# TYPE prometheus_engine_queries gauge
prometheus_engine_queries 0
# HELP prometheus_engine_queries_concurrent_max The max number of concurrent queries.
# TYPE prometheus_engine_queries_concurrent_max gauge
prometheus_engine_queries_concurrent_max 20
# HELP prometheus_engine_query_duration_seconds Query timings
# TYPE prometheus_engine_query_duration_seconds summary
prometheus_engine_query_duration_seconds{slice="inner_eval",quantile="0.5"} 5.3591e-05
prometheus_engine_query_duration_seconds{slice="inner_eval",quantile="0.9"} 0.00042511
prometheus_engine_query_duration_seconds{slice="inner_eval",quantile="0.99"} 0.021718858
prometheus_engine_query_duration_seconds_sum{slice="inner_eval"} 1500.2650119610184
prometheus_engine_query_duration_seconds_count{slice="inner_eval"} 807192
prometheus_engine_query_duration_seconds{slice="prepare_time",quantile="0.5"} 0.000100928
prometheus_engine_query_duration_seconds{slice="prepare_time",quantile="0.9"} 0.000380495
prometheus_engine_query_duration_seconds{slice="prepare_time",quantile="0.99"} 0.004744683
prometheus_engine_query_duration_seconds_sum{slice="prepare_time"} 490.4031671379878
prometheus_engine_query_duration_seconds_count{slice="prepare_time"} 807192
prometheus_engine_query_duration_seconds{slice="queue_time",quantile="0.5"} 8.832e-06
prometheus_engine_query_duration_seconds{slice="queue_time",quantile="0.9"} 1.9715e-05
prometheus_engine_query_duration_seconds{slice="queue_time",quantile="0.99"} 4.485e-05
prometheus_engine_query_duration_seconds_sum{slice="queue_time"} 8.823439537999821
prometheus_engine_query_duration_seconds_count{slice="queue_time"} 807192
prometheus_engine_query_duration_seconds{slice="result_sort",quantile="0.5"} NaN
prometheus_engine_query_duration_seconds{slice="result_sort",quantile="0.9"} NaN
prometheus_engine_query_duration_seconds{slice="result_sort",quantile="0.99"} NaN
prometheus_engine_query_duration_seconds_sum{slice="result_sort"} 0
prometheus_engine_query_duration_seconds_count{slice="result_sort"} 0
# HELP prometheus_engine_query_log_enabled State of the query log.
# TYPE prometheus_engine_query_log_enabled gauge
prometheus_engine_query_log_enabled 0
# HELP prometheus_engine_query_log_failures_total The number of query log failures.
# TYPE prometheus_engine_query_log_failures_total counter
prometheus_engine_query_log_failures_total 0
# HELP prometheus_http_request_duration_seconds Histogram of latencies for HTTP requests.
# TYPE prometheus_http_request_duration_seconds histogram
prometheus_http_request_duration_seconds_bucket{handler="/-/healthy",le="0.1"} 9135
prometheus_http_request_duration_seconds_bucket{handler="/-/healthy",le="0.2"} 9135
prometheus_http_request_duration_seconds_bucket{handler="/-/healthy",le="0.4"} 9135
prometheus_http_request_duration_seconds_bucket{handler="/-/healthy",le="1"} 9135
prometheus_http_request_duration_seconds_bucket{handler="/-/healthy",le="3"} 9135
prometheus_http_request_duration_seconds_bucket{handler="/-/healthy",le="8"} 9135
prometheus_http_request_duration_seconds_bucket{handler="/-/healthy",le="20"} 9135
prometheus_http_request_duration_seconds_bucket{handler="/-/healthy",le="60"} 9135
prometheus_http_request_duration_seconds_bucket{handler="/-/healthy",le="120"} 9135
prometheus_http_request_duration_seconds_bucket{handler="/-/healthy",le="+Inf"} 9135
prometheus_http_request_duration_seconds_sum{handler="/-/healthy"} 0.16481428200000045
prometheus_http_request_duration_seconds_count{handler="/-/healthy"} 9135
prometheus_http_request_duration_seconds_bucket{handler="/-/ready",le="0.1"} 9139
prometheus_http_request_duration_seconds_bucket{handler="/-/ready",le="0.2"} 9139
prometheus_http_request_duration_seconds_bucket{handler="/-/ready",le="0.4"} 9139
prometheus_http_request_duration_seconds_bucket{handler="/-/ready",le="1"} 9139
prometheus_http_request_duration_seconds_bucket{handler="/-/ready",le="3"} 9139
prometheus_http_request_duration_seconds_bucket{handler="/-/ready",le="8"} 9139
prometheus_http_request_duration_seconds_bucket{handler="/-/ready",le="20"} 9139
prometheus_http_request_duration_seconds_bucket{handler="/-/ready",le="60"} 9139
prometheus_http_request_duration_seconds_bucket{handler="/-/ready",le="120"} 9139
prometheus_http_request_duration_seconds_bucket{handler="/-/ready",le="+Inf"} 9139
prometheus_http_request_duration_seconds_sum{handler="/-/ready"} 0.15617847400000054
prometheus_http_request_duration_seconds_count{handler="/-/ready"} 9139
prometheus_http_request_duration_seconds_bucket{handler="/-/reload",le="0.1"} 0
prometheus_http_request_duration_seconds_bucket{handler="/-/reload",le="0.2"} 0
prometheus_http_request_duration_seconds_bucket{handler="/-/reload",le="0.4"} 1
prometheus_http_request_duration_seconds_bucket{handler="/-/reload",le="1"} 1
prometheus_http_request_duration_seconds_bucket{handler="/-/reload",le="3"} 1
prometheus_http_request_duration_seconds_bucket{handler="/-/reload",le="8"} 1
prometheus_http_request_duration_seconds_bucket{handler="/-/reload",le="20"} 1
prometheus_http_request_duration_seconds_bucket{handler="/-/reload",le="60"} 1
prometheus_http_request_duration_seconds_bucket{handler="/-/reload",le="120"} 1
prometheus_http_request_duration_seconds_bucket{handler="/-/reload",le="+Inf"} 1
prometheus_http_request_duration_seconds_sum{handler="/-/reload"} 0.366002603
prometheus_http_request_duration_seconds_count{handler="/-/reload"} 1
prometheus_http_request_duration_seconds_bucket{handler="/alerts",le="0.1"} 3
prometheus_http_request_duration_seconds_bucket{handler="/alerts",le="0.2"} 3
prometheus_http_request_duration_seconds_bucket{handler="/alerts",le="0.4"} 3
prometheus_http_request_duration_seconds_bucket{handler="/alerts",le="1"} 3
prometheus_http_request_duration_seconds_bucket{handler="/alerts",le="3"} 3
prometheus_http_request_duration_seconds_bucket{handler="/alerts",le="8"} 3
prometheus_http_request_duration_seconds_bucket{handler="/alerts",le="20"} 3
prometheus_http_request_duration_seconds_bucket{handler="/alerts",le="60"} 3
prometheus_http_request_duration_seconds_bucket{handler="/alerts",le="120"} 3
prometheus_http_request_duration_seconds_bucket{handler="/alerts",le="+Inf"} 3
prometheus_http_request_duration_seconds_sum{handler="/alerts"} 0.012007298
prometheus_http_request_duration_seconds_count{handler="/alerts"} 3
prometheus_http_request_duration_seconds_bucket{handler="/api/v1/label/:name/values",le="0.1"} 2
prometheus_http_request_duration_seconds_bucket{handler="/api/v1/label/:name/values",le="0.2"} 2
prometheus_http_request_duration_seconds_bucket{handler="/api/v1/label/:name/values",le="0.4"} 2
prometheus_http_request_duration_seconds_bucket{handler="/api/v1/label/:name/values",le="1"} 2
prometheus_http_request_duration_seconds_bucket{handler="/api/v1/label/:name/values",le="3"} 2
prometheus_http_request_duration_seconds_bucket{handler="/api/v1/label/:name/values",le="8"} 2
prometheus_http_request_duration_seconds_bucket{handler="/api/v1/label/:name/values",le="20"} 2
prometheus_http_request_duration_seconds_bucket{handler="/api/v1/label/:name/values",le="60"} 2
prometheus_http_request_duration_seconds_bucket{handler="/api/v1/label/:name/values",le="120"} 2
prometheus_http_request_duration_seconds_bucket{handler="/api/v1/label/:name/values",le="+Inf"} 2
prometheus_http_request_duration_seconds_sum{handler="/api/v1/label/:name/values"} 0.044899833
prometheus_http_request_duration_seconds_count{handler="/api/v1/label/:name/values"} 2
prometheus_http_request_duration_seconds_bucket{handler="/api/v1/query",le="0.1"} 2
prometheus_http_request_duration_seconds_bucket{handler="/api/v1/query",le="0.2"} 2
prometheus_http_request_duration_seconds_bucket{handler="/api/v1/query",le="0.4"} 2
prometheus_http_request_duration_seconds_bucket{handler="/api/v1/query",le="1"} 2
prometheus_http_request_duration_seconds_bucket{handler="/api/v1/query",le="3"} 2
prometheus_http_request_duration_seconds_bucket{handler="/api/v1/query",le="8"} 2
prometheus_http_request_duration_seconds_bucket{handler="/api/v1/query",le="20"} 2
prometheus_http_request_duration_seconds_bucket{handler="/api/v1/query",le="60"} 2
prometheus_http_request_duration_seconds_bucket{handler="/api/v1/query",le="120"} 2
prometheus_http_request_duration_seconds_bucket{handler="/api/v1/query",le="+Inf"} 2
prometheus_http_request_duration_seconds_sum{handler="/api/v1/query"} 0.001431704
prometheus_http_request_duration_seconds_count{handler="/api/v1/query"} 2
prometheus_http_request_duration_seconds_bucket{handler="/api/v1/rules",le="0.1"} 3
prometheus_http_request_duration_seconds_bucket{handler="/api/v1/rules",le="0.2"} 3
prometheus_http_request_duration_seconds_bucket{handler="/api/v1/rules",le="0.4"} 3
prometheus_http_request_duration_seconds_bucket{handler="/api/v1/rules",le="1"} 3
prometheus_http_request_duration_seconds_bucket{handler="/api/v1/rules",le="3"} 3
prometheus_http_request_duration_seconds_bucket{handler="/api/v1/rules",le="8"} 3
prometheus_http_request_duration_seconds_bucket{handler="/api/v1/rules",le="20"} 3
prometheus_http_request_duration_seconds_bucket{handler="/api/v1/rules",le="60"} 3
prometheus_http_request_duration_seconds_bucket{handler="/api/v1/rules",le="120"} 3
prometheus_http_request_duration_seconds_bucket{handler="/api/v1/rules",le="+Inf"} 3
prometheus_http_request_duration_seconds_sum{handler="/api/v1/rules"} 0.06705794300000001
prometheus_http_request_duration_seconds_count{handler="/api/v1/rules"} 3
prometheus_http_request_duration_seconds_bucket{handler="/api/v1/targets",le="0.1"} 1
prometheus_http_request_duration_seconds_bucket{handler="/api/v1/targets",le="0.2"} 1
prometheus_http_request_duration_seconds_bucket{handler="/api/v1/targets",le="0.4"} 1
prometheus_http_request_duration_seconds_bucket{handler="/api/v1/targets",le="1"} 1
prometheus_http_request_duration_seconds_bucket{handler="/api/v1/targets",le="3"} 1
prometheus_http_request_duration_seconds_bucket{handler="/api/v1/targets",le="8"} 1
prometheus_http_request_duration_seconds_bucket{handler="/api/v1/targets",le="20"} 1
prometheus_http_request_duration_seconds_bucket{handler="/api/v1/targets",le="60"} 1
prometheus_http_request_duration_seconds_bucket{handler="/api/v1/targets",le="120"} 1
prometheus_http_request_duration_seconds_bucket{handler="/api/v1/targets",le="+Inf"} 1
prometheus_http_request_duration_seconds_sum{handler="/api/v1/targets"} 0.018461351
prometheus_http_request_duration_seconds_count{handler="/api/v1/targets"} 1
prometheus_http_request_duration_seconds_bucket{handler="/favicon.ico",le="0.1"} 3
prometheus_http_request_duration_seconds_bucket{handler="/favicon.ico",le="0.2"} 3
prometheus_http_request_duration_seconds_bucket{handler="/favicon.ico",le="0.4"} 3
prometheus_http_request_duration_seconds_bucket{handler="/favicon.ico",le="1"} 3
prometheus_http_request_duration_seconds_bucket{handler="/favicon.ico",le="3"} 3
prometheus_http_request_duration_seconds_bucket{handler="/favicon.ico",le="8"} 3
prometheus_http_request_duration_seconds_bucket{handler="/favicon.ico",le="20"} 3
prometheus_http_request_duration_seconds_bucket{handler="/favicon.ico",le="60"} 3
prometheus_http_request_duration_seconds_bucket{handler="/favicon.ico",le="120"} 3
prometheus_http_request_duration_seconds_bucket{handler="/favicon.ico",le="+Inf"} 3
prometheus_http_request_duration_seconds_sum{handler="/favicon.ico"} 0.0017768419999999998
prometheus_http_request_duration_seconds_count{handler="/favicon.ico"} 3
prometheus_http_request_duration_seconds_bucket{handler="/graph",le="0.1"} 1
prometheus_http_request_duration_seconds_bucket{handler="/graph",le="0.2"} 1
prometheus_http_request_duration_seconds_bucket{handler="/graph",le="0.4"} 1
prometheus_http_request_duration_seconds_bucket{handler="/graph",le="1"} 1
prometheus_http_request_duration_seconds_bucket{handler="/graph",le="3"} 1
prometheus_http_request_duration_seconds_bucket{handler="/graph",le="8"} 1
prometheus_http_request_duration_seconds_bucket{handler="/graph",le="20"} 1
prometheus_http_request_duration_seconds_bucket{handler="/graph",le="60"} 1
prometheus_http_request_duration_seconds_bucket{handler="/graph",le="120"} 1
prometheus_http_request_duration_seconds_bucket{handler="/graph",le="+Inf"} 1
prometheus_http_request_duration_seconds_sum{handler="/graph"} 0.000257512
prometheus_http_request_duration_seconds_count{handler="/graph"} 1
prometheus_http_request_duration_seconds_bucket{handler="/manifest.json",le="0.1"} 3
prometheus_http_request_duration_seconds_bucket{handler="/manifest.json",le="0.2"} 3
prometheus_http_request_duration_seconds_bucket{handler="/manifest.json",le="0.4"} 3
prometheus_http_request_duration_seconds_bucket{handler="/manifest.json",le="1"} 3
prometheus_http_request_duration_seconds_bucket{handler="/manifest.json",le="3"} 3
prometheus_http_request_duration_seconds_bucket{handler="/manifest.json",le="8"} 3
prometheus_http_request_duration_seconds_bucket{handler="/manifest.json",le="20"} 3
prometheus_http_request_duration_seconds_bucket{handler="/manifest.json",le="60"} 3
prometheus_http_request_duration_seconds_bucket{handler="/manifest.json",le="120"} 3
prometheus_http_request_duration_seconds_bucket{handler="/manifest.json",le="+Inf"} 3
prometheus_http_request_duration_seconds_sum{handler="/manifest.json"} 0.000425453
prometheus_http_request_duration_seconds_count{handler="/manifest.json"} 3
prometheus_http_request_duration_seconds_bucket{handler="/metrics",le="0.1"} 1523
prometheus_http_request_duration_seconds_bucket{handler="/metrics",le="0.2"} 1523
prometheus_http_request_duration_seconds_bucket{handler="/metrics",le="0.4"} 1523
prometheus_http_request_duration_seconds_bucket{handler="/metrics",le="1"} 1523
prometheus_http_request_duration_seconds_bucket{handler="/metrics",le="3"} 1523
prometheus_http_request_duration_seconds_bucket{handler="/metrics",le="8"} 1523
prometheus_http_request_duration_seconds_bucket{handler="/metrics",le="20"} 1523
prometheus_http_request_duration_seconds_bucket{handler="/metrics",le="60"} 1523
prometheus_http_request_duration_seconds_bucket{handler="/metrics",le="120"} 1523
prometheus_http_request_duration_seconds_bucket{handler="/metrics",le="+Inf"} 1523
prometheus_http_request_duration_seconds_sum{handler="/metrics"} 31.192399906000045
prometheus_http_request_duration_seconds_count{handler="/metrics"} 1523
prometheus_http_request_duration_seconds_bucket{handler="/static/*filepath",le="0.1"} 7
prometheus_http_request_duration_seconds_bucket{handler="/static/*filepath",le="0.2"} 7
prometheus_http_request_duration_seconds_bucket{handler="/static/*filepath",le="0.4"} 7
prometheus_http_request_duration_seconds_bucket{handler="/static/*filepath",le="1"} 7
prometheus_http_request_duration_seconds_bucket{handler="/static/*filepath",le="3"} 8
prometheus_http_request_duration_seconds_bucket{handler="/static/*filepath",le="8"} 10
prometheus_http_request_duration_seconds_bucket{handler="/static/*filepath",le="20"} 11
prometheus_http_request_duration_seconds_bucket{handler="/static/*filepath",le="60"} 11
prometheus_http_request_duration_seconds_bucket{handler="/static/*filepath",le="120"} 11
prometheus_http_request_duration_seconds_bucket{handler="/static/*filepath",le="+Inf"} 11
prometheus_http_request_duration_seconds_sum{handler="/static/*filepath"} 25.022674474
prometheus_http_request_duration_seconds_count{handler="/static/*filepath"} 11
# HELP prometheus_http_requests_total Counter of HTTP requests.
# TYPE prometheus_http_requests_total counter
prometheus_http_requests_total{code="200",handler="/-/healthy"} 9135
prometheus_http_requests_total{code="200",handler="/-/ready"} 9139
prometheus_http_requests_total{code="200",handler="/-/reload"} 1
prometheus_http_requests_total{code="200",handler="/alerts"} 3
prometheus_http_requests_total{code="200",handler="/api/v1/label/:name/values"} 2
prometheus_http_requests_total{code="200",handler="/api/v1/query"} 2
prometheus_http_requests_total{code="200",handler="/api/v1/rules"} 3
prometheus_http_requests_total{code="200",handler="/api/v1/targets"} 1
prometheus_http_requests_total{code="200",handler="/favicon.ico"} 3
prometheus_http_requests_total{code="200",handler="/graph"} 1
prometheus_http_requests_total{code="200",handler="/manifest.json"} 3
prometheus_http_requests_total{code="200",handler="/metrics"} 1523
prometheus_http_requests_total{code="200",handler="/static/*filepath"} 11
# HELP prometheus_http_response_size_bytes Histogram of response size for HTTP requests.
# TYPE prometheus_http_response_size_bytes histogram
prometheus_http_response_size_bytes_bucket{handler="/-/healthy",le="100"} 9135
prometheus_http_response_size_bytes_bucket{handler="/-/healthy",le="1000"} 9135
prometheus_http_response_size_bytes_bucket{handler="/-/healthy",le="10000"} 9135
prometheus_http_response_size_bytes_bucket{handler="/-/healthy",le="100000"} 9135
prometheus_http_response_size_bytes_bucket{handler="/-/healthy",le="1e+06"} 9135
prometheus_http_response_size_bytes_bucket{handler="/-/healthy",le="1e+07"} 9135
prometheus_http_response_size_bytes_bucket{handler="/-/healthy",le="1e+08"} 9135
prometheus_http_response_size_bytes_bucket{handler="/-/healthy",le="1e+09"} 9135
prometheus_http_response_size_bytes_bucket{handler="/-/healthy",le="+Inf"} 9135
prometheus_http_response_size_bytes_sum{handler="/-/healthy"} 274050
prometheus_http_response_size_bytes_count{handler="/-/healthy"} 9135
prometheus_http_response_size_bytes_bucket{handler="/-/ready",le="100"} 9139
prometheus_http_response_size_bytes_bucket{handler="/-/ready",le="1000"} 9139
prometheus_http_response_size_bytes_bucket{handler="/-/ready",le="10000"} 9139
prometheus_http_response_size_bytes_bucket{handler="/-/ready",le="100000"} 9139
prometheus_http_response_size_bytes_bucket{handler="/-/ready",le="1e+06"} 9139
prometheus_http_response_size_bytes_bucket{handler="/-/ready",le="1e+07"} 9139
prometheus_http_response_size_bytes_bucket{handler="/-/ready",le="1e+08"} 9139
prometheus_http_response_size_bytes_bucket{handler="/-/ready",le="1e+09"} 9139
prometheus_http_response_size_bytes_bucket{handler="/-/ready",le="+Inf"} 9139
prometheus_http_response_size_bytes_sum{handler="/-/ready"} 255892
prometheus_http_response_size_bytes_count{handler="/-/ready"} 9139
prometheus_http_response_size_bytes_bucket{handler="/-/reload",le="100"} 1
prometheus_http_response_size_bytes_bucket{handler="/-/reload",le="1000"} 1
prometheus_http_response_size_bytes_bucket{handler="/-/reload",le="10000"} 1
prometheus_http_response_size_bytes_bucket{handler="/-/reload",le="100000"} 1
prometheus_http_response_size_bytes_bucket{handler="/-/reload",le="1e+06"} 1
prometheus_http_response_size_bytes_bucket{handler="/-/reload",le="1e+07"} 1
prometheus_http_response_size_bytes_bucket{handler="/-/reload",le="1e+08"} 1
prometheus_http_response_size_bytes_bucket{handler="/-/reload",le="1e+09"} 1
prometheus_http_response_size_bytes_bucket{handler="/-/reload",le="+Inf"} 1
prometheus_http_response_size_bytes_sum{handler="/-/reload"} 0
prometheus_http_response_size_bytes_count{handler="/-/reload"} 1
prometheus_http_response_size_bytes_bucket{handler="/alerts",le="100"} 0
prometheus_http_response_size_bytes_bucket{handler="/alerts",le="1000"} 3
prometheus_http_response_size_bytes_bucket{handler="/alerts",le="10000"} 3
prometheus_http_response_size_bytes_bucket{handler="/alerts",le="100000"} 3
prometheus_http_response_size_bytes_bucket{handler="/alerts",le="1e+06"} 3
prometheus_http_response_size_bytes_bucket{handler="/alerts",le="1e+07"} 3
prometheus_http_response_size_bytes_bucket{handler="/alerts",le="1e+08"} 3
prometheus_http_response_size_bytes_bucket{handler="/alerts",le="1e+09"} 3
prometheus_http_response_size_bytes_bucket{handler="/alerts",le="+Inf"} 3
prometheus_http_response_size_bytes_sum{handler="/alerts"} 2142
prometheus_http_response_size_bytes_count{handler="/alerts"} 3
prometheus_http_response_size_bytes_bucket{handler="/api/v1/label/:name/values",le="100"} 0
prometheus_http_response_size_bytes_bucket{handler="/api/v1/label/:name/values",le="1000"} 0
prometheus_http_response_size_bytes_bucket{handler="/api/v1/label/:name/values",le="10000"} 0
prometheus_http_response_size_bytes_bucket{handler="/api/v1/label/:name/values",le="100000"} 2
prometheus_http_response_size_bytes_bucket{handler="/api/v1/label/:name/values",le="1e+06"} 2
prometheus_http_response_size_bytes_bucket{handler="/api/v1/label/:name/values",le="1e+07"} 2
prometheus_http_response_size_bytes_bucket{handler="/api/v1/label/:name/values",le="1e+08"} 2
prometheus_http_response_size_bytes_bucket{handler="/api/v1/label/:name/values",le="1e+09"} 2
prometheus_http_response_size_bytes_bucket{handler="/api/v1/label/:name/values",le="+Inf"} 2
prometheus_http_response_size_bytes_sum{handler="/api/v1/label/:name/values"} 75442
prometheus_http_response_size_bytes_count{handler="/api/v1/label/:name/values"} 2
prometheus_http_response_size_bytes_bucket{handler="/api/v1/query",le="100"} 0
prometheus_http_response_size_bytes_bucket{handler="/api/v1/query",le="1000"} 2
prometheus_http_response_size_bytes_bucket{handler="/api/v1/query",le="10000"} 2
prometheus_http_response_size_bytes_bucket{handler="/api/v1/query",le="100000"} 2
prometheus_http_response_size_bytes_bucket{handler="/api/v1/query",le="1e+06"} 2
prometheus_http_response_size_bytes_bucket{handler="/api/v1/query",le="1e+07"} 2
prometheus_http_response_size_bytes_bucket{handler="/api/v1/query",le="1e+08"} 2
prometheus_http_response_size_bytes_bucket{handler="/api/v1/query",le="1e+09"} 2
prometheus_http_response_size_bytes_bucket{handler="/api/v1/query",le="+Inf"} 2
prometheus_http_response_size_bytes_sum{handler="/api/v1/query"} 208
prometheus_http_response_size_bytes_count{handler="/api/v1/query"} 2
prometheus_http_response_size_bytes_bucket{handler="/api/v1/rules",le="100"} 0
prometheus_http_response_size_bytes_bucket{handler="/api/v1/rules",le="1000"} 0
prometheus_http_response_size_bytes_bucket{handler="/api/v1/rules",le="10000"} 0
prometheus_http_response_size_bytes_bucket{handler="/api/v1/rules",le="100000"} 3
prometheus_http_response_size_bytes_bucket{handler="/api/v1/rules",le="1e+06"} 3
prometheus_http_response_size_bytes_bucket{handler="/api/v1/rules",le="1e+07"} 3
prometheus_http_response_size_bytes_bucket{handler="/api/v1/rules",le="1e+08"} 3
prometheus_http_response_size_bytes_bucket{handler="/api/v1/rules",le="1e+09"} 3
prometheus_http_response_size_bytes_bucket{handler="/api/v1/rules",le="+Inf"} 3
prometheus_http_response_size_bytes_sum{handler="/api/v1/rules"} 131717
prometheus_http_response_size_bytes_count{handler="/api/v1/rules"} 3
prometheus_http_response_size_bytes_bucket{handler="/api/v1/targets",le="100"} 0
prometheus_http_response_size_bytes_bucket{handler="/api/v1/targets",le="1000"} 0
prometheus_http_response_size_bytes_bucket{handler="/api/v1/targets",le="10000"} 0
prometheus_http_response_size_bytes_bucket{handler="/api/v1/targets",le="100000"} 1
prometheus_http_response_size_bytes_bucket{handler="/api/v1/targets",le="1e+06"} 1
prometheus_http_response_size_bytes_bucket{handler="/api/v1/targets",le="1e+07"} 1
prometheus_http_response_size_bytes_bucket{handler="/api/v1/targets",le="1e+08"} 1
prometheus_http_response_size_bytes_bucket{handler="/api/v1/targets",le="1e+09"} 1
prometheus_http_response_size_bytes_bucket{handler="/api/v1/targets",le="+Inf"} 1
prometheus_http_response_size_bytes_sum{handler="/api/v1/targets"} 22776
prometheus_http_response_size_bytes_count{handler="/api/v1/targets"} 1
prometheus_http_response_size_bytes_bucket{handler="/favicon.ico",le="100"} 0
prometheus_http_response_size_bytes_bucket{handler="/favicon.ico",le="1000"} 0
prometheus_http_response_size_bytes_bucket{handler="/favicon.ico",le="10000"} 0
prometheus_http_response_size_bytes_bucket{handler="/favicon.ico",le="100000"} 3
prometheus_http_response_size_bytes_bucket{handler="/favicon.ico",le="1e+06"} 3
prometheus_http_response_size_bytes_bucket{handler="/favicon.ico",le="1e+07"} 3
prometheus_http_response_size_bytes_bucket{handler="/favicon.ico",le="1e+08"} 3
prometheus_http_response_size_bytes_bucket{handler="/favicon.ico",le="1e+09"} 3
prometheus_http_response_size_bytes_bucket{handler="/favicon.ico",le="+Inf"} 3
prometheus_http_response_size_bytes_sum{handler="/favicon.ico"} 45258
prometheus_http_response_size_bytes_count{handler="/favicon.ico"} 3
prometheus_http_response_size_bytes_bucket{handler="/graph",le="100"} 0
prometheus_http_response_size_bytes_bucket{handler="/graph",le="1000"} 1
prometheus_http_response_size_bytes_bucket{handler="/graph",le="10000"} 1
prometheus_http_response_size_bytes_bucket{handler="/graph",le="100000"} 1
prometheus_http_response_size_bytes_bucket{handler="/graph",le="1e+06"} 1
prometheus_http_response_size_bytes_bucket{handler="/graph",le="1e+07"} 1
prometheus_http_response_size_bytes_bucket{handler="/graph",le="1e+08"} 1
prometheus_http_response_size_bytes_bucket{handler="/graph",le="1e+09"} 1
prometheus_http_response_size_bytes_bucket{handler="/graph",le="+Inf"} 1
prometheus_http_response_size_bytes_sum{handler="/graph"} 714
prometheus_http_response_size_bytes_count{handler="/graph"} 1
prometheus_http_response_size_bytes_bucket{handler="/manifest.json",le="100"} 0
prometheus_http_response_size_bytes_bucket{handler="/manifest.json",le="1000"} 3
prometheus_http_response_size_bytes_bucket{handler="/manifest.json",le="10000"} 3
prometheus_http_response_size_bytes_bucket{handler="/manifest.json",le="100000"} 3
prometheus_http_response_size_bytes_bucket{handler="/manifest.json",le="1e+06"} 3
prometheus_http_response_size_bytes_bucket{handler="/manifest.json",le="1e+07"} 3
prometheus_http_response_size_bytes_bucket{handler="/manifest.json",le="1e+08"} 3
prometheus_http_response_size_bytes_bucket{handler="/manifest.json",le="1e+09"} 3
prometheus_http_response_size_bytes_bucket{handler="/manifest.json",le="+Inf"} 3
prometheus_http_response_size_bytes_sum{handler="/manifest.json"} 954
prometheus_http_response_size_bytes_count{handler="/manifest.json"} 3
prometheus_http_response_size_bytes_bucket{handler="/metrics",le="100"} 0
prometheus_http_response_size_bytes_bucket{handler="/metrics",le="1000"} 0
prometheus_http_response_size_bytes_bucket{handler="/metrics",le="10000"} 0
prometheus_http_response_size_bytes_bucket{handler="/metrics",le="100000"} 1523
prometheus_http_response_size_bytes_bucket{handler="/metrics",le="1e+06"} 1523
prometheus_http_response_size_bytes_bucket{handler="/metrics",le="1e+07"} 1523
prometheus_http_response_size_bytes_bucket{handler="/metrics",le="1e+08"} 1523
prometheus_http_response_size_bytes_bucket{handler="/metrics",le="1e+09"} 1523
prometheus_http_response_size_bytes_bucket{handler="/metrics",le="+Inf"} 1523
prometheus_http_response_size_bytes_sum{handler="/metrics"} 4.2759237e+07
prometheus_http_response_size_bytes_count{handler="/metrics"} 1523
prometheus_http_response_size_bytes_bucket{handler="/static/*filepath",le="100"} 0
prometheus_http_response_size_bytes_bucket{handler="/static/*filepath",le="1000"} 0
prometheus_http_response_size_bytes_bucket{handler="/static/*filepath",le="10000"} 3
prometheus_http_response_size_bytes_bucket{handler="/static/*filepath",le="100000"} 3
prometheus_http_response_size_bytes_bucket{handler="/static/*filepath",le="1e+06"} 7
prometheus_http_response_size_bytes_bucket{handler="/static/*filepath",le="1e+07"} 11
prometheus_http_response_size_bytes_bucket{handler="/static/*filepath",le="1e+08"} 11
prometheus_http_response_size_bytes_bucket{handler="/static/*filepath",le="1e+09"} 11
prometheus_http_response_size_bytes_bucket{handler="/static/*filepath",le="+Inf"} 11
prometheus_http_response_size_bytes_sum{handler="/static/*filepath"} 1.0837855e+07
prometheus_http_response_size_bytes_count{handler="/static/*filepath"} 11
# HELP prometheus_notifications_alertmanagers_discovered The number of alertmanagers discovered and active.
# TYPE prometheus_notifications_alertmanagers_discovered gauge
prometheus_notifications_alertmanagers_discovered 2
# HELP prometheus_notifications_dropped_total Total number of alerts dropped due to errors when sending to Alertmanager.
# TYPE prometheus_notifications_dropped_total counter
prometheus_notifications_dropped_total 0
# HELP prometheus_notifications_errors_total Total number of errors sending alert notifications.
# TYPE prometheus_notifications_errors_total counter
prometheus_notifications_errors_total{alertmanager="http://10.129.252.159:9093/api/v2/alerts"} 0
prometheus_notifications_errors_total{alertmanager="https://system-monitor-staging.k8s.test.org/api/v2/alerts"} 20666
# HELP prometheus_notifications_latency_seconds Latency quantiles for sending alert notifications.
# TYPE prometheus_notifications_latency_seconds summary
prometheus_notifications_latency_seconds{alertmanager="http://10.129.252.159:9093/api/v2/alerts",quantile="0.5"} 0.001921053
prometheus_notifications_latency_seconds{alertmanager="http://10.129.252.159:9093/api/v2/alerts",quantile="0.9"} 0.005416296
prometheus_notifications_latency_seconds{alertmanager="http://10.129.252.159:9093/api/v2/alerts",quantile="0.99"} 0.010897484
prometheus_notifications_latency_seconds_sum{alertmanager="http://10.129.252.159:9093/api/v2/alerts"} 56.228670395999565
prometheus_notifications_latency_seconds_count{alertmanager="http://10.129.252.159:9093/api/v2/alerts"} 20666
prometheus_notifications_latency_seconds{alertmanager="https://system-monitor-staging.k8s.test.org/api/v2/alerts",quantile="0.5"} 0.003821372
prometheus_notifications_latency_seconds{alertmanager="https://system-monitor-staging.k8s.test.org/api/v2/alerts",quantile="0.9"} 0.006877826
prometheus_notifications_latency_seconds{alertmanager="https://system-monitor-staging.k8s.test.org/api/v2/alerts",quantile="0.99"} 0.011607642
prometheus_notifications_latency_seconds_sum{alertmanager="https://system-monitor-staging.k8s.test.org/api/v2/alerts"} 92.34316748099977
prometheus_notifications_latency_seconds_count{alertmanager="https://system-monitor-staging.k8s.test.org/api/v2/alerts"} 20666
# HELP prometheus_notifications_queue_capacity The capacity of the alert notifications queue.
# TYPE prometheus_notifications_queue_capacity gauge
prometheus_notifications_queue_capacity 10000
# HELP prometheus_notifications_queue_length The number of alert notifications in the queue.
# TYPE prometheus_notifications_queue_length gauge
prometheus_notifications_queue_length 0
# HELP prometheus_notifications_sent_total Total number of alerts sent.
# TYPE prometheus_notifications_sent_total counter
prometheus_notifications_sent_total{alertmanager="http://10.129.252.159:9093/api/v2/alerts"} 53156
prometheus_notifications_sent_total{alertmanager="https://system-monitor-staging.k8s.test.org/api/v2/alerts"} 53156
# HELP prometheus_ready Whether Prometheus startup was fully completed and the server is ready for normal operation.
# TYPE prometheus_ready gauge
prometheus_ready 1
# HELP prometheus_remote_storage_bytes_total The total number of bytes of data (not metadata) sent by the queue after compression. Note that when exemplars over remote write is enabled the exemplars included in a remote write request count towards this metric.
# TYPE prometheus_remote_storage_bytes_total counter
prometheus_remote_storage_bytes_total{remote_name="28d637",url="http://10.106.2.6:8187/api/v3/prometheus/write"} 2.5383190411e+10
prometheus_remote_storage_bytes_total{remote_name="4bfe0a",url="https://system-monitor-staging.k8s.test.org:8427/api/v1/write"} 0
prometheus_remote_storage_bytes_total{remote_name="c93768",url="http://10.101.17.191:1429/api/v1/write"} 0
# HELP prometheus_remote_storage_enqueue_retries_total Total number of times enqueue has failed because a shards queue was full.
# TYPE prometheus_remote_storage_enqueue_retries_total counter
prometheus_remote_storage_enqueue_retries_total{remote_name="28d637",url="http://10.106.2.6:8187/api/v3/prometheus/write"} 543127
prometheus_remote_storage_enqueue_retries_total{remote_name="4bfe0a",url="https://system-monitor-staging.k8s.test.org:8427/api/v1/write"} 9145
prometheus_remote_storage_enqueue_retries_total{remote_name="c93768",url="http://10.101.17.191:1429/api/v1/write"} 9145
# HELP prometheus_remote_storage_exemplars_dropped_total Total number of exemplars which were dropped after being read from the WAL before being sent via remote write, either via relabelling or unintentionally because of an unknown reference ID.
# TYPE prometheus_remote_storage_exemplars_dropped_total counter
prometheus_remote_storage_exemplars_dropped_total{remote_name="28d637",url="http://10.106.2.6:8187/api/v3/prometheus/write"} 0
prometheus_remote_storage_exemplars_dropped_total{remote_name="4bfe0a",url="https://system-monitor-staging.k8s.test.org:8427/api/v1/write"} 0
prometheus_remote_storage_exemplars_dropped_total{remote_name="c93768",url="http://10.101.17.191:1429/api/v1/write"} 0
# HELP prometheus_remote_storage_exemplars_failed_total Total number of exemplars which failed on send to remote storage, non-recoverable errors.
# TYPE prometheus_remote_storage_exemplars_failed_total counter
prometheus_remote_storage_exemplars_failed_total{remote_name="28d637",url="http://10.106.2.6:8187/api/v3/prometheus/write"} 0
prometheus_remote_storage_exemplars_failed_total{remote_name="4bfe0a",url="https://system-monitor-staging.k8s.test.org:8427/api/v1/write"} 0
prometheus_remote_storage_exemplars_failed_total{remote_name="c93768",url="http://10.101.17.191:1429/api/v1/write"} 0
# HELP prometheus_remote_storage_exemplars_in_total Exemplars in to remote storage, compare to exemplars out for queue managers.
# TYPE prometheus_remote_storage_exemplars_in_total counter
prometheus_remote_storage_exemplars_in_total 0
# HELP prometheus_remote_storage_exemplars_pending The number of exemplars pending in the queues shards to be sent to the remote storage.
# TYPE prometheus_remote_storage_exemplars_pending gauge
prometheus_remote_storage_exemplars_pending{remote_name="28d637",url="http://10.106.2.6:8187/api/v3/prometheus/write"} 0
prometheus_remote_storage_exemplars_pending{remote_name="4bfe0a",url="https://system-monitor-staging.k8s.test.org:8427/api/v1/write"} 0
prometheus_remote_storage_exemplars_pending{remote_name="c93768",url="http://10.101.17.191:1429/api/v1/write"} 0
# HELP prometheus_remote_storage_exemplars_retried_total Total number of exemplars which failed on send to remote storage but were retried because the send error was recoverable.
# TYPE prometheus_remote_storage_exemplars_retried_total counter
prometheus_remote_storage_exemplars_retried_total{remote_name="28d637",url="http://10.106.2.6:8187/api/v3/prometheus/write"} 0
prometheus_remote_storage_exemplars_retried_total{remote_name="4bfe0a",url="https://system-monitor-staging.k8s.test.org:8427/api/v1/write"} 0
prometheus_remote_storage_exemplars_retried_total{remote_name="c93768",url="http://10.101.17.191:1429/api/v1/write"} 0
# HELP prometheus_remote_storage_exemplars_total Total number of exemplars sent to remote storage.
# TYPE prometheus_remote_storage_exemplars_total counter
prometheus_remote_storage_exemplars_total{remote_name="28d637",url="http://10.106.2.6:8187/api/v3/prometheus/write"} 0
prometheus_remote_storage_exemplars_total{remote_name="4bfe0a",url="https://system-monitor-staging.k8s.test.org:8427/api/v1/write"} 0
prometheus_remote_storage_exemplars_total{remote_name="c93768",url="http://10.101.17.191:1429/api/v1/write"} 0
# HELP prometheus_remote_storage_highest_timestamp_in_seconds Highest timestamp that has come into the remote storage via the Appender interface, in seconds since epoch.
# TYPE prometheus_remote_storage_highest_timestamp_in_seconds gauge
prometheus_remote_storage_highest_timestamp_in_seconds 1.67264724e+09
# HELP prometheus_remote_storage_histograms_dropped_total Total number of histograms which were dropped after being read from the WAL before being sent via remote write, either via relabelling or unintentionally because of an unknown reference ID.
# TYPE prometheus_remote_storage_histograms_dropped_total counter
prometheus_remote_storage_histograms_dropped_total{remote_name="28d637",url="http://10.106.2.6:8187/api/v3/prometheus/write"} 0
prometheus_remote_storage_histograms_dropped_total{remote_name="4bfe0a",url="https://system-monitor-staging.k8s.test.org:8427/api/v1/write"} 0
prometheus_remote_storage_histograms_dropped_total{remote_name="c93768",url="http://10.101.17.191:1429/api/v1/write"} 0
# HELP prometheus_remote_storage_histograms_failed_total Total number of histograms which failed on send to remote storage, non-recoverable errors.
# TYPE prometheus_remote_storage_histograms_failed_total counter
prometheus_remote_storage_histograms_failed_total{remote_name="28d637",url="http://10.106.2.6:8187/api/v3/prometheus/write"} 0
prometheus_remote_storage_histograms_failed_total{remote_name="4bfe0a",url="https://system-monitor-staging.k8s.test.org:8427/api/v1/write"} 0
prometheus_remote_storage_histograms_failed_total{remote_name="c93768",url="http://10.101.17.191:1429/api/v1/write"} 0
# HELP prometheus_remote_storage_histograms_in_total HistogramSamples in to remote storage, compare to histograms out for queue managers.
# TYPE prometheus_remote_storage_histograms_in_total counter
prometheus_remote_storage_histograms_in_total 0
# HELP prometheus_remote_storage_histograms_pending The number of histograms pending in the queues shards to be sent to the remote storage.
# TYPE prometheus_remote_storage_histograms_pending gauge
prometheus_remote_storage_histograms_pending{remote_name="28d637",url="http://10.106.2.6:8187/api/v3/prometheus/write"} 0
prometheus_remote_storage_histograms_pending{remote_name="4bfe0a",url="https://system-monitor-staging.k8s.test.org:8427/api/v1/write"} 0
prometheus_remote_storage_histograms_pending{remote_name="c93768",url="http://10.101.17.191:1429/api/v1/write"} 0
# HELP prometheus_remote_storage_histograms_retried_total Total number of histograms which failed on send to remote storage but were retried because the send error was recoverable.
# TYPE prometheus_remote_storage_histograms_retried_total counter
prometheus_remote_storage_histograms_retried_total{remote_name="28d637",url="http://10.106.2.6:8187/api/v3/prometheus/write"} 0
prometheus_remote_storage_histograms_retried_total{remote_name="4bfe0a",url="https://system-monitor-staging.k8s.test.org:8427/api/v1/write"} 0
prometheus_remote_storage_histograms_retried_total{remote_name="c93768",url="http://10.101.17.191:1429/api/v1/write"} 0
# HELP prometheus_remote_storage_histograms_total Total number of histograms sent to remote storage.
# TYPE prometheus_remote_storage_histograms_total counter
prometheus_remote_storage_histograms_total{remote_name="28d637",url="http://10.106.2.6:8187/api/v3/prometheus/write"} 0
prometheus_remote_storage_histograms_total{remote_name="4bfe0a",url="https://system-monitor-staging.k8s.test.org:8427/api/v1/write"} 0
prometheus_remote_storage_histograms_total{remote_name="c93768",url="http://10.101.17.191:1429/api/v1/write"} 0
# HELP prometheus_remote_storage_max_samples_per_send The maximum number of samples to be sent, in a single request, to the remote storage. Note that, when sending of exemplars over remote write is enabled, exemplars count towards this limt.
# TYPE prometheus_remote_storage_max_samples_per_send gauge
prometheus_remote_storage_max_samples_per_send{remote_name="28d637",url="http://10.106.2.6:8187/api/v3/prometheus/write"} 500
prometheus_remote_storage_max_samples_per_send{remote_name="4bfe0a",url="https://system-monitor-staging.k8s.test.org:8427/api/v1/write"} 500
prometheus_remote_storage_max_samples_per_send{remote_name="c93768",url="http://10.101.17.191:1429/api/v1/write"} 500
# HELP prometheus_remote_storage_metadata_bytes_total The total number of bytes of metadata sent by the queue after compression.
# TYPE prometheus_remote_storage_metadata_bytes_total counter
prometheus_remote_storage_metadata_bytes_total{remote_name="28d637",url="http://10.106.2.6:8187/api/v3/prometheus/write"} 1.27046925e+08
prometheus_remote_storage_metadata_bytes_total{remote_name="4bfe0a",url="https://system-monitor-staging.k8s.test.org:8427/api/v1/write"} 0
prometheus_remote_storage_metadata_bytes_total{remote_name="c93768",url="http://10.101.17.191:1429/api/v1/write"} 0
# HELP prometheus_remote_storage_metadata_failed_total Total number of metadata entries which failed on send to remote storage, non-recoverable errors.
# TYPE prometheus_remote_storage_metadata_failed_total counter
prometheus_remote_storage_metadata_failed_total{remote_name="28d637",url="http://10.106.2.6:8187/api/v3/prometheus/write"} 0
prometheus_remote_storage_metadata_failed_total{remote_name="4bfe0a",url="https://system-monitor-staging.k8s.test.org:8427/api/v1/write"} 0
prometheus_remote_storage_metadata_failed_total{remote_name="c93768",url="http://10.101.17.191:1429/api/v1/write"} 0
# HELP prometheus_remote_storage_metadata_retried_total Total number of metadata entries which failed on send to remote storage but were retried because the send error was recoverable.
# TYPE prometheus_remote_storage_metadata_retried_total counter
prometheus_remote_storage_metadata_retried_total{remote_name="28d637",url="http://10.106.2.6:8187/api/v3/prometheus/write"} 0
prometheus_remote_storage_metadata_retried_total{remote_name="4bfe0a",url="https://system-monitor-staging.k8s.test.org:8427/api/v1/write"} 4.555e+06
prometheus_remote_storage_metadata_retried_total{remote_name="c93768",url="http://10.101.17.191:1429/api/v1/write"} 4.564e+06
# HELP prometheus_remote_storage_metadata_total Total number of metadata entries sent to remote storage.
# TYPE prometheus_remote_storage_metadata_total counter
prometheus_remote_storage_metadata_total{remote_name="28d637",url="http://10.106.2.6:8187/api/v3/prometheus/write"} 4.644748e+06
prometheus_remote_storage_metadata_total{remote_name="4bfe0a",url="https://system-monitor-staging.k8s.test.org:8427/api/v1/write"} 0
prometheus_remote_storage_metadata_total{remote_name="c93768",url="http://10.101.17.191:1429/api/v1/write"} 0
# HELP prometheus_remote_storage_queue_highest_sent_timestamp_seconds Timestamp from a WAL sample, the highest timestamp successfully sent by this queue, in seconds since epoch.
# TYPE prometheus_remote_storage_queue_highest_sent_timestamp_seconds gauge
prometheus_remote_storage_queue_highest_sent_timestamp_seconds{remote_name="28d637",url="http://10.106.2.6:8187/api/v3/prometheus/write"} 1.67264724e+09
# HELP prometheus_remote_storage_samples_dropped_total Total number of samples which were dropped after being read from the WAL before being sent via remote write, either via relabelling or unintentionally because of an unknown reference ID.
# TYPE prometheus_remote_storage_samples_dropped_total counter
prometheus_remote_storage_samples_dropped_total{remote_name="28d637",url="http://10.106.2.6:8187/api/v3/prometheus/write"} 0
prometheus_remote_storage_samples_dropped_total{remote_name="4bfe0a",url="https://system-monitor-staging.k8s.test.org:8427/api/v1/write"} 0
prometheus_remote_storage_samples_dropped_total{remote_name="c93768",url="http://10.101.17.191:1429/api/v1/write"} 0
# HELP prometheus_remote_storage_samples_failed_total Total number of samples which failed on send to remote storage, non-recoverable errors.
# TYPE prometheus_remote_storage_samples_failed_total counter
prometheus_remote_storage_samples_failed_total{remote_name="28d637",url="http://10.106.2.6:8187/api/v3/prometheus/write"} 0
prometheus_remote_storage_samples_failed_total{remote_name="4bfe0a",url="https://system-monitor-staging.k8s.test.org:8427/api/v1/write"} 0
prometheus_remote_storage_samples_failed_total{remote_name="c93768",url="http://10.101.17.191:1429/api/v1/write"} 0
# HELP prometheus_remote_storage_samples_in_total Samples in to remote storage, compare to samples out for queue managers.
# TYPE prometheus_remote_storage_samples_in_total counter
prometheus_remote_storage_samples_in_total 4.82905572e+08
# HELP prometheus_remote_storage_samples_pending The number of samples pending in the queues shards to be sent to the remote storage.
# TYPE prometheus_remote_storage_samples_pending gauge
prometheus_remote_storage_samples_pending{remote_name="28d637",url="http://10.106.2.6:8187/api/v3/prometheus/write"} 268
prometheus_remote_storage_samples_pending{remote_name="4bfe0a",url="https://system-monitor-staging.k8s.test.org:8427/api/v1/write"} 3499
prometheus_remote_storage_samples_pending{remote_name="c93768",url="http://10.101.17.191:1429/api/v1/write"} 3499
# HELP prometheus_remote_storage_samples_retried_total Total number of samples which failed on send to remote storage but were retried because the send error was recoverable.
# TYPE prometheus_remote_storage_samples_retried_total counter
prometheus_remote_storage_samples_retried_total{remote_name="28d637",url="http://10.106.2.6:8187/api/v3/prometheus/write"} 0
prometheus_remote_storage_samples_retried_total{remote_name="4bfe0a",url="https://system-monitor-staging.k8s.test.org:8427/api/v1/write"} 4.56e+06
prometheus_remote_storage_samples_retried_total{remote_name="c93768",url="http://10.101.17.191:1429/api/v1/write"} 4.5695e+06
# HELP prometheus_remote_storage_samples_total Total number of samples sent to remote storage.
# TYPE prometheus_remote_storage_samples_total counter
prometheus_remote_storage_samples_total{remote_name="28d637",url="http://10.106.2.6:8187/api/v3/prometheus/write"} 4.82902e+08
prometheus_remote_storage_samples_total{remote_name="4bfe0a",url="https://system-monitor-staging.k8s.test.org:8427/api/v1/write"} 4.5605e+06
prometheus_remote_storage_samples_total{remote_name="c93768",url="http://10.101.17.191:1429/api/v1/write"} 4.57e+06
# HELP prometheus_remote_storage_sent_batch_duration_seconds Duration of send calls to the remote storage.
# TYPE prometheus_remote_storage_sent_batch_duration_seconds histogram
prometheus_remote_storage_sent_batch_duration_seconds_bucket{remote_name="28d637",url="http://10.106.2.6:8187/api/v3/prometheus/write",le="0.005"} 621455
prometheus_remote_storage_sent_batch_duration_seconds_bucket{remote_name="28d637",url="http://10.106.2.6:8187/api/v3/prometheus/write",le="0.01"} 834448
prometheus_remote_storage_sent_batch_duration_seconds_bucket{remote_name="28d637",url="http://10.106.2.6:8187/api/v3/prometheus/write",le="0.025"} 953546
prometheus_remote_storage_sent_batch_duration_seconds_bucket{remote_name="28d637",url="http://10.106.2.6:8187/api/v3/prometheus/write",le="0.05"} 972694
prometheus_remote_storage_sent_batch_duration_seconds_bucket{remote_name="28d637",url="http://10.106.2.6:8187/api/v3/prometheus/write",le="0.1"} 975179
prometheus_remote_storage_sent_batch_duration_seconds_bucket{remote_name="28d637",url="http://10.106.2.6:8187/api/v3/prometheus/write",le="0.25"} 975603
prometheus_remote_storage_sent_batch_duration_seconds_bucket{remote_name="28d637",url="http://10.106.2.6:8187/api/v3/prometheus/write",le="0.5"} 975657
prometheus_remote_storage_sent_batch_duration_seconds_bucket{remote_name="28d637",url="http://10.106.2.6:8187/api/v3/prometheus/write",le="1"} 975674
prometheus_remote_storage_sent_batch_duration_seconds_bucket{remote_name="28d637",url="http://10.106.2.6:8187/api/v3/prometheus/write",le="2.5"} 975691
prometheus_remote_storage_sent_batch_duration_seconds_bucket{remote_name="28d637",url="http://10.106.2.6:8187/api/v3/prometheus/write",le="5"} 975692
prometheus_remote_storage_sent_batch_duration_seconds_bucket{remote_name="28d637",url="http://10.106.2.6:8187/api/v3/prometheus/write",le="10"} 975694
prometheus_remote_storage_sent_batch_duration_seconds_bucket{remote_name="28d637",url="http://10.106.2.6:8187/api/v3/prometheus/write",le="25"} 975694
prometheus_remote_storage_sent_batch_duration_seconds_bucket{remote_name="28d637",url="http://10.106.2.6:8187/api/v3/prometheus/write",le="60"} 975694
prometheus_remote_storage_sent_batch_duration_seconds_bucket{remote_name="28d637",url="http://10.106.2.6:8187/api/v3/prometheus/write",le="120"} 975694
prometheus_remote_storage_sent_batch_duration_seconds_bucket{remote_name="28d637",url="http://10.106.2.6:8187/api/v3/prometheus/write",le="300"} 975694
prometheus_remote_storage_sent_batch_duration_seconds_bucket{remote_name="28d637",url="http://10.106.2.6:8187/api/v3/prometheus/write",le="+Inf"} 975694
prometheus_remote_storage_sent_batch_duration_seconds_sum{remote_name="28d637",url="http://10.106.2.6:8187/api/v3/prometheus/write"} 5905.673532565018
prometheus_remote_storage_sent_batch_duration_seconds_count{remote_name="28d637",url="http://10.106.2.6:8187/api/v3/prometheus/write"} 975694
prometheus_remote_storage_sent_batch_duration_seconds_bucket{remote_name="4bfe0a",url="https://system-monitor-staging.k8s.test.org:8427/api/v1/write",le="0.005"} 1
prometheus_remote_storage_sent_batch_duration_seconds_bucket{remote_name="4bfe0a",url="https://system-monitor-staging.k8s.test.org:8427/api/v1/write",le="0.01"} 8116
prometheus_remote_storage_sent_batch_duration_seconds_bucket{remote_name="4bfe0a",url="https://system-monitor-staging.k8s.test.org:8427/api/v1/write",le="0.025"} 17914
prometheus_remote_storage_sent_batch_duration_seconds_bucket{remote_name="4bfe0a",url="https://system-monitor-staging.k8s.test.org:8427/api/v1/write",le="0.05"} 18149
prometheus_remote_storage_sent_batch_duration_seconds_bucket{remote_name="4bfe0a",url="https://system-monitor-staging.k8s.test.org:8427/api/v1/write",le="0.1"} 18200
prometheus_remote_storage_sent_batch_duration_seconds_bucket{remote_name="4bfe0a",url="https://system-monitor-staging.k8s.test.org:8427/api/v1/write",le="0.25"} 18227
prometheus_remote_storage_sent_batch_duration_seconds_bucket{remote_name="4bfe0a",url="https://system-monitor-staging.k8s.test.org:8427/api/v1/write",le="0.5"} 18232
prometheus_remote_storage_sent_batch_duration_seconds_bucket{remote_name="4bfe0a",url="https://system-monitor-staging.k8s.test.org:8427/api/v1/write",le="1"} 18232
prometheus_remote_storage_sent_batch_duration_seconds_bucket{remote_name="4bfe0a",url="https://system-monitor-staging.k8s.test.org:8427/api/v1/write",le="2.5"} 18232
prometheus_remote_storage_sent_batch_duration_seconds_bucket{remote_name="4bfe0a",url="https://system-monitor-staging.k8s.test.org:8427/api/v1/write",le="5"} 18232
prometheus_remote_storage_sent_batch_duration_seconds_bucket{remote_name="4bfe0a",url="https://system-monitor-staging.k8s.test.org:8427/api/v1/write",le="10"} 18232
prometheus_remote_storage_sent_batch_duration_seconds_bucket{remote_name="4bfe0a",url="https://system-monitor-staging.k8s.test.org:8427/api/v1/write",le="25"} 18232
prometheus_remote_storage_sent_batch_duration_seconds_bucket{remote_name="4bfe0a",url="https://system-monitor-staging.k8s.test.org:8427/api/v1/write",le="60"} 18232
prometheus_remote_storage_sent_batch_duration_seconds_bucket{remote_name="4bfe0a",url="https://system-monitor-staging.k8s.test.org:8427/api/v1/write",le="120"} 18232
prometheus_remote_storage_sent_batch_duration_seconds_bucket{remote_name="4bfe0a",url="https://system-monitor-staging.k8s.test.org:8427/api/v1/write",le="300"} 18232
prometheus_remote_storage_sent_batch_duration_seconds_bucket{remote_name="4bfe0a",url="https://system-monitor-staging.k8s.test.org:8427/api/v1/write",le="+Inf"} 18232
prometheus_remote_storage_sent_batch_duration_seconds_sum{remote_name="4bfe0a",url="https://system-monitor-staging.k8s.test.org:8427/api/v1/write"} 213.30928544100072
prometheus_remote_storage_sent_batch_duration_seconds_count{remote_name="4bfe0a",url="https://system-monitor-staging.k8s.test.org:8427/api/v1/write"} 18232
prometheus_remote_storage_sent_batch_duration_seconds_bucket{remote_name="c93768",url="http://10.101.17.191:1429/api/v1/write",le="0.005"} 18062
prometheus_remote_storage_sent_batch_duration_seconds_bucket{remote_name="c93768",url="http://10.101.17.191:1429/api/v1/write",le="0.01"} 18181
prometheus_remote_storage_sent_batch_duration_seconds_bucket{remote_name="c93768",url="http://10.101.17.191:1429/api/v1/write",le="0.025"} 18250
prometheus_remote_storage_sent_batch_duration_seconds_bucket{remote_name="c93768",url="http://10.101.17.191:1429/api/v1/write",le="0.05"} 18263
prometheus_remote_storage_sent_batch_duration_seconds_bucket{remote_name="c93768",url="http://10.101.17.191:1429/api/v1/write",le="0.1"} 18266
prometheus_remote_storage_sent_batch_duration_seconds_bucket{remote_name="c93768",url="http://10.101.17.191:1429/api/v1/write",le="0.25"} 18267
prometheus_remote_storage_sent_batch_duration_seconds_bucket{remote_name="c93768",url="http://10.101.17.191:1429/api/v1/write",le="0.5"} 18267
prometheus_remote_storage_sent_batch_duration_seconds_bucket{remote_name="c93768",url="http://10.101.17.191:1429/api/v1/write",le="1"} 18267
prometheus_remote_storage_sent_batch_duration_seconds_bucket{remote_name="c93768",url="http://10.101.17.191:1429/api/v1/write",le="2.5"} 18269
prometheus_remote_storage_sent_batch_duration_seconds_bucket{remote_name="c93768",url="http://10.101.17.191:1429/api/v1/write",le="5"} 18269
prometheus_remote_storage_sent_batch_duration_seconds_bucket{remote_name="c93768",url="http://10.101.17.191:1429/api/v1/write",le="10"} 18269
prometheus_remote_storage_sent_batch_duration_seconds_bucket{remote_name="c93768",url="http://10.101.17.191:1429/api/v1/write",le="25"} 18269
prometheus_remote_storage_sent_batch_duration_seconds_bucket{remote_name="c93768",url="http://10.101.17.191:1429/api/v1/write",le="60"} 18269
prometheus_remote_storage_sent_batch_duration_seconds_bucket{remote_name="c93768",url="http://10.101.17.191:1429/api/v1/write",le="120"} 18269
prometheus_remote_storage_sent_batch_duration_seconds_bucket{remote_name="c93768",url="http://10.101.17.191:1429/api/v1/write",le="300"} 18269
prometheus_remote_storage_sent_batch_duration_seconds_bucket{remote_name="c93768",url="http://10.101.17.191:1429/api/v1/write",le="+Inf"} 18269
prometheus_remote_storage_sent_batch_duration_seconds_sum{remote_name="c93768",url="http://10.101.17.191:1429/api/v1/write"} 27.462209213000108
prometheus_remote_storage_sent_batch_duration_seconds_count{remote_name="c93768",url="http://10.101.17.191:1429/api/v1/write"} 18269
# HELP prometheus_remote_storage_shard_capacity The capacity of each shard of the queue used for parallel sending to the remote storage.
# TYPE prometheus_remote_storage_shard_capacity gauge
prometheus_remote_storage_shard_capacity{remote_name="28d637",url="http://10.106.2.6:8187/api/v3/prometheus/write"} 2500
prometheus_remote_storage_shard_capacity{remote_name="4bfe0a",url="https://system-monitor-staging.k8s.test.org:8427/api/v1/write"} 2500
prometheus_remote_storage_shard_capacity{remote_name="c93768",url="http://10.101.17.191:1429/api/v1/write"} 2500
# HELP prometheus_remote_storage_shards The number of shards used for parallel sending to the remote storage.
# TYPE prometheus_remote_storage_shards gauge
prometheus_remote_storage_shards{remote_name="28d637",url="http://10.106.2.6:8187/api/v3/prometheus/write"} 1
prometheus_remote_storage_shards{remote_name="4bfe0a",url="https://system-monitor-staging.k8s.test.org:8427/api/v1/write"} 1
prometheus_remote_storage_shards{remote_name="c93768",url="http://10.101.17.191:1429/api/v1/write"} 1
# HELP prometheus_remote_storage_shards_desired The number of shards that the queues shard calculation wants to run based on the rate of samples in vs. samples out.
# TYPE prometheus_remote_storage_shards_desired gauge
prometheus_remote_storage_shards_desired{remote_name="28d637",url="http://10.106.2.6:8187/api/v3/prometheus/write"} 0.17990070963374152
prometheus_remote_storage_shards_desired{remote_name="4bfe0a",url="https://system-monitor-staging.k8s.test.org:8427/api/v1/write"} 1
prometheus_remote_storage_shards_desired{remote_name="c93768",url="http://10.101.17.191:1429/api/v1/write"} 1
# HELP prometheus_remote_storage_shards_max The maximum number of shards that the queue is allowed to run.
# TYPE prometheus_remote_storage_shards_max gauge
prometheus_remote_storage_shards_max{remote_name="28d637",url="http://10.106.2.6:8187/api/v3/prometheus/write"} 200
prometheus_remote_storage_shards_max{remote_name="4bfe0a",url="https://system-monitor-staging.k8s.test.org:8427/api/v1/write"} 200
prometheus_remote_storage_shards_max{remote_name="c93768",url="http://10.101.17.191:1429/api/v1/write"} 200
# HELP prometheus_remote_storage_shards_min The minimum number of shards that the queue is allowed to run.
# TYPE prometheus_remote_storage_shards_min gauge
prometheus_remote_storage_shards_min{remote_name="28d637",url="http://10.106.2.6:8187/api/v3/prometheus/write"} 1
prometheus_remote_storage_shards_min{remote_name="4bfe0a",url="https://system-monitor-staging.k8s.test.org:8427/api/v1/write"} 1
prometheus_remote_storage_shards_min{remote_name="c93768",url="http://10.101.17.191:1429/api/v1/write"} 1
# HELP prometheus_remote_storage_string_interner_zero_reference_releases_total The number of times release has been called for strings that are not interned.
# TYPE prometheus_remote_storage_string_interner_zero_reference_releases_total counter
prometheus_remote_storage_string_interner_zero_reference_releases_total 0
# HELP prometheus_rule_evaluation_duration_seconds The duration for a rule to execute.
# TYPE prometheus_rule_evaluation_duration_seconds summary
prometheus_rule_evaluation_duration_seconds{quantile="0.5"} 0.000354314
prometheus_rule_evaluation_duration_seconds{quantile="0.9"} 0.001542137
prometheus_rule_evaluation_duration_seconds{quantile="0.99"} 0.031095665
prometheus_rule_evaluation_duration_seconds_sum 2186.9782992800046
prometheus_rule_evaluation_duration_seconds_count 807190
# HELP prometheus_rule_evaluation_failures_total The total number of rule evaluation failures.
# TYPE prometheus_rule_evaluation_failures_total counter
prometheus_rule_evaluation_failures_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/anko0714airflowbkp-airflow-anko0714airflowbkp-15b95c5c-90e3-4ec6-b92e-08c69e578a6e.yaml;general.kubeRules"} 0
prometheus_rule_evaluation_failures_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/arango-db-a-prometheus-arangodb-service-rules-78cd09de-3e41-4795-a7ad-aba41c42d51d.yaml;general.rules"} 0
prometheus_rule_evaluation_failures_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/arangodb-graph-test-prometheus-arangodb-service-rules-20e72473-b95d-4df3-915b-16a9ff774836.yaml;general.rules"} 0
prometheus_rule_evaluation_failures_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/arangodb-prometheus-arangodb-service-rules-75974246-3cf7-474e-adfd-688c798c25d1.yaml;general.rules"} 0
prometheus_rule_evaluation_failures_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/consul-service-prometheus-consul-rules-7ec3c416-9681-4630-832e-65285b6200cd.yaml;general.rules"} 0
prometheus_rule_evaluation_failures_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/core-mistral-alert-7e644bba-9d46-40b9-9ca1-47189e9aa887.yaml;core-mistral-service"} 0
prometheus_rule_evaluation_failures_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/core-prometheus-cdc-streaming-service-rules-39068d15-ab3e-49c3-85c3-97089f4a831f.yaml;core-cdc-streaming-platform"} 0
prometheus_rule_evaluation_failures_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/core-prometheus-cloud-streaming-service-rules-b6c45b5b-4656-4240-8b08-cf538881cbed.yaml;core-streaming-service-core"} 0
prometheus_rule_evaluation_failures_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/default-prometheus-arangodb-service-rules-7f343f20-16b7-4f90-a284-8e244507be69.yaml;general.rules"} 0
prometheus_rule_evaluation_failures_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/dev-helm-install-mistral-alert-abbb4231-aba7-4ee8-9a40-25abba628086.yaml;dev-helm-install-mistral-operator-dev-helm-install"} 0
prometheus_rule_evaluation_failures_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/dev-mistral-demo-mistral-alert-c30038fd-6be0-45a6-8d1a-402e9bd863d2.yaml;dev-mistral-demo-mistral-service"} 0
prometheus_rule_evaluation_failures_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/dev-mistral-mistral-alert-682c78ec-8d06-43ca-b7cf-936694c0fae2.yaml;dev-mistral-mistral-service"} 0
prometheus_rule_evaluation_failures_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/domain-search-prometheus-opensearch-service-rules-0a349e4b-db79-42f1-8265-0f1e104a3d41.yaml;opensearch.rules"} 0
prometheus_rule_evaluation_failures_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/elasticsearch-cluster-prometheus-elasticsearch-service-rules-ba16a3c6-bd4f-445c-830c-2e0c8e600e70.yaml;elasticsearch.rules"} 0
prometheus_rule_evaluation_failures_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/elasticsearch-service-prometheus-elasticsearch-service-rules-2be23428-fad5-4626-afc0-d5601aca6f01.yaml;ElasticsearchAlerts"} 0
prometheus_rule_evaluation_failures_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/hadoop-monitoring-hadoop-monitoring-prometheus-rules-88f89818-d5ed-4b63-90b5-5e80f0cdf4ee.yaml;general.rules"} 0
prometheus_rule_evaluation_failures_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/hadoop-monitoring-hadoop-monitoring-prometheus-rules2-af3957d3-435b-4f5f-9719-09c7487bd66b.yaml;general.rules"} 0
prometheus_rule_evaluation_failures_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/hadoopmonitoring-hadoop-monitoring-prometheus-rules-4d478506-60d0-49c1-b5f5-35f440f2f495.yaml;hadoopmonitoring-monitor.general.rules"} 0
prometheus_rule_evaluation_failures_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/kafka-service-prometheus-kafka-service-rules-aa79acbc-51be-49ca-ad42-3b59bcc63696.yaml;kafka-service-kafka-service-kafka-service"} 0
prometheus_rule_evaluation_failures_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/mistral-mistral-alert-8dd98d92-e2ec-4c81-b9c8-8ca237a49094.yaml;mistral-mistral-operator"} 0
prometheus_rule_evaluation_failures_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/mongodb-prometheus-mongodb-rules-662a31e8-0419-44e0-9a6e-1cd71ecfe13e.yaml;mongodb-mongodb-operator"} 0
prometheus_rule_evaluation_failures_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-cert-exporter-prometheus-rules-be00bda8-3ce7-4de6-aa54-60de240d847f.yaml;cert-exporter"} 0
prometheus_rule_evaluation_failures_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;AlertManager"} 0
prometheus_rule_evaluation_failures_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;CoreDnsAlerts"} 0
prometheus_rule_evaluation_failures_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;DockerContainers"} 0
prometheus_rule_evaluation_failures_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;Etcd"} 0
prometheus_rule_evaluation_failures_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;KubebernetesAlerts"} 0
prometheus_rule_evaluation_failures_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;NginxIngressAlerts"} 0
prometheus_rule_evaluation_failures_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;NodeExporters"} 0
prometheus_rule_evaluation_failures_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;NodeProcesses"} 0
prometheus_rule_evaluation_failures_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;SelfMonitoring"} 0
prometheus_rule_evaluation_failures_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/opensearch-service-prometheus-opensearch-service-rules-d7a62102-5e19-427b-8cec-9f49c9520a9d.yaml;opensearch.rules"} 0
prometheus_rule_evaluation_failures_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/pg-metr-prometheus-postgres-service-rules-dc530ea7-cf17-4e3d-bf63-ada04b027880.yaml;pg-metr-postgres-operator"} 0
prometheus_rule_evaluation_failures_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/postgres-upgrade-prometheus-clickhouse-service-rules-80ab245c-653e-4fe4-952b-3fc62c8ddee0.yaml;postgres-upgrade-clickhouse-operator"} 0
prometheus_rule_evaluation_failures_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/profiler-cassandra-prometheus-cassandra-rules-3b7ef1e2-2e09-4776-a7c9-3af374492907.yaml;profiler-cassandra-cassandra-operator"} 0
prometheus_rule_evaluation_failures_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/rabbitmq-prometheus-rabbitmq-service-rules-59b077e6-5aa6-49ca-995b-55225bdfa9a6.yaml;general.rules"} 0
prometheus_rule_evaluation_failures_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/redis-prometheus-redis-rules-8aceaa70-6b23-4d0c-be5e-d445627aaa0e.yaml;redis-redis-operator"} 0
prometheus_rule_evaluation_failures_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/spark-operator-gcp-spark-operator-prometheus-rules-28d8b3a4-83b9-457a-9dcd-4b915ff141b0.yaml;spark-operator-gcp-"} 0
prometheus_rule_evaluation_failures_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/streaming-service-prometheus-streaming-service-rules-ca22df5a-e469-4c53-a00e-4ce9d2b8d1b1.yaml;streaming-service-streaming-service-streaming-service"} 0
prometheus_rule_evaluation_failures_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/vault-service-prometheus-vault-service-rules-e526a2e8-d9ed-4748-a535-6bedefae75e8.yaml;vault-service-vault-operator-vault-service"} 0
prometheus_rule_evaluation_failures_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/zookeeper-service-prometheus-zookeeper-service-rules-8646f671-901d-40f4-9ef7-8b7535ce08c3.yaml;zookeeper-service-zookeeper-service-zookeeper-service"} 0
# HELP prometheus_rule_evaluations_total The total number of rule evaluations.
# TYPE prometheus_rule_evaluations_total counter
prometheus_rule_evaluations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/anko0714airflowbkp-airflow-anko0714airflowbkp-15b95c5c-90e3-4ec6-b92e-08c69e578a6e.yaml;general.kubeRules"} 15230
prometheus_rule_evaluations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/arango-db-a-prometheus-arangodb-service-rules-78cd09de-3e41-4795-a7ad-aba41c42d51d.yaml;general.rules"} 10661
prometheus_rule_evaluations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/arangodb-graph-test-prometheus-arangodb-service-rules-20e72473-b95d-4df3-915b-16a9ff774836.yaml;general.rules"} 10661
prometheus_rule_evaluations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/arangodb-prometheus-arangodb-service-rules-75974246-3cf7-474e-adfd-688c798c25d1.yaml;general.rules"} 9138
prometheus_rule_evaluations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/consul-service-prometheus-consul-rules-7ec3c416-9681-4630-832e-65285b6200cd.yaml;general.rules"} 7615
prometheus_rule_evaluations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/core-mistral-alert-7e644bba-9d46-40b9-9ca1-47189e9aa887.yaml;core-mistral-service"} 36552
prometheus_rule_evaluations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/core-prometheus-cdc-streaming-service-rules-39068d15-ab3e-49c3-85c3-97089f4a831f.yaml;core-cdc-streaming-platform"} 27414
prometheus_rule_evaluations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/core-prometheus-cloud-streaming-service-rules-b6c45b5b-4656-4240-8b08-cf538881cbed.yaml;core-streaming-service-core"} 27414
prometheus_rule_evaluations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/default-prometheus-arangodb-service-rules-7f343f20-16b7-4f90-a284-8e244507be69.yaml;general.rules"} 9138
prometheus_rule_evaluations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/dev-helm-install-mistral-alert-abbb4231-aba7-4ee8-9a40-25abba628086.yaml;dev-helm-install-mistral-operator-dev-helm-install"} 36552
prometheus_rule_evaluations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/dev-mistral-demo-mistral-alert-c30038fd-6be0-45a6-8d1a-402e9bd863d2.yaml;dev-mistral-demo-mistral-service"} 36552
prometheus_rule_evaluations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/dev-mistral-mistral-alert-682c78ec-8d06-43ca-b7cf-936694c0fae2.yaml;dev-mistral-mistral-service"} 36552
prometheus_rule_evaluations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/domain-search-prometheus-opensearch-service-rules-0a349e4b-db79-42f1-8265-0f1e104a3d41.yaml;opensearch.rules"} 13707
prometheus_rule_evaluations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/elasticsearch-cluster-prometheus-elasticsearch-service-rules-ba16a3c6-bd4f-445c-830c-2e0c8e600e70.yaml;elasticsearch.rules"} 13707
prometheus_rule_evaluations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/elasticsearch-service-prometheus-elasticsearch-service-rules-2be23428-fad5-4626-afc0-d5601aca6f01.yaml;ElasticsearchAlerts"} 13707
prometheus_rule_evaluations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/hadoop-monitoring-hadoop-monitoring-prometheus-rules-88f89818-d5ed-4b63-90b5-5e80f0cdf4ee.yaml;general.rules"} 44167
prometheus_rule_evaluations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/hadoop-monitoring-hadoop-monitoring-prometheus-rules2-af3957d3-435b-4f5f-9719-09c7487bd66b.yaml;general.rules"} 1523
prometheus_rule_evaluations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/hadoopmonitoring-hadoop-monitoring-prometheus-rules-4d478506-60d0-49c1-b5f5-35f440f2f495.yaml;hadoopmonitoring-monitor.general.rules"} 50259
prometheus_rule_evaluations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/kafka-service-prometheus-kafka-service-rules-aa79acbc-51be-49ca-ad42-3b59bcc63696.yaml;kafka-service-kafka-service-kafka-service"} 10661
prometheus_rule_evaluations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/mistral-mistral-alert-8dd98d92-e2ec-4c81-b9c8-8ca237a49094.yaml;mistral-mistral-operator"} 36552
prometheus_rule_evaluations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/mongodb-prometheus-mongodb-rules-662a31e8-0419-44e0-9a6e-1cd71ecfe13e.yaml;mongodb-mongodb-operator"} 13707
prometheus_rule_evaluations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-cert-exporter-prometheus-rules-be00bda8-3ce7-4de6-aa54-60de240d847f.yaml;cert-exporter"} 4569
prometheus_rule_evaluations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;AlertManager"} 4569
prometheus_rule_evaluations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;CoreDnsAlerts"} 13707
prometheus_rule_evaluations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;DockerContainers"} 6092
prometheus_rule_evaluations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;Etcd"} 15230
prometheus_rule_evaluations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;KubebernetesAlerts"} 45690
prometheus_rule_evaluations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;NginxIngressAlerts"} 4569
prometheus_rule_evaluations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;NodeExporters"} 21322
prometheus_rule_evaluations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;NodeProcesses"} 1523
prometheus_rule_evaluations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;SelfMonitoring"} 33506
prometheus_rule_evaluations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/opensearch-service-prometheus-opensearch-service-rules-d7a62102-5e19-427b-8cec-9f49c9520a9d.yaml;opensearch.rules"} 13707
prometheus_rule_evaluations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/pg-metr-prometheus-postgres-service-rules-dc530ea7-cf17-4e3d-bf63-ada04b027880.yaml;pg-metr-postgres-operator"} 51782
prometheus_rule_evaluations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/postgres-upgrade-prometheus-clickhouse-service-rules-80ab245c-653e-4fe4-952b-3fc62c8ddee0.yaml;postgres-upgrade-clickhouse-operator"} 48736
prometheus_rule_evaluations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/profiler-cassandra-prometheus-cassandra-rules-3b7ef1e2-2e09-4776-a7c9-3af374492907.yaml;profiler-cassandra-cassandra-operator"} 13707
prometheus_rule_evaluations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/rabbitmq-prometheus-rabbitmq-service-rules-59b077e6-5aa6-49ca-995b-55225bdfa9a6.yaml;general.rules"} 7615
prometheus_rule_evaluations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/redis-prometheus-redis-rules-8aceaa70-6b23-4d0c-be5e-d445627aaa0e.yaml;redis-redis-operator"} 7615
prometheus_rule_evaluations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/spark-operator-gcp-spark-operator-prometheus-rules-28d8b3a4-83b9-457a-9dcd-4b915ff141b0.yaml;spark-operator-gcp-"} 6092
prometheus_rule_evaluations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/streaming-service-prometheus-streaming-service-rules-ca22df5a-e469-4c53-a00e-4ce9d2b8d1b1.yaml;streaming-service-streaming-service-streaming-service"} 27414
prometheus_rule_evaluations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/vault-service-prometheus-vault-service-rules-e526a2e8-d9ed-4748-a535-6bedefae75e8.yaml;vault-service-vault-operator-vault-service"} 10661
prometheus_rule_evaluations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/zookeeper-service-prometheus-zookeeper-service-rules-8646f671-901d-40f4-9ef7-8b7535ce08c3.yaml;zookeeper-service-zookeeper-service-zookeeper-service"} 7615
# HELP prometheus_rule_group_duration_seconds The duration of rule group evaluations.
# TYPE prometheus_rule_group_duration_seconds summary
prometheus_rule_group_duration_seconds{quantile="0.01"} 0.000432025
prometheus_rule_group_duration_seconds{quantile="0.05"} 0.001266848
prometheus_rule_group_duration_seconds{quantile="0.5"} 0.007111784
prometheus_rule_group_duration_seconds{quantile="0.9"} 0.034159516
prometheus_rule_group_duration_seconds{quantile="0.99"} 0.905638164
prometheus_rule_group_duration_seconds_sum 2190.8614194229935
prometheus_rule_group_duration_seconds_count 62443
# HELP prometheus_rule_group_interval_seconds The interval of a rule group.
# TYPE prometheus_rule_group_interval_seconds gauge
prometheus_rule_group_interval_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/anko0714airflowbkp-airflow-anko0714airflowbkp-15b95c5c-90e3-4ec6-b92e-08c69e578a6e.yaml;general.kubeRules"} 30
prometheus_rule_group_interval_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/arango-db-a-prometheus-arangodb-service-rules-78cd09de-3e41-4795-a7ad-aba41c42d51d.yaml;general.rules"} 30
prometheus_rule_group_interval_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/arangodb-graph-test-prometheus-arangodb-service-rules-20e72473-b95d-4df3-915b-16a9ff774836.yaml;general.rules"} 30
prometheus_rule_group_interval_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/arangodb-prometheus-arangodb-service-rules-75974246-3cf7-474e-adfd-688c798c25d1.yaml;general.rules"} 30
prometheus_rule_group_interval_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/consul-service-prometheus-consul-rules-7ec3c416-9681-4630-832e-65285b6200cd.yaml;general.rules"} 30
prometheus_rule_group_interval_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/core-mistral-alert-7e644bba-9d46-40b9-9ca1-47189e9aa887.yaml;core-mistral-service"} 30
prometheus_rule_group_interval_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/core-prometheus-cdc-streaming-service-rules-39068d15-ab3e-49c3-85c3-97089f4a831f.yaml;core-cdc-streaming-platform"} 30
prometheus_rule_group_interval_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/core-prometheus-cloud-streaming-service-rules-b6c45b5b-4656-4240-8b08-cf538881cbed.yaml;core-streaming-service-core"} 30
prometheus_rule_group_interval_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/default-prometheus-arangodb-service-rules-7f343f20-16b7-4f90-a284-8e244507be69.yaml;general.rules"} 30
prometheus_rule_group_interval_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/dev-helm-install-mistral-alert-abbb4231-aba7-4ee8-9a40-25abba628086.yaml;dev-helm-install-mistral-operator-dev-helm-install"} 30
prometheus_rule_group_interval_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/dev-mistral-demo-mistral-alert-c30038fd-6be0-45a6-8d1a-402e9bd863d2.yaml;dev-mistral-demo-mistral-service"} 30
prometheus_rule_group_interval_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/dev-mistral-mistral-alert-682c78ec-8d06-43ca-b7cf-936694c0fae2.yaml;dev-mistral-mistral-service"} 30
prometheus_rule_group_interval_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/domain-search-prometheus-opensearch-service-rules-0a349e4b-db79-42f1-8265-0f1e104a3d41.yaml;opensearch.rules"} 30
prometheus_rule_group_interval_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/elasticsearch-cluster-prometheus-elasticsearch-service-rules-ba16a3c6-bd4f-445c-830c-2e0c8e600e70.yaml;elasticsearch.rules"} 30
prometheus_rule_group_interval_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/elasticsearch-service-prometheus-elasticsearch-service-rules-2be23428-fad5-4626-afc0-d5601aca6f01.yaml;ElasticsearchAlerts"} 30
prometheus_rule_group_interval_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/hadoop-monitoring-hadoop-monitoring-prometheus-rules-88f89818-d5ed-4b63-90b5-5e80f0cdf4ee.yaml;general.rules"} 30
prometheus_rule_group_interval_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/hadoop-monitoring-hadoop-monitoring-prometheus-rules2-af3957d3-435b-4f5f-9719-09c7487bd66b.yaml;general.rules"} 30
prometheus_rule_group_interval_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/hadoopmonitoring-hadoop-monitoring-prometheus-rules-4d478506-60d0-49c1-b5f5-35f440f2f495.yaml;hadoopmonitoring-monitor.general.rules"} 30
prometheus_rule_group_interval_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/kafka-service-prometheus-kafka-service-rules-aa79acbc-51be-49ca-ad42-3b59bcc63696.yaml;kafka-service-kafka-service-kafka-service"} 30
prometheus_rule_group_interval_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/mistral-mistral-alert-8dd98d92-e2ec-4c81-b9c8-8ca237a49094.yaml;mistral-mistral-operator"} 30
prometheus_rule_group_interval_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/mongodb-prometheus-mongodb-rules-662a31e8-0419-44e0-9a6e-1cd71ecfe13e.yaml;mongodb-mongodb-operator"} 30
prometheus_rule_group_interval_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-cert-exporter-prometheus-rules-be00bda8-3ce7-4de6-aa54-60de240d847f.yaml;cert-exporter"} 30
prometheus_rule_group_interval_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;AlertManager"} 30
prometheus_rule_group_interval_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;CoreDnsAlerts"} 30
prometheus_rule_group_interval_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;DockerContainers"} 30
prometheus_rule_group_interval_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;Etcd"} 30
prometheus_rule_group_interval_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;KubebernetesAlerts"} 30
prometheus_rule_group_interval_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;NginxIngressAlerts"} 30
prometheus_rule_group_interval_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;NodeExporters"} 30
prometheus_rule_group_interval_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;NodeProcesses"} 30
prometheus_rule_group_interval_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;SelfMonitoring"} 30
prometheus_rule_group_interval_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/opensearch-service-prometheus-opensearch-service-rules-d7a62102-5e19-427b-8cec-9f49c9520a9d.yaml;opensearch.rules"} 30
prometheus_rule_group_interval_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/pg-metr-prometheus-postgres-service-rules-dc530ea7-cf17-4e3d-bf63-ada04b027880.yaml;pg-metr-postgres-operator"} 30
prometheus_rule_group_interval_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/postgres-upgrade-prometheus-clickhouse-service-rules-80ab245c-653e-4fe4-952b-3fc62c8ddee0.yaml;postgres-upgrade-clickhouse-operator"} 30
prometheus_rule_group_interval_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/profiler-cassandra-prometheus-cassandra-rules-3b7ef1e2-2e09-4776-a7c9-3af374492907.yaml;profiler-cassandra-cassandra-operator"} 30
prometheus_rule_group_interval_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/rabbitmq-prometheus-rabbitmq-service-rules-59b077e6-5aa6-49ca-995b-55225bdfa9a6.yaml;general.rules"} 30
prometheus_rule_group_interval_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/redis-prometheus-redis-rules-8aceaa70-6b23-4d0c-be5e-d445627aaa0e.yaml;redis-redis-operator"} 30
prometheus_rule_group_interval_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/spark-operator-gcp-spark-operator-prometheus-rules-28d8b3a4-83b9-457a-9dcd-4b915ff141b0.yaml;spark-operator-gcp-"} 30
prometheus_rule_group_interval_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/streaming-service-prometheus-streaming-service-rules-ca22df5a-e469-4c53-a00e-4ce9d2b8d1b1.yaml;streaming-service-streaming-service-streaming-service"} 30
prometheus_rule_group_interval_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/vault-service-prometheus-vault-service-rules-e526a2e8-d9ed-4748-a535-6bedefae75e8.yaml;vault-service-vault-operator-vault-service"} 30
prometheus_rule_group_interval_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/zookeeper-service-prometheus-zookeeper-service-rules-8646f671-901d-40f4-9ef7-8b7535ce08c3.yaml;zookeeper-service-zookeeper-service-zookeeper-service"} 30
# HELP prometheus_rule_group_iterations_missed_total The total number of rule group evaluations missed due to slow rule group evaluation.
# TYPE prometheus_rule_group_iterations_missed_total counter
prometheus_rule_group_iterations_missed_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/anko0714airflowbkp-airflow-anko0714airflowbkp-15b95c5c-90e3-4ec6-b92e-08c69e578a6e.yaml;general.kubeRules"} 0
prometheus_rule_group_iterations_missed_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/arango-db-a-prometheus-arangodb-service-rules-78cd09de-3e41-4795-a7ad-aba41c42d51d.yaml;general.rules"} 0
prometheus_rule_group_iterations_missed_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/arangodb-graph-test-prometheus-arangodb-service-rules-20e72473-b95d-4df3-915b-16a9ff774836.yaml;general.rules"} 0
prometheus_rule_group_iterations_missed_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/arangodb-prometheus-arangodb-service-rules-75974246-3cf7-474e-adfd-688c798c25d1.yaml;general.rules"} 0
prometheus_rule_group_iterations_missed_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/consul-service-prometheus-consul-rules-7ec3c416-9681-4630-832e-65285b6200cd.yaml;general.rules"} 0
prometheus_rule_group_iterations_missed_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/core-mistral-alert-7e644bba-9d46-40b9-9ca1-47189e9aa887.yaml;core-mistral-service"} 0
prometheus_rule_group_iterations_missed_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/core-prometheus-cdc-streaming-service-rules-39068d15-ab3e-49c3-85c3-97089f4a831f.yaml;core-cdc-streaming-platform"} 0
prometheus_rule_group_iterations_missed_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/core-prometheus-cloud-streaming-service-rules-b6c45b5b-4656-4240-8b08-cf538881cbed.yaml;core-streaming-service-core"} 0
prometheus_rule_group_iterations_missed_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/default-prometheus-arangodb-service-rules-7f343f20-16b7-4f90-a284-8e244507be69.yaml;general.rules"} 0
prometheus_rule_group_iterations_missed_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/dev-helm-install-mistral-alert-abbb4231-aba7-4ee8-9a40-25abba628086.yaml;dev-helm-install-mistral-operator-dev-helm-install"} 0
prometheus_rule_group_iterations_missed_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/dev-mistral-demo-mistral-alert-c30038fd-6be0-45a6-8d1a-402e9bd863d2.yaml;dev-mistral-demo-mistral-service"} 0
prometheus_rule_group_iterations_missed_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/dev-mistral-mistral-alert-682c78ec-8d06-43ca-b7cf-936694c0fae2.yaml;dev-mistral-mistral-service"} 0
prometheus_rule_group_iterations_missed_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/domain-search-prometheus-opensearch-service-rules-0a349e4b-db79-42f1-8265-0f1e104a3d41.yaml;opensearch.rules"} 0
prometheus_rule_group_iterations_missed_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/elasticsearch-cluster-prometheus-elasticsearch-service-rules-ba16a3c6-bd4f-445c-830c-2e0c8e600e70.yaml;elasticsearch.rules"} 0
prometheus_rule_group_iterations_missed_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/elasticsearch-service-prometheus-elasticsearch-service-rules-2be23428-fad5-4626-afc0-d5601aca6f01.yaml;ElasticsearchAlerts"} 0
prometheus_rule_group_iterations_missed_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/hadoop-monitoring-hadoop-monitoring-prometheus-rules-88f89818-d5ed-4b63-90b5-5e80f0cdf4ee.yaml;general.rules"} 0
prometheus_rule_group_iterations_missed_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/hadoop-monitoring-hadoop-monitoring-prometheus-rules2-af3957d3-435b-4f5f-9719-09c7487bd66b.yaml;general.rules"} 0
prometheus_rule_group_iterations_missed_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/hadoopmonitoring-hadoop-monitoring-prometheus-rules-4d478506-60d0-49c1-b5f5-35f440f2f495.yaml;hadoopmonitoring-monitor.general.rules"} 0
prometheus_rule_group_iterations_missed_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/kafka-service-prometheus-kafka-service-rules-aa79acbc-51be-49ca-ad42-3b59bcc63696.yaml;kafka-service-kafka-service-kafka-service"} 0
prometheus_rule_group_iterations_missed_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/mistral-mistral-alert-8dd98d92-e2ec-4c81-b9c8-8ca237a49094.yaml;mistral-mistral-operator"} 0
prometheus_rule_group_iterations_missed_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/mongodb-prometheus-mongodb-rules-662a31e8-0419-44e0-9a6e-1cd71ecfe13e.yaml;mongodb-mongodb-operator"} 0
prometheus_rule_group_iterations_missed_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-cert-exporter-prometheus-rules-be00bda8-3ce7-4de6-aa54-60de240d847f.yaml;cert-exporter"} 0
prometheus_rule_group_iterations_missed_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;AlertManager"} 0
prometheus_rule_group_iterations_missed_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;CoreDnsAlerts"} 0
prometheus_rule_group_iterations_missed_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;DockerContainers"} 0
prometheus_rule_group_iterations_missed_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;Etcd"} 0
prometheus_rule_group_iterations_missed_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;KubebernetesAlerts"} 0
prometheus_rule_group_iterations_missed_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;NginxIngressAlerts"} 0
prometheus_rule_group_iterations_missed_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;NodeExporters"} 0
prometheus_rule_group_iterations_missed_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;NodeProcesses"} 0
prometheus_rule_group_iterations_missed_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;SelfMonitoring"} 0
prometheus_rule_group_iterations_missed_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/opensearch-service-prometheus-opensearch-service-rules-d7a62102-5e19-427b-8cec-9f49c9520a9d.yaml;opensearch.rules"} 0
prometheus_rule_group_iterations_missed_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/pg-metr-prometheus-postgres-service-rules-dc530ea7-cf17-4e3d-bf63-ada04b027880.yaml;pg-metr-postgres-operator"} 0
prometheus_rule_group_iterations_missed_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/postgres-upgrade-prometheus-clickhouse-service-rules-80ab245c-653e-4fe4-952b-3fc62c8ddee0.yaml;postgres-upgrade-clickhouse-operator"} 0
prometheus_rule_group_iterations_missed_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/profiler-cassandra-prometheus-cassandra-rules-3b7ef1e2-2e09-4776-a7c9-3af374492907.yaml;profiler-cassandra-cassandra-operator"} 0
prometheus_rule_group_iterations_missed_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/rabbitmq-prometheus-rabbitmq-service-rules-59b077e6-5aa6-49ca-995b-55225bdfa9a6.yaml;general.rules"} 0
prometheus_rule_group_iterations_missed_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/redis-prometheus-redis-rules-8aceaa70-6b23-4d0c-be5e-d445627aaa0e.yaml;redis-redis-operator"} 0
prometheus_rule_group_iterations_missed_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/spark-operator-gcp-spark-operator-prometheus-rules-28d8b3a4-83b9-457a-9dcd-4b915ff141b0.yaml;spark-operator-gcp-"} 0
prometheus_rule_group_iterations_missed_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/streaming-service-prometheus-streaming-service-rules-ca22df5a-e469-4c53-a00e-4ce9d2b8d1b1.yaml;streaming-service-streaming-service-streaming-service"} 0
prometheus_rule_group_iterations_missed_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/vault-service-prometheus-vault-service-rules-e526a2e8-d9ed-4748-a535-6bedefae75e8.yaml;vault-service-vault-operator-vault-service"} 0
prometheus_rule_group_iterations_missed_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/zookeeper-service-prometheus-zookeeper-service-rules-8646f671-901d-40f4-9ef7-8b7535ce08c3.yaml;zookeeper-service-zookeeper-service-zookeeper-service"} 0
# HELP prometheus_rule_group_iterations_total The total number of scheduled rule group evaluations, whether executed or missed.
# TYPE prometheus_rule_group_iterations_total counter
prometheus_rule_group_iterations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/anko0714airflowbkp-airflow-anko0714airflowbkp-15b95c5c-90e3-4ec6-b92e-08c69e578a6e.yaml;general.kubeRules"} 1523
prometheus_rule_group_iterations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/arango-db-a-prometheus-arangodb-service-rules-78cd09de-3e41-4795-a7ad-aba41c42d51d.yaml;general.rules"} 1523
prometheus_rule_group_iterations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/arangodb-graph-test-prometheus-arangodb-service-rules-20e72473-b95d-4df3-915b-16a9ff774836.yaml;general.rules"} 1523
prometheus_rule_group_iterations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/arangodb-prometheus-arangodb-service-rules-75974246-3cf7-474e-adfd-688c798c25d1.yaml;general.rules"} 1523
prometheus_rule_group_iterations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/consul-service-prometheus-consul-rules-7ec3c416-9681-4630-832e-65285b6200cd.yaml;general.rules"} 1523
prometheus_rule_group_iterations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/core-mistral-alert-7e644bba-9d46-40b9-9ca1-47189e9aa887.yaml;core-mistral-service"} 1523
prometheus_rule_group_iterations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/core-prometheus-cdc-streaming-service-rules-39068d15-ab3e-49c3-85c3-97089f4a831f.yaml;core-cdc-streaming-platform"} 1523
prometheus_rule_group_iterations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/core-prometheus-cloud-streaming-service-rules-b6c45b5b-4656-4240-8b08-cf538881cbed.yaml;core-streaming-service-core"} 1523
prometheus_rule_group_iterations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/default-prometheus-arangodb-service-rules-7f343f20-16b7-4f90-a284-8e244507be69.yaml;general.rules"} 1523
prometheus_rule_group_iterations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/dev-helm-install-mistral-alert-abbb4231-aba7-4ee8-9a40-25abba628086.yaml;dev-helm-install-mistral-operator-dev-helm-install"} 1523
prometheus_rule_group_iterations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/dev-mistral-demo-mistral-alert-c30038fd-6be0-45a6-8d1a-402e9bd863d2.yaml;dev-mistral-demo-mistral-service"} 1523
prometheus_rule_group_iterations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/dev-mistral-mistral-alert-682c78ec-8d06-43ca-b7cf-936694c0fae2.yaml;dev-mistral-mistral-service"} 1523
prometheus_rule_group_iterations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/domain-search-prometheus-opensearch-service-rules-0a349e4b-db79-42f1-8265-0f1e104a3d41.yaml;opensearch.rules"} 1523
prometheus_rule_group_iterations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/elasticsearch-cluster-prometheus-elasticsearch-service-rules-ba16a3c6-bd4f-445c-830c-2e0c8e600e70.yaml;elasticsearch.rules"} 1523
prometheus_rule_group_iterations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/elasticsearch-service-prometheus-elasticsearch-service-rules-2be23428-fad5-4626-afc0-d5601aca6f01.yaml;ElasticsearchAlerts"} 1523
prometheus_rule_group_iterations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/hadoop-monitoring-hadoop-monitoring-prometheus-rules-88f89818-d5ed-4b63-90b5-5e80f0cdf4ee.yaml;general.rules"} 1523
prometheus_rule_group_iterations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/hadoop-monitoring-hadoop-monitoring-prometheus-rules2-af3957d3-435b-4f5f-9719-09c7487bd66b.yaml;general.rules"} 1523
prometheus_rule_group_iterations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/hadoopmonitoring-hadoop-monitoring-prometheus-rules-4d478506-60d0-49c1-b5f5-35f440f2f495.yaml;hadoopmonitoring-monitor.general.rules"} 1523
prometheus_rule_group_iterations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/kafka-service-prometheus-kafka-service-rules-aa79acbc-51be-49ca-ad42-3b59bcc63696.yaml;kafka-service-kafka-service-kafka-service"} 1523
prometheus_rule_group_iterations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/mistral-mistral-alert-8dd98d92-e2ec-4c81-b9c8-8ca237a49094.yaml;mistral-mistral-operator"} 1523
prometheus_rule_group_iterations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/mongodb-prometheus-mongodb-rules-662a31e8-0419-44e0-9a6e-1cd71ecfe13e.yaml;mongodb-mongodb-operator"} 1523
prometheus_rule_group_iterations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-cert-exporter-prometheus-rules-be00bda8-3ce7-4de6-aa54-60de240d847f.yaml;cert-exporter"} 1523
prometheus_rule_group_iterations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;AlertManager"} 1523
prometheus_rule_group_iterations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;CoreDnsAlerts"} 1523
prometheus_rule_group_iterations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;DockerContainers"} 1523
prometheus_rule_group_iterations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;Etcd"} 1523
prometheus_rule_group_iterations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;KubebernetesAlerts"} 1523
prometheus_rule_group_iterations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;NginxIngressAlerts"} 1523
prometheus_rule_group_iterations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;NodeExporters"} 1523
prometheus_rule_group_iterations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;NodeProcesses"} 1523
prometheus_rule_group_iterations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;SelfMonitoring"} 1523
prometheus_rule_group_iterations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/opensearch-service-prometheus-opensearch-service-rules-d7a62102-5e19-427b-8cec-9f49c9520a9d.yaml;opensearch.rules"} 1523
prometheus_rule_group_iterations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/pg-metr-prometheus-postgres-service-rules-dc530ea7-cf17-4e3d-bf63-ada04b027880.yaml;pg-metr-postgres-operator"} 1523
prometheus_rule_group_iterations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/postgres-upgrade-prometheus-clickhouse-service-rules-80ab245c-653e-4fe4-952b-3fc62c8ddee0.yaml;postgres-upgrade-clickhouse-operator"} 1523
prometheus_rule_group_iterations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/profiler-cassandra-prometheus-cassandra-rules-3b7ef1e2-2e09-4776-a7c9-3af374492907.yaml;profiler-cassandra-cassandra-operator"} 1523
prometheus_rule_group_iterations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/rabbitmq-prometheus-rabbitmq-service-rules-59b077e6-5aa6-49ca-995b-55225bdfa9a6.yaml;general.rules"} 1523
prometheus_rule_group_iterations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/redis-prometheus-redis-rules-8aceaa70-6b23-4d0c-be5e-d445627aaa0e.yaml;redis-redis-operator"} 1523
prometheus_rule_group_iterations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/spark-operator-gcp-spark-operator-prometheus-rules-28d8b3a4-83b9-457a-9dcd-4b915ff141b0.yaml;spark-operator-gcp-"} 1523
prometheus_rule_group_iterations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/streaming-service-prometheus-streaming-service-rules-ca22df5a-e469-4c53-a00e-4ce9d2b8d1b1.yaml;streaming-service-streaming-service-streaming-service"} 1523
prometheus_rule_group_iterations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/vault-service-prometheus-vault-service-rules-e526a2e8-d9ed-4748-a535-6bedefae75e8.yaml;vault-service-vault-operator-vault-service"} 1523
prometheus_rule_group_iterations_total{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/zookeeper-service-prometheus-zookeeper-service-rules-8646f671-901d-40f4-9ef7-8b7535ce08c3.yaml;zookeeper-service-zookeeper-service-zookeeper-service"} 1523
# HELP prometheus_rule_group_last_duration_seconds The duration of the last rule group evaluation.
# TYPE prometheus_rule_group_last_duration_seconds gauge
prometheus_rule_group_last_duration_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/anko0714airflowbkp-airflow-anko0714airflowbkp-15b95c5c-90e3-4ec6-b92e-08c69e578a6e.yaml;general.kubeRules"} 0.005813925
prometheus_rule_group_last_duration_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/arango-db-a-prometheus-arangodb-service-rules-78cd09de-3e41-4795-a7ad-aba41c42d51d.yaml;general.rules"} 0.000908235
prometheus_rule_group_last_duration_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/arangodb-graph-test-prometheus-arangodb-service-rules-20e72473-b95d-4df3-915b-16a9ff774836.yaml;general.rules"} 0.050567025
prometheus_rule_group_last_duration_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/arangodb-prometheus-arangodb-service-rules-75974246-3cf7-474e-adfd-688c798c25d1.yaml;general.rules"} 0.001391681
prometheus_rule_group_last_duration_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/consul-service-prometheus-consul-rules-7ec3c416-9681-4630-832e-65285b6200cd.yaml;general.rules"} 0.003557188
prometheus_rule_group_last_duration_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/core-mistral-alert-7e644bba-9d46-40b9-9ca1-47189e9aa887.yaml;core-mistral-service"} 0.007384662
prometheus_rule_group_last_duration_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/core-prometheus-cdc-streaming-service-rules-39068d15-ab3e-49c3-85c3-97089f4a831f.yaml;core-cdc-streaming-platform"} 0.009362665
prometheus_rule_group_last_duration_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/core-prometheus-cloud-streaming-service-rules-b6c45b5b-4656-4240-8b08-cf538881cbed.yaml;core-streaming-service-core"} 0.014248832
prometheus_rule_group_last_duration_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/default-prometheus-arangodb-service-rules-7f343f20-16b7-4f90-a284-8e244507be69.yaml;general.rules"} 0.004207313
prometheus_rule_group_last_duration_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/dev-helm-install-mistral-alert-abbb4231-aba7-4ee8-9a40-25abba628086.yaml;dev-helm-install-mistral-operator-dev-helm-install"} 0.013613844
prometheus_rule_group_last_duration_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/dev-mistral-demo-mistral-alert-c30038fd-6be0-45a6-8d1a-402e9bd863d2.yaml;dev-mistral-demo-mistral-service"} 0.005658234
prometheus_rule_group_last_duration_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/dev-mistral-mistral-alert-682c78ec-8d06-43ca-b7cf-936694c0fae2.yaml;dev-mistral-mistral-service"} 0.013111994
prometheus_rule_group_last_duration_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/domain-search-prometheus-opensearch-service-rules-0a349e4b-db79-42f1-8265-0f1e104a3d41.yaml;opensearch.rules"} 0.002380437
prometheus_rule_group_last_duration_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/elasticsearch-cluster-prometheus-elasticsearch-service-rules-ba16a3c6-bd4f-445c-830c-2e0c8e600e70.yaml;elasticsearch.rules"} 0.002503483
prometheus_rule_group_last_duration_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/elasticsearch-service-prometheus-elasticsearch-service-rules-2be23428-fad5-4626-afc0-d5601aca6f01.yaml;ElasticsearchAlerts"} 0.00167689
prometheus_rule_group_last_duration_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/hadoop-monitoring-hadoop-monitoring-prometheus-rules-88f89818-d5ed-4b63-90b5-5e80f0cdf4ee.yaml;general.rules"} 0.010886867
prometheus_rule_group_last_duration_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/hadoop-monitoring-hadoop-monitoring-prometheus-rules2-af3957d3-435b-4f5f-9719-09c7487bd66b.yaml;general.rules"} 0.000447767
prometheus_rule_group_last_duration_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/hadoopmonitoring-hadoop-monitoring-prometheus-rules-4d478506-60d0-49c1-b5f5-35f440f2f495.yaml;hadoopmonitoring-monitor.general.rules"} 0.005778229
prometheus_rule_group_last_duration_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/kafka-service-prometheus-kafka-service-rules-aa79acbc-51be-49ca-ad42-3b59bcc63696.yaml;kafka-service-kafka-service-kafka-service"} 0.005417022
prometheus_rule_group_last_duration_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/mistral-mistral-alert-8dd98d92-e2ec-4c81-b9c8-8ca237a49094.yaml;mistral-mistral-operator"} 0.010154882
prometheus_rule_group_last_duration_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/mongodb-prometheus-mongodb-rules-662a31e8-0419-44e0-9a6e-1cd71ecfe13e.yaml;mongodb-mongodb-operator"} 0.006847614
prometheus_rule_group_last_duration_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-cert-exporter-prometheus-rules-be00bda8-3ce7-4de6-aa54-60de240d847f.yaml;cert-exporter"} 0.011564945
prometheus_rule_group_last_duration_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;AlertManager"} 0.001574915
prometheus_rule_group_last_duration_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;CoreDnsAlerts"} 0.006278217
prometheus_rule_group_last_duration_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;DockerContainers"} 0.024314154
prometheus_rule_group_last_duration_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;Etcd"} 0.46934642
prometheus_rule_group_last_duration_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;KubebernetesAlerts"} 1.670591264
prometheus_rule_group_last_duration_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;NginxIngressAlerts"} 0.091868264
prometheus_rule_group_last_duration_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;NodeExporters"} 0.053536307
prometheus_rule_group_last_duration_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;NodeProcesses"} 0.017646139
prometheus_rule_group_last_duration_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;SelfMonitoring"} 0.013488786
prometheus_rule_group_last_duration_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/opensearch-service-prometheus-opensearch-service-rules-d7a62102-5e19-427b-8cec-9f49c9520a9d.yaml;opensearch.rules"} 0.003121454
prometheus_rule_group_last_duration_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/pg-metr-prometheus-postgres-service-rules-dc530ea7-cf17-4e3d-bf63-ada04b027880.yaml;pg-metr-postgres-operator"} 0.011012717
prometheus_rule_group_last_duration_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/postgres-upgrade-prometheus-clickhouse-service-rules-80ab245c-653e-4fe4-952b-3fc62c8ddee0.yaml;postgres-upgrade-clickhouse-operator"} 0.016774905
prometheus_rule_group_last_duration_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/profiler-cassandra-prometheus-cassandra-rules-3b7ef1e2-2e09-4776-a7c9-3af374492907.yaml;profiler-cassandra-cassandra-operator"} 0.015382333
prometheus_rule_group_last_duration_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/rabbitmq-prometheus-rabbitmq-service-rules-59b077e6-5aa6-49ca-995b-55225bdfa9a6.yaml;general.rules"} 0.003930287
prometheus_rule_group_last_duration_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/redis-prometheus-redis-rules-8aceaa70-6b23-4d0c-be5e-d445627aaa0e.yaml;redis-redis-operator"} 0.006462573
prometheus_rule_group_last_duration_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/spark-operator-gcp-spark-operator-prometheus-rules-28d8b3a4-83b9-457a-9dcd-4b915ff141b0.yaml;spark-operator-gcp-"} 0.003137044
prometheus_rule_group_last_duration_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/streaming-service-prometheus-streaming-service-rules-ca22df5a-e469-4c53-a00e-4ce9d2b8d1b1.yaml;streaming-service-streaming-service-streaming-service"} 0.025970085
prometheus_rule_group_last_duration_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/vault-service-prometheus-vault-service-rules-e526a2e8-d9ed-4748-a535-6bedefae75e8.yaml;vault-service-vault-operator-vault-service"} 0.006429117
prometheus_rule_group_last_duration_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/zookeeper-service-prometheus-zookeeper-service-rules-8646f671-901d-40f4-9ef7-8b7535ce08c3.yaml;zookeeper-service-zookeeper-service-zookeeper-service"} 0.006975709
# HELP prometheus_rule_group_last_evaluation_samples The number of samples returned during the last rule group evaluation.
# TYPE prometheus_rule_group_last_evaluation_samples gauge
prometheus_rule_group_last_evaluation_samples{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/anko0714airflowbkp-airflow-anko0714airflowbkp-15b95c5c-90e3-4ec6-b92e-08c69e578a6e.yaml;general.kubeRules"} 6
prometheus_rule_group_last_evaluation_samples{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/arango-db-a-prometheus-arangodb-service-rules-78cd09de-3e41-4795-a7ad-aba41c42d51d.yaml;general.rules"} 0
prometheus_rule_group_last_evaluation_samples{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/arangodb-graph-test-prometheus-arangodb-service-rules-20e72473-b95d-4df3-915b-16a9ff774836.yaml;general.rules"} 0
prometheus_rule_group_last_evaluation_samples{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/arangodb-prometheus-arangodb-service-rules-75974246-3cf7-474e-adfd-688c798c25d1.yaml;general.rules"} 0
prometheus_rule_group_last_evaluation_samples{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/consul-service-prometheus-consul-rules-7ec3c416-9681-4630-832e-65285b6200cd.yaml;general.rules"} 0
prometheus_rule_group_last_evaluation_samples{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/core-mistral-alert-7e644bba-9d46-40b9-9ca1-47189e9aa887.yaml;core-mistral-service"} 0
prometheus_rule_group_last_evaluation_samples{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/core-prometheus-cdc-streaming-service-rules-39068d15-ab3e-49c3-85c3-97089f4a831f.yaml;core-cdc-streaming-platform"} 0
prometheus_rule_group_last_evaluation_samples{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/core-prometheus-cloud-streaming-service-rules-b6c45b5b-4656-4240-8b08-cf538881cbed.yaml;core-streaming-service-core"} 4
prometheus_rule_group_last_evaluation_samples{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/default-prometheus-arangodb-service-rules-7f343f20-16b7-4f90-a284-8e244507be69.yaml;general.rules"} 4
prometheus_rule_group_last_evaluation_samples{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/dev-helm-install-mistral-alert-abbb4231-aba7-4ee8-9a40-25abba628086.yaml;dev-helm-install-mistral-operator-dev-helm-install"} 0
prometheus_rule_group_last_evaluation_samples{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/dev-mistral-demo-mistral-alert-c30038fd-6be0-45a6-8d1a-402e9bd863d2.yaml;dev-mistral-demo-mistral-service"} 0
prometheus_rule_group_last_evaluation_samples{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/dev-mistral-mistral-alert-682c78ec-8d06-43ca-b7cf-936694c0fae2.yaml;dev-mistral-mistral-service"} 0
prometheus_rule_group_last_evaluation_samples{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/domain-search-prometheus-opensearch-service-rules-0a349e4b-db79-42f1-8265-0f1e104a3d41.yaml;opensearch.rules"} 2
prometheus_rule_group_last_evaluation_samples{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/elasticsearch-cluster-prometheus-elasticsearch-service-rules-ba16a3c6-bd4f-445c-830c-2e0c8e600e70.yaml;elasticsearch.rules"} 2
prometheus_rule_group_last_evaluation_samples{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/elasticsearch-service-prometheus-elasticsearch-service-rules-2be23428-fad5-4626-afc0-d5601aca6f01.yaml;ElasticsearchAlerts"} 0
prometheus_rule_group_last_evaluation_samples{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/hadoop-monitoring-hadoop-monitoring-prometheus-rules-88f89818-d5ed-4b63-90b5-5e80f0cdf4ee.yaml;general.rules"} 4
prometheus_rule_group_last_evaluation_samples{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/hadoop-monitoring-hadoop-monitoring-prometheus-rules2-af3957d3-435b-4f5f-9719-09c7487bd66b.yaml;general.rules"} 0
prometheus_rule_group_last_evaluation_samples{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/hadoopmonitoring-hadoop-monitoring-prometheus-rules-4d478506-60d0-49c1-b5f5-35f440f2f495.yaml;hadoopmonitoring-monitor.general.rules"} 4
prometheus_rule_group_last_evaluation_samples{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/kafka-service-prometheus-kafka-service-rules-aa79acbc-51be-49ca-ad42-3b59bcc63696.yaml;kafka-service-kafka-service-kafka-service"} 0
prometheus_rule_group_last_evaluation_samples{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/mistral-mistral-alert-8dd98d92-e2ec-4c81-b9c8-8ca237a49094.yaml;mistral-mistral-operator"} 0
prometheus_rule_group_last_evaluation_samples{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/mongodb-prometheus-mongodb-rules-662a31e8-0419-44e0-9a6e-1cd71ecfe13e.yaml;mongodb-mongodb-operator"} 0
prometheus_rule_group_last_evaluation_samples{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-cert-exporter-prometheus-rules-be00bda8-3ce7-4de6-aa54-60de240d847f.yaml;cert-exporter"} 0
prometheus_rule_group_last_evaluation_samples{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;AlertManager"} 0
prometheus_rule_group_last_evaluation_samples{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;CoreDnsAlerts"} 8
prometheus_rule_group_last_evaluation_samples{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;DockerContainers"} 28
prometheus_rule_group_last_evaluation_samples{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;Etcd"} 0
prometheus_rule_group_last_evaluation_samples{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;KubebernetesAlerts"} 62
prometheus_rule_group_last_evaluation_samples{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;NginxIngressAlerts"} 0
prometheus_rule_group_last_evaluation_samples{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;NodeExporters"} 6
prometheus_rule_group_last_evaluation_samples{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;NodeProcesses"} 0
prometheus_rule_group_last_evaluation_samples{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;SelfMonitoring"} 44
prometheus_rule_group_last_evaluation_samples{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/opensearch-service-prometheus-opensearch-service-rules-d7a62102-5e19-427b-8cec-9f49c9520a9d.yaml;opensearch.rules"} 0
prometheus_rule_group_last_evaluation_samples{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/pg-metr-prometheus-postgres-service-rules-dc530ea7-cf17-4e3d-bf63-ada04b027880.yaml;pg-metr-postgres-operator"} 2
prometheus_rule_group_last_evaluation_samples{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/postgres-upgrade-prometheus-clickhouse-service-rules-80ab245c-653e-4fe4-952b-3fc62c8ddee0.yaml;postgres-upgrade-clickhouse-operator"} 42
prometheus_rule_group_last_evaluation_samples{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/profiler-cassandra-prometheus-cassandra-rules-3b7ef1e2-2e09-4776-a7c9-3af374492907.yaml;profiler-cassandra-cassandra-operator"} 0
prometheus_rule_group_last_evaluation_samples{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/rabbitmq-prometheus-rabbitmq-service-rules-59b077e6-5aa6-49ca-995b-55225bdfa9a6.yaml;general.rules"} 0
prometheus_rule_group_last_evaluation_samples{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/redis-prometheus-redis-rules-8aceaa70-6b23-4d0c-be5e-d445627aaa0e.yaml;redis-redis-operator"} 2
prometheus_rule_group_last_evaluation_samples{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/spark-operator-gcp-spark-operator-prometheus-rules-28d8b3a4-83b9-457a-9dcd-4b915ff141b0.yaml;spark-operator-gcp-"} 0
prometheus_rule_group_last_evaluation_samples{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/streaming-service-prometheus-streaming-service-rules-ca22df5a-e469-4c53-a00e-4ce9d2b8d1b1.yaml;streaming-service-streaming-service-streaming-service"} 4
prometheus_rule_group_last_evaluation_samples{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/vault-service-prometheus-vault-service-rules-e526a2e8-d9ed-4748-a535-6bedefae75e8.yaml;vault-service-vault-operator-vault-service"} 0
prometheus_rule_group_last_evaluation_samples{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/zookeeper-service-prometheus-zookeeper-service-rules-8646f671-901d-40f4-9ef7-8b7535ce08c3.yaml;zookeeper-service-zookeeper-service-zookeeper-service"} 0
# HELP prometheus_rule_group_last_evaluation_timestamp_seconds The timestamp of the last rule group evaluation in seconds.
# TYPE prometheus_rule_group_last_evaluation_timestamp_seconds gauge
prometheus_rule_group_last_evaluation_timestamp_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/anko0714airflowbkp-airflow-anko0714airflowbkp-15b95c5c-90e3-4ec6-b92e-08c69e578a6e.yaml;general.kubeRules"} 1.6726472151840675e+09
prometheus_rule_group_last_evaluation_timestamp_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/arango-db-a-prometheus-arangodb-service-rules-78cd09de-3e41-4795-a7ad-aba41c42d51d.yaml;general.rules"} 1.6726472311602674e+09
prometheus_rule_group_last_evaluation_timestamp_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/arangodb-graph-test-prometheus-arangodb-service-rules-20e72473-b95d-4df3-915b-16a9ff774836.yaml;general.rules"} 1.6726472220745614e+09
prometheus_rule_group_last_evaluation_timestamp_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/arangodb-prometheus-arangodb-service-rules-75974246-3cf7-474e-adfd-688c798c25d1.yaml;general.rules"} 1.6726472115420032e+09
prometheus_rule_group_last_evaluation_timestamp_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/consul-service-prometheus-consul-rules-7ec3c416-9681-4630-832e-65285b6200cd.yaml;general.rules"} 1.6726472348605375e+09
prometheus_rule_group_last_evaluation_timestamp_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/core-mistral-alert-7e644bba-9d46-40b9-9ca1-47189e9aa887.yaml;core-mistral-service"} 1.672647223757945e+09
prometheus_rule_group_last_evaluation_timestamp_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/core-prometheus-cdc-streaming-service-rules-39068d15-ab3e-49c3-85c3-97089f4a831f.yaml;core-cdc-streaming-platform"} 1.6726472314125328e+09
prometheus_rule_group_last_evaluation_timestamp_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/core-prometheus-cloud-streaming-service-rules-b6c45b5b-4656-4240-8b08-cf538881cbed.yaml;core-streaming-service-core"} 1.6726472360350692e+09
prometheus_rule_group_last_evaluation_timestamp_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/default-prometheus-arangodb-service-rules-7f343f20-16b7-4f90-a284-8e244507be69.yaml;general.rules"} 1.6726472351953046e+09
prometheus_rule_group_last_evaluation_timestamp_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/dev-helm-install-mistral-alert-abbb4231-aba7-4ee8-9a40-25abba628086.yaml;dev-helm-install-mistral-operator-dev-helm-install"} 1.6726472227152605e+09
prometheus_rule_group_last_evaluation_timestamp_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/dev-mistral-demo-mistral-alert-c30038fd-6be0-45a6-8d1a-402e9bd863d2.yaml;dev-mistral-demo-mistral-service"} 1.672647219527874e+09
prometheus_rule_group_last_evaluation_timestamp_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/dev-mistral-mistral-alert-682c78ec-8d06-43ca-b7cf-936694c0fae2.yaml;dev-mistral-mistral-service"} 1.6726472327427883e+09
prometheus_rule_group_last_evaluation_timestamp_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/domain-search-prometheus-opensearch-service-rules-0a349e4b-db79-42f1-8265-0f1e104a3d41.yaml;opensearch.rules"} 1.672647225965844e+09
prometheus_rule_group_last_evaluation_timestamp_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/elasticsearch-cluster-prometheus-elasticsearch-service-rules-ba16a3c6-bd4f-445c-830c-2e0c8e600e70.yaml;elasticsearch.rules"} 1.672647227503709e+09
prometheus_rule_group_last_evaluation_timestamp_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/elasticsearch-service-prometheus-elasticsearch-service-rules-2be23428-fad5-4626-afc0-d5601aca6f01.yaml;ElasticsearchAlerts"} 1.6726472122678807e+09
prometheus_rule_group_last_evaluation_timestamp_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/hadoop-monitoring-hadoop-monitoring-prometheus-rules-88f89818-d5ed-4b63-90b5-5e80f0cdf4ee.yaml;general.rules"} 1.6726472160430403e+09
prometheus_rule_group_last_evaluation_timestamp_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/hadoop-monitoring-hadoop-monitoring-prometheus-rules2-af3957d3-435b-4f5f-9719-09c7487bd66b.yaml;general.rules"} 1.672647223427377e+09
prometheus_rule_group_last_evaluation_timestamp_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/hadoopmonitoring-hadoop-monitoring-prometheus-rules-4d478506-60d0-49c1-b5f5-35f440f2f495.yaml;hadoopmonitoring-monitor.general.rules"} 1.672647231622138e+09
prometheus_rule_group_last_evaluation_timestamp_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/kafka-service-prometheus-kafka-service-rules-aa79acbc-51be-49ca-ad42-3b59bcc63696.yaml;kafka-service-kafka-service-kafka-service"} 1.6726472143960621e+09
prometheus_rule_group_last_evaluation_timestamp_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/mistral-mistral-alert-8dd98d92-e2ec-4c81-b9c8-8ca237a49094.yaml;mistral-mistral-operator"} 1.6726472294604716e+09
prometheus_rule_group_last_evaluation_timestamp_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/mongodb-prometheus-mongodb-rules-662a31e8-0419-44e0-9a6e-1cd71ecfe13e.yaml;mongodb-mongodb-operator"} 1.6726472282698576e+09
prometheus_rule_group_last_evaluation_timestamp_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-cert-exporter-prometheus-rules-be00bda8-3ce7-4de6-aa54-60de240d847f.yaml;cert-exporter"} 1.6726472273356035e+09
prometheus_rule_group_last_evaluation_timestamp_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;AlertManager"} 1.6726472317202709e+09
prometheus_rule_group_last_evaluation_timestamp_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;CoreDnsAlerts"} 1.672647235354644e+09
prometheus_rule_group_last_evaluation_timestamp_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;DockerContainers"} 1.6726472348562593e+09
prometheus_rule_group_last_evaluation_timestamp_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;Etcd"} 1.6726472216685116e+09
prometheus_rule_group_last_evaluation_timestamp_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;KubebernetesAlerts"} 1.6726472208085809e+09
prometheus_rule_group_last_evaluation_timestamp_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;NginxIngressAlerts"} 1.6726472366395535e+09
prometheus_rule_group_last_evaluation_timestamp_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;NodeExporters"} 1.672647217973402e+09
prometheus_rule_group_last_evaluation_timestamp_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;NodeProcesses"} 1.672647228405819e+09
prometheus_rule_group_last_evaluation_timestamp_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;SelfMonitoring"} 1.6726472113582654e+09
prometheus_rule_group_last_evaluation_timestamp_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/opensearch-service-prometheus-opensearch-service-rules-d7a62102-5e19-427b-8cec-9f49c9520a9d.yaml;opensearch.rules"} 1.6726472241294677e+09
prometheus_rule_group_last_evaluation_timestamp_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/pg-metr-prometheus-postgres-service-rules-dc530ea7-cf17-4e3d-bf63-ada04b027880.yaml;pg-metr-postgres-operator"} 1.6726472393907888e+09
prometheus_rule_group_last_evaluation_timestamp_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/postgres-upgrade-prometheus-clickhouse-service-rules-80ab245c-653e-4fe4-952b-3fc62c8ddee0.yaml;postgres-upgrade-clickhouse-operator"} 1.6726472309098163e+09
prometheus_rule_group_last_evaluation_timestamp_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/profiler-cassandra-prometheus-cassandra-rules-3b7ef1e2-2e09-4776-a7c9-3af374492907.yaml;profiler-cassandra-cassandra-operator"} 1.6726472247488592e+09
prometheus_rule_group_last_evaluation_timestamp_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/rabbitmq-prometheus-rabbitmq-service-rules-59b077e6-5aa6-49ca-995b-55225bdfa9a6.yaml;general.rules"} 1.672647220430341e+09
prometheus_rule_group_last_evaluation_timestamp_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/redis-prometheus-redis-rules-8aceaa70-6b23-4d0c-be5e-d445627aaa0e.yaml;redis-redis-operator"} 1.6726472263562508e+09
prometheus_rule_group_last_evaluation_timestamp_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/spark-operator-gcp-spark-operator-prometheus-rules-28d8b3a4-83b9-457a-9dcd-4b915ff141b0.yaml;spark-operator-gcp-"} 1.6726472290790093e+09
prometheus_rule_group_last_evaluation_timestamp_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/streaming-service-prometheus-streaming-service-rules-ca22df5a-e469-4c53-a00e-4ce9d2b8d1b1.yaml;streaming-service-streaming-service-streaming-service"} 1.672647227892915e+09
prometheus_rule_group_last_evaluation_timestamp_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/vault-service-prometheus-vault-service-rules-e526a2e8-d9ed-4748-a535-6bedefae75e8.yaml;vault-service-vault-operator-vault-service"} 1.6726472163301296e+09
prometheus_rule_group_last_evaluation_timestamp_seconds{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/zookeeper-service-prometheus-zookeeper-service-rules-8646f671-901d-40f4-9ef7-8b7535ce08c3.yaml;zookeeper-service-zookeeper-service-zookeeper-service"} 1.6726472169148827e+09
# HELP prometheus_rule_group_rules The number of rules.
# TYPE prometheus_rule_group_rules gauge
prometheus_rule_group_rules{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/anko0714airflowbkp-airflow-anko0714airflowbkp-15b95c5c-90e3-4ec6-b92e-08c69e578a6e.yaml;general.kubeRules"} 10
prometheus_rule_group_rules{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/arango-db-a-prometheus-arangodb-service-rules-78cd09de-3e41-4795-a7ad-aba41c42d51d.yaml;general.rules"} 7
prometheus_rule_group_rules{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/arangodb-graph-test-prometheus-arangodb-service-rules-20e72473-b95d-4df3-915b-16a9ff774836.yaml;general.rules"} 7
prometheus_rule_group_rules{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/arangodb-prometheus-arangodb-service-rules-75974246-3cf7-474e-adfd-688c798c25d1.yaml;general.rules"} 6
prometheus_rule_group_rules{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/consul-service-prometheus-consul-rules-7ec3c416-9681-4630-832e-65285b6200cd.yaml;general.rules"} 5
prometheus_rule_group_rules{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/core-mistral-alert-7e644bba-9d46-40b9-9ca1-47189e9aa887.yaml;core-mistral-service"} 24
prometheus_rule_group_rules{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/core-prometheus-cdc-streaming-service-rules-39068d15-ab3e-49c3-85c3-97089f4a831f.yaml;core-cdc-streaming-platform"} 18
prometheus_rule_group_rules{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/core-prometheus-cloud-streaming-service-rules-b6c45b5b-4656-4240-8b08-cf538881cbed.yaml;core-streaming-service-core"} 18
prometheus_rule_group_rules{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/default-prometheus-arangodb-service-rules-7f343f20-16b7-4f90-a284-8e244507be69.yaml;general.rules"} 6
prometheus_rule_group_rules{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/dev-helm-install-mistral-alert-abbb4231-aba7-4ee8-9a40-25abba628086.yaml;dev-helm-install-mistral-operator-dev-helm-install"} 24
prometheus_rule_group_rules{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/dev-mistral-demo-mistral-alert-c30038fd-6be0-45a6-8d1a-402e9bd863d2.yaml;dev-mistral-demo-mistral-service"} 24
prometheus_rule_group_rules{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/dev-mistral-mistral-alert-682c78ec-8d06-43ca-b7cf-936694c0fae2.yaml;dev-mistral-mistral-service"} 24
prometheus_rule_group_rules{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/domain-search-prometheus-opensearch-service-rules-0a349e4b-db79-42f1-8265-0f1e104a3d41.yaml;opensearch.rules"} 9
prometheus_rule_group_rules{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/elasticsearch-cluster-prometheus-elasticsearch-service-rules-ba16a3c6-bd4f-445c-830c-2e0c8e600e70.yaml;elasticsearch.rules"} 9
prometheus_rule_group_rules{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/elasticsearch-service-prometheus-elasticsearch-service-rules-2be23428-fad5-4626-afc0-d5601aca6f01.yaml;ElasticsearchAlerts"} 9
prometheus_rule_group_rules{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/hadoop-monitoring-hadoop-monitoring-prometheus-rules-88f89818-d5ed-4b63-90b5-5e80f0cdf4ee.yaml;general.rules"} 29
prometheus_rule_group_rules{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/hadoop-monitoring-hadoop-monitoring-prometheus-rules2-af3957d3-435b-4f5f-9719-09c7487bd66b.yaml;general.rules"} 1
prometheus_rule_group_rules{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/hadoopmonitoring-hadoop-monitoring-prometheus-rules-4d478506-60d0-49c1-b5f5-35f440f2f495.yaml;hadoopmonitoring-monitor.general.rules"} 33
prometheus_rule_group_rules{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/kafka-service-prometheus-kafka-service-rules-aa79acbc-51be-49ca-ad42-3b59bcc63696.yaml;kafka-service-kafka-service-kafka-service"} 7
prometheus_rule_group_rules{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/mistral-mistral-alert-8dd98d92-e2ec-4c81-b9c8-8ca237a49094.yaml;mistral-mistral-operator"} 24
prometheus_rule_group_rules{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/mongodb-prometheus-mongodb-rules-662a31e8-0419-44e0-9a6e-1cd71ecfe13e.yaml;mongodb-mongodb-operator"} 9
prometheus_rule_group_rules{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-cert-exporter-prometheus-rules-be00bda8-3ce7-4de6-aa54-60de240d847f.yaml;cert-exporter"} 3
prometheus_rule_group_rules{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;AlertManager"} 3
prometheus_rule_group_rules{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;CoreDnsAlerts"} 9
prometheus_rule_group_rules{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;DockerContainers"} 4
prometheus_rule_group_rules{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;Etcd"} 10
prometheus_rule_group_rules{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;KubebernetesAlerts"} 30
prometheus_rule_group_rules{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;NginxIngressAlerts"} 3
prometheus_rule_group_rules{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;NodeExporters"} 14
prometheus_rule_group_rules{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;NodeProcesses"} 1
prometheus_rule_group_rules{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/monitoring-prometheus-rules-3aee8cfe-f1ae-473b-8b01-7ef2b8ecf714.yaml;SelfMonitoring"} 22
prometheus_rule_group_rules{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/opensearch-service-prometheus-opensearch-service-rules-d7a62102-5e19-427b-8cec-9f49c9520a9d.yaml;opensearch.rules"} 9
prometheus_rule_group_rules{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/pg-metr-prometheus-postgres-service-rules-dc530ea7-cf17-4e3d-bf63-ada04b027880.yaml;pg-metr-postgres-operator"} 34
prometheus_rule_group_rules{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/postgres-upgrade-prometheus-clickhouse-service-rules-80ab245c-653e-4fe4-952b-3fc62c8ddee0.yaml;postgres-upgrade-clickhouse-operator"} 32
prometheus_rule_group_rules{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/profiler-cassandra-prometheus-cassandra-rules-3b7ef1e2-2e09-4776-a7c9-3af374492907.yaml;profiler-cassandra-cassandra-operator"} 9
prometheus_rule_group_rules{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/rabbitmq-prometheus-rabbitmq-service-rules-59b077e6-5aa6-49ca-995b-55225bdfa9a6.yaml;general.rules"} 5
prometheus_rule_group_rules{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/redis-prometheus-redis-rules-8aceaa70-6b23-4d0c-be5e-d445627aaa0e.yaml;redis-redis-operator"} 5
prometheus_rule_group_rules{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/spark-operator-gcp-spark-operator-prometheus-rules-28d8b3a4-83b9-457a-9dcd-4b915ff141b0.yaml;spark-operator-gcp-"} 4
prometheus_rule_group_rules{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/streaming-service-prometheus-streaming-service-rules-ca22df5a-e469-4c53-a00e-4ce9d2b8d1b1.yaml;streaming-service-streaming-service-streaming-service"} 18
prometheus_rule_group_rules{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/vault-service-prometheus-vault-service-rules-e526a2e8-d9ed-4748-a535-6bedefae75e8.yaml;vault-service-vault-operator-vault-service"} 7
prometheus_rule_group_rules{rule_group="/etc/prometheus/rules/prometheus-k8s-rulefiles-0/zookeeper-service-prometheus-zookeeper-service-rules-8646f671-901d-40f4-9ef7-8b7535ce08c3.yaml;zookeeper-service-zookeeper-service-zookeeper-service"} 5
# HELP prometheus_sd_azure_failures_total Number of Azure service discovery refresh failures.
# TYPE prometheus_sd_azure_failures_total counter
prometheus_sd_azure_failures_total 0
# HELP prometheus_sd_consul_rpc_duration_seconds The duration of a Consul RPC call in seconds.
# TYPE prometheus_sd_consul_rpc_duration_seconds summary
prometheus_sd_consul_rpc_duration_seconds{call="service",endpoint="catalog",quantile="0.5"} NaN
prometheus_sd_consul_rpc_duration_seconds{call="service",endpoint="catalog",quantile="0.9"} NaN
prometheus_sd_consul_rpc_duration_seconds{call="service",endpoint="catalog",quantile="0.99"} NaN
prometheus_sd_consul_rpc_duration_seconds_sum{call="service",endpoint="catalog"} 0
prometheus_sd_consul_rpc_duration_seconds_count{call="service",endpoint="catalog"} 0
prometheus_sd_consul_rpc_duration_seconds{call="services",endpoint="catalog",quantile="0.5"} NaN
prometheus_sd_consul_rpc_duration_seconds{call="services",endpoint="catalog",quantile="0.9"} NaN
prometheus_sd_consul_rpc_duration_seconds{call="services",endpoint="catalog",quantile="0.99"} NaN
prometheus_sd_consul_rpc_duration_seconds_sum{call="services",endpoint="catalog"} 0
prometheus_sd_consul_rpc_duration_seconds_count{call="services",endpoint="catalog"} 0
# HELP prometheus_sd_consul_rpc_failures_total The number of Consul RPC call failures.
# TYPE prometheus_sd_consul_rpc_failures_total counter
prometheus_sd_consul_rpc_failures_total 0
# HELP prometheus_sd_discovered_targets Current number of discovered targets.
# TYPE prometheus_sd_discovered_targets gauge
prometheus_sd_discovered_targets{config="config-0",name="notify"} 51
prometheus_sd_discovered_targets{config="config-1",name="notify"} 1
prometheus_sd_discovered_targets{config="podMonitor/monitoring-examples/fake-log-generator-pod-monitor/0",name="scrape"} 0
prometheus_sd_discovered_targets{config="podMonitor/monitoring-examples/monitoring-example-pod-monitor/0",name="scrape"} 0
prometheus_sd_discovered_targets{config="podMonitor/monitoring/monitoring-alertmanager-pod-monitor/0",name="scrape"} 24
prometheus_sd_discovered_targets{config="podMonitor/monitoring/monitoring-grafana-operator-pod-monitor/0",name="scrape"} 24
prometheus_sd_discovered_targets{config="podMonitor/monitoring/monitoring-grafana-pod-monitor/0",name="scrape"} 24
prometheus_sd_discovered_targets{config="podMonitor/monitoring/monitoring-nginx-ingress-pod-monitor/0",name="scrape"} 621
prometheus_sd_discovered_targets{config="podMonitor/monitoring/monitoring-prometheus-operator-pod-monitor/0",name="scrape"} 24
prometheus_sd_discovered_targets{config="podMonitor/monitoring/monitoring-prometheus-pod-monitor/0",name="scrape"} 24
prometheus_sd_discovered_targets{config="podMonitor/spark-operator-gcp/spark-operator-podmonitor/0",name="scrape"} 2
prometheus_sd_discovered_targets{config="podMonitor/vault-service/vault-operator-pod-monitor/0",name="scrape"} 8
prometheus_sd_discovered_targets{config="probe/monitoring-examples/blackbox-ingress-probe",name="scrape"} 1
prometheus_sd_discovered_targets{config="probe/monitoring-examples/blackbox-static-probe",name="scrape"} 2
prometheus_sd_discovered_targets{config="serviceMonitor/anko0714airflowbkp/airflow-anko0714airflowbkp/0",name="scrape"} 0
prometheus_sd_discovered_targets{config="serviceMonitor/arango-db-a/arangodb-service-monitor/0",name="scrape"} 0
prometheus_sd_discovered_targets{config="serviceMonitor/arangodb-graph-test/arangodb-service-monitor/0",name="scrape"} 0
prometheus_sd_discovered_targets{config="serviceMonitor/arangodb/arangodb-service-monitor/0",name="scrape"} 12
prometheus_sd_discovered_targets{config="serviceMonitor/consul-service/consul-service-monitor/0",name="scrape"} 152
prometheus_sd_discovered_targets{config="serviceMonitor/core/cdc-control-monitor/0",name="scrape"} 64
prometheus_sd_discovered_targets{config="serviceMonitor/core/cdc-stream-processor-monitor/0",name="scrape"} 64
prometheus_sd_discovered_targets{config="serviceMonitor/core/cdc-streaming-service-monitor-jmx-exporter/0",name="scrape"} 64
prometheus_sd_discovered_targets{config="serviceMonitor/core/cdc-test-api-monitor/0",name="scrape"} 64
prometheus_sd_discovered_targets{config="serviceMonitor/core/cdc-test-db-monitor/0",name="scrape"} 64
prometheus_sd_discovered_targets{config="serviceMonitor/core/cloud-streaming-service-monitor-jmx-exporter/0",name="scrape"} 64
prometheus_sd_discovered_targets{config="serviceMonitor/core/mistral-service-monitor/0",name="scrape"} 64
prometheus_sd_discovered_targets{config="serviceMonitor/default/arangodb-service-monitor/0",name="scrape"} 3
prometheus_sd_discovered_targets{config="serviceMonitor/dev-helm-install/mistral-service-monitor/0",name="scrape"} 9
prometheus_sd_discovered_targets{config="serviceMonitor/dev-mistral-demo/mistral-service-monitor/0",name="scrape"} 4
prometheus_sd_discovered_targets{config="serviceMonitor/dev-mistral/mistral-service-monitor/0",name="scrape"} 4
prometheus_sd_discovered_targets{config="serviceMonitor/domain-search/opensearch-service-monitor/0",name="scrape"} 6
prometheus_sd_discovered_targets{config="serviceMonitor/elasticsearch-cluster/elasticsearch-service-monitor/0",name="scrape"} 43
prometheus_sd_discovered_targets{config="serviceMonitor/elasticsearch-service/elasticsearch-service-monitor/0",name="scrape"} 25
prometheus_sd_discovered_targets{config="serviceMonitor/envoy/monitoring-envoy/0",name="scrape"} 2
prometheus_sd_discovered_targets{config="serviceMonitor/envoy/monitoring-keycloak/0",name="scrape"} 2
prometheus_sd_discovered_targets{config="serviceMonitor/kafka-service/kafka-lag-exporter-service-monitor/0",name="scrape"} 87
prometheus_sd_discovered_targets{config="serviceMonitor/kafka-service/kafka-service-monitor-jmx-exporter/0",name="scrape"} 87
prometheus_sd_discovered_targets{config="serviceMonitor/kafka-service/kafka-service-monitor/0",name="scrape"} 87
prometheus_sd_discovered_targets{config="serviceMonitor/logging/logging-service-operator-metrics/0",name="scrape"} 0
prometheus_sd_discovered_targets{config="serviceMonitor/mistral/mistral-service-monitor/0",name="scrape"} 4
prometheus_sd_discovered_targets{config="serviceMonitor/mongodb/mongo-backup-exporter-service-monitor/0",name="scrape"} 9
prometheus_sd_discovered_targets{config="serviceMonitor/mongodb/mongo-dbaas-exporter-service-monitor/0",name="scrape"} 9
prometheus_sd_discovered_targets{config="serviceMonitor/mongodb/mongodb-prometheus-exporter-service-monitor/0",name="scrape"} 9
prometheus_sd_discovered_targets{config="serviceMonitor/monitoring-examples/go-example-service-service-monitor/0",name="scrape"} 0
prometheus_sd_discovered_targets{config="serviceMonitor/monitoring-examples/go-example-service-service-monitor/1",name="scrape"} 0
prometheus_sd_discovered_targets{config="serviceMonitor/monitoring-examples/java-example-service-monitor/0",name="scrape"} 0
prometheus_sd_discovered_targets{config="serviceMonitor/monitoring-examples/monitoring-example-service-monitor/0",name="scrape"} 0
prometheus_sd_discovered_targets{config="serviceMonitor/monitoring/monitoring-cert-exporter/0",name="scrape"} 51
prometheus_sd_discovered_targets{config="serviceMonitor/monitoring/monitoring-core-dns-service-monitor/0",name="scrape"} 9
prometheus_sd_discovered_targets{config="serviceMonitor/monitoring/monitoring-etcd-service-monitor/0",name="scrape"} 9
prometheus_sd_discovered_targets{config="serviceMonitor/monitoring/monitoring-kube-apiserver-service-monitor/0",name="scrape"} 3
prometheus_sd_discovered_targets{config="serviceMonitor/monitoring/monitoring-kube-state-metrics/0",name="scrape"} 51
prometheus_sd_discovered_targets{config="serviceMonitor/monitoring/monitoring-kube-state-metrics/1",name="scrape"} 51
prometheus_sd_discovered_targets{config="serviceMonitor/monitoring/monitoring-kube-state-metrics/2",name="scrape"} 51
prometheus_sd_discovered_targets{config="serviceMonitor/monitoring/monitoring-kubelet-service-monitor/0",name="scrape"} 51
prometheus_sd_discovered_targets{config="serviceMonitor/monitoring/monitoring-kubelet-service-monitor/1",name="scrape"} 51
prometheus_sd_discovered_targets{config="serviceMonitor/monitoring/monitoring-node-exporter/0",name="scrape"} 51
prometheus_sd_discovered_targets{config="serviceMonitor/opensearch-service/opensearch-service-monitor/0",name="scrape"} 55
prometheus_sd_discovered_targets{config="serviceMonitor/pg-metr/postgres-service-monitor/0",name="scrape"} 13
prometheus_sd_discovered_targets{config="serviceMonitor/postgres-upgrade/clickhouse-backup-daemon-service-monitor/0",name="scrape"} 1
prometheus_sd_discovered_targets{config="serviceMonitor/postgres-upgrade/clickhouse-service-monitor/0",name="scrape"} 1
prometheus_sd_discovered_targets{config="serviceMonitor/profiler-cassandra/cassandra-exporter-service-monitor/0",name="scrape"} 69
prometheus_sd_discovered_targets{config="serviceMonitor/profiler-test/esc-collector-service-monitor/0",name="scrape"} 5
prometheus_sd_discovered_targets{config="serviceMonitor/profiler-test/esc-test-service-monitor/0",name="scrape"} 5
prometheus_sd_discovered_targets{config="serviceMonitor/profiler-test/esc-ui-service-monitor/0",name="scrape"} 5
prometheus_sd_discovered_targets{config="serviceMonitor/profiler/esc-collector-service-monitor/0",name="scrape"} 5
prometheus_sd_discovered_targets{config="serviceMonitor/profiler/esc-test-service-monitor/0",name="scrape"} 5
prometheus_sd_discovered_targets{config="serviceMonitor/profiler/esc-ui-service-monitor/0",name="scrape"} 5
prometheus_sd_discovered_targets{config="serviceMonitor/rabbitmq/rabbitmq-per-object-service-monitor/0",name="scrape"} 83
prometheus_sd_discovered_targets{config="serviceMonitor/rabbitmq/rabbitmq-service-monitor-old/0",name="scrape"} 83
prometheus_sd_discovered_targets{config="serviceMonitor/rabbitmq/rabbitmq-service-monitor/0",name="scrape"} 83
prometheus_sd_discovered_targets{config="serviceMonitor/redis/redis-prometheus-exporter-service-monitor/0",name="scrape"} 5
prometheus_sd_discovered_targets{config="serviceMonitor/streaming-service/streaming-service-monitor-jmx-exporter/0",name="scrape"} 14
prometheus_sd_discovered_targets{config="serviceMonitor/streaming-service/streaming-service-monitor/0",name="scrape"} 14
prometheus_sd_discovered_targets{config="serviceMonitor/vault-service/vault-service-monitor/0",name="scrape"} 19
prometheus_sd_discovered_targets{config="serviceMonitor/zookeeper-service/zookeeper-service-monitor-jmx-exporter/0",name="scrape"} 66
prometheus_sd_discovered_targets{config="serviceMonitor/zookeeper-service/zookeeper-service-monitor/0",name="scrape"} 66
# HELP prometheus_sd_dns_lookup_failures_total The number of DNS-SD lookup failures.
# TYPE prometheus_sd_dns_lookup_failures_total counter
prometheus_sd_dns_lookup_failures_total 0
# HELP prometheus_sd_dns_lookups_total The number of DNS-SD lookups.
# TYPE prometheus_sd_dns_lookups_total counter
prometheus_sd_dns_lookups_total 0
# HELP prometheus_sd_failed_configs Current number of service discovery configurations that failed to load.
# TYPE prometheus_sd_failed_configs gauge
prometheus_sd_failed_configs{name="notify"} 0
prometheus_sd_failed_configs{name="scrape"} 0
# HELP prometheus_sd_file_read_errors_total The number of File-SD read errors.
# TYPE prometheus_sd_file_read_errors_total counter
prometheus_sd_file_read_errors_total 0
# HELP prometheus_sd_file_scan_duration_seconds The duration of the File-SD scan in seconds.
# TYPE prometheus_sd_file_scan_duration_seconds summary
prometheus_sd_file_scan_duration_seconds{quantile="0.5"} NaN
prometheus_sd_file_scan_duration_seconds{quantile="0.9"} NaN
prometheus_sd_file_scan_duration_seconds{quantile="0.99"} NaN
prometheus_sd_file_scan_duration_seconds_sum 0
prometheus_sd_file_scan_duration_seconds_count 0
# HELP prometheus_sd_http_failures_total Number of HTTP service discovery refresh failures.
# TYPE prometheus_sd_http_failures_total counter
prometheus_sd_http_failures_total 0
# HELP prometheus_sd_kubernetes_events_total The number of Kubernetes events handled.
# TYPE prometheus_sd_kubernetes_events_total counter
prometheus_sd_kubernetes_events_total{event="add",role="endpoints"} 219
prometheus_sd_kubernetes_events_total{event="add",role="endpointslice"} 0
prometheus_sd_kubernetes_events_total{event="add",role="ingress"} 1
prometheus_sd_kubernetes_events_total{event="add",role="node"} 0
prometheus_sd_kubernetes_events_total{event="add",role="pod"} 283
prometheus_sd_kubernetes_events_total{event="add",role="service"} 208
prometheus_sd_kubernetes_events_total{event="delete",role="endpoints"} 0
prometheus_sd_kubernetes_events_total{event="delete",role="endpointslice"} 0
prometheus_sd_kubernetes_events_total{event="delete",role="ingress"} 0
prometheus_sd_kubernetes_events_total{event="delete",role="node"} 0
prometheus_sd_kubernetes_events_total{event="delete",role="pod"} 0
prometheus_sd_kubernetes_events_total{event="delete",role="service"} 0
prometheus_sd_kubernetes_events_total{event="update",role="endpoints"} 15941
prometheus_sd_kubernetes_events_total{event="update",role="endpointslice"} 0
prometheus_sd_kubernetes_events_total{event="update",role="ingress"} 76
prometheus_sd_kubernetes_events_total{event="update",role="node"} 0
prometheus_sd_kubernetes_events_total{event="update",role="pod"} 32072
prometheus_sd_kubernetes_events_total{event="update",role="service"} 15656
# HELP prometheus_sd_kubernetes_http_request_duration_seconds Summary of latencies for HTTP requests to the Kubernetes API by endpoint.
# TYPE prometheus_sd_kubernetes_http_request_duration_seconds summary
prometheus_sd_kubernetes_http_request_duration_seconds_sum{endpoint="/api/v1/namespaces/%7Bnamespace%7D/endpoints"} 28.416563796
prometheus_sd_kubernetes_http_request_duration_seconds_count{endpoint="/api/v1/namespaces/%7Bnamespace%7D/endpoints"} 66
prometheus_sd_kubernetes_http_request_duration_seconds_sum{endpoint="/api/v1/namespaces/%7Bnamespace%7D/pods"} 32.367052503000004
prometheus_sd_kubernetes_http_request_duration_seconds_count{endpoint="/api/v1/namespaces/%7Bnamespace%7D/pods"} 74
prometheus_sd_kubernetes_http_request_duration_seconds_sum{endpoint="/api/v1/namespaces/%7Bnamespace%7D/services"} 28.777541830999994
prometheus_sd_kubernetes_http_request_duration_seconds_count{endpoint="/api/v1/namespaces/%7Bnamespace%7D/services"} 66
prometheus_sd_kubernetes_http_request_duration_seconds_sum{endpoint="/api/v1/pods"} 0.947325647
prometheus_sd_kubernetes_http_request_duration_seconds_count{endpoint="/api/v1/pods"} 2
prometheus_sd_kubernetes_http_request_duration_seconds_sum{endpoint="/apis/networking.k8s.io/v1/namespaces/%7Bnamespace%7D/ingresses"} 0.020346129
prometheus_sd_kubernetes_http_request_duration_seconds_count{endpoint="/apis/networking.k8s.io/v1/namespaces/%7Bnamespace%7D/ingresses"} 1
prometheus_sd_kubernetes_http_request_duration_seconds_sum{endpoint="/version"} 1.309936846
prometheus_sd_kubernetes_http_request_duration_seconds_count{endpoint="/version"} 2
# HELP prometheus_sd_kubernetes_http_request_total Total number of HTTP requests to the Kubernetes API by status code.
# TYPE prometheus_sd_kubernetes_http_request_total counter
prometheus_sd_kubernetes_http_request_total{status_code="200"} 10892
prometheus_sd_kubernetes_http_request_total{status_code="<error>"} 62
# HELP prometheus_sd_kubernetes_workqueue_depth Current depth of the work queue.
# TYPE prometheus_sd_kubernetes_workqueue_depth gauge
prometheus_sd_kubernetes_workqueue_depth{queue_name="endpoints"} 11
prometheus_sd_kubernetes_workqueue_depth{queue_name="ingress"} 0
prometheus_sd_kubernetes_workqueue_depth{queue_name="pod"} 0
# HELP prometheus_sd_kubernetes_workqueue_items_total Total number of items added to the work queue.
# TYPE prometheus_sd_kubernetes_workqueue_items_total counter
prometheus_sd_kubernetes_workqueue_items_total{queue_name="endpoints"} 29870
prometheus_sd_kubernetes_workqueue_items_total{queue_name="ingress"} 77
prometheus_sd_kubernetes_workqueue_items_total{queue_name="pod"} 32355
# HELP prometheus_sd_kubernetes_workqueue_latency_seconds How long an item stays in the work queue.
# TYPE prometheus_sd_kubernetes_workqueue_latency_seconds summary
prometheus_sd_kubernetes_workqueue_latency_seconds_sum{queue_name="endpoints"} 78.7271875099994
prometheus_sd_kubernetes_workqueue_latency_seconds_count{queue_name="endpoints"} 29859
prometheus_sd_kubernetes_workqueue_latency_seconds_sum{queue_name="ingress"} 0.091597666
prometheus_sd_kubernetes_workqueue_latency_seconds_count{queue_name="ingress"} 77
prometheus_sd_kubernetes_workqueue_latency_seconds_sum{queue_name="pod"} 130.72835804399975
prometheus_sd_kubernetes_workqueue_latency_seconds_count{queue_name="pod"} 32355
# HELP prometheus_sd_kubernetes_workqueue_longest_running_processor_seconds Duration of the longest running processor in the work queue.
# TYPE prometheus_sd_kubernetes_workqueue_longest_running_processor_seconds gauge
prometheus_sd_kubernetes_workqueue_longest_running_processor_seconds{queue_name="endpoints"} 0
prometheus_sd_kubernetes_workqueue_longest_running_processor_seconds{queue_name="ingress"} 0
prometheus_sd_kubernetes_workqueue_longest_running_processor_seconds{queue_name="pod"} 0
# HELP prometheus_sd_kubernetes_workqueue_unfinished_work_seconds How long an item has remained unfinished in the work queue.
# TYPE prometheus_sd_kubernetes_workqueue_unfinished_work_seconds gauge
prometheus_sd_kubernetes_workqueue_unfinished_work_seconds{queue_name="endpoints"} 0
prometheus_sd_kubernetes_workqueue_unfinished_work_seconds{queue_name="ingress"} 0
prometheus_sd_kubernetes_workqueue_unfinished_work_seconds{queue_name="pod"} 0
# HELP prometheus_sd_kubernetes_workqueue_work_duration_seconds How long processing an item from the work queue takes.
# TYPE prometheus_sd_kubernetes_workqueue_work_duration_seconds summary
prometheus_sd_kubernetes_workqueue_work_duration_seconds_sum{queue_name="endpoints"} 8.697643657000052
prometheus_sd_kubernetes_workqueue_work_duration_seconds_count{queue_name="endpoints"} 29859
prometheus_sd_kubernetes_workqueue_work_duration_seconds_sum{queue_name="ingress"} 0.0019781339999999994
prometheus_sd_kubernetes_workqueue_work_duration_seconds_count{queue_name="ingress"} 77
prometheus_sd_kubernetes_workqueue_work_duration_seconds_sum{queue_name="pod"} 2.145918401000005
prometheus_sd_kubernetes_workqueue_work_duration_seconds_count{queue_name="pod"} 32355
# HELP prometheus_sd_kuma_fetch_duration_seconds The duration of a Kuma MADS fetch call.
# TYPE prometheus_sd_kuma_fetch_duration_seconds summary
prometheus_sd_kuma_fetch_duration_seconds{quantile="0.5"} NaN
prometheus_sd_kuma_fetch_duration_seconds{quantile="0.9"} NaN
prometheus_sd_kuma_fetch_duration_seconds{quantile="0.99"} NaN
prometheus_sd_kuma_fetch_duration_seconds_sum 0
prometheus_sd_kuma_fetch_duration_seconds_count 0
# HELP prometheus_sd_kuma_fetch_failures_total The number of Kuma MADS fetch call failures.
# TYPE prometheus_sd_kuma_fetch_failures_total counter
prometheus_sd_kuma_fetch_failures_total 0
# HELP prometheus_sd_kuma_fetch_skipped_updates_total The number of Kuma MADS fetch calls that result in no updates to the targets.
# TYPE prometheus_sd_kuma_fetch_skipped_updates_total counter
prometheus_sd_kuma_fetch_skipped_updates_total 0
# HELP prometheus_sd_linode_failures_total Number of Linode service discovery refresh failures.
# TYPE prometheus_sd_linode_failures_total counter
prometheus_sd_linode_failures_total 0
# HELP prometheus_sd_nomad_failures_total Number of nomad service discovery refresh failures.
# TYPE prometheus_sd_nomad_failures_total counter
prometheus_sd_nomad_failures_total 0
# HELP prometheus_sd_received_updates_total Total number of update events received from the SD providers.
# TYPE prometheus_sd_received_updates_total counter
prometheus_sd_received_updates_total{name="notify"} 1309
prometheus_sd_received_updates_total{name="scrape"} 60990
# HELP prometheus_sd_updates_total Total number of update events sent to the SD consumers.
# TYPE prometheus_sd_updates_total counter
prometheus_sd_updates_total{name="notify"} 78
prometheus_sd_updates_total{name="scrape"} 4856
# HELP prometheus_target_interval_length_seconds Actual intervals between scrapes.
# TYPE prometheus_target_interval_length_seconds summary
prometheus_target_interval_length_seconds{interval="1m0s",quantile="0.01"} 59.97834156
prometheus_target_interval_length_seconds{interval="1m0s",quantile="0.05"} 59.999014794
prometheus_target_interval_length_seconds{interval="1m0s",quantile="0.5"} 60.00001158
prometheus_target_interval_length_seconds{interval="1m0s",quantile="0.9"} 60.000789905
prometheus_target_interval_length_seconds{interval="1m0s",quantile="0.99"} 60.032832484
prometheus_target_interval_length_seconds_sum{interval="1m0s"} 1.0038602439958253e+06
prometheus_target_interval_length_seconds_count{interval="1m0s"} 16731
prometheus_target_interval_length_seconds{interval="20s",quantile="0.01"} 19.989795005
prometheus_target_interval_length_seconds{interval="20s",quantile="0.05"} 19.994497927
prometheus_target_interval_length_seconds{interval="20s",quantile="0.5"} 20.000041729
prometheus_target_interval_length_seconds{interval="20s",quantile="0.9"} 20.000839167
prometheus_target_interval_length_seconds{interval="20s",quantile="0.99"} 20.009598724
prometheus_target_interval_length_seconds_sum{interval="20s"} 182720.1146782967
prometheus_target_interval_length_seconds_count{interval="20s"} 9136
prometheus_target_interval_length_seconds{interval="2m0s",quantile="0.01"} 119.999855193
prometheus_target_interval_length_seconds{interval="2m0s",quantile="0.05"} 119.999855193
prometheus_target_interval_length_seconds{interval="2m0s",quantile="0.5"} 120.000114285
prometheus_target_interval_length_seconds{interval="2m0s",quantile="0.9"} 120.000534528
prometheus_target_interval_length_seconds{interval="2m0s",quantile="0.99"} 120.000534528
prometheus_target_interval_length_seconds_sum{interval="2m0s"} 45600.00456756995
prometheus_target_interval_length_seconds_count{interval="2m0s"} 380
prometheus_target_interval_length_seconds{interval="30s",quantile="0.01"} 29.991807138
prometheus_target_interval_length_seconds{interval="30s",quantile="0.05"} 29.99912114
prometheus_target_interval_length_seconds{interval="30s",quantile="0.5"} 30.000053864
prometheus_target_interval_length_seconds{interval="30s",quantile="0.9"} 30.000778896
prometheus_target_interval_length_seconds{interval="30s",quantile="0.99"} 30.011451223
prometheus_target_interval_length_seconds_sum{interval="30s"} 3.971883504577789e+06
prometheus_target_interval_length_seconds_count{interval="30s"} 132396
prometheus_target_interval_length_seconds{interval="5s",quantile="0.01"} 4.986981897
prometheus_target_interval_length_seconds{interval="5s",quantile="0.05"} 4.998066942
prometheus_target_interval_length_seconds{interval="5s",quantile="0.5"} 4.999967736
prometheus_target_interval_length_seconds{interval="5s",quantile="0.9"} 5.000875523
prometheus_target_interval_length_seconds{interval="5s",quantile="0.99"} 5.01310281
prometheus_target_interval_length_seconds_sum{interval="5s"} 45680.12913268399
prometheus_target_interval_length_seconds_count{interval="5s"} 9136
# HELP prometheus_target_metadata_cache_bytes The number of bytes that are currently used for storing metric metadata in the cache
# TYPE prometheus_target_metadata_cache_bytes gauge
prometheus_target_metadata_cache_bytes{scrape_job="podMonitor/monitoring/monitoring-alertmanager-pod-monitor/0"} 7152
prometheus_target_metadata_cache_bytes{scrape_job="podMonitor/monitoring/monitoring-grafana-operator-pod-monitor/0"} 2788
prometheus_target_metadata_cache_bytes{scrape_job="podMonitor/monitoring/monitoring-grafana-pod-monitor/0"} 7177
prometheus_target_metadata_cache_bytes{scrape_job="podMonitor/monitoring/monitoring-nginx-ingress-pod-monitor/0"} 18707
prometheus_target_metadata_cache_bytes{scrape_job="podMonitor/monitoring/monitoring-prometheus-operator-pod-monitor/0"} 3435
prometheus_target_metadata_cache_bytes{scrape_job="podMonitor/monitoring/monitoring-prometheus-pod-monitor/0"} 15446
prometheus_target_metadata_cache_bytes{scrape_job="podMonitor/spark-operator-gcp/spark-operator-podmonitor/0"} 7095
prometheus_target_metadata_cache_bytes{scrape_job="podMonitor/vault-service/vault-operator-pod-monitor/0"} 1945
prometheus_target_metadata_cache_bytes{scrape_job="probe/monitoring-examples/blackbox-ingress-probe"} 0
prometheus_target_metadata_cache_bytes{scrape_job="probe/monitoring-examples/blackbox-static-probe"} 0
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/anko0714airflowbkp/airflow-anko0714airflowbkp/0"} 0
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/arango-db-a/arangodb-service-monitor/0"} 0
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/arangodb-graph-test/arangodb-service-monitor/0"} 0
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/arangodb/arangodb-service-monitor/0"} 0
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/consul-service/consul-service-monitor/0"} 57070
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/core/cdc-control-monitor/0"} 3227
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/core/cdc-stream-processor-monitor/0"} 2843
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/core/cdc-streaming-service-monitor-jmx-exporter/0"} 12049
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/core/cdc-test-api-monitor/0"} 1886
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/core/cdc-test-db-monitor/0"} 0
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/core/cloud-streaming-service-monitor-jmx-exporter/0"} 14388
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/core/mistral-service-monitor/0"} 391
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/default/arangodb-service-monitor/0"} 0
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/dev-helm-install/mistral-service-monitor/0"} 0
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/dev-mistral-demo/mistral-service-monitor/0"} 391
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/dev-mistral/mistral-service-monitor/0"} 0
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/domain-search/opensearch-service-monitor/0"} 24924
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/elasticsearch-cluster/elasticsearch-service-monitor/0"} 22236
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/elasticsearch-service/elasticsearch-service-monitor/0"} 20156
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/envoy/monitoring-envoy/0"} 0
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/envoy/monitoring-keycloak/0"} 0
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/kafka-service/kafka-lag-exporter-service-monitor/0"} 2212
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/kafka-service/kafka-service-monitor-jmx-exporter/0"} 36294
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/kafka-service/kafka-service-monitor/0"} 5226
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/mistral/mistral-service-monitor/0"} 391
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/mongodb/mongo-backup-exporter-service-monitor/0"} 736
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/mongodb/mongo-dbaas-exporter-service-monitor/0"} 5899
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/mongodb/mongodb-prometheus-exporter-service-monitor/0"} 62168
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/monitoring-examples/go-example-service-service-monitor/0"} 0
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/monitoring-examples/go-example-service-service-monitor/1"} 0
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/monitoring-examples/java-example-service-monitor/0"} 0
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/monitoring-examples/monitoring-example-service-monitor/0"} 0
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/monitoring/monitoring-cert-exporter/0"} 2002
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/monitoring/monitoring-core-dns-service-monitor/0"} 6336
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/monitoring/monitoring-etcd-service-monitor/0"} 20928
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/monitoring/monitoring-kube-apiserver-service-monitor/0"} 29701
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/monitoring/monitoring-kube-state-metrics/0"} 12575
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/monitoring/monitoring-kube-state-metrics/1"} 12575
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/monitoring/monitoring-kube-state-metrics/2"} 2201
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/monitoring/monitoring-kubelet-service-monitor/0"} 57623
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/monitoring/monitoring-kubelet-service-monitor/1"} 29592
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/monitoring/monitoring-node-exporter/0"} 144271
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/opensearch-service/opensearch-service-monitor/0"} 34428
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/pg-metr/postgres-service-monitor/0"} 3693
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/postgres-upgrade/clickhouse-backup-daemon-service-monitor/0"} 0
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/postgres-upgrade/clickhouse-service-monitor/0"} 1852
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/profiler-cassandra/cassandra-exporter-service-monitor/0"} 2367
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/profiler-test/esc-collector-service-monitor/0"} 4615
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/profiler-test/esc-test-service-monitor/0"} 2811
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/profiler-test/esc-ui-service-monitor/0"} 2811
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/profiler/esc-collector-service-monitor/0"} 4057
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/profiler/esc-test-service-monitor/0"} 2811
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/profiler/esc-ui-service-monitor/0"} 2811
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/rabbitmq/rabbitmq-per-object-service-monitor/0"} 40071
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/rabbitmq/rabbitmq-service-monitor-old/0"} 39960
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/rabbitmq/rabbitmq-service-monitor/0"} 39960
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/redis/redis-prometheus-exporter-service-monitor/0"} 1805
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/streaming-service/streaming-service-monitor-jmx-exporter/0"} 47152
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/streaming-service/streaming-service-monitor/0"} 7786
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/vault-service/vault-service-monitor/0"} 12208
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/zookeeper-service/zookeeper-service-monitor-jmx-exporter/0"} 30729
prometheus_target_metadata_cache_bytes{scrape_job="serviceMonitor/zookeeper-service/zookeeper-service-monitor/0"} 5450
# HELP prometheus_target_metadata_cache_entries Total number of metric metadata entries in the cache
# TYPE prometheus_target_metadata_cache_entries gauge
prometheus_target_metadata_cache_entries{scrape_job="podMonitor/monitoring/monitoring-alertmanager-pod-monitor/0"} 102
prometheus_target_metadata_cache_entries{scrape_job="podMonitor/monitoring/monitoring-grafana-operator-pod-monitor/0"} 49
prometheus_target_metadata_cache_entries{scrape_job="podMonitor/monitoring/monitoring-grafana-pod-monitor/0"} 152
prometheus_target_metadata_cache_entries{scrape_job="podMonitor/monitoring/monitoring-nginx-ingress-pod-monitor/0"} 359
prometheus_target_metadata_cache_entries{scrape_job="podMonitor/monitoring/monitoring-prometheus-operator-pod-monitor/0"} 58
prometheus_target_metadata_cache_entries{scrape_job="podMonitor/monitoring/monitoring-prometheus-pod-monitor/0"} 229
prometheus_target_metadata_cache_entries{scrape_job="podMonitor/spark-operator-gcp/spark-operator-podmonitor/0"} 90
prometheus_target_metadata_cache_entries{scrape_job="podMonitor/vault-service/vault-operator-pod-monitor/0"} 40
prometheus_target_metadata_cache_entries{scrape_job="probe/monitoring-examples/blackbox-ingress-probe"} 0
prometheus_target_metadata_cache_entries{scrape_job="probe/monitoring-examples/blackbox-static-probe"} 0
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/anko0714airflowbkp/airflow-anko0714airflowbkp/0"} 0
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/arango-db-a/arangodb-service-monitor/0"} 0
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/arangodb-graph-test/arangodb-service-monitor/0"} 0
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/arangodb/arangodb-service-monitor/0"} 0
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/consul-service/consul-service-monitor/0"} 1521
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/core/cdc-control-monitor/0"} 56
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/core/cdc-stream-processor-monitor/0"} 49
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/core/cdc-streaming-service-monitor-jmx-exporter/0"} 244
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/core/cdc-test-api-monitor/0"} 35
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/core/cdc-test-db-monitor/0"} 0
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/core/cloud-streaming-service-monitor-jmx-exporter/0"} 292
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/core/mistral-service-monitor/0"} 12
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/default/arangodb-service-monitor/0"} 0
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/dev-helm-install/mistral-service-monitor/0"} 0
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/dev-mistral-demo/mistral-service-monitor/0"} 12
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/dev-mistral/mistral-service-monitor/0"} 0
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/domain-search/opensearch-service-monitor/0"} 690
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/elasticsearch-cluster/elasticsearch-service-monitor/0"} 606
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/elasticsearch-service/elasticsearch-service-monitor/0"} 541
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/envoy/monitoring-envoy/0"} 0
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/envoy/monitoring-keycloak/0"} 0
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/kafka-service/kafka-lag-exporter-service-monitor/0"} 45
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/kafka-service/kafka-service-monitor-jmx-exporter/0"} 367
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/kafka-service/kafka-service-monitor/0"} 70
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/mistral/mistral-service-monitor/0"} 12
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/mongodb/mongo-backup-exporter-service-monitor/0"} 23
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/mongodb/mongo-dbaas-exporter-service-monitor/0"} 73
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/mongodb/mongodb-prometheus-exporter-service-monitor/0"} 1477
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/monitoring-examples/go-example-service-service-monitor/0"} 0
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/monitoring-examples/go-example-service-service-monitor/1"} 0
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/monitoring-examples/java-example-service-monitor/0"} 0
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/monitoring-examples/monitoring-example-service-monitor/0"} 0
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/monitoring/monitoring-cert-exporter/0"} 40
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/monitoring/monitoring-core-dns-service-monitor/0"} 118
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/monitoring/monitoring-etcd-service-monitor/0"} 381
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/monitoring/monitoring-kube-apiserver-service-monitor/0"} 359
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/monitoring/monitoring-kube-state-metrics/0"} 220
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/monitoring/monitoring-kube-state-metrics/1"} 220
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/monitoring/monitoring-kube-state-metrics/2"} 41
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/monitoring/monitoring-kubelet-service-monitor/0"} 825
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/monitoring/monitoring-kubelet-service-monitor/1"} 594
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/monitoring/monitoring-node-exporter/0"} 2903
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/opensearch-service/opensearch-service-monitor/0"} 987
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/pg-metr/postgres-service-monitor/0"} 97
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/postgres-upgrade/clickhouse-backup-daemon-service-monitor/0"} 0
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/postgres-upgrade/clickhouse-service-monitor/0"} 37
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/profiler-cassandra/cassandra-exporter-service-monitor/0"} 393
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/profiler-test/esc-collector-service-monitor/0"} 119
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/profiler-test/esc-test-service-monitor/0"} 48
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/profiler-test/esc-ui-service-monitor/0"} 48
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/profiler/esc-collector-service-monitor/0"} 70
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/profiler/esc-test-service-monitor/0"} 48
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/profiler/esc-ui-service-monitor/0"} 48
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/rabbitmq/rabbitmq-per-object-service-monitor/0"} 663
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/rabbitmq/rabbitmq-service-monitor-old/0"} 663
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/rabbitmq/rabbitmq-service-monitor/0"} 663
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/redis/redis-prometheus-exporter-service-monitor/0"} 38
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/streaming-service/streaming-service-monitor-jmx-exporter/0"} 953
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/streaming-service/streaming-service-monitor/0"} 150
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/vault-service/vault-service-monitor/0"} 271
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/zookeeper-service/zookeeper-service-monitor-jmx-exporter/0"} 307
prometheus_target_metadata_cache_entries{scrape_job="serviceMonitor/zookeeper-service/zookeeper-service-monitor/0"} 77
# HELP prometheus_target_scrape_pool_exceeded_label_limits_total Total number of times scrape pools hit the label limits, during sync or config reload.
# TYPE prometheus_target_scrape_pool_exceeded_label_limits_total counter
prometheus_target_scrape_pool_exceeded_label_limits_total 0
# HELP prometheus_target_scrape_pool_exceeded_target_limit_total Total number of times scrape pools hit the target limit, during sync or config reload.
# TYPE prometheus_target_scrape_pool_exceeded_target_limit_total counter
prometheus_target_scrape_pool_exceeded_target_limit_total 0
# HELP prometheus_target_scrape_pool_reloads_failed_total Total number of failed scrape pool reloads.
# TYPE prometheus_target_scrape_pool_reloads_failed_total counter
prometheus_target_scrape_pool_reloads_failed_total 0
# HELP prometheus_target_scrape_pool_reloads_total Total number of scrape pool reloads.
# TYPE prometheus_target_scrape_pool_reloads_total counter
prometheus_target_scrape_pool_reloads_total 0
# HELP prometheus_target_scrape_pool_sync_total Total number of syncs that were executed on a scrape pool.
# TYPE prometheus_target_scrape_pool_sync_total counter
prometheus_target_scrape_pool_sync_total{scrape_job="podMonitor/monitoring/monitoring-alertmanager-pod-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="podMonitor/monitoring/monitoring-grafana-operator-pod-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="podMonitor/monitoring/monitoring-grafana-pod-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="podMonitor/monitoring/monitoring-nginx-ingress-pod-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="podMonitor/monitoring/monitoring-prometheus-operator-pod-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="podMonitor/monitoring/monitoring-prometheus-pod-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="podMonitor/spark-operator-gcp/spark-operator-podmonitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="podMonitor/vault-service/vault-operator-pod-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="probe/monitoring-examples/blackbox-ingress-probe"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="probe/monitoring-examples/blackbox-static-probe"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/anko0714airflowbkp/airflow-anko0714airflowbkp/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/arango-db-a/arangodb-service-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/arangodb-graph-test/arangodb-service-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/arangodb/arangodb-service-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/consul-service/consul-service-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/core/cdc-control-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/core/cdc-stream-processor-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/core/cdc-streaming-service-monitor-jmx-exporter/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/core/cdc-test-api-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/core/cdc-test-db-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/core/cloud-streaming-service-monitor-jmx-exporter/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/core/mistral-service-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/default/arangodb-service-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/dev-helm-install/mistral-service-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/dev-mistral-demo/mistral-service-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/dev-mistral/mistral-service-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/domain-search/opensearch-service-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/elasticsearch-cluster/elasticsearch-service-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/elasticsearch-service/elasticsearch-service-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/envoy/monitoring-envoy/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/envoy/monitoring-keycloak/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/kafka-service/kafka-lag-exporter-service-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/kafka-service/kafka-service-monitor-jmx-exporter/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/kafka-service/kafka-service-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/mistral/mistral-service-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/mongodb/mongo-backup-exporter-service-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/mongodb/mongo-dbaas-exporter-service-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/mongodb/mongodb-prometheus-exporter-service-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/monitoring-examples/go-example-service-service-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/monitoring-examples/go-example-service-service-monitor/1"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/monitoring-examples/java-example-service-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/monitoring-examples/monitoring-example-service-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/monitoring/monitoring-cert-exporter/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/monitoring/monitoring-core-dns-service-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/monitoring/monitoring-etcd-service-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/monitoring/monitoring-kube-apiserver-service-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/monitoring/monitoring-kube-state-metrics/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/monitoring/monitoring-kube-state-metrics/1"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/monitoring/monitoring-kube-state-metrics/2"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/monitoring/monitoring-kubelet-service-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/monitoring/monitoring-kubelet-service-monitor/1"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/monitoring/monitoring-node-exporter/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/opensearch-service/opensearch-service-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/pg-metr/postgres-service-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/postgres-upgrade/clickhouse-backup-daemon-service-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/postgres-upgrade/clickhouse-service-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/profiler-cassandra/cassandra-exporter-service-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/profiler-test/esc-collector-service-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/profiler-test/esc-test-service-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/profiler-test/esc-ui-service-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/profiler/esc-collector-service-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/profiler/esc-test-service-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/profiler/esc-ui-service-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/rabbitmq/rabbitmq-per-object-service-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/rabbitmq/rabbitmq-service-monitor-old/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/rabbitmq/rabbitmq-service-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/redis/redis-prometheus-exporter-service-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/streaming-service/streaming-service-monitor-jmx-exporter/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/streaming-service/streaming-service-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/vault-service/vault-service-monitor/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/zookeeper-service/zookeeper-service-monitor-jmx-exporter/0"} 4856
prometheus_target_scrape_pool_sync_total{scrape_job="serviceMonitor/zookeeper-service/zookeeper-service-monitor/0"} 4856
# HELP prometheus_target_scrape_pool_targets Current number of targets in this scrape pool.
# TYPE prometheus_target_scrape_pool_targets gauge
prometheus_target_scrape_pool_targets{scrape_job="podMonitor/monitoring/monitoring-alertmanager-pod-monitor/0"} 1
prometheus_target_scrape_pool_targets{scrape_job="podMonitor/monitoring/monitoring-grafana-operator-pod-monitor/0"} 1
prometheus_target_scrape_pool_targets{scrape_job="podMonitor/monitoring/monitoring-grafana-pod-monitor/0"} 1
prometheus_target_scrape_pool_targets{scrape_job="podMonitor/monitoring/monitoring-nginx-ingress-pod-monitor/0"} 6
prometheus_target_scrape_pool_targets{scrape_job="podMonitor/monitoring/monitoring-prometheus-operator-pod-monitor/0"} 1
prometheus_target_scrape_pool_targets{scrape_job="podMonitor/monitoring/monitoring-prometheus-pod-monitor/0"} 1
prometheus_target_scrape_pool_targets{scrape_job="podMonitor/spark-operator-gcp/spark-operator-podmonitor/0"} 1
prometheus_target_scrape_pool_targets{scrape_job="podMonitor/vault-service/vault-operator-pod-monitor/0"} 1
prometheus_target_scrape_pool_targets{scrape_job="probe/monitoring-examples/blackbox-ingress-probe"} 0
prometheus_target_scrape_pool_targets{scrape_job="probe/monitoring-examples/blackbox-static-probe"} 2
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/anko0714airflowbkp/airflow-anko0714airflowbkp/0"} 0
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/arango-db-a/arangodb-service-monitor/0"} 0
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/arangodb-graph-test/arangodb-service-monitor/0"} 0
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/arangodb/arangodb-service-monitor/0"} 0
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/consul-service/consul-service-monitor/0"} 3
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/core/cdc-control-monitor/0"} 1
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/core/cdc-stream-processor-monitor/0"} 1
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/core/cdc-streaming-service-monitor-jmx-exporter/0"} 1
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/core/cdc-test-api-monitor/0"} 1
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/core/cdc-test-db-monitor/0"} 0
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/core/cloud-streaming-service-monitor-jmx-exporter/0"} 1
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/core/mistral-service-monitor/0"} 1
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/default/arangodb-service-monitor/0"} 0
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/dev-helm-install/mistral-service-monitor/0"} 0
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/dev-mistral-demo/mistral-service-monitor/0"} 1
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/dev-mistral/mistral-service-monitor/0"} 1
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/domain-search/opensearch-service-monitor/0"} 1
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/elasticsearch-cluster/elasticsearch-service-monitor/0"} 1
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/elasticsearch-service/elasticsearch-service-monitor/0"} 1
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/envoy/monitoring-envoy/0"} 0
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/envoy/monitoring-keycloak/0"} 0
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/kafka-service/kafka-lag-exporter-service-monitor/0"} 1
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/kafka-service/kafka-service-monitor-jmx-exporter/0"} 3
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/kafka-service/kafka-service-monitor/0"} 1
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/mistral/mistral-service-monitor/0"} 1
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/mongodb/mongo-backup-exporter-service-monitor/0"} 1
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/mongodb/mongo-dbaas-exporter-service-monitor/0"} 1
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/mongodb/mongodb-prometheus-exporter-service-monitor/0"} 5
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/monitoring-examples/go-example-service-service-monitor/0"} 0
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/monitoring-examples/go-example-service-service-monitor/1"} 0
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/monitoring-examples/java-example-service-monitor/0"} 0
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/monitoring-examples/monitoring-example-service-monitor/0"} 0
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/monitoring/monitoring-cert-exporter/0"} 1
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/monitoring/monitoring-core-dns-service-monitor/0"} 2
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/monitoring/monitoring-etcd-service-monitor/0"} 3
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/monitoring/monitoring-kube-apiserver-service-monitor/0"} 3
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/monitoring/monitoring-kube-state-metrics/0"} 1
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/monitoring/monitoring-kube-state-metrics/1"} 1
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/monitoring/monitoring-kube-state-metrics/2"} 1
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/monitoring/monitoring-kubelet-service-monitor/0"} 9
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/monitoring/monitoring-kubelet-service-monitor/1"} 9
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/monitoring/monitoring-node-exporter/0"} 9
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/opensearch-service/opensearch-service-monitor/0"} 1
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/pg-metr/postgres-service-monitor/0"} 1
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/postgres-upgrade/clickhouse-backup-daemon-service-monitor/0"} 0
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/postgres-upgrade/clickhouse-service-monitor/0"} 1
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/profiler-cassandra/cassandra-exporter-service-monitor/0"} 3
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/profiler-test/esc-collector-service-monitor/0"} 1
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/profiler-test/esc-test-service-monitor/0"} 1
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/profiler-test/esc-ui-service-monitor/0"} 1
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/profiler/esc-collector-service-monitor/0"} 1
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/profiler/esc-test-service-monitor/0"} 1
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/profiler/esc-ui-service-monitor/0"} 1
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/rabbitmq/rabbitmq-per-object-service-monitor/0"} 3
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/rabbitmq/rabbitmq-service-monitor-old/0"} 3
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/rabbitmq/rabbitmq-service-monitor/0"} 3
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/redis/redis-prometheus-exporter-service-monitor/0"} 1
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/streaming-service/streaming-service-monitor-jmx-exporter/0"} 3
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/streaming-service/streaming-service-monitor/0"} 1
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/vault-service/vault-service-monitor/0"} 4
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/zookeeper-service/zookeeper-service-monitor-jmx-exporter/0"} 3
prometheus_target_scrape_pool_targets{scrape_job="serviceMonitor/zookeeper-service/zookeeper-service-monitor/0"} 1
# HELP prometheus_target_scrape_pools_failed_total Total number of scrape pool creations that failed.
# TYPE prometheus_target_scrape_pools_failed_total counter
prometheus_target_scrape_pools_failed_total 0
# HELP prometheus_target_scrape_pools_total Total number of scrape pool creation attempts.
# TYPE prometheus_target_scrape_pools_total counter
prometheus_target_scrape_pools_total 72
# HELP prometheus_target_scrapes_cache_flush_forced_total How many times a scrape cache was flushed due to getting big while scrapes are failing.
# TYPE prometheus_target_scrapes_cache_flush_forced_total counter
prometheus_target_scrapes_cache_flush_forced_total 0
# HELP prometheus_target_scrapes_exceeded_body_size_limit_total Total number of scrapes that hit the body size limit
# TYPE prometheus_target_scrapes_exceeded_body_size_limit_total counter
prometheus_target_scrapes_exceeded_body_size_limit_total 0
# HELP prometheus_target_scrapes_exceeded_sample_limit_total Total number of scrapes that hit the sample limit and were rejected.
# TYPE prometheus_target_scrapes_exceeded_sample_limit_total counter
prometheus_target_scrapes_exceeded_sample_limit_total 0
# HELP prometheus_target_scrapes_exemplar_out_of_order_total Total number of exemplar rejected due to not being out of the expected order.
# TYPE prometheus_target_scrapes_exemplar_out_of_order_total counter
prometheus_target_scrapes_exemplar_out_of_order_total 0
# HELP prometheus_target_scrapes_sample_duplicate_timestamp_total Total number of samples rejected due to duplicate timestamps but different values.
# TYPE prometheus_target_scrapes_sample_duplicate_timestamp_total counter
prometheus_target_scrapes_sample_duplicate_timestamp_total 0
# HELP prometheus_target_scrapes_sample_out_of_bounds_total Total number of samples rejected due to timestamp falling outside of the time bounds.
# TYPE prometheus_target_scrapes_sample_out_of_bounds_total counter
prometheus_target_scrapes_sample_out_of_bounds_total 0
# HELP prometheus_target_scrapes_sample_out_of_order_total Total number of samples rejected due to not being out of the expected order.
# TYPE prometheus_target_scrapes_sample_out_of_order_total counter
prometheus_target_scrapes_sample_out_of_order_total 0
# HELP prometheus_target_sync_failed_total Total number of target sync failures.
# TYPE prometheus_target_sync_failed_total counter
prometheus_target_sync_failed_total{scrape_job="podMonitor/monitoring/monitoring-alertmanager-pod-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="podMonitor/monitoring/monitoring-grafana-operator-pod-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="podMonitor/monitoring/monitoring-grafana-pod-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="podMonitor/monitoring/monitoring-nginx-ingress-pod-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="podMonitor/monitoring/monitoring-prometheus-operator-pod-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="podMonitor/monitoring/monitoring-prometheus-pod-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="podMonitor/spark-operator-gcp/spark-operator-podmonitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="podMonitor/vault-service/vault-operator-pod-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="probe/monitoring-examples/blackbox-ingress-probe"} 0
prometheus_target_sync_failed_total{scrape_job="probe/monitoring-examples/blackbox-static-probe"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/anko0714airflowbkp/airflow-anko0714airflowbkp/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/arango-db-a/arangodb-service-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/arangodb-graph-test/arangodb-service-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/arangodb/arangodb-service-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/consul-service/consul-service-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/core/cdc-control-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/core/cdc-stream-processor-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/core/cdc-streaming-service-monitor-jmx-exporter/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/core/cdc-test-api-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/core/cdc-test-db-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/core/cloud-streaming-service-monitor-jmx-exporter/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/core/mistral-service-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/default/arangodb-service-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/dev-helm-install/mistral-service-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/dev-mistral-demo/mistral-service-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/dev-mistral/mistral-service-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/domain-search/opensearch-service-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/elasticsearch-cluster/elasticsearch-service-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/elasticsearch-service/elasticsearch-service-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/envoy/monitoring-envoy/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/envoy/monitoring-keycloak/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/kafka-service/kafka-lag-exporter-service-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/kafka-service/kafka-service-monitor-jmx-exporter/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/kafka-service/kafka-service-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/mistral/mistral-service-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/mongodb/mongo-backup-exporter-service-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/mongodb/mongo-dbaas-exporter-service-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/mongodb/mongodb-prometheus-exporter-service-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/monitoring-examples/go-example-service-service-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/monitoring-examples/go-example-service-service-monitor/1"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/monitoring-examples/java-example-service-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/monitoring-examples/monitoring-example-service-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/monitoring/monitoring-cert-exporter/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/monitoring/monitoring-core-dns-service-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/monitoring/monitoring-etcd-service-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/monitoring/monitoring-kube-apiserver-service-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/monitoring/monitoring-kube-state-metrics/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/monitoring/monitoring-kube-state-metrics/1"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/monitoring/monitoring-kube-state-metrics/2"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/monitoring/monitoring-kubelet-service-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/monitoring/monitoring-kubelet-service-monitor/1"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/monitoring/monitoring-node-exporter/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/opensearch-service/opensearch-service-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/pg-metr/postgres-service-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/postgres-upgrade/clickhouse-backup-daemon-service-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/postgres-upgrade/clickhouse-service-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/profiler-cassandra/cassandra-exporter-service-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/profiler-test/esc-collector-service-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/profiler-test/esc-test-service-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/profiler-test/esc-ui-service-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/profiler/esc-collector-service-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/profiler/esc-test-service-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/profiler/esc-ui-service-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/rabbitmq/rabbitmq-per-object-service-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/rabbitmq/rabbitmq-service-monitor-old/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/rabbitmq/rabbitmq-service-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/redis/redis-prometheus-exporter-service-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/streaming-service/streaming-service-monitor-jmx-exporter/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/streaming-service/streaming-service-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/vault-service/vault-service-monitor/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/zookeeper-service/zookeeper-service-monitor-jmx-exporter/0"} 0
prometheus_target_sync_failed_total{scrape_job="serviceMonitor/zookeeper-service/zookeeper-service-monitor/0"} 0
# HELP prometheus_target_sync_length_seconds Actual interval to sync the scrape pool.
# TYPE prometheus_target_sync_length_seconds summary
prometheus_target_sync_length_seconds{scrape_job="podMonitor/monitoring/monitoring-alertmanager-pod-monitor/0",quantile="0.01"} 0.000478394
prometheus_target_sync_length_seconds{scrape_job="podMonitor/monitoring/monitoring-alertmanager-pod-monitor/0",quantile="0.05"} 0.0004875
prometheus_target_sync_length_seconds{scrape_job="podMonitor/monitoring/monitoring-alertmanager-pod-monitor/0",quantile="0.5"} 0.000591915
prometheus_target_sync_length_seconds{scrape_job="podMonitor/monitoring/monitoring-alertmanager-pod-monitor/0",quantile="0.9"} 0.001435112
prometheus_target_sync_length_seconds{scrape_job="podMonitor/monitoring/monitoring-alertmanager-pod-monitor/0",quantile="0.99"} 0.004571312
prometheus_target_sync_length_seconds_sum{scrape_job="podMonitor/monitoring/monitoring-alertmanager-pod-monitor/0"} 4.248507244999982
prometheus_target_sync_length_seconds_count{scrape_job="podMonitor/monitoring/monitoring-alertmanager-pod-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="podMonitor/monitoring/monitoring-grafana-operator-pod-monitor/0",quantile="0.01"} 0.00046987
prometheus_target_sync_length_seconds{scrape_job="podMonitor/monitoring/monitoring-grafana-operator-pod-monitor/0",quantile="0.05"} 0.000502218
prometheus_target_sync_length_seconds{scrape_job="podMonitor/monitoring/monitoring-grafana-operator-pod-monitor/0",quantile="0.5"} 0.000570261
prometheus_target_sync_length_seconds{scrape_job="podMonitor/monitoring/monitoring-grafana-operator-pod-monitor/0",quantile="0.9"} 0.001223753
prometheus_target_sync_length_seconds{scrape_job="podMonitor/monitoring/monitoring-grafana-operator-pod-monitor/0",quantile="0.99"} 0.002011532
prometheus_target_sync_length_seconds_sum{scrape_job="podMonitor/monitoring/monitoring-grafana-operator-pod-monitor/0"} 3.566350142000012
prometheus_target_sync_length_seconds_count{scrape_job="podMonitor/monitoring/monitoring-grafana-operator-pod-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="podMonitor/monitoring/monitoring-grafana-pod-monitor/0",quantile="0.01"} 0.000480834
prometheus_target_sync_length_seconds{scrape_job="podMonitor/monitoring/monitoring-grafana-pod-monitor/0",quantile="0.05"} 0.000491202
prometheus_target_sync_length_seconds{scrape_job="podMonitor/monitoring/monitoring-grafana-pod-monitor/0",quantile="0.5"} 0.000559493
prometheus_target_sync_length_seconds{scrape_job="podMonitor/monitoring/monitoring-grafana-pod-monitor/0",quantile="0.9"} 0.001042526
prometheus_target_sync_length_seconds{scrape_job="podMonitor/monitoring/monitoring-grafana-pod-monitor/0",quantile="0.99"} 0.001668062
prometheus_target_sync_length_seconds_sum{scrape_job="podMonitor/monitoring/monitoring-grafana-pod-monitor/0"} 3.8678524220000035
prometheus_target_sync_length_seconds_count{scrape_job="podMonitor/monitoring/monitoring-grafana-pod-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="podMonitor/monitoring/monitoring-nginx-ingress-pod-monitor/0",quantile="0.01"} 0.01207253
prometheus_target_sync_length_seconds{scrape_job="podMonitor/monitoring/monitoring-nginx-ingress-pod-monitor/0",quantile="0.05"} 0.012281226
prometheus_target_sync_length_seconds{scrape_job="podMonitor/monitoring/monitoring-nginx-ingress-pod-monitor/0",quantile="0.5"} 0.015810135
prometheus_target_sync_length_seconds{scrape_job="podMonitor/monitoring/monitoring-nginx-ingress-pod-monitor/0",quantile="0.9"} 0.033147187
prometheus_target_sync_length_seconds{scrape_job="podMonitor/monitoring/monitoring-nginx-ingress-pod-monitor/0",quantile="0.99"} 0.044897557
prometheus_target_sync_length_seconds_sum{scrape_job="podMonitor/monitoring/monitoring-nginx-ingress-pod-monitor/0"} 100.61268193500011
prometheus_target_sync_length_seconds_count{scrape_job="podMonitor/monitoring/monitoring-nginx-ingress-pod-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="podMonitor/monitoring/monitoring-prometheus-operator-pod-monitor/0",quantile="0.01"} 0.000473378
prometheus_target_sync_length_seconds{scrape_job="podMonitor/monitoring/monitoring-prometheus-operator-pod-monitor/0",quantile="0.05"} 0.000485465
prometheus_target_sync_length_seconds{scrape_job="podMonitor/monitoring/monitoring-prometheus-operator-pod-monitor/0",quantile="0.5"} 0.000585082
prometheus_target_sync_length_seconds{scrape_job="podMonitor/monitoring/monitoring-prometheus-operator-pod-monitor/0",quantile="0.9"} 0.001061299
prometheus_target_sync_length_seconds{scrape_job="podMonitor/monitoring/monitoring-prometheus-operator-pod-monitor/0",quantile="0.99"} 0.002071929
prometheus_target_sync_length_seconds_sum{scrape_job="podMonitor/monitoring/monitoring-prometheus-operator-pod-monitor/0"} 3.660628124999989
prometheus_target_sync_length_seconds_count{scrape_job="podMonitor/monitoring/monitoring-prometheus-operator-pod-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="podMonitor/monitoring/monitoring-prometheus-pod-monitor/0",quantile="0.01"} 0.000462457
prometheus_target_sync_length_seconds{scrape_job="podMonitor/monitoring/monitoring-prometheus-pod-monitor/0",quantile="0.05"} 0.0004946
prometheus_target_sync_length_seconds{scrape_job="podMonitor/monitoring/monitoring-prometheus-pod-monitor/0",quantile="0.5"} 0.000563848
prometheus_target_sync_length_seconds{scrape_job="podMonitor/monitoring/monitoring-prometheus-pod-monitor/0",quantile="0.9"} 0.001012904
prometheus_target_sync_length_seconds{scrape_job="podMonitor/monitoring/monitoring-prometheus-pod-monitor/0",quantile="0.99"} 0.002326904
prometheus_target_sync_length_seconds_sum{scrape_job="podMonitor/monitoring/monitoring-prometheus-pod-monitor/0"} 4.088855916
prometheus_target_sync_length_seconds_count{scrape_job="podMonitor/monitoring/monitoring-prometheus-pod-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="podMonitor/spark-operator-gcp/spark-operator-podmonitor/0",quantile="0.01"} 0.000102681
prometheus_target_sync_length_seconds{scrape_job="podMonitor/spark-operator-gcp/spark-operator-podmonitor/0",quantile="0.05"} 0.000105296
prometheus_target_sync_length_seconds{scrape_job="podMonitor/spark-operator-gcp/spark-operator-podmonitor/0",quantile="0.5"} 0.000130066
prometheus_target_sync_length_seconds{scrape_job="podMonitor/spark-operator-gcp/spark-operator-podmonitor/0",quantile="0.9"} 0.000185886
prometheus_target_sync_length_seconds{scrape_job="podMonitor/spark-operator-gcp/spark-operator-podmonitor/0",quantile="0.99"} 0.000248187
prometheus_target_sync_length_seconds_sum{scrape_job="podMonitor/spark-operator-gcp/spark-operator-podmonitor/0"} 0.7775190260000037
prometheus_target_sync_length_seconds_count{scrape_job="podMonitor/spark-operator-gcp/spark-operator-podmonitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="podMonitor/vault-service/vault-operator-pod-monitor/0",quantile="0.01"} 0.00021489
prometheus_target_sync_length_seconds{scrape_job="podMonitor/vault-service/vault-operator-pod-monitor/0",quantile="0.05"} 0.000218761
prometheus_target_sync_length_seconds{scrape_job="podMonitor/vault-service/vault-operator-pod-monitor/0",quantile="0.5"} 0.000251168
prometheus_target_sync_length_seconds{scrape_job="podMonitor/vault-service/vault-operator-pod-monitor/0",quantile="0.9"} 0.000509995
prometheus_target_sync_length_seconds{scrape_job="podMonitor/vault-service/vault-operator-pod-monitor/0",quantile="0.99"} 0.001001754
prometheus_target_sync_length_seconds_sum{scrape_job="podMonitor/vault-service/vault-operator-pod-monitor/0"} 1.9917822390000075
prometheus_target_sync_length_seconds_count{scrape_job="podMonitor/vault-service/vault-operator-pod-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="probe/monitoring-examples/blackbox-ingress-probe",quantile="0.01"} 2.7632e-05
prometheus_target_sync_length_seconds{scrape_job="probe/monitoring-examples/blackbox-ingress-probe",quantile="0.05"} 2.8524e-05
prometheus_target_sync_length_seconds{scrape_job="probe/monitoring-examples/blackbox-ingress-probe",quantile="0.5"} 3.7808e-05
prometheus_target_sync_length_seconds{scrape_job="probe/monitoring-examples/blackbox-ingress-probe",quantile="0.9"} 6.7035e-05
prometheus_target_sync_length_seconds{scrape_job="probe/monitoring-examples/blackbox-ingress-probe",quantile="0.99"} 9.5256e-05
prometheus_target_sync_length_seconds_sum{scrape_job="probe/monitoring-examples/blackbox-ingress-probe"} 0.5018445669999995
prometheus_target_sync_length_seconds_count{scrape_job="probe/monitoring-examples/blackbox-ingress-probe"} 4856
prometheus_target_sync_length_seconds{scrape_job="probe/monitoring-examples/blackbox-static-probe",quantile="0.01"} 6.0641e-05
prometheus_target_sync_length_seconds{scrape_job="probe/monitoring-examples/blackbox-static-probe",quantile="0.05"} 6.174e-05
prometheus_target_sync_length_seconds{scrape_job="probe/monitoring-examples/blackbox-static-probe",quantile="0.5"} 7.6537e-05
prometheus_target_sync_length_seconds{scrape_job="probe/monitoring-examples/blackbox-static-probe",quantile="0.9"} 0.000166787
prometheus_target_sync_length_seconds{scrape_job="probe/monitoring-examples/blackbox-static-probe",quantile="0.99"} 0.000363752
prometheus_target_sync_length_seconds_sum{scrape_job="probe/monitoring-examples/blackbox-static-probe"} 0.45891808200000056
prometheus_target_sync_length_seconds_count{scrape_job="probe/monitoring-examples/blackbox-static-probe"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/anko0714airflowbkp/airflow-anko0714airflowbkp/0",quantile="0.01"} 5.232e-06
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/anko0714airflowbkp/airflow-anko0714airflowbkp/0",quantile="0.05"} 6.304e-06
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/anko0714airflowbkp/airflow-anko0714airflowbkp/0",quantile="0.5"} 8.87e-06
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/anko0714airflowbkp/airflow-anko0714airflowbkp/0",quantile="0.9"} 1.614e-05
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/anko0714airflowbkp/airflow-anko0714airflowbkp/0",quantile="0.99"} 3.79e-05
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/anko0714airflowbkp/airflow-anko0714airflowbkp/0"} 0.05361370900000013
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/anko0714airflowbkp/airflow-anko0714airflowbkp/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/arango-db-a/arangodb-service-monitor/0",quantile="0.01"} 6.538e-06
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/arango-db-a/arangodb-service-monitor/0",quantile="0.05"} 8.024e-06
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/arango-db-a/arangodb-service-monitor/0",quantile="0.5"} 1.2188e-05
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/arango-db-a/arangodb-service-monitor/0",quantile="0.9"} 2.1879e-05
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/arango-db-a/arangodb-service-monitor/0",quantile="0.99"} 7.638e-05
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/arango-db-a/arangodb-service-monitor/0"} 0.06778865499999981
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/arango-db-a/arangodb-service-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/arangodb-graph-test/arangodb-service-monitor/0",quantile="0.01"} 7.617e-06
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/arangodb-graph-test/arangodb-service-monitor/0",quantile="0.05"} 8.412e-06
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/arangodb-graph-test/arangodb-service-monitor/0",quantile="0.5"} 1.1112e-05
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/arangodb-graph-test/arangodb-service-monitor/0",quantile="0.9"} 1.6108e-05
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/arangodb-graph-test/arangodb-service-monitor/0",quantile="0.99"} 3.8364e-05
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/arangodb-graph-test/arangodb-service-monitor/0"} 0.0668932089999999
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/arangodb-graph-test/arangodb-service-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/arangodb/arangodb-service-monitor/0",quantile="0.01"} 0.000272077
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/arangodb/arangodb-service-monitor/0",quantile="0.05"} 0.000276573
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/arangodb/arangodb-service-monitor/0",quantile="0.5"} 0.000322863
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/arangodb/arangodb-service-monitor/0",quantile="0.9"} 0.000544532
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/arangodb/arangodb-service-monitor/0",quantile="0.99"} 0.001297346
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/arangodb/arangodb-service-monitor/0"} 2.2879743389999962
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/arangodb/arangodb-service-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/consul-service/consul-service-monitor/0",quantile="0.01"} 0.005311447
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/consul-service/consul-service-monitor/0",quantile="0.05"} 0.005358427
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/consul-service/consul-service-monitor/0",quantile="0.5"} 0.006100132
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/consul-service/consul-service-monitor/0",quantile="0.9"} 0.009174307
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/consul-service/consul-service-monitor/0",quantile="0.99"} 0.01271137
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/consul-service/consul-service-monitor/0"} 40.101817825999944
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/consul-service/consul-service-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/core/cdc-control-monitor/0",quantile="0.01"} 0.001629762
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/core/cdc-control-monitor/0",quantile="0.05"} 0.001634939
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/core/cdc-control-monitor/0",quantile="0.5"} 0.001887065
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/core/cdc-control-monitor/0",quantile="0.9"} 0.002655032
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/core/cdc-control-monitor/0",quantile="0.99"} 0.003515496
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/core/cdc-control-monitor/0"} 13.033920575999993
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/core/cdc-control-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/core/cdc-stream-processor-monitor/0",quantile="0.01"} 0.001647296
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/core/cdc-stream-processor-monitor/0",quantile="0.05"} 0.001666247
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/core/cdc-stream-processor-monitor/0",quantile="0.5"} 0.001940454
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/core/cdc-stream-processor-monitor/0",quantile="0.9"} 0.003145055
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/core/cdc-stream-processor-monitor/0",quantile="0.99"} 0.004482992
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/core/cdc-stream-processor-monitor/0"} 13.111877266000013
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/core/cdc-stream-processor-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/core/cdc-streaming-service-monitor-jmx-exporter/0",quantile="0.01"} 0.001638564
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/core/cdc-streaming-service-monitor-jmx-exporter/0",quantile="0.05"} 0.001657035
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/core/cdc-streaming-service-monitor-jmx-exporter/0",quantile="0.5"} 0.00191909
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/core/cdc-streaming-service-monitor-jmx-exporter/0",quantile="0.9"} 0.00307788
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/core/cdc-streaming-service-monitor-jmx-exporter/0",quantile="0.99"} 0.00478486
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/core/cdc-streaming-service-monitor-jmx-exporter/0"} 12.858484265999987
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/core/cdc-streaming-service-monitor-jmx-exporter/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/core/cdc-test-api-monitor/0",quantile="0.01"} 0.001627927
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/core/cdc-test-api-monitor/0",quantile="0.05"} 0.001639766
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/core/cdc-test-api-monitor/0",quantile="0.5"} 0.00197116
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/core/cdc-test-api-monitor/0",quantile="0.9"} 0.002621091
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/core/cdc-test-api-monitor/0",quantile="0.99"} 0.006462561
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/core/cdc-test-api-monitor/0"} 13.154285158000011
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/core/cdc-test-api-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/core/cdc-test-db-monitor/0",quantile="0.01"} 0.001551199
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/core/cdc-test-db-monitor/0",quantile="0.05"} 0.001573617
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/core/cdc-test-db-monitor/0",quantile="0.5"} 0.001856868
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/core/cdc-test-db-monitor/0",quantile="0.9"} 0.003037767
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/core/cdc-test-db-monitor/0",quantile="0.99"} 0.004899238
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/core/cdc-test-db-monitor/0"} 13.735398495000007
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/core/cdc-test-db-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/core/cloud-streaming-service-monitor-jmx-exporter/0",quantile="0.01"} 0.00164124
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/core/cloud-streaming-service-monitor-jmx-exporter/0",quantile="0.05"} 0.00168045
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/core/cloud-streaming-service-monitor-jmx-exporter/0",quantile="0.5"} 0.001837757
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/core/cloud-streaming-service-monitor-jmx-exporter/0",quantile="0.9"} 0.003561561
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/core/cloud-streaming-service-monitor-jmx-exporter/0",quantile="0.99"} 0.005446043
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/core/cloud-streaming-service-monitor-jmx-exporter/0"} 13.054817863000022
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/core/cloud-streaming-service-monitor-jmx-exporter/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/core/mistral-service-monitor/0",quantile="0.01"} 0.001592515
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/core/mistral-service-monitor/0",quantile="0.05"} 0.001628933
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/core/mistral-service-monitor/0",quantile="0.5"} 0.001854922
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/core/mistral-service-monitor/0",quantile="0.9"} 0.003067183
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/core/mistral-service-monitor/0",quantile="0.99"} 0.006439849
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/core/mistral-service-monitor/0"} 13.235468017999969
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/core/mistral-service-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/default/arangodb-service-monitor/0",quantile="0.01"} 4.2027e-05
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/default/arangodb-service-monitor/0",quantile="0.05"} 4.3535e-05
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/default/arangodb-service-monitor/0",quantile="0.5"} 5.2744e-05
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/default/arangodb-service-monitor/0",quantile="0.9"} 0.00010349
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/default/arangodb-service-monitor/0",quantile="0.99"} 0.000138272
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/default/arangodb-service-monitor/0"} 0.5609047000000017
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/default/arangodb-service-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/dev-helm-install/mistral-service-monitor/0",quantile="0.01"} 0.000189645
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/dev-helm-install/mistral-service-monitor/0",quantile="0.05"} 0.000194057
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/dev-helm-install/mistral-service-monitor/0",quantile="0.5"} 0.00022645
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/dev-helm-install/mistral-service-monitor/0",quantile="0.9"} 0.000441465
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/dev-helm-install/mistral-service-monitor/0",quantile="0.99"} 0.000918007
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/dev-helm-install/mistral-service-monitor/0"} 1.4724237409999994
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/dev-helm-install/mistral-service-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/dev-mistral-demo/mistral-service-monitor/0",quantile="0.01"} 0.000157137
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/dev-mistral-demo/mistral-service-monitor/0",quantile="0.05"} 0.000160755
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/dev-mistral-demo/mistral-service-monitor/0",quantile="0.5"} 0.000187936
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/dev-mistral-demo/mistral-service-monitor/0",quantile="0.9"} 0.000324403
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/dev-mistral-demo/mistral-service-monitor/0",quantile="0.99"} 0.000500688
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/dev-mistral-demo/mistral-service-monitor/0"} 1.404593423000001
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/dev-mistral-demo/mistral-service-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/dev-mistral/mistral-service-monitor/0",quantile="0.01"} 0.000155141
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/dev-mistral/mistral-service-monitor/0",quantile="0.05"} 0.000163861
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/dev-mistral/mistral-service-monitor/0",quantile="0.5"} 0.000189181
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/dev-mistral/mistral-service-monitor/0",quantile="0.9"} 0.000337919
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/dev-mistral/mistral-service-monitor/0",quantile="0.99"} 0.000849282
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/dev-mistral/mistral-service-monitor/0"} 1.5090701700000018
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/dev-mistral/mistral-service-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/domain-search/opensearch-service-monitor/0",quantile="0.01"} 0.00029543
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/domain-search/opensearch-service-monitor/0",quantile="0.05"} 0.000300471
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/domain-search/opensearch-service-monitor/0",quantile="0.5"} 0.000335916
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/domain-search/opensearch-service-monitor/0",quantile="0.9"} 0.000545944
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/domain-search/opensearch-service-monitor/0",quantile="0.99"} 0.001322947
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/domain-search/opensearch-service-monitor/0"} 2.1872240069999997
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/domain-search/opensearch-service-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/elasticsearch-cluster/elasticsearch-service-monitor/0",quantile="0.01"} 0.001450387
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/elasticsearch-cluster/elasticsearch-service-monitor/0",quantile="0.05"} 0.001468103
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/elasticsearch-cluster/elasticsearch-service-monitor/0",quantile="0.5"} 0.001762719
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/elasticsearch-cluster/elasticsearch-service-monitor/0",quantile="0.9"} 0.003420111
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/elasticsearch-cluster/elasticsearch-service-monitor/0",quantile="0.99"} 0.004747689
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/elasticsearch-cluster/elasticsearch-service-monitor/0"} 10.809264451000017
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/elasticsearch-cluster/elasticsearch-service-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/elasticsearch-service/elasticsearch-service-monitor/0",quantile="0.01"} 0.000657921
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/elasticsearch-service/elasticsearch-service-monitor/0",quantile="0.05"} 0.000668411
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/elasticsearch-service/elasticsearch-service-monitor/0",quantile="0.5"} 0.000743501
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/elasticsearch-service/elasticsearch-service-monitor/0",quantile="0.9"} 0.001207929
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/elasticsearch-service/elasticsearch-service-monitor/0",quantile="0.99"} 0.001292946
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/elasticsearch-service/elasticsearch-service-monitor/0"} 5.781920525999996
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/elasticsearch-service/elasticsearch-service-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/envoy/monitoring-envoy/0",quantile="0.01"} 5.6873e-05
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/envoy/monitoring-envoy/0",quantile="0.05"} 5.9744e-05
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/envoy/monitoring-envoy/0",quantile="0.5"} 7.0543e-05
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/envoy/monitoring-envoy/0",quantile="0.9"} 0.000111803
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/envoy/monitoring-envoy/0",quantile="0.99"} 0.000261867
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/envoy/monitoring-envoy/0"} 0.9140477089999997
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/envoy/monitoring-envoy/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/envoy/monitoring-keycloak/0",quantile="0.01"} 5.561e-05
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/envoy/monitoring-keycloak/0",quantile="0.05"} 5.9233e-05
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/envoy/monitoring-keycloak/0",quantile="0.5"} 7.118e-05
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/envoy/monitoring-keycloak/0",quantile="0.9"} 0.000130075
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/envoy/monitoring-keycloak/0",quantile="0.99"} 0.000264035
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/envoy/monitoring-keycloak/0"} 0.8167953069999995
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/envoy/monitoring-keycloak/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/kafka-service/kafka-lag-exporter-service-monitor/0",quantile="0.01"} 0.001936797
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/kafka-service/kafka-lag-exporter-service-monitor/0",quantile="0.05"} 0.001947707
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/kafka-service/kafka-lag-exporter-service-monitor/0",quantile="0.5"} 0.002212634
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/kafka-service/kafka-lag-exporter-service-monitor/0",quantile="0.9"} 0.004118746
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/kafka-service/kafka-lag-exporter-service-monitor/0",quantile="0.99"} 0.007346671
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/kafka-service/kafka-lag-exporter-service-monitor/0"} 14.467423377999962
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/kafka-service/kafka-lag-exporter-service-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/kafka-service/kafka-service-monitor-jmx-exporter/0",quantile="0.01"} 0.00202804
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/kafka-service/kafka-service-monitor-jmx-exporter/0",quantile="0.05"} 0.00209348
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/kafka-service/kafka-service-monitor-jmx-exporter/0",quantile="0.5"} 0.002336979
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/kafka-service/kafka-service-monitor-jmx-exporter/0",quantile="0.9"} 0.003840073
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/kafka-service/kafka-service-monitor-jmx-exporter/0",quantile="0.99"} 0.004831824
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/kafka-service/kafka-service-monitor-jmx-exporter/0"} 16.44201469000003
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/kafka-service/kafka-service-monitor-jmx-exporter/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/kafka-service/kafka-service-monitor/0",quantile="0.01"} 0.001913026
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/kafka-service/kafka-service-monitor/0",quantile="0.05"} 0.001941557
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/kafka-service/kafka-service-monitor/0",quantile="0.5"} 0.002210508
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/kafka-service/kafka-service-monitor/0",quantile="0.9"} 0.003682428
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/kafka-service/kafka-service-monitor/0",quantile="0.99"} 0.006919145
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/kafka-service/kafka-service-monitor/0"} 14.253522554999986
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/kafka-service/kafka-service-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/mistral/mistral-service-monitor/0",quantile="0.01"} 0.000153102
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/mistral/mistral-service-monitor/0",quantile="0.05"} 0.00016014
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/mistral/mistral-service-monitor/0",quantile="0.5"} 0.000187498
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/mistral/mistral-service-monitor/0",quantile="0.9"} 0.000353892
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/mistral/mistral-service-monitor/0",quantile="0.99"} 0.00061555
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/mistral/mistral-service-monitor/0"} 1.3540015120000035
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/mistral/mistral-service-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/mongodb/mongo-backup-exporter-service-monitor/0",quantile="0.01"} 0.000255455
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/mongodb/mongo-backup-exporter-service-monitor/0",quantile="0.05"} 0.000258978
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/mongodb/mongo-backup-exporter-service-monitor/0",quantile="0.5"} 0.000304557
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/mongodb/mongo-backup-exporter-service-monitor/0",quantile="0.9"} 0.000599796
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/mongodb/mongo-backup-exporter-service-monitor/0",quantile="0.99"} 0.00086663
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/mongodb/mongo-backup-exporter-service-monitor/0"} 2.0929943230000014
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/mongodb/mongo-backup-exporter-service-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/mongodb/mongo-dbaas-exporter-service-monitor/0",quantile="0.01"} 0.000229466
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/mongodb/mongo-dbaas-exporter-service-monitor/0",quantile="0.05"} 0.000254034
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/mongodb/mongo-dbaas-exporter-service-monitor/0",quantile="0.5"} 0.000332418
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/mongodb/mongo-dbaas-exporter-service-monitor/0",quantile="0.9"} 0.000597611
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/mongodb/mongo-dbaas-exporter-service-monitor/0",quantile="0.99"} 0.001070133
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/mongodb/mongo-dbaas-exporter-service-monitor/0"} 2.204699052999998
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/mongodb/mongo-dbaas-exporter-service-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/mongodb/mongodb-prometheus-exporter-service-monitor/0",quantile="0.01"} 0.000391549
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/mongodb/mongodb-prometheus-exporter-service-monitor/0",quantile="0.05"} 0.000396667
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/mongodb/mongodb-prometheus-exporter-service-monitor/0",quantile="0.5"} 0.000466785
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/mongodb/mongodb-prometheus-exporter-service-monitor/0",quantile="0.9"} 0.000929338
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/mongodb/mongodb-prometheus-exporter-service-monitor/0",quantile="0.99"} 0.001711006
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/mongodb/mongodb-prometheus-exporter-service-monitor/0"} 3.2805540889999945
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/mongodb/mongodb-prometheus-exporter-service-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring-examples/go-example-service-service-monitor/0",quantile="0.01"} 5.914e-06
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring-examples/go-example-service-service-monitor/0",quantile="0.05"} 6.47e-06
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring-examples/go-example-service-service-monitor/0",quantile="0.5"} 8.766e-06
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring-examples/go-example-service-service-monitor/0",quantile="0.9"} 1.3216e-05
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring-examples/go-example-service-service-monitor/0",quantile="0.99"} 3.067e-05
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/monitoring-examples/go-example-service-service-monitor/0"} 0.055482270999999986
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/monitoring-examples/go-example-service-service-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring-examples/go-example-service-service-monitor/1",quantile="0.01"} 5.489e-06
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring-examples/go-example-service-service-monitor/1",quantile="0.05"} 6.183e-06
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring-examples/go-example-service-service-monitor/1",quantile="0.5"} 9.513e-06
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring-examples/go-example-service-service-monitor/1",quantile="0.9"} 1.7886e-05
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring-examples/go-example-service-service-monitor/1",quantile="0.99"} 3.3249e-05
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/monitoring-examples/go-example-service-service-monitor/1"} 0.058730856999999866
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/monitoring-examples/go-example-service-service-monitor/1"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring-examples/java-example-service-monitor/0",quantile="0.01"} 5.743e-06
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring-examples/java-example-service-monitor/0",quantile="0.05"} 6.176e-06
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring-examples/java-example-service-monitor/0",quantile="0.5"} 9.233e-06
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring-examples/java-example-service-monitor/0",quantile="0.9"} 1.5979e-05
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring-examples/java-example-service-monitor/0",quantile="0.99"} 4.24e-05
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/monitoring-examples/java-example-service-monitor/0"} 0.053928117999999976
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/monitoring-examples/java-example-service-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring-examples/monitoring-example-service-monitor/0",quantile="0.01"} 5.488e-06
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring-examples/monitoring-example-service-monitor/0",quantile="0.05"} 5.521e-06
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring-examples/monitoring-example-service-monitor/0",quantile="0.5"} 8.761e-06
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring-examples/monitoring-example-service-monitor/0",quantile="0.9"} 1.64e-05
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring-examples/monitoring-example-service-monitor/0",quantile="0.99"} 3.091e-05
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/monitoring-examples/monitoring-example-service-monitor/0"} 0.053624127999999945
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/monitoring-examples/monitoring-example-service-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring/monitoring-cert-exporter/0",quantile="0.01"} 0.000905937
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring/monitoring-cert-exporter/0",quantile="0.05"} 0.000946646
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring/monitoring-cert-exporter/0",quantile="0.5"} 0.00107334
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring/monitoring-cert-exporter/0",quantile="0.9"} 0.002098538
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring/monitoring-cert-exporter/0",quantile="0.99"} 0.003570429
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/monitoring/monitoring-cert-exporter/0"} 8.043134464999973
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/monitoring/monitoring-cert-exporter/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring/monitoring-core-dns-service-monitor/0",quantile="0.01"} 0.000294241
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring/monitoring-core-dns-service-monitor/0",quantile="0.05"} 0.000317041
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring/monitoring-core-dns-service-monitor/0",quantile="0.5"} 0.000363809
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring/monitoring-core-dns-service-monitor/0",quantile="0.9"} 0.000787971
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring/monitoring-core-dns-service-monitor/0",quantile="0.99"} 0.001534544
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/monitoring/monitoring-core-dns-service-monitor/0"} 2.3516832030000008
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/monitoring/monitoring-core-dns-service-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring/monitoring-etcd-service-monitor/0",quantile="0.01"} 0.000338985
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring/monitoring-etcd-service-monitor/0",quantile="0.05"} 0.000345221
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring/monitoring-etcd-service-monitor/0",quantile="0.5"} 0.000392073
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring/monitoring-etcd-service-monitor/0",quantile="0.9"} 0.000697672
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring/monitoring-etcd-service-monitor/0",quantile="0.99"} 0.000936188
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/monitoring/monitoring-etcd-service-monitor/0"} 2.6201977329999933
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/monitoring/monitoring-etcd-service-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring/monitoring-kube-apiserver-service-monitor/0",quantile="0.01"} 0.000126909
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring/monitoring-kube-apiserver-service-monitor/0",quantile="0.05"} 0.000128822
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring/monitoring-kube-apiserver-service-monitor/0",quantile="0.5"} 0.000157159
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring/monitoring-kube-apiserver-service-monitor/0",quantile="0.9"} 0.000263725
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring/monitoring-kube-apiserver-service-monitor/0",quantile="0.99"} 0.00074259
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/monitoring/monitoring-kube-apiserver-service-monitor/0"} 0.9873157450000023
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/monitoring/monitoring-kube-apiserver-service-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring/monitoring-kube-state-metrics/0",quantile="0.01"} 0.00091894
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring/monitoring-kube-state-metrics/0",quantile="0.05"} 0.000943349
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring/monitoring-kube-state-metrics/0",quantile="0.5"} 0.001119356
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring/monitoring-kube-state-metrics/0",quantile="0.9"} 0.002145136
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring/monitoring-kube-state-metrics/0",quantile="0.99"} 0.003992436
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/monitoring/monitoring-kube-state-metrics/0"} 7.825166292000034
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/monitoring/monitoring-kube-state-metrics/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring/monitoring-kube-state-metrics/1",quantile="0.01"} 0.000890396
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring/monitoring-kube-state-metrics/1",quantile="0.05"} 0.000920505
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring/monitoring-kube-state-metrics/1",quantile="0.5"} 0.001063892
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring/monitoring-kube-state-metrics/1",quantile="0.9"} 0.001708559
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring/monitoring-kube-state-metrics/1",quantile="0.99"} 0.002373427
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/monitoring/monitoring-kube-state-metrics/1"} 6.891797427999996
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/monitoring/monitoring-kube-state-metrics/1"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring/monitoring-kube-state-metrics/2",quantile="0.01"} 0.00093833
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring/monitoring-kube-state-metrics/2",quantile="0.05"} 0.000948207
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring/monitoring-kube-state-metrics/2",quantile="0.5"} 0.001119391
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring/monitoring-kube-state-metrics/2",quantile="0.9"} 0.002219745
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring/monitoring-kube-state-metrics/2",quantile="0.99"} 0.004177795
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/monitoring/monitoring-kube-state-metrics/2"} 7.911408235000012
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/monitoring/monitoring-kube-state-metrics/2"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring/monitoring-kubelet-service-monitor/0",quantile="0.01"} 0.001145175
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring/monitoring-kubelet-service-monitor/0",quantile="0.05"} 0.001161143
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring/monitoring-kubelet-service-monitor/0",quantile="0.5"} 0.001307908
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring/monitoring-kubelet-service-monitor/0",quantile="0.9"} 0.002347787
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring/monitoring-kubelet-service-monitor/0",quantile="0.99"} 0.002902761
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/monitoring/monitoring-kubelet-service-monitor/0"} 8.910772621000001
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/monitoring/monitoring-kubelet-service-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring/monitoring-kubelet-service-monitor/1",quantile="0.01"} 0.00117235
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring/monitoring-kubelet-service-monitor/1",quantile="0.05"} 0.001178484
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring/monitoring-kubelet-service-monitor/1",quantile="0.5"} 0.001306261
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring/monitoring-kubelet-service-monitor/1",quantile="0.9"} 0.002092551
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring/monitoring-kubelet-service-monitor/1",quantile="0.99"} 0.004732772
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/monitoring/monitoring-kubelet-service-monitor/1"} 9.406532065000002
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/monitoring/monitoring-kubelet-service-monitor/1"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring/monitoring-node-exporter/0",quantile="0.01"} 0.001239099
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring/monitoring-node-exporter/0",quantile="0.05"} 0.001264822
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring/monitoring-node-exporter/0",quantile="0.5"} 0.001454474
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring/monitoring-node-exporter/0",quantile="0.9"} 0.002772277
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/monitoring/monitoring-node-exporter/0",quantile="0.99"} 0.003961403
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/monitoring/monitoring-node-exporter/0"} 9.10701441000001
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/monitoring/monitoring-node-exporter/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/opensearch-service/opensearch-service-monitor/0",quantile="0.01"} 0.001842619
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/opensearch-service/opensearch-service-monitor/0",quantile="0.05"} 0.001850347
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/opensearch-service/opensearch-service-monitor/0",quantile="0.5"} 0.002182367
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/opensearch-service/opensearch-service-monitor/0",quantile="0.9"} 0.004049521
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/opensearch-service/opensearch-service-monitor/0",quantile="0.99"} 0.005522806
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/opensearch-service/opensearch-service-monitor/0"} 13.130813329999949
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/opensearch-service/opensearch-service-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/pg-metr/postgres-service-monitor/0",quantile="0.01"} 0.00033742
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/pg-metr/postgres-service-monitor/0",quantile="0.05"} 0.000341171
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/pg-metr/postgres-service-monitor/0",quantile="0.5"} 0.000387811
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/pg-metr/postgres-service-monitor/0",quantile="0.9"} 0.000722043
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/pg-metr/postgres-service-monitor/0",quantile="0.99"} 0.001133087
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/pg-metr/postgres-service-monitor/0"} 2.937938130999996
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/pg-metr/postgres-service-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/postgres-upgrade/clickhouse-backup-daemon-service-monitor/0",quantile="0.01"} 4.3197e-05
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/postgres-upgrade/clickhouse-backup-daemon-service-monitor/0",quantile="0.05"} 4.3937e-05
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/postgres-upgrade/clickhouse-backup-daemon-service-monitor/0",quantile="0.5"} 5.4217e-05
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/postgres-upgrade/clickhouse-backup-daemon-service-monitor/0",quantile="0.9"} 9.2495e-05
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/postgres-upgrade/clickhouse-backup-daemon-service-monitor/0",quantile="0.99"} 0.000200568
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/postgres-upgrade/clickhouse-backup-daemon-service-monitor/0"} 0.431835727999997
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/postgres-upgrade/clickhouse-backup-daemon-service-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/postgres-upgrade/clickhouse-service-monitor/0",quantile="0.01"} 9.0801e-05
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/postgres-upgrade/clickhouse-service-monitor/0",quantile="0.05"} 9.4869e-05
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/postgres-upgrade/clickhouse-service-monitor/0",quantile="0.5"} 0.0001178
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/postgres-upgrade/clickhouse-service-monitor/0",quantile="0.9"} 0.000200512
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/postgres-upgrade/clickhouse-service-monitor/0",quantile="0.99"} 0.000263936
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/postgres-upgrade/clickhouse-service-monitor/0"} 0.8251322040000002
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/postgres-upgrade/clickhouse-service-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/profiler-cassandra/cassandra-exporter-service-monitor/0",quantile="0.01"} 0.001576188
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/profiler-cassandra/cassandra-exporter-service-monitor/0",quantile="0.05"} 0.001609517
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/profiler-cassandra/cassandra-exporter-service-monitor/0",quantile="0.5"} 0.001868444
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/profiler-cassandra/cassandra-exporter-service-monitor/0",quantile="0.9"} 0.00299062
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/profiler-cassandra/cassandra-exporter-service-monitor/0",quantile="0.99"} 0.006206651
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/profiler-cassandra/cassandra-exporter-service-monitor/0"} 12.42091707899997
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/profiler-cassandra/cassandra-exporter-service-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/profiler-test/esc-collector-service-monitor/0",quantile="0.01"} 0.000219238
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/profiler-test/esc-collector-service-monitor/0",quantile="0.05"} 0.000235353
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/profiler-test/esc-collector-service-monitor/0",quantile="0.5"} 0.000291355
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/profiler-test/esc-collector-service-monitor/0",quantile="0.9"} 0.000550685
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/profiler-test/esc-collector-service-monitor/0",quantile="0.99"} 0.001037912
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/profiler-test/esc-collector-service-monitor/0"} 1.7423707749999977
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/profiler-test/esc-collector-service-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/profiler-test/esc-test-service-monitor/0",quantile="0.01"} 0.00024284
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/profiler-test/esc-test-service-monitor/0",quantile="0.05"} 0.000244121
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/profiler-test/esc-test-service-monitor/0",quantile="0.5"} 0.000287298
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/profiler-test/esc-test-service-monitor/0",quantile="0.9"} 0.000482574
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/profiler-test/esc-test-service-monitor/0",quantile="0.99"} 0.00053525
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/profiler-test/esc-test-service-monitor/0"} 1.9087618999999982
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/profiler-test/esc-test-service-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/profiler-test/esc-ui-service-monitor/0",quantile="0.01"} 0.000232176
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/profiler-test/esc-ui-service-monitor/0",quantile="0.05"} 0.000239561
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/profiler-test/esc-ui-service-monitor/0",quantile="0.5"} 0.000275902
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/profiler-test/esc-ui-service-monitor/0",quantile="0.9"} 0.000515432
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/profiler-test/esc-ui-service-monitor/0",quantile="0.99"} 0.000783985
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/profiler-test/esc-ui-service-monitor/0"} 1.8061487619999974
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/profiler-test/esc-ui-service-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/profiler/esc-collector-service-monitor/0",quantile="0.01"} 0.000225117
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/profiler/esc-collector-service-monitor/0",quantile="0.05"} 0.000235426
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/profiler/esc-collector-service-monitor/0",quantile="0.5"} 0.000303163
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/profiler/esc-collector-service-monitor/0",quantile="0.9"} 0.000556708
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/profiler/esc-collector-service-monitor/0",quantile="0.99"} 0.000839055
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/profiler/esc-collector-service-monitor/0"} 2.1593644690000082
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/profiler/esc-collector-service-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/profiler/esc-test-service-monitor/0",quantile="0.01"} 0.000232162
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/profiler/esc-test-service-monitor/0",quantile="0.05"} 0.000236436
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/profiler/esc-test-service-monitor/0",quantile="0.5"} 0.000266392
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/profiler/esc-test-service-monitor/0",quantile="0.9"} 0.000458488
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/profiler/esc-test-service-monitor/0",quantile="0.99"} 0.001085369
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/profiler/esc-test-service-monitor/0"} 1.811757474999998
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/profiler/esc-test-service-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/profiler/esc-ui-service-monitor/0",quantile="0.01"} 0.000230681
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/profiler/esc-ui-service-monitor/0",quantile="0.05"} 0.000234531
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/profiler/esc-ui-service-monitor/0",quantile="0.5"} 0.000275282
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/profiler/esc-ui-service-monitor/0",quantile="0.9"} 0.000582265
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/profiler/esc-ui-service-monitor/0",quantile="0.99"} 0.000802806
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/profiler/esc-ui-service-monitor/0"} 1.9110734129999993
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/profiler/esc-ui-service-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/rabbitmq/rabbitmq-per-object-service-monitor/0",quantile="0.01"} 0.001808076
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/rabbitmq/rabbitmq-per-object-service-monitor/0",quantile="0.05"} 0.001826941
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/rabbitmq/rabbitmq-per-object-service-monitor/0",quantile="0.5"} 0.002099063
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/rabbitmq/rabbitmq-per-object-service-monitor/0",quantile="0.9"} 0.00320976
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/rabbitmq/rabbitmq-per-object-service-monitor/0",quantile="0.99"} 0.005932554
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/rabbitmq/rabbitmq-per-object-service-monitor/0"} 14.088320785999969
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/rabbitmq/rabbitmq-per-object-service-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/rabbitmq/rabbitmq-service-monitor-old/0",quantile="0.01"} 0.001805863
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/rabbitmq/rabbitmq-service-monitor-old/0",quantile="0.05"} 0.001823064
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/rabbitmq/rabbitmq-service-monitor-old/0",quantile="0.5"} 0.002193137
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/rabbitmq/rabbitmq-service-monitor-old/0",quantile="0.9"} 0.005082731
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/rabbitmq/rabbitmq-service-monitor-old/0",quantile="0.99"} 0.006331504
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/rabbitmq/rabbitmq-service-monitor-old/0"} 14.023586161000003
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/rabbitmq/rabbitmq-service-monitor-old/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/rabbitmq/rabbitmq-service-monitor/0",quantile="0.01"} 0.001810987
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/rabbitmq/rabbitmq-service-monitor/0",quantile="0.05"} 0.001843473
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/rabbitmq/rabbitmq-service-monitor/0",quantile="0.5"} 0.002335543
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/rabbitmq/rabbitmq-service-monitor/0",quantile="0.9"} 0.003318585
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/rabbitmq/rabbitmq-service-monitor/0",quantile="0.99"} 0.00499894
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/rabbitmq/rabbitmq-service-monitor/0"} 14.132577602000016
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/rabbitmq/rabbitmq-service-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/redis/redis-prometheus-exporter-service-monitor/0",quantile="0.01"} 0.000173236
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/redis/redis-prometheus-exporter-service-monitor/0",quantile="0.05"} 0.000174362
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/redis/redis-prometheus-exporter-service-monitor/0",quantile="0.5"} 0.000239822
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/redis/redis-prometheus-exporter-service-monitor/0",quantile="0.9"} 0.0004626
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/redis/redis-prometheus-exporter-service-monitor/0",quantile="0.99"} 0.00074719
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/redis/redis-prometheus-exporter-service-monitor/0"} 1.7029760710000035
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/redis/redis-prometheus-exporter-service-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/streaming-service/streaming-service-monitor-jmx-exporter/0",quantile="0.01"} 0.000494362
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/streaming-service/streaming-service-monitor-jmx-exporter/0",quantile="0.05"} 0.000498822
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/streaming-service/streaming-service-monitor-jmx-exporter/0",quantile="0.5"} 0.000552546
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/streaming-service/streaming-service-monitor-jmx-exporter/0",quantile="0.9"} 0.000937077
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/streaming-service/streaming-service-monitor-jmx-exporter/0",quantile="0.99"} 0.002203832
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/streaming-service/streaming-service-monitor-jmx-exporter/0"} 3.830265409000001
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/streaming-service/streaming-service-monitor-jmx-exporter/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/streaming-service/streaming-service-monitor/0",quantile="0.01"} 0.000387545
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/streaming-service/streaming-service-monitor/0",quantile="0.05"} 0.000393465
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/streaming-service/streaming-service-monitor/0",quantile="0.5"} 0.000446366
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/streaming-service/streaming-service-monitor/0",quantile="0.9"} 0.000738123
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/streaming-service/streaming-service-monitor/0",quantile="0.99"} 0.001379595
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/streaming-service/streaming-service-monitor/0"} 3.325650028999999
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/streaming-service/streaming-service-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/vault-service/vault-service-monitor/0",quantile="0.01"} 0.000757234
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/vault-service/vault-service-monitor/0",quantile="0.05"} 0.000772424
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/vault-service/vault-service-monitor/0",quantile="0.5"} 0.000839064
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/vault-service/vault-service-monitor/0",quantile="0.9"} 0.001390742
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/vault-service/vault-service-monitor/0",quantile="0.99"} 0.002725271
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/vault-service/vault-service-monitor/0"} 5.902253167000005
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/vault-service/vault-service-monitor/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/zookeeper-service/zookeeper-service-monitor-jmx-exporter/0",quantile="0.01"} 0.001561905
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/zookeeper-service/zookeeper-service-monitor-jmx-exporter/0",quantile="0.05"} 0.00164344
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/zookeeper-service/zookeeper-service-monitor-jmx-exporter/0",quantile="0.5"} 0.001866401
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/zookeeper-service/zookeeper-service-monitor-jmx-exporter/0",quantile="0.9"} 0.003348544
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/zookeeper-service/zookeeper-service-monitor-jmx-exporter/0",quantile="0.99"} 0.005122901
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/zookeeper-service/zookeeper-service-monitor-jmx-exporter/0"} 12.529299536000027
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/zookeeper-service/zookeeper-service-monitor-jmx-exporter/0"} 4856
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/zookeeper-service/zookeeper-service-monitor/0",quantile="0.01"} 0.001440073
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/zookeeper-service/zookeeper-service-monitor/0",quantile="0.05"} 0.001446722
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/zookeeper-service/zookeeper-service-monitor/0",quantile="0.5"} 0.001656306
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/zookeeper-service/zookeeper-service-monitor/0",quantile="0.9"} 0.003027863
prometheus_target_sync_length_seconds{scrape_job="serviceMonitor/zookeeper-service/zookeeper-service-monitor/0",quantile="0.99"} 0.004579822
prometheus_target_sync_length_seconds_sum{scrape_job="serviceMonitor/zookeeper-service/zookeeper-service-monitor/0"} 11.92986080900003
prometheus_target_sync_length_seconds_count{scrape_job="serviceMonitor/zookeeper-service/zookeeper-service-monitor/0"} 4856
# HELP prometheus_template_text_expansion_failures_total The total number of template text expansion failures.
# TYPE prometheus_template_text_expansion_failures_total counter
prometheus_template_text_expansion_failures_total 0
# HELP prometheus_template_text_expansions_total The total number of template text expansions.
# TYPE prometheus_template_text_expansions_total counter
prometheus_template_text_expansions_total 619133
# HELP prometheus_treecache_watcher_goroutines The current number of watcher goroutines.
# TYPE prometheus_treecache_watcher_goroutines gauge
prometheus_treecache_watcher_goroutines 0
# HELP prometheus_treecache_zookeeper_failures_total The total number of ZooKeeper failures.
# TYPE prometheus_treecache_zookeeper_failures_total counter
prometheus_treecache_zookeeper_failures_total 0
# HELP prometheus_tsdb_blocks_loaded Number of currently loaded data blocks
# TYPE prometheus_tsdb_blocks_loaded gauge
prometheus_tsdb_blocks_loaded 6
# HELP prometheus_tsdb_checkpoint_creations_failed_total Total number of checkpoint creations that failed.
# TYPE prometheus_tsdb_checkpoint_creations_failed_total counter
prometheus_tsdb_checkpoint_creations_failed_total 0
# HELP prometheus_tsdb_checkpoint_creations_total Total number of checkpoint creations attempted.
# TYPE prometheus_tsdb_checkpoint_creations_total counter
prometheus_tsdb_checkpoint_creations_total 5
# HELP prometheus_tsdb_checkpoint_deletions_failed_total Total number of checkpoint deletions that failed.
# TYPE prometheus_tsdb_checkpoint_deletions_failed_total counter
prometheus_tsdb_checkpoint_deletions_failed_total 0
# HELP prometheus_tsdb_checkpoint_deletions_total Total number of checkpoint deletions attempted.
# TYPE prometheus_tsdb_checkpoint_deletions_total counter
prometheus_tsdb_checkpoint_deletions_total 5
# HELP prometheus_tsdb_clean_start -1: lockfile is disabled. 0: a lockfile from a previous execution was replaced. 1: lockfile creation was clean
# TYPE prometheus_tsdb_clean_start gauge
prometheus_tsdb_clean_start 1
# HELP prometheus_tsdb_compaction_chunk_range_seconds Final time range of chunks on their first compaction
# TYPE prometheus_tsdb_compaction_chunk_range_seconds histogram
prometheus_tsdb_compaction_chunk_range_seconds_bucket{le="100"} 1302
prometheus_tsdb_compaction_chunk_range_seconds_bucket{le="400"} 1302
prometheus_tsdb_compaction_chunk_range_seconds_bucket{le="1600"} 1302
prometheus_tsdb_compaction_chunk_range_seconds_bucket{le="6400"} 1302
prometheus_tsdb_compaction_chunk_range_seconds_bucket{le="25600"} 1704
prometheus_tsdb_compaction_chunk_range_seconds_bucket{le="102400"} 6492
prometheus_tsdb_compaction_chunk_range_seconds_bucket{le="409600"} 7162
prometheus_tsdb_compaction_chunk_range_seconds_bucket{le="1.6384e+06"} 442409
prometheus_tsdb_compaction_chunk_range_seconds_bucket{le="6.5536e+06"} 3.313629e+06
prometheus_tsdb_compaction_chunk_range_seconds_bucket{le="2.62144e+07"} 3.508072e+06
prometheus_tsdb_compaction_chunk_range_seconds_bucket{le="+Inf"} 3.508072e+06
prometheus_tsdb_compaction_chunk_range_seconds_sum 1.2072884925365e+13
prometheus_tsdb_compaction_chunk_range_seconds_count 3.508072e+06
# HELP prometheus_tsdb_compaction_chunk_samples Final number of samples on their first compaction
# TYPE prometheus_tsdb_compaction_chunk_samples histogram
prometheus_tsdb_compaction_chunk_samples_bucket{le="4"} 6633
prometheus_tsdb_compaction_chunk_samples_bucket{le="6"} 6975
prometheus_tsdb_compaction_chunk_samples_bucket{le="9"} 7047
prometheus_tsdb_compaction_chunk_samples_bucket{le="13.5"} 7491
prometheus_tsdb_compaction_chunk_samples_bucket{le="20.25"} 8394
prometheus_tsdb_compaction_chunk_samples_bucket{le="30.375"} 47751
prometheus_tsdb_compaction_chunk_samples_bucket{le="45.5625"} 49518
prometheus_tsdb_compaction_chunk_samples_bucket{le="68.34375"} 323517
prometheus_tsdb_compaction_chunk_samples_bucket{le="102.515625"} 334056
prometheus_tsdb_compaction_chunk_samples_bucket{le="153.7734375"} 3.503897e+06
prometheus_tsdb_compaction_chunk_samples_bucket{le="230.66015625"} 3.508046e+06
prometheus_tsdb_compaction_chunk_samples_bucket{le="345.990234375"} 3.508072e+06
prometheus_tsdb_compaction_chunk_samples_bucket{le="+Inf"} 3.508072e+06
prometheus_tsdb_compaction_chunk_samples_sum 3.97404056e+08
prometheus_tsdb_compaction_chunk_samples_count 3.508072e+06
# HELP prometheus_tsdb_compaction_chunk_size_bytes Final size of chunks on their first compaction
# TYPE prometheus_tsdb_compaction_chunk_size_bytes histogram
prometheus_tsdb_compaction_chunk_size_bytes_bucket{le="32"} 16627
prometheus_tsdb_compaction_chunk_size_bytes_bucket{le="48"} 198084
prometheus_tsdb_compaction_chunk_size_bytes_bucket{le="72"} 1.463891e+06
prometheus_tsdb_compaction_chunk_size_bytes_bucket{le="108"} 2.028826e+06
prometheus_tsdb_compaction_chunk_size_bytes_bucket{le="162"} 2.308438e+06
prometheus_tsdb_compaction_chunk_size_bytes_bucket{le="243"} 2.63113e+06
prometheus_tsdb_compaction_chunk_size_bytes_bucket{le="364.5"} 3.280042e+06
prometheus_tsdb_compaction_chunk_size_bytes_bucket{le="546.75"} 3.404723e+06
prometheus_tsdb_compaction_chunk_size_bytes_bucket{le="820.125"} 3.471002e+06
prometheus_tsdb_compaction_chunk_size_bytes_bucket{le="1230.1875"} 3.507873e+06
prometheus_tsdb_compaction_chunk_size_bytes_bucket{le="1845.28125"} 3.508072e+06
prometheus_tsdb_compaction_chunk_size_bytes_bucket{le="2767.921875"} 3.508072e+06
prometheus_tsdb_compaction_chunk_size_bytes_bucket{le="+Inf"} 3.508072e+06
prometheus_tsdb_compaction_chunk_size_bytes_sum 5.58908691e+08
prometheus_tsdb_compaction_chunk_size_bytes_count 3.508072e+06
# HELP prometheus_tsdb_compaction_duration_seconds Duration of compaction runs
# TYPE prometheus_tsdb_compaction_duration_seconds histogram
prometheus_tsdb_compaction_duration_seconds_bucket{le="1"} 0
prometheus_tsdb_compaction_duration_seconds_bucket{le="2"} 0
prometheus_tsdb_compaction_duration_seconds_bucket{le="4"} 0
prometheus_tsdb_compaction_duration_seconds_bucket{le="8"} 1
prometheus_tsdb_compaction_duration_seconds_bucket{le="16"} 6
prometheus_tsdb_compaction_duration_seconds_bucket{le="32"} 6
prometheus_tsdb_compaction_duration_seconds_bucket{le="64"} 6
prometheus_tsdb_compaction_duration_seconds_bucket{le="128"} 6
prometheus_tsdb_compaction_duration_seconds_bucket{le="256"} 6
prometheus_tsdb_compaction_duration_seconds_bucket{le="512"} 6
prometheus_tsdb_compaction_duration_seconds_bucket{le="1024"} 6
prometheus_tsdb_compaction_duration_seconds_bucket{le="2048"} 6
prometheus_tsdb_compaction_duration_seconds_bucket{le="4096"} 6
prometheus_tsdb_compaction_duration_seconds_bucket{le="8192"} 6
prometheus_tsdb_compaction_duration_seconds_bucket{le="+Inf"} 6
prometheus_tsdb_compaction_duration_seconds_sum 54.148245326
prometheus_tsdb_compaction_duration_seconds_count 6
# HELP prometheus_tsdb_compaction_populating_block Set to 1 when a block is currently being written to the disk.
# TYPE prometheus_tsdb_compaction_populating_block gauge
prometheus_tsdb_compaction_populating_block 0
# HELP prometheus_tsdb_compactions_failed_total Total number of compactions that failed for the partition.
# TYPE prometheus_tsdb_compactions_failed_total counter
prometheus_tsdb_compactions_failed_total 0
# HELP prometheus_tsdb_compactions_skipped_total Total number of skipped compactions due to disabled auto compaction.
# TYPE prometheus_tsdb_compactions_skipped_total counter
prometheus_tsdb_compactions_skipped_total 0
# HELP prometheus_tsdb_compactions_total Total number of compactions that were executed for the partition.
# TYPE prometheus_tsdb_compactions_total counter
prometheus_tsdb_compactions_total 6
# HELP prometheus_tsdb_compactions_triggered_total Total number of triggered compactions for the partition.
# TYPE prometheus_tsdb_compactions_triggered_total counter
prometheus_tsdb_compactions_triggered_total 767
# HELP prometheus_tsdb_data_replay_duration_seconds Time taken to replay the data on disk.
# TYPE prometheus_tsdb_data_replay_duration_seconds gauge
prometheus_tsdb_data_replay_duration_seconds 0.002696783
# HELP prometheus_tsdb_exemplar_exemplars_appended_total Total number of appended exemplars.
# TYPE prometheus_tsdb_exemplar_exemplars_appended_total counter
prometheus_tsdb_exemplar_exemplars_appended_total 0
# HELP prometheus_tsdb_exemplar_exemplars_in_storage Number of exemplars currently in circular storage.
# TYPE prometheus_tsdb_exemplar_exemplars_in_storage gauge
prometheus_tsdb_exemplar_exemplars_in_storage 0
# HELP prometheus_tsdb_exemplar_last_exemplars_timestamp_seconds The timestamp of the oldest exemplar stored in circular storage. Useful to check for what timerange the current exemplar buffer limit allows. This usually means the last timestampfor all exemplars for a typical setup. This is not true though if one of the series timestamp is in future compared to rest series.
# TYPE prometheus_tsdb_exemplar_last_exemplars_timestamp_seconds gauge
prometheus_tsdb_exemplar_last_exemplars_timestamp_seconds 0
# HELP prometheus_tsdb_exemplar_max_exemplars Total number of exemplars the exemplar storage can store, resizeable.
# TYPE prometheus_tsdb_exemplar_max_exemplars gauge
prometheus_tsdb_exemplar_max_exemplars 0
# HELP prometheus_tsdb_exemplar_out_of_order_exemplars_total Total number of out of order exemplar ingestion failed attempts.
# TYPE prometheus_tsdb_exemplar_out_of_order_exemplars_total counter
prometheus_tsdb_exemplar_out_of_order_exemplars_total 0
# HELP prometheus_tsdb_exemplar_series_with_exemplars_in_storage Number of series with exemplars currently in circular storage.
# TYPE prometheus_tsdb_exemplar_series_with_exemplars_in_storage gauge
prometheus_tsdb_exemplar_series_with_exemplars_in_storage 0
# HELP prometheus_tsdb_head_active_appenders Number of currently active appender transactions
# TYPE prometheus_tsdb_head_active_appenders gauge
prometheus_tsdb_head_active_appenders 2
# HELP prometheus_tsdb_head_chunks Total number of chunks in the head block.
# TYPE prometheus_tsdb_head_chunks gauge
prometheus_tsdb_head_chunks 963506
# HELP prometheus_tsdb_head_chunks_created_total Total number of chunks created in the head
# TYPE prometheus_tsdb_head_chunks_created_total counter
prometheus_tsdb_head_chunks_created_total 4.471578e+06
# HELP prometheus_tsdb_head_chunks_removed_total Total number of chunks removed in the head
# TYPE prometheus_tsdb_head_chunks_removed_total counter
prometheus_tsdb_head_chunks_removed_total 3.508072e+06
# HELP prometheus_tsdb_head_gc_duration_seconds Runtime of garbage collection in the head block.
# TYPE prometheus_tsdb_head_gc_duration_seconds summary
prometheus_tsdb_head_gc_duration_seconds_sum 2.4343483850000003
prometheus_tsdb_head_gc_duration_seconds_count 6
# HELP prometheus_tsdb_head_max_time Maximum timestamp of the head block. The unit is decided by the library consumer.
# TYPE prometheus_tsdb_head_max_time gauge
prometheus_tsdb_head_max_time 1.672647240392e+12
# HELP prometheus_tsdb_head_max_time_seconds Maximum timestamp of the head block.
# TYPE prometheus_tsdb_head_max_time_seconds gauge
prometheus_tsdb_head_max_time_seconds 1.67264724e+09
# HELP prometheus_tsdb_head_min_time Minimum time bound of the head block. The unit is decided by the library consumer.
# TYPE prometheus_tsdb_head_min_time gauge
prometheus_tsdb_head_min_time 1.67263920006e+12
# HELP prometheus_tsdb_head_min_time_seconds Minimum time bound of the head block.
# TYPE prometheus_tsdb_head_min_time_seconds gauge
prometheus_tsdb_head_min_time_seconds 1.6726392e+09
# HELP prometheus_tsdb_head_out_of_order_samples_appended_total Total number of appended out of order samples.
# TYPE prometheus_tsdb_head_out_of_order_samples_appended_total counter
prometheus_tsdb_head_out_of_order_samples_appended_total 0
# HELP prometheus_tsdb_head_samples_appended_total Total number of appended samples.
# TYPE prometheus_tsdb_head_samples_appended_total counter
prometheus_tsdb_head_samples_appended_total{type="float"} 4.82529504e+08
prometheus_tsdb_head_samples_appended_total{type="histogram"} 0
# HELP prometheus_tsdb_head_series Total number of series in the head block.
# TYPE prometheus_tsdb_head_series gauge
prometheus_tsdb_head_series 326598
# HELP prometheus_tsdb_head_series_created_total Total number of series created in the head
# TYPE prometheus_tsdb_head_series_created_total counter
prometheus_tsdb_head_series_created_total 335399
# HELP prometheus_tsdb_head_series_not_found_total Total number of requests for series that were not found.
# TYPE prometheus_tsdb_head_series_not_found_total counter
prometheus_tsdb_head_series_not_found_total 0
# HELP prometheus_tsdb_head_series_removed_total Total number of series removed in the head
# TYPE prometheus_tsdb_head_series_removed_total counter
prometheus_tsdb_head_series_removed_total 8801
# HELP prometheus_tsdb_head_truncations_failed_total Total number of head truncations that failed.
# TYPE prometheus_tsdb_head_truncations_failed_total counter
prometheus_tsdb_head_truncations_failed_total 0
# HELP prometheus_tsdb_head_truncations_total Total number of head truncations attempted.
# TYPE prometheus_tsdb_head_truncations_total counter
prometheus_tsdb_head_truncations_total 6
# HELP prometheus_tsdb_isolation_high_watermark The highest TSDB append ID that has been given out.
# TYPE prometheus_tsdb_isolation_high_watermark gauge
prometheus_tsdb_isolation_high_watermark 975050
# HELP prometheus_tsdb_isolation_low_watermark The lowest TSDB append ID that is still referenced.
# TYPE prometheus_tsdb_isolation_low_watermark gauge
prometheus_tsdb_isolation_low_watermark 974477
# HELP prometheus_tsdb_lowest_timestamp Lowest timestamp value stored in the database. The unit is decided by the library consumer.
# TYPE prometheus_tsdb_lowest_timestamp gauge
prometheus_tsdb_lowest_timestamp 1.67260153102e+12
# HELP prometheus_tsdb_lowest_timestamp_seconds Lowest timestamp value stored in the database.
# TYPE prometheus_tsdb_lowest_timestamp_seconds gauge
prometheus_tsdb_lowest_timestamp_seconds 1.672601531e+09
# HELP prometheus_tsdb_mmap_chunk_corruptions_total Total number of memory-mapped chunk corruptions.
# TYPE prometheus_tsdb_mmap_chunk_corruptions_total counter
prometheus_tsdb_mmap_chunk_corruptions_total 0
# HELP prometheus_tsdb_out_of_bound_samples_total Total number of out of bound samples ingestion failed attempts with out of order support disabled.
# TYPE prometheus_tsdb_out_of_bound_samples_total counter
prometheus_tsdb_out_of_bound_samples_total{type="float"} 0
# HELP prometheus_tsdb_out_of_order_samples_total Total number of out of order samples ingestion failed attempts due to out of order being disabled.
# TYPE prometheus_tsdb_out_of_order_samples_total counter
prometheus_tsdb_out_of_order_samples_total{type="float"} 0
prometheus_tsdb_out_of_order_samples_total{type="histogram"} 0
# HELP prometheus_tsdb_reloads_failures_total Number of times the database failed to reloadBlocks block data from disk.
# TYPE prometheus_tsdb_reloads_failures_total counter
prometheus_tsdb_reloads_failures_total 0
# HELP prometheus_tsdb_reloads_total Number of times the database reloaded block data from disk.
# TYPE prometheus_tsdb_reloads_total counter
prometheus_tsdb_reloads_total 762
# HELP prometheus_tsdb_retention_limit_bytes Max number of bytes to be retained in the tsdb blocks, configured 0 means disabled
# TYPE prometheus_tsdb_retention_limit_bytes gauge
prometheus_tsdb_retention_limit_bytes 0
# HELP prometheus_tsdb_size_retentions_total The number of times that blocks were deleted because the maximum number of bytes was exceeded.
# TYPE prometheus_tsdb_size_retentions_total counter
prometheus_tsdb_size_retentions_total 0
# HELP prometheus_tsdb_snapshot_replay_error_total Total number snapshot replays that failed.
# TYPE prometheus_tsdb_snapshot_replay_error_total counter
prometheus_tsdb_snapshot_replay_error_total 0
# HELP prometheus_tsdb_storage_blocks_bytes The number of bytes that are currently used for local storage by all blocks.
# TYPE prometheus_tsdb_storage_blocks_bytes gauge
prometheus_tsdb_storage_blocks_bytes 8.67104887e+08
# HELP prometheus_tsdb_symbol_table_size_bytes Size of symbol table in memory for loaded blocks
# TYPE prometheus_tsdb_symbol_table_size_bytes gauge
prometheus_tsdb_symbol_table_size_bytes 26136
# HELP prometheus_tsdb_time_retentions_total The number of times that blocks were deleted because the maximum time limit was exceeded.
# TYPE prometheus_tsdb_time_retentions_total counter
prometheus_tsdb_time_retentions_total 0
# HELP prometheus_tsdb_tombstone_cleanup_seconds The time taken to recompact blocks to remove tombstones.
# TYPE prometheus_tsdb_tombstone_cleanup_seconds histogram
prometheus_tsdb_tombstone_cleanup_seconds_bucket{le="0.005"} 0
prometheus_tsdb_tombstone_cleanup_seconds_bucket{le="0.01"} 0
prometheus_tsdb_tombstone_cleanup_seconds_bucket{le="0.025"} 0
prometheus_tsdb_tombstone_cleanup_seconds_bucket{le="0.05"} 0
prometheus_tsdb_tombstone_cleanup_seconds_bucket{le="0.1"} 0
prometheus_tsdb_tombstone_cleanup_seconds_bucket{le="0.25"} 0
prometheus_tsdb_tombstone_cleanup_seconds_bucket{le="0.5"} 0
prometheus_tsdb_tombstone_cleanup_seconds_bucket{le="1"} 0
prometheus_tsdb_tombstone_cleanup_seconds_bucket{le="2.5"} 0
prometheus_tsdb_tombstone_cleanup_seconds_bucket{le="5"} 0
prometheus_tsdb_tombstone_cleanup_seconds_bucket{le="10"} 0
prometheus_tsdb_tombstone_cleanup_seconds_bucket{le="+Inf"} 0
prometheus_tsdb_tombstone_cleanup_seconds_sum 0
prometheus_tsdb_tombstone_cleanup_seconds_count 0
# HELP prometheus_tsdb_too_old_samples_total Total number of out of order samples ingestion failed attempts with out of support enabled, but sample outside of time window.
# TYPE prometheus_tsdb_too_old_samples_total counter
prometheus_tsdb_too_old_samples_total{type="float"} 0
# HELP prometheus_tsdb_vertical_compactions_total Total number of compactions done on overlapping blocks.
# TYPE prometheus_tsdb_vertical_compactions_total counter
prometheus_tsdb_vertical_compactions_total 0
# HELP prometheus_tsdb_wal_completed_pages_total Total number of completed pages.
# TYPE prometheus_tsdb_wal_completed_pages_total counter
prometheus_tsdb_wal_completed_pages_total 80531
# HELP prometheus_tsdb_wal_corruptions_total Total number of WAL corruptions.
# TYPE prometheus_tsdb_wal_corruptions_total counter
prometheus_tsdb_wal_corruptions_total 0
# HELP prometheus_tsdb_wal_fsync_duration_seconds Duration of write log fsync.
# TYPE prometheus_tsdb_wal_fsync_duration_seconds summary
prometheus_tsdb_wal_fsync_duration_seconds{quantile="0.5"} NaN
prometheus_tsdb_wal_fsync_duration_seconds{quantile="0.9"} NaN
prometheus_tsdb_wal_fsync_duration_seconds{quantile="0.99"} NaN
prometheus_tsdb_wal_fsync_duration_seconds_sum 2.5175242759999996
prometheus_tsdb_wal_fsync_duration_seconds_count 23
# HELP prometheus_tsdb_wal_page_flushes_total Total number of page flushes.
# TYPE prometheus_tsdb_wal_page_flushes_total counter
prometheus_tsdb_wal_page_flushes_total 307536
# HELP prometheus_tsdb_wal_segment_current Write log segment index that TSDB is currently writing to.
# TYPE prometheus_tsdb_wal_segment_current gauge
prometheus_tsdb_wal_segment_current 23
# HELP prometheus_tsdb_wal_truncate_duration_seconds Duration of WAL truncation.
# TYPE prometheus_tsdb_wal_truncate_duration_seconds summary
prometheus_tsdb_wal_truncate_duration_seconds_sum 32.33774713
prometheus_tsdb_wal_truncate_duration_seconds_count 5
# HELP prometheus_tsdb_wal_truncations_failed_total Total number of write log truncations that failed.
# TYPE prometheus_tsdb_wal_truncations_failed_total counter
prometheus_tsdb_wal_truncations_failed_total 0
# HELP prometheus_tsdb_wal_truncations_total Total number of write log truncations attempted.
# TYPE prometheus_tsdb_wal_truncations_total counter
prometheus_tsdb_wal_truncations_total 5
# HELP prometheus_tsdb_wal_writes_failed_total Total number of write log writes that failed.
# TYPE prometheus_tsdb_wal_writes_failed_total counter
prometheus_tsdb_wal_writes_failed_total 0
# HELP prometheus_wal_watcher_current_segment Current segment the WAL watcher is reading records from.
# TYPE prometheus_wal_watcher_current_segment gauge
prometheus_wal_watcher_current_segment{consumer="28d637"} 23
prometheus_wal_watcher_current_segment{consumer="4bfe0a"} 0
prometheus_wal_watcher_current_segment{consumer="c93768"} 0
# HELP prometheus_wal_watcher_record_decode_failures_total Number of records read by the WAL watcher that resulted in an error when decoding.
# TYPE prometheus_wal_watcher_record_decode_failures_total counter
prometheus_wal_watcher_record_decode_failures_total{consumer="28d637"} 0
prometheus_wal_watcher_record_decode_failures_total{consumer="4bfe0a"} 0
prometheus_wal_watcher_record_decode_failures_total{consumer="c93768"} 0
# HELP prometheus_wal_watcher_records_read_total Number of records read by the WAL watcher from the WAL.
# TYPE prometheus_wal_watcher_records_read_total counter
prometheus_wal_watcher_records_read_total{consumer="28d637",type="samples"} 251412
prometheus_wal_watcher_records_read_total{consumer="28d637",type="series"} 4200
prometheus_wal_watcher_records_read_total{consumer="4bfe0a",type="samples"} 4
prometheus_wal_watcher_records_read_total{consumer="4bfe0a",type="series"} 4
prometheus_wal_watcher_records_read_total{consumer="c93768",type="samples"} 4
prometheus_wal_watcher_records_read_total{consumer="c93768",type="series"} 4
# HELP prometheus_wal_watcher_samples_sent_pre_tailing_total Number of sample records read by the WAL watcher and sent to remote write during replay of existing WAL.
# TYPE prometheus_wal_watcher_samples_sent_pre_tailing_total counter
prometheus_wal_watcher_samples_sent_pre_tailing_total{consumer="28d637"} 0
prometheus_wal_watcher_samples_sent_pre_tailing_total{consumer="4bfe0a"} 0
prometheus_wal_watcher_samples_sent_pre_tailing_total{consumer="c93768"} 0
# HELP prometheus_web_federation_errors_total Total number of errors that occurred while sending federation responses.
# TYPE prometheus_web_federation_errors_total counter
prometheus_web_federation_errors_total 0
# HELP prometheus_web_federation_warnings_total Total number of warnings that occurred while sending federation responses.
# TYPE prometheus_web_federation_warnings_total counter
prometheus_web_federation_warnings_total 0
# HELP promhttp_metric_handler_requests_in_flight Current number of scrapes being served.
# TYPE promhttp_metric_handler_requests_in_flight gauge
promhttp_metric_handler_requests_in_flight 1
# HELP promhttp_metric_handler_requests_total Total number of scrapes by HTTP status code.
# TYPE promhttp_metric_handler_requests_total counter
promhttp_metric_handler_requests_total{code="200"} 1523
promhttp_metric_handler_requests_total{code="500"} 0
promhttp_metric_handler_requests_total{code="503"} 0
```
