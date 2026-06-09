### cert-exporter

**Note**: Pay attention to the specifics of deploy with `cert-exporter`. For more information, refer to
the [Deploy with cert-exporter](#deploy-with-cert-exporter) section.

<!-- markdownlint-disable line-length -->
| Field                                   | Description                                                                                                                                                                                                            | Scheme                                                                                                                       |
| --------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| install                                 | Allows to enable or disable deploy cert-exporter.                                                                                                                                                                      | bool                                                                                                                         |
| name                                    | A name of the microservice to deploy with. This name is used as the name of the microservice deployment and in labels.                                                                                                 | string                                                                                                                       |
| setupSecurityContext                    | Allows to create PodSecurityPolicy or SecurityContextConstraints.                                                                                                                                                      | bool                                                                                                                         |
| setupGrafanaDashboard                   | Allows to create Grafana dashboard for cert-exporter.                                                                                                                                                                  | bool                                                                                                                         |
| setupAlertingRules                      | Allows to create Prometheus alerting rules for cert-exporter.                                                                                                                                                          | bool                                                                                                                         |
| additionalHostPathVolumes               | Allows to mount additional directories or files from the host file system to container. Should be used when the exporter needs to collect information from files from non-default paths.                               | list[object]                                                                                                                 |
| additionalHostPathVolumes[N].volumeName | Name of volume, must be unique. Reserved names (shouldn't be used): kube, openshift-origin, openshift-etcd, kubelet-pki, root-kube-config.                                                                             | string                                                                                                                       |
| additionalHostPathVolumes[N].volumePath | Path to the directory or file on the host system. The same path is used in the container.                                                                                                                              | string                                                                                                                       |
| certsInFiles                            | Settings for parsing certificates from host file system.                                                                                                                                                               | object                                                                                                                       |
| certsInFiles.enabled                    | Enables parsing certificates from host file system. If true, the part of the exporter will be deployed as a DaemonSet.                                                                                                 | bool                                                                                                                         |
| certsInFiles.defaultCerts               | Allows to check internal certificates from default paths which is default for Kubernetes or Openshift clusters.                                                                                                        | bool                                                                                                                         |
| certsInFiles.includeCerts               | Allows to check certificates with custom regex.                                                                                                                                                                        | string                                                                                                                       |
| certsInFiles.excludeCerts               | Allows to check only certificates that NOT match regex.                                                                                                                                                                | string                                                                                                                       |
| certsInKubeconfig                       | Settings for parsing certificates from Kubeconfig files on the host file system.                                                                                                                                       | object                                                                                                                       |
| certsInKubeconfig.enabled               | Enables parsing certificates from Kubeconfig. If true, the part of the exporter will be deployed as a DaemonSet.                                                                                                       | bool                                                                                                                         |
| certsInKubeconfig.defaultCerts          | Allows to check Kubeconfig files from default paths on Kubernetes or Openshift clusters.                                                                                                                               | bool                                                                                                                         |
| certsInKubeconfig.includeCerts          | Allows to check Kubeconfig files with custom regex.                                                                                                                                                                    | string                                                                                                                       |
| certsInKubeconfig.excludeCerts          | Allows to check only Kubeconfig files that NOT match regex.                                                                                                                                                            | string                                                                                                                       |
| certsInSecrets                          | Settings for parsing certificates from Kubernetes secrets.                                                                                                                                                             | object                                                                                                                       |
| certsInSecrets.enabled                  | Enables parsing certificates from Kubernetes secrets. If true, the part of the exporter will be deployed as a Deployment.                                                                                              | bool                                                                                                                         |
| certsInSecrets.includeCerts             | Allows to check fields in the secrets with keys that match regex.                                                                                                                                                      | string                                                                                                                       |
| certsInSecrets.excludeCerts             | Allows to check only fields in the secrets with keys NOT match regex.                                                                                                                                                  | string                                                                                                                       |
| certsInSecrets.annotationSelector       | Allows to match secrets by annotation.                                                                                                                                                                                 | string                                                                                                                       |
| certsInSecrets.labelSelector            | Allows to match secrets by label.                                                                                                                                                                                      | string                                                                                                                       |
| certsInSecrets.namespaces               | Allows to find secrets in the selected namespaces. Kubernetes comma-delimited list of namespaces to search for secrets. Empty string specifies that the exporter checks all available namespaces.                      | string                                                                                                                       |
| certsInSecrets.types                    | Allows to select only specific secret type. An empty list specifies that the exporter checks all available secrets.                                                                                                    | list[string]                                                                                                                 |
| certsInSecrets.kubeconfigPath           | Allows to specify path to kubeconfig file for getting access to secrets via kubectl. Only required if out-of-cluster installation.                                                                                     | string                                                                                                                       |
| pollingPeriod                           | Periodic interval in which to check certs. Format: [time.Duration](https://golang.org/pkg/time/#ParseDuration) from GoLang                                                                                             | string                                                                                                                       |
| image                                   | A Docker image to deploy the cert-exporter.                                                                                                                                                                            | string                                                                                                                       |
| serviceMonitor                          | Service monitor for pulling metrics.                                                                                                                                                                                   | object                                                                                                                       |
| serviceMonitor.install                  | Allows to install serviceMonitor.                                                                                                                                                                                      | bool                                                                                                                         |
| serviceMonitor.interval                 | Allow change metrics scrape interval.                                                                                                                                                                                  | string                                                                                                                       |
| servicePort                             | Port for cert-exporter service.                                                                                                                                                                                        | int                                                                                                                          |
| daemonset                               | Pod-specific settings for cert-exporter as a daemonset. Works if collecting from files and/or kubeconfig is enabled.                                                                                                   | object                                                                                                                       |
| daemonset.resources                     | The resources that describe the compute resource requests and limits for single pods. Affects daemonset pods.                                                                                                          | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcerequirements-v1-core) |
| daemonset.extraArgs                     | Additional arguments for cert-exporter containers from daemonset.                                                                                                                                                      | list[string]                                                                                                                 |
| daemonset.securityContext               | SecurityContext holds pod-level security attributes.Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                      | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#podsecuritycontext-v1-core)    |
| daemonset.tolerations                   | Tolerations allow the pods to schedule onto nodes with matching taints.                                                                                                                                                | []v1.Toleration                                                                                                              |
| daemonset.nodeSelector                  | Defines which nodes the pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\"                                                                                                     | map[string]string                                                                                                            |
| daemonset.affinity                      | If specified, the pod's scheduling constraints                                                                                                                                                                         | *v1.Affinity                                                                                                                 |
| daemonset.annotations                   | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value" | map[string]string                                                                                                            |
| daemonset.labels                        | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                       | map[string]string                                                                                                            |
| daemonset.priorityClassName             | PriorityClassName assigned to the Pods to prevent them from evicting.                                                                                                                                                  | string                                                                                                                       |
| deployment                              | Settings for parsing certificates from Kubernetes secrets.                                                                                                                                                             | object                                                                                                                       |
| deployment.resources                    | The resources that describe the compute resource requests and limits for single pods. Affects deployment pods.                                                                                                         | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcerequirements-v1-core) |
| deployment.extraArgs                    | Additional arguments for cert-exporter containers from deployment.                                                                                                                                                     | list[string]                                                                                                                 |
| deployment.securityContext              | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                     | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#podsecuritycontext-v1-core)    |
| deployment.tolerations                  | Tolerations allow the pods to schedule onto nodes with matching taints.                                                                                                                                                | []v1.Toleration                                                                                                              |
| deployment.nodeSelector                 | Defines which nodes the pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\"                                                                                                     | map[string]string                                                                                                            |
| deployment.affinity                     | If specified, the pod's scheduling constraints                                                                                                                                                                         | *v1.Affinity                                                                                                                 |
| deployment.annotations                  | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value" | map[string]string                                                                                                            |
| deployment.labels                       | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                       | map[string]string                                                                                                            |
| deployment.priorityClassName            | PriorityClassName assigned to the Pods to prevent them from evicting.                                                                                                                                                  | string                                                                                                                       |
<!-- markdownlint-enable line-length -->

Example of standard installation without root rights:

```yaml
certExporter:
  ...
  install: true
  name: cert-exporter
  setupSecurityContext: true
  setupGrafanaDashboard: true
  setupAlertingRules: true
  additionalHostPathVolumes: []
  certsInFiles:
    enabled: true
    defaultCerts: true
  certsInKubeconfig:
    enabled: false
  certsInSecrets:
    enabled: true
  pollingPeriod: 1h0m0s
  image: joeelliott/cert-exporter:v2.7.0
  serviceMonitor:
    install: true
    interval: 30s
  servicePort: 9219
  daemonset:
    resources:
      limits:
        cpu: 20m
        memory: 50Mi
      requests:
        cpu: 10m
        memory: 25Mi
    extraArgs: []
    securityContext:
      runAsUser: "0"
      fsGroup: "0"
    tolerations:
      - operator: "Exists"
    nodeSelector:
      node-role.kubernetes.io/worker: worker
    affinity: {}
    labels:
      label.key: label-value
    annotations:
      annotation.key: annotation-value
    priorityClassName: priority-class
  deployment:
    resources:
      limits:
        cpu: 20m
        memory: 150Mi
      requests:
        cpu: 10m
        memory: 50Mi
    extraArgs: []
    securityContext:
      runAsUser: 2000
      fsGroup: 2000
    tolerations: []
    nodeSelector:
      node-role.kubernetes.io/worker: worker
    affinity: {}
    labels:
      label.key: label-value
    annotations:
      annotation.key: annotation-value
    priorityClassName: priority-class
```

Example of installation with root rights:

```yaml
certExporter:
  ...
  install: true
  name: cert-exporter
  setupSecurityContext: true
  setupGrafanaDashboard: true
  setupAlertingRules: true
  certsInFiles:
    enabled: true
    defaultCerts: true
  certsInKubeconfig:
    enabled: true
    defaultCerts: true
  certsInSecrets:
    enabled: true
    kubeconfigPath: "/root/.kube/config"
  pollingPeriod: 1h0m0s
  image: joeelliott/cert-exporter:v2.7.0
  securityContext:
    runAsUser: "0"
```

Example of installation with non-default paths of files:

```yaml
certExporter:
  ...
  additionalHostPathVolumes:
            - volumeName: unique-volume-name-1
              volumePath: /path/to/certificates
                            - volumeName: unique-volume-name-2
                            volumePath: /path/to/kubeconfig
  certsInFiles:
    enabled: true
    defaultCerts: true
    includeCerts: "/path/to/certificates/*.{crt,cert}"
  certsInKubeconfig:
    enabled: true
    defaultCerts: false
    includeCerts: "/path/to/kubeconfig/**/*.conf"
  certsInSecrets:
    enabled: true
    kubeconfigPath: "/path/to/kubeconfig/config"
  ...
```

Example of installation with specifying parameters for parsing secrets:

```yaml
certExporter:
  ...
  certsInSecrets:
    enabled: true
    includeCerts: "*.crt"
    excludeCerts: "do-not-check.crt"
    annotationSelector: "some-annotation"
    labelSelector: "some-label"
    namespaces: "monitoring, namespace_1"
    types:
      - "Opaque"
      - "kubernetes.io/tls"
    kubeconfigPath: "/root/.kube/config"
  ...
  deployment:
    extraArgs:
      - '--secrets-include-glob=*.cer'
      - '--secrets-include-glob=*.pem'
  ...
```


