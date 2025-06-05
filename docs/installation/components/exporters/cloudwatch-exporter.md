### cloudwatch-exporter

cloudwatchExporter is a specification of the desired deployment of cloudwatch-exporter.

**Note**: Pay attention to settings that can be specified under the `cloudwatchExporter.config` parameter, such as,

```yaml
cloudwatchExporter:
  serviceMonitor:
    interval: 2m
  config:
    period_seconds: 120
    delay_seconds: 60
```

These parameters are very important to get the actual metrics from CloudWatch without a delay. The default values of
`delay_seconds` is `600 s`, and this parameter is used to avoid collecting data that has not fully converged.
`600 s` specifies that the exporter fetches data only after `600 s = 10 m`. To decrease the delay, it is recommended
to decrease the value to `60 s`. This time is enough for CloudWatch to collect the metrics and allows
cloudwatch_exporter to fetch them.

Refer to the official documentation of cloudwatch_exporter for full descriptions of all parameters at
[https://github.com/prometheus/cloudwatch_exporter#configuration](https://github.com/prometheus/cloudwatch_exporter#configuration).

<!-- markdownlint-disable line-length -->
| Field                                       | Description                                                                                                                                                                                                                                                                                                           | Scheme                                                                                                                       |
| ------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------- |
| install                                     | Allows to disable deploy cloudwatch-exporter.                                                                                                                                                                                                                                                                         | bool                                                                                                                         |
| replicas                                    | Number of created pods.                                                                                                                                                                                                                                                                                               | int                                                                                                                          |
| name                                        | A deployment name for cloudwatch-exporter                                                                                                                                                                                                                                                                             | string                                                                                                                       |
| image                                       | A docker image to use for cloudwatch-exporter deployment                                                                                                                                                                                                                                                              | string                                                                                                                       |
| imagePullPolicy                             | Image pull policy to use for cloudwatch-exporter deployment                                                                                                                                                                                                                                                           | string                                                                                                                       |
| command                                     | Allow override command to run docker container                                                                                                                                                                                                                                                                        | []string                                                                                                                     |
| resources                                   | The resources that describe the compute resource requests and limits for single pods.                                                                                                                                                                                                                                 | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcerequirements-v1-core) |
| serviceAccount.install                      | Allow to disable create ServiceAccount during deploy                                                                                                                                                                                                                                                                  | bool                                                                                                                         |
| serviceAccount.name                         | Provide a name in place of cloudwatch-exporter for ServiceAccount                                                                                                                                                                                                                                                     | bool                                                                                                                         |
| serviceAccount.automountServiceAccountToken | Specifies whether to automount API credentials for the ServiceAccount.                                                                                                                                                                                                                                                | bool                                                                                                                         |
| rbac.createClusterRole                      | Allow creating ClusterRole. If set to `false`, ClusterRole must be created manually. Default: `true`                                                                                                                                                                                                                  | bool                                                                                                                         |
| rbac.createClusterRoleBinding               | Allow creating ClusterRoleBinding. If set to `false`, ClusterRoleBinding must be created manually. Default: `true`                                                                                                                                                                                                    | bool                                                                                                                         |
| nodeSelector                                | Defines which nodes the pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\"                                                                                                                                                                                                    | map[string]string                                                                                                            |
| annotations                                 | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value"                                                                                                | map[string]string                                                                                                            |
| labels                                      | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                                                                                                                      | map[string]string                                                                                                            |
| securityContext                             | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 65534, fsGroup: 65534 }`.                                                                                                                                                                                  | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#podsecuritycontext-v1-core)    |
| tolerations                                 | Tolerations allow the pods to schedule onto nodes with matching taints.                                                                                                                                                                                                                                               | []v1.Toleration                                                                                                              |
| affinity                                    | It specifies the pod's scheduling constraints. For more information, refer to [https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core)                                                          | *v1.Affinity                                                                                                                 |
| serviceMonitor                              | ServiceMonitor holds configuration attributes for cloudwatch-exporter.                                                                                                                                                                                                                                                | object                                                                                                                       |
| aws.aws_access_key_id                       | AWS Access Key ID for programmatic access. Do not specify the `aws_access_key_id` and `aws_secret_access_key` if you specified role or `.aws.secret.name` before.                                                                                                                                                     | string                                                                                                                       |
| aws.aws_secret_access_key                   | AWS Secret Access Key for programmatic access. Do not specify the `aws_access_key_id` and `aws_secret_access_key` if you specified role or `.aws.secret.name` before.                                                                                                                                                 | string                                                                                                                       |
| aws.secret.name                             | The name of a pre-created secret in which AWS credentials are stored. When set, `aws_access_key_id` is assumed to be in a field called `access_key`, `aws_secret_access_key` is assumed to be in a field called `secret_key`, and the session token, if it exists, is assumed to be in a field called security_token. | string                                                                                                                       |
| aws.secret.includesSessionToken             | Allow specify manually generated token and exporter will not try to use STS for get token. When set, token is assumed to be in a field called `security_token`.                                                                                                                                                       | bool                                                                                                                         |
| config                                      | Configuration is rendered with `tpl` function, therefore you can use any Helm variables and/or templates.                                                                                                                                                                                                             | string                                                                                                                       |
| priorityClassName                           | PriorityClassName assigned to the Pods to prevent them from evicting.                                                                                                                                                                                                                                                 | string                                                                                                                       |
<!-- markdownlint-enable line-length -->

Example:

```yaml
cloudwatchExporter:
  install: true
  replicas: 1
  name: cloudwatch-exporter
  image: prom/cloudwatch-exporter:cloudwatch_exporter-0.11.0
  imagePullPolicy: IfNotPresent
  command:
    - 'java'
    - '-Dhttp.proxyHost=proxy.example.com'
    - '-Dhttp.proxyPort=3128'
    - '-Dhttps.proxyHost=proxy.example.com'
    - '-Dhttps.proxyPort=3128'
    - '-jar'
    - '/cloudwatch_exporter.jar'
    - '9106'
    - '/config/config.yml'
  resources:
    limits:
      cpu: 200m
      memory: 256Mi
    requests:
      cpu: 100m
      memory: 128Mi
  serviceAccount:
    install: true
    name: cloudwatch-exporter
    annotations:
      eks.amazonaws.com/role-arn: arn:aws:iam::1234567890:role/prom-cloudwatch-exporter-oidc
    automountServiceAccountToken: true
  rbac:
    createClusterRole: true
    createClusterRoleBinding: true
  nodeSelector:
    node-role.kubernetes.io/worker: worker
  labels:
    label.key: label-value
  annotations:
    annotation.key: annotation-value
  securityContext:
    runAsUser: 65534  # run as nobody user instead of root
    fsGroup: 65534  # necessary to be able to read the EKS IAM token
  tolerations:
    - key: "example-key"
      operator: "Exists"
      effect: "NoSchedule"
  affinity: {}
  serviceMonitor:
    install: true
    interval: 5m
    telemetryPath: /metrics
    labels: {}
    timeout: 10s
    relabelings: []
    metricRelabelings:
      - sourceLabels: [dbinstance_identifier]
        action: replace
        replacement: mydbname
        targetLabel: dbname
  aws:
    secret:
      name: <secret_name>
      includesSessionToken: false
    aws_access_key_id: <access_key_id>
    aws_secret_access_key: <secret_access_key>
  config: |-
