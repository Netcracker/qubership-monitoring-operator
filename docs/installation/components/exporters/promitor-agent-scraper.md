### promitor-agent-scraper

<!-- markdownlint-disable line-length -->
| Field                                                         | Description                                                                                                                                                                                                                                                  | Scheme                                                                                                                       |
| ------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ---------------------------------------------------------------------------------------------------------------------------- |
| install                                                       | Allows to enable or disable deploy promitor-agent-scraper.                                                                                                                                                                                                   | boolean                                                                                                                      |
| replicas                                                      | Number of created pods.                                                                                                                                                                                                                                      | integer                                                                                                                      |
| name                                                          | A name of the microservice to deploy with. This name is used as the name of the microservice deployment and in labels.                                                                                                                                       | string                                                                                                                       |
| image                                                         | A Docker image to deploy the promitor-agent-scraper.                                                                                                                                                                                                         | string                                                                                                                       |
| imagePullPolicy                                               | Image pull policy to use for promitor-agent-scraper deployment.                                                                                                                                                                                              | string                                                                                                                       |
| azureAuthentication.mode                                      | Authentication type to use to authenticate. Options are ServicePrincipal (default), UserAssignedManagedIdentity or SystemAssignedManagedIdentity.                                                                                                            | string                                                                                                                       |
| azureAuthentication.identity.id                               | Id of the Azure AD entity to authenticate with.                                                                                                                                                                                                              | string                                                                                                                       |
| azureAuthentication.identity.key                              | Secret of the Azure AD entity to authenticate with. Sets the environment variable PROMITOR_AUTH_APPKEY through the secrets.appKeySecret field in the Secret.                                                                                                 | string                                                                                                                       |
| azureAuthentication.identity.binding                          | Aad Pod Identity name, when using UserAssignedManagedIdentity or SystemAssignedManagedIdentity as mode.                                                                                                                                                      | string                                                                                                                       |
| metricSinks.prometheusScrapingEndpoint.baseUriPath            | Controls the path where the scraping endpoint for Prometheus is being exposed.                                                                                                                                                                               | string                                                                                                                       |
| metricSinks.prometheusScrapingEndpoint.enableMetricTimestamps | Defines whether or not a timestamp should be included when the value was scraped on Azure Monitor. Supported values are True to opt-in & False to opt-out.                                                                                                   | boolean                                                                                                                      |
| metricSinks.prometheusScrapingEndpoint.metricUnavailableValue | Defines the value that will be reported if a metric is unavailable.                                                                                                                                                                                          | string                                                                                                                       |
| metricSinks.prometheusScrapingEndpoint.labelTransformation    | Controls how label values are reported to Prometheus by using transformation. Options are None & Lowercase.                                                                                                                                                  | string                                                                                                                       |
| resourceDiscovery.enabled                                     | Indication whether or not resource discovery is enabled through the Promitor Resource Discovery agent.                                                                                                                                                       | boolean                                                                                                                      |
| resourceDiscovery.host                                        | DNS name or IP address of the Promitor Resource Discovery agent.                                                                                                                                                                                             | string                                                                                                                       |
| resourceDiscovery.port                                        | Port (UDP) address of the Promitor Resource Discovery agent.                                                                                                                                                                                                 | integer                                                                                                                      |
| telemetry.defaultLogLevel                                     | Defines the default minimum log level that should be logged if a sink does not provide one. Allowed values are Trace, Debug, Information, Warning, Error, Critical, None ordered from most to least verbose.                                                 | string                                                                                                                       |
| telemetry.applicationInsights.enabled                         | Determines if the sink is used or not.                                                                                                                                                                                                                       | boolean                                                                                                                      |
| telemetry.applicationInsights.key                             | Defines the instrumentation key to use when sending telemetry to Azure Application Insights.                                                                                                                                                                 | string                                                                                                                       |
| telemetry.applicationInsights.logLevel                        | Verbosity to use for this sink, if not specified then the telemetry.defaultLogLevel will be used.                                                                                                                                                            | string                                                                                                                       |
| telemetry.containerLogs.enabled                               | Determines if the sink is used or not.                                                                                                                                                                                                                       | boolean                                                                                                                      |
| telemetry.containerLogs.logLevel                              | Verbosity to use for this sink, if not specified then the telemetry.defaultLogLevel will be used.                                                                                                                                                            | string                                                                                                                       |
| azureMetadata.tenantId                                        | The id of the Azure tenant that will be queried.                                                                                                                                                                                                             | string                                                                                                                       |
| azureMetadata.subscriptionId                                  | The id of the default subscription to query.                                                                                                                                                                                                                 | string                                                                                                                       |
| azureMetadata.resourceGroupName                               | The name of the default resource group to query.                                                                                                                                                                                                             | string                                                                                                                       |
| azureMetadata.cloud                                           | The name of the Azure cloud to use. Options are Global (default), China, UsGov & Germany.                                                                                                                                                                    | string                                                                                                                       |
| metricDefaults.aggregation.interval                           | The default interval which defines over what period measurements of a metric should be aggregated. A cron that fits your needs.                                                                                                                              | string                                                                                                                       |
| metricDefaults.scraping.schedule                              | A cron expression that controls the frequency of which all the configured metrics will be scraped from Azure Monitor. You can use crontab-generator.org to generate a cron that fits your needs.                                                             | string                                                                                                                       |
| metrics                                                       | List of metrics to scrape.                                                                                                                                                                                                                                   | list[[Metric](https://docs.promitor.io/latest/scraping/overview/)]                                                         |
| secrets.createSecret                                          | To use your own secret, set createSecret to false and define the name/keys that your secret uses. Indication if you want to bring your own secret level of logging.                                                                                          | boolean                                                                                                                      |
| secrets.secretName                                            | Name of the secret for Azure AD identity secret.                                                                                                                                                                                                             | string                                                                                                                       |
| secrets.appKeySecret                                          | Name of the field for Azure AD identity secret in the Secret resource.                                                                                                                                                                                       | string                                                                                                                       |
| service.type                                                  | Type of promitor-agent-scraper service.                                                                                                                                                                                                                      | string                                                                                                                       |
| service.port                                                  | Port of promitor-agent-scraper which use in service.                                                                                                                                                                                                         | integer                                                                                                                      |
| service.targetPort                                            | Target port of promitor-agent-scraper which use in service and container.                                                                                                                                                                                    | integer                                                                                                                      |
| service.labels                                                | Labels set which will create in service.                                                                                                                                                                                                                     | map[string]string                                                                                                            |
| service.azureLoadBalancer.dnsPrefix                           | Prefix for DNS name to expose the service on using \<name\>.\<location\>.cloudapp.azure.com format. This setting is specific to Azure Kubernetes Service.                                                                                                    | string                                                                                                                       |
| service.azureLoadBalancer.exposeInternally                    | To restrict access to Promitor by exposing it through an internal load balancer. This setting is specific to Azure Kubernetes Service.                                                                                                                       | string                                                                                                                       |
| serviceAccount.install                                        | Allow to disable create ServiceAccount during deploy.                                                                                                                                                                                                        | boolean                                                                                                                      |
| serviceAccount.name                                           | Provide a name in place of promitor-agent-scraper for ServiceAccount.                                                                                                                                                                                        | string                                                                                                                       |
| serviceAccount.annotations                                    | Will add the provided map to the annotations for the created serviceAccount.                                                                                                                                                                                 | map[string]string                                                                                                            |
| serviceAccount.labels                                         | Will add the provided map to the labels for the created serviceAccount.                                                                                                                                                                                      | map[string]string                                                                                                            |
| serviceMonitor                                                | ServiceMonitor holds configuration attributes for promitor-agent-scraper.                                                                                                                                                                                    | object                                                                                                                       |
| resources                                                     | The resources that describe the compute resource requests and limits for single pods.                                                                                                                                                                        | [v1.ResourceRequirements](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#resourcerequirements-v1-core) |
| securityContext                                               | SecurityContext holds pod-level security attributes. Default for Kubernetes, `securityContext:{ runAsUser: 2000, fsGroup: 2000 }`.                                                                                                                           | [*v1.PodSecurityContext](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.32/#podsecuritycontext-v1-core)    |
| affinity                                                      | It specifies the pod's scheduling constraints. For more information, refer to [https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.30/#affinity-v1-core) | *v1.Affinity                                                                                                                 |
| nodeSelector                                                  | Defines which nodes the pods are scheduled on. Specified just as map[string]string. For example: \"type: compute\"                                                                                                                                           | map[string]string                                                                                                            |
| tolerations                                                   | Tolerations allow the pods to schedule onto nodes with matching taints.                                                                                                                                                                                      | []v1.Toleration                                                                                                              |
| annotations                                                   | Map of string keys and values stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. Specified just as map[string]string. For example: "annotations-key: annotation-value"                                       | map[string]string                                                                                                            |
| labels                                                        | Map of string keys and values that can be used to organize and categorize (scope and select) objects. Specified just as map[string]string. For example: "label-key: label-value"                                                                             | map[string]string                                                                                                            |
| priorityClassName                                             | PriorityClassName assigned to the Pods to prevent them from evicting.                                                                                                                                                                                        | string                                                                                                                       |
<!-- markdownlint-enable line-length -->

Example:

```yaml
promitorAgentScraper:
  install: true
  replicaCount: 1
  name: promitor-agent-scraper
  image: ghcr.io/tomkerkhove/promitor-agent-scraper:2.5.1
  imagePullPolicy: IfNotPresent
  azureAuthentication:
    mode: "ServicePrincipal"
    identity:
      id: "<azure-ad-id>"
      key: "<azure-ad-secret>"
      binding: ""
  metricSinks:
    prometheusScrapingEndpoint:
      baseUriPath: /metrics
      enableMetricTimestamps: true
      metricUnavailableValue: NaN
      labelTransformation: None
      enableServiceDiscovery: true
  resourceDiscovery:
    enabled: false
    host: ""
    port: 80
  telemetry:
    defaultLogLevel: "Error"
    applicationInsights:
      enabled: false
      key: ""
      logLevel: ""
    containerLogs:
      enabled: true
      logLevel: ""
  azureMetadata:
    tenantId: "<tenant-id>"
    subscriptionId: "<subscription-id>"
    resourceGroupName: "<resource-group-name>"
    cloud: "Global"
  metricDefaults:
    aggregation:
      interval: 00:05:00
    scraping:
      schedule: "*/5 * * * *"
  metrics:
    - name: azure_node_cpu_usage_millicores
      description: "Node CPU usage millicores"
      resourceType: KubernetesService
      azureMetricConfiguration:
        metricName: node_cpu_usage_millicores
        aggregation:
          type: Average
      resources:
        - clusterName: <cluster-name>
  secrets:
    createSecret: true
    secretName: promitor-agent-scraper
    appKeySecret: azure-app-key
  service:
    type: ClusterIP
    port: 8888
    targetPort: 5000
    labels: {}
    azureLoadBalancer:
      dnsPrefix:
      exposeInternally: false
  resources:
    limits:
      cpu: 200m
      memory: 256Mi
    requests:
      cpu: 100m
      memory: 128Mi
  serviceAccount:
    install: true
    name: promitor-agent-scraper
    annotations: {}
    labels: {}
  serviceMonitor:
    install: true
    interval: 5m
    telemetryPath: /metrics
    labels:
      app.kubernetes.io/component: monitoring
  securityContext:
    runAsUser: 2000
    fsGroup: 2000
  affinity: {}
  nodeSelector: {}
  tolerations: []
  annotations: {}
  labels: {}
  priorityClassName: priority-class
```


