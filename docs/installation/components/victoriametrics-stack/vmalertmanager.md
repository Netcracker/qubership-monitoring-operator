#### vmalertmanager

<!-- markdownlint-disable line-length -->
| Field                         | Description                                                                                                                                                                                                                                                                                        | Scheme                                                                                                                       |
| ----------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| install                       | Allows to enable or disable deploy vmalertmanager.                                                                                                                                                                                                                                                 | boolean                                                                                                                      |
| image                         | A Docker image to deploy the vmalertmanager.                                                                                                                                                                                                                                                       | string                                                                                                                       |
| ingress                       | Ingress allows to create Ingress for the vmalertmanager UI.                                                                                                                                                                                                                                        | *[Ingress](#ingress)                                                                                                         |
| secrets                       | Secrets is a list of Secrets in the same namespace as the VMSingle object, which shall be mounted into the vmalertmanager pods                                                                                                                                                                     | []string                                                                                                                     |
| configRawYaml                 | Raw configuration for vmalertmanager, it helps it to start without secret.                                                                                                                                                                                                                         | string                                                                                                                       |
| configSecret                  | The name of a Kubernetes Secret in the same namespace as the VMAlertmanager object, which contains configuration for this VMAlertmanager, configuration must be inside secret key: alertmanager.yaml.                                                                                              | string                                                                                                                       |
| retention                     | Time duration VMAlertmanager shall retain data for. Default is '120h'                                                                                                                                                                                                                              | string                                                                                                                       |
| paused                        | Set paused to reconciliation for vmalertmanager                                                                                                                                                                                                                                                    | boolean                                                                                                                      |
| nodeSelector                  | Defines which nodes the pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\"                                                                                                                                                                                 | map[string]string                                                                                                            |
| resources                     | The resources that describe the compute resource requests and limits for single pods.                                                                                                                                                                                                              | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcerequirements-v1-core) |
| affinity                      | It specifies the pod's scheduling constraints. For more information, refer to [https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core)                                       | *v1.Affinity                                                                                                                 |
| tolerations                   | Tolerations allow the pods to schedule onto nodes with matching taints.                                                                                                                                                                                                                            | []v1.Toleration                                                                                                              |
| securityContext               | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                                                                                                 | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#podsecuritycontext-v1-core)    |
| containers                    | Allows injecting additional containers or patching existing containers. This is meant to allow adding an authentication proxy to an VMAlertmanager pod.                                                                                                                                            | [[]v1.Container](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#container-v1-core)                     |
| selectAllByDefault            | Changes default behavior for empty CRD selectors, such ConfigSelector.                                                                                                                                                                                                                             | boolean                                                                                                                      |
| configSelector                | Defines selector for VMAlertmanagerConfig, result config will be merged with with Raw or Secret config.                                                                                                                                                                                            | *metav1.LabelSelector                                                                                                        |
| configNamespaceSelector       | Defines namespace selector for VMAlertmanagerConfig.                                                                                                                                                                                                                                               | *metav1.LabelSelector                                                                                                        |
| terminationGracePeriodSeconds | Defines period for container graceful termination.                                                                                                                                                                                                                                                 | string                                                                                                                       |
| labels                        | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                                                                                                   | map[string]string                                                                                                            |
| annotations                   | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value"                                                                             | map[string]string                                                                                                            |
| storage                       | Definition of how storage will be used by the VMAlertmanager                                                                                                                                                                                                                                       | *v1beta1.StorageSpec                                                                                                         |
| extraArgs                     | List of key-value with additional arguments which will be passed to VMAlertmanager pod [https://docs.victoriametrics.com/#list-of-command-line-flags](https://docs.victoriametrics.com/#list-of-command-line-flags)                                                                                | map[string]string                                                                                                            |
| extraEnvs                     | Allow to set extra system environment variables for VMAlertmanager.                                                                                                                                                                                                                                | map[string]string                                                                                                            |
| volumes                       | Volumes allows configuration of additional volumes on the output deploy definition. Volumes specified will be appended to other volumes that are generated as a result of StorageSpec objects                                                                                                      | []v1.Volume                                                                                                                  |
| volumeMounts                  | VolumeMounts allows configuration of additional VolumeMounts on the output deploy definition. VolumeMounts specified will be appended to other VolumeMounts in the vmsingle container that are generated as a result of StorageSpec objects                                                        | []v1.VolumeMount                                                                                                             |
| priorityClassName             | PriorityClassName assigned to the Pods to prevent them from evicting                                                                                                                                                                                                                               | string                                                                                                                       |
| tlsConfig                     | TLS configuration for VMAlertManager. Must be specified if `victoriametrics.tlsEnabled` is set to `true`                                                                                                                                                                                           | [TLSConfig](#tls-config)                                                                                                     |
| webConfig                     | Web configuration for VMAlertManager. Parameter is optional and this configuration is auto-populated by the operator. However, it can be used to customize/override TLS settings. More details [here](https://github.com/prometheus/alertmanager/blob/main/docs/https.md#https-and-authentication) | object                                                                                                                       |
| gossipConfig                  | Gossip configuration for VMAlertManager. Can be used to specify whether to use mutual TLS for gossip. More details [here](https://github.com/prometheus/alertmanager/blob/main/docs/https.md#gossip-traffic)                                                                                       | object                                                                                                                       |
<!-- markdownlint-enable line-length -->

Example:

```yaml
  vmAlertManager:
    install: true
    ingress:
      host: vmalertmanager.test.org
      install: true
    paused: false
    selectAllByDefault: true
    securityContext:
      runAsUser: 2001
      fsGroup: 2001
    secrets:
      - kube-etcd-client-certs
    volumes: {}
    volumeMounts: {}
```

If TLS is enabled for Victoriametrics:

```yaml
victoriametrics:
  tlsEnabled: true
  vmAlertManager:
    install: true
    ingress:
      host: vmalertmanager.test.org
      install: true
    paused: false
    selectAllByDefault: true
    securityContext:
      runAsUser: 2001
      fsGroup: 2001
    secrets:
      - kube-etcd-client-certs
    volumes: {}
    volumeMounts: {}
    tlsConfig:
      existingSecret: "vmalertmanager-tls-secret"
    webConfig: # optional
      tls_server_config:
        cert_secret_ref:
          key: tls.crt
          name: vmalertmanager-tls-secret
        key_secret_ref:
          key: tls.key
          name: vmalertmanager-tls-secret
```

**Note**
In above case, it is assumed that there is already a secret with the name `vmalertmanager-tls-secret`
present in the server with valid certificates.

