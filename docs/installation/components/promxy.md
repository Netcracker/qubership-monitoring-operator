### promxy

<!-- markdownlint-disable line-length -->
| Field                          | Description                                                                                                                                                                                                                                                  | Scheme                                                                                                                       |
| ------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ---------------------------------------------------------------------------------------------------------------------------- |
| install                        | Allows to enable or disable deploy promxy.                                                                                                                                                                                                                   | boolean                                                                                                                      |
| name                           | A name of the microservice to deploy with. This name is used as the name of the microservice deployment and in labels.                                                                                                                                       | string                                                                                                                       |
| image                          | A Docker image to deploy the promxy.                                                                                                                                                                                                                         | string                                                                                                                       |
| extraArgs                      | Additional arguments for promxy container.                                                                                                                                                                                                                   | list[string]                                                                                                                 |
| config                         | Configuration for promxy.                                                                                                                                                                                                                                    | object                                                                                                                       |
| config.serverGroups            | Base configuration of targets (server groups) for promxy. Every group should contain address (service or ingress for scraping Prometheus metrics), scheme (http or https) and unique label for identification in Grafana. Scheme "http" is using by default. | list[object]                                                                                                                 |
| config.serverGroups[N].address | URL (service or ingress) for scraping Prometheus metrics.                                                                                                                                                                                                    | string                                                                                                                       |
| config.serverGroups[N].scheme  | Configures the protocol scheme used for requests (http or https). Defaults to http.                                                                                                                                                                          | string                                                                                                                       |
| config.serverGroups[N].label   | Unique label to be added to metrics retrieved from this server group for identification in Grafana.                                                                                                                                                          | string                                                                                                                       |
| config.detailedConfig          | Detailed configuration for promxy in the YAML format. If this parameter is specified, the rest of parameters from the config section will not be used.                                                                                                       | [Promxy Configuration](https://github.com/jacksontj/promxy/blob/v0.0.71/cmd/promxy/config.yaml)                              |
| resources                      | The resources that describe the compute resource requests and limits for single pods.                                                                                                                                                                        | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcerequirements-v1-core) |
| port                           | Port for promxy container and service.                                                                                                                                                                                                                       | integer                                                                                                                      |
| configmapReload                | Additional container with small sidecar that reload promxy when configmap with configuration is changing.                                                                                                                                                    | object                                                                                                                       |
| configmapReload.install        | Allow to disable deploy promxy-configmap-reload container.                                                                                                                                                                                                   | boolean                                                                                                                      |
| configmapReload.image          | A docker image to use for promxy-configmap-reload deployment.                                                                                                                                                                                                | string                                                                                                                       |
| configmapReload.resources      | The resources describes the compute resource requests and limits for promxy-configmap-reload.                                                                                                                                                                | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcerequirements-v1-core) |
| securityContext                | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                                                           | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#podsecuritycontext-v1-core)    |
| tolerations                    | Tolerations allow the pods to schedule onto nodes with matching taints.                                                                                                                                                                                      | []v1.Toleration                                                                                                              |
| nodeSelector                   | Defines which nodes the pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\"                                                                                                                                           | map[string]string                                                                                                            |
| affinity                                            | If specified, the pod's scheduling constraints                                                                                                                                                                                      | *v1.Affinity                                                                                                                   |
| annotations                    | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value"                                       | map[string]string                                                                                                            |
| labels                         | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                                                             | map[string]string                                                                                                            |
| priorityClassName              | PriorityClassName assigned to the Pods to prevent them from evicting.                                                                                                                                                                                        | string                                                                                                                       |
<!-- markdownlint-enable line-length -->

Example:

```yaml
promxy:
  install: true
  name: "promxy"
  image: quay.io/jacksontj/promxy:v0.0.92
  extraArgs:
    - '--log-level=info'
  config:
    serverGroups:
      - address: "prometheus-operated:9090"
        label: "k8s-1"
      - address: "prometheus.k8s-2.cloud.com"
        scheme: https
        label: "k8s-2"
      - address: "victoriametrics.cloud.com"
        label: "k8s-itdpl"
    detailedConfig: ""
  resources:
    limits:
      cpu: 150m
      memory: 256Mi
    requests:
      cpu: 50m
      memory: 128Mi
  port: 9090
  configmapReload:
    install: true
    image: jimmidyson/configmap-reload:v0.5.0
    resources:
      limits:
        cpu: 10m
        memory: 20Mi
      requests:
        cpu: 5m
        memory: 3Mi
  labels:
    label.key: label-value
  annotations:
    annotation.key: annotation-value
  priorityClassName: priority-class
```


