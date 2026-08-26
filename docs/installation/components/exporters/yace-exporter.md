# Yace-exporter

`yaceExporter` is a specification of the desired deployment of `yace-exporter`
(Yet Another CloudWatch Exporter).

`yace-exporter` collects AWS CloudWatch metrics through the `GetMetricData` API, with tag-based auto-discovery,
cross-account roles, and static jobs. It is an alternative to the [cloudwatch-exporter](cloudwatch-exporter.md), it is
disabled by default, and you can enable it independently of `cloudwatch-exporter`.

Refer to the official documentation of YACE for full descriptions of all configuration parameters at
[https://github.com/prometheus-community/yet-another-cloudwatch-exporter/blob/master/docs/configuration.md](https://github.com/prometheus-community/yet-another-cloudwatch-exporter/blob/master/docs/configuration.md).

## Authentication

The exporter authenticates with one home identity. Pick one of the following options:

* **IRSA or instance profile.** Leave `aws.secret.name`, `aws.aws_access_key_id`, and `aws.aws_secret_access_key`
  unset, and attach the IAM role through `serviceAccount.annotations`.
* **Pre-created Secret.** Set `aws.secret.name`. The Secret must contain the keys `access_key` and `secret_key`, plus
  `security_token` when `aws.secret.includesSessionToken` is `true`.
* **Chart-managed keys.** Set both `aws.aws_access_key_id` and `aws.aws_secret_access_key`. The chart creates the
  Secret. Setting only one of the two keys fails the render with an explicit message.

Extra AWS accounts are not extra credentials. List them as `roles[].roleArn` on a job inside `config`, and allow the
home identity to call `sts:AssumeRole` on those ARNs.

The IAM policy of the home identity needs `cloudwatch:GetMetricData`, `cloudwatch:GetMetricStatistics`,
`cloudwatch:ListMetrics`, and `tag:GetResources`.

## Parameters

<!-- markdownlint-disable line-length -->
| Field                                       | Description                                                                                                                                                                                                                                          | Scheme                                                                                                                       |
|---------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------|
| install                                     | Allows to disable deploy yace-exporter. Default: `false`.                                                                                                                                                                                            | bool                                                                                                                         |
| replicas                                    | Number of created pods.                                                                                                                                                                                                                              | int                                                                                                                          |
| name                                        | A deployment name for yace-exporter.                                                                                                                                                                                                                 | string                                                                                                                       |
| image                                       | A Docker image to use for yace-exporter deployment.                                                                                                                                                                                                  | string                                                                                                                       |
| imagePullPolicy                             | Image pull policy to use for yace-exporter deployment.                                                                                                                                                                                               | string                                                                                                                       |
| command                                     | Allow override command to run Docker container. When set, it replaces the entrypoint, so pass `-listen-address` and `-config.file` yourself.                                                                                                          | []string                                                                                                                     |
| extraArgs                                   | Additional command-line flags appended to the default arguments, for example `-scraping-interval` or `-cloudwatch-concurrency`.                                                                                                                       | []string                                                                                                                     |
| listenAddress                               | Network address the HTTP server binds to. Default: `0.0.0.0:5000`. YACE binds to `127.0.0.1` by default, which is unreachable from Prometheus.                                                                                                        | string                                                                                                                       |
| containerPort                               | Container port. Must match the port in `listenAddress`. Default: `5000`.                                                                                                                                                                             | int                                                                                                                          |
| configFile                                  | Path passed to `-config.file`. Default: `/config/config.yml`.                                                                                                                                                                                        | string                                                                                                                       |
| resources                                   | The resources that describe the compute resource requests and limits for single pods.                                                                                                                                                                | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcerequirements-v1-core) |
| serviceAccount.install                      | Allow to disable create ServiceAccount during deploy.                                                                                                                                                                                                | bool                                                                                                                         |
| serviceAccount.name                         | Provide a name in place of yace-exporter for ServiceAccount.                                                                                                                                                                                         | string                                                                                                                       |
| serviceAccount.annotations                  | Annotations for the ServiceAccount, for example the IRSA annotation `eks.amazonaws.com/role-arn`.                                                                                                                                                     | map[string]string                                                                                                            |
| serviceAccount.automountServiceAccountToken | Specifies whether to automount API credentials for the ServiceAccount.                                                                                                                                                                               | bool                                                                                                                         |
| rbac.createClusterRole                      | Allow creating ClusterRole. The exporter does not call the Kubernetes API, so leave it disabled unless an external process needs `get` on the yace-exporter Secret and ConfigMap. Default: `false`.                                                   | bool                                                                                                                         |
| rbac.createClusterRoleBinding               | Allow creating ClusterRoleBinding. Set it together with `rbac.createClusterRole`. Default: `false`.                                                                                                                                                   | bool                                                                                                                         |
| nodeSelector                                | Defines which nodes the pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\"                                                                                                                                    | map[string]string                                                                                                            |
| annotations                                 | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value"                                | map[string]string                                                                                                            |
| labels                                      | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                                                      | map[string]string                                                                                                            |
| securityContext                             | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 65534, fsGroup: 65534 }`. `fsGroup` is required to read the EKS IAM token.                                                                 | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#podsecuritycontext-v1-core)    |
| tolerations                                 | Tolerations allow the pods to schedule onto nodes with matching taints.                                                                                                                                                                              | []v1.Toleration                                                                                                              |
| affinity                                    | It specifies the pod's scheduling constraints. For more information, refer to [https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core) | *v1.Affinity                                                                                                                 |
| serviceMonitor                              | ServiceMonitor holds configuration attributes for yace-exporter.                                                                                                                                                                                     | object                                                                                                                       |
| aws.secret.name                             | Name of a pre-created Secret with AWS credentials in the keys `access_key`, `secret_key`, and optionally `security_token`.                                                                                                                            | string                                                                                                                       |
| aws.secret.includesSessionToken             | Whether the pre-created Secret also holds an STS session token under the `security_token` key.                                                                                                                                                        | bool                                                                                                                         |
| aws.aws_access_key_id                       | AWS Access Key ID for chart-managed credentials. Do not set it when using `aws.secret.name` or IRSA.                                                                                                                                                  | string                                                                                                                       |
| aws.aws_secret_access_key                   | AWS Secret Access Key for chart-managed credentials. Do not set it when using `aws.secret.name` or IRSA.                                                                                                                                              | string                                                                                                                       |
| config                                      | The YACE `v1alpha1` configuration, rendered with the `tpl` function. It defines `discovery`, `static`, and `customNamespace` jobs. Add extra AWS accounts as `roles[].roleArn` on a job and extra regions as `regions` on a job. `sts-region` is a single STS endpoint. | string                                                                                                                       |
| priorityClassName                           | PriorityClassName assigned to the Pods to prevent them from evicting.                                                                                                                                                                                | string                                                                                                                       |
<!-- markdownlint-enable line-length -->

## Example

```yaml
yaceExporter:
  install: true
  replicas: 1
  name: yace-exporter
  image: docker.io/prometheuscommunity/yet-another-cloudwatch-exporter:v0.67.0
  imagePullPolicy: IfNotPresent
  resources:
    limits:
      cpu: 200m
      memory: 256Mi
    requests:
      cpu: 100m
      memory: 128Mi
  serviceAccount:
    install: true
    name: yace-exporter
    annotations:
      eks.amazonaws.com/role-arn: arn:aws:iam::1234567890:role/yace-exporter-oidc
    automountServiceAccountToken: true
  securityContext:
    runAsUser: 65534  # run as nobody user instead of root
    fsGroup: 65534  # necessary to be able to read the EKS IAM token
  serviceMonitor:
    install: true
    interval: 5m
    telemetryPath: /metrics
  config: |-
    apiVersion: v1alpha1
    sts-region: us-east-1
    discovery:
      jobs:
        - type: AWS/RDS
          regions:
            - us-east-1
          searchTags:
            - key: environment
              value: production
          period: 300
          length: 600
          metrics:
            - name: DatabaseConnections
              statistics: [Average]
            - name: CPUUtilization
              statistics: [Average]
