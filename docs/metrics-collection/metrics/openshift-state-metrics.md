This document describes the metrics list and how to collect them from Openshift state metrics.

# Metrics

Consul already can exposes its metrics in Prometheus format and doesn't require to use of specific exporters.

| Name       | Metrics Port | Metrics Endpoint    | Need Exporter? | Auth?          | Is Exporter Third Party? |
| ---------- | ------------ | ------------------- | -------------- | -------------- | ------------------------ |
| Prometheus | `8443`       | `/metrics`          | No             | Require, token | N/A                      |

## How to Collect

Metrics expose on port `8443` and endpoint `/metrics`. By default, Openshift state metrics has authentication by token.

Config `ServiceMonitor` for `prometheus-operator` to collect Openshift state metrics:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: openshift-state-metrics
  labels:
    k8s-app: openshift-state-metrics
    app.kubernetes.io/name: openshift-apiserver-operator-check-endpoints
    app.kubernetes.io/component: monitoring
    app.kubernetes.io/managed-by: monitoring-operator
spec:
  endpoints:
    - bearerTokenFile: /var/run/secrets/kubernetes.io/serviceaccount/token
      honorLabels: true
      interval: 30s
      port: https-main
      scheme: https
      scrapeTimeout: 2m
      tlsConfig:
        insecureSkipVerify: true
    - bearerTokenFile: /var/run/secrets/kubernetes.io/serviceaccount/token
      interval: 30s
      port: https-self
      scheme: https
      scrapeTimeout: 2m
      tlsConfig:
        insecureSkipVerify: true
  namespaceSelector:
    matchNames:
      - openshift-monitoring
  selector:
    matchLabels:
      k8s-app: openshift-state-metrics
