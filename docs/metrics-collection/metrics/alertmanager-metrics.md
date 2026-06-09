This document describes the metrics list and how to collect them from Alertmanager.

# Metrics

Alertmanager already exposes its metrics in Prometheus format and doesn't require to use of specific exporters.

| Name       | Metrics Port | Metrics Endpoint | Need Exporter? | Auth? | Is Exporter Third Party? |
| ---------- | ------------ | ---------------- | -------------- | ----- | ------------------------ |
| Prometheus | `9093`       | `/metrics`       | No             | No    | N/A                      |

## How to Collect

Metrics expose on port `9093` and endpoint `/metrics`. By default, Alertmanager has no authentication.

Config `PodMonitor` for `prometheus-operator` to collect Alertmanager metrics:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: PodMonitor
metadata:
  labels:
    app.kubernetes.io/component: monitoring
    app.kubernetes.io/managed-by: monitoring-operator
    app.kubernetes.io/name: alertmanager-pod-monitor
    k8s-app: alertmanager-pod-monitor
  name: alertmanager-pod-monitor
spec:
  podMetricsEndpoints:
  - interval: 30s
    port: web
    scheme: http
  jobLabel: k8s-app
  namespaceSelector: {}
  selector:
    matchLabels:
      app.kubernetes.io/name: alertmanager
