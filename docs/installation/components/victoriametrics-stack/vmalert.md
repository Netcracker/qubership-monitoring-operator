#### vmalert

<!-- markdownlint-disable line-length -->
| Field                         | Description                                                                                                                                                                                                                                                  | Scheme                                                                                                                       |
| ----------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ---------------------------------------------------------------------------------------------------------------------------- |
| install                       | Allows to enable or disable deploy vmalert.                                                                                                                                                                                                                  | boolean                                                                                                                      |
| image                         | A Docker image to deploy the vmalert.                                                                                                                                                                                                                        | string                                                                                                                       |
| ingress                       | Ingress allows to create Ingress for the vmalert UI.                                                                                                                                                                                                         | *[Ingress](#ingress)                                                                                                         |
| secrets                       | Secrets is a list of Secrets in the same namespace as the VMAlert object, which shall be mounted into the vmalert pods                                                                                                                                       | []string                                                                                                                     |
| resources                     | The resources that describe the compute resource requests and limits for single pods.                                                                                                                                                                        | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcerequirements-v1-core) |
| affinity                      | It specifies the pod's scheduling constraints. For more information, refer to [https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core) | *v1.Affinity                                                                                                                 |
| tolerations                   | Tolerations allow the pods to schedule onto nodes with matching taints.                                                                                                                                                                                      | []v1.Toleration                                                                                                              |
| securityContext               | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                                                           | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#podsecuritycontext-v1-core)    |
| containers                    | Containers property allows to inject additions sidecars or to patch existing containers. It can be useful for proxies, backup, etc.                                                                                                                          | [[]v1.Container](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#container-v1-core)                     |
| evaluationInterval            | How often evalute rules by default.                                                                                                                                                                                                                          | string                                                                                                                       |
| selectAllByDefault            | Changes default behavior for empty CRD selectors, such RuleSelector.                                                                                                                                                                                         | boolean                                                                                                                      |
| ruleSelector                  | Defines which VMRules to mount for loading alerting                                                                                                                                                                                                          | *metav1.LabelSelector                                                                                                        |
| ruleNamespaceSelector         | Defines namespace selector for VMRules discovery.                                                                                                                                                                                                            | *metav1.LabelSelector                                                                                                        |
| port                          | Port for listen.                                                                                                                                                                                                                                             | string                                                                                                                       |
| remoteWrite                   | Optional URL to remote-write compatible storage to persist vmalert state and rule results to.                                                                                                                                                                | []v1beta1.VMAlertRemoteWriteSpec                                                                                             |
| remoteWrite.url               | URL of the endpoint to send samples to.                                                                                                                                                                                                                      | string                                                                                                                       |
| remoteWrite.flushInterval     | Defines interval of flushes to remote write endpoint (default 5s).                                                                                                                                                                                           | string                                                                                                                       |
| remoteWrite.basicAuth         | Allow an endpoint to authenticate over basic authentication.                                                                                                                                                                                                 | *v1beta1.BasicAuth                                                                                                           |
| remoteWrite.tlsConfig         | Describes tls configuration for remote write target.                                                                                                                                                                                                         | *v1beta1.TLSConfig                                                                                                           |
| remoteRead                    | Optional URL to read vmalert state (persisted via RemoteWrite) This configuration only makes sense if alerts state has been successfully persisted (via RemoteWrite) before.                                                                                 | []v1beta1.VMAlertRemoteReadSpec                                                                                              |
| remoteRead.url                | URL of the endpoint to send samples to.                                                                                                                                                                                                                      | string                                                                                                                       |
| remoteRead.basicAuth          | Allow an endpoint to authenticate over basic authentication.                                                                                                                                                                                                 | *v1beta1.BasicAuth                                                                                                           |
| remoteRead.tlsConfig          | Describes tls configuration for remote write target.                                                                                                                                                                                                         | *v1beta1.TLSConfig                                                                                                           |
| rulePath                      | RulePath to the file with alert rules. Supports patterns. Flag can be specified multiple times.                                                                                                                                                              | []string                                                                                                                     |
| notifier                      | Notifier prometheus alertmanager endpoint spec.                                                                                                                                                                                                              | *v1beta1.VMAlertNotifierSpec                                                                                                 |
| notifier.url                  | AlertManager url.                                                                                                                                                                                                                                            | string                                                                                                                       |
| notifier.basicAuth            | Allow notifier to authenticate over basic authentication.                                                                                                                                                                                                    | *v1beta1.BasicAuth                                                                                                           |
| notifier.tlsConfig            | Describes describes tls configuration for notifier.                                                                                                                                                                                                          | *v1beta1.TLSConfig                                                                                                           |
| notifier.selector             | Allows service discovery for alertmanager in this case all matched vmalertmanager replicas will be added into vmalert notifier.url as statefulset pod.fqdn                                                                                                   | *metav1.LabelSelector                                                                                                        |
| notifiers                     | Notifier prometheus alertmanager endpoint specs.                                                                                                                                                                                                             | []v1beta1.VMAlertNotifierSpec                                                                                                |
| notifierConfigRef             | NotifierConfigRef reference for secret with notifier configuration for vmalert. Only one of notifier options could be chosen: notifierConfigRef or notifiers +  notifier.                                                                                    | *v1.SecretKeySelector                                                                                                        |
| datasource                    | Datasource Victoria Metrics.                                                                                                                                                                                                                                 | *vmetricsv1b1.VMAlertDatasourceSpec                                                                                          |
| datasource.url                | Victoria Metrics url.                                                                                                                                                                                                                                        | string                                                                                                                       |
| datasource.basicAuth          | Allow datasource to authenticate over basic authentication.                                                                                                                                                                                                  | *v1beta1.BasicAuth                                                                                                           |
| datasource.tlsConfig          | Describes describes tls configuration for datasource target.                                                                                                                                                                                                 | *v1beta1.TLSConfig                                                                                                           |
| paused                        | Set paused to reconciliation for vmalert                                                                                                                                                                                                                     | boolean                                                                                                                      |
| terminationGracePeriodSeconds | Defines period for container graceful termination.                                                                                                                                                                                                           | string                                                                                                                       |
| externalLabels                | The labels to add to any time series. Example: ["team=monitoring", "project=development", "environment=dr311qa", "cluster=dr311qa"]                                                                                                                          | list                                                                                                                         |
| nodeSelector                  | Defines which nodes the pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\"                                                                                                                                           | map[string]string                                                                                                            |
| labels                        | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                                                             | map[string]string                                                                                                            |
| annotations                   | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value"                                       | map[string]string                                                                                                            |
| extraArgs                     | List of key-value with additional arguments which will be passed to VMAgent pod [https://docs.victoriametrics.com/#list-of-command-line-flags](https://docs.victoriametrics.com/#list-of-command-line-flags)                                                 | map[string]string                                                                                                            |
| extraEnvs                     | Allow to set extra system environment variables for vmagent.                                                                                                                                                                                                 | map[string]string                                                                                                            |
| volumes                       | Volumes allows configuration of additional volumes on the output deploy definition. Volumes specified will be appended to other volumes that are generated as a result of StorageSpec objects                                                                | []v1.Volume                                                                                                                  |
| volumeMounts                  | VolumeMounts allows configuration of additional VolumeMounts on the output deploy definition. VolumeMounts specified will be appended to other VolumeMounts in the vmsingle container that are generated as a result of StorageSpec objects                  | []v1.VolumeMount                                                                                                             |
| priorityClassName             | PriorityClassName assigned to the Pods to prevent them from evicting                                                                                                                                                                                         | string                                                                                                                       |
| tlsConfig                     | TLS configuration for VMAlert. Must be specified if `victoriametrics.tlsEnabled` is set to `true`                                                                                                                                                            | [TLSConfig](#tls-config)                                                                                                     |
<!-- markdownlint-enable line-length -->

```yaml
  vmAlert:
    install: true
    ingress:
      host: vmalert.test.org
      install: true
    paused: false
    selectAllByDefault: true
    datasource:
      url: http://vmsingle-vmsingle.monitoring:8429
    notifier:
      url: http://vmalertmanager-vmalertmanager.monitoring:9093
    remoteWrite:
      url: http://vmsingle-vmsingle.monitoring:8429
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
  clusterIssuerName: dev-cluster-issuer
  vmAlert:
    install: true
    ingress:
      host: vmalert.test.org
      install: true
    paused: false
    selectAllByDefault: true
    datasource:
      url: http://vmsingle-vmsingle.monitoring:8429
    notifier:
      url: http://vmalertmanager-vmalertmanager.monitoring:9093
    remoteWrite:
      url: http://vmsingle-vmsingle.monitoring:8429
    securityContext:
      runAsUser: 2001
      fsGroup: 2001
    secrets:
      - kube-etcd-client-certs
    volumes: {}
    volumeMounts: {}
    tlsConfig:
      generateCerts:
        enabled: true
        duration: 365
        renewBefore: 15
        secretName: "vmalert-tls-secret"
```

