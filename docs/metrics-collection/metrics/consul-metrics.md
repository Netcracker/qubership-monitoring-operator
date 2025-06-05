This document describes the metrics list and how to collect them from Consul.

# Metrics

Consul already can exposes its metrics in Prometheus format and doesn't require to use of specific exporters.

| Name       | Metrics Port | Metrics Endpoint    | Need Exporter? | Auth?          | Is Exporter Third Party? |
| ---------- | ------------ | ------------------- | -------------- | -------------- | ------------------------ |
| Prometheus | `8500`       | `/v1/agent/metrics` | No             | Require, token | N/A                      |

## How to Collect

Metrics expose on port `9093` and endpoint `/metrics`. By default, Consul has authentication by token.

Config `ServiceMonitor` for `prometheus-operator` to collect Consul Server and Client metrics:

`Server` (usually deploy as StatefulSet):

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  labels:
    app.kubernetes.io/component: monitoring
  name: consul-server-service-monitor
spec:
  endpoints:
    - bearerTokenSecret:
        key: token
        name: consul-bootstrap-acl-token
      interval: 30s
      params:
        format:
          - prometheus
      path: /v1/agent/metrics
      port: http
      relabelings:
        - action: replace
          replacement: Server
          targetLabel: role
        - action: replace
          replacement: dc1
          targetLabel: datacenter
      scheme: http
  jobLabel: k8s-app
  selector:
    matchLabels:
      component: server
      name: consul-server
```

`Client` (usually deploy as DaemonSet):

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  labels:
    app.kubernetes.io/component: monitoring
  name: consul-client-service-monitor
spec:
  endpoints:
    - bearerTokenSecret:
        key: token
        name: consul-bootstrap-acl-token
      interval: 30s
      params:
        format:
          - prometheus
      path: /v1/agent/metrics
      port: http
      relabelings:
        - action: replace
          replacement: Client
          targetLabel: role
        - action: replace
          replacement: dc1
          targetLabel: datacenter
      scheme: http
  jobLabel: k8s-app
  selector:
    matchLabels:
      component: client
```

To collect (or just to check) metrics manually you can use the following command:

```bash
curl -v -k -L -H "Authorization: Bearer <token>" "http://<consul_ip_or_dns>:8500/v1/agent/metrics?format=prometheus"
```

Token usually you can find in the Secret with the name `consul-bootstrap-acl-token` in the Consul namespace.

You can't use `wget` because it doesn't allow to add headers for authorization.

## Metrics List

### Consul Server

