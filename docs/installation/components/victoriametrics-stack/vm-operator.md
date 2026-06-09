#### vmoperator

<!-- markdownlint-disable line-length -->
| Field                            | Description                                                                                                                                                                                                            | Scheme                                                                                                                       |
| -------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| install                          | Allow to enable or disable deploy vmOperator via monitoring-operator.                                                                                                                                                  | boolean                                                                                                                      |
| paused                           | Set paused to reconciliation for vmOperator                                                                                                                                                                            | boolean                                                                                                                      |
| serviceMonitor                   | ServiceMonitor holds configuration attributes for vmoperator.                                                                                                                                                          | object                                                                                                                       |
| serviceMonitor.install           | Allow to enable or disable deploy vmOperator service monitor                                                                                                                                                           | boolean                                                                                                                      |
| serviceMonitor.interval          | Allow to change metrics scrape interval.                                                                                                                                                                               | string                                                                                                                       |
| serviceMonitor.scrapeTimeout     | Allow to change metrics scrape timeout. Note that scrapeTimeout must be less the interval                                                                                                                              | string                                                                                                                       |
| serviceMonitor.metricRelabelings | Set metricRelabelings for the ServiceMonitor, use to apply to samples for ingestion                                                                                                                                    | []*promv1.RelabelConfig                                                                                                      |
| serviceMonitor.relabelings       | Set relabelings for the ServiceMonitor, use to apply to samples before scraping                                                                                                                                        | []*promv1.RelabelConfig                                                                                                      |
| extraEnvs                        | Allow to set extra system environment variables for vmoperator, refer to [https://docs.victoriametrics.com/operator/vars/](https://docs.victoriametrics.com/operator/vars/)                                    | map[string]string                                                                                                            |
| image                            | A Docker image to deploy the vmoperator.                                                                                                                                                                               | string                                                                                                                       |
| containerSecurityContext         | Holds container-level security attributes. Default for Kubernetes, `containerSecurityContext: { capabilities.drop: ALL, allowPrivilegeEscalation: false }`.                                                            | [core/v1.SecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#securitycontext-v1-core)      |
| securityContext                  | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                     | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#podsecuritycontext-v1-core)    |
| resources                        | The resources that describe the compute resource requests and limits for single pods.                                                                                                                                  | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcerequirements-v1-core) |
| nodeSelector                     | Define which Nodes the Pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\"                                                                                                      | map[string]string                                                                                                            |
| tolerations                      | Allow the pods to schedule onto nodes with matching taints.                                                                                                                                                            | []v1.Toleration                                                                                                              |
| affinity                                            | If specified, the pod's scheduling constraints                                                                                                                                                                                      | *v1.Affinity                                                                                                                   |
| annotations                      | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value" | map[string]string                                                                                                            |
| labels                           | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                       | map[string]string                                                                                                            |
| priorityClassName                | PriorityClassName assigned to the Pods to prevent them from evicting.                                                                                                                                                  | string                                                                                                                       |
| tlsConfig                        | TLS configuration for VMOperator. Must be specified if `victoriametrics.tlsEnabled` is set to `true`                                                                                                                   | [TLSConfig](#tls-config)                                                                                                     |
<!-- markdownlint-enable line-length -->

Example:

```yaml
victoriametrics:
  tlsEnabled: false
  vmOperator:
    install: true
    paused: false
    serviceMonitor:
      install: true
      interval: 30s
      scrapeTimeout: 30s
      metricRelabelings: []
      relabelings: []
    extraEnvs: {}
    securityContext:
      runAsUser: 2000
      fsGroup: 2000
    resources:
      limits:
        cpu: 400m
        memory: 200Mi
      requests:
        cpu: 200m
        memory: 100Mi
    priorityClassName: priority-class
```

If TLS is enabled for Victoriametrics:

```yaml
victoriametrics:
  tlsEnabled: true
  vmOperator:
    install: true
    paused: false
    serviceMonitor:
      install: true
      interval: 30s
      scrapeTimeout: 30s
      metricRelabelings: []
      relabelings: []
    extraEnvs: {}
    securityContext:
      runAsUser: 2000
      fsGroup: 2000
    resources:
      limits:
        cpu: 400m
        memory: 200Mi
      requests:
        cpu: 200m
        memory: 100Mi
    priorityClassName: priority-class
    tlsConfig:
      generateCerts:
        enabled: true
        duration: 365
        renewBefore: 15
        secretName: "vmoperator-tls-secret"
```