```

## Multiple accounts and regions

One Deployment covers every account and region. YACE assumes each role listed in `roles` and scrapes it in each region
listed in `regions`, so the job below produces four scrape targets.

```yaml
yaceExporter:
  install: true
  serviceAccount:
    annotations:
      eks.amazonaws.com/role-arn: arn:aws:iam::<monitoring-account>:role/<irsa-role>
      eks.amazonaws.com/sts-regional-endpoints: "true"
  config: |-
    apiVersion: v1alpha1
    sts-region: us-east-1
    customNamespace:
      - name: rds-all
        namespace: AWS/RDS
        regions:
          - us-east-1
          - eu-central-1
        roles:
          - roleArn: "arn:aws:iam::111111111111:role/yace-cloudwatch"
          - roleArn: "arn:aws:iam::222222222222:role/yace-cloudwatch"
        dimensionNameRequirements: [DBInstanceIdentifier]
        period: 120
        length: 180
        delay: 60
        metrics:
          - name: CPUUtilization
            statistics: [Average]
```

On clusters without IRSA, keep the same `config` and set `yaceExporter.aws` instead of `serviceAccount.annotations`.

For longer examples, see
[multi-account-multi-region.yaml](../../../examples/components/yace-exporter-config/multi-account-multi-region.yaml)
and the [yace-exporter deploy parameters](../../../examples/deploy-parameters/yace-exporter).