```

To collect (or just to check) metrics manually you can use the following command:

```bash
curl -v -k -L -H "Authorization: Bearer <token>" "http://<api_server_ip_or_dns>:8443/metrics"
```

Token usually you can find in the prometheus/vmagent pod by path `/var/run/secrets/kubernetes.io/serviceaccount/token`.

You can't use `wget` because it doesn't allow to add headers for authorization.

## Metrics List

### Openshift state metrics

```prometheus
# HELP openshift_buildconfig_created Unix creation timestamp
# TYPE openshift_buildconfig_created gauge
# HELP openshift_buildconfig_metadata_generation Sequence number representing a specific generation of the desired state.
# TYPE openshift_buildconfig_metadata_generation gauge
# HELP openshift_buildconfig_labels Kubernetes labels converted to Prometheus labels.
# TYPE openshift_buildconfig_labels gauge
# HELP openshift_buildconfig_status_latest_version The latest version of buildconfig.
# TYPE openshift_buildconfig_status_latest_version gauge
# HELP openshift_build_created_timestamp_seconds Unix creation timestamp
# TYPE openshift_build_created_timestamp_seconds gauge
# HELP openshift_build_metadata_generation_info Sequence number representing a specific generation of the desired state.
# TYPE openshift_build_metadata_generation_info gauge
# HELP openshift_build_labels Kubernetes labels converted to Prometheus labels.
# TYPE openshift_build_labels gauge
# HELP openshift_build_status_phase_total The build phase.
# TYPE openshift_build_status_phase_total gauge
# HELP openshift_build_start_timestamp_seconds Start time of the build
# TYPE openshift_build_start_timestamp_seconds gauge
# HELP openshift_build_completed_timestamp_seconds Completion time of the build
# TYPE openshift_build_completed_timestamp_seconds gauge
# HELP openshift_build_duration_seconds Duration of the build
# TYPE openshift_build_duration_seconds gauge
# HELP openshift_clusterresourcequota_created Unix creation timestamp
# TYPE openshift_clusterresourcequota_created gauge
# HELP openshift_clusterresourcequota_labels Kubernetes labels converted to Prometheus labels.
# TYPE openshift_clusterresourcequota_labels gauge
# HELP openshift_clusterresourcequota_usage Usage about resource quota.
# TYPE openshift_clusterresourcequota_usage gauge
# HELP openshift_clusterresourcequota_namespace_usage Usage about clusterresource quota per namespace.
# TYPE openshift_clusterresourcequota_namespace_usage gauge
# HELP openshift_clusterresourcequota_selector Selector of clusterresource quota, which defines the affected namespaces.
# TYPE openshift_clusterresourcequota_selector gauge
# HELP openshift_deploymentconfig_created Unix creation timestamp
# TYPE openshift_deploymentconfig_created gauge
# HELP openshift_deploymentconfig_status_replicas The number of replicas per deployment.
# TYPE openshift_deploymentconfig_status_replicas gauge
# HELP openshift_deploymentconfig_status_replicas_available The number of available replicas per deployment.
# TYPE openshift_deploymentconfig_status_replicas_available gauge
# HELP openshift_deploymentconfig_status_replicas_unavailable The number of unavailable replicas per deployment.
# TYPE openshift_deploymentconfig_status_replicas_unavailable gauge
# HELP openshift_deploymentconfig_status_replicas_updated The number of updated replicas per deployment.
# TYPE openshift_deploymentconfig_status_replicas_updated gauge
# HELP openshift_deploymentconfig_status_observed_generation The generation observed by the deployment controller.
# TYPE openshift_deploymentconfig_status_observed_generation gauge
# HELP openshift_deploymentconfig_spec_replicas Number of desired pods for a deployment.
# TYPE openshift_deploymentconfig_spec_replicas gauge
# HELP openshift_deploymentconfig_spec_paused Whether the deployment is paused and will not be processed by the deployment controller.
# TYPE openshift_deploymentconfig_spec_paused gauge
# HELP openshift_deploymentconfig_spec_strategy_rollingupdate_max_unavailable Maximum number of unavailable replicas during a rolling update of a deployment.
# TYPE openshift_deploymentconfig_spec_strategy_rollingupdate_max_unavailable gauge
# HELP openshift_deploymentconfig_spec_strategy_rollingupdate_max_surge Maximum number of replicas that can be scheduled above the desired number of replicas during a rolling update of a deployment.
# TYPE openshift_deploymentconfig_spec_strategy_rollingupdate_max_surge gauge
# HELP openshift_deploymentconfig_metadata_generation Sequence number representing a specific generation of the desired state.
# TYPE openshift_deploymentconfig_metadata_generation gauge
# HELP openshift_deploymentconfig_labels Kubernetes labels converted to Prometheus labels.
# TYPE openshift_deploymentconfig_labels gauge
# HELP openshift_group_created Unix creation timestamp
# TYPE openshift_group_created gauge
# HELP openshift_group_user_account User account in a group.
# TYPE openshift_group_user_account gauge
# HELP openshift_route_created Unix creation timestamp
# TYPE openshift_route_created gauge
openshift_route_created{namespace="openshift-monitoring",route="prometheus-k8s"} 1.689348583e+09
openshift_route_created{namespace="opensearch-service",route="opensearch-9hz5d"} 1.69875453e+09
openshift_route_created{namespace="opensearch-service",route="opensearch-dashboards-s7jwz"} 1.691580768e+09
openshift_route_created{namespace="ops-profiler",route="esc-collector-service"} 1.690880862e+09
openshift_route_created{namespace="openshift-console",route="downloads"} 1.689348757e+09
openshift_route_created{namespace="openshift-monitoring",route="alertmanager-main"} 1.689348584e+09
openshift_route_created{namespace="ops-monitoring",route="ops-monitoring-alertmanager"} 1.690876138e+09
openshift_route_created{namespace="openshift-authentication",route="oauth-openshift"} 1.689348573e+09
openshift_route_created{namespace="profiler-qa",route="esc-test-service-25bnf"} 1.699524266e+09
openshift_route_created{namespace="openshift-monitoring",route="prometheus-k8s-federate"} 1.689348584e+09
openshift_route_created{namespace="openshift-console",route="console"} 1.689348756e+09
openshift_route_created{namespace="openshift-monitoring",route="thanos-querier"} 1.689348584e+09
openshift_route_created{namespace="profiler-qa",route="esc-ui-service-xngbb"} 1.699524242e+09
openshift_route_created{namespace="profiler-qa",route="esc-collector-service-x8vqp"} 1.699524253e+09
openshift_route_created{namespace="ops-profiler",route="cdt-ui-service"} 1.693195673e+09
openshift_route_created{namespace="ops-profiler",route="esc-static-service"} 1.690880859e+09
openshift_route_created{namespace="ops-monitoring",route="ops-monitoring-grafana"} 1.690879833e+09
openshift_route_created{namespace="ops-monitoring",route="ops-monitoring-vmauth"} 1.690876094e+09
openshift_route_created{namespace="ops-profiler",route="esc-test-service"} 1.690880865e+09
openshift_route_created{namespace="profiler-qa",route="esc-static-service-mr7qm"} 1.699524248e+09
openshift_route_created{namespace="openshift-ingress-canary",route="canary"} 1.68934855e+09
openshift_route_created{namespace="ops-profiler",route="esc-ui-service"} 1.690880852e+09
# HELP openshift_route_info Information about route.
# TYPE openshift_route_info gauge
openshift_route_info{namespace="openshift-monitoring",route="prometheus-k8s",host="prometheus-k8s-openshift-monitoring.k8s.test.org",path="/api",tls_termination="reencrypt",to_kind="Service",to_name="prometheus-k8s",to_weight="100"} 1
openshift_route_info{namespace="opensearch-service",route="opensearch-9hz5d",host="opensearch-elasticsearch-cluster.k8s.test.org",path="/",tls_termination="",to_kind="Service",to_name="opensearch-internal",to_weight="100"} 1
openshift_route_info{namespace="opensearch-service",route="opensearch-dashboards-s7jwz",host="dashboards-opensearch-service.k8s.test.org",path="/",tls_termination="",to_kind="Service",to_name="opensearch-dashboards",to_weight="100"} 1
openshift_route_info{namespace="ops-profiler",route="esc-collector-service",host="esc-collector-service-ops-profiler.k8s.test.org",path="",tls_termination="",to_kind="Service",to_name="esc-collector-service",to_weight="100"} 1
openshift_route_info{namespace="openshift-console",route="downloads",host="downloads-openshift-console.k8s.test.org",path="",tls_termination="edge",to_kind="Service",to_name="downloads",to_weight="100"} 1
openshift_route_info{namespace="openshift-monitoring",route="alertmanager-main",host="alertmanager-main-openshift-monitoring.k8s.test.org",path="/api",tls_termination="reencrypt",to_kind="Service",to_name="alertmanager-main",to_weight="100"} 1
openshift_route_info{namespace="ops-monitoring",route="ops-monitoring-alertmanager",host="alertmanager-ops-monitoring.k8s.test.org",path="",tls_termination="",to_kind="Service",to_name="alertmanager-operated",to_weight="100"} 1
openshift_route_info{namespace="openshift-authentication",route="oauth-openshift",host="oauth-openshift.k8s.test.org",path="",tls_termination="passthrough",to_kind="Service",to_name="oauth-openshift",to_weight="100"} 1
openshift_route_info{namespace="profiler-qa",route="esc-test-service-25bnf",host="esc-test-service-profiler-qa.k8s.test.org",path="/",tls_termination="",to_kind="Service",to_name="esc-test-service",to_weight="100"} 1
openshift_route_info{namespace="openshift-monitoring",route="prometheus-k8s-federate",host="prometheus-k8s-federate-openshift-monitoring.k8s.test.org",path="/federate",tls_termination="reencrypt",to_kind="Service",to_name="prometheus-k8s",to_weight="100"} 1
openshift_route_info{namespace="openshift-console",route="console",host="console-openshift-console.k8s.test.org",path="",tls_termination="reencrypt",to_kind="Service",to_name="console",to_weight="100"} 1
openshift_route_info{namespace="openshift-monitoring",route="thanos-querier",host="thanos-querier-openshift-monitoring.k8s.test.org",path="/api",tls_termination="reencrypt",to_kind="Service",to_name="thanos-querier",to_weight="100"} 1
openshift_route_info{namespace="profiler-qa",route="esc-ui-service-xngbb",host="esc-ui-service-profiler-qa.k8s.test.org",path="/",tls_termination="",to_kind="Service",to_name="esc-ui-service",to_weight="100"} 1
openshift_route_info{namespace="profiler-qa",route="esc-collector-service-x8vqp",host="esc-collector-service-profiler-qa.k8s.test.org",path="/",tls_termination="",to_kind="Service",to_name="esc-collector-service",to_weight="100"} 1
openshift_route_info{namespace="ops-profiler",route="cdt-ui-service",host="cdt-ui-service-ops-profiler.k8s.test.org",path="",tls_termination="",to_kind="Service",to_name="cdt-collector-deployment",to_weight="100"} 1
openshift_route_info{namespace="ops-profiler",route="esc-static-service",host="esc-static-service-ops-profiler.k8s.test.org",path="",tls_termination="",to_kind="Service",to_name="esc-static-service",to_weight="100"} 1
openshift_route_info{namespace="ops-monitoring",route="ops-monitoring-grafana",host="grafana-ops-monitoring.k8s.test.org",path="",tls_termination="",to_kind="Service",to_name="grafana-service",to_weight="100"} 1
openshift_route_info{namespace="ops-monitoring",route="ops-monitoring-vmauth",host="vmauth-ops-monitoring.k8s.test.org",path="/",tls_termination="",to_kind="Service",to_name="vmauth-k8s",to_weight="100"} 1
openshift_route_info{namespace="ops-profiler",route="esc-test-service",host="esc-test-service-ops-profiler.k8s.test.org",path="",tls_termination="",to_kind="Service",to_name="esc-test-service",to_weight="100"} 1
openshift_route_info{namespace="profiler-qa",route="esc-static-service-mr7qm",host="esc-static-service-profiler-qa.k8s.test.org",path="/",tls_termination="",to_kind="Service",to_name="esc-static-service",to_weight="100"} 1
openshift_route_info{namespace="openshift-ingress-canary",route="canary",host="canary-openshift-ingress-canary.k8s.test.org",path="",tls_termination="edge",to_kind="Service",to_name="ingress-canary",to_weight="100"} 1
openshift_route_info{namespace="ops-profiler",route="esc-ui-service",host="esc-ui-service-ops-profiler.k8s.test.org",path="",tls_termination="",to_kind="Service",to_name="esc-ui-service",to_weight="100"} 1
# HELP openshift_route_status Information about route status.
# TYPE openshift_route_status gauge
openshift_route_status{namespace="openshift-authentication",route="oauth-openshift",status="True",type="Admitted",host="oauth-openshift.k8s.test.org",router_name="default"} 1
openshift_route_status{namespace="profiler-qa",route="esc-test-service-25bnf",status="True",type="Admitted",host="esc-test-service-profiler-qa.k8s.test.org",router_name="default"} 1
openshift_route_status{namespace="openshift-monitoring",route="prometheus-k8s-federate",status="True",type="Admitted",host="prometheus-k8s-federate-openshift-monitoring.k8s.test.org",router_name="default"} 1
openshift_route_status{namespace="openshift-console",route="console",status="True",type="Admitted",host="console-openshift-console.k8s.test.org",router_name="default"} 1
openshift_route_status{namespace="openshift-console",route="downloads",status="True",type="Admitted",host="downloads-openshift-console.k8s.test.org",router_name="default"} 1
openshift_route_status{namespace="openshift-monitoring",route="alertmanager-main",status="True",type="Admitted",host="alertmanager-main-openshift-monitoring.k8s.test.org",router_name="default"} 1
openshift_route_status{namespace="ops-monitoring",route="ops-monitoring-alertmanager",status="True",type="Admitted",host="alertmanager-ops-monitoring.k8s.test.org",router_name="default"} 1
openshift_route_status{namespace="ops-profiler",route="cdt-ui-service",status="True",type="Admitted",host="cdt-ui-service-ops-profiler.k8s.test.org",router_name="default"} 1
openshift_route_status{namespace="ops-profiler",route="esc-static-service",status="True",type="Admitted",host="esc-static-service-ops-profiler.k8s.test.org",router_name="default"} 1
openshift_route_status{namespace="openshift-monitoring",route="thanos-querier",status="True",type="Admitted",host="thanos-querier-openshift-monitoring.k8s.test.org",router_name="default"} 1
openshift_route_status{namespace="profiler-qa",route="esc-ui-service-xngbb",status="True",type="Admitted",host="esc-ui-service-profiler-qa.k8s.test.org",router_name="default"} 1
openshift_route_status{namespace="profiler-qa",route="esc-collector-service-x8vqp",status="True",type="Admitted",host="esc-collector-service-profiler-qa.k8s.test.org",router_name="default"} 1
openshift_route_status{namespace="profiler-qa",route="esc-static-service-mr7qm",status="True",type="Admitted",host="esc-static-service-profiler-qa.k8s.test.org",router_name="default"} 1
openshift_route_status{namespace="openshift-ingress-canary",route="canary",status="True",type="Admitted",host="canary-openshift-ingress-canary.k8s.test.org",router_name="default"} 1
openshift_route_status{namespace="ops-profiler",route="esc-ui-service",status="True",type="Admitted",host="esc-ui-service-ops-profiler.k8s.test.org",router_name="default"} 1
openshift_route_status{namespace="ops-monitoring",route="ops-monitoring-grafana",status="True",type="Admitted",host="grafana-ops-monitoring.k8s.test.org",router_name="default"} 1
openshift_route_status{namespace="ops-monitoring",route="ops-monitoring-vmauth",status="True",type="Admitted",host="vmauth-ops-monitoring.k8s.test.org",router_name="default"} 1
openshift_route_status{namespace="ops-profiler",route="esc-test-service",status="True",type="Admitted",host="esc-test-service-ops-profiler.k8s.test.org",router_name="default"} 1
openshift_route_status{namespace="opensearch-service",route="opensearch-dashboards-s7jwz",status="True",type="Admitted",host="dashboards-opensearch-service.k8s.test.org",router_name="default"} 1
openshift_route_status{namespace="ops-profiler",route="esc-collector-service",status="True",type="Admitted",host="esc-collector-service-ops-profiler.k8s.test.org",router_name="default"} 1
openshift_route_status{namespace="openshift-monitoring",route="prometheus-k8s",status="True",type="Admitted",host="prometheus-k8s-openshift-monitoring.k8s.test.org",router_name="default"} 1
openshift_route_status{namespace="opensearch-service",route="opensearch-9hz5d",status="True",type="Admitted",host="opensearch-elasticsearch-cluster.k8s.test.org",router_name="default"} 1
# HELP openshift_route_labels Kubernetes labels converted to Prometheus labels.
# TYPE openshift_route_labels gauge
openshift_route_labels{namespace="ops-profiler",route="cdt-ui-service"} 1
openshift_route_labels{namespace="ops-profiler",route="esc-static-service",label_app_kubernetes_io_version="release-9.3.2.64-20230727.060643-2-RELEASE",label_app_kubernetes_io_component="backend",label_app_kubernetes_io_managed_by="Helm",label_app_kubernetes_io_name="esc-static-service",label_app_kubernetes_io_part_of="esc"} 1
openshift_route_labels{namespace="openshift-monitoring",route="thanos-querier",label_app_kubernetes_io_component="query-layer",label_app_kubernetes_io_instance="thanos-querier",label_app_kubernetes_io_name="thanos-query",label_app_kubernetes_io_part_of="openshift-monitoring",label_app_kubernetes_io_version="0.30.2"} 1
openshift_route_labels{namespace="profiler-qa",route="esc-ui-service-xngbb",label_app_kubernetes_io_part_of="esc",label_app_kubernetes_io_version="release-9.3.2.70-20231102.083351-3-RELEASE",label_app_kubernetes_io_component="backend",label_app_kubernetes_io_managed_by="Helm",label_app_kubernetes_io_name="esc-ui-service"} 1
openshift_route_labels{namespace="profiler-qa",route="esc-collector-service-x8vqp",label_app_kubernetes_io_component="backend",label_app_kubernetes_io_managed_by="Helm",label_app_kubernetes_io_name="esc-collector-service",label_app_kubernetes_io_part_of="esc",label_app_kubernetes_io_version="release-9.3.2.70-20231102.083351-3-RELEASE"} 1
openshift_route_labels{namespace="openshift-ingress-canary",route="canary",label_ingress_openshift_io_canary="canary_controller"} 1
openshift_route_labels{namespace="ops-profiler",route="esc-ui-service",label_app_kubernetes_io_component="backend",label_app_kubernetes_io_managed_by="Helm",label_app_kubernetes_io_name="esc-ui-service",label_app_kubernetes_io_part_of="esc",label_app_kubernetes_io_version="release-9.3.2.64-20230727.060643-2-RELEASE"} 1
openshift_route_labels{namespace="ops-monitoring",route="ops-monitoring-grafana",label_app_kubernetes_io_component="monitoring",label_app_kubernetes_io_managed_by="grafana-operator",label_app_kubernetes_io_name="grafana"} 1
openshift_route_labels{namespace="ops-monitoring",route="ops-monitoring-vmauth",label_app_kubernetes_io_component="monitoring",label_app_kubernetes_io_instance="k8s",label_app_kubernetes_io_managed_by="monitoring-operator",label_app_kubernetes_io_name="vmauth"} 1
openshift_route_labels{namespace="ops-profiler",route="esc-test-service",label_app_kubernetes_io_managed_by="Helm",label_app_kubernetes_io_name="esc-test-service",label_app_kubernetes_io_part_of="esc",label_app_kubernetes_io_version="release-9.3.2.64-20230727.060643-2-RELEASE",label_app_kubernetes_io_component="backend"} 1
openshift_route_labels{namespace="profiler-qa",route="esc-static-service-mr7qm",label_app_kubernetes_io_component="backend",label_app_kubernetes_io_managed_by="Helm",label_app_kubernetes_io_name="esc-static-service",label_app_kubernetes_io_part_of="esc",label_app_kubernetes_io_version="release-9.3.2.70-20231102.083351-3-RELEASE"} 1
openshift_route_labels{namespace="opensearch-service",route="opensearch-dashboards-s7jwz",label_app="opensearch",label_app_kubernetes_io_managed_by="Helm",label_chart="opensearch-0.1.0",label_heritage="Helm",label_release="opensearch-service"} 1
openshift_route_labels{namespace="ops-profiler",route="esc-collector-service",label_app_kubernetes_io_component="backend",label_app_kubernetes_io_managed_by="Helm",label_app_kubernetes_io_name="esc-collector-service",label_app_kubernetes_io_part_of="esc",label_app_kubernetes_io_version="release-9.3.2.64-20230727.060643-2-RELEASE"} 1
openshift_route_labels{namespace="openshift-monitoring",route="prometheus-k8s"} 1
openshift_route_labels{namespace="opensearch-service",route="opensearch-9hz5d",label_app_kubernetes_io_managed_by="Helm",label_chart="opensearch-0.1.0",label_heritage="Helm",label_release="opensearch-service",label_role="client",label_app="opensearch"} 1
openshift_route_labels{namespace="openshift-monitoring",route="prometheus-k8s-federate"} 1
openshift_route_labels{namespace="openshift-console",route="console",label_app="console"} 1
openshift_route_labels{namespace="openshift-console",route="downloads"} 1
openshift_route_labels{namespace="openshift-monitoring",route="alertmanager-main"} 1
openshift_route_labels{namespace="ops-monitoring",route="ops-monitoring-alertmanager",label_app_kubernetes_io_component="monitoring",label_app_kubernetes_io_managed_by="prometheus-operator",label_app_kubernetes_io_name="alertmanager"} 1
openshift_route_labels{namespace="openshift-authentication",route="oauth-openshift",label_app="oauth-openshift"} 1
openshift_route_labels{namespace="profiler-qa",route="esc-test-service-25bnf",label_app_kubernetes_io_component="backend",label_app_kubernetes_io_managed_by="Helm",label_app_kubernetes_io_name="esc-test-service",label_app_kubernetes_io_part_of="esc",label_app_kubernetes_io_version="release-9.3.2.70-20231102.083351-3-RELEASE"} 1
```
