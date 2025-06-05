#### vmsingle

<!-- markdownlint-disable line-length -->
| Field                      | Description                                                                                                                                                                                                                                 | Scheme                                                                                                                       |
| -------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| install                    | Allow to enable or disable deploy vmSingle via monitoring-operator.                                                                                                                                                                         | boolean                                                                                                                      |
| image                      | A Docker image to deploy the vmoperator.                                                                                                                                                                                                    | string                                                                                                                       |
| resources                  | The resources that describe the compute resource requests and limits for single pods.                                                                                                                                                       | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcerequirements-v1-core) |
| securityContext            | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000, runAsGroup: 2000 }`.                                                                                        | [*v1.PodSecurityContext](https://kubernetes.io/docs/tasks/configure-pod-container/security-context/)                         |
| containers                 | Containers property allows to inject additions sidecars or to patch existing containers. It can be useful for proxies, backup, etc.                                                                                                         | [[]v1.Container](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#container-v1-core)                     |
| securityContext.runAsUser  | Specifies that for any Containers in the Pod, all processes run with user ID.                                                                                                                                                               | [*v1.PodSecurityContext](https://kubernetes.io/docs/tasks/configure-pod-container/security-context/)                         |
| securityContext.runAsGroup | Specifies the primary group ID for all processes within any containers of the Pod.                                                                                                                                                          | [*v1.PodSecurityContext](https://kubernetes.io/docs/tasks/configure-pod-container/security-context/)                         |
| securityContext.fsGroup    | Specifies that all processes of the container are also part of the supplementary group ID. The owner for volume and any files created in that volume will be Group ID.                                                                      | [*v1.PodSecurityContext](https://kubernetes.io/docs/tasks/configure-pod-container/security-context/)                         |
| paused                     | Set paused to reconciliation for vmsingle                                                                                                                                                                                                   | boolean                                                                                                                      |
| tolerations                | Allow the pods to schedule onto nodes with matching taints.                                                                                                                                                                                 | []v1.Toleration                                                                                                              |
| nodeSelector               | Define which Nodes the Pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\"                                                                                                                           | map[string]string                                                                                                            |
| affinity                                            | If specified, the pod's scheduling constraints                                                                                                                                                                                      | *v1.Affinity                                                                                                                   |
| labels                     | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                                            | map[string]string                                                                                                            |
| annotations                | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value"                      | map[string]string                                                                                                            |
| ingress                    | Ingress allows to create Ingress.                                                                                                                                                                                                           | *[Ingress](#ingress)                                                                                                         |
| retentionPeriod            | RetentionPeriod for the stored metrics. [https://docs.victoriametrics.com/Single-server-VictoriaMetrics.html#retention](https://docs.victoriametrics.com/Single-server-VictoriaMetrics.html#retention)                                      | string                                                                                                                       |
| extraArgs                  | List of key-value with additional arguments for run vmsingle [https://docs.victoriametrics.com/#list-of-command-line-flags](https://docs.victoriametrics.com/#list-of-command-line-flags)                                                   | map[string]string                                                                                                            |
| extraEnvs                  | Allow to set extra system environment variables for vmsingle.                                                                                                                                                                               | map[string]string                                                                                                            |
| secrets                    | Secrets is a list of Secrets in the same namespace as the VMSingle object, which shall be mounted into the VMSingle Pods                                                                                                                    | []string                                                                                                                     |
| storageDataPath            | Disables spec.storage option and overrides arg for victoria-metrics binary --storageDataPath, its users responsibility to mount proper device into given path.                                                                              | string                                                                                                                       |
| storage                    | Storage is the definition of how storage will be used by the VMSingle by default it`s empty dir                                                                                                                                             | *v1.PersistentVolumeClaimSpec                                                                                                |
| storageMetadata            | Defines annotations and labels attached to PVC for given vmsingle CR                                                                                                                                                                        | *github.com/VictoriaMetrics/operator/api/operator/v1beta1.EmbeddedObjectMetadata                                             |
| volumes                    | Volumes allows configuration of additional volumes on the output deploy definition. Volumes specified will be appended to other volumes that are generated as a result of StorageSpec objects                                               | []v1.Volume                                                                                                                  |
| volumeMounts               | VolumeMounts allows configuration of additional VolumeMounts on the output deploy definition. VolumeMounts specified will be appended to other VolumeMounts in the vmsingle container that are generated as a result of StorageSpec objects | []v1.VolumeMount                                                                                                             |
| priorityClassName          | PriorityClassName assigned to the Pods to prevent them from evicting                                                                                                                                                                        | string                                                                                                                       |
| tlsConfig                  | TLS configuration for VMSingle. Must be specified if `victoriametrics.tlsEnabled` is set to `true`                                                                                                                                          | [TLSConfig](#tls-config)                                                                                                     |
<!-- markdownlint-enable line-length -->

Example:

```yaml
victoriametrics:
  vmSingle:
    install: true
    ingress:
      host: vmsingle.test.org
      install: true
    extraArgs:
      search.maxPointsPerTimeseries: "150000"
    paused: false
    retentionPeriod: 1d
    securityContext:
      runAsUser: 2000
      runAsGroup: 2000
      fsGroup: 2000
    resources:
      limits:
        cpu: 1000m
        memory: 2000Mi
      requests:
        cpu: 500m
        memory: 1000Mi
    storage:
      accessModes:
        - ReadWriteOnce
      resources:
        requests:
          storage: 1Gi
      selector:
        matchLabels:
          app.kubernetes.io/name: vmsingle
      storageClassName: openstack-cinder
```

If TLS is enabled for Victoriametrics:

```yaml
victoriametrics:
  tlsEnabled: true
  clusterIssuerName: dev-cluster-issuer
  vmSingle:
    install: true
    ingress:
      host: vmsingle.test.org
      install: true
    extraArgs:
      search.maxPointsPerTimeseries: "150000"
    paused: false
    retentionPeriod: 1d
    securityContext:
      runAsUser: 2000
      runAsGroup: 2000
      fsGroup: 2000
    resources:
      limits:
        cpu: 1000m
        memory: 2000Mi
      requests:
        cpu: 500m
        memory: 1000Mi
    storage:
      accessModes:
        - ReadWriteOnce
      resources:
        requests:
          storage: 1Gi
      selector:
        matchLabels:
          app.kubernetes.io/name: vmsingle
      storageClassName: openstack-cinder
    tlsConfig:
      generateCerts:
        enabled: true
        duration: 365
        renewBefore: 15
        secretName: "vmsingle-tls-secret"
```

