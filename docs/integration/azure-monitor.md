This document provides information about various integration options with Azure Monitor tools
and managed services.

# Azure Monitor

Supported features matrix:

| **Monitoring System** | **Metrics** | **Dashboards** | **Alerting** | **Autoscaling** |
| --------------------- | ----------- | -------------- | ------------ | --------------- |
| Azure Monitor         | ✓ Yes       | ✗ No           | ✗ No         | ✗ No            |

**Legend**:

* `Yes` - The feature is supported and implemented.
* `No` - The feature is not supported by the target monitoring system.
* `-` - The feature is not applicable.

Integration between Azure monitor and Monitoring Stack should be
configured **before** deploying Monitoring, and it is required to set correct deployment parameters to enable
necessary features.


## Integration with Azure Monitor

_Analysis is in progress._


## Promitor Agent Scraper

The Promitor Agent Scraper is used to scrape metrics from Azure managed services in the Prometheus format.

### How to Install Promitor Agent Scraper

Configure the Promitor correctly to collect metrics from Azure Monitor.
For the full list of installation parameters, refer to the
[Platform Monitoring Installation Procedure](../installation/README.md#promitor-agent-scraper).

Set up some parameters for Azure authentication and configure scraping of metrics
from the right tenant.

### How to Configure Azure Authentication

Promitor Agent Scraper provides the following authentication mechanisms:

* `Service principle` - Use the application ID & secret of the Azure AD entity
  that has been pre-created to authenticate with (default).
* `Managed Identity` - Use zero-secret authentication by letting Microsoft handle the authentication for you.
  For more information, see [https://docs.microsoft.com/en-us/azure/active-directory/managed-identities-azure-resources/overview](https://docs.microsoft.com/en-us/azure/active-directory/managed-identities-azure-resources/overview).

You can change the authentication mode through the `azureAuthentication.mode` parameter.

#### Service Principle Authentication Mode

For the Service principle method, the agent uses the application ID & secret of the Azure AD entity.
You have to create the App Registration in the Azure AD.

To create the application in the Azure Active Directory for Promitor in the Azure UI:

1. Navigate to `Azure Active Directory > App registrations > New registration` and create the application.
2. Navigate to the created application (for example, "Promitor").
3. Copy `Application (client) ID` and use it as the `azureAuthentication.identity.id` parameter in the configuration.
4. Navigate to `Certificates & secrets > Client secrets`.
5. Click `+ New client secret` and create a new client secret.
6. Copy `Value` and use it as the `azureAuthentication.identity.key` parameter in the configuration.

To assign a role to the created application:

1. Navigate to `Resource Groups > <resource-group-name> > Access control (IAM) > + Add > Add Role Assignment`.
2. Add the `Monitoring Reader` role to the earlier created application (for example, "Promitor").
3. Add the `Reader` role if you want to use a `Promitor Resource Discovery`

If you set the `azureAuthentication.identity.key` parameter and set `secrets.createSecret`
to `true` (`true` by default), a secret with the name `secrets.secretName`
(`promitor-agent-scraper` by default) is created during the deployment. This secret contains
a field with the name `secrets.appKeySecret` (`azure-app-key` by default) and
`azureAuthentication.identity.key` as the value for this field.

If you want to use the previously created secret instead of using the `azureAuthentication.identity.key`
parameter, you should set `secrets.createSecret` to `false`, set the name of the secret as
`secrets.secretName`, and set the name of the field that contains the application secret as
`secrets.appKeySecret`.


#### Managed Identity Authentication Mode

_Analysis is in progress._


### Metric Declaration for Promitor Agent Scraper

Enter the `azureMetadata` section in the configuration correctly to use Promitor Agent Scraper.

To get information for the `azureMetadata` section from the Azure UI:

1. Navigate to `Azure Active Directory > Overview`.
2. Copy `Tenant ID` and use it as the `azureMetadata.tenantId` parameter in the configuration.
3. Navigate to `Resource Groups > <resource-group-name> > Overview`.
4. Copy `Subscription ID` and use it as the `azureMetadata.subscriptionId` parameter in the configuration.
5. Copy `<resource-group-name>` and use it as the `azureMetadata.resourceGroupName` parameter
   in the configuration.
6. (Optional) You can change the `azureMetadata.cloud` parameter from `Global` to `China`, `UsGov`, or `Germany`.

You can set `metricDefaults.aggregation.interval` (`00:05:00` by default)
and `metricDefaults.scraping.schedule` (`*/5 * * * *` by default - one scrape per 5 minute).
These values are used for metrics by default and can override each metric separately.
For more information about the parameters, refer to the [Platform Monitoring Installation Procedure](../installation/README.md#promitor-agent-scraper).

To collect metrics from Azure Monitor, set up the `metrics` parameter.
This parameter is a list of objects. Each item in the list describe one metric that should be
collected from Azure.

Every metric that is being declared needs to define the following fields:

* `name` - Name of the metric that is reported.
* `description` - Description for the metric that is reported.
* `resourceType` - Defines what type of resource needs to be queried.
* `azureMetricConfiguration.metricName` - The name of the metric in Azure Monitor to query.
* `azureMetricConfiguration.aggregation.type` - The aggregation that needs to be used when
  querying Azure Monitor.
* `azureMetricConfiguration.aggregation.interval` - Overrides `metricDefaults.aggregation.interval`
  `metricDefaults.aggregation.interval` with a new interval.
* `resources` - An array of one or more resources to get the metrics for.
  The fields required vary depending on the `resourceType` being created, and are documented
  for each resource in the official documentation at [https://docs.promitor.io/latest/scraping/overview/](https://docs.promitor.io/latest/scraping/overview/).
* `azureMetricConfiguration.limit` - The maximum amount of resources to scrape when using
  dimensions or filters.
* `resourceDiscoveryGroups` - An array of one or more resource discovery groups
  that are used to automatically discover all resources through Promitor Resource Discovery.
  For every found resource, it gets the metrics and reports them.

All resources provide the capability to override the default Azure metadata:

* `subscriptionId` - Changes the subscription ID to which the resource belongs. Overrides `azureMetadata.subscriptionId`.
* `resourceGroupName` - Changes the resource group that contains the resource. Overrides `azureMetadata.resourceGroupName`.

Additionally, the following fields are optional:

* `azureMetricConfiguration.dimension.name` - The name of the dimension that should be used to scrape
  a multi-dimensional metric in Azure Monitor.
* `labels` - Defines a set of custom labels to include for a given metric.
* `scraping.schedule` - A scraping schedule for the individual metric. Overrides `metricDefaults.scraping.schedule`.

For more information about metrics, see the list of supported Azure resources for scraping at
[https://docs.promitor.io/latest/scraping/overview/](https://docs.promitor.io/latest/scraping/overview/).

For all metrics supported by Azure monitor, see
[https://docs.microsoft.com/en-us/azure/azure-monitor/essentials/metrics-supported](https://docs.microsoft.com/en-us/azure/azure-monitor/essentials/metrics-supported).

**Note**: The configuration of each metric must contain either `resources`
(if you installed Promitor Scraper without Promitor Resource Discovery),
or `resourceDiscoveryGroups` (if you use Promitor Resource Discovery).

Following are examples of metric declaration for some Azure resources. Replace the words
inside the `<>` symbols with names of your own resources and leave only one of the fields:
`resources` or `resourceDiscoveryGroups` to use these configurations.

#### Azure Kubernetes Service

```yaml
name: azure_kubernetes_available_cpu_cores
description: "Available CPU cores in cluster"
resourceType: KubernetesService
azureMetricConfiguration:
  metricName: kube_node_status_allocatable_cpu_cores
  aggregation:
    type: Average
# Required when no resource discovery is configured
resources:
- clusterName: <cluster_name>
# Requires Promitor Resource Discovery agent
# (https://promitor.io/concepts/how-it-works#using-resource-discovery)
resourceDiscoveryGroups:
- name: <discovery_group_name>
```

#### Azure Event Hubs (Kafka)

```yaml
name: azure_event_hubs_incoming_messages
description: "The number of incoming messages on an Azure Event Hubs topic"
resourceType: EventHubs
azureMetricConfiguration:
  metricName: IncomingMessages
  aggregation:
    type: Total
# Required when no resource discovery is configured
resources:
  - namespace: <namespace_1>
    topicName: <topic_name_1>
  - namespace: <namespace_2>
    topicName: <topic_name_2>
# Requires Promitor Resource Discovery agent
# (https://promitor.io/concepts/how-it-works#using-resource-discovery)
resourceDiscoveryGroups:
  - name: <discovery_group_name>
```

#### Azure Database for PostgreSQL

```yaml
name: azure_postgre_sql_cpu_percent
description: "The CPU percentage on the server"
resourceType: PostgreSql
scraping:
  schedule: "0 */2 * ? * *"
azureMetricConfiguration:
  metricName: cpu_percent
  aggregation:
    type: Average
    interval: 00:01:00
# Required when no resource discovery is configured
resources:
  - serverName: <server_name_1>
  - serverName: <server_name_2>
# Requires Promitor Resource Discovery agent
# (https://promitor.io/concepts/how-it-works#using-resource-discovery)
resourceDiscoveryGroups:
  - name: <discovery_group_name>
```

#### Azure Cache for Redis

```yaml
name: azure_redis_cache_cache_hits
description: "The number of successful key lookups during the specified reporting interval. This maps to keyspace_hits from the Redis INFO command."
resourceType: RedisCache
scraping:
  schedule: "0 */2 * ? * *"
azureMetricConfiguration:
  metricName: CacheHits
  aggregation:
    type: Total
    interval: 00:01:00
# Required when no resource discovery is configured
resources:
  - cacheName: <cache_name_1>
  - cacheName: <cache_name_2>
# Requires Promitor Resource Discovery agent
# (https://promitor.io/concepts/how-it-works#using-resource-discovery)
resourceDiscoveryGroups:
  - name: <discovery_group_name>
```

#### Azure SQL Elastic Pool

```yaml
name: promitor_demo_sql_elastic_pool_cpu
description: "CPU percentage used for a Azure SQL Elastic Pool"
resourceType: SqlElasticPool
labels:
  app: promitor
azureMetricConfiguration:
  metricName: cpu_percent
  aggregation:
    type: Average
# Required when no resource discovery is configured
resources:
  - serverName: <server_name>
    poolName: <pool_name>
# Requires Promitor Resource Discovery agent
# (https://promitor.io/concepts/how-it-works#using-resource-discovery)
resourceDiscoveryGroups:
  - name: <discovery_group_name>
```

#### Azure Cosmos DB (Cassandra/MongoDB)

```yaml
name: azure_cosmos_db_total_requests
description: "Demo cosmos query"
resourceType: CosmosDb
azureMetricConfiguration:
  metricName: TotalRequests
  aggregation:
    type: Count
# Required when no resource discovery is configured
resources:
  - dbName: <db_name_1>
  - dbName: <db_name_1>
# Requires Promitor Resource Discovery agent
# (https://promitor.io/concepts/how-it-works#using-resource-discovery)
resourceDiscoveryGroups:
  - name: <discovery_group_name>
```

## Promitor Agent Resource Discovery

### How to Install Promitor Agent Resource Discovery

Promitor Agent Resource Discovery can be installed [manually](https://docs.promitor.io/v2.9/deployment/resource-discovery/)
or as a part of Platform Monitoring.
For the full list of installation parameters, refer to the
[Platform Monitoring Installation Procedure](../installation/README.md#promitor-agent-resource-discovery).
It requires Azure Authentication that must be configured the same as promitor-agent-scraper, see
[below](#how-to-configure-azure-authentication).

### How to Configure Resource Discovery Groups

The explaining how to configure resource groups can be found
[here](https://docs.promitor.io/v2.9/resource-discovery/declaring-resource-discovery-groups/).