```prometheus
# HELP consul_client_rpc Telegraf collected metric
# TYPE consul_client_rpc counter
consul_client_rpc{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="counter",project_name="consul-service",role="Server"} 12
# HELP consul_cluster_health_expected_count Telegraf collected metric
# TYPE consul_cluster_health_expected_count untyped
consul_cluster_health_expected_count{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 3
# HELP consul_cluster_health_servers_count Telegraf collected metric
# TYPE consul_cluster_health_servers_count untyped
consul_cluster_health_servers_count{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 3
# HELP consul_cluster_health_status_code Telegraf collected metric
# TYPE consul_cluster_health_status_code untyped
consul_cluster_health_status_code{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
# HELP consul_cluster_servers_is_healthy Telegraf collected metric
# TYPE consul_cluster_servers_is_healthy untyped
consul_cluster_servers_is_healthy{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server",server_name="consul-server-0"} 1
consul_cluster_servers_is_healthy{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server",server_name="consul-server-1"} 1
consul_cluster_servers_is_healthy{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server",server_name="consul-server-2"} 1
# HELP consul_cluster_servers_is_leader Telegraf collected metric
# TYPE consul_cluster_servers_is_leader untyped
consul_cluster_servers_is_leader{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server",server_name="consul-server-0"} 0
consul_cluster_servers_is_leader{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server",server_name="consul-server-1"} 1
consul_cluster_servers_is_leader{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server",server_name="consul-server-2"} 0
# HELP consul_fsm_coordinate_batch_update_90_percentile Telegraf collected metric
# TYPE consul_fsm_coordinate_batch_update_90_percentile untyped
consul_fsm_coordinate_batch_update_90_percentile{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0.123045
# HELP consul_fsm_coordinate_batch_update_count Telegraf collected metric
# TYPE consul_fsm_coordinate_batch_update_count untyped
consul_fsm_coordinate_batch_update_count{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 5
# HELP consul_fsm_coordinate_batch_update_lower Telegraf collected metric
# TYPE consul_fsm_coordinate_batch_update_lower untyped
consul_fsm_coordinate_batch_update_lower{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0.074822
# HELP consul_fsm_coordinate_batch_update_mean Telegraf collected metric
# TYPE consul_fsm_coordinate_batch_update_mean untyped
consul_fsm_coordinate_batch_update_mean{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0.0983936
# HELP consul_fsm_coordinate_batch_update_stddev Telegraf collected metric
# TYPE consul_fsm_coordinate_batch_update_stddev untyped
consul_fsm_coordinate_batch_update_stddev{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0.017036326946850956
# HELP consul_fsm_coordinate_batch_update_sum Telegraf collected metric
# TYPE consul_fsm_coordinate_batch_update_sum untyped
consul_fsm_coordinate_batch_update_sum{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0.491968
# HELP consul_fsm_coordinate_batch_update_upper Telegraf collected metric
# TYPE consul_fsm_coordinate_batch_update_upper untyped
consul_fsm_coordinate_batch_update_upper{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0.123045
# HELP consul_http_GET_v1_operator_autopilot_health_90_percentile Telegraf collected metric
# TYPE consul_http_GET_v1_operator_autopilot_health_90_percentile untyped
consul_http_GET_v1_operator_autopilot_health_90_percentile{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 2.043135
# HELP consul_http_GET_v1_operator_autopilot_health_count Telegraf collected metric
# TYPE consul_http_GET_v1_operator_autopilot_health_count untyped
consul_http_GET_v1_operator_autopilot_health_count{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 1
# HELP consul_http_GET_v1_operator_autopilot_health_lower Telegraf collected metric
# TYPE consul_http_GET_v1_operator_autopilot_health_lower untyped
consul_http_GET_v1_operator_autopilot_health_lower{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 2.043135
# HELP consul_http_GET_v1_operator_autopilot_health_mean Telegraf collected metric
# TYPE consul_http_GET_v1_operator_autopilot_health_mean untyped
consul_http_GET_v1_operator_autopilot_health_mean{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 2.043135
# HELP consul_http_GET_v1_operator_autopilot_health_stddev Telegraf collected metric
# TYPE consul_http_GET_v1_operator_autopilot_health_stddev untyped
consul_http_GET_v1_operator_autopilot_health_stddev{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0
# HELP consul_http_GET_v1_operator_autopilot_health_sum Telegraf collected metric
# TYPE consul_http_GET_v1_operator_autopilot_health_sum untyped
consul_http_GET_v1_operator_autopilot_health_sum{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 2.043135
# HELP consul_http_GET_v1_operator_autopilot_health_upper Telegraf collected metric
# TYPE consul_http_GET_v1_operator_autopilot_health_upper untyped
consul_http_GET_v1_operator_autopilot_health_upper{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 2.043135
# HELP consul_http_GET_v1_status_leader_90_percentile Telegraf collected metric
# TYPE consul_http_GET_v1_status_leader_90_percentile untyped
consul_http_GET_v1_status_leader_90_percentile{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0.148626
# HELP consul_http_GET_v1_status_leader_count Telegraf collected metric
# TYPE consul_http_GET_v1_status_leader_count untyped
consul_http_GET_v1_status_leader_count{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 10
# HELP consul_http_GET_v1_status_leader_lower Telegraf collected metric
# TYPE consul_http_GET_v1_status_leader_lower untyped
consul_http_GET_v1_status_leader_lower{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0.093697
# HELP consul_http_GET_v1_status_leader_mean Telegraf collected metric
# TYPE consul_http_GET_v1_status_leader_mean untyped
consul_http_GET_v1_status_leader_mean{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0.1167306
# HELP consul_http_GET_v1_status_leader_stddev Telegraf collected metric
# TYPE consul_http_GET_v1_status_leader_stddev untyped
consul_http_GET_v1_status_leader_stddev{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0.018598949310108893
# HELP consul_http_GET_v1_status_leader_sum Telegraf collected metric
# TYPE consul_http_GET_v1_status_leader_sum untyped
consul_http_GET_v1_status_leader_sum{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 1.167306
# HELP consul_http_GET_v1_status_leader_upper Telegraf collected metric
# TYPE consul_http_GET_v1_status_leader_upper untyped
consul_http_GET_v1_status_leader_upper{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0.148626
# HELP consul_memberlist_gossip_90_percentile Telegraf collected metric
# TYPE consul_memberlist_gossip_90_percentile untyped
consul_memberlist_gossip_90_percentile{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0.028696
# HELP consul_memberlist_gossip_count Telegraf collected metric
# TYPE consul_memberlist_gossip_count untyped
consul_memberlist_gossip_count{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 210
# HELP consul_memberlist_gossip_lower Telegraf collected metric
# TYPE consul_memberlist_gossip_lower untyped
consul_memberlist_gossip_lower{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0.007545
# HELP consul_memberlist_gossip_mean Telegraf collected metric
# TYPE consul_memberlist_gossip_mean untyped
consul_memberlist_gossip_mean{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0.020084280952380955
# HELP consul_memberlist_gossip_stddev Telegraf collected metric
# TYPE consul_memberlist_gossip_stddev untyped
consul_memberlist_gossip_stddev{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0.011358932219355963
# HELP consul_memberlist_gossip_sum Telegraf collected metric
# TYPE consul_memberlist_gossip_sum untyped
consul_memberlist_gossip_sum{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 4.217698999999999
# HELP consul_memberlist_gossip_upper Telegraf collected metric
# TYPE consul_memberlist_gossip_upper untyped
consul_memberlist_gossip_upper{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0.115167
# HELP consul_memberlist_probeNode_90_percentile Telegraf collected metric
# TYPE consul_memberlist_probeNode_90_percentile untyped
consul_memberlist_probeNode_90_percentile{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 1.643698
# HELP consul_memberlist_probeNode_count Telegraf collected metric
# TYPE consul_memberlist_probeNode_count untyped
consul_memberlist_probeNode_count{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 35
# HELP consul_memberlist_probeNode_lower Telegraf collected metric
# TYPE consul_memberlist_probeNode_lower untyped
consul_memberlist_probeNode_lower{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0.672314
# HELP consul_memberlist_probeNode_mean Telegraf collected metric
# TYPE consul_memberlist_probeNode_mean untyped
consul_memberlist_probeNode_mean{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 1.2349229714285714
# HELP consul_memberlist_probeNode_stddev Telegraf collected metric
# TYPE consul_memberlist_probeNode_stddev untyped
consul_memberlist_probeNode_stddev{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0.7059119420187311
# HELP consul_memberlist_probeNode_sum Telegraf collected metric
# TYPE consul_memberlist_probeNode_sum untyped
consul_memberlist_probeNode_sum{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 43.22230400000001
# HELP consul_memberlist_probeNode_upper Telegraf collected metric
# TYPE consul_memberlist_probeNode_upper untyped
consul_memberlist_probeNode_upper{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 5.010332
# HELP consul_memberlist_pushPullNode_90_percentile Telegraf collected metric
# TYPE consul_memberlist_pushPullNode_90_percentile untyped
consul_memberlist_pushPullNode_90_percentile{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 3.244158
# HELP consul_memberlist_pushPullNode_count Telegraf collected metric
# TYPE consul_memberlist_pushPullNode_count untyped
consul_memberlist_pushPullNode_count{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 1
# HELP consul_memberlist_pushPullNode_lower Telegraf collected metric
# TYPE consul_memberlist_pushPullNode_lower untyped
consul_memberlist_pushPullNode_lower{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 3.244158
# HELP consul_memberlist_pushPullNode_mean Telegraf collected metric
# TYPE consul_memberlist_pushPullNode_mean untyped
consul_memberlist_pushPullNode_mean{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 3.244158
# HELP consul_memberlist_pushPullNode_stddev Telegraf collected metric
# TYPE consul_memberlist_pushPullNode_stddev untyped
consul_memberlist_pushPullNode_stddev{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0
# HELP consul_memberlist_pushPullNode_sum Telegraf collected metric
# TYPE consul_memberlist_pushPullNode_sum untyped
consul_memberlist_pushPullNode_sum{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 3.244158
# HELP consul_memberlist_pushPullNode_upper Telegraf collected metric
# TYPE consul_memberlist_pushPullNode_upper untyped
consul_memberlist_pushPullNode_upper{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 3.244158
# HELP consul_memberlist_tcp_accept Telegraf collected metric
# TYPE consul_memberlist_tcp_accept counter
consul_memberlist_tcp_accept{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="counter",project_name="consul-service",role="Server"} 1
# HELP consul_memberlist_tcp_connect Telegraf collected metric
# TYPE consul_memberlist_tcp_connect counter
consul_memberlist_tcp_connect{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="counter",project_name="consul-service",role="Server"} 1
# HELP consul_memberlist_tcp_sent Telegraf collected metric
# TYPE consul_memberlist_tcp_sent counter
consul_memberlist_tcp_sent{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="counter",project_name="consul-service",role="Server"} 2319
# HELP consul_memberlist_udp_received Telegraf collected metric
# TYPE consul_memberlist_udp_received counter
consul_memberlist_udp_received{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="counter",project_name="consul-service",role="Server"} 8465
# HELP consul_memberlist_udp_sent Telegraf collected metric
# TYPE consul_memberlist_udp_sent counter
consul_memberlist_udp_sent{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="counter",project_name="consul-service",role="Server"} 8449
# HELP consul_raft_applied_index Telegraf collected metric
# TYPE consul_raft_applied_index gauge
consul_raft_applied_index{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="gauge",project_name="consul-service",role="Server"} 94001
# HELP consul_raft_fsm_apply_90_percentile Telegraf collected metric
# TYPE consul_raft_fsm_apply_90_percentile untyped
consul_raft_fsm_apply_90_percentile{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0.449836
# HELP consul_raft_fsm_apply_count Telegraf collected metric
# TYPE consul_raft_fsm_apply_count untyped
consul_raft_fsm_apply_count{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 5
# HELP consul_raft_fsm_apply_lower Telegraf collected metric
# TYPE consul_raft_fsm_apply_lower untyped
consul_raft_fsm_apply_lower{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0.169107
# HELP consul_raft_fsm_apply_mean Telegraf collected metric
# TYPE consul_raft_fsm_apply_mean untyped
consul_raft_fsm_apply_mean{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0.24918040000000002
# HELP consul_raft_fsm_apply_stddev Telegraf collected metric
# TYPE consul_raft_fsm_apply_stddev untyped
consul_raft_fsm_apply_stddev{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0.10423432709928146
# HELP consul_raft_fsm_apply_sum Telegraf collected metric
# TYPE consul_raft_fsm_apply_sum untyped
consul_raft_fsm_apply_sum{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 1.2459019999999998
# HELP consul_raft_fsm_apply_upper Telegraf collected metric
# TYPE consul_raft_fsm_apply_upper untyped
consul_raft_fsm_apply_upper{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0.449836
# HELP consul_raft_last_index Telegraf collected metric
# TYPE consul_raft_last_index gauge
consul_raft_last_index{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="gauge",project_name="consul-service",role="Server"} 94001
# HELP consul_raft_rpc_appendEntries_90_percentile Telegraf collected metric
# TYPE consul_raft_rpc_appendEntries_90_percentile untyped
consul_raft_rpc_appendEntries_90_percentile{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0.013055
# HELP consul_raft_rpc_appendEntries_count Telegraf collected metric
# TYPE consul_raft_rpc_appendEntries_count untyped
consul_raft_rpc_appendEntries_count{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 448
# HELP consul_raft_rpc_appendEntries_lower Telegraf collected metric
# TYPE consul_raft_rpc_appendEntries_lower untyped
consul_raft_rpc_appendEntries_lower{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0.001753
# HELP consul_raft_rpc_appendEntries_mean Telegraf collected metric
# TYPE consul_raft_rpc_appendEntries_mean untyped
consul_raft_rpc_appendEntries_mean{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0.11128883928571424
# HELP consul_raft_rpc_appendEntries_processLogs_90_percentile Telegraf collected metric
# TYPE consul_raft_rpc_appendEntries_processLogs_90_percentile untyped
consul_raft_rpc_appendEntries_processLogs_90_percentile{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0.013267
# HELP consul_raft_rpc_appendEntries_processLogs_count Telegraf collected metric
# TYPE consul_raft_rpc_appendEntries_processLogs_count untyped
consul_raft_rpc_appendEntries_processLogs_count{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 6
# HELP consul_raft_rpc_appendEntries_processLogs_lower Telegraf collected metric
# TYPE consul_raft_rpc_appendEntries_processLogs_lower untyped
consul_raft_rpc_appendEntries_processLogs_lower{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0.008781
# HELP consul_raft_rpc_appendEntries_processLogs_mean Telegraf collected metric
# TYPE consul_raft_rpc_appendEntries_processLogs_mean untyped
consul_raft_rpc_appendEntries_processLogs_mean{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0.010390666666666666
# HELP consul_raft_rpc_appendEntries_processLogs_stddev Telegraf collected metric
# TYPE consul_raft_rpc_appendEntries_processLogs_stddev untyped
consul_raft_rpc_appendEntries_processLogs_stddev{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0.0014072485763203153
# HELP consul_raft_rpc_appendEntries_processLogs_sum Telegraf collected metric
# TYPE consul_raft_rpc_appendEntries_processLogs_sum untyped
consul_raft_rpc_appendEntries_processLogs_sum{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0.062344
# HELP consul_raft_rpc_appendEntries_processLogs_upper Telegraf collected metric
# TYPE consul_raft_rpc_appendEntries_processLogs_upper untyped
consul_raft_rpc_appendEntries_processLogs_upper{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0.013267
# HELP consul_raft_rpc_appendEntries_stddev Telegraf collected metric
# TYPE consul_raft_rpc_appendEntries_stddev untyped
consul_raft_rpc_appendEntries_stddev{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0.9393835241674746
# HELP consul_raft_rpc_appendEntries_storeLogs_90_percentile Telegraf collected metric
# TYPE consul_raft_rpc_appendEntries_storeLogs_90_percentile untyped
consul_raft_rpc_appendEntries_storeLogs_90_percentile{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 11.927401
# HELP consul_raft_rpc_appendEntries_storeLogs_count Telegraf collected metric
# TYPE consul_raft_rpc_appendEntries_storeLogs_count untyped
consul_raft_rpc_appendEntries_storeLogs_count{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 6
# HELP consul_raft_rpc_appendEntries_storeLogs_lower Telegraf collected metric
# TYPE consul_raft_rpc_appendEntries_storeLogs_lower untyped
consul_raft_rpc_appendEntries_storeLogs_lower{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 4.976229
# HELP consul_raft_rpc_appendEntries_storeLogs_mean Telegraf collected metric
# TYPE consul_raft_rpc_appendEntries_storeLogs_mean untyped
consul_raft_rpc_appendEntries_storeLogs_mean{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 7.5146991666666665
# HELP consul_raft_rpc_appendEntries_storeLogs_stddev Telegraf collected metric
# TYPE consul_raft_rpc_appendEntries_storeLogs_stddev untyped
consul_raft_rpc_appendEntries_storeLogs_stddev{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 3.083714780644357
# HELP consul_raft_rpc_appendEntries_storeLogs_sum Telegraf collected metric
# TYPE consul_raft_rpc_appendEntries_storeLogs_sum untyped
consul_raft_rpc_appendEntries_storeLogs_sum{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 45.088195000000006
# HELP consul_raft_rpc_appendEntries_storeLogs_upper Telegraf collected metric
# TYPE consul_raft_rpc_appendEntries_storeLogs_upper untyped
consul_raft_rpc_appendEntries_storeLogs_upper{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 11.927401
# HELP consul_raft_rpc_appendEntries_sum Telegraf collected metric
# TYPE consul_raft_rpc_appendEntries_sum untyped
consul_raft_rpc_appendEntries_sum{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 49.8574
# HELP consul_raft_rpc_appendEntries_upper Telegraf collected metric
# TYPE consul_raft_rpc_appendEntries_upper untyped
consul_raft_rpc_appendEntries_upper{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 11.976887
# HELP consul_raft_rpc_processHeartbeat_90_percentile Telegraf collected metric
# TYPE consul_raft_rpc_processHeartbeat_90_percentile untyped
consul_raft_rpc_processHeartbeat_90_percentile{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0.056779
# HELP consul_raft_rpc_processHeartbeat_count Telegraf collected metric
# TYPE consul_raft_rpc_processHeartbeat_count untyped
consul_raft_rpc_processHeartbeat_count{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 41
# HELP consul_raft_rpc_processHeartbeat_lower Telegraf collected metric
# TYPE consul_raft_rpc_processHeartbeat_lower untyped
consul_raft_rpc_processHeartbeat_lower{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0.022834
# HELP consul_raft_rpc_processHeartbeat_mean Telegraf collected metric
# TYPE consul_raft_rpc_processHeartbeat_mean untyped
consul_raft_rpc_processHeartbeat_mean{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0.04517173170731707
# HELP consul_raft_rpc_processHeartbeat_stddev Telegraf collected metric
# TYPE consul_raft_rpc_processHeartbeat_stddev untyped
consul_raft_rpc_processHeartbeat_stddev{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0.012594230049422443
# HELP consul_raft_rpc_processHeartbeat_sum Telegraf collected metric
# TYPE consul_raft_rpc_processHeartbeat_sum untyped
consul_raft_rpc_processHeartbeat_sum{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 1.8520410000000003
# HELP consul_raft_rpc_processHeartbeat_upper Telegraf collected metric
# TYPE consul_raft_rpc_processHeartbeat_upper untyped
consul_raft_rpc_processHeartbeat_upper{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0.094538
# HELP consul_rpc_query Telegraf collected metric
# TYPE consul_rpc_query counter
consul_rpc_query{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="counter",project_name="consul-service",role="Server"} 4
# HELP consul_rpc_request Telegraf collected metric
# TYPE consul_rpc_request counter
consul_rpc_request{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="counter",project_name="consul-service",role="Server"} 25
# HELP consul_runtime_alloc_bytes Telegraf collected metric
# TYPE consul_runtime_alloc_bytes gauge
consul_runtime_alloc_bytes{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="gauge",project_name="consul-service",role="Server"} 9.891424e+06
# HELP consul_runtime_free_count Telegraf collected metric
# TYPE consul_runtime_free_count gauge
consul_runtime_free_count{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="gauge",project_name="consul-service",role="Server"} 1.40498608e+08
# HELP consul_runtime_gc_pause_ns_90_percentile Telegraf collected metric
# TYPE consul_runtime_gc_pause_ns_90_percentile untyped
consul_runtime_gc_pause_ns_90_percentile{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 62098
# HELP consul_runtime_gc_pause_ns_count Telegraf collected metric
# TYPE consul_runtime_gc_pause_ns_count untyped
consul_runtime_gc_pause_ns_count{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 2
# HELP consul_runtime_gc_pause_ns_lower Telegraf collected metric
# TYPE consul_runtime_gc_pause_ns_lower untyped
consul_runtime_gc_pause_ns_lower{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 37366
# HELP consul_runtime_gc_pause_ns_mean Telegraf collected metric
# TYPE consul_runtime_gc_pause_ns_mean untyped
consul_runtime_gc_pause_ns_mean{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 49732
# HELP consul_runtime_gc_pause_ns_stddev Telegraf collected metric
# TYPE consul_runtime_gc_pause_ns_stddev untyped
consul_runtime_gc_pause_ns_stddev{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 12366
# HELP consul_runtime_gc_pause_ns_sum Telegraf collected metric
# TYPE consul_runtime_gc_pause_ns_sum untyped
consul_runtime_gc_pause_ns_sum{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 99464
# HELP consul_runtime_gc_pause_ns_upper Telegraf collected metric
# TYPE consul_runtime_gc_pause_ns_upper untyped
consul_runtime_gc_pause_ns_upper{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 62098
# HELP consul_runtime_heap_objects Telegraf collected metric
# TYPE consul_runtime_heap_objects gauge
consul_runtime_heap_objects{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="gauge",project_name="consul-service",role="Server"} 47037
# HELP consul_runtime_malloc_count Telegraf collected metric
# TYPE consul_runtime_malloc_count gauge
consul_runtime_malloc_count{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="gauge",project_name="consul-service",role="Server"} 1.40545648e+08
# HELP consul_runtime_num_goroutines Telegraf collected metric
# TYPE consul_runtime_num_goroutines gauge
consul_runtime_num_goroutines{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="gauge",project_name="consul-service",role="Server"} 96
# HELP consul_runtime_sys_bytes Telegraf collected metric
# TYPE consul_runtime_sys_bytes gauge
consul_runtime_sys_bytes{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="gauge",project_name="consul-service",role="Server"} 7.505536e+07
# HELP consul_runtime_total_gc_pause_ns Telegraf collected metric
# TYPE consul_runtime_total_gc_pause_ns gauge
consul_runtime_total_gc_pause_ns{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="gauge",project_name="consul-service",role="Server"} 6.95344832e+08
# HELP consul_runtime_total_gc_runs Telegraf collected metric
# TYPE consul_runtime_total_gc_runs gauge
consul_runtime_total_gc_runs{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="gauge",project_name="consul-service",role="Server"} 7724
# HELP consul_serf_coordinate_adjustment_ms_90_percentile Telegraf collected metric
# TYPE consul_serf_coordinate_adjustment_ms_90_percentile untyped
consul_serf_coordinate_adjustment_ms_90_percentile{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0.086559
# HELP consul_serf_coordinate_adjustment_ms_count Telegraf collected metric
# TYPE consul_serf_coordinate_adjustment_ms_count untyped
consul_serf_coordinate_adjustment_ms_count{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 35
# HELP consul_serf_coordinate_adjustment_ms_lower Telegraf collected metric
# TYPE consul_serf_coordinate_adjustment_ms_lower untyped
consul_serf_coordinate_adjustment_ms_lower{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0.020771
# HELP consul_serf_coordinate_adjustment_ms_mean Telegraf collected metric
# TYPE consul_serf_coordinate_adjustment_ms_mean untyped
consul_serf_coordinate_adjustment_ms_mean{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0.049225114285714286
# HELP consul_serf_coordinate_adjustment_ms_stddev Telegraf collected metric
# TYPE consul_serf_coordinate_adjustment_ms_stddev untyped
consul_serf_coordinate_adjustment_ms_stddev{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0.028809318708433247
# HELP consul_serf_coordinate_adjustment_ms_sum Telegraf collected metric
# TYPE consul_serf_coordinate_adjustment_ms_sum untyped
consul_serf_coordinate_adjustment_ms_sum{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 1.7228790000000003
# HELP consul_serf_coordinate_adjustment_ms_upper Telegraf collected metric
# TYPE consul_serf_coordinate_adjustment_ms_upper untyped
consul_serf_coordinate_adjustment_ms_upper{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0.139714
# HELP consul_serf_queue_Event_90_percentile Telegraf collected metric
# TYPE consul_serf_queue_Event_90_percentile untyped
consul_serf_queue_Event_90_percentile{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0
# HELP consul_serf_queue_Event_count Telegraf collected metric
# TYPE consul_serf_queue_Event_count untyped
consul_serf_queue_Event_count{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 2
# HELP consul_serf_queue_Event_lower Telegraf collected metric
# TYPE consul_serf_queue_Event_lower untyped
consul_serf_queue_Event_lower{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0
# HELP consul_serf_queue_Event_mean Telegraf collected metric
# TYPE consul_serf_queue_Event_mean untyped
consul_serf_queue_Event_mean{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0
# HELP consul_serf_queue_Event_stddev Telegraf collected metric
# TYPE consul_serf_queue_Event_stddev untyped
consul_serf_queue_Event_stddev{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0
# HELP consul_serf_queue_Event_sum Telegraf collected metric
# TYPE consul_serf_queue_Event_sum untyped
consul_serf_queue_Event_sum{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0
# HELP consul_serf_queue_Event_upper Telegraf collected metric
# TYPE consul_serf_queue_Event_upper untyped
consul_serf_queue_Event_upper{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0
# HELP consul_serf_queue_Intent_90_percentile Telegraf collected metric
# TYPE consul_serf_queue_Intent_90_percentile untyped
consul_serf_queue_Intent_90_percentile{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0
# HELP consul_serf_queue_Intent_count Telegraf collected metric
# TYPE consul_serf_queue_Intent_count untyped
consul_serf_queue_Intent_count{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 2
# HELP consul_serf_queue_Intent_lower Telegraf collected metric
# TYPE consul_serf_queue_Intent_lower untyped
consul_serf_queue_Intent_lower{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0
# HELP consul_serf_queue_Intent_mean Telegraf collected metric
# TYPE consul_serf_queue_Intent_mean untyped
consul_serf_queue_Intent_mean{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0
# HELP consul_serf_queue_Intent_stddev Telegraf collected metric
# TYPE consul_serf_queue_Intent_stddev untyped
consul_serf_queue_Intent_stddev{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0
# HELP consul_serf_queue_Intent_sum Telegraf collected metric
# TYPE consul_serf_queue_Intent_sum untyped
consul_serf_queue_Intent_sum{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0
# HELP consul_serf_queue_Intent_upper Telegraf collected metric
# TYPE consul_serf_queue_Intent_upper untyped
consul_serf_queue_Intent_upper{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0
# HELP consul_serf_queue_Query_90_percentile Telegraf collected metric
# TYPE consul_serf_queue_Query_90_percentile untyped
consul_serf_queue_Query_90_percentile{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0
# HELP consul_serf_queue_Query_count Telegraf collected metric
# TYPE consul_serf_queue_Query_count untyped
consul_serf_queue_Query_count{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 2
# HELP consul_serf_queue_Query_lower Telegraf collected metric
# TYPE consul_serf_queue_Query_lower untyped
consul_serf_queue_Query_lower{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0
# HELP consul_serf_queue_Query_mean Telegraf collected metric
# TYPE consul_serf_queue_Query_mean untyped
consul_serf_queue_Query_mean{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0
# HELP consul_serf_queue_Query_stddev Telegraf collected metric
# TYPE consul_serf_queue_Query_stddev untyped
consul_serf_queue_Query_stddev{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0
# HELP consul_serf_queue_Query_sum Telegraf collected metric
# TYPE consul_serf_queue_Query_sum untyped
consul_serf_queue_Query_sum{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0
# HELP consul_serf_queue_Query_upper Telegraf collected metric
# TYPE consul_serf_queue_Query_upper untyped
consul_serf_queue_Query_upper{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="timing",project_name="consul-service",role="Server"} 0
# HELP consul_session_ttl_active Telegraf collected metric
# TYPE consul_session_ttl_active gauge
consul_session_ttl_active{datacenter="dc1",fullname="consul",host="consul-server-0",metric_type="gauge",project_name="consul-service",role="Server"} 0
# HELP cpu_usage_guest Telegraf collected metric
# TYPE cpu_usage_guest gauge
cpu_usage_guest{cpu="cpu-total",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
cpu_usage_guest{cpu="cpu0",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
cpu_usage_guest{cpu="cpu1",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
cpu_usage_guest{cpu="cpu2",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
cpu_usage_guest{cpu="cpu3",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
cpu_usage_guest{cpu="cpu4",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
cpu_usage_guest{cpu="cpu5",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
cpu_usage_guest{cpu="cpu6",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
cpu_usage_guest{cpu="cpu7",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
# HELP cpu_usage_guest_nice Telegraf collected metric
# TYPE cpu_usage_guest_nice gauge
cpu_usage_guest_nice{cpu="cpu-total",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
cpu_usage_guest_nice{cpu="cpu0",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
cpu_usage_guest_nice{cpu="cpu1",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
cpu_usage_guest_nice{cpu="cpu2",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
cpu_usage_guest_nice{cpu="cpu3",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
cpu_usage_guest_nice{cpu="cpu4",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
cpu_usage_guest_nice{cpu="cpu5",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
cpu_usage_guest_nice{cpu="cpu6",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
cpu_usage_guest_nice{cpu="cpu7",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
# HELP cpu_usage_idle Telegraf collected metric
# TYPE cpu_usage_idle gauge
cpu_usage_idle{cpu="cpu-total",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 85.83460949443992
cpu_usage_idle{cpu="cpu0",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 85.7878475794261
cpu_usage_idle{cpu="cpu1",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 84.98131158703607
cpu_usage_idle{cpu="cpu2",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 83.92250169888257
cpu_usage_idle{cpu="cpu3",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 85.44712682747318
cpu_usage_idle{cpu="cpu4",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 86.51494565234987
cpu_usage_idle{cpu="cpu5",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 86.46718803130665
cpu_usage_idle{cpu="cpu6",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 86.70068027235182
cpu_usage_idle{cpu="cpu7",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 86.79565512592507
# HELP cpu_usage_iowait Telegraf collected metric
# TYPE cpu_usage_iowait gauge
cpu_usage_iowait{cpu="cpu-total",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0.212693551129777
cpu_usage_iowait{cpu="cpu0",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0.03432887057452017
cpu_usage_iowait{cpu="cpu1",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0.06795786612463765
cpu_usage_iowait{cpu="cpu2",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0.7817811012859472
cpu_usage_iowait{cpu="cpu3",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0.1020061203635544
cpu_usage_iowait{cpu="cpu4",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0.20380434781695334
cpu_usage_iowait{cpu="cpu5",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0.10200612037592427
cpu_usage_iowait{cpu="cpu6",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0.13605442175973417
cpu_usage_iowait{cpu="cpu7",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0.33944331296222424
# HELP cpu_usage_irq Telegraf collected metric
# TYPE cpu_usage_irq gauge
cpu_usage_irq{cpu="cpu-total",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
cpu_usage_irq{cpu="cpu0",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
cpu_usage_irq{cpu="cpu1",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
cpu_usage_irq{cpu="cpu2",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
cpu_usage_irq{cpu="cpu3",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
cpu_usage_irq{cpu="cpu4",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
cpu_usage_irq{cpu="cpu5",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
cpu_usage_irq{cpu="cpu6",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
cpu_usage_irq{cpu="cpu7",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
# HELP cpu_usage_nice Telegraf collected metric
# TYPE cpu_usage_nice gauge
cpu_usage_nice{cpu="cpu-total",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0.008507742045183344
cpu_usage_nice{cpu="cpu0",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
cpu_usage_nice{cpu="cpu1",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
cpu_usage_nice{cpu="cpu2",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0.03399048266468831
cpu_usage_nice{cpu="cpu3",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0.03400204012247333
cpu_usage_nice{cpu="cpu4",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
cpu_usage_nice{cpu="cpu5",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
cpu_usage_nice{cpu="cpu6",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
cpu_usage_nice{cpu="cpu7",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
# HELP cpu_usage_softirq Telegraf collected metric
# TYPE cpu_usage_softirq gauge
cpu_usage_softirq{cpu="cpu-total",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0.21694742215323917
cpu_usage_softirq{cpu="cpu0",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0.9612083762364301
cpu_usage_softirq{cpu="cpu1",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0.2038735983692774
cpu_usage_softirq{cpu="cpu2",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0.13596193065856002
cpu_usage_softirq{cpu="cpu3",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0.06800408024494665
cpu_usage_softirq{cpu="cpu4",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0.10190217391233831
cpu_usage_softirq{cpu="cpu5",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0.06800408024494665
cpu_usage_softirq{cpu="cpu6",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0.10204081632676104
cpu_usage_softirq{cpu="cpu7",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0.10183299389005653
# HELP cpu_usage_steal Telegraf collected metric
# TYPE cpu_usage_steal gauge
cpu_usage_steal{cpu="cpu-total",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0.05955419431608996
cpu_usage_steal{cpu="cpu0",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0.10298661173995204
cpu_usage_steal{cpu="cpu1",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0.033978933061546236
cpu_usage_steal{cpu="cpu2",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0.033990482664640005
cpu_usage_steal{cpu="cpu3",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0.06800408024494665
cpu_usage_steal{cpu="cpu4",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0.06793478260822554
cpu_usage_steal{cpu="cpu5",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0.06800408024494665
cpu_usage_steal{cpu="cpu6",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0.03401360544225368
cpu_usage_steal{cpu="cpu7",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0.033944331296685504
# HELP cpu_usage_system Telegraf collected metric
# TYPE cpu_usage_system gauge
cpu_usage_system{cpu="cpu-total",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 4.33044070102405
cpu_usage_system{cpu="cpu0",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 3.810504634383689
cpu_usage_system{cpu="cpu1",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 4.553177030239469
cpu_usage_system{cpu="cpu2",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 5.064581917053774
cpu_usage_system{cpu="cpu3",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 4.726283577026112
cpu_usage_system{cpu="cpu4",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 4.415760869548562
cpu_usage_system{cpu="cpu5",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 4.148248894949477
cpu_usage_system{cpu="cpu6",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 3.945578231316895
cpu_usage_system{cpu="cpu7",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 3.937542430430955
# HELP cpu_usage_user Telegraf collected metric
# TYPE cpu_usage_user gauge
cpu_usage_user{cpu="cpu-total",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 9.337246894577403
cpu_usage_user{cpu="cpu0",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 9.303123927256063
cpu_usage_user{cpu="cpu1",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 10.159700985379919
cpu_usage_user{cpu="cpu2",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 10.027192386117491
cpu_usage_user{cpu="cpu3",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 9.554573274514738
cpu_usage_user{cpu="cpu4",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 8.695652173852869
cpu_usage_user{cpu="cpu5",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 9.14654879296156
cpu_usage_user{cpu="cpu6",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 9.08163265313355
cpu_usage_user{cpu="cpu7",datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 8.791581805837687
# HELP diskio_io_time Telegraf collected metric
# TYPE diskio_io_time counter
diskio_io_time{datacenter="dc1",fullname="consul",host="consul-server-0",name="vda",project_name="consul-service",role="Server"} 7.1901406e+07
diskio_io_time{datacenter="dc1",fullname="consul",host="consul-server-0",name="vda1",project_name="consul-service",role="Server"} 5.4476926e+07
diskio_io_time{datacenter="dc1",fullname="consul",host="consul-server-0",name="vdb",project_name="consul-service",role="Server"} 143
diskio_io_time{datacenter="dc1",fullname="consul",host="consul-server-0",name="vdc",project_name="consul-service",role="Server"} 4.1604763e+07
# HELP diskio_iops_in_progress Telegraf collected metric
# TYPE diskio_iops_in_progress counter
diskio_iops_in_progress{datacenter="dc1",fullname="consul",host="consul-server-0",name="vda",project_name="consul-service",role="Server"} 0
diskio_iops_in_progress{datacenter="dc1",fullname="consul",host="consul-server-0",name="vda1",project_name="consul-service",role="Server"} 0
diskio_iops_in_progress{datacenter="dc1",fullname="consul",host="consul-server-0",name="vdb",project_name="consul-service",role="Server"} 0
diskio_iops_in_progress{datacenter="dc1",fullname="consul",host="consul-server-0",name="vdc",project_name="consul-service",role="Server"} 0
# HELP diskio_read_bytes Telegraf collected metric
# TYPE diskio_read_bytes counter
diskio_read_bytes{datacenter="dc1",fullname="consul",host="consul-server-0",name="vda",project_name="consul-service",role="Server"} 4.610237327872e+12
diskio_read_bytes{datacenter="dc1",fullname="consul",host="consul-server-0",name="vda1",project_name="consul-service",role="Server"} 4.610236246528e+12
diskio_read_bytes{datacenter="dc1",fullname="consul",host="consul-server-0",name="vdb",project_name="consul-service",role="Server"} 2.12992e+06
diskio_read_bytes{datacenter="dc1",fullname="consul",host="consul-server-0",name="vdc",project_name="consul-service",role="Server"} 4.31769924608e+11
# HELP diskio_read_time Telegraf collected metric
# TYPE diskio_read_time counter
diskio_read_time{datacenter="dc1",fullname="consul",host="consul-server-0",name="vda",project_name="consul-service",role="Server"} 1.758329184e+09
diskio_read_time{datacenter="dc1",fullname="consul",host="consul-server-0",name="vda1",project_name="consul-service",role="Server"} 1.758329178e+09
diskio_read_time{datacenter="dc1",fullname="consul",host="consul-server-0",name="vdb",project_name="consul-service",role="Server"} 330
diskio_read_time{datacenter="dc1",fullname="consul",host="consul-server-0",name="vdc",project_name="consul-service",role="Server"} 5.24353375e+08
# HELP diskio_reads Telegraf collected metric
# TYPE diskio_reads counter
diskio_reads{datacenter="dc1",fullname="consul",host="consul-server-0",name="vda",project_name="consul-service",role="Server"} 6.0892955e+07
diskio_reads{datacenter="dc1",fullname="consul",host="consul-server-0",name="vda1",project_name="consul-service",role="Server"} 6.0892923e+07
diskio_reads{datacenter="dc1",fullname="consul",host="consul-server-0",name="vdb",project_name="consul-service",role="Server"} 88
diskio_reads{datacenter="dc1",fullname="consul",host="consul-server-0",name="vdc",project_name="consul-service",role="Server"} 1.153379e+06
# HELP diskio_weighted_io_time Telegraf collected metric
# TYPE diskio_weighted_io_time counter
diskio_weighted_io_time{datacenter="dc1",fullname="consul",host="consul-server-0",name="vda",project_name="consul-service",role="Server"} 1.261213773e+09
diskio_weighted_io_time{datacenter="dc1",fullname="consul",host="consul-server-0",name="vda1",project_name="consul-service",role="Server"} 1.225543593e+09
diskio_weighted_io_time{datacenter="dc1",fullname="consul",host="consul-server-0",name="vdb",project_name="consul-service",role="Server"} 330
diskio_weighted_io_time{datacenter="dc1",fullname="consul",host="consul-server-0",name="vdc",project_name="consul-service",role="Server"} 3.101505045e+09
# HELP diskio_write_bytes Telegraf collected metric
# TYPE diskio_write_bytes counter
diskio_write_bytes{datacenter="dc1",fullname="consul",host="consul-server-0",name="vda",project_name="consul-service",role="Server"} 1.0426337536e+12
diskio_write_bytes{datacenter="dc1",fullname="consul",host="consul-server-0",name="vda1",project_name="consul-service",role="Server"} 1.0426337536e+12
diskio_write_bytes{datacenter="dc1",fullname="consul",host="consul-server-0",name="vdb",project_name="consul-service",role="Server"} 0
diskio_write_bytes{datacenter="dc1",fullname="consul",host="consul-server-0",name="vdc",project_name="consul-service",role="Server"} 4.51002793984e+11
# HELP diskio_write_time Telegraf collected metric
# TYPE diskio_write_time counter
diskio_write_time{datacenter="dc1",fullname="consul",host="consul-server-0",name="vda",project_name="consul-service",role="Server"} 3.66855249e+08
diskio_write_time{datacenter="dc1",fullname="consul",host="consul-server-0",name="vda1",project_name="consul-service",role="Server"} 3.28260783e+08
diskio_write_time{datacenter="dc1",fullname="consul",host="consul-server-0",name="vdb",project_name="consul-service",role="Server"} 0
diskio_write_time{datacenter="dc1",fullname="consul",host="consul-server-0",name="vdc",project_name="consul-service",role="Server"} 3.668256164e+09
# HELP diskio_writes Telegraf collected metric
# TYPE diskio_writes counter
diskio_writes{datacenter="dc1",fullname="consul",host="consul-server-0",name="vda",project_name="consul-service",role="Server"} 7.72827e+07
diskio_writes{datacenter="dc1",fullname="consul",host="consul-server-0",name="vda1",project_name="consul-service",role="Server"} 7.2802132e+07
diskio_writes{datacenter="dc1",fullname="consul",host="consul-server-0",name="vdb",project_name="consul-service",role="Server"} 0
diskio_writes{datacenter="dc1",fullname="consul",host="consul-server-0",name="vdc",project_name="consul-service",role="Server"} 971160
# HELP go_gc_duration_seconds A summary of the GC invocation durations.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 9.3165e-05
go_gc_duration_seconds{quantile="0.25"} 0.000166623
go_gc_duration_seconds{quantile="0.5"} 0.000221478
go_gc_duration_seconds{quantile="0.75"} 0.000572069
go_gc_duration_seconds{quantile="1"} 0.100084546
go_gc_duration_seconds_sum 23.231454054
go_gc_duration_seconds_count 2242
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 24
# HELP go_memstats_alloc_bytes Number of bytes allocated and still in use.
# TYPE go_memstats_alloc_bytes gauge
go_memstats_alloc_bytes 6.074456e+06
# HELP go_memstats_alloc_bytes_total Total number of bytes allocated, even if freed.
# TYPE go_memstats_alloc_bytes_total counter
go_memstats_alloc_bytes_total 8.058774896e+09
# HELP go_memstats_buck_hash_sys_bytes Number of bytes used by the profiling bucket hash table.
# TYPE go_memstats_buck_hash_sys_bytes gauge
go_memstats_buck_hash_sys_bytes 1.551764e+06
# HELP go_memstats_frees_total Total number of frees.
# TYPE go_memstats_frees_total counter
go_memstats_frees_total 1.14698487e+08
# HELP go_memstats_gc_sys_bytes Number of bytes used for garbage collection system metadata.
# TYPE go_memstats_gc_sys_bytes gauge
go_memstats_gc_sys_bytes 2.383872e+06
# HELP go_memstats_heap_alloc_bytes Number of heap bytes allocated and still in use.
# TYPE go_memstats_heap_alloc_bytes gauge
go_memstats_heap_alloc_bytes 6.074456e+06
# HELP go_memstats_heap_idle_bytes Number of heap bytes waiting to be used.
# TYPE go_memstats_heap_idle_bytes gauge
go_memstats_heap_idle_bytes 5.6778752e+07
# HELP go_memstats_heap_inuse_bytes Number of heap bytes that are in use.
# TYPE go_memstats_heap_inuse_bytes gauge
go_memstats_heap_inuse_bytes 8.691712e+06
# HELP go_memstats_heap_objects Number of allocated objects.
# TYPE go_memstats_heap_objects gauge
go_memstats_heap_objects 28674
# HELP go_memstats_heap_released_bytes_total Total number of heap bytes released to OS.
# TYPE go_memstats_heap_released_bytes_total counter
go_memstats_heap_released_bytes_total 5.4042624e+07
# HELP go_memstats_heap_sys_bytes Number of heap bytes obtained from system.
# TYPE go_memstats_heap_sys_bytes gauge
go_memstats_heap_sys_bytes 6.5470464e+07
# HELP go_memstats_last_gc_time_seconds Number of seconds since 1970 of last garbage collection.
# TYPE go_memstats_last_gc_time_seconds gauge
go_memstats_last_gc_time_seconds 1.613405778253354e+09
# HELP go_memstats_lookups_total Total number of pointer lookups.
# TYPE go_memstats_lookups_total counter
go_memstats_lookups_total 0
# HELP go_memstats_mallocs_total Total number of mallocs.
# TYPE go_memstats_mallocs_total counter
go_memstats_mallocs_total 1.14727161e+08
# HELP go_memstats_mcache_inuse_bytes Number of bytes in use by mcache structures.
# TYPE go_memstats_mcache_inuse_bytes gauge
go_memstats_mcache_inuse_bytes 13824
# HELP go_memstats_mcache_sys_bytes Number of bytes used for mcache structures obtained from system.
# TYPE go_memstats_mcache_sys_bytes gauge
go_memstats_mcache_sys_bytes 16384
# HELP go_memstats_mspan_inuse_bytes Number of bytes in use by mspan structures.
# TYPE go_memstats_mspan_inuse_bytes gauge
go_memstats_mspan_inuse_bytes 122968
# HELP go_memstats_mspan_sys_bytes Number of bytes used for mspan structures obtained from system.
# TYPE go_memstats_mspan_sys_bytes gauge
go_memstats_mspan_sys_bytes 163840
# HELP go_memstats_next_gc_bytes Number of heap bytes when next garbage collection will take place.
# TYPE go_memstats_next_gc_bytes gauge
go_memstats_next_gc_bytes 1.0578416e+07
# HELP go_memstats_other_sys_bytes Number of bytes used for other system allocations.
# TYPE go_memstats_other_sys_bytes gauge
go_memstats_other_sys_bytes 2.10826e+06
# HELP go_memstats_stack_inuse_bytes Number of bytes in use by the stack allocator.
# TYPE go_memstats_stack_inuse_bytes gauge
go_memstats_stack_inuse_bytes 1.6384e+06
# HELP go_memstats_stack_sys_bytes Number of bytes obtained from system for stack allocator.
# TYPE go_memstats_stack_sys_bytes gauge
go_memstats_stack_sys_bytes 1.6384e+06
# HELP go_memstats_sys_bytes Number of bytes obtained by system. Sum of all system allocations.
# TYPE go_memstats_sys_bytes gauge
go_memstats_sys_bytes 7.3332984e+07
# HELP linux_sysctl_fs_aio_max_nr Telegraf collected metric
# TYPE linux_sysctl_fs_aio_max_nr untyped
linux_sysctl_fs_aio_max_nr{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 65536
# HELP linux_sysctl_fs_aio_nr Telegraf collected metric
# TYPE linux_sysctl_fs_aio_nr untyped
linux_sysctl_fs_aio_nr{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 320
# HELP linux_sysctl_fs_dentry_age_limit Telegraf collected metric
# TYPE linux_sysctl_fs_dentry_age_limit untyped
linux_sysctl_fs_dentry_age_limit{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 45
# HELP linux_sysctl_fs_dentry_nr Telegraf collected metric
# TYPE linux_sysctl_fs_dentry_nr untyped
linux_sysctl_fs_dentry_nr{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 145533
# HELP linux_sysctl_fs_dentry_unused_nr Telegraf collected metric
# TYPE linux_sysctl_fs_dentry_unused_nr untyped
linux_sysctl_fs_dentry_unused_nr{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 76122
# HELP linux_sysctl_fs_dentry_want_pages Telegraf collected metric
# TYPE linux_sysctl_fs_dentry_want_pages untyped
linux_sysctl_fs_dentry_want_pages{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
# HELP linux_sysctl_fs_file_max Telegraf collected metric
# TYPE linux_sysctl_fs_file_max untyped
linux_sysctl_fs_file_max{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 1.19709e+06
# HELP linux_sysctl_fs_file_nr Telegraf collected metric
# TYPE linux_sysctl_fs_file_nr untyped
linux_sysctl_fs_file_nr{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 8320
# HELP linux_sysctl_fs_inode_free_nr Telegraf collected metric
# TYPE linux_sysctl_fs_inode_free_nr untyped
linux_sysctl_fs_inode_free_nr{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 24550
# HELP linux_sysctl_fs_inode_nr Telegraf collected metric
# TYPE linux_sysctl_fs_inode_nr untyped
linux_sysctl_fs_inode_nr{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 125647
# HELP linux_sysctl_fs_inode_preshrink_nr Telegraf collected metric
# TYPE linux_sysctl_fs_inode_preshrink_nr untyped
linux_sysctl_fs_inode_preshrink_nr{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
# HELP mem_active Telegraf collected metric
# TYPE mem_active gauge
mem_active{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 9.731948544e+09
# HELP mem_available Telegraf collected metric
# TYPE mem_available gauge
mem_available{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 2.810503168e+09
# HELP mem_available_percent Telegraf collected metric
# TYPE mem_available_percent gauge
mem_available_percent{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 22.613826653848523
# HELP mem_buffered Telegraf collected metric
# TYPE mem_buffered gauge
mem_buffered{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 180224
# HELP mem_cached Telegraf collected metric
# TYPE mem_cached gauge
mem_cached{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 3.164741632e+09
# HELP mem_commit_limit Telegraf collected metric
# TYPE mem_commit_limit gauge
mem_commit_limit{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 6.21412352e+09
# HELP mem_committed_as Telegraf collected metric
# TYPE mem_committed_as gauge
mem_committed_as{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 2.2322089984e+10
# HELP mem_dirty Telegraf collected metric
# TYPE mem_dirty gauge
mem_dirty{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 581632
# HELP mem_free Telegraf collected metric
# TYPE mem_free gauge
mem_free{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 4.57601024e+08
# HELP mem_high_free Telegraf collected metric
# TYPE mem_high_free gauge
mem_high_free{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
# HELP mem_high_total Telegraf collected metric
# TYPE mem_high_total gauge
mem_high_total{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
# HELP mem_huge_page_size Telegraf collected metric
# TYPE mem_huge_page_size gauge
mem_huge_page_size{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 2.097152e+06
# HELP mem_huge_pages_free Telegraf collected metric
# TYPE mem_huge_pages_free gauge
mem_huge_pages_free{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
# HELP mem_huge_pages_total Telegraf collected metric
# TYPE mem_huge_pages_total gauge
mem_huge_pages_total{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
# HELP mem_inactive Telegraf collected metric
# TYPE mem_inactive gauge
mem_inactive{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 1.488584704e+09
# HELP mem_low_free Telegraf collected metric
# TYPE mem_low_free gauge
mem_low_free{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
# HELP mem_low_total Telegraf collected metric
# TYPE mem_low_total gauge
mem_low_total{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
# HELP mem_mapped Telegraf collected metric
# TYPE mem_mapped gauge
mem_mapped{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 4.0486912e+08
# HELP mem_page_tables Telegraf collected metric
# TYPE mem_page_tables gauge
mem_page_tables{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 4.2258432e+07
# HELP mem_shared Telegraf collected metric
# TYPE mem_shared gauge
mem_shared{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 7.29657344e+08
# HELP mem_slab Telegraf collected metric
# TYPE mem_slab gauge
mem_slab{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 4.46738432e+08
# HELP mem_swap_cached Telegraf collected metric
# TYPE mem_swap_cached gauge
mem_swap_cached{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
# HELP mem_swap_free Telegraf collected metric
# TYPE mem_swap_free gauge
mem_swap_free{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
# HELP mem_swap_total Telegraf collected metric
# TYPE mem_swap_total gauge
mem_swap_total{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
# HELP mem_total Telegraf collected metric
# TYPE mem_total gauge
mem_total{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 1.2428251136e+10
# HELP mem_used Telegraf collected metric
# TYPE mem_used gauge
mem_used{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 8.805728256e+09
# HELP mem_used_percent Telegraf collected metric
# TYPE mem_used_percent gauge
mem_used_percent{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 70.85251303373727
# HELP mem_vmalloc_chunk Telegraf collected metric
# TYPE mem_vmalloc_chunk gauge
mem_vmalloc_chunk{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 3.5184189190144e+13
# HELP mem_vmalloc_total Telegraf collected metric
# TYPE mem_vmalloc_total gauge
mem_vmalloc_total{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 3.5184372087808e+13
# HELP mem_vmalloc_used Telegraf collected metric
# TYPE mem_vmalloc_used gauge
mem_vmalloc_used{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 5.4390784e+07
# HELP mem_wired Telegraf collected metric
# TYPE mem_wired gauge
mem_wired{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
# HELP mem_write_back Telegraf collected metric
# TYPE mem_write_back gauge
mem_write_back{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
# HELP mem_write_back_tmp Telegraf collected metric
# TYPE mem_write_back_tmp gauge
mem_write_back_tmp{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
# HELP process_cpu_seconds_total Total user and system CPU time spent in seconds.
# TYPE process_cpu_seconds_total counter
process_cpu_seconds_total 258.21
# HELP process_max_fds Maximum number of open file descriptors.
# TYPE process_max_fds gauge
process_max_fds 1.048576e+06
# HELP process_open_fds Number of open file descriptors.
# TYPE process_open_fds gauge
process_open_fds 8
# HELP process_resident_memory_bytes Resident memory size in bytes.
# TYPE process_resident_memory_bytes gauge
process_resident_memory_bytes 1.5876096e+07
# HELP process_start_time_seconds Start time of the process since unix epoch in seconds.
# TYPE process_start_time_seconds gauge
process_start_time_seconds 1.61331344633e+09
# HELP process_virtual_memory_bytes Virtual memory size in bytes.
# TYPE process_virtual_memory_bytes gauge
process_virtual_memory_bytes 1.66879232e+08
# HELP swap_free Telegraf collected metric
# TYPE swap_free gauge
swap_free{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
# HELP swap_in Telegraf collected metric
# TYPE swap_in counter
swap_in{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
# HELP swap_out Telegraf collected metric
# TYPE swap_out counter
swap_out{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
# HELP swap_total Telegraf collected metric
# TYPE swap_total gauge
swap_total{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
# HELP swap_used Telegraf collected metric
# TYPE swap_used gauge
swap_used{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
# HELP swap_used_percent Telegraf collected metric
# TYPE swap_used_percent gauge
swap_used_percent{datacenter="dc1",fullname="consul",host="consul-server-0",project_name="consul-service",role="Server"} 0
```

