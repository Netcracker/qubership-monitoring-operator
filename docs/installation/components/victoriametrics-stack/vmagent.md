#### vmagent

<!-- markdownlint-disable line-length -->
| Field                                         | Description                                                                                                                                                                                                                                                                                                                                  | Scheme                                                                                                                       |
| --------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| install                                       | Allows to enable or disable deploy vmagent.                                                                                                                                                                                                                                                                                                  | boolean                                                                                                                      |
| image                                         | A Docker image to deploy the vmagent.                                                                                                                                                                                                                                                                                                        | string                                                                                                                       |
| resources                                     | The resources that describe the compute resource requests and limits for single pods.                                                                                                                                                                                                                                                        | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcerequirements-v1-core) |
| securityContext                               | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                                                                                                                                           | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#podsecuritycontext-v1-core)    |
| containers                                    | Containers property allows to inject additions sidecars or to patch existing containers. It can be useful for proxies, backup, etc.                                                                                                                                                                                                          | [[]v1.Container](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#container-v1-core)                     |
| paused                                        | Set paused to reconciliation for vmagent                                                                                                                                                                                                                                                                                                     | boolean                                                                                                                      |
| tolerations                                   | Tolerations allow the pods to schedule onto nodes with matching taints.                                                                                                                                                                                                                                                                      | []v1.Toleration                                                                                                              |
| nodeSelector                                  | Defines which nodes the pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\"                                                                                                                                                                                                                           | map[string]string                                                                                                            |
| affinity                                            | If specified, the pod's scheduling constraints                                                                                                                                                                                      | *v1.Affinity                                                                                                                   |
| labels                                        | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                                                                                                                                             | map[string]string                                                                                                            |
| annotations                                   | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value"                                                                                                                       | map[string]string                                                                                                            |
| ingress                                       | Ingress allows to create Ingress for the vmagent UI.                                                                                                                                                                                                                                                                                         | *[Ingress](#ingress)                                                                                                         |
| scrapeInterval                                | Defines how often scrape targets by default.                                                                                                                                                                                                                                                                                                 | string                                                                                                                       |
| additionalScrape                              | The additionalScrape allows you to monitor hosts outside the cloud (balancers, graylog server, jenkins...)                                                                                                                                                                                                                                   | *v1.SecretKeySelector, [Victoria Metrics Config](https://docs.victoriametrics.com/operator/additional-scrape.html)           |
| vmAgentExternalLabelName                      | Name of vmAgent external label used to denote VmAgent instance name. Defaults to the value of `vmagent`. External label will _not_ be added when value is set to empty string (`""`).                                                                                                                                                        | string                                                                                                                       |
| externalLabels                                | The labels to add to any time series scraped by vmagent. It doesn't affect metrics ingested directly by push API's                                                                                                                                                                                                                           | map[string]string                                                                                                            |
| remoteWriteSettings                           | Defines global settings for all remoteWrite urls. For more information, refer to [https://docs.victoriametrics.com/vmagent.html](https://docs.victoriametrics.com/vmagent.html)                                                                                                                                                              | *v1beta1.VMAgentRemoteWriteSettings                                                                                          |
| remoteWriteSettings.maxDiskUsagePerURL        | The maximum file-based buffer size in bytes at -remoteWrite.tmpDataPath for each -remoteWrite.url. Disk usage is unlimited if the value is set to 0. For more information, refer to [https://docs.victoriametrics.com/vmagent.html](https://docs.victoriametrics.com/vmagent.html)                                                           | string                                                                                                                       |
| remoteWriteSettings.showURL                   | Whether to show -remoteWrite.url in the exported metrics. It is hidden by default, since it can contain sensitive auth info                                                                                                                                                                                                                  | string                                                                                                                       |
| remoteWriteSettings.label                     | Optional labels in the form 'name=value' to add to all the metrics before sending them                                                                                                                                                                                                                                                       | string                                                                                                                       |
| remoteWrite                                   | RemoteWriteSpec defines the remote_write configuration for vmagent. The `remote_write` allows to transparently send samples to a long term storage. For more information, refer to [https://docs.victoriametrics.com/vmagent.html](https://docs.victoriametrics.com/vmagent.html)                                                            | []v1beta1.VMAgentRemoteWriteSpec                                                                                             |
| remoteWrite.url                               | URL of the endpoint to send samples to.                                                                                                                                                                                                                                                                                                      | string                                                                                                                       |
| remoteWrite.basicAuth                         | Allow an endpoint to authenticate over basic authentication.                                                                                                                                                                                                                                                                                 | *v1beta1.BasicAuth                                                                                                           |
| remoteWrite.basicAuth.createSecret            | Allow create secret for basic authentification automatically.                                                                                                                                                                                                                                                                                | map[string]string                                                                                                            |
| remoteWrite.basicAuth.createSecret.secretName | Name of the secret with which will be created.                                                                                                                                                                                                                                                                                               | string                                                                                                                       |
| remoteWrite.basicAuth.createSecret.username   | Username for basic authentification.                                                                                                                                                                                                                                                                                                         | string                                                                                                                       |
| remoteWrite.basicAuth.createSecret.password   | Password for basic authentification.                                                                                                                                                                                                                                                                                                         | string                                                                                                                       |
| remoteWrite.bearerTokenSecret                 | Optional bearer auth token to use for -remoteWrite.url.                                                                                                                                                                                                                                                                                      | *v1.SecretKeySelector                                                                                                        |
| remoteWrite.oauth2                            | Defines auth configuration.                                                                                                                                                                                                                                                                                                                  | *v1beta1.OAuth2                                                                                                              |
| remoteWrite.tlsConfig                         | Describes tls configuration for remote write target.                                                                                                                                                                                                                                                                                         | *v1beta1.OAuth2                                                                                                              |
| remoteWrite.urlRelabelConfig                  | ConfigMap with relabeling config which is applied to metrics before sending them to the corresponding -remoteWrite.url [https://github.com/cybozu-go/VictoriaMetrics-operator/blob/master/docs/relabeling.MD](https://github.com/cybozu-go/VictoriaMetrics-operator/blob/master/docs/relabeling.MD)                                          | *v1.ConfigMapKeySelector                                                                                                     |
| remoteWrite.inlineUrlRelabelConfig            | Defines relabeling config for remoteWriteURL, it can be defined at crd spec. [https://github.com/cybozu-go/VictoriaMetrics-operator/blob/master/docs/relabeling.MD](https://github.com/cybozu-go/VictoriaMetrics-operator/blob/master/docs/relabeling.MD)                                                                                    | []RelabelConfig                                                                                                              |
| urlRelabelConfig                              | ConfigMap with global relabel config -remoteWrite.relabelConfig. This relabeling is applied to all the collected metrics before sending them to remote storage. [https://github.com/cybozu-go/VictoriaMetrics-operator/blob/master/docs/relabeling.MD](https://github.com/cybozu-go/VictoriaMetrics-operator/blob/master/docs/relabeling.MD) | *v1.ConfigMapKeySelector                                                                                                     |
| inlineUrlRelabelConfig                        | Defines GlobalRelabelConfig for vmagent, can be defined directly at CRD. [https://github.com/cybozu-go/VictoriaMetrics-operator/blob/master/docs/relabeling.MD](https://github.com/cybozu-go/VictoriaMetrics-operator/blob/master/docs/relabeling.MD)                                                                                        | []RelabelConfig                                                                                                              |
| podMonitorNamespaceSelector                   | Defines Namespaces to be selected for VMPodScrape discovery.                                                                                                                                                                                                                                                                                 | *metav1.LabelSelector                                                                                                        |
| serviceMonitorNamespaceSelector               | Defines Namespaces to be selected for VMServiceScrape discovery.                                                                                                                                                                                                                                                                             | *metav1.LabelSelector                                                                                                        |
| podMonitorSelector                            | Defines PodScrapes to be selected for target discovery.                                                                                                                                                                                                                                                                                      | *metav1.LabelSelector                                                                                                        |
| serviceMonitorSelector                        | Defines ServiceScrapes to be selected for target discovery.                                                                                                                                                                                                                                                                                  | *metav1.LabelSelector                                                                                                        |
| extraArgs                                     | List of key-value with additional arguments which will be passed to VMAgent pod [https://docs.victoriametrics.com/#list-of-command-line-flags](https://docs.victoriametrics.com/#list-of-command-line-flags)                                                                                                                                 | map[string]string                                                                                                            |
| extraEnvs                                     | Allow to set extra system environment variables for vmagent.                                                                                                                                                                                                                                                                                 | []v1.EnvVar                                                                                                                  |
| secrets                                       | Secrets is a list of Secrets in the same namespace as the VMSingle object, which shall be mounted into the vmagent pods                                                                                                                                                                                                                      | []string                                                                                                                     |
| volumes                                       | Volumes allows configuration of additional volumes on the output deploy definition. Volumes specified will be appended to other volumes that are generated as a result of StorageSpec objects                                                                                                                                                | []v1.Volume                                                                                                                  |
| volumeMounts                                  | VolumeMounts allows configuration of additional VolumeMounts on the output deploy definition. VolumeMounts specified will be appended to other VolumeMounts in the vmsingle container that are generated as a result of StorageSpec objects                                                                                                  | []v1.VolumeMount                                                                                                             |
| maxScrapeInterval                             | MaxScrapeInterval allows limiting maximum scrape interval for VMServiceScrape, VMPodScrape and other scrapes                                                                                                                                                                                                                                 | string                                                                                                                       |
| minScrapeInterval                             | MinScrapeInterval allows limiting minimal scrape interval for VMServiceScrape, VMPodScrape and other scrapes                                                                                                                                                                                                                                 | string                                                                                                                       |
| priorityClassName                             | PriorityClassName assigned to the Pods to prevent them from evicting                                                                                                                                                                                                                                                                         | string                                                                                                                       |
| tlsConfig                                     | TLS configuration for VMAgent. Must be specified if `victoriametrics.tlsEnabled` is set to `true`                                                                                                                                                                                                                                            | [TLSConfig](#tls-config)                                                                                                     |
<!-- markdownlint-enable line-length -->

Example:

```yaml
  vmAgent:
    install: true
    route:
      install: true
      host: vmagent.test.org
    logLevel: WARN
    serviceMonitor:
      interval: 29s
      scrapeTimeout: 9s
    extraArgs:
      loggerFormat: default
    remoteWriteSettings:
      label:
        project: development
        environment: qa-kuber-test
        cluster: qa-kuber-test
        team: test_team-test
      showURL: true
    remoteWrite:
      - url: https://remotehost-1:8429/api/v1/write
        basicAuth:
          username:
            key: username
            name: vmsecret
          password:
            key: password
            name: vmsecret
        tlsConfig:
          insecureSkipVerify: true
        urlRelabelConfig:
          name: "vmagent-relabel"
          key: "target-1-relabel.yaml"          
        inlineUrlRelabelConfig:
          - source_labels: [__name__]
            separator: ;
            regex: kube.*
            replacement: $1
            action: keep
      - url: https://remotehost-2:8429/api/v1/write
        basicAuth:
          createSecret:
            username: user
            password: pass
            secretName: vmsecret-basicauth
        tlsConfig:
          insecureSkipVerify: true
        urlRelabelConfig:
          name: "vmagent-relabel"
          key: "target-1-relabel.yaml"          
        inlineUrlRelabelConfig:
          - source_labels: [__name__]
            separator: ;
            regex: kube.*
            replacement: $1
            action: keep
    nodeSelector:
      role: compute
    labels:
      label: label_vmagent
    annotations:
      annotation: test_vmagent
    secrets:
      - kube-etcd-client-certs
```

If TLS is enabled for Victoriametrics:

```yaml
victoriametrics:
  tlsEnabled: true
  clusterIssuerName: dev-cluster-issuer
  vmAgent:
    install: true
    route:
      install: true
      host: vmagent.test.org
    logLevel: WARN
    serviceMonitor:
      interval: 29s
      scrapeTimeout: 9s
    extraArgs:
      loggerFormat: default
    remoteWriteSettings:
      label:
        project: development
        environment: qa-kuber-test
        cluster: qa-kuber-test
        team: test_team-test
      showURL: true
    remoteWrite:
      - url: https://remotehost-1:8429/api/v1/write
        basicAuth:
          username:
            key: username
            name: vmsecret
          password:
            key: password
            name: vmsecret
        tlsConfig:
          insecureSkipVerify: true
        urlRelabelConfig:
          name: "vmagent-relabel"
          key: "target-1-relabel.yaml"
        inlineUrlRelabelConfig:
          - source_labels: [__name__]
            separator: ;
            regex: kube.*
            replacement: $1
            action: keep
      - url: https://remotehost-2:8429/api/v1/write
        basicAuth:
          createSecret:
            username: user
            password: pass
            secretName: vmsecret-basicauth
        tlsConfig:
          insecureSkipVerify: true
        urlRelabelConfig:
          name: "vmagent-relabel"
          key: "target-1-relabel.yaml"          
        inlineUrlRelabelConfig:
          - source_labels: [__name__]
            separator: ;
            regex: kube.*
            replacement: $1
            action: keep
    nodeSelector:
      role: compute
    labels:
      label: label_vmagent
    annotations:
      annotation: test_vmagent
    secrets:
      - kube-etcd-client-certs
    tlsConfig:
      generateCerts:
        enabled: false
      createSecret:
        secretName: "vmagent-tls-secret"
        ca: |-
          -----BEGIN CERTIFICATE-----
          ...
          -----END CERTIFICATE-----
        key: |-
          -----BEGIN PRIVATE KEY-----
          ...
          -----END PRIVATE KEY-----
        cert: |-
          -----BEGIN CERTIFICATE-----
          ...
          ------END CERTIFICATE-----
```