```

To collect (or just to check) metrics manually you can use the following command:

```bash
curl -v -k -L http://<alertmanager_ip_or_dns>:9093/metrics
```

or

```bash
wget -O - http://<alertmanager_ip_or_dns>:9093/metrics
```

## Metrics List

```prometheus
# HELP alertmanager_alerts How many alerts by state.
# TYPE alertmanager_alerts gauge
alertmanager_alerts{state="active"} 110
alertmanager_alerts{state="suppressed"} 0
# HELP alertmanager_alerts_invalid_total The total number of received alerts that were invalid.
# TYPE alertmanager_alerts_invalid_total counter
alertmanager_alerts_invalid_total{version="v1"} 0
alertmanager_alerts_invalid_total{version="v2"} 0
# HELP alertmanager_alerts_received_total The total number of received alerts.
# TYPE alertmanager_alerts_received_total counter
alertmanager_alerts_received_total{status="firing",version="v1"} 0
alertmanager_alerts_received_total{status="firing",version="v2"} 52501
alertmanager_alerts_received_total{status="resolved",version="v1"} 0
alertmanager_alerts_received_total{status="resolved",version="v2"} 503
# HELP alertmanager_build_info A metric with a constant '1' value labeled by version, revision, branch, and goversion from which alertmanager was built.
# TYPE alertmanager_build_info gauge
alertmanager_build_info{branch="HEAD",goversion="go1.17.8",revision="f484b17fa3c583ed1b2c8bbcec20ba1db2aa5f11",version="0.24.0"} 1
# HELP alertmanager_cluster_enabled Indicates whether the clustering is enabled or not.
# TYPE alertmanager_cluster_enabled gauge
alertmanager_cluster_enabled 0
# HELP alertmanager_config_hash Hash of the currently loaded alertmanager configuration.
# TYPE alertmanager_config_hash gauge
alertmanager_config_hash 1.50698929575612e+14
# HELP alertmanager_config_last_reload_success_timestamp_seconds Timestamp of the last successful configuration reload.
# TYPE alertmanager_config_last_reload_success_timestamp_seconds gauge
alertmanager_config_last_reload_success_timestamp_seconds 1.6726013963134277e+09
# HELP alertmanager_config_last_reload_successful Whether the last configuration reload attempt was successful.
# TYPE alertmanager_config_last_reload_successful gauge
alertmanager_config_last_reload_successful 1
# HELP alertmanager_dispatcher_aggregation_groups Number of active aggregation groups
# TYPE alertmanager_dispatcher_aggregation_groups gauge
alertmanager_dispatcher_aggregation_groups 1
# HELP alertmanager_dispatcher_alert_processing_duration_seconds Summary of latencies for the processing of alerts.
# TYPE alertmanager_dispatcher_alert_processing_duration_seconds summary
alertmanager_dispatcher_alert_processing_duration_seconds_sum 0.6176038409999952
alertmanager_dispatcher_alert_processing_duration_seconds_count 53004
# HELP alertmanager_http_concurrency_limit_exceeded_total Total number of times an HTTP request failed because the concurrency limit was reached.
# TYPE alertmanager_http_concurrency_limit_exceeded_total counter
alertmanager_http_concurrency_limit_exceeded_total{method="get"} 0
# HELP alertmanager_http_request_duration_seconds Histogram of latencies for HTTP requests.
# TYPE alertmanager_http_request_duration_seconds histogram
alertmanager_http_request_duration_seconds_bucket{handler="/-/healthy",method="get",le="0.05"} 4570
alertmanager_http_request_duration_seconds_bucket{handler="/-/healthy",method="get",le="0.1"} 4570
alertmanager_http_request_duration_seconds_bucket{handler="/-/healthy",method="get",le="0.25"} 4570
alertmanager_http_request_duration_seconds_bucket{handler="/-/healthy",method="get",le="0.5"} 4570
alertmanager_http_request_duration_seconds_bucket{handler="/-/healthy",method="get",le="0.75"} 4570
alertmanager_http_request_duration_seconds_bucket{handler="/-/healthy",method="get",le="1"} 4570
alertmanager_http_request_duration_seconds_bucket{handler="/-/healthy",method="get",le="2"} 4570
alertmanager_http_request_duration_seconds_bucket{handler="/-/healthy",method="get",le="5"} 4570
alertmanager_http_request_duration_seconds_bucket{handler="/-/healthy",method="get",le="20"} 4570
alertmanager_http_request_duration_seconds_bucket{handler="/-/healthy",method="get",le="60"} 4570
alertmanager_http_request_duration_seconds_bucket{handler="/-/healthy",method="get",le="+Inf"} 4570
alertmanager_http_request_duration_seconds_sum{handler="/-/healthy",method="get"} 0.09827384700000014
alertmanager_http_request_duration_seconds_count{handler="/-/healthy",method="get"} 4570
alertmanager_http_request_duration_seconds_bucket{handler="/-/ready",method="get",le="0.05"} 9140
alertmanager_http_request_duration_seconds_bucket{handler="/-/ready",method="get",le="0.1"} 9140
alertmanager_http_request_duration_seconds_bucket{handler="/-/ready",method="get",le="0.25"} 9140
alertmanager_http_request_duration_seconds_bucket{handler="/-/ready",method="get",le="0.5"} 9140
alertmanager_http_request_duration_seconds_bucket{handler="/-/ready",method="get",le="0.75"} 9140
alertmanager_http_request_duration_seconds_bucket{handler="/-/ready",method="get",le="1"} 9140
alertmanager_http_request_duration_seconds_bucket{handler="/-/ready",method="get",le="2"} 9140
alertmanager_http_request_duration_seconds_bucket{handler="/-/ready",method="get",le="5"} 9140
alertmanager_http_request_duration_seconds_bucket{handler="/-/ready",method="get",le="20"} 9140
alertmanager_http_request_duration_seconds_bucket{handler="/-/ready",method="get",le="60"} 9140
alertmanager_http_request_duration_seconds_bucket{handler="/-/ready",method="get",le="+Inf"} 9140
alertmanager_http_request_duration_seconds_sum{handler="/-/ready",method="get"} 0.19480094899999995
alertmanager_http_request_duration_seconds_count{handler="/-/ready",method="get"} 9140
alertmanager_http_request_duration_seconds_bucket{handler="/-/reload",method="post",le="0.05"} 1
alertmanager_http_request_duration_seconds_bucket{handler="/-/reload",method="post",le="0.1"} 1
alertmanager_http_request_duration_seconds_bucket{handler="/-/reload",method="post",le="0.25"} 1
alertmanager_http_request_duration_seconds_bucket{handler="/-/reload",method="post",le="0.5"} 1
alertmanager_http_request_duration_seconds_bucket{handler="/-/reload",method="post",le="0.75"} 1
alertmanager_http_request_duration_seconds_bucket{handler="/-/reload",method="post",le="1"} 1
alertmanager_http_request_duration_seconds_bucket{handler="/-/reload",method="post",le="2"} 1
alertmanager_http_request_duration_seconds_bucket{handler="/-/reload",method="post",le="5"} 1
alertmanager_http_request_duration_seconds_bucket{handler="/-/reload",method="post",le="20"} 1
alertmanager_http_request_duration_seconds_bucket{handler="/-/reload",method="post",le="60"} 1
alertmanager_http_request_duration_seconds_bucket{handler="/-/reload",method="post",le="+Inf"} 1
alertmanager_http_request_duration_seconds_sum{handler="/-/reload",method="post"} 0.004383523
alertmanager_http_request_duration_seconds_count{handler="/-/reload",method="post"} 1
alertmanager_http_request_duration_seconds_bucket{handler="/metrics",method="get",le="0.05"} 1515
alertmanager_http_request_duration_seconds_bucket{handler="/metrics",method="get",le="0.1"} 1518
alertmanager_http_request_duration_seconds_bucket{handler="/metrics",method="get",le="0.25"} 1518
alertmanager_http_request_duration_seconds_bucket{handler="/metrics",method="get",le="0.5"} 1518
alertmanager_http_request_duration_seconds_bucket{handler="/metrics",method="get",le="0.75"} 1518
alertmanager_http_request_duration_seconds_bucket{handler="/metrics",method="get",le="1"} 1518
alertmanager_http_request_duration_seconds_bucket{handler="/metrics",method="get",le="2"} 1518
alertmanager_http_request_duration_seconds_bucket{handler="/metrics",method="get",le="5"} 1518
alertmanager_http_request_duration_seconds_bucket{handler="/metrics",method="get",le="20"} 1518
alertmanager_http_request_duration_seconds_bucket{handler="/metrics",method="get",le="60"} 1518
alertmanager_http_request_duration_seconds_bucket{handler="/metrics",method="get",le="+Inf"} 1518
alertmanager_http_request_duration_seconds_sum{handler="/metrics",method="get"} 6.987856180999997
alertmanager_http_request_duration_seconds_count{handler="/metrics",method="get"} 1518
# HELP alertmanager_http_requests_in_flight Current number of HTTP requests being processed.
# TYPE alertmanager_http_requests_in_flight gauge
alertmanager_http_requests_in_flight{method="get"} 1
# HELP alertmanager_http_response_size_bytes Histogram of response size for HTTP requests.
# TYPE alertmanager_http_response_size_bytes histogram
alertmanager_http_response_size_bytes_bucket{handler="/-/healthy",method="get",le="100"} 4570
alertmanager_http_response_size_bytes_bucket{handler="/-/healthy",method="get",le="1000"} 4570
alertmanager_http_response_size_bytes_bucket{handler="/-/healthy",method="get",le="10000"} 4570
alertmanager_http_response_size_bytes_bucket{handler="/-/healthy",method="get",le="100000"} 4570
alertmanager_http_response_size_bytes_bucket{handler="/-/healthy",method="get",le="1e+06"} 4570
alertmanager_http_response_size_bytes_bucket{handler="/-/healthy",method="get",le="1e+07"} 4570
alertmanager_http_response_size_bytes_bucket{handler="/-/healthy",method="get",le="1e+08"} 4570
alertmanager_http_response_size_bytes_bucket{handler="/-/healthy",method="get",le="+Inf"} 4570
alertmanager_http_response_size_bytes_sum{handler="/-/healthy",method="get"} 9140
alertmanager_http_response_size_bytes_count{handler="/-/healthy",method="get"} 4570
alertmanager_http_response_size_bytes_bucket{handler="/-/ready",method="get",le="100"} 9140
alertmanager_http_response_size_bytes_bucket{handler="/-/ready",method="get",le="1000"} 9140
alertmanager_http_response_size_bytes_bucket{handler="/-/ready",method="get",le="10000"} 9140
alertmanager_http_response_size_bytes_bucket{handler="/-/ready",method="get",le="100000"} 9140
alertmanager_http_response_size_bytes_bucket{handler="/-/ready",method="get",le="1e+06"} 9140
alertmanager_http_response_size_bytes_bucket{handler="/-/ready",method="get",le="1e+07"} 9140
alertmanager_http_response_size_bytes_bucket{handler="/-/ready",method="get",le="1e+08"} 9140
alertmanager_http_response_size_bytes_bucket{handler="/-/ready",method="get",le="+Inf"} 9140
alertmanager_http_response_size_bytes_sum{handler="/-/ready",method="get"} 18280
alertmanager_http_response_size_bytes_count{handler="/-/ready",method="get"} 9140
alertmanager_http_response_size_bytes_bucket{handler="/-/reload",method="post",le="100"} 1
alertmanager_http_response_size_bytes_bucket{handler="/-/reload",method="post",le="1000"} 1
alertmanager_http_response_size_bytes_bucket{handler="/-/reload",method="post",le="10000"} 1
alertmanager_http_response_size_bytes_bucket{handler="/-/reload",method="post",le="100000"} 1
alertmanager_http_response_size_bytes_bucket{handler="/-/reload",method="post",le="1e+06"} 1
alertmanager_http_response_size_bytes_bucket{handler="/-/reload",method="post",le="1e+07"} 1
alertmanager_http_response_size_bytes_bucket{handler="/-/reload",method="post",le="1e+08"} 1
alertmanager_http_response_size_bytes_bucket{handler="/-/reload",method="post",le="+Inf"} 1
alertmanager_http_response_size_bytes_sum{handler="/-/reload",method="post"} 0
alertmanager_http_response_size_bytes_count{handler="/-/reload",method="post"} 1
alertmanager_http_response_size_bytes_bucket{handler="/metrics",method="get",le="100"} 0
alertmanager_http_response_size_bytes_bucket{handler="/metrics",method="get",le="1000"} 0
alertmanager_http_response_size_bytes_bucket{handler="/metrics",method="get",le="10000"} 1518
alertmanager_http_response_size_bytes_bucket{handler="/metrics",method="get",le="100000"} 1518
alertmanager_http_response_size_bytes_bucket{handler="/metrics",method="get",le="1e+06"} 1518
alertmanager_http_response_size_bytes_bucket{handler="/metrics",method="get",le="1e+07"} 1518
alertmanager_http_response_size_bytes_bucket{handler="/metrics",method="get",le="1e+08"} 1518
alertmanager_http_response_size_bytes_bucket{handler="/metrics",method="get",le="+Inf"} 1518
alertmanager_http_response_size_bytes_sum{handler="/metrics",method="get"} 8.917092e+06
alertmanager_http_response_size_bytes_count{handler="/metrics",method="get"} 1518
# HELP alertmanager_integrations Number of configured integrations.
# TYPE alertmanager_integrations gauge
alertmanager_integrations 0
# HELP alertmanager_nflog_gc_duration_seconds Duration of the last notification log garbage collection cycle.
# TYPE alertmanager_nflog_gc_duration_seconds summary
alertmanager_nflog_gc_duration_seconds_sum 6.1342e-05
alertmanager_nflog_gc_duration_seconds_count 50
# HELP alertmanager_nflog_gossip_messages_propagated_total Number of received gossip messages that have been further gossiped.
# TYPE alertmanager_nflog_gossip_messages_propagated_total counter
alertmanager_nflog_gossip_messages_propagated_total 0
# HELP alertmanager_nflog_queries_total Number of notification log queries were received.
# TYPE alertmanager_nflog_queries_total counter
alertmanager_nflog_queries_total 0
# HELP alertmanager_nflog_query_duration_seconds Duration of notification log query evaluation.
# TYPE alertmanager_nflog_query_duration_seconds histogram
alertmanager_nflog_query_duration_seconds_bucket{le="0.005"} 0
alertmanager_nflog_query_duration_seconds_bucket{le="0.01"} 0
alertmanager_nflog_query_duration_seconds_bucket{le="0.025"} 0
alertmanager_nflog_query_duration_seconds_bucket{le="0.05"} 0
alertmanager_nflog_query_duration_seconds_bucket{le="0.1"} 0
alertmanager_nflog_query_duration_seconds_bucket{le="0.25"} 0
alertmanager_nflog_query_duration_seconds_bucket{le="0.5"} 0
alertmanager_nflog_query_duration_seconds_bucket{le="1"} 0
alertmanager_nflog_query_duration_seconds_bucket{le="2.5"} 0
alertmanager_nflog_query_duration_seconds_bucket{le="5"} 0
alertmanager_nflog_query_duration_seconds_bucket{le="10"} 0
alertmanager_nflog_query_duration_seconds_bucket{le="+Inf"} 0
alertmanager_nflog_query_duration_seconds_sum 0
alertmanager_nflog_query_duration_seconds_count 0
# HELP alertmanager_nflog_query_errors_total Number notification log received queries that failed.
# TYPE alertmanager_nflog_query_errors_total counter
alertmanager_nflog_query_errors_total 0
# HELP alertmanager_nflog_snapshot_duration_seconds Duration of the last notification log snapshot.
# TYPE alertmanager_nflog_snapshot_duration_seconds summary
alertmanager_nflog_snapshot_duration_seconds_sum 0.00017687599999999996
alertmanager_nflog_snapshot_duration_seconds_count 50
# HELP alertmanager_nflog_snapshot_size_bytes Size of the last notification log snapshot in bytes.
# TYPE alertmanager_nflog_snapshot_size_bytes gauge
alertmanager_nflog_snapshot_size_bytes 0
# HELP alertmanager_notification_latency_seconds The latency of notifications in seconds.
# TYPE alertmanager_notification_latency_seconds histogram
alertmanager_notification_latency_seconds_bucket{integration="email",le="1"} 0
alertmanager_notification_latency_seconds_bucket{integration="email",le="5"} 0
alertmanager_notification_latency_seconds_bucket{integration="email",le="10"} 0
alertmanager_notification_latency_seconds_bucket{integration="email",le="15"} 0
alertmanager_notification_latency_seconds_bucket{integration="email",le="20"} 0
alertmanager_notification_latency_seconds_bucket{integration="email",le="+Inf"} 0
alertmanager_notification_latency_seconds_sum{integration="email"} 0
alertmanager_notification_latency_seconds_count{integration="email"} 0
alertmanager_notification_latency_seconds_bucket{integration="opsgenie",le="1"} 0
alertmanager_notification_latency_seconds_bucket{integration="opsgenie",le="5"} 0
alertmanager_notification_latency_seconds_bucket{integration="opsgenie",le="10"} 0
alertmanager_notification_latency_seconds_bucket{integration="opsgenie",le="15"} 0
alertmanager_notification_latency_seconds_bucket{integration="opsgenie",le="20"} 0
alertmanager_notification_latency_seconds_bucket{integration="opsgenie",le="+Inf"} 0
alertmanager_notification_latency_seconds_sum{integration="opsgenie"} 0
alertmanager_notification_latency_seconds_count{integration="opsgenie"} 0
alertmanager_notification_latency_seconds_bucket{integration="pagerduty",le="1"} 0
alertmanager_notification_latency_seconds_bucket{integration="pagerduty",le="5"} 0
alertmanager_notification_latency_seconds_bucket{integration="pagerduty",le="10"} 0
alertmanager_notification_latency_seconds_bucket{integration="pagerduty",le="15"} 0
alertmanager_notification_latency_seconds_bucket{integration="pagerduty",le="20"} 0
alertmanager_notification_latency_seconds_bucket{integration="pagerduty",le="+Inf"} 0
alertmanager_notification_latency_seconds_sum{integration="pagerduty"} 0
alertmanager_notification_latency_seconds_count{integration="pagerduty"} 0
alertmanager_notification_latency_seconds_bucket{integration="pushover",le="1"} 0
alertmanager_notification_latency_seconds_bucket{integration="pushover",le="5"} 0
alertmanager_notification_latency_seconds_bucket{integration="pushover",le="10"} 0
alertmanager_notification_latency_seconds_bucket{integration="pushover",le="15"} 0
alertmanager_notification_latency_seconds_bucket{integration="pushover",le="20"} 0
alertmanager_notification_latency_seconds_bucket{integration="pushover",le="+Inf"} 0
alertmanager_notification_latency_seconds_sum{integration="pushover"} 0
alertmanager_notification_latency_seconds_count{integration="pushover"} 0
alertmanager_notification_latency_seconds_bucket{integration="slack",le="1"} 0
alertmanager_notification_latency_seconds_bucket{integration="slack",le="5"} 0
alertmanager_notification_latency_seconds_bucket{integration="slack",le="10"} 0
alertmanager_notification_latency_seconds_bucket{integration="slack",le="15"} 0
alertmanager_notification_latency_seconds_bucket{integration="slack",le="20"} 0
alertmanager_notification_latency_seconds_bucket{integration="slack",le="+Inf"} 0
alertmanager_notification_latency_seconds_sum{integration="slack"} 0
alertmanager_notification_latency_seconds_count{integration="slack"} 0
alertmanager_notification_latency_seconds_bucket{integration="sns",le="1"} 0
alertmanager_notification_latency_seconds_bucket{integration="sns",le="5"} 0
alertmanager_notification_latency_seconds_bucket{integration="sns",le="10"} 0
alertmanager_notification_latency_seconds_bucket{integration="sns",le="15"} 0
alertmanager_notification_latency_seconds_bucket{integration="sns",le="20"} 0
alertmanager_notification_latency_seconds_bucket{integration="sns",le="+Inf"} 0
alertmanager_notification_latency_seconds_sum{integration="sns"} 0
alertmanager_notification_latency_seconds_count{integration="sns"} 0
alertmanager_notification_latency_seconds_bucket{integration="telegram",le="1"} 0
alertmanager_notification_latency_seconds_bucket{integration="telegram",le="5"} 0
alertmanager_notification_latency_seconds_bucket{integration="telegram",le="10"} 0
alertmanager_notification_latency_seconds_bucket{integration="telegram",le="15"} 0
alertmanager_notification_latency_seconds_bucket{integration="telegram",le="20"} 0
alertmanager_notification_latency_seconds_bucket{integration="telegram",le="+Inf"} 0
alertmanager_notification_latency_seconds_sum{integration="telegram"} 0
alertmanager_notification_latency_seconds_count{integration="telegram"} 0
alertmanager_notification_latency_seconds_bucket{integration="victorops",le="1"} 0
alertmanager_notification_latency_seconds_bucket{integration="victorops",le="5"} 0
alertmanager_notification_latency_seconds_bucket{integration="victorops",le="10"} 0
alertmanager_notification_latency_seconds_bucket{integration="victorops",le="15"} 0
alertmanager_notification_latency_seconds_bucket{integration="victorops",le="20"} 0
alertmanager_notification_latency_seconds_bucket{integration="victorops",le="+Inf"} 0
alertmanager_notification_latency_seconds_sum{integration="victorops"} 0
alertmanager_notification_latency_seconds_count{integration="victorops"} 0
alertmanager_notification_latency_seconds_bucket{integration="webhook",le="1"} 0
alertmanager_notification_latency_seconds_bucket{integration="webhook",le="5"} 0
alertmanager_notification_latency_seconds_bucket{integration="webhook",le="10"} 0
alertmanager_notification_latency_seconds_bucket{integration="webhook",le="15"} 0
alertmanager_notification_latency_seconds_bucket{integration="webhook",le="20"} 0
alertmanager_notification_latency_seconds_bucket{integration="webhook",le="+Inf"} 0
alertmanager_notification_latency_seconds_sum{integration="webhook"} 0
alertmanager_notification_latency_seconds_count{integration="webhook"} 0
alertmanager_notification_latency_seconds_bucket{integration="wechat",le="1"} 0
alertmanager_notification_latency_seconds_bucket{integration="wechat",le="5"} 0
alertmanager_notification_latency_seconds_bucket{integration="wechat",le="10"} 0
alertmanager_notification_latency_seconds_bucket{integration="wechat",le="15"} 0
alertmanager_notification_latency_seconds_bucket{integration="wechat",le="20"} 0
alertmanager_notification_latency_seconds_bucket{integration="wechat",le="+Inf"} 0
alertmanager_notification_latency_seconds_sum{integration="wechat"} 0
alertmanager_notification_latency_seconds_count{integration="wechat"} 0
# HELP alertmanager_notification_requests_failed_total The total number of failed notification requests.
# TYPE alertmanager_notification_requests_failed_total counter
alertmanager_notification_requests_failed_total{integration="email"} 0
alertmanager_notification_requests_failed_total{integration="opsgenie"} 0
alertmanager_notification_requests_failed_total{integration="pagerduty"} 0
alertmanager_notification_requests_failed_total{integration="pushover"} 0
alertmanager_notification_requests_failed_total{integration="slack"} 0
alertmanager_notification_requests_failed_total{integration="sns"} 0
alertmanager_notification_requests_failed_total{integration="telegram"} 0
alertmanager_notification_requests_failed_total{integration="victorops"} 0
alertmanager_notification_requests_failed_total{integration="webhook"} 0
alertmanager_notification_requests_failed_total{integration="wechat"} 0
# HELP alertmanager_notification_requests_total The total number of attempted notification requests.
# TYPE alertmanager_notification_requests_total counter
alertmanager_notification_requests_total{integration="email"} 0
alertmanager_notification_requests_total{integration="opsgenie"} 0
alertmanager_notification_requests_total{integration="pagerduty"} 0
alertmanager_notification_requests_total{integration="pushover"} 0
alertmanager_notification_requests_total{integration="slack"} 0
alertmanager_notification_requests_total{integration="sns"} 0
alertmanager_notification_requests_total{integration="telegram"} 0
alertmanager_notification_requests_total{integration="victorops"} 0
alertmanager_notification_requests_total{integration="webhook"} 0
alertmanager_notification_requests_total{integration="wechat"} 0
# HELP alertmanager_notifications_failed_total The total number of failed notifications.
# TYPE alertmanager_notifications_failed_total counter
alertmanager_notifications_failed_total{integration="email"} 0
alertmanager_notifications_failed_total{integration="opsgenie"} 0
alertmanager_notifications_failed_total{integration="pagerduty"} 0
alertmanager_notifications_failed_total{integration="pushover"} 0
alertmanager_notifications_failed_total{integration="slack"} 0
alertmanager_notifications_failed_total{integration="sns"} 0
alertmanager_notifications_failed_total{integration="telegram"} 0
alertmanager_notifications_failed_total{integration="victorops"} 0
alertmanager_notifications_failed_total{integration="webhook"} 0
alertmanager_notifications_failed_total{integration="wechat"} 0
# HELP alertmanager_notifications_total The total number of attempted notifications.
# TYPE alertmanager_notifications_total counter
alertmanager_notifications_total{integration="email"} 0
alertmanager_notifications_total{integration="opsgenie"} 0
alertmanager_notifications_total{integration="pagerduty"} 0
alertmanager_notifications_total{integration="pushover"} 0
alertmanager_notifications_total{integration="slack"} 0
alertmanager_notifications_total{integration="sns"} 0
alertmanager_notifications_total{integration="telegram"} 0
alertmanager_notifications_total{integration="victorops"} 0
alertmanager_notifications_total{integration="webhook"} 0
alertmanager_notifications_total{integration="wechat"} 0
# HELP alertmanager_receivers Number of configured receivers.
# TYPE alertmanager_receivers gauge
alertmanager_receivers 1
# HELP alertmanager_silences How many silences by state.
# TYPE alertmanager_silences gauge
alertmanager_silences{state="active"} 0
alertmanager_silences{state="expired"} 0
alertmanager_silences{state="pending"} 0
# HELP alertmanager_silences_gc_duration_seconds Duration of the last silence garbage collection cycle.
# TYPE alertmanager_silences_gc_duration_seconds summary
alertmanager_silences_gc_duration_seconds_sum 7.882100000000001e-05
alertmanager_silences_gc_duration_seconds_count 50
# HELP alertmanager_silences_gossip_messages_propagated_total Number of received gossip messages that have been further gossiped.
# TYPE alertmanager_silences_gossip_messages_propagated_total counter
alertmanager_silences_gossip_messages_propagated_total 0
# HELP alertmanager_silences_queries_total How many silence queries were received.
# TYPE alertmanager_silences_queries_total counter
alertmanager_silences_queries_total 4765
# HELP alertmanager_silences_query_duration_seconds Duration of silence query evaluation.
# TYPE alertmanager_silences_query_duration_seconds histogram
alertmanager_silences_query_duration_seconds_bucket{le="0.005"} 4767
alertmanager_silences_query_duration_seconds_bucket{le="0.01"} 4767
alertmanager_silences_query_duration_seconds_bucket{le="0.025"} 4767
alertmanager_silences_query_duration_seconds_bucket{le="0.05"} 4767
alertmanager_silences_query_duration_seconds_bucket{le="0.1"} 4767
alertmanager_silences_query_duration_seconds_bucket{le="0.25"} 4767
alertmanager_silences_query_duration_seconds_bucket{le="0.5"} 4767
alertmanager_silences_query_duration_seconds_bucket{le="1"} 4767
alertmanager_silences_query_duration_seconds_bucket{le="2.5"} 4767
alertmanager_silences_query_duration_seconds_bucket{le="5"} 4767
alertmanager_silences_query_duration_seconds_bucket{le="10"} 4767
alertmanager_silences_query_duration_seconds_bucket{le="+Inf"} 4767
alertmanager_silences_query_duration_seconds_sum 0.009317313000000039
alertmanager_silences_query_duration_seconds_count 4767
# HELP alertmanager_silences_query_errors_total How many silence received queries did not succeed.
# TYPE alertmanager_silences_query_errors_total counter
alertmanager_silences_query_errors_total 0
# HELP alertmanager_silences_snapshot_duration_seconds Duration of the last silence snapshot.
# TYPE alertmanager_silences_snapshot_duration_seconds summary
alertmanager_silences_snapshot_duration_seconds_sum 0.000175152
alertmanager_silences_snapshot_duration_seconds_count 50
# HELP alertmanager_silences_snapshot_size_bytes Size of the last silence snapshot in bytes.
# TYPE alertmanager_silences_snapshot_size_bytes gauge
alertmanager_silences_snapshot_size_bytes 0
# HELP go_gc_cycles_automatic_gc_cycles_total Count of completed GC cycles generated by the Go runtime.
# TYPE go_gc_cycles_automatic_gc_cycles_total counter
go_gc_cycles_automatic_gc_cycles_total 381
# HELP go_gc_cycles_forced_gc_cycles_total Count of completed GC cycles forced by the application.
# TYPE go_gc_cycles_forced_gc_cycles_total counter
go_gc_cycles_forced_gc_cycles_total 0
# HELP go_gc_cycles_total_gc_cycles_total Count of all completed GC cycles.
# TYPE go_gc_cycles_total_gc_cycles_total counter
go_gc_cycles_total_gc_cycles_total 381
# HELP go_gc_duration_seconds A summary of the pause duration of garbage collection cycles.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 5.6209e-05
go_gc_duration_seconds{quantile="0.25"} 0.000115146
go_gc_duration_seconds{quantile="0.5"} 0.000164741
go_gc_duration_seconds{quantile="0.75"} 0.000243889
go_gc_duration_seconds{quantile="1"} 0.012260785
go_gc_duration_seconds_sum 0.136202424
go_gc_duration_seconds_count 381
# HELP go_gc_heap_allocs_by_size_bytes_total Distribution of heap allocations by approximate size. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks.
# TYPE go_gc_heap_allocs_by_size_bytes_total histogram
go_gc_heap_allocs_by_size_bytes_total_bucket{le="8.999999999999998"} 315860
go_gc_heap_allocs_by_size_bytes_total_bucket{le="24.999999999999996"} 7.193908e+06
go_gc_heap_allocs_by_size_bytes_total_bucket{le="64.99999999999999"} 9.781051e+06
go_gc_heap_allocs_by_size_bytes_total_bucket{le="144.99999999999997"} 1.1864115e+07
go_gc_heap_allocs_by_size_bytes_total_bucket{le="320.99999999999994"} 1.2794599e+07
go_gc_heap_allocs_by_size_bytes_total_bucket{le="704.9999999999999"} 1.3154449e+07
go_gc_heap_allocs_by_size_bytes_total_bucket{le="1536.9999999999998"} 1.3337969e+07
go_gc_heap_allocs_by_size_bytes_total_bucket{le="3200.9999999999995"} 1.3345581e+07
go_gc_heap_allocs_by_size_bytes_total_bucket{le="6528.999999999999"} 1.3364997e+07
go_gc_heap_allocs_by_size_bytes_total_bucket{le="13568.999999999998"} 1.3374579e+07
go_gc_heap_allocs_by_size_bytes_total_bucket{le="27264.999999999996"} 1.3378925e+07
go_gc_heap_allocs_by_size_bytes_total_bucket{le="+Inf"} 1.3381466e+07
go_gc_heap_allocs_by_size_bytes_total_sum 1.878615392e+09
go_gc_heap_allocs_by_size_bytes_total_count 1.3381466e+07
# HELP go_gc_heap_allocs_bytes_total Cumulative sum of memory allocated to the heap by the application.
# TYPE go_gc_heap_allocs_bytes_total counter
go_gc_heap_allocs_bytes_total 1.878615392e+09
# HELP go_gc_heap_allocs_objects_total Cumulative count of heap allocations triggered by the application. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks.
# TYPE go_gc_heap_allocs_objects_total counter
go_gc_heap_allocs_objects_total 1.3381466e+07
# HELP go_gc_heap_frees_by_size_bytes_total Distribution of freed heap allocations by approximate size. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks.
# TYPE go_gc_heap_frees_by_size_bytes_total histogram
go_gc_heap_frees_by_size_bytes_total_bucket{le="8.999999999999998"} 311622
go_gc_heap_frees_by_size_bytes_total_bucket{le="24.999999999999996"} 7.165815e+06
go_gc_heap_frees_by_size_bytes_total_bucket{le="64.99999999999999"} 9.739541e+06
go_gc_heap_frees_by_size_bytes_total_bucket{le="144.99999999999997"} 1.1809094e+07
go_gc_heap_frees_by_size_bytes_total_bucket{le="320.99999999999994"} 1.2735153e+07
go_gc_heap_frees_by_size_bytes_total_bucket{le="704.9999999999999"} 1.3092536e+07
go_gc_heap_frees_by_size_bytes_total_bucket{le="1536.9999999999998"} 1.3275355e+07
go_gc_heap_frees_by_size_bytes_total_bucket{le="3200.9999999999995"} 1.3282714e+07
go_gc_heap_frees_by_size_bytes_total_bucket{le="6528.999999999999"} 1.3302018e+07
go_gc_heap_frees_by_size_bytes_total_bucket{le="13568.999999999998"} 1.3311533e+07
go_gc_heap_frees_by_size_bytes_total_bucket{le="27264.999999999996"} 1.331584e+07
go_gc_heap_frees_by_size_bytes_total_bucket{le="+Inf"} 1.3318361e+07
go_gc_heap_frees_by_size_bytes_total_sum 1.86770628e+09
go_gc_heap_frees_by_size_bytes_total_count 1.3318361e+07
# HELP go_gc_heap_frees_bytes_total Cumulative sum of heap memory freed by the garbage collector.
# TYPE go_gc_heap_frees_bytes_total counter
go_gc_heap_frees_bytes_total 1.86770628e+09
# HELP go_gc_heap_frees_objects_total Cumulative count of heap allocations whose storage was freed by the garbage collector. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks.
# TYPE go_gc_heap_frees_objects_total counter
go_gc_heap_frees_objects_total 1.3318361e+07
# HELP go_gc_heap_goal_bytes Heap size target for the end of the GC cycle.
# TYPE go_gc_heap_goal_bytes gauge
go_gc_heap_goal_bytes 1.5724144e+07
# HELP go_gc_heap_objects_objects Number of objects, live or unswept, occupying heap memory.
# TYPE go_gc_heap_objects_objects gauge
go_gc_heap_objects_objects 63105
# HELP go_gc_heap_tiny_allocs_objects_total Count of small allocations that are packed together into blocks. These allocations are counted separately from other allocations because each individual allocation is not tracked by the runtime, only their block. Each block is already accounted for in allocs-by-size and frees-by-size.
# TYPE go_gc_heap_tiny_allocs_objects_total counter
go_gc_heap_tiny_allocs_objects_total 1.266107e+06
# HELP go_gc_pauses_seconds_total Distribution individual GC-related stop-the-world pause latencies.
# TYPE go_gc_pauses_seconds_total histogram
go_gc_pauses_seconds_total_bucket{le="-5e-324"} 0
go_gc_pauses_seconds_total_bucket{le="9.999999999999999e-10"} 0
go_gc_pauses_seconds_total_bucket{le="9.999999999999999e-09"} 0
go_gc_pauses_seconds_total_bucket{le="9.999999999999998e-08"} 0
go_gc_pauses_seconds_total_bucket{le="1.0239999999999999e-06"} 0
go_gc_pauses_seconds_total_bucket{le="1.0239999999999999e-05"} 338
go_gc_pauses_seconds_total_bucket{le="0.00010239999999999998"} 467
go_gc_pauses_seconds_total_bucket{le="0.0010485759999999998"} 742
go_gc_pauses_seconds_total_bucket{le="0.010485759999999998"} 760
go_gc_pauses_seconds_total_bucket{le="0.10485759999999998"} 762
go_gc_pauses_seconds_total_bucket{le="+Inf"} 762
go_gc_pauses_seconds_total_sum NaN
go_gc_pauses_seconds_total_count 762
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 19
# HELP go_info Information about the Go environment.
# TYPE go_info gauge
go_info{version="go1.17.8"} 1
# HELP go_memory_classes_heap_free_bytes Memory that is completely free and eligible to be returned to the underlying system, but has not been. This metric is the runtime's estimate of free address space that is backed by physical memory.
# TYPE go_memory_classes_heap_free_bytes gauge
go_memory_classes_heap_free_bytes 2.924544e+06
# HELP go_memory_classes_heap_objects_bytes Memory occupied by live objects and dead objects that have not yet been marked free by the garbage collector.
# TYPE go_memory_classes_heap_objects_bytes gauge
go_memory_classes_heap_objects_bytes 1.0909112e+07
# HELP go_memory_classes_heap_released_bytes Memory that is completely free and has been returned to the underlying system. This metric is the runtime's estimate of free address space that is still mapped into the process, but is not backed by physical memory.
# TYPE go_memory_classes_heap_released_bytes gauge
go_memory_classes_heap_released_bytes 7.872512e+06
# HELP go_memory_classes_heap_stacks_bytes Memory allocated from the heap that is reserved for stack space, whether or not it is currently in-use.
# TYPE go_memory_classes_heap_stacks_bytes gauge
go_memory_classes_heap_stacks_bytes 983040
# HELP go_memory_classes_heap_unused_bytes Memory that is reserved for heap objects but is not currently used to hold heap objects.
# TYPE go_memory_classes_heap_unused_bytes gauge
go_memory_classes_heap_unused_bytes 2.476616e+06
# HELP go_memory_classes_metadata_mcache_free_bytes Memory that is reserved for runtime mcache structures, but not in-use.
# TYPE go_memory_classes_metadata_mcache_free_bytes gauge
go_memory_classes_metadata_mcache_free_bytes 10384
# HELP go_memory_classes_metadata_mcache_inuse_bytes Memory that is occupied by runtime mcache structures that are currently being used.
# TYPE go_memory_classes_metadata_mcache_inuse_bytes gauge
go_memory_classes_metadata_mcache_inuse_bytes 6000
# HELP go_memory_classes_metadata_mspan_free_bytes Memory that is reserved for runtime mspan structures, but not in-use.
# TYPE go_memory_classes_metadata_mspan_free_bytes gauge
go_memory_classes_metadata_mspan_free_bytes 33400
# HELP go_memory_classes_metadata_mspan_inuse_bytes Memory that is occupied by runtime mspan structures that are currently being used.
# TYPE go_memory_classes_metadata_mspan_inuse_bytes gauge
go_memory_classes_metadata_mspan_inuse_bytes 195976
# HELP go_memory_classes_metadata_other_bytes Memory that is reserved for or used to hold runtime metadata.
# TYPE go_memory_classes_metadata_other_bytes gauge
go_memory_classes_metadata_other_bytes 5.710896e+06
# HELP go_memory_classes_os_stacks_bytes Stack memory allocated by the underlying operating system.
# TYPE go_memory_classes_os_stacks_bytes gauge
go_memory_classes_os_stacks_bytes 0
# HELP go_memory_classes_other_bytes Memory used by execution trace buffers, structures for debugging the runtime, finalizer and profiler specials, and more.
# TYPE go_memory_classes_other_bytes gauge
go_memory_classes_other_bytes 930438
# HELP go_memory_classes_profiling_buckets_bytes Memory that is used by the stack trace hash map used for profiling.
# TYPE go_memory_classes_profiling_buckets_bytes gauge
go_memory_classes_profiling_buckets_bytes 1.584466e+06
# HELP go_memory_classes_total_bytes All memory mapped by the Go runtime into the current process as read-write. Note that this does not include memory mapped by code called via cgo or via the syscall package. Sum of all metrics in /memory/classes.
# TYPE go_memory_classes_total_bytes gauge
go_memory_classes_total_bytes 3.3637384e+07
# HELP go_memstats_alloc_bytes Number of bytes allocated and still in use.
# TYPE go_memstats_alloc_bytes gauge
go_memstats_alloc_bytes 1.0909112e+07
# HELP go_memstats_alloc_bytes_total Total number of bytes allocated, even if freed.
# TYPE go_memstats_alloc_bytes_total counter
go_memstats_alloc_bytes_total 1.878615392e+09
# HELP go_memstats_buck_hash_sys_bytes Number of bytes used by the profiling bucket hash table.
# TYPE go_memstats_buck_hash_sys_bytes gauge
go_memstats_buck_hash_sys_bytes 1.584466e+06
# HELP go_memstats_frees_total Total number of frees.
# TYPE go_memstats_frees_total counter
go_memstats_frees_total 1.4584468e+07
# HELP go_memstats_gc_cpu_fraction The fraction of this program's available CPU time used by the GC since the program started.
# TYPE go_memstats_gc_cpu_fraction gauge
go_memstats_gc_cpu_fraction 0
# HELP go_memstats_gc_sys_bytes Number of bytes used for garbage collection system metadata.
# TYPE go_memstats_gc_sys_bytes gauge
go_memstats_gc_sys_bytes 5.710896e+06
# HELP go_memstats_heap_alloc_bytes Number of heap bytes allocated and still in use.
# TYPE go_memstats_heap_alloc_bytes gauge
go_memstats_heap_alloc_bytes 1.0909112e+07
# HELP go_memstats_heap_idle_bytes Number of heap bytes waiting to be used.
# TYPE go_memstats_heap_idle_bytes gauge
go_memstats_heap_idle_bytes 1.0797056e+07
# HELP go_memstats_heap_inuse_bytes Number of heap bytes that are in use.
# TYPE go_memstats_heap_inuse_bytes gauge
go_memstats_heap_inuse_bytes 1.3385728e+07
# HELP go_memstats_heap_objects Number of allocated objects.
# TYPE go_memstats_heap_objects gauge
go_memstats_heap_objects 63105
# HELP go_memstats_heap_released_bytes Number of heap bytes released to OS.
# TYPE go_memstats_heap_released_bytes gauge
go_memstats_heap_released_bytes 7.872512e+06
# HELP go_memstats_heap_sys_bytes Number of heap bytes obtained from system.
# TYPE go_memstats_heap_sys_bytes gauge
go_memstats_heap_sys_bytes 2.4182784e+07
# HELP go_memstats_last_gc_time_seconds Number of seconds since 1970 of last garbage collection.
# TYPE go_memstats_last_gc_time_seconds gauge
go_memstats_last_gc_time_seconds 1.6726470292988563e+09
# HELP go_memstats_lookups_total Total number of pointer lookups.
# TYPE go_memstats_lookups_total counter
go_memstats_lookups_total 0
# HELP go_memstats_mallocs_total Total number of mallocs.
# TYPE go_memstats_mallocs_total counter
go_memstats_mallocs_total 1.4647573e+07
# HELP go_memstats_mcache_inuse_bytes Number of bytes in use by mcache structures.
# TYPE go_memstats_mcache_inuse_bytes gauge
go_memstats_mcache_inuse_bytes 6000
# HELP go_memstats_mcache_sys_bytes Number of bytes used for mcache structures obtained from system.
# TYPE go_memstats_mcache_sys_bytes gauge
go_memstats_mcache_sys_bytes 16384
# HELP go_memstats_mspan_inuse_bytes Number of bytes in use by mspan structures.
# TYPE go_memstats_mspan_inuse_bytes gauge
go_memstats_mspan_inuse_bytes 195976
# HELP go_memstats_mspan_sys_bytes Number of bytes used for mspan structures obtained from system.
# TYPE go_memstats_mspan_sys_bytes gauge
go_memstats_mspan_sys_bytes 229376
# HELP go_memstats_next_gc_bytes Number of heap bytes when next garbage collection will take place.
# TYPE go_memstats_next_gc_bytes gauge
go_memstats_next_gc_bytes 1.5724144e+07
# HELP go_memstats_other_sys_bytes Number of bytes used for other system allocations.
# TYPE go_memstats_other_sys_bytes gauge
go_memstats_other_sys_bytes 930438
# HELP go_memstats_stack_inuse_bytes Number of bytes in use by the stack allocator.
# TYPE go_memstats_stack_inuse_bytes gauge
go_memstats_stack_inuse_bytes 983040
# HELP go_memstats_stack_sys_bytes Number of bytes obtained from system for stack allocator.
# TYPE go_memstats_stack_sys_bytes gauge
go_memstats_stack_sys_bytes 983040
# HELP go_memstats_sys_bytes Number of bytes obtained from system.
# TYPE go_memstats_sys_bytes gauge
go_memstats_sys_bytes 3.3637384e+07
# HELP go_sched_goroutines_goroutines Count of live goroutines.
# TYPE go_sched_goroutines_goroutines gauge
go_sched_goroutines_goroutines 19
# HELP go_sched_latencies_seconds Distribution of the time goroutines have spent in the scheduler in a runnable state before actually running.
# TYPE go_sched_latencies_seconds histogram
go_sched_latencies_seconds_bucket{le="-5e-324"} 0
go_sched_latencies_seconds_bucket{le="9.999999999999999e-10"} 53075
go_sched_latencies_seconds_bucket{le="9.999999999999999e-09"} 53075
go_sched_latencies_seconds_bucket{le="9.999999999999998e-08"} 53130
go_sched_latencies_seconds_bucket{le="1.0239999999999999e-06"} 106674
go_sched_latencies_seconds_bucket{le="1.0239999999999999e-05"} 112346
go_sched_latencies_seconds_bucket{le="0.00010239999999999998"} 122940
go_sched_latencies_seconds_bucket{le="0.0010485759999999998"} 134507
go_sched_latencies_seconds_bucket{le="0.010485759999999998"} 135117
go_sched_latencies_seconds_bucket{le="0.10485759999999998"} 135119
go_sched_latencies_seconds_bucket{le="+Inf"} 135119
go_sched_latencies_seconds_sum NaN
go_sched_latencies_seconds_count 135119
# HELP go_threads Number of OS threads created.
# TYPE go_threads gauge
go_threads 11
# HELP process_cpu_seconds_total Total user and system CPU time spent in seconds.
# TYPE process_cpu_seconds_total counter
process_cpu_seconds_total 51.95
# HELP process_max_fds Maximum number of open file descriptors.
# TYPE process_max_fds gauge
process_max_fds 1.048576e+06
# HELP process_open_fds Number of open file descriptors.
# TYPE process_open_fds gauge
process_open_fds 11
# HELP process_resident_memory_bytes Resident memory size in bytes.
# TYPE process_resident_memory_bytes gauge
process_resident_memory_bytes 3.1952896e+07
# HELP process_start_time_seconds Start time of the process since unix epoch in seconds.
# TYPE process_start_time_seconds gauge
process_start_time_seconds 1.67260139191e+09
# HELP process_virtual_memory_bytes Virtual memory size in bytes.
# TYPE process_virtual_memory_bytes gauge
process_virtual_memory_bytes 7.46610688e+08
# HELP process_virtual_memory_max_bytes Maximum amount of virtual memory available in bytes.
# TYPE process_virtual_memory_max_bytes gauge
process_virtual_memory_max_bytes 1.8446744073709552e+19
# HELP promhttp_metric_handler_requests_in_flight Current number of scrapes being served.
# TYPE promhttp_metric_handler_requests_in_flight gauge
promhttp_metric_handler_requests_in_flight 1
# HELP promhttp_metric_handler_requests_total Total number of scrapes by HTTP status code.
# TYPE promhttp_metric_handler_requests_total counter
promhttp_metric_handler_requests_total{code="200"} 1518
promhttp_metric_handler_requests_total{code="500"} 0
promhttp_metric_handler_requests_total{code="503"} 0
```
