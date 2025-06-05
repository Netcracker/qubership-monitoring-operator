### stackdriver-exporter

<!-- markdownlint-disable line-length -->
| Field                               | Description                                                                                                                                                                                                                                                  | Scheme                                                                                                                       |
| ----------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ---------------------------------------------------------------------------------------------------------------------------- |
| install                             | Allows to disable deploy stackdriver-exporter.                                                                                                                                                                                                               | bool                                                                                                                         |
| replicas                            | Number of created pods.                                                                                                                                                                                                                                      | int                                                                                                                          |
| name                                | A deployment name for stackdriver-exporter                                                                                                                                                                                                                   | string                                                                                                                       |
| image                               | A docker image to use for stackdriver-exporter deployment                                                                                                                                                                                                    | string                                                                                                                       |
| imagePullPolicy                     | Image pull policy to use for stackdriver-exporter deployment                                                                                                                                                                                                 | string                                                                                                                       |
| resources                           | The resources that describe the compute resource requests and limits for single pods.                                                                                                                                                                        | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcerequirements-v1-core) |
| serviceAccount.create               | Allow to disable create ServiceAccount during deploy                                                                                                                                                                                                         | bool                                                                                                                         |
| serviceAccount.name                 | Provide a name in place of stackdriver-exporter for ServiceAccount                                                                                                                                                                                           | string                                                                                                                       |
| serviceAccount.annotations          | Will add the provided map to the annotations for the created serviceAccount.                                                                                                                                                                                 | map[string]string                                                                                                            |
| service.type                        | Type of promitor-agent-scraper service.                                                                                                                                                                                                                      | string                                                                                                                       |
| service.port                        | Port of promitor-agent-scraper which use in service.                                                                                                                                                                                                         | integer                                                                                                                      |
| service.labels                      | Annotations set which will be created in service.                                                                                                                                                                                                            | map[string]string                                                                                                            |
| nodeSelector                        | Defines which nodes the pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\"                                                                                                                                           | map[string]string                                                                                                            |
| annotations                         | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value"                                       | map[string]string                                                                                                            |
| labels                              | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                                                             | map[string]string                                                                                                            |
| securityContext                     | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                                                           | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#podsecuritycontext-v1-core)    |
| tolerations                         | Tolerations allow the pods to schedule onto nodes with matching taints.                                                                                                                                                                                      | []v1.Toleration                                                                                                              |
| affinity                            | It specifies the pod's scheduling constraints. For more information, refer to [https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core) | *v1.Affinity                                                                                                                 |
| priorityClassName                   | PriorityClassName assigned to the Pods to prevent them from evicting.                                                                                                                                                                                        | string                                                                                                                       |
| serviceMonitor                      | ServiceMonitor holds configuration attributes for stackdriver-exporter.                                                                                                                                                                                      | object                                                                                                                       |
| stackdriver.projectId               | The Google Project ID to gather metrics for                                                                                                                                                                                                                  | string                                                                                                                       |
| stackdriver.serviceAccountSecret    | An existing secret which contains credentials.json                                                                                                                                                                                                           | string                                                                                                                       |
| stackdriver.serviceAccountSecretKey | Provide custom key for the existing secret to load credentials.json from                                                                                                                                                                                     | string                                                                                                                       |
| stackdriver.serviceAccountKey       | A service account key JSON file. Must be provided when no existing secret is used, in this case a new secret will be created holding this service account                                                                                                    | string                                                                                                                       |
| stackdriver.maxRetries              | Max number of retries that should be attempted on 503 errors from Stackdriver                                                                                                                                                                                | string                                                                                                                       |
| stackdriver.httpTimeout             | How long should Stackdriver_exporter wait for a result from the Stackdriver API                                                                                                                                                                              | string                                                                                                                       |
| stackdriver.maxBackoff              | Max time between each request in an exp backoff scenario                                                                                                                                                                                                     | string                                                                                                                       |
| stackdriver.backoffJitter           | The amount of jitter to introduce in an exp backoff scenario                                                                                                                                                                                                 | string                                                                                                                       |
| stackdriver.retryStatuses           | The HTTP statuses that should trigger a retry                                                                                                                                                                                                                | string                                                                                                                       |
| stackdriver.dropDelegatedProjects   | Drop metrics from attached projects and fetch `project_id` only                                                                                                                                                                                              | string                                                                                                                       |
| stackdriver.metrics.typePrefixes    | The prefixes to gather metrics for, we default to just CPU metrics                                                                                                                                                                                           | string                                                                                                                       |
| stackdriver.metrics.interval        | The frequency to request                                                                                                                                                                                                                                     | string                                                                                                                       |
| stackdriver.metrics.offset          | How far into the past to offset                                                                                                                                                                                                                              | string                                                                                                                       |
| stackdriver.metrics.ingestDelay     | Offset for the Google Stackdriver Monitoring Metrics interval into the past by the ingest delay from the metric's metadata.                                                                                                                                  | string                                                                                                                       |
<!-- markdownlint-enable line-length -->

Example:

<!-- markdownlint-disable line-length -->
```yaml
stackdriverExporter:
  install: true
  name: "stackdriver-exporter"
  replicas: 1
  image: "prometheuscommunity/stackdriver-exporter:v0.12.0"
  imagePullPolicy: IfNotPresent
  resources:
    requests:
      cpu: 100m
      memory: 128Mi
    limits:
      cpu: 100m
      memory: 128Mi
  securityContext:
    runAsUser: 2001
    fsGroup: 2001
  service:
    type: ClusterIP
    port: 9255
    annotations: {}
  stackdriver:
    projectId: "my-test-project"
    serviceAccountSecret: ""
    serviceAccountSecretKey: ""
    serviceAccountKey: |
      {
        "type": "service_account",
        "project_id": "dummy-acc",
        "private_key_id": "00000000000000000000000000000000000",
        "private_key": "",
        "client_email": "dummy@dummy-acc.iam.gserviceaccount.com",
        "client_id": "0000000000000000000000000",
        "auth_uri": "https://accounts.google.com/o/oauth2/auth",
        "token_uri": "https://oauth2.googleapis.com/token",
        "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
        "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/dummy%dummy-acc.iam.gserviceaccount.com"
      }
    maxRetries: 0
    httpTimeout: 10s
    maxBackoff: 5s
    backoffJitter: 1s
    retryStatuses: 503
    dropDelegatedProjects: false
    metrics:
      typePrefixes: 'compute.googleapis.com/instance/cpu,compute.googleapis.com/instance/disk,pubsub.googleapis.com/subscription'
      interval: '5m'
      offset: '0s'
      ingestDelay: false
  affinity: {}
  extraArgs: []
    - --monitoring.filters='pubsub.googleapis.com/subscription:resource.labels.subscription_id=monitoring.regex.full_match("us-west4.*my-team-subs.*")'
  nodeSelector: {}
  tolerations: []
  serviceAccount:
    create: true
    name: stackdriver-exporter
    annotations: {}
  serviceMonitor:
    install: true
    interval: 2m
    telemetryPath: /metrics
    labels:
      app.kubernetes.io/component: monitoring
    timeout: 30s
    relabelings: []
    metricRelabelings: []
  annotations: {}
  labels: {}
  priorityClassName: priority-class
```
<!-- markdownlint-enable line-length -->