### Consul Client

```prometheus
# HELP consul_acl_ResolveToken This measures the time it takes to resolve an ACL token.
# TYPE consul_acl_ResolveToken summary
consul_acl_ResolveToken{quantile="0.5"} 0.06769300252199173
consul_acl_ResolveToken{quantile="0.9"} 0.06769300252199173
consul_acl_ResolveToken{quantile="0.99"} 0.06769300252199173
consul_acl_ResolveToken_sum 16835.869077105075
consul_acl_ResolveToken_count 9370
# HELP consul_acl_authmethod_delete 
# TYPE consul_acl_authmethod_delete summary
consul_acl_authmethod_delete{quantile="0.5"} NaN
consul_acl_authmethod_delete{quantile="0.9"} NaN
consul_acl_authmethod_delete{quantile="0.99"} NaN
consul_acl_authmethod_delete_sum 0
consul_acl_authmethod_delete_count 0
# HELP consul_acl_authmethod_upsert 
# TYPE consul_acl_authmethod_upsert summary
consul_acl_authmethod_upsert{quantile="0.5"} NaN
consul_acl_authmethod_upsert{quantile="0.9"} NaN
consul_acl_authmethod_upsert{quantile="0.99"} NaN
consul_acl_authmethod_upsert_sum 0
consul_acl_authmethod_upsert_count 0
# HELP consul_acl_bindingrule_delete 
# TYPE consul_acl_bindingrule_delete summary
consul_acl_bindingrule_delete{quantile="0.5"} NaN
consul_acl_bindingrule_delete{quantile="0.9"} NaN
consul_acl_bindingrule_delete{quantile="0.99"} NaN
consul_acl_bindingrule_delete_sum 0
consul_acl_bindingrule_delete_count 0
# HELP consul_acl_bindingrule_upsert 
# TYPE consul_acl_bindingrule_upsert summary
consul_acl_bindingrule_upsert{quantile="0.5"} NaN
consul_acl_bindingrule_upsert{quantile="0.9"} NaN
consul_acl_bindingrule_upsert{quantile="0.99"} NaN
consul_acl_bindingrule_upsert_sum 0
consul_acl_bindingrule_upsert_count 0
# HELP consul_acl_blocked_check_deregistration Increments whenever a deregistration fails for a check (blocked by an ACL)
# TYPE consul_acl_blocked_check_deregistration counter
consul_acl_blocked_check_deregistration 0
# HELP consul_acl_blocked_check_registration Increments whenever a registration fails for a check (blocked by an ACL)
# TYPE consul_acl_blocked_check_registration counter
consul_acl_blocked_check_registration 0
# HELP consul_acl_blocked_node_registration Increments whenever a registration fails for a node (blocked by an ACL)
# TYPE consul_acl_blocked_node_registration counter
consul_acl_blocked_node_registration 0
# HELP consul_acl_blocked_service_deregistration Increments whenever a deregistration fails for a service (blocked by an ACL)
# TYPE consul_acl_blocked_service_deregistration counter
consul_acl_blocked_service_deregistration 0
# HELP consul_acl_blocked_service_registration Increments whenever a registration fails for a service (blocked by an ACL)
# TYPE consul_acl_blocked_service_registration counter
consul_acl_blocked_service_registration 0
# HELP consul_acl_login 
# TYPE consul_acl_login summary
consul_acl_login{quantile="0.5"} NaN
consul_acl_login{quantile="0.9"} NaN
consul_acl_login{quantile="0.99"} NaN
consul_acl_login_sum 0
consul_acl_login_count 0
# HELP consul_acl_logout 
# TYPE consul_acl_logout summary
consul_acl_logout{quantile="0.5"} NaN
consul_acl_logout{quantile="0.9"} NaN
consul_acl_logout{quantile="0.99"} NaN
consul_acl_logout_sum 0
consul_acl_logout_count 0
# HELP consul_acl_policy_delete 
# TYPE consul_acl_policy_delete summary
consul_acl_policy_delete{quantile="0.5"} NaN
consul_acl_policy_delete{quantile="0.9"} NaN
consul_acl_policy_delete{quantile="0.99"} NaN
consul_acl_policy_delete_sum 0
consul_acl_policy_delete_count 0
# HELP consul_acl_policy_upsert 
# TYPE consul_acl_policy_upsert summary
consul_acl_policy_upsert{quantile="0.5"} NaN
consul_acl_policy_upsert{quantile="0.9"} NaN
consul_acl_policy_upsert{quantile="0.99"} NaN
consul_acl_policy_upsert_sum 0
consul_acl_policy_upsert_count 0
# HELP consul_acl_role_delete 
# TYPE consul_acl_role_delete summary
consul_acl_role_delete{quantile="0.5"} NaN
consul_acl_role_delete{quantile="0.9"} NaN
consul_acl_role_delete{quantile="0.99"} NaN
consul_acl_role_delete_sum 0
consul_acl_role_delete_count 0
# HELP consul_acl_role_upsert 
# TYPE consul_acl_role_upsert summary
consul_acl_role_upsert{quantile="0.5"} NaN
consul_acl_role_upsert{quantile="0.9"} NaN
consul_acl_role_upsert{quantile="0.99"} NaN
consul_acl_role_upsert_sum 0
consul_acl_role_upsert_count 0
# HELP consul_acl_token_cache_hit Increments if Consul is able to resolve a token's identity, or a legacy token, from the cache.
# TYPE consul_acl_token_cache_hit counter
consul_acl_token_cache_hit 4617
# HELP consul_acl_token_cache_miss Increments if Consul cannot resolve a token's identity, or a legacy token, from the cache.
# TYPE consul_acl_token_cache_miss counter
consul_acl_token_cache_miss 4753
# HELP consul_acl_token_clone 
# TYPE consul_acl_token_clone summary
consul_acl_token_clone{quantile="0.5"} NaN
consul_acl_token_clone{quantile="0.9"} NaN
consul_acl_token_clone{quantile="0.99"} NaN
consul_acl_token_clone_sum 0
consul_acl_token_clone_count 0
# HELP consul_acl_token_delete 
# TYPE consul_acl_token_delete summary
consul_acl_token_delete{quantile="0.5"} NaN
consul_acl_token_delete{quantile="0.9"} NaN
consul_acl_token_delete{quantile="0.99"} NaN
consul_acl_token_delete_sum 0
consul_acl_token_delete_count 0
# HELP consul_acl_token_upsert 
# TYPE consul_acl_token_upsert summary
consul_acl_token_upsert{quantile="0.5"} NaN
consul_acl_token_upsert{quantile="0.9"} NaN
consul_acl_token_upsert{quantile="0.99"} NaN
consul_acl_token_upsert_sum 0
consul_acl_token_upsert_count 0
# HELP consul_agent_tls_cert_expiry Seconds until the agent tls certificate expires. Updated every hour
# TYPE consul_agent_tls_cert_expiry gauge
consul_agent_tls_cert_expiry 3.1255196e+07
# HELP consul_api_http Samples how long it takes to service the given HTTP request for the given verb and path.
# TYPE consul_api_http summary
consul_api_http{quantile="0.5"} NaN
consul_api_http{quantile="0.9"} NaN
consul_api_http{quantile="0.99"} NaN
consul_api_http_sum 0
consul_api_http_count 0
consul_api_http{method="GET",path="v1_agent_metrics",quantile="0.5"} 0.3934749960899353
consul_api_http{method="GET",path="v1_agent_metrics",quantile="0.9"} 0.3934749960899353
consul_api_http{method="GET",path="v1_agent_metrics",quantile="0.99"} 0.3934749960899353
consul_api_http_sum{method="GET",path="v1_agent_metrics"} 122243.81384652853
consul_api_http_count{method="GET",path="v1_agent_metrics"} 9369
consul_api_http{method="GET",path="v1_status_leader",quantile="0.5"} NaN
consul_api_http{method="GET",path="v1_status_leader",quantile="0.9"} NaN
consul_api_http{method="GET",path="v1_status_leader",quantile="0.99"} NaN
consul_api_http_sum{method="GET",path="v1_status_leader"} 67605.39766699076
consul_api_http_count{method="GET",path="v1_status_leader"} 28105
# HELP consul_cache_bypass Counts how many times a request bypassed the cache because no cache-key was provided.
# TYPE consul_cache_bypass counter
consul_cache_bypass 0
# HELP consul_cache_entries_count Represents the number of entries in this cache.
# TYPE consul_cache_entries_count gauge
consul_cache_entries_count 0
# HELP consul_cache_evict_expired Counts the number of expired entries that are evicted.
# TYPE consul_cache_evict_expired counter
consul_cache_evict_expired 0
# HELP consul_cache_fetch_error Counts the number of failed fetches by the cache.
# TYPE consul_cache_fetch_error counter
consul_cache_fetch_error 0
# HELP consul_cache_fetch_success Counts the number of successful fetches by the cache.
# TYPE consul_cache_fetch_success counter
consul_cache_fetch_success 0
# HELP consul_catalog_connect_not_found Increments for each connect-based catalog query where the given service could not be found.
# TYPE consul_catalog_connect_not_found counter
consul_catalog_connect_not_found 0
# HELP consul_catalog_connect_query Increments for each connect-based catalog query for the given service.
# TYPE consul_catalog_connect_query counter
consul_catalog_connect_query 0
# HELP consul_catalog_connect_query_tag Increments for each connect-based catalog query for the given service with the given tag.
# TYPE consul_catalog_connect_query_tag counter
consul_catalog_connect_query_tag 0
# HELP consul_catalog_connect_query_tags Increments for each connect-based catalog query for the given service with the given tags.
# TYPE consul_catalog_connect_query_tags counter
consul_catalog_connect_query_tags 0
# HELP consul_catalog_deregister Measures the time it takes to complete a catalog deregister operation.
# TYPE consul_catalog_deregister summary
consul_catalog_deregister{quantile="0.5"} NaN
consul_catalog_deregister{quantile="0.9"} NaN
consul_catalog_deregister{quantile="0.99"} NaN
consul_catalog_deregister_sum 0
consul_catalog_deregister_count 0
# HELP consul_catalog_register Measures the time it takes to complete a catalog register operation.
# TYPE consul_catalog_register summary
consul_catalog_register{quantile="0.5"} NaN
consul_catalog_register{quantile="0.9"} NaN
consul_catalog_register{quantile="0.99"} NaN
consul_catalog_register_sum 0
consul_catalog_register_count 0
# HELP consul_catalog_service_not_found Increments for each catalog query where the given service could not be found.
# TYPE consul_catalog_service_not_found counter
consul_catalog_service_not_found 0
# HELP consul_catalog_service_query Increments for each catalog query for the given service.
# TYPE consul_catalog_service_query counter
consul_catalog_service_query 0
# HELP consul_catalog_service_query_tag Increments for each catalog query for the given service with the given tag.
# TYPE consul_catalog_service_query_tag counter
consul_catalog_service_query_tag 0
# HELP consul_catalog_service_query_tags Increments for each catalog query for the given service with the given tags.
# TYPE consul_catalog_service_query_tags counter
consul_catalog_service_query_tags 0
# HELP consul_client_api_catalog_datacenters Increments whenever a Consul agent receives a request to list datacenters in the catalog.
# TYPE consul_client_api_catalog_datacenters counter
consul_client_api_catalog_datacenters 0
# HELP consul_client_api_catalog_deregister Increments whenever a Consul agent receives a catalog deregister request.
# TYPE consul_client_api_catalog_deregister counter
consul_client_api_catalog_deregister 0
# HELP consul_client_api_catalog_gateway_services Increments whenever a Consul agent receives a request to list services associated with a gateway.
# TYPE consul_client_api_catalog_gateway_services counter
consul_client_api_catalog_gateway_services 0
# HELP consul_client_api_catalog_node_service_list Increments whenever a Consul agent receives a request to list a node's registered services.
# TYPE consul_client_api_catalog_node_service_list counter
consul_client_api_catalog_node_service_list 0
# HELP consul_client_api_catalog_node_services Increments whenever a Consul agent successfully responds to a request to list nodes offering a service.
# TYPE consul_client_api_catalog_node_services counter
consul_client_api_catalog_node_services 0
# HELP consul_client_api_catalog_nodes Increments whenever a Consul agent receives a request to list nodes from the catalog.
# TYPE consul_client_api_catalog_nodes counter
consul_client_api_catalog_nodes 0
# HELP consul_client_api_catalog_register Increments whenever a Consul agent receives a catalog register request.
# TYPE consul_client_api_catalog_register counter
consul_client_api_catalog_register 0
# HELP consul_client_api_catalog_service_nodes Increments whenever a Consul agent receives a request to list nodes offering a service.
# TYPE consul_client_api_catalog_service_nodes counter
consul_client_api_catalog_service_nodes 0
# HELP consul_client_api_catalog_services Increments whenever a Consul agent receives a request to list services from the catalog.
# TYPE consul_client_api_catalog_services counter
consul_client_api_catalog_services 0
# HELP consul_client_api_error_catalog_service_nodes Increments whenever a Consul agent receives an RPC error for request to list nodes offering a service.
# TYPE consul_client_api_error_catalog_service_nodes counter
consul_client_api_error_catalog_service_nodes 0
# HELP consul_client_api_success_catalog_datacenters Increments whenever a Consul agent successfully responds to a request to list datacenters.
# TYPE consul_client_api_success_catalog_datacenters counter
consul_client_api_success_catalog_datacenters 0
# HELP consul_client_api_success_catalog_deregister Increments whenever a Consul agent successfully responds to a catalog deregister request.
# TYPE consul_client_api_success_catalog_deregister counter
consul_client_api_success_catalog_deregister 0
# HELP consul_client_api_success_catalog_gateway_services Increments whenever a Consul agent successfully responds to a request to list services associated with a gateway.
# TYPE consul_client_api_success_catalog_gateway_services counter
consul_client_api_success_catalog_gateway_services 0
# HELP consul_client_api_success_catalog_node_service_list Increments whenever a Consul agent successfully responds to a request to list a node's registered services.
# TYPE consul_client_api_success_catalog_node_service_list counter
consul_client_api_success_catalog_node_service_list 0
# HELP consul_client_api_success_catalog_node_services Increments whenever a Consul agent successfully responds to a request to list services in a node.
# TYPE consul_client_api_success_catalog_node_services counter
consul_client_api_success_catalog_node_services 0
# HELP consul_client_api_success_catalog_nodes Increments whenever a Consul agent successfully responds to a request to list nodes.
# TYPE consul_client_api_success_catalog_nodes counter
consul_client_api_success_catalog_nodes 0
# HELP consul_client_api_success_catalog_register Increments whenever a Consul agent successfully responds to a catalog register request.
# TYPE consul_client_api_success_catalog_register counter
consul_client_api_success_catalog_register 0
# HELP consul_client_api_success_catalog_service_nodes Increments whenever a Consul agent successfully responds to a request to list nodes offering a service.
# TYPE consul_client_api_success_catalog_service_nodes counter
consul_client_api_success_catalog_service_nodes 0
# HELP consul_client_api_success_catalog_services Increments whenever a Consul agent successfully responds to a request to list services.
# TYPE consul_client_api_success_catalog_services counter
consul_client_api_success_catalog_services 0
# HELP consul_client_rpc Increments whenever a Consul agent in client mode makes an RPC request to a Consul server.
# TYPE consul_client_rpc counter
consul_client_rpc 56334
# HELP consul_client_rpc_error_catalog_datacenters Increments whenever a Consul agent receives an RPC error for a request to list datacenters.
# TYPE consul_client_rpc_error_catalog_datacenters counter
consul_client_rpc_error_catalog_datacenters 0
# HELP consul_client_rpc_error_catalog_deregister Increments whenever a Consul agent receives an RPC error for a catalog deregister request.
# TYPE consul_client_rpc_error_catalog_deregister counter
consul_client_rpc_error_catalog_deregister 0
# HELP consul_client_rpc_error_catalog_gateway_services Increments whenever a Consul agent receives an RPC error for a request to list services associated with a gateway.
# TYPE consul_client_rpc_error_catalog_gateway_services counter
consul_client_rpc_error_catalog_gateway_services 0
# HELP consul_client_rpc_error_catalog_node_service_list Increments whenever a Consul agent receives an RPC error for request to list a node's registered services.
# TYPE consul_client_rpc_error_catalog_node_service_list counter
consul_client_rpc_error_catalog_node_service_list 0
# HELP consul_client_rpc_error_catalog_node_services Increments whenever a Consul agent receives an RPC error for a request to list services in a node.
# TYPE consul_client_rpc_error_catalog_node_services counter
consul_client_rpc_error_catalog_node_services 0
# HELP consul_client_rpc_error_catalog_nodes Increments whenever a Consul agent receives an RPC error for a request to list nodes.
# TYPE consul_client_rpc_error_catalog_nodes counter
consul_client_rpc_error_catalog_nodes 0
# HELP consul_client_rpc_error_catalog_register Increments whenever a Consul agent receives an RPC error for a catalog register request.
# TYPE consul_client_rpc_error_catalog_register counter
consul_client_rpc_error_catalog_register 0
# HELP consul_client_rpc_error_catalog_service_nodes Increments whenever a Consul agent receives an RPC error for a request to list nodes offering a service.
# TYPE consul_client_rpc_error_catalog_service_nodes counter
consul_client_rpc_error_catalog_service_nodes 0
# HELP consul_client_rpc_error_catalog_services Increments whenever a Consul agent receives an RPC error for a request to list services.
# TYPE consul_client_rpc_error_catalog_services counter
consul_client_rpc_error_catalog_services 0
# HELP consul_client_rpc_exceeded Increments whenever a Consul agent in client mode makes an RPC request to a Consul server gets rate limited by that agent's limits configuration.
# TYPE consul_client_rpc_exceeded counter
consul_client_rpc_exceeded 0
# HELP consul_client_rpc_failed Increments whenever a Consul agent in client mode makes an RPC request to a Consul server and fails.
# TYPE consul_client_rpc_failed counter
consul_client_rpc_failed 0
# HELP consul_consul_cache_bypass Deprecated - please use cache_bypass instead.
# TYPE consul_consul_cache_bypass counter
consul_consul_cache_bypass 0
# HELP consul_consul_cache_entries_count Deprecated - please use cache_entries_count instead.
# TYPE consul_consul_cache_entries_count gauge
consul_consul_cache_entries_count 0
# HELP consul_consul_cache_evict_expired Deprecated - please use cache_evict_expired instead.
# TYPE consul_consul_cache_evict_expired counter
consul_consul_cache_evict_expired 0
# HELP consul_consul_cache_fetch_error Deprecated - please use cache_fetch_error instead.
# TYPE consul_consul_cache_fetch_error counter
consul_consul_cache_fetch_error 0
# HELP consul_consul_cache_fetch_success Deprecated - please use cache_fetch_success instead.
# TYPE consul_consul_cache_fetch_success counter
consul_consul_cache_fetch_success 0
# HELP consul_consul_fsm_ca Deprecated - use fsm_ca instead
# TYPE consul_consul_fsm_ca summary
consul_consul_fsm_ca{quantile="0.5"} NaN
consul_consul_fsm_ca{quantile="0.9"} NaN
consul_consul_fsm_ca{quantile="0.99"} NaN
consul_consul_fsm_ca_sum 0
consul_consul_fsm_ca_count 0
# HELP consul_consul_fsm_intention Deprecated - use fsm_intention instead
# TYPE consul_consul_fsm_intention summary
consul_consul_fsm_intention{quantile="0.5"} NaN
consul_consul_fsm_intention{quantile="0.9"} NaN
consul_consul_fsm_intention{quantile="0.99"} NaN
consul_consul_fsm_intention_sum 0
consul_consul_fsm_intention_count 0
# HELP consul_consul_intention_apply Deprecated - please use intention_apply
# TYPE consul_consul_intention_apply summary
consul_consul_intention_apply{quantile="0.5"} NaN
consul_consul_intention_apply{quantile="0.9"} NaN
consul_consul_intention_apply{quantile="0.99"} NaN
consul_consul_intention_apply_sum 0
consul_consul_intention_apply_count 0
# HELP consul_consul_members_clients Deprecated - please use members_clients instead.
# TYPE consul_consul_members_clients gauge
consul_consul_members_clients 0
# HELP consul_consul_members_servers Deprecated - please use members_servers instead.
# TYPE consul_consul_members_servers gauge
consul_consul_members_servers 0
# HELP consul_consul_state_config_entries Deprecated - please use state_config_entries instead.
# TYPE consul_consul_state_config_entries gauge
consul_consul_state_config_entries 0
# HELP consul_consul_state_connect_instances Deprecated - please use state_connect_instances instead.
# TYPE consul_consul_state_connect_instances gauge
consul_consul_state_connect_instances 0
# HELP consul_consul_state_kv_entries Deprecated - please use kv_entries instead.
# TYPE consul_consul_state_kv_entries gauge
consul_consul_state_kv_entries 0
# HELP consul_consul_state_nodes Deprecated - please use state_nodes instead.
# TYPE consul_consul_state_nodes gauge
consul_consul_state_nodes 0
# HELP consul_consul_state_peerings Deprecated - please use state_peerings instead.
# TYPE consul_consul_state_peerings gauge
consul_consul_state_peerings 0
# HELP consul_consul_state_service_instances Deprecated - please use state_service_instances instead.
# TYPE consul_consul_state_service_instances gauge
consul_consul_state_service_instances 0
# HELP consul_consul_state_services Deprecated - please use state_services instead.
# TYPE consul_consul_state_services gauge
consul_consul_state_services 0
# HELP consul_federation_state_apply 
# TYPE consul_federation_state_apply summary
consul_federation_state_apply{quantile="0.5"} NaN
consul_federation_state_apply{quantile="0.9"} NaN
consul_federation_state_apply{quantile="0.99"} NaN
consul_federation_state_apply_sum 0
consul_federation_state_apply_count 0
# HELP consul_federation_state_get 
# TYPE consul_federation_state_get summary
consul_federation_state_get{quantile="0.5"} NaN
consul_federation_state_get{quantile="0.9"} NaN
consul_federation_state_get{quantile="0.99"} NaN
consul_federation_state_get_sum 0
consul_federation_state_get_count 0
# HELP consul_federation_state_list 
# TYPE consul_federation_state_list summary
consul_federation_state_list{quantile="0.5"} NaN
consul_federation_state_list{quantile="0.9"} NaN
consul_federation_state_list{quantile="0.99"} NaN
consul_federation_state_list_sum 0
consul_federation_state_list_count 0
# HELP consul_federation_state_list_mesh_gateways 
# TYPE consul_federation_state_list_mesh_gateways summary
consul_federation_state_list_mesh_gateways{quantile="0.5"} NaN
consul_federation_state_list_mesh_gateways{quantile="0.9"} NaN
consul_federation_state_list_mesh_gateways{quantile="0.99"} NaN
consul_federation_state_list_mesh_gateways_sum 0
consul_federation_state_list_mesh_gateways_count 0
# HELP consul_fsm_acl Measures the time it takes to apply the given ACL operation to the FSM.
# TYPE consul_fsm_acl summary
consul_fsm_acl{quantile="0.5"} NaN
consul_fsm_acl{quantile="0.9"} NaN
consul_fsm_acl{quantile="0.99"} NaN
consul_fsm_acl_sum 0
consul_fsm_acl_count 0
# HELP consul_fsm_acl_authmethod Measures the time it takes to apply an ACL authmethod operation to the FSM.
# TYPE consul_fsm_acl_authmethod summary
consul_fsm_acl_authmethod{quantile="0.5"} NaN
consul_fsm_acl_authmethod{quantile="0.9"} NaN
consul_fsm_acl_authmethod{quantile="0.99"} NaN
consul_fsm_acl_authmethod_sum 0
consul_fsm_acl_authmethod_count 0
# HELP consul_fsm_acl_bindingrule Measures the time it takes to apply an ACL binding rule operation to the FSM.
# TYPE consul_fsm_acl_bindingrule summary
consul_fsm_acl_bindingrule{quantile="0.5"} NaN
consul_fsm_acl_bindingrule{quantile="0.9"} NaN
consul_fsm_acl_bindingrule{quantile="0.99"} NaN
consul_fsm_acl_bindingrule_sum 0
consul_fsm_acl_bindingrule_count 0
# HELP consul_fsm_acl_policy Measures the time it takes to apply an ACL policy operation to the FSM.
# TYPE consul_fsm_acl_policy summary
consul_fsm_acl_policy{quantile="0.5"} NaN
consul_fsm_acl_policy{quantile="0.9"} NaN
consul_fsm_acl_policy{quantile="0.99"} NaN
consul_fsm_acl_policy_sum 0
consul_fsm_acl_policy_count 0
# HELP consul_fsm_acl_token Measures the time it takes to apply an ACL token operation to the FSM.
# TYPE consul_fsm_acl_token summary
consul_fsm_acl_token{quantile="0.5"} NaN
consul_fsm_acl_token{quantile="0.9"} NaN
consul_fsm_acl_token{quantile="0.99"} NaN
consul_fsm_acl_token_sum 0
consul_fsm_acl_token_count 0
# HELP consul_fsm_autopilot Measures the time it takes to apply the given autopilot update to the FSM.
# TYPE consul_fsm_autopilot summary
consul_fsm_autopilot{quantile="0.5"} NaN
consul_fsm_autopilot{quantile="0.9"} NaN
consul_fsm_autopilot{quantile="0.99"} NaN
consul_fsm_autopilot_sum 0
consul_fsm_autopilot_count 0
# HELP consul_fsm_ca Measures the time it takes to apply CA configuration operations to the FSM.
# TYPE consul_fsm_ca summary
consul_fsm_ca{quantile="0.5"} NaN
consul_fsm_ca{quantile="0.9"} NaN
consul_fsm_ca{quantile="0.99"} NaN
consul_fsm_ca_sum 0
consul_fsm_ca_count 0
# HELP consul_fsm_ca_leaf Measures the time it takes to apply an operation while signing a leaf certificate.
# TYPE consul_fsm_ca_leaf summary
consul_fsm_ca_leaf{quantile="0.5"} NaN
consul_fsm_ca_leaf{quantile="0.9"} NaN
consul_fsm_ca_leaf{quantile="0.99"} NaN
consul_fsm_ca_leaf_sum 0
consul_fsm_ca_leaf_count 0
# HELP consul_fsm_coordinate_batch_update Measures the time it takes to apply the given batch coordinate update to the FSM.
# TYPE consul_fsm_coordinate_batch_update summary
consul_fsm_coordinate_batch_update{quantile="0.5"} NaN
consul_fsm_coordinate_batch_update{quantile="0.9"} NaN
consul_fsm_coordinate_batch_update{quantile="0.99"} NaN
consul_fsm_coordinate_batch_update_sum 0
consul_fsm_coordinate_batch_update_count 0
# HELP consul_fsm_deregister Measures the time it takes to apply a catalog deregister operation to the FSM.
# TYPE consul_fsm_deregister summary
consul_fsm_deregister{quantile="0.5"} NaN
consul_fsm_deregister{quantile="0.9"} NaN
consul_fsm_deregister{quantile="0.99"} NaN
consul_fsm_deregister_sum 0
consul_fsm_deregister_count 0
# HELP consul_fsm_intention Measures the time it takes to apply an intention operation to the FSM.
# TYPE consul_fsm_intention summary
consul_fsm_intention{quantile="0.5"} NaN
consul_fsm_intention{quantile="0.9"} NaN
consul_fsm_intention{quantile="0.99"} NaN
consul_fsm_intention_sum 0
consul_fsm_intention_count 0
# HELP consul_fsm_kvs Measures the time it takes to apply the given KV operation to the FSM.
# TYPE consul_fsm_kvs summary
consul_fsm_kvs{quantile="0.5"} NaN
consul_fsm_kvs{quantile="0.9"} NaN
consul_fsm_kvs{quantile="0.99"} NaN
consul_fsm_kvs_sum 0
consul_fsm_kvs_count 0
# HELP consul_fsm_peering Measures the time it takes to apply a peering operation to the FSM.
# TYPE consul_fsm_peering summary
consul_fsm_peering{quantile="0.5"} NaN
consul_fsm_peering{quantile="0.9"} NaN
consul_fsm_peering{quantile="0.99"} NaN
consul_fsm_peering_sum 0
consul_fsm_peering_count 0
# HELP consul_fsm_persist Measures the time it takes to persist the FSM to a raft snapshot.
# TYPE consul_fsm_persist summary
consul_fsm_persist{quantile="0.5"} NaN
consul_fsm_persist{quantile="0.9"} NaN
consul_fsm_persist{quantile="0.99"} NaN
consul_fsm_persist_sum 0
consul_fsm_persist_count 0
# HELP consul_fsm_prepared_query Measures the time it takes to apply the given prepared query update operation to the FSM.
# TYPE consul_fsm_prepared_query summary
consul_fsm_prepared_query{quantile="0.5"} NaN
consul_fsm_prepared_query{quantile="0.9"} NaN
consul_fsm_prepared_query{quantile="0.99"} NaN
consul_fsm_prepared_query_sum 0
consul_fsm_prepared_query_count 0
# HELP consul_fsm_register Measures the time it takes to apply a catalog register operation to the FSM.
# TYPE consul_fsm_register summary
consul_fsm_register{quantile="0.5"} NaN
consul_fsm_register{quantile="0.9"} NaN
consul_fsm_register{quantile="0.99"} NaN
consul_fsm_register_sum 0
consul_fsm_register_count 0
# HELP consul_fsm_session Measures the time it takes to apply the given session operation to the FSM.
# TYPE consul_fsm_session summary
consul_fsm_session{quantile="0.5"} NaN
consul_fsm_session{quantile="0.9"} NaN
consul_fsm_session{quantile="0.99"} NaN
consul_fsm_session_sum 0
consul_fsm_session_count 0
# HELP consul_fsm_system_metadata Measures the time it takes to apply a system metadata operation to the FSM.
# TYPE consul_fsm_system_metadata summary
consul_fsm_system_metadata{quantile="0.5"} NaN
consul_fsm_system_metadata{quantile="0.9"} NaN
consul_fsm_system_metadata{quantile="0.99"} NaN
consul_fsm_system_metadata_sum 0
consul_fsm_system_metadata_count 0
# HELP consul_fsm_tombstone Measures the time it takes to apply the given tombstone operation to the FSM.
# TYPE consul_fsm_tombstone summary
consul_fsm_tombstone{quantile="0.5"} NaN
consul_fsm_tombstone{quantile="0.9"} NaN
consul_fsm_tombstone{quantile="0.99"} NaN
consul_fsm_tombstone_sum 0
consul_fsm_tombstone_count 0
# HELP consul_fsm_txn Measures the time it takes to apply the given transaction update to the FSM.
# TYPE consul_fsm_txn summary
consul_fsm_txn{quantile="0.5"} NaN
consul_fsm_txn{quantile="0.9"} NaN
consul_fsm_txn{quantile="0.99"} NaN
consul_fsm_txn_sum 0
consul_fsm_txn_count 0
# HELP consul_grpc_client_connection_count Counts the number of new gRPC connections opened by the client agent to a Consul server.
# TYPE consul_grpc_client_connection_count counter
consul_grpc_client_connection_count 0
consul_grpc_client_connection_count{server_type="internal"} 1264
# HELP consul_grpc_client_connections Measures the number of active gRPC connections open from the client agent to any Consul servers.
# TYPE consul_grpc_client_connections gauge
consul_grpc_client_connections 0
consul_grpc_client_connections{server_type="internal"} 1
# HELP consul_grpc_client_request_count Counts the number of gRPC requests made by the client agent to a Consul server.
# TYPE consul_grpc_client_request_count counter
consul_grpc_client_request_count 0
# HELP consul_grpc_server_connection_count Counts the number of new gRPC connections received by the server.
# TYPE consul_grpc_server_connection_count counter
consul_grpc_server_connection_count 0
# HELP consul_grpc_server_connections Measures the number of active gRPC connections open on the server.
# TYPE consul_grpc_server_connections gauge
consul_grpc_server_connections 0
# HELP consul_grpc_server_request_count Counts the number of gRPC requests received by the server.
# TYPE consul_grpc_server_request_count counter
consul_grpc_server_request_count 0
# HELP consul_grpc_server_stream_count Counts the number of new gRPC streams received by the server.
# TYPE consul_grpc_server_stream_count counter
consul_grpc_server_stream_count 0
# HELP consul_grpc_server_streams Measures the number of active gRPC streams handled by the server.
# TYPE consul_grpc_server_streams gauge
consul_grpc_server_streams 0
# HELP consul_intention_apply 
# TYPE consul_intention_apply summary
consul_intention_apply{quantile="0.5"} NaN
consul_intention_apply{quantile="0.9"} NaN
consul_intention_apply{quantile="0.99"} NaN
consul_intention_apply_sum 0
consul_intention_apply_count 0
# HELP consul_kvs_apply Measures the time it takes to complete an update to the KV store.
# TYPE consul_kvs_apply summary
consul_kvs_apply{quantile="0.5"} NaN
consul_kvs_apply{quantile="0.9"} NaN
consul_kvs_apply{quantile="0.99"} NaN
consul_kvs_apply_sum 0
consul_kvs_apply_count 0
# HELP consul_leader_barrier Measures the time spent waiting for the raft barrier upon gaining leadership.
# TYPE consul_leader_barrier summary
consul_leader_barrier{quantile="0.5"} NaN
consul_leader_barrier{quantile="0.9"} NaN
consul_leader_barrier{quantile="0.99"} NaN
consul_leader_barrier_sum 0
consul_leader_barrier_count 0
# HELP consul_leader_reapTombstones Measures the time spent clearing tombstones.
# TYPE consul_leader_reapTombstones summary
consul_leader_reapTombstones{quantile="0.5"} NaN
consul_leader_reapTombstones{quantile="0.9"} NaN
consul_leader_reapTombstones{quantile="0.99"} NaN
consul_leader_reapTombstones_sum 0
consul_leader_reapTombstones_count 0
# HELP consul_leader_reconcile Measures the time spent updating the raft store from the serf member information.
# TYPE consul_leader_reconcile summary
consul_leader_reconcile{quantile="0.5"} NaN
consul_leader_reconcile{quantile="0.9"} NaN
consul_leader_reconcile{quantile="0.99"} NaN
consul_leader_reconcile_sum 0
consul_leader_reconcile_count 0
# HELP consul_leader_reconcileMember Measures the time spent updating the raft store for a single serf member's information.
# TYPE consul_leader_reconcileMember summary
consul_leader_reconcileMember{quantile="0.5"} NaN
consul_leader_reconcileMember{quantile="0.9"} NaN
consul_leader_reconcileMember{quantile="0.99"} NaN
consul_leader_reconcileMember_sum 0
consul_leader_reconcileMember_count 0
# HELP consul_leader_replication_acl_policies_index Tracks the index of ACL policies in the primary that the secondary has successfully replicated
# TYPE consul_leader_replication_acl_policies_index gauge
consul_leader_replication_acl_policies_index 0
# HELP consul_leader_replication_acl_policies_status Tracks the current health of ACL policy replication on the leader
# TYPE consul_leader_replication_acl_policies_status gauge
consul_leader_replication_acl_policies_status 0
# HELP consul_leader_replication_acl_roles_index Tracks the index of ACL roles in the primary that the secondary has successfully replicated
# TYPE consul_leader_replication_acl_roles_index gauge
consul_leader_replication_acl_roles_index 0
# HELP consul_leader_replication_acl_roles_status Tracks the current health of ACL role replication on the leader
# TYPE consul_leader_replication_acl_roles_status gauge
consul_leader_replication_acl_roles_status 0
# HELP consul_leader_replication_acl_tokens_index Tracks the index of ACL tokens in the primary that the secondary has successfully replicated
# TYPE consul_leader_replication_acl_tokens_index gauge
consul_leader_replication_acl_tokens_index 0
# HELP consul_leader_replication_acl_tokens_status Tracks the current health of ACL token replication on the leader
# TYPE consul_leader_replication_acl_tokens_status gauge
consul_leader_replication_acl_tokens_status 0
# HELP consul_leader_replication_config_entries_index Tracks the index of config entries in the primary that the secondary has successfully replicated
# TYPE consul_leader_replication_config_entries_index gauge
consul_leader_replication_config_entries_index 0
# HELP consul_leader_replication_config_entries_status Tracks the current health of config entry replication on the leader
# TYPE consul_leader_replication_config_entries_status gauge
consul_leader_replication_config_entries_status 0
# HELP consul_leader_replication_federation_state_index Tracks the index of federation states in the primary that the secondary has successfully replicated
# TYPE consul_leader_replication_federation_state_index gauge
consul_leader_replication_federation_state_index 0
# HELP consul_leader_replication_federation_state_status Tracks the current health of federation state replication on the leader
# TYPE consul_leader_replication_federation_state_status gauge
consul_leader_replication_federation_state_status 0
# HELP consul_leader_replication_namespaces_index Tracks the index of federation states in the primary that the secondary has successfully replicated
# TYPE consul_leader_replication_namespaces_index gauge
consul_leader_replication_namespaces_index 0
# HELP consul_leader_replication_namespaces_status Tracks the current health of federation state replication on the leader
# TYPE consul_leader_replication_namespaces_status gauge
consul_leader_replication_namespaces_status 0
# HELP consul_memberlist_gossip consul_memberlist_gossip
# TYPE consul_memberlist_gossip summary
consul_memberlist_gossip{network="lan",quantile="0.5"} 0.024581000208854675
consul_memberlist_gossip{network="lan",quantile="0.9"} 0.03684600070118904
consul_memberlist_gossip{network="lan",quantile="0.99"} 0.05559699982404709
consul_memberlist_gossip_sum{network="lan"} 38488.5026860768
consul_memberlist_gossip_count{network="lan"} 1.405334e+06
# HELP consul_memberlist_node_instances consul_memberlist_node_instances
# TYPE consul_memberlist_node_instances gauge
consul_memberlist_node_instances{network="lan",node_state="alive"} 9
consul_memberlist_node_instances{network="lan",node_state="dead"} 0
consul_memberlist_node_instances{network="lan",node_state="left"} 0
consul_memberlist_node_instances{network="lan",node_state="suspect"} 0
# HELP consul_memberlist_probeNode consul_memberlist_probeNode
# TYPE consul_memberlist_probeNode summary
consul_memberlist_probeNode{network="lan",quantile="0.5"} 1.5850770473480225
consul_memberlist_probeNode{network="lan",quantile="0.9"} 2.4520740509033203
consul_memberlist_probeNode{network="lan",quantile="0.99"} 9.382650375366211
consul_memberlist_probeNode_sum{network="lan"} 535711.760569483
consul_memberlist_probeNode_count{network="lan"} 281066
# HELP consul_memberlist_pushPullNode consul_memberlist_pushPullNode
# TYPE consul_memberlist_pushPullNode summary
consul_memberlist_pushPullNode{network="lan",quantile="0.5"} NaN
consul_memberlist_pushPullNode{network="lan",quantile="0.9"} NaN
consul_memberlist_pushPullNode{network="lan",quantile="0.99"} NaN
consul_memberlist_pushPullNode_sum{network="lan"} 36623.60043847561
consul_memberlist_pushPullNode_count{network="lan"} 9369
# HELP consul_memberlist_queue_broadcasts consul_memberlist_queue_broadcasts
# TYPE consul_memberlist_queue_broadcasts summary
consul_memberlist_queue_broadcasts{network="lan",quantile="0.5"} NaN
consul_memberlist_queue_broadcasts{network="lan",quantile="0.9"} NaN
consul_memberlist_queue_broadcasts{network="lan",quantile="0.99"} NaN
consul_memberlist_queue_broadcasts_sum{network="lan"} 0
consul_memberlist_queue_broadcasts_count{network="lan"} 9368
# HELP consul_memberlist_size_local consul_memberlist_size_local
# TYPE consul_memberlist_size_local gauge
consul_memberlist_size_local{network="lan"} 2.208582144e+09
# HELP consul_memberlist_tcp_accept consul_memberlist_tcp_accept
# TYPE consul_memberlist_tcp_accept counter
consul_memberlist_tcp_accept{network="lan"} 18624
# HELP consul_memberlist_tcp_connect consul_memberlist_tcp_connect
# TYPE consul_memberlist_tcp_connect counter
consul_memberlist_tcp_connect{network="lan"} 9369
# HELP consul_memberlist_tcp_sent consul_memberlist_tcp_sent
# TYPE consul_memberlist_tcp_sent counter
consul_memberlist_tcp_sent{network="lan"} 2.7873918e+07
# HELP consul_memberlist_udp_received consul_memberlist_udp_received
# TYPE consul_memberlist_udp_received counter
consul_memberlist_udp_received{network="lan"} 6.9972902e+07
# HELP consul_memberlist_udp_sent consul_memberlist_udp_sent
# TYPE consul_memberlist_udp_sent counter
consul_memberlist_udp_sent{network="lan"} 6.997255e+07
# HELP consul_members_clients Measures the current number of client agents registered with Consul. It is only emitted by Consul servers. Added in v1.9.6.
# TYPE consul_members_clients gauge
consul_members_clients 0
# HELP consul_members_servers Measures the current number of server agents registered with Consul. It is only emitted by Consul servers. Added in v1.9.6.
# TYPE consul_members_servers gauge
consul_members_servers 0
# HELP consul_prepared_query_apply Measures the time it takes to apply a prepared query update.
# TYPE consul_prepared_query_apply summary
consul_prepared_query_apply{quantile="0.5"} NaN
consul_prepared_query_apply{quantile="0.9"} NaN
consul_prepared_query_apply{quantile="0.99"} NaN
consul_prepared_query_apply_sum 0
consul_prepared_query_apply_count 0
# HELP consul_prepared_query_execute Measures the time it takes to process a prepared query execute request.
# TYPE consul_prepared_query_execute summary
consul_prepared_query_execute{quantile="0.5"} NaN
consul_prepared_query_execute{quantile="0.9"} NaN
consul_prepared_query_execute{quantile="0.99"} NaN
consul_prepared_query_execute_sum 0
consul_prepared_query_execute_count 0
# HELP consul_prepared_query_execute_remote Measures the time it takes to process a prepared query execute request that was forwarded to another datacenter.
# TYPE consul_prepared_query_execute_remote summary
consul_prepared_query_execute_remote{quantile="0.5"} NaN
consul_prepared_query_execute_remote{quantile="0.9"} NaN
consul_prepared_query_execute_remote{quantile="0.99"} NaN
consul_prepared_query_execute_remote_sum 0
consul_prepared_query_execute_remote_count 0
# HELP consul_prepared_query_explain Measures the time it takes to process a prepared query explain request.
# TYPE consul_prepared_query_explain summary
consul_prepared_query_explain{quantile="0.5"} NaN
consul_prepared_query_explain{quantile="0.9"} NaN
consul_prepared_query_explain{quantile="0.99"} NaN
consul_prepared_query_explain_sum 0
consul_prepared_query_explain_count 0
# HELP consul_raft_applied_index Represents the raft applied index.
# TYPE consul_raft_applied_index gauge
consul_raft_applied_index 0
# HELP consul_raft_apply This counts the number of Raft transactions occurring over the interval.
# TYPE consul_raft_apply counter
consul_raft_apply 0
# HELP consul_raft_commitTime This measures the time it takes to commit a new entry to the Raft log on the leader.
# TYPE consul_raft_commitTime summary
consul_raft_commitTime{quantile="0.5"} NaN
consul_raft_commitTime{quantile="0.9"} NaN
consul_raft_commitTime{quantile="0.99"} NaN
consul_raft_commitTime_sum 0
consul_raft_commitTime_count 0
# HELP consul_raft_fsm_lastRestoreDuration This measures how long the last FSM restore (from disk or leader) took.
# TYPE consul_raft_fsm_lastRestoreDuration gauge
consul_raft_fsm_lastRestoreDuration 0
# HELP consul_raft_last_index Represents the raft last index.
# TYPE consul_raft_last_index gauge
consul_raft_last_index 0
# HELP consul_raft_leader_lastContact Measures the time since the leader was last able to contact the follower nodes when checking its leader lease.
# TYPE consul_raft_leader_lastContact summary
consul_raft_leader_lastContact{quantile="0.5"} NaN
consul_raft_leader_lastContact{quantile="0.9"} NaN
consul_raft_leader_lastContact{quantile="0.99"} NaN
consul_raft_leader_lastContact_sum 0
consul_raft_leader_lastContact_count 0
# HELP consul_raft_leader_oldestLogAge This measures how old the oldest log in the leader's log store is.
# TYPE consul_raft_leader_oldestLogAge gauge
consul_raft_leader_oldestLogAge 0
# HELP consul_raft_rpc_installSnapshot Measures the time it takes the raft leader to install a snapshot on a follower that is catching up after being down or has just joined the cluster.
# TYPE consul_raft_rpc_installSnapshot summary
consul_raft_rpc_installSnapshot{quantile="0.5"} NaN
consul_raft_rpc_installSnapshot{quantile="0.9"} NaN
consul_raft_rpc_installSnapshot{quantile="0.99"} NaN
consul_raft_rpc_installSnapshot_sum 0
consul_raft_rpc_installSnapshot_count 0
# HELP consul_raft_snapshot_persist Measures the time it takes raft to write a new snapshot to disk.
# TYPE consul_raft_snapshot_persist summary
consul_raft_snapshot_persist{quantile="0.5"} NaN
consul_raft_snapshot_persist{quantile="0.9"} NaN
consul_raft_snapshot_persist{quantile="0.99"} NaN
consul_raft_snapshot_persist_sum 0
consul_raft_snapshot_persist_count 0
# HELP consul_raft_state_candidate This increments whenever a Consul server starts an election.
# TYPE consul_raft_state_candidate counter
consul_raft_state_candidate 0
# HELP consul_raft_state_leader This increments whenever a Consul server becomes a leader.
# TYPE consul_raft_state_leader counter
consul_raft_state_leader 0
# HELP consul_rpc_accept_conn Increments when a server accepts an RPC connection.
# TYPE consul_rpc_accept_conn counter
consul_rpc_accept_conn 0
# HELP consul_rpc_consistentRead Measures the time spent confirming that a consistent read can be performed.
# TYPE consul_rpc_consistentRead summary
consul_rpc_consistentRead{quantile="0.5"} NaN
consul_rpc_consistentRead{quantile="0.9"} NaN
consul_rpc_consistentRead{quantile="0.99"} NaN
consul_rpc_consistentRead_sum 0
consul_rpc_consistentRead_count 0
# HELP consul_rpc_cross_dc Increments when a server sends a (potentially blocking) cross datacenter RPC query.
# TYPE consul_rpc_cross_dc counter
consul_rpc_cross_dc 0
# HELP consul_rpc_queries_blocking Shows the current number of in-flight blocking queries the server is handling.
# TYPE consul_rpc_queries_blocking gauge
consul_rpc_queries_blocking 0
# HELP consul_rpc_query Increments when a server receives a read request, indicating the rate of new read queries.
# TYPE consul_rpc_query counter
consul_rpc_query 0
# HELP consul_rpc_raft_handoff Increments when a server accepts a Raft-related RPC connection.
# TYPE consul_rpc_raft_handoff counter
consul_rpc_raft_handoff 0
# HELP consul_rpc_request Increments when a server receives a Consul-related RPC request.
# TYPE consul_rpc_request counter
consul_rpc_request 0
# HELP consul_rpc_request_error Increments when a server returns an error from an RPC request.
# TYPE consul_rpc_request_error counter
consul_rpc_request_error 0
# HELP consul_runtime_alloc_bytes consul_runtime_alloc_bytes
# TYPE consul_runtime_alloc_bytes gauge
consul_runtime_alloc_bytes 1.7193976e+07
# HELP consul_runtime_free_count consul_runtime_free_count
# TYPE consul_runtime_free_count gauge
consul_runtime_free_count 3.33539072e+08
# HELP consul_runtime_gc_pause_ns consul_runtime_gc_pause_ns
# TYPE consul_runtime_gc_pause_ns summary
consul_runtime_gc_pause_ns{quantile="0.5"} 135596
consul_runtime_gc_pause_ns{quantile="0.9"} 135596
consul_runtime_gc_pause_ns{quantile="0.99"} 135596
consul_runtime_gc_pause_ns_sum 1.913592356e+09
consul_runtime_gc_pause_ns_count 7093
# HELP consul_runtime_heap_objects consul_runtime_heap_objects
# TYPE consul_runtime_heap_objects gauge
consul_runtime_heap_objects 43677
# HELP consul_runtime_malloc_count consul_runtime_malloc_count
# TYPE consul_runtime_malloc_count gauge
consul_runtime_malloc_count 3.33582752e+08
# HELP consul_runtime_num_goroutines consul_runtime_num_goroutines
# TYPE consul_runtime_num_goroutines gauge
consul_runtime_num_goroutines 55
# HELP consul_runtime_sys_bytes consul_runtime_sys_bytes
# TYPE consul_runtime_sys_bytes gauge
consul_runtime_sys_bytes 5.6210696e+07
# HELP consul_runtime_total_gc_pause_ns consul_runtime_total_gc_pause_ns
# TYPE consul_runtime_total_gc_pause_ns gauge
consul_runtime_total_gc_pause_ns 1.91359232e+09
# HELP consul_runtime_total_gc_runs consul_runtime_total_gc_runs
# TYPE consul_runtime_total_gc_runs gauge
consul_runtime_total_gc_runs 7093
# HELP consul_serf_coordinate_adjustment_ms consul_serf_coordinate_adjustment_ms
# TYPE consul_serf_coordinate_adjustment_ms summary
consul_serf_coordinate_adjustment_ms{network="lan",quantile="0.5"} 0.06316500157117844
consul_serf_coordinate_adjustment_ms{network="lan",quantile="0.9"} 0.10257299989461899
consul_serf_coordinate_adjustment_ms{network="lan",quantile="0.99"} 0.13280600309371948
consul_serf_coordinate_adjustment_ms_sum{network="lan"} 54553.18788221612
consul_serf_coordinate_adjustment_ms_count{network="lan"} 281066
# HELP consul_serf_queue_Event consul_serf_queue_Event
# TYPE consul_serf_queue_Event summary
consul_serf_queue_Event{network="lan",quantile="0.5"} NaN
consul_serf_queue_Event{network="lan",quantile="0.9"} NaN
consul_serf_queue_Event{network="lan",quantile="0.99"} NaN
consul_serf_queue_Event_sum{network="lan"} 0
consul_serf_queue_Event_count{network="lan"} 9368
# HELP consul_serf_queue_Intent consul_serf_queue_Intent
# TYPE consul_serf_queue_Intent summary
consul_serf_queue_Intent{network="lan",quantile="0.5"} NaN
consul_serf_queue_Intent{network="lan",quantile="0.9"} NaN
consul_serf_queue_Intent{network="lan",quantile="0.99"} NaN
consul_serf_queue_Intent_sum{network="lan"} 0
consul_serf_queue_Intent_count{network="lan"} 9368
# HELP consul_serf_queue_Query consul_serf_queue_Query
# TYPE consul_serf_queue_Query summary
consul_serf_queue_Query{network="lan",quantile="0.5"} NaN
consul_serf_queue_Query{network="lan",quantile="0.9"} NaN
consul_serf_queue_Query{network="lan",quantile="0.99"} NaN
consul_serf_queue_Query_sum{network="lan"} 0
consul_serf_queue_Query_count{network="lan"} 9368
# HELP consul_server_isLeader Tracks if the server is a leader.
# TYPE consul_server_isLeader gauge
consul_server_isLeader 0
# HELP consul_session_apply Measures the time spent applying a session update.
# TYPE consul_session_apply summary
consul_session_apply{quantile="0.5"} NaN
consul_session_apply{quantile="0.9"} NaN
consul_session_apply{quantile="0.99"} NaN
consul_session_apply_sum 0
consul_session_apply_count 0
# HELP consul_session_renew Measures the time spent renewing a session.
# TYPE consul_session_renew summary
consul_session_renew{quantile="0.5"} NaN
consul_session_renew{quantile="0.9"} NaN
consul_session_renew{quantile="0.99"} NaN
consul_session_renew_sum 0
consul_session_renew_count 0
# HELP consul_session_ttl_active Tracks the active number of sessions being tracked.
# TYPE consul_session_ttl_active gauge
consul_session_ttl_active 0
# HELP consul_session_ttl_invalidate Measures the time spent invalidating an expired session.
# TYPE consul_session_ttl_invalidate summary
consul_session_ttl_invalidate{quantile="0.5"} NaN
consul_session_ttl_invalidate{quantile="0.9"} NaN
consul_session_ttl_invalidate{quantile="0.99"} NaN
consul_session_ttl_invalidate_sum 0
consul_session_ttl_invalidate_count 0
# HELP consul_state_config_entries Measures the current number of unique configuration entries registered with Consul, labeled by Kind. It is only emitted by Consul servers. Added in v1.10.4.
# TYPE consul_state_config_entries gauge
consul_state_config_entries 0
# HELP consul_state_connect_instances Measures the current number of unique connect service instances registered with Consul, labeled by Kind. It is only emitted by Consul servers. Added in v1.10.4.
# TYPE consul_state_connect_instances gauge
consul_state_connect_instances 0
# HELP consul_state_kv_entries Measures the current number of entries in the Consul KV store. It is only emitted by Consul servers. Added in v1.10.3.
# TYPE consul_state_kv_entries gauge
consul_state_kv_entries 0
# HELP consul_state_nodes Measures the current number of nodes registered with Consul. It is only emitted by Consul servers. Added in v1.9.0.
# TYPE consul_state_nodes gauge
consul_state_nodes 0
# HELP consul_state_peerings Measures the current number of peerings registered with Consul. It is only emitted by Consul servers. Added in v1.13.0.
# TYPE consul_state_peerings gauge
consul_state_peerings 0
# HELP consul_state_service_instances Measures the current number of unique services registered with Consul, based on service name. It is only emitted by Consul servers. Added in v1.9.0.
# TYPE consul_state_service_instances gauge
consul_state_service_instances 0
# HELP consul_state_services Measures the current number of unique services registered with Consul, based on service name. It is only emitted by Consul servers. Added in v1.9.0.
# TYPE consul_state_services gauge
consul_state_services 0
# HELP consul_txn_apply Measures the time spent applying a transaction operation.
# TYPE consul_txn_apply summary
consul_txn_apply{quantile="0.5"} NaN
consul_txn_apply{quantile="0.9"} NaN
consul_txn_apply{quantile="0.99"} NaN
consul_txn_apply_sum 0
consul_txn_apply_count 0
# HELP consul_txn_read Measures the time spent returning a read transaction.
# TYPE consul_txn_read summary
consul_txn_read{quantile="0.5"} NaN
consul_txn_read{quantile="0.9"} NaN
consul_txn_read{quantile="0.99"} NaN
consul_txn_read_sum 0
consul_txn_read_count 0
# HELP consul_version Represents the Consul version.
# TYPE consul_version gauge
consul_version 0
# HELP consul_xds_server_streamDrained Counts the number of xDS streams that are drained when rebalancing the load between servers.
# TYPE consul_xds_server_streamDrained counter
consul_xds_server_streamDrained 0
# HELP consul_xds_server_streamStart Measures the time in milliseconds after an xDS stream is opened until xDS resources are first generated for the stream.
# TYPE consul_xds_server_streamStart summary
consul_xds_server_streamStart{quantile="0.5"} NaN
consul_xds_server_streamStart{quantile="0.9"} NaN
consul_xds_server_streamStart{quantile="0.99"} NaN
consul_xds_server_streamStart_sum 0
consul_xds_server_streamStart_count 0
# HELP consul_xds_server_streams Measures the number of active xDS streams handled by the server split by protocol version.
# TYPE consul_xds_server_streams gauge
consul_xds_server_streams 0
# HELP go_gc_duration_seconds A summary of the pause duration of garbage collection cycles.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 6.8354e-05
go_gc_duration_seconds{quantile="0.25"} 0.000108932
go_gc_duration_seconds{quantile="0.5"} 0.000136016
go_gc_duration_seconds{quantile="0.75"} 0.000192746
go_gc_duration_seconds{quantile="1"} 0.015105933
go_gc_duration_seconds_sum 1.913592357
go_gc_duration_seconds_count 7093
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 61
# HELP go_info Information about the Go environment.
# TYPE go_info gauge
go_info{version="go1.19.4"} 1
# HELP go_memstats_alloc_bytes Number of bytes allocated and still in use.
# TYPE go_memstats_alloc_bytes gauge
go_memstats_alloc_bytes 1.743812e+07
# HELP go_memstats_alloc_bytes_total Total number of bytes allocated, even if freed.
# TYPE go_memstats_alloc_bytes_total counter
go_memstats_alloc_bytes_total 1.0422934356e+11
# HELP go_memstats_buck_hash_sys_bytes Number of bytes used by the profiling bucket hash table.
# TYPE go_memstats_buck_hash_sys_bytes gauge
go_memstats_buck_hash_sys_bytes 2.277207e+06
# HELP go_memstats_frees_total Total number of frees.
# TYPE go_memstats_frees_total counter
go_memstats_frees_total 3.33539112e+08
# HELP go_memstats_gc_cpu_fraction The fraction of this program's available CPU time used by the GC since the program started.
# TYPE go_memstats_gc_cpu_fraction gauge
go_memstats_gc_cpu_fraction 8.146458759393804e-05
# HELP go_memstats_gc_sys_bytes Number of bytes used for garbage collection system metadata.
# TYPE go_memstats_gc_sys_bytes gauge
go_memstats_gc_sys_bytes 1.0240104e+07
# HELP go_memstats_heap_alloc_bytes Number of heap bytes allocated and still in use.
# TYPE go_memstats_heap_alloc_bytes gauge
go_memstats_heap_alloc_bytes 1.743812e+07
# HELP go_memstats_heap_idle_bytes Number of heap bytes waiting to be used.
# TYPE go_memstats_heap_idle_bytes gauge
go_memstats_heap_idle_bytes 1.9775488e+07
# HELP go_memstats_heap_inuse_bytes Number of heap bytes that are in use.
# TYPE go_memstats_heap_inuse_bytes gauge
go_memstats_heap_inuse_bytes 2.1086208e+07
# HELP go_memstats_heap_objects Number of allocated objects.
# TYPE go_memstats_heap_objects gauge
go_memstats_heap_objects 44377
# HELP go_memstats_heap_released_bytes Number of heap bytes released to OS.
# TYPE go_memstats_heap_released_bytes gauge
go_memstats_heap_released_bytes 9.05216e+06
# HELP go_memstats_heap_sys_bytes Number of heap bytes obtained from system.
# TYPE go_memstats_heap_sys_bytes gauge
go_memstats_heap_sys_bytes 4.0861696e+07
# HELP go_memstats_last_gc_time_seconds Number of seconds since 1970 of last garbage collection.
# TYPE go_memstats_last_gc_time_seconds gauge
go_memstats_last_gc_time_seconds 1.6835651035864115e+09
# HELP go_memstats_lookups_total Total number of pointer lookups.
# TYPE go_memstats_lookups_total counter
go_memstats_lookups_total 0
# HELP go_memstats_mallocs_total Total number of mallocs.
# TYPE go_memstats_mallocs_total counter
go_memstats_mallocs_total 3.33583489e+08
# HELP go_memstats_mcache_inuse_bytes Number of bytes in use by mcache structures.
# TYPE go_memstats_mcache_inuse_bytes gauge
go_memstats_mcache_inuse_bytes 6000
# HELP go_memstats_mcache_sys_bytes Number of bytes used for mcache structures obtained from system.
# TYPE go_memstats_mcache_sys_bytes gauge
go_memstats_mcache_sys_bytes 15600
# HELP go_memstats_mspan_inuse_bytes Number of bytes in use by mspan structures.
# TYPE go_memstats_mspan_inuse_bytes gauge
go_memstats_mspan_inuse_bytes 258048
# HELP go_memstats_mspan_sys_bytes Number of bytes used for mspan structures obtained from system.
# TYPE go_memstats_mspan_sys_bytes gauge
go_memstats_mspan_sys_bytes 341712
# HELP go_memstats_next_gc_bytes Number of heap bytes when next garbage collection will take place.
# TYPE go_memstats_next_gc_bytes gauge
go_memstats_next_gc_bytes 3.1364352e+07
# HELP go_memstats_other_sys_bytes Number of bytes used for other system allocations.
# TYPE go_memstats_other_sys_bytes gauge
go_memstats_other_sys_bytes 1.393033e+06
# HELP go_memstats_stack_inuse_bytes Number of bytes in use by the stack allocator.
# TYPE go_memstats_stack_inuse_bytes gauge
go_memstats_stack_inuse_bytes 1.081344e+06
# HELP go_memstats_stack_sys_bytes Number of bytes obtained from system for stack allocator.
# TYPE go_memstats_stack_sys_bytes gauge
go_memstats_stack_sys_bytes 1.081344e+06
# HELP go_memstats_sys_bytes Number of bytes obtained from system.
# TYPE go_memstats_sys_bytes gauge
go_memstats_sys_bytes 5.6210696e+07
# HELP go_threads Number of OS threads created.
# TYPE go_threads gauge
go_threads 11
# HELP process_cpu_seconds_total Total user and system CPU time spent in seconds.
# TYPE process_cpu_seconds_total counter
process_cpu_seconds_total 1797.79
# HELP process_max_fds Maximum number of open file descriptors.
# TYPE process_max_fds gauge
process_max_fds 1.048576e+06
# HELP process_open_fds Number of open file descriptors.
# TYPE process_open_fds gauge
process_open_fds 20
# HELP process_resident_memory_bytes Resident memory size in bytes.
# TYPE process_resident_memory_bytes gauge
process_resident_memory_bytes 6.991872e+07
# HELP process_start_time_seconds Start time of the process since unix epoch in seconds.
# TYPE process_start_time_seconds gauge
process_start_time_seconds 1.68328404099e+09
# HELP process_virtual_memory_bytes Virtual memory size in bytes.
# TYPE process_virtual_memory_bytes gauge
process_virtual_memory_bytes 8.30472192e+08
# HELP process_virtual_memory_max_bytes Maximum amount of virtual memory available in bytes.
# TYPE process_virtual_memory_max_bytes gauge
process_virtual_memory_max_bytes -1
```
