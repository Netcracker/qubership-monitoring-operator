### network-latency-exporter

**Important**: In some cases, TCP and UDP probes can show large packet loss not due to network problems, but due to system
settings. The `net.ipv4.icmp_ratelimit` kernel parameter with a default value of `1000` prevents a node from sending "
ICMP time in-transit exceeded" reply immediately, so `mtr` utility which is used in `network-latency-exporter` shows
losses.

Address **Important**: Please note that memory consumption by network-latency-exporter depends on number of nodes where
exporter will be installed and depends on number of using protocols.

The main formula is:

```bash
10Mi + (5Mi * (<number_of_nodes_in_cluster> - 1) * <number_of_protocols>)
```

_To avoid this behavior, set kernel variable `net.ipv4.icmp_ratelimit = 0` on all nodes._

<!-- markdownlint-disable line-length -->
| Field                         | Description                                                                                                                                                                                                               | Scheme                                                                                                                       |
| ----------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| install                       | Allows to disable create network-latency-exporter during the deployment.                                                                                                                                                  | bool                                                                                                                         |
| name                          | A name of the microservice to deploy with. This name is used as the name of the microservice deployment and in labels.                                                                                                    | string                                                                                                                       |
| image                         | A docker image to be used for the network-latency-exporter deployment.                                                                                                                                                    | string                                                                                                                       |
| rbac.createClusterRole        | Allow creating ClusterRole. If set to `false`, ClusterRole must be created manually.                                                                                                                                      | bool                                                                                                                         |
| rbac.createClusterRoleBinding | Allow creating ClusterRoleBinding. If set to `false`, ClusterRoleBinding must be created manually.                                                                                                                        | bool                                                                                                                         |
| rbac.setupSecurityContext     | Allow creating PodSecurityPolicy or SecurityContextConstraints. If set to `false`, PodSecurityPolicy / SecurityContextConstraints must be created manually.                                                               | bool                                                                                                                         |
| rbac.privileged               | If `true`, set parameters in PSP or SCC for ability to running container in the privileged mode and set `privileged: true` to the security context in the exporter's container.                                           | bool                                                                                                                         |
| createGrafanaDashboards       | Allow creating Grafana Dashboards `Network Latency Overview` and `Network Latency Details`.                                                                                                                               | bool                                                                                                                         |
| serviceAccount.create         | Allow creating ServiceAccount. If set to `false`, ServiceAccount must be created manually.                                                                                                                                | bool                                                                                                                         |
| serviceAccount.name           | Provide a name in place of network-latency-exporter for ServiceAccount.                                                                                                                                                   | string                                                                                                                       |
| extraArgs                     | The level of application logging may be set using extraArgs. Set `INFO` or `DEBUG` to see informative messages about metrics collection.                                                                                  | object                                                                                                                       |
| discoverEnable                | Allow enabling/disabling script for discovering nodes IP.                                                                                                                                                                 | bool                                                                                                                         |
| targets                       | List nodes with name and IP address. For example, [{"name": "node1", "ipAddress": "1.2.3.1"}, {"name": "node2", "ipAddress": "1.2.3.3"}]                                                                                  | object                                                                                                                       |
| requestTimeout                | The response time for each packet sent which the application waits response, in seconds.                                                                                                                                  | integer                                                                                                                      |
| timeout                       | The metrics collection timeout. Can be calculated as TIMEOUT = 10s + (REQUEST_TIMEOUT * PACKETS_NUM * <NUMBER_OF_PROTOCOLS>)                                                                                              | string                                                                                                                       |
| packetsNum                    | The number of packets to send per probe.                                                                                                                                                                                  | integer                                                                                                                      |
| packetSize                    | The size of packet to sent in bytes.                                                                                                                                                                                      | integer                                                                                                                      |
| checkTarget                   | The comma-separated list of network protocols and ports (separated by ':') through which packets are sent. The supported protocols are: `UDP`, `TCP`, `ICMP`. If no port is specified for the protocol, port `1` is used. | string                                                                                                                       |
| serviceMonitor.enabled        | If true, a ServiceMonitor is created for a `prometheus-operator`.                                                                                                                                                         | string                                                                                                                       |
| serviceMonitor.interval       | Scraping interval for Prometheus.                                                                                                                                                                                         | bool                                                                                                                         |
| additionalLabels              | Allows specifying custom labels for DaemonSet of `network-latency-exporter`.                                                                                                                                              | object                                                                                                                       |
| setupSecurityContext          | Allows to create PodSecurityPolicy or SecurityContextConstraints.                                                                                                                                                         | string                                                                                                                       |
| port                          | The port for node-exporter daemonset and service.                                                                                                                                                                         | int                                                                                                                          |
| resources                     | The resources that describe the compute resource requests and limits for single pods.                                                                                                                                     | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcerequirements-v1-core) |
| securityContext               | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 0, fsGroup: 2000 }`.                                                                                           | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#podsecuritycontext-v1-core)    |
| tolerations                   | Tolerations allow the pods to schedule onto nodes with matching taints.                                                                                                                                                   | []v1.Toleration                                                                                                              |
| nodeSelector                  | Defines which nodes the pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\"                                                                                                        | map[string]string                                                                                                            |
| affinity                      | If specified, the pod's scheduling constraints                                                                                                                                                                            | *v1.Affinity                                                                                                                 |
| annotations                   | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value"    | map[string]string                                                                                                            |
| labels                        | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                          | map[string]string                                                                                                            |
| priorityClassName             | PriorityClassName assigned to the Pods to prevent them from evicting.                                                                                                                                                     | string                                                                                                                       |
<!-- markdownlint-enable line-length -->

Example:

```yaml
networkLatencyExporter:
  install: true
  name: "network-latency-exporter"
  rbac:
    createClusterRole: true
    createClusterRoleBinding: true
    setupSecurityContext: true
    privileged: false
  createGrafanaDashboards: true
  serviceAccount:
    create: true
    name: "network-latency-exporter"
  image: ghcr.io/netcracker/qubership-network-latency-exporter:latest
  resources:
    limits:
      cpu: 200m
      memory: 256Mi
    requests:
      cpu: 100m
      memory: 128Mi
  securityContext:
    runAsUser: "0"
  tolerations: []
  nodeSelector: {}
  affinity: {}
  labels:
    label.key: label-value
  annotations:
    annotation.key: annotation-value
  priorityClassName: priority-class
  extraArgs:
    - "--log.level=debug"
  requestTimeout: 5
  timeout: 60s
  packetsNum: 10
  packetSize: 64
  checkTarget: "UDP:80,TCP:80,ICMP"
  serviceMonitor:
    enabled: true
    interval: 30s
```


