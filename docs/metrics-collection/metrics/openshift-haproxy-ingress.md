This document describes the metrics list and how to collect them from Openshift HAProxy metrics.

# Metrics

Consul already can exposes its metrics in Prometheus format and doesn't require to use of specific exporters.

| Name       | Metrics Port | Metrics Endpoint    | Need Exporter? | Auth?          | Is Exporter Third Party? |
| ---------- | ------------ | ------------------- | -------------- | -------------- | ------------------------ |
| Prometheus | `80`         | `/metrics`          | No             | Require, token | N/A                      |

## How to Collect

Metrics expose on port `80` and endpoint `/metrics`. By default, Openshift HAProxy has authentication by token.

Config `ServiceMonitor` for `prometheus-operator` to collect Openshift HAProxy metrics:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: openshift-haproxy-router-service-monitor
  labels:
    k8s-app: openshift-haproxy-router
    app.kubernetes.io/name: openshift-haproxy-router
    app.kubernetes.io/component: monitoring
    app.kubernetes.io/managed-by: monitoring-operator
spec:
  endpoints:
    - bearerTokenFile: /var/run/secrets/kubernetes.io/serviceaccount/token
      interval: 30s
      path: /metrics
      port: metrics
      scheme: https
      tlsConfig:
        insecureSkipVerify: true
  namespaceSelector:
    matchNames:
      - openshift-ingress
  selector:
    matchLabels:
      ingresscontroller.operator.openshift.io/owning-ingresscontroller: default
```

To collect (or just to check) metrics manually you can use the following command:

```bash
curl -v -k -L -H "Authorization: Bearer <token>" "http://<api_server_ip_or_dns>/metrics"
```

Token usually you can find in the prometheus/vmagent pod by path `/var/run/secrets/kubernetes.io/serviceaccount/token`.

You can't use `wget` because it doesn't allow to add headers for authorization.

## Metrics List

### Openshift HAProxy metrics

```prometheus
# HELP go_gc_duration_seconds A summary of the pause duration of garbage collection cycles.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 5.8668e-05
go_gc_duration_seconds{quantile="0.25"} 0.000101402
go_gc_duration_seconds{quantile="0.5"} 0.000129657
go_gc_duration_seconds{quantile="0.75"} 0.000203654
go_gc_duration_seconds{quantile="1"} 0.007646407
go_gc_duration_seconds_sum 5.033045947
go_gc_duration_seconds_count 11799
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 48
# HELP go_info Information about the Go environment.
# TYPE go_info gauge
go_info{version="go1.20.10 X:strictfipsruntime"} 1
# HELP go_memstats_alloc_bytes Number of bytes allocated and still in use.
# TYPE go_memstats_alloc_bytes gauge
go_memstats_alloc_bytes 1.855736e+07
# HELP go_memstats_alloc_bytes_total Total number of bytes allocated, even if freed.
# TYPE go_memstats_alloc_bytes_total counter
go_memstats_alloc_bytes_total 1.39527723384e+11
# HELP go_memstats_buck_hash_sys_bytes Number of bytes used by the profiling bucket hash table.
# TYPE go_memstats_buck_hash_sys_bytes gauge
go_memstats_buck_hash_sys_bytes 2.020231e+06
# HELP go_memstats_frees_total Total number of frees.
# TYPE go_memstats_frees_total counter
go_memstats_frees_total 2.123773695e+09
# HELP go_memstats_gc_sys_bytes Number of bytes used for garbage collection system metadata.
# TYPE go_memstats_gc_sys_bytes gauge
go_memstats_gc_sys_bytes 9.361464e+06
# HELP go_memstats_heap_alloc_bytes Number of heap bytes allocated and still in use.
# TYPE go_memstats_heap_alloc_bytes gauge
go_memstats_heap_alloc_bytes 1.855736e+07
# HELP go_memstats_heap_idle_bytes Number of heap bytes waiting to be used.
# TYPE go_memstats_heap_idle_bytes gauge
go_memstats_heap_idle_bytes 1.5474688e+07
# HELP go_memstats_heap_inuse_bytes Number of heap bytes that are in use.
# TYPE go_memstats_heap_inuse_bytes gauge
go_memstats_heap_inuse_bytes 2.4600576e+07
# HELP go_memstats_heap_objects Number of allocated objects.
# TYPE go_memstats_heap_objects gauge
go_memstats_heap_objects 193449
# HELP go_memstats_heap_released_bytes Number of heap bytes released to OS.
# TYPE go_memstats_heap_released_bytes gauge
go_memstats_heap_released_bytes 8.8064e+06
# HELP go_memstats_heap_sys_bytes Number of heap bytes obtained from system.
# TYPE go_memstats_heap_sys_bytes gauge
go_memstats_heap_sys_bytes 4.0075264e+07
# HELP go_memstats_last_gc_time_seconds Number of seconds since 1970 of last garbage collection.
# TYPE go_memstats_last_gc_time_seconds gauge
go_memstats_last_gc_time_seconds 1.7035056587277122e+09
# HELP go_memstats_lookups_total Total number of pointer lookups.
# TYPE go_memstats_lookups_total counter
go_memstats_lookups_total 0
# HELP go_memstats_mallocs_total Total number of mallocs.
# TYPE go_memstats_mallocs_total counter
go_memstats_mallocs_total 2.123967144e+09
# HELP go_memstats_mcache_inuse_bytes Number of bytes in use by mcache structures.
# TYPE go_memstats_mcache_inuse_bytes gauge
go_memstats_mcache_inuse_bytes 9600
# HELP go_memstats_mcache_sys_bytes Number of bytes used for mcache structures obtained from system.
# TYPE go_memstats_mcache_sys_bytes gauge
go_memstats_mcache_sys_bytes 15600
# HELP go_memstats_mspan_inuse_bytes Number of bytes in use by mspan structures.
# TYPE go_memstats_mspan_inuse_bytes gauge
go_memstats_mspan_inuse_bytes 418080
# HELP go_memstats_mspan_sys_bytes Number of bytes used for mspan structures obtained from system.
# TYPE go_memstats_mspan_sys_bytes gauge
go_memstats_mspan_sys_bytes 538560
# HELP go_memstats_next_gc_bytes Number of heap bytes when next garbage collection will take place.
# TYPE go_memstats_next_gc_bytes gauge
go_memstats_next_gc_bytes 2.7018416e+07
# HELP go_memstats_other_sys_bytes Number of bytes used for other system allocations.
# TYPE go_memstats_other_sys_bytes gauge
go_memstats_other_sys_bytes 1.496217e+06
# HELP go_memstats_stack_inuse_bytes Number of bytes in use by the stack allocator.
# TYPE go_memstats_stack_inuse_bytes gauge
go_memstats_stack_inuse_bytes 1.867776e+06
# HELP go_memstats_stack_sys_bytes Number of bytes obtained from system for stack allocator.
# TYPE go_memstats_stack_sys_bytes gauge
go_memstats_stack_sys_bytes 1.867776e+06
# HELP go_memstats_sys_bytes Number of bytes obtained from system.
# TYPE go_memstats_sys_bytes gauge
go_memstats_sys_bytes 5.5375112e+07
# HELP go_threads Number of OS threads created.
# TYPE go_threads gauge
go_threads 16
# HELP haproxy_backend_bytes_in_total Current total of incoming bytes.
# TYPE haproxy_backend_bytes_in_total gauge
haproxy_backend_bytes_in_total{backend="http",namespace="airflow",route="airflow-ingress-nzbtm"} 0
haproxy_backend_bytes_in_total{backend="http",namespace="airflow",route="airflow-site-manager-hpl66"} 0
haproxy_backend_bytes_in_total{backend="http",namespace="arangodb",route="arangodb-leader-ingress-zw82p"} 0
haproxy_backend_bytes_in_total{backend="http",namespace="consul-service",route="consul-ingress-hmmvf"} 0
haproxy_backend_bytes_in_total{backend="http",namespace="dbaas",route="aggregator-fw5ft"} 0
haproxy_backend_bytes_in_total{backend="http",namespace="infra-keycloak",route="infra-keycloak-5q4rj"} 0
haproxy_backend_bytes_in_total{backend="http",namespace="jaeger",route="jaeger-hotrod-qq5cq"} 0
haproxy_backend_bytes_in_total{backend="http",namespace="jaeger",route="jaeger-query-pc6w5"} 0
haproxy_backend_bytes_in_total{backend="http",namespace="kafka-service",route="akhq-ingress-lhz5z"} 26307
haproxy_backend_bytes_in_total{backend="http",namespace="license-server",route="license-server-9txj9"} 0
haproxy_backend_bytes_in_total{backend="http",namespace="mistral",route="mistral-vslnb"} 0
haproxy_backend_bytes_in_total{backend="http",namespace="opensearch",route="opensearch-dashboards-l6vdv"} 0
haproxy_backend_bytes_in_total{backend="http",namespace="profiler",route="esc-collector-service-kdvdw"} 0
haproxy_backend_bytes_in_total{backend="http",namespace="profiler",route="esc-static-service-j9sgv"} 0
haproxy_backend_bytes_in_total{backend="http",namespace="profiler",route="esc-test-service-lbp8r"} 0
haproxy_backend_bytes_in_total{backend="http",namespace="profiler",route="esc-ui-service-r5tfm"} 0
haproxy_backend_bytes_in_total{backend="http",namespace="prometheus-operator",route="prometheus-operator-grafana"} 0
haproxy_backend_bytes_in_total{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmagent-k8s"} 0
haproxy_backend_bytes_in_total{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmalert-k8s"} 0
haproxy_backend_bytes_in_total{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmalertmanager-k8s"} 0
haproxy_backend_bytes_in_total{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmsingle-k8s"} 9319
haproxy_backend_bytes_in_total{backend="http",namespace="rabbitmq",route="rabbitmq-fnzhp"} 0
haproxy_backend_bytes_in_total{backend="http",namespace="shared-license-distributor",route="shared-license-distributor-wnj68"} 0
haproxy_backend_bytes_in_total{backend="http",namespace="spark",route="spark-history-server-hqshb"} 0
haproxy_backend_bytes_in_total{backend="http",namespace="spark",route="spark-site-manager-p4nqh"} 0
haproxy_backend_bytes_in_total{backend="https",namespace="openshift-console",route="console"} 8.5294893e+07
haproxy_backend_bytes_in_total{backend="https",namespace="openshift-monitoring",route="alertmanager-main"} 0
haproxy_backend_bytes_in_total{backend="https",namespace="openshift-monitoring",route="prometheus-k8s"} 0
haproxy_backend_bytes_in_total{backend="https",namespace="openshift-monitoring",route="prometheus-k8s-federate"} 0
haproxy_backend_bytes_in_total{backend="https",namespace="openshift-monitoring",route="thanos-querier"} 0
haproxy_backend_bytes_in_total{backend="https-edge",namespace="openshift-console",route="downloads"} 0
haproxy_backend_bytes_in_total{backend="https-edge",namespace="openshift-ingress-canary",route="canary"} 0
haproxy_backend_bytes_in_total{backend="other/be_no_sni",namespace="",route=""} 2610
haproxy_backend_bytes_in_total{backend="other/be_sni",namespace="",route=""} 8.846819e+07
haproxy_backend_bytes_in_total{backend="other/openshift_default",namespace="",route=""} 10608
haproxy_backend_bytes_in_total{backend="tcp",namespace="openshift-authentication",route="oauth-openshift"} 214811
haproxy_backend_bytes_in_total{backend="tcp",namespace="site-manager",route="site-manager"} 0
# HELP haproxy_backend_bytes_out_total Current total of outgoing bytes.
# TYPE haproxy_backend_bytes_out_total gauge
haproxy_backend_bytes_out_total{backend="http",namespace="airflow",route="airflow-ingress-nzbtm"} 0
haproxy_backend_bytes_out_total{backend="http",namespace="airflow",route="airflow-site-manager-hpl66"} 0
haproxy_backend_bytes_out_total{backend="http",namespace="arangodb",route="arangodb-leader-ingress-zw82p"} 0
haproxy_backend_bytes_out_total{backend="http",namespace="consul-service",route="consul-ingress-hmmvf"} 0
haproxy_backend_bytes_out_total{backend="http",namespace="dbaas",route="aggregator-fw5ft"} 0
haproxy_backend_bytes_out_total{backend="http",namespace="infra-keycloak",route="infra-keycloak-5q4rj"} 0
haproxy_backend_bytes_out_total{backend="http",namespace="jaeger",route="jaeger-hotrod-qq5cq"} 0
haproxy_backend_bytes_out_total{backend="http",namespace="jaeger",route="jaeger-query-pc6w5"} 0
haproxy_backend_bytes_out_total{backend="http",namespace="kafka-service",route="akhq-ingress-lhz5z"} 581056
haproxy_backend_bytes_out_total{backend="http",namespace="license-server",route="license-server-9txj9"} 0
haproxy_backend_bytes_out_total{backend="http",namespace="mistral",route="mistral-vslnb"} 0
haproxy_backend_bytes_out_total{backend="http",namespace="opensearch",route="opensearch-dashboards-l6vdv"} 0
haproxy_backend_bytes_out_total{backend="http",namespace="profiler",route="esc-collector-service-kdvdw"} 0
haproxy_backend_bytes_out_total{backend="http",namespace="profiler",route="esc-static-service-j9sgv"} 0
haproxy_backend_bytes_out_total{backend="http",namespace="profiler",route="esc-test-service-lbp8r"} 0
haproxy_backend_bytes_out_total{backend="http",namespace="profiler",route="esc-ui-service-r5tfm"} 0
haproxy_backend_bytes_out_total{backend="http",namespace="prometheus-operator",route="prometheus-operator-grafana"} 0
haproxy_backend_bytes_out_total{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmagent-k8s"} 0
haproxy_backend_bytes_out_total{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmalert-k8s"} 0
haproxy_backend_bytes_out_total{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmalertmanager-k8s"} 0
haproxy_backend_bytes_out_total{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmsingle-k8s"} 851946
haproxy_backend_bytes_out_total{backend="http",namespace="rabbitmq",route="rabbitmq-fnzhp"} 0
haproxy_backend_bytes_out_total{backend="http",namespace="shared-license-distributor",route="shared-license-distributor-wnj68"} 0
haproxy_backend_bytes_out_total{backend="http",namespace="spark",route="spark-history-server-hqshb"} 0
haproxy_backend_bytes_out_total{backend="http",namespace="spark",route="spark-site-manager-p4nqh"} 0
haproxy_backend_bytes_out_total{backend="https",namespace="openshift-console",route="console"} 6.373787468e+09
haproxy_backend_bytes_out_total{backend="https",namespace="openshift-monitoring",route="alertmanager-main"} 0
haproxy_backend_bytes_out_total{backend="https",namespace="openshift-monitoring",route="prometheus-k8s"} 0
haproxy_backend_bytes_out_total{backend="https",namespace="openshift-monitoring",route="prometheus-k8s-federate"} 0
haproxy_backend_bytes_out_total{backend="https",namespace="openshift-monitoring",route="thanos-querier"} 0
haproxy_backend_bytes_out_total{backend="https-edge",namespace="openshift-console",route="downloads"} 0
haproxy_backend_bytes_out_total{backend="https-edge",namespace="openshift-ingress-canary",route="canary"} 0
haproxy_backend_bytes_out_total{backend="other/be_no_sni",namespace="",route=""} 0
haproxy_backend_bytes_out_total{backend="other/be_sni",namespace="",route=""} 4.734176415e+09
haproxy_backend_bytes_out_total{backend="other/openshift_default",namespace="",route=""} 23712
haproxy_backend_bytes_out_total{backend="tcp",namespace="openshift-authentication",route="oauth-openshift"} 2.921621e+06
haproxy_backend_bytes_out_total{backend="tcp",namespace="site-manager",route="site-manager"} 0
# HELP haproxy_backend_connection_errors_total Total of connection errors.
# TYPE haproxy_backend_connection_errors_total gauge
haproxy_backend_connection_errors_total{backend="http",namespace="airflow",route="airflow-ingress-nzbtm"} 0
haproxy_backend_connection_errors_total{backend="http",namespace="airflow",route="airflow-site-manager-hpl66"} 0
haproxy_backend_connection_errors_total{backend="http",namespace="arangodb",route="arangodb-leader-ingress-zw82p"} 0
haproxy_backend_connection_errors_total{backend="http",namespace="consul-service",route="consul-ingress-hmmvf"} 0
haproxy_backend_connection_errors_total{backend="http",namespace="dbaas",route="aggregator-fw5ft"} 0
haproxy_backend_connection_errors_total{backend="http",namespace="infra-keycloak",route="infra-keycloak-5q4rj"} 0
haproxy_backend_connection_errors_total{backend="http",namespace="jaeger",route="jaeger-hotrod-qq5cq"} 0
haproxy_backend_connection_errors_total{backend="http",namespace="jaeger",route="jaeger-query-pc6w5"} 0
haproxy_backend_connection_errors_total{backend="http",namespace="kafka-service",route="akhq-ingress-lhz5z"} 0
haproxy_backend_connection_errors_total{backend="http",namespace="license-server",route="license-server-9txj9"} 0
haproxy_backend_connection_errors_total{backend="http",namespace="mistral",route="mistral-vslnb"} 0
haproxy_backend_connection_errors_total{backend="http",namespace="opensearch",route="opensearch-dashboards-l6vdv"} 0
haproxy_backend_connection_errors_total{backend="http",namespace="profiler",route="esc-collector-service-kdvdw"} 0
haproxy_backend_connection_errors_total{backend="http",namespace="profiler",route="esc-static-service-j9sgv"} 0
haproxy_backend_connection_errors_total{backend="http",namespace="profiler",route="esc-test-service-lbp8r"} 0
haproxy_backend_connection_errors_total{backend="http",namespace="profiler",route="esc-ui-service-r5tfm"} 0
haproxy_backend_connection_errors_total{backend="http",namespace="prometheus-operator",route="prometheus-operator-grafana"} 0
haproxy_backend_connection_errors_total{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmagent-k8s"} 0
haproxy_backend_connection_errors_total{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmalert-k8s"} 0
haproxy_backend_connection_errors_total{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmalertmanager-k8s"} 0
haproxy_backend_connection_errors_total{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmsingle-k8s"} 0
haproxy_backend_connection_errors_total{backend="http",namespace="rabbitmq",route="rabbitmq-fnzhp"} 0
haproxy_backend_connection_errors_total{backend="http",namespace="shared-license-distributor",route="shared-license-distributor-wnj68"} 0
haproxy_backend_connection_errors_total{backend="http",namespace="spark",route="spark-history-server-hqshb"} 0
haproxy_backend_connection_errors_total{backend="http",namespace="spark",route="spark-site-manager-p4nqh"} 0
haproxy_backend_connection_errors_total{backend="https",namespace="openshift-console",route="console"} 0
haproxy_backend_connection_errors_total{backend="https",namespace="openshift-monitoring",route="alertmanager-main"} 0
haproxy_backend_connection_errors_total{backend="https",namespace="openshift-monitoring",route="prometheus-k8s"} 0
haproxy_backend_connection_errors_total{backend="https",namespace="openshift-monitoring",route="prometheus-k8s-federate"} 0
haproxy_backend_connection_errors_total{backend="https",namespace="openshift-monitoring",route="thanos-querier"} 0
haproxy_backend_connection_errors_total{backend="https-edge",namespace="openshift-console",route="downloads"} 0
haproxy_backend_connection_errors_total{backend="https-edge",namespace="openshift-ingress-canary",route="canary"} 0
haproxy_backend_connection_errors_total{backend="other/be_no_sni",namespace="",route=""} 0
haproxy_backend_connection_errors_total{backend="other/be_sni",namespace="",route=""} 0
haproxy_backend_connection_errors_total{backend="other/openshift_default",namespace="",route=""} 156
haproxy_backend_connection_errors_total{backend="tcp",namespace="openshift-authentication",route="oauth-openshift"} 0
haproxy_backend_connection_errors_total{backend="tcp",namespace="site-manager",route="site-manager"} 0
# HELP haproxy_backend_connections_reused_total Total number of connections reused.
# TYPE haproxy_backend_connections_reused_total gauge
haproxy_backend_connections_reused_total{backend="http",namespace="airflow",route="airflow-ingress-nzbtm"} 0
haproxy_backend_connections_reused_total{backend="http",namespace="airflow",route="airflow-site-manager-hpl66"} 0
haproxy_backend_connections_reused_total{backend="http",namespace="arangodb",route="arangodb-leader-ingress-zw82p"} 0
haproxy_backend_connections_reused_total{backend="http",namespace="consul-service",route="consul-ingress-hmmvf"} 0
haproxy_backend_connections_reused_total{backend="http",namespace="dbaas",route="aggregator-fw5ft"} 0
haproxy_backend_connections_reused_total{backend="http",namespace="infra-keycloak",route="infra-keycloak-5q4rj"} 0
haproxy_backend_connections_reused_total{backend="http",namespace="jaeger",route="jaeger-hotrod-qq5cq"} 0
haproxy_backend_connections_reused_total{backend="http",namespace="jaeger",route="jaeger-query-pc6w5"} 0
haproxy_backend_connections_reused_total{backend="http",namespace="kafka-service",route="akhq-ingress-lhz5z"} 19
haproxy_backend_connections_reused_total{backend="http",namespace="license-server",route="license-server-9txj9"} 0
haproxy_backend_connections_reused_total{backend="http",namespace="mistral",route="mistral-vslnb"} 0
haproxy_backend_connections_reused_total{backend="http",namespace="opensearch",route="opensearch-dashboards-l6vdv"} 0
haproxy_backend_connections_reused_total{backend="http",namespace="profiler",route="esc-collector-service-kdvdw"} 0
haproxy_backend_connections_reused_total{backend="http",namespace="profiler",route="esc-static-service-j9sgv"} 0
haproxy_backend_connections_reused_total{backend="http",namespace="profiler",route="esc-test-service-lbp8r"} 0
haproxy_backend_connections_reused_total{backend="http",namespace="profiler",route="esc-ui-service-r5tfm"} 0
haproxy_backend_connections_reused_total{backend="http",namespace="prometheus-operator",route="prometheus-operator-grafana"} 0
haproxy_backend_connections_reused_total{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmagent-k8s"} 0
haproxy_backend_connections_reused_total{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmalert-k8s"} 0
haproxy_backend_connections_reused_total{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmalertmanager-k8s"} 0
haproxy_backend_connections_reused_total{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmsingle-k8s"} 27
haproxy_backend_connections_reused_total{backend="http",namespace="rabbitmq",route="rabbitmq-fnzhp"} 0
haproxy_backend_connections_reused_total{backend="http",namespace="shared-license-distributor",route="shared-license-distributor-wnj68"} 0
haproxy_backend_connections_reused_total{backend="http",namespace="spark",route="spark-history-server-hqshb"} 0
haproxy_backend_connections_reused_total{backend="http",namespace="spark",route="spark-site-manager-p4nqh"} 0
haproxy_backend_connections_reused_total{backend="https",namespace="openshift-console",route="console"} 46902
haproxy_backend_connections_reused_total{backend="https",namespace="openshift-monitoring",route="alertmanager-main"} 0
haproxy_backend_connections_reused_total{backend="https",namespace="openshift-monitoring",route="prometheus-k8s"} 0
haproxy_backend_connections_reused_total{backend="https",namespace="openshift-monitoring",route="prometheus-k8s-federate"} 0
haproxy_backend_connections_reused_total{backend="https",namespace="openshift-monitoring",route="thanos-querier"} 0
haproxy_backend_connections_reused_total{backend="https-edge",namespace="openshift-console",route="downloads"} 0
haproxy_backend_connections_reused_total{backend="https-edge",namespace="openshift-ingress-canary",route="canary"} 0
haproxy_backend_connections_reused_total{backend="other/be_no_sni",namespace="",route=""} 0
haproxy_backend_connections_reused_total{backend="other/be_sni",namespace="",route=""} 0
haproxy_backend_connections_reused_total{backend="other/openshift_default",namespace="",route=""} 0
haproxy_backend_connections_reused_total{backend="tcp",namespace="openshift-authentication",route="oauth-openshift"} 0
haproxy_backend_connections_reused_total{backend="tcp",namespace="site-manager",route="site-manager"} 0
# HELP haproxy_backend_connections_total Total number of connections.
# TYPE haproxy_backend_connections_total gauge
haproxy_backend_connections_total{backend="http",namespace="airflow",route="airflow-ingress-nzbtm"} 0
haproxy_backend_connections_total{backend="http",namespace="airflow",route="airflow-site-manager-hpl66"} 0
haproxy_backend_connections_total{backend="http",namespace="arangodb",route="arangodb-leader-ingress-zw82p"} 0
haproxy_backend_connections_total{backend="http",namespace="consul-service",route="consul-ingress-hmmvf"} 0
haproxy_backend_connections_total{backend="http",namespace="dbaas",route="aggregator-fw5ft"} 0
haproxy_backend_connections_total{backend="http",namespace="infra-keycloak",route="infra-keycloak-5q4rj"} 0
haproxy_backend_connections_total{backend="http",namespace="jaeger",route="jaeger-hotrod-qq5cq"} 0
haproxy_backend_connections_total{backend="http",namespace="jaeger",route="jaeger-query-pc6w5"} 0
haproxy_backend_connections_total{backend="http",namespace="kafka-service",route="akhq-ingress-lhz5z"} 22
haproxy_backend_connections_total{backend="http",namespace="license-server",route="license-server-9txj9"} 0
haproxy_backend_connections_total{backend="http",namespace="mistral",route="mistral-vslnb"} 0
haproxy_backend_connections_total{backend="http",namespace="opensearch",route="opensearch-dashboards-l6vdv"} 0
haproxy_backend_connections_total{backend="http",namespace="profiler",route="esc-collector-service-kdvdw"} 0
haproxy_backend_connections_total{backend="http",namespace="profiler",route="esc-static-service-j9sgv"} 0
haproxy_backend_connections_total{backend="http",namespace="profiler",route="esc-test-service-lbp8r"} 0
haproxy_backend_connections_total{backend="http",namespace="profiler",route="esc-ui-service-r5tfm"} 0
haproxy_backend_connections_total{backend="http",namespace="prometheus-operator",route="prometheus-operator-grafana"} 0
haproxy_backend_connections_total{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmagent-k8s"} 0
haproxy_backend_connections_total{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmalert-k8s"} 0
haproxy_backend_connections_total{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmalertmanager-k8s"} 0
haproxy_backend_connections_total{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmsingle-k8s"} 36
haproxy_backend_connections_total{backend="http",namespace="rabbitmq",route="rabbitmq-fnzhp"} 0
haproxy_backend_connections_total{backend="http",namespace="shared-license-distributor",route="shared-license-distributor-wnj68"} 0
haproxy_backend_connections_total{backend="http",namespace="spark",route="spark-history-server-hqshb"} 0
haproxy_backend_connections_total{backend="http",namespace="spark",route="spark-site-manager-p4nqh"} 0
haproxy_backend_connections_total{backend="https",namespace="openshift-console",route="console"} 60949
haproxy_backend_connections_total{backend="https",namespace="openshift-monitoring",route="alertmanager-main"} 0
haproxy_backend_connections_total{backend="https",namespace="openshift-monitoring",route="prometheus-k8s"} 0
haproxy_backend_connections_total{backend="https",namespace="openshift-monitoring",route="prometheus-k8s-federate"} 0
haproxy_backend_connections_total{backend="https",namespace="openshift-monitoring",route="thanos-querier"} 0
haproxy_backend_connections_total{backend="https-edge",namespace="openshift-console",route="downloads"} 0
haproxy_backend_connections_total{backend="https-edge",namespace="openshift-ingress-canary",route="canary"} 0
haproxy_backend_connections_total{backend="other/be_no_sni",namespace="",route=""} 6
haproxy_backend_connections_total{backend="other/be_sni",namespace="",route=""} 24941
haproxy_backend_connections_total{backend="other/openshift_default",namespace="",route=""} 156
haproxy_backend_connections_total{backend="tcp",namespace="openshift-authentication",route="oauth-openshift"} 74
haproxy_backend_connections_total{backend="tcp",namespace="site-manager",route="site-manager"} 0
# HELP haproxy_backend_current_queue Current number of queued requests not assigned to any server.
# TYPE haproxy_backend_current_queue gauge
haproxy_backend_current_queue{backend="http",namespace="airflow",route="airflow-ingress-nzbtm"} 0
haproxy_backend_current_queue{backend="http",namespace="airflow",route="airflow-site-manager-hpl66"} 0
haproxy_backend_current_queue{backend="http",namespace="arangodb",route="arangodb-leader-ingress-zw82p"} 0
haproxy_backend_current_queue{backend="http",namespace="consul-service",route="consul-ingress-hmmvf"} 0
haproxy_backend_current_queue{backend="http",namespace="dbaas",route="aggregator-fw5ft"} 0
haproxy_backend_current_queue{backend="http",namespace="infra-keycloak",route="infra-keycloak-5q4rj"} 0
haproxy_backend_current_queue{backend="http",namespace="jaeger",route="jaeger-hotrod-qq5cq"} 0
haproxy_backend_current_queue{backend="http",namespace="jaeger",route="jaeger-query-pc6w5"} 0
haproxy_backend_current_queue{backend="http",namespace="kafka-service",route="akhq-ingress-lhz5z"} 0
haproxy_backend_current_queue{backend="http",namespace="license-server",route="license-server-9txj9"} 0
haproxy_backend_current_queue{backend="http",namespace="mistral",route="mistral-vslnb"} 0
haproxy_backend_current_queue{backend="http",namespace="opensearch",route="opensearch-dashboards-l6vdv"} 0
haproxy_backend_current_queue{backend="http",namespace="profiler",route="esc-collector-service-kdvdw"} 0
haproxy_backend_current_queue{backend="http",namespace="profiler",route="esc-static-service-j9sgv"} 0
haproxy_backend_current_queue{backend="http",namespace="profiler",route="esc-test-service-lbp8r"} 0
haproxy_backend_current_queue{backend="http",namespace="profiler",route="esc-ui-service-r5tfm"} 0
haproxy_backend_current_queue{backend="http",namespace="prometheus-operator",route="prometheus-operator-grafana"} 0
haproxy_backend_current_queue{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmagent-k8s"} 0
haproxy_backend_current_queue{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmalert-k8s"} 0
haproxy_backend_current_queue{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmalertmanager-k8s"} 0
haproxy_backend_current_queue{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmsingle-k8s"} 0
haproxy_backend_current_queue{backend="http",namespace="rabbitmq",route="rabbitmq-fnzhp"} 0
haproxy_backend_current_queue{backend="http",namespace="shared-license-distributor",route="shared-license-distributor-wnj68"} 0
haproxy_backend_current_queue{backend="http",namespace="spark",route="spark-history-server-hqshb"} 0
haproxy_backend_current_queue{backend="http",namespace="spark",route="spark-site-manager-p4nqh"} 0
haproxy_backend_current_queue{backend="https",namespace="openshift-console",route="console"} 0
haproxy_backend_current_queue{backend="https",namespace="openshift-monitoring",route="alertmanager-main"} 0
haproxy_backend_current_queue{backend="https",namespace="openshift-monitoring",route="prometheus-k8s"} 0
haproxy_backend_current_queue{backend="https",namespace="openshift-monitoring",route="prometheus-k8s-federate"} 0
haproxy_backend_current_queue{backend="https",namespace="openshift-monitoring",route="thanos-querier"} 0
haproxy_backend_current_queue{backend="https-edge",namespace="openshift-console",route="downloads"} 0
haproxy_backend_current_queue{backend="https-edge",namespace="openshift-ingress-canary",route="canary"} 0
haproxy_backend_current_queue{backend="other/be_no_sni",namespace="",route=""} 0
haproxy_backend_current_queue{backend="other/be_sni",namespace="",route=""} 0
haproxy_backend_current_queue{backend="other/openshift_default",namespace="",route=""} 0
haproxy_backend_current_queue{backend="tcp",namespace="openshift-authentication",route="oauth-openshift"} 0
haproxy_backend_current_queue{backend="tcp",namespace="site-manager",route="site-manager"} 0
# HELP haproxy_backend_current_session_rate Current number of sessions per second over last elapsed second.
# TYPE haproxy_backend_current_session_rate gauge
haproxy_backend_current_session_rate{backend="http",namespace="airflow",route="airflow-ingress-nzbtm"} 0
haproxy_backend_current_session_rate{backend="http",namespace="airflow",route="airflow-site-manager-hpl66"} 0
haproxy_backend_current_session_rate{backend="http",namespace="arangodb",route="arangodb-leader-ingress-zw82p"} 0
haproxy_backend_current_session_rate{backend="http",namespace="consul-service",route="consul-ingress-hmmvf"} 0
haproxy_backend_current_session_rate{backend="http",namespace="dbaas",route="aggregator-fw5ft"} 0
haproxy_backend_current_session_rate{backend="http",namespace="infra-keycloak",route="infra-keycloak-5q4rj"} 0
haproxy_backend_current_session_rate{backend="http",namespace="jaeger",route="jaeger-hotrod-qq5cq"} 0
haproxy_backend_current_session_rate{backend="http",namespace="jaeger",route="jaeger-query-pc6w5"} 0
haproxy_backend_current_session_rate{backend="http",namespace="kafka-service",route="akhq-ingress-lhz5z"} 0
haproxy_backend_current_session_rate{backend="http",namespace="license-server",route="license-server-9txj9"} 0
haproxy_backend_current_session_rate{backend="http",namespace="mistral",route="mistral-vslnb"} 0
haproxy_backend_current_session_rate{backend="http",namespace="opensearch",route="opensearch-dashboards-l6vdv"} 0
haproxy_backend_current_session_rate{backend="http",namespace="profiler",route="esc-collector-service-kdvdw"} 0
haproxy_backend_current_session_rate{backend="http",namespace="profiler",route="esc-static-service-j9sgv"} 0
haproxy_backend_current_session_rate{backend="http",namespace="profiler",route="esc-test-service-lbp8r"} 0
haproxy_backend_current_session_rate{backend="http",namespace="profiler",route="esc-ui-service-r5tfm"} 0
haproxy_backend_current_session_rate{backend="http",namespace="prometheus-operator",route="prometheus-operator-grafana"} 0
haproxy_backend_current_session_rate{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmagent-k8s"} 0
haproxy_backend_current_session_rate{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmalert-k8s"} 0
haproxy_backend_current_session_rate{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmalertmanager-k8s"} 0
haproxy_backend_current_session_rate{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmsingle-k8s"} 0
haproxy_backend_current_session_rate{backend="http",namespace="rabbitmq",route="rabbitmq-fnzhp"} 0
haproxy_backend_current_session_rate{backend="http",namespace="shared-license-distributor",route="shared-license-distributor-wnj68"} 0
haproxy_backend_current_session_rate{backend="http",namespace="spark",route="spark-history-server-hqshb"} 0
haproxy_backend_current_session_rate{backend="http",namespace="spark",route="spark-site-manager-p4nqh"} 0
haproxy_backend_current_session_rate{backend="https",namespace="openshift-console",route="console"} 0
haproxy_backend_current_session_rate{backend="https",namespace="openshift-monitoring",route="alertmanager-main"} 0
haproxy_backend_current_session_rate{backend="https",namespace="openshift-monitoring",route="prometheus-k8s"} 0
haproxy_backend_current_session_rate{backend="https",namespace="openshift-monitoring",route="prometheus-k8s-federate"} 0
haproxy_backend_current_session_rate{backend="https",namespace="openshift-monitoring",route="thanos-querier"} 0
haproxy_backend_current_session_rate{backend="https-edge",namespace="openshift-console",route="downloads"} 0
haproxy_backend_current_session_rate{backend="https-edge",namespace="openshift-ingress-canary",route="canary"} 0
haproxy_backend_current_session_rate{backend="other/be_no_sni",namespace="",route=""} 0
haproxy_backend_current_session_rate{backend="other/be_sni",namespace="",route=""} 0
haproxy_backend_current_session_rate{backend="other/openshift_default",namespace="",route=""} 0
haproxy_backend_current_session_rate{backend="tcp",namespace="openshift-authentication",route="oauth-openshift"} 0
haproxy_backend_current_session_rate{backend="tcp",namespace="site-manager",route="site-manager"} 0
# HELP haproxy_backend_current_sessions Current number of active sessions.
# TYPE haproxy_backend_current_sessions gauge
haproxy_backend_current_sessions{backend="http",namespace="airflow",route="airflow-ingress-nzbtm"} 0
haproxy_backend_current_sessions{backend="http",namespace="airflow",route="airflow-site-manager-hpl66"} 0
haproxy_backend_current_sessions{backend="http",namespace="arangodb",route="arangodb-leader-ingress-zw82p"} 0
haproxy_backend_current_sessions{backend="http",namespace="consul-service",route="consul-ingress-hmmvf"} 0
haproxy_backend_current_sessions{backend="http",namespace="dbaas",route="aggregator-fw5ft"} 0
haproxy_backend_current_sessions{backend="http",namespace="infra-keycloak",route="infra-keycloak-5q4rj"} 0
haproxy_backend_current_sessions{backend="http",namespace="jaeger",route="jaeger-hotrod-qq5cq"} 0
haproxy_backend_current_sessions{backend="http",namespace="jaeger",route="jaeger-query-pc6w5"} 0
haproxy_backend_current_sessions{backend="http",namespace="kafka-service",route="akhq-ingress-lhz5z"} 0
haproxy_backend_current_sessions{backend="http",namespace="license-server",route="license-server-9txj9"} 0
haproxy_backend_current_sessions{backend="http",namespace="mistral",route="mistral-vslnb"} 0
haproxy_backend_current_sessions{backend="http",namespace="opensearch",route="opensearch-dashboards-l6vdv"} 0
haproxy_backend_current_sessions{backend="http",namespace="profiler",route="esc-collector-service-kdvdw"} 0
haproxy_backend_current_sessions{backend="http",namespace="profiler",route="esc-static-service-j9sgv"} 0
haproxy_backend_current_sessions{backend="http",namespace="profiler",route="esc-test-service-lbp8r"} 0
haproxy_backend_current_sessions{backend="http",namespace="profiler",route="esc-ui-service-r5tfm"} 0
haproxy_backend_current_sessions{backend="http",namespace="prometheus-operator",route="prometheus-operator-grafana"} 0
haproxy_backend_current_sessions{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmagent-k8s"} 0
haproxy_backend_current_sessions{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmalert-k8s"} 0
haproxy_backend_current_sessions{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmalertmanager-k8s"} 0
haproxy_backend_current_sessions{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmsingle-k8s"} 0
haproxy_backend_current_sessions{backend="http",namespace="rabbitmq",route="rabbitmq-fnzhp"} 0
haproxy_backend_current_sessions{backend="http",namespace="shared-license-distributor",route="shared-license-distributor-wnj68"} 0
haproxy_backend_current_sessions{backend="http",namespace="spark",route="spark-history-server-hqshb"} 0
haproxy_backend_current_sessions{backend="http",namespace="spark",route="spark-site-manager-p4nqh"} 0
haproxy_backend_current_sessions{backend="https",namespace="openshift-console",route="console"} 0
haproxy_backend_current_sessions{backend="https",namespace="openshift-monitoring",route="alertmanager-main"} 0
haproxy_backend_current_sessions{backend="https",namespace="openshift-monitoring",route="prometheus-k8s"} 0
haproxy_backend_current_sessions{backend="https",namespace="openshift-monitoring",route="prometheus-k8s-federate"} 0
haproxy_backend_current_sessions{backend="https",namespace="openshift-monitoring",route="thanos-querier"} 0
haproxy_backend_current_sessions{backend="https-edge",namespace="openshift-console",route="downloads"} 0
haproxy_backend_current_sessions{backend="https-edge",namespace="openshift-ingress-canary",route="canary"} 0
haproxy_backend_current_sessions{backend="other/be_no_sni",namespace="",route=""} 0
haproxy_backend_current_sessions{backend="other/be_sni",namespace="",route=""} 0
haproxy_backend_current_sessions{backend="other/openshift_default",namespace="",route=""} 0
haproxy_backend_current_sessions{backend="tcp",namespace="openshift-authentication",route="oauth-openshift"} 0
haproxy_backend_current_sessions{backend="tcp",namespace="site-manager",route="site-manager"} 0
# HELP haproxy_backend_http_average_connect_latency_milliseconds Average connect latency of the last 1024 requests in milliseconds.
# TYPE haproxy_backend_http_average_connect_latency_milliseconds gauge
haproxy_backend_http_average_connect_latency_milliseconds{backend="http",namespace="airflow",route="airflow-ingress-nzbtm"} 0
haproxy_backend_http_average_connect_latency_milliseconds{backend="http",namespace="airflow",route="airflow-site-manager-hpl66"} 0
haproxy_backend_http_average_connect_latency_milliseconds{backend="http",namespace="arangodb",route="arangodb-leader-ingress-zw82p"} 0
haproxy_backend_http_average_connect_latency_milliseconds{backend="http",namespace="consul-service",route="consul-ingress-hmmvf"} 0
haproxy_backend_http_average_connect_latency_milliseconds{backend="http",namespace="dbaas",route="aggregator-fw5ft"} 0
haproxy_backend_http_average_connect_latency_milliseconds{backend="http",namespace="infra-keycloak",route="infra-keycloak-5q4rj"} 0
haproxy_backend_http_average_connect_latency_milliseconds{backend="http",namespace="jaeger",route="jaeger-hotrod-qq5cq"} 0
haproxy_backend_http_average_connect_latency_milliseconds{backend="http",namespace="jaeger",route="jaeger-query-pc6w5"} 0
haproxy_backend_http_average_connect_latency_milliseconds{backend="http",namespace="kafka-service",route="akhq-ingress-lhz5z"} 0
haproxy_backend_http_average_connect_latency_milliseconds{backend="http",namespace="license-server",route="license-server-9txj9"} 0
haproxy_backend_http_average_connect_latency_milliseconds{backend="http",namespace="mistral",route="mistral-vslnb"} 0
haproxy_backend_http_average_connect_latency_milliseconds{backend="http",namespace="opensearch",route="opensearch-dashboards-l6vdv"} 0
haproxy_backend_http_average_connect_latency_milliseconds{backend="http",namespace="profiler",route="esc-collector-service-kdvdw"} 0
haproxy_backend_http_average_connect_latency_milliseconds{backend="http",namespace="profiler",route="esc-static-service-j9sgv"} 0
haproxy_backend_http_average_connect_latency_milliseconds{backend="http",namespace="profiler",route="esc-test-service-lbp8r"} 0
haproxy_backend_http_average_connect_latency_milliseconds{backend="http",namespace="profiler",route="esc-ui-service-r5tfm"} 0
haproxy_backend_http_average_connect_latency_milliseconds{backend="http",namespace="prometheus-operator",route="prometheus-operator-grafana"} 0
haproxy_backend_http_average_connect_latency_milliseconds{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmagent-k8s"} 0
haproxy_backend_http_average_connect_latency_milliseconds{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmalert-k8s"} 0
haproxy_backend_http_average_connect_latency_milliseconds{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmalertmanager-k8s"} 0
haproxy_backend_http_average_connect_latency_milliseconds{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmsingle-k8s"} 0
haproxy_backend_http_average_connect_latency_milliseconds{backend="http",namespace="rabbitmq",route="rabbitmq-fnzhp"} 0
haproxy_backend_http_average_connect_latency_milliseconds{backend="http",namespace="shared-license-distributor",route="shared-license-distributor-wnj68"} 0
haproxy_backend_http_average_connect_latency_milliseconds{backend="http",namespace="spark",route="spark-history-server-hqshb"} 0
haproxy_backend_http_average_connect_latency_milliseconds{backend="http",namespace="spark",route="spark-site-manager-p4nqh"} 0
haproxy_backend_http_average_connect_latency_milliseconds{backend="https",namespace="openshift-console",route="console"} 0
haproxy_backend_http_average_connect_latency_milliseconds{backend="https",namespace="openshift-monitoring",route="alertmanager-main"} 0
haproxy_backend_http_average_connect_latency_milliseconds{backend="https",namespace="openshift-monitoring",route="prometheus-k8s"} 0
haproxy_backend_http_average_connect_latency_milliseconds{backend="https",namespace="openshift-monitoring",route="prometheus-k8s-federate"} 0
haproxy_backend_http_average_connect_latency_milliseconds{backend="https",namespace="openshift-monitoring",route="thanos-querier"} 0
haproxy_backend_http_average_connect_latency_milliseconds{backend="https-edge",namespace="openshift-console",route="downloads"} 0
haproxy_backend_http_average_connect_latency_milliseconds{backend="https-edge",namespace="openshift-ingress-canary",route="canary"} 0
haproxy_backend_http_average_connect_latency_milliseconds{backend="other/be_no_sni",namespace="",route=""} 0
haproxy_backend_http_average_connect_latency_milliseconds{backend="other/be_sni",namespace="",route=""} 0
haproxy_backend_http_average_connect_latency_milliseconds{backend="other/openshift_default",namespace="",route=""} 0
haproxy_backend_http_average_connect_latency_milliseconds{backend="tcp",namespace="openshift-authentication",route="oauth-openshift"} 0
haproxy_backend_http_average_connect_latency_milliseconds{backend="tcp",namespace="site-manager",route="site-manager"} 0
# HELP haproxy_backend_http_average_queue_latency_milliseconds Average latency to be dequeued of the last 1024 requests in milliseconds.
# TYPE haproxy_backend_http_average_queue_latency_milliseconds gauge
haproxy_backend_http_average_queue_latency_milliseconds{backend="http",namespace="airflow",route="airflow-ingress-nzbtm"} 0
haproxy_backend_http_average_queue_latency_milliseconds{backend="http",namespace="airflow",route="airflow-site-manager-hpl66"} 0
haproxy_backend_http_average_queue_latency_milliseconds{backend="http",namespace="arangodb",route="arangodb-leader-ingress-zw82p"} 0
haproxy_backend_http_average_queue_latency_milliseconds{backend="http",namespace="consul-service",route="consul-ingress-hmmvf"} 0
haproxy_backend_http_average_queue_latency_milliseconds{backend="http",namespace="dbaas",route="aggregator-fw5ft"} 0
haproxy_backend_http_average_queue_latency_milliseconds{backend="http",namespace="infra-keycloak",route="infra-keycloak-5q4rj"} 0
haproxy_backend_http_average_queue_latency_milliseconds{backend="http",namespace="jaeger",route="jaeger-hotrod-qq5cq"} 0
haproxy_backend_http_average_queue_latency_milliseconds{backend="http",namespace="jaeger",route="jaeger-query-pc6w5"} 0
haproxy_backend_http_average_queue_latency_milliseconds{backend="http",namespace="kafka-service",route="akhq-ingress-lhz5z"} 0
haproxy_backend_http_average_queue_latency_milliseconds{backend="http",namespace="license-server",route="license-server-9txj9"} 0
haproxy_backend_http_average_queue_latency_milliseconds{backend="http",namespace="mistral",route="mistral-vslnb"} 0
haproxy_backend_http_average_queue_latency_milliseconds{backend="http",namespace="opensearch",route="opensearch-dashboards-l6vdv"} 0
haproxy_backend_http_average_queue_latency_milliseconds{backend="http",namespace="profiler",route="esc-collector-service-kdvdw"} 0
haproxy_backend_http_average_queue_latency_milliseconds{backend="http",namespace="profiler",route="esc-static-service-j9sgv"} 0
haproxy_backend_http_average_queue_latency_milliseconds{backend="http",namespace="profiler",route="esc-test-service-lbp8r"} 0
haproxy_backend_http_average_queue_latency_milliseconds{backend="http",namespace="profiler",route="esc-ui-service-r5tfm"} 0
haproxy_backend_http_average_queue_latency_milliseconds{backend="http",namespace="prometheus-operator",route="prometheus-operator-grafana"} 0
haproxy_backend_http_average_queue_latency_milliseconds{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmagent-k8s"} 0
haproxy_backend_http_average_queue_latency_milliseconds{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmalert-k8s"} 0
haproxy_backend_http_average_queue_latency_milliseconds{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmalertmanager-k8s"} 0
haproxy_backend_http_average_queue_latency_milliseconds{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmsingle-k8s"} 0
haproxy_backend_http_average_queue_latency_milliseconds{backend="http",namespace="rabbitmq",route="rabbitmq-fnzhp"} 0
haproxy_backend_http_average_queue_latency_milliseconds{backend="http",namespace="shared-license-distributor",route="shared-license-distributor-wnj68"} 0
haproxy_backend_http_average_queue_latency_milliseconds{backend="http",namespace="spark",route="spark-history-server-hqshb"} 0
haproxy_backend_http_average_queue_latency_milliseconds{backend="http",namespace="spark",route="spark-site-manager-p4nqh"} 0
haproxy_backend_http_average_queue_latency_milliseconds{backend="https",namespace="openshift-console",route="console"} 0
haproxy_backend_http_average_queue_latency_milliseconds{backend="https",namespace="openshift-monitoring",route="alertmanager-main"} 0
haproxy_backend_http_average_queue_latency_milliseconds{backend="https",namespace="openshift-monitoring",route="prometheus-k8s"} 0
haproxy_backend_http_average_queue_latency_milliseconds{backend="https",namespace="openshift-monitoring",route="prometheus-k8s-federate"} 0
haproxy_backend_http_average_queue_latency_milliseconds{backend="https",namespace="openshift-monitoring",route="thanos-querier"} 0
haproxy_backend_http_average_queue_latency_milliseconds{backend="https-edge",namespace="openshift-console",route="downloads"} 0
haproxy_backend_http_average_queue_latency_milliseconds{backend="https-edge",namespace="openshift-ingress-canary",route="canary"} 0
haproxy_backend_http_average_queue_latency_milliseconds{backend="other/be_no_sni",namespace="",route=""} 0
haproxy_backend_http_average_queue_latency_milliseconds{backend="other/be_sni",namespace="",route=""} 0
haproxy_backend_http_average_queue_latency_milliseconds{backend="other/openshift_default",namespace="",route=""} 0
haproxy_backend_http_average_queue_latency_milliseconds{backend="tcp",namespace="openshift-authentication",route="oauth-openshift"} 0
haproxy_backend_http_average_queue_latency_milliseconds{backend="tcp",namespace="site-manager",route="site-manager"} 0
# HELP haproxy_backend_http_average_response_latency_milliseconds Average response latency of the last 1024 requests in milliseconds.
# TYPE haproxy_backend_http_average_response_latency_milliseconds gauge
haproxy_backend_http_average_response_latency_milliseconds{backend="http",namespace="airflow",route="airflow-ingress-nzbtm"} 0
haproxy_backend_http_average_response_latency_milliseconds{backend="http",namespace="airflow",route="airflow-site-manager-hpl66"} 0
haproxy_backend_http_average_response_latency_milliseconds{backend="http",namespace="arangodb",route="arangodb-leader-ingress-zw82p"} 0
haproxy_backend_http_average_response_latency_milliseconds{backend="http",namespace="consul-service",route="consul-ingress-hmmvf"} 0
haproxy_backend_http_average_response_latency_milliseconds{backend="http",namespace="dbaas",route="aggregator-fw5ft"} 0
haproxy_backend_http_average_response_latency_milliseconds{backend="http",namespace="infra-keycloak",route="infra-keycloak-5q4rj"} 0
haproxy_backend_http_average_response_latency_milliseconds{backend="http",namespace="jaeger",route="jaeger-hotrod-qq5cq"} 0
haproxy_backend_http_average_response_latency_milliseconds{backend="http",namespace="jaeger",route="jaeger-query-pc6w5"} 0
haproxy_backend_http_average_response_latency_milliseconds{backend="http",namespace="kafka-service",route="akhq-ingress-lhz5z"} 0
haproxy_backend_http_average_response_latency_milliseconds{backend="http",namespace="license-server",route="license-server-9txj9"} 0
haproxy_backend_http_average_response_latency_milliseconds{backend="http",namespace="mistral",route="mistral-vslnb"} 0
haproxy_backend_http_average_response_latency_milliseconds{backend="http",namespace="opensearch",route="opensearch-dashboards-l6vdv"} 0
haproxy_backend_http_average_response_latency_milliseconds{backend="http",namespace="profiler",route="esc-collector-service-kdvdw"} 0
haproxy_backend_http_average_response_latency_milliseconds{backend="http",namespace="profiler",route="esc-static-service-j9sgv"} 0
haproxy_backend_http_average_response_latency_milliseconds{backend="http",namespace="profiler",route="esc-test-service-lbp8r"} 0
haproxy_backend_http_average_response_latency_milliseconds{backend="http",namespace="profiler",route="esc-ui-service-r5tfm"} 0
haproxy_backend_http_average_response_latency_milliseconds{backend="http",namespace="prometheus-operator",route="prometheus-operator-grafana"} 0
haproxy_backend_http_average_response_latency_milliseconds{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmagent-k8s"} 0
haproxy_backend_http_average_response_latency_milliseconds{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmalert-k8s"} 0
haproxy_backend_http_average_response_latency_milliseconds{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmalertmanager-k8s"} 0
haproxy_backend_http_average_response_latency_milliseconds{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmsingle-k8s"} 0
haproxy_backend_http_average_response_latency_milliseconds{backend="http",namespace="rabbitmq",route="rabbitmq-fnzhp"} 0
haproxy_backend_http_average_response_latency_milliseconds{backend="http",namespace="shared-license-distributor",route="shared-license-distributor-wnj68"} 0
haproxy_backend_http_average_response_latency_milliseconds{backend="http",namespace="spark",route="spark-history-server-hqshb"} 0
haproxy_backend_http_average_response_latency_milliseconds{backend="http",namespace="spark",route="spark-site-manager-p4nqh"} 0
haproxy_backend_http_average_response_latency_milliseconds{backend="https",namespace="openshift-console",route="console"} 0
haproxy_backend_http_average_response_latency_milliseconds{backend="https",namespace="openshift-monitoring",route="alertmanager-main"} 0
haproxy_backend_http_average_response_latency_milliseconds{backend="https",namespace="openshift-monitoring",route="prometheus-k8s"} 0
haproxy_backend_http_average_response_latency_milliseconds{backend="https",namespace="openshift-monitoring",route="prometheus-k8s-federate"} 0
haproxy_backend_http_average_response_latency_milliseconds{backend="https",namespace="openshift-monitoring",route="thanos-querier"} 0
haproxy_backend_http_average_response_latency_milliseconds{backend="https-edge",namespace="openshift-console",route="downloads"} 0
haproxy_backend_http_average_response_latency_milliseconds{backend="https-edge",namespace="openshift-ingress-canary",route="canary"} 0
haproxy_backend_http_average_response_latency_milliseconds{backend="other/be_no_sni",namespace="",route=""} 0
haproxy_backend_http_average_response_latency_milliseconds{backend="other/be_sni",namespace="",route=""} 0
haproxy_backend_http_average_response_latency_milliseconds{backend="other/openshift_default",namespace="",route=""} 0
haproxy_backend_http_average_response_latency_milliseconds{backend="tcp",namespace="openshift-authentication",route="oauth-openshift"} 0
haproxy_backend_http_average_response_latency_milliseconds{backend="tcp",namespace="site-manager",route="site-manager"} 0
# HELP haproxy_backend_http_responses_total Total of HTTP responses.
# TYPE haproxy_backend_http_responses_total gauge
haproxy_backend_http_responses_total{backend="http",code="1xx",namespace="airflow",route="airflow-ingress-nzbtm"} 0
haproxy_backend_http_responses_total{backend="http",code="1xx",namespace="airflow",route="airflow-site-manager-hpl66"} 0
haproxy_backend_http_responses_total{backend="http",code="1xx",namespace="arangodb",route="arangodb-leader-ingress-zw82p"} 0
haproxy_backend_http_responses_total{backend="http",code="1xx",namespace="consul-service",route="consul-ingress-hmmvf"} 0
haproxy_backend_http_responses_total{backend="http",code="1xx",namespace="dbaas",route="aggregator-fw5ft"} 0
haproxy_backend_http_responses_total{backend="http",code="1xx",namespace="infra-keycloak",route="infra-keycloak-5q4rj"} 0
haproxy_backend_http_responses_total{backend="http",code="1xx",namespace="jaeger",route="jaeger-hotrod-qq5cq"} 0
haproxy_backend_http_responses_total{backend="http",code="1xx",namespace="jaeger",route="jaeger-query-pc6w5"} 0
haproxy_backend_http_responses_total{backend="http",code="1xx",namespace="kafka-service",route="akhq-ingress-lhz5z"} 0
haproxy_backend_http_responses_total{backend="http",code="1xx",namespace="license-server",route="license-server-9txj9"} 0
haproxy_backend_http_responses_total{backend="http",code="1xx",namespace="mistral",route="mistral-vslnb"} 0
haproxy_backend_http_responses_total{backend="http",code="1xx",namespace="opensearch",route="opensearch-dashboards-l6vdv"} 0
haproxy_backend_http_responses_total{backend="http",code="1xx",namespace="profiler",route="esc-collector-service-kdvdw"} 0
haproxy_backend_http_responses_total{backend="http",code="1xx",namespace="profiler",route="esc-static-service-j9sgv"} 0
haproxy_backend_http_responses_total{backend="http",code="1xx",namespace="profiler",route="esc-test-service-lbp8r"} 0
haproxy_backend_http_responses_total{backend="http",code="1xx",namespace="profiler",route="esc-ui-service-r5tfm"} 0
haproxy_backend_http_responses_total{backend="http",code="1xx",namespace="prometheus-operator",route="prometheus-operator-grafana"} 0
haproxy_backend_http_responses_total{backend="http",code="1xx",namespace="prometheus-operator",route="prometheus-operator-vmagent-k8s"} 0
haproxy_backend_http_responses_total{backend="http",code="1xx",namespace="prometheus-operator",route="prometheus-operator-vmalert-k8s"} 0
haproxy_backend_http_responses_total{backend="http",code="1xx",namespace="prometheus-operator",route="prometheus-operator-vmalertmanager-k8s"} 0
haproxy_backend_http_responses_total{backend="http",code="1xx",namespace="prometheus-operator",route="prometheus-operator-vmsingle-k8s"} 0
haproxy_backend_http_responses_total{backend="http",code="1xx",namespace="rabbitmq",route="rabbitmq-fnzhp"} 0
haproxy_backend_http_responses_total{backend="http",code="1xx",namespace="shared-license-distributor",route="shared-license-distributor-wnj68"} 0
haproxy_backend_http_responses_total{backend="http",code="1xx",namespace="spark",route="spark-history-server-hqshb"} 0
haproxy_backend_http_responses_total{backend="http",code="1xx",namespace="spark",route="spark-site-manager-p4nqh"} 0
haproxy_backend_http_responses_total{backend="http",code="2xx",namespace="airflow",route="airflow-ingress-nzbtm"} 0
haproxy_backend_http_responses_total{backend="http",code="2xx",namespace="airflow",route="airflow-site-manager-hpl66"} 0
haproxy_backend_http_responses_total{backend="http",code="2xx",namespace="arangodb",route="arangodb-leader-ingress-zw82p"} 0
haproxy_backend_http_responses_total{backend="http",code="2xx",namespace="consul-service",route="consul-ingress-hmmvf"} 0
haproxy_backend_http_responses_total{backend="http",code="2xx",namespace="dbaas",route="aggregator-fw5ft"} 0
haproxy_backend_http_responses_total{backend="http",code="2xx",namespace="infra-keycloak",route="infra-keycloak-5q4rj"} 0
haproxy_backend_http_responses_total{backend="http",code="2xx",namespace="jaeger",route="jaeger-hotrod-qq5cq"} 0
haproxy_backend_http_responses_total{backend="http",code="2xx",namespace="jaeger",route="jaeger-query-pc6w5"} 0
haproxy_backend_http_responses_total{backend="http",code="2xx",namespace="kafka-service",route="akhq-ingress-lhz5z"} 20
haproxy_backend_http_responses_total{backend="http",code="2xx",namespace="license-server",route="license-server-9txj9"} 0
haproxy_backend_http_responses_total{backend="http",code="2xx",namespace="mistral",route="mistral-vslnb"} 0
haproxy_backend_http_responses_total{backend="http",code="2xx",namespace="opensearch",route="opensearch-dashboards-l6vdv"} 0
haproxy_backend_http_responses_total{backend="http",code="2xx",namespace="profiler",route="esc-collector-service-kdvdw"} 0
haproxy_backend_http_responses_total{backend="http",code="2xx",namespace="profiler",route="esc-static-service-j9sgv"} 0
haproxy_backend_http_responses_total{backend="http",code="2xx",namespace="profiler",route="esc-test-service-lbp8r"} 0
haproxy_backend_http_responses_total{backend="http",code="2xx",namespace="profiler",route="esc-ui-service-r5tfm"} 0
haproxy_backend_http_responses_total{backend="http",code="2xx",namespace="prometheus-operator",route="prometheus-operator-grafana"} 0
haproxy_backend_http_responses_total{backend="http",code="2xx",namespace="prometheus-operator",route="prometheus-operator-vmagent-k8s"} 0
haproxy_backend_http_responses_total{backend="http",code="2xx",namespace="prometheus-operator",route="prometheus-operator-vmalert-k8s"} 0
haproxy_backend_http_responses_total{backend="http",code="2xx",namespace="prometheus-operator",route="prometheus-operator-vmalertmanager-k8s"} 0
haproxy_backend_http_responses_total{backend="http",code="2xx",namespace="prometheus-operator",route="prometheus-operator-vmsingle-k8s"} 36
haproxy_backend_http_responses_total{backend="http",code="2xx",namespace="rabbitmq",route="rabbitmq-fnzhp"} 0
haproxy_backend_http_responses_total{backend="http",code="2xx",namespace="shared-license-distributor",route="shared-license-distributor-wnj68"} 0
haproxy_backend_http_responses_total{backend="http",code="2xx",namespace="spark",route="spark-history-server-hqshb"} 0
haproxy_backend_http_responses_total{backend="http",code="2xx",namespace="spark",route="spark-site-manager-p4nqh"} 0
haproxy_backend_http_responses_total{backend="http",code="3xx",namespace="airflow",route="airflow-ingress-nzbtm"} 0
haproxy_backend_http_responses_total{backend="http",code="3xx",namespace="airflow",route="airflow-site-manager-hpl66"} 0
haproxy_backend_http_responses_total{backend="http",code="3xx",namespace="arangodb",route="arangodb-leader-ingress-zw82p"} 0
haproxy_backend_http_responses_total{backend="http",code="3xx",namespace="consul-service",route="consul-ingress-hmmvf"} 0
haproxy_backend_http_responses_total{backend="http",code="3xx",namespace="dbaas",route="aggregator-fw5ft"} 0
haproxy_backend_http_responses_total{backend="http",code="3xx",namespace="infra-keycloak",route="infra-keycloak-5q4rj"} 0
haproxy_backend_http_responses_total{backend="http",code="3xx",namespace="jaeger",route="jaeger-hotrod-qq5cq"} 0
haproxy_backend_http_responses_total{backend="http",code="3xx",namespace="jaeger",route="jaeger-query-pc6w5"} 0
haproxy_backend_http_responses_total{backend="http",code="3xx",namespace="kafka-service",route="akhq-ingress-lhz5z"} 2
haproxy_backend_http_responses_total{backend="http",code="3xx",namespace="license-server",route="license-server-9txj9"} 0
haproxy_backend_http_responses_total{backend="http",code="3xx",namespace="mistral",route="mistral-vslnb"} 0
haproxy_backend_http_responses_total{backend="http",code="3xx",namespace="opensearch",route="opensearch-dashboards-l6vdv"} 0
haproxy_backend_http_responses_total{backend="http",code="3xx",namespace="profiler",route="esc-collector-service-kdvdw"} 0
haproxy_backend_http_responses_total{backend="http",code="3xx",namespace="profiler",route="esc-static-service-j9sgv"} 0
haproxy_backend_http_responses_total{backend="http",code="3xx",namespace="profiler",route="esc-test-service-lbp8r"} 0
haproxy_backend_http_responses_total{backend="http",code="3xx",namespace="profiler",route="esc-ui-service-r5tfm"} 0
haproxy_backend_http_responses_total{backend="http",code="3xx",namespace="prometheus-operator",route="prometheus-operator-grafana"} 0
haproxy_backend_http_responses_total{backend="http",code="3xx",namespace="prometheus-operator",route="prometheus-operator-vmagent-k8s"} 0
haproxy_backend_http_responses_total{backend="http",code="3xx",namespace="prometheus-operator",route="prometheus-operator-vmalert-k8s"} 0
haproxy_backend_http_responses_total{backend="http",code="3xx",namespace="prometheus-operator",route="prometheus-operator-vmalertmanager-k8s"} 0
haproxy_backend_http_responses_total{backend="http",code="3xx",namespace="prometheus-operator",route="prometheus-operator-vmsingle-k8s"} 0
haproxy_backend_http_responses_total{backend="http",code="3xx",namespace="rabbitmq",route="rabbitmq-fnzhp"} 0
haproxy_backend_http_responses_total{backend="http",code="3xx",namespace="shared-license-distributor",route="shared-license-distributor-wnj68"} 0
haproxy_backend_http_responses_total{backend="http",code="3xx",namespace="spark",route="spark-history-server-hqshb"} 0
haproxy_backend_http_responses_total{backend="http",code="3xx",namespace="spark",route="spark-site-manager-p4nqh"} 0
haproxy_backend_http_responses_total{backend="http",code="4xx",namespace="airflow",route="airflow-ingress-nzbtm"} 0
haproxy_backend_http_responses_total{backend="http",code="4xx",namespace="airflow",route="airflow-site-manager-hpl66"} 0
haproxy_backend_http_responses_total{backend="http",code="4xx",namespace="arangodb",route="arangodb-leader-ingress-zw82p"} 0
haproxy_backend_http_responses_total{backend="http",code="4xx",namespace="consul-service",route="consul-ingress-hmmvf"} 0
haproxy_backend_http_responses_total{backend="http",code="4xx",namespace="dbaas",route="aggregator-fw5ft"} 0
haproxy_backend_http_responses_total{backend="http",code="4xx",namespace="infra-keycloak",route="infra-keycloak-5q4rj"} 0
haproxy_backend_http_responses_total{backend="http",code="4xx",namespace="jaeger",route="jaeger-hotrod-qq5cq"} 0
haproxy_backend_http_responses_total{backend="http",code="4xx",namespace="jaeger",route="jaeger-query-pc6w5"} 0
haproxy_backend_http_responses_total{backend="http",code="4xx",namespace="kafka-service",route="akhq-ingress-lhz5z"} 0
haproxy_backend_http_responses_total{backend="http",code="4xx",namespace="license-server",route="license-server-9txj9"} 0
haproxy_backend_http_responses_total{backend="http",code="4xx",namespace="mistral",route="mistral-vslnb"} 0
haproxy_backend_http_responses_total{backend="http",code="4xx",namespace="opensearch",route="opensearch-dashboards-l6vdv"} 0
haproxy_backend_http_responses_total{backend="http",code="4xx",namespace="profiler",route="esc-collector-service-kdvdw"} 0
haproxy_backend_http_responses_total{backend="http",code="4xx",namespace="profiler",route="esc-static-service-j9sgv"} 0
haproxy_backend_http_responses_total{backend="http",code="4xx",namespace="profiler",route="esc-test-service-lbp8r"} 0
haproxy_backend_http_responses_total{backend="http",code="4xx",namespace="profiler",route="esc-ui-service-r5tfm"} 0
haproxy_backend_http_responses_total{backend="http",code="4xx",namespace="prometheus-operator",route="prometheus-operator-grafana"} 0
haproxy_backend_http_responses_total{backend="http",code="4xx",namespace="prometheus-operator",route="prometheus-operator-vmagent-k8s"} 0
haproxy_backend_http_responses_total{backend="http",code="4xx",namespace="prometheus-operator",route="prometheus-operator-vmalert-k8s"} 0
haproxy_backend_http_responses_total{backend="http",code="4xx",namespace="prometheus-operator",route="prometheus-operator-vmalertmanager-k8s"} 0
haproxy_backend_http_responses_total{backend="http",code="4xx",namespace="prometheus-operator",route="prometheus-operator-vmsingle-k8s"} 0
haproxy_backend_http_responses_total{backend="http",code="4xx",namespace="rabbitmq",route="rabbitmq-fnzhp"} 0
haproxy_backend_http_responses_total{backend="http",code="4xx",namespace="shared-license-distributor",route="shared-license-distributor-wnj68"} 0
haproxy_backend_http_responses_total{backend="http",code="4xx",namespace="spark",route="spark-history-server-hqshb"} 0
haproxy_backend_http_responses_total{backend="http",code="4xx",namespace="spark",route="spark-site-manager-p4nqh"} 0
haproxy_backend_http_responses_total{backend="http",code="5xx",namespace="airflow",route="airflow-ingress-nzbtm"} 0
haproxy_backend_http_responses_total{backend="http",code="5xx",namespace="airflow",route="airflow-site-manager-hpl66"} 0
haproxy_backend_http_responses_total{backend="http",code="5xx",namespace="arangodb",route="arangodb-leader-ingress-zw82p"} 0
haproxy_backend_http_responses_total{backend="http",code="5xx",namespace="consul-service",route="consul-ingress-hmmvf"} 0
haproxy_backend_http_responses_total{backend="http",code="5xx",namespace="dbaas",route="aggregator-fw5ft"} 0
haproxy_backend_http_responses_total{backend="http",code="5xx",namespace="infra-keycloak",route="infra-keycloak-5q4rj"} 0
haproxy_backend_http_responses_total{backend="http",code="5xx",namespace="jaeger",route="jaeger-hotrod-qq5cq"} 0
haproxy_backend_http_responses_total{backend="http",code="5xx",namespace="jaeger",route="jaeger-query-pc6w5"} 0
haproxy_backend_http_responses_total{backend="http",code="5xx",namespace="kafka-service",route="akhq-ingress-lhz5z"} 0
haproxy_backend_http_responses_total{backend="http",code="5xx",namespace="license-server",route="license-server-9txj9"} 0
haproxy_backend_http_responses_total{backend="http",code="5xx",namespace="mistral",route="mistral-vslnb"} 0
haproxy_backend_http_responses_total{backend="http",code="5xx",namespace="opensearch",route="opensearch-dashboards-l6vdv"} 0
haproxy_backend_http_responses_total{backend="http",code="5xx",namespace="profiler",route="esc-collector-service-kdvdw"} 0
haproxy_backend_http_responses_total{backend="http",code="5xx",namespace="profiler",route="esc-static-service-j9sgv"} 0
haproxy_backend_http_responses_total{backend="http",code="5xx",namespace="profiler",route="esc-test-service-lbp8r"} 0
haproxy_backend_http_responses_total{backend="http",code="5xx",namespace="profiler",route="esc-ui-service-r5tfm"} 0
haproxy_backend_http_responses_total{backend="http",code="5xx",namespace="prometheus-operator",route="prometheus-operator-grafana"} 0
haproxy_backend_http_responses_total{backend="http",code="5xx",namespace="prometheus-operator",route="prometheus-operator-vmagent-k8s"} 0
haproxy_backend_http_responses_total{backend="http",code="5xx",namespace="prometheus-operator",route="prometheus-operator-vmalert-k8s"} 0
haproxy_backend_http_responses_total{backend="http",code="5xx",namespace="prometheus-operator",route="prometheus-operator-vmalertmanager-k8s"} 0
haproxy_backend_http_responses_total{backend="http",code="5xx",namespace="prometheus-operator",route="prometheus-operator-vmsingle-k8s"} 0
haproxy_backend_http_responses_total{backend="http",code="5xx",namespace="rabbitmq",route="rabbitmq-fnzhp"} 0
haproxy_backend_http_responses_total{backend="http",code="5xx",namespace="shared-license-distributor",route="shared-license-distributor-wnj68"} 0
haproxy_backend_http_responses_total{backend="http",code="5xx",namespace="spark",route="spark-history-server-hqshb"} 0
haproxy_backend_http_responses_total{backend="http",code="5xx",namespace="spark",route="spark-site-manager-p4nqh"} 0
haproxy_backend_http_responses_total{backend="http",code="other",namespace="airflow",route="airflow-ingress-nzbtm"} 0
haproxy_backend_http_responses_total{backend="http",code="other",namespace="airflow",route="airflow-site-manager-hpl66"} 0
haproxy_backend_http_responses_total{backend="http",code="other",namespace="arangodb",route="arangodb-leader-ingress-zw82p"} 0
haproxy_backend_http_responses_total{backend="http",code="other",namespace="consul-service",route="consul-ingress-hmmvf"} 0
haproxy_backend_http_responses_total{backend="http",code="other",namespace="dbaas",route="aggregator-fw5ft"} 0
haproxy_backend_http_responses_total{backend="http",code="other",namespace="infra-keycloak",route="infra-keycloak-5q4rj"} 0
haproxy_backend_http_responses_total{backend="http",code="other",namespace="jaeger",route="jaeger-hotrod-qq5cq"} 0
haproxy_backend_http_responses_total{backend="http",code="other",namespace="jaeger",route="jaeger-query-pc6w5"} 0
haproxy_backend_http_responses_total{backend="http",code="other",namespace="kafka-service",route="akhq-ingress-lhz5z"} 0
haproxy_backend_http_responses_total{backend="http",code="other",namespace="license-server",route="license-server-9txj9"} 0
haproxy_backend_http_responses_total{backend="http",code="other",namespace="mistral",route="mistral-vslnb"} 0
haproxy_backend_http_responses_total{backend="http",code="other",namespace="opensearch",route="opensearch-dashboards-l6vdv"} 0
haproxy_backend_http_responses_total{backend="http",code="other",namespace="profiler",route="esc-collector-service-kdvdw"} 0
haproxy_backend_http_responses_total{backend="http",code="other",namespace="profiler",route="esc-static-service-j9sgv"} 0
haproxy_backend_http_responses_total{backend="http",code="other",namespace="profiler",route="esc-test-service-lbp8r"} 0
haproxy_backend_http_responses_total{backend="http",code="other",namespace="profiler",route="esc-ui-service-r5tfm"} 0
haproxy_backend_http_responses_total{backend="http",code="other",namespace="prometheus-operator",route="prometheus-operator-grafana"} 0
haproxy_backend_http_responses_total{backend="http",code="other",namespace="prometheus-operator",route="prometheus-operator-vmagent-k8s"} 0
haproxy_backend_http_responses_total{backend="http",code="other",namespace="prometheus-operator",route="prometheus-operator-vmalert-k8s"} 0
haproxy_backend_http_responses_total{backend="http",code="other",namespace="prometheus-operator",route="prometheus-operator-vmalertmanager-k8s"} 0
haproxy_backend_http_responses_total{backend="http",code="other",namespace="prometheus-operator",route="prometheus-operator-vmsingle-k8s"} 0
haproxy_backend_http_responses_total{backend="http",code="other",namespace="rabbitmq",route="rabbitmq-fnzhp"} 0
haproxy_backend_http_responses_total{backend="http",code="other",namespace="shared-license-distributor",route="shared-license-distributor-wnj68"} 0
haproxy_backend_http_responses_total{backend="http",code="other",namespace="spark",route="spark-history-server-hqshb"} 0
haproxy_backend_http_responses_total{backend="http",code="other",namespace="spark",route="spark-site-manager-p4nqh"} 0
haproxy_backend_http_responses_total{backend="https",code="1xx",namespace="openshift-console",route="console"} 453
haproxy_backend_http_responses_total{backend="https",code="1xx",namespace="openshift-monitoring",route="alertmanager-main"} 0
haproxy_backend_http_responses_total{backend="https",code="1xx",namespace="openshift-monitoring",route="prometheus-k8s"} 0
haproxy_backend_http_responses_total{backend="https",code="1xx",namespace="openshift-monitoring",route="prometheus-k8s-federate"} 0
haproxy_backend_http_responses_total{backend="https",code="1xx",namespace="openshift-monitoring",route="thanos-querier"} 0
haproxy_backend_http_responses_total{backend="https",code="2xx",namespace="openshift-console",route="console"} 59914
haproxy_backend_http_responses_total{backend="https",code="2xx",namespace="openshift-monitoring",route="alertmanager-main"} 0
haproxy_backend_http_responses_total{backend="https",code="2xx",namespace="openshift-monitoring",route="prometheus-k8s"} 0
haproxy_backend_http_responses_total{backend="https",code="2xx",namespace="openshift-monitoring",route="prometheus-k8s-federate"} 0
haproxy_backend_http_responses_total{backend="https",code="2xx",namespace="openshift-monitoring",route="thanos-querier"} 0
haproxy_backend_http_responses_total{backend="https",code="3xx",namespace="openshift-console",route="console"} 15
haproxy_backend_http_responses_total{backend="https",code="3xx",namespace="openshift-monitoring",route="alertmanager-main"} 0
haproxy_backend_http_responses_total{backend="https",code="3xx",namespace="openshift-monitoring",route="prometheus-k8s"} 0
haproxy_backend_http_responses_total{backend="https",code="3xx",namespace="openshift-monitoring",route="prometheus-k8s-federate"} 0
haproxy_backend_http_responses_total{backend="https",code="3xx",namespace="openshift-monitoring",route="thanos-querier"} 0
haproxy_backend_http_responses_total{backend="https",code="4xx",namespace="openshift-console",route="console"} 305
haproxy_backend_http_responses_total{backend="https",code="4xx",namespace="openshift-monitoring",route="alertmanager-main"} 0
haproxy_backend_http_responses_total{backend="https",code="4xx",namespace="openshift-monitoring",route="prometheus-k8s"} 0
haproxy_backend_http_responses_total{backend="https",code="4xx",namespace="openshift-monitoring",route="prometheus-k8s-federate"} 0
haproxy_backend_http_responses_total{backend="https",code="4xx",namespace="openshift-monitoring",route="thanos-querier"} 0
haproxy_backend_http_responses_total{backend="https",code="5xx",namespace="openshift-console",route="console"} 0
haproxy_backend_http_responses_total{backend="https",code="5xx",namespace="openshift-monitoring",route="alertmanager-main"} 0
haproxy_backend_http_responses_total{backend="https",code="5xx",namespace="openshift-monitoring",route="prometheus-k8s"} 0
haproxy_backend_http_responses_total{backend="https",code="5xx",namespace="openshift-monitoring",route="prometheus-k8s-federate"} 0
haproxy_backend_http_responses_total{backend="https",code="5xx",namespace="openshift-monitoring",route="thanos-querier"} 0
haproxy_backend_http_responses_total{backend="https",code="other",namespace="openshift-console",route="console"} 0
haproxy_backend_http_responses_total{backend="https",code="other",namespace="openshift-monitoring",route="alertmanager-main"} 0
haproxy_backend_http_responses_total{backend="https",code="other",namespace="openshift-monitoring",route="prometheus-k8s"} 0
haproxy_backend_http_responses_total{backend="https",code="other",namespace="openshift-monitoring",route="prometheus-k8s-federate"} 0
haproxy_backend_http_responses_total{backend="https",code="other",namespace="openshift-monitoring",route="thanos-querier"} 0
haproxy_backend_http_responses_total{backend="https-edge",code="1xx",namespace="openshift-console",route="downloads"} 0
haproxy_backend_http_responses_total{backend="https-edge",code="1xx",namespace="openshift-ingress-canary",route="canary"} 0
haproxy_backend_http_responses_total{backend="https-edge",code="2xx",namespace="openshift-console",route="downloads"} 0
haproxy_backend_http_responses_total{backend="https-edge",code="2xx",namespace="openshift-ingress-canary",route="canary"} 0
haproxy_backend_http_responses_total{backend="https-edge",code="3xx",namespace="openshift-console",route="downloads"} 0
haproxy_backend_http_responses_total{backend="https-edge",code="3xx",namespace="openshift-ingress-canary",route="canary"} 0
haproxy_backend_http_responses_total{backend="https-edge",code="4xx",namespace="openshift-console",route="downloads"} 0
haproxy_backend_http_responses_total{backend="https-edge",code="4xx",namespace="openshift-ingress-canary",route="canary"} 0
haproxy_backend_http_responses_total{backend="https-edge",code="5xx",namespace="openshift-console",route="downloads"} 0
haproxy_backend_http_responses_total{backend="https-edge",code="5xx",namespace="openshift-ingress-canary",route="canary"} 0
haproxy_backend_http_responses_total{backend="https-edge",code="other",namespace="openshift-console",route="downloads"} 0
haproxy_backend_http_responses_total{backend="https-edge",code="other",namespace="openshift-ingress-canary",route="canary"} 0
haproxy_backend_http_responses_total{backend="other/openshift_default",code="1xx",namespace="",route=""} 0
haproxy_backend_http_responses_total{backend="other/openshift_default",code="2xx",namespace="",route=""} 0
haproxy_backend_http_responses_total{backend="other/openshift_default",code="3xx",namespace="",route=""} 0
haproxy_backend_http_responses_total{backend="other/openshift_default",code="4xx",namespace="",route=""} 0
haproxy_backend_http_responses_total{backend="other/openshift_default",code="5xx",namespace="",route=""} 156
haproxy_backend_http_responses_total{backend="other/openshift_default",code="other",namespace="",route=""} 0
# HELP haproxy_backend_max_session_rate Maximum number of sessions per second.
# TYPE haproxy_backend_max_session_rate gauge
haproxy_backend_max_session_rate{backend="http",namespace="airflow",route="airflow-ingress-nzbtm"} 0
haproxy_backend_max_session_rate{backend="http",namespace="airflow",route="airflow-site-manager-hpl66"} 0
haproxy_backend_max_session_rate{backend="http",namespace="arangodb",route="arangodb-leader-ingress-zw82p"} 0
haproxy_backend_max_session_rate{backend="http",namespace="consul-service",route="consul-ingress-hmmvf"} 0
haproxy_backend_max_session_rate{backend="http",namespace="dbaas",route="aggregator-fw5ft"} 0
haproxy_backend_max_session_rate{backend="http",namespace="infra-keycloak",route="infra-keycloak-5q4rj"} 0
haproxy_backend_max_session_rate{backend="http",namespace="jaeger",route="jaeger-hotrod-qq5cq"} 0
haproxy_backend_max_session_rate{backend="http",namespace="jaeger",route="jaeger-query-pc6w5"} 0
haproxy_backend_max_session_rate{backend="http",namespace="kafka-service",route="akhq-ingress-lhz5z"} 0
haproxy_backend_max_session_rate{backend="http",namespace="license-server",route="license-server-9txj9"} 0
haproxy_backend_max_session_rate{backend="http",namespace="mistral",route="mistral-vslnb"} 0
haproxy_backend_max_session_rate{backend="http",namespace="opensearch",route="opensearch-dashboards-l6vdv"} 0
haproxy_backend_max_session_rate{backend="http",namespace="profiler",route="esc-collector-service-kdvdw"} 0
haproxy_backend_max_session_rate{backend="http",namespace="profiler",route="esc-static-service-j9sgv"} 0
haproxy_backend_max_session_rate{backend="http",namespace="profiler",route="esc-test-service-lbp8r"} 0
haproxy_backend_max_session_rate{backend="http",namespace="profiler",route="esc-ui-service-r5tfm"} 0
haproxy_backend_max_session_rate{backend="http",namespace="prometheus-operator",route="prometheus-operator-grafana"} 0
haproxy_backend_max_session_rate{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmagent-k8s"} 0
haproxy_backend_max_session_rate{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmalert-k8s"} 0
haproxy_backend_max_session_rate{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmalertmanager-k8s"} 0
haproxy_backend_max_session_rate{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmsingle-k8s"} 0
haproxy_backend_max_session_rate{backend="http",namespace="rabbitmq",route="rabbitmq-fnzhp"} 0
haproxy_backend_max_session_rate{backend="http",namespace="shared-license-distributor",route="shared-license-distributor-wnj68"} 0
haproxy_backend_max_session_rate{backend="http",namespace="spark",route="spark-history-server-hqshb"} 0
haproxy_backend_max_session_rate{backend="http",namespace="spark",route="spark-site-manager-p4nqh"} 0
haproxy_backend_max_session_rate{backend="https",namespace="openshift-console",route="console"} 0
haproxy_backend_max_session_rate{backend="https",namespace="openshift-monitoring",route="alertmanager-main"} 0
haproxy_backend_max_session_rate{backend="https",namespace="openshift-monitoring",route="prometheus-k8s"} 0
haproxy_backend_max_session_rate{backend="https",namespace="openshift-monitoring",route="prometheus-k8s-federate"} 0
haproxy_backend_max_session_rate{backend="https",namespace="openshift-monitoring",route="thanos-querier"} 0
haproxy_backend_max_session_rate{backend="https-edge",namespace="openshift-console",route="downloads"} 0
haproxy_backend_max_session_rate{backend="https-edge",namespace="openshift-ingress-canary",route="canary"} 0
haproxy_backend_max_session_rate{backend="other/be_no_sni",namespace="",route=""} 0
haproxy_backend_max_session_rate{backend="other/be_sni",namespace="",route=""} 0
haproxy_backend_max_session_rate{backend="other/openshift_default",namespace="",route=""} 1
haproxy_backend_max_session_rate{backend="tcp",namespace="openshift-authentication",route="oauth-openshift"} 0
haproxy_backend_max_session_rate{backend="tcp",namespace="site-manager",route="site-manager"} 0
# HELP haproxy_backend_max_sessions Maximum observed number of active sessions.
# TYPE haproxy_backend_max_sessions gauge
haproxy_backend_max_sessions{backend="http",namespace="airflow",route="airflow-ingress-nzbtm"} 0
haproxy_backend_max_sessions{backend="http",namespace="airflow",route="airflow-site-manager-hpl66"} 0
haproxy_backend_max_sessions{backend="http",namespace="arangodb",route="arangodb-leader-ingress-zw82p"} 0
haproxy_backend_max_sessions{backend="http",namespace="consul-service",route="consul-ingress-hmmvf"} 0
haproxy_backend_max_sessions{backend="http",namespace="dbaas",route="aggregator-fw5ft"} 0
haproxy_backend_max_sessions{backend="http",namespace="infra-keycloak",route="infra-keycloak-5q4rj"} 0
haproxy_backend_max_sessions{backend="http",namespace="jaeger",route="jaeger-hotrod-qq5cq"} 0
haproxy_backend_max_sessions{backend="http",namespace="jaeger",route="jaeger-query-pc6w5"} 0
haproxy_backend_max_sessions{backend="http",namespace="kafka-service",route="akhq-ingress-lhz5z"} 0
haproxy_backend_max_sessions{backend="http",namespace="license-server",route="license-server-9txj9"} 0
haproxy_backend_max_sessions{backend="http",namespace="mistral",route="mistral-vslnb"} 0
haproxy_backend_max_sessions{backend="http",namespace="opensearch",route="opensearch-dashboards-l6vdv"} 0
haproxy_backend_max_sessions{backend="http",namespace="profiler",route="esc-collector-service-kdvdw"} 0
haproxy_backend_max_sessions{backend="http",namespace="profiler",route="esc-static-service-j9sgv"} 0
haproxy_backend_max_sessions{backend="http",namespace="profiler",route="esc-test-service-lbp8r"} 0
haproxy_backend_max_sessions{backend="http",namespace="profiler",route="esc-ui-service-r5tfm"} 0
haproxy_backend_max_sessions{backend="http",namespace="prometheus-operator",route="prometheus-operator-grafana"} 0
haproxy_backend_max_sessions{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmagent-k8s"} 0
haproxy_backend_max_sessions{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmalert-k8s"} 0
haproxy_backend_max_sessions{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmalertmanager-k8s"} 0
haproxy_backend_max_sessions{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmsingle-k8s"} 0
haproxy_backend_max_sessions{backend="http",namespace="rabbitmq",route="rabbitmq-fnzhp"} 0
haproxy_backend_max_sessions{backend="http",namespace="shared-license-distributor",route="shared-license-distributor-wnj68"} 0
haproxy_backend_max_sessions{backend="http",namespace="spark",route="spark-history-server-hqshb"} 0
haproxy_backend_max_sessions{backend="http",namespace="spark",route="spark-site-manager-p4nqh"} 0
haproxy_backend_max_sessions{backend="https",namespace="openshift-console",route="console"} 0
haproxy_backend_max_sessions{backend="https",namespace="openshift-monitoring",route="alertmanager-main"} 0
haproxy_backend_max_sessions{backend="https",namespace="openshift-monitoring",route="prometheus-k8s"} 0
haproxy_backend_max_sessions{backend="https",namespace="openshift-monitoring",route="prometheus-k8s-federate"} 0
haproxy_backend_max_sessions{backend="https",namespace="openshift-monitoring",route="thanos-querier"} 0
haproxy_backend_max_sessions{backend="https-edge",namespace="openshift-console",route="downloads"} 0
haproxy_backend_max_sessions{backend="https-edge",namespace="openshift-ingress-canary",route="canary"} 0
haproxy_backend_max_sessions{backend="other/be_no_sni",namespace="",route=""} 0
haproxy_backend_max_sessions{backend="other/be_sni",namespace="",route=""} 0
haproxy_backend_max_sessions{backend="other/openshift_default",namespace="",route=""} 1
haproxy_backend_max_sessions{backend="tcp",namespace="openshift-authentication",route="oauth-openshift"} 0
haproxy_backend_max_sessions{backend="tcp",namespace="site-manager",route="site-manager"} 0
# HELP haproxy_backend_response_errors_total Total of response errors.
# TYPE haproxy_backend_response_errors_total gauge
haproxy_backend_response_errors_total{backend="http",namespace="airflow",route="airflow-ingress-nzbtm"} 0
haproxy_backend_response_errors_total{backend="http",namespace="airflow",route="airflow-site-manager-hpl66"} 0
haproxy_backend_response_errors_total{backend="http",namespace="arangodb",route="arangodb-leader-ingress-zw82p"} 0
haproxy_backend_response_errors_total{backend="http",namespace="consul-service",route="consul-ingress-hmmvf"} 0
haproxy_backend_response_errors_total{backend="http",namespace="dbaas",route="aggregator-fw5ft"} 0
haproxy_backend_response_errors_total{backend="http",namespace="infra-keycloak",route="infra-keycloak-5q4rj"} 0
haproxy_backend_response_errors_total{backend="http",namespace="jaeger",route="jaeger-hotrod-qq5cq"} 0
haproxy_backend_response_errors_total{backend="http",namespace="jaeger",route="jaeger-query-pc6w5"} 0
haproxy_backend_response_errors_total{backend="http",namespace="kafka-service",route="akhq-ingress-lhz5z"} 0
haproxy_backend_response_errors_total{backend="http",namespace="license-server",route="license-server-9txj9"} 0
haproxy_backend_response_errors_total{backend="http",namespace="mistral",route="mistral-vslnb"} 0
haproxy_backend_response_errors_total{backend="http",namespace="opensearch",route="opensearch-dashboards-l6vdv"} 0
haproxy_backend_response_errors_total{backend="http",namespace="profiler",route="esc-collector-service-kdvdw"} 0
haproxy_backend_response_errors_total{backend="http",namespace="profiler",route="esc-static-service-j9sgv"} 0
haproxy_backend_response_errors_total{backend="http",namespace="profiler",route="esc-test-service-lbp8r"} 0
haproxy_backend_response_errors_total{backend="http",namespace="profiler",route="esc-ui-service-r5tfm"} 0
haproxy_backend_response_errors_total{backend="http",namespace="prometheus-operator",route="prometheus-operator-grafana"} 0
haproxy_backend_response_errors_total{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmagent-k8s"} 0
haproxy_backend_response_errors_total{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmalert-k8s"} 0
haproxy_backend_response_errors_total{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmalertmanager-k8s"} 0
haproxy_backend_response_errors_total{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmsingle-k8s"} 0
haproxy_backend_response_errors_total{backend="http",namespace="rabbitmq",route="rabbitmq-fnzhp"} 0
haproxy_backend_response_errors_total{backend="http",namespace="shared-license-distributor",route="shared-license-distributor-wnj68"} 0
haproxy_backend_response_errors_total{backend="http",namespace="spark",route="spark-history-server-hqshb"} 0
haproxy_backend_response_errors_total{backend="http",namespace="spark",route="spark-site-manager-p4nqh"} 0
haproxy_backend_response_errors_total{backend="https",namespace="openshift-console",route="console"} 0
haproxy_backend_response_errors_total{backend="https",namespace="openshift-monitoring",route="alertmanager-main"} 0
haproxy_backend_response_errors_total{backend="https",namespace="openshift-monitoring",route="prometheus-k8s"} 0
haproxy_backend_response_errors_total{backend="https",namespace="openshift-monitoring",route="prometheus-k8s-federate"} 0
haproxy_backend_response_errors_total{backend="https",namespace="openshift-monitoring",route="thanos-querier"} 0
haproxy_backend_response_errors_total{backend="https-edge",namespace="openshift-console",route="downloads"} 0
haproxy_backend_response_errors_total{backend="https-edge",namespace="openshift-ingress-canary",route="canary"} 0
haproxy_backend_response_errors_total{backend="other/be_no_sni",namespace="",route=""} 0
haproxy_backend_response_errors_total{backend="other/be_sni",namespace="",route=""} 13
haproxy_backend_response_errors_total{backend="other/openshift_default",namespace="",route=""} 0
haproxy_backend_response_errors_total{backend="tcp",namespace="openshift-authentication",route="oauth-openshift"} 0
haproxy_backend_response_errors_total{backend="tcp",namespace="site-manager",route="site-manager"} 0
# HELP haproxy_backend_up Current health status of the backend (1 = UP, 0 = DOWN).
# TYPE haproxy_backend_up gauge
haproxy_backend_up{backend="http",namespace="airflow",route="airflow-ingress-nzbtm"} 1
haproxy_backend_up{backend="http",namespace="airflow",route="airflow-site-manager-hpl66"} 1
haproxy_backend_up{backend="http",namespace="arangodb",route="arangodb-leader-ingress-zw82p"} 1
haproxy_backend_up{backend="http",namespace="consul-service",route="consul-ingress-hmmvf"} 1
haproxy_backend_up{backend="http",namespace="dbaas",route="aggregator-fw5ft"} 1
haproxy_backend_up{backend="http",namespace="infra-keycloak",route="infra-keycloak-5q4rj"} 1
haproxy_backend_up{backend="http",namespace="jaeger",route="jaeger-hotrod-qq5cq"} 1
haproxy_backend_up{backend="http",namespace="jaeger",route="jaeger-query-pc6w5"} 1
haproxy_backend_up{backend="http",namespace="kafka-service",route="akhq-ingress-lhz5z"} 1
haproxy_backend_up{backend="http",namespace="license-server",route="license-server-9txj9"} 1
haproxy_backend_up{backend="http",namespace="mistral",route="mistral-vslnb"} 1
haproxy_backend_up{backend="http",namespace="opensearch",route="opensearch-dashboards-l6vdv"} 1
haproxy_backend_up{backend="http",namespace="profiler",route="esc-collector-service-kdvdw"} 1
haproxy_backend_up{backend="http",namespace="profiler",route="esc-static-service-j9sgv"} 1
haproxy_backend_up{backend="http",namespace="profiler",route="esc-test-service-lbp8r"} 1
haproxy_backend_up{backend="http",namespace="profiler",route="esc-ui-service-r5tfm"} 1
haproxy_backend_up{backend="http",namespace="prometheus-operator",route="prometheus-operator-grafana"} 1
haproxy_backend_up{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmagent-k8s"} 1
haproxy_backend_up{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmalert-k8s"} 1
haproxy_backend_up{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmalertmanager-k8s"} 1
haproxy_backend_up{backend="http",namespace="prometheus-operator",route="prometheus-operator-vmsingle-k8s"} 1
haproxy_backend_up{backend="http",namespace="rabbitmq",route="rabbitmq-fnzhp"} 1
haproxy_backend_up{backend="http",namespace="shared-license-distributor",route="shared-license-distributor-wnj68"} 1
haproxy_backend_up{backend="http",namespace="spark",route="spark-history-server-hqshb"} 1
haproxy_backend_up{backend="http",namespace="spark",route="spark-site-manager-p4nqh"} 1
haproxy_backend_up{backend="https",namespace="openshift-console",route="console"} 1
haproxy_backend_up{backend="https",namespace="openshift-monitoring",route="alertmanager-main"} 1
haproxy_backend_up{backend="https",namespace="openshift-monitoring",route="prometheus-k8s"} 1
haproxy_backend_up{backend="https",namespace="openshift-monitoring",route="prometheus-k8s-federate"} 1
haproxy_backend_up{backend="https",namespace="openshift-monitoring",route="thanos-querier"} 1
haproxy_backend_up{backend="https-edge",namespace="openshift-console",route="downloads"} 1
haproxy_backend_up{backend="https-edge",namespace="openshift-ingress-canary",route="canary"} 1
haproxy_backend_up{backend="other/be_no_sni",namespace="",route=""} 1
haproxy_backend_up{backend="other/be_sni",namespace="",route=""} 1
haproxy_backend_up{backend="other/openshift_default",namespace="",route=""} 1
haproxy_backend_up{backend="tcp",namespace="openshift-authentication",route="oauth-openshift"} 1
haproxy_backend_up{backend="tcp",namespace="site-manager",route="site-manager"} 1
# HELP haproxy_exporter_csv_parse_failures Number of errors while parsing CSV.
# TYPE haproxy_exporter_csv_parse_failures counter
haproxy_exporter_csv_parse_failures 0
# HELP haproxy_exporter_scrape_interval The time in seconds before another scrape is allowed, proportional to size of data.
# TYPE haproxy_exporter_scrape_interval gauge
haproxy_exporter_scrape_interval 5
# HELP haproxy_exporter_server_threshold Number of servers tracked and the current threshold value.
# TYPE haproxy_exporter_server_threshold gauge
haproxy_exporter_server_threshold{type="current"} 50
haproxy_exporter_server_threshold{type="limit"} 500
# HELP haproxy_exporter_total_scrapes Current total HAProxy scrapes.
# TYPE haproxy_exporter_total_scrapes counter
haproxy_exporter_total_scrapes 43519
# HELP haproxy_frontend_bytes_in_total Current total of incoming bytes.
# TYPE haproxy_frontend_bytes_in_total gauge
haproxy_frontend_bytes_in_total{frontend="fe_no_sni"} 0
haproxy_frontend_bytes_in_total{frontend="fe_sni"} 8.5294893e+07
haproxy_frontend_bytes_in_total{frontend="public"} 2.1227866e+07
haproxy_frontend_bytes_in_total{frontend="public_ssl"} 8.8685611e+07
# HELP haproxy_frontend_bytes_out_total Current total of outgoing bytes.
# TYPE haproxy_frontend_bytes_out_total gauge
haproxy_frontend_bytes_out_total{frontend="fe_no_sni"} 0
haproxy_frontend_bytes_out_total{frontend="fe_sni"} 6.373787468e+09
haproxy_frontend_bytes_out_total{frontend="public"} 3.2554298e+07
haproxy_frontend_bytes_out_total{frontend="public_ssl"} 4.737098036e+09
# HELP haproxy_frontend_connections_total Total number of connections.
# TYPE haproxy_frontend_connections_total gauge
haproxy_frontend_connections_total{frontend="fe_no_sni"} 6
haproxy_frontend_connections_total{frontend="fe_sni"} 24941
haproxy_frontend_connections_total{frontend="public"} 173661
haproxy_frontend_connections_total{frontend="public_ssl"} 25021
# HELP haproxy_frontend_current_session_rate Current number of sessions per second over last elapsed second.
# TYPE haproxy_frontend_current_session_rate gauge
haproxy_frontend_current_session_rate{frontend="fe_no_sni"} 0
haproxy_frontend_current_session_rate{frontend="fe_sni"} 0
haproxy_frontend_current_session_rate{frontend="public"} 0
haproxy_frontend_current_session_rate{frontend="public_ssl"} 0
# HELP haproxy_frontend_current_sessions Current number of active sessions.
# TYPE haproxy_frontend_current_sessions gauge
haproxy_frontend_current_sessions{frontend="fe_no_sni"} 0
haproxy_frontend_current_sessions{frontend="fe_sni"} 0
haproxy_frontend_current_sessions{frontend="public"} 0
haproxy_frontend_current_sessions{frontend="public_ssl"} 0
# HELP haproxy_frontend_http_responses_total Total of HTTP responses.
# TYPE haproxy_frontend_http_responses_total gauge
haproxy_frontend_http_responses_total{code="1xx",frontend="fe_no_sni"} 0
haproxy_frontend_http_responses_total{code="1xx",frontend="fe_sni"} 453
haproxy_frontend_http_responses_total{code="1xx",frontend="public"} 0
haproxy_frontend_http_responses_total{code="2xx",frontend="fe_no_sni"} 0
haproxy_frontend_http_responses_total{code="2xx",frontend="fe_sni"} 59914
haproxy_frontend_http_responses_total{code="2xx",frontend="public"} 173558
haproxy_frontend_http_responses_total{code="3xx",frontend="fe_no_sni"} 0
haproxy_frontend_http_responses_total{code="3xx",frontend="fe_sni"} 15
haproxy_frontend_http_responses_total{code="3xx",frontend="public"} 4
haproxy_frontend_http_responses_total{code="4xx",frontend="fe_no_sni"} 0
haproxy_frontend_http_responses_total{code="4xx",frontend="fe_sni"} 2976
haproxy_frontend_http_responses_total{code="4xx",frontend="public"} 0
haproxy_frontend_http_responses_total{code="5xx",frontend="fe_no_sni"} 0
haproxy_frontend_http_responses_total{code="5xx",frontend="fe_sni"} 0
haproxy_frontend_http_responses_total{code="5xx",frontend="public"} 156
haproxy_frontend_http_responses_total{code="other",frontend="fe_no_sni"} 0
haproxy_frontend_http_responses_total{code="other",frontend="fe_sni"} 0
haproxy_frontend_http_responses_total{code="other",frontend="public"} 0
# HELP haproxy_frontend_max_session_rate Maximum observed number of sessions per second.
# TYPE haproxy_frontend_max_session_rate gauge
haproxy_frontend_max_session_rate{frontend="fe_no_sni"} 0
haproxy_frontend_max_session_rate{frontend="fe_sni"} 0
haproxy_frontend_max_session_rate{frontend="public"} 2
haproxy_frontend_max_session_rate{frontend="public_ssl"} 0
# HELP haproxy_frontend_max_sessions Maximum observed number of active sessions.
# TYPE haproxy_frontend_max_sessions gauge
haproxy_frontend_max_sessions{frontend="fe_no_sni"} 0
haproxy_frontend_max_sessions{frontend="fe_sni"} 0
haproxy_frontend_max_sessions{frontend="public"} 3
haproxy_frontend_max_sessions{frontend="public_ssl"} 0
# HELP haproxy_process_cpu_seconds_total Total user and system CPU time spent in seconds.
# TYPE haproxy_process_cpu_seconds_total counter
haproxy_process_cpu_seconds_total 4.83
# HELP haproxy_process_max_fds Maximum number of open file descriptors.
# TYPE haproxy_process_max_fds gauge
haproxy_process_max_fds 1.048576e+06
# HELP haproxy_process_resident_memory_bytes Resident memory size in bytes.
# TYPE haproxy_process_resident_memory_bytes gauge
haproxy_process_resident_memory_bytes 1.7432576e+07
# HELP haproxy_process_start_time_seconds Start time of the process since unix epoch in seconds.
# TYPE haproxy_process_start_time_seconds gauge
haproxy_process_start_time_seconds 1.70350474265e+09
# HELP haproxy_process_virtual_memory_bytes Virtual memory size in bytes.
# TYPE haproxy_process_virtual_memory_bytes gauge
haproxy_process_virtual_memory_bytes 2.8657664e+08
# HELP haproxy_process_virtual_memory_max_bytes Maximum amount of virtual memory available in bytes.
# TYPE haproxy_process_virtual_memory_max_bytes gauge
haproxy_process_virtual_memory_max_bytes 1.8446744073709552e+19
# HELP haproxy_server_bytes_in_total Current total of incoming bytes.
# TYPE haproxy_server_bytes_in_total gauge
haproxy_server_bytes_in_total{namespace="",pod="",route="",server="fe_no_sni",service=""} 2610
haproxy_server_bytes_in_total{namespace="",pod="",route="",server="fe_sni",service=""} 8.846819e+07
haproxy_server_bytes_in_total{namespace="airflow",pod="airflow-site-manager-68ccf87889-bwbdf",route="airflow-site-manager-hpl66",server="10.132.149.80:8080",service="airflow-site-manager"} 0
haproxy_server_bytes_in_total{namespace="airflow",pod="airflow-webserver-695b5c545b-rxqc4",route="airflow-ingress-nzbtm",server="10.134.173.137:8080",service="airflow-webserver"} 0
haproxy_server_bytes_in_total{namespace="arangodb",pod="main-arangodb-1-87669b688-q98z4",route="arangodb-leader-ingress-zw82p",server="10.132.149.185:8529",service="arangodb-leader"} 0
haproxy_server_bytes_in_total{namespace="consul-service",pod="consul-server-0",route="consul-ingress-hmmvf",server="10.134.248.50:8500",service="consul-ui"} 0
haproxy_server_bytes_in_total{namespace="consul-service",pod="consul-server-1",route="consul-ingress-hmmvf",server="10.132.149.160:8500",service="consul-ui"} 0
haproxy_server_bytes_in_total{namespace="consul-service",pod="consul-server-2",route="consul-ingress-hmmvf",server="10.134.173.136:8500",service="consul-ui"} 0
haproxy_server_bytes_in_total{namespace="dbaas",pod="dbaas-aggregator-5f744d74f5-r7tqm",route="aggregator-fw5ft",server="10.134.248.185:8080",service="dbaas-aggregator"} 0
haproxy_server_bytes_in_total{namespace="infra-keycloak",pod="infra-keycloak-7fd57cd59c-d6c5b",route="infra-keycloak-5q4rj",server="10.132.148.223:8080",service="infra-keycloak"} 0
haproxy_server_bytes_in_total{namespace="infra-keycloak",pod="infra-keycloak-7fd57cd59c-qtqv9",route="infra-keycloak-5q4rj",server="10.134.173.182:8080",service="infra-keycloak"} 0
haproxy_server_bytes_in_total{namespace="jaeger",pod="jaeger-hotrod-85b766c49d-s6vz7",route="jaeger-hotrod-qq5cq",server="10.134.249.205:8080",service="jaeger-hotrod"} 0
haproxy_server_bytes_in_total{namespace="jaeger",pod="jaeger-query-65d9446fc6-jx5qk",route="jaeger-query-pc6w5",server="10.134.248.54:16686",service="jaeger-query"} 0
haproxy_server_bytes_in_total{namespace="kafka-service",pod="akhq-58744f76db-h2mtf",route="akhq-ingress-lhz5z",server="10.134.249.133:8080",service="akhq"} 26307
haproxy_server_bytes_in_total{namespace="license-server",pod="license-server-7b8b6f765f-wvcx5",route="license-server-9txj9",server="10.132.148.47:8080",service="license-server"} 0
haproxy_server_bytes_in_total{namespace="opensearch",pod="opensearch-dashboards-764bc8548b-zngcr",route="opensearch-dashboards-l6vdv",server="10.134.173.229:5601",service="opensearch-dashboards"} 0
haproxy_server_bytes_in_total{namespace="openshift-authentication",pod="oauth-openshift-75965c7487-kc8m7",route="oauth-openshift",server="10.132.153.132:6443",service="oauth-openshift"} 0
haproxy_server_bytes_in_total{namespace="openshift-authentication",pod="oauth-openshift-75965c7487-nsfxw",route="oauth-openshift",server="10.133.43.151:6443",service="oauth-openshift"} 214811
haproxy_server_bytes_in_total{namespace="openshift-authentication",pod="oauth-openshift-75965c7487-xmmmq",route="oauth-openshift",server="10.135.207.187:6443",service="oauth-openshift"} 0
haproxy_server_bytes_in_total{namespace="openshift-console",pod="console-66fcc4f4fc-crbz6",route="console",server="10.133.43.162:8443",service="console"} 3.8698283e+07
haproxy_server_bytes_in_total{namespace="openshift-console",pod="console-66fcc4f4fc-j5dvz",route="console",server="10.132.153.185:8443",service="console"} 4.659661e+07
haproxy_server_bytes_in_total{namespace="openshift-console",pod="downloads-6fb844fb47-cwnvw",route="downloads",server="10.135.207.221:8080",service="downloads"} 0
haproxy_server_bytes_in_total{namespace="openshift-console",pod="downloads-6fb844fb47-l485s",route="downloads",server="10.132.153.222:8080",service="downloads"} 0
haproxy_server_bytes_in_total{namespace="openshift-ingress-canary",pod="ingress-canary-7rvzj",route="canary",server="10.134.173.247:8080",service="ingress-canary"} 0
haproxy_server_bytes_in_total{namespace="openshift-ingress-canary",pod="ingress-canary-lwpkn",route="canary",server="10.132.149.137:8080",service="ingress-canary"} 0
haproxy_server_bytes_in_total{namespace="openshift-ingress-canary",pod="ingress-canary-z2b4x",route="canary",server="10.134.248.101:8080",service="ingress-canary"} 0
haproxy_server_bytes_in_total{namespace="openshift-monitoring",pod="alertmanager-main-0",route="alertmanager-main",server="10.134.173.99:9095",service="alertmanager-main"} 0
haproxy_server_bytes_in_total{namespace="openshift-monitoring",pod="alertmanager-main-1",route="alertmanager-main",server="10.134.248.127:9095",service="alertmanager-main"} 0
haproxy_server_bytes_in_total{namespace="openshift-monitoring",pod="prometheus-k8s-0",route="prometheus-k8s",server="10.134.173.51:9091",service="prometheus-k8s"} 0
haproxy_server_bytes_in_total{namespace="openshift-monitoring",pod="prometheus-k8s-0",route="prometheus-k8s-federate",server="10.134.173.51:9091",service="prometheus-k8s"} 0
haproxy_server_bytes_in_total{namespace="openshift-monitoring",pod="prometheus-k8s-1",route="prometheus-k8s",server="10.134.249.153:9091",service="prometheus-k8s"} 0
haproxy_server_bytes_in_total{namespace="openshift-monitoring",pod="prometheus-k8s-1",route="prometheus-k8s-federate",server="10.134.249.153:9091",service="prometheus-k8s"} 0
haproxy_server_bytes_in_total{namespace="openshift-monitoring",pod="thanos-querier-59d7db9787-hwjzt",route="thanos-querier",server="10.134.249.65:9091",service="thanos-querier"} 0
haproxy_server_bytes_in_total{namespace="openshift-monitoring",pod="thanos-querier-59d7db9787-p8m9l",route="thanos-querier",server="10.134.173.112:9091",service="thanos-querier"} 0
haproxy_server_bytes_in_total{namespace="profiler",pod="esc-static-service-654f8bc589-xchjk",route="esc-static-service-j9sgv",server="10.134.173.123:8080",service="esc-static-service"} 0
haproxy_server_bytes_in_total{namespace="profiler",pod="esc-test-service-75849d7c6c-czvsw",route="esc-test-service-lbp8r",server="10.134.172.182:8080",service="esc-test-service"} 0
haproxy_server_bytes_in_total{namespace="profiler",pod="esc-ui-service-66977c9bcb-kg9mc",route="esc-ui-service-r5tfm",server="10.134.249.215:8180",service="esc-ui-service"} 0
haproxy_server_bytes_in_total{namespace="prometheus-operator",pod="grafana-deployment-88d46bc5f-5gx7h",route="prometheus-operator-grafana",server="10.132.149.3:3000",service="grafana-service"} 0
haproxy_server_bytes_in_total{namespace="prometheus-operator",pod="vmagent-k8s-67654d6f97-98nf5",route="prometheus-operator-vmagent-k8s",server="10.132.149.156:8429",service="vmagent-k8s"} 0
haproxy_server_bytes_in_total{namespace="prometheus-operator",pod="vmalert-k8s-6bcf966dd6-jvpz5",route="prometheus-operator-vmalert-k8s",server="10.132.149.175:8080",service="vmalert-k8s"} 0
haproxy_server_bytes_in_total{namespace="prometheus-operator",pod="vmalertmanager-k8s-0",route="prometheus-operator-vmalertmanager-k8s",server="10.134.173.133:9093",service="vmalertmanager-k8s"} 0
haproxy_server_bytes_in_total{namespace="prometheus-operator",pod="vmsingle-k8s-f4844bfc4-5l8sf",route="prometheus-operator-vmsingle-k8s",server="10.134.249.110:8429",service="vmsingle-k8s"} 9319
haproxy_server_bytes_in_total{namespace="rabbitmq",pod="rmqlocal-0",route="rabbitmq-fnzhp",server="10.134.249.141:15672",service="rabbitmq"} 0
haproxy_server_bytes_in_total{namespace="rabbitmq",pod="rmqlocal-1",route="rabbitmq-fnzhp",server="10.134.173.127:15672",service="rabbitmq"} 0
haproxy_server_bytes_in_total{namespace="rabbitmq",pod="rmqlocal-2",route="rabbitmq-fnzhp",server="10.132.149.46:15672",service="rabbitmq"} 0
haproxy_server_bytes_in_total{namespace="shared-license-distributor",pod="shared-license-distributor-64ccc89449-8kw8p",route="shared-license-distributor-wnj68",server="10.134.249.239:8080",service="shared-license-distributor"} 0
haproxy_server_bytes_in_total{namespace="site-manager",pod="site-manager-547546bdd5-hd85n",route="site-manager",server="10.134.248.81:8443",service="site-manager"} 0
haproxy_server_bytes_in_total{namespace="spark",pod="spark-history-server-bd755d9f5-9ccsz",route="spark-history-server-hqshb",server="10.132.149.178:18080",service="spark-history-server"} 0
haproxy_server_bytes_in_total{namespace="spark",pod="spark-history-server-bd755d9f5-zf2dt",route="spark-history-server-hqshb",server="10.134.173.147:18080",service="spark-history-server"} 0
haproxy_server_bytes_in_total{namespace="spark",pod="spark-site-manager-66c44c7b78-cc7jr",route="spark-site-manager-p4nqh",server="10.132.149.148:8080",service="spark-site-manager"} 0
# HELP haproxy_server_bytes_out_total Current total of outgoing bytes.
# TYPE haproxy_server_bytes_out_total gauge
haproxy_server_bytes_out_total{namespace="",pod="",route="",server="fe_no_sni",service=""} 0
haproxy_server_bytes_out_total{namespace="",pod="",route="",server="fe_sni",service=""} 4.734176415e+09
haproxy_server_bytes_out_total{namespace="airflow",pod="airflow-site-manager-68ccf87889-bwbdf",route="airflow-site-manager-hpl66",server="10.132.149.80:8080",service="airflow-site-manager"} 0
haproxy_server_bytes_out_total{namespace="airflow",pod="airflow-webserver-695b5c545b-rxqc4",route="airflow-ingress-nzbtm",server="10.134.173.137:8080",service="airflow-webserver"} 0
haproxy_server_bytes_out_total{namespace="arangodb",pod="main-arangodb-1-87669b688-q98z4",route="arangodb-leader-ingress-zw82p",server="10.132.149.185:8529",service="arangodb-leader"} 0
haproxy_server_bytes_out_total{namespace="consul-service",pod="consul-server-0",route="consul-ingress-hmmvf",server="10.134.248.50:8500",service="consul-ui"} 0
haproxy_server_bytes_out_total{namespace="consul-service",pod="consul-server-1",route="consul-ingress-hmmvf",server="10.132.149.160:8500",service="consul-ui"} 0
haproxy_server_bytes_out_total{namespace="consul-service",pod="consul-server-2",route="consul-ingress-hmmvf",server="10.134.173.136:8500",service="consul-ui"} 0
haproxy_server_bytes_out_total{namespace="dbaas",pod="dbaas-aggregator-5f744d74f5-r7tqm",route="aggregator-fw5ft",server="10.134.248.185:8080",service="dbaas-aggregator"} 0
haproxy_server_bytes_out_total{namespace="infra-keycloak",pod="infra-keycloak-7fd57cd59c-d6c5b",route="infra-keycloak-5q4rj",server="10.132.148.223:8080",service="infra-keycloak"} 0
haproxy_server_bytes_out_total{namespace="infra-keycloak",pod="infra-keycloak-7fd57cd59c-qtqv9",route="infra-keycloak-5q4rj",server="10.134.173.182:8080",service="infra-keycloak"} 0
haproxy_server_bytes_out_total{namespace="jaeger",pod="jaeger-hotrod-85b766c49d-s6vz7",route="jaeger-hotrod-qq5cq",server="10.134.249.205:8080",service="jaeger-hotrod"} 0
haproxy_server_bytes_out_total{namespace="jaeger",pod="jaeger-query-65d9446fc6-jx5qk",route="jaeger-query-pc6w5",server="10.134.248.54:16686",service="jaeger-query"} 0
haproxy_server_bytes_out_total{namespace="kafka-service",pod="akhq-58744f76db-h2mtf",route="akhq-ingress-lhz5z",server="10.134.249.133:8080",service="akhq"} 581056
haproxy_server_bytes_out_total{namespace="license-server",pod="license-server-7b8b6f765f-wvcx5",route="license-server-9txj9",server="10.132.148.47:8080",service="license-server"} 0
haproxy_server_bytes_out_total{namespace="opensearch",pod="opensearch-dashboards-764bc8548b-zngcr",route="opensearch-dashboards-l6vdv",server="10.134.173.229:5601",service="opensearch-dashboards"} 0
haproxy_server_bytes_out_total{namespace="openshift-authentication",pod="oauth-openshift-75965c7487-kc8m7",route="oauth-openshift",server="10.132.153.132:6443",service="oauth-openshift"} 0
haproxy_server_bytes_out_total{namespace="openshift-authentication",pod="oauth-openshift-75965c7487-nsfxw",route="oauth-openshift",server="10.133.43.151:6443",service="oauth-openshift"} 2.921621e+06
haproxy_server_bytes_out_total{namespace="openshift-authentication",pod="oauth-openshift-75965c7487-xmmmq",route="oauth-openshift",server="10.135.207.187:6443",service="oauth-openshift"} 0
haproxy_server_bytes_out_total{namespace="openshift-console",pod="console-66fcc4f4fc-crbz6",route="console",server="10.133.43.162:8443",service="console"} 2.526271709e+09
haproxy_server_bytes_out_total{namespace="openshift-console",pod="console-66fcc4f4fc-j5dvz",route="console",server="10.132.153.185:8443",service="console"} 3.847515759e+09
haproxy_server_bytes_out_total{namespace="openshift-console",pod="downloads-6fb844fb47-cwnvw",route="downloads",server="10.135.207.221:8080",service="downloads"} 0
haproxy_server_bytes_out_total{namespace="openshift-console",pod="downloads-6fb844fb47-l485s",route="downloads",server="10.132.153.222:8080",service="downloads"} 0
haproxy_server_bytes_out_total{namespace="openshift-ingress-canary",pod="ingress-canary-7rvzj",route="canary",server="10.134.173.247:8080",service="ingress-canary"} 0
haproxy_server_bytes_out_total{namespace="openshift-ingress-canary",pod="ingress-canary-lwpkn",route="canary",server="10.132.149.137:8080",service="ingress-canary"} 0
haproxy_server_bytes_out_total{namespace="openshift-ingress-canary",pod="ingress-canary-z2b4x",route="canary",server="10.134.248.101:8080",service="ingress-canary"} 0
haproxy_server_bytes_out_total{namespace="openshift-monitoring",pod="alertmanager-main-0",route="alertmanager-main",server="10.134.173.99:9095",service="alertmanager-main"} 0
haproxy_server_bytes_out_total{namespace="openshift-monitoring",pod="alertmanager-main-1",route="alertmanager-main",server="10.134.248.127:9095",service="alertmanager-main"} 0
haproxy_server_bytes_out_total{namespace="openshift-monitoring",pod="prometheus-k8s-0",route="prometheus-k8s",server="10.134.173.51:9091",service="prometheus-k8s"} 0
haproxy_server_bytes_out_total{namespace="openshift-monitoring",pod="prometheus-k8s-0",route="prometheus-k8s-federate",server="10.134.173.51:9091",service="prometheus-k8s"} 0
haproxy_server_bytes_out_total{namespace="openshift-monitoring",pod="prometheus-k8s-1",route="prometheus-k8s",server="10.134.249.153:9091",service="prometheus-k8s"} 0
haproxy_server_bytes_out_total{namespace="openshift-monitoring",pod="prometheus-k8s-1",route="prometheus-k8s-federate",server="10.134.249.153:9091",service="prometheus-k8s"} 0
haproxy_server_bytes_out_total{namespace="openshift-monitoring",pod="thanos-querier-59d7db9787-hwjzt",route="thanos-querier",server="10.134.249.65:9091",service="thanos-querier"} 0
haproxy_server_bytes_out_total{namespace="openshift-monitoring",pod="thanos-querier-59d7db9787-p8m9l",route="thanos-querier",server="10.134.173.112:9091",service="thanos-querier"} 0
haproxy_server_bytes_out_total{namespace="profiler",pod="esc-static-service-654f8bc589-xchjk",route="esc-static-service-j9sgv",server="10.134.173.123:8080",service="esc-static-service"} 0
haproxy_server_bytes_out_total{namespace="profiler",pod="esc-test-service-75849d7c6c-czvsw",route="esc-test-service-lbp8r",server="10.134.172.182:8080",service="esc-test-service"} 0
haproxy_server_bytes_out_total{namespace="profiler",pod="esc-ui-service-66977c9bcb-kg9mc",route="esc-ui-service-r5tfm",server="10.134.249.215:8180",service="esc-ui-service"} 0
haproxy_server_bytes_out_total{namespace="prometheus-operator",pod="grafana-deployment-88d46bc5f-5gx7h",route="prometheus-operator-grafana",server="10.132.149.3:3000",service="grafana-service"} 0
haproxy_server_bytes_out_total{namespace="prometheus-operator",pod="vmagent-k8s-67654d6f97-98nf5",route="prometheus-operator-vmagent-k8s",server="10.132.149.156:8429",service="vmagent-k8s"} 0
haproxy_server_bytes_out_total{namespace="prometheus-operator",pod="vmalert-k8s-6bcf966dd6-jvpz5",route="prometheus-operator-vmalert-k8s",server="10.132.149.175:8080",service="vmalert-k8s"} 0
haproxy_server_bytes_out_total{namespace="prometheus-operator",pod="vmalertmanager-k8s-0",route="prometheus-operator-vmalertmanager-k8s",server="10.134.173.133:9093",service="vmalertmanager-k8s"} 0
haproxy_server_bytes_out_total{namespace="prometheus-operator",pod="vmsingle-k8s-f4844bfc4-5l8sf",route="prometheus-operator-vmsingle-k8s",server="10.134.249.110:8429",service="vmsingle-k8s"} 851946
haproxy_server_bytes_out_total{namespace="rabbitmq",pod="rmqlocal-0",route="rabbitmq-fnzhp",server="10.134.249.141:15672",service="rabbitmq"} 0
haproxy_server_bytes_out_total{namespace="rabbitmq",pod="rmqlocal-1",route="rabbitmq-fnzhp",server="10.134.173.127:15672",service="rabbitmq"} 0
haproxy_server_bytes_out_total{namespace="rabbitmq",pod="rmqlocal-2",route="rabbitmq-fnzhp",server="10.132.149.46:15672",service="rabbitmq"} 0
haproxy_server_bytes_out_total{namespace="shared-license-distributor",pod="shared-license-distributor-64ccc89449-8kw8p",route="shared-license-distributor-wnj68",server="10.134.249.239:8080",service="shared-license-distributor"} 0
haproxy_server_bytes_out_total{namespace="site-manager",pod="site-manager-547546bdd5-hd85n",route="site-manager",server="10.134.248.81:8443",service="site-manager"} 0
haproxy_server_bytes_out_total{namespace="spark",pod="spark-history-server-bd755d9f5-9ccsz",route="spark-history-server-hqshb",server="10.132.149.178:18080",service="spark-history-server"} 0
haproxy_server_bytes_out_total{namespace="spark",pod="spark-history-server-bd755d9f5-zf2dt",route="spark-history-server-hqshb",server="10.134.173.147:18080",service="spark-history-server"} 0
haproxy_server_bytes_out_total{namespace="spark",pod="spark-site-manager-66c44c7b78-cc7jr",route="spark-site-manager-p4nqh",server="10.132.149.148:8080",service="spark-site-manager"} 0
# HELP haproxy_server_check_failures_total Total number of failed health checks.
# TYPE haproxy_server_check_failures_total gauge
haproxy_server_check_failures_total{namespace="consul-service",pod="consul-server-0",route="consul-ingress-hmmvf",server="10.134.248.50:8500",service="consul-ui"} 0
haproxy_server_check_failures_total{namespace="consul-service",pod="consul-server-1",route="consul-ingress-hmmvf",server="10.132.149.160:8500",service="consul-ui"} 0
haproxy_server_check_failures_total{namespace="consul-service",pod="consul-server-2",route="consul-ingress-hmmvf",server="10.134.173.136:8500",service="consul-ui"} 0
haproxy_server_check_failures_total{namespace="infra-keycloak",pod="infra-keycloak-7fd57cd59c-d6c5b",route="infra-keycloak-5q4rj",server="10.132.148.223:8080",service="infra-keycloak"} 0
haproxy_server_check_failures_total{namespace="infra-keycloak",pod="infra-keycloak-7fd57cd59c-qtqv9",route="infra-keycloak-5q4rj",server="10.134.173.182:8080",service="infra-keycloak"} 0
haproxy_server_check_failures_total{namespace="openshift-authentication",pod="oauth-openshift-75965c7487-kc8m7",route="oauth-openshift",server="10.132.153.132:6443",service="oauth-openshift"} 0
haproxy_server_check_failures_total{namespace="openshift-authentication",pod="oauth-openshift-75965c7487-nsfxw",route="oauth-openshift",server="10.133.43.151:6443",service="oauth-openshift"} 0
haproxy_server_check_failures_total{namespace="openshift-authentication",pod="oauth-openshift-75965c7487-xmmmq",route="oauth-openshift",server="10.135.207.187:6443",service="oauth-openshift"} 0
haproxy_server_check_failures_total{namespace="openshift-console",pod="console-66fcc4f4fc-crbz6",route="console",server="10.133.43.162:8443",service="console"} 0
haproxy_server_check_failures_total{namespace="openshift-console",pod="console-66fcc4f4fc-j5dvz",route="console",server="10.132.153.185:8443",service="console"} 0
haproxy_server_check_failures_total{namespace="openshift-console",pod="downloads-6fb844fb47-cwnvw",route="downloads",server="10.135.207.221:8080",service="downloads"} 0
haproxy_server_check_failures_total{namespace="openshift-console",pod="downloads-6fb844fb47-l485s",route="downloads",server="10.132.153.222:8080",service="downloads"} 0
haproxy_server_check_failures_total{namespace="openshift-ingress-canary",pod="ingress-canary-7rvzj",route="canary",server="10.134.173.247:8080",service="ingress-canary"} 0
haproxy_server_check_failures_total{namespace="openshift-ingress-canary",pod="ingress-canary-lwpkn",route="canary",server="10.132.149.137:8080",service="ingress-canary"} 0
haproxy_server_check_failures_total{namespace="openshift-ingress-canary",pod="ingress-canary-z2b4x",route="canary",server="10.134.248.101:8080",service="ingress-canary"} 0
haproxy_server_check_failures_total{namespace="openshift-monitoring",pod="alertmanager-main-0",route="alertmanager-main",server="10.134.173.99:9095",service="alertmanager-main"} 0
haproxy_server_check_failures_total{namespace="openshift-monitoring",pod="alertmanager-main-1",route="alertmanager-main",server="10.134.248.127:9095",service="alertmanager-main"} 0
haproxy_server_check_failures_total{namespace="openshift-monitoring",pod="prometheus-k8s-0",route="prometheus-k8s",server="10.134.173.51:9091",service="prometheus-k8s"} 0
haproxy_server_check_failures_total{namespace="openshift-monitoring",pod="prometheus-k8s-0",route="prometheus-k8s-federate",server="10.134.173.51:9091",service="prometheus-k8s"} 0
haproxy_server_check_failures_total{namespace="openshift-monitoring",pod="prometheus-k8s-1",route="prometheus-k8s",server="10.134.249.153:9091",service="prometheus-k8s"} 0
haproxy_server_check_failures_total{namespace="openshift-monitoring",pod="prometheus-k8s-1",route="prometheus-k8s-federate",server="10.134.249.153:9091",service="prometheus-k8s"} 0
haproxy_server_check_failures_total{namespace="openshift-monitoring",pod="thanos-querier-59d7db9787-hwjzt",route="thanos-querier",server="10.134.249.65:9091",service="thanos-querier"} 0
haproxy_server_check_failures_total{namespace="openshift-monitoring",pod="thanos-querier-59d7db9787-p8m9l",route="thanos-querier",server="10.134.173.112:9091",service="thanos-querier"} 0
haproxy_server_check_failures_total{namespace="rabbitmq",pod="rmqlocal-0",route="rabbitmq-fnzhp",server="10.134.249.141:15672",service="rabbitmq"} 6
haproxy_server_check_failures_total{namespace="rabbitmq",pod="rmqlocal-1",route="rabbitmq-fnzhp",server="10.134.173.127:15672",service="rabbitmq"} 6
haproxy_server_check_failures_total{namespace="rabbitmq",pod="rmqlocal-2",route="rabbitmq-fnzhp",server="10.132.149.46:15672",service="rabbitmq"} 2
haproxy_server_check_failures_total{namespace="spark",pod="spark-history-server-bd755d9f5-9ccsz",route="spark-history-server-hqshb",server="10.132.149.178:18080",service="spark-history-server"} 1
haproxy_server_check_failures_total{namespace="spark",pod="spark-history-server-bd755d9f5-zf2dt",route="spark-history-server-hqshb",server="10.134.173.147:18080",service="spark-history-server"} 0
# HELP haproxy_server_connection_errors_total Total of connection errors.
# TYPE haproxy_server_connection_errors_total gauge
haproxy_server_connection_errors_total{namespace="",pod="",route="",server="fe_no_sni",service=""} 0
haproxy_server_connection_errors_total{namespace="",pod="",route="",server="fe_sni",service=""} 0
haproxy_server_connection_errors_total{namespace="airflow",pod="airflow-site-manager-68ccf87889-bwbdf",route="airflow-site-manager-hpl66",server="10.132.149.80:8080",service="airflow-site-manager"} 0
haproxy_server_connection_errors_total{namespace="airflow",pod="airflow-webserver-695b5c545b-rxqc4",route="airflow-ingress-nzbtm",server="10.134.173.137:8080",service="airflow-webserver"} 0
haproxy_server_connection_errors_total{namespace="arangodb",pod="main-arangodb-1-87669b688-q98z4",route="arangodb-leader-ingress-zw82p",server="10.132.149.185:8529",service="arangodb-leader"} 0
haproxy_server_connection_errors_total{namespace="consul-service",pod="consul-server-0",route="consul-ingress-hmmvf",server="10.134.248.50:8500",service="consul-ui"} 0
haproxy_server_connection_errors_total{namespace="consul-service",pod="consul-server-1",route="consul-ingress-hmmvf",server="10.132.149.160:8500",service="consul-ui"} 0
haproxy_server_connection_errors_total{namespace="consul-service",pod="consul-server-2",route="consul-ingress-hmmvf",server="10.134.173.136:8500",service="consul-ui"} 0
haproxy_server_connection_errors_total{namespace="dbaas",pod="dbaas-aggregator-5f744d74f5-r7tqm",route="aggregator-fw5ft",server="10.134.248.185:8080",service="dbaas-aggregator"} 0
haproxy_server_connection_errors_total{namespace="infra-keycloak",pod="infra-keycloak-7fd57cd59c-d6c5b",route="infra-keycloak-5q4rj",server="10.132.148.223:8080",service="infra-keycloak"} 0
haproxy_server_connection_errors_total{namespace="infra-keycloak",pod="infra-keycloak-7fd57cd59c-qtqv9",route="infra-keycloak-5q4rj",server="10.134.173.182:8080",service="infra-keycloak"} 0
haproxy_server_connection_errors_total{namespace="jaeger",pod="jaeger-hotrod-85b766c49d-s6vz7",route="jaeger-hotrod-qq5cq",server="10.134.249.205:8080",service="jaeger-hotrod"} 0
haproxy_server_connection_errors_total{namespace="jaeger",pod="jaeger-query-65d9446fc6-jx5qk",route="jaeger-query-pc6w5",server="10.134.248.54:16686",service="jaeger-query"} 0
haproxy_server_connection_errors_total{namespace="kafka-service",pod="akhq-58744f76db-h2mtf",route="akhq-ingress-lhz5z",server="10.134.249.133:8080",service="akhq"} 0
haproxy_server_connection_errors_total{namespace="license-server",pod="license-server-7b8b6f765f-wvcx5",route="license-server-9txj9",server="10.132.148.47:8080",service="license-server"} 0
haproxy_server_connection_errors_total{namespace="opensearch",pod="opensearch-dashboards-764bc8548b-zngcr",route="opensearch-dashboards-l6vdv",server="10.134.173.229:5601",service="opensearch-dashboards"} 0
haproxy_server_connection_errors_total{namespace="openshift-authentication",pod="oauth-openshift-75965c7487-kc8m7",route="oauth-openshift",server="10.132.153.132:6443",service="oauth-openshift"} 0
haproxy_server_connection_errors_total{namespace="openshift-authentication",pod="oauth-openshift-75965c7487-nsfxw",route="oauth-openshift",server="10.133.43.151:6443",service="oauth-openshift"} 0
haproxy_server_connection_errors_total{namespace="openshift-authentication",pod="oauth-openshift-75965c7487-xmmmq",route="oauth-openshift",server="10.135.207.187:6443",service="oauth-openshift"} 0
haproxy_server_connection_errors_total{namespace="openshift-console",pod="console-66fcc4f4fc-crbz6",route="console",server="10.133.43.162:8443",service="console"} 0
haproxy_server_connection_errors_total{namespace="openshift-console",pod="console-66fcc4f4fc-j5dvz",route="console",server="10.132.153.185:8443",service="console"} 0
haproxy_server_connection_errors_total{namespace="openshift-console",pod="downloads-6fb844fb47-cwnvw",route="downloads",server="10.135.207.221:8080",service="downloads"} 0
haproxy_server_connection_errors_total{namespace="openshift-console",pod="downloads-6fb844fb47-l485s",route="downloads",server="10.132.153.222:8080",service="downloads"} 0
haproxy_server_connection_errors_total{namespace="openshift-ingress-canary",pod="ingress-canary-7rvzj",route="canary",server="10.134.173.247:8080",service="ingress-canary"} 0
haproxy_server_connection_errors_total{namespace="openshift-ingress-canary",pod="ingress-canary-lwpkn",route="canary",server="10.132.149.137:8080",service="ingress-canary"} 0
haproxy_server_connection_errors_total{namespace="openshift-ingress-canary",pod="ingress-canary-z2b4x",route="canary",server="10.134.248.101:8080",service="ingress-canary"} 0
haproxy_server_connection_errors_total{namespace="openshift-monitoring",pod="alertmanager-main-0",route="alertmanager-main",server="10.134.173.99:9095",service="alertmanager-main"} 0
haproxy_server_connection_errors_total{namespace="openshift-monitoring",pod="alertmanager-main-1",route="alertmanager-main",server="10.134.248.127:9095",service="alertmanager-main"} 0
haproxy_server_connection_errors_total{namespace="openshift-monitoring",pod="prometheus-k8s-0",route="prometheus-k8s",server="10.134.173.51:9091",service="prometheus-k8s"} 0
haproxy_server_connection_errors_total{namespace="openshift-monitoring",pod="prometheus-k8s-0",route="prometheus-k8s-federate",server="10.134.173.51:9091",service="prometheus-k8s"} 0
haproxy_server_connection_errors_total{namespace="openshift-monitoring",pod="prometheus-k8s-1",route="prometheus-k8s",server="10.134.249.153:9091",service="prometheus-k8s"} 0
haproxy_server_connection_errors_total{namespace="openshift-monitoring",pod="prometheus-k8s-1",route="prometheus-k8s-federate",server="10.134.249.153:9091",service="prometheus-k8s"} 0
haproxy_server_connection_errors_total{namespace="openshift-monitoring",pod="thanos-querier-59d7db9787-hwjzt",route="thanos-querier",server="10.134.249.65:9091",service="thanos-querier"} 0
haproxy_server_connection_errors_total{namespace="openshift-monitoring",pod="thanos-querier-59d7db9787-p8m9l",route="thanos-querier",server="10.134.173.112:9091",service="thanos-querier"} 0
haproxy_server_connection_errors_total{namespace="profiler",pod="esc-static-service-654f8bc589-xchjk",route="esc-static-service-j9sgv",server="10.134.173.123:8080",service="esc-static-service"} 0
haproxy_server_connection_errors_total{namespace="profiler",pod="esc-test-service-75849d7c6c-czvsw",route="esc-test-service-lbp8r",server="10.134.172.182:8080",service="esc-test-service"} 0
haproxy_server_connection_errors_total{namespace="profiler",pod="esc-ui-service-66977c9bcb-kg9mc",route="esc-ui-service-r5tfm",server="10.134.249.215:8180",service="esc-ui-service"} 0
haproxy_server_connection_errors_total{namespace="prometheus-operator",pod="grafana-deployment-88d46bc5f-5gx7h",route="prometheus-operator-grafana",server="10.132.149.3:3000",service="grafana-service"} 0
haproxy_server_connection_errors_total{namespace="prometheus-operator",pod="vmagent-k8s-67654d6f97-98nf5",route="prometheus-operator-vmagent-k8s",server="10.132.149.156:8429",service="vmagent-k8s"} 0
haproxy_server_connection_errors_total{namespace="prometheus-operator",pod="vmalert-k8s-6bcf966dd6-jvpz5",route="prometheus-operator-vmalert-k8s",server="10.132.149.175:8080",service="vmalert-k8s"} 0
haproxy_server_connection_errors_total{namespace="prometheus-operator",pod="vmalertmanager-k8s-0",route="prometheus-operator-vmalertmanager-k8s",server="10.134.173.133:9093",service="vmalertmanager-k8s"} 0
haproxy_server_connection_errors_total{namespace="prometheus-operator",pod="vmsingle-k8s-f4844bfc4-5l8sf",route="prometheus-operator-vmsingle-k8s",server="10.134.249.110:8429",service="vmsingle-k8s"} 0
haproxy_server_connection_errors_total{namespace="rabbitmq",pod="rmqlocal-0",route="rabbitmq-fnzhp",server="10.134.249.141:15672",service="rabbitmq"} 0
haproxy_server_connection_errors_total{namespace="rabbitmq",pod="rmqlocal-1",route="rabbitmq-fnzhp",server="10.134.173.127:15672",service="rabbitmq"} 0
haproxy_server_connection_errors_total{namespace="rabbitmq",pod="rmqlocal-2",route="rabbitmq-fnzhp",server="10.132.149.46:15672",service="rabbitmq"} 0
haproxy_server_connection_errors_total{namespace="shared-license-distributor",pod="shared-license-distributor-64ccc89449-8kw8p",route="shared-license-distributor-wnj68",server="10.134.249.239:8080",service="shared-license-distributor"} 0
haproxy_server_connection_errors_total{namespace="site-manager",pod="site-manager-547546bdd5-hd85n",route="site-manager",server="10.134.248.81:8443",service="site-manager"} 0
haproxy_server_connection_errors_total{namespace="spark",pod="spark-history-server-bd755d9f5-9ccsz",route="spark-history-server-hqshb",server="10.132.149.178:18080",service="spark-history-server"} 0
haproxy_server_connection_errors_total{namespace="spark",pod="spark-history-server-bd755d9f5-zf2dt",route="spark-history-server-hqshb",server="10.134.173.147:18080",service="spark-history-server"} 0
haproxy_server_connection_errors_total{namespace="spark",pod="spark-site-manager-66c44c7b78-cc7jr",route="spark-site-manager-p4nqh",server="10.132.149.148:8080",service="spark-site-manager"} 0
# HELP haproxy_server_connections_reused_total Total number of connections reused.
# TYPE haproxy_server_connections_reused_total gauge
haproxy_server_connections_reused_total{namespace="",pod="",route="",server="fe_no_sni",service=""} 0
haproxy_server_connections_reused_total{namespace="",pod="",route="",server="fe_sni",service=""} 0
haproxy_server_connections_reused_total{namespace="airflow",pod="airflow-site-manager-68ccf87889-bwbdf",route="airflow-site-manager-hpl66",server="10.132.149.80:8080",service="airflow-site-manager"} 0
haproxy_server_connections_reused_total{namespace="airflow",pod="airflow-webserver-695b5c545b-rxqc4",route="airflow-ingress-nzbtm",server="10.134.173.137:8080",service="airflow-webserver"} 0
haproxy_server_connections_reused_total{namespace="arangodb",pod="main-arangodb-1-87669b688-q98z4",route="arangodb-leader-ingress-zw82p",server="10.132.149.185:8529",service="arangodb-leader"} 0
haproxy_server_connections_reused_total{namespace="consul-service",pod="consul-server-0",route="consul-ingress-hmmvf",server="10.134.248.50:8500",service="consul-ui"} 0
haproxy_server_connections_reused_total{namespace="consul-service",pod="consul-server-1",route="consul-ingress-hmmvf",server="10.132.149.160:8500",service="consul-ui"} 0
haproxy_server_connections_reused_total{namespace="consul-service",pod="consul-server-2",route="consul-ingress-hmmvf",server="10.134.173.136:8500",service="consul-ui"} 0
haproxy_server_connections_reused_total{namespace="dbaas",pod="dbaas-aggregator-5f744d74f5-r7tqm",route="aggregator-fw5ft",server="10.134.248.185:8080",service="dbaas-aggregator"} 0
haproxy_server_connections_reused_total{namespace="infra-keycloak",pod="infra-keycloak-7fd57cd59c-d6c5b",route="infra-keycloak-5q4rj",server="10.132.148.223:8080",service="infra-keycloak"} 0
haproxy_server_connections_reused_total{namespace="infra-keycloak",pod="infra-keycloak-7fd57cd59c-qtqv9",route="infra-keycloak-5q4rj",server="10.134.173.182:8080",service="infra-keycloak"} 0
haproxy_server_connections_reused_total{namespace="jaeger",pod="jaeger-hotrod-85b766c49d-s6vz7",route="jaeger-hotrod-qq5cq",server="10.134.249.205:8080",service="jaeger-hotrod"} 0
haproxy_server_connections_reused_total{namespace="jaeger",pod="jaeger-query-65d9446fc6-jx5qk",route="jaeger-query-pc6w5",server="10.134.248.54:16686",service="jaeger-query"} 0
haproxy_server_connections_reused_total{namespace="kafka-service",pod="akhq-58744f76db-h2mtf",route="akhq-ingress-lhz5z",server="10.134.249.133:8080",service="akhq"} 19
haproxy_server_connections_reused_total{namespace="license-server",pod="license-server-7b8b6f765f-wvcx5",route="license-server-9txj9",server="10.132.148.47:8080",service="license-server"} 0
haproxy_server_connections_reused_total{namespace="opensearch",pod="opensearch-dashboards-764bc8548b-zngcr",route="opensearch-dashboards-l6vdv",server="10.134.173.229:5601",service="opensearch-dashboards"} 0
haproxy_server_connections_reused_total{namespace="openshift-authentication",pod="oauth-openshift-75965c7487-kc8m7",route="oauth-openshift",server="10.132.153.132:6443",service="oauth-openshift"} 0
haproxy_server_connections_reused_total{namespace="openshift-authentication",pod="oauth-openshift-75965c7487-nsfxw",route="oauth-openshift",server="10.133.43.151:6443",service="oauth-openshift"} 0
haproxy_server_connections_reused_total{namespace="openshift-authentication",pod="oauth-openshift-75965c7487-xmmmq",route="oauth-openshift",server="10.135.207.187:6443",service="oauth-openshift"} 0
haproxy_server_connections_reused_total{namespace="openshift-console",pod="console-66fcc4f4fc-crbz6",route="console",server="10.133.43.162:8443",service="console"} 23599
haproxy_server_connections_reused_total{namespace="openshift-console",pod="console-66fcc4f4fc-j5dvz",route="console",server="10.132.153.185:8443",service="console"} 23303
haproxy_server_connections_reused_total{namespace="openshift-console",pod="downloads-6fb844fb47-cwnvw",route="downloads",server="10.135.207.221:8080",service="downloads"} 0
haproxy_server_connections_reused_total{namespace="openshift-console",pod="downloads-6fb844fb47-l485s",route="downloads",server="10.132.153.222:8080",service="downloads"} 0
haproxy_server_connections_reused_total{namespace="openshift-ingress-canary",pod="ingress-canary-7rvzj",route="canary",server="10.134.173.247:8080",service="ingress-canary"} 0
haproxy_server_connections_reused_total{namespace="openshift-ingress-canary",pod="ingress-canary-lwpkn",route="canary",server="10.132.149.137:8080",service="ingress-canary"} 0
haproxy_server_connections_reused_total{namespace="openshift-ingress-canary",pod="ingress-canary-z2b4x",route="canary",server="10.134.248.101:8080",service="ingress-canary"} 0
haproxy_server_connections_reused_total{namespace="openshift-monitoring",pod="alertmanager-main-0",route="alertmanager-main",server="10.134.173.99:9095",service="alertmanager-main"} 0
haproxy_server_connections_reused_total{namespace="openshift-monitoring",pod="alertmanager-main-1",route="alertmanager-main",server="10.134.248.127:9095",service="alertmanager-main"} 0
haproxy_server_connections_reused_total{namespace="openshift-monitoring",pod="prometheus-k8s-0",route="prometheus-k8s",server="10.134.173.51:9091",service="prometheus-k8s"} 0
haproxy_server_connections_reused_total{namespace="openshift-monitoring",pod="prometheus-k8s-0",route="prometheus-k8s-federate",server="10.134.173.51:9091",service="prometheus-k8s"} 0
haproxy_server_connections_reused_total{namespace="openshift-monitoring",pod="prometheus-k8s-1",route="prometheus-k8s",server="10.134.249.153:9091",service="prometheus-k8s"} 0
haproxy_server_connections_reused_total{namespace="openshift-monitoring",pod="prometheus-k8s-1",route="prometheus-k8s-federate",server="10.134.249.153:9091",service="prometheus-k8s"} 0
haproxy_server_connections_reused_total{namespace="openshift-monitoring",pod="thanos-querier-59d7db9787-hwjzt",route="thanos-querier",server="10.134.249.65:9091",service="thanos-querier"} 0
haproxy_server_connections_reused_total{namespace="openshift-monitoring",pod="thanos-querier-59d7db9787-p8m9l",route="thanos-querier",server="10.134.173.112:9091",service="thanos-querier"} 0
haproxy_server_connections_reused_total{namespace="profiler",pod="esc-static-service-654f8bc589-xchjk",route="esc-static-service-j9sgv",server="10.134.173.123:8080",service="esc-static-service"} 0
haproxy_server_connections_reused_total{namespace="profiler",pod="esc-test-service-75849d7c6c-czvsw",route="esc-test-service-lbp8r",server="10.134.172.182:8080",service="esc-test-service"} 0
haproxy_server_connections_reused_total{namespace="profiler",pod="esc-ui-service-66977c9bcb-kg9mc",route="esc-ui-service-r5tfm",server="10.134.249.215:8180",service="esc-ui-service"} 0
haproxy_server_connections_reused_total{namespace="prometheus-operator",pod="grafana-deployment-88d46bc5f-5gx7h",route="prometheus-operator-grafana",server="10.132.149.3:3000",service="grafana-service"} 0
haproxy_server_connections_reused_total{namespace="prometheus-operator",pod="vmagent-k8s-67654d6f97-98nf5",route="prometheus-operator-vmagent-k8s",server="10.132.149.156:8429",service="vmagent-k8s"} 0
haproxy_server_connections_reused_total{namespace="prometheus-operator",pod="vmalert-k8s-6bcf966dd6-jvpz5",route="prometheus-operator-vmalert-k8s",server="10.132.149.175:8080",service="vmalert-k8s"} 0
haproxy_server_connections_reused_total{namespace="prometheus-operator",pod="vmalertmanager-k8s-0",route="prometheus-operator-vmalertmanager-k8s",server="10.134.173.133:9093",service="vmalertmanager-k8s"} 0
haproxy_server_connections_reused_total{namespace="prometheus-operator",pod="vmsingle-k8s-f4844bfc4-5l8sf",route="prometheus-operator-vmsingle-k8s",server="10.134.249.110:8429",service="vmsingle-k8s"} 27
haproxy_server_connections_reused_total{namespace="rabbitmq",pod="rmqlocal-0",route="rabbitmq-fnzhp",server="10.134.249.141:15672",service="rabbitmq"} 0
haproxy_server_connections_reused_total{namespace="rabbitmq",pod="rmqlocal-1",route="rabbitmq-fnzhp",server="10.134.173.127:15672",service="rabbitmq"} 0
haproxy_server_connections_reused_total{namespace="rabbitmq",pod="rmqlocal-2",route="rabbitmq-fnzhp",server="10.132.149.46:15672",service="rabbitmq"} 0
haproxy_server_connections_reused_total{namespace="shared-license-distributor",pod="shared-license-distributor-64ccc89449-8kw8p",route="shared-license-distributor-wnj68",server="10.134.249.239:8080",service="shared-license-distributor"} 0
haproxy_server_connections_reused_total{namespace="site-manager",pod="site-manager-547546bdd5-hd85n",route="site-manager",server="10.134.248.81:8443",service="site-manager"} 0
haproxy_server_connections_reused_total{namespace="spark",pod="spark-history-server-bd755d9f5-9ccsz",route="spark-history-server-hqshb",server="10.132.149.178:18080",service="spark-history-server"} 0
haproxy_server_connections_reused_total{namespace="spark",pod="spark-history-server-bd755d9f5-zf2dt",route="spark-history-server-hqshb",server="10.134.173.147:18080",service="spark-history-server"} 0
haproxy_server_connections_reused_total{namespace="spark",pod="spark-site-manager-66c44c7b78-cc7jr",route="spark-site-manager-p4nqh",server="10.132.149.148:8080",service="spark-site-manager"} 0
# HELP haproxy_server_connections_total Total number of connections.
# TYPE haproxy_server_connections_total gauge
haproxy_server_connections_total{namespace="",pod="",route="",server="fe_no_sni",service=""} 6
haproxy_server_connections_total{namespace="",pod="",route="",server="fe_sni",service=""} 24941
haproxy_server_connections_total{namespace="airflow",pod="airflow-site-manager-68ccf87889-bwbdf",route="airflow-site-manager-hpl66",server="10.132.149.80:8080",service="airflow-site-manager"} 0
haproxy_server_connections_total{namespace="airflow",pod="airflow-webserver-695b5c545b-rxqc4",route="airflow-ingress-nzbtm",server="10.134.173.137:8080",service="airflow-webserver"} 0
haproxy_server_connections_total{namespace="arangodb",pod="main-arangodb-1-87669b688-q98z4",route="arangodb-leader-ingress-zw82p",server="10.132.149.185:8529",service="arangodb-leader"} 0
haproxy_server_connections_total{namespace="consul-service",pod="consul-server-0",route="consul-ingress-hmmvf",server="10.134.248.50:8500",service="consul-ui"} 0
haproxy_server_connections_total{namespace="consul-service",pod="consul-server-1",route="consul-ingress-hmmvf",server="10.132.149.160:8500",service="consul-ui"} 0
haproxy_server_connections_total{namespace="consul-service",pod="consul-server-2",route="consul-ingress-hmmvf",server="10.134.173.136:8500",service="consul-ui"} 0
haproxy_server_connections_total{namespace="dbaas",pod="dbaas-aggregator-5f744d74f5-r7tqm",route="aggregator-fw5ft",server="10.134.248.185:8080",service="dbaas-aggregator"} 0
haproxy_server_connections_total{namespace="infra-keycloak",pod="infra-keycloak-7fd57cd59c-d6c5b",route="infra-keycloak-5q4rj",server="10.132.148.223:8080",service="infra-keycloak"} 0
haproxy_server_connections_total{namespace="infra-keycloak",pod="infra-keycloak-7fd57cd59c-qtqv9",route="infra-keycloak-5q4rj",server="10.134.173.182:8080",service="infra-keycloak"} 0
haproxy_server_connections_total{namespace="jaeger",pod="jaeger-hotrod-85b766c49d-s6vz7",route="jaeger-hotrod-qq5cq",server="10.134.249.205:8080",service="jaeger-hotrod"} 0
haproxy_server_connections_total{namespace="jaeger",pod="jaeger-query-65d9446fc6-jx5qk",route="jaeger-query-pc6w5",server="10.134.248.54:16686",service="jaeger-query"} 0
haproxy_server_connections_total{namespace="kafka-service",pod="akhq-58744f76db-h2mtf",route="akhq-ingress-lhz5z",server="10.134.249.133:8080",service="akhq"} 22
haproxy_server_connections_total{namespace="license-server",pod="license-server-7b8b6f765f-wvcx5",route="license-server-9txj9",server="10.132.148.47:8080",service="license-server"} 0
haproxy_server_connections_total{namespace="opensearch",pod="opensearch-dashboards-764bc8548b-zngcr",route="opensearch-dashboards-l6vdv",server="10.134.173.229:5601",service="opensearch-dashboards"} 0
haproxy_server_connections_total{namespace="openshift-authentication",pod="oauth-openshift-75965c7487-kc8m7",route="oauth-openshift",server="10.132.153.132:6443",service="oauth-openshift"} 0
haproxy_server_connections_total{namespace="openshift-authentication",pod="oauth-openshift-75965c7487-nsfxw",route="oauth-openshift",server="10.133.43.151:6443",service="oauth-openshift"} 74
haproxy_server_connections_total{namespace="openshift-authentication",pod="oauth-openshift-75965c7487-xmmmq",route="oauth-openshift",server="10.135.207.187:6443",service="oauth-openshift"} 0
haproxy_server_connections_total{namespace="openshift-console",pod="console-66fcc4f4fc-crbz6",route="console",server="10.133.43.162:8443",service="console"} 27590
haproxy_server_connections_total{namespace="openshift-console",pod="console-66fcc4f4fc-j5dvz",route="console",server="10.132.153.185:8443",service="console"} 33359
haproxy_server_connections_total{namespace="openshift-console",pod="downloads-6fb844fb47-cwnvw",route="downloads",server="10.135.207.221:8080",service="downloads"} 0
haproxy_server_connections_total{namespace="openshift-console",pod="downloads-6fb844fb47-l485s",route="downloads",server="10.132.153.222:8080",service="downloads"} 0
haproxy_server_connections_total{namespace="openshift-ingress-canary",pod="ingress-canary-7rvzj",route="canary",server="10.134.173.247:8080",service="ingress-canary"} 0
haproxy_server_connections_total{namespace="openshift-ingress-canary",pod="ingress-canary-lwpkn",route="canary",server="10.132.149.137:8080",service="ingress-canary"} 0
haproxy_server_connections_total{namespace="openshift-ingress-canary",pod="ingress-canary-z2b4x",route="canary",server="10.134.248.101:8080",service="ingress-canary"} 0
haproxy_server_connections_total{namespace="openshift-monitoring",pod="alertmanager-main-0",route="alertmanager-main",server="10.134.173.99:9095",service="alertmanager-main"} 0
haproxy_server_connections_total{namespace="openshift-monitoring",pod="alertmanager-main-1",route="alertmanager-main",server="10.134.248.127:9095",service="alertmanager-main"} 0
haproxy_server_connections_total{namespace="openshift-monitoring",pod="prometheus-k8s-0",route="prometheus-k8s",server="10.134.173.51:9091",service="prometheus-k8s"} 0
haproxy_server_connections_total{namespace="openshift-monitoring",pod="prometheus-k8s-0",route="prometheus-k8s-federate",server="10.134.173.51:9091",service="prometheus-k8s"} 0
haproxy_server_connections_total{namespace="openshift-monitoring",pod="prometheus-k8s-1",route="prometheus-k8s",server="10.134.249.153:9091",service="prometheus-k8s"} 0
haproxy_server_connections_total{namespace="openshift-monitoring",pod="prometheus-k8s-1",route="prometheus-k8s-federate",server="10.134.249.153:9091",service="prometheus-k8s"} 0
haproxy_server_connections_total{namespace="openshift-monitoring",pod="thanos-querier-59d7db9787-hwjzt",route="thanos-querier",server="10.134.249.65:9091",service="thanos-querier"} 0
haproxy_server_connections_total{namespace="openshift-monitoring",pod="thanos-querier-59d7db9787-p8m9l",route="thanos-querier",server="10.134.173.112:9091",service="thanos-querier"} 0
haproxy_server_connections_total{namespace="profiler",pod="esc-static-service-654f8bc589-xchjk",route="esc-static-service-j9sgv",server="10.134.173.123:8080",service="esc-static-service"} 0
haproxy_server_connections_total{namespace="profiler",pod="esc-test-service-75849d7c6c-czvsw",route="esc-test-service-lbp8r",server="10.134.172.182:8080",service="esc-test-service"} 0
haproxy_server_connections_total{namespace="profiler",pod="esc-ui-service-66977c9bcb-kg9mc",route="esc-ui-service-r5tfm",server="10.134.249.215:8180",service="esc-ui-service"} 0
haproxy_server_connections_total{namespace="prometheus-operator",pod="grafana-deployment-88d46bc5f-5gx7h",route="prometheus-operator-grafana",server="10.132.149.3:3000",service="grafana-service"} 0
haproxy_server_connections_total{namespace="prometheus-operator",pod="vmagent-k8s-67654d6f97-98nf5",route="prometheus-operator-vmagent-k8s",server="10.132.149.156:8429",service="vmagent-k8s"} 0
haproxy_server_connections_total{namespace="prometheus-operator",pod="vmalert-k8s-6bcf966dd6-jvpz5",route="prometheus-operator-vmalert-k8s",server="10.132.149.175:8080",service="vmalert-k8s"} 0
haproxy_server_connections_total{namespace="prometheus-operator",pod="vmalertmanager-k8s-0",route="prometheus-operator-vmalertmanager-k8s",server="10.134.173.133:9093",service="vmalertmanager-k8s"} 0
haproxy_server_connections_total{namespace="prometheus-operator",pod="vmsingle-k8s-f4844bfc4-5l8sf",route="prometheus-operator-vmsingle-k8s",server="10.134.249.110:8429",service="vmsingle-k8s"} 36
haproxy_server_connections_total{namespace="rabbitmq",pod="rmqlocal-0",route="rabbitmq-fnzhp",server="10.134.249.141:15672",service="rabbitmq"} 0
haproxy_server_connections_total{namespace="rabbitmq",pod="rmqlocal-1",route="rabbitmq-fnzhp",server="10.134.173.127:15672",service="rabbitmq"} 0
haproxy_server_connections_total{namespace="rabbitmq",pod="rmqlocal-2",route="rabbitmq-fnzhp",server="10.132.149.46:15672",service="rabbitmq"} 0
haproxy_server_connections_total{namespace="shared-license-distributor",pod="shared-license-distributor-64ccc89449-8kw8p",route="shared-license-distributor-wnj68",server="10.134.249.239:8080",service="shared-license-distributor"} 0
haproxy_server_connections_total{namespace="site-manager",pod="site-manager-547546bdd5-hd85n",route="site-manager",server="10.134.248.81:8443",service="site-manager"} 0
haproxy_server_connections_total{namespace="spark",pod="spark-history-server-bd755d9f5-9ccsz",route="spark-history-server-hqshb",server="10.132.149.178:18080",service="spark-history-server"} 0
haproxy_server_connections_total{namespace="spark",pod="spark-history-server-bd755d9f5-zf2dt",route="spark-history-server-hqshb",server="10.134.173.147:18080",service="spark-history-server"} 0
haproxy_server_connections_total{namespace="spark",pod="spark-site-manager-66c44c7b78-cc7jr",route="spark-site-manager-p4nqh",server="10.132.149.148:8080",service="spark-site-manager"} 0
# HELP haproxy_server_current_queue Current number of queued requests assigned to this server.
# TYPE haproxy_server_current_queue gauge
haproxy_server_current_queue{namespace="",pod="",route="",server="fe_no_sni",service=""} 0
haproxy_server_current_queue{namespace="",pod="",route="",server="fe_sni",service=""} 0
haproxy_server_current_queue{namespace="airflow",pod="airflow-site-manager-68ccf87889-bwbdf",route="airflow-site-manager-hpl66",server="10.132.149.80:8080",service="airflow-site-manager"} 0
haproxy_server_current_queue{namespace="airflow",pod="airflow-webserver-695b5c545b-rxqc4",route="airflow-ingress-nzbtm",server="10.134.173.137:8080",service="airflow-webserver"} 0
haproxy_server_current_queue{namespace="arangodb",pod="main-arangodb-1-87669b688-q98z4",route="arangodb-leader-ingress-zw82p",server="10.132.149.185:8529",service="arangodb-leader"} 0
haproxy_server_current_queue{namespace="consul-service",pod="consul-server-0",route="consul-ingress-hmmvf",server="10.134.248.50:8500",service="consul-ui"} 0
haproxy_server_current_queue{namespace="consul-service",pod="consul-server-1",route="consul-ingress-hmmvf",server="10.132.149.160:8500",service="consul-ui"} 0
haproxy_server_current_queue{namespace="consul-service",pod="consul-server-2",route="consul-ingress-hmmvf",server="10.134.173.136:8500",service="consul-ui"} 0
haproxy_server_current_queue{namespace="dbaas",pod="dbaas-aggregator-5f744d74f5-r7tqm",route="aggregator-fw5ft",server="10.134.248.185:8080",service="dbaas-aggregator"} 0
haproxy_server_current_queue{namespace="infra-keycloak",pod="infra-keycloak-7fd57cd59c-d6c5b",route="infra-keycloak-5q4rj",server="10.132.148.223:8080",service="infra-keycloak"} 0
haproxy_server_current_queue{namespace="infra-keycloak",pod="infra-keycloak-7fd57cd59c-qtqv9",route="infra-keycloak-5q4rj",server="10.134.173.182:8080",service="infra-keycloak"} 0
haproxy_server_current_queue{namespace="jaeger",pod="jaeger-hotrod-85b766c49d-s6vz7",route="jaeger-hotrod-qq5cq",server="10.134.249.205:8080",service="jaeger-hotrod"} 0
haproxy_server_current_queue{namespace="jaeger",pod="jaeger-query-65d9446fc6-jx5qk",route="jaeger-query-pc6w5",server="10.134.248.54:16686",service="jaeger-query"} 0
haproxy_server_current_queue{namespace="kafka-service",pod="akhq-58744f76db-h2mtf",route="akhq-ingress-lhz5z",server="10.134.249.133:8080",service="akhq"} 0
haproxy_server_current_queue{namespace="license-server",pod="license-server-7b8b6f765f-wvcx5",route="license-server-9txj9",server="10.132.148.47:8080",service="license-server"} 0
haproxy_server_current_queue{namespace="opensearch",pod="opensearch-dashboards-764bc8548b-zngcr",route="opensearch-dashboards-l6vdv",server="10.134.173.229:5601",service="opensearch-dashboards"} 0
haproxy_server_current_queue{namespace="openshift-authentication",pod="oauth-openshift-75965c7487-kc8m7",route="oauth-openshift",server="10.132.153.132:6443",service="oauth-openshift"} 0
haproxy_server_current_queue{namespace="openshift-authentication",pod="oauth-openshift-75965c7487-nsfxw",route="oauth-openshift",server="10.133.43.151:6443",service="oauth-openshift"} 0
haproxy_server_current_queue{namespace="openshift-authentication",pod="oauth-openshift-75965c7487-xmmmq",route="oauth-openshift",server="10.135.207.187:6443",service="oauth-openshift"} 0
haproxy_server_current_queue{namespace="openshift-console",pod="console-66fcc4f4fc-crbz6",route="console",server="10.133.43.162:8443",service="console"} 0
haproxy_server_current_queue{namespace="openshift-console",pod="console-66fcc4f4fc-j5dvz",route="console",server="10.132.153.185:8443",service="console"} 0
haproxy_server_current_queue{namespace="openshift-console",pod="downloads-6fb844fb47-cwnvw",route="downloads",server="10.135.207.221:8080",service="downloads"} 0
haproxy_server_current_queue{namespace="openshift-console",pod="downloads-6fb844fb47-l485s",route="downloads",server="10.132.153.222:8080",service="downloads"} 0
haproxy_server_current_queue{namespace="openshift-ingress-canary",pod="ingress-canary-7rvzj",route="canary",server="10.134.173.247:8080",service="ingress-canary"} 0
haproxy_server_current_queue{namespace="openshift-ingress-canary",pod="ingress-canary-lwpkn",route="canary",server="10.132.149.137:8080",service="ingress-canary"} 0
haproxy_server_current_queue{namespace="openshift-ingress-canary",pod="ingress-canary-z2b4x",route="canary",server="10.134.248.101:8080",service="ingress-canary"} 0
haproxy_server_current_queue{namespace="openshift-monitoring",pod="alertmanager-main-0",route="alertmanager-main",server="10.134.173.99:9095",service="alertmanager-main"} 0
haproxy_server_current_queue{namespace="openshift-monitoring",pod="alertmanager-main-1",route="alertmanager-main",server="10.134.248.127:9095",service="alertmanager-main"} 0
haproxy_server_current_queue{namespace="openshift-monitoring",pod="prometheus-k8s-0",route="prometheus-k8s",server="10.134.173.51:9091",service="prometheus-k8s"} 0
haproxy_server_current_queue{namespace="openshift-monitoring",pod="prometheus-k8s-0",route="prometheus-k8s-federate",server="10.134.173.51:9091",service="prometheus-k8s"} 0
haproxy_server_current_queue{namespace="openshift-monitoring",pod="prometheus-k8s-1",route="prometheus-k8s",server="10.134.249.153:9091",service="prometheus-k8s"} 0
haproxy_server_current_queue{namespace="openshift-monitoring",pod="prometheus-k8s-1",route="prometheus-k8s-federate",server="10.134.249.153:9091",service="prometheus-k8s"} 0
haproxy_server_current_queue{namespace="openshift-monitoring",pod="thanos-querier-59d7db9787-hwjzt",route="thanos-querier",server="10.134.249.65:9091",service="thanos-querier"} 0
haproxy_server_current_queue{namespace="openshift-monitoring",pod="thanos-querier-59d7db9787-p8m9l",route="thanos-querier",server="10.134.173.112:9091",service="thanos-querier"} 0
haproxy_server_current_queue{namespace="profiler",pod="esc-static-service-654f8bc589-xchjk",route="esc-static-service-j9sgv",server="10.134.173.123:8080",service="esc-static-service"} 0
haproxy_server_current_queue{namespace="profiler",pod="esc-test-service-75849d7c6c-czvsw",route="esc-test-service-lbp8r",server="10.134.172.182:8080",service="esc-test-service"} 0
haproxy_server_current_queue{namespace="profiler",pod="esc-ui-service-66977c9bcb-kg9mc",route="esc-ui-service-r5tfm",server="10.134.249.215:8180",service="esc-ui-service"} 0
haproxy_server_current_queue{namespace="prometheus-operator",pod="grafana-deployment-88d46bc5f-5gx7h",route="prometheus-operator-grafana",server="10.132.149.3:3000",service="grafana-service"} 0
haproxy_server_current_queue{namespace="prometheus-operator",pod="vmagent-k8s-67654d6f97-98nf5",route="prometheus-operator-vmagent-k8s",server="10.132.149.156:8429",service="vmagent-k8s"} 0
haproxy_server_current_queue{namespace="prometheus-operator",pod="vmalert-k8s-6bcf966dd6-jvpz5",route="prometheus-operator-vmalert-k8s",server="10.132.149.175:8080",service="vmalert-k8s"} 0
haproxy_server_current_queue{namespace="prometheus-operator",pod="vmalertmanager-k8s-0",route="prometheus-operator-vmalertmanager-k8s",server="10.134.173.133:9093",service="vmalertmanager-k8s"} 0
haproxy_server_current_queue{namespace="prometheus-operator",pod="vmsingle-k8s-f4844bfc4-5l8sf",route="prometheus-operator-vmsingle-k8s",server="10.134.249.110:8429",service="vmsingle-k8s"} 0
haproxy_server_current_queue{namespace="rabbitmq",pod="rmqlocal-0",route="rabbitmq-fnzhp",server="10.134.249.141:15672",service="rabbitmq"} 0
haproxy_server_current_queue{namespace="rabbitmq",pod="rmqlocal-1",route="rabbitmq-fnzhp",server="10.134.173.127:15672",service="rabbitmq"} 0
haproxy_server_current_queue{namespace="rabbitmq",pod="rmqlocal-2",route="rabbitmq-fnzhp",server="10.132.149.46:15672",service="rabbitmq"} 0
haproxy_server_current_queue{namespace="shared-license-distributor",pod="shared-license-distributor-64ccc89449-8kw8p",route="shared-license-distributor-wnj68",server="10.134.249.239:8080",service="shared-license-distributor"} 0
haproxy_server_current_queue{namespace="site-manager",pod="site-manager-547546bdd5-hd85n",route="site-manager",server="10.134.248.81:8443",service="site-manager"} 0
haproxy_server_current_queue{namespace="spark",pod="spark-history-server-bd755d9f5-9ccsz",route="spark-history-server-hqshb",server="10.132.149.178:18080",service="spark-history-server"} 0
haproxy_server_current_queue{namespace="spark",pod="spark-history-server-bd755d9f5-zf2dt",route="spark-history-server-hqshb",server="10.134.173.147:18080",service="spark-history-server"} 0
haproxy_server_current_queue{namespace="spark",pod="spark-site-manager-66c44c7b78-cc7jr",route="spark-site-manager-p4nqh",server="10.132.149.148:8080",service="spark-site-manager"} 0
# HELP haproxy_server_current_session_rate Current number of sessions per second over last elapsed second.
# TYPE haproxy_server_current_session_rate gauge
haproxy_server_current_session_rate{namespace="",pod="",route="",server="fe_no_sni",service=""} 0
haproxy_server_current_session_rate{namespace="",pod="",route="",server="fe_sni",service=""} 0
haproxy_server_current_session_rate{namespace="airflow",pod="airflow-site-manager-68ccf87889-bwbdf",route="airflow-site-manager-hpl66",server="10.132.149.80:8080",service="airflow-site-manager"} 0
haproxy_server_current_session_rate{namespace="airflow",pod="airflow-webserver-695b5c545b-rxqc4",route="airflow-ingress-nzbtm",server="10.134.173.137:8080",service="airflow-webserver"} 0
haproxy_server_current_session_rate{namespace="arangodb",pod="main-arangodb-1-87669b688-q98z4",route="arangodb-leader-ingress-zw82p",server="10.132.149.185:8529",service="arangodb-leader"} 0
haproxy_server_current_session_rate{namespace="consul-service",pod="consul-server-0",route="consul-ingress-hmmvf",server="10.134.248.50:8500",service="consul-ui"} 0
haproxy_server_current_session_rate{namespace="consul-service",pod="consul-server-1",route="consul-ingress-hmmvf",server="10.132.149.160:8500",service="consul-ui"} 0
haproxy_server_current_session_rate{namespace="consul-service",pod="consul-server-2",route="consul-ingress-hmmvf",server="10.134.173.136:8500",service="consul-ui"} 0
haproxy_server_current_session_rate{namespace="dbaas",pod="dbaas-aggregator-5f744d74f5-r7tqm",route="aggregator-fw5ft",server="10.134.248.185:8080",service="dbaas-aggregator"} 0
haproxy_server_current_session_rate{namespace="infra-keycloak",pod="infra-keycloak-7fd57cd59c-d6c5b",route="infra-keycloak-5q4rj",server="10.132.148.223:8080",service="infra-keycloak"} 0
haproxy_server_current_session_rate{namespace="infra-keycloak",pod="infra-keycloak-7fd57cd59c-qtqv9",route="infra-keycloak-5q4rj",server="10.134.173.182:8080",service="infra-keycloak"} 0
haproxy_server_current_session_rate{namespace="jaeger",pod="jaeger-hotrod-85b766c49d-s6vz7",route="jaeger-hotrod-qq5cq",server="10.134.249.205:8080",service="jaeger-hotrod"} 0
haproxy_server_current_session_rate{namespace="jaeger",pod="jaeger-query-65d9446fc6-jx5qk",route="jaeger-query-pc6w5",server="10.134.248.54:16686",service="jaeger-query"} 0
haproxy_server_current_session_rate{namespace="kafka-service",pod="akhq-58744f76db-h2mtf",route="akhq-ingress-lhz5z",server="10.134.249.133:8080",service="akhq"} 0
haproxy_server_current_session_rate{namespace="license-server",pod="license-server-7b8b6f765f-wvcx5",route="license-server-9txj9",server="10.132.148.47:8080",service="license-server"} 0
haproxy_server_current_session_rate{namespace="opensearch",pod="opensearch-dashboards-764bc8548b-zngcr",route="opensearch-dashboards-l6vdv",server="10.134.173.229:5601",service="opensearch-dashboards"} 0
haproxy_server_current_session_rate{namespace="openshift-authentication",pod="oauth-openshift-75965c7487-kc8m7",route="oauth-openshift",server="10.132.153.132:6443",service="oauth-openshift"} 0
haproxy_server_current_session_rate{namespace="openshift-authentication",pod="oauth-openshift-75965c7487-nsfxw",route="oauth-openshift",server="10.133.43.151:6443",service="oauth-openshift"} 0
haproxy_server_current_session_rate{namespace="openshift-authentication",pod="oauth-openshift-75965c7487-xmmmq",route="oauth-openshift",server="10.135.207.187:6443",service="oauth-openshift"} 0
haproxy_server_current_session_rate{namespace="openshift-console",pod="console-66fcc4f4fc-crbz6",route="console",server="10.133.43.162:8443",service="console"} 0
haproxy_server_current_session_rate{namespace="openshift-console",pod="console-66fcc4f4fc-j5dvz",route="console",server="10.132.153.185:8443",service="console"} 0
haproxy_server_current_session_rate{namespace="openshift-console",pod="downloads-6fb844fb47-cwnvw",route="downloads",server="10.135.207.221:8080",service="downloads"} 0
haproxy_server_current_session_rate{namespace="openshift-console",pod="downloads-6fb844fb47-l485s",route="downloads",server="10.132.153.222:8080",service="downloads"} 0
haproxy_server_current_session_rate{namespace="openshift-ingress-canary",pod="ingress-canary-7rvzj",route="canary",server="10.134.173.247:8080",service="ingress-canary"} 0
haproxy_server_current_session_rate{namespace="openshift-ingress-canary",pod="ingress-canary-lwpkn",route="canary",server="10.132.149.137:8080",service="ingress-canary"} 0
haproxy_server_current_session_rate{namespace="openshift-ingress-canary",pod="ingress-canary-z2b4x",route="canary",server="10.134.248.101:8080",service="ingress-canary"} 0
haproxy_server_current_session_rate{namespace="openshift-monitoring",pod="alertmanager-main-0",route="alertmanager-main",server="10.134.173.99:9095",service="alertmanager-main"} 0
haproxy_server_current_session_rate{namespace="openshift-monitoring",pod="alertmanager-main-1",route="alertmanager-main",server="10.134.248.127:9095",service="alertmanager-main"} 0
haproxy_server_current_session_rate{namespace="openshift-monitoring",pod="prometheus-k8s-0",route="prometheus-k8s",server="10.134.173.51:9091",service="prometheus-k8s"} 0
haproxy_server_current_session_rate{namespace="openshift-monitoring",pod="prometheus-k8s-0",route="prometheus-k8s-federate",server="10.134.173.51:9091",service="prometheus-k8s"} 0
haproxy_server_current_session_rate{namespace="openshift-monitoring",pod="prometheus-k8s-1",route="prometheus-k8s",server="10.134.249.153:9091",service="prometheus-k8s"} 0
haproxy_server_current_session_rate{namespace="openshift-monitoring",pod="prometheus-k8s-1",route="prometheus-k8s-federate",server="10.134.249.153:9091",service="prometheus-k8s"} 0
haproxy_server_current_session_rate{namespace="openshift-monitoring",pod="thanos-querier-59d7db9787-hwjzt",route="thanos-querier",server="10.134.249.65:9091",service="thanos-querier"} 0
haproxy_server_current_session_rate{namespace="openshift-monitoring",pod="thanos-querier-59d7db9787-p8m9l",route="thanos-querier",server="10.134.173.112:9091",service="thanos-querier"} 0
haproxy_server_current_session_rate{namespace="profiler",pod="esc-static-service-654f8bc589-xchjk",route="esc-static-service-j9sgv",server="10.134.173.123:8080",service="esc-static-service"} 0
haproxy_server_current_session_rate{namespace="profiler",pod="esc-test-service-75849d7c6c-czvsw",route="esc-test-service-lbp8r",server="10.134.172.182:8080",service="esc-test-service"} 0
haproxy_server_current_session_rate{namespace="profiler",pod="esc-ui-service-66977c9bcb-kg9mc",route="esc-ui-service-r5tfm",server="10.134.249.215:8180",service="esc-ui-service"} 0
haproxy_server_current_session_rate{namespace="prometheus-operator",pod="grafana-deployment-88d46bc5f-5gx7h",route="prometheus-operator-grafana",server="10.132.149.3:3000",service="grafana-service"} 0
haproxy_server_current_session_rate{namespace="prometheus-operator",pod="vmagent-k8s-67654d6f97-98nf5",route="prometheus-operator-vmagent-k8s",server="10.132.149.156:8429",service="vmagent-k8s"} 0
haproxy_server_current_session_rate{namespace="prometheus-operator",pod="vmalert-k8s-6bcf966dd6-jvpz5",route="prometheus-operator-vmalert-k8s",server="10.132.149.175:8080",service="vmalert-k8s"} 0
haproxy_server_current_session_rate{namespace="prometheus-operator",pod="vmalertmanager-k8s-0",route="prometheus-operator-vmalertmanager-k8s",server="10.134.173.133:9093",service="vmalertmanager-k8s"} 0
haproxy_server_current_session_rate{namespace="prometheus-operator",pod="vmsingle-k8s-f4844bfc4-5l8sf",route="prometheus-operator-vmsingle-k8s",server="10.134.249.110:8429",service="vmsingle-k8s"} 0
haproxy_server_current_session_rate{namespace="rabbitmq",pod="rmqlocal-0",route="rabbitmq-fnzhp",server="10.134.249.141:15672",service="rabbitmq"} 0
haproxy_server_current_session_rate{namespace="rabbitmq",pod="rmqlocal-1",route="rabbitmq-fnzhp",server="10.134.173.127:15672",service="rabbitmq"} 0
haproxy_server_current_session_rate{namespace="rabbitmq",pod="rmqlocal-2",route="rabbitmq-fnzhp",server="10.132.149.46:15672",service="rabbitmq"} 0
haproxy_server_current_session_rate{namespace="shared-license-distributor",pod="shared-license-distributor-64ccc89449-8kw8p",route="shared-license-distributor-wnj68",server="10.134.249.239:8080",service="shared-license-distributor"} 0
haproxy_server_current_session_rate{namespace="site-manager",pod="site-manager-547546bdd5-hd85n",route="site-manager",server="10.134.248.81:8443",service="site-manager"} 0
haproxy_server_current_session_rate{namespace="spark",pod="spark-history-server-bd755d9f5-9ccsz",route="spark-history-server-hqshb",server="10.132.149.178:18080",service="spark-history-server"} 0
haproxy_server_current_session_rate{namespace="spark",pod="spark-history-server-bd755d9f5-zf2dt",route="spark-history-server-hqshb",server="10.134.173.147:18080",service="spark-history-server"} 0
haproxy_server_current_session_rate{namespace="spark",pod="spark-site-manager-66c44c7b78-cc7jr",route="spark-site-manager-p4nqh",server="10.132.149.148:8080",service="spark-site-manager"} 0
# HELP haproxy_server_current_sessions Current number of active sessions.
# TYPE haproxy_server_current_sessions gauge
haproxy_server_current_sessions{namespace="",pod="",route="",server="fe_no_sni",service=""} 0
haproxy_server_current_sessions{namespace="",pod="",route="",server="fe_sni",service=""} 0
haproxy_server_current_sessions{namespace="airflow",pod="airflow-site-manager-68ccf87889-bwbdf",route="airflow-site-manager-hpl66",server="10.132.149.80:8080",service="airflow-site-manager"} 0
haproxy_server_current_sessions{namespace="airflow",pod="airflow-webserver-695b5c545b-rxqc4",route="airflow-ingress-nzbtm",server="10.134.173.137:8080",service="airflow-webserver"} 0
haproxy_server_current_sessions{namespace="arangodb",pod="main-arangodb-1-87669b688-q98z4",route="arangodb-leader-ingress-zw82p",server="10.132.149.185:8529",service="arangodb-leader"} 0
haproxy_server_current_sessions{namespace="consul-service",pod="consul-server-0",route="consul-ingress-hmmvf",server="10.134.248.50:8500",service="consul-ui"} 0
haproxy_server_current_sessions{namespace="consul-service",pod="consul-server-1",route="consul-ingress-hmmvf",server="10.132.149.160:8500",service="consul-ui"} 0
haproxy_server_current_sessions{namespace="consul-service",pod="consul-server-2",route="consul-ingress-hmmvf",server="10.134.173.136:8500",service="consul-ui"} 0
haproxy_server_current_sessions{namespace="dbaas",pod="dbaas-aggregator-5f744d74f5-r7tqm",route="aggregator-fw5ft",server="10.134.248.185:8080",service="dbaas-aggregator"} 0
haproxy_server_current_sessions{namespace="infra-keycloak",pod="infra-keycloak-7fd57cd59c-d6c5b",route="infra-keycloak-5q4rj",server="10.132.148.223:8080",service="infra-keycloak"} 0
haproxy_server_current_sessions{namespace="infra-keycloak",pod="infra-keycloak-7fd57cd59c-qtqv9",route="infra-keycloak-5q4rj",server="10.134.173.182:8080",service="infra-keycloak"} 0
haproxy_server_current_sessions{namespace="jaeger",pod="jaeger-hotrod-85b766c49d-s6vz7",route="jaeger-hotrod-qq5cq",server="10.134.249.205:8080",service="jaeger-hotrod"} 0
haproxy_server_current_sessions{namespace="jaeger",pod="jaeger-query-65d9446fc6-jx5qk",route="jaeger-query-pc6w5",server="10.134.248.54:16686",service="jaeger-query"} 0
haproxy_server_current_sessions{namespace="kafka-service",pod="akhq-58744f76db-h2mtf",route="akhq-ingress-lhz5z",server="10.134.249.133:8080",service="akhq"} 0
haproxy_server_current_sessions{namespace="license-server",pod="license-server-7b8b6f765f-wvcx5",route="license-server-9txj9",server="10.132.148.47:8080",service="license-server"} 0
haproxy_server_current_sessions{namespace="opensearch",pod="opensearch-dashboards-764bc8548b-zngcr",route="opensearch-dashboards-l6vdv",server="10.134.173.229:5601",service="opensearch-dashboards"} 0
haproxy_server_current_sessions{namespace="openshift-authentication",pod="oauth-openshift-75965c7487-kc8m7",route="oauth-openshift",server="10.132.153.132:6443",service="oauth-openshift"} 0
haproxy_server_current_sessions{namespace="openshift-authentication",pod="oauth-openshift-75965c7487-nsfxw",route="oauth-openshift",server="10.133.43.151:6443",service="oauth-openshift"} 0
haproxy_server_current_sessions{namespace="openshift-authentication",pod="oauth-openshift-75965c7487-xmmmq",route="oauth-openshift",server="10.135.207.187:6443",service="oauth-openshift"} 0
haproxy_server_current_sessions{namespace="openshift-console",pod="console-66fcc4f4fc-crbz6",route="console",server="10.133.43.162:8443",service="console"} 0
haproxy_server_current_sessions{namespace="openshift-console",pod="console-66fcc4f4fc-j5dvz",route="console",server="10.132.153.185:8443",service="console"} 0
haproxy_server_current_sessions{namespace="openshift-console",pod="downloads-6fb844fb47-cwnvw",route="downloads",server="10.135.207.221:8080",service="downloads"} 0
haproxy_server_current_sessions{namespace="openshift-console",pod="downloads-6fb844fb47-l485s",route="downloads",server="10.132.153.222:8080",service="downloads"} 0
haproxy_server_current_sessions{namespace="openshift-ingress-canary",pod="ingress-canary-7rvzj",route="canary",server="10.134.173.247:8080",service="ingress-canary"} 0
haproxy_server_current_sessions{namespace="openshift-ingress-canary",pod="ingress-canary-lwpkn",route="canary",server="10.132.149.137:8080",service="ingress-canary"} 0
haproxy_server_current_sessions{namespace="openshift-ingress-canary",pod="ingress-canary-z2b4x",route="canary",server="10.134.248.101:8080",service="ingress-canary"} 0
haproxy_server_current_sessions{namespace="openshift-monitoring",pod="alertmanager-main-0",route="alertmanager-main",server="10.134.173.99:9095",service="alertmanager-main"} 0
haproxy_server_current_sessions{namespace="openshift-monitoring",pod="alertmanager-main-1",route="alertmanager-main",server="10.134.248.127:9095",service="alertmanager-main"} 0
haproxy_server_current_sessions{namespace="openshift-monitoring",pod="prometheus-k8s-0",route="prometheus-k8s",server="10.134.173.51:9091",service="prometheus-k8s"} 0
haproxy_server_current_sessions{namespace="openshift-monitoring",pod="prometheus-k8s-0",route="prometheus-k8s-federate",server="10.134.173.51:9091",service="prometheus-k8s"} 0
haproxy_server_current_sessions{namespace="openshift-monitoring",pod="prometheus-k8s-1",route="prometheus-k8s",server="10.134.249.153:9091",service="prometheus-k8s"} 0
haproxy_server_current_sessions{namespace="openshift-monitoring",pod="prometheus-k8s-1",route="prometheus-k8s-federate",server="10.134.249.153:9091",service="prometheus-k8s"} 0
haproxy_server_current_sessions{namespace="openshift-monitoring",pod="thanos-querier-59d7db9787-hwjzt",route="thanos-querier",server="10.134.249.65:9091",service="thanos-querier"} 0
haproxy_server_current_sessions{namespace="openshift-monitoring",pod="thanos-querier-59d7db9787-p8m9l",route="thanos-querier",server="10.134.173.112:9091",service="thanos-querier"} 0
haproxy_server_current_sessions{namespace="profiler",pod="esc-static-service-654f8bc589-xchjk",route="esc-static-service-j9sgv",server="10.134.173.123:8080",service="esc-static-service"} 0
haproxy_server_current_sessions{namespace="profiler",pod="esc-test-service-75849d7c6c-czvsw",route="esc-test-service-lbp8r",server="10.134.172.182:8080",service="esc-test-service"} 0
haproxy_server_current_sessions{namespace="profiler",pod="esc-ui-service-66977c9bcb-kg9mc",route="esc-ui-service-r5tfm",server="10.134.249.215:8180",service="esc-ui-service"} 0
haproxy_server_current_sessions{namespace="prometheus-operator",pod="grafana-deployment-88d46bc5f-5gx7h",route="prometheus-operator-grafana",server="10.132.149.3:3000",service="grafana-service"} 0
haproxy_server_current_sessions{namespace="prometheus-operator",pod="vmagent-k8s-67654d6f97-98nf5",route="prometheus-operator-vmagent-k8s",server="10.132.149.156:8429",service="vmagent-k8s"} 0
haproxy_server_current_sessions{namespace="prometheus-operator",pod="vmalert-k8s-6bcf966dd6-jvpz5",route="prometheus-operator-vmalert-k8s",server="10.132.149.175:8080",service="vmalert-k8s"} 0
haproxy_server_current_sessions{namespace="prometheus-operator",pod="vmalertmanager-k8s-0",route="prometheus-operator-vmalertmanager-k8s",server="10.134.173.133:9093",service="vmalertmanager-k8s"} 0
haproxy_server_current_sessions{namespace="prometheus-operator",pod="vmsingle-k8s-f4844bfc4-5l8sf",route="prometheus-operator-vmsingle-k8s",server="10.134.249.110:8429",service="vmsingle-k8s"} 0
haproxy_server_current_sessions{namespace="rabbitmq",pod="rmqlocal-0",route="rabbitmq-fnzhp",server="10.134.249.141:15672",service="rabbitmq"} 0
haproxy_server_current_sessions{namespace="rabbitmq",pod="rmqlocal-1",route="rabbitmq-fnzhp",server="10.134.173.127:15672",service="rabbitmq"} 0
haproxy_server_current_sessions{namespace="rabbitmq",pod="rmqlocal-2",route="rabbitmq-fnzhp",server="10.132.149.46:15672",service="rabbitmq"} 0
haproxy_server_current_sessions{namespace="shared-license-distributor",pod="shared-license-distributor-64ccc89449-8kw8p",route="shared-license-distributor-wnj68",server="10.134.249.239:8080",service="shared-license-distributor"} 0
haproxy_server_current_sessions{namespace="site-manager",pod="site-manager-547546bdd5-hd85n",route="site-manager",server="10.134.248.81:8443",service="site-manager"} 0
haproxy_server_current_sessions{namespace="spark",pod="spark-history-server-bd755d9f5-9ccsz",route="spark-history-server-hqshb",server="10.132.149.178:18080",service="spark-history-server"} 0
haproxy_server_current_sessions{namespace="spark",pod="spark-history-server-bd755d9f5-zf2dt",route="spark-history-server-hqshb",server="10.134.173.147:18080",service="spark-history-server"} 0
haproxy_server_current_sessions{namespace="spark",pod="spark-site-manager-66c44c7b78-cc7jr",route="spark-site-manager-p4nqh",server="10.132.149.148:8080",service="spark-site-manager"} 0
# HELP haproxy_server_downtime_seconds_total Total downtime in seconds.
# TYPE haproxy_server_downtime_seconds_total gauge
haproxy_server_downtime_seconds_total{namespace="consul-service",pod="consul-server-0",route="consul-ingress-hmmvf",server="10.134.248.50:8500",service="consul-ui"} 0
haproxy_server_downtime_seconds_total{namespace="consul-service",pod="consul-server-1",route="consul-ingress-hmmvf",server="10.132.149.160:8500",service="consul-ui"} 0
haproxy_server_downtime_seconds_total{namespace="consul-service",pod="consul-server-2",route="consul-ingress-hmmvf",server="10.134.173.136:8500",service="consul-ui"} 0
haproxy_server_downtime_seconds_total{namespace="infra-keycloak",pod="infra-keycloak-7fd57cd59c-d6c5b",route="infra-keycloak-5q4rj",server="10.132.148.223:8080",service="infra-keycloak"} 0
haproxy_server_downtime_seconds_total{namespace="infra-keycloak",pod="infra-keycloak-7fd57cd59c-qtqv9",route="infra-keycloak-5q4rj",server="10.134.173.182:8080",service="infra-keycloak"} 0
haproxy_server_downtime_seconds_total{namespace="openshift-authentication",pod="oauth-openshift-75965c7487-kc8m7",route="oauth-openshift",server="10.132.153.132:6443",service="oauth-openshift"} 0
haproxy_server_downtime_seconds_total{namespace="openshift-authentication",pod="oauth-openshift-75965c7487-nsfxw",route="oauth-openshift",server="10.133.43.151:6443",service="oauth-openshift"} 0
haproxy_server_downtime_seconds_total{namespace="openshift-authentication",pod="oauth-openshift-75965c7487-xmmmq",route="oauth-openshift",server="10.135.207.187:6443",service="oauth-openshift"} 0
haproxy_server_downtime_seconds_total{namespace="openshift-console",pod="console-66fcc4f4fc-crbz6",route="console",server="10.133.43.162:8443",service="console"} 0
haproxy_server_downtime_seconds_total{namespace="openshift-console",pod="console-66fcc4f4fc-j5dvz",route="console",server="10.132.153.185:8443",service="console"} 0
haproxy_server_downtime_seconds_total{namespace="openshift-console",pod="downloads-6fb844fb47-cwnvw",route="downloads",server="10.135.207.221:8080",service="downloads"} 0
haproxy_server_downtime_seconds_total{namespace="openshift-console",pod="downloads-6fb844fb47-l485s",route="downloads",server="10.132.153.222:8080",service="downloads"} 0
haproxy_server_downtime_seconds_total{namespace="openshift-ingress-canary",pod="ingress-canary-7rvzj",route="canary",server="10.134.173.247:8080",service="ingress-canary"} 0
haproxy_server_downtime_seconds_total{namespace="openshift-ingress-canary",pod="ingress-canary-lwpkn",route="canary",server="10.132.149.137:8080",service="ingress-canary"} 0
haproxy_server_downtime_seconds_total{namespace="openshift-ingress-canary",pod="ingress-canary-z2b4x",route="canary",server="10.134.248.101:8080",service="ingress-canary"} 0
haproxy_server_downtime_seconds_total{namespace="openshift-monitoring",pod="alertmanager-main-0",route="alertmanager-main",server="10.134.173.99:9095",service="alertmanager-main"} 0
haproxy_server_downtime_seconds_total{namespace="openshift-monitoring",pod="alertmanager-main-1",route="alertmanager-main",server="10.134.248.127:9095",service="alertmanager-main"} 0
haproxy_server_downtime_seconds_total{namespace="openshift-monitoring",pod="prometheus-k8s-0",route="prometheus-k8s",server="10.134.173.51:9091",service="prometheus-k8s"} 0
haproxy_server_downtime_seconds_total{namespace="openshift-monitoring",pod="prometheus-k8s-0",route="prometheus-k8s-federate",server="10.134.173.51:9091",service="prometheus-k8s"} 0
haproxy_server_downtime_seconds_total{namespace="openshift-monitoring",pod="prometheus-k8s-1",route="prometheus-k8s",server="10.134.249.153:9091",service="prometheus-k8s"} 0
haproxy_server_downtime_seconds_total{namespace="openshift-monitoring",pod="prometheus-k8s-1",route="prometheus-k8s-federate",server="10.134.249.153:9091",service="prometheus-k8s"} 0
haproxy_server_downtime_seconds_total{namespace="openshift-monitoring",pod="thanos-querier-59d7db9787-hwjzt",route="thanos-querier",server="10.134.249.65:9091",service="thanos-querier"} 0
haproxy_server_downtime_seconds_total{namespace="openshift-monitoring",pod="thanos-querier-59d7db9787-p8m9l",route="thanos-querier",server="10.134.173.112:9091",service="thanos-querier"} 0
haproxy_server_downtime_seconds_total{namespace="rabbitmq",pod="rmqlocal-0",route="rabbitmq-fnzhp",server="10.134.249.141:15672",service="rabbitmq"} 46
haproxy_server_downtime_seconds_total{namespace="rabbitmq",pod="rmqlocal-1",route="rabbitmq-fnzhp",server="10.134.173.127:15672",service="rabbitmq"} 46
haproxy_server_downtime_seconds_total{namespace="rabbitmq",pod="rmqlocal-2",route="rabbitmq-fnzhp",server="10.132.149.46:15672",service="rabbitmq"} 11
haproxy_server_downtime_seconds_total{namespace="spark",pod="spark-history-server-bd755d9f5-9ccsz",route="spark-history-server-hqshb",server="10.132.149.178:18080",service="spark-history-server"} 20
haproxy_server_downtime_seconds_total{namespace="spark",pod="spark-history-server-bd755d9f5-zf2dt",route="spark-history-server-hqshb",server="10.134.173.147:18080",service="spark-history-server"} 0
# HELP haproxy_server_http_average_connect_latency_milliseconds Average connect latency of the last 1024 requests in milliseconds.
# TYPE haproxy_server_http_average_connect_latency_milliseconds gauge
haproxy_server_http_average_connect_latency_milliseconds{namespace="",pod="",route="",server="fe_no_sni",service=""} 0
haproxy_server_http_average_connect_latency_milliseconds{namespace="",pod="",route="",server="fe_sni",service=""} 0
haproxy_server_http_average_connect_latency_milliseconds{namespace="airflow",pod="airflow-site-manager-68ccf87889-bwbdf",route="airflow-site-manager-hpl66",server="10.132.149.80:8080",service="airflow-site-manager"} 0
haproxy_server_http_average_connect_latency_milliseconds{namespace="airflow",pod="airflow-webserver-695b5c545b-rxqc4",route="airflow-ingress-nzbtm",server="10.134.173.137:8080",service="airflow-webserver"} 0
haproxy_server_http_average_connect_latency_milliseconds{namespace="arangodb",pod="main-arangodb-1-87669b688-q98z4",route="arangodb-leader-ingress-zw82p",server="10.132.149.185:8529",service="arangodb-leader"} 0
haproxy_server_http_average_connect_latency_milliseconds{namespace="consul-service",pod="consul-server-0",route="consul-ingress-hmmvf",server="10.134.248.50:8500",service="consul-ui"} 0
haproxy_server_http_average_connect_latency_milliseconds{namespace="consul-service",pod="consul-server-1",route="consul-ingress-hmmvf",server="10.132.149.160:8500",service="consul-ui"} 0
haproxy_server_http_average_connect_latency_milliseconds{namespace="consul-service",pod="consul-server-2",route="consul-ingress-hmmvf",server="10.134.173.136:8500",service="consul-ui"} 0
haproxy_server_http_average_connect_latency_milliseconds{namespace="dbaas",pod="dbaas-aggregator-5f744d74f5-r7tqm",route="aggregator-fw5ft",server="10.134.248.185:8080",service="dbaas-aggregator"} 0
haproxy_server_http_average_connect_latency_milliseconds{namespace="infra-keycloak",pod="infra-keycloak-7fd57cd59c-d6c5b",route="infra-keycloak-5q4rj",server="10.132.148.223:8080",service="infra-keycloak"} 0
haproxy_server_http_average_connect_latency_milliseconds{namespace="infra-keycloak",pod="infra-keycloak-7fd57cd59c-qtqv9",route="infra-keycloak-5q4rj",server="10.134.173.182:8080",service="infra-keycloak"} 0
haproxy_server_http_average_connect_latency_milliseconds{namespace="jaeger",pod="jaeger-hotrod-85b766c49d-s6vz7",route="jaeger-hotrod-qq5cq",server="10.134.249.205:8080",service="jaeger-hotrod"} 0
haproxy_server_http_average_connect_latency_milliseconds{namespace="jaeger",pod="jaeger-query-65d9446fc6-jx5qk",route="jaeger-query-pc6w5",server="10.134.248.54:16686",service="jaeger-query"} 0
haproxy_server_http_average_connect_latency_milliseconds{namespace="kafka-service",pod="akhq-58744f76db-h2mtf",route="akhq-ingress-lhz5z",server="10.134.249.133:8080",service="akhq"} 0
haproxy_server_http_average_connect_latency_milliseconds{namespace="license-server",pod="license-server-7b8b6f765f-wvcx5",route="license-server-9txj9",server="10.132.148.47:8080",service="license-server"} 0
haproxy_server_http_average_connect_latency_milliseconds{namespace="opensearch",pod="opensearch-dashboards-764bc8548b-zngcr",route="opensearch-dashboards-l6vdv",server="10.134.173.229:5601",service="opensearch-dashboards"} 0
haproxy_server_http_average_connect_latency_milliseconds{namespace="openshift-authentication",pod="oauth-openshift-75965c7487-kc8m7",route="oauth-openshift",server="10.132.153.132:6443",service="oauth-openshift"} 0
haproxy_server_http_average_connect_latency_milliseconds{namespace="openshift-authentication",pod="oauth-openshift-75965c7487-nsfxw",route="oauth-openshift",server="10.133.43.151:6443",service="oauth-openshift"} 0
haproxy_server_http_average_connect_latency_milliseconds{namespace="openshift-authentication",pod="oauth-openshift-75965c7487-xmmmq",route="oauth-openshift",server="10.135.207.187:6443",service="oauth-openshift"} 0
haproxy_server_http_average_connect_latency_milliseconds{namespace="openshift-console",pod="console-66fcc4f4fc-crbz6",route="console",server="10.133.43.162:8443",service="console"} 0
haproxy_server_http_average_connect_latency_milliseconds{namespace="openshift-console",pod="console-66fcc4f4fc-j5dvz",route="console",server="10.132.153.185:8443",service="console"} 0
haproxy_server_http_average_connect_latency_milliseconds{namespace="openshift-console",pod="downloads-6fb844fb47-cwnvw",route="downloads",server="10.135.207.221:8080",service="downloads"} 0
haproxy_server_http_average_connect_latency_milliseconds{namespace="openshift-console",pod="downloads-6fb844fb47-l485s",route="downloads",server="10.132.153.222:8080",service="downloads"} 0
haproxy_server_http_average_connect_latency_milliseconds{namespace="openshift-ingress-canary",pod="ingress-canary-7rvzj",route="canary",server="10.134.173.247:8080",service="ingress-canary"} 0
haproxy_server_http_average_connect_latency_milliseconds{namespace="openshift-ingress-canary",pod="ingress-canary-lwpkn",route="canary",server="10.132.149.137:8080",service="ingress-canary"} 0
haproxy_server_http_average_connect_latency_milliseconds{namespace="openshift-ingress-canary",pod="ingress-canary-z2b4x",route="canary",server="10.134.248.101:8080",service="ingress-canary"} 0
haproxy_server_http_average_connect_latency_milliseconds{namespace="openshift-monitoring",pod="alertmanager-main-0",route="alertmanager-main",server="10.134.173.99:9095",service="alertmanager-main"} 0
haproxy_server_http_average_connect_latency_milliseconds{namespace="openshift-monitoring",pod="alertmanager-main-1",route="alertmanager-main",server="10.134.248.127:9095",service="alertmanager-main"} 0
haproxy_server_http_average_connect_latency_milliseconds{namespace="openshift-monitoring",pod="prometheus-k8s-0",route="prometheus-k8s",server="10.134.173.51:9091",service="prometheus-k8s"} 0
haproxy_server_http_average_connect_latency_milliseconds{namespace="openshift-monitoring",pod="prometheus-k8s-0",route="prometheus-k8s-federate",server="10.134.173.51:9091",service="prometheus-k8s"} 0
haproxy_server_http_average_connect_latency_milliseconds{namespace="openshift-monitoring",pod="prometheus-k8s-1",route="prometheus-k8s",server="10.134.249.153:9091",service="prometheus-k8s"} 0
haproxy_server_http_average_connect_latency_milliseconds{namespace="openshift-monitoring",pod="prometheus-k8s-1",route="prometheus-k8s-federate",server="10.134.249.153:9091",service="prometheus-k8s"} 0
haproxy_server_http_average_connect_latency_milliseconds{namespace="openshift-monitoring",pod="thanos-querier-59d7db9787-hwjzt",route="thanos-querier",server="10.134.249.65:9091",service="thanos-querier"} 0
haproxy_server_http_average_connect_latency_milliseconds{namespace="openshift-monitoring",pod="thanos-querier-59d7db9787-p8m9l",route="thanos-querier",server="10.134.173.112:9091",service="thanos-querier"} 0
haproxy_server_http_average_connect_latency_milliseconds{namespace="profiler",pod="esc-static-service-654f8bc589-xchjk",route="esc-static-service-j9sgv",server="10.134.173.123:8080",service="esc-static-service"} 0
haproxy_server_http_average_connect_latency_milliseconds{namespace="profiler",pod="esc-test-service-75849d7c6c-czvsw",route="esc-test-service-lbp8r",server="10.134.172.182:8080",service="esc-test-service"} 0
haproxy_server_http_average_connect_latency_milliseconds{namespace="profiler",pod="esc-ui-service-66977c9bcb-kg9mc",route="esc-ui-service-r5tfm",server="10.134.249.215:8180",service="esc-ui-service"} 0
haproxy_server_http_average_connect_latency_milliseconds{namespace="prometheus-operator",pod="grafana-deployment-88d46bc5f-5gx7h",route="prometheus-operator-grafana",server="10.132.149.3:3000",service="grafana-service"} 0
haproxy_server_http_average_connect_latency_milliseconds{namespace="prometheus-operator",pod="vmagent-k8s-67654d6f97-98nf5",route="prometheus-operator-vmagent-k8s",server="10.132.149.156:8429",service="vmagent-k8s"} 0
haproxy_server_http_average_connect_latency_milliseconds{namespace="prometheus-operator",pod="vmalert-k8s-6bcf966dd6-jvpz5",route="prometheus-operator-vmalert-k8s",server="10.132.149.175:8080",service="vmalert-k8s"} 0
haproxy_server_http_average_connect_latency_milliseconds{namespace="prometheus-operator",pod="vmalertmanager-k8s-0",route="prometheus-operator-vmalertmanager-k8s",server="10.134.173.133:9093",service="vmalertmanager-k8s"} 0
haproxy_server_http_average_connect_latency_milliseconds{namespace="prometheus-operator",pod="vmsingle-k8s-f4844bfc4-5l8sf",route="prometheus-operator-vmsingle-k8s",server="10.134.249.110:8429",service="vmsingle-k8s"} 0
haproxy_server_http_average_connect_latency_milliseconds{namespace="rabbitmq",pod="rmqlocal-0",route="rabbitmq-fnzhp",server="10.134.249.141:15672",service="rabbitmq"} 0
haproxy_server_http_average_connect_latency_milliseconds{namespace="rabbitmq",pod="rmqlocal-1",route="rabbitmq-fnzhp",server="10.134.173.127:15672",service="rabbitmq"} 0
haproxy_server_http_average_connect_latency_milliseconds{namespace="rabbitmq",pod="rmqlocal-2",route="rabbitmq-fnzhp",server="10.132.149.46:15672",service="rabbitmq"} 0
haproxy_server_http_average_connect_latency_milliseconds{namespace="shared-license-distributor",pod="shared-license-distributor-64ccc89449-8kw8p",route="shared-license-distributor-wnj68",server="10.134.249.239:8080",service="shared-license-distributor"} 0
haproxy_server_http_average_connect_latency_milliseconds{namespace="site-manager",pod="site-manager-547546bdd5-hd85n",route="site-manager",server="10.134.248.81:8443",service="site-manager"} 0
haproxy_server_http_average_connect_latency_milliseconds{namespace="spark",pod="spark-history-server-bd755d9f5-9ccsz",route="spark-history-server-hqshb",server="10.132.149.178:18080",service="spark-history-server"} 0
haproxy_server_http_average_connect_latency_milliseconds{namespace="spark",pod="spark-history-server-bd755d9f5-zf2dt",route="spark-history-server-hqshb",server="10.134.173.147:18080",service="spark-history-server"} 0
haproxy_server_http_average_connect_latency_milliseconds{namespace="spark",pod="spark-site-manager-66c44c7b78-cc7jr",route="spark-site-manager-p4nqh",server="10.132.149.148:8080",service="spark-site-manager"} 0
# HELP haproxy_server_http_average_queue_latency_milliseconds Average latency to be dequeued of the last 1024 requests in milliseconds.
# TYPE haproxy_server_http_average_queue_latency_milliseconds gauge
haproxy_server_http_average_queue_latency_milliseconds{namespace="",pod="",route="",server="fe_no_sni",service=""} 0
haproxy_server_http_average_queue_latency_milliseconds{namespace="",pod="",route="",server="fe_sni",service=""} 0
haproxy_server_http_average_queue_latency_milliseconds{namespace="airflow",pod="airflow-site-manager-68ccf87889-bwbdf",route="airflow-site-manager-hpl66",server="10.132.149.80:8080",service="airflow-site-manager"} 0
haproxy_server_http_average_queue_latency_milliseconds{namespace="airflow",pod="airflow-webserver-695b5c545b-rxqc4",route="airflow-ingress-nzbtm",server="10.134.173.137:8080",service="airflow-webserver"} 0
haproxy_server_http_average_queue_latency_milliseconds{namespace="arangodb",pod="main-arangodb-1-87669b688-q98z4",route="arangodb-leader-ingress-zw82p",server="10.132.149.185:8529",service="arangodb-leader"} 0
haproxy_server_http_average_queue_latency_milliseconds{namespace="consul-service",pod="consul-server-0",route="consul-ingress-hmmvf",server="10.134.248.50:8500",service="consul-ui"} 0
haproxy_server_http_average_queue_latency_milliseconds{namespace="consul-service",pod="consul-server-1",route="consul-ingress-hmmvf",server="10.132.149.160:8500",service="consul-ui"} 0
haproxy_server_http_average_queue_latency_milliseconds{namespace="consul-service",pod="consul-server-2",route="consul-ingress-hmmvf",server="10.134.173.136:8500",service="consul-ui"} 0
haproxy_server_http_average_queue_latency_milliseconds{namespace="dbaas",pod="dbaas-aggregator-5f744d74f5-r7tqm",route="aggregator-fw5ft",server="10.134.248.185:8080",service="dbaas-aggregator"} 0
haproxy_server_http_average_queue_latency_milliseconds{namespace="infra-keycloak",pod="infra-keycloak-7fd57cd59c-d6c5b",route="infra-keycloak-5q4rj",server="10.132.148.223:8080",service="infra-keycloak"} 0
haproxy_server_http_average_queue_latency_milliseconds{namespace="infra-keycloak",pod="infra-keycloak-7fd57cd59c-qtqv9",route="infra-keycloak-5q4rj",server="10.134.173.182:8080",service="infra-keycloak"} 0
haproxy_server_http_average_queue_latency_milliseconds{namespace="jaeger",pod="jaeger-hotrod-85b766c49d-s6vz7",route="jaeger-hotrod-qq5cq",server="10.134.249.205:8080",service="jaeger-hotrod"} 0
haproxy_server_http_average_queue_latency_milliseconds{namespace="jaeger",pod="jaeger-query-65d9446fc6-jx5qk",route="jaeger-query-pc6w5",server="10.134.248.54:16686",service="jaeger-query"} 0
haproxy_server_http_average_queue_latency_milliseconds{namespace="kafka-service",pod="akhq-58744f76db-h2mtf",route="akhq-ingress-lhz5z",server="10.134.249.133:8080",service="akhq"} 0
haproxy_server_http_average_queue_latency_milliseconds{namespace="license-server",pod="license-server-7b8b6f765f-wvcx5",route="license-server-9txj9",server="10.132.148.47:8080",service="license-server"} 0
haproxy_server_http_average_queue_latency_milliseconds{namespace="opensearch",pod="opensearch-dashboards-764bc8548b-zngcr",route="opensearch-dashboards-l6vdv",server="10.134.173.229:5601",service="opensearch-dashboards"} 0
haproxy_server_http_average_queue_latency_milliseconds{namespace="openshift-authentication",pod="oauth-openshift-75965c7487-kc8m7",route="oauth-openshift",server="10.132.153.132:6443",service="oauth-openshift"} 0
haproxy_server_http_average_queue_latency_milliseconds{namespace="openshift-authentication",pod="oauth-openshift-75965c7487-nsfxw",route="oauth-openshift",server="10.133.43.151:6443",service="oauth-openshift"} 0
haproxy_server_http_average_queue_latency_milliseconds{namespace="openshift-authentication",pod="oauth-openshift-75965c7487-xmmmq",route="oauth-openshift",server="10.135.207.187:6443",service="oauth-openshift"} 0
haproxy_server_http_average_queue_latency_milliseconds{namespace="openshift-console",pod="console-66fcc4f4fc-crbz6",route="console",server="10.133.43.162:8443",service="console"} 0
haproxy_server_http_average_queue_latency_milliseconds{namespace="openshift-console",pod="console-66fcc4f4fc-j5dvz",route="console",server="10.132.153.185:8443",service="console"} 0
haproxy_server_http_average_queue_latency_milliseconds{namespace="openshift-console",pod="downloads-6fb844fb47-cwnvw",route="downloads",server="10.135.207.221:8080",service="downloads"} 0
haproxy_server_http_average_queue_latency_milliseconds{namespace="openshift-console",pod="downloads-6fb844fb47-l485s",route="downloads",server="10.132.153.222:8080",service="downloads"} 0
haproxy_server_http_average_queue_latency_milliseconds{namespace="openshift-ingress-canary",pod="ingress-canary-7rvzj",route="canary",server="10.134.173.247:8080",service="ingress-canary"} 0
haproxy_server_http_average_queue_latency_milliseconds{namespace="openshift-ingress-canary",pod="ingress-canary-lwpkn",route="canary",server="10.132.149.137:8080",service="ingress-canary"} 0
haproxy_server_http_average_queue_latency_milliseconds{namespace="openshift-ingress-canary",pod="ingress-canary-z2b4x",route="canary",server="10.134.248.101:8080",service="ingress-canary"} 0
haproxy_server_http_average_queue_latency_milliseconds{namespace="openshift-monitoring",pod="alertmanager-main-0",route="alertmanager-main",server="10.134.173.99:9095",service="alertmanager-main"} 0
haproxy_server_http_average_queue_latency_milliseconds{namespace="openshift-monitoring",pod="alertmanager-main-1",route="alertmanager-main",server="10.134.248.127:9095",service="alertmanager-main"} 0
haproxy_server_http_average_queue_latency_milliseconds{namespace="openshift-monitoring",pod="prometheus-k8s-0",route="prometheus-k8s",server="10.134.173.51:9091",service="prometheus-k8s"} 0
haproxy_server_http_average_queue_latency_milliseconds{namespace="openshift-monitoring",pod="prometheus-k8s-0",route="prometheus-k8s-federate",server="10.134.173.51:9091",service="prometheus-k8s"} 0
haproxy_server_http_average_queue_latency_milliseconds{namespace="openshift-monitoring",pod="prometheus-k8s-1",route="prometheus-k8s",server="10.134.249.153:9091",service="prometheus-k8s"} 0
haproxy_server_http_average_queue_latency_milliseconds{namespace="openshift-monitoring",pod="prometheus-k8s-1",route="prometheus-k8s-federate",server="10.134.249.153:9091",service="prometheus-k8s"} 0
haproxy_server_http_average_queue_latency_milliseconds{namespace="openshift-monitoring",pod="thanos-querier-59d7db9787-hwjzt",route="thanos-querier",server="10.134.249.65:9091",service="thanos-querier"} 0
haproxy_server_http_average_queue_latency_milliseconds{namespace="openshift-monitoring",pod="thanos-querier-59d7db9787-p8m9l",route="thanos-querier",server="10.134.173.112:9091",service="thanos-querier"} 0
haproxy_server_http_average_queue_latency_milliseconds{namespace="profiler",pod="esc-static-service-654f8bc589-xchjk",route="esc-static-service-j9sgv",server="10.134.173.123:8080",service="esc-static-service"} 0
haproxy_server_http_average_queue_latency_milliseconds{namespace="profiler",pod="esc-test-service-75849d7c6c-czvsw",route="esc-test-service-lbp8r",server="10.134.172.182:8080",service="esc-test-service"} 0
haproxy_server_http_average_queue_latency_milliseconds{namespace="profiler",pod="esc-ui-service-66977c9bcb-kg9mc",route="esc-ui-service-r5tfm",server="10.134.249.215:8180",service="esc-ui-service"} 0
haproxy_server_http_average_queue_latency_milliseconds{namespace="prometheus-operator",pod="grafana-deployment-88d46bc5f-5gx7h",route="prometheus-operator-grafana",server="10.132.149.3:3000",service="grafana-service"} 0
haproxy_server_http_average_queue_latency_milliseconds{namespace="prometheus-operator",pod="vmagent-k8s-67654d6f97-98nf5",route="prometheus-operator-vmagent-k8s",server="10.132.149.156:8429",service="vmagent-k8s"} 0
haproxy_server_http_average_queue_latency_milliseconds{namespace="prometheus-operator",pod="vmalert-k8s-6bcf966dd6-jvpz5",route="prometheus-operator-vmalert-k8s",server="10.132.149.175:8080",service="vmalert-k8s"} 0
haproxy_server_http_average_queue_latency_milliseconds{namespace="prometheus-operator",pod="vmalertmanager-k8s-0",route="prometheus-operator-vmalertmanager-k8s",server="10.134.173.133:9093",service="vmalertmanager-k8s"} 0
haproxy_server_http_average_queue_latency_milliseconds{namespace="prometheus-operator",pod="vmsingle-k8s-f4844bfc4-5l8sf",route="prometheus-operator-vmsingle-k8s",server="10.134.249.110:8429",service="vmsingle-k8s"} 0
haproxy_server_http_average_queue_latency_milliseconds{namespace="rabbitmq",pod="rmqlocal-0",route="rabbitmq-fnzhp",server="10.134.249.141:15672",service="rabbitmq"} 0
haproxy_server_http_average_queue_latency_milliseconds{namespace="rabbitmq",pod="rmqlocal-1",route="rabbitmq-fnzhp",server="10.134.173.127:15672",service="rabbitmq"} 0
haproxy_server_http_average_queue_latency_milliseconds{namespace="rabbitmq",pod="rmqlocal-2",route="rabbitmq-fnzhp",server="10.132.149.46:15672",service="rabbitmq"} 0
haproxy_server_http_average_queue_latency_milliseconds{namespace="shared-license-distributor",pod="shared-license-distributor-64ccc89449-8kw8p",route="shared-license-distributor-wnj68",server="10.134.249.239:8080",service="shared-license-distributor"} 0
haproxy_server_http_average_queue_latency_milliseconds{namespace="site-manager",pod="site-manager-547546bdd5-hd85n",route="site-manager",server="10.134.248.81:8443",service="site-manager"} 0
haproxy_server_http_average_queue_latency_milliseconds{namespace="spark",pod="spark-history-server-bd755d9f5-9ccsz",route="spark-history-server-hqshb",server="10.132.149.178:18080",service="spark-history-server"} 0
haproxy_server_http_average_queue_latency_milliseconds{namespace="spark",pod="spark-history-server-bd755d9f5-zf2dt",route="spark-history-server-hqshb",server="10.134.173.147:18080",service="spark-history-server"} 0
haproxy_server_http_average_queue_latency_milliseconds{namespace="spark",pod="spark-site-manager-66c44c7b78-cc7jr",route="spark-site-manager-p4nqh",server="10.132.149.148:8080",service="spark-site-manager"} 0
# HELP haproxy_server_http_average_response_latency_milliseconds Average response latency of the last 1024 requests in milliseconds.
# TYPE haproxy_server_http_average_response_latency_milliseconds gauge
haproxy_server_http_average_response_latency_milliseconds{namespace="",pod="",route="",server="fe_no_sni",service=""} 0
haproxy_server_http_average_response_latency_milliseconds{namespace="",pod="",route="",server="fe_sni",service=""} 0
haproxy_server_http_average_response_latency_milliseconds{namespace="airflow",pod="airflow-site-manager-68ccf87889-bwbdf",route="airflow-site-manager-hpl66",server="10.132.149.80:8080",service="airflow-site-manager"} 0
haproxy_server_http_average_response_latency_milliseconds{namespace="airflow",pod="airflow-webserver-695b5c545b-rxqc4",route="airflow-ingress-nzbtm",server="10.134.173.137:8080",service="airflow-webserver"} 0
haproxy_server_http_average_response_latency_milliseconds{namespace="arangodb",pod="main-arangodb-1-87669b688-q98z4",route="arangodb-leader-ingress-zw82p",server="10.132.149.185:8529",service="arangodb-leader"} 0
haproxy_server_http_average_response_latency_milliseconds{namespace="consul-service",pod="consul-server-0",route="consul-ingress-hmmvf",server="10.134.248.50:8500",service="consul-ui"} 0
haproxy_server_http_average_response_latency_milliseconds{namespace="consul-service",pod="consul-server-1",route="consul-ingress-hmmvf",server="10.132.149.160:8500",service="consul-ui"} 0
haproxy_server_http_average_response_latency_milliseconds{namespace="consul-service",pod="consul-server-2",route="consul-ingress-hmmvf",server="10.134.173.136:8500",service="consul-ui"} 0
haproxy_server_http_average_response_latency_milliseconds{namespace="dbaas",pod="dbaas-aggregator-5f744d74f5-r7tqm",route="aggregator-fw5ft",server="10.134.248.185:8080",service="dbaas-aggregator"} 0
haproxy_server_http_average_response_latency_milliseconds{namespace="infra-keycloak",pod="infra-keycloak-7fd57cd59c-d6c5b",route="infra-keycloak-5q4rj",server="10.132.148.223:8080",service="infra-keycloak"} 0
haproxy_server_http_average_response_latency_milliseconds{namespace="infra-keycloak",pod="infra-keycloak-7fd57cd59c-qtqv9",route="infra-keycloak-5q4rj",server="10.134.173.182:8080",service="infra-keycloak"} 0
haproxy_server_http_average_response_latency_milliseconds{namespace="jaeger",pod="jaeger-hotrod-85b766c49d-s6vz7",route="jaeger-hotrod-qq5cq",server="10.134.249.205:8080",service="jaeger-hotrod"} 0
haproxy_server_http_average_response_latency_milliseconds{namespace="jaeger",pod="jaeger-query-65d9446fc6-jx5qk",route="jaeger-query-pc6w5",server="10.134.248.54:16686",service="jaeger-query"} 0
haproxy_server_http_average_response_latency_milliseconds{namespace="kafka-service",pod="akhq-58744f76db-h2mtf",route="akhq-ingress-lhz5z",server="10.134.249.133:8080",service="akhq"} 0
haproxy_server_http_average_response_latency_milliseconds{namespace="license-server",pod="license-server-7b8b6f765f-wvcx5",route="license-server-9txj9",server="10.132.148.47:8080",service="license-server"} 0
haproxy_server_http_average_response_latency_milliseconds{namespace="opensearch",pod="opensearch-dashboards-764bc8548b-zngcr",route="opensearch-dashboards-l6vdv",server="10.134.173.229:5601",service="opensearch-dashboards"} 0
haproxy_server_http_average_response_latency_milliseconds{namespace="openshift-authentication",pod="oauth-openshift-75965c7487-kc8m7",route="oauth-openshift",server="10.132.153.132:6443",service="oauth-openshift"} 0
haproxy_server_http_average_response_latency_milliseconds{namespace="openshift-authentication",pod="oauth-openshift-75965c7487-nsfxw",route="oauth-openshift",server="10.133.43.151:6443",service="oauth-openshift"} 0
haproxy_server_http_average_response_latency_milliseconds{namespace="openshift-authentication",pod="oauth-openshift-75965c7487-xmmmq",route="oauth-openshift",server="10.135.207.187:6443",service="oauth-openshift"} 0
haproxy_server_http_average_response_latency_milliseconds{namespace="openshift-console",pod="console-66fcc4f4fc-crbz6",route="console",server="10.133.43.162:8443",service="console"} 0
haproxy_server_http_average_response_latency_milliseconds{namespace="openshift-console",pod="console-66fcc4f4fc-j5dvz",route="console",server="10.132.153.185:8443",service="console"} 0
haproxy_server_http_average_response_latency_milliseconds{namespace="openshift-console",pod="downloads-6fb844fb47-cwnvw",route="downloads",server="10.135.207.221:8080",service="downloads"} 0
haproxy_server_http_average_response_latency_milliseconds{namespace="openshift-console",pod="downloads-6fb844fb47-l485s",route="downloads",server="10.132.153.222:8080",service="downloads"} 0
haproxy_server_http_average_response_latency_milliseconds{namespace="openshift-ingress-canary",pod="ingress-canary-7rvzj",route="canary",server="10.134.173.247:8080",service="ingress-canary"} 0
haproxy_server_http_average_response_latency_milliseconds{namespace="openshift-ingress-canary",pod="ingress-canary-lwpkn",route="canary",server="10.132.149.137:8080",service="ingress-canary"} 0
haproxy_server_http_average_response_latency_milliseconds{namespace="openshift-ingress-canary",pod="ingress-canary-z2b4x",route="canary",server="10.134.248.101:8080",service="ingress-canary"} 0
haproxy_server_http_average_response_latency_milliseconds{namespace="openshift-monitoring",pod="alertmanager-main-0",route="alertmanager-main",server="10.134.173.99:9095",service="alertmanager-main"} 0
haproxy_server_http_average_response_latency_milliseconds{namespace="openshift-monitoring",pod="alertmanager-main-1",route="alertmanager-main",server="10.134.248.127:9095",service="alertmanager-main"} 0
haproxy_server_http_average_response_latency_milliseconds{namespace="openshift-monitoring",pod="prometheus-k8s-0",route="prometheus-k8s",server="10.134.173.51:9091",service="prometheus-k8s"} 0
haproxy_server_http_average_response_latency_milliseconds{namespace="openshift-monitoring",pod="prometheus-k8s-0",route="prometheus-k8s-federate",server="10.134.173.51:9091",service="prometheus-k8s"} 0
haproxy_server_http_average_response_latency_milliseconds{namespace="openshift-monitoring",pod="prometheus-k8s-1",route="prometheus-k8s",server="10.134.249.153:9091",service="prometheus-k8s"} 0
haproxy_server_http_average_response_latency_milliseconds{namespace="openshift-monitoring",pod="prometheus-k8s-1",route="prometheus-k8s-federate",server="10.134.249.153:9091",service="prometheus-k8s"} 0
haproxy_server_http_average_response_latency_milliseconds{namespace="openshift-monitoring",pod="thanos-querier-59d7db9787-hwjzt",route="thanos-querier",server="10.134.249.65:9091",service="thanos-querier"} 0
haproxy_server_http_average_response_latency_milliseconds{namespace="openshift-monitoring",pod="thanos-querier-59d7db9787-p8m9l",route="thanos-querier",server="10.134.173.112:9091",service="thanos-querier"} 0
haproxy_server_http_average_response_latency_milliseconds{namespace="profiler",pod="esc-static-service-654f8bc589-xchjk",route="esc-static-service-j9sgv",server="10.134.173.123:8080",service="esc-static-service"} 0
haproxy_server_http_average_response_latency_milliseconds{namespace="profiler",pod="esc-test-service-75849d7c6c-czvsw",route="esc-test-service-lbp8r",server="10.134.172.182:8080",service="esc-test-service"} 0
haproxy_server_http_average_response_latency_milliseconds{namespace="profiler",pod="esc-ui-service-66977c9bcb-kg9mc",route="esc-ui-service-r5tfm",server="10.134.249.215:8180",service="esc-ui-service"} 0
haproxy_server_http_average_response_latency_milliseconds{namespace="prometheus-operator",pod="grafana-deployment-88d46bc5f-5gx7h",route="prometheus-operator-grafana",server="10.132.149.3:3000",service="grafana-service"} 0
haproxy_server_http_average_response_latency_milliseconds{namespace="prometheus-operator",pod="vmagent-k8s-67654d6f97-98nf5",route="prometheus-operator-vmagent-k8s",server="10.132.149.156:8429",service="vmagent-k8s"} 0
haproxy_server_http_average_response_latency_milliseconds{namespace="prometheus-operator",pod="vmalert-k8s-6bcf966dd6-jvpz5",route="prometheus-operator-vmalert-k8s",server="10.132.149.175:8080",service="vmalert-k8s"} 0
haproxy_server_http_average_response_latency_milliseconds{namespace="prometheus-operator",pod="vmalertmanager-k8s-0",route="prometheus-operator-vmalertmanager-k8s",server="10.134.173.133:9093",service="vmalertmanager-k8s"} 0
haproxy_server_http_average_response_latency_milliseconds{namespace="prometheus-operator",pod="vmsingle-k8s-f4844bfc4-5l8sf",route="prometheus-operator-vmsingle-k8s",server="10.134.249.110:8429",service="vmsingle-k8s"} 0
haproxy_server_http_average_response_latency_milliseconds{namespace="rabbitmq",pod="rmqlocal-0",route="rabbitmq-fnzhp",server="10.134.249.141:15672",service="rabbitmq"} 0
haproxy_server_http_average_response_latency_milliseconds{namespace="rabbitmq",pod="rmqlocal-1",route="rabbitmq-fnzhp",server="10.134.173.127:15672",service="rabbitmq"} 0
haproxy_server_http_average_response_latency_milliseconds{namespace="rabbitmq",pod="rmqlocal-2",route="rabbitmq-fnzhp",server="10.132.149.46:15672",service="rabbitmq"} 0
haproxy_server_http_average_response_latency_milliseconds{namespace="shared-license-distributor",pod="shared-license-distributor-64ccc89449-8kw8p",route="shared-license-distributor-wnj68",server="10.134.249.239:8080",service="shared-license-distributor"} 0
haproxy_server_http_average_response_latency_milliseconds{namespace="site-manager",pod="site-manager-547546bdd5-hd85n",route="site-manager",server="10.134.248.81:8443",service="site-manager"} 0
haproxy_server_http_average_response_latency_milliseconds{namespace="spark",pod="spark-history-server-bd755d9f5-9ccsz",route="spark-history-server-hqshb",server="10.132.149.178:18080",service="spark-history-server"} 0
haproxy_server_http_average_response_latency_milliseconds{namespace="spark",pod="spark-history-server-bd755d9f5-zf2dt",route="spark-history-server-hqshb",server="10.134.173.147:18080",service="spark-history-server"} 0
haproxy_server_http_average_response_latency_milliseconds{namespace="spark",pod="spark-site-manager-66c44c7b78-cc7jr",route="spark-site-manager-p4nqh",server="10.132.149.148:8080",service="spark-site-manager"} 0
# HELP haproxy_server_http_responses_total Total of HTTP responses.
# TYPE haproxy_server_http_responses_total gauge
haproxy_server_http_responses_total{code="1xx",namespace="airflow",pod="airflow-site-manager-68ccf87889-bwbdf",route="airflow-site-manager-hpl66",server="10.132.149.80:8080",service="airflow-site-manager"} 0
haproxy_server_http_responses_total{code="1xx",namespace="airflow",pod="airflow-webserver-695b5c545b-rxqc4",route="airflow-ingress-nzbtm",server="10.134.173.137:8080",service="airflow-webserver"} 0
haproxy_server_http_responses_total{code="1xx",namespace="arangodb",pod="main-arangodb-1-87669b688-q98z4",route="arangodb-leader-ingress-zw82p",server="10.132.149.185:8529",service="arangodb-leader"} 0
haproxy_server_http_responses_total{code="1xx",namespace="consul-service",pod="consul-server-0",route="consul-ingress-hmmvf",server="10.134.248.50:8500",service="consul-ui"} 0
haproxy_server_http_responses_total{code="1xx",namespace="consul-service",pod="consul-server-1",route="consul-ingress-hmmvf",server="10.132.149.160:8500",service="consul-ui"} 0
haproxy_server_http_responses_total{code="1xx",namespace="consul-service",pod="consul-server-2",route="consul-ingress-hmmvf",server="10.134.173.136:8500",service="consul-ui"} 0
haproxy_server_http_responses_total{code="1xx",namespace="dbaas",pod="dbaas-aggregator-5f744d74f5-r7tqm",route="aggregator-fw5ft",server="10.134.248.185:8080",service="dbaas-aggregator"} 0
haproxy_server_http_responses_total{code="1xx",namespace="infra-keycloak",pod="infra-keycloak-7fd57cd59c-d6c5b",route="infra-keycloak-5q4rj",server="10.132.148.223:8080",service="infra-keycloak"} 0
haproxy_server_http_responses_total{code="1xx",namespace="infra-keycloak",pod="infra-keycloak-7fd57cd59c-qtqv9",route="infra-keycloak-5q4rj",server="10.134.173.182:8080",service="infra-keycloak"} 0
haproxy_server_http_responses_total{code="1xx",namespace="jaeger",pod="jaeger-hotrod-85b766c49d-s6vz7",route="jaeger-hotrod-qq5cq",server="10.134.249.205:8080",service="jaeger-hotrod"} 0
haproxy_server_http_responses_total{code="1xx",namespace="jaeger",pod="jaeger-query-65d9446fc6-jx5qk",route="jaeger-query-pc6w5",server="10.134.248.54:16686",service="jaeger-query"} 0
haproxy_server_http_responses_total{code="1xx",namespace="kafka-service",pod="akhq-58744f76db-h2mtf",route="akhq-ingress-lhz5z",server="10.134.249.133:8080",service="akhq"} 0
haproxy_server_http_responses_total{code="1xx",namespace="license-server",pod="license-server-7b8b6f765f-wvcx5",route="license-server-9txj9",server="10.132.148.47:8080",service="license-server"} 0
haproxy_server_http_responses_total{code="1xx",namespace="opensearch",pod="opensearch-dashboards-764bc8548b-zngcr",route="opensearch-dashboards-l6vdv",server="10.134.173.229:5601",service="opensearch-dashboards"} 0
haproxy_server_http_responses_total{code="1xx",namespace="openshift-console",pod="console-66fcc4f4fc-crbz6",route="console",server="10.133.43.162:8443",service="console"} 339
haproxy_server_http_responses_total{code="1xx",namespace="openshift-console",pod="console-66fcc4f4fc-j5dvz",route="console",server="10.132.153.185:8443",service="console"} 368
haproxy_server_http_responses_total{code="1xx",namespace="openshift-console",pod="downloads-6fb844fb47-cwnvw",route="downloads",server="10.135.207.221:8080",service="downloads"} 0
haproxy_server_http_responses_total{code="1xx",namespace="openshift-console",pod="downloads-6fb844fb47-l485s",route="downloads",server="10.132.153.222:8080",service="downloads"} 0
haproxy_server_http_responses_total{code="1xx",namespace="openshift-ingress-canary",pod="ingress-canary-7rvzj",route="canary",server="10.134.173.247:8080",service="ingress-canary"} 0
haproxy_server_http_responses_total{code="1xx",namespace="openshift-ingress-canary",pod="ingress-canary-lwpkn",route="canary",server="10.132.149.137:8080",service="ingress-canary"} 0
haproxy_server_http_responses_total{code="1xx",namespace="openshift-ingress-canary",pod="ingress-canary-z2b4x",route="canary",server="10.134.248.101:8080",service="ingress-canary"} 0
haproxy_server_http_responses_total{code="1xx",namespace="openshift-monitoring",pod="alertmanager-main-0",route="alertmanager-main",server="10.134.173.99:9095",service="alertmanager-main"} 0
haproxy_server_http_responses_total{code="1xx",namespace="openshift-monitoring",pod="alertmanager-main-1",route="alertmanager-main",server="10.134.248.127:9095",service="alertmanager-main"} 0
haproxy_server_http_responses_total{code="1xx",namespace="openshift-monitoring",pod="prometheus-k8s-0",route="prometheus-k8s",server="10.134.173.51:9091",service="prometheus-k8s"} 0
haproxy_server_http_responses_total{code="1xx",namespace="openshift-monitoring",pod="prometheus-k8s-0",route="prometheus-k8s-federate",server="10.134.173.51:9091",service="prometheus-k8s"} 0
haproxy_server_http_responses_total{code="1xx",namespace="openshift-monitoring",pod="prometheus-k8s-1",route="prometheus-k8s",server="10.134.249.153:9091",service="prometheus-k8s"} 0
haproxy_server_http_responses_total{code="1xx",namespace="openshift-monitoring",pod="prometheus-k8s-1",route="prometheus-k8s-federate",server="10.134.249.153:9091",service="prometheus-k8s"} 0
haproxy_server_http_responses_total{code="1xx",namespace="openshift-monitoring",pod="thanos-querier-59d7db9787-hwjzt",route="thanos-querier",server="10.134.249.65:9091",service="thanos-querier"} 0
haproxy_server_http_responses_total{code="1xx",namespace="openshift-monitoring",pod="thanos-querier-59d7db9787-p8m9l",route="thanos-querier",server="10.134.173.112:9091",service="thanos-querier"} 0
haproxy_server_http_responses_total{code="1xx",namespace="profiler",pod="esc-static-service-654f8bc589-xchjk",route="esc-static-service-j9sgv",server="10.134.173.123:8080",service="esc-static-service"} 0
haproxy_server_http_responses_total{code="1xx",namespace="profiler",pod="esc-test-service-75849d7c6c-czvsw",route="esc-test-service-lbp8r",server="10.134.172.182:8080",service="esc-test-service"} 0
haproxy_server_http_responses_total{code="1xx",namespace="profiler",pod="esc-ui-service-66977c9bcb-kg9mc",route="esc-ui-service-r5tfm",server="10.134.249.215:8180",service="esc-ui-service"} 0
haproxy_server_http_responses_total{code="1xx",namespace="prometheus-operator",pod="grafana-deployment-88d46bc5f-5gx7h",route="prometheus-operator-grafana",server="10.132.149.3:3000",service="grafana-service"} 0
haproxy_server_http_responses_total{code="1xx",namespace="prometheus-operator",pod="vmagent-k8s-67654d6f97-98nf5",route="prometheus-operator-vmagent-k8s",server="10.132.149.156:8429",service="vmagent-k8s"} 0
haproxy_server_http_responses_total{code="1xx",namespace="prometheus-operator",pod="vmalert-k8s-6bcf966dd6-jvpz5",route="prometheus-operator-vmalert-k8s",server="10.132.149.175:8080",service="vmalert-k8s"} 0
haproxy_server_http_responses_total{code="1xx",namespace="prometheus-operator",pod="vmalertmanager-k8s-0",route="prometheus-operator-vmalertmanager-k8s",server="10.134.173.133:9093",service="vmalertmanager-k8s"} 0
haproxy_server_http_responses_total{code="1xx",namespace="prometheus-operator",pod="vmsingle-k8s-f4844bfc4-5l8sf",route="prometheus-operator-vmsingle-k8s",server="10.134.249.110:8429",service="vmsingle-k8s"} 0
haproxy_server_http_responses_total{code="1xx",namespace="rabbitmq",pod="rmqlocal-0",route="rabbitmq-fnzhp",server="10.134.249.141:15672",service="rabbitmq"} 0
haproxy_server_http_responses_total{code="1xx",namespace="rabbitmq",pod="rmqlocal-1",route="rabbitmq-fnzhp",server="10.134.173.127:15672",service="rabbitmq"} 0
haproxy_server_http_responses_total{code="1xx",namespace="rabbitmq",pod="rmqlocal-2",route="rabbitmq-fnzhp",server="10.132.149.46:15672",service="rabbitmq"} 0
haproxy_server_http_responses_total{code="1xx",namespace="shared-license-distributor",pod="shared-license-distributor-64ccc89449-8kw8p",route="shared-license-distributor-wnj68",server="10.134.249.239:8080",service="shared-license-distributor"} 0
haproxy_server_http_responses_total{code="1xx",namespace="spark",pod="spark-history-server-bd755d9f5-9ccsz",route="spark-history-server-hqshb",server="10.132.149.178:18080",service="spark-history-server"} 0
haproxy_server_http_responses_total{code="1xx",namespace="spark",pod="spark-history-server-bd755d9f5-zf2dt",route="spark-history-server-hqshb",server="10.134.173.147:18080",service="spark-history-server"} 0
haproxy_server_http_responses_total{code="1xx",namespace="spark",pod="spark-site-manager-66c44c7b78-cc7jr",route="spark-site-manager-p4nqh",server="10.132.149.148:8080",service="spark-site-manager"} 0
haproxy_server_http_responses_total{code="2xx",namespace="airflow",pod="airflow-site-manager-68ccf87889-bwbdf",route="airflow-site-manager-hpl66",server="10.132.149.80:8080",service="airflow-site-manager"} 0
haproxy_server_http_responses_total{code="2xx",namespace="airflow",pod="airflow-webserver-695b5c545b-rxqc4",route="airflow-ingress-nzbtm",server="10.134.173.137:8080",service="airflow-webserver"} 0
haproxy_server_http_responses_total{code="2xx",namespace="arangodb",pod="main-arangodb-1-87669b688-q98z4",route="arangodb-leader-ingress-zw82p",server="10.132.149.185:8529",service="arangodb-leader"} 0
haproxy_server_http_responses_total{code="2xx",namespace="consul-service",pod="consul-server-0",route="consul-ingress-hmmvf",server="10.134.248.50:8500",service="consul-ui"} 0
haproxy_server_http_responses_total{code="2xx",namespace="consul-service",pod="consul-server-1",route="consul-ingress-hmmvf",server="10.132.149.160:8500",service="consul-ui"} 0
haproxy_server_http_responses_total{code="2xx",namespace="consul-service",pod="consul-server-2",route="consul-ingress-hmmvf",server="10.134.173.136:8500",service="consul-ui"} 0
haproxy_server_http_responses_total{code="2xx",namespace="dbaas",pod="dbaas-aggregator-5f744d74f5-r7tqm",route="aggregator-fw5ft",server="10.134.248.185:8080",service="dbaas-aggregator"} 0
haproxy_server_http_responses_total{code="2xx",namespace="infra-keycloak",pod="infra-keycloak-7fd57cd59c-d6c5b",route="infra-keycloak-5q4rj",server="10.132.148.223:8080",service="infra-keycloak"} 0
haproxy_server_http_responses_total{code="2xx",namespace="infra-keycloak",pod="infra-keycloak-7fd57cd59c-qtqv9",route="infra-keycloak-5q4rj",server="10.134.173.182:8080",service="infra-keycloak"} 0
haproxy_server_http_responses_total{code="2xx",namespace="jaeger",pod="jaeger-hotrod-85b766c49d-s6vz7",route="jaeger-hotrod-qq5cq",server="10.134.249.205:8080",service="jaeger-hotrod"} 0
haproxy_server_http_responses_total{code="2xx",namespace="jaeger",pod="jaeger-query-65d9446fc6-jx5qk",route="jaeger-query-pc6w5",server="10.134.248.54:16686",service="jaeger-query"} 0
haproxy_server_http_responses_total{code="2xx",namespace="kafka-service",pod="akhq-58744f76db-h2mtf",route="akhq-ingress-lhz5z",server="10.134.249.133:8080",service="akhq"} 20
haproxy_server_http_responses_total{code="2xx",namespace="license-server",pod="license-server-7b8b6f765f-wvcx5",route="license-server-9txj9",server="10.132.148.47:8080",service="license-server"} 0
haproxy_server_http_responses_total{code="2xx",namespace="opensearch",pod="opensearch-dashboards-764bc8548b-zngcr",route="opensearch-dashboards-l6vdv",server="10.134.173.229:5601",service="opensearch-dashboards"} 0
haproxy_server_http_responses_total{code="2xx",namespace="openshift-console",pod="console-66fcc4f4fc-crbz6",route="console",server="10.133.43.162:8443",service="console"} 27030
haproxy_server_http_responses_total{code="2xx",namespace="openshift-console",pod="console-66fcc4f4fc-j5dvz",route="console",server="10.132.153.185:8443",service="console"} 32888
haproxy_server_http_responses_total{code="2xx",namespace="openshift-console",pod="downloads-6fb844fb47-cwnvw",route="downloads",server="10.135.207.221:8080",service="downloads"} 0
haproxy_server_http_responses_total{code="2xx",namespace="openshift-console",pod="downloads-6fb844fb47-l485s",route="downloads",server="10.132.153.222:8080",service="downloads"} 0
haproxy_server_http_responses_total{code="2xx",namespace="openshift-ingress-canary",pod="ingress-canary-7rvzj",route="canary",server="10.134.173.247:8080",service="ingress-canary"} 0
haproxy_server_http_responses_total{code="2xx",namespace="openshift-ingress-canary",pod="ingress-canary-lwpkn",route="canary",server="10.132.149.137:8080",service="ingress-canary"} 0
haproxy_server_http_responses_total{code="2xx",namespace="openshift-ingress-canary",pod="ingress-canary-z2b4x",route="canary",server="10.134.248.101:8080",service="ingress-canary"} 0
haproxy_server_http_responses_total{code="2xx",namespace="openshift-monitoring",pod="alertmanager-main-0",route="alertmanager-main",server="10.134.173.99:9095",service="alertmanager-main"} 0
haproxy_server_http_responses_total{code="2xx",namespace="openshift-monitoring",pod="alertmanager-main-1",route="alertmanager-main",server="10.134.248.127:9095",service="alertmanager-main"} 0
haproxy_server_http_responses_total{code="2xx",namespace="openshift-monitoring",pod="prometheus-k8s-0",route="prometheus-k8s",server="10.134.173.51:9091",service="prometheus-k8s"} 0
haproxy_server_http_responses_total{code="2xx",namespace="openshift-monitoring",pod="prometheus-k8s-0",route="prometheus-k8s-federate",server="10.134.173.51:9091",service="prometheus-k8s"} 0
haproxy_server_http_responses_total{code="2xx",namespace="openshift-monitoring",pod="prometheus-k8s-1",route="prometheus-k8s",server="10.134.249.153:9091",service="prometheus-k8s"} 0
haproxy_server_http_responses_total{code="2xx",namespace="openshift-monitoring",pod="prometheus-k8s-1",route="prometheus-k8s-federate",server="10.134.249.153:9091",service="prometheus-k8s"} 0
haproxy_server_http_responses_total{code="2xx",namespace="openshift-monitoring",pod="thanos-querier-59d7db9787-hwjzt",route="thanos-querier",server="10.134.249.65:9091",service="thanos-querier"} 0
haproxy_server_http_responses_total{code="2xx",namespace="openshift-monitoring",pod="thanos-querier-59d7db9787-p8m9l",route="thanos-querier",server="10.134.173.112:9091",service="thanos-querier"} 0
haproxy_server_http_responses_total{code="2xx",namespace="profiler",pod="esc-static-service-654f8bc589-xchjk",route="esc-static-service-j9sgv",server="10.134.173.123:8080",service="esc-static-service"} 0
haproxy_server_http_responses_total{code="2xx",namespace="profiler",pod="esc-test-service-75849d7c6c-czvsw",route="esc-test-service-lbp8r",server="10.134.172.182:8080",service="esc-test-service"} 0
haproxy_server_http_responses_total{code="2xx",namespace="profiler",pod="esc-ui-service-66977c9bcb-kg9mc",route="esc-ui-service-r5tfm",server="10.134.249.215:8180",service="esc-ui-service"} 0
haproxy_server_http_responses_total{code="2xx",namespace="prometheus-operator",pod="grafana-deployment-88d46bc5f-5gx7h",route="prometheus-operator-grafana",server="10.132.149.3:3000",service="grafana-service"} 0
haproxy_server_http_responses_total{code="2xx",namespace="prometheus-operator",pod="vmagent-k8s-67654d6f97-98nf5",route="prometheus-operator-vmagent-k8s",server="10.132.149.156:8429",service="vmagent-k8s"} 0
haproxy_server_http_responses_total{code="2xx",namespace="prometheus-operator",pod="vmalert-k8s-6bcf966dd6-jvpz5",route="prometheus-operator-vmalert-k8s",server="10.132.149.175:8080",service="vmalert-k8s"} 0
haproxy_server_http_responses_total{code="2xx",namespace="prometheus-operator",pod="vmalertmanager-k8s-0",route="prometheus-operator-vmalertmanager-k8s",server="10.134.173.133:9093",service="vmalertmanager-k8s"} 0
haproxy_server_http_responses_total{code="2xx",namespace="prometheus-operator",pod="vmsingle-k8s-f4844bfc4-5l8sf",route="prometheus-operator-vmsingle-k8s",server="10.134.249.110:8429",service="vmsingle-k8s"} 36
haproxy_server_http_responses_total{code="2xx",namespace="rabbitmq",pod="rmqlocal-0",route="rabbitmq-fnzhp",server="10.134.249.141:15672",service="rabbitmq"} 0
haproxy_server_http_responses_total{code="2xx",namespace="rabbitmq",pod="rmqlocal-1",route="rabbitmq-fnzhp",server="10.134.173.127:15672",service="rabbitmq"} 0
haproxy_server_http_responses_total{code="2xx",namespace="rabbitmq",pod="rmqlocal-2",route="rabbitmq-fnzhp",server="10.132.149.46:15672",service="rabbitmq"} 0
haproxy_server_http_responses_total{code="2xx",namespace="shared-license-distributor",pod="shared-license-distributor-64ccc89449-8kw8p",route="shared-license-distributor-wnj68",server="10.134.249.239:8080",service="shared-license-distributor"} 0
haproxy_server_http_responses_total{code="2xx",namespace="spark",pod="spark-history-server-bd755d9f5-9ccsz",route="spark-history-server-hqshb",server="10.132.149.178:18080",service="spark-history-server"} 0
haproxy_server_http_responses_total{code="2xx",namespace="spark",pod="spark-history-server-bd755d9f5-zf2dt",route="spark-history-server-hqshb",server="10.134.173.147:18080",service="spark-history-server"} 0
haproxy_server_http_responses_total{code="2xx",namespace="spark",pod="spark-site-manager-66c44c7b78-cc7jr",route="spark-site-manager-p4nqh",server="10.132.149.148:8080",service="spark-site-manager"} 0
haproxy_server_http_responses_total{code="3xx",namespace="airflow",pod="airflow-site-manager-68ccf87889-bwbdf",route="airflow-site-manager-hpl66",server="10.132.149.80:8080",service="airflow-site-manager"} 0
haproxy_server_http_responses_total{code="3xx",namespace="airflow",pod="airflow-webserver-695b5c545b-rxqc4",route="airflow-ingress-nzbtm",server="10.134.173.137:8080",service="airflow-webserver"} 0
haproxy_server_http_responses_total{code="3xx",namespace="arangodb",pod="main-arangodb-1-87669b688-q98z4",route="arangodb-leader-ingress-zw82p",server="10.132.149.185:8529",service="arangodb-leader"} 0
haproxy_server_http_responses_total{code="3xx",namespace="consul-service",pod="consul-server-0",route="consul-ingress-hmmvf",server="10.134.248.50:8500",service="consul-ui"} 0
haproxy_server_http_responses_total{code="3xx",namespace="consul-service",pod="consul-server-1",route="consul-ingress-hmmvf",server="10.132.149.160:8500",service="consul-ui"} 0
haproxy_server_http_responses_total{code="3xx",namespace="consul-service",pod="consul-server-2",route="consul-ingress-hmmvf",server="10.134.173.136:8500",service="consul-ui"} 0
haproxy_server_http_responses_total{code="3xx",namespace="dbaas",pod="dbaas-aggregator-5f744d74f5-r7tqm",route="aggregator-fw5ft",server="10.134.248.185:8080",service="dbaas-aggregator"} 0
haproxy_server_http_responses_total{code="3xx",namespace="infra-keycloak",pod="infra-keycloak-7fd57cd59c-d6c5b",route="infra-keycloak-5q4rj",server="10.132.148.223:8080",service="infra-keycloak"} 0
haproxy_server_http_responses_total{code="3xx",namespace="infra-keycloak",pod="infra-keycloak-7fd57cd59c-qtqv9",route="infra-keycloak-5q4rj",server="10.134.173.182:8080",service="infra-keycloak"} 0
haproxy_server_http_responses_total{code="3xx",namespace="jaeger",pod="jaeger-hotrod-85b766c49d-s6vz7",route="jaeger-hotrod-qq5cq",server="10.134.249.205:8080",service="jaeger-hotrod"} 0
haproxy_server_http_responses_total{code="3xx",namespace="jaeger",pod="jaeger-query-65d9446fc6-jx5qk",route="jaeger-query-pc6w5",server="10.134.248.54:16686",service="jaeger-query"} 0
haproxy_server_http_responses_total{code="3xx",namespace="kafka-service",pod="akhq-58744f76db-h2mtf",route="akhq-ingress-lhz5z",server="10.134.249.133:8080",service="akhq"} 2
haproxy_server_http_responses_total{code="3xx",namespace="license-server",pod="license-server-7b8b6f765f-wvcx5",route="license-server-9txj9",server="10.132.148.47:8080",service="license-server"} 0
haproxy_server_http_responses_total{code="3xx",namespace="opensearch",pod="opensearch-dashboards-764bc8548b-zngcr",route="opensearch-dashboards-l6vdv",server="10.134.173.229:5601",service="opensearch-dashboards"} 0
haproxy_server_http_responses_total{code="3xx",namespace="openshift-console",pod="console-66fcc4f4fc-crbz6",route="console",server="10.133.43.162:8443",service="console"} 9
haproxy_server_http_responses_total{code="3xx",namespace="openshift-console",pod="console-66fcc4f4fc-j5dvz",route="console",server="10.132.153.185:8443",service="console"} 6
haproxy_server_http_responses_total{code="3xx",namespace="openshift-console",pod="downloads-6fb844fb47-cwnvw",route="downloads",server="10.135.207.221:8080",service="downloads"} 0
haproxy_server_http_responses_total{code="3xx",namespace="openshift-console",pod="downloads-6fb844fb47-l485s",route="downloads",server="10.132.153.222:8080",service="downloads"} 0
haproxy_server_http_responses_total{code="3xx",namespace="openshift-ingress-canary",pod="ingress-canary-7rvzj",route="canary",server="10.134.173.247:8080",service="ingress-canary"} 0
haproxy_server_http_responses_total{code="3xx",namespace="openshift-ingress-canary",pod="ingress-canary-lwpkn",route="canary",server="10.132.149.137:8080",service="ingress-canary"} 0
haproxy_server_http_responses_total{code="3xx",namespace="openshift-ingress-canary",pod="ingress-canary-z2b4x",route="canary",server="10.134.248.101:8080",service="ingress-canary"} 0
haproxy_server_http_responses_total{code="3xx",namespace="openshift-monitoring",pod="alertmanager-main-0",route="alertmanager-main",server="10.134.173.99:9095",service="alertmanager-main"} 0
haproxy_server_http_responses_total{code="3xx",namespace="openshift-monitoring",pod="alertmanager-main-1",route="alertmanager-main",server="10.134.248.127:9095",service="alertmanager-main"} 0
haproxy_server_http_responses_total{code="3xx",namespace="openshift-monitoring",pod="prometheus-k8s-0",route="prometheus-k8s",server="10.134.173.51:9091",service="prometheus-k8s"} 0
haproxy_server_http_responses_total{code="3xx",namespace="openshift-monitoring",pod="prometheus-k8s-0",route="prometheus-k8s-federate",server="10.134.173.51:9091",service="prometheus-k8s"} 0
haproxy_server_http_responses_total{code="3xx",namespace="openshift-monitoring",pod="prometheus-k8s-1",route="prometheus-k8s",server="10.134.249.153:9091",service="prometheus-k8s"} 0
haproxy_server_http_responses_total{code="3xx",namespace="openshift-monitoring",pod="prometheus-k8s-1",route="prometheus-k8s-federate",server="10.134.249.153:9091",service="prometheus-k8s"} 0
haproxy_server_http_responses_total{code="3xx",namespace="openshift-monitoring",pod="thanos-querier-59d7db9787-hwjzt",route="thanos-querier",server="10.134.249.65:9091",service="thanos-querier"} 0
haproxy_server_http_responses_total{code="3xx",namespace="openshift-monitoring",pod="thanos-querier-59d7db9787-p8m9l",route="thanos-querier",server="10.134.173.112:9091",service="thanos-querier"} 0
haproxy_server_http_responses_total{code="3xx",namespace="profiler",pod="esc-static-service-654f8bc589-xchjk",route="esc-static-service-j9sgv",server="10.134.173.123:8080",service="esc-static-service"} 0
haproxy_server_http_responses_total{code="3xx",namespace="profiler",pod="esc-test-service-75849d7c6c-czvsw",route="esc-test-service-lbp8r",server="10.134.172.182:8080",service="esc-test-service"} 0
haproxy_server_http_responses_total{code="3xx",namespace="profiler",pod="esc-ui-service-66977c9bcb-kg9mc",route="esc-ui-service-r5tfm",server="10.134.249.215:8180",service="esc-ui-service"} 0
haproxy_server_http_responses_total{code="3xx",namespace="prometheus-operator",pod="grafana-deployment-88d46bc5f-5gx7h",route="prometheus-operator-grafana",server="10.132.149.3:3000",service="grafana-service"} 0
haproxy_server_http_responses_total{code="3xx",namespace="prometheus-operator",pod="vmagent-k8s-67654d6f97-98nf5",route="prometheus-operator-vmagent-k8s",server="10.132.149.156:8429",service="vmagent-k8s"} 0
haproxy_server_http_responses_total{code="3xx",namespace="prometheus-operator",pod="vmalert-k8s-6bcf966dd6-jvpz5",route="prometheus-operator-vmalert-k8s",server="10.132.149.175:8080",service="vmalert-k8s"} 0
haproxy_server_http_responses_total{code="3xx",namespace="prometheus-operator",pod="vmalertmanager-k8s-0",route="prometheus-operator-vmalertmanager-k8s",server="10.134.173.133:9093",service="vmalertmanager-k8s"} 0
haproxy_server_http_responses_total{code="3xx",namespace="prometheus-operator",pod="vmsingle-k8s-f4844bfc4-5l8sf",route="prometheus-operator-vmsingle-k8s",server="10.134.249.110:8429",service="vmsingle-k8s"} 0
haproxy_server_http_responses_total{code="3xx",namespace="rabbitmq",pod="rmqlocal-0",route="rabbitmq-fnzhp",server="10.134.249.141:15672",service="rabbitmq"} 0
haproxy_server_http_responses_total{code="3xx",namespace="rabbitmq",pod="rmqlocal-1",route="rabbitmq-fnzhp",server="10.134.173.127:15672",service="rabbitmq"} 0
haproxy_server_http_responses_total{code="3xx",namespace="rabbitmq",pod="rmqlocal-2",route="rabbitmq-fnzhp",server="10.132.149.46:15672",service="rabbitmq"} 0
haproxy_server_http_responses_total{code="3xx",namespace="shared-license-distributor",pod="shared-license-distributor-64ccc89449-8kw8p",route="shared-license-distributor-wnj68",server="10.134.249.239:8080",service="shared-license-distributor"} 0
haproxy_server_http_responses_total{code="3xx",namespace="spark",pod="spark-history-server-bd755d9f5-9ccsz",route="spark-history-server-hqshb",server="10.132.149.178:18080",service="spark-history-server"} 0
haproxy_server_http_responses_total{code="3xx",namespace="spark",pod="spark-history-server-bd755d9f5-zf2dt",route="spark-history-server-hqshb",server="10.134.173.147:18080",service="spark-history-server"} 0
haproxy_server_http_responses_total{code="3xx",namespace="spark",pod="spark-site-manager-66c44c7b78-cc7jr",route="spark-site-manager-p4nqh",server="10.132.149.148:8080",service="spark-site-manager"} 0
haproxy_server_http_responses_total{code="4xx",namespace="airflow",pod="airflow-site-manager-68ccf87889-bwbdf",route="airflow-site-manager-hpl66",server="10.132.149.80:8080",service="airflow-site-manager"} 0
haproxy_server_http_responses_total{code="4xx",namespace="airflow",pod="airflow-webserver-695b5c545b-rxqc4",route="airflow-ingress-nzbtm",server="10.134.173.137:8080",service="airflow-webserver"} 0
haproxy_server_http_responses_total{code="4xx",namespace="arangodb",pod="main-arangodb-1-87669b688-q98z4",route="arangodb-leader-ingress-zw82p",server="10.132.149.185:8529",service="arangodb-leader"} 0
haproxy_server_http_responses_total{code="4xx",namespace="consul-service",pod="consul-server-0",route="consul-ingress-hmmvf",server="10.134.248.50:8500",service="consul-ui"} 0
haproxy_server_http_responses_total{code="4xx",namespace="consul-service",pod="consul-server-1",route="consul-ingress-hmmvf",server="10.132.149.160:8500",service="consul-ui"} 0
haproxy_server_http_responses_total{code="4xx",namespace="consul-service",pod="consul-server-2",route="consul-ingress-hmmvf",server="10.134.173.136:8500",service="consul-ui"} 0
haproxy_server_http_responses_total{code="4xx",namespace="dbaas",pod="dbaas-aggregator-5f744d74f5-r7tqm",route="aggregator-fw5ft",server="10.134.248.185:8080",service="dbaas-aggregator"} 0
haproxy_server_http_responses_total{code="4xx",namespace="infra-keycloak",pod="infra-keycloak-7fd57cd59c-d6c5b",route="infra-keycloak-5q4rj",server="10.132.148.223:8080",service="infra-keycloak"} 0
haproxy_server_http_responses_total{code="4xx",namespace="infra-keycloak",pod="infra-keycloak-7fd57cd59c-qtqv9",route="infra-keycloak-5q4rj",server="10.134.173.182:8080",service="infra-keycloak"} 0
haproxy_server_http_responses_total{code="4xx",namespace="jaeger",pod="jaeger-hotrod-85b766c49d-s6vz7",route="jaeger-hotrod-qq5cq",server="10.134.249.205:8080",service="jaeger-hotrod"} 0
haproxy_server_http_responses_total{code="4xx",namespace="jaeger",pod="jaeger-query-65d9446fc6-jx5qk",route="jaeger-query-pc6w5",server="10.134.248.54:16686",service="jaeger-query"} 0
haproxy_server_http_responses_total{code="4xx",namespace="kafka-service",pod="akhq-58744f76db-h2mtf",route="akhq-ingress-lhz5z",server="10.134.249.133:8080",service="akhq"} 0
haproxy_server_http_responses_total{code="4xx",namespace="license-server",pod="license-server-7b8b6f765f-wvcx5",route="license-server-9txj9",server="10.132.148.47:8080",service="license-server"} 0
haproxy_server_http_responses_total{code="4xx",namespace="opensearch",pod="opensearch-dashboards-764bc8548b-zngcr",route="opensearch-dashboards-l6vdv",server="10.134.173.229:5601",service="opensearch-dashboards"} 0
haproxy_server_http_responses_total{code="4xx",namespace="openshift-console",pod="console-66fcc4f4fc-crbz6",route="console",server="10.133.43.162:8443",service="console"} 209
haproxy_server_http_responses_total{code="4xx",namespace="openshift-console",pod="console-66fcc4f4fc-j5dvz",route="console",server="10.132.153.185:8443",service="console"} 96
haproxy_server_http_responses_total{code="4xx",namespace="openshift-console",pod="downloads-6fb844fb47-cwnvw",route="downloads",server="10.135.207.221:8080",service="downloads"} 0
haproxy_server_http_responses_total{code="4xx",namespace="openshift-console",pod="downloads-6fb844fb47-l485s",route="downloads",server="10.132.153.222:8080",service="downloads"} 0
haproxy_server_http_responses_total{code="4xx",namespace="openshift-ingress-canary",pod="ingress-canary-7rvzj",route="canary",server="10.134.173.247:8080",service="ingress-canary"} 0
haproxy_server_http_responses_total{code="4xx",namespace="openshift-ingress-canary",pod="ingress-canary-lwpkn",route="canary",server="10.132.149.137:8080",service="ingress-canary"} 0
haproxy_server_http_responses_total{code="4xx",namespace="openshift-ingress-canary",pod="ingress-canary-z2b4x",route="canary",server="10.134.248.101:8080",service="ingress-canary"} 0
haproxy_server_http_responses_total{code="4xx",namespace="openshift-monitoring",pod="alertmanager-main-0",route="alertmanager-main",server="10.134.173.99:9095",service="alertmanager-main"} 0
haproxy_server_http_responses_total{code="4xx",namespace="openshift-monitoring",pod="alertmanager-main-1",route="alertmanager-main",server="10.134.248.127:9095",service="alertmanager-main"} 0
haproxy_server_http_responses_total{code="4xx",namespace="openshift-monitoring",pod="prometheus-k8s-0",route="prometheus-k8s",server="10.134.173.51:9091",service="prometheus-k8s"} 0
haproxy_server_http_responses_total{code="4xx",namespace="openshift-monitoring",pod="prometheus-k8s-0",route="prometheus-k8s-federate",server="10.134.173.51:9091",service="prometheus-k8s"} 0
haproxy_server_http_responses_total{code="4xx",namespace="openshift-monitoring",pod="prometheus-k8s-1",route="prometheus-k8s",server="10.134.249.153:9091",service="prometheus-k8s"} 0
haproxy_server_http_responses_total{code="4xx",namespace="openshift-monitoring",pod="prometheus-k8s-1",route="prometheus-k8s-federate",server="10.134.249.153:9091",service="prometheus-k8s"} 0
haproxy_server_http_responses_total{code="4xx",namespace="openshift-monitoring",pod="thanos-querier-59d7db9787-hwjzt",route="thanos-querier",server="10.134.249.65:9091",service="thanos-querier"} 0
haproxy_server_http_responses_total{code="4xx",namespace="openshift-monitoring",pod="thanos-querier-59d7db9787-p8m9l",route="thanos-querier",server="10.134.173.112:9091",service="thanos-querier"} 0
haproxy_server_http_responses_total{code="4xx",namespace="profiler",pod="esc-static-service-654f8bc589-xchjk",route="esc-static-service-j9sgv",server="10.134.173.123:8080",service="esc-static-service"} 0
haproxy_server_http_responses_total{code="4xx",namespace="profiler",pod="esc-test-service-75849d7c6c-czvsw",route="esc-test-service-lbp8r",server="10.134.172.182:8080",service="esc-test-service"} 0
haproxy_server_http_responses_total{code="4xx",namespace="profiler",pod="esc-ui-service-66977c9bcb-kg9mc",route="esc-ui-service-r5tfm",server="10.134.249.215:8180",service="esc-ui-service"} 0
haproxy_server_http_responses_total{code="4xx",namespace="prometheus-operator",pod="grafana-deployment-88d46bc5f-5gx7h",route="prometheus-operator-grafana",server="10.132.149.3:3000",service="grafana-service"} 0
haproxy_server_http_responses_total{code="4xx",namespace="prometheus-operator",pod="vmagent-k8s-67654d6f97-98nf5",route="prometheus-operator-vmagent-k8s",server="10.132.149.156:8429",service="vmagent-k8s"} 0
haproxy_server_http_responses_total{code="4xx",namespace="prometheus-operator",pod="vmalert-k8s-6bcf966dd6-jvpz5",route="prometheus-operator-vmalert-k8s",server="10.132.149.175:8080",service="vmalert-k8s"} 0
haproxy_server_http_responses_total{code="4xx",namespace="prometheus-operator",pod="vmalertmanager-k8s-0",route="prometheus-operator-vmalertmanager-k8s",server="10.134.173.133:9093",service="vmalertmanager-k8s"} 0
haproxy_server_http_responses_total{code="4xx",namespace="prometheus-operator",pod="vmsingle-k8s-f4844bfc4-5l8sf",route="prometheus-operator-vmsingle-k8s",server="10.134.249.110:8429",service="vmsingle-k8s"} 0
haproxy_server_http_responses_total{code="4xx",namespace="rabbitmq",pod="rmqlocal-0",route="rabbitmq-fnzhp",server="10.134.249.141:15672",service="rabbitmq"} 0
haproxy_server_http_responses_total{code="4xx",namespace="rabbitmq",pod="rmqlocal-1",route="rabbitmq-fnzhp",server="10.134.173.127:15672",service="rabbitmq"} 0
haproxy_server_http_responses_total{code="4xx",namespace="rabbitmq",pod="rmqlocal-2",route="rabbitmq-fnzhp",server="10.132.149.46:15672",service="rabbitmq"} 0
haproxy_server_http_responses_total{code="4xx",namespace="shared-license-distributor",pod="shared-license-distributor-64ccc89449-8kw8p",route="shared-license-distributor-wnj68",server="10.134.249.239:8080",service="shared-license-distributor"} 0
haproxy_server_http_responses_total{code="4xx",namespace="spark",pod="spark-history-server-bd755d9f5-9ccsz",route="spark-history-server-hqshb",server="10.132.149.178:18080",service="spark-history-server"} 0
haproxy_server_http_responses_total{code="4xx",namespace="spark",pod="spark-history-server-bd755d9f5-zf2dt",route="spark-history-server-hqshb",server="10.134.173.147:18080",service="spark-history-server"} 0
haproxy_server_http_responses_total{code="4xx",namespace="spark",pod="spark-site-manager-66c44c7b78-cc7jr",route="spark-site-manager-p4nqh",server="10.132.149.148:8080",service="spark-site-manager"} 0
haproxy_server_http_responses_total{code="5xx",namespace="airflow",pod="airflow-site-manager-68ccf87889-bwbdf",route="airflow-site-manager-hpl66",server="10.132.149.80:8080",service="airflow-site-manager"} 0
haproxy_server_http_responses_total{code="5xx",namespace="airflow",pod="airflow-webserver-695b5c545b-rxqc4",route="airflow-ingress-nzbtm",server="10.134.173.137:8080",service="airflow-webserver"} 0
haproxy_server_http_responses_total{code="5xx",namespace="arangodb",pod="main-arangodb-1-87669b688-q98z4",route="arangodb-leader-ingress-zw82p",server="10.132.149.185:8529",service="arangodb-leader"} 0
haproxy_server_http_responses_total{code="5xx",namespace="consul-service",pod="consul-server-0",route="consul-ingress-hmmvf",server="10.134.248.50:8500",service="consul-ui"} 0
haproxy_server_http_responses_total{code="5xx",namespace="consul-service",pod="consul-server-1",route="consul-ingress-hmmvf",server="10.132.149.160:8500",service="consul-ui"} 0
haproxy_server_http_responses_total{code="5xx",namespace="consul-service",pod="consul-server-2",route="consul-ingress-hmmvf",server="10.134.173.136:8500",service="consul-ui"} 0
haproxy_server_http_responses_total{code="5xx",namespace="dbaas",pod="dbaas-aggregator-5f744d74f5-r7tqm",route="aggregator-fw5ft",server="10.134.248.185:8080",service="dbaas-aggregator"} 0
haproxy_server_http_responses_total{code="5xx",namespace="infra-keycloak",pod="infra-keycloak-7fd57cd59c-d6c5b",route="infra-keycloak-5q4rj",server="10.132.148.223:8080",service="infra-keycloak"} 0
haproxy_server_http_responses_total{code="5xx",namespace="infra-keycloak",pod="infra-keycloak-7fd57cd59c-qtqv9",route="infra-keycloak-5q4rj",server="10.134.173.182:8080",service="infra-keycloak"} 0
haproxy_server_http_responses_total{code="5xx",namespace="jaeger",pod="jaeger-hotrod-85b766c49d-s6vz7",route="jaeger-hotrod-qq5cq",server="10.134.249.205:8080",service="jaeger-hotrod"} 0
haproxy_server_http_responses_total{code="5xx",namespace="jaeger",pod="jaeger-query-65d9446fc6-jx5qk",route="jaeger-query-pc6w5",server="10.134.248.54:16686",service="jaeger-query"} 0
haproxy_server_http_responses_total{code="5xx",namespace="kafka-service",pod="akhq-58744f76db-h2mtf",route="akhq-ingress-lhz5z",server="10.134.249.133:8080",service="akhq"} 0
haproxy_server_http_responses_total{code="5xx",namespace="license-server",pod="license-server-7b8b6f765f-wvcx5",route="license-server-9txj9",server="10.132.148.47:8080",service="license-server"} 0
haproxy_server_http_responses_total{code="5xx",namespace="opensearch",pod="opensearch-dashboards-764bc8548b-zngcr",route="opensearch-dashboards-l6vdv",server="10.134.173.229:5601",service="opensearch-dashboards"} 0
haproxy_server_http_responses_total{code="5xx",namespace="openshift-console",pod="console-66fcc4f4fc-crbz6",route="console",server="10.133.43.162:8443",service="console"} 0
haproxy_server_http_responses_total{code="5xx",namespace="openshift-console",pod="console-66fcc4f4fc-j5dvz",route="console",server="10.132.153.185:8443",service="console"} 0
haproxy_server_http_responses_total{code="5xx",namespace="openshift-console",pod="downloads-6fb844fb47-cwnvw",route="downloads",server="10.135.207.221:8080",service="downloads"} 0
haproxy_server_http_responses_total{code="5xx",namespace="openshift-console",pod="downloads-6fb844fb47-l485s",route="downloads",server="10.132.153.222:8080",service="downloads"} 0
haproxy_server_http_responses_total{code="5xx",namespace="openshift-ingress-canary",pod="ingress-canary-7rvzj",route="canary",server="10.134.173.247:8080",service="ingress-canary"} 0
haproxy_server_http_responses_total{code="5xx",namespace="openshift-ingress-canary",pod="ingress-canary-lwpkn",route="canary",server="10.132.149.137:8080",service="ingress-canary"} 0
haproxy_server_http_responses_total{code="5xx",namespace="openshift-ingress-canary",pod="ingress-canary-z2b4x",route="canary",server="10.134.248.101:8080",service="ingress-canary"} 0
haproxy_server_http_responses_total{code="5xx",namespace="openshift-monitoring",pod="alertmanager-main-0",route="alertmanager-main",server="10.134.173.99:9095",service="alertmanager-main"} 0
haproxy_server_http_responses_total{code="5xx",namespace="openshift-monitoring",pod="alertmanager-main-1",route="alertmanager-main",server="10.134.248.127:9095",service="alertmanager-main"} 0
haproxy_server_http_responses_total{code="5xx",namespace="openshift-monitoring",pod="prometheus-k8s-0",route="prometheus-k8s",server="10.134.173.51:9091",service="prometheus-k8s"} 0
haproxy_server_http_responses_total{code="5xx",namespace="openshift-monitoring",pod="prometheus-k8s-0",route="prometheus-k8s-federate",server="10.134.173.51:9091",service="prometheus-k8s"} 0
haproxy_server_http_responses_total{code="5xx",namespace="openshift-monitoring",pod="prometheus-k8s-1",route="prometheus-k8s",server="10.134.249.153:9091",service="prometheus-k8s"} 0
haproxy_server_http_responses_total{code="5xx",namespace="openshift-monitoring",pod="prometheus-k8s-1",route="prometheus-k8s-federate",server="10.134.249.153:9091",service="prometheus-k8s"} 0
haproxy_server_http_responses_total{code="5xx",namespace="openshift-monitoring",pod="thanos-querier-59d7db9787-hwjzt",route="thanos-querier",server="10.134.249.65:9091",service="thanos-querier"} 0
haproxy_server_http_responses_total{code="5xx",namespace="openshift-monitoring",pod="thanos-querier-59d7db9787-p8m9l",route="thanos-querier",server="10.134.173.112:9091",service="thanos-querier"} 0
haproxy_server_http_responses_total{code="5xx",namespace="profiler",pod="esc-static-service-654f8bc589-xchjk",route="esc-static-service-j9sgv",server="10.134.173.123:8080",service="esc-static-service"} 0
haproxy_server_http_responses_total{code="5xx",namespace="profiler",pod="esc-test-service-75849d7c6c-czvsw",route="esc-test-service-lbp8r",server="10.134.172.182:8080",service="esc-test-service"} 0
haproxy_server_http_responses_total{code="5xx",namespace="profiler",pod="esc-ui-service-66977c9bcb-kg9mc",route="esc-ui-service-r5tfm",server="10.134.249.215:8180",service="esc-ui-service"} 0
haproxy_server_http_responses_total{code="5xx",namespace="prometheus-operator",pod="grafana-deployment-88d46bc5f-5gx7h",route="prometheus-operator-grafana",server="10.132.149.3:3000",service="grafana-service"} 0
haproxy_server_http_responses_total{code="5xx",namespace="prometheus-operator",pod="vmagent-k8s-67654d6f97-98nf5",route="prometheus-operator-vmagent-k8s",server="10.132.149.156:8429",service="vmagent-k8s"} 0
haproxy_server_http_responses_total{code="5xx",namespace="prometheus-operator",pod="vmalert-k8s-6bcf966dd6-jvpz5",route="prometheus-operator-vmalert-k8s",server="10.132.149.175:8080",service="vmalert-k8s"} 0
haproxy_server_http_responses_total{code="5xx",namespace="prometheus-operator",pod="vmalertmanager-k8s-0",route="prometheus-operator-vmalertmanager-k8s",server="10.134.173.133:9093",service="vmalertmanager-k8s"} 0
haproxy_server_http_responses_total{code="5xx",namespace="prometheus-operator",pod="vmsingle-k8s-f4844bfc4-5l8sf",route="prometheus-operator-vmsingle-k8s",server="10.134.249.110:8429",service="vmsingle-k8s"} 0
haproxy_server_http_responses_total{code="5xx",namespace="rabbitmq",pod="rmqlocal-0",route="rabbitmq-fnzhp",server="10.134.249.141:15672",service="rabbitmq"} 0
haproxy_server_http_responses_total{code="5xx",namespace="rabbitmq",pod="rmqlocal-1",route="rabbitmq-fnzhp",server="10.134.173.127:15672",service="rabbitmq"} 0
haproxy_server_http_responses_total{code="5xx",namespace="rabbitmq",pod="rmqlocal-2",route="rabbitmq-fnzhp",server="10.132.149.46:15672",service="rabbitmq"} 0
haproxy_server_http_responses_total{code="5xx",namespace="shared-license-distributor",pod="shared-license-distributor-64ccc89449-8kw8p",route="shared-license-distributor-wnj68",server="10.134.249.239:8080",service="shared-license-distributor"} 0
haproxy_server_http_responses_total{code="5xx",namespace="spark",pod="spark-history-server-bd755d9f5-9ccsz",route="spark-history-server-hqshb",server="10.132.149.178:18080",service="spark-history-server"} 0
haproxy_server_http_responses_total{code="5xx",namespace="spark",pod="spark-history-server-bd755d9f5-zf2dt",route="spark-history-server-hqshb",server="10.134.173.147:18080",service="spark-history-server"} 0
haproxy_server_http_responses_total{code="5xx",namespace="spark",pod="spark-site-manager-66c44c7b78-cc7jr",route="spark-site-manager-p4nqh",server="10.132.149.148:8080",service="spark-site-manager"} 0
haproxy_server_http_responses_total{code="other",namespace="airflow",pod="airflow-site-manager-68ccf87889-bwbdf",route="airflow-site-manager-hpl66",server="10.132.149.80:8080",service="airflow-site-manager"} 0
haproxy_server_http_responses_total{code="other",namespace="airflow",pod="airflow-webserver-695b5c545b-rxqc4",route="airflow-ingress-nzbtm",server="10.134.173.137:8080",service="airflow-webserver"} 0
haproxy_server_http_responses_total{code="other",namespace="arangodb",pod="main-arangodb-1-87669b688-q98z4",route="arangodb-leader-ingress-zw82p",server="10.132.149.185:8529",service="arangodb-leader"} 0
haproxy_server_http_responses_total{code="other",namespace="consul-service",pod="consul-server-0",route="consul-ingress-hmmvf",server="10.134.248.50:8500",service="consul-ui"} 0
haproxy_server_http_responses_total{code="other",namespace="consul-service",pod="consul-server-1",route="consul-ingress-hmmvf",server="10.132.149.160:8500",service="consul-ui"} 0
haproxy_server_http_responses_total{code="other",namespace="consul-service",pod="consul-server-2",route="consul-ingress-hmmvf",server="10.134.173.136:8500",service="consul-ui"} 0
haproxy_server_http_responses_total{code="other",namespace="dbaas",pod="dbaas-aggregator-5f744d74f5-r7tqm",route="aggregator-fw5ft",server="10.134.248.185:8080",service="dbaas-aggregator"} 0
haproxy_server_http_responses_total{code="other",namespace="infra-keycloak",pod="infra-keycloak-7fd57cd59c-d6c5b",route="infra-keycloak-5q4rj",server="10.132.148.223:8080",service="infra-keycloak"} 0
haproxy_server_http_responses_total{code="other",namespace="infra-keycloak",pod="infra-keycloak-7fd57cd59c-qtqv9",route="infra-keycloak-5q4rj",server="10.134.173.182:8080",service="infra-keycloak"} 0
haproxy_server_http_responses_total{code="other",namespace="jaeger",pod="jaeger-hotrod-85b766c49d-s6vz7",route="jaeger-hotrod-qq5cq",server="10.134.249.205:8080",service="jaeger-hotrod"} 0
haproxy_server_http_responses_total{code="other",namespace="jaeger",pod="jaeger-query-65d9446fc6-jx5qk",route="jaeger-query-pc6w5",server="10.134.248.54:16686",service="jaeger-query"} 0
haproxy_server_http_responses_total{code="other",namespace="kafka-service",pod="akhq-58744f76db-h2mtf",route="akhq-ingress-lhz5z",server="10.134.249.133:8080",service="akhq"} 0
haproxy_server_http_responses_total{code="other",namespace="license-server",pod="license-server-7b8b6f765f-wvcx5",route="license-server-9txj9",server="10.132.148.47:8080",service="license-server"} 0
haproxy_server_http_responses_total{code="other",namespace="opensearch",pod="opensearch-dashboards-764bc8548b-zngcr",route="opensearch-dashboards-l6vdv",server="10.134.173.229:5601",service="opensearch-dashboards"} 0
haproxy_server_http_responses_total{code="other",namespace="openshift-console",pod="console-66fcc4f4fc-crbz6",route="console",server="10.133.43.162:8443",service="console"} 0
haproxy_server_http_responses_total{code="other",namespace="openshift-console",pod="console-66fcc4f4fc-j5dvz",route="console",server="10.132.153.185:8443",service="console"} 0
haproxy_server_http_responses_total{code="other",namespace="openshift-console",pod="downloads-6fb844fb47-cwnvw",route="downloads",server="10.135.207.221:8080",service="downloads"} 0
haproxy_server_http_responses_total{code="other",namespace="openshift-console",pod="downloads-6fb844fb47-l485s",route="downloads",server="10.132.153.222:8080",service="downloads"} 0
haproxy_server_http_responses_total{code="other",namespace="openshift-ingress-canary",pod="ingress-canary-7rvzj",route="canary",server="10.134.173.247:8080",service="ingress-canary"} 0
haproxy_server_http_responses_total{code="other",namespace="openshift-ingress-canary",pod="ingress-canary-lwpkn",route="canary",server="10.132.149.137:8080",service="ingress-canary"} 0
haproxy_server_http_responses_total{code="other",namespace="openshift-ingress-canary",pod="ingress-canary-z2b4x",route="canary",server="10.134.248.101:8080",service="ingress-canary"} 0
haproxy_server_http_responses_total{code="other",namespace="openshift-monitoring",pod="alertmanager-main-0",route="alertmanager-main",server="10.134.173.99:9095",service="alertmanager-main"} 0
haproxy_server_http_responses_total{code="other",namespace="openshift-monitoring",pod="alertmanager-main-1",route="alertmanager-main",server="10.134.248.127:9095",service="alertmanager-main"} 0
haproxy_server_http_responses_total{code="other",namespace="openshift-monitoring",pod="prometheus-k8s-0",route="prometheus-k8s",server="10.134.173.51:9091",service="prometheus-k8s"} 0
haproxy_server_http_responses_total{code="other",namespace="openshift-monitoring",pod="prometheus-k8s-0",route="prometheus-k8s-federate",server="10.134.173.51:9091",service="prometheus-k8s"} 0
haproxy_server_http_responses_total{code="other",namespace="openshift-monitoring",pod="prometheus-k8s-1",route="prometheus-k8s",server="10.134.249.153:9091",service="prometheus-k8s"} 0
haproxy_server_http_responses_total{code="other",namespace="openshift-monitoring",pod="prometheus-k8s-1",route="prometheus-k8s-federate",server="10.134.249.153:9091",service="prometheus-k8s"} 0
haproxy_server_http_responses_total{code="other",namespace="openshift-monitoring",pod="thanos-querier-59d7db9787-hwjzt",route="thanos-querier",server="10.134.249.65:9091",service="thanos-querier"} 0
haproxy_server_http_responses_total{code="other",namespace="openshift-monitoring",pod="thanos-querier-59d7db9787-p8m9l",route="thanos-querier",server="10.134.173.112:9091",service="thanos-querier"} 0
haproxy_server_http_responses_total{code="other",namespace="profiler",pod="esc-static-service-654f8bc589-xchjk",route="esc-static-service-j9sgv",server="10.134.173.123:8080",service="esc-static-service"} 0
haproxy_server_http_responses_total{code="other",namespace="profiler",pod="esc-test-service-75849d7c6c-czvsw",route="esc-test-service-lbp8r",server="10.134.172.182:8080",service="esc-test-service"} 0
haproxy_server_http_responses_total{code="other",namespace="profiler",pod="esc-ui-service-66977c9bcb-kg9mc",route="esc-ui-service-r5tfm",server="10.134.249.215:8180",service="esc-ui-service"} 0
haproxy_server_http_responses_total{code="other",namespace="prometheus-operator",pod="grafana-deployment-88d46bc5f-5gx7h",route="prometheus-operator-grafana",server="10.132.149.3:3000",service="grafana-service"} 0
haproxy_server_http_responses_total{code="other",namespace="prometheus-operator",pod="vmagent-k8s-67654d6f97-98nf5",route="prometheus-operator-vmagent-k8s",server="10.132.149.156:8429",service="vmagent-k8s"} 0
haproxy_server_http_responses_total{code="other",namespace="prometheus-operator",pod="vmalert-k8s-6bcf966dd6-jvpz5",route="prometheus-operator-vmalert-k8s",server="10.132.149.175:8080",service="vmalert-k8s"} 0
haproxy_server_http_responses_total{code="other",namespace="prometheus-operator",pod="vmalertmanager-k8s-0",route="prometheus-operator-vmalertmanager-k8s",server="10.134.173.133:9093",service="vmalertmanager-k8s"} 0
haproxy_server_http_responses_total{code="other",namespace="prometheus-operator",pod="vmsingle-k8s-f4844bfc4-5l8sf",route="prometheus-operator-vmsingle-k8s",server="10.134.249.110:8429",service="vmsingle-k8s"} 0
haproxy_server_http_responses_total{code="other",namespace="rabbitmq",pod="rmqlocal-0",route="rabbitmq-fnzhp",server="10.134.249.141:15672",service="rabbitmq"} 0
haproxy_server_http_responses_total{code="other",namespace="rabbitmq",pod="rmqlocal-1",route="rabbitmq-fnzhp",server="10.134.173.127:15672",service="rabbitmq"} 0
haproxy_server_http_responses_total{code="other",namespace="rabbitmq",pod="rmqlocal-2",route="rabbitmq-fnzhp",server="10.132.149.46:15672",service="rabbitmq"} 0
haproxy_server_http_responses_total{code="other",namespace="shared-license-distributor",pod="shared-license-distributor-64ccc89449-8kw8p",route="shared-license-distributor-wnj68",server="10.134.249.239:8080",service="shared-license-distributor"} 0
haproxy_server_http_responses_total{code="other",namespace="spark",pod="spark-history-server-bd755d9f5-9ccsz",route="spark-history-server-hqshb",server="10.132.149.178:18080",service="spark-history-server"} 0
haproxy_server_http_responses_total{code="other",namespace="spark",pod="spark-history-server-bd755d9f5-zf2dt",route="spark-history-server-hqshb",server="10.134.173.147:18080",service="spark-history-server"} 0
haproxy_server_http_responses_total{code="other",namespace="spark",pod="spark-site-manager-66c44c7b78-cc7jr",route="spark-site-manager-p4nqh",server="10.132.149.148:8080",service="spark-site-manager"} 0
# HELP haproxy_server_max_session_rate Maximum observed number of sessions per second.
# TYPE haproxy_server_max_session_rate gauge
haproxy_server_max_session_rate{namespace="",pod="",route="",server="fe_no_sni",service=""} 0
haproxy_server_max_session_rate{namespace="",pod="",route="",server="fe_sni",service=""} 0
haproxy_server_max_session_rate{namespace="airflow",pod="airflow-site-manager-68ccf87889-bwbdf",route="airflow-site-manager-hpl66",server="10.132.149.80:8080",service="airflow-site-manager"} 0
haproxy_server_max_session_rate{namespace="airflow",pod="airflow-webserver-695b5c545b-rxqc4",route="airflow-ingress-nzbtm",server="10.134.173.137:8080",service="airflow-webserver"} 0
haproxy_server_max_session_rate{namespace="arangodb",pod="main-arangodb-1-87669b688-q98z4",route="arangodb-leader-ingress-zw82p",server="10.132.149.185:8529",service="arangodb-leader"} 0
haproxy_server_max_session_rate{namespace="consul-service",pod="consul-server-0",route="consul-ingress-hmmvf",server="10.134.248.50:8500",service="consul-ui"} 0
haproxy_server_max_session_rate{namespace="consul-service",pod="consul-server-1",route="consul-ingress-hmmvf",server="10.132.149.160:8500",service="consul-ui"} 0
haproxy_server_max_session_rate{namespace="consul-service",pod="consul-server-2",route="consul-ingress-hmmvf",server="10.134.173.136:8500",service="consul-ui"} 0
haproxy_server_max_session_rate{namespace="dbaas",pod="dbaas-aggregator-5f744d74f5-r7tqm",route="aggregator-fw5ft",server="10.134.248.185:8080",service="dbaas-aggregator"} 0
haproxy_server_max_session_rate{namespace="infra-keycloak",pod="infra-keycloak-7fd57cd59c-d6c5b",route="infra-keycloak-5q4rj",server="10.132.148.223:8080",service="infra-keycloak"} 0
haproxy_server_max_session_rate{namespace="infra-keycloak",pod="infra-keycloak-7fd57cd59c-qtqv9",route="infra-keycloak-5q4rj",server="10.134.173.182:8080",service="infra-keycloak"} 0
haproxy_server_max_session_rate{namespace="jaeger",pod="jaeger-hotrod-85b766c49d-s6vz7",route="jaeger-hotrod-qq5cq",server="10.134.249.205:8080",service="jaeger-hotrod"} 0
haproxy_server_max_session_rate{namespace="jaeger",pod="jaeger-query-65d9446fc6-jx5qk",route="jaeger-query-pc6w5",server="10.134.248.54:16686",service="jaeger-query"} 0
haproxy_server_max_session_rate{namespace="kafka-service",pod="akhq-58744f76db-h2mtf",route="akhq-ingress-lhz5z",server="10.134.249.133:8080",service="akhq"} 0
haproxy_server_max_session_rate{namespace="license-server",pod="license-server-7b8b6f765f-wvcx5",route="license-server-9txj9",server="10.132.148.47:8080",service="license-server"} 0
haproxy_server_max_session_rate{namespace="opensearch",pod="opensearch-dashboards-764bc8548b-zngcr",route="opensearch-dashboards-l6vdv",server="10.134.173.229:5601",service="opensearch-dashboards"} 0
haproxy_server_max_session_rate{namespace="openshift-authentication",pod="oauth-openshift-75965c7487-kc8m7",route="oauth-openshift",server="10.132.153.132:6443",service="oauth-openshift"} 0
haproxy_server_max_session_rate{namespace="openshift-authentication",pod="oauth-openshift-75965c7487-nsfxw",route="oauth-openshift",server="10.133.43.151:6443",service="oauth-openshift"} 0
haproxy_server_max_session_rate{namespace="openshift-authentication",pod="oauth-openshift-75965c7487-xmmmq",route="oauth-openshift",server="10.135.207.187:6443",service="oauth-openshift"} 0
haproxy_server_max_session_rate{namespace="openshift-console",pod="console-66fcc4f4fc-crbz6",route="console",server="10.133.43.162:8443",service="console"} 0
haproxy_server_max_session_rate{namespace="openshift-console",pod="console-66fcc4f4fc-j5dvz",route="console",server="10.132.153.185:8443",service="console"} 0
haproxy_server_max_session_rate{namespace="openshift-console",pod="downloads-6fb844fb47-cwnvw",route="downloads",server="10.135.207.221:8080",service="downloads"} 0
haproxy_server_max_session_rate{namespace="openshift-console",pod="downloads-6fb844fb47-l485s",route="downloads",server="10.132.153.222:8080",service="downloads"} 0
haproxy_server_max_session_rate{namespace="openshift-ingress-canary",pod="ingress-canary-7rvzj",route="canary",server="10.134.173.247:8080",service="ingress-canary"} 0
haproxy_server_max_session_rate{namespace="openshift-ingress-canary",pod="ingress-canary-lwpkn",route="canary",server="10.132.149.137:8080",service="ingress-canary"} 0
haproxy_server_max_session_rate{namespace="openshift-ingress-canary",pod="ingress-canary-z2b4x",route="canary",server="10.134.248.101:8080",service="ingress-canary"} 0
haproxy_server_max_session_rate{namespace="openshift-monitoring",pod="alertmanager-main-0",route="alertmanager-main",server="10.134.173.99:9095",service="alertmanager-main"} 0
haproxy_server_max_session_rate{namespace="openshift-monitoring",pod="alertmanager-main-1",route="alertmanager-main",server="10.134.248.127:9095",service="alertmanager-main"} 0
haproxy_server_max_session_rate{namespace="openshift-monitoring",pod="prometheus-k8s-0",route="prometheus-k8s",server="10.134.173.51:9091",service="prometheus-k8s"} 0
haproxy_server_max_session_rate{namespace="openshift-monitoring",pod="prometheus-k8s-0",route="prometheus-k8s-federate",server="10.134.173.51:9091",service="prometheus-k8s"} 0
haproxy_server_max_session_rate{namespace="openshift-monitoring",pod="prometheus-k8s-1",route="prometheus-k8s",server="10.134.249.153:9091",service="prometheus-k8s"} 0
haproxy_server_max_session_rate{namespace="openshift-monitoring",pod="prometheus-k8s-1",route="prometheus-k8s-federate",server="10.134.249.153:9091",service="prometheus-k8s"} 0
haproxy_server_max_session_rate{namespace="openshift-monitoring",pod="thanos-querier-59d7db9787-hwjzt",route="thanos-querier",server="10.134.249.65:9091",service="thanos-querier"} 0
haproxy_server_max_session_rate{namespace="openshift-monitoring",pod="thanos-querier-59d7db9787-p8m9l",route="thanos-querier",server="10.134.173.112:9091",service="thanos-querier"} 0
haproxy_server_max_session_rate{namespace="profiler",pod="esc-static-service-654f8bc589-xchjk",route="esc-static-service-j9sgv",server="10.134.173.123:8080",service="esc-static-service"} 0
haproxy_server_max_session_rate{namespace="profiler",pod="esc-test-service-75849d7c6c-czvsw",route="esc-test-service-lbp8r",server="10.134.172.182:8080",service="esc-test-service"} 0
haproxy_server_max_session_rate{namespace="profiler",pod="esc-ui-service-66977c9bcb-kg9mc",route="esc-ui-service-r5tfm",server="10.134.249.215:8180",service="esc-ui-service"} 0
haproxy_server_max_session_rate{namespace="prometheus-operator",pod="grafana-deployment-88d46bc5f-5gx7h",route="prometheus-operator-grafana",server="10.132.149.3:3000",service="grafana-service"} 0
haproxy_server_max_session_rate{namespace="prometheus-operator",pod="vmagent-k8s-67654d6f97-98nf5",route="prometheus-operator-vmagent-k8s",server="10.132.149.156:8429",service="vmagent-k8s"} 0
haproxy_server_max_session_rate{namespace="prometheus-operator",pod="vmalert-k8s-6bcf966dd6-jvpz5",route="prometheus-operator-vmalert-k8s",server="10.132.149.175:8080",service="vmalert-k8s"} 0
haproxy_server_max_session_rate{namespace="prometheus-operator",pod="vmalertmanager-k8s-0",route="prometheus-operator-vmalertmanager-k8s",server="10.134.173.133:9093",service="vmalertmanager-k8s"} 0
haproxy_server_max_session_rate{namespace="prometheus-operator",pod="vmsingle-k8s-f4844bfc4-5l8sf",route="prometheus-operator-vmsingle-k8s",server="10.134.249.110:8429",service="vmsingle-k8s"} 0
haproxy_server_max_session_rate{namespace="rabbitmq",pod="rmqlocal-0",route="rabbitmq-fnzhp",server="10.134.249.141:15672",service="rabbitmq"} 0
haproxy_server_max_session_rate{namespace="rabbitmq",pod="rmqlocal-1",route="rabbitmq-fnzhp",server="10.134.173.127:15672",service="rabbitmq"} 0
haproxy_server_max_session_rate{namespace="rabbitmq",pod="rmqlocal-2",route="rabbitmq-fnzhp",server="10.132.149.46:15672",service="rabbitmq"} 0
haproxy_server_max_session_rate{namespace="shared-license-distributor",pod="shared-license-distributor-64ccc89449-8kw8p",route="shared-license-distributor-wnj68",server="10.134.249.239:8080",service="shared-license-distributor"} 0
haproxy_server_max_session_rate{namespace="site-manager",pod="site-manager-547546bdd5-hd85n",route="site-manager",server="10.134.248.81:8443",service="site-manager"} 0
haproxy_server_max_session_rate{namespace="spark",pod="spark-history-server-bd755d9f5-9ccsz",route="spark-history-server-hqshb",server="10.132.149.178:18080",service="spark-history-server"} 0
haproxy_server_max_session_rate{namespace="spark",pod="spark-history-server-bd755d9f5-zf2dt",route="spark-history-server-hqshb",server="10.134.173.147:18080",service="spark-history-server"} 0
haproxy_server_max_session_rate{namespace="spark",pod="spark-site-manager-66c44c7b78-cc7jr",route="spark-site-manager-p4nqh",server="10.132.149.148:8080",service="spark-site-manager"} 0
# HELP haproxy_server_max_sessions Maximum observed number of active sessions.
# TYPE haproxy_server_max_sessions gauge
haproxy_server_max_sessions{namespace="",pod="",route="",server="fe_no_sni",service=""} 0
haproxy_server_max_sessions{namespace="",pod="",route="",server="fe_sni",service=""} 0
haproxy_server_max_sessions{namespace="airflow",pod="airflow-site-manager-68ccf87889-bwbdf",route="airflow-site-manager-hpl66",server="10.132.149.80:8080",service="airflow-site-manager"} 0
haproxy_server_max_sessions{namespace="airflow",pod="airflow-webserver-695b5c545b-rxqc4",route="airflow-ingress-nzbtm",server="10.134.173.137:8080",service="airflow-webserver"} 0
haproxy_server_max_sessions{namespace="arangodb",pod="main-arangodb-1-87669b688-q98z4",route="arangodb-leader-ingress-zw82p",server="10.132.149.185:8529",service="arangodb-leader"} 0
haproxy_server_max_sessions{namespace="consul-service",pod="consul-server-0",route="consul-ingress-hmmvf",server="10.134.248.50:8500",service="consul-ui"} 0
haproxy_server_max_sessions{namespace="consul-service",pod="consul-server-1",route="consul-ingress-hmmvf",server="10.132.149.160:8500",service="consul-ui"} 0
haproxy_server_max_sessions{namespace="consul-service",pod="consul-server-2",route="consul-ingress-hmmvf",server="10.134.173.136:8500",service="consul-ui"} 0
haproxy_server_max_sessions{namespace="dbaas",pod="dbaas-aggregator-5f744d74f5-r7tqm",route="aggregator-fw5ft",server="10.134.248.185:8080",service="dbaas-aggregator"} 0
haproxy_server_max_sessions{namespace="infra-keycloak",pod="infra-keycloak-7fd57cd59c-d6c5b",route="infra-keycloak-5q4rj",server="10.132.148.223:8080",service="infra-keycloak"} 0
haproxy_server_max_sessions{namespace="infra-keycloak",pod="infra-keycloak-7fd57cd59c-qtqv9",route="infra-keycloak-5q4rj",server="10.134.173.182:8080",service="infra-keycloak"} 0
haproxy_server_max_sessions{namespace="jaeger",pod="jaeger-hotrod-85b766c49d-s6vz7",route="jaeger-hotrod-qq5cq",server="10.134.249.205:8080",service="jaeger-hotrod"} 0
haproxy_server_max_sessions{namespace="jaeger",pod="jaeger-query-65d9446fc6-jx5qk",route="jaeger-query-pc6w5",server="10.134.248.54:16686",service="jaeger-query"} 0
haproxy_server_max_sessions{namespace="kafka-service",pod="akhq-58744f76db-h2mtf",route="akhq-ingress-lhz5z",server="10.134.249.133:8080",service="akhq"} 0
haproxy_server_max_sessions{namespace="license-server",pod="license-server-7b8b6f765f-wvcx5",route="license-server-9txj9",server="10.132.148.47:8080",service="license-server"} 0
haproxy_server_max_sessions{namespace="opensearch",pod="opensearch-dashboards-764bc8548b-zngcr",route="opensearch-dashboards-l6vdv",server="10.134.173.229:5601",service="opensearch-dashboards"} 0
haproxy_server_max_sessions{namespace="openshift-authentication",pod="oauth-openshift-75965c7487-kc8m7",route="oauth-openshift",server="10.132.153.132:6443",service="oauth-openshift"} 0
haproxy_server_max_sessions{namespace="openshift-authentication",pod="oauth-openshift-75965c7487-nsfxw",route="oauth-openshift",server="10.133.43.151:6443",service="oauth-openshift"} 0
haproxy_server_max_sessions{namespace="openshift-authentication",pod="oauth-openshift-75965c7487-xmmmq",route="oauth-openshift",server="10.135.207.187:6443",service="oauth-openshift"} 0
haproxy_server_max_sessions{namespace="openshift-console",pod="console-66fcc4f4fc-crbz6",route="console",server="10.133.43.162:8443",service="console"} 0
haproxy_server_max_sessions{namespace="openshift-console",pod="console-66fcc4f4fc-j5dvz",route="console",server="10.132.153.185:8443",service="console"} 0
haproxy_server_max_sessions{namespace="openshift-console",pod="downloads-6fb844fb47-cwnvw",route="downloads",server="10.135.207.221:8080",service="downloads"} 0
haproxy_server_max_sessions{namespace="openshift-console",pod="downloads-6fb844fb47-l485s",route="downloads",server="10.132.153.222:8080",service="downloads"} 0
haproxy_server_max_sessions{namespace="openshift-ingress-canary",pod="ingress-canary-7rvzj",route="canary",server="10.134.173.247:8080",service="ingress-canary"} 0
haproxy_server_max_sessions{namespace="openshift-ingress-canary",pod="ingress-canary-lwpkn",route="canary",server="10.132.149.137:8080",service="ingress-canary"} 0
haproxy_server_max_sessions{namespace="openshift-ingress-canary",pod="ingress-canary-z2b4x",route="canary",server="10.134.248.101:8080",service="ingress-canary"} 0
haproxy_server_max_sessions{namespace="openshift-monitoring",pod="alertmanager-main-0",route="alertmanager-main",server="10.134.173.99:9095",service="alertmanager-main"} 0
haproxy_server_max_sessions{namespace="openshift-monitoring",pod="alertmanager-main-1",route="alertmanager-main",server="10.134.248.127:9095",service="alertmanager-main"} 0
haproxy_server_max_sessions{namespace="openshift-monitoring",pod="prometheus-k8s-0",route="prometheus-k8s",server="10.134.173.51:9091",service="prometheus-k8s"} 0
haproxy_server_max_sessions{namespace="openshift-monitoring",pod="prometheus-k8s-0",route="prometheus-k8s-federate",server="10.134.173.51:9091",service="prometheus-k8s"} 0
haproxy_server_max_sessions{namespace="openshift-monitoring",pod="prometheus-k8s-1",route="prometheus-k8s",server="10.134.249.153:9091",service="prometheus-k8s"} 0
haproxy_server_max_sessions{namespace="openshift-monitoring",pod="prometheus-k8s-1",route="prometheus-k8s-federate",server="10.134.249.153:9091",service="prometheus-k8s"} 0
haproxy_server_max_sessions{namespace="openshift-monitoring",pod="thanos-querier-59d7db9787-hwjzt",route="thanos-querier",server="10.134.249.65:9091",service="thanos-querier"} 0
haproxy_server_max_sessions{namespace="openshift-monitoring",pod="thanos-querier-59d7db9787-p8m9l",route="thanos-querier",server="10.134.173.112:9091",service="thanos-querier"} 0
haproxy_server_max_sessions{namespace="profiler",pod="esc-static-service-654f8bc589-xchjk",route="esc-static-service-j9sgv",server="10.134.173.123:8080",service="esc-static-service"} 0
haproxy_server_max_sessions{namespace="profiler",pod="esc-test-service-75849d7c6c-czvsw",route="esc-test-service-lbp8r",server="10.134.172.182:8080",service="esc-test-service"} 0
haproxy_server_max_sessions{namespace="profiler",pod="esc-ui-service-66977c9bcb-kg9mc",route="esc-ui-service-r5tfm",server="10.134.249.215:8180",service="esc-ui-service"} 0
haproxy_server_max_sessions{namespace="prometheus-operator",pod="grafana-deployment-88d46bc5f-5gx7h",route="prometheus-operator-grafana",server="10.132.149.3:3000",service="grafana-service"} 0
haproxy_server_max_sessions{namespace="prometheus-operator",pod="vmagent-k8s-67654d6f97-98nf5",route="prometheus-operator-vmagent-k8s",server="10.132.149.156:8429",service="vmagent-k8s"} 0
haproxy_server_max_sessions{namespace="prometheus-operator",pod="vmalert-k8s-6bcf966dd6-jvpz5",route="prometheus-operator-vmalert-k8s",server="10.132.149.175:8080",service="vmalert-k8s"} 0
haproxy_server_max_sessions{namespace="prometheus-operator",pod="vmalertmanager-k8s-0",route="prometheus-operator-vmalertmanager-k8s",server="10.134.173.133:9093",service="vmalertmanager-k8s"} 0
haproxy_server_max_sessions{namespace="prometheus-operator",pod="vmsingle-k8s-f4844bfc4-5l8sf",route="prometheus-operator-vmsingle-k8s",server="10.134.249.110:8429",service="vmsingle-k8s"} 0
haproxy_server_max_sessions{namespace="rabbitmq",pod="rmqlocal-0",route="rabbitmq-fnzhp",server="10.134.249.141:15672",service="rabbitmq"} 0
haproxy_server_max_sessions{namespace="rabbitmq",pod="rmqlocal-1",route="rabbitmq-fnzhp",server="10.134.173.127:15672",service="rabbitmq"} 0
haproxy_server_max_sessions{namespace="rabbitmq",pod="rmqlocal-2",route="rabbitmq-fnzhp",server="10.132.149.46:15672",service="rabbitmq"} 0
haproxy_server_max_sessions{namespace="shared-license-distributor",pod="shared-license-distributor-64ccc89449-8kw8p",route="shared-license-distributor-wnj68",server="10.134.249.239:8080",service="shared-license-distributor"} 0
haproxy_server_max_sessions{namespace="site-manager",pod="site-manager-547546bdd5-hd85n",route="site-manager",server="10.134.248.81:8443",service="site-manager"} 0
haproxy_server_max_sessions{namespace="spark",pod="spark-history-server-bd755d9f5-9ccsz",route="spark-history-server-hqshb",server="10.132.149.178:18080",service="spark-history-server"} 0
haproxy_server_max_sessions{namespace="spark",pod="spark-history-server-bd755d9f5-zf2dt",route="spark-history-server-hqshb",server="10.134.173.147:18080",service="spark-history-server"} 0
haproxy_server_max_sessions{namespace="spark",pod="spark-site-manager-66c44c7b78-cc7jr",route="spark-site-manager-p4nqh",server="10.132.149.148:8080",service="spark-site-manager"} 0
# HELP haproxy_server_response_errors_total Total of response errors.
# TYPE haproxy_server_response_errors_total gauge
haproxy_server_response_errors_total{namespace="",pod="",route="",server="fe_no_sni",service=""} 0
haproxy_server_response_errors_total{namespace="",pod="",route="",server="fe_sni",service=""} 13
haproxy_server_response_errors_total{namespace="airflow",pod="airflow-site-manager-68ccf87889-bwbdf",route="airflow-site-manager-hpl66",server="10.132.149.80:8080",service="airflow-site-manager"} 0
haproxy_server_response_errors_total{namespace="airflow",pod="airflow-webserver-695b5c545b-rxqc4",route="airflow-ingress-nzbtm",server="10.134.173.137:8080",service="airflow-webserver"} 0
haproxy_server_response_errors_total{namespace="arangodb",pod="main-arangodb-1-87669b688-q98z4",route="arangodb-leader-ingress-zw82p",server="10.132.149.185:8529",service="arangodb-leader"} 0
haproxy_server_response_errors_total{namespace="consul-service",pod="consul-server-0",route="consul-ingress-hmmvf",server="10.134.248.50:8500",service="consul-ui"} 0
haproxy_server_response_errors_total{namespace="consul-service",pod="consul-server-1",route="consul-ingress-hmmvf",server="10.132.149.160:8500",service="consul-ui"} 0
haproxy_server_response_errors_total{namespace="consul-service",pod="consul-server-2",route="consul-ingress-hmmvf",server="10.134.173.136:8500",service="consul-ui"} 0
haproxy_server_response_errors_total{namespace="dbaas",pod="dbaas-aggregator-5f744d74f5-r7tqm",route="aggregator-fw5ft",server="10.134.248.185:8080",service="dbaas-aggregator"} 0
haproxy_server_response_errors_total{namespace="infra-keycloak",pod="infra-keycloak-7fd57cd59c-d6c5b",route="infra-keycloak-5q4rj",server="10.132.148.223:8080",service="infra-keycloak"} 0
haproxy_server_response_errors_total{namespace="infra-keycloak",pod="infra-keycloak-7fd57cd59c-qtqv9",route="infra-keycloak-5q4rj",server="10.134.173.182:8080",service="infra-keycloak"} 0
haproxy_server_response_errors_total{namespace="jaeger",pod="jaeger-hotrod-85b766c49d-s6vz7",route="jaeger-hotrod-qq5cq",server="10.134.249.205:8080",service="jaeger-hotrod"} 0
haproxy_server_response_errors_total{namespace="jaeger",pod="jaeger-query-65d9446fc6-jx5qk",route="jaeger-query-pc6w5",server="10.134.248.54:16686",service="jaeger-query"} 0
haproxy_server_response_errors_total{namespace="kafka-service",pod="akhq-58744f76db-h2mtf",route="akhq-ingress-lhz5z",server="10.134.249.133:8080",service="akhq"} 0
haproxy_server_response_errors_total{namespace="license-server",pod="license-server-7b8b6f765f-wvcx5",route="license-server-9txj9",server="10.132.148.47:8080",service="license-server"} 0
haproxy_server_response_errors_total{namespace="opensearch",pod="opensearch-dashboards-764bc8548b-zngcr",route="opensearch-dashboards-l6vdv",server="10.134.173.229:5601",service="opensearch-dashboards"} 0
haproxy_server_response_errors_total{namespace="openshift-authentication",pod="oauth-openshift-75965c7487-kc8m7",route="oauth-openshift",server="10.132.153.132:6443",service="oauth-openshift"} 0
haproxy_server_response_errors_total{namespace="openshift-authentication",pod="oauth-openshift-75965c7487-nsfxw",route="oauth-openshift",server="10.133.43.151:6443",service="oauth-openshift"} 0
haproxy_server_response_errors_total{namespace="openshift-authentication",pod="oauth-openshift-75965c7487-xmmmq",route="oauth-openshift",server="10.135.207.187:6443",service="oauth-openshift"} 0
haproxy_server_response_errors_total{namespace="openshift-console",pod="console-66fcc4f4fc-crbz6",route="console",server="10.133.43.162:8443",service="console"} 0
haproxy_server_response_errors_total{namespace="openshift-console",pod="console-66fcc4f4fc-j5dvz",route="console",server="10.132.153.185:8443",service="console"} 0
haproxy_server_response_errors_total{namespace="openshift-console",pod="downloads-6fb844fb47-cwnvw",route="downloads",server="10.135.207.221:8080",service="downloads"} 0
haproxy_server_response_errors_total{namespace="openshift-console",pod="downloads-6fb844fb47-l485s",route="downloads",server="10.132.153.222:8080",service="downloads"} 0
haproxy_server_response_errors_total{namespace="openshift-ingress-canary",pod="ingress-canary-7rvzj",route="canary",server="10.134.173.247:8080",service="ingress-canary"} 0
haproxy_server_response_errors_total{namespace="openshift-ingress-canary",pod="ingress-canary-lwpkn",route="canary",server="10.132.149.137:8080",service="ingress-canary"} 0
haproxy_server_response_errors_total{namespace="openshift-ingress-canary",pod="ingress-canary-z2b4x",route="canary",server="10.134.248.101:8080",service="ingress-canary"} 0
haproxy_server_response_errors_total{namespace="openshift-monitoring",pod="alertmanager-main-0",route="alertmanager-main",server="10.134.173.99:9095",service="alertmanager-main"} 0
haproxy_server_response_errors_total{namespace="openshift-monitoring",pod="alertmanager-main-1",route="alertmanager-main",server="10.134.248.127:9095",service="alertmanager-main"} 0
haproxy_server_response_errors_total{namespace="openshift-monitoring",pod="prometheus-k8s-0",route="prometheus-k8s",server="10.134.173.51:9091",service="prometheus-k8s"} 0
haproxy_server_response_errors_total{namespace="openshift-monitoring",pod="prometheus-k8s-0",route="prometheus-k8s-federate",server="10.134.173.51:9091",service="prometheus-k8s"} 0
haproxy_server_response_errors_total{namespace="openshift-monitoring",pod="prometheus-k8s-1",route="prometheus-k8s",server="10.134.249.153:9091",service="prometheus-k8s"} 0
haproxy_server_response_errors_total{namespace="openshift-monitoring",pod="prometheus-k8s-1",route="prometheus-k8s-federate",server="10.134.249.153:9091",service="prometheus-k8s"} 0
haproxy_server_response_errors_total{namespace="openshift-monitoring",pod="thanos-querier-59d7db9787-hwjzt",route="thanos-querier",server="10.134.249.65:9091",service="thanos-querier"} 0
haproxy_server_response_errors_total{namespace="openshift-monitoring",pod="thanos-querier-59d7db9787-p8m9l",route="thanos-querier",server="10.134.173.112:9091",service="thanos-querier"} 0
haproxy_server_response_errors_total{namespace="profiler",pod="esc-static-service-654f8bc589-xchjk",route="esc-static-service-j9sgv",server="10.134.173.123:8080",service="esc-static-service"} 0
haproxy_server_response_errors_total{namespace="profiler",pod="esc-test-service-75849d7c6c-czvsw",route="esc-test-service-lbp8r",server="10.134.172.182:8080",service="esc-test-service"} 0
haproxy_server_response_errors_total{namespace="profiler",pod="esc-ui-service-66977c9bcb-kg9mc",route="esc-ui-service-r5tfm",server="10.134.249.215:8180",service="esc-ui-service"} 0
haproxy_server_response_errors_total{namespace="prometheus-operator",pod="grafana-deployment-88d46bc5f-5gx7h",route="prometheus-operator-grafana",server="10.132.149.3:3000",service="grafana-service"} 0
haproxy_server_response_errors_total{namespace="prometheus-operator",pod="vmagent-k8s-67654d6f97-98nf5",route="prometheus-operator-vmagent-k8s",server="10.132.149.156:8429",service="vmagent-k8s"} 0
haproxy_server_response_errors_total{namespace="prometheus-operator",pod="vmalert-k8s-6bcf966dd6-jvpz5",route="prometheus-operator-vmalert-k8s",server="10.132.149.175:8080",service="vmalert-k8s"} 0
haproxy_server_response_errors_total{namespace="prometheus-operator",pod="vmalertmanager-k8s-0",route="prometheus-operator-vmalertmanager-k8s",server="10.134.173.133:9093",service="vmalertmanager-k8s"} 0
haproxy_server_response_errors_total{namespace="prometheus-operator",pod="vmsingle-k8s-f4844bfc4-5l8sf",route="prometheus-operator-vmsingle-k8s",server="10.134.249.110:8429",service="vmsingle-k8s"} 0
haproxy_server_response_errors_total{namespace="rabbitmq",pod="rmqlocal-0",route="rabbitmq-fnzhp",server="10.134.249.141:15672",service="rabbitmq"} 0
haproxy_server_response_errors_total{namespace="rabbitmq",pod="rmqlocal-1",route="rabbitmq-fnzhp",server="10.134.173.127:15672",service="rabbitmq"} 0
haproxy_server_response_errors_total{namespace="rabbitmq",pod="rmqlocal-2",route="rabbitmq-fnzhp",server="10.132.149.46:15672",service="rabbitmq"} 0
haproxy_server_response_errors_total{namespace="shared-license-distributor",pod="shared-license-distributor-64ccc89449-8kw8p",route="shared-license-distributor-wnj68",server="10.134.249.239:8080",service="shared-license-distributor"} 0
haproxy_server_response_errors_total{namespace="site-manager",pod="site-manager-547546bdd5-hd85n",route="site-manager",server="10.134.248.81:8443",service="site-manager"} 0
haproxy_server_response_errors_total{namespace="spark",pod="spark-history-server-bd755d9f5-9ccsz",route="spark-history-server-hqshb",server="10.132.149.178:18080",service="spark-history-server"} 0
haproxy_server_response_errors_total{namespace="spark",pod="spark-history-server-bd755d9f5-zf2dt",route="spark-history-server-hqshb",server="10.134.173.147:18080",service="spark-history-server"} 0
haproxy_server_response_errors_total{namespace="spark",pod="spark-site-manager-66c44c7b78-cc7jr",route="spark-site-manager-p4nqh",server="10.132.149.148:8080",service="spark-site-manager"} 0
# HELP haproxy_server_up Current health status of the server (1 = UP, 0 = DOWN).
# TYPE haproxy_server_up gauge
haproxy_server_up{namespace="",pod="",route="",server="fe_no_sni",service=""} 1
haproxy_server_up{namespace="",pod="",route="",server="fe_sni",service=""} 1
haproxy_server_up{namespace="airflow",pod="airflow-site-manager-68ccf87889-bwbdf",route="airflow-site-manager-hpl66",server="10.132.149.80:8080",service="airflow-site-manager"} 1
haproxy_server_up{namespace="airflow",pod="airflow-webserver-695b5c545b-rxqc4",route="airflow-ingress-nzbtm",server="10.134.173.137:8080",service="airflow-webserver"} 1
haproxy_server_up{namespace="arangodb",pod="main-arangodb-1-87669b688-q98z4",route="arangodb-leader-ingress-zw82p",server="10.132.149.185:8529",service="arangodb-leader"} 1
haproxy_server_up{namespace="consul-service",pod="consul-server-0",route="consul-ingress-hmmvf",server="10.134.248.50:8500",service="consul-ui"} 1
haproxy_server_up{namespace="consul-service",pod="consul-server-1",route="consul-ingress-hmmvf",server="10.132.149.160:8500",service="consul-ui"} 1
haproxy_server_up{namespace="consul-service",pod="consul-server-2",route="consul-ingress-hmmvf",server="10.134.173.136:8500",service="consul-ui"} 1
haproxy_server_up{namespace="dbaas",pod="dbaas-aggregator-5f744d74f5-r7tqm",route="aggregator-fw5ft",server="10.134.248.185:8080",service="dbaas-aggregator"} 1
haproxy_server_up{namespace="infra-keycloak",pod="infra-keycloak-7fd57cd59c-d6c5b",route="infra-keycloak-5q4rj",server="10.132.148.223:8080",service="infra-keycloak"} 1
haproxy_server_up{namespace="infra-keycloak",pod="infra-keycloak-7fd57cd59c-qtqv9",route="infra-keycloak-5q4rj",server="10.134.173.182:8080",service="infra-keycloak"} 1
haproxy_server_up{namespace="jaeger",pod="jaeger-hotrod-85b766c49d-s6vz7",route="jaeger-hotrod-qq5cq",server="10.134.249.205:8080",service="jaeger-hotrod"} 1
haproxy_server_up{namespace="jaeger",pod="jaeger-query-65d9446fc6-jx5qk",route="jaeger-query-pc6w5",server="10.134.248.54:16686",service="jaeger-query"} 1
haproxy_server_up{namespace="kafka-service",pod="akhq-58744f76db-h2mtf",route="akhq-ingress-lhz5z",server="10.134.249.133:8080",service="akhq"} 1
haproxy_server_up{namespace="license-server",pod="license-server-7b8b6f765f-wvcx5",route="license-server-9txj9",server="10.132.148.47:8080",service="license-server"} 1
haproxy_server_up{namespace="opensearch",pod="opensearch-dashboards-764bc8548b-zngcr",route="opensearch-dashboards-l6vdv",server="10.134.173.229:5601",service="opensearch-dashboards"} 1
haproxy_server_up{namespace="openshift-authentication",pod="oauth-openshift-75965c7487-kc8m7",route="oauth-openshift",server="10.132.153.132:6443",service="oauth-openshift"} 1
haproxy_server_up{namespace="openshift-authentication",pod="oauth-openshift-75965c7487-nsfxw",route="oauth-openshift",server="10.133.43.151:6443",service="oauth-openshift"} 1
haproxy_server_up{namespace="openshift-authentication",pod="oauth-openshift-75965c7487-xmmmq",route="oauth-openshift",server="10.135.207.187:6443",service="oauth-openshift"} 1
haproxy_server_up{namespace="openshift-console",pod="console-66fcc4f4fc-crbz6",route="console",server="10.133.43.162:8443",service="console"} 1
haproxy_server_up{namespace="openshift-console",pod="console-66fcc4f4fc-j5dvz",route="console",server="10.132.153.185:8443",service="console"} 1
haproxy_server_up{namespace="openshift-console",pod="downloads-6fb844fb47-cwnvw",route="downloads",server="10.135.207.221:8080",service="downloads"} 1
haproxy_server_up{namespace="openshift-console",pod="downloads-6fb844fb47-l485s",route="downloads",server="10.132.153.222:8080",service="downloads"} 1
haproxy_server_up{namespace="openshift-ingress-canary",pod="ingress-canary-7rvzj",route="canary",server="10.134.173.247:8080",service="ingress-canary"} 1
haproxy_server_up{namespace="openshift-ingress-canary",pod="ingress-canary-lwpkn",route="canary",server="10.132.149.137:8080",service="ingress-canary"} 1
haproxy_server_up{namespace="openshift-ingress-canary",pod="ingress-canary-z2b4x",route="canary",server="10.134.248.101:8080",service="ingress-canary"} 1
haproxy_server_up{namespace="openshift-monitoring",pod="alertmanager-main-0",route="alertmanager-main",server="10.134.173.99:9095",service="alertmanager-main"} 1
haproxy_server_up{namespace="openshift-monitoring",pod="alertmanager-main-1",route="alertmanager-main",server="10.134.248.127:9095",service="alertmanager-main"} 1
haproxy_server_up{namespace="openshift-monitoring",pod="prometheus-k8s-0",route="prometheus-k8s",server="10.134.173.51:9091",service="prometheus-k8s"} 1
haproxy_server_up{namespace="openshift-monitoring",pod="prometheus-k8s-0",route="prometheus-k8s-federate",server="10.134.173.51:9091",service="prometheus-k8s"} 1
haproxy_server_up{namespace="openshift-monitoring",pod="prometheus-k8s-1",route="prometheus-k8s",server="10.134.249.153:9091",service="prometheus-k8s"} 1
haproxy_server_up{namespace="openshift-monitoring",pod="prometheus-k8s-1",route="prometheus-k8s-federate",server="10.134.249.153:9091",service="prometheus-k8s"} 1
haproxy_server_up{namespace="openshift-monitoring",pod="thanos-querier-59d7db9787-hwjzt",route="thanos-querier",server="10.134.249.65:9091",service="thanos-querier"} 1
haproxy_server_up{namespace="openshift-monitoring",pod="thanos-querier-59d7db9787-p8m9l",route="thanos-querier",server="10.134.173.112:9091",service="thanos-querier"} 1
haproxy_server_up{namespace="profiler",pod="esc-static-service-654f8bc589-xchjk",route="esc-static-service-j9sgv",server="10.134.173.123:8080",service="esc-static-service"} 1
haproxy_server_up{namespace="profiler",pod="esc-test-service-75849d7c6c-czvsw",route="esc-test-service-lbp8r",server="10.134.172.182:8080",service="esc-test-service"} 1
haproxy_server_up{namespace="profiler",pod="esc-ui-service-66977c9bcb-kg9mc",route="esc-ui-service-r5tfm",server="10.134.249.215:8180",service="esc-ui-service"} 1
haproxy_server_up{namespace="prometheus-operator",pod="grafana-deployment-88d46bc5f-5gx7h",route="prometheus-operator-grafana",server="10.132.149.3:3000",service="grafana-service"} 1
haproxy_server_up{namespace="prometheus-operator",pod="vmagent-k8s-67654d6f97-98nf5",route="prometheus-operator-vmagent-k8s",server="10.132.149.156:8429",service="vmagent-k8s"} 1
haproxy_server_up{namespace="prometheus-operator",pod="vmalert-k8s-6bcf966dd6-jvpz5",route="prometheus-operator-vmalert-k8s",server="10.132.149.175:8080",service="vmalert-k8s"} 1
haproxy_server_up{namespace="prometheus-operator",pod="vmalertmanager-k8s-0",route="prometheus-operator-vmalertmanager-k8s",server="10.134.173.133:9093",service="vmalertmanager-k8s"} 1
haproxy_server_up{namespace="prometheus-operator",pod="vmsingle-k8s-f4844bfc4-5l8sf",route="prometheus-operator-vmsingle-k8s",server="10.134.249.110:8429",service="vmsingle-k8s"} 1
haproxy_server_up{namespace="rabbitmq",pod="rmqlocal-0",route="rabbitmq-fnzhp",server="10.134.249.141:15672",service="rabbitmq"} 1
haproxy_server_up{namespace="rabbitmq",pod="rmqlocal-1",route="rabbitmq-fnzhp",server="10.134.173.127:15672",service="rabbitmq"} 1
haproxy_server_up{namespace="rabbitmq",pod="rmqlocal-2",route="rabbitmq-fnzhp",server="10.132.149.46:15672",service="rabbitmq"} 1
haproxy_server_up{namespace="shared-license-distributor",pod="shared-license-distributor-64ccc89449-8kw8p",route="shared-license-distributor-wnj68",server="10.134.249.239:8080",service="shared-license-distributor"} 1
haproxy_server_up{namespace="site-manager",pod="site-manager-547546bdd5-hd85n",route="site-manager",server="10.134.248.81:8443",service="site-manager"} 1
haproxy_server_up{namespace="spark",pod="spark-history-server-bd755d9f5-9ccsz",route="spark-history-server-hqshb",server="10.132.149.178:18080",service="spark-history-server"} 1
haproxy_server_up{namespace="spark",pod="spark-history-server-bd755d9f5-zf2dt",route="spark-history-server-hqshb",server="10.134.173.147:18080",service="spark-history-server"} 1
haproxy_server_up{namespace="spark",pod="spark-site-manager-66c44c7b78-cc7jr",route="spark-site-manager-p4nqh",server="10.132.149.148:8080",service="spark-site-manager"} 1
# HELP haproxy_up Was the last scrape of haproxy successful.
# TYPE haproxy_up gauge
haproxy_up 1
# HELP process_cpu_seconds_total Total user and system CPU time spent in seconds.
# TYPE process_cpu_seconds_total counter
process_cpu_seconds_total 2132.14
# HELP process_max_fds Maximum number of open file descriptors.
# TYPE process_max_fds gauge
process_max_fds 1.048576e+06
# HELP process_open_fds Number of open file descriptors.
# TYPE process_open_fds gauge
process_open_fds 14
# HELP process_resident_memory_bytes Resident memory size in bytes.
# TYPE process_resident_memory_bytes gauge
process_resident_memory_bytes 6.1546496e+07
# HELP process_start_time_seconds Start time of the process since unix epoch in seconds.
# TYPE process_start_time_seconds gauge
process_start_time_seconds 1.70263821541e+09
# HELP process_virtual_memory_bytes Virtual memory size in bytes.
# TYPE process_virtual_memory_bytes gauge
process_virtual_memory_bytes 1.86044416e+09
# HELP process_virtual_memory_max_bytes Maximum amount of virtual memory available in bytes.
# TYPE process_virtual_memory_max_bytes gauge
process_virtual_memory_max_bytes 1.8446744073709552e+19
# HELP promhttp_metric_handler_requests_in_flight Current number of scrapes being served.
# TYPE promhttp_metric_handler_requests_in_flight gauge
promhttp_metric_handler_requests_in_flight 1
# HELP promhttp_metric_handler_requests_total Total number of scrapes by HTTP status code.
# TYPE promhttp_metric_handler_requests_total counter
promhttp_metric_handler_requests_total{code="200"} 72325
promhttp_metric_handler_requests_total{code="500"} 0
promhttp_metric_handler_requests_total{code="503"} 0
# HELP template_router_reload_failure Metric to track the status of the most recent HAProxy reload
# TYPE template_router_reload_failure gauge
template_router_reload_failure 0
# HELP template_router_reload_seconds Measures the time spent reloading the router in seconds.
# TYPE template_router_reload_seconds summary
template_router_reload_seconds_sum 11.640913091000009
template_router_reload_seconds_count 156
# HELP template_router_write_config_seconds Measures the time spent writing out the router configuration to disk in seconds.
# TYPE template_router_write_config_seconds summary
template_router_write_config_seconds_sum 2.336337956999999
template_router_write_config_seconds_count 156
```
