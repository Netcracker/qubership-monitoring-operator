### alertmanager

AlertManager is a specification of the desired behavior of the AlertManager cluster.

<!-- markdownlint-disable line-length -->
| Field             | Description                                                                                                                                                                                                            | Scheme                                                                                                                       |
| ----------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| install           | Allows to disable deploy AlertManager. If AlertManager was not deployed during the deployment using helm, it can be deployed using change custom resource PlatformMonitoring.                                          | bool                                                                                                                         |
| paused            | Set paused to reconciliation.                                                                                                                                                                                          | bool                                                                                                                         |
| image             | A docker image to use for AlertManager deployment.                                                                                                                                                                     | string                                                                                                                       |
| ingress           | Ingress allows to create Ingress for the AlertManager UI.                                                                                                                                                              | *[Ingress](#ingress)                                                                                                         |
| nodeSelector      | Defines which nodes the pods are scheduled on. Specified just as map[string]string. For example: "type: compute"                                                                                                       | map[string]string                                                                                                            |
| affinity                                            | If specified, the pod's scheduling constraints                                                                                                                                                                                      | *v1.Affinity                                                                                                                   |
| annotations       | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value" | map[string]string                                                                                                            |
| labels            | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                       | map[string]string                                                                                                            |
| port              | The port for AlertManager service.                                                                                                                                                                                     | int                                                                                                                          |
| replicas          | Set replicas.                                                                                                                                                                                                          | *int32                                                                                                                       |
| resources         | The resources that describe the compute resource requests and limits for single Pods.                                                                                                                                  | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcerequirements-v1-core) |
| securityContext   | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                     | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#podsecuritycontext-v1-core)    |
| podMonitor        | Pod monitor for self monitoring.                                                                                                                                                                                       | *[Monitor](#monitor)                                                                                                         |
| priorityClassName | PriorityClassName assigned to the Pods to prevent them from evicting.                                                                                                                                                  | string                                                                                                                       |
<!-- markdownlint-enable line-length -->

Example:

```yaml
alertManager:
  install: true
  paused: false
  image: prom/alertmanager:v0.19.0
  port: 30903
  resources:
    limits:
      cpu: 200m
      memory: 200Mi
    requests:
      cpu: 100m
      memory: 100Mi
  ingress:
    ...see example by link...
  nodeSelector:
    node-role.kubernetes.io/worker: worker
  labels:
    label.key: label-value
  annotations:
    annotation.key: annotation-value 
  replicas: 1
  securityContext:
    runAsUser: 2000
    fsGroup: 2000
  podMonitor:
    ...see example by link...
```


