#### vmauth

<!-- markdownlint-disable line-length -->
| Field                         | Description                                                                                                                                                                                                                                                  | Scheme                                                                                                                       | Required |
| ----------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ---------------------------------------------------------------------------------------------------------------------------- | -------- |
| install                       | Allows to enable or disable deploy vmauth.                                                                                                                                                                                                                   | boolean                                                                                                                      | false    |
| image                         | A Docker image to deploy the vmauth.                                                                                                                                                                                                                         | string                                                                                                                       | false    |
| paused                        | Set paused to reconciliation for vmauth.                                                                                                                                                                                                                     | boolean                                                                                                                      | false    |
| ingress                       | Allows Ingress configuration for the vmauth.                                                                                                                                                                                                                 | *v1beta1.EmbeddedIngress                                                                                                     | false    |
| secrets                       | Secrets is a list of Secrets in the same namespace as the VMAuth object, which shall be mounted into the vmauth pods                                                                                                                                         | []string                                                                                                                     | false    |
| configMaps                    | ConfigMaps is a list of ConfigMaps in the same namespace as the VMAuth object, which shall be mounted into the VMAuth Pods.                                                                                                                                  | []string                                                                                                                     | false    |
| volumes                       | Volumes allows configuration of additional volumes on the output deploy definition. Volumes specified will be appended to other volumes that are generated as a result of StorageSpec objects                                                                | []v1.Volume                                                                                                                  | false    |
| volumeMounts                  | VolumeMounts allows configuration of additional VolumeMounts on the output deploy definition. VolumeMounts specified will be appended to other VolumeMounts in the vmauth container that are generated as a result of StorageSpec objects                    | []v1.VolumeMount                                                                                                             | false    |
| resources                     | The resources that describe the compute resource requests and limits for single pods.                                                                                                                                                                        | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcerequirements-v1-core) | false    |
| affinity                      | It specifies the pod's scheduling constraints. For more information, refer to [https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core) | *v1.Affinity                                                                                                                 | false    |
| tolerations                   | Tolerations allow the pods to schedule onto nodes with matching taints.                                                                                                                                                                                      | []v1.Toleration                                                                                                              | false    |
| securityContext               | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                                                           | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#podsecuritycontext-v1-core)    | false    |
| containers                    | Containers property allows to inject additions sidecars or to patch existing containers. It can be useful for proxies, backup, etc.                                                                                                                          | [[]v1.Container](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#container-v1-core)                     | false    |
| port                          | Port for listen.                                                                                                                                                                                                                                             | string                                                                                                                       | false    |
| selectAllByDefault            | Changes default behavior for empty CRD selectors, such UserSelector. Default - true.                                                                                                                                                                         | boolean                                                                                                                      | false    |
| userSelector                  | Defines VMUser to be selected for config file generation                                                                                                                                                                                                     | *metav1.LabelSelector                                                                                                        | false    |
| userNamespaceSelector         | Defines namespace selector for VMAuth discovery.                                                                                                                                                                                                             | *metav1.LabelSelector                                                                                                        | false    |
| extraArgs                     | List of key-value with additional arguments which will be passed to VMAuth pod [https://docs.victoriametrics.com/#list-of-command-line-flags](https://docs.victoriametrics.com/#list-of-command-line-flags)                                                  | map[string]string                                                                                                            | false    |
| extraEnvs                     | Allow to set extra system environment variables for vmauth.                                                                                                                                                                                                  | map[string]string                                                                                                            | false    |
| labels                        | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                                                             | map[string]string                                                                                                            | false    |
| annotations                   | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value"                                       | map[string]string                                                                                                            | false    |
| nodeSelector                  | Defines which nodes the pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\"                                                                                                                                           | map[string]string                                                                                                            | false    |
| terminationGracePeriodSeconds | Defines period for container graceful termination.                                                                                                                                                                                                           | string                                                                                                                       | false    |
| priorityClassName             | PriorityClassName assigned to the Pods to prevent them from evicting                                                                                                                                                                                         | string                                                                                                                       | false    |
| tlsConfig                     | TLS configuration for VMAuth. Must be specified if `victoriametrics.tlsEnabled` is set to `true`                                                                                                                                                             | [TLSConfig](#tls-config)                                                                                                     | false    |
<!-- markdownlint-enable line-length -->

```yaml
  vmAuth:
    install: true
    ingress:
      host: vmauth.test.org
    paused: false
    secretName: vmauth-secret
    extraVarsSecret:
      pass: vmauth
    securityContext:
      runAsUser: 2000
      fsGroup: 2000
```

If TLS is enabled for Victoriametrics:

```yaml
victoriametrics:
  tlsEnabled: true
  clusterIssuerName: dev-cluster-issuer
  vmAuth:
    install: true
    ingress:
      host: vmauth.test.org
    paused: false
    secretName: vmauth-secret
    extraVarsSecret:
      pass: vmauth
    securityContext:
      runAsUser: 2000
      fsGroup: 2000
    tlsConfig:
      generateCerts:
        enabled: true
        duration: 365
        renewBefore: 15
        secretName: "vmauth-tls-secret"
```

