This document provides information about troubleshooting Monitoring deployment and its configurations.
It also provides the configurations of monitoring services.

# Deploy

## Manifests Contain a New Resource that Already Exists

Sometimes the following error can be seen during a run deploy using a job or manually using a Helm.

```bash
helm.go:75: [debug] existing resource conflict: kind: <kind>, namespace: <namespace>, name: <name>
rendered manifests contain a new resource that already exists. Unable to continue with [install|update]
```

This error can contain various conflict resources (for example, PlatformMonitoring, Role, ClusterRole, and so on). Such
an error may occur because a resource of this kind already exists in the namespace and with this name.

But it doesn't contain the Helm annotation, so Helm does not track it. For example, annotations:

```yaml
annotations:
  meta.helm.sh/release-name: monitoring-operator-monitoring
  meta.helm.sh/release-namespace: monitoring
```

In this case, Helm cannot install or upgrade this resource because it did not create it and should not change.

## Set Integer Parameter Value to Zero

Due to the nature of work, the Helm templates with any integer value parameters that are equal to zero
(for example, `securityContext.runAsUser` or `securityContext.fsGroup`) must be passed as a string (with quotes).
For example:

```yaml
securityContext:
  runAsUser: "0"
  fsGroup: "0"
```

Otherwise, the default values are set.

# Runtime Issues

## Common Issues

### Prometheus or Grafana Operators are Restarting

Prometheus-operator and Grafana-operator can use a lot of memory during the start if a lot of Custom Resources (CR)
with which the operator works are deployed on the Cloud.

For example, similar issues were observed on the Cloud which had:

* `ServiceMonitor` CRs = **835**
* `GrafanaDashboard` CRs = **347**

In the Cloud with so many CRs, both operators try to allocate:

* `grafana-operator` - more that `400Mi`
* `prometheus-operator` - more than `200Mi`

But after processing all the CRs, their memory usage decreased to:

* `grafana-operator` - `~143M`
* `prometheus-operator` - `~49Mi`

**Warning**: OpenShift clouds with versions `v4.x` require more resources for `prometheus-operator` than the same
configuration in Kubernetes or previous OpenShift versions, so this issue is more likely to occur.

Signs of a problem:

* The operator restarts shortly after starting.
* In operator status or events, records like `OOMKilled` exist.
* Monitoring shows that the operator's memory usage is near the memory limit.

To solve this issue, you need to increase the memory limit for the operator. There is no clear recommendation
on how
much to raise the memory limit since it depends on the CRs count that tries to handle the operator.

However, you can increase the limit to `100Mi`. If necessary, you can raise the limit by `100Mi` several times
until
this solves the problem.

There are two ways to do it:

1. Redeploy and specify increased memory limits in deploy parameters. For more information, refer to
   the [Platform Monitoring Installation Procedure](installation/README.md) in the _Cloud Platform Installation_ guide.
2. Change memory limits directly in the Platform Monitoring CR. For more information, refer to
   the [Platform Monitoring Maintenance](maintenance.md) in the _Cloud Platform Maintenance Guide_.

The memory parameters which can be used to change the memory limits:

```yaml
grafana:
  operator:
    resources:
      limits:
        cpu: 100m
        memory: 100Mi
      requests:
        cpu: 50m
        memory: 50Mi
prometheus:
  operator:
    resources:
      requests:
        cpu: 100m
        memory: 256Mi
      limits:
        cpu: 50m
        memory: 50Mi
```

The parameters' descriptions can be found in the following section of the _Cloud Platform Installation_ guide:

* [Platform Monitoring Installation Procedure](installation/README.md)

### Metrics absent or errors during metrics collection

By default `prometheus/victoriametrics` tries to fetch metrics and endpoints from all the namespaces cluster-wide which contain
`serviceMonitors/podMonitors/probeMonitors` and monitor them. However in some cases filters can be used to specify particular
namespaces for monitoring using labels like `serviceMonitorNamespaceSelector/podMonitorNamespaceSelector/probeNamespaceSelector`
in prometheus specifications. Details for these labels can be found in
[https://github.com/prometheus-operator/prometheus-operator/blob/v0.79.2/Documentation/api.md](https://github.com/prometheus-operator/prometheus-operator/blob/v0.79.2/Documentation/api.md).

Now when our namespace is being monitored, usually we are using the following flow.

First, check the VictoriaMetrics/Prometheus config.

In case of VictoriaMetrics, open Ingress for VMAuth or VMAgent and navigate to `/config`.
In case of Prometheus, open Ingress for Prometheus and navigate to `Status -> Configuration`.

Use the search on the page and try to find the ServiceMonitor/PodMonitor CR name, for example: `postgres-exporter-service-monitor`.
Also, you can check the relabeling rules under the job name. Some of them usually apply as a filter criteria.
While executing discovery it created special metadata labels with following syntax name:

```yaml
relabel_configs:
- action: keep
  source_labels: [__meta_discoveryname_typeofresource_label_labelname]
  regex: <label_value>
```

For example:

```yaml
relabel_configs:
- action: keep
  source_labels: [__meta_kubernetes_service_label_component]
  regex: backup-daemon
```

will filter all containers only with the label on the Service:

```yaml
metadata:
  labels:
    component: backup-daemon
```

Remember or save the job name for your service.

If you already found the problem in this step, then the root cause may be the following:

* You forgot to specify the mandatory label `app.kubernetes.io/component: monitoring` in ServiceMonitor/PodMonitor
* Sometimes configuration will apply during 30s - 3 minutes. It may be necessary to wait a little longer.
* Various errors in the ServiceMonitor/PodMonitor
* Other errors (usually it's very useful to check victoriametrics-operator/prometheus-operator logs)

Second, check the ServiceDiscovery if you found a configuration for your ServiceMonitor/PodMonitor.
VictoriaMetrics and Prometheus use Kubernetes ServiceDiscovery to discover pods in the Kubernetes.

In the case of VictoriaMetrics, open Ingress for VMAuth or VMAgent and navigate to `/service-discovery`.
In the case of Prometheus, open Ingress for Prometheus and navigate to `Status -> ServiceDiscovery`.

You can search by the job name and check that your job:

* Have at least one target
* Check the status of targets

If all targets for your job have the `Dropped` state it means that you incorrectly specified labels or label selectors
in the ServiceMonitor/PodMonitor.
And now your configuration has the filtration criteria that don't allow to use of these targets.

**Note:** The common issue is when labels specified in the ServiceMonitor label selector are not specified on the
Pod and Service

Third, if the pod/targets exist in ServiceDiscovery and have the `UP` status but metrics still do not exist,
it is necessary to check the target status.

In the case of VictoriaMetrics, open Ingress for VMAuth or VMAgent and navigate to `/targets`.
In the case of Prometheus, open Ingress for Prometheus and navigate to `Status -> Targets`.

Need to find your job and check the metrics collection status. If VMAgent or Prometheus can't collect metrics you should
analyze the errors due to which they can't collect metrics.
Usually, the problem at this step means:

* Incorrect port for metrics
* Incorrect endpoint for metrics
* Less timeout value to collect metrics
* The service can't connect to this endpoint/port
* Other errors on the service/application side are out of the scope of our product team.

Usually, the errors at this point are visible in the `Last error` column in active targets and more details can be found
in the `response` link in the case of VictoriaMetrics.

### Unable to update Etcd certificates

In case of restricted privileges monitoring-operator/etcd_monitor_reconciler can throw warn message like:

```bash
"{"lvl":"WARN", ..., "details":"Unable to update etcd certificates due to a lack of permission to access the requested etcd resource."}"
```

It means there is no access to Etcd resources.
The solutions are described in the `monitoring-operator` RBAC document, `the-monitoring-operator-rbac` section.

## Prometheus Stack

### Prometheus Target is Down

If a pod with an application that should expose metrics to Prometheus is running, the corresponding `ServiceMonitor`
(or `PodMonitor`) is created and the corresponding target is present in the Prometheus UI, but that target is down.
You can increase the `scrapeTimeout` parameter for `ServiceMonitor` (or `PodMonitor`). The current scrape duration can be
found in the Prometheus UI on the `targets` page. If the scraped duration approaches the one specified in the parameter,
increasing the timeout may help.

**Note that scrapeTimeout must be less than the interval.**

## VictoriaMetrics Stack

### VictoriaMetrics Operator Can't Remove ClusterRole/ClusterRoleBinding

**When it can occur:**

* You already had deployed Monitoring with VictoriaMetrics stack in the Cloud
* Deploy failed with errors

**Errors in the `victoriametrics-operator`**

```bash
Failed to watch *v1.ClusterRoleBinding: failed to list *v1.ClusterRoleBinding: clusterrolebindings.rbac.authorization.k8s.io is forbidden:
  User "system:serviceaccount:system-monitor:system-monitor-victoriametrics-operator" cannot list resource "clusterrolebindings" in API group "rbac.authorization.k8s.io" at the cluster scope
```

**Root-cause:**

The `victoriametrics-operator` **must handle** the deletion of objects with the `finalizer`. It means that:

* the operator will receive the event that somebody wants to remove the object with the `finalizer`
* the operator must handle it by its logic
* after handling it must remove the `finalizer` from the list of `finalizers`

The run job in the `Clean Install` mode led to the situation when:

* the operator already removed and can't handle finalizers
* some resources were removed after and there is no operator to handle the finalizer
* resources are stuck because there is no operator to handle and remove the finalizer

Kubernetes can't use objects that contain:

```yaml
metadata:
  deletionTimestamp: ....
```

because this object has already been removed and will be physically removed after all controllers handle its remove event.

So to fix this issue we just remove `finalizers` from the objects that were stuck in the delete state and
weren't removed physically.

**How to fix it:**

To fix this issue you must manually remove the `finalizer` for victoriametrics-operator in each object that is stuck
in the delete/termination state.

You need to open editing the object from which you want to remove a finalizer and remove it from the metadata section:

```yaml
metadata:
  finalizers:
    - apps.victoriametrics.com/finalizer
```

Or if an object has no other `finalizers` you can use a patch command to set an empty value
for the whole `finalizers` section:

```bash
kubectl -n <namespace_name> patch <object_type> <object_name> -p '{"metadata":{"finalizers":null}}' --type=merge
```

### EFS Persistence Volume overflow

**When it can occur:**

* VictoriaMetrics fully utilizes EFS Persistence Volume (PV) throughput
* VictoriaMetrics slowly selects data and/or returns query results very long time
* AWS CloudWatch or AWS EFS shows that the throughput of EFS volume is fully utilized

**Root-cause:**

AWS provides different types of EFS divers with different throughputs.

* [https://docs.aws.amazon.com/efs/latest/ug/performance.html](https://docs.aws.amazon.com/efs/latest/ug/performance.html)

By default, many people created EFS volume with `Bursting` throughput type. But usually, it's not enough for
VictoriaMetrics.

**How to fix it:**

Need to select `Provisioned Throughput` mode for EFS volume or at least `Elastic Throughput`. The first option will
charge some additional money. The minimal value of `Provisioned Throughput` is 10 Mb/s.

Also, you can replace AWS EFS with AWS EBS volume. EBS volumes will be cheap and will have more throughput. If you want
to use AWS EBS, please make sure that you will use at least one HDD with type `Throughput Optimized HDD volumes`.

You can use faster disks (like SSD), but `Throughput Optimized HDD volumes` should be enough for VictoriaMetrics.

About AWS EBS disk types you can read in the documentation
[https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ebs-volume-types.html](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ebs-volume-types.html).

### Cannot read stream body

**When it can occur:**

* Big the response size
* Small `promscrape.maxScrapeSize`
* You can face with error message
  
  ```bash
  Cannot read stream body in 1 seconds: the response from "https://x.x.x.x:6443/metrics" exceeds -promscrape.maxScrapeSize=16777216; either reduce the response size for the target or increase -promscrape.maxScrapeSize
  ```

**Root-cause:**

VMAgent has a small default value for `promscrape.maxScrapeSize` parameter. It is `16777216` bytes (or `16MB`).

```shell
-promscrape.maxScrapeSize size
     The maximum size of scrape response in bytes to process from Prometheus targets. Bigger responses are rejected
     Supports the following optional suffixes for size values: KB, MB, GB, TB, KiB, MiB, GiB, TiB (default 16777216)
```

**How to fix it:**

Need to redeploy and specify increased max scrape size in deploy parameters.

```yaml
victoriametrics:
  vmAgent:
    extraArgs:
      promscrape.maxScrapeSize: 256MiB
```

### VictoriaMetrics pods continuously restart

**When it can occur:**

* VMAgent or VMAlertmanager pods may be continuously restarting
* In logs you can see error messages like:

  ```bash
  cannot read "/etc/vmagent/config_out/vmagent.env.yaml": cannot parse Prometheus config from "/etc/vmagent/config_out/vmagent.env.yaml": cannot parse `scrape_config`: cannot parse auth config for `job_name` "serviceScrape/kafka/kafka-service-monitor-jmx-exporter/0": missing `username` in `basic_auth` section
  ```

  or

  ```bash
  ts=2024-03-26T06:58:37.006Z caller=coordinator.go:118 level=error component=configuration msg="Loading configuration file failed" file=/etc/alertmanager/config/alertmanager.yaml err="missing to address in email config"
  ```

**Root-cause:**

The configuration for the VictoriaMetrics pod is incorrect and these pods can't start with an incorrect configuration.
VMAgent, VMAlert and VMAlertmanager are components most susceptible to this issue because users can provide custom
settings.

Also, the reason and part of the incorrect configuration are printed in the log.

**How to fix it:**

Need to read the error message carefully and search for information about which part of the config is incorrect.

For example, for the message:

```bash
cannot parse auth config for `job_name` "serviceScrape/kafka/kafka-service-monitor-jmx-exporter/0": missing `username` in `basic_auth` section
```

this error means that the problem:

* In `ServiceMonitor` (in case of using prometheus-operator CRs) or in `VMServiceScrape`
  (in case of using victoriametrics-operator CRs)
* In the namespace with the name `kafka`
* In the CR with the name `kafka-service-monitor-jmx-exporter`
* The section `basicAuth` has no `username` field

In this case, need to find this CR and fix incorrect configurations in it.

Another example, of the message:

```bash
msg="Loading configuration file failed" file=/etc/alertmanager/config/alertmanager.yaml err="missing to address in email config"
```

this error means that the problem:

* Either in `AlertmanagerConfig` CRs (in case of using prometheus-operator CRs) or in `VMAlertmanagerConfig` CRs
  (in case of using victoriametrics-operator CRs)
* Or in the Secret with additional configs (that allow specifying the part of config in the native format)
* The email address was not specified in the email notification channel config

In this case, need to find where was specified email notification channel config was and add the lost email address.

### VictoriaMetrics config secret does not contain inhibit rules

**When it can occur:**

* VM configuration secret doesn't contain inhibit_rules

**Root-cause:**

The AlertmanagerConfig is incorrect and inhibit rules can't be added into VM configuration secret.
Route and receivers are not mandatory fields for Prometheus AlertmanagerConfig but required for VMAlertmanagerConfig.

Alertmanager incorrect config in case of using VictoriaMetrics:

```yaml
apiVersion: monitoring.coreos.com/v1alpha1
kind: AlertmanagerConfig
metadata:
  name: slack-config-example
  labels:
    app.kubernetes.io/component: monitoring  # Mandatory label
spec:
  inhibitRules:
    - equal:
        - namespace
        - service
      sourceMatch:
        - matchType: '='
          name: alertname
          value: StreamingPlatformIsDownscaledAlert
        - matchType: '='
          name: severity
          value: information
      targetMatch:
        - matchType: '='
          name: service
          value: streaming-service-streaming-service
        - matchType: '=~'
          name: severity
          value: (high|warning|critical)
  receivers: []
```

**How to fix it:**

As VictoriaMetrics operator convert Prometheus AlertmanagerConfig into VMAlertmanagerConfig,
route and receivers fields become mandatory for Prometheus AlertmanagerConfig too.

Alertmanager correct config in case of using VictoriaMetrics:

```yaml
apiVersion: monitoring.coreos.com/v1alpha1
kind: AlertmanagerConfig
metadata:
  name: slack-config-example
  labels:
    app.kubernetes.io/component: monitoring  # Mandatory label
spec:
  inhibitRules:
    - equal:
        - namespace
        - service
      sourceMatch:
        - matchType: '='
          name: alertname
          value: StreamingPlatformIsDownscaledAlert
        - matchType: '='
          name: severity
          value: information
      targetMatch:
        - matchType: '='
          name: service
          value: streaming-service-streaming-service
        - matchType: '=~'
          name: severity
          value: (high|warning|critical)
  receivers:
    - name: base
  route:
    receiver: base
```

For more information about AlertmanagerConfig, refer to the [Configuration documentation](configuration.md#alertmanagerconfig).

## Grafana

### Grafana Data Source not found

The `grafana-operator v4.x` has a problem ([https://github.com/grafana/grafana-operator/issues/652](https://github.com/grafana/grafana-operator/issues/652))
when the Prometheus datasource is not imported until the controller-manager
is restarted. The problem is that both resources are created at the same time and the Grafana pod is started before
the operator updates the config map for data sources.

However, you can update the GrafanaDataSource object (for example, add a new label) so that grafana-deployment is
recreated and datasource is imported.

### Grafana-operator stop discover new dashboards or locked

**Description:**

The `grafana-operator` can be run in 2 or more replicas and to synchronize actions between some replicas
it has a special mechanism to select the master replica. All other replicas should just wait until
the master releases the lock.

The `grafana-operator` uses the ConfigMap to save which replica currently has the master role
and lock other replicas.

Sometimes during the unexpected restart of `grafana-operator` it can lose this ConfigMap and all
other replicas (if we have only one replica of the operator) stay locked.

**Stack trace(s):**

These logs in the `grafana-operator` indicate the the current replicas is locked and won't execute
any actions:

```bash
1.7067316980451612e+09  INFO  leader  Trying to become the leader.
1.7067316996522064e+09  INFO  leader  Found existing lock {"LockOwner": "grafana-operator-75bc7d8469-pv289"}
1.7067316996705446e+09  INFO  leader  Leader pod has been deleted, waiting for garbage collection to remove the lock.
```

**How to solve:**

To unlock the `grafana-operator` need:

1. Login in Kubernetes using `kubectl` or other ways
2. Find in the namespace with monitoring the ConfigMap with the name `grafana-operator-lock`,
   for example using the command:

    ```bash
    kubectl -n <monitoring_namespace> get configmaps grafana-operator-lock
    ```

    example of output:

    ```bash
    ❯ kubectl -n monitoring get configmaps grafana-operator-lock
    NAME                    DATA   AGE
    grafana-operator-lock   0      15d
    ```

3. Remove this ConfigMap, for example using the command:

    ```bash
    kubectl -n <monitoring_namespace> delete configmaps grafana-operator-lock
    ```

4. Restart the grafana-operator, for example using the command:

    ```bash
    kubectl -n <monitoring_namespace> delete pod <grafana_operator_pod_name>
    ```

## Exporters

### NodeExporter

#### Not all node-exporter Endpoints have UP Status in OpenShift

Iptables service specifies a port's range in OpenShift. By default, the range is from `30000` to `32999`.
The `NodeExporter` runs on the 9900 port, which is outside the specified default range.

So the iptables service rejects all the requests to this port.

**Warning**: The port for NodeExporter should be opened on each virtual machine. For OpenStack, the port should be in
the Security Group.

To fix this problem, you can edit the list of opened ports on each virtual machine in OpenShift (except balancer nodes).

1. Login to the virtual machine by SSH.
2. Open the file `etc/sysconfig/iptables`:

    ```bash
    vi /etc/sysconfig/iptables
    ```

3. Find the word "COMMIT" and add the following line before it:

    ```bash
    -A OS_FIREWALL_ALLOW -p tcp -m state --state NEW -m tcp --dport 9900 -j ACCEPT
    ```

4. Save the file and exit.
5. Restart the iptables service:

    ```bash
    systemctl restart iptables
    ```

`monitoring-operator` will use the built-in SecurityContextConstraint for node-exporter that is provided
by Openshift.

### Network-latency-exporter

#### Network Latency Dashboards do not Contain Data, but Exporter Pods are Healthy

This may be due to insufficient rights for the exporter. Make sure that RBAC resources were created correctly.

Network-latency-exporter requires root permissions to work, so make sure that `securityContext` for exporter pods has
the `runAsUser: "0"` parameter.

**OpenShift v4.x only**: For security reasons, giving root user rights is not enough for the correct working
of the exporter in the OpenShift v4.x version. If you want to use the exporter in this case, you can set the
`.Values.networkLatencyExporter.rbac.privileged` parameter to `true` during the deployment.

### Kubelet

#### VMSingle/Prometheus cannot handle container_start_time_seconds metric with too small a timestamp

**This issue can be fixed on the Kubernetes side since v1.29**: <https://github.com/kubernetes/kubernetes/pull/120518>

Example of warning log message from VMSingle:

```bash
2024-01-08T09:42:33.094Z    warn    VictoriaMetrics/lib/storage/storage.go:1733    warn occurred during rows addition:
cannot insert row with too small timestamp 1696452826250 outside the retention; minimum allowed timestamp is
1704447752000; probably you need updating -retentionPeriod command-line flag; metricName:
container_start_time_seconds{...}
```

Kubelet Service Monitor has 3 endpoints with different paths (`/metrics/resource`, `/metrics/cadvisor` and default path)
that have different sets of metrics, but `container_start_time_seconds` is present on several endpoints. However, only
metrics from `/metrics/resource` endpoint may cause such a problem.

It's a known issue in the Monitoring, so the default set of parameters already has relabeling and
metricRelabelings rules to drop the metric only from the particular endpoint. Therefore, you can face this issue if
you override the default configuration of Kubernetes monitors, manually or during deployment.

You can add the following relabeling and metricRelabeling rules for kubelet service monitor to `kubernetesMonitors`
section of configuration to fix this issue:

```yaml
kubernetesMonitors:
  ...
  kubeletServiceMonitor:
    ...
    metricRelabelings:
      # We drop the container_start_time_seconds metric from /metrics/resource endpoint, because we faced problems with
      # this particular metric before, and it can be found on other kubelet endpoints
      - action: drop
        regex: container_start_time_seconds;\/metrics\/resource
        separator: ;
        sourceLabels: ['__name__', 'metrics_path']
      # We drop entire label set "metrics_path", because it was created only for the drop of metric above
      - action: labeldrop
        regex: metrics_path
    ...
    relabelings:
      # This relabeling is used to drop metric container_start_time_seconds in the "metricRelabelings" section, because:
      # 1. "metricRelabelings" happens later than "relabelings"
      # 2. __metrics_path__ drops after "relabelings" step, so it's not in the "metricRelabelings" section
      # 3. we have to use this value, because we must filter out the metric container_start_time_seconds only from
      # /metrics/resource endpoint, but there are no ways to configure relabelings for each endpoint separately
      - action: replace
        regex: (\/metrics\/resource)
        replacement: $1
        sourceLabels:
          - __metrics_path__
        targetLabel: metrics_path
```

Also, you can add the relabeling above manually directly to the service monitor, but only in cases when the monitor is
NOT managed by monitoring-operator. Otherwise, more likely these changes will be overridden during the next
reconciliation cycle.

## Integrations

### Graphite-remote-adapter

This section describes various issues that can occur with the graphite-remote-adapter component and how to troubleshoot
them.

#### Graphite-remote-adapter Regularly Restarts with OOM

In some cases, the graphite-remote-adapter pod regularly restarts with an OOM error in
the status.

This issue can occur in different ways:

* The adaptor works normally in one environment but restarts in another environment.
* The adapter works without any issues for a long time. However, after a restart, it begins to restart with OOM.
* The adapter was scaled to 0 replicas for a long time, and after it was scaled up, it started to restart with OOM.

This issue occurs due to the nature of the working of the Prometheus RemoteWrite protocol.

RemoteWrite is a protocol based on `protobuf` for sending metrics from Prometheus to remote storage pr to remote
integration. Each remote write destination starts a queue, which reads from the write-ahead log (WAL), writes the
samples into a memory queue owned by a shard, which then sends a request to the configured endpoint.

The flow of data looks like this:

```bash
      |-->  queue (shard_1)   --> remote endpoint
WAL --|-->  queue (shard_...) --> remote endpoint
      |-->  queue (shard_n)   --> remote endpoint
```

When one shard backs up and fills its queue, Prometheus blocks reading from the WAL into any shards.
The failures are retried without loss of data unless the remote endpoint remains down for more than 2 hours.
After 2 hours, the WAL is compacted and the data that has not been sent is lost.

In other words, if you configure RemoteWrite to graphite-remote-adapter, but for any reason, the adapter was stopped,
Prometheus sends the data for the last 2 hours when it comes back.

It leads to a lot of points that are sent to graphite-remote-adapter, and it has no resources to process
it immediately. So a lot of points are written in the cache.
The cache size (which usually specifies the counts of records) may read the memory limit
and the pods are killed with OOMKiller.

To solve this issue, do either of the following:

* Increase CPU to 1 vCPU (in some cases to 2 vCPU) and memory to 2-4 Gb (recommended way).
* Decrease the cache size (there is no information on how much the cache size should be decreased).

To increase the allocated resources, redeploy graphite-remote-adapter using the following parameters:

```yaml
graphite_remote_adapter:
  resources:
    limits:
      cpu: 2000m
      memory: 4000Mi
    requests:
      cpu: 1000m
      memory: 2000Mi
```

Alternatively, you can manually change the deployment of graphite-remote-adapter.

To decrease the cache size, redeploy the adapter with the following parameters:

```yaml
graphite_remote_adapter:
  additionalGraphiteConfig:
    write:
      timeout: 5m
    graphite:
      write:
        enable_paths_cache: true
        paths_cache_ttl: 7m
        paths_cache_purge_interval: 8m
```

Alternatively, you can add/change these parameters into ConfigMap with the name, `graphite-config`. For example:

```yaml
write:
  timeout: 5m
graphite:
  write:
    enable_paths_cache: true
    paths_cache_ttl: 7m
    paths_cache_purge_interval: 8m
```

Remember to restart the Graphite pod.

#### Graphite-remote-adapter send metrics with delay

**Description:**

In some cases, the `graphite-remote-adapter` pod can send metrics to Graphite/carbon-c-relay/other graphite-compatible
target with delay.

For example, the delay to send metrics can be from 5m to 1h.

It is not always an issue with `graphite-remote-adapter` and tuning the adapter
parameters and increase replica count or resources, the target, and its settings may not help.

There are a lot of cases when `carbon-c-relay` that should receive metrics in the Graphite format can't process
the entire metrics stream that is generated by `graphite-remote-adapter`. As a result, the adapter store does not send data
to the cache.

But in case of lot of metrics sent through the `graphite-remote-adapter`, it may require a lot of resources.
Without these resources, adapter can send converted metrics with delay.

**How to solve:**

First of all, check `carbon-c-relay` and `Graphite` to make sure that they can process the entire metrics stream
that is generated by the `graphite-remote-adapter`.

Second, find or turn on metric scrape for `graphite-remote-adapter` to understand what exactly happened with it.
If metrics weren't enabled you can enable it using the deployment parameters:

```yaml
graphite_remote_adapter:
  servicemonitor:
    install: true
  grafanaDashboard: true
```

**Note:** These parameters by default set as `true` and monitoring for `graphite-remote-adapter`
should be enabled by default.

Third, if by metrics you see, that the `graphite-remote-adapter` does not have enough resources,
in this case better to scale it to some replicas.

There are two ways how you can scale the `graphite-remote-adapter` replicas count:

* Update deployment parameters and execute `Rolling Update`, parameters:

  ```yaml
  graphite_remote_adapter:
    replicas: 2
  ```

* Manually using `kubectl` CLI or other ways to scale replicas in deployment, for example, using the CLI:

  ```bash
  kubectl --namespace=monitoring scale deployment --replicas=2 --selector="app.kubernetes.io/component=graphite-remote-adapter"
  ```

Fourth, if by metrics you see, that the `graphite-remote-adapter` has a high throttling rate, you can manually
add a new ENV in the Deployment:

```yaml
env:
  - name: GOMAXPROCS
    value: "<core_count>"
```

where:

* `<core_count>` should align with a CPU limit, for example, if the CPU limit set a `2` or `2000m`,
  the `GOMAXPROCS` should be `2`

**Note:** Since release `0.70.0` the graphite-remote-adapter can set `GOMAXPROCS` automatically based on the set
pod's limits during the start. So if you are using this or a higher version, you can skip this recommendation.

## Prometheus-adapter

### ServiceNotFound for `v1.custom.metrics.k8s.io` API Service

When `prometheus-adapter` is not available (the pod was removed, or the pod continuously restarted),
Kubernetes API Server could not handle some requests.

Signs of a problem:

* Kubernetes namespace can't be removed (it is stuck after the delete command)
* Errors related to API in the Kubernetes API Server

Or logs that look like:

```bash
"Error while check hasAPI","error":"unable to retrieve the complete list of server APIs: custom.metrics.k8s.io/v1: the server is currently unable to handle the request"
```

To check that `prometheus-adapter` is not available, execute the following command:

```bash
❯ kubectl get apiservices
NAME                       SERVICE                         AVAILABLE                 AGE
...
v1.custom.metrics.k8s.io   prometheus/prometheus-adapter   False (ServiceNotFound)   12s
...
```

If you see that the `AVAILABLE` column contains `False`, it means that the adapter in the namespace `prometheus`
and with the name `prometheus-adapter` couldn't handle requests.

To fix this issue, you have two options:

* Restore the `prometheus-adapter` pod, for example, run it or add more resources.
* Remove the early registered `APIService`.

To remove the early registered `APIService`, execute the following command:

```bash
kubectl delete apiservice v1.custom.metrics.k8s.io
```

**Note:** Old Kubernetes versions can use `v1beta1` API version and you have to work
with `v1beta1.custom.metrics.k8s.io`.

### FailedDiscoveryCheck for `v1.custom.metrics.k8s.io` API Service

When `prometheus-adapter` is available, but didn't discover no one rules, it won't run handler for registered API.
It means that Kubernetes API Server will can't process HPA requests for custom metrics.

Signs of a problem:

* HPA by custom metrics doesn't work
* Horizontal Pod Autoscaler (HPA) objects generate event with error like:

    ```bash
    Warning  FailedGetPodsMetric  <invalid> (x63342 over 10d)  horizontal-pod-autoscaler
      unable to get metric tm_busyness: unable to fetch metrics from custom metrics API:
        the server could not find the metric <metric_name> for pods
    ```

To check th `prometheus-adapter` status and availability you can use a command:

```bash
❯ kubectl get apiservices
NAME                       SERVICE                         AVAILABLE                      AGE
...
v1.custom.metrics.k8s.io   monitoring/prometheus-adapter   False (FailedDiscoveryCheck)   85m
...
```

Usually, this problem means one of two:

* either somebody specifies incorrect label selectors to discover Custom Resources with metrics declaration
* or you have no one Custom Resources with declared custom metrics

To fix the first issue you have two options:

* either remove conditions from the label selector for prometheus-adapter:

    ```yaml
    prometheusAdapter:
      customScaleMetricRulesSelector: []   # set empty labelSelector
    ```

* or add expected labels in `CustomScaleMetricRule` CR:
  * for example, you have the following label selector:

    ```yaml
    prometheusAdapter:
      customScaleMetricRulesSelector:
        - matchExpressions:
          - key: app.kubernetes.io/component
            operator: In
            values:
              - monitoring
    ```
  
  * in this case, you have to add label `app.kubernetes.io/component: monitoring` to CR:

    ```yaml
    apiVersion: monitoring.qubership.org/v1alpha1
    kind: CustomScaleMetricRule
    metadata:
      name: you-custom-rules
      labels:
        app.kubernetes.io/component: monitoring   # add expected label from label selector
    ...
    ```

To solve the second issue need to add at least one `CustomScaleMetricRule` CR with at least one rule.
The `prometheus-adapter` will enable the handler and will process API requests for `v1.custom.metrics.k8s.io`
only in case when it has at least one rule for `v1.custom.metrics.k8s.io`.

### prometheus-adapter Pod is Down or Restarting

When `prometheus-adapter` is down or the `prometheus-adapter` pod is continuously restarted, it can
have some root cases:

* `prometheus-adapter` did not have enough resources to work and was killed by OOMKiller.
* `Prometheus URL` contains an incorrect link to Prometheus.
* Prometheus, specified as a source for metrics, is not available.

To solve the issue when `prometheus-adapter` does not have enough resources, increase CPU
or memory limits for the adapter. To change the limits, change the deployment parameters or change the value
directly in the `Prometheus Adapter` CR. For example:

```yaml
prometheusAdapter:
  resources:
    requests:
      cpu: 500m
      memory: 1Gi
    limits:
      cpu: 1000m
      memory: 2Gi
```

After redeploying or changing the CR, the operator redeploys the `prometheus-adapter` pod.

**Note**: The amount of required resources depends on the count of custom metrics and the frequency of
Prometheus metrics scraping.

To fix an incorrect link on Prometheus, change the deployment parameters or change the value
directly in the `Prometheus Adapter` CR. For example:

```yaml
prometheusAdapter:
  prometheusUrl: http://prometheus.monitoring.svc:9090
```

**Note**: Use a full URL to Prometheus with schema and port.

To fix issues with Prometheus, refer to the sections that describe how to fix Prometheus.
This document provides information about troubleshooting Monitoring deployment and its configurations.
It also provides the configurations of monitoring services.
